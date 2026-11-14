/**
 * Character Form
 *
 * Popup open/close for character creation and user details,
 * data collection, saveCharacter, saveUserDetails, loadUserDetails.
 *
 * Depends on: character-state.js, character-images.js,
 *             character-lorebooks.js, character-grid.js
 *
 * @module character-form
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// Character Popup Lifecycle
// ---------------------------------------------------------------------------

function openCharacterCreationPopup() {
    document.getElementById('character-creation-overlay').style.display = 'flex';
	setCharacterFormReadOnly(false);
    loadAvailableCharacters();
    const cs = window.characterState;
    if (cs && typeof cs.resetDirtyFields === 'function') cs.resetDirtyFields();
    resetCharacterVoiceControls();
}

function closeCharacterCreationPopup() {
    document.getElementById('character-creation-overlay').style.display = 'none';
	setCharacterFormReadOnly(false);
    const cs = window.characterState;
    cs.currentEditingIndex  = -1;
    cs.characterImages      = [];
    cs.selectedLorebooks.length = 0;
    if (typeof cs.resetDirtyFields === 'function') cs.resetDirtyFields();
    const form = document.querySelector('.character-creation-form');
    if (form) form.reset();
}

function setCharacterFormReadOnly(readonly, character) {
	const overlay = document.getElementById('character-creation-overlay');
	if (!overlay) return;
	overlay.classList.toggle('read-only', !!readonly);
	const banner = document.getElementById('character-external-banner');
	if (banner) {
		banner.style.display = readonly ? 'block' : 'none';
		banner.textContent = readonly
			? `Read-only external card${character?._externalFormat === 'sillytavern' ? ' (SillyTavern)' : ''} · ${character?._externalPath || ''}. Howl is displaying the original file; it cannot edit or delete it here.`
			: '';
	}
	const save = document.getElementById('character-creation-save');
	if (save) save.style.display = readonly ? 'none' : '';
	overlay.querySelectorAll('input, textarea, select').forEach(field => { field.disabled = !!readonly; });
	const title = overlay.querySelector('.character-creation-title');
	if (title) title.textContent = readonly ? 'External Character' : (window.characterState?.currentEditingIndex >= 0 ? 'Edit Character' : 'Create New Character');
}

function viewExternalCharacter(index, event) {
	if (event) event.stopPropagation();
	const cs = window.characterState;
	const character = (cs.characterCache || [])[index];
	if (!character || !character._external) return;
	cs.currentEditingIndex = index;
	cs.currentEditingId = character.id;
	const set = (id, value) => { const el = document.getElementById(id); if (el) el.value = value || ''; };
	set('character-name', character.name);
	set('character-age', '');
	set('character-role', character.metadata?.role || '');
	set('character-priority', 3);
	set('character-description', character.description);
	set('character-personality', character.personality);
	set('character-background', '');
	set('character-intro', character.first_mes);
	set('character-dialogue-examples', character.mes_example);
	cs.characterImages = [];
	cs.selectedLorebooks = [];
	if (typeof updateImagesList === 'function') updateImagesList();
	if (typeof updateSelectedLorebooksList === 'function') updateSelectedLorebooksList();
	openCharacterCreationPopup();
	setCharacterFormReadOnly(true, character);
}

// ---------------------------------------------------------------------------
// User Details Popup Lifecycle
// ---------------------------------------------------------------------------

function openUserDetailsPopup() {
    document.getElementById('user-details-overlay').style.display = 'flex';
    loadUserDetails();
}

function closeUserDetailsPopup() {
    document.getElementById('user-details-overlay').style.display = 'none';
}

// ---------------------------------------------------------------------------
// Data Collection
// ---------------------------------------------------------------------------

function _readCharacterValue(id, fallback) {
    const el = document.getElementById(id);
    return el ? el.value : fallback;
}

/**
 * Only write voice fields the user actually changed. A fixed-length select can
 * display an empty value when a stored engine/voice is absent from its option
 * list; treating that empty value as an edit would erase valid stored settings.
 */
