/**
 * Chat Streaming
 *
 * All token-streaming logic for single and multi-character
 * responses, think-block handling, lorebook sync, and the
 * simulated local streamer (development fallback).
 *
 * Depends on: chat-constants.js, chat-render.js
 *
 * @module chat-stream
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Module-level streaming state
// ---------------------------------------------------------------------------

let currentStreamingBubble = null;
let thinkStartTime  = null;
let thinkContent    = '';
let mainContent     = '';
let isInThinkMode   = false;

// ---------------------------------------------------------------------------
// Lorebook Sync
// ---------------------------------------------------------------------------

async function syncLorebooks() {
    if (!window.go || !window.go.main || !window.go.main.App) {
        console.error('Wails runtime not available');
        return;
    }
    try {
        const lorebooks = JSON.parse(localStorage.getItem('howl-chat-lorebooks') || '[]');
        window.go.main.App.SetLorebooks(Array.isArray(lorebooks) ? lorebooks : []).catch(error => {
            console.error('Failed to sync lorebooks:', error);
        });
    } catch (error) {
        console.error('Failed to sync lorebooks:', error);
    }
}

// Keep legacy alias used by streamAssistantResponseFromBackendWithImage
const syncLorebooksToBackend = syncLorebooks;

// ---------------------------------------------------------------------------
// Character Context Builder (for single-stream fallback)
// ---------------------------------------------------------------------------

function getRecentSystemMessages(maxCount = 3) {
    const chatHistory = document.querySelector('.chat-messages');
    if (!chatHistory) return [];
    const systemEls = chatHistory.querySelectorAll('.chat-message.system');
    const messages = [];
    const start = Math.max(0, systemEls.length - maxCount);
    for (let i = start; i < systemEls.length; i++) {
        const text = systemEls[i].querySelector('.chat-text');
        if (text) messages.push(text.textContent.trim());
    }
    return messages;
}

function isCharacterActive(character) {
    return character.status !== 'dead' && character.status !== 'stepped-out';
}

function buildCharacterContext() {
    const cd = (typeof currentChatData !== 'undefined') ? currentChatData
             : (window.sidebarState ? window.sidebarState.currentChatData : null);
    if (!cd || !cd.characters || cd.characters.length === 0) return '';

    const activeChars = cd.characters.filter(isCharacterActive);
    if (activeChars.length === 0) return '';

    let context = 'Active Characters:\n';
    activeChars.forEach((character, index) => {
        context += `${index + 1}. ${character.name}`;
        if (character.description) context += ` - ${character.description}`;
        if (character.personality) context += `\n   Personality: ${character.personality}`;
        if (character.background)  context += `\n   Background: ${character.background}`;
        context += '\n';
    });

    const inactive = cd.characters.filter(c => !isCharacterActive(c));
    if (inactive.length > 0) {
        context += '\n[INACTIVE / ABSENT]\n';
        inactive.forEach(c => {
            const reason = c.stepOutReason ? ` (${c.stepOutReason})` : '';
            context += `- ${c.name} is ${c.status}${reason}\n`;
        });
    }

    return context;
}

function buildCharacterSpecificMessage(userMessage, character, priorResponses = []) {
    let msg = `You are ${character.name}, a distinct participant in a group roleplay conversation.\n`;
    msg += `You speak only when someone addresses you directly by name, role, or context, or when it is naturally your turn to respond.\n`;
    msg += `If others are speaking to each other, you listen quietly and stay in character, reacting only when your involvement makes sense.\n\n`;
    if (character.description) msg += `Description: ${character.description}\n`;
    if (character.personality)  msg += `Personality: ${character.personality}\n`;
    if (character.background)   msg += `Background: ${character.background}\n`;

    // Recent system events (death, step-out) so characters are aware
    const recentSystem = getRecentSystemMessages();
    if (recentSystem.length > 0) {
        msg += `\n[RECENT EVENTS]\n`;
        for (const event of recentSystem) {
            msg += `- ${event}\n`;
        }
    }

    // Group context: what others already said this turn, so this character can respond as if listening
    if (priorResponses.length > 0) {
        msg += `\n[GROUP CONTEXT — what others said just now]\n`;
        for (const r of priorResponses) {
            msg += `${r.name}: ${r.text}\n`;
        }
    }

    // Note inactive characters so this character knows who is present
    const cd = (typeof currentChatData !== 'undefined') ? currentChatData
             : (window.sidebarState ? window.sidebarState.currentChatData : null);
    if (cd && cd.characters) {
        const inactive = cd.characters.filter(c => c.id !== character.id && !isCharacterActive(c));
        if (inactive.length > 0) {
            msg += `\n[NOTE: The following characters are no longer present]\n`;
            for (const c of inactive) {
                const reason = c.stepOutReason ? ` — ${c.stepOutReason}` : '';
                msg += `- ${c.name} is ${c.status}${reason}\n`;
            }
        }
    }

    // Scoped lorebook entries: world → all; scenario → all; character → owner only
    const lorebooks = (() => {
        try { return JSON.parse(localStorage.getItem('howl-chat-lorebooks') || '[]'); }
        catch (_) { return []; }
    })();
    const activeLore = lorebooks.filter(e => e.enabled !== false);
    const worldLore   = activeLore.filter(e => e.scope === 'world'   || e.type === 'world');
    const scenarioLore= activeLore.filter(e => e.scope === 'scenario'|| e.type === 'scenario');
    const charLore    = activeLore.filter(e =>
        (e.scope === 'character' || e.type === 'character') &&
        (e.linkedCharacters?.includes?.(character.id) || e.ownerId === character.id || e.ownerId === character.name)
    );
    const scopedLore = [...worldLore, ...scenarioLore, ...charLore];
    if (scopedLore.length > 0) {
        msg += `\n[RELEVANT LORE]\n`;
        for (const entry of scopedLore.slice(0, 6)) {
            msg += `- ${entry.title || 'Untitled'}: ${entry.content?.slice?.(0, 300) || entry.content || ''}\n`;
        }
    }

    // Automated memories: character-scoped memories for this character only
    const charMemories = window.MemoryStore
        ? window.MemoryStore.getCharacterMemories(character.id)
        : [];
    if (charMemories.length > 0) {
        msg += `\n[MEMORIES]\n`;
        for (const mem of charMemories.slice(0, 4)) {
            msg += `- ${mem.title || 'Memory'}: ${mem.summary || mem.content?.slice?.(0, 200) || ''}\n`;
        }
    }

    msg += `\nUser: ${userMessage}\nRespond as ${character.name} only, aware of what was just said above.`;
    return msg;
}

// ---------------------------------------------------------------------------
// Think Area Rendering
// ---------------------------------------------------------------------------

function updateThinkArea(thinkArea, content, startTime) {
    if (!content) {
        thinkArea.innerHTML = '';
        thinkArea.style.display = 'none';
        return;
    }
    const elapsed = startTime ? ((Date.now() - startTime) / 1000).toFixed(1) : 0;
    thinkArea.style.display = 'block';
    thinkArea.innerHTML = `
        <div class="thinking-toggle" onclick="toggleThink(this)">
            <span>Show reasoning (${elapsed}s)</span>
        </div>
        <div class="thinking-block">
            <pre>${content}</pre>
        </div>
    `;
}

function toggleThink(toggleElement) {
    const thinkArea   = toggleElement.closest('.chat-think');
    const block       = thinkArea.querySelector('.thinking-block');
    const toggleText  = toggleElement.querySelector('span');
    if (block.classList.contains('open')) {
        block.classList.remove('open');
        toggleText.textContent = toggleText.textContent.replace('Hide', 'Show');
    } else {
        block.classList.add('open');
        toggleText.textContent = toggleText.textContent.replace('Show', 'Hide');
    }
}

// ---------------------------------------------------------------------------
// Send Button State Helpers
// ---------------------------------------------------------------------------

function setSendButtonAbort() {
    const sendButton = document.querySelector('.send-button');
    if (sendButton) { sendButton.textContent = 'Abort'; sendButton.dataset.isAbort = 'true'; }
}

function setSendButtonSend() {
    const sendButton = document.querySelector('.send-button');
    if (sendButton) { sendButton.textContent = 'Send'; sendButton.dataset.isAbort = 'false'; }
}

function cleanupStreamListeners() {
    if (window.runtime) {
        window.runtime.EventsOff('chat:chunk');
        window.runtime.EventsOff('chat:complete');
        window.runtime.EventsOff('chat:error');
        window.runtime.EventsOff('chat:aborted');
    }
}

// ---------------------------------------------------------------------------
// Multi-Character Response Orchestration
// ---------------------------------------------------------------------------

async function streamMultiCharacterResponses(userMessage) {
    const cd         = (typeof currentChatData !== 'undefined') ? currentChatData
                     : (window.sidebarState ? window.sidebarState.currentChatData : null);
    const characters = (cd && cd.characters) ? cd.characters : [];

    setSendButtonAbort();

    if (characters.length === 0) {
        console.log('No characters loaded, using generic AI response');
        if (typeof hideGenerateBanner === 'function') hideGenerateBanner();
        await streamAssistantResponseFromBackend(userMessage);
        setSendButtonSend();
        return;
    }

    // Exclude dead or stepped-out characters from the response queue
    const activeChars = characters.filter(isCharacterActive);

    // Sort by metadata.priority (1 = responds first, 5 = responds last, default 3)
    const ordered = [...activeChars].sort((a, b) => {
        const pa = a.metadata?.priority ?? 3;
        const pb = b.metadata?.priority ?? 3;
        return pa - pb;
    });

    const skipped = characters.length - activeChars.length;
    if (skipped > 0) {
        console.log(`Skipped ${skipped} inactive character(s) from response queue`);
    }

    if (activeChars.length === 0) {
        console.log('All characters are inactive — no responses generated');
        if (typeof hideGenerateBanner === 'function') hideGenerateBanner();
        addSystemMessage('All characters are currently inactive. Add new characters or revive existing ones to continue.');
        setSendButtonSend();
        return;
    }

    console.log(`Starting multi-character response sequence with ${ordered.length} character(s) in priority order: ${ordered.map(c => `${c.name}(${c.metadata?.priority ?? 3})`).join(', ')}`);
    window.CharacterQueue.initialize(ordered);

    // Accumulate what each character says so subsequent characters hear the full room
    const priorResponses = [];

    let character = window.CharacterQueue.getCurrentCharacter();
    while (character) {
        console.log(`Generating response for: ${character.name} (priority ${character.metadata?.priority ?? 3})`);
        const characterSpecificMessage = buildCharacterSpecificMessage(userMessage, character, priorResponses);
        let responseText;
        try {
            responseText = await streamSingleCharacterResponse(characterSpecificMessage, character);
        } catch (err) {
            console.error(`Unexpected error for ${character.name}:`, err);
            responseText = null;
        }
        // null means the stream was aborted or errored — stop the queue immediately
        if (responseText === null) {
            console.log('Stream stopped — stopping character queue');
            break;
        }
        priorResponses.push({ name: character.name, text: responseText });
        window.CharacterQueue.markResponded(character.id);
        character = window.CharacterQueue.rotateToNext();
    }

    console.log('Character queue complete');
    window.CharacterQueue.resetRound();
    setSendButtonSend();
}

// ---------------------------------------------------------------------------
// Single Character Stream
// ---------------------------------------------------------------------------

async function streamSingleCharacterResponse(message, character) {
    if (!isCharacterActive(character)) {
        console.log(`Skipping ${character.name} — status is ${character.status || 'unknown'}`);
        return '';
    }
    return new Promise((resolve, reject) => {
        const chatHistory = document.querySelector('.chat-messages');

        const messageDiv = document.createElement('div');
        messageDiv.className = 'chat-message assistant streaming';
        messageDiv.dataset.characterId   = character.id;
        messageDiv.dataset.characterName = character.name;

        const bubble = document.createElement('div');
        bubble.className = 'chat-bubble';
        bubble.innerHTML = `
            <div class="chat-sender">${character.name}<span class="typing-indicator"><span></span><span></span><span></span></span></div>
            <div class="chat-text"></div>
            <div class="chat-think"></div>
            <div class="chat-bubble-buttons">
                <button class="bubble-button icon-button" onclick="regenerateMessage(this)" title="Regenerate">${ICONS.regenerate}</button>
                <button class="bubble-button icon-button" onclick="previousMessage(this)" title="Previous" id="prev-btn">${ICONS.previous}</button>
                <button class="bubble-button icon-button" onclick="nextMessage(this)" title="Next" id="next-btn">${ICONS.next}</button>
                <button class="bubble-button icon-button" onclick="continueMessage(this)" title="Continue">${ICONS.continue}</button>
                <button class="bubble-button icon-button" onclick="editMessage(this)" title="Edit">${ICONS.edit}</button>
                <button class="bubble-button icon-button" onclick="copyMessage(this)" title="Copy">${ICONS.copy}</button>
                <button class="bubble-button icon-button" onclick="deleteMessage(this)" title="Delete">${ICONS.delete}</button>
            </div>
        `;

        messageDiv.appendChild(bubble);
        chatHistory.appendChild(messageDiv);
        chatHistory.scrollTop = chatHistory.scrollHeight;

        const textElement = bubble.querySelector('.chat-text');
        const thinkArea   = bubble.querySelector('.chat-think');

        let streamContent   = '';
        let localThink      = '';
        let localThinkMode  = false;
        let localThinkStart = null;

        const onChunk = (chunk) => {
            if (THINK_START_REGEX.test(chunk)) {
                localThinkMode  = true;
                localThinkStart = Date.now();
                chunk = chunk.replace(THINK_START_REGEX, '');
            }
            if (THINK_END_REGEX.test(chunk)) {
                localThinkMode = false;
                chunk = chunk.replace(THINK_END_REGEX, '');
                if (localThink) updateThinkArea(thinkArea, localThink, localThinkStart);
            }
            if (localThinkMode) {
                localThink += chunk;
            } else {
                streamContent += chunk;
                textElement.textContent = streamContent;
            }
            chatHistory.scrollTop = chatHistory.scrollHeight;
        };

        const onComplete = () => {
            console.log(`${character.name} response complete`);
            const indicator = bubble.querySelector('.typing-indicator');
            if (indicator) indicator.remove();
            messageDiv.classList.remove('streaming');
            if (localThink) updateThinkArea(thinkArea, localThink, localThinkStart);
            cleanup();
            updateNavigationButtons();
            if (typeof hideGenerateBanner === 'function') hideGenerateBanner();

            if (window.settingsState?.current?.general?.ttsAutoplay && window.AudioService?.playTTS) {
                window.AudioService.playTTS(streamContent, character.id);
            }

            resolve(streamContent); // return completed text for group context accumulation
        };

        const onError = (err) => {
            console.error(`Error streaming ${character.name} response:`, err);
            textElement.textContent += `\n[Error: ${err}]`;
            const indicator = bubble.querySelector('.typing-indicator');
            if (indicator) indicator.remove();
            messageDiv.classList.remove('streaming');
            cleanup();
            if (typeof hideGenerateBanner === 'function') hideGenerateBanner();
            resolve(null); // treat error as abort — stops queue, resets send button
        };

        const onAborted = () => {
            console.log(`Stream aborted for ${character.name}`);
            const indicator = bubble.querySelector('.typing-indicator');
            if (indicator) indicator.remove();
            messageDiv.classList.remove('streaming');
            cleanup();
            if (typeof hideGenerateBanner === 'function') hideGenerateBanner();
            resolve(null); // null signals the caller to stop the queue
        };

        function cleanup() {
            if (!window.runtime) return;
            window.runtime.EventsOff(`chat:chunk:${character.id}`);
            window.runtime.EventsOff(`chat:complete:${character.id}`);
            window.runtime.EventsOff(`chat:error:${character.id}`);
            window.runtime.EventsOff(`chat:aborted:${character.id}`);
            window.runtime.EventsOff('chat:aborted');
        }

        if (window.runtime) {
            window.runtime.EventsOn(`chat:chunk:${character.id}`,    onChunk);
            window.runtime.EventsOn(`chat:complete:${character.id}`, onComplete);
            window.runtime.EventsOn(`chat:error:${character.id}`,    onError);
            window.runtime.EventsOn(`chat:aborted:${character.id}`,  onAborted);
            window.runtime.EventsOn('chat:aborted',                  onAborted);
        }

        if (window.go && window.go.main && window.go.main.App) {
            window.go.main.App.GenerateCharacterResponse(message, character.id).catch(err => {
                console.error('Failed to request character response:', err);
                textElement.textContent = `[Error: ${err}]`;
                const indicator = bubble.querySelector('.typing-indicator');
                if (indicator) indicator.remove();
                messageDiv.classList.remove('streaming');
                cleanup();
                resolve(null); // treat as abort — stops queue, resets send button
            });
        }
    });
}

// ---------------------------------------------------------------------------
// Generic Backend Stream (no image)
// ---------------------------------------------------------------------------

async function streamAssistantResponseFromBackend(message) {
    const characterContext = buildCharacterContext();
    if (characterContext) message = characterContext + '\n\n' + message;

    const chatHistory = document.querySelector('.chat-messages');

    const messageDiv = document.createElement('div');
    messageDiv.className = 'chat-message assistant streaming';

    const bubble = document.createElement('div');
    bubble.className = 'chat-bubble';
    bubble.innerHTML = `
        <div class="chat-sender">AI<span class="typing-indicator"><span></span><span></span><span></span></span></div>
        <div class="chat-text"></div>
        <div class="chat-think"></div>
        <div class="chat-bubble-buttons">
            <button class="bubble-button icon-button" onclick="regenerateMessage(this)" title="Regenerate">${ICONS.regenerate}</button>
            <button class="bubble-button icon-button" onclick="previousMessage(this)" title="Previous" id="prev-btn">${ICONS.previous}</button>
            <button class="bubble-button icon-button" onclick="nextMessage(this)" title="Next" id="next-btn">${ICONS.next}</button>
            <button class="bubble-button icon-button" onclick="continueMessage(this)" title="Continue">${ICONS.continue}</button>
            <button class="bubble-button icon-button" onclick="editMessage(this)" title="Edit">${ICONS.edit}</button>
            <button class="bubble-button icon-button" onclick="copyMessage(this)" title="Copy">${ICONS.copy}</button>
            <button class="bubble-button icon-button" onclick="undoMessage(this)" title="Undo (remove this and everything after)">${ICONS.undo}</button>
            <button class="bubble-button icon-button" onclick="deleteMessage(this)" title="Delete">${ICONS.delete}</button>
        </div>
    `;

    messageDiv.appendChild(bubble);
    chatHistory.appendChild(messageDiv);
    chatHistory.scrollTop = chatHistory.scrollHeight;

    currentStreamingBubble = bubble;

    const textElement = bubble.querySelector('.chat-text');
    const thinkArea   = bubble.querySelector('.chat-think');

    thinkContent = ''; mainContent = ''; isInThinkMode = false; thinkStartTime = null;

    const onStart = () => { console.log('Stream started'); setSendButtonAbort(); };

    const onChunk = (chunk) => {
        if (!currentStreamingBubble) return;

        if (chunk.startsWith('CHARACTER_INFO:')) {
            const parts = chunk.substring('CHARACTER_INFO:'.length).split('|');
            if (parts.length >= 3) updateBubbleWithCharacterInfo(bubble, parts[0], parts[1], parts[2]);
            return;
        }

        if (THINK_START_REGEX.test(chunk)) { isInThinkMode = true; thinkStartTime = Date.now(); chunk = chunk.replace(THINK_START_REGEX, ''); }
        if (THINK_END_REGEX.test(chunk))   { isInThinkMode = false; chunk = chunk.replace(THINK_END_REGEX, ''); if (thinkContent) updateThinkArea(thinkArea, thinkContent, thinkStartTime); }

        if (isInThinkMode) { thinkContent += chunk; }
        else               { mainContent += chunk; textElement.textContent = mainContent; }
        chatHistory.scrollTop = chatHistory.scrollHeight;
    };

    const onComplete = () => {
        const indicator = bubble.querySelector('.typing-indicator');
        if (indicator) indicator.remove();
        messageDiv.classList.remove('streaming');
        if (thinkContent) updateThinkArea(thinkArea, thinkContent, thinkStartTime);
        setSendButtonSend();
        cleanupStreamListeners();
        currentStreamingBubble = null;
        updateNavigationButtons();
        if (typeof hideGenerateBanner === 'function') hideGenerateBanner();

        if (window.settingsState?.current?.general?.ttsAutoplay && window.AudioService?.playTTS) {
            window.AudioService.playTTS(mainContent, null);
        }
    };

    const onError = (err) => {
        const errorMsg = err || 'Unknown error occurred';
        textElement.textContent += '\n[Error: ' + errorMsg + ']';
        const indicator = bubble.querySelector('.typing-indicator');
        if (indicator) indicator.remove();
        messageDiv.classList.remove('streaming');
        messageDiv.classList.add('error');
        setSendButtonSend();
        currentStreamingBubble = null;
        if (typeof hideGenerateBanner === 'function') hideGenerateBanner();
        
        // Show a persistent notification
        showNotification('Chat generation failed: ' + errorMsg, 'error', 10000);
        
        // Also log to console for debugging
        console.error('Chat stream error:', errorMsg);
    };

    const onAborted = () => {
        console.log('Stream aborted');
        const indicator = bubble.querySelector('.typing-indicator');
        if (indicator) indicator.remove();
        messageDiv.classList.remove('streaming');
        setSendButtonSend();
        cleanupStreamListeners();
        currentStreamingBubble = null;
        if (typeof hideGenerateBanner === 'function') hideGenerateBanner();
    };

    if (window.runtime) {
        window.runtime.EventsOn('chat:start',   onStart);
        window.runtime.EventsOn('chat:chunk',   onChunk);
        window.runtime.EventsOn('chat:complete', onComplete);
        window.runtime.EventsOn('chat:error',   onError);
        window.runtime.EventsOn('chat:aborted', onAborted);
    }

    window.go.main.App.SendMessageStream(message).catch((err) => {
        console.error('Error starting stream:', err);
        onError(err.message || err);
    });
}

// ---------------------------------------------------------------------------
// Backend Stream With Image
// ---------------------------------------------------------------------------

async function streamAssistantResponseFromBackendWithImage(message, imageData) {
    await syncLorebooksToBackend();

    const chatHistory = document.querySelector('.chat-messages');

    const messageDiv = document.createElement('div');
    messageDiv.className = 'chat-message assistant streaming';

    const bubble = document.createElement('div');
    bubble.className = 'chat-bubble';
    bubble.innerHTML = `
        <div class="chat-sender">AI<span class="typing-indicator"><span></span><span></span><span></span></span></div>
        <div class="chat-text"></div>
        <div class="chat-think"></div>
        <div class="chat-bubble-buttons">
            <button class="bubble-button icon-button" onclick="regenerateMessage(this)" title="Regenerate">${ICONS.regenerate}</button>
            <button class="bubble-button icon-button" onclick="previousMessage(this)" title="Previous" id="prev-btn">${ICONS.previous}</button>
            <button class="bubble-button icon-button" onclick="nextMessage(this)" title="Next" id="next-btn">${ICONS.next}</button>
            <button class="bubble-button icon-button" onclick="continueMessage(this)" title="Continue">${ICONS.continue}</button>
            <button class="bubble-button icon-button" onclick="editMessage(this)" title="Edit">${ICONS.edit}</button>
            <button class="bubble-button icon-button" onclick="copyMessage(this)" title="Copy">${ICONS.copy}</button>
            <button class="bubble-button icon-button" onclick="undoMessage(this)" title="Undo (remove this and everything after)">${ICONS.undo}</button>
            <button class="bubble-button icon-button" onclick="deleteMessage(this)" title="Delete">${ICONS.delete}</button>
        </div>
    `;

    messageDiv.appendChild(bubble);
    chatHistory.appendChild(messageDiv);
    chatHistory.scrollTop = chatHistory.scrollHeight;

    currentStreamingBubble = bubble;

    const textElement = bubble.querySelector('.chat-text');
    const thinkArea   = bubble.querySelector('.chat-think');

    thinkContent = ''; mainContent = ''; isInThinkMode = false; thinkStartTime = null;

    const onStart = () => { console.log('Stream started with image'); setSendButtonAbort(); };

    const onChunk = (chunk) => {
        if (!currentStreamingBubble) return;
        if (THINK_START_REGEX.test(chunk)) { isInThinkMode = true; thinkStartTime = Date.now(); chunk = chunk.replace(THINK_START_REGEX, ''); }
        if (THINK_END_REGEX.test(chunk))   { isInThinkMode = false; chunk = chunk.replace(THINK_END_REGEX, ''); if (thinkContent) updateThinkArea(thinkArea, thinkContent, thinkStartTime); }
        if (isInThinkMode) { thinkContent += chunk; }
        else               { mainContent += chunk; textElement.textContent = mainContent; }
        chatHistory.scrollTop = chatHistory.scrollHeight;
    };

    const onComplete = () => {
        const indicator = bubble.querySelector('.typing-indicator');
        if (indicator) indicator.remove();
        messageDiv.classList.remove('streaming');
        if (thinkContent) updateThinkArea(thinkArea, thinkContent, thinkStartTime);
        setSendButtonSend();
        cleanupStreamListeners();
        currentStreamingBubble = null;
        updateNavigationButtons();
        if (typeof hideGenerateBanner === 'function') hideGenerateBanner();

        if (window.settingsState?.current?.general?.ttsAutoplay && window.AudioService?.playTTS) {
            window.AudioService.playTTS(mainContent, null);
        }
    };

    const onError = (err) => {
        textElement.textContent += '\n[Error: ' + err + ']';
        const indicator = bubble.querySelector('.typing-indicator');
        if (indicator) indicator.remove();
        messageDiv.classList.remove('streaming');
        setSendButtonSend();
        currentStreamingBubble = null;
        if (typeof hideGenerateBanner === 'function') hideGenerateBanner();
    };

    const onAborted = () => {
        const indicator = bubble.querySelector('.typing-indicator');
        if (indicator) indicator.remove();
        messageDiv.classList.remove('streaming');
        setSendButtonSend();
        cleanupStreamListeners();
        currentStreamingBubble = null;
        if (typeof hideGenerateBanner === 'function') hideGenerateBanner();
    };

    if (window.runtime) {
        window.runtime.EventsOn('chat:start',    onStart);
        window.runtime.EventsOn('chat:chunk',    onChunk);
        window.runtime.EventsOn('chat:complete', onComplete);
        window.runtime.EventsOn('chat:error',    onError);
        window.runtime.EventsOn('chat:aborted',  onAborted);
    }

    if (window.go && window.go.main && window.go.main.App) {
        window.go.main.App.SendMessageStreamWithImage(message, imageData).catch((err) => {
            console.error('Error starting stream with image:', err);
            onError(err.message || err);
        });
    }
}
