package llama

import "testing"

func TestCleanTranscriptStripsLanguageWrapper(t *testing.T) {
	text, language := CleanTranscript("language English<asr_text>The quick brown fox.")
	if text != "The quick brown fox." {
		t.Fatalf("text = %q", text)
	}
	if language != "english" {
		t.Fatalf("language = %q", language)
	}
}

func TestCleanTranscriptLeavesPlainTextAlone(t *testing.T) {
	text, language := CleanTranscript("  plain transcript  ")
	if text != "plain transcript" || language != "" {
		t.Fatalf("text = %q language = %q", text, language)
	}
}

func TestCleanTranscriptHandlesBareMarker(t *testing.T) {
	text, language := CleanTranscript("language zh<asr_text>你好")
	if text != "你好" || language != "zh" {
		t.Fatalf("text = %q language = %q", text, language)
	}
}
