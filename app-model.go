package main

import (
	"encoding/json"
	"fmt"
	"howl-chat/internal/paths"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"howl-chat/internal/backend/chatpolicy"
	"howl-chat/internal/backend/gguf"
	"howl-chat/internal/backend/loader"
	"howl-chat/internal/backend/types"
	"howl-chat/internal/errors"
	"howl-chat/internal/resource"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ---------------------------------------------------------------------------
// Model Loading Wails Bindings
// ---------------------------------------------------------------------------

// LoadModel starts llama-server with the given model path.
func (a *App) LoadModel(modelPath string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic in LoadModel: %v", r)
			a.logError(errors.SYSPanicRecovered, errors.Critical, "llm.load", "app-model.go:LoadModel", err)
		}
	}()
	fmt.Printf("[LoadModel] START path=%q\n", modelPath)
	if modelPath == "" {
		fmt.Println("[LoadModel] FAIL: empty path")
		return fmt.Errorf("model path cannot be empty")
	}
	if modelPath == "." || modelPath == ".." {
		fmt.Println("[LoadModel] FAIL: relative path")
		return fmt.Errorf("model path cannot be a relative directory reference: %s", modelPath)
	}
	if info, err := os.Stat(modelPath); err != nil {
		fmt.Printf("[LoadModel] FAIL: stat error: %v\n", err)
		return fmt.Errorf("model path does not exist: %w", err)
	} else if info.IsDir() {
		fmt.Println("[LoadModel] FAIL: path is directory")
		return fmt.Errorf("model path is a directory, not a file: %s", modelPath)
	}
	fmt.Println("[LoadModel] path validated")

	// Admission is decided from the file's own GGUF metadata before any
	// running server is stopped. A rejected asset must not cost the user the
	// model that is currently loaded.
	admission := gguf.AdmitModel(modelPath)
	fmt.Printf("[LoadModel] admission: class=%s arch=%q projector=%q evidence=%v\n",
		admission.Class, admission.Arch, admission.Projector, admission.Evidence)
	if !admission.AllowedAsChat() {
		rec := loader.FailureRecord{
			Stage:     loader.StageAdmission,
			ModelPath: modelPath,
			ModelName: filepath.Base(modelPath),
			Arch:      admission.Arch,
			Detail:    admission.ChatRejectionMessage(),
		}
		if path, writeErr := loader.WriteFailureRecord(a.loaderLogDir(), rec); writeErr == nil {
			rec.RecordPath = path
		} else {
			fmt.Printf("[LoadModel] WARN: could not persist diagnostics: %v\n", writeErr)
		}
		userMsg := rec.UserMessage()
		a.logError(errors.LLMArchRejected, errors.Critical, "llm.load.admission", "app-model.go:LoadModel", fmt.Errorf("%s", userMsg))
		a.loadingProgress.Store(0)
		a.loadingStage.Store("Rejected")
		a.emitEvent("model:error", userMsg)
		fmt.Printf("[LoadModel] FAIL: %s\n", userMsg)
		return fmt.Errorf("%s", userMsg)
	}

	// Architecture admission is not task eligibility. Qwen3-ASR and
	// Ultravox/Voxtral share chat architectures but have different intended
	// tasks (ASR and audio understanding), and their chat path is not
	// implemented or verified. Refuse them before projector resolution and
	// before the running server is stopped.
	eligibility := chatpolicy.Evaluate(modelPath, admission)
	fmt.Printf("[LoadModel] chat eligibility: task=%s eligible=%t basis=%s identity=%q aliases=%v text=%t vision=%s audio=%s\n",
		eligibility.Task, eligibility.Eligible, eligibility.Basis, eligibility.Identity, eligibility.Aliases,
		eligibility.TextChat, eligibility.Vision, eligibility.AudioInput)
	if !eligibility.Eligible {
		rec := loader.FailureRecord{
			Stage:     loader.StageAdmission,
			ModelPath: modelPath,
			ModelName: filepath.Base(modelPath),
			Arch:      admission.Arch,
			Detail:    eligibility.RejectionMessage(modelPath),
		}
		if path, writeErr := loader.WriteFailureRecord(a.loaderLogDir(), rec); writeErr == nil {
			rec.RecordPath = path
		} else {
			fmt.Printf("[LoadModel] WARN: could not persist diagnostics: %v\n", writeErr)
		}
		userMsg := rec.UserMessage()
		a.logError(errors.LLMChatTaskRejected, errors.Critical, "llm.load.task", "app-model.go:LoadModel", fmt.Errorf("%s", userMsg))
		a.loadingProgress.Store(0)
		a.loadingStage.Store("Rejected")
		a.emitEvent("model:error", userMsg)
		fmt.Printf("[LoadModel] FAIL: %s\n", userMsg)
		return fmt.Errorf("%s", userMsg)
	}

	// The companion projector is resolved by exact identity before the current
	// model is stopped, so a pairing failure also leaves the running model
	// alone.
	pairing, pairErr := gguf.ResolveProjector(gguf.PairingRequest{
		ModelPath:        modelPath,
		Arch:             admission.Arch,
		EmbeddingLength:  admission.EmbeddingLength,
		Require:          admission.RequiresProjector(),
		RequiredModality: admission.RequiredModality(),
	})
	if pairErr != nil {
		rec := loader.FailureRecord{
			Stage:     loader.StageProjector,
			ModelPath: modelPath,
			ModelName: filepath.Base(modelPath),
			Arch:      admission.Arch,
			Detail:    pairErr.Error(),
		}
		if path, writeErr := loader.WriteFailureRecord(a.loaderLogDir(), rec); writeErr == nil {
			rec.RecordPath = path
		} else {
			fmt.Printf("[LoadModel] WARN: could not persist diagnostics: %v\n", writeErr)
		}
		userMsg := rec.UserMessage()
		a.logError(errors.LLMProjectorPairing, errors.Critical, "llm.load.projector", "app-model.go:LoadModel", fmt.Errorf("%s", userMsg))
		a.loadingProgress.Store(0)
		a.loadingStage.Store("Rejected")
		a.emitEvent("model:error", userMsg)
		fmt.Printf("[LoadModel] FAIL: %s\n", userMsg)
		return fmt.Errorf("%s", userMsg)
	}
	if pairing.Path != "" {
		fmt.Printf("[LoadModel] projector paired by identity: %s (kind=%s identity=%q projection_dim=%d model_embedding=%d mismatch=%t notes=%v)\n",
			pairing.Path, pairing.Info.Kind, pairing.Identity, pairing.ProjectionDim, pairing.ModelEmbedding, pairing.DimensionMismatch, pairing.Notes)
	} else {
		fmt.Println("[LoadModel] no companion projector required or matched")
	}

	// Prevent redundant reloads: if the same model is already loaded, skip.
	if a.serverProc != nil && a.serverProc.Process != nil {
		fmt.Println("[LoadModel] existing serverProc found")
		if a.httpClient != nil {
			if loaded, _ := a.httpClient.GetLoadedModel(); loaded != nil {
				if data, ok := loaded["data"].([]interface{}); ok && len(data) > 0 {
					if first, ok := data[0].(map[string]interface{}); ok {
						if id, ok := first["id"].(string); ok && id == filepath.Base(modelPath) {
							fmt.Println("[LoadModel] same model already loaded, skipping")
							return nil
						}
					}
				}
			}
		}
		progress := a.loadingProgress.Load()
		if progress > 0 && progress < 1000 {
			fmt.Printf("[LoadModel] load in progress (progress=%d), skipping\n", progress)
			return nil
		}
		fmt.Println("[LoadModel] killing existing server")
		_ = a.serverProc.Process.Kill()
		_ = a.serverProc.Wait()
		fmt.Printf("[LoadModel] stopped existing llama-server (pid %d)\n", a.serverProc.Process.Pid)
	}

	fmt.Println("[LoadModel] checking port 8080")
	if !waitForPortFree("8080", 5*time.Second) {
		a.logError(errors.LLMPortConflict, errors.Critical, "llm.load", "app-model.go:56", fmt.Errorf("port 8080 occupied"))
		fmt.Println("[LoadModel] FAIL: port 8080 occupied")
		return fmt.Errorf("port 8080 is occupied by another application; cannot start llama-server")
	}
	fmt.Println("[LoadModel] port 8080 free")

	fmt.Println("[LoadModel] validating llama-server binary")
	if err := ValidateLlamaServer(); err != nil {
		a.logError(errors.LLMServerNotFound, errors.High, "llm.load", "app-model.go:59", err)
		fmt.Printf("[LoadModel] FAIL: llama-server validation error: %v\n", err)
		return err
	}
	fmt.Println("[LoadModel] llama-server binary validated")

	fmt.Printf("[LoadModel] reading GGUF metadata from: %s\n", modelPath)
	profile, err := gguf.ExtractModelProfile(modelPath)
	if err != nil {
		fmt.Printf("[LoadModel] WARN: GGUF read failed: %v\n", err)
		a.logErrorSilent(errors.LLMGGUFReadFail, errors.Medium, "llm.load", "app-model.go:64", err)
	} else {
		fmt.Printf("[LoadModel] profile: Family=%s Layers=%d Params=%d Size=%.1fMB Quant=%s ArchSpec=%v\n",
			profile.Family, profile.LayerCount, profile.ParameterCount,
			float64(profile.FileSizeBytes)/(1024*1024), profile.FileType,
			profile.ArchitectureSpec != nil)
		if profile.Template != "" {
			if a.modelSettings == nil {
				a.modelSettings = make(map[string]interface{})
			}
			// The template is NOT written into the saved settings here. Doing
			// that made the first family's template stick forever: the second
			// family loaded with if _, exists… still true and was prompted in
			// the wrong format. It is injected per load in spawnSettingsFor().
		}
		fmt.Println("[LoadModel] applying safe defaults")
		a.applySafeDefaults(profile)
		fmt.Println("[LoadModel] safe defaults applied")

		// Projector capability is read from the projector's own GGUF metadata.
		// The companion itself was resolved by identity above; directory order
		// is never consulted. Presence must never enable image input on its own:
		// an audio projector and a vision projector share the mmproj file name
		// (DEBT-T1).
		// A projector declaration is not a verified chat capability. Advertise
		// image/audio input only when the task policy marks it verified; until
		// then the projector remains loaded for diagnostics and future use but
		// the request path does not claim the capability.
		profile.HasVision = eligibility.AllowsVision()
		profile.HasAudioInput = eligibility.AllowsAudio()
		profile.ProjectorKind = string(gguf.ProjectorUnknown)
		if pairing.Path == "" {
			fmt.Println("[LoadModel] no companion projector for this model")
		} else {
			info := pairing.Info
			profile.MMPProjPath = pairing.Path
			profile.ProjectorKind = string(info.Kind)
			fmt.Printf("[LoadModel] projector %s declares %q (declared_vision=%t declared_audio=%t advertised_vision=%t advertised_audio=%t dim=%d type=%s",
				pairing.Path, info.Kind, info.HasVision(), info.HasAudio(), profile.HasVision, profile.HasAudioInput, info.ProjectionDim, info.ProjectorType)
			if len(info.Evidence) > 0 {
				fmt.Printf(" evidence=%s", strings.Join(info.Evidence, ","))
			}
			if info.Err != "" {
				fmt.Printf(" error=%q", info.Err)
			}
			fmt.Println(")")
		}
		a.currentModelProfile = profile
	}

	// Validate spawn settings before starting server
	if profile != nil && profile.ArchitectureSpec != nil {
		fmt.Println("[LoadModel] validating spawn settings")
		hw := resource.GetManager().HardwareSnapshot()
		freeVRAM := int64(float64(hw.FreeVRAMBytes) * 0.85 / 1024 / 1024)
		freeRAM := int64(float64(hw.FreeRAMBytes) * 0.90 / 1024 / 1024)
		ctxSize, _ := strconv.Atoi(intSetting(a.modelSettings, "context_size", 4096))
		gpuLayers, _ := strconv.Atoi(intSetting(a.modelSettings, "gpu_layers", 0))
		batchSize, _ := strconv.Atoi(intSetting(a.modelSettings, "batch_size", 512))

		// Sanity check: if we have a GPU but detected < 1GB free VRAM, detection is likely wrong.
		if hw.TotalVRAMBytes > 1024*1024*1024 && freeVRAM < 1024 {
			fmt.Printf("[LoadModel] WARNING: Detected free VRAM (%d MB) is suspiciously low for a GPU with %d MB total. "+
				"Detection may have failed. Proceeding with total*0.7 as fallback.\n",
				freeVRAM, hw.TotalVRAMBytes/1024/1024)
			freeVRAM = int64(float64(hw.TotalVRAMBytes) * 0.70 / 1024 / 1024)
		}

		fmt.Printf("[LoadModel] spawn params: ctx=%d gpuLayers=%d batch=%d freeVRAM=%d freeRAM=%d\n",
			ctxSize, gpuLayers, batchSize, freeVRAM, freeRAM)

		validation := gguf.ValidateSpawn(profile.ArchitectureSpec, ctxSize, gpuLayers, batchSize, freeVRAM, freeRAM)
		fmt.Printf("[LoadModel] validation result: level=%d msg=%q\n", validation.Level, validation.Message)

		forceLoad := false
		a.modelSettingsMu.RLock()
		if a.modelSettings != nil {
			if v, ok := a.modelSettings["force_load"].(bool); ok {
				forceLoad = v
			}
		}
		a.modelSettingsMu.RUnlock()

		switch validation.Level {
		case gguf.ValidationSafe:
			fmt.Println("[LoadModel] spawn validation passed")
		case gguf.ValidationWarning:
			fmt.Printf("[LoadModel] spawn validation warning: %s\n", validation.Message)
			// Emit event so UI can show the warning, but don't block
			a.emitEvent("model-load-warning", validation.Message+" "+validation.Suggestion)
		case gguf.ValidationCritical:
			if forceLoad {
				fmt.Printf("[LoadModel] FORCE LOAD: ignoring critical validation due to force_load=true: %s\n", validation.Message)
				a.emitEvent("model-load-warning",
					"Force load enabled. Model may OOM or crash. "+validation.Message)
			} else {
				a.logError(errors.LLMNoModelLoaded, errors.Critical, "llm.load.validate", "app-model.go:97",
					fmt.Errorf("%s %s", validation.Message, validation.Suggestion))
				fmt.Printf("[LoadModel] FAIL: spawn validation critical: %s\n", validation.Message)
				return fmt.Errorf("%s\n\n%s\n\nTo bypass this check, set force_load=true in model settings.",
					validation.Message, validation.Suggestion)
			}
		}
	} else {
		fmt.Println("[LoadModel] skipping spawn validation (no profile or arch spec)")
	}

	fitParams := true
	if profile != nil && profile.ArchitectureSpec != nil && profile.ArchitectureSpec.UseSlidingWindow {
		fitParams = false
		fmt.Println("[LoadModel] SWA model detected, disabling --fit to avoid scheduler crash")
	}
	spawnSettings := a.spawnSettingsFor(profile)
	fmt.Println("[LoadModel] spawning llama-server")
	serverOutput := loader.NewCapture(loader.DefaultCaptureBytes)
	cmd, err := spawnLlamaServerWithModel(modelPath, pairing.Path, spawnSettings, serverOutput, fitParams)
	if err != nil {
		rec := loader.FailureRecord{
			Stage:      loader.StageSpawn,
			ModelPath:  modelPath,
			ModelName:  filepath.Base(modelPath),
			Arch:       admission.Arch,
			Projector:  pairing.Path,
			Executable: findLlamaServer(),
			Args:       buildLlamaServerArgs(modelPath, pairing.Path, spawnSettings, fitParams),
			Detail:     err.Error(),
		}
		if path, writeErr := loader.WriteFailureRecord(a.loaderLogDir(), rec); writeErr == nil {
			rec.RecordPath = path
		} else {
			fmt.Printf("[LoadModel] WARN: could not persist diagnostics: %v\n", writeErr)
		}
		userMsg := rec.UserMessage()
		a.logLLMLoadError(fmt.Errorf("%s", userMsg), "app-model.go:84")
		fmt.Printf("[LoadModel] FAIL: spawn error: %v\n", err)
		a.loadingProgress.Store(0)
		a.loadingStage.Store("Failed")
		a.emitEvent("model:error", userMsg)
		return fmt.Errorf("failed to spawn llama-server: %w", err)
	}
	activeCmd := cmd
	fmt.Printf("[LoadModel] llama-server spawned (pid %d)\n", cmd.Process.Pid)
	a.serverProc = cmd
	pm := waitMonitor(cmd)

	a.loadingProgress.Store(0)
	a.loadingStage.Store("Starting")

	go func() {
		const maxTicks = 120
		defer func() {
			if r := recover(); r != nil {
				a.logError(errors.SYSPanicRecovered, errors.Critical, "llm.load.monitor", "app-model.go:goroutine", fmt.Errorf("monitor goroutine panic: %v", r))
				fmt.Printf("[LoadModel] MONITOR PANIC: %v\n", r)
				a.emitEvent("model:error", "internal error during model load monitoring")
			}
		}()
		fmt.Printf("[LoadModel] monitor started (pid %d, timeout %ds)\n", cmd.Process.Pid, maxTicks/2)
		fitRetries := 0
		const maxFitRetries = 1
		for tick := 0; tick < maxTicks; tick++ {
			time.Sleep(500 * time.Millisecond)
			pct := int(float64(tick) / float64(maxTicks) * 90.0)
			a.loadingProgress.Store(int64(pct) * 10)
			switch {
			case pct < 20:
				a.loadingStage.Store("Initializing")
			case pct < 50:
				a.loadingStage.Store("Loading weights")
			case pct < 80:
				a.loadingStage.Store("Allocating tensors")
			default:
				a.loadingStage.Store("Warming up")
			}
			a.emitEvent("model:progress", pct)

			if tick%20 == 0 && tick > 0 {
				fmt.Printf("[LoadModel] monitor tick %d/%d (%.0f%%), stage=%s, alive=%v\n",
					tick, maxTicks, float64(pct), a.loadingStage.Load(), pm.isAlive())
			}

			if !pm.isAlive() {
				captured := serverOutput.Tail(0)
				fmt.Printf("[LoadModel] monitor: process died at tick %d, captured output len=%d\n", tick, len(captured))
				if strings.Contains(captured, "GGML_SCHED_MAX_SPLIT_INPUTS") && fitRetries < maxFitRetries && tick < maxTicks-20 {
					fitRetries++
					fmt.Printf("[LoadModel] monitor: GGML_SCHED_MAX_SPLIT_INPUTS crash, retry %d/%d\n", fitRetries, maxFitRetries)
					a.logLLMCrash(fmt.Errorf("GGML_SCHED_MAX_SPLIT_INPUTS crash, retry %d/%d", fitRetries, maxFitRetries), "app-model.go:118")
					if a.serverProc != nil && a.serverProc.Process != nil {
						_ = a.serverProc.Process.Kill()
						_ = a.serverProc.Wait()
					}
					a.serverProc = nil
					serverOutput.Reset()
					if !waitForPortFree("8080", 5*time.Second) {
						a.logError(errors.LLMPortConflict, errors.Critical, "llm.load", "app-model.go:127", fmt.Errorf("port 8080 occupied after crash"))
						a.emitEvent("model:error", "port 8080 occupied after crash; cannot restart llama-server")
						return
					}
					cmd2, err2 := spawnLlamaServerWithModel(modelPath, pairing.Path, spawnSettings, serverOutput, false)
					if err2 != nil {
						rec := loader.FailureRecord{
							Stage:      loader.StageSpawn,
							ModelPath:  modelPath,
							ModelName:  filepath.Base(modelPath),
							Arch:       admission.Arch,
							Projector:  pairing.Path,
							Executable: findLlamaServer(),
							Args:       buildLlamaServerArgs(modelPath, pairing.Path, spawnSettings, false),
							Detail:     "retry after GGML_SCHED_MAX_SPLIT_INPUTS: " + err2.Error(),
						}
						if path, writeErr := loader.WriteFailureRecord(a.loaderLogDir(), rec); writeErr == nil {
							rec.RecordPath = path
						}
						userMsg := rec.UserMessage()
						a.logLLMLoadError(fmt.Errorf("%s", userMsg), "app-model.go:133")
						a.emitEvent("model:error", userMsg)
						return
					}
					activeCmd = cmd2
					a.serverProc = cmd2
					pm = waitMonitor(cmd2)
					tick = -1
					continue
				}

				rec := loader.FailureRecord{
					Stage:      loader.StageStartup,
					ModelPath:  modelPath,
					ModelName:  filepath.Base(modelPath),
					Arch:       admission.Arch,
					Projector:  pairing.Path,
					Executable: activeCmd.Path,
					Args:       activeCmd.Args,
					ExitStatus: loader.ExitStatus(pm.exitErr),
					Output:     serverOutput.Tail(loader.FailureTailBytes),
					Detail:     "llama-server exited before becoming healthy",
				}
				rec.OutputBytes, _, rec.OutputTruncated = serverOutput.Stats()
				if path, writeErr := loader.WriteFailureRecord(a.loaderLogDir(), rec); writeErr == nil {
					rec.RecordPath = path
				} else {
					fmt.Printf("[LoadModel] WARN: could not persist diagnostics: %v\n", writeErr)
				}
				errMsg := rec.UserMessage()
				fmt.Printf("[LoadModel] monitor: %s\n", rec.Format())
				a.logLLMCrash(fmt.Errorf("%s", errMsg), "app-model.go:146")
				a.loadingProgress.Store(0)
				a.loadingStage.Store("Failed")
				a.emitEvent("model:error", errMsg)
				return
			}

			if a.httpClient == nil {
				fmt.Println("[LoadModel] monitor: httpClient is nil, cannot health check")
				continue
			}
			if err := a.httpClient.Health(); err == nil {
				fmt.Println("[LoadModel] monitor: server healthy, model loaded")
				a.loadingProgress.Store(1000)
				a.loadingStage.Store("Ready")
				a.emitEvent("model:progress", 100)
				a.emitEvent("model:loaded", modelPath)
				return
			}
		}
		fmt.Println("[LoadModel] monitor: timeout after 60s")
		rec := loader.FailureRecord{
			Stage:      loader.StageTimeout,
			ModelPath:  modelPath,
			ModelName:  filepath.Base(modelPath),
			Arch:       admission.Arch,
			Projector:  pairing.Path,
			Executable: activeCmd.Path,
			Args:       activeCmd.Args,
			Output:     serverOutput.Tail(loader.FailureTailBytes),
			Detail:     "llama-server did not become ready after 60s",
		}
		rec.OutputBytes, _, rec.OutputTruncated = serverOutput.Stats()
		if path, writeErr := loader.WriteFailureRecord(a.loaderLogDir(), rec); writeErr == nil {
			rec.RecordPath = path
		} else {
			fmt.Printf("[LoadModel] WARN: could not persist diagnostics: %v\n", writeErr)
		}
		userMsg := rec.UserMessage()
		a.logError(errors.LLMServerTimeout, errors.High, "llm.load", "app-model.go:172", fmt.Errorf("%s", userMsg))
		a.loadingProgress.Store(0)
		a.loadingStage.Store("Timeout")
		a.emitEvent("model:error", userMsg)
	}()

	// Register with resource manager using actual memory estimates
	modelName := filepath.Base(modelPath)
	gpuLayers := 0
	if gl, ok := a.modelSettings["gpu_layers"].(int); ok {
		gpuLayers = gl
	}
	ctxSize := 4096
	if cs, ok := a.modelSettings["context_size"].(int); ok {
		ctxSize = cs
	}
	var gpuMB, ramMB int64 = int64(gpuLayers) * 100, 4000
	if profile != nil {
		est := gguf.EstimateMemory(profile, gpuLayers, ctxSize)
		gpuMB = est.GPUMB
		ramMB = est.RAMMB
		fmt.Printf("INFO: Registering model with estimated GPU:%dMB RAM:%dMB\n", gpuMB, ramMB)
	}
	resource.GetManager().RegisterModel(
		resource.ResourceLLM,
		modelName,
		modelPath,
		gpuMB,
		ramMB,
		resource.PriorityLLM,
		gpuLayers > 0,
	)

	return nil
}

