/**
 * Character Images
 *
 * Image selection, removal, thumbnail management, and list rendering
 * for the character creation/edit form.
 *
 * Depends on: character-state.js
 *
 * @module character-images
 * @version 1.0.0
 */

'use strict';

function handleCharacterImageSelect(event) {
    const cs    = window.characterState;
    const files = Array.from(event.target.files).filter(f => f.type.startsWith('image/'));
    event.target.value = '';
    if (files.length === 0) return;

    Promise.all(files.map(async (file) => {
        const result = await window.ImageProcessor.process(file);
        cs.characterImages.push({
            id:          `img_${Date.now()}_${Math.random().toString(36).substr(2, 6)}`,
            name:        result.meta.name,
            size:        cs.formatFileSize(result.meta.size),
            data:        result.thumbnail,
            medium:      result.medium,
            original:    result.original,
            isThumbnail: cs.characterImages.length === 0
        });
    })).then(() => updateImagesList())
       .catch(err => { console.error('Image processing failed:', err); alert('Failed to process image: ' + err.message); });
}

function updateImagesList() {
    const list = document.getElementById('character-images-list');
    if (!list) return;
    const cs = window.characterState;

    if (cs.characterImages.length === 0) {
        list.innerHTML = `
            <div class="no-images-placeholder">
                <span class="no-images-icon">🖼️</span>
                <span class="no-images-text">No images added yet</span>
            </div>
        `;
        return;
    }

    list.innerHTML = cs.characterImages.map((image, index) => `
        <div class="character-image-item ${image.isThumbnail ? 'thumbnail' : ''}" onclick="setThumbnail(${index})">
            <img src="${image.data}" alt="${image.name}" class="character-image-thumbnail">
            <div class="character-image-info">
                <div class="character-image-name">${image.name}</div>
                <div class="character-image-size">${image.size}</div>
            </div>
            <div class="character-image-actions">
                ${image.isThumbnail ? '<button class="thumbnail-indicator">Thumbnail</button>' : ''}
                <button class="remove-image-button" onclick="removeImage(${index}, event)">Remove</button>
            </div>
        </div>
    `).join('');
}

function setThumbnail(index) {
    const cs = window.characterState;
    cs.characterImages.forEach((img, i) => { img.isThumbnail = i === index; });
    cs.thumbnailIndex = index;
    updateImagesList();
}

function removeImage(index, event) {
    event.stopPropagation();
    const cs = window.characterState;
    cs.characterImages.splice(index, 1);
    if (cs.characterImages.length > 0 && cs.thumbnailIndex >= cs.characterImages.length) {
        cs.thumbnailIndex = 0;
        cs.characterImages[0].isThumbnail = true;
    }
    updateImagesList();
}
