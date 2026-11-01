/*
 * Wails runtime shim for the UI harness.
 *
 * The packaged app runs inside WebView2, where Wails injects window.go and
 * window.runtime before any page script executes. In a normal browser those
 * do not exist, so every binding call and event subscription would throw.
 *
 * This shim provides the same shapes, proxying each binding call to the harness
 * HTTP endpoint, which invokes the real Go method by reflection. Events are
 * delivered in-page instead of across the native bridge, which is enough for
 * anything the frontend emits to itself; backend-originated events (streaming
 * tokens) have no channel here and simply never arrive.
 */
(function () {
    'use strict';

    if (window.__HOWL_HARNESS__) return;
    window.__HOWL_HARNESS__ = true;

    function rpc(method, args) {
        return fetch('/__rpc', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ method: method, args: args })
        }).then(function (res) {
            return res.text().then(function (text) {
                var payload = null;
                try { payload = JSON.parse(text); } catch (e) { /* non-JSON body */ }
                if (!res.ok) {
                    var message = (payload && payload.error) || text || ('HTTP ' + res.status);
                    console.error('[harness] ' + method + ' failed: ' + message);
                    throw new Error(message);
                }
                return payload ? payload.result : null;
            });
        });
    }

    var appProxy = new Proxy({}, {
        get: function (_target, prop) {
            if (typeof prop !== 'string') return undefined;
            if (prop === 'then' || prop === 'toJSON') return undefined;
            return function () {
                return rpc(prop, Array.prototype.slice.call(arguments));
            };
        }
    });

    window.go = { main: { App: appProxy } };

    var listeners = {};
    // Backend-originated events (streamed chunks, reasoning, completion) arrive
    // over Server-Sent Events and are dispatched exactly like Wails EventsOn
    // callbacks, so the real frontend handlers run. Reconnects automatically.
    if (typeof EventSource !== 'undefined') {
        var stream = new EventSource('/__events');
        stream.onmessage = function (e) {
            var msg;
            try { msg = JSON.parse(e.data); } catch (err) { return; }
            if (!msg || !msg.name) return;
            (listeners[msg.name] || []).forEach(function (cb) {
                try { cb.apply(null, msg.args || []); } catch (err) { console.error('[harness] listener ' + msg.name, err); }
            });
        };
        stream.onerror = function () { /* EventSource retries on its own */ };
    }

    window.runtime = {
        EventsOn: function (name, cb) {
            (listeners[name] = listeners[name] || []).push(cb);
        },
        EventsOnMultiple: function (name, cb) { window.runtime.EventsOn(name, cb); },
        EventsOnce: function (name, cb) {
            var wrapped = function () { cb.apply(null, arguments); window.runtime.EventsOff(name); };
            window.runtime.EventsOn(name, wrapped);
        },
        EventsOff: function (name) { delete listeners[name]; },
        EventsEmit: function (name) {
            var args = Array.prototype.slice.call(arguments, 1);
            (listeners[name] || []).forEach(function (cb) { cb.apply(null, args); });
        },
        LogPrint: function (m) { console.log('[wails]', m); },
        LogTrace: function (m) { console.debug('[wails]', m); },
        LogDebug: function (m) { console.debug('[wails]', m); },
        LogInfo: function (m) { console.info('[wails]', m); },
        LogWarning: function (m) { console.warn('[wails]', m); },
        LogError: function (m) { console.error('[wails]', m); },
        LogFatal: function (m) { console.error('[wails]', m); },
        WindowShow: function () {},
        WindowHide: function () {},
        WindowCenter: function () {},
        WindowReload: function () { location.reload(); },
        WindowSetTitle: function (t) { document.title = t; },
        WindowSetSize: function () {},
        WindowGetSize: function () { return Promise.resolve({ w: 1400, h: 900 }); },
        WindowSetPosition: function () {},
        WindowGetPosition: function () { return Promise.resolve({ x: 0, y: 0 }); },
        WindowMaximise: function () {},
        WindowUnmaximise: function () {},
        WindowIsMaximised: function () { return Promise.resolve(false); },
        WindowFullscreen: function () {},
        WindowUnfullscreen: function () {},
        WindowIsFullscreen: function () { return Promise.resolve(false); },
        WindowMinimise: function () {},
        WindowSetAlwaysOnTop: function () {},
        BrowserOpenURL: function (u) { window.open(u, '_blank', 'noopener'); },
        ClipboardSetText: function () { return Promise.resolve(); },
        ClipboardGetText: function () { return Promise.resolve(''); },
        Environment: function () {
            return Promise.resolve({ buildType: 'harness', platform: 'windows', arch: 'amd64' });
        },
        ScreenGetAll: function () { return Promise.resolve([]); },
        Quit: function () { console.warn('[harness] Quit ignored'); },
        OnMultipleInstances: function () {}
    };

    window.wails = {
        EventsEmit: window.runtime.EventsEmit,
        EventsOn: window.runtime.EventsOn,
        EventsOff: window.runtime.EventsOff
    };

    console.log('[harness] Wails shim active — bindings proxy to the Go backend');
})();
