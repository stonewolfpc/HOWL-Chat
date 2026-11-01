/**
 * Theme Management Module
 * 
 * This module handles theme cycling through the Artherius color palette
 * with smooth transitions and animations, plus light/dark theme switching.
 * 
 * @module theme
 */

// Theme definitions for Artherius color palette
const themes = [
    {
        color: '#14b8a6',
        glow: 'rgba(20, 184, 166, 0.5)',
        border: 'rgba(20, 184, 166, 0.3)',
        textGlow: 'rgba(20, 184, 166, 0.4)',
        secondary: '#06b6d4'
    },
    {
        color: '#06b6d4',
        glow: 'rgba(6, 182, 212, 0.5)',
        border: 'rgba(6, 182, 212, 0.3)',
        textGlow: 'rgba(6, 182, 212, 0.4)',
        secondary: '#14b8a6'
    },
    {
        color: '#a855f7',
        glow: 'rgba(168, 85, 247, 0.5)',
        border: 'rgba(168, 85, 247, 0.3)',
        textGlow: 'rgba(168, 85, 247, 0.4)',
        secondary: '#ec4899'
    }
];

let currentThemeIndex = 2; // Start with cyan
const targetFPS = 30; // Target 30fps for smooth animation
const frameInterval = 1000 / targetFPS; // ~33.33ms per frame
let animationStartTime = 0;
let isAnimating = false;

// Light/Dark mode state
const THEME_MODE_KEY = 'howl-theme-mode';
let currentThemeMode = 'dark'; // 'dark' or 'light'

// Easing function for smooth transitions
function easeInOutCubic(t) {
    return t < 0.5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2;
}

// Smoother easing function for borders
function easeInOutQuad(t) {
    return t < 0.5 ? 2 * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2;
}

// Fade color interpolation
function fadeColor(color1, color2, progress) {
    const r1 = parseInt(color1.slice(1, 3), 16);
    const g1 = parseInt(color1.slice(3, 5), 16);
    const b1 = parseInt(color1.slice(5, 7), 16);
    
    const r2 = parseInt(color2.slice(1, 3), 16);
    const g2 = parseInt(color2.slice(3, 5), 16);
    const b2 = parseInt(color2.slice(5, 7), 16);
    
    const easedProgress = easeInOutCubic(progress);
    
    const r = Math.round(r1 + (r2 - r1) * easedProgress);
    const g = Math.round(g1 + (g2 - g1) * easedProgress);
    const b = Math.round(b1 + (b2 - b1) * easedProgress);
    
    return `#${r.toString(16).padStart(2, '0')}${g.toString(16).padStart(2, '0')}${b.toString(16).padStart(2, '0')}`;
}

// Fade RGBA interpolation
function fadeRGBA(rgba1, rgba2, progress) {
    // Parse RGBA values using regex to handle variable digit lengths
    const parseRGBA = (rgba) => {
        const match = rgba.match(/rgba?\((\d+),\s*(\d+),\s*(\d+),\s*([\d.]+)\)/);
        if (!match) return { r: 0, g: 0, b: 0, a: 0 };
        return {
            r: parseInt(match[1], 10),
            g: parseInt(match[2], 10),
            b: parseInt(match[3], 10),
            a: parseFloat(match[4])
        };
    };
    
    const color1 = parseRGBA(rgba1);
    const color2 = parseRGBA(rgba2);
    
    const easedProgress = easeInOutCubic(progress);
    
    const r = Math.round(color1.r + (color2.r - color1.r) * easedProgress);
    const g = Math.round(color1.g + (color2.g - color1.g) * easedProgress);
    const b = Math.round(color1.b + (color2.b - color1.b) * easedProgress);
    const a = (color1.a + (color2.a - color1.a) * easedProgress).toFixed(2);
    
    return `rgba(${r}, ${g}, ${b}, ${a})`;
}

// Background color for fading to
const backgroundColor = '#0a0a0f';

