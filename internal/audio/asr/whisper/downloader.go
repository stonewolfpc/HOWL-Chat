package whisper

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// WhisperBinaries contains paths to downloaded executables
type WhisperBinaries struct {
	CliPath  string
	DllPaths []string
	ModelDir string
}

// BinaryDownloader handles downloading whisper.cpp binaries
type BinaryDownloader struct {
	targetDir string
}

// NewBinaryDownloader creates a downloader for the specified directory
func NewBinaryDownloader(targetDir string) *BinaryDownloader {
	return &BinaryDownloader{
		targetDir: targetDir,
	}
}

// DownloadLatest downloads the latest prebuilt binaries for Windows
// Note: For production use, you should bundle these with your installer
func (d *BinaryDownloader) DownloadLatest() (*WhisperBinaries, error) {
	// Create target directory
	if err := os.MkdirAll(d.targetDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create target directory: %w", err)
	}

	// For Windows, download from SourceForge mirror
	if runtime.GOOS == "windows" {
		return d.downloadWindowsBinaries()
	}

	// For Linux/Mac, user needs to build from source
	return nil, fmt.Errorf("automatic download only supported on Windows. Please build from source on %s", runtime.GOOS)
}

// downloadWindowsBinaries downloads Windows prebuilt binaries
func (d *BinaryDownloader) downloadWindowsBinaries() (*WhisperBinaries, error) {
	// SourceForge mirror download URL
	downloadURL := "https://sourceforge.net/projects/whisper-cpp.mirror/files/latest/download"

	zipPath := filepath.Join(d.targetDir, "whisper-temp.zip")

	// Download ZIP
	fmt.Printf("Downloading whisper.cpp binaries from %s...\n", downloadURL)
	if err := downloadFile(downloadURL, zipPath); err != nil {
		return nil, fmt.Errorf("failed to download: %w", err)
	}
	defer os.Remove(zipPath)

	// Extract ZIP
	if err := extractZip(zipPath, d.targetDir); err != nil {
		return nil, fmt.Errorf("failed to extract: %w", err)
	}

	// Find the executables
	binaries := &WhisperBinaries{
		ModelDir: filepath.Join(d.targetDir, "models"),
	}

	// Search for whisper-cli.exe
	cliPath := filepath.Join(d.targetDir, "whisper-cli.exe")
	if _, err := os.Stat(cliPath); os.IsNotExist(err) {
		// Try in subdirectories
		err := filepath.Walk(d.targetDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() == "whisper-cli.exe" {
				binaries.CliPath = path
			}
			if filepath.Ext(path) == ".dll" {
				binaries.DllPaths = append(binaries.DllPaths, path)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to find whisper-cli.exe: %w", err)
		}
	} else {
		binaries.CliPath = cliPath
	}

	if binaries.CliPath == "" {
		return nil, fmt.Errorf("whisper-cli.exe not found in downloaded archive")
	}

	// Create models directory
	os.MkdirAll(binaries.ModelDir, 0755)

	return binaries, nil
}

// downloadFile downloads a file from URL to local path
func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// extractZip extracts a ZIP file to the target directory
func extractZip(zipPath, targetDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(targetDir, file.Name)

		// Security: Prevent zip slip
		if !strings.HasPrefix(path, filepath.Clean(targetDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path in zip: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

// GetInstallationInstructions returns manual install steps
func GetInstallationInstructions() string {
	return `
WHISPER.CPP INSTALLATION GUIDE
==============================

Option 1: Prebuilt Binaries (Windows)
-------------------------------------
1. Download from: https://sourceforge.net/projects/whisper-cpp.mirror/files/latest/download
2. Extract to: C:\Program Files\whisper\ or your app directory
3. Add to PATH or place in: bin/ subdirectory of your app

Option 2: Build from Source
---------------------------
Requirements:
- CMake 3.18+
- Visual Studio 2022 (Windows) or GCC (Linux/Mac)
- Git

Build Steps:
  git clone https://github.com/ggml-org/whisper.cpp.git
  cd whisper.cpp
  cmake -B build -DWHISPER_BUILD_EXAMPLES=ON
  cmake --build build --config Release

Result:
- build/bin/Release/whisper-cli.exe (Windows)
- build/bin/whisper-cli (Linux/Mac)

Option 3: Package with Your App
-------------------------------
For distribution, bundle whisper-cli.exe with your installer:
- Place whisper-cli.exe in: bin/
- Place ggml.dll in: bin/ (same directory)
- Place models in: models/

Download Models:
----------------
Models are separate from the binary. Download GGUF models from:
- https://huggingface.co/ggerganov/whisper.cpp

Recommended starter models:
- ggml-base.bin (74MB) - Good balance
- ggml-small.bin (244MB) - Better quality

Place models in: models/ subdirectory

Dependencies:
-------------
Windows: No additional dependencies (single .exe + .dlls included)
Linux:   May need: sudo apt install libgomp1
Mac:     No additional dependencies

Verification:
-------------
Test installation:
  whisper-cli.exe --help
  
Test transcription:
  whisper-cli.exe --model models/ggml-base.bin --file test.wav
`
}

// CheckInstallation verifies whisper-cli is properly installed
func CheckInstallation(cliPath string) error {
	if cliPath == "" {
		return fmt.Errorf("whisper-cli path not specified")
	}

	if _, err := os.Stat(cliPath); os.IsNotExist(err) {
		return fmt.Errorf("whisper-cli not found at: %s", cliPath)
	}

	// Try to run --help to verify it works
	cmd := execCommand(cliPath, "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("whisper-cli failed to execute: %w\nOutput: %s", err, string(output))
	}

	return nil
}
