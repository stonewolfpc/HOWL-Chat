/**
 * Model Settings Module
 * 
 * This module handles the Model Settings flyout panel toggle
 * and manages the UI for model configuration options.
 * 
 * @module model-settings
 */

// Model settings state — synced to unified store under the 'model' key
let modelSettings = {
    // Core Model Settings (keys match app.go/spawnLlamaServerWithModel)
    model_family: 'auto',
    model_variant: 'auto',
    threads: 8,
    batch_size: 512,
    gpu_layers: 0,
    context_size: 4096,
    kv_cache_type: 'f16',

    // Performance Settings
    rope_mode: 'auto',
    rope_factor: 1.0,
    rope_base: 0,
    tensor_split: 0,
    flash_attention_mode: 'auto',
    offload_kqv: true,
    use_mmap: true,
    use_mlock: false,

    // Prompt & Template Settings (handled per-request)
    prompt_template: 'auto',
    custom_jinja_template: '',
    system_prompt_override: '',
    user_prefix: '',
    assistant_prefix: '',
    stop_sequences: '',

    // App-level only (no llama-server equivalent)
    context_rollover_mode: 'memory-guided',
    memory_integration: true,
    system_prompt_cache: true,
    tokenizer_override: 'auto'
};

// Toggle Model Settings flyout
function toggleModelSettingsFlyout() {
    const overlay = document.getElementById('model-settings-overlay');
    const flyout = document.getElementById('model-settings-flyout');
    
    if (overlay && flyout) {
        overlay.classList.toggle('show');
        flyout.classList.toggle('show');
        
        // Load settings when opening
        if (flyout.classList.contains('show')) {
            loadModelSettings();
        }
    }
}

// Toggle Memory Integration switch
function toggleMemoryIntegration() {
    const toggle = document.getElementById('memory-integration-toggle');
    const label = toggle.nextElementSibling;
    
    if (toggle) {
        toggle.classList.toggle('active');
        modelSettings.memory_integration = toggle.classList.contains('active');
        label.textContent = modelSettings.memory_integration ? 'Enabled' : 'Disabled';
        saveModelSettings();
    }
}

// Toggle System Prompt Cache switch
function toggleSystemPromptCache() {
    const toggle = document.getElementById('system-prompt-cache-toggle');
    const label = toggle.nextElementSibling;
    
    if (toggle) {
        toggle.classList.toggle('active');
        modelSettings.system_prompt_cache = toggle.classList.contains('active');
        label.textContent = modelSettings.system_prompt_cache ? 'Enabled' : 'Disabled';
        saveModelSettings();
    }
}

// Initialize slider-value pairs
function initSliderValuePairs() {
    const sliderValuePairs = [
        { slider: 'thread-count', value: 'thread-count-value', setting: 'threads' },
        { slider: 'batch-size', value: 'batch-size-value', setting: 'batch_size' },
        { slider: 'gpu-layers', value: 'gpu-layers-value', setting: 'gpu_layers' },
        { slider: 'rope-factor', value: 'rope-factor-value', setting: 'rope_factor' },
        { slider: 'tensor-split', value: 'tensor-split-value', setting: 'tensor_split' },
        { slider: 'max-context', value: 'max-context-value', setting: 'context_size' }
    ];
    
    sliderValuePairs.forEach(pair => {
        const slider = document.getElementById(pair.slider);
        const valueInput = document.getElementById(pair.value);
        
        if (slider && valueInput) {
            // Slider change updates number input
            slider.addEventListener('input', function() {
                valueInput.value = this.value;
                modelSettings[pair.setting] = parseFloat(this.value);
                saveModelSettings();
            });
            
            // Number input change updates slider
            valueInput.addEventListener('input', function() {
                const min = parseFloat(slider.min);
                const max = parseFloat(slider.max);
                let value = parseFloat(this.value);
                
                // Clamp value to min/max range
                if (isNaN(value)) {
                    value = min;
                }
                if (value < min) value = min;
                if (value > max) value = max;
                
                slider.value = value;
                modelSettings[pair.setting] = value;
                saveModelSettings();
            });
        }
    });
}

// Initialize dropdown selects
function initDropdowns() {
    const dropdowns = [
        { id: 'model-family', setting: 'model_family' },
        { id: 'model-variant', setting: 'model_variant' },
        { id: 'kv-cache-type', setting: 'kv_cache_type' },
        { id: 'flash-attention-mode', setting: 'flash_attention_mode' },
        { id: 'rope-scaling-mode', setting: 'rope_mode' },
        { id: 'context-rollover-mode', setting: 'context_rollover_mode' },
        { id: 'prompt-template', setting: 'prompt_template' },
        { id: 'tokenizer-override', setting: 'tokenizer_override' }
    ];
    
    dropdowns.forEach(dropdown => {
        const element = document.getElementById(dropdown.id);
        if (element) {
            element.addEventListener('change', function() {
                modelSettings[dropdown.setting] = this.value;
                saveModelSettings();
            });
        }
    });
}

