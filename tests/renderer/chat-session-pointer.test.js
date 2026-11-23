/**
 * Chat session pointer tests.
 *
 * The pointer is allowed to hold the active chat id and selection ids/tracking
 * arrays only. Full records must come from the domain repositories.
 */

const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');

const JS = path.resolve(__dirname, '../../frontend/dist/js');
const SOURCE = fs.readFileSync(path.join(JS, 'chat-session-pointer.js'), 'utf8');
const KEY = 'howl-chat-current-session';

function createHarness(options = {}) {
    const store = new Map();
    for (const [key, value] of Object.entries(options.storage || {})) {
        store.set(key, typeof value === 'string' ? value : JSON.stringify(value));
    }

    const worlds = options.worlds || [];
    const scenarios = options.scenarios || [];
    const characters = options.characters || [];
    const lorebooks = options.lorebooks || [];
    const window = {
        localStorage: {
            getItem(key) { return store.has(key) ? store.get(key) : null; },
            setItem(key, value) { store.set(key, String(value)); },
            removeItem(key) { store.delete(key); },
        },
        dataSync: { currentChatId: options.chatId === undefined ? 'chat-1' : options.chatId },
        sidebarState: { currentChatData: { sessionId: 'chat-1', chatName: 'Test Chat' } },
        chatBridge: options.bridgeState ? { getState: () => options.bridgeState } : null,
        WorldRepository: { ready: Promise.resolve(), all: () => worlds },
        ScenarioRepository: { ready: Promise.resolve(), all: () => scenarios },
        CharacterStorage: { ready: Promise.resolve(), all: () => characters },
        LorebookRepo: {
            state: { ready: true, entries: lorebooks },
            all: () => lorebooks,
            readCache: () => lorebooks,
        },
    };
    window.window = window;

    const sandbox = {
        window,
        localStorage: window.localStorage,
        console: { log() {}, warn() {}, error() {} },
        Promise,
    };
    vm.createContext(sandbox);
    vm.runInContext(SOURCE, sandbox, { filename: 'chat-session-pointer.js' });
    return { pointer: window.ChatSessionPointer, store };
}

function fullChatData() {
    return {
        sessionId: 'chat-1',
        chatName: 'Chat Name',
        world: { id: 'w1', name: 'World Name', title: 'World Title' },
        scenario: { id: 's1', name: 'Scenario Name', title: 'Scenario Title' },
        characters: [{
            id: 'c1',
            name: 'Character Name',
            status: 'idle',
            active: true,
            stepOutReason: 'busy',
            stepOutMode: 'manual',
        }],
        lorebooks: [{ id: 'l1', name: 'Lore Name', title: 'Lore Title' }],
        activeCharacterIds: ['c1'],
        removedCharacterIds: [],
        addedCharacterIds: ['c1'],
        activeLorebookIds: ['l1'],
        removedLorebookIds: [],
        addedLorebookIds: ['l1'],
    };
}

test('write stores only the chat id and selection ids, never catalog records', () => {
    const harness = createHarness();
    const written = harness.pointer.write(fullChatData());

    const raw = JSON.parse(harness.store.get(KEY));
    assert.equal(raw.version, 2);
    assert.equal(raw.chatId, 'chat-1');
    assert.equal(raw.selection.worldId, 'w1');
    assert.equal(raw.selection.scenarioId, 's1');
    assert.deepEqual(raw.selection.characters, [{
        id: 'c1',
        status: 'idle',
        active: true,
        stepOutReason: 'busy',
        stepOutMode: 'manual',
    }]);
    assert.deepEqual(raw.selection.lorebookIds, ['l1']);
    assert.equal(raw.selection.lorebooks, undefined);
    assert.equal(raw.selection.world, undefined);
    assert.equal(raw.selection.chatName, undefined);
    assert.equal(JSON.stringify(raw).includes('World Title'), false);
    assert.equal(JSON.stringify(raw).includes('Character Name'), false);
    assert.deepEqual(written.selection.activeCharacterIds, ['c1']);
});

