package chat

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"howl-chat/internal/backend/lorebook"
	"howl-chat/internal/library"
	"howl-chat/internal/worldstore"
)

func newWorldService(t *testing.T) (*Service, *library.Registry, *worldstore.Store, string) {
	t.Helper()
	reg, err := library.New(t.TempDir())
	if err != nil {
		t.Fatalf("library registry: %v", err)
	}
	root := t.TempDir()
	if _, err := reg.Add(library.Registration{
		Category: library.Worlds,
		Path:     root,
		Role:     library.RoleRead | library.RoleWrite | library.RoleSave,
		Origin:   library.OriginHowl,
		Format:   library.FormatHowl,
	}); err != nil {
		t.Fatalf("register world root: %v", err)
	}
	store, err := worldstore.New(reg, library.Worlds)
	if err != nil {
		t.Fatalf("world store: %v", err)
	}
	service := &Service{}
	if err := service.SetContentLibraries(reg); err != nil {
		t.Fatalf("attach store: %v", err)
	}
	return service, reg, store, root
}

func saveWorld(t *testing.T, store *worldstore.Store, id string) {
	t.Helper()
	if _, err := store.Save(map[string]interface{}{
		"id":          id,
		"name":        "Ashfall",
		"description": "A city under falling ash",
		"key_facts":   []interface{}{"ash"},
	}); err != nil {
		t.Fatalf("save world: %v", err)
	}
}

func writeLore(t *testing.T, root, id, filename, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, id, filename), []byte(body), 0644); err != nil {
		t.Fatalf("write %s: %v", filename, err)
	}
}

func TestLoadCanonicalWorldUsesStableIDAndClassifiesLegacyLore(t *testing.T) {
	service, _, store, root := newWorldService(t)
	saveWorld(t, store, "wld_context")
	writeLore(t, root, "wld_context", "world_lore.json",
		`{"entries":[{"id":"lb_manual","scope":"world","content":"manual truth"}]}`)
	writeLore(t, root, "wld_context", "automated_world_lore.json",
		`{"entries":[{"id":"lb_generated","scope":"world","content":"generated truth"}]}`)
	automatedPath := filepath.Join(root, "wld_context", "automated_world_lore.json")
	before, err := os.ReadFile(automatedPath)
	if err != nil {
		t.Fatalf("read lore before load: %v", err)
	}

	world, manual, automated, status := service.loadCanonicalWorld("wld_context")
	if status != ContextAvailable {
		t.Fatalf("status = %q, want %q", status, ContextAvailable)
	}
	if world == nil || world.Name != "Ashfall" {
		t.Fatalf("world = %+v", world)
	}
	if len(manual) != 1 || manual[0].OwnerID != "wld_context" {
		t.Fatalf("manual lore = %+v, want owner stamped from canonical folder", manual)
	}
	if len(automated) != 1 {
		t.Fatalf("automated lore = %+v", automated)
	}
	if !automated[0].Automated || automated[0].OriginKind != lorebook.OriginLegacyGeneratedUnattributed {
		t.Fatalf("legacy generated classification = %+v", automated[0])
	}
	if automated[0].OwnerID != "wld_context" {
		t.Fatalf("legacy owner = %q, want stable ID", automated[0].OwnerID)
	}
	after, err := os.ReadFile(automatedPath)
	if err != nil {
		t.Fatalf("read lore after load: %v", err)
	}
	if !bytes.Equal(after, before) {
		t.Fatal("reading canonical lore rewrote the file; classification must stay in memory")
	}
}

