/**
 * Sidebar Session Management
 *
 * Handles loading, saving, and syncing chat sessions including:
 * - Auto-discovery of the most recent session
 * - Loading sessions by ID from the Wails backend
 * - Persisting session state to HOWLStorage / backend
 * - Rebuilding the chat message UI from history
 * - Syncing world/scenario/character/lorebook context to backend
 *
 * Depends on: sidebar-state.js
 *
 * @module sidebar-session
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Session Discovery
// ---------------------------------------------------------------------------

/**
 * Queries the backend for all chat sessions and returns the most recently
 * updated one, or null if none are found or the backend is unavailable.
 * @returns {Promise<Object|null>}
 */
function _safeDate(val) {
    const d = new Date(val || 0);
    return isNaN(d.getTime()) ? new Date(0) : d;
}

async function discoverMostRecentChatSession() {
    if (!window.go || !window.go.main || !window.go.main.App) {
        console.log('Backend not available for chat discovery');
        return null;
    }

    try {
        const chatSummaries = await window.go.main.App.GetChatList();
        if (!chatSummaries || chatSummaries.length === 0) {
            console.log('No chat sessions found');
            return null;
        }

        const mostRecent = chatSummaries.reduce((best, current) => {
            const bestDate    = _safeDate(best.updatedAt    || best.createdAt);
            const currentDate = _safeDate(current.updatedAt || current.createdAt);
            return currentDate > bestDate ? current : best;
        });

        console.log('[discoverMostRecentChatSession] found:', mostRecent.name || mostRecent.Name, mostRecent.id || mostRecent.ID);
        return mostRecent;

    } catch (error) {
        console.error('Failed to discover chat sessions:', error);
        return null;
    }
}

// ---------------------------------------------------------------------------
// Session Loading
// ---------------------------------------------------------------------------

/**
 * Waits for the Wails Go bridge to become available.
 * @param {number} timeout Max ms to wait (default 10000).
 * @returns {Promise<boolean>}
 */
function waitForWails(timeout = 10000) {
    return new Promise((resolve) => {
        if (window.go && window.go.main && window.go.main.App) {
            resolve(true);
            return;
        }
        let attempts = 0;
        const maxAttempts = Math.max(1, Math.floor(timeout / 100));
        const interval = setInterval(() => {
            if (window.go && window.go.main && window.go.main.App) {
                clearInterval(interval);
                resolve(true);
            } else if (++attempts >= maxAttempts) {
                clearInterval(interval);
                resolve(false);
            }
        }, 100);
    });
}

/**
 * Top-level entry point. Restores the backend's explicit LastChatID first;
 * falls back to discovering the most recently updated session, then to
 * LoadChatData / HOWLStorage.
 */
async function loadChatData() {
    console.log('[loadChatData] starting startup load sequence...');
    const wailsReady = await waitForWails(10000);
    console.log('[loadChatData] Wails ready:', wailsReady);

    if (wailsReady) {
        // 1. Try the backend's explicit LastChatID first (most accurate)
        try {
            const activeMeta = await window.go.main.App.GetActiveChatMeta();
            console.log('[loadChatData] GetActiveChatMeta returned:', activeMeta);
            const activeId   = activeMeta && (activeMeta.id   || activeMeta.ID);
            const activeName = activeMeta && (activeMeta.name || activeMeta.Name);
            if (activeId) {
                console.log('[loadChatData] restoring last active chat:', activeName || activeId);
                await loadChatSessionById(activeId);
                return;
            }
            console.log('[loadChatData] no active chat ID in meta');
        } catch (err) {
            console.warn('[loadChatData] GetActiveChatMeta failed, falling back to discovery:', err);
        }

        // 2. Fallback: discover the most recently updated session
        const recentSession = await discoverMostRecentChatSession();
        if (recentSession) {
            const recentId   = recentSession.id   || recentSession.ID;
            const recentName = recentSession.name || recentSession.Name;
            console.log('[loadChatData] auto-loading most recent session:', recentName || recentId);
            await loadChatSessionById(recentId);
            return;
        }
    }

    // 3. Final fallback: generic LoadChatData or frontend cache
    if (window.go && window.go.main && window.go.main.App) {
        try {
            console.log('[loadChatData] falling back to LoadChatData');
            const data = await window.go.main.App.LoadChatData();
            if (data) {
                window.sidebarState.currentChatData = data;
                await updateContextDisplay();
            } else {
                window.sidebarState.resetToDefaultState();
                await updateContextDisplay();
            }
        } catch (err) {
            console.error('[loadChatData] LoadChatData failed:', err);
            await loadChatDataFromStorage();
        }
    } else {
        console.log('[loadChatData] backend unavailable, loading from storage');
        await loadChatDataFromStorage();
    }
}

