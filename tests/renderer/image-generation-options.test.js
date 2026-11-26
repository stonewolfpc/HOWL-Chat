/**
 * Generation request identity tests.
 *
 * A generation request is confirmed once. The LoRA set used for the fit check,
 * for the unload re-admission, and for the backend job must therefore be one
 * frozen snapshot, not a live reference to window.activeLoRAs.
 */

const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');

const CORE = path.resolve(__dirname, '../../frontend/dist/js/image-generator-core.js');

function loadCore() {
    const source = fs.readFileSync(CORE, 'utf8');
    const window = {
        addEventListener() {},
        currentModelFamily: '',
    };
    window.window = window;

    const fields = {
        'ig-prompt': { value: 'a prompt' },
        'ig-model': { value: 'model-a' },
        'ig-generate': {
            disabled: false,
            innerHTML: '',
            classList: { add() {}, remove() {}, toggle() {} },
        },
        'ig-negative': { value: '' },
        'ig-width': { value: '512' },
        'ig-height': { value: '512' },
        'ig-steps': { value: '4' },
        'ig-seed': { value: '-1' },
        'ig-scheduler': { value: 'normal' },
        'ig-cfg': { value: '7' },
        'ig-batch': { value: '1' },
        'ig-hires-enable': { checked: false },
    };
    const document = {
        addEventListener() {},
        getElementById(id) {
            return fields[id] || { value: '', checked: false, textContent: '' };
        },
        querySelector() { return null; },
    };

    const sandbox = {
        window,
        document,
        console: { log() {}, warn() {}, error() {} },
        alert() {},
        confirm() { return true; },
        setTimeout() { return 1; },
        clearTimeout() {},
        setInterval(fn, ms) {
            if (ms === 500) queueMicrotask(fn);
            return 1;
        },
        clearInterval() {},
    };
    vm.createContext(sandbox);
    vm.runInContext(source, sandbox, { filename: CORE });

    // Keep the test focused on request identity: the DOM/progress seams are
    // not part of the workload being checked.
    // models, lastUsedSeed and seedHistory are lexical bindings in the core
    // module, so they must be assigned inside the same VM context.
    vm.runInContext(
        'models = [{ name: "model-a", files: {} }]; currentModelFamily = ""; lastUsedSeed = 1; seedHistory = [];',
        sandbox,
    );
    sandbox.currentMode = 'txt2img';
    sandbox.seedLocked = false;
    sandbox.readSamplerValue = () => 'euler';
    sandbox.setStatus = () => {};
    sandbox.showSpinner = () => {};
    sandbox.hideSpinner = () => {};
    sandbox.startProgressBar = () => {};
    sandbox.stopProgressBar = () => {};
    sandbox.updateProgressBar = () => {};
    sandbox.updateDiagnostics = () => {};
    sandbox.displayResult = () => {};
    sandbox.resetGenerateButton = () => {};
    return { window, sandbox };
}

test('generation LoRAs are snapshotted and survive live-array mutation', () => {
    const { window } = loadCore();
    const live = [
        { name: 'alpha', strength: 0.75, enabled: true },
        { name: 'beta', strength: -0.5, enabled: false },
    ];
    const snapshot = window.snapshotGenerationLoRAs(live);

    live[0].name = 'mutated';
    live[1].strength = 2;
    live.push({ name: 'late', strength: 1, enabled: true });

    assert.deepEqual(JSON.parse(JSON.stringify(snapshot)), [
        { name: 'alpha', strength: 0.75, enabled: true },
        { name: 'beta', strength: -0.5, enabled: false },
    ]);
    assert.throws(() => snapshot.push({ name: 'x', strength: 1, enabled: true }), /TypeError/);
});

test('preflight, unload re-admission and backend job receive one frozen request', async () => {
    const { window, sandbox } = loadCore();
    const live = [{ name: 'alpha', strength: 0.75, enabled: true }];
    window.activeLoRAs = live;

    const seen = [];
    let fitChecks = 0;
    let confirmMessage = null;
    sandbox.confirm = message => {
        confirmMessage = message;
        return true;
    };
    sandbox.window.go = {
        main: {
            App: {
                async CheckImageFit(opts) {
                    seen.push({ stage: 'preflight', opts });
                    fitChecks++;
                    // Mutate the source array in flight, exactly as a UI edit
                    // can while the confirmation dialog is open.
                    live[0].name = 'mutated-in-flight';
                    live.push({ name: 'late', strength: 1, enabled: true });
                    return { can_fit: false, blocking: ['chat'], message: 'Chat is using the GPU.' };
                },
                async FreeVRAMForImage(opts, confirmed) {
                    seen.push({ stage: 'unload', opts, confirmed });
                    return { can_fit: true };
                },
                async StartImageJob(opts) {
                    seen.push({ stage: 'backend', opts });
                    return 'job-1';
                },
                async GetImageJobStatus() {
                    return { status: 'completed', result: {} };
                },
            },
        },
    };

    await vm.runInContext('onGenerate()', sandbox);

    assert.equal(fitChecks, 1);
    assert.equal(confirmMessage, 'Chat is using the GPU.');
    assert.deepEqual(seen.map(entry => entry.stage), ['preflight', 'unload', 'backend']);
    assert.strictEqual(seen[0].opts, seen[1].opts);
    assert.strictEqual(seen[1].opts, seen[2].opts);
    // The unload is handed the user's answer, not a constant.
    assert.equal(seen[1].confirmed, true);
    assert.deepEqual(JSON.parse(JSON.stringify(seen[0].opts.loras)), [{ name: 'alpha', strength: 0.75, enabled: true }]);
    assert.deepEqual(JSON.parse(JSON.stringify(live[0])), { name: 'mutated-in-flight', strength: 0.75, enabled: true });
    assert.equal(live.length, 2);
});

test('a declined confirmation never unloads and never queues the job', async () => {
    const { window, sandbox } = loadCore();
    let unloads = 0;
    let jobs = 0;
    sandbox.confirm = () => false;
    window.go = {
        main: {
            App: {
                async CheckImageFit() {
                    return { can_fit: false, blocking: ['chat'], message: 'Chat is using the GPU.' };
                },
                async FreeVRAMForImage() {
                    unloads++;
                    return { can_fit: true };
                },
                async StartImageJob() {
                    jobs++;
                    return 'job-1';
                },
            },
        },
    };

    const allowed = await vm.runInContext(
        'requestImageVRAMArbitration({ prompt: "p", model_name: "model-a" })',
        sandbox,
    );

    assert.equal(allowed, false);
    assert.equal(unloads, 0, 'a declined dialog must not unload the chat model');
    assert.equal(jobs, 0, 'a declined dialog must not queue a generation');
});