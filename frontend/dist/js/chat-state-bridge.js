/**
 * Chat State Bridge
 *
 * Single source of truth that keeps contextDisplay and sidebarState
 * in sync.  Every load, new-chat, and save must go through here.
 *
 * Public API (window.chatBridge):
 *   hydrateFromBackend(chatData)  — after LoadChat; populates all panels
 *   clearForNewChat(chatName)     — after NewChat; wipes all panels
 *   collectForSave()              — gathers full save payload
 *   refreshAllPanels()            — re-renders every sidebar panel
 *
 * @module chat-state-bridge
 * @version 1.0.0
 */

// ---------------------------------------------------------------------------
// Internal canonical state
// ---------------------------------------------------------------------------

function _defaultState(chatName) {
    return {
        chatName:            chatName || 'New Chat',
        sessionId:           null,
        world:               null,
        scenario:            null,
        characters:          [],
        lorebooks:           [],
        // Tracking arrays — chat is the sole source of truth
        activeCharacterIds:  [],
        removedCharacterIds: [],
        addedCharacterIds:   [],
        activeLorebookIds:   [],
        removedLorebookIds:  [],
        addedLorebookIds:    []
    };
}

if (!window._chatBridgeState) window._chatBridgeState = _defaultState();
var _state = window._chatBridgeState;

// ---------------------------------------------------------------------------
// Image resolution (async — handles HOWLStorage fallback)
// ---------------------------------------------------------------------------

async function _resolveImage(imageRef, storeType) {
    if (!imageRef) return null;
    if (typeof imageRef === 'string' && imageRef) return imageRef;
    if (typeof imageRef === 'object') {
        const quick = imageRef.data || imageRef.url || imageRef.path || imageRef.thumbnail;
        if (quick) return quick;

        // Character images live in a separate CharacterStorage DB, not HOWLStorage
        if (storeType === 'character' && imageRef.id && window.CharacterStorage) {
            try {
                // imageRef.id is the image entry id; load all images for this character id
                // The character id may be on the parent object, so this path is best-effort
                // A more reliable path is handled in hydrateFromBackend via CharacterStorage.loadAll
            } catch (_) {}
        }

        if (imageRef.id && window.HOWLStorage) {
            try {
                const stored = await window.HOWLStorage.getImage(imageRef.id, storeType);
                if (stored && stored.data) return stored.data;
            } catch (_) {}
        }
    }
    return null;
}

/**
 * Resolve a multi-image array, returning the same array with .data populated.
 * For {id, stored:true} refs, fetches from IndexedDB. For objects with .data, passes through.
 */
async function _resolveImages(images, storeType) {
    if (!Array.isArray(images) || images.length === 0) return images;
    const resolved = [];
    for (const img of images) {
        if (!img) { resolved.push(img); continue; }
        if (typeof img === 'string') { resolved.push({ data: img }); continue; }
        if (img.data || img.url || img.path || img.thumbnail) {
            resolved.push(img);
            continue;
        }
        // {id, stored:true} reference — try IndexedDB
        if (img.id && window.HOWLStorage) {
            try {
                const stored = await window.HOWLStorage.getImage(img.id, storeType);
                if (stored && stored.data) {
                    resolved.push({ ...img, data: stored.data });
                    continue;
                }
            } catch (_) {}
        }
        resolved.push(img);
    }
    return resolved;
}

// ---------------------------------------------------------------------------
// Hydrate from backend payload  (called after LoadChat)
// ---------------------------------------------------------------------------

