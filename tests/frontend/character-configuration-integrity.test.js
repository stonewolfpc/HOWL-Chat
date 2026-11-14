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

function createElement(id) {
    return {
        id,
        value: '',
        checked: false,
        style: {},
        addEventListener() {},
        reset() {}
    };
}

function fillForm(page, record, overrides = {}) {
    const values = {
        'character-name': record.name || '',
        'character-age': record.metadata?.age || '',
        'character-role': record.metadata?.role || '',
        'character-priority': String(record.metadata?.priority ?? 3),
        'character-description': record.description || '',
        'character-personality': record.personality || '',
        'character-background': record.metadata?.background || '',
        'character-intro': record.first_mes || '',
        'character-dialogue-examples': record.mes_example || '',
        'char-voice-enabled': !!record.howl?.voice?.enabled,
        'char-voice-engine': record.howl?.voice?.engine || '',
        'char-voice-model': record.howl?.voice?.voice || '',
        'char-voice-speed': String(record.howl?.voice?.speed ?? 1.0),
        'char-voice-volume': String(record.howl?.voice?.volume ?? 0.8),
        ...overrides
    };

    for (const [id, value] of Object.entries(values)) {
        if (!page.elements.has(id)) page.elements.set(id, createElement(id));
        const el = page.elements.get(id);
        if (id === 'char-voice-enabled') el.checked = !!value;
        else el.value = value;
    }
}

function setStateForEdit(page, record, overrides = {}) {
    page.state.characterCache = [clone(record)];
    page.state.currentEditingIndex = 0;
    page.state.characterImages = clone(record.howl?.images || []);
    page.state.selectedLorebooks = clone(record.howl?.lorebooks || record.lorebooks || []);
    page.state.resetDirtyFields();
    fillForm(page, record, overrides);
}

function createCharacterPage(seedCharacters, shared = {}) {
    const localStorage = shared.localStorage || new MemoryStorage();
    const indexedDB  = shared.indexedDB || createIndexedDB();

    if (seedCharacters && !localStorage.getItem(CHARACTER_KEY)) {
        localStorage.setItem(CHARACTER_KEY, JSON.stringify(seedCharacters));
    }

    const elements = new Map();
    elements.set('character-creation-overlay', createElement('character-creation-overlay'));
    const form = { reset() {} };
    const document = {
        getElementById(id) { return elements.get(id) || null; },
        querySelector(selector) { return selector === '.character-creation-form' ? form : null; },
        addEventListener() {}
    };

    const state = {
        characterImages: [],
        selectedLorebooks: [],
        currentEditingIndex: -1,
        characterCache: null,
        dirty: new Set(),
        markFieldDirty(id) { this.dirty.add(id); },
        isFieldDirty(id) { return this.dirty.has(id); },
        resetDirtyFields() { this.dirty.clear(); },
        lorebookKey(lorebook) {
            const id = (lorebook?.id || '').toString().trim();
            return id || this.lorebookDisplayName(lorebook).toLowerCase();
        },
        lorebookDisplayName(lorebook) {
            return (lorebook?.title || lorebook?.name || '').trim() || 'Untitled Lorebook';
        }
    };

    const alerts = [];
    const sandbox = {
        console,
        document,
        localStorage,
        indexedDB,
        state,
        alert(message) { alerts.push(message); },
        setTimeout,
        loadCharacters() {},
        parent: null
    };
    sandbox.window = sandbox;
    sandbox.globalThis = sandbox;
    sandbox.characterState = state;

    vm.createContext(sandbox);
    vm.runInContext(
        fs.readFileSync(path.join(__dirname, '..', '..', 'frontend', 'dist', 'js', 'character-storage.js'), 'utf8'),
        sandbox,
        { filename: 'character-storage.js' }
    );
    pageStorageSetup(sandbox);
    vm.runInContext(
        fs.readFileSync(path.join(__dirname, '..', '..', 'frontend', 'dist', 'js', 'character-form.js'), 'utf8'),
        sandbox,
        { filename: 'character-form.js' }
    );

    return {
        sandbox,
        storage: sandbox.CharacterStorage,
        state,
        elements,
        alerts,
        localStorage,
        indexedDB
    };
}

