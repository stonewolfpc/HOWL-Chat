package sdcpp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"howl-chat/internal/image/types"
)

func TestDGLUpscaleScriptUsesForwardSlashes(t *testing.T) {
	source, err := buildDGLUpscaleScript(`C:\models\input.png`, `C:\models\2x-ESRGAN.gguf`, 128)
	if err != nil {
		t.Fatalf("build script: %v", err)
	}
	if !strings.Contains(source, `input_path = "C:/models/input.png"`) {
		t.Fatalf("input path must be emitted with forward slashes, got:\n%s", source)
	}
	if !strings.Contains(source, `upscaler_path = "C:/models/2x-ESRGAN.gguf"`) {
		t.Fatalf("upscaler path must be emitted with forward slashes, got:\n%s", source)
	}
	if strings.Contains(source, `\U`) || strings.Contains(source, `\m`) {
		t.Fatalf("script contains a backslash escape DGL rejected:\n%s", source)
	}
	if !strings.Contains(source, "load_upscaler(upscaler_path, tile_size = tile_size)") {
		t.Fatalf("script must load the upscaler explicitly:\n%s", source)
	}
}

func TestDGLUpscaleScriptRejectsQuoteInPath(t *testing.T) {
	if _, err := buildDGLUpscaleScript(`C:\bad"path.png`, `C:\up.gguf`, 128); err == nil {
		t.Fatal("a quote in a path must be refused rather than emitted into a script")
	}
}

func TestDGLShapeDimensions(t *testing.T) {
	if w, h := dglShapeWidthHeight([]int{3, 1024, 512}); w != 512 || h != 1024 {
		t.Fatalf("CHW shape gave %dx%d, want 512x1024", w, h)
	}
	if w, h := dglShapeWidthHeight([]int{512, 1024, 3}); w != 1024 || h != 512 {
		t.Fatalf("HWC shape gave %dx%d, want 1024x512", w, h)
	}
}

