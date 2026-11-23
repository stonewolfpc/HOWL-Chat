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

// The active conversation has one owner in the frontend: dataSync.currentChatId.
// It was previously read from a separate window.currentChatId that nothing ever
// assigned, so "Delete" could only ever report "No active chat session to delete".
function activeChatId() {
    if (window.dataSync && window.dataSync.currentChatId) {
        return window.dataSync.currentChatId;
    }
    return null;
}

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

    const chatId = activeChatId();
    if (!chatId || chatId === 'default') {
        showNotification('No active chat session to delete', 'warning');
        return;
    }

    // Deletion is best-effort: the chat leaves the list immediately, and the folder
    // is removed only if Windows allows it. Say where the files are before the user
    // confirms, rather than reporting a deletion failure afterwards.
    let directory = '';
    try {
        directory = await window.go.main.App.GetChatDirectory(chatId);
    } catch (error) {
        console.warn('Could not resolve the chat directory for the deletion notice:', error);
    }

    const escapeHtml = (value) => String(value).replace(/[&<>"']/g, (ch) => ({
        '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;'
    }[ch]));

    // The whole path is one block so every pixel of it is a hit target, not only
    // the glyphs of a wrapped inline link.
    const note = directory
        ? 'The chat folder is ' +
          '<a href="#" id="dc-open-dir" style="display:block;margin-top:6px;padding:7px 9px;' +
          'border-radius:6px;background:rgba(127,184,255,0.10);color:#7fb8ff;' +
          'text-decoration:underline;word-break:break-all;">' + escapeHtml(directory) + '</a>' +
          'If Windows blocks the move, the files stay there. Open the folder to check.'
        : '';

    window.DeleteConfirm.show({
        title: 'Delete current chat?',
        body: 'This chat leaves HOWL immediately. Its folder is moved to the Windows Recycle Bin, where you can normally restore it. World, scenario, and character data are preserved.',
        warning: 'The chat is removed from HOWL now. Windows keeps the folder in the Recycle Bin until you empty it.',
        note: note,
        onOpenDirectory: async () => {
            try {
                await window.go.main.App.OpenChatDirectory(chatId);
            } catch (error) {
                // A folder that cannot be opened is reported, never silent: the
                // whole point of the link is that the user can check by hand.
                console.error('Failed to open chat directory:', error);
                showNotification('Could not open the folder: ' + error, 'error');
            }
        },
        onConfirm: async () => {
            try {
                showLoading(true);
                await window.go.main.App.DeleteChat(chatId);
                showNotification('Chat session deleted successfully', 'success');

                window.dataSync.currentChatId = null;
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
