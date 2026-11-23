/**
 * Sampler Settings Module
 *
 * Handles synchronization between sliders/number inputs/toggles for sampler
 * settings. All reads and writes go through settings-state.js, the renderer's
 * single seam to the backend settings authority.
 *
 * @module sampler-settings
 */

const SAMPLER_PRESETS_KEY = 'howl-sampler-presets';

// Built-in sampler presets
const SAMPLER_PRESETS = {
    balanced: {
        temperature: 0.7, top_p: 0.9, top_k: 40, min_p: 0.05,
        typical_p: 1.0, repeat_penalty: 1.1, repeat_last_n: 64,
        presence_penalty: 0.0, frequency_penalty: 0.0,        mirostat_mode: 0, mirostat_tau: 5.0, mirostat_eta: 0.1,
        dyn_temp_range: 0.0, dyn_temp_exp: 1.0, dry_multiplier: 0.0,
        dry_allowed_length: 2, dry_base: 1.0,        smoothing_curve: 1.0, top_a: 0.0,        eta_cutoff: 0.0, encoder_repeat_penalty: 1.0, seed: -1
    },
    creative: {
        temperature: 0.9, top_p: 0.95, top_k: 80, min_p: 0.02,
        typical_p: 1.0, repeat_penalty: 1.05, repeat_last_n: 32,
        presence_penalty: 0.0, frequency_penalty: 0.0,        mirostat_mode: 0, mirostat_tau: 5.0, mirostat_eta: 0.1,
        dyn_temp_range: 0.5, dyn_temp_exp: 1.0, dry_multiplier: 0.0,
        dry_allowed_length: 2, dry_base: 1.0,        smoothing_curve: 1.0, top_a: 0.0,        eta_cutoff: 0.0, encoder_repeat_penalty: 1.0, seed: -1
    },
    deterministic: {
        temperature: 0.2, top_p: 0.5, top_k: 10, min_p: 0.1,
        typical_p: 1.0, repeat_penalty: 1.2, repeat_last_n: 128,
        presence_penalty: 0.0, frequency_penalty: 0.0,        mirostat_mode: 0, mirostat_tau: 5.0, mirostat_eta: 0.1,
        dyn_temp_range: 0.0, dyn_temp_exp: 1.0, dry_multiplier: 0.0,
        dry_allowed_length: 2, dry_base: 1.0,        smoothing_curve: 1.0, top_a: 0.0,        eta_cutoff: 0.0, encoder_repeat_penalty: 1.0, seed: -1
    },
    precise: {
        temperature: 0.2, top_p: 0.5, top_k: 10, min_p: 0.1,
        typical_p: 1.0, repeat_penalty: 1.2, repeat_last_n: 128,
        presence_penalty: 0.0, frequency_penalty: 0.0,        mirostat_mode: 0, mirostat_tau: 5.0, mirostat_eta: 0.1,
        dyn_temp_range: 0.0, dyn_temp_exp: 1.0, dry_multiplier: 0.0,
        dry_allowed_length: 2, dry_base: 1.0,        smoothing_curve: 1.0, top_a: 0.0,        eta_cutoff: 0.0, encoder_repeat_penalty: 1.0, seed: -1
    }
};

// Sampler settings configuration
const samplerSettings = [
    { id: 'temperature', default: 0.7 },
    { id: 'top-p', default: 0.9, hasToggle: true },
    { id: 'top-k', default: 40 },
    { id: 'min-p', default: 0.05 },
    { id: 'typical-p', default: 1.0, hasToggle: true },
    { id: 'repeat-penalty', default: 1.1 },
    { id: 'repeat-last-n', default: 64 },
    { id: 'presence-penalty', default: 0.0, hasToggle: true },
    { id: 'frequency-penalty', default: 0.0, hasToggle: true },
    { id: 'mirostat-mode', default: 0, hasToggle: true },
    { id: 'mirostat-tau', default: 5.0 },
    { id: 'mirostat-eta', default: 0.1 },
    { id: 'dyn-temp-range', default: 0.0, hasToggle: true },
    { id: 'dyn-temp-exp', default: 1.0 },
    { id: 'dry-multiplier', default: 0.0 },
    { id: 'dry-allowed-length', default: 2 },
    { id: 'dry-base', default: 1.0 },
    { id: 'seed', default: -1 }
];

