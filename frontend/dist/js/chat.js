/**
 * Chat — Orchestrator & Session Operations
 *
 * Wires together the modular chat components:
 *   chat-constants.js  — ICONS, regex, Wails readiness
 *   chat-render.js     — addUserMessage / addAssistantMessage / addSystemMessage
 *   chat-stream.js     — all streaming logic and think-block handling
 *   chat-actions.js    — copy, edit, delete, regenerate, continue, navigate
 *   chat-input.js      — send/abort handler, image/audio/file attachments
 *   chat-model.js      — loadModel, unloadModel, checkModelStatus
 *
 * This file owns:
 *   - newChat / deleteChatSession session lifecycle
 *   - clearChatHistory DOM helper
 *   - window.chat public API surface
 *
 * @module chat
 * @version 2.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Session Lifecycle
// ---------------------------------------------------------------------------

window.currentChatId = null;

window.newChat = async function () {
    if (!window.go || !window.go.main || !window.go.main.App) {
        console.error('Wails runtime not available');
        showNotification('Backend not available', 'error');
        return;
    }

    if (typeof dataSync === 'undefined' || typeof dataSync.createNewChat !== 'function') {
        console.error('dataSync.createNewChat is not available');
        showNotification('Chat sync module not loaded', 'error');
        return;
    }

    try {
        showLoading(true);
        const chatName = `Chat ${new Date().toLocaleString()}`;
        const chatMeta = await dataSync.createNewChat(chatName);

        if (!chatMeta) throw new Error('Backend did not return chat metadata');

        console.log('New chat created:', chatMeta);
        showLoading(false);
        showNotification('New chat created', 'success');
    } catch (error) {
        console.error('Failed to create new chat:', error);
        showLoading(false);
        showNotification('Failed to create new chat: ' + error, 'error');
    }
};

async function deleteActiveChatSession() {
    if (!window.go || !window.go.main || !window.go.main.App || !window.go.main.App.DeleteChat) {
        showNotification('Backend not available', 'error');
        return;
    }

    if (!window.currentChatId || window.currentChatId === 'default') {
        showNotification('No active chat session to delete', 'warning');
        return;
    }

    window.DeleteConfirm.show({
        title: 'Delete current chat?',
        body: 'This chat\'s history and automated lore entries will be permanently deleted. World, scenario, and character data are preserved.',
        onConfirm: async () => {
            try {
                showLoading(true);
                await window.go.main.App.DeleteChat(window.currentChatId);
                showNotification('Chat session deleted successfully', 'success');

                window.currentChatId = null;
                if (window.sidebarState) window.sidebarState.resetToDefaultState();

                const chatHistory = document.querySelector('.chat-messages');
                if (chatHistory) chatHistory.innerHTML = '';

                const characterList = document.querySelector('.character-list');
                if (characterList) characterList.innerHTML = '';

                const worldThumbnail = document.getElementById('world-image');
                const thumbnailPlaceholder = document.querySelector('.thumbnail-placeholder');
                if (worldThumbnail && thumbnailPlaceholder) {
                    worldThumbnail.style.display    = 'none';
                    thumbnailPlaceholder.style.display = 'block';
                }

                hideLoading();
                showNotification('Chat session deleted. Starting new session...', 'info');
                setTimeout(() => { window.location.reload(); }, 1000);

            } catch (error) {
                console.error('Failed to delete chat session:', error);
                showNotification('Failed to delete chat session', 'error');
            } finally {
                hideLoading();
            }
        }
    });
}

// ---------------------------------------------------------------------------
// DOM Helper
// ---------------------------------------------------------------------------

function clearChatHistory() {
    const chatHistory = document.querySelector('.chat-messages');
    if (chatHistory) chatHistory.innerHTML = '';
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

window.chat = {
    newChat: window.newChat,
    deleteChatSession: deleteActiveChatSession,
    clearChatHistory
};
