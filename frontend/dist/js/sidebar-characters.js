/**
 * Sidebar Character Panel
 *
 * Renders and manages the left-sidebar character list including:
 * - Character list rendering and item creation
 * - Character selection modal (add to chat)
 * - Context menu (step-out, edit, replace image, memories, remove)
 * - Step-out / death dialog
 * - Character image replacement
 *
 * Depends on: sidebar-state.js, sidebar-session.js
 *
 * @module sidebar-characters
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Panel Render
// ---------------------------------------------------------------------------



// ---------------------------------------------------------------------------
// Character Item Creation
// ---------------------------------------------------------------------------

async function createCharacterItem(character) {
    const item = document.createElement('div');
    item.className = 'character-item';
    item.dataset.characterId = character.id;

    const statusClass = character.status || window.sidebarState.CharacterStatus.OFFLINE;
    const statusText  = statusClass.charAt(0).toUpperCase() + statusClass.slice(1);

    // Resolve avatar: _resolvedImage first, then avatar/image field, then CharacterStorage
    let avatarSrc = character._resolvedImage
                  || window.sidebarState.resolveImageSrc(character.avatar || character.image);

    if (!avatarSrc && character.id && window.CharacterStorage) {
        try {
            const images = await window.CharacterStorage.loadImages(character.id);
            const thumb = images.find(img => img.isThumbnail) || images[0];
            if (thumb && thumb.data) {
                avatarSrc = thumb.data;
                character._resolvedImage = thumb.data;
            }
        } catch (_) {}
    }

    if (!avatarSrc && character.howl?.images?.length > 0) {
        const thumb = character.howl.images.find(img => img.isThumbnail) || character.howl.images[0];
        if (thumb && thumb.data) {
            avatarSrc = thumb.data;
            character._resolvedImage = thumb.data;
        }
    }

    const role = character.role || character.metadata?.role || '';

    item.innerHTML = `
        <div class="character-avatar" onclick="openCharacterSheet('${character.id}')">
            ${avatarSrc
                ? `<img src="${avatarSrc}" alt="${character.name}" onerror="this.style.display='none'; this.nextElementSibling.style.display='flex';">
                   <div class="avatar-placeholder" style="display:none;">👤</div>`
                : `<div class="avatar-placeholder">👤</div>`
            }
            <div class="status-indicator ${statusClass}"></div>
        </div>
        <div class="character-info" onclick="openCharacterSheet('${character.id}')">
            <div class="character-name">${character.name}</div>
            ${role ? `<div class="character-role">${role}</div>` : ''}
        </div>
        <div class="character-actions">
            <button class="character-menu-btn" onclick="showCharacterMenu(event, '${character.id}')">⋮</button>
        </div>
    `;

    if (character.status === 'stepped-out' || character.status === 'dead') {
        item.classList.add('character-inactive');
    }

    item.addEventListener('contextmenu', (e) => {
        e.preventDefault();
        showCharacterMenu(e, character.id);
    });

    return item;
}

// ---------------------------------------------------------------------------
// Character Sheet / Edit
// ---------------------------------------------------------------------------

function openCharacterSheet(characterId) {
    console.log('Opening character sheet for:', characterId);
    if (window.go && window.go.main && window.go.main.App) {
        window.go.main.App.OpenCharacterSheet(characterId)
            .then(() => console.log('Character sheet opened'))
            .catch(err => { console.error('Failed to open character sheet:', err); navigateTo('character.html?id=' + characterId); });
    } else {
        navigateTo('character.html?id=' + characterId);
    }
}

function editCharacter(characterId) {
    console.log('Editing character:', characterId);
    if (window.go && window.go.main && window.go.main.App && window.go.main.App.OpenCharacterSheet) {
        window.go.main.App.OpenCharacterSheet(characterId)
            .then(() => console.log('Character sheet opened:', characterId))
            .catch(err => { console.error('Failed to open character sheet:', err); navigateTo('character.html?id=' + characterId); });
    } else {
        navigateTo('character.html?id=' + characterId);
    }
}

// ---------------------------------------------------------------------------
// Context Menu
// ---------------------------------------------------------------------------

function showCharacterMenu(event, characterId) {
    event.preventDefault();
    event.stopPropagation();

    const existingMenu = document.querySelector('.character-context-menu');
    if (existingMenu) existingMenu.remove();

    const cd = window.sidebarState.currentChatData;
    const character = cd.characters.find(c => c.id === characterId);
    const isDead = character && character.status === 'dead';
    const isSteppedOut = character && character.status === 'stepped-out';
    const isInactive = isDead || isSteppedOut;

    const menu = document.createElement('div');
    menu.className = 'character-context-menu';

    let menuHtml = '';
    if (isInactive) {
        if (isDead) {
            menuHtml += `<div class="menu-item" onclick="openReturnDialog('${characterId}', 'revive')"><span class="icon icon-heart"></span> Revive / Resurrect</div>`;
        } else {
            menuHtml += `<div class="menu-item" onclick="openReturnDialog('${characterId}', 'return')"><span class="icon icon-walk"></span> Step Back In</div>`;
        }
    } else {
        menuHtml += `<div class="menu-item" onclick="openStepOutDialog('${characterId}')"><span class="icon icon-walk"></span> Step Out</div>`;
    }
    menuHtml += `
        <div class="menu-item" onclick="editCharacter('${characterId}')"><span class="icon icon-edit"></span> Edit Character</div>
        <div class="menu-item" onclick="replaceCharacterImage('${characterId}')"><span class="icon icon-image"></span> Replace Image</div>
        <div class="menu-item" onclick="viewCharacterMemories('${characterId}')"><span class="icon icon-brain"></span> View Memories</div>
        <div class="menu-separator"></div>
        <div class="menu-item danger" onclick="removeCharacterFromChat('${characterId}')"><span class="icon icon-trash"></span> Remove from Chat</div>
    `;
    menu.innerHTML = menuHtml;

    menu.style.cssText = `position:fixed; left:${event.clientX}px; top:${event.clientY}px; z-index:10000;`;
    document.body.appendChild(menu);

    setTimeout(() => {
        document.addEventListener('click', function closeMenu(e) {
            if (!menu.contains(e.target)) {
                menu.remove();
                document.removeEventListener('click', closeMenu);
            }
        });
    }, 100);
}

// ---------------------------------------------------------------------------
// Step-Out Dialog
// ---------------------------------------------------------------------------

function openStepOutDialog(characterId) {
    const cd = window.sidebarState.currentChatData;
    const character = cd.characters.find(c => c.id === characterId);
    if (!character) return;

    const existingMenu = document.querySelector('.character-context-menu');
    if (existingMenu) existingMenu.remove();

    const modal = document.createElement('div');
    modal.className = 'stepout-modal-overlay';
    modal.style.cssText = `position:fixed;top:0;left:0;right:0;bottom:0;background:var(--overlay-scrim);display:flex;align-items:center;justify-content:center;z-index:10001;backdrop-filter:blur(4px);`;

    modal.innerHTML = `
        <div style="background:var(--shell-panel-bg);border-radius:14px;width:min(520px,95%);max-width:520px;padding:24px;border:1px solid rgba(255,255,255,0.12);">
            <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:16px;">
                <div>
                    <div style="font-size:18px;font-weight:700;color:#fff;">Step Out</div>
                    <div style="font-size:13px;color:rgba(255,255,255,0.65);">${character.name} will be marked as inactive.</div>
                </div>
                <button onclick="document.body.removeChild(document.querySelector('.stepout-modal-overlay'))" style="background:var(--surface-hover);border:none;color:#fff;width:32px;height:32px;border-radius:8px;cursor:pointer;">✕</button>
            </div>
            <div style="display:flex;gap:12px;margin-bottom:16px;">
                <button id="stepout-mode-regular" class="stepout-mode-btn active">Stepped Out</button>
                <button id="stepout-mode-death" class="stepout-mode-btn">Death</button>
            </div>
            <div style="margin-bottom:16px;">
                <label for="stepout-reason" style="display:block;margin-bottom:8px;font-size:13px;color:rgba(255,255,255,0.75);">Reason</label>
                <textarea id="stepout-reason" rows="4" placeholder="Enter the reason for ${character.name} stepping out..." style="width:100%;background:var(--surface-hover);border:1px solid rgba(255,255,255,0.12);color:#fff;border-radius:10px;padding:12px;resize:vertical;"></textarea>
            </div>
            <div style="display:flex;gap:12px;justify-content:flex-end;">
                <button onclick="document.body.removeChild(document.querySelector('.stepout-modal-overlay'))" style="padding:10px 16px;background:var(--surface-hover);border:1px solid rgba(255,255,255,0.14);color:#fff;border-radius:8px;cursor:pointer;">Cancel</button>
                <button id="stepout-confirm-btn" style="padding:10px 16px;background:rgba(0,184,115,0.18);border:1px solid rgba(0,184,115,0.35);color:#00b873;border-radius:8px;cursor:pointer;">Confirm</button>
            </div>
        </div>
    `;

    document.body.appendChild(modal);

    const regularBtn = modal.querySelector('#stepout-mode-regular');
    const deathBtn   = modal.querySelector('#stepout-mode-death');
    const reasonInput = modal.querySelector('#stepout-reason');
    const confirmBtn  = modal.querySelector('#stepout-confirm-btn');

    const updateMode = (mode) => {
        if (mode === 'death') {
            regularBtn.classList.remove('active');
            deathBtn.classList.add('active');
            reasonInput.placeholder = `Enter the death details for ${character.name}...`;
        } else {
            regularBtn.classList.add('active');
            deathBtn.classList.remove('active');
            reasonInput.placeholder = `Enter the reason for ${character.name} stepping out...`;
        }
        reasonInput.dataset.mode = mode;
    };

    regularBtn.addEventListener('click', () => updateMode('stepout'));
    deathBtn.addEventListener('click',   () => updateMode('death'));
    updateMode('stepout');

    confirmBtn.addEventListener('click', async () => {
        const mode   = reasonInput.dataset.mode || 'stepout';
        const reason = reasonInput.value.trim() || (mode === 'death' ? 'died unexpectedly' : 'stepped away temporarily');
        await applyCharacterStepOut(characterId, mode, reason);
        document.body.removeChild(modal);
    });
}

async function applyCharacterStepOut(characterId, mode, reason) {
    const cd = window.sidebarState.currentChatData;
    const character = cd.characters.find(c => c.id === characterId);
    if (!character) return;

    character.status        = mode === 'death' ? 'dead' : 'stepped-out';
    character.active        = false;
    character.stepOutReason = reason;
    character.stepOutMode   = mode;

    // Also update the bridge state so refreshAllPanels picks up the change
    if (window.chatBridge) {
        const bridgeChars = window.chatBridge.getState().characters;
        const bc = bridgeChars.find(c => c.id === characterId);
        if (bc) {
            bc.status        = character.status;
            bc.active        = false;
            bc.stepOutReason = reason;
            bc.stepOutMode   = mode;
        }
        await window.chatBridge.refreshAllPanels();
    }

    // Persist the status change so it survives reloads
    saveChatData();
    if (window.updateCharacterPanel) await window.updateCharacterPanel();

    const eventText = mode === 'death'
        ? `${character.name} has died: ${reason}`
        : `${character.name} has stepped out: ${reason}`;

    if (typeof queueNarratorEvent === 'function') {
        queueNarratorEvent(eventText);
    } else {
        addSystemMessage(eventText);
    }

    if (window.go && window.go.main && window.go.main.App && window.go.main.App.AddSystemMessage) {
        window.go.main.App.AddSystemMessage(eventText).catch(err => {
            console.error('Failed to add system event to backend:', err);
        });
    }

    // Auto-create memory entry for death events
    if (mode === 'death') {
        createDeathMemoryEntry(character, reason);
    }

    showNotification(`Event queued: ${character.name} ${mode === 'death' ? 'died' : 'stepped out'}`, 'success');
}

// ---------------------------------------------------------------------------
// Return Dialog (Revive / Step Back In)
// ---------------------------------------------------------------------------

function openReturnDialog(characterId, returnType) {
    const cd = window.sidebarState.currentChatData;
    const character = cd.characters.find(c => c.id === characterId);
    if (!character) return;

    const existingMenu = document.querySelector('.character-context-menu');
    if (existingMenu) existingMenu.remove();

    const isRevive = returnType === 'revive';
    const title = isRevive ? 'Revive / Resurrect' : 'Step Back In';
    const subtitle = isRevive
        ? `${character.name} will return to the living.`
        : `${character.name} will rejoin the scene.`;
    const placeholder = isRevive
        ? `Enter how ${character.name} is revived...`
        : `Enter why ${character.name} is returning...`;
    const defaultReason = isRevive ? 'resurrected by unknown means' : 'returned to the scene';

    const modal = document.createElement('div');
    modal.className = 'return-modal-overlay';
    modal.style.cssText = `position:fixed;top:0;left:0;right:0;bottom:0;background:var(--overlay-scrim);display:flex;align-items:center;justify-content:center;z-index:10001;backdrop-filter:blur(4px);`;

    modal.innerHTML = `
        <div style="background:var(--shell-panel-bg);border-radius:14px;width:min(520px,95%);max-width:520px;padding:24px;border:1px solid rgba(255,255,255,0.12);">
            <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:16px;">
                <div>
                    <div style="font-size:18px;font-weight:700;color:#fff;">${title}</div>
                    <div style="font-size:13px;color:rgba(255,255,255,0.65);">${subtitle}</div>
                </div>
                <button onclick="document.body.removeChild(document.querySelector('.return-modal-overlay'))" style="background:var(--surface-hover);border:none;color:#fff;width:32px;height:32px;border-radius:8px;cursor:pointer;">&#10005;</button>
            </div>
            <div style="margin-bottom:16px;">
                <label for="return-reason" style="display:block;margin-bottom:8px;font-size:13px;color:rgba(255,255,255,0.75);">${isRevive ? 'Method / Reason' : 'Reason'}</label>
                <textarea id="return-reason" rows="4" placeholder="${placeholder}" style="width:100%;background:var(--surface-hover);border:1px solid rgba(255,255,255,0.12);color:#fff;border-radius:10px;padding:12px;resize:vertical;"></textarea>
            </div>
            <div style="display:flex;gap:12px;justify-content:flex-end;">
                <button onclick="document.body.removeChild(document.querySelector('.return-modal-overlay'))" style="padding:10px 16px;background:var(--surface-hover);border:1px solid rgba(255,255,255,0.14);color:#fff;border-radius:8px;cursor:pointer;">Cancel</button>
                <button id="return-confirm-btn" style="padding:10px 16px;background:rgba(0,184,115,0.18);border:1px solid rgba(0,184,115,0.35);color:#00b873;border-radius:8px;cursor:pointer;">Confirm</button>
            </div>
        </div>
    `;

    document.body.appendChild(modal);

    const reasonInput = modal.querySelector('#return-reason');
    const confirmBtn  = modal.querySelector('#return-confirm-btn');

    confirmBtn.addEventListener('click', async () => {
        const reason = reasonInput.value.trim() || defaultReason;
        await applyCharacterReturn(characterId, returnType, reason);
        document.body.removeChild(modal);
    });
}

async function applyCharacterReturn(characterId, returnType, reason) {
    const cd = window.sidebarState.currentChatData;
    const character = cd.characters.find(c => c.id === characterId);
    if (!character) return;

    const wasDead = character.status === 'dead';
    character.status        = 'ready';
    character.active        = true;
    character.stepOutReason = null;
    character.stepOutMode   = null;

    // Also update the bridge state so refreshAllPanels picks up the change
    if (window.chatBridge) {
        const bridgeChars = window.chatBridge.getState().characters;
        const bc = bridgeChars.find(c => c.id === characterId);
        if (bc) {
            bc.status        = 'ready';
            bc.active        = true;
            bc.stepOutReason = null;
            bc.stepOutMode   = null;
        }
        await window.chatBridge.refreshAllPanels();
    }

    // Persist the status change so it survives reloads
    saveChatData();
    if (window.updateCharacterPanel) await window.updateCharacterPanel();

    const eventText = wasDead
        ? `${character.name} has been revived: ${reason}`
        : `${character.name} has stepped back in: ${reason}`;

    if (typeof queueNarratorEvent === 'function') {
        queueNarratorEvent(eventText);
    } else {
        addSystemMessage(eventText);
    }

    if (window.go && window.go.main && window.go.main.App && window.go.main.App.AddSystemMessage) {
        window.go.main.App.AddSystemMessage(eventText).catch(err => {
            console.error('Failed to add return event to backend:', err);
        });
    }

    showNotification(`${character.name} is back!`, 'success');
}

// ---------------------------------------------------------------------------
// Death Lorebook Entry
// ---------------------------------------------------------------------------

function createDeathMemoryEntry(character, reason) {
    const cd = window.sidebarState?.currentChatData || {};
    const worldName    = cd.world?.name    || cd.world?.title    || 'Unknown World';
    const scenarioName = cd.scenario?.name || cd.scenario?.title || 'Unknown Scenario';
    const now          = new Date();
    const dateStr      = now.toLocaleDateString();
    const timeStr      = now.toLocaleTimeString();

    const summary = `${character.name} died on ${dateStr} at ${timeStr}. Reason: ${reason}`;
    const content = [
        `On ${dateStr} at ${timeStr}, ${character.name} died.`,
        `Reason: ${reason}`,
        `World: ${worldName}`,
        `Scenario: ${scenarioName}`,
        `Characters present: ${(cd.characters || []).filter(c => c.id !== character.id).map(c => c.name).join(', ') || 'None'}`,
        `Active lorebooks: ${(cd.lorebooks || []).map(l => l.title || l.name).join(', ') || 'None'}`
    ].join('\n');

    if (window.MemoryStore) {
        window.MemoryStore.createMemoryEntry({
            title:       `Death of ${character.name}`,
            summary:     summary,
            content:     content,
            scope:       'character',
            ownerId:     character.id,
            ownerKind:   'character_id',
            tags:        ['death', 'event', character.name.toLowerCase().replace(/\s+/g, '-')],
            source:      'auto',
            trigger:     'death',
            confidence:  100
        });
    } else {
        console.warn('MemoryStore not available — death event not persisted to memory system');
    }

    console.log(`Created automated memory entry for death of ${character.name}`);
}

// ---------------------------------------------------------------------------
// Image Replacement
// ---------------------------------------------------------------------------

/**
 * Update a character's image from a pre-processed base64 data string.
 * Used when image data is already in hand (e.g. from a drag-drop handler).
 * For file-picker initiated replacements use replaceCharacterImage instead.
 * @param {string} characterId
 * @param {string} imageData - base64 data URL
 */
