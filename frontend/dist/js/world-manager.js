/**
 * World Manager Module
 *
 * Handles world selection, management, and auto-population of related data.
 * Responsible for world hot-swapping and maintaining world-scenario-character-lorebook relationships.
 *
 * @module world-manager
 * @version 1.0.0
 */

class WorldManager {
    constructor() {
        this.worlds = [];
        this.selectedScenarios = [];
        this.selectedCharacters = [];
        this.selectedLorebooks = [];
        this.editingWorldId = null;
    }

    /**
     * Initialize world manager
     */
    initialize() {
        this.loadWorlds();
    }

    /**
     * Load worlds from localStorage
     */
    loadWorlds() {
        try {
            this.worlds = window.WorldRepository.all();
        } catch (error) {
            console.error('Failed to load worlds:', error);
            this.worlds = [];
        }
    }

    /**
     * Get all worlds
     * @returns {Array} Array of world objects
     */
    getWorlds() {
        return this.worlds;
    }

    /**
     * Get world by ID
     * @param {string} worldId - World ID to find
     * @returns {Object|null} World object or null if not found
     */
    getWorldById(worldId) {
        return this.worlds.find(world => world.id === worldId) || null;
    }

    /**
     * Open world creation popup
     * @param {string} worldId - World ID for editing (null for creation)
     */
    openWorldCreationPopup(worldId = null) {
        const overlay = document.getElementById('world-creation-overlay');
        overlay.classList.add('show');

        if (worldId) {
            this.editingWorldId = worldId;
            this.populateEditForm(worldId);
        } else {
            this.editingWorldId = null;
            this.clearWorldForm();
        }
    }

    /**
     * Close world creation popup
     */
    closeWorldCreationPopup() {
        const overlay = document.getElementById('world-creation-overlay');
        overlay.classList.remove('show');
        this.clearWorldForm();
    }

    /**
     * Clear world creation form
     */
    clearWorldForm() {
        document.getElementById('world-title').value = '';
        document.getElementById('world-details').value = '';
        document.getElementById('world-image-preview').style.display = 'none';
        document.getElementById('world-image-placeholder').style.display = 'flex';
        document.getElementById('world-image-remove').style.display = 'none';
        document.getElementById('world-image-input').value = '';
        this.selectedScenarios = [];
        this.selectedCharacters = [];
        this.selectedLorebooks = [];
        this.updateSelectedItemsDisplay('scenario');
        this.updateSelectedItemsDisplay('character');
        this.updateSelectedItemsDisplay('lorebook');
    }

    /**
     * Populate edit form with world data
     * @param {string} worldId - World ID to edit
     */
    populateEditForm(worldId) {
        const world = this.getWorldById(worldId);
        if (!world) return;

        document.getElementById('world-title').value = world.title || '';
        document.getElementById('world-details').value = world.details || '';

        if (world.image) {
            document.getElementById('world-image-preview').src = world.image;
            document.getElementById('world-image-preview').style.display = 'block';
            document.getElementById('world-image-placeholder').style.display = 'none';
            document.getElementById('world-image-remove').style.display = 'block';
        }

        this.selectedScenarios = world.scenarios || [];
        this.selectedCharacters = world.characters || [];
        this.selectedLorebooks = world.howl?.lorebooks || [];

        this.updateSelectedItemsDisplay('scenario');
        this.updateSelectedItemsDisplay('character');
        this.updateSelectedItemsDisplay('lorebook');
    }

    /**
     * Update display of selected items
     * @param {string} type - Type of items ('scenario', 'character', 'lorebook')
     */
    updateSelectedItemsDisplay(type) {
        const containerId = `selected-${type}s`;
        const container = document.getElementById(containerId);
        if (!container) return;

        const selectedItems = this.getSelectedItems(type);
        container.innerHTML = '';

        if (selectedItems.length === 0) {
            container.innerHTML = `<div class="no-items">No ${type}s selected</div>`;
            return;
        }

        selectedItems.forEach(item => {
            const itemDiv = document.createElement('div');
            itemDiv.className = 'selected-item';
            itemDiv.innerHTML = `
                <span class="item-name">${item.name || item.title || 'Unnamed'}</span>
                <button class="remove-item" onclick="worldManager.removeSelectedItem('${type}', '${item.id}')">×</button>
            `;
            container.appendChild(itemDiv);
        });
    }

    /**
     * Get selected items by type
     * @param {string} type - Type of items
     * @returns {Array} Selected items
     */
    getSelectedItems(type) {
        switch (type) {
            case 'scenario':
                return this.selectedScenarios;
            case 'character':
                return this.selectedCharacters;
            case 'lorebook':
                return this.selectedLorebooks;
            default:
                return [];
        }
    }

    /**
     * Remove selected item
     * @param {string} type - Type of item
     * @param {string} itemId - Item ID to remove
     */
    removeSelectedItem(type, itemId) {
        switch (type) {
            case 'scenario':
                this.selectedScenarios = this.selectedScenarios.filter(s => s.id !== itemId);
                break;
            case 'character':
                this.selectedCharacters = this.selectedCharacters.filter(c => c.id !== itemId);
                break;
            case 'lorebook':
                this.selectedLorebooks = this.selectedLorebooks.filter(l => l.id !== itemId);
                break;
        }
        this.updateSelectedItemsDisplay(type);
    }

