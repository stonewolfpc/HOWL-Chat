/**
 * HOWL Chat Lorebooks
 *
 * A lorebook entry is a scoped memory rule. It decides when lore wakes up,
 * how important it is, and how safely it enters the model context.
 */

const LOREBOOK_STORAGE_KEY = 'howl-chat-lorebooks';

const LOREBOOK_DEFAULTS = {
    title: '',
    content: '',
    tags: [],
    scope: 'world',
    enabled: true,
    triggerPhrases: [],
    secondaryTriggers: [],
    triggerMode: 'exact',
    triggerDirection: 'both',
    triggerFrequency: 'always',
    priorityLevel: 3,
    conflictRule: 'higher_priority',
    visibility: 'public',
    ownerId: '',
    revealCondition: 'never',
    injectionStyle: 'summary',
    injectionPosition: 'system_summary_only',
    maxLength: 'medium',
    linkedCharacters: [],
    linkedLoreEntries: [],
    inheritance: 'none',
    contextBudget: 160,
    scanWindow: 4,
    recursion: 'limited'
};

const SCOPE_META = {
    world: {
        label: 'World',
        description: 'Global rules, places, history, and truths for the active world.'
    },
    scenario: {
        label: 'Scenario',
        description: 'Local events, scene rules, NPCs, and situation-specific facts.'
    },
    character: {
        label: 'Character',
        description: 'Personal memories, secrets, traits, and backstory for one AI.'
    }
};

const SELECT_LABELS = {
    triggerMode: {
        loose: 'Loose Match',
        exact: 'Exact Match',
        semantic: 'Semantic Match'
    },
    triggerDirection: {
        user: 'User Only',
        character: 'Characters Only',
        both: 'Both'
    },
    triggerFrequency: {
        always: 'Always',
        once_per_chat: 'Once Per Chat',
        once_per_scene: 'Once Per Scene',
        once_per_message: 'Once Per Message'
    },
    conflictRule: {
        higher_priority: 'Use Higher Priority',
        newer_entry: 'Use Newer Entry',
        use_both: 'Use Both'
    },
    visibility: {
        public: 'Public',
        private: 'Private',
        hidden: 'Hidden'
    },
    revealCondition: {
        character_says: 'Character Says It',
        user_reveals: 'User Reveals It',
        scenario_reveals: 'Scenario Reveals It',
        never: 'Never Reveal'
    },
    injectionStyle: {
        inline: 'Inline',
        summary: 'Summary',
        memory: 'Memory',
        background: 'Background'
    },
    injectionPosition: {
        before_message: 'Before Message',
        after_message: 'After Message',
        system_summary_only: 'System Summary Only',
        hidden: 'Hidden'
    },
    maxLength: {
        short: 'Short',
        medium: 'Medium',
        full: 'Full Entry'
    },
    inheritance: {
        child_scenarios: 'Inherit to Child Scenarios',
        child_worlds: 'Inherit to Child Worlds',
        none: 'Do Not Inherit'
    },
    recursion: {
        off: 'Off',
        limited: 'Limited',
        full: 'Full'
    }
};

let lorebooks = [];
let selectedLorebookType = 'world';
let currentEditingId = null;

function lorebookRepo() {
    return window.LorebookRepo || null;
}

function splitList(value) {
    return (value || '')
        .split(',')
        .map(item => item.trim())
        .filter(Boolean);
}

function joinList(value) {
    return Array.isArray(value) ? value.join(', ') : '';
}

function optionLabel(group, value) {
    return SELECT_LABELS[group]?.[value] || value || '';
}

function normalizeLorebook(raw = {}) {
    const content = raw.content || (Array.isArray(raw.entries) ? raw.entries.map(entry => entry.content).join('\n\n') : '');
    const priority = Number(raw.priorityLevel || raw.priority || LOREBOOK_DEFAULTS.priorityLevel);
    const scope = raw.scope || raw.type || raw.visibility || LOREBOOK_DEFAULTS.scope;

    return {
        ...LOREBOOK_DEFAULTS,
        ...raw,
        id: raw.id || Date.now().toString(),
        title: raw.title || raw.name || '',
        content,
        tags: Array.isArray(raw.tags) ? raw.tags : splitList(raw.tags),
        scope,
        type: scope,
        enabled: raw.enabled !== false,
        triggerPhrases: Array.isArray(raw.triggerPhrases) ? raw.triggerPhrases : splitList(raw.triggerPhrases || raw.keys),
        secondaryTriggers: Array.isArray(raw.secondaryTriggers) ? raw.secondaryTriggers : splitList(raw.secondaryTriggers || raw.secondaryKeys),
        priorityLevel: Math.min(5, Math.max(1, priority)),
        contextBudget: Number(raw.contextBudget || LOREBOOK_DEFAULTS.contextBudget),
        scanWindow: Number(raw.scanWindow || LOREBOOK_DEFAULTS.scanWindow),
        createdAt: raw.createdAt || new Date().toISOString(),
        updatedAt: raw.updatedAt || new Date().toISOString()
    };
}

