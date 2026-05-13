/**
 * Story Tools Module
 * 
 * This module provides comprehensive memory management functionality
 * for viewing, editing, and managing automated lorebook entries
 * captured during chat sessions. Features include search, filtering,
 * export/import, and real-time updates.
 * 
 * @module story-tools
 * @version 1.0.0
 */

// Global state
let currentChatId = null;
let currentMemoryId = null;
let memoriesData = [];
let currentTab = 'all';

// Initialize on page load
document.addEventListener('DOMContentLoaded', function() {
    console.log('Story Tools: DOMContentLoaded event fired');
    initializeStoryTools();
});

async function initializeStoryTools() {
    console.log('Story Tools: initializeStoryTools() called');
    
    // Get current chat ID from URL or localStorage
    currentChatId = getCurrentChatId();
    console.log('Story Tools: currentChatId =', currentChatId);
    
    if (!currentChatId) {
        console.log('Story Tools: No active chat session found');
        showError('No active chat session found');
        return;
    }
    
    console.log('Story Tools: About to load recent memories');
    // Load initial memories data
    await loadRecentMemories();
    
    // Set up event listeners
    setupEventListeners();
    
    console.log('Story Tools: initialization complete');
    console.log('Story Tools initialized for chat:', currentChatId);
}

// Get current chat ID
function getCurrentChatId() {
    // Try to get from URL parameter first
    const urlParams = new URLSearchParams(window.location.search);
    const chatId = urlParams.get('chatId');
    
    if (chatId) {
        return chatId;
    }
    
    // Fallback to localStorage
    const currentChat = localStorage.getItem('howl-chat-current-chat');
    if (currentChat) {
        try {
            const chatData = JSON.parse(currentChat);
            return chatData.chatName || 'default';
        } catch (e) {
            console.error('Failed to parse current chat data:', e);
        }
    }
    
    return 'default';
}

// Setup event listeners
function setupEventListeners() {
    const searchInput = document.getElementById('memories-search');
    const filterSelect = document.getElementById('memories-filter-scope');
    
    if (searchInput) {
        searchInput.addEventListener('input', debounce(handleSearch, 300));
    }
    
    if (filterSelect) {
        filterSelect.addEventListener('change', handleFilter);
    }
}

// Load recent memories (last 5)
async function loadRecentMemories() {
    try {
        showLoading(true);
        const memories = await getMemories({ limit: 5, sort: 'newest' });
        renderRecentList(memories);
        showLoading(false);
    } catch (error) {
        console.error('Failed to load recent memories:', error);
        showError('Failed to load memories');
        showLoading(false);
    }
}

// Toggle memories flyout
async function toggleMemoriesFlyout() {
    const flyout = document.getElementById('memories-flyout');
    const overlay = document.getElementById('memories-flyout-overlay');
    const expandBtn = document.getElementById('memories-expand');
    
    if (!flyout || !overlay) return;
    
    const isOpen = flyout.classList.contains('show');
    
    if (isOpen) {
        // Close flyout
        flyout.classList.remove('show');
        overlay.classList.remove('show');
        if (expandBtn) {
            expandBtn.innerHTML = `
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <polyline points="9 18 15 12 12 5 6 5"></polyline>
                    <polyline points="9 6 21 6"></polyline>
                </svg>
            `;
        }
    } else {
        // Open flyout
        flyout.classList.add('show');
        overlay.classList.add('show');
        if (expandBtn) {
            expandBtn.innerHTML = `
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <polyline points="7 10 12 17 10 17"></polyline>
                    <polyline points="7 6 21 6"></polyline>
                </svg>
            `;
        }
        
        // Load full memories list
        await loadFullMemoriesList();
    }
}

// Load full memories list
async function loadFullMemoriesList() {
    try {
        showLoading(true);
        const memories = await getMemories({});
        renderFullList(memories);
        showLoading(false);
    } catch (error) {
        console.error('Failed to load full memories list:', error);
        showError('Failed to load memories');
        showLoading(false);
    }
}

