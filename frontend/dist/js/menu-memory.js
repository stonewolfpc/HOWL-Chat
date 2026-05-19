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
    container.innerHTML = memories.map(m => `
        <div class="memory-list-item" onclick="editMemoryEntry('${m.id}')">
            <div class="memory-item-time">${formatTime(m.createdAt || m.created_at)}</div>
            <div class="memory-item-title">${m.title || 'Untitled'}</div>
        </div>
    `).join('');
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
    container.innerHTML = memories.map(m => `
        <div class="memory-entry">
            <div class="memory-title">${m.title || 'Untitled'}</div>
            <div class="memory-time">${formatTime(m.createdAt || m.created_at)}</div>
        </div>
    `).join('');
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
