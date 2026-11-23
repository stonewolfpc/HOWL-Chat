package chat

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"howl-chat/internal/backend/lorebook"
)

// SessionEntry is one turn of the active conversation.
//
// It replaces the frontend-shaped map the application used to keep in memory.
// Sender, Character and CharacterID are three different things and stay apart:
// Sender is the display attribution written to history, Character is the name
// of the character whose voice produced an assistant turn, and CharacterID is
// the durable identity that authorizes delivery of character-scoped content.
type SessionEntry struct {
	Role        string
	Content     string
	Sender      string
	Character   string
	CharacterID string
	HasImage    bool
	Timestamp   string
}

// Session owns the active conversation: its turns, the world, scenario,
// characters and lorebooks selected for it, and the once-per-chat trigger
// bookkeeping that lore resolution updates.
//
// It is the single in-memory owner of that state. Persistence writes it through
// Project and reads it back through Load; the Wails layer only converts to and
// from renderer shapes. It holds no files, no paths and no process handles, so a
// test can construct one without a disk.
type Session struct {
	mu         sync.RWMutex
	entries    []SessionEntry
	world      *WorldContext
	scenario   *ScenarioContext
	characters []*CharacterContext
	lorebooks  []lorebook.Entry
	triggered  map[string]bool
}

// NewSession returns an empty session.
func NewSession() *Session {
	return &Session{triggered: map[string]bool{}}
}

// Entries returns a copy of the turns in order. Callers that build prompts must
// not be able to mutate stored history by editing the returned slice.
func (s *Session) Entries() []SessionEntry {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]SessionEntry, len(s.entries))
	copy(out, s.entries)
	return out
}

// Len reports how many turns are stored, including ones with empty content.
func (s *Session) Len() int {
	if s == nil {
		return 0
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}

// AppendEntry adds one turn.
func (s *Session) AppendEntry(entry SessionEntry) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, entry)
}

// AppendUserTurn adds the user's turn unless the immediately preceding turn is
// the identical user message.
//
// The frontend records the user's message before it starts inference. Without
// this guard the same turn is appended a second time and buildPromptMessages
// sends it to the model twice.
func (s *Session) AppendUserTurn(content string) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if n := len(s.entries); n > 0 {
		last := s.entries[n-1]
		if last.Role == "user" && last.Content == content {
			return
		}
	}
	s.entries = append(s.entries, SessionEntry{
		Role:      "user",
		Content:   content,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

// ReplaceEntries replaces the whole history. This is the delete/edit path: the
// surviving turns are the new history, so a deleted turn cannot be replayed from
// anything left behind.
func (s *Session) ReplaceEntries(entries []SessionEntry) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = make([]SessionEntry, len(entries))
	copy(s.entries, entries)
}

// ClearEntries empties the history and the trigger bookkeeping. The chat's
// world, scenario, characters and lorebooks are untouched.
func (s *Session) ClearEntries() {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = nil
	s.triggered = map[string]bool{}
}

// Reset empties everything the session holds.
func (s *Session) Reset() {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = nil
	s.world = nil
	s.scenario = nil
	s.characters = nil
	s.lorebooks = nil
	s.triggered = map[string]bool{}
}

// World returns the selected world, or nil.
func (s *Session) World() *WorldContext {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.world
}

// SetWorld selects the world.
func (s *Session) SetWorld(world *WorldContext) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.world = world
}

// ClearWorld deselects the world.
func (s *Session) ClearWorld() {
	s.SetWorld(nil)
}

// Scenario returns the selected scenario, or nil.
func (s *Session) Scenario() *ScenarioContext {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.scenario
}

// SetScenario selects the scenario.
func (s *Session) SetScenario(scenario *ScenarioContext) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scenario = scenario
}

// ClearScenario deselects the scenario.
func (s *Session) ClearScenario() {
	s.SetScenario(nil)
}

// Characters returns the selected characters. The slice is a copy; the entries
// are the same pointers the caller supplied.
func (s *Session) Characters() []*CharacterContext {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.characters == nil {
		return nil
	}
	out := make([]*CharacterContext, len(s.characters))
	copy(out, s.characters)
	return out
}

