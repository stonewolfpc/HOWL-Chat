'use strict';

const IG_STATE_KEY = 'howl-image-gen-state';
let currentResult = null;
let models = [];       // full ModelInfo list (diffusion only)
let scanIndex = null;  // full ScanResult from ScanImageModels
let seedHistory = [];
let lastUsedSeed = -1;
let currentModelFamily = '';
let generationInFlight = false;

// VRAM Check for Generation
async function checkVRAMForGeneration(opts) {
    try {
        // Get GPU info from backend
        const gpuInfo = window.go?.main?.App?.GetGPUInfo?.();
        if (!gpuInfo || !gpuInfo.TotalMemory) {
            return { willExceedVRAM: false, estimatedVRAM: 0, availableVRAM: 0 };
        }
        
        const availableVRAM = gpuInfo.TotalMemory; // in MB
        const freeVRAM = gpuInfo.FreeMemory || (availableVRAM * 0.7); // estimate if not provided
        
        // Estimate VRAM needed based on resolution and model family
        // Base estimate: ~2MB per pixel for high-quality generation
        const pixels = opts.width * opts.height;
        const baseVRAM = (pixels * 2) / (1024 * 1024); // Convert to MB
        
        // Model family multiplier
        const family = (window.currentModelFamily || '').toLowerCase();
        const modelMultipliers = {
            'flux': 1.5,
            'sdxl': 1.3,
            'sd3': 1.4,
            'chroma': 1.5,
            'sd2.1': 1.0,
            'sd21': 1.0,
            'default': 1.0
        };
        const multiplier = modelMultipliers[family] || modelMultipliers['default'];
        
        // Batch size multiplier
        const batchMultiplier = opts.batch_count || 1;
        
        // Steps multiplier (more steps = more VRAM for intermediate activations)
        const stepsMultiplier = 1 + (opts.steps / 100) * 0.2;
        
        // CFG scale multiplier
        const cfgMultiplier = 1 + (opts.cfg_scale / 10) * 0.1;
        
        const estimatedVRAM = Math.floor(baseVRAM * multiplier * batchMultiplier * stepsMultiplier * cfgMultiplier);
        
        // Will exceed if estimated > free VRAM with some buffer
        const willExceedVRAM = estimatedVRAM > (freeVRAM * 0.9); // 10% buffer
        
        console.log(`[VRAM Check] Estimated: ${estimatedVRAM} MB, Available: ${freeVRAM} MB, Exceed: ${willExceedVRAM}`);
        
        return {
            willExceedVRAM,
            estimatedVRAM,
            availableVRAM: freeVRAM
        };
    } catch (e) {
        console.warn('VRAM check failed:', e);
        return { willExceedVRAM: false, estimatedVRAM: 0, availableVRAM: 0 };
    }
}
window.currentModelFamily = currentModelFamily; // keep in sync below via setter

// ---------------------------------------------------------------------------
// Diagnostics
// ---------------------------------------------------------------------------
function showRuntimeDebug() {
    const dbg = document.getElementById('ig-runtime-debug');
    if (!dbg) return;
    const hasGo = !!window.go;
    const hasMain = hasGo && !!window.go.main;
    const hasApp = hasMain && !!window.go.main.App;
    const hasDirPicker = hasApp && !!window.go.main.App.OpenImageModelsDirPicker;
    const hasFilePicker = hasApp && !!window.go.main.App.OpenSDServerFilePicker;
    dbg.textContent = 'Runtime: ' + (hasGo ? 'go ' : 'no-go ') +
        (hasMain ? 'main ' : '') + (hasApp ? 'App ' : '') +
        (hasDirPicker ? 'DirPicker' : 'NoDirPicker') + ' ' +
        (hasFilePicker ? 'FilePicker' : 'NoFilePicker');
}

// ---------------------------------------------------------------------------
// Path Configuration
// ---------------------------------------------------------------------------

let _pathApplyTimer = null;

// Called once on page load. Reads persisted config from the Go backend,
// populates the inputs, and immediately triggers a model scan.
// No user action required — it just works.
async function initPaths() {
    if (!window.go?.main?.App) { loadModels(); return; }
    try {
        const saved = await window.go.main.App.GetImageSettings();
        const modelsDir = saved?.engine_config?.models_dir || '';
        const sdPath    = saved?.engine_config?.binary_path || '';
        if (modelsDir) document.getElementById('ig-models-dir').value = modelsDir;
        if (sdPath)    document.getElementById('ig-sd-server').value  = sdPath;
        if (modelsDir || sdPath) {
            setStatus('Paths loaded from saved config');
        }
    } catch (err) {
        console.warn('[IG] Could not load saved image settings:', err);
    }
    loadModels();
}

// Persist both paths to the Go backend (which writes image-settings.json).
// Also mirrors into settingsState for the frontend settings panel.
async function _persistPaths(modelsDir, sdPath) {
    const app = window.go?.main?.App;
    if (!app) return;
    try {
        if (modelsDir) await app.SetImageModelsDir(modelsDir);
        if (sdPath)    await app.SetSDServerPath(sdPath);
        if (window.settingsState?.current) {
            window.settingsState.current.imageGen = window.settingsState.current.imageGen || {};
            window.settingsState.current.imageGen.modelsDir    = modelsDir;
            window.settingsState.current.imageGen.sdServerPath = sdPath;
            window.settingsState.saveSettings?.();
        }
    } catch (err) {
        setStatus('Path error: ' + (err.message || err));
        console.error('[IG] _persistPaths failed:', err);
    }
}

// Auto-apply when the user finishes typing (500 ms debounce), then re-scan.
function _schedulePathApply() {
    clearTimeout(_pathApplyTimer);
    _pathApplyTimer = setTimeout(async () => {
        const modelsDir = document.getElementById('ig-models-dir')?.value.trim() || '';
        const sdPath    = document.getElementById('ig-sd-server')?.value.trim()  || '';
        await _persistPaths(modelsDir, sdPath);
        await loadModels();
    }, 500);
}

// Wire input change listeners — called once after the DOM is ready.
function initPathListeners() {
    document.getElementById('ig-models-dir')?.addEventListener('change', _schedulePathApply);
    document.getElementById('ig-sd-server')?.addEventListener('change', _schedulePathApply);
}

async function browseModelsDir() {
    const input = document.getElementById('ig-models-dir');
    if (!window.go?.main?.App?.OpenImageModelsDirPicker) { input.focus(); return; }
    try {
        const dir = await window.go.main.App.OpenImageModelsDirPicker();
        if (dir) {
            input.value = dir;
            await _persistPaths(dir, document.getElementById('ig-sd-server')?.value.trim() || '');
            await loadModels();
        }
    } catch (err) { console.error('Browse models dir failed:', err); }
}

async function browseSDServer() {
    const input = document.getElementById('ig-sd-server');
    if (!window.go?.main?.App?.OpenSDServerFilePicker) { input.focus(); return; }
    try {
        const path = await window.go.main.App.OpenSDServerFilePicker();
        if (path) {
            input.value = path;
            await _persistPaths(document.getElementById('ig-models-dir')?.value.trim() || '', path);
            await loadModels();
        }
    } catch (err) { console.error('Browse sd-server failed:', err); }
}

// applyPaths kept for the explicit "Apply" button — now just delegates.
async function applyPaths() {
    const modelsDir = document.getElementById('ig-models-dir')?.value.trim() || '';
    const sdPath    = document.getElementById('ig-sd-server')?.value.trim()  || '';
    if (!modelsDir && !sdPath) { setStatus('Set at least one path.'); return; }
    setStatus('Applying paths...');
    await _persistPaths(modelsDir, sdPath);
    await loadModels();
    setStatus('Ready');
}

// ---------------------------------------------------------------------------
// Model Loading  (full-index, family-aware)
// ---------------------------------------------------------------------------

