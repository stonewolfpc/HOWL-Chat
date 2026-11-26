package piper

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"howl-chat/internal/audio/tts"
	"howl-chat/internal/audio/types"
	"howl-chat/internal/paths"
	"howl-chat/internal/process"
)

// Synthesizer implements TTSSynthesizer using Piper CLI.
type Synthesizer struct {
	*tts.BaseSynthesizer
	piperPath string
	modelPath string
	outputDir string
	registry  *process.Registry
	childSeq  atomic.Uint64

	// Test seam for process execution. Production always starts a job-owned,
	// registered child; only provider unit tests replace the runner.
	resolveBinary func() (string, error)
	runCommand    func(context.Context, string, []string, string) ([]byte, error)
	readFile      func(string) ([]byte, error)
}

// New returns a Piper synthesizer instance. A nil registry falls back to the
// process-wide one so ownership is never optional.
func New(registries ...*process.Registry) *Synthesizer {
	var registry *process.Registry
	if len(registries) > 0 {
		registry = registries[0]
	}
	return &Synthesizer{
		BaseSynthesizer: tts.NewBaseSynthesizer(types.TTSPiper),
		registry:        registry,
	}
}

// Initialize prepares Piper binary/model and output directories.
func (s *Synthesizer) Initialize(ctx context.Context, config types.TTSConfig) error {
	s.LockCalls()
	defer s.UnlockCalls()

	if config.ModelPath == "" {
		return types.NewAudioError(types.ErrCodeModelNotFound, "piper model path is required", nil)
	}
	if _, err := os.Stat(config.ModelPath); err != nil {
		return types.NewAudioError(types.ErrCodeModelNotFound, "piper model not found", err)
	}

	resolver := resolvePiperBinary
	if s.resolveBinary != nil {
		resolver = s.resolveBinary
	}
	piperPath, err := resolver()
	if err != nil {
		return types.NewAudioError(types.ErrCodeModelNotFound, "piper binary not found in PATH", err)
	}

	outputDir := config.CacheDir
	if outputDir == "" {
		outputDir = filepath.Join(paths.TempDir(), "howl-tts")
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return types.NewAudioError(types.ErrCodeModelLoadFailed, "failed to create tts output directory", err)
	}

	s.piperPath = piperPath
	s.modelPath = config.ModelPath
	s.outputDir = outputDir
	s.SetInitialized(config)
	return nil
}

// Synthesize converts text to a WAV file.
func (s *Synthesizer) Synthesize(ctx context.Context, req types.SynthesisRequest, callback types.ProgressCallback) (*types.SynthesisResult, error) {
	s.LockCalls()
	defer s.UnlockCalls()
	return s.synthesizeLocked(ctx, req, callback)
}

func (s *Synthesizer) synthesizeLocked(ctx context.Context, req types.SynthesisRequest, callback types.ProgressCallback) (*types.SynthesisResult, error) {
	if err := s.CheckInitialized(); err != nil {
		return nil, err
	}
	if err := s.ValidateRequest(req); err != nil {
		return nil, err
	}

	if callback != nil {
		callback(5, "Preparing synthesis request")
	}

	filename := fmt.Sprintf("tts_%d.wav", time.Now().UnixNano())
	outPath := filepath.Join(s.outputDir, filename)
	if err := s.runPiper(ctx, req, outPath); err != nil {
		return nil, err
	}

	if callback != nil {
		callback(90, "Reading output metadata")
	}

	meta, err := types.GetAudioMetadata(outPath)
	if err != nil {
		meta = &types.AudioMetadata{
			Format:     types.FormatWAV,
			SampleRate: 22050,
			Channels:   1,
		}
	}

	if callback != nil {
		callback(100, "Complete")
	}
	return &types.SynthesisResult{
		AudioPath: outPath,
		Metadata:  *meta,
		Duration:  meta.Duration,
	}, nil
}

