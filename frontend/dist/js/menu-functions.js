/**
 * Menu Functions Implementation
 * 
 * This file contains all the missing dropdown menu functions
 * for the HOWL Chat application top menu bar.
 * 
 * @module menu-functions
 * @version 1.0.0
 */

// Edit Menu Functions
function editUndo() {
    if (window.wails && window.wails.main && window.wails.main.App) {
        // For now, just show a notification
        showNotification('Undo functionality coming soon', 'info');
    }
    console.log('Undo action requested');
}

function editRedo() {
    if (window.wails && window.wails.main && window.wails.main.App) {
        showNotification('Redo functionality coming soon', 'info');
    }
    console.log('Redo action requested');
}

function editCut() {
    if (window.getSelection) {
        const selection = window.getSelection();
        if (selection.toString()) {
            document.execCommand('cut');
            showNotification('Text cut to clipboard', 'success');
        } else {
            showNotification('No text selected', 'warning');
        }
    }
    console.log('Cut action requested');
}

function editCopy() {
    if (window.getSelection) {
        const selection = window.getSelection();
        if (selection.toString()) {
            document.execCommand('copy');
            showNotification('Text copied to clipboard', 'success');
        } else {
            showNotification('No text selected', 'warning');
        }
    }
    console.log('Copy action requested');
}

function editPaste() {
    if (navigator.clipboard && navigator.clipboard.readText) {
        navigator.clipboard.readText().then(text => {
            if (window.getSelection) {
                const selection = window.getSelection();
                if (selection.rangeCount > 0) {
                    const range = selection.getRangeAt(0);
                    range.deleteContents();
                    range.insertNode(document.createTextNode(text));
                    showNotification('Text pasted from clipboard', 'success');
                }
            }
        }).catch(err => {
            showNotification('Failed to paste from clipboard', 'error');
        });
    } else {
        document.execCommand('paste');
        showNotification('Paste attempted', 'info');
    }
    console.log('Paste action requested');
}

function editSelectAll() {
    if (window.getSelection) {
        const selection = window.getSelection();
        selection.selectAllChildren(document.body);
        showNotification('All text selected', 'success');
    }
    console.log('Select All action requested');
}

function editFind() {
    const searchTerm = prompt('Enter text to find:');
    if (searchTerm) {
        if (window.find) {
            const found = window.find(searchTerm);
            if (found) {
                showNotification(`Found: "${searchTerm}"`, 'success');
            } else {
                showNotification(`Not found: "${searchTerm}"`, 'warning');
            }
        } else {
            showNotification('Find functionality not supported in this browser', 'error');
        }
    }
    console.log('Find action requested');
}

function editReplace() {
    const findTerm = prompt('Enter text to find:');
    const replaceTerm = prompt('Enter replacement text:');
    
    if (findTerm && replaceTerm) {
        if (window.getSelection) {
            const selection = window.getSelection();
            const selectedText = selection.toString();
            
            if (selectedText === findTerm) {
                const range = selection.getRangeAt(0);
                range.deleteContents();
                range.insertNode(document.createTextNode(replaceTerm));
                showNotification('Text replaced', 'success');
            } else {
                showNotification('Select the text to replace first', 'warning');
            }
        }
    }
    console.log('Replace action requested');
}

// View Menu Functions
function toggleColorblindMode() {
    document.body.classList.toggle('colorblind-mode');
    const isActive = document.body.classList.contains('colorblind-mode');
    showNotification(`Colorblind mode ${isActive ? 'enabled' : 'disabled'}`, 'success');
    console.log('Colorblind mode toggled');
}

function toggleHighContrastMode() {
    document.body.classList.toggle('high-contrast-mode');
    const isActive = document.body.classList.contains('high-contrast-mode');
    showNotification(`High contrast mode ${isActive ? 'enabled' : 'disabled'}`, 'success');
    console.log('High contrast mode toggled');
}

function setTheme(theme) {
    document.body.className = document.body.className.replace(/theme-\w+/g, '');
    document.body.classList.add(`theme-${theme}`);
    localStorage.setItem('howl-chat-theme', theme);
    showNotification(`${theme.charAt(0).toUpperCase() + theme.slice(1)} theme applied`, 'success');
    console.log(`Theme set to: ${theme}`);
}