// Initialize slider-value synchronization
function initSamplerSettings() {
    samplerSettings.forEach(setting => {
        const slider = document.getElementById(setting.id);
        const valueInput = document.getElementById(`${setting.id}-value`);
        const toggle = document.getElementById(`${setting.id}-enabled`);
        const card = document.getElementById(`${setting.id}-card`);
        
        if (slider && valueInput) {
            // Slider change updates number input
            slider.addEventListener('input', function() {
                valueInput.value = this.value;
                saveSamplerSettings();
            });
            
            // Number input change updates slider
            valueInput.addEventListener('input', function() {
                const min = parseFloat(slider.min);
                const max = parseFloat(slider.max);
                let value = parseFloat(this.value);
                
                // Clamp value to min/max range
                if (value < min) value = min;
                if (value > max) value = max;
                
                slider.value = value;
                saveSamplerSettings();
            });
            
            // Number input blur validates value
            valueInput.addEventListener('blur', function() {
                const min = parseFloat(slider.min);
                const max = parseFloat(slider.max);
                let value = parseFloat(this.value);
                
                if (isNaN(value)) {
                    this.value = setting.default;
                    slider.value = setting.default;
                    saveSamplerSettings();
                } else if (value < min) {
                    this.value = min;
                    slider.value = min;
                    saveSamplerSettings();
                } else if (value > max) {
                    this.value = max;
                    slider.value = max;
                    saveSamplerSettings();
                }
            });
        }
        
        // Toggle change
        if (toggle && card) {
            toggle.addEventListener('change', function() {
                if (this.checked) {
                    card.classList.remove('disabled');
                    handleConflictDetection(setting.id, true);
                } else {
                    card.classList.add('disabled');
                    handleConflictDetection(setting.id, false);
                }
                saveSamplerSettings();
            });
        }
    });
    
    // Initialize default reset buttons
    initResetButtons();
    
    // Load settings once the authoritative snapshot is in memory.
    (window.settingsState?.ready || Promise.resolve(true)).then(loadSamplerSettings);
}

// Initialize default reset button functionality
function initResetButtons() {
    const resetButtons = document.querySelectorAll('.setting-reset');
    resetButtons.forEach(button => {
        button.addEventListener('click', function() {
            const settingId = this.getAttribute('data-setting');
            const defaultValue = this.getAttribute('data-default');
            
            const slider = document.getElementById(settingId);
            const valueInput = document.getElementById(`${settingId}-value`);
            
            if (slider && valueInput) {
                slider.value = defaultValue;
                valueInput.value = defaultValue;
                saveSamplerSettings();
            }
        });
    });
}

// Handle conflict detection for mutually exclusive settings
function handleConflictDetection(settingId, isEnabled) {
    // Define conflict groups
    const conflictGroups = {
        'top-p': ['typical-p', 'min-p'],
        'min-p': ['top-p', 'typical-p'],
        'mirostat-mode': ['dyn-temp-range'],
        'dyn-temp-range': ['mirostat-mode'],
        'presence-penalty': ['frequency-penalty'],
        'frequency-penalty': ['presence-penalty']
    };

    // Get conflicts for this setting
    const conflicts = conflictGroups[settingId];
    if (!conflicts) return;

    // Handle each conflicting setting
    conflicts.forEach(conflictId => {
        const conflictToggle = document.getElementById(`${conflictId}-enabled`);
        const conflictCard = document.getElementById(`${conflictId}-card`);
        
        if (conflictToggle && conflictCard) {
            if (isEnabled) {
                // Disable conflicting setting when this one is enabled
                conflictToggle.checked = false;
                conflictCard.classList.add('disabled');
            }
        }
    });
}

// Frontend slider ID → backend key mapping for IDs that don't follow the simple dash-to-underscore rule.
const SAMPLER_ID_TO_BACKEND = {
    'mirostat-mode':        'mirostat',
    'dyn-temp-range':       'dynamic_temp_range',
    'dyn-temp-exp':         'dynamic_temp_exponent'
};

// Backend key → frontend slider ID (reverse of above)
const SAMPLER_BACKEND_TO_ID = {};
for (const [fid, bid] of Object.entries(SAMPLER_ID_TO_BACKEND)) {
    SAMPLER_BACKEND_TO_ID[bid] = fid;
}

// Get all sampler settings as an object
function getSamplerSettings() {
    const settings = {};
    samplerSettings.forEach(setting => {
        const slider = document.getElementById(setting.id);
        if (slider) {
            // Use explicit mapping if available, otherwise dash → underscore
            const backendId = SAMPLER_ID_TO_BACKEND[setting.id] || setting.id.replace(/-/g, '_');
            settings[backendId] = parseFloat(slider.value);

            // Include enabled state if toggle exists
            if (setting.hasToggle) {
                const toggle = document.getElementById(`${setting.id}-enabled`);
                if (toggle) {
                    settings[`${backendId}_enabled`] = toggle.checked;
                }
            }
        }
    });
    return settings;
}

