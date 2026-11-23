/**
 * Error Reporting settings — controls in the Privacy & Data panel.
 *
 * The backend is the single authority: `IsErrorReportingEnabled` reads the
 * preference persisted in error-reporting.json, `SetErrorReportingEnabled`
 * writes it, and `GetErrorLog` reads the same log the diagnostics owner writes.
 * There is deliberately no localStorage mirror, so a restart always shows the
 * persisted value rather than a stale browser copy.
 */

'use strict';

window.ErrorSettings = {
    refresh,
    setEnabled,
    preview,
    send
};

function appBinding(name) {
    return window.go?.main?.App?.[name];
}

function reportCheckbox() {
    return document.getElementById('error-reporting-enabled');
}

function sendButton() {
    return document.getElementById('error-report-send');
}

function setStatus(text) {
    const el = document.getElementById('error-report-status');
    if (el) el.textContent = text || '';
}

function setSendEnabled(enabled) {
    const btn = sendButton();
    if (btn) btn.disabled = !enabled;
}

// ---------------------------------------------------------------------------
// Toggle
// ---------------------------------------------------------------------------

async function refresh() {
    await hydrateToggle();
    await refreshLog();
}

async function hydrateToggle() {
    const checkbox = reportCheckbox();
    if (!checkbox) return;
    const read = appBinding('IsErrorReportingEnabled');
    if (!read) {
        checkbox.checked = false;
        checkbox.disabled = true;
        setSendEnabled(false);
        setStatus('Backend unavailable');
        return;
    }
    try {
        const enabled = !!(await read());
        checkbox.checked = enabled;
        checkbox.disabled = false;
        setSendEnabled(enabled);
        setStatus(enabled ? 'Error reporting enabled' : 'Error reporting disabled');
    } catch (err) {
        checkbox.checked = false;
        setSendEnabled(false);
        setStatus('Could not read the preference: ' + (err?.message || err));
    }
}

async function setEnabled(enabled) {
    const checkbox = reportCheckbox();
    const write = appBinding('SetErrorReportingEnabled');
    if (!write) {
        if (checkbox) checkbox.checked = !enabled;
        setStatus('Backend unavailable');
        return;
    }
    try {
        await write(enabled);
        setSendEnabled(enabled);
        setStatus(enabled ? 'Error reporting enabled' : 'Error reporting disabled');
        await refreshLog();
    } catch (err) {
        if (checkbox) checkbox.checked = !enabled;
        setSendEnabled(checkbox ? checkbox.checked : false);
        setStatus('Could not save the preference: ' + (err?.message || err));
    }
}

// ---------------------------------------------------------------------------
// Report actions
// ---------------------------------------------------------------------------

async function preview() {
    const list = document.getElementById('error-log-list');
    if (!list) return;
    const getPreview = appBinding('GetErrorReportPreview');
    if (!getPreview) {
        setStatus('Backend unavailable');
        return;
    }
    setStatus('Generating preview...');
    try {
        const preview = await getPreview();
        list.innerHTML = '<pre style="white-space:pre-wrap;word-break:break-all;">' +
            escapeHtml(preview) + '</pre>';
        setStatus('Preview only — nothing was uploaded');
    } catch (err) {
        setStatus('Preview failed: ' + (err?.message || err));
    }
}

async function send() {
    const btn = sendButton();
    if (!btn || btn.disabled) return;
    const sendReport = appBinding('SendErrorReport');
    if (!sendReport) {
        setStatus('Backend unavailable');
        return;
    }
    const original = btn.textContent;
    btn.textContent = 'Sending...';
    btn.disabled = true;
    try {
        const reportId = await sendReport();
        btn.textContent = 'Queued: ' + String(reportId).slice(-8);
        setStatus('Report queued for upload — failures appear under Recent Errors');
        setTimeout(() => {
            btn.textContent = original;
            setSendEnabled(true);
        }, 3000);
    } catch (err) {
        btn.textContent = 'Failed';
        setStatus('Send failed: ' + (err?.message || err));
        setTimeout(() => {
            btn.textContent = original;
            setSendEnabled(true);
        }, 3000);
    }
}

// ---------------------------------------------------------------------------
// Log
// ---------------------------------------------------------------------------

async function refreshLog() {
    const list = document.getElementById('error-log-list');
    if (!list) return;
    const getLog = appBinding('GetErrorLog');
    if (!getLog) {
        list.innerHTML = '<em>Backend unavailable</em>';
        return;
    }
    try {
        const entries = await getLog(50);
        if (!entries || entries.length === 0) {
            list.innerHTML = '<em>No errors logged</em>';
            return;
        }
        list.innerHTML = entries.map((e) => `
            <div style="padding:4px 0;border-bottom:1px solid var(--border);">
                <span style="color:${e.severity === 'critical' ? '#ef4444' : e.severity === 'high' ? '#f59e0b' : '#6b7280'};font-weight:600;">${escapeHtml(e.code)}</span>
                <span style="margin-left:8px;">${escapeHtml(e.operation)}</span>
                <span style="margin-left:8px;font-size:11px;">${escapeHtml(e.location || '')}</span>
            </div>
        `).join('');
    } catch (err) {
        list.innerHTML = '<em>Unable to load errors</em>';
    }
}

// ---------------------------------------------------------------------------
// Wiring
// ---------------------------------------------------------------------------

// The panel exists in the page from the start, but the log should reflect the
// moment the user looks at it, so refresh whenever Privacy & Data is selected.
document.addEventListener('click', (event) => {
    const nav = event.target?.closest?.('.real-settings-nav-item');
    if (!nav || !/Privacy/.test(nav.textContent || '')) return;
    setTimeout(refresh, 0);
});

document.addEventListener('DOMContentLoaded', () => {
    if (document.getElementById('real-settings-overlay')) {
        refresh();
    }
});
