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
    const CACHE_VERSION = 2;

    let worlds = [];
    let source = 'defaults';
    let loadPromise = null;

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
        const legacy = _readLegacy();
        let legacyFailed = false;
        if (legacy && typeof app?.ImportWorlds === 'function') {
            try {
                const imported = await app.ImportWorlds(legacy);
                _emit('howl:worlds-imported', { count: imported });
            } catch (e) {
                legacyFailed = true;
                console.warn('[worldRepository] legacy world import failed:', e);
            }
        }

        let backendList = null;
        if (typeof app?.GetWorlds === 'function') {
            try {
                const list = await app.GetWorlds();
                if (Array.isArray(list)) backendList = list;
            } catch (e) {
                console.warn('[worldRepository] backend load failed:', e);
            }
        }

        if (backendList) {
            // A failed import must not erase data the user can still see: keep the
            // raw array as a documented migration fallback and do not rewrite it.
            if (legacyFailed && legacy && backendList.length === 0) {
                worlds = _normalize(legacy);
                source = 'legacy';
                _emit('howl:worlds-loaded', { source });
                return worlds.slice();
            }
            worlds = _normalize(backendList);
            source = 'backend';
            _writeCache();
            _emit('howl:worlds-loaded', { source });
            return worlds.slice();
        }

        const cached = _readCache();
        worlds = _normalize(cached);
        source = cached ? 'cache' : 'defaults';
        _emit('howl:worlds-loaded', { source });
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

        const saved = await app.SaveWorld(merged);
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
        await app.DeleteWorld(id);
        worlds = worlds.filter(w => w.id !== id);
        _writeCache();
        _emit('howl:world-deleted', { id });
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
        all() { return worlds; },
        get(id) { return _find(id); },
        loadAll,
        save,
        remove,
        mergeRecord,
    };
})();

window.WorldStorage = WorldStorage;
window.WorldRepository = WorldStorage;