// loadModels performs a full scan via ScanImageModels, then populates all dropdowns.
async function loadModels() {
    const app = window.go?.main?.App;
    if (!app) { console.warn('[IG] Wails runtime not available'); return; }
    setStatus('Scanning models...');
    try {
        // Try the fast cached registry first (no disk scan), then fall back
        // to a full scan only if the fast path is unavailable or empty.
        if (app.GetRegisteredModels) {
            scanIndex = await app.GetRegisteredModels();
        }
        if (!scanIndex || (!scanIndex.models || scanIndex.models.length === 0) && app.ScanImageModels) {
            scanIndex = await app.ScanImageModels();
        }
        if (!scanIndex || (!scanIndex.models || scanIndex.models.length === 0) && app.GetModelIndex) {
            scanIndex = await app.GetModelIndex();
        }
        if (!scanIndex) {
            // Legacy fallback
            const legacy = await app.ListImageModels();
            scanIndex = { models: legacy || [], vaes: [], loras: [], controlnets: [], embeddings: [], upscalers: [], t5xxls: [], clip_ls: [], clip_gs: [] };
        }
        if (!scanIndex) scanIndex = { models: [], vaes: [], loras: [], controlnets: [], embeddings: [], upscalers: [], t5xxls: [], clip_ls: [], clip_gs: [] };

        models = scanIndex.models || [];
        console.log('[IG] Scan:', models.length, 'models,', (scanIndex.vaes||[]).length, 'VAEs,',
                    (scanIndex.loras||[]).length, 'LoRAs,', (scanIndex.controlnets||[]).length, 'ControlNets',
                    (scanIndex.t5xxls||[]).length, 'T5,', (scanIndex.clip_ls||[]).length, 'CLIP-L,', (scanIndex.clip_gs||[]).length, 'CLIP-G');

        _populateDiffusionDropdown();
        _populateVAEDropdown(scanIndex.vaes || []);
        _populateLoRAList(scanIndex.loras || []);
        _populateControlNetDropdown(scanIndex.controlnets || []);
        _populateT5XXLDropdown(scanIndex.t5xxls || []);
        _populateClipLDropdown(scanIndex.clip_ls || []);
        _populateClipGDropdown(scanIndex.clip_gs || []);

        if (models.length === 0) {
            setStatus('No models found. Check Engine Paths.');
            return;
        }
        const total = models.length + (scanIndex.vaes||[]).length + (scanIndex.loras||[]).length;
        setStatus(`${models.length} model(s) · ${(scanIndex.vaes||[]).length} VAE · ${(scanIndex.loras||[]).length} LoRA discovered`);

        // Restore previously selected model
        const saved = JSON.parse(localStorage.getItem(IG_STATE_KEY) || '{}');
        if (saved.model && models.some(m => m.name === saved.model)) {
            document.getElementById('ig-model').value = saved.model;
            onModelChange();
        }
    } catch (err) {
        console.error('[IG] Failed to load models:', err);
        setStatus('Error loading models: ' + (err.message || err));
    }
}

// Refresh family-filtered secondary dropdowns without a full rescan.
async function refreshCompatibleAssets(family) {
    const app = window.go?.main?.App;
    if (!app?.GetCompatibleAssets || !family) return;
    try {
        const filtered = await app.GetCompatibleAssets(family);
        _populateVAEDropdown(filtered.vaes || []);
        _populateLoRAList(filtered.loras || []);
        _populateControlNetDropdown(filtered.controlnets || []);
    } catch (err) {
        console.warn('[IG] GetCompatibleAssets failed, using full index:', err);
    }
}

// Exposed for image-model-library.js — called after a successful install to
// refresh all dropdowns without requiring a separate ScanImageModels call.
window.refreshModelDropdownsFromScan = loadModels;

// ---- individual dropdown populators ----------------------------------------

function _populateDiffusionDropdown() {
    const sel = document.getElementById('ig-model');
    if (!sel) return;
    sel.innerHTML = '';
    // Only files stable-diffusion.cpp can actually load belong in this dropdown.
    // LoRAs, VAEs, ControlNets, embeddings and workflow JSON are classified
    // separately (Model Library) and must never appear as if they were a model.
    const split = (typeof window.splitCompatibleImageModels === 'function')
        ? window.splitCompatibleImageModels(models, (typeof currentMode !== 'undefined' ? currentMode : 'txt2img'))
        : { usable: models, hidden: 0, total: models.length };
    const usable = split.usable;
    const ph = document.createElement('option');
    ph.value = '';
    ph.disabled = true;
    ph.selected = true;
    ph.textContent = usable.length === 0 ? 'No compatible models found' : 'Select a model...';
    sel.appendChild(ph);
    if (split.hidden > 0) {
        // Use an optgroup label so the count is visible but never selectable.
        // A disabled <option> can still be selected in some browsers via click
        // or script, which sends the literal label string to the backend and
        // produces a confusing "model not found" error.
        const og = document.createElement('optgroup');
        og.label = '(' + split.hidden + ' incompatible file' + (split.hidden === 1 ? '' : 's') + ' hidden)';
        sel.appendChild(og);
    }
    usable.forEach(m => {
        const opt = document.createElement('option');
        opt.value = m.name;
        const tag = m.format === 'bundle' ? '[bundle]' : m.format ? `[${m.format}]` : '';
        const warn = m.compat === 'unsupported' ? '[⚠ outdated binary]' : '';
        opt.textContent = `${m.name} (${m.family || '?'}) ${tag} ${warn}`.trim();
        opt.dataset.family = m.family || '';
        opt.dataset.compat = m.compat || '';
        opt.dataset.recommended = JSON.stringify(m.recommended || {});
        if (m.compat === 'unsupported') {
            opt.style.color = '#ff6b6b';
            opt.title = 'This model requires a newer sd-server binary';
        }
        sel.appendChild(opt);
    });
}

function _populateVAEDropdown(vaes) {
    const sel = document.getElementById('ig-vae-model');
    if (!sel) return;
    sel.innerHTML = '<option value="">None (built-in)</option>';
    vaes.forEach(v => {
        const opt = document.createElement('option');
        opt.value = v.path;
        opt.textContent = `${v.name} [${v.family || 'any'}]`;
        sel.appendChild(opt);
    });
}

function _populateLoRAList(loras) {
    const sel = document.getElementById('ig-lora-selector');
    if (!sel) return;
    sel.innerHTML = '<option value="">Select a LoRA to add...</option>';
    loras.forEach(l => {
        const opt = document.createElement('option');
        opt.value = l.name;
        opt.textContent = `${l.name} [${l.family || 'any'}]`;
        sel.appendChild(opt);
    });
}

function _populateControlNetDropdown(cns) {
    const sel = document.getElementById('ig-control-model');
    if (!sel) return;
    const prev = sel.value;
    sel.innerHTML = '<option value="">None</option>';
    cns.forEach(c => {
        const opt = document.createElement('option');
        opt.value = c.path;
        opt.textContent = `${c.name} [${c.family || 'any'}]`;
        sel.appendChild(opt);
    });
    if (prev) sel.value = prev; // preserve selection if still present

    // A ControlNet can also be bound to the checkpoint bundle, where sd-server
    // loads it at startup. Say so rather than showing an empty list.
    const modelSel = document.getElementById('ig-model');
    const boundModel = (models || []).find(m => m.name === (modelSel && modelSel.value)) || {};
    const bound = (boundModel.files || {}).control_net || '';
    const hint = document.getElementById('ig-controlnet-status');
    if (bound && hint) {
        const fname = String(bound).split(/[\/]/).pop();
        hint.textContent = 'Bound to this model: ' + fname + ' — pose applies automatically.';
    } else if (hint && !prev) {
        hint.textContent = '';
    }
}

// onModelChange updates UI presets and re-filters secondary dropdowns by family.
function onModelChange() {
    const name = document.getElementById('ig-model').value;
    if (!name) return;
    const model = models.find(m => m.name === name);
    if (model) {
        autoPopulateModelSpecificFields(model);
        currentModelFamily = model.family || '';
        window.currentModelFamily = currentModelFamily;
        applyModelPreset(model.family, model.recommended || null);
        updateLiveModelInfo(model);
        // Re-filter VAE / LoRA / ControlNet for this family
        if (currentModelFamily) refreshCompatibleAssets(currentModelFamily);
    }
    const statusEl = document.getElementById('ig-model-status');
    if (statusEl) {
        if (model && model.compat === 'unsupported') {
            statusEl.textContent = '⚠ Requires newer sd-server';
            statusEl.className = 'ig-panel-status ig-status-error';
        } else {
            statusEl.textContent = 'Auto-load on Generate';
            statusEl.className = 'ig-panel-status';
        }
    }
    saveState();
}

