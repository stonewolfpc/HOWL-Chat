/**
 * Character Form
 *
 * Popup open/close for character creation and user details,
 * data collection, saveCharacter, saveUserDetails, loadUserDetails.
 *
 * Depends on: character-state.js, character-images.js,
 *             character-lorebooks.js, character-grid.js
 *
 * @module character-form
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Character Popup Lifecycle
// ---------------------------------------------------------------------------

function openCharacterCreationPopup() {
    document.getElementById('character-creation-overlay').style.display = 'flex';
    loadAvailableCharacters();
}

function closeCharacterCreationPopup() {
    document.getElementById('character-creation-overlay').style.display = 'none';
    const cs = window.characterState;
    cs.currentEditingIndex  = -1;
    cs.characterImages      = [];
    cs.selectedLorebooks.length = 0;
    const form = document.querySelector('.character-creation-form');
    if (form) form.reset();
}

// ---------------------------------------------------------------------------
// User Details Popup Lifecycle
// ---------------------------------------------------------------------------

function openUserDetailsPopup() {
    document.getElementById('user-details-overlay').style.display = 'flex';
    loadUserDetails();
}

function closeUserDetailsPopup() {
    document.getElementById('user-details-overlay').style.display = 'none';
}

// ---------------------------------------------------------------------------
// Data Collection
// ---------------------------------------------------------------------------

function collectCharacterData() {
    const cs        = window.characterState;
    const thumbnail = cs.characterImages.find(img => img.isThumbnail);
    const voiceEnabled = document.getElementById('char-voice-enabled')?.checked || false;
    return {
        id:          `char_${Date.now()}_${Math.floor(Math.random() * 1000)}`,
        name:        document.getElementById('character-name')?.value              || '',
        description: document.getElementById('character-description')?.value       || '',
        personality: document.getElementById('character-personality')?.value       || '',
        scenario:    document.getElementById('character-scenario')?.value          || '',
        first_mes:   document.getElementById('character-intro')?.value             || '',
        mes_example: document.getElementById('character-dialogue-examples')?.value || '',
        avatar:      thumbnail ? thumbnail.data : '',
        tags:        [],
        metadata: {
            version:    1,
            created:    new Date().toISOString(),
            updated:    new Date().toISOString(),
            seed:       Math.floor(Math.random() * 1000000000),
            age:        document.getElementById('character-age')?.value        || '',
            role:       document.getElementById('character-role')?.value       || '',
            background: document.getElementById('character-background')?.value || '',
            priority:   Math.min(5, Math.max(1, parseInt(document.getElementById('character-priority')?.value || '3', 10)))
        },
        howl: {
            world:     '',
            lorebooks: cs.selectedLorebooks.map(l => ({
                id:          l.id || cs.lorebookKey(l),
                name:        cs.lorebookDisplayName(l),
                title:       cs.lorebookDisplayName(l),
                description: l.description || ''
            })),
            memories:    [],
            notes:       '',
            image_style: 'fantasy-realistic',
            images:      cs.characterImages.map(img => ({
                id: img.id, name: img.name, size: img.size,
                data: img.data, isThumbnail: img.isThumbnail
            })),
            voice: voiceEnabled ? {
                enabled: true,
                engine:  document.getElementById('char-voice-engine')?.value  || '',
                voice:   document.getElementById('char-voice-model')?.value  || '',
                speed:   parseFloat(document.getElementById('char-voice-speed')?.value  || '1.0'),
                volume:  parseFloat(document.getElementById('char-voice-volume')?.value || '0.8')
            } : { enabled: false }
        }
    };
}

// ---------------------------------------------------------------------------
// Save Character
// ---------------------------------------------------------------------------

function saveCharacter() {
    const characterData = collectCharacterData();
    if (!characterData.name) { alert('Character name is required!'); return; }

    const cs       = window.characterState;
    const isEdit   = cs.currentEditingIndex >= 0;
    const existing = isEdit
        ? (cs.characterCache || [])[cs.currentEditingIndex]
        : null;

    // Preserve the original id when editing so IndexedDB images stay linked
    if (isEdit && existing?.id) characterData.id = existing.id;

    window.CharacterStorage.save(characterData)
        .then(() => {
            if (window.parent && window.parent !== window) {
                const msgType = isEdit ? 'character-updated' : 'character-created';
                const payload = isEdit
                    ? { type: msgType, character: characterData, index: cs.currentEditingIndex }
                    : { type: msgType, character: characterData };
                window.parent.postMessage(payload, '*');
            }
            closeCharacterCreationPopup();
            cs.currentEditingIndex = -1;
            loadCharacters();
        })
        .catch(err => {
            console.error('Failed to save character:', err);
            alert('Failed to save character. Please try again.');
        });
}

// ---------------------------------------------------------------------------
// User Details
// ---------------------------------------------------------------------------

function saveUserDetails() {
    const details = {
        name:        document.getElementById('user-name').value,
        description: document.getElementById('user-description').value,
        notes:       document.getElementById('user-notes').value
    };
    localStorage.setItem('howl-chat-user-details', JSON.stringify(details));
    closeUserDetailsPopup();
}

function loadUserDetails() {
    const saved = localStorage.getItem('howl-chat-user-details');
    if (!saved) return;
    const details = JSON.parse(saved);
    document.getElementById('user-name').value        = details.name        || '';
    document.getElementById('user-description').value = details.description || '';
    document.getElementById('user-notes').value       = details.notes       || '';
}

// ---------------------------------------------------------------------------
// Voice Settings UI Wiring
// ---------------------------------------------------------------------------

function _updateVoiceControlsVisible() {
    const enabled = document.getElementById('char-voice-enabled')?.checked || false;
    const controls = document.getElementById('char-voice-controls');
    if (controls) controls.style.display = enabled ? 'flex' : 'none';
}

function _updateSliderDisplay(id, valId, suffix) {
    const slider = document.getElementById(id);
    const valSpan = document.getElementById(valId);
    if (!slider || !valSpan) return;
    valSpan.textContent = slider.value + (suffix || '');
}

document.addEventListener('DOMContentLoaded', () => {
    const voiceToggle = document.getElementById('char-voice-enabled');
    if (voiceToggle) {
        voiceToggle.addEventListener('change', _updateVoiceControlsVisible);
        _updateVoiceControlsVisible();
    }
    const speedSlider = document.getElementById('char-voice-speed');
    if (speedSlider) {
        speedSlider.addEventListener('input', () => _updateSliderDisplay('char-voice-speed', 'char-voice-speed-val', 'x'));
    }
    const volSlider = document.getElementById('char-voice-volume');
    if (volSlider) {
        volSlider.addEventListener('input', () => _updateSliderDisplay('char-voice-volume', 'char-voice-volume-val', ''));
    }
});
