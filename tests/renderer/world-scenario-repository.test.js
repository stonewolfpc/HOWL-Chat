/**
 * World and scenario repository tests.
 *
 * These execute the real frontend modules in a VM with a fake Wails bridge and
 * localStorage. They prove the cutover contract for both categories:
 *   - the backend list beats a conflicting cache,
 *   - pre-cutover raw arrays are imported and replaced by the envelope,
 *   - saves/deletes reach the backend before the projection or cache changes,
 *   - a refused write leaves the projection and cache untouched,
 *   - no other renderer module reads or writes either key directly.
 */

const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');

const JS = path.resolve(__dirname, '../../frontend/dist/js');

const CASES = [
    {
        label: 'world',
        script: 'world-repository.js',
        global: 'WorldRepository',
        key: 'howl-worlds',
        list: 'worlds',
        get: 'GetWorlds',
        save: 'SaveWorld',
        remove: 'DeleteWorld',
        import: 'ImportWorlds',
        id: 'wld_mira',
        prefix: 'wld',
    },
    {
        label: 'scenario',
        script: 'scenario-repository.js',
        global: 'ScenarioRepository',
        key: 'howl-chat-scenarios',
        list: 'scenarios',
        get: 'GetScenarios',
        save: 'SaveScenario',
        remove: 'DeleteScenario',
        import: 'ImportScenarios',
        id: 'scn_mira',
        prefix: 'scn',
    },
];

