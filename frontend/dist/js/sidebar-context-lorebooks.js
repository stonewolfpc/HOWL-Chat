/**
 * Sidebar Context — Lorebook Section
 *
 * Lorebook list display, lorebook selector modal (with world-attachment
 * filtering and search), lorebook toggle selection, remove, and manager nav.
 *
 * Depends on: sidebar-state.js, sidebar-session.js
 *
 * @module sidebar-context-lorebooks
 * @version 2.0.0
 */

/**
 * Determines the source type of a lorebook ('world', 'scenario', 'character', 'standalone')
 * by cross-referencing the current chat's world, scenario, and characters.
 * @param {Object} lorebook
 * @param {Object} cd - currentChatData
 * @returns {'world'|'scenario'|'character'|'standalone'}
 */
function getLorebookSourceType(lorebook, cd) {
    const id = (lorebook.id || '').toString().trim().toLowerCase();

    const inRefs = (refs) => Array.isArray(refs) && refs.some(r => (r.id || '').toString().trim().toLowerCase() === id);

    if (cd.world) {
        if (inRefs(cd.world.howl?.lorebooks) || inRefs(cd.world.lorebooks)) return 'world';
    }
    if (cd.scenario) {
        if (inRefs(cd.scenario.howl?.lorebooks) || inRefs(cd.scenario.lorebooks)) return 'scenario';
    }
    if (Array.isArray(cd.characters) && cd.characters.length > 0) {
        const charIds = cd.characters.map(c => (c.id || '').toString().trim().toLowerCase()).filter(Boolean);
        const catalog = JSON.parse(localStorage.getItem('howl-chat-characters') || '[]');
        for (const c of cd.characters) {
            if (inRefs(c.howl?.lorebooks) || inRefs(c.lorebooks)) return 'character';
        }
        for (const cid of charIds) {
            const full = catalog.find(c => (c.id || '').toString().trim().toLowerCase() === cid);
            if (full && (inRefs(full.howl?.lorebooks) || inRefs(full.lorebooks))) return 'character';
        }
    }
    return 'standalone';
}

/**
 * Renders the sidebar lorebook panel, grouped by World / Scenario / Character.
 * Automated lorebooks (automated: true) are never shown here.
 */
function updateLorebookDisplay() {
    const lorebookList = document.getElementById('lorebook-list');
    if (!lorebookList) return;
    if (!window.sidebarState) return;

    const cd = window.sidebarState.currentChatData;
    lorebookList.innerHTML = '';

    const manualLorebooks = (cd.lorebooks || []).filter(l => !l.automated);

    if (manualLorebooks.length === 0) {
        lorebookList.innerHTML = `
            <div class="lorebook-empty">
                <div class="empty-icon">📚</div>
                <div class="empty-text">No lorebooks selected</div>
                <button class="add-lorebook-btn" onclick="openLorebookSelector()">Add Lorebook</button>
            </div>
        `;
        return;
    }

    const groups = { world: [], scenario: [], character: [], standalone: [] };
    manualLorebooks.forEach(lb => groups[getLorebookSourceType(lb, cd)].push(lb));

    const groupConfig = [
        { key: 'world',      label: 'World Lorebooks',      color: '#00b873' },
        { key: 'scenario',   label: 'Scenario Lorebooks',   color: '#007acc' },
        { key: 'character',  label: 'Character Lorebooks',  color: '#9c27b0' },
        { key: 'standalone', label: 'Additional Lorebooks', color: '#ff9500' }
    ];

    groupConfig.forEach(({ key, label, color }) => {
        if (groups[key].length === 0) return;
        const header = document.createElement('div');
        header.className = 'lorebook-group-header';
        header.style.cssText = `color:${color};font-size:11px;font-weight:600;text-transform:uppercase;letter-spacing:0.5px;padding:6px 4px 2px;opacity:0.8;`;
        header.textContent = label;
        lorebookList.appendChild(header);
        groups[key].forEach(lb => lorebookList.appendChild(createLorebookItem(lb, key)));
    });

    const addBtn = document.createElement('div');
    addBtn.className = 'add-lorebook-section';
    addBtn.innerHTML = `<button class="add-lorebook-btn" onclick="openLorebookSelector()">Add Lorebook</button>`;
    lorebookList.appendChild(addBtn);
}

