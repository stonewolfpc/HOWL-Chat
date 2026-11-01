/**
 * Menu — Edit Actions
 *
 * All Edit menu handler functions.
 *
 * @module menu-edit
 * @version 1.0.0
 */

'use strict';

// Undo/redo act on the focused text field; the browser keeps that field's edit history.
function _activeTextTarget() {
    const el = document.activeElement;
    if (!el) return null;
    if (el.isContentEditable) return el;
    if (el.tagName === 'TEXTAREA') return el;
    if (el.tagName === 'INPUT') {
        const t = (el.type || 'text').toLowerCase();
        if (['text', 'search', 'url', 'tel', 'password', 'email', 'number'].includes(t)) return el;
    }
    return null;
}

function editUndo() {
    const el = _activeTextTarget();
    if (!el) { showNotification('Click a text field first, then Undo', 'warning'); return; }
    el.focus();
    const done = document.execCommand('undo');
    showNotification(done ? 'Undo' : 'Nothing to undo', done ? 'success' : 'info');
}
function editRedo() {
    const el = _activeTextTarget();
    if (!el) { showNotification('Click a text field first, then Redo', 'warning'); return; }
    el.focus();
    const done = document.execCommand('redo');
    showNotification(done ? 'Redo' : 'Nothing to redo', done ? 'success' : 'info');
}

function editCut() {
    const sel = window.getSelection?.();
    if (sel?.toString()) { document.execCommand('cut');  showNotification('Text cut to clipboard', 'success'); }
    else                  { showNotification('No text selected', 'warning'); }
}

function editCopy() {
    const sel = window.getSelection?.();
    if (sel?.toString()) { document.execCommand('copy'); showNotification('Text copied to clipboard', 'success'); }
    else                  { showNotification('No text selected', 'warning'); }
}

function editPaste() {
    if (navigator.clipboard?.readText) {
        navigator.clipboard.readText().then(text => {
            const sel = window.getSelection?.();
            if (sel?.rangeCount > 0) {
                const range = sel.getRangeAt(0);
                range.deleteContents();
                range.insertNode(document.createTextNode(text));
                showNotification('Text pasted from clipboard', 'success');
            }
        }).catch(() => showNotification('Failed to paste from clipboard', 'error'));
    } else {
        document.execCommand('paste');
        showNotification('Paste attempted', 'info');
    }
}

function editSelectAll() {
    window.getSelection?.()?.selectAllChildren(document.body);
    showNotification('All text selected', 'success');
}

function editFind() {
    const term = prompt('Enter text to find:');
    if (!term) return;
    if (window.find) {
        showNotification(window.find(term) ? `Found: "${term}"` : `Not found: "${term}"`, window.find(term) ? 'success' : 'warning');
    } else {
        showNotification('Find functionality not supported in this browser', 'error');
    }
}

function editReplace() {
    const findTerm    = prompt('Enter text to find:');
    const replaceTerm = prompt('Enter replacement text:');
    if (!findTerm || !replaceTerm) return;
    const sel = window.getSelection?.();
    if (sel?.toString() === findTerm) {
        const range = sel.getRangeAt(0);
        range.deleteContents();
        range.insertNode(document.createTextNode(replaceTerm));
        showNotification('Text replaced', 'success');
    } else {
        showNotification('Select the text to replace first', 'warning');
    }
}
