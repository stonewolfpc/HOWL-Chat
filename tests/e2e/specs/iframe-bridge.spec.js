/**
 * D1 iframe-to-host bridge coverage.
 *
 * Serves frontend/dist over a real HTTP origin (production trust path: origin
 * equality, not the file:// exception), with a fixture page that owns the fake
 * backend and the host while the audio-settings panel runs in an iframe.
 */

const { test, expect } = require('@playwright/test');
const http = require('http');
const fs = require('fs');
const path = require('path');

const DIST = path.resolve(__dirname, '../../../frontend/dist');

const FIXTURE = `<!doctype html>
<html><body>
<script>
window.__be = { calls: [], applied: null, inProgress: false, fail: false, delay: 0, publishApplied: true };
window.go = { main: { App: {
  GetSettings: async () => ({ audio: { asrEnabled: true, asrEngine: 'whisper-base', ttsEngine: 'piper' } }),
  SetSettings: async () => {},
  GetAudioStatus: async () => Object.assign({
    asr_readiness: 'local-initialized',
    asr_ready: true,
    tts_ready: false,
    apply_in_progress: window.__be.inProgress,
    applied: window.__be.applied
  }, window.__be.status || {}),
  ApplyAudioSettings: async (...args) => {
    window.__be.calls.push(args);
    if (window.__be.delay) await new Promise(r => setTimeout(r, window.__be.delay));
    if (window.__be.fail) throw new Error('activation failed');
    if (window.__be.publishApplied) {
      window.__be.applied = {
        asr_enabled: !!args[0], tts_enabled: !!args[1],
        asr_engine: args[2], tts_engine: args[3]
      };
    }
  },
  GetTTSModels: async () => [{ name: 'voice-a', size_mb: 12 }],
  OpenDirectoryPicker: async () => 'C:/voices',
  ScanTTSDir: async (dir) => { window.__be.scanned = dir; }
} } };
</script>
<script src="/js/iframe-bridge-host.js"></script>
<iframe id="f" src="/audio-settings.html#general" style="width:900px;height:620px"></iframe>
<script>
window.howlBridgeHost.attach(document.getElementById('f'));
</script>
</body></html>`;

let server;
let base;

test.beforeAll(async () => {
    server = http.createServer((req, res) => {
        const url = req.url.split('?')[0];
        if (url === '/__fixture.html') {
            res.setHeader('content-type', 'text/html');
            res.end(FIXTURE);
            return;
        }
        const file = path.join(DIST, url === '/' ? 'index.html' : url);
        if (!file.startsWith(DIST)) { res.statusCode = 403; res.end(); return; }
        fs.readFile(file, (err, buf) => {
            if (err) { res.statusCode = 404; res.end('not found'); return; }
            res.setHeader('content-type', url.endsWith('.js') ? 'text/javascript' : 'text/html');
            res.end(buf);
        });
    });
    await new Promise(resolve => server.listen(0, '127.0.0.1', resolve));
    base = 'http://127.0.0.1:' + server.address().port;
});

test.afterAll(async () => { if (server) await new Promise(r => server.close(r)); });

async function openFixture(page) {
    await page.goto(base + '/__fixture.html', { waitUntil: 'load' });
    await expect.poll(
        async () => page.frames().map(f => f.url()).join(' '),
        { timeout: 15000 }
    ).toContain('audio-settings.html');
    const panel = page.frames().find(f => f.url().includes('audio-settings.html'));
    expect(panel).toBeTruthy();
    await expect.poll(async () => panel.evaluate(() => window.__howlBridgeIframe?.mode)).toBe('ready');
    await expect(panel.locator('#audio-settings-save')).toBeEnabled({ timeout: 10000 });
    return panel;
}

