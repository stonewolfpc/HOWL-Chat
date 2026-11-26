package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"howl-chat/internal/library"
	"howl-chat/internal/paths"
	"howl-chat/internal/registries"
)

type failingProbeFile struct{ writeErr error }

func (f failingProbeFile) Write(p []byte) (int, error) {
	if f.writeErr != nil {
		return 0, f.writeErr
	}
	return len(p), nil
}
func (f failingProbeFile) Sync() error  { return nil }
func (f failingProbeFile) Close() error { return nil }

func testRegistry(t *testing.T) *library.Registry {
	t.Helper()
	reg, err := library.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return reg
}

func addLibrary(t *testing.T, reg *library.Registry, category library.Category, dir string) library.Registration {
	t.Helper()
	registration, err := reg.Add(library.Registration{
		Category:    category,
		Path:        dir,
		Role:        library.RoleRead | library.RoleWrite | library.RoleSave,
		Origin:      library.OriginHowl,
		Format:      library.FormatHowl,
		DefaultSave: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	return registration
}

func readyProbe(category library.Category, dir string) ProvisionReport {
	return ProvisionReport{
		Category: string(category), Path: dir, RootCreated: true, CategoryCreated: true,
		FileCreated: true, BytesWritten: 16, ReadBack: true, Cleanup: true,
		Readable: true, Writable: true, Status: ContentReady,
	}
}

func TestProbeProvesRootCategoryFileReadbackAndCleanup(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "Howl", "Worlds")

	report := probeContentDir(defaultContentFSOps(), string(library.Worlds), dir)
	if report.Err != nil {
		t.Fatalf("probe failed: %v (%s: %s)", report.Err, report.Step, report.Reason)
	}
	if !report.RootCreated || !report.CategoryCreated || !report.FileCreated || !report.ReadBack || !report.Cleanup {
		t.Fatalf("probe did not report every step: %+v", report)
	}
	if report.BytesWritten != len("howl-write-probe") {
		t.Fatalf("bytes written = %d", report.BytesWritten)
	}
	if !report.Writable || !report.Readable || report.Status != ContentReady {
		t.Fatalf("probe state = %+v", report)
	}
	// Non-destructive cleanup removed only what this call created.
	if !report.RootRemoved || !report.CategoryRemoved {
		t.Fatalf("probe left created directories behind: %+v", report)
	}
	if _, err := os.Stat(filepath.Join(base, "Howl")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("probe left %q behind: %v", dir, err)
	}
}

func TestProbeKeepsAnExistingDirectory(t *testing.T) {
	dir := t.TempDir()
	report := probeContentDir(defaultContentFSOps(), string(library.Worlds), dir)
	if report.Err != nil {
		t.Fatalf("probe failed: %v", report.Err)
	}
	if report.CategoryCreated || report.CategoryRemoved {
		t.Fatalf("existing directory was treated as created: %+v", report)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("probe left artifacts: %v", entries)
	}
}

func TestProbeFailsOnARegularFileWhereADirectoryBelongs(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "Worlds")
	if err := os.WriteFile(dir, []byte("not a directory"), 0o600); err != nil {
		t.Fatal(err)
	}
	report := probeContentDir(defaultContentFSOps(), string(library.Worlds), dir)
	if report.Err == nil {
		t.Fatal("probe accepted a file as a category directory")
	}
	if report.Step != "inspect category directory" || report.Status != ContentUnavailable {
		t.Fatalf("report = %+v", report)
	}
	if report.Retryable {
		t.Fatal("a blocked path was reported as retryable")
	}
	if !strings.Contains(report.Reason, "is a file, not a directory") {
		t.Fatalf("reason = %q", report.Reason)
	}
}

