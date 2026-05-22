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

function _addModelSelectOption(path, text) {
    const opts = _getModelSelectOptionsContainer();
    if (!opts) return;
    const div = document.createElement('div');
    div.className = 'model-select-option';
    div.dataset.value = path;
    div.textContent = text;
    div.addEventListener('click', () => {
        _setModelSelectValue(path, text);
        _closeModelSelect();
    });
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

async function autoLoadDefaultModel() {
    const defaultModel = window.settingsState?.current?.general?.defaultModel;
    if (defaultModel !== 'last-used') return;
    const savedPath = localStorage.getItem('howl-last-model-path');
    if (!savedPath) return;
    if (!window.go?.main?.App?.LoadModel) return;
    if (isModelLoading()) {
        console.log('Auto-load skipped: model load already in progress');
        return;
    }

    updateModelDisplay('Loading...');
    setModelLoading(true);
    try {
        await window.go.main.App.LoadModel(savedPath);
        const modelName = savedPath.split(/[\\/]/).pop();
        updateModelDisplay(modelName);
        console.log('Auto-loaded model:', modelName);
    } catch (err) {
        console.error('Auto-load model failed:', err);
        updateModelDisplay('No model loaded');
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
                updateModelDisplay('No model loaded');
                autoLoadDefaultModel();
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

async function populateLLMDropdown() {
    const app = window.go?.main?.App;
    _setupModelSelect();
    if (!app?.GetLLMModels) return;
    try {
        const models = await app.GetLLMModels();
        const prev = _getModelSelectValue();
        _clearModelSelectOptions();
        (models || []).forEach(m => {
            const sizeTxt = m.size_mb > 0 ? ` (${(m.size_mb / 1024).toFixed(1)} GB)` : '';
            _addModelSelectOption(m.path, `${m.name}${sizeTxt}`);
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
        const dir = await app.OpenDirectoryPicker();
        if (!dir) return;
        updateModelDisplay('Scanning…');
        await app.ScanLLMDir(dir);
        await populateLLMDropdown();
        updateModelDisplay('Scan complete');
    } catch (err) {
        console.error('[LLM] Directory scan failed:', err);
        updateModelDisplay('Scan failed');
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

function openModelBrowser() {
    const modal = document.getElementById('model-browser-modal');
    if (modal) modal.style.display = 'flex';
    populateModelBrowser();
}

function closeModelBrowser() {
    const modal = document.getElementById('model-browser-modal');
    if (modal) modal.style.display = 'none';
    _browserSelectedPath = null;
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
        _browserModels = models || [];
        renderModelBrowserList(_browserModels);
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
        item.innerHTML = `
            <span class="mb-item-name">${escapeHtml(m.name)}</span>
            <span class="mb-item-meta">${escapeHtml(m.format)} ${sizeGB}</span>
        `;
        list.appendChild(item);
    });
}

function filterModelBrowser() {
    const query = document.getElementById('model-browser-search')?.value.toLowerCase() || '';
    const filtered = _browserModels.filter(m => m.name.toLowerCase().includes(query));
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
    } catch (err) {
        console.error('[Browser] Failed to get model details:', err);
        document.getElementById('mb-detail-name').textContent = 'Error loading details';
    }
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

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}
