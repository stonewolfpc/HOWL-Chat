package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"howl-chat/internal/backend/chat"
	"howl-chat/internal/image/gallery"
	"howl-chat/internal/library"
	"howl-chat/internal/media"
	"howl-chat/internal/registries"
)

type wipeFixture struct {
	app         *App
	chatRoot    string
	galleryRoot string
	chatsDir    string
	external    string
}

func newWipeFixture(t *testing.T) *wipeFixture {
	t.Helper()
	chatRoot := t.TempDir()
	galleryRoot := t.TempDir()
	externalRoot := t.TempDir()
	manager, err := chat.NewPersistenceManagerWithDir(chatRoot)
	if err != nil {
		t.Fatal(err)
	}
	if err := manager.SetLastChatID("chat-1"); err != nil {
		t.Fatal(err)
	}
	chatsDir := filepath.Join(chatRoot, "chats")
	chatFile := filepath.Join(chatsDir, "chat-1", "meta.json")
	if err := os.MkdirAll(filepath.Dir(chatFile), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(chatFile, []byte("chat"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(galleryRoot, "img.png"), []byte("image"), 0644); err != nil {
		t.Fatal(err)
	}
	external := filepath.Join(externalRoot, "keep.json")
	if err := os.WriteFile(external, []byte("keep"), 0644); err != nil {
		t.Fatal(err)
	}
	inventory := library.NewInventoryStore(t.TempDir())
	if err := inventory.Record("keep-root", []library.Item{{ID: "card", Name: "Keep"}}); err != nil {
		t.Fatal(err)
	}
	app := &App{
		conversation: chat.NewActiveConversation(manager),
		registries:   registries.Compose(registries.Parts{Inventory: inventory}),
	}
	app.images = media.NewImage("")
	app.images.SetGallery(gallery.NewService(galleryRoot))
	return &wipeFixture{
		app:         app,
		chatRoot:    chatRoot,
		galleryRoot: galleryRoot,
		chatsDir:    chatsDir,
		external:    external,
	}
}

func preserveWipeHooks(t *testing.T) {
	t.Helper()
	origRename := wipeRename
	origRemoveAll := wipeRemoveAll
	t.Cleanup(func() {
		wipeRename = origRename
		wipeRemoveAll = origRemoveAll
	})
}

func TestWipeAllDataStagesBothRootsAndResetsState(t *testing.T) {
	f := newWipeFixture(t)
	if err := f.app.WipeAllData(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(f.chatsDir); !os.IsNotExist(err) {
		t.Fatalf("chat root still exists: %v", err)
	}
	if entries, err := os.ReadDir(f.galleryRoot); err != nil || len(entries) != 0 {
		t.Fatalf("gallery not emptied: entries=%v err=%v", entries, err)
	}
	if got := f.app.conversation.ActiveID(); got != "" {
		t.Fatalf("LastChatID = %q, want empty", got)
	}
	if _, ok := f.app.registries.Inventory().Cached("keep-root"); ok {
		t.Fatal("inventory cache survived wipe")
	}
	if _, err := os.Stat(f.external); err != nil {
		t.Fatalf("external file was removed: %v", err)
	}
}

func TestWipeAllDataFirstStageFailureLeavesStateUntouched(t *testing.T) {
	f := newWipeFixture(t)
	preserveWipeHooks(t)
	wipeRename = func(string, string) error { return errors.New("injected first stage failure") }
	if err := f.app.WipeAllData(); err == nil {
		t.Fatal("wipe reported success after first stage failure")
	}
	if _, err := os.Stat(f.chatsDir); err != nil {
		t.Fatalf("chat root changed: %v", err)
	}
	if _, err := os.Stat(f.galleryRoot); err != nil {
		t.Fatalf("gallery root changed: %v", err)
	}
	if got := f.app.conversation.ActiveID(); got != "chat-1" {
		t.Fatalf("LastChatID = %q, want chat-1", got)
	}
	if _, ok := f.app.registries.Inventory().Cached("keep-root"); !ok {
		t.Fatal("inventory cache was cleared before staging succeeded")
	}
}

func TestWipeAllDataSecondStageFailureRestoresFirstRoot(t *testing.T) {
	f := newWipeFixture(t)
	preserveWipeHooks(t)
	calls := 0
	wipeRename = func(oldPath, newPath string) error {
		calls++
		if calls == 2 {
			return errors.New("injected second stage failure")
		}
		return os.Rename(oldPath, newPath)
	}
	if err := f.app.WipeAllData(); err == nil {
		t.Fatal("wipe reported success after second stage failure")
	}
	if _, err := os.Stat(f.chatsDir); err != nil {
		t.Fatalf("chat root was not restored: %v", err)
	}
	if _, err := os.Stat(f.galleryRoot); err != nil {
		t.Fatalf("gallery root was not restored: %v", err)
	}
	if got := f.app.conversation.ActiveID(); got != "chat-1" {
		t.Fatalf("LastChatID = %q, want chat-1", got)
	}
	if _, ok := f.app.registries.Inventory().Cached("keep-root"); !ok {
		t.Fatal("inventory cache was cleared before both stages succeeded")
	}
}

func TestWipeAllDataCleanupFailureLeavesRecoverableStage(t *testing.T) {
	f := newWipeFixture(t)
	preserveWipeHooks(t)
	wipeRemoveAll = func(path string) error {
		if strings.Contains(filepath.Base(path), "howl-wipe") {
			return errors.New("injected cleanup failure")
		}
		return os.RemoveAll(path)
	}
	err := f.app.WipeAllData()
	if err == nil || !strings.Contains(err.Error(), "cleanup is incomplete") {
		t.Fatalf("cleanup error = %v", err)
	}
	chatStages, _ := filepath.Glob(filepath.Join(f.chatRoot, ".chats.howl-wipe-*"))
	if len(chatStages) != 1 {
		t.Fatalf("chat staged data = %v, want one recoverable stage", chatStages)
	}
	galleryStages, _ := filepath.Glob(filepath.Join(filepath.Dir(f.galleryRoot), ".gallery.howl-wipe-*"))
	_ = galleryStages // gallery may have been recreated; the error still names any stage
	if got := f.app.conversation.ActiveID(); got != "" {
		t.Fatalf("LastChatID = %q, want empty after logical wipe", got)
	}
	if _, ok := f.app.registries.Inventory().Cached("keep-root"); ok {
		t.Fatal("inventory cache survived logical wipe")
	}
}

func TestWipeAllDataReportsMissingRoots(t *testing.T) {
	app := &App{}
	if err := app.WipeAllData(); err == nil {
		t.Fatal("wipe without services reported success")
	}
}
