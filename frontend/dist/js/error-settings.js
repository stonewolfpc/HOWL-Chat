/**
 * Error Settings - Vercel upload toggle and log viewer
 * Lives under File > Settings in the chat page topbar.
 */

'use strict';

// ---------------------------------------------------------------------------
// Settings UI Integration
// ---------------------------------------------------------------------------

function injectErrorSettings() {
    const settingsPanel = document.getElementById('settings-panel');
    if (!settingsPanel) return;

    // Avoid duplicate injection
    if (document.getElementById('error-settings-section')) return;

    const section = document.createElement('div');
    section.id = 'error-settings-section';
    section.className = 'settings-section';
    section.innerHTML = `
        <h3>Error Reporting</h3>
        <div class="settings-group">
            <label class="setting-row">
                <span>Send error reports to Vercel</span>
                <input type="checkbox" id="toggle-error-reporting" onchange="onErrorReportingToggle(this)">
            </label>
            <p class="setting-desc" style="font-size:12px;color:var(--text-secondary);margin-top:4px;">
                Reports contain only error codes and line numbers. No personal data, no messages, no model names.
            </p>
        </div>
        <div class="settings-group" style="margin-top:12px;">
            <button class="btn btn-small" onclick="previewErrorReport()" id="btn-preview-report">Preview Report</button>
            <button class="btn btn-small" onclick="sendErrorReportNow()" id="btn-send-report" style="margin-left:8px;" disabled>Send Report</button>
        </div>
        <div class="settings-group" style="margin-top:12px;">
            <details>
                <summary style="cursor:pointer;font-size:13px;color:var(--text-secondary);">View Recent Errors (${ErrorSystem.activeErrors?.size || 0})</summary>
                <div id="error-log-list" style="margin-top:8px;max-height:200px;overflow-y:auto;font-size:12px;font-family:monospace;background:var(--bg-tertiary);padding:8px;border-radius:4px;">
                    <em style="color:var(--text-secondary);">Loading...</em>
                </div>
            </details>
        </div>
    `;

    settingsPanel.appendChild(section);
    loadErrorReportingState();
    refreshErrorLogList();
}

// ---------------------------------------------------------------------------
// Toggle handling
// ---------------------------------------------------------------------------

async function onErrorReportingToggle(checkbox) {
    const enabled = checkbox.checked;
    try {
        if (window.go?.main?.App?.SetErrorReportingEnabled) {
            await window.go.main.App.SetErrorReportingEnabled(enabled);
        }
        localStorage.setItem('howl-error-reporting', enabled ? '1' : '0');
        updateSendButton(enabled);
    } catch (err) {
        console.error('Failed to save error reporting preference:', err);
        checkbox.checked = !enabled;
    }
}

async function loadErrorReportingState() {
    const checkbox = document.getElementById('toggle-error-reporting');
    if (!checkbox) return;

    try {
        if (window.go?.main?.App?.IsErrorReportingEnabled) {
            const enabled = await window.go.main.App.IsErrorReportingEnabled();
            checkbox.checked = enabled;
            updateSendButton(enabled);
        } else {
            const cached = localStorage.getItem('howl-error-reporting') === '1';
            checkbox.checked = cached;
            updateSendButton(cached);
        }
    } catch {
        checkbox.checked = false;
        updateSendButton(false);
    }
}

function updateSendButton(enabled) {
    const btn = document.getElementById('btn-send-report');
    if (btn) btn.disabled = !enabled;
}

// ---------------------------------------------------------------------------
// Report actions
// ---------------------------------------------------------------------------

async function previewErrorReport() {
    const list = document.getElementById('error-log-list');
    if (!list) return;

    list.innerHTML = '<em>Generating preview...</em>';
    try {
        if (window.go?.main?.App?.GetErrorReportPreview) {
            const preview = await window.go.main.App.GetErrorReportPreview();
            list.innerHTML = `<pre style="white-space:pre-wrap;word-break:break-all;">${escapeHtml(preview)}</pre>`;
        } else {
            list.innerHTML = '<em>Backend not available</em>';
        }
    } catch (err) {
        list.innerHTML = `<em style="color:#ef4444;">Error: ${escapeHtml(err.message || String(err))}</em>`;
    }
}

async function sendErrorReportNow() {
    const btn = document.getElementById('btn-send-report');
    if (!btn || btn.disabled) return;

    const originalText = btn.textContent;
    btn.textContent = 'Sending...';
    btn.disabled = true;

    try {
        if (window.go?.main?.App?.SendErrorReport) {
            const reportId = await window.go.main.App.SendErrorReport();
            btn.textContent = 'Sent: ' + reportId.slice(-8);
            setTimeout(() => { btn.textContent = originalText; btn.disabled = false; }, 3000);
        }
    } catch (err) {
        btn.textContent = 'Failed';
        setTimeout(() => { btn.textContent = originalText; btn.disabled = false; }, 3000);
    }
}

// ---------------------------------------------------------------------------
// Error log list
// ---------------------------------------------------------------------------

async function refreshErrorLogList() {
    const list = document.getElementById('error-log-list');
    if (!list) return;

    try {
        if (window.go?.main?.App?.GetErrorLog) {
            const entries = await window.go.main.App.GetErrorLog(50);
            if (!entries || entries.length === 0) {
                list.innerHTML = '<em style="color:var(--text-secondary);">No errors logged</em>';
                return;
            }
            list.innerHTML = entries.map(e => `
                <div style="padding:4px 0;border-bottom:1px solid var(--border);">
                    <span style="color:${e.severity === 'critical' ? '#ef4444' : e.severity === 'high' ? '#f59e0b' : '#6b7280'};font-weight:600;">${escapeHtml(e.code)}</span>
                    <span style="color:var(--text-secondary);margin-left:8px;">${escapeHtml(e.operation)}</span>
                    <span style="color:var(--text-muted);margin-left:8px;font-size:11px;">${escapeHtml(e.location || '')}</span>
                </div>
            `).join('');
        } else {
            // Fallback: show frontend-only errors
            const entries = Array.from(ErrorSystem.activeErrors.values());
            if (entries.length === 0) {
                list.innerHTML = '<em style="color:var(--text-secondary);">No frontend errors</em>';
                return;
            }
            list.innerHTML = entries.map(e => `
                <div style="padding:4px 0;border-bottom:1px solid var(--border);">
                    <span style="color:${e.severity === 'critical' ? '#ef4444' : '#f59e0b'};font-weight:600;">${escapeHtml(e.code)}</span>
                    <span style="color:var(--text-secondary);margin-left:8px;">${escapeHtml(e.operation)}</span>
                </div>
            `).join('');
        }
    } catch {
        list.innerHTML = '<em style="color:var(--text-secondary);">Unable to load errors</em>';
    }
}


// ---------------------------------------------------------------------------
// Hook into settings panel open
// ---------------------------------------------------------------------------

const originalOpenSettings = window.openSettings || function(){};
window.openSettings = function() {
    originalOpenSettings();
    setTimeout(injectErrorSettings, 50);
};

// Also inject on DOM ready if settings panel exists
document.addEventListener('DOMContentLoaded', () => {
    if (document.getElementById('settings-panel')) {
        injectErrorSettings();
    }
});
