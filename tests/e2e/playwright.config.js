// @ts-check
const { defineConfig, devices } = require('@playwright/test');

module.exports = defineConfig({
    testDir:    './specs',
    timeout:    30_000,
    retries:    0,
    reporter:   [['list'], ['html', { outputFolder: 'playwright-report', open: 'never' }]],

    use: {
        // Wails dev server runs on this port when launched with `wails dev`
        baseURL:           'http://localhost:34115',
        headless:          false,   // set true in CI
        viewport:          { width: 1280, height: 800 },
        screenshot:        'only-on-failure',
        video:             'retain-on-failure',
        trace:             'retain-on-failure',
    },

    projects: [
        {
            name:    'chromium',
            use:     { ...devices['Desktop Chrome'] },
        },
    ],
});
