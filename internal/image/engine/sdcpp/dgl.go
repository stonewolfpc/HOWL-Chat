package sdcpp

// DGL is sd-server's small scripting surface. The only operation Howl needs
// from it right now is ESRGAN upscaling: the native HTTP API exposes upscaling
// only as part of a generation request, while DGL can run it as its own job
// between two native generations.
//
// The contract below was measured against the installed sd-server:
//   POST /sdcpp/v1/dgl/jobs            {"source": "..."} -> 202 + job
//   GET  /sdcpp/v1/dgl/jobs/{id}       job status/progress/result
//   GET  /sdcpp/v1/dgl/jobs/{id}/outputs/{index}?frame=N -> image/png
//
// A completed job's result.outputs carries the metadata needed to pick the
// image output and its dimensions. Jobs expire after 600s, so the bytes are
// fetched as soon as completion is observed and the caller owns them after
// that.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// DGLJobError is the error object in a failed job.
type DGLJobError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// DGLJobOutput describes one value returned by a DGL script.
type DGLJobOutput struct {
	Name   string `json:"name"`
	Shape  []int  `json:"shape"`
	Dtype  int    `json:"dtype"`
	Bytes  int    `json:"bytes"`
	Frames int    `json:"frames"`
}

// DGLJobResult is the result object of a completed job.
type DGLJobResult struct {
	Type    string           `json:"type"`
	Outputs []DGLJobOutput   `json:"outputs"`
	Media   []map[string]any `json:"media"`
}

// DGLJob is the status document returned by submit and poll.
type DGLJob struct {
	ID        string        `json:"id"`
	Kind      string        `json:"kind"`
	Status    string        `json:"status"`
	Error     *DGLJobError  `json:"error"`
	Result    *DGLJobResult `json:"result"`
	Created   int64         `json:"created"`
	Started   int64         `json:"started"`
	Completed int64         `json:"completed"`
}

// DGLUpscaleReport records what the job actually did, for workflow stage
// reporting. Width/Height come from the job's own output shape, not from a
// planned scale factor.
type DGLUpscaleReport struct {
	JobID         string `json:"job_id"`
	InputWidth    int    `json:"input_width,omitempty"`
	InputHeight   int    `json:"input_height,omitempty"`
	OutputWidth   int    `json:"output_width"`
	OutputHeight  int    `json:"output_height"`
	DurationMS    int64  `json:"duration_ms"`
	OutputBytes   int    `json:"output_bytes"`
	UpscalerPath  string `json:"upscaler_path"`
	UpscalerTile  int    `json:"upscaler_tile"`
	Executor      string `json:"executor"`
	JobStatus     string `json:"job_status"`
	BackendReport string `json:"backend_report,omitempty"`
}

// DGLUpscale runs one ESRGAN upscale as a DGL job and returns the PNG bytes.
func (e *Engine) DGLUpscale(ctx context.Context, inputPath, upscalerPath string, tileSize int) ([]byte, DGLUpscaleReport, error) {
	report := DGLUpscaleReport{
		UpscalerPath: upscalerPath,
		UpscalerTile: tileSize,
		Executor:     "dgl",
	}
	if e.httpClient == nil {
		return nil, report, fmt.Errorf("engine not initialized")
	}
	if inputPath == "" {
		return nil, report, fmt.Errorf("DGL upscale requires an input image")
	}
	if upscalerPath == "" {
		return nil, report, fmt.Errorf("DGL upscale requires an upscaler model")
	}
	if tileSize <= 0 {
		tileSize = 128
	}
	source, err := dglUpscaleSource(inputPath, upscalerPath, tileSize)
	if err != nil {
		return nil, report, err
	}

	h, err := e.beginGenerationWithContext("", ctx)
	if err != nil {
		return nil, report, err
	}
	started := time.Now()
	job, err := e.httpClient.dglSubmit(h.ctx, source)
	if err == nil {
		report.JobID = job.ID
		job, err = e.httpClient.dglWait(h.ctx, job.ID)
	}
	report.JobStatus = job.Status
	report.DurationMS = time.Since(started).Milliseconds()
	if err == nil {
		var outBytes []byte
		outBytes, job, err = e.httpClient.dglOutput(h.ctx, job)
		if err == nil {
			report.OutputBytes = len(outBytes)
			report.OutputWidth, report.OutputHeight = dglImageDimensions(job.Result)
			return outBytes, report, e.settleGeneration(h, nil)
		}
	}
	// A terminal job state is the backend's own completion, even when it is a
	// failure. Only an unknown/transport outcome leaves the generation slot
	// abandoned and makes reclamation consider a recycle.
	switch job.Status {
	case "completed", "failed", "cancelled":
		return nil, report, e.settleGenerationKnown(h, err)
	default:
		return nil, report, e.settleGeneration(h, err)
	}
}

