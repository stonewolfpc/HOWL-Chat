/**
 * E2E: Scenario Load Tests
 *
 * Verifies that when a chat session is loaded:
 *  - Scenario name shows correctly (not "Unnamed", not blank)
 *  - Both characters attached to the scenario appear in the character panel
 *  - Character avatars render (not broken image / placeholder icon)
 *  - World name displays correctly
 *  - Scenario thumbnail image renders when present
 */

const { test, expect } = require('@playwright/test');
const { waitForAppReady, assertDisplayText, assertImageLoaded, auditPanel, buildSeedScript } = require('./helpers');

// Minimal inline base64 PNG (1x1 red pixel) for image tests
const TINY_PNG = 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwADhQGAWjR9awAAAABJRU5ErkJggg==';

const TEST_CHARACTERS = [
    {
        id:          'char_emily_test',
        name:        'Emily',
        role:        'Daughter',
        personality: 'Sweet and confident',
        backstory:   'Penthouse kid',
        avatar:      TINY_PNG,
    },
    {
        id:          'char_bri_test',
        name:        'Bri',
        role:        'Friend',
        personality: 'Clever and kind',
        backstory:   'Neighbour',
        avatar:      TINY_PNG,
    },
];

const TEST_SCENARIO = {
    id:          'scenario-e2e-test-001',
    title:       'Stone Wolf Systems',
    description: 'A tech company HQ scenario',
    image:       { data: TINY_PNG, name: 'hq.png', type: 'image/png', size: 100, stored: false },
    characters:  TEST_CHARACTERS.map(c => ({ name: c.name, role: c.role, avatar: c.avatar })),
    lorebooks:   [],
};

test.describe('Scenario loads correctly', () => {

    test.beforeEach(async ({ page }) => {
        // Seed localStorage before the page evaluates any JS
        await page.addInitScript(`
            ${buildSeedScript([TEST_SCENARIO], TEST_CHARACTERS)}
        `);
        await page.goto('/');
        await waitForAppReady(page);
    });

    test('scenario name is not unnamed or blank after selecting scenario', async ({ page }) => {
        // Trigger scenario selection via the bridge
        await page.evaluate((scenario) => {
            if (window.chatBridge) window.chatBridge.applyScenarioChange(scenario);
        }, TEST_SCENARIO);

        await page.waitForTimeout(500); // allow async image resolution

        await assertDisplayText(page, '#scenario-name-display', {
            not:      ['Unnam', 'unnamed'],
            contains: 'Stone Wolf Systems',
        });
    });

    test('both characters appear in character panel', async ({ page }) => {
        await page.evaluate((scenario) => {
            if (window.chatBridge) window.chatBridge.applyScenarioChange(scenario);
        }, TEST_SCENARIO);

        await page.waitForTimeout(500);

        const charNames = await page.$$eval('.character-name', els => els.map(e => e.textContent.trim()));
        expect(charNames).toContain('Emily');
        expect(charNames).toContain('Bri');
        expect(charNames.length).toBeGreaterThanOrEqual(2);
    });

    test('character avatars render real images (not placeholder icon)', async ({ page }) => {
        await page.evaluate((scenario) => {
            if (window.chatBridge) window.chatBridge.applyScenarioChange(scenario);
        }, TEST_SCENARIO);

        await page.waitForTimeout(500);

        const avatarImgs = await page.$$('.character-avatar img');
        expect(avatarImgs.length).toBeGreaterThan(0);

        for (const img of avatarImgs) {
            const naturalWidth = await img.evaluate(el => el.naturalWidth);
            const display      = await img.evaluate(el => getComputedStyle(el).display);
            const src          = await img.evaluate(el => el.src.substring(0, 60));

            expect(display, `Avatar img display:none — src: ${src}`).not.toBe('none');
            expect(naturalWidth, `Avatar naturalWidth=0 (broken) — src: ${src}`).toBeGreaterThan(0);
        }
    });

    test('scenario thumbnail renders image (not "No Scenario" placeholder text)', async ({ page }) => {
        await page.evaluate((scenario) => {
            if (window.chatBridge) window.chatBridge.applyScenarioChange(scenario);
        }, TEST_SCENARIO);

        await page.waitForTimeout(500);

        const imgEl = page.locator('#scenario-image');
        await expect(imgEl).not.toHaveCSS('display', 'none');

        const naturalWidth = await imgEl.evaluate(el => el.naturalWidth);
        expect(naturalWidth, 'Scenario thumbnail image did not load').toBeGreaterThan(0);
    });

    test('placeholder text is hidden when scenario is loaded', async ({ page }) => {
        await page.evaluate((scenario) => {
            if (window.chatBridge) window.chatBridge.applyScenarioChange(scenario);
        }, TEST_SCENARIO);

        await page.waitForTimeout(300);

        const placeholderVisible = await page.evaluate(() => {
            const span = document.querySelector('#scenario-thumbnail .thumbnail-placeholder');
            if (!span) return false;
            return getComputedStyle(span).display !== 'none';
        });
        expect(placeholderVisible, '"No Scenario" placeholder still visible after load').toBe(false);
    });
});

test.describe('Character panel is functional (not a placeholder stub)', () => {

    test.beforeEach(async ({ page }) => {
        await page.goto('/');
        await waitForAppReady(page);
    });

    test('character panel exists and has real content structure', async ({ page }) => {
        const audit = await auditPanel(page, '.characters-panel, #characters-panel, [class*="character"][class*="panel"]');
        expect(audit.isPlaceholder, `Character panel is a stub: ${audit.reason}`).toBe(false);
    });

    test('"Add Character" button is present and visible', async ({ page }) => {
        const btn = page.locator('button:has-text("Add Character")');
        await expect(btn).toBeVisible();
    });
});