function createHarness(kind, options = {}) {
    const store = new Map();
    if (options.cache !== undefined) {
        store.set(kind.key, typeof options.cache === 'string' ? options.cache : JSON.stringify(options.cache));
    }

    const operations = [];
    let backendList = options.backend ? options.backend.map(item => ({ ...item })) : null;
    let failWrites = options.failWrites === true;
    let onWrite = options.onWrite || null;

    const app = {
        async [kind.get]() {
            operations.push(kind.get);
            if (backendList === null) throw new Error('backend unavailable');
            return JSON.parse(JSON.stringify(backendList));
        },
        async [kind.save](item) {
            operations.push(kind.save);
            if (onWrite) onWrite(item);
            if (failWrites) throw new Error('disk full');
            const saved = { ...item, id: item.id || `${kind.prefix}_created`, revision: 'rev-' + (backendList || []).length };
            backendList = (backendList || []).filter(entry => entry.id !== saved.id).concat(saved);
            operations.push(kind.save + ':ack');
            return JSON.parse(JSON.stringify(saved));
        },
        async [kind.remove](id) {
            operations.push(kind.remove);
            if (failWrites) throw new Error('disk full');
            backendList = (backendList || []).filter(entry => entry.id !== id);
            operations.push(kind.remove + ':ack');
        },
        async [kind.import](items) {
            operations.push(kind.import);
            backendList = backendList || [];
            let imported = 0;
            for (const item of items) {
                if (!backendList.some(entry => entry.id === item.id)) {
                    backendList.push({ ...item, id: item.id || `${kind.prefix}_imported` });
                    imported++;
                }
            }
            return imported;
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
        Math,
    };
    vm.createContext(sandbox);
    const source = fs.readFileSync(path.join(JS, kind.script), 'utf8');
    vm.runInContext(source, sandbox, { filename: kind.script });

    return {
        state: window[kind.global],
        store,
        operations,
        events,
        setFail(value) { failWrites = value; },
        setOnWrite(fn) { onWrite = fn; },
    };
}

for (const kind of CASES) {
    test(`${kind.label}: a conflicting versioned cache cannot override the backend list`, async () => {
        const harness = createHarness(kind, {
            backend: [{ id: kind.id, title: 'backend' }],
            cache: { version: 2, savedAt: '2026-01-01T00:00:00.000Z', [kind.list]: [{ id: kind.id, title: 'cache' }] },
        });

        await harness.state.ready;

        assert.equal(harness.state.source, 'backend');
        assert.equal(harness.state.all()[0].title, 'backend');
        const cached = JSON.parse(harness.store.get(kind.key));
        assert.equal(cached.version, 2);
        assert.equal(cached[kind.list][0].title, 'backend');
    });

    test(`${kind.label}: a pre-cutover raw array stays browser-owned until an explicit import`, async () => {
        const harness = createHarness(kind, {
            backend: [],
            cache: JSON.stringify([{ id: `${kind.prefix}_legacy`, title: 'Legacy' }]),
        });

        await harness.state.ready;

        assert.equal(harness.operations.includes(kind.import), false);
        assert.equal(harness.state.source, 'legacy');
        assert.equal(harness.state.all().some(item => item.id === `${kind.prefix}_legacy`), true);
        assert.equal(Array.isArray(JSON.parse(harness.store.get(kind.key))), true);

        await harness.state.importLegacy();

        assert.ok(harness.operations.indexOf(kind.import) > 0);
        assert.equal(harness.state.all().some(item => item.id === `${kind.prefix}_legacy`), true);
        assert.equal(JSON.parse(harness.store.get(kind.key)).version, 2);
    });

    test(`${kind.label}: a save reaches the backend before the projection or cache changes`, async () => {
        let observedDuringWrite = null;
        let cacheWritesDuringWrite = 0;
        const harness = createHarness(kind, {
            backend: [{ id: kind.id, title: 'old', revision: 'rev-old' }],
        });
        harness.setOnWrite(() => {
            observedDuringWrite = harness.state.all()[0].title;
            cacheWritesDuringWrite = harness.operations.filter(op => op.startsWith('cache:')).length;
        });

        await harness.state.ready;
        harness.operations.length = 0;

        await harness.state.save({ id: kind.id, title: 'new' });

        assert.equal(observedDuringWrite, 'old');
        assert.equal(cacheWritesDuringWrite, 0);

        const saveIndex = harness.operations.indexOf(kind.save);
        const ackIndex = harness.operations.indexOf(kind.save + ':ack');
        const cacheIndex = harness.operations.findIndex(op => op.startsWith('cache:'));
        assert.ok(saveIndex >= 0 && ackIndex > saveIndex && cacheIndex > ackIndex);
        assert.equal(harness.state.get(kind.id).title, 'new');
        assert.equal(JSON.parse(harness.store.get(kind.key))[kind.list][0].title, 'new');
    });

    test(`${kind.label}: a refused save leaves the projection and cache untouched`, async () => {
        const harness = createHarness(kind, {
            backend: [{ id: kind.id, title: 'old', revision: 'rev-old' }],
        });
        await harness.state.ready;
        const cacheBefore = harness.store.get(kind.key);
        harness.operations.length = 0;
        harness.setFail(true);

        await assert.rejects(() => harness.state.save({ id: kind.id, title: 'new' }));

        assert.equal(harness.state.get(kind.id).title, 'old');
        assert.equal(harness.store.get(kind.key), cacheBefore);
        assert.equal(harness.operations.some(op => op.startsWith('cache:')), false);
    });

    test(`${kind.label}: delete reaches the backend before the projection changes`, async () => {
        const harness = createHarness(kind, {
            backend: [{ id: kind.id, title: 'Mira', revision: 'rev-old' }],
        });
        await harness.state.ready;
        harness.operations.length = 0;

        await harness.state.remove(kind.id);

        const deleteIndex = harness.operations.indexOf(kind.remove);
        const ackIndex = harness.operations.indexOf(kind.remove + ':ack');
        const cacheIndex = harness.operations.findIndex(op => op.startsWith('cache:'));
        assert.ok(deleteIndex >= 0 && ackIndex > deleteIndex && cacheIndex > ackIndex);
        assert.equal(harness.state.get(kind.id), null);
    });
}

test('a refused world save cannot mutate the repository projection', async () => {
    const existing = { id: 'wld_existing', title: 'old', revision: 'rv-old' };
    const window = {
        WorldRepository: {
            all() { return [existing]; },
            async save() { throw new Error('disk full'); },
        },
    };
    window.window = window;
    class StorageManager {
        async initialize() {}
        generateId() { return 'wld_generated'; }
    }
    const sandbox = {
        window,
        StorageManager,
        console: { log() {}, warn() {}, error() {} },
        localStorage: { getItem() { return null; }, setItem() {}, removeItem() {} },
        Date,
        Math,
    };
    vm.createContext(sandbox);
    vm.runInContext(
        fs.readFileSync(path.join(JS, 'storage-scenarios.js'), 'utf8'),
        sandbox,
        { filename: 'storage-scenarios.js' }
    );

    const manager = new StorageManager();
    await assert.rejects(() => manager.storeWorld({ id: 'wld_existing', title: 'new' }));

    assert.equal(existing.title, 'old');
    assert.equal(existing.revision, 'rv-old');
});

test('no renderer module touches world or scenario localStorage outside the repositories', () => {
    const repositoryFiles = new Set(['world-repository.js', 'scenario-repository.js']);
    const offenders = [];
    for (const name of fs.readdirSync(JS).filter(entry => entry.endsWith('.js'))) {
        if (repositoryFiles.has(name)) continue;
        const source = fs.readFileSync(path.join(JS, name), 'utf8');
        if (/localStorage\.(?:get|set|remove)Item\(\s*['"]howl-(?:worlds|chat-scenarios)['"]/.test(source)) {
            offenders.push(name);
        }
    }
    assert.deepEqual(offenders, []);
});
