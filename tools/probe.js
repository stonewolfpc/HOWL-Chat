#!/usr/bin/env node
/**
 * tools/probe.js — measure panel geometry at several window sizes.
 *
 * CSS viewport units reflow live, but "it looks fine" is not evidence. This
 * loads the harness in headless Edge at a series of window sizes, opens each
 * panel, and prints the measured panel size, the page's horizontal overflow
 * and the resolved column count, so full-screen behaviour can be checked at
 * every resolution rather than assumed from one screenshot.
 *
 * Usage: node tools/probe.js tools/probes.json
 */

'use strict';

const { spawn } = require('child_process');
const fs = require('fs');
const os = require('os');
const path = require('path');

const BROWSERS = [
    'C:/Program Files (x86)/Microsoft/Edge/Application/msedge.exe',
    'C:/Program Files/Microsoft/Edge/Application/msedge.exe',
    'C:/Program Files/Google/Chrome/Application/chrome.exe'
];

const sleep = ms => new Promise(r => setTimeout(r, ms));

async function waitForTarget(port, timeoutMs = 20000) {
    const deadline = Date.now() + timeoutMs;
    while (Date.now() < deadline) {
        try {
            const res = await fetch(`http://127.0.0.1:${port}/json/list`);
            const targets = await res.json();
            const page = targets.find(t => t.type === 'page' && t.webSocketDebuggerUrl);
            if (page) return page;
        } catch (e) { /* not up yet */ }
        await sleep(250);
    }
    throw new Error('DevTools target never appeared');
}

class CDP {
    constructor(ws) {
        this.ws = ws;
        this.nextId = 1;
        this.pending = new Map();
        this.loadWaiters = [];
        ws.addEventListener('message', ev => {
            const msg = JSON.parse(ev.data);
            if (msg.id && this.pending.has(msg.id)) {
                const { resolve, reject } = this.pending.get(msg.id);
                this.pending.delete(msg.id);
                msg.error ? reject(new Error(JSON.stringify(msg.error))) : resolve(msg.result);
            } else if (msg.method === 'Page.loadEventFired') {
                this.loadWaiters.splice(0).forEach(fn => fn());
            }
        });
    }

    send(method, params = {}) {
        const id = this.nextId++;
        return new Promise((resolve, reject) => {
            this.pending.set(id, { resolve, reject });
            this.ws.send(JSON.stringify({ id, method, params }));
            setTimeout(() => {
                if (this.pending.has(id)) {
                    this.pending.delete(id);
                    reject(new Error('CDP timeout: ' + method));
                }
            }, 20000);
        });
    }

    loaded() {
        return new Promise(resolve => {
            this.loadWaiters.push(resolve);
            setTimeout(resolve, 15000);
        });
    }

    async eval(expression) {
        const result = await this.send('Runtime.evaluate', {
            expression, awaitPromise: true, returnByValue: true
        });
        if (result.exceptionDetails) {
            throw new Error('Page error: ' + (result.exceptionDetails.exception?.description || result.exceptionDetails.text));
        }
        return result.result?.value;
    }
}

async function main() {
    const spec = JSON.parse(fs.readFileSync(process.argv[2], 'utf8'));
    const url = spec.url;
    const port = 9500 + Math.floor(Math.random() * 300);
    const profile = fs.mkdtempSync(path.join(os.tmpdir(), 'howl-probe-'));
    const browser = spawn(BROWSERS.find(p => fs.existsSync(p)), [
        '--headless=new', '--disable-gpu', '--hide-scrollbars', '--no-first-run',
        `--window-size=1600,1000`, `--remote-debugging-port=${port}`,
        `--user-data-dir=${profile}`, 'about:blank'
    ], { stdio: 'ignore' });

    let ws;
    const cleanup = () => {
        try { ws?.close(); } catch (e) {}
        try { browser.kill(); } catch (e) {}
    };

    try {
        const target = await waitForTarget(port);
        ws = new WebSocket(target.webSocketDebuggerUrl);
        await new Promise((res, rej) => {
            ws.addEventListener('open', res, { once: true });
            ws.addEventListener('error', rej, { once: true });
        });
        const cdp = new CDP(ws);
        await cdp.send('Page.enable');
        await cdp.send('Runtime.enable');

        const rows = [];
        for (const check of spec.checks) {
            await cdp.send('Emulation.setDeviceMetricsOverride', {
                width: check.width, height: check.height, deviceScaleFactor: 1, mobile: false
            });
            const loaded = cdp.loaded();
            await cdp.send('Page.navigate', { url });
            await loaded;
            await sleep(900);
            await cdp.eval("new Promise(r => { const done = () => (document.readyState === 'complete' && document.styleSheets.length > 5) ? r(true) : setTimeout(done, 40); done(); })");

            await cdp.eval(`
                (() => {
                    if (document.getElementById('_probe_style')) return;
                    const s = document.createElement('style');
                    s.id = '_probe_style';
                    s.textContent = '*{transition:none!important;animation:none!important}' +
                        '[class*=toast],[id*=toast],[class*=notification],[id*=notification]{display:none!important}';
                    document.head.appendChild(s);
                })()
            `);

            if (check.setup) await cdp.eval(`(async () => { ${check.setup} })()`);
            await sleep(check.settle || 350);

            let out, prevOut;
            for (let attempt = 0; attempt < 8; attempt++) {
                out = await cdp.eval(check.js);
                if (prevOut && JSON.stringify(prevOut) === JSON.stringify(out)) break;
                prevOut = out;
                await sleep(200);
            }
            rows.push({ name: check.name, size: `${check.width}x${check.height}`, ...out });
        }

        for (const r of rows) {
            const geo = r.panel ? 'panel ' + r.panel[0] + 'x' + r.panel[1] : '';
            const extras = Object.keys(r)
                .filter(k => !['name', 'size', 'pane', 'panel', 'overflowX'].includes(k))
                .map(k => k + '=' + r[k]).join(' ');
            const overflow = r.overflowX ? ' OVERFLOW-X ' + r.overflowX + 'px' : ' overflow=0';
            console.log(r.name.padEnd(12) + ' ' + r.size.padEnd(9) + ' ' + geo.padEnd(18) + ' ' + extras.padEnd(30) + overflow);
        }
        // A 'pane' check measures a sub-panel, so viewport width does not apply.
        const bad = rows.filter(r => r.overflowX > 0 || (r.panel && r.pane !== true && Math.abs(r.panel[0] - parseInt(r.size.split('x')[0], 10)) > 1));
        console.log(bad.length ? `\nFAIL: ${bad.map(b => b.name + '@' + b.size).join(', ')}` : '\nPASS: every panel matched its viewport with no horizontal overflow');
        process.exitCode = bad.length ? 1 : 0;
    } finally {
        cleanup();
    }
}

main().catch(e => { console.error('probe failed:', e.message); process.exit(1); });
