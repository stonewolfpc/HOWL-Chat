package main

import (
	"howl-chat/internal/backend/chat"
	"howl-chat/internal/backend/lorebook"
)

// The active conversation belongs to the chat domain. App holds one handle to
// chat.Conversation, which owns both the session (turns, world, scenario,
// characters, lorebooks, triggers) and the collection (membership, append log,
// snapshots). These accessors are what the Wails layer reads; none of them is
// storage and none of them may be written to. Writes go through the session's or
// the conversation's own methods.

// session returns the active conversation's session, or nil before boot binds one.
func (a *App) session() *chat.Session {
	if a.conversation == nil {
		return nil
	}
	return a.conversation.Session()
}

// sessionMessages projects the session's turns into the renderer's message shape.
func (a *App) sessionMessages() []map[string]interface{} {
	return sessionEntriesToMaps(a.session().Entries())
}

func (a *App) sessionWorld() *chat.WorldContext {
	return a.session().World()
}

func (a *App) sessionScenario() *chat.ScenarioContext {
	return a.session().Scenario()
}

func (a *App) sessionCharacters() []*chat.CharacterContext {
	return a.session().Characters()
}

func (a *App) sessionLorebooks() []lorebook.Entry {
	return a.session().Lorebooks()
}