/**
 * Loads a specific chat session by ID, hydrates currentChatData,
 * persists to storage, rebuilds the chat history UI, and syncs to backend.
 * @param {string} sessionId
 * @returns {Promise<boolean>}
 */
async function loadChatSessionById(sessionId) {
    if (!window.go || !window.go.main || !window.go.main.App) {
        console.error('Wails runtime not available');
        return false;
    }

    try {
        const chatData = await window.go.main.App.LoadChat(sessionId);
        console.log('Chat session loaded:', chatData);

        // Normalize: backend may return PascalCase or camelCase
        const world      = chatData.World      || chatData.world      || null;
        const scenario   = chatData.Scenario   || chatData.scenario   || null;
        const characters = chatData.Characters || chatData.characters || [];
        const lorebooks  = chatData.Lorebooks  || chatData.lorebooks  || [];
        const history    = chatData.History    || chatData.history    || [];
        const meta       = chatData.Meta       || chatData.meta       || {};

        if (window.sidebarState) window.sidebarState.logChatContract('load:backend', chatData);

        if (history.length > 0) {
            const first = history[0];
            const last  = history[history.length - 1];
            console.log('[TRACE:loadChatSessionById] first msg keys:', Object.keys(first),
                'role=', first.role || first.Role || '(missing)',
                'sender=', first.sender || first.Sender || '(missing)');
            console.log('[TRACE:loadChatSessionById] last  msg keys:', Object.keys(last),
                'role=', last.role || last.Role || '(missing)',
                'sender=', last.sender || last.Sender || '(missing)');
        }
        console.log('Loaded chat data:', {
            hasWorld:       !!world,
            hasScenario:    !!scenario,
            characters:     characters.length,
            lorebooks:      lorebooks.length,
            historyCount:   history.length,
            firstMessage:   history.length > 0 ? `${history[0].Sender || history[0].sender}: ${(history[0].Text || history[0].text || '').slice(0, 60)}` : '(none)',
            lastMessage:    history.length > 1 ? `${history[history.length-1].Sender || history[history.length-1].sender}: ${(history[history.length-1].Text || history[history.length-1].text || '').slice(0, 60)}` : '(none)'
        });

        if (window.sidebarState) {
            window.sidebarState.currentChatData = {
                world, scenario, characters, lorebooks, sessionId,
                sessionName:      meta.Name      || meta.name      || 'New Chat',
                sessionCreatedAt: meta.CreatedAt || meta.createdAt || null,
                sessionUpdatedAt: meta.UpdatedAt || meta.updatedAt || null,
                chatName:         meta.Name      || meta.name      || 'New Chat'
            };
        }

        if (window.dataSync) {
            window.dataSync.currentChatId = sessionId;
        }

        if (window.HOWLStorage) {
            await window.HOWLStorage.setChatData('howl-chat-current-chat',
                window.sidebarState ? window.sidebarState.currentChatData : { world, scenario, characters, lorebooks });
        }
        if (window.sidebarState) window.sidebarState.logChatContract('load:hydrated', window.sidebarState.currentChatData);

        // Hydrate both state systems and all panels via the bridge
        if (window.chatBridge) {
            await window.chatBridge.hydrateFromBackend(chatData);
        } else {
            if (history.length > 0) {
                rebuildChatFromHistory(history);
            } else {
                const chatMessages = document.querySelector('.chat-messages');
                if (chatMessages) chatMessages.innerHTML = '';
            }
            await updateAllChatPanels();
        }

        // Do NOT syncAllDataToBackend here — LoadChat already set the backend state.
        // Re-syncing would overwrite chat-specific character/lorebook selections
        // with the full scenario defaults (e.g. all scenario characters).
        updateSessionInfoDisplay(meta);
        const loadedName = meta.Name || meta.name || 'New Chat';
        showNotification(`Chat session loaded: ${loadedName}`, 'success');
        return true;

    } catch (error) {
        console.error('Failed to load chat session:', error);
        return false;
    }
}

