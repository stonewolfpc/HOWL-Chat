// Headless config for the file-based audio-settings lifecycle spec. It uses the
// Edge build installed on this machine, so the spec does not need a running
// Wails dev server.
const { defineConfig } = require('@playwright/test');

module.exports = defineConfig({
    testDir: './specs',
    timeout: 30_000,
    retries: 0,
    reporter: [['list']],
    use: {
        headless: true,
        viewport: { width: 1280, height: 800 },
        channel: 'msedge',
    },
});