/**
 * World Page Module
 * 
 * This module handles all World page interactions including:
 * - World creation popup management
 * - Image selection and preview
 * - Scenario and character dropdown management
 * - World bubble actions (edit, copy, delete)
 * - World saving and hotswap functionality
 * 
 * @module world
 */

// State management
let worlds = [];
let selectedScenarios = [];
let selectedCharacters = [];
let selectedLorebooks = [];
let editingWorldId = null;

// Open world creation popup
function openWorldCreationPopup(worldId = null) {
    const overlay = document.getElementById('world-creation-overlay');
    overlay.classList.add('show');
    
    if (worldId) {
        // Edit mode - populate form
        editingWorldId = worldId;
        const world = worlds.find(w => w.id === worldId);
        if (world) {
            document.getElementById('world-title').value = world.title;
            document.getElementById('world-details').value = world.details;
            
            if (world.image) {
                document.getElementById('world-image-preview').src = world.image;
                document.getElementById('world-image-preview').style.display = 'block';
                document.getElementById('world-image-placeholder').style.display = 'none';
                document.getElementById('world-image-remove').style.display = 'block';
            }
            
            selectedScenarios = world.scenarios || [];
            selectedCharacters = world.characters || [];
            selectedLorebooks = world.howl?.lorebooks || [];
            updateSelectedItemsDisplay('scenario');
            updateSelectedItemsDisplay('character');
            updateWorldSelectedLorebooksList();
        }
    } else {
        // Create mode - clear form
        editingWorldId = null;
        clearWorldForm();
    }
}

// Close world creation popup
function closeWorldCreationPopup() {
    const overlay = document.getElementById('world-creation-overlay');
    overlay.classList.remove('show');
    clearWorldForm();
}

// Clear world creation form
function clearWorldForm() {
    document.getElementById('world-title').value = '';
    document.getElementById('world-details').value = '';
    document.getElementById('world-image-preview').style.display = 'none';
    document.getElementById('world-image-placeholder').style.display = 'flex';
    document.getElementById('world-image-remove').style.display = 'none';
    document.getElementById('world-image-input').value = '';
    selectedScenarios = [];
    selectedCharacters = [];
    updateSelectedItemsDisplay('scenario');
    updateSelectedItemsDisplay('character');
}

// Handle world image selection
function handleWorldImageSelect(event) {
    const file = event.target.files[0];
    if (file) {
        const reader = new FileReader();
        reader.onload = function(e) {
            const preview = document.getElementById('world-image-preview');
            preview.src = e.target.result;
            preview.style.display = 'block';
            document.getElementById('world-image-placeholder').style.display = 'none';
            document.getElementById('world-image-remove').style.display = 'block';
        };
        reader.readAsDataURL(file);
    }
}

// Remove world image
function removeWorldImage(event) {
    event.stopPropagation();
    document.getElementById('world-image-preview').style.display = 'none';
    document.getElementById('world-image-placeholder').style.display = 'flex';
    document.getElementById('world-image-remove').style.display = 'none';
    document.getElementById('world-image-input').value = '';
}

// Toggle dropdown
function toggleDropdown(dropdownId) {
    const dropdown = document.getElementById(dropdownId);
    dropdown.classList.toggle('open');
    
    // Close other dropdowns
    document.querySelectorAll('.world-dropdown-container').forEach(d => {
        if (d.id !== dropdownId) {
            d.classList.remove('open');
        }
    });
}

// Update selected items display
function updateSelectedItemsDisplay(type) {
    const container = document.getElementById(`${type}-selected-items`);
    const items = type === 'scenario' ? selectedScenarios : selectedCharacters;
    
    container.innerHTML = items.map(item => `
        <div class="world-selected-item">
            <span>${item.name}</span>
            <button class="world-selected-item-remove" onclick="removeSelectedItem('${type}', '${item.id}')">×</button>
        </div>
    `).join('');
    
    // Update dropdown trigger text
    const triggerText = document.getElementById(`${type}-dropdown-text`);
    if (items.length > 0) {
        triggerText.textContent = `${items.length} ${type}${items.length === 1 ? '' : 's'} selected`;
    } else {
        triggerText.textContent = `Select ${type}s...`;
    }
}

// Remove selected item
function removeSelectedItem(type, itemId) {
    if (type === 'scenario') {
        selectedScenarios = selectedScenarios.filter(s => s.id !== itemId);
    } else {
        selectedCharacters = selectedCharacters.filter(c => c.id !== itemId);
    }
    updateSelectedItemsDisplay(type);
}

// Save world
function saveWorld() {
    const title = document.getElementById('world-title').value.trim();
    const details = document.getElementById('world-details').value.trim();
    const imagePreview = document.getElementById('world-image-preview');
    
    if (!title) {
        alert('Please enter a world title');
        return;
    }
    
    const worldData = {
        id: editingWorldId || Date.now().toString(),
        title: title,
        details: details,
        image: imagePreview.style.display === 'block' ? imagePreview.src : null,
        scenarios: selectedScenarios,
        characters: selectedCharacters,
        createdAt: editingWorldId ? worlds.find(w => w.id === editingWorldId)?.createdAt : new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        howl: {
            lorebooks: selectedLorebooks.map(lorebook => ({
                name: lorebook.name,
                description: lorebook.description
            }))
        }
    };
    
    if (editingWorldId) {
        // Update existing world
        const index = worlds.findIndex(w => w.id === editingWorldId);
        if (index !== -1) {
            worlds[index] = worldData;
        }
    } else {
        // Create new world
        worlds.push(worldData);
    }
    
    renderWorldGrid();
    closeWorldCreationPopup();
    
    // Auto-save to Documents\HOWL_Chat\Worlds (backend integration needed)
    console.log('World saved:', worldData);
}

