package main

import (
	"fmt"

	"howl-chat/internal/backend/types"
	"howl-chat/internal/errors"
)

// ---------------------------------------------------------------------------
// Memory DTOs (frontend-visible types)
// ---------------------------------------------------------------------------

// MemoryFilter defines filtering options for memory queries.
type MemoryFilter struct {
	Query  string `json:"query,omitempty"`
	Scope  string `json:"scope,omitempty"`
	Status string `json:"status,omitempty"`
}

// MemoryEdit defines editable fields for memory updates.
type MemoryEdit struct {
	Title   string   `json:"title,omitempty"`
	Summary string   `json:"summary,omitempty"`
	Tags    []string `json:"tags,omitempty"`
	Pinned  bool     `json:"pinned,omitempty"`
}

// MemoryRecord represents a complete memory entry with metadata.
type MemoryRecord struct {
	ID              string                 `json:"id"`
	CreatedAt       string                 `json:"createdAt"`
	LastTriggeredAt string                 `json:"lastTriggeredAt,omitempty"`
	SourceMessages  []string               `json:"sourceMessages,omitempty"`
	RawText         string                 `json:"rawText,omitempty"`
	Summary         string                 `json:"summary,omitempty"`
	Title           string                 `json:"title"`
	Type            string                 `json:"type"`
	Scope           string                 `json:"scope"`
	ScopeID         string                 `json:"scopeId,omitempty"`
	Confidence      int                    `json:"confidence,omitempty"`
	Status          string                 `json:"status,omitempty"`
	Pinned          bool                   `json:"pinned,omitempty"`
	Tags            []string               `json:"tags,omitempty"`
	Author          string                 `json:"author,omitempty"`
	Trigger         string                 `json:"trigger,omitempty"`
	ScoreBreakdown  map[string]interface{} `json:"scoreBreakdown,omitempty"`
	ResolverVersion int                    `json:"resolverVersion,omitempty"`
	Version         int                    `json:"version,omitempty"`
}

// AutoCaptureConfig defines automatic memory capture settings.
type AutoCaptureConfig struct {
	Enabled       bool     `json:"enabled"`
	MinConfidence int      `json:"minConfidence"`
	MaxPerMessage int      `json:"maxPerMessage"`
	TriggerWords  []string `json:"triggerWords,omitempty"`
}

// ---------------------------------------------------------------------------
// Memory Wails Bindings
// ---------------------------------------------------------------------------

// GetMemories retrieves memories for a chat with optional filtering.
func (a *App) GetMemories(chatID string, filter MemoryFilter) ([]MemoryRecord, error) {
	if a.memoryService == nil {
		return nil, fmt.Errorf("memory service not initialized")
	}
	if chatID != "" {
		a.memoryService.SetChatID(chatID)
	}
	typesFilter := types.MemoryFilter{
		Query:  filter.Query,
		Scope:  filter.Scope,
		Status: filter.Status,
	}
	typesMemories, err := a.memoryService.GetMemories(typesFilter)
	if err != nil {
		a.logError(errors.MEMGetFail, errors.High, "memory.get", "app-memory.go:76", err)
		return nil, err
	}
	var memories []MemoryRecord
	for _, tm := range typesMemories {
		memories = append(memories, typesMemoryToRecord(tm))
	}
	return memories, nil
}

// GetMemoryDetail retrieves detailed information about a specific memory.
func (a *App) GetMemoryDetail(chatID string, memoryID string) (*MemoryRecord, error) {
	if a.memoryService == nil {
		return nil, fmt.Errorf("memory service not initialized")
	}
	if chatID != "" {
		a.memoryService.SetChatID(chatID)
	}
	tm, err := a.memoryService.GetMemoryDetail(memoryID)
	if err != nil {
		a.logError(errors.MEMGetFail, errors.High, "memory.detail", "app-memory.go:95", err)
		return nil, err
	}
	mem := typesMemoryToRecord(*tm)
	return &mem, nil
}

// EditMemory updates an existing memory with new values.
func (a *App) EditMemory(chatID string, memoryID string, edits MemoryEdit) (*MemoryRecord, error) {
	if a.memoryService == nil {
		return nil, fmt.Errorf("memory service not initialized")
	}
	if chatID != "" {
		a.memoryService.SetChatID(chatID)
	}
	typesEdits := types.MemoryEdit{
		Title:   edits.Title,
		Summary: edits.Summary,
		Tags:    edits.Tags,
		Pinned:  edits.Pinned,
	}
	tm, err := a.memoryService.EditMemory(memoryID, typesEdits)
	if err != nil {
		a.logError(errors.MEMEditFail, errors.High, "memory.edit", "app-memory.go:117", err)
		return nil, err
	}
	mem := typesMemoryToRecord(*tm)
	return &mem, nil
}

