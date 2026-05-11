/**
 * Character Page JavaScript
 * 
 * This file contains all character-related functionality including
 * character creation, user details management, and character grid management.
 * 
 * @module character
 */

// Character Creation Popup Functions
function openCharacterCreationPopup() {
    const overlay = document.getElementById('character-creation-overlay');
    overlay.style.display = 'flex';
    loadAvailableCharacters();
    loadAvailableLorebooks();
}

function closeCharacterCreationPopup() {
    const overlay = document.getElementById('character-creation-overlay');
    overlay.style.display = 'none';
}

// User Details Popup Functions
function openUserDetailsPopup() {
    const overlay = document.getElementById('user-details-overlay');
    overlay.style.display = 'flex';
    // Load existing user details if any
    loadUserDetails();
}

function closeUserDetailsPopup() {
    const overlay = document.getElementById('user-details-overlay');
    overlay.style.display = 'none';
}

// Character Image Management
let characterImages = [];
let thumbnailIndex = 0;

function handleCharacterImageSelect(event) {
    const files = Array.from(event.target.files);
    
    files.forEach(file => {
        if (file.type.startsWith('image/')) {
            const reader = new FileReader();
            reader.onload = function(e) {
                const imageData = {
                    id: Date.now() + Math.random(),
                    name: file.name,
                    size: formatFileSize(file.size),
                    data: e.target.result,
                    isThumbnail: characterImages.length === 0 // First image is automatically thumbnail
                };
                characterImages.push(imageData);
                updateImagesList();
            };
            reader.readAsDataURL(file);
        }
    });
    
    // Clear the input to allow selecting the same file again
    event.target.value = '';
}

function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

function updateImagesList() {
    const imagesList = document.getElementById('character-images-list');
    
    if (characterImages.length === 0) {
        imagesList.innerHTML = `
            <div class="no-images-placeholder">
                <span class="no-images-icon">🖼️</span>
                <span class="no-images-text">No images added yet</span>
            </div>
        `;
        return;
    }
    
    imagesList.innerHTML = characterImages.map((image, index) => `
        <div class="character-image-item ${image.isThumbnail ? 'thumbnail' : ''}" onclick="setThumbnail(${index})">
            <img src="${image.data}" alt="${image.name}" class="character-image-thumbnail">
            <div class="character-image-info">
                <div class="character-image-name">${image.name}</div>
                <div class="character-image-size">${image.size}</div>
            </div>
            <div class="character-image-actions">
                ${image.isThumbnail ? '<button class="thumbnail-indicator">Thumbnail</button>' : ''}
                <button class="remove-image-button" onclick="removeImage(${index}, event)">Remove</button>
            </div>
        </div>
    `).join('');
}

function setThumbnail(index) {
    characterImages.forEach((image, i) => {
        image.isThumbnail = i === index;
    });
    thumbnailIndex = index;
    updateImagesList();
}

function removeImage(index, event) {
    event.stopPropagation();
    characterImages.splice(index, 1);
    
    // If we removed the thumbnail, make the first image the new thumbnail
    if (characterImages.length > 0 && thumbnailIndex >= characterImages.length) {
        thumbnailIndex = 0;
        characterImages[0].isThumbnail = true;
    }
    
    updateImagesList();
}

// Character Management Functions
function saveCharacter() {
    const characterData = collectCharacterData();
    
    if (!characterData.name) {
        alert('Character name is required!');
        return;
    }
    
    // Save to cache and localStorage
    const characters = characterCache || JSON.parse(localStorage.getItem('howl-chat-characters') || '[]');
    characters.push(characterData);
    characterCache = characters;
    localStorage.setItem('howl-chat-characters', JSON.stringify(characters));
    
    console.log('Character saved:', characterData);
    closeCharacterCreationPopup();
    
    // Refresh the character grid
    loadCharacters();
}

// Lorebook Management
let availableLorebooks = [];
let selectedLorebooks = [];
let currentEditingId = null;
let characterCache = null;

function loadAvailableLorebooks() {
    // For now, create some mock lorebooks
    // In future, this will load from backend
    availableLorebooks = [
        { name: 'Main Lore', description: 'Core world knowledge' },
        { name: 'Magic System', description: 'Rules and mechanics of magic' },
        { name: 'Character Backstories', description: 'Detailed character histories' },
        { name: 'World History', description: 'Timeline of major events' }
    ];
    updateLorebookDropdown();
}

