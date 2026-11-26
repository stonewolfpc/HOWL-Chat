package main

import (
	"fmt"
	"howl-chat/internal/paths"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"howl-chat/internal/backend/chatpolicy"
	"howl-chat/internal/backend/gguf"
	llamabackend "howl-chat/internal/backend/llama"
	"howl-chat/internal/backend/loader"
	"howl-chat/internal/backend/types"
	"howl-chat/internal/errors"
	"howl-chat/internal/proc"
	"howl-chat/internal/process"
	"howl-chat/internal/resource"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ---------------------------------------------------------------------------
// Model Loading Wails Bindings
// ---------------------------------------------------------------------------

// LoadModel starts llama-server with the given model path.
func (a *App) LoadModel(modelPath string) (err error) {
	// One lifecycle transition at a time: the decision and the stop/spawn that
	// follows it must see the same ledger, otherwise two loads can both pass on
	// one free value. The transition is handed to the monitor so it stays held
	// through the health check and registration, not only through the spawn.
	transition, gateErr := a.beginLifecycleTransition(resource.ResourceLLM, "loading", lifecycleStartWait)
	if gateErr != nil {
		return gateErr
	}
	monitorOwnsTransition := false
	defer func() {
		if !monitorOwnsTransition {
			transition.Release()
		}
	}()
	if err := a.ensureAcceptingWork("load model"); err != nil {
		return err
	}
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
		a.engine.SetProgressMilli(0)
		a.engine.SetStage("Rejected")
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
		a.engine.SetProgressMilli(0)
		a.engine.SetStage("Rejected")
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
		a.engine.SetProgressMilli(0)
		a.engine.SetStage("Rejected")
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
	if existing := a.level0().Processes.Current(); existing != nil && existing.Process != nil {
		fmt.Println("[LoadModel] existing serverProc found")
		if a.engine.Client() != nil {
			if loaded, _ := a.engine.Client().GetLoadedModel(); loaded != nil {
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
		progress := a.engine.ProgressMilli()
		if progress > 0 && progress < 1000 {
			fmt.Printf("[LoadModel] load in progress (progress=%d), skipping\n", progress)
			return nil
		}
		fmt.Println("[LoadModel] stopping existing server")
		if err := a.level0().Processes.Stop(5 * time.Second); err != nil {
			fmt.Printf("[LoadModel] WARN: existing llama-server did not stop cleanly: %v\n", err)
			return fmt.Errorf("existing llama-server did not stop: %w", err)
		}
		a.unregisterChatModels()
		// The old placement belonged to the process just stopped. Keeping it
		// would let the next load's notification quote a footprint that no
		// longer exists.
		a.engine.SetChatPlacement(nil)
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
		// The template is NOT written into the saved settings here. Doing that
		// made the first family's template stick forever: the second family
		// loaded with if _, exists… still true and was prompted in the wrong
		// format. It is injected per load in spawnSettingsFor().
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
		a.engine.SetProfile(profile)
	}

	// Howl supplies an explicit context, batch and GPU-layer plan, so llama.cpp's
	// fit pass has nothing it is allowed to adjust. Bypass it explicitly and say
	// why; leaving it on looked like the engine was sizing the model when it was
	// not. SWA models additionally need it off to avoid a scheduler crash.
	fitParams := false
	fmt.Println("[LoadModel] --fit off: Howl admission supplies the explicit context, batch and GPU-layer plan")
	if profile != nil && profile.ArchitectureSpec != nil && profile.ArchitectureSpec.UseSlidingWindow {
		fmt.Println("[LoadModel] SWA model detected; --fit off is also required to avoid the scheduler crash")
	}
	if err := a.ensureAcceptingWork("load model"); err != nil {
		return err
	}

	// Probe the shipped binary once per load. Optional flags are only sent when
	// this binary advertises them; an older build gets an honest, smaller set
	// instead of an argument it rejects at startup.
	caps := llamabackend.ProbeCLI(findLlamaServerForRuntime())
	fmt.Printf("[LoadModel] llama capability probe: probed=%t version=%q build=%q commit=%q flags=%d err=%q\n",
		caps.Probed, caps.Version, caps.Build, caps.Commit, len(caps.Flags), caps.Error)

	modelName := filepath.Base(modelPath)
	if a.admission == nil {
		a.admission = resource.DefaultAuthority()
	}
	loadRequest := a.chatAdmissionRequest(profile, modelPath, true)
	// The user asked to load this model. If the automatic plan cannot fit or the
	// hardware cannot be measured, admission returns a warned reduced plan
	// instead of locking the load out.
	loadRequest.AllowManual = true
	// One primary GPU workload at a time: a chat start that wants the GPU hands
	// the device over instead of being demoted to the CPU behind the user's
	// back. The handoff stops the image engine, waits for the driver to report
	// the freed memory, and only then lets admission run against fresh numbers.
	if loadRequest.RequireGPU {
		if name, _ := a.loadedImageModelInfo(); name != "" {
			fmt.Printf("[VRAM] automatic handoff: unloading image model %q to load the chat model\n", name)
			fit, err := a.FreeVRAMFor("llm", modelPath, true)
			if err != nil {
				a.engine.SetProgressMilli(0)
				a.engine.SetStage("Refused")
				a.emitEvent("model:error", err.Error())
				return err
			}
			fmt.Printf("[VRAM] handoff complete: %s\n", fit.Message)
		}
	}
	loadDecision, loadLease := a.admission.Admit(loadRequest)
	fmt.Printf("[LoadModel] admission: admitted=%t manual=%t %s\n", loadDecision.Admitted, loadDecision.Manual, loadDecision.Reason)
	// Live device usage at decision time. The difference after the health check
	// is this start's measured device cost, which is compared with the engine's
	// own claim before anything is written to the ledger.
	residencyBeforeMB := loadDecision.Estimate.TotalVRAMMB - loadDecision.Estimate.FreeVRAMMB
	if !loadDecision.Admitted {
		a.logError(errors.LLMNoModelLoaded, errors.Critical, "llm.load.admission", "app-model.go:LoadModel", fmt.Errorf("%s", loadDecision.Reason))
		a.engine.SetProgressMilli(0)
		a.engine.SetStage("Refused")
		a.emitEvent("model:error", loadDecision.Reason)
		return fmt.Errorf("%s", loadDecision.Reason)
	}

	spawnSettings := a.spawnSettingsFor(profile)
	if kvType, warning := resolveKVType(spawnSettings); warning != "" {
		fmt.Printf("[LoadModel] WARNING: %s\n", warning)
		a.emitEvent("model:warning", warning)
		spawnSettings["kv_cache_type"] = kvType
	}
	if adjustment := loadDecision.Adjustment; adjustment != nil {
		spawnSettings["context_size"] = adjustment.ContextSize
		spawnSettings["gpu_layers"] = adjustment.GPULayers
		spawnSettings["batch_size"] = adjustment.BatchSize
		if adjustment.Auto {
			fmt.Printf("[LoadModel] automatic placement: context=%d gpu_layers=%d batch=%d\n",
				adjustment.ContextSize, adjustment.GPULayers, adjustment.BatchSize)
		} else {
			fmt.Printf("[LoadModel] fitted plan applied: context=%d gpu_layers=%d batch=%d\n",
				adjustment.ContextSize, adjustment.GPULayers, adjustment.BatchSize)
		}
		if loadDecision.Manual {
			a.engine.SetStage("Starting (manual plan)")
		} else {
			a.engine.SetStage("Starting (fitted plan)")
		}
	}
	// A negative layer count must never reach the engine: llama.cpp reads -1 as
	// "all layers", which is not what an unplanned auto request means. Every
	// admitted auto request carries an adjustment; this is the guard that makes
	// that invariant observable instead of assumed.
	if gl, _ := strconv.Atoi(intSetting(spawnSettings, "gpu_layers", 0)); gl < 0 {
		spawnSettings["gpu_layers"] = 0
		warning := "automatic placement was requested but admission returned no layer plan; starting on the CPU rather than sending -1 (all layers) to the engine."
		fmt.Printf("[LoadModel] WARNING: %s\n", warning)
		a.emitEvent("model:warning", warning)
	}
	if loadDecision.Warning != "" {
		fmt.Printf("[LoadModel] WARNING: %s\n", loadDecision.Warning)
		a.emitEvent("model:warning", loadDecision.Warning)
	}
	fmt.Println("[LoadModel] spawning llama-server")
	serverOutput := loader.NewCapture(loader.DefaultCaptureBytes)
	cmd, err := spawnLlamaServerWithModel(modelPath, pairing.Path, spawnSettings, serverOutput, fitParams, caps)
	if err != nil {
		loadLease.Release()
		rec := loader.FailureRecord{
			Stage:      loader.StageSpawn,
			ModelPath:  modelPath,
			ModelName:  filepath.Base(modelPath),
			Arch:       admission.Arch,
			Projector:  pairing.Path,
			Executable: findLlamaServerForRuntime(),
			Args:       buildLlamaServerArgs(modelPath, pairing.Path, spawnSettings, fitParams, caps),
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
		a.engine.SetProgressMilli(0)
		a.engine.SetStage("Failed")
		a.emitEvent("model:error", userMsg)
		return fmt.Errorf("failed to spawn llama-server: %w", err)
	}
	activeCmd := cmd
	fmt.Printf("[LoadModel] llama-server spawned (pid %d)\n", cmd.Process.Pid)
	pm := waitMonitor(cmd)
	activeHandle, adoptErr := a.level0().Processes.Adopt(cmd, pm.wait)
	if adoptErr != nil {
		proc.KillTree(cmd.Process.Pid)
		_ = cmd.Process.Kill()
		loadLease.Release()
		a.engine.SetProgressMilli(0)
		a.engine.SetStage("Failed")
		a.emitEvent("model:error", "could not establish ownership of llama-server: "+adoptErr.Error())
		return fmt.Errorf("failed to adopt llama-server: %w", adoptErr)
	}
	// Shutdown may have passed its stop step while the model was being
	// validated or spawned. Adopt first, then stop the child if that happened:
	// no child may outlive the shutdown that started during this load.
	if a.level0().Gate.Stopping() {
		loadLease.Release()
		_ = a.level0().Processes.Stop(llamaShutdownTimeout)
		return fmt.Errorf("model load cancelled because Howl is shutting down")
	}

	a.engine.SetProgressMilli(0)
	a.engine.SetStage("Starting")

	monitorOwnsTransition = true
	go func() {
		// The transition covers the whole start, including the health check and
		// registration: a stop that landed between spawn and health used to be
		// able to kill the server the user was waiting on.
		defer transition.Release()
		const maxTicks = 120
		defer func() {
			if r := recover(); r != nil {
				a.failChatLoad(loadLease)
				a.logError(errors.SYSPanicRecovered, errors.Critical, "llm.load.monitor", "app-model.go:goroutine", fmt.Errorf("monitor goroutine panic: %v", r))
				fmt.Printf("[LoadModel] MONITOR PANIC: %v\n", r)
				a.emitEvent("model:error", "internal error during model load monitoring")
			}
		}()
		fmt.Printf("[LoadModel] monitor started (pid %d, timeout %ds)\n", cmd.Process.Pid, maxTicks/2)
		fitRetries := 0
		const maxFitRetries = 1
		kvFallbackRetries := 0

		// One retry path for both recoverable startup failures. The caller has
		// already decided what to change; this only stops the dead child, frees
		// the port, and adopts the replacement with the same load identity.
		retryAfterFailure := func(detail string) bool {
			if err := a.level0().Processes.Stop(5 * time.Second); err != nil {
				fmt.Printf("[LoadModel] WARN: crashed llama-server did not stop cleanly: %v\n", err)
				a.failChatLoad(loadLease)
				a.logError(errors.LLMServerTimeout, errors.High, "llm.load.stop", "app-model.go:118", err)
				a.emitEvent("model:error", "crashed llama-server did not stop; cannot retry the load")
				return false
			}
			serverOutput.Reset()
			if a.level0().Gate.Stopping() {
				fmt.Println("[LoadModel] retry cancelled because Howl is shutting down")
				a.failChatLoad(loadLease)
				return false
			}
			if !waitForPortFree("8080", 5*time.Second) {
				a.failChatLoad(loadLease)
				a.logError(errors.LLMPortConflict, errors.Critical, "llm.load", "app-model.go:127", fmt.Errorf("port 8080 occupied after crash"))
				a.emitEvent("model:error", "port 8080 occupied after crash; cannot restart llama-server")
				return false
			}
			cmd2, err2 := spawnLlamaServerWithModel(modelPath, pairing.Path, spawnSettings, serverOutput, false, caps)
			if err2 != nil {
				rec := loader.FailureRecord{
					Stage:      loader.StageSpawn,
					ModelPath:  modelPath,
					ModelName:  filepath.Base(modelPath),
					Arch:       admission.Arch,
					Projector:  pairing.Path,
					Executable: findLlamaServerForRuntime(),
					Args:       buildLlamaServerArgs(modelPath, pairing.Path, spawnSettings, false, caps),
					Detail:     detail + ": " + err2.Error(),
				}
				if path, writeErr := loader.WriteFailureRecord(a.loaderLogDir(), rec); writeErr == nil {
					rec.RecordPath = path
				}
				userMsg := rec.UserMessage()
				a.failChatLoad(loadLease)
				a.logLLMLoadError(fmt.Errorf("%s", userMsg), "app-model.go:133")
				a.emitEvent("model:error", userMsg)
				return false
			}
			activeCmd = cmd2
			pm = waitMonitor(cmd2)
			nextHandle, adoptErr := a.level0().Processes.Adopt(cmd2, pm.wait)
			if adoptErr != nil {
				proc.KillTree(cmd2.Process.Pid)
				_ = cmd2.Process.Kill()
				a.failChatLoad(loadLease)
				a.emitEvent("model:error", "could not establish ownership of the replacement llama-server: "+adoptErr.Error())
				return false
			}
			activeHandle = nextHandle
			if a.level0().Gate.Stopping() {
				a.failChatLoad(loadLease)
				return false
			}
			return true
		}

		for tick := 0; tick < maxTicks; tick++ {
			time.Sleep(500 * time.Millisecond)
			pct := int(float64(tick) / float64(maxTicks) * 90.0)
			a.engine.SetProgressMilli(int64(pct) * 10)
			switch {
			case pct < 20:
				a.engine.SetStage("Initializing")
			case pct < 50:
				a.engine.SetStage("Loading weights")
			case pct < 80:
				a.engine.SetStage("Allocating tensors")
			default:
				a.engine.SetStage("Warming up")
			}
			a.emitEvent("model:progress", pct)

			if tick%20 == 0 && tick > 0 {
				fmt.Printf("[LoadModel] monitor tick %d/%d (%.0f%%), stage=%s, alive=%v\n",
					tick, maxTicks, float64(pct), a.engine.Stage(), pm.isAlive())
			}

			if !pm.isAlive() {
				captured := serverOutput.Tail(0)
				fmt.Printf("[LoadModel] monitor: process died at tick %d, captured output len=%d\n", tick, len(captured))
				if tick < maxTicks-20 {
					if requested, warning, ok := prepareQuantizedKVRetry(captured, spawnSettings, &kvFallbackRetries); ok {
						fmt.Printf("[LoadModel] monitor: quantized KV rejected (requested %s), retrying once with f16\n", requested)
						a.logLLMCrash(fmt.Errorf("quantized KV rejection: requested %s, retrying with f16", requested), "app-model.go:118")
						a.emitEvent("model:warning", warning)
						if retryAfterFailure("retry after quantized KV rejection") {
							tick = -1
							continue
						}
						return
					}
					if strings.Contains(captured, "GGML_SCHED_MAX_SPLIT_INPUTS") && fitRetries < maxFitRetries {
						fitRetries++
						fmt.Printf("[LoadModel] monitor: GGML_SCHED_MAX_SPLIT_INPUTS crash, retry %d/%d\n", fitRetries, maxFitRetries)
						a.logLLMCrash(fmt.Errorf("GGML_SCHED_MAX_SPLIT_INPUTS crash, retry %d/%d", fitRetries, maxFitRetries), "app-model.go:118")
						if retryAfterFailure("retry after GGML_SCHED_MAX_SPLIT_INPUTS") {
							tick = -1
							continue
						}
						return
					}
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
				a.failChatLoad(loadLease)
				a.engine.SetProgressMilli(0)
				a.engine.SetStage("Failed")
				a.emitEvent("model:error", errMsg)
				return
			}

			if a.engine.Client() == nil {
				fmt.Println("[LoadModel] monitor: no inference client, cannot health check")
				continue
			}
			if err := a.engine.Client().Health(); err == nil {
				fmt.Println("[LoadModel] monitor: server healthy, model loaded")
				// The health check is evidence for the process we adopted, not for
				// whichever llama-server happens to own the port now. A stop or
				// replacement during the check invalidates this load's registration.
				if current := a.level0().Processes.Current(); current != activeCmd || !a.level0().Processes.Owns(activeHandle) {
					fmt.Println("[LoadModel] monitor: process identity changed during health check; not registering")
					a.failChatLoad(loadLease)
					return
				}
				// The plan was Howl's arithmetic. Parse the engine's own load
				// statements now that the model is healthy, and report that
				// separately instead of presenting the estimate as the result.
				placement := llamabackend.ParsePlacement(serverOutput.Tail(0))
				placementReport := chatPlacementReport(modelPath, spawnSettings, loadDecision, placement)
				a.engine.SetChatPlacement(placementReport)
				// Reconcile the ledger with observation: the larger of the engine's
				// reported total and the device's own change. The plan stays the
				// planning bound; it is not the resident truth.
				ledgerVRAMMB := loadDecision.Estimate.RequiredVRAMMB
				if snap, snapErr := resource.GetManager().RefreshSnapshot(); snapErr == nil {
					const mb = int64(1024 * 1024)
					afterMB := (snap.TotalVRAMBytes - snap.FreeVRAMBytes) / mb
					observation := resource.NewResidencyObservation(
						resource.ResourceLLM, modelName,
						loadDecision.Estimate.RequiredVRAMMB,
						int64(placement.TotalGPUMB()),
						afterMB-residencyBeforeMB,
					)
					ledgerVRAMMB = observation.LedgerVRAMMB()
					a.recordResidency(observation)
				}
				fmt.Printf("[LoadModel] actual placement: known=%t layers=%d/%d vram=%.0fMB ram=%.0fMB context=%d matches_plan=%t\n",
					placement.Known, placement.GPULayers, placement.TotalLayers, placement.GPUMB, placement.RAMMB, placement.Context, placementReport.MatchesPlan)
				// The health check has consumed all the placement evidence. Drop the
				// load-phase capture now so no later inference output can accumulate
				// in memory through a writer that failed to observe its marker.
				serverOutput.Reset()
				// Register only after llama-server has answered a health check;
				// a spawned child is not a loaded model.
				if !a.registerChatModelAfterHealth(
					modelName, modelPath,
					ledgerVRAMMB,
					loadDecision.Estimate.RequiredRAMMB,
					activeCmd, activeHandle,
				) {
					fmt.Println("[LoadModel] monitor: process identity changed before registration; not registering")
					a.failChatLoad(loadLease)
					return
				}
				loadLease.Release()
				a.engine.SetProgressMilli(1000)
				a.engine.SetStage("Ready")
				a.emitEvent("model:progress", 100)
				// A success notification is emitted before the loaded event so
				// the UI can show the observed arrangement, not merely "Ready".
				a.emitEvent("model:placement", placementReport.Message)
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
		// The timeout branch owns the child now: unregister any optimistic
		// entry, stop the process, and return the reservation.
		a.failChatLoad(loadLease)
		a.logError(errors.LLMServerTimeout, errors.High, "llm.load", "app-model.go:172", fmt.Errorf("%s", userMsg))
		a.engine.SetProgressMilli(0)
		a.engine.SetStage("Timeout")
		a.emitEvent("model:error", userMsg)
	}()

	return nil
}

// failChatLoad is the single failure cleanup for a chat start that has not
// reached a healthy server. It is safe to call more than once.
func (a *App) failChatLoad(lease *resource.AdmissionLease) {
	a.engine.SetChatPlacement(nil)
	if err := a.level0().Processes.Stop(llamaShutdownTimeout); err != nil {
		// A child that will not stop must not also poison admission: the
		// reservation is for this attempt, not for a process that failed to
		// start. Release and unregister below even when Stop errors, otherwise
		// the next load is refused against memory nobody is using.
		fmt.Printf("[LoadModel] WARN: failed child did not stop cleanly: %v\n", err)
	}
	lease.Release()
	a.unregisterChatModels()
}

func (a *App) unregisterChatModels() {
	for _, m := range resource.GetManager().ModelsByType(resource.ResourceLLM) {
		resource.GetManager().UnregisterModel(resource.ResourceLLM, m.Name)
	}
}

// registerChatModelAfterHealth is the one route from a successful health check
// to a chat registration. It refuses when the process that answered is no
// longer the process this load adopted, so a stop or replacement during the
// check cannot leave a registration for a dead child.
func (a *App) registerChatModelAfterHealth(
	modelName, modelPath string,
	vramMB, ramMB int64,
	activeCmd *exec.Cmd,
	activeHandle *process.Handle,
) bool {
	if current := a.level0().Processes.Current(); current != activeCmd || !a.level0().Processes.Owns(activeHandle) {
		return false
	}
	a.unregisterChatModels()
	resource.GetManager().RegisterModel(
		resource.ResourceLLM,
		modelName,
		modelPath,
		vramMB,
		ramMB,
		resource.PriorityLLM,
		vramMB > 0,
	)
	return true
}

// chatAdmissionRequest describes the chat workload from the settings actually
// passed to llama-server. Replacing is true only for the preflight that may stop
// an existing chat model; the final decision after the stop must not claim that
// credit.
func (a *App) chatAdmissionRequest(profile *types.ModelProfile, modelPath string, replacing bool) resource.AdmissionRequest {
	live := a.settings.Model()
	ctxSize, _ := strconv.Atoi(intSetting(live, "context_size", 4096))
	// -1 is automatic placement (the planner picks the largest plan that fits);
	// 0 remains an explicit CPU-only choice and saved zeroes are not migrated.
	gpuLayers, _ := strconv.Atoi(intSetting(live, "gpu_layers", -1))
	batchSize, _ := strconv.Atoi(intSetting(live, "batch_size", 512))
	deviceIndex := 0
	switch v := live["main_gpu"].(type) {
	case int:
		if v >= 0 {
			deviceIndex = v
		}
	case float64:
		if v >= 0 {
			deviceIndex = int(v)
		}
	case string:
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			deviceIndex = n
		}
	}
	req := resource.AdmissionRequest{
		Kind:        resource.ResourceLLM,
		Name:        filepath.Base(modelPath),
		Path:        modelPath,
		ContextSize: ctxSize,
		GPULayers:   gpuLayers,
		BatchSize:   batchSize,
		// Auto (-1) asks for GPU placement; 0 is the explicit CPU-only choice.
		RequireGPU: gpuLayers != 0,
		// The KV cache is costed at the type the engine will actually build.
		KVCacheBytesPerElement: resource.KVCacheBytesPerElement(kvCacheType(live)),
		DeviceIndex:            deviceIndex,
	}
	if profile != nil {
		req.Arch = profile.ArchitectureSpec
	}
	if replacing {
		req.ReplacingKind = resource.ResourceLLM
	}
	return req
}

// UnloadModel kills the llama-server process, freeing all model memory.
func (a *App) UnloadModel() error {
	// A stop is a transition too. If a chat start owns the slot, this refuses
	// with the conflicting transition named instead of killing its child.
	transition, err := a.beginLifecycleTransition(resource.ResourceLLM, "stop", lifecycleStopWait)
	if err != nil {
		return err
	}
	defer transition.Release()
	if err := a.level0().Processes.Stop(llamaShutdownTimeout); err != nil {
		return fmt.Errorf("stop llama-server: %w", err)
	}
	fmt.Println("INFO: llama-server stopped, model memory freed")
	a.engine.SetChatPlacement(nil)
	a.engine.SetProgressMilli(0)
	a.engine.SetStage("Unloaded")
	a.emitEvent("model:unloaded")

	// Unregister from resource manager
	for _, m := range resource.GetManager().ModelsByType(resource.ResourceLLM) {
		resource.GetManager().UnregisterModel(resource.ResourceLLM, m.Name)
	}

	return nil
}

// IsModelLoaded returns true if a model is currently loaded.
func (a *App) IsModelLoaded() bool {
	if a.engine.Client() == nil {
		return false
	}
	return a.engine.Client().IsLoaded()
}

// GetLoadingProgress returns the current model loading progress (0.0–1.0).
func (a *App) GetLoadingProgress() float64 {
	return a.engine.Progress()
}

// GetLoadingStage returns the current loading stage description.
func (a *App) GetLoadingStage() string {
	return a.engine.Stage()
}

// GetLoadedModelName returns the name of the currently loaded model.
func (a *App) GetLoadedModelName() string {
	if a.engine.Client() == nil {
		return "No model loaded"
	}
	info, err := a.engine.Client().GetLoadedModel()
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
	return a.sessionMessages()
}

// ClearChat clears the current chat context and triggered lore state.
func (a *App) ClearChat() error {
	a.session().ClearEntries()
	return nil
}

// ---------------------------------------------------------------------------
// Sampler / Model Settings Wails Bindings
// ---------------------------------------------------------------------------

// GetSamplerSettings returns the current sampler settings.
func (a *App) GetSamplerSettings() map[string]interface{} {
	return a.settings.Sampler()
}

// SetSamplerSettings updates and persists sampler settings.
func (a *App) SetSamplerSettings(settings map[string]interface{}) error {
	a.settings.SetSampler(settings)
	if err := a.settings.SaveSamplerFile(); err != nil {
		fmt.Printf("Warning: failed to save sampler settings to disk: %v\n", err)
	}
	return nil
}

// SetModelOptions merges options into the model section and persists them.
func (a *App) SetModelOptions(options map[string]interface{}) error {
	a.settings.MergeModel(options)
	if err := a.settings.SaveModelFile(); err != nil {
		fmt.Printf("Warning: failed to save model settings to disk: %v\n", err)
	}
	return nil
}

// GetModelOptions returns the current model options.
func (a *App) GetModelOptions() map[string]interface{} {
	return a.settings.Model()
}

// ---------------------------------------------------------------------------
// Settings persistence
// ---------------------------------------------------------------------------

// getSettingsDir returns (and creates) the config directory next to the exe.
func (a *App) getSettingsDir() string {
	settingsDir := paths.InternalConfigDir()
	_ = os.MkdirAll(settingsDir, 0755)
	return settingsDir
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

// applySafeDefaults only normalizes model metadata and caps the context to the
// model's usable window. It deliberately does not estimate memory or clamp GPU
// layers: the admission authority owns that decision and refuses rather than
// silently changing the user's requested workload.
func (a *App) applySafeDefaults(profile *types.ModelProfile) {
	if profile == nil {
		return
	}
	live := a.settings.Model()
	ctxSize, _ := strconv.Atoi(intSetting(live, "context_size", 4096))
	if ctxSize > profile.UsableContext && profile.UsableContext > 0 {
		ctxSize = profile.UsableContext
	}

	if profile.ArchitectureSpec == nil {
		spec := &types.ArchitectureSpec{
			Family:         profile.Family,
			Layers:         profile.LayerCount,
			ContextLength:  profile.MaxContext,
			FileSizeBytes:  profile.FileSizeBytes,
			ParameterCount: profile.ParameterCount,
			QuantType:      profile.FileType,
			BitsPerWeight:  profile.BitsPerWeight,
		}
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
		profile.ArchitectureSpec = spec
	}

	a.settings.MergeModel(map[string]interface{}{"context_size": ctxSize})
}

// ---------------------------------------------------------------------------
// Model directory + file pickers
// ---------------------------------------------------------------------------

// loaderLogDir returns where model-load failure records are written.
func (a *App) loaderLogDir() string {
	if a.loaderDiagnosticsDir != "" {
		return a.loaderDiagnosticsDir
	}
	return loader.LogDir(paths.StateRoot())
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
	a.settings.SetModel(map[string]interface{}{
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
	})
	return nil
}

// UpdatePromptSettings stores prompt pipeline settings.
func (a *App) UpdatePromptSettings(settings *types.PromptSettings) error {
	a.settings.MergeModel(map[string]interface{}{
		"prompt_template":        settings.PromptTemplate,
		"custom_jinja_template":  settings.CustomJinjaTemplate,
		"system_prompt_override": settings.SystemPromptOverride,
		"user_prefix":            settings.UserPrefix,
		"assistant_prefix":       settings.AssistantPrefix,
		"stop_sequences":         settings.StopSequences,
	})
	return nil
}
