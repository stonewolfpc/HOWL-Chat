package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"howl-chat/internal/image/registry"
	"howl-chat/internal/paths"
)

var sdServerPath string

// resolveSDServerPath finds the sd-server binary the same way the app does:
// project root first, then next to the binary, then PATH. An explicit -sd flag
// always wins so this harness works against a build kept anywhere.
func resolveSDServerPath(explicit string) string {
	if explicit != "" {
		return explicit
	}
	candidates := []string{filepath.Join(paths.ProjectRoot(), "sd-server.exe")}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "sd-server.exe"))
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	if p, err := exec.LookPath("sd-server.exe"); err == nil {
		return p
	}
	return ""
}

// resolveScanDirs defaults to the models directory the app is configured to
// use, so the harness reports on the same folder the Image Generator scans.
func resolveScanDirs(explicit string) []string {
	if strings.TrimSpace(explicit) != "" {
		var dirs []string
		for _, p := range strings.Split(explicit, ",") {
			if p = strings.TrimSpace(p); p != "" {
				dirs = append(dirs, p)
			}
		}
		return dirs
	}

	var dirs []string
	seen := make(map[string]bool)
	add := func(dir string) {
		dir = strings.TrimSpace(dir)
		if dir == "" {
			return
		}
		key := strings.ToLower(filepath.Clean(dir))
		if seen[key] {
			return
		}
		seen[key] = true
		dirs = append(dirs, dir)
	}

	configDir := filepath.Join(paths.ProjectRoot(), "config")

	// Follow the same directory ledger the app uses, so this report covers
	// exactly the folders the Image Generator scans.
	if reg, err := registry.New(configDir); err == nil {
		for _, dir := range reg.Directories(registry.KindImage) {
			add(dir)
		}
	}

	if len(dirs) == 0 {
		if data, err := os.ReadFile(filepath.Join(configDir, "image-settings.json")); err == nil {
			var cfg struct {
				Engine struct {
					ModelsDir string `json:"models_dir"`
				} `json:"engine_config"`
			}
			if err := json.Unmarshal(data, &cfg); err == nil {
				add(cfg.Engine.ModelsDir)
			}
		}
	}
	if len(dirs) == 0 {
		add(filepath.Join(paths.ProjectRoot(), "models"))
	}
	return dirs
}

// classify determines if a file is a full model, LoRA, VAE, or embedding by size and name
func classify(name string, size int64) string {
	lower := strings.ToLower(name)
	mb := float64(size) / (1024 * 1024)

	// Name-based detection (most reliable)
	if strings.Contains(lower, "lora") {
		return "LoRA"
	}
	if strings.Contains(lower, "vae") {
		return "VAE"
	}
	if strings.Contains(lower, "embedding") || strings.Contains(lower, "neg") && mb < 2 {
		return "embedding"
	}
	if strings.Contains(lower, "control") || strings.Contains(lower, "controlnet") {
		return "ControlNet"
	}
	if strings.Contains(lower, "upscaler") || strings.Contains(lower, "esrgan") {
		return "upscaler"
	}

	// Size-based detection
	if mb < 50 {
		return "embedding"
	}
	if mb < 300 {
		return "LoRA"
	}
	if mb < 800 {
		return "VAE_or_LoRA"
	}
	if mb >= 1500 {
		return "model"
	}
	return "unknown"
}

// testModel tries to start sd-server with the given model and checks if API responds
func testModel(path string) (bool, string) {
	// Kill any existing sd-server
	exec.Command("taskkill", "/f", "/im", "sd-server.exe").Run()
	time.Sleep(500 * time.Millisecond)

	cmd := exec.Command(sdServerPath,
		"--listen-port", "8082", // use different port to avoid conflicts
		"-m", path,
	)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard

	if err := cmd.Start(); err != nil {
		return false, fmt.Sprintf("failed to start: %v", err)
	}

	// Wait for model to load and API to respond (up to 90s)
	client := &http.Client{Timeout: 5 * time.Second}
	url := "http://127.0.0.1:8082/sdapi/v1/samplers"
	deadline := time.Now().Add(90 * time.Second)
	var lastErr string

	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			cmd.Process.Kill()
			cmd.Wait()
			return true, "API responded OK"
		}
		lastErr = err.Error()
		time.Sleep(2 * time.Second)
	}

	// Capture stderr for diagnosis
	cmd.Process.Kill()
	cmd.Wait()

	return false, fmt.Sprintf("timed out (last: %s)", lastErr)
}

