package sdcpp

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"howl-chat/internal/image/types"
	"howl-chat/internal/paths"
	"howl-chat/internal/proc"
)

// videoFamilies lists the architectures sd-cli can render video from in
// -M vid_gen mode. Anything else is an image model.
var videoFamilies = map[string]bool{
	"wan": true, "wan2.1": true, "wan2.2": true, "wan21": true, "wan22": true,
	"ltx": true, "ltx2": true, "ltx-video": true,
	"hunyuan": true, "hunyuanvideo": true, "hunyuan-video": true,
	"animatediff": true,
	"minimax":     true, "minimax-h3": true,
}

// videoDeps lists companion-file families that need a text encoder + VAE.
var videoDeps = map[string]bool{
	"wan": true, "wan2.1": true, "wan2.2": true, "wan21": true, "wan22": true,
	"ltx": true, "ltx2": true, "ltx-video": true,
	"hunyuan": true, "hunyuanvideo": true, "hunyuan-video": true,
	"minimax": true, "minimax-h3": true,
}

// IsVideoFamily reports whether a model family is a video architecture.
func IsVideoFamily(family string) bool {
	return videoFamilies[strings.ToLower(strings.TrimSpace(family))]
}

// videoDefaults returns frame-count and fps defaults tuned per architecture.
func videoDefaults(family string) (frames, fps int) {
	switch strings.ToLower(strings.TrimSpace(family)) {
	case "ltx", "ltx2", "ltx-video":
		return 97, 24
	case "hunyuan", "hunyuanvideo", "hunyuan-video":
		return 45, 24
	default: // Wan and AnimateDiff
		return 33, 16
	}
}

// videoContainer maps a requested output format onto a container sd-cli can
// write by itself. MP4/GIF need ffmpeg, which is not bundled.
func videoContainer(format string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "", "webm":
		return "webm", nil
	case "avi":
		return "avi", nil
	case "webp":
		return "webp", nil
	case "mp4", "gif":
		return "", fmt.Errorf("%s output needs ffmpeg, which is not bundled yet — choose WebM, AVI or animated WebP", strings.ToUpper(format))
	default:
		return "", fmt.Errorf("unsupported video format %q (use webm, avi or webp)", format)
	}
}

// stepRe matches sd.cpp's sampling progress lines, e.g. "  4/20 - 1.20s/it".
var stepRe = regexp.MustCompile(`(\d+)\s*/\s*(\d+)\s*-\s*[\d.]+s/it`)

