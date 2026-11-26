/**
 * Context Selectors — ContextDisplay Modal Pickers
 *
 * Mixin methods patched onto ContextDisplay.prototype for the
 * world, scenario, and character selection modals, plus the
 * hierarchical auto-population logic triggered on selection.
 *
 * Depends on: context-data.js, context-panels.js
 *
 * @module context-selectors
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Shared Modal Builder
// ---------------------------------------------------------------------------

function _buildSelectorModal(innerHtml) {
    const modal = document.createElement('div');
    modal.className  = 'selector-modal';
    modal.style.cssText = [
        'position:fixed', 'top:0', 'left:0', 'right:0', 'bottom:0',
        'background:var(--overlay-scrim)', 'display:flex', 'align-items:center',
        'justify-content:center', 'z-index:9999', 'backdrop-filter:blur(4px)'
    ].join(';');
    modal.innerHTML = innerHtml;
    document.body.appendChild(modal);
    modal.addEventListener('click', e => { if (e.target === modal) modal.remove(); });
    return modal;
}

function _itemRow(onclick, isActive, bodyHtml) {
    const bg     = isActive ? 'rgba(0,184,115,0.2)'  : 'rgba(255,255,255,0.05)';
    const border = isActive ? 'rgba(0,184,115,0.3)'  : 'rgba(255,255,255,0.1)';
    const bgOut  = isActive ? 'rgba(0,184,115,0.2)'  : 'rgba(255,255,255,0.05)';
    const bdOut  = isActive ? 'rgba(0,184,115,0.3)'  : 'rgba(255,255,255,0.1)';
    return `<div onclick="${onclick}"
        style="padding:12px;background:${bg};margin-bottom:8px;border-radius:6px;cursor:pointer;border:1px solid ${border};color:#fff;transition:all 0.2s ease;"
        onmouseover="this.style.background='var(--surface-hover-strong)';this.style.borderColor='var(--border-secondary)';"
        onmouseout="this.style.background='${bgOut}';this.style.borderColor='${bdOut}';">${bodyHtml}</div>`;
}

function _modalShell(title, listHtml, cancelOnclick, newLabel, newOnclick) {
    return `<div style="background:var(--shell-panel-bg);border-radius:12px;padding:24px;max-width:480px;width:90%;border:1px solid var(--surface-hover-strong);">
        <h2 style="color:#fff;margin-top:0;margin-bottom:16px;">${title}</h2>
        <div style="max-height:300px;overflow-y:auto;margin-bottom:16px;">${listHtml}</div>
        <div style="display:flex;gap:8px;">
            <button onclick="${cancelOnclick}" style="flex:1;padding:10px;background:var(--surface-hover-strong);border:1px solid var(--border-secondary);color:#fff;border-radius:6px;cursor:pointer;">Cancel</button>
            <button onclick="${newOnclick}"   style="flex:1;padding:10px;background:var(--success-bg);border:1px solid var(--success);color:#00b873;border-radius:6px;cursor:pointer;">${newLabel}</button>
        </div>
    </div>`;
}

// ---------------------------------------------------------------------------
// Filterable Modal Builder
// ---------------------------------------------------------------------------

function _filterBarHtml() {
    return `
        <div style="display:flex;gap:8px;margin-bottom:12px;align-items:center;">
            <input type="text" class="selector-search-input" placeholder="Search..."
                style="flex:1;padding:8px 12px;background:var(--surface-hover);border:1px solid var(--surface-hover-strong);border-radius:6px;color:#fff;outline:none;"
                oninput="window._selectorFilter && window._selectorFilter(this.value)">
            <select class="selector-sort-select"
                style="padding:8px 12px;border:1px solid var(--surface-hover-strong);border-radius:6px;outline:none;cursor:pointer;"
                onchange="window._selectorSort && window._selectorSort(this.value)">
                <option value="name:asc">Name A-Z</option>
                <option value="name:desc">Name Z-A</option>
                <option value="date:desc">Newest First</option>
                <option value="date:asc">Oldest First</option>
            </select>
        </div>`;
}

function _filterableModalShell(title, listHtml, cancelOnclick, newLabel, newOnclick) {
    return `<div style="background:var(--shell-panel-bg);border-radius:12px;padding:24px;max-width:480px;width:90%;border:1px solid var(--surface-hover-strong);">
        <h2 style="color:#fff;margin-top:0;margin-bottom:16px;">${title}</h2>
        ${_filterBarHtml()}
        <div class="selector-list-container" style="max-height:260px;overflow-y:auto;margin-bottom:16px;">${listHtml}</div>
        <div style="display:flex;gap:8px;">
            <button onclick="${cancelOnclick}" style="flex:1;padding:10px;background:var(--surface-hover-strong);border:1px solid var(--border-secondary);color:#fff;border-radius:6px;cursor:pointer;">Cancel</button>
            <button onclick="${newOnclick}"   style="flex:1;padding:10px;background:var(--success-bg);border:1px solid var(--success);color:#00b873;border-radius:6px;cursor:pointer;">${newLabel}</button>
        </div>
    </div>`;
}

function _filterBarWithAttachmentHtml(label) {
    return `
        <div style="display:flex;gap:8px;margin-bottom:12px;align-items:center;flex-wrap:wrap;">
            <input type="text" class="selector-search-input" placeholder="Search..."
                style="flex:1;min-width:120px;padding:8px 12px;background:var(--surface-hover);border:1px solid var(--surface-hover-strong);border-radius:6px;color:#fff;outline:none;"
                oninput="window._selectorFilter && window._selectorFilter(this.value)">
            <select class="selector-sort-select"
                style="padding:8px 12px;border:1px solid var(--surface-hover-strong);border-radius:6px;outline:none;cursor:pointer;"
                onchange="window._selectorSort && window._selectorSort(this.value)">
                <option value="name:asc">Name A-Z</option>
                <option value="name:desc">Name Z-A</option>
                <option value="date:desc">Newest First</option>
                <option value="date:asc">Oldest First</option>
            </select>
            <label style="display:flex;align-items:center;gap:6px;color:rgba(255,255,255,0.7);font-size:12px;white-space:nowrap;cursor:pointer;">
                <input type="checkbox" class="selector-attach-check" style="cursor:pointer;"
                    onchange="window._selectorAttach && window._selectorAttach(this.checked)">
                ${label}
            </label>
        </div>`;
}

function _renderWorldList(worlds, activeWorldId) {
    if (!worlds.length) return '<div style="color:var(--surface-hover);text-align:center;padding:20px;">No worlds available</div>';
    return worlds.map(w => {
        const isActive = activeWorldId === w.id;
        return _itemRow(
            `contextDisplay.selectWorld('${w.id}');document.querySelector('.selector-modal').remove();`,
            isActive,
            `<strong style="display:block;margin-bottom:4px;">${w.title || w.name}</strong><small style="color:var(--surface-hover);">${w.details || w.description || 'No description'}</small>`
        );
    }).join('');
}

// ---------------------------------------------------------------------------
// World Selector
// ---------------------------------------------------------------------------

ContextDisplay.prototype.openWorldSelector = function () {
    const allWorlds = window.WorldRepository.all();
    const activeWorldId = this.currentChatData.world?.id;
    const fe = window.FilterEngine;

    const render = (items) => {
        const container = document.querySelector('.selector-modal .selector-list-container');
        if (container) container.innerHTML = _renderWorldList(items, activeWorldId);
    };

    let currentItems = fe ? fe.processItems(allWorlds, { sortField: 'name', sortDirection: 'asc' }) : allWorlds;

    window._selectorFilter = (query) => {
        currentItems = fe ? fe.processItems(allWorlds, { query, searchFields: ['name','title','details','description'], sortField: window._selectorSortField || 'name', sortDirection: window._selectorSortDir || 'asc' }) : allWorlds;
        render(currentItems);
    };
    window._selectorSort = (val) => {
        const [field, dir] = val.split(':');
        window._selectorSortField = field;
        window._selectorSortDir = dir;
        const query = document.querySelector('.selector-modal .selector-search-input')?.value || '';
        currentItems = fe ? fe.processItems(allWorlds, { query, searchFields: ['name','title','details','description'], sortField: field, sortDirection: dir }) : allWorlds;
        render(currentItems);
    };

    _buildSelectorModal(_filterableModalShell(
        'Select a World', _renderWorldList(currentItems, activeWorldId),
        "document.querySelector('.selector-modal').remove();",
        '+ New World',
        "navigateTo('world.html?from=chat');document.querySelector('.selector-modal').remove();"
    ));
};

ContextDisplay.prototype.selectWorld = async function (worldId) {
    const worlds = window.WorldRepository.all();
    let world  = worlds.find(w => w.id === worldId);
    if (!world) return;

    // Rehydrate world images from IndexedDB before storing in state
    if (window.HOWLStorage) {
        try {
            world = await window.HOWLStorage.restoreWorldImages(world);
        } catch (error) {
            console.warn('Failed to rehydrate world images for selector:', error);
        }
    }

    console.log('[TRACE:selectWorld] start — world:', world?.id, 'name:', world?.title || world?.name);
    this.currentChatData.world     = world;
    this.currentChatData.scenario  = null;
    this.currentChatData.characters = [];
    this.currentChatData.lorebooks  = [];

    const scenariosCatalog  = window.ScenarioRepository.all();
    const charactersCatalog = (window.CharacterStorage?.all?.() || []);
    const lorebooksCatalog  = (window.LorebookRepo ? window.LorebookRepo.all() : []);

    const lorebookRefSets   = [];
    const characterIds      = new Set();
    const attachedChars     = [];

    const worldScenarioRefs = [
        ...(Array.isArray(world.scenarios)        ? world.scenarios        : []),
        ...(Array.isArray(world.howl?.scenarios)  ? world.howl.scenarios   : [])
    ];

    let firstScenario = null;
    if (worldScenarioRefs.length > 0) {
        const matched = this.resolveLinkedItems(scenariosCatalog, worldScenarioRefs);
        if (matched.length > 0) {
            firstScenario = matched[0];
            // Rehydrate scenario images from IndexedDB
            if (window.HOWLStorage) {
                try {
                    firstScenario = await window.HOWLStorage.restoreScenarioImages(firstScenario);
                } catch (error) {
                    console.warn('Failed to rehydrate scenario images for selector:', error);
                }
            }
            this.currentChatData.scenario = firstScenario;
            const refs = this.getScenarioLorebookRefs(firstScenario);
            if (refs.length > 0) lorebookRefSets.push(refs);
        }
    }

    const worldCharRefs = Array.isArray(world.characters) ? world.characters : [];
    this.resolveLinkedItems(charactersCatalog, worldCharRefs).forEach(c => {
        if (!characterIds.has(c.id)) { characterIds.add(c.id); attachedChars.push(c); }
    });

    if (firstScenario?.characters?.length > 0) {
        this.resolveLinkedItems(charactersCatalog, firstScenario.characters).forEach(c => {
            if (!characterIds.has(c.id)) { characterIds.add(c.id); attachedChars.push(c); }
        });
    }

    console.log('[TRACE:selectWorld] attachedChars count:', attachedChars.length, 'ids:', attachedChars.map(c => c.id));
    attachedChars.forEach(c => {
        if (Array.isArray(c.howl?.lorebooks)) lorebookRefSets.push(c.howl.lorebooks);
        if (Array.isArray(c.lorebooks))       lorebookRefSets.push(c.lorebooks);
    });

    const worldLbRefs = [
        ...(Array.isArray(world.howl?.lorebooks) ? world.howl.lorebooks : []),
        ...(Array.isArray(world.lorebooks)        ? world.lorebooks      : [])
    ];
    if (worldLbRefs.length > 0) lorebookRefSets.push(worldLbRefs);

    const mergedLorebooks = this.mergeLorebookRefs(lorebooksCatalog, lorebookRefSets);
    this.currentChatData.lorebooks  = mergedLorebooks;
    this.currentChatData.characters = attachedChars;

    // Character swap: use bridge batch update to avoid flicker.
    // First update bridge lorebooks so _push() propagates them.
    // Then batch-update characters (single panel refresh + backend sync).
    if (window.chatBridge) {
        // The bridge owns the state the panels render from. Writing
        // this.currentChatData.world directly was undone by the next _push(),
        // which rebuilt currentChatData from _state.world (still the previous
        // world) and re-rendered the World panel from it, so selecting a world
        // left the box unchanged and never persisted.
        await window.chatBridge.applyWorldChange(world);
        await window.chatBridge.applyScenarioChange(firstScenario || null);
        window.chatBridge.setLorebooks(mergedLorebooks || []);
        const bridgeChars = window.chatBridge.getState()?.characters || [];
        const attachedIds = new Set(attachedChars.map(c => c.id));
        const removeIds   = bridgeChars.filter(c => !attachedIds.has(c.id)).map(c => c.id);
        await window.chatBridge.applyCharacterBatchUpdate(removeIds, attachedChars);
    } else {
        this.saveChatData();
        this.updateWorldDisplay();
        await this.updateScenarioDisplay();
        this.updateCharacterPanel();
        this.updateLorebookDisplay();
    }

    // Sync world, scenario, and lorebooks to backend (characters already synced by bridge)
    await this.syncWorldDataToBackend(world);
    if (firstScenario)                             await this.syncScenarioDataToBackend(firstScenario);
    if (this.currentChatData.lorebooks.length)     await this.syncLorebooksDataToBackend(this.currentChatData.lorebooks);
    console.log('[TRACE:selectWorld] end — backend synced');

    if (typeof showNotification === 'function') {
        showNotification(`World loaded: ${world.title || world.name} (${attachedChars.length} characters, ${this.currentChatData.lorebooks.length} lorebooks)`, 'success');
    }
};

// ---------------------------------------------------------------------------
// Scenario Selector
// ---------------------------------------------------------------------------

function _renderScenarioList(scenarios, activeScenarioId) {
    if (!scenarios.length) return '<div style="color:var(--surface-hover);text-align:center;padding:20px;">No scenarios available</div>';
    return scenarios.map(s => {
        const isActive = activeScenarioId === s.id;
        return _itemRow(
            `contextDisplay.selectScenario('${s.id}');document.querySelector('.selector-modal').remove();`,
            isActive,
            `<strong style="display:block;margin-bottom:4px;">${s.title || s.name}</strong><small style="color:var(--surface-hover);">${s.description || 'No description'}</small>`
        );
    }).join('');
}

function _filterableModalShellWithAttach(title, listHtml, cancelOnclick, newLabel, newOnclick, attachLabel) {
    return `<div style="background:var(--shell-panel-bg);border-radius:12px;padding:24px;max-width:480px;width:90%;border:1px solid var(--surface-hover-strong);">
        <h2 style="color:#fff;margin-top:0;margin-bottom:16px;">${title}</h2>
        ${_filterBarWithAttachmentHtml(attachLabel)}
        <div class="selector-list-container" style="max-height:240px;overflow-y:auto;margin-bottom:16px;">${listHtml}</div>
        <div style="display:flex;gap:8px;">
            <button onclick="${cancelOnclick}" style="flex:1;padding:10px;background:var(--surface-hover-strong);border:1px solid var(--border-secondary);color:#fff;border-radius:6px;cursor:pointer;">Cancel</button>
            <button onclick="${newOnclick}"   style="flex:1;padding:10px;background:var(--success-bg);border:1px solid var(--success);color:#00b873;border-radius:6px;cursor:pointer;">${newLabel}</button>
        </div>
    </div>`;
}

ContextDisplay.prototype.openScenarioSelector = function () {
    const allScenarios = window.ScenarioRepository.all();
    const activeScenarioId = this.currentChatData.scenario?.id;
    const currentWorld = this.currentChatData.world;
    const fe = window.FilterEngine;

    const render = (items) => {
        const container = document.querySelector('.selector-modal .selector-list-container');
        if (container) container.innerHTML = _renderScenarioList(items, activeScenarioId);
    };

    let currentItems = fe ? fe.processItems(allScenarios, { sortField: 'name', sortDirection: 'asc' }) : allScenarios;
    let attachOnly = false;

    const apply = () => {
        const query = document.querySelector('.selector-modal .selector-search-input')?.value || '';
        const sortVal = document.querySelector('.selector-modal .selector-sort-select')?.value || 'name:asc';
        const [field, dir] = sortVal.split(':');
        const opts = { query, searchFields: ['name','title','description'], sortField: field, sortDirection: dir };
        if (attachOnly && currentWorld) {
            opts.attachTo = currentWorld;
            opts.attachField = 'scenarios';
        }
        currentItems = fe ? fe.processItems(allScenarios, opts) : allScenarios;
        render(currentItems);
    };

    window._selectorFilter = () => apply();
    window._selectorSort = () => apply();
    window._selectorAttach = (checked) => { attachOnly = checked; apply(); };

    _buildSelectorModal(_filterableModalShellWithAttach(
        'Select a Scenario', _renderScenarioList(currentItems, activeScenarioId),
        "document.querySelector('.selector-modal').remove();",
        '+ New Scenario',
        "navigateTo('scenario.html');document.querySelector('.selector-modal').remove();",
        'Only attached to world'
    ));
};

ContextDisplay.prototype.selectScenario = async function (scenarioId) {
    const scenarios = window.ScenarioRepository.all();
    let scenario  = scenarios.find(s => s.id === scenarioId);
    if (!scenario) return;

    // Rehydrate scenario images from IndexedDB before storing in state
    if (window.HOWLStorage) {
        try {
            scenario = await window.HOWLStorage.restoreScenarioImages(scenario);
        } catch (error) {
            console.warn('Failed to rehydrate scenario images for selector:', error);
        }
    }

    console.log('[TRACE:selectScenario] start — scenario:', scenario?.id, 'name:', scenario?.title || scenario?.name);
    this.currentChatData.scenario = scenario;

    // Resolve scenario characters from localStorage
    let scenarioChars = [];
    if (scenario.characters?.length > 0) {
        const all = (window.CharacterStorage?.all?.() || []);
        scenarioChars = all.filter(c => scenario.characters.some(sc => sc.id === c.id));
        if (scenarioChars.length > 0 && this.currentChatData.characters?.length > 0 && typeof dataSync !== 'undefined') {
            await dataSync.createMemoryLorebook(this.currentChatData.characters, scenarioChars);
        }
    }

    // Lorebook merge
    const lbRefs = this.getScenarioLorebookRefs(scenario);
    if (lbRefs.length > 0) {
        const catalog    = (window.LorebookRepo ? window.LorebookRepo.all() : []);
        const existRefs  = (this.currentChatData.lorebooks || []).map(lb => ({ id: lb.id, title: lb.title, name: lb.name, description: lb.description || '' }));
        this.currentChatData.lorebooks = this.mergeLorebookRefs(catalog, [existRefs, lbRefs]);
    }

    // Character swap: use bridge batch update to avoid flicker.
    // First update bridge lorebooks + scenario so _push() propagates them.
    // Then batch-update characters (single panel refresh + backend sync).
    if (window.chatBridge) {
        window.chatBridge.setLorebooks(this.currentChatData.lorebooks || []);
        await window.chatBridge.applyScenarioChange(scenario);
        const bridgeChars = window.chatBridge.getState()?.characters || [];
        const scenarioIds = new Set(scenarioChars.map(c => c.id));
        const removeIds   = bridgeChars.filter(c => !scenarioIds.has(c.id)).map(c => c.id);
        await window.chatBridge.applyCharacterBatchUpdate(removeIds, scenarioChars);
    } else {
        this.currentChatData.characters = scenarioChars;
        this.saveChatData();
        await this.updateScenarioDisplay();
        this.updateCharacterPanel();
        this.updateLorebookDisplay();
    }

    // Sync scenario and lorebooks to backend (characters already synced by bridge)
    await this.syncScenarioDataToBackend(scenario);
    if (this.currentChatData.lorebooks?.length) await this.syncLorebooksDataToBackend(this.currentChatData.lorebooks);

    console.log('[TRACE:selectScenario] end — chars:', (this.currentChatData.characters || []).length, 'lorebooks:', (this.currentChatData.lorebooks || []).length);
    if (typeof showNotification === 'function') showNotification(`Scenario set: ${scenario.title || scenario.name}`, 'success');
};

// ---------------------------------------------------------------------------
// Character Selector
// ---------------------------------------------------------------------------

function _renderCharacterList(characters, selectedIds) {
    if (!characters.length) return '<div style="color:var(--surface-hover);text-align:center;padding:20px;">No characters available</div>';
    return characters.map(c => {
        const isSel = selectedIds.has(c.id);
        const avatar = c.image
            ? `<img src="${c.image}" style="width:100%;height:100%;border-radius:50%;object-fit:cover;">`
            : '👤';
        const dot = isSel ? '<div style="width:12px;height:12px;border-radius:50%;background:#00b873;"></div>' : '';
        const bg     = isSel ? 'rgba(0,184,115,0.2)' : 'rgba(255,255,255,0.05)';
        const border = isSel ? 'rgba(0,184,115,0.3)' : 'rgba(255,255,255,0.1)';
        return `<div onclick="contextDisplay.toggleCharacter('${c.id}');"
            style="padding:12px;background:${bg};margin-bottom:8px;border-radius:6px;cursor:pointer;border:1px solid ${border};color:#fff;display:flex;align-items:center;gap:12px;">
            <div style="width:32px;height:32px;border-radius:50%;background:var(--border-secondary);display:flex;align-items:center;justify-content:center;flex-shrink:0;">${avatar}</div>
            <div style="flex:1;"><strong style="display:block;margin-bottom:2px;">${c.name || 'Unnamed'}</strong><small style="color:var(--surface-hover);">${c.role || 'No role'}</small></div>
            <div style="width:20px;height:20px;border-radius:50%;border:2px solid var(--surface-hover);display:flex;align-items:center;justify-content:center;">${dot}</div>
        </div>`;
    }).join('');
}

ContextDisplay.prototype.openCharacterSelector = function () {
    const allCharacters = (window.CharacterStorage?.all?.() || []);
    const selectedIds = new Set((this.currentChatData.characters || []).map(c => c.id));
    const currentWorld = this.currentChatData.world;
    const currentScenario = this.currentChatData.scenario;
    const fe = window.FilterEngine;

    const render = (items) => {
        const container = document.querySelector('.selector-modal .selector-list-container');
        if (container) container.innerHTML = _renderCharacterList(items, selectedIds);
    };

    let currentItems = fe ? fe.processItems(allCharacters, { sortField: 'name', sortDirection: 'asc' }) : allCharacters;
    let attachOnly = false;

    const apply = () => {
        const query = document.querySelector('.selector-modal .selector-search-input')?.value || '';
        const sortVal = document.querySelector('.selector-modal .selector-sort-select')?.value || 'name:asc';
        const [field, dir] = sortVal.split(':');
        const opts = { query, searchFields: ['name','role','description'], sortField: field, sortDirection: dir };
        if (attachOnly) {
            if (currentWorld) {
                opts.attachTo = currentWorld;
                opts.attachField = 'characters';
            } else if (currentScenario) {
                opts.attachTo = currentScenario;
                opts.attachField = 'characters';
            }
        }
        currentItems = fe ? fe.processItems(allCharacters, opts) : allCharacters;
        render(currentItems);
    };

    window._selectorFilter = () => apply();
    window._selectorSort = () => apply();
    window._selectorAttach = (checked) => { attachOnly = checked; apply(); };

    _buildSelectorModal(_filterableModalShellWithAttach(
        'Select Characters', _renderCharacterList(currentItems, selectedIds),
        "document.querySelector('.selector-modal').remove();",
        '+ New Character',
        "navigateTo('character.html');document.querySelector('.selector-modal').remove();",
        'Only attached'
    ));
};

ContextDisplay.prototype.toggleCharacter = async function (characterId) {
    console.log('[TRACE:toggleCharacter] id:', characterId);
    const all       = (window.CharacterStorage?.all?.() || []);
    const character = all.find(c => c.id === characterId);
    if (!character) return;

    const isAdding = !(this.currentChatData.characters || []).some(c => c.id === characterId);

    if (window.chatBridge) {
        if (isAdding) await window.chatBridge.applyCharacterAdd(character);
        else          await window.chatBridge.applyCharacterRemove(characterId);
    } else {
        if (!this.currentChatData.characters) this.currentChatData.characters = [];
        const idx = this.currentChatData.characters.findIndex(c => c.id === characterId);
        if (idx >= 0) this.currentChatData.characters.splice(idx, 1);
        else          this.currentChatData.characters.push(character);
    }

    this.saveChatData();
    await this.updateCharacterPanel();
    if ((this.currentChatData.characters || []).length) await this.syncCharactersDataToBackend(this.currentChatData.characters);

    const modal = document.querySelector('.selector-modal');
    if (modal) {
        modal.querySelectorAll('div[onclick*="toggleCharacter"]').forEach(div => {
            if (!div.getAttribute('onclick')?.includes(`'${characterId}'`)) return;
            const isNow = this.currentChatData.characters.some(c => c.id === characterId);
            div.style.background  = isNow ? 'var(--success-bg)' : 'var(--surface-hover)';
            div.style.borderColor = isNow ? 'var(--success)' : 'var(--surface-hover-strong)';
            const checkbox = div.querySelector('div:last-child');
            if (checkbox) checkbox.innerHTML = isNow ? '<div style="width:12px;height:12px;border-radius:50%;background:#00b873;"></div>' : '';
        });
    }
};

ContextDisplay.prototype.removeCharacterFromChat = async function (characterId) {
    if (!confirm('Remove this character from the chat?')) return;
    if (window.chatBridge) {
        await window.chatBridge.applyCharacterRemove(characterId);
    } else {
        this.currentChatData.characters = this.currentChatData.characters.filter(c => c.id !== characterId);
        if ((this.currentChatData.activeCharacterIds || []).includes(characterId)) {
            if (!this.currentChatData.removedCharacterIds) this.currentChatData.removedCharacterIds = [];
            if (!this.currentChatData.removedCharacterIds.includes(characterId)) {
                this.currentChatData.removedCharacterIds.push(characterId);
            }
        }
        this.currentChatData.addedCharacterIds = (this.currentChatData.addedCharacterIds || []).filter(id => id !== characterId);
    }
    this.saveChatData();
    await this.updateCharacterPanel();
    if (typeof showNotification === 'function') showNotification('Character removed from chat', 'success');
};
