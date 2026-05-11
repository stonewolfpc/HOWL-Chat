package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DetectAudioFormat identifies format from file extension and header
func DetectAudioFormat(path string) (AudioFormat, error) {
	ext := strings.ToLower(filepath.Ext(path))
	
	switch ext {
	case ".wav":
		return FormatWAV, nil
	case ".mp3":
		return FormatMP3, nil
	case ".ogg":
		return FormatOGG, nil
	case ".flac":
		return FormatFLAC, nil
	case ".m4a", ".aac":
		return FormatM4A, nil
	default:
		// Try to detect from file header
		return detectFromHeader(path)
	}
}

// detectFromHeader reads file magic bytes to identify format
func detectFromHeader(path string) (AudioFormat, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	// Read first 12 bytes
	header := make([]byte, 12)
	if _, err := file.Read(header); err != nil {
		return "", fmt.Errorf("cannot read header: %w", err)
	}

	// Check magic bytes
	if strings.HasPrefix(string(header), "RIFF") && strings.Contains(string(header[8:12]), "WAVE") {
		return FormatWAV, nil
	}
	if strings.HasPrefix(string(header), "ID3") || (header[0] == 0xFF && (header[1]&0xE0) == 0xE0) {
		return FormatMP3, nil
	}
	if strings.HasPrefix(string(header), "OggS") {
		return FormatOGG, nil
	}
	if strings.HasPrefix(string(header), "fLaC") {
		return FormatFLAC, nil
	}

	return "", fmt.Errorf("unknown audio format")
}

// ConvertAudioFormat uses ffmpeg to convert between formats
func ConvertAudioFormat(inputPath string, targetFormat AudioFormat, outputPath string) error {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found in PATH: %w", err)
	}

	// Build ffmpeg command
	args := []string{
		"-i", inputPath,
		"-y", // Overwrite output
	}

	// Add format-specific options
	switch targetFormat {
	case FormatWAV:
		args = append(args, "-ar", "16000", "-ac", "1", "-c:a", "pcm_s16le")
	case FormatMP3:
		args = append(args, "-ar", "44100", "-ac", "2", "-b:a", "128k")
	case FormatOGG:
		args = append(args, "-ar", "44100", "-ac", "2", "-c:a", "libvorbis", "-q:a", "4")
	case FormatFLAC:
		args = append(args, "-ar", "44100", "-ac", "2", "-c:a", "flac")
	}

	args = append(args, outputPath)

	cmd := exec.Command("ffmpeg", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ffmpeg conversion failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// GetAudioMetadata extracts metadata using ffprobe
func GetAudioMetadata(path string) (*AudioMetadata, error) {
	if _, err := exec.LookPath("ffprobe"); err != nil {
		return nil, fmt.Errorf("ffprobe not found in PATH: %w", err)
	}

	// Get format info
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration,bit_rate",
		"-show_entries", "stream=sample_rate,channels",
		"-of", "default=noprint_wrappers=1",
		path,
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe failed: %w", err)
	}

	// Parse output
	metadata := &AudioMetadata{}
	lines := strings.Split(string(output), "\n")
	
	for _, line := range lines {
		if strings.HasPrefix(line, "duration=") {
			fmt.Sscanf(line, "duration=%f", &metadata.Duration)
		} else if strings.HasPrefix(line, "bit_rate=") {
			var br int
			fmt.Sscanf(line, "bit_rate=%d", &br)
			metadata.BitRate = br / 1000 // Convert to kbps
		} else if strings.HasPrefix(line, "sample_rate=") {
			fmt.Sscanf(line, "sample_rate=%d", &metadata.SampleRate)
		} else if strings.HasPrefix(line, "channels=") {
			fmt.Sscanf(line, "channels=%d", &metadata.Channels)
		}
	}

	// Detect format
	format, err := DetectAudioFormat(path)
	if err != nil {
		return nil, err
	}
	metadata.Format = format

	// Estimate bit depth (not directly available from ffprobe)
	if metadata.BitRate > 0 && metadata.SampleRate > 0 && metadata.Channels > 0 {
		// bit_depth = (bit_rate * 1000) / (sample_rate * channels)
		// This is approximate
		metadata.BitDepth = 16 // Assume 16-bit for most formats
	}

	return metadata, nil
}

