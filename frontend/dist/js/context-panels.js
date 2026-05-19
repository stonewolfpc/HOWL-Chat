/**
 * Context Panels — ContextDisplay Display Methods
 *
 * Mixin methods patched onto ContextDisplay.prototype for
 * rendering the world, scenario, character, and lorebook
 * panels in the sidebar. Also owns DOM clear helpers.
 *
 * Depends on: context-data.js
 *
 * @module context-panels
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Orchestrator
// ---------------------------------------------------------------------------

ContextDisplay.prototype.updateContextDisplay = async function () {
    await this.updateWorldDisplay();
    await this.updateScenarioDisplay();
    await this.updateCharacterPanel();
    this.updateLorebookDisplay();
};

// ---------------------------------------------------------------------------
// World Panel
// ---------------------------------------------------------------------------

ContextDisplay.prototype.updateWorldDisplay = async function () {
    const worldNameEl      = document.getElementById('world-name-display');
    const worldImageEl     = document.getElementById('world-image');
    const worldThumbnailEl = document.getElementById('world-thumbnail');
    if (!worldNameEl) return;

    if (this.currentChatData.world) {
        worldNameEl.textContent = this.currentChatData.world.title || this.currentChatData.world.name || 'Unnamed World';
        let src = this.currentChatData.world._resolvedImage || '';

        // New multi-image array format
        if (!src && this.currentChatData.world.images && Array.isArray(this.currentChatData.world.images) && this.currentChatData.world.images.length > 0) {
            const thumbIdx = this.currentChatData.world.thumbnailIndex || 0;
            const thumb = this.currentChatData.world.images[thumbIdx] || this.currentChatData.world.images[0];
            src = thumb.data || thumb.thumbnail || '';
        }

        // Legacy single-image format
        if (!src && this.currentChatData.world.image) {
            src = this.resolveImageSrc(this.currentChatData.world.image);
        }

        if (!src && this.currentChatData.world.image && window.HOWLStorage) {
            try {
                const ref     = this.currentChatData.world.image;
                const imageId = (typeof ref === 'object' && ref.id) ? ref.id : null;
                if (imageId) {
                    const imageData = await window.HOWLStorage.getImage(imageId, 'world');
                    if (imageData && imageData.data) {
                        src = imageData.data;
                        this.currentChatData.world._resolvedImage = src;
                    }
                }
            } catch (error) {
                console.warn('Failed to restore world image:', error);
            }
        }

        if (src && worldImageEl) {
            worldImageEl.src           = src;
            worldImageEl.style.display = 'block';
            const ph = worldThumbnailEl?.querySelector('.thumbnail-placeholder');
            if (ph) ph.style.display = 'none';
        } else if (worldImageEl) {
            worldImageEl.style.display = 'none';
            const ph = worldThumbnailEl?.querySelector('.thumbnail-placeholder');
            if (ph) ph.style.display = 'flex';
        }
    } else {
        worldNameEl.textContent = 'No World Selected';
        if (worldImageEl) worldImageEl.style.display = 'none';
        const ph = worldThumbnailEl?.querySelector('.thumbnail-placeholder');
        if (ph) ph.style.display = 'flex';
    }
};

// ---------------------------------------------------------------------------
// Scenario Panel
// ---------------------------------------------------------------------------

ContextDisplay.prototype.updateScenarioDisplay = async function () {
    const scenarioNameEl      = document.getElementById('scenario-name-display');
    const scenarioImageEl     = document.getElementById('scenario-image');
    const scenarioThumbnailEl = document.getElementById('scenario-thumbnail');
    if (!scenarioNameEl) return;

    if (this.currentChatData.scenario) {
        scenarioNameEl.textContent = this.currentChatData.scenario.title || this.currentChatData.scenario.name || 'Unnamed Scenario';
        let src = this.currentChatData.scenario._resolvedImage || '';

        // New multi-image array format
        if (!src && this.currentChatData.scenario.images && Array.isArray(this.currentChatData.scenario.images) && this.currentChatData.scenario.images.length > 0) {
            const thumbIdx = this.currentChatData.scenario.thumbnailIndex || 0;
            const thumb = this.currentChatData.scenario.images[thumbIdx] || this.currentChatData.scenario.images[0];
            src = thumb.data || thumb.thumbnail || '';
        }

        // Legacy single-image format
        if (!src && this.currentChatData.scenario.image) {
            src = this.resolveImageSrc(this.currentChatData.scenario.image);
        }

        if (!src && this.currentChatData.scenario.image && window.HOWLStorage) {
            try {
                const ref     = this.currentChatData.scenario.image;
                const imageId = (typeof ref === 'object' && ref.id) ? ref.id : null;
                if (imageId) {
                    const imageData = await window.HOWLStorage.getImage(imageId, 'scenario');
                    if (imageData && imageData.data) {
                        src = imageData.data;
                        this.currentChatData.scenario._resolvedImage = src;
                    }
                }
            } catch (error) {
                console.warn('Failed to restore scenario image:', error);
            }
        }

        if (src && scenarioImageEl) {
            scenarioImageEl.src           = src;
            scenarioImageEl.style.display = 'block';
            const ph = scenarioThumbnailEl?.querySelector('.thumbnail-placeholder');
            if (ph) ph.style.display = 'none';
        } else if (scenarioImageEl) {
            scenarioImageEl.style.display = 'none';
            const ph = scenarioThumbnailEl?.querySelector('.thumbnail-placeholder');
            if (ph) ph.style.display = 'flex';
        }
    } else {
        scenarioNameEl.textContent = 'No Scenario Selected';
        if (scenarioImageEl) scenarioImageEl.style.display = 'none';
        const ph = scenarioThumbnailEl?.querySelector('.thumbnail-placeholder');
        if (ph) ph.style.display = 'flex';
    }
};

// ---------------------------------------------------------------------------
// Character Panel
// ---------------------------------------------------------------------------

ContextDisplay.prototype.updateCharacterPanel = async function () {
    const characterList = document.getElementById('character-list');
    if (!characterList) return;
    characterList.innerHTML = '';

    if (!this.currentChatData.characters || this.currentChatData.characters.length === 0) {
        characterList.innerHTML = `
            <div class="character-empty">
                <div class="empty-icon">👥</div>
                <div class="empty-text">No characters selected</div>
                <button class="add-character-btn" onclick="contextDisplay.openCharacterSelector()">Add Character</button>
            </div>
        `;
        return;
    }

    for (const character of this.currentChatData.characters) {
        const item = await this.createCharacterItem(character);
        characterList.appendChild(item);
    }

    const addSection = document.createElement('div');
    addSection.className = 'add-character-section';
    addSection.innerHTML = '<button class="add-character-btn" onclick="contextDisplay.openCharacterSelector()">Add Character</button>';
    characterList.appendChild(addSection);
};

ContextDisplay.prototype.createCharacterItem = async function (character) {
    const item = document.createElement('div');
    item.className = 'character-item';
    item.dataset.characterId = character.id;

    const statusClass = character.status || (window.sidebarState?.CharacterStatus?.OFFLINE || 'offline');

    let avatarSrc = character._resolvedImage || this.resolveImageSrc(character.image || character.avatar);

    if (!avatarSrc && character.id && window.CharacterStorage) {
        try {
            const images = await window.CharacterStorage.loadImages(character.id);
            const thumb = images.find(img => img.isThumbnail) || images[0];
            if (thumb && thumb.data) {
                avatarSrc = thumb.data;
                character._resolvedImage = thumb.data;
            }
        } catch (_) {}
    }

    if (!avatarSrc && character.howl?.images?.length > 0) {
        const thumb = character.howl.images.find(img => img.isThumbnail) || character.howl.images[0];
        if (thumb && thumb.data) {
            avatarSrc = thumb.data;
            character._resolvedImage = thumb.data;
        }
    }

    const role = character.role || character.metadata?.role || '';

    item.innerHTML = `
        <div class="character-avatar" onclick="openCharacterSheet('${character.id}')">
            ${avatarSrc
                ? `<img src="${avatarSrc}" alt="${character.name}" onerror="this.style.display='none'; this.nextElementSibling.style.display='flex';">
                   <div class="avatar-placeholder" style="display:none;">👤</div>`
                : `<div class="avatar-placeholder">👤</div>`
            }
            <div class="status-indicator ${statusClass}"></div>
        </div>
        <div class="character-info" onclick="openCharacterSheet('${character.id}')">
            <div class="character-name">${character.name || 'Unnamed Character'}</div>
            ${role ? `<div class="character-role">${role}</div>` : ''}
        </div>
        <div class="character-actions">
            <button class="character-menu-btn" onclick="showCharacterMenu(event, '${character.id}')">⋮</button>
        </div>
    `;

    if (character.status === 'stepped-out' || character.status === 'dead') {
        item.classList.add('character-inactive');
    }

    item.addEventListener('contextmenu', (e) => {
        e.preventDefault();
        showCharacterMenu(e, character.id);
    });

    return item;
};

// ---------------------------------------------------------------------------
// Lorebook Panel
// ---------------------------------------------------------------------------

ContextDisplay.prototype.updateLorebookDisplay = function () {
    if (typeof lorebookManager !== 'undefined') lorebookManager.updateLorebookDisplay();
};

// ---------------------------------------------------------------------------
// DOM Helpers
// ---------------------------------------------------------------------------

ContextDisplay.prototype.clearChatMessagesDom = function () {
    const chatHistory = document.querySelector('.chat-messages');
    if (chatHistory) chatHistory.innerHTML = '';
    const input = document.querySelector('.input-field');
    if (input) input.value = '';
    document.getElementById('generate-banner')?.remove();
};