function _pickFirstPath(list) {
    if (!Array.isArray(list) || list.length === 0) return '';
    const item = list[0] || {};
    return item.path || '';
}

function autoPopulateModelSpecificFields(model) {
    if (!model) return;

    const files = model.files || {};
    const family = String(model.family || '').toLowerCase();

    // Intentionally excludes VAE/LoRA/ControlNet/workflow fields per product behavior.
    const clipL = files.clip_l || _pickFirstPath(scanIndex?.clip_ls);
    const clipG = files.clip_g || _pickFirstPath(scanIndex?.clip_gs);
    const t5xxl = files.t5xxl || _pickFirstPath(scanIndex?.t5xxls);

    const clipLSel = document.getElementById('ig-clip-l');
    const clipGSel = document.getElementById('ig-clip-g');
    const t5Sel = document.getElementById('ig-t5xxl');

    const usesT5 = ['flux', 'flux2', 'sd3', 'auraflow', 'hidream'].includes(family);
    const usesClipL = ['flux', 'flux2', 'sd3', 'auraflow', 'hidream', 'sdxl', 'sdxl-turbo'].includes(family);
    const usesClipG = ['sd3', 'sdxl', 'sdxl-turbo'].includes(family);

    if (clipLSel) {
        if (usesClipL && clipL) clipLSel.value = clipL;
        if (!usesClipL) clipLSel.value = '';
    }
    if (clipGSel) {
        if (usesClipG && clipG) clipGSel.value = clipG;
        if (!usesClipG) clipGSel.value = '';
    }
    if (t5Sel) {
        if (usesT5 && t5xxl) t5Sel.value = t5xxl;
        if (!usesT5) t5Sel.value = '';
    }
}

// setModel manually loads a model into the backend (used by Reload button).
async function setModel(name) {
    const model = models.find(m => m.name === name);
    if (model && model.compat === 'unsupported') {
        setStatus('⚠ This model requires a newer sd-server binary. Please update and try again.');
        const statusEl = document.getElementById('ig-model-status');
        if (statusEl) { statusEl.textContent = '⚠ Unsupported'; statusEl.className = 'ig-panel-status ig-status-error'; }
        updateDiagnostics('error', 'Model requires newer sd-server');
        return;
    }
    // A split-component model cannot load without its text encoder and VAE.
    // Ask for consent before going online, and stop quietly if the user says no.
    if (window.ImageDeps && !(await window.ImageDeps.ensureReady(name))) {
        const depsStatus = document.getElementById('ig-model-status');
        if (depsStatus) { depsStatus.textContent = 'Not loaded — dependencies missing'; depsStatus.className = 'ig-panel-status ig-status-error'; }
        setStatus('Load cancelled — the model is missing required files.');
        return;
    }
    const statusEl = document.getElementById('ig-model-status');
    if (statusEl) { statusEl.textContent = 'Loading...'; statusEl.className = 'ig-panel-status ig-status-loading'; }
    updateDiagnostics('loadStart');
    setStatus('Loading model: ' + name + '...');
    try {
        autoPopulateModelSpecificFields(model);

        const opts = {
            model: name,
            vae: document.getElementById('ig-vae-model')?.value || '',
            clip_l: document.getElementById('ig-clip-l')?.value || '',
            clip_g: document.getElementById('ig-clip-g')?.value || '',
            t5xxl: document.getElementById('ig-t5xxl')?.value || '',
            diffusion_model: (model?.files && model.files.diffusion_model) ? model.files.diffusion_model : '',
            backend: document.getElementById('ig-backend')?.value || '',
            params_backend: document.getElementById('ig-params-backend')?.value || '',
            vae_tiling: document.getElementById('ig-vae-tiling')?.checked || false,
            temporal_tiling: document.getElementById('ig-temporal-tiling')?.checked || false,
            circular: document.getElementById('ig-circular')?.checked || false,
            flash_attention: document.getElementById('ig-fa')?.checked || false,
            offload_to_cpu: document.getElementById('ig-offload-cpu')?.checked || false,
            disable_metadata: document.getElementById('ig-disable-metadata')?.checked || false,
        };
        console.log('[IG] Loading model with config:', opts);
        await window.go.main.App.SetActiveImageModel(opts);
        if (model) {
            currentModelFamily = model.family || '';
            window.currentModelFamily = currentModelFamily;
            applyModelPreset(model.family);
            updateLiveModelInfo(model);
        }
        setStatus('Model ready: ' + name);
        if (statusEl) { statusEl.textContent = 'Loaded: ' + name.replace(/^.*[\\/]/, ''); statusEl.className = 'ig-panel-status ig-status-loaded'; }
        updateDiagnostics('loadDone');
    } catch (err) {
        const msg = err.message || String(err);
        let displayMsg = 'Error loading model: ' + msg;
        if (msg.includes('Chroma-family GGUF')) {
            displayMsg = 'Chroma GGUF may not be supported by older stable-diffusion.cpp builds.\nTry the .safetensors version or update your sd-server binary.';
        } else if (msg.includes('get sd version from file failed') || msg.includes('cannot be loaded')) {
            displayMsg = 'Model failed to load — your sd-server binary is outdated.\nPlease download the latest stable-diffusion.cpp build and update your Engine Paths.';
        }
        setStatus(displayMsg);
        if (statusEl) { statusEl.textContent = 'Error'; statusEl.className = 'ig-panel-status ig-status-error'; }
        updateDiagnostics('error', msg);
    }
}

function _populateT5XXLDropdown(items) {
    const sel = document.getElementById('ig-t5xxl');
    if (!sel) return;
    const val = sel.value;
    sel.innerHTML = '<option value="">None</option>';
    items.forEach(m => {
        const opt = document.createElement('option');
        opt.value = m.path;
        opt.textContent = m.name || m.path.replace(/^.*[\\/]/, '');
        sel.appendChild(opt);
    });
    if (val) sel.value = val;
}

function _populateClipLDropdown(items) {
    const sel = document.getElementById('ig-clip-l');
    if (!sel) return;
    const val = sel.value;
    sel.innerHTML = '<option value="">None</option>';
    items.forEach(m => {
        const opt = document.createElement('option');
        opt.value = m.path;
        opt.textContent = m.name || m.path.replace(/^.*[\\/]/, '');
        sel.appendChild(opt);
    });
    if (val) sel.value = val;
}

function _populateClipGDropdown(items) {
    const sel = document.getElementById('ig-clip-g');
    if (!sel) return;
    const val = sel.value;
    sel.innerHTML = '<option value="">None</option>';
    items.forEach(m => {
        const opt = document.createElement('option');
        opt.value = m.path;
        opt.textContent = m.name || m.path.replace(/^.*[\\/]/, '');
        sel.appendChild(opt);
    });
    if (val) sel.value = val;
}

function applyModelPreset(family, recommended) {
    // Prefer scan-manifest recommended values; fall back to family defaults.
    const rec = recommended || _familyDefaults(family);
    const set = (id, val) => { const el = document.getElementById(id); if (el && val !== undefined && val !== null) el.value = val; };
    const setChecked = (id, val) => { const el = document.getElementById(id); if (el && val !== undefined) el.checked = val; };
    set('ig-steps',     rec.steps);
    set('ig-cfg',       rec.cfg);
    set('ig-guidance',  rec.guidance);
    set('ig-width',     rec.width);
    set('ig-height',    rec.height);
    set('ig-clip-skip', rec.clipSkip);
    if (rec.sampler)   set('ig-sampler',   rec.sampler);
    if (rec.scheduler) set('ig-scheduler', rec.scheduler);
    if (rec.flowShift !== undefined) set('ig-flow-shift', rec.flowShift);
    updateLabels();
    saveState();
}

