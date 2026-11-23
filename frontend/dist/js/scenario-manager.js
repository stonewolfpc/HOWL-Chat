/**
 * Scenario Manager Module
 *
 * Handles scenario selection, management, and auto-population of related characters and lorebooks.
 * Manages scenario-specific data and relationships.
 *
 * @module scenario-manager
 * @version 1.0.0
 */

class ScenarioManager {
    constructor() {
        this.scenarios = [];
        this.selectedCharacters = [];
        this.selectedLorebooks = [];
        this.editingScenarioId = null;
    }

    /**
     * Initialize scenario manager
     */
    initialize() {
        this.loadScenarios();
    }

    /**
     * Load scenarios from localStorage
     */
    loadScenarios() {
        try {
            this.scenarios = JSON.parse(localStorage.getItem('howl-chat-scenarios') || '[]');
        } catch (error) {
            console.error('Failed to load scenarios:', error);
            this.scenarios = [];
        }
    }

    /**
     * Save scenarios to localStorage
     */
    saveScenarios() {
        try {
            localStorage.setItem('howl-chat-scenarios', JSON.stringify(this.scenarios));
        } catch (error) {
            console.error('Failed to save scenarios:', error);
        }
    }

    /**
     * Get all scenarios
     * @returns {Array} Array of scenario objects
     */
    getScenarios() {
        return this.scenarios;
    }

    /**
     * Get scenario by ID
     * @param {string} scenarioId - Scenario ID to find
     * @returns {Object|null} Scenario object or null if not found
     */
    getScenarioById(scenarioId) {
        return this.scenarios.find(scenario => scenario.id === scenarioId) || null;
    }

    /**
     * Open scenario creation popup
     * @param {string} scenarioId - Scenario ID for editing (null for creation)
     */
    openScenarioCreationPopup(scenarioId = null) {
        const overlay = document.getElementById('scenario-creation-overlay');
        if (overlay) {
            overlay.classList.add('show');
        }

        if (scenarioId) {
            this.editingScenarioId = scenarioId;
            this.populateEditForm(scenarioId);
        } else {
            this.editingScenarioId = null;
            this.clearScenarioForm();
        }
    }

    /**
     * Close scenario creation popup
     */
    closeScenarioCreationPopup() {
        const overlay = document.getElementById('scenario-creation-overlay');
        if (overlay) {
            overlay.classList.remove('show');
        }
        this.clearScenarioForm();
    }

    /**
     * Clear scenario creation form
     */
    clearScenarioForm() {
        const titleEl = document.getElementById('scenario-title');
        const detailsEl = document.getElementById('scenario-details');
        const objectiveEl = document.getElementById('scenario-objective');

        if (titleEl) titleEl.value = '';
        if (detailsEl) detailsEl.value = '';
        if (objectiveEl) objectiveEl.value = '';

        this.selectedCharacters = [];
        this.selectedLorebooks = [];
        this.updateSelectedItemsDisplay('character');
        this.updateSelectedItemsDisplay('lorebook');
    }

    /**
     * Populate edit form with scenario data
     * @param {string} scenarioId - Scenario ID to edit
     */
    populateEditForm(scenarioId) {
        const scenario = this.getScenarioById(scenarioId);
        if (!scenario) return;

        const titleEl = document.getElementById('scenario-title');
        const detailsEl = document.getElementById('scenario-details');
        const objectiveEl = document.getElementById('scenario-objective');

        if (titleEl) titleEl.value = scenario.title || '';
        if (detailsEl) detailsEl.value = scenario.details || '';
        if (objectiveEl) objectiveEl.value = scenario.objective || '';

        this.selectedCharacters = scenario.characters || [];
        this.selectedLorebooks = scenario.lorebooks || [];

        this.updateSelectedItemsDisplay('character');
        this.updateSelectedItemsDisplay('lorebook');
    }

    /**
     * Update display of selected items
     * @param {string} type - Type of items ('character', 'lorebook')
     */
    updateSelectedItemsDisplay(type) {
        const containerId = `selected-scenario-${type}s`;
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
                <button class="remove-item" onclick="scenarioManager.removeSelectedItem('${type}', '${item.id}')">×</button>
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
                    <div onclick="scenarioManager.toggleItemSelection('${type}', '${item.id}');" style="
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
                    <button onclick="scenarioManager.confirmItemSelection('${type}'); this.closest('.item-selector-modal').remove();" style="
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
     * Save scenario
     */
    saveScenario() {
        const titleEl = document.getElementById('scenario-title');
        const detailsEl = document.getElementById('scenario-details');
        const objectiveEl = document.getElementById('scenario-objective');

        const title = titleEl ? titleEl.value.trim() : '';
        const details = detailsEl ? detailsEl.value.trim() : '';
        const objective = objectiveEl ? objectiveEl.value.trim() : '';

        if (!title) {
            alert('Scenario title is required');
            return;
        }

        const scenarioData = {
            id: this.editingScenarioId || Date.now().toString(),
            title: title,
            details: details,
            objective: objective,
            characters: this.selectedCharacters,
            lorebooks: this.selectedLorebooks,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString()
        };

        if (this.editingScenarioId) {
            const index = this.scenarios.findIndex(s => s.id === this.editingScenarioId);
            if (index !== -1) {
                this.scenarios[index] = scenarioData;
            }
        } else {
            this.scenarios.push(scenarioData);
        }

        this.saveScenarios();
        this.closeScenarioCreationPopup();

        // Show success message
        if (typeof showNotification === 'function') {
            showNotification(`Scenario ${this.editingScenarioId ? 'updated' : 'created'}: ${title}`, 'success');
        }
    }

    /**
     * Delete scenario
     * @param {string} scenarioId - Scenario ID to delete
     */
    deleteScenario(scenarioId) {
        const scenario = this.getScenarioById(scenarioId);
        if (!scenario) return;

        if (confirm(`Delete scenario "${scenario.title}"? This action cannot be undone.`)) {
            this.scenarios = this.scenarios.filter(s => s.id !== scenarioId);
            this.saveScenarios();

            if (typeof showNotification === 'function') {
                showNotification(`Scenario deleted: ${scenario.title}`, 'success');
            }
        }
    }
}

// Create global instance
const scenarioManager = new ScenarioManager();

// Initialize when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    scenarioManager.initialize();
});