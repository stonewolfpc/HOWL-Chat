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
        const catalog = (window.CharacterStorage?.all?.() || []);
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
    const fe           = window.FilterEngine;

    const manualLorebooks = allLorebooks.filter(l => !l.automated);

    const inWorldRefs    = (id) => cd.world    && (
        (Array.isArray(cd.world.howl?.lorebooks)    && cd.world.howl.lorebooks.some(r => r.id === id)) ||
        (Array.isArray(cd.world.lorebooks)           && cd.world.lorebooks.some(r => r.id === id))
    );
    const inScenarioRefs = (id) => cd.scenario && (
        (Array.isArray(cd.scenario.howl?.lorebooks) && cd.scenario.howl.lorebooks.some(r => r.id === id)) ||
        (Array.isArray(cd.scenario.lorebooks)        && cd.scenario.lorebooks.some(r => r.id === id))
    );
    const _charCatalog   = (window.CharacterStorage?.all?.() || []);
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
    modal.style.cssText = `position:fixed;top:0;left:0;right:0;bottom:0;background:var(--overlay-scrim);display:flex;align-items:center;justify-content:center;z-index:9999;backdrop-filter:blur(4px);`;

    let sortField = 'name';
    let sortDir   = 'asc';

    const buildListHtml = (searchTerm) => {
        let html = '';
        let anyItems = false;
        const term = (searchTerm || '').toLowerCase();

        selectorGroups.forEach(({ key, label, color, items }) => {
            let groupItems = items.slice();
            if (fe && (sortField || sortDir)) {
                groupItems = fe.sortItems(groupItems, sortField, sortDir);
            }
            if (term) {
                groupItems = groupItems.filter(l => {
                    const name = (l.name || l.title || '').toLowerCase();
                    const desc = (l.description || '').toLowerCase();
                    return name.includes(term) || desc.includes(term);
                });
            }
            if (groupItems.length === 0) return;
            anyItems = true;
            html += `<div style="color:${color};font-size:11px;font-weight:600;text-transform:uppercase;letter-spacing:0.5px;padding:8px 0 4px;opacity:0.9;">${label}</div>`;
            groupItems.forEach(l => { html += buildLorebookItem(l, cd, key); });
        });

        if (!anyItems) {
            html += `<div style="color:var(--surface-hover);text-align:center;padding:24px;font-style:italic;">No lorebooks available</div>`;
        }
        return html;
    };

    const content = `
        <div style="background:var(--shell-panel-bg);border-radius:12px;padding:24px;max-width:600px;width:90%;max-height:80vh;border:1px solid var(--surface-hover-strong);overflow-y:auto;">
            <h2 style="color:#fff;margin-top:0;margin-bottom:4px;">Select Lorebooks</h2>
            <p style="color:var(--surface-hover);font-size:12px;margin-top:0;margin-bottom:16px;">Manual lorebooks only. Automated memories are in Story Tools.</p>
            <div style="display:flex;gap:8px;margin-bottom:12px;align-items:center;">
                <input type="text" id="lorebook-search" placeholder="Search lorebooks..."
                    style="flex:1;box-sizing:border-box;padding:8px 12px;background:var(--surface-hover-strong);border:1px solid var(--border-secondary);border-radius:6px;color:#fff;font-size:14px;"
                    onkeyup="window._lorebookSelectorSearch && window._lorebookSelectorSearch(this.value)">
                <select id="lorebook-sort"
                    style="padding:8px 12px;border:1px solid var(--border-secondary);border-radius:6px;font-size:14px;cursor:pointer;"
                    onchange="window._lorebookSelectorSort && window._lorebookSelectorSort(this.value)">
                    <option value="name:asc">Name A-Z</option>
                    <option value="name:desc">Name Z-A</option>
                    <option value="date:desc">Newest First</option>
                    <option value="date:asc">Oldest First</option>
                </select>
            </div>
            <div id="available-lorebooks-list">${buildListHtml('')}</div>
            <div style="display:flex;gap:8px;margin-top:16px;">
                <button onclick="document.querySelector('.lorebook-selector-modal').remove();" style="flex:1;padding:10px;background:var(--surface-hover-strong);border:1px solid var(--border-secondary);color:#fff;border-radius:6px;cursor:pointer;">Done</button>
                <button onclick="navigateTo('lorebooks.html'); document.querySelector('.lorebook-selector-modal').remove();" style="flex:1;padding:10px;background:var(--success-bg);border:1px solid var(--success);color:#00b873;border-radius:6px;cursor:pointer;">+ New Lorebook</button>
            </div>
        </div>
    `;

    window._lorebookSelectorSearch = (val) => {
        const list = document.getElementById('available-lorebooks-list');
        if (list) list.innerHTML = buildListHtml(val);
    };
    window._lorebookSelectorSort = (val) => {
        const [field, dir] = val.split(':');
        sortField = field;
        sortDir = dir;
        const searchVal = document.getElementById('lorebook-search')?.value || '';
        const list = document.getElementById('available-lorebooks-list');
        if (list) list.innerHTML = buildListHtml(searchVal);
    };

    modal.innerHTML = content;
    document.body.appendChild(modal);
    setTimeout(() => { const ls = document.getElementById('lorebook-search'); if (ls) ls.focus(); }, 100);
}

function buildLorebookItem(lorebook, cd, sourceKey) {
    const isActive  = cd.lorebooks && cd.lorebooks.some(l => l.id === lorebook.id);
    const lbName    = (lorebook.name || lorebook.title || '').toLowerCase();
    return `
        <div class="lorebook-item available-lorebook" data-lorebook-id="${lorebook.id}" data-lorebook-name="${lbName}" data-source="${sourceKey}" onclick="toggleChatLorebook('${lorebook.id}');" style="
            padding:12px;background:${isActive ? 'var(--success-bg)' : 'var(--surface-hover)'};
            margin-bottom:8px;border-radius:6px;cursor:pointer;position:relative;
            border:1px solid ${isActive ? 'var(--success)' : 'rgba(255,255,255,0.1)'};color:#fff;transition:all 0.2s ease;">
            ${isActive ? '<span style="position:absolute;top:8px;right:8px;color:#00b873;">✓</span>' : ''}
            <strong style="display:block;margin-bottom:4px;">${lorebook.name || lorebook.title || 'Untitled'}</strong>
            <small style="color:var(--surface-hover);">${lorebook.description || 'No description'}</small>
        </div>
    `;
}

function filterLorebooks(searchTerm) {
    const term = searchTerm.toLowerCase();
    document.querySelectorAll('.available-lorebook').forEach(item => {
        item.style.display = item.dataset.lorebookName.includes(term) ? 'block' : 'none';
    });
}

// Chat-level toggle (by id). The world/scenario dropdowns have their own
// toggleLorebookSelection(lorebook, element); same name, different job,
// which is why this one is named for what it does.
async function toggleChatLorebook(lorebookId) {
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
