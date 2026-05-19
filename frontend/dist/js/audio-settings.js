/**
 * Audio Settings — Controls & Modal Lifecycle
 *
 * Wires all audio setting form controls to the live currentSettings object,
 * handles tab switching, modal open/close, device enumeration,
 * and backend status polling for ASR/TTS readiness.
 *
 * @module audio-settings
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

window.AudioSettings = {
    show,
    close,
    showTab,
    refreshAudioDevices,
    saveAudioSettings,
    resetAudioSettings
};

// ---------------------------------------------------------------------------
// Modal Lifecycle
// ---------------------------------------------------------------------------

function _notify(msg, type) {
    try {
        if (typeof window.parent?.showNotification === 'function') {
            window.parent.showNotification(msg, type);
            return;
        }
    } catch (e) {}
    if (typeof showNotification === 'function') {
        showNotification(msg, type);
    } else {
        console.log('[' + (type || 'info') + '] ' + msg);
    }
}

function show(initialTab) {
    const overlay = document.getElementById('audio-settings-overlay');
    if (!overlay) return;
    overlay.classList.add('show');
    window.settingsState.loadSettings();
    _hydrateAllControls();
    _pollBackendStatus();
    if (initialTab) showTab(initialTab);
    document.addEventListener('keydown', _onEsc);
}

function close() {
    const overlay = document.getElementById('audio-settings-overlay');
    if (overlay) overlay.classList.remove('show');
    document.removeEventListener('keydown', _onEsc);
    window.settingsState.saveSettings();
    // Remove parent iframe modal if running inside one
    try {
        const parentModal = window.parent.document.getElementById('audio-settings-modal');
        if (parentModal) parentModal.remove();
    } catch (e) { /* cross-origin or no parent */ }
}

function _onEsc(e) { if (e.key === 'Escape') close(); }

// Auto-init when loaded as standalone page (iframe)
window.addEventListener('DOMContentLoaded', () => {
    const hash = window.location.hash.replace('#', '');
    if (hash && ['general', 'asr', 'tts', 'vad', 'hardware'].includes(hash)) {
        showTab(hash);
    }
    window.settingsState.loadSettings();
    _hydrateAllControls();
    _pollBackendStatus();
});

function showTab(tabId) {
    document.querySelectorAll('.audio-panel').forEach(p => p.classList.remove('active'));
    document.querySelectorAll('.audio-tab').forEach(t => t.classList.remove('active'));
    document.getElementById('panel-' + tabId)?.classList.add('active');
    document.getElementById('tab-' + tabId)?.classList.add('active');
}

// ---------------------------------------------------------------------------
// Hydrate / Wire Controls
// ---------------------------------------------------------------------------