function updateCharacterImage(characterId, imageData) {
    const cd        = window.sidebarState?.currentChatData;
    const character = cd?.characters?.find(c => c.id === characterId);
    if (!character) return;

    character.image          = imageData;
    character._resolvedImage = imageData;

    if (window.CharacterStorage) {
        window.CharacterStorage.save(character).catch(err => console.warn('CharacterStorage save failed:', err));
    }

    if (window.chatBridge) {
        const bc = window.chatBridge.getState().characters.find(c => c.id === characterId);
        if (bc) { bc.image = imageData; bc._resolvedImage = imageData; }
        window.chatBridge.refreshAllPanels();
    } else {
        if (typeof window.updateCharacterPanel === 'function') window.updateCharacterPanel();
    }

    if (typeof showNotification === 'function') showNotification('Character image updated', 'success');
}

function replaceCharacterImage(characterId) {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'image/*';
    input.onchange = async (e) => {
        const file = e.target.files[0];
        if (!file) return;
        try {
            const result    = await window.ImageProcessor.process(file);
            const cd        = window.sidebarState.currentChatData;
            const character = cd.characters.find(c => c.id === characterId);
            if (!character) return;

            const imageEntry = {
                id:          `img_${Date.now()}_${Math.random().toString(36).substr(2, 6)}`,
                name:        result.meta.name,
                size:        result.meta.size,
                data:        result.thumbnail,
                medium:      result.medium,
                original:    result.original,
                isThumbnail: true
            };
            character.image          = result.thumbnail;
            character._resolvedImage = result.thumbnail;
            character.howl           = { ...(character.howl || {}), images: [imageEntry] };

            if (window.CharacterStorage) {
                await window.CharacterStorage.save(character).catch(err => console.warn('CharacterStorage save failed:', err));
            }
            if (window.chatBridge) {
                const bc = window.chatBridge.getState().characters.find(c => c.id === characterId);
                if (bc) { bc.image = result.thumbnail; bc._resolvedImage = result.thumbnail; }
                window.chatBridge.refreshAllPanels();
            } else {
                saveChatData();
                if (window.updateCharacterPanel) window.updateCharacterPanel();
            }
            showNotification('Character image updated', 'success');
        } catch (err) {
            console.error('Image replacement failed:', err);
            showNotification('Failed to update image: ' + err.message, 'error');
        }
    };
    input.click();
}

