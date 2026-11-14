package helper

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestWindowsReparseComponentValidation(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows junction evidence only")
	}
	root := t.TempDir()
	target := filepath.Join(root, "target")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "junction")
	cmd := exec.Command("cmd", "/c", "mklink", "/J", link, target)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("could not create a test junction: %v (%s)", err, out)
	}
	nested := filepath.Join(link, "nested")
	if err := os.MkdirAll(filepath.Join(target, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got, err := hasReparseComponent(nested); err != nil || !got {
		t.Fatalf("intermediate junction not detected: got=%v err=%v", got, err)
	}
	helper := filepath.Join(target, "helper.exe")
	if err := os.WriteFile(helper, []byte("test double"), 0o644); err != nil {
		t.Fatal(err)
	}
	if regularHelperFile(filepath.Join(link, "helper.exe")) {
		t.Fatal("helper beneath an intermediate junction was accepted")
	}
	if err := validateOutput(link, "out.wav"); err == nil {
		t.Fatal("output directory beneath an intermediate junction was accepted")
	}
	output := filepath.Join(target, "nested", "out.wav")
	if _, _, err := writeTestWAV(output, 24000, 240); err != nil {
		t.Fatal(err)
	}
	if _, failure := validateWAVFile(synthesisResult{}, link+string(filepath.Separator)+"nested", "out.wav"); failure == nil {
		t.Fatal("output beneath an intermediate junction was accepted")
	}
	managed := filepath.Join(target, "out.wav")
	if _, _, err := writeTestWAV(managed, 24000, 240); err != nil {
		t.Fatal(err)
	}
	cleanupExact(link, "out.wav")
	if _, err := os.Stat(managed); err != nil {
		t.Fatalf("cleanup traversed an intermediate junction: %v", err)
	}
}

func TestWindowsFinalSymlinkValidation(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows symlink evidence only")
	}
	root := t.TempDir()
	target := filepath.Join(root, "target.wav")
	if _, _, err := writeTestWAV(target, 24000, 240); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "linked.wav")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("could not create a test symlink: %v", err)
	}
	if regularHelperFile(link) {
		t.Fatal("final symlink helper was accepted")
	}
	if got, err := hasReparseComponent(link); err != nil || !got {
		t.Fatalf("final symlink not detected: got=%v err=%v", got, err)
	}
	if _, failure := validateWAVFile(synthesisResult{}, root, "linked.wav"); failure == nil {
		t.Fatal("final symlink output was accepted")
	}
}
