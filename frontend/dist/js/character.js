/**
 * Character Page — Orchestrator
 *
 * Wires DOMContentLoaded initialization, keyboard shortcuts,
 * outside-click dismissal, and exposes window.characterModule.
 *
 * Depends on (must load first):
 *   character-state.js     — shared state + utilities
 *   character-images.js    — image select/remove/thumbnail
 *   character-lorebooks.js — lorebook dropdown + toggleDropdown
 *   character-form.js      — popup lifecycle, data collect, save
 *   character-grid.js      — grid render + edit/delete
 *
 * @module character
 * @version 2.0.0
 */

'use strict';

document.addEventListener('DOMContentLoaded', function () {
    setTimeout(async () => {
        try {
            await window.LorebookRepo.hydrate();
        } catch (error) {
            console.error('Failed to refresh the lorebook library:', error);
        }

        const load = loadCharacters();
        loadAvailableLorebooks();
        Promise.resolve(load).then(() => {
            const characterId = new URLSearchParams(window.location.search).get('id');
            if (!characterId) return;
            const characters = window.characterState?.characterCache || [];
            const index = characters.findIndex(character => character.id === characterId);
            if (index >= 0) {
                editCharacter(index);
            } else if (typeof showNotification === 'function') {
                showNotification('The requested character was not found.', 'error');
            }
        });
    }, 100);

    document.addEventListener('keydown', function (event) {
        if (event.key === 'Escape') {
            closeCharacterCreationPopup();
            closeUserDetailsPopup();
        }
        if (event.ctrlKey && event.key === 'n') { event.preventDefault(); openCharacterCreationPopup(); }
        if (event.ctrlKey && event.key === 'u') { event.preventDefault(); openUserDetailsPopup(); }
    });

    document.addEventListener('click', function (event) {
        const charOverlay = document.getElementById('character-creation-overlay');
        const userOverlay = document.getElementById('user-details-overlay');
        if (event.target === charOverlay) closeCharacterCreationPopup();
        if (event.target === userOverlay) closeUserDetailsPopup();

        if (!event.target.closest('.world-dropdown-container')) {
            document.querySelectorAll('.world-dropdown-menu').forEach(m => m.classList.remove('show'));
            document.querySelectorAll('.world-dropdown-container').forEach(c => c.classList.remove('open'));
        }
    });
});

window.characterModule = {
    openCharacterCreationPopup,
    closeCharacterCreationPopup,
    openUserDetailsPopup,
    closeUserDetailsPopup,
    saveCharacter,
    saveUserDetails,
    loadUserDetails,
    loadCharacters,
    loadAvailableLorebooks
};
