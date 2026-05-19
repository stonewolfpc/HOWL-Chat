/**
 * Sidebar Panels â€” Orchestrator
 *
 * Entry point that wires together the modular sidebar components:
 *   sidebar-state.js       â€” shared state, constants, utilities
 *   sidebar-session.js     â€” session load/save/sync
 *   sidebar-characters.js  â€” character panel rendering & interactions
 *   sidebar-context.js     â€” world/scenario/lorebook panel
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

function initializeSidebarPanels() {
    if (window._chatDataInitialized) return;
    window._chatDataInitialized = true;
    console.log('Initializing Sidebar Panels (v2.1)...');
    initializeCharacterPanel();
    window.sidebarContext.initializeContextPanel();
    loadChatData();
}

function initializeCharacterPanel() {
    console.log('Initializing Character Panel...');
}

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initializeSidebarPanels);
} else {
    setTimeout(initializeSidebarPanels, 0);
}

// ---------------------------------------------------------------------------
// Delegation shims â€” route bare function calls to the correct module
// ---------------------------------------------------------------------------

function loadChatData() { return window.sidebarSession.loadChatData(); }
function saveChatData() { return window.sidebarSession.saveChatData(); }
// Re-establish canonical panel shims after sidebar modules load (they overwrite context-display.js shims)
function updateCharacterPanel()  { return contextDisplay.updateCharacterPanel(); }
function updateContextDisplay()  { return contextDisplay.updateContextDisplay(); }
function updateWorldDisplay()    { return contextDisplay.updateWorldDisplay(); }
function updateScenarioDisplay() { return contextDisplay.updateScenarioDisplay(); }
function updateLorebookDisplay() { return contextDisplay.updateLorebookDisplay(); }

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
                updateCharacterPanel();
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
    updateCharacterPanel,
    updateContextDisplay,
    selectWorld:          (id)    => window.sidebarContext.selectWorld(id),
    selectScenario:       (id)    => window.sidebarContext.selectScenario(id),
    openCharacterSheet:   (id)    => window.sidebarCharacters.openCharacterSheet(id),
    showCharacterMenu:    (e, id) => window.sidebarCharacters.showCharacterMenu(e, id),
    openWorldSelector:    ()      => window.sidebarContext.openWorldSelector(),
    openScenarioSelector: ()      => window.sidebarContext.openScenarioSelector(),
    openLorebookManager:  ()      => window.sidebarContext.openLorebookManager(),
    handleCharacterCreation
};

// Expose legacy global that some pages bind to directly
window.logChatContract = (...args) => window.sidebarState.logChatContract(...args);
