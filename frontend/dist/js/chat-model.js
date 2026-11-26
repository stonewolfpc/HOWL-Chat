/**
 * Chat Model Management
 *
 * Load/unload the LLM model via the Wails backend,
 * read model metadata to update UI sliders, and
 * display the currently loaded model name.
 *
 * @module chat-model
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Model Display
// ---------------------------------------------------------------------------

function updateModelDisplay(modelName) {
    const display = document.getElementById('model-name-display');
    if (display) display.textContent = modelName;
}

// ---------------------------------------------------------------------------
// Model Loading Guard
// ---------------------------------------------------------------------------

let _modelLoading = false;
let _modelLoadTimeout = null;
const MODEL_LOAD_TIMEOUT_MS = 60000; // 60s safety reset

function isModelLoading() {
    return _modelLoading;
}

function setModelLoading(value) {
    _modelLoading = value;
    if (value) {
        // Safety: auto-reset if load hangs
        if (_modelLoadTimeout) clearTimeout(_modelLoadTimeout);
        _modelLoadTimeout = setTimeout(() => {
            console.warn('[LLM] Model load timeout — resetting loading flag');
            _modelLoading = false;
        }, MODEL_LOAD_TIMEOUT_MS);
    } else {
        if (_modelLoadTimeout) {
            clearTimeout(_modelLoadTimeout);
            _modelLoadTimeout = null;
        }
    }
}

// ---------------------------------------------------------------------------
// Custom Model Dropdown Helpers
// ---------------------------------------------------------------------------

function _getModelSelect() {
    return document.getElementById('llm-model-select');
}

function _getModelSelectTrigger() {
    const sel = _getModelSelect();
    return sel?.querySelector('.model-select-trigger');
}

function _getModelSelectOptionsContainer() {
    const sel = _getModelSelect();
    return sel?.querySelector('.model-select-options');
}

function _getModelSelectDropdown() {
    const sel = _getModelSelect();
    return sel?.querySelector('.model-select-dropdown');
}

function _getModelSelectValue() {
    const trigger = _getModelSelectTrigger();
    return trigger?.dataset.value || '';
}

function _setModelSelectValue(path, text) {
    const trigger = _getModelSelectTrigger();
    if (!trigger) return;
    trigger.dataset.value = path || '';
    trigger.textContent = text || (path ? path : '— select model —');
    trigger.title = path ? `${text || path}\n${path}` : 'Select a registered LLM model';
    const opts = _getModelSelectOptionsContainer();
    if (!opts) return;
    opts.querySelectorAll('.model-select-option').forEach(opt => {
        opt.classList.toggle('active', opt.dataset.value === path);
    });
}

function _clearModelSelectOptions() {
    const opts = _getModelSelectOptionsContainer();
    if (opts) opts.innerHTML = '';
}

function _addModelSelectOption(path, text, title) {
    const opts = _getModelSelectOptionsContainer();
    if (!opts) return;
    const div = document.createElement('div');
    div.className = 'model-select-option';
    div.dataset.value = path;
    div.textContent = text;
    div.title = title || path || text;
    div.addEventListener('click', () => {
        _setModelSelectValue(path, text);
        _closeModelSelect();
    });
    opts.appendChild(div);
}

function _addModelSelectEmpty(text) {
    const opts = _getModelSelectOptionsContainer();
    if (!opts) return;
    const div = document.createElement('div');
    div.className = 'model-select-empty';
    div.textContent = text;
    opts.appendChild(div);
}

function _modelSelectHasOption(path) {
    const opts = _getModelSelectOptionsContainer();
    if (!opts) return false;
    return [...opts.children].some(c => c.dataset.value === path);
}

function _openModelSelect() {
    const dd = _getModelSelectDropdown();
    if (dd) dd.classList.add('open');
}

function _closeModelSelect() {
    const dd = _getModelSelectDropdown();
    if (dd) dd.classList.remove('open');
}

function _toggleModelSelect() {
    const dd = _getModelSelectDropdown();
    if (!dd) return;
    dd.classList.toggle('open');
}

function _setupModelSelect() {
    const sel = _getModelSelect();
    if (!sel || sel.dataset.customSetup) return;
    sel.dataset.customSetup = 'true';

    const trigger = sel.querySelector('.model-select-trigger');
    if (trigger) {
        trigger.addEventListener('click', (e) => {
            e.stopPropagation();
            _toggleModelSelect();
        });
    }

    document.addEventListener('click', (e) => {
        if (!sel.contains(e.target)) _closeModelSelect();
    });
}

// ---------------------------------------------------------------------------
// Model Status Check
// ---------------------------------------------------------------------------

// ensureChatModelReadyForSend is the demand-driven activation point. Opening the
// chat page used to load (and possibly reclaim) a model before the user asked
// for anything; the send itself now owns loading, handoff and readiness, which
// is why a message can no longer race the engine it just started.
async function ensureChatModelReadyForSend() {
    const app = window.go?.main?.App;
    if (!app?.EnsureChatModelReady) return true;
    const modelPath = localStorage.getItem('howl-last-model-path') || '';
    setModelLoading(true);
    try {
        await app.EnsureChatModelReady(modelPath);
        const name = await app.GetLoadedModelName().catch(() => '');
        if (name) updateModelDisplay(name);
        return true;
    } catch (err) {
        const msg = (err && err.message) ? err.message : String(err);
        console.error('[Chat] demand load failed:', msg);
        updateModelDisplay('No model loaded');
        if (typeof showNotification === 'function') {
            showNotification('Could not start the chat model: ' + msg, 'error', 10000);
        } else {
            alert('Could not start the chat model: ' + msg);
        }
        return false;
    } finally {
        setModelLoading(false);
    }
}

function checkModelStatus() {
    if (!window.go || !window.go.main || !window.go.main.App) {
        console.log('Wails runtime not available, skipping model status check');
        return;
    }
    if (isModelLoading()) {
        console.log('Model status check skipped: load in progress');
        return;
    }
    try {
        window.go.main.App.IsModelLoaded().then(isLoaded => {
            if (isLoaded) {
                window.go.main.App.GetLoadedModelName().then(modelName => {
                    updateModelDisplay(modelName);
                }).catch(error => {
                    console.error('Error getting model name:', error);
                    updateModelDisplay('Error getting model name');
                });
            } else {
                // Navigation only shows the chat UI. Loading is the send's job,
                // and the explicit Load control still preloads on request.
                updateModelDisplay('No model loaded');
            }
        }).catch(error => {
            console.error('Error checking model status:', error);
            updateModelDisplay('Error checking status');
        });
    } catch (error) {
        console.error('Error checking model status:', error);
        updateModelDisplay('No model loaded');
    }
}

// ---------------------------------------------------------------------------
// Registry-based LLM dropdown
// ---------------------------------------------------------------------------

// A chat picker must only offer files llama-server can actually load. The
// registry also stores vision projectors and, from older scans, non-GGUF
// assets; those used to appear as selectable "models" and then fail on Load.
function _isChatModel(m) {
    if (!m || !m.path || !/\.(gguf|bin)$/i.test(m.path)) return false;
    const file = String(m.path.split(/[\\/]/).pop() || '').toLowerCase();
    const base = String(m.name || file).toLowerCase();
    return !/^mmproj(?:[-_.]|$)/.test(base) && !/^mmproj(?:[-_.]|$)/.test(file);
}

function _chatModels(models) {
    return (models || [])
        .filter(_isChatModel)
        .sort((a, b) => String(a.name || '').localeCompare(
            String(b.name || ''), undefined, { numeric: true, sensitivity: 'base' }
        ));
}

async function populateLLMDropdown() {
    const app = window.go?.main?.App;
    _setupModelSelect();
    if (!app?.GetLLMModels) return;
    try {
        const models = _chatModels(await app.GetLLMModels());
        const prev = _getModelSelectValue();
        _clearModelSelectOptions();
        if (models.length === 0) {
            _addModelSelectEmpty('No chat models registered');
        }
        models.forEach(m => {
            const sizeTxt = m.size_mb > 0 ? ` (${(m.size_mb / 1024).toFixed(1)} GB)` : '';
            _addModelSelectOption(m.path, `${m.name}${sizeTxt}`, m.path);
        });
        if (prev && _modelSelectHasOption(prev)) {
            const opt = [..._getModelSelectOptionsContainer().children].find(c => c.dataset.value === prev);
            _setModelSelectValue(prev, opt?.textContent || prev);
        }
        const saved = localStorage.getItem('howl-last-model-path');
        if (!_getModelSelectValue() && saved && _modelSelectHasOption(saved)) {
            const opt = [..._getModelSelectOptionsContainer().children].find(c => c.dataset.value === saved);
            _setModelSelectValue(saved, opt?.textContent || saved);
        }
    } catch (err) {
        console.warn('[LLM] Failed to populate model dropdown:', err);
    }
}

async function loadSelectedModel() {
    const path = _getModelSelectValue();
    if (!path) { alert('Please select a model first.'); return; }
    await _loadModelByPath(path);
}

async function browseLLMModel() {
    const app = window.go?.main?.App;
    if (!app?.OpenModelFilePicker) return;
    try {
        const path = await app.OpenModelFilePicker();
        if (!path) return;
        await _loadModelByPath(path);
        await populateLLMDropdown();
        if (_modelSelectHasOption(path)) {
            const opt = [..._getModelSelectOptionsContainer().children].find(c => c.dataset.value === path);
            _setModelSelectValue(path, opt?.textContent || path);
        }
    } catch (err) {
        updateModelDisplay('Error loading model');
        console.error('Failed to browse model:', err);
    }
}

async function scanLLMDirectory() {
    const app = window.go?.main?.App;
    if (!app) return;
    try {
        // The themed locator is the product path; the native dialog only backs
        // a build where the locator binding is unavailable.
        const dir = typeof openDirectoryLocator === 'function'
            ? await openDirectoryLocator('')
            : await app.OpenDirectoryPicker();
        if (!dir) {
            if (typeof showNotification === 'function') showNotification('No folder selected.', 'error');
            return;
        }
        updateModelDisplay('Scanning…');
        await app.ScanLLMDir(dir);
        await populateLLMDropdown();
        updateModelDisplay('Scan complete');
    } catch (err) {
        console.error('[LLM] Directory scan failed:', err);
        updateModelDisplay('Scan failed');
    }
}

// VRAM arbitration before a chat model load: if a loaded image model has to go
// first, say so and let the user decide. Returns false when the user cancels.
async function requestVRAMArbitration(kind, path) {
    const app = window.go?.main?.App;
    if (!app?.CheckModelFit) {
        alert('This build cannot verify VRAM fit; refusing to start a model blindly.');
        return false;
    }
    try {
        let fit = await app.CheckModelFit(kind, path || '');
        if (fit.can_fit) return true;
        if (!fit.blocking || fit.blocking.length === 0) {
            alert(fit.message);
            return false;
        }
        // The backend unloads only on this answer, so hand it the user's
        // decision rather than a constant.
        const approved = confirm(fit.message);
        if (!approved) return false;
        updateModelDisplay('Unloading the other model…');
        fit = await app.FreeVRAMFor(kind, path || '', approved);
        if (!fit.can_fit) { alert(fit.message); updateModelDisplay('No model loaded'); return false; }
        return true;
    } catch (err) {
        // A failed fit check is not permission to start. The backend authority
        // will still refuse if the allocation is unsafe; do not silently claim
        // the model fits when the check could not run.
        console.error('[VRAM] arbitration check failed:', err);
        alert('Could not verify available VRAM: ' + err);
        return false;
    }
}

async function _loadModelByPath(modelPath) {
    const app = window.go?.main?.App;
    if (!app?.LoadModel) return;
    if (isModelLoading()) {
        console.log('Model load skipped: already loading a model');
        return;
    }
    // Validate path before calling backend
    if (!modelPath || modelPath === '.' || modelPath === '..') {
        console.error('Invalid model path:', modelPath);
        alert('Please select a valid model file, not a directory.');
        return;
    }
    if (!modelPath.toLowerCase().endsWith('.gguf') && !modelPath.toLowerCase().endsWith('.bin')) {
        console.error('Invalid model file extension:', modelPath);
        alert('Please select a .gguf or .bin model file.');
        return;
    }
    // Manual selection is a model start too: ask before the backend stops the
    // model that is currently holding the VRAM.
    if (!(await requestVRAMArbitration('llm', modelPath))) {
        console.log('[VRAM] chat model load cancelled by user');
        updateModelDisplay('Load cancelled');
        return;
    }
    setModelLoading(true);
    updateModelDisplay('Loading...');
    try {
        try {
            const metadata = await app.LoadModelMetadata(modelPath);
            if (metadata) {
                const gpuSlider = document.getElementById('gpu-layers');
                const gpuVal    = document.getElementById('gpu-layers-value');
                if (gpuSlider && gpuVal && metadata.LayerCount > 0) {
                    gpuSlider.max = metadata.LayerCount;
                    gpuVal.max    = metadata.LayerCount;
                }
                const ctxSlider = document.getElementById('max-context');
                const ctxVal    = document.getElementById('max-context-value');
                if (ctxSlider && ctxVal && metadata.MaxContext > 0) {
                    ctxSlider.max = metadata.MaxContext;
                    ctxVal.max    = metadata.MaxContext;
                    if (parseInt(ctxSlider.value) > metadata.UsableContext) {
                        ctxSlider.value = metadata.UsableContext;
                        ctxVal.value    = metadata.UsableContext;
                    }
                }
            }
        } catch (e) { console.warn('Metadata read failed:', e); }

        await app.LoadModel(modelPath);
        const modelName = modelPath.split(/[\\/]/).pop();
        updateModelDisplay(modelName);
        localStorage.setItem('howl-last-model-path', modelPath);
    } catch (err) {
        updateModelDisplay('Error loading model');
        console.error('Failed to load model:', err);
        alert('Failed to load model: ' + err);
    } finally {
        setModelLoading(false);
    }
}

// ---------------------------------------------------------------------------
// Unload Model
// ---------------------------------------------------------------------------

async function unloadModel() {
    if (!confirm('⚠️ AGGRESSIVE UNLOAD ⚠️\n\nThis will forcefully unload the current model from memory.\nAll unsaved context will be lost.\n\nAre you sure you want to proceed?')) return;
    setModelLoading(false);
    try {
        await window.go.main.App.UnloadModel();
        updateModelDisplay('No model loaded');
        alert('🔴 MODEL AGGRESSIVELY UNLOADED\n\nGPU cache cleared.\nMemory freed.\nModel process terminated.');
    } catch (error) {
        console.error('Error unloading model:', error);
        alert(`Failed to unload model: ${error}`);
    }
}

// ---------------------------------------------------------------------------
// Init on DOM Ready
// ---------------------------------------------------------------------------

document.addEventListener('DOMContentLoaded', () => {
    checkModelStatus();
    populateLLMDropdown();
});

// ---------------------------------------------------------------------------
// Model Browser
// ---------------------------------------------------------------------------

let _browserModels = [];
let _browserSelectedPath = null;
let _browserSelectedDir = null;

function openModelBrowser() {
    const modal = document.getElementById('model-browser-modal');
    if (modal) modal.classList.add('show');
    document.body.style.overflow = 'hidden';
    populateModelBrowser();
}

function closeModelBrowser() {
    const modal = document.getElementById('model-browser-modal');
    if (modal) modal.classList.remove('show');
    document.body.style.overflow = '';
    _browserSelectedPath = null;
    _browserSelectedDir = null;
}

// ---------------------------------------------------------------------------
// Model Browser Directory Dropdown
// ---------------------------------------------------------------------------

function _getDirsSelect() {
    return document.getElementById('model-browser-dirs-select');
}

function _getDirsTrigger() {
    const sel = _getDirsSelect();
    return sel?.querySelector('.mb-dirs-trigger');
}

function _getDirsOptionsContainer() {
    return document.getElementById('model-browser-dirs-options');
}

function _getDirsDropdown() {
    const sel = _getDirsSelect();
    return sel?.querySelector('.mb-dirs-dropdown');
}

function _toggleDirsDropdown() {
    const dd = _getDirsDropdown();
    if (!dd) return;
    dd.classList.toggle('open');
}

function _closeDirsDropdown() {
    const dd = _getDirsDropdown();
    if (dd) dd.classList.remove('open');
}

function _setupDirsDropdown() {
    const sel = _getDirsSelect();
    if (!sel || sel.dataset.customSetup) return;
    sel.dataset.customSetup = 'true';

    const trigger = sel.querySelector('.mb-dirs-trigger');
    if (trigger) {
        trigger.addEventListener('click', (e) => {
            e.stopPropagation();
            _toggleDirsDropdown();
        });
    }

    document.addEventListener('click', (e) => {
        if (!sel.contains(e.target)) _closeDirsDropdown();
    });
}

function _setDirsTriggerText(text) {
    const trigger = _getDirsTrigger();
    if (trigger) trigger.textContent = text || 'All directories';
}

function _clearDirOptions() {
    const opts = _getDirsOptionsContainer();
    if (opts) opts.innerHTML = '';
}

function _addDirOption(path, isActive) {
    const opts = _getDirsOptionsContainer();
    if (!opts) return;
    const div = document.createElement('div');
    div.className = 'mb-dirs-option' + (isActive ? ' active' : '');
    div.title = path;
    div.innerHTML = `
        <span class="mb-dirs-option-text">${escapeHtml(path)}</span>
        <span class="mb-dirs-remove" title="Remove directory">&times;</span>
    `;
    div.querySelector('.mb-dirs-option-text').addEventListener('click', () => {
        _browserSelectedDir = path;
        _setDirsTriggerText(path);
        _closeDirsDropdown();
        filterModelBrowser();
    });
    div.querySelector('.mb-dirs-remove').addEventListener('click', (e) => {
        e.stopPropagation();
        removeModelDirectory(path);
    });
    opts.appendChild(div);
}

async function populateModelBrowserDirs() {
    _setupDirsDropdown();
    const app = window.go?.main?.App;
    if (!app?.GetRegisteredLLMDirs) return;
    try {
        const dirs = await app.GetRegisteredLLMDirs();
        _clearDirOptions();
        // "All directories" option
        const opts = _getDirsOptionsContainer();
        if (opts) {
            const allDiv = document.createElement('div');
            allDiv.className = 'mb-dirs-option' + (!_browserSelectedDir ? ' active' : '');
            allDiv.textContent = 'All directories';
            allDiv.addEventListener('click', () => {
                _browserSelectedDir = null;
                _setDirsTriggerText('All directories');
                _closeDirsDropdown();
                filterModelBrowser();
            });
            opts.appendChild(allDiv);
        }
        (dirs || []).forEach(d => _addDirOption(d, d === _browserSelectedDir));
    } catch (err) {
        console.warn('[Browser] Failed to load directories:', err);
    }
}

async function addModelDirectory() {
    const app = window.go?.main?.App;
    if (!app?.ScanLLMDir) return;
    try {
        const dir = typeof openDirectoryLocator === 'function'
            ? await openDirectoryLocator('')
            : await app.OpenDirectoryPicker();
        if (!dir) {
            // The native picker returns an empty string when no folder is
            // actually selected — for example while a search filter hides
            // every subfolder. Point at the path box instead of failing
            // silently, which looked like the button doing nothing.
            const dd = _getDirsDropdown();
            if (dd) dd.classList.add('open');
            document.getElementById('mb-add-dir-path')?.focus();
            if (typeof showNotification === 'function') {
                showNotification('No folder selected. Paste the path instead.', 'error');
            }
            return;
        }
        await addModelDirectoryPath(dir);
    } catch (err) {
        console.warn('[Browser] Failed to add directory:', err);
        if (typeof showNotification === 'function') {
            showNotification('Could not add directory: ' + (err?.message || err), 'error');
        }
    }
}

// addModelDirectoryPath adds one directory and refreshes everything that shows
// it. Shared by the folder picker and the path box so both behave identically.
async function addModelDirectoryPath(dir) {
    const app = window.go?.main?.App;
    if (!app?.ScanLLMDir) return;
    dir = (dir || '').trim();
    if (!dir) {
        if (typeof showNotification === 'function') showNotification('Enter a folder path first', 'error');
        return;
    }
    await app.ScanLLMDir(dir);
    // Reload the model list as well: this used to refresh only the directory
    // dropdown and the chat dropdown, so the browser stayed on "No models
    // found" until it was reopened.
    _browserSelectedDir = null;
    _setDirsTriggerText('All directories');
    await populateModelBrowser();
    await populateLLMDropdown();
    if (typeof showNotification === 'function') showNotification('Added ' + dir, 'success');
}

async function addModelDirectoryFromInput() {
    const input = document.getElementById('mb-add-dir-path');
    await addModelDirectoryPath(input?.value || '');
    if (input) input.value = '';
}

async function rescanModelDirectories() {
    const app = window.go?.main?.App;
    if (!app?.RescanAllDirs) return;
    try {
        app.RescanAllDirs();
        await populateModelBrowser();
        if (typeof showNotification === 'function') showNotification('Model directories rescanned', 'success');
    } catch (err) {
        console.warn('[Browser] Failed to rescan directories:', err);
        if (typeof showNotification === 'function') showNotification('Rescan failed', 'error');
    }
}

async function removeLLMModel(path, name) {
    if (!confirm(`Remove "${name}" from the model registry?\n\nThe file will NOT be deleted from disk.`)) return;
    const app = window.go?.main?.App;
    if (!app?.RemoveLLMModel) return;
    try {
        await app.RemoveLLMModel(path);
        _browserModels = _browserModels.filter(m => m.path !== path);
        filterModelBrowser();
        if (typeof showNotification === 'function') showNotification(`Removed "${name}"`, 'success');
    } catch (err) {
        console.warn('[Browser] Failed to remove model:', err);
        if (typeof showNotification === 'function') showNotification('Failed to remove model', 'error');
    }
}

async function removeModelDirectory(dir) {
    if (!confirm(`Remove directory "${dir}" and all its models from the registry?\n\nNo files will be deleted from disk.`)) return;
    const app = window.go?.main?.App;
    if (!app?.UnregisterLLMDir) return;
    try {
        await app.UnregisterLLMDir(dir);
        if (_browserSelectedDir === dir) {
            _browserSelectedDir = null;
            _setDirsTriggerText('All directories');
        }
        await populateModelBrowserDirs();
        _browserModels = _browserModels.filter(m => !m.path.toLowerCase().startsWith(dir.toLowerCase()));
        filterModelBrowser();
        if (typeof showNotification === 'function') showNotification('Directory removed', 'success');
    } catch (err) {
        console.warn('[Browser] Failed to remove directory:', err);
        if (typeof showNotification === 'function') showNotification('Failed to remove directory', 'error');
    }
}

async function populateModelBrowser() {
    const list = document.getElementById('model-browser-list');
    const detail = document.getElementById('model-browser-detail');
    if (!list) return;

    list.innerHTML = '<div class="mb-no-selection">Loading models...</div>';
    if (detail) detail.style.display = 'none';

    try {
        const app = window.go?.main?.App;
        if (!app?.GetLLMModels) {
            list.innerHTML = '<div class="mb-no-selection">Backend not available</div>';
            return;
        }
        const models = await app.GetLLMModels();
        _browserModels = _chatModels(models);
        renderModelBrowserList(_browserModels);
        await populateModelBrowserDirs();
    } catch (err) {
        console.error('[Browser] Failed to load models:', err);
        list.innerHTML = '<div class="mb-no-selection">Error loading models</div>';
    }
}

function renderModelBrowserList(models) {
    const list = document.getElementById('model-browser-list');
    if (!list) return;

    if (models.length === 0) {
        list.innerHTML = '<div class="mb-no-selection">No models found. Scan a directory first.</div>';
        return;
    }

    list.innerHTML = '';
    models.forEach(m => {
        const item = document.createElement('div');
        item.className = 'model-browser-item';
        item.dataset.path = m.path;
        item.onclick = () => selectModelBrowserItem(m.path);

        const sizeGB = m.size_mb > 0 ? (m.size_mb / 1024).toFixed(1) + ' GB' : '';
        const folder = String(m.path || '').replace(/[\\/][^\\/]+$/, '');
        item.title = m.path || m.name;
        item.innerHTML = `
            <div class="mb-item-info">
                <span class="mb-item-name">${escapeHtml(m.name)}</span>
                <span class="mb-item-meta">${escapeHtml(m.format)} ${sizeGB}</span>
                <span class="mb-item-path" title="${escapeHtml(m.path)}">${escapeHtml(folder)}</span>
            </div>
            <span class="mb-item-remove" title="Remove model from registry">&times;</span>
        `;
        item.querySelector('.mb-item-remove').addEventListener('click', (e) => {
            e.stopPropagation();
            removeLLMModel(m.path, m.name);
        });
        list.appendChild(item);
    });
}

function filterModelBrowser() {
    const query = document.getElementById('model-browser-search')?.value.toLowerCase() || '';
    let filtered = _browserModels;
    if (_browserSelectedDir) {
        filtered = filtered.filter(m => m.path && m.path.toLowerCase().startsWith(_browserSelectedDir.toLowerCase()));
    }
    if (query) {
        filtered = filtered.filter(m => m.name.toLowerCase().includes(query));
    }
    renderModelBrowserList(filtered);
}

async function selectModelBrowserItem(path) {
    _browserSelectedPath = path;

    document.querySelectorAll('.model-browser-item').forEach(el => {
        el.classList.toggle('active', el.dataset.path === path);
    });

    const detail = document.getElementById('model-browser-detail');
    if (!detail) return;
    detail.style.display = 'block';

    // Show loading state
    document.getElementById('mb-detail-name').textContent = 'Loading...';

    try {
        const app = window.go?.main?.App;
        if (!app?.GetLLMModelDetails) {
            document.getElementById('mb-detail-name').textContent = 'Details unavailable';
            return;
        }
        const d = await app.GetLLMModelDetails(path);
        document.getElementById('mb-detail-name').textContent = d.name || 'Unknown';
        document.getElementById('mb-detail-family').textContent = d.family || '—';
        document.getElementById('mb-detail-params').textContent = d.parameter_count || '—';
        document.getElementById('mb-detail-quant').textContent = d.quant || '—';
        document.getElementById('mb-detail-layers').textContent = d.layer_count || '—';
        document.getElementById('mb-detail-context').textContent = d.max_context ? d.max_context.toLocaleString() : '—';
        document.getElementById('mb-detail-size').textContent = d.size_mb ? d.size_mb.toFixed(1) + ' MB' : '—';
        document.getElementById('mb-detail-vram').textContent = d.gpu_mb ? d.gpu_mb + ' MB' : '—';
        document.getElementById('mb-detail-ram').textContent = d.ram_mb ? d.ram_mb + ' MB' : '—';
		await populateModelIntelligence(path);
    } catch (err) {
        console.error('[Browser] Failed to get model details:', err);
        document.getElementById('mb-detail-name').textContent = 'Error loading details';
    }
}

async function populateModelIntelligence(path) {
	const app = window.go?.main?.App;
	const set = (id, value) => { const el = document.getElementById(id); if (el) el.textContent = value || '—'; };
	if (typeof app?.GetModelIntelligence !== 'function') {
		set('mb-detail-evidence', 'Unavailable');
		set('mb-detail-capabilities', '—');
		set('mb-detail-runtimes', '—');
		set('mb-detail-tokenizer', '—');
		return;
	}
	try {
		const profile = await app.GetModelIntelligence(path);
		set('mb-detail-evidence', evidenceLabel(profile.evidence));
		set('mb-detail-capabilities', (profile.capabilities || []).join(', '));
		set('mb-detail-runtimes', (profile.runtimes || []).join(', '));
		set('mb-detail-tokenizer', profile.tokenizer || 'Not found');
		const source = document.getElementById('mb-detail-source');
		if (source) source.textContent = profile.source ? 'Identification source: ' + profile.source : '';
		const notes = document.getElementById('mb-detail-notes');
		if (notes) notes.textContent = (profile.notes || []).join(' · ');
	} catch (err) {
		console.warn('[Browser] Model intelligence unavailable:', err);
		set('mb-detail-evidence', 'Unavailable');
	}
}

function evidenceLabel(value) {
	return ({ verified: 'Verified', 'source-derived': 'Source-derived', inferred: 'Inferred', unknown: 'Unknown', conflicting: 'Conflicting' })[value] || value || 'Unknown';
}

async function loadModelFromBrowser() {
    if (!_browserSelectedPath) return;
    const path = _browserSelectedPath; // capture before closeModelBrowser nulls it
    closeModelBrowser();
    console.log('[Browser] Loading model from path:', path);
    await _loadModelByPath(path);
    if (_modelSelectHasOption(path)) {
        const opt = [..._getModelSelectOptionsContainer().children].find(c => c.dataset.value === path);
        _setModelSelectValue(path, opt?.textContent || path);
    }
}