// Initialize text inputs
function initTextInputs() {
    const textInputs = [
        { id: 'custom-jinja-template', setting: 'custom_jinja_template' },
        { id: 'system-prompt-override', setting: 'system_prompt_override' },
        { id: 'user-prompt-prefix', setting: 'user_prefix' },
        { id: 'assistant-prompt-prefix', setting: 'assistant_prefix' },
        { id: 'stop-sequences', setting: 'stop_sequences' }
    ];
    
    textInputs.forEach(input => {
        const element = document.getElementById(input.id);
        if (element) {
            element.addEventListener('input', function() {
                modelSettings[input.setting] = this.value;
                saveModelSettings();
            });
        }
    });
}

// The KV location select is a boolean in the engine (offload_kqv). Keeping the
// mapping here means the saved document still has one authoritative key.
function initKVLocationControl() {
    const element = document.getElementById('kv-cache-location');
    if (!element) return;
    element.addEventListener('change', function () {
        modelSettings.offload_kqv = this.value === 'gpu';
        saveModelSettings();
    });
}

// Migrate values written by older builds before the UI reads them. The old KV
// size duplicated the context control and was never sent; the old KV type
// offered kernel names the engine does not accept as cache types.
function normalizeModelSettings() {
    if (!['f16', 'q8_0', 'q4_0'].includes(modelSettings.kv_cache_type)) {
        modelSettings.kv_cache_type = 'f16';
    }
    if (!['auto', 'on', 'off'].includes(modelSettings.flash_attention_mode)) {
        modelSettings.flash_attention_mode = modelSettings.flash_attention ? 'on' : 'auto';
    }
    delete modelSettings.kv_cache_size;
    delete modelSettings.flash_attention;
}

// An option the shipped binary does not advertise becomes unavailable with a
// visible reason instead of a control that silently does nothing.
function markSettingUnavailable(element, reason) {
    if (!element) return;
    element.disabled = true;
    element.title = reason;
    const card = element.closest('.setting-card');
    if (card && !card.querySelector('[data-capability-reason]')) {
        const note = document.createElement('div');
        note.className = 'setting-detail warning';
        note.setAttribute('data-capability-reason', 'true');
        note.textContent = reason;
        card.appendChild(note);
    }
}

async function applyLlamaCapabilityGates() {
    const app = window.go?.main?.App;
    if (typeof app?.GetLlamaCapabilities !== 'function') return;
    let caps;
    try {
        caps = await app.GetLlamaCapabilities();
    } catch (e) {
        console.warn('[model-settings] llama capability probe unavailable:', e);
        return;
    }
    const supports = caps?.supports || {};
    const build = caps?.version ? `this llama-server (${caps.version})` : 'the shipped llama-server';

    const flash = document.getElementById('flash-attention-mode');
    if (flash && !supports.flash_attn) {
        markSettingUnavailable(flash, `${build} does not advertise --flash-attn; Flash Attention cannot be selected.`);
    }

    const kvType = document.getElementById('kv-cache-type');
    if (kvType && (!supports.cache_type_k || !supports.cache_type_v)) {
        for (const option of kvType.options) {
            if (option.value !== 'f16') {
                option.disabled = true;
            }
        }
        kvType.title = `${build} does not advertise both --cache-type-k and --cache-type-v; only the standard f16 cache is available.`;
        const card = kvType.closest('.setting-card');
        if (card && !card.querySelector('[data-capability-reason]')) {
            const note = document.createElement('div');
            note.className = 'setting-detail warning';
            note.setAttribute('data-capability-reason', 'true');
            note.textContent = kvType.title;
            card.appendChild(note);
        }
    }

    const kvLocation = document.getElementById('kv-cache-location');
    if (kvLocation && !supports.kv_offload && !supports.no_kv_offload) {
        markSettingUnavailable(kvLocation, `${build} does not advertise KV offload control; the engine default is used.`);
    }
}

// Save model settings through the settings seam. The seam writes the backend
// document before any renderer state changes, so there is no second path.
async function saveModelSettings() {
    if (!window.settingsState?.set) {
        console.error('Error saving model settings: settings seam unavailable');
        return;
    }
    const saved = await window.settingsState.set('model', modelSettings);
    if (!saved) console.error('Error saving model settings: backend rejected the write');
}

