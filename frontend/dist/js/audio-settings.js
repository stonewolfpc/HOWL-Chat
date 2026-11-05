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
    retryAudioActivation,
    resetAudioSettings
};

let _hydrated = false;
let _hydrationPromise = null;
let _saving = false;
let _lastOutcome = null;

const REMOTE_ASR_ENGINES = new Set(['qwen3-asr-0.6b', 'qwen3-asr-1.7b']);

// Tri-state: true = activation active, false = authority says idle,
// null = unknown. Unknown is never treated as idle, and no stale value is kept
// across a failed read.
let _applyInProgress = null;

const RESTORED_LOCK_MESSAGE =
    'An activation was started in an earlier page of this app and its outcome cannot be ' +
    'established here. A request that has not yet reached the audio backend cannot be ruled ' +
    'out, so this panel will not start a second activation. Close HOWL Chat completely and ' +
    'reopen it to recover.';

// The host is a parent-document singleton. The panel must read it through the
// same-origin parent; window.howlBridgeHost is undefined inside the iframe.
function _hostBridge() {
    try {
        if (window.howlBridgeHost) return window.howlBridgeHost;
        if (window.parent && window.parent !== window) {
            return window.parent.howlBridgeHost || null;
        }
    } catch (e) { /* cross-origin or detached */ }
    return null;
}

function _restoredLockHeld() {
    try {
        const host = _hostBridge();
        if (!host || typeof host.activation !== 'function') return false;
        const a = host.activation();
        return !!(a && a.held && a.restored);
    } catch (e) {
        return false;
    }
}

function _showRestoredLockState() {
    _lastOutcome = { saved: false, applied: false };
    _setOutcome('Not saved', 'State unknown — restart required', 'Readiness unchecked',
        RESTORED_LOCK_MESSAGE, 'warn');
    _setSaveEnabled(false);
    _setRetryVisible(false);
}

function _bridgeReady() {
    const bridge = window.__howlBridgeIframe;
    if (!bridge) return true;                 // native window.go path
    return bridge.mode === 'ready';
}

function _activationBusy() {
    if (_applyInProgress !== false) return true;   // true or unknown
    try {
        const host = _hostBridge();
        return !!(host && typeof host.isActivationLocked === 'function' && host.isActivationLocked());
    } catch (e) {
        return true;                          // unknown: fail closed
    }
}

// apply_in_progress is authoritative and must be re-read before deciding that
// activation is allowed. Only an explicit boolean is evidence; a failed or
// malformed read resets to unknown rather than preserving a stale value.
async function _refreshApplyInProgress() {
    const app = window.go?.main?.App;
    _applyInProgress = null;
    if (typeof app?.GetAudioStatus !== 'function') return null;
    try {
        const status = await app.GetAudioStatus();
        _applyInProgress = (status && typeof status.apply_in_progress === 'boolean')
            ? status.apply_in_progress
            : null;
    } catch (e) {
        _applyInProgress = null;
    }
    return _applyInProgress;
}

window.addEventListener('howl:bridge-ready', () => {
    _setSaveEnabled(_hydrated);
    if (_restoredLockHeld()) _showRestoredLockState();
});
window.addEventListener('howl:bridge-unavailable', () => { _setSaveEnabled(false); });
// A settlement notification is not evidence of idleness: it only prompts a
// fresh authoritative read.
window.addEventListener('howl:activation-settled', () => {
    _applyInProgress = null;
    _refreshApplyInProgress().then(() => {
        if (_restoredLockHeld()) { _showRestoredLockState(); return; }
        _setRetryVisible(true);
    });
});

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
    if (initialTab) showTab(initialTab);
    document.removeEventListener('keydown', _onEsc);
    document.addEventListener('keydown', _onEsc);
    _ensureHydrated().then(() => {
        _pollBackendStatus();
        populateTTSDropdown();
    });
}