async function hydrateFromBackend(chatData) {
    if (!chatData) return;

    const world     = chatData.World     || chatData.world     || null;
    const scenario  = chatData.Scenario  || chatData.scenario  || null;
    let   characters= chatData.Characters|| chatData.characters|| [];
    let   lorebooks = chatData.Lorebooks || chatData.lorebooks || [];
    const meta      = chatData.Meta      || chatData.meta      || {};
    const sessionId = meta.ID || meta.id || chatData.sessionId || null;
    const chatName  = meta.Name || meta.name || 'New Chat';

    // -----------------------------------------------------------------------
    // SOURCE OF TRUTH: tracking arrays in meta
    // Final lists = (active - removed) + added
    // -----------------------------------------------------------------------
    const activeChars  = meta.activeCharacters  || meta.ActiveCharacters  || [];
    const removedChars = meta.removedCharacters || meta.RemovedCharacters || [];
    const addedChars   = meta.addedCharacters   || meta.AddedCharacters   || [];

    const activeBooks  = meta.activeLorebooks  || meta.ActiveLorebooks  || [];
    const removedBooks = meta.removedLorebooks || meta.RemovedLorebooks || [];
    const addedBooks   = meta.addedLorebooks   || meta.AddedLorebooks   || [];

    console.log('[hydrateFromBackend] raw characters:', characters.length, 'raw lorebooks:', lorebooks.length);
    console.log('[hydrateFromBackend] activeChars:', activeChars.length, 'removedChars:', removedChars.length, 'addedChars:', addedChars.length);
    console.log('[hydrateFromBackend] activeBooks:', activeBooks.length, 'removedBooks:', removedBooks.length, 'addedBooks:', addedBooks.length);

    const removedCharSet = new Set(removedChars);
    const finalCharIds   = new Set([
        ...activeChars.filter(id => !removedCharSet.has(id)),
        ...addedChars
    ]);

    const removedBookSet = new Set(removedBooks);
    const finalBookIds   = new Set([
        ...activeBooks.filter(id => !removedBookSet.has(id)),
        ...addedBooks
    ]);

    console.log('[hydrateFromBackend] finalCharIds:', [...finalCharIds], 'finalBookIds:', [...finalBookIds]);

    // Ensure every character/lorebook in the raw data is accounted for in the tracking.
    // If an entity was saved but never made it into the tracking arrays
    // (e.g. added via world/scenario load before tracking existed), add it now.
    const rawCharIds = new Set((characters || []).map(c => c.id).filter(Boolean));
    for (const id of rawCharIds) {
        if (!finalCharIds.has(id)) {
            finalCharIds.add(id);
            if (!activeChars.includes(id) && !addedChars.includes(id)) {
                addedChars.push(id);
            }
        }
    }
    const rawBookIds = new Set((lorebooks || []).map(l => l.id).filter(Boolean));
    for (const id of rawBookIds) {
        if (!finalBookIds.has(id)) {
            finalBookIds.add(id);
            if (!activeBooks.includes(id) && !addedBooks.includes(id)) {
                addedBooks.push(id);
            }
        }
    }

    // Filter loaded characters/lorebooks to only those in the computed sets
    characters = characters.filter(c => finalCharIds.has(c.id));
    lorebooks  = lorebooks.filter(l => finalBookIds.has(l.id));
    console.log('[hydrateFromBackend] after filter characters:', characters.length, 'lorebooks:', lorebooks.length);

    // If the backend didn't return a character that's supposed to be active,
    // look it up from localStorage (handles characters added after last save)
    const loadedCharIds = new Set(characters.map(c => c.id));
    // Only look up characters that are explicitly ADDED (not stale active ones).
    // Stale active IDs that aren't loaded should NOT be recovered — they were
    // removed or the save is old.  Only added-after-save IDs need recovery.
    const addedCharSet = new Set(addedChars);
    const missingCharIds = [...finalCharIds].filter(id => !loadedCharIds.has(id) && addedCharSet.has(id));
    console.log('[hydrateFromBackend] missingCharIds (added-only):', missingCharIds);
    if (missingCharIds.length > 0) {
        try {
            const allChars = (window.CharacterStorage?.all?.() || []);
            for (const id of missingCharIds) {
                const found = allChars.find(c => c.id === id);
                if (found) {
                    console.log('[hydrateFromBackend] adding missing character from localStorage:', id, found.name);
                    characters.push(_enrichCharacterFromRepository(found));
                } else {
                    console.log('[hydrateFromBackend] character not found in localStorage:', id);
                }
            }
        } catch (_) {}
    }

    // Same for lorebooks
    const loadedBookIds = new Set(lorebooks.map(l => l.id));
    const addedBookSet = new Set(addedBooks);
    const missingBookIds = [...finalBookIds].filter(id => !loadedBookIds.has(id) && addedBookSet.has(id));
    console.log('[hydrateFromBackend] missingBookIds (added-only):', missingBookIds);
    if (missingBookIds.length > 0) {
        try {
            const allBooks = (window.LorebookRepo ? window.LorebookRepo.all() : []);
            for (const id of missingBookIds) {
                const found = allBooks.find(b => b.id === id);
                if (found) {
                    console.log('[hydrateFromBackend] adding missing lorebook from localStorage:', id, found.name || found.title);
                    lorebooks.push(found);
                } else {
                    console.log('[hydrateFromBackend] lorebook not found in localStorage:', id);
                }
            }
        } catch (_) {}
    }

    console.log('[hydrateFromBackend] final characters:', characters.length, 'final lorebooks:', lorebooks.length);

    _state = {
        chatName,
        sessionId,
        world,
        scenario,
        characters,
        lorebooks,
        activeCharacterIds:  activeChars,
        removedCharacterIds: removedChars,
        addedCharacterIds:   addedChars,
        activeLorebookIds:   activeBooks,
        removedLorebookIds:  removedBooks,
        addedLorebookIds:    addedBooks
    };
    window._chatBridgeState = _state;

    // Resolve images asynchronously so displays fill in correctly
    if (world    && world.image)    world._resolvedImage    = await _resolveImage(world.image,    'world');
    if (scenario && scenario.image) scenario._resolvedImage = await _resolveImage(scenario.image, 'scenario');

    // Resolve multi-image arrays
    if (world    && world.images    && Array.isArray(world.images))    world.images    = await _resolveImages(world.images,    'world');
    if (scenario && scenario.images && Array.isArray(scenario.images)) scenario.images = await _resolveImages(scenario.images, 'scenario');

    // Enrich from localStorage — recover title/image the backend doesn't store
    if (_state.world)    _enrichWorldFromLocalStorage(_state.world);
    if (_state.scenario) _enrichScenarioFromLocalStorage(_state.scenario);
    _state.characters = _state.characters.map(c => _enrichCharacterFromRepository(c));
    _state.lorebooks  = _enrichLorebooksFromLocalStorage(_state.lorebooks);

    // Resolve images async for anything enrichment added
    if (_state.world && _state.world.image && !_state.world._resolvedImage) {
        _state.world._resolvedImage = await _resolveImage(_state.world.image, 'world');
    }
    if (_state.scenario && _state.scenario.image && !_state.scenario._resolvedImage) {
        _state.scenario._resolvedImage = await _resolveImage(_state.scenario.image, 'scenario');
    }
    if (_state.world && _state.world.images && Array.isArray(_state.world.images)) {
        _state.world.images = await _resolveImages(_state.world.images, 'world');
    }
    if (_state.scenario && _state.scenario.images && Array.isArray(_state.scenario.images)) {
        _state.scenario.images = await _resolveImages(_state.scenario.images, 'scenario');
    }
    for (const c of _state.characters) {
        if (c.image && !c._resolvedImage) {
            c._resolvedImage = await _resolveImage(c.image, 'character');
        }
        if (!c._resolvedImage && c.id && window.CharacterStorage) {
            try {
                const images = await window.CharacterStorage.loadImages(c.id);
                const thumb = images.find(img => img.isThumbnail) || images[0];
                if (thumb && thumb.data) {
                    c._resolvedImage = thumb.data;
                    c.image = thumb.data;
                }
            } catch (_) {}
        }
        if (!c._resolvedImage && c.howl?.images?.length > 0) {
            const thumb = c.howl.images.find(img => img.isThumbnail) || c.howl.images[0];
            if (thumb && thumb.data) {
                c._resolvedImage = thumb.data;
                c.image = thumb.data;
            }
        }
    }

    _pushToContextDisplay();
    _pushToSidebarState();

    // Re-sync enriched world/scenario/characters to Go backend so a.world/a.scenario
    // and a.characters are populated for buildPrompt and SaveChat.
    if (_state.world && window.go?.main?.App?.AddWorldToChat) {
        window.go.main.App.AddWorldToChat(_state.world).catch(() => {});
    }
    if (_state.scenario && window.go?.main?.App?.AddScenarioToChat) {
        window.go.main.App.AddScenarioToChat(_state.scenario).catch(() => {});
    }
    if (_state.characters.length > 0 && window.go?.main?.App?.AddCharacterToChat) {
        _state.characters.forEach(c => {
            window.go.main.App.AddCharacterToChat(c).catch(() => {});
        });
    }

    // Rebuild chat history DOM
    const history = chatData.History || chatData.history || [];
    if (history.length > 0 && typeof rebuildChatFromHistory === 'function') {
        rebuildChatFromHistory(history);
    } else {
        const msgs = document.querySelector('.chat-messages');
        if (msgs) msgs.innerHTML = '';
    }

    await refreshAllPanels();

    if (typeof updateSessionInfoDisplay === 'function') {
        updateSessionInfoDisplay(meta);
    }
}

