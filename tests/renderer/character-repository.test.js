/**
 * Character repository tests.
 *
 * These execute the real frontend/dist/js/character-repository.js in a VM with
 * a fake Wails bridge and localStorage. They prove the cutover contract:
 *   - the backend list beats a conflicting cache,
 *   - a pre-cutover raw array is imported through ImportCharacters and replaced,
 *   - a save/delete reaches the backend before the projection or cache changes,
 *   - a refused write leaves both the projection and the cache untouched,
 *   - the versioned cache is used only when the backend is unreachable.
 */

const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');

const SOURCE = fs.readFileSync(
    path.resolve(__dirname, '../../frontend/dist/js/character-repository.js'),
    'utf8'
);
const CACHE_KEY = 'howl-chat-characters';

function createHarness(options = {}) {
    const store = new Map();
    if (options.cache !== undefined) {
        store.set(CACHE_KEY, typeof options.cache === 'string' ? options.cache : JSON.stringify(options.cache));
    }

    const operations = [];
    let backendList = options.backend ? options.backend.map(c => ({ ...c })) : null;
    let failWrites = options.failWrites === true;
    let onWrite = options.onWrite || null;

    const app = {
        async GetCharacters() {
            operations.push('GetCharacters');
            if (backendList === null) throw new Error('backend unavailable');
            return JSON.parse(JSON.stringify(backendList));
        },
        async SaveCharacter(character) {
            operations.push('SaveCharacter');
            if (onWrite) onWrite(character);
            if (failWrites) throw new Error('disk full');
            const saved = { ...character, id: character.id || 'chr_created', revision: 'rev-' + (backendList || []).length };
            backendList = (backendList || []).filter(c => c.id !== saved.id).concat(saved);
            operations.push('SaveCharacter:ack');
            return JSON.parse(JSON.stringify(saved));
        },
        async DeleteCharacter(id) {
            operations.push('DeleteCharacter');
            if (failWrites) throw new Error('disk full');
            backendList = (backendList || []).filter(c => c.id !== id);
            operations.push('DeleteCharacter:ack');
        },
        async ImportCharacters(list) {
            operations.push('ImportCharacters');
            backendList = (backendList || []).slice();
            for (const item of list) {
                const id = item.id || 'chr_imported';
                if (!backendList.some(c => c.id === id)) backendList.push({ ...item, id });
            }
            return [];
        },
    };

    const localStorage = {
        getItem(key) { return store.has(key) ? store.get(key) : null; },
        setItem(key, value) { operations.push('cache:' + key); store.set(key, String(value)); },
        removeItem(key) { store.delete(key); },
    };

    const events = [];
    const window = {
        go: options.noBackend ? {} : { main: { App: app } },
        addEventListener() {},
        dispatchEvent(event) { events.push(event.type); },
        showNotification() {},
    };
    window.window = window;

    const sandbox = {
        window,
        localStorage,
        console: { log() {}, warn() {}, error() {} },
        CustomEvent: class CustomEvent {
            constructor(type, init) { this.type = type; this.detail = init && init.detail; }
        },
        setTimeout,
        clearTimeout,
        Date,
    };
    vm.createContext(sandbox);
    vm.runInContext(SOURCE, sandbox, { filename: 'character-repository.js' });

    return {
        state: window.CharacterStorage,
        store,
        operations,
        events,
        setFail(value) { failWrites = value; },
        setOnWrite(fn) { onWrite = fn; },
    };
}

test('a conflicting versioned cache cannot override the backend list', async () => {
    const harness = createHarness({
        backend: [{ id: 'chr_mira', name: 'Mira', description: 'backend' }],
        cache: {
            version: 2,
            savedAt: '2026-01-01T00:00:00.000Z',
            characters: [{ id: 'chr_mira', name: 'Mira', description: 'cache' }],
        },
    });

    await harness.state.ready;

    assert.equal(harness.state.source, 'backend');
    assert.equal(harness.state.all()[0].description, 'backend');
    const cached = JSON.parse(harness.store.get(CACHE_KEY));
    assert.equal(cached.version, 2);
    assert.equal(cached.characters[0].description, 'backend');
});