// ---------------------------------------------------------------------------
// Misc Character Actions
// ---------------------------------------------------------------------------

function viewCharacterMemories(characterId) {
    if (window.CharacterMemories && window.CharacterMemories.view) {
        window.CharacterMemories.view(characterId);
    } else {
        console.warn('CharacterMemories module not loaded');
        showNotification('Memory viewer loading...', 'info');
    }
}

async function removeCharacterFromChat(characterId) {
    if (confirm('Remove this character from the chat?')) {
        if (window.chatBridge) {
            await window.chatBridge.applyCharacterRemove(characterId);
        } else {
            const cd = window.sidebarState.currentChatData;
            cd.characters = cd.characters.filter(c => c.id !== characterId);
            if ((cd.activeCharacterIds || []).includes(characterId)) {
                if (!cd.removedCharacterIds) cd.removedCharacterIds = [];
                if (!cd.removedCharacterIds.includes(characterId)) {
                    cd.removedCharacterIds.push(characterId);
                }
            }
            cd.addedCharacterIds = (cd.addedCharacterIds || []).filter(id => id !== characterId);
            saveChatData();
            if (window.updateCharacterPanel) window.updateCharacterPanel();
        }
        showNotification('Character removed from chat', 'success');
    }
}

// ---------------------------------------------------------------------------
// Add Character Modal
// ---------------------------------------------------------------------------

