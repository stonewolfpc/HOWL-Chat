/**
 * Sidebar Context — Scenario Section
 *
 * Scenario display update, scenario selector modal (with world-attachment
 * filtering and search), scenario selection, and auto-switch toggle.
 *
 * Depends on: sidebar-state.js, sidebar-session.js
 *
 * @module sidebar-context-scenario
 * @version 1.0.0
 */

async function updateScenarioDisplay() {
    const scenarioNameEl        = document.getElementById('scenario-name-display');
    const scenarioDescEl        = document.getElementById('scenario-description-display');
    const scenarioImageEl       = document.getElementById('scenario-image');
    const scenarioThumbnailEl   = document.getElementById('scenario-thumbnail');
    if (!scenarioNameEl) return;
    if (!window.sidebarState) return;

    const cd = window.sidebarState.currentChatData;

    if (cd.scenario) {
        scenarioNameEl.textContent = cd.scenario.title || cd.scenario.name;
        if (scenarioDescEl) {
            scenarioDescEl.textContent = cd.scenario.description || '';
        }

        if (scenarioImageEl) {
            // New multi-image array format: use thumbnail from images array
            let imageSrc = cd.scenario._resolvedImage || '';

            if (!imageSrc && cd.scenario.images && Array.isArray(cd.scenario.images) && cd.scenario.images.length > 0) {
                const thumbIdx = cd.scenario.thumbnailIndex || 0;
                const thumb = cd.scenario.images[thumbIdx] || cd.scenario.images[0];
                imageSrc = thumb.data || thumb.thumbnail || '';
            }

            // Legacy single-image format fallback
            if (!imageSrc && cd.scenario.image) {
                imageSrc = window.sidebarState.resolveImageSrc(cd.scenario.image) || '';
            }

            // Async fallback: resolve {id, stored:true} reference from IndexedDB (legacy)
            if (!imageSrc && cd.scenario.image && window.HOWLStorage) {
                try {
                    const ref       = cd.scenario.image;
                    const imageId   = (typeof ref === 'object' && ref.id) ? ref.id : null;
                    if (imageId) {
                        const imageData = await window.HOWLStorage.getImage(imageId, 'scenario');
                        if (imageData && imageData.data) {
                            imageSrc = imageData.data;
                            cd.scenario._resolvedImage = imageSrc;
                        } else {
                            console.warn('Scenario image not found in IndexedDB:', imageId);
                        }
                    }
                } catch (error) {
                    console.warn('Failed to restore scenario image:', error);
                }
            }

            const placeholder = scenarioThumbnailEl ? scenarioThumbnailEl.querySelector('.thumbnail-placeholder') : null;
            if (imageSrc) {
                scenarioImageEl.src = imageSrc;
                scenarioImageEl.style.display = 'block';
                if (placeholder) placeholder.style.display = 'none';
            } else {
                scenarioImageEl.style.display = 'none';
                if (placeholder) placeholder.style.display = '';
            }
        }
    } else {
        scenarioNameEl.textContent = 'No Scenario Selected';
        if (scenarioDescEl) scenarioDescEl.textContent = '';
        if (scenarioImageEl)     scenarioImageEl.style.display = 'none';
        const placeholder = scenarioThumbnailEl ? scenarioThumbnailEl.querySelector('.thumbnail-placeholder') : null;
        if (placeholder) placeholder.style.display = '';
    }
}

