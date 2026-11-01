/**
 * Filter Engine
 *
 * Reusable search, sort, and attachment-filter utilities for all
 * entity lists (worlds, scenarios, characters, lorebooks, chats).
 *
 * Zero backend dependencies — operates on in-memory arrays.
 *
 * @module filter-engine
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Utilities
// ---------------------------------------------------------------------------

function _normalizeText(value) {
    if (value === null || value === undefined) return '';
    return String(value).toLowerCase().trim();
}

function _getField(item, field) {
    if (!item) return '';
    const val = item[field];
    if (val !== undefined && val !== null) return val;
    // Fallback aliases
    if (field === 'name')  return item.title || item.label || '';
    if (field === 'title') return item.name || item.label || '';
    if (field === 'date')  return item.createdAt || item.updatedAt || item.timestamp || '';
    return '';
}

function _parseDate(value) {
    if (!value) return 0;
    if (typeof value === 'number') return value;
    const ts = Date.parse(value);
    return isNaN(ts) ? 0 : ts;
}

function _extractIdSet(entity, relationField) {
    if (!entity || !relationField) return new Set();
    const refs = entity[relationField] || entity.howl?.[relationField];
    if (!Array.isArray(refs)) return new Set();
    return new Set(refs.map(r => (typeof r === 'string' ? r : r?.id)).filter(Boolean));
}

// ---------------------------------------------------------------------------
// Search
// ---------------------------------------------------------------------------

/**
 * Filter an array of items by a text query across given fields.
 *
 * @param {Array}  items  Array of objects to filter.
 * @param {string} query  Search text (case-insensitive).
 * @param {string[]} fields Property names to search. Defaults to ['name','title'].
 * @returns {Array} Filtered items.
 */
function filterItems(items, query, fields) {
    if (!Array.isArray(items)) return [];
    if (!query || !_normalizeText(query)) return items.slice();

    const searchFields = Array.isArray(fields) && fields.length > 0
        ? fields
        : ['name', 'title'];

    const terms = _normalizeText(query).split(/\s+/).filter(Boolean);

    return items.filter(item => {
        const text = searchFields.map(f => _normalizeText(_getField(item, f))).join(' ');
        return terms.every(term => text.includes(term));
    });
}

// ---------------------------------------------------------------------------
// Sort
// ---------------------------------------------------------------------------

/**
 * Sort an array of items by a field.
 *
 * @param {Array}  items     Array to sort (mutates copy).
 * @param {string} field     'name' | 'title' | 'date' | 'count' | 'messageCount' | 'characterCount'
 * @param {string} direction 'asc' | 'desc'  (default 'asc')
 * @returns {Array} New sorted array.
 */
function sortItems(items, field, direction) {
    if (!Array.isArray(items)) return [];
    const dir = direction === 'desc' ? -1 : 1;
    const sorted = items.slice();

    sorted.sort((a, b) => {
        let va, vb;

        if (field === 'date' || field === 'createdAt' || field === 'updatedAt') {
            va = _parseDate(_getField(a, field));
            vb = _parseDate(_getField(b, field));
        } else if (field === 'count' || field === 'messageCount') {
            va = a.metadata?.messageCount ?? a.messageCount ?? 0;
            vb = b.metadata?.messageCount ?? b.messageCount ?? 0;
        } else if (field === 'characterCount') {
            va = a.metadata?.characterCount ?? a.characterCount ?? 0;
            vb = b.metadata?.characterCount ?? b.characterCount ?? 0;
        } else {
            va = _normalizeText(_getField(a, field));
            vb = _normalizeText(_getField(b, field));
        }

        if (va < vb) return -1 * dir;
        if (va > vb) return 1 * dir;
        return 0;
    });

    return sorted;
}

// ---------------------------------------------------------------------------
// Attachment Filter
// ---------------------------------------------------------------------------

/**
 * Filter items to only those referenced by a parent entity.
 *
 * Example: filterByAttachment(scenarios, currentWorld, 'scenarios')
 * keeps only scenarios whose id is in currentWorld.scenarios.
 *
 * @param {Array}  items          Items to filter.
 * @param {Object} parentEntity   The entity that holds references.
 * @param {string} relationField  Property name on parent that holds refs.
 * @returns {Array} Attached items.
 */
function filterByAttachment(items, parentEntity, relationField) {
    if (!Array.isArray(items) || !parentEntity || !relationField) return items.slice();
    const idSet = _extractIdSet(parentEntity, relationField);
    if (idSet.size === 0) return [];
    return items.filter(item => idSet.has(item.id));
}

/**
 * Inverse: filter items to those NOT referenced by the parent.
 */
function filterByDetachment(items, parentEntity, relationField) {
    if (!Array.isArray(items) || !parentEntity || !relationField) return items.slice();
    const idSet = _extractIdSet(parentEntity, relationField);
    if (idSet.size === 0) return items.slice();
    return items.filter(item => !idSet.has(item.id));
}

// ---------------------------------------------------------------------------
// Combined Pipeline
// ---------------------------------------------------------------------------

/**
 * Run the full filter → sort pipeline in one call.
 *
 * @param {Array}  items   Raw items.
 * @param {Object} opts
 *   @param {string}   opts.query          Search text.
 *   @param {string[]} opts.searchFields   Fields to search.
 *   @param {string}   opts.sortField      Field to sort by.
 *   @param {string}   opts.sortDirection  'asc' | 'desc'.
 *   @param {Object}   opts.attachTo       Parent entity for attachment filter.
 *   @param {string}   opts.attachField    Relation field on parent.
 *   @param {boolean}  opts.attachInvert   If true, show detached instead of attached.
 * @returns {Array} Processed items.
 */
function processItems(items, opts) {
    if (!Array.isArray(items)) return [];
    opts = opts || {};

    let result = items.slice();

    // Attachment filter
    if (opts.attachTo && opts.attachField) {
        result = opts.attachInvert
            ? filterByDetachment(result, opts.attachTo, opts.attachField)
            : filterByAttachment(result, opts.attachTo, opts.attachField);
    }

    // Text search
    if (opts.query) {
        result = filterItems(result, opts.query, opts.searchFields);
    }

    // Sort
    if (opts.sortField) {
        result = sortItems(result, opts.sortField, opts.sortDirection);
    }

    return result;
}

// ---------------------------------------------------------------------------
// Exports
// ---------------------------------------------------------------------------

window.FilterEngine = {
    filterItems,
    sortItems,
    filterByAttachment,
    filterByDetachment,
    processItems
};
