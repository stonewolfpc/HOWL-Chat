/**
 * Menu — Edit Actions
 *
 * All Edit menu handler functions.
 *
 * @module menu-edit
 * @version 1.0.0
 */

'use strict';

function editUndo()      { showNotification('Undo functionality coming soon', 'info'); }
function editRedo()      { showNotification('Redo functionality coming soon', 'info'); }

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
