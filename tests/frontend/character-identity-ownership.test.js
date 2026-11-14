'use strict';

const test = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');

const CHARACTER_KEY = 'howl-chat-characters';

function clone(value) {
    return value === undefined ? undefined : JSON.parse(JSON.stringify(value));
}

class MemoryStorage {
    constructor() {
        this.values = new Map();
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

function createRequest(result) {
    const req = {};
    setTimeout(() => {
        if (req.onsuccess) req.onsuccess({ target: { result } });
    }, 0);
    return req;
}

function createIndexedDB() {
    const stores = new Map();
    const db = {
        objectStoreNames: {
            contains(name) { return stores.has(name); }
        },
        createObjectStore(name) {
            if (!stores.has(name)) stores.set(name, new Map());
        },
        transaction(name) {
            if (!stores.has(name)) stores.set(name, new Map());
            const store = stores.get(name);
            return {
                objectStore() {
                    return {
                        put(entry) {
                            store.set(entry.charId, clone(entry));
                            return createRequest(undefined);
                        },
                        get(charId) {
                            return createRequest(store.get(charId));
                        },
                        delete(charId) {
                            store.delete(charId);
                            return createRequest(undefined);
                        }
                    };
                }
            };
        }
    };
    return {
        open() {
            const req = {};
            setTimeout(() => {
                if (!stores.has('images')) {
                    stores.set('images', new Map());
                    if (req.onupgradeneeded) req.onupgradeneeded({ target: { result: db } });
                }
                if (req.onsuccess) req.onsuccess({ target: { result: db } });
            }, 0);
            return req;
        }
    };
}

function createStoragePage(seed, shared = {}) {
    const localStorage = shared.localStorage || new MemoryStorage();
    const indexedDB = shared.indexedDB || createIndexedDB();
    if (seed && !localStorage.getItem(CHARACTER_KEY)) {
        localStorage.setItem(CHARACTER_KEY, JSON.stringify(seed));
    }
    const sandbox = {
        console,
        localStorage,
        indexedDB,
        setTimeout,
        CustomEvent: class { constructor(type, init) { this.type = type; this.detail = init?.detail; } },
        dispatchEvent() {}
    };
    sandbox.window = sandbox;
    vm.createContext(sandbox);
    vm.runInContext(
        fs.readFileSync(path.join(__dirname, '..', '..', 'frontend', 'dist', 'js', 'character-storage.js'), 'utf8'),
        sandbox,
        { filename: 'character-storage.js' }
    );
    return { storage: sandbox.CharacterStorage, localStorage, indexedDB };
}

function createOwnerSelect() {
    const select = {
        dataset: {},
        innerHTML: '',
        selectedOptions: [],
        _value: '',
        type: '',
        addEventListener() {},
        get value() { return this._value; },
        set value(next) {
            this._value = next;
            const matches = [...this.innerHTML.matchAll(/<option value="([^"]*)" data-owner-name="([^"]*)"/g)];
            const found = matches.find(match => match[1] === next);
            this.selectedOptions = found ? [{ dataset: { ownerName: found[2] } }] : [];
        }
    };
    return select;
}

function field(value = '', extra = {}) {
    return {
        value,
        checked: false,
        type: '',
        dataset: {},
        style: {},
        textContent: '',
        selectedOptions: [],
        addEventListener() {},
        ...extra
    };
}

function createLorebooksPage({ characters, migration, localStorage = new MemoryStorage() }) {
    const fields = new Map();
    const owner = createOwnerSelect();
    fields.set('lorebook-owner', owner);
    fields.set('lorebook-scope', field('character'));
    fields.set('lorebook-title', field('Entry'));
    fields.set('lorebook-content', field('body'));
    fields.set('lorebook-visibility', field('private'));
    fields.set('lorebook-enabled', field('', { type: 'checkbox', checked: true }));
    fields.set('lorebook-store-status', field(''));

    const sandbox = {
        console,
        localStorage,
        setTimeout,
        alert() {},
        confirm() { return true; },
        escapeHtml(value) {
            return String(value ?? '').replace(/[&<>"']/g, character => ({
                '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;'
            }[character]));
        },
        CharacterStorage: { loadAll: async () => characters },
        go: {
            main: {
                App: {
                    MigrateLorebookOwners: async () => migration
                }
            }
        },
        document: {
            getElementById(id) { return fields.get(id) || null; },
            addEventListener() {},
            querySelector() { return null; }
        }
    };
    sandbox.window = sandbox;
    sandbox.globalThis = sandbox;
    vm.createContext(sandbox);
    vm.runInContext(
        fs.readFileSync(path.join(__dirname, '..', '..', 'frontend', 'dist', 'js', 'lorebooks.js'), 'utf8'),
        sandbox,
        { filename: 'lorebooks.js' }
    );
    return { sandbox, owner, fields };
}

test('stable character IDs survive edit, rename, restart and deletion', async () => {
    const page = createStoragePage([
        { id: 'char_a', name: 'Aria', howl: {} },
        { id: 'char_b', name: 'Aria', howl: {} }
    ]);

    let loaded = await page.storage.loadAll();
    assert.equal(loaded.map(character => character.id).sort().join(','), 'char_a,char_b');
    assert.equal(loaded.filter(character => character.name === 'Aria').length, 2);

    await page.storage.save({ id: 'char_a', name: 'Aria Renamed' });
    const restarted = createStoragePage(null, {
        localStorage: page.localStorage,
        indexedDB: page.indexedDB
    });
    loaded = await restarted.storage.loadAll();
    const renamed = loaded.find(character => character.id === 'char_a');
    assert.equal(renamed.name, 'Aria Renamed');
    assert.equal(loaded.length, 2, 'a rename must not create a second record');

    await restarted.storage.remove('char_b');
    await restarted.storage.save({ id: 'char_c', name: 'Aria' });
    loaded = await restarted.storage.loadAll();
    assert.equal(loaded.map(character => character.id).sort().join(','), 'char_a,char_c');
    assert.notEqual(loaded.find(character => character.id === 'char_c').id, 'char_b',
        'a recreated name must receive a new identity');
});

test('the lorebook editor stores stable IDs and keeps duplicate names distinguishable', async () => {
    const { sandbox, owner } = createLorebooksPage({
        characters: [
            { id: 'char_aria_111111', name: 'Aria' },
            { id: 'char_aria_222222', name: 'Aria' },
            { id: 'char_bram_333333', name: 'Bram' }
        ],
        migration: { migrated: 0, ambiguous: 0, unknown: 0, failed: 0, items: [] }
    });
    owner.dataset.pendingOwner = 'char_aria_222222';

    await sandbox.loadLorebookOwnerOptions();

    assert.match(owner.innerHTML, /value="char_aria_111111" data-owner-name="Aria"/);
    assert.match(owner.innerHTML, /Aria \(111111\)/);
    assert.match(owner.innerHTML, /Aria \(222222\)/);
    assert.equal(owner.value, 'char_aria_222222');

    const lorebook = sandbox.buildLorebookFromForm();
    assert.equal(lorebook.ownerId, 'char_aria_222222');
    assert.equal(lorebook.ownerName, 'Aria');
});

test('a legacy or deleted owner stays visibly unresolved instead of becoming unowned', async () => {
    const { sandbox, owner } = createLorebooksPage({
        characters: [{ id: 'char_aria_111111', name: 'Aria' }],
        migration: { migrated: 0, ambiguous: 0, unknown: 0, failed: 0, items: [] }
    });
    owner.dataset.pendingOwner = 'Aria';

    await sandbox.loadLorebookOwnerOptions();

    assert.match(owner.innerHTML, /Unresolved owner: Aria/);
    assert.equal(owner.value, 'Aria');
    const lorebook = sandbox.buildLorebookFromForm();
    assert.equal(lorebook.ownerId, 'Aria');
    assert.equal(lorebook.ownerName, '', 'an unresolved label is not promoted to a display name');
});

test('a successful migration reselects the stable ID in the editor', async () => {
    const { sandbox, owner } = createLorebooksPage({
        characters: [{ id: 'char_aria_111111', name: 'Aria' }],
        migration: {
            migrated: 1,
            ambiguous: 0,
            unknown: 0,
            failed: 0,
            items: [{ action: 'migrated', owner_name: 'Aria', owner_id: 'char_aria_111111' }]
        }
    });
    owner.dataset.pendingOwner = 'Aria';

    await sandbox.loadLorebookOwnerOptions();

    assert.equal(owner.value, 'char_aria_111111');
    assert.equal(owner.selectedOptions[0].dataset.ownerName, 'Aria');
    assert.match(sandbox.document.getElementById('lorebook-store-status').textContent, /1 updated/);
});
