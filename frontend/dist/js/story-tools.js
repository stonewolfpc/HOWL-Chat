/**
 * Story Tools — Orchestrator
 *
 * Initializes the page, wires the flyout toggle, tab switching,
 * search/filter handlers, and exposes window.storyTools.
 *
 * Depends on (must load first):
 *   story-tools-state.js — shared state + utilities
 *   story-tools-data.js  — getMemories, save, delete, export
 *   story-tools-ui.js    — render, editor, notifications
 *
 * @module story-tools
 * @version 2.0.0
 */

'use strict';

document.addEventListener('DOMContentLoaded', function () {
    initializeStoryTools();
});

async function initializeStoryTools() {
    const st = window.storyToolsState;
    st.currentChatId = st.getCurrentChatId();
    if (!st.currentChatId) { showError('No active chat session found'); return; }
    await loadRecentMemories();
    setupEventListeners();
}

function setupEventListeners() {
    const searchInput  = document.getElementById('memories-search');
    const filterSelect = document.getElementById('memories-filter-scope');
    if (searchInput)  searchInput.addEventListener('input', window.storyToolsState.debounce(handleSearch, 300));
    if (filterSelect) filterSelect.addEventListener('change', handleFilter);
}

// ---------------------------------------------------------------------------
// Load Helpers
// ---------------------------------------------------------------------------

async function loadRecentMemories() {
    try {
        showLoading(true);
        const memories = await getMemories({ limit: 5, sort: 'newest' });
        renderRecentList(memories);
        showLoading(false);
    } catch (error) {
        console.error('Failed to load recent memories:', error);
        showError('Failed to load memories');
        showLoading(false);
    }
}

async function loadFullMemoriesList() {
    try {
        showLoading(true);
        const memories = await getMemories({});
        renderFullList(memories);
        showLoading(false);
    } catch (error) {
        console.error('Failed to load full memories list:', error);
        showError('Failed to load memories');
        showLoading(false);
    }
}

// ---------------------------------------------------------------------------
// Flyout Toggle
// ---------------------------------------------------------------------------

async function toggleMemoriesFlyout() {
    const flyout    = document.getElementById('memories-flyout');
    const overlay   = document.getElementById('memories-flyout-overlay');
    const expandBtn = document.getElementById('memories-expand');
    if (!flyout || !overlay) return;

    const isOpen = flyout.classList.contains('show');
    if (isOpen) {
        flyout.classList.remove('show');
        overlay.classList.remove('show');
        if (expandBtn) expandBtn.innerHTML = `<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 12 5 6 5"></polyline><polyline points="9 6 21 6"></polyline></svg>`;
    } else {
        flyout.classList.add('show');
        overlay.classList.add('show');
        if (expandBtn) expandBtn.innerHTML = `<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="7 10 12 17 10 17"></polyline><polyline points="7 6 21 6"></polyline></svg>`;
        await loadFullMemoriesList();
    }
}

// ---------------------------------------------------------------------------
// Tab & Filter
// ---------------------------------------------------------------------------

function switchMemoriesTab(tab) {
    window.storyToolsState.currentTab = tab;
    document.querySelectorAll('.memories-tab').forEach(t => t.classList.remove('active'));
    document.getElementById(`memories-tab-${tab}`)?.classList.add('active');
    loadFullMemoriesList();
}

async function handleSearch(event) {
    const filter = { query: event.target.value.trim() || undefined, scope: document.getElementById('memories-filter-scope').value };
    try {
        showLoading(true);
        renderFullList(await getMemories(filter));
        showLoading(false);
    } catch { showError('Search failed'); showLoading(false); }
}

async function handleFilter(event) {
    const scope  = event.target.value;
    const filter = { scope: scope === 'all' ? undefined : scope, query: document.getElementById('memories-search').value.trim() || undefined };
    try {
        showLoading(true);
        renderFullList(await getMemories(filter));
        showLoading(false);
    } catch { showError('Filter failed'); showLoading(false); }
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

window.storyTools = {
    toggleMemoriesFlyout,
    switchMemoriesTab,
    saveMemory,
    deleteMemory,
    exportMemories,
    loadRecentMemories,
    selectMemory
};