function _familyDefaults(family) {
    const f = (family || '').toLowerCase();
    switch (f) {
        case 'flux':
            return { steps: 20, cfg: 1,   guidance: 3.5, width: 1024, height: 1024, sampler: 'euler', scheduler: 'simple', clipSkip: -1, flowShift: 0 };
        case 'hidream':
            return { steps: 20, cfg: 4.5, guidance: 3.5, width: 1024, height: 1024, sampler: 'euler', scheduler: 'simple', clipSkip: -1, flowShift: 0 };
        case 'sd1.5': case 'sd15':
            return { steps: 20, cfg: 7,   guidance: 3.5, width: 512,  height: 512,  sampler: 'euler_a', scheduler: 'karras', clipSkip: 1 };
        case 'sd2.1': case 'sd21':
            return { steps: 20, cfg: 7.5, guidance: 3.5, width: 768,  height: 768,  sampler: 'dpmpp_2m', scheduler: 'karras', clipSkip: 2 };
        case 'sdxl-turbo': case 'turbo':
            return { steps: 4,  cfg: 1,   guidance: 3.5, width: 512,  height: 512,  sampler: 'euler_a', scheduler: 'karras', clipSkip: 2 };
        case 'sd3':
            return { steps: 20, cfg: 7,   guidance: 3.5, width: 1024, height: 1024, sampler: 'euler', scheduler: 'simple', clipSkip: -1, flowShift: 0 };
        case 'chroma':
            return { steps: 20, cfg: 7,   guidance: 3.5, width: 1024, height: 1024, sampler: 'euler', scheduler: 'simple', clipSkip: -1, flowShift: 0 };
        case 'kolors':
            return { steps: 20, cfg: 7,   guidance: 3.5, width: 1024, height: 1024, sampler: 'euler', scheduler: 'karras', clipSkip: -1 };
        case 'anima':
            return { steps: 20, cfg: 7,   guidance: 3.5, width: 1024, height: 1024, sampler: 'euler', scheduler: 'simple', clipSkip: -1 };
        // Video families keep their native landscape frame instead of a square:
        // these match sdcpp.videoResolution, so the UI and engine agree.
        case 'wan': case 'wan2.1': case 'wan2.2': case 'wan21': case 'wan22':
            // cfg 7 / guidance 3.5 are sd.cpp's own defaults and the only values
            // verified to produce a coherent clip; cfg 1 washed the frames out.
            return { steps: 20, cfg: 7,   guidance: 3.5, width: 832,  height: 480,  sampler: 'euler', scheduler: 'simple', clipSkip: -1, flowShift: 0 };
        case 'ltx': case 'ltx2': case 'ltx-video':
            return { steps: 20, cfg: 7,   guidance: 3.5, width: 768,  height: 512,  sampler: 'euler', scheduler: 'simple', clipSkip: -1, flowShift: 0 };
        default: // sdxl, unknown
            return { steps: 30, cfg: 7,   guidance: 3.5, width: 1024, height: 1024, sampler: 'dpmpp_2m', scheduler: 'karras', clipSkip: 2 };
    }
}

// ---------------------------------------------------------------------------
// Generation
// ---------------------------------------------------------------------------
// VRAM arbitration before generation: sd-server needs GPU room for the
// checkpoint, and a loaded chat model is holding most of the card.
async function requestImageVRAMArbitration() {
    const app = window.go?.main?.App;
    if (!app?.CheckModelFit) return true;
    try {
        let fit = await app.CheckModelFit('image', '');
        if (fit.can_fit || !fit.blocking || fit.blocking.length === 0) return true;
        if (!confirm(fit.message)) return false;
        setStatus('Unloading the chat model…');
        fit = await app.FreeVRAMFor('image', '');
        if (!fit.can_fit) { alert(fit.message); setStatus('Not enough VRAM'); return false; }
        return true;
    } catch (err) {
        console.warn('[VRAM] arbitration check failed, continuing:', err);
        return true;
    }
}

