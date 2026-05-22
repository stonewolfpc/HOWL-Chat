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
    if (!window.go || !window.go.main || !window.go.main.App) {
        console.error('Wails runtime not available');
        return;
    }

    // Ensure an active chat exists before saving
    const hasActiveChat = window.dataSync?.currentChatId ||
                          window.chatBridge?.getState()?.sessionId ||
                          window.currentChatId;
    if (!hasActiveChat && window.dataSync?.createNewChat) {
        console.log('[saveChatSession] No active chat — auto-creating before save');
        const chatMeta = await window.dataSync.createNewChat(`Chat ${new Date().toLocaleString()}`);
        if (!chatMeta) {
            console.error('Failed to auto-create chat before save');
            return;
        }
    }

    try {
        await window.go.main.App.SaveChat();
        console.log('Chat saved successfully with all context data');
    } catch (error) {
        console.error('Failed to save chat:', error);
        alert('Failed to save chat: ' + error);
    }
}

// Collect complete chat data including all context for saving
function collectCompleteChatData() {
    if (typeof dataSync !== 'undefined' && typeof dataSync.collectCompleteChatData === 'function') {
        return dataSync.collectCompleteChatData();
    }
    if (typeof contextDisplay !== 'undefined') {
        return {
            meta: {
                id: Date.now().toString(),
                name: contextDisplay.currentChatData?.chatName || 'Chat Session',
                createdAt: new Date().toISOString(),
                updatedAt: new Date().toISOString(),
                worldName: '',
                scenarioName: '',
                characterNames: [],
                lorebookNames: [],
                messageCount: 0
            },
            history: [],
            world: null,
            scenario: null,
            characters: [],
            lorebooks: [],
            memories: []
        };
    }
    console.error('No chat data collector available');
    return null;
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
function deleteChatSession(sessionId, sessionName, event) {
    // Prevent event bubbling to avoid triggering session selection
    if (event) {
        event.stopPropagation();
        event.preventDefault();
    }
    
    if (!sessionId) {
        alert('No session selected to delete');
        return;
    }
    
    if (confirm(`Are you sure you want to delete "${sessionName}"?\n\nThis will permanently delete the chat and all its associated automated lorebook entries.`)) {
        // Stop frontend autosave so it doesn't recreate the directory after backend removal
        if (window.dataSync) {
            window.dataSync.stopAutosave();
        }
        // Call backend to delete chat
        if (window.go && window.go.main && window.go.main.App) {
            window.go.main.App.DeleteChat(sessionId).then(() => {
                console.log('Chat session deleted:', sessionId);

                // Remove from local list and refresh
                savedSessions = savedSessions.filter(s => s.id !== sessionId);

                // Clear selection if deleted session was selected
                if (selectedSessionId === sessionId) {
                    selectedSessionId = null;
                    const loadButton = document.querySelector('.session-list-load');
                    if (loadButton) {
                        loadButton.disabled = true;
                    }
                }

                // If this was the active session, reset bridge state
                const activeId = window.chatBridge?.getState()?.sessionId;
                if (activeId === sessionId && window.chatBridge) {
                    window.chatBridge.clearForNewChat('New Chat');
                }

                renderSessionList();

                // Force re-fetch from backend to ensure sync
                if (typeof loadSavedSessions === 'function') {
                    loadSavedSessions().catch(e => console.warn('Post-delete refresh failed:', e));
                }

                showNotification(`Chat session deleted: ${sessionName}`, 'success');
            }).catch(err => {
                console.error('Failed to delete chat session:', err);
                alert('Failed to delete chat session: ' + err);
            });
        } else {
            // Fallback for development
            savedSessions = savedSessions.filter(s => s.id !== sessionId);
            if (selectedSessionId === sessionId) {
                selectedSessionId = null;
            }
            renderSessionList();
            console.log('Session deleted (fallback mode)');
        }
    }
}

// Collect chat history from DOM
function collectChatHistory() {
    if (!currentSession) return;
    
    const messages = [];
    document.querySelectorAll('.chat-message').forEach(msg => {
        const sender = msg.querySelector('.chat-sender')?.textContent || '';
        const text = msg.querySelector('.chat-text')?.textContent || '';
        const role = msg.classList.contains('system') ? 'system' : msg.classList.contains('user') ? 'user' : 'assistant';
        
        messages.push({
            id: Date.now().toString() + Math.random().toString(36).substr(2, 9),
            sender: sender,
            text: text,
            role: role,
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
    if (!window.go || !window.go.main || !window.go.main.App) {
        console.error('Wails runtime not available');
        savedSessions = [];
        renderSessionList();
        return;
    }
    
    try {
        const chatSummaries = await window.go.main.App.GetChatList();
        
        // Handle case where chatSummaries might be null or undefined
        if (!chatSummaries || !Array.isArray(chatSummaries)) {
            console.log('No chat summaries available');
            savedSessions = [];
            renderSessionList();
            return;
        }
        
        savedSessions = chatSummaries.map(summary => ({
            id: summary.id || '',
            name: summary.name || 'Untitled Chat',
            createdAt: summary.createdAt || new Date().toISOString(),
            updatedAt: summary.updatedAt || new Date().toISOString(),
            metadata: {
                messageCount: summary.messageCount || 0,
                characterCount: (summary.characterNames && Array.isArray(summary.characterNames)) ? summary.characterNames.length : 0,
                hasImages: false
            },
            worldName: summary.worldName || '',
            scenarioName: summary.scenarioName || '',
            characterNames: (summary.characterNames && Array.isArray(summary.characterNames)) ? summary.characterNames : [],
            lorebookNames: (summary.lorebookNames && Array.isArray(summary.lorebookNames)) ? summary.lorebookNames : []
        }));
        
        renderSessionList();
    } catch (error) {
        console.error('Failed to load chat list:', error);
        savedSessions = [];
        renderSessionList();
    }
}

// Render session list
function renderSessionList() {
    const grid = document.getElementById('session-list-grid');
    if (!grid) return;

    if (savedSessions.length === 0) {
        grid.innerHTML = `
            <div class="session-list-empty">
                <div class="session-list-empty-icon"><span class="icon icon-book" style="font-size:2em;"></span></div>
                <div class="session-list-empty-text">No saved sessions</div>
                <div class="session-list-empty-subtext">Save a chat session to see it here</div>
            </div>
        `;
        return;
    }

    const maxRecent = window.settingsState?.current?.general?.maxRecentChats ?? 50;
    const archivedIds = getArchivedIds();
    const sorted = savedSessions.slice().sort((a, b) => {
        const da = new Date(a.updatedAt || a.createdAt).getTime();
        const db = new Date(b.updatedAt || b.createdAt).getTime();
        return db - da;
    });
    const display = sorted
        .filter(s => !archivedIds.includes(s.id))
        .slice(0, Math.max(1, maxRecent));

    grid.innerHTML = display.map(session => `
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
            <button class="session-delete-btn" onclick="deleteChatSession('${session.id}', '${session.name || 'Unnamed Session'}', event)" title="Delete chat session">
                <span class="icon icon-trash"></span>
            </button>
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
async function loadSelectedSession() {
    if (!selectedSessionId) {
        return;
    }
    
    if (!window.go || !window.go.main || !window.go.main.App) {
        console.error('Wails runtime not available');
        return;
    }
    
    try {
        const chatData = await window.go.main.App.LoadChat(selectedSessionId);
        console.log('Chat loaded:', chatData);
        if (typeof window.logChatContract === 'function') {
            window.logChatContract('filemenu-load:backend', chatData);
        }
        
        // Clear current chat UI
        const chatMessages = document.querySelector('.chat-messages');
        if (chatMessages) {
            chatMessages.innerHTML = '';
        }
        
        // Rebuild visible chat history from loaded backend data
        if (chatData && Array.isArray(chatData.History)) {
            rebuildChatFromLoadedHistory(chatData.History);
        }

        // Update session info display
        if (chatData && chatData.Meta) {
            updateSessionInfoDisplayFromMeta({
                name: chatData.Meta.Name,
                worldName: chatData.Meta.WorldName,
                scenarioName: chatData.Meta.ScenarioName,
                lorebookNames: chatData.Meta.LorebookNames || []
            });
        }

        // Refresh context panels using the loaded chat session.
        // Prefer sidebar loader first so shared currentChatData stays in sync.
        if (typeof window.loadChatSessionById === 'function') {
            await window.loadChatSessionById(selectedSessionId);
        } else if (window.contextDisplay && typeof window.contextDisplay.loadChatSessionById === 'function') {
            await window.contextDisplay.loadChatSessionById(selectedSessionId);
        }
        if (typeof window.logChatContract === 'function' && typeof currentChatData !== 'undefined') {
            window.logChatContract('filemenu-load:post-hydrate', currentChatData);
        }
        
        closeSessionList();

    } catch (error) {
        console.error('Failed to load chat:', error);
        alert('Failed to load chat: ' + error);
    }
}

function rebuildChatFromLoadedHistory(history) {
    const chatMessages = document.querySelector('.chat-messages');
    if (!chatMessages) return;

    chatMessages.innerHTML = '';

    history.forEach(msg => {
        const role = (msg.Role || msg.role || '').toLowerCase();
        const sender = msg.Sender || msg.sender || (role === 'assistant' ? 'AI' : role === 'user' ? 'You' : 'System');
        const text = msg.Text || msg.text || '';

        const messageDiv = document.createElement('div');
        messageDiv.className = `chat-message ${role === 'user' ? 'user' : 'assistant'}`;
        messageDiv.innerHTML = `
            <div class="chat-bubble">
                <div class="chat-sender">${sender}</div>
                <div class="chat-text">${text}</div>
            </div>
        `;
        chatMessages.appendChild(messageDiv);
    });
}
window.rebuildChatFromLoadedHistory = rebuildChatFromLoadedHistory;

// Update session info display from chat meta
function updateSessionInfoDisplayFromMeta(chatMeta) {
    // Update the info items in the Chats panel
    const currentInfo = document.querySelector('.info-text');
    if (currentInfo) {
        currentInfo.textContent = `Current: ${chatMeta.name || 'New Chat'}`;
    }
    
    // Update world info
    const worldInfo = document.querySelectorAll('.info-text')[1];
    if (worldInfo) {
        worldInfo.textContent = `World: ${chatMeta.worldName || 'None'}`;
    }
    
    // Update scenario info
    const scenarioInfo = document.querySelectorAll('.info-text')[2];
    if (scenarioInfo) {
        scenarioInfo.textContent = `Scenario: ${chatMeta.scenarioName || 'None'}`;
    }
    
    // Update lorebooks info
    const lorebooksInfo = document.querySelectorAll('.info-text')[3];
    if (lorebooksInfo) {
        lorebooksInfo.textContent = `Lorebooks: ${chatMeta.lorebookNames?.length || 0}`;
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
async function confirmSaveAs() {
    const nameInput = document.getElementById('session-name-input');

    if (!nameInput || !nameInput.value.trim()) {
        alert('Please enter a session name');
        return;
    }

    if (!window.go || !window.go.main || !window.go.main.App) {
        console.error('Wails runtime not available');
        return;
    }

    // Ensure an active chat exists before saving as
    const hasActiveChat = window.dataSync?.currentChatId ||
                          window.chatBridge?.getState()?.sessionId ||
                          window.currentChatId;
    if (!hasActiveChat && window.dataSync?.createNewChat) {
        console.log('[confirmSaveAs] No active chat — auto-creating before save-as');
        const chatMeta = await window.dataSync.createNewChat(nameInput.value.trim());
        if (!chatMeta) {
            console.error('Failed to auto-create chat before save-as');
            return;
        }
        closeSessionModal();
        nameInput.value = '';
        showNotification('Chat created and saved', 'success');
        return;
    }

    try {
        const chatMeta = await window.go.main.App.SaveChatAs(nameInput.value.trim());
        console.log('Chat saved as:', chatMeta);

        closeSessionModal();

        // Clear input
        nameInput.value = '';

        // Refresh chat list if visible
        if (document.getElementById('session-list-overlay').classList.contains('show')) {
            loadSavedSessions();
        }

    } catch (error) {
        console.error('Failed to save chat as:', error);
        alert('Failed to save chat: ' + error);
    }
}

// Autosave system
let _autosaveTimer = null;

function performAutosave() {
    const settings = window.settingsState?.current;
    if (settings && settings.storage.chatAutosave === false) return;
    if (!currentSession || !currentSession.chatHistory.length) return;
    if (!window.go?.main?.App?.SaveChat) return;

    collectChatHistory();
    updateSessionMetadata();

    window.go.main.App.SaveChat()
        .then(() => console.log('Autosave completed'))
        .catch(err => console.error('Autosave failed:', err));
}

function restartAutosaveTimer() {
    if (_autosaveTimer) clearInterval(_autosaveTimer);
    const minutes = window.settingsState?.current?.general?.autoSaveInterval ?? 5;
    _autosaveTimer = setInterval(performAutosave, minutes * 60 * 1000);
}
window.restartAutosaveTimer = restartAutosaveTimer;

// ---------------------------------------------------------------------------
// Auto-Archive
// ---------------------------------------------------------------------------

const ARCHIVED_KEY = 'howl-archived-sessions';
let _archiveTimer = null;

function getArchivedIds() {
    try { return JSON.parse(localStorage.getItem(ARCHIVED_KEY) || '[]'); }
    catch (e) { return []; }
}

function setArchivedIds(ids) {
    localStorage.setItem(ARCHIVED_KEY, JSON.stringify(ids));
}

function isArchived(sessionId) {
    return getArchivedIds().includes(sessionId);
}

async function performAutoArchive() {
    const s = window.settingsState?.current;
    if (!s || s.storage.autoArchive !== true) return;
    const days = s.storage.archiveDays ?? 30;
    const cutoff = Date.now() - (days * 24 * 60 * 60 * 1000);

    const archived = getArchivedIds();
    let changed = false;

    savedSessions.forEach(session => {
        const updated = new Date(session.updatedAt || session.createdAt).getTime();
        if (updated < cutoff && !archived.includes(session.id)) {
            archived.push(session.id);
            changed = true;
        }
    });

    if (changed) {
        setArchivedIds(archived);
        renderSessionList();
        console.log('Auto-archived sessions older than', days, 'days');
    }
}

function restartArchiveTimer() {
    if (_archiveTimer) clearInterval(_archiveTimer);
    _archiveTimer = setInterval(performAutoArchive, 60 * 60 * 1000); // hourly check
}
window.restartArchiveTimer = restartArchiveTimer;

// Initialize autosave on page load (only if dataSync is not present — dataSync handles autosave)
window.addEventListener('load', () => {
    if (typeof dataSync === 'undefined' || !dataSync.setupAutosave) {
        restartAutosaveTimer();
    }
    restartArchiveTimer();

    // Autosave on page unload (only if dataSync is not present)
    window.addEventListener('beforeunload', () => {
        if (typeof dataSync === 'undefined' || !dataSync.setupAutosave) {
            performAutosave();
        }
    });
});

console.log('Chat session management initialized');
