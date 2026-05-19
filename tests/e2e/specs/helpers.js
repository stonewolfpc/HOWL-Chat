/**
 * Shared helpers for HOWL Chat E2E tests.
 * Provides utilities for seeding localStorage, waiting for Wails runtime,
 * and asserting image visibility.
 */

const SCENARIO_KEY  = 'howl-chat-scenarios';
const CHARACTER_KEY = 'howl-chat-characters';
const WORLD_KEY     = 'howl-worlds';

/**
 * Wait for the Wails runtime to signal readiness.
 * Falls back to waiting for the chat container to be visible.
 */
async function waitForAppReady(page) {
    await page.waitForFunction(() => {
        return typeof window.go !== 'undefined' && window.go?.main?.App;
    }, { timeout: 15_000 });
    await page.waitForSelector('.chat-messages', { timeout: 10_000 });
}

/**
 * Seed localStorage with test scenarios and characters before the page loads.
 * Must be called inside page.addInitScript or before navigation.
 */
function buildSeedScript(scenarios, characters, worlds = []) {
    return `
        localStorage.setItem(${JSON.stringify(SCENARIO_KEY)},  ${JSON.stringify(JSON.stringify(scenarios))});
        localStorage.setItem(${JSON.stringify(CHARACTER_KEY)}, ${JSON.stringify(JSON.stringify(characters))});
        localStorage.setItem(${JSON.stringify(WORLD_KEY)},     ${JSON.stringify(JSON.stringify(worlds))});
    `;
}

/**
 * Assert that an <img> element is actually visible and its naturalWidth > 0
 * (i.e. it loaded a real image, not a broken src).
 */
async function assertImageLoaded(page, selector) {
    const result = await page.evaluate((sel) => {
        const img = document.querySelector(sel);
        if (!img) return { found: false };
        return {
            found:        true,
            display:      getComputedStyle(img).display,
            naturalWidth: img.naturalWidth,
            src:          img.src.substring(0, 80),
        };
    }, selector);

    if (!result.found)   throw new Error(`Image element not found: ${selector}`);
    if (result.display === 'none') throw new Error(`Image is display:none — not visible: ${selector}`);
    if (result.naturalWidth === 0) throw new Error(`Image has naturalWidth=0 (broken/empty src): ${selector} src=${result.src}`);
    return result;
}

/**
 * Assert an element is visible AND contains non-placeholder text.
 * Placeholder values: '', 'No Scenario Selected', 'No World Selected', 'Unnamed', 'undefined'.
 */
async function assertDisplayText(page, selector, opts = {}) {
    const { not: notValues = [], contains } = opts;
    const defaultBad = ['', 'No Scenario Selected', 'No World Selected', 'unnamed', 'undefined', 'null'];
    const badValues  = [...defaultBad, ...notValues].map(v => v.toLowerCase());

    const text = await page.textContent(selector);
    const trimmed = (text || '').trim();

    if (badValues.includes(trimmed.toLowerCase())) {
        throw new Error(`Element ${selector} shows placeholder/bad value: "${trimmed}"`);
    }
    if (contains && !trimmed.toLowerCase().includes(contains.toLowerCase())) {
        throw new Error(`Element ${selector} text "${trimmed}" does not contain expected "${contains}"`);
    }
    return trimmed;
}

/**
 * Get the current localStorage quota usage summary.
 */
async function getStorageStats(page) {
    return page.evaluate(() => {
        let total = 0;
        const keys = {};
        for (let i = 0; i < localStorage.length; i++) {
            const key = localStorage.key(i);
            const val = localStorage.getItem(key) || '';
            keys[key] = val.length;
            total += val.length;
        }
        return { totalChars: total, totalKB: (total / 1024).toFixed(1), keys };
    });
}

/**
 * Check whether a panel element is a real panel or a placeholder stub.
 * Returns { isPlaceholder: bool, reason: string }
 */
async function auditPanel(page, panelSelector) {
    return page.evaluate((sel) => {
        const el = document.querySelector(sel);
        if (!el) return { isPlaceholder: true, reason: 'element not found' };

        const text     = el.textContent || '';
        const children = el.children.length;
        const hasImages = el.querySelectorAll('img').length;
        const placeholderPhrases = ['coming soon', 'todo', 'placeholder', 'not implemented', 'tbd'];
        const lowerText = text.toLowerCase();
        const foundPhrase = placeholderPhrases.find(p => lowerText.includes(p));

        if (foundPhrase) return { isPlaceholder: true, reason: `contains "${foundPhrase}"` };
        if (children === 0 && text.trim() === '') return { isPlaceholder: true, reason: 'empty panel' };

        return { isPlaceholder: false, hasImages, childCount: children, textLength: text.length };
    }, panelSelector);
}

module.exports = { waitForAppReady, buildSeedScript, assertImageLoaded, assertDisplayText, getStorageStats, auditPanel, SCENARIO_KEY, CHARACTER_KEY, WORLD_KEY };