async function onGenerate() {
    if (generationInFlight) {
        setStatus('Generation already in progress...');
        return;
    }
    generationInFlight = true;
    const prompt = document.getElementById('ig-prompt').value.trim();
    if (!prompt) { setStatus('Please enter a prompt.'); generationInFlight = false; return; }
    const modelSel = document.getElementById('ig-model');
    if (!modelSel.value) { setStatus('Please select a model.'); generationInFlight = false; return; }
    const model = models.find(m => m.name === modelSel.value);
    if (model && model.compat === 'unsupported') {
        setStatus('⚠ This model requires a newer sd-server binary. Please update and try again.');
        updateDiagnostics('error', 'Model requires newer sd-server');
        generationInFlight = false;
        return;
    }

    // Generating auto-loads the model, so check its files before the run starts.
    if (window.ImageDeps && !(await window.ImageDeps.ensureReady(modelSel.value))) {
        setStatus('Generation cancelled — the model is missing required files.');
        generationInFlight = false;
        return;
    }

    // The chat model is holding the VRAM; ask before unloading it.
    if (!(await requestImageVRAMArbitration())) {
        setStatus('Cancelled — not enough VRAM with the chat model loaded.');
        generationInFlight = false;
        return;
    }

    const btn = document.getElementById('ig-generate');
    btn.disabled = true;
    setStatus('Generating...');
    showSpinner();
    updateDiagnostics('genStart');

    const steps = parseInt(document.getElementById('ig-steps').value, 10) || 20;
    startProgressBar(steps, currentModelFamily || 'sdxl');

    const seedVal = parseInt(document.getElementById('ig-seed').value, 10);
    if (seedLocked && !isNaN(seedVal) && seedVal !== -1) {
        // Seed is locked to a specific value
        lastUsedSeed = seedVal;
    } else if (seedLocked && (isNaN(seedVal) || seedVal === -1)) {
        // Seed is locked but was random: generate once and lock it
        lastUsedSeed = Math.floor(Math.random() * 2147483647);
        document.getElementById('ig-seed').value = lastUsedSeed;
    } else if (seedVal === -1 || isNaN(seedVal)) {
        // Not locked, use random seed
        lastUsedSeed = Math.floor(Math.random() * 2147483647);
        document.getElementById('ig-seed').value = lastUsedSeed;
    } else {
        lastUsedSeed = seedVal;
    }
    // Upscale mode has its own seed field; when locked it wins over the main seed.
    if (currentMode === 'upscale' && typeof upscaleSeedLocked !== 'undefined' && upscaleSeedLocked) {
        const upscaleSeedEl = document.getElementById('ig-upscale-seed');
        const lockedSeed = upscaleSeedEl ? parseInt(upscaleSeedEl.value, 10) : NaN;
        if (!isNaN(lockedSeed) && lockedSeed >= 0) {
            document.getElementById('ig-seed').value = lockedSeed;
            lastUsedSeed = lockedSeed;
        }
    }
    seedHistory.unshift(lastUsedSeed);
    if (seedHistory.length > 10) seedHistory.pop();

    // Handle img2img / inpaint init image upload
    let initImagePath = '';
    let maskPath = '';
    let controlNetPath = '';
    
    const needsInitImage = ['img2img', 'inpaint', 'outpaint', 'upscale'].includes(currentMode);
    if (needsInitImage) {
        if (!initImageData) {
            setStatus('Please upload an image first.');
            btn.disabled = false;
            hideSpinner();
            return;
        }
        try {
            setStatus('Uploading image...');
            initImagePath = await window.go.main.App.SaveTempImage(initImageData);
        } catch (err) {
            console.error('Upload failed:', err);
            setStatus('Upload failed: ' + (err.message || err));
            btn.disabled = false;
            hideSpinner();
            return;
        }
    }

    // Handle ControlNet image if present
    if (window.controlNetImageData) {
        try {
            controlNetPath = await window.go.main.App.SaveTempImage(window.controlNetImageData);
        } catch (err) {
            console.error('ControlNet image upload failed:', err);
        }
    }

    // A control image without a ControlNet model is ignored by the engine: the
    // run burns a full sampling pass and looks unresponsive. Stop and say why.
    const controlModel = (document.getElementById('ig-control-model') || {}).value || '';
    // A ControlNet may already be bound to the checkpoint bundle, in which case
    // sd-server loads it at startup (--control_net) and no dropdown pick is needed.
    const boundModel = (models || []).find(m => m.name === modelSel.value) || {};
    const bundleControlNet = (boundModel.files || {}).control_net || '';
    if (window.controlNetImageData && !controlModel && !bundleControlNet) {
        setStatus('ControlNet: choose a Control Net model, or clear the control image.');
        btn.disabled = false;
        generationInFlight = false;
        hideSpinner();
        stopProgressBar();
        return;
    }

    // Handle video end image if present
    let endImagePath = '';
    if (currentMode === 'video' && window.videoEndImageData) {
        try {
            endImagePath = await window.go.main.App.SaveTempImage(window.videoEndImageData);
        } catch (err) {
            console.error('Video end image upload failed:', err);
        }
    }

    // Handle img2img mask if present
    if (currentMode === 'img2img' && window.img2imgMaskData) {
        try {
            maskPath = await window.go.main.App.SaveTempImage(window.img2imgMaskData);
        } catch (err) {
            console.error('Mask upload failed:', err);
        }
    }

    const opts = {
        prompt:          prompt,
        negative_prompt: document.getElementById('ig-negative').value,
        width:           parseInt(document.getElementById('ig-width').value, 10),
        height:          parseInt(document.getElementById('ig-height').value, 10),
        steps:           parseInt(document.getElementById('ig-steps').value, 10),
        seed:            lastUsedSeed,
        sampler:         readSamplerValue(),
        scheduler:       document.getElementById('ig-scheduler').value,
        cfg_scale:       parseFloat(document.getElementById('ig-cfg').value),
        guidance:        parseFloat(document.getElementById('ig-guidance')?.value || '3.5'),
        slg_scale:       parseFloat(document.getElementById('ig-slg-scale')?.value || '0'),
        eta:             parseFloat(document.getElementById('ig-eta')?.value || '0'),
        batch_count:     parseInt(document.getElementById('ig-batch').value, 10) || 1,
        mode:            currentMode || 'txt2img',
        model_name:      modelSel.value,
        vae_path:        document.getElementById('ig-vae-model')?.value || '',
        clip_l_path:     document.getElementById('ig-clip-l')?.value || '',
        clip_g_path:     document.getElementById('ig-clip-g')?.value || '',
        t5xxl_path:      document.getElementById('ig-t5xxl')?.value || '',
        control_net_path: document.getElementById('ig-control-model')?.value || '',
        control_image_path: controlNetPath,
        control_strength: parseFloat(document.getElementById('ig-control-strength')?.value || '0.9'),
        loras:           window.activeLoRAs || [],
        enable_hr:       document.getElementById('ig-hires-enable')?.checked || false,
        hr_scale:        parseFloat(document.getElementById('ig-hr-scale')?.value || '2'),
        hr_upscaler:     document.getElementById('ig-hr-upscaler')?.value || '',
        clip_skip:       parseInt(document.getElementById('ig-clip-skip')?.value || '-1', 10),
        video_frames:    parseInt(document.getElementById('ig-video-frames')?.value || '33', 10),
        fps:             parseInt(document.getElementById('ig-video-fps')?.value || '16', 10),
        strength:        parseFloat(
            (currentMode === 'inpaint'
                ? document.getElementById('ig-strength-inpaint')
                : (currentMode === 'outpaint'
                    ? document.getElementById('ig-outpaint-strength')
                    : (currentMode === 'upscale'
                        ? document.getElementById('ig-upscale-denoise')
                        : document.getElementById('ig-strength')))
            )?.value || '0.75'
        ),
        outpaint_margin:    parseInt(document.getElementById('ig-outpaint-margin')?.value || '128', 10),
        outpaint_direction: document.querySelector('.ig-dir-btn.active')?.dataset?.dir || 'all',
        upscale_factor:     parseInt(document.getElementById('ig-upscale-factor')?.value || '2', 10),
        init_image_path: initImagePath,
        end_image_path:  endImagePath,
        mask_path:       maskPath,
        resize_mode:     document.getElementById(
            currentMode === 'inpaint' ? 'ig-resize-mode-inpaint' :
            currentMode === 'outpaint' ? 'ig-resize-mode-outpaint' : 'ig-resize-mode'
        )?.value || 'just_resize',
        style_strength:  parseFloat(document.getElementById(
            currentMode === 'inpaint' ? 'ig-style-strength-inpaint' :
            currentMode === 'outpaint' ? 'ig-style-strength-outpaint' : 'ig-style-strength'
        )?.value || '0.5'),
        preserve_colors: document.getElementById(
            currentMode === 'inpaint' ? 'ig-preserve-colors-inpaint' :
            currentMode === 'outpaint' ? 'ig-preserve-colors-outpaint' : 'ig-preserve-colors'
        )?.checked || false,
        match_histogram: document.getElementById(
            currentMode === 'inpaint' ? 'ig-match-histogram-inpaint' :
            currentMode === 'outpaint' ? 'ig-match-histogram-outpaint' : 'ig-match-histogram'
        )?.checked || false,
        latent_noise:    parseFloat(document.getElementById(
            currentMode === 'inpaint' ? 'ig-latent-noise-inpaint' :
            currentMode === 'outpaint' ? 'ig-latent-noise-outpaint' : 'ig-latent-noise'
        )?.value || '0'),
        mask_blur:       parseInt(document.getElementById(
            currentMode === 'inpaint' ? 'ig-mask-blur-inpaint' : 'ig-mask-blur'
        )?.value || '4', 10),
        mask_invert:     document.getElementById(
            currentMode === 'inpaint' ? 'ig-mask-invert' : 'ig-mask-invert'
        )?.checked || false,
        masked_content:  document.getElementById('ig-masked-content')?.value || 'fill',
        edge_behavior:   document.getElementById('ig-edge-behavior')?.value || 'mirror',
        upscale_method:  document.getElementById('ig-upscale-method')?.value || 'latent_sd',
        upscale_sharpen: document.getElementById('ig-upscale-sharpen')?.value || 'medium',
        upscale_anti_artifact: document.getElementById('ig-upscale-anti-artifact')?.value || 'off',
        prediction:      document.getElementById('ig-prediction')?.value || '',
        rng:             document.getElementById('ig-rng')?.value || '',
        sampler_rng:     document.getElementById('ig-sampler-rng')?.value || '',
        flow_shift:      parseFloat(document.getElementById('ig-flow-shift')?.value || '0'),
        vae_tile_overlap: parseFloat(document.getElementById('ig-vae-tile-overlap')?.value || '0.5'),
        disable_metadata: document.getElementById('ig-disable-metadata')?.checked || false,
        skip_layer_start: parseFloat(document.getElementById('ig-skip-layer-start')?.value || '0.01'),
        skip_layer_end:   parseFloat(document.getElementById('ig-skip-layer-end')?.value || '0.2'),
        high_noise_steps:          parseInt(document.getElementById('ig-high-noise-steps')?.value || '-1', 10),
        high_noise_cfg_scale:      parseFloat(document.getElementById('ig-high-noise-cfg')?.value || '7'),
        high_noise_guidance:       parseFloat(document.getElementById('ig-high-noise-guidance')?.value || '3.5'),
        high_noise_slg_scale:      parseFloat(document.getElementById('ig-high-noise-slg')?.value || '0'),
        high_noise_eta:            parseFloat(document.getElementById('ig-high-noise-eta')?.value || '0'),
        high_noise_skip_layer_start: parseFloat(document.getElementById('ig-high-noise-skip-start')?.value || '0.01'),
        high_noise_skip_layer_end:   parseFloat(document.getElementById('ig-high-noise-skip-end')?.value || '0.2'),
        // Video-specific
        output_format:     document.getElementById('ig-video-format')?.value || 'webm',
        video_seed_mode:  document.getElementById('ig-video-seed-mode')?.value || 'fixed',
        video_motion_strength:     parseFloat(document.getElementById('ig-video-motion')?.value || '0.5'),
        video_denoise:    parseFloat(document.getElementById('ig-video-denoise')?.value || '0.6'),
        video_loop:       document.getElementById('ig-video-loop')?.checked || false,
        // Global expert controls
        cfg_rescale:      parseFloat(document.getElementById('ig-cfg-rescale')?.value || '0'),
        face_restore:     parseFloat(document.getElementById('ig-face-restore')?.value || '0'),
        s_churn:          parseFloat(document.getElementById('ig-s-churn')?.value || '0'),
        s_tmin:           parseFloat(document.getElementById('ig-s-tmin')?.value || '0'),
        s_tmax:           parseFloat(document.getElementById('ig-s-tmax')?.value || '50'),
        s_noise:          parseFloat(document.getElementById('ig-s-noise')?.value || '1'),
        // Model startup flags forwarded to lazy-load path
        vae_tiling:       document.getElementById('ig-vae-tiling')?.checked || false,
        temporal_tiling:  document.getElementById('ig-temporal-tiling')?.checked || false,
        flash_attention:  document.getElementById('ig-fa')?.checked || false,
        offload_to_cpu:   document.getElementById('ig-offload-cpu')?.checked || false,
    };

    // VRAM check before generation
    const vramCheck = await checkVRAMForGeneration(opts);
    if (vramCheck.willExceedVRAM) {
        const proceed = confirm(
            'You don\'t have enough VRAM for this task.\n\n' +
            'Some of this will use your system RAM instead.\n' +
            'The image will take longer to appear, please be patient.\n\n' +
            `Estimated VRAM needed: ${vramCheck.estimatedVRAM} MB\n` +
            `Available VRAM: ${vramCheck.availableVRAM} MB\n\n` +
            'Click OK to proceed with CPU/RAM fallback, or Cancel to adjust settings.'
        );
        if (!proceed) {
            setStatus('Cancelled');
            return;
        }
    }

    let activeJobId = null;
    let pollInterval = null;
    let pollWatchdog = null;
    let lastPollMs = 0;
    let abortReject = null;
    progressReal = false;

    try {
        if (!window.go?.main?.App?.StartImageJob) {
            // Fallback to synchronous GenerateImage for older backends
            currentResult = await window.go.main.App.GenerateImage(opts);
            displayResult(currentResult);
            setStatus('Done!');
            updateDiagnostics('genDone');
            if (window.loadHistoryPanel) window.loadHistoryPanel();
            return;
        }

        // Async job flow — wrap in its own try so callback-registration
        // failures fall back to the synchronous path.
        try {
            activeJobId = await window.go.main.App.StartImageJob(opts);
        } catch (cbErr) {
            const msg = String(cbErr?.message || cbErr);
            throw cbErr;
        }
        setStatus('Generating...');
        showSpinner();

        // Update button to Abort
        btn.classList.add('generating');
        btn.innerHTML = '<span class="icon icon-close ig-fab-icon"></span><span class="ig-fab-text">Abort</span>';
        btn.disabled = false;
        btn.onclick = async () => {
            if (activeJobId && window.go?.main?.App?.CancelImageJob) {
                try { await window.go.main.App.CancelImageJob(activeJobId); } catch (e) { /* ignore */ }
            }
            if (pollInterval) { clearInterval(pollInterval); pollInterval = null; }
            if (abortReject) { abortReject(new Error('Cancelled')); abortReject = null; }
            resetGenerateButton();
            setStatus('Aborted');
            hideSpinner();
        };

        // Poll for status
        await new Promise((resolve, reject) => {
            abortReject = reject;
            pollInterval = setInterval(async () => {
                try {
                    const job = await window.go.main.App.GetImageJobStatus(activeJobId);
                    lastPollMs = Date.now();
                    if (!job) return;
                    if (job.status === 'completed') {
                        clearInterval(pollInterval);
                        pollInterval = null;
                        resolve(job.result);
                    } else if (job.status === 'failed') {
                        clearInterval(pollInterval);
                        pollInterval = null;
                        reject(new Error(job.error || 'Generation failed'));
                    } else if (job.status === 'cancelled') {
                        clearInterval(pollInterval);
                        pollInterval = null;
                        reject(new Error('Cancelled'));
                    } else if (job.status === 'loading') {
                        // Loading is a phase of the job now. A model can take minutes
                        // to come off disk, and saying so is the difference between a
                        // working app and one that looks hung.
                        setStatus('Loading model...');
                        updateProgressBar(0, 'Loading model...');
                    } else if (job.status === 'generating') {
                        const step = job.step || 0;
                        const total = job.total_steps || 1;
                        const pct = job.progress > 0 ? Math.min(job.progress * 100, 99) : Math.min((step / total) * 100, 99);
                        const label = step > 0 ? 'Step ' + step + ' of ' + total : 'Generating...';
                        updateProgressBar(pct, label);
                        if (job.preview_image) {
                            showPreviewImage(job.preview_image);
                        }
                    }
                    // else: queued / generating — keep polling
                } catch (pollErr) {
                    clearInterval(pollInterval);
                    pollInterval = null;
                    reject(pollErr);
                }
            }, 500);

            // Watchdog: one stalled status request must never leave the button on
            // "Abort" and the bar on "Generating..." forever. If nothing has come
            // back for 20s, fail with a message that says where the image went.
            lastPollMs = Date.now();
            pollWatchdog = setInterval(() => {
                if (Date.now() - lastPollMs < 20000) return;
                if (pollInterval) { clearInterval(pollInterval); pollInterval = null; }
                if (abortReject) {
                    const fail = abortReject;
                    abortReject = null;
                    fail(new Error('Lost contact with the image engine. The job may still finish in the background — check the gallery.'));
                }
            }, 2000);
        });

        // Job completed successfully — result is the resolved value
        const jobResult = await window.go.main.App.GetImageJobStatus(activeJobId);
        currentResult = jobResult.result;
        displayResult(currentResult);
        setStatus('Done!');
        updateDiagnostics('genDone');
        if (window.loadHistoryPanel) window.loadHistoryPanel();
    } catch (err) {
        console.error('Generation failed:', err);
        if (err.message !== 'Cancelled') {
            setStatus('Error: ' + (err.message || err));
            updateDiagnostics('error', err.message || String(err));
        }
    } finally {
        if (pollInterval) { clearInterval(pollInterval); pollInterval = null; }
        if (pollWatchdog) { clearInterval(pollWatchdog); pollWatchdog = null; }
        stopProgressBar();
        resetGenerateButton();
        hideSpinner();
        generationInFlight = false;
    }
}