// Get memories from backend
async function getMemories(filter = {}) {
    if (!window.go || !window.go.main || !window.go.main.App) {
        console.warn('Wails runtime not available, using mock data');
        // Return mock data for demonstration when Wails isn't available
        return [
            {
                id: '1',
                title: 'Character Introduction',
                summary: 'First meeting with the main character',
                tags: ['character', 'introduction'],
                pinned: true,
                trigger: 'system',
                created_at: new Date().toISOString()
            },
            {
                id: '2', 
                title: 'World Discovery',
                summary: 'Exploring the main setting and environment',
                tags: ['world', 'exploration'],
                pinned: false,
                trigger: 'user',
                created_at: new Date().toISOString()
            },
            {
                id: '3',
                title: 'Plot Development',
                summary: 'Key story progression moment',
                tags: ['plot', 'story'],
                pinned: false,
                trigger: 'system',
                created_at: new Date().toISOString()
            }
        ];
    }
    
    try {
        const memories = await window.go.main.App.GetMemories(currentChatId, filter);
        return memories || [];
    } catch (error) {
        console.error('Failed to get memories:', error);
        return [];
    }
}

// Get memory detail
async function getMemoryDetail(memoryId) {
    if (!window.go || !window.go.main || !window.go.main.App) {
        console.warn('Wails runtime not available, using mock data');
        // Return mock memory detail for demonstration
        const mockMemories = {
            '1': {
                id: '1',
                title: 'Character Introduction',
                summary: 'First meeting with the main character. This memory captures the initial encounter and establishes the character\'s personality and background.',
                tags: ['character', 'introduction'],
                pinned: true,
                trigger: 'system',
                created_at: new Date().toISOString()
            },
            '2': {
                id: '2',
                title: 'World Discovery',
                summary: 'Exploring the main setting and environment. This memory documents key locations and environmental details.',
                tags: ['world', 'exploration'],
                pinned: false,
                trigger: 'user',
                created_at: new Date().toISOString()
            },
            '3': {
                id: '3',
                title: 'Plot Development',
                summary: 'Key story progression moment. This memory captures important plot twists and story advancement.',
                tags: ['plot', 'story'],
                pinned: false,
                trigger: 'system',
                created_at: new Date().toISOString()
            }
        };
        return mockMemories[memoryId] || null;
    }
    
    try {
        const memory = await window.go.main.App.GetMemoryDetail(currentChatId, memoryId);
        return memory;
    } catch (error) {
        console.error('Error calling GetMemoryDetail:', error);
        throw error;
    }
}

// Render recent memories list
function renderRecentList(memories) {
    const container = document.getElementById('memories-recent');
    if (!container) return;
    
    if (!memories || memories.length === 0) {
        container.innerHTML = '<div class="memories-empty">No memories loaded</div>';
        return;
    }
    
    container.innerHTML = '';
    memories.forEach(memory => {
        const item = createRecentItem(memory);
        container.appendChild(item);
    });
}

// Create recent memory item
function createRecentItem(memory) {
    const item = document.createElement('div');
    item.className = 'memory-recent-item';
    item.dataset.memoryId = memory.id;
    
    const date = new Date(memory.lastTriggeredAt || memory.createdAt);
    const dateStr = date.toLocaleDateString();
    
    item.innerHTML = `
        <div class="memory-recent-title">${escapeHtml(memory.title)}</div>
        <div class="memory-recent-date">${dateStr}</div>
    `;
    
    item.addEventListener('click', () => selectMemory(memory.id));
    return item;
}

// Select memory (open in flyout)
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

// Open flyout
function openFlyout() {
    const flyout = document.getElementById('memories-flyout');
    if (flyout) {
        flyout.classList.add('show');
    }
}

// Render full memories list
function renderFullList(memories) {
    const container = document.getElementById('memories-list');
    if (!container) return;
    
    if (!memories || memories.length === 0) {
        container.innerHTML = '<div class="memories-empty">No memories found</div>';
        return;
    }
    
    container.innerHTML = '';
    memories.forEach(memory => {
        const item = createFullListItem(memory);
        container.appendChild(item);
    });
}

