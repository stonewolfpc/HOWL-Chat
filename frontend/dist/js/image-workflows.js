/**
 * image-workflows.js — saved generation recipes for the Image Generator.
 *
 * A workflow is every control on this page, the mode it belongs to, the LoRA
 * stack, and (optionally) a pose skeleton. It is stored as one JSON file by the
 * Go backend so recipes survive updates and can be copied between machines.
 *
 * stable-diffusion.cpp is not a graph engine, so this is deliberately a recipe
 * and not a ComfyUI-style node graph: it restores exactly what the UI can do.
 */
(function () {
    'use strict';

    const FORMAT = 'howl-workflow';
    const VERSION = 1;

    // Fields the recipe must not capture: the workflow picker itself, file
    // inputs, and the folder manager (all local UI state, not generation state).
    const SKIP_PREFIXES = ['ig-workflow', 'ig-wf-', 'ig-model-dirs', 'ig-dir-', 'ig-add-model-dir'];
    const SKIP_IDS = new Set(['ig-upload-input', 'ig-model-dir-input', 'ig-upload-controlnet-file']);

    function shouldSkip(id) {
        if (!id || SKIP_IDS.has(id)) return true;
        return SKIP_PREFIXES.some(p => id.startsWith(p));
    }

    function fieldElements() {
        return document.querySelectorAll('.ig-panel input[id], .ig-panel select[id], .ig-panel textarea[id]');
    }

    /** Capture every control as {type, value} so apply can restore it faithfully. */
    function collectFields() {
        const out = {};
        fieldElements().forEach(el => {
            const id = el.id;
            if (shouldSkip(id) || el.type === 'file' || el.disabled) return;
            if (el.type === 'checkbox') { out[id] = { t: 'c', v: !!el.checked }; return; }
            if (el.tagName === 'SELECT') { out[id] = { t: 's', v: el.value }; return; }
            if (el.type === 'number' || el.type === 'range') { out[id] = { t: 'n', v: el.value }; return; }
            out[id] = { t: 'v', v: el.value };
        });
        return out;
    }

    /** Restore one field; selects only adopt a value they actually offer. */
    function applyField(id, entry) {
        const el = document.getElementById(id);
        if (!el || !entry) return false;
        if (el.type === 'checkbox') {
            el.checked = !!entry.v;
        } else if (el.tagName === 'SELECT') {
            const wanted = String(entry.v ?? '');
            if (wanted === '' || [...el.options].some(o => o.value === wanted)) el.value = wanted;
            else return false;
        } else {
            el.value = entry.v;
        }
        // Sliders update their value labels and the engine preview elsewhere via
        // these listeners, so fire them the way a real drag would.
        el.dispatchEvent(new Event('input', { bubbles: true }));
        el.dispatchEvent(new Event('change', { bubbles: true }));
        return true;
    }

    function captureSettings() {
        return {
            format: FORMAT,
            version: VERSION,
            name: '',
            saved_at: new Date().toISOString(),
            mode: (typeof currentMode !== 'undefined' ? currentMode : 'txt2img'),
            model: (document.getElementById('ig-model') || {}).value || '',
            settings: collectFields(),
            loras: Array.isArray(window.activeLoRAs) ? window.activeLoRAs : [],
            pose: (window.igPoseJoints && window.igPoseJoints.length)
                ? { joints: window.igPoseJoints, size: window.igPoseSize || 768 }
                : null
        };
    }

    function setStatus(text, kind) {
        const el = document.getElementById('ig-workflow-status');
        if (!el) return;
        el.textContent = text || '';
        el.dataset.kind = kind || '';
        if (!text) return;
        clearTimeout(setStatus._t);
        setStatus._t = setTimeout(() => { el.textContent = ''; el.dataset.kind = ''; }, 6000);
    }

    function go() { return window.go && window.go.main && window.go.main.App; }

    let workflows = [];

    async function refresh() {
        const sel = document.getElementById('ig-workflow-select');
        if (!sel) return;
        const previous = sel.value;
        sel.innerHTML = '<option value="">-- Saved workflows --</option>';
        const api = go();
        if (!api || typeof api.ListImageWorkflows !== 'function') {
            setStatus('Workflow storage is unavailable in this build.', 'error');
            return;
        }
        try {
            workflows = (await api.ListImageWorkflows()) || [];
        } catch (err) {
            setStatus('Could not read workflows: ' + (err && err.message ? err.message : err), 'error');
            return;
        }
        workflows.forEach(w => {
            const opt = document.createElement('option');
            opt.value = w.file_name;
            const bits = [w.mode, w.model, w.saved_at ? w.saved_at.slice(0, 10) : ''].filter(Boolean);
            opt.textContent = w.name + (bits.length ? ' — ' + bits.join(' · ') : '');
            opt.dataset.name = w.name;
            sel.appendChild(opt);
        });
        if (previous && workflows.some(w => w.file_name === previous)) sel.value = previous;
        const count = document.getElementById('ig-workflow-count');
        if (count) count.textContent = workflows.length ? String(workflows.length) : '';
    }

    async function saveWorkflow() {
        const nameEl = document.getElementById('ig-workflow-name');
        const sel = document.getElementById('ig-workflow-select');
        const api = go();
        if (!api || typeof api.SaveImageWorkflow !== 'function') {
            setStatus('Workflow storage is unavailable in this build.', 'error');
            return;
        }
        let name = (nameEl && nameEl.value || '').trim();
        if (!name && sel && sel.value) {
            const current = workflows.find(w => w.file_name === sel.value);
            name = current ? current.name : '';
        }
        if (!name) {
            setStatus('Give the workflow a name first.', 'error');
            if (nameEl) nameEl.focus();
            return;
        }
        const doc = captureSettings();
        doc.name = name;
        try {
            await api.SaveImageWorkflow(JSON.stringify(doc));
            if (nameEl) nameEl.value = name;
            await refresh();
            if (sel) {
                const match = workflows.find(w => w.name === name);
                if (match) sel.value = match.file_name;
            }
            setStatus('Saved “' + name + '”', 'ok');
        } catch (err) {
            setStatus('Save failed: ' + (err && err.message ? err.message : err), 'error');
        }
    }

    async function loadWorkflow() {
        const sel = document.getElementById('ig-workflow-select');
        const api = go();
        if (!sel || !sel.value || !api) return;
        let doc;
        try {
            const raw = await api.GetImageWorkflow(sel.value);
            doc = JSON.parse(raw);
        } catch (err) {
            setStatus('Could not open that workflow: ' + (err && err.message ? err.message : err), 'error');
            return;
        }

        // Mode first: it scopes which panels and model families are valid.
        if (doc.mode && typeof setWorkflowMode === 'function') {
            try { setWorkflowMode(doc.mode); } catch (_) { /* mode may be unknown in older files */ }
        }
        const fields = doc.settings || {};
        let restored = 0;
        let skipped = 0;
        Object.keys(fields).forEach(id => {
            if (applyField(id, fields[id])) restored++; else skipped++;
        });

        if (Array.isArray(doc.loras) && doc.loras.length) {
            window.activeLoRAs = doc.loras;
            if (typeof window.renderActiveLoRAs === 'function') window.renderActiveLoRAs();
        }

        if (doc.pose && Array.isArray(doc.pose.joints) && doc.pose.joints.length && typeof window.igApplyPose === 'function') {
            window.igApplyPose(doc.pose.joints, doc.pose.size || 768);
        }

        const nameEl = document.getElementById('ig-workflow-name');
        if (nameEl) nameEl.value = doc.name || '';
        setStatus('Loaded “' + (doc.name || sel.value) + '” · ' + restored + ' settings' +
            (skipped ? ' · ' + skipped + ' left untouched' : ''), 'ok');
    }

    async function deleteWorkflow() {
        const sel = document.getElementById('ig-workflow-select');
        const api = go();
        if (!sel || !sel.value || !api) return;
        const current = workflows.find(w => w.file_name === sel.value);
        const label = current ? current.name : sel.value;
        if (!confirm('Delete the workflow “' + label + '”?\n\nThe image files it produced are not affected.')) return;
        try {
            await api.DeleteImageWorkflow(sel.value);
            sel.value = '';
            await refresh();
            const nameEl = document.getElementById('ig-workflow-name');
            if (nameEl) nameEl.value = '';
            setStatus('Deleted “' + label + '”', 'ok');
        } catch (err) {
            setStatus('Delete failed: ' + (err && err.message ? err.message : err), 'error');
        }
    }

    function openFolder() {
        const api = go();
        if (!api || typeof api.OpenImageWorkflowsFolder !== 'function') return;
        api.OpenImageWorkflowsFolder().catch(err => setStatus('Could not open folder: ' + (err && err.message ? err.message : err), 'error'));
    }

    async function copyFromFile(file) {
        // Import by reading the JSON, then saving it through the backend so the
        // same validation (format, version, name) applies to shared recipes.
        const text = await file.text();
        let doc;
        try { doc = JSON.parse(text); } catch (_) {
            setStatus('That file is not a workflow (invalid JSON).', 'error');
            return;
        }
        if (doc.format && doc.format !== FORMAT) {
            setStatus('That is a ' + doc.format + ' file, not a HOWL workflow.', 'error');
            return;
        }
        const api = go();
        if (!api) return;
        try {
            await api.SaveImageWorkflow(JSON.stringify(doc));
            await refresh();
            const nameEl = document.getElementById('ig-workflow-name');
            if (nameEl) nameEl.value = doc.name || '';
            setStatus('Imported “' + (doc.name || file.name) + '”', 'ok');
        } catch (err) {
            setStatus('Import failed: ' + (err && err.message ? err.message : err), 'error');
        }
    }

    function exportCurrent() {
        const doc = captureSettings();
        const nameEl = document.getElementById('ig-workflow-name');
        doc.name = (nameEl && nameEl.value.trim()) || 'workflow';
        const blob = new Blob([JSON.stringify(doc, null, 2)], { type: 'application/json' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = doc.name.replace(/[^A-Za-z0-9 _.-]+/g, '-') + '.howl.json';
        document.body.appendChild(a);
        a.click();
        a.remove();
        setTimeout(() => URL.revokeObjectURL(url), 2000);
        setStatus('Exported “' + doc.name + '”', 'ok');
    }

    // -----------------------------------------------------------------------
    // Sequential pipeline runner
    // -----------------------------------------------------------------------

    function pipelineStatus(text, kind) {
        const el = document.getElementById('ig-wf-pipeline-status');
        if (!el) return;
        el.textContent = text || '';
        el.dataset.kind = kind || '';
    }

    function pipelineValue(id) {
        const el = document.getElementById(id);
        return el ? el.value : '';
    }

    function pipelineNumber(id, fallback) {
        const n = Number(pipelineValue(id));
        return Number.isFinite(n) ? n : fallback;
    }

    function renderPipelineStages(stages) {
        const box = document.getElementById('ig-wf-pipeline-stages');
        if (!box) return;
        box.innerHTML = '';
        (stages || []).forEach(stage => {
            const row = document.createElement('div');
            row.className = 'ig-wf-stage ig-wf-stage-' + (stage.status || 'unknown');
            const title = document.createElement('div');
            title.className = 'ig-wf-stage-title';
            title.textContent = `${stage.index}. ${stage.name} — ${stage.status}`;
            row.appendChild(title);
            const detail = document.createElement('div');
            detail.className = 'ig-wf-stage-detail';
            const bits = [];
            if (stage.executor) bits.push(stage.executor);
            if (stage.duration_ms != null) bits.push(`${stage.duration_ms} ms`);
            if (stage.output_width && stage.output_height) bits.push(`${stage.output_width}×${stage.output_height}`);
            if (stage.job_id) bits.push(stage.job_id);
            detail.textContent = bits.join(' · ');
            row.appendChild(detail);
            if (stage.admission) {
                const adm = document.createElement('div');
                adm.className = 'ig-wf-stage-detail';
                adm.textContent = stage.admission;
                row.appendChild(adm);
            }
            if (stage.placement) {
                const placement = document.createElement('div');
                placement.className = 'ig-wf-stage-detail';
                placement.textContent = 'placement: ' + (stage.placement.summary || 'see placement report');
                row.appendChild(placement);
            }
            if (stage.error) {
                const err = document.createElement('div');
                err.className = 'ig-wf-stage-error';
                err.textContent = stage.error;
                row.appendChild(err);
            }
            box.appendChild(row);
        });
    }

    async function runPipeline() {
        const api = go();
        if (!api || typeof api.RunImageWorkflow !== 'function') {
            pipelineStatus('Sequential pipeline is unavailable in this build.', 'error');
            return;
        }
        const button = document.getElementById('ig-wf-run');
        const cancelButton = document.getElementById('ig-wf-cancel');
        const upscaler = pipelineValue('ig-wf-upscaler').trim();
        const controlnet = pipelineValue('ig-wf-controlnet').trim();
        let pose = pipelineValue('ig-wf-pose').trim();
        if (!pose && window.controlNetImageData) {
            try {
                pose = await api.SaveTempImage(window.controlNetImageData);
                const poseEl = document.getElementById('ig-wf-pose');
                if (poseEl) poseEl.value = pose;
            } catch (err) {
                pipelineStatus('Could not stage the current pose image: ' + (err && err.message ? err.message : err), 'error');
                return;
            }
        }
        const ipAdapter = pipelineValue('ig-ip-model').trim();
        const clipVision = pipelineValue('ig-ip-vision').trim();
        let ipReference = '';
        if (ipAdapter || clipVision) {
            if (!ipAdapter || !clipVision || !window.ipAdapterImageData) {
                pipelineStatus('IP-Adapter needs an adapter, a CLIP-Vision encoder and a reference image.', 'error');
                return;
            }
            try {
                ipReference = await api.SaveTempImage(window.ipAdapterImageData);
            } catch (err) {
                pipelineStatus('Could not stage the IP-Adapter reference: ' + (err && err.message ? err.message : err), 'error');
                return;
            }
        }
        const opts = {
            model_name: pipelineValue('ig-model'),
            prompt: pipelineValue('ig-prompt'),
            negative_prompt: pipelineValue('ig-negative'),
            seed: Math.trunc(pipelineNumber('ig-seed', -1)),
            sampler: pipelineValue('ig-sampler'),
            scheduler: pipelineValue('ig-scheduler'),
            cfg_scale: pipelineNumber('ig-cfg', 7),
            base_width: Math.trunc(pipelineNumber('ig-width', 512)),
            base_height: Math.trunc(pipelineNumber('ig-height', 512)),
            base_steps: Math.trunc(pipelineNumber('ig-steps', 20)),
            upscaler_path: upscaler,
            upscaler_tile: 128,
            control_net_path: controlnet,
            pose_image_path: pose,
            control_strength: 0.9,
            refine_width: Math.trunc(pipelineNumber('ig-wf-refine-width', 768)),
            refine_height: Math.trunc(pipelineNumber('ig-wf-refine-height', 768)),
            refine_steps: Math.trunc(pipelineNumber('ig-wf-refine-steps', 20)),
            refine_strength: pipelineNumber('ig-wf-refine-strength', 0.5),
            batch_count: Math.max(1, Math.trunc(pipelineNumber('ig-batch', 1))),
            // The workflow uses the same Performance controls as a single
            // generation. Without this the panel could say 5 GiB while the
            // resident engine kept its previous budget.
            max_vram: (typeof readImageMaxVRAM === 'function') ? readImageMaxVRAM() : undefined,
            ip_adapter_path: ipAdapter,
            clip_vision_path: clipVision,
            ip_adapter_image_path: ipReference,
            ip_adapter_strength: pipelineNumber('ig-ip-strength', 0.7),
            clip_vision_cpu: (document.getElementById('ig-ip-vision-cpu') || {}).checked || false,
            engine_profile: ipAdapter ? 'sdcpp-c6dfa9d-ip' : ''
        };
        if (!opts.prompt.trim()) { pipelineStatus('Enter a prompt first.', 'error'); return; }
        if (!upscaler) { pipelineStatus('Enter the ESRGAN upscaler path.', 'error'); return; }
        if (!controlnet) { pipelineStatus('Enter the OpenPose ControlNet path.', 'error'); return; }
        if (!pose) { pipelineStatus('Add a pose image, or open the Pose Editor first. Pose extraction is external (OpenPose/DWPose).', 'error'); return; }

        button.disabled = true;
        if (cancelButton) cancelButton.disabled = false;
        pipelineStatus('Running base generation…', '');
        renderPipelineStages([]);
        try {
            const result = await api.RunImageWorkflow(opts);
            renderPipelineStages(result && result.stages);
            const failed = result && result.error;
            const cancelled = !!failed && failed.indexOf('image workflow cancelled') !== -1 && failed.indexOf('engine termination failed') === -1;
            const passes = (result && result.batch_runs && result.batch_runs.length) || 1;
            const done = passes > 1 ? `Pipeline complete (${passes} passes). Weights retained unless reclamation reported otherwise.` : 'Pipeline complete. Weights retained unless reclamation reported otherwise.';
            const status = cancelled ? 'Pipeline cancelled; the active engine process reached its termination boundary.' : (failed ? 'Pipeline failed: ' + failed : done);
            pipelineStatus(status, cancelled ? 'ok' : (failed ? 'error' : 'ok'));
            if (result && result.final_image_b64) {
                const fake = {
                    image_data: result.final_image_b64,
                    metadata: {
                        model: opts.model_name,
                        family: '',
                        width: result.width,
                        height: result.height,
                        steps: opts.refine_steps,
                        seed: opts.seed,
                        sampler: opts.sampler,
                        cfg_scale: opts.cfg_scale,
                        mode: 'workflow',
                        file: result.final_image_path || ''
                    }
                };
                if (typeof displayResult === 'function') displayResult(fake);
                try { currentResult = fake; } catch (e) { /* top-level binding is optional */ }
            }
        } catch (err) {
            pipelineStatus('Pipeline failed: ' + (err && err.message ? err.message : err), 'error');
        } finally {
            button.disabled = false;
            if (cancelButton) cancelButton.disabled = true;
        }
    }

    async function cancelPipeline() {
        const api = go();
        const button = document.getElementById('ig-wf-cancel');
        if (!api || typeof api.CancelImageWorkflow !== 'function') {
            pipelineStatus('Pipeline cancellation is unavailable in this build.', 'error');
            return;
        }
        if (button) button.disabled = true;
        pipelineStatus('Cancelling; waiting for the active engine operation and its process to terminate…', '');
        try {
            await api.CancelImageWorkflow();
        } catch (err) {
            pipelineStatus('Cancel failed: ' + (err && err.message ? err.message : err), 'error');
            if (button) button.disabled = false;
        }
    }

    function useCurrentPose() {
        if (!window.controlNetImageData) {
            pipelineStatus('No current ControlNet image. Draw a pose or load one first.', 'error');
            return;
        }
        const api = go();
        if (!api || typeof api.SaveTempImage !== 'function') {
            pipelineStatus('Temporary image storage is unavailable.', 'error');
            return;
        }
        api.SaveTempImage(window.controlNetImageData).then(path => {
            const el = document.getElementById('ig-wf-pose');
            if (el) el.value = path || '';
            pipelineStatus('Current pose image staged for the pipeline.', 'ok');
        }).catch(err => pipelineStatus('Could not stage pose: ' + (err && err.message ? err.message : err), 'error'));
    }

    function bind() {
        const on = (id, fn) => {
            const el = document.getElementById(id);
            if (el) el.addEventListener('click', fn);
        };
        on('ig-workflow-save', saveWorkflow);
        on('ig-workflow-load', loadWorkflow);
        on('ig-workflow-delete', deleteWorkflow);
        on('ig-workflow-open-folder', openFolder);
        on('ig-workflow-export', exportCurrent);
        on('ig-wf-run', runPipeline);
        on('ig-wf-cancel', cancelPipeline);
        on('ig-wf-use-pose', useCurrentPose);

        const file = document.getElementById('ig-workflow-file');
        if (file) {
            file.addEventListener('change', e => {
                const f = e.target.files && e.target.files[0];
                if (f) copyFromFile(f).finally(() => { e.target.value = ''; });
            });
        }
        const sel = document.getElementById('ig-workflow-select');
        if (sel) sel.addEventListener('dblclick', loadWorkflow);

        refresh();
    }

    window.howlWorkflows = { capture: captureSettings, refresh, save: saveWorkflow, load: loadWorkflow };

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', bind);
    } else {
        bind();
    }
})();