// Theme cycling function
function cycleTheme() {
    if (isAnimating) return;
    
    isAnimating = true;
    animationStartTime = performance.now();
    
    const currentTheme = themes[currentThemeIndex];
    const nextThemeIndex = (currentThemeIndex + 1) % themes.length;
    const nextTheme = themes[nextThemeIndex];
    
    const fadeDuration = 4000; // 4 seconds for fade
    const pauseDuration = 1000; // 1 second pause
    
    function animate(timestamp) {
        const elapsed = timestamp - animationStartTime;
        const progress = Math.min(elapsed / fadeDuration, 1);
        
        // Fade to background first
        if (progress < 0.5) {
            const fadeProgress = progress * 2;
            document.documentElement.style.setProperty('--theme-color', fadeColor(currentTheme.color, backgroundColor, fadeProgress));
            document.documentElement.style.setProperty('--theme-glow', fadeRGBA(currentTheme.glow, 'rgba(15, 21, 31, 0.5)', fadeProgress));
            // Keep border static - don't change to prevent pulsing
            document.documentElement.style.setProperty('--theme-text-glow', fadeRGBA(currentTheme.textGlow, 'rgba(100, 100, 100, 0.4)', fadeProgress));
        } else {
            // Fade to new theme
            const fadeProgress = (progress - 0.5) * 2;
            document.documentElement.style.setProperty('--theme-color', fadeColor(backgroundColor, nextTheme.color, fadeProgress));
            document.documentElement.style.setProperty('--theme-glow', fadeRGBA('rgba(15, 21, 31, 0.5)', nextTheme.glow, fadeProgress));
            // Keep border static - don't change to prevent pulsing
            document.documentElement.style.setProperty('--theme-text-glow', fadeRGBA('rgba(100, 100, 100, 0.4)', nextTheme.textGlow, fadeProgress));
        }

        if (progress < 1) {
            requestAnimationFrame(animate);
        } else {
            // Set final theme (keep border static)
            document.documentElement.style.setProperty('--theme-color', nextTheme.color);
            document.documentElement.style.setProperty('--theme-glow', nextTheme.glow);
            // Don't change border - keep it static to prevent pulsing
            document.documentElement.style.setProperty('--theme-text-glow', nextTheme.textGlow);
            document.documentElement.style.setProperty('--theme-secondary', nextTheme.secondary);

            currentThemeIndex = nextThemeIndex;
            isAnimating = false;

            // Schedule next cycle
            setTimeout(cycleTheme, pauseDuration);
        }
    }
    
    requestAnimationFrame(animate);
}

// ---------------------------------------------------------------------------
// Light/Dark Theme Switching
// ---------------------------------------------------------------------------

/**
 * Switch between light and dark theme modes
 * @param {string} mode - 'light' or 'dark'
 */
function setThemeMode(mode, persistStore) {
    if (mode !== 'light' && mode !== 'dark') {
        console.error('Invalid theme mode. Use "light" or "dark".');
        return;
    }
    
    currentThemeMode = mode;
    const body = document.body;
    
    if (mode === 'light') {
        body.classList.remove('theme-dark');
        body.classList.add('theme-light');
        // Light mode uses reduced glow
        document.documentElement.style.setProperty('--theme-glow', `rgba(20, 184, 166, 0.25)`);
        document.documentElement.style.setProperty('--theme-text-glow', `rgba(20, 184, 166, 0.2)`);
    } else {
        body.classList.remove('theme-light');
        body.classList.add('theme-dark');
        // Dark mode uses full glow
        const currentTheme = themes[currentThemeIndex];
        document.documentElement.style.setProperty('--theme-glow', currentTheme.glow);
        document.documentElement.style.setProperty('--theme-text-glow', currentTheme.textGlow);
    }
    
    // Save to the unified settings store, and mirror it to localStorage.
    // Only three of the nine pages load settings-state.js, so without the
    // mirror every other page had nothing to restore from and stayed dark.
    try {
        localStorage.setItem(THEME_MODE_KEY, mode);
    } catch (e) {
        // private mode / disabled storage: the in-memory class still applies
    }
    if (persistStore && window.settingsState?.set) {
        window.settingsState.set('general.theme', mode);
    }
}