// ---------------------------------------------------------------------------
// Clear for new chat  (called after NewChat backend response)
// ---------------------------------------------------------------------------

async function clearForNewChat(chatName, sessionId) {
    _state = _defaultState(chatName || 'New Chat');
    _state.sessionId = sessionId || null;
    _state.removedCharacterIds = []; // fresh chat = fresh removals
    window._chatBridgeState = _state;

    _pushToContextDisplay();
    _pushToSidebarState();

    const msgs = document.querySelector('.chat-messages');
    if (msgs) msgs.innerHTML = '';

    await refreshAllPanels();

    const preset = window.settingsState?.current?.general?.defaultSampler;
    if (preset && typeof applySamplerPreset === 'function') {
        applySamplerPreset(preset);
    }
}

// ---------------------------------------------------------------------------
// Apply a context change  (world select, scenario swap, character add, etc.)
// ---------------------------------------------------------------------------

async function applyWorldChange(world) {
    _state.world = world || null;
    if (world && world.image) world._resolvedImage = await _resolveImage(world.image, 'world');
    if (world && world.images && Array.isArray(world.images)) {
        world.images = await _resolveImages(world.images, 'world');
    }
    _push(); await refreshAllPanels();
}

async function applyScenarioChange(scenario) {
    const prev = _state.scenario;
    _state.scenario = scenario || null;

    if (scenario && scenario.image) {
        scenario._resolvedImage = await _resolveImage(scenario.image, 'scenario');
    }
    if (scenario && scenario.images && Array.isArray(scenario.images)) {
        scenario.images = await _resolveImages(scenario.images, 'scenario');
    }

    // Scenarios are metadata only — they NEVER override chat state.
    // Characters and lorebooks remain exactly as the chat defines them.
    if (prev && scenario && prev.id !== scenario.id) {
        _emitScenarioSwapEvent(prev, scenario);
    }

    _push(); await refreshAllPanels();
    _syncScenarioAndCharactersToBackend();
}

