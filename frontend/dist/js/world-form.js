/**
 * World Creation Form
 *
 * Manages the world creation/edit popup: open, close, clear,
 * image select/remove, and saving world data to localStorage.
 *
 * Depends on: world-state.js, world-grid.js, world-dropdowns.js
 *
 * @module world-form
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Popup Open / Close / Clear
// ---------------------------------------------------------------------------

async function openWorldCreationPopup(worldId = null) {
    const overlay = document.getElementById('world-creation-overlay');
    overlay.classList.add('show');

    if (worldId) {
        window.worldState.editingWorldId = worldId;
        let world = window.worldState.worlds.find(w => w.id === worldId);

        // Rehydrate images from IndexedDB before editing
        if (window.HOWLStorage && world) {
            try {
                world = await window.HOWLStorage.restoreWorldImages(world);
            } catch (error) {
                console.warn('Failed to rehydrate world images for edit:', error);
            }
        }

        if (world) {
            document.getElementById('world-title').value   = world.title;
            document.getElementById('world-details').value = world.details;

            window.worldState.worldImages = [];
            window.worldState.worldThumbnailIndex = 0;
            if (world.images && Array.isArray(world.images) && world.images.length > 0) {
                window.worldState.worldImages = world.images.map((img, i) => ({
                    id: img.id || `img_${Date.now()}_${Math.random().toString(36).substr(2, 6)}`,
                    name: img.name || 'Image',
                    type: img.type,
                    size: img.size,
                    data: img.data || img.thumbnail || '',
                    medium: img.medium,
                    isThumbnail: i === (world.thumbnailIndex || 0)
                }));
                window.worldState.worldThumbnailIndex = world.thumbnailIndex || 0;
            } else if (world.image) {
                const imgData = typeof world.image === 'string'
                    ? world.image
                    : (world.image.data || world.image.thumbnail || '');
                if (imgData) {
                    window.worldState.worldImages = [{
                        id: `img_${Date.now()}_${Math.random().toString(36).substr(2, 6)}`,
                        name: world.image.name || 'Image',
                        type: world.image.type,
                        size: world.image.size,
                        data: imgData,
                        medium: world.image.medium,
                        isThumbnail: true
                    }];
                    window.worldState.worldThumbnailIndex = 0;
                }
            }
            updateWorldImagesList();

            window.worldState.selectedScenarios  = world.scenarios || [];
            window.worldState.selectedCharacters = world.characters || [];
            window.worldState.selectedLorebooks  = window.worldState.resolveLorebooksFromRefs(
                world.howl?.lorebooks || world.lorebooks || []
            );

            updateSelectedScenariosDisplay();
            updateSelectedCharactersDisplay();
            updateSelectedLorebooksDisplay();
        }
    } else {
        window.worldState.editingWorldId = null;
        clearWorldForm();
    }
}

function closeWorldCreationPopup() {
    document.getElementById('world-creation-overlay').classList.remove('show');
    clearWorldForm();
}

function clearWorldForm() {
    document.getElementById('world-title').value               = '';
    document.getElementById('world-details').value             = '';
    window.worldState.worldImages = [];
    window.worldState.worldThumbnailIndex = 0;
    updateWorldImagesList();
    window.worldState.selectedScenarios  = [];
    window.worldState.selectedCharacters = [];
    window.worldState.selectedLorebooks  = [];
    updateSelectedScenariosDisplay();
    updateSelectedCharactersDisplay();
    updateSelectedLorebooksDisplay();
}

// ---------------------------------------------------------------------------
// Image Handling (Multi-Image Array with Thumbnail Selection)
// ---------------------------------------------------------------------------

function handleWorldImageSelect(event) {
    const files = Array.from(event.target.files).filter(f => f.type.startsWith('image/'));
    event.target.value = '';
    if (files.length === 0) return;

    Promise.all(files.map(async (file) => {
        const result = await window.ImageProcessor.process(file);
        window.worldState.worldImages.push({
            id:          `img_${Date.now()}_${Math.random().toString(36).substr(2, 6)}`,
            name:        result.meta.name,
            type:        result.meta.type,
            size:        result.meta.size,
            data:        result.thumbnail,
            medium:      result.medium,
            original:    result.original,
            isThumbnail: window.worldState.worldImages.length === 0
        });
    })).then(() => updateWorldImagesList())
       .catch(err => { console.error('Image processing failed:', err); alert('Failed to process image: ' + err.message); });
}

function updateWorldImagesList() {
    const list = document.getElementById('world-images-list');
    if (!list) return;
    const ws = window.worldState;

    if (ws.worldImages.length === 0) {
        list.innerHTML = `
            <div class="no-images-placeholder">
                <span class="no-images-icon icon icon-image icon-lg"></span>
                <span class="no-images-text">No images added yet</span>
            </div>
        `;
        return;
    }

    list.innerHTML = ws.worldImages.map((image, index) => `
        <div class="world-image-item ${image.isThumbnail ? 'thumbnail' : ''}" onclick="event.stopPropagation(); setWorldThumbnail(${index})">
            <img src="${image.data}" alt="${image.name}" class="world-image-item-thumb">
            <div class="world-image-item-info">
                <div class="world-image-item-name">${image.name}</div>
            </div>
            <div class="world-image-item-actions">
                ${image.isThumbnail ? '<span class="world-image-item-badge">Thumbnail</span>' : ''}
                <button class="world-image-item-remove" onclick="removeWorldImage(${index}, event)" title="Remove"><span class="icon icon-trash"></span></button>
            </div>
        </div>
    `).join('');
}

function setWorldThumbnail(index) {
    const ws = window.worldState;
    ws.worldImages.forEach((img, i) => { img.isThumbnail = i === index; });
    ws.worldThumbnailIndex = index;
    updateWorldImagesList();
}

function removeWorldImage(index, event) {
    event.stopPropagation();
    const ws = window.worldState;
    ws.worldImages.splice(index, 1);
    if (ws.worldImages.length > 0 && ws.worldThumbnailIndex >= ws.worldImages.length) {
        ws.worldThumbnailIndex = 0;
        ws.worldImages[0].isThumbnail = true;
    }
    updateWorldImagesList();
}

// ---------------------------------------------------------------------------
// Save
// ---------------------------------------------------------------------------

async function saveWorld() {
    const title   = document.getElementById('world-title').value.trim();
    const details = document.getElementById('world-details').value.trim();

    if (!title) { alert('Please enter a world title'); return; }

    const ws = window.worldState;
    const images = ws.worldImages.map(img => ({
        id: img.id,
        name: img.name,
        type: img.type,
        size: img.size,
        data: img.data,
        medium: img.medium,
        original: img.original
    }));
    const worldData = {
        id:       ws.editingWorldId || Date.now().toString(),
        title,
        details,
        images,
        thumbnailIndex: ws.worldThumbnailIndex,
        scenarios: ws.selectedScenarios.map(s => ({
            id: s.id, title: s.title || s.name, name: s.name || s.title, description: s.description
        })),
        characters: ws.selectedCharacters.map(c => ({
            id: c.id, name: c.name, description: c.description
        })),
        createdAt: ws.editingWorldId
            ? ws.worlds.find(w => w.id === ws.editingWorldId)?.createdAt
            : new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        howl: {
            lorebooks: ws.selectedLorebooks.map(l => ({
                id:          l.id || ws.lorebookKey(l),
                name:        ws.lorebookDisplayName(l),
                title:       ws.lorebookDisplayName(l),
                description: l.description
            }))
        }
    };

    try {
        if (!window.HOWLStorage) throw new Error('Storage manager not available');
        const saved = await window.HOWLStorage.storeWorld(worldData);
        const lsWorlds = JSON.parse(localStorage.getItem('howl-worlds') || '[]');
        const idx = lsWorlds.findIndex(w => w.id === saved.id);
        if (idx >= 0) ws.worlds[idx] = lsWorlds[idx]; else ws.worlds.push(lsWorlds[lsWorlds.length - 1]);
        ws.worlds = JSON.parse(localStorage.getItem('howl-worlds') || '[]');
    } catch (err) {
        console.error('Failed to save world:', err);
        alert('Failed to save world: ' + err.message);
        return;
    }

    ws.worldImages = [];
    ws.worldThumbnailIndex = 0;
    renderWorldGrid();
    closeWorldCreationPopup();

    console.log('World saved:', worldData.title, {
        scenarios:  worldData.scenarios.length,
        characters: worldData.characters.length,
        lorebooks:  worldData.howl.lorebooks.length
    });

    alert(`World "${worldData.title}" saved successfully!\n\nAttached:\n- ${worldData.scenarios.length} scenarios\n- ${worldData.characters.length} characters\n- ${worldData.howl.lorebooks.length} lorebooks`);

    // Creation started from the chat's world picker: that world is meant for
    // the chat. navigateTo() is a full page load, so the picker cannot be told
    // by postMessage — mark the world and return, and the chat attaches it on
    // load. Without this the created world never reaches the chat, and the
    // World panel keeps showing "No World Selected".
    if (new URLSearchParams(window.location.search).get('from') === 'chat') {
        try {
            sessionStorage.setItem('howl-pending-chat-world', worldData.id);
        } catch (err) {
            console.warn('Could not record pending world:', err);
        }
        window.location.href = 'index.html';
        return;
    }

    if (window.parent && window.parent !== window) {
        window.parent.postMessage({ type: 'world-created', world: worldData }, '*');
    }
}
