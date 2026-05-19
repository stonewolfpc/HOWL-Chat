/**
 * Storage — Scenario CRUD
 *
 * Scenario store/get/update/delete methods plus the image
 * compaction helper and reference builder.
 *
 * Patched onto StorageManager.prototype — must load after
 * storage-manager.js and storage-indexeddb.js.
 *
 * @module storage-scenarios
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Image Reference Builder
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Strip all blob data from a scenario object before localStorage write.
// Images go to IndexedDB; avatars are on the character record.
// ---------------------------------------------------------------------------

function _stripScenarioBlobs(scenario) {
    const s = { ...scenario };
    // Strip legacy single-image blob
    if (s.image && s.image.data) {
        s.image = { id: s.image.id || null, name: s.image.name, type: s.image.type, size: s.image.size, stored: !!s.image.stored };
    }
    // Strip multi-image array blobs
    if (s.images && Array.isArray(s.images)) {
        s.images = s.images.map(img => {
            if (!img) return img;
            if (img.data || img.medium || img.original) {
                return { id: img.id || null, name: img.name, type: img.type, size: img.size, stored: !!img.stored };
            }
            return img;
        });
    }
    // Strip character avatar blobs — these can be 100s of KB each
    if (Array.isArray(s.characters)) {
        s.characters = s.characters.map(c => {
            if (c.avatar && c.avatar.startsWith('data:')) {
                const { avatar, ...rest } = c;
                return rest;
            }
            return c;
        });
    }
    return s;
}

StorageManager.prototype.scenarioImageReference = function (imageMeta) {
    return { id: imageMeta.id, name: imageMeta.name, type: imageMeta.type, size: imageMeta.size, stored: true };
};

// ---------------------------------------------------------------------------
// Image Compaction
// ---------------------------------------------------------------------------

StorageManager.prototype.compactScenarioImages = async function (scenarios) {
    let changed = false;
    for (let i = 0; i < scenarios.length; i++) {
        const scenario = scenarios[i];
        // Compact legacy single image
        if (scenario.image && scenario.image.data) {
            try {
                const imageId = await this.storeImage(scenario.image, 'scenario', scenario.image.id || null);
                scenarios[i] = {
                    ...scenario,
                    image: { id: imageId, name: scenario.image.name, type: scenario.image.type, size: scenario.image.size, stored: true }
                };
                changed = true;
            } catch (error) {
                console.warn(`Failed to compact image for scenario "${scenario.title || i}":`, error);
            }
        }
        // Compact multi-image array
        if (scenario.images && Array.isArray(scenario.images)) {
            const compactedImages = [];
            let imagesChanged = false;
            for (const img of scenario.images) {
                if (!img) { compactedImages.push(img); continue; }
                if (img.data) {
                    try {
                        const imageId = await this.storeImage(img, 'scenario', img.id || null);
                        compactedImages.push({ id: imageId, name: img.name, type: img.type, size: img.size, stored: true });
                        imagesChanged = true;
                    } catch (error) {
                        console.warn(`Failed to compact image for scenario "${scenario.title || i}":`, error);
                        compactedImages.push(img);
                    }
                } else {
                    compactedImages.push(img);
                }
            }
            if (imagesChanged) {
                scenarios[i] = { ...scenarios[i], images: compactedImages };
                changed = true;
            }
        }
    }
    if (changed) {
        try { localStorage.setItem('howl-chat-scenarios', JSON.stringify(scenarios)); }
        catch (error) { console.warn('Could not persist compacted scenarios:', error); }
    }
    return changed;
};

// ---------------------------------------------------------------------------
// Scenario Store
// ---------------------------------------------------------------------------

StorageManager.prototype.storeScenario = async function (scenarioData) {
    await this.initialize();
    const scenarios = JSON.parse(localStorage.getItem('howl-chat-scenarios') || '[]');
    await this.compactScenarioImages(scenarios);

    const processedData = { ...scenarioData };
    if (!processedData.id) processedData.id = this.generateId();

    // Process legacy single image
    if (scenarioData.image && scenarioData.image.data) {
        try {
            const imageId = await this.storeImage(scenarioData.image, 'scenario');
            processedData.image = { id: imageId, name: scenarioData.image.name, type: scenarioData.image.type, size: scenarioData.image.size, stored: true };
        } catch (error) {
            console.warn('Failed to store scenario image:', error);
            processedData.image = null;
        }
    }

    // Process multi-image array
    if (scenarioData.images && Array.isArray(scenarioData.images)) {
        const storedImages = [];
        for (const img of scenarioData.images) {
            if (!img) continue;
            if (img.data) {
                try {
                    const imageId = await this.storeImage(img, 'scenario');
                    storedImages.push({ id: imageId, name: img.name, type: img.type, size: img.size, stored: true });
                } catch (error) {
                    console.warn('Failed to store scenario image:', img.name, error);
                    storedImages.push(img);
                }
            } else {
                storedImages.push(img);
            }
        }
        processedData.images = storedImages;
    }

    // Strip all blob data (thumbnail + character avatars) before localStorage write
    const lsData = _stripScenarioBlobs(processedData);

    try {
        const fresh = JSON.parse(localStorage.getItem('howl-chat-scenarios') || '[]');
        await this.compactScenarioImages(fresh);
        fresh.push(lsData);
        try {
            localStorage.setItem('howl-chat-scenarios', JSON.stringify(fresh));
        } catch (quotaErr) {
            if (quotaErr.name === 'QuotaExceededError') {
                // Strip blobs from all existing entries and retry
                localStorage.setItem('howl-chat-scenarios', JSON.stringify(fresh.map(_stripScenarioBlobs)));
            } else {
                throw quotaErr;
            }
        }
        return true;
    } catch (error) {
        console.error('Failed to store scenario:', error);
        throw error;
    }
};

// ---------------------------------------------------------------------------
// Scenario Get
// ---------------------------------------------------------------------------

StorageManager.prototype.getScenario = async function (scenarioId) {
    const scenarios = JSON.parse(localStorage.getItem('howl-chat-scenarios') || '[]');
    const scenario  = scenarios.find(s => s.id === scenarioId);
    if (!scenario) return null;

    if (scenario.image && scenario.image.stored) {
        try {
            const imageData = await this.getImage(scenario.image.id, 'scenario');
            scenario.image  = imageData || null;
        } catch (error) {
            console.error('Failed to restore scenario image:', error);
            scenario.image = null;
        }
    }
    return scenario;
};

// ---------------------------------------------------------------------------
// Scenario Update
// ---------------------------------------------------------------------------

StorageManager.prototype.updateScenario = async function (index, scenarioData) {
    await this.initialize();
    const scenarios = JSON.parse(localStorage.getItem('howl-chat-scenarios') || '[]');
    if (index < 0 || index >= scenarios.length) throw new Error('Invalid scenario index for update');

    await this.compactScenarioImages(scenarios);

    const processedData = { ...scenarioData };
    const oldScenario   = scenarios[index];

    if (scenarioData.image && scenarioData.image.data) {
        const sameStored = oldScenario.image?.stored && scenarioData.image.id && oldScenario.image.id === scenarioData.image.id;
        if (sameStored) {
            processedData.image = this.scenarioImageReference(oldScenario.image);
        } else {
            try {
                if (oldScenario.image?.stored && oldScenario.image.id !== scenarioData.image.id) {
                    await this.deleteImage(oldScenario.image.id, 'scenario');
                }
                const imageId = await this.storeImage(scenarioData.image, 'scenario', scenarioData.image.id || null);
                processedData.image = this.scenarioImageReference({ id: imageId, name: scenarioData.image.name, type: scenarioData.image.type, size: scenarioData.image.size });
            } catch (error) {
                console.warn('Failed to store updated scenario image:', error);
                processedData.image = oldScenario.image?.stored ? this.scenarioImageReference(oldScenario.image) : null;
            }
        }
    } else if (scenarioData.image?.stored) {
        processedData.image = this.scenarioImageReference(scenarioData.image);
    } else {
        if (oldScenario.image?.stored) await this.deleteImage(oldScenario.image.id, 'scenario');
        processedData.image = null;
    }

    // Process multi-image array
    if (scenarioData.images && Array.isArray(scenarioData.images)) {
        const oldImages = oldScenario.images || [];
        const newImages = [];
        const newIds = new Set();

        for (const img of scenarioData.images) {
            if (!img) continue;
            if (img.data) {
                // New or updated image with base64 data
                try {
                    const imageId = await this.storeImage(img, 'scenario', img.id || null);
                    newImages.push({ id: imageId, name: img.name, type: img.type, size: img.size, stored: true });
                    newIds.add(imageId);
                } catch (error) {
                    console.warn('Failed to store updated scenario image:', img.name, error);
                    newImages.push(img);
                }
            } else if (img.stored && img.id) {
                // Existing stored reference — keep it
                newImages.push(img);
                newIds.add(img.id);
            } else {
                // Plain reference or string
                newImages.push(img);
            }
        }

        // Delete orphaned old images no longer in the new array
        for (const oldImg of oldImages) {
            if (oldImg && oldImg.stored && oldImg.id && !newIds.has(oldImg.id)) {
                await this.deleteImage(oldImg.id, 'scenario').catch(() => {});
            }
        }

        processedData.images = newImages;
    } else if (oldScenario.images && Array.isArray(oldScenario.images)) {
        // All images removed — delete old ones
        for (const oldImg of oldScenario.images) {
            if (oldImg && oldImg.stored && oldImg.id) {
                await this.deleteImage(oldImg.id, 'scenario').catch(() => {});
            }
        }
        processedData.images = [];
    }

    // Strip all blob data (thumbnail + character avatars) before localStorage write
    scenarios[index] = _stripScenarioBlobs(processedData);

    try {
        try {
            localStorage.setItem('howl-chat-scenarios', JSON.stringify(scenarios));
        } catch (quotaErr) {
            if (quotaErr.name === 'QuotaExceededError') {
                // Strip blobs from every entry and retry
                localStorage.setItem('howl-chat-scenarios', JSON.stringify(scenarios.map(_stripScenarioBlobs)));
            } else {
                throw quotaErr;
            }
        }
        return true;
    } catch (error) {
        console.error('Failed to update scenario:', error);
        throw error;
    }
};

// ---------------------------------------------------------------------------
// World Image Helpers
// ---------------------------------------------------------------------------

function _stripWorldBlobs(world) {
    const w = { ...world };
    // Strip legacy single-image blob
    if (w.image && w.image.data) {
        w.image = { id: w.image.id || null, name: w.image.name, type: w.image.type, size: w.image.size, stored: !!w.image.stored };
    }
    // Strip multi-image array blobs
    if (w.images && Array.isArray(w.images)) {
        w.images = w.images.map(img => {
            if (!img) return img;
            if (img.data || img.medium || img.original) {
                return { id: img.id || null, name: img.name, type: img.type, size: img.size, stored: !!img.stored };
            }
            return img;
        });
    }
    return w;
}

// ---------------------------------------------------------------------------
// World Store
// ---------------------------------------------------------------------------

StorageManager.prototype.storeWorld = async function (worldData) {
    await this.initialize();
    const worlds = JSON.parse(localStorage.getItem('howl-worlds') || '[]');

    const processedData = { ...worldData };
    if (!processedData.id) processedData.id = this.generateId();

    // Process legacy single image
    if (worldData.image && worldData.image.data) {
        try {
            const imageId = await this.storeImage(worldData.image, 'world');
            processedData.image = { id: imageId, name: worldData.image.name, type: worldData.image.type, size: worldData.image.size, stored: true };
        } catch (error) {
            console.warn('Failed to store world image:', error);
            processedData.image = null;
        }
    }

    // Process multi-image array
    if (worldData.images && Array.isArray(worldData.images)) {
        const storedImages = [];
        for (const img of worldData.images) {
            if (!img) continue;
            if (img.data) {
                try {
                    const imageId = await this.storeImage(img, 'world');
                    storedImages.push({ id: imageId, name: img.name, type: img.type, size: img.size, stored: true });
                } catch (error) {
                    console.warn('Failed to store world image:', img.name, error);
                    storedImages.push(img);
                }
            } else {
                storedImages.push(img);
            }
        }
        processedData.images = storedImages;
    }

    const lsData = _stripWorldBlobs(processedData);
    const existingIndex = worlds.findIndex(w => w.id === lsData.id);
    if (existingIndex >= 0) {
        if (worlds[existingIndex].image?.stored && worlds[existingIndex].image.id !== lsData.image?.id) {
            await this.deleteImage(worlds[existingIndex].image.id, 'world').catch(() => {});
        }
        // Delete orphaned old multi-image array images no longer referenced
        const oldImages = worlds[existingIndex].images || [];
        const newIds = new Set((lsData.images || []).filter(img => img && img.id).map(img => img.id));
        for (const oldImg of oldImages) {
            if (oldImg && oldImg.stored && oldImg.id && !newIds.has(oldImg.id)) {
                await this.deleteImage(oldImg.id, 'world').catch(() => {});
            }
        }
        worlds[existingIndex] = lsData;
    } else {
        worlds.push(lsData);
    }

    try {
        localStorage.setItem('howl-worlds', JSON.stringify(worlds));
    } catch (quotaErr) {
        if (quotaErr.name === 'QuotaExceededError') {
            localStorage.setItem('howl-worlds', JSON.stringify(worlds.map(_stripWorldBlobs)));
        } else { throw quotaErr; }
    }
    return processedData;
};

// ---------------------------------------------------------------------------
// World Delete
// ---------------------------------------------------------------------------

StorageManager.prototype.deleteWorld = async function (worldId) {
    await this.initialize();
    const worlds = JSON.parse(localStorage.getItem('howl-worlds') || '[]');
    const idx    = worlds.findIndex(w => w.id === worldId);
    if (idx < 0) return;

    if (worlds[idx].image?.stored) {
        await this.deleteImage(worlds[idx].image.id, 'world').catch(() => {});
    }
    // Delete multi-image array images from IndexedDB
    if (worlds[idx].images && Array.isArray(worlds[idx].images)) {
        for (const img of worlds[idx].images) {
            if (img && img.stored && img.id) {
                await this.deleteImage(img.id, 'world').catch(() => {});
            }
        }
    }
    worlds.splice(idx, 1);
    localStorage.setItem('howl-worlds', JSON.stringify(worlds));
};

// ---------------------------------------------------------------------------
// Scenario Delete
// ---------------------------------------------------------------------------

StorageManager.prototype.deleteScenario = async function (index) {
    await this.initialize();
    const scenarios = JSON.parse(localStorage.getItem('howl-chat-scenarios') || '[]');
    if (index < 0 || index >= scenarios.length) throw new Error('Invalid scenario index for deletion');

    const scenario = scenarios[index];
    if (scenario.image && scenario.image.stored) {
        try { await this.deleteImage(scenario.image.id, 'scenario'); }
        catch (error) { console.warn('Failed to delete scenario image:', error); }
    }

    // Delete multi-image array images from IndexedDB
    if (scenario.images && Array.isArray(scenario.images)) {
        for (const img of scenario.images) {
            if (img && img.stored && img.id) {
                try { await this.deleteImage(img.id, 'scenario'); }
                catch (error) { console.warn('Failed to delete scenario image:', error); }
            }
        }
    }

    scenarios.splice(index, 1);
    try {
        localStorage.setItem('howl-chat-scenarios', JSON.stringify(scenarios));
        return true;
    } catch (error) {
        console.error('Failed to delete scenario:', error);
        throw error;
    }
};

// ---------------------------------------------------------------------------
// Scenario Image Restore
// ---------------------------------------------------------------------------

StorageManager.prototype.restoreScenarioImages = async function (scenario) {
    if (!scenario) return scenario;
    const s = { ...scenario };

    // Restore legacy single image
    if (s.image && s.image.stored && s.image.id) {
        try {
            const img = await this.getImage(s.image.id, 'scenario');
            if (img && img.data) s.image = { ...s.image, data: img.data, medium: img.medium };
        } catch (error) { console.warn('Failed to restore scenario image:', error); }
    }

    // Restore multi-image array
    if (s.images && Array.isArray(s.images)) {
        const restoredImages = [];
        for (const img of s.images) {
            if (!img) { restoredImages.push(img); continue; }
            if (img.stored && img.id) {
                try {
                    const stored = await this.getImage(img.id, 'scenario');
                    if (stored && stored.data) {
                        restoredImages.push({ ...img, data: stored.data, medium: stored.medium });
                        continue;
                    }
                } catch (error) { console.warn('Failed to restore scenario image:', img.id, error); }
            }
            restoredImages.push(img);
        }
        s.images = restoredImages;
    }
    return s;
};

// ---------------------------------------------------------------------------
// World Image Restore (multi-image array support)
// ---------------------------------------------------------------------------

StorageManager.prototype.restoreWorldImages = async function (world) {
    if (!world) return world;
    const w = { ...world };

    // Restore legacy single image
    if (w.image && w.image.stored && w.image.id) {
        try {
            const img = await this.getImage(w.image.id, 'world');
            if (img && img.data) w.image = { ...w.image, data: img.data, medium: img.medium };
        } catch (error) { console.warn('Failed to restore world image:', error); }
    }

    // Restore multi-image array
    if (w.images && Array.isArray(w.images)) {
        const restoredImages = [];
        for (const img of w.images) {
            if (!img) { restoredImages.push(img); continue; }
            if (img.stored && img.id) {
                try {
                    const stored = await this.getImage(img.id, 'world');
                    if (stored && stored.data) {
                        restoredImages.push({ ...img, data: stored.data, medium: stored.medium });
                        continue;
                    }
                } catch (error) { console.warn('Failed to restore world image:', img.id, error); }
            }
            restoredImages.push(img);
        }
        w.images = restoredImages;
    }
    return w;
};