async function applyCharacterAdd(character) {
    console.log('[TRACE:applyCharacterAdd] id:', character?.id, 'name:', character?.name, 'existing:', _state.characters.some(c => c.id === character?.id));
    if (!character || _state.characters.some(c => c.id === character.id)) return;

    // Update tracking arrays:
    // If the character was previously removed, clear the removal.
    // Otherwise, if it's not in the active base list, mark it as added.
    const wasRemoved = (_state.removedCharacterIds || []).includes(character.id);
    if (wasRemoved) {
        _state.removedCharacterIds = _state.removedCharacterIds.filter(id => id !== character.id);
    } else if (!(_state.activeCharacterIds || []).includes(character.id)) {
        if (!_state.addedCharacterIds.includes(character.id)) {
            _state.addedCharacterIds.push(character.id);
        }
    }

    _state.characters.push(character);
    _push(); await refreshAllPanels();
    if (window.go?.main?.App?.AddCharacterToChat) {
        window.go.main.App.AddCharacterToChat(character).catch(() => {});
    }
}

async function applyCharacterRemove(characterId) {
    console.log('[TRACE:applyCharacterRemove] id:', characterId, 'beforeCount:', _state.characters.length);

    _state.characters = _state.characters.filter(c => c.id !== characterId);

    // Update tracking arrays:
    // If the character was in the active base list, mark it as removed.
    // If it was in the added list, remove it from there.
    if ((_state.activeCharacterIds || []).includes(characterId)) {
        if (!_state.removedCharacterIds.includes(characterId)) {
            _state.removedCharacterIds.push(characterId);
        }
    }
    _state.addedCharacterIds = (_state.addedCharacterIds || []).filter(id => id !== characterId);

    _push(); await refreshAllPanels();

    // Sync removal to backend so a.characters stays clean
    if (window.go?.main?.App?.RemoveCharacterFromChat) {
        window.go.main.App.RemoveCharacterFromChat(characterId).catch(() => {});
    }
    if (typeof saveChatData === 'function') saveChatData();
}

