/**
 * D1 activation-lock coverage.
 *
 * The fixture loads the real iframe-bridge-host.js and the real
 * menu-functions.js, and opens the modal only through openAudioSettings(), so
 * the overlay listener, the parent Escape handler, attach/detach and the
 * howl:audio-settings-closed path are the production implementations.
 *
 * Unrelated menu-action symbols are stubbed only so menu-functions.js can finish
 * its export object; openPanel asserts the audio entry points are the real ones.
 *
 * These are harness tests. They are not packaged-app or restart proof.
 */

const { test, expect } = require('@playwright/test');
const http = require('http');
const fs = require('fs');
const path = require('path');

const DIST = path.resolve(__dirname, '../../../frontend/dist');
const MARKER_KEY = 'howl.audio.activation';

const MENU_STUBS = [
    'editUndo', 'editRedo', 'editCut', 'editCopy', 'editPaste', 'editSelectAll',
    'editFind', 'editReplace', 'toggleColorblindMode', 'toggleHighContrastMode',
    'setTheme', 'increaseUIScale', 'decreaseUIScale', 'resetUIScale',
    'toggleSidebar', 'toggleCharacterPanel', 'toggleWorldPanel', 'toggleFullscreen'
];

const FIXTURE = `<!doctype html>
<html><body>
<script>
const q = new URLSearchParams(location.search);
window.__be = {
  calls: [],
  applied: null,
  inProgress: false,
  fail: false,
  delay: 0,
  publishApplied: true,
  statusReads: 0,
  setSettingsCalls: 0,
  statusFail: q.get('statusFail') === '1',
  desiredGen: 1,
  appliedGen: 1
};
window.go = { main: { App: {
  GetSettings: async () => ({ audio: { asrEnabled: true, asrEngine: 'whisper-base', ttsEngine: 'piper' } }),
  SetSettings: async () => { window.__be.setSettingsCalls++; },
  GetAudioStatus: async () => {
    window.__be.statusReads++;
    if (window.__be.statusFail) throw new Error('status unavailable');
    return {
      asr_readiness: 'local-initialized',
      asr_ready: true,
      tts_ready: false,
      apply_in_progress: window.__be.inProgress,
      applied: window.__be.applied,
      desired_gen: window.__be.desiredGen,
      applied_gen: window.__be.appliedGen
    };
  },
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
  ScanTTSDir: async () => {}
} } };
window.__stubNames = [];
window.__listenerStats = { keydown: 0 };
(function () {
  const names = ${JSON.stringify(MENU_STUBS)};
  names.forEach(function (name) {
    window.__stubNames.push(name);
    window[name] = function () {};
  });
  const add = document.addEventListener.bind(document);
  const remove = document.removeEventListener.bind(document);
  document.addEventListener = function (type, fn, opts) {
    if (type === 'keydown') window.__listenerStats.keydown++;
    return add(type, fn, opts);
  };
  document.removeEventListener = function (type, fn, opts) {
    if (type === 'keydown' && window.__listenerStats.keydown > 0) window.__listenerStats.keydown--;
    return remove(type, fn, opts);
  };
})();
</script>
<script src="/js/iframe-bridge-host.js"></script>
<script src="/js/menu-functions.js"></script>
</body></html>`;

let server;
let base;

