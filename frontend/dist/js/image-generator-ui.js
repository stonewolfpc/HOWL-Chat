'use strict';

// ---------------------------------------------------------------------------
// State Persistence
// ---------------------------------------------------------------------------
function saveState() {
    const state = {
        model:     document.getElementById('ig-model').value,
        vae:       document.getElementById('ig-vae-model').value,
        clipL:     document.getElementById('ig-clip-l').value,
        clipG:     document.getElementById('ig-clip-g').value,
        t5xxl:     document.getElementById('ig-t5xxl').value,
        taesd:     document.getElementById('ig-taesd').checked,
        prompt:    document.getElementById('ig-prompt').value,
        negative:  document.getElementById('ig-negative').value,
        width:     document.getElementById('ig-width').value,
        height:    document.getElementById('ig-height').value,
        steps:     document.getElementById('ig-steps').value,
        seed:      document.getElementById('ig-seed').value,
        sampler:   document.getElementById('ig-sampler').value,
        scheduler: document.getElementById('ig-scheduler').value,
        cfg:       document.getElementById('ig-cfg').value,
        guidance:  document.getElementById('ig-guidance').value,
        slgScale:  document.getElementById('ig-slg-scale').value,
        eta:       document.getElementById('ig-eta').value,
        batch:     document.getElementById('ig-batch').value,
        clipSkip:  document.getElementById('ig-clip-skip').value,
        lora:      Array.from(document.getElementById('ig-lora')?.selectedOptions || []).map(o => o.value).filter(v => v),
        tiling:    document.getElementById('ig-tiling').checked,
        hires:     document.getElementById('ig-hires-enable').checked,
        ctrlModel: document.getElementById('ig-control-model').value,
        ctrlStrength: document.getElementById('ig-control-strength').value,
        cacheMode: document.getElementById('ig-cache-mode').value,
        cacheOption: document.getElementById('ig-cache-option').value,
        videoFrames: document.getElementById('ig-video-frames')?.value,
        videoFPS:    document.getElementById('ig-video-fps')?.value,
    };
    localStorage.setItem(IG_STATE_KEY, JSON.stringify(state));
}

function restoreState() {
    const raw = localStorage.getItem(IG_STATE_KEY);
    if (!raw) return;
    try {
        const s = JSON.parse(raw);
        if (s.prompt)    document.getElementById('ig-prompt').value    = s.prompt;
        if (s.negative)  document.getElementById('ig-negative').value  = s.negative;
        if (s.width)     document.getElementById('ig-width').value    = s.width;
        if (s.height)    document.getElementById('ig-height').value   = s.height;
        if (s.steps)     document.getElementById('ig-steps').value     = s.steps;
        if (s.seed)      document.getElementById('ig-seed').value      = s.seed;
        if (s.sampler)   document.getElementById('ig-sampler').value    = s.sampler;
        if (s.scheduler) document.getElementById('ig-scheduler').value  = s.scheduler;
        if (s.cfg)       document.getElementById('ig-cfg').value       = s.cfg;
        if (s.guidance)  document.getElementById('ig-guidance').value  = s.guidance;
        if (s.slgScale)  document.getElementById('ig-slg-scale').value  = s.slgScale;
        if (s.eta)       document.getElementById('ig-eta').value       = s.eta;
        if (s.batch)     document.getElementById('ig-batch').value       = s.batch;
        if (s.clipSkip)  document.getElementById('ig-clip-skip').value   = s.clipSkip;
        if (s.vae)       document.getElementById('ig-vae-model').value = s.vae;
        if (s.clipL)     document.getElementById('ig-clip-l').value    = s.clipL;
        if (s.clipG)     document.getElementById('ig-clip-g').value    = s.clipG;
        if (s.t5xxl)     document.getElementById('ig-t5xxl').value    = s.t5xxl;
        if (s.taesd !== undefined) document.getElementById('ig-taesd').checked = s.taesd;
        if (s.lora && Array.isArray(s.lora)) {
            const loraSel = document.getElementById('ig-lora');
            if (loraSel) Array.from(loraSel.options).forEach(opt => { opt.selected = s.lora.includes(opt.value); });
        }
        if (s.tiling !== undefined) document.getElementById('ig-tiling').checked = s.tiling;
        if (s.hires !== undefined)  document.getElementById('ig-hires-enable').checked = s.hires;
        if (s.ctrlModel) document.getElementById('ig-control-model').value = s.ctrlModel;
        if (s.ctrlStrength) document.getElementById('ig-control-strength').value = s.ctrlStrength;
        if (s.cacheMode) document.getElementById('ig-cache-mode').value = s.cacheMode;
        if (s.cacheOption) document.getElementById('ig-cache-option').value = s.cacheOption;
        if (s.videoFrames && document.getElementById('ig-video-frames')) document.getElementById('ig-video-frames').value = s.videoFrames;
        if (s.videoFPS && document.getElementById('ig-video-fps')) document.getElementById('ig-video-fps').value = s.videoFPS;
        updateLabels();
    } catch (e) { /* ignore */ }
}