// Closing the dialog is a save action, matching the existing auto-save
// behavior. Once a save and its apply attempt have completed, the close button
// can dismiss the dialog without writing again.
async function close() {
    if (_saving) return;

    // A restored lock has no resolvable outcome. Dismissal must not write,
    // apply, retry, or release anything; the lock outlives the panel.
    if (_restoredLockHeld()) {
        _finishClose();
        return;
    }

    if (_lastOutcome?.saved && _lastOutcome?.applied) {
        _finishClose();
        return;
    }
    await _ensureHydrated();
    if (!_hydrated) {
        _setOutcome('Not saved', 'Not applied', 'Readiness unchecked',
            'Saved audio settings could not be loaded, so closing without overwriting them is not available yet.', 'error');
        _finishClose();
        return;
    }
    return saveAudioSettings();
}

function _finishClose() {
    const overlay = document.getElementById('audio-settings-overlay');
    if (overlay) overlay.classList.remove('show');
    document.removeEventListener('keydown', _onEsc);
    try {
        window.parent.document.dispatchEvent(new window.parent.CustomEvent('howl:audio-settings-closed'));
    } catch (e) { /* cross-origin or no parent */ }
    // Remove parent iframe modal if running inside one
    try {
        const parentModal = window.parent.document.getElementById('audio-settings-modal');
        if (parentModal) parentModal.remove();
    } catch (e) { /* cross-origin or no parent */ }
}

function _onEsc(e) { if (e.key === 'Escape') { e.preventDefault(); close(); } }

// Auto-init when loaded as standalone page (iframe)
window.addEventListener('DOMContentLoaded', () => {
    const hash = window.location.hash.replace('#', '');
    if (hash && ['general', 'asr', 'tts', 'vad', 'hardware'].includes(hash)) {
        showTab(hash);
    }
    document.addEventListener('keydown', _onEsc);
    _ensureHydrated().then(() => {
        _pollBackendStatus();
        populateTTSDropdown();
    });
});

function _ensureHydrated() {
    if (_hydrationPromise) return _hydrationPromise;
    _hydrationPromise = _hydrateSettings();
    return _hydrationPromise;
}

async function _hydrateSettings() {
    _setSaveEnabled(false);
    _setRetryVisible(false);
    _setOutcome('Not saved', 'Not applied', 'Readiness unchecked', 'Loading saved audio settings…', 'busy');
    try {
        if (!window.settingsState?.loadSettings) {
            throw new Error('settings state is unavailable');
        }
        await window.settingsState.loadSettings();
        _hydrateAllControls();
        _hydrated = true;
        _setOutcome('Not saved', 'Not applied', 'Readiness unchecked', '', '');
    } catch (err) {
        _hydrated = false;
        console.error('[audio-settings] hydration failed:', err);
        _setOutcome('Not saved', 'Not applied', 'Readiness unchecked',
            'Saved audio settings could not be loaded. Save stays disabled rather than overwriting them with defaults.', 'error');
    } finally {
        _setSaveEnabled(_hydrated);
        if (_restoredLockHeld()) _showRestoredLockState();
    }
}

function _setSaveEnabled(enabled) {
    const button = document.getElementById('audio-settings-save');
    // Save is available only after hydration, a completed bridge handshake and
    // no unresolved restored activation.
    if (button) button.disabled = !(enabled && _bridgeReady()) || _restoredLockHeld();
}

function _setRetryVisible(visible) {
    const button = document.getElementById('audio-retry-activation');
    if (!button) return;
    if (_restoredLockHeld()) {
        button.hidden = true;
        button.disabled = true;
        return;
    }
    button.hidden = !visible;
    // Retry must not start a second activation while one may still be running.
    button.disabled = !visible || _activationBusy();
}