// Set sampler settings from an object
function setSamplerSettings(settings) {
    Object.keys(settings).forEach(backendId => {
        // Toggle state is stored as `<key>_enabled`. The suffix is part of the
        // storage key, not of the element id, so it has to be stripped before
        // mapping or every toggle resolves to `<id>-enabled-enabled` and is
        // silently skipped — which is why disabled samplers came back enabled
        // after a restart.
        const isEnabledKey = backendId.endsWith('_enabled');
        const baseKey = isEnabledKey ? backendId.slice(0, -'_enabled'.length) : backendId;

        // Use explicit reverse mapping if available, otherwise underscore → dash
        let frontendId = SAMPLER_BACKEND_TO_ID[baseKey];
        if (!frontendId) {
            frontendId = baseKey.replace(/_/g, '-');
        }
        const slider = document.getElementById(frontendId);
        const valueInput = document.getElementById(`${frontendId}-value`);
        const toggle = document.getElementById(`${frontendId}-enabled`);
        const card = document.getElementById(`${frontendId}-card`);

        if (slider && valueInput) {
            const value = settings[backendId];
            slider.value = value;
            valueInput.value = value;
        }

        // Handle toggle state
        if (isEnabledKey && toggle && card) {
            const isEnabled = settings[backendId];
            toggle.checked = isEnabled;
            if (isEnabled) {
                card.classList.remove('disabled');
            } else {
                card.classList.add('disabled');
            }
        }
    });
}

// Save sampler settings through the settings seam. The seam writes the whole
// authoritative document to the backend before it updates any renderer state,
// so this module never talks to a second backend binding.
async function saveSamplerSettings() {
    if (!window.settingsState?.set) {
        console.error('Failed to save sampler settings: settings seam unavailable');
        return;
    }
    const saved = await window.settingsState.set('sampler', getSamplerSettings());
    if (!saved) console.error('Failed to save sampler settings: backend rejected the write');
}

// Load sampler settings from the authoritative settings snapshot, or defaults.
async function loadSamplerSettings() {
    const settings = window.settingsState?.get?.('sampler');
    if (settings && Object.keys(settings).length > 0) {
        setSamplerSettings(settings);
        return;
    }

    samplerSettings.forEach(setting => {
        const slider = document.getElementById(setting.id);
        const valueInput = document.getElementById(`${setting.id}-value`);
        if (slider && valueInput) {
            slider.value = setting.default;
            valueInput.value = setting.default;
        }
    });
}

// Apply a named preset to all sampler controls
function applySamplerPreset(presetName) {
    const preset = SAMPLER_PRESETS[presetName];
    if (!preset) return;
    setSamplerSettings(preset);
    saveSamplerSettings();
    console.log('Applied sampler preset:', presetName);
}

// Reset all sampler settings to defaults
function resetSamplerSettings() {
    samplerSettings.forEach(setting => {
        const slider = document.getElementById(setting.id);
        const valueInput = document.getElementById(`${setting.id}-value`);

        if (slider && valueInput) {
            slider.value = setting.default;
            valueInput.value = setting.default;
        }
    });
    saveSamplerSettings();
}

// After a model loads, show the values the engine is actually using. The
// family defaults live in the backend, so without this the sliders would keep
// displaying the last saved number while inference ran on something else.
async function refreshEffectiveSampler() {
    const app = window.go?.main?.App;
    if (!app?.GetEffectiveSampler) return;
    try {
        const eff = await app.GetEffectiveSampler();
        if (!eff) return;
        const shown = {};
        ["temperature", "top_p", "top_k", "min_p", "repeat_penalty", "repeat_last_n",
         "presence_penalty", "frequency_penalty"].forEach(k => {
            if (typeof eff[k] === "number") shown[k] = eff[k];
        });
        setSamplerSettings(shown);
        console.log('[sampler] showing effective values for family:', eff.family || 'unknown');
    } catch (err) {
        console.warn('[sampler] could not read effective values:', err);
    }
}

// Initialize on DOM load
document.addEventListener('DOMContentLoaded', function() {
    initSamplerSettings();
    refreshEffectiveSampler();
});

// ---------------------------------------------------------------------------
// Clamped descriptions keep their full sentence reachable
//
// The full-screen sheet caps .setting-description at two lines so one long
// description cannot stretch its whole grid row, which would undo the compact
// layout. The complete sentence is exposed as a native tooltip so nothing is
// hidden, only shortened.
// ---------------------------------------------------------------------------

document.addEventListener('DOMContentLoaded', function () {
    document.querySelectorAll('#sampler-settings-flyout .setting-description').forEach(function (el) {
        const full = el.textContent.trim();
        if (full && !el.title) el.title = full;
    });
});
