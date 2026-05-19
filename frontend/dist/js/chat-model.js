/**
 * Chat Model Management
 *
 * Load/unload the LLM model via the Wails backend,
 * read model metadata to update UI sliders, and
 * display the currently loaded model name.
 *
 * @module chat-model
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Model Display
// ---------------------------------------------------------------------------

function updateModelDisplay(modelName) {
    const display = document.getElementById('model-name-display');
    if (display) display.textContent = modelName;
}

// ---------------------------------------------------------------------------
// Model Status Check
// ---------------------------------------------------------------------------

async function autoLoadDefaultModel() {
    const defaultModel = window.settingsState?.current?.general?.defaultModel;
    if (defaultModel !== 'last-used') return;
    const savedPath = localStorage.getItem('howl-last-model-path');
    if (!savedPath) return;
    if (!window.go?.main?.App?.LoadModel) return;

    updateModelDisplay('Loading...');
    try {
        await window.go.main.App.LoadModel(savedPath);
        const modelName = savedPath.split(/[\\/]/).pop();
        updateModelDisplay(modelName);
        console.log('Auto-loaded model:', modelName);
    } catch (err) {
        console.error('Auto-load model failed:', err);
        updateModelDisplay('No model loaded');
    }
}

function checkModelStatus() {
    if (!window.go || !window.go.main || !window.go.main.App) {
        console.log('Wails runtime not available, skipping model status check');
        return;
    }
    try {
        window.go.main.App.IsModelLoaded().then(isLoaded => {
            if (isLoaded) {
                window.go.main.App.GetLoadedModelName().then(modelName => {
                    updateModelDisplay(modelName);
                }).catch(error => {
                    console.error('Error getting model name:', error);
                    updateModelDisplay('Error getting model name');
                });
            } else {
                updateModelDisplay('No model loaded');
                autoLoadDefaultModel();
            }
        }).catch(error => {
            console.error('Error checking model status:', error);
            updateModelDisplay('Error checking status');
        });
    } catch (error) {
        console.error('Error checking model status:', error);
        updateModelDisplay('No model loaded');
    }
}

// ---------------------------------------------------------------------------
// Load Model
// ---------------------------------------------------------------------------

async function loadModel() {
    if (!window.go || !window.go.main || !window.go.main.App) {
        console.error('Wails runtime not available');
        return;
    }

    try {
        const modelPath = await window.go.main.App.OpenModelFilePicker();
        if (!modelPath) return;

        updateModelDisplay('Loading...');

        try {
            const metadata = await window.go.main.App.LoadModelMetadata(modelPath);
            if (metadata) {
                const gpuLayersSlider = document.getElementById('gpu-layers');
                const gpuLayersValue  = document.getElementById('gpu-layers-value');
                if (gpuLayersSlider && gpuLayersValue && metadata.LayerCount > 0) {
                    gpuLayersSlider.max = metadata.LayerCount;
                    gpuLayersValue.max  = metadata.LayerCount;
                    console.log(`Updated GPU layers slider max to ${metadata.LayerCount}`);
                }

                const contextSlider = document.getElementById('max-context');
                const contextValue  = document.getElementById('max-context-value');
                if (contextSlider && contextValue && metadata.MaxContext > 0) {
                    contextSlider.max = metadata.MaxContext;
                    contextValue.max  = metadata.MaxContext;
                    if (parseInt(contextSlider.value) > metadata.UsableContext) {
                        contextSlider.value = metadata.UsableContext;
                        contextValue.value  = metadata.UsableContext;
                    }
                    console.log(`Updated context slider max to ${metadata.MaxContext} (usable: ${metadata.UsableContext})`);
                }
            }
        } catch (error) {
            console.warn('Failed to load model metadata:', error);
        }

        await window.go.main.App.LoadModel(modelPath);
        const modelName = modelPath.split(/[\\/]/).pop();
        updateModelDisplay(modelName);
        localStorage.setItem('howl-last-model-path', modelPath);
        console.log('Model loaded successfully:', modelName);

    } catch (error) {
        updateModelDisplay('Error loading model');
        console.error('Failed to load model:', error);
        alert('Failed to load model: ' + error);
    }
}

// ---------------------------------------------------------------------------
// Unload Model
// ---------------------------------------------------------------------------

async function unloadModel() {
    if (!confirm('⚠️ AGGRESSIVE UNLOAD ⚠️\n\nThis will forcefully unload the current model from memory.\nAll unsaved context will be lost.\n\nAre you sure you want to proceed?')) return;

    console.log('AGGRESSIVE UNLOAD initiated');
    try {
        await window.go.main.App.UnloadModel();
        updateModelDisplay('No model loaded');
        alert('🔴 MODEL AGGRESSIVELY UNLOADED\n\nGPU cache cleared.\nMemory freed.\nModel process terminated.');
    } catch (error) {
        console.error('Error unloading model:', error);
        alert(`Failed to unload model: ${error}`);
    }
}

// ---------------------------------------------------------------------------
// Init on DOM Ready
// ---------------------------------------------------------------------------

document.addEventListener('DOMContentLoaded', () => { checkModelStatus(); });