/**
 * Loads chat data from HOWLStorage (IndexedDB/localStorage fallback).
 */
async function loadChatDataFromStorage() {
    const savedData = await window.HOWLStorage.getChatData('howl-chat-current-chat');
    if (!savedData) {
        window.sidebarState.resetToDefaultState();
        await updateContextDisplay();
        return;
    }

    // The backend membership authority decides whether a cached session may be
    // restored. A missing or orphaned entry is in the list but is not loadable,
    // and an unbound cache is not evidence of a conversation at all: pushing
    // either back would let this shadow copy mint a chat the authority never
    // created.
    const cachedSessionId = savedData.sessionId || (window.chatBridge && window.chatBridge.getState()?.sessionId);
    if (!cachedSessionId) {
        console.warn('[loadChatDataFromStorage] cached session has no conversation id — discarding it');
        await window.HOWLStorage.deleteChatData('howl-chat-current-chat').catch(() => {});
        localStorage.removeItem('howl-chat-current-session');
        window.sidebarState.resetToDefaultState();
        await updateContextDisplay();
        return;
    }
    // Fail closed. If the authority cannot confirm the session — no bridge yet,
    // or a failed call — the cache is discarded. Restoring on an unconfirmed
    // cache is exactly how a shadow copy mints conversations the backend never
    // created.
    let confirmed = false;
    if (window.go?.main?.App?.GetChatList) {
        try {
            const chats = await window.go.main.App.GetChatList();
            confirmed = (chats || []).some(c => c.id === cachedSessionId && (!c.state || c.state === 'present'));
        } catch (err) {
            console.warn('[loadChatDataFromStorage] backend validation failed:', err);
        }
    }
    if (!confirmed) {
        console.warn('[loadChatDataFromStorage] cached session', cachedSessionId, 'is not a present conversation — discarding cache');
        await window.HOWLStorage.deleteChatData('howl-chat-current-chat').catch(() => {});
        localStorage.removeItem('howl-chat-current-session');
        window.sidebarState.resetToDefaultState();
        await updateContextDisplay();
        return;
    }

    if (window.chatBridge) {
        await window.chatBridge.hydrateFromBackend(savedData);
    } else {
        window.sidebarState.currentChatData = savedData;
        await updateContextDisplay();
    }

    // Ensure backend state is in sync when loading from frontend storage
    await syncAllDataToBackend();
}

// ---------------------------------------------------------------------------
// Chat History Reconstruction
// ---------------------------------------------------------------------------

/**
 * Clears the chat message container and rebuilds it from a history array.
 * @param {Array} history
 */