test.beforeAll(async () => {
    server = http.createServer((req, res) => {
        const url = req.url.split('?')[0];
        if (url === '/__lock-fixture.html') {
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

test.afterAll(async () => {
    if (!server) return;
    if (typeof server.closeAllConnections === 'function') server.closeAllConnections();
    await new Promise(r => server.close(r));
});

function trackPageErrors(page) {
    if (!page.__pageErrors) {
        page.__pageErrors = [];
        page.on('pageerror', e => page.__pageErrors.push(String(e && e.message ? e.message : e)));
    }
    return page.__pageErrors;
}

async function seedMarker(page) {
    await page.addInitScript(key => {
        try {
            window.sessionStorage.setItem(key, JSON.stringify({
                activationId: '1:1',
                startedAt: Date.now()
            }));
        } catch (e) { /* storage unavailable */ }
    }, MARKER_KEY);
}

async function waitForPanel(page, { expectSaveEnabled = true } = {}) {
    await expect.poll(
        async () => page.frames().map(f => f.url()).join(' '),
        { timeout: 15000 }
    ).toContain('audio-settings.html');
    const panel = page.frames().find(f => f.url().includes('audio-settings.html'));
    expect(panel).toBeTruthy();
    await expect.poll(async () => panel.evaluate(() => window.__howlBridgeIframe?.mode)).toBe('ready');
    if (expectSaveEnabled) {
        await expect(panel.locator('#audio-settings-save')).toBeEnabled({ timeout: 10000 });
    }
    return panel;
}

async function openPanel(page, { expectSaveEnabled = true, query = '' } = {}) {
    trackPageErrors(page);
    await page.goto(base + '/__lock-fixture.html' + (query ? '?' + query : ''), { waitUntil: 'load' });

    // Stub isolation: the exported audio entry points must be the real ones.
    const realness = await page.evaluate(() => ({
        sameOpen: window.menuFunctions.openAudioSettings === window.openAudioSettings,
        sameTts: window.menuFunctions.openTTSSettings === window.openTTSSettings,
        body: String(window.openAudioSettings),
        stubHit: window.__stubNames.some(n => /audio/i.test(n))
    }));
    expect(realness.sameOpen).toBe(true);
    expect(realness.sameTts).toBe(true);
    expect(realness.body).toContain('_openAudioSettingsIframe');
    expect(realness.stubHit).toBe(false);

    await page.evaluate(() => { window.__keydownBaseline = window.__listenerStats.keydown; });
    await page.evaluate(() => window.openAudioSettings());

    const panel = await waitForPanel(page, { expectSaveEnabled });
    expect(page.__pageErrors).toEqual([]);
    return panel;
}

async function openLockedPanel(page) {
    await seedMarker(page);
    return openPanel(page, { expectSaveEnabled: false });
}

async function expectModalOpen(page) {
    await expect(page.locator('#audio-settings-modal')).toHaveCount(1);
    await expect(page.locator('#audio-settings-modal iframe')).toHaveCount(1);
    expect(await page.evaluate(() => window.howlBridgeHost.listening)).toBe(true);
    expect(await page.evaluate(() => window.__listenerStats.keydown))
        .toBe(await page.evaluate(() => window.__keydownBaseline) + 1);
}

async function expectClosedWithoutSideEffects(page, before) {
    await expect(page.locator('#audio-settings-modal')).toHaveCount(0);
    expect(await page.evaluate(() => window.howlBridgeHost.listening)).toBe(false);
    await expect.poll(
        async () => page.frames().filter(f => f.url().includes('audio-settings.html')).length
    ).toBe(0);
    expect(await page.evaluate(() => window.__be.setSettingsCalls)).toBe(0);
    expect(await page.evaluate(() => window.__be.calls.length)).toBe(0);
    expect(await page.evaluate(() => window.howlBridgeHost.isActivationLocked())).toBe(true);
    expect(await page.evaluate(() => window.howlBridgeHost.activation())).toEqual(before);
    expect(await page.evaluate(key => window.sessionStorage.getItem(key), MARKER_KEY)).not.toBeNull();
    expect(await page.evaluate(() => window.__listenerStats.keydown))
        .toBe(await page.evaluate(() => window.__keydownBaseline));
    expect(page.__pageErrors).toEqual([]);
}

const activate = (panel) => panel.evaluate(() =>
    window.go.main.App.ApplyAudioSettings(true, false, 'whisper-base', 'piper'));

test('second activation while locked is refused without reaching the backend', async ({ page }) => {
    const panel = await openPanel(page);
    await page.evaluate(() => { window.__be.delay = 1200; });

    await panel.evaluate(() => {
        window.__first = window.go.main.App.ApplyAudioSettings(true, false, 'whisper-base', 'piper');
    });
    await expect.poll(() => page.evaluate(() => window.__be.calls.length)).toBe(1);
    await expect.poll(() => page.evaluate(() => window.howlBridgeHost.isActivationLocked())).toBe(true);

    const second = await panel.evaluate(() => window.go.main.App
        .ApplyAudioSettings(true, false, 'whisper-base', 'piper')
        .then(() => ({ ok: true }), e => ({ ok: false, name: e.name, message: e.message })));

    expect(second.ok).toBe(false);
    expect(second.name).toBe('ActivationLockedError');
    expect(second.message).toContain('did not reach the backend');
    expect(await page.evaluate(() => window.__be.calls.length)).toBe(1);
    expect(await page.evaluate(() => window.howlBridgeHost.__stats.refusedActivations)).toBe(1);

    await panel.evaluate(() => window.__first);
});

test('activation is accepted again after the live lock settles', async ({ page }) => {
    const panel = await openPanel(page);

    await activate(panel);
    await expect.poll(() => page.evaluate(() => window.howlBridgeHost.isActivationLocked())).toBe(false);

    await activate(panel);
    expect(await page.evaluate(() => window.__be.calls.length)).toBe(2);
});

test('a failed status read is unknown, not idle', async ({ page }) => {
    const panel = await openPanel(page);

    await page.evaluate(() => { window.__be.statusFail = true; });
    await panel.evaluate(() => window.dispatchEvent(new CustomEvent('howl:activation-settled')));
    await page.waitForTimeout(150);

    await panel.evaluate(() => window.AudioSettings.retryAudioActivation());

    expect(await page.evaluate(() => window.__be.calls.length)).toBe(0);
    await expect(panel.locator('#audio-apply-state')).toContainText('may still be running');
});

test('an authoritative idle read permits activation', async ({ page }) => {
    const panel = await openPanel(page);

    await page.evaluate(() => { window.__be.statusFail = false; window.__be.inProgress = false; });
    await panel.evaluate(() => window.dispatchEvent(new CustomEvent('howl:activation-settled')));
    await page.waitForTimeout(150);

    // Fire-and-forget: a successful activation closes the panel, detaching the
    // frame before evaluate could deliver a return value.
    await panel.evaluate(() => { window.AudioSettings.retryAudioActivation(); });

    await expect.poll(() => page.evaluate(() => window.__be.calls.length)).toBe(1);
});

test('restored marker stays held when builds are zero and generations are equal', async ({ page }) => {
    const panel = await openLockedPanel(page);
    await expectModalOpen(page);

    expect(await page.evaluate(() => window.howlBridgeHost.isActivationLocked())).toBe(true);
    expect((await page.evaluate(() => window.howlBridgeHost.activation())).restored).toBe(true);

    await expect(panel.locator('#audio-apply-state')).toHaveText('State unknown — restart required');
    await expect(panel.locator('#audio-outcome-detail')).toContainText('Close HOWL Chat completely');
    await expect(panel.locator('#audio-settings-save')).toBeDisabled();
    await expect(panel.locator('#audio-retry-activation')).toBeHidden();

    await page.waitForTimeout(3200);

    expect(await page.evaluate(() => window.__be.calls.length)).toBe(0);
    expect(await page.evaluate(() => window.howlBridgeHost.isActivationLocked())).toBe(true);
});

test('navigation during an activation keeps the lock after the backend settles', async ({ page }) => {
    const panel = await openPanel(page);
    await page.evaluate(() => { window.__be.delay = 5000; });
    await panel.evaluate(() => {
        window.go.main.App.ApplyAudioSettings(true, false, 'whisper-base', 'piper');
    });
    await expect.poll(() => page.evaluate(() => window.__be.calls.length)).toBe(1);
    await expect.poll(() => page.evaluate(() => window.howlBridgeHost.isActivationLocked())).toBe(true);

    await page.goto(base + '/__lock-fixture.html?phase=2', { waitUntil: 'load' });
    await page.evaluate(() => { window.__keydownBaseline = window.__listenerStats.keydown; });
    await page.evaluate(() => window.openAudioSettings());
    const next = await waitForPanel(page, { expectSaveEnabled: false });

    expect(await page.evaluate(() => window.howlBridgeHost.isActivationLocked())).toBe(true);
    expect((await page.evaluate(() => window.howlBridgeHost.activation())).restored).toBe(true);
    await expect(next.locator('#audio-apply-state')).toHaveText('State unknown — restart required');

    await page.waitForTimeout(3000);

    expect(await page.evaluate(() => window.__be.calls.length)).toBe(0);
    expect(await page.evaluate(() => window.howlBridgeHost.isActivationLocked())).toBe(true);
});

test('restored marker stays held when status reads fail', async ({ page }) => {
    await seedMarker(page);
    const panel = await openPanel(page, { expectSaveEnabled: false, query: 'statusFail=1' });
    await expectModalOpen(page);

    expect(await page.evaluate(() => window.howlBridgeHost.isActivationLocked())).toBe(true);
    await expect(panel.locator('#audio-apply-state')).toHaveText('State unknown — restart required');
    await page.waitForTimeout(1500);

    expect(await page.evaluate(() => window.__be.statusReads)).toBeGreaterThan(0);
    expect(await page.evaluate(() => window.__be.calls.length)).toBe(0);
});

test('activation-settled with a restored lock triggers a read and never unlocks', async ({ page }) => {
    const panel = await openLockedPanel(page);
    await expectModalOpen(page);

    const before = await page.evaluate(() => window.__be.statusReads);
    await panel.evaluate(() => window.dispatchEvent(new CustomEvent('howl:activation-settled')));
    await expect.poll(() => page.evaluate(() => window.__be.statusReads)).toBeGreaterThan(before);

    expect(await page.evaluate(() => window.howlBridgeHost.isActivationLocked())).toBe(true);
    await expect(panel.locator('#audio-apply-state')).toHaveText('State unknown — restart required');
});

test('Save and Retry do nothing while a restored lock is unresolved', async ({ page }) => {
    const panel = await openLockedPanel(page);
    await expectModalOpen(page);

    await panel.evaluate(() => window.AudioSettings.retryAudioActivation());
    await panel.evaluate(() => { document.getElementById('audio-settings-save').click(); });
    await page.waitForTimeout(300);

    expect(await page.evaluate(() => window.__be.calls.length)).toBe(0);
    expect(await page.evaluate(() => window.howlBridgeHost.isActivationLocked())).toBe(true);
    await expect(panel.locator('#audio-settings-save')).toBeDisabled();
});

test('a fresh browsing context carries no restored marker (restart proxy)', async ({ browser }) => {
    const ctx = await browser.newContext();
    const page = await ctx.newPage();
    const panel = await openPanel(page, { expectSaveEnabled: true });

    expect(await page.evaluate(() => window.howlBridgeHost.isActivationLocked())).toBe(false);

    await panel.evaluate(() => { window.AudioSettings.retryAudioActivation(); });
    await expect.poll(() => page.evaluate(() => window.__be.calls.length)).toBe(1);

    await ctx.close();
});

// ---------------------------------------------------------------------------
// Dismissal correction: the real parent integration must close the panel with no
// settings write, no activation and no lock release.
// ---------------------------------------------------------------------------

test('Escape dismisses a restored-lock panel without writing or releasing', async ({ page }) => {
    const panel = await openLockedPanel(page);
    await expectModalOpen(page);
    const before = await page.evaluate(() => window.howlBridgeHost.activation());

    await panel.locator('body').press('Escape');

    await expectClosedWithoutSideEffects(page, before);
});

test('the panel close button dismisses without writing or releasing', async ({ page }) => {
    const panel = await openLockedPanel(page);
    await expectModalOpen(page);
    const before = await page.evaluate(() => window.howlBridgeHost.activation());

    await panel.locator('button.real-settings-close').click();

    await expectClosedWithoutSideEffects(page, before);
});

test('an overlay click dismisses without writing or releasing', async ({ page }) => {
    await openLockedPanel(page);
    await expectModalOpen(page);
    const before = await page.evaluate(() => window.howlBridgeHost.activation());

    const point = await page.evaluate(() => {
        const modal = document.getElementById('audio-settings-modal');
        const container = modal.querySelector('.real-settings-container');
        const rect = container.getBoundingClientRect();
        const outside = (x, y) => x < rect.left || x > rect.right || y < rect.top || y > rect.bottom;
        for (let y = 8; y < innerHeight - 8; y += 24) {
            for (let x = innerWidth - 8; x > 8; x -= 24) {
                const el = document.elementFromPoint(x, y);
                if (el && el.id === 'audio-settings-modal' && outside(x, y)) return { x, y };
            }
        }
        return null;
    });
    expect(point).not.toBeNull();

    await page.mouse.click(point.x, point.y);

    await expectClosedWithoutSideEffects(page, before);
});

test('parent-window Escape dismisses without writing or releasing', async ({ page }) => {
    await openLockedPanel(page);
    await expectModalOpen(page);
    const before = await page.evaluate(() => window.howlBridgeHost.activation());

    await page.evaluate(() => {
        document.body.setAttribute('tabindex', '-1');
        document.body.focus();
    });
    await page.keyboard.press('Escape');

    await expectClosedWithoutSideEffects(page, before);
});

test('reopening after dismissal still reports the restored lock', async ({ page }) => {
    const panel = await openLockedPanel(page);
    await expectModalOpen(page);
    const before = await page.evaluate(() => window.howlBridgeHost.activation());
    const generationBefore = await page.evaluate(() => window.howlBridgeHost.generation);

    await panel.locator('button.real-settings-close').click();
    await expectClosedWithoutSideEffects(page, before);

    await page.evaluate(() => window.openAudioSettings());
    const next = await waitForPanel(page, { expectSaveEnabled: false });
    await expect(page.locator('#audio-settings-modal')).toHaveCount(1);

    await expect(next.locator('#audio-apply-state')).toHaveText('State unknown — restart required');
    expect(await page.evaluate(() => window.howlBridgeHost.generation)).toBeGreaterThan(generationBefore);
    expect(await page.evaluate(() => window.howlBridgeHost.isActivationLocked())).toBe(true);
    expect(await page.evaluate(() => window.__be.calls.length)).toBe(0);
    expect(await page.evaluate(() => window.__be.setSettingsCalls)).toBe(0);
    expect(await page.evaluate(key => window.sessionStorage.getItem(key), MARKER_KEY)).not.toBeNull();
    expect(page.__pageErrors).toEqual([]);
});

test('a programmatic save is refused before writing', async ({ page }) => {
    const panel = await openLockedPanel(page);
    await expectModalOpen(page);
    const before = await page.evaluate(() => window.howlBridgeHost.activation());

    const result = await panel.evaluate(() => window.AudioSettings.saveAudioSettings());

    expect(result).toBe(false);
    expect(await page.evaluate(() => window.__be.setSettingsCalls)).toBe(0);
    expect(await page.evaluate(() => window.__be.calls.length)).toBe(0);
    expect(await page.evaluate(() => window.howlBridgeHost.isActivationLocked())).toBe(true);
    expect(await page.evaluate(() => window.howlBridgeHost.activation())).toEqual(before);
    expect(page.__pageErrors).toEqual([]);
});