function increaseUIScale() {
    const currentScale = parseFloat(getComputedStyle(document.documentElement).getPropertyValue('--ui-scale') || '1.0');
    const newScale = Math.min(currentScale + 0.1, 2.0);
    document.documentElement.style.setProperty('--ui-scale', newScale.toString());
    localStorage.setItem('howl-chat-ui-scale', newScale.toString());
    showNotification(`UI scale increased to ${(newScale * 100).toFixed(0)}%`, 'success');
    console.log(`UI scale increased to: ${newScale}`);
}

function decreaseUIScale() {
    const currentScale = parseFloat(getComputedStyle(document.documentElement).getPropertyValue('--ui-scale') || '1.0');
    const newScale = Math.max(currentScale - 0.1, 0.5);
    document.documentElement.style.setProperty('--ui-scale', newScale.toString());
    localStorage.setItem('howl-chat-ui-scale', newScale.toString());
    showNotification(`UI scale decreased to ${(newScale * 100).toFixed(0)}%`, 'success');
    console.log(`UI scale decreased to: ${newScale}`);
}

function resetUIScale() {
    document.documentElement.style.setProperty('--ui-scale', '1.0');
    localStorage.setItem('howl-chat-ui-scale', '1.0');
    showNotification('UI scale reset to 100%', 'success');
    console.log('UI scale reset to default');
}

function toggleSidebar() {
    const sidebar = document.querySelector('.nav-sidebar');
    if (sidebar) {
        sidebar.classList.toggle('hidden');
        const isHidden = sidebar.classList.contains('hidden');
        showNotification(`Sidebar ${isHidden ? 'hidden' : 'shown'}`, 'success');
        console.log(`Sidebar ${isHidden ? 'hidden' : 'shown'}`);
    }
}

function toggleCharacterPanel() {
    const characterPanel = document.querySelector('.character-sidebar');
    if (characterPanel) {
        characterPanel.classList.toggle('hidden');
        const isHidden = characterPanel.classList.contains('hidden');
        showNotification(`Character panel ${isHidden ? 'hidden' : 'shown'}`, 'success');
        console.log(`Character panel ${isHidden ? 'hidden' : 'shown'}`);
    }
}

function toggleWorldPanel() {
    const worldPanel = document.querySelector('.world-sidebar');
    if (worldPanel) {
        worldPanel.classList.toggle('hidden');
        const isHidden = worldPanel.classList.contains('hidden');
        showNotification(`World panel ${isHidden ? 'hidden' : 'shown'}`, 'success');
        console.log(`World panel ${isHidden ? 'hidden' : 'shown'}`);
    }
}

function toggleFullscreen() {
    if (!document.fullscreenElement) {
        document.documentElement.requestFullscreen().then(() => {
            showNotification('Entered fullscreen mode', 'success');
        }).catch(err => {
            showNotification('Failed to enter fullscreen mode', 'error');
        });
    } else {
        document.exitFullscreen().then(() => {
            showNotification('Exited fullscreen mode', 'success');
        }).catch(err => {
            showNotification('Failed to exit fullscreen mode', 'error');
        });
    }
}

// Settings Menu Functions
function openStoryTools() {
    console.log('openStoryTools() called');
    toggleFlyout('story-tools-flyout');
}

// Story Tools Functions
function toggleMemoryList() {
    const memoryList = document.getElementById('memory-list-flyout');
    if (memoryList) {
        memoryList.classList.toggle('show');
        // Load memories when opening the list
        if (memoryList.classList.contains('show')) {
            loadMemoryList();
        }
    }
}

function editMemoryEntry(memoryId) {
    // Close memory list
    const memoryList = document.getElementById('memory-list-flyout');
    if (memoryList) {
        memoryList.classList.remove('show');
    }
    
    // Load memory data and open edit flyout
    loadMemoryForEdit(memoryId);
}

