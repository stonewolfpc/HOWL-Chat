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
    document.body.className = document.body.className.replace(/theme-\w+/g, '');
    let effectiveTheme = theme;
    if (theme === 'system') {
        effectiveTheme = window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark';
    }
    document.body.classList.add(`theme-${effectiveTheme}`);
    localStorage.setItem('howl-chat-theme', theme);
}

function applyUIScale(scale) {
    document.documentElement.style.setProperty('--ui-scale', scale.toString());
    localStorage.setItem('howl-chat-ui-scale', scale.toString());
}

function applyAnimationSpeed(speed) {
    document.documentElement.style.setProperty('--animation-speed', speed.toString());
    localStorage.setItem('howl-chat-animation-speed', speed.toString());
}

// ---------------------------------------------------------------------------
// Storage Path Actions
// ---------------------------------------------------------------------------

function _browseDirectory(inputId, settingsKey) {
    const el = document.getElementById(inputId);
    if (!el) return;

    const input = document.createElement('input');
    input.type = 'file';
    input.webkitdirectory = true;
    input.style.display = 'none';
    document.body.appendChild(input);

    input.addEventListener('change', () => {
        if (input.files && input.files.length > 0) {
            const first = input.files[0];
            const path = first.path || first.webkitRelativePath.split('/')[0] || first.name;
            el.value = path;
            window.settingsState.current.storage[settingsKey] = path;
        }
        document.body.removeChild(input);
    });
    input.click();
}

function _resetLocation(inputId, settingsKey, defaultValue) {
    window.settingsState.current.storage[settingsKey] = defaultValue;
    const el = document.getElementById(inputId);
    if (el) el.value = defaultValue;
    window.settingsState.saveSettings();
}

function browseChatLocation()     { _browseDirectory('chat-storage-location',     'chatStorageLocation'); }
function resetChatLocation()      { _resetLocation('chat-storage-location',       'chatStorageLocation', 'AppData/Chats'); }
function browseContentLocation()  { _browseDirectory('content-storage-location',  'contentStorageLocation'); }
function resetContentLocation()   { _resetLocation('content-storage-location',  'contentStorageLocation', 'AppData'); }
function browseWorldLocation()    { _browseDirectory('world-storage-location',    'worldStorageLocation'); }
function resetWorldLocation()     { _resetLocation('world-storage-location',    'worldStorageLocation', 'AppData/Worlds'); }
function browseScenarioLocation() { _browseDirectory('scenario-storage-location', 'scenarioStorageLocation'); }
function resetScenarioLocation()  { _resetLocation('scenario-storage-location',   'scenarioStorageLocation', 'AppData/Scenarios'); }
function browseCharacterLocation(){ _browseDirectory('character-storage-location','characterStorageLocation'); }
function resetCharacterLocation() { _resetLocation('character-storage-location',  'characterStorageLocation', 'AppData/Characters'); }
function browseImageLocation()    { _browseDirectory('image-storage-location',    'imageStorageLocation'); }
function resetImageLocation()     { _resetLocation('image-storage-location',      'imageStorageLocation', 'AppData/Images'); }

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
// Permission Actions
// ---------------------------------------------------------------------------

function managePermissions() {
    const s = window.settingsState.current;
    const perms = s.privacy.permissions || {};
    const permList = Object.entries(perms).map(([k, v]) => `${k}: ${v ? 'allowed' : 'denied'}`).join('\n') || 'No custom permissions set.';
    const action = prompt(`Current Permissions:\n\n${permList}\n\nEnter action: (add, remove, clear, or cancel)`);
    if (!action) return;
    const act = action.trim().toLowerCase();
    if (act === 'add') {
        const name = prompt('Enter permission name:');
        if (name) {
            perms[name] = true;
            s.privacy.permissions = perms;
            if (typeof showNotification === 'function') showNotification(`Permission "${name}" allowed`, 'success');
        }
    } else if (act === 'remove') {
        const name = prompt(`Enter permission to deny:\n${Object.keys(perms).join(', ')}`);
        if (name && perms.hasOwnProperty(name)) {
            perms[name] = false;
            s.privacy.permissions = perms;
            if (typeof showNotification === 'function') showNotification(`Permission "${name}" denied`, 'success');
        }
    } else if (act === 'clear') {
        s.privacy.permissions = {};
        if (typeof showNotification === 'function') showNotification('All permissions reset', 'success');
    }
}

function resetPermissions() {
    window.settingsState.current.privacy.permissions = {};
    if (typeof showNotification === 'function') showNotification('Permissions reset to defaults', 'success');
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
