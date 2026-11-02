#!/usr/bin/env node
/**
 * tools/shot.js — element-accurate screenshots of the running UI harness.
 *
 * The desktop preview panel captures at a fixed ~480px surface, which is far
 * too small for documentation. This drives headless Edge/Chrome over the
 * DevTools Protocol instead and clips each capture to a specific element at
 * 2x, so tutorial images are pixel-accurate and independent of window size.
 *
 * No npm dependencies: Node 22+ ships a global WebSocket.
 *
 * Usage:  node tools/shot.js tools/shots.json
 *
 * Spec:
 *   {
 *     "url":    "http://127.0.0.1:8799/index.html",
 *     "outDir": "frontend/dist/img/tutorials",
 *     "width":  1600,
 *     "height": 1000,
 *     "shots": [
 *       { "name": "sampler-sheet", "setup": "toggleFlyout('sampler-settings-flyout')",
 *         "selector": "#sampler-settings-flyout", "scale": 2 }
 *     ]
 *   }
 */

'use strict';

const { spawn } = require('child_process');
const fs = require('fs');
const os = require('os');
const path = require('path');

const BROWSERS = [
    'C:/Program Files (x86)/Microsoft/Edge/Application/msedge.exe',
    'C:/Program Files/Microsoft/Edge/Application/msedge.exe',
    'C:/Program Files/Google/Chrome/Application/chrome.exe',
    'C:/Program Files (x86)/Google/Chrome/Application/chrome.exe'
];

const sleep = ms => new Promise(r => setTimeout(r, ms));

function findBrowser() {
    const found = BROWSERS.find(p => fs.existsSync(p));
    if (!found) throw new Error('No Edge/Chrome binary found');
    return found;
}

async function waitForTarget(port, timeoutMs = 20000) {
    const deadline = Date.now() + timeoutMs;
    while (Date.now() < deadline) {
        try {
            const res = await fetch(`http://127.0.0.1:${port}/json/list`);
            const targets = await res.json();
            const page = targets.find(t => t.type === 'page' && t.webSocketDebuggerUrl);
            if (page) return page;
        } catch (e) { /* browser not up yet */ }
        await sleep(250);
    }
    throw new Error('DevTools target never appeared');
}

class CDP {
    constructor(ws) {
        this.ws = ws;
        this.nextId = 1;
        this.pending = new Map();
        this.listeners = new Set();
        ws.addEventListener('message', ev => {
            const msg = JSON.parse(ev.data);
            if (msg.id && this.pending.has(msg.id)) {
                const { resolve, reject } = this.pending.get(msg.id);
                this.pending.delete(msg.id);
                msg.error ? reject(new Error(JSON.stringify(msg.error))) : resolve(msg.result);
            } else if (msg.method) {
                this.listeners.forEach(fn => fn(msg.method, msg.params));
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

    on(fn) { this.listeners.add(fn); }

    async eval(expression) {
        const result = await this.send('Runtime.evaluate', {
            expression,
            awaitPromise: true,
            returnByValue: true
        });
        if (result.exceptionDetails) {
            throw new Error('Page error: ' + (result.exceptionDetails.text || '') +
                ' ' + (result.exceptionDetails.exception?.description || ''));
        }
        return result.result?.value;
    }
}

async function main() {
    const specPath = process.argv[2];
    if (!specPath) throw new Error('Usage: node tools/shot.js <spec.json>');
    const spec = JSON.parse(fs.readFileSync(specPath, 'utf8'));

    const url = spec.url || 'http://127.0.0.1:8799/index.html';
    const outDir = spec.outDir || 'frontend/dist/img/tutorials';
    const width = spec.width || 1600;
    const height = spec.height || 1000;
    const port = 9200 + Math.floor(Math.random() * 300);
    const profile = fs.mkdtempSync(path.join(os.tmpdir(), 'howl-shot-'));

    fs.mkdirSync(outDir, { recursive: true });

    const browser = spawn(findBrowser(), [
        '--headless=new',
        '--disable-gpu',
        '--hide-scrollbars',
        '--no-first-run',
        '--no-default-browser-check',
        '--force-device-scale-factor=1',
        `--window-size=${width},${height}`,
        `--remote-debugging-port=${port}`,
        `--user-data-dir=${profile}`,
        'about:blank'
    ], { stdio: 'ignore' });

    let ws, cdp;
    const cleanup = () => {
        try { ws?.close(); } catch (e) {}
        try { browser.kill(); } catch (e) {}
    };

    try {
        const target = await waitForTarget(port);
        ws = new WebSocket(target.webSocketDebuggerUrl);
        await new Promise((resolve, reject) => {
            ws.addEventListener('open', resolve, { once: true });
            ws.addEventListener('error', reject, { once: true });
        });
        cdp = new CDP(ws);

        await cdp.send('Page.enable');
        await cdp.send('Runtime.enable');
        await cdp.send('Emulation.setDeviceMetricsOverride', {
            width, height, deviceScaleFactor: 1, mobile: false
        });

        const loaded = new Promise(resolve => {
            cdp.on((method) => { if (method === 'Page.loadEventFired') resolve(); });
        });
        await cdp.send('Page.navigate', { url });
        await Promise.race([loaded, sleep(15000)]);
        await sleep(1200);

        // Kill transitions once: opening a sheet must not be caught mid-animation.
        await cdp.eval(`
            (() => {
                if (document.getElementById('_shot_style')) return;
                const s = document.createElement('style');
                s.id = '_shot_style';
        s.textContent = '*{transition:none!important;animation:none!important;scroll-behavior:auto!important}' + '[class*=toast],[id*=toast],[class*=notification],[id*=notification],[role=status],.snackbar,.banner{display:none!important}';
                document.head.appendChild(s);
            })()
        `);

        const report = [];
        for (const shot of spec.shots) {
            if (shot.width || shot.height) {
                await cdp.send('Emulation.setDeviceMetricsOverride', {
                    width: shot.width || spec.width || 1600,
                    height: shot.height || spec.height || 1000,
                    deviceScaleFactor: 1,
                    mobile: false
                });
                await sleep(450);
            }

            if (shot.setup) {
                await cdp.eval(`(async () => { ${shot.setup} })()`);
                await sleep(shot.settle || 400);
            }

            const rect = await cdp.eval(`
                (() => {
                    const el = document.querySelector(${JSON.stringify(shot.selector)});
                    if (!el) return null;
                    const r = el.getBoundingClientRect();
                    return { x: r.x + window.scrollX, y: r.y + window.scrollY, width: r.width, height: r.height };
                })()
            `);
            if (!rect || !rect.width) throw new Error('Selector not found or zero-size: ' + shot.selector);

            const pad = shot.pad == null ? 0 : shot.pad;
            const { data } = await cdp.send('Page.captureScreenshot', {
                format: 'png',
                captureBeyondViewport: true,
                clip: {
                    x: Math.max(0, rect.x - pad),
                    y: Math.max(0, rect.y - pad),
                    width: rect.width + pad * 2,
                    height: rect.height + pad * 2,
                    scale: shot.scale || spec.scale || 2
                }
            });

            const file = path.join(outDir, shot.name + '.png');
            fs.writeFileSync(file, Buffer.from(data, 'base64'));
            const bytes = fs.statSync(file).size;
            report.push(`${shot.name}.png  ${Math.round(rect.width)}x${Math.round(rect.height)} css  ${bytes} bytes`);
        }

        console.log(report.join('\n'));
    } finally {
        cleanup();
    }
}

main().catch(err => {
    console.error('shot failed:', err.message);
    process.exit(1);
});
