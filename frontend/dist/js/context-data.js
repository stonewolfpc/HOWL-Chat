/**
 * Context Data — ContextDisplay Class (Data Layer)
 *
 * Defines the ContextDisplay class with all data loading,
 * persistence, backend sync, and resolution utility methods.
 * Display and selector logic live in context-panels.js and
 * context-selectors.js respectively.
 *
 * @module context-data
 * @version 1.0.0
 */

'use strict';

class ContextDisplay {
    constructor() {
        this.currentChatData = {
            characters: [],
            world:      null,
            scenario:   null,
            lorebooks:  [],
            chatName:   'New Chat'
        };
    }

    // -----------------------------------------------------------------------
    // Resolution Utilities
    // -----------------------------------------------------------------------

    resolveImageSrc(imageValue) {
        if (!imageValue) return '';
        if (typeof imageValue === 'string') return imageValue;
        if (typeof imageValue === 'object') {
            if (typeof imageValue.path      === 'string') return imageValue.path;
            if (typeof imageValue.url       === 'string') return imageValue.url;
            if (typeof imageValue.data      === 'string') return imageValue.data;
            if (typeof imageValue.thumbnail === 'string') return imageValue.thumbnail;
            // {id, stored:true} references need async HOWLStorage lookup
            if (imageValue.id && imageValue.stored) return '';
        }
        return '';
    }

    refMatches(item, ref) {
        if (!item || !ref) return false;
        const itemId   = (item.id   || '').toString().trim().toLowerCase();
        const refId    = (ref.id    || '').toString().trim().toLowerCase();
        if (itemId && refId && itemId === refId) return true;
        const itemName = (item.name || item.title || '').toString().trim().toLowerCase();
        const refName  = (ref.name  || ref.title  || '').toString().trim().toLowerCase();
        return !!(itemName && refName && itemName === refName);
    }

    resolveLinkedItems(allItems, refs) {
        if (!Array.isArray(allItems) || !Array.isArray(refs) || refs.length === 0) return [];
        return allItems.filter(item => refs.some(ref => this.refMatches(item, ref)));
    }

    getScenarioLorebookRefs(scenario) {
        const howlRefs = scenario?.howl?.lorebooks;
        const topRefs  = scenario?.lorebooks;
        if (Array.isArray(howlRefs) && howlRefs.length > 0) return howlRefs;
        if (Array.isArray(topRefs)  && topRefs.length  > 0) return topRefs;
        return [];
    }

    mergeLorebookRefs(catalog, refsList) {
        const merged = [];
        const seen   = new Set();
        refsList.forEach(refs => {
            this.resolveLinkedItems(catalog, refs).forEach(lb => {
                const key = (lb?.id || lb?.title || lb?.name || '').toString().trim().toLowerCase();
                if (!key || seen.has(key)) return;
                seen.add(key);
                merged.push(lb);
            });
        });
        return merged;
    }

    // -----------------------------------------------------------------------
    // Initialization
    // -----------------------------------------------------------------------

    initialize() {
        // Sidebar-panels.js is the canonical startup orchestrator.
        // Only set up event listeners here; do not independently load chats.
        this.setupEventListeners();
    }

    setupEventListeners() {
        const worldSelector = document.getElementById('world-selector');
        if (worldSelector) worldSelector.addEventListener('click', () => this.openWorldSelector());

        const scenarioSelector = document.getElementById('scenario-selector');
        if (scenarioSelector) scenarioSelector.addEventListener('click', () => this.openScenarioSelector());

        const characterSelector = document.getElementById('character-selector');
        if (characterSelector) characterSelector.addEventListener('click', () => this.openCharacterSelector());
    }

    // -----------------------------------------------------------------------
    // Data Loading
    // -----------------------------------------------------------------------

    async loadChatData() {
        // Delegate to the canonical sidebar loader if available
        if (window.sidebarSession && typeof window.sidebarSession.loadChatData === 'function') {
            return window.sidebarSession.loadChatData();
        }

        const recentSession = await this.discoverMostRecentChatSession();
        if (recentSession) {
            console.log('[context-data] auto-loading most recent session:', recentSession.name || recentSession.Name);
            await this.loadChatSessionById(recentSession.id || recentSession.ID);
            return;
        }

        if (window.go && window.go.main && window.go.main.App) {
            try {
                const data = await window.go.main.App.LoadChatData();
                if (data) {
                    this.currentChatData = data;
                    await this.updateContextDisplay();
                } else {
                    this.resetToDefaultState();
                    await this.updateContextDisplay();
                }
            } catch (err) {
                console.error('Failed to load chat data:', err);
                await this.loadChatDataFromStorage();
            }
        } else {
            await this.loadChatDataFromStorage();
        }
    }

    async discoverMostRecentChatSession() {
        if (!window.go || !window.go.main || !window.go.main.App) return null;
        try {
            const chatSummaries = await window.go.main.App.GetChatList();
            if (chatSummaries && chatSummaries.length > 0) {
                const safeDate = (v) => { const d = new Date(v || 0); return isNaN(d.getTime()) ? new Date(0) : d; };
                return chatSummaries.sort((a, b) => safeDate(b.updatedAt || b.UpdatedAt) - safeDate(a.updatedAt || a.UpdatedAt))[0];
            }
        } catch (error) {
            console.error('Failed to discover recent chat session:', error);
        }
        return null;
    }

