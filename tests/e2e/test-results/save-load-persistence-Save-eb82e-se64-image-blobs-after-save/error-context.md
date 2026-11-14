# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: save-load-persistence.spec.js >> Save → reload persistence >> localStorage does not contain raw base64 image blobs after save
- Location: specs\save-load-persistence.spec.js:87:5

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/", waiting until "load"

```

# Test source

```ts
  1   | /**
  2   |  * E2E: Save → Reload Persistence Tests
  3   |  *
  4   |  * Verifies that chat state (scenario, characters, world) survives
  5   |  * a full page reload and that the backend correctly persists and
  6   |  * restores it.
  7   |  */
  8   | 
  9   | const { test, expect } = require('@playwright/test');
  10  | const { waitForAppReady, assertDisplayText, buildSeedScript } = require('./helpers');
  11  | 
  12  | const TINY_PNG = 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwADhQGAWjR9awAAAABJRU5ErkJggg==';
  13  | 
  14  | const SEED_SCENARIO = {
  15  |     id:         'scenario-persist-001',
  16  |     title:      'Persistence Test Scenario',
  17  |     characters: [
  18  |         { name: 'Alice', role: 'hero',     avatar: TINY_PNG },
  19  |         { name: 'Bob',   role: 'sidekick', avatar: TINY_PNG },
  20  |     ],
  21  | };
  22  | 
  23  | const SEED_CHARS = [
  24  |     { id: 'char_persist_alice', name: 'Alice', personality: 'Brave',  avatar: TINY_PNG },
  25  |     { id: 'char_persist_bob',   name: 'Bob',   personality: 'Clever', avatar: TINY_PNG },
  26  | ];
  27  | 
  28  | test.describe('Save → reload persistence', () => {
  29  | 
  30  |     test('scenario name survives page reload', async ({ page }) => {
  31  |         await page.addInitScript(`${buildSeedScript([SEED_SCENARIO], SEED_CHARS)}`);
  32  |         await page.goto('/');
  33  |         await waitForAppReady(page);
  34  | 
  35  |         // Apply scenario and trigger save
  36  |         await page.evaluate((scenario) => {
  37  |             if (window.chatBridge) window.chatBridge.applyScenarioChange(scenario);
  38  |         }, SEED_SCENARIO);
  39  |         await page.waitForTimeout(400);
  40  | 
  41  |         // Save via the session helper
  42  |         await page.evaluate(async () => {
  43  |             if (window.sidebarSession?.saveChatData) {
  44  |                 await window.sidebarSession.saveChatData();
  45  |             } else if (window.go?.main?.App?.SaveChat) {
  46  |                 await window.go.main.App.SaveChat();
  47  |             }
  48  |         });
  49  |         await page.waitForTimeout(600);
  50  | 
  51  |         // Reload
  52  |         await page.reload();
  53  |         await waitForAppReady(page);
  54  |         await page.waitForTimeout(800);
  55  | 
  56  |         await assertDisplayText(page, '#scenario-name-display', {
  57  |             not:      ['unnamed', 'no scenario selected', ''],
  58  |             contains: 'Persistence Test Scenario',
  59  |         });
  60  |     });
  61  | 
  62  |     test('both characters survive page reload', async ({ page }) => {
  63  |         await page.addInitScript(`${buildSeedScript([SEED_SCENARIO], SEED_CHARS)}`);
  64  |         await page.goto('/');
  65  |         await waitForAppReady(page);
  66  | 
  67  |         await page.evaluate((scenario) => {
  68  |             if (window.chatBridge) window.chatBridge.applyScenarioChange(scenario);
  69  |         }, SEED_SCENARIO);
  70  |         await page.waitForTimeout(400);
  71  | 
  72  |         await page.evaluate(async () => {
  73  |             if (window.sidebarSession?.saveChatData) await window.sidebarSession.saveChatData();
  74  |             else if (window.go?.main?.App?.SaveChat)  await window.go.main.App.SaveChat();
  75  |         });
  76  |         await page.waitForTimeout(600);
  77  | 
  78  |         await page.reload();
  79  |         await waitForAppReady(page);
  80  |         await page.waitForTimeout(800);
  81  | 
  82  |         const charNames = await page.$$eval('.character-name', els => els.map(e => e.textContent.trim()));
  83  |         expect(charNames, 'Expected Alice in character panel after reload').toContain('Alice');
  84  |         expect(charNames, 'Expected Bob in character panel after reload').toContain('Bob');
  85  |     });
  86  | 
  87  |     test('localStorage does not contain raw base64 image blobs after save', async ({ page }) => {
  88  |         await page.addInitScript(`${buildSeedScript([SEED_SCENARIO], SEED_CHARS)}`);
> 89  |         await page.goto('/');
      |                    ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  90  |         await waitForAppReady(page);
  91  | 
  92  |         await page.evaluate((scenario) => {
  93  |             if (window.chatBridge) window.chatBridge.applyScenarioChange(scenario);
  94  |         }, SEED_SCENARIO);
  95  |         await page.waitForTimeout(400);
  96  | 
  97  |         await page.evaluate(async () => {
  98  |             if (window.sidebarSession?.saveChatData) await window.sidebarSession.saveChatData();
  99  |         });
  100 |         await page.waitForTimeout(400);
  101 | 
  102 |         const blobsFound = await page.evaluate(() => {
  103 |             const results = [];
  104 |             for (let i = 0; i < localStorage.length; i++) {
  105 |                 const key = localStorage.key(i);
  106 |                 const val = localStorage.getItem(key) || '';
  107 |                 // base64 image data URLs are > 1 KB inline
  108 |                 if (val.includes('data:image') && val.length > 2000) {
  109 |                     results.push({ key, sampleLength: val.length });
  110 |                 }
  111 |             }
  112 |             return results;
  113 |         });
  114 | 
  115 |         expect(
  116 |             blobsFound,
  117 |             `Raw base64 blobs found in localStorage — will cause QuotaExceededError: ${JSON.stringify(blobsFound)}`
  118 |         ).toHaveLength(0);
  119 |     });
  120 | });
  121 | 
```