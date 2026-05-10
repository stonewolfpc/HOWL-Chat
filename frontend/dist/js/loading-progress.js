/**
 * Loading Progress Bar Module
 * 
 * This module handles the creation and management of the loading progress bar
 * for model loading. It dynamically injects the progress bar into the DOM
 * to avoid modifying the HTML structure directly.
 * 
 * @module loading-progress
 */

// Create and inject loading progress bar
function createLoadingProgressBar() {
    const topMenu = document.querySelector('.top-menu');
    if (!topMenu) {
        console.error('Top menu not found');
        return;
    }

    // Create loading progress container
    const container = document.createElement('div');
    container.className = 'loading-progress-container';
    container.id = 'loading-progress-container';

    // Create progress bar
    const progressBar = document.createElement('div');
    progressBar.className = 'loading-progress-bar';
    progressBar.id = 'loading-progress-bar';

    // Create progress fill
    const progressFill = document.createElement('div');
    progressFill.className = 'loading-progress-fill';
    progressFill.id = 'loading-progress-fill';
    progressFill.style.width = '0%';

    // Create status text
    const status = document.createElement('div');
    status.className = 'loading-status';
    status.id = 'loading-status';
    status.textContent = 'no model loaded';

    // Assemble elements
    progressBar.appendChild(progressFill);
    container.appendChild(progressBar);
    container.appendChild(status);

    // Inject into top menu as a flex item
    topMenu.appendChild(container);

    console.log('Loading progress bar created and injected');
}

// Update progress bar
function updateLoadingProgress(progress) {
    const progressFill = document.getElementById('loading-progress-fill');
    const progressBar = document.getElementById('loading-progress-bar');
    
    if (progressFill && progressBar) {
        progressFill.style.width = progress + '%';
        
        if (progress > 0 && progress < 100) {
            progressBar.classList.add('loading');
        }
    }
}

// Set loading status
function setLoadingStatus(status, state = 'default') {
    const statusElement = document.getElementById('loading-status');
    const progressBar = document.getElementById('loading-progress-bar');
    
    if (statusElement) {
        statusElement.textContent = status;
    }
    
    if (progressBar) {
        progressBar.classList.remove('loading', 'completed', 'error');
        if (state !== 'default') {
            progressBar.classList.add(state);
        }
        
        if (statusElement) {
            statusElement.classList.remove('completed', 'error');
            if (state !== 'default') {
                statusElement.classList.add(state);
            }
        }
    }
}

// Initialize loading progress bar on page load
document.addEventListener('DOMContentLoaded', () => {
    createLoadingProgressBar();
    initModelEventListeners();
});

// Listen for Wails backend events
function initModelEventListeners() {
    if (window.runtime) {
        // Model loading started
        window.runtime.EventsOn('model:loading', () => {
            updateLoadingProgress(0);
            setLoadingStatus('Loading...', 'loading');
        });

        // Progress updates (0-100)
        window.runtime.EventsOn('model:progress', (progress) => {
            updateLoadingProgress(progress);
        });

        // Model ready
        window.runtime.EventsOn('model:loaded', () => {
            updateLoadingProgress(100);
            setLoadingStatus('Ready', 'completed');
        });

        // Model unloaded
        window.runtime.EventsOn('model:unloaded', () => {
            updateLoadingProgress(0);
            setLoadingStatus('no model loaded', 'default');
        });

        // Error during loading
        window.runtime.EventsOn('model:error', (err) => {
            updateLoadingProgress(0);
            setLoadingStatus('Error', 'error');
            console.error('Model loading error:', err);
        });
    }
}
