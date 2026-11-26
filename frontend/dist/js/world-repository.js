/**
 * World Repository — the renderer's single seam to world state
 *
 * `internal/worldstore` is canonical. The host bindings GetWorlds / SaveWorld /
 * DeleteWorld / ImportWorlds are the only way this module reads or writes a
 * world. `all()` returns the last backend-acknowledged projection; localStorage
 * holds a versioned cache rewritten only after an acknowledgement.
 *
 * No other renderer module writes `howl-worlds`.
 *
 * @module world-repository
 * @version 1.0.0
 */

'use strict';

const WorldStorage = (() => {
    const CACHE_KEY = 'howl-worlds';
    const LEGACY_BACKUP_KEY = CACHE_KEY + '-cutover-backup';
    const CACHE_VERSION = 2;

    let worlds = [];
    let source = 'defaults';
    let loadPromise = null;
    let cutover = false;
    let cutoverStatus = null;

    function _app() {
        return window.go?.main?.App;
    }

    function _emit(name, detail) {
        if (typeof window.dispatchEvent !== 'function' || typeof CustomEvent !== 'function') return;
        window.dispatchEvent(new CustomEvent(name, { detail }));
    }

    function mergeRecord(existing, patch) {
        if (!patch || typeof patch !== 'object' || Array.isArray(patch)) return patch;
        if (!existing || typeof existing !== 'object' || Array.isArray(existing)) return patch;
        const merged = { ...existing };
        for (const key of Object.keys(patch)) {
            const value = patch[key];
            if (value === undefined) continue;
            if (value && typeof value === 'object' && !Array.isArray(value)
                && existing[key] && typeof existing[key] === 'object' && !Array.isArray(existing[key])) {
                merged[key] = mergeRecord(existing[key], value);
            } else {
                merged[key] = value;
            }
        }
        return merged;
    }

    function _newID() {
        return `wld_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`;
    }

    function _normalizeOne(raw, index) {
        const record = { ...(raw || {}) };
        if (!record.id) record.id = index === undefined ? _newID() : `wld_${Date.now()}_${index}`;
        if (!record.title) record.title = record.name || record.id;
        return record;
    }

    function _normalize(list) {
        if (!Array.isArray(list)) return [];
        return list.map((item, index) => _normalizeOne(item, index));
    }

    function _find(id) {
        if (!id) return null;
        return worlds.find(w => w.id === id) || null;
    }

    function _replace(record) {
        const index = worlds.findIndex(w => w.id === record.id);
        if (index >= 0) worlds[index] = record;
        else worlds.push(record);
    }

    function _readCache() {
        try {
            const raw = localStorage.getItem(CACHE_KEY);
            if (!raw) return null;
            const parsed = JSON.parse(raw);
            if (!parsed || parsed.version !== CACHE_VERSION || !Array.isArray(parsed.worlds)) return null;
            return parsed.worlds;
        } catch (e) {
            console.warn('[worldRepository] cache read failed:', e);
            return null;
        }
    }

    function _readLegacy() {
        try {
            const raw = localStorage.getItem(CACHE_KEY);
            if (!raw) return null;
            const parsed = JSON.parse(raw);
            return Array.isArray(parsed) ? parsed : null;
        } catch {
            return null;
        }
    }

    // The canonical document remains the authority for images; the derived cache
    // keeps only references so a large world image cannot consume the quota.
    function _cacheSafe(world) {
        const copy = { ...world };
        if (copy.image && copy.image.data) {
            copy.image = { id: copy.image.id || null, name: copy.image.name, type: copy.image.type, size: copy.image.size, stored: !!copy.image.stored };
        }
        if (Array.isArray(copy.images)) {
            copy.images = copy.images.map(img => {
                if (!img) return img;
                if (img.data || img.medium || img.original) {
                    return { id: img.id || null, name: img.name, type: img.type, size: img.size, stored: !!img.stored };
                }
                return img;
            });
        }
        return copy;
    }

    function _writeCache() {
        try {
            localStorage.setItem(CACHE_KEY, JSON.stringify({
                version: CACHE_VERSION,
                savedAt: new Date().toISOString(),
                worlds: worlds.map(_cacheSafe),
            }));
        } catch (e) {
            console.warn('[worldRepository] cache write failed:', e);
        }
    }

    async function _load() {
        const app = _app();
        let backendList = null;

        if (typeof app?.GetWorldCutoverStatus === 'function') {
            try {
                cutoverStatus = await app.GetWorldCutoverStatus();
                cutover = Boolean(cutoverStatus?.cutover);
            } catch (e) {
                console.warn('[worldRepository] cutover status failed:', e);
            }
        }

        if (typeof app?.GetWorlds === 'function') {
            try {
                const list = await app.GetWorlds();
                if (Array.isArray(list)) backendList = list;
            } catch (e) {
                console.warn('[worldRepository] backend load failed:', e);
            }
        }

        // Storage state is reported separately from the document list: an
        // uninitialized or read-only library is a per-library state, not an
        // error and not an empty library.
        window.ContentReadiness?.refresh?.();

        if (backendList) {
            // Before the explicit cutover an empty canonical list must not erase
            // the only copy of the user's pre-cutover documents.
            if (!cutover && backendList.length === 0) {
                const legacy = _readLegacy();
                if (legacy) {
                    worlds = _normalize(legacy);
                    source = 'legacy';
                    _emit('howl:worlds-loaded', { source });
                    return worlds.slice();
                }
            }
            worlds = _normalize(backendList);
            source = 'backend';
            _writeCache();
            _emit('howl:worlds-loaded', { source });
            return worlds.slice();
        }

        const cached = _readCache();
        if (cached) {
            worlds = _normalize(cached);
            source = 'cache';
        } else {
            const legacy = _readLegacy();
            worlds = _normalize(legacy);
            source = legacy ? 'legacy' : 'defaults';
        }
        _emit('howl:worlds-loaded', {
            source,
            stale: window.ContentReadiness?.canRead?.('worlds') === false,
        });
        return worlds.slice();
    }

    function loadAll() {
        if (loadPromise) return loadPromise;
        loadPromise = _load().finally(() => { loadPromise = null; });
        return loadPromise;
    }

    async function save(world) {
        if (!world || typeof world !== 'object') throw new Error('World is required');
        const app = _app();
        if (typeof app?.SaveWorld !== 'function') throw new Error('world backend unavailable');

        const existing = world.id ? _find(world.id) : null;
        const merged = mergeRecord(existing, world);
        if (!merged.id) merged.id = _newID();
        if (existing?.revision && !merged.revision) merged.revision = existing.revision;
        if (!merged.createdAt) merged.createdAt = new Date().toISOString();
        merged.updatedAt = new Date().toISOString();

        let saved;
        try {
            saved = await app.SaveWorld(merged);
        } catch (error) {
            // The backend refused the write and said why. Refresh the reported
            // state so the panel explains it instead of looking empty.
            window.ContentReadiness?.refresh?.();
            throw error;
        }
        if (!saved || typeof saved !== 'object') throw new Error('world backend returned no document');
        const record = _normalizeOne(saved);
        _replace(record);
        _writeCache();
        _emit('howl:world-saved', { id: record.id });
        return record;
    }

    async function remove(id) {
        const app = _app();
        if (typeof app?.DeleteWorld !== 'function') throw new Error('world backend unavailable');
        try {
            await app.DeleteWorld(id);
        } catch (error) {
            window.ContentReadiness?.refresh?.();
            throw error;
        }
        worlds = worlds.filter(w => w.id !== id);
        _writeCache();
        _emit('howl:world-deleted', { id });
    }

    async function importLegacy() {
        const app = _app();
        if (typeof app?.ImportWorlds !== 'function') {
            throw new Error('world backend unavailable');
        }
        const legacy = _readLegacy();
        if (!legacy || legacy.length === 0) {
            throw new Error('No pre-cutover worlds were found in this browser.');
        }
        localStorage.setItem(LEGACY_BACKUP_KEY, JSON.stringify(legacy));
        let result;
        try {
            result = await app.ImportWorlds(legacy);
        } catch (error) {
            localStorage.removeItem(LEGACY_BACKUP_KEY);
            throw error;
        }
        await loadAll();
        return result;
    }

    async function rollback() {
        const app = _app();
        if (typeof app?.RollbackWorlds !== 'function') {
            throw new Error('world backend unavailable');
        }
        const result = await app.RollbackWorlds();
        const backup = localStorage.getItem(LEGACY_BACKUP_KEY);
        if (backup !== null) {
            localStorage.setItem(CACHE_KEY, backup);
            localStorage.removeItem(LEGACY_BACKUP_KEY);
        }
        await loadAll();
        return result;
    }

    const seeded = _readCache();
    if (seeded) {
        worlds = _normalize(seeded);
        source = 'cache';
    }

    const ready = loadAll();

    return {
        CACHE_KEY,
        CACHE_VERSION,
        get source() { return source; },
        get ready() { return ready; },
        get cutover() { return cutover; },
        get cutoverStatus() { return cutoverStatus; },
        all() { return worlds; },
        get(id) { return _find(id); },
        loadAll,
        importLegacy,
        rollback,
        save,
        remove,
        mergeRecord,
    };
})();

window.WorldStorage = WorldStorage;
window.WorldRepository = WorldStorage;