function updateLorebookDropdown() {
    const dropdownItems = document.getElementById('character-lorebook-dropdown-items');
    
    if (availableLorebooks.length === 0) {
        dropdownItems.innerHTML = '<div class="world-dropdown-item">No lorebooks available</div>';
        return;
    }
    
    dropdownItems.innerHTML = availableLorebooks
        .filter(lorebook => !selectedLorebooks.find(selected => selected.name === lorebook.name))
        .map(lorebook => `
            <div class="world-dropdown-item" onclick="addLorebookToCharacter('${lorebook.name}')">
                <div class="dropdown-item-content">
                    <div class="dropdown-item-title">${lorebook.name}</div>
                    <div class="dropdown-item-subtitle">${lorebook.description}</div>
                </div>
            </div>
        `).join('');
}

function searchCharacterLorebooks(searchTerm) {
    const filtered = availableLorebooks.filter(lorebook => 
        !selectedLorebooks.find(selected => selected.name === lorebook.name) &&
        (lorebook.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
         lorebook.description.toLowerCase().includes(searchTerm.toLowerCase()))
    );
    
    const dropdownItems = document.getElementById('character-lorebook-dropdown-items');
    
    if (filtered.length === 0) {
        dropdownItems.innerHTML = '<div class="world-dropdown-item">No lorebooks found</div>';
        return;
    }
    
    dropdownItems.innerHTML = filtered.map(lorebook => `
        <div class="world-dropdown-item" onclick="addLorebookToCharacter('${lorebook.name}')">
            <div class="dropdown-item-content">
                <div class="dropdown-item-title">${lorebook.name}</div>
                <div class="dropdown-item-subtitle">${lorebook.description}</div>
            </div>
        </div>
    `).join('');
}

function addLorebookToCharacter(lorebookName) {
    const lorebook = availableLorebooks.find(l => l.name === lorebookName);
    if (lorebook && !selectedLorebooks.find(l => l.name === lorebookName)) {
        selectedLorebooks.push(lorebook);
        updateSelectedLorebooksList();
        updateLorebookDropdown();
    }
}

function removeLorebookFromCharacter(lorebookName) {
    selectedLorebooks = selectedLorebooks.filter(l => l.name !== lorebookName);
    updateSelectedLorebooksList();
    updateLorebookDropdown();
}

function updateSelectedLorebooksList() {
    const selectedItems = document.getElementById('character-lorebook-selected-items');
    
    if (selectedLorebooks.length === 0) {
        selectedItems.innerHTML = '';
        return;
    }
    
    selectedItems.innerHTML = selectedLorebooks.map(lorebook => `
        <div class="world-selected-item">
            <span class="selected-item-name">${lorebook.name}</span>
            <button class="selected-item-remove" onclick="removeLorebookFromCharacter('${lorebook.name}')">×</button>
        </div>
    `).join('');
}

function collectCharacterData() {
    const thumbnail = characterImages.find(img => img.isThumbnail);
    
    return {
        // SillyTavern compatible fields
        name: document.getElementById('character-name').value,
        description: document.getElementById('character-description').value,
        personality: document.getElementById('character-personality').value,
        scenario: document.getElementById('character-intro').value,
        first_mes: document.getElementById('character-intro').value,
        mes_example: document.getElementById('character-dialogue-examples').value,
        avatar: thumbnail ? thumbnail.data : '',
        tags: [],
        
        // Metadata
        metadata: {
            version: 1,
            created: new Date().toISOString(),
            updated: new Date().toISOString(),
            seed: Math.floor(Math.random() * 1000000000),
            age: document.getElementById('character-age').value,
            role: document.getElementById('character-role').value,
            background: document.getElementById('character-background').value
        },
        
        // HOWL specific fields
        howl: {
            world: '',
            lorebooks: selectedLorebooks.map(lorebook => ({
                name: lorebook.name,
                description: lorebook.description
            })),
            memories: [],
            notes: '',
            image_style: 'fantasy-realistic',
            images: characterImages.map(img => ({
                id: img.id,
                name: img.name,
                size: img.size,
                data: img.data,
                isThumbnail: img.isThumbnail
            }))
        }
    };
}

function saveUserDetails() {
    const userName = document.getElementById('user-name').value;
    const userDescription = document.getElementById('user-description').value;
    const userNotes = document.getElementById('user-notes').value;
    
    // Save user details to localStorage for now
    const userDetails = {
        name: userName,
        description: userDescription,
        notes: userNotes
    };
    
    localStorage.setItem('howl-chat-user-details', JSON.stringify(userDetails));
    
    console.log('User details saved:', userDetails);
    closeUserDetailsPopup();
}

function loadUserDetails() {
    // Load existing user details from localStorage
    const savedDetails = localStorage.getItem('howl-chat-user-details');
    
    if (savedDetails) {
        const userDetails = JSON.parse(savedDetails);
        document.getElementById('user-name').value = userDetails.name || '';
        document.getElementById('user-description').value = userDetails.description || '';
        document.getElementById('user-notes').value = userDetails.notes || '';
    }
}

