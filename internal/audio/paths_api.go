package audio

import (
	"os"
	"path/filepath"

	"howl-chat/internal/audio/capability"
	"howl-chat/internal/audio/tts/melotts"
	"howl-chat/internal/audio/types"
	"howl-chat/internal/paths"
)

// ASRModelPathFor reports the model path the pipeline would hand to the ASR
// adapter for a model type. Discovery goes through the resolver (nested, flat,
// or HOWL_ASR_MODEL override); when nothing is installed the expected path is
// returned so the capability contract can name where the model belongs. The
// pipeline loads through this same function, so status and loading cannot
// disagree.
func ASRModelPathFor(model types.ASRModelType) string {
	if installed := resolveInstalledASRModel(model); installed != "" {
		return installed
	}
	return defaultASRModelPath(model)
}

// ASRProjectorPathFor reports the audio projector paired with a model. Empty
// means no matching projector was discovered; audio models must not load
// without it.
func ASRProjectorPathFor(model types.ASRModelType) string {
	return resolveInstalledASRProjector(model)
}

// ASRServerPathFor reports the llama-server executable used by the ASR
// adapter. HOWL_ASR_SERVER overrides the bundled runtime.
func ASRServerPathFor() string {
	if override := os.Getenv("HOWL_ASR_SERVER"); override != "" {
		return override
	}
	root := paths.ProjectRoot()
	for _, name := range []string{"llama-server.exe", "llama-server"} {
		candidate := filepath.Join(root, "runtime", "llama", name)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return filepath.Join(root, "runtime", "llama", "llama-server.exe")
}

// TTSModelPathFor reports the model path the pipeline would hand to the TTS
// adapter for a model type.
func TTSModelPathFor(model types.TTSModelType) string {
	return ttsConfigForModel(model, defaultTTSCacheDir()).ModelPath
}

// TTSCacheDir reports the directory synthesized audio is written to. Exported
// so the renderer bridge can confine reads to it; the pipeline writes through
// the same function, so the two cannot disagree.
func TTSCacheDir() string { return defaultTTSCacheDir() }

// KokoroVoicesPathFor and KokoroConfigPathFor expose the same archive and
// vocabulary paths Initialize resolves, so a catalog can report which asset is
// missing without loading a model or starting Python.
func KokoroVoicesPathFor() string { return defaultKokoroVoicesPath() }

func KokoroConfigPathFor() string { return defaultKokoroConfigPath() }

func KokoroScriptPathFor() string { return defaultKokoroScriptPath() }

// ProductionCapabilityRecords is the capability contract the application and
// the pipeline report. It is capability.ProductionEnv with one dimension made
// truthful: MeloTTS readiness is the managed service's ownership plus liveness,
// never a TCP handshake from a process Howl does not own.
func ProductionCapabilityRecords() []capability.Record {
	env := capability.ProductionEnv(ASRModelPathFor, ASRProjectorPathFor, TTSModelPathFor)
	env.MeloTTSReachable = melotts.Readiness
	return capability.Catalog(env)
}