// AddMemory creates a new memory record in the chat session.
func (a *App) AddMemory(chatID string, record MemoryRecord) (*MemoryRecord, error) {
	if a.memoryService == nil {
		return nil, fmt.Errorf("memory service not initialized")
	}
	if chatID != "" {
		a.memoryService.SetChatID(chatID)
	}
	typesRecord := types.MemoryRecord{
		ID:         record.ID,
		Title:      record.Title,
		Summary:    record.Summary,
		Content:    record.RawText,
		Scope:      record.Scope,
		OwnerID:    record.ScopeID,
		Tags:       record.Tags,
		Pinned:     record.Pinned,
		Status:     record.Status,
		Confidence: float64(record.Confidence),
		Source:     record.Trigger,
		CreatedAt:  record.CreatedAt,
		UpdatedAt:  record.LastTriggeredAt,
	}
	tm, err := a.memoryService.AddMemory(typesRecord)
	if err != nil {
		a.logError(errors.MEMAddFail, errors.High, "memory.add", "app-memory.go:148", err)
		return nil, err
	}
	mem := typesMemoryToRecord(*tm)
	return &mem, nil
}

// PinMemory toggles the pinned status of a memory.
func (a *App) PinMemory(chatID string, memoryID string, pinned bool) error {
	if a.memoryService == nil {
		return fmt.Errorf("memory service not initialized")
	}
	if chatID != "" {
		a.memoryService.SetChatID(chatID)
	}
	err := a.memoryService.PinMemory(memoryID, pinned)
	if err != nil {
		a.logError(errors.MEMEditFail, errors.Medium, "memory.pin", "app-memory.go:164", err)
	}
	return err
}

// DeleteMemory removes a memory from the chat session.
func (a *App) DeleteMemory(chatID string, memoryID string) error {
	if a.memoryService == nil {
		return fmt.Errorf("memory service not initialized")
	}
	if chatID != "" {
		a.memoryService.SetChatID(chatID)
	}
	err := a.memoryService.DeleteMemory(memoryID)
	if err != nil {
		a.logError(errors.MEMDeleteFail, errors.Medium, "memory.delete", "app-memory.go:175", err)
	}
	return err
}

// ExportMemories exports memories with optional filtering.
func (a *App) ExportMemories(chatID string, filter MemoryFilter) ([]MemoryRecord, error) {
	if a.memoryService == nil {
		return nil, fmt.Errorf("memory service not initialized")
	}
	if chatID != "" {
		a.memoryService.SetChatID(chatID)
	}
	typesFilter := types.MemoryFilter{
		Query:  filter.Query,
		Scope:  filter.Scope,
		Status: filter.Status,
	}
	typesMemories, err := a.memoryService.ExportMemories(typesFilter)
	if err != nil {
		a.logError(errors.MEMExportFail, errors.High, "memory.export", "app-memory.go:191", err)
		return nil, err
	}
	var memories []MemoryRecord
	for _, tm := range typesMemories {
		memories = append(memories, typesMemoryToRecord(tm))
	}
	return memories, nil
}

// GetAutoCaptureConfig retrieves automatic memory capture settings.
func (a *App) GetAutoCaptureConfig(chatID string) (*AutoCaptureConfig, error) {
	if a.memoryService == nil {
		return nil, fmt.Errorf("memory service not initialized")
	}
	if chatID != "" {
		a.memoryService.SetChatID(chatID)
	}
	typesConfig, err := a.memoryService.GetAutoCaptureConfig()
	if err != nil {
		a.logErrorSilent(errors.MEMGetFail, errors.Medium, "memory.config.get", "app-memory.go:210", err)
		return nil, err
	}
	config := AutoCaptureConfig{
		Enabled:       typesConfig.Enabled,
		MinConfidence: int(typesConfig.MinConfidence),
		MaxPerMessage: typesConfig.MaxPerMessage,
		TriggerWords:  typesConfig.TriggerWords,
	}
	return &config, nil
}

// SetAutoCaptureConfig updates automatic memory capture settings.
func (a *App) SetAutoCaptureConfig(chatID string, config AutoCaptureConfig) error {
	if a.memoryService == nil {
		return fmt.Errorf("memory service not initialized")
	}
	if chatID != "" {
		a.memoryService.SetChatID(chatID)
	}
	typesConfig := types.AutoCaptureConfig{
		Enabled:       config.Enabled,
		MinConfidence: float64(config.MinConfidence),
		MaxPerMessage: config.MaxPerMessage,
		TriggerWords:  config.TriggerWords,
	}
	err := a.memoryService.SetAutoCaptureConfig(typesConfig)
	if err != nil {
		a.logErrorSilent(errors.MEMEditFail, errors.Medium, "memory.config.set", "app-memory.go:237", err)
	}
	return err
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// typesMemoryToRecord converts an internal types.MemoryRecord to the frontend DTO.
func typesMemoryToRecord(tm types.MemoryRecord) MemoryRecord {
	return MemoryRecord{
		ID:              tm.ID,
		CreatedAt:       tm.CreatedAt,
		LastTriggeredAt: tm.UpdatedAt,
		SourceMessages:  []string{},
		RawText:         tm.Content,
		Summary:         tm.Summary,
		Title:           tm.Title,
		Type:            "automated",
		Scope:           tm.Scope,
		ScopeID:         tm.OwnerID,
		Confidence:      int(tm.Confidence),
		Status:          tm.Status,
		Pinned:          tm.Pinned,
		Tags:            tm.Tags,
		Author:          "system",
		Trigger:         tm.Source,
		ScoreBreakdown:  map[string]interface{}{},
		ResolverVersion: 1,
		Version:         1,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