function getField(id) {
    return document.getElementById(id);
}

function setValue(id, value) {
    const field = getField(id);
    if (!field) return;

    if (field.type === 'checkbox') {
        field.checked = Boolean(value);
    } else {
        field.value = value ?? '';
    }
}

function getValue(id) {
    const field = getField(id);
    if (!field) return '';
    return field.type === 'checkbox' ? field.checked : field.value;
}

function loadLorebooksFromStorage() {
    try {
        const saved = JSON.parse(localStorage.getItem(LOREBOOK_STORAGE_KEY) || '[]');
        lorebooks = Array.isArray(saved) ? saved.map(normalizeLorebook) : [];
    } catch (error) {
        console.error('Failed to load lorebooks:', error);
        lorebooks = [];
    }
}

function saveLorebooksToStorage() {
    localStorage.setItem(LOREBOOK_STORAGE_KEY, JSON.stringify(lorebooks));
}

function selectLorebookType(type) {
    selectedLorebookType = type;
    setValue('lorebook-scope', type);

    document.querySelectorAll('.lorebook-type-item').forEach(item => {
        item.classList.toggle('selected', item.dataset.scope === type);
    });

    const helper = getField('lorebook-scope-helper');
    if (helper) {
        helper.textContent = SCOPE_META[type]?.description || '';
    }

    updateLorebookOwnerVisibility();
}

// Ownership is a separate relationship from linkedCharacters. Only the owner
// receives private character lore; links never grant that authority.
function updateLorebookOwnerVisibility() {
    const group = getField('lorebook-owner-group');
    const owner = getField('lorebook-owner');
    if (!group || !owner) return;
    const isCharacterScope = (getValue('lorebook-scope') || selectedLorebookType) === 'character';
    group.style.display = isCharacterScope ? '' : 'none';
    owner.disabled = !isCharacterScope;
}

const OWNER_MIGRATION_VERSION = 'howl-character-owner-v1';

function ownerIdentities(characters) {
    const seen = new Set();
    return (characters || [])
        .map(character => ({
            id: String(character.id || character.ID || '').trim(),
            name: String(character.name || character.Name || character.title || character.Title || '').trim()
        }))
        .filter(character => {
            if (!character.id || !character.name || seen.has(character.id)) return false;
            seen.add(character.id);
            return true;
        });
}

// Migrates legacy display-name owners to stable character IDs. The backend only
// rewrites exact, unique names and leaves ambiguous or unknown values untouched;
// this page must not guess either. Re-running is safe.
async function migrateLorebookOwners(records) {
    const app = window.go?.main?.App;
    if (!app || typeof app.MigrateLorebookOwners !== 'function') return null;
    const identities = ownerIdentities(records);
    if (!identities.length) return null;
    try {
        const result = await app.MigrateLorebookOwners(identities, OWNER_MIGRATION_VERSION);
        const items = Array.isArray(result?.items) ? result.items : [];
        const remap = new Map(items
            .filter(item => item.action === 'migrated' && item.owner_name && item.owner_id)
            .map(item => [item.owner_name, item.owner_id]));
        if (remap.size) {
            lorebooks.forEach(lorebook => {
                if (remap.has(lorebook.ownerId)) lorebook.ownerId = remap.get(lorebook.ownerId);
            });
        }
        if (result && (result.migrated || result.ambiguous || result.unknown || result.failed)) {
            const status = getField('lorebook-store-status');
            if (status) {
                const parts = [];
                if (result.migrated) parts.push(`${result.migrated} updated`);
                if (result.ambiguous) parts.push(`${result.ambiguous} ambiguous`);
                if (result.unknown) parts.push(`${result.unknown} unknown`);
                if (result.failed) parts.push(`${result.failed} failed`);
                status.textContent = `Owner identity migration: ${parts.join(', ')}.`;
            }
        }
        return result;
    } catch (error) {
        console.warn('Owner identity migration did not run:', error);
        return null;
    }
}