// UnloadModel kills the llama-server process, freeing all model memory.
func (a *App) UnloadModel() error {
	if a.serverProc != nil && a.serverProc.Process != nil {
		_ = a.serverProc.Process.Kill()
		_ = a.serverProc.Wait()
		fmt.Println("INFO: llama-server killed, model memory freed")
	}
	a.serverProc = nil
	a.loadingProgress.Store(0)
	a.loadingStage.Store("Unloaded")
	a.emitEvent("model:unloaded")

	// Unregister from resource manager
	for _, m := range resource.GetManager().ModelsByType(resource.ResourceLLM) {
		resource.GetManager().UnregisterModel(resource.ResourceLLM, m.Name)
	}

	return nil
}

// IsModelLoaded returns true if a model is currently loaded.
func (a *App) IsModelLoaded() bool {
	if a.httpClient == nil {
		return false
	}
	return a.httpClient.IsLoaded()
}

// GetLoadingProgress returns the current model loading progress (0.0–1.0).
func (a *App) GetLoadingProgress() float64 {
	return float64(a.loadingProgress.Load()) / 1000.0
}

// GetLoadingStage returns the current loading stage description.
func (a *App) GetLoadingStage() string {
	if v := a.loadingStage.Load(); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetLoadedModelName returns the name of the currently loaded model.
func (a *App) GetLoadedModelName() string {
	if a.httpClient == nil {
		return "No model loaded"
	}
	info, err := a.httpClient.GetLoadedModel()
	if err != nil || info == nil {
		return "No model loaded"
	}
	if data, ok := info["data"].([]interface{}); ok && len(data) > 0 {
		if m, ok := data[0].(map[string]interface{}); ok {
			if id, ok := m["id"].(string); ok {
				return filepath.Base(id)
			}
		}
	}
	return "No model loaded"
}

// GetChatMessages returns all messages in the current chat context.
func (a *App) GetChatMessages() []map[string]interface{} {
	return a.chatMessages
}

// ClearChat clears the current chat context and triggered lore state.
func (a *App) ClearChat() error {
	a.chatMessages = []map[string]interface{}{}
	a.triggeredLore = map[string]bool{}
	return nil
}

// ---------------------------------------------------------------------------
// Sampler / Model Settings Wails Bindings
// ---------------------------------------------------------------------------

// GetSamplerSettings returns the current sampler settings.
func (a *App) GetSamplerSettings() map[string]interface{} {
	if a.samplerSettings == nil {
		return map[string]interface{}{}
	}
	return a.samplerSettings
}

// SetSamplerSettings updates and persists sampler settings.
func (a *App) SetSamplerSettings(settings map[string]interface{}) error {
	a.samplerSettings = settings
	if err := a.saveSamplerSettingsToFile(); err != nil {
		fmt.Printf("Warning: failed to save sampler settings to disk: %v\n", err)
	}
	return nil
}

// SetModelOptions merges options into modelSettings and persists them.
func (a *App) SetModelOptions(options map[string]interface{}) error {
	if a.modelSettings == nil {
		a.modelSettings = make(map[string]interface{})
	}
	for k, v := range options {
		a.modelSettings[k] = v
	}
	if err := a.saveModelSettingsToFile(); err != nil {
		fmt.Printf("Warning: failed to save model settings to disk: %v\n", err)
	}
	return nil
}

// GetModelOptions returns the current model options.
func (a *App) GetModelOptions() map[string]interface{} {
	if a.modelSettings == nil {
		return map[string]interface{}{}
	}
	return a.modelSettings
}

// ---------------------------------------------------------------------------
// Settings persistence
// ---------------------------------------------------------------------------

// getSettingsDir returns (and creates) the config directory next to the exe.
func (a *App) getSettingsDir() string {
	settingsDir := filepath.Join(paths.ProjectRoot(), "config")
	_ = os.MkdirAll(settingsDir, 0755)
	return settingsDir
}

func (a *App) saveSamplerSettingsToFile() error {
	if a.samplerSettings == nil {
		return nil
	}
	data, err := json.MarshalIndent(a.samplerSettings, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal sampler settings: %w", err)
	}
	path := filepath.Join(a.getSettingsDir(), "sampler-settings.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write sampler settings: %w", err)
	}
	return nil
}

