package memory_test

import (
	"reflect"
	"sort"
	"testing"

	"howl-chat/internal/backend/chat"
	"howl-chat/internal/backend/types/entities"
)

// g1Envelope is the logical capture contract exercised by the fixture. It is
// test-only: no production type or path is changed.
type g1Envelope struct {
	SourceMessageID   string
	SourceSpan        string
	SpeakerRole       string
	SourceCharacterID *string // nil = missing; pointer to "" = explicit none
	ChatID            string
	TurnID            string
	EventType         string
	Evidence          string
	WorldID           string
	ScenarioID        string
	Attribution       string
	AuthoredAssetRef  string
	Content           string
	CaptureState      string
	MissingFields     []string
	Determination     string
	EligibleClasses   []string
	NonClassProcess   bool
	AssignedClass     string
}

type g1Case struct {
	name                string
	data                chat.ChatData
	attribution         string
	authoredAssetRef    string
	determination       string
	eventType           string
	evidence            string
	wantState           string
	wantMissing         []string
	wantEligible        []string
	wantSourceCharacter *string
	// referenceClass is deliberately separate from the source facts and from
	// EligibleClasses. It is empty throughout G1: the fixture must not encode
	// the class it expects the classifier to reproduce.
	referenceClass string
}

func strPtr(s string) *string { return &s }

func g1BaseData(h chat.ChatHistory) chat.ChatData {
	return chat.ChatData{
		Meta: chat.ChatMeta{
			ID:               "chat-1",
			WorldID:          "world-1",
			ScenarioID:       "scenario-1",
			ActiveCharacters: []string{"char-a", "char-b"},
		},
		History: []chat.ChatHistory{h},
		World: &chat.WorldContext{
			World: entities.World{ID: "world-1", Name: "North"},
		},
		Scenario: &chat.ScenarioContext{
			Scenario: entities.Scenario{ID: "scenario-1", Name: "Winter"},
		},
		Characters: []*chat.CharacterContext{
			{Character: entities.Character{ID: "char-a", Name: "Aria"}},
			{Character: entities.Character{ID: "char-b", Name: "Aria"}},
		},
	}
}

func sortedCopy(in []string) []string {
	out := append([]string(nil), in...)
	sort.Strings(out)
	return out
}

func sameStrings(a, b []string) bool {
	return reflect.DeepEqual(sortedCopy(a), sortedCopy(b))
}

func sameStringPtr(a, b *string) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

// buildG1Envelope models the contract from existing app data. It is test-only
// and does not call or change the production memory API.
func buildG1Envelope(data chat.ChatData, attribution, authoredAssetRef, determination, eventType, evidence string) g1Envelope {
	h := data.History[0]
	missing := map[string]bool{}

	add := func(field, value string) {
		if value == "" {
			missing[field] = true
		}
	}

	add("source_message_id", h.ID)
	add("source_span", h.Text)
	add("speaker_role", h.Role)
	add("chat_id", data.Meta.ID)
	add("event_type", eventType)
	add("evidence", evidence)
	add("world_id", data.Meta.WorldID)
	add("scenario_id", data.Meta.ScenarioID)
	add("content", h.Text)

	var sourceCharacterID *string
	switch h.Role {
	case "assistant":
		if h.CharacterID != "" {
			sourceCharacterID = strPtr(h.CharacterID)
		} else {
			missing["source_character_id"] = true
		}
	case "user", "system", "narrator":
		sourceCharacterID = strPtr("") // explicit none
	default:
		// speaker_role is already missing; do not invent an identity.
	}

	if determination == "experience_belief" && attribution == "" {
		missing["attribution"] = true
	}

	missingFields := make([]string, 0, len(missing))
	for field := range missing {
		missingFields = append(missingFields, field)
	}
	sort.Strings(missingFields)

	state := "complete"
	if len(missingFields) > 0 {
		state = "incomplete"
	}

	eligible := []string{}
	if state == "complete" {
		switch determination {
		case "experience_belief":
			if sourceCharacterID != nil && *sourceCharacterID != "" {
				eligible = append(eligible, "experiential", "character_specific")
			}
		case "ambient_restricted":
			if authoredAssetRef != "" {
				eligible = append(eligible, "ambient", "restricted")
			}
		}
	}
	sort.Strings(eligible)

	return g1Envelope{
		SourceMessageID:   h.ID,
		SourceSpan:        h.Text,
		SpeakerRole:       h.Role,
		SourceCharacterID: sourceCharacterID,
		ChatID:            data.Meta.ID,
		TurnID:            "0",
		EventType:         eventType,
		Evidence:          evidence,
		WorldID:           data.Meta.WorldID,
		ScenarioID:        data.Meta.ScenarioID,
		Attribution:       attribution,
		AuthoredAssetRef:  authoredAssetRef,
		Content:           h.Text,
		CaptureState:      state,
		MissingFields:     missingFields,
		Determination:     determination,
		EligibleClasses:   eligible,
		NonClassProcess:   true,
		AssignedClass:     "",
	}
}

