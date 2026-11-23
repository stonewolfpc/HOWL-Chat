package chat

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ErrNoCollection means the conversation has no persistence root: a caller asked
// a collection question before one was bound.
var ErrNoCollection = errors.New("chat persistence not initialized")

// ActiveConversation is the active conversation: one Session bound to one collection.
//
// It owns the persistence wiring the application used to hold — creating and
// loading conversations, membership checks, the append log, the snapshot commit
// and the history rewrite — so that the turns, the log and the rendered transcript
// come from a single object. It holds no window, no event sink and no process
// handles, so it can be constructed and tested without a host.
type ActiveConversation struct {
	mu      sync.Mutex
	session *Session
	store   *PersistenceManager
}

// NewConversation binds a fresh session to a collection. A nil store is allowed:
// the conversation then works in memory only and every persistence call reports
// ErrNoCollection.
func NewActiveConversation(store *PersistenceManager) *ActiveConversation {
	return &ActiveConversation{session: NewSession(), store: store}
}

// Session returns the conversation's turns and content selection.
func (c *ActiveConversation) Session() *Session {
	if c == nil {
		return nil
	}
	return c.session
}

// Store returns the backing collection. It is for collection-level operations the
// host owns (a content wipe staging the root); every ordinary write goes through
// this type's own methods.
func (c *ActiveConversation) Store() *PersistenceManager {
	if c == nil {
		return nil
	}
	return c.store
}

// Bound reports whether a collection is attached.
func (c *ActiveConversation) Bound() bool {
	return c != nil && c.store != nil
}

// RootDir reports the collection root, or "" when unbound.
func (c *ActiveConversation) RootDir() string {
	if !c.Bound() {
		return ""
	}
	return c.store.RootDir()
}

// ActiveID reports the persisted active conversation identifier.
func (c *ActiveConversation) ActiveID() string {
	if !c.Bound() {
		return ""
	}
	return c.store.GetLastChatID()
}

// SetActiveID repoints the collection at a conversation, including the empty
// string that clears the pointer.
func (c *ActiveConversation) SetActiveID(id string) error {
	if !c.Bound() {
		return ErrNoCollection
	}
	return c.store.SetLastChatID(id)
}

// List returns every conversation summary in the collection.
func (c *ActiveConversation) List() ([]ChatSummary, error) {
	if !c.Bound() {
		return []ChatSummary{}, ErrNoCollection
	}
	return c.store.GetChatList()
}

// Meta returns the active conversation's metadata.
func (c *ActiveConversation) Meta() (*ChatMeta, error) {
	if !c.Bound() {
		return nil, ErrNoCollection
	}
	return c.store.GetActiveChatMeta()
}

// Directory reports where a conversation's files live.
func (c *ActiveConversation) Directory(id string) (string, error) {
	if !c.Bound() {
		return "", ErrNoCollection
	}
	return c.store.ChatDirectory(id)
}

// New creates a conversation in the collection and makes it the active one. The
// session is emptied: the new conversation has no turns yet.
func (c *ActiveConversation) New(name string) (*ChatMeta, error) {
	if !c.Bound() {
		return nil, ErrNoCollection
	}
	data, err := c.store.NewChat(name)
	if err != nil {
		return nil, err
	}
	c.session.ClearEntries()
	return &data.Meta, nil
}

// Load makes a stored conversation active. The conversation currently bound is
// saved first so switching never discards the turns that were just sent.
func (c *ActiveConversation) Load(id string) (*ChatData, error) {
	if !c.Bound() {
		return nil, ErrNoCollection
	}
	if last := c.store.GetLastChatID(); last != "" && last != id {
		if err := c.Save(); err != nil {
			// Switching conversations is more important than the previous one's
			// autosave, but the refusal is reported rather than swallowed.
			fmt.Printf("Warning: Failed to save current chat: %v\n", err)
		}
	}
	data, err := c.store.LoadChat(id)
	if err != nil {
		return nil, err
	}
	c.session.Load(data)
	if err := c.store.SetLastChatID(id); err != nil {
		return nil, fmt.Errorf("failed to set active chat: %w", err)
	}
	return data, nil
}

// SaveAs copies the active conversation under a new name.
func (c *ActiveConversation) SaveAs(name string) (*ChatMeta, error) {
	if !c.Bound() {
		return nil, ErrNoCollection
	}
	last := c.store.GetLastChatID()
	if last == "" {
		return nil, fmt.Errorf("no active chat to save")
	}
	data, err := c.store.SaveChatAs(last, name)
	if err != nil {
		return nil, err
	}
	return &data.Meta, nil
}

