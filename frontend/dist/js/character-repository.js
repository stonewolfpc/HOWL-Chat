/**
 * Character Repository — the renderer's single seam to character state
 *
 * Ownership:
 *   - `internal/characterstore` is canonical. The host bindings
 *     GetCharacters / SaveCharacter / DeleteCharacter / ImportCharacters are
 *     the only way this module reads or writes characters.
 *   - `characters` is the last backend-acknowledged list. `all()` returns that
 *     projection for synchronous display code; it is not a second authority.
 *   - localStorage holds a versioned, derived cache that is rewritten only
 *     after an acknowledged backend read or write. A pre-cutover raw array is
 *     imported once through ImportCharacters, then replaced by the envelope.
 *   - Legacy IndexedDB character images are read only as a display fallback;
 *     the next save writes them into the canonical document.
 *
 * There is no other writer of `howl-chat-characters`.
 *
 * @module character-repository
 * @version 3.0.0
 */

'use strict';

const CharacterStorage = (() => {
    const CACHE_KEY = 'howl-chat-characters';
    const LEGACY_BACKUP_KEY = CACHE_KEY + '-cutover-backup';
    const CACHE_VERSION = 2;
    const LEGACY_DB_NAME = 'howl-character-images';
    const LEGACY_DB_VERSION = 1;
    const LEGACY_STORE = 'images';

    let characters = [];
    let source = 'defaults';
    let loadPromise = null;
    let cutover = false;
    let cutoverStatus = null;

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

    /**
     * Merge a character patch into an existing record. Objects merge
     * recursively so one edited voice field cannot erase its siblings; arrays
     * and scalars replace, which keeps explicit clears authoritative.
     */
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

    function _normalizeOne(raw, index) {
        const record = { ...(raw || {}) };
        if (!record.id) {
            record.id = `char_${Date.now()}_${index}_${Math.floor(Math.random() * 1000)}`;
        }
        if (!record.name) record.name = record.id;
        return record;
    }

    function _normalize(list) {
        if (!Array.isArray(list)) return [];
        return list.map((item, index) => _normalizeOne(item, index));
    }

    function _find(id) {
        if (!id) return null;
        return characters.find(c => c.id === id) || null;
    }

    function _replace(record) {
        const index = characters.findIndex(c => c.id === record.id);
        if (index >= 0) characters[index] = record;
        else characters.push(record);
    }

    // -----------------------------------------------------------------------
    // Derived cache
    // -----------------------------------------------------------------------

    function _readCache() {
        try {
            const raw = localStorage.getItem(CACHE_KEY);
            if (!raw) return null;
            const parsed = JSON.parse(raw);
            if (!parsed || parsed.version !== CACHE_VERSION || !Array.isArray(parsed.characters)) {
                return null;
            }
            return parsed.characters;
        } catch (e) {
            console.warn('[characterRepository] cache read failed:', e);
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

    // The canonical document may carry base64 image data. The derived cache is
    // for display and must not duplicate blobs into a small quota, so it keeps
    // image metadata only; the backend document remains the image authority.
    function _cacheSafe(character) {
        const copy = { ...character };
        if (Array.isArray(copy.howl?.images)) {
            copy.howl = {
                ...copy.howl,
                images: copy.howl.images.map(img => ({
                    id: img.id,
                    name: img.name,
                    size: img.size,
                    isThumbnail: img.isThumbnail,
                })),
            };
        }
        if (typeof copy.avatar === 'string' && copy.avatar.startsWith('data:')) {
            copy.avatar = '';
        }
        return copy;
    }

    function _writeCache() {
        try {
            localStorage.setItem(CACHE_KEY, JSON.stringify({
                version: CACHE_VERSION,
                savedAt: new Date().toISOString(),
                characters: characters.map(_cacheSafe),
            }));
        } catch (e) {
            console.warn('[characterRepository] cache write failed:', e);
        }
    }

    // -----------------------------------------------------------------------
    // Legacy image fallback (read-only; canonical images live in the document)
    // -----------------------------------------------------------------------

    function _openLegacyImageDB() {
        if (typeof indexedDB === 'undefined') return Promise.resolve(null);
        return new Promise(resolve => {
            const req = indexedDB.open(LEGACY_DB_NAME, LEGACY_DB_VERSION);
            req.onupgradeneeded = (e) => {
                const db = e.target.result;
                if (!db.objectStoreNames.contains(LEGACY_STORE)) {
                    db.createObjectStore(LEGACY_STORE, { keyPath: 'charId' });
                }
            };
            req.onsuccess = (e) => resolve(e.target.result);
            req.onerror = () => resolve(null);
        });
    }

    async function _loadLegacyImages(charId) {
        const db = await _openLegacyImageDB();
        if (!db) return [];
        return new Promise(resolve => {
            try {
                const tx = db.transaction(LEGACY_STORE, 'readonly');
                const req = tx.objectStore(LEGACY_STORE).get(charId);
                req.onsuccess = (e) => resolve(e.target.result?.images || []);
                req.onerror = () => resolve([]);
            } catch {
                resolve([]);
            }
        });
    }

    async function _hydrateLegacyImages() {
        if (typeof indexedDB === 'undefined') return;
        await Promise.all(characters.map(async (character) => {
            if (Array.isArray(character.howl?.images) && character.howl.images.length > 0) return;
            const images = await _loadLegacyImages(character.id);
            if (images.length === 0) return;
            character.howl = { ...(character.howl || {}), images };
            if (!character.avatar) {
                const thumb = images.find(img => img.isThumbnail) || images[0];
                character.avatar = thumb?.data || character.avatar || '';
            }
        }));
    }

    // -----------------------------------------------------------------------
    // Load — backend first, versioned cache only when the backend is unreachable
    // -----------------------------------------------------------------------

    async function _load() {
        const app = _app();
        let backendList = null;

        if (typeof app?.GetCharacterCutoverStatus === 'function') {
            try {
                cutoverStatus = await app.GetCharacterCutoverStatus();
                cutover = Boolean(cutoverStatus?.cutover);
            } catch (e) {
                console.warn('[characterRepository] cutover status failed:', e);
            }
        }

        if (typeof app?.GetCharacters === 'function') {
            try {
                const list = await app.GetCharacters();
                if (Array.isArray(list)) backendList = list;
            } catch (e) {
                console.warn('[characterRepository] backend load failed:', e);
            }
        }

        if (backendList) {
            // Before the explicit cutover an empty canonical list must not erase
            // the only copy of the user's pre-cutover characters.
            if (!cutover && backendList.length === 0) {
                const legacy = _readLegacy();
                if (legacy) {
                    characters = _normalize(legacy);
                    source = 'legacy';
                    _emit('howl:characters-loaded', { source });
                    return characters.slice();
                }
            }
            characters = _normalize(backendList);
            source = 'backend';
            await _hydrateLegacyImages();
            _writeCache();
            _emit('howl:characters-loaded', { source });
            return characters.slice();
        }

        const cached = _readCache();
        if (cached) {
            characters = _normalize(cached);
            source = 'cache';
        } else {
            const legacy = _readLegacy();
            characters = _normalize(legacy);
            source = legacy ? 'legacy' : 'defaults';
        }
        _emit('howl:characters-loaded', { source });
        return characters.slice();
    }

    function loadAll() {
        if (loadPromise) return loadPromise;
        loadPromise = _load().finally(() => { loadPromise = null; });
        return loadPromise;
    }

    // -----------------------------------------------------------------------
    // Writes — backend first; cache and projection only after acknowledgement
    // -----------------------------------------------------------------------

    async function save(character) {
        if (!character || typeof character !== 'object') throw new Error('Character is required');
        const app = _app();
        if (typeof app?.SaveCharacter !== 'function') {
            throw new Error('character backend unavailable');
        }

        const existing = character.id ? _find(character.id) : null;
        const merged = mergeRecord(existing, character);
        if (existing?.revision && !merged.revision) merged.revision = existing.revision;

        const saved = await app.SaveCharacter(merged);
        if (!saved || typeof saved !== 'object') {
            throw new Error('character backend returned no document');
        }
        const record = _normalizeOne(saved, characters.length);
        _replace(record);

        // Derived cache, written only after the authority acknowledged.
        _writeCache();
        _emit('howl:character-saved', { id: record.id });
        return record;
    }

    async function remove(id) {
        const app = _app();
        if (typeof app?.DeleteCharacter !== 'function') {
            throw new Error('character backend unavailable');
        }
        await app.DeleteCharacter(id);
        characters = characters.filter(c => c.id !== id);
        _writeCache();
        _emit('howl:character-deleted', { id });
    }

    async function loadImages(charId) {
        const character = _find(charId);
        if (Array.isArray(character?.howl?.images) && character.howl.images.length > 0) {
            return character.howl.images;
        }
        const legacy = await _loadLegacyImages(charId);
        if (legacy.length > 0 && character) {
            character.howl = { ...(character.howl || {}), images: legacy };
        }
        return legacy;
    }

    async function importLegacy() {
        const app = _app();
        if (typeof app?.ImportCharacters !== 'function') {
            throw new Error('character backend unavailable');
        }
        const legacy = _readLegacy();
        if (!legacy || legacy.length === 0) {
            throw new Error('No pre-cutover characters were found in this browser.');
        }
        localStorage.setItem(LEGACY_BACKUP_KEY, JSON.stringify(legacy));
        let result;
        try {
            result = await app.ImportCharacters(legacy);
        } catch (error) {
            localStorage.removeItem(LEGACY_BACKUP_KEY);
            throw error;
        }
        await loadAll();
        return result;
    }

    async function rollback() {
        const app = _app();
        if (typeof app?.RollbackCharacters !== 'function') {
            throw new Error('character backend unavailable');
        }
        await app.RollbackCharacters();
        const backup = localStorage.getItem(LEGACY_BACKUP_KEY);
        if (backup !== null) {
            localStorage.setItem(CACHE_KEY, backup);
            localStorage.removeItem(LEGACY_BACKUP_KEY);
        }
        await loadAll();
    }

    // Seed synchronously from the derived cache so pages that read during
    // parsing see the last acknowledged list; the backend refresh follows.
    const seeded = _readCache();
    if (seeded) {
        characters = _normalize(seeded);
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
        all() { return characters; },
        get(id) { return _find(id); },
        loadAll,
        importLegacy,
        rollback,
        save,
        remove,
        loadImages,
        mergeRecord,
    };
})();

window.CharacterStorage = CharacterStorage;
window.CharacterRepository = CharacterStorage;
