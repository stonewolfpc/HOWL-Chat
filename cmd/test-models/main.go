package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var sdServerPath = `D:\StableDiffusion\sd-server.exe`

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
	dirs := []string{`D:\AI GGUF\Image`, `D:\AI GGUF\Image Models`}

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
