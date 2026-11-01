package memory

import "testing"

func TestParseAutoCaptureConfigReadsSavedValues(t *testing.T) {
	raw := []byte(`{"enabled": false, "min_confidence": 42, "max_per_message": 7, "trigger_words": ["log"]}`)

	config, err := parseAutoCaptureConfig(raw)
	if err != nil {
		t.Fatalf("parseAutoCaptureConfig error: %v", err)
	}
	if config.Enabled {
		t.Error("Enabled = true, want false from the saved config")
	}
	if config.MinConfidence != 42 {
		t.Errorf("MinConfidence = %v, want 42", config.MinConfidence)
	}
	if config.MaxPerMessage != 7 {
		t.Errorf("MaxPerMessage = %d, want 7", config.MaxPerMessage)
	}
	if len(config.TriggerWords) != 1 || config.TriggerWords[0] != "log" {
		t.Errorf("TriggerWords = %v", config.TriggerWords)
	}
}

func TestParseAutoCaptureConfigFillsMissingTriggerWords(t *testing.T) {
	config, err := parseAutoCaptureConfig([]byte(`{"enabled": true}`))
	if err != nil {
		t.Fatalf("parseAutoCaptureConfig error: %v", err)
	}
	if len(config.TriggerWords) == 0 {
		t.Fatal("TriggerWords empty; defaults should fill in")
	}
	if config.MinConfidence != 75 {
		t.Errorf("MinConfidence = %v, want the default 75", config.MinConfidence)
	}
}

func TestParseAutoCaptureConfigRejectsGarbage(t *testing.T) {
	if _, err := parseAutoCaptureConfig([]byte("not json")); err == nil {
		t.Fatal("expected an error for invalid JSON")
	}
}
