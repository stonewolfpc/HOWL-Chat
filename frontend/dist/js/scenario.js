/**
 * Scenario Page JavaScript
 * 
 * This file contains all scenario-related functionality including
 * scenario creation, character management, and scenario grid management.
 * 
 * @module scenario
 */

// Scenario Creation Popup Functions
function openScenarioCreationPopup() {
    const overlay = document.getElementById('scenario-creation-overlay');
    overlay.style.display = 'flex';
    loadAvailableCharacters();
}

function closeScenarioCreationPopup() {
    const overlay = document.getElementById('scenario-creation-overlay');
    overlay.style.display = 'none';
}

// Scenario Image Management
let scenarioImageData = null;

function handleScenarioImageSelect(event) {
    const file = event.target.files[0];
    
    if (file && file.type.startsWith('image/')) {
        const reader = new FileReader();
        reader.onload = function(e) {
            scenarioImageData = {
                name: file.name,
                data: e.target.result
            };
            updateScenarioImagePreview();
        };
        reader.readAsDataURL(file);
    }
    
    // Clear input to allow selecting same file again
    event.target.value = '';
}

function updateScenarioImagePreview() {
    const placeholder = document.getElementById('scenario-image-placeholder');
    const preview = document.getElementById('scenario-image-preview');
    const removeBtn = document.getElementById('scenario-image-remove');
    
    if (scenarioImageData) {
        placeholder.style.display = 'none';
        preview.src = scenarioImageData.data;
        preview.style.display = 'block';
        removeBtn.style.display = 'block';
    } else {
        placeholder.style.display = 'flex';
        preview.style.display = 'none';
        removeBtn.style.display = 'none';
    }
}

function removeScenarioImage(event) {
    event.stopPropagation();
    scenarioImageData = null;
    updateScenarioImagePreview();
}

// Character Dropdown Management
let availableCharacters = [];
let selectedCharacters = [];

function loadAvailableCharacters() {
    // Load characters from localStorage for now
    const characters = JSON.parse(localStorage.getItem('howl-chat-characters') || '[]');
    availableCharacters = characters;
    updateCharacterDropdown();
}

function updateCharacterDropdown() {
    const dropdownItems = document.getElementById('scenario-character-dropdown-items');
    
    if (availableCharacters.length === 0) {
        dropdownItems.innerHTML = '<div class="world-dropdown-item">No characters available</div>';
        return;
    }
    
    dropdownItems.innerHTML = availableCharacters
        .filter(character => !selectedCharacters.find(selected => selected.name === character.name))
        .map(character => `
            <div class="world-dropdown-item" onclick="addCharacterToScenario('${character.name}')">
                <div class="dropdown-item-content">
                    <div class="dropdown-item-title">${character.name}</div>
                    <div class="dropdown-item-subtitle">${character.metadata?.role || 'No role'}</div>
                </div>
            </div>
        `).join('');
}

function searchScenarioCharacters(searchTerm) {
    const filtered = availableCharacters.filter(character => 
        !selectedCharacters.find(selected => selected.name === character.name) &&
        (character.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
         (character.metadata?.role || '').toLowerCase().includes(searchTerm.toLowerCase()))
    );
    
    const dropdownItems = document.getElementById('scenario-character-dropdown-items');
    
    if (filtered.length === 0) {
        dropdownItems.innerHTML = '<div class="world-dropdown-item">No characters found</div>';
        return;
    }
    
    dropdownItems.innerHTML = filtered.map(character => `
        <div class="world-dropdown-item" onclick="addCharacterToScenario('${character.name}')">
            <div class="dropdown-item-content">
                <div class="dropdown-item-title">${character.name}</div>
                <div class="dropdown-item-subtitle">${character.metadata?.role || 'No role'}</div>
            </div>
        </div>
    `).join('');
}

function addCharacterToScenario(characterName) {
    const character = availableCharacters.find(c => c.name === characterName);
    if (character && !selectedCharacters.find(c => c.name === characterName)) {
        selectedCharacters.push(character);
        updateSelectedCharactersList();
        updateCharacterDropdown();
    }
}

function removeCharacterFromScenario(characterName) {
    selectedCharacters = selectedCharacters.filter(c => c.name !== characterName);
    updateSelectedCharactersList();
    updateCharacterDropdown();
}

function updateSelectedCharactersList() {
    const selectedItems = document.getElementById('scenario-character-selected-items');
    
    if (selectedCharacters.length === 0) {
        selectedItems.innerHTML = '';
        return;
    }
    
    selectedItems.innerHTML = selectedCharacters.map(character => `
        <div class="world-selected-item">
            <span class="selected-item-name">${character.name}</span>
            <button class="selected-item-remove" onclick="removeCharacterFromScenario('${character.name}')">×</button>
        </div>
    `).join('');
}

