/**
 * Character Lorebooks
 *
 * Lorebook dropdown loading, toggle selection, chip display,
 * and the generic toggleDropdown function for the character page.
 *
 * Depends on: character-state.js
 *
 * @module character-lorebooks
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Dropdown Toggle
// ---------------------------------------------------------------------------

function toggleDropdown(dropdownId) {
    const container = document.getElementById(dropdownId);
    if (!container) { console.error(`Dropdown container not found: ${dropdownId}`); return; }

    const menu = container.querySelector('.world-dropdown-menu');
    if (!menu) { console.error(`Dropdown menu not found in: ${dropdownId}`); return; }

    document.querySelectorAll('.world-dropdown-menu').forEach(d => {
        if (d.parentElement.id !== dropdownId) d.classList.remove('show');
    });

    const isNowShowing = !menu.classList.contains('show');
    menu.classList.toggle('show');
    container.classList.toggle('open', menu.classList.contains('show'));

    if (dropdownId === 'character-lorebook-dropdown' && isNowShowing) loadAvailableLorebooks();
}

// ---------------------------------------------------------------------------
// Load Dropdown
// ---------------------------------------------------------------------------

async function loadAvailableLorebooks() {
    const dropdownItems = document.getElementById('character-lorebook-dropdown-items');
    if (!dropdownItems) return;

    try {
        const cs       = window.characterState;
        const lorebooks = JSON.parse(localStorage.getItem('howl-chat-lorebooks') || '[]');
        cs.availableLorebooks = lorebooks;
        const fe = window.FilterEngine;
        const sorted = fe ? fe.sortItems(lorebooks, 'name', 'asc') : lorebooks;
        renderCharacterLorebookDropdown(sorted);

    } catch (error) {
        console.error('Failed to load lorebooks for character:', error);
        dropdownItems.innerHTML = `<div class="world-dropdown-item disabled"><span class="world-dropdown-item-text">Error loading lorebooks</span></div>`;
    }
}

function renderCharacterLorebookDropdown(items) {
    const dropdownItems = document.getElementById('character-lorebook-dropdown-items');
    if (!dropdownItems) return;
    dropdownItems.innerHTML = '';

    const cs = window.characterState;
    if (items.length === 0) {
        dropdownItems.innerHTML = `<div class="world-dropdown-item disabled"><span class="world-dropdown-item-text">No lorebooks found</span></div>`;
        return;
    }

    items.forEach(lorebook => {
        if (!lorebook.id) lorebook.id = cs.lorebookKey(lorebook);
        const displayName = cs.lorebookDisplayName(lorebook);
        const isSelected  = cs.selectedLorebooks.some(l => cs.lorebookKey(l) === cs.lorebookKey(lorebook));

        const item = document.createElement('div');
        item.className = 'world-dropdown-item';
        if (isSelected) item.classList.add('selected');
        item.innerHTML = `<span class="world-dropdown-item-text">${displayName}</span>`;
        item.dataset.lorebookId = lorebook.id;
        item.addEventListener('click', () => toggleLorebookSelection(lorebook, item));
        dropdownItems.appendChild(item);
    });
}

function searchCharacterLorebooks(query) {
    const cs = window.characterState;
    const fe = window.FilterEngine;
    const sortVal = document.querySelector('#character-lorebook-dropdown .world-dropdown-sort')?.value || 'name:asc';
    const [field, dir] = sortVal.split(':');
    const items = fe ? fe.processItems(cs.availableLorebooks, { query, searchFields: ['name','title','description'], sortField: field, sortDirection: dir }) : cs.availableLorebooks;
    renderCharacterLorebookDropdown(items);
}

window._characterLorebookSort = (val) => {
    const [field, dir] = val.split(':');
    const query = document.querySelector('#character-lorebook-dropdown .world-dropdown-search-input')?.value || '';
    const cs = window.characterState;
    const fe = window.FilterEngine;
    const items = fe ? fe.processItems(cs.availableLorebooks, { query, searchFields: ['name','title','description'], sortField: field, sortDirection: dir }) : cs.availableLorebooks;
    renderCharacterLorebookDropdown(items);
};

// ---------------------------------------------------------------------------
// Toggle & Remove
// ---------------------------------------------------------------------------

function toggleLorebookSelection(lorebook, itemElement) {
    const cs  = window.characterState;
    const key = cs.lorebookKey(lorebook);
    const idx = cs.selectedLorebooks.findIndex(l => cs.lorebookKey(l) === key);

    if (idx > -1) {
        cs.selectedLorebooks.splice(idx, 1);
        itemElement.classList.remove('selected');
    } else {
        if (!cs.selectedLorebooks.some(l => cs.lorebookKey(l) === key)) cs.selectedLorebooks.push(lorebook);
        itemElement.classList.add('selected');
    }

    updateSelectedLorebooksList();
    const dropdown = document.getElementById('character-lorebook-dropdown');
    const menu     = dropdown?.querySelector('.world-dropdown-menu');
    if (menu) { menu.classList.remove('show'); dropdown.classList.remove('open'); }
}

function removeLorebookFromCharacter(lorebookId) {
    const cs  = window.characterState;
    const idx = cs.selectedLorebooks.findIndex(l => cs.lorebookKey(l) === lorebookId);
    if (idx > -1) {
        cs.selectedLorebooks.splice(idx, 1);
        updateSelectedLorebooksList();
        loadAvailableLorebooks();
    }
}

function updateSelectedLorebooksList() {
    const container = document.getElementById('character-lorebook-selected-items');
    if (!container) return;
    const cs = window.characterState;
    container.innerHTML = cs.selectedLorebooks.length === 0
        ? ''
        : cs.selectedLorebooks.map(l => `
            <div class="world-selected-item">
                <span class="selected-item-name">${cs.lorebookDisplayName(l)}</span>
                <button type="button" class="world-selected-item-remove" onclick="removeLorebookFromCharacter('${cs.lorebookKey(l)}')">×</button>
            </div>
        `).join('');
}

// ---------------------------------------------------------------------------
// Cross-page Compat
// ---------------------------------------------------------------------------

window.characterLoadAvailableLorebooks = loadAvailableLorebooks;
