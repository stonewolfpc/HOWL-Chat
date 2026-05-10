/**
 * Sampler Settings Module
 * 
 * This module handles the synchronization between sliders and number inputs
 * for the sampler settings dropdown menu, and communicates with the backend.
 * 
 * @module sampler-settings
 */

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
    { id: 'no-repeat-ngram', default: 0 },
    { id: 'mirostat-mode', default: 0, hasToggle: true },
    { id: 'mirostat-tau', default: 5.0 },
    { id: 'mirostat-eta', default: 0.1 },
    { id: 'dyn-temp-range', default: 0.0, hasToggle: true },
    { id: 'dyn-temp-exp', default: 1.0 },
    { id: 'dry-multiplier', default: 0.0 },
    { id: 'dry-allowed-length', default: 2 },
    { id: 'dry-base', default: 1.0 },
    { id: 'smoothing-factor', default: 0.0 },
    { id: 'smoothing-curve', default: 1.0 },
    { id: 'top-a', default: 0.0, hasToggle: true },
    { id: 'epsilon-cutoff', default: 0.0 },
    { id: 'eta-cutoff', default: 0.0 },
    { id: 'encoder-repeat-penalty', default: 1.0 },
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
    
    // Load settings from backend on initialization
    loadSamplerSettings();
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
        'top-p': ['typical-p', 'top-a'],
        'typical-p': ['top-p'],
        'top-a': ['top-p'],
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

// Get all sampler settings as an object
function getSamplerSettings() {
    const settings = {};
    samplerSettings.forEach(setting => {
        const slider = document.getElementById(setting.id);
        if (slider) {
            // Convert ID format (e.g., 'top-p' to 'top_p')
            const backendId = setting.id.replace(/-/g, '_');
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
        // Convert backend ID format (e.g., 'top_p' to 'top-p')
        const frontendId = backendId.replace(/_/g, '-');
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
        if (backendId.endsWith('_enabled') && toggle && card) {
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

// Save sampler settings to backend
async function saveSamplerSettings() {
    try {
        const settings = getSamplerSettings();
        if (window.go && window.go.main && window.go.main.App) {
            await window.go.main.App.SetSamplerSettings(settings);
        }
    } catch (error) {
        console.error('Failed to save sampler settings:', error);
    }
}

// Load sampler settings from backend
async function loadSamplerSettings() {
    try {
        if (window.go && window.go.main && window.go.main.App) {
            const settings = await window.go.main.App.GetSamplerSettings();
            if (settings) {
                setSamplerSettings(settings);
            }
        }
    } catch (error) {
        console.error('Failed to load sampler settings:', error);
        // Use defaults if backend call fails
        samplerSettings.forEach(setting => {
            const slider = document.getElementById(setting.id);
            const valueInput = document.getElementById(`${setting.id}-value`);
            if (slider && valueInput) {
                slider.value = setting.default;
                valueInput.value = setting.default;
            }
        });
    }
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

// Initialize on DOM load
document.addEventListener('DOMContentLoaded', function() {
    initSamplerSettings();
});