function updateLabels() {
    const ids = ['ig-steps', 'ig-cfg', 'ig-strength', 'ig-strength-inpaint', 'ig-outpaint-strength', 
                 'ig-upscale-denoise', 'ig-guidance', 'ig-slg-scale', 'ig-eta', 
                 'ig-control-strength', 'ig-video-frames', 'ig-video-fps'];
    ids.forEach(id => {
        const slider = document.getElementById(id);
        const valEl = document.getElementById(id + '-val');
        if (slider && valEl) valEl.textContent = slider.value;
    });
}

// ---------------------------------------------------------------------------
// Event Binding
// ---------------------------------------------------------------------------
function bindEvents() {
    // Sliders
    const sliderIds = ['ig-steps', 'ig-cfg', 'ig-strength', 'ig-strength-inpaint', 'ig-outpaint-strength', 
                       'ig-upscale-denoise', 'ig-guidance', 'ig-slg-scale', 'ig-eta', 
                       'ig-control-strength', 'ig-video-frames', 'ig-video-fps', 'ig-brush-size'];
    sliderIds.forEach(id => {
        const slider = document.getElementById(id);
        if (slider) {
            slider.addEventListener('input', () => {
                const valEl = document.getElementById(id + '-val');
                if (valEl) valEl.textContent = slider.value;
                saveState();
            });
        }
    });

    // Aspect ratio buttons
    document.querySelectorAll('.ig-aspect-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            document.querySelectorAll('.ig-aspect-btn').forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
            document.getElementById('ig-width').value  = btn.dataset.w;
            document.getElementById('ig-height').value = btn.dataset.h;
            saveState();
        });
    });

    // Inputs that trigger save
    ['ig-model', 'ig-vae-model', 'ig-clip-l', 'ig-clip-g', 'ig-t5xxl', 'ig-taesd',
     'ig-prompt', 'ig-negative', 'ig-width', 'ig-height', 'ig-seed',
     'ig-sampler', 'ig-scheduler', 'ig-batch', 'ig-clip-skip', 'ig-lora', 'ig-tiling',
     'ig-hires-enable', 'ig-control-model', 'ig-cache-mode', 'ig-cache-option'].forEach(id => {
        const el = document.getElementById(id);
        if (el) el.addEventListener('change', saveState);
    });

    // File upload for img2img
    const uploadInput = document.getElementById('ig-upload-input');
    const uploadZone = document.getElementById('ig-upload-zone');
    if (uploadInput && uploadZone) {
        uploadZone.addEventListener('click', () => {
            uploadInput.dataset.context = 'ig-upload-zone';
            uploadInput.click();
        });
        uploadInput.addEventListener('change', (e) => {
            if (e.target.files?.[0]) {
                const context = uploadInput.dataset.context;
                if (context === 'ig-upload-controlnet') {
                    handleControlNetUpload(e.target.files[0]);
                } else {
                    handleImageUpload(e.target.files[0]);
                }
            }
        });
        // Drag and drop
        uploadZone.addEventListener('dragover', (e) => {
            e.preventDefault();
            uploadZone.classList.add('dragover');
        });
        uploadZone.addEventListener('dragleave', () => uploadZone.classList.remove('dragover'));
        uploadZone.addEventListener('drop', (e) => {
            e.preventDefault();
            uploadZone.classList.remove('dragover');
            if (e.dataTransfer.files?.[0]) handleImageUpload(e.dataTransfer.files[0]);
        });
    }

    // Mini upload zones in Mode Tools panel
    document.querySelectorAll('.ig-upload-mini').forEach(zone => {
        zone.addEventListener('click', () => {
            if (uploadInput) {
                // Set a context on the input so we know what we're uploading for
                input.dataset.context = zone.id;
                input.click();
            }
        });
    });

    // Brush controls for inpaint (now in left panel Mode Tools)
    const brushSize = document.getElementById('ig-brush-size');
    if (brushSize) {
        brushSize.addEventListener('input', () => {
            const valEl = document.getElementById('ig-brush-val');
            if (valEl) valEl.textContent = brushSize.value;
        });
    }
    const clearMaskBtn = document.getElementById('ig-clear-mask');
    if (clearMaskBtn) clearMaskBtn.addEventListener('click', clearMask);
    initMaskDrawing();

    // Generate
    document.getElementById('ig-generate').addEventListener('click', onGenerate);

    // Load Model button
    const loadBtn = document.getElementById('ig-load-model');
    if (loadBtn) loadBtn.addEventListener('click', onLoadModelClick);

    // Unload Model button
    const unloadBtn = document.getElementById('ig-unload-model');
    if (unloadBtn) unloadBtn.addEventListener('click', onUnloadModel);

    // Reload Model button
    const reloadBtn = document.getElementById('ig-reload-model');
    if (reloadBtn) reloadBtn.addEventListener('click', onReloadModel);

    // Panel toggles
    document.querySelectorAll('.ig-panel-header').forEach(header => {
        header.addEventListener('click', () => togglePanel(header.closest('.ig-panel')));
    });

    // Workflow mode buttons
    document.querySelectorAll('.ig-workflow-btn').forEach(btn => {
        btn.addEventListener('click', () => setWorkflowMode(btn.dataset.mode));
    });

    // Seed lock
    const lockBtn = document.getElementById('ig-seed-lock');
    if (lockBtn) lockBtn.addEventListener('click', toggleSeedLock);

    // Seed reuse
    const reuseBtn = document.getElementById('ig-seed-reuse');
    if (reuseBtn) reuseBtn.addEventListener('click', reuseLastSeed);

    // Keyboard shortcuts
    document.addEventListener('keydown', handleKeydown);
}

