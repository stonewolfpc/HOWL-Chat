/**
 * Story Tools — State & Utilities
 *
 * Shared mutable state and pure utility functions used across
 * all story-tools sub-modules.
 *
 * @module story-tools-state
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------

let currentChatId  = null;
let currentMemoryId = null;
let memoriesData   = [];
let currentTab     = 'all';

// ---------------------------------------------------------------------------
// Chat ID Resolution
// ---------------------------------------------------------------------------

function getCurrentChatId() {
    const urlParams = new URLSearchParams(window.location.search);
    const chatId    = urlParams.get('chatId');
    if (chatId) return chatId;

    const stored = localStorage.getItem('howl-chat-current-chat');
    if (stored) {
        try { return JSON.parse(stored).chatName || 'default'; }
        catch { /* fall through */ }
    }
    return 'default';
}

// ---------------------------------------------------------------------------
// DOM Utilities
// ---------------------------------------------------------------------------

function debounce(func, wait) {
    let timeout;
    return function (...args) {
        clearTimeout(timeout);
        timeout = setTimeout(() => { clearTimeout(timeout); func(...args); }, wait);
    };
}

// ---------------------------------------------------------------------------
// Public Surface
// ---------------------------------------------------------------------------

window.storyToolsState = {
    get currentChatId()    { return currentChatId; },
    set currentChatId(v)   { currentChatId = v; },
    get currentMemoryId()  { return currentMemoryId; },
    set currentMemoryId(v) { currentMemoryId = v; },
    get memoriesData()     { return memoriesData; },
    set memoriesData(v)    { memoriesData = v; },
    get currentTab()       { return currentTab; },
    set currentTab(v)      { currentTab = v; },
    getCurrentChatId,
    escapeHtml,
    debounce
};
