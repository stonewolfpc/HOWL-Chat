/**
 * Lorebook Repository
 *
 * One path for lorebook data. Canonical entries are JSON files in registered
 * library folders, owned by the backend store. The browser store is a derived
 * cache so existing synchronous readers keep working. Before the explicit
 * import the browser store is still the source; nothing changes silently.
 *
 * @module lorebook-repository
 * @version 1.1.0
 */

(function () {
    'use strict';

    const CACHE_KEY = 'howl-chat-lorebooks';

    const state = {
        ready: false,
        cutover: false,
        entries: [],
        status: null,
        error: null
    };

    let hydratePromise = null;

    function api() {
        return (window.go && window.go.main) ? window.go.main.App : null;
    }

    function readCache() {
        try {
            const raw = JSON.parse(localStorage.getItem(CACHE_KEY) || '[]');
            return Array.isArray(raw) ? raw : [];
        } catch (error) {
            console.error('Failed to read the lorebook cache:', error);
            return [];
        }
    }

    function writeCache(list) {
        try {
            localStorage.setItem(CACHE_KEY, JSON.stringify(list));
        } catch (error) {
            console.error('Failed to write the lorebook cache:', error);
        }
    }

    function sortEntries(list) {
        return [...list].sort((a, b) =>
            String(a.title || a.name || '').localeCompare(String(b.title || b.name || '')));
    }

    // Seed synchronously so a page that renders before the backend answers sees
    // the last acknowledged cache instead of an empty catalog. Hydration then
    // replaces it with canonical entries once the backend has answered.
    state.entries = readCache();

    async function hydrate() {
        if (hydratePromise) return hydratePromise;
        hydratePromise = (async () => {
            const app = api();
            if (!app || typeof app.ListLorebooks !== 'function') {
                state.ready = true;
                state.entries = readCache();
                return state;
            }
            try {
                const result = await app.ListLorebooks();
                state.status = result || null;
                state.cutover = Boolean(result && result.cutover);
                if (state.cutover) {
                    const items = (result.items || [])
                        .map(item => item && item.entry)
                        .filter(Boolean);
                    state.entries = sortEntries(items);
                    writeCache(state.entries);
                } else {
                    state.entries = readCache();
                }
                state.error = null;
            } catch (error) {
                console.error('Failed to read the lorebook library:', error);
                state.error = error;
                state.entries = readCache();
            }
            state.ready = true;
            return state;
        })();
        try {
            return await hydratePromise;
        } finally {
            hydratePromise = null;
        }
    }

    async function save(entry, isExisting) {
        const app = api();
        if (!state.cutover || !app || typeof app.SaveLorebook !== 'function') {
            const list = readCache();
            const index = list.findIndex(item => item.id === entry.id);
            if (index >= 0) {
                list[index] = entry;
            } else {
                list.push(entry);
            }
            writeCache(list);
            state.entries = list;
            return { status: 'saved' };
        }
        const result = await app.SaveLorebook(entry, Boolean(isExisting));
        if (!result || (result.status !== 'saved' && result.status !== 'replaced')) {
            throw new Error((result && result.message) || 'The lorebook could not be saved to its library.');
        }
        await hydrate();
        return result;
    }

    async function remove(id) {
        const app = api();
        if (!state.cutover || !app || typeof app.DeleteLorebook !== 'function') {
            writeCache(readCache().filter(item => item.id !== id));
            await hydrate();
            return { status: 'deleted' };
        }
        const result = await app.DeleteLorebook(id);
        if (!result || result.status !== 'deleted') {
            throw new Error((result && result.message) || 'The lorebook could not be deleted.');
        }
        await hydrate();
        return result;
    }

    async function importLegacy() {
        const app = api();
        if (!app || typeof app.ImportLorebooks !== 'function') {
            throw new Error('The lorebook library is not available in this window.');
        }
        const result = await app.ImportLorebooks(readCache());
        await hydrate();
        return result;
    }

    async function rollback() {
        const app = api();
        if (!app || typeof app.RollbackLorebookImport !== 'function') {
            throw new Error('The lorebook library is not available in this window.');
        }
        const result = await app.RollbackLorebookImport();
        await hydrate();
        return result;
    }

    // Begin hydration once the document is ready. Pages may await the same
    // promise without creating a second backend read.
    if (typeof document !== 'undefined' && typeof document.addEventListener === 'function') {
        document.addEventListener('DOMContentLoaded', () => {
            hydrate().catch(error => console.error('Failed to refresh the lorebook library:', error));
        });
    }

    window.LorebookRepo = {
        hydrate,
        save,
        remove,
        importLegacy,
        rollback,
        readCache,
        all() { return state.entries; },
        state
    };
})();
