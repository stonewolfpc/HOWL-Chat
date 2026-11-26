package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	apperrors "howl-chat/internal/errors"
	"howl-chat/internal/library"
	"howl-chat/internal/paths"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Content-library bindings. These expose registration and discovery only;
// removing a registration never removes a directory or a file.
func (a *App) GetContentLibraries(category string) ([]library.View, error) {
	if a.registries.Libraries() == nil {
		return []library.View{}, fmt.Errorf("content library registry not initialised")
	}
	cat := library.Category(strings.TrimSpace(category))
	if category != "" && !cat.Valid() {
		return nil, fmt.Errorf("unknown content category %q", category)
	}
	return a.registries.Libraries().List(cat), nil
}

func (a *App) AddContentLibrary(path, category, origin, format string, role int) (library.Registration, error) {
	if a.registries.Libraries() == nil {
		return library.Registration{}, fmt.Errorf("content library registry not initialised")
	}
	cat := library.Category(strings.TrimSpace(category))
	registration, err := a.registries.Libraries().Add(library.Registration{
		Category: cat,
		Path:     path,
		Role:     library.NormalizeRole(library.Role(role)),
		Origin:   library.Origin(strings.TrimSpace(origin)),
		Format:   library.Format(strings.TrimSpace(format)),
	})
	if err != nil {
		return library.Registration{}, err
	}
	// A user-added root is probed once, here, rather than on every read. A
	// readable but unwritable root is still registered: its content stays
	// visible and only mutations are refused.
	if registration.Role.CanWrite() || registration.Role.CanSave() {
		report := probeContentDir(a.contentOpsOrDefault(), string(cat), registration.Path, false)
		a.recordProbe(registration.ID, report)
	}
	return registration, nil
}

// AddDefaultContentLibrary creates the approved Documents/Pictures/Videos
// location for a category and registers it as a Howl-owned save root. It is an
// explicit user action and delegates to the one provisioning authority
// (app-content-roots.go); merely reading content never creates or moves data.
func (a *App) AddDefaultContentLibrary(category string, role int) (library.Registration, error) {
	cat := library.Category(strings.TrimSpace(category))
	if !cat.Valid() {
		return library.Registration{}, fmt.Errorf("unknown content category %q", category)
	}
	report, registration := a.provisionDefaultRoot(cat, library.Role(role))
	if report.Err != nil {
		a.recordInit(cat, report)
		return library.Registration{}, fmt.Errorf("%s", report.Reason)
	}
	return registration, nil
}

func (a *App) RemoveContentLibrary(id string) error {
	if a.registries.Libraries() == nil {
		return fmt.Errorf("content library registry not initialised")
	}
	if err := a.registries.Libraries().Remove(id); err != nil {
		return err
	}
	if a.registries.Inventory() != nil {
		a.registries.Inventory().Forget(id)
	}
	return nil
}

func (a *App) SetContentLibraryRole(id string, role int) error {
	if a.registries.Libraries() == nil {
		return fmt.Errorf("content library registry not initialised")
	}
	return a.registries.Libraries().SetRole(id, library.Role(role))
}

func (a *App) SetDefaultContentLibrary(id string) error {
	if a.registries.Libraries() == nil {
		return fmt.Errorf("content library registry not initialised")
	}
	return a.registries.Libraries().SetDefaultSave(id)
}

// GetContentLibraryCharacters returns character cards from every readable
// character root. A previously indexed item from an unavailable root remains
// visible as unavailable rather than disappearing. The same physical card
// discovered through overlapping roots appears once. No file is modified and
// no card is imported.
func (a *App) GetContentLibraryCharacters() ([]library.Item, error) {
	if a.registries.Libraries() == nil {
		return []library.Item{}, fmt.Errorf("content library registry not initialised")
	}
	seen := map[string]bool{}
	out := []library.Item{}
	for _, view := range a.registries.Libraries().List(library.Characters) {
		var items []library.Item
		live := false
		if view.Readable {
			scanned, err := library.ScanCharacters(view.Registration)
			if err == nil {
				items = scanned
				live = true
				if a.registries.Inventory() != nil {
					if recordErr := a.registries.Inventory().Record(view.ID, scanned); recordErr != nil {
						a.logErrorSilent("SYS-CONTENT-LIBRARY", apperrors.Low, "library.inventory.record", "app-libraries.go", recordErr)
					}
				}
			} else {
				a.logErrorSilent("SYS-CONTENT-LIBRARY", apperrors.Medium, "library.characters.scan", "app-libraries.go", err)
			}
		}
		if !live && a.registries.Inventory() != nil {
			if cached, ok := a.registries.Inventory().Cached(view.ID); ok {
				items = cached.Items
				for i := range items {
					items[i].Available = false
					items[i].Editable = false
				}
			}
		}
		for _, item := range items {
			key := strings.ToLower(filepath.Clean(item.Path))
			if seen[key] {
				continue
			}
			seen[key] = true
			out = append(out, item)
		}
	}
	return out, nil
}

