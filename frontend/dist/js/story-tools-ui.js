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
        <div class="memory-recent-title">${window.storyToolsState.escapeHtml(memory.title)}</div>
        <div class="memory-recent-date">${dateStr}</div>
    `;
    item.addEventListener('click', () => selectMemory(memory.id));
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
        <div class="memory-item-title">${esc(memory.title)}</div>
        <div class="memory-item-summary">${esc(memory.summary || '')}</div>
        <div class="memory-item-meta">
            <div class="memory-item-tags">
                ${(memory.tags || []).map(tag => `<span class="memory-tag">${esc(tag)}</span>`).join('')}
            </div>
            <div class="memory-item-date">${dateStr}</div>
        </div>
    `;
    item.addEventListener('click', () => selectMemory(memory.id));
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

function showLoading(show) {
    const el = document.getElementById('loading-overlay');
    if (el) el.style.display = show ? 'flex' : 'none';
}

function showSuccess(message) { showNotification(message, 'success'); }
function showError(message)   { showNotification(message, 'error'); }

