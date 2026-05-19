package entities

type World struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	KeyFacts    []string `json:"key_facts"`
	References  []string `json:"references"`
}

type Scenario struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	KeyDetails  []string `json:"key_details"`
	References  []string `json:"references"`
}

// VoiceConfig stores per-character TTS voice overrides.
// Null/empty fields inherit the global application default.
type VoiceConfig struct {
	Enabled         bool    `json:"enabled"`
	Engine          string  `json:"engine,omitempty"`           // "piper", "melotts", "browser", "" = inherit
	Voice           string  `json:"voice,omitempty"`            // e.g. "en_US-lessac-medium"
	Speed           float64 `json:"speed,omitempty"`            // 0.5 - 2.0
	Volume          float64 `json:"volume,omitempty"`           // 0.0 - 1.0
	NoiseScale      float64 `json:"noise_scale,omitempty"`      // 0.0 - 1.0
	NoiseW          float64 `json:"noise_w,omitempty"`          // 0.0 - 1.0
	SentenceSilence float64 `json:"sentence_silence,omitempty"` // seconds
}

type Character struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Personality   string   `json:"personality"`
	Backstory     string   `json:"backstory"`
	Goals         []string `json:"goals"`
	Relationships []string `json:"relationships"`
	References    []string `json:"references"`
	// Metadata fields hoisted from the frontend metadata{} sub-object by ParseCharacterContextFromMap
	Age        string `json:"age,omitempty"`
	Role       string `json:"role,omitempty"`
	Background string `json:"background,omitempty"`
	Priority   int    `json:"priority,omitempty"`
	Version    int    `json:"version,omitempty"`
	// Per-character TTS voice override (map for Wails binding compatibility)
	VoiceConfig map[string]interface{} `json:"voice_config,omitempty"`
}
