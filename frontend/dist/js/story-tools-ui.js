/**
 * Story Tools — UI Layer
 *
 * All render, editor-populate, preview-update, and notification
 * functions for the memories flyout and editor panel.
 *
 * Depends on: story-tools-state.js, story-tools-data.js
 *
 * @module story-tools-ui
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// List Renderers
// ---------------------------------------------------------------------------

function renderRecentList(memories) {
    const container = document.getElementById('memories-recent');
    if (!container) return;
    if (!memories || memories.length === 0) {
        container.innerHTML = '<div class="memories-empty">No memories loaded</div>';
        return;
    }
    container.innerHTML = '';
    memories.forEach(m => container.appendChild(createRecentItem(m)));
}

function createRecentItem(memory) {
    const item = document.createElement('div');
    item.className      = 'memory-recent-item';
    item.dataset.memoryId = memory.id;
    const dateStr       = new Date(memory.lastTriggeredAt || memory.createdAt).toLocaleDateString();
    item.innerHTML      = `
        <div class="memory-recent-info">
            <div class="memory-recent-title">${window.storyToolsState.escapeHtml(memory.title)}</div>
            <div class="memory-recent-date">${dateStr}</div>
        </div>
        <span class="memory-delete-btn" title="Delete memory"><span class="icon icon-trash"></span></span>
    `;
    item.querySelector('.memory-recent-info').addEventListener('click', () => selectMemory(memory.id));
    item.querySelector('.memory-delete-btn').addEventListener('click', (e) => {
        e.stopPropagation();
        deleteMemory(memory.id, memory.title);
    });
    return item;
}

function renderFullList(memories) {
    const container = document.getElementById('memories-list');
    if (!container) return;
    if (!memories || memories.length === 0) {
        container.innerHTML = '<div class="memories-empty">No memories found</div>';
        return;
    }
    container.innerHTML = '';
    memories.forEach(m => container.appendChild(createFullListItem(m)));
}

function createFullListItem(memory) {
    const item = document.createElement('div');
    item.className        = 'memory-item';
    item.dataset.memoryId = memory.id;
    const dateStr         = new Date(memory.lastTriggeredAt || memory.createdAt).toLocaleDateString();
    const esc             = window.storyToolsState.escapeHtml;
    item.innerHTML = `
        <div class="memory-item-content">
            <div class="memory-item-title">${esc(memory.title)}</div>
            <div class="memory-item-summary">${esc(memory.summary || '')}</div>
            <div class="memory-item-meta">
                <div class="memory-item-tags">
                    ${(memory.tags || []).map(tag => `<span class="memory-tag">${esc(tag)}</span>`).join('')}
                </div>
                <div class="memory-item-date">${dateStr}</div>
            </div>
        </div>
        <span class="memory-delete-btn" title="Delete memory"><span class="icon icon-trash"></span></span>
    `;
    item.querySelector('.memory-item-content').addEventListener('click', () => selectMemory(memory.id));
    item.querySelector('.memory-delete-btn').addEventListener('click', (e) => {
        e.stopPropagation();
        deleteMemory(memory.id, memory.title);
    });
    return item;
}

// ---------------------------------------------------------------------------
// Editor
// ---------------------------------------------------------------------------

function populateEditor(memory) {
    window.storyToolsState.currentMemoryId = memory.id;
    document.getElementById('mem-title').value   = memory.title   || '';
    document.getElementById('mem-summary').value = memory.summary || '';
    document.getElementById('mem-tags').value    = (memory.tags || []).join(', ');
    document.getElementById('mem-pinned').checked = memory.pinned || false;
    updateTriggerInfo(memory);
    updatePreview(memory);
}

function clearEditor() {
    window.storyToolsState.currentMemoryId = null;
    document.getElementById('mem-title').value    = '';
    document.getElementById('mem-summary').value  = '';
    document.getElementById('mem-tags').value     = '';
    document.getElementById('mem-pinned').checked = false;
    updateTriggerInfo({});
    updatePreview({});
}

function updateTriggerInfo(memory) {
    const triggerInfo = document.getElementById('mem-trigger-info');
    const triggerType = document.querySelector('.trigger-type');
    if (!triggerType) return;
    if (memory.author === 'auto') {
        triggerType.textContent = 'System trigger - not editable';
        if (triggerInfo) triggerInfo.className = 'trigger-info';
    } else {
        triggerType.textContent = 'Manual entry - editable';
        if (triggerInfo) triggerInfo.className = 'trigger-info manual-trigger';
    }
}

function updatePreview(memory) {
    const preview = document.getElementById('memories-preview');
    if (!preview) return;
    if (!memory.pinned) {
        preview.innerHTML = '<div class="preview-empty">Select a memory to view details</div>';
        return;
    }
    const esc = window.storyToolsState.escapeHtml;
    preview.innerHTML = `
        <div class="preview-title">${esc(memory.title)}</div>
        <div class="preview-summary">${esc(memory.summary || '')}</div>
        <div class="preview-tags">
            ${(memory.tags || []).map(tag => `<span class="memory-tag">${esc(tag)}</span>`).join('')}
        </div>
        <div class="preview-pinned">📌 Pinned to context</div>
    `;
}

// ---------------------------------------------------------------------------
// Select Memory
// ---------------------------------------------------------------------------

async function selectMemory(memoryId) {
    try {
        showLoading(true);
        const memory = await getMemoryDetail(memoryId);
        populateEditor(memory);
        openFlyout();
        showLoading(false);
    } catch (error) {
        console.error('Failed to select memory:', error);
        showError('Failed to load memory details');
        showLoading(false);
    }
}

function openFlyout() {
    document.getElementById('memories-flyout')?.classList.add('show');
}

// ---------------------------------------------------------------------------
// Status / Notifications
// ---------------------------------------------------------------------------

// One implementation for two call styles that used to be two competing
// globals: rag-knowledge.js called showLoading('Uploading document...') while
// chat.js calls showLoading(true). Whichever file loaded last decided which
// meaning `true` had, so on story-tools the overlay could be shown with the
// text "true" and never hidden again. Both address #loading-overlay.
function showLoading(show) {
    const el = document.getElementById('loading-overlay');
    if (!el) return;
    if (typeof show === 'string') {
        const label = el.querySelector('.loading-text');
        if (label) label.textContent = show;
        el.style.display = 'flex';
        return;
    }
    el.style.display = show ? 'flex' : 'none';
}

async function deleteMemory(memoryId, memoryTitle) {
    const st = window.storyToolsState;
    // Two call styles, one implementation: the editor's Delete button passes
    // nothing (it acts on the selected memory), list rows pass id and title.
    const id = memoryId || st?.currentMemoryId;
    if (!id) { showError('No memory selected to delete'); return; }

    const known = (st?.memoriesData || []).find(m => m.id === id);
    const title = memoryTitle || known?.title;
    if (!confirm(`Delete memory "${title || 'Untitled'}"? This cannot be undone.`)) return;

    try {
        const deleted = window.MemoryStore?.deleteMemoryEntry
            ? window.MemoryStore.deleteMemoryEntry(id)
            : false;

        // The backend copy stays authoritative for the chat; a failure here is
        // logged, not fatal (the local store already dropped it).
        if (st?.currentChatId && window.go?.main?.App?.DeleteMemory) {
            try { await window.go.main.App.DeleteMemory(st.currentChatId, id); }
            catch (err) { console.warn('Backend memory delete failed:', err); }
        }

        if (!deleted) { showError('Memory not found'); return; }

        renderRecentList(await getMemories({ limit: 5, sort: 'newest' }));
        renderFullList(await getMemories({}));
        if (st?.currentMemoryId === id) clearEditor();
        showSuccess('Memory deleted');
    } catch (error) {
        console.error('Failed to delete memory:', error);
        showError('Failed to delete memory');
    }
}

function showSuccess(message) { showNotification(message, 'success'); }
function showError(message)   { showNotification(message, 'error'); }

