/**
 * Chat Message Actions
 *
 * All post-render interactions on chat bubbles:
 * copy, edit (in-place), save/cancel edit, delete,
 * regenerate, continue, prev/next navigation,
 * and the "generate banner" that appears after deletion.
 *
 * Depends on: chat-render.js, chat-stream.js
 *
 * @module chat-actions
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Copy
// ---------------------------------------------------------------------------

function copyMessage(button) {
    const bubble = button.closest('.chat-bubble');
    const text   = bubble.querySelector('.chat-text').textContent;
    navigator.clipboard.writeText(text).then(() => {
        const original = button.innerHTML;
        button.textContent = 'Copied!';
        setTimeout(() => { button.innerHTML = original; }, 1500);
    });
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

function deleteMessage(button) {
    const messageDiv  = button.closest('.chat-message');
    const isAssistant = messageDiv.classList.contains('assistant');
    messageDiv.remove();
    _syncDomHistoryToBackend();
    if (isAssistant) showGenerateBanner();
    updateNavigationButtons();
}

// ---------------------------------------------------------------------------
// Shared: sync the current DOM message list to backend history
// ---------------------------------------------------------------------------

async function _syncDomHistoryToBackend() {
    if (!window.go?.main?.App?.ClearChatHistory) return;
    try {
        await window.go.main.App.ClearChatHistory();
        if (!window.go.main.App.AppendChatMessage) return;
        const msgs = document.querySelectorAll('.chat-messages .chat-message');
        for (const msgDiv of msgs) {
            const role   = msgDiv.classList.contains('user')      ? 'user'
                         : msgDiv.classList.contains('assistant') ? 'assistant'
                         : 'system';
            const text   = msgDiv.querySelector('.chat-text')?.textContent?.trim() || '';
            const sender = msgDiv.querySelector('.chat-sender')?.textContent?.trim() || '';
            if (text) {
                await window.go.main.App.AppendChatMessage(role, text, sender).catch(() => {});
            }
        }
    } catch (err) {
        console.warn('Failed to sync DOM history to backend:', err);
    }
}

// ---------------------------------------------------------------------------
// Shared: prune a message and everything after it, then sync backend
// ---------------------------------------------------------------------------

function _pruneMessagesAfter(startDiv, inclusive = true) {
    const chatHistory = document.querySelector('.chat-messages');
    if (!chatHistory) return;
    const all = Array.from(chatHistory.querySelectorAll('.chat-message'));
    const idx = all.indexOf(startDiv);
    if (idx === -1) return;
    const removeFrom = inclusive ? idx : idx + 1;
    for (let i = all.length - 1; i >= removeFrom; i--) {
        all[i].remove();
    }
    _syncDomHistoryToBackend();
    updateNavigationButtons();
}

// ---------------------------------------------------------------------------
// Undo — remove this message and everything after it (no re-send)
// ---------------------------------------------------------------------------

function undoMessage(button) {
    const messageDiv = button.closest('.chat-message');
    _pruneMessagesAfter(messageDiv, true);
}

// ---------------------------------------------------------------------------
// In-Place Edit
// On confirm: update text, remove everything after this message, re-send if user turn
// ---------------------------------------------------------------------------

function editMessage(button) {
    const bubble      = button.closest('.chat-bubble');
    const messageDiv  = button.closest('.chat-message');
    const textElement = bubble.querySelector('.chat-text');
    const currentText = textElement.textContent;
    const isUser      = messageDiv.classList.contains('user');

    // Prevent double-editing
    if (bubble.querySelector('.chat-text-edit')) return;

    const textarea = document.createElement('textarea');
    textarea.className = 'chat-text-edit';
    textarea.value = currentText;
    textarea.rows  = Math.min(10, Math.max(3, Math.ceil(currentText.length / 60)));
    textarea.dataset.originalText = currentText;
    textarea.style.width = '100%';
    textarea.style.boxSizing = 'border-box';

    textElement.style.display = 'none';
    bubble.insertBefore(textarea, textElement.nextSibling);
    
    // Add Save/Cancel buttons
    const buttonContainer = document.createElement('div');
    buttonContainer.className = 'chat-edit-buttons';
    buttonContainer.innerHTML = `
        <button class="chat-edit-save" title="Save (Ctrl+Enter)">✓ Save</button>
        <button class="chat-edit-cancel" title="Cancel (Escape)">✕ Cancel</button>
    `;
    bubble.insertBefore(buttonContainer, textarea.nextSibling);
    
    textarea.focus();
    textarea.select();

    let committed = false;

    const confirmEdit = () => {
        if (committed) return;
        committed = true;

        const newText = textarea.value.trim();
        const buttonContainer = bubble.querySelector('.chat-edit-buttons');
        if (buttonContainer) buttonContainer.remove();
        textarea.remove();
        textElement.style.display = 'block';

        if (!newText || newText === currentText) return; // nothing changed

        textElement.textContent = newText;

        // Remove all messages after this one, sync backend
        _pruneMessagesAfter(messageDiv, false);

        // Re-send as user message if this is a user bubble
        if (isUser) {
            if (typeof hideGenerateBanner === 'function') hideGenerateBanner();
            if (window.go?.main?.App?.AppendChatMessage) {
                window.go.main.App.AppendChatMessage('user', newText, 'User').catch(() => {});
            }
            if (typeof streamMultiCharacterResponses === 'function') {
                streamMultiCharacterResponses(newText).catch((err) => {
                    console.error('Failed to stream response after edit:', err);
                    showNotification('Failed to generate response: ' + (err.message || err), 'error', 5000);
                });
            }
        }
    };

    const cancelEdit = () => {
        if (committed) return;
        committed = true;
        const buttonContainer = bubble.querySelector('.chat-edit-buttons');
        if (buttonContainer) buttonContainer.remove();
        textarea.remove();
        textElement.style.display = 'block';
        textElement.textContent = currentText;
    };

    // Add button event listeners
    const saveBtn = buttonContainer.querySelector('.chat-edit-save');
    const cancelBtn = buttonContainer.querySelector('.chat-edit-cancel');
    if (saveBtn) saveBtn.addEventListener('click', confirmEdit);
    if (cancelBtn) cancelBtn.addEventListener('click', cancelEdit);

    // Confirm on Ctrl+Enter; cancel on Escape; blur = confirm
    textarea.addEventListener('keydown', (e) => {
        if (e.key === 'Enter' && e.ctrlKey) { e.preventDefault(); confirmEdit(); }
        if (e.key === 'Escape')             { e.preventDefault(); cancelEdit();  }
    });
    textarea.addEventListener('blur', () => {
        // Short delay so that a keydown=Escape fires before blur
        setTimeout(confirmEdit, 150);
    });
}

// ---------------------------------------------------------------------------
// Save / Cancel Edit (legacy button-based path — kept for compatibility)
// ---------------------------------------------------------------------------

function saveEdit(button, isUserMessage) {
    const bubble      = button.closest('.chat-bubble');
    const messageDiv  = button.closest('.chat-message');
    const textarea    = bubble.querySelector('.chat-text-edit');
    const editedText  = textarea.value;

    if (isUserMessage) {
        const chatHistory = document.querySelector('.chat-messages');
        const messages    = Array.from(chatHistory.querySelectorAll('.chat-message'));
        const messageIndex = messages.indexOf(messageDiv);

        const textElement = bubble.querySelector('.chat-text');
        textElement.textContent   = editedText;
        textElement.style.display = 'block';
        textarea.style.display    = 'none';
        bubble.removeChild(textarea);

        for (let i = messages.length - 1; i > messageIndex; i--) { messages[i].remove(); }

        const input = document.querySelector('.input-field');
        if (input) input.value = '';

        clearChatHistory();
        if (typeof hideGenerateBanner === 'function') hideGenerateBanner();
        streamAssistantResponseFromBackend(editedText);
    } else {
        const newTextElement = document.createElement('div');
        newTextElement.className   = 'chat-text';
        newTextElement.textContent = editedText;
        textarea.replaceWith(newTextElement);
    }

    const buttonsContainer = bubble.querySelector('.chat-bubble-buttons');
    if (buttonsContainer && buttonsContainer.dataset.originalButtons) {
        buttonsContainer.innerHTML = buttonsContainer.dataset.originalButtons;
    }
    updateNavigationButtons();
}

function cancelEdit(button) {
    const bubble       = button.closest('.chat-bubble');
    const textarea     = bubble.querySelector('.chat-text-edit');
    const originalText = textarea.dataset.originalText || textarea.value;

    const newTextElement = document.createElement('div');
    newTextElement.className   = 'chat-text';
    newTextElement.textContent = originalText;
    textarea.replaceWith(newTextElement);

    const buttonsContainer = bubble.querySelector('.chat-bubble-buttons');
    if (buttonsContainer && buttonsContainer.dataset.originalButtons) {
        buttonsContainer.innerHTML = buttonsContainer.dataset.originalButtons;
    }
}

// ---------------------------------------------------------------------------
// Regenerate
// ---------------------------------------------------------------------------

function regenerateMessage(button) {
    try {
        const chatHistory       = document.querySelector('.chat-messages');
        const assistantMessages = Array.from(chatHistory.querySelectorAll('.chat-message.assistant'));

        if (assistantMessages.length === 0) {
            showNotification('No AI messages to regenerate', 'warning');
            return;
        }

        const lastAIMessage = assistantMessages[assistantMessages.length - 1];
        if (lastAIMessage) lastAIMessage.remove();

        const userMessages  = Array.from(chatHistory.querySelectorAll('.chat-message.user'));
        const lastUserMsg   = userMessages[userMessages.length - 1];
        if (lastUserMsg) {
            const userText = lastUserMsg.querySelector('.chat-text').textContent;
            if (typeof hideGenerateBanner === 'function') hideGenerateBanner();
            
            // Add error handling around the stream call
            streamAssistantResponseFromBackend(userText).catch((err) => {
                console.error('Regenerate failed:', err);
                showNotification('Failed to regenerate message: ' + (err.message || err), 'error', 10000);
            });
        } else {
            showNotification('No user message found to regenerate from', 'warning');
        }
    } catch (err) {
        console.error('Regenerate error:', err);
        showNotification('Error during regenerate: ' + (err.message || err), 'error', 10000);
    }
}

// ---------------------------------------------------------------------------
// Continue
// ---------------------------------------------------------------------------

function continueMessage(button) {
    const chatHistory       = document.querySelector('.chat-messages');
    const assistantMessages = Array.from(chatHistory.querySelectorAll('.chat-message.assistant'));

    if (assistantMessages.length === 0) {
        console.log('No AI messages to continue from');
        return;
    }

    const lastAssistant = assistantMessages[assistantMessages.length - 1];
    const lastText      = lastAssistant.querySelector('.chat-text').textContent;
    if (typeof hideGenerateBanner === 'function') hideGenerateBanner();
    streamAssistantResponseFromBackend(lastText);
}

// ---------------------------------------------------------------------------
// Navigation (Prev / Next)
// ---------------------------------------------------------------------------

function previousMessage(button) {
    const chatHistory       = document.querySelector('.chat-messages');
    const assistantMessages = Array.from(chatHistory.querySelectorAll('.chat-message.assistant'));
    const currentIndex      = assistantMessages.findIndex(msg => msg.contains(button.closest('.chat-bubble')));
    if (currentIndex > 0) {
        const previous = assistantMessages[currentIndex - 1];
        chatHistory.scrollTop = previous.offsetTop - chatHistory.offsetTop;
    }
}

function nextMessage(button) {
    const chatHistory       = document.querySelector('.chat-messages');
    const assistantMessages = Array.from(chatHistory.querySelectorAll('.chat-message.assistant'));
    const currentIndex      = assistantMessages.findIndex(msg => msg.contains(button.closest('.chat-bubble')));
    if (currentIndex < assistantMessages.length - 1) {
        const next = assistantMessages[currentIndex + 1];
        chatHistory.scrollTop = next.offsetTop - chatHistory.offsetTop;
    }
}

function updateNavigationButtons() {
    const chatHistory       = document.querySelector('.chat-messages');
    if (!chatHistory) return;
    const assistantMessages = Array.from(chatHistory.querySelectorAll('.chat-message.assistant'));
    const hasPrevious       = assistantMessages.length > 1;

    chatHistory.querySelectorAll('#prev-btn').forEach(btn => btn.classList.toggle('disabled', !hasPrevious));
    chatHistory.querySelectorAll('#next-btn').forEach(btn => btn.classList.add('disabled'));
}

// ---------------------------------------------------------------------------
// Generate Banner
// ---------------------------------------------------------------------------

function showGenerateBanner() {
    const chatPanel = document.querySelector('.chat-panel');
    if (!chatPanel || document.getElementById('generate-banner')) return;

    const banner = document.createElement('div');
    banner.id        = 'generate-banner';
    banner.className = 'generate-banner';
    banner.innerHTML = 'Generate AI Response';
    banner.onclick   = () => {
        const userMessages  = Array.from(document.querySelectorAll('.chat-message.user'));
        const lastUserMsg   = userMessages[userMessages.length - 1];
        if (lastUserMsg) {
            streamAssistantResponseFromBackend(lastUserMsg.querySelector('.chat-text').textContent);
        }
        banner.remove();
    };
    chatPanel.appendChild(banner);
}

function hideGenerateBanner() {
    const banner = document.getElementById('generate-banner');
    if (banner) banner.remove();
}
