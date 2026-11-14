/**
 * Character Memories Viewer
 *
 * Displays and manages memories scoped to a specific character.
 * Integrated with sidebar character context menu.
 *
 * @module character-memories
 * @version 1.0.0
 * @requires memory-store.js, sidebar-state.js
 */

'use strict';

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

/**
 * Open character memories viewer modal
 * @param {string} characterId - Character ID to view memories for
 */
// ---------------------------------------------------------------------------
// Data Layer
// ---------------------------------------------------------------------------

function getCharacterMemoriesData(characterId) {
    const allMemories = window.MemoryStore ?
        window.MemoryStore.getMemoryEntries() : [];

    return allMemories.filter(m =>
        m.scope === 'character' &&
        m.ownerKind === 'character_id' &&
        m.ownerId === characterId
    );
}

// ---------------------------------------------------------------------------
// UI Rendering
// ---------------------------------------------------------------------------

function renderCharacterMemoriesModal(character, memories) {
    // Remove existing modal
    const existing = document.querySelector('.character-memories-modal');
    if (existing) existing.remove();

    const modal = document.createElement('div');
    modal.className = 'character-memories-modal modal-overlay';
    modal.style.cssText = `
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0,0,0,0.8);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 10000;
        backdrop-filter: blur(4px);
    `;

    const pinned = memories.filter(m => m.pinned);
    const active = memories.filter(m => !m.pinned && m.status === 'active');
    const archived = memories.filter(m => m.status === 'archived');

    modal.innerHTML = `
        <div class="modal-container" style="
            background: var(--shell-panel-bg);
            border-radius: 16px;
            width: min(720px, 95%);
            max-height: 85vh;
            display: flex;
            flex-direction: column;
            border: 1px solid rgba(255,255,255,0.12);
        ">
            <div class="modal-header" style="
                padding: 20px 24px;
                border-bottom: 1px solid rgba(255,255,255,0.08);
                display: flex;
                align-items: center;
                justify-content: space-between;
            ">
                <div>
                    <div style="font-size: 20px; font-weight: 700; color: #fff;">
                        ${character.name}'s Memories
                    </div>
                    <div style="font-size: 13px; color: var(--surface-hover); margin-top: 4px;">
                        ${memories.length} memory${memories.length !== 1 ? 'ies' : 'y'} •
                        ${pinned.length} pinned •
                        ${active.length} active
                    </div>
                </div>
                <button onclick="closeCharacterMemoriesModal()" style="
                    background: var(--surface-hover);
                    border: none;
                    color: #fff;
                    width: 36px;
                    height: 36px;
                    border-radius: 10px;
                    cursor: pointer;
                    font-size: 18px;
                ">✕</button>
            </div>

            <div class="modal-body" style="
                padding: 20px 24px;
                overflow-y: auto;
                flex: 1;
            ">
                ${renderMemorySection('Pinned Memories', pinned, character.id)}
                ${renderMemorySection('Active Memories', active, character.id)}
                ${archived.length > 0 ? renderMemorySection('Archived Memories', archived, character.id, true) : ''}

                ${memories.length === 0 ? renderEmptyState(character.name) : ''}
            </div>

            <div class="modal-footer" style="
                padding: 16px 24px;
                border-top: 1px solid rgba(255,255,255,0.08);
                display: flex;
                justify-content: flex-end;
                gap: 12px;
            ">
                <button onclick="closeCharacterMemoriesModal()" style="
                    padding: 10px 20px;
                    background: var(--surface-hover);
                    border: 1px solid rgba(255,255,255,0.14);
                    color: #fff;
                    border-radius: 8px;
                    cursor: pointer;
                ">Close</button>
            </div>
        </div>
    `;

    document.body.appendChild(modal);

    // Close on backdrop click
    modal.addEventListener('click', (e) => {
        if (e.target === modal) closeCharacterMemoriesModal();
    });
}

function renderMemorySection(title, memories, characterId, isArchived = false) {
    if (memories.length === 0) return '';

    const opacity = isArchived ? '0.6' : '1';

    return `
        <div style="margin-bottom: 24px; opacity: ${opacity};">
            <div style="
                font-size: 13px;
                font-weight: 600;
                color: rgba(255,255,255,0.5);
                text-transform: uppercase;
                letter-spacing: 0.5px;
                margin-bottom: 12px;
            ">${title} (${memories.length})</div>
            <div style="display: flex; flex-direction: column; gap: 8px;">
                ${memories.map(m => renderMemoryCard(m, characterId)).join('')}
            </div>
        </div>
    `;
}

