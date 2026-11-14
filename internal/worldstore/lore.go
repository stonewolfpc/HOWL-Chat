package worldstore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"howl-chat/internal/library"
)

// maxLoreBytes bounds one lore file. Lore is item content, not a second
// document store, so it stays comfortably below the document limit.
const maxLoreBytes = 4 << 20

// LoreFile is one canonical lore file read through a registered root. Entries
// stay raw so this package does not take ownership of the lorebook schema.
// AuthorityIndeterminate is copied from the item lookup: a found copy whose
// authority cannot be established must not enter a prompt.
type LoreFile struct {
	Entries                []json.RawMessage `json:"entries"`
	Item                   Item              `json:"item"`
	AuthorityIndeterminate bool              `json:"authority_indeterminate,omitempty"`
	UnavailableRoots       []string          `json:"unavailable_roots,omitempty"`
	Errors                 []string          `json:"errors,omitempty"`
}

// LoadLore reads one lore file from the canonical item folder. The file must be
// an allowlisted name for this store's category, so this method cannot be used
// to read arbitrary item-folder content. A missing lore file is empty, not an
// error. An unavailable or indeterminate item returns the authority facts and
// no entries.
func (s *Store) LoadLore(id, filename string) (LoreFile, error) {
	if s == nil || s.reg == nil {
		return LoreFile{}, fmt.Errorf("content library registry is not initialised")
	}
	id = strings.TrimSpace(id)
	if !validID.MatchString(id) {
		return LoreFile{}, fmt.Errorf("%w: identity %q", ErrInvalidDocument, id)
	}
	if filename != filepath.Base(filename) || strings.ContainsAny(filename, `/\`) {
		return LoreFile{}, fmt.Errorf("%w: lore file %q", ErrInvalidDocument, filename)
	}
	if !s.loreFileAllowed(filename) {
		return LoreFile{}, fmt.Errorf("%w: lore file %q is not part of %s", ErrInvalidDocument, filename, s.category)
	}

	search, err := s.Lookup(id)
	if err != nil {
		return LoreFile{}, err
	}
	if !search.Found || search.Item == nil {
		if !search.Complete {
			return LoreFile{}, fmt.Errorf("%w: an item root is unavailable", ErrUnavailable)
		}
		return LoreFile{}, ErrNotFound
	}
	if search.AuthorityIndeterminate {
		return LoreFile{
			Item:                   *search.Item,
			AuthorityIndeterminate: true,
			UnavailableRoots:       search.UnavailableRoots,
			Errors:                 search.Errors,
		}, nil
	}

	path := filepath.Join(search.Item.Root, id, filename)
	payload, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return LoreFile{Item: *search.Item, Entries: []json.RawMessage{}}, nil
		}
		return LoreFile{}, fmt.Errorf("%w: read lore: %v", ErrUnavailable, err)
	}
	if len(payload) > maxLoreBytes {
		return LoreFile{}, fmt.Errorf("%w: lore file exceeds %d bytes", ErrInvalidDocument, maxLoreBytes)
	}
	var doc struct {
		Entries []json.RawMessage `json:"entries"`
	}
	if err := json.Unmarshal(payload, &doc); err != nil {
		return LoreFile{}, fmt.Errorf("%w: decode lore: %v", ErrInvalidDocument, err)
	}
	if doc.Entries == nil {
		doc.Entries = []json.RawMessage{}
	}
	return LoreFile{Item: *search.Item, Entries: doc.Entries}, nil
}

func (s *Store) loreFileAllowed(filename string) bool {
	switch s.category {
	case library.Worlds:
		return filename == "world_lore.json" || filename == "automated_world_lore.json"
	case library.Scenarios:
		return filename == "scenario_lore.json" || filename == "automated_scenario_lore.json"
	default:
		return false
	}
}
