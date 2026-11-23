/**
 * Delete Service
 *
 * Detach helpers for permanent deletion of worlds, scenarios, characters,
 * and lorebooks. Each helper removes cross-references from every related
 * localStorage store before the caller deletes the record itself.
 *
 * Public API (window.DeleteService):
 *   detachWorldFromAll(worldId)      — removes world from chat current-data
 *   detachScenarioFromAll(scenarioId)— removes scenario refs from worlds + chat
 *   detachCharacterFromAll(charId)   — removes char refs from scenarios + worlds + chat
 *   detachLorebookFromAll(lbId)      — removes lorebook refs from worlds, scenarios, characters
 *
 * @module delete-service
 * @version 1.0.0
 */

'use strict';

const DeleteService = (() => {

    // -----------------------------------------------------------------------
    // localStorage helpers
    // -----------------------------------------------------------------------

    function _readLS(key)       { try { return JSON.parse(localStorage.getItem(key) || '[]'); } catch { return []; } }
    function _writeLS(key, val) { try { localStorage.setItem(key, JSON.stringify(val)); } catch (e) { console.warn('DeleteService _writeLS:', e); } }

    // -----------------------------------------------------------------------
    // World detach
    // -----------------------------------------------------------------------

    /**
     * Remove world cross-references when a world is permanently deleted.
     * - Clears world from the in-memory chat bridge state
     * - Clears world from the HOWLStorage current-chat record
     * @param {string} worldId
     */
    function detachWorldFromAll(worldId) {
        // Clear from active chat bridge state
        if (window.chatBridge) {
            const st = window.chatBridge.getState();
            if (st.world && st.world.id === worldId) {
                if (window.go?.main?.App?.ClearWorldFromSession) {
                    window.go.main.App.ClearWorldFromSession().catch(() => {});
                }
            }
        }

        // Clear from HOWLStorage current-chat record
        _clearEntityFromCurrentChat('world', worldId);
    }

    // -----------------------------------------------------------------------
    // Scenario detach
    // -----------------------------------------------------------------------

    /**
     * Remove scenario cross-references when a scenario is permanently deleted.
     * - Removes scenario reference from any world's howl.scenarios list
     * - Clears scenario from active chat state
     * @param {string} scenarioId
     */
    function detachScenarioFromAll(scenarioId) {
        // Remove from worlds
        const worlds = _readLS('howl-worlds');
        let worldsChanged = false;
        worlds.forEach(w => {
            if (w.howl && Array.isArray(w.howl.scenarios)) {
                const before = w.howl.scenarios.length;
                w.howl.scenarios = w.howl.scenarios.filter(s => s.id !== scenarioId);
                if (w.howl.scenarios.length !== before) worldsChanged = true;
            }
        });
        if (worldsChanged) _writeLS('howl-worlds', worlds);

        // Clear from active chat
        if (window.chatBridge) {
            const st = window.chatBridge.getState();
            if (st.scenario && st.scenario.id === scenarioId) {
                if (window.go?.main?.App?.ClearScenarioFromSession) {
                    window.go.main.App.ClearScenarioFromSession().catch(() => {});
                }
            }
        }

        _clearEntityFromCurrentChat('scenario', scenarioId);
    }

    // -----------------------------------------------------------------------
    // Character detach
    // -----------------------------------------------------------------------

    /**
     * Remove character cross-references when a character is permanently deleted.
     * - Removes character from all scenario character lists
     * - Removes character from all world character lists (howl.characters)
     * - Clears character from active chat state + calls backend
     * @param {string} charId
     */
    async function detachCharacterFromAll(charId) {
        // Remove from scenarios
        const scenarios = _readLS('howl-chat-scenarios');
        let scenariosChanged = false;
        scenarios.forEach(s => {
            if (Array.isArray(s.characters)) {
                const before = s.characters.length;
                s.characters = s.characters.filter(c => c.id !== charId);
                if (s.characters.length !== before) scenariosChanged = true;
            }
        });
        if (scenariosChanged) _writeLS('howl-chat-scenarios', scenarios);

        // Remove from worlds
        const worlds = _readLS('howl-worlds');
        let worldsChanged = false;
        worlds.forEach(w => {
            if (w.howl && Array.isArray(w.howl.characters)) {
                const before = w.howl.characters.length;
                w.howl.characters = w.howl.characters.filter(c => c.id !== charId);
                if (w.howl.characters.length !== before) worldsChanged = true;
            }
        });
        if (worldsChanged) _writeLS('howl-worlds', worlds);

        // Remove from active chat bridge + backend
        if (window.chatBridge) {
            await window.chatBridge.applyCharacterRemove(charId);
        }
        if (window.go?.main?.App?.ResetSessionCharacters) {
            const st = window.chatBridge ? window.chatBridge.getState() : null;
            const remaining = st ? st.characters.filter(c => c.id !== charId) : [];
            window.go.main.App.ResetSessionCharacters().then(() => {
                remaining.forEach(c => window.go.main.App.AddCharacterToChat(c).catch(() => {}));
            }).catch(() => {});
        }

        _clearCharacterFromCurrentChat(charId);
    }

    // -----------------------------------------------------------------------
    // Lorebook detach
    // -----------------------------------------------------------------------

    /**
     * Remove lorebook cross-references when a lorebook is permanently deleted.
     * - Removes lorebook refs from all worlds (howl.lorebooks)
     * - Removes lorebook refs from all scenarios (lorebooks)
     * - Removes lorebook refs from all characters (lorebooks / howl.lorebooks)
     * - Removes lorebook from active chat state
     * @param {string} lbId
     */
    async function detachLorebookFromAll(lbId) {
        // Remove from worlds
        const worlds = _readLS('howl-worlds');
        let worldsChanged = false;
        worlds.forEach(w => {
            if (w.howl && Array.isArray(w.howl.lorebooks)) {
                const before = w.howl.lorebooks.length;
                w.howl.lorebooks = w.howl.lorebooks.filter(l => l.id !== lbId);
                if (w.howl.lorebooks.length !== before) worldsChanged = true;
            }
        });
        if (worldsChanged) _writeLS('howl-worlds', worlds);

        // Remove from scenarios
        const scenarios = _readLS('howl-chat-scenarios');
        let scenariosChanged = false;
        scenarios.forEach(s => {
            if (Array.isArray(s.lorebooks)) {
                const before = s.lorebooks.length;
                s.lorebooks = s.lorebooks.filter(l => l.id !== lbId);
                if (s.lorebooks.length !== before) scenariosChanged = true;
            }
        });
        if (scenariosChanged) _writeLS('howl-chat-scenarios', scenarios);

        // Remove from characters through the canonical store. Work on a copy so
        // a failed save cannot leave the repository projection mutated.
        const charSource = typeof window.CharacterStorage?.all === 'function'
            ? JSON.parse(JSON.stringify((window.CharacterStorage?.all?.() || [])))
            : [];
        let charsChanged = false;
        charSource.forEach(c => {
            if (Array.isArray(c.lorebooks)) {
                const before = c.lorebooks.length;
                c.lorebooks = c.lorebooks.filter(l => l.id !== lbId);
                if (c.lorebooks.length !== before) charsChanged = true;
            }
            if (c.howl && Array.isArray(c.howl.lorebooks)) {
                const before = c.howl.lorebooks.length;
                c.howl.lorebooks = c.howl.lorebooks.filter(l => l.id !== lbId);
                if (c.howl.lorebooks.length !== before) charsChanged = true;
            }
        });
        if (charsChanged && typeof window.CharacterStorage?.save === 'function') {
            for (const character of charSource) {
                try {
                    await window.CharacterStorage.save(character);
                } catch (e) {
                    console.warn('[DeleteService] character lorebook detach save failed:', e);
                }
            }
        }

        // Remove from active chat via bridge
        if (window.chatBridge) {
            await window.chatBridge.applyLorebookToggle({ id: lbId }); // toggle removes if present
        }

        _clearLorebookFromCurrentChat(lbId);
    }

    // -----------------------------------------------------------------------
    // Current-chat HOWLStorage helpers
    // -----------------------------------------------------------------------

    function _clearEntityFromCurrentChat(field, id) {
        if (!window.HOWLStorage) return;
        window.HOWLStorage.getChatData('howl-chat-current-chat').then(data => {
            if (!data) return;
            if (data[field] && data[field].id === id) {
                data[field] = null;
                window.HOWLStorage.setChatData('howl-chat-current-chat', data).catch(() => {});
            }
        }).catch(() => {});
    }

    function _clearCharacterFromCurrentChat(charId) {
        if (!window.HOWLStorage) return;
        window.HOWLStorage.getChatData('howl-chat-current-chat').then(data => {
            if (!data || !Array.isArray(data.characters)) return;
            const before = data.characters.length;
            data.characters = data.characters.filter(c => c.id !== charId);
            if (data.characters.length !== before) {
                window.HOWLStorage.setChatData('howl-chat-current-chat', data).catch(() => {});
            }
        }).catch(() => {});
    }

    function _clearLorebookFromCurrentChat(lbId) {
        if (!window.HOWLStorage) return;
        window.HOWLStorage.getChatData('howl-chat-current-chat').then(data => {
            if (!data || !Array.isArray(data.lorebooks)) return;
            const before = data.lorebooks.length;
            data.lorebooks = data.lorebooks.filter(l => l.id !== lbId);
            if (data.lorebooks.length !== before) {
                window.HOWLStorage.setChatData('howl-chat-current-chat', data).catch(() => {});
            }
        }).catch(() => {});
    }

    return {
        detachWorldFromAll,
        detachScenarioFromAll,
        detachCharacterFromAll,
        detachLorebookFromAll
    };

})();

window.DeleteService = DeleteService;
