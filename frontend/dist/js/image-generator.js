'use strict';

// ---------------------------------------------------------------------------
// Main Orchestrator
// ---------------------------------------------------------------------------
document.addEventListener('DOMContentLoaded', async () => {
    showRuntimeDebug();
    restoreState();
    await initPaths();
    initPathListeners();
    bindEvents();
    loadHistoryPanel();
    updateDiagnostics('init');
    startDiagnosticsPoll();
});

async function startDiagnosticsPoll() {
    if (!window.go?.main?.App?.GetAppDiagnostics) return;
    const poll = async () => {
        try {
            const d = await window.go.main.App.GetAppDiagnostics();
            document.getElementById('diag-engine').textContent = d.model_loaded ? d.model_name : 'Idle';
            document.getElementById('diag-heap').textContent = d.go_heap_alloc_mb + ' MB';
            document.getElementById('diag-sys').textContent = d.go_heap_sys_mb + ' MB';
            document.getElementById('diag-goroutines').textContent = d.go_goroutines;
            document.getElementById('diag-gc').textContent = d.go_num_gc;
        } catch (e) { /* ignore */ }
    };
    poll();
    setInterval(poll, 2000);
}

// --- History Panel ---
async function loadHistoryPanel() {
    if (!window.go?.main?.App?.ListGalleryImages) return;
    try {
        const images = await window.go.main.App.ListGalleryImages();
        const list = document.getElementById('ig-history-list');
        if (!images || images.length === 0) {
            list.innerHTML = '<div class="ig-history-placeholder">Generated images appear here</div>';
            return;
        }
        list.innerHTML = '';
        images.reverse().forEach(img => {
            const item = document.createElement('div');
            item.className = 'ig-history-item';
            item.dataset.file = img.file;
            item.innerHTML = `
                <button class="ig-history-delete" title="Delete"
                    onclick="event.stopPropagation();deleteHistoryImage('${escapeJS(img.file)}',this)">&#x1F5D1;</button>
                <img id="hist-img-${escapeId(img.file)}" alt="" loading="lazy"
                     style="width:100%;aspect-ratio:1;object-fit:cover;border-radius:6px;">
                <div class="ig-history-actions">
                    <button class="ig-history-send-btn" title="Send to Character"
                        onclick="event.stopPropagation();sendHistoryToEntity('${escapeJS(img.file)}','character')">Char</button>
                    <button class="ig-history-send-btn" title="Send to World"
                        onclick="event.stopPropagation();sendHistoryToEntity('${escapeJS(img.file)}','world')">World</button>
                    <button class="ig-history-send-btn" title="Send to Scenario"
                        onclick="event.stopPropagation();sendHistoryToEntity('${escapeJS(img.file)}','scenario')">Scen</button>
                </div>
            `;
            item.addEventListener('click', () => viewHistoryImage(img.file));
            list.appendChild(item);
            loadGalleryImage(img.file, `hist-img-${escapeId(img.file)}`);
        });
    } catch (err) {
        console.error('Failed to load history:', err);
    }
}

async function viewHistoryImage(filePath) {
    if (!window.go?.main?.App?.GetGalleryImage) return;
    try {
        setStatus('Loading image...');
        const b64 = await window.go.main.App.GetGalleryImage(filePath);
        const area = document.getElementById('ig-preview-area');
        area.innerHTML = `<img src="data:image/png;base64,${b64}" alt="History">`;
        document.querySelectorAll('.ig-history-item').forEach(t => t.classList.remove('active'));
        const active = document.querySelector(`.ig-history-item[data-file="${CSS.escape(filePath)}"]`);
        if (active) active.classList.add('active');
        document.getElementById('ig-meta').innerHTML = `<strong>File:</strong> ${filePath.replace(/^.*[\\/]/, '')}`;
        document.getElementById('ig-meta').style.display = 'block';
        document.getElementById('ig-actions').style.display = 'flex';
        // Set currentResult so Save/Delete work for history images
        currentResult = { image_data: b64, metadata: { file: filePath } };
        setStatus('');
    } catch (err) {
        console.error('View history failed:', err);
        setStatus('Failed to load image.');
    }
}

async function deleteHistoryImage(filePath, btnElement) {
    const baseName = filePath.replace(/^.*[\\/]/, '').replace(/\.[^.]+$/, '');
    try {
        await window.go.main.App.DeleteGalleryImage(baseName);
        const item = btnElement.closest('.ig-history-item');
        if (item) item.remove();
        const list = document.getElementById('ig-history-list');
        if (list.children.length === 0) {
            list.innerHTML = '<div class="ig-history-placeholder">Generated images appear here</div>';
        }
    } catch (err) {
        console.error('Delete failed:', err);
    }
}