    /**
     * Open item selector modal
     * @param {string} type - Type of items to select
     */
    openItemSelector(type) {
        const modal = document.createElement('div');
        modal.className = 'item-selector-modal';
        modal.style.cssText = 'position:fixed;top:0;left:0;right:0;bottom:0;background:var(--overlay-scrim);display:flex;align-items:center;justify-content:center;z-index:9999;backdrop-filter:blur(4px);';
        modal.innerHTML = this.createItemSelectorContent(type);
        document.body.appendChild(modal);

        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                modal.remove();
            }
        });
    }

    /**
     * Create item selector content
     * @param {string} type - Type of items
     * @returns {string} HTML content
     */
    createItemSelectorContent(type) {
        const items = this.getAvailableItems(type);
        const selectedIds = this.getSelectedItems(type).map(item => item.id);

        let content = `
            <div style="background: var(--shell-panel-bg); border-radius: 12px; padding: 24px; max-width: 400px; border: 1px solid var(--surface-hover-strong);">
                <h2 style="color: #fff; margin-top: 0; margin-bottom: 16px;">Select ${type.charAt(0).toUpperCase() + type.slice(1)}s</h2>
                <div style="max-height: 300px; overflow-y: auto; margin-bottom: 16px;">
        `;

        if (items.length === 0) {
            content += `<div style="color: var(--surface-hover); text-align: center; padding: 20px;">No ${type}s available</div>`;
        } else {
            items.forEach(item => {
                const isSelected = selectedIds.includes(item.id);
                content += `
                    <div onclick="worldManager.toggleItemSelection('${type}', '${item.id}');" style="
                        padding: 12px;
                        background: ${isSelected ? 'var(--success-bg)' : 'var(--surface-hover)'};
                        margin-bottom: 8px;
                        border-radius: 6px;
                        cursor: pointer;
                        border: 1px solid ${isSelected ? 'rgba(0, 184, 115, 0.3)' : 'rgba(255, 255, 255, 0.1)'};
                        color: #fff;
                        transition: all 0.2s ease;
                    ">
                        <strong>${item.title || item.name || 'Unnamed'}</strong>
                        ${isSelected ? '<span style="float: right; color: #00b873;">✓</span>' : ''}
                    </div>
                `;
            });
        }

        content += `
                </div>
                <div style="display: flex; gap: 8px;">
                    <button onclick="this.closest('.item-selector-modal').remove();" style="
                        flex: 1;
                        padding: 10px;
                        background: var(--surface-hover-strong);
                        border: 1px solid rgba(255, 255, 255, 0.2);
                        color: #fff;
                        border-radius: 6px;
                        cursor: pointer;
                    ">Cancel</button>
                    <button onclick="worldManager.confirmItemSelection('${type}'); this.closest('.item-selector-modal').remove();" style="
                        flex: 1;
                        padding: 10px;
                        background: var(--success-bg);
                        border: 1px solid rgba(0, 184, 115, 0.3);
                        color: #00b873;
                        border-radius: 6px;
                        cursor: pointer;
                    ">Done</button>
                </div>
            </div>
        `;

        return content;
    }

    /**
     * Get available items for selection
     * @param {string} type - Type of items
     * @returns {Array} Available items
     */
    getAvailableItems(type) {
        switch (type) {
            case 'scenario':
                return window.ScenarioRepository.all();
            case 'character':
                return (window.CharacterStorage?.all?.() || []);
            case 'lorebook':
                return JSON.parse(localStorage.getItem('howl-chat-lorebooks') || '[]');
            default:
                return [];
        }
    }

    /**
     * Toggle item selection
     * @param {string} type - Type of item
     * @param {string} itemId - Item ID
     */
    toggleItemSelection(type, itemId) {
        const items = this.getAvailableItems(type);
        const item = items.find(i => i.id === itemId);
        if (!item) return;

        const selectedItems = this.getSelectedItems(type);
        const isSelected = selectedItems.some(s => s.id === itemId);

        if (isSelected) {
            this.removeSelectedItem(type, itemId);
        } else {
            selectedItems.push(item);
            this.updateSelectedItemsDisplay(type);
        }
    }

    /**
     * Confirm item selection
     * @param {string} type - Type of items
     */
    confirmItemSelection(type) {
        this.updateSelectedItemsDisplay(type);
    }

    /**
     * Save world
     */
    async saveWorld() {
        const title = document.getElementById('world-title').value.trim();
        const details = document.getElementById('world-details').value.trim();

        if (!title) {
            alert('World title is required');
            return;
        }

        const existing = this.editingWorldId ? window.WorldRepository.get(this.editingWorldId) : null;
        const worldData = {
            ...(existing || {}),
            id: this.editingWorldId || undefined,
            title: title,
            details: details,
            image: document.getElementById('world-image-preview').src || null,
            scenarios: this.selectedScenarios,
            characters: this.selectedCharacters,
            howl: {
                lorebooks: this.selectedLorebooks
            },
            createdAt: existing?.createdAt || new Date().toISOString(),
            updatedAt: new Date().toISOString()
        };

        try {
            await window.WorldRepository.save(worldData);
            this.loadWorlds();
        } catch (error) {
            console.error('Failed to save world:', error);
            alert('Failed to save world: ' + error.message);
            return;
        }

        this.closeWorldCreationPopup();

        // Show success message
        if (typeof showNotification === 'function') {
            showNotification(`World ${this.editingWorldId ? 'updated' : 'created'}: ${title}`, 'success');
        }
    }

    /**
     * Delete world
     * @param {string} worldId - World ID to delete
     */
    async deleteWorld(worldId) {
        const world = this.getWorldById(worldId);
        if (!world) return;

        if (confirm(`Delete world "${world.title}"? This action cannot be undone.`)) {
            try {
                await window.WorldRepository.remove(worldId);
                this.loadWorlds();
            } catch (error) {
                console.error('Failed to delete world:', error);
                alert('Failed to delete world: ' + error.message);
                return;
            }

            if (typeof showNotification === 'function') {
                showNotification(`World deleted: ${world.title}`, 'success');
            }
        }
    }
}

// Create global instance
const worldManager = new WorldManager();

// Initialize when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    worldManager.initialize();
});