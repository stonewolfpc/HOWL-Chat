/**
 * Character State & Shared Utilities
 *
 * Mutable page-level state and helper functions shared across
 * all character sub-modules.
 *
 * @module character-state
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Shared State
// ---------------------------------------------------------------------------

let characterImages     = [];
let thumbnailIndex      = 0;
let availableLorebooks  = [];
let selectedLorebooks   = [];
let currentEditingId    = null;
let currentEditingIndex = -1;
let characterCache      = null;
let dirtyFields         = new Set();

// ---------------------------------------------------------------------------
// Format Utility
// ---------------------------------------------------------------------------

function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k     = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i     = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// ---------------------------------------------------------------------------
// Lorebook Utilities
// ---------------------------------------------------------------------------

function lorebookDisplayName(lorebook) {
    return (lorebook?.title || lorebook?.name || '').trim() || 'Untitled Lorebook';
}

function lorebookKey(lorebook) {
    const id = (lorebook?.id || '').toString().trim();
    if (id) return id;
    return lorebookDisplayName(lorebook).toLowerCase();
}

function resolveLorebooksFromRefs(refs) {
    const all = JSON.parse(localStorage.getItem('howl-chat-lorebooks') || '[]');
    if (!Array.isArray(refs)) return [];
    return refs.map(ref => {
        const refId   = (ref?.id   || '').toString().trim().toLowerCase();
        const refName = lorebookDisplayName(ref).toLowerCase();
        const canonical = all.find(lb => {
            const lbId   = (lb?.id   || '').toString().trim().toLowerCase();
            const lbName = lorebookDisplayName(lb).toLowerCase();
            return (refId && lbId && refId === lbId) || (refName && lbName && refName === lbName);
        });
        return canonical || ref;
    }).filter(Boolean).reduce((acc, lb) => {
        const key = lorebookKey(lb);
        if (!acc.some(e => lorebookKey(e) === key)) acc.push(lb);
        return acc;
    }, []);
}

// ---------------------------------------------------------------------------
// ID Migration
// ---------------------------------------------------------------------------

function markCharacterFieldDirty(id) {
    if (id) dirtyFields.add(id);
}

function isCharacterFieldDirty(id) {
    return dirtyFields.has(id);
}

function resetCharacterDirtyFields() {
    dirtyFields = new Set();
}

function migrateCharactersWithIds(characters) {
    let needsUpdate = false;
    const migrated  = characters.map((c, i) => {
        if (!c.id) {
            needsUpdate = true;
            return { ...c, id: `char_${Date.now()}_${i}_${Math.floor(Math.random() * 1000)}` };
        }
        return c;
    });
    if (needsUpdate) localStorage.setItem('howl-chat-characters', JSON.stringify(migrated));
    return migrated;
}

// ---------------------------------------------------------------------------
// Public Surface
// ---------------------------------------------------------------------------

window.characterState = {
    get characterImages()     { return characterImages; },
    set characterImages(v)    { characterImages = v; },
    get thumbnailIndex()      { return thumbnailIndex; },
    set thumbnailIndex(v)     { thumbnailIndex = v; },
    get availableLorebooks()  { return availableLorebooks; },
    set availableLorebooks(v) { availableLorebooks = v; },
    get selectedLorebooks()   { return selectedLorebooks; },
    set selectedLorebooks(v)  { selectedLorebooks = v; },
    get currentEditingId()    { return currentEditingId; },
    set currentEditingId(v)   { currentEditingId = v; },
    get currentEditingIndex() { return currentEditingIndex; },
    set currentEditingIndex(v){ currentEditingIndex = v; },
    get characterCache()      { return characterCache; },
    set characterCache(v)     { characterCache = v; },
    markFieldDirty: markCharacterFieldDirty,
    isFieldDirty: isCharacterFieldDirty,
    resetDirtyFields: resetCharacterDirtyFields,
    formatFileSize,
    lorebookDisplayName,
    lorebookKey,
    resolveLorebooksFromRefs,
    migrateCharactersWithIds
};
