/**
 * Stage 4 F-01 regression: a blank character engine must inherit the active
 * backend engine without injecting the Piper-specific default voice.
 */

const { test, expect } = require('@playwright/test');
const fs = require('fs');
const path = require('path');

const AUDIO_SERVICE_SOURCE = fs.readFileSync(
    path.resolve(__dirname, '../../../frontend/dist/js/audio-service.js'),
    'utf8'
);

async function preparePage(page, audio, character) {
    await page.setContent('<html><body></body></html>');
    await page.evaluate(({ audio, character }) => {
        window.settingsState = { current: { audio } };
        window.sidebarState = { currentChatData: { characters: [character] } };
        window.__voiceRequest = null;
        window.go = {
            main: {
                App: {
                    SynthesizeSpeech: async () => '',
                    SynthesizeSpeechWithVoice: async (text, voiceJSON) => {
                        window.__voiceRequest = { text, voice: JSON.parse(voiceJSON) };
                        return '';
                    },
                },
            },
        };
    }, { audio, character });
    await page.addScriptTag({ content: AUDIO_SERVICE_SOURCE });
}

test('blank engine and voice inherit a non-Piper active engine without a Piper default', async ({ page }) => {
    await preparePage(page, {
        ttsEnabled: true,
        ttsEngine: 'kokoro',
        ttsVoice: '',
        ttsSpeed: 1,
        ttsVolume: 1,
        ttsNoiseScale: 0.667,
        ttsNoiseW: 0.8,
        ttsSentenceSilence: 0.2,
    }, {
        id: 'character-inherit',
        howl: { voice: { enabled: true, engine: '', voice: '' } },
    });

    await page.evaluate(() => window.AudioService.playTTS('hello', 'character-inherit'));

    const request = await page.evaluate(() => window.__voiceRequest);
    expect(request.voice.character_id).toBe('character-inherit');
    expect(request.voice.engine).toBe('');
    expect(request.voice.voice).toBe('');
});

test('blank engine preserves a configured global voice inheritance', async ({ page }) => {
    await preparePage(page, {
        ttsEnabled: true,
        ttsEngine: 'piper',
        ttsVoice: 'global-voice',
        ttsSpeed: 1,
        ttsVolume: 1,
        ttsNoiseScale: 0.667,
        ttsNoiseW: 0.8,
        ttsSentenceSilence: 0.2,
    }, {
        id: 'character-global',
        howl: { voice: { enabled: true, engine: '', voice: '' } },
    });

    await page.evaluate(() => window.AudioService.playTTS('hello', 'character-global'));

    const request = await page.evaluate(() => window.__voiceRequest);
    expect(request.voice.engine).toBe('');
    expect(request.voice.voice).toBe('global-voice');
});