// Character Grid Management
function loadCharacters() {
    const grid = document.getElementById('character-grid');
    
    // Load characters from cache or localStorage for now
    const characters = characterCache || JSON.parse(localStorage.getItem('howl-chat-characters') || '[]');
    characterCache = characters; // Update cache
    
    if (characters.length === 0) {
        grid.innerHTML = `
            <div class="character-grid-empty">
                <div class="character-grid-empty-icon">👥</div>
                <div class="character-grid-empty-text">No characters created yet</div>
                <div class="character-grid-empty-subtext">Click "Create New" to get started</div>
            </div>
        `;
        return;
    }
    
    grid.innerHTML = characters.map((character, index) => {
        const thumbnail = character.howl?.images?.find(img => img.isThumbnail);
        const avatarSrc = thumbnail ? thumbnail.data : (character.avatar || 'assets/Character.png');
        
        return `
            <div class="character-bubble" onclick="editCharacter(${index})">
                <div class="character-avatar">
                    <img src="${avatarSrc}" alt="${character.name}" class="character-avatar-img">
                </div>
                <div class="character-info">
                    <div class="character-name">${character.name}</div>
                    <div class="character-meta">
                        <span class="character-role">${character.metadata?.role || 'No role'}</span>
                    </div>
                </div>
                <div class="character-actions">
                    <button class="character-edit-btn" onclick="editCharacter(${index}, event)">✏️</button>
                    <button class="character-delete-btn" onclick="deleteCharacter(${index}, event)">🗑️</button>
                </div>
            </div>
        `;
    }).join('');
}

function editCharacter(index, event) {
    if (event) event.stopPropagation();
    
    // Load characters from cache to avoid repeated localStorage calls
    const characters = window.characterCache || JSON.parse(localStorage.getItem('howl-chat-characters') || '[]');
    const character = characters[index];
    
    if (!character) return;
    
    // Load character data into form
    document.getElementById('character-name').value = character.name || '';
    document.getElementById('character-age').value = character.metadata?.age || '';
    document.getElementById('character-role').value = character.metadata?.role || '';
    document.getElementById('character-description').value = character.description || '';
    document.getElementById('character-personality').value = character.personality || '';
    document.getElementById('character-background').value = character.metadata?.background || '';
    document.getElementById('character-intro').value = character.first_mes || '';
    document.getElementById('character-dialogue-examples').value = character.mes_example || '';
    
    // Load images
    characterImages = character.howl?.images || [];
    updateImagesList();
    
    // Load lorebooks
    selectedLorebooks = character.howl?.lorebooks || [];
    updateSelectedLorebooksList();
    loadAvailableLorebooks();
    
    // Open popup for editing
    openCharacterCreationPopup();
}

function deleteCharacter(index, event) {
    if (event) event.stopPropagation();
    
    if (confirm('Are you sure you want to delete this character?')) {
        const characters = characterCache || JSON.parse(localStorage.getItem('howl-chat-characters') || '[]');
        characters.splice(index, 1);
        characterCache = characters;
        localStorage.setItem('howl-chat-characters', JSON.stringify(characters));
        loadCharacters();
    }
}

// Initialize character page when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    loadCharacters();
    loadAvailableLorebooks();
    
    // Add keyboard shortcuts
    document.addEventListener('keydown', function(event) {
        // Escape key closes popups
        if (event.key === 'Escape') {
            closeCharacterCreationPopup();
            closeUserDetailsPopup();
        }
        
        // Ctrl+N for new character
        if (event.ctrlKey && event.key === 'n') {
            event.preventDefault();
            openCharacterCreationPopup();
        }
        
        // Ctrl+U for user details
        if (event.ctrlKey && event.key === 'u') {
            event.preventDefault();
            openUserDetailsPopup();
        }
    });
    
    // Close popups when clicking outside - use event delegation to prevent blocking
    document.addEventListener('click', function(event) {
        const characterOverlay = document.getElementById('character-creation-overlay');
        const userOverlay = document.getElementById('user-details-overlay');
        
        if (event.target === characterOverlay) {
            closeCharacterCreationPopup();
        }
        if (event.target === userOverlay) {
            closeUserDetailsPopup();
        }
    });
});

// Export functions for potential use by other modules
window.characterModule = {
    openCharacterCreationPopup,
    closeCharacterCreationPopup,
    openUserDetailsPopup,
    closeUserDetailsPopup,
    saveCharacter,
    saveUserDetails,
    loadUserDetails,
    loadCharacters
};
