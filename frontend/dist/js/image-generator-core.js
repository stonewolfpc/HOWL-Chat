'use strict';

const IG_STATE_KEY = 'howl-image-gen-state';
let currentResult = null;
let models = [];
let seedHistory = [];
let lastUsedSeed = -1;
let currentModelFamily = '';

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
async function initPaths() {
    const settings = window.settingsState?.current?.imageGen || {};
    const modelsDir = settings.modelsDir || '';
    const sdPath    = settings.sdServerPath || '';

    document.getElementById('ig-models-dir').value = modelsDir;
    document.getElementById('ig-sd-server').value = sdPath;

    if (modelsDir && sdPath && window.go?.main?.App) {
        try {
            await window.go.main.App.SetImageModelsDir(modelsDir);
            await window.go.main.App.SetSDServerPath(sdPath);
            setStatus('Engine initialized with saved paths');
            await loadModels();
            if (window.go.main.App.GetImageEngineConfig) {
                const cfg = await window.go.main.App.GetImageEngineConfig();
                console.log('[IG] Engine config after init:', cfg);
            }
            return;
        } catch (err) {
            console.warn('Saved paths failed, falling back to auto-detect:', err);
        }
    }
    loadModels();
}

async function browseModelsDir() {
    const input = document.getElementById('ig-models-dir');
    if (window.go?.main?.App?.OpenImageModelsDirPicker) {
        try {
            const dir = await window.go.main.App.OpenImageModelsDirPicker();
            if (dir) input.value = dir;
            return;
        } catch (err) {
            console.error('Browse failed:', err);
        }
    }
    input.focus();
    input.select();
}

async function browseSDServer() {
    const input = document.getElementById('ig-sd-server');
    if (window.go?.main?.App?.OpenSDServerFilePicker) {
        try {
            const path = await window.go.main.App.OpenSDServerFilePicker();
            if (path) input.value = path;
            return;
        } catch (err) {
            console.error('Browse failed:', err);
        }
    }
    input.focus();
    input.select();
}

async function applyPaths() {
    const modelsDir = document.getElementById('ig-models-dir').value.trim();
    const sdPath    = document.getElementById('ig-sd-server').value.trim();
    if (!modelsDir || !sdPath) {
        setStatus('Please set both paths.');
        return;
    }
    setStatus('Applying paths...');
    try {
        await window.go.main.App.SetImageModelsDir(modelsDir);
        await window.go.main.App.SetSDServerPath(sdPath);
        if (window.settingsState?.current) {
            window.settingsState.current.imageGen = window.settingsState.current.imageGen || {};
            window.settingsState.current.imageGen.modelsDir = modelsDir;
            window.settingsState.current.imageGen.sdServerPath = sdPath;
            if (window.go?.main?.App?.SaveSettings) {
                await window.go.main.App.SaveSettings(window.settingsState.current);
            }
        }
        if (window.go.main.App.GetImageEngineConfig) {
            const cfg = await window.go.main.App.GetImageEngineConfig();
            console.log('[IG] Engine config after apply:', cfg);
        }
        await loadModels();
        setStatus('Ready');
    } catch (err) {
        setStatus('Error: ' + (err.message || err));
        console.error('Apply paths failed:', err);
    }
}

// ---------------------------------------------------------------------------
// Model Loading
// ---------------------------------------------------------------------------
async function loadModels() {
    if (!window.go?.main?.App?.ListImageModels) {
        console.warn('[IG] ListImageModels binding not available');
        return;
    }
    if (window.go?.main?.App?.DebugImageEngine) {
        try {
            const dbg = await window.go.main.App.DebugImageEngine();
            console.log('[IG] DebugImageEngine:', dbg);
        } catch (e) {
            console.warn('[IG] DebugImageEngine failed:', e);
        }
    }
    try {
        models = await window.go.main.App.ListImageModels();
        console.log('[IG] ListImageModels returned:', models);
        if (!models) models = [];
        populateModelDropdowns();
        await loadLoRAs();
        if (models.length === 0) {
            setStatus('No models found. Check Engine Paths.');
            return;
        }
        setStatus(`${models.length} model(s) discovered. Select one and click Load.`);
        const saved = JSON.parse(localStorage.getItem(IG_STATE_KEY) || '{}');
        if (saved.model && models.some(m => m.name === saved.model)) {
            document.getElementById('ig-model').value = saved.model;
        }
    } catch (err) {
        console.error('Failed to load models:', err);
        setStatus('Error loading models: ' + err.message);
        document.getElementById('ig-model').innerHTML = '<option value="">Error loading models</option>';
    }
}

async function loadLoRAs() {
    const sel = document.getElementById('ig-lora');
    if (!sel) return; // LoRA select not present in current UI
    if (!window.go?.main?.App?.ListImageLoRAs) {
        sel.innerHTML = '<option value="">LoRA scanning not available</option>';
        return;
    }
    try {
        const loras = await window.go.main.App.ListImageLoRAs();
        if (!loras) { sel.innerHTML = '<option value="">None</option>'; return; }
        sel.innerHTML = '<option value="">None</option>';
        loras.forEach(name => {
            const opt = document.createElement('option');
            opt.value = name;
            opt.textContent = name;
            sel.appendChild(opt);
        });
    } catch (err) {
        console.error('Failed to load LoRAs:', err);
        sel.innerHTML = '<option value="">Error loading LoRAs</option>';
    }
}

