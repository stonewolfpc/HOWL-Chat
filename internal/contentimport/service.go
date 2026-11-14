package contentimport

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"howl-chat/internal/library"
	"howl-chat/internal/worldstore"
)

// Service owns import orchestration for one world/scenario category. Canonical
// writes still belong exclusively to worldstore.
type Service struct {
	reg       *library.Registry
	store     *worldstore.Store
	category  library.Category
	configDir string
	journalMu sync.Mutex
	markerMu  sync.Mutex
	faultMu   sync.Mutex
	fault     func(string) error
}

// New builds an import service for worlds or scenarios.
func New(reg *library.Registry, category library.Category, configDir string) (*Service, error) {
	store, err := worldstore.New(reg, category)
	if err != nil {
		return nil, err
	}
	return &Service{reg: reg, store: store, category: category, configDir: configDir}, nil
}

// SetFaultHook installs a deterministic test hook. It is never set by the
// production binary.
func (s *Service) SetFaultHook(hook func(string) error) {
	s.faultMu.Lock()
	s.fault = hook
	s.faultMu.Unlock()
}

func (s *Service) inject(boundary string) error {
	s.faultMu.Lock()
	hook := s.fault
	s.faultMu.Unlock()
	if hook == nil {
		return nil
	}
	return hook(boundary)
}

// Preview is observational. It reads sources and destination state, validates
// and hashes, and creates no journal, staging directory, or marker.
func (s *Service) Preview(req Request) (Plan, error) {
	built, err := s.buildPlan(req)
	if err != nil {
		return Plan{}, err
	}
	return built.Plan, nil
}

func (s *Service) buildPlan(req Request) (builtPlan, error) {
	category := categoryOf(req.Category)
	if category != library.Worlds && category != library.Scenarios {
		return builtPlan{}, fmt.Errorf("category must be worlds or scenarios")
	}
	if category != s.category {
		return builtPlan{}, fmt.Errorf("import service category %q cannot handle %q", s.category, category)
	}
	destinationID := strings.TrimSpace(req.DestinationLibraryID)
	plan := Plan{
		Category:             string(category),
		DestinationLibraryID: destinationID,
		DestinationStatus:    "destination_unavailable",
		ConversionVersion:    ConversionVersion,
		Items:                []PlanItem{},
		Errors:               []string{},
	}
	var destinationView library.View
	for _, view := range s.reg.List(category) {
		if view.ID == destinationID {
			destinationView = view
			break
		}
	}
	if destinationView.ID == "" {
		plan.Errors = append(plan.Errors, "The destination library is not registered.")
	} else {
		plan.DestinationRoot = destinationView.Path
		switch {
		case !destinationView.Exists || !destinationView.Readable:
			plan.DestinationStatus = "destination_unavailable"
		case !destinationView.Role.CanSave() || !destinationView.Role.CanWrite():
			plan.DestinationStatus = "destination_not_rollback_safe"
		default:
			plan.DestinationStatus = "available"
		}
	}

	browser, err := s.loadBrowserSources(req.BrowserRecords)
	if err != nil {
		return builtPlan{}, err
	}
	legacy, legacyErrors, err := s.loadLegacySources(req.LegacyRoots, plan.DestinationRoot)
	if err != nil {
		return builtPlan{}, err
	}
	plan.Errors = append(plan.Errors, legacyErrors...)
	sources := append(browser, legacy...)
	resolutions := map[string]Resolution{}
	for _, resolution := range req.Resolutions {
		key := strings.TrimSpace(resolution.SourceKey)
		if key == "" {
			plan.Errors = append(plan.Errors, "A resolution had an empty source key.")
			continue
		}
		if _, exists := resolutions[key]; exists {
			plan.Errors = append(plan.Errors, "Duplicate resolution for "+key)
			continue
		}
		resolutions[key] = resolution
	}

	idCounts := map[string][]int{}
	nameCounts := map[string][]int{}
	for i, source := range sources {
		if source.id != "" {
			idCounts[source.id] = append(idCounts[source.id], i)
		}
		if nameKey := sourceNameKey(source.name); nameKey != "" {
			nameCounts[nameKey] = append(nameCounts[nameKey], i)
		}
	}
	for i := range sources {
		if len(idCounts[sources[i].id]) > 1 && sources[i].id != "" {
			sources[i].blockReason = "duplicate_identity"
			sources[i].blockMessage = "More than one supplied source uses this identity."
		}
		if nameKey := sourceNameKey(sources[i].name); nameKey != "" && len(nameCounts[nameKey]) > 1 {
			if sources[i].blockReason == "" {
				sources[i].blockReason = "duplicate_name"
				sources[i].blockMessage = "More than one supplied source uses this display name."
			}
		}
	}

	built := builtPlan{Plan: plan, prepared: make([]preparedItem, len(sources))}
	for i, source := range sources {
		planItem := PlanItem{
			SourceKey:  source.key,
			SourceKind: source.kind,
			SourceID:   source.id,
			SourceName: source.name,
			SourceRoot: source.root,
			Folder:     source.folder,
			SourceHash: sourceFingerprint(source),
			Action:     ActionQuarantine,
		}
		resolution, hasResolution := resolutions[source.key]
		evidence := ""
		if hasResolution {
			evidence = strings.TrimSpace(resolution.Evidence)
		}
		planItem.Evidence = evidence
		if source.blockReason != "" {
			planItem.Reason = source.blockReason
			planItem.Message = source.blockMessage
		} else if plan.DestinationStatus != "available" {
			planItem.Reason = plan.DestinationStatus
			planItem.Message = "The destination cannot be used for a reversible import."
		} else {
			destinationItemID := source.id
			switch {
			case hasResolution && strings.TrimSpace(resolution.DestinationID) != "":
				destinationItemID = strings.TrimSpace(resolution.DestinationID)
				if !worldstore.ValidID(destinationItemID) {
					planItem.Reason = "invalid_document_id"
					planItem.Message = "The explicit destination identity is not safe for a canonical item."
					destinationItemID = ""
				}
			case hasResolution && resolution.NewID:
				destinationItemID = s.store.ProposeID()
				evidence = "explicit new identity"
			case source.id == "":
				planItem.Reason = "missing_document_id"
				planItem.Message = "The source has no stable identity; an explicit mapping is required."
			case !worldstore.ValidID(source.id):
				planItem.Reason = "invalid_document_id"
				planItem.Message = "The source identity is not safe for a canonical item; an explicit new identity is required."
			}
			if planItem.Reason == "" {
				prepared, reason, message := s.prepareItem(source, destinationItemID, evidence)
				if reason != "" {
					planItem.Reason = reason
					planItem.Message = message
				} else {
					built.prepared[i] = prepared
					planItem.ProposedID = destinationItemID
					planItem.DocumentHash = hashJSON(prepared.document)
					planItem.LoreHashes = loreHashes(prepared.lore)
					planItem.MediaHash = hashBytes(prepared.media)
					planItem.References = append([]ReferenceIssue(nil), prepared.references...)
					plan.ReferenceIssues += len(prepared.references)
					planItem.Action, planItem.Reason, planItem.Message = s.destinationDisposition(destinationItemID, prepared)
					if planItem.Action == ActionReady {
						planItem.Writes = writeNames(s, prepared)
					}
				}
			}
		}
		switch planItem.Action {
		case ActionReady:
			plan.Ready++
			plan.Writes++
		case ActionAdopt:
			plan.Adoptable++
		default:
			plan.Quarantined++
			if planItem.Reason == "unsupported_source_file" || planItem.Reason == "unsupported_media_manifest" {
				plan.Unsupported++
			}
			if planItem.Reason == "identity_conflict" || planItem.Reason == "lore_conflict" || planItem.Reason == "media_conflict" {
				plan.Conflicted++
			}
		}
		plan.Items = append(plan.Items, planItem)
	}
	built.Plan = plan
	built.Plan.Hash = planHash(plan)
	return built, nil
}