function _hydrateAllControls() {
    const s = window.settingsState.current.audio;

    // -- General --
    _setToggle('asr-enabled',          s.asrEnabled);
    _setToggle('tts-enabled',            s.ttsEnabled);
    _setToggle('tts-autoplay-message',   s.ttsAutoplayMessage);
    _setToggle('push-to-talk',           s.pushToTalk);
    _setSlider('tts-volume',             s.ttsVolume, 0, 1, 0.05);

    // -- ASR --
    _setSelect('asr-engine',             s.asrEngine);
    _setSelect('asr-language',           s.asrLanguage);
    _setSelect('asr-device',             s.asrDevice);
    _setSlider('asr-threads',            s.asrThreads, 1, 16, 1);
    _setSlider('asr-beam-size',          s.asrBeamSize, 1, 10, 1);
    _setSlider('asr-best-of',            s.asrBestOf, 1, 10, 1);
    _setSlider('asr-temperature',        s.asrTemperature, 0, 1, 0.05);
    _setToggle('asr-translate',          s.asrTranslateToEnglish);
    _setToggle('asr-no-timestamps',      s.asrNoTimestamps);
    _setToggle('asr-flash-attn',         s.asrFlashAttention);

    // -- TTS --
    _setSelect('tts-engine',             s.ttsEngine);
    _setSelect('tts-voice',              s.ttsVoice);
    _setSlider('tts-speed',              s.ttsSpeed, 0.5, 2.0, 0.05);
    _setSlider('tts-noise-scale',        s.ttsNoiseScale, 0, 1, 0.01);
    _setSlider('tts-noise-w',            s.ttsNoiseW, 0, 1, 0.01);
    _setSlider('tts-sentence-silence',   s.ttsSentenceSilence, 0, 2, 0.05);
    _setToggle('tts-cache-enabled',      s.ttsCacheEnabled);
    _setSlider('tts-cache-size',         s.ttsCacheSizeMB, 50, 1024, 50);

    // -- VAD --
    _setToggle('vad-enabled',            s.vadEnabled);
    _setSlider('vad-threshold',          s.vadThreshold, 0, 1, 0.05);
    _setSlider('vad-min-speech',         s.vadMinSpeechMs, 0, 1000, 50);
    _setSlider('vad-min-silence',        s.vadMinSilenceMs, 0, 1000, 50);
    _setSlider('vad-padding',            s.vadPaddingMs, 0, 500, 10);

    // -- Hardware --
    _setSelect('audio-input-device',     s.audioInputDevice);
    _setSelect('audio-output-device',    s.audioOutputDevice);
    _setSelect('audio-sample-rate',      s.audioSampleRate);
    _setSelect('audio-channels',         s.audioChannels);
    _setSelect('audio-bit-depth',        s.audioBitDepth);
    _setSlider('mic-sensitivity',        s.micSensitivity, 0, 1, 0.05);

    // Live wiring (two-way binding)
    _wireToggle('asr-enabled',          'asrEnabled');
    _wireToggle('tts-enabled',          'ttsEnabled');
    _wireToggle('tts-autoplay-message', 'ttsAutoplayMessage');
    _wireToggle('push-to-talk',         'pushToTalk');
    _wireSlider('tts-volume',           'ttsVolume');

    _wireSelect('asr-engine',           'asrEngine');
    _wireSelect('asr-language',         'asrLanguage');
    _wireSelect('asr-device',           'asrDevice');
    _wireSlider('asr-threads',          'asrThreads');
    _wireSlider('asr-beam-size',          'asrBeamSize');
    _wireSlider('asr-best-of',            'asrBestOf');
    _wireSlider('asr-temperature',        'asrTemperature');
    _wireToggle('asr-translate',          'asrTranslateToEnglish');
    _wireToggle('asr-no-timestamps',      'asrNoTimestamps');
    _wireToggle('asr-flash-attn',         'asrFlashAttention');

    _wireSelect('tts-engine',           'ttsEngine');
    _wireSelect('tts-voice',            'ttsVoice');
    _wireSlider('tts-speed',            'ttsSpeed');
    _wireSlider('tts-noise-scale',      'ttsNoiseScale');
    _wireSlider('tts-noise-w',          'ttsNoiseW');
    _wireSlider('tts-sentence-silence', 'ttsSentenceSilence');
    _wireToggle('tts-cache-enabled',    'ttsCacheEnabled');
    _wireSlider('tts-cache-size',       'ttsCacheSizeMB');

    _wireToggle('vad-enabled',          'vadEnabled');
    _wireSlider('vad-threshold',        'vadThreshold');
    _wireSlider('vad-min-speech',       'vadMinSpeechMs');
    _wireSlider('vad-min-silence',      'vadMinSilenceMs');
    _wireSlider('vad-padding',          'vadPaddingMs');

    _wireSelect('audio-sample-rate',    'audioSampleRate');
    _wireSelect('audio-channels',       'audioChannels');
    _wireSelect('audio-bit-depth',      'audioBitDepth');
    _wireSlider('mic-sensitivity',      'micSensitivity');
    _wireSelect('audio-input-device',   'audioInputDevice');
    _wireSelect('audio-output-device',  'audioOutputDevice');
}

// ---------------------------------------------------------------------------
// Low-level Helpers
// ---------------------------------------------------------------------------

function _setToggle(id, val) {
    const el = document.getElementById(id);
    if (el) el.checked = val;
}

function _setSelect(id, val) {
    const el = document.getElementById(id);
    if (el) el.value = String(val);
}

function _setSlider(id, val, min, max, step) {
    const slider = document.getElementById(id);
    const input  = document.getElementById(id + '-value');
    if (slider) { slider.min = min; slider.max = max; slider.step = step; slider.value = val; }
    if (input)  { input.min = min; input.max = max; input.step = step; input.value = val; }
}

function _wireToggle(id, key) {
    const el = document.getElementById(id);
    if (!el || el._wired) return;
    el._wired = true;
    el.addEventListener('change', e => { window.settingsState.current.audio[key] = e.target.checked; });
}

function _wireSelect(id, key) {
    const el = document.getElementById(id);
    if (!el || el._wired) return;
    el._wired = true;
    el.addEventListener('change', e => { window.settingsState.current.audio[key] = e.target.value; });
}

function _wireSlider(id, key) {
    const slider = document.getElementById(id);
    const input  = document.getElementById(id + '-value');
    if (!slider || slider._wired) return;
    slider._wired = true;
    slider.addEventListener('input', e => {
        const v = parseFloat(e.target.value);
        window.settingsState.current.audio[key] = v;
        if (input) input.value = v;
    });
    if (input && !input._wired) {
        input._wired = true;
        input.addEventListener('input', e => {
            const v = parseFloat(e.target.value);
            const min = parseFloat(slider.min), max = parseFloat(slider.max);
            if (v >= min && v <= max) {
                window.settingsState.current.audio[key] = v;
                slider.value = v;
            }
        });
    }
}

