/**
 * Character Manager Module
 *
 * Handles character selection, management, and character-specific lorebook integration.
 * Manages character data and relationships with scenarios and worlds.
 *
 * @module character-manager
 * @version 1.0.0
 */

class CharacterManager {
    constructor() {
        this.characters = [];
        this.selectedLorebooks = [];
        this.editingCharacterId = null;
    }

    /**
     * Initialize character manager
     */
    initialize() {
        this.loadCharacters();
    }

    /**
     * Load characters from localStorage
     */
    loadCharacters() {
        try {
            this.characters = JSON.parse(localStorage.getItem('howl-chat-characters') || '[]');
        } catch (error) {
            console.error('Failed to load characters:', error);
            this.characters = [];
        }
    }

    /**
     * Save characters to localStorage
     */
    saveCharacters() {
        try {
            localStorage.setItem('howl-chat-characters', JSON.stringify(this.characters));
        } catch (error) {
            console.error('Failed to save characters:', error);
        }
    }

    /**
     * Get all characters
     * @returns {Array} Array of character objects
     */
    getCharacters() {
        return this.characters;
    }

    /**
     * Get character by ID
     * @param {string} characterId - Character ID to find
     * @returns {Object|null} Character object or null if not found
     */
    getCharacterById(characterId) {
        return this.characters.find(character => character.id === characterId) || null;
    }

    /**
     * Get characters attached to a scenario
     * @param {string} scenarioId - Scenario ID
     * @returns {Array} Array of attached characters
     */
    getCharactersForScenario(scenarioId) {
        const scenario = JSON.parse(localStorage.getItem('howl-chat-scenarios') || '[]')
            .find(s => s.id === scenarioId);
        if (!scenario || !scenario.characters) return [];

        return this.characters.filter(character =>
            scenario.characters.some(sc => sc.id === character.id)
        );
    }

    /**
     * Get characters attached to a world
     * @param {string} worldId - World ID
     * @returns {Array} Array of attached characters
     */
    getCharactersForWorld(worldId) {
        const world = JSON.parse(localStorage.getItem('howl-worlds') || '[]')
            .find(w => w.id === worldId);
        if (!world || !world.characters) return [];

        return this.characters.filter(character =>
            world.characters.some(wc => wc.id === character.id)
        );
    }

    /**
     * Open character creation popup
     * @param {string} characterId - Character ID for editing (null for creation)
     */
    openCharacterCreationPopup(characterId = null) {
        const overlay = document.getElementById('character-creation-overlay');
        if (overlay) {
            overlay.classList.add('show');
        }

        if (characterId) {
            this.editingCharacterId = characterId;
            this.populateEditForm(characterId);
        } else {
            this.editingCharacterId = null;
            this.clearCharacterForm();
        }
    }

    /**
     * Close character creation popup
     */
    closeCharacterCreationPopup() {
        const overlay = document.getElementById('character-creation-overlay');
        if (overlay) {
            overlay.classList.remove('show');
        }
        this.clearCharacterForm();
    }

    /**
     * Clear character creation form
     */
    clearCharacterForm() {
        const nameEl = document.getElementById('character-name');
        const roleEl = document.getElementById('character-role');
        const personalityEl = document.getElementById('character-personality');
        const backstoryEl = document.getElementById('character-backstory');

        if (nameEl) nameEl.value = '';
        if (roleEl) roleEl.value = '';
        if (personalityEl) personalityEl.value = '';
        if (backstoryEl) backstoryEl.value = '';

        this.selectedLorebooks = [];
        this.updateSelectedLorebooksDisplay();
    }

    /**
     * Populate edit form with character data
     * @param {string} characterId - Character ID to edit
     */
    populateEditForm(characterId) {
        const character = this.getCharacterById(characterId);
        if (!character) return;

        const nameEl = document.getElementById('character-name');
        const roleEl = document.getElementById('character-role');
        const personalityEl = document.getElementById('character-personality');
        const backstoryEl = document.getElementById('character-backstory');

        if (nameEl) nameEl.value = character.name || '';
        if (roleEl) roleEl.value = character.role || '';
        if (personalityEl) personalityEl.value = character.personality || '';
        if (backstoryEl) backstoryEl.value = character.backstory || '';

        this.selectedLorebooks = character.lorebooks || [];
        this.updateSelectedLorebooksDisplay();
    }

    /**
     * Update display of selected lorebooks
     */
    updateSelectedLorebooksDisplay() {
        const container = document.getElementById('selected-character-lorebooks');
        if (!container) return;

        container.innerHTML = '';

        if (this.selectedLorebooks.length === 0) {
            container.innerHTML = '<div class="no-items">No lorebooks selected</div>';
            return;
        }

        this.selectedLorebooks.forEach(lorebook => {
            const itemDiv = document.createElement('div');
            itemDiv.className = 'selected-item';
            itemDiv.innerHTML = `
                <span class="item-name">${lorebook.name || lorebook.title || 'Unnamed'}</span>
                <button class="remove-item" onclick="characterManager.removeSelectedLorebook('${lorebook.id}')">×</button>
            `;
            container.appendChild(itemDiv);
        });
    }

    /**
     * Remove selected lorebook
     * @param {string} lorebookId - Lorebook ID to remove
     */
    removeSelectedLorebook(lorebookId) {
        this.selectedLorebooks = this.selectedLorebooks.filter(l => l.id !== lorebookId);
        this.updateSelectedLorebooksDisplay();
    }