async function loadLorebookOwnerOptions() {
    const select = getField('lorebook-owner');
    if (!select) return;
    let desired = select.dataset.pendingOwner ?? select.value;
    try {
        // Chat adds characters from the character repository. Ownership is
        // the stable id from that same library; the name is only a label.
        let characters = null;
        if (window.CharacterStorage && typeof window.CharacterStorage.loadAll === 'function') {
            characters = await window.CharacterStorage.loadAll();
        } else {
            const app = window.go?.main?.App;
            if (!app || typeof app.GetCharacters !== 'function') return;
            characters = await app.GetCharacters();
        }
        const records = ownerIdentities(characters);
        const migration = await migrateLorebookOwners(records);
        const remap = new Map((migration?.items || [])
            .filter(item => item.action === 'migrated' && item.owner_name && item.owner_id)
            .map(item => [item.owner_name, item.owner_id]));
        if (remap.has(desired)) desired = remap.get(desired);

        const nameCounts = new Map();
        records.forEach(record => nameCounts.set(record.name, (nameCounts.get(record.name) || 0) + 1));
        records.sort((a, b) => a.name.localeCompare(b.name));
        const known = new Set(records.map(record => record.id));
        select.innerHTML = '<option value="">Unowned</option>' + records
            .map(record => {
                const suffix = nameCounts.get(record.name) > 1 ? ` (${record.id.slice(-6)})` : '';
                return `<option value="${escapeHtml(record.id)}" data-owner-name="${escapeHtml(record.name)}">${escapeHtml(record.name + suffix)}</option>`;
            })
            .join('');
        if (desired && !known.has(desired)) {
            // Legacy name or a deleted character: keep the value and show that it
            // is unresolved. Never silently turn it into an unowned entry.
            select.innerHTML += `<option value="${escapeHtml(desired)}" data-owner-name="">Unresolved owner: ${escapeHtml(desired)}</option>`;
        }
        select.value = desired;
        delete select.dataset.pendingOwner;
    } catch (error) {
        console.warn('Could not load character owners:', error);
    }
}

function clearLorebookForm() {
    const defaults = { ...LOREBOOK_DEFAULTS, scope: selectedLorebookType };

    setValue('lorebook-title', defaults.title);
    setValue('lorebook-content', defaults.content);
    setValue('lorebook-tags', '');
    setValue('lorebook-enabled', defaults.enabled);
    setValue('lorebook-trigger-phrases', '');
    setValue('lorebook-secondary-triggers', '');
    setValue('lorebook-trigger-mode', defaults.triggerMode);
    setValue('lorebook-trigger-direction', defaults.triggerDirection);
    setValue('lorebook-trigger-frequency', defaults.triggerFrequency);
    setValue('lorebook-priority-level', defaults.priorityLevel);
    setValue('lorebook-conflict-rule', defaults.conflictRule);
    setValue('lorebook-visibility', defaults.visibility);
    setValue('lorebook-reveal-condition', defaults.revealCondition);
    setValue('lorebook-injection-style', defaults.injectionStyle);
    setValue('lorebook-injection-position', defaults.injectionPosition);
    setValue('lorebook-max-length', defaults.maxLength);
    setValue('lorebook-context-budget', defaults.contextBudget);
    setValue('lorebook-scan-window', defaults.scanWindow);
    setValue('lorebook-recursion', defaults.recursion);
    setValue('lorebook-linked-characters', '');
    setValue('lorebook-linked-entries', '');
    setValue('lorebook-owner', '');
    setValue('lorebook-inheritance', defaults.inheritance);

    selectLorebookType(selectedLorebookType);
}

function openLorebookCreationPopup(resetForm = true) {
    const overlay = getField('lorebook-creation-overlay');
    if (!overlay) return;

    if (resetForm) {
        currentEditingId = null;
        selectedLorebookType = selectedLorebookType || 'world';
        clearLorebookForm();
    }

    overlay.style.display = 'flex';
    loadLorebookOwnerOptions();
    const title = getField('lorebook-editor-title');
    if (title) {
        title.textContent = currentEditingId ? 'Edit Lorebook Entry' : 'Create Lorebook Entry';
    }
}

function closeLorebookCreationPopup() {
    const overlay = getField('lorebook-creation-overlay');
    if (overlay) {
        overlay.style.display = 'none';
    }
    currentEditingId = null;
    clearLorebookForm();
}