    async loadChatSessionById(chatId) {
        // Delegate to the canonical sidebar loader if available
        if (window.sidebarSession && typeof window.sidebarSession.loadChatSessionById === 'function') {
            return window.sidebarSession.loadChatSessionById(chatId);
        }

        if (!window.go || !window.go.main || !window.go.main.App) return;
        try {
            const chatData = await window.go.main.App.LoadChat(chatId);
            if (chatData) {
                const characters = chatData.characters || chatData.Characters || [];
                const world      = chatData.world      || chatData.World      || null;
                const scenario   = chatData.scenario   || chatData.Scenario   || null;
                const lorebooks  = chatData.lorebooks  || chatData.Lorebooks  || [];
                const history    = chatData.history    || chatData.History    || [];
                const meta       = chatData.meta       || chatData.Meta       || {};

                this.currentChatData = { characters, world, scenario, lorebooks, chatName: meta.name || meta.Name || 'New Chat' };
                await this.updateContextDisplay();
                if (Array.isArray(history) && typeof window.rebuildChatFromHistory === 'function') {
                    window.rebuildChatFromHistory(history);
                }
            }
        } catch (error) {
            console.error('Failed to load chat session:', error);
        }
    }

    async loadChatDataFromStorage() {
        try {
            const stored = localStorage.getItem('howl-chat-current-session');
            if (stored) {
                const light          = JSON.parse(stored);
                const allWorlds      = JSON.parse(localStorage.getItem('howl-worlds')            || '[]');
                const allScenarios   = JSON.parse(localStorage.getItem('howl-chat-scenarios')    || '[]');
                const allCharacters  = JSON.parse(localStorage.getItem('howl-chat-characters')   || '[]');
                const allLorebooks   = JSON.parse(localStorage.getItem('howl-chat-lorebooks')    || '[]');
                this.currentChatData = {
                    characters: light.characters ? light.characters.map(lc => {
                        const full = allCharacters.find(ac => ac.id === lc.id);
                        if (!full) return null;
                        return {
                            ...full,
                            status:        lc.status        || full.status,
                            active:        typeof lc.active === 'boolean' ? lc.active : full.active,
                            stepOutReason: lc.stepOutReason || full.stepOutReason,
                            stepOutMode:   lc.stepOutMode   || full.stepOutMode
                        };
                    }).filter(Boolean) : [],
                    world:      light.world      ? allWorlds.find(aw => aw.id === light.world.id)           : null,
                    scenario:   light.scenario   ? allScenarios.find(as => as.id === light.scenario.id)     : null,
                    lorebooks:  light.lorebooks  ? light.lorebooks.map(ll => allLorebooks.find(al => al.id === ll.id)).filter(Boolean) : [],
                    chatName:   light.chatName || 'New Chat'
                };
            } else {
                this.resetToDefaultState();
            }
            await this.updateContextDisplay();
        } catch (error) {
            console.error('Failed to load from localStorage:', error);
            this.resetToDefaultState();
            await this.updateContextDisplay();
        }
    }

    // -----------------------------------------------------------------------
    // State Reset & Persistence
    // -----------------------------------------------------------------------

    resetToDefaultState() {
        this.currentChatData = {
            characters: [], world: null, scenario: null, lorebooks: [], chatName: 'New Chat', removedCharacterIds: []
        };
    }

    saveChatData() {
        try {
            const light = {
                characters: this.currentChatData.characters
                    ? this.currentChatData.characters.map(c => ({ id: c.id, name: c.name })) : [],
                world: this.currentChatData.world
                    ? { id: this.currentChatData.world.id, name: this.currentChatData.world.name, title: this.currentChatData.world.title } : null,
                scenario: this.currentChatData.scenario
                    ? { id: this.currentChatData.scenario.id, name: this.currentChatData.scenario.name, title: this.currentChatData.scenario.title } : null,
                lorebooks: this.currentChatData.lorebooks
                    ? this.currentChatData.lorebooks.map(lb => ({ id: lb.id, name: lb.name, title: lb.title })) : [],
                removedCharacterIds: this.currentChatData.removedCharacterIds || [],
                chatName: this.currentChatData.chatName
            };
            localStorage.setItem('howl-chat-current-session', JSON.stringify(light));
        } catch (error) {
            console.error('Failed to save chat data:', error);
        }
    }

    // -----------------------------------------------------------------------
    // Backend Sync
    // -----------------------------------------------------------------------

    async syncClearedContextToBackend() {
        await this.syncLorebooksDataToBackend([]);
        await this.syncCharactersDataToBackend([]);
    }

    async syncWorldDataToBackend(world) {
        if (window.go?.main?.App?.AddWorldToChat) {
            try   { await window.go.main.App.AddWorldToChat(world); }
            catch (e) { console.warn('Backend world sync failed:', e); }
        }
    }

    async syncScenarioDataToBackend(scenario) {
        if (window.go?.main?.App?.AddScenarioToChat) {
            try   { await window.go.main.App.AddScenarioToChat(scenario); }
            catch (e) { console.warn('Backend scenario sync failed:', e); }
        }
    }

    async syncCharactersDataToBackend(characters) {
        if (window.go?.main?.App?.AddCharacterToChat) {
            try {
                for (const c of characters) await window.go.main.App.AddCharacterToChat(c);
            } catch (e) { console.warn('Backend characters sync failed:', e); }
        }
    }

    async syncLorebooksDataToBackend(lorebooks) {
        if (window.go?.main?.App?.SetLorebooks) {
            try   { await window.go.main.App.SetLorebooks(lorebooks); }
            catch (e) { console.warn('Backend lorebooks sync failed:', e); }
        }
    }
}
