/**
 * Context Display — Orchestrator
 *
 * Instantiates ContextDisplay (defined across the modules below),
 * exposes the legacy global shims consumed by chat.js / chat-session.js,
 * and fires initialization on DOMContentLoaded.
 *
 * Depends on (must load first):
 *   context-data.js      — class definition + data/sync methods
 *   context-panels.js    — updateContextDisplay / panel render methods
 *   context-selectors.js — selector modals + selectWorld/Scenario/Character
 *
 * @module context-display
 * @version 2.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Global Instance
// ---------------------------------------------------------------------------

const contextDisplay = new ContextDisplay();
window.contextDisplay = contextDisplay;

// ---------------------------------------------------------------------------
// Legacy Global Shims
// ---------------------------------------------------------------------------

Object.defineProperty(window, 'currentChatData', {
    get()      { return contextDisplay.currentChatData; },
    set(value) { contextDisplay.currentChatData = value; }
});

window.resetToDefaultState  = () => contextDisplay.resetToDefaultState();
window.updateContextDisplay = () => contextDisplay.updateContextDisplay();
window.updateCharacterPanel = () => contextDisplay.updateCharacterPanel();

// ---------------------------------------------------------------------------
// Init
// ---------------------------------------------------------------------------

document.addEventListener('DOMContentLoaded', () => { contextDisplay.initialize(); });
