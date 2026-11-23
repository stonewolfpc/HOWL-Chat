/**
 * Lorebook Manager Module
 *
 * Handles lorebook selection, management, and source tracking.
 * Manages lorebooks from different sources (world, scenario, character) with proper visual indicators.
 *
 * @module lorebook-manager
 * @version 1.0.0
 */

class LorebookManager {
    constructor() {
        this.lorebooks = [];
        this.editingLorebookId = null;
        this.sourceColors = {
            world: '#00b873',      // Green
            scenario: '#007acc',   // Blue
            character: '#9c27b0'   // Purple
        };
    }

    /**
     * Initialize lorebook manager
     */
    initialize() {
        this.loadLorebooks();
    }

    /**
     * Load lorebooks from localStorage
     */
    loadLorebooks() {
        try {
            this.lorebooks = JSON.parse(localStorage.getItem('howl-chat-lorebooks') || '[]');
        } catch (error) {
            console.error('Failed to load lorebooks:', error);
            this.lorebooks = [];
        }
    }

    /**
     * Save lorebooks to localStorage
     */
    saveLorebooks() {
        try {
            localStorage.setItem('howl-chat-lorebooks', JSON.stringify(this.lorebooks));
        } catch (error) {
            console.error('Failed to save lorebooks:', error);
        }
    }

    /**
     * Get all lorebooks
     * @returns {Array} Array of lorebook objects
     */
    getLorebooks() {
        return this.lorebooks;
    }

    /**
     * Get lorebook by ID
     * @param {string} lorebookId - Lorebook ID to find
     * @returns {Object|null} Lorebook object or null if not found
     */
    getLorebookById(lorebookId) {
        return this.lorebooks.find(lorebook => lorebook.id === lorebookId) || null;
    }

    /**
     * Get lorebooks by source
     * @param {string} source - Source type ('world', 'scenario', 'character')
     * @param {string} sourceId - Source ID
     * @returns {Array} Array of lorebooks from the source
     */
    getLorebooksBySource(source, sourceId) {
        switch (source) {
            case 'world':
                const world = window.WorldRepository.all()
                    .find(w => w.id === sourceId);
                return world && world.howl && world.howl.lorebooks ?
                    this.lorebooks.filter(l => world.howl.lorebooks.some(wl => wl.id === l.id)) : [];

            case 'scenario':
                const scenario = window.ScenarioRepository.all()
                    .find(s => s.id === sourceId);
                return scenario && scenario.lorebooks ?
                    this.lorebooks.filter(l => scenario.lorebooks.some(sl => sl.id === l.id)) : [];

            case 'character':
                const character = (window.CharacterStorage?.all?.() || [])
                    .find(c => c.id === sourceId);
                return character && character.lorebooks ?
                    this.lorebooks.filter(l => character.lorebooks.some(cl => cl.id === l.id)) : [];

            default:
                return [];
        }
    }

    /**
     * Determine lorebook source and return source info
     * @param {string} lorebookId - Lorebook ID
     * @returns {Object} Source information {type, id, name}
     */
    getLorebookSource(lorebookId) {
        const worlds = window.WorldRepository.all();
        const scenarios = window.ScenarioRepository.all();
        const characters = (window.CharacterStorage?.all?.() || []);

        // Check worlds
        for (const world of worlds) {
            if (world.howl && world.howl.lorebooks &&
                world.howl.lorebooks.some(l => l.id === lorebookId)) {
                return { type: 'world', id: world.id, name: world.title || world.name };
            }
        }

        // Check scenarios
        for (const scenario of scenarios) {
            if (scenario.lorebooks && scenario.lorebooks.some(l => l.id === lorebookId)) {
                return { type: 'scenario', id: scenario.id, name: scenario.title || scenario.name };
            }
        }

        // Check characters
        for (const character of characters) {
            if (character.lorebooks && character.lorebooks.some(l => l.id === lorebookId)) {
                return { type: 'character', id: character.id, name: character.name };
            }
        }

        return { type: 'standalone', id: null, name: 'Standalone' };
    }

    /**
     * Get source color for lorebook
     * @param {string} lorebookId - Lorebook ID
     * @returns {string} Color hex code
     */
    getLorebookSourceColor(lorebookId) {
        const source = this.getLorebookSource(lorebookId);
        return this.sourceColors[source.type] || '#666666';
    }

    /**
     * Get source badge text for lorebook
     * @param {string} lorebookId - Lorebook ID
     * @returns {string} Badge text
     */
    getLorebookSourceBadge(lorebookId) {
        const source = this.getLorebookSource(lorebookId);
        return source.type.charAt(0).toUpperCase() + source.type.slice(1);
    }