/**
 * Batch update characters: remove a set of IDs, then add a set of characters.
 * Only refreshes panels once and syncs to backend once — eliminates flicker.
 */
async function applyCharacterBatchUpdate(removeIds, addCharacters) {
    const removeSet = new Set(removeIds || []);
    const addMap    = new Map((addCharacters || []).map(c => [c.id, c]));

    // Apply all state mutations first
    _state.characters = _state.characters.filter(c => !removeSet.has(c.id) && !addMap.has(c.id));
    for (const c of addMap.values()) {
        _state.characters.push(c);
    }

    // Update tracking arrays
    for (const id of removeSet) {
        if ((_state.activeCharacterIds || []).includes(id)) {
            if (!_state.removedCharacterIds.includes(id)) _state.removedCharacterIds.push(id);
        }
        _state.addedCharacterIds = (_state.addedCharacterIds || []).filter(aid => aid !== id);
    }
    for (const c of addMap.values()) {
        const wasRemoved = (_state.removedCharacterIds || []).includes(c.id);
        if (wasRemoved) {
            _state.removedCharacterIds = _state.removedCharacterIds.filter(id => id !== c.id);
        } else if (!(_state.activeCharacterIds || []).includes(c.id)) {
            if (!_state.addedCharacterIds.includes(c.id)) _state.addedCharacterIds.push(c.id);
        }
    }

    console.log('[TRACE:applyCharacterBatchUpdate] removed:', removeIds?.length || 0, 'added:', addCharacters?.length || 0, 'final:', _state.characters.length);
    _push(); await refreshAllPanels();

    // Single backend sync: reset then add all final characters
    if (window.go?.main?.App?.ResetSessionCharacters) {
        await window.go.main.App.ResetSessionCharacters().catch(() => {});
    }
    if (window.go?.main?.App?.AddCharacterToChat) {
        for (const c of _state.characters) {
            await window.go.main.App.AddCharacterToChat(c).catch(() => {});
        }
    }
    if (typeof saveChatData === 'function') saveChatData();
}

function setLorebooks(lorebooks) {
    _state.lorebooks = lorebooks || [];
    _push();
    console.log('[TRACE:setLorebooks] count:', _state.lorebooks.length);
}