// Render world grid
function renderWorldGrid() {
    const grid = document.getElementById('world-grid');
    
    if (worlds.length === 0) {
        grid.innerHTML = `
            <div class="world-grid-empty">
                <div class="world-grid-empty-icon">🌍</div>
                <div class="world-grid-empty-text">No worlds created yet</div>
                <div class="world-grid-empty-subtext">Click "Create New World" to get started</div>
            </div>
        `;
        return;
    }
    
    grid.innerHTML = worlds.map(world => `
        <div class="world-bubble" data-world-id="${world.id}">
            ${world.image ? 
                `<img class="world-bubble-image" src="${world.image}" alt="${world.title}">` :
                `<div class="world-bubble-placeholder">
                    <div class="world-bubble-placeholder-icon">🌍</div>
                </div>`
            }
            <div class="world-bubble-name">${world.title}</div>
            <div class="world-bubble-actions">
                <button class="world-bubble-action" onclick="editWorld('${world.id}')" title="Edit">✏️</button>
                <button class="world-bubble-action" onclick="copyWorld('${world.id}')" title="Copy">📋</button>
                <button class="world-bubble-action" onclick="deleteWorld('${world.id}')" title="Delete">🗑️</button>
            </div>
        </div>
    `).join('');
}

// Edit world
function editWorld(worldId) {
    openWorldCreationPopup(worldId);
}

// Copy world
function copyWorld(worldId) {
    const world = worlds.find(w => w.id === worldId);
    if (world) {
    }
    
    dropdownItems.innerHTML = lorebooks
        .filter(lorebook => !selectedLorebooks.find(selected => selected.name === lorebook.name))
        .map(lorebook => `
            <div class="world-dropdown-item" onclick="addLorebookToWorld('${lorebook.name}')">
                <div class="dropdown-item-content">
                    <div class="dropdown-item-title">${lorebook.name}</div>
                    <div class="dropdown-item-subtitle">${lorebook.description}</div>
                </div>
            </div>
        `).join('');
}

function searchWorldLorebooks(searchTerm) {
    const mockLorebooks = [
        { name: 'Main Lore', description: 'Core world knowledge' },
        { name: 'Magic System', description: 'Rules and mechanics of magic' },
        { name: 'Character Backstories', description: 'Detailed character histories' },
        { name: 'World History', description: 'Timeline of major events' }
    ];
    
    const filtered = mockLorebooks.filter(lorebook => 
        !selectedLorebooks.find(selected => selected.name === lorebook.name) &&
        (lorebook.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
         lorebook.description.toLowerCase().includes(searchTerm.toLowerCase()))
    );
    
    const dropdownItems = document.getElementById('world-lorebook-dropdown-items');
    
    if (filtered.length === 0) {
        dropdownItems.innerHTML = '<div class="world-dropdown-item">No lorebooks found</div>';
        return;
    }
    
    dropdownItems.innerHTML = filtered.map(lorebook => `
        <div class="world-dropdown-item" onclick="addLorebookToWorld('${lorebook.name}')">
            <div class="dropdown-item-content">
                <div class="dropdown-item-title">${lorebook.name}</div>
                <div class="dropdown-item-subtitle">${lorebook.description}</div>
            </div>
        </div>
    `).join('');
}

function addLorebookToWorld(lorebookName) {
    const mockLorebooks = [
        { name: 'Main Lore', description: 'Core world knowledge' },
        { name: 'Magic System', description: 'Rules and mechanics of magic' },
        { name: 'Character Backstories', description: 'Detailed character histories' },
        { name: 'World History', description: 'Timeline of major events' }
    ];
    
    const lorebook = mockLorebooks.find(l => l.name === lorebookName);
    if (lorebook && !selectedLorebooks.find(l => l.name === lorebookName)) {
        selectedLorebooks.push(lorebook);
        updateWorldSelectedLorebooksList();
        updateWorldLorebookDropdown();
    }
}

function removeLorebookFromWorld(lorebookName) {
    selectedLorebooks = selectedLorebooks.filter(l => l.name !== lorebookName);
    updateWorldSelectedLorebooksList();
    updateWorldLorebookDropdown();
}

function updateWorldSelectedLorebooksList() {
    const selectedItems = document.getElementById('world-lorebook-selected-items');
    
    if (selectedLorebooks.length === 0) {
        selectedItems.innerHTML = '';
        return;
    }
    
    selectedItems.innerHTML = selectedLorebooks.map(lorebook => `
        <div class="world-selected-item">
            <span class="selected-item-name">${lorebook.name}</span>
            <button class="selected-item-remove" onclick="removeLorebookFromWorld('${lorebook.name}')">×</button>
        </div>
    `).join('');
}

// Initialize world page when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    loadWorlds();
    loadAvailableScenarios();
    loadAvailableCharacters();
    loadAvailableLorebooks();
    
    // Add keyboard shortcuts
    document.addEventListener('keydown', function(event) {
        // Escape key closes popup
        if (event.key === 'Escape') {
            closeWorldCreationPopup();
        }
        
        // Ctrl+N for new world
        if (event.ctrlKey && event.key === 'n') {
            event.preventDefault();
            openWorldCreationPopup();
        }
    });
    
    // Close dropdowns when clicking outside
    document.addEventListener('click', function(event) {
        if (!event.target.closest('.world-dropdown-container')) {
            document.querySelectorAll('.world-dropdown-menu').forEach(dropdown => {
                dropdown.classList.remove('show');
            });
        }
    });
    
    // Close popup when clicking outside
    document.getElementById('world-creation-overlay').addEventListener('click', function(event) {
        if (event.target === this) {
            closeWorldCreationPopup();
        }
    });
});