func g1Cases() []g1Case {
	assistant := func(id, text string) chat.ChatData {
		return g1BaseData(chat.ChatHistory{
			ID:          "msg-1",
			Sender:      "Aria",
			Text:        text,
			Role:        "assistant",
			CharacterID: id,
		})
	}
	user := g1BaseData(chat.ChatHistory{
		ID: "msg-user", Sender: "User", Text: "I remember the northern pass.", Role: "user",
	})
	system := g1BaseData(chat.ChatHistory{
		ID: "msg-system", Sender: "Narrator", Text: "The northern gate is sealed.", Role: "system",
	})

	f10 := assistant("char-a", "The northern pass is closed.")
	f10.Meta.WorldID = ""
	f10.Meta.ScenarioID = ""

	return []g1Case{
		{
			name: "F1 stable character message", data: assistant("char-a", "I feel safe here."),
			attribution: "observed", determination: "experience_belief", eventType: "emotion",
			evidence: "feel safe", wantState: "complete",
			wantEligible: []string{"character_specific", "experiential"},
			wantSourceCharacter: strPtr("char-a"),
		},
		{
			name: "F2 user message explicit none", data: user,
			determination: "ownership", eventType: "custom", evidence: "remember",
			wantState: "complete", wantEligible: []string{}, wantSourceCharacter: strPtr(""),
		},
		{
			name: "F3 narrator message", data: system,
			determination: "ambient_restricted", eventType: "custom", evidence: "sealed",
			wantState: "complete", wantEligible: []string{}, wantSourceCharacter: strPtr(""),
		},
		{
			name: "F4 empty character id", data: assistant("", "I feel safe here."),
			attribution: "observed", determination: "experience_belief", eventType: "emotion",
			evidence: "feel safe", wantState: "incomplete",
			wantMissing: []string{"source_character_id"}, wantEligible: []string{}, wantSourceCharacter: nil,
		},
		{
			name: "F5 duplicate display name stable ids", data: assistant("char-b", "We shared this."),
			attribution: "observed", determination: "experience_belief", eventType: "shared",
			evidence: "shared", wantState: "complete",
			wantEligible: []string{"character_specific", "experiential"},
			wantSourceCharacter: strPtr("char-b"),
		},
		{
			name: "F6 missing attribution", data: assistant("char-a", "I grew up in the north."),
			determination: "experience_belief", eventType: "custom", evidence: "grew up",
			wantState: "incomplete", wantMissing: []string{"attribution"},
			wantEligible: []string{}, wantSourceCharacter: strPtr("char-a"),
		},
		{
			name: "F7 authored world definition", data: system,
			authoredAssetRef: "asset:world-1", determination: "ambient_restricted",
			eventType: "custom", evidence: "authored", wantState: "complete",
			wantEligible: []string{"ambient", "restricted"}, wantSourceCharacter: strPtr(""),
		},
		{
			name: "F8 world-related character speech", data: assistant("char-a", "The northern kingdom is ruled by Queen Elara."),
			attribution: "observed", determination: "experience_belief", eventType: "discovery",
			evidence: "northern kingdom", wantState: "complete",
			wantEligible: []string{"character_specific", "experiential"},
			wantSourceCharacter: strPtr("char-a"),
		},
		{
			name: "F9 legacy record without source", data: g1BaseData(chat.ChatHistory{
				Text: "legacy memory with no source", Role: "",
			}),
			determination: "ownership", eventType: "custom", evidence: "legacy",
			wantState: "incomplete", wantMissing: []string{"source_message_id", "speaker_role"},
			wantEligible: []string{}, wantSourceCharacter: nil,
		},
		{
			name: "F10 missing world and scenario", data: f10,
			attribution: "observed", determination: "experience_belief", eventType: "emotion",
			evidence: "northern pass", wantState: "incomplete",
			wantMissing: []string{"scenario_id", "world_id"}, wantEligible: []string{},
			wantSourceCharacter: strPtr("char-a"),
		},
		{
			name: "character tell", data: assistant("char-a", "I grew up in the northern mountains."),
			attribution: "told", determination: "experience_belief", eventType: "shared",
			evidence: "grew up", wantState: "complete",
			wantEligible: []string{"character_specific", "experiential"},
			wantSourceCharacter: strPtr("char-a"),
		},
	}
}

