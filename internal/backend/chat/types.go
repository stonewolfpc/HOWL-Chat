package chat

import (
	entities "howl-chat/internal/backend/types/entities"
)

// Image represents image data with a path and optional thumbnail
type Image struct {
	Path      string `json:"path"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

// WorldContext extends entities.World with additional context-specific fields
type WorldContext struct {
	entities.World
	Image *Image `json:"image,omitempty"`
}

// ScenarioContextCharacter is a lightweight character reference stored inside a scenario
type ScenarioContextCharacter struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Role   string `json:"role,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

// ScenarioContext extends entities.Scenario with additional context-specific fields
type ScenarioContext struct {
	entities.Scenario
	Image      *Image                     `json:"image,omitempty"`
	Characters []ScenarioContextCharacter `json:"characters,omitempty"`
	Lorebooks  []interface{}              `json:"lorebooks,omitempty"`
}

// CharacterContext extends entities.Character with additional context-specific fields
type CharacterContext struct {
	entities.Character
	Image *Image `json:"image,omitempty"`
}

// ChatMeta represents the metadata for a chat session
type ChatMeta struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
	WorldName      string   `json:"worldName"`
	ScenarioName   string   `json:"scenarioName"`
	CharacterNames []string `json:"characterNames"`
	LorebookNames  []string `json:"lorebookNames"`
}

// ChatHistory represents a single message in the chat history
type ChatHistory struct {
	ID        string `json:"id"`
	Sender    string `json:"sender"`
	Text      string `json:"text"`
	Role      string `json:"role"` // "system", "user", "assistant"
	Timestamp string `json:"timestamp"`
	HasImage  bool   `json:"hasImage"`
}

// ChatData represents the complete chat data
type ChatData struct {
	Meta       ChatMeta            `json:"meta"`
	History    []ChatHistory       `json:"history"`
	World      *WorldContext       `json:"world"`
	Scenario   *ScenarioContext    `json:"scenario"`
	Characters []*CharacterContext `json:"characters"`
	Lorebooks  []interface{}       `json:"lorebooks"` // Keep as interface{} for now, will address lorebooks later
	Memories   []interface{}       `json:"memories"`
}

// AppConfig represents the application configuration
type AppConfig struct {
	LastChatID string `json:"lastChatId"`
}

// ChatSummary represents a lightweight summary for chat lists
type ChatSummary struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
	WorldName      string   `json:"worldName"`
	ScenarioName   string   `json:"scenarioName"`
	CharacterNames []string `json:"characterNames"`
	LorebookNames  []string `json:"lorebookNames"`
	MessageCount   int      `json:"messageCount"`
}
