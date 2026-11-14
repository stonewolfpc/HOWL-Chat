package memory

import (
	"encoding/json"
	"testing"

	"howl-chat/internal/backend/types"
)

func memoryRecord(id, owner, kind string) types.MemoryRecord {
	return types.MemoryRecord{
		ID:         id,
		Title:      id,
		Summary:    id,
		Content:    "content-" + id,
		Scope:      "character",
		OwnerID:    owner,
		OwnerKind:  kind,
		Status:     "active",
		Confidence: 90,
	}
}

func TestMemoryCaptureUsesStableIDAndFailsClosedOnAmbiguity(t *testing.T) {
	detection := types.EventDetectionResult{
		EventType:  types.EventRomance,
		Context:    "Aria smiles and takes your hand.",
		Confidence: 90,
	}

	one := BatchClassify([]types.EventDetectionResult{detection}, "msg",
		[]CharacterIdentity{{ID: "char_aria", Name: "Aria"}}, "chat_1")
	if len(one) != 1 {
		t.Fatalf("one identity produced %d candidates", len(one))
	}
	if one[0].OwnerID != "char_aria" || one[0].OwnerKind != types.OwnerKindCharacterID {
		t.Fatalf("ownership = %q/%q, want char_aria/character_id", one[0].OwnerID, one[0].OwnerKind)
	}

	// Two characters with the same display name are not distinguishable from the
	// message text, so the memory must remain chat-owned, never guessed.
	two := BatchClassify([]types.EventDetectionResult{detection}, "msg",
		[]CharacterIdentity{
			{ID: "char_aria_1", Name: "Aria"},
			{ID: "char_aria_2", Name: "Aria"},
		}, "chat_1")
	if two[0].OwnerID != "chat_1" || two[0].OwnerKind != types.OwnerKindChat {
		t.Fatalf("ambiguous ownership = %q/%q, want chat_1/chat", two[0].OwnerID, two[0].OwnerKind)
	}

	// No identity at all must also remain chat-owned rather than inventing one.
	none := BatchClassify([]types.EventDetectionResult{detection}, "msg", nil, "chat_1")
	if none[0].OwnerID != "chat_1" || none[0].OwnerKind != types.OwnerKindChat {
		t.Fatalf("missing identity ownership = %q/%q, want chat_1/chat", none[0].OwnerID, none[0].OwnerKind)
	}
}

func TestContinuityTransferRequiresSameStableIdentity(t *testing.T) {
	config := ContinuityConfig{Enabled: true, Rule: InheritAll, MaxMemories: 20, MinConfidence: 0}
	cc := NewCharacterContinuity()

	legacyName := memoryRecord("legacy-name", "Aria", "")
	oldID := memoryRecord("old-id", "char_old", types.OwnerKindCharacterID)
	newID := memoryRecord("new-id", "char_new", types.OwnerKindCharacterID)
	foreignID := memoryRecord("foreign-id", "char_other", types.OwnerKindCharacterID)
	all := []types.MemoryRecord{legacyName, oldID, newID, foreignID}

	// A deleted and recreated same-name character has a new ID and must not
	// inherit the legacy name-owned record or the old identity's memories.
	exported := cc.ExportCharacterMemories("char_new", all, config)
	if len(exported) != 1 || exported[0].ID != "new-id" {
		t.Fatalf("export = %+v, want only new-id", exported)
	}

	existing := types.MemoryRecord{ID: "existing", Scope: "chat", OwnerID: "chat_1"}
	merged := cc.ImportCharacterMemories("char_new", all, []types.MemoryRecord{existing}, config)
	if len(merged) != 2 {
		t.Fatalf("import produced %d records, want existing + new-id: %+v", len(merged), merged)
	}
	if merged[0].ID != "existing" {
		t.Fatalf("existing record was not preserved: %+v", merged[0])
	}
	if merged[1].ID != "new-id" || merged[1].OwnerID != "char_new" || merged[1].OwnerKind != types.OwnerKindCharacterID {
		t.Fatalf("imported ownership changed or widened: %+v", merged[1])
	}

	for _, rejected := range []string{"legacy-name", "old-id", "foreign-id"} {
		for _, record := range merged {
			if record.ID == rejected {
				t.Fatalf("%s crossed the character identity boundary", rejected)
			}
		}
	}
}