function loadMemoryForEdit(memoryId) {
    const editFlyout = document.getElementById('memory-edit-flyout');
    const titleField = document.getElementById('memory-title');
    const contentField = document.getElementById('memory-content');
    
    if (editFlyout && titleField && contentField) {
        // Store the current memory ID for saving
        titleField.setAttribute('data-memory-id', memoryId);
        
        // Load memory from backend
        if (window.go && window.go.main && window.go.main.App) {
            window.go.main.App.GetMemoryDetail('default', memoryId).then(memory => {
                if (memory) {
                    titleField.value = memory.title || '';
                    contentField.value = memory.summary || memory.content || '';
                }
                editFlyout.classList.add('show');
            }).catch(error => {
                console.error('Failed to load memory:', error);
                showNotification('Failed to load memory', 'error');
            });
        } else {
            console.warn('Wails runtime not available');
            showNotification('Backend not available', 'error');
        }
    }
}

function closeMemoryEdit() {
    const editFlyout = document.getElementById('memory-edit-flyout');
    if (editFlyout) {
        editFlyout.classList.remove('show');
    }
}

function saveMemoryEdit() {
    // Get the values
    const titleField = document.getElementById('memory-title');
    const contentField = document.getElementById('memory-content');
    
    if (titleField && contentField) {
        const title = titleField.value.trim();
        const content = contentField.value.trim();
        
        if (title && content) {
            // Save to backend
            if (window.go && window.go.main && window.go.main.App) {
                // Get current memory ID (you'd need to track this)
                const currentMemoryId = titleField.getAttribute('data-memory-id');
                if (currentMemoryId) {
                    const edits = { title, summary: content };
                    window.go.main.App.EditMemory('default', currentMemoryId, edits).then(() => {
                        showNotification('Memory saved successfully', 'success');
                        closeMemoryEdit();
                        // Refresh the memory list
                        loadMemoryList();
                    }).catch(error => {
                        console.error('Failed to save memory:', error);
                        showNotification('Failed to save memory', 'error');
                    });
                }
            } else {
                console.warn('Wails runtime not available');
                showNotification('Backend not available', 'error');
            }
        } else {
            showNotification('Please fill in both title and content', 'error');
        }
    }
}

function loadMemoryList() {
    const fullMemoryList = document.getElementById('full-memory-list');
    if (fullMemoryList) {
        fullMemoryList.innerHTML = '<div class="memory-loading">Loading memories...</div>';
        
        if (window.go && window.go.main && window.go.main.App) {
            window.go.main.App.GetMemories('default', {}).then(memories => {
                if (memories && memories.length > 0) {
                    renderMemoryList(memories);
                } else {
                    fullMemoryList.innerHTML = '<div class="memory-empty">No memories found</div>';
                }
            }).catch(error => {
                console.error('Failed to load memories:', error);
                fullMemoryList.innerHTML = '<div class="memory-error">Failed to load memories</div>';
            });
        } else {
            console.warn('Wails runtime not available');
            fullMemoryList.innerHTML = '<div class="memory-error">Backend not available</div>';
        }
    }
}

function renderMemoryList(memories) {
    const fullMemoryList = document.getElementById('full-memory-list');
    if (!fullMemoryList) return;
    
    fullMemoryList.innerHTML = memories.map(memory => `
        <div class="memory-list-item" onclick="editMemoryEntry('${memory.id}')">
            <div class="memory-item-time">${formatTime(memory.created_at)}</div>
            <div class="memory-item-title">${memory.title || 'Untitled'}</div>
        </div>
    `).join('');
}

function loadRecentMemories() {
    const recentMemories = document.getElementById('recent-memories');
    if (recentMemories) {
        recentMemories.innerHTML = '<div class="memory-loading">Loading memories...</div>';
        
        if (window.go && window.go.main && window.go.main.App) {
            window.go.main.App.GetMemories('default', {}).then(memories => {
                if (memories && memories.length > 0) {
                    // Show last 3 memories
                    const recent = memories.slice(-3).reverse();
                    renderRecentMemories(recent);
                } else {
                    recentMemories.innerHTML = '<div class="memory-empty">No memories found</div>';
                }
            }).catch(error => {
                console.error('Failed to load recent memories:', error);
                recentMemories.innerHTML = '<div class="memory-error">Failed to load memories</div>';
            });
        } else {
            console.warn('Wails runtime not available');
            recentMemories.innerHTML = '<div class="memory-error">Backend not available</div>';
        }
    }
}

function renderRecentMemories(memories) {
    const recentMemories = document.getElementById('recent-memories');
    if (!recentMemories) return;
    
    recentMemories.innerHTML = memories.map(memory => `
        <div class="memory-entry">
            <div class="memory-title">${memory.title || 'Untitled'}</div>
            <div class="memory-time">${formatTime(memory.created_at)}</div>
        </div>
    `).join('');
}

