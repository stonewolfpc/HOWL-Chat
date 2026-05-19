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
    enableASR,
    enableTTS,
    isTTSEnabled,
    isASREnabled
};

// ---------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------

let _mediaRecorder = null;
let _recordedChunks = [];
let _recording = false;
let _ttsAudio = null;
let _asrEnabled = window.settingsState?.current?.audio?.asrEnabled || false;
let _ttsEnabled = window.settingsState?.current?.audio?.ttsEnabled || false;

// Sync with persisted settings whenever they change
window.addEventListener('DOMContentLoaded', () => {
    const s = window.settingsState?.current?.audio;
    if (s) {
        _asrEnabled = s.asrEnabled;
        _ttsEnabled = s.ttsEnabled;
    }
});

// ---------------------------------------------------------------------------
// Microphone Recording (ASR Input)
// ---------------------------------------------------------------------------

function isRecording() {
    return _recording;
}

async function startRecording() {
    if (_recording) return;
    if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
        showNotification('Microphone not supported in this browser', 'error');
        return;
    }

    try {
        const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
        _mediaRecorder = new MediaRecorder(stream);
        _recordedChunks = [];

        _mediaRecorder.ondataavailable = (e) => {
            if (e.data.size > 0) _recordedChunks.push(e.data);
        };

        _mediaRecorder.onstop = async () => {
            _recording = false;
            updateMicButtonState(false);
            const blob = new Blob(_recordedChunks, { type: 'audio/webm' });
            await processRecordedAudio(blob);
            // Stop all tracks to release mic
            stream.getTracks().forEach(t => t.stop());
        };

        _mediaRecorder.start();
        _recording = true;
        updateMicButtonState(true);
        showNotification('Recording started...', 'info');
    } catch (err) {
        console.error('Failed to start recording:', err);
        showNotification('Microphone access denied or unavailable', 'error');
        _recording = false;
        updateMicButtonState(false);
    }
}

function stopRecording() {
    if (!_recording || !_mediaRecorder) return;
    _mediaRecorder.stop();
    showNotification('Processing speech...', 'info');
}

function toggleRecording() {
    if (_recording) stopRecording();
    else startRecording();
}

async function processRecordedAudio(blob) {
    try {
        // Convert blob to file and save temporarily
        const arrayBuffer = await blob.arrayBuffer();
        const bytes = new Uint8Array(arrayBuffer);

        // Try to send to backend via Wails
        if (!window.go?.main?.App?.TranscribeAudio) {
            showNotification('ASR not available - backend binding missing', 'warning');
            return;
        }

        // Write blob to a temp file path that Go can read
        // We use a Data URI approach - let the backend handle it
        const tempPath = await saveBlobToTempFile(bytes, 'recording.webm');
        if (!tempPath) {
            showNotification('Failed to save audio for transcription', 'error');
            return;
        }

        const transcript = await window.go.main.App.TranscribeAudio(tempPath);
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

// Saves blob bytes to a temporary file path accessible by the backend.
// On Wails, we can write to the app's working directory.
async function saveBlobToTempFile(bytes, filename) {
    // Wails runtime doesn't give us direct file access from JS.
    // We return a data URL that the backend can decode, or rely on
    // the frontend sending the audio data directly.
    // For now, we create an Object URL and pass it.
    try {
        const blob = new Blob([bytes]);
        const url = URL.createObjectURL(blob);
        return url; // Backend will need to handle this
    } catch (e) {
        console.error('Failed to create audio URL:', e);
        return null;
    }
}

function insertTranscriptIntoChat(text) {
    const input = document.getElementById('chat-input');
    if (input) {
        input.value = text;
        input.focus();
        // Trigger input event for any listeners
        input.dispatchEvent(new Event('input', { bubbles: true }));
    }
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
        return await window.go.main.App.GetAudioStatus();
    } catch (err) {
        console.warn('GetAudioStatus failed:', err);
        return { asr_ready: false, tts_ready: false, active: false };
    }
}

async function enableASR(enabled) {
    _asrEnabled = enabled;
    if (!window.go?.main?.App?.EnableAudioASR) return;
    try {
        await window.go.main.App.EnableAudioASR(enabled);
        showNotification(enabled ? 'Speech recognition enabled' : 'Speech recognition disabled', 'info');
    } catch (err) {
        console.error('EnableAudioASR failed:', err);
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
