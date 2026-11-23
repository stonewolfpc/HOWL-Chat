/**
 * Image generation repository tests.
 *
 * They execute the real frontend module in a VM with a fake Wails bridge and
 * localStorage, proving the §21.5 item-4 contract:
 *   - generation settings come from the backend image store, not IG_STATE_KEY,
 *   - a pre-cutover IG_STATE_KEY is imported once and reduced to UI state,
 *   - cache/projection/revision change only after an acknowledged save,
 *   - a refused save leaves the revision intact and the retry succeeds,
 *   - a cleared cache reloads the settings from the backend.
 */

const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');

const JS = path.resolve(__dirname, '../../frontend/dist/js');
const SOURCE = fs.readFileSync(path.join(JS, 'image-generation-repository.js'), 'utf8');
const UI_KEY = 'howl-image-gen-state';
const CACHE_KEY = 'howl-image-generation-cache';

function createHarness(options = {}) {
    const store = new Map();
    for (const [key, value] of Object.entries(options.storage || {})) {
        store.set(key, typeof value === 'string' ? value : JSON.stringify(value));
    }

    const operations = [];
    let backend = options.backend ? JSON.parse(JSON.stringify(options.backend)) : { revision: '' };
    let failWrites = options.failWrites === true;
    let onWrite = options.onWrite || null;
    let savedBases = [];

    const app = {
        async GetImageGenerationSettings() {
            operations.push('GetImageGenerationSettings');
            return { ...JSON.parse(JSON.stringify(backend)) };
        },
        async SaveImageGenerationSettings(settings, base) {
            operations.push('SaveImageGenerationSettings');
            savedBases.push(base);
            if (onWrite) onWrite(settings);
            if (failWrites) throw new Error('disk full');
            backend = { ...settings, revision: 'rv1-' + savedBases.length };
            operations.push('SaveImageGenerationSettings:ack');
            return JSON.parse(JSON.stringify(backend));
        },
        async ImportImageGenerationSettings(settings) {
            operations.push('ImportImageGenerationSettings');
            if (backend.revision || Object.keys(backend).some(key => key !== 'revision')) {
                return JSON.parse(JSON.stringify(backend));
            }
            backend = { ...settings, revision: 'rv1-imported' };
            return JSON.parse(JSON.stringify(backend));
        },
    };

    const localStorage = {
        getItem(key) { return store.has(key) ? store.get(key) : null; },
        setItem(key, value) { operations.push('cache:' + key); store.set(key, String(value)); },
        removeItem(key) { store.delete(key); },
    };

    const window = {
        go: options.noBackend ? {} : { main: { App: app } },
        addEventListener() {},
        dispatchEvent() {},
    };
    window.window = window;

    const sandbox = {
        window,
        localStorage,
        console: { log() {}, warn() {}, error() {} },
        CustomEvent: class CustomEvent {
            constructor(type, init) { this.type = type; this.detail = init && init.detail; }
        },
        Date,
    };
    vm.createContext(sandbox);
    vm.runInContext(SOURCE, sandbox, { filename: 'image-generation-repository.js' });

    return {
        state: window.ImageGenerationRepository,
        store,
        operations,
        backend: () => JSON.parse(JSON.stringify(backend)),
        savedBases,
        setFail(value) { failWrites = value; },
        setOnWrite(fn) { onWrite = fn; },
    };
}

test('backend generation settings beat a conflicting derived cache', async () => {
    const harness = createHarness({
        backend: { steps: '30', width: '1024', revision: 'rv1-backend' },
        storage: {
            [CACHE_KEY]: { version: 2, settings: { steps: '4', width: '512' }, revision: 'rv1-cache' },
        },
    });

    await harness.state.ready;

    assert.equal(harness.state.source, 'backend');
    assert.equal(harness.state.get('steps'), '30');
    assert.equal(harness.state.revision, 'rv1-backend');
    assert.equal(JSON.parse(harness.store.get(CACHE_KEY)).settings.steps, '30');
});

