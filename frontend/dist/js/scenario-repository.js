/**
 * Scenario Repository — the renderer's single seam to scenario state
 *
 * `internal/scenariostore` is canonical. The host bindings GetScenarios / SaveScenario /
 * DeleteScenario / ImportScenarios are the only way this module reads or writes a
 * scenario. `all()` returns the last backend-acknowledged projection; localStorage
 * holds a versioned cache rewritten only after an acknowledgement.
 *
 * No other renderer module writes `howl-chat-scenarios`.
 *
 * @module scenario-repository
 * @version 1.0.0
 */

'use strict';

const ScenarioStorage = (() => {
    const CACHE_KEY = 'howl-chat-scenarios';
    const LEGACY_BACKUP_KEY = CACHE_KEY + '-cutover-backup';
    const CACHE_VERSION = 2;

    let scenarios = [];
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
        return `scn_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`;
    }

    function _normalizeOne(raw, index) {
        const record = { ...(raw || {}) };
        if (!record.id) record.id = index === undefined ? _newID() : `scn_${Date.now()}_${index}`;
        if (!record.title) record.title = record.name || record.id;
        return record;
    }

    function _normalize(list) {
        if (!Array.isArray(list)) return [];
        return list.map((item, index) => _normalizeOne(item, index));
    }

    function _find(id) {
        if (!id) return null;
        return scenarios.find(w => w.id === id) || null;
    }

    function _replace(record) {
        const index = scenarios.findIndex(w => w.id === record.id);
        if (index >= 0) scenarios[index] = record;
        else scenarios.push(record);
    }

    function _readCache() {
        try {
            const raw = localStorage.getItem(CACHE_KEY);
            if (!raw) return null;
            const parsed = JSON.parse(raw);
            if (!parsed || parsed.version !== CACHE_VERSION || !Array.isArray(parsed.scenarios)) return null;
            return parsed.scenarios;
        } catch (e) {
            console.warn('[scenarioRepository] cache read failed:', e);
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
    // keeps only references so a large scenario image cannot consume the quota.
    function _cacheSafe(scenario) {
        const copy = { ...scenario };
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
                scenarios: scenarios.map(_cacheSafe),
            }));
        } catch (e) {
            console.warn('[scenarioRepository] cache write failed:', e);
        }
    }

    async function _load() {
        const app = _app();
        let backendList = null;

        if (typeof app?.GetScenarioCutoverStatus === 'function') {
            try {
                cutoverStatus = await app.GetScenarioCutoverStatus();
                cutover = Boolean(cutoverStatus?.cutover);
            } catch (e) {
                console.warn('[scenarioRepository] cutover status failed:', e);
            }
        }

        if (typeof app?.GetScenarios === 'function') {
            try {
                const list = await app.GetScenarios();
                if (Array.isArray(list)) backendList = list;
            } catch (e) {
                console.warn('[scenarioRepository] backend load failed:', e);
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
                    scenarios = _normalize(legacy);
                    source = 'legacy';
                    _emit('howl:scenarios-loaded', { source });
                    return scenarios.slice();
                }
            }
            scenarios = _normalize(backendList);
            source = 'backend';
            _writeCache();
            _emit('howl:scenarios-loaded', { source });
            return scenarios.slice();
        }

        const cached = _readCache();
        if (cached) {
            scenarios = _normalize(cached);
            source = 'cache';
        } else {
            const legacy = _readLegacy();
            scenarios = _normalize(legacy);
            source = legacy ? 'legacy' : 'defaults';
        }
        _emit('howl:scenarios-loaded', {
            source,
            stale: window.ContentReadiness?.canRead?.('scenarios') === false,
        });
        return scenarios.slice();
    }

    function loadAll() {
        if (loadPromise) return loadPromise;
        loadPromise = _load().finally(() => { loadPromise = null; });
        return loadPromise;
    }

    async function save(scenario) {
        if (!scenario || typeof scenario !== 'object') throw new Error('Scenario is required');
        const app = _app();
        if (typeof app?.SaveScenario !== 'function') throw new Error('scenario backend unavailable');

        const existing = scenario.id ? _find(scenario.id) : null;
        const merged = mergeRecord(existing, scenario);
        if (!merged.id) merged.id = _newID();
        if (existing?.revision && !merged.revision) merged.revision = existing.revision;
        if (!merged.createdAt) merged.createdAt = new Date().toISOString();
        merged.updatedAt = new Date().toISOString();

        let saved;
        try {
            saved = await app.SaveScenario(merged);
        } catch (error) {
            // The backend refused the write and said why. Refresh the reported
            // state so the panel explains it instead of looking empty.
            window.ContentReadiness?.refresh?.();
            throw error;
        }
        if (!saved || typeof saved !== 'object') throw new Error('scenario backend returned no document');
        const record = _normalizeOne(saved);
        _replace(record);
        _writeCache();
        _emit('howl:scenario-saved', { id: record.id });
        return record;
    }

    async function remove(id) {
        const app = _app();
        if (typeof app?.DeleteScenario !== 'function') throw new Error('scenario backend unavailable');
        try {
            await app.DeleteScenario(id);
        } catch (error) {
            window.ContentReadiness?.refresh?.();
            throw error;
        }
        scenarios = scenarios.filter(w => w.id !== id);
        _writeCache();
        _emit('howl:scenario-deleted', { id });
    }

    async function importLegacy() {
        const app = _app();
        if (typeof app?.ImportScenarios !== 'function') {
            throw new Error('scenario backend unavailable');
        }
        const legacy = _readLegacy();
        if (!legacy || legacy.length === 0) {
            throw new Error('No pre-cutover scenarios were found in this browser.');
        }
        localStorage.setItem(LEGACY_BACKUP_KEY, JSON.stringify(legacy));
        let result;
        try {
            result = await app.ImportScenarios(legacy);
        } catch (error) {
            localStorage.removeItem(LEGACY_BACKUP_KEY);
            throw error;
        }
        await loadAll();
        return result;
    }

    async function rollback() {
        const app = _app();
        if (typeof app?.RollbackScenarios !== 'function') {
            throw new Error('scenario backend unavailable');
        }
        const result = await app.RollbackScenarios();
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
        scenarios = _normalize(seeded);
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
        all() { return scenarios; },
        get(id) { return _find(id); },
        loadAll,
        importLegacy,
        rollback,
        save,
        remove,
        mergeRecord,
    };
})();

window.ScenarioStorage = ScenarioStorage;
window.ScenarioRepository = ScenarioStorage;
