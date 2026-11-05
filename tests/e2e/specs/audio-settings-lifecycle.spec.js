/**
 * Unit B lifecycle coverage for the audio-settings iframe.
 *
 * These tests load the real page and scripts, with a fake Wails bridge. They
 * assert the ordering the panel promises: hydration before Save, awaited
 * persistence before Apply, and independent Saved / Applied / Ready outcomes.
 */

const { test, expect } = require('@playwright/test');
const fs = require('fs');
const path = require('path');
const { pathToFileURL } = require('url');

const PAGE = pathToFileURL(path.resolve(__dirname, '../../../frontend/dist/audio-settings.html')).href;
const MENU_SCRIPT = path.resolve(__dirname, '../../../frontend/dist/js/menu-functions.js');
const MENU_SOURCE = fs.readFileSync(MENU_SCRIPT, 'utf8');
const SETTINGS_SOURCE = fs.readFileSync(path.resolve(__dirname, '../../../frontend/dist/js/settings-state.js'), 'utf8');
const AUDIO_SERVICE_SOURCE = fs.readFileSync(path.resolve(__dirname, '../../../frontend/dist/js/audio-service.js'), 'utf8');
const AUDIO_WAV_SOURCE = fs.readFileSync(path.resolve(__dirname, '../../../frontend/dist/js/audio-wav.js'), 'utf8');

const LOCAL_READY = {
    asr_ready: true,
    asr_readiness: 'local-initialized',
    asr_error: '',
    tts_ready: false,
    tts_error: '',
    capabilities: [],
};

async function installBackend(page, options = {}) {
    await page.addInitScript((opts) => {
        window.__audioBackend = {
            settings: opts.settings || {},
            settingsDelay: opts.settingsDelay || 0,
            saveError: opts.saveError || '',
            applyError: opts.applyError || '',
            status: opts.status || {
                asr_ready: false,
                asr_readiness: 'disabled',
                asr_error: 'ASR is disabled',
                tts_ready: false,
                tts_error: '',
                capabilities: [],
            },
            setCalls: [],
            applyCalls: [],
            statusCalls: 0,
            notifications: [],
        };
        window.showNotification = (message, type) => {
            window.__audioBackend.notifications.push({ message, type });
        };
        window.go = {
            main: {
                App: {
                    GetSettings: async () => {
                        if (window.__audioBackend.settingsDelay > 0) {
                            await new Promise(resolve => setTimeout(resolve, window.__audioBackend.settingsDelay));
                        }
                        return window.__audioBackend.settings;
                    },
                    SetSettings: async (settings) => {
                        if (window.__audioBackend.saveError) {
                            throw new Error(window.__audioBackend.saveError);
                        }
                        window.__audioBackend.setCalls.push(JSON.parse(JSON.stringify(settings)));
                    },
                    ApplyAudioSettings: async (...args) => {
                        window.__audioBackend.applyCalls.push(args);
                        if (window.__audioBackend.applyError) {
                            throw new Error(window.__audioBackend.applyError);
                        }
                    },
                    GetAudioStatus: async () => {
                        window.__audioBackend.statusCalls += 1;
                        const last = window.__audioBackend.applyCalls[window.__audioBackend.applyCalls.length - 1];
                        const applied = last ? {
                            asr_enabled: !!last[0],
                            tts_enabled: !!last[1],
                            asr_engine: last[2] || '',
                            tts_engine: last[3] || ''
                        } : null;
                        return { ...window.__audioBackend.status, applied, apply_in_progress: false };
                    },
                    GetTTSModels: async () => [],
                },
            },
        };
    }, options);
}

async function chooseAudio(page, engine) {
    await page.evaluate((value) => {
        const toggle = document.getElementById('asr-enabled');
        toggle.checked = true;
        toggle.dispatchEvent(new Event('change', { bubbles: true }));
        const select = document.getElementById('asr-engine');
        select.value = value;
        select.dispatchEvent(new Event('change', { bubbles: true }));
    }, engine);
}

