/**
 * Content Libraries settings UI.
 *
 * Backend failures are surfaced; this file never reports a registration or
 * permission change that the backend did not accept.
 */

'use strict';

(function () {
    const api = () => window.go?.main?.App || null;
    const CATEGORIES = [
        { id: 'chats', label: 'Chats' },
        { id: 'characters', label: 'Characters' },
        { id: 'worlds', label: 'Worlds' },
        { id: 'scenarios', label: 'Scenarios' },
        { id: 'lorebooks', label: 'Lorebooks' },
        { id: 'images', label: 'Generated images' },
        { id: 'videos', label: 'Generated videos' }
    ];
    const ROLES = [
        { value: 1, label: 'Read only' },
        { value: 3, label: 'Read / write' },
        { value: 5, label: 'Read / save' },
        { value: 7, label: 'Read / write / save' }
    ];
    let busy = false;

    function notify(message, kind) {
        if (typeof showNotification === 'function') showNotification(message, kind || 'info');
    }

    function roleLabel(role) {
        return ROLES.find(r => r.value === Number(role))?.label || 'Read only';
    }

    function node(tag, className, text) {
        const el = document.createElement(tag);
        if (className) el.className = className;
        if (text !== undefined) el.textContent = text;
        return el;
    }

    function button(label, className, handler) {
        const el = node('button', className || 'cl-mini', label);
        el.type = 'button';
        el.addEventListener('click', handler);
        return el;
    }

    async function refresh() {
        const root = document.getElementById('content-libraries-manager');
        if (!root) return;
        const app = api();
        if (typeof app?.GetContentLibraries !== 'function') {
            root.replaceChildren(node('div', 'cl-empty', 'Content libraries are unavailable in this build.'));
            return;
        }
        root.replaceChildren(node('div', 'cl-empty', 'Loading libraries...'));
        try {
            const libraries = await app.GetContentLibraries('');
            let inventory = [];
            if (typeof app.GetContentLibraryInventory === 'function') {
                try {
                    inventory = await app.GetContentLibraryInventory('');
                } catch (err) {
                    console.warn('[ContentLibraries] Inventory unavailable:', err);
                }
            }
            render(root, Array.isArray(libraries) ? libraries : [], Array.isArray(inventory) ? inventory : []);
        } catch (err) {
            root.replaceChildren(node('div', 'cl-empty', 'Could not load libraries: ' + (err?.message || err)));
        }
    }

    function render(root, libraries, inventory) {
        root.replaceChildren();
        const inventoryByRoot = new Map(inventory.map(entry => [entry.id, entry]));
        for (const category of CATEGORIES) {
            const group = node('div', 'cl-group');
            const head = node('div', 'cl-group-head');
            head.appendChild(node('span', 'cl-group-title', category.label));
            const actions = node('div', 'cl-group-actions');
            actions.appendChild(button('Add default', 'cl-mini', () => addDefault(category.id)));
            actions.appendChild(button('Add existing folder', 'cl-mini', () => addExisting(category.id)));
            head.appendChild(actions);
            group.appendChild(head);

            const rows = libraries.filter(lib => lib.category === category.id);
            if (rows.length === 0) {
                group.appendChild(node('div', 'cl-note', 'No folders registered. Add the default or point Howl at an existing collection.'));
            } else {
                for (const lib of rows) group.appendChild(renderRow(lib, inventoryByRoot.get(lib.id)));
            }
            root.appendChild(group);
        }
    }

    function renderRow(lib, inventory) {
        const row = node('div', 'cl-row');
        const main = node('div');
        const path = node('div', 'cl-path', lib.path || '');
        path.title = lib.path || '';
        main.appendChild(path);

        const meta = node('div', 'cl-meta');
        meta.appendChild(node('span', 'cl-pill', lib.origin === 'howl' ? 'Howl-managed' : 'External'));
        meta.appendChild(node('span', 'cl-pill', lib.format === 'sillytavern' ? 'SillyTavern' : (lib.format === 'howl' ? 'Howl format' : 'Unknown format')));
        if (lib.default_save) meta.appendChild(node('span', 'cl-pill', 'Save destination'));
        if (!lib.exists) meta.appendChild(node('span', 'cl-pill warn', 'Unavailable'));
        if (inventory?.stale) meta.appendChild(node('span', 'cl-pill warn', 'Last-known snapshot'));
        main.appendChild(meta);
        if (inventory) {
            const count = Array.isArray(inventory.items) ? inventory.items.length : 0;
            let note = '';
            if (inventory.stale) {
                note = `Last-known: ${count} item${count === 1 ? '' : 's'} from the most recent successful scan. This is not a live scan.`;
            } else if (inventory.scanned) {
                note = `${count} item${count === 1 ? '' : 's'} discovered now.`;
            } else if (inventory.scan_error) {
                note = 'Inventory scan failed; the registration is retained. ' + inventory.scan_error;
            } else {
                note = 'No live inventory scanner is available for this category.';
            }
            main.appendChild(node('div', 'cl-note', note));
        }
        row.appendChild(main);

        const controls = node('div', 'cl-controls');
        const select = node('select', 'cl-select');
        select.setAttribute('aria-label', 'Library permission for ' + lib.path);
        for (const role of ROLES) {
            const option = document.createElement('option');
            option.value = String(role.value);
            option.textContent = role.label;
            option.selected = Number(lib.role) === role.value;
            select.appendChild(option);
        }
        select.disabled = busy || !lib.exists;
        select.addEventListener('change', () => changeRole(lib, Number(select.value), select));
        controls.appendChild(select);

        if (Number(lib.role) & 4) {
            const makeDefault = button(lib.default_save ? 'Default save' : 'Make save destination', 'cl-mini', () => setDefault(lib));
            makeDefault.disabled = busy || !lib.exists || !!lib.default_save;
            controls.appendChild(makeDefault);
        }

        const remove = button('Remove', 'cl-mini danger', () => remove(lib));
        remove.disabled = busy;
        controls.appendChild(remove);
        row.appendChild(controls);
        return row;
    }

    async function addDefault(category) {
        const app = api();
        if (busy || typeof app?.AddDefaultContentLibrary !== 'function') return;
        busy = true;
        try {
            await app.AddDefaultContentLibrary(category, 7);
            notify('Default ' + category + ' library added.', 'success');
            await refresh();
        } catch (err) {
            notify('Could not add default library: ' + (err?.message || err), 'error');
        } finally {
            busy = false;
        }
    }

    async function addExisting(category) {
        const app = api();
        if (busy || typeof app?.AddContentLibrary !== 'function') return;
        let path = '';
        try {
            path = typeof app.OpenContentDirectoryPicker === 'function'
                ? await app.OpenContentDirectoryPicker()
                : (typeof app.OpenDirectoryPicker === 'function' ? await app.OpenDirectoryPicker() : '');
        } catch (err) {
            notify('Could not open the folder picker: ' + (err?.message || err), 'error');
            return;
        }
        if (!path) return;
        busy = true;
        try {
            await app.AddContentLibrary(path, category, 'external', 'unknown', 1);
            notify('External library added read-only. Its files were not modified.', 'success');
            await refresh();
        } catch (err) {
            notify('Could not add library: ' + (err?.message || err), 'error');
        } finally {
            busy = false;
        }
    }

    async function changeRole(lib, role, select) {
        const app = api();
        if (busy || typeof app?.SetContentLibraryRole !== 'function') return;
        busy = true;
        try {
            await app.SetContentLibraryRole(lib.id, role);
            notify('Library permission updated to ' + roleLabel(role) + '.', 'success');
            await refresh();
        } catch (err) {
            select.value = String(lib.role);
            notify('Could not update permission: ' + (err?.message || err), 'error');
        } finally {
            busy = false;
        }
    }

    async function setDefault(lib) {
        const app = api();
        if (busy || typeof app?.SetDefaultContentLibrary !== 'function') return;
        busy = true;
        try {
            await app.SetDefaultContentLibrary(lib.id);
            notify('New content will be saved in this library.', 'success');
            await refresh();
        } catch (err) {
            notify('Could not set save destination: ' + (err?.message || err), 'error');
        } finally {
            busy = false;
        }
    }

    async function remove(lib) {
        const app = api();
        if (busy || typeof app?.RemoveContentLibrary !== 'function') return;
        if (!confirm('Remove this library from Howl?\n\n' + lib.path + '\n\nOnly the registration is removed. The folder and its files stay exactly where they are.')) return;
        busy = true;
        try {
            await app.RemoveContentLibrary(lib.id);
            notify('Library removed. The folder was not deleted.', 'success');
            await refresh();
        } catch (err) {
            notify('Could not remove library: ' + (err?.message || err), 'error');
        } finally {
            busy = false;
        }
    }

    window.contentLibraries = { refresh };

    document.addEventListener('DOMContentLoaded', () => {
        if (document.getElementById('content-libraries-manager')) refresh();
    });
})();