    /**
     * Open lorebook selector modal
     */
    openLorebookSelector() {
        const modal = document.createElement('div');
        modal.className = 'lorebook-selector-modal';
        modal.style.cssText = 'position:fixed;top:0;left:0;right:0;bottom:0;background:rgba(0,0,0,0.7);display:flex;align-items:center;justify-content:center;z-index:9999;backdrop-filter:blur(4px);';
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
        const lorebooks = JSON.parse(localStorage.getItem('howl-chat-lorebooks') || '[]');
        const selectedIds = this.selectedLorebooks.map(l => l.id);

        let content = `
            <div style="background: rgba(30, 30, 40, 0.98); border-radius: 12px; padding: 24px; max-width: 400px; border: 1px solid rgba(255, 255, 255, 0.1);">
                <h2 style="color: #fff; margin-top: 0; margin-bottom: 16px;">Select Lorebooks</h2>
                <div style="max-height: 300px; overflow-y: auto; margin-bottom: 16px;">
        `;

        if (lorebooks.length === 0) {
            content += '<div style="color: rgba(255, 255, 255, 0.5); text-align: center; padding: 20px;">No lorebooks available</div>';
        } else {
            lorebooks.forEach(lorebook => {
                const isSelected = selectedIds.includes(lorebook.id);
                content += `
                    <div onclick="characterManager.toggleLorebookSelection('${lorebook.id}');" style="
                        padding: 12px;
                        background: ${isSelected ? 'rgba(0, 184, 115, 0.2)' : 'rgba(255, 255, 255, 0.05)'};
                        margin-bottom: 8px;
                        border-radius: 6px;
                        cursor: pointer;
                        border: 1px solid ${isSelected ? 'rgba(0, 184, 115, 0.3)' : 'rgba(255, 255, 255, 0.1)'};
                        color: #fff;
                        transition: all 0.2s ease;
                    ">
                        <strong>${lorebook.title || lorebook.name || 'Unnamed'}</strong>
                        ${isSelected ? '<span style="float: right; color: #00b873;">✓</span>' : ''}
                    </div>
                `;
            });
        }

        content += `
                </div>
                <div style="display: flex; gap: 8px;">
                    <button onclick="this.closest('.lorebook-selector-modal').remove();" style="
                        flex: 1;
                        padding: 10px;
                        background: rgba(255, 255, 255, 0.1);
                        border: 1px solid rgba(255, 255, 255, 0.2);
                        color: #fff;
                        border-radius: 6px;
                        cursor: pointer;
                    ">Cancel</button>
                    <button onclick="characterManager.confirmLorebookSelection(); this.closest('.lorebook-selector-modal').remove();" style="
                        flex: 1;
                        padding: 10px;
                        background: rgba(0, 184, 115, 0.2);
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
     * Toggle lorebook selection
     * @param {string} lorebookId - Lorebook ID
     */
    toggleLorebookSelection(lorebookId) {
        const lorebooks = JSON.parse(localStorage.getItem('howl-chat-lorebooks') || '[]');
        const lorebook = lorebooks.find(l => l.id === lorebookId);
        if (!lorebook) return;

        const isSelected = this.selectedLorebooks.some(l => l.id === lorebookId);

        if (isSelected) {
            this.removeSelectedLorebook(lorebookId);
        } else {
            this.selectedLorebooks.push(lorebook);
            this.updateSelectedLorebooksDisplay();
        }
    }

    /**
     * Confirm lorebook selection
     */
    confirmLorebookSelection() {
        this.updateSelectedLorebooksDisplay();
    }

    /**
     * Save character
     */
    saveCharacter() {
        const nameEl = document.getElementById('character-name');
        const roleEl = document.getElementById('character-role');
        const personalityEl = document.getElementById('character-personality');
        const backstoryEl = document.getElementById('character-backstory');

        const name = nameEl ? nameEl.value.trim() : '';
        const role = roleEl ? roleEl.value.trim() : '';
        const personality = personalityEl ? personalityEl.value.trim() : '';
        const backstory = backstoryEl ? backstoryEl.value.trim() : '';

        if (!name) {
            alert('Character name is required');
            return;
        }

        const characterData = {
            id: this.editingCharacterId || Date.now().toString(),
            name: name,
            role: role,
            personality: personality,
            backstory: backstory,
            lorebooks: this.selectedLorebooks,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString()
        };

        if (this.editingCharacterId) {
            const index = this.characters.findIndex(c => c.id === this.editingCharacterId);
            if (index !== -1) {
                this.characters[index] = characterData;
            }
        } else {
            this.characters.push(characterData);
        }

        this.saveCharacters();
        this.closeCharacterCreationPopup();

        // Show success message
        if (typeof showNotification === 'function') {
            showNotification(`Character ${this.editingCharacterId ? 'updated' : 'created'}: ${name}`, 'success');
        }
    }

    /**
     * Delete character
     * @param {string} characterId - Character ID to delete
     */
    deleteCharacter(characterId) {
        const character = this.getCharacterById(characterId);
        if (!character) return;

        if (confirm(`Delete character "${character.name}"? This action cannot be undone.`)) {
            this.characters = this.characters.filter(c => c.id !== characterId);
            this.saveCharacters();

            if (typeof showNotification === 'function') {
                showNotification(`Character deleted: ${character.name}`, 'success');
            }
        }
    }

    /**
     * Get character-specific lorebooks
     * @param {string} characterId - Character ID
     * @returns {Array} Character's lorebooks
     */
    getCharacterLorebooks(characterId) {
        const character = this.getCharacterById(characterId);
        return character ? character.lorebooks || [] : [];
    }
}

// Create global instance
const characterManager = new CharacterManager();

// Initialize when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    characterManager.initialize();
});