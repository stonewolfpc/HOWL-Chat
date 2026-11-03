'use strict';

// ---------------------------------------------------------------------------
// Board System
// ---------------------------------------------------------------------------
const BOARD_STORAGE_KEY = 'howl-image-boards';
const IMAGE_BOARDS_KEY = 'howl-image-board-assignments';
let boards = [];
let imageBoardAssignments = {};

function loadBoards() {
    try {
        const stored = localStorage.getItem(BOARD_STORAGE_KEY);
        if (stored) {
            boards = JSON.parse(stored);
        } else {
            boards = [{ id: 'uncategorized', name: 'Uncategorized' }];
        }
    } catch (e) {
        console.error('Failed to load boards:', e);
        boards = [{ id: 'uncategorized', name: 'Uncategorized' }];
    }
    updateBoardSelect();
}

function loadImageBoardAssignments() {
    try {
        const stored = localStorage.getItem(IMAGE_BOARDS_KEY);
        if (stored) {
            imageBoardAssignments = JSON.parse(stored);
        }
    } catch (e) {
        console.error('Failed to load image board assignments:', e);
        imageBoardAssignments = {};
    }
}

function saveBoards() {
    try {
        localStorage.setItem(BOARD_STORAGE_KEY, JSON.stringify(boards));
    } catch (e) {
        console.error('Failed to save boards:', e);
    }
}

function saveImageBoardAssignments() {
    try {
        localStorage.setItem(IMAGE_BOARDS_KEY, JSON.stringify(imageBoardAssignments));
    } catch (e) {
        console.error('Failed to save image board assignments:', e);
    }
}

function updateBoardSelect() {
    const select = document.getElementById('ig-board-select');
    if (!select) return;
    const currentValue = select.value;
    select.innerHTML = '';
    boards.forEach(board => {
        const option = document.createElement('option');
        option.value = board.id;
        option.textContent = board.name;
        select.appendChild(option);
    });
    if (currentValue && boards.find(b => b.id === currentValue)) {
        select.value = currentValue;
    }
}

function addBoard() {
    const name = prompt('Enter board name:');
    if (!name || name.trim() === '') return;
    const id = name.toLowerCase().replace(/\s+/g, '-').replace(/[^a-z0-9-]/g, '');
    if (boards.find(b => b.id === id)) {
        alert('A board with this name already exists');
        return;
    }
    boards.push({ id, name });
    saveBoards();
    updateBoardSelect();
    // Select the new board
    const select = document.getElementById('ig-board-select');
    if (select) select.value = id;
}

function assignImageToBoard(imageFile, boardId) {
    imageBoardAssignments[imageFile] = boardId;
    saveImageBoardAssignments();
    loadHistoryPanel(); // Refresh the display
}

function showBoardContextMenu(imageFile, x, y) {
    // Remove any existing context menu
    const existing = document.getElementById('ig-board-context-menu');
    if (existing) existing.remove();

    const menu = document.createElement('div');
    menu.id = 'ig-board-context-menu';
    menu.className = 'ig-board-context-menu';
    menu.style.left = x + 'px';
    menu.style.top = y + 'px';

    let menuHtml = '<div class="ig-board-context-title">Send to Board</div>';
    
    // Add existing boards
    boards.forEach(board => {
        menuHtml += `<button class="ig-board-context-item" onclick="assignImageToBoard('${escapeJS(imageFile)}', '${board.id}')">${board.name}</button>`;
    });
    
    // Add "Create New" option
    menuHtml += `<button class="ig-board-context-item ig-board-context-create" onclick="createBoardAndAssign('${escapeJS(imageFile)}')">+ Create New Board</button>`;
    
    menu.innerHTML = menuHtml;
    document.body.appendChild(menu);

    // Close menu when clicking outside
    setTimeout(() => {
        document.addEventListener('click', hideBoardContextMenu, { once: true });
    }, 100);
}

function hideBoardContextMenu() {
    const menu = document.getElementById('ig-board-context-menu');
    if (menu) menu.remove();
}