func TestLegacyNameFolderIsNotAReadFallback(t *testing.T) {
	service, _, store, root := newWorldService(t)
	saveWorld(t, store, "wld_context")
	nameDir := filepath.Join(root, "Ashfall")
	if err := os.MkdirAll(nameDir, 0755); err != nil {
		t.Fatalf("create name folder: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nameDir, "world.json"), []byte(`{"name":"Legacy Folder World"}`), 0644); err != nil {
		t.Fatalf("write decoy world: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nameDir, "automated_world_lore.json"), []byte(`{"entries":[{"id":"lb_decoy","content":"name-folder lore"}]}`), 0644); err != nil {
		t.Fatalf("write decoy lore: %v", err)
	}

	world, _, automated, status := service.loadCanonicalWorld("wld_context")
	if status != ContextAvailable || world == nil || world.Name != "Ashfall" {
		t.Fatalf("stable ID did not load canonical document: status=%q world=%+v", status, world)
	}
	if len(automated) != 0 {
		t.Fatalf("name-folder lore entered the canonical prompt pool: %+v", automated)
	}

	world, _, _, status = service.loadCanonicalWorld("Ashfall")
	if world != nil || (status != ContextUnavailable && status != ContextNotFound) {
		t.Fatalf("name was used as an identity: status=%q world=%+v", status, world)
	}
}

func TestLoadCanonicalWorldExcludesForeignOwner(t *testing.T) {
	service, _, store, root := newWorldService(t)
	saveWorld(t, store, "wld_owner")
	writeLore(t, root, "wld_owner", "world_lore.json",
		`{"entries":[{"id":"lb_foreign","scope":"world","ownerId":"Ashfall","content":"name-owned"}]}`)

	_, manual, _, status := service.loadCanonicalWorld("wld_owner")
	if status != ContextAvailable {
		t.Fatalf("status = %q, want %q", status, ContextAvailable)
	}
	if len(manual) != 0 {
		t.Fatalf("foreign-owner lore entered the prompt pool: %+v", manual)
	}
}

func TestLoadCanonicalWorldFailsClosedWhenAuthorityIsIndeterminate(t *testing.T) {
	service, reg, store, root := newWorldService(t)
	saveWorld(t, store, "wld_shadowed")
	writeLore(t, root, "wld_shadowed", "world_lore.json",
		`{"entries":[{"id":"lb_1","scope":"world","content":"stale truth"}]}`)

	unavailable := t.TempDir()
	if _, err := reg.Add(library.Registration{
		Category: library.Worlds,
		Path:     unavailable,
		Role:     library.RoleRead,
		Origin:   library.OriginExternal,
		Format:   library.FormatUnknown,
	}); err != nil {
		t.Fatalf("register unavailable root: %v", err)
	}
	if err := os.RemoveAll(unavailable); err != nil {
		t.Fatalf("remove root: %v", err)
	}

	world, manual, automated, status := service.loadCanonicalWorld("wld_shadowed")
	if status != ContextAuthorityIndeterminate {
		t.Fatalf("status = %q, want %q", status, ContextAuthorityIndeterminate)
	}
	if world != nil || manual != nil || automated != nil {
		t.Fatalf("indeterminate content reached the prompt: world=%v manual=%v automated=%v", world, manual, automated)
	}
}

func TestLoadCanonicalWorldDistinguishesUnsetAndAbsent(t *testing.T) {
	service, _, store, _ := newWorldService(t)
	saveWorld(t, store, "wld_present")

	if _, _, _, status := service.loadCanonicalWorld(""); status != ContextNotSet {
		t.Fatalf("empty ID status = %q, want %q", status, ContextNotSet)
	}
	if _, _, _, status := service.loadCanonicalWorld("wld_absent"); status != ContextNotFound {
		t.Fatalf("absent ID status = %q, want %q", status, ContextNotFound)
	}
}

func TestAutomatedLoreSyncTargetIsCharacterOnly(t *testing.T) {
	service := &Service{}
	for _, entry := range []lorebook.Entry{
		{Scope: lorebook.ScopeWorld, OwnerID: "wld_1", Automated: true},
		{Scope: lorebook.ScopeScenario, OwnerID: "scn_1", Automated: true},
		{Scope: lorebook.ScopeCharacter, OwnerID: "char_1", Automated: false},
		{Scope: lorebook.ScopeCharacter, OwnerID: "", Automated: true},
	} {
		if target, ok := service.automatedLoreSyncTarget("root", entry); ok {
			t.Errorf("scope %q with automated=%v produced write target %q", entry.Scope, entry.Automated, target)
		}
	}

	target, ok := service.automatedLoreSyncTarget("root", lorebook.Entry{
		Scope: lorebook.ScopeCharacter, OwnerID: "char_1", Automated: true,
	})
	if !ok {
		t.Fatal("character generated lore lost its existing write target")
	}
	if !strings.HasSuffix(filepath.ToSlash(target), "Characters/char_1/automated_character_lore.json") {
		t.Fatalf("target = %q", target)
	}
}
