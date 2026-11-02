package deps

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"howl-chat/internal/image/engine/sdcpp"
)

// The report must name every role the engine will demand at load time. This is
// the same table the loader reads, so the two cannot drift.
func TestResolveNamesEveryRequiredCompanion(t *testing.T) {
	for family := range map[string]struct{}{"flux": {}, "flux2": {}, "z_image": {}, "qwen_image": {}, "anima": {}, "sd3.5": {}} {
		cap, ok := sdcpp.Capability(family)
		if !ok {
			t.Fatalf("%s missing from capability table", family)
		}
		r := Resolve("model", family, nil, "/models/companions")
		if len(r.Missing) != len(cap.Companions) {
			t.Errorf("%s: missing = %v, want %v", family, r.Missing, cap.Companions)
		}
		for _, role := range cap.Companions {
			found := false
			for _, it := range r.Items {
				if it.Role == role {
					found = true
					if it.State != StateMissing {
						t.Errorf("%s/%s: state = %q, want missing", family, role, it.State)
					}
				}
			}
			if !found {
				t.Errorf("%s: no item for required role %q", family, role)
			}
		}
	}
}

// A file the loader would use is reported as present with its path, never as
// something to download again.
func TestResolveMarksProvidedFilesPresent(t *testing.T) {
	present := map[string]string{
		"llm": `E:\models\qwen_3_4b.safetensors`,
		"vae": `E:\models\ae.safetensors`,
	}
	r := Resolve("z_image_turbo", "z_image", present, "/models/companions")
	if r.BlocksLoad() {
		t.Fatalf("missing = %v, want none", r.Missing)
	}
	for _, it := range r.Items {
		if it.Role == "llm" && (it.State != StatePresent || it.Path != present["llm"]) {
			t.Errorf("llm item = %+v", it)
		}
	}
	if n := len(r.Downloadable()); n != 0 {
		t.Errorf("downloadable = %d, want 0", n)
	}
}

// Optional roles must not block a load, and must not be offered as a download.
func TestResolveTreatsOptionalRolesAsOptional(t *testing.T) {
	r := Resolve("sd15", "sd1.5", nil, "/models/companions")
	if r.BlocksLoad() {
		t.Fatalf("sd1.5 should load with no companions, missing = %v", r.Missing)
	}
	for _, it := range r.Items {
		if it.Role == "clip_l" && it.State != StateOptional {
			t.Errorf("clip_l state = %q, want optional", it.State)
		}
	}
}

// Every required role must lead somewhere useful: a direct link, or a page plus
// an explanation. A missing file with neither is a dead end for the user.
func TestEveryRequiredRoleHasAWayToGetIt(t *testing.T) {
	roles := []string{"vae", "clip_l", "clip_g", "t5xxl", "llm"}
	for family := range Catalog {
		cap, known := sdcpp.Capability(family)
		if !known {
			t.Errorf("%s is in the catalogue but not a supported family", family)
			continue
		}
		for _, role := range cap.Companions {
			if !contains(roles, role) {
				t.Errorf("%s: required role %q has no acquisition path in this package", family, role)
				continue
			}
			a := artifactFor(family, role)
			if a.Source == "" {
				t.Errorf("%s/%s: no source page", family, role)
			}
			if a.URL == "" && a.Note == "" {
				t.Errorf("%s/%s: no direct link and no note explaining what to do", family, role)
			}
			if a.URL != "" && a.FileName == "" {
				t.Errorf("%s/%s: has a URL but no file name", family, role)
			}
		}
	}
}

func TestDownloadWritesFileAndReportsProgress(t *testing.T) {
	payload := strings.Repeat("weights", 2000)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", itoa(len(payload)))
		_, _ = w.Write([]byte(payload))
	}))
	defer srv.Close()

	dir := t.TempDir()
	var calls int
	path, err := Download(context.Background(), srv.Client(), Artifact{FileName: "ae.safetensors", URL: srv.URL}, dir, func(string, int64, int64) { calls++ })
	if err != nil {
		t.Fatal(err)
	}
	if path != filepath.Join(dir, "ae.safetensors") {
		t.Fatalf("path = %q", path)
	}
	got, err := os.ReadFile(path)
	if err != nil || string(got) != payload {
		t.Fatalf("content mismatch (err %v)", err)
	}
	if calls == 0 {
		t.Error("progress was never reported")
	}
	if _, err := os.Stat(path + ".part"); !os.IsNotExist(err) {
		t.Error("partial file left behind")
	}
}

func TestDownloadLeavesNoPartialFileOnServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusNotFound)
	}))
	defer srv.Close()

	dir := t.TempDir()
	_, err := Download(context.Background(), srv.Client(), Artifact{FileName: "vae.safetensors", URL: srv.URL}, dir, nil)
	if err == nil {
		t.Fatal("expected an error for HTTP 404")
	}
	if IsOffline(err) {
		t.Errorf("a 404 is a server answer, not an offline failure: %v", err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Errorf("dir should be empty after a failed download, has %d entries", len(entries))
	}
}

func TestDownloadClassifiesUnreachableHostAsOffline(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := srv.URL
	srv.Close() // nothing is listening now

	_, err := Download(context.Background(), &http.Client{}, Artifact{FileName: "vae.safetensors", URL: url}, t.TempDir(), nil)
	if err == nil {
		t.Fatal("expected an error against a closed host")
	}
	if !IsOffline(err) {
		t.Fatalf("want ErrOffline, got %v", err)
	}
	if !errors.Is(err, ErrOffline) {
		t.Fatalf("error must wrap ErrOffline, got %v", err)
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}

func contains(list []string, want string) bool {
	for _, s := range list {
		if s == want {
			return true
		}
	}
	return false
}
