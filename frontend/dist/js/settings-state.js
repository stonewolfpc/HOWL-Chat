/**
 * Settings State — the renderer's single seam to the settings authority
 *
 * Ownership:
 *   - `runtime.Settings` (GetSettings/SetSettings) is authoritative. This
 *     module is the only renderer code that talks to it for settings.
 *   - `current` is the last backend-acknowledged document. UI code may mutate
 *     it as a draft (the form bindings do), but only `saveSettings()` commits:
 *     the whole document goes to SetSettings FIRST, and only after it is
 *     acknowledged is the versioned localStorage cache rewritten.
 *   - The cache is derived state. It carries a version and is always rebuilt
 *     from the backend on load; a stale, legacy or missing cache can never
 *     override backend values. It is only read when the backend is unreachable.
 *   - Genuinely UI-only preferences (open panels, last viewed tab) belong in
 *     ui-prefs.js, never here and never in the backend settings document.
 *
 * @module settings-state
 * @version 3.0.0
 */

'use strict';

const SETTINGS_CACHE_KEY = 'howl-chat-settings';
const SETTINGS_CACHE_VERSION = 2;

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
        sampler: {},     // live sampler section from the authority
        model:   {}      // live model section from the authority
    };
}

// ---------------------------------------------------------------------------
// Deep Merge / Clone
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

function _clone(value) {
    return JSON.parse(JSON.stringify(value));
}

// Replace the contents of the stable currentSettings object without changing
// its identity: the form bindings captured a reference to it at init time.
function _replaceContents(target, source) {
    for (const key of Object.keys(target)) delete target[key];
    Object.assign(target, _clone(source));
}

// ---------------------------------------------------------------------------
// Mutable State
// ---------------------------------------------------------------------------

let currentSettings = _defaultSettings();
let committedSettings = _clone(currentSettings);
let settingsSource = 'defaults';
let writeInFlight = false;
let pendingWrite = null;
let writeSeq = 0;
let waiters = [];

// ---------------------------------------------------------------------------
// Derived cache — versioned, rebuildable, written only after an acknowledged
// backend write (or after an authoritative backend read).
// ---------------------------------------------------------------------------

function _readCache() {
    try {
        const raw = localStorage.getItem(SETTINGS_CACHE_KEY);
        if (!raw) return null;
        const parsed = JSON.parse(raw);
        if (!parsed || parsed.version !== SETTINGS_CACHE_VERSION ||
            typeof parsed.settings !== 'object' || parsed.settings === null) {
            console.warn('[settingsState] ignoring settings cache with an unknown version');
            return null;
        }
        return parsed.settings;
    } catch (e) {
        console.warn('[settingsState] localStorage read failed:', e);
        return null;
    }
}

function _writeCache(settings) {
    try {
        localStorage.setItem(SETTINGS_CACHE_KEY, JSON.stringify({
            version: SETTINGS_CACHE_VERSION,
            savedAt: new Date().toISOString(),
            settings
        }));
    } catch (e) {
        console.warn('[settingsState] localStorage save failed:', e);
    }
}

// ---------------------------------------------------------------------------
// Backend write queue — serialized, last-write-wins, snapshot per save
// ---------------------------------------------------------------------------

async function _writeBackendSettings(snapshot) {
    const app = window.go?.main?.App;
    if (typeof app?.SetSettings !== 'function') {
        // Never return silently: a missing backend must fail the save path so
        // the caller cannot report "Saved" for a write that never happened.
        throw new Error('settings backend unavailable');
    }
    await app.SetSettings(snapshot);
}

function _commit(snapshot) {
    _replaceContents(currentSettings, snapshot);
    committedSettings = _clone(snapshot);
    settingsSource = 'backend';
    // Cache is derived from what the authority acknowledged, never from the
    // draft the user was editing.
    _writeCache(committedSettings);
}

function _resolveWaiters(upToSeq, ok) {
    const settled = [];
    waiters = waiters.filter(w => {
        if (w.seq <= upToSeq) { settled.push(w); return false; }
        return true;
    });
    settled.forEach(w => w.resolve(ok));
}

function _notifyWriteFailure(error) {
    console.error('[settingsState] settings write failed; reverted to the last saved document:', error);
    if (typeof window.dispatchEvent === 'function' && typeof CustomEvent === 'function') {
        window.dispatchEvent(new CustomEvent('howl:settings-write-failed', {
            detail: { error: String(error?.message || error) }
        }));
    }
    if (typeof window.showNotification === 'function') {
        window.showNotification('Settings were not saved: ' + (error?.message || error), 'error');
    }
}

async function _pumpWrites() {
    if (writeInFlight) return;
    writeInFlight = true;
    while (pendingWrite) {
        const job = pendingWrite;
        pendingWrite = null;
        try {
            await _writeBackendSettings(job.snapshot);
            _commit(job.snapshot);
            _resolveWaiters(job.seq, true);
        } catch (error) {
            // A failed write must leave no trace: the in-memory document goes
            // back to the last acknowledged one and the cache is untouched.
            _replaceContents(currentSettings, committedSettings);
            pendingWrite = null;
            _resolveWaiters(Infinity, false);
            _notifyWriteFailure(error);
        }
    }
    writeInFlight = false;
}

function _enqueueSnapshot(snapshot) {
    const seq = ++writeSeq;
    pendingWrite = { snapshot, seq };
    const settled = new Promise(resolve => waiters.push({ seq, resolve }));
    _pumpWrites();
    return settled;
}