function addCharacterToChat() {
    // The canonical store assigns identity; this list is display-only.
    const characters = (window.CharacterStorage?.all?.() || []);

    const cd = window.sidebarState.currentChatData;
    const modal = document.createElement('div');
    modal.className = 'selector-modal';
    modal.style.cssText = `position:fixed;top:0;left:0;right:0;bottom:0;background:var(--overlay-scrim);display:flex;align-items:center;justify-content:center;z-index:9999;backdrop-filter:blur(4px);`;

    let content = `
        <div style="background:var(--shell-panel-bg);border-radius:12px;padding:24px;max-width:400px;border:1px solid var(--surface-hover-strong);">
            <h2 style="color:#fff;margin-top:0;margin-bottom:16px;">Select a Character</h2>
            <div style="max-height:300px;overflow-y:auto;margin-bottom:16px;">
    `;

    if (characters.length === 0) {
        content += `<div style="color:var(--surface-hover);text-align:center;padding:20px;">No characters available</div>`;
    } else {
        characters.forEach(character => {
            const isActive = cd.characters && cd.characters.find(c => c.id === character.id);
            content += `
                <div onclick="selectCharacterForChat('${character.id}'); document.querySelector('.selector-modal').remove();" style="
                    padding:12px;background:${isActive ? 'var(--success-bg)' : 'var(--surface-hover)'};
                    margin-bottom:8px;border-radius:6px;cursor:pointer;
                    border:1px solid ${isActive ? 'var(--success)' : 'rgba(255,255,255,0.1)'};color:#fff;transition:all 0.2s ease;"
                    onmouseover="this.style.background='var(--surface-hover-strong)'"
                    onmouseout="this.style.background='${isActive ? 'var(--success-bg)' : 'var(--surface-hover)'}'">
                    <strong style="display:block;margin-bottom:4px;">${character.name}</strong>
                    <small style="color:var(--surface-hover);">${character.description || 'No description'}</small>
                </div>
            `;
        });
    }

    content += `
            </div>
            <div style="display:flex;gap:8px;">
                <button onclick="document.querySelector('.selector-modal').remove();" style="flex:1;padding:10px;background:var(--surface-hover-strong);border:1px solid var(--border-secondary);color:#fff;border-radius:6px;cursor:pointer;">Cancel</button>
                <button onclick="navigateTo('character.html'); document.querySelector('.selector-modal').remove();" style="flex:1;padding:10px;background:var(--success-bg);border:1px solid var(--success);color:#00b873;border-radius:6px;cursor:pointer;">+ New Character</button>
            </div>
        </div>
    `;

    modal.innerHTML = content;
    document.body.appendChild(modal);
    modal.addEventListener('click', (e) => { if (e.target === modal) modal.remove(); });
}