function buildLorebookFromForm() {
    const now = new Date().toISOString();
    const existing = currentEditingId ? lorebooks.find(item => item.id === currentEditingId) : null;
    const scope = getValue('lorebook-scope') || selectedLorebookType;

    return normalizeLorebook({
        id: currentEditingId || Date.now().toString(),
        title: getValue('lorebook-title').trim(),
        content: getValue('lorebook-content').trim(),
        tags: splitList(getValue('lorebook-tags')),
        scope,
        type: scope,
        ownerId: scope === 'character' ? getValue('lorebook-owner').trim() : '',
        ownerName: scope === 'character'
            ? (getField('lorebook-owner')?.selectedOptions?.[0]?.dataset.ownerName || '')
            : '',
        enabled: getValue('lorebook-enabled'),
        triggerPhrases: splitList(getValue('lorebook-trigger-phrases')),
        secondaryTriggers: splitList(getValue('lorebook-secondary-triggers')),
        triggerMode: getValue('lorebook-trigger-mode'),
        triggerDirection: getValue('lorebook-trigger-direction'),
        triggerFrequency: getValue('lorebook-trigger-frequency'),
        priorityLevel: Number(getValue('lorebook-priority-level')),
        conflictRule: getValue('lorebook-conflict-rule'),
        visibility: getValue('lorebook-visibility'),
        revealCondition: getValue('lorebook-reveal-condition'),
        injectionStyle: getValue('lorebook-injection-style'),
        injectionPosition: getValue('lorebook-injection-position'),
        maxLength: getValue('lorebook-max-length'),
        contextBudget: Number(getValue('lorebook-context-budget')),
        scanWindow: Number(getValue('lorebook-scan-window')),
        recursion: getValue('lorebook-recursion'),
        linkedCharacters: splitList(getValue('lorebook-linked-characters')),
        linkedLoreEntries: splitList(getValue('lorebook-linked-entries')),
        inheritance: getValue('lorebook-inheritance'),
        createdAt: existing?.createdAt || now,
        updatedAt: now
    });
}

async function saveLorebook() {
    const lorebook = buildLorebookFromForm();

    if (!lorebook.title) {
        alert('Lorebook title is required.');
        return;
    }

    if (!lorebook.content) {
        alert('Lorebook content is required.');
        return;
    }

    if (lorebook.scope === 'character' && lorebook.visibility === 'private' && !lorebook.ownerId) {
        const proceed = confirm('This private character entry has no owner. No character will receive it. Save it as unowned?');
        if (!proceed) return;
    }

    if (lorebook.triggerPhrases.length === 0 && lorebook.injectionStyle !== 'background') {
        const proceed = confirm('This entry has no trigger phrases. It may never activate unless called manually. Save it anyway?');
        if (!proceed) return;
    }

    const previous = lorebooks.slice();
    const index = lorebooks.findIndex(item => item.id === lorebook.id);
    if (index >= 0) {
        lorebooks[index] = lorebook;
    } else {
        lorebooks.push(lorebook);
    }

    try {
        const repo = lorebookRepo();
        if (repo) {
            await repo.save(lorebook, index >= 0);
        } else {
            saveLorebooksToStorage();
        }
    } catch (error) {
        lorebooks = previous;
        console.error('Failed to save lorebook:', error);
        alert('Failed to save lorebook. ' + (error?.message || error));
        return;
    }
    closeLorebookCreationPopup();
    await loadLorebooks();
}

function renderTagList(items, className = 'lorebook-chip') {
    if (!items || items.length === 0) return '';
    return items.slice(0, 4).map(item => `<span class="${className}">${item}</span>`).join('');
}

