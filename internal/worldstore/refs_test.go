package worldstore

import "testing"

func TestNormalizeReferencesKeepsIdentityAndLegacyName(t *testing.T) {
	raw := []interface{}{
		map[string]interface{}{"id": "1794319623163", "title": "Tide Oath"},
		map[string]interface{}{"name": "Old Harbor"},
		"Plain Name",
		map[string]interface{}{"id": "1794319623163", "title": "Tide Oath (duplicate)"},
		nil,
	}

	refs := NormalizeReferences(RefLorebook, raw)
	if len(refs) != 3 {
		t.Fatalf("refs = %d, want 3: %#v", len(refs), refs)
	}
	if refs[0].ID != "1794319623163" || refs[0].Title != "Tide Oath" {
		t.Fatalf("identity reference = %#v", refs[0])
	}
	if refs[0].LegacyName != "" {
		t.Fatalf("identity reference carried a legacy name: %#v", refs[0])
	}
	if refs[1].ID != "" || refs[1].LegacyName != "Old Harbor" {
		t.Fatalf("name-only reference = %#v", refs[1])
	}
	if refs[2].LegacyName != "Plain Name" {
		t.Fatalf("string reference = %#v", refs[2])
	}
}

func TestResolveReferencesIsIDOnly(t *testing.T) {
	catalog := Catalog{RefLorebook: {"lb_1": true}}
	refs := []Reference{
		{Kind: RefLorebook, ID: "lb_1", Title: "Tide Oath"},
		{Kind: RefLorebook, ID: "lb_2", Title: "Tide Oath"},
		{Kind: RefLorebook, LegacyName: "Tide Oath"},
	}

	resolutions := ResolveReferences(refs, catalog)
	if len(resolutions) != 3 {
		t.Fatalf("resolutions = %d", len(resolutions))
	}
	if !resolutions[0].Resolved {
		t.Fatal("known identity did not resolve")
	}
	if resolutions[1].Resolved || resolutions[1].Reason != "not_found" {
		t.Fatalf("same-title unknown identity resolved: %#v", resolutions[1])
	}
	if resolutions[2].Resolved || resolutions[2].Reason != "legacy_name_unresolved" {
		t.Fatalf("legacy name resolved by title: %#v", resolutions[2])
	}
}

func TestNormalizeReferencesRejectsExplicitKindMismatch(t *testing.T) {
	refs := NormalizeReferences(RefLorebook, []interface{}{
		map[string]interface{}{"kind": "character", "id": "char_1"},
		map[string]interface{}{"kind": "lorebook", "id": "lb_1"},
	})
	if len(refs) != 1 || refs[0].ID != "lb_1" {
		t.Fatalf("refs = %#v", refs)
	}
}

func TestResolveReferencesKeepsKindsSeparate(t *testing.T) {
	catalog := Catalog{RefWorld: {"shared_id": true}}
	refs := []Reference{{Kind: RefLorebook, ID: "shared_id"}}

	resolutions := ResolveReferences(refs, catalog)
	if resolutions[0].Resolved {
		t.Fatal("an identity resolved across categories")
	}
}

func TestUnresolvedReferencesPreservesOrder(t *testing.T) {
	catalog := Catalog{RefCharacter: {"char_aria": true}}
	refs := []Reference{
		{Kind: RefCharacter, ID: "char_bram"},
		{Kind: RefCharacter, ID: "char_aria"},
		{Kind: RefCharacter, LegacyName: "Old Bram"},
	}

	unresolved := UnresolvedReferences(refs, catalog)
	if len(unresolved) != 2 {
		t.Fatalf("unresolved = %d, want 2", len(unresolved))
	}
	if unresolved[0].ID != "char_bram" || unresolved[1].LegacyName != "Old Bram" {
		t.Fatalf("order changed: %#v", unresolved)
	}
}