test('Save waits for hydration, persists first, applies, and closes only when locally ready', async ({ page }) => {
    await installBackend(page, { settingsDelay: 1000, status: LOCAL_READY });
    await page.goto(PAGE, { waitUntil: 'domcontentloaded' });

    const save = page.locator('#audio-settings-save');
    await expect(save).toBeDisabled();
    await expect(save).toBeEnabled({ timeout: 5000 });

    await chooseAudio(page, 'whisper-base');
    await save.click();

    await expect(page.locator('#audio-settings-overlay')).not.toHaveClass(/show/);
    await expect(page.locator('#audio-save-state')).toHaveText('Saved');
    await expect(page.locator('#audio-apply-state')).toHaveText('Applied');
    await expect(page.locator('#audio-ready-state')).toHaveText('Ready (local initialized)');

    const calls = await page.evaluate(() => ({ ...window.__audioBackend, status: undefined }));
    expect(calls.setCalls).toHaveLength(1);
    expect(calls.setCalls[0].audio.asrEnabled).toBe(true);
    expect(calls.setCalls[0].audio.asrEngine).toBe('whisper-base');
    expect(calls.applyCalls).toEqual([[true, false, 'whisper-base', 'piper']]);
    expect(calls.notifications.some(n => n.type === 'error')).toBe(false);
});

test('a remote engine is never rendered as locally ready, even if the status flag is wrong', async ({ page }) => {
    await installBackend(page, { status: LOCAL_READY });
    await page.goto(PAGE);
    await expect(page.locator('#audio-settings-save')).toBeEnabled();

    await chooseAudio(page, 'qwen3-asr-1.7b');
    await page.locator('#audio-settings-save').click();

    await expect(page.locator('#audio-save-state')).toHaveText('Saved');
    await expect(page.locator('#audio-apply-state')).toHaveText('Applied');
    await expect(page.locator('#audio-ready-state')).toContainText('Remote engine');
    await expect(page.locator('#audio-ready-state')).not.toHaveText('Ready (local initialized)');
    await expect(page.locator('#audio-settings-overlay')).toHaveClass(/show/);

    const notifications = await page.evaluate(() => window.__audioBackend.notifications);
    expect(notifications.some(n => n.type === 'warning' && /not verified/i.test(n.message))).toBe(true);
});

test('a settings-file failure is reported as Not saved and Apply is not attempted', async ({ page }) => {
    await installBackend(page, { saveError: 'disk is read-only' });
    await page.goto(PAGE);
    await expect(page.locator('#audio-settings-save')).toBeEnabled();

    await chooseAudio(page, 'whisper-base');
    await page.locator('#audio-settings-save').click();

    await expect(page.locator('#audio-save-state')).toHaveText('Not saved');
    await expect(page.locator('#audio-apply-state')).toHaveText('Not applied');
    await expect(page.locator('#audio-settings-overlay')).toHaveClass(/show/);
    const applyCalls = await page.evaluate(() => window.__audioBackend.applyCalls.length);
    expect(applyCalls).toBe(0);
});

test('activation failure keeps Saved visible and Retry applies without saving or reselecting', async ({ page }) => {
    await installBackend(page, { applyError: 'model failed to initialize', status: LOCAL_READY });
    await page.goto(PAGE);
    await expect(page.locator('#audio-settings-save')).toBeEnabled();

    await chooseAudio(page, 'whisper-base');
    await page.locator('#audio-settings-save').click();

    await expect(page.locator('#audio-save-state')).toHaveText('Saved');
    await expect(page.locator('#audio-apply-state')).toHaveText('Apply failed');
    await expect(page.locator('#audio-retry-activation')).toBeVisible();
    const firstSetCalls = await page.evaluate(() => window.__audioBackend.setCalls.length);
    expect(firstSetCalls).toBe(1);

    await page.evaluate(() => { window.__audioBackend.applyError = ''; });
    await page.locator('#audio-retry-activation').click();

    await expect(page.locator('#audio-apply-state')).toHaveText('Applied');
    await expect(page.locator('#audio-ready-state')).toHaveText('Ready (local initialized)');
    await expect(page.locator('#audio-settings-overlay')).not.toHaveClass(/show/);
    const calls = await page.evaluate(() => ({
        setCalls: window.__audioBackend.setCalls.length,
        applyCalls: window.__audioBackend.applyCalls.length,
    }));
    expect(calls.setCalls).toBe(1);
    expect(calls.applyCalls).toBe(2);
});

