package main

import (
	"strings"
	"testing"

	"howl-chat/internal/backend/chat"
	"howl-chat/internal/backend/lorebook"
	entities "howl-chat/internal/backend/types/entities"
)

func privacyEntry(id, title, content string, scope lorebook.Scope, owner, visibility string) lorebook.Entry {
	return lorebook.Entry{
		ID:               id,
		Title:            title,
		Content:          content,
		Enabled:          true,
		Scope:            scope,
		OwnerID:          owner,
		Visibility:       visibility,
		TriggerPhrases:   []string{"sigil"},
		TriggerMode:      lorebook.TriggerLoose,
		TriggerDirection: lorebook.TriggerBoth,
		TriggerFrequency: lorebook.FrequencyAlways,
		PriorityLevel:    3,
	}
}

func privacyLoreApp() *App {
	return &App{
		characters: []*chat.CharacterContext{
			{Character: entities.Character{ID: "char_aria", Name: "Aria"}},
			{Character: entities.Character{ID: "char_bram", Name: "Bram"}},
		},
		lorebooks: []lorebook.Entry{
			privacyEntry("aria-private", "Aria Secret", "ARIA-SECRET", lorebook.ScopeCharacter, "char_aria", lorebook.VisibilityPrivate),
			privacyEntry("bram-private", "Bram Secret", "BRAM-SECRET", lorebook.ScopeCharacter, "char_bram", lorebook.VisibilityPrivate),
			privacyEntry("unowned-private", "Unowned Secret", "UNOWNED-SECRET", lorebook.ScopeCharacter, "", lorebook.VisibilityPrivate),
			privacyEntry("hidden-world", "Hidden World", "HIDDEN-SECRET", lorebook.ScopeWorld, "", lorebook.VisibilityHidden),
			privacyEntry("public-unowned", "Public Lore", "PUBLIC-LORE", lorebook.ScopeCharacter, "", lorebook.VisibilityPublic),
			privacyEntry("legacy-unowned", "Legacy Lore", "LEGACY-LORE", lorebook.ScopeCharacter, "", ""),
		},
		chatMessages: []map[string]interface{}{{"role": "user", "content": "I speak of the sigil"}},
	}
}

// Routing helpers can be right while the assembled prompt is wrong. This checks
// the block the assembler actually receives, for two characters in turn.
func TestCharacterLorePrivacyInAssembledPrompt(t *testing.T) {
	app := privacyLoreApp()

	aria := app.resolveCharacterLorePromptBlock("char_aria")
	if !strings.Contains(aria, "ARIA-SECRET") {
		t.Fatalf("Aria's own lore missing:\n%s", aria)
	}
	for _, leaked := range []string{"BRAM-SECRET", "UNOWNED-SECRET", "HIDDEN-SECRET"} {
		if strings.Contains(aria, leaked) {
			t.Errorf("%s leaked into Aria's block:\n%s", leaked, aria)
		}
	}
	for _, want := range []string{"PUBLIC-LORE", "LEGACY-LORE"} {
		if !strings.Contains(aria, want) {
			t.Errorf("%s missing from Aria's block:\n%s", want, aria)
		}
	}

	bram := app.resolveCharacterLorePromptBlock("char_bram")
	if !strings.Contains(bram, "BRAM-SECRET") {
		t.Fatalf("Bram's own lore missing:\n%s", bram)
	}
	for _, leaked := range []string{"ARIA-SECRET", "UNOWNED-SECRET", "HIDDEN-SECRET"} {
		if strings.Contains(bram, leaked) {
			t.Errorf("%s leaked into Bram's block:\n%s", leaked, bram)
		}
	}

	// Rebuilding for Bram after Aria must not reuse anything from Aria's turn.
	if again := app.resolveCharacterLorePromptBlock("char_bram"); strings.Contains(again, "ARIA-SECRET") {
		t.Errorf("Aria's private lore survived a character switch:\n%s", again)
	}

	// The assembled prompt must also refuse hidden and unowned private lore, and
	// keep public and legacy lore working.
	assembled := promptSystem(t, app)
	for _, leaked := range []string{"UNOWNED-SECRET", "HIDDEN-SECRET", "ARIA-SECRET", "BRAM-SECRET"} {
		if strings.Contains(assembled, leaked) {
			t.Errorf("%s reached the assembled prompt:\n%s", leaked, assembled)
		}
	}
	for _, want := range []string{"PUBLIC-LORE", "LEGACY-LORE"} {
		if !strings.Contains(assembled, want) {
			t.Errorf("%s missing from the assembled prompt:\n%s", want, assembled)
		}
	}
	if len(app.chatMessages) != 1 {
		t.Errorf("prompt assembly mutated chat history: %d messages", len(app.chatMessages))
	}
}

// The character-turn binding must apply the same policy as the general path.
func TestBuildCharacterLoreBlockSharesPolicy(t *testing.T) {
	app := privacyLoreApp()

	aria := app.BuildCharacterLoreBlock("char_aria", "I speak of the sigil")
	for _, want := range []string{"ARIA-SECRET", "PUBLIC-LORE", "LEGACY-LORE"} {
		if !strings.Contains(aria, want) {
			t.Errorf("%s missing from Aria's character block:\n%s", want, aria)
		}
	}
	for _, leaked := range []string{"BRAM-SECRET", "UNOWNED-SECRET", "HIDDEN-SECRET"} {
		if strings.Contains(aria, leaked) {
			t.Errorf("%s leaked into Aria's character block:\n%s", leaked, aria)
		}
	}

	bram := app.BuildCharacterLoreBlock("char_bram", "I speak of the sigil")
	if strings.Contains(bram, "ARIA-SECRET") {
		t.Errorf("Aria's private lore leaked into Bram's character block:\n%s", bram)
	}
	if !strings.Contains(bram, "BRAM-SECRET") {
		t.Errorf("Bram's own lore missing:\n%s", bram)
	}

	// A turn without a trigger must not pull the whole library in. Clear the
	// fixture history first: recent turns are part of the scan, so a previous
	// mention of the phrase would legitimately match.
	app.chatMessages = []map[string]interface{}{{"role": "user", "content": "earlier unrelated line"}}
	if block := app.BuildCharacterLoreBlock("char_aria", "I look at the sky"); strings.Contains(block, "ARIA-SECRET") {
		t.Errorf("untriggered lore was injected:\n%s", block)
	}
}

// A stale entry in the pool must never reach the assembled prompt either.
func TestStaleLoreIsExcludedFromAssembledPrompt(t *testing.T) {
	app := privacyLoreApp()
	entry := privacyEntry("stale-1", "Offline", "STALE-SECRET", lorebook.ScopeWorld, "", lorebook.VisibilityPublic)
	entry.Stale = true
	app.lorebooks = append(app.lorebooks, entry)

	system := promptSystem(t, app)
	if strings.Contains(system, "STALE-SECRET") {
		t.Fatalf("stale lore reached the prompt:\n%s", system)
	}
}
