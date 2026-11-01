/**
 * Settings Actions
 *
 * Apply helpers (theme, UI scale) and all settings action stubs
 * (storage paths, backup, cache, API keys, permissions).
 *
 * Depends on: settings-state.js
 *
 * @module settings-actions
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Apply Helpers
// ---------------------------------------------------------------------------

function applyTheme(theme) {
    // Themes have exactly one owner: theme.js. This used to reimplement the
    // body-class swap and write the settings store directly, which skipped the
    // localStorage mirror that the other eight pages restore from — so picking
    // a theme here changed this page and nothing else.
    if (typeof window.setThemeMode === 'function') {
        // setThemeMode only accepts a concrete mode and only persists when told
        // to, so 'system' has to be resolved here and the preference recorded
        // separately — otherwise Settings' choice was lost on the next load.
        const effective = theme === 'system'
            ? (window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark')
            : theme;
        window.setThemeMode(effective, false);
        if (window.settingsState?.set) {
            window.settingsState.set('general.theme', theme);
        }
        return;
    }
    document.body.classList.remove('theme-light', 'theme-dark');
    const effective = theme === 'system'
        ? (window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark')
        : theme;
    document.body.classList.add(`theme-${effective}`);
}


function applyAnimationSpeed(speed) {
    document.documentElement.style.setProperty('--animation-speed', speed.toString());
    if (window.settingsState?.set) {
        window.settingsState.set('general.animationSpeed', speed);
    }
}

// ---------------------------------------------------------------------------
// Backup Actions
// ---------------------------------------------------------------------------

function createBackup() {
    try {
        const data = {};
        for (let i = 0; i < localStorage.length; i++) {
            const key = localStorage.key(i);
            data[key] = localStorage.getItem(key);
        }
        const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
        const url  = URL.createObjectURL(blob);
        const a    = document.createElement('a');
        a.href     = url;
        a.download = `howl-backup-${new Date().toISOString().replace(/[:.]/g, '-')}.json`;
        a.click();
        URL.revokeObjectURL(url);
        if (typeof showNotification === 'function') showNotification('Backup created and downloaded', 'success');
    } catch (err) {
        console.error('Backup failed:', err);
        if (typeof showNotification === 'function') showNotification('Backup failed: ' + err.message, 'error');
    }
}

function restoreBackup() {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'application/json';
    input.onchange = async (e) => {
        const file = e.target.files[0];
        if (!file) return;
        try {
            const text = await file.text();
            const data = JSON.parse(text);
            if (!confirm('Restore backup? This will overwrite all current local data.')) return;
            localStorage.clear();
            Object.entries(data).forEach(([key, value]) => { localStorage.setItem(key, value); });
            if (typeof showNotification === 'function') showNotification('Backup restored. Reloading...', 'success');
            setTimeout(() => location.reload(), 1000);
        } catch (err) {
            console.error('Restore failed:', err);
            if (typeof showNotification === 'function') showNotification('Restore failed: ' + err.message, 'error');
        }
    };
    input.click();
}

function manageBackups() {
    const backups = [];
    for (let i = 0; i < localStorage.length; i++) {
        const key = localStorage.key(i);
        if (key && key.startsWith('howl-')) {
            const item = localStorage.getItem(key);
            backups.push(`${key}: ${(item.length / 1024).toFixed(1)} KB`);
        }
    }
    const list = backups.length ? backups.join('\n') : 'No HOWL data found.';
    alert(`Local Storage Summary:\n\n${list}\n\nTotal items: ${backups.length}`);
}

// ---------------------------------------------------------------------------
// Data Actions
// ---------------------------------------------------------------------------

function clearCache() {
    if (!confirm('Clear cache? This will remove all temporary files and cached data.')) return;
    try {
        sessionStorage.clear();
        const keysToRemove = [];
        for (let i = 0; i < localStorage.length; i++) {
            const key = localStorage.key(i);
            if (key && (key.includes('cache') || key.includes('temp') || key.includes('preview'))) {
                keysToRemove.push(key);
            }
        }
        keysToRemove.forEach(key => localStorage.removeItem(key));
        if (typeof showNotification === 'function') showNotification('Cache cleared successfully', 'success');
    } catch (err) {
        console.error('Clear cache failed:', err);
    }
}

function clearAllData() {
    // Chats, worlds, characters, scenarios and lorebooks live as JSON files
    // under AppData and are removed by the backend. Clearing localStorage here
    // left every one of them on disk while telling the user everything was gone.
    // resetApplication() is the real wipe: backend WipeAllData, IndexedDB, then
    // localStorage, with its own two-step confirmation.
    if (typeof window.resetApplication === 'function') {
        return window.resetApplication();
    }
    if (!confirm('Clear all local data? This will permanently delete ALL application data including chats, settings, and cached files. This action cannot be undone.')) return;
    if (!confirm('Are you absolutely sure? This cannot be undone.')) return;
    try {
        localStorage.clear();
        sessionStorage.clear();
        if (typeof showNotification === 'function') showNotification('All local data cleared. Reloading...', 'success');
        setTimeout(() => location.reload(), 1000);
    } catch (err) {
        console.error('Clear all data failed:', err);
    }
}

// ---------------------------------------------------------------------------
// Privacy / Local-Only Mode
// ---------------------------------------------------------------------------

function applyPrivacySettings() {
    const s = window.settingsState.current.privacy;

    // Local-only mode: block external network requests from the frontend
    if (s.localOnlyMode) {
        if (!window._originalFetch) window._originalFetch = window.fetch;
        window.fetch = function(...args) {
            const url = typeof args[0] === 'string' ? args[0] : args[0]?.url || '';
            const isLocal = !url || url.startsWith('file:') || url.startsWith('data:') || url.startsWith('blob:') || url.startsWith('http://localhost') || url.startsWith('http://127.0.0.1') || url.startsWith('https://localhost') || url.startsWith('https://127.0.0.1');
            if (!isLocal) {
                console.warn('[Local-Only] Blocked fetch:', url);
                return Promise.reject(new Error('Local-only mode: external requests are blocked'));
            }
            return window._originalFetch.apply(this, args);
        };
        if (!window._originalXHROpen) window._originalXHROpen = XMLHttpRequest.prototype.open;
        XMLHttpRequest.prototype.open = function(method, url, ...rest) {
            const isLocal = !url || url.startsWith('file:') || url.startsWith('data:') || url.startsWith('blob:') || url.startsWith('http://localhost') || url.startsWith('http://127.0.0.1') || url.startsWith('https://localhost') || url.startsWith('https://127.0.0.1');
            if (!isLocal) {
                console.warn('[Local-Only] Blocked XHR:', url);
                throw new Error('Local-only mode: external requests are blocked');
            }
            return window._originalXHROpen.call(this, method, url, ...rest);
        };
    } else {
        if (window._originalFetch) window.fetch = window._originalFetch;
        if (window._originalXHROpen) XMLHttpRequest.prototype.open = window._originalXHROpen;
    }

    // Disable telemetry: set a flag that backend and error handlers can check
    window.__telemetryDisabled = s.disableTelemetry === true;
}
