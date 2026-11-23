/**
 * Settings authority tests.
 *
 * These execute the real frontend/dist/js/settings-state.js in a VM with a
 * fake Wails bridge and a fake localStorage. They prove the contract, not the
 * implementation shape:
 *   - backend values beat a stale or conflicting cache,
 *   - an empty/legacy cache never overrides the backend,
 *   - the backend receives a mutation before in-memory state or cache changes,
 *   - a failed backend write reverts and leaves the cache untouched,
 *   - a versioned cache is only used when the backend is unreachable.
 */

const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');

const SOURCE = fs.readFileSync(
    path.resolve(__dirname, '../../frontend/dist/js/settings-state.js'),
    'utf8'
);
const CACHE_KEY = 'howl-chat-settings';

function createHarness(options = {}) {
    const store = new Map();
    if (options.cache !== undefined) {
        store.set(
            CACHE_KEY,
            typeof options.cache === 'string' ? options.cache : JSON.stringify(options.cache)
        );
    }

    const operations = [];
    let backend = options.backend === undefined ? {} : options.backend;
    let failWrites = options.failWrites === true;
    let onSet = options.onSet || null;

    const app = {
        async GetSettings() {
            operations.push('GetSettings');
            return JSON.parse(JSON.stringify(backend));
        },
        async SetSettings(doc) {
            operations.push('SetSettings');
            if (onSet) onSet(doc);
            if (failWrites) throw new Error('disk full');
            backend = JSON.parse(JSON.stringify(doc));
            operations.push('SetSettings:ack');
        },
    };

    const localStorage = {
        getItem(key) {
            return store.has(key) ? store.get(key) : null;
        },
        setItem(key, value) {
            operations.push('cache:' + key);
            store.set(key, String(value));
        },
        removeItem(key) {
            store.delete(key);
        },
        key(index) {
            return Array.from(store.keys())[index] ?? null;
        },
        get length() {
            return store.size;
        },
    };

    const events = [];
    const window = {
        go: options.noBackend ? {} : { main: { App: app } },
        addEventListener() {},
        dispatchEvent(event) {
            events.push(event.type);
        },
        showNotification() {},
    };
    window.window = window;

    const sandbox = {
        window,
        localStorage,
        console: { log() {}, warn() {}, error() {} },
        CustomEvent: class CustomEvent {
            constructor(type, init) {
                this.type = type;
                this.detail = init && init.detail;
            }
        },
        setTimeout,
        clearTimeout,
        document: { getElementById() { return null; } },
    };
    vm.createContext(sandbox);
    vm.runInContext(SOURCE, sandbox, { filename: 'settings-state.js' });

    return {
        state: window.settingsState,
        store,
        operations,
        events,
        setFail(value) { failWrites = value; },
    };
}

test('a conflicting versioned cache cannot override backend values', async () => {
    const harness = createHarness({
        backend: {
            general: { theme: 'dark', uiScale: 1.25 },
            audio: { ttsEngine: 'kokoro' },
        },
        cache: {
            version: 2,
            savedAt: '2026-01-01T00:00:00.000Z',
            settings: {
                general: { theme: 'light', uiScale: 2 },
                audio: { ttsEngine: 'piper' },
            },
        },
    });

    await harness.state.ready;

    assert.equal(harness.state.source, 'backend');
    assert.equal(harness.state.get('general.theme'), 'dark');
    assert.equal(harness.state.get('general.uiScale'), 1.25);
    assert.equal(harness.state.get('audio.ttsEngine'), 'kokoro');

    // The cache is rebuilt from the authority, not left holding the conflict.
    const cached = JSON.parse(harness.store.get(CACHE_KEY));
    assert.equal(cached.version, 2);
    assert.equal(cached.settings.general.uiScale, 1.25);
    assert.equal(cached.settings.audio.ttsEngine, 'kokoro');
});

test('an empty backend is authoritative over a populated cache', async () => {
    const harness = createHarness({
        backend: {},
        cache: {
            version: 2,
            savedAt: '2026-01-01T00:00:00.000Z',
            settings: { general: { theme: 'light', uiScale: 2 } },
        },
    });

    await harness.state.ready;

    assert.equal(harness.state.source, 'backend');
    assert.equal(harness.state.get('general.theme'), 'dark');
    assert.equal(harness.state.get('general.uiScale'), 1);
});