// ---------------------------------------------------------------------------
// Panel Toggles
// ---------------------------------------------------------------------------
function togglePanel(panel) {
    if (!panel) return;
    const body = panel.querySelector('.ig-panel-body');
    const toggle = panel.querySelector('.ig-panel-toggle');
    if (!body) return;
    const isCollapsed = panel.classList.contains('ig-panel-collapsed');
    if (isCollapsed) {
        panel.classList.remove('ig-panel-collapsed');
        body.style.display = '';
        if (toggle) toggle.innerHTML = '&#x25BC;';
    } else {
        panel.classList.add('ig-panel-collapsed');
        body.style.display = 'none';
        if (toggle) toggle.innerHTML = '&#x25B6;';
    }
}

// ---------------------------------------------------------------------------
// Workflow Modes
// ---------------------------------------------------------------------------
let currentMode = 'txt2img';
let initImageData = null; // base64 string of uploaded image for img2img

function setWorkflowMode(mode) {
    currentMode = mode;
    document.querySelectorAll('.ig-workflow-btn').forEach(btn => {
        btn.classList.toggle('active', btn.dataset.mode === mode);
    });

    const isImg2Img = mode === 'img2img';
    const isInpaint = mode === 'inpaint';
    const isTxt2Img = mode === 'txt2img';
    const isOutpaint = mode === 'outpaint';
    const isUpscale = mode === 'upscale';
    const isVideo = mode === 'video';
    const needsImage = isImg2Img || isInpaint || isOutpaint || isUpscale;

    // Mode Tools panel — show correct section
    document.querySelectorAll('.ig-mode-section').forEach(el => el.style.display = 'none');
    const activeSection = document.getElementById('ig-mode-' + mode);
    if (activeSection) activeSection.style.display = '';

    // Canvas placeholder / upload zone
    const placeholder = document.getElementById('ig-canvas-placeholder');
    const uploadZone = document.getElementById('ig-upload-zone');
    const previewArea = document.getElementById('ig-preview-area');
    const maskCanvas = document.getElementById('ig-mask-canvas');

    if (placeholder) {
        if (needsImage) {
            placeholder.style.display = 'none';
            if (uploadZone) uploadZone.style.display = 'flex';
        } else {
            placeholder.style.display = '';
            if (uploadZone) uploadZone.style.display = 'none';
            initImageData = null;
            clearMask();
            // Reset mini upload zones
            document.querySelectorAll('.ig-upload-mini').forEach(z => {
                if (z.id === 'ig-upload-controlnet') return; // Don't reset ControlNet from mode change
                z.classList.remove('ig-uploaded');
                const span = z.querySelector('span');
                if (span) span.textContent = 'Drop image or click to browse';
            });
        }
    }

    if (maskCanvas) {
        maskCanvas.style.display = isInpaint ? 'block' : 'none';
        maskCanvas.classList.toggle('active', isInpaint);
        if (isInpaint) resizeMaskCanvas();
    }

    if ((isTxt2Img || isVideo) && currentResult?.image_data) {
        previewArea.innerHTML = `<img src="data:image/png;base64,${currentResult.image_data}" alt="Generated">`;
    } else if (!needsImage && !initImageData) {
        previewArea.innerHTML = '';
        if (placeholder) placeholder.style.display = '';
    }

    // Update labels for new mode sliders
    updateLabels();
}

