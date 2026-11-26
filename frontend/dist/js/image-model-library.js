/**
 * image-model-library.js
 *
 * Model Library panel for the Image Generator page.
 *
 * Workflow (mirrors InvokeAI):
 *   1. User opens the Model Library panel.
 *   2. Clicks "Scan Folder" — backend walks the models directory, classifies
 *      every file, and returns a categorised index.
 *   3. The results are shown in a read-only preview table (type, family, size).
 *   4. User clicks "Install All" or selects individual rows and clicks "Install
 *      Selected". The backend physically moves/copies files into organised
 *      subdirectories and re-scans.
 *   5. On completion the main model/VAE/LoRA/ControlNet dropdowns are
 *      auto-populated from the fresh index.
 *
 * Depends on: window.go bindings — ScanImageModels, InstallAssets, GetCompatibleAssets
 */

'use strict';

// ── state ────────────────────────────────────────────────────────────────────

const ModelLibrary = (() => {
  let _scanResult   = null;   // last ScanResult from backend
  let _selectedPaths = new Set(); // paths ticked in the preview table

  // ── public API ─────────────────────────────────────────────────────────────

  function init() {
    _bindButtons();
    _loadFromRegistry();
  }

  /** Load the persistent registry on startup — no scan needed. */
  async function _loadFromRegistry() {
    try {
      const result = await window.go.main.App.GetRegisteredModels();
      const total = _totalCount(result);
      if (total > 0) {
        _scanResult = result;
        _renderTable(result);
        _setScanStatus(`${total} registered model${total !== 1 ? 's' : ''} — Scan to discover more.`, 'done');
        _setInstallEnabled(true);
      }
    } catch (_) {
      // Registry empty or not yet initialised — silent, user can scan manually.
    }
  }

  /** Force-rescan: clears stale .howl.json manifests then does a fresh scan. */
  async function rescan() {
    _setScanStatus('Clearing cached manifests…', 'scanning');
    _clearTable();
    try {
      _scanResult = await window.go.main.App.RescanImageModels();
      _renderTable(_scanResult);
      const total = _totalCount(_scanResult);
      _setScanStatus(`Rescan complete — ${total} asset${total !== 1 ? 's' : ''} found`, 'done');
      _setInstallEnabled(total > 0);
    } catch (err) {
      _setScanStatus(`Rescan failed: ${err}`, 'error');
    }
  }

  /** Trigger a backend scan and render the preview table. */
  async function scan() {
    _setScanStatus('Scanning…', 'scanning');
    _clearTable();
    try {
      _scanResult = await window.go.main.App.ScanImageModels();
      _renderTable(_scanResult);
      const total = _totalCount(_scanResult);
      _setScanStatus(`Found ${total} asset${total !== 1 ? 's' : ''}`, 'done');
      _setInstallEnabled(total > 0);
    } catch (err) {
      _setScanStatus(`Scan failed: ${err}`, 'error');
    }
  }

  /** Install all scanned assets (move mode = true). */
  async function installAll() {
    await _install(true);
  }

  /** Install only the user-selected rows. */
  async function installSelected() {
    if (_selectedPaths.size === 0) {
      _setScanStatus('No assets selected.', 'warn');
      return;
    }
    await _install(true);
  }

  // ── private ────────────────────────────────────────────────────────────────

  function _bindButtons() {
    const btnScan   = document.getElementById('ml-btn-scan');
    const btnInstAll = document.getElementById('ml-btn-register-all');
    const btnInstSel = document.getElementById('ml-btn-register-selected');
    const btnSelAll  = document.getElementById('ml-btn-select-all');
    const btnSelNone = document.getElementById('ml-btn-select-none');

    const btnRescan = document.getElementById('ml-btn-rescan');

    if (btnScan)    btnScan.addEventListener('click', scan);
    if (btnRescan)  btnRescan.addEventListener('click', rescan);
    if (btnInstAll) btnInstAll.addEventListener('click', installAll);
    if (btnInstSel) btnInstSel.addEventListener('click', installSelected);
    if (btnSelAll)  btnSelAll.addEventListener('click', _selectAll);
    if (btnSelNone) btnSelNone.addEventListener('click', _selectNone);
  }

  async function _install(moveFiles) {
    _setScanStatus('Registering…', 'scanning');
    _setInstallEnabled(false);
    try {
      const summary = await window.go.main.App.InstallAssets(moveFiles);
      const n = summary.registered_count    || 0;
      const s = summary.already_known_count || 0;
      const f = summary.failed_count        || 0;
      _setScanStatus(
        `Registered ${n} model${n !== 1 ? 's' : ''}` +
        (s ? `, ${s} already known` : '') +
        (f ? `, ${f} failed` : '') + '. Files were not moved.',
        f > 0 ? 'error' : 'done'
      );
      _renderInstallSummary(summary);
      if (n > 0 && typeof window.refreshModelDropdownsFromScan === 'function') {
        window.refreshModelDropdownsFromScan();
      }
    } catch (err) {
      _setScanStatus(`Install failed: ${err}`, 'error');
      _setInstallEnabled(true);
    }
  }

  async function _refreshDropdownsAfterInstall() {
    try {
      // Full re-scan to pick up moved files
      _scanResult = await window.go.main.App.ScanImageModels();
      _renderTable(_scanResult);
      // Delegate to the core module which owns the dropdowns
      if (typeof window.refreshModelDropdownsFromScan === 'function') {
        window.refreshModelDropdownsFromScan(_scanResult);
      }
    } catch (_) { /* non-fatal */ }
  }

  // ── table rendering ────────────────────────────────────────────────────────

  function _clearTable() {
    const tbody = document.getElementById('ml-table-body');
    if (tbody) tbody.innerHTML = '';
    _selectedPaths.clear();
    _updateSelectionLabel();
  }

  function _renderTable(result) {
    const tbody = document.getElementById('ml-table-body');
    if (!tbody) return;
    tbody.innerHTML = '';
    _selectedPaths.clear();

    const groups = [
      { label: 'Diffusion Models', key: 'models',      type: 'model'      },
      { label: 'VAE',              key: 'vaes',         type: 'vae'        },
      { label: 'LoRA / LyCORIS',   key: 'loras',        type: 'lora'       },
      { label: 'ControlNet',       key: 'controlnets',  type: 'controlnet' },
      { label: 'IP-Adapter',       key: 'ip_adapters',  type: 'ip_adapter' },
      { label: 'CLIP-Vision',      key: 'clip_visions', type: 'clip_vision' },
      { label: 'Embeddings',       key: 'embeddings',   type: 'embedding'  },
      { label: 'Upscalers',        key: 'upscalers',    type: 'upscaler'   },
    ];

    let anyRows = false;
    for (const g of groups) {
      const items = result[g.key] || [];
      if (items.length === 0) continue;
      anyRows = true;

      // Group header row
      const hdr = document.createElement('tr');
      hdr.className = 'ml-group-header';
      hdr.innerHTML = `<td colspan="6">${g.label} <span class="ml-count">${items.length}</span></td>`;
      tbody.appendChild(hdr);

      for (const asset of items) {
        const row = _buildRow(asset);
        tbody.appendChild(row);
        // Pre-select all rows
        _selectedPaths.add(asset.path);
      }
    }

    if (!anyRows) {
      const empty = document.createElement('tr');
      empty.innerHTML = '<td colspan="6" class="ml-empty">No assets found. Check your models directory in Settings.</td>';
      tbody.appendChild(empty);
    }

    _updateSelectionLabel();
  }

  function _buildRow(asset) {
    const row = document.createElement('tr');
    row.className = 'ml-asset-row';
    row.dataset.path = asset.path;

    const checked = _selectedPaths.has(asset.path) ? 'checked' : '';
    const familyBadge = asset.family
      ? `<span class="ml-badge ml-badge-${asset.family}">${asset.family}</span>`
      : '';
    const sizeFmt = asset.size_mb >= 1000
      ? `${(asset.size_mb / 1024).toFixed(1)} GB`
      : `${Math.round(asset.size_mb)} MB`;

    const compatWarning = asset.compatibility === 'unsupported'
      ? ' <span class="ml-compat-warn" title="Not supported by current sd-server build — try a .safetensors version or update sd-server">&#x26A0;</span>'
      : '';

    row.innerHTML = `
      <td><input type="checkbox" class="ml-row-check" ${checked}></td>
      <td class="ml-name" title="${_esc(asset.path)}">${_esc(asset.name)}${compatWarning}</td>
      <td>${_typeLabel(asset.type)}</td>
      <td>${familyBadge}</td>
      <td>${asset.format}</td>
      <td>${sizeFmt}</td>`;

    row.querySelector('.ml-row-check').addEventListener('change', (e) => {
      if (e.target.checked) {
        _selectedPaths.add(asset.path);
      } else {
        _selectedPaths.delete(asset.path);
      }
      _updateSelectionLabel();
    });
    return row;
  }

  function _renderInstallSummary(summary) {
    const tbody = document.getElementById('ml-table-body');
    if (!tbody) return;
    // Annotate rows with installed/skipped/failed status
    const statusMap = new Map();
    for (const a of (summary.registered    || [])) statusMap.set(a.original_path, 'installed');
    for (const a of (summary.already_known || [])) statusMap.set(a.original_path, 'skipped');
    for (const a of (summary.failed        || [])) statusMap.set(a.original_path, 'failed');

    tbody.querySelectorAll('.ml-asset-row').forEach(row => {
      const status = statusMap.get(row.dataset.path);
      if (status) row.classList.add(`ml-status-${status}`);
    });
  }

  // ── selection helpers ──────────────────────────────────────────────────────

  function _selectAll() {
    if (!_scanResult) return;
    document.querySelectorAll('.ml-row-check').forEach(cb => { cb.checked = true; });
    for (const group of ['models','vaes','loras','controlnets','ip_adapters','clip_visions','embeddings','upscalers']) {
      (_scanResult[group] || []).forEach(a => _selectedPaths.add(a.path));
    }
    _updateSelectionLabel();
  }

  function _selectNone() {
    document.querySelectorAll('.ml-row-check').forEach(cb => { cb.checked = false; });
    _selectedPaths.clear();
    _updateSelectionLabel();
  }

  function _updateSelectionLabel() {
    const el = document.getElementById('ml-selection-count');
    if (el) el.textContent = `${_selectedPaths.size} selected`;
  }

  // ── UI helpers ─────────────────────────────────────────────────────────────

  function _setScanStatus(msg, state) {
    const el = document.getElementById('ml-status');
    if (!el) return;
    el.textContent = msg;
    el.className = `ml-status ml-status-${state}`;
  }

  function _setInstallEnabled(enabled) {
    ['ml-btn-register-all', 'ml-btn-register-selected'].forEach(id => {
      const btn = document.getElementById(id);
      if (btn) btn.disabled = !enabled;
    });
  }

  function _totalCount(result) {
    return ['models','vaes','loras','controlnets','ip_adapters','clip_visions','embeddings','upscalers']
      .reduce((n, k) => n + (result[k] || []).length, 0);
  }

  function _typeLabel(type) {
    const map = {
      diffusion: 'Diffusion Model', vae: 'VAE',
      lora: 'LoRA', controlnet: 'ControlNet',
      ip_adapter: 'IP-Adapter', clip_vision: 'CLIP-Vision',
      embedding: 'Embedding', upscaler: 'Upscaler',
      clip: 'CLIP Encoder', t5: 'T5 Encoder',
      unknown: 'Unknown',
    };
    return map[type] || type;
  }

  function _esc(s) {
    return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');
  }

  return { init, scan, installAll, installSelected };
})();

// ── bootstrap ─────────────────────────────────────────────────────────────────
document.addEventListener('DOMContentLoaded', () => ModelLibrary.init());
