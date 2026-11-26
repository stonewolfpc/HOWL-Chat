package main

import (
	"context"
	"fmt"

	"howl-chat/internal/audio"
	"howl-chat/internal/backend/chat"
	"howl-chat/internal/backend/llama"
	"howl-chat/internal/backend/memory"
	"howl-chat/internal/backend/model"
	"howl-chat/internal/backend/rag"
	"howl-chat/internal/host"
	"howl-chat/internal/image/registry"
	"howl-chat/internal/media"
	"howl-chat/internal/paths"
	"howl-chat/internal/registries"
	appruntime "howl-chat/internal/runtime"
	appsettings "howl-chat/internal/settings"
)

// bootStages is the application's one ordered composition.
//
// Each stage constructs or wires exactly one owner, in dependency order, and
// Level 0 (internal/host) runs the list, decides readiness and reports which
// stage failed. App supplies the stage bodies because it is the composition
// root that knows which domains exist; it does not decide their order, and it
// stores none of what they build beyond the owner reference itself.
//
// Ordering constraints that are real, not stylistic:
//   - the audio owner exists before the settings stage applies its desired
//     contract, because the contract is stored on that owner;
//   - the registry owner exists before the image domain boots, because the
//     domain asks the model registry to record its models directory;
//   - the registry owner exists before chat is given the content libraries.
func (a *App) bootStages() []host.Stage {
	return []host.Stage{
		{Name: "core services", Run: a.bootCoreServices},
		{Name: "audio owner", Run: a.bootAudioOwner},
		{Name: "settings store", Run: a.bootSettings},
		{Name: "registries", Run: a.bootRegistries},
		{Name: "image domain", Run: a.bootImageDomain},
		{Name: "inference capture", Run: a.bootInferenceCapture},
	}
}

// bootCoreServices builds the inference client, the model manager and the chat,
// memory and RAG services those depend on.
func (a *App) bootCoreServices(context.Context) error {
	client := llama.NewHTTPClient("localhost", 8080, true)
	models := model.NewManager(client)
	a.engine.SetClient(client)
	a.engine.SetModels(models)

	persistence, err := chat.NewPersistenceManager()
	if err != nil {
		fmt.Printf("Warning: Failed to initialize chat persistence: %v\n", err)
		persistence = nil
	}
	a.conversation = chat.NewActiveConversation(persistence)

	a.memoryService = memory.NewService("default") // Default chat ID for initial setup

	service, err := chat.NewService(models, client)
	if err != nil {
		fmt.Printf("Warning: Failed to initialize chat service: %v\n", err)
		service = nil
	}
	a.chatService = service

	// RAG uses the llama-server embedding endpoint when a client is present.
	var embedder rag.Embedder = rag.NewKeywordEmbedder()
	if client != nil {
		embedder = rag.NewLlamaEmbedder(client)
	}
	a.ragService = rag.NewService(paths.StateRoot(), embedder)
	if a.chatService != nil {
		a.chatService.SetRAGService(a.ragService)
		a.chatService.SetRAGEnabled(true)
	}
	return nil
}

// bootAudioOwner creates the audio domain. The gate is the one admission route;
// the probe never opens a capture stream; failures escalate through the
// application logger.
func (a *App) bootAudioOwner(context.Context) error {
	a.audio = media.NewAudio(audio.NewProductionManager(a.level0().Processes.Registry()),
		func(operation string) error { return a.level0().Gate.Ensure(operation) },
		a.captureProbe,
		a.reportAudioFailure)
	return nil
}

// bootSettings constructs the settings authority, restores the unified document
// and its sampler/model sections (falling back to the legacy section files), and
// applies the desired audio contract without activating anything.
func (a *App) bootSettings(context.Context) error {
	dir := a.getSettingsDir()
	a.settings = appruntime.NewSettings(appsettings.NewStore(dir), dir, defaultSamplerSettings())
	if err := a.settings.Load(); err != nil {
		fmt.Printf("Warning: failed to load settings: %v\n", err)
	}

	a.audio.SetDesired(audio.DefaultDesiredSettings())
	if rawAudio := a.settings.Section("audio"); len(rawAudio) > 0 {
		a.audio.SetDesired(audio.DesiredSettingsFromMap(rawAudio))
	}
	return nil
}

// bootRegistries builds the one registry owner: content roots, inventory,
// world/scenario import services, canonical lorebooks, the model-profile cache
// and the model directory/type registry. A corrupt store is reported and left
// absent; the remaining authorities stay available.
func (a *App) bootRegistries(context.Context) error {
	regs, errs := registries.New(a.getSettingsDir())
	for _, err := range errs {
		fmt.Printf("Warning: %v\n", err)
	}
	a.registries = regs
	if a.chatService != nil && regs.Libraries() != nil {
		if err := a.chatService.SetContentLibraries(regs.Libraries()); err != nil {
			fmt.Printf("Warning: failed to attach world/scenario store: %v\n", err)
		}
	}
	return nil
}

// bootImageDomain creates the image owner and lets it construct its own engine,
// gallery, job manager and last-model restore. App no longer builds any part of
// that session; it only supplies the directories and the registry hook.
func (a *App) bootImageDomain(context.Context) error {
	startupMark("image init: begin")
	a.images = media.NewImage(a.getSettingsDir())
	err := a.images.Boot(media.BootOptions{
		TempDir: paths.TempDir(),
		RegisterModelsDir: func(dir string) error {
			models := a.registries.Models()
			if models == nil {
				return fmt.Errorf("model registry not loaded")
			}
			return models.RegisterDir(dir, registry.KindImage)
		},
		Logf:            func(format string, args ...interface{}) { fmt.Printf(format+"\n", args...) },
		ProcessRegistry: a.level0().Processes.Registry(),
		OnEngineExit:    a.unregisterImageModels,
	})
	if err == nil && a.images != nil {
		if jobs := a.images.Jobs(); jobs != nil {
			jobs.SetOnFinished(a.reclaimImageEngineAfterGeneration)
		}
	}
	startupMark("image init: done")
	return err
}

// bootInferenceCapture wraps the HTTP client so tests can sit at the same call
// point in production builds. It runs last: every service already holds the
// client it will use.
func (a *App) bootInferenceCapture(context.Context) error {
	a.engine.SetClient(wrapInferenceCapture(a.engine.Client()))
	return nil
}

func defaultSamplerSettings() map[string]interface{} {
	return map[string]interface{}{
		"temperature":                0.7,
		"top_p_enabled":              true,
		"top_p":                      0.9,
		"top_k":                      40,
		"min_p":                      0.05,
		"repeat_penalty":             1.1,
		"repeat_last_n":              64,
		"frequency_penalty_enabled":  true,
		"frequency_penalty":          0.0,
		"presence_penalty_enabled":   true,
		"presence_penalty":           0.0,
		"typical_p_enabled":          true,
		"typical_p":                  1.0,
		"mirostat_enabled":           true,
		"mirostat":                   0,
		"mirostat_tau":               5.0,
		"mirostat_eta":               0.1,
		"dynamic_temp_range_enabled": true,
		"dynamic_temp_range":         0.0,
		"dynamic_temp_exponent":      1.0,
		"dry_multiplier":             0.0,
		"dry_allowed_length":         2,
		"dry_base":                   1.0,
		"smoothing_factor":           0.0,
		"smoothing_curve":            1.0,
		"top_a_enabled":              true,
		"top_a":                      0.0,
		"epsilon_cutoff":             0.0,
		"eta_cutoff":                 0.0,
		"encoder_repeat_penalty":     1.0,
		"no_repeat_ngram_size":       0,
		"seed":                       -1,
	}
}