function collectVoicePatch(existing) {
    const cs             = window.characterState || {};
    const existingVoice  = existing?.howl?.voice;
    const forceFullPatch = !existingVoice || typeof existingVoice !== 'object' || Array.isArray(existingVoice);
    const fields = [
        { key: 'enabled', id: 'char-voice-enabled', read: () => document.getElementById('char-voice-enabled')?.checked || false },
        { key: 'engine',  id: 'char-voice-engine',  read: () => _readCharacterValue('char-voice-engine', '') || '' },
        { key: 'voice',   id: 'char-voice-model',   read: () => _readCharacterValue('char-voice-model', '') || '' },
        { key: 'speed',   id: 'char-voice-speed',   read: () => parseFloat(_readCharacterValue('char-voice-speed', '1.0') || '1.0') },
        { key: 'volume',  id: 'char-voice-volume',  read: () => parseFloat(_readCharacterValue('char-voice-volume', '0.8') || '0.8') }
    ];

    const patch = {};
    for (const field of fields) {
        const dirty = typeof cs.isFieldDirty === 'function' && cs.isFieldDirty(field.id);
        if (forceFullPatch || dirty) patch[field.key] = field.read();
    }
    return Object.keys(patch).length > 0 ? patch : null;
}

function collectCharacterData(existing) {
    const cs        = window.characterState;
    const thumbnail = cs.characterImages.find(img => img.isThumbnail);
    const now       = new Date().toISOString();
    const isEdit    = !!existing;

    const lorebooks = cs.selectedLorebooks.map(l => ({
        ...l,
        id:    l.id || cs.lorebookKey(l),
        name:  cs.lorebookDisplayName(l),
        title: cs.lorebookDisplayName(l)
    }));
    const images = cs.characterImages.map(img => ({ ...img }));
    const voice  = collectVoicePatch(existing);

    const data = {
        id:          existing?.id || `char_${Date.now()}_${Math.floor(Math.random() * 1000)}`,
        name:        _readCharacterValue('character-name', existing?.name || ''),
        description: _readCharacterValue('character-description', existing?.description || ''),
        personality: _readCharacterValue('character-personality', existing?.personality || ''),
        scenario:    _readCharacterValue('character-scenario', existing?.scenario || ''),
        first_mes:   _readCharacterValue('character-intro', existing?.first_mes || ''),
        mes_example: _readCharacterValue('character-dialogue-examples', existing?.mes_example || ''),
        avatar:      thumbnail ? thumbnail.data : (existing?.avatar || ''),
        tags:        Array.isArray(existing?.tags) ? existing.tags : [],
        metadata: {
            updated:    now,
            age:        _readCharacterValue('character-age', existing?.metadata?.age || ''),
            role:       _readCharacterValue('character-role', existing?.metadata?.role || ''),
            background: _readCharacterValue('character-background', existing?.metadata?.background || ''),
            priority:   Math.min(5, Math.max(1, parseInt(_readCharacterValue('character-priority', existing?.metadata?.priority ?? 3) || '3', 10)))
        },
        howl: {
            lorebooks,
            images
        }
    };

    if (!isEdit) {
        data.metadata.version = 1;
        data.metadata.created = now;
        data.metadata.seed    = Math.floor(Math.random() * 1000000000);
        data.howl.world        = '';
        data.howl.memories     = [];
        data.howl.notes        = '';
        data.howl.image_style  = 'fantasy-realistic';
    }
    if (voice) data.howl.voice = voice;

    return data;
}

// ---------------------------------------------------------------------------
// Save Character
// ---------------------------------------------------------------------------

