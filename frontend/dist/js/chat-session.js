/**
 * Chat Session Management Module
 * 
 * This module handles all chat session management including:
 * - Save, Save As, Load, Delete operations
 * - Chat session data structure
 * - Autosave system
 * - Image handling with thumbnails
 * - Modular file format (JSON + separate folders)
 * 
 * @module chat-session
 */

// Chat session state
let currentSession = null;
let savedSessions = [];
let selectedSessionId = null;

// Chat session data structure
class ChatSession {
    constructor() {
        this.id = Date.now().toString();
        this.name = null;
        this.createdAt = new Date().toISOString();
        this.updatedAt = new Date().toISOString();
        this.chatHistory = [];
        this.characters = [];
        this.world = null;
        this.scenario = null;
        this.lorebooks = [];
        this.metadata = {
            version: '1.0.0',
            messageCount: 0,
            characterCount: 0,
            hasImages: false
        };
    }
}

// Save current chat session
async function saveChatSession() {
    if (!currentSession) {
        currentSession = new ChatSession();
    }
    
    if (!currentSession.name) {
        saveChatSessionAs();
        return;
    }
    
    // Collect chat history from DOM
    collectChatHistory();
    
    // Update metadata
    updateSessionMetadata();
    
    // Save to backend (placeholder for now)
    console.log('Saving session:', currentSession);
    
    // Backend integration needed here
    // File path: Documents\HOWL_Chat\Chats\<sessionName>.json
    // Images: Documents\HOWL_Chat\Chats\<sessionName>\images\
    
    alert('Session saved successfully');
}

// Save chat session with new name
function saveChatSessionAs() {
    const overlay = document.getElementById('session-modal-overlay');
    if (overlay) {
        overlay.classList.add('show');
    }
}

// Load saved chat session
function loadChatSession() {
    const overlay = document.getElementById('session-list-overlay');
    if (overlay) {
        overlay.classList.add('show');
        loadSavedSessions();
    }
}

// Delete saved chat session
function deleteChatSession() {
    if (!selectedSessionId) {
        alert('No session selected to delete');
        return;
    }
    
    if (confirm('Are you sure you want to delete this chat session?')) {
        savedSessions = savedSessions.filter(s => s.id !== selectedSessionId);
        selectedSessionId = null;
        console.log('Session deleted');
        
        // Backend integration needed here
        // Delete: Documents\HOWL_Chat\Chats\<sessionName>.json
        // Delete folder: Documents\HOWL_Chat\Chats\<sessionName>\
        
        alert('Session deleted successfully');
    }
}

// Collect chat history from DOM
function collectChatHistory() {
    if (!currentSession) return;
    
    const messages = [];
    document.querySelectorAll('.chat-message').forEach(msg => {
        const sender = msg.querySelector('.chat-sender')?.textContent || '';
        const text = msg.querySelector('.chat-text')?.textContent || '';
        const isUser = msg.classList.contains('user');
        
        messages.push({
            id: Date.now().toString() + Math.random().toString(36).substr(2, 9),
            sender: sender,
            text: text,
            isUser: isUser,
            timestamp: new Date().toISOString(),
            hasImage: false // Future: check for images
        });
    });
    
    currentSession.chatHistory = messages;
}

// Update session metadata
function updateSessionMetadata() {
    if (!currentSession) return;
    
    currentSession.updatedAt = new Date().toISOString();
    currentSession.metadata.messageCount = currentSession.chatHistory.length;
    currentSession.metadata.characterCount = currentSession.characters.length;
    currentSession.metadata.hasImages = currentSession.chatHistory.some(msg => msg.hasImage);
}

// Load saved sessions from backend
async function loadSavedSessions() {
    // Backend integration needed here
    // Load from: Documents\HOWL_Chat\Chats\
    
    // Placeholder data
    savedSessions = [];
    
    renderSessionList();
}

// Render session list
function renderSessionList() {
    const grid = document.getElementById('session-list-grid');
    if (!grid) return;
    
    if (savedSessions.length === 0) {
        grid.innerHTML = `
            <div class="session-list-empty">
                <div class="session-list-empty-icon">💬</div>
                <div class="session-list-empty-text">No saved sessions</div>
                <div class="session-list-empty-subtext">Save a chat session to see it here</div>
            </div>
        `;
        return;
    }
    
    grid.innerHTML = savedSessions.map(session => `
        <div class="session-list-item ${session.id === selectedSessionId ? 'selected' : ''}" 
             onclick="selectSession('${session.id}')">
            <div class="session-item-name">${session.name || 'Unnamed Session'}</div>
            <div class="session-item-meta">
                Messages: ${session.metadata.messageCount} | 
                Characters: ${session.metadata.characterCount}
            </div>
            <div class="session-item-date">
                ${new Date(session.updatedAt).toLocaleString()}
            </div>
        </div>
    `).join('');
}