function rebuildChatFromHistory(history) {
    const chatMessages = document.querySelector('.chat-messages');
    if (!chatMessages) return;

    chatMessages.innerHTML = '';

    console.log('[rebuildChatFromHistory] rebuilding', history.length, 'messages');
    history.forEach((msg, index) => {
        const role        = (msg.Role   || msg.role   || '').toLowerCase();
        const isUser      = msg.isUser === true || role === 'user';
        const isSystem    = role === 'system';
        const sender      = msg.Sender || msg.sender || (isUser ? 'User' : 'Assistant');
        const text        = msg.Text   || msg.text   || '';
        const characterId = msg.characterId || msg.CharacterId || null;

        if (index < 3 || index >= history.length - 2) {
            console.log(`[rebuildChatFromHistory] msg ${index}: role=${role} isUser=${isUser} sender=${sender} charId=${characterId} text=${text.slice(0,40)}`);
        }

        if (isSystem) {
            if (typeof addSystemMessage === 'function') addSystemMessage(text);
        } else if (isUser) {
            if (typeof addUserMessage === 'function') addUserMessage(text);
        } else {
            if (typeof addAssistantMessage === 'function') addAssistantMessage(sender, text, characterId);
        }

        // Preserve original metadata on the last inserted DOM element
        const chatHistory = document.querySelector('.chat-messages');
        const lastMsg = chatHistory ? chatHistory.lastElementChild : null;
        if (lastMsg) {
            if (msg.id || msg.ID) lastMsg.dataset.id = msg.id || msg.ID;
            if (msg.timestamp || msg.Timestamp) lastMsg.dataset.timestamp = msg.timestamp || msg.Timestamp;
            lastMsg.dataset.characterId = characterId || '';
            lastMsg.dataset.hasImage = (msg.hasImage || msg.HasImage) ? 'true' : '';
        }
    });
    console.log('[rebuildChatFromHistory] done');

    chatMessages.scrollTop = chatMessages.scrollHeight;
}

// ---------------------------------------------------------------------------
// Panel Orchestration
// ---------------------------------------------------------------------------

/**
 * Refreshes all sidebar panels with the current chat data.
 */
async function updateAllChatPanels() {
    // updateContextDisplay subsumes all individual panel updates
    await updateContextDisplay();
}

/**
 * Updates the session info display element with metadata.
 * @param {Object} chatMeta
 */
function updateSessionInfoDisplay(chatMeta) {
    const sessionInfo = document.querySelector('.session-info');
    if (!sessionInfo) return;

    const name           = chatMeta.name           || chatMeta.Name           || 'New Chat';
    const worldName      = chatMeta.worldName      || chatMeta.WorldName      || '';
    const scenarioName   = chatMeta.scenarioName   || chatMeta.ScenarioName   || '';
    const characterNames = chatMeta.characterNames || chatMeta.CharacterNames || [];

    sessionInfo.innerHTML = `
        <div class="session-name">${name}</div>
        <div class="session-details">
            ${worldName      ? `World: ${worldName}`                     : ''}
            ${scenarioName   ? `Scenario: ${scenarioName}`               : ''}
            ${characterNames.length > 0 ? `Characters: ${characterNames.join(', ')}` : ''}
        </div>
    `;
}

// ---------------------------------------------------------------------------
// Persistence
// ---------------------------------------------------------------------------

/**
 * Saves a minimized version of currentChatData to HOWLStorage.
 * Falls back to essential-only data if quota is exceeded.
 */