async function applyLorebookToggle(lorebook) {
    if (!lorebook) return;
    const idx = _state.lorebooks.findIndex(l => l.id === lorebook.id);
    if (idx === -1) {
        // Adding a lorebook
        _state.lorebooks.push(lorebook);
        const wasRemoved = (_state.removedLorebookIds || []).includes(lorebook.id);
        if (wasRemoved) {
            _state.removedLorebookIds = _state.removedLorebookIds.filter(id => id !== lorebook.id);
        } else if (!(_state.activeLorebookIds || []).includes(lorebook.id)) {
            if (!_state.addedLorebookIds.includes(lorebook.id)) {
                _state.addedLorebookIds.push(lorebook.id);
            }
        }
    } else {
        // Removing a lorebook
        _state.lorebooks.splice(idx, 1);
        if ((_state.activeLorebookIds || []).includes(lorebook.id)) {
            if (!_state.removedLorebookIds.includes(lorebook.id)) {
                _state.removedLorebookIds.push(lorebook.id);
            }
        }
        _state.addedLorebookIds = (_state.addedLorebookIds || []).filter(id => id !== lorebook.id);
    }
    _push(); await refreshAllPanels();
}

// ---------------------------------------------------------------------------
// Collect full save payload
// ---------------------------------------------------------------------------

function collectForSave() {
    const msgs = document.querySelector('.chat-messages');
    const history = msgs ? Array.from(msgs.querySelectorAll('.chat-message')).map((el, i) => {
        const sender = el.querySelector('.chat-sender')?.textContent?.trim() || '';
        const text   = el.querySelector('.chat-text')?.textContent?.trim()   || '';
        // Role comes from the CSS class, never from the sender label. Inferring it
        // from text labeled every player message "assistant" because the default
        // display name is "You", which does not contain the word "user".
        const role   = el.classList.contains('user')      ? 'user'
                     : el.classList.contains('assistant') ? 'assistant'
                     : 'system';
        return {
            id: i.toString(), sender, text, role,
            characterId: el.dataset.characterId || '',
            timestamp: el.dataset.timestamp || new Date().toISOString(),
            hasImage: el.dataset.hasImage === 'true'
        };
    }) : [];

    return {
        meta: {
            id:             _state.sessionId || Date.now().toString(),
            name:           _state.chatName,
            updatedAt:      new Date().toISOString(),
            worldName:      _state.world?.title    || _state.world?.name    || '',
            scenarioName:   _state.scenario?.title || _state.scenario?.name || '',
            characterNames: _state.characters.map(c => c.name || '').filter(Boolean),
            lorebookNames:  _state.lorebooks.map(l => l.name || l.title || '').filter(Boolean),
            messageCount:   history.length,
            // Source-of-truth tracking arrays
            worldId:           _state.world?.id    || _state.world?.ID    || '',
            scenarioId:        _state.scenario?.id || _state.scenario?.ID || '',
            activeCharacters:  _state.activeCharacterIds  || [],
            removedCharacters: _state.removedCharacterIds || [],
            addedCharacters:   _state.addedCharacterIds   || [],
            activeLorebooks:   _state.activeLorebookIds   || [],
            removedLorebooks:  _state.removedLorebookIds  || [],
            addedLorebooks:    _state.addedLorebookIds    || []
        },
        history,
        world:      _state.world,
        scenario:   _state.scenario,
        characters: _state.characters,
        lorebooks:  _state.lorebooks
    };
}

// ---------------------------------------------------------------------------
// Refresh all panels
// ---------------------------------------------------------------------------

async function refreshAllPanels() {
    // contextDisplay.updateContextDisplay subsumes all individual panel updates
    try { if (typeof contextDisplay !== 'undefined') await contextDisplay.updateContextDisplay(); } catch (_) {}
}

// ---------------------------------------------------------------------------
// Getters
// ---------------------------------------------------------------------------

