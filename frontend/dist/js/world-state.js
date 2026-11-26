/**
 * World State & Lorebook Utilities
 *
 * Owns the mutable world page state arrays and provides
 * pure utility functions for lorebook key/name resolution
 * shared across all world sub-modules.
 *
 * @module world-state
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Shared State
// ---------------------------------------------------------------------------

let worlds              = [];
let selectedScenarios   = [];
let selectedCharacters  = [];
let selectedLorebooks   = [];
let editingWorldId      = null;
let worldImages         = [];
let worldThumbnailIndex = 0;

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
    const allLorebooks = (window.LorebookRepo ? window.LorebookRepo.all() : []);
    if (!Array.isArray(refs)) return [];

    const resolved = refs.map(ref => {
        const refId   = (ref?.id || '').toString().trim().toLowerCase();
        const refName = lorebookDisplayName(ref).toLowerCase();
        return allLorebooks.find(lb => {
            const lbId   = (lb?.id || '').toString().trim().toLowerCase();
            const lbName = lorebookDisplayName(lb).toLowerCase();
            return (refId && lbId && refId === lbId) || (refName && lbName && refName === lbName);
        }) || ref;
    });

    return resolved.filter(Boolean).reduce((acc, lorebook) => {
        const key = lorebookKey(lorebook);
        if (!acc.some(existing => lorebookKey(existing) === key)) acc.push(lorebook);
        return acc;
    }, []);
}

// ---------------------------------------------------------------------------
// Exports
// ---------------------------------------------------------------------------

window.worldState = {
    get worlds()              { return worlds; },
    set worlds(v)             { worlds = v; },
    get selectedScenarios()   { return selectedScenarios; },
    set selectedScenarios(v)  { selectedScenarios = v; },
    get selectedCharacters()  { return selectedCharacters; },
    set selectedCharacters(v) { selectedCharacters = v; },
    get selectedLorebooks()   { return selectedLorebooks; },
    set selectedLorebooks(v)  { selectedLorebooks = v; },
    get editingWorldId()      { return editingWorldId; },
    set editingWorldId(v)     { editingWorldId = v; },
    get worldImages()         { return worldImages; },
    set worldImages(v)        { worldImages = v; },
    get worldThumbnailIndex() { return worldThumbnailIndex; },
    set worldThumbnailIndex(v){ worldThumbnailIndex = v; },
    lorebookDisplayName,
    lorebookKey,
    resolveLorebooksFromRefs
};
