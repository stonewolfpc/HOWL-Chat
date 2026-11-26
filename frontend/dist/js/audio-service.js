/**
 * Audio Service Module
 *
 * Wraps Wails audio bindings for ASR (speech-to-text) and TTS (text-to-speech).
 * Provides microphone recording, transcription, and AI voice playback.
 *
 * @module audio-service
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

window.AudioService = {
    isRecording,
    startRecording,
    stopRecording,
    toggleRecording,
    playTTS,
    speakMessage,
    speakComposer,
    stopTTS,
    invalidateVoiceCache,
    getStatus,
    syncFromSettings,
    getASRReadiness,
    getMicrophoneStatus,
    enableASR,
    enableTTS,
    isTTSEnabled,
    isASREnabled
};

// ---------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------

let _audioContext = null;
let _mediaStream = null;
let _sourceNode = null;
let _processorNode = null;
let _recordedChunks = [];
let _recordingSampleRate = 16000;
let _recording = false;
let _starting = false;
let _processing = false;
let _maxRecordingTimer = null;
let _ttsAudio = null;
let _asrEnabled = window.settingsState?.current?.audio?.asrEnabled || false;
let _ttsEnabled = window.settingsState?.current?.audio?.ttsEnabled || false;
let _asrReadiness = _asrEnabled ? 'not-ready' : 'disabled';
let _asrReadinessDetail = '';
let _micStatus = null;

const MIC_NO_DEVICE_MESSAGE = "Windows shows a mic isn't selected. Please connect a microphone to use this option";
const MIC_UNAVAILABLE_MESSAGE = 'The selected microphone is unavailable. Reconnect it or choose another microphone in Windows.';
const MIC_ACCESS_DENIED_MESSAGE = 'Windows is blocking microphone access. Allow desktop apps to use the microphone in Windows Settings.';
const MIC_INIT_FAILED_MESSAGE = 'The microphone is present but could not be opened for recording.';

const MAX_RECORDING_SECONDS = 120;
const REMOTE_ASR_ENGINES = new Set(['qwen3-asr-0.6b', 'qwen3-asr-1.7b']);

let _restoreApplied = false;

function syncFromSettings() {
    const s = window.settingsState?.current?.audio;
    if (!s) return;
    _asrEnabled = !!s.asrEnabled;
    _ttsEnabled = !!s.ttsEnabled;
    if (!_asrEnabled) {
        _asrReadiness = 'disabled';
        _asrReadinessDetail = 'ASR is disabled.';
    } else if (_asrReadiness === 'disabled') {
        _asrReadiness = 'not-ready';
        _asrReadinessDetail = '';
    }
}

async function _refreshASRReadiness() {
    const app = window.go?.main?.App;
    if (typeof app?.GetAudioStatus !== 'function') return null;
    try {
        const status = await app.GetAudioStatus();
        _micStatus = status?.mic || null;
        _asrReadiness = _readinessFromStatus(status);
        _asrReadinessDetail = status.asr_error || (_asrReadiness === 'remote-unverified' ?
            'remote engine availability is not locally verified' : '');
        return status;
    } catch (err) {
        console.warn('GetAudioStatus failed while refreshing ASR readiness:', err);
        _micStatus = null;
        return null;
    }
}

// Read capture availability without touching readiness: readiness comes from the
// applied generation, while the microphone is re-probed on every explicit
// recording attempt so a device connected after startup is picked up.
async function _refreshMicrophoneStatus() {
    const app = window.go?.main?.App;
    if (typeof app?.GetAudioStatus !== 'function') {
        _micStatus = null;
        return null;
    }
    try {
        const status = await app.GetAudioStatus();
        _micStatus = status?.mic || null;
        return _micStatus;
    } catch (err) {
        console.warn('GetAudioStatus failed while checking the microphone:', err);
        _micStatus = null;
        return null;
    }
}

function _microphoneErrorMessage(err) {
    const name = err && err.name;
    if (name === 'NotFoundError' || name === 'DevicesNotFoundError') return MIC_NO_DEVICE_MESSAGE;
    if (name === 'NotAllowedError' || name === 'SecurityError' || name === 'PermissionDeniedError') return MIC_ACCESS_DENIED_MESSAGE;
    if (name === 'NotReadableError' || name === 'TrackStartError') return MIC_UNAVAILABLE_MESSAGE;
    if (_micStatus && _micStatus.message) return _micStatus.message;
    return MIC_INIT_FAILED_MESSAGE;
}

function _readinessFromStatus(status) {
    if (status?.asr_readiness === 'disabled') return 'disabled';
    const engine = window.settingsState?.current?.audio?.asrEngine;
    if (REMOTE_ASR_ENGINES.has(engine)) return 'remote-unverified';
    if (status?.asr_readiness) return status.asr_readiness;
    return status?.asr_ready ? 'local-initialized' : 'not-ready';
}

// Restore persisted selections into the UI and mirror the backend's current
// status. This deliberately does not call ApplyAudioSettings: a saved selection
// must not load a model or reserve VRAM at startup. Activation stays explicit
// through Save or Retry.
async function restorePersistedAudio() {
    if (_restoreApplied) return;
    const s = window.settingsState?.current?.audio;
    if (!s) return;
    _restoreApplied = true;
    syncFromSettings();
    await _refreshASRReadiness();
}

// Sync with persisted settings whenever they change
window.addEventListener('DOMContentLoaded', () => {
    syncFromSettings();
    restorePersistedAudio();
});
window.addEventListener('howl:settings-loaded', () => {
    syncFromSettings();
    restorePersistedAudio();
});

// The audio-settings iframe saves and applies in one user action. Receive the
// result immediately so the parent microphone button does not need a reload to
// see the new enabled state or readiness scope.
window.addEventListener('howl:audio-settings-applied', (event) => {
    const detail = event?.detail;
    if (!detail?.audio) return;
    if (window.settingsState?.current?.audio) {
        Object.assign(window.settingsState.current.audio, detail.audio);
    }
    // Global engine/voice/speed defaults are part of every inherited
    // character config, so the per-character merge cache cannot survive them.
    _invalidateVoiceCache();
    syncFromSettings();
    _restoreApplied = true;
    _asrReadiness = detail.readiness || 'not-ready';
    _asrReadinessDetail = detail.detail || '';
    window.dispatchEvent(new CustomEvent('howl:audio-readiness-changed', {
        detail: { state: _asrReadiness, detail: _asrReadinessDetail }
    }));
});

// ---------------------------------------------------------------------------
// Microphone Recording (ASR Input)
// ---------------------------------------------------------------------------

function isRecording() {
    return _recording;
}

async function startRecording() {
    if (_recording || _starting || _processing) return;
    if (!_asrEnabled) {
        showNotification('Speech recognition is disabled', 'warning');
        return;
    }
    // Capture availability is checked before the readiness gate so a missing or
    // unavailable microphone reports its own state instead of a generic
    // "not ready", matching the backend that refuses ASR allocation.
    await _refreshMicrophoneStatus();
    if (_micStatus && _micStatus.state && _micStatus.state !== 'ready') {
        showNotification(_micStatus.message || MIC_UNAVAILABLE_MESSAGE, 'warning');
        return;
    }
    const blocked = _asrBlockedReason();
    if (blocked) {
        showNotification(blocked, 'warning');
        return;
    }
    if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
        showNotification('Microphone not supported in this browser', 'error');
        return;
    }
    if (!window.AudioWav || typeof window.AudioWav.encodeWav16 !== 'function') {
        showNotification('WAV encoder is unavailable', 'error');
        return;
    }

    _starting = true;
    try {
        const stream = await navigator.mediaDevices.getUserMedia({
            audio: {
                channelCount: 1,
                echoCancellation: true,
                noiseSuppression: true,
                autoGainControl: true
            }
        });
        _mediaStream = stream;

        const AudioContextCtor = window.AudioContext || window.webkitAudioContext;
        if (!AudioContextCtor) throw new Error('Web Audio is not available');

        let context;
        try {
            context = new AudioContextCtor({ sampleRate: 16000 });
        } catch (_) {
            context = new AudioContextCtor();
        }

        _audioContext = context;
        _sourceNode = context.createMediaStreamSource(stream);
        _processorNode = context.createScriptProcessor(4096, 1, 1);
        _recordedChunks = [];
        _recordingSampleRate = context.sampleRate || 16000;

        _processorNode.onaudioprocess = (event) => {
            if (!_recording) return;
            const channel = event.inputBuffer.getChannelData(0);
            _recordedChunks.push(new Float32Array(channel));
        };

        _sourceNode.connect(_processorNode);
        _processorNode.connect(context.destination);
        await context.resume();

        _recording = true;
        updateMicButtonState(true);
        showNotification('Recording started...', 'info');
        _maxRecordingTimer = window.setTimeout(() => {
            if (_recording) {
                showNotification('Recording limit reached', 'info');
                stopRecording();
            }
        }, MAX_RECORDING_SECONDS * 1000);
    } catch (err) {
        console.error('Failed to start recording:', err);
        releaseCapture();
        _recording = false;
        updateMicButtonState(false);
        showNotification(_microphoneErrorMessage(err), 'error');
    } finally {
        _starting = false;
    }
}

function _asrBlockedReason() {
    if (_asrReadiness === 'local-initialized' || _asrReadiness === 'remote-unverified') return '';
    if (_asrReadiness === 'not-ready') {
        return 'Speech recognition is saved but not ready' +
            (_asrReadinessDetail ? ': ' + _asrReadinessDetail : '');
    }
    if (_asrReadiness === 'apply-failed') {
        return 'Speech recognition settings could not be applied' +
            (_asrReadinessDetail ? ': ' + _asrReadinessDetail : '');
    }
    if (_asrReadiness === 'unconfigured') {
        return 'Speech recognition needs an explicit ASR engine';
    }
    if (_asrReadiness === 'unchecked') {
        return 'Speech recognition readiness has not been checked';
    }
    return '';
}

function stopRecording() {
    if (!_recording) return;

    _recording = false;
    updateMicButtonState(false);
    if (_maxRecordingTimer) {
        clearTimeout(_maxRecordingTimer);
        _maxRecordingTimer = null;
    }

    const chunks = _recordedChunks;
    const sampleRate = _recordingSampleRate;
    _recordedChunks = [];
    releaseCapture();

    showNotification('Processing speech...', 'info');
    _processing = true;
    (async () => {
        try {
            const samples = mergeSamples(chunks);
            if (!hasSignal(samples)) {
                showNotification('No speech detected', 'warning');
                return;
            }
            const wav = window.AudioWav.encodeWav16(samples, sampleRate);
            await processRecordedAudio(wav);
        } catch (err) {
            console.error('Microphone processing error:', err);
            showNotification('Transcription failed: ' + err, 'error');
        } finally {
            _processing = false;
        }
    })();
}

function toggleRecording() {
    if (_recording) stopRecording();
    else startRecording();
}

function releaseCapture() {
    if (_processorNode) {
        _processorNode.onaudioprocess = null;
        try { _processorNode.disconnect(); } catch (_) {}
        _processorNode = null;
    }
    if (_sourceNode) {
        try { _sourceNode.disconnect(); } catch (_) {}
        _sourceNode = null;
    }
    if (_audioContext) {
        const context = _audioContext;
        _audioContext = null;
        if (context.state !== 'closed') context.close().catch(() => {});
    }
    if (_mediaStream) {
        _mediaStream.getTracks().forEach((track) => {
            try { track.stop(); } catch (_) {}
        });
        _mediaStream = null;
    }
}

function mergeSamples(chunks) {
    let length = 0;
    for (const chunk of chunks) length += chunk.length;
    const merged = new Float32Array(length);
    let offset = 0;
    for (const chunk of chunks) {
        merged.set(chunk, offset);
        offset += chunk.length;
    }
    return merged;
}

function hasSignal(samples) {
    if (!samples || samples.length === 0) return false;
    let sum = 0;
    for (let i = 0; i < samples.length; i++) sum += samples[i] * samples[i];
    return Math.sqrt(sum / samples.length) > 0.002;
}

function bytesToDataURL(bytes) {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = () => resolve(String(reader.result));
        reader.onerror = () => reject(reader.error || new Error('failed to read audio'));
        reader.readAsDataURL(new Blob([bytes], { type: 'audio/wav' }));
    });
}

async function processRecordedAudio(wavBytes) {
    try {
        const app = window.go?.main?.App;
        if (!app?.TranscribeAudioData) {
            showNotification('ASR not available - backend binding missing', 'warning');
            return;
        }

        const dataURL = await bytesToDataURL(wavBytes);
        const transcript = await app.TranscribeAudioData(dataURL);
        if (transcript && transcript.trim()) {
            insertTranscriptIntoChat(transcript.trim());
            showNotification('Transcription complete', 'success');
        } else {
            showNotification('No speech detected', 'warning');
        }
    } catch (err) {
        console.error('Transcription error:', err);
        showNotification('Transcription failed: ' + err, 'error');
    }
}

function insertTranscriptIntoChat(text) {
    const input = document.getElementById('message-input') || document.querySelector('.input-field');
    if (!input) return;
    // Preserve the existing composer semantics: transcription fills the draft;
    // the user still reviews and sends it explicitly.
    input.value = text;
    input.focus();
    input.dispatchEvent(new Event('input', { bubbles: true }));
}

function updateMicButtonState(recording) {
    const btn = document.querySelector('.mic-button');
    if (!btn) return;
    if (recording) {
        btn.classList.add('recording');
        btn.style.color = '#f87171';
        btn.textContent = '⏹️';
    } else {
        btn.classList.remove('recording');
        btn.style.color = '';
        btn.textContent = '🎙️';
    }
}

// ---------------------------------------------------------------------------
// TTS Playback
// ---------------------------------------------------------------------------

// Per-character merged voice config cache
const _voiceConfigCache = new Map();

function _getMergedVoiceConfig(characterId) {
    if (!characterId) return null;
    if (_voiceConfigCache.has(characterId)) return _voiceConfigCache.get(characterId);

    const global = window.settingsState?.current?.audio || {};
    const cd = (typeof currentChatData !== 'undefined') ? currentChatData
             : (window.sidebarState ? window.sidebarState.currentChatData : null);
    const character = cd?.characters?.find(c => c.id === characterId);
    const voice = character?.howl?.voice;

    // A character with no per-character voice remains a no-override request.
    // The backend applies the active global engine/voice defaults in that case.
    if (!voice?.enabled) { _voiceConfigCache.set(characterId, null); return null; }

    const engine = voice.engine || '';
    const merged = {
        character_id:     characterId,
        // Keep the character's inheritance decision intact. The backend
        // resolves a blank engine against the active global generation.
        engine:           engine,
        // A voice belongs to an engine. An explicit engine with no chosen
        // model inherits that engine's default, not another engine's global
        // voice (for example a Piper name sent to Kokoro).
        voice:            voice.voice || (engine ? '' : (global.ttsVoice || '')),
        speed:            voice.speed   ?? global.ttsSpeed    ?? 1.0,
        volume:           voice.volume  ?? global.ttsVolume   ?? 0.8,
        noiseScale:       voice.noiseScale       ?? global.ttsNoiseScale      ?? 0.667,
        noiseW:           voice.noiseW           ?? global.ttsNoiseW          ?? 0.8,
        sentenceSilence:  voice.sentenceSilence  ?? global.ttsSentenceSilence ?? 0.2
    };
    _voiceConfigCache.set(characterId, merged);
    return merged;
}

function _invalidateVoiceCache(characterId) {
    if (characterId) _voiceConfigCache.delete(characterId);
    else _voiceConfigCache.clear();
}

function invalidateVoiceCache(characterId) {
    _invalidateVoiceCache(characterId);
}

window.addEventListener('howl:character-saved', (event) => {
    _invalidateVoiceCache(event?.detail?.id || null);
});

function _selectedOutputDevice() {
    const audio = window.settingsState?.current?.audio || {};
    return String(audio.audioOutputDevice || 'default');
}

// Resolve an explicitly selected output endpoint. The default setting is
// intentionally allowed to use the Windows default output. An explicit device
// is pinned by ID and never silently replaced if it disappears.
async function _resolveExplicitOutputSink() {
    const selected = _selectedOutputDevice();
    if (selected === 'none') {
        throw new Error('Audio output is disabled in Audio Settings.');
    }
    if (selected === 'default') return null;
    if (!navigator.mediaDevices?.enumerateDevices) {
        throw new Error('This WebView cannot enumerate audio outputs, so the selected output cannot be resolved.');
    }
    const devices = await navigator.mediaDevices.enumerateDevices();
    const outputs = devices.filter(device => device.kind === 'audiooutput');
    if (outputs.length === 0) {
        throw new Error('No audio output endpoints are available. Refresh Audio Settings.');
    }
    const chosen = outputs.find(device => device.deviceId === selected);
    if (!chosen) {
        throw new Error('The selected output device is no longer available. Choose it again in Audio Settings.');
    }
    return chosen.deviceId;
}

async function _pinOutputDevice(audioElement) {
    const sinkId = await _resolveExplicitOutputSink();
    if (!sinkId) return;
    if (!audioElement || typeof audioElement.setSinkId !== 'function') {
        throw new Error('This WebView cannot select the requested output endpoint.');
    }
    await audioElement.setSinkId(sinkId);
}

let _activeSpeakButton = null;
let _ttsObjectUrl = '';

function _notify(message, type = 'info', duration = 5000) {
    if (typeof showNotification === 'function') {
        showNotification(message, type, duration);
    } else if (window.showNotification) {
        window.showNotification(message, type, duration);
    } else {
        console.warn('[tts]', message);
    }
}

// The saved setting is the authority. The module copy is refreshed when
// settings load, but a click that arrives first must not be judged by a stale
// value and then ignored without a word.
function _ttsEnabledNow() {
    const live = window.settingsState?.current?.audio?.ttsEnabled;
    return typeof live === 'boolean' ? live : !!_ttsEnabled;
}

function _markSpeaking(button, speaking) {
    if (_activeSpeakButton && _activeSpeakButton !== button) {
        _activeSpeakButton.classList.remove('playing');
    }
    if (!button) {
        _activeSpeakButton = null;
        return;
    }
    button.classList.toggle('playing', !!speaking);
    _activeSpeakButton = speaking ? button : null;
}

function _releaseTTSObjectUrl() {
    if (_ttsObjectUrl) {
        URL.revokeObjectURL(_ttsObjectUrl);
        _ttsObjectUrl = '';
    }
}

// The backend returns a cache path, and a page cannot play a Windows path as
// media: the browser resolves it against the page origin and the load fails.
// Read the bytes through the binding and play a local object URL instead.
async function _ttsSourceFor(audioPath) {
    const app = window.go?.main?.App;
    if (typeof app?.GetTTSAudioData !== 'function') return audioPath;
    const dataUrl = await app.GetTTSAudioData(audioPath);
    if (!dataUrl) throw new Error('the synthesized audio could not be read');
    const comma = dataUrl.indexOf(',');
    const meta = comma >= 0 ? dataUrl.slice(0, comma) : '';
    if (!meta.includes(';base64')) return dataUrl;
    const binary = atob(dataUrl.slice(comma + 1));
    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) bytes[i] = binary.charCodeAt(i);
    return URL.createObjectURL(new Blob([bytes], { type: 'audio/wav' }));
}

async function playTTS(text, characterId, sourceButton) {
    const speakable = typeof text === 'string' ? text.trim() : '';
    if (!speakable) {
        _notify('Nothing to speak yet — type a message, or use the speaker on an AI reply.', 'info');
        return;
    }
    if (!_ttsEnabledNow()) {
        _notify('Text-to-speech is off. Turn it on in Settings → Audio → Text-to-Speech.', 'warning', 7000);
        return;
    }

    stopTTS(); // Stop any current playback

    const voiceCfg = _getMergedVoiceConfig(characterId);
    const engine = voiceCfg?.engine || window.settingsState?.current?.audio?.ttsEngine || 'piper';

    // Browser TTS follows the default output. It cannot target a specific
    // endpoint, so an explicit non-default selection is reported as a genuine
    // technical limitation rather than silently using another output.
    if (engine === 'browser') {
        useBrowserTTS(speakable, voiceCfg);
        return;
    }

    // Backend TTS path with optional voice override
    if (!window.go?.main?.App?.SynthesizeSpeech) {
        useBrowserTTS(speakable, voiceCfg);
        return;
    }
    try {
        const app = window.go.main.App;
        const audioPath = voiceCfg && app.SynthesizeSpeechWithVoice
            ? await app.SynthesizeSpeechWithVoice(speakable, JSON.stringify(voiceCfg))
            : await app.SynthesizeSpeech(speakable);
        if (!audioPath) throw new Error('the speech engine returned no audio file');

        const src = await _ttsSourceFor(audioPath);
        if (src.startsWith('blob:')) _ttsObjectUrl = src;

        const audio = new Audio(src);
        await _pinOutputDevice(audio);
        audio.volume = voiceCfg?.volume ?? (window.settingsState?.current?.audio?.ttsVolume || 1.0);
        _ttsAudio = audio;
        _markSpeaking(sourceButton, true);

        const finish = () => {
            if (_ttsAudio === audio) _ttsAudio = null;
            _releaseTTSObjectUrl();
            _markSpeaking(sourceButton, false);
        };
        audio.onended = finish;
        audio.onerror = () => {
            finish();
            _notify('The speech audio could not be played.', 'error');
        };
        await audio.play();
    } catch (err) {
        _releaseTTSObjectUrl();
        _markSpeaking(sourceButton, false);
        console.error('TTS error:', err);
        _notify('TTS unavailable: ' + (err?.message || err), 'error', 8000);
    }
}

// The composer speaker is a standalone action: it reads the draft, speaks it,
// and never sends, saves, or creates a message. It needs no chat model — only
// the TTS engine — so it keeps working when no model is loaded or the chat
// backend is down. The draft is left untouched in the input.
function speakComposer() {
    const input = document.getElementById('message-input') || document.querySelector('.input-field');
    const button = document.getElementById('composer-speak-btn') || document.querySelector('.tts-icon');
    if (button && button.classList.contains('playing')) { stopTTS(); return; }

    const draft = input ? String(input.value || '').trim() : '';
    if (!draft) {
        _notify('Nothing to read yet — type a message, then use the speaker. It reads your draft without sending it.', 'info', 6000);
        if (input) input.focus();
        return;
    }
    return playTTS(draft, null, button);
}

// Speak one rendered message. The bubble owns the text and the message owns the
// character, so a per-character voice applies exactly as it does on autoplay.
function speakMessage(button) {
    const bubble = button?.closest?.('.chat-bubble');
    if (!bubble) return;
    if (button.classList.contains('playing')) { stopTTS(); return; }
    const text = (bubble.querySelector('.chat-text')?.textContent || '').trim();
    if (!text) { _notify('This message has no text to speak.', 'info'); return; }
    const characterId = button.closest('.chat-message')?.dataset?.characterId || null;
    return playTTS(text, characterId, button);
}

function stopTTS() {
    if (_ttsAudio) {
        _ttsAudio.onended = null;
        _ttsAudio.onerror = null;
        _ttsAudio.pause();
        _ttsAudio.currentTime = 0;
        _ttsAudio = null;
    }
    _releaseTTSObjectUrl();
    _markSpeaking(null, false);
    // Also stop browser TTS
    if (window.speechSynthesis) {
        window.speechSynthesis.cancel();
    }
}

function useBrowserTTS(text, voiceCfg) {
    const selected = _selectedOutputDevice();
    if (selected === 'none') {
        if (window.showNotification) window.showNotification('Audio output is disabled in Audio Settings.', 'error');
        return;
    }
    if (selected !== 'default') {
        if (window.showNotification) {
            window.showNotification('Browser TTS cannot target a specific output device. Select the Windows default output or use a backend TTS engine.', 'error');
        }
        return;
    }
    if (!window.speechSynthesis) {
        if (window.showNotification) window.showNotification('Browser TTS is unavailable in this WebView.', 'error');
        return;
    }

    const utterance = new SpeechSynthesisUtterance(text);
    const voices = window.speechSynthesis.getVoices();
    let preferred = null;
    if (voiceCfg?.voice) {
        const search = voiceCfg.voice.toLowerCase();
        preferred = voices.find(v => v.name.toLowerCase().includes(search) || v.lang.toLowerCase().includes(search));
    }
    if (!preferred) {
        preferred = voices.find(v => v.name.includes('Google') || v.name.includes('Natural') || v.name.includes('Samantha'));
    }
    if (preferred) utterance.voice = preferred;
    utterance.rate = voiceCfg?.speed ?? (window.settingsState?.current?.audio?.ttsSpeed || 1.0);
    utterance.volume = voiceCfg?.volume ?? (window.settingsState?.current?.audio?.ttsVolume || 1.0);
    window.speechSynthesis.speak(utterance);
}

// ---------------------------------------------------------------------------
// Status & Settings
// ---------------------------------------------------------------------------

async function getStatus() {
    if (!window.go?.main?.App?.GetAudioStatus) {
        return { asr_ready: false, tts_ready: false, active: false };
    }
    try {
        const status = await window.go.main.App.GetAudioStatus();
        _asrReadiness = _readinessFromStatus(status);
        _asrReadinessDetail = status.asr_error || '';
        return status;
    } catch (err) {
        console.warn('GetAudioStatus failed:', err);
        return { asr_ready: false, tts_ready: false, active: false };
    }
}

async function enableASR(enabled) {
    const app = window.go?.main?.App;
    if (!app?.EnableAudioASR) {
        _asrEnabled = enabled;
        _asrReadiness = enabled ? 'not-ready' : 'disabled';
        return;
    }
    try {
        await app.EnableAudioASR(enabled);
        _asrEnabled = enabled;
        await _refreshASRReadiness();
        showNotification(enabled ? 'Speech recognition enabled' : 'Speech recognition disabled', 'info');
    } catch (err) {
        console.error('EnableAudioASR failed:', err);
        _asrReadiness = 'apply-failed';
        _asrReadinessDetail = String(err);
        showNotification('Failed to toggle ASR: ' + err, 'error');
    }
}

async function enableTTS(enabled) {
    _ttsEnabled = enabled;
    if (!window.go?.main?.App?.EnableAudioTTS) return;
    try {
        await window.go.main.App.EnableAudioTTS(enabled);
        showNotification(enabled ? 'Text-to-speech enabled' : 'Text-to-speech disabled', 'info');
    } catch (err) {
        console.error('EnableAudioTTS failed:', err);
        showNotification('Failed to toggle TTS: ' + err, 'error');
    }
}

function isTTSEnabled() {
    return _ttsEnabled;
}

function isASREnabled() {
    return _asrEnabled;
}

function getASRReadiness() {
    return { state: _asrReadiness, detail: _asrReadinessDetail };
}

function getMicrophoneStatus() {
    return _micStatus;
}

// ---------------------------------------------------------------------------
// Auto-play TTS for AI responses (hooked into chat stream)
// ---------------------------------------------------------------------------

function initTTSAutoplay() {
    // Hook into the chat completion event
    if (typeof window.EventsOn === 'function') {
        window.EventsOn('chat:complete', (data) => {
            const shouldAutoplay = window.settingsState?.current?.audio?.ttsAutoplayMessage;
            if (shouldAutoplay && data?.message) {
                // Pass characterId so per-character voice overrides are applied
                const text = extractSpeakableText(data.message);
                if (text) playTTS(text, data?.characterId || null);
            }
        });
    }
}

function extractSpeakableText(htmlOrText) {
    if (!htmlOrText) return '';
    // Strip HTML tags if present
    const tmp = document.createElement('div');
    tmp.innerHTML = htmlOrText;
    const text = tmp.textContent || tmp.innerText || '';
    // Limit length to avoid super long TTS
    const maxLen = 500;
    return text.length > maxLen ? text.substring(0, maxLen) + '...' : text;
}

// Initialize autoplay hook when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initTTSAutoplay);
} else {
    setTimeout(initTTSAutoplay, 1000);
}
