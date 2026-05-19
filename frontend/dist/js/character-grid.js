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
    window.CharacterStorage.loadAll()
        .then(characters => {
            const cs = window.characterState;
            characters = cs.migrateCharactersWithIds(characters);
            cs.characterCache = characters;
            displayCharacters(characters);
        })
        .catch(err => {
            console.error('Failed to load characters:', err);
            loadCharactersFromStorage();
        });
}

function loadCharactersFromStorage() {
    const cs   = window.characterState;
    let chars  = cs.characterCache || JSON.parse(localStorage.getItem('howl-chat-characters') || '[]');
    chars      = cs.migrateCharactersWithIds(chars);
    cs.characterCache = chars;
    displayCharacters(chars);
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
        return `
            <div class="character-bubble" onclick="editCharacter(${index})">
                <div class="character-avatar">
                    <img src="${avatarSrc}" alt="${character.name}" class="character-avatar-img">
                </div>
                <div class="character-info">
                    <div class="character-name">${character.name}</div>
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
                    <button class="character-edit-btn"   onclick="editCharacter(${index}, event)">✏️</button>
                    <button class="character-delete-btn" onclick="deleteCharacter(${index}, event)">🗑️</button>
                </div>
            </div>
        `;
    }).join('');
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