function _setOutcome(saved, applied, readiness, detail, tone) {
    const values = [
        ['audio-save-state', saved, tone === 'busy' ? 'busy' : (saved === 'Saved' ? 'ok' : (saved === 'Not saved' ? '' : 'warn'))],
        ['audio-apply-state', applied, applied === 'Applied' ? 'ok' : (applied === 'Apply failed' ? 'error' : (tone === 'busy' ? 'busy' : ''))],
        ['audio-ready-state', readiness, readiness === 'Ready (local initialized)' ? 'ok' : (readiness === 'Readiness unchecked' || !readiness ? '' : 'warn')]
    ];
    values.forEach(([id, text, state]) => {
        const el = document.getElementById(id);
        if (!el) return;
        el.textContent = text;
        el.classList.remove('ok', 'warn', 'error', 'busy');
        if (state) el.classList.add(state);
    });
    const detailEl = document.getElementById('audio-outcome-detail');
    if (detailEl) {
        detailEl.textContent = detail || '';
        detailEl.hidden = !detail;
    }
}

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
        _applyInProgress = (status && typeof status.apply_in_progress === 'boolean')
            ? status.apply_in_progress
            : null;
        if (asrText) asrText.textContent = _asrStatusText(status);
        if (ttsText) ttsText.textContent = status.tts_ready ? 'Status: ready' : 'Status: not initialized (enable TTS below)';
        _renderCapabilities(status, asrText);
    } catch (err) {
        _applyInProgress = null;
        if (asrText) asrText.textContent = 'Status: error checking backend';
        if (ttsText) ttsText.textContent = 'Status: error checking backend';
    }
}

function _asrStatusText(status) {
    const scope = status?.asr_readiness;
    if (scope === 'local-initialized' && status.asr_ready) {
        return 'Status: ready (local engine initialized; transcription not verified)';
    }
    const selectedEngine = window.settingsState?.current?.audio?.asrEngine;
    if (scope === 'remote-unverified' || REMOTE_ASR_ENGINES.has(selectedEngine)) {
        return 'Status: remote engine applied — availability not verified';
    }
    if (scope === 'disabled') return 'Status: disabled';
    if (scope === 'unconfigured') return 'Status: not ready — choose an explicit engine';
    if (scope === 'not-ready') {
        return 'Status: not ready' + (status.asr_error ? ' — ' + status.asr_error : '');
    }
    return status?.asr_ready ? 'Status: ready' : 'Status: not initialized (enable ASR below)';
}

// Render the capability contract. The UI shows the same reasons the resolver
// uses, so it never offers an engine the backend would reject without saying
// why.
function _renderCapabilities(status, asrText) {
    const records = Array.isArray(status?.capabilities) ? status.capabilities : [];
    if (!records.length) return;

    const warning = document.getElementById('asr-engine-warning');
    const transcription = records.filter(r => r.task === 'file_transcription');
    const usable = transcription.filter(r => !r.reason);

    if (warning) {
        if (usable.length) {
            warning.textContent = usable.map(r => r.model).join(', ') + ' available.';
        } else {
            const first = transcription[0];
            warning.textContent = 'No file-transcription engine is available. ' +
                (first && first.reason ? first.reason : 'Install whisper-cli and a Whisper model.');
        }
    }

    if (asrText && !status.asr_ready && usable.length === 0 &&
        (status.asr_readiness === 'not-ready' || status.asr_readiness === 'unconfigured')) {
        asrText.textContent = 'Status: unavailable — ' +
            (transcription[0] && transcription[0].reason ? transcription[0].reason : 'no engine installed');
    }
}

// ---------------------------------------------------------------------------
// Save / Reset
// ---------------------------------------------------------------------------

function _audioSnapshot() {
    const audio = window.settingsState?.current?.audio || {};
    return {
        asrEnabled: !!audio.asrEnabled,
        ttsEnabled: !!audio.ttsEnabled,
        asrEngine: audio.asrEngine || 'auto',
        ttsEngine: audio.ttsEngine || 'auto'
    };
}

async function _applySnapshot(audio) {
    const app = window.go?.main?.App;
    if (typeof app?.ApplyAudioSettings !== 'function') {
        throw new Error('audio backend is unavailable');
    }
    await app.ApplyAudioSettings(
        audio.asrEnabled,
        audio.ttsEnabled,
        audio.asrEngine,
        audio.ttsEngine
    );
}