// ValidateAudioFile checks if file exists and is readable
func ValidateAudioFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewAudioError(ErrCodeModelNotFound, "audio file not found", err)
		}
		return NewAudioError(ErrCodeModelLoadFailed, "cannot access audio file", err)
	}

	if info.IsDir() {
		return NewAudioError(ErrCodeFormatNotSupported, "path is a directory, not a file", nil)
	}

	if info.Size() == 0 {
		return NewAudioError(ErrCodeFormatNotSupported, "audio file is empty", nil)
	}

	// Check format support
	format, err := DetectAudioFormat(path)
	if err != nil {
		return NewAudioError(ErrCodeFormatNotSupported, "unsupported audio format", err)
	}

	// Check format is in supported list
	supportedFormats := []AudioFormat{FormatWAV, FormatMP3, FormatOGG, FormatFLAC, FormatM4A}
	found := false
	for _, f := range supportedFormats {
		if f == format {
			found = true
			break
		}
	}
	if !found {
		return NewAudioError(ErrCodeFormatNotSupported, 
			fmt.Sprintf("format %s not in supported list", format), nil)
	}

	return nil
}

// GetProcessingFormat returns the internal processing format
func GetProcessingFormat() AudioMetadata {
	return AudioMetadata{
		Format:     FormatWAV,
		SampleRate: 16000,
		Channels:   1,
		BitDepth:   16,
	}
}

// CalculateAudioTokens estimates token count for audio
// 1 second ≈ 50-100 tokens depending on model
func CalculateAudioTokens(duration float64, modelType string) int {
	// Different models have different token densities
	tokensPerSecond := 50.0 // Default

	switch modelType {
	case string(ASRWhisperTiny), string(ASRWhisperBase):
		tokensPerSecond = 50
	case string(ASRWhisperSmall), string(ASRWhisperMed):
		tokensPerSecond = 75
	case string(ASRQwen3ASR06B), string(ASRQwen3ASR17B):
		tokensPerSecond = 100
	case string(MultimodalQwen25Omni7B), string(MultimodalQwen3Omni30B):
		tokensPerSecond = 100
	}

	return int(duration * tokensPerSecond)
}

// EstimateMemoryUsage calculates expected memory for audio processing
func EstimateMemoryUsage(modelType string, audioDuration float64) (int64, error) {
	profile, err := GetModelInfo(modelType)
	if err != nil {
		return 0, err
	}

	// Base model memory
	memory := profile.RAMBytes

	// Add audio buffer memory (assume 16-bit mono 16kHz)
	audioBytes := int64(audioDuration * 16000 * 2) // 2 bytes per sample

	// Add working memory (varies by model)
	workingMemory := int64(100 * 1024 * 1024) // 100MB default

	switch profile.Type {
	case "asr":
		workingMemory = int64(200 * 1024 * 1024) // 200MB for ASR
	case "tts":
		workingMemory = int64(100 * 1024 * 1024) // 100MB for TTS
	case "multimodal":
		workingMemory = int64(500 * 1024 * 1024) // 500MB for multimodal
	}

	return memory + audioBytes + workingMemory, nil
}

// CheckHardwareCompatibility verifies if model can run on this hardware
func CheckHardwareCompatibility(modelType string, availableRAM int64) (bool, error) {
	requiredMemory, err := EstimateMemoryUsage(modelType, 300) // 5 min audio
	if err != nil {
		return false, err
	}

	// Add 20% safety margin
	requiredWithMargin := int64(float64(requiredMemory) * 1.2)

	return availableRAM >= requiredWithMargin, nil
}

// GenerateCacheKey creates a deterministic cache key for TTS
func GenerateCacheKey(text string, voice string, speed float64) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%.6f", voice, text, speed)))
	return hex.EncodeToString(sum[:])
}

// FormatDuration converts seconds to human-readable string
func FormatDuration(seconds float64) string {
	if seconds < 60 {
		return fmt.Sprintf("%.1fs", seconds)
	}
	minutes := int(seconds) / 60
	secs := int(seconds) % 60
	return fmt.Sprintf("%dm%ds", minutes, secs)
}

// FormatBytes converts bytes to human-readable string
func FormatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.0f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// IsRealTimeCapable checks if model can process audio in real-time
func IsRealTimeCapable(modelType string) bool {
	profile, err := GetModelInfo(modelType)
	if err != nil {
		return false
	}

	// Real-time factor < 1.0 means faster than realtime
	return profile.RealTimeFactor < 1.0
}

// GetSupportedLanguages returns languages supported by a model
func GetSupportedLanguages(modelType string) ([]string, error) {
	profile, err := GetModelInfo(modelType)
	if err != nil {
		return nil, err
	}

	return profile.Languages, nil
}

// SupportsLanguage checks if model supports a specific language
func SupportsLanguage(modelType string, language string) bool {
	languages, err := GetSupportedLanguages(modelType)
	if err != nil {
		return false
	}

	// Normalize language code
	lang := strings.ToLower(language)
	if len(lang) > 2 {
		lang = lang[:2]
	}

	for _, l := range languages {
		l = strings.ToLower(l)
		if len(l) > 2 {
			l = l[:2]
		}
		if l == lang {
			return true
		}
	}

	return false
}
