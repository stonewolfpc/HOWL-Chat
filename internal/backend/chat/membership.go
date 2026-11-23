package chat

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// Membership is the single authority for which conversations belong to the
// collection and what happened to them. It is not an index of a directory
// listing: the listing is the payload, and this log is what decides whether a
// directory is a conversation, a leftover, or a stranger.
//
// The log is append-only, one complete JSON object per line, written through the
// same appendLine helper as every other write in this package. Folding it gives
// the membership; nothing else is consulted for identity.
//
// Lifecycle order is deliberate:
//   - created is recorded after the directory has been claimed, before the
//     caller is told the conversation exists. A crash in between leaves a
//     directory with no membership record, which reconciliation reports as
//     orphaned rather than silently adopting it.
//   - deleted is recorded before anything is removed. A crash in between leaves
//     a directory whose membership record says deleted, which reconciliation
//     reports as orphaned and which can never be listed or recreated.
//
// deleted.ndjson is absorbed: if it exists and no lifecycle log does, its
// records are migrated once into the lifecycle log as deleted events. After
// that this program never writes or reads it again.

const (
	membershipEventCreated = "created"
	membershipEventDeleted = "deleted"

	lifecycleLogName = "lifecycle.ndjson"
	legacyDeletedLog = "deleted.ndjson"
)

// ConversationState is the reconciled condition of one conversation.
type ConversationState string

const (
	// StatePresent means the collection claims it and its directory is there.
	StatePresent ConversationState = "present"
	// StateMissing means the collection claims it and its directory is gone.
	StateMissing ConversationState = "missing"
	// StateOrphaned means a directory exists that the collection does not claim,
	// either because it was never created through the lifecycle or because its
	// deletion was recorded and the files survived.
	StateOrphaned ConversationState = "orphaned"
)

// Conversation is one reconciled member of the collection.
type Conversation struct {
	ID        string
	Name      string
	State     ConversationState
	Directory string
	CreatedAt string
}

type membershipEvent struct {
	Event     string `json:"event"`
	ID        string `json:"id"`
	Name      string `json:"name,omitempty"`
	At        string `json:"at,omitempty"`
	Directory string `json:"directory,omitempty"`
}

type memberRecord struct {
	Name string
	At   string
}

// Membership owns created and deleted identifiers.
type Membership struct {
	mu      sync.Mutex
	logPath string
	created map[string]memberRecord
	deleted map[string]string
}

// LoadMembership reads the lifecycle log, migrating a legacy deleted.ndjson once
// when the log does not exist yet.
func LoadMembership(root string) (*Membership, error) {
	m := &Membership{
		logPath: filepath.Join(root, lifecycleLogName),
		created: make(map[string]memberRecord),
		deleted: make(map[string]string),
	}
	_, statErr := os.Stat(m.logPath)
	firstRun := os.IsNotExist(statErr)
	if statErr != nil && !os.IsNotExist(statErr) {
		return nil, fmt.Errorf("membership: stat %s: %w", m.logPath, statErr)
	}

	if err := m.foldLog(); err != nil {
		return nil, err
	}
	if err := m.absorbLegacyDeletions(root); err != nil {
		return nil, err
	}
	// First run of this authority on an existing install: the directories that
	// are already there are the collection, and they are written into the log
	// once so later reconciliation has something to compare against. This is a
	// recorded migration, not adoption: a directory added after the log exists
	// is never claimed, it is reported as orphaned.
	if firstRun {
		if err := m.bootstrapFromDirectories(filepath.Join(root, "chats")); err != nil {
			return nil, err
		}
	}
	return m, nil
}

// bootstrapFromDirectories records the conversations an existing install already
// has. Only directories that identify themselves with readable metadata and a
// usable identifier are recorded; anything else stays orphaned and is reported.
func (m *Membership) bootstrapFromDirectories(chatsDir string) error {
	entries, err := os.ReadDir(chatsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("membership: read %s: %w", chatsDir, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		id := entry.Name()
		if !validChatID(id) {
			continue
		}
		if _, gone := m.deleted[id]; gone {
			continue
		}
		if _, known := m.created[id]; known {
			continue
		}

		name := orphanName(filepath.Join(chatsDir, id))
		if name == "" {
			// No readable meta.json: the folder does not identify itself as a
			// conversation, so it is left unclaimed and reported as orphaned.
			continue
		}
		ev := membershipEvent{Event: membershipEventCreated, ID: id, Name: name, At: nowRFC3339()}
		if err := m.append(ev); err != nil {
			return err
		}
		m.apply(ev)
	}
	return nil
}

// foldLog reads every complete line. A torn final line is ignored: the record it
// would have added was never durable, and refusing to start would be worse than
// losing an event the writer never confirmed.
func (m *Membership) foldLog() error {
	f, err := os.Open(m.logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("membership: open %s: %w", m.logPath, err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	for i, line := range lines {
		var ev membershipEvent
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			// A torn final line is the record a crash never finished writing: it
			// was never confirmed, so it is ignored. Damage anywhere else is a
			// real problem and is reported rather than papered over.
			if i == len(lines)-1 {
				break
			}
			return fmt.Errorf("membership: malformed record in %s: %w", m.logPath, err)
		}
		m.apply(ev)
	}
	return nil
}

// apply folds one event. Deletion wins permanently: a later created event for a
// deleted identifier is ignored, so a resurrected directory cannot re-enter the
// collection by writing itself a creation record.
func (m *Membership) apply(ev membershipEvent) {
	switch ev.Event {
	case membershipEventCreated:
		if _, gone := m.deleted[ev.ID]; gone {
			return
		}
		m.created[ev.ID] = memberRecord{Name: ev.Name, At: ev.At}
	case membershipEventDeleted:
		if ev.ID == "" {
			return
		}
		m.deleted[ev.ID] = ev.At
		delete(m.created, ev.ID)
	}
}

// absorbLegacyDeletions migrates deleted.ndjson once, and only when this
// authority has never written its own log. It is a compatibility read followed
// by a normal append; the legacy file is left untouched as user data.
func (m *Membership) absorbLegacyDeletions(root string) error {
	if _, err := os.Stat(m.logPath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("membership: stat %s: %w", m.logPath, err)
	}

	legacy := filepath.Join(root, legacyDeletedLog)
	f, err := os.Open(legacy)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("membership: open %s: %w", legacy, err)
	}
	defer f.Close()

	var events []membershipEvent
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var record struct {
			ID        string `json:"id"`
			At        string `json:"at"`
			Directory string `json:"directory"`
		}
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			return fmt.Errorf("membership: malformed legacy deletion in %s: %w", legacy, err)
		}
		if record.ID == "" {
			continue
		}
		events = append(events, membershipEvent{
			Event:     membershipEventDeleted,
			ID:        record.ID,
			At:        record.At,
			Directory: record.Directory,
		})
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	for _, ev := range events {
		if err := m.append(ev); err != nil {
			return err
		}
		m.apply(ev)
	}
	return nil
}