// Lorebooks Management
function loadAvailableLorebooks() {
    // For now, create some mock lorebooks
    // In future, this will load from backend
    const mockLorebooks = [
        { name: 'Main Lore', description: 'Core world knowledge' },
        { name: 'Magic System', description: 'Rules and mechanics of magic' },
        { name: 'Character Backstories', description: 'Detailed character histories' },
        { name: 'World History', description: 'Timeline of major events' }
    ];
    
    updateScenarioLorebookDropdown(mockLorebooks);
}

function updateScenarioLorebookDropdown(lorebooks) {
    const dropdownItems = document.getElementById('scenario-lorebook-dropdown-items');
    
    if (lorebooks.length === 0) {
        dropdownItems.innerHTML = '<div class="world-dropdown-item">No lorebooks available</div>';
        return;
    }
    
    dropdownItems.innerHTML = lorebooks
        .filter(lorebook => !selectedLorebooks.find(selected => selected.name === lorebook.name))
        .map(lorebook => `
            <div class="world-dropdown-item" onclick="addLorebookToScenario('${lorebook.name}')">
                <div class="dropdown-item-content">
                    <div class="dropdown-item-title">${lorebook.name}</div>
                    <div class="dropdown-item-subtitle">${lorebook.description}</div>
                </div>
            </div>
        `).join('');
}

function searchScenarioLorebooks(searchTerm) {
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
    
    const dropdownItems = document.getElementById('scenario-lorebook-dropdown-items');
    
    if (filtered.length === 0) {
        dropdownItems.innerHTML = '<div class="world-dropdown-item">No lorebooks found</div>';
        return;
    }
    
    dropdownItems.innerHTML = filtered.map(lorebook => `
        <div class="world-dropdown-item" onclick="addLorebookToScenario('${lorebook.name}')">
            <div class="dropdown-item-content">
                <div class="dropdown-item-title">${lorebook.name}</div>
                <div class="dropdown-item-subtitle">${lorebook.description}</div>
            </div>
        </div>
    `).join('');
}

function addLorebookToScenario(lorebookName) {
    const mockLorebooks = [
        { name: 'Main Lore', description: 'Core world knowledge' },
        { name: 'Magic System', description: 'Rules and mechanics of magic' },
        { name: 'Character Backstories', description: 'Detailed character histories' },
        { name: 'World History', description: 'Timeline of major events' }
    ];
    
    const lorebook = mockLorebooks.find(l => l.name === lorebookName);
    if (lorebook && !selectedLorebooks.find(l => l.name === lorebookName)) {
        selectedLorebooks.push(lorebook);
        updateSelectedLorebooksList();
        updateScenarioLorebookDropdown();
    }
}

function removeLorebookFromScenario(lorebookName) {
    selectedLorebooks = selectedLorebooks.filter(l => l.name !== lorebookName);
    updateSelectedLorebooksList();
    updateScenarioLorebookDropdown();
}

function updateSelectedLorebooksList() {
    const selectedItems = document.getElementById('scenario-lorebook-selected-items');
    
    if (selectedLorebooks.length === 0) {
        selectedItems.innerHTML = '';
        return;
    }
    
    selectedItems.innerHTML = selectedLorebooks.map(lorebook => `
        <div class="world-selected-item">
            <span class="selected-item-name">${lorebook.name}</span>
            <button class="selected-item-remove" onclick="removeLorebookFromScenario('${lorebook.name}')">×</button>
        </div>
    `).join('');
}

// Dropdown Toggle Function
function toggleDropdown(dropdownId) {
    const dropdown = document.getElementById(dropdownId);
    const allDropdowns = document.querySelectorAll('.world-dropdown-menu');
    
    // Close all other dropdowns
    allDropdowns.forEach(d => {
        if (d.id !== dropdownId) {
            d.classList.remove('show');
        }
    });
    
    // Toggle current dropdown
    dropdown.classList.toggle('show');
}

// Scenario Management Functions
function saveScenario() {
    const scenarioData = collectScenarioData();
    
    if (!scenarioData.title) {
        alert('Scenario title is required!');
        return;
    }
    
    // Save to localStorage for now (will be replaced with backend call)
    const scenarios = JSON.parse(localStorage.getItem('howl-chat-scenarios') || '[]');
    scenarios.push(scenarioData);
    localStorage.setItem('howl-chat-scenarios', JSON.stringify(scenarios));
    
    console.log('Scenario saved:', scenarioData);
    closeScenarioCreationPopup();
    
    // Refresh scenario grid
    loadScenarios();
}