// Applied is never inferred from a resolved promise. It requires the backend to
// publish the ACTIVE generation's effective configuration and that it equals
// what this panel requested.
function _confirmApplied(audio, status) {
    const applied = status && status.applied;
    if (!applied) return false;
    return !!applied.asr_enabled === !!audio.asrEnabled &&
           !!applied.tts_enabled === !!audio.ttsEnabled &&
           String(applied.asr_engine || '') === String(audio.asrEngine || '') &&
           String(applied.tts_engine || '') === String(audio.ttsEngine || '');
}

async function _readReadiness(audio) {
    const app = window.go?.main?.App;
    let status = null;
    if (typeof app?.GetAudioStatus === 'function') {
        try {
            status = await app.GetAudioStatus();
        } catch (err) {
            console.warn('[audio-settings] could not read status:', err);
        }
    }

    _applyInProgress = (status && typeof status.apply_in_progress === 'boolean')
        ? status.apply_in_progress
        : null;
    const appliedConfirmed = _confirmApplied(audio, status);
    const context = { appliedConfirmed, applyInProgress: _applyInProgress, status };

    if (!audio.asrEnabled) {
        return { ...context, state: 'disabled', label: 'Disabled (ASR off)', detail: 'Speech recognition is disabled.', closeable: true };
    }
    if (audio.asrEngine === 'auto') {
        return {
            ...context,
            state: 'unconfigured',
            label: 'Not ready (engine unconfigured)',
            detail: 'Auto-select does not identify a specific engine. Choose an explicit ASR engine before enabling speech recognition.',
            closeable: false
        };
    }

    // Apply success is not remote availability. Check this before trusting a
    // generic status flag from an older backend.
    if (REMOTE_ASR_ENGINES.has(audio.asrEngine) || status?.asr_readiness === 'remote-unverified') {
        return {
            ...context,
            state: 'remote-unverified',
            label: 'Remote engine — availability not verified',
            detail: status?.asr_error || 'The remote ASR configuration was applied, but this build has not verified that the service is reachable or can transcribe.',
            closeable: false
        };
    }
    if (!status) {
        return {
            ...context,
            state: 'unchecked',
            label: 'Readiness unchecked',
            detail: 'The configuration was applied, but the backend status could not be read.',
            closeable: false
        };
    }
    if (status.asr_readiness === 'local-initialized' && status.asr_ready) {
        return {
            ...context,
            state: 'local-initialized',
            label: 'Ready (local initialized)',
            detail: 'The local engine initialized. This does not prove a transcription succeeded.',
            closeable: true
        };
    }
    return {
        ...context,
        state: status.asr_readiness || 'not-ready',
        label: 'Not ready',
        detail: status.asr_error || 'The selected ASR engine did not initialize.',
        closeable: false
    };
}

function _dispatchParent(audio, readiness) {
    try {
        window.parent.dispatchEvent(new window.parent.CustomEvent('howl:audio-settings-applied', {
            detail: {
                audio: { ...audio },
                readiness: readiness.state,
                readinessLabel: readiness.label,
                detail: readiness.detail || ''
            }
        }));
    } catch (e) { /* cross-origin or no parent */ }
}

function _showApplied(readiness) {
    const confirmed = readiness.appliedConfirmed === true && readiness.applyInProgress === false;
    _lastOutcome = { saved: true, applied: confirmed };

    if (confirmed) {
        const tone = readiness.closeable ? '' : 'warn';
        _setOutcome('Saved', 'Applied', readiness.label, readiness.detail || '', tone);
    } else {
        _setOutcome('Saved', 'Indeterminate — effective configuration unconfirmed', readiness.label,
            'Activation reported success, but the effective configuration could not be confirmed.' +
            (readiness.detail ? ' ' + readiness.detail : ''), 'warn');
    }
    _setRetryVisible(!confirmed || readiness.state === 'not-ready' || readiness.state === 'unchecked' || readiness.state === 'apply-failed');

    if (readiness.state === 'local-initialized') {
        _notify('Audio settings saved and applied. ASR is initialized.', 'success');
    } else if (readiness.state === 'disabled') {
        _notify('Audio settings saved and applied.', 'success');
    } else if (readiness.state === 'remote-unverified') {
        _notify('Audio settings saved. Remote ASR availability is not verified.', 'warning');
    } else {
        _notify('Audio settings saved, but speech recognition is not ready.', 'warning');
    }
}

