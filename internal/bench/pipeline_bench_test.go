// Pipeline Invariant & Chaos Bench
//
// This is NOT a test suite. It is a structural integrity scanner.
//
// It enforces three core invariants across every object type:
//   1. IF it can be created  → it MUST be saveable
//   2. IF it can be saved    → it MUST be loadable
//   3. IF it can be loaded   → it MUST match the saved object
//
// BREAK-IT MODE (set HOWL_BREAK_IT=1):
//   Injects faults at every pipeline stage and maps exactly where
//   failures propagate, hide, and resurface. Produces a failure map,
//   not a stack trace.
//
// Run clean:    go test -v ./tests/bench/... -timeout 120s
// Run chaos:    HOWL_BREAK_IT=1 go test -v ./tests/bench/... -timeout 120s
// Run race:     go test -v -race ./tests/bench/... -timeout 120s
// Run all:      HOWL_BREAK_IT=1 go test -v -race ./tests/bench/... -timeout 120s

package bench_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"howl-chat/internal/backend/chat"
)

// ---------------------------------------------------------------------------
// Bench configuration
// ---------------------------------------------------------------------------

// breakItMode enables fault injection. Set HOWL_BREAK_IT=1 to activate.
var breakItMode = os.Getenv("HOWL_BREAK_IT") == "1"

// faultLog records every injected fault and where it propagated.
type faultEntry struct {
	Stage      string
	Fault      string
	Propagated bool
	HiddenIn   string
	Caught     bool
	CaughtAt   string
}

var (
	faultMu  sync.Mutex
	faultLog []faultEntry
)

func logFault(stage, fault string) {
	faultMu.Lock()
	defer faultMu.Unlock()
	faultLog = append(faultLog, faultEntry{Stage: stage, Fault: fault})
}

func markPropagated(stage, hiddenIn string) {
	faultMu.Lock()
	defer faultMu.Unlock()
	for i := len(faultLog) - 1; i >= 0; i-- {
		if faultLog[i].Stage == stage && !faultLog[i].Propagated {
			faultLog[i].Propagated = true
			faultLog[i].HiddenIn = hiddenIn
			return
		}
	}
}

func markCaught(stage, where string) {
	faultMu.Lock()
	defer faultMu.Unlock()
	for i := len(faultLog) - 1; i >= 0; i-- {
		if faultLog[i].Stage == stage && !faultLog[i].Caught {
			faultLog[i].Caught = true
			faultLog[i].CaughtAt = where
			return
		}
	}
}

func printFailureMap(t *testing.T) {
	t.Helper()
	faultMu.Lock()
	defer faultMu.Unlock()
	if len(faultLog) == 0 {
		return
	}
	t.Log("\n╔══════════════════════════════════════════════════════╗")
	t.Log("║              FAILURE MAP (break-it mode)             ║")
	t.Log("╠══════════════════════════════════════════════════════╣")
	for _, e := range faultLog {
		propagated := "contained"
		if e.Propagated {
			propagated = fmt.Sprintf("PROPAGATED → hid in: %s", e.HiddenIn)
		}
		caught := "UNCAUGHT"
		if e.Caught {
			caught = fmt.Sprintf("caught at: %s", e.CaughtAt)
		}
		t.Logf("║  stage=%-20s fault=%-30s", e.Stage, e.Fault)
		t.Logf("║    propagation: %-40s", propagated)
		t.Logf("║    detection:   %-40s", caught)
		t.Log("║")
	}
	t.Log("╚══════════════════════════════════════════════════════╝")
}

// ---------------------------------------------------------------------------
// Fixture factories — produce canonical test objects
// ---------------------------------------------------------------------------

func makeScenario(id, name string, charCount int) *chat.ScenarioContext {
	sc := &chat.ScenarioContext{}
	sc.ID = id
	sc.Name = name
	sc.Description = "Bench-generated scenario"
	for i := 0; i < charCount; i++ {
		sc.Characters = append(sc.Characters, chat.ScenarioContextCharacter{
			ID:   fmt.Sprintf("char_%s_%d", id, i),
			Name: fmt.Sprintf("Character%d", i),
			Role: "bench role",
		})
	}
	return sc
}

