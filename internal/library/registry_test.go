package library

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMultipleRootsRolesAndRemovalPreserveFiles(t *testing.T) {
	configDir := t.TempDir()
	rootA := filepath.Join(t.TempDir(), "chars-a")
	rootB := filepath.Join(t.TempDir(), "chars-b")
	for _, dir := range []string{rootA, rootB} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}
	fileA := filepath.Join(rootA, "a.json")
	if err := os.WriteFile(fileA, []byte(`{"name":"A"}`), 0644); err != nil {
		t.Fatal(err)
	}

	r, err := New(configDir)
	if err != nil {
		t.Fatal(err)
	}
	a, err := r.Add(Registration{Category: Characters, Path: rootA, Role: RoleRead, Origin: OriginExternal, Format: FormatSillyTavern})
	if err != nil {
		t.Fatal(err)
	}
	b, err := r.Add(Registration{Category: Characters, Path: rootB, Role: RoleRead | RoleSave, Origin: OriginHowl, Format: FormatHowl})
	if err != nil {
		t.Fatal(err)
	}
	if !b.DefaultSave {
		t.Fatal("first save root should become the default")
	}
	if err := r.SetRole(a.ID, RoleRead|RoleSave); err != nil {
		t.Fatal(err)
	}
	if err := r.SetDefaultSave(a.ID); err != nil {
		t.Fatal(err)
	}
	if got := r.List(Characters); len(got) != 2 {
		t.Fatalf("got %d roots, want 2", len(got))
	}
	if def, ok := r.DefaultSave(Characters); !ok || def.ID != a.ID {
		t.Fatalf("default save = %+v, want %s", def, a.ID)
	}
	if err := r.Remove(a.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(fileA); err != nil {
		t.Fatalf("removing registration removed the file: %v", err)
	}
	reloaded, err := New(configDir)
	if err != nil {
		t.Fatal(err)
	}
	if got := reloaded.List(Characters); len(got) != 1 || got[0].ID != b.ID {
		t.Fatalf("persisted roots = %+v, want only %s", got, b.ID)
	}
}

func TestSameDirectoryCanServeTwoCategories(t *testing.T) {
	dir := t.TempDir()
	r, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := r.Add(Registration{Category: Characters, Path: dir, Role: RoleRead}); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Add(Registration{Category: Lorebooks, Path: dir, Role: RoleRead}); err != nil {
		t.Fatal(err)
	}
	if len(r.List(Characters)) != 1 || len(r.List(Lorebooks)) != 1 {
		t.Fatal("independent category registrations were not preserved")
	}
}

func TestWriteNewFileRejectsReadOnlyAndEscapes(t *testing.T) {
	root := t.TempDir()
	r, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	reg, err := r.Add(Registration{Category: Characters, Path: root, Role: RoleRead})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := r.WriteNewFile(reg.ID, "x.json", []byte("x")); err == nil {
		t.Fatal("read-only library accepted a write")
	}
	if err := r.SetRole(reg.ID, RoleRead|RoleWrite); err != nil {
		t.Fatal(err)
	}
	if _, err := r.WriteNewFile(reg.ID, "../escape.json", []byte("x")); err == nil {
		t.Fatal("path traversal was accepted")
	}
	if _, err := r.WriteNewFile(reg.ID, "nested/card.json", []byte("x")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "nested", "card.json")); err != nil {
		t.Fatalf("write did not create file: %v", err)
	}
}

func TestSaveDestinationCanCreateButNotReplaceOwnedContent(t *testing.T) {
	root := t.TempDir()
	r, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	reg, err := r.Add(Registration{Category: Worlds, Path: root, Role: RoleSave, Origin: OriginHowl, Format: FormatHowl})
	if err != nil {
		t.Fatal(err)
	}
	if !reg.DefaultSave {
		t.Fatal("save destination was not selected as default")
	}
	path, err := r.WriteNewFile(reg.ID, "Example/world.json", []byte(`{"title":"one"}`))
	if err != nil {
		t.Fatalf("save-only library rejected new content: %v", err)
	}
	if _, err := r.WriteFile(reg.ID, "Example/world.json", []byte(`{"title":"two"}`)); !errors.Is(err, ErrReadOnly) {
		t.Fatalf("save-only replace error = %v, want ErrReadOnly", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"title":"one"}` {
		t.Fatalf("existing content changed without write permission: %s", data)
	}
	if err := r.SetRole(reg.ID, RoleRead|RoleWrite|RoleSave); err != nil {
		t.Fatal(err)
	}
	if _, err := r.WriteFile(reg.ID, "Example/world.json", []byte(`{"title":"two"}`)); err != nil {
		t.Fatalf("explicit replacement failed: %v", err)
	}
	data, err = os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"title":"two"}` {
		t.Fatalf("replaced content = %s", data)
	}
}

func TestWriteFileRejectsDirectoryTarget(t *testing.T) {
	root := t.TempDir()
	r, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	reg, err := r.Add(Registration{Category: Worlds, Path: root, Role: RoleSave | RoleWrite, Origin: OriginHowl, Format: FormatHowl})
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "world.json")
	if err := os.Mkdir(target, 0755); err != nil {
		t.Fatal(err)
	}
	if _, err := r.WriteFile(reg.ID, "world.json", []byte("new")); !errors.Is(err, ErrUnsafeTarget) {
		t.Fatalf("directory target error = %v, want ErrUnsafeTarget", err)
	}
	info, err := os.Stat(target)
	if err != nil || !info.IsDir() {
		t.Fatalf("directory target was changed: info=%v err=%v", info, err)
	}
}

