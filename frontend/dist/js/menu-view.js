/**
 * Menu — View Actions
 *
 * View menu handler functions (theme, UI scale, panel toggles,
 * fullscreen, accessibility modes) and DOMContentLoaded
 * theme/scale restoration.
 *
 * @module menu-view
 * @version 1.0.0
 */

'use strict';

function toggleColorblindMode() {
    document.body.classList.toggle('colorblind-mode');
    const enabled = document.body.classList.contains('colorblind-mode');
    if (window.settingsState?.set) {
        window.settingsState.set('accessibility.colorblindMode', enabled);
    }
    showNotification(`Colorblind mode ${enabled ? 'enabled' : 'disabled'}`, 'success');
}

function toggleHighContrastMode() {
    document.body.classList.toggle('high-contrast-mode');
    const enabled = document.body.classList.contains('high-contrast-mode');
    if (window.settingsState?.set) {
        window.settingsState.set('accessibility.highContrastMode', enabled);
    }
    showNotification(`High contrast mode ${enabled ? 'enabled' : 'disabled'}`, 'success');
}

function increaseUIScale() {
    const cur = parseFloat(getComputedStyle(document.documentElement).getPropertyValue('--ui-scale') || '1.0');
    const val = Math.min(cur + 0.1, 2.0);
    document.documentElement.style.setProperty('--ui-scale', val.toString());
    if (window.settingsState?.set) {
        window.settingsState.set('general.uiScale', val);
    }
    showNotification(`UI scale increased to ${(val * 100).toFixed(0)}%`, 'success');
}

function decreaseUIScale() {
    const cur = parseFloat(getComputedStyle(document.documentElement).getPropertyValue('--ui-scale') || '1.0');
    const val = Math.max(cur - 0.1, 0.5);
    document.documentElement.style.setProperty('--ui-scale', val.toString());
    if (window.settingsState?.set) {
        window.settingsState.set('general.uiScale', val);
    }
    showNotification(`UI scale decreased to ${(val * 100).toFixed(0)}%`, 'success');
}

function resetUIScale() {
    document.documentElement.style.setProperty('--ui-scale', '1.0');
    if (window.settingsState?.set) {
        window.settingsState.set('general.uiScale', 1.0);
    }
    showNotification('UI scale reset to 100%', 'success');
}

function toggleSidebar() {
    const el = document.querySelector('.nav-sidebar');
    if (el) {
        el.classList.toggle('hidden');
        const hidden = el.classList.contains('hidden');
        if (window.settingsState?.set) {
            window.settingsState.set('ui.sidebarHidden', hidden);
        }
        showNotification(`Sidebar ${hidden ? 'hidden' : 'shown'}`, 'success');
    }
}

function toggleCharacterPanel() {
    const el = document.querySelector('.character-sidebar');
    if (el) {
        el.classList.toggle('hidden');
        const hidden = el.classList.contains('hidden');
        if (window.settingsState?.set) {
            window.settingsState.set('ui.characterPanelHidden', hidden);
        }
        showNotification(`Character panel ${hidden ? 'hidden' : 'shown'}`, 'success');
    }
}

function toggleWorldPanel() {
    const el = document.querySelector('.world-sidebar');
    if (el) {
        el.classList.toggle('hidden');
        const hidden = el.classList.contains('hidden');
        if (window.settingsState?.set) {
            window.settingsState.set('ui.worldPanelHidden', hidden);
        }
        showNotification(`World panel ${hidden ? 'hidden' : 'shown'}`, 'success');
    }
}

function toggleFullscreen() {
    if (!document.fullscreenElement) {
        document.documentElement.requestFullscreen()
            .then(() => showNotification('Entered fullscreen mode', 'success'))
            .catch(() => showNotification('Failed to enter fullscreen mode', 'error'));
    } else {
        document.exitFullscreen()
            .then(() => showNotification('Exited fullscreen mode', 'success'))
            .catch(() => showNotification('Failed to exit fullscreen mode', 'error'));
    }
}

// ---------------------------------------------------------------------------
// Restore persisted theme + scale on load
// ---------------------------------------------------------------------------

document.addEventListener('DOMContentLoaded', function () {
    const s = window.settingsState?.current;

    // Theme is restored by theme.js - single owner. It used to be duplicated
    // here with a third storage key ('howl-chat-theme') and a hardcoded 'dark'
    // default, so on any page that does not load settings-state.js it overrode
    // the restored mode and wrote 'dark' back to the mirror.

    // UI Scale
    const savedScale = s?.general?.uiScale || localStorage.getItem('howl-chat-ui-scale');
    if (savedScale) document.documentElement.style.setProperty('--ui-scale', savedScale);

    // Animation Speed
    const savedAnimSpeed = s?.general?.animationSpeed || localStorage.getItem('howl-chat-animation-speed');
    if (savedAnimSpeed) document.documentElement.style.setProperty('--animation-speed', savedAnimSpeed);

    // Accessibility modes
    if (s?.accessibility?.colorblindMode) document.body.classList.add('colorblind-mode');
    if (s?.accessibility?.highContrastMode) document.body.classList.add('high-contrast-mode');

    // Panel visibility
    const sidebar = document.querySelector('.nav-sidebar');
    if (sidebar && s?.ui?.sidebarHidden) sidebar.classList.add('hidden');

    const charPanel = document.querySelector('.character-sidebar');
    if (charPanel && s?.ui?.characterPanelHidden) charPanel.classList.add('hidden');

    const worldPanel = document.querySelector('.world-sidebar');
    if (worldPanel && s?.ui?.worldPanelHidden) worldPanel.classList.add('hidden');
});