function openScenarioSelector() {
    const scenarios = JSON.parse(localStorage.getItem('howl-chat-scenarios') || '[]');
    const cd        = window.sidebarState.currentChatData;

    let attachedScenarios  = [];
    let availableScenarios = [];

    if (cd.world) {
        attachedScenarios  = scenarios.filter(s => {
            const inScenarios = cd.world.scenarios && cd.world.scenarios.some(ws => ws.id === s.id);
            const inHowl      = cd.world.howl && cd.world.howl.scenarios && cd.world.howl.scenarios.some(ws => ws.id === s.id);
            return inScenarios || inHowl;
        });
        availableScenarios = scenarios.filter(s => !attachedScenarios.some(a => a.id === s.id));
    } else {
        availableScenarios = scenarios;
    }

    const modal = document.createElement('div');
    modal.className   = 'scenario-selector-modal';
    modal.style.cssText = `position:fixed;top:0;left:0;right:0;bottom:0;background:rgba(0,0,0,0.7);display:flex;align-items:center;justify-content:center;z-index:9999;backdrop-filter:blur(4px);`;

    let content = `
        <div style="background:rgba(30,30,40,0.98);border-radius:12px;padding:24px;max-width:600px;max-height:80vh;border:1px solid rgba(255,255,255,0.1);overflow-y:auto;">
            <h2 style="color:#fff;margin-top:0;margin-bottom:16px;">Select Scenario</h2>
            <div style="margin-bottom:24px;">
                <h3 style="color:#00b873;margin-bottom:12px;font-size:14px;">Attached to ${cd.world ? cd.world.title || cd.world.name : 'World'}</h3>
                <div style="max-height:200px;overflow-y:auto;">
    `;

    if (attachedScenarios.length === 0) {
        content += `<div style="color:rgba(255,255,255,0.3);text-align:center;padding:16px;font-style:italic;">No scenarios attached to this world</div>`;
    } else {
        attachedScenarios.forEach(s => { content += buildScenarioItem(s, cd, '.scenario-selector-modal'); });
    }

    content += `
                </div>
            </div>
            <div style="margin-bottom:24px;">
                <h3 style="color:#ff9500;margin-bottom:12px;font-size:14px;">Available Scenarios</h3>
                <div style="margin-bottom:12px;">
                    <input type="text" id="scenario-search" placeholder="Search scenarios..." style="width:100%;padding:8px 12px;background:rgba(255,255,255,0.1);border:1px solid rgba(255,255,255,0.2);border-radius:6px;color:#fff;font-size:14px;" onkeyup="filterScenarios(this.value)">
                </div>
                <div style="max-height:200px;overflow-y:auto;" id="available-scenarios-list">
    `;

    if (availableScenarios.length === 0) {
        content += `<div style="color:rgba(255,255,255,0.3);text-align:center;padding:16px;font-style:italic;">${cd.world ? 'No additional scenarios available' : 'No scenarios available'}</div>`;
    } else {
        availableScenarios.forEach(s => { content += buildScenarioItem(s, cd, '.scenario-selector-modal', true); });
    }

    content += `
                </div>
            </div>
            <div style="display:flex;gap:8px;">
                <button onclick="document.querySelector('.scenario-selector-modal').remove();" style="flex:1;padding:10px;background:rgba(255,255,255,0.1);border:1px solid rgba(255,255,255,0.2);color:#fff;border-radius:6px;cursor:pointer;">Cancel</button>
                <button onclick="navigateTo('scenario.html'); document.querySelector('.scenario-selector-modal').remove();" style="flex:1;padding:10px;background:rgba(0,184,115,0.2);border:1px solid rgba(0,184,115,0.3);color:#00b873;border-radius:6px;cursor:pointer;">+ New Scenario</button>
            </div>
        </div>
    `;

    modal.innerHTML = content;
    document.body.appendChild(modal);
    setTimeout(() => { const si = document.getElementById('scenario-search'); if (si) si.focus(); }, 100);
}

function buildScenarioItem(scenario, cd, modalClass, isAvailable = false) {
    const isActive = cd.scenario && cd.scenario.id === scenario.id;
    const dataAttr = isAvailable ? `class="scenario-item" data-scenario-name="${(scenario.title || scenario.name).toLowerCase()}"` : '';
    return `
        <div ${dataAttr} onclick="selectScenario('${scenario.id}'); document.querySelector('${modalClass}').remove();" style="
            padding:12px;background:${isActive ? 'rgba(0,184,115,0.2)' : 'rgba(255,255,255,0.05)'};
            margin-bottom:8px;border-radius:6px;cursor:pointer;position:relative;
            border:1px solid ${isActive ? 'rgba(0,184,115,0.3)' : 'rgba(255,255,255,0.1)'};color:#fff;transition:all 0.2s ease;">
            ${isActive ? '<span style="position:absolute;top:8px;right:8px;color:#00b873;">✓</span>' : ''}
            <strong style="display:block;margin-bottom:4px;">${scenario.title || scenario.name}</strong>
            <small style="color:rgba(255,255,255,0.5);">${scenario.description || 'No description'}</small>
        </div>
    `;
}

function filterScenarios(searchTerm) {
    const term = searchTerm.toLowerCase();
    document.querySelectorAll('.scenario-item').forEach(item => {
        item.style.display = item.dataset.scenarioName.includes(term) ? 'block' : 'none';
    });
}

async function selectScenario(scenarioId) {
    const scenarios = JSON.parse(localStorage.getItem('howl-chat-scenarios') || '[]');
    const scenario  = scenarios.find(s => s.id === scenarioId);
    if (!scenario) return;

    const cd           = window.sidebarState.currentChatData;
    const prevScenario = cd.scenario;

    console.log('[selectScenario] using chatBridge:', !!window.chatBridge, 'scenario:', scenario.id, 'chars:', (scenario.characters || []).length, 'books:', (scenario.lorebooks || []).length);
    if (window.chatBridge) {
        await window.chatBridge.applyScenarioChange(scenario);
    } else {
        cd.scenario = scenario;
        saveChatData();
        await updateScenarioDisplay();
        updateLorebookDisplay();
    }

    if (prevScenario && prevScenario.id !== scenarioId) {
        showNotification(`Hot-swapped scenario: "${prevScenario.title || prevScenario.name}" → "${scenario.title || scenario.name}"`, 'info');
    } else {
        showNotification(`Scenario set: ${scenario.title || scenario.name}`, 'success');
    }

    if (window.go?.main?.App?.AddScenarioToChat) {
        try { window.go.main.App.AddScenarioToChat(scenario); }
        catch (error) { console.warn('Backend scenario update failed:', error); }
    }
}

function toggleScenarioAutoSwitch() {
    const checkbox = document.getElementById('scenario-auto-switch');
    console.log('Scenario auto-switch:', checkbox.checked);
}
