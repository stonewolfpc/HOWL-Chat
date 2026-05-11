package whisper

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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

// DownloadLatest installs whisper-cli: official Windows ZIP from GitHub releases;
// UNIX-like systems build whisper-cli via git plus cmake plus a toolchain.
func (d *BinaryDownloader) DownloadLatest() (*WhisperBinaries, error) {
	return d.Acquire(context.Background())
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

Option 1: Automatic install (recommended)
------------------------------------------
Windows installs the official whisper-bin ZIP from ggml-org/whisper.cpp GitHub Releases.
Darwin/Linux clone the same release tag and build whisper-cli with git plus cmake plus clang/gcc.

Environment:
- HOWL_WHISPER_HOME  optional install root (defaults to ./bin/whisper under the working directory)

Option 2: Manual build from Source
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

	cmd := exec.Command(cliPath, "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("whisper-cli failed to execute: %w\nOutput: %s", err, string(output))
	}

	return nil
}
