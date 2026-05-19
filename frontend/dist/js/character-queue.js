/**
 * Character Queue Manager
 * 
 * Handles character turn-taking for multi-character conversations.
 * Manages character rotation, selection, and state tracking.
 * 
 * Features:
 * - FIFO character rotation
 * - Current character tracking
 * - Queue pause/resume
 * - Character filtering and exclusion
 * 
 * @module CharacterQueue
 * @version 1.0.0
 */

class CharacterQueue {
    /**
     * Initialize character queue
     */
    constructor() {
        this.queue = [];
        this.currentIndex = 0;
        this.currentCharacter = null;
        this.isActive = false;
        this.pausedCharacters = new Set();
        this.respondedCharacters = new Set();
    }

    /**
     * Initialize queue from loaded characters
     * @param {Array} characters - Array of character objects
     */
    initialize(characters = []) {
        const active = characters.filter(c => c.status !== 'dead' && c.status !== 'stepped-out');
        this.queue = active;
        this.currentIndex = 0;
        this.currentCharacter = this.queue.length > 0 ? this.queue[0] : null;
        this.respondedCharacters.clear();
        const skipped = characters.length - active.length;
        if (skipped > 0) {
            console.log(`⏭️ CharacterQueue filtered out ${skipped} inactive character(s)`);
        }
        console.log(`📋 Character queue initialized with ${this.queue.length} characters`);
        this.logQueueState();
    }

    /**
     * Add character to queue
     * @param {Object} character - Character object to add
     */
    addCharacter(character) {
        if (character.status === 'dead' || character.status === 'stepped-out') {
            console.log(`⏭️ Skipped adding inactive character to queue: ${character.name} (${character.status})`);
            return;
        }
        if (!this.queue.find(c => c.id === character.id)) {
            this.queue.push(character);
            console.log(`➕ Added character to queue: ${character.name}`);
            this.logQueueState();
        }
    }

    /**
     * Remove character from queue
     * @param {String} characterId - ID of character to remove
     */
    removeCharacter(characterId) {
        const index = this.queue.findIndex(c => c.id === characterId);
        if (index !== -1) {
            const removed = this.queue[index];
            this.queue.splice(index, 1);
            this.respondedCharacters.delete(removed.id);
            this.pausedCharacters.delete(removed.id);
            
            // Adjust current index if needed
            if (this.currentIndex >= this.queue.length) {
                this.currentIndex = Math.max(0, this.queue.length - 1);
            }
            
            this.currentCharacter = this.queue[this.currentIndex] || null;
            console.log(`➖ Removed character from queue: ${removed.name}`);
            this.logQueueState();
        }
    }

    /**
     * Get current character
     * @returns {Object|null} Current character object
     */
    getCurrentCharacter() {
        return this.currentCharacter;
    }

    /**
     * Get current character ID
     * @returns {String|null} Current character ID
     */
    getCurrentCharacterId() {
        return this.currentCharacter ? this.currentCharacter.id : null;
    }

    /**
     * Get current character name
     * @returns {String|null} Current character name
     */
    getCurrentCharacterName() {
        return this.currentCharacter ? this.currentCharacter.name : null;
    }

    /**
     * Rotate to next character in queue
     * @returns {Object|null} Next character, or null if all have responded
     */
    rotateToNext() {
        if (this.queue.length === 0) {
            console.warn('⚠️ No characters in queue');
            return null;
        }

        // Find next character that hasn't responded yet
        let attempts = 0;
        const maxAttempts = this.queue.length;

        while (attempts < maxAttempts) {
            this.currentIndex = (this.currentIndex + 1) % this.queue.length;
            const candidate = this.queue[this.currentIndex];

            // Skip if character has already responded, is paused, or is dead/stepped-out
            const isInactive = candidate.status === 'dead' || candidate.status === 'stepped-out';
            if (!this.respondedCharacters.has(candidate.id) && 
                !this.pausedCharacters.has(candidate.id) &&
                !isInactive) {
                this.currentCharacter = candidate;
                console.log(`🔄 Rotated to: ${this.getCurrentCharacterName()}`);
                return this.currentCharacter;
            }
            if (isInactive) {
                console.log(`⏭️ Skipping inactive character: ${candidate.name} (${candidate.status})`);
            }

            attempts++;
        }

        // All characters have responded
        console.log('✅ All characters have responded in this round');
        return null;
    }

    /**
     * Mark character as responded
     * @param {String} characterId - ID of responding character
     */
    markResponded(characterId) {
        this.respondedCharacters.add(characterId);
        console.log(`✅ Marked as responded: ${characterId}`);
    }

    /**
     * Reset response tracking for next round
     */
    resetRound() {
        this.respondedCharacters.clear();
        this.currentIndex = 0;
        this.currentCharacter = this.queue.length > 0 ? this.queue[0] : null;
        console.log('🔁 Response tracking reset for new round');
    }

    /**
     * Pause character from responding
     * @param {String} characterId - ID of character to pause
     */
    pauseCharacter(characterId) {
        this.pausedCharacters.add(characterId);
        console.log(`⏸️ Paused character: ${characterId}`);
    }

    /**
     * Resume paused character
     * @param {String} characterId - ID of character to resume
     */
    resumeCharacter(characterId) {
        this.pausedCharacters.delete(characterId);
        console.log(`▶️ Resumed character: ${characterId}`);
    }

    /**
     * Pause all characters
     */
    pauseAll() {
        this.isActive = false;
        console.log('⏸️ All characters paused');
    }

    /**
     * Resume all characters
     */
    resumeAll() {
        this.isActive = true;
        console.log('▶️ All characters resumed');
    }

    /**
     * Get list of all characters
     * @returns {Array} Array of character objects
     */
    getAllCharacters() {
        return [...this.queue];
    }

    /**
     * Get list of active (not paused) characters
     * @returns {Array} Array of active character objects
     */
    getActiveCharacters() {
        return this.queue.filter(c => !this.pausedCharacters.has(c.id));
    }

    /**
     * Get list of characters that have responded
     * @returns {Array} Array of responded character objects
     */
    getRespondedCharacters() {
        return this.queue.filter(c => this.respondedCharacters.has(c.id));
    }

    /**
     * Get list of characters that haven't responded yet
     * @returns {Array} Array of pending character objects
     */
    getPendingCharacters() {
        return this.queue.filter(c => !this.respondedCharacters.has(c.id) &&
                                       !this.pausedCharacters.has(c.id));
    }

    /**
     * Check if all characters have responded
     * @returns {Boolean} True if all have responded
     */
    allHaveResponded() {
        return this.getPendingCharacters().length === 0;
    }

    /**
     * Get queue size
     * @returns {Number} Number of characters in queue
     */
    size() {
        return this.queue.length;
    }

    /**
     * Get active queue size
     * @returns {Number} Number of active characters
     */
    activeSize() {
        return this.getActiveCharacters().length;
    }

    /**
     * Clear entire queue
     */
    clear() {
        this.queue = [];
        this.currentIndex = 0;
        this.currentCharacter = null;
        this.respondedCharacters.clear();
        this.pausedCharacters.clear();
        console.log('🗑️ Character queue cleared');
    }

    /**
     * Log current queue state for debugging
     */
    logQueueState() {
        console.log(`📊 Queue State:`, {
            total: this.queue.length,
            active: this.getActiveCharacters().length,
            responded: this.respondedCharacters.size,
            pending: this.getPendingCharacters().length,
            current: this.getCurrentCharacterName(),
            paused: this.pausedCharacters.size
        });
    }
}

// Export singleton instance
window.CharacterQueue = new CharacterQueue();