function populateModelDropdowns() {
    const ids = ['ig-model', 'ig-vae-model', 'ig-clip-l', 'ig-clip-g', 'ig-t5xxl', 'ig-control-model'];
    ids.forEach(id => {
        const sel = document.getElementById(id);
        if (!sel) return;
        sel.innerHTML = '';
        const placeholder = document.createElement('option');
        placeholder.value = '';
        placeholder.textContent = models.length === 0 ? 'No models found' : (id === 'ig-model' ? 'Select a model...' : 'None (Internal)');
        if (id === 'ig-model') {
            placeholder.disabled = true;
            placeholder.selected = true;
        }
        sel.appendChild(placeholder);
        models.forEach(m => {
            const opt = document.createElement('option');
            opt.value = m.name;
            const format = m.format || 'unknown';
            const compat = m.compatibility || 'unknown';
            let label = `${m.name} (${m.family})`;
            if (compat === 'unknown') label += ' [GGUF]';
            if (compat === 'incompatible') label += ' [!]';
            if (compat === 'verified') label += ` [${format}]`;
            opt.textContent = label;
            sel.appendChild(opt);
        });
    });
}

async function setModel(name) {
    const statusEl = document.getElementById('ig-model-status');
    if (statusEl) { statusEl.textContent = 'Loading...'; statusEl.className = 'ig-panel-status ig-status-loading'; }
    updateDiagnostics('loadStart');
    setStatus('Loading model: ' + name + '...');
    try {
        // Prepare advanced model config if present
        const opts = {
            model: name,
            vae: document.getElementById('ig-vae-model')?.value || '',
            clip_l: document.getElementById('ig-clip-l')?.value || '',
            clip_g: document.getElementById('ig-clip-g')?.value || '',
            t5xxl: document.getElementById('ig-t5xxl')?.value || '',
        };
        console.log('[IG] Loading model with config:', opts);
        
        await window.go.main.App.SetActiveImageModel(opts);
        
        const model = models.find(m => m.name === name);
        if (model) {
            currentModelFamily = model.family || '';
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

async function onLoadModelClick() {
    const name = document.getElementById('ig-model').value;
    if (!name) { setStatus('Please select a model first.'); return; }
    const model = models.find(m => m.name === name);
    if (model && model.compatibility === 'unknown') {
        if (!confirm('This model is a GGUF file with unknown compatibility. It may fail to load.\n\nTry loading it anyway?')) return;
    }
    if (model && model.compatibility === 'incompatible') {
        setStatus('This model format is incompatible. Use .safetensors or .ckpt instead.');
        return;
    }
    await setModel(name);
}

function applyModelPreset(family) {
    if (family === 'flux') {
        document.getElementById('ig-steps').value = 4;
        document.getElementById('ig-cfg').value = 1;
        document.getElementById('ig-width').value = 1024;
        document.getElementById('ig-height').value = 1024;
    } else if (family === 'sd15') {
        document.getElementById('ig-steps').value = 20;
        document.getElementById('ig-cfg').value = 7;
        document.getElementById('ig-width').value = 512;
        document.getElementById('ig-height').value = 512;
    } else if (family === 'sdxl' || family === 'sd3') {
        document.getElementById('ig-steps').value = 20;
        document.getElementById('ig-cfg').value = 7;
        document.getElementById('ig-width').value = 1024;
        document.getElementById('ig-height').value = 1024;
    }
    updateLabels();
    saveState();
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
    if (seedVal === -1 || isNaN(seedVal)) {
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
        tiling:          document.getElementById('ig-tiling')?.checked || false,
        clip_skip:       parseInt(document.getElementById('ig-clip-skip')?.value || '-1', 10),
        cache_mode:      document.getElementById('ig-cache-mode')?.value || 'none',
        cache_option:    document.getElementById('ig-cache-option')?.value || '',
        video_frames:    parseInt(document.getElementById('ig-video-frames')?.value || '24', 10),
        fps:             parseInt(document.getElementById('ig-video-fps')?.value || '12', 10),
        strength:        parseFloat(
            (currentMode === 'inpaint'
                ? document.getElementById('ig-strength-inpaint')
                : (currentMode === 'outpaint' ? document.getElementById('ig-outpaint-strength') : document.getElementById('ig-strength'))
            )?.value || '0.75'
        ),
        init_image_path: initImagePath,
        mask_path:       maskPath,
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

        // Async job flow
        activeJobId = await window.go.main.App.StartImageJob(opts);
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

    if (result?.image_data) {
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
        if (window.addHistoryItem) window.addHistoryItem(result.metadata.file, result.image_data);
    }
}

// ---------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------
async function saveImageAs() {
    if (!currentResult?.image_data) { setStatus('No image data to save.'); return; }
    if (!window.go?.main?.App?.SaveImageAs) { setStatus('Save dialog not available.'); return; }
    try {
        setStatus('Saving...');
        const savedPath = await window.go.main.App.SaveImageAs(currentResult.image_data);
        if (savedPath) setStatus('Saved to: ' + savedPath);
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
    document.getElementById('ig-actions').style.display = 'none';
    document.querySelectorAll('.ig-history-item').forEach(t => t.classList.remove('active'));
}

async function sendToEntity(entityType) {
    if (!currentResult?.image_data) { setStatus('No image to send.'); return; }
    if (!window.go?.main?.App?.SaveEntityImage) { setStatus('Entity image save not available.'); return; }

    // Ask user for entity name/ID (or use a default placeholder)
    const entityName = prompt(`Enter ${entityType} name (or ID):`)?.trim();
    if (!entityName) { setStatus('Cancelled.'); return; }

    const safeID = entityName.toLowerCase().replace(/[^a-z0-9_-]/g, '_').substring(0, 40);

    try {
        setStatus(`Saving image for ${entityType}...`);
        const path = await window.go.main.App.SaveEntityImage(currentResult.image_data, entityType, safeID);
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
