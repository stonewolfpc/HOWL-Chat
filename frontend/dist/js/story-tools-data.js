/**
 * Story Tools — Data Layer
 *
 * Primary storage is MemoryStore (localStorage `howl-chat-memories`).
 * Backend syncs are attempted when Wails runtime is available.
 *
 * Depends on: story-tools-state.js, memory-store.js
 *
 * @module story-tools-data
 * @version 2.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Read
// ---------------------------------------------------------------------------

async function getMemories(filter = {}) {
    // Use MemoryStore as the primary source of truth
    let memories = window.MemoryStore ? window.MemoryStore.getMemoryEntries(filter) : [];

    // If backend is available, attempt to merge any backend-only memories
    if (window.go?.main?.App?.GetMemories) {
        try {
            const backendMemories = await window.go.main.App.GetMemories(
                window.storyToolsState.currentChatId, filter
            );
            if (backendMemories?.length > 0) {
                const localIds = new Set(memories.map(m => m.id));
                for (const bm of backendMemories) {
                    if (!localIds.has(bm.id)) {
                        memories.push(_normalizeBackendMemory(bm));
                    }
                }
            }
        } catch (err) {
            console.warn('Backend memory sync failed, using local:', err);
        }
    }

    return memories;
}

async function getMemoryDetail(memoryId) {
    // Primary: MemoryStore
    if (window.MemoryStore) {
        const local = window.MemoryStore.getMemoryEntry(memoryId);
        if (local) return local;
    }

    // Fallback: backend
    if (window.go?.main?.App?.GetMemoryDetail) {
        try {
            const bm = await window.go.main.App.GetMemoryDetail(
                window.storyToolsState.currentChatId, memoryId
            );
            return bm ? _normalizeBackendMemory(bm) : null;
        } catch (err) {
            console.warn('Backend GetMemoryDetail failed:', err);
        }
    }

    return null;
}

// ---------------------------------------------------------------------------
// Write
// ---------------------------------------------------------------------------

async function saveMemory() {
    const st = window.storyToolsState;
    if (!st.currentMemoryId) { showError('No memory selected to save'); return; }

    const title   = document.getElementById('mem-title').value.trim();
    const summary = document.getElementById('mem-summary').value.trim();
    const tags    = document.getElementById('mem-tags').value.split(',').map(t => t.trim()).filter(Boolean);
    const pinned  = document.getElementById('mem-pinned').checked;

    if (!title) { showError('Memory title is required'); return; }

    try {
        showLoading(true);
        const updates = { title, summary, tags, pinned };

        if (window.MemoryStore) {
            window.MemoryStore.updateMemoryEntry(st.currentMemoryId, updates);
        }

        if (window.go?.main?.App?.EditMemory) {
            await window.go.main.App.EditMemory(st.currentChatId, st.currentMemoryId, updates);
        }

        // Refresh in-memory data
        st.memoriesData = window.MemoryStore ? window.MemoryStore.getMemoryEntries() : st.memoriesData;
        renderFullList(st.memoriesData);
        showSuccess('Memory saved successfully');
        showLoading(false);
    } catch (error) {
        console.error('Save failed:', error);
        showError('Failed to save memory');
        showLoading(false);
    }
}



// ---------------------------------------------------------------------------
// Export
// ---------------------------------------------------------------------------

async function exportMemories() {
    try {
        showLoading(true);
        const filter = {
            scope: document.getElementById('memories-filter-scope').value,
            query: document.getElementById('memories-search').value.trim() || undefined
        };
        const memories = await getMemories(filter);
        const blob = new Blob([JSON.stringify(memories, null, 2)], { type: 'application/json' });
        const url  = URL.createObjectURL(blob);
        const a    = document.createElement('a');
        a.href     = url;
        a.download = `memories_${window.storyToolsState.currentChatId}_${new Date().toISOString().split('T')[0]}.json`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
        showSuccess(`Exported ${memories.length} memories`);
        showLoading(false);
    } catch (error) {
        console.error('Export failed:', error);
        showError('Failed to export memories');
        showLoading(false);
    }
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function _normalizeBackendMemory(bm) {
    // Backend uses snake_case; normalize to camelCase for frontend consistency
    return {
        id:         bm.id          || bm.ID,
        title:      bm.title       || bm.Title,
        summary:    bm.summary     || bm.Summary,
        content:    bm.content     || bm.Content,
        scope:      bm.scope       || bm.Scope       || 'global',
        ownerId:    bm.owner_id    || bm.ownerId     || '',
        ownerKind:  bm.owner_kind  || bm.ownerKind   || '',
        tags:       bm.tags        || bm.Tags        || [],
        pinned:     bm.pinned      || bm.Pinned      || false,
        status:     bm.status      || bm.Status      || 'active',
        confidence: bm.confidence  || bm.Confidence  || 100,
        source:     bm.source      || bm.Source      || 'auto',
        createdAt:  bm.created_at  || bm.createdAt   || new Date().toISOString(),
        updatedAt:  bm.updated_at  || bm.updatedAt   || new Date().toISOString()
    };
}
