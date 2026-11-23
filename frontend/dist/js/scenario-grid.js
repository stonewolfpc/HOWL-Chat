/**
 * Scenario Grid
 *
 * Renders the scenario card grid and handles edit/delete actions,
 * including image restoration from HOWLStorage.
 *
 * Depends on: scenario-state.js, scenario-form.js, scenario-dropdowns.js
 *
 * @module scenario-grid
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Grid Render
// ---------------------------------------------------------------------------

async function loadScenarios() {
    const grid = document.getElementById('scenario-grid');
    if (!grid) return;

    if (window.ScenarioRepository?.ready) await window.ScenarioRepository.ready;
    let scenarios = window.ScenarioRepository.all();

    // Rehydrate images from IndexedDB before rendering
    if (window.HOWLStorage) {
        try {
            scenarios = await Promise.all(scenarios.map(s => window.HOWLStorage.restoreScenarioImages(s)));
        } catch (error) {
            console.warn('Failed to rehydrate scenario images:', error);
        }
    }

    if (scenarios.length === 0) {
        grid.innerHTML = `
            <div class="world-grid-empty">
                <div class="world-grid-empty-icon">📜</div>
                <div class="world-grid-empty-text">No scenarios created yet</div>
                <div class="world-grid-empty-subtext">Click "Create New Scenario" to get started</div>
            </div>
        `;
        return;
    }

    grid.innerHTML = scenarios.map((scenario, index) => {
        let imageSrc = '';
        if (scenario.images && Array.isArray(scenario.images) && scenario.images.length > 0) {
            const thumbIdx = scenario.thumbnailIndex || 0;
            const thumb = scenario.images[thumbIdx] || scenario.images[0];
            imageSrc = thumb.data || thumb.thumbnail || '';
        } else if (scenario.image) {
            imageSrc = typeof scenario.image === 'string'
                ? scenario.image
                : (scenario.image.data || scenario.image.thumbnail || '');
        }
        const hasImage = !!imageSrc;
        return `
            <div class="world-bubble" data-scenario-id="${scenario.id || index}">
                ${hasImage
                    ? `<img class="world-bubble-image" src="${imageSrc}" alt="${scenario.title}" onerror="this.style.display='none';this.nextElementSibling.style.display='flex';">
                       <div class="world-bubble-placeholder" style="display:none;"><div class="world-bubble-placeholder-icon icon icon-scenario">📜</div></div>`
                    : `<div class="world-bubble-placeholder"><div class="world-bubble-placeholder-icon icon icon-scenario">📜</div></div>`
                }
                <div class="world-bubble-name">${scenario.title}</div>
                <div class="world-bubble-meta">${scenario.description || ''}</div>
                <div class="world-bubble-actions">
                    <button class="world-bubble-action" onclick="editScenario(${index}, event)" title="Edit"><span class="icon icon-edit"></span></button>
                    <button class="world-bubble-action" onclick="deleteScenario(${index}, event)" title="Delete"><span class="icon icon-trash"></span></button>
                </div>
            </div>
        `;
    }).join('');
}

// ---------------------------------------------------------------------------
// Edit
// ---------------------------------------------------------------------------

function editScenarioWrapper(index) { editScenario(index, null); }

async function editScenario(index, event) {
    if (event) event.stopPropagation();

    let scenarios = window.ScenarioRepository.all();
    let scenario  = scenarios[index];

    // Rehydrate images from IndexedDB before editing
    if (window.HOWLStorage && scenario) {
        try {
            scenario = await window.HOWLStorage.restoreScenarioImages(scenario);
        } catch (error) {
            console.warn('Failed to rehydrate scenario images for edit:', error);
        }
    }

    const ss        = window.scenarioState;

    ss.editingScenarioIndex = index;

    document.getElementById('scenario-title').value   = scenario.title       || '';
    document.getElementById('scenario-details').value = scenario.description || '';

    ss.scenarioImages = [];
    ss.scenarioThumbnailIndex = 0;
    if (scenario.images && Array.isArray(scenario.images) && scenario.images.length > 0) {
        ss.scenarioImages = scenario.images.map((img, i) => ({
            id: img.id || `img_${Date.now()}_${Math.random().toString(36).substr(2, 6)}`,
            name: img.name || 'Image',
            type: img.type,
            size: img.size,
            data: img.data || img.thumbnail || '',
            medium: img.medium,
            isThumbnail: i === (scenario.thumbnailIndex || 0)
        }));
        ss.scenarioThumbnailIndex = scenario.thumbnailIndex || 0;
    } else if (scenario.image) {
        const imgData = typeof scenario.image === 'string'
            ? scenario.image
            : (scenario.image.data || scenario.image.thumbnail || '');
        if (imgData) {
            ss.scenarioImages = [{
                id: `img_${Date.now()}_${Math.random().toString(36).substr(2, 6)}`,
                name: scenario.image.name || 'Image',
                type: scenario.image.type,
                size: scenario.image.size,
                data: imgData,
                medium: scenario.image.medium,
                isThumbnail: true
            }];
            ss.scenarioThumbnailIndex = 0;
        }
    }
    updateScenarioImagesList();

    // Re-hydrate from the full character catalog — stored scenarios have avatars stripped
    const allChars = (window.CharacterStorage?.all?.() || []);
    ss.selectedCharacters = (scenario.characters || []).map(sc => {
        const full = allChars.find(c => c.name === sc.name || c.id === sc.id);
        return full || sc; // fall back to stored (name+role only) if not in catalog
    });
    updateSelectedCharactersList();
    loadAvailableCharacters();

    ss.selectedLorebooks = ss.resolveLorebooksFromRefs(ss.getScenarioLorebookRefs(scenario));
    updateSelectedLorebooksList();
    loadAvailableLorebooks();

    openScenarioCreationPopup();
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

function deleteScenarioWrapper(index) { deleteScenario(index, null); }

async function deleteScenario(index, event) {
    if (event) event.stopPropagation();

    const scenarios    = window.ScenarioRepository.all();
    const scenario     = scenarios[index];
    const scenarioName = scenario?.title || 'Unnamed scenario';
    const scenarioId   = scenario?.id;

    window.DeleteConfirm.show({
        title: `Delete "${scenarioName}"?`,
        body:  'The scenario, its image, and all world/chat references will be removed. Lorebooks are kept.',
        onConfirm: async () => {
            if (window.DeleteService && scenarioId) await window.DeleteService.detachScenarioFromAll(scenarioId);
            try {
                if (!scenarioId) throw new Error('Scenario identity is missing');
                await window.ScenarioRepository.remove(scenarioId);
                loadScenarios();
            } catch (error) {
                console.error('Failed to delete scenario:', error);
            }
        }
    });
}
