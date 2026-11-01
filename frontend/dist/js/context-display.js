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

window.updateContextDisplay = () => contextDisplay.updateContextDisplay();
window.updateCharacterPanel = () => contextDisplay.updateCharacterPanel();
// World and scenario panel rendering is owned by sidebar-context-world.js and
// sidebar-context-scenario.js: they write the description elements as well, and
// they load after this file. Publishing a second global here only meant those
// two modules silently replaced it a few hundred lines later.
window.updateLorebookDisplay = () => contextDisplay.updateLorebookDisplay();

// ---------------------------------------------------------------------------
// Init
// ---------------------------------------------------------------------------

document.addEventListener('DOMContentLoaded', () => { contextDisplay.initialize(); });
