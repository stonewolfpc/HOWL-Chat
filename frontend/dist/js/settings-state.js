/**
 * Settings State — Single Source of Truth
 *
 * Professional settings store pattern:
 *   - One in-memory object (currentSettings)
 *   - One backend file (settings.json)  — survives localStorage clears
 *   - One localStorage key (howl-chat-settings) — fallback when backend unavailable
 *   - Every write syncs to both instantly
 *   - UI never stores state independently
 *
 * @module settings-state
 * @version 2.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Default Schema
// ---------------------------------------------------------------------------

function _defaultSettings() {
    return {
        general: {
            theme:             'dark',
            animationSpeed:    1,
            uiScale:           1,
            soundEffects:      true,
            startupSound:      true,
            ttsAutoplay:       true,
            defaultModel:      '',
            defaultSampler:    'balanced',
            autoSaveInterval:  5
        },
        storage: {
            chatAutosave:           true,
            autoArchive:            true,
            archiveDays:            30,
            maxRecentChats:         10,
            chatStorageLocation:    'AppData/Chats',
            contentStorageLocation: 'AppData',
            worldStorageLocation:   'AppData/Worlds',
            scenarioStorageLocation:'AppData/Scenarios',
            characterStorageLocation:'AppData/Characters',
            imageStorageLocation:   'AppData/Images'
        },
        privacy: {
            localOnlyMode:    false,
            disableTelemetry: true,
            permissions:      {}
        },
        audio: {
            asrEnabled:          false,
            ttsEnabled:          false,
            asrEngine:           'auto',
            asrLanguage:         'auto',
            asrDevice:           'auto',
            asrThreads:          4,
            asrBeamSize:         5,
            asrBestOf:           5,
            asrTemperature:      0.0,
            asrTranslateToEnglish: false,
            asrNoTimestamps:     false,
            asrFlashAttention:   true,
            vadEnabled:          false,
            vadThreshold:        0.5,
            vadMinSpeechMs:      250,
            vadMinSilenceMs:     100,
            vadPaddingMs:        30,
            ttsEngine:           'piper',
            ttsVoice:            'en_US-lessac-medium',
            ttsSpeed:            1.0,
            ttsVolume:           0.8,
            ttsNoiseScale:       0.667,
            ttsNoiseW:           0.8,
            ttsSentenceSilence:  0.2,
            ttsAutoplayMessage:  true,
            ttsCacheEnabled:     true,
            ttsCacheSizeMB:      256,
            audioInputDevice:    'default',
            audioOutputDevice:   'default',
            audioSampleRate:     16000,
            audioChannels:       'mono',
            audioBitDepth:       16,
            micSensitivity:      0.5,
            pushToTalk:          true
        },
        imageGen: {
            modelsDir:      '',
            sdServerPath:   '',
            outputDir:      '',
            defaultWidth:   1024,
            defaultHeight:  1024,
            defaultSteps:   20,
            defaultSampler: 'euler_a',
            defaultCfg:     7.0
        },
        sampler: {},     // populated from backend / legacy files
        model:   {}      // populated from backend / legacy files
    };
}

// ---------------------------------------------------------------------------
// Deep Merge
// ---------------------------------------------------------------------------

function deepMerge(target, source) {
    if (typeof source !== 'object' || source === null) return;
    for (const key of Object.keys(source)) {
        if (typeof source[key] === 'object' && source[key] !== null && !Array.isArray(source[key])) {
            if (typeof target[key] !== 'object' || target[key] === null) target[key] = {};
            deepMerge(target[key], source[key]);
        } else {
            target[key] = source[key];
        }
    }
}

// ---------------------------------------------------------------------------
// Mutable State
// ---------------------------------------------------------------------------

let currentSettings = _defaultSettings();

// ---------------------------------------------------------------------------
// Persistence — backend first, localStorage fallback
// ---------------------------------------------------------------------------

async function loadSettings() {
    let loaded = false;

    // 1. Try backend unified settings
    if (window.go?.main?.App?.GetSettings) {
        try {
            const backend = await window.go.main.App.GetSettings();
            if (backend && Object.keys(backend).length > 0) {
                deepMerge(currentSettings, backend);
                loaded = true;
                console.log('[settingsState] loaded from backend');
            }
        } catch (e) {
            console.warn('[settingsState] backend load failed:', e);
        }
    }

    // 2. Fallback to localStorage
    if (!loaded) {
        try {
            const saved = localStorage.getItem('howl-chat-settings');
            if (saved) {
                deepMerge(currentSettings, JSON.parse(saved));
                loaded = true;
                console.log('[settingsState] loaded from localStorage');
            }
        } catch (e) {
            console.error('[settingsState] localStorage load failed:', e);
        }
    }

    return loaded;
}

let saveDebounceTimer = null;

function _writeLocalSettings() {
    try {
        localStorage.setItem('howl-chat-settings', JSON.stringify(currentSettings));
    } catch (e) {
        console.warn('[settingsState] localStorage save failed:', e);
    }
}

async function _writeBackendSettings() {
    const app = window.go?.main?.App;
    if (typeof app?.SetSettings !== 'function') {
        // Never return silently: a missing backend must fail the save path so
        // the caller cannot report "Saved" for a write that never happened.
        throw new Error('settings backend unavailable');
    }
    await app.SetSettings(currentSettings);
}

// saveSettingsNow is the explicit save path. It clears any pending debounce and
// awaits the backend write, so callers that are about to destroy their document
// (the audio-settings iframe) can prove the preference reached disk first.
async function saveSettingsNow() {
    _writeLocalSettings();
    if (saveDebounceTimer) {
        clearTimeout(saveDebounceTimer);
        saveDebounceTimer = null;
    }
    await _writeBackendSettings();
}

async function saveSettings() {
    // 1. Always write to localStorage immediately (cheap, synchronous)
    _writeLocalSettings();

    // 2. Debounce backend sync — batch rapid changes (e.g. slider drags)
    //    into a single call 150ms after the user pauses.
    if (saveDebounceTimer) clearTimeout(saveDebounceTimer);
    saveDebounceTimer = setTimeout(async () => {
        saveDebounceTimer = null;
        try {
            await _writeBackendSettings();
            console.log('[settingsState] saved to backend');
        } catch (e) {
            console.warn('[settingsState] backend save failed:', e);
        }
    }, 150);
}

// Convenience: set a value by dot-path (e.g. 'general.theme') and auto-save
function setSetting(path, value) {
    const keys = path.split('.');
    let obj = currentSettings;
    for (let i = 0; i < keys.length - 1; i++) {
        if (!obj[keys[i]]) obj[keys[i]] = {};
        obj = obj[keys[i]];
    }
    obj[keys[keys.length - 1]] = value;
    saveSettings();
}

// Convenience: get a value by dot-path
function getSetting(path, fallback) {
    const keys = path.split('.');
    let obj = currentSettings;
    for (const key of keys) {
        if (obj == null || typeof obj !== 'object') return fallback;
        obj = obj[key];
    }
    return obj !== undefined ? obj : fallback;
}

function saveAllSettings() {
    saveSettings();
    const overlay = document.getElementById('real-settings-overlay');
    if (overlay) overlay.classList.remove('show');
}

function resetAllSettings() {
    if (!confirm('Reset all settings to default values? This will restore all settings to their original state.')) return;
    currentSettings = _defaultSettings();
    initializeSettingsControls();
    saveSettings();
    if (typeof showNotification === 'function') showNotification('Settings reset to defaults', 'success');
}

// ---------------------------------------------------------------------------
// Init
// ---------------------------------------------------------------------------

(async function _initSettings() {
    await loadSettings();
    // Announce completion so audio-service can re-apply a saved engine after
    // the persisted settings are actually in memory. DOMContentLoaded alone can
    // fire before the backend settings promise resolves.
    window.dispatchEvent(new CustomEvent('howl:settings-loaded'));
})();

// ---------------------------------------------------------------------------
// Public Surface
// ---------------------------------------------------------------------------

window.settingsState = {
    get current() { return currentSettings; },
    loadSettings,
    saveSettings,
    flushSettings: saveSettingsNow,
    set: setSetting,
    get: getSetting,
    resetAllSettings
};