async function saveChatData() {
    const cd = window.sidebarState.currentChatData;
    try {
        window.sidebarState.logChatContract('save:pre-minimalize', cd);

        const minimalChatData = {
            world: cd.world ? {
                id:    cd.world.id,
                title: cd.world.title || cd.world.name,
                name:  cd.world.name  || cd.world.title
            } : null,
            scenario: cd.scenario ? {
                id:    cd.scenario.id,
                title: cd.scenario.title || cd.scenario.name,
                name:  cd.scenario.name  || cd.scenario.title
            } : null,
            lorebooks:  (cd.lorebooks  || []).map(l => ({ id: l.id, name: l.name || l.title, title: l.title || l.name, description: l.description || '', enabled: l.enabled !== false })),
            characters: (cd.characters || []).map(c => ({
                id: c.id,
                name: c.name,
                status: c.status,
                active: c.active,
                stepOutReason: c.stepOutReason,
                stepOutMode: c.stepOutMode
            })),
            chatName: cd.chatName || 'New Chat',
            // Tracking arrays — chat is the sole source of truth
            activeCharacterIds:  cd.activeCharacterIds  || [],
            removedCharacterIds: cd.removedCharacterIds || [],
            addedCharacterIds:   cd.addedCharacterIds   || [],
            activeLorebookIds:   cd.activeLorebookIds   || [],
            removedLorebookIds:  cd.removedLorebookIds  || [],
            addedLorebookIds:    cd.addedLorebookIds    || []
        };

        const dataString = JSON.stringify(minimalChatData);
        console.log('Saving chat data:', {
            size:       `${(dataString.length / 1024).toFixed(2)} KB`,
            world:      minimalChatData.world     ? minimalChatData.world.title     : 'None',
            scenario:   minimalChatData.scenario  ? minimalChatData.scenario.title  : 'None',
            lorebooks:  minimalChatData.lorebooks.length,
            characters: minimalChatData.characters.length
        });

        await window.HOWLStorage.setChatData('howl-chat-current-chat', JSON.parse(dataString));
        console.log('Chat data saved successfully to storage');

        // Sync full state to backend via SaveChatData (carries tracking arrays)
        // Merge cd into bridgeState so we don't drop characters that were recovered
        // from localStorage during hydrateFromBackend but aren't in cd yet.
        if (window.chatBridge?.getState && window.chatBridge?.collectForSave) {
            const bridgeState = window.chatBridge.getState();
            if (bridgeState) {
                const mergeById = (a, b) => {
                    const map = new Map();
                    (a || []).forEach(x => { if (x?.id) map.set(x.id, x); });
                    (b || []).forEach(x => { if (x?.id) map.set(x.id, x); });
                    return Array.from(map.values());
                };
                bridgeState.world      = cd.world      || bridgeState.world;
                bridgeState.scenario   = cd.scenario   || bridgeState.scenario;
                bridgeState.characters = mergeById(bridgeState.characters, cd.characters);
                bridgeState.lorebooks  = mergeById(bridgeState.lorebooks,  cd.lorebooks);
                bridgeState.activeCharacterIds  = [...new Set([...(bridgeState.activeCharacterIds  || []), ...(cd.activeCharacterIds  || [])])];
                bridgeState.removedCharacterIds = [...new Set([...(bridgeState.removedCharacterIds || []), ...(cd.removedCharacterIds || [])])];
                bridgeState.addedCharacterIds   = [...new Set([...(bridgeState.addedCharacterIds   || []), ...(cd.addedCharacterIds   || [])])];
                bridgeState.activeLorebookIds   = [...new Set([...(bridgeState.activeLorebookIds   || []), ...(cd.activeLorebookIds   || [])])];
                bridgeState.removedLorebookIds  = [...new Set([...(bridgeState.removedLorebookIds  || []), ...(cd.removedLorebookIds  || [])])];
                bridgeState.addedLorebookIds    = [...new Set([...(bridgeState.addedLorebookIds    || []), ...(cd.addedLorebookIds    || [])])];
            }
        }

        if (window.chatBridge?.collectForSave && window.go?.main?.App?.SaveChatData) {
            const payload = window.chatBridge.collectForSave();
            window.go.main.App.SaveChatData(payload).then(() => {
                console.log('Chat data saved to backend via SaveChatData');
            }).catch(err => {
                console.warn('Backend SaveChatData failed, falling back to SaveChat:', err);
                if (window.go?.main?.App?.SaveChat) {
                    window.go.main.App.SaveChat().catch(() => {});
                }
            });
        } else if (window.go?.main?.App?.SaveChat) {
            window.go.main.App.SaveChat().then(() => {
                console.log('Chat data also saved to backend');
            }).catch(err => {
                console.warn('Backend save failed, localStorage save succeeded:', err);
            });
        }

    } catch (error) {
        if (error.name === 'QuotaExceededError') {
            console.error('localStorage quota exceeded. Attempting minimal save...');
            ['howl-chat-current-chat', 'howl-chat-scenarios-cache', 'howl-chat-lorebooks-cache'].forEach(key => {
                try { localStorage.removeItem(key); } catch (e) { console.warn('Failed to clear', key, e); }
            });

            try {
                const essentialData = {
                    world:    cd.world    ? { id: cd.world.id,    title: cd.world.title    || cd.world.name    } : null,
                    scenario: cd.scenario ? { id: cd.scenario.id, title: cd.scenario.title || cd.scenario.name } : null,
                    chatName: cd.chatName || 'New Chat',
                    // Tracking arrays
                    activeCharacterIds:  cd.activeCharacterIds  || [],
                    removedCharacterIds: cd.removedCharacterIds || [],
                    addedCharacterIds:   cd.addedCharacterIds   || [],
                    activeLorebookIds:   cd.activeLorebookIds   || [],
                    removedLorebookIds:  cd.removedLorebookIds  || [],
                    addedLorebookIds:    cd.addedLorebookIds    || []
                };
                await window.HOWLStorage.setChatData('howl-chat-current-chat', essentialData);
                console.log('Essential chat data saved after cleanup');
            } catch (retryError) {
                console.error('Even minimal save failed:', retryError);
                alert('Storage quota exceeded. Please clear browser data or use a different browser.');
            }
        } else {
            console.error('Failed to save chat data:', error);
            throw error;
        }
    }
}