async function saveCharacter() {
    const cs       = window.characterState;
    const isEdit   = cs.currentEditingIndex >= 0;
    const existing = isEdit
        ? (cs.characterCache || [])[cs.currentEditingIndex]
        : null;
	if (existing?._external) {
		alert('This is a read-only external card. Howl is displaying the original file and will not overwrite it.');
		return;
	}
    const voicePatch = collectVoicePatch(existing);
    if (voicePatch && !(await validateCharacterVoice(voicePatch))) return;
    const characterData = collectCharacterData(existing);
    if (!characterData.name) { alert('Character name is required!'); return; }

    // Preserve the original id when editing so IndexedDB images stay linked.
    // Storage performs the same safe merge for the other save entry points.
    if (isEdit && existing?.id) characterData.id = existing.id;
    const merged = window.CharacterStorage?.mergeRecord
        ? window.CharacterStorage.mergeRecord(existing, characterData)
        : characterData;

    return window.CharacterStorage.save(merged)
        .then(async () => {
            if (window.AudioService?.invalidateVoiceCache) {
                window.AudioService.invalidateVoiceCache(merged.id);
            }
            if (window.parent && window.parent !== window) {
                const msgType = isEdit ? 'character-updated' : 'character-created';
                const payload = isEdit
                    ? { type: msgType, character: merged, index: cs.currentEditingIndex }
                    : { type: msgType, character: merged };
                window.parent.postMessage(payload, '*');
            }
            closeCharacterCreationPopup();
            cs.currentEditingIndex = -1;
            loadCharacters();
        })
        .catch(err => {
            console.error('Failed to save character:', err);
            alert('Failed to save character. Please try again.');
        });
}

// ---------------------------------------------------------------------------
// User Details
// ---------------------------------------------------------------------------

function saveUserDetails() {
    const details = {
        name:        document.getElementById('user-name').value,
        description: document.getElementById('user-description').value,
        notes:       document.getElementById('user-notes').value
    };
    localStorage.setItem('howl-chat-user-details', JSON.stringify(details));
    closeUserDetailsPopup();
}

function loadUserDetails() {
    const saved = localStorage.getItem('howl-chat-user-details');
    if (!saved) return;
    const details = JSON.parse(saved);
    document.getElementById('user-name').value        = details.name        || '';
    document.getElementById('user-description').value = details.description || '';
    document.getElementById('user-notes').value       = details.notes       || '';
}

// ---------------------------------------------------------------------------
// Voice Settings UI Wiring
// ---------------------------------------------------------------------------

function _updateVoiceControlsVisible() {
    const enabled = document.getElementById('char-voice-enabled')?.checked || false;
    const controls = document.getElementById('char-voice-controls');
    if (controls) controls.style.display = enabled ? 'flex' : 'none';
}

function _updateSliderDisplay(id, valId, suffix) {
    const slider = document.getElementById(id);
    const valSpan = document.getElementById(valId);
    if (!slider || !valSpan) return;
    valSpan.textContent = slider.value + (suffix || '');
}

// ---------------------------------------------------------------------------
// Engine-aware voice catalog
// ---------------------------------------------------------------------------

let _voiceCatalogCache = new Map();
let _voiceCatalogRequest = 0;
let _activeVoiceCatalog = null;

function _globalAudioSettings() {
    return window.settingsState?.current?.audio || {};
}

function _voiceEngineLabel(engine) {
    const labels = {
        piper: 'Piper', kokoro: 'Kokoro', melotts: 'MeloTTS',
        browser: 'Browser TTS', auto: 'Auto-select'
    };
    return labels[engine] || engine || 'Inherit Global';
}

function _reportVoiceError(message) {
    if (typeof window.showNotification === 'function') {
        window.showNotification(message, 'error');
    } else {
        alert(message);
    }
}

function _emptyVoiceCatalog(engine, detail) {
    return {
        engine: engine || '',
        inherited: !engine,
        available: false,
        speed_supported: false,
        min_speed: 0.5,
        max_speed: 2.0,
        detail: detail,
        voices: []
    };
}

async function _getVoiceCatalog(engine, force) {
    const key = String(engine || '');
    if (!force && _voiceCatalogCache.has(key)) return _voiceCatalogCache.get(key);

    const app = window.go?.main?.App;
    if (!app?.GetTTSVoiceCatalog) {
        const missing = _emptyVoiceCatalog(key, 'The voice catalog binding is unavailable.');
        _voiceCatalogCache.set(key, missing);
        return missing;
    }
    try {
        const catalog = await app.GetTTSVoiceCatalog(key);
        if (!catalog || typeof catalog !== 'object') throw new Error('empty catalog response');
        _voiceCatalogCache.set(key, catalog);
        return catalog;
    } catch (err) {
        const failed = _emptyVoiceCatalog(key, 'Voice catalog unavailable: ' + (err?.message || err));
        _voiceCatalogCache.set(key, failed);
        return failed;
    }
}