test('a pre-cutover raw array is imported once and replaced by the envelope', async () => {
    const harness = createHarness({
        backend: [],
        cache: JSON.stringify([{ id: 'chr_legacy', name: 'Legacy' }]),
    });

    await harness.state.ready;

    assert.ok(harness.operations.indexOf('ImportCharacters') < harness.operations.indexOf('GetCharacters'));
    assert.equal(harness.state.all().some(c => c.id === 'chr_legacy'), true);
    const cached = JSON.parse(harness.store.get(CACHE_KEY));
    assert.equal(cached.version, 2);
    assert.equal(Array.isArray(cached), false);
});

test('a save reaches the backend before the projection or cache changes', async () => {
    let observedDuringWrite = null;
    let cacheWritesDuringWrite = 0;
    const harness = createHarness({
        backend: [{ id: 'chr_mira', name: 'Mira', description: 'old', revision: 'rev-old' }],
        onWrite() {
            observedDuringWrite = harness.state.all()[0].description;
            cacheWritesDuringWrite = harness.operations.filter(op => op.startsWith('cache:')).length;
        },
    });
    harness.setOnWrite(() => {
        observedDuringWrite = harness.state.all()[0].description;
        cacheWritesDuringWrite = harness.operations.filter(op => op.startsWith('cache:')).length;
    });

    await harness.state.ready;
    harness.operations.length = 0;

    await harness.state.save({ id: 'chr_mira', description: 'new' });

    assert.equal(observedDuringWrite, 'old', 'the projection must still show the acknowledged value during the write');
    assert.equal(cacheWritesDuringWrite, 0, 'no cache write may precede the acknowledgement');

    const saveIndex = harness.operations.indexOf('SaveCharacter');
    const ackIndex = harness.operations.indexOf('SaveCharacter:ack');
    const cacheIndex = harness.operations.findIndex(op => op.startsWith('cache:'));
    assert.ok(saveIndex >= 0 && ackIndex > saveIndex && cacheIndex > ackIndex);
    assert.equal(harness.state.get('chr_mira').description, 'new');
    assert.equal(JSON.parse(harness.store.get(CACHE_KEY)).characters[0].description, 'new');
});

test('a refused save leaves the projection and cache untouched', async () => {
    const harness = createHarness({
        backend: [{ id: 'chr_mira', name: 'Mira', description: 'old', revision: 'rev-old' }],
    });
    await harness.state.ready;
    const cacheBefore = harness.store.get(CACHE_KEY);
    harness.operations.length = 0;
    harness.setFail(true);

    await assert.rejects(() => harness.state.save({ id: 'chr_mira', description: 'new' }));

    assert.equal(harness.state.get('chr_mira').description, 'old');
    assert.equal(harness.store.get(CACHE_KEY), cacheBefore);
    assert.equal(harness.operations.some(op => op.startsWith('cache:')), false);
});

test('delete reaches the backend before the projection changes', async () => {
    const harness = createHarness({
        backend: [{ id: 'chr_mira', name: 'Mira', revision: 'rev-old' }],
    });
    await harness.state.ready;
    harness.operations.length = 0;

    await harness.state.remove('chr_mira');

    const deleteIndex = harness.operations.indexOf('DeleteCharacter');
    const ackIndex = harness.operations.indexOf('DeleteCharacter:ack');
    const cacheIndex = harness.operations.findIndex(op => op.startsWith('cache:'));
    assert.ok(deleteIndex >= 0 && ackIndex > deleteIndex && cacheIndex > ackIndex);
    assert.equal(harness.state.get('chr_mira'), null);
});

test('no renderer module writes the character cache outside the repository', () => {
    const directory = path.resolve(__dirname, '../../frontend/dist/js');
    const offenders = fs.readdirSync(directory)
        .filter(name => name.endsWith('.js'))
        .filter(name => fs.readFileSync(path.join(directory, name), 'utf8')
            .includes("localStorage.setItem('howl-chat-characters'"));
    assert.deepEqual(offenders, []);
    assert.equal(fs.existsSync(path.join(directory, 'character-storage.js')), false);
});

test('the versioned cache is used only when the backend is unreachable', async () => {
    const harness = createHarness({
        noBackend: true,
        cache: {
            version: 2,
            savedAt: '2026-01-01T00:00:00.000Z',
            characters: [{ id: 'chr_cached', name: 'Cached' }],
        },
    });

    await harness.state.ready;

    assert.equal(harness.state.source, 'cache');
    assert.equal(harness.state.all()[0].name, 'Cached');
});
