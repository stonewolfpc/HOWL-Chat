package library

// OperationStatus describes the actual outcome of a content operation. It is
// deliberately more precise than a boolean: a browser-store save can succeed
// while an optional export fails, and the caller must be able to say so.
type OperationStatus string

const (
	StatusSaveSucceeded          OperationStatus = "save_succeeded"
	StatusExportSucceeded        OperationStatus = "export_succeeded"
	StatusPartial                OperationStatus = "partial"
	StatusFailed                 OperationStatus = "failed"
	StatusDestinationUnavailable OperationStatus = "destination_unavailable"
	StatusPermissionDenied       OperationStatus = "permission_denied"
	StatusNameCollision          OperationStatus = "name_collision"
	StatusReplaced               OperationStatus = "replaced"
	StatusCancelled              OperationStatus = "cancelled"
)

// OperationResult is returned by export/copy operations. Target is set only
// when a file was actually written; Message is safe to show to the user.
type OperationResult struct {
	Status   OperationStatus `json:"status"`
	Category Category        `json:"category,omitempty"`
	ItemID   string          `json:"item_id,omitempty"`
	Target   string          `json:"target,omitempty"`
	Message  string          `json:"message,omitempty"`
}

// OK reports whether the requested operation completed in full.
func (r OperationResult) OK() bool {
	return r.Status == StatusSaveSucceeded ||
		r.Status == StatusExportSucceeded ||
		r.Status == StatusReplaced
}