function renderMemoryCard(memory, characterId) {
    const timeAgo = formatTimeAgo(memory.createdAt);
    const tagsHtml = (memory.tags || []).map(t =>
        `<span style="
            font-size: 11px;
            padding: 2px 8px;
            background: var(--surface-hover);
            border-radius: 4px;
            color: rgba(255,255,255,0.7);
        ">${t}</span>`
    ).join('');

    return `
        <div class="memory-card" style="
            background: var(--surface-hover);
            border: 1px solid rgba(255,255,255,0.06);
            border-radius: 10px;
            padding: 14px 16px;
            cursor: pointer;
            transition: all 0.2s ease;
        " onmouseenter="this.style.background='var(--surface-hover)'"
           onmouseleave="this.style.background='var(--surface-hover)'"
           onclick="editCharacterMemory('${memory.id}', '${characterId}')">
            <div style="display: flex; align-items: flex-start; gap: 12px;">
                <div style="font-size: 20px; line-height: 1;">
                    ${memory.pinned ? '📌' : memory.source === 'auto' ? '⚡' : '📝'}
                </div>
                <div style="flex: 1; min-width: 0;">
                    <div style="
                        font-size: 15px;
                        font-weight: 600;
                        color: #fff;
                        margin-bottom: 4px;
                        overflow: hidden;
                        text-overflow: ellipsis;
                        white-space: nowrap;
                    ">${memory.title || 'Untitled Memory'}</div>
                    <div style="
                        font-size: 13px;
                        color: rgba(255,255,255,0.6);
                        line-height: 1.4;
                        display: -webkit-box;
                        -webkit-line-clamp: 2;
                        -webkit-box-orient: vertical;
                        overflow: hidden;
                    ">${memory.summary || memory.content || ''}</div>
                    <div style="
                        display: flex;
                        align-items: center;
                        gap: 8px;
                        margin-top: 8px;
                        flex-wrap: wrap;
                    ">
                        <span style="font-size: 12px; color: var(--surface-hover);">
                            ${timeAgo}
                        </span>
                        ${memory.confidence ? `<span style="font-size: 12px; color: var(--surface-hover);">• ${memory.confidence}% confidence</span>` : ''}
                        ${tagsHtml}
                    </div>
                </div>
            </div>
        </div>
    `;
}

function renderEmptyState(characterName) {
    return `
        <div style="text-align: center; padding: 40px 20px; color: var(--surface-hover);">
            <div style="font-size: 48px; margin-bottom: 16px;">🧠</div>
            <div style="font-size: 16px; font-weight: 600; margin-bottom: 8px;">
                No memories for ${characterName}
            </div>
            <div style="font-size: 14px; line-height: 1.5;">
                Memories are created automatically when significant events occur,<br>
                or manually through Story Tools.
            </div>
        </div>
    `;
}

// ---------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------

function closeCharacterMemoriesModal() {
    const modal = document.querySelector('.character-memories-modal');
    if (modal) modal.remove();
}

function editCharacterMemory(memoryId, characterId) {
    // Navigate to Story Tools with this memory pre-selected
    closeCharacterMemoriesModal();

    // Store selection in session storage for Story Tools to pick up
    sessionStorage.setItem('story-tools-selected-memory', memoryId);
    sessionStorage.setItem('story-tools-filter-character', characterId);

    navigateTo('story-tools.html');
}

function formatTimeAgo(timestamp) {
    if (!timestamp) return 'Unknown';

    const date = new Date(timestamp);
    const now = new Date();
    const diffMs = now - date;
    const diffMins = Math.floor(diffMs / 60000);
    const diffHrs = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHrs < 24) return `${diffHrs}h ago`;
    if (diffDays < 30) return `${diffDays}d ago`;
    return date.toLocaleDateString();
}

// ---------------------------------------------------------------------------
// Exports
// ---------------------------------------------------------------------------

window.CharacterMemories = {
    view: viewCharacterMemories,
    close: closeCharacterMemoriesModal
};