// SetCharacters replaces the selected characters.
func (s *Session) SetCharacters(characters []*CharacterContext) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if characters == nil {
		s.characters = nil
		return
	}
	s.characters = make([]*CharacterContext, len(characters))
	copy(s.characters, characters)
}

// UpsertCharacter adds a character, replacing the existing entry with the same
// durable ID.
func (s *Session) UpsertCharacter(character *CharacterContext) {
	if s == nil || character == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, existing := range s.characters {
		if existing != nil && existing.ID == character.ID {
			s.characters[i] = character
			return
		}
	}
	s.characters = append(s.characters, character)
}

// RemoveCharacter drops the character with this durable ID.
func (s *Session) RemoveCharacter(characterID string) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	filtered := make([]*CharacterContext, 0, len(s.characters))
	for _, character := range s.characters {
		if character != nil && character.ID != characterID {
			filtered = append(filtered, character)
		}
	}
	s.characters = filtered
}

// Lorebooks returns the active lorebook pool. The slice is a copy, so enabling
// an entry requires SetLorebooks with the edited slice.
func (s *Session) Lorebooks() []lorebook.Entry {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.lorebooks == nil {
		return nil
	}
	out := make([]lorebook.Entry, len(s.lorebooks))
	copy(out, s.lorebooks)
	return out
}

// SetLorebooks replaces the active lorebook pool.
func (s *Session) SetLorebooks(entries []lorebook.Entry) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if entries == nil {
		s.lorebooks = nil
		return
	}
	s.lorebooks = make([]lorebook.Entry, len(entries))
	copy(s.lorebooks, entries)
}

// ResolveLore resolves the active lorebooks for one turn and records which
// once-per-chat/scene/message entries fired.
//
// Trigger bookkeeping is session state, so it is updated here rather than by a
// caller holding a map. req.Triggered is filled from the session and ignored on
// the way in.
func (s *Session) ResolveLore(req lorebook.ResolveRequest) []lorebook.ResolvedEntry {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.triggered == nil {
		s.triggered = map[string]bool{}
	}
	req.Triggered = s.triggered
	resolved := lorebook.Resolve(s.lorebooks, req)
	for _, entry := range resolved {
		if entry.TriggerFrequency != lorebook.FrequencyAlways {
			s.triggered[entry.ID] = true
		}
	}
	return resolved
}

// History projects the turns into the persisted history shape.
func (s *Session) History() []ChatHistory {
	entries := s.Entries()
	history := make([]ChatHistory, 0, len(entries))
	for _, entry := range entries {
		if strings.TrimSpace(entry.Content) == "" {
			continue
		}
		sender := "System"
		switch entry.Role {
		case "user":
			sender = "User"
			if entry.Sender != "" {
				sender = entry.Sender
			}
		case "assistant":
			switch {
			case entry.Character != "":
				sender = entry.Character
			case entry.Sender != "":
				sender = entry.Sender
			default:
				sender = "Assistant"
			}
		}
		timestamp := entry.Timestamp
		if timestamp == "" {
			timestamp = time.Now().Format(time.RFC3339)
		}
		history = append(history, ChatHistory{
			ID:          strconv.Itoa(len(history)),
			Sender:      sender,
			Text:        entry.Content,
			Role:        entry.Role,
			Timestamp:   timestamp,
			HasImage:    entry.HasImage,
			CharacterID: entry.CharacterID,
		})
	}
	return history
}