    /**
     * Open lorebook creation popup
     * @param {string} lorebookId - Lorebook ID for editing (null for creation)
     */
    openLorebookCreationPopup(lorebookId = null) {
        const overlay = document.getElementById('lorebook-creation-overlay');
        if (overlay) {
            overlay.classList.add('show');
        }

        if (lorebookId) {
            this.editingLorebookId = lorebookId;
            this.populateEditForm(lorebookId);
        } else {
            this.editingLorebookId = null;
            this.clearLorebookForm();
        }
    }

    /**
     * Close lorebook creation popup
     */
    closeLorebookCreationPopup() {
        const overlay = document.getElementById('lorebook-creation-overlay');
        if (overlay) {
            overlay.classList.remove('show');
        }
        this.clearLorebookForm();
    }

    /**
     * Clear lorebook creation form
     */
    clearLorebookForm() {
        const titleEl = document.getElementById('lorebook-title');
        const contentEl = document.getElementById('lorebook-content');
        const tagsEl = document.getElementById('lorebook-tags');

        if (titleEl) titleEl.value = '';
        if (contentEl) contentEl.value = '';
        if (tagsEl) tagsEl.value = '';
    }

    /**
     * Populate edit form with lorebook data
     * @param {string} lorebookId - Lorebook ID to edit
     */
    populateEditForm(lorebookId) {
        const lorebook = this.getLorebookById(lorebookId);
        if (!lorebook) return;

        const titleEl = document.getElementById('lorebook-title');
        const contentEl = document.getElementById('lorebook-content');
        const tagsEl = document.getElementById('lorebook-tags');

        if (titleEl) titleEl.value = lorebook.title || '';
        if (contentEl) contentEl.value = lorebook.content || '';
        if (tagsEl) tagsEl.value = (lorebook.tags || []).join(', ');
    }

    /**
     * Save lorebook
     */
    saveLorebook() {
        const titleEl = document.getElementById('lorebook-title');
        const contentEl = document.getElementById('lorebook-content');
        const tagsEl = document.getElementById('lorebook-tags');

        const title = titleEl ? titleEl.value.trim() : '';
        const content = contentEl ? contentEl.value.trim() : '';
        const tagsString = tagsEl ? tagsEl.value.trim() : '';

        if (!title) {
            alert('Lorebook title is required');
            return;
        }

        const tags = tagsString ? tagsString.split(',').map(tag => tag.trim()).filter(tag => tag) : [];

        const lorebookData = {
            id: this.editingLorebookId || Date.now().toString(),
            title: title,
            content: content,
            tags: tags,
            scope: 'world', // Default scope
            enabled: true,
            triggerPhrases: [],
            secondaryTriggers: [],
            triggerMode: 'exact',
            triggerDirection: 'both',
            triggerFrequency: 'always',
            priorityLevel: 3,
            conflictRule: 'higher_priority',
            visibility: 'public',
            revealCondition: 'never',
            injectionStyle: 'summary',
            injectionPosition: 'system_summary_only',
            maxLength: 'medium',
            linkedCharacters: [],
            linkedLoreEntries: [],
            inheritance: 'none',
            automated: false,
            memoryType: 'world',
            triggerConfidence: 'medium',
            decayRate: 'never',
            contextBudget: 160,
            scanWindow: 4,
            recursion: 'limited',
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString()
        };

        if (this.editingLorebookId) {
            const index = this.lorebooks.findIndex(l => l.id === this.editingLorebookId);
            if (index !== -1) {
                this.lorebooks[index] = lorebookData;
            }
        } else {
            this.lorebooks.push(lorebookData);
        }

        this.saveLorebooks();
        this.closeLorebookCreationPopup();

        // Show success message
        if (typeof showNotification === 'function') {
            showNotification(`Lorebook ${this.editingLorebookId ? 'updated' : 'created'}: ${title}`, 'success');
        }
    }

    /**
     * Delete lorebook
     * @param {string} lorebookId - Lorebook ID to delete
     */
    deleteLorebook(lorebookId) {
        const lorebook = this.getLorebookById(lorebookId);
        if (!lorebook) return;

        window.DeleteConfirm.show({
            title: `Delete "${lorebook.title}"?`,
            body:  'The lorebook entry will be permanently removed and detached from all worlds, scenarios, and characters. Chat sessions will no longer include it.',
            onConfirm: async () => {
                if (window.DeleteService) await window.DeleteService.detachLorebookFromAll(lorebookId);
                this.lorebooks = this.lorebooks.filter(l => l.id !== lorebookId);
                this.saveLorebooks();
                if (typeof showNotification === 'function') {
                    showNotification(`Lorebook deleted: ${lorebook.title}`, 'success');
                }
            }
        });
    }

