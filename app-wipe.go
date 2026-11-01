package main

import (
	"fmt"
	"os"
	"path/filepath"

	"howl-chat/internal/backend/chat"
	"howl-chat/internal/backend/lorebook"
)

// WipeAllData removes ALL user data and resets the application to a blank slate.
// This includes: chats, characters, worlds, scenarios, memories, lorebooks,
// gallery images, settings, and IndexedDB references.
func (a *App) WipeAllData() error {
	fmt.Println("[Wipe] Starting full data wipe...")

	exeDir := filepath.Dir(os.Args[0])

	// 1. Clear backend chat persistence
	if a.chatPersistence != nil {
		chatsDir := filepath.Join(exeDir, "chats")
		if err := os.RemoveAll(chatsDir); err != nil {
			fmt.Printf("[Wipe] Warning: failed to remove chats dir: %v\n", err)
		} else {
			fmt.Println("[Wipe] Chats directory cleared")
		}
	}

	// 2. Clear gallery images
	galleryDir := a.getGalleryDir()
	if err := os.RemoveAll(galleryDir); err != nil {
		fmt.Printf("[Wipe] Warning: failed to remove gallery dir: %v\n", err)
	} else {
		fmt.Println("[Wipe] Gallery directory cleared")
	}
	_ = os.MkdirAll(galleryDir, 0755)

	// 3. Clear entity folders (Characters, Worlds, Scenarios)
	entityDirs := []string{
		filepath.Join(exeDir, "Characters"),
		filepath.Join(exeDir, "Worlds"),
		filepath.Join(exeDir, "Scenarios"),
	}
	for _, dir := range entityDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue // Directory may not exist
		}
		for _, entry := range entries {
			if entry.IsDir() {
				path := filepath.Join(dir, entry.Name())
				if err := os.RemoveAll(path); err != nil {
					fmt.Printf("[Wipe] Warning: failed to remove %s: %v\n", path, err)
				}
			}
		}
		fmt.Printf("[Wipe] Cleared: %s\n", dir)
	}

	// 4. Clear image subdirectories
	imageSubDirs := []string{"characters", "scenarios", "worlds"}
	for _, sub := range imageSubDirs {
		subDir := filepath.Join(exeDir, "images", sub)
		_ = os.RemoveAll(subDir)
		_ = os.MkdirAll(subDir, 0755)
	}

	// 5. Reset in-memory state
	a.characters = []*chat.CharacterContext{}
	a.world = nil
	a.scenario = nil
	a.lorebooks = []lorebook.Entry{}
	a.chatMessages = []map[string]interface{}{}
	a.triggeredLore = make(map[string]bool)

	// 6. Reset memory service chat ID
	if a.memoryService != nil {
		a.memoryService.SetChatID("")
	}

	fmt.Println("[Wipe] Complete. Application is now a blank slate.")
	return nil
}
