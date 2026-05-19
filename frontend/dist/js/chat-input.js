/**
 * Chat Input & Attachments
 *
 * Manages the send/abort button, Enter-key shortcut,
 * narrator event flushing before send, and media attachment
 * handlers (image, audio, file).
 *
 * Depends on: chat-render.js, chat-stream.js
 *
 * @module chat-input
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Image Attachment State
// ---------------------------------------------------------------------------

let currentImage = null;

const VISION_MAX_DIM = 1024;
const VISION_JPEG_QUALITY = 0.85;

function addImage() {
    const input    = document.createElement('input');
    input.type     = 'file';
    input.accept   = 'image/*';
    input.style.display = 'none';
    document.body.appendChild(input);

    input.onchange = async (event) => {
        const file = event.target.files[0];
        if (!file) { document.body.removeChild(input); return; }

        // Resize and encode as JPEG to keep payloads small and ensure
        // llama-server compatibility (stb_image does not support WebP).
        try {
            const bitmap = await createImageBitmap(file);
            const scale  = Math.min(1, VISION_MAX_DIM / Math.max(bitmap.width, bitmap.height));
            const w      = Math.round(bitmap.width  * scale);
            const h      = Math.round(bitmap.height * scale);

            const canvas = document.createElement('canvas');
            canvas.width  = w;
            canvas.height = h;
            canvas.getContext('2d').drawImage(bitmap, 0, 0, w, h);
            bitmap.close();

            currentImage = canvas.toDataURL('image/jpeg', VISION_JPEG_QUALITY);
            showImagePreview(currentImage);
        } catch (err) {
            console.error('Vision image resize failed:', err);
            // Fallback: read raw (legacy behavior)
            const reader = new FileReader();
            reader.onload = (e) => { currentImage = e.target.result; showImagePreview(currentImage); };
            reader.readAsDataURL(file);
        }
        document.body.removeChild(input);
    };
    input.click();
}

function showImagePreview(imageData) {
    const existingPreview = document.querySelector('.image-preview');
    if (existingPreview) existingPreview.remove();

    const preview   = document.createElement('div');
    preview.className = 'image-preview';
    preview.innerHTML = `
        <img src="${imageData}" alt="Preview">
        <button class="remove-image" onclick="removeImage()">×</button>
    `;

    const inputBar = document.querySelector('.input-bar');
    if (inputBar) inputBar.parentNode.insertBefore(preview, inputBar);
}

function removeImage() {
    currentImage = null;
    const preview = document.querySelector('.image-preview');
    if (preview) preview.remove();
}

// ---------------------------------------------------------------------------
// Audio Attachment (stub — TTS integration placeholder)
// ---------------------------------------------------------------------------

let ttsAudio = null;

function addAudio() {
    console.log('Add audio — not yet implemented');
}

function textToSpeech() {
    const fileInput   = document.createElement('input');
    fileInput.type    = 'file';
    fileInput.accept  = 'audio/*';
    fileInput.style.display = 'none';

    fileInput.onchange = function (e) {
        const file = e.target.files[0];
        if (!file) return;

        const audioUrl = URL.createObjectURL(file);
        if (ttsAudio) { ttsAudio.pause(); ttsAudio = null; }

        ttsAudio = new Audio(audioUrl);
        const ttsButton = document.querySelector('.tts-icon');
        if (ttsButton) ttsButton.classList.add('active');

        ttsAudio.onended = function () {
            URL.revokeObjectURL(audioUrl);
            if (ttsButton) ttsButton.classList.remove('active');
        };
        ttsAudio.onerror = function () {
            console.error('Error playing TTS audio');
            URL.revokeObjectURL(audioUrl);
            if (ttsButton) ttsButton.classList.remove('active');
        };
        ttsAudio.play().catch(err => {
            console.error('Failed to play TTS audio:', err);
            if (ttsButton) ttsButton.classList.remove('active');
        });
    };

    document.body.appendChild(fileInput);
    fileInput.click();
    document.body.removeChild(fileInput);
}

// ---------------------------------------------------------------------------
// File Attachment (stub — RAG placeholder)
// ---------------------------------------------------------------------------

function addFile() {
    console.log('Add file — not yet implemented');
}

// ---------------------------------------------------------------------------
// Narrator Event Flush Helper
// ---------------------------------------------------------------------------

async function flushPendingNarratorEvents() {
    if (!window.pendingNarratorEvents || window.pendingNarratorEvents.length === 0) return;
    for (const eventItem of window.pendingNarratorEvents) {
        // Sync to backend only — queueNarratorEvent already added to DOM
        if (window.go && window.go.main && window.go.main.App && window.go.main.App.AddSystemMessage) {
            try { await window.go.main.App.AddSystemMessage(eventItem.text); }
            catch (error) { console.error('Failed to sync narrator event to backend:', error); }
        }
    }
    window.pendingNarratorEvents = [];
}

// ---------------------------------------------------------------------------
// Abort Helper
// ---------------------------------------------------------------------------

async function handleAbort() {
    if (!window.go || !window.go.main || !window.go.main.App) return;
    try {
        // Always abort the active stream first, regardless of narrator state
        await window.go.main.App.AbortStream().catch(err => {
            console.error('Error aborting stream:', err);
        });

        // Drain any pending narrator events that were queued before the abort
        for (const eventItem of (window.pendingNarratorEvents || [])) {
            if (eventItem.type === 'narrator_think') {
                isInThinkMode  = true;
                thinkStartTime = Date.now();
                thinkContent   = eventItem.text;
                showThinkingIndicator && showThinkingIndicator();
            } else if (eventItem.type === 'narrator_speak') {
                isInThinkMode  = false;
                thinkStartTime = null;
                hideThinkingIndicator && hideThinkingIndicator();
                if (thinkContent) { addSystemMessage(`<em>${thinkContent}</em>`); thinkContent = ''; }
                addSystemMessage(eventItem.text);
            } else if (eventItem.type === 'narrator_event') {
                addSystemMessage(eventItem.text);
            }
            window.go.main.App.AddSystemMessage(eventItem.text).catch(err => {
                console.error('Failed to sync narrator event:', err);
            });
        }
        window.pendingNarratorEvents = [];
    } catch (error) {
        console.error('Failed to handle abort:', error);
    }
}

// ---------------------------------------------------------------------------
// Send Button & Enter Key
// ---------------------------------------------------------------------------

document.addEventListener('DOMContentLoaded', () => {
    const sendButton = document.querySelector('.send-button');
    if (sendButton) {
        sendButton.addEventListener('click', async () => {
            if (sendButton.dataset.isAbort === 'true') {
                await handleAbort();
                return;
            }

            const input   = document.querySelector('.input-field');
            const message = input ? input.value.trim() : '';
            if (!message) return;

            await flushPendingNarratorEvents();

            addUserMessage(message);
            if (input) input.value = '';

            let imageData = '';
            if (currentImage) {
                imageData = currentImage.split(',')[1];
                removeImage();
            }

            if (!window.go || !window.go.main || !window.go.main.App) {
                console.error('Wails runtime not available');
                addAssistantMessage('AI', 'Error: Backend not available. Please run this application through Wails (wails dev) to enable backend functionality.');
                return;
            }

            // Record user turn in backend history before inference begins
            if (window.go.main.App.AppendChatMessage) {
                window.go.main.App.AppendChatMessage('user', message, 'User').catch(() => {});
            }

            if (imageData) {
                await streamAssistantResponseFromBackendWithImage(message, imageData);
            } else {
                await streamMultiCharacterResponses(message);
            }

            // Persist history to disk after every complete turn
            if (window.go.main.App.SaveChat) {
                window.go.main.App.SaveChat().catch(err => {
                    console.warn('Post-turn save failed:', err);
                });
            }
        });
    }
});

document.addEventListener('DOMContentLoaded', () => {
    const inputField = document.querySelector('.input-field');
    if (inputField) {
        inputField.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                const sendButton = document.querySelector('.send-button');
                if (sendButton) sendButton.click();
            }
        });
    }
});
