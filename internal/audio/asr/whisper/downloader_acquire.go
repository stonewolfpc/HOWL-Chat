package whisper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const whisperGitHubReleaseAPI = "https://api.github.com/repos/ggml-org/whisper.cpp/releases/latest"

type githubReleasePayload struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// Acquire obtains whisper-cli: Windows ZIP from ggml-org releases;
// Darwin/Linux build from tagged source via git + cmake when tools exist.
func (d *BinaryDownloader) Acquire(ctx context.Context) (*WhisperBinaries, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := os.MkdirAll(d.targetDir, 0o755); err != nil {
		return nil, fmt.Errorf("create whisper install dir: %w", err)
	}

	switch runtime.GOOS {
	case "windows":
		return d.acquireWindowsZIP(ctx)
	default:
		return d.acquireUnixBuild(ctx)
	}
}

func (d *BinaryDownloader) acquireWindowsZIP(ctx context.Context) (*WhisperBinaries, error) {
	assetURL, zipName, err := pickWindowsZIPAsset(ctx)
	if err != nil {
		return nil, err
	}
	zipPath := filepath.Join(d.targetDir, zipName)
	if err := downloadHTTPFile(ctx, assetURL, zipPath); err != nil {
		return nil, fmt.Errorf("download whisper release zip: %w", err)
	}
	defer func() { _ = os.Remove(zipPath) }()

	if err := extractZip(zipPath, d.targetDir); err != nil {
		return nil, fmt.Errorf("extract whisper zip: %w", err)
	}

	modelDir := filepath.Join(d.targetDir, "models")
	_ = os.MkdirAll(modelDir, 0o755)

	cliPath := filepath.Join(d.targetDir, "whisper-cli.exe")
	if _, err := os.Stat(cliPath); os.IsNotExist(err) {
		var found string
		_ = filepath.Walk(d.targetDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info == nil {
				return nil
			}
			if filepath.Base(path) == "whisper-cli.exe" {
				found = path
				return filepath.SkipAll
			}
			return nil
		})
		if found == "" {
			return nil, fmt.Errorf("whisper-cli.exe not found after extracting %s", zipName)
		}
		cliPath = found
	}

	var dllPaths []string
	_ = filepath.Walk(filepath.Dir(cliPath), func(path string, info os.FileInfo, err error) error {
		if err == nil && info != nil && strings.EqualFold(filepath.Ext(path), ".dll") {
			dllPaths = append(dllPaths, path)
		}
		return nil
	})

	return &WhisperBinaries{CliPath: cliPath, ModelDir: modelDir, DllPaths: dllPaths}, nil
}

func pickWindowsZIPAsset(ctx context.Context) (url string, filename string, err error) {
	name := pickWindowsZIPName(runtime.GOARCH)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, whisperGitHubReleaseAPI, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "HOWL_Chat-whisper-fetch")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", "", fmt.Errorf("GitHub whisper release API HTTP %s: %s", resp.Status, strings.TrimSpace(string(b)))
	}
	var payload githubReleasePayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", "", fmt.Errorf("parse GitHub whisper release payload: %w", err)
	}
	for _, a := range payload.Assets {
		if strings.EqualFold(a.Name, name) {
			return a.BrowserDownloadURL, a.Name, nil
		}
	}
	return "", "", fmt.Errorf("no whisper release asset named %s in latest release", name)
}

func pickWindowsZIPName(goarch string) string {
	switch goarch {
	case "386":
		return "whisper-bin-Win32.zip"
	default:
		return "whisper-bin-x64.zip"
	}
}