async function saveAudioSettings() {
    if (_saving) return false;
    await _ensureHydrated();
    if (!_hydrated) return false;

    // Refuse before any settings write; a programmatic caller must not bypass
    // the disabled Save state.
    if (_restoredLockHeld()) {
        _showRestoredLockState();
        return false;
    }

    const audio = _audioSnapshot();
    _saving = true;
    _setSaveEnabled(false);
    _setRetryVisible(false);

    try {
        _setOutcome('Saving…', 'Not applied', 'Readiness unchecked', 'Waiting for the settings file write…', 'busy');
        await window.settingsState.flushSettings();
    } catch (err) {
        console.error('[audio-settings] save failed:', err);
        _lastOutcome = { saved: false, applied: false };
        _setOutcome('Not saved', 'Not applied', 'Not ready',
            'The preference did not reach the settings file. Activation was not attempted. ' + err, 'error');
        _saving = false;
        _setSaveEnabled(true);
        return false;
    }

    await _refreshApplyInProgress();
    if (_activationBusy()) {
        _lastOutcome = { saved: true, applied: false };
        _setOutcome('Saved', 'Indeterminate — activation may still be running', 'Readiness unchecked',
            'The preference is saved, but a previous activation has not settled, so no second activation was started.', 'warn');
        _setRetryVisible(true);
        _saving = false;
        _setSaveEnabled(true);
        return false;
    }

    _setOutcome('Saved', 'Applying…', 'Readiness unchecked', 'Applying the full desired audio configuration…', 'busy');
    try {
        await _applySnapshot(audio);
    } catch (err) {
        if (err && err.name === 'ActivationLockedError') {
            _lastOutcome = { saved: true, applied: false };
            _setOutcome('Saved', 'Indeterminate — activation may still be running', 'Readiness unchecked',
                'Another activation is still in progress, so this request was not sent to the backend.', 'warn');
            _setRetryVisible(true);
            _saving = false;
            _setSaveEnabled(true);
            return false;
        }
        console.error('[audio-settings] apply failed:', err);
        _lastOutcome = { saved: true, applied: false };
        _setOutcome('Saved', 'Apply failed', 'Not ready',
            'The preference is saved. Activation failed: ' + err, 'error');
        _setRetryVisible(true);
        _dispatchParent(audio, { state: 'apply-failed', label: 'Apply failed', detail: String(err) });
        _notify('Audio settings saved, but activation failed. Use Retry activation.', 'error');
        _saving = false;
        _setSaveEnabled(true);
        return false;
    }

    const readiness = await _readReadiness(audio);
    _showApplied(readiness);
    _dispatchParent(audio, readiness);
    _saving = false;
    _setSaveEnabled(true);

    if (readiness.closeable) _finishClose();
    return true;
}

