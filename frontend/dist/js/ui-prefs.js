/**
 * UI Preferences — renderer-only, deliberately NOT settings
 *
 * This namespace is for ephemeral presentation state: which settings panel was
 * open last, expanded/collapsed sections, and similar view preferences. These
 * values never belong in the backend settings document and must never be
 * written by settings-state.js.
 *
 * Domain settings (theme, audio, sampler, scale, paths, ...) live in
 * settings-state.js, whose backend is authoritative. If a value changes
 * behavior that must survive a reinstall or follow the user across pages, it
 * belongs there, not here.
 *
 * @module ui-prefs
 * @version 1.0.0
 */

'use strict';

const UI_PREF_PREFIX = 'howl-ui:';

function _uiPrefKey(key) {
    return String(key).startsWith(UI_PREF_PREFIX) ? String(key) : UI_PREF_PREFIX + key;
}

window.UIPrefs = {
    prefix: UI_PREF_PREFIX,

    get(key, fallback = null) {
        try {
            const raw = localStorage.getItem(_uiPrefKey(key));
            return raw === null ? fallback : JSON.parse(raw);
        } catch (e) {
            console.warn('[UIPrefs] read failed:', key, e);
            return fallback;
        }
    },

    set(key, value) {
        try {
            localStorage.setItem(_uiPrefKey(key), JSON.stringify(value));
            return true;
        } catch (e) {
            console.warn('[UIPrefs] write failed:', key, e);
            return false;
        }
    },

    remove(key) {
        try {
            localStorage.removeItem(_uiPrefKey(key));
        } catch (e) {
            console.warn('[UIPrefs] remove failed:', key, e);
        }
    },

    isUIPrefKey(key) {
        return typeof key === 'string' && key.startsWith(UI_PREF_PREFIX);
    }
};
