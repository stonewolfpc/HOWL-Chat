package main

import (
	"os"
	"regexp"
	"testing"

	"howl-chat/internal/audio/capability"
	audiotypes "howl-chat/internal/audio/types"
)

// The status the settings UI reads must carry the same capability contract that
// decides availability, including before any audio manager exists.
func TestGetAudioStatusReportsCapabilityContract(t *testing.T) {
	app := &App{}
	status := app.GetAudioStatus()

	if status["asr_ready"] != false || status["tts_ready"] != false {
		t.Fatalf("status = %+v, want nothing reported ready before the manager exists", status)
	}

	records, ok := status["capabilities"].([]capability.Record)
	if !ok || len(records) == 0 {
		t.Fatalf("capabilities = %T, want a populated capability contract", status["capabilities"])
	}

	tasks := make(map[capability.Task]bool)
	for _, record := range records {
		tasks[record.Task] = true
		if record.Task == capability.TaskStreamingTranscription && record.Selectable() {
			t.Fatalf("%q advertises streaming", record.ID)
		}
		if record.Verified() && record.ID != string(audiotypes.ASRQwen3ASR17BLocal) {
			t.Fatalf("%q claims end-to-end verification without a run", record.ID)
		}
	}
	for _, want := range capability.Tasks() {
		if !tasks[want] {
			t.Fatalf("status contract is missing task %q", want)
		}
	}
}

// DEBT-T3: the engine list the UI offers must be exactly the registry, so the
// UI cannot present an identity the resolver rejects without knowing it.
func TestSettingsUIOffersOnlyRegistryASREngines(t *testing.T) {
	data, err := os.ReadFile("frontend/dist/audio-settings.html")
	if err != nil {
		t.Fatalf("read settings UI: %v", err)
	}

	selectBlock := regexp.MustCompile(`(?s)id="asr-engine".*?</select>`).FindString(string(data))
	if selectBlock == "" {
		t.Fatal("no asr-engine select found in the settings UI")
	}

	options := regexp.MustCompile(`<option value="([^"]+)"`).FindAllStringSubmatch(selectBlock, -1)
	if len(options) < 2 {
		t.Fatalf("asr-engine select has %d options, want auto plus registry engines", len(options))
	}

	seenAuto := false
	for _, option := range options {
		value := option[1]
		if value == "auto" {
			seenAuto = true
			continue
		}
		if _, ok := audiotypes.GetASRProfile(audiotypes.ASRModelType(value)); !ok {
			t.Fatalf("settings UI offers %q, which is not a registry model", value)
		}
	}
	if !seenAuto {
		t.Fatal("settings UI has no auto option")
	}
}
