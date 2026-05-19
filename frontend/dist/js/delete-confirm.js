/**
 * Delete Confirm
 *
 * Reusable "Deletion is permanent, proceed?" modal.
 * Replaces all raw confirm() calls for destructive actions.
 *
 * Usage:
 *   DeleteConfirm.show({ title, body, onConfirm })
 *
 * @module delete-confirm
 * @version 1.0.0
 */

'use strict';

const DeleteConfirm = (() => {

    const OVERLAY_ID = 'delete-confirm-overlay';

    function _removeExisting() {
        const existing = document.getElementById(OVERLAY_ID);
        if (existing) existing.remove();
    }

    /**
     * Show the permanent-deletion confirmation modal.
     * @param {object} opts
     * @param {string}   opts.title      - Bold heading, e.g. 'Delete "My World"?'
     * @param {string}   [opts.body]     - Optional detail line
     * @param {function} opts.onConfirm  - Called when user clicks Delete
     * @param {function} [opts.onCancel] - Called when user cancels (optional)
     */
    function show({ title, body, onConfirm, onCancel }) {
        _removeExisting();

        const overlay = document.createElement('div');
        overlay.id = OVERLAY_ID;
        overlay.style.cssText = `
            position: fixed; inset: 0; z-index: 10000;
            background: rgba(0,0,0,0.7);
            display: flex; align-items: center; justify-content: center;
        `;

        const bodyText = body
            ? `<p style="color:rgba(255,255,255,0.65);margin:8px 0 0;font-size:14px;">${body}</p>`
            : '';

        overlay.innerHTML = `
            <div style="
                background:#1e1e28; border:1px solid rgba(255,80,80,0.3);
                border-radius:12px; padding:28px 32px; max-width:420px; width:90%;
                box-shadow:0 8px 32px rgba(0,0,0,0.6);
            ">
                <h3 style="color:#fff;margin:0 0 4px;font-size:18px;">${title}</h3>
                ${bodyText}
                <p style="color:#ff6b6b;font-weight:600;margin:16px 0 0;font-size:13px; display:flex; align-items:center; gap:6px;">
                    <span class="icon icon-warning" style="width:14px;height:14px;"></span> Deletion is permanent and cannot be undone.
                </p>
                <div style="display:flex;gap:12px;margin-top:24px;">
                    <button id="dc-cancel" style="
                        flex:1;padding:10px;background:rgba(255,255,255,0.08);
                        border:1px solid rgba(255,255,255,0.15);color:#fff;
                        border-radius:8px;cursor:pointer;font-size:14px;
                    ">Cancel</button>
                    <button id="dc-confirm" style="
                        flex:1;padding:10px;background:rgba(220,50,50,0.85);
                        border:1px solid rgba(255,80,80,0.4);color:#fff;
                        border-radius:8px;cursor:pointer;font-size:14px;font-weight:600;
                    ">Delete</button>
                </div>
            </div>
        `;

        document.body.appendChild(overlay);

        overlay.querySelector('#dc-cancel').addEventListener('click', () => {
            overlay.remove();
            if (typeof onCancel === 'function') onCancel();
        });

        overlay.querySelector('#dc-confirm').addEventListener('click', () => {
            overlay.remove();
            if (typeof onConfirm === 'function') onConfirm();
        });

        overlay.addEventListener('click', (e) => {
            if (e.target === overlay) {
                overlay.remove();
                if (typeof onCancel === 'function') onCancel();
            }
        });
    }

    return { show };

})();

window.DeleteConfirm = DeleteConfirm;
