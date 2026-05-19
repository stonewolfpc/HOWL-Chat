/**
 * E2E: Save → Reload Persistence Tests
 *
 * Verifies that chat state (scenario, characters, world) survives
 * a full page reload and that the backend correctly persists and
 * restores it.
 */

const { test, expect } = require('@playwright/test');
const { waitForAppReady, assertDisplayText, buildSeedScript } = require('./helpers');

const TINY_PNG = 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwADhQGAWjR9awAAAABJRU5ErkJggg==';

const SEED_SCENARIO = {
    id:         'scenario-persist-001',
    title:      'Persistence Test Scenario',
    characters: [
        { name: 'Alice', role: 'hero',     avatar: TINY_PNG },
        { name: 'Bob',   role: 'sidekick', avatar: TINY_PNG },
    ],
};

const SEED_CHARS = [
    { id: 'char_persist_alice', name: 'Alice', personality: 'Brave',  avatar: TINY_PNG },
    { id: 'char_persist_bob',   name: 'Bob',   personality: 'Clever', avatar: TINY_PNG },
];

test.describe('Save → reload persistence', () => {

    test('scenario name survives page reload', async ({ page }) => {
        await page.addInitScript(`${buildSeedScript([SEED_SCENARIO], SEED_CHARS)}`);
        await page.goto('/');
        await waitForAppReady(page);

        // Apply scenario and trigger save
        await page.evaluate((scenario) => {
            if (window.chatBridge) window.chatBridge.applyScenarioChange(scenario);
        }, SEED_SCENARIO);
        await page.waitForTimeout(400);

        // Save via the session helper
        await page.evaluate(async () => {
            if (window.sidebarSession?.saveChatData) {
                await window.sidebarSession.saveChatData();
            } else if (window.go?.main?.App?.SaveChat) {
                await window.go.main.App.SaveChat();
            }
        });
        await page.waitForTimeout(600);

        // Reload
        await page.reload();
        await waitForAppReady(page);
        await page.waitForTimeout(800);

        await assertDisplayText(page, '#scenario-name-display', {
            not:      ['unnamed', 'no scenario selected', ''],
            contains: 'Persistence Test Scenario',
        });
    });

    test('both characters survive page reload', async ({ page }) => {
        await page.addInitScript(`${buildSeedScript([SEED_SCENARIO], SEED_CHARS)}`);
        await page.goto('/');
        await waitForAppReady(page);

        await page.evaluate((scenario) => {
            if (window.chatBridge) window.chatBridge.applyScenarioChange(scenario);
        }, SEED_SCENARIO);
        await page.waitForTimeout(400);

        await page.evaluate(async () => {
            if (window.sidebarSession?.saveChatData) await window.sidebarSession.saveChatData();
            else if (window.go?.main?.App?.SaveChat)  await window.go.main.App.SaveChat();
        });
        await page.waitForTimeout(600);

        await page.reload();
        await waitForAppReady(page);
        await page.waitForTimeout(800);

        const charNames = await page.$$eval('.character-name', els => els.map(e => e.textContent.trim()));
        expect(charNames, 'Expected Alice in character panel after reload').toContain('Alice');
        expect(charNames, 'Expected Bob in character panel after reload').toContain('Bob');
    });

    test('localStorage does not contain raw base64 image blobs after save', async ({ page }) => {
        await page.addInitScript(`${buildSeedScript([SEED_SCENARIO], SEED_CHARS)}`);
        await page.goto('/');
        await waitForAppReady(page);

        await page.evaluate((scenario) => {
            if (window.chatBridge) window.chatBridge.applyScenarioChange(scenario);
        }, SEED_SCENARIO);
        await page.waitForTimeout(400);

        await page.evaluate(async () => {
            if (window.sidebarSession?.saveChatData) await window.sidebarSession.saveChatData();
        });
        await page.waitForTimeout(400);

        const blobsFound = await page.evaluate(() => {
            const results = [];
            for (let i = 0; i < localStorage.length; i++) {
                const key = localStorage.key(i);
                const val = localStorage.getItem(key) || '';
                // base64 image data URLs are > 1 KB inline
                if (val.includes('data:image') && val.length > 2000) {
                    results.push({ key, sampleLength: val.length });
                }
            }
            return results;
        });

        expect(
            blobsFound,
            `Raw base64 blobs found in localStorage — will cause QuotaExceededError: ${JSON.stringify(blobsFound)}`
        ).toHaveLength(0);
    });
});