// Create full list item
function createFullListItem(memory) {
    const item = document.createElement('div');
    item.className = 'memory-item';
    item.dataset.memoryId = memory.id;
    
    const date = new Date(memory.lastTriggeredAt || memory.createdAt);
    const dateStr = date.toLocaleDateString();
    const isPinned = memory.pinned || false;
    const isSystem = memory.author === 'auto';
    
    item.innerHTML = `
        <div class="memory-item-title">${escapeHtml(memory.title)}</div>
        <div class="memory-item-summary">${escapeHtml(memory.summary || '')}</div>
        <div class="memory-item-meta">
            <div class="memory-item-tags">
                ${(memory.tags || []).map(tag => `<span class="memory-tag">${escapeHtml(tag)}</span>`).join('')}
            </div>
            <div class="memory-item-date">${dateStr}</div>
        </div>
    `;
    
    item.addEventListener('click', () => selectMemory(memory.id));
    return item;
}

// Populate editor with memory data
function populateEditor(memory) {
    currentMemoryId = memory.id;
    
    // Update form fields
    document.getElementById('mem-title').value = memory.title || '';
    document.getElementById('mem-summary').value = memory.summary || '';
    document.getElementById('mem-tags').value = (memory.tags || []).join(', ');
    document.getElementById('mem-pinned').checked = memory.pinned || false;
    
    // Update trigger info
    updateTriggerInfo(memory);
    
    // Update preview
    updatePreview(memory);
}

// Update trigger info display
function updateTriggerInfo(memory) {
    const triggerInfo = document.getElementById('mem-trigger-info');
    const triggerType = document.querySelector('.trigger-type');
    
    if (memory.author === 'auto') {
        triggerType.textContent = 'System trigger - not editable';
        triggerInfo.className = 'trigger-info';
    } else {
        triggerType.textContent = 'Manual entry - editable';
        triggerInfo.className = 'trigger-info manual-trigger';
    }
}

// Update preview
function updatePreview(memory) {
    const preview = document.getElementById('memories-preview');
    if (!preview) return;
    
    if (!memory.pinned) {
        preview.innerHTML = '<div class="preview-empty">Select a memory to view details</div>';
        return;
    }
    
    preview.innerHTML = `
        <div class="preview-title">${escapeHtml(memory.title)}</div>
        <div class="preview-summary">${escapeHtml(memory.summary || '')}</div>
        <div class="preview-tags">
            ${(memory.tags || []).map(tag => `<span class="memory-tag">${escapeHtml(tag)}</span>`).join('')}
        </div>
        <div class="preview-pinned">📌 Pinned to context</div>
    `;
}

// Switch memories tab
function switchMemoriesTab(tab) {
    currentTab = tab;
    
    // Update tab states
    document.querySelectorAll('.memories-tab').forEach(tabBtn => {
        tabBtn.classList.remove('active');
    });
    document.getElementById(`memories-tab-${tab}`).classList.add('active');
    
    // Reload list with filter
    loadFullMemoriesList();
}

// Handle search
async function handleSearch(event) {
    const query = event.target.value.trim();
    const filter = {
        query: query || undefined,
        scope: document.getElementById('memories-filter-scope').value
    };
    
    try {
        showLoading(true);
        const memories = await getMemories(filter);
        renderFullList(memories);
        showLoading(false);
    } catch (error) {
        console.error('Search failed:', error);
        showError('Search failed');
        showLoading(false);
    }
}

// Handle filter
async function handleFilter(event) {
    const scope = event.target.value;
    const filter = {
        scope: scope === 'all' ? undefined : scope,
        query: document.getElementById('memories-search').value.trim() || undefined
    };
    
    try {
        showLoading(true);
        const memories = await getMemories(filter);
        renderFullList(memories);
        showLoading(false);
    } catch (error) {
        console.error('Filter failed:', error);
        showError('Filter failed');
        showLoading(false);
    }
}

