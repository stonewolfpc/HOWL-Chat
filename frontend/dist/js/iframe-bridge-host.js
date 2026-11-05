/**
 * iframe-bridge-host.js — parent-side owner of the audio-settings bridge.
 *
 * Loaded by every document that can open the audio-settings modal, before
 * menu-functions.js. Owns everything that must outlive the panel:
 *   - the generation counter and the per-modal message listener
 *   - in-flight requests, timeouts and retired ids
 *   - the activation lock and the real ApplyAudioSettings promise
 *
 * The iframe owns only its own request handles and presentation state; it never
 * holds the lock and never decides whether an activation may start.
 *
 * Trust rules (see build/unit-c/d1-plan.md 3.2.1):
 *   - production: event.source must be the attached window AND the origin must
 *     equal this document's non-opaque origin. Opaque/empty origins fail closed.
 *   - file:// (standalone test path only): source identity is authoritative.
 */

(function () {
    'use strict';

    if (window.howlBridgeHost) return;

    var ALLOWED = new Set([
        'GetAudioStatus',
        'ApplyAudioSettings',
        'GetTTSModels',
        'OpenDirectoryPicker',
        'ScanTTSDir',
        'GetSettings',
        'SetSettings'
    ]);

    var BUDGETS = {
        GetSettings: 20000,
        GetAudioStatus: 20000,
        GetTTSModels: 20000,
        SetSettings: 60000,
        ScanTTSDir: 120000,
        ApplyAudioSettings: 180000,
        OpenDirectoryPicker: 300000
    };

    var STORAGE_KEY = 'howl.audio.activation';

    var generation = 0;
    var session = null;
    var listenerBound = false;

    var activation = { held: false, id: null, startedAt: 0 };

    var stats = {
        generation: 0,
        rejected: 0,
        lateReply: 0,
        unknownMethod: 0,
        timeouts: 0,
        calls: 0,
        results: 0,
        refusedActivations: 0
    };

    // A lock mirrored across same-origin navigation. It is a UI hint: the
    // authoritative serialization is the backend's applyMu.
    try {
        var restored = JSON.parse(sessionStorage.getItem(STORAGE_KEY) || 'null');
        if (restored && restored.activationId) {
            // A restored marker is not proof that the activation is still
            // running, and it is never released in this document: a call that
            // has not yet reached Go cannot be ruled out. Recovery is a full
            // application restart.
            activation = {
                held: true,
                id: restored.activationId,
                startedAt: restored.startedAt || 0,
                restored: true
            };
        }
    } catch (e) { /* storage unavailable */ }

    function writeMarker() {
        try {
            if (activation.held) {
                sessionStorage.setItem(STORAGE_KEY, JSON.stringify({
                    activationId: activation.id,
                    startedAt: activation.startedAt
                }));
            } else {
                sessionStorage.removeItem(STORAGE_KEY);
            }
        } catch (e) { /* storage unavailable */ }
    }

    function isProductionOrigin() {
        return window.location.protocol !== 'file:';
    }

    function originOk(event, sess) {
        if (!sess || event.source !== sess.win) return false;
        if (!isProductionOrigin()) return true;          // test path only
        var origin = window.location.origin;
        if (!origin || origin === 'null') return false;  // fail closed
        return event.origin === origin;
    }

    function postTo(sess, msg) {
        if (!sess || !sess.active || !sess.win) return;
        try {
            sess.win.postMessage(msg, isProductionOrigin() ? window.location.origin : '*');
        } catch (e) { /* window destroyed */ }
    }

    function beginActivation(id) {
        if (activation.held) return false;
        activation = { held: true, id: id, startedAt: Date.now(), restored: false };
        writeMarker();
        return true;
    }

    function endActivation() {
        // A restored lock has no promise in this document. Only the live
        // document that issued an activation may release it.
        if (activation.restored) return;
        var wasHeld = activation.held;
        activation = { held: false, id: null, startedAt: 0 };
        writeMarker();
        // Uncorrelated notification only: no id, no result, no error (3.7).
        if (wasHeld && session && session.active) {
            postTo(session, { __howlBridge: 1, type: 'activation-settled' });
        }
    }

    function reply(sess, id, ok, payload) {
        if (ok) {
            postTo(sess, { __howlBridge: 1, type: 'result', id: id, ok: true, result: payload });
        } else {
            postTo(sess, {
                __howlBridge: 1,
                type: 'result',
                id: id,
                ok: false,
                error: {
                    name: (payload && payload.name) || 'BridgeBackendError',
                    message: String((payload && payload.message) || payload || 'backend error')
                }
            });
        }
    }

    function handleCall(sess, data) {
        var id = typeof data.id === 'string' ? data.id : '';
        if (!id) return;

        if (sess.pending.has(id)) {
            reply(sess, id, false, 'duplicate request id');
            return;
        }
        if (!ALLOWED.has(data.method)) {
            stats.unknownMethod++;
            reply(sess, id, false, 'method not allowed: ' + String(data.method));
            return;
        }

        var method = data.method;
        var args = Array.isArray(data.args) ? data.args : [];
        var budget = BUDGETS[method] || 20000;
        var isActivation = method === 'ApplyAudioSettings';

        // The lock is enforcement, not an indicator: refuse before a pending
        // entry, a timer or any backend call exists.
        if (isActivation && activation.held) {
            stats.refusedActivations++;
            reply(sess, id, false, {
                name: 'ActivationLockedError',
                message: 'an activation is already in progress; this request did not reach the backend'
            });
            return;
        }

        var entry = { id: id, method: method, timer: null };
        sess.pending.set(id, entry);
        stats.calls++;

        entry.timer = setTimeout(function () {
            if (!sess.pending.has(id)) return;
            sess.pending.delete(id);
            sess.retired.add(id);
            stats.timeouts++;
            // The backend call is NOT cancelled; the activation lock stays held
            // until the real promise settles.
            reply(sess, id, false, 'no answer within ' + budget + 'ms; outcome indeterminate');
        }, budget);

        if (isActivation && !beginActivation(id)) {
            clearTimeout(entry.timer);
            sess.pending.delete(id);
            reply(sess, id, false, {
                name: 'ActivationLockedError',
                message: 'an activation is already in progress; this request did not reach the backend'
            });
            return;
        }

        var settle = function (ok, payload) {
            clearTimeout(entry.timer);
            var live = session === sess && sess.active && sess.pending.has(id);
            if (live) {
                sess.pending.delete(id);
                stats.results++;
                reply(sess, id, ok, ok ? payload : (payload && payload.message) || payload);
            } else {
                // Retired or closed generation: drop the payload entirely (3.7).
                stats.lateReply++;
            }
            if (isActivation) endActivation();
        };

        var app = window.go && window.go.main && window.go.main.App;
        if (!app || typeof app[method] !== 'function') {
            settle(false, 'backend binding unavailable: ' + method);
            return;
        }

        var promise;
        try {
            promise = app[method].apply(null, args);
        } catch (err) {
            settle(false, err);
            return;
        }
        Promise.resolve(promise).then(function (result) {
            settle(true, result);
        }, function (err) {
            settle(false, err);
        });
    }

    function onMessage(event) {
        var sess = session;
        if (!sess) return;
        if (event.source !== sess.win) { stats.rejected++; return; }
        if (!originOk(event, sess)) { stats.rejected++; return; }

        var data = event.data;
        if (!data || data.__howlBridge !== 1) return;

        if (data.type === 'hello') {
            postTo(sess, { __howlBridge: 1, type: 'ready', generation: sess.gen });
            return;
        }
        if (data.type === 'call') handleCall(sess, data);
    }

    function attach(iframe) {
        if (!iframe || !iframe.contentWindow) return null;
        if (session) detach();

        generation += 1;
        session = {
            gen: generation,
            win: iframe.contentWindow,
            pending: new Map(),
            retired: new Set(),
            active: true
        };
        if (!listenerBound) {
            window.addEventListener('message', onMessage);
            listenerBound = true;
        }
        stats.generation = generation;
        return generation;
    }

    function detach(iframe) {
        var sess = session;
        if (!sess) return;
        if (iframe && iframe.contentWindow && iframe.contentWindow !== sess.win) return;

        sess.active = false;
        sess.pending.forEach(function (entry) {
            clearTimeout(entry.timer);
            sess.retired.add(entry.id);
        });
        sess.pending.clear();
        session = null;
        generation += 1;
        stats.generation = generation;

        if (listenerBound) {
            window.removeEventListener('message', onMessage);
            listenerBound = false;
        }
    }

    window.howlBridgeHost = {
        attach: attach,
        detach: detach,
        isActivationLocked: function () { return activation.held === true; },
        activation: function () {
            return {
                held: activation.held,
                id: activation.id,
                startedAt: activation.startedAt,
                restored: activation.restored === true
            };
        },
        get generation() { return generation; },
        get listening() { return listenerBound; },
        __stats: stats
    };
})();
