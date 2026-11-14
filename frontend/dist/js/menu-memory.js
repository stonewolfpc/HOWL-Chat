/**
 * Menu — Memory Flyout
 *
 * Memory list and edit flyout functions called from the menu bar,
 * plus the formatTime utility.
 *
 * @module menu-memory
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Time Formatter
// ---------------------------------------------------------------------------

function formatTime(timestamp) {
    const diffMs   = Date.now() - new Date(timestamp);
    const diffMins = Math.floor(diffMs / 60000);
    const diffHrs  = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);
    if (diffMins < 1)  return 'Just now';
    if (diffMins < 60) return `${diffMins} minute${diffMins > 1 ? 's' : ''} ago`;
    if (diffHrs  < 24) return `${diffHrs} hour${diffHrs > 1 ? 's' : ''} ago`;
    return `${diffDays} day${diffDays > 1 ? 's' : ''} ago`;
}

// ---------------------------------------------------------------------------
// Memory List Flyout
// ---------------------------------------------------------------------------

async function deleteMemoryMenu(memoryId, memoryTitle) {
    if (!confirm(`Delete memory "${memoryTitle || 'Untitled'}"?\n\nThis cannot be undone.`)) return;
    try {
        const deleted = window.MemoryStore?.deleteMemoryEntry
            ? window.MemoryStore.deleteMemoryEntry(memoryId)
            : false;
        if (!deleted) {
            showNotification('Memory not found', 'error');
            return;
        }
        // Re-render both lists
        loadRecentMemoriesMenu();
        loadMemoryList();
        showNotification('Memory deleted', 'success');
    } catch (error) {
        console.error('Failed to delete memory:', error);
        showNotification('Failed to delete memory', 'error');
    }
}

function toggleMemoryList() {
    const memoryList = document.getElementById('memory-list-flyout');
    if (!memoryList) return;
    memoryList.classList.toggle('show');
    if (memoryList.classList.contains('show')) loadMemoryList();
}

function loadMemoryList() {
    const container = document.getElementById('full-memory-list');
    if (!container) return;
    container.innerHTML = '<div class="memory-loading">Loading memories...</div>';

    const chatId = window.sidebarState?.currentChatData?.sessionId || null;
    let memories = [];
    if (window.MemoryStore) {
        memories = window.MemoryStore.getMemoryEntries();
        // Note: Memories are filtered by scope/ownerId on backend, show all local memories
    }

    if (window.go?.main?.App?.GetMemories && chatId) {
        window.go.main.App.GetMemories(chatId, {}).then(backendMemories => {
            if (backendMemories?.length > 0) {
                const localIds = new Set(memories.map(m => m.id));
                for (const bm of backendMemories) {
                    if (!localIds.has(bm.id)) memories.push(bm);
                }
            }
            renderMemoryList(memories);
        }).catch(() => renderMemoryList(memories));
    } else {
        renderMemoryList(memories);
    }
}

function renderMemoryList(memories) {
    const container = document.getElementById('full-memory-list');
    if (!container) return;
    if (!memories || memories.length === 0) {
        container.innerHTML = '<div class="memory-empty">No memories found</div>';
        return;
    }
    const unresolved = memories.filter(m =>
        (m.scope || m.Scope) === 'character' &&
        (m.ownerKind || m.owner_kind) !== 'character_id'
    ).length;
    container.innerHTML = unresolved > 0
        ? `<div class="memory-ownership-warning" style="padding:10px 12px;margin-bottom:10px;border:1px solid rgba(248,113,113,0.35);border-radius:8px;color:#fca5a5;font-size:12px;">${unresolved} character memor${unresolved === 1 ? 'y has' : 'ies have'} unresolved legacy ownership and ${unresolved === 1 ? 'is' : 'are'} excluded from character prompts and transfer.</div>`
        : '';
    memories.forEach(m => {
        const item = document.createElement('div');
        item.className = 'memory-list-item';
        item.innerHTML = `
            <div class="memory-list-item-content">
                <div class="memory-item-time">${formatTime(m.createdAt || m.created_at)}</div>
                <div class="memory-item-title">${m.title || 'Untitled'}</div>
            </div>
            <span class="memory-delete-btn" title="Delete memory"><span class="icon icon-trash"></span></span>
        `;
        item.querySelector('.memory-list-item-content').addEventListener('click', () => editMemoryEntry(m.id));
        item.querySelector('.memory-delete-btn').addEventListener('click', (e) => {
            e.stopPropagation();
            deleteMemoryMenu(m.id, m.title);
        });
        container.appendChild(item);
    });
}

function loadRecentMemoriesMenu() {
    const container = document.getElementById('recent-memories');
    if (!container) return;
    container.innerHTML = '<div class="memory-loading">Loading memories...</div>';

    const chatId = window.sidebarState?.currentChatData?.sessionId || null;
    let memories = [];
    if (window.MemoryStore) {
        memories = window.MemoryStore.getMemoryEntries();
        // Note: Memories are filtered by scope/ownerId on backend, show all local memories
    }

    if (window.go?.main?.App?.GetMemories && chatId) {
        window.go.main.App.GetMemories(chatId, {}).then(backendMemories => {
            if (backendMemories?.length > 0) {
                const localIds = new Set(memories.map(m => m.id));
                for (const bm of backendMemories) {
                    if (!localIds.has(bm.id)) memories.push(bm);
                }
            }
            renderRecentMemoriesMenu(memories.slice(0, 5));
        }).catch(() => renderRecentMemoriesMenu(memories.slice(0, 5)));
    } else {
        renderRecentMemoriesMenu(memories.slice(0, 5));
    }
}

function renderRecentMemoriesMenu(memories) {
    const container = document.getElementById('recent-memories');
    if (!container) return;
    if (!memories || memories.length === 0) {
        container.innerHTML = '<div class="memory-empty">No memories</div>';
        return;
    }
    container.innerHTML = '';
    memories.forEach(m => {
        const entry = document.createElement('div');
        entry.className = 'memory-entry';
        entry.innerHTML = `
            <div class="memory-entry-content">
                <div class="memory-title">${m.title || 'Untitled'}</div>
                <div class="memory-time">${formatTime(m.createdAt || m.created_at)}</div>
            </div>
            <span class="memory-delete-btn" title="Delete memory"><span class="icon icon-trash"></span></span>
        `;
        entry.querySelector('.memory-entry-content').addEventListener('click', () => editMemoryEntry(m.id));
        entry.querySelector('.memory-delete-btn').addEventListener('click', (e) => {
            e.stopPropagation();
            deleteMemoryMenu(m.id, m.title);
        });
        container.appendChild(entry);
    });
}

// ---------------------------------------------------------------------------
// Memory Edit Flyout
// ---------------------------------------------------------------------------

function editMemoryEntry(memoryId) {
    document.getElementById('memory-list-flyout')?.classList.remove('show');
    loadMemoryForEdit(memoryId);
}

function loadMemoryForEdit(memoryId) {
    const editFlyout   = document.getElementById('memory-edit-flyout');
    const titleField   = document.getElementById('memory-title');
    const contentField = document.getElementById('memory-content');
    if (!editFlyout || !titleField || !contentField) return;

    titleField.setAttribute('data-memory-id', memoryId);

    let memory = null;
    if (window.MemoryStore) {
        memory = window.MemoryStore.getMemoryEntry(memoryId);
    }

    if (memory) {
        titleField.value = memory.title || '';
        contentField.value = memory.summary || memory.content || '';
        editFlyout.classList.add('show');
        return;
    }

    if (window.go?.main?.App?.GetMemoryDetail) {
        const chatId = window.sidebarState?.currentChatData?.sessionId || 'default';
        window.go.main.App.GetMemoryDetail(chatId, memoryId).then(bm => {
            if (bm) { titleField.value = bm.title || ''; contentField.value = bm.summary || bm.content || ''; }
            editFlyout.classList.add('show');
        }).catch(() => showNotification('Failed to load memory', 'error'));
    } else {
        showNotification('Memory not found', 'error');
    }
}

function closeMemoryEdit() {
    document.getElementById('memory-edit-flyout')?.classList.remove('show');
}

function saveMemoryEdit() {
    const titleField   = document.getElementById('memory-title');
    const contentField = document.getElementById('memory-content');
    if (!titleField || !contentField) return;

    const title   = titleField.value.trim();
    const content = contentField.value.trim();

    if (!title || !content) { showNotification('Please fill in both title and content', 'error'); return; }

    const memId = titleField.getAttribute('data-memory-id');
    if (!memId) return;

    if (window.MemoryStore) {
        window.MemoryStore.updateMemoryEntry(memId, { title, summary: content });
    }

    if (window.go?.main?.App?.EditMemory) {
        const chatId = window.sidebarState?.currentChatData?.sessionId || 'default';
        window.go.main.App.EditMemory(chatId, memId, { title, summary: content })
            .then(() => { showNotification('Memory saved successfully', 'success'); closeMemoryEdit(); loadMemoryList(); })
            .catch(() => showNotification('Failed to save memory', 'error'));
    } else {
        showNotification('Memory saved locally', 'success');
        closeMemoryEdit();
        loadMemoryList();
    }
}