function pageStorageSetup(sandbox) {
    // CharacterStorage assigns window.CharacterStorage itself; this hook keeps the
    // harness explicit about the handoff between the two page modules.
    if (!sandbox.CharacterStorage && sandbox.window.CharacterStorage) {
        sandbox.CharacterStorage = sandbox.window.CharacterStorage;
    }
}

const FIXTURE = {
    id: 'char_fixture_001',
    name: 'Aria',
    description: 'Aria description',
    personality: 'Calm and watchful',
    scenario: 'A scenario that must survive',
    first_mes: 'Hello there',
    mes_example: '<START>',
    avatar: 'data:image/png;base64,avatar',
    tags: ['alpha', 'beta'],
    metadata: {
        version: 2,
        created: '2025-01-01T00:00:00Z',
        updated: '2025-01-02T00:00:00Z',
        seed: 123456,
        age: '30',
        role: 'Ranger',
        background: 'Forest',
        priority: 4,
        extraMetadata: { keep: true }
    },
    howl: {
        world: { id: 'world_1', title: 'World One', nested: { keep: 'yes' } },
        lorebooks: [{
            id: 'lb_1',
            name: 'Lore One',
            title: 'Lore One',
            description: 'Lore description',
            extraRef: { keep: true }
        }],
        memories: [{ id: 'mem_1', text: 'A memory' }],
        notes: 'Important notes',
        image_style: 'custom-style',
        images: [{
            id: 'img_1',
            name: 'aria.png',
            size: '1 KB',
            data: 'data:image/png;base64,AAAA',
            medium: 'data:image/png;base64,BBBB',
            original: 'data:image/png;base64,CCCC',
            isThumbnail: true,
            extraImage: { keep: true }
        }],
        voice: {
            enabled: true,
            engine: 'piper',
            voice: 'custom_voice_not_in_the_fixed_select',
            speed: 1.1,
            volume: 0.7,
            noiseScale: 0.5,
            noiseW: 0.9,
            sentenceSilence: 0.25,
            extraVoice: { keep: true }
        },
        extraHowl: { keep: true }
    },
    extraTop: { keep: true }
};

