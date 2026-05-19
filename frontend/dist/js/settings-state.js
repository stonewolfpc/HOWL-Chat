/**
 * Settings State
 *
 * Default settings object, load/save persistence, and full reset.
 * All other settings modules read/write through window.settingsState.
 *
 * @module settings-state
 * @version 1.0.0
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
            // General toggles
            asrEnabled:          false,
            ttsEnabled:          false,

            // ASR Engine Settings
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

            // VAD Settings
            vadEnabled:          false,
            vadThreshold:        0.5,
            vadMinSpeechMs:      250,
            vadMinSilenceMs:     100,
            vadPaddingMs:        30,

            // TTS Engine Settings
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

            // Audio Hardware / I/O
            audioInputDevice:    'default',
            audioOutputDevice:   'default',
            audioSampleRate:     16000,
            audioChannels:       'mono',
            audioBitDepth:       16,
            micSensitivity:      0.5,
            pushToTalk:          true
        },
        imageGen: {
            modelsDir:      '',  // empty = auto-detect next to exe
            sdServerPath:   '',  // empty = auto-detect next to exe
            outputDir:      '',  // empty = auto-detect
            defaultWidth:   1024,
            defaultHeight:  1024,
            defaultSteps:   20,
            defaultSampler: 'euler_a',
            defaultCfg:     7.0
        }
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
// Persistence
// ---------------------------------------------------------------------------

function loadSettings() {
    try {
        const saved = localStorage.getItem('howl-chat-settings');
        if (saved) deepMerge(currentSettings, JSON.parse(saved));
    } catch (error) {
        console.error('Failed to load settings:', error);
    }
}

function saveSettings() {
    try {
        localStorage.setItem('howl-chat-settings', JSON.stringify(currentSettings));
        if (typeof showNotification === 'function') showNotification('Settings saved successfully', 'success');
    } catch (error) {
        console.error('Failed to save settings:', error);
        if (typeof showNotification === 'function') showNotification('Failed to save settings', 'error');
    }
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
// Public Surface
// ---------------------------------------------------------------------------

window.settingsState = {
    get current() { return currentSettings; },
    loadSettings,
    saveSettings,
    resetAllSettings
};
