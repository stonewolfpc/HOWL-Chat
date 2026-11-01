/**
 * App Wipe — Frontend Data Reset
 *
 * Clears all localStorage, IndexedDB, and reloads the page.
 * Called by the backend WipeAllData flow.
 *
 * @module app-wipe
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// localStorage Keys Used by HOWL Chat
// ---------------------------------------------------------------------------
const HOWL_LOCALSTORAGE_KEYS = [
    'howl-chat-settings',
    'howl-worlds',
    'howl-chat-scenarios',
    'howl-chat-characters',
    'howl-chat-lorebooks',
    'howl-chat-memories',
    'howl-chat-current-chat',
    'howl-chat-image-generator-state',
    'howl-theme-mode' // legacy
];

// ---------------------------------------------------------------------------
// Wipe Frontend Data
// ---------------------------------------------------------------------------

async function wipeFrontendData() {
    console.log('[Wipe] Clearing frontend data...');

    // 1. Clear all HOWL localStorage keys
    let lsCleared = 0;
    HOWL_LOCALSTORAGE_KEYS.forEach(key => {
        if (localStorage.getItem(key) !== null) {
            localStorage.removeItem(key);
            lsCleared++;
        }
    });
    console.log(`[Wipe] Cleared ${lsCleared} localStorage entries`);

    // 2. Clear IndexedDB
    try {
        if (window.HOWLStorage && window.HOWLStorage.db) {
            const db = window.HOWLStorage.db;
            const stores = ['scenario-images', 'character-images', 'world-images', 'chat-data'];
            for (const storeName of stores) {
                try {
                    const tx = db.transaction(storeName, 'readwrite');
                    const store = tx.objectStore(storeName);
                    await new Promise((resolve, reject) => {
                        const req = store.clear();
                        req.onsuccess = () => resolve();
                        req.onerror = () => reject(req.error);
                    });
                    console.log(`[Wipe] Cleared IndexedDB store: ${storeName}`);
                } catch (e) {
                    console.warn(`[Wipe] Could not clear store ${storeName}:`, e);
                }
            }
        } else {
            // Fallback: delete and recreate the entire DB
            await new Promise((resolve, reject) => {
                const req = indexedDB.deleteDatabase('HOWLChatStorage');
                req.onsuccess = () => { console.log('[Wipe] IndexedDB deleted'); resolve(); };
                req.onerror = () => { console.warn('[Wipe] IndexedDB delete failed:', req.error); resolve(); };
                req.onblocked = () => { console.warn('[Wipe] IndexedDB delete blocked'); resolve(); };
            });
        }
    } catch (e) {
        console.warn('[Wipe] IndexedDB wipe failed:', e);
    }

    // 3. Clear sessionStorage
    sessionStorage.clear();
    console.log('[Wipe] sessionStorage cleared');

    console.log('[Wipe] Frontend data cleared. Reloading...');
    window.location.reload();
}

// ---------------------------------------------------------------------------
// Full App Reset Entry Point
// ---------------------------------------------------------------------------

async function resetApplication() {
    if (!confirm('WARNING: This will delete ALL data:\n\n' +
        '- All chat sessions\n' +
        '- All characters, worlds, and scenarios\n' +
        '- All memories and lorebooks\n' +
        '- All generated gallery images\n' +
        '- All settings (except model paths)\n\n' +
        'This CANNOT be undone. Are you sure?')) {
        return;
    }

    if (!confirm('FINAL CONFIRMATION: Every saved item will be permanently erased. Proceed?')) {
        return;
    }

    try {
        // Call backend wipe first
        if (window.go?.main?.App?.WipeAllData) {
            await window.go.main.App.WipeAllData();
        }
    } catch (err) {
        console.error('[Wipe] Backend wipe failed:', err);
    }

    // Then wipe frontend and reload
    await wipeFrontendData();
}

// Expose globally
window.resetApplication = resetApplication;
window.wipeFrontendData = wipeFrontendData;
