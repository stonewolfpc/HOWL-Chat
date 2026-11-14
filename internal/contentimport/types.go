// Package contentimport implements the authorized F3-C explicit import slice:
// browser-record import and the LF-1 one-time legacy-folder migration. It plans
// from user-supplied copies, writes canonical items only through worldstore, and
// keeps an explicit journal outside content roots.
package contentimport

import (
	"time"

	"howl-chat/internal/library"
	"howl-chat/internal/worldstore"
)

const (
	// ConversionVersion changes whenever canonical conversion changes.
	ConversionVersion = "f3c-import-v1"

	ActionReady      = "ready"
	ActionAdopt      = "adopt"
	ActionPending    = "pending"
	ActionQuarantine = "quarantine"

	StatePlanned       = "planned"
	StatePrepared      = "prepared"
	StateCommitted     = "committed"
	StateAdopted       = "adopted"
	StateQuarantined   = "quarantined"
	StateRolledBack    = "rolled_back"
	StateFailed        = "failed"
	StateAmbiguous     = "ambiguous"
	StateCleanupFailed = "cleanup_failed"
	StateCleanupDone   = "cleanup_complete"

	OperationExecuting        = "executing"
	OperationCompleted        = "completed"
	OperationCompletedPartial = "completed_partial"
	OperationRolledBack       = "rolled_back"
	OperationRollbackPartial  = "rollback_partial"
	OperationAmbiguous        = "ambiguous"
	OperationIntegrityFailure = "integrity_failure"
)

// Request is the renderer's explicit import request. BrowserRecords are deep
// copies supplied by the renderer; Go never reads localStorage. LegacyRoots are
// explicit, ephemeral source selections and are never registered.
type Request struct {
	Category             string                   `json:"category"`
	DestinationLibraryID string                   `json:"destination_library_id"`
	BrowserRecords       []map[string]interface{} `json:"browser_records"`
	LegacyRoots          []string                 `json:"legacy_roots"`
	Resolutions          []Resolution             `json:"resolutions"`
	OperationID          string                   `json:"operation_id,omitempty"`
}

// Resolution is an explicit identity decision. It is never inferred from a
// title, folder name, list order, or content similarity.
type Resolution struct {
	SourceKey     string `json:"source_key"`
	DestinationID string `json:"destination_id,omitempty"`
	NewID         bool   `json:"new_id,omitempty"`
	Actor         string `json:"actor,omitempty"`
	Evidence      string `json:"evidence,omitempty"`
}

// Plan is the preflight report and the authorization token's visible content.
// Hash is returned to the renderer and must be supplied again to Execute.
type Plan struct {
	Hash                 string     `json:"hash"`
	Category             string     `json:"category"`
	DestinationLibraryID string     `json:"destination_library_id"`
	DestinationRoot      string     `json:"destination_root,omitempty"`
	DestinationStatus    string     `json:"destination_status"`
	ConversionVersion    string     `json:"conversion_version"`
	Items                []PlanItem `json:"items"`
	Ready                int        `json:"ready"`
	Adoptable            int        `json:"adoptable"`
	Conflicted           int        `json:"conflicted"`
	Quarantined          int        `json:"quarantined"`
	Unsupported          int        `json:"unsupported"`
	ReferenceIssues      int        `json:"reference_issues"`
	Writes               int        `json:"writes"`
	Errors               []string   `json:"errors,omitempty"`
}

// ReferenceIssue records one supported, legacy, malformed, or duplicate
// reference encountered during import. The source value is preserved in the
// document; this report explains how it was classified.
type ReferenceIssue struct {
	Field      string `json:"field"`
	Index      int    `json:"index"`
	Reason     string `json:"reason"`
	Kind       string `json:"kind,omitempty"`
	ID         string `json:"id,omitempty"`
	LegacyName string `json:"legacy_name,omitempty"`
	Message    string `json:"message,omitempty"`
}

// PlanItem is one source's disposition. Reasons are explicit and are never
// converted into a success by the renderer.
type PlanItem struct {
	SourceKey    string            `json:"source_key"`
	SourceKind   string            `json:"source_kind"`
	SourceID     string            `json:"source_id,omitempty"`
	SourceName   string            `json:"source_name,omitempty"`
	SourceRoot   string            `json:"source_root,omitempty"`
	Folder       string            `json:"folder,omitempty"`
	ProposedID   string            `json:"proposed_id,omitempty"`
	Action       string            `json:"action"`
	Reason       string            `json:"reason,omitempty"`
	Message      string            `json:"message,omitempty"`
	Evidence     string            `json:"evidence,omitempty"`
	SourceHash   string            `json:"source_hash"`
	DocumentHash string            `json:"document_hash,omitempty"`
	LoreHashes   map[string]string `json:"lore_hashes,omitempty"`
	MediaHash    string            `json:"media_hash,omitempty"`
	Writes       []string          `json:"writes,omitempty"`
	References   []ReferenceIssue  `json:"references,omitempty"`
}

