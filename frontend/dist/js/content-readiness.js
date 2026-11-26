/**
 * Content Readiness — per-library storage state for the renderer
 *
 * The backend is authoritative: GetContentReadiness() reports one state per
 * content library (category + library identity, readable, writable, reason,
 * failed step, retryable). This module never invents a status. It renders what
 * the backend said, keeps the mutation controls in step with it, and refreshes
 * after a refused write.
 *
 * A read-only library is displayed as read-only, with its content visible; it is
 * never called unavailable. A category summary is only a presentation rollup: it
 * carries "N of M libraries need attention" so a healthy library cannot hide a
 * failing sibling, and the per-library lines stay visible.
 *
 * @module content-readiness
 * @version 1.0.0
 */

'use strict';

const ContentReadiness = (() => {
    const PAGE_CATEGORIES = [
        { test: /(^|\/)world\.html$/i, category: 'worlds', label: 'Worlds' },
        { test: /(^|\/)scenario\.html$/i, category: 'scenarios', label: 'Scenarios' },
    ];

    let state = { libraries: [], categories: [] };
    let refreshing = null;

    function _app() {
        return window.go?.main?.App;
    }

    function pageCategory() {
        const path = window.location.pathname;
        for (const page of PAGE_CATEGORIES) {
            if (page.test.test(path)) return page;
        }
        return null;
    }

    function category(cat) {
        return state.categories.find(entry => entry.category === cat) || null;
    }

    function libraries(cat) {
        return (state.libraries || []).filter(entry => entry.category === cat);
    }

    function canRead(cat) {
        const summary = category(cat);
        if (summary) return Boolean(summary.readable);
        return true;
    }

    // Editing an existing item is possible while any readable library is
    // writable. Creating is not: new documents go to the default save root, so
    // create follows that root's own writability. A read-only default root must
    // not make a healthy sibling's content look unusable, and a writable sibling
    // must not make create look possible when it would be refused.
    function canMutate(cat) {
        const summary = category(cat);
        return Boolean(summary && summary.writable);
    }

    function canCreate(cat) {
        const summary = category(cat);
        if (!summary) return false;
        if (typeof summary.saveWritable === 'boolean') return summary.saveWritable;
        return Boolean(summary.writable);
    }

    function reasonFor(cat) {
        const summary = category(cat);
        if (!summary || !summary.reason) return '';
        return summary.reason;
    }

    function _statusLabel(status) {
        switch (status) {
            case 'ready': return 'ready';
            case 'read-only': return 'read-only — content visible, new items refused';
            case 'not-readable': return 'not readable — cached items only';
            case 'not-initialized': return 'not set up yet';
            case 'unavailable': return 'unavailable';
            default: return status || 'unknown';
        }
    }

    function _bannerHost() {
        return document.querySelector('.container') ||
            document.querySelector('main') ||
            document.body;
    }

    function _ensureBanner() {
        let banner = document.getElementById('content-readiness-banner');
        if (banner) return banner;
        const host = _bannerHost();
        if (!host) return null;
        banner = document.createElement('div');
        banner.id = 'content-readiness-banner';
        banner.setAttribute('role', 'status');
        banner.style.cssText = [
            'display:none',
            'margin:0 0 16px 0',
            'padding:12px 16px',
            'border-radius:10px',
            'border:1px solid var(--border-primary, rgba(255,255,255,0.18))',
            'background:var(--shell-panel-bg, rgba(30,30,36,0.92))',
            'color:var(--text-primary, #e8e8ef)',
            'font-size:13px',
            'line-height:1.5',
        ].join(';');
        if (host.firstChild) host.insertBefore(banner, host.firstChild);
        else host.appendChild(banner);
        return banner;
    }

    function _renderBanner(page) {
        const banner = _ensureBanner();
        if (!banner) return;
        const summary = category(page.category);
        const entries = libraries(page.category);
        if (!summary || entries.length === 0 || summary.failingLibraries === 0) {
            banner.style.display = 'none';
            banner.innerHTML = '';
            return;
        }

        const heading = summary.failingLibraries === 1
            ? `${page.label}: 1 of ${summary.libraries} libraries needs attention`
            : `${page.label}: ${summary.failingLibraries} of ${summary.libraries} libraries need attention`;

        const rows = entries.map(entry => {
            const label = _statusLabel(entry.status);
            const detail = entry.reason ? ` ${entry.reason}` : '';
            return `<li style="margin:2px 0"><code style="opacity:0.85">${_escape(entry.path)}</code> — <strong>${_escape(label)}</strong>.${_escape(detail)}</li>`;
        }).join('');

        const needsSetup = entries.some(entry => entry.status === 'not-initialized' || entry.status === 'unavailable');

        banner.innerHTML = `
            <div style="display:flex;gap:12px;align-items:flex-start;justify-content:space-between;flex-wrap:wrap">
                <div>
                    <div style="font-weight:600">${_escape(heading)}</div>
                    <ul style="margin:6px 0 0 18px;padding:0">${rows}</ul>
                </div>
                ${needsSetup ? `<button type="button" data-content-readiness="setup" style="padding:8px 12px;border-radius:8px;border:1px solid var(--border-primary, rgba(255,255,255,0.25));background:transparent;color:inherit;cursor:pointer">Set up default location</button>` : ''}
            </div>`;
        banner.style.display = 'block';

        const setup = banner.querySelector('[data-content-readiness="setup"]');
        if (setup) {
            setup.addEventListener('click', async () => {
                const app = _app();
                if (typeof app?.EnsureDefaultContentLibrary !== 'function') return;
                setup.disabled = true;
                setup.textContent = 'Setting up…';
                try {
                    await app.EnsureDefaultContentLibrary(page.category);
                } catch (error) {
                    console.warn('[contentReadiness] setup failed:', error);
                }
                setup.disabled = false;
                setup.textContent = 'Set up default location';
                await refresh();
                window.dispatchEvent(new CustomEvent('howl:content-readiness-changed', { detail: { category: page.category } }));
            });
        }
    }

    function _escape(value) {
        return String(value == null ? '' : value)
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;');
    }

    function _applyControls(page) {
        const reason = reasonFor(page.category);
        const rules = [
            { selector: '.create-world-button', enabled: canCreate(page.category) },
            { selector: '.world-creation-save', enabled: canMutate(page.category) },
        ];
        for (const rule of rules) {
            document.querySelectorAll(rule.selector).forEach(element => {
                element.disabled = !rule.enabled;
                element.style.opacity = rule.enabled ? '' : '0.55';
                element.style.cursor = rule.enabled ? '' : 'not-allowed';
                if (rule.enabled) {
                    element.removeAttribute('title');
                } else {
                    element.title = reason || `This ${page.label.toLowerCase()} library cannot be written.`;
                }
            });
        }
    }

    function render() {
        const page = pageCategory();
        if (!page) return;
        _renderBanner(page);
        _applyControls(page);
    }

    async function refresh() {
        if (refreshing) return refreshing;
        const app = _app();
        if (typeof app?.GetContentReadiness !== 'function') return state;
        refreshing = (async () => {
            try {
                const next = await app.GetContentReadiness();
                if (next && Array.isArray(next.libraries)) state = next;
            } catch (error) {
                console.warn('[contentReadiness] readiness unavailable:', error);
            } finally {
                refreshing = null;
            }
            render();
            window.dispatchEvent(new CustomEvent('howl:content-readiness', { detail: state }));
            return state;
        })();
        return refreshing;
    }

    document.addEventListener('DOMContentLoaded', () => {
        refresh();
    });

    return {
        get state() { return state; },
        refresh,
        render,
        pageCategory,
        category,
        libraries,
        canRead,
        canMutate,
        canCreate,
        reasonFor,
    };
})();

window.ContentReadiness = ContentReadiness;
