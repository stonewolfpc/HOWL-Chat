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
    // Worlds are loaded by world-grid.js (it rehydrates images and renders the
    // grid). The log-only copy that lived here shadowed it.
    loadScenarios();
    // Panel rendering is handled by sidebar-panels.js init flow via contextDisplay
}


function loadScenarios() {
    const scenarios = window.ScenarioRepository.all();
    console.log('Loaded scenarios:', scenarios);
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

// Every panel entry is resolved at call time. This module is loaded before
// sidebar-panels.js, and on story-tools.html it loads without context-display.js
// at all, so a direct reference here was a ReferenceError that left
// window.sidebarContext undefined on that page.
function _callGlobal(name) {
    return function (...args) {
        const fn = window[name];
        return typeof fn === 'function' ? fn.apply(window, args) : undefined;
    };
}

window.sidebarContext = {
    initializeContextPanel,
    loadContextData,
    updateContextDisplay:  _callGlobal('updateContextDisplay'),
    updateWorldDisplay:    _callGlobal('updateWorldDisplay'),
    updateScenarioDisplay: _callGlobal('updateScenarioDisplay'),
    updateLorebookDisplay: _callGlobal('updateLorebookDisplay'),
    openWorldSelector:     _callGlobal('openWorldSelector'),
    selectWorld:           _callGlobal('selectWorld'),
    toggleWorldAutoSwitch: _callGlobal('toggleWorldAutoSwitch'),
    openScenarioSelector:  _callGlobal('openScenarioSelector'),
    filterScenarios:       _callGlobal('filterScenarios'),
    selectScenario:        _callGlobal('selectScenario'),
    toggleScenarioAutoSwitch: _callGlobal('toggleScenarioAutoSwitch'),
    createLorebookItem:    _callGlobal('createLorebookItem'),
    removeLorebookFromChat: _callGlobal('removeLorebookFromChat'),
    openLorebookSelector:  _callGlobal('openLorebookSelector'),
    filterLorebooks:       _callGlobal('filterLorebooks'),
    toggleChatLorebook:    _callGlobal('toggleChatLorebook'),
    openLorebookManager:   _callGlobal('openLorebookManager')
};
