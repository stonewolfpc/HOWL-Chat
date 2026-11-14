package helper

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGoldenFixtures(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("testdata", "golden", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no golden protocol fixtures")
	}
	for _, path := range files {
		path := path
		t.Run(filepath.Base(path), func(t *testing.T) {
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			line := bytes.TrimSuffix(raw, []byte("\n"))
			base := filepath.Base(path)
			switch {
			case strings.HasSuffix(base, ".request.json"):
				if _, err := decodeRequestFrame(line); err != nil {
					t.Fatalf("decodeRequestFrame: %v", err)
				}
			case strings.HasSuffix(base, ".response.json"):
				resp, err := decodeResponseFrame(line)
				if err != nil {
					t.Fatalf("decodeResponseFrame: %v", err)
				}
				validateGoldenResult(t, base, resp)
			default:
				t.Fatalf("unrecognized fixture name %q", base)
			}
		})
	}
}

func validateGoldenResult(t *testing.T, base string, resp responseEnvelope) {
	t.Helper()
	switch base {
	case "hello.response.json":
		var h helloResult
		if err := decodeResult(resp.Result, &h); err != nil {
			t.Fatal(err)
		}
		if err := validateHello(h); err != nil {
			t.Fatal(err)
		}
	case "catalog.response.json":
		var c catalogResult
		if err := decodeResult(resp.Result, &c); err != nil {
			t.Fatal(err)
		}
		if err := validateCatalog(c, catalogParams{ProviderID: "piper", Engine: "piper"}); err != nil {
			t.Fatal(err)
		}
	case "load.response.json":
		var l loadResult
		if err := decodeResult(resp.Result, &l); err != nil {
			t.Fatal(err)
		}
		if err := validateLoad(l, loadParams{
			ProviderID: "piper", Engine: "piper", ModelIdentity: "piper:en_US-amy-medium:v1",
		}); err != nil {
			t.Fatal(err)
		}
	case "synthesize.response.json":
		var s synthesisResult
		if err := decodeResult(resp.Result, &s); err != nil {
			t.Fatal(err)
		}
		if err := validateSynthesisResult(s, synthesisParams{
			ModelIdentity: "piper:en_US-amy-medium:v1",
			Output:        outputSpec{Directory: "C:/Howl/cache/tts", Filename: "c0ffee.wav"},
		}); err != nil {
			t.Fatal(err)
		}
	case "health.response.json":
		var h healthResult
		if err := decodeResult(resp.Result, &h); err != nil {
			t.Fatal(err)
		}
		if err := validateHealth(h); err != nil {
			t.Fatal(err)
		}
	}
}