func (a *App) loadSamplerSettingsFromFile() error {
	path := filepath.Join(a.getSettingsDir(), "sampler-settings.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read sampler settings: %w", err)
	}
	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return fmt.Errorf("unmarshal sampler settings: %w", err)
	}
	a.samplerSettings = settings
	return nil
}

func (a *App) saveModelSettingsToFile() error {
	if a.modelSettings == nil {
		return nil
	}
	data, err := json.MarshalIndent(a.modelSettings, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal model settings: %w", err)
	}
	path := filepath.Join(a.getSettingsDir(), "model-settings.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write model settings: %w", err)
	}
	return nil
}

func (a *App) loadModelSettingsFromFile() error {
	path := filepath.Join(a.getSettingsDir(), "model-settings.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read model settings: %w", err)
	}
	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return fmt.Errorf("unmarshal model settings: %w", err)
	}
	a.modelSettings = settings

	return nil
}

// ---------------------------------------------------------------------------
// Memory estimation + safe defaults
// ---------------------------------------------------------------------------

// GetModelMemoryEstimate returns estimated GPU/RAM for a model with given settings.
func (a *App) GetModelMemoryEstimate(path string, gpuLayers, contextSize int) (*types.MemoryEstimate, error) {
	profile, err := gguf.ExtractModelProfile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read model metadata: %w", err)
	}
	est := gguf.EstimateMemory(profile, gpuLayers, contextSize)
	return &est, nil
}