// Delete removes a conversation. Deleting the active one empties the session, so
// a late autosave cannot resurrect it from memory.
func (c *ActiveConversation) Delete(id string) error {
	if !c.Bound() {
		return ErrNoCollection
	}
	if c.store.GetLastChatID() == id {
		c.session.Reset()
	}
	return c.store.DeleteChat(id)
}

// Save projects the session onto the active conversation and writes it.
func (c *ActiveConversation) Save() error {
	if !c.Bound() {
		return ErrNoCollection
	}
	id := c.store.GetLastChatID()
	if id == "" {
		return fmt.Errorf("no active chat to save")
	}
	data, err := c.store.LoadChat(id)
	if err != nil {
		return err
	}
	c.session.Project(data)
	return c.store.SaveChat(data)
}

// FinalSave commits the active conversation once, on shutdown, and reports how
// many history entries were committed.
//
// Companions, then the log, then the snapshot are written under one lock hold, so
// a stale log line, a deleted message or a legacy DOM-scrape entry cannot survive
// the rewrite, and a conversation that was never loaded keeps its stored history
// instead of being truncated.
func (c *ActiveConversation) FinalSave() (int, error) {
	if !c.Bound() {
		return 0, nil
	}
	id := c.store.GetLastChatID()
	if id == "" {
		return 0, nil
	}
	data, err := c.store.LoadChat(id)
	if err != nil {
		return 0, fmt.Errorf("final save: %w", err)
	}
	c.session.Project(data)
	if err := c.store.CommitChat(data); err != nil {
		return 0, fmt.Errorf("final save: %w", err)
	}
	return len(data.History), nil
}

// AppendTurn records one turn in the session and, when the conversation is bound,
// appends it to the log so a reload replays it.
func (c *ActiveConversation) AppendTurn(entry SessionEntry) error {
	if strings.TrimSpace(entry.Content) == "" {
		return nil
	}
	if entry.Timestamp == "" {
		entry.Timestamp = time.Now().Format(time.RFC3339)
	}
	c.session.AppendEntry(entry)
	if !c.Bound() {
		return nil
	}
	id := c.store.GetLastChatID()
	if id == "" {
		return nil
	}
	return c.store.AppendHistoryEntry(id, ChatHistory{
		ID:          strconv.Itoa(c.session.Len() - 1),
		Sender:      entry.Sender,
		Text:        entry.Content,
		Role:        entry.Role,
		Timestamp:   entry.Timestamp,
		HasImage:    entry.HasImage,
		CharacterID: entry.CharacterID,
	})
}

// RecordUserTurn adds the user's turn to the session without touching the log.
//
// The renderer records the user's message before it starts inference, so the send
// path only needs the duplicate guard here; the log entry was already appended.
func (c *ActiveConversation) RecordUserTurn(content string) {
	c.session.AppendUserTurn(content)
}

// ReplaceHistory replaces the whole history and rewrites the log in place. This
// is the delete/edit path: appending the survivors would leave deleted turns on
// disk to be replayed on the next load.
func (c *ActiveConversation) ReplaceHistory(entries []SessionEntry) error {
	c.session.ReplaceEntries(entries)
	if !c.Bound() {
		return nil
	}
	id := c.store.GetLastChatID()
	if id == "" {
		return nil
	}
	return c.store.RewriteHistory(id, c.session.History())
}

// ClearHistory empties the session and truncates the stored log.
func (c *ActiveConversation) ClearHistory() error {
	c.session.ClearEntries()
	if !c.Bound() {
		return nil
	}
	id := c.store.GetLastChatID()
	if id == "" {
		return nil
	}
	return c.store.RewriteHistory(id, []ChatHistory{})
}

// SaveRecord persists a record supplied by the renderer, refusing an identifier
// the collection does not claim: a cached or stale payload must never become a
// conversation.
func (c *ActiveConversation) SaveRecord(data *ChatData) error {
	if data == nil {
		return &Refusal{Code: RefusedUnbound, Op: "save-data", Detail: "no chat data"}
	}
	if !c.Bound() {
		return nil
	}
	if !c.store.Present(data.Meta.ID) {
		return &Refusal{
			Code:   RefusedUnbound,
			Op:     "save-data",
			ChatID: data.Meta.ID,
			Detail: "conversation is not present in the collection",
		}
	}
	return c.store.SaveChat(data)
}

// Present reports whether the collection claims this conversation.
func (c *ActiveConversation) Present(id string) bool {
	return c.Bound() && c.store.Present(id)
}
