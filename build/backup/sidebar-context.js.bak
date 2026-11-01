/**
 * Sidebar Context — Orchestrator
 *
 * Initializes the context panel and exposes window.sidebarContext.
 * All logic is delegated to sub-modules:
 *
 *   sidebar-context-world.js     — world display + selector
 *   sidebar-context-scenario.js  — scenario display + selector + filter
 *   sidebar-context-lorebooks.js — lorebook display + selector + toggle
 *
 * Depends on: sidebar-state.js, sidebar-session.js
 *
 * @module sidebar-context
 * @version 2.0.0
 */

'use strict';

async function initializeContextPanel() {
    console.log('Initializing Context Panel (v2)...');
    await loadContextData();
}

async function loadContextData() {
    loadWorlds();
    loadScenarios();
    // Panel rendering is handled by sidebar-panels.js init flow via contextDisplay
}

function loadWorlds() {
    const worlds = JSON.parse(localStorage.getItem('howl-worlds') || '[]');
    console.log('Loaded worlds:', worlds);
}

function loadScenarios() {
    const scenarios = JSON.parse(localStorage.getItem('howl-chat-scenarios') || '[]');
    console.log('Loaded scenarios:', scenarios);
}

function updateContextDisplay() {
    updateWorldDisplay();
    updateScenarioDisplay();
    updateLorebookDisplay();
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

window.sidebarContext = {
    initializeContextPanel,
    loadContextData,
    updateContextDisplay,
    updateWorldDisplay,
    openWorldSelector,
    selectWorld,
    toggleWorldAutoSwitch,
    updateScenarioDisplay,
    openScenarioSelector,
    filterScenarios,
    selectScenario,
    toggleScenarioAutoSwitch,
    updateLorebookDisplay,
    createLorebookItem,
    removeLorebookFromChat,
    openLorebookSelector,
    filterLorebooks,
    toggleLorebookSelection,
    openLorebookManager
};
