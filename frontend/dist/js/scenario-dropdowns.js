/**
 * Scenario Dropdowns
 *
 * Character and lorebook multi-select dropdowns for the
 * scenario creation form, including search and chip displays.
 *
 * Depends on: scenario-state.js
 *
 * @module scenario-dropdowns
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Generic Dropdown Toggle
// ---------------------------------------------------------------------------

function toggleDropdown(dropdownId) {
    const container = document.getElementById(dropdownId);
    if (!container) { console.warn(`Dropdown container ${dropdownId} not found`); return; }

    const menu = container.querySelector('.world-dropdown-menu');
    if (!menu) { console.warn(`Dropdown menu not found in ${dropdownId}`); return; }

    document.querySelectorAll('.world-dropdown-menu').forEach(d => {
        if (d.parentElement.id !== dropdownId) d.classList.remove('show');
    });

    const isNowShowing = !menu.classList.contains('show');
    menu.classList.toggle('show');

    if (isNowShowing) {
        if (dropdownId === 'character-lorebook-dropdown') {
            if (typeof window.characterLoadAvailableLorebooks === 'function') window.characterLoadAvailableLorebooks();
        } else if (dropdownId === 'scenario-lorebook-dropdown') {
            loadAvailableLorebooks();
        }
    }
}

// ---------------------------------------------------------------------------
// Character Dropdown
// ---------------------------------------------------------------------------

function loadAvailableCharacters() {
    const characters = JSON.parse(localStorage.getItem('howl-chat-characters') || '[]');
    window.scenarioState.availableCharacters = characters;
    updateCharacterDropdown();
}

function updateCharacterDropdown() {
    const dropdownItems = document.getElementById('scenario-character-dropdown-items');
    if (!dropdownItems) return;

    const ss = window.scenarioState;
    if (ss.availableCharacters.length === 0) {
        dropdownItems.innerHTML = '<div class="world-dropdown-item">No characters available</div>';
        return;
    }

    dropdownItems.innerHTML = ss.availableCharacters
        .filter(c => !ss.selectedCharacters.find(s => s.name === c.name))
        .map(c => `
            <div class="world-dropdown-item" onclick="addCharacterToScenario('${c.name}')">
                <div class="dropdown-item-content">
                    <div class="dropdown-item-title">${c.name}</div>
                    <div class="dropdown-item-subtitle">${c.metadata?.role || 'No role'}</div>
                </div>
            </div>
        `).join('');
}

function searchScenarioCharacters(searchTerm) {
    const ss = window.scenarioState;
    const filtered = ss.availableCharacters.filter(c =>
        !ss.selectedCharacters.find(s => s.name === c.name) &&
        (c.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
         (c.metadata?.role || '').toLowerCase().includes(searchTerm.toLowerCase()))
    );

    const dropdownItems = document.getElementById('scenario-character-dropdown-items');
    if (!dropdownItems) return;

    dropdownItems.innerHTML = filtered.length === 0
        ? '<div class="world-dropdown-item">No characters found</div>'
        : filtered.map(c => `
            <div class="world-dropdown-item" onclick="addCharacterToScenario('${c.name}')">
                <div class="dropdown-item-content">
                    <div class="dropdown-item-title">${c.name}</div>
                    <div class="dropdown-item-subtitle">${c.metadata?.role || 'No role'}</div>
                </div>
            </div>
        `).join('');
}

function addCharacterToScenario(characterName) {
    const ss = window.scenarioState;
    const character = ss.availableCharacters.find(c => c.name === characterName);
    if (character && !ss.selectedCharacters.find(c => c.name === characterName)) {
        ss.selectedCharacters.push(character);
        updateSelectedCharactersList();
        updateCharacterDropdown();
    }
}

function removeCharacterFromScenario(characterName) {
    const ss = window.scenarioState;
    ss.selectedCharacters = ss.selectedCharacters.filter(c => c.name !== characterName);
    updateSelectedCharactersList();
    updateCharacterDropdown();
}

function updateSelectedCharactersList() {
    const selectedItems = document.getElementById('scenario-character-selected-items');
    if (!selectedItems) return;
    const items = window.scenarioState.selectedCharacters;
    selectedItems.innerHTML = items.length === 0
        ? ''
        : items.map(c => `
            <div class="world-selected-item">
                <span class="selected-item-name">${c.name}</span>
                <button class="selected-item-remove" onclick="removeCharacterFromScenario('${c.name}')">×</button>
            </div>
        `).join('');
}

// ---------------------------------------------------------------------------
// Lorebook Dropdown
// ---------------------------------------------------------------------------

async function loadAvailableLorebooks() {
    const dropdownItems = document.getElementById('scenario-lorebook-dropdown-items');
    if (!dropdownItems) return;

    try {
        const ss = window.scenarioState;
        const lorebooks = JSON.parse(localStorage.getItem('howl-chat-lorebooks') || '[]')
            .map(ss.normalizeStoredLorebook).filter(Boolean);

        dropdownItems.innerHTML = '';

        if (lorebooks.length === 0) {
            dropdownItems.innerHTML = `<div class="world-dropdown-item disabled"><span class="world-dropdown-item-text">No lorebooks created yet</span></div>`;
            return;
        }

        lorebooks.forEach(lorebook => {
            if (!lorebook.id) lorebook.id = ss.lorebookKey(lorebook);
            const displayName = ss.lorebookDisplayName(lorebook);
            const isSelected  = ss.selectedLorebooks.some(l => ss.lorebookKey(l) === ss.lorebookKey(lorebook));

            const item = document.createElement('div');
            item.className = 'world-dropdown-item';
            if (isSelected) item.classList.add('selected');
            item.innerHTML = `<span class="world-dropdown-item-text">${displayName}</span>`;
            item.dataset.lorebookId = lorebook.id;
            item.addEventListener('click', () => toggleLorebookSelection(lorebook, item));
            dropdownItems.appendChild(item);
        });

    } catch (error) {
        console.error('Failed to load lorebooks:', error);
        dropdownItems.innerHTML = `<div class="world-dropdown-item disabled"><span class="world-dropdown-item-text">Error loading lorebooks</span></div>`;
    }
}

function searchScenarioLorebooks(searchTerm) {
    const ss = window.scenarioState;
    const all = JSON.parse(localStorage.getItem('howl-chat-lorebooks') || '[]')
        .map(ss.normalizeStoredLorebook).filter(Boolean);
    const term = searchTerm.toLowerCase();
    const filtered = all.filter(l => {
        const name = ss.lorebookDisplayName(l).toLowerCase();
        return name.includes(term);
    });

    const dropdownItems = document.getElementById('scenario-lorebook-dropdown-items');
    if (!dropdownItems) return;
    dropdownItems.innerHTML = '';

    if (filtered.length === 0) {
        dropdownItems.innerHTML = '<div class="world-dropdown-item">No lorebooks found</div>';
        return;
    }
    filtered.forEach(lorebook => {
        if (!lorebook.id) lorebook.id = ss.lorebookKey(lorebook);
        const displayName = ss.lorebookDisplayName(lorebook);
        const isSelected  = ss.selectedLorebooks.some(l => ss.lorebookKey(l) === ss.lorebookKey(lorebook));
        const item = document.createElement('div');
        item.className = 'world-dropdown-item';
        if (isSelected) item.classList.add('selected');
        item.innerHTML = `<span class="world-dropdown-item-text">${displayName}</span>`;
        item.dataset.lorebookId = lorebook.id;
        item.addEventListener('click', () => toggleLorebookSelection(lorebook, item));
        dropdownItems.appendChild(item);
    });
}

function toggleLorebookSelection(lorebook, itemElement) {
    const ss  = window.scenarioState;
    const key = ss.lorebookKey(lorebook);
    const idx = ss.selectedLorebooks.findIndex(l => ss.lorebookKey(l) === key);

    if (idx > -1) {
        ss.selectedLorebooks.splice(idx, 1);
        itemElement.classList.remove('selected');
    } else {
        if (!ss.selectedLorebooks.some(l => ss.lorebookKey(l) === key)) ss.selectedLorebooks.push(lorebook);
        itemElement.classList.add('selected');
    }

    updateSelectedLorebooksList();
    document.getElementById('scenario-lorebook-dropdown')?.classList.remove('show');
}

function updateSelectedLorebooksList() {
    const selectedItems = document.getElementById('scenario-lorebook-selected-items');
    if (!selectedItems) return;
    const ss = window.scenarioState;
    selectedItems.innerHTML = ss.selectedLorebooks.length === 0
        ? ''
        : ss.selectedLorebooks.map(l => `
            <div class="world-selected-item">
                <span class="selected-item-name">${ss.lorebookDisplayName(l)}</span>
                <button class="selected-item-remove" onclick="removeLorebookFromScenario('${ss.lorebookKey(l)}')">×</button>
            </div>
        `).join('');
}

function removeLorebookFromScenario(lorebookId) {
    const ss  = window.scenarioState;
    const idx = ss.selectedLorebooks.findIndex(l => ss.lorebookKey(l) === lorebookId);
    if (idx > -1) {
        ss.selectedLorebooks.splice(idx, 1);
        updateSelectedLorebooksList();
        loadAvailableLorebooks();
    }
}
