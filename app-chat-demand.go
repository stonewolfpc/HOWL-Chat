package main

import (
	"fmt"
	"strings"
	"time"
)

// app-chat-demand.go — intent-driven chat activation.
//
// Navigation and model selection express interest; they must not move GPU
// memory. A submitted message is the intent to execute, so the send path owns
// the whole sequence here: load (admission decides, the lifecycle gate
// serialises, a needed handoff happens through the same reclaim route), verify
// the engine answers, and only then let the request go out. One owner, one
// order, no page-return race.

// chatReadyTimeout bounds how long a demand-driven send waits for the engine it
// just started to accept connections.
var chatReadyTimeout = 120 * time.Second

const chatReadyPoll = 250 * time.Millisecond

// EnsureChatModelReady prepares the engine for a message that is about to be
// sent. It is a no-op for a loaded model that already answers.
//
// modelPath is the model the user selected; an empty path with nothing loaded
// is a refusal, not a silent default.
func (a *App) EnsureChatModelReady(modelPath string) error {
	if err := a.ensureAcceptingWork("load chat model"); err != nil {
		return err
	}
	if !a.IsModelLoaded() {
		path := strings.TrimSpace(modelPath)
		if path == "" {
			return fmt.Errorf("no chat model selected, so there is nothing to load for this message")
		}
		if a.engine.Client() == nil {
			return fmt.Errorf("llama client not initialized")
		}
		fmt.Printf("[Chat] demand load: no model is loaded, loading %q for this message\n", path)
		if err := a.LoadModel(path); err != nil {
			return err
		}
	}
	client := a.engine.Client()
	if client == nil {
		return fmt.Errorf("llama client not initialized")
	}
	if err := waitForHealth(client.Health, chatReadyTimeout); err != nil {
		return fmt.Errorf("the chat engine did not become ready: %w", err)
	}
	return nil
}

// waitForHealth waits until health reports success or the bound expires. It
// never treats "we started a process" as readiness.
func waitForHealth(health func() error, timeout time.Duration) error {
	if health == nil {
		return fmt.Errorf("no health check available")
	}
	deadline := time.Now().Add(timeout)
	var last error
	for {
		last = health()
		if last == nil {
			return nil
		}
		if time.Now().After(deadline) {
			return last
		}
		time.Sleep(chatReadyPoll)
	}
}
