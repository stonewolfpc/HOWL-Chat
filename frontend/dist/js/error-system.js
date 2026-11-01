/**
 * Error System - Core
 * Silent until failure. Fail loud, recover smooth.
 * No console spam. Logs backend errors with minimal UI footprint.
 */

'use strict';

// ---------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------

const ErrorSystem = {
    initialized: false,
    activeErrors: new Map(),
    suppressed: new Set(),
    circuitBreaker: new Map(),
    maxRetries: 3,
    listenerBound: false
};

// ---------------------------------------------------------------------------
// Initialization
// ---------------------------------------------------------------------------

function initErrorSystem() {
    if (ErrorSystem.initialized) return;
    if (!window.runtime) {
        setTimeout(initErrorSystem, 500);
        return;
    }

    bindErrorEvents();
    detectHardwareEdgeCases();
    ErrorSystem.initialized = true;
}

function bindErrorEvents() {
    if (ErrorSystem.listenerBound) return;
    ErrorSystem.listenerBound = true;

    window.runtime.EventsOn('app:error', (payload) => {
        handleBackendError(parsePayload(payload));
    });

    window.runtime.EventsOn('app:error:silent', (payload) => {
        handleSilentError(parsePayload(payload));
    });

    window.runtime.EventsOn('app:recovered', (payload) => {
        handleRecoveryNotice(parsePayload(payload));
    });
}

// ---------------------------------------------------------------------------
// Payload parsing
// ---------------------------------------------------------------------------

function parsePayload(raw) {
    if (typeof raw === 'string') {
        try { return JSON.parse(raw); } catch { return { code: raw }; }
    }
    return raw || {};
}

// ---------------------------------------------------------------------------
// Error handling logic
// ---------------------------------------------------------------------------

function handleBackendError(payload) {
    const code = payload.code || 'UNKNOWN';
    const severity = payload.severity || 'medium';
    const operation = payload.operation || 'unknown';

    // Circuit breaker: if same code failed too many times, suppress UI
    if (shouldSuppress(code)) {
        ErrorSystem.suppressed.add(code);
        return;
    }

    // Track active error
    const errId = `${code}-${Date.now()}`;
    ErrorSystem.activeErrors.set(errId, {
        code, severity, operation, time: Date.now(), payload
    });

    // Fail loud for critical/high
    if (severity === 'critical' || severity === 'high') {
        showErrorBanner(code, operation, severity);
    }

    // Medium/low: subtle indicator only
    if (severity === 'medium' || severity === 'low') {
        updateErrorIndicator();
    }

    incrementCircuitBreaker(code);

    // Recovery hooks (error-handler.js) run from here now. It used to wrap this
    // function after load, which only worked while the load order held.
    if (typeof onErrorReported === 'function') onErrorReported(payload);
}

function handleSilentError(payload) {
    // Logged to file but never shown in UI unless user opens settings
    const code = payload.code || 'UNKNOWN';
    console.debug(`[ErrorSystem] Silent: ${code}`);
}

function handleRecoveryNotice(payload) {
    const operation = payload.operation || 'unknown';
    showRecoveryToast(operation);
}

// ---------------------------------------------------------------------------
// Circuit breaker
// ---------------------------------------------------------------------------

function shouldSuppress(code) {
    const count = ErrorSystem.circuitBreaker.get(code) || 0;
    return count >= ErrorSystem.maxRetries;
}

function incrementCircuitBreaker(code) {
    const count = (ErrorSystem.circuitBreaker.get(code) || 0) + 1;
    ErrorSystem.circuitBreaker.set(code, count);

    // Reset after 60 seconds
    setTimeout(() => {
        const current = ErrorSystem.circuitBreaker.get(code) || 0;
        if (current > 0) {
            ErrorSystem.circuitBreaker.set(code, current - 1);
        }
    }, 60000);
}

function resetCircuitBreaker(code) {
    ErrorSystem.circuitBreaker.delete(code);
    ErrorSystem.suppressed.delete(code);
}

// ---------------------------------------------------------------------------
// UI helpers (minimal)
// ---------------------------------------------------------------------------