// FindSDCLI locates sd-cli.exe, which is the only sd.cpp binary that can run
// -M vid_gen (sd-server exposes the image API only).
func FindSDCLI() string {
	candidates := []string{
		filepath.Join(filepath.Dir(os.Args[0]), "runtime", "stablediffusion", "sd-cli.exe"),
		filepath.Join(paths.ProjectRoot(), "sd-cli.exe"),
		filepath.Join(filepath.Dir(os.Args[0]), "sd-cli.exe"),
		"sd-cli.exe",
		filepath.Join(paths.ProjectRoot(), "sd.exe"),
		filepath.Join(filepath.Dir(os.Args[0]), "sd.exe"),
		"sd.exe",
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return ""
}

// findSDCLIFor resolves sd-cli next to the configured sd-server first, so a
// bundled runtime layout keeps working from any launch directory.
func (e *Engine) findSDCLIFor() string {
	if e.config.BinaryPath != "" {
		dir := filepath.Dir(e.config.BinaryPath)
		base := filepath.Base(e.config.BinaryPath)
		cliName := "sd-cli.exe"
		if strings.EqualFold(filepath.Ext(base), "") {
			cliName = "sd-cli"
		}
		near := filepath.Join(dir, cliName)
		if _, err := os.Stat(near); err == nil {
			return near
		}
	}
	return FindSDCLI()
}

// stopServer terminates the HTTP server so the CLI run gets the whole GPU.
// The next image request relaunches it lazily through EnsureReady.
func (e *Engine) stopServer() {
	if e.serverProc == nil {
		return
	}
	proc.KillTree(e.serverProc.Pid)
	_ = e.serverProc.Kill()
	_, _ = e.serverProc.Wait()
	e.serverProc = nil
}

// SelectModel resolves a model without starting sd-server. Video generation
// uses it so the only process that touches the GPU is sd-cli itself.
func (e *Engine) SelectModel(opts types.ModelLoadOptions) error {
	target, err := e.resolveTarget(opts)
	if err != nil {
		return err
	}
	e.currentModel = target
	return nil
}

// generateVideo renders a clip with sd-cli and returns the written file.
func (e *Engine) generateVideo(opts types.ImageOptions) (*types.ImageResult, error) {
	if e.currentModel == nil {
		return nil, fmt.Errorf("no model loaded")
	}
	family := strings.ToLower(e.currentModel.Family)
	if !IsVideoFamily(family) {
		return nil, fmt.Errorf("model %q (%s) cannot generate video — load a Wan, LTX, HunyuanVideo or AnimateDiff model", e.currentModel.Name, e.currentModel.Family)
	}
	cli := e.findSDCLIFor()
	if cli == "" {
		return nil, fmt.Errorf("sd-cli.exe not found — video generation needs sd-cli next to sd-server")
	}
	ext, err := videoContainer(opts.OutputFormat)
	if err != nil {
		return nil, err
	}
	if !e.runMu.TryLock() {
		return nil, fmt.Errorf("a video generation is already running")
	}
	defer e.runMu.Unlock()

	frames, fps := videoDefaults(family)
	if opts.VideoFrames > 0 {
		frames = opts.VideoFrames
	}
	if opts.FPS > 0 {
		fps = opts.FPS
	}
	width, height := videoResolution(e.currentModel, family)
	if opts.Width > 0 {
		width = opts.Width
	}
	if opts.Height > 0 {
		height = opts.Height
	}
	if opts.Steps <= 0 {
		opts.Steps = 20
	}

	outPath := filepath.Join(e.config.OutputDir, fmt.Sprintf("video_%d.%s", nowUnixNano(), ext))
	args := e.buildVideoArgs(opts, frames, fps, width, height, outPath)

	// Free the GPU first: sd-server would otherwise hold the model in VRAM.
	e.stopServer()

	fmt.Printf("[ImageEngine] sd-cli %s\n", strings.Join(args, " "))
	err = e.runVideoCLI(cli, args, outPath, opts.JobID)
	if e.consumeCancel() {
		return nil, fmt.Errorf("video generation cancelled")
	}
	if err != nil {
		return nil, err
	}

	info, statErr := os.Stat(outPath)
	if statErr != nil || info.Size() == 0 {
		_ = os.Remove(outPath)
		return nil, fmt.Errorf("sd-cli finished without writing %s", filepath.Base(outPath))
	}

	videoB64, encErr := readAsBase64(outPath)
	if encErr != nil {
		return nil, fmt.Errorf("read generated video: %w", encErr)
	}

	md := types.ImageMetadata{
		File:           outPath,
		Prompt:         opts.Prompt,
		NegativePrompt: opts.NegativePrompt,
		Model:          e.currentModel.Name,
		Family:         e.currentModel.Family,
		Width:          width,
		Height:         height,
		Steps:          opts.Steps,
		Seed:           opts.Seed,
		Sampler:        opts.Sampler,
		CfgScale:       opts.CfgScale,
		Mode:           "video",
		CreatedAt:      createdStamp(),
		Options:        opts,
	}
	return &types.ImageResult{
		Path:      outPath,
		VideoData: videoB64,
		VideoMIME: videoMIME(ext),
		Metadata:  md,
	}, nil
}

// buildVideoArgs assembles the sd-cli command line for a video run.
func (e *Engine) buildVideoArgs(opts types.ImageOptions, frames, fps, width, height int, outPath string) []string {
	m := e.currentModel
	family := strings.ToLower(m.Family)

	// Wan/LTX/Hunyuan/MiniMax ship as separate diffusion-model + encoders + VAE.
	split := videoDeps[family]
	args := []string{"-M", "vid_gen"}
	if p := m.Files["model"]; p != "" {
		if split {
			args = append(args, "--diffusion-model", p)
		} else {
			args = append(args, "-m", p)
		}
	}
	if p := m.Files["vae"]; p != "" {
		args = append(args, "--vae", p)
	}
	if p := m.Files["t5xxl"]; p != "" {
		args = append(args, "--t5xxl", p)
	}
	if p := m.Files["clip_l"]; p != "" {
		args = append(args, "--clip_l", p)
	}
	if p := m.Files["clip_g"]; p != "" {
		args = append(args, "--clip_g", p)
	}
	if p := m.Files["clip_vision"]; p != "" {
		args = append(args, "--clip_vision", p)
	}
	if p := m.Files["llm"]; p != "" {
		args = append(args, "--llm", p)
	}
	if p := m.Files["motion_module"]; p != "" {
		args = append(args, "--motion-module", p)
	}

	args = append(args, "--prompt", opts.Prompt)
	if opts.NegativePrompt != "" {
		args = append(args, "--negative-prompt", opts.NegativePrompt)
	}
	args = append(args,
		"-W", strconv.Itoa(width),
		"-H", strconv.Itoa(height),
		"--steps", strconv.Itoa(opts.Steps),
		"--seed", strconv.FormatInt(opts.Seed, 10),
		"--video-frames", strconv.Itoa(frames),
		"--fps", strconv.Itoa(fps),
		"-o", outPath,
	)
	if opts.Sampler != "" {
		args = append(args, "--sampling-method", opts.Sampler)
	}
	if opts.Scheduler != "" {
		args = append(args, "--scheduler", opts.Scheduler)
	}
	if opts.CfgScale > 0 {
		args = append(args, "--cfg-scale", trimFloat(opts.CfgScale))
	}
	if opts.Guidance > 0 {
		args = append(args, "--guidance", trimFloat(opts.Guidance))
	}
	if opts.FlowShift > 0 {
		args = append(args, "--flow-shift", trimFloat(opts.FlowShift))
	}
	if opts.InitImagePath != "" {
		args = append(args, "-i", opts.InitImagePath)
	}
	if opts.EndImagePath != "" {
		args = append(args, "--end-img", opts.EndImagePath)
	}
	if opts.MaskPath != "" {
		args = append(args, "--mask", opts.MaskPath)
	}
	// Always tile the VAE for video: the identical 9-frame 384x256 decode took
	// 112 s untiled (the VAE spilled to CPU) and 9.4 s tiled, with sampling
	// unchanged. It is the single biggest speed win for video on 8 GB cards.
	args = append(args, "--vae-tiling")
	if opts.TemporalTiling {
		args = append(args, "--temporal-tiling")
	}
	if opts.FlashAttention {
		args = append(args, "--fa")
	}
	if opts.OffloadToCPU {
		args = append(args, "--offload-to-cpu")
	}
	if opts.CacheMode != "" {
		args = append(args, "--cache-mode", opts.CacheMode)
		if opts.CacheOption != "" {
			args = append(args, "--cache-option", opts.CacheOption)
		}
	}
	if opts.DisableMetadata {
		args = append(args, "--disable-image-metadata")
	}
	return args
}

// runVideoCLI executes sd-cli, streaming progress and keeping the tail of the
// log so a failure can be explained to the user.
func (e *Engine) runVideoCLI(cli string, args []string, outPath string, owner string) error {
	cmd := exec.Command(cli, args...)
	cmd.Dir = filepath.Dir(cli)
	proc.HideWindow(cmd)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("capture sd-cli output: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start sd-cli: %w", err)
	}
	proc.Attach(cmd)

	e.videoMu.Lock()
	e.videoProc = cmd.Process
	e.videoOwner = owner
	e.cancelled = false
	e.progress = types.ProgressInfo{}
	e.videoMu.Unlock()

	var tail []string
	scanDone := make(chan struct{})
	go func() {
		defer close(scanDone)
		sc := bufio.NewScanner(stderr)
		sc.Buffer(make([]byte, 32*1024), 512*1024)
		sc.Split(scanAnyEOL)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" {
				continue
			}
			e.recordVideoProgress(line)
			tail = append(tail, line)
			if len(tail) > 8 {
				tail = tail[1:]
			}
		}
	}()

	waitErr := cmd.Wait()
	<-scanDone

	e.videoMu.Lock()
	e.videoProc = nil
	e.videoOwner = ""
	e.videoMu.Unlock()

	if waitErr != nil {
		if f, statErr := os.Stat(outPath); statErr == nil && f.Size() > 0 {
			return nil // wrote the file, then exited non-zero: keep the result
		}
		detail := "no output captured"
		if len(tail) > 0 {
			detail = strings.Join(tail, " | ")
		}
		return fmt.Errorf("sd-cli failed: %v (%s)", waitErr, detail)
	}
	return nil
}

