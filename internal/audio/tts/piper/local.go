package piper

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"howl-chat/internal/audio/tts"
	"howl-chat/internal/audio/types"
)

// Synthesizer implements TTSSynthesizer using Piper CLI.
type Synthesizer struct {
	*tts.BaseSynthesizer
	piperPath string
	modelPath string
	outputDir string
}

// New returns a Piper synthesizer instance.
func New() *Synthesizer {
	return &Synthesizer{
		BaseSynthesizer: tts.NewBaseSynthesizer(types.TTSPiper),
	}
}

// Initialize prepares Piper binary/model and output directories.
func (s *Synthesizer) Initialize(ctx context.Context, config types.TTSConfig) error {
	if config.ModelPath == "" {
		return types.NewAudioError(types.ErrCodeModelNotFound, "piper model path is required", nil)
	}
	if _, err := os.Stat(config.ModelPath); err != nil {
		return types.NewAudioError(types.ErrCodeModelNotFound, "piper model not found", err)
	}

	piperPath, err := resolvePiperBinary()
	if err != nil {
		return types.NewAudioError(types.ErrCodeModelNotFound, "piper binary not found in PATH", err)
	}

	outputDir := config.CacheDir
	if outputDir == "" {
		outputDir = filepath.Join(os.TempDir(), "howl-tts")
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

	result, err := s.Synthesize(ctx, req, nil)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(result.AudioPath)
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
	return nil
}

func (s *Synthesizer) runPiper(ctx context.Context, req types.SynthesisRequest, outPath string) error {
	args := []string{
		"--model", s.modelPath,
		"--output_file", outPath,
	}
	if req.Speed > 0 {
		ls := 1.0 / clampSpeed(req.Speed)
		args = append(args, "--length_scale", fmt.Sprintf("%.4f", ls))
	}

	cmd := exec.CommandContext(ctx, s.piperPath, args...)
	cmd.Stdin = strings.NewReader(req.Text)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return types.NewAudioError(
			types.ErrCodeSynthesisFailed,
			fmt.Sprintf("piper synthesis failed: %s", strings.TrimSpace(string(output))),
			err,
		)
	}
	return nil
}

func resolvePiperBinary() (string, error) {
	if p, err := exec.LookPath("piper"); err == nil {
		return p, nil
	}
	return exec.LookPath("piper.exe")
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
