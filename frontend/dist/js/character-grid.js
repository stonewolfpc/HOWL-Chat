/**
 * Character Grid
 *
 * Loads, renders, edits, and deletes characters in the card grid.
 * Also owns the available-characters cache loader used by the form.
 *
 * Depends on: character-state.js, character-images.js,
 *             character-lorebooks.js, character-form.js
 *
 * @module character-grid
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Cache Loader (used by form popup open)
// ---------------------------------------------------------------------------

function loadAvailableCharacters() {
    const cs     = window.characterState;
    let chars    = cs.characterCache || JSON.parse(localStorage.getItem('howl-chat-characters') || '[]');
    chars        = cs.migrateCharactersWithIds(chars);
    cs.characterCache = chars;
}

// ---------------------------------------------------------------------------
// Grid Load
// ---------------------------------------------------------------------------

function loadCharacters() {
	return (async () => {
		try {
			const cs = window.characterState;
			let characters = await window.CharacterStorage.loadAll();
			characters = cs.migrateCharactersWithIds(characters);
			const external = await loadExternalCharacters();
			cs.characterCache = characters.concat(external);
			displayCharacters(cs.characterCache);
		} catch (err) {
			console.error('Failed to load characters:', err);
			loadCharactersFromStorage();
		}
	})();
}

async function loadExternalCharacters() {
	const app = window.go?.main?.App;
	if (typeof app?.GetContentLibraryCharacters !== 'function') return [];
	try {
		const items = await app.GetContentLibraryCharacters();
		return (Array.isArray(items) ? items : []).map(item => ({
			id: `external_${item.id}`,
			name: item.name || 'Unnamed external character',
			description: item.description || '',
			personality: item.personality || '',
			scenario: item.scenario || '',
			first_mes: item.first_mes || '',
			mes_example: item.mes_example || '',
			metadata: { role: item.format === 'sillytavern' ? 'SillyTavern card' : 'External card' },
			_external: true,
			_externalPath: item.path || '',
			_externalSource: item.source || 'external',
			_externalFormat: item.format || 'unknown',
			_externalEditable: !!item.editable,
			_libraryId: item.library_id || ''
		})).filter(item => item._externalSource !== 'howl');
	} catch (err) {
		console.warn('External character libraries unavailable:', err);
		return [];
	}
}

function loadCharactersFromStorage() {
    const cs   = window.characterState;
    let chars  = cs.characterCache || JSON.parse(localStorage.getItem('howl-chat-characters') || '[]');
    chars      = cs.migrateCharactersWithIds(chars);
    cs.characterCache = chars;
    return displayCharacters(chars);
}

// ---------------------------------------------------------------------------
// Grid Render
// ---------------------------------------------------------------------------

function displayCharacters(characters) {
    const grid = document.getElementById('character-grid');
    if (!grid) return;

    if (!characters || characters.length === 0) {
        grid.innerHTML = `
            <div class="character-grid-empty">
                <div class="character-grid-empty-icon">👥</div>
                <div class="character-grid-empty-text">No characters created yet</div>
                <div class="character-grid-empty-subtext">Click "Create New" to get started</div>
            </div>
        `;
        return;
    }

	grid.innerHTML = characters.map((character, index) => {
        const thumbnail = character.howl?.images?.find(img => img.isThumbnail);
        const avatarSrc = thumbnail ? thumbnail.data : (character.avatar || 'assets/Character.png');
		const external = !!character._external;
		const sourceLabel = external
			? (character._externalFormat === 'sillytavern' ? 'SillyTavern' : 'External library')
			: '';
		const name = escapeCharacterText(character.name || 'Unnamed');
        return `
			<div class="character-bubble${external ? ' character-bubble-external' : ''}" onclick="${external ? `viewExternalCharacter(${index})` : `editCharacter(${index})`}" title="${external ? escapeCharacterText(character._externalPath || '') : ''}">
                <div class="character-avatar">
					<img src="${avatarSrc}" alt="${name}" class="character-avatar-img">
                </div>
                <div class="character-info">
					<div class="character-name">${name}</div>
					${external ? `<div class="character-external-badge">${escapeCharacterText(sourceLabel)} · Read only</div>` : ''}
                    <div class="character-meta">
                        ${(() => {
                            const role = character.metadata?.role || '';
                            const age  = character.metadata?.age  || '';
                            if (role) {
                                const display = role.length > 40 ? role.slice(0, 38) + '…' : role;
                                return `<span class="character-role" title="${role}">${display}</span>`;
                            }
                            if (age) return `<span class="character-role">Age ${age}</span>`;
                            return `<span class="character-role character-role-empty">No role specified</span>`;
                        })()}
                    </div>
                </div>
				<div class="character-actions">
					${external
						? `<button class="character-edit-btn" onclick="viewExternalCharacter(${index}, event)" title="View external card">👁</button>`
						: `<button class="character-edit-btn" onclick="editCharacter(${index}, event)">✏️</button>
						   <button class="character-delete-btn" onclick="deleteCharacter(${index}, event)">🗑️</button>`}
				</div>
            </div>
        `;
    }).join('');
}

function escapeCharacterText(value) {
	return String(value ?? '').replace(/[&<>"']/g, ch => ({
		'&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;'
	}[ch]));
}

// ---------------------------------------------------------------------------
// Edit
// ---------------------------------------------------------------------------

function editCharacter(index, event) {
    if (event) event.stopPropagation();

    const cs         = window.characterState;
    cs.currentEditingIndex = index;

    const characters = cs.characterCache || JSON.parse(localStorage.getItem('howl-chat-characters') || '[]');
    const character  = characters[index];
    if (!character) { console.error('Character not found at index:', index); return; }

    const set = (id, val) => { const el = document.getElementById(id); if (el) el.value = val; };
    const setChecked = (id, val) => { const el = document.getElementById(id); if (el) el.checked = !!val; };
    set('character-name',              character.name                     || '');
    set('character-age',               character.metadata?.age            || '');
    set('character-role',              character.metadata?.role           || '');
    set('character-priority',          character.metadata?.priority       ?? 3);
    set('character-description',       character.description              || '');
    set('character-personality',       character.personality              || '');
    set('character-background',        character.metadata?.background     || '');
    set('character-intro',             character.first_mes                || '');
    set('character-dialogue-examples', character.mes_example              || '');

    // Voice override fields
    const voice = character.howl?.voice;
    setChecked('char-voice-enabled', voice?.enabled || false);
    set('char-voice-engine', voice?.engine || '');
    set('char-voice-model',  voice?.voice  || '');
    set('char-voice-speed',  voice?.speed  ?? 1.0);
    set('char-voice-volume', voice?.volume ?? 0.8);
    if (typeof _updateVoiceControlsVisible === 'function') _updateVoiceControlsVisible();
    if (typeof _updateSliderDisplay === 'function') {
        _updateSliderDisplay('char-voice-speed', 'char-voice-speed-val', 'x');
        _updateSliderDisplay('char-voice-volume', 'char-voice-volume-val', '');
    }

    cs.characterImages    = character.howl?.images || [];
    updateImagesList();

    cs.selectedLorebooks  = cs.resolveLorebooksFromRefs(character.howl?.lorebooks || character.lorebooks || []);
    updateSelectedLorebooksList();
    loadAvailableLorebooks();

    if (window.go?.main?.App?.SetLorebooks) {
        window.go.main.App.SetLorebooks(Array.isArray(cs.selectedLorebooks) ? cs.selectedLorebooks : []);
    }

    openCharacterCreationPopup();
    if (typeof loadCharacterVoiceControls === 'function') {
        loadCharacterVoiceControls(voice?.voice || '', true);
    }
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

function deleteCharacter(index, event) {
    if (event) event.stopPropagation();

    const cs     = window.characterState;
    const chars  = cs.characterCache || [];
    const target = chars[index];
    if (!target) return;

    window.DeleteConfirm.show({
        title: `Delete "${target.name || 'this character'}"?`,
        body:  'The character and its images will be removed. They will be detached from all scenarios and the active chat. Lorebooks are kept.',
        onConfirm: () => {
            if (window.DeleteService) window.DeleteService.detachCharacterFromAll(target.id);
            window.CharacterStorage.remove(target.id)
                .then(() => {
                    cs.characterCache = null;
                    loadCharacters();
                })
                .catch(err => {
                    console.error('Failed to delete character:', err);
                    chars.splice(index, 1);
                    cs.characterCache = chars;
                    localStorage.setItem('howl-chat-characters', JSON.stringify(chars));
                    loadCharacters();
                });
        }
    });
}
