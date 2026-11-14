# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: scenario-load.spec.js >> Scenario loads correctly >> scenario name is not unnamed or blank after selecting scenario
- Location: specs\scenario-load.spec.js:57:5

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/", waiting until "load"

```

# Test source

```ts
  1   | /**
  2   |  * E2E: Scenario Load Tests
  3   |  *
  4   |  * Verifies that when a chat session is loaded:
  5   |  *  - Scenario name shows correctly (not "Unnamed", not blank)
  6   |  *  - Both characters attached to the scenario appear in the character panel
  7   |  *  - Character avatars render (not broken image / placeholder icon)
  8   |  *  - World name displays correctly
  9   |  *  - Scenario thumbnail image renders when present
  10  |  */
  11  | 
  12  | const { test, expect } = require('@playwright/test');
  13  | const { waitForAppReady, assertDisplayText, assertImageLoaded, auditPanel, buildSeedScript } = require('./helpers');
  14  | 
  15  | // Minimal inline base64 PNG (1x1 red pixel) for image tests
  16  | const TINY_PNG = 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwADhQGAWjR9awAAAABJRU5ErkJggg==';
  17  | 
  18  | const TEST_CHARACTERS = [
  19  |     {
  20  |         id:          'char_emily_test',
  21  |         name:        'Emily',
  22  |         role:        'Daughter',
  23  |         personality: 'Sweet and confident',
  24  |         backstory:   'Penthouse kid',
  25  |         avatar:      TINY_PNG,
  26  |     },
  27  |     {
  28  |         id:          'char_bri_test',
  29  |         name:        'Bri',
  30  |         role:        'Friend',
  31  |         personality: 'Clever and kind',
  32  |         backstory:   'Neighbour',
  33  |         avatar:      TINY_PNG,
  34  |     },
  35  | ];
  36  | 
  37  | const TEST_SCENARIO = {
  38  |     id:          'scenario-e2e-test-001',
  39  |     title:       'Stone Wolf Systems',
  40  |     description: 'A tech company HQ scenario',
  41  |     image:       { data: TINY_PNG, name: 'hq.png', type: 'image/png', size: 100, stored: false },
  42  |     characters:  TEST_CHARACTERS.map(c => ({ name: c.name, role: c.role, avatar: c.avatar })),
  43  |     lorebooks:   [],
  44  | };
  45  | 
  46  | test.describe('Scenario loads correctly', () => {
  47  | 
  48  |     test.beforeEach(async ({ page }) => {
  49  |         // Seed localStorage before the page evaluates any JS
  50  |         await page.addInitScript(`
  51  |             ${buildSeedScript([TEST_SCENARIO], TEST_CHARACTERS)}
  52  |         `);
> 53  |         await page.goto('/');
      |                    ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  54  |         await waitForAppReady(page);
  55  |     });
  56  | 
  57  |     test('scenario name is not unnamed or blank after selecting scenario', async ({ page }) => {
  58  |         // Trigger scenario selection via the bridge
  59  |         await page.evaluate((scenario) => {
  60  |             if (window.chatBridge) window.chatBridge.applyScenarioChange(scenario);
  61  |         }, TEST_SCENARIO);
  62  | 
  63  |         await page.waitForTimeout(500); // allow async image resolution
  64  | 
  65  |         await assertDisplayText(page, '#scenario-name-display', {
  66  |             not:      ['Unnam', 'unnamed'],
  67  |             contains: 'Stone Wolf Systems',
  68  |         });
  69  |     });
  70  | 
  71  |     test('both characters appear in character panel', async ({ page }) => {
  72  |         await page.evaluate((scenario) => {
  73  |             if (window.chatBridge) window.chatBridge.applyScenarioChange(scenario);
  74  |         }, TEST_SCENARIO);
  75  | 
  76  |         await page.waitForTimeout(500);
  77  | 
  78  |         const charNames = await page.$$eval('.character-name', els => els.map(e => e.textContent.trim()));
  79  |         expect(charNames).toContain('Emily');
  80  |         expect(charNames).toContain('Bri');
  81  |         expect(charNames.length).toBeGreaterThanOrEqual(2);
  82  |     });
  83  | 
  84  |     test('character avatars render real images (not placeholder icon)', async ({ page }) => {
  85  |         await page.evaluate((scenario) => {
  86  |             if (window.chatBridge) window.chatBridge.applyScenarioChange(scenario);
  87  |         }, TEST_SCENARIO);
  88  | 
  89  |         await page.waitForTimeout(500);
  90  | 
  91  |         const avatarImgs = await page.$$('.character-avatar img');
  92  |         expect(avatarImgs.length).toBeGreaterThan(0);
  93  | 
  94  |         for (const img of avatarImgs) {
  95  |             const naturalWidth = await img.evaluate(el => el.naturalWidth);
  96  |             const display      = await img.evaluate(el => getComputedStyle(el).display);
  97  |             const src          = await img.evaluate(el => el.src.substring(0, 60));
  98  | 
  99  |             expect(display, `Avatar img display:none — src: ${src}`).not.toBe('none');
  100 |             expect(naturalWidth, `Avatar naturalWidth=0 (broken) — src: ${src}`).toBeGreaterThan(0);
  101 |         }
  102 |     });
  103 | 
  104 |     test('scenario thumbnail renders image (not "No Scenario" placeholder text)', async ({ page }) => {
  105 |         await page.evaluate((scenario) => {
  106 |             if (window.chatBridge) window.chatBridge.applyScenarioChange(scenario);
  107 |         }, TEST_SCENARIO);
  108 | 
  109 |         await page.waitForTimeout(500);
  110 | 
  111 |         const imgEl = page.locator('#scenario-image');
  112 |         await expect(imgEl).not.toHaveCSS('display', 'none');
  113 | 
  114 |         const naturalWidth = await imgEl.evaluate(el => el.naturalWidth);
  115 |         expect(naturalWidth, 'Scenario thumbnail image did not load').toBeGreaterThan(0);
  116 |     });
  117 | 
  118 |     test('placeholder text is hidden when scenario is loaded', async ({ page }) => {
  119 |         await page.evaluate((scenario) => {
  120 |             if (window.chatBridge) window.chatBridge.applyScenarioChange(scenario);
  121 |         }, TEST_SCENARIO);
  122 | 
  123 |         await page.waitForTimeout(300);
  124 | 
  125 |         const placeholderVisible = await page.evaluate(() => {
  126 |             const span = document.querySelector('#scenario-thumbnail .thumbnail-placeholder');
  127 |             if (!span) return false;
  128 |             return getComputedStyle(span).display !== 'none';
  129 |         });
  130 |         expect(placeholderVisible, '"No Scenario" placeholder still visible after load').toBe(false);
  131 |     });
  132 | });
  133 | 
  134 | test.describe('Character panel is functional (not a placeholder stub)', () => {
  135 | 
  136 |     test.beforeEach(async ({ page }) => {
  137 |         await page.goto('/');
  138 |         await waitForAppReady(page);
  139 |     });
  140 | 
  141 |     test('character panel exists and has real content structure', async ({ page }) => {
  142 |         const audit = await auditPanel(page, '.characters-panel, #characters-panel, [class*="character"][class*="panel"]');
  143 |         expect(audit.isPlaceholder, `Character panel is a stub: ${audit.reason}`).toBe(false);
  144 |     });
  145 | 
  146 |     test('"Add Character" button is present and visible', async ({ page }) => {
  147 |         const btn = page.locator('button:has-text("Add Character")');
  148 |         await expect(btn).toBeVisible();
  149 |     });
  150 | });
  151 | 
```