// SynthesizeStream synthesizes then streams bytes as chunks.
func (s *Synthesizer) SynthesizeStream(ctx context.Context, req types.SynthesisRequest, chunks chan<- types.AudioChunk) error {
	defer close(chunks)

	s.LockCalls()
	result, err := s.synthesizeLocked(ctx, req, nil)
	if err != nil {
		s.UnlockCalls()
		return err
	}
	data, err := s.readOutput(result.AudioPath)
	s.UnlockCalls()
	if err != nil {
		return types.NewAudioError(types.ErrCodeSynthesisFailed, "failed to read synthesized audio", err)
	}

	const chunkSize = 32 * 1024
	for i, offset := 0, 0; offset < len(data); i, offset = i+1, offset+chunkSize {
		end := offset + chunkSize
		if end > len(data) {
			end = len(data)
		}
		select {
		case <-ctx.Done():
			return types.NewAudioError(types.ErrCodeCancelled, "stream synthesis cancelled", ctx.Err())
		case chunks <- types.AudioChunk{
			Data:     data[offset:end],
			Metadata: result.Metadata,
			Index:    i,
			IsLast:   end == len(data),
		}:
		}
	}

	return nil
}

// GetVoices returns profile voices from model registry.
func (s *Synthesizer) GetVoices() []types.Voice {
	s.LockCalls()
	defer s.UnlockCalls()

	profile, ok := types.GetTTSProfile(types.TTSPiper)
	if !ok {
		return nil
	}
	voices := make([]types.Voice, 0, len(profile.SupportedVoices))
	for _, v := range profile.SupportedVoices {
		voices = append(voices, types.Voice{
			ID:       v,
			Name:     v,
			Language: "en",
		})
	}
	return voices
}

// SupportsSSML reports SSML support.
func (s *Synthesizer) SupportsSSML() bool {
	return false
}

// Release performs synthesizer cleanup.
func (s *Synthesizer) Release() error {
	s.LockCalls()
	defer s.UnlockCalls()
	return nil
}

func (s *Synthesizer) runPiper(ctx context.Context, req types.SynthesisRequest, outPath string) error {
	modelPath := s.modelPath
	if req.VoicePath != "" {
		modelPath = req.VoicePath
	}
	args := []string{
		"--model", modelPath,
		"--output_file", outPath,
	}
	if req.Speed > 0 {
		ls := 1.0 / clampSpeed(req.Speed)
		args = append(args, "--length_scale", fmt.Sprintf("%.4f", ls))
	}

	if s.runCommand != nil {
		_, err := s.runCommand(ctx, s.piperPath, args, req.Text)
		return err
	}

	// One registered, job-owned child per synthesis. Cancellation force-kills
	// the whole tree and the provider does not return before the registered
	// waiter observed exit; an unconfirmed exit is an error, not success.
	child, err := process.StartChild(s.registry, process.ChildSpec{
		Role:     process.RoleTTSChild,
		Instance: "piper|" + modelPath + "|" + strconv.FormatUint(s.childSeq.Add(1), 10),
		Path:     s.piperPath,
		Args:     args,
		Stdin:    []byte(req.Text),
		Capture:  true,
	})
	if err != nil {
		return types.NewAudioError(
			types.ErrCodeSynthesisFailed,
			fmt.Sprintf("piper could not be started under ownership: %v", err),
			err,
		)
	}
	if err := child.Wait(ctx, process.DefaultGrace); err != nil {
		return types.NewAudioError(
			types.ErrCodeSynthesisFailed,
			fmt.Sprintf("piper synthesis failed: %v: %s", err, strings.TrimSpace(child.Output())),
			err,
		)
	}
	return nil
}

func resolvePiperBinary() (string, error) {
	if p, err := exec.LookPath("piper"); err == nil {
		return p, nil
	}
	if p, err := exec.LookPath("piper.exe"); err == nil {
		return p, nil
	}
	// Fallback: check project bin/ directory (common in dev builds)
	candidates := []string{
		"bin/piper.exe",
		"../bin/piper.exe",
		"../../bin/piper.exe",
		"bin/piper",
		"../bin/piper",
		"../../bin/piper",
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			abs, _ := filepath.Abs(c)
			return abs, nil
		}
	}
	return "", fmt.Errorf("piper binary not found in PATH or project bin/ directory")
}

func (s *Synthesizer) readOutput(path string) ([]byte, error) {
	if s.readFile != nil {
		return s.readFile(path)
	}
	return os.ReadFile(path)
}

func clampSpeed(speed float64) float64 {
	if speed < 0.5 {
		return 0.5
	}
	if speed > 2.0 {
		return 2.0
	}
	return speed
}
