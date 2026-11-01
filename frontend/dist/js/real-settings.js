/**
 * Real Settings — Orchestrator
 *
 * Modal show/close/panel-switch, DOMContentLoaded init,
 * and window.realSettings public API.
 *
 * Depends on (must load first):
 *   settings-state.js    — currentSettings, load/save/reset
 *   settings-actions.js  — applyTheme, backup/cache actions
 *   settings-controls.js — initializeSettingsControls + section inits
 *
 * @module real-settings
 * @version 2.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Modal Lifecycle
// ---------------------------------------------------------------------------

function showRealSettings() {
    const overlay = document.getElementById('real-settings-overlay');
    if (!overlay) return;
    overlay.classList.add('show');
    window.settingsState.loadSettings();
}

function openRealSettings() {
    showRealSettings();
}

function closeRealSettings() {
    const overlay = document.getElementById('real-settings-overlay');
    if (!overlay) return;
    overlay.classList.remove('show');
    window.settingsState.saveSettings();
}

function showSettingsPanel(panelId) {
    document.querySelectorAll('.real-settings-panel').forEach(p => p.classList.remove('active'));
    document.querySelectorAll('.real-settings-nav-item').forEach(i => i.classList.remove('active'));
    document.getElementById(`${panelId}-panel`)?.classList.add('active');
    document.querySelector(`[onclick^="showSettingsPanel('${panelId}')"]`)?.classList.add('active');
}

// ---------------------------------------------------------------------------
// Init
// ---------------------------------------------------------------------------

document.addEventListener('DOMContentLoaded', function () {
    window.settingsState.loadSettings();
    initializeSettingsControls();
});

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

window.realSettings = {
    show:      showRealSettings,
    close:     closeRealSettings,
    showPanel: showSettingsPanel,
    reset:     () => window.settingsState.resetAllSettings(),
    save:      () => window.settingsState.saveSettings()
};