func makeCharacter(id, name string) *chat.CharacterContext {
	cc := &chat.CharacterContext{}
	cc.ID = id
	cc.Name = name
	cc.Personality = "bench personality"
	cc.Backstory = "bench backstory"
	return cc
}

func makeWorld(id, name string) *chat.WorldContext {
	wc := &chat.WorldContext{}
	wc.ID = id
	wc.Name = name
	wc.Description = "bench world"
	return wc
}

func makeChatData(id string, scenario *chat.ScenarioContext, chars []*chat.CharacterContext, world *chat.WorldContext) *chat.ChatData {
	return &chat.ChatData{
		Meta: chat.ChatMeta{
			ID:        id,
			Name:      "Bench Chat " + id,
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		},
		World:      world,
		Scenario:   scenario,
		Characters: chars,
		Lorebooks:  []interface{}{},
		Memories:   []interface{}{},
		History:    []chat.ChatHistory{},
	}
}

func newBenchPM(t *testing.T) (*chat.PersistenceManager, string) {
	t.Helper()
	dir := t.TempDir()
	pm, err := chat.NewPersistenceManagerWithDir(dir)
	if err != nil {
		t.Fatalf("NewPersistenceManagerWithDir: %v", err)
	}
	return pm, dir
}

// ---------------------------------------------------------------------------
// Core invariant: create → save → load → match
// ---------------------------------------------------------------------------

func assertRoundTrip(t *testing.T, pm *chat.PersistenceManager, in *chat.ChatData) *chat.ChatData {
	t.Helper()

	// SAVE
	if err := pm.SaveChat(in); err != nil {
		t.Fatalf("[INVARIANT BROKEN] save failed for chat %q: %v\n  → save pipeline has a hole", in.Meta.ID, err)
	}

	// LOAD
	out, err := pm.LoadChat(in.Meta.ID)
	if err != nil {
		t.Fatalf("[INVARIANT BROKEN] load failed for chat %q: %v\n  → load pipeline cannot reconstruct what save wrote", in.Meta.ID, err)
	}

	// MATCH — scenario
	if in.Scenario != nil {
		if out.Scenario == nil {
			t.Errorf("[INVARIANT BROKEN] scenario lost in transit\n  saved name=%q → loaded=nil\n  → scenario pipeline drops the object", in.Scenario.Name)
		} else {
			if out.Scenario.Name != in.Scenario.Name {
				t.Errorf("[INVARIANT BROKEN] scenario name mutated\n  saved=%q → loaded=%q\n  → name field is being overwritten or not serialised", in.Scenario.Name, out.Scenario.Name)
			}
			if len(out.Scenario.Characters) != len(in.Scenario.Characters) {
				t.Errorf("[INVARIANT BROKEN] scenario characters count mutated\n  saved=%d → loaded=%d\n  → character array is not being serialised; check ScenarioContext.Characters json tag",
					len(in.Scenario.Characters), len(out.Scenario.Characters))
			} else {
				for i, c := range in.Scenario.Characters {
					if out.Scenario.Characters[i].Name != c.Name {
						t.Errorf("[INVARIANT BROKEN] scenario character[%d] name mutated: saved=%q loaded=%q", i, c.Name, out.Scenario.Characters[i].Name)
					}
				}
			}
		}
	}

	// MATCH — characters
	if len(out.Characters) != len(in.Characters) {
		t.Errorf("[INVARIANT BROKEN] character count mutated\n  saved=%d → loaded=%d\n  → characters pipeline drops entries", len(in.Characters), len(out.Characters))
	} else {
		inNames := map[string]bool{}
		for _, c := range in.Characters {
			inNames[c.Name] = true
		}
		for _, c := range out.Characters {
			if !inNames[c.Name] {
				t.Errorf("[INVARIANT BROKEN] unexpected character after load: %q", c.Name)
			}
		}
	}

	// MATCH — world
	if in.World != nil {
		if out.World == nil {
			t.Errorf("[INVARIANT BROKEN] world lost in transit: saved=%q loaded=nil", in.World.Name)
		} else if out.World.Name != in.World.Name {
			t.Errorf("[INVARIANT BROKEN] world name mutated: saved=%q loaded=%q", in.World.Name, out.World.Name)
		}
	}

	return out
}

