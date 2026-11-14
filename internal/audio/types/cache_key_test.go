package types

import "testing"

func keyForRequest(req SynthesisRequest) string {
	return GenerateCacheKeyForRequest(req)
}

func TestGenerateCacheKeyForRequestSeparatesMaterialFields(t *testing.T) {
	base := SynthesisRequest{
		Text:          "hello",
		Engine:        "kokoro",
		ModelIdentity: "model-a",
		Voice:         "af_heart",
		Language:      "en",
		VoicePath:     "voices-a.bin",
		Speed:         1.0,
	}
	baseKey := keyForRequest(base)
	if baseKey == "" {
		t.Fatal("cache key is empty")
	}

	tests := []struct {
		name   string
		mutate func(*SynthesisRequest)
	}{
		{name: "text", mutate: func(r *SynthesisRequest) { r.Text = "hello " }},
		{name: "engine", mutate: func(r *SynthesisRequest) { r.Engine = "melotts" }},
		{name: "model", mutate: func(r *SynthesisRequest) { r.ModelIdentity = "model-b" }},
		{name: "voice", mutate: func(r *SynthesisRequest) { r.Voice = "am_echo" }},
		{name: "language", mutate: func(r *SynthesisRequest) { r.Language = "zh" }},
		{name: "voice path", mutate: func(r *SynthesisRequest) { r.VoicePath = "voices-b.bin" }},
		{name: "speed", mutate: func(r *SynthesisRequest) { r.Speed = 1.25 }},
		{name: "sub-microsecond speed", mutate: func(r *SynthesisRequest) { r.Speed = 1.0000001 }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changed := base
			tt.mutate(&changed)
			if got := keyForRequest(changed); got == baseKey {
				t.Fatalf("%s did not participate in the cache key", tt.name)
			}
		})
	}
}

func TestGenerateCacheKeyForRequestDoesNotPartitionByCharacter(t *testing.T) {
	a := SynthesisRequest{
		Text:          "hello",
		Engine:        "piper",
		ModelIdentity: "model-a",
		Voice:         "voice-a",
		Speed:         1,
		CharacterID:   "character-a",
	}
	b := a
	b.CharacterID = "character-b"

	if keyForRequest(a) != keyForRequest(b) {
		t.Fatal("character identity changed an otherwise equivalent cache key")
	}
}

func TestTTSModelIdentitySeparatesEngineAssets(t *testing.T) {
	base := TTSConfig{
		ModelPath:  "model-a.onnx",
		VoicesPath: "voices-a.bin",
		ConfigPath: "config-a.json",
		ScriptPath: "bridge-a.py",
		Device:     "cpu",
	}
	baseIdentity := TTSModelIdentity("kokoro", base, "voice-a")
	for name, mutate := range map[string]func(*TTSConfig){
		"model":      func(c *TTSConfig) { c.ModelPath = "model-b.onnx" },
		"model name": func(c *TTSConfig) { c.ModelName = "tts-2" },
		"voices":     func(c *TTSConfig) { c.VoicesPath = "voices-b.bin" },
		"config":     func(c *TTSConfig) { c.ConfigPath = "config-b.json" },
		"script":     func(c *TTSConfig) { c.ScriptPath = "bridge-b.py" },
		"device":     func(c *TTSConfig) { c.Device = "cuda" },
	} {
		t.Run(name, func(t *testing.T) {
			changed := base
			mutate(&changed)
			if got := TTSModelIdentity("kokoro", changed, "voice-a"); got == baseIdentity {
				t.Fatalf("%s did not participate in the model identity", name)
			}
		})
	}
	if got := TTSModelIdentity("kokoro", base, "voice-b"); got == baseIdentity {
		t.Fatal("voice path did not participate in the model identity")
	}
}

func TestGenerateCacheKeyForRequestIsVersioned(t *testing.T) {
	req := SynthesisRequest{
		Text:          "hello",
		Engine:        "piper",
		ModelIdentity: "model-a",
		Voice:         "voice-a",
		Speed:         1,
	}
	legacy := GenerateCacheKey(req.Text, req.Voice, req.Speed)
	if keyForRequest(req) == legacy {
		t.Fatal("request cache key matches the legacy text/voice/speed key")
	}
	if keyForRequest(req) != keyForRequest(req) {
		t.Fatal("equivalent requests produced different cache keys")
	}
}