// --- ControlNet Upload Handler ---
let controlNetImageData = null;
function handleControlNetUpload(file) {
    if (!file.type.startsWith('image/')) {
        setStatus('Please upload an image file for ControlNet.');
        return;
    }
    const reader = new FileReader();
    reader.onload = (e) => {
        const dataUrl = e.target.result;
        const commaIdx = dataUrl.indexOf(',');
        controlNetImageData = commaIdx >= 0 ? dataUrl.slice(commaIdx + 1) : dataUrl;

        const zone = document.getElementById('ig-upload-controlnet');
        if (zone) {
            zone.classList.add('ig-uploaded');
            zone.querySelector('span').textContent = 'Control image loaded';
        }
        setStatus('ControlNet image ready.');
    };
    reader.readAsDataURL(file);
}

// --- File Upload Handler ---
function handleImageUpload(file) {
    if (!file.type.startsWith('image/')) {
        setStatus('Please upload an image file (PNG, JPG, WEBP).');
        return;
    }
    const reader = new FileReader();
    reader.onload = (e) => {
        const dataUrl = e.target.result;
        // Extract base64 portion after comma
        const commaIdx = dataUrl.indexOf(',');
        initImageData = commaIdx >= 0 ? dataUrl.slice(commaIdx + 1) : dataUrl;

        const previewArea = document.getElementById('ig-preview-area');
        const uploadZone = document.getElementById('ig-upload-zone');
        previewArea.innerHTML = `<img id="ig-init-image" src="${dataUrl}" alt="Init Image" style="max-width:100%;max-height:100%;object-fit:contain;border-radius:8px;">`;
        if (uploadZone) uploadZone.style.display = 'none';

        // Mark mini upload zones as loaded
        document.querySelectorAll('.ig-upload-mini').forEach(z => {
            z.classList.add('ig-uploaded');
            z.querySelector('span').textContent = 'Image loaded — click to replace';
        });

        // Size mask canvas to match the new image
        requestAnimationFrame(() => {
            resizeMaskCanvas();
            clearMask();
        });

        setStatus('Image loaded. Adjust strength and click Generate.');
    };
    reader.readAsDataURL(file);
}

