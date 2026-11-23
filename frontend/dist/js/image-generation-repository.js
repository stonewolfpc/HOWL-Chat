/**
 * Image Generation Repository — the renderer's single seam to generation state
 *
 * Ownership:
 *   - `internal/image/settings` owns the persisted generation document inside
 *     the existing image-settings record.
 *   - GetImageGenerationSettings / SaveImageGenerationSettings /
 *     ImportImageGenerationSettings are the only backend paths.
 *   - `settings` is the last backend-acknowledged projection and `revision` is
 *     the base revision the next save must cite.
 *   - `IG_STATE_KEY` holds only transient UI state (workflow mode and locks).
 *     Generation settings are cached under a separate derived key, written only
 *     after an acknowledged backend read or write.
 *
 * No other renderer module reads or writes generation settings through
 * localStorage.
 *
 * @module image-generation-repository
 * @version 1.0.0
 */

'use strict';

const ImageGenerationStorage = (() => {
    const UI_KEY = 'howl-image-gen-state';
    const UI_VERSION = 3;
    const CACHE_KEY = 'howl-image-generation-cache';
    const CACHE_VERSION = 2;
    const UI_FIELDS = ['currentMode', 'seedLocked', 'aspectLocked'];

    let settings = {};
    let revision = '';
    let ui = {};
    let source = 'defaults';
    let loaded = false;
    let loadPromise = null;
    let writeChain = Promise.resolve();

    function _app() {
        return window.go?.main?.App;
    }

    function _clone(value) {
        return JSON.parse(JSON.stringify(value));
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

    function _withoutRevision(document) {
        const copy = { ...(document || {}) };
        delete copy.revision;
        return copy;
    }

    function _readUiRaw() {
        try {
            const raw = localStorage.getItem(UI_KEY);
            if (!raw) return null;
            const parsed = JSON.parse(raw);
            return parsed && typeof parsed === 'object' && !Array.isArray(parsed) ? parsed : null;
        } catch {
            return null;
        }
    }

    function _splitLegacy(raw) {
        if (!raw) return { ui: {}, settings: null };
        if (raw.version === UI_VERSION && raw.ui && typeof raw.ui === 'object') {
            return { ui: _clone(raw.ui), settings: null };
        }
        const uiPart = {};
        for (const key of UI_FIELDS) {
            if (raw[key] !== undefined) uiPart[key] = raw[key];
        }
        let generation = null;
        if (raw.generation && typeof raw.generation === 'object' && !Array.isArray(raw.generation)) {
            generation = _clone(raw.generation);
        } else {
            generation = {};
            for (const key of Object.keys(raw)) {
                if (key === '_version' || key === 'version' || key === 'ui' || UI_FIELDS.includes(key)) continue;
                generation[key] = raw[key];
            }
            if (Object.keys(generation).length === 0) generation = null;
        }
        return { ui: uiPart, settings: generation };
    }

    function _writeUi(next) {
        ui = { ...(next || {}) };
        try {
            localStorage.setItem(UI_KEY, JSON.stringify({ version: UI_VERSION, ui }));
        } catch (e) {
            console.warn('[imageGenerationRepository] UI state write failed:', e);
        }
    }

    function _readCache() {
        try {
            const raw = localStorage.getItem(CACHE_KEY);
            if (!raw) return null;
            const parsed = JSON.parse(raw);
            if (!parsed || parsed.version !== CACHE_VERSION || !parsed.settings || typeof parsed.settings !== 'object') {
                return null;
            }
            return { settings: parsed.settings, revision: parsed.revision || '' };
        } catch (e) {
            console.warn('[imageGenerationRepository] cache read failed:', e);
            return null;
        }
    }

    function _writeCache() {
        try {
            localStorage.setItem(CACHE_KEY, JSON.stringify({
                version: CACHE_VERSION,
                savedAt: new Date().toISOString(),
                settings,
                revision,
            }));
        } catch (e) {
            console.warn('[imageGenerationRepository] cache write failed:', e);
        }
    }

    async function _load() {
        const app = _app();
        const rawUi = _readUiRaw();
        const legacy = _splitLegacy(rawUi);
        let legacyFailed = false;
        let imported = false;

        if (legacy.settings && typeof app?.ImportImageGenerationSettings === 'function') {
            try {
                await app.ImportImageGenerationSettings(legacy.settings);
                imported = true;
                _emit('howl:image-generation-imported', { count: Object.keys(legacy.settings).length });
            } catch (e) {
                legacyFailed = true;
                console.warn('[imageGenerationRepository] legacy generation import failed:', e);
            }
        }

        let backend = null;
        if (typeof app?.GetImageGenerationSettings === 'function') {
            try {
                const document = await app.GetImageGenerationSettings();
                if (document && typeof document === 'object' && !Array.isArray(document)) backend = document;
            } catch (e) {
                console.warn('[imageGenerationRepository] backend load failed:', e);
            }
        }

        if (backend) {
            const backendSettings = _withoutRevision(backend);
            // A failed import must not destroy the only copy of legacy settings.
            if (legacyFailed && legacy.settings && Object.keys(backendSettings).length === 0) {
                settings = _clone(legacy.settings);
                revision = '';
                source = 'legacy';
                _emit('howl:image-generation-loaded', { source });
                return settings;
            }
            settings = backendSettings;
            revision = typeof backend.revision === 'string' ? backend.revision : '';
            source = 'backend';
            if (!legacyFailed) {
                _writeCache();
                if (imported || legacy.settings) _writeUi(legacy.ui);
            }
            _emit('howl:image-generation-loaded', { source });
            return _clone(settings);
        }

        const cached = _readCache();
        if (cached) {
            settings = _clone(cached.settings);
            revision = cached.revision;
            source = 'cache';
        } else if (legacy.settings) {
            settings = _clone(legacy.settings);
            revision = '';
            source = 'legacy';
        } else {
            settings = {};
            revision = '';
            source = 'defaults';
        }
        ui = { ...legacy.ui };
        _emit('howl:image-generation-loaded', { source });
        return _clone(settings);
    }

    function loadAll() {
        if (loaded) return Promise.resolve(_clone(settings));
        if (loadPromise) return loadPromise;
        loadPromise = _load().then(result => {
            loaded = true;
            return result;
        }).finally(() => { loadPromise = null; });
        return loadPromise;
    }

    async function _save(patch) {
        await loadAll();
        const app = _app();
        if (typeof app?.SaveImageGenerationSettings !== 'function') {
            throw new Error('image generation backend unavailable');
        }
        if (!patch || typeof patch !== 'object' || Array.isArray(patch)) {
            throw new Error('Image generation settings are required');
        }
        const merged = mergeRecord(settings, _withoutRevision(patch));
        const saved = await app.SaveImageGenerationSettings(merged, revision);
        if (!saved || typeof saved !== 'object') {
            throw new Error('image generation backend returned no document');
        }
        settings = _withoutRevision(saved);
        revision = typeof saved.revision === 'string' ? saved.revision : '';
        _writeCache();
        _emit('howl:image-generation-saved', { revision });
        return _clone(settings);
    }

    // Writes are serialized so a second change cannot cite the revision from
    // before the first acknowledgement. A failed save leaves the projection,
    // revision and cache untouched, so a retry is not wedged into a conflict.
    function save(patch) {
        const run = writeChain.then(() => _save(patch));
        writeChain = run.catch(() => {});
        return run;
    }

    function saveUi(next) {
        // An unmigrated legacy blob is the only copy of generation settings;
        // never replace it with the UI-only envelope until import succeeds.
        if (source === 'legacy') return;
        _writeUi(next);
    }

    function seedUi() {
        const raw = _readUiRaw();
        const split = _splitLegacy(raw);
        ui = { ...split.ui };
        return { ...ui };
    }

    seedUi();
    const ready = loadAll();

    return {
        UI_KEY,
        UI_VERSION,
        CACHE_KEY,
        CACHE_VERSION,
        get source() { return source; },
        get ready() { return ready; },
        get revision() { return revision; },
        all() { return _clone(settings); },
        get(key) { return settings[key]; },
        ui() { return { ...ui }; },
        loadAll,
        save,
        saveUi,
        mergeRecord,
    };
})();

window.ImageGenerationStorage = ImageGenerationStorage;
window.ImageGenerationRepository = ImageGenerationStorage;