function showErrorBanner(code, operation, severity) {
    const banner = getOrCreateBanner();
    const color = severity === 'critical' ? '#ef4444' : '#f59e0b';

    banner.innerHTML = `
        <div style="display:flex;align-items:center;gap:8px;padding:10px 14px;background:${color};color:white;border-radius:6px;font-size:13px;max-width:400px;box-shadow:0 4px 12px rgba(0,0,0,0.3);">
            <span style="font-weight:600;">${code}</span>
            <span style="opacity:0.9;">${operation} failed</span>
            <button onclick="dismissErrorBanner('${code}')" style="margin-left:auto;background:var(--border-secondary);border:none;color:white;padding:2px 8px;border-radius:4px;cursor:pointer;font-size:11px;">Dismiss</button>
        </div>
    `;
    banner.style.display = 'block';

    // Auto-dismiss critical after 10s, high after 6s
    const timeout = severity === 'critical' ? 10000 : 6000;
    setTimeout(() => dismissErrorBanner(code), timeout);
}

function showRecoveryToast(operation) {
    const banner = getOrCreateBanner();
    banner.innerHTML = `
        <div style="display:flex;align-items:center;gap:8px;padding:10px 14px;background:#10b981;color:white;border-radius:6px;font-size:13px;max-width:400px;box-shadow:0 4px 12px rgba(0,0,0,0.3);">
            <span style="font-weight:600;">Recovered</span>
            <span style="opacity:0.9;">${operation} stabilized</span>
        </div>
    `;
    banner.style.display = 'block';
    setTimeout(() => { banner.style.display = 'none'; }, 4000);
}

function getOrCreateBanner() {
    let banner = document.getElementById('error-system-banner');
    if (!banner) {
        banner = document.createElement('div');
        banner.id = 'error-system-banner';
        banner.style.cssText = 'position:fixed;top:60px;right:16px;z-index:10000;display:none;';
        document.body.appendChild(banner);
    }
    return banner;
}

function dismissErrorBanner(code) {
    const banner = document.getElementById('error-system-banner');
    if (banner) banner.style.display = 'none';
    resetCircuitBreaker(code);
}

function updateErrorIndicator() {
    // Subtle dot in top bar - only visible when errors exist
    let indicator = document.getElementById('error-indicator');
    if (!indicator) {
        indicator = document.createElement('span');
        indicator.id = 'error-indicator';
        indicator.style.cssText = 'display:inline-block;width:6px;height:6px;background:#f59e0b;border-radius:50%;margin-left:6px;cursor:pointer;';
        indicator.title = 'Errors logged - open settings to view';
        indicator.onclick = () => openSettingsPane('error-log');

        const display = document.getElementById('model-name-display');
        if (display && display.parentNode) {
            display.parentNode.appendChild(indicator);
        }
    }
    indicator.style.display = ErrorSystem.activeErrors.size > 0 ? 'inline-block' : 'none';
}

// ---------------------------------------------------------------------------
// Hardware detection (defensive - detects edge case platforms)
// ---------------------------------------------------------------------------

function detectHardwareEdgeCases() {
    const ua = navigator.userAgent.toLowerCase();
    const memory = navigator.deviceMemory || 4;
    const cores = navigator.hardwareConcurrency || 4;

    // Mobile / ARM detection
    if (/android|iphone|ipad|ipod|windows phone/i.test(ua)) {
        logFrontendError('SYS-HW-MOBILE', 'hardware.detect', new Error('Mobile platform detected'));
    }

    // macOS detection
    if (/macintosh|mac os x/i.test(ua) && !/iphone|ipad|ipod/i.test(ua)) {
        logFrontendError('SYS-HW-MAC', 'hardware.detect', new Error('macOS platform detected'));
    }

    // Potato PC detection (< 4GB RAM or < 4 cores)
    if (memory < 4 || cores < 4) {
        logFrontendError('SYS-HW-POTATO', 'hardware.detect', new Error(`Low spec: ${memory}GB RAM, ${cores} cores`));
    }

    // No GPU detection (canvas test)
    try {
        const canvas = document.createElement('canvas');
        const gl = canvas.getContext('webgl') || canvas.getContext('experimental-webgl');
        if (!gl) {
            logFrontendError('SYS-HW-NO-GPU', 'hardware.detect', new Error('WebGL unavailable'));
        }
    } catch {
        logFrontendError('SYS-HW-NO-GPU', 'hardware.detect', new Error('WebGL test failed'));
    }
}