func TestWriteFileRejectsSymlinkTargetWhenSupported(t *testing.T) {
	root := t.TempDir()
	r, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	reg, err := r.Add(Registration{Category: Worlds, Path: root, Role: RoleSave | RoleWrite, Origin: OriginHowl, Format: FormatHowl})
	if err != nil {
		t.Fatal(err)
	}
	real := filepath.Join(root, "real.json")
	if err := os.WriteFile(real, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "world.json")
	if err := os.Symlink(real, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if _, err := r.WriteFile(reg.ID, "world.json", []byte("new")); !errors.Is(err, ErrUnsafeTarget) {
		t.Fatalf("symlink target error = %v, want ErrUnsafeTarget", err)
	}
	data, err := os.ReadFile(real)
	if err != nil || string(data) != "old" {
		t.Fatalf("symlink target changed: data=%q err=%v", data, err)
	}
}

func TestReplacementCommitAndRestoreFailureAreReported(t *testing.T) {
	root := t.TempDir()
	r, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	reg, err := r.Add(Registration{Category: Worlds, Path: root, Role: RoleSave | RoleWrite, Origin: OriginHowl, Format: FormatHowl})
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "world.json")
	if err := os.WriteFile(target, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	calls := 0
	r.renameFile = func(oldPath, newPath string) error {
		calls++
		switch calls {
		case 2:
			return errors.New("injected commit failure")
		case 3:
			return errors.New("injected restore failure")
		default:
			return os.Rename(oldPath, newPath)
		}
	}
	_, err = r.WriteFile(reg.ID, "world.json", []byte("new"))
	if err == nil || !strings.Contains(err.Error(), "restore failed") {
		t.Fatalf("replacement error = %v, want restore failure", err)
	}
	if _, statErr := os.Stat(target); !os.IsNotExist(statErr) {
		t.Fatalf("target exists after failed restore: %v", statErr)
	}
	backups, err := filepath.Glob(target + ".bak-*")
	if err != nil || len(backups) != 1 {
		t.Fatalf("backups = %v, err=%v, want one recoverable backup", backups, err)
	}
	data, err := os.ReadFile(backups[0])
	if err != nil || string(data) != "old" {
		t.Fatalf("backup content = %q, err=%v, want old", data, err)
	}
}

func TestReplacementCommitFailureRestoresOriginal(t *testing.T) {
	root := t.TempDir()
	r, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	reg, err := r.Add(Registration{Category: Worlds, Path: root, Role: RoleSave | RoleWrite, Origin: OriginHowl, Format: FormatHowl})
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "world.json")
	if err := os.WriteFile(target, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	calls := 0
	r.renameFile = func(oldPath, newPath string) error {
		calls++
		if calls == 2 {
			return errors.New("injected commit failure")
		}
		return os.Rename(oldPath, newPath)
	}
	if _, err := r.WriteFile(reg.ID, "world.json", []byte("new")); err == nil {
		t.Fatal("commit failure was reported as success")
	}
	data, err := os.ReadFile(target)
	if err != nil || string(data) != "old" {
		t.Fatalf("restored content = %q, err=%v, want old", data, err)
	}
}

func TestWriteFailurePreservesExistingContent(t *testing.T) {
	root := t.TempDir()
	r, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	reg, err := r.Add(Registration{Category: Worlds, Path: root, Role: RoleSave | RoleWrite, Origin: OriginHowl, Format: FormatHowl})
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "world.json")
	if err := os.WriteFile(target, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(target+".tmp", 0755); err != nil {
		t.Fatal(err)
	}
	if _, err := r.WriteFile(reg.ID, "world.json", []byte("new")); err == nil {
		t.Fatal("write with blocked staging path succeeded")
	}
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "old" {
		t.Fatalf("failed write changed existing content: %q", data)
	}
}

func TestWriteNewFileRejectsCollision(t *testing.T) {
	root := t.TempDir()
	r, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	reg, err := r.Add(Registration{Category: Lorebooks, Path: root, Role: RoleSave, Origin: OriginHowl, Format: FormatHowl})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := r.WriteNewFile(reg.ID, "same.json", []byte("one")); err != nil {
		t.Fatal(err)
	}
	if _, err := r.WriteNewFile(reg.ID, "same.json", []byte("two")); !errors.Is(err, ErrCollision) {
		t.Fatalf("collision error = %v, want ErrCollision", err)
	}
	data, err := os.ReadFile(filepath.Join(root, "same.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "one" {
		t.Fatalf("collision overwrote existing content: %s", data)
	}
}
