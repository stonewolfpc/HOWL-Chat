/**
 * World Grid
 *
 * Renders the world card grid and handles grid-level
 * actions: edit, copy, delete, and load-for-use (hotswap).
 *
 * Depends on: world-state.js
 *
 * @module world-grid
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Grid Render
// ---------------------------------------------------------------------------

function _resolveWorldImageSrc(world) {
    if (!world) return '';
    // New multi-image array format
    if (world.images && Array.isArray(world.images) && world.images.length > 0) {
        const thumbIdx = world.thumbnailIndex || 0;
        const thumb = world.images[thumbIdx] || world.images[0];
        return thumb.data || thumb.thumbnail || '';
    }
    // Legacy single-image format
    const image = world.image;
    if (!image) return '';
    if (typeof image === 'string') return image;
    if (typeof image === 'object') {
        return image.data || image.thumbnail || image.url || image.path || '';
    }
    return '';
}

async function renderWorldGrid() {
    const grid = document.getElementById('world-grid');
    if (!grid) return;

    const worlds = window.worldState.worlds;

    if (worlds.length === 0) {
        grid.innerHTML = `
            <div class="world-grid-empty">
                <div class="world-grid-empty-icon icon icon-world icon-xl"></div>
                <div class="world-grid-empty-text">No worlds created yet</div>
                <div class="world-grid-empty-subtext">Click "Create New World" to get started</div>
            </div>
        `;
        return;
    }

    // Rehydrate images from IndexedDB before rendering
    let hydratedWorlds = worlds;
    if (window.HOWLStorage) {
        try {
            hydratedWorlds = await Promise.all(worlds.map(w => window.HOWLStorage.restoreWorldImages(w)));
        } catch (error) {
            console.warn('Failed to rehydrate world images:', error);
        }
    }

    // Resolve all images before rendering
    const resolved = await Promise.all(hydratedWorlds.map(async (world) => {
        const w = { ...world };
        w._resolvedImageSrc = _resolveWorldImageSrc(w);
        return w;
    }));

    grid.innerHTML = resolved.map(world => `
        <div class="world-bubble" data-world-id="${world.id}">
            ${world._resolvedImageSrc
                ? `<img class="world-bubble-image" src="${world._resolvedImageSrc}" alt="${world.title}" onerror="this.style.display='none';this.nextElementSibling.style.display='flex';">
                   <div class="world-bubble-placeholder" style="display:none;"><div class="world-bubble-placeholder-icon icon icon-world"></div></div>`
                : `<div class="world-bubble-placeholder"><div class="world-bubble-placeholder-icon icon icon-world"></div></div>`
            }
            <div class="world-bubble-name">${world.title}</div>
            <div class="world-bubble-meta">${world.details || world.description || ''}</div>
            <div class="world-bubble-actions">
                <button class="world-bubble-action" onclick="loadWorldForUse('${world.id}')" title="Load for Chat"><span class="icon icon-rocket"></span></button>
                <button class="world-bubble-action" onclick="editWorld('${world.id}')" title="Edit"><span class="icon icon-edit"></span></button>
                <button class="world-bubble-action" onclick="copyWorld('${world.id}')" title="Copy"><span class="icon icon-clipboard"></span></button>
                <button class="world-bubble-action" onclick="deleteWorld('${world.id}')" title="Delete"><span class="icon icon-trash"></span></button>
            </div>
        </div>
    `).join('');
}

// ---------------------------------------------------------------------------
// Grid Actions
// ---------------------------------------------------------------------------

function editWorld(worldId) {
    openWorldCreationPopup(worldId);
}

async function copyWorld(worldId) {
    const ws    = window.worldState;
    const world = ws.worlds.find(w => w.id === worldId);
    if (!world) return;
    const copy = {
        ...world,
        id:        Date.now().toString(),
        title:     `${world.title} (Copy)`,
        images:    world.images ? world.images.map(img => ({ ...img })) : [],
        thumbnailIndex: world.thumbnailIndex || 0,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString()
    };
    try {
        await window.WorldRepository.save(copy);
        ws.worlds = window.WorldRepository.all();
    } catch (err) {
        console.warn('copyWorld storage failed:', err);
        return;
    }
    renderWorldGrid();
}

async function deleteWorld(worldId) {
    const ws     = window.worldState;
    const target = (ws.worlds || []).find(w => w.id === worldId);
    const label  = target ? (target.title || target.name || 'this world') : 'this world';

    window.DeleteConfirm.show({
        title: `Delete "${label}"?`,
        body:  'The world, its image, and all associated chat references will be removed.',
        onConfirm: async () => {
            if (window.DeleteService) window.DeleteService.detachWorldFromAll(worldId);
            try {
                await window.WorldRepository.remove(worldId);
            } catch (err) {
                console.error('deleteWorld storage failed:', err);
                return;
            }
            await loadWorlds();
        }
    });
}

async function loadWorldForUse(worldId) {
    const world = window.worldState.worlds.find(w => w.id === worldId);
    if (!world) { console.error('World not found:', worldId); return; }

    console.log(`Loading world "${world.title}" for active use...`);

    if (world.scenarios && world.scenarios.length > 0) {
        world.scenarios.forEach(s => console.log(`  - Scenario: ${s.title || 'Untitled'}`));
        alert(`World "${world.title}" loaded with ${world.scenarios.length} scenarios for active chat use.`);
    } else {
        alert(`World "${world.title}" has no attached scenarios.`);
    }
}

// ---------------------------------------------------------------------------
// Load Worlds from Storage
// ---------------------------------------------------------------------------

async function loadWorlds() {
    if (window.WorldRepository?.ready) await window.WorldRepository.ready;
    let worlds = window.WorldRepository.all();

    // Rehydrate images from IndexedDB before storing in state
    if (window.HOWLStorage) {
        try {
            worlds = await Promise.all(worlds.map(w => window.HOWLStorage.restoreWorldImages(w)));
        } catch (error) {
            console.warn('Failed to rehydrate world images on load:', error);
        }
    }

    window.worldState.worlds = worlds;
    await renderWorldGrid();
}