func TestProbeReportsEachFailingStep(t *testing.T) {
	t.Run("create library root", func(t *testing.T) {
		base := t.TempDir()
		ops := defaultContentFSOps()
		ops.MkdirAll = func(string, os.FileMode) error { return fs.ErrPermission }
		report := probeContentDir(ops, string(library.Worlds), filepath.Join(base, "Howl", "Worlds"))
		if report.Step != "create library root" || report.Status != ContentUnavailable {
			t.Fatalf("report = %+v", report)
		}
	})

	t.Run("create category directory", func(t *testing.T) {
		base := t.TempDir()
		root := filepath.Join(base, "Howl")
		if err := os.MkdirAll(root, 0o755); err != nil {
			t.Fatal(err)
		}
		ops := defaultContentFSOps()
		ops.MkdirAll = func(string, os.FileMode) error { return fs.ErrPermission }
		report := probeContentDir(ops, string(library.Worlds), filepath.Join(root, "Worlds"))
		if report.Step != "create category directory" || report.Status != ContentUnavailable {
			t.Fatalf("report = %+v", report)
		}
	})

	t.Run("create probe file in an existing readable directory is read-only", func(t *testing.T) {
		dir := t.TempDir()
		ops := defaultContentFSOps()
		ops.OpenFile = func(string) (probeFile, error) { return nil, fs.ErrPermission }
		report := probeContentDir(ops, string(library.Worlds), dir)
		if report.Step != "create probe file" || report.Status != ContentReadOnly {
			t.Fatalf("report = %+v", report)
		}
		if !report.Readable || report.Writable {
			t.Fatalf("read-only state = %+v", report)
		}
		if !strings.Contains(report.Reason, "can be read but not written") {
			t.Fatalf("reason = %q", report.Reason)
		}
	})

	t.Run("write probe file", func(t *testing.T) {
		dir := t.TempDir()
		ops := defaultContentFSOps()
		ops.OpenFile = func(string) (probeFile, error) {
			return failingProbeFile{writeErr: fs.ErrPermission}, nil
		}
		report := probeContentDir(ops, string(library.Worlds), dir)
		if report.Step != "write probe file" || report.Writable {
			t.Fatalf("report = %+v", report)
		}
		entries, _ := os.ReadDir(dir)
		if len(entries) != 0 {
			t.Fatalf("failed probe left the probe file behind: %v", entries)
		}
	})

	t.Run("read back probe file", func(t *testing.T) {
		dir := t.TempDir()
		ops := defaultContentFSOps()
		ops.ReadFile = func(string) ([]byte, error) { return nil, fs.ErrPermission }
		report := probeContentDir(ops, string(library.Worlds), dir)
		if report.Step != "read back probe file" || report.ReadBack || report.Writable {
			t.Fatalf("report = %+v", report)
		}
	})

	t.Run("verify probe file", func(t *testing.T) {
		dir := t.TempDir()
		ops := defaultContentFSOps()
		ops.ReadFile = func(string) ([]byte, error) { return []byte("different"), nil }
		report := probeContentDir(ops, string(library.Worlds), dir)
		if report.Step != "verify probe file" || report.Writable {
			t.Fatalf("report = %+v", report)
		}
	})

	t.Run("remove probe file", func(t *testing.T) {
		dir := t.TempDir()
		ops := defaultContentFSOps()
		ops.Remove = func(string) error { return fs.ErrPermission }
		report := probeContentDir(ops, string(library.Worlds), dir)
		if report.Step != "remove probe file" || report.Cleanup || report.Writable {
			t.Fatalf("report = %+v", report)
		}
	})
}

