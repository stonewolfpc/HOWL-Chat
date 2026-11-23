package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"howl-chat/internal/audio"
	"howl-chat/internal/audio/capability"
	"howl-chat/internal/audio/tts/kokoro"
	audiotypes "howl-chat/internal/audio/types"
	"howl-chat/internal/image/registry"
)

// TTSVoiceOption is one selectable voice plus its actual availability. A voice
// that is listed by a profile but not installed is returned with Available
// false and a detail explaining what is missing; the UI must not present it as
// a working choice.
type TTSVoiceOption struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Language    string `json:"language,omitempty"`
	Description string `json:"description,omitempty"`
	Available   bool   `json:"available"`
	Detail      string `json:"detail,omitempty"`
}

// TTSVoiceCatalog is the engine-aware voice menu used by the Character tab.
// Inherited is true when the caller asked for the global engine rather than an
// explicit per-character engine.
type TTSVoiceCatalog struct {
	Engine         string           `json:"engine"`
	Inherited      bool             `json:"inherited"`
	Available      bool             `json:"available"`
	DefaultVoice   string           `json:"default_voice,omitempty"`
	SpeedSupported bool             `json:"speed_supported"`
	MinSpeed       float64          `json:"min_speed"`
	MaxSpeed       float64          `json:"max_speed"`
	Detail         string           `json:"detail,omitempty"`
	Voices         []TTSVoiceOption `json:"voices"`
}

// GetTTSVoiceCatalog returns the real voices for one engine. A blank engine
// means "inherit global": the global engine is resolved first, then its
// catalog is returned with Inherited set. An auto global engine cannot provide
// a truthful per-character catalog and is reported as unavailable.
func (a *App) GetTTSVoiceCatalog(engine string) (TTSVoiceCatalog, error) {
	engine = strings.ToLower(strings.TrimSpace(engine))
	inherited := engine == ""
	if inherited {
		engine = strings.ToLower(strings.TrimSpace(a.audio.Desired().TTSEngine))
	}
	if engine == "" || engine == "auto" {
		return TTSVoiceCatalog{
			Engine:    "auto",
			Inherited: inherited,
			Detail:    "Inherit Global is selected, but the global TTS engine is auto. Choose an explicit engine before selecting a voice.",
			Voices:    []TTSVoiceOption{},
		}, nil
	}
	if engine == "browser" {
		return TTSVoiceCatalog{
			Engine:    engine,
			Inherited: inherited,
			Available: true,
			Detail:    "Browser TTS uses the WebView's system voices and the Windows default output. It cannot target a specific output device.",
			Voices:    []TTSVoiceOption{},
		}, nil
	}

	profile, ok := audiotypes.GetTTSProfile(audiotypes.TTSModelType(engine))
	if !ok {
		return TTSVoiceCatalog{}, fmt.Errorf("unsupported TTS engine %q", engine)
	}
	catalog := TTSVoiceCatalog{
		Engine:         engine,
		Inherited:      inherited,
		SpeedSupported: profile.SupportsSpeed,
		MinSpeed:       profile.MinSpeed,
		MaxSpeed:       profile.MaxSpeed,
		Voices:         []TTSVoiceOption{},
	}

	switch engine {
	case "piper":
		catalog = buildPiperVoiceCatalog(catalog, a.registeredPiperEntries(), *profile)
	case "kokoro":
		catalog = buildKokoroVoiceCatalog(catalog, *profile)
	case "melotts":
		catalog = buildMeloVoiceCatalog(catalog, *profile)
	default:
		catalog.Detail = "No voice catalog is implemented for this engine."
	}
	return catalog, nil
}

func (a *App) registeredPiperEntries() []registry.Entry {
	if a.registries.Models() == nil {
		return nil
	}
	return a.registries.Models().ByTypeInDirs("tts", registry.KindTTS)
}

