/**
 * Scenario State & Lorebook Utilities
 *
 * Owns mutable scenario page state and all lorebook key/name/
 * resolution helpers shared across the scenario sub-modules.
 *
 * @module scenario-state
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Shared State
// ---------------------------------------------------------------------------

let scenarioImages         = [];
let scenarioThumbnailIndex = 0;
let editingScenarioIndex   = -1;
let availableCharacters    = [];
let selectedCharacters     = [];
let selectedLorebooks      = [];

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

function normalizeStoredLorebook(raw) {
    if (!raw || typeof raw !== 'object') return null;
    const title = (raw.title || raw.name || '').trim();
    return {
        ...raw,
        id:    (raw.id || '').toString().trim() || lorebookKey(raw),
        title,
        name:  (raw.name || raw.title || title).trim()
    };
}

function refDisplayName(ref) {
    return (ref?.title || ref?.name || '').trim();
}

function mergeLorebookWithRef(canonical, ref) {
    const refTitle = refDisplayName(ref);
    if (!canonical) {
        const title = refTitle || lorebookDisplayName(ref);
        return normalizeStoredLorebook({ ...ref, id: (ref?.id || '').toString().trim() || lorebookKey(ref), title, name: title });
    }
    const title = (canonical.title || canonical.name || refTitle || '').trim();
    return normalizeStoredLorebook({
        ...canonical,
        id:    (canonical.id || ref?.id || '').toString().trim(),
        title: title || lorebookDisplayName(canonical),
        name:  title || lorebookDisplayName(canonical)
    });
}

function getScenarioLorebookRefs(scenario) {
    const howlRefs = scenario?.howl?.lorebooks;
    const topRefs  = scenario?.lorebooks;
    if (Array.isArray(howlRefs) && howlRefs.length > 0) return howlRefs;
    if (Array.isArray(topRefs)  && topRefs.length  > 0) return topRefs;
    return [];
}

function lorebookToRef(lorebook) {
    const all = (window.LorebookRepo ? window.LorebookRepo.all() : [])
        .map(normalizeStoredLorebook).filter(Boolean);
    const refId    = (lorebook?.id || '').toString().trim().toLowerCase();
    const canonical = refId ? all.find(lb => (lb.id || '').toString().trim().toLowerCase() === refId) : null;
    const source   = canonical || normalizeStoredLorebook(lorebook) || lorebook;
    const title    = (source.title || source.name || refDisplayName(lorebook) || '').trim();
    return {
        id:          (source.id || lorebook?.id || lorebookKey(source)).toString().trim(),
        title:       title || lorebookDisplayName(source),
        name:        title || lorebookDisplayName(source),
        description: source.description || source.content || lorebook?.description || ''
    };
}

function resolveLorebooksFromRefs(refs) {
    const all = (window.LorebookRepo ? window.LorebookRepo.all() : [])
        .map(normalizeStoredLorebook).filter(Boolean);
    if (!Array.isArray(refs)) return [];
    return refs.map(rawRef => {
        const ref = typeof rawRef === 'string' ? { id: rawRef }
                  : (rawRef && typeof rawRef === 'object' ? rawRef : null);
        if (!ref) return null;
        const refId   = (ref.id   || '').toString().trim().toLowerCase();
        const refName = refDisplayName(ref).toLowerCase();
        let canonical = refId ? all.find(lb => (lb.id || '').toString().trim().toLowerCase() === refId) : null;
        if (!canonical && refName) {
            canonical = all.find(lb => (lb.title || lb.name || '').trim().toLowerCase() === refName);
        }
        return mergeLorebookWithRef(canonical, ref);
    }).filter(Boolean).reduce((acc, lb) => {
        const key = lorebookKey(lb);
        if (!acc.some(e => lorebookKey(e) === key)) acc.push(lb);
        return acc;
    }, []);
}

// ---------------------------------------------------------------------------
// Public Surface
// ---------------------------------------------------------------------------

window.scenarioState = {
    get scenarioImages()         { return scenarioImages; },
    set scenarioImages(v)        { scenarioImages = v; },
    get scenarioThumbnailIndex() { return scenarioThumbnailIndex; },
    set scenarioThumbnailIndex(v){ scenarioThumbnailIndex = v; },
    get editingScenarioIndex()   { return editingScenarioIndex; },
    set editingScenarioIndex(v) { editingScenarioIndex = v; },
    get availableCharacters()  { return availableCharacters; },
    set availableCharacters(v) { availableCharacters = v; },
    get selectedCharacters()   { return selectedCharacters; },
    set selectedCharacters(v)  { selectedCharacters = v; },
    get selectedLorebooks()    { return selectedLorebooks; },
    set selectedLorebooks(v)   { selectedLorebooks = v; },
    lorebookDisplayName,
    lorebookKey,
    normalizeStoredLorebook,
    refDisplayName,
    mergeLorebookWithRef,
    getScenarioLorebookRefs,
    lorebookToRef,
    resolveLorebooksFromRefs
};
