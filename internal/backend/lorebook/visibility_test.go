package lorebook

import "testing"

// "hidden" is not a hint: the entry must never enter a prompt, whatever its scope.
func TestResolveSkipsHiddenVisibility(t *testing.T) {
	entry := Entry{
		ID:               "hidden-1",
		Title:            "Vault",
		Content:          "The vault code is 4417.",
		Enabled:          true,
		Scope:            ScopeWorld,
		TriggerPhrases:   []string{"vault"},
		TriggerMode:      TriggerExact,
		TriggerDirection: TriggerBoth,
		TriggerFrequency: FrequencyAlways,
		Visibility:       VisibilityHidden,
	}

	if resolved := Resolve([]Entry{entry}, ResolveRequest{Message: "the vault"}); len(resolved) != 0 {
		t.Fatalf("hidden entry resolved: %+v", resolved)
	}
}

// Character lore belongs to one stable character ID. An unowned entry comes from
// the Lorebooks form, which does not write ownerId, so it may only be shared when
// it is public; private unowned lore must not reach another character.
func TestRouteCharacterLoreVisibilityBoundary(t *testing.T) {
	entries := []ResolvedEntry{
		{Entry: Entry{ID: "private-unowned", Scope: ScopeCharacter, Visibility: VisibilityPrivate}},
		{Entry: Entry{ID: "public-unowned", Scope: ScopeCharacter, Visibility: VisibilityPublic}},
		{Entry: Entry{ID: "legacy-unowned", Scope: ScopeCharacter}},
		{Entry: Entry{ID: "owned-aria", Scope: ScopeCharacter, OwnerID: "char_aria", OwnerName: "Aria", Visibility: VisibilityPrivate}},
		{Entry: Entry{ID: "owned-bram", Scope: ScopeCharacter, OwnerID: "char_bram", OwnerName: "Bram", Visibility: VisibilityPublic}},
		{Entry: Entry{ID: "world-1", Scope: ScopeWorld}},
	}

	routed := RouteResolvedEntries(entries, "char_aria", "", "")
	seen := make(map[string]bool, len(routed))
	for _, entry := range routed {
		seen[entry.ID] = true
	}

	for _, want := range []string{"public-unowned", "legacy-unowned", "owned-aria", "world-1"} {
		if !seen[want] {
			t.Errorf("%s was not routed to Aria", want)
		}
	}
	for _, leaked := range []string{"private-unowned", "owned-bram"} {
		if seen[leaked] {
			t.Errorf("%s leaked into Aria's prompt", leaked)
		}
	}
	if len(routed) != 4 {
		t.Fatalf("routed %d entries, want 4: %+v", len(routed), routed)
	}

	bram := RouteResolvedEntries(entries, "char_bram", "", "")
	for _, entry := range bram {
		if entry.ID == "owned-aria" {
			t.Error("Aria's private lore reached Bram")
		}
	}
}

// A display name is not an identity: duplicate names are distinct characters, a
// rename must not change ownership, and ownerName must never authorize delivery.
func TestRouteUsesStableIDNotDisplayName(t *testing.T) {
	entries := []ResolvedEntry{
		{Entry: Entry{ID: "aria-one", Scope: ScopeCharacter, OwnerID: "char_aria_1", OwnerName: "Aria"}},
		{Entry: Entry{ID: "aria-two", Scope: ScopeCharacter, OwnerID: "char_aria_2", OwnerName: "Aria"}},
		{Entry: Entry{ID: "renamed", Scope: ScopeCharacter, OwnerID: "char_aria_2", OwnerName: "Aria the Elder"}},
		{Entry: Entry{ID: "misleading-label", Scope: ScopeCharacter, OwnerID: "char_bram", OwnerName: "Aria"}},
	}

	first := RouteResolvedEntries(entries, "char_aria_1", "", "")
	if len(first) != 1 || first[0].ID != "aria-one" {
		t.Fatalf("first duplicate routed to %+v", first)
	}

	second := RouteResolvedEntries(entries, "char_aria_2", "", "")
	seen := map[string]bool{}
	for _, entry := range second {
		seen[entry.ID] = true
	}
	if !seen["aria-two"] || !seen["renamed"] || seen["aria-one"] || seen["misleading-label"] {
		t.Fatalf("stable ownership routing = %+v", seen)
	}

	// Passing the display name must not match any owner.
	if byName := RouteResolvedEntries(entries, "Aria", "", ""); len(byName) != 0 {
		t.Fatalf("a display name authorized ownership: %+v", byName)
	}
}