// GetContentLibraryInventory returns every registration, including roots whose
// drive or folder is unavailable. Successful character scans become the
// last-known snapshot; an unavailable root returns that snapshot marked stale.
// It never scans a drive recursively beyond the registered root.
func (a *App) GetContentLibraryInventory(category string) ([]library.InventoryRoot, error) {
	if a.registries.Libraries() == nil {
		return []library.InventoryRoot{}, fmt.Errorf("content library registry not initialised")
	}
	cat := library.Category(strings.TrimSpace(category))
	if category != "" && !cat.Valid() {
		return nil, fmt.Errorf("unknown content category %q", category)
	}
	views := a.registries.Libraries().List(cat)
	out := make([]library.InventoryRoot, 0, len(views))
	for _, view := range views {
		root := library.InventoryRoot{
			Registration: view.Registration,
			Available:    view.Readable,
		}
		live := false
		if view.Readable && view.Category == library.Characters {
			items, err := library.ScanCharacters(view.Registration)
			if err == nil {
				root.Items = items
				root.Scanned = true
				root.ScannedAt = time.Now().UTC().Format(time.RFC3339)
				live = true
				if a.registries.Inventory() != nil {
					if recordErr := a.registries.Inventory().Record(view.ID, items); recordErr != nil {
						a.logErrorSilent("SYS-CONTENT-LIBRARY", apperrors.Low, "library.inventory.record", "app-libraries.go", recordErr)
					}
				}
			} else {
				root.ScanError = err.Error()
				a.logErrorSilent("SYS-CONTENT-LIBRARY", apperrors.Medium, "library.inventory.scan", "app-libraries.go", err)
			}
		}
		if !live && a.registries.Inventory() != nil {
			if cached, ok := a.registries.Inventory().Cached(view.ID); ok {
				root.Items = cached.Items
				root.ScannedAt = cached.ScannedAt
				root.Stale = true
				for i := range root.Items {
					root.Items[i].Available = false
					root.Items[i].Editable = false
				}
			}
		}
		out = append(out, root)
	}
	return out, nil
}

// DefaultContentLibraryPath reports the approved location without creating it.
func (a *App) DefaultContentLibraryPath(category string) (string, error) {
	return paths.DefaultContentDir(strings.TrimSpace(category))
}

// OpenContentDirectoryPicker opens the native folder chooser for a content
// registration. It does not create or modify anything by itself.
func (a *App) OpenContentDirectoryPicker() (string, error) {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Content Folder",
	})
	if err != nil {
		return "", fmt.Errorf("content directory dialog: %w", err)
	}
	return dir, nil
}

