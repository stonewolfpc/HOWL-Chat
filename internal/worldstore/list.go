package worldstore

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// List returns every readable canonical document in the registered roots. A
// duplicate identity is returned once and marked authority-indeterminate, and
// an unreadable root does not make a readable copy disappear: it marks that
// copy indeterminate rather than silently treating it as canon.
func (s *Store) List() ([]Item, error) {
	if s == nil || s.reg == nil {
		return nil, fmt.Errorf("content library registry is not initialised")
	}
	views := s.reg.List(s.category)
	byID := map[string][]Item{}
	incomplete := false
	for _, view := range views {
		if !view.Exists || !view.Readable {
			incomplete = true
			continue
		}
		entries, err := os.ReadDir(view.Path)
		if err != nil {
			incomplete = true
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			path := filepath.Join(view.Path, entry.Name(), s.file)
			payload, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			doc, err := decodeDocument(payload, s.format, entry.Name())
			if err != nil {
				continue
			}
			computed := Revision(doc.Document)
			byID[doc.ID] = append(byID[doc.ID], Item{
				ID:               doc.ID,
				Format:           doc.Format,
				Revision:         doc.Revision,
				ComputedRevision: computed,
				RevisionMismatch: computed != doc.Revision,
				CreatedAt:        doc.CreatedAt,
				UpdatedAt:        doc.UpdatedAt,
				Origin:           doc.Origin,
				LibraryID:        view.ID,
				Root:             view.Path,
				Path:             path,
				Available:        true,
				Editable:         view.Writable,
				Document:         doc.Document,
			})
		}
	}

	ids := make([]string, 0, len(byID))
	for id := range byID {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	items := make([]Item, 0, len(ids))
	for _, id := range ids {
		copies := byID[id]
		item := copies[0]
		if len(copies) > 1 || incomplete {
			item.AuthorityIndeterminate = true
			item.Editable = false
		}
		items = append(items, item)
	}
	return items, nil
}

// Delete removes one canonical document folder. Absence must be proven, the
// copy must be authoritative, and the root must grant write permission; every
// refusal is returned as a status rather than attempted.
func (s *Store) Delete(id string) (SaveResult, error) {
	res := SaveResult{Status: "write_failed", ID: id}
	if s == nil || s.reg == nil {
		res.Status = "destination_unavailable"
		res.Message = "content library registry is not initialised"
		return res, nil
	}
	search, err := s.Lookup(id)
	if err != nil {
		return res, err
	}
	if !search.Found || search.Item == nil {
		if !search.Complete {
			res.Status = "destination_unavailable"
			res.Message = "No readable library could prove whether this item exists; nothing was removed."
			return res, nil
		}
		res.Status = "not_found"
		res.Message = "The world/scenario is not in a canonical library."
		return res, nil
	}
	item := *search.Item
	if search.AuthorityIndeterminate {
		res.Status = "authority_indeterminate"
		res.Message = "The item was found, but its authority cannot be established, so nothing was removed."
		return res, nil
	}
	if !item.Editable {
		res.Status = "permission_denied"
		res.Message = "The library holding this item does not grant write permission."
		return res, nil
	}
	dir := filepath.Dir(item.Path)
	rel, err := filepath.Rel(item.Root, dir)
	if err != nil || rel == "." || rel == ".." || filepath.IsAbs(rel) {
		res.Status = "unsafe_target"
		res.Message = "The item folder is outside its registered library."
		return res, nil
	}
	if err := os.RemoveAll(dir); err != nil {
		if errors.Is(err, os.ErrPermission) {
			res.Status = "permission_denied"
			res.Message = "The destination denied permission to remove the item."
			return res, nil
		}
		return res, s.mapWriteError(&res, err)
	}
	res.Status = "deleted"
	res.Path = dir
	return res, nil
}
