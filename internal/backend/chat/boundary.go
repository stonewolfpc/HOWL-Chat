package chat

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"howl-chat/internal/durable"
	"howl-chat/internal/paths"
)

// This file is the controlled mutation boundary for conversation content.
//
// Every exported mutator on PersistenceManager enters through withChatLock and
// writes only through the helpers in writers.go. The analysis test
// (boundary_analysis_test.go) fails the build if that stops being true, which is
// what makes the boundary structural rather than a convention.

// RefusalCode names why a mutation did not happen.
type RefusalCode string

const (
	// RefusedBusy means another writer holds the conversation. Nothing was written.
	RefusedBusy RefusalCode = "refused_busy"
	// RefusedReentrant means one operation entered the boundary twice for the
	// same conversation. That is a composition bug, never contention.
	RefusedReentrant RefusalCode = "refused_reentrant"
	// RefusedUnbound means the operation has no conversation bound to it, or the
	// conversation was deleted and must not be recreated.
	RefusedUnbound RefusalCode = "refused_unbound"
	// RefusedContainment means the identifier or resolved path is not inside the
	// conversation root. Fail-closed: the write does not happen.
	RefusedContainment RefusalCode = "refused_containment"
	// RefusedNotWritten means the boundary itself failed before any content was
	// written (lock file, lock directory).
	RefusedNotWritten RefusalCode = "refused_not_written"
)

// Refusal is the typed "nothing was written" result. Callers must surface it;
// a discarded refusal is non-conforming (V3-27).
type Refusal struct {
	Code   RefusalCode
	Op     string
	ChatID string
	Detail string
	Err    error
}

func (r *Refusal) Error() string {
	parts := []string{fmt.Sprintf("chat %s refused", r.Op)}
	if r.ChatID != "" {
		parts = append(parts, "id="+r.ChatID)
	}
	parts = append(parts, string(r.Code))
	if r.Detail != "" {
		parts = append(parts, r.Detail)
	}
	if r.Err != nil {
		parts = append(parts, r.Err.Error())
	}
	return strings.Join(parts, ": ")
}

func (r *Refusal) Unwrap() error { return r.Err }

// IsRefusal reports whether err is a refusal with the given code.
func IsRefusal(err error, code RefusalCode) bool {
	var r *Refusal
	return errors.As(err, &r) && r.Code == code
}

// DeletionIncomplete reports that a conversation left the visible collection but
// its directory could not be removed. It is a warning, not a failure: deletion is
// best-effort by product decision, and leftovers are the user's to inspect.
type DeletionIncomplete struct {
	ChatID    string
	Directory string
	Err       error
}

func (d *DeletionIncomplete) Error() string {
	return fmt.Sprintf("conversation %s was removed from the list; its directory may still exist at %s: %v",
		d.ChatID, d.Directory, d.Err)
}

func (d *DeletionIncomplete) Unwrap() error { return d.Err }

// locksDir is the sibling namespace for per-conversation locks. It is outside
// chats/ so a lock file is never a candidate conversation and so deletion and
// wipe never have to remove a directory that contains an open handle.
func (pm *PersistenceManager) locksDir() string {
	if pm == nil {
		return ""
	}
	return filepath.Join(pm.rootDir, ".chat-locks")
}

// lockPath derives a fixed-length lock file name from the logical conversation
// identifier, so an identifier can never become a path.
func (pm *PersistenceManager) lockPath(chatID string) string {
	sum := sha256.Sum256([]byte(chatID))
	return filepath.Join(pm.locksDir(), hex.EncodeToString(sum[:16])+".lock")
}

// chatLocks tracks locks held by this process. Its only job is to turn a
// same-process double entry into a diagnosed refusal instead of a deadlock or a
// misleading "busy".
type chatLocks struct {
	mu   sync.Mutex
	held map[string]bool
}

var processChatLocks = &chatLocks{held: make(map[string]bool)}