// A DGL stage binds its generation context to the workflow caller, so cancelling
// that context reaches the active backend request instead of being ignored.
func TestWorkflowActiveStageObservesCallerCancellation(t *testing.T) {
	entered := make(chan struct{}, 1)
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/sdcpp/v1/dgl/jobs":
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte(`{"id":"job_cancel","status":"queued"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/sdcpp/v1/dgl/jobs/job_cancel":
			select {
			case entered <- struct{}{}:
			default:
			}
			<-release
			_, _ = w.Write([]byte(`{"id":"job_cancel","status":"completed","result":{"outputs":[{"name":"image","shape":[3,64,64],"dtype":1,"bytes":8,"frames":1}]}}`))
		case r.Method == http.MethodGet && r.URL.Path == "/sdcpp/v1/dgl/jobs/job_cancel/outputs/0":
			_, _ = w.Write([]byte("\x89PNG\r\n\x1a\n"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	engine := &Engine{httpClient: newClient(server.URL)}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, _, err := engine.DGLUpscale(ctx, `C:\in.png`, `C:\up.gguf`, 128)
		done <- err
	}()

	select {
	case <-entered:
	case <-time.After(2 * time.Second):
		t.Fatal("DGL stage never polled its job")
	}
	cancel()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("cancelled DGL stage returned no error")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("cancelling the caller context did not reach the active DGL stage")
	}

	close(release)
}

func TestDGLUpscaleRunsJobAndFetchesOutput(t *testing.T) {
	var submitted string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/sdcpp/v1/dgl/jobs":
			var body struct {
				Source string `json:"source"`
			}
			_ = json.NewDecoder(r.Body).Decode(&body)
			submitted = body.Source
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte(`{"id":"job_test","status":"queued"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/sdcpp/v1/dgl/jobs/job_test":
			_, _ = w.Write([]byte(`{"id":"job_test","status":"completed","result":{"outputs":[{"name":"image","shape":[3,1024,1024],"dtype":1,"bytes":8,"frames":1}]}}`))
		case r.Method == http.MethodGet && r.URL.Path == "/sdcpp/v1/dgl/jobs/job_test/outputs/0":
			_, _ = w.Write([]byte("\x89PNG\r\n\x1a\n"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	engine := &Engine{httpClient: newClient(server.URL)}
	data, report, err := engine.DGLUpscale(context.Background(), `C:\in.png`, `C:\up.gguf`, 128)
	if err != nil {
		t.Fatalf("DGLUpscale: %v", err)
	}
	if string(data) != "\x89PNG\r\n\x1a\n" {
		t.Fatalf("unexpected output bytes: %q", data)
	}
	if report.JobID != "job_test" || report.OutputWidth != 1024 || report.OutputHeight != 1024 {
		t.Fatalf("unexpected report: %+v", report)
	}
	if !strings.Contains(submitted, "C:/in.png") || strings.Contains(submitted, `C:\`) {
		t.Fatalf("submitted script did not normalize the path:\n%s", submitted)
	}
	if engine.GenerationState() != "idle" {
		t.Fatalf("generation slot was not released: %s", engine.GenerationState())
	}
}

func TestSDCPPJobProgressAcceptsObjectShape(t *testing.T) {
	var job sdcppJob
	body := `{"id":"job_1","status":"generating","progress":{"completed":2,"total":8,"stage":"sampling"},"result":{"images":[]}}`
	if err := json.Unmarshal([]byte(body), &job); err != nil {
		t.Fatalf("object progress must decode: %v", err)
	}
	if got, want := float64(job.Progress), 0.25; got != want {
		t.Fatalf("progress = %v, want %v", got, want)
	}
	number := `{"id":"job_2","status":"generating","progress":0.5}`
	if err := json.Unmarshal([]byte(number), &job); err != nil {
		t.Fatalf("numeric progress must still decode: %v", err)
	}
	if got, want := float64(job.Progress), 0.5; got != want {
		t.Fatalf("numeric progress = %v, want %v", got, want)
	}
}

func TestImg2ImgSyncUsesCompatibilityRouteAndHonorsSteps(t *testing.T) {
	var gotPath string
	var gotSteps int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		var body struct {
			Steps int `json:"steps"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		gotSteps = body.Steps
		_, _ = w.Write([]byte(`{"images":["aGVsbG8="],"info":""}`))
	}))
	defer server.Close()

	c := newClient(server.URL)
	result, err := c.img2imgSync(context.Background(), types.ImageOptions{
		Prompt: "test", Width: 512, Height: 512, Steps: 8, Strength: 0.5,
		InitImagePath: "dgl_test.go",
	})
	if err != nil {
		t.Fatalf("img2imgSync: %v", err)
	}
	if gotPath != "/sdapi/v1/img2img" {
		t.Fatalf("sync route = %q, want /sdapi/v1/img2img", gotPath)
	}
	if gotSteps != 8 {
		t.Fatalf("requested steps = %d, want 8", gotSteps)
	}
	if result.ImageData != "aGVsbG8=" {
		t.Fatalf("unexpected image data: %q", result.ImageData)
	}
}

func TestNativeImg2ImgCarriesControlAndNestedSteps(t *testing.T) {
	dir := t.TempDir()
	initPath := filepath.Join(dir, "init.png")
	controlPath := filepath.Join(dir, "pose.png")
	if err := os.WriteFile(initPath, []byte("init-bytes"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(controlPath, []byte("control-bytes"), 0o600); err != nil {
		t.Fatal(err)
	}

	var gotPath string
	var got nativeImg2ImgRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/sdcpp/v1/img_gen":
			gotPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&got)
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte(`{"id":"job_native","status":"queued"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/sdcpp/v1/jobs/job_native":
			_, _ = w.Write([]byte(`{"id":"job_native","status":"completed","result":{"images":[{"b64_json":"aGVsbG8=","index":0}]}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	c := newClient(server.URL)
	result, err := c.img2imgNativeControl(context.Background(), types.ImageOptions{
		Prompt: "test", Width: 512, Height: 512, Steps: 8, Sampler: "euler_a",
		Scheduler: "discrete", CfgScale: 7, Strength: 0.5, BatchCount: 1,
		InitImagePath: initPath, ControlImagePath: controlPath,
		ControlNetPath: filepath.Join(dir, "control.safetensors"), ControlStrength: 0.9,
	})
	if err != nil {
		t.Fatalf("img2imgNativeControl: %v", err)
	}
	if gotPath != "/sdcpp/v1/img_gen" {
		t.Fatalf("native route = %q, want /sdcpp/v1/img_gen", gotPath)
	}
	if got.SampleParams.SampleSteps != 8 {
		t.Fatalf("sample_params.sample_steps = %d, want 8", got.SampleParams.SampleSteps)
	}
	if got.SampleParams.SampleMethod != "euler_a" || got.SampleParams.Scheduler != "discrete" {
		t.Fatalf("sample params were not forwarded: %+v", got.SampleParams)
	}
	if got.SampleParams.Guidance.TxtCFG != 7 {
		t.Fatalf("cfg was not nested under guidance: %+v", got.SampleParams.Guidance)
	}
	if want := base64.StdEncoding.EncodeToString([]byte("control-bytes")); got.ControlImage != want {
		t.Fatalf("control image was not sent as base64")
	}
	if !strings.HasSuffix(got.ControlNetPath, "/control.safetensors") || strings.Contains(got.ControlNetPath, `\`) {
		t.Fatalf("control net path was not normalized: %q", got.ControlNetPath)
	}
	if result.ImageData != "aGVsbG8=" {
		t.Fatalf("b64_json image was not decoded: %q", result.ImageData)
	}
}

func TestDGLFailedJobReleasesGenerationSlot(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost:
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte(`{"id":"job_fail","status":"queued"}`))
		case r.URL.Path == "/sdcpp/v1/dgl/jobs/job_fail":
			_, _ = w.Write([]byte(`{"id":"job_fail","status":"failed","error":{"code":"generation_failed","message":"upscaler loading failed"}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	engine := &Engine{httpClient: newClient(server.URL)}
	_, report, err := engine.DGLUpscale(context.Background(), `C:\in.png`, `C:\up.gguf`, 128)
	if err == nil {
		t.Fatal("failed DGL job must return an error")
	}
	if report.JobStatus != "failed" {
		t.Fatalf("job status = %q, want failed", report.JobStatus)
	}
	if engine.GenerationState() != "idle" {
		t.Fatalf("a definite DGL failure must not leave the engine abandoned: %s", engine.GenerationState())
	}
}