// Load replaces the session with a stored conversation. Only meta.json is
// required by the store, so history, companions and lorebooks are all optional.
func (s *Session) Load(data *ChatData) {
	if s == nil || data == nil {
		return
	}
	entries := make([]SessionEntry, 0, len(data.History))
	for _, h := range data.History {
		entry := SessionEntry{
			Role:        h.Role,
			Content:     h.Text,
			Timestamp:   h.Timestamp,
			HasImage:    h.HasImage,
			CharacterID: h.CharacterID,
		}
		// Preserve attribution: a named speaker becomes the character for an
		// assistant turn and the sender otherwise. The generic labels carry no
		// ownership and are dropped.
		if h.Sender != "" && h.Sender != "User" && h.Sender != "Assistant" && h.Sender != "System" {
			if h.Role == "assistant" {
				entry.Character = h.Sender
			} else {
				entry.Sender = h.Sender
			}
		} else if h.Sender != "" && h.Sender != "System" {
			entry.Sender = h.Sender
		}
		entries = append(entries, entry)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = entries
	s.world = data.World
	s.scenario = data.Scenario
	s.characters = data.Characters
	s.triggered = map[string]bool{}
}

// Project writes the session's state into a stored conversation record.
//
// A history is written only when the session has turns or the stored record has
// none: an unloaded chat keeps its stored history instead of being truncated by
// an empty session. Character and lorebook tracking arrays are recomputed from
// what the session currently holds, and an entry the user removed is dropped
// from the removed list only when it is active again.
func (s *Session) Project(data *ChatData) {
	if s == nil || data == nil {
		return
	}
	history := s.History()
	s.mu.RLock()
	world := s.world
	scenario := s.scenario
	characters := make([]*CharacterContext, len(s.characters))
	copy(characters, s.characters)
	lorebooks := make([]lorebook.Entry, len(s.lorebooks))
	copy(lorebooks, s.lorebooks)
	s.mu.RUnlock()

	if len(history) > 0 || len(data.History) == 0 {
		data.History = history
	}
	data.Meta.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if world != nil {
		data.World = world
	}
	if scenario != nil {
		data.Scenario = scenario
	}
	if characters != nil {
		data.Characters = characters
	}

	lorebookItems := make([]interface{}, 0, len(lorebooks))
	lorebookNames := make([]string, 0, len(lorebooks))
	for _, lb := range lorebooks {
		lorebookItems = append(lorebookItems, map[string]interface{}{
			"id": lb.ID, "title": lb.Title, "content": lb.Content, "tags": lb.Tags,
		})
		if lb.Title != "" {
			lorebookNames = append(lorebookNames, lb.Title)
		}
	}
	data.Lorebooks = lorebookItems
	if data.Memories == nil {
		data.Memories = []interface{}{}
	}
	// Membership is recomputed from the session, which is the authority for what
	// this conversation currently holds. An empty session leaves the stored
	// selection alone: a conversation that was never loaded must not be emptied
	// by a save.
	if world != nil {
		if data.Meta.WorldName == "" {
			data.Meta.WorldName = world.Name
		}
		data.Meta.WorldID = world.ID
	}
	if scenario != nil {
		if data.Meta.ScenarioName == "" {
			data.Meta.ScenarioName = scenario.Name
		}
		data.Meta.ScenarioID = scenario.ID
	}
	if len(characters) > 0 {
		names := make([]string, 0, len(characters))
		for _, c := range characters {
			if c != nil && c.Name != "" {
				names = append(names, c.Name)
			}
		}
		data.Meta.CharacterNames = names
	}
	data.Meta.LorebookNames = lorebookNames

	activeCharIDs := make([]string, 0, len(characters))
	for _, c := range characters {
		if c != nil && c.ID != "" {
			activeCharIDs = append(activeCharIDs, c.ID)
		}
	}
	removedChars := make([]string, 0, len(data.Meta.RemovedCharacters))
	for _, id := range data.Meta.RemovedCharacters {
		if !contains(activeCharIDs, id) {
			removedChars = append(removedChars, id)
		}
	}
	data.Meta.ActiveCharacters = activeCharIDs
	data.Meta.RemovedCharacters = removedChars
	data.Meta.AddedCharacters = []string{}

	activeBookIDs := make([]string, 0, len(lorebooks))
	for _, lb := range lorebooks {
		if lb.ID != "" {
			activeBookIDs = append(activeBookIDs, lb.ID)
		}
	}
	removedBooks := make([]string, 0, len(data.Meta.RemovedLorebooks))
	for _, id := range data.Meta.RemovedLorebooks {
		if !contains(activeBookIDs, id) {
			removedBooks = append(removedBooks, id)
		}
	}
	data.Meta.ActiveLorebooks = activeBookIDs
	data.Meta.RemovedLorebooks = removedBooks
	data.Meta.AddedLorebooks = []string{}
}