// withChatLock runs fn while holding the conversation's exclusive lock.
//
// Refusal, not waiting, is the specified outcome on contention: a caller that
// cannot write must learn that immediately. The lock is held for one operation
// and released even if fn panics.
func (pm *PersistenceManager) withChatLock(op, chatID string, fn func() error) error {
	if pm == nil {
		return &Refusal{Code: RefusedUnbound, Op: op, Detail: "no persistence manager"}
	}
	if chatID == "" {
		return &Refusal{Code: RefusedUnbound, Op: op, Detail: "no conversation bound"}
	}
	lockPath := pm.lockPath(chatID)

	processChatLocks.mu.Lock()
	if processChatLocks.held[lockPath] {
		processChatLocks.mu.Unlock()
		return &Refusal{
			Code:   RefusedReentrant,
			Op:     op,
			ChatID: chatID,
			Detail: "the same operation entered the mutation boundary twice",
		}
	}
	processChatLocks.held[lockPath] = true
	processChatLocks.mu.Unlock()

	release := func() {
		processChatLocks.mu.Lock()
		delete(processChatLocks.held, lockPath)
		processChatLocks.mu.Unlock()
	}

	if err := ensureDir(pm.locksDir()); err != nil {
		release()
		return &Refusal{Code: RefusedNotWritten, Op: op, ChatID: chatID, Detail: "lock directory", Err: err}
	}
	unlock, err := durable.LockFile(lockPath, 0)
	if err != nil {
		release()
		if errors.Is(err, durable.ErrBusy) {
			return &Refusal{Code: RefusedBusy, Op: op, ChatID: chatID, Detail: "another writer holds this conversation"}
		}
		return &Refusal{Code: RefusedNotWritten, Op: op, ChatID: chatID, Detail: "lock", Err: err}
	}
	defer func() {
		unlock()
		release()
	}()

	return fn()
}

// chatsRoot is the only place the conversation root is spelled.
func (pm *PersistenceManager) chatsRoot() (string, error) {
	if pm == nil || pm.rootDir == "" {
		return "", &Refusal{Code: RefusedContainment, Op: "resolve", Detail: "empty conversation root"}
	}
	root := filepath.Join(pm.rootDir, "chats")
	if err := paths.CheckContained(pm.rootDir, root); err != nil {
		return "", &Refusal{Code: RefusedContainment, Op: "resolve", Detail: "conversation root outside the history root", Err: err}
	}
	return root, nil
}

// conversationDir is the only place a conversation directory is constructed.
// It refuses an identifier that is not a single path component and a directory
// that does not physically resolve inside the history root (a junction planted
// inside chats/ is caught here).
func (pm *PersistenceManager) conversationDir(chatID string) (string, error) {
	if !validChatID(chatID) {
		return "", &Refusal{Code: RefusedContainment, Op: "resolve", ChatID: chatID, Detail: "identifier is not a single path component"}
	}
	root, err := pm.chatsRoot()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(root, chatID)
	if err := paths.CheckContained(pm.rootDir, dir); err != nil {
		return "", &Refusal{Code: RefusedContainment, Op: "resolve", ChatID: chatID, Detail: "conversation directory resolves outside the history root", Err: err}
	}
	return dir, nil
}

func validChatID(chatID string) bool {
	if chatID == "" || chatID == "." || chatID == ".." {
		return false
	}
	if strings.ContainsAny(chatID, `/\`) {
		return false
	}
	if filepath.VolumeName(chatID) != "" {
		return false
	}
	return filepath.Base(chatID) == chatID
}

// deletedSet reports the identifiers the membership authority recorded as
// deleted. Deletion intent is no longer a second file of its own: it is a
// "deleted" event in the lifecycle log (membership.go), which absorbed
// deleted.ndjson once on first load. The existing file is left on disk as user
// data and is never written again.
func (pm *PersistenceManager) deletedSet() map[string]bool {
	if pm == nil || pm.membership == nil {
		return nil
	}
	return pm.membership.DeletedSet()
}
