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

// Session list filters
let _sessionFilterQuery = '';
let _sessionSortField = 'date';
let _sessionSortDir = 'desc';
let _sessionWorldFilter = '';
let _sessionScenarioFilter = '';
let _loadingSession = false;

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
                messageCount: 0,
                // Tracking arrays
                worldId: '', scenarioId: '',
                activeCharacters: [], removedCharacters: [], addedCharacters: [],
                activeLorebooks: [], removedLorebooks: [], addedLorebooks: []
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
            id: summary.id || summary.ID || '',
            name: summary.name || summary.Name || 'Untitled Chat',
            createdAt: summary.createdAt || summary.CreatedAt || new Date().toISOString(),
            updatedAt: summary.updatedAt || summary.UpdatedAt || new Date().toISOString(),
            metadata: {
                messageCount: summary.messageCount || summary.MessageCount || 0,
                characterCount: (summary.characterNames || summary.CharacterNames) ? (summary.characterNames || summary.CharacterNames).length : 0,
                hasImages: false
            },
            worldName: summary.worldName || summary.WorldName || '',
            scenarioName: summary.scenarioName || summary.ScenarioName || '',
            characterNames: Array.isArray(summary.characterNames || summary.CharacterNames) ? (summary.characterNames || summary.CharacterNames) : [],
            lorebookNames: Array.isArray(summary.lorebookNames || summary.LorebookNames) ? (summary.lorebookNames || summary.LorebookNames) : []
        })).filter(s => s.id !== '');
        
        renderSessionList();
        _populateSessionFilterDropdowns();
    } catch (error) {
        console.error('Failed to load chat list:', error);
        savedSessions = [];
        renderSessionList();
        _populateSessionFilterDropdowns();
    }
}