// Result is the explicit outcome of execute/resume. A nil error does not mean
// every item committed.
type Result struct {
	OperationID string       `json:"operation_id"`
	Status      string       `json:"status"`
	PlanHash    string       `json:"plan_hash"`
	Plan        Plan         `json:"plan"`
	Items       []ItemResult `json:"items"`
	Message     string       `json:"message,omitempty"`
}

// ItemResult reports one item after execution or recovery.
type ItemResult struct {
	SourceKey string `json:"source_key"`
	ItemID    string `json:"item_id,omitempty"`
	State     string `json:"state"`
	Reason    string `json:"reason,omitempty"`
	Path      string `json:"path,omitempty"`
	Message   string `json:"message,omitempty"`
}

// RollbackResult reports the guarded two-phase rollback.
type RollbackResult struct {
	OperationID string         `json:"operation_id"`
	Status      string         `json:"status"`
	Items       []RollbackItem `json:"items"`
	Message     string         `json:"message,omitempty"`
}

// RollbackItem reports removal or the exact reason removal did not occur.
type RollbackItem struct {
	ItemID  string `json:"item_id"`
	State   string `json:"state"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

// Operation is the durable journal view exposed to the page.
type Operation struct {
	Version              int           `json:"version"`
	Category             string        `json:"category"`
	OperationID          string        `json:"operation_id"`
	PlanHash             string        `json:"plan_hash"`
	DestinationLibraryID string        `json:"destination_library_id"`
	DestinationRoot      string        `json:"destination_root,omitempty"`
	Status               string        `json:"status"`
	StartedAt            string        `json:"started_at"`
	UpdatedAt            string        `json:"updated_at"`
	Items                []JournalItem `json:"items"`
	Message              string        `json:"message,omitempty"`
}

// JournalItem is one durable per-item state and its recovery evidence.
type JournalItem struct {
	SourceKey    string            `json:"source_key"`
	SourceKind   string            `json:"source_kind"`
	SourceID     string            `json:"source_id,omitempty"`
	SourceName   string            `json:"source_name,omitempty"`
	SourceRoot   string            `json:"source_root,omitempty"`
	Folder       string            `json:"folder,omitempty"`
	ItemID       string            `json:"item_id,omitempty"`
	LibraryID    string            `json:"library_id,omitempty"`
	Root         string            `json:"root,omitempty"`
	FinalPath    string            `json:"final_path,omitempty"`
	StagingPath  string            `json:"staging_path,omitempty"`
	OwnerPath    string            `json:"owner_path,omitempty"`
	Revision     string            `json:"revision,omitempty"`
	PlanHash     string            `json:"plan_hash,omitempty"`
	State        string            `json:"state"`
	Reason       string            `json:"reason,omitempty"`
	Message      string            `json:"message,omitempty"`
	Evidence     string            `json:"evidence,omitempty"`
	SourceHash   string            `json:"source_hash,omitempty"`
	DocumentHash string            `json:"document_hash,omitempty"`
	LoreHashes   map[string]string `json:"lore_hashes,omitempty"`
	MediaHash    string            `json:"media_hash,omitempty"`
}

// preparedItem is re-derived at execution; it is never accepted from a
// renderer-supplied plan.
type preparedItem struct {
	plan       PlanItem
	source     sourceInfo
	document   map[string]interface{}
	origin     worldstore.Origin
	createdAt  string
	lore       map[string][]byte
	media      []byte
	references []ReferenceIssue
}

type builtPlan struct {
	Plan
	prepared []preparedItem
}

type journalFile struct {
	Version    int         `json:"version"`
	Operations []Operation `json:"operations"`
}

type marker struct {
	Version              int      `json:"version"`
	Category             string   `json:"category"`
	OperationID          string   `json:"operation_id"`
	PlanHash             string   `json:"plan_hash"`
	DestinationLibraryID string   `json:"destination_library_id"`
	DestinationRoot      string   `json:"destination_root,omitempty"`
	ItemIDs              []string `json:"item_ids"`
	Status               string   `json:"status"`
	CompletedAt          string   `json:"completed_at"`
}

type sourceInfo struct {
	key          string
	kind         string
	id           string
	name         string
	root         string
	folder       string
	raw          map[string]interface{}
	body         map[string]interface{}
	rawHash      string
	docHash      string
	lore         map[string][]byte
	media        []byte
	sourceHashes map[string]string
	blockReason  string
	blockMessage string
	unsupported  []string
}

func categoryOf(raw string) library.Category {
	switch raw {
	case "world", "worlds":
		return library.Worlds
	case "scenario", "scenarios":
		return library.Scenarios
	default:
		return library.Category(raw)
	}
}

func nowRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}