func TestReadinessKeepsPerLibraryStatesAndCountsFailures(t *testing.T) {
	reg := testRegistry(t)
	writable := addLibrary(t, reg, library.Worlds, t.TempDir())
	restricted := addLibrary(t, reg, library.Worlds, t.TempDir())

	readOnly := ProvisionReport{
		Category: string(library.Worlds), Path: restricted.Path,
		Readable: true, Writable: false, Step: "create probe file", Operation: "create probe file",
		Status:    ContentReadOnly,
		Reason:    "This worlds library can be read but not written: Windows denied access.",
		Retryable: true,
		Err:       fs.ErrPermission,
	}
	readiness := buildContentReadiness(reg, map[string]ProvisionReport{
		writable.ID:   readyProbe(library.Worlds, writable.Path),
		restricted.ID: readOnly,
	}, nil)

	byStatus := map[string]ContentLibraryState{}
	worldsLibraries := 0
	for _, state := range readiness.Libraries {
		if state.Status == ContentNotInitialized {
			continue // synthesized entry for a category with no registration (scenarios)
		}
		worldsLibraries++
		byStatus[state.Status] = state
		if state.Category != string(library.Worlds) || state.LibraryID == "" || state.Path == "" {
			t.Fatalf("library state lost its identity: %+v", state)
		}
	}
	if worldsLibraries != 2 {
		t.Fatalf("libraries = %#v", readiness.Libraries)
	}
	if byStatus[ContentReady].LibraryID != writable.ID {
		t.Fatalf("writable library lost its ready state: %+v", byStatus)
	}
	if byStatus[ContentReadOnly].LibraryID != restricted.ID {
		t.Fatalf("read-only library was not reported read-only: %+v", byStatus)
	}
	if !byStatus[ContentReadOnly].Readable {
		t.Fatal("read-only library must stay readable")
	}
	var worlds *ContentCategoryState
	for i := range readiness.Categories {
		if readiness.Categories[i].Category == string(library.Worlds) {
			worlds = &readiness.Categories[i]
		}
	}
	if worlds == nil {
		t.Fatalf("categories = %#v", readiness.Categories)
	}
	if worlds.Status != ContentReadOnly || worlds.FailingLibraries != 1 ||
		worlds.Libraries != 2 || !worlds.Writable {
		t.Fatalf("rollup hid or misreported a failing library: %+v", *worlds)
	}
	if !strings.Contains(worlds.Reason, "1 of 2") {
		t.Fatalf("rollup reason = %q", worlds.Reason)
	}
}

func TestReadinessDoesNotShareStatusAcrossCategories(t *testing.T) {
	reg := testRegistry(t)
	scenarios := addLibrary(t, reg, library.Scenarios, t.TempDir())

	readiness := buildContentReadiness(reg, map[string]ProvisionReport{
		scenarios.ID: readyProbe(library.Scenarios, scenarios.Path),
	}, map[string]ProvisionReport{
		string(library.Worlds): {
			Category: string(library.Worlds), Path: filepath.Join(t.TempDir(), "Worlds"),
			Status: ContentUnavailable, Operation: "create category directory", Retryable: true,
			Reason: "Windows denied access while trying to create category directory.",
			Err:    fs.ErrPermission,
		},
	})

	statuses := map[string]string{}
	for _, category := range readiness.Categories {
		statuses[category.Category] = category.Status
	}
	if statuses[string(library.Worlds)] != ContentUnavailable {
		t.Fatalf("worlds = %q", statuses[string(library.Worlds)])
	}
	if statuses[string(library.Scenarios)] != ContentReady {
		t.Fatalf("scenarios inherited the worlds failure: %q", statuses[string(library.Scenarios)])
	}
}

func TestReadinessSynthesizesNotInitializedPerCategory(t *testing.T) {
	readiness := buildContentReadiness(testRegistry(t), nil, nil)
	if len(readiness.Libraries) != 2 {
		t.Fatalf("libraries = %#v", readiness.Libraries)
	}
	for _, state := range readiness.Libraries {
		if state.Status != ContentNotInitialized {
			t.Fatalf("state = %+v", state)
		}
		if state.Category == "" || state.Path == "" || state.Reason == "" || !state.Retryable {
			t.Fatalf("synthesized state is missing identity or remediation: %+v", state)
		}
	}
	if len(readiness.Categories) != 2 {
		t.Fatalf("categories = %#v", readiness.Categories)
	}
}

