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
    showNotification(`Colorblind mode ${document.body.classList.contains('colorblind-mode') ? 'enabled' : 'disabled'}`, 'success');
}

function toggleHighContrastMode() {
    document.body.classList.toggle('high-contrast-mode');
    showNotification(`High contrast mode ${document.body.classList.contains('high-contrast-mode') ? 'enabled' : 'disabled'}`, 'success');
}

function setTheme(theme) {
    document.body.className = document.body.className.replace(/theme-\w+/g, '');
    let effective = theme;
    if (theme === 'system') {
        effective = window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark';
    }
    document.body.classList.add(`theme-${effective}`);
    localStorage.setItem('howl-chat-theme', theme);
    showNotification(`${effective.charAt(0).toUpperCase() + effective.slice(1)} theme applied`, 'success');
}

function increaseUIScale() {
    const cur = parseFloat(getComputedStyle(document.documentElement).getPropertyValue('--ui-scale') || '1.0');
    const val = Math.min(cur + 0.1, 2.0);
    document.documentElement.style.setProperty('--ui-scale', val.toString());
    localStorage.setItem('howl-chat-ui-scale', val.toString());
    showNotification(`UI scale increased to ${(val * 100).toFixed(0)}%`, 'success');
}

function decreaseUIScale() {
    const cur = parseFloat(getComputedStyle(document.documentElement).getPropertyValue('--ui-scale') || '1.0');
    const val = Math.max(cur - 0.1, 0.5);
    document.documentElement.style.setProperty('--ui-scale', val.toString());
    localStorage.setItem('howl-chat-ui-scale', val.toString());
    showNotification(`UI scale decreased to ${(val * 100).toFixed(0)}%`, 'success');
}

function resetUIScale() {
    document.documentElement.style.setProperty('--ui-scale', '1.0');
    localStorage.setItem('howl-chat-ui-scale', '1.0');
    showNotification('UI scale reset to 100%', 'success');
}

function toggleSidebar() {
    const el = document.querySelector('.nav-sidebar');
    if (el) { el.classList.toggle('hidden'); showNotification(`Sidebar ${el.classList.contains('hidden') ? 'hidden' : 'shown'}`, 'success'); }
}

function toggleCharacterPanel() {
    const el = document.querySelector('.character-sidebar');
    if (el) { el.classList.toggle('hidden'); showNotification(`Character panel ${el.classList.contains('hidden') ? 'hidden' : 'shown'}`, 'success'); }
}

function toggleWorldPanel() {
    const el = document.querySelector('.world-sidebar');
    if (el) { el.classList.toggle('hidden'); showNotification(`World panel ${el.classList.contains('hidden') ? 'hidden' : 'shown'}`, 'success'); }
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
    const settingsTheme = window.settingsState?.current?.general?.theme;
    const theme = settingsTheme || localStorage.getItem('howl-chat-theme') || 'dark';
    setTheme(theme);

    const savedScale = localStorage.getItem('howl-chat-ui-scale');
    if (savedScale) document.documentElement.style.setProperty('--ui-scale', savedScale);

    const savedAnimSpeed = localStorage.getItem('howl-chat-animation-speed');
    if (savedAnimSpeed) document.documentElement.style.setProperty('--animation-speed', savedAnimSpeed);

    if (theme === 'system') {
        window.matchMedia('(prefers-color-scheme: light)').addEventListener('change', e => {
            document.body.className = document.body.className.replace(/theme-\w+/g, '');
            document.body.classList.add(e.matches ? 'theme-light' : 'theme-dark');
        });
    }
});
