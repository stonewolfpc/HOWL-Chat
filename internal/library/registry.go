// Package library owns user-selected content roots. A root is a registration,
// not a copy: removing one forgets the association and never touches the files.
package library

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Category string

const (
	Chats      Category = "chats"
	Characters Category = "characters"
	Worlds     Category = "worlds"
	Scenarios  Category = "scenarios"
	Lorebooks  Category = "lorebooks"
	Images     Category = "images"
	Videos     Category = "videos"
)

func (c Category) Valid() bool {
	switch c {
	case Chats, Characters, Worlds, Scenarios, Lorebooks, Images, Videos:
		return true
	}
	return false
}

type Role uint8

const (
	RoleRead  Role = 1 << 0
	RoleWrite Role = 1 << 1
	RoleSave  Role = 1 << 2
)

func NormalizeRole(role Role) Role {
	if role == 0 {
		return RoleRead
	}
	if role&(RoleWrite|RoleSave) != 0 {
		role |= RoleRead
	}
	return role
}

func (r Role) CanRead() bool  { return r&RoleRead != 0 }
func (r Role) CanWrite() bool { return r&RoleWrite != 0 }
func (r Role) CanSave() bool  { return r&RoleSave != 0 }

type Origin string

const (
	OriginHowl     Origin = "howl"
	OriginExternal Origin = "external"
	OriginImported Origin = "imported"
)

type Format string

const (
	FormatHowl        Format = "howl"
	FormatSillyTavern Format = "sillytavern"
	FormatUnknown     Format = "unknown"
)

var (
	ErrNotRegistered = errors.New("content library is not registered")
	ErrReadOnly      = errors.New("content library is read-only")
	ErrCollision     = errors.New("content path already exists")
	ErrInvalidPath   = errors.New("content path is invalid")
	ErrUnsafeTarget  = errors.New("content replacement target is not a regular file")
)

type Registration struct {
	ID          string   `json:"id"`
	Category    Category `json:"category"`
	Path        string   `json:"path"`
	Role        Role     `json:"role"`
	Origin      Origin   `json:"origin"`
	Format      Format   `json:"format"`
	DefaultSave bool     `json:"default_save"`
	CreatedAt   string   `json:"created_at"`
}

type View struct {
	Registration
	Exists    bool `json:"exists"`
	Readable  bool `json:"readable"`
	Writable  bool `json:"writable"`
	SaveRoot  bool `json:"save_root"`
	Available bool `json:"available"`
}

type Registry struct {
	mu         sync.RWMutex
	path       string
	items      map[string]Registration
	renameFile func(string, string) error
}

type persistedRegistry struct {
	Version int            `json:"version"`
	Items   []Registration `json:"items"`
}

func New(configDir string) (*Registry, error) {
	r := &Registry{
		path:       filepath.Join(configDir, "content-libraries.json"),
		items:      map[string]Registration{},
		renameFile: os.Rename,
	}
	data, err := os.ReadFile(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return r, nil
		}
		return nil, fmt.Errorf("read content libraries: %w", err)
	}
	var stored persistedRegistry
	if err := json.Unmarshal(data, &stored); err != nil {
		return nil, fmt.Errorf("parse content libraries: %w", err)
	}
	for _, item := range stored.Items {
		if !item.Category.Valid() || strings.TrimSpace(item.Path) == "" {
			continue
		}
		item.Role = NormalizeRole(item.Role)
		if item.ID == "" {
			item.ID = registrationID(item.Category, item.Path)
		}
		r.items[item.ID] = item
	}
	return r, nil
}

func (r *Registry) List(category Category) []View {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]View, 0, len(r.items))
	for _, item := range r.items {
		if category != "" && item.Category != category {
			continue
		}
		out = append(out, viewFor(item))
	}
	sortViews(out)
	return out
}

