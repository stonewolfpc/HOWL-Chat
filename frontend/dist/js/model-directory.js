/**
 * Model Directory UI
 *
 * Replaces the native Windows folder dialog with a themed, resolution-aware
 * locator, and presents the "new model folder detected" prompt.
 *
 * Nothing in this file loads a model. "Scan directory" classifies files and
 * fills the pickers; it never starts inference or reserves VRAM.
 *
 * @module model-directory
 */

'use strict';

(function () {
    const api = () => window.go?.main?.App || null;

    let _locatorResolve = null;
    let _locatorBrowse = null;
    let _promptBusy = false;

    function el(id) {
        return document.getElementById(id);
    }

    function show(id) {
        const node = el(id);
        if (node) node.classList.add('show');
    }

    function hide(id) {
        const node = el(id);
        if (node) node.classList.remove('show');
    }

    function notify(message, kind) {
        if (typeof showNotification === 'function') {
            showNotification(message, kind || 'info');
        }
    }

    function kindLabel(kind) {
        switch (kind) {
            case 'llm': return 'chat models';
            case 'tts': return 'voice models';
            case 'image': return 'image models';
            default: return 'models';
        }
    }

    /* ------------------------------------------------------------------ */
    /* Folder locator                                                     */
    /* ------------------------------------------------------------------ */

    async function openDirectoryLocator(startPath) {
        const app = api();

        // The in-app locator is the product path. The native dialog stays only
        // as a fallback for a build where the binding is missing.
        if (typeof app?.BrowseModelDirectory !== 'function') {
            if (typeof app?.OpenDirectoryPicker === 'function') {
                return app.OpenDirectoryPicker();
            }
            notify('Folder browser is unavailable in this build.', 'error');
            return '';
        }

        show('md-locator-overlay');
        setLocatorStatus('Looking for model files…');
        await browse(startPath || '');

        return new Promise((resolve) => {
            _locatorResolve = resolve;
        });
    }

    function finishLocator(result) {
        hide('md-locator-overlay');
        const resolve = _locatorResolve;
        _locatorResolve = null;
        if (typeof resolve === 'function') resolve(result || '');
    }

    async function browse(path) {
        const app = api();
        if (typeof app?.BrowseModelDirectory !== 'function') return;

        try {
            const result = await app.BrowseModelDirectory(path || '');
            _locatorBrowse = result;
            renderLocator(result);
        } catch (err) {
            setLocatorStatus('Cannot open that folder: ' + (err?.message || err), true);
        }
    }

    function renderLocator(result) {
        const pathInput = el('md-path');
        if (pathInput) pathInput.value = result.path || '';

        const up = el('md-up');
        if (up) {
            up.disabled = !!result.at_root;
            up.dataset.parent = result.parent || '';
        }

        const mounts = el('md-mounts');
        if (mounts) {
            mounts.innerHTML = '';
            if (result.home) {
                mounts.appendChild(mountButton('Home', result.home, result.path === result.home));
            }
            for (const drive of result.drives || []) {
                mounts.appendChild(mountButton(drive, drive, result.path === drive));
            }
        }

        const list = el('md-folders');
        if (list) {
            list.innerHTML = '';
            const folders = result.folders || [];
            if (folders.length === 0) {
                const empty = document.createElement('div');
                empty.className = 'md-empty';
                empty.textContent = 'No subfolders. This folder itself can be selected.';
                list.appendChild(empty);
            }
            for (const folder of folders) {
                list.appendChild(folderRow(folder));
            }
        }

        const files = result.model_files || 0;
        if (files > 0) {
            setLocatorStatus(files + ' model file' + (files === 1 ? '' : 's') + ' detected · will populate ' + kindLabel(result.kind));
        } else {
            setLocatorStatus('No model files detected in this folder.');
        }

        const select = el('md-select');
        if (select) select.disabled = !result.path;
    }

    function mountButton(label, path, active) {
        const button = document.createElement('button');
        button.type = 'button';
        button.className = 'md-mount' + (active ? ' active' : '');
        button.textContent = label;
        button.title = path;
        button.addEventListener('click', () => browse(path));
        return button;
    }

    function folderRow(folder) {
        const row = document.createElement('button');
        row.type = 'button';
        row.className = 'md-row';
        row.title = 'Open ' + folder.path;

        const name = document.createElement('span');
        name.className = 'md-row-name';
        name.textContent = folder.name;
        row.appendChild(name);

        if (folder.model_files > 0) {
            const badge = document.createElement('span');
            badge.className = 'md-badge';
            badge.textContent = folder.model_files + ' model' + (folder.model_files === 1 ? '' : 's');
            row.appendChild(badge);
        }

        row.addEventListener('click', () => browse(folder.path));
        return row;
    }

    function setLocatorStatus(text, isError) {
        const status = el('md-status-text');
        if (!status) return;
        status.textContent = text;
        status.style.color = isError ? 'var(--danger, #ef4444)' : '';
    }

    function wireLocator() {
        const overlay = el('md-locator-overlay');
        if (!overlay) return;

        el('md-close')?.addEventListener('click', () => finishLocator(''));
        el('md-cancel')?.addEventListener('click', () => finishLocator(''));
        el('md-select')?.addEventListener('click', () => {
            finishLocator(_locatorBrowse?.path || '');
        });

        el('md-up')?.addEventListener('click', () => {
            if (_locatorBrowse && !_locatorBrowse.at_root) browse(_locatorBrowse.parent);
        });

        el('md-go')?.addEventListener('click', () => {
            const value = el('md-path')?.value || '';
            browse(value.trim());
        });

        el('md-path')?.addEventListener('keydown', (event) => {
            if (event.key === 'Enter') {
                event.preventDefault();
                browse((el('md-path')?.value || '').trim());
            }
        });

        overlay.addEventListener('click', (event) => {
            if (event.target === overlay) finishLocator('');
        });

        document.addEventListener('keydown', (event) => {
            if (event.key === 'Escape' && overlay.classList.contains('show')) {
                finishLocator('');
            }
        });
    }

    /* ------------------------------------------------------------------ */
    /* New model folder prompt                                            */
    /* ------------------------------------------------------------------ */

    async function refreshDiscoveryPrompt() {
        const app = api();
        if (typeof app?.GetModelDiscoveryState !== 'function') return;

        try {
            const state = await app.GetModelDiscoveryState();
            syncDiscoveryToggle(state.prompting);
            if (!state.prompting || !state.candidate) {
                hide('md-prompt-overlay');
                return;
            }
            renderPrompt(state.candidate);
            show('md-prompt-overlay');
        } catch (err) {
            console.warn('[discovery] state unavailable:', err);
        }
    }

    function renderPrompt(candidate) {
        const pathNode = el('md-prompt-path');
        if (pathNode) pathNode.textContent = candidate.path || '';

        const summary = el('md-prompt-summary');
        if (summary) {
            const kinds = (candidate.kinds || []).join(', ');
            const files = candidate.model_files || 0;
            summary.textContent = files + ' model file' + (files === 1 ? '' : 's')
                + (kinds ? ' · ' + kinds : '');
        }

        const scan = el('md-prompt-scan');
        if (scan) scan.disabled = false;
    }

    async function promptScan() {
        const app = api();
        if (_promptBusy || typeof app?.ScanDiscoveredDirectory !== 'function') return;
        const path = el('md-prompt-path')?.textContent || '';
        if (!path) return;

        _promptBusy = true;
        const scan = el('md-prompt-scan');
        if (scan) {
            scan.disabled = true;
            scan.textContent = 'Scanning…';
        }
        try {
            const dir = await app.ScanDiscoveredDirectory(path);
            notify('Added ' + (dir?.path || path) + ' (' + (dir?.kind || 'models') + ')', 'success');
            hide('md-prompt-overlay');
            await refreshAfterDirectoryChange();
        } catch (err) {
            notify('Scan failed: ' + (err?.message || err), 'error');
        } finally {
            _promptBusy = false;
            if (scan) {
                scan.disabled = false;
                scan.textContent = 'Scan directory';
            }
        }
    }

    async function promptDefer() {
        const app = api();
        const path = el('md-prompt-path')?.textContent || '';
        if (typeof app?.DeferDiscoveredDirectory === 'function' && path) {
            try {
                await app.DeferDiscoveredDirectory(path);
            } catch (err) {
                console.warn('[discovery] defer failed:', err);
            }
        }
        hide('md-prompt-overlay');
    }

    async function promptSuppress() {
        const app = api();
        if (typeof app?.SetDirectoryDiscovery === 'function') {
            try {
                // Prompts off; discovery keeps running silently in the background
                // so the pickers stay accurate without interruptions.
                await app.SetDirectoryDiscovery(false);
                syncDiscoveryToggle(false);
            } catch (err) {
                console.warn('[discovery] preference save failed:', err);
            }
        }
        hide('md-prompt-overlay');
    }

    function wirePrompt() {
        el('md-prompt-scan')?.addEventListener('click', promptScan);
        el('md-prompt-later')?.addEventListener('click', promptDefer);
        el('md-prompt-never')?.addEventListener('click', promptSuppress);
        el('md-prompt-close')?.addEventListener('click', promptDefer);

        el('discovery-prompts')?.addEventListener('change', async (event) => {
            const app = api();
            if (typeof app?.SetDirectoryDiscovery !== 'function') return;
            try {
                await app.SetDirectoryDiscovery(!!event.target.checked);
                if (event.target.checked) await refreshDiscoveryPrompt();
            } catch (err) {
                notify('Could not save that preference.', 'error');
            }
        });
    }

    function syncDiscoveryToggle(enabled) {
        const toggle = el('discovery-prompts');
        if (toggle) toggle.checked = !!enabled;
    }

    async function refreshAfterDirectoryChange() {
        if (typeof populateLLMDropdown === 'function') {
            try { await populateLLMDropdown(); } catch (err) { console.warn(err); }
        }
        if (typeof populateModelBrowserDirs === 'function') {
            try { await populateModelBrowserDirs(); } catch (err) { console.warn(err); }
        }
        if (typeof populateModelBrowser === 'function') {
            try { await populateModelBrowser(); } catch (err) { console.warn(err); }
        }
    }

    /* ------------------------------------------------------------------ */

    document.addEventListener('DOMContentLoaded', () => {
        wireLocator();
        wirePrompt();
        refreshDiscoveryPrompt();
    });

    // The model browser opens later; ask again then so a folder found while the
    // app was running is offered without restarting.
    window.openModelBrowser = (function (original) {
        return function () {
            if (typeof original === 'function') original.apply(this, arguments);
            refreshDiscoveryPrompt();
        };
    })(window.openModelBrowser);

    window.openDirectoryLocator = openDirectoryLocator;
    window.refreshModelDiscoveryPrompt = refreshDiscoveryPrompt;
})();
