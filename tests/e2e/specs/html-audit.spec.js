/**
 * E2E: HTML Structural Audit Tests
 *
 * Catches:
 *  - <label for="X"> where X doesn't exist in the DOM
 *  - <img> elements with non-empty src that fail to load (naturalWidth=0)
 *  - Panels that are empty stubs / placeholders
 *  - Sidebar sections that have no real content
 *  - Wails bindings that are expected but missing
 */

const { test, expect } = require('@playwright/test');
const { waitForAppReady, auditPanel } = require('./helpers');

const HTML_PAGES = [
    { name: 'Main app',   path: '/' },
    { name: 'Scenario',   path: '/scenario.html' },
    { name: 'Character',  path: '/character.html' },
    { name: 'World',      path: '/world.html' },
    { name: 'Lorebooks',  path: '/lorebooks.html' },
    { name: 'Settings',   path: '/settings.html' },
];

// ---------------------------------------------------------------------------
// Label / form field association audit — runs on every HTML page
// ---------------------------------------------------------------------------

for (const { name, path } of HTML_PAGES) {
    test(`[${name}] All <label for="X"> reference a real form element`, async ({ page }) => {
        await page.goto(path, { waitUntil: 'domcontentloaded' });

        const violations = await page.evaluate(() => {
            const labels = Array.from(document.querySelectorAll('label[for]'));
            return labels
                .filter(label => {
                    const target = document.getElementById(label.htmlFor);
                    return !target;
                })
                .map(label => ({
                    for:  label.htmlFor,
                    text: label.textContent.trim().substring(0, 60),
                }));
        });

        expect(
            violations,
            `[${name}] Labels with broken 'for' attribute:\n${JSON.stringify(violations, null, 2)}`
        ).toHaveLength(0);
    });
}

// ---------------------------------------------------------------------------
// Broken image audit — any <img src="non-empty"> that fails to load
// ---------------------------------------------------------------------------

test('Main app: no broken <img> elements on load', async ({ page }) => {
    await page.goto('/');
    await waitForAppReady(page);

    const broken = await page.evaluate(() => {
        return Array.from(document.querySelectorAll('img'))
            .filter(img => {
                const src = img.src;
                if (!src || src === window.location.href || src.endsWith('/')) return false;
                if (getComputedStyle(img).display === 'none') return false;
                return img.naturalWidth === 0;
            })
            .map(img => ({ src: img.src.substring(0, 100), id: img.id, className: img.className }));
    });

    expect(
        broken,
        `Broken <img> elements found (visible, non-empty src, naturalWidth=0):\n${JSON.stringify(broken, null, 2)}`
    ).toHaveLength(0);
});

// ---------------------------------------------------------------------------
// Wails bindings audit — critical bindings must be present at runtime
// ---------------------------------------------------------------------------

test('All critical Wails bindings are registered', async ({ page }) => {
    await page.goto('/');
    await waitForAppReady(page);

    const REQUIRED_BINDINGS = [
        'SaveChat',
        'LoadChat',
        'AddScenarioToChat',
        'AddCharacterToChat',
        'GetChatList',
        'NewChat',
    ];

    const missing = await page.evaluate((bindings) => {
        return bindings.filter(name => {
            return typeof window.go?.main?.App?.[name] !== 'function';
        });
    }, REQUIRED_BINDINGS);

    expect(
        missing,
        `Missing Wails bindings — these will silently fail at runtime: ${missing.join(', ')}`
    ).toHaveLength(0);
});

// ---------------------------------------------------------------------------
// Sidebar panel audit — panels must not be empty stubs
// ---------------------------------------------------------------------------

test('Character panel is not a placeholder stub', async ({ page }) => {
    await page.goto('/');
    await waitForAppReady(page);

    // Try several possible selectors for the character panel
    const characterPanelSelectors = [
        '.sidebar-panel[data-panel="characters"]',
        '#character-list-panel',
        '.characters-panel',
        '[id*="character"][class*="panel"]',
    ];

    let found = false;
    for (const sel of characterPanelSelectors) {
        const count = await page.locator(sel).count();
        if (count > 0) {
            const audit = await auditPanel(page, sel);
            expect(audit.isPlaceholder, `Character panel stub detected: ${audit.reason}`).toBe(false);
            found = true;
            break;
        }
    }

    // If we can't find a character panel selector, fail with a diagnostic
    if (!found) {
        const allIds = await page.evaluate(() =>
            Array.from(document.querySelectorAll('[id*="character"]')).map(e => e.id).slice(0, 20)
        );
        throw new Error(`Could not find character panel. Elements with 'character' in id: ${allIds.join(', ')}`);
    }
});

test('Scenario context panel is not a placeholder stub', async ({ page }) => {
    await page.goto('/');
    await waitForAppReady(page);

    const el = page.locator('#scenario-name-display');
    await expect(el).toBeVisible();
});

test('World context panel is not a placeholder stub', async ({ page }) => {
    await page.goto('/');
    await waitForAppReady(page);

    const el = page.locator('#world-name-display');
    await expect(el).toBeVisible();
});

// ---------------------------------------------------------------------------
// localStorage quota audit
// ---------------------------------------------------------------------------

test('localStorage usage is under 4 MB after app init', async ({ page }) => {
    await page.goto('/');
    await waitForAppReady(page);

    const stats = await page.evaluate(() => {
        let total = 0;
        const keys = {};
        for (let i = 0; i < localStorage.length; i++) {
            const key = localStorage.key(i);
            const val = localStorage.getItem(key) || '';
            keys[key] = Math.round(val.length / 1024) + ' KB';
            total += val.length;
        }
        return { totalKB: total / 1024, keys };
    });

    expect(
        stats.totalKB,
        `localStorage usage ${stats.totalKB.toFixed(1)} KB exceeds 4 MB — risk of QuotaExceededError.\nBreakdown: ${JSON.stringify(stats.keys, null, 2)}`
    ).toBeLessThan(4096);
});
