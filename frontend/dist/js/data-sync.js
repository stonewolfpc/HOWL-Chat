/**
 * Data Synchronization Module
 *
 * Handles synchronization between frontend state and backend persistence.
 * Manages chat data saving, loading, and state consistency.
 *
 * @module data-sync
 * @version 1.0.0
 */

class DataSync {
    constructor() {
        this.isInitialized = false;
        this.currentChatId = null;
        this.autosaveInterval = null;
        this.autosaveDelay = 30000; // 30 seconds
    }

    /**
     * Initialize data synchronization
     */
    initialize() {
        if (this.isInitialized) return;

        this.setupAutosave();
        this.loadLastChatId();
        this.isInitialized = true;

        console.log('Data synchronization initialized');
    }

    /**
     * Setup autosave functionality
     */
    setupAutosave() {
        // Clear any existing interval
        if (this.autosaveInterval) {
            clearInterval(this.autosaveInterval);
        }

        // Set up autosave interval
        this.autosaveInterval = setInterval(() => {
            this.performAutosave();
        }, this.autosaveDelay);
    }

    /**
     * Load last chat ID from config
     */
    async loadLastChatId() {
        if (!window.go || !window.go.main || !window.go.main.App) {
            console.log('Backend not available for chat ID loading');
            return;
        }

        try {
            const activeChatMeta = await window.go.main.App.GetActiveChatMeta();
            if (activeChatMeta) {
                this.currentChatId = activeChatMeta.ID;
                console.log('Loaded active chat ID:', this.currentChatId);
            }
        } catch (error) {
            console.error('Failed to load active chat meta:', error);
        }
    }

    /**
     * Perform autosave
     */
    async performAutosave() {
        if (!this.currentChatId || !this.hasUnsavedChanges()) {
            return;
        }

        try {
            await this.saveCurrentChatData();
            console.log('Autosave completed');
        } catch (error) {
            console.error('Autosave failed:', error);
        }
    }

    /**
     * Check if there are unsaved changes
     */
    hasUnsavedChanges() {
        // Check if contextDisplay has current chat data
        if (typeof contextDisplay !== 'undefined' && contextDisplay.currentChatData) {
            // Compare with last saved state (simplified check)
            return true; // For now, always assume changes exist
        }
        return false;
    }

    /**
     * Save current chat data to backend
     */
    async saveCurrentChatData() {
        if (!window.go || !window.go.main || !window.go.main.App) {
            console.error('Backend not available for saving');
            return;
        }

        try {
            await window.go.main.App.SaveChat();
            console.log('Chat data saved successfully');
        } catch (error) {
            console.error('Failed to save chat data:', error);
            throw error;
        }
    }

