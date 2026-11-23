package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// These are variables so failure-injection tests can exercise the wipe
// transaction without touching real user content.
var (
	wipeRename    = os.Rename
	wipeRemoveAll = os.RemoveAll
)

// WipeAllData removes the application-owned content roots it can identify:
// the configured chat root and the generated gallery. It does not scan or
// delete arbitrary drives, external content-library roots, or user files.
// Settings reset is a separate contract and is never invoked here.
//
// The two roots are staged to same-filesystem sibling paths before any logical
// state is reset. If the second stage fails, the first is restored. If staging
// succeeds, the active chat pointer and inventory cache are reset before the
// staged data is removed, so a completed logical wipe cannot leave a stale
// active chat or stale cached cards behind.
func (a *App) WipeAllData() error {
	fmt.Println("[Wipe] Starting application content wipe...")

	if !a.conversation.Bound() {
		return fmt.Errorf("content wipe not started: chat persistence is not initialized")
	}
	if a.images.Gallery() == nil {
		return fmt.Errorf("content wipe not started: gallery service is not initialized")
	}
	chatRoot := a.conversation.RootDir()
	galleryDir := a.images.Gallery().BaseDir()
	if chatRoot == "" || galleryDir == "" {
		return fmt.Errorf("content wipe not started: application content root is empty")
	}

	chatsDir := filepath.Join(chatRoot, "chats")
	oldLastChatID := a.conversation.ActiveID()

	chatStage, err := stageContentRoot(chatsDir)
	if err != nil {
		return fmt.Errorf("content wipe not started: stage chat root: %w", err)
	}
	galleryStage, err := stageContentRoot(galleryDir)
	if err != nil {
		if restoreErr := restoreContentRoot(chatStage, chatsDir); restoreErr != nil {
			return fmt.Errorf("content wipe not started: stage gallery root: %w; chat restore failed: %v", err, restoreErr)
		}
		return fmt.Errorf("content wipe not started: stage gallery root: %w", err)
	}

	// At this point both roots are recoverable if a state reset fails.
	if err := a.conversation.SetActiveID(""); err != nil {
		restoreErr := restoreStagedRoots(chatStage, chatsDir, galleryStage, galleryDir)
		if restoreErr != nil {
			return fmt.Errorf("content wipe not completed: reset active chat: %w; root restore: %v", err, restoreErr)
		}
		return fmt.Errorf("content wipe not completed: reset active chat: %w", err)
	}
	if a.registries.Inventory() != nil {
		if err := a.registries.Inventory().Clear(); err != nil {
			pointerErr := a.conversation.SetActiveID(oldLastChatID)
			restoreErr := restoreStagedRoots(chatStage, chatsDir, galleryStage, galleryDir)
			return fmt.Errorf("content wipe not completed: clear inventory: %w; pointer restore: %v; root restore: %v", err, pointerErr, restoreErr)
		}
	}

	// The logical wipe has succeeded. Reset in-memory state only now.
	a.session().Reset()
	if a.memoryService != nil {
		a.memoryService.SetChatID("")
	}

	var cleanupFailures []error
	if chatStage != "" {
		if err := wipeRemoveAll(chatStage); err != nil {
			cleanupFailures = append(cleanupFailures, fmt.Errorf("chat data remains staged at %s: %w", chatStage, err))
		}
	}
	if galleryStage != "" {
		if err := wipeRemoveAll(galleryStage); err != nil {
			cleanupFailures = append(cleanupFailures, fmt.Errorf("gallery data remains staged at %s: %w", galleryStage, err))
		}
	}
	if err := os.MkdirAll(galleryDir, 0755); err != nil {
		cleanupFailures = append(cleanupFailures, fmt.Errorf("recreate gallery root %s: %w", galleryDir, err))
	}

	if len(cleanupFailures) > 0 {
		return fmt.Errorf("content wipe completed, but cleanup is incomplete: %w", errors.Join(cleanupFailures...))
	}
	fmt.Println("[Wipe] Complete.")
	return nil
}

func stageContentRoot(path string) (string, error) {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory", path)
	}
	staged := filepath.Join(filepath.Dir(path), fmt.Sprintf(".%s.howl-wipe-%d", filepath.Base(path), time.Now().UnixNano()))
	if err := wipeRename(path, staged); err != nil {
		return "", err
	}
	return staged, nil
}

func restoreContentRoot(staged, original string) error {
	if staged == "" {
		return nil
	}
	return wipeRename(staged, original)
}

func restoreStagedRoots(chatStage, chatsDir, galleryStage, galleryDir string) error {
	var failures []error
	if err := restoreContentRoot(chatStage, chatsDir); err != nil {
		failures = append(failures, fmt.Errorf("restore chat root: %w", err))
	}
	if err := restoreContentRoot(galleryStage, galleryDir); err != nil {
		failures = append(failures, fmt.Errorf("restore gallery root: %w", err))
	}
	if len(failures) > 0 {
		return errors.Join(failures...)
	}
	return nil
}
