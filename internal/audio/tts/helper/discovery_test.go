package helper

import (
	"os"
	"path/filepath"
	"testing"
)

func helperFile(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "runtime", "tts", "howl-tts-helper.exe")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("helper"), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestDiscoveryUsesEnvironmentFirst(t *testing.T) {
	root := t.TempDir()
	envPath := filepath.Join(root, "env-helper.exe")
	if err := os.WriteFile(envPath, []byte("env"), 0o644); err != nil {
		t.Fatal(err)
	}
	appPath := helperFile(t, filepath.Join(root, "app"))
	devPath := helperFile(t, filepath.Join(root, "dev"))

	d := Discoverer{
		LookupEnv:   func(string) (string, bool) { return envPath, true },
		Executable:  func() (string, error) { return filepath.Join(root, "app", "howl.exe"), nil },
		ProjectRoot: filepath.Join(root, "dev"),
	}
	got, failure := d.Discover()
	if failure != nil {
		t.Fatal(failure)
	}
	if got.Path != envPath || got.Source != "HOWL_TTS_HELPER_PATH" {
		t.Fatalf("discovered %+v, want env path %s", got, envPath)
	}
	_ = appPath
	_ = devPath
}

func TestDiscoveryUsesApplicationThenDevelopment(t *testing.T) {
	root := t.TempDir()
	appRoot := filepath.Join(root, "app")
	devRoot := filepath.Join(root, "dev")
	appPath := helperFile(t, appRoot)
	devPath := helperFile(t, devRoot)

	d := Discoverer{
		LookupEnv:   func(string) (string, bool) { return "", false },
		Executable:  func() (string, error) { return filepath.Join(appRoot, "howl.exe"), nil },
		ProjectRoot: devRoot,
	}
	got, failure := d.Discover()
	if failure != nil {
		t.Fatal(failure)
	}
	if got.Path != appPath {
		t.Fatalf("discovered %s, want application path %s", got.Path, appPath)
	}

	if err := os.Remove(appPath); err != nil {
		t.Fatal(err)
	}
	got, failure = d.Discover()
	if failure != nil {
		t.Fatal(failure)
	}
	if got.Path != devPath {
		t.Fatalf("discovered %s, want development path %s", got.Path, devPath)
	}
}

func TestDiscoveryRejectsRelativeEnvironmentPath(t *testing.T) {
	d := Discoverer{LookupEnv: func(string) (string, bool) { return `relative\helper.exe`, true }}
	_, failure := d.Discover()
	if failure == nil || failure.Code != "helper_missing" {
		t.Fatalf("failure = %v, want helper_missing", failure)
	}
}

func TestDiscoveryMissingDoesNotUsePATH(t *testing.T) {
	root := t.TempDir()
	d := Discoverer{
		LookupEnv:  func(string) (string, bool) { return "", false },
		Executable: func() (string, error) { return filepath.Join(root, "howl.exe"), nil },
	}
	_, failure := d.Discover()
	if failure == nil || failure.Code != "helper_missing" {
		t.Fatalf("failure = %v, want helper_missing", failure)
	}
}