func buildPiperVoiceCatalog(catalog TTSVoiceCatalog, entries []registry.Entry, profile audiotypes.TTSModelProfile) TTSVoiceCatalog {
	defaultPath := audio.TTSModelPathFor(audiotypes.TTSPiper)
	paths := make(map[string]string)
	for _, entry := range entries {
		paths[entry.Name] = entry.Path
	}
	for _, voice := range profile.SupportedVoices {
		if _, ok := paths[voice]; !ok {
			paths[voice] = filepath.Join("models", voice+".onnx")
		}
	}
	if defaultPath != "" {
		name := strings.TrimSuffix(filepath.Base(defaultPath), filepath.Ext(defaultPath))
		if name != "" {
			paths[name] = defaultPath
		}
	}

	for name, path := range paths {
		available := fileExists(path)
		detail := ""
		if !available {
			detail = "Voice model is missing: " + path
		}
		catalog.Voices = append(catalog.Voices, TTSVoiceOption{
			ID:        name,
			Name:      name,
			Language:  voiceLanguage(name),
			Available: available,
			Detail:    detail,
		})
	}
	sort.Slice(catalog.Voices, func(i, j int) bool {
		if catalog.Voices[i].Available != catalog.Voices[j].Available {
			return catalog.Voices[i].Available
		}
		return catalog.Voices[i].Name < catalog.Voices[j].Name
	})
	for _, voice := range catalog.Voices {
		if voice.Available {
			catalog.Available = true
			catalog.DefaultVoice = voice.ID
			break
		}
	}
	if catalog.Available {
		catalog.Detail = "Piper voices are installed. Piper binary availability is checked when synthesis starts."
	} else {
		catalog.Detail = "No Piper voice models are installed. Add a directory in Audio Settings or set HOWL_TTS_PIPER_MODEL."
	}
	return catalog
}

func buildKokoroVoiceCatalog(catalog TTSVoiceCatalog, profile audiotypes.TTSModelProfile) TTSVoiceCatalog {
	modelPath := audio.TTSModelPathFor(audiotypes.TTSKokoro)
	voicesPath := audio.KokoroVoicesPathFor()
	configPath := audio.KokoroConfigPathFor()
	scriptPath := audio.KokoroScriptPathFor()

	missing := make([]string, 0, 4)
	for label, path := range map[string]string{
		"model": modelPath, "voices archive": voicesPath, "config": configPath, "bridge script": scriptPath,
	} {
		if path == "" || !fileExists(path) {
			missing = append(missing, label+" ("+emptyPath(path)+")")
		}
	}
	sort.Strings(missing)
	assetsAvailable := len(missing) == 0

	voices, err := kokoro.ListVoicesFromArchive(voicesPath)
	if err != nil || len(voices) == 0 {
		if err != nil {
			catalog.Detail = "Kokoro voices archive could not be read: " + err.Error()
		} else {
			catalog.Detail = "Kokoro voices archive contains no voices."
		}
		if len(missing) > 0 {
			catalog.Detail += " Missing: " + strings.Join(missing, ", ") + "."
		}
		return catalog
	}

	for _, voice := range voices {
		catalog.Voices = append(catalog.Voices, TTSVoiceOption{
			ID:        voice,
			Name:      voice,
			Available: assetsAvailable,
			Detail:    kokoroAvailabilityDetail(assetsAvailable, missing),
		})
	}
	if assetsAvailable {
		catalog.Available = true
		catalog.DefaultVoice = voices[0]
		catalog.Detail = "Kokoro assets are present. Python is checked when synthesis starts."
	} else {
		catalog.Detail = "Kokoro assets are incomplete. Missing: " + strings.Join(missing, ", ") + "."
	}
	_ = profile
	return catalog
}

func buildMeloVoiceCatalog(catalog TTSVoiceCatalog, profile audiotypes.TTSModelProfile) TTSVoiceCatalog {
	available := false
	detail := ""
	record, err := capability.Select(audioCapabilityRecords(), capability.TaskSpeechSynthesis, string(audiotypes.TTSMeloTTS))
	if err != nil {
		detail = "MeloTTS is unavailable: " + err.Error()
	} else {
		available = true
		detail = "MeloTTS is reachable. Availability is a service check; synthesis is not verified by the catalog."
		_ = record
	}
	for _, voice := range profile.SupportedVoices {
		catalog.Voices = append(catalog.Voices, TTSVoiceOption{
			ID:        voice,
			Name:      voice,
			Language:  strings.ToLower(voice),
			Available: available,
			Detail:    detail,
		})
	}
	catalog.Available = available
	catalog.Detail = detail
	if available && len(profile.SupportedVoices) > 0 {
		catalog.DefaultVoice = profile.SupportedVoices[0]
	}
	return catalog
}

func kokoroAvailabilityDetail(available bool, missing []string) string {
	if available {
		return "Kokoro assets are present."
	}
	return "Missing: " + strings.Join(missing, ", ")
}

func emptyPath(path string) string {
	if path == "" {
		return "not resolved"
	}
	return path
}

func fileExists(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func voiceLanguage(id string) string {
	base := id
	if i := strings.IndexByte(base, '-'); i > 0 {
		base = base[:i]
	}
	if strings.Contains(base, "_") {
		return strings.ToLower(strings.ReplaceAll(base, "_", "-"))
	}
	return ""
}
