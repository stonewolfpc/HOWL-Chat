/**
 * Sidebar Panels Implementation
 * 
 * This file handles all sidebar functionality including:
 * - Character Cards Panel (left sidebar)
 * - World/Scenario/New Chat/Lorebooks Plugins (right sidebar)
 * 
 * @module sidebar-panels
 * @version 1.0.0
 */

// Global state for sidebar data
let currentChatData = {
    characters: [],
    world: null,
    scenario: null,
    lorebooks: [],
    chatName: 'New Chat'
};

// Character status types
const CharacterStatus = {
    ONLINE: 'online',
    READY: 'ready',
    IDLE: 'idle',
    OFFLINE: 'offline'
};

// Initialize sidebar panels with multiple fallbacks
function initializeSidebarPanels() {
    console.log('Initializing Sidebar Panels...');
    initializeCharacterPanel();
    initializeRightSidebarPlugins();
    loadChatData();
}

// Try multiple initialization methods for robustness
document.addEventListener('DOMContentLoaded', initializeSidebarPanels);
window.addEventListener('load', initializeSidebarPanels);

// Also try immediate execution if DOM is already ready
if (document.readyState === 'loading') {
    // DOM is still loading, wait for events
} else {
    // DOM is already ready, initialize immediately
    setTimeout(initializeSidebarPanels, 100);
}

// ==================== CHARACTER PANEL (LEFT SIDEBAR) ====================

function initializeCharacterPanel() {
    console.log('Initializing Character Panel...');
    // Panel will be populated when chat data is loaded
}

function loadChatData() {
    // Load current chat data from backend
    if (window.go && window.go.main && window.go.main.App) {
        window.go.main.App.LoadChatData().then(data => {
            if (data) {
                currentChatData = data;
                updateCharacterPanel();
                updateRightSidebarInfo();
            } else {
                // Initialize empty state if backend returns null
                resetToDefaultState();
                updateCharacterPanel();
                updateRightSidebarInfo();
            }
        }).catch(err => {
            console.error('Failed to load chat data:', err);
            // Fallback to localStorage for demo
            loadChatDataFromStorage();
        });
    } else {
        // Fallback for development
        loadChatDataFromStorage();
    }
}

function resetToDefaultState() {
    currentChatData = {
        characters: [],
        world: null,
        scenario: null,
        lorebooks: [],
        chatName: 'New Chat'
    };
}

function loadChatDataFromStorage() {
    const savedData = localStorage.getItem('howl-chat-current-chat');
    if (savedData) {
        currentChatData = JSON.parse(savedData);
    } else {
        // Default empty chat data
        resetToDefaultState();
    }
    updateCharacterPanel();
    updateRightSidebarInfo();
}

function updateCharacterPanel() {
    const characterList = document.querySelector('.character-list');
    if (!characterList) return;

    characterList.innerHTML = '';

    if (!currentChatData.characters || currentChatData.characters.length === 0) {
        characterList.innerHTML = `
            <div class="character-empty">
                <div class="empty-icon">👥</div>
                <div class="empty-text">No characters in chat</div>
                <button class="add-character-btn" onclick="addCharacterToChat()">Add Character</button>
            </div>
        `;
        return;
    }

    currentChatData.characters.forEach(character => {
        const characterItem = createCharacterItem(character);
        characterList.appendChild(characterItem);
    });
}

function createCharacterItem(character) {
    const item = document.createElement('div');
    item.className = 'character-item';
    item.dataset.characterId = character.id;
    
    // Determine status display
    const statusClass = character.status || CharacterStatus.OFFLINE;
    const statusText = statusClass.charAt(0).toUpperCase() + statusClass.slice(1);
    
    item.innerHTML = `
        <div class="character-avatar" onclick="openCharacterSheet('${character.id}')">
            ${character.image ? 
                `<img src="${character.image}" alt="${character.name}" onerror="this.style.display='none'; this.nextElementSibling.style.display='flex';">
                 <div class="avatar-placeholder" style="display:none;">👤</div>` :
                `<div class="avatar-placeholder">👤</div>`
            }
            <div class="status-indicator ${statusClass}"></div>
        </div>
        <div class="character-info" onclick="openCharacterSheet('${character.id}')">
            <div class="character-name">${character.name}</div>
            <div class="character-status">${statusText}</div>
        </div>
        <div class="character-actions">
            <button class="character-menu-btn" onclick="showCharacterMenu(event, '${character.id}')">⋮</button>
        </div>
    `;

    // Add right-click menu
    item.addEventListener('contextmenu', (e) => {
        e.preventDefault();
        showCharacterMenu(e, character.id);
    });

    return item;
}