func (s *Service) destinationDisposition(id string, prepared preparedItem) (string, string, string) {
	search, err := s.store.Lookup(id)
	if err != nil {
		return ActionQuarantine, "destination_unavailable", err.Error()
	}
	if !search.Complete {
		return ActionQuarantine, "destination_unavailable", "Not every registered root could be searched; absence and identity cannot be proven."
	}
	if !search.Found || search.Item == nil {
		return ActionReady, "", ""
	}
	if search.AuthorityIndeterminate || search.Item.AuthorityIndeterminate {
		return ActionQuarantine, "authority_indeterminate", "Another root is unavailable or holds a duplicate; no import is allowed."
	}
	if hashJSON(search.Item.Document) != hashJSON(prepared.document) {
		return ActionQuarantine, "identity_conflict", "A canonical item with this identity has different content; it will not be overwritten."
	}
	gotMedia, mediaExists, err := s.store.LoadMedia(id)
	if err != nil {
		return ActionQuarantine, "media_conflict", "The destination item exists but its media manifest could not be verified."
	}
	wantMedia := prepared.media
	if len(wantMedia) == 0 && mediaExists {
		return ActionQuarantine, "media_conflict", "The source has no media manifest but the canonical item has one; the item will not be adopted."
	}
	if len(wantMedia) > 0 && !mediaExists {
		return ActionQuarantine, "media_conflict", "The source has a media manifest but the canonical item has none; the item will not be adopted."
	}
	if len(wantMedia) > 0 && !bytes.Equal(wantMedia, gotMedia) {
		return ActionQuarantine, "media_conflict", "The source and canonical media manifests differ; the item will not be adopted."
	}
	for _, name := range s.loreNames() {
		want := prepared.lore[name]
		got, err := s.store.LoadLore(id, name)
		if err != nil {
			return ActionQuarantine, "lore_conflict", "The destination item exists but its lore could not be verified."
		}
		if semanticLoreHash(want) != hashJSON(map[string]interface{}{"entries": got.Entries}) {
			return ActionQuarantine, "lore_conflict", "The destination item has different lore; manual and automated lore are never merged."
		}
	}
	_ = prepared
	return ActionAdopt, "", "The canonical item is identical; it will be adopted without a write."
}

