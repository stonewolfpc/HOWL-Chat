package kokoro

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

// Synthesizer implements TTSSynthesizer using Kokoro ONNX via Python bridge.
type Synthesizer struct {
	*tts.BaseSynthesizer
	pythonPath string
	scriptPath string
	modelPath  string
	voicesPath string
	configPath string
	outputDir  string

	// Test seams for interpreter discovery, process execution, and output reads.
	resolvePython func() (string, error)
	runCommand    func(context.Context, string, []string, string) ([]byte, error)
	readFile      func(string) ([]byte, error)
}

// New returns a Kokoro synthesizer instance.
func New() *Synthesizer {
	return &Synthesizer{
		BaseSynthesizer: tts.NewBaseSynthesizer(types.TTSKokoro),
	}
}

// Initialize prepares Kokoro model, voices, config, and output directories.
func (s *Synthesizer) Initialize(ctx context.Context, config types.TTSConfig) error {
	s.LockCalls()
	defer s.UnlockCalls()

	if config.ModelPath == "" {
		return types.NewAudioError(types.ErrCodeModelNotFound, "kokoro model path is required", nil)
	}
	if _, err := os.Stat(config.ModelPath); err != nil {
		return types.NewAudioError(types.ErrCodeModelNotFound, "kokoro model not found", err)
	}

	scriptPath := config.ScriptPath
	if scriptPath == "" {
		scriptPath = defaultKokoroScriptPath()
	}
	if _, err := os.Stat(scriptPath); err != nil {
		return types.NewAudioError(types.ErrCodeModelNotFound, "kokoro-tts.py script not found", err)
	}

	resolver := resolvePythonBinary
	if s.resolvePython != nil {
		resolver = s.resolvePython
	}
	pythonPath, err := resolver()
	if err != nil {
		return types.NewAudioError(types.ErrCodeModelNotFound, "python binary not found in PATH", err)
	}

	outputDir := config.CacheDir
	if outputDir == "" {
		outputDir = filepath.Join(os.TempDir(), "howl-kokoro-tts")
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return types.NewAudioError(types.ErrCodeModelLoadFailed, "failed to create kokoro output directory", err)
	}

	// Derive voices and config paths from model path if not explicitly set
	modelDir := filepath.Dir(config.ModelPath)
	voicesPath := config.VoicesPath
	if voicesPath == "" {
		voicesPath = filepath.Join(modelDir, "voices-v1.0.bin")
	}
	configPath := config.ConfigPath
	if configPath == "" {
		configPath = filepath.Join(modelDir, "hf_cache", "config.json")
	}

	// Validate voices and config exist
	if _, err := os.Stat(voicesPath); err != nil {
		return types.NewAudioError(types.ErrCodeModelNotFound, "kokoro voices archive not found", err)
	}
	if _, err := os.Stat(configPath); err != nil {
		return types.NewAudioError(types.ErrCodeModelNotFound, "kokoro config.json not found; run: python scripts/kokoro-tts.py --list-voices to auto-download", err)
	}

	s.pythonPath = pythonPath
	s.scriptPath = scriptPath
	s.modelPath = config.ModelPath
	s.voicesPath = voicesPath
	s.configPath = configPath
	s.outputDir = outputDir
	s.SetInitialized(config)
	return nil
}

// Synthesize converts text to a WAV file via the Python bridge.
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
		callback(5, "Preparing Kokoro synthesis request")
	}

	voice := req.Voice
	if voice == "" {
		voice = "af_heart"
	}
	speed := req.Speed
	if speed <= 0 {
		speed = 1.0
	}

	filename := fmt.Sprintf("kokoro_%d.wav", time.Now().UnixNano())
	outPath := filepath.Join(s.outputDir, filename)

	if callback != nil {
		callback(10, "Running Kokoro ONNX inference")
	}

	if err := s.runKokoro(ctx, req.Text, voice, speed, outPath); err != nil {
		return nil, err
	}

	if callback != nil {
		callback(90, "Reading output metadata")
	}

	meta, err := types.GetAudioMetadata(outPath)
	if err != nil {
		meta = &types.AudioMetadata{
			Format:     types.FormatWAV,
			SampleRate: 24000,
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

// GetVoices returns available voices from the voices archive.
func (s *Synthesizer) GetVoices() []types.Voice {
	s.LockCalls()
	defer s.UnlockCalls()

	if s.voicesPath == "" {
		return nil
	}
	voices, err := s.listVoices()
	if err != nil {
		return nil
	}
	result := make([]types.Voice, 0, len(voices))
	for _, v := range voices {
		result = append(result, types.Voice{
			ID:       v,
			Name:     v,
			Language: "en",
		})
	}
	return result
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

func (s *Synthesizer) runKokoro(ctx context.Context, text, voice string, speed float64, outPath string) error {
	args := []string{
		s.scriptPath,
		"--text", text,
		"--voice", voice,
		"--speed", fmt.Sprintf("%.2f", speed),
		"--output", outPath,
		"--model", s.modelPath,
		"--voices", s.voicesPath,
		"--config", s.configPath,
	}
	if s.runCommand != nil {
		_, err := s.runCommand(ctx, s.pythonPath, args, "")
		return err
	}

	cmd := exec.CommandContext(ctx, s.pythonPath,
		s.scriptPath,
		"--text", text,
		"--voice", voice,
		"--speed", fmt.Sprintf("%.2f", speed),
		"--output", outPath,
		"--model", s.modelPath,
		"--voices", s.voicesPath,
		"--config", s.configPath,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return types.NewAudioError(
			types.ErrCodeSynthesisFailed,
			fmt.Sprintf("kokoro synthesis failed: %s", strings.TrimSpace(string(output))),
			err,
		)
	}
	return nil
}

func (s *Synthesizer) listVoices() ([]string, error) {
	args := []string{s.scriptPath, "--list-voices", "--voices", s.voicesPath}
	var output []byte
	var err error
	if s.runCommand != nil {
		output, err = s.runCommand(context.Background(), s.pythonPath, args, "")
	} else {
		cmd := exec.Command(s.pythonPath, args...)
		output, err = cmd.Output()
	}
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	return lines, nil
}

func (s *Synthesizer) readOutput(path string) ([]byte, error) {
	if s.readFile != nil {
		return s.readFile(path)
	}
	return os.ReadFile(path)
}

func resolvePythonBinary() (string, error) {
	candidates := []string{"python3", "python", "py"}
	if os.Getenv("GOOS") == "windows" || os.PathSeparator == '\\' {
		candidates = []string{"python.exe", "python", "py.exe", "py"}
	}
	for _, name := range candidates {
		if p, err := exec.LookPath(name); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("python not found in PATH")
}

func defaultKokoroScriptPath() string {
	candidates := []string{
		"scripts/kokoro-tts.py",
		"../scripts/kokoro-tts.py",
		"../../scripts/kokoro-tts.py",
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			abs, _ := filepath.Abs(c)
			return abs
		}
	}
	return "scripts/kokoro-tts.py"
}