function openCharacterSheet(characterId) {
    console.log('Opening character sheet for:', characterId);

    if (window.go && window.go.main && window.go.main.App) {
        window.go.main.App.OpenCharacterSheet(characterId).then(() => {
            console.log('Character sheet opened');
        }).catch(err => {
            console.error('Failed to open character sheet:', err);
            navigateTo('character.html?id=' + characterId);
        });
    } else {
        navigateTo('character.html?id=' + characterId);
    }
}

function showCharacterMenu(event, characterId) {
    event.preventDefault();
    event.stopPropagation();

    // Remove existing menus
    const existingMenu = document.querySelector('.character-context-menu');
    if (existingMenu) {
        existingMenu.remove();
    }

    const menu = document.createElement('div');
    menu.className = 'character-context-menu';
    menu.innerHTML = `
        <div class="menu-item" onclick="editCharacter('${characterId}')">✏️ Edit Character</div>
        <div class="menu-item" onclick="replaceCharacterImage('${characterId}')">🖼️ Replace Image</div>
        <div class="menu-item" onclick="viewCharacterMemories('${characterId}')">🧠 View Memories</div>
        <div class="menu-separator"></div>
        <div class="menu-item danger" onclick="removeCharacterFromChat('${characterId}')">🗑️ Remove from Chat</div>
    `;

    // Position menu
    menu.style.position = 'absolute';
    menu.style.left = event.pageX + 'px';
    menu.style.top = event.pageY + 'px';
    menu.style.zIndex = '10000';

    document.body.appendChild(menu);

    // Close menu when clicking outside
    setTimeout(() => {
        document.addEventListener('click', function closeMenu(e) {
            if (!menu.contains(e.target)) {
                menu.remove();
                document.removeEventListener('click', closeMenu);
            }
        });
    }, 100);
}

function editCharacter(characterId) {
    console.log('Editing character:', characterId);
    navigateTo('character.html?id=' + characterId);
}

function replaceCharacterImage(characterId) {
    console.log('Replacing image for character:', characterId);
    
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'image/*';
    input.onchange = (e) => {
        const file = e.target.files[0];
        if (file) {
            // Handle image upload
            const reader = new FileReader();
            reader.onload = (e) => {
                const imageData = e.target.result;
                updateCharacterImage(characterId, imageData);
            };
            reader.readAsDataURL(file);
        }
    };
    input.click();
}

function updateCharacterImage(characterId, imageData) {
    const character = currentChatData.characters.find(c => c.id === characterId);
    if (character) {
        character.image = imageData;
        saveChatData();
        updateCharacterPanel();
        showNotification('Character image updated', 'success');
    }
}

function viewCharacterMemories(characterId) {
    if (window.CharacterMemories && typeof window.CharacterMemories.view === 'function') {
        window.CharacterMemories.view(characterId);
    } else {
        console.log('Viewing memories for character:', characterId);
    }
}

function removeCharacterFromChat(characterId) {
    if (confirm('Remove this character from the chat?')) {
        currentChatData.characters = currentChatData.characters.filter(c => c.id !== characterId);
        saveChatData();
        updateCharacterPanel();
        showNotification('Character removed from chat', 'success');
    }
}

function addCharacterToChat() {
    console.log('Adding character to chat...');
    navigateTo('character.html');
}

// ==================== RIGHT SIDEBAR PLUGINS ====================

function initializeRightSidebarPlugins() {
    console.log('Initializing Right Sidebar Plugins...');
    // Plugins are now hardcoded in HTML, don't add them dynamically
    loadPluginData();
}