    /**
     * Create lorebook item element with source indicators
     * @param {Object} lorebook - Lorebook data
     * @returns {HTMLElement} Lorebook item element
     */
    createLorebookItem(lorebook) {
        const item = document.createElement('div');
        const isEnabled = lorebook.enabled !== false;
        item.className = 'lorebook-item' + (isEnabled ? '' : ' disabled');
        item.dataset.lorebookId = lorebook.id;

        const source = this.getLorebookSource(lorebook.id);
        const sourceColor = this.getLorebookSourceColor(lorebook.id);
        const sourceBadge = this.getLorebookSourceBadge(lorebook.id);
        const desc = this._resolveLorebookDescription(lorebook);

        item.innerHTML = `
            <div class="lorebook-icon" style="color: ${sourceColor}">📚</div>
            <div class="lorebook-info">
                <div class="lorebook-name">${lorebook.name || lorebook.title || 'Untitled Lorebook'}</div>
                ${desc ? `<div class="lorebook-description">${desc}</div>` : ''}
                <div class="lorebook-source" style="color: ${sourceColor}; font-size: 11px; margin-top: 4px;">
                    ${sourceBadge}: ${source.name}
                </div>
            </div>
            <button class="lorebook-menu-btn" onclick="lorebookManager.showLorebookMenu(event, '${lorebook.id}')">⋮</button>
            <button class="lorebook-remove" onclick="lorebookManager.removeLorebookFromChat('${lorebook.id}')">×</button>
        `;

        item.addEventListener('contextmenu', (e) => {
            e.preventDefault();
            this.showLorebookMenu(e, lorebook.id);
        });

        return item;
    }

    _resolveLorebookDescription(lorebook) {
        if (lorebook.description) return lorebook.description;
        try {
            const catalog = JSON.parse(localStorage.getItem('howl-chat-lorebooks') || '[]');
            const full = catalog.find(l => l.id === lorebook.id);
            return (full && full.description) ? full.description : '';
        } catch (_) { return ''; }
    }

    /**
     * Remove lorebook from chat
     * @param {string} lorebookId - Lorebook ID to remove
     */
    removeLorebookFromChat(lorebookId) {
        if (confirm('Remove this lorebook from the chat?')) {
            const currentChatData = window.contextDisplay?.currentChatData;
            if (currentChatData && currentChatData.lorebooks) {
                currentChatData.lorebooks = currentChatData.lorebooks.filter(l => l.id !== lorebookId);
                this.updateLorebookDisplay();
                if (window.contextDisplay && typeof window.contextDisplay.saveChatData === 'function') {
                    window.contextDisplay.saveChatData();
                }
                if (typeof showNotification === 'function') {
                    showNotification('Lorebook removed from chat', 'success');
                }
            }
        }
    }

    /**
     * Update lorebook display
     */
    updateLorebookDisplay() {
        const currentChatData = window.contextDisplay?.currentChatData;
        const lorebookList = document.getElementById('lorebook-list');
        console.log('[lorebookManager.updateLorebookDisplay] lorebooks:', (currentChatData?.lorebooks || []).length, 'list found:', !!lorebookList);
        if (!lorebookList) return;

        lorebookList.innerHTML = '';

        if (!currentChatData || !currentChatData.lorebooks || currentChatData.lorebooks.length === 0) {
            lorebookList.innerHTML = `
                <div class="lorebook-empty">
                    <div class="empty-icon">📚</div>
                    <div class="empty-text">No lorebooks selected</div>
                    <button class="add-lorebook-btn" onclick="lorebookManager.openLorebookSelector()">Add Lorebook</button>
                </div>
            `;
            return;
        }

        currentChatData.lorebooks.forEach(lorebook => {
            const lorebookItem = this.createLorebookItem(lorebook);
            lorebookList.appendChild(lorebookItem);
        });

        const addLorebookBtn = document.createElement('div');
        addLorebookBtn.className = 'add-lorebook-section';
        addLorebookBtn.innerHTML = '<button class="add-lorebook-btn" onclick="lorebookManager.openLorebookSelector()">Add Lorebook</button>';
        lorebookList.appendChild(addLorebookBtn);
    }

