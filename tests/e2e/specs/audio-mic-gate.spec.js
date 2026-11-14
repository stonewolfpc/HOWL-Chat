/**
 * Microphone gate at the renderer boundary.
 *
 * Loads the real settings-state and audio-service scripts with a fake Wails
 * backend, and asserts:
 *   1. a persisted ASR selection is not applied at startup;
 *   2. starting capture with no Windows microphone shows the exact copy and
 *      never calls getUserMedia;
 *   3. a browser-level NotFoundError maps to the same no-device copy.
 */

const { test, expect } = require('@playwright/test');
const path = require('path');

const SETTINGS_STATE = path.resolve(__dirname, '../../../frontend/dist/js/settings-state.js');
const AUDIO_SERVICE = path.resolve(__dirname, '../../../frontend/dist/js/audio-service.js');

const NO_DEVICE_MESSAGE = "Windows shows a mic isn't selected. Please connect a microphone to use this option";

async function installBackend(page, options = {}) {
    await page.addInitScript((opts) => {
        window.__micBackend = {
            applyCalls: [],
            getUserMediaCalls: 0,
            notifications: [],
            settingsLoaded: false,
            mic: opts.mic || {
                state: 'no-device',
                message: "Windows shows a mic isn't selected. Please connect a microphone to use this option",
            },
            asrReady: !!opts.asrReady,
            asrReadiness: opts.asrReadiness || 'mic-unavailable',
            getUserMediaError: opts.getUserMediaError || '',
        };

        window.addEventListener('howl:settings-loaded', () => {
            window.__micBackend.settingsLoaded = true;
        });
        window.showNotification = (message, type) => {
            window.__micBackend.notifications.push({ message, type });
        };
        window.AudioWav = { encodeWav16: () => new Uint8Array([1, 2, 3]) };

        const mediaDevices = {
            getUserMedia: async () => {
                window.__micBackend.getUserMediaCalls += 1;
                if (window.__micBackend.getUserMediaError) {
                    throw new DOMException('test failure', window.__micBackend.getUserMediaError);
                }
                return { getTracks: () => [] };
            },
        };
        Object.defineProperty(navigator, 'mediaDevices', { value: mediaDevices, configurable: true });

        const audioSettings = {
            asrEnabled: true,
            asrEngine: 'whisper-base',
            ttsEnabled: false,
            ttsEngine: 'piper',
        };

        window.go = {
            main: {
                App: {
                    GetSettings: async () => ({ audio: audioSettings }),
                    SetSettings: async () => {},
                    ApplyAudioSettings: async (...args) => {
                        window.__micBackend.applyCalls.push(args);
                    },
                    GetAudioStatus: async () => ({
                        asr_ready: window.__micBackend.asrReady,
                        asr_readiness: window.__micBackend.asrReadiness,
                        asr_error: '',
                        tts_ready: false,
                        tts_error: '',
                        capabilities: [],
                        mic: window.__micBackend.mic,
                        apply_in_progress: false,
                        applied: null,
                    }),
                    GetTTSModels: async () => [],
                },
            },
        };
    }, options);
}

async function loadScripts(page) {
    await page.goto('about:blank');
    await page.addScriptTag({ path: SETTINGS_STATE });
    await page.addScriptTag({ path: AUDIO_SERVICE });
    await page.waitForFunction(() => window.__micBackend.settingsLoaded === true, null, { timeout: 5000 });
    await page.waitForTimeout(200);
}

test('persisted ASR selection is restored but not applied at startup', async ({ page }) => {
    await installBackend(page);
    await loadScripts(page);

    const state = await page.evaluate(() => ({
        enabled: window.AudioService.isASREnabled(),
        applyCalls: window.__micBackend.applyCalls.length,
    }));
    expect(state.enabled).toBe(true);
    expect(state.applyCalls).toBe(0);
});

test('capture is blocked with the exact no-device copy and never opens a stream', async ({ page }) => {
    await installBackend(page);
    await loadScripts(page);

    await page.evaluate(() => window.AudioService.startRecording());
    await page.waitForFunction(() => window.__micBackend.notifications.length > 0, null, { timeout: 5000 });

    const state = await page.evaluate(() => ({
        notifications: window.__micBackend.notifications,
        getUserMediaCalls: window.__micBackend.getUserMediaCalls,
    }));
    expect(state.notifications[0].message).toBe(NO_DEVICE_MESSAGE);
    expect(state.getUserMediaCalls).toBe(0);
});

test('browser NotFoundError maps to the no-device copy', async ({ page }) => {
    await installBackend(page, {
        mic: { state: 'ready', message: '' },
        asrReady: true,
        asrReadiness: 'local-initialized',
        getUserMediaError: 'NotFoundError',
    });
    await loadScripts(page);

    // Establish the applied state the way the settings iframe does after a real
    // activation, so this test isolates the getUserMedia error mapping.
    await page.evaluate(() => window.dispatchEvent(new CustomEvent('howl:audio-settings-applied', {
        detail: {
            audio: { asrEnabled: true, ttsEnabled: false, asrEngine: 'whisper-base', ttsEngine: 'piper' },
            readiness: 'local-initialized',
            detail: '',
        },
    })));

    await page.evaluate(() => window.AudioService.startRecording());
    await page.waitForFunction(() => window.__micBackend.notifications.length > 0, null, { timeout: 5000 });

    const notification = await page.evaluate(() => window.__micBackend.notifications[0]);
    expect(notification.message).toBe(NO_DEVICE_MESSAGE);
});
