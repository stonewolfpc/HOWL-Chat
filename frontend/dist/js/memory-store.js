/**
 * Memory Store — Frontend storage for automated memories
 *
 * Separated from manual lorebook entries. Uses localStorage key
 * `howl-chat-memories` as the primary source, syncs to backend
 * when available via AddMemory / GetMemories bindings.
 *
 * @module memory-store
 * @version 1.0.0
 */

'use strict';

const MEMORY_STORAGE_KEY = 'howl-chat-memories';

// ---------------------------------------------------------------------------
// CRUD Operations
// ---------------------------------------------------------------------------

function _loadMemories() {
    try {
        const raw = localStorage.getItem(MEMORY_STORAGE_KEY);
        const parsed = raw ? JSON.parse(raw) : [];
        return Array.isArray(parsed) ? parsed : [];
    } catch (_) {
        return [];
    }
}

function _saveMemories(memories) {
    localStorage.setItem(MEMORY_STORAGE_KEY, JSON.stringify(memories));
}

function _generateMemoryId() {
    return `mem-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`;
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

/**
 * Create a new automated memory entry.
 *
 * @param {Object} data
 * @param {string} data.title        — Memory title
 * @param {string} data.summary      — Short summary
 * @param {string} [data.content]    — Full content
 * @param {string} [data.scope]      — 'global' | 'world' | 'scenario' | 'character'
 * @param {string} [data.ownerId]    — Character/world/scenario ID for scoped memories
 * @param {string[]} [data.tags]     — Categorization tags
 * @param {string} [data.source]     — 'auto' | 'manual' | 'system'
 * @param {string} [data.trigger]    — Trigger description (e.g. 'death', 'anniversary')
 * @param {number} [data.confidence] — 0-100
 * @returns {Object} The created memory record
 */
function createMemoryEntry(data) {
    const now = new Date().toISOString();
    const entry = {
        id:          data.id || _generateMemoryId(),
        title:       data.title,
        summary:     data.summary || '',
        content:     data.content || data.summary || '',
        scope:       data.scope || 'global',
        ownerId:     data.ownerId || '',
        ownerKind:   data.ownerKind || '',
        tags:        data.tags || [],
        pinned:      data.pinned || false,
        status:      'active',
        confidence:  data.confidence || 100,
        source:      data.source || 'auto',
        trigger:     data.trigger || null,
        createdAt:   now,
        updatedAt:   now
    };

    const memories = _loadMemories();
    memories.push(entry);
    _saveMemories(memories);

    // Sync to backend if available
    if (window.go?.main?.App?.AddMemory) {
        const chatId = window.sidebarState?.currentChatData?.sessionId || 'default';
        const backendEntry = {
            id:              entry.id,
            title:           entry.title,
            summary:         entry.summary,
            rawText:         entry.content,
            scope:           entry.scope,
            scopeId:         entry.ownerId,
            ownerKind:       entry.ownerKind,
            tags:            entry.tags,
            pinned:          entry.pinned,
            status:          entry.status,
            confidence:      entry.confidence,
            trigger:         entry.source,
            createdAt:       entry.createdAt,
            lastTriggeredAt: entry.updatedAt,
            type:            'automated',
            author:          'system'
        };
        window.go.main.App.AddMemory(chatId, backendEntry).catch(err => {
            console.warn('Failed to sync memory to backend:', err);
        });
    }

    console.log(`Created memory entry: ${entry.title} (${entry.id})`);
    return entry;
}

/**
 * Get all memories, optionally filtered.
 *
 * @param {Object} [filter]
 * @param {string} [filter.scope]  — Filter by scope
 * @param {string} [filter.query]  — Search in title/summary/tags
 * @param {string} [filter.status] — Filter by status
 * @returns {Object[]} Filtered memory list
 */
function getMemoryEntries(filter = {}) {
    let memories = _loadMemories();

    if (filter.scope && filter.scope !== 'all') {
        memories = memories.filter(m => m.scope === filter.scope);
    }
    if (filter.status) {
        memories = memories.filter(m => m.status === filter.status);
    }
    if (filter.query) {
        const q = filter.query.toLowerCase();
        memories = memories.filter(m =>
            (m.title || '').toLowerCase().includes(q) ||
            (m.summary || '').toLowerCase().includes(q) ||
            (m.tags || []).some(t => t.toLowerCase().includes(q))
        );
    }

    // Sort newest first
    memories.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));
    return memories;
}

/**
 * Get a single memory by ID.
 *
 * @param {string} id
 * @returns {Object|null}
 */
function getMemoryEntry(id) {
    return _loadMemories().find(m => m.id === id) || null;
}

/**
 * Update an existing memory.
 *
 * @param {string} id
 * @param {Object} updates
 * @returns {Object|null} Updated memory or null if not found
 */
function updateMemoryEntry(id, updates) {
    const memories = _loadMemories();
    const idx = memories.findIndex(m => m.id === id);
    if (idx === -1) return null;

    memories[idx] = { ...memories[idx], ...updates, updatedAt: new Date().toISOString() };
    _saveMemories(memories);

    if (window.go?.main?.App?.EditMemory) {
        const chatId = window.sidebarState?.currentChatData?.sessionId || 'default';
        window.go.main.App.EditMemory(chatId, id, updates).catch(() => {});
    }

    return memories[idx];
}

/**
 * Delete a memory by ID.
 *
 * @param {string} id
 * @returns {boolean}
 */
function deleteMemoryEntry(id) {
    const memories = _loadMemories();
    const filtered = memories.filter(m => m.id !== id);
    if (filtered.length === memories.length) return false;

    _saveMemories(filtered);

    if (window.go?.main?.App?.DeleteMemory) {
        const chatId = window.sidebarState?.currentChatData?.sessionId || 'default';
        window.go.main.App.DeleteMemory(chatId, id).catch(() => {});
    }

    return true;
}

/**
 * Get memories scoped to a specific character.
 *
 * @param {string} characterId
 * @returns {Object[]}
 */
function getCharacterMemories(characterId) {
    // Only stable-ID-owned records are character-private. Legacy name-owned or
    // unmarked records stay in storage but fail closed here and in the backend.
    return getMemoryEntries({ scope: 'character' }).filter(m =>
        m.ownerKind === 'character_id' && m.ownerId === characterId
    );
}

/**
 * Get memories scoped to the current scenario.
 *
 * @returns {Object[]}
 */
function getScenarioMemories() {
    return getMemoryEntries({ scope: 'scenario' });
}

/**
 * Get global (unscoped) memories.
 *
 * @returns {Object[]}
 */
function getGlobalMemories() {
    return getMemoryEntries({ scope: 'global' });
}

// ---------------------------------------------------------------------------
// Global Exports
// ---------------------------------------------------------------------------

window.MemoryStore = {
    createMemoryEntry,
    getMemoryEntries,
    getMemoryEntry,
    updateMemoryEntry,
    deleteMemoryEntry,
    getCharacterMemories,
    getScenarioMemories,
    getGlobalMemories
};