func main() {
	sdFlag := flag.String("sd", "", "path to sd-server executable (default: auto-discover from the project root)")
	dirsFlag := flag.String("dirs", "", "comma-separated model directories to scan (default: every directory added to HOWL as an image directory)")
	flag.Parse()

	sdServerPath = resolveSDServerPath(*sdFlag)
	dirs := resolveScanDirs(*dirsFlag)
	if sdServerPath == "" {
		fmt.Println("ERROR: sd-server not found. Pass -sd <path> or place sd-server.exe in the project root.")
	} else {
		fmt.Printf("Using sd-server: %s\n", sdServerPath)
	}

	fmt.Println("╔══════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    HOWL Image Model Diagnostic Report                       ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	var modelsToTest []struct {
		path  string
		name  string
		size  int64
		class string
		dir   string
	}

	// Phase 1: Scan and classify
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			fmt.Printf("ERROR reading %s: %v\n", dir, err)
			continue
		}

		fmt.Printf("\n📁 %s\n", dir)
		fmt.Println(strings.Repeat("─", 80))

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			info, _ := entry.Info()
			name := entry.Name()
			path := filepath.Join(dir, name)
			class := classify(name, info.Size())
			mb := float64(info.Size()) / (1024 * 1024)

			modelsToTest = append(modelsToTest, struct {
				path  string
				name  string
				size  int64
				class string
				dir   string
			}{path, name, info.Size(), class, dir})

			icon := "🎨"
			switch class {
			case "LoRA":
				icon = "🔧"
			case "VAE", "VAE_or_LoRA":
				icon = "🎭"
			case "embedding":
				icon = "🏷️"
			case "ControlNet":
				icon = "🎮"
			case "upscaler":
				icon = "🔍"
			}
			fmt.Printf("  %s %-50s %8.1f MB  [%s]\n", icon, name, mb, class)
		}
	}

	// Phase 2: Test each classified "model" (skip LoRAs, VAEs, embeddings for now)
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════════════════════")
	fmt.Println("                           MODEL LOAD TESTS")
	fmt.Println("═══════════════════════════════════════════════════════════════════════════════")

	for _, m := range modelsToTest {
		if m.class != "model" && m.class != "VAE_or_LoRA" {
			continue // skip LoRAs, embeddings, etc.
		}

		mb := float64(m.size) / (1024 * 1024)
		fmt.Printf("\n🧪 Testing: %-45s (%6.1f MB) ", m.name, mb)

		// For GGUF files, we already know they don't work
		if strings.HasSuffix(strings.ToLower(m.name), ".gguf") {
			fmt.Println("❌ SKIP (GGUF - incompatible with stable-diffusion.cpp)")
			continue
		}

		// For very small files, skip
		if mb < 800 {
			fmt.Println("⚠️  SKIP (too small to be a full model)")
			continue
		}

		ok, msg := testModel(m.path)
		if ok {
			fmt.Printf("✅ WORKS - %s\n", msg)
		} else {
			fmt.Printf("❌ FAIL - %s\n", msg)
		}
	}

	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════════════════════")
	fmt.Println("                              SUMMARY")
	fmt.Println("═══════════════════════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("CATEGORY BREAKDOWN:")

	counts := make(map[string]int)
	for _, m := range modelsToTest {
		counts[m.class]++
	}
	for class, count := range counts {
		fmt.Printf("  %-15s: %d file(s)\n", class, count)
	}

	fmt.Println()
	fmt.Println("RECOMMENDATIONS:")
	fmt.Println("  • Full models (.safetensors, .ckpt) ≥1.5GB → should work with sd-server")
	fmt.Println("  • GGUF files → INCOMPATIBLE with stable-diffusion.cpp (use .safetensors instead)")
	fmt.Println("  • LoRAs (50-300MB) → Keep but don't load as main model")
	fmt.Println("  • VAEs (100-700MB) → Keep but don't load as main model")
	fmt.Println("  • Embeddings (<10MB) → Keep but don't load as main model")
	fmt.Println()
	fmt.Println("SAFE TO DELETE (confirmed incompatible):")
	for _, m := range modelsToTest {
		if strings.HasSuffix(strings.ToLower(m.name), ".gguf") && m.class == "model" {
			fmt.Printf("  • %s\n", m.path)
		}
	}
}
