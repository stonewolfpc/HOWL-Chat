package sdcpp

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"howl-chat/internal/image/types"
)

func TestVideoContainer(t *testing.T) {
	ok := map[string]string{"": "webm", "webm": "webm", "AVI": "avi", "webp": "webp"}
	for in, want := range ok {
		got, err := videoContainer(in)
		if err != nil || got != want {
			t.Fatalf("videoContainer(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	for _, bad := range []string{"mp4", "gif", "mkv"} {
		if _, err := videoContainer(bad); err == nil {
			t.Fatalf("videoContainer(%q) should be rejected", bad)
		}
	}
}

func TestIsVideoFamily(t *testing.T) {
	for _, f := range []string{"wan", "WAN2.2", "ltx2", "hunyuanvideo", "animatediff"} {
		if !IsVideoFamily(f) {
			t.Errorf("IsVideoFamily(%q) = false, want true", f)
		}
	}
	for _, f := range []string{"sdxl", "flux", "sd1.5", ""} {
		if IsVideoFamily(f) {
			t.Errorf("IsVideoFamily(%q) = true, want false", f)
		}
	}
}

func TestBuildVideoArgsWan(t *testing.T) {
	e := &Engine{currentModel: &types.ModelInfo{
		Name:   "wan21",
		Family: "wan",
		Files: map[string]string{
			"model": "m.gguf",
			"vae":   "v.safetensors",
			"t5xxl": "t.gguf",
		},
	}}
	args := strings.Join(e.buildVideoArgs(types.ImageOptions{Prompt: "a fox", Seed: 7}, 33, 16, 832, 480, "out.webm"), " ")
	wants := []string{
		"-M vid_gen",
		"--diffusion-model m.gguf",
		"--vae v.safetensors",
		"--t5xxl t.gguf",
		"--prompt a fox",
		"-W 832",
		"-H 480",
		"--video-frames 33",
		"--fps 16",
		"--seed 7",
		"-o out.webm",
	}
	for _, want := range wants {
		if !strings.Contains(args, want) {
			t.Fatalf("args missing %q\nfull: %s", want, args)
		}
	}
}

func TestBuildVideoArgsAnimateDiffUsesSingleModel(t *testing.T) {
	e := &Engine{currentModel: &types.ModelInfo{
		Name:   "sd15-animatediff",
		Family: "animatediff",
		Files: map[string]string{
			"model":         "sd15.safetensors",
			"motion_module": "mm.ckpt",
		},
	}}
	args := strings.Join(e.buildVideoArgs(types.ImageOptions{Prompt: "wave"}, 16, 8, 512, 512, "o.webm"), " ")
	if !strings.Contains(args, "-m sd15.safetensors") {
		t.Fatalf("AnimateDiff should load with -m, got: %s", args)
	}
	if !strings.Contains(args, "--motion-module mm.ckpt") {
		t.Fatalf("missing motion module: %s", args)
	}
	if strings.Contains(args, "--diffusion-model") {
		t.Fatalf("AnimateDiff must not use --diffusion-model: %s", args)
	}
}

func TestParseResolution(t *testing.T) {
	if w, h, ok := parseResolution("832x480"); !ok || w != 832 || h != 480 {
		t.Fatalf("parseResolution(832x480) = %d,%d,%v", w, h, ok)
	}
	if _, _, ok := parseResolution("nonsense"); ok {
		t.Fatal("parseResolution should reject junk")
	}
	if _, _, ok := parseResolution("10x480"); ok {
		t.Fatal("parseResolution should reject absurd sizes")
	}
}

func TestBuildVideoArgsAlwaysTilesVAE(t *testing.T) {
	e := &Engine{currentModel: &types.ModelInfo{
		Name:   "wan22-ti2v-5b",
		Family: "wan",
		Files:  map[string]string{"model": "m.gguf", "vae": "v.safetensors", "t5xxl": "t.gguf"},
	}}
	args := strings.Join(e.buildVideoArgs(types.ImageOptions{Prompt: "a fox"}, 33, 16, 832, 480, "o.webm"), " ")
	if !strings.Contains(args, "--vae-tiling") {
		t.Fatalf("video args must always tile the VAE: %s", args)
	}
	if n := strings.Count(args, "--vae-tiling"); n != 1 {
		t.Fatalf("--vae-tiling emitted %d times, want 1: %s", n, args)
	}
}

func TestVideoResolutionFallsBackPerFamily(t *testing.T) {
	w, h := videoResolution(&types.ModelInfo{Family: "wan"}, "wan")
	if w != 832 || h != 480 {
		t.Fatalf("wan default = %dx%d, want 832x480", w, h)
	}
	w, h = videoResolution(&types.ModelInfo{Family: "wan", ResolutionHint: "640x384"}, "wan")
	if w != 640 || h != 384 {
		t.Fatalf("hint ignored: %dx%d", w, h)
	}
}

func TestRecordVideoProgress(t *testing.T) {
	e := &Engine{}
	e.recordVideoProgress("  4/20 - 1.20s/it")
	if got := e.progress.Progress; got < 0.19 || got > 0.21 {
		t.Fatalf("progress = %v, want ~0.2", got)
	}
	if e.progress.Step != 4 || e.progress.TotalSteps != 20 {
		t.Fatalf("step = %d/%d, want 4/20", e.progress.Step, e.progress.TotalSteps)
	}
	e.recordVideoProgress("loading model...")
	if e.progress.Step != 4 {
		t.Fatal("unrelated log lines must not change progress")
	}
}

func TestScanAnyEOLHandlesCarriageReturns(t *testing.T) {
	input := "first\rsecond\nthird"
	var got []string
	rest := []byte(input)
	for {
		adv, tok, err := scanAnyEOL(rest, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if adv == 0 {
			break
		}
		got = append(got, string(tok))
		rest = rest[adv:]
	}
	if len(got) != 2 || got[0] != "first" || got[1] != "second" {
		t.Fatalf("scanAnyEOL tokens = %q", got)
	}
}

func TestReadAsBase64RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "clip.webm")
	if err := os.WriteFile(path, []byte("hello-video"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	got, err := readAsBase64(path)
	if err != nil {
		t.Fatalf("readAsBase64: %v", err)
	}
	raw, err := base64.StdEncoding.DecodeString(got)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if string(raw) != "hello-video" {
		t.Fatalf("round trip mismatch: %q", raw)
	}
}

func TestVideoMIME(t *testing.T) {
	cases := map[string]string{"webm": "video/webm", "avi": "video/x-msvideo", "webp": "image/webp", "": "video/webm"}
	for ext, want := range cases {
		if got := videoMIME(ext); got != want {
			t.Errorf("videoMIME(%q) = %q, want %q", ext, got, want)
		}
	}
}