test('Character Tab edit preserves unrelated, nested, unknown, and image data through save and reload', async () => {
    const page = createCharacterPage([clone(FIXTURE)]);
    setStateForEdit(page, FIXTURE, {
        'character-name': 'Aria Updated',
        'char-voice-model': ''
    });

    await page.sandbox.saveCharacter();

    const raw = JSON.parse(page.localStorage.getItem(CHARACTER_KEY))[0];
    assert.equal(raw.id, FIXTURE.id);
    assert.equal(raw.name, 'Aria Updated');
    assert.equal(raw.description, FIXTURE.description);
    assert.equal(raw.personality, FIXTURE.personality);
    assert.equal(raw.scenario, FIXTURE.scenario);
    assert.deepEqual(raw.tags, FIXTURE.tags);
    assert.deepEqual(raw.metadata.extraMetadata, { keep: true });
    assert.equal(raw.metadata.created, FIXTURE.metadata.created);
    assert.equal(raw.metadata.seed, FIXTURE.metadata.seed);
    assert.deepEqual(raw.howl.world, FIXTURE.howl.world);
    assert.deepEqual(raw.howl.memories, FIXTURE.howl.memories);
    assert.equal(raw.howl.notes, FIXTURE.howl.notes);
    assert.equal(raw.howl.image_style, FIXTURE.howl.image_style);
    assert.deepEqual(raw.howl.extraHowl, { keep: true });
    assert.deepEqual(raw.extraTop, { keep: true });
    assert.deepEqual(raw.howl.lorebooks[0].extraRef, { keep: true });
    assert.equal(raw.howl.voice.voice, FIXTURE.howl.voice.voice);
    assert.equal(raw.howl.voice.noiseScale, FIXTURE.howl.voice.noiseScale);
    assert.equal(raw.howl.voice.noiseW, FIXTURE.howl.voice.noiseW);
    assert.equal(raw.howl.voice.sentenceSilence, FIXTURE.howl.voice.sentenceSilence);
    assert.deepEqual(raw.howl.voice.extraVoice, { keep: true });
    assert.equal(raw.howl.images[0].id, FIXTURE.howl.images[0].id);
    assert.equal(raw.howl.images[0].data, undefined);

    // A fresh CharacterStorage/page instance reads the persisted record and
    // rehydrates the image blobs from IndexedDB.
    const restarted = createCharacterPage(null, {
        localStorage: page.localStorage,
        indexedDB: page.indexedDB
    });
    const loaded = await restarted.storage.loadAll();
    assert.equal(loaded.length, 1);
    assert.equal(loaded[0].name, 'Aria Updated');
    assert.equal(loaded[0].howl.world.nested.keep, 'yes');
    assert.equal(loaded[0].howl.voice.noiseScale, 0.5);
    assert.equal(loaded[0].howl.voice.extraVoice.keep, true);
    assert.equal(loaded[0].howl.images[0].data, FIXTURE.howl.images[0].data);
    assert.equal(loaded[0].howl.images[0].medium, FIXTURE.howl.images[0].medium);
    assert.equal(loaded[0].howl.images[0].extraImage.keep, true);
});

test('a partial sidebar-style storage save preserves the persisted character and image blobs', async () => {
    const page = createCharacterPage([clone(FIXTURE)]);
    await page.storage.save(clone(FIXTURE));

    // This is the shape used by sidebar image updates: a derived character view
    // with no metadata/howl sections and no images field.
    await page.storage.save({ id: FIXTURE.id, image: 'data:image/png;base64,NEW' });

    const raw = JSON.parse(page.localStorage.getItem(CHARACTER_KEY))[0];
    assert.equal(raw.name, FIXTURE.name);
    assert.equal(raw.howl.notes, FIXTURE.howl.notes);
    assert.equal(raw.howl.voice.noiseScale, FIXTURE.howl.voice.noiseScale);
    assert.deepEqual(raw.extraTop, { keep: true });

    const loaded = await page.storage.loadAll();
    assert.equal(loaded[0].howl.images[0].data, FIXTURE.howl.images[0].data);
    assert.equal(loaded[0].howl.images[0].medium, FIXTURE.howl.images[0].medium);
});

test('an omitted nested property is preserved while an explicit null is applied', async () => {
    const page = createCharacterPage([clone(FIXTURE)]);
    await page.storage.save(clone(FIXTURE));

    // Omission must leave the existing non-null nested value untouched.
    await page.storage.save({
        id: FIXTURE.id,
        name: 'Null Semantics',
        howl: { voice: { speed: 1.1 } }
    });
    let loaded = await page.storage.loadAll();
    assert.equal(loaded[0].howl.voice.noiseScale, FIXTURE.howl.voice.noiseScale);
    assert.equal(loaded[0].howl.voice.speed, 1.1);
    assert.equal(loaded[0].howl.voice.extraVoice.keep, true);

    // An explicit null is authoritative, not equivalent to omission.
    await page.storage.save({ id: FIXTURE.id, howl: { voice: { noiseScale: null } } });
    loaded = await page.storage.loadAll();
    assert.equal(loaded[0].howl.voice.noiseScale, null);
    assert.equal(loaded[0].howl.voice.speed, 1.1);
});