function resetCharacterVoiceControls() {
    _activeVoiceCatalog = null;
    const voiceSelect = document.getElementById('char-voice-model');
    if (voiceSelect) {
        voiceSelect.innerHTML = '<option value="">Inherit Global</option>';
        voiceSelect.disabled = false;
    }
    const status = document.getElementById('char-voice-status');
    if (status) status.textContent = 'Enable a custom voice to load the available catalog.';
    const hint = document.getElementById('char-voice-speed-hint');
    if (hint) hint.textContent = '';
}

function _renderVoiceOptions(catalog, selectedVoice, inheritingEngine) {
    const select = document.getElementById('char-voice-model');
    if (!select) return;

    select.innerHTML = '';
    const inherited = document.createElement('option');
    inherited.value = '';
    inherited.textContent = inheritingEngine ? 'Inherit Global Voice' : 'Inherit Engine Default';
    select.appendChild(inherited);

    const voices = Array.isArray(catalog.voices) ? catalog.voices : [];
    voices.forEach(voice => {
        const option = document.createElement('option');
        option.value = voice.id;
        option.textContent = (voice.name || voice.id)
            + (voice.language ? ` · ${voice.language}` : '')
            + (voice.available ? '' : ' (unavailable)');
        option.disabled = !voice.available;
        option.title = voice.detail || '';
        select.appendChild(option);
    });

    const wanted = String(selectedVoice || '');
    if (wanted && !voices.some(voice => voice.id === wanted)) {
        const missing = document.createElement('option');
        missing.value = wanted;
        missing.textContent = wanted + ' (not in catalog)';
        missing.disabled = true;
        missing.title = 'The saved voice is not present in the selected engine catalog.';
        select.appendChild(missing);
    }
    select.value = wanted;
    select.disabled = voices.length === 0 && !wanted;
}

function _renderVoiceStatus(catalog, selectedVoice, inheritingEngine) {
    const status = document.getElementById('char-voice-status');
    if (!status) return;
    const label = _voiceEngineLabel(catalog.engine);
    const available = (Array.isArray(catalog.voices) ? catalog.voices : []).filter(v => v.available).length;
    const missing = (Array.isArray(catalog.voices) ? catalog.voices : []).filter(v => !v.available).length;
    let prefix = '';
    if (inheritingEngine) prefix = `Inheriting global engine (${label}). `;
    else prefix = `Using ${label}. `;

    let detail = catalog.detail || '';
    if (catalog.available) {
        detail = `${available} available voice${available === 1 ? '' : 's'}`
            + (missing ? `; ${missing} unavailable` : '') + '. ' + detail;
    } else if (!detail) {
        detail = 'No available voices were reported for this engine.';
    }
    const selected = selectedVoice
        ? (catalog.voices || []).find(voice => voice.id === selectedVoice)
        : null;
    if (selected && !selected.available) {
        detail = `Saved voice “${selectedVoice}” is unavailable. ` + (selected.detail || detail);
    }
    status.textContent = prefix + detail;
    status.style.color = catalog.available && (!selected || selected.available)
        ? 'var(--text-secondary, #888)' : '#f59e0b';
}

function _renderSpeedGuidance(catalog) {
    const hint = document.getElementById('char-voice-speed-hint');
    const slider = document.getElementById('char-voice-speed');
    if (!hint) return;
    if (!catalog.speed_supported) {
        hint.textContent = 'This engine does not expose a reliable speed control.';
        hint.style.color = 'var(--text-secondary, #888)';
        if (slider) slider.dataset.voiceSpeedSupported = 'false';
        return;
    }
    const min = Number(catalog.min_speed);
    const max = Number(catalog.max_speed);
    hint.textContent = `${_voiceEngineLabel(catalog.engine)} accepts ${min}–${max}x. Values outside this range are rejected when saving.`;
    if (slider) {
        slider.dataset.voiceSpeedSupported = 'true';
        slider.dataset.voiceMinSpeed = String(min);
        slider.dataset.voiceMaxSpeed = String(max);
    }
    const value = parseFloat(slider?.value || '1.0');
    const invalid = Number.isFinite(value) && (value < min || value > max);
    hint.style.color = invalid ? '#f59e0b' : 'var(--text-secondary, #888)';
    if (invalid) hint.textContent += ' Current value is outside this range.';
}

