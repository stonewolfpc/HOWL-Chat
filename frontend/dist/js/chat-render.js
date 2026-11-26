/**
 * Chat Message Rendering
 *
 * Creates and appends chat message DOM elements for user,
 * assistant, and system (narrator) messages. Also handles
 * updating an existing assistant bubble with character metadata.
 *
 * Depends on: chat-constants.js
 *
 * @module chat-render
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// User Message
// ---------------------------------------------------------------------------

function addUserMessage(text) {
    const chatHistory = document.querySelector('.chat-messages');
    const messageDiv  = document.createElement('div');
    messageDiv.className = 'chat-message user';

    const msgId = 'msg-' + Date.now() + '-' + Math.random().toString(36).slice(2, 8);
    messageDiv.dataset.id        = msgId;
    messageDiv.dataset.timestamp = new Date().toISOString();

    let userName = 'You';
    const cd = (typeof currentChatData !== 'undefined') ? currentChatData
             : (window.sidebarState ? window.sidebarState.currentChatData : null);

    if (cd && cd.playerCharacter && cd.playerCharacter.name) {
        userName = cd.playerCharacter.name;
    } else {
        const userProfile = localStorage.getItem('howl-user-profile');
        if (userProfile) {
            try { userName = JSON.parse(userProfile).name || 'You'; } catch (e) { /* use default */ }
        }
    }

    messageDiv.innerHTML = `
        <div class="chat-bubble">
            <div class="chat-sender">${userName}</div>
            <div class="chat-text">${text}</div>
            <div class="chat-bubble-buttons">
                <button class="bubble-button icon-button" onclick="copyMessage(this)" title="Copy">${ICONS.copy}</button>
                <button class="bubble-button icon-button" onclick="editMessage(this)" title="Edit">${ICONS.edit}</button>
                <button class="bubble-button icon-button" onclick="undoMessage(this)" title="Undo (remove this and everything after)">${ICONS.undo}</button>
                <button class="bubble-button icon-button" onclick="deleteMessage(this)" title="Delete">${ICONS.delete}</button>
            </div>
        </div>
    `;
    chatHistory.appendChild(messageDiv);
    chatHistory.scrollTop = chatHistory.scrollHeight;
}

// ---------------------------------------------------------------------------
// Assistant Message
// ---------------------------------------------------------------------------

function addAssistantMessage(sender, text, characterId = null) {
    const chatHistory = document.querySelector('.chat-messages');
    const messageDiv  = document.createElement('div');
    messageDiv.className = 'chat-message assistant';

    const msgId = 'msg-' + Date.now() + '-' + Math.random().toString(36).slice(2, 8);
    messageDiv.dataset.id          = msgId;
    messageDiv.dataset.characterId = characterId || '';
    messageDiv.dataset.timestamp   = new Date().toISOString();

    const cd = (typeof currentChatData !== 'undefined') ? currentChatData
             : (window.sidebarState ? window.sidebarState.currentChatData : null);

    let character = null;
    if (cd && cd.characters && characterId) {
        character = cd.characters.find(c => c.id === characterId) || null;
    }

    // Resolve image — character.image may be a string, an object {data,path,url}, or absent
    const resolveImg = (val) => {
        if (!val) return '';
        if (typeof val === 'string') return val;
        if (typeof val === 'object') return val.data || val.url || val.path || val.thumbnail || '';
        return '';
    };
    const imgSrc = resolveImg(character?.image) || resolveImg(character?.avatar)
                || resolveImg(character?._resolvedImage) || '';

    let senderHTML = `<div class="chat-sender">`;
    if (imgSrc) {
        senderHTML += `<img src="${imgSrc}" alt="${character?.name || sender}" class="character-thumbnail" onerror="this.style.display='none'">` ;
    }
    senderHTML += `<span class="character-name">${character ? character.name : sender}</span></div>`;

    messageDiv.innerHTML = `
        <div class="chat-bubble">
            ${senderHTML}
            <div class="chat-text">${text}</div>
            <div class="chat-bubble-buttons">
                <button class="bubble-button icon-button" onclick="window.AudioService && AudioService.speakMessage(this)" title="Speak this message">${ICONS.speaker}</button>
                <button class="bubble-button icon-button" onclick="regenerateMessage(this)" title="Regenerate">${ICONS.regenerate}</button>
                <button class="bubble-button icon-button" onclick="previousMessage(this)" title="Previous" id="prev-btn">${ICONS.previous}</button>
                <button class="bubble-button icon-button" onclick="nextMessage(this)" title="Next" id="next-btn">${ICONS.next}</button>
                <button class="bubble-button icon-button" onclick="continueMessage(this)" title="Continue">${ICONS.continue}</button>
                <button class="bubble-button icon-button" onclick="editMessage(this)" title="Edit">${ICONS.edit}</button>
                <button class="bubble-button icon-button" onclick="copyMessage(this)" title="Copy">${ICONS.copy}</button>
                <button class="bubble-button icon-button" onclick="undoMessage(this)" title="Undo (remove this and everything after)">${ICONS.undo}</button>
                <button class="bubble-button icon-button" onclick="deleteMessage(this)" title="Delete">${ICONS.delete}</button>
            </div>
        </div>
    `;
    chatHistory.appendChild(messageDiv);
    chatHistory.scrollTop = chatHistory.scrollHeight;
    updateNavigationButtons();
}

// ---------------------------------------------------------------------------
// System / Narrator Message
// ---------------------------------------------------------------------------

function addSystemMessage(text) {
    const chatHistory = document.querySelector('.chat-messages');
    const messageDiv  = document.createElement('div');
    messageDiv.className = 'chat-message system';

    const msgId = 'msg-' + Date.now() + '-' + Math.random().toString(36).slice(2, 8);
    messageDiv.dataset.id        = msgId;
    messageDiv.dataset.timestamp = new Date().toISOString();

    messageDiv.innerHTML = `
        <div class="chat-bubble">
            <div class="chat-sender">Narrator</div>
            <div class="chat-text">${text}</div>
        </div>
    `;
    chatHistory.appendChild(messageDiv);
    chatHistory.scrollTop = chatHistory.scrollHeight;
}

function queueNarratorEvent(text) {
    if (!window.pendingNarratorEvents) window.pendingNarratorEvents = [];
    window.pendingNarratorEvents.push({ text });
    addSystemMessage(text);
}

// ---------------------------------------------------------------------------
// Bubble Character Info Update (called during streaming via CHARACTER_INFO token)
// ---------------------------------------------------------------------------

function updateBubbleWithCharacterInfo(bubble, characterName, personality, backstory) {
    const senderElement = bubble.querySelector('.chat-sender');
    if (!senderElement) return;

    const cd = (typeof currentChatData !== 'undefined') ? currentChatData
             : (window.sidebarState ? window.sidebarState.currentChatData : null);
    const character = cd && cd.characters ? cd.characters.find(c => c.name === characterName) : null;
    const avatar    = character ? (character.avatar || '') : '';

    senderElement.innerHTML = `
        <div class="character-sender-info">
            ${avatar
                ? `<img src="${avatar}" alt="${characterName}" class="character-avatar-small" onerror="this.style.display='none'; this.nextElementSibling.style.display='flex';">
                   <div class="character-avatar-placeholder" style="display:none;">👤</div>`
                : `<div class="character-avatar-placeholder">👤</div>`
            }
            <div class="character-name">${characterName}</div>
        </div>
        <span class="typing-indicator"><span></span><span></span><span></span></span>
    `;

    bubble.dataset.characterName        = characterName;
    bubble.dataset.characterPersonality = personality;
    bubble.dataset.characterBackstory   = backstory;
}