// Execute re-derives the plan and refuses a stale or edited plan before any
// journal or staging write. Atomicity is per item, never transaction-wide.
func (s *Service) Execute(req Request, planHash string) (Result, error) {
	built, err := s.buildPlan(req)
	if err != nil {
		return Result{}, err
	}
	if strings.TrimSpace(planHash) == "" || planHash != built.Hash {
		return Result{Status: "plan_stale", Plan: built.Plan, Message: "The sources or destination changed after preview; no write occurred."}, nil
	}
	operationID := strings.TrimSpace(req.OperationID)
	if operationID == "" {
		operationID = newOperationID()
	}
	if !worldstore.ValidID(operationID) {
		return Result{}, fmt.Errorf("operation identity is invalid")
	}
	if existing, ok := s.loadOperation(operationID); ok {
		if err := s.verifyMarker(existing); err != nil {
			return resultFromOperation(existing, OperationIntegrityFailure, err.Error()), nil
		}
		if existing.PlanHash != built.Hash {
			return Result{OperationID: operationID, Status: "plan_mismatch", Plan: built.Plan, Message: "The operation belongs to a different plan; no write occurred."}, nil
		}
		if operationTerminal(existing.Status) {
			return resultFromOperation(existing, "already_imported", "This operation is already terminal."), nil
		}
		return resultFromOperation(existing, "operation_active", "The operation is already active; use explicit resume."), nil
	}

	operation := Operation{
		Version:              1,
		Category:             string(s.category),
		OperationID:          operationID,
		PlanHash:             built.Hash,
		DestinationLibraryID: built.DestinationLibraryID,
		DestinationRoot:      built.DestinationRoot,
		Status:               OperationExecuting,
		StartedAt:            nowRFC3339(),
		UpdatedAt:            nowRFC3339(),
		Items:                make([]JournalItem, 0, len(built.Items)),
	}
	for _, item := range built.Items {
		state := StateQuarantined
		if item.Action == ActionReady {
			state = StatePlanned
		} else if item.Action == ActionAdopt {
			state = StateAdopted
		}
		operation.Items = append(operation.Items, JournalItem{
			SourceKey:  item.SourceKey,
			SourceKind: item.SourceKind,
			SourceID:   item.SourceID,
			SourceName: item.SourceName,
			SourceRoot: item.SourceRoot,
			Folder:     item.Folder,
			ItemID:     item.ProposedID,
			LibraryID:  built.DestinationLibraryID,
			Root:       built.DestinationRoot,
			PlanHash:   built.Hash,
			State:      state,
			Reason:     item.Reason,
			Message:    item.Message,
			Evidence:   item.Evidence,
			SourceHash: item.SourceHash,
		})
	}
	if err := s.writeOperation(operation); err != nil {
		return Result{}, err
	}
	if err := s.inject("after_journal"); err != nil {
		return resultFromOperation(operation, "resume_required", err.Error()), nil
	}

	for i := range operation.Items {
		item := &operation.Items[i]
		if item.State != StatePlanned {
			continue
		}
		prepared := built.prepared[i]
		if prepared.document == nil {
			item.State = StateFailed
			item.Reason = "write_failed"
			item.Message = "The re-derived plan had no prepared item."
			operation.UpdatedAt = nowRFC3339()
			_ = s.writeOperation(operation)
			continue
		}
		staged, err := s.store.StageImport(worldstore.ImportItem{
			ID:                   item.ItemID,
			DestinationLibraryID: built.DestinationLibraryID,
			Document:             prepared.document,
			Origin:               prepared.origin,
			CreatedAt:            prepared.createdAt,
			Lore:                 prepared.lore,
			Media:                prepared.media,
			PlanHash:             built.Hash,
		}, operationID)
		if err != nil {
			item.State = StateFailed
			item.Reason = "write_failed"
			if errors.Is(err, worldstore.ErrImportStagingExists) {
				item.State = StateAmbiguous
				item.Reason = "ambiguous_state"
			}
			item.Message = err.Error()
			operation.UpdatedAt = nowRFC3339()
			_ = s.writeOperation(operation)
			continue
		}
		if err := s.inject("after_stage"); err != nil {
			return resultFromOperation(operation, "resume_required", err.Error()), nil
		}
		item.State = StatePrepared
		item.StagingPath = staged.StagingPath
		item.OwnerPath = staged.OwnerPath
		item.FinalPath = staged.FinalPath
		item.Revision = staged.Revision
		item.DocumentHash = staged.DocumentHash
		item.LoreHashes = staged.LoreHashes
		item.MediaHash = staged.MediaHash
		operation.UpdatedAt = nowRFC3339()
		if err := s.writeOperation(operation); err != nil {
			return resultFromOperation(operation, "resume_required", err.Error()), err
		}
		if err := s.inject("after_prepared"); err != nil {
			return resultFromOperation(operation, "resume_required", err.Error()), nil
		}

		commit, commitErr := s.store.CommitImport(staged)
		if commitErr != nil {
			item.State = StateFailed
			item.Reason = "write_failed"
			item.Message = commitErr.Error()
			if abortErr := s.store.AbortImport(staged); abortErr != nil {
				item.State = StateAmbiguous
				item.Reason = "cleanup_failed"
				item.Message += "; cleanup failed: " + abortErr.Error()
			}
			operation.UpdatedAt = nowRFC3339()
			_ = s.writeOperation(operation)
			continue
		}
		if commit.Status != "saved" {
			item.State = StateQuarantined
			item.Reason = commitReason(commit.Status)
			item.Message = commit.Message
			if abortErr := s.store.AbortImport(staged); abortErr != nil {
				item.State = StateAmbiguous
				item.Reason = "cleanup_failed"
				item.Message += "; cleanup failed: " + abortErr.Error()
			}
			operation.UpdatedAt = nowRFC3339()
			_ = s.writeOperation(operation)
			continue
		}
		if err := s.inject("after_rename"); err != nil {
			return resultFromOperation(operation, "resume_required", err.Error()), nil
		}
		item.State = StateCommitted
		item.FinalPath = commit.Path
		item.Revision = commit.Revision
		operation.UpdatedAt = nowRFC3339()
		if err := s.writeOperation(operation); err != nil {
			return resultFromOperation(operation, "resume_required", err.Error()), err
		}
		if err := s.inject("after_committed"); err != nil {
			return resultFromOperation(operation, "resume_required", err.Error()), nil
		}
	}
	s.finalize(&operation)
	if err := s.writeOperation(operation); err != nil {
		return Result{}, err
	}
	if err := s.syncMarker(operation); err != nil {
		return resultFromOperation(operation, OperationIntegrityFailure, err.Error()), nil
	}
	return resultFromOperation(operation, operation.Status, operation.Message), nil
}

