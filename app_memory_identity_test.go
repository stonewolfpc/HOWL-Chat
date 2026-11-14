package main

import (
	"testing"

	"howl-chat/internal/backend/types"
)

func TestTypesMemoryToRecordPreservesOwnerKind(t *testing.T) {
	stable := typesMemoryToRecord(types.MemoryRecord{
		ID:        "stable",
		Scope:     "character",
		OwnerID:   "char_aria",
		OwnerKind: types.OwnerKindCharacterID,
	})
	if stable.OwnerKind != types.OwnerKindCharacterID || stable.ScopeID != "char_aria" {
		t.Fatalf("stable export ownership = %q/%q", stable.ScopeID, stable.OwnerKind)
	}

	legacy := typesMemoryToRecord(types.MemoryRecord{
		ID:      "legacy",
		Scope:   "character",
		OwnerID: "char_aria",
	})
	if legacy.OwnerKind != "" {
		t.Fatalf("legacy export was reclassified as %q", legacy.OwnerKind)
	}
}