// Save memory
async function saveMemory() {
    if (!currentMemoryId) {
        showError('No memory selected to save');
        return;
    }
    
    const title = document.getElementById('mem-title').value.trim();
    const summary = document.getElementById('mem-summary').value.trim();
    const tags = document.getElementById('mem-tags').value.split(',').map(tag => tag.trim()).filter(tag => tag);
    const pinned = document.getElementById('mem-pinned').checked;
    
    if (!title) {
        showError('Memory title is required');
        return;
    }
    
    try {
        showLoading(true);
        const edits = {
            title,
            summary,
            tags,
            pinned
        };
        
        let updated;
        if (!window.go || !window.go.main || !window.go.main.App) {
            console.warn('Wails runtime not available, simulating save with mock data');
            // Simulate saving with mock data
            updated = { ...edits, id: currentMemoryId, created_at: new Date().toISOString() };
        } else {
            updated = await window.go.main.App.EditMemory(currentChatId, currentMemoryId, edits);
        }
        
        // Update local data
        const memoryIndex = memoriesData.findIndex(m => m.id === currentMemoryId);
        if (memoryIndex !== -1) {
            memoriesData[memoryIndex] = { ...memoriesData[memoryIndex], ...updated };
        }
        
        // Update UI
        renderFullList(memoriesData);
        showSuccess('Memory saved successfully');
        showLoading(false);
    } catch (error) {
        console.error('Save failed:', error);
        showError('Failed to save memory');
        showLoading(false);
    }
}

// Delete memory
async function deleteMemory() {
    if (!currentMemoryId) {
        showError('No memory selected to delete');
        return;
    }
    
    if (!confirm('Are you sure you want to delete this memory? This action cannot be undone.')) {
        return;
    }
    
    try {
        showLoading(true);
        if (!window.go || !window.go.main || !window.go.main.App) {
            console.warn('Wails runtime not available, simulating delete with mock data');
            // Simulate deletion by removing from local data only
        } else {
            await window.go.main.App.DeleteMemory(currentChatId, currentMemoryId);
        }
        
        // Update local data
        memoriesData = memoriesData.filter(m => m.id !== currentMemoryId);
        
        // Clear editor
        clearEditor();
        
        // Update UI
        renderFullList(memoriesData);
        showSuccess('Memory deleted successfully');
        showLoading(false);
    } catch (error) {
        console.error('Delete failed:', error);
        showError('Failed to delete memory');
        showLoading(false);
    }
}

// Export memories
async function exportMemories() {
    try {
        showLoading(true);
        const filter = {
            scope: document.getElementById('memories-filter-scope').value,
            query: document.getElementById('memories-search').value.trim() || undefined
        };
        
        const memories = await getMemories(filter);
        
        // Create download
        const dataStr = JSON.stringify(memories, null, 2);
        const blob = new Blob([dataStr], { type: 'application/json' });
        const url = URL.createObjectURL(blob);
        
        const a = document.createElement('a');
        a.href = url;
        a.download = `memories_${currentChatId}_${new Date().toISOString().split('T')[0]}.json`;
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

// Clear editor
function clearEditor() {
    currentMemoryId = null;
    document.getElementById('mem-title').value = '';
    document.getElementById('mem-summary').value = '';
    document.getElementById('mem-tags').value = '';
    document.getElementById('mem-pinned').checked = false;
    updateTriggerInfo({});
    updatePreview({});
}

// Show loading
function showLoading(show) {
    const overlay = document.getElementById('loading-overlay');
    if (overlay) {
        overlay.style.display = show ? 'flex' : 'none';
    }
}

// Show success message
function showSuccess(message) {
    showNotification(message, 'success');
}

// Show error message
function showError(message) {
    showNotification(message, 'error');
}

// Show notification
function showNotification(message, type = 'info') {
    // Create notification element
    const notification = document.createElement('div');
    notification.className = `notification notification-${type}`;
    notification.textContent = message;
    
    // Add to page
    document.body.appendChild(notification);
    
    // Auto remove after 3 seconds
    setTimeout(() => {
        if (notification.parentNode) {
            notification.parentNode.removeChild(notification);
        }
    }, 3000);
}

// Escape HTML to prevent XSS
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Debounce function
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// Export functions for global access
window.storyTools = {
    toggleMemoriesFlyout,
    switchMemoriesTab,
    saveMemory,
    deleteMemory,
    exportMemories,
    loadRecentMemories,
    selectMemory
};