func TestExportRoundTripDoesNotReclassifyLegacyMemory(t *testing.T) {
	// A generic JSON export carries the unresolved owner string verbatim. Even
	// when that string looks exactly like the destination character ID, import
	// must rely on owner_kind and reject it rather than infer ownership.
	legacyIDShape := memoryRecord("legacy-id-shape", "char_new", "")
	legacyName := memoryRecord("legacy-name", "Aria", "")
	stable := memoryRecord("stable", "char_new", types.OwnerKindCharacterID)

	data, err := json.Marshal([]types.MemoryRecord{legacyIDShape, legacyName, stable})
	if err != nil {
		t.Fatal(err)
	}
	var decoded []types.MemoryRecord
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded) != 3 {
		t.Fatalf("round trip changed record count: %d", len(decoded))
	}
	if decoded[0].OwnerKind != "" || decoded[1].OwnerKind != "" {
		t.Fatalf("unresolved ownership was reclassified: %+v", decoded[:2])
	}

	config := ContinuityConfig{Enabled: true, Rule: InheritAll, MaxMemories: 20, MinConfidence: 0}
	cc := NewCharacterContinuity()
	exported := cc.ExportCharacterMemories("char_new", decoded, config)
	if len(exported) != 1 || exported[0].ID != "stable" {
		t.Fatalf("export after JSON round trip = %+v, want only stable", exported)
	}

	merged := cc.ImportCharacterMemories("char_new", decoded, nil, config)
	if len(merged) != 1 || merged[0].ID != "stable" {
		t.Fatalf("import after JSON round trip = %+v, want only stable", merged)
	}
	if merged[0].OwnerKind != types.OwnerKindCharacterID || merged[0].OwnerID != "char_new" {
		t.Fatalf("stable classification changed on import: %+v", merged[0])
	}
}

func TestLegacyMemoryIsNeverReassignedByExportOrImport(t *testing.T) {
	config := ContinuityConfig{Enabled: true, Rule: InheritAll, MaxMemories: 20, MinConfidence: 0}
	cc := NewCharacterContinuity()
	legacy := memoryRecord("legacy", "Aria", "")

	if exported := cc.ExportCharacterMemories("char_new", []types.MemoryRecord{legacy}, config); len(exported) != 0 {
		t.Fatalf("legacy memory exported: %+v", exported)
	}
	merged := cc.ImportCharacterMemories("char_new", []types.MemoryRecord{legacy}, nil, config)
	if len(merged) != 0 {
		t.Fatalf("legacy memory imported: %+v", merged)
	}
	if legacy.OwnerID != "Aria" || legacy.OwnerKind != "" {
		t.Fatalf("legacy ownership mutated: %+v", legacy)
	}
}

func TestStableOwnershipRequiresMatchingCharacterID(t *testing.T) {
	stable := memoryRecord("stable", "char_aria", types.OwnerKindCharacterID)
	if !IsStableCharacterOwnership(stable, "char_aria") {
		t.Fatal("stable identity was rejected")
	}
	if IsStableCharacterOwnership(stable, "char_other") {
		t.Fatal("a foreign character ID was accepted")
	}
	if IsStableCharacterOwnership(stable, "") {
		t.Fatal("an empty character ID was accepted")
	}

	nameOwned := memoryRecord("name-owned", "char_aria", "")
	if IsStableCharacterOwnership(nameOwned, "char_aria") {
		t.Fatal("a legacy record that happens to equal an ID was accepted")
	}
}

func TestPromptOwnershipFilterRejectsLegacyAndForeignRecords(t *testing.T) {
	legacy := memoryRecord("legacy", "char_aria", "")
	stable := memoryRecord("stable", "char_aria", types.OwnerKindCharacterID)
	foreign := memoryRecord("foreign", "char_other", types.OwnerKindCharacterID)

	if IsStableCharacterOwnershipForAny(legacy, []string{"char_aria"}) {
		t.Fatal("legacy record passed the prompt filter")
	}
	if !IsStableCharacterOwnershipForAny(stable, []string{"char_aria"}) {
		t.Fatal("stable record failed the prompt filter")
	}
	if IsStableCharacterOwnershipForAny(foreign, []string{"char_aria"}) {
		t.Fatal("foreign record passed the prompt filter")
	}
}

func TestMemoryIdentityDiagnosticsReportsUnresolvedWithoutMutation(t *testing.T) {
	legacy := memoryRecord("legacy", "Aria", "")
	stable := memoryRecord("stable", "char_aria", types.OwnerKindCharacterID)
	chatScoped := types.MemoryRecord{ID: "chat", Scope: "chat", OwnerID: "chat_1"}

	diag := summarizeMemoryIdentities([]types.MemoryRecord{legacy, stable, chatScoped})
	if diag.Total != 2 || diag.Stable != 1 || diag.Unresolved != 1 {
		t.Fatalf("diagnostics = %+v", diag)
	}
	if len(diag.UnresolvedOwners) != 1 || diag.UnresolvedOwners[0] != "Aria" {
		t.Fatalf("unresolved owners = %v", diag.UnresolvedOwners)
	}
	if legacy.OwnerID != "Aria" || legacy.OwnerKind != "" {
		t.Fatalf("diagnostics mutated legacy ownership: %+v", legacy)
	}
}

func TestMemoryToLorebookEntryCarriesOwnerKind(t *testing.T) {
	entry := ToLorebookEntry(memoryRecord("stable", "char_aria", types.OwnerKindCharacterID))
	if entry.OwnerID != "char_aria" || entry.OwnerKind != types.OwnerKindCharacterID {
		t.Fatalf("converted entry ownership = %q/%q", entry.OwnerID, entry.OwnerKind)
	}
}
