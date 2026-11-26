/**
 * Scenario Page — Orchestrator
 *
 * Initializes the scenario page: loads grid and lorebooks,
 * wires keyboard shortcuts and outside-click dismissal,
 * and exposes the public window.scenarioModule API.
 *
 * Depends on (must load first):
 *   scenario-state.js     — shared state + lorebook utilities
 *   scenario-dropdowns.js — character/lorebook dropdowns
 *   scenario-form.js      — popup lifecycle, image handling, save
 *   scenario-grid.js      — grid render + edit/delete
 *
 * @module scenario
 * @version 2.0.0
 */

'use strict';

document.addEventListener('DOMContentLoaded', async function () {
    const scenarioGrid    = document.getElementById('scenario-grid');
    const scenarioOverlay = document.getElementById('scenario-creation-overlay');

    if (scenarioGrid && scenarioOverlay) {
        try {
            await window.LorebookRepo.hydrate();
        } catch (error) {
            console.error('Failed to refresh the lorebook library:', error);
        }

        loadScenarios();
        loadAvailableLorebooks();

        document.addEventListener('keydown', function (event) {
            if (event.key === 'Escape' && scenarioOverlay.style.display === 'flex') closeScenarioCreationPopup();
            if (event.ctrlKey && event.key === 'n') { event.preventDefault(); openScenarioCreationPopup(); }
        });

        scenarioOverlay.addEventListener('click', function (event) {
            if (event.target === this) closeScenarioCreationPopup();
        });
    }

    document.addEventListener('click', function (event) {
        if (!event.target.closest('.world-dropdown-container')) {
            document.querySelectorAll('.world-dropdown-menu').forEach(d => d.classList.remove('show'));
        }
    });
});

window.scenarioModule = {
    openScenarioCreationPopup,
    closeScenarioCreationPopup,
    saveScenario,
    loadScenarios,
    loadAvailableCharacters,
    addCharacterToScenario,
    removeCharacterFromScenario
};