// Render session list
function renderSessionList() {
    console.log('[renderSessionList] called, selectedSessionId:', selectedSessionId, 'savedSessions:', savedSessions.length);
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
    const fe = window.FilterEngine;

    let display = savedSessions.slice().filter(s => !archivedIds.includes(s.id));

    // World filter
    if (_sessionWorldFilter) {
        display = display.filter(s => s.worldName === _sessionWorldFilter);
    }
    // Scenario filter
    if (_sessionScenarioFilter) {
        display = display.filter(s => s.scenarioName === _sessionScenarioFilter);
    }
    // Text search
    if (_sessionFilterQuery) {
        display = fe ? fe.filterItems(display, _sessionFilterQuery, ['name','worldName','scenarioName','characterNames']) : display;
    }
    // Sort
    if (fe) {
        display = fe.sortItems(display, _sessionSortField, _sessionSortDir);
    } else {
        display.sort((a, b) => {
            const da = new Date(a.updatedAt || a.createdAt).getTime();
            const db = new Date(b.updatedAt || b.createdAt).getTime();
            return db - da;
        });
    }

    display = display.slice(0, Math.max(1, maxRecent));

    grid.innerHTML = display.map(session => {
        const safeId = String(session.id).replace(/'/g, "\\'");
        const safeName = String(session.name || 'Unnamed Session').replace(/'/g, "\\'");
        const isSelected = String(session.id) === String(selectedSessionId);
        return `
        <div class="session-list-item ${isSelected ? 'selected' : ''}"
             data-session-id="${session.id}"
             role="button"
             tabindex="0"
             onclick="selectSession('${safeId}')">
            <div class="session-item-name">${session.name || 'Unnamed Session'}</div>
            <div class="session-item-meta">
                Messages: ${session.metadata.messageCount} |
                Characters: ${session.metadata.characterCount}
            </div>
            <div class="session-item-date">
                ${new Date(session.updatedAt).toLocaleString()}
            </div>
            <button class="session-delete-btn" type="button"
                onclick="deleteChatSession('${safeId}', '${safeName}', event)"
                title="Delete chat session">
                <span class="icon icon-trash"></span>
            </button>
        </div>`;
    }).join('');
}

function _populateSessionFilterDropdowns() {
    const worldSelect = document.getElementById('session-world-filter');
    const scenarioSelect = document.getElementById('session-scenario-filter');
    if (!worldSelect || !scenarioSelect) return;

    const worlds = new Set(savedSessions.map(s => s.worldName).filter(Boolean));
    const scenarios = new Set(savedSessions.map(s => s.scenarioName).filter(Boolean));

    const currentWorld = worldSelect.value;
    const currentScenario = scenarioSelect.value;

    worldSelect.innerHTML = '<option value="">All Worlds</option>' + Array.from(worlds).sort().map(w => `<option value="${w}">${w}</option>`).join('');
    scenarioSelect.innerHTML = '<option value="">All Scenarios</option>' + Array.from(scenarios).sort().map(s => `<option value="${s}">${s}</option>`).join('');

    worldSelect.value = currentWorld;
    scenarioSelect.value = currentScenario;
}

window._sessionFilter = (val) => { _sessionFilterQuery = val; renderSessionList(); };
window._sessionSort = (val) => {
    const [field, dir] = val.split(':');
    _sessionSortField = field;
    _sessionSortDir = dir;
    renderSessionList();
};
window._sessionWorldFilter = (val) => { _sessionWorldFilter = val; renderSessionList(); };
window._sessionScenarioFilter = (val) => { _sessionScenarioFilter = val; renderSessionList(); };

// Select session (loads immediately on click)
function selectSession(sessionId) {
    console.log('[selectSession] clicked:', sessionId);
    if (_loadingSession) { console.log('[selectSession] already loading, ignored'); return; }
    selectedSessionId = sessionId;

    // Toggle selected class without re-rendering entire grid
    document.querySelectorAll('.session-list-item').forEach(el => {
        el.classList.toggle('selected', el.dataset.sessionId === sessionId);
    });

    // Load immediately
    loadSelectedSession();
}

// Load selected session
async function loadSelectedSession() {
    console.log('[loadSelectedSession] called with id:', selectedSessionId);
    if (_loadingSession) { console.log('[loadSelectedSession] already loading, ignored'); return; }
    if (!selectedSessionId) {
        console.warn('[loadSelectedSession] no session selected');
        return;
    }
    
    if (!window.go || !window.go.main || !window.go.main.App) {
        console.error('[loadSelectedSession] Wails runtime not available');
        return;
    }
    
    _loadingSession = true;
    try {
        console.log('[loadSelectedSession] calling LoadChat...');
        const chatData = await window.go.main.App.LoadChat(selectedSessionId);
        console.log('[loadSelectedSession] chatData:', chatData);
        if (typeof window.logChatContract === 'function') {
            window.logChatContract('filemenu-load:backend', chatData);
        }
        
        // Clear current chat UI
        const chatMessages = document.querySelector('.chat-messages');
        if (chatMessages) {
            chatMessages.innerHTML = '';
        }
        
        // Rebuild visible chat history from loaded backend data
        const history = chatData?.History || chatData?.history;
        if (Array.isArray(history)) {
            rebuildChatFromLoadedHistory(history);
        } else {
            console.warn('[loadSelectedSession] no history array in chatData');
        }

        // Update session info display
        const meta = chatData?.Meta || chatData?.meta;
        if (meta) {
            updateSessionInfoDisplayFromMeta({
                name: meta.Name || meta.name,
                worldName: meta.WorldName || meta.worldName,
                scenarioName: meta.ScenarioName || meta.scenarioName,
                lorebookNames: meta.LorebookNames || meta.lorebookNames || []
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
    } finally {
        _loadingSession = false;
    }
}

function rebuildChatFromLoadedHistory(history) {
    console.log('[rebuildChatFromLoadedHistory] history length:', history?.length);
    const chatMessages = document.querySelector('.chat-messages');
    if (!chatMessages) {
        console.warn('[rebuildChatFromLoadedHistory] no .chat-messages element');
        return;
    }

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
