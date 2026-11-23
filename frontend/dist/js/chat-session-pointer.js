/**
 * Chat Session Pointer — the only renderer owner of `howl-chat-current-session`
 *
 * The backend conversation is authoritative for the active chat id. This module
 * stores only that id and the renderer's context selection. Full catalog records
 * are hydrated from the character, world, scenario and lorebook repositories.
 *
 * @module chat-session-pointer
 * @version 1.0.0
 */

'use strict';

const ChatSessionPointer = (() => {
    const KEY = 'howl-chat-current-session';
    const VERSION = 2;
    const TRACKING_FIELDS = [
        'activeCharacterIds',
        'removedCharacterIds',
        'addedCharacterIds',
        'activeLorebookIds',
        'removedLorebookIds',
        'addedLorebookIds',
    ];
    const CHARACTER_SELECTION_FIELDS = ['status', 'active', 'stepOutReason', 'stepOutMode'];

    function _currentChatId() {
        const bridge = window.chatBridge?.getState?.();
        return bridge?.sessionId
            || window.dataSync?.currentChatId
            || window.sidebarState?.currentChatData?.sessionId
            || null;
    }

    function _idOf(value) {
        if (!value) return null;
        if (typeof value === 'string') return value;
        return value.id || value.ID || null;
    }

    function _characterSelection(character) {
        const selected = { id: _idOf(character) };
        if (!selected.id) return null;
        for (const field of CHARACTER_SELECTION_FIELDS) {
            if (character[field] !== undefined) selected[field] = character[field];
        }
        return selected;
    }

    function _array(value) {
        return Array.isArray(value) ? value : [];
    }

    function _tracking(chatData) {
        const tracking = {};
        for (const field of TRACKING_FIELDS) {
            tracking[field] = _array(chatData[field]).map(String);
        }
        return tracking;
    }

    function _selection(chatData) {
        const data = chatData || {};
        return {
            worldId: _idOf(data.world),
            scenarioId: _idOf(data.scenario),
            characters: _array(data.characters).map(_characterSelection).filter(Boolean),
            lorebookIds: _array(data.lorebooks).map(_idOf).filter(Boolean),
            ..._tracking(data),
        };
    }

    function build(chatData, chatId) {
        return {
            version: VERSION,
            chatId: chatId || _currentChatId(),
            selection: _selection(chatData),
        };
    }

    function _legacySelection(raw) {
        return _selection({
            world: raw?.world,
            scenario: raw?.scenario,
            characters: raw?.characters,
            lorebooks: raw?.lorebooks,
            activeCharacterIds: raw?.activeCharacterIds,
            removedCharacterIds: raw?.removedCharacterIds,
            addedCharacterIds: raw?.addedCharacterIds,
            activeLorebookIds: raw?.activeLorebookIds,
            removedLorebookIds: raw?.removedLorebookIds,
            addedLorebookIds: raw?.addedLorebookIds,
        });
    }

    /**
     * Parse the pointer. A pre-pointer blob is converted to the id-only shape
     * once and returned with an in-memory `legacy` copy so the first hydration
     * can still fall back to records that the repository cache may not have.
     */
    function read() {
        let raw;
        try {
            raw = localStorage.getItem(KEY);
        } catch (error) {
            console.warn('[ChatSessionPointer] read failed:', error);
            return null;
        }
        if (!raw) return null;

        let parsed;
        try {
            parsed = JSON.parse(raw);
        } catch (error) {
            console.warn('[ChatSessionPointer] malformed pointer ignored:', error);
            return null;
        }
        if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) return null;

        if (parsed.version === VERSION && parsed.selection && typeof parsed.selection === 'object') {
            return {
                version: VERSION,
                chatId: parsed.chatId || _currentChatId(),
                selection: parsed.selection,
            };
        }

        return {
            version: VERSION,
            chatId: _currentChatId(),
            selection: _legacySelection(parsed),
            legacy: {
                world: parsed.world || null,
                scenario: parsed.scenario || null,
                characters: _array(parsed.characters),
                lorebooks: _array(parsed.lorebooks),
            },
        };
    }

    function write(chatData, chatId) {
        const pointer = build(chatData, chatId);
        try {
            localStorage.setItem(KEY, JSON.stringify(pointer));
        } catch (error) {
            console.warn('[ChatSessionPointer] write failed:', error);
        }
        return pointer;
    }

    function clear() {
        try {
            localStorage.removeItem(KEY);
        } catch (error) {
            console.warn('[ChatSessionPointer] clear failed:', error);
        }
    }

    function _find(list, id) {
        if (!id) return null;
        return _array(list).find(item => String(_idOf(item)) === String(id)) || null;
    }

    async function _worldCatalog() {
        const repo = window.WorldRepository;
        if (repo?.ready) await repo.ready;
        return repo?.all?.() || [];
    }

    async function _scenarioCatalog() {
        const repo = window.ScenarioRepository;
        if (repo?.ready) await repo.ready;
        return repo?.all?.() || [];
    }

    async function _characterCatalog() {
        const repo = window.CharacterStorage || window.CharacterRepository;
        if (repo?.ready) await repo.ready;
        return repo?.all?.() || [];
    }

    async function _lorebookCatalog() {
        const repo = window.LorebookRepo;
        if (!repo) return [];
        if (!repo.state?.ready && typeof repo.hydrate === 'function') {
            await repo.hydrate();
        }
        return repo.all?.() || repo.state?.entries || (typeof repo.readCache === 'function' ? repo.readCache() : []);
    }

    function _characterWithSelection(character, selected) {
        const merged = { ...character };
        for (const field of CHARACTER_SELECTION_FIELDS) {
            if (selected[field] !== undefined) merged[field] = selected[field];
        }
        return merged;
    }

    async function hydrate(pointer) {
        const source = pointer || read();
        if (!source) return null;

        const selection = source.selection || {};
        const legacy = source.legacy || {};
        const [worlds, scenarios, characters, lorebooks] = await Promise.all([
            _worldCatalog(),
            _scenarioCatalog(),
            _characterCatalog(),
            _lorebookCatalog(),
        ]);

        const world = _find(worlds, selection.worldId) || _find([legacy.world], selection.worldId);
        const scenario = _find(scenarios, selection.scenarioId) || _find([legacy.scenario], selection.scenarioId);

        const selectedCharacters = _array(selection.characters)
            .map(selected => {
                const full = _find(characters, selected.id) || _find(legacy.characters, selected.id);
                return full ? _characterWithSelection(full, selected) : null;
            })
            .filter(Boolean);

        const selectedLorebooks = _array(selection.lorebookIds)
            .map(id => _find(lorebooks, id) || _find(legacy.lorebooks, id))
            .filter(Boolean);

        return {
            characters: selectedCharacters,
            world,
            scenario,
            lorebooks: selectedLorebooks,
            chatName: window.chatBridge?.getState?.()?.chatName
                || window.sidebarState?.currentChatData?.chatName
                || 'New Chat',
            ..._tracking(selection),
        };
    }

    return {
        KEY,
        VERSION,
        build,
        read,
        write,
        clear,
        hydrate,
        currentChatId: _currentChatId,
    };
})();

window.ChatSessionPointer = ChatSessionPointer;
