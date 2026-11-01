/**
 * Image Recall Menu
 *
 * Right-clicking the staged image (or a history thumbnail) opens the parameters
 * that produced it. "Recall all" restores every one; each row restores only that
 * field. Values come from the image's own metadata, never from whatever happens
 * to be on screen, so a recall cannot silently mix two images together.
 *
 * Gallery metadata keeps the interesting half of the recipe in `options`, so
 * every lookup reads through to it before giving up.
 */
(function () {
    'use strict';

    if (window.__howlRecallInstalled) return;
    window.__howlRecallInstalled = true;

    const STYLE_ID = 'ig-recall-styles';
    let menu = null;

    function ensureStyles() {
        if (document.getElementById(STYLE_ID)) return;
        const style = document.createElement('style');
        style.id = STYLE_ID;
        style.textContent = [
            '.ig-recall-menu{position:fixed;z-index:99999;min-width:280px;max-height:70vh;overflow-y:auto;',
            'background:#12121a;border:1px solid rgba(20,184,166,.35);border-radius:10px;padding:6px;',
            'box-shadow:0 12px 32px rgba(0,0,0,.55);font-size:12px;color:#d4d4d4}',
            '.ig-recall-menu .ig-recall-head{padding:6px 10px 8px;color:var(--text-tertiary);',
            'font-size:11px;text-transform:uppercase;letter-spacing:.6px;border-bottom:1px solid rgba(100,100,100,.25);',
            'margin-bottom:4px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}',
            '.ig-recall-menu button{display:block;width:100%;text-align:left;background:transparent;border:0;',
            'color:inherit;padding:7px 10px;border-radius:6px;cursor:pointer;font-size:12px;font-family:inherit}',
            '.ig-recall-menu button:hover{background:rgba(20,184,166,.15);color:#14b8a6}',
            '.ig-recall-menu hr{border:0;border-top:1px solid rgba(100,100,100,.25);margin:5px 4px}',
            '.ig-recall-menu .ig-recall-all{font-weight:700;color:#14b8a6}',
            '.ig-recall-menu .ig-recall-value{margin-left:6px;color:var(--text-tertiary);font-weight:400}'
        ].join('');
        document.head.appendChild(style);
    }

    function closeMenu() {
        if (menu) {
            menu.remove();
            menu = null;
            document.removeEventListener('click', closeMenu, true);
            document.removeEventListener('keydown', onKey, true);
        }
    }

    function onKey(e) {
        if (e.key === 'Escape') closeMenu();
    }

    function setInput(id, value) {
        if (value === undefined || value === null || value === '') return false;
        const el = document.getElementById(id);
        if (!el) return false;
        el.value = String(value);
        el.dispatchEvent(new Event('input', { bubbles: true }));
        el.dispatchEvent(new Event('change', { bubbles: true }));
        return true;
    }

    function setSelect(id, value) {
        if (value === undefined || value === null || value === '') return false;
        const el = document.getElementById(id);
        if (!el || !el.options) return false;
        const want = String(value).trim().toLowerCase();
        const opt = [...el.options].find(o =>
            (o.value || '').trim().toLowerCase() === want ||
            (o.textContent || '').trim().toLowerCase() === want);
        if (!opt) return false;
        el.value = opt.value;
        el.dispatchEvent(new Event('change', { bubbles: true }));
        return true;
    }

    function capitalize(key) {
        return key.replace(/(^|_)([a-z])/g, (_, __, c) => c.toUpperCase());
    }

    // Top level first, then the nested options block, in snake_case or PascalCase.
    function pick(meta, key) {
        if (!meta) return null;
        const sources = [meta, meta.options, meta.Options];
        const names = [key, capitalize(key)];
        for (const src of sources) {
            if (!src || typeof src !== 'object') continue;
            for (const name of names) {
                const v = src[name];
                if (v !== undefined && v !== null && v !== '') return v;
            }
        }
        return null;
    }

    const FIELDS = [
        { label: 'Prompt',          key: 'prompt',          apply: (m) => setInput('ig-prompt', pick(m, 'prompt')) },
        { label: 'Negative Prompt', key: 'negative_prompt', apply: (m) => setInput('ig-negative', pick(m, 'negative_prompt')) },
        { label: 'Seed',            key: 'seed',            apply: (m) => setInput('ig-seed', pick(m, 'seed')) },
        { label: 'Steps',           key: 'steps',           apply: (m) => setInput('ig-steps', pick(m, 'steps')) },
        { label: 'CFG Scale',       key: 'cfg_scale',       apply: (m) => setInput('ig-cfg', pick(m, 'cfg_scale')) },
        { label: 'Width',           key: 'width',           apply: (m) => setInput('ig-width', pick(m, 'width')) },
        { label: 'Height',          key: 'height',          apply: (m) => setInput('ig-height', pick(m, 'height')) },
        { label: 'Sampler',         key: 'sampler',         apply: (m) => setSelect('ig-sampler', pick(m, 'sampler')) },
        { label: 'Scheduler',       key: 'scheduler',       apply: (m) => setSelect('ig-scheduler', pick(m, 'scheduler')) },
        { label: 'Model',           key: 'model',           apply: (m) => setSelect('ig-model', pick(m, 'model')) }
    ];

    function displayValue(value) {
        const text = String(value);
        return text.length > 42 ? text.slice(0, 42) + '…' : text;
    }

    function recallAll(meta) {
        const applied = FIELDS.filter(f => f.apply(meta)).map(f => f.label);
        closeMenu();
        if (typeof setStatus === 'function') {
            setStatus(applied.length ? 'Recalled: ' + applied.join(', ') : 'Nothing to recall from this image');
        }
    }

    function recallOne(meta, field) {
        const ok = field.apply(meta);
        closeMenu();
        if (typeof setStatus === 'function') {
            setStatus(ok ? 'Recalled ' + field.label : field.label + ' is not stored on this image');
        }
    }

    function openMenu(meta, x, y, title) {
        ensureStyles();
        closeMenu();
        if (!meta || typeof meta !== 'object') return;

        menu = document.createElement('div');
        menu.className = 'ig-recall-menu';

        const head = document.createElement('div');
        head.className = 'ig-recall-head';
        head.textContent = title || meta.file || 'Image';
        menu.appendChild(head);

        const all = document.createElement('button');
        all.className = 'ig-recall-all';
        all.textContent = 'Recall all settings';
        all.addEventListener('click', (e) => { e.stopPropagation(); recallAll(meta); });
        menu.appendChild(all);
        menu.appendChild(document.createElement('hr'));

        FIELDS.forEach(field => {
            const value = pick(meta, field.key);
            if (value === null) return;
            const btn = document.createElement('button');
            btn.textContent = field.label;
            const span = document.createElement('span');
            span.className = 'ig-recall-value';
            span.textContent = displayValue(value);
            btn.appendChild(span);
            btn.addEventListener('click', (e) => { e.stopPropagation(); recallOne(meta, field); });
            menu.appendChild(btn);
        });

        if (typeof showBoardContextMenu === 'function' && meta.file) {
            menu.appendChild(document.createElement('hr'));
            const board = document.createElement('button');
            board.textContent = 'Send to Board…';
            board.addEventListener('click', (e) => {
                e.stopPropagation();
                closeMenu();
                showBoardContextMenu(meta.file, x, y);
            });
            menu.appendChild(board);
        }

        document.body.appendChild(menu);

        // Stay on screen and scroll rather than overflow.
        const rect = menu.getBoundingClientRect();
        menu.style.left = Math.max(8, Math.min(x, window.innerWidth - rect.width - 8)) + 'px';
        menu.style.top = Math.max(8, Math.min(y, window.innerHeight - rect.height - 8)) + 'px';

        setTimeout(() => {
            document.addEventListener('click', closeMenu, true);
            document.addEventListener('keydown', onKey, true);
        }, 0);
    }

    function metaFromPreview() {
        if (window.currentResult && window.currentResult.metadata) return window.currentResult.metadata;
        const item = document.querySelector('.ig-history-item.active');
        if (item && item.dataset.meta) {
            try { return JSON.parse(item.dataset.meta); } catch (_) { return null; }
        }
        return null;
    }

    document.addEventListener('contextmenu', (event) => {
        const preview = event.target.closest('#ig-preview-area');
        const historyItem = event.target.closest('.ig-history-item');
        if (!preview && !historyItem) return;

        let meta = null;
        let title = null;
        if (historyItem && historyItem.dataset.meta) {
            try { meta = JSON.parse(historyItem.dataset.meta); } catch (_) { meta = null; }
            title = meta && (meta.file || meta.filename);
        }
        if (!meta && preview) meta = metaFromPreview();
        if (!meta) return;

        event.preventDefault();
        event.stopPropagation();
        openMenu(meta, event.clientX, event.clientY, title);
    });

    window.openImageRecallMenu = function () {
        const meta = metaFromPreview();
        if (!meta) return false;
        openMenu(meta, window.innerWidth / 2 - 140, 160);
        return true;
    };
})();
