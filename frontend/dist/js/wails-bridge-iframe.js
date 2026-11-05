/**
 * wails-bridge-iframe.js — iframe-side audio-settings bridge client.
 *
 * Loaded as the first script in audio-settings.html, before settings-state.js.
 * Installs a window.go shim ONLY when a native window.go is absent, so a future
 * Wails version that injects into frames wins automatically and the standalone
 * file:// spec is unaffected.
 *
 * This file owns only the panel's request handles (sequence, promise/timers,
 * retired ids) and its presentation state. It never owns the activation lock,
 * never holds the real backend promise, and never decides whether an activation
 * may start.
 */

(function () {
    'use strict';

    if (window.go) return;                    // native injection wins
    if (window.__howlBridgeIframe) return;

    var METHODS = [
        'GetAudioStatus',
        'ApplyAudioSettings',
        'GetTTSModels',
        'OpenDirectoryPicker',
        'ScanTTSDir',
        'GetSettings',
        'SetSettings'
    ];

    var BUDGETS = {
        GetSettings: 20000,
        GetAudioStatus: 20000,
        GetTTSModels: 20000,
        SetSettings: 60000,
        ScanTTSDir: 120000,
        ApplyAudioSettings: 180000,
        OpenDirectoryPicker: 300000
    };

    var HANDSHAKE_TIMEOUT = 5000;

    var parentOrigin = window.location.origin;
    var TARGET = (parentOrigin && parentOrigin !== 'null') ? parentOrigin : '*';

    var bridge = {
        mode: 'connecting',
        generation: null,
        seq: 0,
        pending: new Map(),
        retired: new Set(),
        queue: [],
        lateReplies: 0,
        rejects: 0,
        reason: ''
    };

    function failure(name, message) {
        var err = new Error(message);
        err.name = name;
        return err;
    }

    function post(msg) {
        try {
            window.parent.postMessage(msg, TARGET);
        } catch (e) { /* no parent */ }
    }

    function dispatch(entry) {
        entry.id = String(bridge.generation) + ':' + (++bridge.seq);
        entry.queued = false;
        bridge.pending.set(entry.id, entry);
        post({ __howlBridge: 1, type: 'call', id: entry.id, method: entry.method, args: entry.args });
    }

    function call(method, args) {
        return new Promise(function (resolve, reject) {
            if (bridge.mode === 'unavailable') {
                bridge.rejects++;
                reject(failure('BridgeUnavailableError', bridge.reason || 'backend bridge unavailable'));
                return;
            }
            if (bridge.mode === 'closed') {
                bridge.rejects++;
                reject(failure('BridgeClosedError', 'panel closed'));
                return;
            }

            var budget = BUDGETS[method] || 20000;
            var entry = {
                method: method,
                args: args,
                resolve: resolve,
                reject: reject,
                id: null,
                queued: bridge.mode !== 'ready',
                timer: null
            };

            entry.timer = setTimeout(function () {
                if (entry.queued) {
                    var i = bridge.queue.indexOf(entry);
                    if (i >= 0) bridge.queue.splice(i, 1);
                } else if (entry.id) {
                    bridge.pending.delete(entry.id);
                    bridge.retired.add(entry.id);
                }
                // Retired, never reversed: a late reply cannot resolve it.
                entry.reject(failure('BridgeTimeoutError',
                    'no answer within ' + budget + 'ms; outcome indeterminate'));
            }, budget);

            if (entry.queued) bridge.queue.push(entry);
            else dispatch(entry);
        });
    }

    function onReady(msg) {
        if (bridge.mode === 'ready') return;
        bridge.generation = (typeof msg.generation === 'number') ? msg.generation : 0;
        bridge.mode = 'ready';
        clearTimeout(handshakeTimer);

        var queued = bridge.queue.splice(0, bridge.queue.length);
        queued.forEach(dispatch);

        window.dispatchEvent(new CustomEvent('howl:bridge-ready', {
            detail: { generation: bridge.generation }
        }));
    }

    function onMessage(event) {
        if (event.source !== window.parent) return;
        if (parentOrigin && parentOrigin !== 'null' &&
            event.origin && event.origin !== 'null' && event.origin !== parentOrigin) {
            return;
        }
        var data = event.data;
        if (!data || data.__howlBridge !== 1) return;

        if (data.type === 'ready') { onReady(data); return; }

        if (data.type === 'activation-settled') {
            window.dispatchEvent(new CustomEvent('howl:activation-settled'));
            return;
        }

        if (data.type === 'result') {
            var entry = data.id ? bridge.pending.get(data.id) : null;
            if (!entry) { bridge.lateReplies++; return; }   // retired or unknown id
            bridge.pending.delete(data.id);
            clearTimeout(entry.timer);
            if (data.ok) {
                entry.resolve(data.result);
            } else {
                var named = (data.error && data.error.name) || 'BridgeBackendError';
                entry.reject(failure(named, (data.error && data.error.message) || 'backend error'));
            }
        }
    }

    var handshakeTimer = setTimeout(function () {
        if (bridge.mode === 'ready') return;
        bridge.mode = 'unavailable';
        bridge.reason = 'no host handshake within ' + HANDSHAKE_TIMEOUT + 'ms';
        var queued = bridge.queue.splice(0, bridge.queue.length);
        queued.forEach(function (entry) {
            clearTimeout(entry.timer);
            entry.reject(failure('BridgeUnavailableError', bridge.reason));
        });
        window.dispatchEvent(new CustomEvent('howl:bridge-unavailable', {
            detail: { reason: bridge.reason }
        }));
    }, HANDSHAKE_TIMEOUT);

    window.addEventListener('message', onMessage);

    window.addEventListener('pagehide', function () {
        bridge.mode = 'closed';
        bridge.pending.forEach(function (entry) {
            clearTimeout(entry.timer);
            entry.reject(failure('BridgeClosedError', 'panel closed'));
        });
        bridge.pending.clear();
        bridge.queue.splice(0, bridge.queue.length).forEach(function (entry) {
            clearTimeout(entry.timer);
            entry.reject(failure('BridgeClosedError', 'panel closed'));
        });
    });

    var App = {};
    METHODS.forEach(function (name) {
        App[name] = function () {
            return call(name, Array.prototype.slice.call(arguments));
        };
    });

    window.go = { main: { App: App } };

    Object.defineProperty(window, '__howlBridgeIframe', {
        value: bridge,
        enumerable: false,
        configurable: false
    });

    post({ __howlBridge: 1, type: 'hello' });
})();