// --- Mask Drawing ---
let maskCtx = null;
let isDrawing = false;
let lastPos = null;

function resizeMaskCanvas() {
    const canvas = document.getElementById('ig-mask-canvas');
    const img = document.getElementById('ig-init-image');
    if (!canvas || !img) return;

    // Match canvas size to the displayed image size
    const rect = img.getBoundingClientRect();
    const wrapper = document.querySelector('.ig-canvas-wrapper');
    if (!wrapper) return;
    const wrapRect = wrapper.getBoundingClientRect();

    canvas.width = rect.width;
    canvas.height = rect.height;
    canvas.style.width = rect.width + 'px';
    canvas.style.height = rect.height + 'px';

    maskCtx = canvas.getContext('2d');
    clearMask();
}

function clearMask() {
    if (!maskCtx) return;
    maskCtx.fillStyle = '#000000';
    maskCtx.fillRect(0, 0, maskCtx.canvas.width, maskCtx.canvas.height);
}

function getMaskBase64() {
    if (!maskCtx) return '';
    return maskCtx.canvas.toDataURL('image/png').split(',')[1];
}

function initMaskDrawing() {
    const canvas = document.getElementById('ig-mask-canvas');
    if (!canvas) return;

    canvas.addEventListener('mousedown', (e) => {
        if (currentMode !== 'inpaint') return;
        isDrawing = true;
        lastPos = { x: e.offsetX, y: e.offsetY };
    });
    canvas.addEventListener('mousemove', (e) => {
        if (!isDrawing || currentMode !== 'inpaint') return;
        drawStroke(lastPos.x, lastPos.y, e.offsetX, e.offsetY);
        lastPos = { x: e.offsetX, y: e.offsetY };
    });
    window.addEventListener('mouseup', () => { isDrawing = false; lastPos = null; });

    // Touch support
    canvas.addEventListener('touchstart', (e) => {
        if (currentMode !== 'inpaint') return;
        e.preventDefault();
        const t = e.touches[0];
        const rect = canvas.getBoundingClientRect();
        isDrawing = true;
        lastPos = { x: t.clientX - rect.left, y: t.clientY - rect.top };
    }, { passive: false });
    canvas.addEventListener('touchmove', (e) => {
        if (!isDrawing || currentMode !== 'inpaint') return;
        e.preventDefault();
        const t = e.touches[0];
        const rect = canvas.getBoundingClientRect();
        const pos = { x: t.clientX - rect.left, y: t.clientY - rect.top };
        drawStroke(lastPos.x, lastPos.y, pos.x, pos.y);
        lastPos = pos;
    }, { passive: false });
    window.addEventListener('touchend', () => { isDrawing = false; lastPos = null; });
}

function drawStroke(x0, y0, x1, y1) {
    if (!maskCtx) return;
    const size = parseInt(document.getElementById('ig-brush-size')?.value || '16', 10);
    maskCtx.globalCompositeOperation = 'source-over';
    maskCtx.strokeStyle = '#ffffff';
    maskCtx.lineWidth = size;
    maskCtx.lineCap = 'round';
    maskCtx.lineJoin = 'round';
    maskCtx.beginPath();
    maskCtx.moveTo(x0, y0);
    maskCtx.lineTo(x1, y1);
    maskCtx.stroke();
}

// ---------------------------------------------------------------------------
// Model Control
// ---------------------------------------------------------------------------
function onUnloadModel() {
    setStatus('Unloading model...');
    if (window.go?.main?.App?.UnloadImageModel) {
        window.go.main.App.UnloadImageModel()
            .then(() => {
                setStatus('Model unloaded');
                const statusEl = document.getElementById('ig-model-status');
                if (statusEl) { statusEl.textContent = 'No model loaded'; statusEl.className = 'ig-panel-status'; }
                updateLiveModelInfo(null);
            })
            .catch(err => { setStatus('Unload failed: ' + err.message); });
    } else {
        setStatus('Unload not available in this backend.');
    }
}