function loadPluginData() {
    loadWorlds();
    loadScenarios();
    loadLorebooks();
}

function loadWorlds() {
    // Load from localStorage since that's where world.js saves data
    const worlds = JSON.parse(localStorage.getItem('howl-worlds') || '[]');
    populateWorldSelector(worlds);
}

function populateWorldSelector(worlds) {
    const selector = document.getElementById('world-selector');
    if (!selector) return;

    // Clear existing options except the first one
    while (selector.children.length > 1) {
        selector.removeChild(selector.lastChild);
    }

    worlds.forEach(world => {
        const option = document.createElement('option');
        option.value = world.id;
        option.textContent = world.title || world.name;
        selector.appendChild(option);
    });

    // Select current world if set
    if (currentChatData.world) {
        selector.value = currentChatData.world.id;
    }
}

function loadScenarios() {
    // Load from localStorage since that's where scenario.js saves data
    const scenarios = JSON.parse(localStorage.getItem('howl-scenarios') || '[]');
    populateScenarioSelector(scenarios);
}

function populateScenarioSelector(scenarios) {
    const selector = document.getElementById('scenario-selector');
    if (!selector) return;

    // Clear existing options except the first one
    while (selector.children.length > 1) {
        selector.removeChild(selector.lastChild);
    }

    scenarios.forEach(scenario => {
        const option = document.createElement('option');
        option.value = scenario.id;
        option.textContent = scenario.title || scenario.name;
        selector.appendChild(option);
    });

    // Select current scenario if set
    if (currentChatData.scenario) {
        selector.value = currentChatData.scenario.id;
    }
}

function loadLorebooks() {
    const lorebookList = document.getElementById('lorebook-list');
    if (!lorebookList) return;

    lorebookList.innerHTML = '';

    // Load world lorebook
    if (currentChatData.world) {
        const worldLorebook = document.createElement('div');
        worldLorebook.className = 'lorebook-item';
        worldLorebook.innerHTML = `
            <span class="lorebook-name">🌍 ${currentChatData.world.title || currentChatData.world.name}</span>
            <button class="lorebook-edit-btn" onclick="editLorebook('world')">Edit</button>
        `;
        lorebookList.appendChild(worldLorebook);
    }

    // Load scenario lorebook
    if (currentChatData.scenario) {
        const scenarioLorebook = document.createElement('div');
        scenarioLorebook.className = 'lorebook-item';
        scenarioLorebook.innerHTML = `
            <span class="lorebook-name">📜 ${currentChatData.scenario.title || currentChatData.scenario.name}</span>
            <button class="lorebook-edit-btn" onclick="editLorebook('scenario')">Edit</button>
        `;
        lorebookList.appendChild(scenarioLorebook);
    }

    // Load character lorebooks
    currentChatData.characters.forEach(character => {
        const characterLorebook = document.createElement('div');
        characterLorebook.className = 'lorebook-item';
        characterLorebook.innerHTML = `
            <span class="lorebook-name">👤 ${character.name}</span>
            <button class="lorebook-edit-btn" onclick="editLorebook('${character.id}')">Edit</button>
        `;
        lorebookList.appendChild(characterLorebook);
    });

    if (lorebookList.children.length === 0) {
        lorebookList.innerHTML = '<div class="no-lorebooks">No lorebooks available</div>';
    }
}

// Plugin action functions
function selectWorld(worldId) {
    if (!worldId) {
        currentChatData.world = null;
    } else {
        // Load world data from localStorage
        const worlds = JSON.parse(localStorage.getItem('howl-worlds') || '[]');
        const world = worlds.find(w => w.id === worldId);
        if (world) {
            currentChatData.world = world;
            saveChatData();
            updateRightSidebarInfo();
            loadLorebooks();
            showNotification(`World set: ${world.title || world.name}`, 'success');
        }
    }
}