// Select session
function selectSession(sessionId) {
    selectedSessionId = sessionId;
    renderSessionList();
    
    // Enable load button
    const loadButton = document.querySelector('.session-list-load');
    if (loadButton) {
        loadButton.disabled = false;
    }
}

// Load selected session
function loadSelectedSession() {
    if (!selectedSessionId) {
        return;
    }
    
    const session = savedSessions.find(s => s.id === selectedSessionId);
    if (!session) {
        alert('Session not found');
        return;
    }
    
    // Load session data into current session
    currentSession = session;
    
    // Rebuild chat UI from session history
    rebuildChatFromHistory(session.chatHistory);
    
    // Update session info display
    updateSessionInfoDisplay(session);
    
    closeSessionList();
    
    alert('Session loaded successfully');
}

// Rebuild chat UI from session history
function rebuildChatFromHistory(history) {
    const chatMessages = document.querySelector('.chat-messages');
    if (!chatMessages) return;
    
    chatMessages.innerHTML = '';
    
    history.forEach(msg => {
        const messageDiv = document.createElement('div');
        messageDiv.className = `chat-message ${msg.isUser ? 'user' : 'assistant'}`;
        messageDiv.innerHTML = `
            <div class="chat-bubble">
                <div class="chat-sender">${msg.sender}</div>
                <div class="chat-text">${msg.text}</div>
                <div class="chat-bubble-buttons">
                    <button class="bubble-button" onclick="copyMessage(this)">Copy</button>
                    ${msg.isUser ? 
                        `<button class="bubble-button" onclick="editMessage(this)">Edit</button>` :
                        `<button class="bubble-button" onclick="regenerateMessage(this)">Regenerate</button>
                    `}
                    ${!msg.isUser ? 
                        `<button class="bubble-button" onclick="continueMessage(this)">Continue</button>` :
                        ''
                    }
                </div>
            </div>
        `;
        chatMessages.appendChild(messageDiv);
    });
}

// Update session info display
function updateSessionInfoDisplay(session) {
    // Update the info items in the Chats panel
    const currentInfo = document.querySelector('.info-text');
    if (currentInfo) {
        currentInfo.textContent = `Current: ${session.name || 'New Chat'}`;
    }
    
    // Update world info
    const worldInfo = document.querySelectorAll('.info-text')[1];
    if (worldInfo) {
        worldInfo.textContent = `World: ${session.world?.name || 'None'}`;
    }
    
    // Update scenario info
    const scenarioInfo = document.querySelectorAll('.info-text')[2];
    if (scenarioInfo) {
        scenarioInfo.textContent = `Scenario: ${session.scenario?.name || 'None'}`;
    }
    
    // Update lorebooks info
    const lorebooksInfo = document.querySelectorAll('.info-text')[3];
    if (lorebooksInfo) {
        lorebooksInfo.textContent = `Lorebooks: ${session.lorebooks?.length || 0}`;
    }
}

// Close session modal
function closeSessionModal() {
    const overlay = document.getElementById('session-modal-overlay');
    if (overlay) {
        overlay.classList.remove('show');
    }
}

// Close session list
function closeSessionList() {
    const overlay = document.getElementById('session-list-overlay');
    if (overlay) {
        overlay.classList.remove('show');
    }
}

// Confirm save as
function confirmSaveAs() {
    const nameInput = document.getElementById('session-name-input');
    const descInput = document.getElementById('session-desc-input');
    
    if (!nameInput || !nameInput.value.trim()) {
        alert('Please enter a session name');
        return;
    }
    
    if (!currentSession) {
        currentSession = new ChatSession();
    }
    
    currentSession.name = nameInput.value.trim();
    currentSession.description = descInput.value.trim();
    
    collectChatHistory();
    updateSessionMetadata();
    
    // Save to backend
    console.log('Saving session as:', currentSession);
    
    closeSessionModal();
    
    // Clear inputs
    nameInput.value = '';
    descInput.value = '';
    
    alert('Session saved successfully');
}

// Autosave system
function performAutosave() {
    if (!currentSession || !currentSession.chatHistory.length) {
        return;
    }
    
    collectChatHistory();
    updateSessionMetadata();
    
    const timestamp = new Date().toISOString();
    const date = timestamp.split('T')[0];
    const time = timestamp.split('T')[1].split('.')[0].replace(/:/g, '-');
    
    const autosaveName = `autosave_${date}_${time}`;
    
    // Backend integration needed here
    // Save to: Documents\HOWL_Chat\Autosaves\autosave_YYYY-MM-DD_HH-MM.json
    // Keep only 3 most recent autosaves
    
    console.log('Autosave performed:', autosaveName);
}

// Initialize autosave on page load
window.addEventListener('load', () => {
    // Autosave every 5 minutes
    setInterval(performAutosave, 5 * 60 * 1000);
    
    // Autosave on page unload
    window.addEventListener('beforeunload', () => {
        performAutosave();
    });
});

console.log('Chat session management initialized');