function onReloadModel() {
    const name = document.getElementById('ig-model').value;
    if (!name) { setStatus('No model selected to reload.'); return; }
    setModel(name);
}

// ---------------------------------------------------------------------------
// Presets
// ---------------------------------------------------------------------------
const PRESETS = {
    default:       { prompt: '', negative: '', steps: 20, cfg: 7,  sampler: 'euler_a', scheduler: 'karras', width: 512, height: 512 },
    realistic:     { prompt: 'masterpiece, best quality, highly detailed, realistic, 8k uhd, professional photography', negative: 'cartoon, anime, 3d render, blurry, low quality, deformed', steps: 25, cfg: 7, sampler: 'dpm++2m', scheduler: 'karras', width: 512, height: 768 },
    cinematic:     { prompt: 'cinematic shot, dramatic lighting, film grain, color graded, movie still', negative: 'overexposed, underexposed, video game, cartoon', steps: 30, cfg: 7, sampler: 'dpm++2m', scheduler: 'karras', width: 1024, height: 576 },
    anime:         { prompt: 'masterpiece, best quality, anime style, detailed background, vibrant colors', negative: 'realistic, 3d, western cartoon, photo, blurry', steps: 25, cfg: 7, sampler: 'euler_a', scheduler: 'karras', width: 512, height: 768 },
    illustration:  { prompt: 'detailed illustration, artistic, stylized, clean linework, concept art', negative: 'photograph, realistic, 3d render, noisy', steps: 25, cfg: 7, sampler: 'euler_a', scheduler: 'normal', width: 768, height: 768 },
    'low-light':   { prompt: 'night scene, moody lighting, dark atmosphere, cinematic shadows', negative: 'bright, overexposed, daylight, flat lighting', steps: 30, cfg: 8, sampler: 'dpm++2m', scheduler: 'karras', width: 768, height: 512 },
    'high-detail': { prompt: 'hyperdetailed, intricate, 8k resolution, sharp focus, ultra detailed', negative: 'simple, low detail, blurry, soft focus', steps: 40, cfg: 7, sampler: 'dpm++2m', scheduler: 'exponential', width: 1024, height: 1024 },
};

function applyPreset() {
    const key = document.getElementById('ig-preset').value;
    if (!key || !PRESETS[key]) { setStatus('Select a preset first.'); return; }
    const p = PRESETS[key];
    if (p.prompt) document.getElementById('ig-prompt').value = p.prompt;
    if (p.negative) document.getElementById('ig-negative').value = p.negative;
    document.getElementById('ig-steps').value = p.steps;
    document.getElementById('ig-cfg').value = p.cfg;
    document.getElementById('ig-sampler').value = p.sampler;
    document.getElementById('ig-scheduler').value = p.scheduler;
    document.getElementById('ig-width').value = p.width;
    document.getElementById('ig-height').value = p.height;
    updateLabels();
    saveState();
    setStatus('Preset applied: ' + key);
}

function saveCurrentPreset() {
    const name = prompt('Enter preset name (no spaces):');
    if (!name) return;
    const key = name.toLowerCase().replace(/\s+/g, '-');
    PRESETS[key] = {
        prompt: document.getElementById('ig-prompt').value,
        negative: document.getElementById('ig-negative').value,
        steps: parseInt(document.getElementById('ig-steps').value, 10),
        cfg: parseFloat(document.getElementById('ig-cfg').value),
        sampler: document.getElementById('ig-sampler').value,
        scheduler: document.getElementById('ig-scheduler').value,
        width: parseInt(document.getElementById('ig-width').value, 10),
        height: parseInt(document.getElementById('ig-height').value, 10),
    };
    const sel = document.getElementById('ig-preset');
    const opt = document.createElement('option');
    opt.value = key;
    opt.textContent = name;
    sel.appendChild(opt);
    sel.value = key;
    setStatus('Preset saved: ' + name);
}

// ---------------------------------------------------------------------------
// Seed Tools
// ---------------------------------------------------------------------------
function randomizeSeed() {
    const seed = Math.floor(Math.random() * 2147483647);
    document.getElementById('ig-seed').value = seed;
    saveState();
}

