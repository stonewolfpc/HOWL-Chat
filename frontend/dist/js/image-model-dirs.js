/**
 * image-model-dirs.js
 *
 * Image-model directories, handled the same way as GGUF loading:
 *   - pick a folder with the native dialog
 *   - the folder is remembered and rescanned on later launches
 *   - several folders can be kept side by side
 *
 * The backend already owns all of this (AddModelDirectory / RemoveModelDirectory /
 * GetModelDirectories / RefreshModelDirectories with kind "image"); this module is
 * the UI for it on the Image Generator page, plus the filter that keeps
 * non-Stable-Diffusion files out of the model dropdown.
 *
 * Depends on: OpenImageModelsDirPicker, AddModelDirectory, RemoveModelDirectory,
 *             GetModelDirectories, RefreshModelDirectories, loadModels
 */
(function () {
    'use strict';

    const KIND = 'image';

    function esc(s) {
        return String(s == null ? '' : s)
            .replace(/&/g, '&amp;').replace(/</g, '&lt;')
            .replace(/>/g, '&gt;').replace(/"/g, '&quot;');
    }

    function app() {
        return window.go && window.go.main && window.go.main.App;
    }

    function shortPath(p) {
        const parts = String(p).split(/[\\/]/).filter(Boolean);
        return parts.length <= 2 ? p : '…' + '/' + parts.slice(-2).join('/');
    }

    async function refreshImageModelDirs() {
        const list = document.getElementById('ig-model-dirs-list');
        const countEl = document.getElementById('ig-model-dirs-count');
        if (!list) return;
        const api = app();
        if (!api || !api.GetModelDirectories) {
            list.innerHTML = '<div class="ig-dir-empty">Directory management unavailable.</div>';
            return;
        }
        let dirs = [];
        try {
            const all = await api.GetModelDirectories();
            dirs = (all || []).filter(d => (d.kind || '').toLowerCase() === KIND);
        } catch (err) {
            list.innerHTML = '<div class="ig-dir-empty">Could not read directories: ' + esc(err) + '</div>';
            return;
        }

        if (countEl) {
            countEl.textContent = dirs.length === 0
                ? 'no folders added'
                : dirs.length + ' folder' + (dirs.length === 1 ? '' : 's');
        }

        if (dirs.length === 0) {
            list.innerHTML = '<div class="ig-dir-empty">No image folders yet — add one to load SD models from anywhere.</div>';
            return;
        }

        list.innerHTML = '';
        dirs.forEach(d => {
            const row = document.createElement('div');
            row.className = 'ig-dir-row' + (d.exists === false ? ' ig-dir-missing' : '');
            row.innerHTML =
                '<span class="ig-dir-path" title="' + esc(d.path) + '">' + esc(shortPath(d.path)) + '</span>' +
                '<span class="ig-dir-meta">' + (d.exists === false ? 'missing' : (d.models || 0) + ' assets') + '</span>' +
                '<button class="ig-dir-remove" title="Forget this folder (files are not deleted)">&#x2715;</button>';
            row.querySelector('.ig-dir-remove').addEventListener('click', (e) => {
                e.stopPropagation();
                removeImageModelDir(d.path);
            });
            list.appendChild(row);
        });
    }

    async function addImageModelDir() {
        const api = app();
        if (!api) return;
        let dir = '';
        try {
            dir = await api.OpenImageModelsDirPicker();
        } catch (err) {
            if (typeof setStatus === 'function') setStatus('Folder picker failed: ' + err);
            return;
        }
        if (!dir) return; // cancelled

        try {
            if (typeof setStatus === 'function') setStatus('Scanning ' + shortPath(dir) + '…');
            const added = await api.AddModelDirectory(dir, KIND);
            await refreshImageModelDirs();
            if (typeof loadModels === 'function') await loadModels();
            if (typeof setStatus === 'function') {
                setStatus('Added ' + shortPath(dir) + ' — ' + (added && added.models != null ? added.models : 0) + ' assets');
            }
        } catch (err) {
            if (typeof setStatus === 'function') setStatus('Could not add folder: ' + err);
        }
    }

    async function removeImageModelDir(dir) {
        const api = app();
        if (!api || !dir) return;
        // WebView2 does not show window.confirm, so the raw call returned false
        // and the folder-row X did nothing at all in the packaged app. The shared
        // DeleteConfirm modal is the in-app confirmation every other destructive
        // action on this frontend already uses.
        const forget = async () => {
            try {
                await api.RemoveModelDirectory(dir);
                await refreshImageModelDirs();
                if (typeof loadModels === 'function') await loadModels();
                if (typeof setStatus === 'function') setStatus('Forgot folder ' + shortPath(dir));
            } catch (err) {
                if (typeof setStatus === 'function') setStatus('Could not remove folder: ' + err);
            }
        };
        if (window.DeleteConfirm && window.DeleteConfirm.show) {
            window.DeleteConfirm.show({
                title: 'Forget this folder?',
                body: dir,
                warning: 'Only the registration is removed. Model files on disk stay where they are.',
                confirmLabel: 'Forget',
                onConfirm: forget,
            });
            return;
        }
        const ok = typeof confirm === 'function'
            ? confirm('Forget this folder? ' + dir + ' � model files on disk are not deleted.')
            : true;
        if (!ok) return;
        await forget();
    }

    async function rescanImageModelDirs() {
        const api = app();
        if (!api) return;
        try {
            if (typeof setStatus === 'function') setStatus('Rescanning image folders…');
            if (api.RefreshModelDirectories) await api.RefreshModelDirectories(KIND);
            await refreshImageModelDirs();
            if (typeof loadModels === 'function') await loadModels();
            if (typeof setStatus === 'function') setStatus('Image folders rescanned');
        } catch (err) {
            if (typeof setStatus === 'function') setStatus('Rescan failed: ' + err);
        }
    }

    /**
     * A file is usable as an image model only when stable-diffusion.cpp can load
     * it: a known diffusion family and a format the current binary accepts.
     * Everything else (LoRAs, VAEs, ControlNets, embeddings, workflows, stray
     * GGUF text models) is classified elsewhere and must not appear here.
     */
    // Video families are real models, but they cannot render a still image.
    // Leaving them in the image dropdown is the exact "looks available, then
    // fails" problem: they belong in Video mode only, and SD checkpoints do not
    // belong in Video mode.
    const VIDEO_ONLY_FAMILIES = ['wan', 'wan2.1', 'wan2.2', 'svd', 'hunyuan-video', 'ltx-video', 'ltxv'];

    function isVideoFamily(family) {
        return VIDEO_ONLY_FAMILIES.some(v => family === v || family.startsWith(v));
    }

    function isUsableImageModel(m, mode) {
        if (!m) return false;
        const family = String(m.family || '').toLowerCase();
        const compat = String(m.compat || m.compatibility || '').toLowerCase();
        if (!family || family === 'unknown') return false;
        if (compat === 'unsupported' || compat === 'incompatible') return false;
        const wantsVideo = String(mode) === 'video';
        if (isVideoFamily(family)) return wantsVideo;
        return !wantsVideo;
    }

    function splitModels(list, mode) {
        const all = Array.isArray(list) ? list : [];
        const usable = all.filter(m => isUsableImageModel(m, mode));
        return { usable: usable, hidden: all.length - usable.length, total: all.length };
    }

    window.refreshImageModelDirs = refreshImageModelDirs;
    window.addImageModelDir = addImageModelDir;
    window.removeImageModelDir = removeImageModelDir;
    window.rescanImageModelDirs = rescanImageModelDirs;
    window.splitCompatibleImageModels = splitModels;

    document.addEventListener('DOMContentLoaded', () => setTimeout(refreshImageModelDirs, 400));
})();