// ---------------------------------------------------------------------------
// Backend Sync Helpers
// ---------------------------------------------------------------------------

async function syncAllDataToBackend() {
    const cd = window.sidebarState.currentChatData;
    if (cd.world)                              await syncWorldDataToBackend(cd.world);
    if (cd.scenario)                           await syncScenarioDataToBackend(cd.scenario);
    if (cd.characters && cd.characters.length) await syncCharactersDataToBackend(cd.characters);
    if (cd.lorebooks  && cd.lorebooks.length)  await syncLorebooksDataToBackend(cd.lorebooks);
}

async function syncWorldDataToBackend(world) {
    if (window.go && window.go.main && window.go.main.App && window.go.main.App.AddWorldToChat) {
        try {
            await window.go.main.App.AddWorldToChat(world);
            console.log(`Synced world to backend: ${world.title || world.name}`);
        } catch (error) {
            console.warn('Backend world sync failed:', error);
        }
    }
}

async function syncScenarioDataToBackend(scenario) {
    if (window.go?.main?.App?.AddScenarioToChat) {
        try {
            await window.go.main.App.AddScenarioToChat(scenario);
            console.log(`Synced scenario to backend: ${scenario.title || scenario.name}`);
        } catch (error) {
            console.warn('Backend scenario sync failed:', error);
        }
    }
}

async function syncCharactersDataToBackend(characters) {
    if (window.go && window.go.main && window.go.main.App && window.go.main.App.AddCharacterToChat) {
        try {
            // Clear stale characters from previous session before re-syncing
            if (window.go.main.App.ResetSessionCharacters) {
                await window.go.main.App.ResetSessionCharacters();
            }
            for (const character of characters) {
                await window.go.main.App.AddCharacterToChat(character);
            }
            console.log(`Synced ${characters.length} characters to backend`);
        } catch (error) {
            console.warn('Backend character sync failed:', error);
        }
    }
}