let seedLocked = false;
function toggleSeedLock() {
    seedLocked = !seedLocked;
    const btn = document.getElementById('ig-seed-lock');
    if (btn) btn.classList.toggle('ig-locked', seedLocked);
}

function reuseLastSeed() {
    if (lastUsedSeed !== -1) {
        document.getElementById('ig-seed').value = lastUsedSeed;
        saveState();
    } else {
        setStatus('No previous seed to reuse.');
    }
}

// ---------------------------------------------------------------------------
// Prompt Weighting
// ---------------------------------------------------------------------------
function wrapPromptWeight() {
    const ta = document.getElementById('ig-prompt');
    const start = ta.selectionStart;
    const end = ta.selectionEnd;
    const text = ta.value;
    const before = text.slice(0, start);
    const selected = text.slice(start, end);
    const after = text.slice(end);
    const newText = before + '(' + (selected || 'word') + ')' + after;
    ta.value = newText;
    ta.selectionStart = start + 1;
    ta.selectionEnd = start + 1 + (selected ? selected.length : 4);
    ta.focus();
    saveState();
}

function adjustWeight(delta) {
    const ta = document.getElementById('ig-prompt');
    const text = ta.value;
    const start = ta.selectionStart;
    const end = ta.selectionEnd;
    const selected = text.slice(start, end);
    if (!selected) { setStatus('Select text to weight first.'); return; }
    let match = selected.match(/^\(([\s\S]+)\)$/);
    if (match) {
        const inner = match[1];
        const numMatch = inner.match(/^(.*):(\d+(\.\d+)?)$/);
        if (numMatch) {
            const newVal = Math.max(0, Math.min(5, parseFloat(numMatch[2]) + delta * 0.1));
            const newText = text.slice(0, start) + '(' + numMatch[1] + ':' + newVal.toFixed(1) + ')' + text.slice(end);
            ta.value = newText;
        } else {
            const newVal = (1 + delta * 0.1).toFixed(1);
            const newText = text.slice(0, start) + '(' + inner + ':' + newVal + ')' + text.slice(end);
            ta.value = newText;
        }
    } else {
        const newVal = (1 + delta * 0.1).toFixed(1);
        const newText = text.slice(0, start) + '(' + selected + ':' + newVal + ')' + text.slice(end);
        ta.value = newText;
    }
    ta.focus();
    saveState();
}

// ---------------------------------------------------------------------------
// Keyboard Shortcuts
// ---------------------------------------------------------------------------
function handleKeydown(e) {
    // Ctrl+Enter = Generate
    if (e.ctrlKey && e.key === 'Enter') {
        e.preventDefault();
        onGenerate();
        return;
    }
    // Ctrl+S = Save Image As
    if (e.ctrlKey && e.key === 's' && !e.shiftKey) {
        e.preventDefault();
        saveImageAs();
        return;
    }
    // Ctrl+L = Load Model
    if (e.ctrlKey && e.key === 'l') {
        e.preventDefault();
        onLoadModelClick();
        return;
    }
    // Ctrl+U = Unload Model
    if (e.ctrlKey && e.key === 'u') {
        e.preventDefault();
        onUnloadModel();
        return;
    }
    // Ctrl+R = Reload Model
    if (e.ctrlKey && e.key === 'r') {
        e.preventDefault();
        onReloadModel();
        return;
    }
    // R = Random Seed (when not in an input)
    if (e.key === 'r' && !e.ctrlKey && !e.altKey && !e.metaKey) {
        const active = document.activeElement;
        if (active && (active.tagName === 'INPUT' || active.tagName === 'TEXTAREA')) return;
        e.preventDefault();
        randomizeSeed();
    }
    // Workflow mode shortcuts
    if (e.ctrlKey) {
        const modeMap = { '1': 'txt2img', '2': 'img2img', '3': 'inpaint', '4': 'outpaint', '5': 'upscale' };
        if (modeMap[e.key]) {
            e.preventDefault();
            setWorkflowMode(modeMap[e.key]);
        }
    }
}
