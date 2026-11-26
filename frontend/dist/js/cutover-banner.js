/**
 * Explicit domain cutover banner.
 *
 * A page never imports legacy data on its own. This banner is the one visible
 * user action that starts or reverses the ownership transition for the domain
 * the page edits.
 *
 * @module cutover-banner
 */

(function () {
    'use strict';

    const TARGETS = [
        { path: 'world.html', label: 'Worlds', repo: 'WorldStorage' },
        { path: 'scenario.html', label: 'Scenarios', repo: 'ScenarioStorage' },
        { path: 'character.html', label: 'Characters', repo: 'CharacterStorage' },
        { path: 'image-generator.html', label: 'Image generation settings', repo: 'ImageGenerationStorage' },
    ];

    function targetForPage() {
        const path = window.location.pathname.replace(/\\/g, '/');
        return TARGETS.find(target => path.endsWith('/' + target.path) || path.endsWith(target.path));
    }

    function styles() {
        const style = document.createElement('style');
        style.textContent = `
            .cutover-banner{position:fixed;left:18px;bottom:18px;z-index:9999;max-width:430px;padding:14px 16px;border:1px solid rgba(154,120,70,.55);border-radius:10px;background:#1d1a16;color:#e8dfcf;box-shadow:0 8px 26px rgba(0,0,0,.45);font:13px/1.45 system-ui,sans-serif}
            .cutover-banner strong{display:block;margin-bottom:4px;color:#f1e6d2;font-size:14px}
            .cutover-banner p{margin:0 0 10px}
            .cutover-banner .cutover-actions{display:flex;gap:8px;flex-wrap:wrap}
            .cutover-banner button{cursor:pointer;border-radius:6px;border:1px solid #8a6a3d;padding:6px 10px;background:#2b241b;color:#f1e6d2;font:inherit}
            .cutover-banner button.primary{background:#9a783f;border-color:#b08c50;color:#1a1611;font-weight:600}
            .cutover-banner button:disabled{opacity:.55;cursor:wait}
        `;
        document.head.appendChild(style);
    }

    function render(repo, target) {
        const status = repo.cutoverStatus || { cutover: Boolean(repo.cutover) };
        const existing = document.getElementById('cutover-banner');
        if (existing) existing.remove();
        if (!status.cutover && sessionStorage.getItem('cutover-dismissed:' + target.path)) return;

        const banner = document.createElement('div');
        banner.id = 'cutover-banner';
        banner.className = 'cutover-banner';
        const title = document.createElement('strong');
        const text = document.createElement('p');
        const actions = document.createElement('div');
        actions.className = 'cutover-actions';

        if (status.cutover) {
            title.textContent = target.label + ' are library-owned';
            const at = status.cutover_at ? ' since ' + status.cutover_at : '';
            text.textContent = 'The canonical library owns this domain' + at + '. The browser copy is no longer authoritative.';
            if (status.rollback_available) {
                const undo = document.createElement('button');
                undo.textContent = 'Undo last import';
                undo.addEventListener('click', async () => {
                    if (!confirm('Undo the last import and return ' + target.label.toLowerCase() + ' to the browser copy?')) return;
                    await run(banner, () => repo.rollback(), 'Import undone. Reloading…');
                });
                actions.appendChild(undo);
            }
        } else {
            title.textContent = target.label + ' are still browser-owned';
            text.textContent = 'Import this browser copy into the library when you are ready. Nothing is moved automatically.';
            const importButton = document.createElement('button');
            importButton.className = 'primary';
            importButton.textContent = 'Import to library';
            importButton.addEventListener('click', async () => {
                if (!confirm('Import this browser copy of ' + target.label.toLowerCase() + ' into the library now?')) return;
                await run(banner, () => repo.importLegacy(), 'Import complete. Reloading…');
            });
            actions.appendChild(importButton);
        }

        const dismiss = document.createElement('button');
        dismiss.textContent = 'Not now';
        dismiss.addEventListener('click', () => {
            sessionStorage.setItem('cutover-dismissed:' + target.path, '1');
            banner.remove();
        });
        actions.appendChild(dismiss);
        banner.append(title, text, actions);
        document.body.appendChild(banner);
    }

    async function run(banner, action, success) {
        banner.querySelectorAll('button').forEach(button => { button.disabled = true; });
        try {
            await action();
            alert(success);
            window.location.reload();
        } catch (error) {
            banner.querySelectorAll('button').forEach(button => { button.disabled = false; });
            alert((error && error.message) || String(error));
        }
    }

    async function boot() {
        const target = targetForPage();
        const repo = target && window[target.repo];
        if (!repo || typeof repo.ready?.then !== 'function') return;
        try {
            await repo.ready;
        } catch (error) {
            console.error('Cutover status could not be read:', error);
            return;
        }
        styles();
        render(repo, target);
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', boot, { once: true });
    } else {
        boot();
    }
})();