test('hydrate resolves the pointer from the domain repositories', async () => {
    const harness = createHarness({
        worlds: [{ id: 'w1', title: 'Canonical World', name: 'Canonical World' }],
        scenarios: [{ id: 's1', title: 'Canonical Scenario', name: 'Canonical Scenario' }],
        characters: [{ id: 'c1', name: 'Canonical Character', status: 'online', active: false }],
        lorebooks: [{ id: 'l1', title: 'Canonical Lore', name: 'Canonical Lore' }],
    });
    harness.pointer.write(fullChatData());
    const pointer = harness.pointer.read();
    const hydrated = await harness.pointer.hydrate(pointer);

    assert.equal(hydrated.world.title, 'Canonical World');
    assert.equal(hydrated.scenario.title, 'Canonical Scenario');
    assert.equal(hydrated.characters[0].name, 'Canonical Character');
    assert.equal(hydrated.characters[0].status, 'idle');
    assert.equal(hydrated.characters[0].stepOutReason, 'busy');
    assert.deepEqual([...hydrated.activeCharacterIds], ['c1']);
    assert.equal(hydrated.lorebooks[0].title, 'Canonical Lore');
    assert.equal(hydrated.chatName, 'Test Chat');
});

test('a pre-pointer projection is migrated once and loses embedded catalog records', async () => {
    const legacy = {
        chatName: 'Legacy Chat',
        world: { id: 'w1', name: 'Legacy World', title: 'Legacy World' },
        scenario: { id: 's1', name: 'Legacy Scenario', title: 'Legacy Scenario' },
        characters: [{ id: 'c1', name: 'Legacy Character', status: 'ready' }],
        lorebooks: [{ id: 'l1', name: 'Legacy Lore', title: 'Legacy Lore' }],
        activeCharacterIds: ['c1'],
        activeLorebookIds: ['l1'],
    };
    const harness = createHarness({ storage: { [KEY]: legacy } });
    const pointer = harness.pointer.read();

    assert.equal(pointer.version, 2);
    assert.equal(pointer.chatId, 'chat-1');
    assert.equal(pointer.selection.worldId, 'w1');
    assert.equal(pointer.selection.scenarioId, 's1');
    assert.equal(pointer.selection.characters[0].id, 'c1');
    assert.deepEqual([...pointer.selection.lorebookIds], ['l1']);
    assert.ok(pointer.legacy);

    const hydrated = await harness.pointer.hydrate(pointer);
    assert.equal(hydrated.world.title, 'Legacy World');
    assert.equal(hydrated.characters[0].name, 'Legacy Character');
    assert.equal(hydrated.lorebooks[0].title, 'Legacy Lore');

    harness.pointer.write(hydrated, pointer.chatId);
    const migrated = JSON.parse(harness.store.get(KEY));
    assert.equal(migrated.version, 2);
    assert.equal(JSON.stringify(migrated).includes('Legacy World'), false);
    assert.equal(JSON.stringify(migrated).includes('Legacy Character'), false);
    assert.equal(migrated.selection.worldId, 'w1');
});

test('the pointer never falls back to a raw lorebook localStorage catalog', () => {
    const harness = createHarness({
        storage: { 'howl-chat-lorebooks': JSON.stringify([{ id: 'l1', title: 'Old Catalog' }]) },
        lorebooks: [],
    });
    harness.pointer.write(fullChatData());
    const pointer = harness.pointer.read();
    return harness.pointer.hydrate(pointer).then(hydrated => {
        assert.equal(hydrated.lorebooks.length, 0);
    });
});

test('clear removes the pointer and no other module owns the key', () => {
    const harness = createHarness();
    harness.pointer.write(fullChatData());
    assert.ok(harness.store.get(KEY));
    harness.pointer.clear();
    assert.equal(harness.store.get(KEY), undefined);

    const files = fs.readdirSync(JS).filter(name => name.endsWith('.js'));
    const offenders = files.filter(name => {
        if (name === 'chat-session-pointer.js' || name === 'app-wipe.js') return false;
        return fs.readFileSync(path.join(JS, name), 'utf8').includes(KEY);
    });
    assert.deepEqual(offenders, []);
});

test('the chat id prefers the canonical bridge state, then dataSync', () => {
    const harness = createHarness({ chatId: 'chat-data-sync', bridgeState: { sessionId: 'chat-bridge' } });
    const written = harness.pointer.write(fullChatData());
    assert.equal(written.chatId, 'chat-bridge');
});