async function loadLorebooks() {
    const grid = getField('lorebooks-grid');
    if (!grid) return;

    const repo = lorebookRepo();
    if (repo) {
        try {
            const state = await repo.hydrate();
            lorebooks = (state.entries || []).map(normalizeLorebook);
        } catch (error) {
            console.error('Failed to load lorebooks from the library:', error);
            loadLorebooksFromStorage();
        }
    } else {
        loadLorebooksFromStorage();
    }
    renderLorebookStoreStatus();

    if (lorebooks.length === 0) {
        grid.innerHTML = `
            <div class="world-grid-empty">
                <div class="world-grid-empty-icon">BOOK</div>
                <div class="world-grid-empty-text">No lorebook entries yet</div>
                <div class="world-grid-empty-subtext">Create the first memory rule for your world, scenario, or character.</div>
            </div>
        `;
        return;
    }

    const sorted = [...lorebooks]
        .filter(lorebook => !lorebook.automated)
        .sort((a, b) => {
            const scopeOrder = { world: 0, scenario: 1, character: 2 };
            return (scopeOrder[a.scope] ?? 9) - (scopeOrder[b.scope] ?? 9) || a.priorityLevel - b.priorityLevel;
        });

    grid.innerHTML = sorted.map(lorebook => `
        <div class="world-bubble lorebook-card ${lorebook.enabled ? '' : 'is-disabled'}" onclick="editLorebookById('${lorebook.id}')">
            <div class="world-info">
                <div class="world-name">${escapeHtml(lorebook.title)}</div>
                <div class="lorebook-card-content">${escapeHtml(lorebook.content).slice(0, 180)}${lorebook.content.length > 180 ? '...' : ''}</div>
                <div class="world-meta lorebook-card-meta">
                    <span class="world-character-count">${SCOPE_META[lorebook.scope]?.label || lorebook.scope}</span>
                    <span class="world-character-count">Priority ${lorebook.priorityLevel}</span>
                    <span class="world-character-count">${optionLabel('injectionStyle', lorebook.injectionStyle)}</span>
                    <span class="world-character-count">${optionLabel('maxLength', lorebook.maxLength)}</span>
                </div>
                <div class="lorebook-chip-row">${renderTagList(lorebook.triggerPhrases)}${renderTagList(lorebook.tags, 'lorebook-chip muted')}</div>
            </div>
            <div class="world-actions">
                <button class="world-bubble-action" onclick="event.stopPropagation(); editLorebookById('${lorebook.id}')" title="Edit"><span class="icon icon-edit"></span></button>
                <button class="world-bubble-action" onclick="event.stopPropagation(); deleteLorebookById('${lorebook.id}')" title="Delete"><span class="icon icon-trash"></span></button>
            </div>
        </div>
    `).join('');
}

function populateForm(lorebook) {
    currentEditingId = lorebook.id;
    selectedLorebookType = lorebook.scope;

    setValue('lorebook-title', lorebook.title);
    setValue('lorebook-content', lorebook.content);
    setValue('lorebook-tags', joinList(lorebook.tags));
    setValue('lorebook-scope', lorebook.scope);
    setValue('lorebook-enabled', lorebook.enabled);
    setValue('lorebook-trigger-phrases', joinList(lorebook.triggerPhrases));
    setValue('lorebook-secondary-triggers', joinList(lorebook.secondaryTriggers));
    setValue('lorebook-trigger-mode', lorebook.triggerMode);
    setValue('lorebook-trigger-direction', lorebook.triggerDirection);
    setValue('lorebook-trigger-frequency', lorebook.triggerFrequency);
    setValue('lorebook-priority-level', lorebook.priorityLevel);
    setValue('lorebook-conflict-rule', lorebook.conflictRule);
    setValue('lorebook-visibility', lorebook.visibility);
    setValue('lorebook-reveal-condition', lorebook.revealCondition);
    setValue('lorebook-injection-style', lorebook.injectionStyle);
    setValue('lorebook-injection-position', lorebook.injectionPosition);
    setValue('lorebook-max-length', lorebook.maxLength);
    setValue('lorebook-context-budget', lorebook.contextBudget);
    setValue('lorebook-scan-window', lorebook.scanWindow);
    setValue('lorebook-recursion', lorebook.recursion);
    setValue('lorebook-linked-characters', joinList(lorebook.linkedCharacters));
    setValue('lorebook-linked-entries', joinList(lorebook.linkedLoreEntries));
    setValue('lorebook-inheritance', lorebook.inheritance);

    const ownerField = getField('lorebook-owner');
    if (ownerField) ownerField.dataset.pendingOwner = lorebook.ownerId || '';
    setValue('lorebook-owner', lorebook.ownerId || '');

    selectLorebookType(lorebook.scope);
}

function editLorebookById(id) {
    const lorebook = lorebooks.find(item => item.id === id);
    if (!lorebook) return;

    populateForm(lorebook);
    openLorebookCreationPopup(false);
}

function editLorebook(index) {
    const lorebook = lorebooks[index];
    if (lorebook) {
        editLorebookById(lorebook.id);
    }
}