test('an explicit image-list clear still removes the persisted image blobs', async () => {
    const page = createCharacterPage([clone(FIXTURE)]);
    await page.storage.save(clone(FIXTURE));

    await page.storage.save({ id: FIXTURE.id, howl: { images: [] } });

    const loaded = await page.storage.loadAll();
    assert.equal(loaded[0].howl.images.length, 0);
    assert.equal(loaded[0].howl.notes, FIXTURE.howl.notes);
});

test('sequential edits preserve other sections and explicit scalar clears remain authoritative', async () => {
    const page = createCharacterPage([clone(FIXTURE)]);
    setStateForEdit(page, FIXTURE, {
        'character-name': 'Aria First Edit',
        'char-voice-model': ''
    });
    await page.sandbox.saveCharacter();

    const afterFirst = JSON.parse(page.localStorage.getItem(CHARACTER_KEY))[0];
    const loaded = await page.storage.loadAll();

    setStateForEdit(page, loaded[0], {
        'character-description': '',
        'character-role': 'Updated Role',
        'character-age': '31',
        'char-voice-model': ''
    });
    await page.sandbox.saveCharacter();

    const afterSecond = JSON.parse(page.localStorage.getItem(CHARACTER_KEY))[0];
    assert.equal(afterSecond.name, 'Aria First Edit');
    assert.equal(afterSecond.description, '');
    assert.equal(afterSecond.metadata.role, 'Updated Role');
    assert.equal(afterSecond.metadata.age, '31');
    assert.deepEqual(afterSecond.metadata.extraMetadata, { keep: true });
    assert.equal(afterSecond.metadata.background, FIXTURE.metadata.background);
    assert.equal(afterSecond.howl.notes, FIXTURE.howl.notes);
    assert.equal(afterSecond.howl.voice.noiseScale, FIXTURE.howl.voice.noiseScale);
    assert.deepEqual(afterSecond.howl.extraHowl, { keep: true });
    assert.equal(afterFirst.id, afterSecond.id);
});

test('an explicit voice-model clear works while an unrepresentable stored voice is preserved', async () => {
    const page = createCharacterPage([clone(FIXTURE)]);
    setStateForEdit(page, FIXTURE, {
        'character-name': 'Voice Preserve',
        'char-voice-model': ''
    });
    await page.sandbox.saveCharacter();

    let raw = JSON.parse(page.localStorage.getItem(CHARACTER_KEY))[0];
    assert.equal(raw.howl.voice.voice, FIXTURE.howl.voice.voice);
    assert.equal(raw.howl.voice.noiseScale, FIXTURE.howl.voice.noiseScale);

    const loaded = await page.storage.loadAll();
    setStateForEdit(page, loaded[0], {
        'character-name': 'Voice Clear',
        'char-voice-model': ''
    });
    page.state.markFieldDirty('char-voice-model');
    await page.sandbox.saveCharacter();

    raw = JSON.parse(page.localStorage.getItem(CHARACTER_KEY))[0];
    assert.equal(raw.howl.voice.voice, '');
    assert.equal(raw.howl.voice.noiseScale, FIXTURE.howl.voice.noiseScale);
    assert.equal(raw.howl.voice.extraVoice.keep, true);
});

test('legacy records without newer optional fields remain editable without destructive migration', async () => {
    const legacy = {
        id: 'legacy_character_001',
        name: 'Legacy',
        howl: { notes: 'keep legacy notes' },
        metadata: { age: '42' }
    };
    const page = createCharacterPage([clone(legacy)]);
    setStateForEdit(page, legacy, {
        'character-name': 'Legacy Updated'
    });
    await page.sandbox.saveCharacter();

    const raw = JSON.parse(page.localStorage.getItem(CHARACTER_KEY))[0];
    assert.equal(raw.id, legacy.id);
    assert.equal(raw.name, 'Legacy Updated');
    assert.equal(raw.metadata.age, '42');
    assert.equal(raw.howl.notes, 'keep legacy notes');
    assert.equal(raw.howl.voice.enabled, false);
});
