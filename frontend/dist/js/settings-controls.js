/**
 * Settings Controls
 *
 * Wires all DOM form controls in the settings modal to the
 * live currentSettings object via event listeners.
 *
 * Depends on: settings-state.js, settings-actions.js
 *
 * @module settings-controls
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Orchestrator
// ---------------------------------------------------------------------------

function initializeSettingsControls() {
    initializeGeneralSettings();
    initializeStorageSettings();
    initializePrivacySettings();
    initializeAudioSettings();
}

// ---------------------------------------------------------------------------
// General
// ---------------------------------------------------------------------------

function initializeGeneralSettings() {
    const s = window.settingsState.current;

    const bind = (id, key, section, parse, extra) => {
        const el = document.getElementById(id);
        if (!el) return;
        el[el.type === 'checkbox' ? 'checked' : 'value'] = s[section][key];
        el.addEventListener(el.type === 'checkbox' ? 'change' : 'input', e => {
            s[section][key] = parse ? parse(e.target.value) : (el.type === 'checkbox' ? e.target.checked : e.target.value);
            if (extra) extra(s[section][key]);
            if (window.settingsState?.saveSettings) window.settingsState.saveSettings();
        });
    };

    const themeSelect = document.getElementById('default-theme');
    if (themeSelect) {
        themeSelect.value = s.general.theme;
        themeSelect.addEventListener('change', e => {
            s.general.theme = e.target.value;
            applyTheme(e.target.value);
            if (window.settingsState?.saveSettings) window.settingsState.saveSettings();
        });
    }

    _bindSliderPair('animation-speed', 'animation-speed-value', 'general', 'animationSpeed', 0, 2, applyAnimationSpeed);

    bind('sound-effects',    'soundEffects',   'general', null);
    bind('tts-autoplay',     'ttsAutoplay',    'general', null);
    bind('default-model',    'defaultModel',   'general', null);
    bind('default-sampler',  'defaultSampler', 'general', null);

    const autoSave = document.getElementById('auto-save-interval');
    if (autoSave) {
        autoSave.value = s.general.autoSaveInterval;
        autoSave.addEventListener('change', e => {
            const v = parseInt(e.target.value);
            if (v >= 1 && v <= 60) {
                s.general.autoSaveInterval = v;
                if (typeof window.restartAutosaveTimer === 'function') window.restartAutosaveTimer();
                if (window.settingsState?.saveSettings) window.settingsState.saveSettings();
            }
        });
    }
}

// ---------------------------------------------------------------------------
// Storage
// ---------------------------------------------------------------------------

function initializeStorageSettings() {
    const s = window.settingsState.current;

    const save = () => { if (window.settingsState?.saveSettings) window.settingsState.saveSettings(); };

    const chatAutosave = document.getElementById('chat-autosave');
    if (chatAutosave) {
        chatAutosave.checked = s.storage.chatAutosave;
        chatAutosave.addEventListener('change', e => { s.storage.chatAutosave = e.target.checked; save(); });
    }

    const autoArchive = document.getElementById('auto-archive');
    const archiveDays = document.getElementById('archive-days');
    if (autoArchive && archiveDays) {
        autoArchive.checked  = s.storage.autoArchive;
        archiveDays.value    = s.storage.archiveDays;
        autoArchive.addEventListener('change', e => {
            s.storage.autoArchive = e.target.checked;
            if (typeof window.restartArchiveTimer === 'function') window.restartArchiveTimer();
            save();
        });
        archiveDays.addEventListener('input',  e => {
            const v = parseInt(e.target.value);
            if (v >= 7 && v <= 365) {
                s.storage.archiveDays = v;
                if (typeof window.restartArchiveTimer === 'function') window.restartArchiveTimer();
                save();
            }
        });
    }

    const maxRecent = document.getElementById('max-recent-chats');
    if (maxRecent) {
        maxRecent.value = s.storage.maxRecentChats;
        maxRecent.addEventListener('change', e => {
            const v = parseInt(e.target.value);
            if (v >= 5 && v <= 50) {
                s.storage.maxRecentChats = v;
                if (typeof renderSessionList === 'function') renderSessionList();
                save();
            }
        });
    }

    // Storage paths are fixed in this build (AppData\HOWL_Chat): the backend owns
    // them and nothing reads these keys, so the old six location cards with a
    // broken webkitdirectory picker were removed rather than left pretending.
}

// ---------------------------------------------------------------------------
// Privacy
// ---------------------------------------------------------------------------

function initializePrivacySettings() {
    const s = window.settingsState.current;

    const localOnly = document.getElementById('local-only-mode');
    if (localOnly) {
        localOnly.checked = s.privacy.localOnlyMode;
        localOnly.addEventListener('change', e => {
            s.privacy.localOnlyMode = e.target.checked;
            if (typeof applyPrivacySettings === 'function') applyPrivacySettings();
            if (window.settingsState?.saveSettings) window.settingsState.saveSettings();
        });
    }

    if (typeof applyPrivacySettings === 'function') applyPrivacySettings();
}

// ---------------------------------------------------------------------------
// Audio
// ---------------------------------------------------------------------------

function initializeAudioSettings() {
    const s = window.settingsState.current;

    // Voice selection belongs to the Audio settings page, which populates the
    // real Piper voice list. The old select here offered invented voices
    // (male-1/female-2) that could never match a stored voice id.
    _bindSliderPair('tts-speed',        'tts-speed-value',        'audio', 'ttsSpeed',       0.5, 2);
    _bindSliderPair('tts-volume',       'tts-volume-value',       'audio', 'ttsVolume',      0,   1);
    _bindSliderPair('mic-sensitivity',  'mic-sensitivity-value',  'audio', 'micSensitivity', 0.1, 1);

    _bindCheckbox('push-to-talk', 'audio', 'pushToTalk');
}

// ---------------------------------------------------------------------------
// Private Helpers
// ---------------------------------------------------------------------------

function _bindCheckbox(id, section, key) {
    const el = document.getElementById(id);
    if (!el) return;
    el.checked = window.settingsState.current[section][key];
    el.addEventListener('change', e => {
        window.settingsState.current[section][key] = e.target.checked;
        if (window.settingsState?.saveSettings) window.settingsState.saveSettings();
    });
}

function _bindIntInput(id, section, key, min, max) {
    const el = document.getElementById(id);
    if (!el) return;
    el.value = window.settingsState.current[section][key];
    el.addEventListener('change', e => {
        const v = parseInt(e.target.value);
        if (v >= min && v <= max) {
            window.settingsState.current[section][key] = v;
            if (window.settingsState?.saveSettings) window.settingsState.saveSettings();
        }
    });
}

function _bindSliderPair(sliderId, valueId, section, key, min, max, onApply) {
    const slider = document.getElementById(sliderId);
    const input  = document.getElementById(valueId);
    if (!slider || !input) return;
    const s = window.settingsState.current;
    slider.value = s[section][key];
    input.value  = s[section][key];
    slider.addEventListener('input', e => {
        const v = parseFloat(e.target.value);
        s[section][key] = v; input.value = v;
        if (onApply) onApply(v);
        if (window.settingsState?.saveSettings) window.settingsState.saveSettings();
    });
    input.addEventListener('input', e => {
        const v = parseFloat(e.target.value);
        if (v >= min && v <= max) { s[section][key] = v; slider.value = v; if (onApply) onApply(v); }
        if (window.settingsState?.saveSettings) window.settingsState.saveSettings();
    });
}