// ---------------------------------------------------------------------------
// API availability guard
// ---------------------------------------------------------------------------

function guardBackendAPI(apiPath, fallback) {
    const parts = apiPath.split('.');
    let obj = window;
    for (const part of parts) {
        obj = obj?.[part];
        if (obj === undefined || obj === null) {
            logFrontendError('FE-API-MISS', `api.guard.${apiPath}`, new Error(`API ${apiPath} not available`));
            return fallback;
        }
    }
    return obj;
}

// ---------------------------------------------------------------------------
// DOM defensive getter
// ---------------------------------------------------------------------------

function safeGetElement(id, operation) {
    const el = document.getElementById(id);
    if (!el) {
        logFrontendError('FE-DOM-MISS', operation || 'dom.get', new Error(`Missing element: #${id}`));
    }
    return el;
}

// ---------------------------------------------------------------------------
// Input size guard
// ---------------------------------------------------------------------------

function guardInputSize(value, maxChars, operation) {
    if (typeof value === 'string' && value.length > maxChars) {
        logFrontendError('FE-INPUT-OVERSIZE', operation || 'input.guard', new Error(`Input exceeded ${maxChars} chars`));
        return value.slice(0, maxChars);
    }
    return value;
}

// ---------------------------------------------------------------------------
// Image validation
// ---------------------------------------------------------------------------

function validateImageForVision(imageData, operation) {
    if (!imageData || imageData.length < 100) {
        logFrontendError('IMG-VISION-BLANK', operation || 'image.validate', new Error('Blank or missing image'));
        return { valid: false, reason: 'blank' };
    }
    // Rough base64 size estimate (1 char ≈ 0.75 bytes)
    const estBytes = imageData.length * 0.75;
    const maxBytes = 20 * 1024 * 1024; // 20MB vision limit
    if (estBytes > maxBytes) {
        logFrontendError('IMG-VISION-OVERSIZE', operation || 'image.validate', new Error(`Image ~${(estBytes/1024/1024).toFixed(1)}MB exceeds ${maxBytes/1024/1024}MB`));
        return { valid: false, reason: 'oversize' };
    }
    return { valid: true };
}

// ---------------------------------------------------------------------------
// Audio device guard
// ---------------------------------------------------------------------------

async function guardAudioDevice(operation) {
    try {
        if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
            logFrontendError('AUD-MIC-NF', operation || 'audio.guard', new Error('getUserMedia not supported'));
            return false;
        }
        await navigator.mediaDevices.getUserMedia({ audio: true });
        return true;
    } catch (err) {
        logFrontendError('AUD-MIC-NF', operation || 'audio.guard', err);
        return false;
    }
}

// ---------------------------------------------------------------------------
// Manual error logging (frontend-only errors)
// ---------------------------------------------------------------------------

function logFrontendError(code, operation, err) {
    const payload = {
        code: code || 'FE-UNKNOWN',
        severity: 'medium',
        operation: operation || 'frontend',
        time: new Date().toISOString()
    };
    ErrorSystem.activeErrors.set(`${code}-${Date.now()}`, payload);
    updateErrorIndicator();
}

// ---------------------------------------------------------------------------
// Export
// ---------------------------------------------------------------------------

window.ErrorSystem = ErrorSystem;
window.initErrorSystem = initErrorSystem;
window.dismissErrorBanner = dismissErrorBanner;
window.logFrontendError = logFrontendError;
window.guardBackendAPI = guardBackendAPI;
window.safeGetElement = safeGetElement;
window.guardInputSize = guardInputSize;
window.validateImageForVision = validateImageForVision;
window.guardAudioDevice = guardAudioDevice;
window.detectHardwareEdgeCases = detectHardwareEdgeCases;

document.addEventListener('DOMContentLoaded', initErrorSystem);
