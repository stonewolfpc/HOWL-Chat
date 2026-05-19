/**
 * World Page — Orchestrator
 *
 * Initializes the world page by loading data and wiring
 * keyboard shortcuts and outside-click dismissal.
 *
 * Depends on:
 *   world-state.js      — shared state + lorebook utilities
 *   world-form.js       — popup open/close/save
 *   world-dropdowns.js  — scenario/character/lorebook dropdowns
 *   world-grid.js       — grid render + edit/copy/delete/load
 *
 * @module world
 * @version 2.0.0
 */

'use strict';

document.addEventListener('DOMContentLoaded', function () {
    const onWorldPage = window.location.pathname.includes('world.html') ||
                        window.location.pathname.endsWith('/world') ||
                        !!document.querySelector('#world-creation-overlay');

    if (!onWorldPage) {
        console.log('Not on world page, skipping world.js initialization');
        return;
    }

    console.log('Initializing world page...');
    loadWorlds();
    loadAvailableScenarios();
    loadAvailableCharacters();
    loadAvailableLorebooks();

    // Keyboard shortcuts
    document.addEventListener('keydown', function (event) {
        if (event.key === 'Escape')             closeWorldCreationPopup();
        if (event.ctrlKey && event.key === 'n') { event.preventDefault(); openWorldCreationPopup(); }
    });

    // Close dropdowns when clicking outside
    document.addEventListener('click', function (event) {
        if (!event.target.closest('.world-dropdown-container')) {
            document.querySelectorAll('.world-dropdown-menu').forEach(d => d.classList.remove('show'));
        }
    });

    // Close popup when clicking outside
    const overlay = document.getElementById('world-creation-overlay');
    if (overlay) {
        overlay.addEventListener('click', function (event) {
            if (event.target === this) closeWorldCreationPopup();
        });
    }
});
