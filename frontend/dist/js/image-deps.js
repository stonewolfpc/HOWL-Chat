/**
 * Image Model Dependencies
 *
 * Explains what a diffusion model needs before it can load, and — with the
 * user's explicit consent — fetches the missing files.
 *
 * HOWL Chat stays offline by default: nothing here touches the network until
 * the consent modal is accepted. If the connection fails, the modal lists the
 * manual link for every file and the rest of the app keeps working.
 *
 * Usage:
 *   if (!(await ImageDeps.ensureReady(modelName))) return; // user cancelled
 *
 * @module image-deps
 */

'use strict';

const ImageDeps = (() => {

    const OVERLAY_ID = 'image-deps-overlay';
    let _listenerAttached = false;
    let _progressByRole = {};

    function _app() {
        return window.go?.main?.App;
    }

    function _removeOverlay() {
        document.getElementById(OVERLAY_ID)?.remove();
    }

    function _fmtMB(mb) {
        if (!mb) return 'size unknown';
        return mb >= 1024 ? (mb / 1024).toFixed(1) + ' GB' : mb + ' MB';
    }

    function _itemRow(item) {
        const present = item.state === 'present';
        const optional = item.state === 'optional';
        const mismatch = item.state === 'mismatch';
        const icon = present ? '✓' : (optional ? '·' : (mismatch ? '!' : '✕'));
        const colour = present ? 'var(--success, #4ade80)'
            : (optional ? 'var(--text-tertiary)' : (mismatch ? 'var(--warning, #f0b429)' : 'var(--error, #ff6b6b)'));
        const detail = present
            ? `<div style="color:var(--text-tertiary);font-size:11px;margin-top:2px;word-break:break-all;">${_esc(item.path)}</div>`
            : `<div style="color:var(--text-secondary);font-size:12px;margin-top:2px;">
                   ${mismatch
                     ? `<div style="color:var(--warning,#f0b429);font-size:12px;margin-bottom:2px;">Found <code>${_esc((item.path || '').split(/[\/]/).pop())}</code>, which cannot be used: ${_esc(item.mismatch || 'wrong kind of file')}</div>`
                     : ''}
                   ${item.fileName ? `Needs <strong>${_esc(item.fileName)}</strong> · ${_fmtMB(item.sizeMB)}` : 'No automatic download available'}
                   ${item.source ? ` · <a href="#" data-dep-link="${_esc(item.source)}" style="color:var(--accent,#7aa2f7);">Get it here</a>` : ''}
               </div>
               ${item.note ? `<div style="color:var(--text-tertiary);font-size:11px;margin-top:2px;">${_esc(item.note)}</div>` : ''}
               ${item.error ? `<div style="color:var(--error,#ff6b6b);font-size:12px;margin-top:4px;">${_esc(item.error)}</div>` : ''}`;
        return `<div style="display:flex;gap:8px;align-items:flex-start;padding:8px 0;border-bottom:1px solid rgba(255,255,255,0.06);">
                    <span style="color:${colour};font-weight:700;width:14px;">${icon}</span>
                    <div style="flex:1;">
                        <div style="color:var(--text-primary,#fff);font-size:13px;">${_esc(item.label)}
                            ${optional ? '<span style="color:var(--text-tertiary);font-size:11px;"> (optional)</span>' : ''}
                        </div>
                        ${detail}
                    </div>
                </div>`;
    }

    function _esc(s) {
        return String(s == null ? '' : s)
            .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;').replace(/'/g, '&#39;');
    }

    function _wireLinks(root) {
        root.querySelectorAll('[data-dep-link]').forEach(a => {
            a.addEventListener('click', (e) => {
                e.preventDefault();
                const url = a.getAttribute('data-dep-link');
                if (window.runtime?.BrowserOpenURL) window.runtime.BrowserOpenURL(url);
                else window.open(url, '_blank', 'noopener');
            });
        });
    }

    function _onProgress(data) {
        if (!data || !data.role) return;
        _progressByRole[data.role] = data;
        const bar = document.getElementById('dep-bar-' + data.role);
        const txt = document.getElementById('dep-txt-' + data.role);
        if (!bar) return;
        if (data.total > 0) {
            bar.style.width = Math.min(100, Math.round((data.received / data.total) * 100)) + '%';
            if (txt) txt.textContent = (data.received / 1048576).toFixed(0) + ' / ' + (data.total / 1048576).toFixed(0) + ' MB';
        } else if (txt) {
            txt.textContent = (data.received / 1048576).toFixed(0) + ' MB';
        }
    }

    function _attachListener() {
        if (_listenerAttached || !window.runtime?.EventsOn) return;
        window.runtime.EventsOn('image-deps:progress', _onProgress);
        _listenerAttached = true;
    }

    /**
     * Show the modal for a report and resolve to true once the model's required
     * files are present, or false if the user declined or the download failed.
     */
    function _prompt(report) {
        return new Promise((resolve) => {
            _removeOverlay();
            _progressByRole = {};

            const missingItems = report.items.filter(i => i.state === 'missing');
            const downloadable = missingItems.filter(i => i.url);
            const rows = report.items.map(_itemRow).join('');

            const overlay = document.createElement('div');
            overlay.id = OVERLAY_ID;
            overlay.style.cssText = `
                position: fixed; inset: 0; z-index: 10000;
                background: var(--overlay-scrim, rgba(0,0,0,0.6));
                display: flex; align-items: center; justify-content: center;
            `;
            overlay.innerHTML = `
                <div style="background:var(--surface-panel,#1e1e28);border:1px solid rgba(122,162,247,0.35);
                            border-radius:12px;padding:24px;max-width:560px;width:92%;
                            max-height:86vh;overflow:auto;box-shadow:0 8px 32px rgba(0,0,0,0.6);">
                    <h3 style="color:#fff;margin:0 0 4px;font-size:17px;">Missing dependencies</h3>
                    <p style="color:var(--text-secondary);margin:0 0 12px;font-size:13px;">
                        <strong>${_esc(report.model)}</strong> is a <strong>${_esc(report.family || 'unknown')}</strong> model and
                        cannot load without the files below.
                    </p>
                    <div id="dep-list">${rows}</div>
                    <p id="dep-consent" style="color:var(--text-secondary);font-size:12px;line-height:1.5;
                            background:rgba(122,162,247,0.08);border-left:2px solid rgba(122,162,247,0.5);
                            padding:10px;border-radius:4px;margin:14px 0 0;">
                        ${_esc(report.consentText)}
                        ${downloadable.length
                            ? `<br><span style="color:var(--text-tertiary);">About ${_fmtMB(report.totalSizeMB)} will be saved to <code>${_esc(report.downloadDir)}</code>.</span>`
                            : `<br><span style="color:var(--text-tertiary);">No automatic download is available for these files; open the links above.</span>`}
                    </p>
                    <div id="dep-error" style="display:none;color:var(--error,#ff6b6b);font-size:12px;
                            margin-top:10px;line-height:1.5;"></div>
                    <div style="display:flex;gap:10px;margin-top:18px;">
                        <button id="dep-cancel" style="flex:1;padding:10px;background:var(--surface-hover,#2a2a38);
                                border:1px solid rgba(255,255,255,0.15);color:#fff;border-radius:8px;
                                cursor:pointer;font-size:14px;">Cancel</button>
                        <button id="dep-continue" ${downloadable.length ? '' : 'disabled'} style="flex:1;padding:10px;
                                background:${downloadable.length ? 'rgba(80,140,255,0.9)' : 'rgba(90,90,110,0.5)'};
                                border:1px solid rgba(122,162,247,0.4);color:#fff;border-radius:8px;
                                cursor:${downloadable.length ? 'pointer' : 'not-allowed'};font-size:14px;font-weight:600;">
                            ${downloadable.length ? 'Continue' : 'Close'}
                        </button>
                    </div>
                </div>
            `;
            document.body.appendChild(overlay);
            _wireLinks(overlay);

            const close = (value) => { _removeOverlay(); resolve(value); };
            overlay.querySelector('#dep-cancel').addEventListener('click', () => close(false));
            overlay.addEventListener('click', (e) => { if (e.target === overlay) close(false); });

            const btn = overlay.querySelector('#dep-continue');
            if (!downloadable.length) {
                btn.addEventListener('click', () => close(false));
                return;
            }

            btn.addEventListener('click', async () => {
                btn.disabled = true;
                btn.textContent = 'Downloading…';
                overlay.querySelector('#dep-cancel').textContent = 'Cancel (stop after this file)';
                overlay.querySelector('#dep-consent').style.display = 'none';

                // One progress row per file being fetched.
                const list = overlay.querySelector('#dep-list');
                list.innerHTML = downloadable.map(it => `
                    <div style="padding:8px 0;">
                        <div style="color:#fff;font-size:13px;">${_esc(it.label)}
                            <span id="dep-txt-${_esc(it.role)}" style="color:var(--text-tertiary);font-size:11px;"></span>
                        </div>
                        <div style="height:5px;background:rgba(255,255,255,0.08);border-radius:3px;margin-top:6px;overflow:hidden;">
                            <div id="dep-bar-${_esc(it.role)}" style="height:100%;width:0;background:var(--accent,#7aa2f7);transition:width 0.2s;"></div>
                        </div>
                    </div>
                `).join('');

                _attachListener();

                let result;
                try {
                    result = await _app().DownloadModelDependencies(report.model);
                } catch (err) {
                    const box = overlay.querySelector('#dep-error');
                    box.style.display = 'block';
                    box.textContent = 'Download failed: ' + (err?.message || err) +
                        ' — the app is still usable; open the links above to get the files manually.';
                    btn.disabled = false;
                    btn.textContent = 'Retry';
                    return;
                }

                const stillMissing = result.missing || [];
                if (stillMissing.length === 0) {
                    btn.textContent = 'Loaded';
                    _removeOverlay();
                    resolve(true);
                    return;
                }

                // Report what is still missing, with the manual links restored.
                list.innerHTML = result.items.map(_itemRow).join('');
                _wireLinks(overlay);
                const box = overlay.querySelector('#dep-error');
                box.style.display = 'block';
                box.textContent = result.message || ('Still missing: ' + stillMissing.join(', '));
                btn.disabled = false;
                btn.textContent = result.canDownload ? 'Retry' : 'Close';
                if (!result.canDownload) {
                    btn.onclick = () => close(false);
                }
                btn.addEventListener('click', () => close(false), { once: true });
            });
        });
    }

    /**
     * Ensure a model's required dependencies are present, offering a download
     * when they are not. Resolves true when the model can be loaded.
     *
     * A failed lookup must never block loading: if the report cannot be built
     * the loader is allowed to try and report its own error.
     */
    async function ensureReady(modelName) {
        if (!modelName || !_app()?.GetModelDependencies) return true;
        let report;
        try {
            report = await _app().GetModelDependencies(modelName);
        } catch (err) {
            console.warn('[deps] lookup failed, letting the loader decide:', err);
            return true;
        }
        if (!report || !report.missing || report.missing.length === 0) return true;
        if (!report.known) return true; // unsupported family: loader gives its own message
        return await _prompt(report);
    }

    return { ensureReady };
})();

window.ImageDeps = ImageDeps;