// Load model settings from the authoritative settings snapshot, or defaults.
async function loadModelSettings() {
    const settings = window.settingsState?.get?.('model');
    if (settings && Object.keys(settings).length > 0) {
        modelSettings = { ...modelSettings, ...settings };
    }
    normalizeModelSettings();
    updateUIFromSettings();
}

// Update UI from settings
function updateUIFromSettings() {
    // Update sliders and values
    const sliderValuePairs = [
        { slider: 'thread-count', value: 'thread-count-value', setting: 'threads' },
        { slider: 'batch-size', value: 'batch-size-value', setting: 'batch_size' },
        { slider: 'gpu-layers', value: 'gpu-layers-value', setting: 'gpu_layers' },
        { slider: 'rope-factor', value: 'rope-factor-value', setting: 'rope_factor' },
        { slider: 'tensor-split', value: 'tensor-split-value', setting: 'tensor_split' },
        { slider: 'max-context', value: 'max-context-value', setting: 'context_size' }
    ];
    
    sliderValuePairs.forEach(pair => {
        const slider = document.getElementById(pair.slider);
        const valueInput = document.getElementById(pair.value);
        
        if (slider && valueInput) {
            slider.value = modelSettings[pair.setting];
            valueInput.value = modelSettings[pair.setting];
        }
    });
    
    // Update dropdowns
    const dropdowns = [
        { id: 'model-family', setting: 'model_family' },
        { id: 'model-variant', setting: 'model_variant' },
        { id: 'kv-cache-type', setting: 'kv_cache_type' },
        { id: 'flash-attention-mode', setting: 'flash_attention_mode' },
        { id: 'rope-scaling-mode', setting: 'rope_mode' },
        { id: 'context-rollover-mode', setting: 'context_rollover_mode' },
        { id: 'prompt-template', setting: 'prompt_template' },
        { id: 'tokenizer-override', setting: 'tokenizer_override' }
    ];
    
    dropdowns.forEach(dropdown => {
        const element = document.getElementById(dropdown.id);
        if (element) {
            element.value = modelSettings[dropdown.setting];
        }
    });
    
    // Update text inputs
    const textInputs = [
        { id: 'custom-jinja-template', setting: 'custom_jinja_template' },
        { id: 'system-prompt-override', setting: 'system_prompt_override' },
        { id: 'user-prompt-prefix', setting: 'user_prefix' },
        { id: 'assistant-prompt-prefix', setting: 'assistant_prefix' },
        { id: 'stop-sequences', setting: 'stop_sequences' }
    ];
    
    textInputs.forEach(input => {
        const element = document.getElementById(input.id);
        if (element) {
            element.value = modelSettings[input.setting];
        }
    });
    
    const kvLocation = document.getElementById('kv-cache-location');
    if (kvLocation) {
        kvLocation.value = modelSettings.offload_kqv === false ? 'ram' : 'gpu';
    }

    // Update toggle switches
    const memoryToggle = document.getElementById('memory-integration-toggle');
    if (memoryToggle) {
        if (modelSettings.memory_integration) {
            memoryToggle.classList.add('active');
            memoryToggle.nextElementSibling.textContent = 'Enabled';
        } else {
            memoryToggle.classList.remove('active');
            memoryToggle.nextElementSibling.textContent = 'Disabled';
        }
    }
    
    const systemPromptToggle = document.getElementById('system-prompt-cache-toggle');
    if (systemPromptToggle) {
        if (modelSettings.system_prompt_cache) {
            systemPromptToggle.classList.add('active');
            systemPromptToggle.nextElementSibling.textContent = 'Enabled';
        } else {
            systemPromptToggle.classList.remove('active');
            systemPromptToggle.nextElementSibling.textContent = 'Disabled';
        }
    }
}

// Initialize on DOM load
document.addEventListener('DOMContentLoaded', function() {
    initSliderValuePairs();
    initDropdowns();
    initTextInputs();
    initKVLocationControl();
    
    // Load settings once the authoritative snapshot is in memory, then hide or
    // disable controls the shipped llama-server does not advertise.
    (window.settingsState?.ready || Promise.resolve(true)).then(async function () {
        await loadModelSettings();
        await applyLlamaCapabilityGates();
    });
});

// ---------------------------------------------------------------------------
// Clamped descriptions keep their full sentence reachable
//
// The full-screen sheet caps .setting-description at two lines so one long
// description cannot stretch its whole grid row. The complete sentence stays
// available as a native tooltip.
// ---------------------------------------------------------------------------

document.addEventListener('DOMContentLoaded', function () {
    document.querySelectorAll('#model-settings-flyout .setting-description').forEach(function (el) {
        const full = el.textContent.trim();
        if (full && !el.title) el.title = full;
    });
});
