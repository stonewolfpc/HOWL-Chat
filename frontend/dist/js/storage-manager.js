/**
 * HOWL Chat Storage Manager — Orchestrator
 *
 * Defines the StorageManager class shell with constructor,
 * availability check, initialization, and the three public
 * image methods (storeImage, getImage, deleteImage) plus
 * getStorageStats. All other methods are patched onto
 * the prototype by the modules listed below.
 *
 * Depends on (must load after this file):
 *   storage-indexeddb.js  — DB init, low-level primitives, utilities
 *   storage-scenarios.js  — scenario CRUD + image compaction
 *   storage-chat.js       — chat-data CRUD
 *
 * @module storage-manager
 * @version 2.0.0
 */

'use strict';

class StorageManager {
    constructor() {
        this.dbName              = 'HOWLChatStorage';
        this.dbVersion           = 2;
        this.db                  = null;
        this.isIndexedDBAvailable = this.checkIndexedDBSupport();
        this.maxLocalStorageSize = 4 * 1024 * 1024;
        this.initialized         = false;
    }

    checkIndexedDBSupport() {
        return 'indexedDB' in window && indexedDB !== null && typeof indexedDB.open === 'function';
    }

    async initialize() {
        if (this.initialized) return;
        try {
            if (this.isIndexedDBAvailable) {
                await this.initializeIndexedDB();
                console.log('IndexedDB initialized successfully');
            } else {
                console.warn('IndexedDB not available, using localStorage fallback');
            }
            this.initialized = true;
        } catch (error) {
            console.error('Storage initialization failed:', error);
            this.isIndexedDBAvailable = false;
        }
    }

    // -----------------------------------------------------------------------
    // Image Public API
    // -----------------------------------------------------------------------

    async storeImage(imageData, type = 'scenario', id = null) {
        await this.initialize();
        const obj = {
            id:        id || this.generateId(),
            type,
            name:      imageData.name,
            mimeType:  imageData.type,
            size:      imageData.size,
            data:      imageData.data,        // thumbnail base64
            medium:    imageData.medium,      // medium base64 (optional)
            timestamp: Date.now()
        };

        // Store original blob separately in IndexedDB if available
        if (imageData.original && imageData.original instanceof Blob) {
            obj.originalBlob = imageData.original;
        }

        if (this.isIndexedDBAvailable) {
            try { await this.storeImageIndexedDB(obj); return obj.id; }
            catch (error) { console.warn('IndexedDB storage failed, trying localStorage:', error); }
        }
        // Fallback to localStorage (strips blob, keeps thumbnail only)
        const fallbackObj = { ...obj };
        delete fallbackObj.originalBlob;
        try { await this.storeImageLocalStorage(fallbackObj); return obj.id; }
        catch (error) { throw new Error('Unable to store image - storage quota exceeded'); }
    }

    async getImage(id, type = 'scenario') {
        await this.initialize();
        if (this.isIndexedDBAvailable) {
            try { const img = await this.getImageIndexedDB(id, type); if (img) return img; }
            catch (error) { console.warn('IndexedDB retrieval failed:', error); }
        }
        try   { return await this.getImageLocalStorage(id, type); }
        catch { return null; }
    }

    async deleteImage(id, type = 'scenario') {
        await this.initialize();
        if (this.isIndexedDBAvailable) {
            try { await this.deleteImageIndexedDB(id, type); }
            catch (error) { console.warn('IndexedDB deletion failed:', error); }
        }
        try   { await this.deleteImageLocalStorage(id, type); }
        catch (error) { console.warn('localStorage deletion failed:', error); }
        this.removeImageReference(id, type);
    }

    async getStorageStats() {
        return {
            localStorage: this.getLocalStorageUsage(),
            indexedDB:    this.isIndexedDBAvailable ? await this.getIndexedDBUsage() : null,
            available:    this.isIndexedDBAvailable
        };
    }
}

// ---------------------------------------------------------------------------
// Singleton
// ---------------------------------------------------------------------------

window.HOWLStorage = new StorageManager();
