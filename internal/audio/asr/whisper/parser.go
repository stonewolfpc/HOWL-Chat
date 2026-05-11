package whisper

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"howl-chat/internal/audio/types"
)

// whisperJSONOutput represents the JSON structure from whisper-cli --output-json
type whisperJSONOutput struct {
	SystemInfo    string                 `json:"systeminfo"`
	Model         whisperModelInfo       `json:"model"`
	Params        whisperParams          `json:"params"`
	Result        whisperResult          `json:"result"`
	Transcription []whisperTranscription `json:"transcription"`
}

// whisperModelInfo contains model metadata
type whisperModelInfo struct {
	Type         string           `json:"type"`
	Multilingual bool             `json:"multilingual"`
	Vocab        int              `json:"vocab"`
	Audio        whisperAudioInfo `json:"audio"`
	Text         whisperAudioInfo `json:"text"`
	Mels         int              `json:"mels"`
	Ftype        int              `json:"ftype"`
}

// whisperAudioInfo contains audio/text context info
type whisperAudioInfo struct {
	Ctx   int `json:"ctx"`
	State int `json:"state"`
	Head  int `json:"head"`
	Layer int `json:"layer"`
}

// whisperParams contains processing parameters
type whisperParams struct {
	Model     string `json:"model"`
	Language  string `json:"language"`
	Translate bool   `json:"translate"`
}

// whisperResult contains detection results
type whisperResult struct {
	Language string `json:"language"`
}

// whisperTranscription represents a single transcription segment
type whisperTranscription struct {
	Timestamps whisperTimestamps `json:"timestamps"`
	Offsets    whisperOffsets    `json:"offsets"`
	Text       string            `json:"text"`
}

// whisperOffsets contains sample-based timing
type whisperOffsets struct {
	From int `json:"from"`
	To   int `json:"to"`
}

// whisperTimestamps contains timing information
type whisperTimestamps struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// whisperSegment represents a single utterance segment
type whisperSegment struct {
	ID               int     `json:"id"`
	Seek             int     `json:"seek"`
	Start            float64 `json:"start"`
	End              float64 `json:"end"`
	Text             string  `json:"text"`
	Tokens           []int   `json:"tokens"`
	Temper           float64 `json:"temperature"`
	AvgLogProb       float64 `json:"avg_logprob"`
	CompressionRatio float64 `json:"compression_ratio"`
	NoSpeechProb     float64 `json:"no_speech_prob"`
}

// parseWhisperOutput reads and parses the JSON output from whisper-cli
func parseWhisperOutput(outputPath string) (*types.RecognitionResult, error) {
	// Read JSON file
	data, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, types.NewAudioError(
			types.ErrCodeModelLoadFailed,
			fmt.Sprintf("Failed to read whisper output: %s", outputPath),
			err,
		)
	}

	// Parse JSON
	var output whisperJSONOutput
	if err := json.Unmarshal(data, &output); err != nil {
		return nil, types.NewAudioError(
			types.ErrCodeModelLoadFailed,
			"Failed to parse whisper JSON output",
			err,
		)
	}

	// Build result from transcription array
	result := &types.RecognitionResult{
		Language: output.Result.Language,
	}

	// Concatenate all transcription text
	var allText strings.Builder
	for _, trans := range output.Transcription {
		cleaned := cleanText(trans.Text)
		if cleaned != "" {
			allText.WriteString(cleaned)
			allText.WriteString(" ")
		}
	}
	result.Text = strings.TrimSpace(allText.String())

	// Calculate duration from timestamps
	if len(output.Transcription) > 0 {
		lastSegment := output.Transcription[len(output.Transcription)-1]
		// Parse timestamp "00:00:06,000" to seconds
		result.Duration = parseTimestampToSeconds(lastSegment.Timestamps.To)

		// Convert transcription to segments
		for _, trans := range output.Transcription {
			cleaned := cleanText(trans.Text)
			if cleaned != "" {
				segment := types.Segment{
					StartTime:  parseTimestampToSeconds(trans.Timestamps.From),
					EndTime:    parseTimestampToSeconds(trans.Timestamps.To),
					Text:       cleaned,
					Confidence: 0.8, // Default confidence since new format doesn't provide it
				}
				result.Segments = append(result.Segments, segment)
			}
		}
		result.Confidence = 0.8 // Default confidence
	} else {
		result.Confidence = 0.8
	}

	return result, nil
}

// parseTimestampToSeconds converts timestamp string "00:00:06,000" to seconds
func parseTimestampToSeconds(ts string) float64 {
	// Format: HH:MM:SS,mmm
	var hours, minutes, seconds, milliseconds int
	fmt.Sscanf(ts, "%d:%d:%d,%d", &hours, &minutes, &seconds, &milliseconds)
	return float64(hours*3600+minutes*60+seconds) + float64(milliseconds)/1000.0
}

// cleanText removes leading/trailing whitespace and normalizes
func cleanText(text string) string {
	// Remove leading/trailing whitespace
	text = strings.TrimSpace(text)

	// Remove leading "0" that whisper sometimes adds
	text = strings.TrimPrefix(text, "0")
	text = strings.TrimSpace(text)

	// Remove duplicate spaces
	text = strings.Join(strings.Fields(text), " ")

	return text
}

// probabilityToConfidence converts log probability to 0-1 confidence score
func probabilityToConfidence(logProb float64) float64 {
	// Whisper log probabilities are typically in range [-2, 0]
	// Map to 0-1 range with sigmoid-like behavior

	// Typical log probs:
	// -0.1 = very confident (~95%)
	// -0.5 = confident (~80%)
	// -1.0 = moderate (~60%)
	// -2.0 = uncertain (~40%)

	if logProb > -0.1 {
		return 0.98
	}
	if logProb > -0.3 {
		return 0.9
	}
	if logProb > -0.6 {
		return 0.8
	}
	if logProb > -1.0 {
		return 0.7
	}
	if logProb > -1.5 {
		return 0.6
	}
	if logProb > -2.0 {
		return 0.5
	}
	return 0.4
}

// calculateAverageConfidence computes overall confidence from segments
func calculateAverageConfidence(segments []types.Segment) float64 {
	if len(segments) == 0 {
		return 0.8
	}

	var total float64
	for _, seg := range segments {
		total += seg.Confidence
	}

	return total / float64(len(segments))
}

