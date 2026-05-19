/**
 * Sidebar State & Shared Utilities
 *
 * Owns the single source of truth for current chat context data
 * and provides pure utility functions used across all sidebar modules.
 *
 * @module sidebar-state
 * @version 1.0.0
 */

// ---------------------------------------------------------------------------
// Global Chat Context State
// ---------------------------------------------------------------------------

if (!window._sidebarCurrentChatData) {
    window._sidebarCurrentChatData = {
        characters: [],
        world: null,
        scenario: null,
        lorebooks: [],
        chatName: 'New Chat'
    };
}
var currentChatData = window._sidebarCurrentChatData;

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

const CharacterStatus = {
    ONLINE:  'online',
    READY:   'ready',
    IDLE:    'idle',
    OFFLINE: 'offline'
};

// ---------------------------------------------------------------------------
// State Helpers
// ---------------------------------------------------------------------------

function resetToDefaultState() {
    currentChatData = {
        characters: [],
        world: null,
        scenario: null,
        lorebooks: [],
        chatName: 'New Chat'
    };
    window._sidebarCurrentChatData = currentChatData;
}

// ---------------------------------------------------------------------------
// Chat Contract Validator (development / debug tool)
// ---------------------------------------------------------------------------

function logChatContract(stage, payload) {
    try {
        const data    = payload || {};
        const history = data.History    || data.history    || [];
        const world   = data.World      || data.world      || null;
        const scenario= data.Scenario   || data.scenario   || null;
        const characters = data.Characters || data.characters || [];
        const lorebooks  = data.Lorebooks  || data.lorebooks  || [];
        const meta    = data.Meta       || data.meta       || {};
        const issues  = [];

        if (!Array.isArray(history))    issues.push('history:not-array');
        if (!Array.isArray(characters)) issues.push('characters:not-array');
        if (!Array.isArray(lorebooks))  issues.push('lorebooks:not-array');
        if (world    && typeof world    !== 'object') issues.push('world:not-object');
        if (scenario && typeof scenario !== 'object') issues.push('scenario:not-object');
        if (meta     && typeof meta     !== 'object') issues.push('meta:not-object');

        if (world    && world.image    && typeof world.image    !== 'string' && typeof world.image    !== 'object') issues.push('world.image:invalid-type');
        if (scenario && scenario.image && typeof scenario.image !== 'string' && typeof scenario.image !== 'object') issues.push('scenario.image:invalid-type');

        const summary = {
            stage,
            historyCount:   Array.isArray(history)    ? history.length    : -1,
            hasWorld:       !!world,
            hasScenario:    !!scenario,
            characterCount: Array.isArray(characters) ? characters.length : -1,
            lorebookCount:  Array.isArray(lorebooks)  ? lorebooks.length  : -1,
            chatName: meta.Name || meta.name || data.chatName || 'New Chat',
            issues
        };

        if (issues.length > 0) {
            console.warn('Chat contract warning:', summary);
        } else {
            console.log('Chat contract OK:', summary);
        }
    } catch (e) {
        console.warn('Chat contract check failed:', e);
    }
}

// ---------------------------------------------------------------------------
// Image Resolution Utility
// ---------------------------------------------------------------------------

/**
 * Accepts a plain URL string or an image object (path/url/data/thumbnail)
 * and returns a usable src string.
 * @param {string|Object} imageValue
 * @returns {string}
 */
function resolveImageSrc(imageValue) {
    if (!imageValue) return '';
    if (typeof imageValue === 'string') return imageValue;
    if (typeof imageValue === 'object') {
        if (typeof imageValue.path      === 'string') return imageValue.path;
        if (typeof imageValue.url       === 'string') return imageValue.url;
        if (typeof imageValue.data      === 'string') return imageValue.data;
        if (typeof imageValue.thumbnail === 'string') return imageValue.thumbnail;
        // Handle {id, stored:true} references — can't resolve async here, but
        // consumer should use HOWLStorage.getImage() for these
        if (imageValue.id && imageValue.stored) return '';
    }
    return '';
}

// ---------------------------------------------------------------------------
// Reference Matching Utilities
// ---------------------------------------------------------------------------

/**
 * Returns true when an item matches a ref by id or name/title.
 */
function refMatches(item, ref) {
    if (!item || !ref) return false;
    const itemId  = (item.id   || '').toString().trim().toLowerCase();
    const refId   = (ref.id    || '').toString().trim().toLowerCase();
    if (itemId && refId && itemId === refId) return true;

    const itemName = (item.name || item.title || '').toString().trim().toLowerCase();
    const refName  = (ref.name  || ref.title  || '').toString().trim().toLowerCase();
    return !!(itemName && refName && itemName === refName);
}

/**
 * Filters allItems to those that match any entry in refs.
 */
function resolveLinkedItems(allItems, refs) {
    if (!Array.isArray(allItems) || !Array.isArray(refs) || refs.length === 0) return [];
    return allItems.filter(item => refs.some(ref => refMatches(item, ref)));
}

// ---------------------------------------------------------------------------
// Global Exports
// ---------------------------------------------------------------------------

window.logChatContract    = logChatContract;
window.sidebarState = {
    get currentChatData()      { return currentChatData; },
    set currentChatData(value) { currentChatData = value; },
    CharacterStatus,
    resetToDefaultState,
    logChatContract,
    resolveImageSrc,
    refMatches,
    resolveLinkedItems
};