function getState() { return _state; }

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

function _push() {
    _pushToContextDisplay();
    _pushToSidebarState();
}

function _pushToSidebarState() {
    if (!window.sidebarState) return;
    window.sidebarState.currentChatData = {
        world:               _state.world,
        scenario:            _state.scenario,
        characters:          _state.characters,
        lorebooks:           _state.lorebooks,
        sessionId:           _state.sessionId,
        chatName:            _state.chatName,
        // Tracking arrays — chat is the sole source of truth
        activeCharacterIds:  _state.activeCharacterIds,
        removedCharacterIds: _state.removedCharacterIds,
        addedCharacterIds:   _state.addedCharacterIds,
        activeLorebookIds:   _state.activeLorebookIds,
        removedLorebookIds:  _state.removedLorebookIds,
        addedLorebookIds:    _state.addedLorebookIds
    };
}

function _pushToContextDisplay() {
    if (typeof contextDisplay === 'undefined') return;
    contextDisplay.currentChatData = {
        world:               _state.world,
        scenario:            _state.scenario,
        characters:          _state.characters,
        lorebooks:           _state.lorebooks,
        sessionId:           _state.sessionId,
        chatName:            _state.chatName,
        // Tracking arrays — chat is the sole source of truth
        activeCharacterIds:  _state.activeCharacterIds,
        removedCharacterIds: _state.removedCharacterIds,
        addedCharacterIds:   _state.addedCharacterIds,
        activeLorebookIds:   _state.activeLorebookIds,
        removedLorebookIds:  _state.removedLorebookIds,
        addedLorebookIds:    _state.addedLorebookIds
    };
    console.log('[_pushToContextDisplay] lorebooks pushed:', (_state.lorebooks || []).length);
}

function _dedupeById(arr) {
    const seen = new Set();
    return arr.filter(item => {
        if (!item) return false;
        const key = item.id || item.name || JSON.stringify(item);
        if (seen.has(key)) return false;
        seen.add(key);
        return true;
    });
}

function _emitScenarioSwapEvent(prev, next) {
    const prevName = prev?.title || prev?.name || 'Unknown';
    const nextName = next?.title || next?.name || 'Unknown';
    const removedChars = (prev.characters || []).filter(c => !(next.characters || []).some(nc => nc.id === c.id));

    if (removedChars.length > 0 && window.go?.main?.App?.AddSystemMessage) {
        const names = removedChars.map(c => c.name).join(', ');
        const msg   = `[System] Scenario changed from "${prevName}" to "${nextName}". Characters removed from scene: ${names}.`;
        window.go.main.App.AddSystemMessage(msg).catch(() => {});
    }
}

// ---------------------------------------------------------------------------
// localStorage enrichment helpers
// ---------------------------------------------------------------------------

function _enrichWorldFromLocalStorage(world) {
    if (!world || !world.id) return;
    try {
        const stored = window.WorldRepository.all();
        const match  = stored.find(w => w.id === world.id);
        if (!match) return;
        if (!world.name  || world.name  === '') world.name  = match.title || match.name  || world.name;
        if (!world.title || world.title === '') world.title = match.title || match.name  || world.title;
        if (!world.description && match.description) world.description = match.description;
        if (!world.details && match.details) world.details = match.details;
        // Always use localStorage image ref (may be {id,stored:true}) — _resolveImage handles IndexedDB lookup
        if (!world.image && match.image) world.image = match.image;
        // If stored record has a better image ref, prefer it
        if (world.image && typeof world.image === 'object' && world.image.stored && match.image) {
            world.image = match.image;
        }
        // Copy multi-image array if present in localStorage but missing on the world object
        if ((!world.images || !Array.isArray(world.images) || world.images.length === 0) && match.images && Array.isArray(match.images) && match.images.length > 0) {
            world.images = match.images;
            world.thumbnailIndex = match.thumbnailIndex || 0;
        }
    } catch (_) {}
}

