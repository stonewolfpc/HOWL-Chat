package main

import (
	"fmt"
	"strings"

	"howl-chat/internal/contentimport"
)

// PreviewContentImport performs the observational import preflight. It is the
// only route that returns a plan hash; it never writes.
func (a *App) PreviewContentImport(req contentimport.Request) (contentimport.Plan, error) {
	service, err := a.importService(req.Category)
	if err != nil {
		return contentimport.Plan{}, err
	}
	return service.Preview(req)
}

// ExecuteContentImport re-derives the plan and requires the reviewed hash. The
// renderer cannot submit an edited write set.
func (a *App) ExecuteContentImport(req contentimport.Request, planHash string) (contentimport.Result, error) {
	service, err := a.importService(req.Category)
	if err != nil {
		return contentimport.Result{}, err
	}
	return service.Execute(req, planHash)
}

// ResumeContentImport performs explicit crash recovery. It is never called at
// startup, on chat load, or on editor open.
func (a *App) ResumeContentImport(category, operationID string) (contentimport.Result, error) {
	service, err := a.importService(category)
	if err != nil {
		return contentimport.Result{}, err
	}
	return service.Resume(category, operationID)
}

// RollbackContentImport verifies all journaled items before removing any
// journal-created content.
func (a *App) RollbackContentImport(category, operationID string) (contentimport.RollbackResult, error) {
	service, err := a.importService(category)
	if err != nil {
		return contentimport.RollbackResult{}, err
	}
	return service.Rollback(category, operationID)
}

// ListContentImportJournal returns durable operations without running recovery.
func (a *App) ListContentImportJournal(category string) ([]contentimport.Operation, error) {
	service, err := a.importService(category)
	if err != nil {
		return nil, err
	}
	return service.Operations(), nil
}

func (a *App) importService(category string) (*contentimport.Service, error) {
	switch strings.ToLower(strings.TrimSpace(category)) {
	case "world", "worlds":
		if a.registries.Worlds() == nil {
			return nil, fmt.Errorf("world import service is not initialised")
		}
		return a.registries.Worlds(), nil
	case "scenario", "scenarios":
		if a.registries.Scenarios() == nil {
			return nil, fmt.Errorf("scenario import service is not initialised")
		}
		return a.registries.Scenarios(), nil
	default:
		return nil, fmt.Errorf("unsupported import category %q", category)
	}
}