    /**
     * Collect complete chat data for saving
     */
    collectCompleteChatData() {
        console.log('🔍 Collecting complete chat data for save...');

        let chatData = {
            meta: {
                id: this.currentChatId || Date.now().toString(),
                name: 'Chat Session',
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

        // Get data from contextDisplay if available
        if (typeof contextDisplay !== 'undefined' && contextDisplay.currentChatData) {
            const data = contextDisplay.currentChatData;

            chatData.meta.name = data.chatName || 'Chat Session';
            chatData.meta.worldName = data.world?.title || data.world?.name || '';
            chatData.meta.scenarioName = data.scenario?.title || data.scenario?.name || '';
            chatData.meta.characterNames = (data.characters || []).map(c => c.name || c.title || '').filter(Boolean);
            chatData.meta.lorebookNames = (data.lorebooks || []).map(l => l.name || l.title || '').filter(Boolean);

            chatData.world = data.world;
            chatData.scenario = data.scenario;
            chatData.characters = data.characters || [];
            chatData.lorebooks = data.lorebooks || [];
        }

        // Get chat history from DOM
        const chatMessages = document.querySelector('.chat-messages');
        if (chatMessages) {
            const messageElements = chatMessages.querySelectorAll('.chat-message');
            chatData.history = Array.from(messageElements).map((msgEl, index) => {
                const senderEl = msgEl.querySelector('.chat-sender');
                const textEl = msgEl.querySelector('.chat-text');

                if (senderEl && textEl) {
                    const sender = senderEl.textContent.trim();
                    const text = textEl.textContent.trim();

                    // Determine role from sender
                    let role = 'system';
                    if (sender.toLowerCase().includes('user')) {
                        role = 'user';
                    } else if (sender.toLowerCase().includes('assistant') || sender.toLowerCase().includes('ai')) {
                        role = 'assistant';
                    }

                    return {
                        id: index.toString(),
                        sender: sender,
                        text: text,
                        role: role,
                        timestamp: new Date().toISOString(),
                        hasImage: false
                    };
                }
                return null;
            }).filter(Boolean);

            chatData.meta.messageCount = chatData.history.length;
        }

        console.log('📋 Collected chat data for save:', {
            hasWorld: !!chatData.world,
            worldName: chatData.meta.worldName,
            hasScenario: !!chatData.scenario,
            scenarioName: chatData.meta.scenarioName,
            characterCount: chatData.characters.length,
            lorebookCount: chatData.lorebooks.length,
            historyCount: chatData.history.length
        });

        return chatData;
    }

    /**
     * Load chat by ID
     */
    async loadChat(chatId) {
        if (!window.go || !window.go.main || !window.go.main.App) {
            console.error('Backend not available for chat loading');
            return null;
        }

        try {
            // Save current chat first
            if (this.currentChatId && this.currentChatId !== chatId) {
                await this.saveCurrentChatData();
            }

            const chatData = await window.go.main.App.LoadChat(chatId);
            if (chatData) {
                this.currentChatId = chatId;

                // Bridge hydrates both state systems and all panels
                if (window.chatBridge) {
                    await window.chatBridge.hydrateFromBackend(chatData);
                } else if (typeof contextDisplay !== 'undefined') {
                    contextDisplay.currentChatData = {
                        characters: chatData.Characters || [],
                        world: chatData.World,
                        scenario: chatData.Scenario,
                        lorebooks: chatData.Lorebooks || [],
                        chatName: chatData.Meta ? chatData.Meta.Name : 'New Chat'
                    };
                    await contextDisplay.updateContextDisplay();
                }

                console.log('Chat loaded successfully:', chatId);
                return chatData;
            }
        } catch (error) {
            console.error('Failed to load chat:', error);
        }
        return null;
    }

    /**
     * Create new chat
     */
    async createNewChat(name) {
        if (!window.go || !window.go.main || !window.go.main.App) {
            console.error('Backend not available for chat creation');
            return null;
        }

        try {
            // Save current chat first
            if (this.currentChatId) {
                await this.saveCurrentChatData();
            }

            const chatMeta = await window.go.main.App.NewChat(name);
            if (chatMeta) {
                this.currentChatId = chatMeta.ID || chatMeta.id;

                // Bridge clears both state systems and all panels
                if (window.chatBridge) {
                    await window.chatBridge.clearForNewChat(chatMeta.Name || chatMeta.name || name, this.currentChatId);
                } else if (typeof contextDisplay !== 'undefined') {
                    contextDisplay.clearChatMessagesDom();
                    contextDisplay.resetToDefaultState();
                    contextDisplay.currentChatData.chatName = chatMeta.Name || chatMeta.name || name;
                    await contextDisplay.updateContextDisplay();
                }

                if (typeof contextDisplay !== 'undefined') {
                    contextDisplay.saveChatData();
                    await contextDisplay.syncClearedContextToBackend();
                }

                await this.saveCurrentChatData();

                console.log('New chat created:', this.currentChatId);
                return chatMeta;
            }
        } catch (error) {
            console.error('Failed to create new chat:', error);
        }
        return null;
    }

    /**
     * Save chat as new chat
     */
    async saveChatAs(name) {
        if (!window.go || !window.go.main || !window.go.main.App) {
            console.error('Backend not available for save as');
            return null;
        }

        try {
            const chatMeta = await window.go.main.App.SaveChatAs(name);
            if (chatMeta) {
                this.currentChatId = chatMeta.ID;
                console.log('Chat saved as:', chatMeta.ID);
                return chatMeta;
            }
        } catch (error) {
            console.error('Failed to save chat as:', error);
        }
        return null;
    }

    /**
     * Delete chat
     */
    async deleteChat(chatId) {
        if (!window.go || !window.go.main || !window.go.main.App) {
            console.error('Backend not available for chat deletion');
            return false;
        }

        try {
            await window.go.main.App.DeleteChat(chatId);

            // If deleted chat was current, reset
            if (this.currentChatId === chatId) {
                this.currentChatId = null;
                if (typeof contextDisplay !== 'undefined') {
                    contextDisplay.resetToDefaultState();
                    await contextDisplay.updateContextDisplay();
                }
            }

            console.log('Chat deleted:', chatId);
            return true;
        } catch (error) {
            console.error('Failed to delete chat:', error);
            return false;
        }
    }

    /**
     * Get chat list
     */
    async getChatList() {
        if (!window.go || !window.go.main || !window.go.main.App) {
            console.log('Backend not available for chat list');
            return [];
        }

        try {
            const chatList = await window.go.main.App.GetChatList();
            return chatList || [];
        } catch (error) {
            console.error('Failed to get chat list:', error);
            return [];
        }
    }

    /**
     * Sync frontend state to backend
     */
    async syncToBackend() {
        try {
            const chatData = this.collectCompleteChatData();

            // Sync world
            if (chatData.world) {
                await this.syncWorldToBackend(chatData.world);
            }

            // Sync scenario
            if (chatData.scenario) {
                await this.syncScenarioToBackend(chatData.scenario);
            }

            // Sync characters
            if (chatData.characters && chatData.characters.length > 0) {
                await this.syncCharactersToBackend(chatData.characters);
            }

            // Sync lorebooks
            if (chatData.lorebooks && chatData.lorebooks.length > 0) {
                await this.syncLorebooksToBackend(chatData.lorebooks);
            }

            console.log('Frontend state synced to backend');
        } catch (error) {
            console.error('Failed to sync to backend:', error);
        }
    }

    /**
     * Sync world to backend
     */
    async syncWorldToBackend(world) {
        if (window.go && window.go.main && window.go.main.App && window.go.main.App.AddWorldToChat) {
            try {
                await window.go.main.App.AddWorldToChat(world);
            } catch (error) {
                console.warn('World sync failed:', error);
            }
        }
    }

    /**
     * Sync scenario to backend
     */
    async syncScenarioToBackend(scenario) {
        if (window.go?.main?.App?.AddScenarioToChat) {
            try {
                await window.go.main.App.AddScenarioToChat(scenario);
            } catch (error) {
                console.warn('Scenario sync failed:', error);
            }
        }
    }

    /**
     * Sync characters to backend
     */
    async syncCharactersToBackend(characters) {
        if (!window.go?.main?.App?.AddCharacterToChat) return;
        // Clear stale characters from previous session before re-syncing
        if (window.go.main.App.ResetSessionCharacters) {
            try { await window.go.main.App.ResetSessionCharacters(); } catch (_) {}
        }
        for (const character of characters) {
            try {
                await window.go.main.App.AddCharacterToChat(character);
            } catch (error) {
                console.warn('Character sync failed:', character?.name, error);
            }
        }
    }

    /**
     * Sync lorebooks to backend
     */
    async syncLorebooksToBackend(lorebooks) {
        if (window.go && window.go.main && window.go.main.App && window.go.main.App.SetLorebooks) {
            try {
                await window.go.main.App.SetLorebooks(lorebooks);
            } catch (error) {
                console.warn('Lorebooks sync failed:', error);
            }
        }
    }

    /**
     * Create memory entry for character swap (stored in memory system, NOT lorebooks)
     */
    async createMemoryLorebook(oldCharacters, newCharacters) {
        if (!oldCharacters || oldCharacters.length === 0) {
            return;
        }

        const cd = window.sidebarState?.currentChatData || {};
        const worldName = cd.world?.name || cd.world?.title || 'Unknown World';
        const scenarioName = cd.scenario?.name || cd.scenario?.title || 'Unknown Scenario';

        oldCharacters.forEach(char => {
            const summary = `${char.name} was removed from the chat`;
            const content = [
                `Character "${char.name}" was removed from the chat.`,
                `Role: ${char.role || 'Unknown'}`,
                `Personality: ${char.personality || 'Unknown'}`,
                `Backstory: ${char.backstory || 'Unknown'}`,
                `World: ${worldName}`,
                `Scenario: ${scenarioName}`,
                `Characters now present: ${(cd.characters || []).filter(c => c.id !== char.id).map(c => c.name).join(', ') || 'None'}`
            ].join('\n');

            if (window.MemoryStore) {
                window.MemoryStore.createMemoryEntry({
                    title: `Character Swap: ${char.name}`,
                    summary: summary,
                    content: content,
                    scope: 'chat',
                    ownerId: cd.sessionId || 'default',
                    tags: ['character-swap', 'continuity', char.name.toLowerCase().replace(/\s+/g, '-')],
                    source: 'auto',
                    trigger: 'character_swap',
                    confidence: 100
                });
            }
        });

        console.log('Memory entries created for character swap');
    }

    /**
     * Cleanup on page unload
     */
    cleanup() {
        if (this.autosaveInterval) {
            clearInterval(this.autosaveInterval);
            this.autosaveInterval = null;
        }

        // Perform final save
        if (this.currentChatId) {
            this.saveCurrentChatData().catch(error => {
                console.error('Final save failed:', error);
            });
        }
    }
}

// Create global instance
const dataSync = new DataSync();
window.dataSync = dataSync;

// Initialize when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    dataSync.initialize();
});

// Cleanup on page unload
window.addEventListener('beforeunload', () => {
    dataSync.cleanup();
});