test('a legacy unversioned cache is ignored', async () => {
    const harness = createHarness({
        noBackend: true,
        cache: JSON.stringify({ general: { theme: 'light', uiScale: 3 } }),
    });

    await harness.state.ready;

    assert.equal(harness.state.source, 'defaults');
    assert.equal(harness.state.get('general.theme'), 'dark');
    assert.equal(harness.state.get('general.uiScale'), 1);
});

test('a versioned cache renders only when the backend is unreachable', async () => {
    const harness = createHarness({
        noBackend: true,
        cache: {
            version: 2,
            savedAt: '2026-01-01T00:00:00.000Z',
            settings: { general: { theme: 'light', uiScale: 1.5 } },
        },
    });

    await harness.state.ready;

    assert.equal(harness.state.source, 'cache');
    assert.equal(harness.state.get('general.theme'), 'light');
    assert.equal(harness.state.get('general.uiScale'), 1.5);
});

test('a mutation reaches the backend before state or cache changes', async () => {
    let valueSeenByBackend = null;
    let cacheWritesDuringBackend = 0;
    const harness = createHarness({
        backend: { general: { uiScale: 1 } },
        onSet() {
            valueSeenByBackend = harness.state.get('general.uiScale');
            cacheWritesDuringBackend = harness.operations.filter(op => op.startsWith('cache:')).length;
        },
    });

    await harness.state.ready;
    harness.operations.length = 0;

    const ok = await harness.state.set('general.uiScale', 1.75);

    assert.equal(ok, true);
    assert.equal(valueSeenByBackend, 1, 'the old value must still be observable during the backend write');
    assert.equal(cacheWritesDuringBackend, 0, 'no cache write may happen before the backend acknowledges');

    const setIndex = harness.operations.indexOf('SetSettings');
    const ackIndex = harness.operations.indexOf('SetSettings:ack');
    const cacheIndex = harness.operations.findIndex(op => op.startsWith('cache:'));
    assert.ok(setIndex >= 0, 'SetSettings was not called');
    assert.ok(ackIndex > setIndex, 'the cache/state update must wait for the acknowledgement');
    assert.ok(cacheIndex > ackIndex, 'the cache must be written only after the acknowledgement');

    assert.equal(harness.state.get('general.uiScale'), 1.75);
    const cached = JSON.parse(harness.store.get(CACHE_KEY));
    assert.equal(cached.version, 2);
    assert.equal(cached.settings.general.uiScale, 1.75);
});

test('the form-binding save path is backend-first too', async () => {
    const harness = createHarness({ backend: { general: { uiScale: 1 } } });
    await harness.state.ready;
    harness.operations.length = 0;

    // settings-controls.js mutates current as a draft, then calls saveSettings.
    harness.state.current.general.uiScale = 1.9;
    const ok = await harness.state.saveSettings();

    assert.equal(ok, true);
    const setIndex = harness.operations.indexOf('SetSettings');
    const cacheIndex = harness.operations.findIndex(op => op.startsWith('cache:'));
    assert.ok(setIndex >= 0 && cacheIndex > setIndex);
    assert.equal(harness.state.get('general.uiScale'), 1.9);
});

test('a rejected write reverts state and leaves the cache untouched', async () => {
    const harness = createHarness({ backend: { general: { uiScale: 1 } } });
    await harness.state.ready;
    const cacheBefore = harness.store.get(CACHE_KEY);
    harness.operations.length = 0;
    harness.setFail(true);

    const ok = await harness.state.set('general.uiScale', 2);

    assert.equal(ok, false);
    assert.equal(harness.state.get('general.uiScale'), 1);
    assert.equal(harness.store.get(CACHE_KEY), cacheBefore);
    assert.equal(harness.operations.some(op => op.startsWith('cache:')), false);
    assert.ok(harness.events.includes('howl:settings-write-failed'));
});