// Resume is the only recovery entry point. It verifies ownership evidence and
// hashes before repairing a prepared item; it is never called at startup.
func (s *Service) Resume(category, operationID string) (Result, error) {
	if categoryOf(category) != s.category {
		return Result{}, fmt.Errorf("import service category %q cannot resume %q", s.category, category)
	}
	operation, ok := s.loadOperation(strings.TrimSpace(operationID))
	if !ok {
		return Result{}, fmt.Errorf("import operation %q was not found", operationID)
	}
	if err := s.verifyMarker(operation); err != nil {
		return resultFromOperation(operation, OperationIntegrityFailure, err.Error()), nil
	}
	if operation.Status == OperationRolledBack || operation.Status == OperationRollbackPartial {
		return resultFromOperation(operation, operation.Status, operation.Message), nil
	}
	for i := range operation.Items {
		item := &operation.Items[i]
		switch item.State {
		case StateCommitted:
			ok, err := s.verifyJournalItem(*item)
			if err != nil || !ok {
				item.State = StateAmbiguous
				item.Reason = "ambiguous_state"
				item.Message = "The committed item no longer matches its journaled revision and hashes."
			}
		case StateAdopted, StateRolledBack, StateCleanupDone:
			continue
		case StateQuarantined, StateFailed, StateCleanupFailed:
			if s.finalExists(*item) {
				item.State = StateAmbiguous
				item.Reason = "ambiguous_state"
				item.Message = "A final item exists for a non-committed journal entry; no automatic cleanup occurred."
				continue
			}
			if !s.stagingPathExists(operation.OperationID, *item) {
				if s.stagingOwnerExists(*item) {
					item.State = StateCleanupFailed
					item.Reason = "cleanup_failed"
					item.Message = "An owned staging marker remains without an owned staging item; inspect it before cleanup."
				}
				continue
			}
			staged, err := s.store.LoadStaging(operation.OperationID, item.ItemID, item.LibraryID, item.PlanHash)
			if err != nil {
				item.State = StateCleanupFailed
				item.Reason = "cleanup_failed"
				item.Message = "Staging ownership cannot be verified for cleanup: " + err.Error()
				continue
			}
			if err := s.store.AbortImport(staged); err != nil {
				item.State = StateCleanupFailed
				item.Reason = "cleanup_failed"
				item.Message = "Owned staging could not be removed: " + err.Error()
				continue
			}
			item.State = StateCleanupDone
			item.StagingPath = ""
			item.OwnerPath = ""
		case StatePrepared:
			finalExists := s.finalExists(*item)
			if finalExists {
				ok, err := s.verifyJournalItem(*item)
				if err != nil || !ok {
					item.State = StateAmbiguous
					item.Reason = "ambiguous_state"
					item.Message = "A final item exists but its journaled hashes cannot be verified."
					continue
				}
				item.State = StateCommitted
				continue
			}
			staged, err := s.store.LoadStaging(operation.OperationID, item.ItemID, item.LibraryID, item.PlanHash)
			if err != nil {
				item.State = StateAmbiguous
				item.Reason = "ambiguous_state"
				item.Message = "Staging ownership or hashes cannot be proven: " + err.Error()
				continue
			}
			if err := s.commitRecovered(item, staged); err != nil {
				item.State = StateAmbiguous
				item.Reason = "ambiguous_state"
				item.Message = err.Error()
			}
		case StatePlanned:
			if s.finalExists(*item) {
				item.State = StateAmbiguous
				item.Reason = "ambiguous_state"
				item.Message = "A final item exists without a prepared journal entry."
				continue
			}
			staged, err := s.store.LoadStaging(operation.OperationID, item.ItemID, item.LibraryID, item.PlanHash)
			if err != nil {
				if s.stagingPathExists(operation.OperationID, *item) {
					item.State = StateAmbiguous
					item.Reason = "ambiguous_state"
					item.Message = "Staging exists but its ownership or hashes are incomplete; no canonical write occurred."
				} else {
					item.State = StateFailed
					item.Reason = "write_failed"
					item.Message = "No complete owned staging is available; no canonical write occurred."
				}
				continue
			}
			item.State = StatePrepared
			item.StagingPath = staged.StagingPath
			item.OwnerPath = staged.OwnerPath
			item.FinalPath = staged.FinalPath
			item.Revision = staged.Revision
			item.DocumentHash = staged.DocumentHash
			item.LoreHashes = staged.LoreHashes
			item.MediaHash = staged.MediaHash
			if err := s.commitRecovered(item, staged); err != nil {
				item.State = StateAmbiguous
				item.Reason = "ambiguous_state"
				item.Message = err.Error()
			}
		case StateAmbiguous:
			continue
		default:
			item.State = StateAmbiguous
			item.Reason = "ambiguous_state"
			item.Message = "The journal contains an unknown item state."
		}
	}
	s.finalize(&operation)
	if err := s.writeOperation(operation); err != nil {
		return Result{}, err
	}
	if err := s.syncMarker(operation); err != nil {
		return resultFromOperation(operation, OperationIntegrityFailure, err.Error()), nil
	}
	return resultFromOperation(operation, operation.Status, operation.Message), nil
}