func TestRequireCreateWritableRefusesAReadOnlyDefaultRoot(t *testing.T) {
	dir := t.TempDir()
	reg := testRegistry(t)
	registration := addLibrary(t, reg, library.Worlds, dir)
	app := &App{registries: registries.Compose(registries.Parts{Libraries: reg})}
	app.recordProbe(registration.ID, ProvisionReport{
		Category: string(library.Worlds), Path: dir, Readable: true,
		Step: "create probe file", Operation: "create probe file", Status: ContentReadOnly,
		Reason: "This worlds library can be read but not written: Windows denied access.", Retryable: true,
		Err: fs.ErrPermission,
	})

	err := app.requireCreateWritable(library.Worlds, "SaveWorld")
	if err == nil {
		t.Fatal("a read-only default root was accepted for a create")
	}
	if !strings.Contains(err.Error(), "can be read but not written") {
		t.Fatalf("refusal = %q", err)
	}

	if _, err := app.SaveWorld(map[string]interface{}{"title": "Blocked"}); err == nil {
		t.Fatal("SaveWorld wrote to a read-only library")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("nothing should have been written: %v", entries)
	}
}

func TestReadsDoNotCreateContentDirectories(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, paths.IsolationMarkerName), []byte(paths.IsolationMarkerContent), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(paths.IsolationEnvVar, root)

	reg := testRegistry(t)
	app := &App{registries: registries.Compose(registries.Parts{Libraries: reg})}
	dir, err := paths.DefaultContentDir(string(library.Worlds))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dir); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("precondition: %q already exists (%v)", dir, err)
	}

	worlds, err := app.GetWorlds()
	if err != nil {
		t.Fatalf("GetWorlds on an uninitialized library: %v", err)
	}
	if len(worlds) != 0 {
		t.Fatalf("worlds = %#v", worlds)
	}
	if _, err := os.Stat(dir); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("GetWorlds created %q", dir)
	}
	if _, err := os.Stat(filepath.Dir(dir)); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("GetWorlds created the library root")
	}
	if len(reg.List(library.Worlds)) != 0 {
		t.Fatal("GetWorlds registered a library")
	}
}

func TestEnsureDefaultContentLibraryProvisionsAndRegisters(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, paths.IsolationMarkerName), []byte(paths.IsolationMarkerContent), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(paths.IsolationEnvVar, root)

	reg := testRegistry(t)
	app := &App{registries: registries.Compose(registries.Parts{Libraries: reg})}

	state, err := app.EnsureDefaultContentLibrary("worlds")
	if err != nil {
		t.Fatalf("ensure failed: %v (%s: %s)", err, state.Operation, state.Reason)
	}
	if state.Status != ContentReady || !state.Readable || !state.Writable || state.LibraryID == "" {
		t.Fatalf("state = %+v", state)
	}
	if len(reg.List(library.Worlds)) != 1 {
		t.Fatalf("registration missing: %#v", reg.List(library.Worlds))
	}
	if _, err := os.Stat(state.Path); err != nil {
		t.Fatalf("provisioned directory missing: %v", err)
	}
}

func TestProbeObservationOverridesViewForUnenumerableLibrary(t *testing.T) {
	reg := testRegistry(t)
	registration := addLibrary(t, reg, library.Worlds, t.TempDir())

	// The registry's view says the directory exists and has a write role; the
	// probe could not enumerate it or write into it. The probe is authoritative.
	readiness := buildContentReadiness(reg, map[string]ProvisionReport{
		registration.ID: {
			Category: string(library.Worlds), Path: registration.Path,
			Readable: false, Writable: false,
			Step: "create probe file", Operation: "create probe file",
			Status: ContentUnavailable, Retryable: true,
			Reason: "Windows denied access while trying to create probe file.",
			Err:    fs.ErrPermission,
		},
	}, nil)

	var worlds *ContentLibraryState
	for i := range readiness.Libraries {
		if readiness.Libraries[i].LibraryID == registration.ID {
			worlds = &readiness.Libraries[i]
		}
	}
	if worlds == nil {
		t.Fatalf("library missing from readiness: %#v", readiness.Libraries)
	}
	if worlds.Readable || worlds.Writable || worlds.Status == ContentReady {
		t.Fatalf("view claimed access the probe disproved: %+v", *worlds)
	}
}

