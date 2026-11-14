package kokoro

import (
	"archive/zip"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestListVoicesFromArchive(t *testing.T) {
	path := filepath.Join(t.TempDir(), "voices.bin")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(f)
	for _, name := range []string{"af_heart.npy", "am_michael.npy", "nested/bf_emma.npy", "readme.txt"} {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte("x")); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	got, err := ListVoicesFromArchive(path)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"af_heart", "am_michael", "bf_emma"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("voices = %v, want %v", got, want)
	}
}