async function loadCharacterVoiceControls(selectedVoice, force) {
    const token = ++_voiceCatalogRequest;
    const engine = _readCharacterValue('char-voice-engine', '');
    const status = document.getElementById('char-voice-status');
    if (status) {
        status.textContent = 'Loading voice catalog...';
        status.style.color = 'var(--text-secondary, #888)';
    }
    const catalog = await _getVoiceCatalog(engine, force);
    if (token !== _voiceCatalogRequest) return catalog;

    _activeVoiceCatalog = catalog;
    _renderVoiceOptions(catalog, selectedVoice, !engine);
    _renderVoiceStatus(catalog, selectedVoice, !engine);
    _renderSpeedGuidance(catalog);
    return catalog;
}

async function validateCharacterVoice(patch) {
    if (!patch || patch.enabled === false) return true;

    const engine = patch.engine !== undefined
        ? patch.engine
        : _readCharacterValue('char-voice-engine', '');
    const voice = patch.voice !== undefined
        ? patch.voice
        : _readCharacterValue('char-voice-model', '');
    const speed = patch.speed !== undefined
        ? patch.speed
        : parseFloat(_readCharacterValue('char-voice-speed', '1.0') || '1.0');

    // Blank engine and voice is the complete inheritance case, but an explicit
    // speed still has to satisfy the inherited engine's range.
    const speedChanged = Object.prototype.hasOwnProperty.call(patch, 'speed');
    if (!engine && !voice && !speedChanged) return true;

    const catalog = await _getVoiceCatalog(engine, true);
    if ((engine || voice) && !catalog.available) {
        _reportVoiceError('Voice settings cannot be saved: ' + (catalog.detail || 'the selected engine has no available voices.'));
        return false;
    }
    if (voice) {
        const option = (catalog.voices || []).find(candidate => candidate.id === voice);
        if (!option || !option.available) {
            _reportVoiceError(`Voice “${voice}” is not available for ${_voiceEngineLabel(catalog.engine)}. Choose an available voice or clear the override.`);
            return false;
        }
    }
    if (catalog.speed_supported && Number.isFinite(speed)) {
        const min = Number(catalog.min_speed);
        const max = Number(catalog.max_speed);
        if (speed < min || speed > max) {
            _reportVoiceError(`${_voiceEngineLabel(catalog.engine)} accepts speed ${min}–${max}x, not ${speed}x.`);
            return false;
        }
    }
    return true;
}

document.addEventListener('DOMContentLoaded', () => {
    const voiceToggle = document.getElementById('char-voice-enabled');
    if (voiceToggle) {
        voiceToggle.addEventListener('change', () => {
            _updateVoiceControlsVisible();
            if (voiceToggle.checked) {
                loadCharacterVoiceControls(_readCharacterValue('char-voice-model', ''), false);
            }
        });
        _updateVoiceControlsVisible();
    }
    ['char-voice-enabled', 'char-voice-engine', 'char-voice-model', 'char-voice-speed', 'char-voice-volume'].forEach(id => {
        const el = document.getElementById(id);
        if (!el) return;
        const markDirty = () => window.characterState?.markFieldDirty?.(id);
        el.addEventListener('input', markDirty);
        el.addEventListener('change', markDirty);
    });
    const speedSlider = document.getElementById('char-voice-speed');
    if (speedSlider) {
        speedSlider.addEventListener('input', () => {
            _updateSliderDisplay('char-voice-speed', 'char-voice-speed-val', 'x');
            if (_activeVoiceCatalog) _renderSpeedGuidance(_activeVoiceCatalog);
        });
    }
    const volSlider = document.getElementById('char-voice-volume');
    if (volSlider) {
        volSlider.addEventListener('input', () => _updateSliderDisplay('char-voice-volume', 'char-voice-volume-val', ''));
    }

    const engineSelect = document.getElementById('char-voice-engine');
    if (engineSelect) {
        engineSelect.addEventListener('change', () => {
            const voiceSelect = document.getElementById('char-voice-model');
            if (voiceSelect) {
                voiceSelect.value = '';
                window.characterState?.markFieldDirty?.('char-voice-model');
            }
            loadCharacterVoiceControls('', true);
        });
    }
});