// Rollback verifies every committed item first, then removes only unchanged,
// journal-created items. Adopted items or any ambiguity stop all deletion.
func (s *Service) Rollback(category, operationID string) (RollbackResult, error) {
	if categoryOf(category) != s.category {
		return RollbackResult{}, fmt.Errorf("import service category %q cannot roll back %q", s.category, category)
	}
	operation, ok := s.loadOperation(strings.TrimSpace(operationID))
	if !ok {
		return RollbackResult{}, fmt.Errorf("import operation %q was not found", operationID)
	}
	if err := s.verifyMarker(operation); err != nil {
		return RollbackResult{OperationID: operation.OperationID, Status: OperationIntegrityFailure, Message: err.Error()}, nil
	}
	res := RollbackResult{OperationID: operation.OperationID, Items: []RollbackItem{}}
	if operation.Status == OperationRolledBack {
		res.Status = OperationRolledBack
		res.Message = "The operation was already rolled back."
		return res, nil
	}
	if operation.Status == OperationExecuting || operation.Status == OperationAmbiguous {
		res.Status = "rollback_conflict"
		res.Message = "The operation is not in a terminal, verifiable state; run explicit resume or inspect the journal first."
		return res, nil
	}
	for _, item := range operation.Items {
		switch item.State {
		case StateAdopted:
			res.Status = "rollback_conflict"
			res.Message = "An adopted item is present; rollback refuses to delete or reinterpret adopted content."
			return s.finishRollbackConflict(operation, res)
		case StateAmbiguous:
			res.Status = "rollback_conflict"
			res.Message = "The operation contains an ambiguous item; no file was removed."
			return s.finishRollbackConflict(operation, res)
		case StatePrepared:
			if s.finalExists(item) {
				res.Status = "rollback_conflict"
				res.Message = "A prepared item already has a final path; no file was removed."
				return s.finishRollbackConflict(operation, res)
			}
		case StateFailed, StateQuarantined, StateCleanupFailed:
			if s.finalExists(item) {
				res.Status = "rollback_conflict"
				res.Message = "A final item exists for a non-committed journal entry; no file was removed."
				return s.finishRollbackConflict(operation, res)
			}
			if s.stagingPathExists(operation.OperationID, item) {
				staged, err := s.store.LoadStaging(operation.OperationID, item.ItemID, item.LibraryID, item.PlanHash)
				if err != nil {
					res.Status = "rollback_conflict"
					res.Message = "Owned staging cannot be verified for cleanup: " + err.Error()
					return s.finishRollbackConflict(operation, res)
				}
				_ = staged
			}
		case StatePlanned:
			if s.finalExists(item) || s.stagingPathExists(operation.OperationID, item) {
				res.Status = "rollback_conflict"
				res.Message = "A planned item has staging or final-path evidence; run explicit resume before rollback."
				return s.finishRollbackConflict(operation, res)
			}
		case StateCommitted:
			if s.stagingPathExists(operation.OperationID, item) || s.stagingOwnerExists(item) {
				res.Status = "rollback_conflict"
				res.Message = "Owned staging or marker evidence remains after commit; no file was removed."
				return s.finishRollbackConflict(operation, res)
			}
			ok, err := s.verifyJournalItem(item)
			if err != nil || !ok {
				res.Status = "rollback_conflict"
				res.Items = append(res.Items, RollbackItem{ItemID: item.ItemID, State: StateCommitted, Reason: "rollback_conflict", Message: "The item changed after import; it was left in place."})
				res.Message = "At least one committed item no longer matches its journaled hashes; no file was removed."
				return s.finishRollbackConflict(operation, res)
			}
		}
	}
	for _, item := range operation.Items {
		switch item.State {
		case StateCommitted:
			removed, err := s.store.RemoveImported(item.ItemID, item.Revision, s.journalHashes(item))
			if err != nil || removed.Status != "rolled_back" {
				res.Items = append(res.Items, RollbackItem{ItemID: item.ItemID, State: StateCommitted, Reason: "rollback_conflict", Message: removalMessage(removed, err)})
				res.Status = "rollback_conflict"
				continue
			}
			item.State = StateRolledBack
			res.Items = append(res.Items, RollbackItem{ItemID: item.ItemID, State: StateRolledBack})
		case StatePrepared, StateFailed, StateQuarantined, StateCleanupFailed:
			if item.State != StatePrepared && !s.stagingPathExists(operation.OperationID, item) {
				if s.stagingOwnerExists(item) {
					res.Items = append(res.Items, RollbackItem{ItemID: item.ItemID, State: StateCleanupFailed, Reason: "cleanup_failed", Message: "An owned staging marker remains without its staging item."})
					res.Status = "rollback_conflict"
					continue
				}
				item.State = StateCleanupDone
				res.Items = append(res.Items, RollbackItem{ItemID: item.ItemID, State: StateCleanupDone})
				continue
			}
			staged, err := s.store.LoadStaging(operation.OperationID, item.ItemID, item.LibraryID, item.PlanHash)
			if err != nil {
				res.Items = append(res.Items, RollbackItem{ItemID: item.ItemID, State: item.State, Reason: "ambiguous_state", Message: err.Error()})
				res.Status = "rollback_conflict"
				continue
			}
			if err := s.store.AbortImport(staged); err != nil {
				item.State = StateCleanupFailed
				res.Items = append(res.Items, RollbackItem{ItemID: item.ItemID, State: StateCleanupFailed, Reason: "cleanup_failed", Message: err.Error()})
				res.Status = "rollback_conflict"
				continue
			}
			if item.State == StatePrepared {
				item.State = StateRolledBack
			} else {
				item.State = StateCleanupDone
			}
			item.StagingPath = ""
			item.OwnerPath = ""
			res.Items = append(res.Items, RollbackItem{ItemID: item.ItemID, State: item.State})
		}
	}
	if res.Status == "" {
		operation.Status = OperationRolledBack
		operation.Message = "Journal-created items were removed; adopted items and sources were untouched."
		res.Status = OperationRolledBack
		res.Message = operation.Message
	} else if res.Status == "rollback_conflict" {
		operation.Status = OperationRollbackPartial
		operation.Message = "Rollback stopped or was partial; inspect the reported items."
	}
	operation.UpdatedAt = nowRFC3339()
	if err := s.writeOperation(operation); err != nil {
		return res, err
	}
	if err := s.removeMarker(operation.OperationID); err != nil {
		return res, err
	}
	return res, nil
}