// ---------------------------------------------------------------------------
// T1: Happy path invariant — every object type, canonical data
// ---------------------------------------------------------------------------

func TestInvariant_HappyPath_ScenarioWithTwoCharacters(t *testing.T) {
	pm, _ := newBenchPM(t)

	sc := makeScenario("sc-bench-001", "Stone Wolf Systems", 2)
	sc.Characters[0].Name = "Emily"
	sc.Characters[1].Name = "Bri"

	chars := []*chat.CharacterContext{
		makeCharacter("char_emily", "Emily"),
		makeCharacter("char_bri", "Bri"),
	}

	data := makeChatData("chat-bench-001", sc, chars, makeWorld("world-001", "SWS HQ"))
	assertRoundTrip(t, pm, data)
}

func TestInvariant_HappyPath_ScenarioNameFromTitle(t *testing.T) {
	m := map[string]interface{}{
		"id":    "sc-title-001",
		"title": "Title Only Scenario",
		"characters": []interface{}{
			map[string]interface{}{"name": "Alice", "role": "hero"},
		},
	}
	sc, err := chat.ParseScenarioContextFromMap(m)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if sc.Name != "Title Only Scenario" {
		t.Errorf("[INVARIANT BROKEN] title→name mapping: got %q", sc.Name)
	}
	if len(sc.Characters) != 1 {
		t.Errorf("[INVARIANT BROKEN] characters not preserved through parse: got %d", len(sc.Characters))
	}
}

// ---------------------------------------------------------------------------
// T2: Stress — rapid create/save/load/delete cycles
// ---------------------------------------------------------------------------

func TestStress_RapidCreateSaveLoad(t *testing.T) {
	pm, _ := newBenchPM(t)
	const cycles = 50

	for i := 0; i < cycles; i++ {
		id := fmt.Sprintf("stress-chat-%04d", i)
		sc := makeScenario(fmt.Sprintf("sc-%04d", i), fmt.Sprintf("Scenario %d", i), 2)
		chars := []*chat.CharacterContext{
			makeCharacter(fmt.Sprintf("char-a-%04d", i), fmt.Sprintf("Alpha%d", i)),
			makeCharacter(fmt.Sprintf("char-b-%04d", i), fmt.Sprintf("Beta%d", i)),
		}
		data := makeChatData(id, sc, chars, nil)
		out := assertRoundTrip(t, pm, data)
		if out == nil {
			t.Errorf("[STRESS] cycle %d: round-trip returned nil", i)
		}
	}
}

// ---------------------------------------------------------------------------
// T3: Concurrent saves — race detector will catch shared-state corruption
// ---------------------------------------------------------------------------

func TestStress_ConcurrentSaves(t *testing.T) {
	const workers = 8
	var wg sync.WaitGroup
	errors := make(chan string, workers*2)

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			pm, _ := newBenchPM(t)
			for i := 0; i < 10; i++ {
				id := fmt.Sprintf("concurrent-w%d-c%d", workerID, i)
				sc := makeScenario(id+"-sc", fmt.Sprintf("ConcurrentScenario-%d-%d", workerID, i), 2)
				data := makeChatData(id, sc, []*chat.CharacterContext{makeCharacter(id+"-char", "ConcurrentChar")}, nil)
				if err := pm.SaveChat(data); err != nil {
					errors <- fmt.Sprintf("worker %d save failed: %v", workerID, err)
					return
				}
				out, err := pm.LoadChat(id)
				if err != nil {
					errors <- fmt.Sprintf("worker %d load failed: %v", workerID, err)
					return
				}
				if out.Scenario == nil || out.Scenario.Name != sc.Name {
					errors <- fmt.Sprintf("worker %d invariant broken: scenario name mutated under concurrency", workerID)
				}
			}
		}(w)
	}
	wg.Wait()
	close(errors)
	for msg := range errors {
		t.Error(msg)
	}
}