/**
 * Toggle between light and dark theme
 */
function toggleThemeMode() {
    const newMode = currentThemeMode === 'dark' ? 'light' : 'dark';
    setThemeMode(newMode);
}

/**
 * Load saved theme preference from unified settings store
 */
/**
 * Resolve a stored preference into a concrete mode.
 * 'system' is decided at load time on every page.
 */
function resolveThemeMode(mode) {
    if (mode === 'system') {
        return window.matchMedia && window.matchMedia('(prefers-color-scheme: light)').matches
            ? 'light' : 'dark';
    }
    return (mode === 'light' || mode === 'dark') ? mode : null;
}

/**
 * Restore the theme on every page.
 *
 * The unified settings store is authoritative when it has an answer, but only
 * index, audio-settings and real-settings load it. localStorage is the durable
 * mirror every page can read synchronously, so it carries the preference to the
 * rest of the app. Restoring never writes back, so 'system' survives.
 */
function loadThemePreference() {
    let fromStore = null;
    if (window.settingsState && window.settingsState.current && window.settingsState.current.general) {
        fromStore = window.settingsState.current.general.theme;
    }

    let stored = null;
    try {
        stored = localStorage.getItem(THEME_MODE_KEY);
    } catch (e) {
        stored = null;
    }

    const chosen = resolveThemeMode(fromStore) || resolveThemeMode(stored);
    if (chosen) {
        setThemeMode(chosen, false);
        return;
    }

    // Nothing saved anywhere yet: keep the default, but record it so the first
    // navigation to another page finds a consistent answer.
    try {
        localStorage.setItem(THEME_MODE_KEY, currentThemeMode);
    } catch (e) {
        // ignore
    }
}

// Initialize theme preference as early as this script is reachable. On most
// pages theme.js sits near the end of the body, so this runs after the markup
// is parsed but before first paint of the theme transition.
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', loadThemePreference);
} else {
    loadThemePreference();
}

// Follow the OS while the preference is 'system'. This listener used to sit in
// menu-view.js and referenced a variable that no longer existed, so switching
// the OS theme never reached the app. theme.js owns the preference, so it owns
// the reaction to it too.
if (window.matchMedia) {
    window.matchMedia('(prefers-color-scheme: light)').addEventListener('change', function () {
        if (currentThemeMode !== 'system') return;
        const resolved = resolveThemeMode('system');
        if (resolved) setThemeMode(resolved, false);
    });
}

// ---------------------------------------------------------------------------
// Legacy Function Support for Menu Items
// ---------------------------------------------------------------------------

/**
 * Legacy function for menu compatibility
 * Maps to setThemeMode for existing menu items
 * @param {string} theme - 'light' or 'dark'
 */
function setTheme(theme) {
    // Single owner. menu-view.js used to define its own setTheme, which
    // shadowed this one on six pages; it resolved 'system' and passed the
    // result to setThemeMode. That resolution lives here now, because
    // setThemeMode only ever accepts light or dark.
    const mode = theme === 'system'
        ? (window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark')
        : theme;
    setThemeMode(mode);
    // The menu hands us the user's choice, not the resolved mode: store
    // 'system' as 'system' so Settings shows what was picked, and so the choice
    // survives a restart. setThemeMode only persists when explicitly asked.
    if (window.settingsState?.set) {
        window.settingsState.set('general.theme', theme);
        if (window.settingsState.saveSettings) window.settingsState.saveSettings();
    }
}

// Start theme cycling after initial delay (disabled to prevent flickering)
// setTimeout(cycleTheme, 2000);
