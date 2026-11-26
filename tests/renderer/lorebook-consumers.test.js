/**
 * Lorebook consumer migration tests.
 *
 * The browser cache may be read only by the repository. Every page that shows
 * or resolves lorebooks must load that repository, and every synchronous
 * consumer must read the repository projection rather than localStorage.
 */

const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');

const ROOT = path.resolve(__dirname, '../..');
const DIST = path.join(ROOT, 'frontend/dist');
const JS = path.join(DIST, 'js');
const KEY = 'howl-chat-lorebooks';

const CONSUMER_PAGES = ['index.html', 'lorebooks.html', 'world.html', 'scenario.html', 'character.html', 'story-tools.html'];

test('no renderer module reads or writes the lorebook cache directly', () => {
    for (const name of fs.readdirSync(JS).filter(file => file.endsWith('.js'))) {
        if (name === 'lorebook-repository.js' || name === 'app-wipe.js') continue;
        const source = fs.readFileSync(path.join(JS, name), 'utf8');
        assert.equal(source.includes(`getItem('${KEY}'`), false, `${name} reads ${KEY} directly`);
        assert.equal(source.includes(`setItem('${KEY}'`), false, `${name} writes ${KEY} directly`);
    }
});

test('every page with lorebook consumers loads the repository', () => {
    for (const name of CONSUMER_PAGES) {
        const html = fs.readFileSync(path.join(DIST, name), 'utf8');
        assert.match(html, /lorebook-repository\.js\?v=3/, `${name} does not load LorebookRepo`);
    }
});

function createRepository(storage, listResult) {
    const store = new Map(Object.entries(storage || {}));
    const operations = [];
    const window = {
        localStorage: {
            getItem(key) { return store.has(key) ? store.get(key) : null; },
            setItem(key, value) { store.set(key, String(value)); },
            removeItem(key) { store.delete(key); },
        },
        go: {
            main: {
                App: {
                    async ListLorebooks() {
                        operations.push('ListLorebooks');
                        return listResult;
                    },
                },
            },
        },
    };
    window.window = window;

    const sandbox = {
        window,
        localStorage: window.localStorage,
        document: { addEventListener() {} },
        console: { log() {}, warn() {}, error() {} },
        Promise,
    };
    vm.createContext(sandbox);
    const source = fs.readFileSync(path.join(JS, 'lorebook-repository.js'), 'utf8');
    vm.runInContext(source, sandbox, { filename: 'lorebook-repository.js' });
    return { repo: window.LorebookRepo, store, operations };
}

test('the cached catalog seeds synchronous readers before hydration', () => {
    const cached = [{ id: 'lb_cached', title: 'Cached' }];
    const { repo } = createRepository({ [KEY]: JSON.stringify(cached) }, { cutover: false, items: [] });
    assert.deepEqual(JSON.parse(JSON.stringify(repo.all())), cached);
});

test('hydration replaces the cache with canonical library entries', async () => {
    const cached = [{ id: 'lb_cached', title: 'Cached' }];
    const canonical = [
        { id: 'lb_library', title: 'Library B' },
        { id: 'lb_alpha', title: 'Library A' },
    ];
    const { repo, store, operations } = createRepository({ [KEY]: JSON.stringify(cached) }, {
        cutover: true,
        items: canonical.map(entry => ({ entry })),
    });

    await repo.hydrate();

    assert.deepEqual(JSON.parse(JSON.stringify(repo.all().map(entry => entry.id))), ['lb_alpha', 'lb_library']);
    assert.deepEqual(JSON.parse(store.get(KEY)).map(entry => entry.id), ['lb_alpha', 'lb_library']);
    assert.deepEqual(operations, ['ListLorebooks']);
});