// applySafeDefaults adjusts gpu_layers and context_size to fit available memory.
// It uses the architecture-aware memory model for accurate estimation.
func (a *App) applySafeDefaults(profile *types.ModelProfile) {
	fmt.Println("[applySafeDefaults] enter")
	if a.modelSettings == nil {
		a.modelSettings = make(map[string]interface{})
	}

	fmt.Println("[applySafeDefaults] getting hardware snapshot")
	hw := resource.GetManager().HardwareSnapshot()
	freeVRAM := int64(float64(hw.FreeVRAMBytes) * 0.85 / 1024 / 1024) // 85% of free VRAM
	freeRAM := int64(float64(hw.FreeRAMBytes) * 0.90 / 1024 / 1024)   // 90% of free RAM
	fmt.Printf("[applySafeDefaults] freeVRAM=%d freeRAM=%d\n", freeVRAM, freeRAM)

	fmt.Println("[applySafeDefaults] reading settings")
	ctxSize, _ := strconv.Atoi(intSetting(a.modelSettings, "context_size", 4096))
	if ctxSize > profile.UsableContext && profile.UsableContext > 0 {
		ctxSize = profile.UsableContext
		a.modelSettings["context_size"] = ctxSize
		fmt.Printf("[applySafeDefaults] capped context to %d\n", ctxSize)
	}

	gpuLayers, _ := strconv.Atoi(intSetting(a.modelSettings, "gpu_layers", 0))
	batchSize, _ := strconv.Atoi(intSetting(a.modelSettings, "batch_size", 512))
	fmt.Printf("[applySafeDefaults] settings: ctx=%d gpuLayers=%d batch=%d\n", ctxSize, gpuLayers, batchSize)

	// Build architecture spec. If profile has the new ArchitectureSpec, use it.
	// Otherwise, fall back to the legacy profile fields.
	spec := profile.ArchitectureSpec
	if spec == nil {
		fmt.Println("[applySafeDefaults] archSpec is nil, building fallback")
		spec = &types.ArchitectureSpec{
			Family:         profile.Family,
			Layers:         profile.LayerCount,
			ContextLength:  profile.MaxContext,
			FileSizeBytes:  profile.FileSizeBytes,
			ParameterCount: profile.ParameterCount,
			QuantType:      profile.FileType,
			BitsPerWeight:  profile.BitsPerWeight,
		}
		// Estimate missing critical fields from parameter count
		if spec.Layers <= 0 && profile.ParameterCount > 0 {
			spec.Layers = int(profile.ParameterCount / 100_000_000)
			if spec.Layers < 1 {
				spec.Layers = 1
			}
			if spec.Layers > 80 {
				spec.Layers = 80
			}
		}
		if spec.DModel <= 0 {
			spec.DModel = gguf.EstimateDModel(float64(profile.ParameterCount))
		}
		if spec.NumHeads <= 0 {
			spec.NumHeads = gguf.EstimateNumHeads(spec.DModel)
		}
		if spec.NumKVHeads <= 0 {
			spec.NumKVHeads = spec.NumHeads
		}
		if spec.FFNSize <= 0 {
			spec.FFNSize = gguf.EstimateFFNSize(spec.DModel)
		}
		if spec.HeadDim <= 0 && spec.DModel > 0 && spec.NumHeads > 0 {
			spec.HeadDim = spec.DModel / spec.NumHeads
		}
	} else {
		fmt.Printf("[applySafeDefaults] archSpec: layers=%d dmodel=%d heads=%d kvheads=%d\n",
			spec.Layers, spec.DModel, spec.NumHeads, spec.NumKVHeads)
	}

	fmt.Println("[applySafeDefaults] creating memory model")
	model := gguf.NewMemoryModel(spec)
	fmt.Println("[applySafeDefaults] calling ClampToHardware")
	ctxSize, gpuLayers, batchSize, finalEst, adjustments := model.ClampToHardware(ctxSize, gpuLayers, batchSize, freeVRAM, freeRAM)
	fmt.Println("[applySafeDefaults] ClampToHardware returned")

	for _, adj := range adjustments {
		fmt.Printf("[applySafeDefaults] adjustment: %s\n", adj)
	}
	fmt.Printf("[applySafeDefaults] final: GPU=%dMB RAM=%dMB ctx=%d gpuLayers=%d batch=%d\n",
		finalEst.GPUMB, finalEst.RAMMB, ctxSize, gpuLayers, batchSize)

	// Apply clamped settings
	a.modelSettings["context_size"] = ctxSize
	a.modelSettings["gpu_layers"] = gpuLayers
	a.modelSettings["batch_size"] = batchSize
	fmt.Println("[applySafeDefaults] done")
}