async function deleteLorebookById(id) {
    const lorebook = lorebooks.find(item => item.id === id);
    if (!lorebook) return;

    if (!confirm(`Delete "${lorebook.title}"?`)) return;

    try {
        const repo = lorebookRepo();
        if (repo) {
            await repo.remove(id);
        } else {
            lorebooks = lorebooks.filter(item => item.id !== id);
            saveLorebooksToStorage();
        }
    } catch (error) {
        console.error('Failed to delete lorebook:', error);
        alert('Failed to delete lorebook. ' + (error?.message || error));
        return;
    }
    await loadLorebooks();
}

function deleteLorebook(index) {
    const lorebook = lorebooks[index];
    if (lorebook) {
        deleteLorebookById(lorebook.id);
    }
}

function renderLorebookStoreStatus() {
    const status = getField('lorebook-store-status');
    const importButton = getField('lorebook-import-button');
    const rollbackButton = getField('lorebook-rollback-button');
    const repo = lorebookRepo();
    const state = repo ? repo.state : null;

    if (status) {
        if (!repo || !state) {
            status.textContent = 'Browser storage only.';
        } else if (state.cutover) {
            const roots = state.status && Array.isArray(state.status.roots) ? state.status.roots.length : 0;
            const conflicts = state.status && Array.isArray(state.status.conflicts) ? state.status.conflicts.length : 0;
            const count = state.entries.length;
            status.textContent = `Library files in ${roots} folder${roots === 1 ? '' : 's'} · ${count} entr${count === 1 ? 'y' : 'ies'}`
                + (conflicts ? ` · ${conflicts} identity conflict${conflicts === 1 ? '' : 's'} need review` : '');
        } else {
            status.textContent = 'Stored in browser storage. Import to keep lorebooks as files in your library folder.';
        }
        status.title = status.textContent;
    }
    if (importButton) importButton.disabled = Boolean(state && state.cutover);
    if (rollbackButton) rollbackButton.disabled = !(state && state.cutover);
}

async function importLorebooksToLibrary() {
    const repo = lorebookRepo();
    if (!repo) {
        alert('The lorebook library is not available in this window.');
        return;
    }
    const entries = repo.readCache();
    if (entries.length === 0) {
        alert('There are no browser-stored lorebooks to import.');
        return;
    }
    const ok = confirm(`Import ${entries.length} lorebook entr${entries.length === 1 ? 'y' : 'ies'} into the configured library folder?\n\nFiles are created beside your other content. The browser copy is kept and this can be undone.`);
    if (!ok) return;
    try {
        const result = await repo.importLegacy();
        await loadLorebooks();
        alert(result.message || 'Import complete.');
    } catch (error) {
        console.error('Lorebook import failed:', error);
        alert('Import failed. ' + (error?.message || error));
    }
}

async function rollbackLorebookLibrary() {
    const repo = lorebookRepo();
    if (!repo) {
        alert('The lorebook library is not available in this window.');
        return;
    }
    const ok = confirm('Remove the imported lorebook files and return to browser storage?\n\nOnly files created by the last import are removed. Your browser lorebooks are not changed.');
    if (!ok) return;
    try {
        const result = await repo.rollback();
        await loadLorebooks();
        alert(result.message || 'Import undone.');
    } catch (error) {
        console.error('Lorebook import rollback failed:', error);
        alert('Undo failed. ' + (error?.message || error));
    }
}

function exportLorebooksForBackend() {
    loadLorebooksFromStorage();
    return lorebooks.map(normalizeLorebook);
}

document.addEventListener('DOMContentLoaded', function() {
    loadLorebooks();
    clearLorebookForm();

    getField('lorebook-scope')?.addEventListener('change', event => selectLorebookType(event.target.value));

    document.addEventListener('keydown', function(event) {
        if (event.key === 'Escape') {
            closeLorebookCreationPopup();
        }

        if (event.ctrlKey && event.key.toLowerCase() === 'n') {
            event.preventDefault();
            openLorebookCreationPopup();
        }
    });

    getField('lorebook-creation-overlay')?.addEventListener('click', function(event) {
        if (event.target === this) {
            closeLorebookCreationPopup();
        }
    });
});

window.howlLorebooks = {
    exportLorebooksForBackend,
    loadLorebooks
};

window.lorebooksModule = {
    openLorebookCreationPopup,
    closeLorebookCreationPopup,
    saveLorebook,
    loadLorebooks,
    editLorebook,
    editLorebookById,
    deleteLorebook,
    deleteLorebookById,
    selectLorebookType
};