async function sendHistoryToEntity(filePath, entityType) {
    if (!window.go?.main?.App?.SaveEntityImage) {
        setStatus('Entity image save not available.');
        return;
    }
    const entityName = prompt(`Enter ${entityType} name (or ID):`)?.trim();
    if (!entityName) { setStatus('Cancelled.'); return; }
    const safeID = entityName.toLowerCase().replace(/[^a-z0-9_-]/g, '_').substring(0, 40);
    try {
        setStatus(`Saving image for ${entityType}...`);
        const b64 = await window.go.main.App.GetGalleryImage(filePath);
        const path = await window.go.main.App.SaveEntityImage(b64, entityType, safeID);
        setStatus(`Image saved for ${entityType}. Navigate to ${entityType} page to assign it.`);
        if (confirm(`Image saved. Open ${entityType} page now?`)) {
            const pageMap = { character: 'character.html', world: 'world.html', scenario: 'scenario.html' };
            window.location.href = `${pageMap[entityType] || 'index.html'}?pendingImage=${encodeURIComponent(path)}`;
        }
    } catch (err) {
        console.error('Send to entity failed:', err);
        setStatus('Failed: ' + (err.message || err));
    }
}

function addHistoryItem(filePath, imageData) {
    const list = document.getElementById('ig-history-list');
    if (list.querySelector('.ig-history-placeholder')) list.innerHTML = '';
    document.querySelectorAll('.ig-history-item').forEach(t => t.classList.remove('active'));
    const item = document.createElement('div');
    item.className = 'ig-history-item active';
    item.dataset.file = filePath;
    const safeId = escapeId(filePath);
    item.innerHTML = `
        <button class="ig-history-delete" title="Delete"
            onclick="event.stopPropagation();deleteHistoryImage('${escapeJS(filePath)}',this)">&#x1F5D1;</button>
        <img id="hist-img-${safeId}" src="data:image/png;base64,${imageData}" alt=""
             style="width:100%;aspect-ratio:1;object-fit:cover;border-radius:6px;">
        <div class="ig-history-actions">
            <button class="ig-history-send-btn" title="Send to Character"
                onclick="event.stopPropagation();sendHistoryToEntity('${escapeJS(filePath)}','character')">Char</button>
            <button class="ig-history-send-btn" title="Send to World"
                onclick="event.stopPropagation();sendHistoryToEntity('${escapeJS(filePath)}','world')">World</button>
            <button class="ig-history-send-btn" title="Send to Scenario"
                onclick="event.stopPropagation();sendHistoryToEntity('${escapeJS(filePath)}','scenario')">Scen</button>
        </div>
    `;
    item.addEventListener('click', () => viewHistoryImage(filePath));
    list.insertBefore(item, list.firstChild);
}

// --- Model Info ---
function updateLiveModelInfo(model) {
    const statusEl = document.getElementById('diag-status');
    const typeEl = document.getElementById('diag-type');
    const familyEl = document.getElementById('diag-family');
    const formatEl = document.getElementById('diag-format');
    if (!model) {
        if (statusEl) statusEl.textContent = 'Idle';
        if (typeEl) typeEl.textContent = '-';
        if (familyEl) familyEl.textContent = '-';
        if (formatEl) formatEl.textContent = '-';
        return;
    }
    let typeStr = 'Unknown';
    if (model.family === 'sd15') typeStr = 'SD 1.5';
    else if (model.family === 'sdxl') typeStr = 'SDXL';
    else if (model.family === 'sd3') typeStr = 'SD3';
    else if (model.family === 'flux') typeStr = 'Flux';
    if (statusEl) statusEl.textContent = 'Loaded';
    if (typeEl) typeEl.textContent = typeStr;
    if (familyEl) familyEl.textContent = model.family || '-';
    if (formatEl) formatEl.textContent = model.format || '-';
}

// --- Diagnostics ---
let diagLoadStart = 0;
let diagGenStart = 0;
let diagErrorCount = 0;

function updateDiagnostics(event, detail) {
    const loadTimeEl = document.getElementById('diag-loadtime');
    const lastGenEl = document.getElementById('diag-lastgen');
    const errorsEl = document.getElementById('diag-errors');

    switch (event) {
        case 'init':
            if (loadTimeEl) loadTimeEl.textContent = '-';
            if (lastGenEl) lastGenEl.textContent = '-';
            if (errorsEl) { errorsEl.textContent = 'None'; errorsEl.className = 'ig-info-value ig-info-ok'; }
            break;
        case 'loadStart':
            diagLoadStart = performance.now();
            break;
        case 'loadDone':
            if (loadTimeEl && diagLoadStart) {
                const ms = Math.round(performance.now() - diagLoadStart);
                loadTimeEl.textContent = (ms / 1000).toFixed(1) + 's';
            }
            break;
        case 'genStart':
            diagGenStart = performance.now();
            break;
        case 'genDone':
            if (lastGenEl && diagGenStart) {
                const ms = Math.round(performance.now() - diagGenStart);
                lastGenEl.textContent = (ms / 1000).toFixed(1) + 's';
            }
            break;
        case 'error':
            diagErrorCount++;
            if (errorsEl) {
                errorsEl.textContent = diagErrorCount + (diagErrorCount === 1 ? ' error' : ' errors');
                errorsEl.className = 'ig-info-value ig-info-error';
            }
            break;
    }
}


document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') closeTutorials();
});

// --- Metadata Viewer ---
function toggleMetadata() {
    const viewer = document.getElementById('ig-metadata-viewer');
    if (!viewer) return;
    if (viewer.style.display === 'none') {
        if (!currentResult?.metadata) {
            setStatus('No metadata available.');
            return;
        }
        viewer.innerHTML = '<pre>' + JSON.stringify(currentResult.metadata, null, 2) + '</pre>';
        viewer.style.display = 'block';
    } else {
        viewer.style.display = 'none';
    }
}