function selectScenario(scenarioId) {
    if (!scenarioId) {
        currentChatData.scenario = null;
    } else {
        // Load scenario data from localStorage
        const scenarios = JSON.parse(localStorage.getItem('howl-scenarios') || '[]');
        const scenario = scenarios.find(s => s.id === scenarioId);
        if (scenario) {
            currentChatData.scenario = scenario;
            saveChatData();
            updateRightSidebarInfo();
            loadLorebooks();
            showNotification(`Scenario set: ${scenario.title || scenario.name}`, 'success');
        }
    }
}

function createNewWorld() {
    navigateTo('world.html');
}

function editCurrentWorld() {
    if (currentChatData.world) {
        navigateTo('world.html');
    } else {
        showNotification('No world selected', 'warning');
    }
}

function createNewScenario() {
    navigateTo('scenario.html');
}

function editCurrentScenario() {
    if (currentChatData.scenario) {
        navigateTo('scenario.html');
    } else {
        showNotification('No scenario selected', 'warning');
    }
}

function startNewChat() {
    if (confirm('Start a new chat? Current chat will be saved.')) {
        // Save current chat
        saveChatData();
        
        // Clear current chat data
        currentChatData = {
            characters: [],
            world: null,
            scenario: null,
            lorebooks: [],
            chatName: 'New Chat'
        };
        
        // Clear chat messages
        const chatMessages = document.querySelector('.chat-messages');
        if (chatMessages) {
            chatMessages.innerHTML = '';
        }
        
        // Update UI
        updateCharacterPanel();
        updateRightSidebarInfo();
        loadPluginData();
        
        showNotification('New chat started', 'success');
    }
}

function loadChatSession() {
    // Use existing loadChatSession function
    if (typeof loadChatSession === 'function') {
        loadChatSession();
    } else {
        showNotification('Load chat functionality coming soon', 'info');
    }
}

function saveChatSessionAs() {
    // Use existing saveChatSessionAs function
    if (typeof saveChatSessionAs === 'function') {
        saveChatSessionAs();
    } else {
        showNotification('Save chat functionality coming soon', 'info');
    }
}

function createLorebookEntry() {
    navigateTo('lorebooks.html');
}

function viewAutomatedMemories() {
    if (typeof toggleMemoryList === 'function') {
        toggleMemoryList();
    }
}

function editLorebook(lorebookType) {
    navigateTo('lorebooks.html');
}

function updateRightSidebarInfo() {
    // Update info items in right sidebar using specific IDs
    const chatName = document.getElementById('current-chat-name');
    const worldName = document.getElementById('current-world-name');
    const scenarioName = document.getElementById('current-scenario-name');
    const lorebooksCount = document.getElementById('current-lorebooks-count');
    
    if (chatName) {
        chatName.textContent = `Current: ${currentChatData.chatName}`;
    }
    if (worldName) {
        worldName.textContent = `World: ${currentChatData.world ? (currentChatData.world.title || currentChatData.world.name) : 'None'}`;
    }
    if (scenarioName) {
        scenarioName.textContent = `Scenario: ${currentChatData.scenario ? (currentChatData.scenario.title || currentChatData.scenario.name) : 'None'}`;
    }
    if (lorebooksCount) {
        lorebooksCount.textContent = `Lorebooks: ${currentChatData.lorebooks.length > 0 ? currentChatData.lorebooks.length : 'None'}`;
    }
}

function saveChatData() {
    // Save to backend
    if (window.go && window.go.main && window.go.main.App) {
        window.go.main.App.SaveChatData(currentChatData).then(() => {
            console.log('Chat data saved to backend');
        }).catch(err => {
            console.error('Failed to save chat data to backend:', err);
            // Fallback to localStorage
            localStorage.setItem('howl-chat-current-chat', JSON.stringify(currentChatData));
        });
    } else {
        // Fallback to localStorage
        localStorage.setItem('howl-chat-current-chat', JSON.stringify(currentChatData));
    }
}

// Export functions for global access
window.sidebarPanels = {
    loadChatData,
    updateCharacterPanel,
    updateRightSidebarInfo,
    selectWorld,
    selectScenario,
    openCharacterSheet,
    showCharacterMenu,
    createNewWorld,
    createNewScenario,
    startNewChat
};
