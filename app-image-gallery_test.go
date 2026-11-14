package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"howl-chat/internal/image/gallery"
)

func TestGetGalleryImageRejectsOutsideGallery(t *testing.T) {
	base := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.png")
	if err := os.WriteFile(outside, []byte("outside"), 0644); err != nil {
		t.Fatal(err)
	}
	app := &App{gallerySvc: gallery.NewService(base)}
	if _, err := app.GetGalleryImage(outside); !errors.Is(err, gallery.ErrOutsideGallery) {
		t.Fatalf("outside read error = %v, want ErrOutsideGallery", err)
	}
}

func TestGetGalleryImageReadsInsideGallery(t *testing.T) {
	base := t.TempDir()
	image := filepath.Join(base, "image.png")
	if err := os.WriteFile(image, []byte("inside"), 0644); err != nil {
		t.Fatal(err)
	}
	app := &App{gallerySvc: gallery.NewService(base)}
	data, err := app.GetGalleryImage(image)
	if err != nil {
		t.Fatal(err)
	}
	if data == "" {
		t.Fatal("empty base64 image")
	}
}
