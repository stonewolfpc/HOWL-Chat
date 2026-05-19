/**
 * Scenario Form
 *
 * Popup open/close/reset, image select/remove/preview,
 * data collection, and save (via HOWLStorage).
 *
 * Depends on: scenario-state.js, scenario-dropdowns.js, scenario-grid.js
 *
 * @module scenario-form
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Popup Lifecycle
// ---------------------------------------------------------------------------

function openScenarioCreationPopup() {
    const overlay = document.getElementById('scenario-creation-overlay');
    overlay.style.display = 'flex';
    if (window.scenarioState.editingScenarioIndex < 0) resetScenarioCreationForm();
    loadAvailableCharacters();
}

function closeScenarioCreationPopup() {
    const overlay = document.getElementById('scenario-creation-overlay');
    overlay.style.display = 'none';
    window.scenarioState.editingScenarioIndex = -1;
    resetScenarioCreationForm();
}

function resetScenarioCreationForm() {
    const titleEl   = document.getElementById('scenario-title');
    const detailsEl = document.getElementById('scenario-details');
    if (titleEl)   titleEl.value   = '';
    if (detailsEl) detailsEl.value = '';
    window.scenarioState.scenarioImages = [];
    window.scenarioState.scenarioThumbnailIndex = 0;
    window.scenarioState.selectedCharacters   = [];
    window.scenarioState.selectedLorebooks    = [];
    updateScenarioImagesList();
    updateSelectedCharactersList();
    updateSelectedLorebooksList();
}

// ---------------------------------------------------------------------------
// Image Handling (Multi-Image Array with Thumbnail Selection)
// ---------------------------------------------------------------------------

function handleScenarioImageSelect(event) {
    const files = Array.from(event.target.files).filter(f => f.type.startsWith('image/'));
    event.target.value = '';
    if (files.length === 0) return;

    Promise.all(files.map(async (file) => {
        const result = await window.ImageProcessor.process(file);
        window.scenarioState.scenarioImages.push({
            id:          `img_${Date.now()}_${Math.random().toString(36).substr(2, 6)}`,
            name:        result.meta.name,
            type:        result.meta.type,
            size:        result.meta.size,
            data:        result.thumbnail,
            medium:      result.medium,
            original:    result.original,
            isThumbnail: window.scenarioState.scenarioImages.length === 0
        });
    })).then(() => updateScenarioImagesList())
       .catch(err => { console.error('Image processing failed:', err); alert('Failed to process image: ' + err.message); });
}

function updateScenarioImagesList() {
    const list = document.getElementById('scenario-images-list');
    if (!list) return;
    const ss = window.scenarioState;

    if (ss.scenarioImages.length === 0) {
        list.innerHTML = `
            <div class="no-images-placeholder">
                <span class="no-images-icon icon icon-image icon-lg"></span>
                <span class="no-images-text">No images added yet</span>
            </div>
        `;
        return;
    }

    list.innerHTML = ss.scenarioImages.map((image, index) => `
        <div class="world-image-item ${image.isThumbnail ? 'thumbnail' : ''}" onclick="event.stopPropagation(); setScenarioThumbnail(${index})">
            <img src="${image.data}" alt="${image.name}" class="world-image-item-thumb">
            <div class="world-image-item-info">
                <div class="world-image-item-name">${image.name}</div>
            </div>
            <div class="world-image-item-actions">
                ${image.isThumbnail ? '<span class="world-image-item-badge">Thumbnail</span>' : ''}
                <button class="world-image-item-remove" onclick="removeScenarioImage(${index}, event)" title="Remove"><span class="icon icon-trash"></span></button>
            </div>
        </div>
    `).join('');
}

function setScenarioThumbnail(index) {
    const ss = window.scenarioState;
    ss.scenarioImages.forEach((img, i) => { img.isThumbnail = i === index; });
    ss.scenarioThumbnailIndex = index;
    updateScenarioImagesList();
}

function removeScenarioImage(index, event) {
    event.stopPropagation();
    const ss = window.scenarioState;
    ss.scenarioImages.splice(index, 1);
    if (ss.scenarioImages.length > 0 && ss.scenarioThumbnailIndex >= ss.scenarioImages.length) {
        ss.scenarioThumbnailIndex = 0;
        ss.scenarioImages[0].isThumbnail = true;
    }
    updateScenarioImagesList();
}

// ---------------------------------------------------------------------------
// Data Collection
// ---------------------------------------------------------------------------

function collectScenarioData() {
    const ss = window.scenarioState;
    const baseData = {
        title:       document.getElementById('scenario-title').value,
        description: document.getElementById('scenario-details').value,
        images:      ss.scenarioImages.map(img => ({
            id: img.id, name: img.name, type: img.type, size: img.size,
            data: img.data, medium: img.medium, original: img.original
        })),
        thumbnailIndex: ss.scenarioThumbnailIndex,
        characters:  ss.selectedCharacters.map(c => ({
            id:     c.id,
            name:   c.name,
            role:   c.metadata?.role || '',
            avatar: c.avatar || ''
        })),
        metadata: {
            version: 1,
            created: new Date().toISOString(),
            updated: new Date().toISOString(),
            seed:    Math.floor(Math.random() * 1000000000)
        },
        howl: {
            worldId:     '',
            world:       '',
            lorebooks:   ss.selectedLorebooks.map(ss.lorebookToRef),
            notes:       '',
            image_style: 'fantasy-realistic'
        }
    };

    baseData.lorebooks = baseData.howl.lorebooks;

    if (ss.editingScenarioIndex >= 0) {
        const scenarios = JSON.parse(localStorage.getItem('howl-chat-scenarios') || '[]');
        const existing  = scenarios[ss.editingScenarioIndex];
        if (existing?.id) {
            baseData.id = existing.id;
            if (existing.metadata?.created) baseData.metadata.created = existing.metadata.created;
        } else {
            baseData.id = `scenario-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
        }
    } else {
        baseData.id = `scenario-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
    }

    return baseData;
}

// ---------------------------------------------------------------------------
// Save
// ---------------------------------------------------------------------------

async function saveScenario() {
    const scenarioData = collectScenarioData();
    if (!scenarioData.title) { alert('Scenario title is required!'); return; }

    try {
        if (!window.HOWLStorage) throw new Error('Storage manager not available');
        const ss = window.scenarioState;
        if (ss.editingScenarioIndex >= 0) {
            await window.HOWLStorage.updateScenario(ss.editingScenarioIndex, scenarioData);
        } else {
            await window.HOWLStorage.storeScenario(scenarioData);
        }
        closeScenarioCreationPopup();
        loadScenarios();
    } catch (error) {
        console.error('Failed to save scenario:', error);
        alert('Failed to save scenario: ' + error.message);
    }
}