func (s *Service) finishRollbackConflict(operation Operation, res RollbackResult) (RollbackResult, error) {
	operation.Status = OperationRollbackPartial
	if res.Message == "" {
		res.Message = "Rollback stopped or was partial; inspect the reported items."
	}
	operation.Message = res.Message
	operation.UpdatedAt = nowRFC3339()
	if err := s.writeOperation(operation); err != nil {
		return res, err
	}
	if err := s.removeMarker(operation.OperationID); err != nil {
		return res, err
	}
	return res, nil
}

func (s *Service) commitRecovered(item *JournalItem, staged worldstore.StagedImport) error {
	commit, err := s.store.CommitImport(staged)
	if err != nil {
		return err
	}
	if commit.Status != "saved" {
		return fmt.Errorf("commit returned %s: %s", commit.Status, commit.Message)
	}
	item.State = StateCommitted
	item.FinalPath = commit.Path
	item.Revision = commit.Revision
	return nil
}

func (s *Service) verifyJournalItem(item JournalItem) (bool, error) {
	return s.store.VerifyImported(item.ItemID, item.Revision, s.journalHashes(item))
}

func (s *Service) finalExists(item JournalItem) bool {
	if item.FinalPath == "" {
		return false
	}
	_, err := statPath(item.FinalPath)
	return err == nil
}

