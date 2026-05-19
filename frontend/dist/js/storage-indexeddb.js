/**
 * Storage — IndexedDB & Primitive Layer
 *
 * Low-level IndexedDB setup and read/write/delete primitives
 * for the image and chat-data object stores, plus localStorage
 * fallback helpers, the image reference registry, and utility methods.
 *
 * Patched onto StorageManager.prototype — must load after
 * storage-manager.js defines the class.
 *
 * @module storage-indexeddb
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// IndexedDB Initialization
// ---------------------------------------------------------------------------

StorageManager.prototype.initializeIndexedDB = function () {
    return new Promise((resolve, reject) => {
        const request = indexedDB.open(this.dbName, this.dbVersion);
        request.onerror     = () => reject(request.error);
        request.onsuccess   = () => { this.db = request.result; resolve(); };
        request.onupgradeneeded = (event) => {
            const db = event.target.result;
            if (!db.objectStoreNames.contains('scenario-images'))  db.createObjectStore('scenario-images',  { keyPath: 'id' });
            if (!db.objectStoreNames.contains('character-images')) db.createObjectStore('character-images', { keyPath: 'id' });
            if (!db.objectStoreNames.contains('world-images'))     db.createObjectStore('world-images',     { keyPath: 'id' });
            if (!db.objectStoreNames.contains('chat-data'))        db.createObjectStore('chat-data',        { keyPath: 'key' });
        };
    });
};

// ---------------------------------------------------------------------------
// Image — IndexedDB
// ---------------------------------------------------------------------------

StorageManager.prototype.storeImageIndexedDB = function (imageData) {
    const storeName   = `${imageData.type}-images`;
    const transaction = this.db.transaction([storeName], 'readwrite');
    const store       = transaction.objectStore(storeName);
    return new Promise((resolve, reject) => {
        const req = store.put(imageData);
        req.onsuccess = () => resolve(req.result);
        req.onerror   = () => reject(req.error);
    });
};

StorageManager.prototype.getImageIndexedDB = function (id, type) {
    const storeName   = `${type}-images`;
    const transaction = this.db.transaction([storeName], 'readonly');
    const store       = transaction.objectStore(storeName);
    return new Promise((resolve, reject) => {
        const req = store.get(id);
        req.onsuccess = () => resolve(req.result);
        req.onerror   = () => reject(req.error);
    });
};

StorageManager.prototype.deleteImageIndexedDB = function (id, type) {
    const storeName   = `${type}-images`;
    const transaction = this.db.transaction([storeName], 'readwrite');
    const store       = transaction.objectStore(storeName);
    return new Promise((resolve, reject) => {
        const req = store.delete(id);
        req.onsuccess = () => resolve();
        req.onerror   = () => reject(req.error);
    });
};

// ---------------------------------------------------------------------------
// Image — localStorage Fallback
// ---------------------------------------------------------------------------

StorageManager.prototype.storeImageLocalStorage = function (imageData) {
    const dataSize = JSON.stringify(imageData).length;
    if (dataSize > this.maxLocalStorageSize) {
        throw new Error(`Image too large for localStorage (${(dataSize / 1024 / 1024).toFixed(2)}MB)`);
    }
    const key = `howl-chat-${imageData.type}-image-${imageData.id}`;
    localStorage.setItem(key, JSON.stringify(imageData));
    this.updateImageReferences(imageData);
};

StorageManager.prototype.getImageLocalStorage = async function (id, type) {
    const key  = `howl-chat-${type}-image-${id}`;
    const data = localStorage.getItem(key);
    return data ? JSON.parse(data) : null;
};

StorageManager.prototype.deleteImageLocalStorage = async function (id, type) {
    localStorage.removeItem(`howl-chat-${type}-image-${id}`);
};

// ---------------------------------------------------------------------------
// Image Reference Registry
// ---------------------------------------------------------------------------

StorageManager.prototype.updateImageReferences = function (imageData) {
    const registry = JSON.parse(localStorage.getItem('howl-chat-image-registry') || '{}');
    registry[`${imageData.type}-${imageData.id}`] = {
        id: imageData.id, type: imageData.type, name: imageData.name,
        size: imageData.size, timestamp: imageData.timestamp
    };
    localStorage.setItem('howl-chat-image-registry', JSON.stringify(registry));
};

StorageManager.prototype.removeImageReference = function (id, type) {
    const registry = JSON.parse(localStorage.getItem('howl-chat-image-registry') || '{}');
    delete registry[`${type}-${id}`];
    localStorage.setItem('howl-chat-image-registry', JSON.stringify(registry));
};

// ---------------------------------------------------------------------------
// Utility
// ---------------------------------------------------------------------------

StorageManager.prototype.generateId = function () {
    return Date.now().toString(36) + Math.random().toString(36).substr(2);
};

StorageManager.prototype.estimateSize = function (obj) {
    return new Blob([typeof obj === 'string' ? obj : JSON.stringify(obj)]).size;
};

// ---------------------------------------------------------------------------
// Storage Stats
// ---------------------------------------------------------------------------

StorageManager.prototype.getLocalStorageUsage = function () {
    let totalSize = 0, itemCount = 0;
    for (const key in localStorage) {
        if (localStorage.hasOwnProperty(key)) { totalSize += localStorage[key].length + key.length; itemCount++; }
    }
    return {
        used: totalSize, usedMB: (totalSize / 1024 / 1024).toFixed(2),
        max: this.maxLocalStorageSize, maxMB: (this.maxLocalStorageSize / 1024 / 1024).toFixed(2),
        percentage: ((totalSize / this.maxLocalStorageSize) * 100).toFixed(1), items: itemCount
    };
};

StorageManager.prototype.getIndexedDBUsage = async function () {
    if (!this.db) return null;
    try {
        const est = await navigator.storage.estimate();
        return {
            quota: est.quota, usage: est.usage,
            usageMB: (est.usage / 1024 / 1024).toFixed(2),
            quotaMB: (est.quota / 1024 / 1024).toFixed(2),
            percentage: ((est.usage / est.quota) * 100).toFixed(1)
        };
    } catch (error) { return { error: error.message }; }
};
