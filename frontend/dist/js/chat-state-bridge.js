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
        chatName:           chatName || 'New Chat',
        sessionId:          null,
        world:              null,
        scenario:           null,
        characters:         [],
        lorebooks:          [],
        removedCharacterIds: []
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
    const lorebooks = chatData.Lorebooks || chatData.lorebooks || [];
    const meta      = chatData.Meta      || chatData.meta      || {};
    const sessionId = meta.ID || meta.id || chatData.sessionId || null;
    const chatName  = meta.Name || meta.name || 'New Chat';

    // Merge removed-character tracking from payload or persisted storage
    let removedIds = chatData.removedCharacterIds || [];
    if (!removedIds.length && window.HOWLStorage) {
        try {
            const saved = await window.HOWLStorage.getChatData('howl-chat-current-chat');
            if (saved && Array.isArray(saved.removedCharacterIds)) {
                removedIds = saved.removedCharacterIds;
            }
        } catch (_) {}
    }

    // Filter out characters the user explicitly removed from this chat
    if (removedIds.length > 0) {
        const removedSet = new Set(removedIds);
        characters = characters.filter(c => !removedSet.has(c.id));
    }

    _state = {
        chatName,
        sessionId,
        world,
        scenario,
        characters,
        lorebooks,
        removedCharacterIds: removedIds
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
    _state.characters = _state.characters.map(c => _enrichCharacterFromLocalStorage(c));
    _state.lorebooks  = _enrichLorebooksFromLocalStorage(_state.lorebooks);

    // Only auto-merge scenario characters if the saved chat has no characters.
    // If a user explicitly removed characters from the chat, the backend will
    // return the authoritative (reduced) list — don't re-inject scenario chars.
    if (_state.scenario && _state.characters.length === 0) {
        const scenarioChars = _state.scenario.characters || _state.scenario.Characters || [];
        for (const sc of scenarioChars) {
            const enriched = _enrichCharacterFromLocalStorage(sc);
            const realId   = enriched.id;
            const already  = _state.characters.some(c =>
                (realId && c.id === realId) ||
                (c.name && enriched.name && c.name.toLowerCase() === enriched.name.toLowerCase())
            );
            if (!already) _state.characters.push(enriched);
        }
    }

    // Auto-merge scenario lorebooks (including character-level lorebooks) if the saved chat has none.
    // This mirrors the logic in applyScenarioChange so hydration and hot-swap behave identically.
    if (_state.scenario && (!_state.lorebooks || _state.lorebooks.length === 0)) {
        const scenarioBooks = _state.scenario.lorebooks || _state.scenario.Lorebooks || [];
        let charLorebooks = [];
        try {
            const scenarioChars = _state.scenario.characters || _state.scenario.Characters || [];
            const allChars = JSON.parse(localStorage.getItem('howl-chat-characters') || '[]');
            charLorebooks = scenarioChars.flatMap(sc => {
                const fullChar = allChars.find(c => c.id === sc.id || c.name === sc.name);
                return fullChar?.howl?.lorebooks || fullChar?.lorebooks || [];
            }).filter(Boolean);
        } catch (_) {}
        const enrichedScenarioBooks = _enrichLorebooksFromLocalStorage(scenarioBooks);
        const enrichedCharBooks     = _enrichLorebooksFromLocalStorage(charLorebooks);
        _state.lorebooks = _dedupeById([...enrichedScenarioBooks, ...enrichedCharBooks]);
    }

    // Resolve images async for anything enrichment added
    if (_state.world && _state.world.image && !_state.world._resolvedImage) {
        _state.world._resolvedImage = await _resolveImage(_state.world.image, 'world');
    }
    if (_state.scenario && _state.scenario.image && !_state.scenario._resolvedImage) {
        _state.scenario._resolvedImage = await _resolveImage(_state.scenario.image, 'scenario');
    }
    // Resolve multi-image arrays after enrichment restored them from localStorage
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
        // CharacterStorage stores full image blobs in a separate IndexedDB.
        // If the backend/localStorage only has stripped refs, load the real images.
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
        // Also try howl.images if present (may have base64 data directly)
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

    // Re-sync enriched world/scenario to Go backend so a.world/a.scenario are populated
    // for buildPrompt and SaveChat (covers the refresh path where LoadChat wasn't called)
    if (_state.world && window.go?.main?.App?.AddWorldToChat) {
        window.go.main.App.AddWorldToChat(_state.world).catch(() => {});
    }
    if (_state.scenario && window.go?.main?.App?.AddScenarioToChat) {
        window.go.main.App.AddScenarioToChat(_state.scenario).catch(() => {});
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

    // Auto-populate characters and lorebooks from scenario
    if (scenario) {
        const rawChars      = scenario.characters || scenario.Characters || [];
        let   scenarioChars = rawChars.map(c => c.id ? c : Object.assign({}, c, { id: `sc-${(c.name || '').toLowerCase().replace(/\s+/g, '-')}-${scenario.id || ''}` }));
        const scenarioBooks = scenario.lorebooks  || scenario.Lorebooks  || [];

        console.log('[applyScenarioChange] scenario:', scenario.id, scenario.title || scenario.name);
        console.log('[applyScenarioChange] rawChars count:', rawChars.length, 'scenarioBooks count:', scenarioBooks.length);

        // Look up full character records from localStorage to get their individual lorebooks
        let charLorebooks = [];
        try {
            const allChars = JSON.parse(localStorage.getItem('howl-chat-characters') || '[]');
            console.log('[applyScenarioChange] allChars count:', allChars.length);
            charLorebooks = scenarioChars.flatMap(sc => {
                const fullChar = allChars.find(c => c.id === sc.id || c.name === sc.name);
                console.log('[applyScenarioChange] lookup char:', sc.id, sc.name, 'found:', !!fullChar);
                if (fullChar) {
                    const books = fullChar?.howl?.lorebooks || fullChar?.lorebooks || [];
                    console.log('[applyScenarioChange] char lorebooks:', books.length, books.map(b => b.id || b.name));
                    return books;
                }
                return [];
            }).filter(Boolean);
        } catch (e) { console.error('[applyScenarioChange] char lookup error:', e); }

        console.log('[applyScenarioChange] charLorebooks total:', charLorebooks.length);

        // Enrich lorebooks from localStorage so descriptions survive the backend round-trip
        const enrichedScenarioBooks = _enrichLorebooksFromLocalStorage(scenarioBooks);
        const enrichedCharBooks     = _enrichLorebooksFromLocalStorage(charLorebooks);

        // Combine scenario-level lorebooks with character-level lorebooks
        const allScenarioBooks = _dedupeById([...enrichedScenarioBooks, ...enrichedCharBooks]);
        console.log('[applyScenarioChange] allScenarioBooks:', allScenarioBooks.length, allScenarioBooks.map(b => b.id || b.name));

        // Never re-inject characters the user explicitly removed from this chat
        const removedSet = new Set(_state.removedCharacterIds || []);
        scenarioChars    = scenarioChars.filter(c => !removedSet.has(c.id));

        // Merge with existing, deduplicate by id
        const mergedChars = _dedupeById([..._state.characters, ...scenarioChars]);
        const mergedBooks = _dedupeById([..._state.lorebooks,  ...allScenarioBooks]);

        // If hot-swapping, remove chars that were ONLY from the previous scenario
        if (prev) {
            const prevCharIds = new Set((prev.characters || []).map(c => c.id));
            const newCharIds  = new Set(scenarioChars.map(c => c.id));
            _state.characters = mergedChars.filter(c => !prevCharIds.has(c.id) || newCharIds.has(c.id));

            const prevBookIds = new Set((prev.lorebooks || []).map(l => l.id));
            const newBookIds  = new Set(allScenarioBooks.map(l => l.id));
            _state.lorebooks  = mergedBooks.filter(l => !prevBookIds.has(l.id) || newBookIds.has(l.id));
        } else {
            _state.characters = mergedChars;
            _state.lorebooks  = mergedBooks;
        }

        console.log('[applyScenarioChange] final _state.lorebooks:', _state.lorebooks.length, _state.lorebooks.map(b => b.id || b.name));

        // Fire lorebook system event if characters changed
        if (prev && prev.id !== scenario.id) {
            _emitScenarioSwapEvent(prev, scenario);
        }
    }

    _push(); await refreshAllPanels();
    _syncScenarioAndCharactersToBackend();
}

async function applyCharacterAdd(character) {
    if (!character || _state.characters.some(c => c.id === character.id)) return;
    // Do not add back a character the user explicitly removed
    if ((_state.removedCharacterIds || []).includes(character.id)) return;
    _state.characters.push(character);
    _push(); await refreshAllPanels();
    if (window.go?.main?.App?.AddCharacterToChat) {
        window.go.main.App.AddCharacterToChat(character).catch(() => {});
    }
}

async function applyCharacterRemove(characterId) {
    console.log('[applyCharacterRemove] called for:', characterId);
    console.log('[applyCharacterRemove] _state before:', { world: !!_state.world, scenario: !!_state.scenario, characters: (_state.characters || []).length, lorebooks: (_state.lorebooks || []).length });
    _state.characters = _state.characters.filter(c => c.id !== characterId);
    if (!_state.removedCharacterIds.includes(characterId)) {
        _state.removedCharacterIds.push(characterId);
    }
    console.log('[applyCharacterRemove] _state after filter:', { world: !!_state.world, scenario: !!_state.scenario, characters: (_state.characters || []).length, lorebooks: (_state.lorebooks || []).length });
    _push(); await refreshAllPanels();
    // Persist removal list so ghosts stay dead after reload
    if (typeof saveChatData === 'function') saveChatData();
    console.log('[applyCharacterRemove] completed');
}

async function applyLorebookToggle(lorebook) {
    if (!lorebook) return;
    const idx = _state.lorebooks.findIndex(l => l.id === lorebook.id);
    if (idx === -1) { _state.lorebooks.push(lorebook); }
    else            { _state.lorebooks.splice(idx, 1); }
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
        const role   = sender.toLowerCase().includes('user') ? 'user' : 'assistant';
        return { id: i.toString(), sender, text, role, timestamp: new Date().toISOString() };
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
            messageCount:   history.length
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
        removedCharacterIds: _state.removedCharacterIds
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
        removedCharacterIds: _state.removedCharacterIds
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
        const stored = JSON.parse(localStorage.getItem('howl-worlds') || '[]');
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
        const stored = JSON.parse(localStorage.getItem('howl-chat-scenarios') || '[]');
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

function _enrichCharacterFromLocalStorage(character) {
    if (!character) return character;
    try {
        const stored = JSON.parse(localStorage.getItem('howl-chat-characters') || '[]');
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
        const stored = JSON.parse(localStorage.getItem('howl-chat-lorebooks') || '[]');
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
// beforeunload — auto-save if there's an active session
// ---------------------------------------------------------------------------

window.addEventListener('beforeunload', function () {
    if (!_state.sessionId) return;
    if (window.go?.main?.App?.SaveChat) {
        try { window.go.main.App.SaveChat(); } catch (_) {}
    }
});

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
    applyLorebookToggle,
    collectForSave,
    refreshAllPanels,
    getState
};
