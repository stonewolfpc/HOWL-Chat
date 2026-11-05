package audio

import (
	"os"
	"testing"

	"howl-chat/internal/audio/capability"
	"howl-chat/internal/audio/types"
)

// TestMain installs a fully selectable catalog for the package's pipeline
// tests. Those tests exercise generations, leases, lifecycle and provider
// outcomes with scripted loaders; they are not evidence about installed audio
// assets. Production reads the probed environment through the same seam.
func TestMain(m *testing.M) {
	pipelineCapabilityCatalog = selectableCatalogForTests
	os.Exit(m.Run())
}

func selectableCatalogForTests() []capability.Record {
	entries := []struct {
		id   string
		task capability.Task
	}{
		{string(types.ASRWhisperTiny), capability.TaskFileTranscription},
		{string(types.ASRWhisperBase), capability.TaskFileTranscription},
		{string(types.ASRWhisperSmall), capability.TaskFileTranscription},
		{string(types.TTSPiper), capability.TaskSpeechSynthesis},
		{string(types.TTSMeloTTS), capability.TaskSpeechSynthesis},
		{string(types.TTSKokoro), capability.TaskSpeechSynthesis},
	}

	records := make([]capability.Record, 0, len(entries))
	for _, entry := range entries {
		records = append(records, capability.Record{
			ID:            entry.id,
			Task:          entry.task,
			Local:         true,
			ProviderState: capability.Avail("scripted test provider"),
			ModelState:    capability.Avail("scripted test model"),
			Projector:     capability.NoRequirement("scripted test needs no projector"),
			Runtime:       capability.Avail("scripted test runtime"),
			TaskSupport:   capability.Avail("scripted test task"),
			Hardware:      capability.Avail("scripted test hardware"),
			Verification:  types.VerificationImplementation,
		})
	}
	return records
}