// ---------------------------------------------------------------------------
// Model directory + file pickers
// ---------------------------------------------------------------------------

// loaderLogDir returns where model-load failure records are written.
func (a *App) loaderLogDir() string {
	if a.loaderDiagnosticsDir != "" {
		return a.loaderDiagnosticsDir
	}
	return loader.LogDir(paths.ProjectRoot())
}

// GetModelsDir returns the absolute path to the models directory.
func (a *App) GetModelsDir() string {
	modelsDir := filepath.Join(paths.ProjectRoot(), "models")
	_ = os.MkdirAll(modelsDir, 0755)
	return modelsDir
}

// OpenModelFilePicker opens a native file dialog and returns the selected model path.
func (a *App) OpenModelFilePicker() (string, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "Select Model File",
		DefaultDirectory: a.GetModelsDir(),
		Filters: []runtime.FileFilter{
			{DisplayName: "GGUF Models (*.gguf)", Pattern: "*.gguf"},
			{DisplayName: "All Model Files (*.gguf;*.bin)", Pattern: "*.gguf;*.bin"},
		},
	})
	if err != nil {
		return "", fmt.Errorf("file dialog error: %w", err)
	}
	return path, nil
}

// LoadModelMetadata reads GGUF metadata and returns a ModelProfile.
func (a *App) LoadModelMetadata(path string) (*types.ModelProfile, error) {
	return gguf.NewReader().ReadProfile(path)
}

