/**
 * Error Handler - Integration layer
 * Connects existing UI components (loading progress, chat model) to ErrorSystem.
 * Provides recovery strategies for each subsystem.
 */

'use strict';

// ---------------------------------------------------------------------------
// Subsystem recovery strategies
// ---------------------------------------------------------------------------

const RecoveryStrategies = {
    'llm.load': async () => {
        setLoadingStatus('Model load failed - retry available', 'error');
        updateModelDisplay('Load failed');
        return { canRetry: true, retryDelay: 5000 };
    },
    'llm.infer': async () => {
        updateModelDisplay('Inference error');
        return { canRetry: true, retryDelay: 2000 };
    },
    'image.load': async () => {
        const status = document.getElementById('image-status');
        if (status) status.textContent = 'Model load failed';
        return { canRetry: true, retryDelay: 3000 };
    },
    'image.generate': async () => {
        const status = document.getElementById('image-status');
        if (status) status.textContent = 'Generation failed';
        return { canRetry: true, retryDelay: 2000 };
    },
    'tts.synthesize': async () => {
        return { canRetry: false };
    },
    'chat.save': async () => {
        showSaveWarning();
        return { canRetry: true, retryDelay: 1000 };
    },
    'resource.allocate': async () => {
        updateModelDisplay('Resource conflict');
        return { canRetry: true, retryDelay: 4000 };
    },
    'default': async () => {
        return { canRetry: false };
    }
};

// ---------------------------------------------------------------------------
// Recovery orchestration
// ---------------------------------------------------------------------------

async function executeRecovery(operation) {
    const strategy = RecoveryStrategies[operation] || RecoveryStrategies['default'];
    const result = await strategy();

    if (result.canRetry && result.retryDelay) {
        setTimeout(() => clearDegradedState(operation), result.retryDelay);
    }

    return result;
}

function clearDegradedState(operation) {
    // Clear visual degraded indicators after successful recovery window
    if (operation.startsWith('llm.')) {
        checkModelStatus();
    }
    if (operation.startsWith('image.')) {
        const status = document.getElementById('image-status');
        if (status && status.textContent.includes('failed')) {
            status.textContent = '';
        }
    }
}

// ---------------------------------------------------------------------------
// UI helpers
// ---------------------------------------------------------------------------

function showSaveWarning() {
    let warn = document.getElementById('chat-save-warning');
    if (!warn) {
        warn = document.createElement('div');
        warn.id = 'chat-save-warning';
        warn.style.cssText = 'position:fixed;bottom:16px;left:50%;transform:translateX(-50%);background:#f59e0b;color:white;padding:8px 16px;border-radius:6px;font-size:13px;z-index:10000;';
        document.body.appendChild(warn);
    }
    warn.textContent = 'Chat save failed - will retry';
    warn.style.display = 'block';
    setTimeout(() => { warn.style.display = 'none'; }, 4000);
}

// ---------------------------------------------------------------------------
// Integration with ErrorSystem
// ---------------------------------------------------------------------------

function onErrorReported(err) {
    // Hook called by ErrorSystem after each error
    if (err && err.operation) {
        executeRecovery(err.operation).catch(() => {});
    }
}

// ErrorSystem calls onErrorReported() itself now, so handleBackendError is not
// monkey-patched. Recovery happens even if this file loads after error-system.

// ---------------------------------------------------------------------------
// Global error boundary
// ---------------------------------------------------------------------------

window.addEventListener('error', (e) => {
    logFrontendError('FE-WINDOW-ERR', 'window.onerror', e.error || e.message);
});

window.addEventListener('unhandledrejection', (e) => {
    logFrontendError('FE-PROMISE-REJ', 'unhandledrejection', e.reason);
});

// ---------------------------------------------------------------------------
// Export
// ---------------------------------------------------------------------------

window.RecoveryStrategies = RecoveryStrategies;
window.executeRecovery = executeRecovery;
window.showSaveWarning = showSaveWarning;
