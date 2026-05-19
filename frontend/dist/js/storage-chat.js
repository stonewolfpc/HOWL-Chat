/**
 * Storage — Chat Data CRUD
 *
 * setChatData, getChatData, deleteChatData, and clearOldChatData
 * with IndexedDB primary and localStorage fallback.
 *
 * Patched onto StorageManager.prototype — must load after
 * storage-manager.js and storage-indexeddb.js.
 *
 * @module storage-chat
 * @version 1.0.0
 */

'use strict';

StorageManager.prototype.setChatData = async function (key, data) {
    await this.initialize();
    const size    = this.estimateSize(data);
    const dataObj = { key, data, size, timestamp: Date.now() };

    if (this.isIndexedDBAvailable && this.db && this.db.objectStoreNames.contains('chat-data')) {
        try {
            return await new Promise((resolve, reject) => {
                const tx    = this.db.transaction(['chat-data'], 'readwrite');
                const store = tx.objectStore('chat-data');
                const req   = store.put(dataObj);
                req.onerror   = () => reject(req.error);
                req.onsuccess = () => resolve(true);
            });
        } catch (error) {
            console.warn('IndexedDB write failed, falling back to localStorage:', error);
        }
    }

    try {
        if (size < this.maxLocalStorageSize) {
            localStorage.setItem(key, JSON.stringify(data));
            return true;
        }
    } catch (error) {
        if (error.name === 'QuotaExceededError') {
            console.warn('localStorage quota exceeded, clearing old entries...');
            await this.clearOldChatData();
            try { localStorage.setItem(key, JSON.stringify(data)); return true; }
            catch { return false; }
        }
    }
    return false;
};

StorageManager.prototype.getChatData = async function (key) {
    await this.initialize();

    if (this.isIndexedDBAvailable && this.db && this.db.objectStoreNames.contains('chat-data')) {
        try {
            const result = await new Promise((resolve, reject) => {
                const tx    = this.db.transaction(['chat-data'], 'readonly');
                const store = tx.objectStore('chat-data');
                const req   = store.get(key);
                req.onerror   = () => reject(req.error);
                req.onsuccess = () => resolve(req.result);
            });
            if (result) return result.data;
        } catch (error) {
            console.warn('IndexedDB read failed:', error);
        }
    }

    try {
        const stored = localStorage.getItem(key);
        return stored ? JSON.parse(stored) : null;
    } catch (error) {
        console.error('Failed to retrieve chat data:', error);
        return null;
    }
};

StorageManager.prototype.deleteChatData = async function (key) {
    await this.initialize();

    if (this.isIndexedDBAvailable && this.db && this.db.objectStoreNames.contains('chat-data')) {
        try {
            await new Promise((resolve, reject) => {
                const tx    = this.db.transaction(['chat-data'], 'readwrite');
                const store = tx.objectStore('chat-data');
                const req   = store.delete(key);
                req.onerror   = () => reject(req.error);
                req.onsuccess = () => resolve(true);
            });
        } catch (error) {
            console.warn('Failed to delete from IndexedDB:', error);
        }
    }

    try { localStorage.removeItem(key); }
    catch (error) { console.warn('Failed to delete from localStorage:', error); }
};

StorageManager.prototype.clearOldChatData = async function () {
    if (!this.isIndexedDBAvailable || !this.db) return;
    try {
        const tx    = this.db.transaction(['chat-data'], 'readwrite');
        const store = tx.objectStore('chat-data');
        const all   = await new Promise((resolve, reject) => {
            const req = store.getAll();
            req.onerror   = () => reject(req.error);
            req.onsuccess = () => resolve(req.result);
        });
        if (all.length > 3) {
            all.sort((a, b) => b.timestamp - a.timestamp);
            for (let i = 3; i < all.length; i++) {
                await new Promise((resolve, reject) => {
                    const req = store.delete(all[i].key);
                    req.onerror   = () => reject(req.error);
                    req.onsuccess = () => resolve(true);
                });
            }
        }
    } catch (error) { console.error('Failed to clear old chat data:', error); }
};
