'use strict';

const IG_STATE_KEY = 'howl-image-gen-state';
let currentResult = null;
let models = [];       // full ModelInfo list (diffusion only)
let scanIndex = null;  // full ScanResult from ScanImageModels
let seedHistory = [];
let lastUsedSeed = -1;
let currentModelFamily = '';
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
        if (app.ScanImageModels) {
            scanIndex = await app.ScanImageModels();
        } else if (app.GetModelIndex) {
            scanIndex = await app.GetModelIndex();
        } else {
            // Legacy fallback
            const legacy = await app.ListImageModels();
            scanIndex = { models: legacy || [], vaes: [], loras: [], controlnets: [], embeddings: [], upscalers: [] };
        }
        if (!scanIndex) scanIndex = { models: [], vaes: [], loras: [], controlnets: [], embeddings: [], upscalers: [] };

        models = scanIndex.models || [];
        console.log('[IG] Scan:', models.length, 'models,', (scanIndex.vaes||[]).length, 'VAEs,',
                    (scanIndex.loras||[]).length, 'LoRAs,', (scanIndex.controlnets||[]).length, 'ControlNets');

        _populateDiffusionDropdown();
        _populateVAEDropdown(scanIndex.vaes || []);
        _populateLoRAList(scanIndex.loras || []);
        _populateControlNetDropdown(scanIndex.controlnets || []);

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
    const ph = document.createElement('option');
    ph.value = '';
    ph.disabled = true;
    ph.selected = true;
    ph.textContent = models.length === 0 ? 'No models found' : 'Select a model...';
    sel.appendChild(ph);
    models.forEach(m => {
        const opt = document.createElement('option');
        opt.value = m.name;
        const tag = m.format === 'bundle' ? '[bundle]' : m.format ? `[${m.format}]` : '';
        opt.textContent = `${m.name} (${m.family || '?'}) ${tag}`.trim();
        opt.dataset.family = m.family || '';
        opt.dataset.recommended = JSON.stringify(m.recommended || {});
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
    const sel = document.getElementById('ig-lora');
    if (!sel) return;
    sel.innerHTML = '<option value="">None</option>';
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
}

// onModelChange updates UI presets and re-filters secondary dropdowns by family.
function onModelChange() {
    const name = document.getElementById('ig-model').value;
    if (!name) return;
    const model = models.find(m => m.name === name);
    if (model) {
        currentModelFamily = model.family || '';
        window.currentModelFamily = currentModelFamily;
        applyModelPreset(model.family, model.recommended || null);
        updateLiveModelInfo(model);
        // Re-filter VAE / LoRA / ControlNet for this family
        if (currentModelFamily) refreshCompatibleAssets(currentModelFamily);
    }
    const statusEl = document.getElementById('ig-model-status');
    if (statusEl) { statusEl.textContent = 'Auto-load on Generate'; statusEl.className = 'ig-panel-status'; }
    saveState();
}

// setModel manually loads a model into the backend (used by Reload button).
async function setModel(name) {
    const statusEl = document.getElementById('ig-model-status');
    if (statusEl) { statusEl.textContent = 'Loading...'; statusEl.className = 'ig-panel-status ig-status-loading'; }
    updateDiagnostics('loadStart');
    setStatus('Loading model: ' + name + '...');
    try {
        const opts = {
            model: name,
            vae: document.getElementById('ig-vae-model')?.value || '',
            clip_l: document.getElementById('ig-clip-l')?.value || '',
            clip_g: document.getElementById('ig-clip-g')?.value || '',
            t5xxl: document.getElementById('ig-t5xxl')?.value || '',
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
        const model = models.find(m => m.name === name);
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
        setStatus('Error loading model: ' + (err.message || err));
        if (statusEl) { statusEl.textContent = 'Error'; statusEl.className = 'ig-panel-status ig-status-error'; }
        updateDiagnostics('error', err.message || String(err));
    }
}

function applyModelPreset(family, recommended) {
    // Prefer scan-manifest recommended values; fall back to family defaults.
    const rec = recommended || _familyDefaults(family);
    const set = (id, val) => { const el = document.getElementById(id); if (el && val) el.value = val; };
    set('ig-steps',  rec.steps);
    set('ig-cfg',    rec.cfg);
    set('ig-width',  rec.width);
    set('ig-height', rec.height);
    if (rec.sampler) set('ig-sampler', rec.sampler);
    updateLabels();
    saveState();
}

function _familyDefaults(family) {
    switch ((family || '').toLowerCase()) {
        case 'flux':
            return { steps: 4,  cfg: 1,   width: 1024, height: 1024, sampler: 'euler' };
        case 'sd1.5': case 'sd15':
            return { steps: 20, cfg: 7,   width: 512,  height: 512,  sampler: 'euler_a' };
        case 'sd2.1': case 'sd21':
            return { steps: 20, cfg: 7.5, width: 768,  height: 768,  sampler: 'dpmpp_2m' };
        case 'sdxl-turbo': case 'turbo':
            return { steps: 4,  cfg: 1,   width: 512,  height: 512,  sampler: 'euler_a' };
        default: // sdxl, sd3, unknown
            return { steps: 20, cfg: 7,   width: 1024, height: 1024, sampler: 'dpmpp_2m' };
    }
}

// ---------------------------------------------------------------------------
// Generation
// ---------------------------------------------------------------------------
async function onGenerate() {
    const prompt = document.getElementById('ig-prompt').value.trim();
    if (!prompt) { setStatus('Please enter a prompt.'); return; }
    const modelSel = document.getElementById('ig-model');
    if (!modelSel.value) { setStatus('Please select a model.'); return; }

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

    const opts = {
        prompt:          prompt,
        negative_prompt: document.getElementById('ig-negative').value,
        width:           parseInt(document.getElementById('ig-width').value, 10),
        height:          parseInt(document.getElementById('ig-height').value, 10),
        steps:           parseInt(document.getElementById('ig-steps').value, 10),
        seed:            lastUsedSeed,
        sampler:         document.getElementById('ig-sampler').value,
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
        control_strength: parseFloat(document.getElementById('ig-control-strength')?.value || '0.9'),
        loras:           Array.from(document.getElementById('ig-lora')?.selectedOptions || []).map(o => o.value).filter(v => v),
        enable_hr:       document.getElementById('ig-hires-enable')?.checked || false,
        hr_scale:        parseFloat(document.getElementById('ig-hr-scale')?.value || '2'),
        hr_upscaler:     document.getElementById('ig-hr-upscaler')?.value || '',
        tiling:          document.getElementById('ig-tiling')?.checked || false,
        clip_skip:       parseInt(document.getElementById('ig-clip-skip')?.value || '-1', 10),
        cache_mode:      document.getElementById('ig-cache-mode')?.value || 'none',
        cache_option:    document.getElementById('ig-cache-option')?.value || '',
        video_frames:    parseInt(document.getElementById('ig-video-frames')?.value || '24', 10),
        fps:             parseInt(document.getElementById('ig-video-fps')?.value || '12', 10),
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
        mask_path:       maskPath,
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
        // Model startup flags forwarded to lazy-load path
        vae_tiling:       document.getElementById('ig-vae-tiling')?.checked || false,
        temporal_tiling:  document.getElementById('ig-temporal-tiling')?.checked || false,
        flash_attention:  document.getElementById('ig-fa')?.checked || false,
        offload_to_cpu:   document.getElementById('ig-offload-cpu')?.checked || false,
    };

    let activeJobId = null;
    let pollInterval = null;
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
            if (msg.includes('not registered') || msg.includes('callback')) {
                console.warn('[IG] Async job callback not registered, falling back to sync GenerateImage');
                currentResult = await window.go.main.App.GenerateImage(opts);
                displayResult(currentResult);
                setStatus('Done!');
                updateDiagnostics('genDone');
                if (window.loadHistoryPanel) window.loadHistoryPanel();
                return;
            }
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
        stopProgressBar();
        resetGenerateButton();
        hideSpinner();
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

function hideSpinner() {}

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

    const msPerStep = { sd15: 200, sdxl: 400, sd3: 500, flux: 800 };
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
    const prefix = b64.startsWith('data:') ? '' : 'data:image/png;base64,';
    area.innerHTML = `<img src="${prefix}${b64}" alt="Preview" style="opacity:0.85;filter:contrast(1.05);">`;
}

function displayResult(result) {
    const area = document.getElementById('ig-preview-area');
    const meta = document.getElementById('ig-meta');
    const actions = document.getElementById('ig-actions');

    if (result?.images && result.images.length > 1) {
        // Batch result: show grid of all images
        let gridHtml = '<div class="ig-batch-grid">';
        result.images.forEach((b64, idx) => {
            gridHtml += `<img src="data:image/png;base64,${b64}" alt="Generated ${idx+1}" class="ig-batch-thumb" data-idx="${idx}" onclick="window.viewBatchImage(${idx})">`;
        });
        gridHtml += '</div>';
        area.innerHTML = gridHtml;
    } else if (result?.image_data) {
        area.innerHTML = `<img src="data:image/png;base64,${result.image_data}" alt="Generated">`;
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
    area.innerHTML = `<img src="data:image/png;base64,${currentResult.images[idx]}" alt="Generated ${idx+1}">`;
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

function toggleMetadata() {
    const viewer = document.getElementById('ig-metadata-viewer');
    if (!viewer) return;
    if (viewer.style.display === 'none') {
        if (currentResult?.metadata) {
            viewer.innerHTML = '<pre>' + JSON.stringify(currentResult.metadata, null, 2) + '</pre>';
        } else {
            viewer.innerHTML = '<div class="ig-mode-note">No metadata available for this image.</div>';
        }
        viewer.style.display = 'block';
    } else {
        viewer.style.display = 'none';
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
