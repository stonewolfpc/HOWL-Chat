package helper

import (
	"os"
	"path/filepath"
)

// Discoverer resolves the helper without searching PATH or System32.
type Discoverer struct {
	LookupEnv   func(string) (string, bool)
	Executable  func() (string, error)
	ProjectRoot string
}

// DiscoveredPath records both the path and which ordered rule selected it.
type DiscoveredPath struct {
	Path   string
	Source string
}

func (d Discoverer) lookupEnv() func(string) (string, bool) {
	if d.LookupEnv != nil {
		return d.LookupEnv
	}
	return os.LookupEnv
}

func (d Discoverer) executable() func() (string, error) {
	if d.Executable != nil {
		return d.Executable
	}
	return os.Executable
}

// Discover returns exactly one path from the frozen three-rule order.
func (d Discoverer) Discover() (DiscoveredPath, *Failure) {
	if value, ok := d.lookupEnv()("HOWL_TTS_HELPER_PATH"); ok {
		if !filepath.IsAbs(value) {
			return DiscoveredPath{}, newFailure("helper_missing", "HOWL_TTS_HELPER_PATH must be an absolute file path")
		}
		if !regularHelperFile(value) {
			return DiscoveredPath{}, newFailure("helper_missing", "configured TTS helper is not a regular file")
		}
		return DiscoveredPath{Path: filepath.Clean(value), Source: "HOWL_TTS_HELPER_PATH"}, nil
	}

	if exe, err := d.executable()(); err == nil && exe != "" {
		candidate := filepath.Join(filepath.Dir(exe), "runtime", "tts", "howl-tts-helper.exe")
		if regularHelperFile(candidate) {
			return DiscoveredPath{Path: candidate, Source: "application runtime"}, nil
		}
	}

	if d.ProjectRoot != "" {
		candidate := filepath.Join(filepath.Clean(d.ProjectRoot), "runtime", "tts", "howl-tts-helper.exe")
		if regularHelperFile(candidate) {
			return DiscoveredPath{Path: candidate, Source: "development runtime"}, nil
		}
	}

	return DiscoveredPath{}, newFailure("helper_missing", "optional TTS helper was not found")
}

func regularHelperFile(path string) bool {
	if path == "" || !filepath.IsAbs(path) {
		return false
	}
	info, err := os.Lstat(path)
	if err != nil || info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
		return false
	}
	reparse, reparseErr := hasReparseComponent(path)
	return reparseErr == nil && !reparse
}