func (s *Service) finalize(operation *Operation) {
	partial := false
	ambiguous := false
	for _, item := range operation.Items {
		switch item.State {
		case StateQuarantined, StateFailed, StateCleanupDone:
			partial = true
		case StateAmbiguous, StateCleanupFailed:
			ambiguous = true
			partial = true
		}
	}
	switch {
	case ambiguous:
		operation.Status = OperationAmbiguous
		operation.Message = "At least one item requires explicit inspection; completed items remain journaled."
	case partial:
		operation.Status = OperationCompletedPartial
		operation.Message = "The operation completed with quarantined or failed items; no unsafe work remains."
	default:
		operation.Status = OperationCompleted
		operation.Message = "All importable items completed."
	}
	operation.UpdatedAt = nowRFC3339()
}

// countQuarantined matches actionForState's quarantine outcome. In-progress
// planned/prepared items are not quarantined; failed, ambiguous, cleanup-failed,
// and cleaned failed/quarantined dispositions are. A prepared item removed
// during rollback is rolled_back, not cleanup_complete.
func countQuarantined(items []JournalItem) int {
	count := 0
	for _, item := range items {
		switch item.State {
		case StateQuarantined, StateFailed, StateAmbiguous, StateCleanupFailed, StateCleanupDone:
			count++
		}
	}
	return count
}

func resultFromOperation(operation Operation, status, message string) Result {
	planItems := make([]PlanItem, 0, len(operation.Items))
	items := make([]ItemResult, 0, len(operation.Items))
	for _, item := range operation.Items {
		planItems = append(planItems, PlanItem{
			SourceKey:  item.SourceKey,
			SourceKind: item.SourceKind,
			SourceID:   item.SourceID,
			SourceName: item.SourceName,
			SourceRoot: item.SourceRoot,
			Folder:     item.Folder,
			ProposedID: item.ItemID,
			Action:     actionForState(item.State),
			Reason:     item.Reason,
			Message:    item.Message,
		})
		items = append(items, ItemResult{
			SourceKey: item.SourceKey,
			ItemID:    item.ItemID,
			State:     item.State,
			Reason:    item.Reason,
			Path:      item.FinalPath,
			Message:   item.Message,
		})
	}
	return Result{
		OperationID: operation.OperationID,
		Status:      status,
		PlanHash:    operation.PlanHash,
		Plan: Plan{
			Hash:                 operation.PlanHash,
			Category:             operation.Category,
			DestinationLibraryID: operation.DestinationLibraryID,
			DestinationRoot:      operation.DestinationRoot,
			ConversionVersion:    ConversionVersion,
			Items:                planItems,
			Quarantined:          countQuarantined(operation.Items),
			Conflicted:           countConflicts(operation.Items),
			Unsupported:          countReason(operation.Items, "unsupported_source_file") + countReason(operation.Items, "unsupported_media_manifest"),
			Ready:                countState(operation.Items, StateCommitted),
			Adoptable:            countState(operation.Items, StateAdopted),
			Writes:               countState(operation.Items, StateCommitted),
		},
		Items:   items,
		Message: message,
	}
}

func actionForState(state string) string {
	switch state {
	case StateCommitted, StateRolledBack:
		return ActionReady
	case StateAdopted:
		return ActionAdopt
	default:
		return ActionQuarantine
	}
}

func countState(items []JournalItem, state string) int {
	count := 0
	for _, item := range items {
		if item.State == state {
			count++
		}
	}
	return count
}

func countReason(items []JournalItem, reason string) int {
	count := 0
	for _, item := range items {
		if item.State == StateQuarantined && item.Reason == reason {
			count++
		}
	}
	return count
}

