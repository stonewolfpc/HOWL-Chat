# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: html-audit.spec.js >> Character panel is not a placeholder stub
- Location: specs\html-audit.spec.js:110:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/", waiting until "load"

```

# Test source

```ts
  11  | 
  12  | const { test, expect } = require('@playwright/test');
  13  | const { waitForAppReady, auditPanel } = require('./helpers');
  14  | 
  15  | const HTML_PAGES = [
  16  |     { name: 'Main app',   path: '/' },
  17  |     { name: 'Scenario',   path: '/scenario.html' },
  18  |     { name: 'Character',  path: '/character.html' },
  19  |     { name: 'World',      path: '/world.html' },
  20  |     { name: 'Lorebooks',  path: '/lorebooks.html' },
  21  |     { name: 'Settings',   path: '/settings.html' },
  22  | ];
  23  | 
  24  | // ---------------------------------------------------------------------------
  25  | // Label / form field association audit — runs on every HTML page
  26  | // ---------------------------------------------------------------------------
  27  | 
  28  | for (const { name, path } of HTML_PAGES) {
  29  |     test(`[${name}] All <label for="X"> reference a real form element`, async ({ page }) => {
  30  |         await page.goto(path, { waitUntil: 'domcontentloaded' });
  31  | 
  32  |         const violations = await page.evaluate(() => {
  33  |             const labels = Array.from(document.querySelectorAll('label[for]'));
  34  |             return labels
  35  |                 .filter(label => {
  36  |                     const target = document.getElementById(label.htmlFor);
  37  |                     return !target;
  38  |                 })
  39  |                 .map(label => ({
  40  |                     for:  label.htmlFor,
  41  |                     text: label.textContent.trim().substring(0, 60),
  42  |                 }));
  43  |         });
  44  | 
  45  |         expect(
  46  |             violations,
  47  |             `[${name}] Labels with broken 'for' attribute:\n${JSON.stringify(violations, null, 2)}`
  48  |         ).toHaveLength(0);
  49  |     });
  50  | }
  51  | 
  52  | // ---------------------------------------------------------------------------
  53  | // Broken image audit — any <img src="non-empty"> that fails to load
  54  | // ---------------------------------------------------------------------------
  55  | 
  56  | test('Main app: no broken <img> elements on load', async ({ page }) => {
  57  |     await page.goto('/');
  58  |     await waitForAppReady(page);
  59  | 
  60  |     const broken = await page.evaluate(() => {
  61  |         return Array.from(document.querySelectorAll('img'))
  62  |             .filter(img => {
  63  |                 const src = img.src;
  64  |                 if (!src || src === window.location.href || src.endsWith('/')) return false;
  65  |                 if (getComputedStyle(img).display === 'none') return false;
  66  |                 return img.naturalWidth === 0;
  67  |             })
  68  |             .map(img => ({ src: img.src.substring(0, 100), id: img.id, className: img.className }));
  69  |     });
  70  | 
  71  |     expect(
  72  |         broken,
  73  |         `Broken <img> elements found (visible, non-empty src, naturalWidth=0):\n${JSON.stringify(broken, null, 2)}`
  74  |     ).toHaveLength(0);
  75  | });
  76  | 
  77  | // ---------------------------------------------------------------------------
  78  | // Wails bindings audit — critical bindings must be present at runtime
  79  | // ---------------------------------------------------------------------------
  80  | 
  81  | test('All critical Wails bindings are registered', async ({ page }) => {
  82  |     await page.goto('/');
  83  |     await waitForAppReady(page);
  84  | 
  85  |     const REQUIRED_BINDINGS = [
  86  |         'SaveChat',
  87  |         'LoadChat',
  88  |         'AddScenarioToChat',
  89  |         'AddCharacterToChat',
  90  |         'GetChatList',
  91  |         'NewChat',
  92  |     ];
  93  | 
  94  |     const missing = await page.evaluate((bindings) => {
  95  |         return bindings.filter(name => {
  96  |             return typeof window.go?.main?.App?.[name] !== 'function';
  97  |         });
  98  |     }, REQUIRED_BINDINGS);
  99  | 
  100 |     expect(
  101 |         missing,
  102 |         `Missing Wails bindings — these will silently fail at runtime: ${missing.join(', ')}`
  103 |     ).toHaveLength(0);
  104 | });
  105 | 
  106 | // ---------------------------------------------------------------------------
  107 | // Sidebar panel audit — panels must not be empty stubs
  108 | // ---------------------------------------------------------------------------
  109 | 
  110 | test('Character panel is not a placeholder stub', async ({ page }) => {
> 111 |     await page.goto('/');
      |                ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  112 |     await waitForAppReady(page);
  113 | 
  114 |     // Try several possible selectors for the character panel
  115 |     const characterPanelSelectors = [
  116 |         '.sidebar-panel[data-panel="characters"]',
  117 |         '#character-list-panel',
  118 |         '.characters-panel',
  119 |         '[id*="character"][class*="panel"]',
  120 |     ];
  121 | 
  122 |     let found = false;
  123 |     for (const sel of characterPanelSelectors) {
  124 |         const count = await page.locator(sel).count();
  125 |         if (count > 0) {
  126 |             const audit = await auditPanel(page, sel);
  127 |             expect(audit.isPlaceholder, `Character panel stub detected: ${audit.reason}`).toBe(false);
  128 |             found = true;
  129 |             break;
  130 |         }
  131 |     }
  132 | 
  133 |     // If we can't find a character panel selector, fail with a diagnostic
  134 |     if (!found) {
  135 |         const allIds = await page.evaluate(() =>
  136 |             Array.from(document.querySelectorAll('[id*="character"]')).map(e => e.id).slice(0, 20)
  137 |         );
  138 |         throw new Error(`Could not find character panel. Elements with 'character' in id: ${allIds.join(', ')}`);
  139 |     }
  140 | });
  141 | 
  142 | test('Scenario context panel is not a placeholder stub', async ({ page }) => {
  143 |     await page.goto('/');
  144 |     await waitForAppReady(page);
  145 | 
  146 |     const el = page.locator('#scenario-name-display');
  147 |     await expect(el).toBeVisible();
  148 | });
  149 | 
  150 | test('World context panel is not a placeholder stub', async ({ page }) => {
  151 |     await page.goto('/');
  152 |     await waitForAppReady(page);
  153 | 
  154 |     const el = page.locator('#world-name-display');
  155 |     await expect(el).toBeVisible();
  156 | });
  157 | 
  158 | // ---------------------------------------------------------------------------
  159 | // localStorage quota audit
  160 | // ---------------------------------------------------------------------------
  161 | 
  162 | test('localStorage usage is under 4 MB after app init', async ({ page }) => {
  163 |     await page.goto('/');
  164 |     await waitForAppReady(page);
  165 | 
  166 |     const stats = await page.evaluate(() => {
  167 |         let total = 0;
  168 |         const keys = {};
  169 |         for (let i = 0; i < localStorage.length; i++) {
  170 |             const key = localStorage.key(i);
  171 |             const val = localStorage.getItem(key) || '';
  172 |             keys[key] = Math.round(val.length / 1024) + ' KB';
  173 |             total += val.length;
  174 |         }
  175 |         return { totalKB: total / 1024, keys };
  176 |     });
  177 | 
  178 |     expect(
  179 |         stats.totalKB,
  180 |         `localStorage usage ${stats.totalKB.toFixed(1)} KB exceeds 4 MB — risk of QuotaExceededError.\nBreakdown: ${JSON.stringify(stats.keys, null, 2)}`
  181 |     ).toBeLessThan(4096);
  182 | });
  183 | 
```