// buildDGLUpscaleScript is separate from the HTTP calls so its escaping can be
// tested without a server. DGL string literals are JSON-like but the parser
// rejected Windows backslash paths, so paths are emitted with forward slashes
// and quoted with strconv.Quote (which also escapes embedded quotes).
func buildDGLUpscaleScript(inputPath, upscalerPath string, tileSize int) (string, error) {
	if tileSize <= 0 {
		tileSize = 128
	}
	input, err := dglPathLiteral(inputPath)
	if err != nil {
		return "", fmt.Errorf("input path: %w", err)
	}
	upscaler, err := dglPathLiteral(upscalerPath)
	if err != nil {
		return "", fmt.Errorf("upscaler path: %w", err)
	}
	return fmt.Sprintf(`# version: 0.0.0

var {
    input_path = %s
    upscaler_path = %s
    tile_size = %d
}

fn main() {
    upscaler = load_upscaler(upscaler_path, tile_size = tile_size)
    image = upscale(upscaler, image::load(input_path))
    return { image = image }
}
`, input, upscaler, tileSize), nil
}

func dglUpscaleSource(inputPath, upscalerPath string, tileSize int) (string, error) {
	return buildDGLUpscaleScript(inputPath, upscalerPath, tileSize)
}

func dglPathLiteral(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("path is empty")
	}
	if strings.ContainsAny(path, "\"\r\n") {
		return "", fmt.Errorf("path contains a character DGL cannot represent: %q", path)
	}
	// DGL's parser rejected C:\... ("invalid scalar literal" on \U), while the
	// same path with forward slashes loaded correctly. Normalizing here keeps
	// the script portable and removes that class of failure.
	return strconv.Quote(filepath.ToSlash(path)), nil
}