function resetGenerateButton() {
    const btn = document.getElementById('ig-generate');
    if (!btn) return;
    btn.classList.remove('generating');
    btn.innerHTML = '<span class="icon icon-sparkle ig-fab-icon"></span><span class="ig-fab-text">Generate</span><span class="ig-key-hint">Ctrl+Enter</span>';
    btn.onclick = null;
    btn.disabled = false;
}

function showSpinner() {
    const area = document.getElementById('ig-preview-area');
    area.innerHTML = '<div class="ig-spinner"></div>';
}

function hideSpinner() {
    const area = document.getElementById('ig-preview-area');
    if (!area) return;
    area.querySelectorAll('.ig-spinner').forEach(el => el.remove());
}

let progressTimer = null;
let progressStart = 0;
let progressDuration = 0;
let progressSteps = 0;
let progressReal = false;

function startProgressBar(steps, family) {
    const wrap = document.getElementById('ig-progress-wrap');
    const bar = document.getElementById('ig-progress-bar');
    const text = document.getElementById('ig-progress-text');
    if (!wrap || !bar || !text) return;

    const msPerStep = { sd15: 200, sdxl: 400, sd3: 500, flux: 800, hidream: 1200 };
    progressSteps = steps;
    progressDuration = steps * (msPerStep[family] || 400);
    progressStart = performance.now();

    bar.style.width = '0%';
    bar.style.transition = 'none';
    wrap.style.display = '';
    text.textContent = 'Initializing...';

    progressTimer = setInterval(() => {
        if (progressReal) return;
        const elapsed = performance.now() - progressStart;
        let pct = Math.min((elapsed / progressDuration) * 100, 99);
        bar.style.width = pct + '%';

        const currentStep = Math.min(Math.floor((elapsed / progressDuration) * steps), steps);
        if (currentStep === 0) {
            text.textContent = 'Initializing...';
        } else {
            text.textContent = 'Step ' + currentStep + ' of ' + steps;
        }
    }, 100);
}