function formatTime(timestamp) {
    const date = new Date(timestamp);
    const now = new Date();
    const diffMs = now - date;
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);
    
    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins} minute${diffMins > 1 ? 's' : ''} ago`;
    if (diffHours < 24) return `${diffHours} hour${diffHours > 1 ? 's' : ''} ago`;
    return `${diffDays} day${diffDays > 1 ? 's' : ''} ago`;
}

// Initialize recent memories when Story Tools is opened
const originalOpenStoryTools = openStoryTools;
openStoryTools = function() {
    originalOpenStoryTools();
    loadRecentMemories();
};

function openAdvancedSettings() {
    showNotification('Advanced settings coming soon', 'info');
    console.log('Advanced settings requested');
}

function showKeyboardShortcuts() {
    const shortcuts = `
Keyboard Shortcuts:

General:
• Ctrl+N - New Chat
• Ctrl+S - Save Chat
• Ctrl+O - Load Chat
• Ctrl+Z - Undo
• Ctrl+Y - Redo
• Ctrl+C - Copy
• Ctrl+V - Paste
• Ctrl+A - Select All
• Ctrl+F - Find
• F11 - Fullscreen
• Esc - Close dialogs

Chat:
• Enter - Send message
• Shift+Enter - New line in message
• Ctrl+Enter - Continue generation
• Alt+R - Regenerate response

Navigation:
• Ctrl+1 - Chat
• Ctrl+2 - World
• Ctrl+3 - Characters
• Ctrl+4 - Scenarios
• Ctrl+5 - Lorebooks
• Ctrl+6 - Image Generator
    `;
    
    alert(shortcuts.trim());
    console.log('Keyboard shortcuts displayed');
}

function openAudioSettings() {
    showNotification('Audio settings coming soon', 'info');
    console.log('Audio settings requested');
}

function openTTSSettings() {
    showNotification('TTS settings coming soon', 'info');
    console.log('TTS settings requested');
}

// Utility Functions
function showNotification(message, type = 'info') {
    // Create notification element if it doesn't exist
    let notification = document.getElementById('notification');
    if (!notification) {
        notification = document.createElement('div');
        notification.id = 'notification';
        notification.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            padding: 12px 20px;
            border-radius: 8px;
            color: white;
            font-weight: 500;
            z-index: 10000;
            opacity: 0;
            transform: translateX(100%);
            transition: all 0.3s ease;
            max-width: 300px;
            word-wrap: break-word;
        `;
        document.body.appendChild(notification);
    }
    
    // Set color based on type
    const colors = {
        success: '#28a745',
        error: '#dc3545',
        warning: '#ffc107',
        info: '#17a2b8'
    };
    
    notification.style.backgroundColor = colors[type] || colors.info;
    notification.textContent = message;
    
    // Show notification
    notification.style.opacity = '1';
    notification.style.transform = 'translateX(0)';
    
    // Hide after 3 seconds
    setTimeout(() => {
        notification.style.opacity = '0';
        notification.style.transform = 'translateX(100%)';
    }, 3000);
}

// Initialize theme and UI scale from localStorage
document.addEventListener('DOMContentLoaded', function() {
    // Load saved theme
    const savedTheme = localStorage.getItem('howl-chat-theme');
    if (savedTheme) {
        setTheme(savedTheme);
    }
    
    // Load saved UI scale
    const savedScale = localStorage.getItem('howl-chat-ui-scale');
    if (savedScale) {
        document.documentElement.style.setProperty('--ui-scale', savedScale);
    }
});

// Export functions for global access
window.menuFunctions = {
    editUndo,
    editRedo,
    editCut,
    editCopy,
    editPaste,
    editSelectAll,
    editFind,
    editReplace,
    toggleColorblindMode,
    toggleHighContrastMode,
    setTheme,
    increaseUIScale,
    decreaseUIScale,
    resetUIScale,
    toggleSidebar,
    toggleCharacterPanel,
    toggleWorldPanel,
    toggleFullscreen,
    openStoryTools,
    openAdvancedSettings,
    showKeyboardShortcuts,
    openAudioSettings,
    openTTSSettings
};