async function syncLorebooksDataToBackend(lorebooks) {
    if (window.go && window.go.main && window.go.main.App && window.go.main.App.UpdateLorebookContext) {
        try {
            const enabledIds = (lorebooks || []).filter(l => l.enabled !== false).map(l => l.id);
            await window.go.main.App.UpdateLorebookContext(enabledIds);
            console.log(`Synced ${enabledIds.length} active lorebooks to backend`);
        } catch (error) {
            console.warn('Backend lorebook sync failed:', error);
        }
    }
}

// ---------------------------------------------------------------------------
// Delete Chat Session
// ---------------------------------------------------------------------------

/**
 * Permanently delete a chat session from disk.
 * Shows a "Deletion is permanent" confirmation first.
 * If the deleted chat is the currently active session, clears bridge state.
 * Also updates the File > Load session list UI if present.
 * @param {string} chatId
 * @param {string} [chatName]
 * @param {Event|function} [onDeletedOrEvent]  - event object from inline onclick, or callback
 */
function deleteChatSession(chatId, chatName, onDeletedOrEvent) {
    if (!chatId) return;

    // Detect if third arg is a click event (from session-list delete button)
    let callback = onDeletedOrEvent;
    if (onDeletedOrEvent && typeof onDeletedOrEvent.stopPropagation === 'function') {
        onDeletedOrEvent.stopPropagation();
        onDeletedOrEvent.preventDefault();
        callback = null;
    }

    const label = chatName || chatId;

    window.DeleteConfirm.show({
        title: `Delete "${label}"?`,
        body:  'This chat\'s history and automated lore entries will be permanently deleted. World, scenario, and character data are preserved.',
        onConfirm: async () => {
            try {
                if (window.go?.main?.App?.DeleteChat) {
                    await window.go.main.App.DeleteChat(chatId);
                }

                // If this was the active session, reset bridge state, stop autosave, and wipe cache
                const activeId = window.chatBridge?.getState()?.sessionId;
                if (activeId === chatId && window.chatBridge) {
                    await window.chatBridge.clearForNewChat('New Chat');
                }
                if (window.dataSync) {
                    window.dataSync.stopAutosave();
                }

                // Wipe frontend caches so the deleted chat cannot resurface as a ghost
                if (window.HOWLStorage) {
                    await window.HOWLStorage.deleteChatData('howl-chat-current-chat').catch(() => {});
                }
                localStorage.removeItem('howl-chat-current-session');

                // Update File > Load session list if it exists
                if (typeof savedSessions !== 'undefined' && Array.isArray(savedSessions)) {
                    savedSessions = savedSessions.filter(s => s.id !== chatId);
                }
                if (typeof selectedSessionId !== 'undefined' && selectedSessionId === chatId) {
                    selectedSessionId = null;
                    const loadButton = document.querySelector('.session-list-load');
                    if (loadButton) loadButton.disabled = true;
                }
                if (typeof renderSessionList === 'function') {
                    renderSessionList();
                }

                // Force re-fetch from backend to ensure sync
                if (typeof loadSavedSessions === 'function') {
                    loadSavedSessions().catch(e => console.warn('Post-delete refresh failed:', e));
                }

                if (typeof callback === 'function') callback(chatId);

                if (typeof showNotification === 'function') {
                    showNotification(`Chat deleted: ${label}`, 'success');
                }
            } catch (err) {
                console.error('Failed to delete chat session:', err);
                if (typeof showNotification === 'function') {
                    showNotification('Failed to delete chat: ' + (err.message || err), 'error');
                }
            }
        }
    });
}

// ---------------------------------------------------------------------------
// Global Exports
// ---------------------------------------------------------------------------

window.sidebarSession = {
    loadChatData,
    loadChatSessionById,
    loadChatDataFromStorage,
    rebuildChatFromHistory,
    updateAllChatPanels,
    updateSessionInfoDisplay,
    saveChatData,
    syncAllDataToBackend,
    syncWorldDataToBackend,
    syncScenarioDataToBackend,
    syncCharactersDataToBackend,
    syncLorebooksDataToBackend,
    deleteChatSession
};