    /**
     * Open lorebook selector modal
     */
    openLorebookSelector() {
        console.log('Opening lorebook selector...');
        const modal = document.createElement('div');
        modal.className = 'lorebook-selector-modal';
        modal.style.cssText = 'position:fixed;top:0;left:0;right:0;bottom:0;background:var(--overlay-scrim);display:flex;align-items:center;justify-content:center;z-index:9999;backdrop-filter:blur(4px);';
        modal.innerHTML = this.createLorebookSelectorContent();
        document.body.appendChild(modal);

        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                modal.remove();
            }
        });
    }

    /**
     * Create lorebook selector content
     * @returns {string} HTML content
     */
    createLorebookSelectorContent() {
        const currentChatData = window.contextDisplay?.currentChatData || { lorebooks: [], world: null };
        const worlds = window.WorldRepository.all();
        const currentWorld = currentChatData.world;

        let attachedLorebooks = [];
        let availableLorebooks = [];

        if (currentWorld) {
            attachedLorebooks = this.getLorebooksBySource('world', currentWorld.id);
            availableLorebooks = this.lorebooks.filter(l =>
                !attachedLorebooks.some(al => al.id === l.id)
            );
        } else {
            availableLorebooks = this.lorebooks;
        }

        let content = `
            <div style="background: var(--shell-panel-bg); border-radius: 12px; padding: 24px; max-width: 600px; max-height: 80vh; border: 1px solid var(--surface-hover-strong); overflow-y: auto;">
                <h2 style="color: #fff; margin-top: 0; margin-bottom: 16px;">Select Lorebooks</h2>
        `;

        // Attached lorebooks section
        if (currentWorld && attachedLorebooks.length > 0) {
            content += `
                <div style="margin-bottom: 24px;">
                    <h3 style="color: #00b873; margin-bottom: 12px; font-size: 14px;">
                        📎 Attached to ${currentWorld.title || currentWorld.name}
                    </h3>
                    <div style="max-height: 200px; overflow-y: auto;">
            `;

            attachedLorebooks.forEach(lorebook => {
                const isActive = currentChatData.lorebooks &&
                    currentChatData.lorebooks.some(l => l.id === lorebook.id);
                content += `
                    <div class="lorebook-item" data-lorebook-id="${lorebook.id}" onclick="lorebookManager.toggleLorebookSelection('${lorebook.id}');" style="
                        padding: 12px;
                        background: ${isActive ? 'var(--success-bg)' : 'var(--surface-hover)'};
                        margin-bottom: 8px;
                        border-radius: 6px;
                        cursor: pointer;
                        border: 1px solid ${isActive ? 'rgba(0, 184, 115, 0.3)' : 'rgba(255, 255, 255, 0.1)'};
                        color: #fff;
                        transition: all 0.2s ease;
                        position: relative;
                    ">
                        ${isActive ? '<span style="position: absolute; top: 8px; right: 8px; color: #00b873;">✓</span>' : ''}
                        <strong>${lorebook.title || lorebook.name || 'Unnamed'}</strong>
                        <div style="color: var(--surface-hover); font-size: 12px; margin-top: 4px;">World Lorebook</div>
                    </div>
                `;
            });

            content += '</div></div>';
        }

        // Available lorebooks section
        content += `
            <div>
                <h3 style="color: #fff; margin-bottom: 12px; font-size: 14px;">
                    📚 Available Lorebooks
                </h3>
                <div style="max-height: 300px; overflow-y: auto;">
        `;

        if (availableLorebooks.length === 0) {
            content += '<div style="color: var(--surface-hover); text-align: center; padding: 16px; font-style: italic;">No additional lorebooks available</div>';
        } else {
            availableLorebooks.forEach(lorebook => {
                const isActive = currentChatData.lorebooks &&
                    currentChatData.lorebooks.some(l => l.id === lorebook.id);
                const source = this.getLorebookSource(lorebook.id);
                const sourceColor = this.getLorebookSourceColor(lorebook.id);

                content += `
                    <div class="lorebook-item" data-lorebook-id="${lorebook.id}" onclick="lorebookManager.toggleLorebookSelection('${lorebook.id}');" style="
                        padding: 12px;
                        background: ${isActive ? 'var(--success-bg)' : 'var(--surface-hover)'};
                        margin-bottom: 8px;
                        border-radius: 6px;
                        cursor: pointer;
                        border: 1px solid ${isActive ? 'rgba(0, 184, 115, 0.3)' : 'rgba(255, 255, 255, 0.1)'};
                        color: #fff;
                        transition: all 0.2s ease;
                        position: relative;
                    ">
                        ${isActive ? '<span style="position: absolute; top: 8px; right: 8px; color: #00b873;">✓</span>' : ''}
                        <strong>${lorebook.title || lorebook.name || 'Unnamed'}</strong>
                        <div style="color: ${sourceColor}; font-size: 12px; margin-top: 4px;">${source.type.charAt(0).toUpperCase() + source.type.slice(1)} Lorebook</div>
                    </div>
                `;
            });
        }

        content += `
                </div>
            </div>
            <div style="display: flex; gap: 8px; margin-top: 16px;">
                <button onclick="this.closest('.lorebook-selector-modal').remove();" style="
                    flex: 1;
                    padding: 10px;
                    background: var(--surface-hover-strong);
                    border: 1px solid rgba(255, 255, 255, 0.2);
                    color: #fff;
                    border-radius: 6px;
                    cursor: pointer;
                ">Cancel</button>
            </div>
        </div>
        `;

        return content;
    }

    /**
     * Toggle lorebook selection
     * @param {string} lorebookId - Lorebook ID
     */
    toggleLorebookSelection(lorebookId) {
        const currentChatData = window.contextDisplay?.currentChatData;
        if (!currentChatData) return;

        if (!currentChatData.lorebooks) {
            currentChatData.lorebooks = [];
        }

        const isSelected = currentChatData.lorebooks.some(l => l.id === lorebookId);

        if (isSelected) {
            currentChatData.lorebooks = currentChatData.lorebooks.filter(l => l.id !== lorebookId);
        } else {
            const lorebook = this.getLorebookById(lorebookId);
            if (lorebook) {
                currentChatData.lorebooks.push(lorebook);
            }
        }

        this.updateLorebookDisplay();
        if (window.contextDisplay && typeof window.contextDisplay.saveChatData === 'function') {
            window.contextDisplay.saveChatData();
        }
    }

    showLorebookMenu(event, lorebookId) {
        event.preventDefault();
        event.stopPropagation();

        const existingMenu = document.querySelector('.lorebook-context-menu');
        if (existingMenu) existingMenu.remove();

        const currentChatData = window.contextDisplay?.currentChatData;
        const lorebook = (currentChatData?.lorebooks || []).find(l => l.id === lorebookId);
        if (!lorebook) return;

        const isEnabled = lorebook.enabled !== false;
        const menu = document.createElement('div');
        menu.className = 'lorebook-context-menu';
        menu.innerHTML = `
            <div class="menu-item" onclick="lorebookManager.toggleLorebookEnabled('${lorebookId}')">
                <span class="icon icon-toggle"></span>
                ${isEnabled ? 'Disable Lorebook' : 'Enable Lorebook'}
            </div>
            <div class="menu-separator"></div>
            <div class="menu-item danger" onclick="lorebookManager.removeLorebookFromChat('${lorebookId}')">
                <span class="icon icon-trash"></span> Remove from Chat
            </div>
        `;

        menu.style.cssText = `position:fixed; left:${event.clientX}px; top:${event.clientY}px; z-index:10000;`;
        document.body.appendChild(menu);

        setTimeout(() => {
            document.addEventListener('click', function closeMenu(e) {
                if (!menu.contains(e.target)) {
                    menu.remove();
                    document.removeEventListener('click', closeMenu);
                }
            });
        }, 100);
    }

    async toggleLorebookEnabled(lorebookId) {
        const currentChatData = window.contextDisplay?.currentChatData;
        const lorebook = (currentChatData?.lorebooks || []).find(l => l.id === lorebookId);
        if (!lorebook) return;

        lorebook.enabled = lorebook.enabled === false;
        this.updateLorebookDisplay();
        if (window.contextDisplay && typeof window.contextDisplay.saveChatData === 'function') {
            window.contextDisplay.saveChatData();
        }

        const enabledIds = (currentChatData?.lorebooks || []).filter(l => l.enabled !== false).map(l => l.id);
        if (window.go?.main?.App?.UpdateLorebookContext) {
            try {
                await window.go.main.App.UpdateLorebookContext(enabledIds);
                console.log('[lorebookManager.toggleLorebookEnabled] synced', enabledIds.length, 'active lorebooks to backend');
            } catch (error) {
                console.warn('Backend lorebook enable sync failed:', error);
            }
        }

        const isEnabled = lorebook.enabled !== false;
        if (typeof showNotification === 'function') {
            showNotification(`Lorebook ${isEnabled ? 'enabled' : 'disabled'}`, 'success');
        }
    }
}

// Create global instance
const lorebookManager = new LorebookManager();

// Initialize when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    lorebookManager.initialize();
});