// ---------------------------------------------------------------------------
// Device Enumeration
// ---------------------------------------------------------------------------

async function refreshAudioDevices() {
    try {
        if (!navigator.mediaDevices || !navigator.mediaDevices.enumerateDevices) {
            _notify('Device enumeration not supported in this browser', 'warning');
            return;
        }
        const devices = await navigator.mediaDevices.enumerateDevices();
        const inputs  = devices.filter(d => d.kind === 'audioinput');
        const outputs = devices.filter(d => d.kind === 'audiooutput');

        const inSelect  = document.getElementById('audio-input-device');
        const outSelect = document.getElementById('audio-output-device');

        if (inSelect) {
            const oldVal = inSelect.value;
            inSelect.innerHTML = '<option value="default">System Default</option>';
            inputs.forEach((d, i) => {
                const opt = document.createElement('option');
                opt.value = d.deviceId;
                opt.textContent = d.label || ('Microphone ' + (i + 1));
                inSelect.appendChild(opt);
            });
            inSelect.value = oldVal;
        }
        if (outSelect) {
            const oldVal = outSelect.value;
            outSelect.innerHTML = '<option value="default">System Default</option>';
            outputs.forEach((d, i) => {
                const opt = document.createElement('option');
                opt.value = d.deviceId;
                opt.textContent = d.label || ('Speaker ' + (i + 1));
                outSelect.appendChild(opt);
            });
            outSelect.value = oldVal;
        }
        _notify('Audio device list refreshed', 'success');
    } catch (err) {
        console.error('Device refresh error:', err);
        _notify('Failed to refresh audio devices', 'error');
    }
}

// ---------------------------------------------------------------------------
// Backend Status Polling
// ---------------------------------------------------------------------------

async function _pollBackendStatus() {
    const asrText = document.getElementById('asr-status-text');
    const ttsText = document.getElementById('tts-status-text');
    if (!window.go?.main?.App?.GetAudioStatus) {
        if (asrText) asrText.textContent = 'Status: backend not available (dev mode)';
        if (ttsText) ttsText.textContent = 'Status: backend not available (dev mode)';
        return;
    }
    try {
        const status = await window.go.main.App.GetAudioStatus();
        if (asrText) asrText.textContent = status.asr_ready ? 'Status: ready' : 'Status: not initialized (enable ASR below)';
        if (ttsText) ttsText.textContent = status.tts_ready ? 'Status: ready' : 'Status: not initialized (enable TTS below)';
    } catch (err) {
        if (asrText) asrText.textContent = 'Status: error checking backend';
        if (ttsText) ttsText.textContent = 'Status: error checking backend';
    }
}

// ---------------------------------------------------------------------------
// Save / Reset
// ---------------------------------------------------------------------------

function saveAudioSettings() {
    window.settingsState.saveSettings();
    close();
    _notify('Audio settings saved', 'success');
}

function resetAudioSettings() {
    if (!confirm('Reset all audio settings to their default values?')) return;
    const defs = {
        asrEnabled: false, ttsEnabled: false, asrEngine: 'auto', asrLanguage: 'auto',
        asrDevice: 'auto', asrThreads: 4, asrBeamSize: 5, asrBestOf: 5,
        asrTemperature: 0.0, asrTranslateToEnglish: false, asrNoTimestamps: false,
        asrFlashAttention: true, vadEnabled: false, vadThreshold: 0.5,
        vadMinSpeechMs: 250, vadMinSilenceMs: 100, vadPaddingMs: 30,
        ttsEngine: 'piper', ttsVoice: 'en_US-lessac-medium', ttsSpeed: 1.0,
        ttsVolume: 0.8, ttsNoiseScale: 0.667, ttsNoiseW: 0.8,
        ttsSentenceSilence: 0.2, ttsAutoplayMessage: true, ttsCacheEnabled: true,
        ttsCacheSizeMB: 256, audioInputDevice: 'default', audioOutputDevice: 'default',
        audioSampleRate: 16000, audioChannels: 'mono', audioBitDepth: 16,
        micSensitivity: 0.5, pushToTalk: true
    };
    Object.assign(window.settingsState.current.audio, defs);
    _hydrateAllControls();
    window.settingsState.saveSettings();
    _notify('Audio settings reset to defaults', 'info');
}

// Global aliases for HTML onclick handlers
function closeAudioSettings() { close(); }
function showAudioTab(tabId) { AudioSettings.showTab(tabId); }