func (r *Registry) Add(item Registration) (Registration, error) {
	if !item.Category.Valid() {
		return Registration{}, fmt.Errorf("unknown content category %q", item.Category)
	}
	if strings.TrimSpace(item.Path) == "" {
		return Registration{}, fmt.Errorf("library path is empty")
	}
	abs, err := filepath.Abs(item.Path)
	if err != nil {
		abs = filepath.Clean(item.Path)
	}
	info, err := os.Stat(abs)
	if err != nil {
		return Registration{}, fmt.Errorf("cannot access %q: %w", abs, err)
	}
	if !info.IsDir() {
		return Registration{}, fmt.Errorf("%q is not a directory", abs)
	}
	item.Path = filepath.Clean(abs)
	item.Role = NormalizeRole(item.Role)
	if item.Origin == "" {
		item.Origin = OriginExternal
	}
	if item.Format == "" {
		item.Format = FormatUnknown
	}
	if item.ID == "" {
		item.ID = registrationID(item.Category, item.Path)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC().Format(time.RFC3339)
	if old, ok := r.items[item.ID]; ok {
		if old.CreatedAt != "" {
			item.CreatedAt = old.CreatedAt
		}
	}
	if item.CreatedAt == "" {
		item.CreatedAt = now
	}
	if item.DefaultSave {
		if !item.Role.CanSave() {
			return Registration{}, fmt.Errorf("default save destination must have the save role")
		}
		for id, other := range r.items {
			if other.Category == item.Category && other.DefaultSave {
				other.DefaultSave = false
				r.items[id] = other
			}
		}
	}
	if item.Role.CanSave() && !r.hasDefaultSaveLocked(item.Category) {
		item.DefaultSave = true
	}
	r.items[item.ID] = item
	if err := r.saveLocked(); err != nil {
		return Registration{}, err
	}
	return item, nil
}

func (r *Registry) Remove(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[id]; !ok {
		return fmt.Errorf("content library %q is not registered", id)
	}
	delete(r.items, id)
	return r.saveLocked()
}

func (r *Registry) SetRole(id string, role Role) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.items[id]
	if !ok {
		return fmt.Errorf("content library %q is not registered", id)
	}
	item.Role = NormalizeRole(role)
	if item.DefaultSave && !item.Role.CanSave() {
		item.DefaultSave = false
	}
	r.items[id] = item
	return r.saveLocked()
}

func (r *Registry) SetDefaultSave(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.items[id]
	if !ok {
		return fmt.Errorf("content library %q is not registered", id)
	}
	if !item.Role.CanSave() {
		return fmt.Errorf("content library %q is not a save destination", id)
	}
	for otherID, other := range r.items {
		if other.Category == item.Category {
			other.DefaultSave = false
			r.items[otherID] = other
		}
	}
	item.DefaultSave = true
	r.items[id] = item
	return r.saveLocked()
}

func (r *Registry) DefaultSave(category Category) (Registration, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, item := range r.items {
		if item.Category == category && item.DefaultSave && item.Role.CanSave() {
			return item, true
		}
	}
	return Registration{}, false
}

func (r *Registry) ForPath(category Category, path string) (Registration, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, item := range r.items {
		if item.Category == category && IsUnderDir(path, item.Path) {
			return item, true
		}
	}
	return Registration{}, false
}

func (r *Registry) Paths(category Category, role Role) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []string{}
	for _, item := range r.items {
		if item.Category == category && item.Role&role == role {
			out = append(out, item.Path)
		}
	}
	return out
}

func (r *Registry) WriteNewFile(id, relative string, data []byte) (string, error) {
	return r.writeFile(id, relative, data, false)
}

// WriteFile safely replaces a file only in a library with the write role. A
// save-only destination can create new files but cannot replace an existing
// one; replacement must be an explicit operation on a writable root.
func (r *Registry) WriteFile(id, relative string, data []byte) (string, error) {
	return r.writeFile(id, relative, data, true)
}

