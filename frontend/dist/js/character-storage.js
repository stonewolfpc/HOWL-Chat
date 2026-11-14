/**
 * Character Storage
 *
 * Two-tier persistence for character data:
 *   - Metadata (all fields except image blobs) → localStorage key 'howl-chat-characters'
 *   - Image blobs (base64 data)                → IndexedDB store 'character-images'
 *
 * Consumers call:
 *   CharacterStorage.save(character)           → Promise<void>
 *   CharacterStorage.loadAll()                 → Promise<character[]>  (images rehydrated)
 *   CharacterStorage.remove(id)                → Promise<void>
 *
 * @module character-storage
 * @version 1.0.0
 */

'use strict';

const CharacterStorage = (() => {

    // -----------------------------------------------------------------------
    // IndexedDB bootstrap
    // -----------------------------------------------------------------------

    const DB_NAME    = 'howl-character-images';
    const DB_VERSION = 1;
    const STORE_NAME = 'images';

    /** @returns {Promise<IDBDatabase>} */
    function openDB() {
        return new Promise((resolve, reject) => {
            const req = indexedDB.open(DB_NAME, DB_VERSION);
            req.onupgradeneeded = (e) => {
                const db = e.target.result;
                if (!db.objectStoreNames.contains(STORE_NAME)) {
                    db.createObjectStore(STORE_NAME, { keyPath: 'charId' });
                }
            };
            req.onsuccess = (e) => resolve(e.target.result);
            req.onerror   = (e) => reject(e.target.error);
        });
    }

    // -----------------------------------------------------------------------
    // Image helpers
    // -----------------------------------------------------------------------

    /**
     * Persist an array of image objects for a character.
     * Stored as { charId, images: [...] } — replaces any prior entry.
     * @param {string}   charId
     * @param {object[]} images   Full image objects including .data blobs
     * @returns {Promise<void>}
     */
    async function saveImages(charId, images) {
        const db = await openDB();
        return new Promise((resolve, reject) => {
            const tx    = db.transaction(STORE_NAME, 'readwrite');
            const store = tx.objectStore(STORE_NAME);
            const req   = store.put({ charId, images });
            req.onsuccess = () => resolve();
            req.onerror   = (e) => reject(e.target.error);
        });
    }

    /**
     * Retrieve image objects for a character.
     * @param {string} charId
     * @returns {Promise<object[]>}
     */
    async function loadImages(charId) {
        const db = await openDB();
        return new Promise((resolve, reject) => {
            const tx    = db.transaction(STORE_NAME, 'readonly');
            const store = tx.objectStore(STORE_NAME);
            const req   = store.get(charId);
            req.onsuccess = (e) => resolve(e.target.result?.images || []);
            req.onerror   = (e) => reject(e.target.error);
        });
    }

    /**
     * Delete stored images for a character.
     * @param {string} charId
     * @returns {Promise<void>}
     */
    async function deleteImages(charId) {
        const db = await openDB();
        return new Promise((resolve, reject) => {
            const tx    = db.transaction(STORE_NAME, 'readwrite');
            const store = tx.objectStore(STORE_NAME);
            const req   = store.delete(charId);
            req.onsuccess = () => resolve();
            req.onerror   = (e) => reject(e.target.error);
        });
    }

    // -----------------------------------------------------------------------
    // Metadata helpers (localStorage — no blobs)
    // -----------------------------------------------------------------------

    const LS_KEY = 'howl-chat-characters';

    /**
     * Merge a character patch into an existing stored record.
     *
     * Objects are merged recursively so an edit to one voice/metadata field
     * cannot erase its siblings. Arrays and scalar values are replaced, which
     * keeps explicit clears and list edits authoritative.
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

    /** Strip image .data and .original blobs before writing to localStorage. */
    function stripImageData(character) {
        const stripped = { ...character };
        if (stripped.howl?.images) {
            stripped.howl = {
                ...stripped.howl,
                images: stripped.howl.images.map(img => ({
                    id:          img.id,
                    name:        img.name,
                    size:        img.size,
                    isThumbnail: img.isThumbnail
                    // 'data' and 'original' intentionally omitted - stored in IndexedDB only
                }))
            };
        }
        // Also strip top-level avatar blob if present
        if (stripped.avatar && stripped.avatar.startsWith('data:')) {
            stripped.avatar = '';
        }
        return stripped;
    }

    function readMetadata() {
        try {
            return JSON.parse(localStorage.getItem(LS_KEY) || '[]');
        } catch {
            return [];
        }
    }

    function writeMetadata(characters) {
        localStorage.setItem(LS_KEY, JSON.stringify(characters));
    }

    // -----------------------------------------------------------------------
    // Public API
    // -----------------------------------------------------------------------

    /**
     * Persist a character: metadata → localStorage, images → IndexedDB.
     * Preserves list order; upserts by character.id.
     * @param {object} character
     * @returns {Promise<void>}
     */
    async function save(character) {
        if (!character || !character.id) throw new Error('Character must have an id');

        const all = readMetadata();
        const idx = all.findIndex(c => c.id === character.id);
        const existing = idx >= 0 ? all[idx] : null;

        // The caller may be a partial view of the record (for example the sidebar
        // image action). Merge into the persisted record before writing anything.
        const merged = mergeRecord(existing, character);

        // Only replace image blobs when the caller supplied an images field.
        // A partial save with no images intent must not erase existing blobs.
        const howlPatch = character.howl;
        const replacesImages = howlPatch && typeof howlPatch === 'object'
            && Object.prototype.hasOwnProperty.call(howlPatch, 'images');
        if (replacesImages) {
            await saveImages(merged.id, Array.isArray(howlPatch.images) ? howlPatch.images : []);
        }

        const images = Array.isArray(merged.howl?.images) ? merged.howl.images : [];

        // Determine thumbnail src for display (first thumbnail or first image)
        const thumb   = images.find(img => img.isThumbnail) || images[0];
        const stripped = stripImageData(merged);
        // Store a lightweight avatarThumb reference for the grid (no full blob)
        stripped._thumbId = thumb ? (thumb.id || null) : null;

        if (idx >= 0) {
            all[idx] = stripped;
        } else {
            all.push(stripped);
        }
        writeMetadata(all);
        if (typeof window !== 'undefined' && typeof window.dispatchEvent === 'function') {
            window.dispatchEvent(new CustomEvent('howl:character-saved', {
                detail: { id: merged.id }
            }));
        }
    }

    /**
     * Load all characters with images rehydrated from IndexedDB.
     * @returns {Promise<object[]>}
     */
    async function loadAll() {
        const all = readMetadata();
        const hydrated = await Promise.all(all.map(async (char) => {
            const images = await loadImages(char.id).catch(() => []);
            if (images.length === 0) return char;
            return {
                ...char,
                howl: {
                    ...(char.howl || {}),
                    images
                },
                avatar: (() => {
                    const thumb = images.find(img => img.isThumbnail) || images[0];
                    return thumb?.data || char.avatar || '';
                })()
            };
        }));
        return hydrated;
    }

    /**
     * Remove a character by id from both localStorage and IndexedDB.
     * @param {string} id
     * @returns {Promise<void>}
     */
    async function remove(id) {
        await deleteImages(id).catch(() => {});
        const all     = readMetadata();
        const updated = all.filter(c => c.id !== id);
        writeMetadata(updated);
    }

    return { save, loadAll, remove, loadImages, mergeRecord };

})();

window.CharacterStorage = CharacterStorage;
