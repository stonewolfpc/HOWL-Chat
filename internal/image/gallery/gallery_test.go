package gallery

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestReadRejectsPathOutsideGallery(t *testing.T) {
	base := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.png")
	if err := os.WriteFile(outside, []byte("outside"), 0644); err != nil {
		t.Fatal(err)
	}
	svc := NewService(base)
	if _, err := svc.Read(outside); !errors.Is(err, ErrOutsideGallery) {
		t.Fatalf("outside read error = %v, want ErrOutsideGallery", err)
	}
	if _, err := svc.Read(filepath.Join(base, "..", filepath.Base(outside))); !errors.Is(err, ErrOutsideGallery) {
		t.Fatalf("traversal read error = %v, want ErrOutsideGallery", err)
	}
}

func TestReadAcceptsFileInsideGallery(t *testing.T) {
	base := t.TempDir()
	image := filepath.Join(base, "img.png")
	if err := os.WriteFile(image, []byte("image"), 0644); err != nil {
		t.Fatal(err)
	}
	svc := NewService(base)
	got, err := svc.Read(image)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "image" {
		t.Fatalf("read = %q", got)
	}
}

func TestReadRejectsSymlinkEscapeWhenSupported(t *testing.T) {
	base := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.png")
	if err := os.WriteFile(outside, []byte("outside"), 0644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(base, "link.png")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	svc := NewService(base)
	if _, err := svc.Read(link); !errors.Is(err, ErrOutsideGallery) {
		t.Fatalf("symlink escape error = %v, want ErrOutsideGallery", err)
	}
}

func TestListIncludesLegitimateInRootImage(t *testing.T) {
	base := t.TempDir()
	image := filepath.Join(base, "entry.png")
	if err := os.WriteFile(image, []byte("image"), 0644); err != nil {
		t.Fatal(err)
	}
	meta, err := json.Marshal(map[string]string{"file": image})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "entry.json"), meta, 0644); err != nil {
		t.Fatal(err)
	}
	items, err := NewService(base).List()
	if err != nil {
		t.Fatal(err)
	}
	expected, err := filepath.EvalSymlinks(image)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].File != expected {
		t.Fatalf("items = %+v, want file %s", items, expected)
	}
}

func TestListRejectsOutsideResolvingImageSymlink(t *testing.T) {
	base := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.png")
	if err := os.WriteFile(outside, []byte("outside"), 0644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(base, "entry.png")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	meta, err := json.Marshal(map[string]string{"file": link})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "entry.json"), meta, 0644); err != nil {
		t.Fatal(err)
	}
	items, err := NewService(base).List()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatalf("outside image symlink was listed: %+v", items)
	}
}

func TestListRejectsOutsideResolvingSidecarSymlink(t *testing.T) {
	base := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.json")
	if err := os.WriteFile(outside, []byte(`{"file":"entry.png"}`), 0644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(base, "entry.json")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	items, err := NewService(base).List()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatalf("outside sidecar symlink was listed: %+v", items)
	}
}

func TestDeleteRejectsTraversalAndReportsMissingItem(t *testing.T) {
	svc := NewService(t.TempDir())
	if err := svc.Delete("../outside"); err == nil {
		t.Fatal("traversal delete was accepted")
	}
	if err := svc.Delete("missing"); err == nil {
		t.Fatal("missing delete was reported as success")
	}
}

func TestDeleteRemovesImageAndSidecar(t *testing.T) {
	base := t.TempDir()
	svc := NewService(base)
	for _, name := range []string{"entry.png", "entry.json"} {
		if err := os.WriteFile(filepath.Join(base, name), []byte("x"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	if err := svc.Delete("entry"); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"entry.png", "entry.json"} {
		if _, err := os.Stat(filepath.Join(base, name)); !os.IsNotExist(err) {
			t.Fatalf("%s still exists: %v", name, err)
		}
	}
}
