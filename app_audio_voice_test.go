package main

import (
	"strings"
	"testing"
)

func TestParseVoiceOverrideJSONRejectsMalformedPayload(t *testing.T) {
	if _, err := parseVoiceOverrideJSON(`{"engine":`); err == nil {
		t.Fatal("malformed voice override was accepted")
	}
}

func TestParseVoiceOverrideJSONPreservesCharacterIdentityAndInheritance(t *testing.T) {
	override, err := parseVoiceOverrideJSON(`{"character_id":"char_a","engine":"","voice":""}`)
	if err != nil {
		t.Fatalf("parseVoiceOverrideJSON: %v", err)
	}
	if override == nil {
		t.Fatal("explicit inheritance payload became nil")
	}
	if override.CharacterID != "char_a" || override.Engine != "" || override.Voice != "" {
		t.Fatalf("override = %+v, want character_a with blank inheritance fields", override)
	}
}

func TestParseVoiceOverrideJSONTreatsEmptyPayloadAsNoOverride(t *testing.T) {
	override, err := parseVoiceOverrideJSON("")
	if err != nil {
		t.Fatalf("parseVoiceOverrideJSON: %v", err)
	}
	if override != nil {
		t.Fatalf("empty payload = %+v, want nil", override)
	}
}

func TestParseVoiceOverrideJSONErrorNamesTheBoundary(t *testing.T) {
	_, err := parseVoiceOverrideJSON(`{"engine":`)
	if err == nil || !strings.Contains(err.Error(), "voice override JSON") {
		t.Fatalf("error = %v, want voice override JSON context", err)
	}
}