async function _drainWrites() {
    while (writeInFlight || pendingWrite) {
        await new Promise(resolve => setTimeout(resolve, 0));
    }
}

// ---------------------------------------------------------------------------
// Load — backend first, cache only when the backend is unreachable
// ---------------------------------------------------------------------------

async function loadSettings() {
    await _drainWrites();

    const app = window.go?.main?.App;
    if (typeof app?.GetSettings === 'function') {
        try {
            const backend = await app.GetSettings();
            if (backend && typeof backend === 'object') {
                const next = _defaultSettings();
                deepMerge(next, backend);
                _replaceContents(currentSettings, next);
                committedSettings = _clone(next);
                settingsSource = 'backend';
                // Refresh the derived cache from the authority. This is the
                // rebuild path for an empty, stale or conflicting cache.
                _writeCache(committedSettings);
                console.log('[settingsState] loaded from backend');
                return true;
            }
        } catch (e) {
            console.warn('[settingsState] backend load failed:', e);
        }
    }

    const cached = _readCache();
    if (cached) {
        const next = _defaultSettings();
        deepMerge(next, cached);
        _replaceContents(currentSettings, next);
        committedSettings = _clone(next);
        settingsSource = 'cache';
        console.warn('[settingsState] backend unavailable; using the versioned cache');
        return true;
    }

    const next = _defaultSettings();
    _replaceContents(currentSettings, next);
    committedSettings = _clone(next);
    settingsSource = 'defaults';
    return false;
}

// ---------------------------------------------------------------------------
// Save — backend first, cache only after acknowledgement
// ---------------------------------------------------------------------------

function saveSettings() {
    return _enqueueSnapshot(_clone(currentSettings));
}

// saveSettingsNow is the explicit flush path: it awaits the backend write, so
// a caller that is about to destroy its document can prove the preference
// reached disk first. It rejects when the write fails, unlike the fire-and-
// forget form-binding path, so a caller that reports "Saved" cannot lie.
function saveSettingsNow() {
    return _enqueueSnapshot(_clone(currentSettings)).then(ok => {
        if (!ok) throw new Error('settings backend rejected the write');
        return true;
    });
}

// Convenience: set a value by dot-path (e.g. 'general.theme') and persist it.
// The value is not committed to `current` until the backend acknowledges it.
function setSetting(path, value) {
    const draft = _clone(currentSettings);
    const keys = path.split('.');
    let obj = draft;
    for (let i = 0; i < keys.length - 1; i++) {
        if (!obj[keys[i]] || typeof obj[keys[i]] !== 'object') obj[keys[i]] = {};
        obj = obj[keys[i]];
    }
    obj[keys[keys.length - 1]] = value;
    return _enqueueSnapshot(draft);
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
    const pending = saveSettings();
    const overlay = document.getElementById('real-settings-overlay');
    if (overlay) overlay.classList.remove('show');
    return pending;
}

async function resetAllSettings() {
    if (typeof confirm === 'function' && !confirm('Reset all settings to default values? This restores application preferences, including image-generation settings. Chats, characters, worlds, scenarios, lore, memories, registered libraries, and model files are preserved.')) return false;

    // A reset must not race a queued edit: settle it first, or the old draft
    // could land after the defaults and silently undo the reset.
    await _drainWrites();

    const previous = _clone(committedSettings);
    const defaults = _defaultSettings();
    const app = window.go?.main?.App;

    try {
        // Backend first. The reset also clears the separate image-settings
        // store; the SetSettings fallback keeps older builds functional.
        if (typeof app?.ResetApplicationSettings === 'function') {
            await app.ResetApplicationSettings(defaults);
        } else if (typeof app?.SetSettings === 'function') {
            await app.SetSettings(defaults);
        } else {
            throw new Error('settings backend unavailable');
        }
        _commit(defaults);
        if (typeof initializeSettingsControls === 'function') initializeSettingsControls();
        if (typeof showNotification === 'function') showNotification('Settings reset to defaults', 'success');
        return true;
    } catch (e) {
        _replaceContents(currentSettings, previous);
        committedSettings = _clone(previous);
        if (typeof initializeSettingsControls === 'function') initializeSettingsControls();
        console.error('[settingsState] reset failed:', e);
        if (typeof showNotification === 'function') showNotification('Failed to reset settings: ' + e.message, 'error');
        return false;
    }
}

// ---------------------------------------------------------------------------
// Public Surface
// ---------------------------------------------------------------------------

window.settingsState = {
    get current() { return currentSettings; },
    get source() { return settingsSource; },
    cacheKey: SETTINGS_CACHE_KEY,
    cacheVersion: SETTINGS_CACHE_VERSION,
    loadSettings,
    saveSettings,
    saveSettingsNow,
    flushSettings: saveSettingsNow,
    set: setSetting,
    get: getSetting,
    resetAllSettings
};

// ---------------------------------------------------------------------------
// Init
// ---------------------------------------------------------------------------

window.settingsState.ready = (async function _initSettings() {
    await loadSettings();
    // Announce completion so audio-service can re-apply a saved engine after
    // the persisted settings are actually in memory. DOMContentLoaded alone can
    // fire before the backend settings promise resolves.
    window.dispatchEvent(new CustomEvent('howl:settings-loaded'));
    return true;
})();