func (r *Registry) writeFile(id, relative string, data []byte, replace bool) (string, error) {
	r.mu.RLock()
	item, ok := r.items[id]
	r.mu.RUnlock()
	if !ok {
		return "", fmt.Errorf("%w: %q", ErrNotRegistered, id)
	}
	if replace && !item.Role.CanWrite() {
		return "", fmt.Errorf("%w: %q has no write role", ErrReadOnly, id)
	}
	if !replace && !item.Role.CanWrite() && !item.Role.CanSave() {
		return "", fmt.Errorf("%w: %q", ErrReadOnly, id)
	}
	rel := filepath.Clean(relative)
	if filepath.IsAbs(rel) || rel == "." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." {
		return "", fmt.Errorf("%w: path must be relative to the library", ErrInvalidPath)
	}
	target := filepath.Join(item.Path, rel)
	if !IsUnderDir(target, item.Path) {
		return "", fmt.Errorf("%w: path escapes the library", ErrInvalidPath)
	}
	info, statErr := os.Lstat(target)
	exists := false
	switch {
	case statErr == nil:
		exists = true
		if !replace {
			return "", fmt.Errorf("%w: %s", ErrCollision, target)
		}
		if !info.Mode().IsRegular() {
			return "", fmt.Errorf("%w: %s", ErrUnsafeTarget, target)
		}
	case os.IsNotExist(statErr):
		// A new target is allowed for save-only roots.
	default:
		return "", fmt.Errorf("inspect content path: %w", statErr)
	}
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return "", fmt.Errorf("create content directory: %w", err)
	}
	tmp := target + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return "", fmt.Errorf("write content: %w", err)
	}
	if exists {
		backup := fmt.Sprintf("%s.bak-%d", target, time.Now().UnixNano())
		if err := r.renameContent(target, backup); err != nil {
			_ = os.Remove(tmp)
			return "", fmt.Errorf("stage existing content: %w", err)
		}
		if err := r.renameContent(tmp, target); err != nil {
			restoreErr := r.renameContent(backup, target)
			_ = os.Remove(tmp)
			if restoreErr != nil {
				return "", fmt.Errorf("commit content: %w; restore failed; original remains at %s: %v", err, backup, restoreErr)
			}
			return "", fmt.Errorf("commit content: %w", err)
		}
		_ = os.Remove(backup)
		return target, nil
	}
	if err := r.renameContent(tmp, target); err != nil {
		_ = os.Remove(tmp)
		return "", fmt.Errorf("commit content: %w", err)
	}
	return target, nil
}

func (r *Registry) renameContent(oldPath, newPath string) error {
	if r.renameFile != nil {
		return r.renameFile(oldPath, newPath)
	}
	return os.Rename(oldPath, newPath)
}

func (r *Registry) hasDefaultSaveLocked(category Category) bool {
	for _, item := range r.items {
		if item.Category == category && item.DefaultSave && item.Role.CanSave() {
			return true
		}
	}
	return false
}

func (r *Registry) saveLocked() error {
	items := make([]Registration, 0, len(r.items))
	for _, item := range r.items {
		items = append(items, item)
	}
	data, err := json.MarshalIndent(persistedRegistry{Version: 1, Items: items}, "", "  ")
	if err != nil {
		return fmt.Errorf("encode content libraries: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(r.path), 0755); err != nil {
		return fmt.Errorf("create content library config directory: %w", err)
	}
	tmp := r.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return fmt.Errorf("write content libraries: %w", err)
	}
	if err := os.Rename(tmp, r.path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("commit content libraries: %w", err)
	}
	return nil
}

func registrationID(category Category, path string) string {
	sum := sha256.Sum256([]byte(string(category) + "\x00" + strings.ToLower(filepath.Clean(path))))
	return hex.EncodeToString(sum[:12])
}

func viewFor(item Registration) View {
	_, err := os.Stat(item.Path)
	exists := err == nil
	return View{
		Registration: item,
		Exists:       exists,
		Readable:     item.Role.CanRead() && exists,
		Writable:     item.Role.CanWrite() && exists,
		SaveRoot:     item.Role.CanSave() && exists,
		Available:    exists,
	}
}

func sortViews(items []View) {
	for i := 1; i < len(items); i++ {
		for j := i; j > 0 && (items[j].Category < items[j-1].Category ||
			(items[j].Category == items[j-1].Category && items[j].Path < items[j-1].Path)); j-- {
			items[j], items[j-1] = items[j-1], items[j]
		}
	}
}

// IsUnderDir reports whether path is dir itself or lives beneath it. It is
// case-insensitive on Windows and checks a separator boundary, so a sibling
// named "Characters2" is not treated as part of "Characters".
func IsUnderDir(path, dir string) bool {
	if path == "" || dir == "" {
		return false
	}
	p := filepath.Clean(path)
	d := filepath.Clean(dir)
	if strings.EqualFold(p, d) {
		return true
	}
	if !strings.HasSuffix(d, string(filepath.Separator)) {
		d += string(filepath.Separator)
	}
	return strings.HasPrefix(strings.ToLower(p), strings.ToLower(d))
}