// recordVideoProgress parses a sampling line from sd.cpp's log.
func (e *Engine) recordVideoProgress(line string) {
	m := stepRe.FindStringSubmatch(line)
	if m == nil {
		return
	}
	step, errStep := strconv.Atoi(m[1])
	total, errTotal := strconv.Atoi(m[2])
	if errStep != nil || errTotal != nil || total <= 0 || step > total {
		return
	}
	e.videoMu.Lock()
	e.progress = types.ProgressInfo{
		Progress:   float64(step) / float64(total),
		Step:       step,
		TotalSteps: total,
	}
	e.videoMu.Unlock()
}

// cancelVideoOwned kills the sd-cli render only when jobID owns it, so one job's
// cancel cannot stop another job's video. This is cancellation identity only:
// the process handling below is unchanged.
func (e *Engine) cancelVideoOwned(jobID string) bool {
	if jobID == "" {
		return false
	}
	e.videoMu.Lock()
	procRef := e.videoProc
	owned := procRef != nil && e.videoOwner == jobID
	if owned {
		e.cancelled = true
	}
	e.videoMu.Unlock()
	if !owned {
		return false
	}
	proc.KillTree(procRef.Pid)
	_ = procRef.Kill()
	return true
}

// cancelVideo kills a running sd-cli render. Returns true when one was active.
func (e *Engine) cancelVideo() bool {
	e.videoMu.Lock()
	procRef := e.videoProc
	if procRef != nil {
		e.cancelled = true
	}
	e.videoMu.Unlock()

	if procRef == nil {
		return false
	}
	proc.KillTree(procRef.Pid)
	_ = procRef.Kill()
	return true
}

