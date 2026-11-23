// Package registries is the composition owner for the application's registry
// authorities.
//
// It does not add state or index anything: every registry here is already the
// single authority for its own truth (registered model directories in
// internal/image/registry, content roots in internal/library, import journals in
// internal/contentimport, canonical lorebooks in internal/lorebookstore, model
// profiles in internal/modelintel). This type constructs them in one place and
// hands out one handle, so App holds no registry fields of its own.
//
// Failure policy matches boot: a corrupt or unreadable store leaves that one
// component nil (its bindings then report "not initialised") without discarding
// the others. Content import services depend on the content-library registry, so
// they are absent when that registry could not be loaded; the inventory,
// lorebook store, intelligence cache and model registry are independent.
package registries

import (
	"fmt"

	"howl-chat/internal/contentimport"
	"howl-chat/internal/image/registry"
	"howl-chat/internal/library"
	"howl-chat/internal/lorebookstore"
	"howl-chat/internal/modelintel"
)

// Set holds the registry authorities. A nil *Set is an uninitialised owner:
// every accessor returns nil, which is the same state the bindings already
// report as "not initialised".
type Set struct {
	models    *registry.Registry
	libraries *library.Registry
	inventory *library.InventoryStore
	worlds    *contentimport.Service
	scenarios *contentimport.Service
	lorebooks *lorebookstore.Store
	intel     *modelintel.Cache
}

// Parts is the explicit composition used by tests and by callers that already
// own one of the components. Unset parts stay nil.
type Parts struct {
	Models    *registry.Registry
	Libraries *library.Registry
	Inventory *library.InventoryStore
	Worlds    *contentimport.Service
	Scenarios *contentimport.Service
	Lorebooks *lorebookstore.Store
	Intel     *modelintel.Cache
}

// Compose builds a Set from explicit parts without touching disk.
func Compose(parts Parts) *Set {
	return &Set{
		models:    parts.Models,
		libraries: parts.Libraries,
		inventory: parts.Inventory,
		worlds:    parts.Worlds,
		scenarios: parts.Scenarios,
		lorebooks: parts.Lorebooks,
		intel:     parts.Intel,
	}
}

// New loads every registry authority rooted at configDir. The returned errors
// are the per-component failures, in boot order; the Set is always usable and
// its accessors report exactly which components failed.
func New(configDir string) (*Set, []error) {
	var errs []error
	set := &Set{}

	libs, err := library.New(configDir)
	if err != nil {
		errs = append(errs, fmt.Errorf("load content libraries: %w", err))
	} else {
		set.libraries = libs
		if service, err := contentimport.New(libs, library.Worlds, configDir); err != nil {
			errs = append(errs, fmt.Errorf("initialise world import: %w", err))
		} else {
			set.worlds = service
		}
		if service, err := contentimport.New(libs, library.Scenarios, configDir); err != nil {
			errs = append(errs, fmt.Errorf("initialise scenario import: %w", err))
		} else {
			set.scenarios = service
		}
	}

	// The inventory ignores an unreadable cache (it is a snapshot, not canon)
	// and is always present so a stale snapshot can be replaced.
	set.inventory = library.NewInventoryStore(configDir)
	// The lorebook store is constructed even without a registry: its List
	// reports no roots rather than inventing a second source of truth.
	set.lorebooks = lorebookstore.New(libs, configDir)

	if intel, err := modelintel.NewCache(configDir); err != nil {
		errs = append(errs, fmt.Errorf("load model intelligence cache: %w", err))
	} else {
		set.intel = intel
	}

	if models, err := registry.New(configDir); err != nil {
		errs = append(errs, fmt.Errorf("load model registry: %w", err))
	} else {
		set.models = models
	}

	return set, errs
}

// Models returns the registered model directories/types authority.
func (s *Set) Models() *registry.Registry {
	if s == nil {
		return nil
	}
	return s.models
}

// Libraries returns the content-root registry.
func (s *Set) Libraries() *library.Registry {
	if s == nil {
		return nil
	}
	return s.libraries
}

// Inventory returns the last-known content snapshot store.
func (s *Set) Inventory() *library.InventoryStore {
	if s == nil {
		return nil
	}
	return s.inventory
}

// Worlds returns the world import service.
func (s *Set) Worlds() *contentimport.Service {
	if s == nil {
		return nil
	}
	return s.worlds
}

// Scenarios returns the scenario import service.
func (s *Set) Scenarios() *contentimport.Service {
	if s == nil {
		return nil
	}
	return s.scenarios
}

// Lorebooks returns the canonical lorebook store.
func (s *Set) Lorebooks() *lorebookstore.Store {
	if s == nil {
		return nil
	}
	return s.lorebooks
}

// Intel returns the shared model-profile cache.
func (s *Set) Intel() *modelintel.Cache {
	if s == nil {
		return nil
	}
	return s.intel
}
