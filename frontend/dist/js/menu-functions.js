/**
 * Menu Functions — Orchestrator
 *
 * Settings launchers, keyboard shortcuts dialog, audio stubs,
 * the openStoryTools override that loads recent memories,
 * the showNotification utility, and window.menuFunctions export.
 *
 * Depends on (must load first):
 *   menu-edit.js   — edit menu actions
 *   menu-view.js   — view menu actions + DOMContentLoaded theme restore
 *   menu-memory.js — memory list/edit flyout helpers
 *
 * @module menu-functions
 * @version 2.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Settings Launchers
// ---------------------------------------------------------------------------

function openStoryTools() {
    toggleFlyout('story-tools-flyout');
    loadRecentMemoriesMenu();
}

let _audioSettingsEscHandler = null;

function _removeAudioSettingsModal(modal) {
    const iframeEl = modal?.querySelector('iframe');
    if (iframeEl && window.howlBridgeHost?.detach) {
        try { window.howlBridgeHost.detach(iframeEl); } catch (e) { /* bridge already gone */ }
    }
    if (modal?.parentNode) modal.remove();
    if (_audioSettingsEscHandler) {
        document.removeEventListener('keydown', _audioSettingsEscHandler);
        _audioSettingsEscHandler = null;
    }
}

function _requestAudioSettingsClose(modal) {
    try {
        const iframe = modal?.querySelector('iframe');
        const close = iframe?.contentWindow?.AudioSettings?.close;
        if (typeof close === 'function') {
            close();
            return;
        }
    } catch (e) { /* fall back to removing the modal */ }
    _removeAudioSettingsModal(modal);
}

function _openAudioSettingsIframe(hash) {
    // Remove existing modal if already open
    const existing = document.getElementById('audio-settings-modal');
    if (existing) _removeAudioSettingsModal(existing);

    const modal = document.createElement('div');
    modal.className = 'real-settings-overlay show';
    modal.id        = 'audio-settings-modal';
    modal.style.zIndex = '2500';
    // rs-container.css turns .real-settings-container into a full-window workspace
    // (height:100vh, no border) for the main settings dialog. This panel is a
    // floating dialog, so every one of those properties is overridden here:
    // otherwise the container fills the window and the area below the 600px
    // iframe renders as an empty grey band.
    modal.innerHTML = `
        <div class="real-settings-container" style="width:900px;max-width:95vw;height:auto;max-height:85vh;border:1px solid var(--border-primary);border-radius:12px;overflow:hidden;background:var(--shell-panel-bg);box-shadow:0 0 30px rgba(0,0,0,0.5);">
            <iframe src="audio-settings.html#${hash}" style="display:block;width:100%;height:min(600px,78vh);border:none;border-radius:8px;" id="audio-settings-iframe"></iframe>
        </div>
    `;
    document.body.appendChild(modal);
    const iframeEl = modal.querySelector('iframe');
    if (iframeEl && window.howlBridgeHost?.attach) {
        try { window.howlBridgeHost.attach(iframeEl); } catch (e) { console.error('[audio-settings] bridge attach failed:', e); }
    }
    modal.addEventListener('click', e => {
        if (e.target === modal) _requestAudioSettingsClose(modal);
    });
    if (_audioSettingsEscHandler) document.removeEventListener('keydown', _audioSettingsEscHandler);
    _audioSettingsEscHandler = e => {
        if (e.key === 'Escape') _requestAudioSettingsClose(modal);
    };
    document.addEventListener('keydown', _audioSettingsEscHandler);
}

// The iframe closes itself after saving/applying. Remove the matching parent
// key handler and modal without adding a second close path.
document.addEventListener('howl:audio-settings-closed', () => {
    const modal = document.getElementById('audio-settings-modal');
    if (modal) _removeAudioSettingsModal(modal);
});

function openAudioSettings() { _openAudioSettingsIframe('general'); }
function openTTSSettings()   { _openAudioSettingsIframe('tts'); }

// openRealSettings() is defined in real-settings.js (loaded earlier)

// ---------------------------------------------------------------------------
// About Modal
// ---------------------------------------------------------------------------

function showAbout() {
    const overlay = document.getElementById('about-overlay');
    if (overlay) {
        overlay.classList.add('show');
    }
    // Close all dropdown menus
    document.querySelectorAll('.dropdown').forEach(menu => menu.classList.remove('show'));
}

function hideAbout() {
    const overlay = document.getElementById('about-overlay');
    if (overlay) {
        overlay.classList.remove('show');
    }
}

// Close About on Escape
document.addEventListener('keydown', function(event) {
    if (event.key === 'Escape') {
        hideAbout();
    }
});

// ---------------------------------------------------------------------------
// Keyboard Shortcuts
// ---------------------------------------------------------------------------

function showKeyboardShortcuts() {
    alert(`Keyboard Shortcuts:

General:
• Ctrl+N - New Chat       • Ctrl+S - Save Chat
• Ctrl+Z - Undo           • Ctrl+Y - Redo
• Ctrl+C - Copy           • Ctrl+V - Paste
• Ctrl+A - Select All     • Ctrl+F - Find
• F11 - Fullscreen        • Esc - Close dialogs

Chat:
• Enter - Send            • Shift+Enter - New line
• Ctrl+Enter - Continue   • Alt+R - Regenerate

Navigation:
• Ctrl+1 to Ctrl+6 - Switch pages`);
}

// ---------------------------------------------------------------------------
// Shared Notification (used by menu-edit, menu-view, menu-memory)
// ---------------------------------------------------------------------------

function showNotification(message, type = 'info', duration = 3000) {
    let n = document.getElementById('notification');
    if (!n) {
        n = document.createElement('div');
        n.id = 'notification';
        n.style.cssText = 'position:fixed;top:20px;right:20px;padding:12px 20px;border-radius:8px;color:white;font-weight:500;z-index:10000;opacity:0;transform:translateX(100%);transition:all 0.3s ease;max-width:300px;word-wrap:break-word;';
        document.body.appendChild(n);
    }
    const colors = { success: '#28a745', error: '#dc3545', warning: '#ffc107', info: '#17a2b8' };
    n.style.backgroundColor = colors[type] || colors.info;
    n.textContent           = message;
    n.style.opacity         = '1';
    n.style.transform       = 'translateX(0)';
    setTimeout(() => { n.style.opacity = '0'; n.style.transform = 'translateX(100%)'; }, duration);
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

window.menuFunctions = {
    editUndo, editRedo, editCut, editCopy, editPaste, editSelectAll, editFind, editReplace,
    toggleColorblindMode, toggleHighContrastMode, setTheme,
    increaseUIScale, decreaseUIScale, resetUIScale,
    toggleSidebar, toggleCharacterPanel, toggleWorldPanel, toggleFullscreen,
    openStoryTools, showKeyboardShortcuts,
    openAudioSettings, openTTSSettings
};
