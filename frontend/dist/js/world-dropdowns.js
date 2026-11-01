/**
 * World Form Dropdowns
 *
 * Loads scenario, character, and lorebook data into their
 * respective multi-select dropdowns, handles toggle/deselect,
 * and updates the selected-items chip displays.
 *
 * Depends on: world-state.js
 *
 * @module world-dropdowns
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Generic Dropdown Toggle
// ---------------------------------------------------------------------------

function toggleDropdown(dropdownId) {
    const dropdown = document.getElementById(dropdownId);
    const isOpen   = dropdown.classList.contains('open');

    document.querySelectorAll('.world-dropdown-container').forEach(d => {
        if (d.id !== dropdownId) d.classList.remove('open');
    });
    dropdown.classList.toggle('open');

    if (dropdownId === 'scenario-dropdown'     && !isOpen) loadAvailableScenarios();
    if (dropdownId === 'character-dropdown'    && !isOpen) loadAvailableCharacters();
    if (dropdownId === 'world-lorebook-dropdown' && !isOpen) loadAvailableLorebooks();
}

// ---------------------------------------------------------------------------
// Generic Selected-Items Display (scenario / character)
// ---------------------------------------------------------------------------

function updateSelectedItemsDisplay(type) {
    const container = document.getElementById(`${type}-selected-items`);
    const items     = type === 'scenario'
        ? window.worldState.selectedScenarios
        : window.worldState.selectedCharacters;

    container.innerHTML = items.map(item => `
        <div class="world-selected-item">
            <span>${item.name || item.title || 'Untitled'}</span>
            <button class="world-selected-item-remove" onclick="removeSelectedItem('${type}', '${item.id}')">×</button>
        </div>
    `).join('');

    const triggerText = document.getElementById(`${type}-dropdown-text`);
    if (triggerText) {
        triggerText.textContent = items.length > 0
            ? `${items.length} ${type}${items.length === 1 ? '' : 's'} selected`
            : `Select ${type}s...`;
    }
}

function removeSelectedItem(type, itemId) {
    if (type === 'scenario') {
        window.worldState.selectedScenarios = window.worldState.selectedScenarios.filter(s => s.id !== itemId);
    } else {
        window.worldState.selectedCharacters = window.worldState.selectedCharacters.filter(c => c.id !== itemId);
    }
    updateSelectedItemsDisplay(type);
}

// ---------------------------------------------------------------------------
// Dropdown Cache & Helpers
// ---------------------------------------------------------------------------

window._worldDropdownCache = { scenarios: [], characters: [], lorebooks: [] };

function _worldDropdownSortSelectHtml(id) {
    return `<select class="world-dropdown-sort" style="padding:6px 8px;border:1px solid var(--surface-hover-strong);border-radius:4px;font-size:12px;cursor:pointer;"
        onchange="window._worldDropdownSort && window._worldDropdownSort('${id}', this.value)">
        <option value="name:asc">Name A-Z</option>
        <option value="name:desc">Name Z-A</option>
        <option value="date:desc">Newest First</option>
        <option value="date:asc">Oldest First</option>
    </select>`;
}

// ---------------------------------------------------------------------------
// Scenario Dropdown
// ---------------------------------------------------------------------------

function renderScenarioDropdown(items) {
    const dropdownItems = document.getElementById('scenario-dropdown-items');
    if (!dropdownItems) return;
    dropdownItems.innerHTML = '';

    const ws = window.worldState;
    if (items.length === 0) {
        dropdownItems.innerHTML = `<div class="world-dropdown-item disabled"><span class="world-dropdown-item-text">No scenarios found</span></div>`;
        return;
    }

    items.forEach((scenario) => {
        const displayTitle = scenario.title || scenario.name || 'Untitled';
        const isSelected = ws.selectedScenarios.some(s => s.id === scenario.id);
        const item = document.createElement('div');
        item.className = 'world-dropdown-item' + (isSelected ? ' selected' : '');
        item.innerHTML = `<span class="world-dropdown-item-text">${displayTitle}</span>`;
        item.dataset.scenarioId = scenario.id;
        item.addEventListener('click', () => toggleScenarioSelection(scenario, item));
        dropdownItems.appendChild(item);
    });
}

async function loadAvailableScenarios() {
    const dropdownItems = document.getElementById('scenario-dropdown-items');
    if (!dropdownItems) return;

    try {
        const scenarios = JSON.parse(localStorage.getItem('howl-chat-scenarios') || '[]');
        if (scenarios.length === 0) {
            dropdownItems.innerHTML = `<div class="world-dropdown-item disabled"><span class="world-dropdown-item-text">No scenarios created yet</span></div>`;
            window._worldDropdownCache.scenarios = [];
            return;
        }

        const processed = await Promise.all(scenarios.map(async (scenario) => {
            const p = { ...scenario };
            if (scenario.image && scenario.image.stored && window.HOWLStorage) {
                try {
                    const img = await window.HOWLStorage.getImage(scenario.image.id, 'scenario');
                    p.image = (img && img.data) ? img : null;
                } catch { p.image = null; }
            }
            return p;
        }));

        window._worldDropdownCache.scenarios = processed;
        const fe = window.FilterEngine;
        const sorted = fe ? fe.sortItems(processed, 'name', 'asc') : processed;
        renderScenarioDropdown(sorted);

    } catch (error) {
        console.error('Failed to load scenarios:', error);
        dropdownItems.innerHTML = `<div class="world-dropdown-item disabled"><span class="world-dropdown-item-text">Error loading scenarios</span></div>`;
    }
}

function searchWorldScenarios(query) {
    const cache = window._worldDropdownCache.scenarios;
    const fe = window.FilterEngine;
    const sortVal = document.querySelector('#scenario-dropdown .world-dropdown-sort')?.value || 'name:asc';
    const [field, dir] = sortVal.split(':');
    const items = fe ? fe.processItems(cache, { query, searchFields: ['name','title','description'], sortField: field, sortDirection: dir }) : cache;
    renderScenarioDropdown(items);
}

function toggleScenarioSelection(scenario, itemElement) {
    const ws    = window.worldState;
    const index = ws.selectedScenarios.findIndex(s => s.id === scenario.id);
    if (index > -1) {
        ws.selectedScenarios.splice(index, 1);
        itemElement.classList.remove('selected');
    } else {
        ws.selectedScenarios.push(scenario);
        itemElement.classList.add('selected');
    }
    updateSelectedScenariosDisplay();
    document.getElementById('scenario-dropdown')?.classList.remove('open');
}

function updateSelectedScenariosDisplay() {
    const container = document.getElementById('scenario-selected-items');
    if (!container) return;
    const items = window.worldState.selectedScenarios;
    container.innerHTML = items.map(s => `
        <div class="world-selected-item">
            <span class="world-selected-item-text">${s.title || s.name || 'Untitled Scenario'}</span>
            <button class="world-selected-item-remove" onclick="removeScenarioSelection('${s.id}')">×</button>
        </div>
    `).join('');
}

function removeScenarioSelection(scenarioId) {
    window.worldState.selectedScenarios = window.worldState.selectedScenarios.filter(s => s.id !== scenarioId);
    document.querySelectorAll('#scenario-dropdown-items .world-dropdown-item').forEach(item => {
        if (item.dataset.scenarioId === scenarioId) item.classList.remove('selected');
    });
    updateSelectedScenariosDisplay();
}

// ---------------------------------------------------------------------------
// Character Dropdown
// ---------------------------------------------------------------------------

function renderCharacterDropdown(items) {
    const dropdownItems = document.getElementById('character-dropdown-items');
    if (!dropdownItems) return;
    dropdownItems.innerHTML = '';

    const ws = window.worldState;
    if (items.length === 0) {
        dropdownItems.innerHTML = `<div class="world-dropdown-item disabled"><span class="world-dropdown-item-text">No characters found</span></div>`;
        return;
    }

    items.forEach((character) => {
        const displayName = character.name || 'Unnamed';
        const isSelected = ws.selectedCharacters.some(c => c.id === character.id);
        const item = document.createElement('div');
        item.className = 'world-dropdown-item' + (isSelected ? ' selected' : '');
        item.innerHTML = `<span class="world-dropdown-item-text">${displayName}</span>`;
        item.dataset.characterId = character.id;
        item.addEventListener('click', () => toggleCharacterSelection(character, item));
        dropdownItems.appendChild(item);
    });
}

async function loadAvailableCharacters() {
    const dropdownItems = document.getElementById('character-dropdown-items');
    if (!dropdownItems) return;

    try {
        const characters = JSON.parse(localStorage.getItem('howl-chat-characters') || '[]');
        if (characters.length === 0) {
            dropdownItems.innerHTML = `<div class="world-dropdown-item disabled"><span class="world-dropdown-item-text">No characters created yet</span></div>`;
            window._worldDropdownCache.characters = [];
            return;
        }

        window._worldDropdownCache.characters = characters;
        const fe = window.FilterEngine;
        const sorted = fe ? fe.sortItems(characters, 'name', 'asc') : characters;
        renderCharacterDropdown(sorted);

    } catch (error) {
        console.error('Failed to load characters:', error);
        dropdownItems.innerHTML = `<div class="world-dropdown-item disabled"><span class="world-dropdown-item-text">Error loading characters</span></div>`;
    }
}

function searchWorldCharacters(query) {
    const cache = window._worldDropdownCache.characters;
    const fe = window.FilterEngine;
    const sortVal = document.querySelector('#character-dropdown .world-dropdown-sort')?.value || 'name:asc';
    const [field, dir] = sortVal.split(':');
    const items = fe ? fe.processItems(cache, { query, searchFields: ['name','role','description'], sortField: field, sortDirection: dir }) : cache;
    renderCharacterDropdown(items);
}

function toggleCharacterSelection(character, itemElement) {
    const ws    = window.worldState;
    const index = ws.selectedCharacters.findIndex(c => c.id === character.id);
    if (index > -1) {
        ws.selectedCharacters.splice(index, 1);
        itemElement.classList.remove('selected');
    } else {
        ws.selectedCharacters.push(character);
        itemElement.classList.add('selected');
    }
    updateSelectedCharactersDisplay();
    document.getElementById('character-dropdown')?.classList.remove('open');
}

function updateSelectedCharactersDisplay() {
    const container = document.getElementById('character-selected-items');
    if (!container) return;
    const items = window.worldState.selectedCharacters;
    container.innerHTML = items.map(c => `
        <div class="world-selected-item">
            <span class="world-selected-item-text">${c.name || c.title || 'Untitled Character'}</span>
            <button class="world-selected-item-remove" onclick="removeCharacterSelection('${c.id}')">×</button>
        </div>
    `).join('');
}

function removeCharacterSelection(characterId) {
    window.worldState.selectedCharacters = window.worldState.selectedCharacters.filter(c => c.id !== characterId);
    document.querySelectorAll('#character-dropdown-items .world-dropdown-item').forEach(item => {
        if (item.dataset.characterId === characterId) item.classList.remove('selected');
    });
    updateSelectedCharactersDisplay();
}

// ---------------------------------------------------------------------------
// Lorebook Dropdown
// ---------------------------------------------------------------------------

function renderLorebookDropdown(items) {
    const dropdownItems = document.getElementById('world-lorebook-dropdown-items');
    if (!dropdownItems) return;
    dropdownItems.innerHTML = '';

    const ws = window.worldState;
    if (items.length === 0) {
        dropdownItems.innerHTML = `<div class="world-dropdown-item disabled"><span class="world-dropdown-item-text">No lorebooks found</span></div>`;
        return;
    }

    items.forEach((lorebook) => {
        if (!lorebook.id) lorebook.id = ws.lorebookKey(lorebook);
        const displayName = ws.lorebookDisplayName(lorebook);
        const isSelected = ws.selectedLorebooks.some(l => ws.lorebookKey(l) === lorebook.id);
        const item = document.createElement('div');
        item.className = 'world-dropdown-item' + (isSelected ? ' selected' : '');
        item.innerHTML = `<span class="world-dropdown-item-text">${displayName}</span>`;
        item.dataset.lorebookId = lorebook.id;
        item.addEventListener('click', () => toggleLorebookSelection(lorebook, item));
        dropdownItems.appendChild(item);
    });
}

async function loadAvailableLorebooks() {
    const dropdownItems = document.getElementById('world-lorebook-dropdown-items');
    if (!dropdownItems) return;

    try {
        const lorebooks = JSON.parse(localStorage.getItem('howl-chat-lorebooks') || '[]');
        if (lorebooks.length === 0) {
            dropdownItems.innerHTML = `<div class="world-dropdown-item disabled"><span class="world-dropdown-item-text">No lorebooks created yet</span></div>`;
            window._worldDropdownCache.lorebooks = [];
            return;
        }

        const ws = window.worldState;
        const processed = lorebooks.map(l => { if (!l.id) l.id = ws.lorebookKey(l); return l; });
        window._worldDropdownCache.lorebooks = processed;
        const fe = window.FilterEngine;
        const sorted = fe ? fe.sortItems(processed, 'name', 'asc') : processed;
        renderLorebookDropdown(sorted);

    } catch (error) {
        console.error('Failed to load lorebooks:', error);
        dropdownItems.innerHTML = `<div class="world-dropdown-item disabled"><span class="world-dropdown-item-text">Error loading lorebooks</span></div>`;
    }
}

function searchWorldLorebooks(query) {
    const cache = window._worldDropdownCache.lorebooks;
    const fe = window.FilterEngine;
    const sortVal = document.querySelector('#world-lorebook-dropdown .world-dropdown-sort')?.value || 'name:asc';
    const [field, dir] = sortVal.split(':');
    const items = fe ? fe.processItems(cache, { query, searchFields: ['name','title','description'], sortField: field, sortDirection: dir }) : cache;
    renderLorebookDropdown(items);
}

window._worldDropdownSort = (dropdownId, val) => {
    const [field, dir] = val.split(':');
    const fe = window.FilterEngine;
    if (dropdownId === 'scenario-dropdown') {
        const query = document.querySelector('#scenario-dropdown .world-dropdown-search-input')?.value || '';
        const items = fe ? fe.processItems(window._worldDropdownCache.scenarios, { query, searchFields: ['name','title','description'], sortField: field, sortDirection: dir }) : window._worldDropdownCache.scenarios;
        renderScenarioDropdown(items);
    } else if (dropdownId === 'character-dropdown') {
        const query = document.querySelector('#character-dropdown .world-dropdown-search-input')?.value || '';
        const items = fe ? fe.processItems(window._worldDropdownCache.characters, { query, searchFields: ['name','role','description'], sortField: field, sortDirection: dir }) : window._worldDropdownCache.characters;
        renderCharacterDropdown(items);
    } else if (dropdownId === 'world-lorebook-dropdown') {
        const query = document.querySelector('#world-lorebook-dropdown .world-dropdown-search-input')?.value || '';
        const items = fe ? fe.processItems(window._worldDropdownCache.lorebooks, { query, searchFields: ['name','title','description'], sortField: field, sortDirection: dir }) : window._worldDropdownCache.lorebooks;
        renderLorebookDropdown(items);
    }
};

function toggleLorebookSelection(lorebook, itemElement) {
    const ws  = window.worldState;
    const key = ws.lorebookKey(lorebook);
    const index = ws.selectedLorebooks.findIndex(l => ws.lorebookKey(l) === key);
    if (index > -1) {
        ws.selectedLorebooks.splice(index, 1);
        itemElement.classList.remove('selected');
    } else {
        if (!ws.selectedLorebooks.some(l => ws.lorebookKey(l) === key)) ws.selectedLorebooks.push(lorebook);
        itemElement.classList.add('selected');
    }
    updateSelectedLorebooksDisplay();
    document.getElementById('world-lorebook-dropdown')?.classList.remove('open');
}

function updateSelectedLorebooksDisplay() {
    const container = document.getElementById('world-lorebook-selected-items');
    if (!container) return;
    const ws = window.worldState;
    container.innerHTML = ws.selectedLorebooks.map(l => `
        <div class="world-selected-item">
            <span class="world-selected-item-text">${ws.lorebookDisplayName(l)}</span>
            <button class="world-selected-item-remove" onclick="removeLorebookSelection('${ws.lorebookKey(l)}')">×</button>
        </div>
    `).join('');
}

function removeLorebookSelection(lorebookId) {
    const ws = window.worldState;
    ws.selectedLorebooks = ws.selectedLorebooks.filter(l => ws.lorebookKey(l) !== lorebookId);
    document.querySelectorAll('#world-lorebook-dropdown-items .world-dropdown-item').forEach(item => {
        item.classList.toggle('selected', ws.selectedLorebooks.some(l => ws.lorebookKey(l) === item.dataset.lorebookId));
    });
    updateSelectedLorebooksDisplay();
}
