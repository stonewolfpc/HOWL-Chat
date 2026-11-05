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
    stopTTS,
    getStatus,
    syncFromSettings,
    getASRReadiness,
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
    if (typeof app?.GetAudioStatus !== 'function') return;
    try {
        const status = await app.GetAudioStatus();
        _asrReadiness = _readinessFromStatus(status);
        _asrReadinessDetail = status.asr_error || (_asrReadiness === 'remote-unverified' ?
            'remote engine availability is not locally verified' : '');
    } catch (err) {
        console.warn('GetAudioStatus failed while refreshing ASR readiness:', err);
    }
}

function _readinessFromStatus(status) {
    if (status?.asr_readiness === 'disabled') return 'disabled';
    const engine = window.settingsState?.current?.audio?.asrEngine;
    if (REMOTE_ASR_ENGINES.has(engine)) return 'remote-unverified';
    if (status?.asr_readiness) return status.asr_readiness;
    return status?.asr_ready ? 'local-initialized' : 'not-ready';
}

// Re-apply the persisted engine choices once settings have loaded. Without this
// the dropdowns could restore an engine the backend had never been told about.
async function restorePersistedAudio() {
    if (_restoreApplied) return;
    const s = window.settingsState?.current?.audio;
    const app = window.go?.main?.App;
    if (!s || (!s.asrEnabled && !s.ttsEnabled)) return;
    if (typeof app?.ApplyAudioSettings !== 'function') return;
    _restoreApplied = true;
    try {
        await app.ApplyAudioSettings(
            !!s.asrEnabled,
            !!s.ttsEnabled,
            s.asrEngine || 'auto',
            s.ttsEngine || 'auto'
        );
        await _refreshASRReadiness();
    } catch (err) {
        console.warn('Audio settings restore failed:', err);
        _asrReadiness = 'apply-failed';
        _asrReadinessDetail = String(err);
    }
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
        showNotification('Microphone access denied or unavailable', 'error');
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

    if (!voice?.enabled) { _voiceConfigCache.set(characterId, null); return null; }

    const merged = {
        engine:           voice.engine  || global.ttsEngine   || 'piper',
        voice:            voice.voice   || global.ttsVoice    || 'en_US-lessac-medium',
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

async function playTTS(text, characterId) {
    if (!text || !_ttsEnabled) return;

    stopTTS(); // Stop any current playback

    const voiceCfg = _getMergedVoiceConfig(characterId);
    const engine = voiceCfg?.engine || window.settingsState?.current?.audio?.ttsEngine || 'piper';

    // Browser TTS path (works with per-character voices natively)
    if (engine === 'browser') {
        useBrowserTTS(text, voiceCfg);
        return;
    }

    // Backend TTS path with optional voice override
    try {
        if (window.go?.main?.App?.SynthesizeSpeech) {
            let audioPath;
            if (voiceCfg && window.go.main.App.SynthesizeSpeechWithVoice) {
                audioPath = await window.go.main.App.SynthesizeSpeechWithVoice(text, JSON.stringify(voiceCfg));
            } else {
                audioPath = await window.go.main.App.SynthesizeSpeech(text);
            }
            if (audioPath) {
                _ttsAudio = new Audio(audioPath);
                _ttsAudio.volume = voiceCfg?.volume ?? (window.settingsState?.current?.audio?.ttsVolume || 1.0);
                await _ttsAudio.play();
            }
        } else {
            useBrowserTTS(text, voiceCfg);
        }
    } catch (err) {
        console.error('TTS error:', err);
        useBrowserTTS(text, voiceCfg);
    }
}

function stopTTS() {
    if (_ttsAudio) {
        _ttsAudio.pause();
        _ttsAudio.currentTime = 0;
        _ttsAudio = null;
    }
    // Also stop browser TTS
    if (window.speechSynthesis) {
        window.speechSynthesis.cancel();
    }
}

function useBrowserTTS(text, voiceCfg) {
    if (!window.speechSynthesis) return;
    const utterance = new SpeechSynthesisUtterance(text);
    const voices = window.speechSynthesis.getVoices();
    // Try to find a voice matching the requested model name or fall back to a good English voice
    let preferred = null;
    if (voiceCfg?.voice) {
        const search = voiceCfg.voice.toLowerCase();
        preferred = voices.find(v => v.name.toLowerCase().includes(search) || v.lang.toLowerCase().includes(search));
    }
    if (!preferred) {
        preferred = voices.find(v => v.name.includes('Google') || v.name.includes('Natural') || v.name.includes('Samantha'));
    }
    if (preferred) utterance.voice = preferred;
    utterance.rate  = voiceCfg?.speed  ?? (window.settingsState?.current?.audio?.ttsSpeed  || 1.0);
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