// ExportContentItem is the only content-library write binding. It creates an
// explicitly identified export copy; it is never called by a normal page save.
// New exports never replace an existing file. Replacement requires a writable
// root and replace=true, and is returned as an explicit status. Worlds and
// scenarios are refused: their canonical envelopes are written only by
// worldstore, and a second writer under those roots would break A10.
func (a *App) ExportContentItem(category, name, itemID string, data map[string]interface{}, replace bool) (library.OperationResult, error) {
	result := library.OperationResult{Category: library.Category(strings.TrimSpace(category)), ItemID: strings.TrimSpace(itemID)}
	if a.registries.Libraries() == nil {
		return result, fmt.Errorf("content library registry not initialised")
	}
	cat := result.Category
	if !cat.Valid() {
		return result, fmt.Errorf("unknown content category %q", category)
	}
	switch cat {
	case library.Images, library.Videos:
		return result, fmt.Errorf("generated %s are written by their generation pipeline", cat)
	case library.Chats:
		return result, fmt.Errorf("chat export is not part of the content-library export contract")
	case library.Worlds, library.Scenarios:
		return result, fmt.Errorf("%s export is not available: canonical items are written only by the world/scenario store", cat)
	}
	reg, ok := a.registries.Libraries().DefaultSave(cat)
	if !ok {
		result.Status = library.StatusDestinationUnavailable
		result.Message = "No save destination is configured for " + string(cat) + ". The content was not exported."
		return result, nil
	}

	if result.ItemID == "" && data == nil {
		return result, fmt.Errorf("export requires an item identity or content data")
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return result, fmt.Errorf("encode %s: %w", cat, err)
	}
	identityBasis := result.ItemID
	if identityBasis == "" {
		identityBasis = string(dataBytes)
	}
	sum := sha256.Sum256([]byte(identityBasis))
	identity := hex.EncodeToString(sum[:6])
	// The directory is identity-only so renaming an item cannot create a second
	// export. The display title is retained in the envelope, not used as identity.
	segment := identity
	file := contentItemFile(cat)
	payload, err := json.MarshalIndent(map[string]interface{}{
		"format":        "howl-content-export-v1",
		"exported_at":   time.Now().UTC().Format(time.RFC3339),
		"category":      cat,
		"item_id":       result.ItemID,
		"display_name":  safeContentSegment(name),
		"source":        "howl",
		"authoritative": false,
		"content":       data,
	}, "", "  ")
	if err != nil {
		return result, fmt.Errorf("encode %s export: %w", cat, err)
	}

	relative := filepath.Join(segment, file)
	var target string
	if replace {
		target, err = a.registries.Libraries().WriteFile(reg.ID, relative, payload)
	} else {
		target, err = a.registries.Libraries().WriteNewFile(reg.ID, relative, payload)
	}
	if err == nil {
		result.Target = target
		if replace {
			result.Status = library.StatusReplaced
			result.Message = "Existing export was explicitly replaced."
		} else {
			result.Status = library.StatusExportSucceeded
			result.Message = "Export copy created."
		}
		return result, nil
	}

	result.Target = target
	switch {
	case errors.Is(err, library.ErrCollision):
		result.Status = library.StatusNameCollision
		result.Message = "An export with this identity already exists. Choose replacement explicitly or export under a different identity."
		return result, nil
	case errors.Is(err, library.ErrReadOnly):
		result.Status = library.StatusPermissionDenied
		result.Message = "The selected library does not grant the permission required for this export."
		return result, nil
	case errors.Is(err, library.ErrUnsafeTarget):
		result.Status = library.StatusFailed
		result.Message = "The export target is a directory, symlink, or other non-regular file. It was not replaced."
		return result, nil
	case errors.Is(err, library.ErrInvalidPath):
		result.Status = library.StatusFailed
		result.Message = "The export path is invalid."
		return result, nil
	case errors.Is(err, os.ErrPermission):
		result.Status = library.StatusPermissionDenied
		result.Message = "The destination denied permission to write the export."
		return result, nil
	default:
		result.Status = library.StatusFailed
		result.Message = "The export failed: " + err.Error()
		return result, nil
	}
}

func contentItemFile(cat library.Category) string {
	switch cat {
	case library.Chats:
		return "chat.json"
	case library.Characters:
		return "character.json"
	case library.Lorebooks:
		return "lorebook.json"
	default:
		return "content.json"
	}
}

func safeContentSegment(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "Untitled"
	}
	var b strings.Builder
	for _, r := range name {
		if strings.ContainsRune(`<>:"/\|?*`, r) || r < 32 {
			b.WriteRune('_')
		} else {
			b.WriteRune(r)
		}
	}
	out := strings.Trim(b.String(), " .")
	if out == "" {
		out = "Untitled"
	}
	upper := strings.ToUpper(out)
	stem := upper
	if i := strings.IndexByte(stem, '.'); i >= 0 {
		stem = stem[:i]
	}
	switch stem {
	case "CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9":
		out = "Item_" + out
	}
	runes := []rune(out)
	if len(runes) > 100 {
		out = string(runes[:100])
	}
	return out
}