function createBoardAndAssign(imageFile) {
    hideBoardContextMenu();
    const name = prompt('Enter new board name:');
    if (!name || name.trim() === '') return;
    const id = name.toLowerCase().replace(/\s+/g, '-').replace(/[^a-z0-9-]/g, '');
    if (boards.find(b => b.id === id)) {
        alert('A board with this name already exists');
        return;
    }
    boards.push({ id, name });
    saveBoards();
    updateBoardSelect();
    assignImageToBoard(imageFile, id);
}

function initBoardSystem() {
    loadBoards();
    loadImageBoardAssignments();
    const addBtn = document.getElementById('ig-add-board');
    if (addBtn) {
        addBtn.addEventListener('click', addBoard);
    }
    
    // Add board select change listener to filter images
    const boardSelect = document.getElementById('ig-board-select');
    if (boardSelect) {
        boardSelect.addEventListener('change', () => loadHistoryPanel());
    }
}

// ---------------------------------------------------------------------------
// Main Orchestrator
// ---------------------------------------------------------------------------
document.addEventListener('DOMContentLoaded', async () => {
    showRuntimeDebug();
    restoreState();
    await initPaths();
    initPathListeners();
    bindEvents();
    initBoardSystem();
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
        try {
            const ver = await window.go.main.App.GetSDServerVersion();
            const el = document.getElementById('diag-sdver');
            if (el) el.textContent = ver;
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
        const boardSelect = document.getElementById('ig-board-select');
        const selectedBoard = boardSelect ? boardSelect.value : 'uncategorized';
        
        if (!images || images.length === 0) {
            list.innerHTML = '<div class="ig-history-placeholder">Generated images appear here</div>';
            return;
        }
        
        // Filter images by selected board
        const filteredImages = images.filter(img => {
            const assignedBoard = imageBoardAssignments[img.file] || 'uncategorized';
            return assignedBoard === selectedBoard;
        });
        
        if (filteredImages.length === 0) {
            list.innerHTML = `<div class="ig-history-placeholder">No images in this board</div>`;
            return;
        }
        
        list.innerHTML = '';
        filteredImages.reverse().forEach(img => {
            const item = document.createElement('div');
            item.className = 'ig-history-item';
            item.dataset.file = img.file;
            item.dataset.meta = JSON.stringify(img);
            item.innerHTML = `
                <button class="ig-history-delete" title="Delete"
                    onclick="event.stopPropagation();deleteHistoryImage('${escapeJS(img.file)}',this)">&#x1F5D1;</button>
                <button class="ig-history-menu-btn" title="Options"
                    onclick="event.stopPropagation();toggleHistoryMenu(this)">&#x22EE;</button>
                <div class="ig-history-menu-popup" style="display:none;">
                    <button onclick="event.stopPropagation();historyMenuAction('${escapeJS(img.file)}','import-seed',this)">Import Seed</button>
                    <button onclick="event.stopPropagation();historyMenuAction('${escapeJS(img.file)}','use-prompt',this)">Use Prompt</button>
                    <button onclick="event.stopPropagation();historyMenuAction('${escapeJS(img.file)}','repeat-new-model',this)">Repeat Generation (New Model)</button>
                    <button onclick="event.stopPropagation();historyMenuAction('${escapeJS(img.file)}','repeat-same-model',this)">Repeat with Same Model</button>
                    <hr style="margin:4px 0;border:1px solid #444;">
                    <button onclick="event.stopPropagation();showBoardContextMenu('${escapeJS(img.file)}', event.clientX, event.clientY)" style="color:#4CAF50;">Send to Board</button>
                </div>
                <img id="hist-img-${escapeId(img.file)}" alt="" loading="lazy"
                     class="ig-history-thumb">
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
            // Right-click is handled once, by image-recall.js, which offers the full
            // parameter set plus a board action. A handler here would stack a second
            // menu underneath it.
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

function toggleHistoryMenu(btn) {
    const popup = btn.nextElementSibling;
    if (!popup) return;
    const isVisible = popup.style.display === 'block';
    document.querySelectorAll('.ig-history-menu-popup').forEach(p => p.style.display = 'none');
    popup.style.display = isVisible ? 'none' : 'block';
}

function historyMenuAction(filePath, action, btn) {
    const popup = btn.closest('.ig-history-menu-popup');
    if (popup) popup.style.display = 'none';

    const item = document.querySelector(`.ig-history-item[data-file="${CSS.escape(filePath)}"]`);
    if (!item) return;
    const metaRaw = item.dataset.meta;
    if (!metaRaw) return;
    let meta;
    try { meta = JSON.parse(metaRaw); } catch (e) { return; }

    const opts = meta.Options || meta.options || {};
    const hasFullOpts = Object.keys(opts).length > 0;

    switch (action) {
        case 'import-seed': {
            const seed = hasFullOpts ? (opts.Seed ?? meta.Seed) : meta.Seed;
            if (seed !== undefined) {
                const el = document.getElementById('ig-seed');
                if (el) el.value = seed;
            }
            break;
        }
        case 'use-prompt': {
            const prompt = hasFullOpts ? (opts.Prompt ?? meta.Prompt) : meta.Prompt;
            const neg = hasFullOpts ? (opts.NegativePrompt ?? meta.NegativePrompt) : meta.NegativePrompt;
            if (prompt) {
                const el = document.getElementById('ig-prompt');
                if (el) el.value = prompt;
            }
            if (neg) {
                const el = document.getElementById('ig-negative');
                if (el) el.value = neg;
            }
            break;
        }
        case 'repeat-new-model': {
            if (hasFullOpts) {
                populateOptionsFromHistory(opts, true);
            } else {
                populateOptionsFromLegacy(meta, true);
            }
            break;
        }
        case 'repeat-same-model': {
            if (hasFullOpts) {
                populateOptionsFromHistory(opts, false);
            } else {
                populateOptionsFromLegacy(meta, false);
            }
            triggerGenerate();
            break;
        }
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
        <img id="hist-img-${safeId}" src="data:image/png;base64,${imageData}" alt="">
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


// --- Modal System ---
function openModal(id) {
    const modal = document.getElementById(id);
    if (modal) modal.style.display = 'flex';
}

function closeModal(id) {
    const modal = document.getElementById(id);
    if (modal) modal.style.display = 'none';
}

// --- Sidebar Toggle ---
function toggleSidebar() {
    const sidebar = document.getElementById('ig-sidebar');
    if (!sidebar) return;
    sidebar.classList.toggle('collapsed');
    const btn = document.getElementById('ig-sidebar-toggle');
    if (btn) {
        btn.innerHTML = sidebar.classList.contains('collapsed') ? '&#x25B6;' : '&#x25C0;';
    }
}

document.addEventListener('DOMContentLoaded', () => {
    const toggleBtn = document.getElementById('ig-sidebar-toggle');
    if (toggleBtn) toggleBtn.addEventListener('click', toggleSidebar);
});

document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
        closeTutorials();
        // Close any open modals
        document.querySelectorAll('.ig-modal').forEach(m => { m.style.display = 'none'; });
    }
});

// Close modal when clicking outside content
document.addEventListener('click', (e) => {
    if (e.target.classList.contains('ig-modal')) {
        e.target.style.display = 'none';
    }
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

// Board search: filters the history grid by filename or prompt text.
// Called from the input's oninput attribute so no load-order or listener timing matters.
function filterHistoryBySearch(value) {
    const list = document.getElementById('ig-history-list');
    if (!list) return;
    const term = (value || '').trim().toLowerCase();
    let shown = 0;
    list.querySelectorAll('.ig-history-item').forEach(item => {
        if (!term) { item.style.display = ''; shown++; return; }
        let haystack = (item.dataset.file || '').toLowerCase();
        try {
            haystack += ' ' + JSON.stringify(JSON.parse(item.dataset.meta || '{}')).toLowerCase();
        } catch (e) { /* meta is optional */ }
        const match = haystack.includes(term);
        item.style.display = match ? '' : 'none';
        if (match) shown++;
    });
    const status = document.getElementById('ig-history-status');
    if (status) {
        if (term) status.textContent = shown + ' match' + (shown === 1 ? '' : 'es') + ' for "' + value + '"';
        else status.textContent = '';
    }
}
