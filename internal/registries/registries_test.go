package registries

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"howl-chat/internal/contentimport"
	"howl-chat/internal/library"
	"howl-chat/internal/modelintel"
)

func TestNewBuildsEveryAuthorityForAnEmptyRoot(t *testing.T) {
	dir := t.TempDir()
	set, errs := New(dir)
	if len(errs) != 0 {
		t.Fatalf("New reported errors for an empty root: %v", errs)
	}
	if set == nil {
		t.Fatal("New returned no owner")
	}
	if set.Libraries() == nil || set.Models() == nil || set.Intel() == nil ||
		set.Inventory() == nil || set.Lorebooks() == nil ||
		set.Worlds() == nil || set.Scenarios() == nil {
		t.Fatalf("owner is incomplete: %+v", set)
	}

	if got := set.Libraries().List(""); len(got) != 0 {
		t.Fatalf("libraries = %+v, want empty", got)
	}
	if got := set.Models().All(); len(got) != 0 {
		t.Fatalf("models = %+v, want empty", got)
	}
	if list := set.Lorebooks().List(); list.Cutover {
		t.Fatalf("lorebooks cut over without an import: %+v", list)
	}
	if _, ok := set.Inventory().Cached("unknown"); ok {
		t.Fatal("inventory returned a snapshot for an unknown root")
	}
	profile := modelintel.Profile{Path: filepath.Join(dir, "model.gguf"), Name: "model"}
	if err := set.Intel().Put(profile); err != nil {
		t.Fatalf("intel Put: %v", err)
	}
	if _, ok := set.Intel().Get(profile.Path, 0, 0, ""); !ok {
		t.Fatal("intel did not return the profile it stored")
	}
}

func TestNilOwnerReportsNoComponents(t *testing.T) {
	var set *Set
	if set.Models() != nil || set.Libraries() != nil || set.Inventory() != nil ||
		set.Worlds() != nil || set.Scenarios() != nil || set.Lorebooks() != nil || set.Intel() != nil {
		t.Fatal("a nil owner handed out a component")
	}
}

func TestComposeCarriesOnlyTheGivenParts(t *testing.T) {
	dir := t.TempDir()
	libs, err := library.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	inventory := library.NewInventoryStore(dir)
	set := Compose(Parts{Libraries: libs, Inventory: inventory})

	if set.Libraries() != libs || set.Inventory() != inventory {
		t.Fatal("Compose did not carry the given parts")
	}
	if set.Models() != nil || set.Worlds() != nil || set.Scenarios() != nil ||
		set.Lorebooks() != nil || set.Intel() != nil {
		t.Fatal("Compose invented an unset part")
	}
}

func TestNewKeepsIndependentAuthoritiesWhenLibrariesAreCorrupt(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "content-libraries.json"), "{not json")
	set, errs := New(dir)

	if len(errs) == 0 {
		t.Fatal("New hid a corrupt content-library registry")
	}
	if set.Libraries() != nil {
		t.Fatal("a corrupt registry was adopted")
	}
	if set.Worlds() != nil || set.Scenarios() != nil {
		t.Fatal("import services were built without their registry")
	}
	// These do not depend on the content registry and must survive.
	if set.Models() == nil || set.Intel() == nil || set.Inventory() == nil || set.Lorebooks() == nil {
		t.Fatal("a content-library failure discarded an independent authority")
	}
}

func TestNewReportsEachCorruptStoreWithoutDroppingOthers(t *testing.T) {
	cases := []struct {
		name    string
		file    string
		content string
		drops   func(*Set) bool
	}{
		{"model registry", "model-registry.json", "{not json", func(s *Set) bool { return s.Models() == nil }},
		{"model intelligence", "model-intelligence.json", "[]", func(s *Set) bool { return s.Intel() == nil }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			writeFile(t, filepath.Join(dir, tc.file), tc.content)
			set, errs := New(dir)
			if len(errs) == 0 {
				t.Fatalf("New hid a corrupt %s", tc.name)
			}
			if !tc.drops(set) {
				t.Fatalf("corrupt %s was adopted", tc.name)
			}
			if set.Libraries() == nil || set.Worlds() == nil || set.Scenarios() == nil ||
				set.Inventory() == nil || set.Lorebooks() == nil {
				t.Fatalf("corrupt %s discarded the content authorities", tc.name)
			}
		})
	}
}

func TestCorruptInventoryIsIgnoredRatherThanAdopted(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "content-library-inventory.json"), "{not json")
	set, errs := New(dir)
	if len(errs) != 0 {
		t.Fatalf("inventory corruption became a boot error: %v", errs)
	}
	if set.Inventory() == nil {
		t.Fatal("inventory store is absent")
	}
	if _, ok := set.Inventory().Cached("root"); ok {
		t.Fatal("inventory adopted a corrupt snapshot")
	}
}

func TestLibrariesCollapseDuplicateRegistrations(t *testing.T) {
	set, errs := New(t.TempDir())
	if len(errs) != 0 {
		t.Fatal(errs)
	}
	root := t.TempDir()
	first, err := set.Libraries().Add(library.Registration{Category: library.Characters, Path: root})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	second, err := set.Libraries().Add(library.Registration{Category: library.Characters, Path: root})
	if err != nil {
		t.Fatalf("Add duplicate: %v", err)
	}
	if first.ID != second.ID {
		t.Fatalf("duplicate registration ids differ: %q vs %q", first.ID, second.ID)
	}
	if got := set.Libraries().List(library.Characters); len(got) != 1 {
		t.Fatalf("libraries = %+v, want one registration", got)
	}
}

func TestWorldImportRefusesAStalePlanThroughTheOwner(t *testing.T) {
	set, errs := New(t.TempDir())
	if len(errs) != 0 {
		t.Fatal(errs)
	}
	result, err := set.Worlds().Execute(contentimport.Request{Category: "worlds"}, "not-the-reviewed-hash")
	if err != nil {
		t.Fatalf("Execute returned an error instead of a surfaced refusal: %v", err)
	}
	if result.Status != "plan_stale" || result.Message == "" {
		t.Fatalf("result = %+v, want a plan_stale refusal with a message", result)
	}
	if _, err := set.Worlds().Execute(contentimport.Request{Category: "scenarios"}, ""); err == nil {
		t.Fatal("the world service accepted a scenario request")
	}
}

func TestAccessorsAreSafeUnderConcurrentUse(t *testing.T) {
	set, errs := New(t.TempDir())
	if len(errs) != 0 {
		t.Fatal(errs)
	}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := 0; n < 100; n++ {
				if set.Libraries() == nil || set.Models() == nil || set.Intel() == nil ||
					set.Inventory() == nil || set.Lorebooks() == nil ||
					set.Worlds() == nil || set.Scenarios() == nil {
					t.Error("accessor returned nil during concurrent use")
					return
				}
			}
		}()
	}
	wg.Wait()
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