func TestRollupKeepsHealthySiblingUsableAndCreateFollowsDefaultSave(t *testing.T) {
	reg := testRegistry(t)
	readOnlySave := addLibrary(t, reg, library.Worlds, t.TempDir())
	writableSibling := addLibrary(t, reg, library.Worlds, t.TempDir())
	// Make the sibling the only writable library; the default-save root is read-only.
	if err := reg.SetDefaultSave(readOnlySave.ID); err != nil {
		t.Fatal(err)
	}

	readiness := buildContentReadiness(reg, map[string]ProvisionReport{
		readOnlySave.ID: {
			Category: string(library.Worlds), Path: readOnlySave.Path, Readable: true,
			Step: "create probe file", Operation: "create probe file", Status: ContentReadOnly,
			Reason: "read-only default root", Retryable: true, Err: fs.ErrPermission,
		},
		writableSibling.ID: readyProbe(library.Worlds, writableSibling.Path),
	}, nil)

	var worlds *ContentCategoryState
	for i := range readiness.Categories {
		if readiness.Categories[i].Category == string(library.Worlds) {
			worlds = &readiness.Categories[i]
		}
	}
	if worlds == nil {
		t.Fatal("worlds category missing")
	}
	if !worlds.Writable {
		t.Fatal("a healthy sibling was disabled by a failing library")
	}
	if worlds.SaveWritable {
		t.Fatal("create was advertised as possible on a read-only default save root")
	}
	if worlds.FailingLibraries != 1 || worlds.Status != ContentReadOnly {
		t.Fatalf("rollup lost the failing library: %+v", *worlds)
	}

	app := &App{registries: registries.Compose(registries.Parts{Libraries: reg})}
	app.recordProbe(readOnlySave.ID, ProvisionReport{
		Category: string(library.Worlds), Path: readOnlySave.Path, Readable: true,
		Status: ContentReadOnly, Step: "create probe file", Operation: "create probe file",
		Reason: "This worlds library can be read but not written.", Retryable: true, Err: fs.ErrPermission,
	})
	if err := app.requireCreateWritable(library.Worlds, "SaveWorld"); err == nil {
		t.Fatal("create was accepted on a read-only default save root")
	}
}

func TestCreateIsAllowedWhenTheDefaultSaveRootIsWritable(t *testing.T) {
	reg := testRegistry(t)
	writableSave := addLibrary(t, reg, library.Worlds, t.TempDir())
	readOnlySibling := addLibrary(t, reg, library.Worlds, t.TempDir())
	if err := reg.SetDefaultSave(writableSave.ID); err != nil {
		t.Fatal(err)
	}

	readiness := buildContentReadiness(reg, map[string]ProvisionReport{
		writableSave.ID: readyProbe(library.Worlds, writableSave.Path),
		readOnlySibling.ID: {
			Category: string(library.Worlds), Path: readOnlySibling.Path, Readable: true,
			Status: ContentReadOnly, Step: "create probe file", Operation: "create probe file",
			Reason: "read-only sibling", Retryable: true, Err: fs.ErrPermission,
		},
	}, nil)

	var worlds *ContentCategoryState
	for i := range readiness.Categories {
		if readiness.Categories[i].Category == string(library.Worlds) {
			worlds = &readiness.Categories[i]
		}
	}
	if worlds == nil || !worlds.SaveWritable || !worlds.Writable {
		t.Fatalf("create disabled by a read-only sibling: %+v", worlds)
	}
	if worlds.FailingLibraries != 1 {
		t.Fatalf("the read-only sibling was hidden: %+v", *worlds)
	}

	app := &App{registries: registries.Compose(registries.Parts{Libraries: reg})}
	app.recordProbe(writableSave.ID, readyProbe(library.Worlds, writableSave.Path))
	if err := app.requireCreateWritable(library.Worlds, "SaveWorld"); err != nil {
		t.Fatalf("create refused although the default save root is writable: %v", err)
	}
}