function stopProgressBar() {
    const wrap = document.getElementById('ig-progress-wrap');
    const bar = document.getElementById('ig-progress-bar');
    const text = document.getElementById('ig-progress-text');
    if (progressTimer) { clearInterval(progressTimer); progressTimer = null; }
    if (wrap) wrap.style.display = 'none';
    if (bar) bar.style.width = '0%';
    if (text) text.textContent = '';
}

function populateOptionsFromHistory(opts, clearModel) {
    const set = (id, val) => {
        const el = document.getElementById(id);
        if (!el || val === undefined || val === null) return;
        if (el.type === 'checkbox') el.checked = !!val;
        else if (el.tagName === 'SELECT') {
            // Only adopt a value the select actually offers. Assigning a value
            // with no matching option leaves selectedIndex at -1: the field
            // looks blank and the request carries '', so the engine silently
            // substitutes its own default for the user's choice.
            const wanted = String(val);
            if ([...el.options].some(o => o.value === wanted)) el.value = wanted;
        }
        else el.value = val;
    };
    set('ig-prompt', opts.prompt);
    set('ig-negative', opts.negative_prompt);
    set('ig-width', opts.width);
    set('ig-height', opts.height);
    set('ig-steps', opts.steps);
    set('ig-seed', opts.seed);
    set('ig-sampler', opts.sampler);
    set('ig-scheduler', opts.scheduler);
    set('ig-cfg', opts.cfg_scale);
    set('ig-guidance', opts.guidance);
    set('ig-slg-scale', opts.slg_scale);
    set('ig-eta', opts.eta);
    set('ig-batch', opts.batch_count);
    set('ig-clip-skip', opts.clip_skip);
    set('ig-vae-model', opts.vae_path);
    set('ig-clip-l', opts.clip_l_path);
    set('ig-clip-g', opts.clip_g_path);
    set('ig-t5xxl', opts.t5xxl_path);
    set('ig-taesd', opts.taesd);
    set('ig-strength', opts.strength);
    set('ig-strength-inpaint', opts.strength);
    set('ig-outpaint-strength', opts.strength);
    set('ig-outpaint-margin', opts.outpaint_margin);
    set('ig-upscale-factor', opts.upscale_factor);
    set('ig-control-model', opts.control_net_path);
    set('ig-control-strength', opts.control_strength);
    set('ig-video-frames', opts.video_frames);
    set('ig-video-fps', opts.fps);
    set('ig-prediction', opts.prediction);
    set('ig-rng', opts.rng);
    set('ig-sampler-rng', opts.sampler_rng);
    set('ig-flow-shift', opts.flow_shift);
    set('ig-vae-tile-overlap', opts.vae_tile_overlap);
    set('ig-disable-metadata', opts.disable_metadata);
    set('ig-skip-layer-start', opts.skip_layer_start);
    set('ig-skip-layer-end', opts.skip_layer_end);
    set('ig-high-noise-steps', opts.high_noise_steps);
    set('ig-high-noise-cfg', opts.high_noise_cfg_scale);
    set('ig-high-noise-guidance', opts.high_noise_guidance);
    set('ig-high-noise-slg', opts.high_noise_slg_scale);
    set('ig-high-noise-eta', opts.high_noise_eta);
    set('ig-high-noise-skip-start', opts.high_noise_skip_layer_start);
    set('ig-high-noise-skip-end', opts.high_noise_skip_layer_end);
    set('ig-backend', opts.backend);
    set('ig-params-backend', opts.params_backend);
    set('ig-vae-tiling', opts.vae_tiling);
    set('ig-temporal-tiling', opts.temporal_tiling);
    set('ig-circular', opts.circular);
    set('ig-fa', opts.flash_attention);
    set('ig-offload-cpu', opts.offload_to_cpu);
    set('ig-hires-enable', opts.enable_hr);
    set('ig-hr-scale', opts.hr_scale);
    set('ig-hr-upscaler', opts.hr_upscaler);
    if (opts.outpaint_direction && document.querySelectorAll('.ig-dir-btn')) {
        document.querySelectorAll('.ig-dir-btn').forEach(btn => {
            btn.classList.toggle('active', btn.dataset.dir === opts.outpaint_direction);
        });
    }
    if (opts.loras && Array.isArray(opts.loras)) {
        window.activeLoRAs = opts.loras.map(lora => {
            if (typeof lora === 'string') return { name: lora, strength: 1.0, enabled: true };
            return {
                name: lora.name || lora.Name || '',
                strength: lora.strength ?? lora.Strength ?? 1.0,
                enabled: lora.enabled !== false
            };
        }).filter(lora => lora.name);
        if (typeof renderActiveLoRAs === 'function') renderActiveLoRAs();
    }
    if (clearModel) set('ig-model', '');
    else if (opts.model_name) set('ig-model', opts.model_name);
    saveState();
}

function populateOptionsFromLegacy(meta, clearModel) {
    const set = (id, val) => {
        const el = document.getElementById(id);
        if (!el || val === undefined || val === null || val === '') return;
        if (el.tagName === 'SELECT' && ![...el.options].some(o => o.value === String(val))) return;
        el.value = val;
    };
    set('ig-prompt', meta.Prompt);
    set('ig-negative', meta.NegativePrompt);
    set('ig-width', meta.Width);
    set('ig-height', meta.Height);
    set('ig-steps', meta.Steps);
    set('ig-seed', meta.Seed);
    set('ig-sampler', meta.Sampler);
    set('ig-cfg', meta.CfgScale);
    if (clearModel) set('ig-model', '');
    else set('ig-model', meta.Model);
    saveState();
}

function triggerGenerate() {
    const btn = document.getElementById('ig-generate');
    if (btn && !btn.disabled) btn.click();
}

function updateProgressBar(percent, label) {
    progressReal = true;
    const bar = document.getElementById('ig-progress-bar');
    const text = document.getElementById('ig-progress-text');
    if (bar) {
        bar.style.transition = 'width 0.3s linear';
        bar.style.width = percent + '%';
    }
    if (text) text.textContent = label || '';
}

function showPreviewImage(b64) {
    const area = document.getElementById('ig-preview-area');
    if (!area || !b64) return;
    // sd-server's current_image is a JPEG, so a hardcoded image/png prefix lies
    // about the payload. Base64 magic is enough to pick the right one.
    let mime = 'image/png';
    if (b64.startsWith('/9j/')) mime = 'image/jpeg';
    else if (b64.startsWith('UklGR')) mime = 'image/webp';
    const prefix = b64.startsWith('data:') ? '' : 'data:' + mime + ';base64,';
    // TAESD frames are tiny (64-128 px for a 512-1024 px render), so without a
    // sizing class they paint at natural size in the middle of the canvas and
    // leave black bands around them. .ig-live-preview scales them to the canvas.
    area.innerHTML = `<img class="ig-live-preview" src="${prefix}${b64}" alt="Preview" style="opacity:0.85;filter:contrast(1.05);">`;
}