test('a pre-cutover IG_STATE_KEY is imported once and reduced to UI state', async () => {
    const harness = createHarness({
        backend: { revision: '' },
        storage: {
            [UI_KEY]: {
                _version: 2,
                steps: '30',
                prompt: 'legacy prompt',
                currentMode: 'inpaint',
                seedLocked: true,
                aspectLocked: false,
            },
        },
    });

    await harness.state.ready;

    assert.ok(harness.operations.indexOf('ImportImageGenerationSettings') < harness.operations.indexOf('GetImageGenerationSettings'));
    assert.equal(harness.state.get('steps'), '30');
    const ui = JSON.parse(harness.store.get(UI_KEY));
    assert.equal(ui.version, 3);
    assert.deepEqual(ui.ui, { currentMode: 'inpaint', seedLocked: true, aspectLocked: false });
    assert.equal(Object.prototype.hasOwnProperty.call(ui, 'steps'), false);
    assert.equal(Object.prototype.hasOwnProperty.call(ui.ui, 'prompt'), false);
    assert.equal(JSON.parse(harness.store.get(CACHE_KEY)).settings.steps, '30');
});

test('a save changes the projection and cache only after acknowledgement', async () => {
    let observedDuringWrite = null;
    let cacheWritesDuringWrite = 0;
    const harness = createHarness({
        backend: { steps: '20', width: '768', revision: 'rv1-old' },
    });
    harness.setOnWrite(() => {
        observedDuringWrite = harness.state.get('steps');
        cacheWritesDuringWrite = harness.operations.filter(op => op === 'cache:' + CACHE_KEY).length;
    });

    await harness.state.ready;
    harness.operations.length = 0;

    await harness.state.save({ steps: '12' });

    assert.equal(observedDuringWrite, '20');
    assert.equal(cacheWritesDuringWrite, 0);
    const saveIndex = harness.operations.indexOf('SaveImageGenerationSettings');
    const ackIndex = harness.operations.indexOf('SaveImageGenerationSettings:ack');
    const cacheIndex = harness.operations.indexOf('cache:' + CACHE_KEY);
    assert.ok(saveIndex >= 0 && ackIndex > saveIndex && cacheIndex > ackIndex);
    assert.equal(harness.state.get('steps'), '12');
    assert.equal(harness.state.get('width'), '768');
    assert.equal(JSON.parse(harness.store.get(CACHE_KEY)).settings.steps, '12');
});

test('a refused save keeps the revision and the retry is not wedged', async () => {
    const harness = createHarness({
        backend: { steps: '20', revision: 'rv1-old' },
    });
    await harness.state.ready;
    const cacheBefore = harness.store.get(CACHE_KEY);
    harness.operations.length = 0;
    harness.setFail(true);

    await assert.rejects(() => harness.state.save({ steps: '9' }));
    assert.equal(harness.state.get('steps'), '20');
    assert.equal(harness.state.revision, 'rv1-old');
    assert.equal(harness.store.get(CACHE_KEY), cacheBefore);

    harness.setFail(false);
    await harness.state.save({ steps: '9' });
    assert.equal(harness.savedBases[0], 'rv1-old');
    assert.equal(harness.savedBases[1], 'rv1-old');
    assert.equal(harness.state.get('steps'), '9');
    assert.equal(harness.state.revision, 'rv1-2');
});

test('generation settings survive a cleared cache and come from the backend', async () => {
    const first = createHarness({
        backend: { steps: '20', prompt: 'first', revision: 'rv1-old' },
    });
    await first.state.ready;
    await first.state.save({ steps: '17', prompt: 'survives' });
    const backend = first.backend();

    const second = createHarness({ backend, storage: {} });
    await second.state.ready;

    assert.equal(second.state.source, 'backend');
    assert.equal(second.state.get('steps'), '17');
    assert.equal(second.state.get('prompt'), 'survives');
    assert.equal(second.state.revision, backend.revision);
});

test('a failed legacy import does not overwrite IG_STATE_KEY', async () => {
    const harness = createHarness({
        backend: { revision: '' },
        storage: {
            [UI_KEY]: { _version: 2, steps: '30', currentMode: 'txt2img' },
        },
        noBackend: true,
    });

    await harness.state.ready;

    assert.equal(harness.state.source, 'legacy');
    const raw = JSON.parse(harness.store.get(UI_KEY));
    assert.equal(raw.steps, '30');
});

test('no other renderer module touches generation settings through localStorage', () => {
    const files = fs.readdirSync(JS).filter(name => name.endsWith('.js') && name !== 'image-generation-repository.js');
    const offenders = files.filter(name => {
        const source = fs.readFileSync(path.join(JS, name), 'utf8');
        return /localStorage\.(?:get|set|remove)Item\(\s*['"](?:howl-image-gen-state|howl-image-generation-cache)['"]/.test(source);
    });
    assert.deepEqual(offenders, []);
});
