/**
 * Sidebar Panels — Orchestrator
 *
 * Entry point that wires together the modular sidebar components:
 *   sidebar-state.js       — shared state, constants, utilities
 *   sidebar-session.js     — session load/save/sync
 *   sidebar-characters.js  — character panel rendering & interactions
 *   sidebar-context.js     — world/scenario/lorebook panel
 *
 * This file only handles:
 *   - Initialization sequencing
 *   - Cross-page postMessage event handlers
 *   - Navigation and notification helpers
 *   - The public window.sidebarPanels API surface
 *
 * @module sidebar-panels
 * @version 2.1.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Initialization
// ---------------------------------------------------------------------------

async function initializeSidebarPanels() {
    if (window._chatDataInitialized) return;
    console.log('Initializing Sidebar Panels (v2.1)...');
    initializeCharacterPanel();
    window.sidebarContext.initializeContextPanel();
    await loadChatData();
    window._chatDataInitialized = true;
    await applyPendingChatWorld();
}

function initializeCharacterPanel() {
    console.log('Initializing Character Panel...');
}

/**
 * A world created through the chat's "Change World -> + New World" flow is
 * recorded in sessionStorage by world-form.js before navigating back here.
 * That navigation is a full page load, so there is no live postMessage channel
 * and the world would otherwise never reach the chat that requested it.
 */
async function applyPendingChatWorld() {
    if (!document.getElementById('world-name-display')) return;
    let worldId = null;
    try { worldId = sessionStorage.getItem('howl-pending-chat-world'); } catch (_) { return; }
    if (!worldId) return;
    try { sessionStorage.removeItem('howl-pending-chat-world'); } catch (_) {}
    if (typeof contextDisplay === 'undefined' || !contextDisplay.selectWorld) return;
    try {
        await contextDisplay.selectWorld(worldId);
        console.log('[applyPendingChatWorld] attached new world:', worldId);
    } catch (err) {
        console.error('[applyPendingChatWorld] failed to attach world:', err);
    }
}

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => initializeSidebarPanels());
} else {
    setTimeout(() => initializeSidebarPanels(), 0);
}

// ---------------------------------------------------------------------------
// Delegation shims — route bare function calls to the correct module
// ---------------------------------------------------------------------------

// loadChatData / saveChatData belong to sidebar-session.js, which owns the load
// sequence and the save contract. No delegating copies live here.

// context-display.js owns updateContextDisplay / updateLorebookDisplay and the
// three panel renderers: it defines them on window and instantiates ContextDisplay.
// story-tools.html loads these sidebar modules without that orchestrator and has
// no world, scenario or lorebook panel at all, so its shared session code used to
// call a function that was either undefined or threw on an undefined contextDisplay.
// Panels that do not exist on the page get an explicit no-op.
if (typeof window.contextDisplay === 'undefined') {
    window.updateContextDisplay  = function () {};
    window.updateLorebookDisplay = function () {};
    window.updateWorldDisplay    = function () {};
    window.updateScenarioDisplay = function () {};
}
// updateCharacterPanel is defined once, by context-display.js. The copy that used
// to live here re-asserted the same body after the sidebar modules loaded, which
// is how the shadowing went unnoticed for so long.
// ---------------------------------------------------------------------------
// Cross-Page postMessage Handlers
// ---------------------------------------------------------------------------

function handleCharacterCreation() {
    if (!window.location.pathname.includes('character.html')) return;
    window.addEventListener('message', (event) => {
        if (event.data && event.data.type === 'character-created') {
            const character = event.data.character;
            if (character && window.go && window.go.main && window.go.main.App && window.go.main.App.AddCharacterToChat) {
                window.go.main.App.AddCharacterToChat(character)
                    .then(() => { console.log('Character added to chat session'); loadChatData(); })
                    .catch(err => console.error('Failed to add character to chat:', err));
            }
        } else if (event.data && event.data.type === 'character-updated') {
            const { character, index } = event.data;
            const cd = window.sidebarState.currentChatData;
            if (character && index !== undefined && cd.characters && cd.characters[index]) {
                cd.characters[index] = character;
                saveChatData();
                if (window.AudioService?.invalidateVoiceCache) {
                    window.AudioService.invalidateVoiceCache(character.id);
                }
                if (window.updateCharacterPanel) window.updateCharacterPanel();
            }
        }
    });
}

function handleWorldCreation() {
    if (!window.location.pathname.includes('world.html')) return;
    window.addEventListener('message', (event) => {
        if (event.data && event.data.type === 'world-created') {
            const world = event.data.world;
            if (world && window.go && window.go.main && window.go.main.App && window.go.main.App.AddWorldToChat) {
                window.go.main.App.AddWorldToChat(world)
                    .then(() => { console.log('World added to chat session'); loadChatData(); })
                    .catch(err => console.error('Failed to add world to chat:', err));
            }
        }
    });
}

function handleScenarioCreation() {
    if (!window.location.pathname.includes('scenario.html')) return;
    window.addEventListener('message', (event) => {
        if (event.data && event.data.type === 'scenario-created') {
            const scenario = event.data.scenario;
            if (scenario && window.go && window.go.main && window.go.main.App && window.go.main.App.AddScenarioToChat) {
                window.go.main.App.AddScenarioToChat(scenario)
                    .then(() => { console.log('Scenario added to chat session'); loadChatData(); })
                    .catch(err => console.error('Failed to add scenario to chat:', err));
            }
        }
    });
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

window.sidebarPanels = {
    loadChatData,
    // Resolved at call time: context-display.js owns both, and this object is
    // built on pages that do not load it (story-tools.html).
    updateCharacterPanel: (...args) => window.updateCharacterPanel?.(...args),
    updateContextDisplay: (...args) => window.updateContextDisplay?.(...args),
    selectWorld:          (id)    => window.sidebarContext.selectWorld(id),
    selectScenario:       (id)    => window.sidebarContext.selectScenario(id),
    openCharacterSheet:   (id)    => window.sidebarCharacters.openCharacterSheet(id),
    showCharacterMenu:    (e, id) => window.sidebarCharacters.showCharacterMenu(e, id),
    openWorldSelector:    ()      => window.sidebarContext.openWorldSelector(),
    openScenarioSelector: ()      => window.sidebarContext.openScenarioSelector(),
    openLorebookManager:  ()      => window.sidebarContext.openLorebookManager(),
    handleCharacterCreation
};

// logChatContract is owned by sidebar-state.js, which already exposes it on window.