async function retryAudioActivation() {
    if (_saving) return false;
    await _ensureHydrated();
    if (!_hydrated) return false;

    await _refreshApplyInProgress();
    if (_restoredLockHeld()) {
        _showRestoredLockState();
        return false;
    }
    if (_activationBusy()) {
        _setOutcome('Saved', 'Indeterminate — activation may still be running', 'Readiness unchecked',
            'The previous activation has not settled. Retry stays disabled until it does.', 'warn');
        _setRetryVisible(true);
        return false;
    }

    const audio = _audioSnapshot();
    _saving = true;
    _setSaveEnabled(false);
    _setRetryVisible(false);
    _setOutcome('Saved', 'Applying…', 'Readiness unchecked', 'Retrying activation with the saved selection…', 'busy');
    try {
        await _applySnapshot(audio);
    } catch (err) {
        if (err && err.name === 'ActivationLockedError') {
            _lastOutcome = { saved: true, applied: false };
            _setOutcome('Saved', 'Indeterminate — activation may still be running', 'Readiness unchecked',
                'Another activation is still in progress, so this request was not sent to the backend.', 'warn');
            _setRetryVisible(true);
            _saving = false;
            _setSaveEnabled(true);
            return false;
        }
        console.error('[audio-settings] retry apply failed:', err);
        _lastOutcome = { saved: true, applied: false };
        _setOutcome('Saved', 'Apply failed', 'Not ready', 'Activation still failed: ' + err, 'error');
        _setRetryVisible(true);
        _dispatchParent(audio, { state: 'apply-failed', label: 'Apply failed', detail: String(err) });
        _saving = false;
        _setSaveEnabled(true);
        return false;
    }

    const readiness = await _readReadiness(audio);
    _showApplied(readiness);
    _dispatchParent(audio, readiness);
    _saving = false;
    _setSaveEnabled(true);
    if (readiness.closeable) _finishClose();
    return true;
}

async function resetAudioSettings() {
    if (_saving) return;
    if (!confirm('Reset all audio settings to their default values?')) return;
    await _ensureHydrated();
    if (!_hydrated) return;
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
    try {
        await window.settingsState.flushSettings();
        _lastOutcome = { saved: true, applied: false };
        _setOutcome('Saved', 'Not applied', 'Readiness unchecked', 'Defaults are saved. Save & Close applies them.', '');
        _notify('Audio settings reset and saved', 'info');
    } catch (err) {
        _lastOutcome = { saved: false, applied: false };
        _setOutcome('Not saved', 'Not applied', 'Not ready', 'Reset values could not be saved: ' + err, 'error');
    }
}

// Global aliases for HTML onclick handlers
function closeAudioSettings() { close(); }
function showAudioTab(tabId) { AudioSettings.showTab(tabId); }

// ---------------------------------------------------------------------------
// TTS Voice Model Registry
// ---------------------------------------------------------------------------

async function populateTTSDropdown() {
    const app = window.go?.main?.App;
    const sel = document.getElementById('tts-voice');
    if (!sel || !app?.GetTTSModels) return;
    try {
        const voices = await app.GetTTSModels();
        const prev = sel.value || window.settingsState?.current?.audio?.ttsVoice || '';
        // Keep the built-in option, then append registered ones
        const builtIn = sel.querySelector('option[value="en_US-lessac-medium"]');
        sel.innerHTML = '';
        if (builtIn) sel.appendChild(builtIn);
        else {
            const opt = document.createElement('option');
            opt.value = 'en_US-lessac-medium';
            opt.textContent = 'en_US-lessac-medium (built-in)';
            sel.appendChild(opt);
        }
        (voices || []).forEach(v => {
            if (v.name === 'en_US-lessac-medium') return; // already present
            const opt = document.createElement('option');
            opt.value = v.name;
            const sizeTxt = v.size_mb > 0 ? ` (${v.size_mb.toFixed(0)} MB)` : '';
            opt.textContent = `${v.name}${sizeTxt}`;
            sel.appendChild(opt);
        });
        if (prev && [...sel.options].some(o => o.value === prev)) sel.value = prev;
    } catch (err) {
        console.warn('[TTS] Failed to populate voice dropdown:', err);
    }
}

async function scanTTSDirectory() {
    const app = window.go?.main?.App;
    if (!app?.OpenDirectoryPicker || !app?.ScanTTSDir) return;
    try {
        const dir = await app.OpenDirectoryPicker();
        if (!dir) return;
        await app.ScanTTSDir(dir);
        await populateTTSDropdown();
    } catch (err) {
        console.error('[TTS] Directory scan failed:', err);
    }
}
