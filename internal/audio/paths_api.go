package audio

import (
	"os"
	"path/filepath"

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
