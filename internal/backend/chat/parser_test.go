package chat

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// ---------------------------------------------------------------------------
// ParseScenarioContextFromMap
// ---------------------------------------------------------------------------

func TestParseScenario_TitleMapsToName(t *testing.T) {
	m := map[string]interface{}{
		"id":    "sc-001",
		"title": "Stone Wolf Systems",
	}
	sc, err := ParseScenarioContextFromMap(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sc.Name != "Stone Wolf Systems" {
		t.Errorf("Name=%q, want 'Stone Wolf Systems'", sc.Name)
	}
}

func TestParseScenario_ExplicitNameBeatsTitle(t *testing.T) {
	m := map[string]interface{}{
		"id":    "sc-002",
		"name":  "Real Name",
		"title": "Should Be Ignored",
	}
	sc, err := ParseScenarioContextFromMap(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sc.Name != "Real Name" {
		t.Errorf("Name=%q, want 'Real Name'", sc.Name)
	}
}

func TestParseScenario_TwoCharactersPreserved(t *testing.T) {
	m := map[string]interface{}{
		"id":    "sc-003",
		"title": "Two Char Test",
		"characters": []interface{}{
			map[string]interface{}{"name": "Emily", "role": "daughter"},
			map[string]interface{}{"name": "Bri", "role": "friend"},
		},
	}
	sc, err := ParseScenarioContextFromMap(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sc.Characters) != 2 {
		t.Fatalf("Characters=%d, want 2 — got: %+v", len(sc.Characters), sc.Characters)
	}
	names := map[string]bool{}
	for _, c := range sc.Characters {
		names[c.Name] = true
	}
	for _, want := range []string{"Emily", "Bri"} {
		if !names[want] {
			t.Errorf("character %q missing from scenario characters", want)
		}
	}
}

func TestParseScenario_EmptyCharacters(t *testing.T) {
	m := map[string]interface{}{
		"id":         "sc-004",
		"title":      "Empty",
		"characters": []interface{}{},
	}
	sc, err := ParseScenarioContextFromMap(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sc.Characters) != 0 {
		t.Errorf("Characters=%d, want 0", len(sc.Characters))
	}
}

func TestParseScenario_ImageStrippedFromBody(t *testing.T) {
	m := map[string]interface{}{
		"id":    "sc-005",
		"title": "Has Image",
		"image": map[string]interface{}{
			"path":      "/some/path.png",
			"thumbnail": "",
		},
	}
	sc, err := ParseScenarioContextFromMap(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sc.Image == nil {
		t.Fatal("expected Image to be set")
	}
	if sc.Image.Path != "/some/path.png" {
		t.Errorf("Image.Path=%q, want '/some/path.png'", sc.Image.Path)
	}
}

func TestParseScenario_Base64ImageNotStoredInPath(t *testing.T) {
	m := map[string]interface{}{
		"id":    "sc-006",
		"title": "Big Image",
		"image": map[string]interface{}{
			"data": "data:image/png;base64,AAAA==",
			"name": "pic.png",
		},
	}
	sc, err := ParseScenarioContextFromMap(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// ParseImageFromValue ignores maps with no path/thumbnail — data-only images return nil
	if sc.Image != nil && sc.Image.Path == "data:image/png;base64,AAAA==" {
		t.Error("raw base64 data must not be stored verbatim in Image.Path")
	}
}

// ---------------------------------------------------------------------------
// ParseCharacterContextFromMap
// ---------------------------------------------------------------------------

func TestParseCharacter_BasicFields(t *testing.T) {
	m := map[string]interface{}{
		"id":          "char_001",
		"name":        "Emily",
		"personality": "Confident and loving",
		"backstory":   "Penthouse kid",
	}
	cc, err := ParseCharacterContextFromMap(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cc.Name != "Emily" {
		t.Errorf("Name=%q, want 'Emily'", cc.Name)
	}
	if cc.Personality != "Confident and loving" {
		t.Errorf("Personality=%q", cc.Personality)
	}
	if cc.ID != "char_001" {
		t.Errorf("ID=%q, want 'char_001'", cc.ID)
	}
}

func TestParseCharacter_NoImagePanic(t *testing.T) {
	m := map[string]interface{}{
		"id":   "char_002",
		"name": "Bri",
	}
	cc, err := ParseCharacterContextFromMap(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cc.Image != nil {
		t.Errorf("expected nil Image when not provided, got %+v", cc.Image)
	}
}

// ---------------------------------------------------------------------------
// ScenarioContext JSON round-trip (pure struct — no external deps)
// ---------------------------------------------------------------------------

func TestScenarioContextJSONRoundTrip(t *testing.T) {
	orig := ScenarioContext{}
	orig.ID = "rt-001"
	orig.Name = "Round Trip"
	orig.Characters = []ScenarioContextCharacter{
		{ID: "c1", Name: "Alice", Role: "hero"},
		{ID: "c2", Name: "Bob", Role: "sidekick"},
	}

	raw, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var restored ScenarioContext
	if err := json.Unmarshal(raw, &restored); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if restored.Name != "Round Trip" {
		t.Errorf("Name=%q after round-trip", restored.Name)
	}
	if len(restored.Characters) != 2 {
		t.Fatalf("Characters=%d after round-trip, want 2", len(restored.Characters))
	}
	if restored.Characters[0].Name != "Alice" || restored.Characters[1].Name != "Bob" {
		t.Errorf("wrong names: %v", restored.Characters)
	}
}

// ---------------------------------------------------------------------------
// Persistence round-trip: save scenario + load, verify name and characters
// ---------------------------------------------------------------------------

func TestPersistenceRoundTrip_ScenarioNameAndCharacters(t *testing.T) {
	dir := t.TempDir()
	pm, err := NewPersistenceManagerWithDir(dir)
	if err != nil {
		t.Fatalf("NewPersistenceManagerWithDir: %v", err)
	}

	chatData, err := pm.NewChat("Test Chat")
	if err != nil {
		t.Fatalf("NewChat: %v", err)
	}

	chatData.Scenario = &ScenarioContext{}
	chatData.Scenario.ID = "sc-persist-001"
	chatData.Scenario.Name = "Persist Test Scenario"
	chatData.Scenario.Characters = []ScenarioContextCharacter{
		{ID: "char_a", Name: "Alice", Role: "hero"},
		{ID: "char_b", Name: "Bob", Role: "sidekick"},
	}

	if err := pm.SaveChat(chatData); err != nil {
		t.Fatalf("SaveChat: %v", err)
	}

	// Read scenario.json directly — verifies the on-disk format
	scenarioPath := filepath.Join(dir, "chats", chatData.Meta.ID, "scenario.json")
	raw, err := os.ReadFile(scenarioPath)
	if err != nil {
		t.Fatalf("reading scenario.json: %v", err)
	}
	var saved ScenarioContext
	if err := json.Unmarshal(raw, &saved); err != nil {
		t.Fatalf("unmarshal scenario.json: %v", err)
	}
	if saved.Name != "Persist Test Scenario" {
		t.Errorf("scenario.json Name=%q, want 'Persist Test Scenario'", saved.Name)
	}
	if len(saved.Characters) != 2 {
		t.Errorf("scenario.json Characters=%d, want 2 — check ScenarioContext.Characters field in types.go", len(saved.Characters))
	}

	// Load it back through the full LoadChat path
	loaded, err := pm.LoadChat(chatData.Meta.ID)
	if err != nil {
		t.Fatalf("LoadChat: %v", err)
	}
	if loaded.Scenario == nil {
		t.Fatal("loaded.Scenario is nil")
	}
	if loaded.Scenario.Name != "Persist Test Scenario" {
		t.Errorf("loaded.Scenario.Name=%q, want 'Persist Test Scenario'", loaded.Scenario.Name)
	}
	if len(loaded.Scenario.Characters) != 2 {
		t.Errorf("loaded.Scenario.Characters=%d, want 2", len(loaded.Scenario.Characters))
	}
}

func TestPersistenceRoundTrip_BothCharactersInCharactersJSON(t *testing.T) {
	dir := t.TempDir()
	pm, err := NewPersistenceManagerWithDir(dir)
	if err != nil {
		t.Fatalf("NewPersistenceManagerWithDir: %v", err)
	}

	chatData, err := pm.NewChat("Test Chat 2")
	if err != nil {
		t.Fatalf("NewChat: %v", err)
	}

	chatData.Characters = []*CharacterContext{
		{},
		{},
	}
	chatData.Characters[0].ID = "char_emily"
	chatData.Characters[0].Name = "Emily"
	chatData.Characters[1].ID = "char_bri"
	chatData.Characters[1].Name = "Bri"

	if err := pm.SaveChat(chatData); err != nil {
		t.Fatalf("SaveChat: %v", err)
	}

	charPath := filepath.Join(dir, "chats", chatData.Meta.ID, "characters.json")
	raw, err := os.ReadFile(charPath)
	if err != nil {
		t.Fatalf("reading characters.json: %v", err)
	}
	var saved []*CharacterContext
	if err := json.Unmarshal(raw, &saved); err != nil {
		t.Fatalf("unmarshal characters.json: %v", err)
	}
	if len(saved) != 2 {
		t.Errorf("characters.json has %d entries, want 2", len(saved))
	}
	names := map[string]bool{}
	for _, c := range saved {
		names[c.Name] = true
	}
	for _, want := range []string{"Emily", "Bri"} {
		if !names[want] {
			t.Errorf("character %q not found in characters.json", want)
		}
	}
}
