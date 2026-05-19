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
        });
    };

    const themeSelect = document.getElementById('default-theme');
    if (themeSelect) {
        themeSelect.value = s.general.theme;
        themeSelect.addEventListener('change', e => { s.general.theme = e.target.value; applyTheme(e.target.value); });
    }

    _bindSliderPair('animation-speed', 'animation-speed-value', 'general', 'animationSpeed', 0, 2, applyAnimationSpeed);
    _bindSliderPair('ui-scale',        'ui-scale-value',        'general', 'uiScale',        0.5, 2, applyUIScale);

    bind('sound-effects',    'soundEffects',   'general', null);
    bind('startup-sound',    'startupSound',   'general', null);
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
            }
        });
    }
}

// ---------------------------------------------------------------------------
// Storage
// ---------------------------------------------------------------------------

function initializeStorageSettings() {
    const s = window.settingsState.current;

    const chatAutosave = document.getElementById('chat-autosave');
    if (chatAutosave) {
        chatAutosave.checked = s.storage.chatAutosave;
        chatAutosave.addEventListener('change', e => { s.storage.chatAutosave = e.target.checked; });
    }

    const autoArchive = document.getElementById('auto-archive');
    const archiveDays = document.getElementById('archive-days');
    if (autoArchive && archiveDays) {
        autoArchive.checked  = s.storage.autoArchive;
        archiveDays.value    = s.storage.archiveDays;
        autoArchive.addEventListener('change', e => {
            s.storage.autoArchive = e.target.checked;
            if (typeof window.restartArchiveTimer === 'function') window.restartArchiveTimer();
        });
        archiveDays.addEventListener('input',  e => {
            const v = parseInt(e.target.value);
            if (v >= 7 && v <= 365) {
                s.storage.archiveDays = v;
                if (typeof window.restartArchiveTimer === 'function') window.restartArchiveTimer();
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
            }
        });
    }

    const chatLoc    = document.getElementById('chat-storage-location');
    const contentLoc = document.getElementById('content-storage-location');
    if (chatLoc)    chatLoc.value    = s.storage.chatStorageLocation;
    if (contentLoc) contentLoc.value = s.storage.contentStorageLocation;

    const worldLoc     = document.getElementById('world-storage-location');
    const scenarioLoc  = document.getElementById('scenario-storage-location');
    const characterLoc = document.getElementById('character-storage-location');
    const imageLoc     = document.getElementById('image-storage-location');
    if (worldLoc)     worldLoc.value     = s.storage.worldStorageLocation;
    if (scenarioLoc)  scenarioLoc.value  = s.storage.scenarioStorageLocation;
    if (characterLoc) characterLoc.value = s.storage.characterStorageLocation;
    if (imageLoc)     imageLoc.value     = s.storage.imageStorageLocation;
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
        });
    }

    const telemetry = document.getElementById('disable-telemetry');
    if (telemetry) {
        telemetry.checked = s.privacy.disableTelemetry;
        telemetry.addEventListener('change', e => {
            s.privacy.disableTelemetry = e.target.checked;
            if (typeof applyPrivacySettings === 'function') applyPrivacySettings();
        });
    }

    if (typeof applyPrivacySettings === 'function') applyPrivacySettings();
}

// ---------------------------------------------------------------------------
// Audio
// ---------------------------------------------------------------------------

function initializeAudioSettings() {
    const s = window.settingsState.current;

    const ttsVoice = document.getElementById('default-tts-voice');
    if (ttsVoice) { ttsVoice.value = s.audio.ttsVoice; ttsVoice.addEventListener('change', e => { s.audio.ttsVoice = e.target.value; }); }

    _bindSliderPair('tts-speed',        'tts-speed-value',        'audio', 'ttsSpeed',       0.5, 2);
    _bindSliderPair('tts-volume',       'tts-volume-value',       'audio', 'ttsVolume',      0,   1);
    _bindSliderPair('mic-sensitivity',  'mic-sensitivity-value',  'audio', 'micSensitivity', 0.1, 1);

    _bindCheckbox('tts-autoplay-message', 'audio', 'ttsAutoplayMessage');
    _bindCheckbox('push-to-talk',         'audio', 'pushToTalk');
}

// ---------------------------------------------------------------------------
// Private Helpers
// ---------------------------------------------------------------------------

function _bindCheckbox(id, section, key) {
    const el = document.getElementById(id);
    if (!el) return;
    el.checked = window.settingsState.current[section][key];
    el.addEventListener('change', e => { window.settingsState.current[section][key] = e.target.checked; });
}

function _bindIntInput(id, section, key, min, max) {
    const el = document.getElementById(id);
    if (!el) return;
    el.value = window.settingsState.current[section][key];
    el.addEventListener('change', e => {
        const v = parseInt(e.target.value);
        if (v >= min && v <= max) window.settingsState.current[section][key] = v;
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
    });
    input.addEventListener('input', e => {
        const v = parseFloat(e.target.value);
        if (v >= min && v <= max) { s[section][key] = v; slider.value = v; if (onApply) onApply(v); }
    });
}