test('repeated parent open/close does not accumulate keydown handlers', async ({ page }) => {
    await page.setContent('<html><body></body></html>');
    await page.addScriptTag({
        content: `
            window.__keydown = { active: 0, added: 0, removed: 0 };
            const add = document.addEventListener.bind(document);
            const remove = document.removeEventListener.bind(document);
            document.addEventListener = (type, fn, opts) => {
                if (type === 'keydown') {
                    window.__keydown.active++;
                    window.__keydown.added++;
                }
                return add(type, fn, opts);
            };
            document.removeEventListener = (type, fn, opts) => {
                if (type === 'keydown') {
                    window.__keydown.active--;
                    window.__keydown.removed++;
                }
                return remove(type, fn, opts);
            };
            [
                'editUndo', 'editRedo', 'editCut', 'editCopy', 'editPaste', 'editSelectAll',
                'editFind', 'editReplace', 'toggleColorblindMode', 'toggleHighContrastMode',
                'setTheme', 'increaseUIScale', 'decreaseUIScale', 'resetUIScale',
                'toggleSidebar', 'toggleCharacterPanel', 'toggleWorldPanel', 'toggleFullscreen',
            ].forEach(name => { window[name] = function () {}; });
        `,
    });
    await page.addScriptTag({ content: MENU_SOURCE });
    const baseline = await page.evaluate(() => window.__keydown.active);

    for (let i = 0; i < 3; i++) {
        await page.evaluate(() => window.menuFunctions.openAudioSettings());
    }
    const afterOpens = await page.evaluate(() => window.__keydown.active);
    expect(afterOpens).toBe(baseline + 1);

    await page.evaluate(() => {
        document.dispatchEvent(new CustomEvent('howl:audio-settings-closed'));
    });
    const afterClose = await page.evaluate(() => window.__keydown.active);
    expect(afterClose).toBe(baseline);
    await expect(page.locator('#audio-settings-modal')).toHaveCount(0);
});

test('parent AudioService receives the saved state and gates the microphone on readiness', async ({ page }) => {
    await page.setContent('<html><body></body></html>');
    await page.addScriptTag({
        content: `
            window.__parentAudio = { micCalls: 0, notifications: [] };
            var showNotification = function (message, type) {
                window.__parentAudio.notifications.push({ message, type });
            };
            window.showNotification = showNotification;
            Object.defineProperty(navigator, 'mediaDevices', {
                configurable: true,
                value: {
                    getUserMedia: async () => {
                        window.__parentAudio.micCalls += 1;
                        throw new Error('test capture stopped');
                    },
                },
            });
            window.go = { main: { App: { GetSettings: async () => ({}) } } };
        `,
    });
    await page.addScriptTag({ content: SETTINGS_SOURCE });
    await page.addScriptTag({ content: AUDIO_WAV_SOURCE });
    await page.addScriptTag({ content: AUDIO_SERVICE_SOURCE });

    const initial = await page.evaluate(() => ({
        enabled: window.AudioService.isASREnabled(),
        readiness: window.AudioService.getASRReadiness().state,
    }));
    expect(initial).toEqual({ enabled: false, readiness: 'disabled' });

    await page.evaluate(() => window.dispatchEvent(new CustomEvent('howl:audio-settings-applied', {
        detail: {
            audio: { asrEnabled: true, ttsEnabled: false, asrEngine: 'whisper-base', ttsEngine: 'piper' },
            readiness: 'not-ready',
            detail: 'model failed to initialize',
        },
    })));
    expect(await page.evaluate(() => ({
        enabled: window.AudioService.isASREnabled(),
        readiness: window.AudioService.getASRReadiness(),
    }))).toEqual({ enabled: true, readiness: { state: 'not-ready', detail: 'model failed to initialize' } });

    await page.evaluate(() => window.AudioService.startRecording());
    expect(await page.evaluate(() => window.__parentAudio.micCalls)).toBe(0);

    await page.evaluate(() => window.dispatchEvent(new CustomEvent('howl:audio-settings-applied', {
        detail: {
            audio: { asrEnabled: true, ttsEnabled: false, asrEngine: 'whisper-base', ttsEngine: 'piper' },
            readiness: 'local-initialized',
            detail: 'local engine initialized',
        },
    })));
    await page.evaluate(() => window.AudioService.startRecording());
    expect(await page.evaluate(() => window.__parentAudio.micCalls)).toBe(1);
});