func (m *Membership) append(ev membershipEvent) error {
	line, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	if err := appendLine(m.logPath, append(line, '\n')); err != nil {
		return fmt.Errorf("membership: append %s: %w", m.logPath, err)
	}
	return nil
}

// Create records a conversation that has just claimed its directory.
func (m *Membership) Create(id, name string) error {
	if id == "" {
		return &Refusal{Code: RefusedUnbound, Op: "membership-create", Detail: "no conversation id"}
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, gone := m.deleted[id]; gone {
		return &Refusal{Code: RefusedUnbound, Op: "membership-create", ChatID: id, Detail: "conversation was deleted"}
	}
	ev := membershipEvent{Event: membershipEventCreated, ID: id, Name: name, At: nowRFC3339()}
	if err := m.append(ev); err != nil {
		return err
	}
	m.apply(ev)
	return nil
}

// Delete records intent. It is written before removal, so a crash can only
// leave files behind, never a conversation that comes back.
func (m *Membership) Delete(id, directory string) error {
	if id == "" {
		return &Refusal{Code: RefusedUnbound, Op: "membership-delete", Detail: "no conversation id"}
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	ev := membershipEvent{Event: membershipEventDeleted, ID: id, At: nowRFC3339(), Directory: directory}
	if err := m.append(ev); err != nil {
		return err
	}
	m.apply(ev)
	return nil
}

// Deleted reports whether the collection has recorded this conversation as
// deleted. It is the resurrection check every mutator uses.
func (m *Membership) Deleted(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, gone := m.deleted[id]
	return gone
}

// DeletedSet returns a copy of the deleted identifiers.
func (m *Membership) DeletedSet() map[string]bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.deleted) == 0 {
		return nil
	}
	set := make(map[string]bool, len(m.deleted))
	for id := range m.deleted {
		set[id] = true
	}
	return set
}

// MemberIDs returns the live members in creation order.
func (m *Membership) MemberIDs() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	ids := make([]string, 0, len(m.created))
	for id := range m.created {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// Member reports the recorded name of a live member.
func (m *Membership) Member(id string) (memberRecord, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	rec, ok := m.created[id]
	return rec, ok
}

// Reconcile compares the authority with the directory listing and writes
// nothing. Missing is a claimed conversation whose directory is absent; orphaned
// is a directory the collection does not claim, including the leftover of a
// deletion whose files survived.
func (m *Membership) Reconcile(chatsDir string) ([]Conversation, error) {
	m.mu.Lock()
	members := make(map[string]memberRecord, len(m.created))
	for id, rec := range m.created {
		members[id] = rec
	}
	deleted := make(map[string]bool, len(m.deleted))
	for id := range m.deleted {
		deleted[id] = true
	}
	m.mu.Unlock()

	entries, err := os.ReadDir(chatsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("membership: read %s: %w", chatsDir, err)
	}

	onDisk := make(map[string]bool, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			onDisk[entry.Name()] = true
		}
	}

	result := make([]Conversation, 0, len(members)+len(onDisk))
	for id, rec := range members {
		dir := filepath.Join(chatsDir, id)
		state := StatePresent
		if !onDisk[id] {
			state = StateMissing
		}
		result = append(result, Conversation{
			ID:        id,
			Name:      rec.Name,
			State:     state,
			Directory: dir,
			CreatedAt: rec.At,
		})
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		id := entry.Name()
		if _, ok := members[id]; ok {
			continue
		}
		dir := filepath.Join(chatsDir, id)
		result = append(result, Conversation{
			ID:        id,
			Name:      orphanName(dir),
			State:     StateOrphaned,
			Directory: dir,
		})
	}

	sort.SliceStable(result, func(i, j int) bool {
		if result[i].State != result[j].State {
			return result[i].State < result[j].State
		}
		return result[i].ID < result[j].ID
	})
	return result, nil
}

// orphanName reads whatever the folder still says about itself. A directory the
// collection does not claim is never trusted with identity, only with a label.
func orphanName(dir string) string {
	data, err := os.ReadFile(filepath.Join(dir, "meta.json"))
	if err != nil {
		return ""
	}
	var meta struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return ""
	}
	return meta.Name
}

func nowRFC3339() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}