function _enrichScenarioFromLocalStorage(scenario) {
    if (!scenario || !scenario.id) return;
    try {
        const stored = window.ScenarioRepository.all();
        const match  = stored.find(s => s.id === scenario.id);
        if (!match) return;
        if (!scenario.name  || scenario.name  === '') scenario.name  = match.title || match.name  || scenario.name;
        if (!scenario.title || scenario.title === '') scenario.title = match.title || match.name  || scenario.title;
        if (!scenario.description && match.description) scenario.description = match.description;
        if ((!scenario.characters || scenario.characters.length === 0) && match.characters) scenario.characters = match.characters;
        if ((!scenario.lorebooks  || scenario.lorebooks.length  === 0) && match.lorebooks)  scenario.lorebooks  = match.lorebooks;
        if ((!scenario.Lorebooks || scenario.Lorebooks.length === 0) && match.Lorebooks) scenario.Lorebooks = match.Lorebooks;
        if (!scenario.image && match.image) scenario.image = match.image;
        if (scenario.image && typeof scenario.image === 'object' && scenario.image.stored && match.image) {
            scenario.image = match.image;
        }
        // Copy multi-image array if present in localStorage but missing on the scenario object
        if ((!scenario.images || !Array.isArray(scenario.images) || scenario.images.length === 0) && match.images && Array.isArray(match.images) && match.images.length > 0) {
            scenario.images = match.images;
            scenario.thumbnailIndex = match.thumbnailIndex || 0;
        }
    } catch (_) {}
}

function _enrichCharacterFromRepository(character) {
    if (!character) return character;
    try {
        const stored = (window.CharacterStorage?.all?.() || []);
        const match  = character.id
            ? stored.find(c => c.id === character.id)
            : stored.find(c => c.name && c.name.toLowerCase() === (character.name || '').toLowerCase());
        if (!match) return character;
        return Object.assign({}, match, character,
            { image: character.image || match.image || match.avatar || null });
    } catch (_) {}
    return character;
}

function _enrichLorebooksFromLocalStorage(lorebooks) {
    if (!Array.isArray(lorebooks) || lorebooks.length === 0) return lorebooks || [];
    try {
        const stored = (window.LorebookRepo ? window.LorebookRepo.all() : []);
        return lorebooks.map(lb => {
            if (!lb || !lb.id) return lb;
            const match = stored.find(s => s.id === lb.id);
            if (!match) return lb;
            return Object.assign({}, match, lb,
                { name: lb.name || match.name || match.title || '',
                  title: lb.title || match.title || match.name || '',
                  description: lb.description || match.description || '' });
        });
    } catch (_) {}
    return lorebooks;
}

// ---------------------------------------------------------------------------
// Backend sync helpers
// ---------------------------------------------------------------------------

function _syncScenarioAndCharactersToBackend() {
    if (_state.scenario && window.go?.main?.App?.AddScenarioToChat) {
        window.go.main.App.AddScenarioToChat(_state.scenario).catch(() => {});
    }
    if (_state.characters.length > 0 && window.go?.main?.App?.AddCharacterToChat) {
        _state.characters.forEach(c => {
            window.go.main.App.AddCharacterToChat(c).catch(() => {});
        });
    }
}

// ---------------------------------------------------------------------------
// beforeunload — the conversation save belongs to DataSync, which already
// listens on beforeunload and serializes overlapping writes. A second raw
// SaveChat here raced that save into the conversation boundary and was refused
// as a competing writer, so it is retired rather than kept as a fallback.
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

window.chatBridge = {
    hydrateFromBackend,
    clearForNewChat,
    applyWorldChange,
    applyScenarioChange,
    applyCharacterAdd,
    applyCharacterRemove,
    applyCharacterBatchUpdate,
    setLorebooks,
    applyLorebookToggle,
    collectForSave,
    refreshAllPanels,
    getState
};