func countConflicts(items []JournalItem) int {
	count := 0
	for _, item := range items {
		if item.State != StateQuarantined {
			continue
		}
		switch item.Reason {
		case "identity_conflict", "lore_conflict", "media_conflict":
			count++
		}
	}
	return count
}

func commitReason(status string) string {
	switch status {
	case "collision":
		return "identity_conflict"
	case "authority_indeterminate":
		return "authority_indeterminate"
	case "destination_unavailable":
		return "destination_unavailable"
	case "permission_denied":
		return "destination_not_rollback_safe"
	default:
		return "write_failed"
	}
}

func operationTerminal(status string) bool {
	switch status {
	case OperationCompleted, OperationCompletedPartial, OperationRolledBack, OperationRollbackPartial:
		return true
	}
	return false
}

func sourceNameKey(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func sourceFingerprint(source sourceInfo) string {
	return hashJSON(map[string]interface{}{
		"id":    source.id,
		"name":  source.name,
		"files": source.sourceHashes,
	})
}

func loreHashes(lore map[string][]byte) map[string]string {
	out := map[string]string{}
	for name, data := range lore {
		out[name] = hashBytes(data)
	}
	return out
}

func writeNames(s *Service, prepared preparedItem) []string {
	out := []string{s.documentName()}
	names := make([]string, 0, len(prepared.lore))
	for name := range prepared.lore {
		names = append(names, name)
	}
	sort.Strings(names)
	out = append(out, names...)
	if len(prepared.media) > 0 {
		out = append(out, "media.json")
	}
	return out
}

func planHash(plan Plan) string {
	type hashItem struct {
		SourceKey    string            `json:"source_key"`
		SourceKind   string            `json:"source_kind"`
		SourceID     string            `json:"source_id"`
		SourceName   string            `json:"source_name"`
		SourceRoot   string            `json:"source_root"`
		Folder       string            `json:"folder"`
		ProposedID   string            `json:"proposed_id"`
		Action       string            `json:"action"`
		Reason       string            `json:"reason"`
		SourceHash   string            `json:"source_hash"`
		DocumentHash string            `json:"document_hash"`
		LoreHashes   map[string]string `json:"lore_hashes"`
		MediaHash    string            `json:"media_hash"`
		Writes       []string          `json:"writes"`
		References   []ReferenceIssue  `json:"references"`
	}
	items := make([]hashItem, 0, len(plan.Items))
	for _, item := range plan.Items {
		items = append(items, hashItem{
			SourceKey:    item.SourceKey,
			SourceKind:   item.SourceKind,
			SourceID:     item.SourceID,
			SourceName:   item.SourceName,
			SourceRoot:   item.SourceRoot,
			Folder:       item.Folder,
			ProposedID:   item.ProposedID,
			Action:       item.Action,
			Reason:       item.Reason,
			SourceHash:   item.SourceHash,
			DocumentHash: item.DocumentHash,
			LoreHashes:   item.LoreHashes,
			MediaHash:    item.MediaHash,
			Writes:       item.Writes,
			References:   item.References,
		})
	}
	return hashJSON(map[string]interface{}{
		"category":           plan.Category,
		"destination_id":     plan.DestinationLibraryID,
		"destination_root":   plan.DestinationRoot,
		"destination_status": plan.DestinationStatus,
		"conversion_version": plan.ConversionVersion,
		"items":              items,
		"errors":             plan.Errors,
		"ready":              plan.Ready,
		"adoptable":          plan.Adoptable,
		"conflicted":         plan.Conflicted,
		"quarantined":        plan.Quarantined,
		"unsupported":        plan.Unsupported,
		"reference_issues":   plan.ReferenceIssues,
		"writes":             plan.Writes,
	})
}

func (s *Service) journalHashes(item JournalItem) map[string]string {
	hashes := map[string]string{}
	if item.DocumentHash != "" {
		hashes[s.documentName()] = item.DocumentHash
	}
	for name, hash := range item.LoreHashes {
		hashes[name] = hash
	}
	if item.MediaHash != "" {
		hashes["media.json"] = item.MediaHash
	}
	return hashes
}

func removalMessage(result worldstore.SaveResult, err error) string {
	if err != nil {
		return err.Error()
	}
	if result.Message != "" {
		return result.Message
	}
	return result.Status
}

func (s *Service) stagingPathExists(operationID string, item JournalItem) bool {
	if item.Root == "" || item.ItemID == "" {
		return false
	}
	_, err := os.Lstat(filepath.Join(item.Root, worldstore.ImportStagingDir, operationID, item.ItemID))
	return err == nil
}

func (s *Service) stagingOwnerExists(item JournalItem) bool {
	if item.OwnerPath == "" {
		return false
	}
	_, err := os.Lstat(item.OwnerPath)
	return err == nil
}

func newOperationID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return "imp_fallback"
	}
	return "imp_" + hex.EncodeToString(buf)
}