async function installParentAudioPage(page, options = {}) {
    await page.setContent('<html><body></body></html>');
    await page.addScriptTag({
        content: `
            window.__parentAudio = {
                settings: ${JSON.stringify(options.settings || {})},
                status: ${JSON.stringify(options.status || {
                    asr_ready: false,
                    asr_readiness: 'disabled',
                    asr_error: 'ASR is disabled',
                    tts_ready: false,
                    tts_error: '',
                })},
                micCalls: 0,
                applyCalls: [],
                statusCalls: 0,
                notifications: [],
                applyError: '',
            };
            var showNotification = function (message, type) {
                window.__parentAudio.notifications.push({ message, type });
            };
            window.showNotification = showNotification;
            Object.defineProperty(navigator, 'mediaDevices', {
                configurable: true,
                value: {
                    getUserMedia: async () => {
                        window.__parentAudio.micCalls += 1;
                        throw new Error('test capture stopped');
                    },
                },
            });
            window.go = {
                main: {
                    App: {
                        GetSettings: async () => window.__parentAudio.settings,
                        ApplyAudioSettings: async (...args) => {
                            window.__parentAudio.applyCalls.push(args);
                            if (window.__parentAudio.applyError) {
                                throw new Error(window.__parentAudio.applyError);
                            }
                        },
                        GetAudioStatus: async () => {
                            window.__parentAudio.statusCalls += 1;
                            return { ...window.__parentAudio.status };
                        },
                    },
                },
            };
        `,
    });
    await page.addScriptTag({ content: AUDIO_WAV_SOURCE });
    if (options.audioServiceFirst) {
        await page.addScriptTag({ content: AUDIO_SERVICE_SOURCE });
        await page.addScriptTag({ content: SETTINGS_SOURCE });
    } else {
        await page.addScriptTag({ content: SETTINGS_SOURCE });
        await page.addScriptTag({ content: AUDIO_SERVICE_SOURCE });
    }
}

test('remote-unverified permits microphone capture and never becomes local-ready', async ({ page }) => {
    await installParentAudioPage(page);
    await page.evaluate(() => window.dispatchEvent(new CustomEvent('howl:audio-settings-applied', {
        detail: {
            audio: { asrEnabled: true, ttsEnabled: false, asrEngine: 'qwen3-asr-1.7b', ttsEngine: 'piper' },
            readiness: 'remote-unverified',
            detail: 'remote engine availability is not locally verified',
        },
    })));

    expect(await page.evaluate(() => window.AudioService.getASRReadiness().state)).toBe('remote-unverified');
    await page.evaluate(() => window.AudioService.startRecording());
    expect(await page.evaluate(() => window.__parentAudio.micCalls)).toBe(1);
    expect(await page.evaluate(() => window.AudioService.getASRReadiness().state)).toBe('remote-unverified');
});

// Ownership decision: startup reconciliation is renderer-owned by AudioService
// after the real howl:settings-loaded event. Go loads settings, but does not run
// a competing startup apply. This test loads the scripts in index.html order.
test('startup reconciles the exact persisted engine through the real settings-loaded event', async ({ page }) => {
    await installParentAudioPage(page, {
        audioServiceFirst: true,
        settings: {
            audio: {
                asrEnabled: true,
                ttsEnabled: false,
                asrEngine: 'whisper-base',
                ttsEngine: 'piper',
            },
        },
        status: LOCAL_READY,
    });

    await page.waitForFunction(() => window.__parentAudio.applyCalls.length > 0);
    const state = await page.evaluate(() => ({
        applyCalls: window.__parentAudio.applyCalls,
        enabled: window.AudioService.isASREnabled(),
        readiness: window.AudioService.getASRReadiness().state,
    }));

    expect(state.applyCalls).toEqual([[true, false, 'whisper-base', 'piper']]);
    expect(state.enabled).toBe(true);
    expect(state.readiness).toBe('local-initialized');
});