test('handshake completes and all seven authorized methods work through the bridge', async ({ page }) => {
    const panel = await openFixture(page);

    const generation = await page.evaluate(() => window.howlBridgeHost.generation);
    expect(generation).toBeGreaterThan(0);
    expect(await panel.evaluate(() => window.__howlBridgeIframe.generation)).toBe(generation);

    const settings = await panel.evaluate(() => window.go.main.App.GetSettings());
    expect(settings.audio.asrEngine).toBe('whisper-base');

    await panel.evaluate(() => window.go.main.App.SetSettings({ audio: { asrEnabled: true } }));
    const status = await panel.evaluate(() => window.go.main.App.GetAudioStatus());
    expect(status.asr_readiness).toBe('local-initialized');

    await panel.evaluate(() => window.go.main.App.ApplyAudioSettings(true, false, 'whisper-base', 'piper'));
    expect(await page.evaluate(() => window.__be.calls.length)).toBe(1);

    const voices = await panel.evaluate(() => window.go.main.App.GetTTSModels());
    expect(voices).toHaveLength(1);
    expect(await panel.evaluate(() => window.go.main.App.OpenDirectoryPicker())).toBe('C:/voices');
    await panel.evaluate(() => window.go.main.App.ScanTTSDir('C:/voices'));
    expect(await page.evaluate(() => window.__be.scanned)).toBe('C:/voices');
});

test('unlisted methods are refused instead of silently succeeding', async ({ page }) => {
    const panel = await openFixture(page);

    const reply = await panel.evaluate(async () => {
        const id = window.__howlBridgeIframe.generation + ':999';
        const wait = new Promise(resolve => {
            window.addEventListener('message', function onMsg(e) {
                if (e.data && e.data.type === 'result' && e.data.id === id) {
                    window.removeEventListener('message', onMsg);
                    resolve(e.data);
                }
            });
        });
        window.parent.postMessage({ __howlBridge: 1, type: 'call', id, method: 'EvalHack', args: [] }, location.origin);
        return wait;
    });

    expect(reply.ok).toBe(false);
    expect(reply.error.message).toContain('method not allowed');
    expect(await page.evaluate(() => window.howlBridgeHost.__stats.unknownMethod)).toBeGreaterThan(0);
    // The shim has no generic forwarding surface.
    expect(await panel.evaluate(() => typeof window.go.main.App.EvalHack)).toBe('undefined');
});

test('honest Applied: mismatched effective configuration renders Indeterminate', async ({ page }) => {
    const panel = await openFixture(page);

    // Non-closeable readiness keeps the panel mounted so the label can be read.
    await page.evaluate(() => {
        window.__be.applied = null;
        window.__be.publishApplied = false;   // uncomparable effective state
        window.__be.status = { asr_readiness: 'not-ready', asr_ready: false };
    });
    await panel.evaluate(() => window.AudioSettings.saveAudioSettings());

    const label = panel.locator('#audio-apply-state');
    await expect(label).toContainText('Indeterminate');
    await expect(label).not.toHaveText('Applied');
});

test('Applied appears only when the active generation matches the request', async ({ page }) => {
    const panel = await openFixture(page);

    await page.evaluate(() => { window.__be.status = { asr_readiness: 'not-ready', asr_ready: false }; });
    await panel.evaluate(() => window.AudioSettings.saveAudioSettings());
    await expect(panel.locator('#audio-apply-state')).toHaveText('Applied');
    await expect(panel.locator('#audio-save-state')).toHaveText('Saved');
});

test('apply_in_progress keeps activation gated', async ({ page }) => {
    const panel = await openFixture(page);

    await page.evaluate(() => { window.__be.inProgress = true; });
    await panel.evaluate(() => window.AudioSettings.retryAudioActivation());

    await expect(panel.locator('#audio-apply-state')).toContainText('may still be running');
    expect(await page.evaluate(() => window.__be.calls.length)).toBe(0);
});

test('detach retires the generation and removes the listener', async ({ page }) => {
    await openFixture(page);

    const before = await page.evaluate(() => window.howlBridgeHost.generation);
    expect(await page.evaluate(() => window.howlBridgeHost.listening)).toBe(true);

    await page.evaluate(() => window.howlBridgeHost.detach(document.getElementById('f')));

    const after = await page.evaluate(() => window.howlBridgeHost.generation);
    expect(after).toBeGreaterThan(before);
    expect(await page.evaluate(() => window.howlBridgeHost.listening)).toBe(false);
});