// ---------------------------------------------------------------------------
// T4: Missing field resilience — objects with only partial data
// ---------------------------------------------------------------------------

func TestResilience_MissingFields(t *testing.T) {
	cases := []struct {
		name string
		m    map[string]interface{}
		want string
	}{
		{"no name no title", map[string]interface{}{"id": "sc-001"}, ""},
		{"empty name", map[string]interface{}{"id": "sc-002", "name": ""}, ""},
		{"title only", map[string]interface{}{"id": "sc-003", "title": "TitleOnly"}, "TitleOnly"},
		{"name and empty title", map[string]interface{}{"id": "sc-004", "name": "N", "title": ""}, "N"},
		{"no id", map[string]interface{}{"title": "NoID"}, "NoID"},
		{"nil characters", map[string]interface{}{"id": "sc-005", "title": "NilChars", "characters": nil}, "NilChars"},
		{"characters missing name", map[string]interface{}{"id": "sc-006", "title": "BlankChar",
			"characters": []interface{}{map[string]interface{}{"role": "ghost"}}}, "BlankChar"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sc, err := chat.ParseScenarioContextFromMap(tc.m)
			if err != nil {
				t.Fatalf("parse must not error on partial data: %v", err)
			}
			if sc.Name != tc.want {
				t.Errorf("Name=%q, want %q", sc.Name, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// T5: Corrupted disk data — what happens when files on disk are broken
// ---------------------------------------------------------------------------

func TestResilience_CorruptedDiskFiles(t *testing.T) {
	corruptionCases := []struct {
		file    string
		content string
		desc    string
	}{
		{"scenario.json", `{broken json`, "truncated JSON"},
		{"scenario.json", `null`, "null scenario — should return nil scenario, not error"},
		{"scenario.json", `{"name":"","characters":[]}`, "empty name — degraded but loadable"},
		{"characters.json", `[`, "truncated characters array"},
		{"characters.json", `null`, "null characters — should return empty array"},
		{"characters.json", `[{"name":""}]`, "character with empty name — degraded but loadable"},
		{"meta.json", `{}`, "empty meta — ID must be recovered from directory name"},
	}

	for _, tc := range corruptionCases {
		t.Run(tc.desc+"_"+tc.file, func(t *testing.T) {
			pm, dir := newBenchPM(t)

			// Create a valid chat first
			data := makeChatData("corrupt-test", makeScenario("sc", "Valid", 1), []*chat.CharacterContext{makeCharacter("c1", "Valid")}, nil)
			if err := pm.SaveChat(data); err != nil {
				t.Fatalf("setup save failed: %v", err)
			}

			// Corrupt the target file
			target := filepath.Join(dir, "chats", "corrupt-test", tc.file)
			if err := os.WriteFile(target, []byte(tc.content), 0644); err != nil {
				t.Fatalf("could not write corrupt file: %v", err)
			}

			if breakItMode {
				logFault("disk:"+tc.file, tc.desc)
			}

			// Attempt load — must not panic, must return coherent error or degraded-but-valid data
			loaded, err := pm.LoadChat("corrupt-test")
			if err != nil {
				if breakItMode {
					markCaught("disk:"+tc.file, "LoadChat returned error")
				}
				t.Logf("  → load correctly returned error for corrupt %s (%s): %v", tc.file, tc.desc, err)
				return
			}

			// If it didn't error, the data must at least be non-nil and internally consistent
			if loaded == nil {
				t.Errorf("[INVARIANT BROKEN] LoadChat returned nil,nil for corrupt %s — silent failure", tc.file)
				return
			}

			// INVARIANT: meta.ID must always be set (directory name recovery must work)
			if loaded.Meta.ID == "" {
				t.Errorf("[INVARIANT BROKEN] corrupt %s: loaded.Meta.ID is empty — chat cannot be found again", tc.file)
			}

			// INVARIANT: Characters must never be nil (empty is fine)
			if loaded.Characters == nil {
				t.Errorf("[INVARIANT BROKEN] corrupt %s: loaded.Characters is nil — should be empty slice", tc.file)
			}

			if breakItMode {
				markPropagated("disk:"+tc.file, "LoadChat returned degraded data without error")
			}
			t.Logf("  → degraded load for corrupt %s (%s): meta.ID=%q scenario=%v chars=%d",
				tc.file, tc.desc, loaded.Meta.ID, loaded.Scenario != nil, len(loaded.Characters))
		})
	}
}

// ---------------------------------------------------------------------------
// T6: Break-it mode — fault injection at every pipeline stage
// ---------------------------------------------------------------------------

func TestBreakIt_ParsePipeline(t *testing.T) {
	if !breakItMode {
		t.Skip("break-it mode disabled — set HOWL_BREAK_IT=1 to run")
	}

	faults := []struct {
		name  string
		input map[string]interface{}
		check func(sc chat.ScenarioContext) string // returns "" if ok, problem description otherwise
	}{
		{
			name:  "unicode bomb in name",
			input: map[string]interface{}{"id": "b1", "title": strings.Repeat("𝕳𝖔𝖜𝖑", 1000)},
			check: func(sc chat.ScenarioContext) string {
				if sc.Name == "" {
					return "name lost after unicode bomb"
				}
				return ""
			},
		},
		{
			name:  "zero-width space in title",
			input: map[string]interface{}{"id": "b2", "title": "Valid\u200BName"},
			check: func(sc chat.ScenarioContext) string {
				if sc.Name == "" {
					return "name empty after zero-width space"
				}
				return ""
			},
		},
		{
			name:  "null character in name",
			input: map[string]interface{}{"id": "b3", "name": "Name\x00WithNull"},
			check: func(sc chat.ScenarioContext) string { return "" }, // must not panic
		},
		{
			name: "10000 characters in scenario",
			input: func() map[string]interface{} {
				chars := make([]interface{}, 10000)
				for i := range chars {
					chars[i] = map[string]interface{}{"name": fmt.Sprintf("Char%d", i), "role": "bench"}
				}
				return map[string]interface{}{"id": "b4", "title": "MassScenario", "characters": chars}
			}(),
			check: func(sc chat.ScenarioContext) string {
				if len(sc.Characters) != 10000 {
					return fmt.Sprintf("lost characters: got %d, want 10000", len(sc.Characters))
				}
				return ""
			},
		},
		{
			name: "deeply nested image object",
			input: map[string]interface{}{"id": "b5", "title": "DeepImg", "image": map[string]interface{}{
				"data": strings.Repeat("A", 500_000), // 500KB fake base64
				"name": "huge.png",
			}},
			check: func(sc chat.ScenarioContext) string { return "" }, // must not panic or OOM
		},
		{
			name:  "characters key is wrong type (string instead of array)",
			input: map[string]interface{}{"id": "b6", "title": "WrongType", "characters": "not-an-array"},
			check: func(sc chat.ScenarioContext) string {
				// should gracefully produce 0 characters, not panic
				return ""
			},
		},
		{
			name: "random garbage map",
			input: func() map[string]interface{} {
				m := map[string]interface{}{}
				r := rand.New(rand.NewSource(42))
				for i := 0; i < 50; i++ {
					key := fmt.Sprintf("key%d", r.Intn(100))
					m[key] = fmt.Sprintf("val%d", r.Intn(1000))
				}
				return m
			}(),
			check: func(sc chat.ScenarioContext) string { return "" }, // must not panic
		},
	}

	for _, f := range faults {
		t.Run(f.name, func(t *testing.T) {
			logFault("parse:scenario", f.name)

			defer func() {
				if r := recover(); r != nil {
					markPropagated("parse:scenario", "PANIC")
					t.Errorf("[BREAK-IT] PANIC in ParseScenarioContextFromMap: %v\n  fault: %s\n  input: %v", r, f.name, f.input)
				}
			}()

			sc, err := chat.ParseScenarioContextFromMap(f.input)
			if err != nil {
				markCaught("parse:scenario", "ParseScenarioContextFromMap returned error")
				t.Logf("  → parser returned error (acceptable): %v", err)
				return
			}

			if problem := f.check(sc); problem != "" {
				markPropagated("parse:scenario", "silent data loss")
				t.Errorf("[BREAK-IT] silent failure in parse pipeline: %s\n  fault: %s", problem, f.name)
			} else {
				markCaught("parse:scenario", "parse+check passed")
			}
		})
	}
}

func TestBreakIt_PersistencePipeline(t *testing.T) {
	if !breakItMode {
		t.Skip("break-it mode disabled — set HOWL_BREAK_IT=1 to run")
	}

	type pipelineFault struct {
		name   string
		mutate func(dir, chatID string) // corrupt the files AFTER save
		expect string                   // what should break
	}

	faults := []pipelineFault{
		{
			name: "delete scenario.json after save",
			mutate: func(dir, id string) {
				_ = os.Remove(filepath.Join(dir, "chats", id, "scenario.json"))
			},
			expect: "load should return nil scenario (file missing is valid)",
		},
		{
			name: "zero-byte scenario.json",
			mutate: func(dir, id string) {
				_ = os.WriteFile(filepath.Join(dir, "chats", id, "scenario.json"), []byte{}, 0644)
			},
			expect: "load should error or return nil scenario",
		},
		{
			name: "scenario.json with only whitespace",
			mutate: func(dir, id string) {
				_ = os.WriteFile(filepath.Join(dir, "chats", id, "scenario.json"), []byte("   \n\t  "), 0644)
			},
			expect: "load should error — whitespace is not valid JSON",
		},
		{
			name: "characters.json replaced with object instead of array",
			mutate: func(dir, id string) {
				_ = os.WriteFile(filepath.Join(dir, "chats", id, "characters.json"), []byte(`{"name":"NotAnArray"}`), 0644)
			},
			expect: "load should error — characters must be an array",
		},
		{
			name: "meta.json with future timestamp",
			mutate: func(dir, id string) {
				meta := map[string]interface{}{
					"id":        id,
					"name":      "Future",
					"createdAt": "2099-01-01T00:00:00Z",
					"updatedAt": "2099-01-01T00:00:00Z",
				}
				b, _ := json.Marshal(meta)
				_ = os.WriteFile(filepath.Join(dir, "chats", id, "meta.json"), b, 0644)
			},
			expect: "load should succeed — future timestamps are valid",
		},
		{
			name: "directory permissions removed",
			mutate: func(dir, id string) {
				chatDir := filepath.Join(dir, "chats", id)
				_ = os.Chmod(chatDir, 0000)
				t.Cleanup(func() { _ = os.Chmod(chatDir, 0755) })
			},
			expect: "load should error — permission denied",
		},
	}

	for _, f := range faults {
		t.Run(f.name, func(t *testing.T) {
			pm, dir := newBenchPM(t)
			chatID := "break-persist-" + strings.ReplaceAll(f.name, " ", "-")

			sc := makeScenario("sc-break", "BreakScenario", 2)
			data := makeChatData(chatID, sc, []*chat.CharacterContext{makeCharacter("c1", "Char1")}, nil)

			if err := pm.SaveChat(data); err != nil {
				t.Fatalf("setup save failed: %v", err)
			}

			logFault("persist:load", f.name)
			f.mutate(dir, chatID)

			defer func() {
				if r := recover(); r != nil {
					markPropagated("persist:load", "PANIC")
					t.Errorf("[BREAK-IT] PANIC in LoadChat: %v\n  fault: %s\n  expect: %s", r, f.name, f.expect)
				}
			}()

			loaded, err := pm.LoadChat(chatID)
			if err != nil {
				markCaught("persist:load", "LoadChat returned error")
				t.Logf("  ✓ LoadChat correctly errored: %v\n    expect: %s", err, f.expect)
				return
			}

			if loaded == nil {
				markCaught("persist:load", "LoadChat returned nil data")
				t.Logf("  ✓ LoadChat returned nil (expected for: %s)", f.expect)
				return
			}

			markPropagated("persist:load", "LoadChat returned degraded data without error")
			t.Logf("  ⚠ LoadChat returned data without error\n    expect: %s\n    got: scenario=%v chars=%d",
				f.expect, loaded.Scenario != nil, len(loaded.Characters))
		})
	}
}

// ---------------------------------------------------------------------------
// T7: Reload loop — same chat saved and loaded N times, data must not drift
// ---------------------------------------------------------------------------

func TestStress_ReloadLoop(t *testing.T) {
	pm, _ := newBenchPM(t)
	const iterations = 20

	sc := makeScenario("sc-reload", "ReloadScenario", 2)
	sc.Characters[0].Name = "Emily"
	sc.Characters[1].Name = "Bri"

	original := makeChatData("reload-loop", sc,
		[]*chat.CharacterContext{
			makeCharacter("char_emily", "Emily"),
			makeCharacter("char_bri", "Bri"),
		},
		makeWorld("world-reload", "SWS HQ"),
	)

	var prev *chat.ChatData = original
	for i := 0; i < iterations; i++ {
		if err := pm.SaveChat(prev); err != nil {
			t.Fatalf("[RELOAD LOOP] iteration %d: save failed: %v", i, err)
		}
		loaded, err := pm.LoadChat(prev.Meta.ID)
		if err != nil {
			t.Fatalf("[RELOAD LOOP] iteration %d: load failed: %v", i, err)
		}

		// Invariant: scenario name must not drift
		if loaded.Scenario == nil {
			t.Fatalf("[RELOAD LOOP] iteration %d: scenario became nil", i)
		}
		if loaded.Scenario.Name != original.Scenario.Name {
			t.Errorf("[RELOAD LOOP] iteration %d: scenario name drifted: %q → %q", i, original.Scenario.Name, loaded.Scenario.Name)
		}
		if len(loaded.Scenario.Characters) != len(original.Scenario.Characters) {
			t.Errorf("[RELOAD LOOP] iteration %d: scenario characters drifted: %d → %d", i, len(original.Scenario.Characters), len(loaded.Scenario.Characters))
		}
		if len(loaded.Characters) != len(original.Characters) {
			t.Errorf("[RELOAD LOOP] iteration %d: characters count drifted: %d → %d", i, len(original.Characters), len(loaded.Characters))
		}

		prev = loaded
	}
	t.Logf("[RELOAD LOOP] completed %d iterations — no drift detected", iterations)
}

// ---------------------------------------------------------------------------
// T8: Structural completeness — every serialised field must be recoverable
// ---------------------------------------------------------------------------

func TestInvariant_FieldCompleteness(t *testing.T) {
	pm, _ := newBenchPM(t)

	sc := &chat.ScenarioContext{}
	sc.ID = "fc-scenario"
	sc.Name = "Completeness Scenario"
	sc.Description = "All fields populated"
	sc.KeyDetails = []string{"detail1", "detail2"}
	sc.References = []string{"ref1"}
	sc.Characters = []chat.ScenarioContextCharacter{
		{ID: "fc-c1", Name: "Alice", Role: "hero", Avatar: "avatar-data"},
		{ID: "fc-c2", Name: "Bob", Role: "sidekick", Avatar: ""},
	}

	world := &chat.WorldContext{}
	world.ID = "fc-world"
	world.Name = "Completeness World"
	world.Description = "Full world"
	world.KeyFacts = []string{"fact1", "fact2"}

	chars := []*chat.CharacterContext{
		makeCharacter("fc-char1", "Alice"),
		makeCharacter("fc-char2", "Bob"),
	}
	chars[0].Personality = "brave"
	chars[0].Backstory = "hero backstory"
	chars[0].Goals = []string{"goal1", "goal2"}
	chars[1].Personality = "clever"

	data := makeChatData("fc-chat", sc, chars, world)
	out := assertRoundTrip(t, pm, data)

	// Deep-check every populated field
	checks := []struct {
		field string
		got   interface{}
		want  interface{}
	}{
		{"scenario.ID", out.Scenario.ID, sc.ID},
		{"scenario.Name", out.Scenario.Name, sc.Name},
		{"scenario.Description", out.Scenario.Description, sc.Description},
		{"scenario.Characters[0].Name", out.Scenario.Characters[0].Name, "Alice"},
		{"scenario.Characters[0].Role", out.Scenario.Characters[0].Role, "hero"},
		{"scenario.Characters[0].Avatar", out.Scenario.Characters[0].Avatar, "avatar-data"},
		{"world.ID", out.World.ID, world.ID},
		{"world.Name", out.World.Name, world.Name},
		{"world.Description", out.World.Description, world.Description},
		{"characters[0].Personality", out.Characters[0].Personality, "brave"},
		{"characters[0].Backstory", out.Characters[0].Backstory, "hero backstory"},
	}

	for _, c := range checks {
		if !reflect.DeepEqual(c.got, c.want) {
			t.Errorf("[FIELD COMPLETENESS] %s: saved=%v loaded=%v", c.field, c.want, c.got)
		}
	}

	// Check slice lengths independently
	if len(out.Scenario.KeyDetails) != 2 {
		t.Errorf("[FIELD COMPLETENESS] scenario.KeyDetails: saved=2 loaded=%d", len(out.Scenario.KeyDetails))
	}
	if len(out.World.KeyFacts) != 2 {
		t.Errorf("[FIELD COMPLETENESS] world.KeyFacts: saved=2 loaded=%d", len(out.World.KeyFacts))
	}
	if len(out.Characters[0].Goals) != 2 {
		t.Errorf("[FIELD COMPLETENESS] characters[0].Goals: saved=2 loaded=%d", len(out.Characters[0].Goals))
	}
}

// ---------------------------------------------------------------------------
// T9: TestMain — print failure map at end if break-it mode
// ---------------------------------------------------------------------------

func TestMain(m *testing.M) {
	code := m.Run()
	if breakItMode && len(faultLog) > 0 {
		fmt.Println()
		fmt.Println("╔══════════════════════════════════════════════════════╗")
		fmt.Println("║              FAILURE MAP (break-it mode)             ║")
		fmt.Println("╠══════════════════════════════════════════════════════╣")
		faultMu.Lock()
		for _, e := range faultLog {
			propagated := "contained"
			if e.Propagated {
				propagated = fmt.Sprintf("PROPAGATED → hid in: %s", e.HiddenIn)
			}
			caught := "UNCAUGHT"
			if e.Caught {
				caught = fmt.Sprintf("caught at: %s", e.CaughtAt)
			}
			fmt.Printf("║  stage=%-22s fault=%-28s\n", e.Stage, e.Fault[:min(len(e.Fault), 28)])
			fmt.Printf("║    propagation: %-43s\n", propagated)
			fmt.Printf("║    detection:   %-43s\n", caught)
			fmt.Println("║")
		}
		faultMu.Unlock()
		fmt.Println("╚══════════════════════════════════════════════════════╝")
	}
	os.Exit(code)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
