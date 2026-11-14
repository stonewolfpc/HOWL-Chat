'use strict';

const test = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');

const MEMORY_KEY = 'howl-chat-memories';

class MemoryStorage {
    constructor(seed = []) {
        this.values = new Map();
        this.values.set(MEMORY_KEY, JSON.stringify(seed));
    }
    getItem(key) {
        return this.values.has(key) ? this.values.get(key) : null;
    }
    setItem(key, value) {
        this.values.set(key, String(value));
    }
    removeItem(key) {
        this.values.delete(key);
    }
}

function loadMemoryPage(seed) {
    const localStorage = new MemoryStorage(seed);
    const sandbox = {
        console,
        localStorage,
        document: {
            querySelector() { return null; },
            createElement() { return {}; }
        }
    };
    sandbox.window = sandbox;
    // The page normally gets this global from sidebar-characters.js; the module
    // itself exports it at load time. Supply the page-order dependency here.
    sandbox.viewCharacterMemories = () => {};
    vm.createContext(sandbox);
    vm.runInContext(
        fs.readFileSync(path.join(__dirname, '..', '..', 'frontend', 'dist', 'js', 'memory-store.js'), 'utf8'),
        sandbox,
        { filename: 'memory-store.js' }
    );
    vm.runInContext(
        fs.readFileSync(path.join(__dirname, '..', '..', 'frontend', 'dist', 'js', 'character-memories.js'), 'utf8'),
        sandbox,
        { filename: 'character-memories.js' }
    );
    return { sandbox, localStorage };
}

test('character memory retrieval accepts only stable-ID ownership', () => {
    const { sandbox } = loadMemoryPage([
        { id: 'stable', scope: 'character', ownerId: 'char_aria', ownerKind: 'character_id' },
        { id: 'legacy-name', scope: 'character', ownerId: 'Aria' },
        { id: 'unmarked-id', scope: 'character', ownerId: 'char_aria' },
        { id: 'foreign', scope: 'character', ownerId: 'char_other', ownerKind: 'character_id' },
        { id: 'chat', scope: 'chat', ownerId: 'chat_1', ownerKind: 'chat' }
    ]);

    const storeMemories = sandbox.MemoryStore.getCharacterMemories('char_aria');
    assert.equal(storeMemories.map(memory => memory.id).join(','), 'stable');

    const viewerMemories = sandbox.getCharacterMemoriesData('char_aria');
    assert.equal(viewerMemories.map(memory => memory.id).join(','), 'stable');
});

test('new character memories retain the explicit ownership kind', () => {
    const { sandbox, localStorage } = loadMemoryPage([]);

    sandbox.MemoryStore.createMemoryEntry({
        title: 'Stable memory',
        scope: 'character',
        ownerId: 'char_aria',
        ownerKind: 'character_id'
    });
    sandbox.MemoryStore.createMemoryEntry({
        title: 'Unmarked memory',
        scope: 'character',
        ownerId: 'char_aria'
    });

    const persisted = JSON.parse(localStorage.getItem(MEMORY_KEY));
    assert.equal(persisted.length, 2);
    assert.equal(persisted[0].ownerKind, 'character_id');
    assert.equal(persisted[1].ownerKind, '');

    const memories = sandbox.MemoryStore.getCharacterMemories('char_aria');
    assert.equal(memories.length, 1);
    assert.equal(memories[0].title, 'Stable memory');
});