function displayResult(result) {
    const area = document.getElementById('ig-preview-area');
    const meta = document.getElementById('ig-meta');
    const actions = document.getElementById('ig-actions');

    if (result?.gallery_error) {
        if (typeof showNotification === 'function') showNotification(result.gallery_error, 'error');
        else if (typeof setStatus === 'function') setStatus(result.gallery_error);
    }

    if (result?.images && result.images.length > 1) {
        // Batch result: show grid of all images
        let gridHtml = '<div class="ig-batch-grid">';
        result.images.forEach((b64, idx) => {
            let mime = 'image/png';
            if (b64.startsWith('/9j/')) mime = 'image/jpeg';
            else if (b64.startsWith('UklGR')) mime = 'image/webp';
            gridHtml += `<img src="data:${mime};base64,${b64}" alt="Generated ${idx+1}" class="ig-batch-thumb" data-idx="${idx}" onclick="window.viewBatchImage(${idx})">`;
        });
        gridHtml += '</div>';
        area.innerHTML = gridHtml;
    } else if (result?.video_data) {
        if ((result.video_mime || '').includes('webp')) {
            area.innerHTML = `<img src="data:image/webp;base64,${result.video_data}" alt="Generated animation">`;
        } else {
            area.innerHTML = `<video controls autoplay loop style="max-width:100%;height:auto;border-radius:8px;" src="data:${result.video_mime || 'video/webm'};base64,${result.video_data}"></video>`;
        }
    } else if (result?.image_data) {
        // Detect MIME from base64 magic bytes, same as showPreviewImage.
        // sd-server can return PNG or JPEG depending on output_format; assume PNG
        // only if the magic bytes match, otherwise probe for JPEG/WebP.
        let mime = 'image/png';
        if (result.image_data.startsWith('/9j/')) mime = 'image/jpeg';
        else if (result.image_data.startsWith('UklGR')) mime = 'image/webp';
        area.innerHTML = `<img src="data:${mime};base64,${result.image_data}" alt="Generated">`;
    } else if (result?.metadata) {
        area.innerHTML = `<div class="ig-canvas-placeholder"><div class="ig-canvas-text">Image saved to gallery</div></div>`;
    } else {
        area.innerHTML = '<div class="ig-canvas-placeholder"><div class="ig-canvas-text">Generation completed</div></div>';
    }

    if (result?.metadata) {
        const m = result.metadata;
        meta.innerHTML = `
            <strong>Model:</strong> ${m.model} (${m.family})<br>
            <strong>Size:</strong> ${m.width}x${m.height} | <strong>Steps:</strong> ${m.steps}<br>
            <strong>Seed:</strong> ${m.seed} | <strong>Sampler:</strong> ${m.sampler}<br>
            <strong>CFG:</strong> ${m.cfg_scale} | <strong>Mode:</strong> ${m.mode}
        `;
        meta.style.display = 'block';
        actions.style.display = 'flex';
        // Add all batch images to history
        if (window.addHistoryItem) {
            if (result.images && result.images.length > 0) {
                result.images.forEach((b64, idx) => {
                    const file = result.metadata.file.replace('.png', '') + (idx > 0 ? `_${idx}.png` : '.png');
                    window.addHistoryItem(file, b64);
                });
            } else if (result.image_data) {
                window.addHistoryItem(result.metadata.file, result.image_data);
            }
        }
    }
}

// Batch image viewer: clicking a thumbnail shows it full-size in the preview area
window.viewBatchImage = function(idx) {
    if (!currentResult?.images?.[idx]) return;
    const area = document.getElementById('ig-preview-area');
    let mime = 'image/png';
    const b64 = currentResult.images[idx];
    if (b64.startsWith('/9j/')) mime = 'image/jpeg';
    else if (b64.startsWith('UklGR')) mime = 'image/webp';
    area.innerHTML = `<img src="data:${mime};base64,${b64}" alt="Generated ${idx+1}">`;
};

// ---------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------
async function saveImageAs() {
    if (!currentResult?.image_data && !(currentResult?.images?.length > 0)) {
        setStatus('No image data to save.'); return;
    }
    if (!window.go?.main?.App?.SaveImageAs) { setStatus('Save dialog not available.'); return; }
    try {
        setStatus('Saving...');
        // If batch, save all images; otherwise save the single image
        const imagesToSave = currentResult.images?.length > 0 ? currentResult.images : [currentResult.image_data];
        let paths = [];
        for (const b64 of imagesToSave) {
            const p = await window.go.main.App.SaveImageAs(b64);
            if (p) paths.push(p);
        }
        if (paths.length > 0) setStatus('Saved ' + paths.length + ' image(s).');
        else setStatus('Save cancelled.');
    } catch (err) {
        console.error('Save failed:', err);
        setStatus('Save failed: ' + (err.message || err));
    }
}

function deleteCurrent() {
    if (!currentResult?.metadata?.file) return;
    const baseName = currentResult.metadata.file.replace(/^.*[\\/]/, '').replace(/\.[^.]+$/, '');
    deleteGalleryImage(currentResult.metadata.file);
    currentResult = null;
    document.getElementById('ig-preview-area').innerHTML = `
        <div class="ig-canvas-placeholder">
            <div class="ig-canvas-icon">&#x2728;</div>
            <div class="ig-canvas-text">Generate images from text prompts.</div>
        </div>`;
    document.getElementById('ig-meta').style.display = 'none';
}

async function sendToEntity(entityType) {
    const imageToSend = currentResult?.images?.[0] || currentResult?.image_data;
    if (!imageToSend) { setStatus('No image to send.'); return; }
    if (!window.go?.main?.App?.SaveEntityImage) { setStatus('Entity image save not available.'); return; }

    // Ask user for entity name/ID (or use a default placeholder)
    const safeID = prompt(`Enter ${entityType} name/ID to save image for:`);
    if (!safeID) { setStatus('Send cancelled.'); return; }

    try {
        setStatus(`Saving image for ${entityType}...`);
        const path = await window.go.main.App.SaveEntityImage(imageToSend, entityType, safeID);
        setStatus(`Image saved for ${entityType}. Navigate to ${entityType} page to assign it.`);

        // Offer to navigate
        if (confirm(`Image saved. Open ${entityType} page now?`)) {
            const pageMap = { character: 'character.html', world: 'world.html', scenario: 'scenario.html' };
            window.location.href = `${pageMap[entityType] || 'index.html'}?pendingImage=${encodeURIComponent(path)}`;
        }
    } catch (err) {
        console.error('Send to entity failed:', err);
        setStatus('Failed: ' + (err.message || err));
    }
}

// ---------------------------------------------------------------------------
// Gallery Helpers
// ---------------------------------------------------------------------------
async function loadGalleryImage(filePath, imgId) {
    try {
        const b64 = await window.go.main.App.GetGalleryImage(filePath);
        const img = document.getElementById(imgId);
        if (img) img.src = 'data:image/png;base64,' + b64;
    } catch (err) {
        console.warn('Failed to load gallery image:', filePath, err);
    }
}

async function deleteGalleryImage(filePath) {
    const baseName = filePath.replace(/^.*[\\/]/, '').replace(/\.[^.]+$/, '');
    try {
        await window.go.main.App.DeleteGalleryImage(baseName);
        if (window.loadHistoryPanel) window.loadHistoryPanel();
    } catch (err) {
        console.error('Delete failed:', err);
    }
}

// ---------------------------------------------------------------------------
// Utilities
// ---------------------------------------------------------------------------
function escapeId(path) {
    return path.replace(/[^a-zA-Z0-9]/g, '_');
}

function escapeJS(str) {
    return str.replace(/\\/g, '\\\\').replace(/'/g, "\\'");
}

function setStatus(msg) {
    const el = document.getElementById('ig-status');
    if (el) el.textContent = msg;
}

// A select with no matching option reports value '' and selectedIndex -1: the
// field looks blank and the request carries an empty sampler, so the engine
// substitutes its own default. Nothing should ever read a sampler-less select,
// so re-establish a valid selection whenever restore paths may have cleared it.
function ensureSamplerSelected() {
    const sel = document.getElementById('ig-sampler');
    if (!sel || !sel.options || sel.options.length === 0) return;
    if (sel.selectedIndex >= 0 && sel.value !== '') return;
    const preferred = [...sel.options].find(o => o.defaultSelected) ||
                      [...sel.options].find(o => o.value === 'euler_a') ||
                      sel.options[0];
    if (preferred) {
        sel.value = preferred.value;
        console.log('[IG] sampler selection repaired to', sel.value);
    }
}

function readSamplerValue() {
    ensureSamplerSelected();
    const el = document.getElementById('ig-sampler');
    return el ? el.value : '';
}

// Settings restore is asynchronous, so a repair at load can be undone a
// moment later. Re-check a few times as restore settles.
function repairSamplerAfterRestore() {
    [0, 250, 1000, 2500].forEach(ms => setTimeout(ensureSamplerSelected, ms));
}
window.addEventListener('load', repairSamplerAfterRestore);
document.addEventListener('DOMContentLoaded', repairSamplerAfterRestore);
document.addEventListener('change', (e) => {
    if (e.target && e.target.id === 'ig-model') setTimeout(ensureSamplerSelected, 0);
});