async function selectCharacterForChat(characterId) {
    const characters = (window.CharacterStorage?.all?.() || []);
    const character  = characters.find(c => c.id === characterId);
    if (!character) { console.error('Character not found with ID:', characterId); return; }

    const cd = window.sidebarState.currentChatData;
    if (cd.characters && cd.characters.some(c => c.id === characterId)) {
        showNotification('Character already in chat', 'info');
        return;
    }

    if (window.chatBridge) {
        await window.chatBridge.applyCharacterAdd(character);
    } else {
        if (!cd.characters) cd.characters = [];
        cd.characters.push(character);
        saveChatData();
        if (window.updateCharacterPanel) window.updateCharacterPanel();
    }
    showNotification(`Character added to chat: ${character.name}`, 'success');
}

// ---------------------------------------------------------------------------
// Global Exports
// ---------------------------------------------------------------------------

window.sidebarCharacters = {
    updateCharacterPanel: (...args) => window.updateCharacterPanel?.(...args),
    createCharacterItem,
    openCharacterSheet,
    editCharacter,
    showCharacterMenu,
    openStepOutDialog,
    applyCharacterStepOut,
    openReturnDialog,
    applyCharacterReturn,
    replaceCharacterImage,
    updateCharacterImage,
    viewCharacterMemories,
    removeCharacterFromChat,
    addCharacterToChat,
    selectCharacterForChat
};
