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
        'background:rgba(0,0,0,0.7)', 'display:flex', 'align-items:center',
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
        onmouseover="this.style.background='rgba(255,255,255,0.1)';this.style.borderColor='rgba(255,255,255,0.2)';"
        onmouseout="this.style.background='${bgOut}';this.style.borderColor='${bdOut}';">${bodyHtml}</div>`;
}

function _modalShell(title, listHtml, cancelOnclick, newLabel, newOnclick) {
    return `<div style="background:rgba(30,30,40,0.98);border-radius:12px;padding:24px;max-width:480px;width:90%;border:1px solid rgba(255,255,255,0.1);">
        <h2 style="color:#fff;margin-top:0;margin-bottom:16px;">${title}</h2>
        <div style="max-height:300px;overflow-y:auto;margin-bottom:16px;">${listHtml}</div>
        <div style="display:flex;gap:8px;">
            <button onclick="${cancelOnclick}" style="flex:1;padding:10px;background:rgba(255,255,255,0.1);border:1px solid rgba(255,255,255,0.2);color:#fff;border-radius:6px;cursor:pointer;">Cancel</button>
            <button onclick="${newOnclick}"   style="flex:1;padding:10px;background:rgba(0,184,115,0.2);border:1px solid rgba(0,184,115,0.3);color:#00b873;border-radius:6px;cursor:pointer;">${newLabel}</button>
        </div>
    </div>`;
}

// ---------------------------------------------------------------------------
// World Selector
// ---------------------------------------------------------------------------

ContextDisplay.prototype.openWorldSelector = function () {
    const worlds = JSON.parse(localStorage.getItem('howl-worlds') || '[]');
    const listHtml = worlds.length === 0
        ? '<div style="color:rgba(255,255,255,0.5);text-align:center;padding:20px;">No worlds available</div>'
        : worlds.map(w => {
            const isActive = this.currentChatData.world?.id === w.id;
            return _itemRow(
                `contextDisplay.selectWorld('${w.id}');document.querySelector('.selector-modal').remove();`,
                isActive,
                `<strong style="display:block;margin-bottom:4px;">${w.title || w.name}</strong><small style="color:rgba(255,255,255,0.5);">${w.details || w.description || 'No description'}</small>`
            );
        }).join('');

    _buildSelectorModal(_modalShell(
        'Select a World', listHtml,
        "document.querySelector('.selector-modal').remove();",
        '+ New World',
        "navigateTo('world.html');document.querySelector('.selector-modal').remove();"
    ));
};

ContextDisplay.prototype.selectWorld = async function (worldId) {
    const worlds = JSON.parse(localStorage.getItem('howl-worlds') || '[]');
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

    this.currentChatData.world     = world;
    this.currentChatData.scenario  = null;
    this.currentChatData.characters = [];
    this.currentChatData.lorebooks  = [];

    const scenariosCatalog  = JSON.parse(localStorage.getItem('howl-chat-scenarios')  || '[]');
    const charactersCatalog = JSON.parse(localStorage.getItem('howl-chat-characters') || '[]');
    const lorebooksCatalog  = JSON.parse(localStorage.getItem('howl-chat-lorebooks')  || '[]');

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

    attachedChars.forEach(c => {
        if (Array.isArray(c.howl?.lorebooks)) lorebookRefSets.push(c.howl.lorebooks);
        if (Array.isArray(c.lorebooks))       lorebookRefSets.push(c.lorebooks);
    });

    const worldLbRefs = [
        ...(Array.isArray(world.howl?.lorebooks) ? world.howl.lorebooks : []),
        ...(Array.isArray(world.lorebooks)        ? world.lorebooks      : [])
    ];
    if (worldLbRefs.length > 0) lorebookRefSets.push(worldLbRefs);

    this.currentChatData.lorebooks  = this.mergeLorebookRefs(lorebooksCatalog, lorebookRefSets);
    this.currentChatData.characters = attachedChars;

    this.saveChatData();
    this.updateWorldDisplay();
    await this.updateScenarioDisplay();
    this.updateCharacterPanel();
    this.updateLorebookDisplay();

    await this.syncWorldDataToBackend(world);
    if (firstScenario)                             await this.syncScenarioDataToBackend(firstScenario);
    if (this.currentChatData.characters.length)    await this.syncCharactersDataToBackend(this.currentChatData.characters);
    if (this.currentChatData.lorebooks.length)     await this.syncLorebooksDataToBackend(this.currentChatData.lorebooks);

    if (typeof showNotification === 'function') {
        showNotification(`World loaded: ${world.title || world.name} (${attachedChars.length} characters, ${this.currentChatData.lorebooks.length} lorebooks)`, 'success');
    }
};

// ---------------------------------------------------------------------------
// Scenario Selector
// ---------------------------------------------------------------------------

ContextDisplay.prototype.openScenarioSelector = function () {
    const scenarios = JSON.parse(localStorage.getItem('howl-chat-scenarios') || '[]');
    const listHtml  = scenarios.length === 0
        ? '<div style="color:rgba(255,255,255,0.5);text-align:center;padding:20px;">No scenarios available</div>'
        : scenarios.map(s => {
            const isActive = this.currentChatData.scenario?.id === s.id;
            return _itemRow(
                `contextDisplay.selectScenario('${s.id}');document.querySelector('.selector-modal').remove();`,
                isActive,
                `<strong style="display:block;margin-bottom:4px;">${s.title || s.name}</strong><small style="color:rgba(255,255,255,0.5);">${s.description || 'No description'}</small>`
            );
        }).join('');

    _buildSelectorModal(_modalShell(
        'Select a Scenario', listHtml,
        "document.querySelector('.selector-modal').remove();",
        '+ New Scenario',
        "navigateTo('scenario.html');document.querySelector('.selector-modal').remove();"
    ));
};

ContextDisplay.prototype.selectScenario = async function (scenarioId) {
    const scenarios = JSON.parse(localStorage.getItem('howl-chat-scenarios') || '[]');
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

    this.currentChatData.scenario = scenario;

    if (scenario.characters?.length > 0) {
        const all = JSON.parse(localStorage.getItem('howl-chat-characters') || '[]');
        const scenarioChars = all.filter(c => scenario.characters.some(sc => sc.id === c.id));
        if (scenarioChars.length > 0) {
            if (this.currentChatData.characters?.length > 0 && typeof dataSync !== 'undefined') {
                await dataSync.createMemoryLorebook(this.currentChatData.characters, scenarioChars);
            }
            this.currentChatData.characters = scenarioChars;
        }
    }

    const lbRefs = this.getScenarioLorebookRefs(scenario);
    if (lbRefs.length > 0) {
        const catalog    = JSON.parse(localStorage.getItem('howl-chat-lorebooks') || '[]');
        const existRefs  = (this.currentChatData.lorebooks || []).map(lb => ({ id: lb.id, title: lb.title, name: lb.name, description: lb.description || '' }));
        this.currentChatData.lorebooks = this.mergeLorebookRefs(catalog, [existRefs, lbRefs]);
    }

    this.saveChatData();
    await this.updateScenarioDisplay();
    this.updateCharacterPanel();
    this.updateLorebookDisplay();

    await this.syncScenarioDataToBackend(scenario);
    if (this.currentChatData.characters?.length) await this.syncCharactersDataToBackend(this.currentChatData.characters);
    if (this.currentChatData.lorebooks?.length)  await this.syncLorebooksDataToBackend(this.currentChatData.lorebooks);

    if (typeof showNotification === 'function') showNotification(`Scenario set: ${scenario.title || scenario.name}`, 'success');
};

// ---------------------------------------------------------------------------
// Character Selector
// ---------------------------------------------------------------------------

ContextDisplay.prototype.openCharacterSelector = function () {
    const characters = JSON.parse(localStorage.getItem('howl-chat-characters') || '[]');
    const listHtml   = characters.length === 0
        ? '<div style="color:rgba(255,255,255,0.5);text-align:center;padding:20px;">No characters available</div>'
        : characters.map(c => {
            const isSel = this.currentChatData.characters?.some(ch => ch.id === c.id);
            const avatar = c.image
                ? `<img src="${c.image}" style="width:100%;height:100%;border-radius:50%;object-fit:cover;">`
                : '👤';
            const dot = isSel ? '<div style="width:12px;height:12px;border-radius:50%;background:#00b873;"></div>' : '';
            const bg     = isSel ? 'rgba(0,184,115,0.2)' : 'rgba(255,255,255,0.05)';
            const border = isSel ? 'rgba(0,184,115,0.3)' : 'rgba(255,255,255,0.1)';
            return `<div onclick="contextDisplay.toggleCharacter('${c.id}');"
                style="padding:12px;background:${bg};margin-bottom:8px;border-radius:6px;cursor:pointer;border:1px solid ${border};color:#fff;display:flex;align-items:center;gap:12px;">
                <div style="width:32px;height:32px;border-radius:50%;background:rgba(255,255,255,0.2);display:flex;align-items:center;justify-content:center;flex-shrink:0;">${avatar}</div>
                <div style="flex:1;"><strong style="display:block;margin-bottom:2px;">${c.name || 'Unnamed'}</strong><small style="color:rgba(255,255,255,0.5);">${c.role || 'No role'}</small></div>
                <div style="width:20px;height:20px;border-radius:50%;border:2px solid rgba(255,255,255,0.3);display:flex;align-items:center;justify-content:center;">${dot}</div>
            </div>`;
        }).join('');

    _buildSelectorModal(_modalShell(
        'Select Characters', listHtml,
        "document.querySelector('.selector-modal').remove();",
        '+ New Character',
        "navigateTo('character.html');document.querySelector('.selector-modal').remove();"
    ));
};

ContextDisplay.prototype.toggleCharacter = async function (characterId) {
    const all       = JSON.parse(localStorage.getItem('howl-chat-characters') || '[]');
    const character = all.find(c => c.id === characterId);
    if (!character) return;

    if (!this.currentChatData.characters) this.currentChatData.characters = [];
    const idx = this.currentChatData.characters.findIndex(c => c.id === characterId);

    if (idx >= 0) {
        this.currentChatData.characters.splice(idx, 1);
    } else {
        this.currentChatData.characters.push(character);
    }

    this.saveChatData();
    await this.updateCharacterPanel();
    if (this.currentChatData.characters.length) await this.syncCharactersDataToBackend(this.currentChatData.characters);

    const modal = document.querySelector('.selector-modal');
    if (modal) {
        modal.querySelectorAll('div[onclick*="toggleCharacter"]').forEach(div => {
            if (!div.getAttribute('onclick')?.includes(`'${characterId}'`)) return;
            const isNow = this.currentChatData.characters.some(c => c.id === characterId);
            div.style.background  = isNow ? 'rgba(0,184,115,0.2)' : 'rgba(255,255,255,0.05)';
            div.style.borderColor = isNow ? 'rgba(0,184,115,0.3)' : 'rgba(255,255,255,0.1)';
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
    }
    if (!this.currentChatData.removedCharacterIds) this.currentChatData.removedCharacterIds = [];
    if (!this.currentChatData.removedCharacterIds.includes(characterId)) {
        this.currentChatData.removedCharacterIds.push(characterId);
    }
    this.saveChatData();
    await this.updateCharacterPanel();
    if (typeof showNotification === 'function') showNotification('Character removed from chat', 'success');
};