function collectScenarioData() {
    return {
        // Basic scenario fields
        name: document.getElementById('scenario-title').value,
        description: document.getElementById('scenario-details').value,
        image: scenarioImageData,
        
        // Character attachments
        characters: selectedCharacters.map(character => ({
            name: character.name,
            role: character.metadata?.role || '',
            avatar: character.avatar || ''
        })),
        
        // Metadata
        metadata: {
            version: 1,
            created: new Date().toISOString(),
            updated: new Date().toISOString(),
            seed: Math.floor(Math.random() * 1000000000)
        },
        
        // HOWL specific fields
        howl: {
            world: '',
            lorebooks: selectedLorebooks.map(lorebook => ({
                name: lorebook.name,
                description: lorebook.description
            })),
            notes: '',
            image_style: 'fantasy-realistic'
        }
    };
}

// Scenario Grid Management
function loadScenarios() {
    const grid = document.getElementById('scenario-grid');
    
    // Load scenarios from localStorage for now
    const scenarios = JSON.parse(localStorage.getItem('howl-chat-scenarios') || '[]');
    
    if (scenarios.length === 0) {
        grid.innerHTML = `
            <div class="world-grid-empty">
                <div class="world-grid-empty-icon">📜</div>
                <div class="world-grid-empty-text">No scenarios created yet</div>
                <div class="world-grid-empty-subtext">Click "Create New Scenario" to get started</div>
            </div>
        `;
        return;
    }
    
    grid.innerHTML = scenarios.map((scenario, index) => {
        const imageSrc = scenario.image ? scenario.image.data : 'assets/scroll.png';
        
        return `
            <div class="world-bubble" onclick="editScenario(${index})">
                <div class="world-image">
                    <img src="${imageSrc}" alt="${scenario.name}" class="world-bubble-img">
                </div>
                <div class="world-info">
                    <div class="world-name">${scenario.name}</div>
                    <div class="world-meta">
                        <span class="world-character-count">${scenario.characters.length} characters</span>
                    </div>
                </div>
                <div class="world-actions">
                    <button class="world-edit-btn" onclick="editScenario(${index}, event)">✏️</button>
                    <button class="world-delete-btn" onclick="deleteScenario(${index}, event)">🗑️</button>
                </div>
            </div>
        `;
    }).join('');
}

function editScenario(index, event) {
    if (event) event.stopPropagation();
    
    const scenarios = JSON.parse(localStorage.getItem('howl-chat-scenarios') || '[]');
    const scenario = scenarios[index];
    
    // Load scenario data into form
    document.getElementById('scenario-title').value = scenario.name || '';
    document.getElementById('scenario-details').value = scenario.description || '';
    
    // Load image
    scenarioImageData = scenario.image || null;
    updateScenarioImagePreview();
    
    // Load characters
    selectedCharacters = scenario.characters || [];
    updateSelectedCharactersList();
    loadAvailableCharacters();
    
    // Load lorebooks
    selectedLorebooks = scenario.howl?.lorebooks || [];
    updateSelectedLorebooksList();
    loadAvailableLorebooks();
    
    // Open popup for editing
    openScenarioCreationPopup();
}

function deleteScenario(index, event) {
    if (event) event.stopPropagation();
    
    if (confirm('Are you sure you want to delete this scenario?')) {
        const scenarios = JSON.parse(localStorage.getItem('howl-chat-scenarios') || '[]');
        scenarios.splice(index, 1);
        localStorage.setItem('howl-chat-scenarios', JSON.stringify(scenarios));
        loadScenarios();
    }
}

// Initialize the scenario page when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    loadScenarios();
    loadAvailableLorebooks();
    
    // Add keyboard shortcuts
    document.addEventListener('keydown', function(event) {
        // Escape key closes popup
        if (event.key === 'Escape') {
            closeScenarioCreationPopup();
        }
        
        // Ctrl+N for new scenario
        if (event.ctrlKey && event.key === 'n') {
            event.preventDefault();
            openScenarioCreationPopup();
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
    document.getElementById('scenario-creation-overlay').addEventListener('click', function(event) {
        if (event.target === this) {
            closeScenarioCreationPopup();
        }
    });
});

// Export functions for potential use by other modules
window.scenarioModule = {
    openScenarioCreationPopup,
    closeScenarioCreationPopup,
    saveScenario,
    loadScenarios,
    loadAvailableCharacters,
    addCharacterToScenario,
    removeCharacterFromScenario
};