func (c *client) dglSubmit(ctx context.Context, source string) (DGLJob, error) {
	payload, err := json.Marshal(map[string]string{"source": source})
	if err != nil {
		return DGLJob{}, notDispatched(fmt.Errorf("encode DGL job: %w", err))
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/sdcpp/v1/dgl/jobs", bytes.NewReader(payload))
	if err != nil {
		return DGLJob{}, notDispatched(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.do(req)
	if err != nil {
		return DGLJob{}, fmt.Errorf("submit DGL job: %w", err)
	}
	defer resp.Body.Close()
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return DGLJob{}, fmt.Errorf("read DGL submit response: %w", readErr)
	}
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return DGLJob{}, fmt.Errorf("submit DGL job: status %d: %s", resp.StatusCode, dglBodySnippet(body))
	}
	var job DGLJob
	if err := json.Unmarshal(body, &job); err != nil {
		return DGLJob{}, fmt.Errorf("decode DGL submit response: %w: %s", err, dglBodySnippet(body))
	}
	if job.ID == "" {
		return DGLJob{}, fmt.Errorf("DGL submit returned no job id: %s", dglBodySnippet(body))
	}
	return job, nil
}

func (c *client) dglWait(ctx context.Context, id string) (DGLJob, error) {
	if id == "" {
		return DGLJob{}, notDispatched(fmt.Errorf("DGL job id is empty"))
	}
	deadline := time.Now().Add(c.pollTimeout)
	for {
		if err := ctx.Err(); err != nil {
			return DGLJob{}, err
		}
		job, err := c.dglGet(ctx, id)
		if err != nil {
			return DGLJob{}, err
		}
		switch job.Status {
		case "completed":
			return job, nil
		case "failed":
			msg := "DGL job failed"
			if job.Error != nil && job.Error.Message != "" {
				msg = job.Error.Message
				if job.Error.Code != "" {
					msg = job.Error.Code + ": " + msg
				}
			}
			return job, fmt.Errorf("%s", msg)
		case "cancelled":
			return job, fmt.Errorf("DGL job was cancelled")
		case "queued", "generating", "executing":
			// Still running; keep polling.
		default:
			return job, fmt.Errorf("DGL job entered unknown status %q", job.Status)
		}
		if time.Now().After(deadline) {
			return job, fmt.Errorf("DGL job %s did not finish within %s", id, c.pollTimeout)
		}
		select {
		case <-ctx.Done():
			return DGLJob{}, ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}
	}
}

func (c *client) dglGet(ctx context.Context, id string) (DGLJob, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/sdcpp/v1/dgl/jobs/"+url.PathEscape(id), nil)
	if err != nil {
		return DGLJob{}, err
	}
	resp, err := c.httpShort.Do(req)
	if err != nil {
		return DGLJob{}, fmt.Errorf("poll DGL job: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return DGLJob{}, fmt.Errorf("read DGL job: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return DGLJob{}, fmt.Errorf("poll DGL job: status %d: %s", resp.StatusCode, dglBodySnippet(body))
	}
	var job DGLJob
	if err := json.Unmarshal(body, &job); err != nil {
		return DGLJob{}, fmt.Errorf("decode DGL job: %w: %s", err, dglBodySnippet(body))
	}
	return job, nil
}

// dglOutput selects the image output and fetches its PNG bytes. A completed
// DGL script may return other values; only the image named "image" is useful
// to the workflow, with index 0 as the fallback for a single unnamed output.
func (c *client) dglOutput(ctx context.Context, job DGLJob) ([]byte, DGLJob, error) {
	if job.Result == nil || len(job.Result.Outputs) == 0 {
		return nil, job, fmt.Errorf("DGL job completed without outputs")
	}
	index := -1
	for i, out := range job.Result.Outputs {
		if strings.EqualFold(out.Name, "image") {
			index = i
			break
		}
	}
	if index < 0 {
		index = 0
	}
	out := job.Result.Outputs[index]
	frame := 0
	reqURL := fmt.Sprintf("%s/sdcpp/v1/dgl/jobs/%s/outputs/%d?frame=%d",
		c.baseURL, url.PathEscape(job.ID), index, frame)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, job, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, job, fmt.Errorf("fetch DGL output: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, job, fmt.Errorf("read DGL output: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, job, fmt.Errorf("fetch DGL output: status %d: %s", resp.StatusCode, dglBodySnippet(data))
	}
	if len(data) < 8 || !bytes.HasPrefix(data, []byte("\x89PNG\r\n\x1a\n")) {
		return nil, job, fmt.Errorf("fetch DGL output: response is not a PNG (%d bytes, %d expected)", len(data), out.Bytes)
	}
	return data, job, nil
}

func dglImageDimensions(result *DGLJobResult) (int, int) {
	if result == nil {
		return 0, 0
	}
	for _, out := range result.Outputs {
		if !strings.EqualFold(out.Name, "image") {
			continue
		}
		return dglShapeWidthHeight(out.Shape)
	}
	if len(result.Outputs) > 0 {
		return dglShapeWidthHeight(result.Outputs[0].Shape)
	}
	return 0, 0
}

// dglShapeWidthHeight accepts both CHW ([3,H,W]) and HWC ([H,W,3]) shapes.
// The installed server returns CHW.
func dglShapeWidthHeight(shape []int) (int, int) {
	if len(shape) < 2 {
		return 0, 0
	}
	if len(shape) == 3 && shape[0] <= 4 {
		return shape[2], shape[1]
	}
	return shape[1], shape[0]
}

func dglBodySnippet(body []byte) string {
	s := strings.TrimSpace(string(body))
	if len(s) > 400 {
		s = s[:400] + "..."
	}
	return s
}