// consumeCancel reports and clears the cancellation flag.
func (e *Engine) consumeCancel() bool {
	e.videoMu.Lock()
	defer e.videoMu.Unlock()
	was := e.cancelled
	e.cancelled = false
	return was
}

// videoProgressInfo returns the live CLI progress while a render is running.
func (e *Engine) videoProgressInfo() (types.ProgressInfo, bool) {
	e.videoMu.Lock()
	defer e.videoMu.Unlock()
	if e.videoProc == nil {
		return types.ProgressInfo{}, false
	}
	return e.progress, true
}

// videoResolution picks the native frame size for a video architecture,
// falling back to the model's resolution hint when it is usable.
func videoResolution(m *types.ModelInfo, family string) (int, int) {
	if w, h, ok := parseResolution(m.ResolutionHint); ok {
		return w, h
	}
	switch family {
	case "ltx", "ltx2", "ltx-video":
		return 768, 512
	case "hunyuan", "hunyuanvideo", "hunyuan-video":
		return 848, 480
	case "minimax", "minimax-h3":
		return 768, 432
	default: // Wan 1.3B/5B and AnimateDiff
		return 832, 480
	}
}

// parseResolution reads a "WxH" hint.
func parseResolution(hint string) (int, int, bool) {
	parts := strings.Split(strings.ToLower(strings.TrimSpace(hint)), "x")
	if len(parts) != 2 {
		return 0, 0, false
	}
	w, errW := strconv.Atoi(parts[0])
	h, errH := strconv.Atoi(parts[1])
	if errW != nil || errH != nil || w < 64 || h < 64 {
		return 0, 0, false
	}
	return w, h, true
}

// scanAnyEOL splits on CR or LF: sd.cpp rewrites its progress line with \r.
func scanAnyEOL(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexAny(data, "\r\n"); i >= 0 {
		return i + 1, data[:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

// trimFloat formats a float without a trailing ".0".
func trimFloat(v float64) string {
	return strconv.FormatFloat(v, 'g', -1, 64)
}

// nowUnixNano and createdStamp keep filenames sortable and metadata RFC3339.
func nowUnixNano() int64   { return time.Now().UnixNano() }
func createdStamp() string { return time.Now().Format("2006-01-02T15:04:05Z07:00") }

// videoMIME maps a container extension to the MIME type used in a data URL.
func videoMIME(ext string) string {
	switch strings.ToLower(ext) {
	case "avi":
		return "video/x-msvideo"
	case "webp":
		return "image/webp"
	default:
		return "video/webm"
	}
}

// readAsBase64 returns file contents encoded for a data: URL.
func readAsBase64(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}