// UpdateRuntimeSettings stores runtime settings for the next model load.
func (a *App) UpdateRuntimeSettings(settings *types.RuntimeSettings) error {
	a.modelSettings = map[string]interface{}{
		"threads":         settings.Threads,
		"batch_size":      settings.BatchSize,
		"context_size":    settings.ContextSize,
		"gpu_layers":      settings.GPULayers,
		"rope_mode":       settings.RopeMode,
		"rope_factor":     settings.RopeFactor,
		"rope_base":       settings.RopeBase,
		"flash_attention": settings.FlashAttention,
		"tensor_split":    settings.TensorSplit,
		"main_gpu":        settings.MainGPU,
		"offload_kqv":     settings.OffloadKQV,
		"use_mmap":        settings.UseMMap,
		"use_mlock":       settings.UseMLock,
		"vocab_only":      settings.VocabOnly,
	}
	return nil
}

// UpdatePromptSettings stores prompt pipeline settings.
func (a *App) UpdatePromptSettings(settings *types.PromptSettings) error {
	if a.modelSettings == nil {
		a.modelSettings = make(map[string]interface{})
	}
	a.modelSettings["prompt_template"] = settings.PromptTemplate
	a.modelSettings["custom_jinja_template"] = settings.CustomJinjaTemplate
	a.modelSettings["system_prompt_override"] = settings.SystemPromptOverride
	a.modelSettings["user_prefix"] = settings.UserPrefix
	a.modelSettings["assistant_prefix"] = settings.AssistantPrefix
	a.modelSettings["stop_sequences"] = settings.StopSequences
	return nil
}