function _resolveLorebookDescription(lorebook) {
    if (lorebook.description) return lorebook.description;
    try {
        const catalog = JSON.parse(localStorage.getItem('howl-chat-lorebooks') || '[]');
        const full = catalog.find(l => l.id === lorebook.id);
        return (full && full.description) ? full.description : '';
    } catch (_) { return ''; }
}

function createLorebookItem(lorebook, sourceType) {
    const item = document.createElement('div');
    const isEnabled = lorebook.enabled !== false;
    item.className = 'lorebook-item' + (isEnabled ? '' : ' disabled');
    item.dataset.lorebookId = lorebook.id;
    const desc = _resolveLorebookDescription(lorebook);
    item.innerHTML = `
        <div class="lorebook-icon">📚</div>
        <div class="lorebook-info">
            <div class="lorebook-name">${lorebook.name || lorebook.title || 'Untitled Lorebook'}</div>
            ${desc ? `<div class="lorebook-description">${desc}</div>` : ''}
        </div>
        <button class="lorebook-menu-btn" onclick="showLorebookMenu(event, '${lorebook.id}')">⋮</button>
        <button class="lorebook-remove" onclick="removeLorebookFromChat('${lorebook.id}')">×</button>
    `;
    item.addEventListener('contextmenu', (e) => {
        e.preventDefault();
        showLorebookMenu(e, lorebook.id);
    });
    return item;
}

async function removeLorebookFromChat(lorebookId) {
    if (confirm('Remove this lorebook from the chat?')) {
        if (window.chatBridge) {
            const cd = window.sidebarState.currentChatData;
            const lb = (cd.lorebooks || []).find(l => l.id === lorebookId);
            if (lb) await window.chatBridge.applyLorebookToggle(lb);
        } else {
            const cd = window.sidebarState.currentChatData;
            cd.lorebooks = cd.lorebooks.filter(l => l.id !== lorebookId);
            saveChatData();
            updateLorebookDisplay();
        }
        showNotification('Lorebook removed from chat', 'success');
    }
}