func (d *BinaryDownloader) acquireUnixBuild(ctx context.Context) (*WhisperBinaries, error) {
	if _, err := exec.LookPath("git"); err != nil {
		return nil, fmt.Errorf("git is required to build whisper.cpp on %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	if _, err := exec.LookPath("cmake"); err != nil {
		return nil, fmt.Errorf("cmake is required to build whisper.cpp on %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	cc := ""
	if runtime.GOOS == "darwin" {
		cc = pickDarwinCompiler()
		if cc == "" {
			return nil, fmt.Errorf("clang or gcc must be installed to build whisper.cpp on Darwin")
		}
	} else if _, err := exec.LookPath("gcc"); err == nil {
		cc = "gcc"
	} else if _, err := exec.LookPath("clang"); err == nil {
		cc = "clang"
	} else {
		return nil, fmt.Errorf("a C compiler (gcc or clang) is required to build whisper.cpp")
	}

	tag, err := latestWhisperReleaseTag(ctx)
	if err != nil {
		return nil, err
	}

	buildRoot := filepath.Join(d.targetDir, "src")
	if err := os.RemoveAll(buildRoot); err != nil {
		return nil, fmt.Errorf("clear build root: %w", err)
	}
	if err := os.MkdirAll(buildRoot, 0o755); err != nil {
		return nil, err
	}
	repoDir := filepath.Join(buildRoot, "whisper.cpp")

	run := func(name string, args ...string) error {
		cmd := exec.CommandContext(ctx, name, args...)
		cmd.Dir = filepath.Dir(repoDir)
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		env := append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
		if cc != "" {
			env = append(env, "CC="+cc, "CXX="+pickCXX(cc))
		}
		cmd.Env = env
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s %v failed: %w", name, args, err)
		}
		return nil
	}

	if err := run("git", "clone", "--depth", "1", "--branch", tag, "--single-branch",
		"https://github.com/ggml-org/whisper.cpp.git", filepath.Base(repoDir)); err != nil {
		return nil, fmt.Errorf("git clone whisper.cpp: %w", err)
	}

	buildDir := filepath.Join(repoDir, "build_rel")
	cmakeExe := cmakeExecutable()
	configureArgs := []string{"-S", repoDir, "-B", buildDir, "-DCMAKE_BUILD_TYPE=Release", "-DWHISPER_BUILD_EXAMPLES=ON"}
	ccEnv := []string{"GIT_TERMINAL_PROMPT=0"}
	if cc != "" {
		ccEnv = append(ccEnv, "CC="+cc, "CXX="+pickCXX(cc))
	}
	wd := filepath.Dir(repoDir)
	if err := execCMake(ctx, wd, cmakeExe, configureArgs, ccEnv); err != nil {
		return nil, fmt.Errorf("cmake configure whisper.cpp failed: %w", err)
	}

	buildArgs := []string{"--build", buildDir, "--config", "Release", "--target", "whisper-cli", "--parallel"}
	if err := execCMake(ctx, wd, cmakeExe, buildArgs, ccEnv); err != nil {
		return nil, fmt.Errorf("cmake build whisper-cli failed: %w", err)
	}

	cliBin, ferr := locateWhisperCLIUnix(buildRoot)
	if ferr != nil {
		return nil, ferr
	}

	modelDir := filepath.Join(d.targetDir, "models")
	_ = os.MkdirAll(modelDir, 0o755)

	return &WhisperBinaries{
		CliPath:  cliBin,
		ModelDir: modelDir,
	}, nil
}

func cmakeExecutable() string {
	v, err := exec.LookPath("cmake")
	if err == nil && v != "" {
		return v
	}
	return "cmake"
}

func pickDarwinCompiler() string {
	if p, err := exec.LookPath("clang"); err == nil && p != "" {
		return p
	}
	if p, err := exec.LookPath("gcc"); err == nil && p != "" {
		return p
	}
	return ""
}

func pickCXX(cc string) string {
	if cc == "" {
		return ""
	}
	base := filepath.Base(strings.TrimSuffix(cc, filepath.Ext(cc)))
	switch base {
	case "gcc":
		if p, err := exec.LookPath("g++"); err == nil {
			return p
		}
	case "clang":
		if p, err := exec.LookPath("clang++"); err == nil {
			return p
		}
	}
	return cc
}

func locateWhisperCLIUnix(root string) (string, error) {
	var last string
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil {
			return nil
		}
		if info.Mode().IsRegular() && filepath.Base(path) == "whisper-cli" {
			last = path
		}
		return nil
	})
	if last == "" {
		return "", fmt.Errorf("whisper-cli binary not located under %s", root)
	}
	return last, nil
}

func latestWhisperReleaseTag(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, whisperGitHubReleaseAPI, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "HOWL_Chat-whisper-fetch")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", fmt.Errorf("GitHub whisper release API HTTP %s: %s", resp.Status, strings.TrimSpace(string(b)))
	}
	var payload githubReleasePayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}
	if payload.TagName == "" {
		return "", fmt.Errorf("empty whisper release tag from GitHub")
	}
	return payload.TagName, nil
}

func execCMake(ctx context.Context, cwd, cmake string, args, extraEnv []string) error {
	cmd := exec.CommandContext(ctx, cmake, args...)
	cmd.Dir = cwd
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	env := append(os.Environ(), extraEnv...)
	cmd.Env = env
	return cmd.Run()
}

func downloadHTTPFile(ctx context.Context, rawURL string, dest string) error {
	ctx, cancel := context.WithTimeout(ctx, 45*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "HOWL_Chat-whisper-fetch")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("HTTP %s: %s", resp.Status, strings.TrimSpace(string(b)))
	}
	tmp := dest + ".partial"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		out.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := out.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, dest)
}
