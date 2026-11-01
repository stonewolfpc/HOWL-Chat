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