function showLorebookMenu(event, lorebookId) {
    event.preventDefault();
    event.stopPropagation();

    const existingMenu = document.querySelector('.lorebook-context-menu');
    if (existingMenu) existingMenu.remove();

    const cd = window.sidebarState.currentChatData;
    const lorebook = (cd.lorebooks || []).find(l => l.id === lorebookId);
    if (!lorebook) return;

    const isEnabled = lorebook.enabled !== false;
    const menu = document.createElement('div');
    menu.className = 'lorebook-context-menu';
    menu.innerHTML = `
        <div class="menu-item" onclick="toggleLorebookEnabled('${lorebookId}')">
            <span class="icon icon-toggle"></span>
            ${isEnabled ? 'Disable Lorebook' : 'Enable Lorebook'}
        </div>
        <div class="menu-separator"></div>
        <div class="menu-item danger" onclick="removeLorebookFromChat('${lorebookId}')">
            <span class="icon icon-trash"></span> Remove from Chat
        </div>
    `;

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

async function toggleLorebookEnabled(lorebookId) {
    const cd = window.sidebarState.currentChatData;
    const lorebook = (cd.lorebooks || []).find(l => l.id === lorebookId);
    if (!lorebook) return;

    lorebook.enabled = lorebook.enabled === false;
    saveChatData();
    updateLorebookDisplay();

    // Sync enabled IDs to backend so UpdateLorebookContext only activates enabled ones.
    const enabledIds = (cd.lorebooks || []).filter(l => l.enabled !== false).map(l => l.id);
    if (window.go?.main?.App?.UpdateLorebookContext) {
        try {
            await window.go.main.App.UpdateLorebookContext(enabledIds);
            console.log('[toggleLorebookEnabled] synced', enabledIds.length, 'active lorebooks to backend');
        } catch (error) {
            console.warn('Backend lorebook enable sync failed:', error);
        }
    }

    const isEnabled = lorebook.enabled !== false;
    showNotification(`Lorebook ${isEnabled ? 'enabled' : 'disabled'}`, 'success');
}

function openLorebookSelector() {
    const allLorebooks = JSON.parse(localStorage.getItem('howl-chat-lorebooks') || '[]');
    const cd           = window.sidebarState.currentChatData;

    const manualLorebooks = allLorebooks.filter(l => !l.automated);

    const inWorldRefs    = (id) => cd.world    && (
        (Array.isArray(cd.world.howl?.lorebooks)    && cd.world.howl.lorebooks.some(r => r.id === id)) ||
        (Array.isArray(cd.world.lorebooks)           && cd.world.lorebooks.some(r => r.id === id))
    );
    const inScenarioRefs = (id) => cd.scenario && (
        (Array.isArray(cd.scenario.howl?.lorebooks) && cd.scenario.howl.lorebooks.some(r => r.id === id)) ||
        (Array.isArray(cd.scenario.lorebooks)        && cd.scenario.lorebooks.some(r => r.id === id))
    );
    const _charCatalog   = JSON.parse(localStorage.getItem('howl-chat-characters') || '[]');
    const _charIds       = Array.isArray(cd.characters) ? cd.characters.map(c => (c.id || '').toString().trim().toLowerCase()) : [];
    const inCharRefs     = (id) => {
        if (!Array.isArray(cd.characters)) return false;
        for (const c of cd.characters) {
            if ((Array.isArray(c.howl?.lorebooks) && c.howl.lorebooks.some(r => r.id === id)) ||
                (Array.isArray(c.lorebooks)       && c.lorebooks.some(r => r.id === id))) return true;
        }
        for (const cid of _charIds) {
            const full = _charCatalog.find(c => (c.id || '').toString().trim().toLowerCase() === cid);
            if (full && (
                (Array.isArray(full.howl?.lorebooks) && full.howl.lorebooks.some(r => r.id === id)) ||
                (Array.isArray(full.lorebooks)       && full.lorebooks.some(r => r.id === id))
            )) return true;
        }
        return false;
    };

    const selectorGroups = [
        {
            key:   'world',
            label: `World Lorebooks${cd.world ? ' — ' + (cd.world.title || cd.world.name || '') : ''}`,
            color: '#00b873',
            items: manualLorebooks.filter(l => inWorldRefs(l.id))
        },
        {
            key:   'scenario',
            label: `Scenario Lorebooks${cd.scenario ? ' — ' + (cd.scenario.title || cd.scenario.name || '') : ''}`,
            color: '#007acc',
            items: manualLorebooks.filter(l => inScenarioRefs(l.id))
        },
        {
            key:   'character',
            label: 'Character Lorebooks',
            color: '#9c27b0',
            items: manualLorebooks.filter(l => inCharRefs(l.id))
        },
        {
            key:   'standalone',
            label: 'Additional Lorebooks',
            color: '#ff9500',
            items: manualLorebooks.filter(l => !inWorldRefs(l.id) && !inScenarioRefs(l.id) && !inCharRefs(l.id))
        }
    ];

    const modal = document.createElement('div');
    modal.className   = 'lorebook-selector-modal';
    modal.style.cssText = `position:fixed;top:0;left:0;right:0;bottom:0;background:rgba(0,0,0,0.7);display:flex;align-items:center;justify-content:center;z-index:9999;backdrop-filter:blur(4px);`;

    let content = `
        <div style="background:rgba(30,30,40,0.98);border-radius:12px;padding:24px;max-width:600px;width:90%;max-height:80vh;border:1px solid rgba(255,255,255,0.1);overflow-y:auto;">
            <h2 style="color:#fff;margin-top:0;margin-bottom:4px;">Select Lorebooks</h2>
            <p style="color:rgba(255,255,255,0.4);font-size:12px;margin-top:0;margin-bottom:16px;">Manual lorebooks only. Automated memories are in Story Tools.</p>
            <div style="margin-bottom:12px;">
                <input type="text" id="lorebook-search" placeholder="Search lorebooks..." style="width:100%;box-sizing:border-box;padding:8px 12px;background:rgba(255,255,255,0.1);border:1px solid rgba(255,255,255,0.2);border-radius:6px;color:#fff;font-size:14px;" onkeyup="filterLorebooks(this.value)">
            </div>
            <div id="available-lorebooks-list">
    `;

    let anyItems = false;
    selectorGroups.forEach(({ key, label, color, items }) => {
        if (items.length === 0) return;
        anyItems = true;
        content += `<div style="color:${color};font-size:11px;font-weight:600;text-transform:uppercase;letter-spacing:0.5px;padding:8px 0 4px;opacity:0.9;">${label}</div>`;
        items.forEach(l => { content += buildLorebookItem(l, cd, key); });
    });

    if (!anyItems) {
        content += `<div style="color:rgba(255,255,255,0.3);text-align:center;padding:24px;font-style:italic;">No lorebooks available</div>`;
    }

    content += `
            </div>
            <div style="display:flex;gap:8px;margin-top:16px;">
                <button onclick="document.querySelector('.lorebook-selector-modal').remove();" style="flex:1;padding:10px;background:rgba(255,255,255,0.1);border:1px solid rgba(255,255,255,0.2);color:#fff;border-radius:6px;cursor:pointer;">Done</button>
                <button onclick="navigateTo('lorebooks.html'); document.querySelector('.lorebook-selector-modal').remove();" style="flex:1;padding:10px;background:rgba(0,184,115,0.2);border:1px solid rgba(0,184,115,0.3);color:#00b873;border-radius:6px;cursor:pointer;">+ New Lorebook</button>
            </div>
        </div>
    `;

    modal.innerHTML = content;
    document.body.appendChild(modal);
    setTimeout(() => { const ls = document.getElementById('lorebook-search'); if (ls) ls.focus(); }, 100);
}

function buildLorebookItem(lorebook, cd, sourceKey) {
    const isActive  = cd.lorebooks && cd.lorebooks.some(l => l.id === lorebook.id);
    const lbName    = (lorebook.name || lorebook.title || '').toLowerCase();
    return `
        <div class="lorebook-item available-lorebook" data-lorebook-id="${lorebook.id}" data-lorebook-name="${lbName}" data-source="${sourceKey}" onclick="toggleLorebookSelection('${lorebook.id}');" style="
            padding:12px;background:${isActive ? 'rgba(0,184,115,0.2)' : 'rgba(255,255,255,0.05)'};
            margin-bottom:8px;border-radius:6px;cursor:pointer;position:relative;
            border:1px solid ${isActive ? 'rgba(0,184,115,0.3)' : 'rgba(255,255,255,0.1)'};color:#fff;transition:all 0.2s ease;">
            ${isActive ? '<span style="position:absolute;top:8px;right:8px;color:#00b873;">✓</span>' : ''}
            <strong style="display:block;margin-bottom:4px;">${lorebook.name || lorebook.title || 'Untitled'}</strong>
            <small style="color:rgba(255,255,255,0.5);">${lorebook.description || 'No description'}</small>
        </div>
    `;
}

function filterLorebooks(searchTerm) {
    const term = searchTerm.toLowerCase();
    document.querySelectorAll('.available-lorebook').forEach(item => {
        item.style.display = item.dataset.lorebookName.includes(term) ? 'block' : 'none';
    });
}

async function toggleLorebookSelection(lorebookId) {
    const lorebooks = JSON.parse(localStorage.getItem('howl-chat-lorebooks') || '[]');
    const lorebook  = lorebooks.find(l => l.id === lorebookId);
    if (!lorebook) return;

    if (window.chatBridge) {
        await window.chatBridge.applyLorebookToggle(lorebook);
    } else {
        const cd = window.sidebarState.currentChatData;
        if (!cd.lorebooks) cd.lorebooks = [];
        const index    = cd.lorebooks.findIndex(l => l.id === lorebookId);
        const isAdding = index === -1;
        if (isAdding) { cd.lorebooks.push(lorebook); }
        else          { cd.lorebooks.splice(index, 1); }
        saveChatData();
        updateLorebookDisplay();
    }

    const cd      = window.sidebarState.currentChatData;
    const isNowOn = (cd.lorebooks || []).some(l => l.id === lorebookId);
    showNotification(`${isNowOn ? 'Added' : 'Removed'} lorebook: ${lorebook.name || lorebook.title}`, isNowOn ? 'success' : 'info');

    if (window.go?.main?.App?.UpdateLorebookContext) {
        const enabledIds = (cd.lorebooks || []).filter(l => l.enabled !== false).map(l => l.id);
        try { window.go.main.App.UpdateLorebookContext(enabledIds); }
        catch (error) { console.warn('Backend lorebook update failed:', error); }
    }
}

function openLorebookManager() {
    navigateTo('lorebooks.html');
}