func TestG1CaptureContract(t *testing.T) {
	for _, tc := range g1Cases() {
		t.Run(tc.name, func(t *testing.T) {
			got := buildG1Envelope(tc.data, tc.attribution, tc.authoredAssetRef, tc.determination, tc.eventType, tc.evidence)

			if got.CaptureState != tc.wantState {
				t.Fatalf("capture state = %q, want %q (missing %v)", got.CaptureState, tc.wantState, got.MissingFields)
			}
			if !sameStrings(got.MissingFields, tc.wantMissing) {
				t.Fatalf("missing fields = %v, want %v", got.MissingFields, tc.wantMissing)
			}
			if !sameStrings(got.EligibleClasses, tc.wantEligible) {
				t.Fatalf("eligible classes = %v, want %v", got.EligibleClasses, tc.wantEligible)
			}
			if !sameStringPtr(got.SourceCharacterID, tc.wantSourceCharacter) {
				t.Fatalf("source character = %v, want %v", got.SourceCharacterID, tc.wantSourceCharacter)
			}
			if got.AssignedClass != "" {
				t.Fatalf("G1 must not assign a class, got %q", got.AssignedClass)
			}
			if tc.referenceClass != "" {
				t.Fatalf("G1 fixture must not encode a reference class, got %q", tc.referenceClass)
			}
			if !got.NonClassProcess {
				t.Fatal("non-class processing must remain permitted")
			}
		})
	}
}

func TestG1ExplicitNoneIsNotMissing(t *testing.T) {
	var f2, f4 g1Case
	for _, tc := range g1Cases() {
		switch tc.name {
		case "F2 user message explicit none":
			f2 = tc
		case "F4 empty character id":
			f4 = tc
		}
	}
	if f2.data.History[0].CharacterID != "" || f4.data.History[0].CharacterID != "" {
		t.Fatal("fixture setup error: both cases should have an empty CharacterID")
	}

	got2 := buildG1Envelope(f2.data, f2.attribution, f2.authoredAssetRef, f2.determination, f2.eventType, f2.evidence)
	got4 := buildG1Envelope(f4.data, f4.attribution, f4.authoredAssetRef, f4.determination, f4.eventType, f4.evidence)

	if got2.SourceCharacterID == nil || *got2.SourceCharacterID != "" {
		t.Fatalf("F2 must carry an explicit none, got %v", got2.SourceCharacterID)
	}
	if got2.CaptureState != "complete" {
		t.Fatalf("F2 must be complete, got %q (%v)", got2.CaptureState, got2.MissingFields)
	}
	if got4.SourceCharacterID != nil {
		t.Fatalf("F4 must carry a missing source character, got %v", got4.SourceCharacterID)
	}
	if got4.CaptureState != "incomplete" {
		t.Fatalf("F4 must be incomplete, got %q", got4.CaptureState)
	}
}

func TestG1WorldSpeechIsNotAmbient(t *testing.T) {
	var f8 g1Case
	for _, tc := range g1Cases() {
		if tc.name == "F8 world-related character speech" {
			f8 = tc
			break
		}
	}
	got := buildG1Envelope(f8.data, f8.attribution, f8.authoredAssetRef, f8.determination, f8.eventType, f8.evidence)
	for _, class := range got.EligibleClasses {
		if class == "ambient" || class == "restricted" {
			t.Fatalf("F8 must not be eligible for %s without an authored asset", class)
		}
	}
	if !sameStrings(got.EligibleClasses, []string{"character_specific", "experiential"}) {
		t.Fatalf("F8 eligible = %v, want experiential/character_specific", got.EligibleClasses)
	}
}

func TestG1CharacterTellRemainsUnclassified(t *testing.T) {
	var ct g1Case
	for _, tc := range g1Cases() {
		if tc.name == "character tell" {
			ct = tc
			break
		}
	}
	got := buildG1Envelope(ct.data, ct.attribution, ct.authoredAssetRef, ct.determination, ct.eventType, ct.evidence)
	if got.CaptureState != "complete" {
		t.Fatalf("character-tell capture = %q, want complete", got.CaptureState)
	}
	if got.Attribution != "told" {
		t.Fatalf("character-tell attribution = %q, want told", got.Attribution)
	}
	if got.AssignedClass != "" {
		t.Fatalf("character-tell must not assign a class, got %q", got.AssignedClass)
	}
	if !got.NonClassProcess {
		t.Fatal("character-tell must permit non-class processing")
	}
}
