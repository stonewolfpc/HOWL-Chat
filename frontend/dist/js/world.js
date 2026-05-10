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
            updateSelectedItemsDisplay('scenario');
            updateSelectedItemsDisplay('character');
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
        updatedAt: new Date().toISOString()
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
        const newWorld = {
            ...world,
            id: Date.now().toString(),
            title: `${world.title} (Copy)`,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString()
        };
        worlds.push(newWorld);
        renderWorldGrid();
        console.log('World copied:', newWorld);
    }
}

// Delete world
function deleteWorld(worldId) {
    if (confirm('Are you sure you want to delete this world?')) {
        worlds = worlds.filter(w => w.id !== worldId);
        renderWorldGrid();
        console.log('World deleted:', worldId);
    }
}

// Close dropdowns when clicking outside
document.addEventListener('click', function(event) {
    if (!event.target.closest('.world-dropdown-container')) {
        document.querySelectorAll('.world-dropdown-container').forEach(d => {
            d.classList.remove('open');
        });
    }
});

// Initialize page
console.log('World page initialized');
