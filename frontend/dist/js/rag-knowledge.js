/**
 * RAG Knowledge Panel
 *
 * Handles document upload, listing, and RAG toggle for the Story Tools page.
 *
 * @module rag-knowledge
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------

let _knowledgeDocs = [];
let _knowledgeScope = 'world';
let _knowledgeScopeID = '';

// ---------------------------------------------------------------------------
// Upload
// ---------------------------------------------------------------------------

async function handleKnowledgeUpload(input) {
    const file = input.files[0];
    if (!file) return;

    // Determine scope and scopeID from current context
    const scope = document.getElementById('knowledge-scope-select')?.value || 'world';
    const scopeID = resolveScopeID(scope);

    try {
        showLoading('Uploading document...');

        // The file input gives us content, not a path, so send the text itself.
        const text = await file.text();
        if (!text.trim()) throw new Error('That file has no readable text.');

        if (window.go?.main?.App?.UploadRAGDocumentText) {
            const docID = await window.go.main.App.UploadRAGDocumentText(
                file.name,
                scope,
                scopeID,
                'file',
                text
            );
            console.log('Document uploaded:', docID);
        } else {
            console.warn('UploadRAGDocumentText not available in Wails runtime');
        }

        await refreshKnowledgeList();
    } catch (err) {
        console.error('Failed to upload document:', err);
        alert('Upload failed: ' + err.message);
    } finally {
        hideLoading();
        input.value = '';
    }
}

// ---------------------------------------------------------------------------
// Document List
// ---------------------------------------------------------------------------

async function refreshKnowledgeList() {
    const scope = document.getElementById('knowledge-scope-select')?.value || 'world';
    const scopeID = resolveScopeID(scope);

    const listEl = document.getElementById('knowledge-list');
    if (!listEl) return;

    try {
        if (window.go?.main?.App?.ListRAGDocuments) {
            _knowledgeDocs = await window.go.main.App.ListRAGDocuments(scope, scopeID);
        } else {
            _knowledgeDocs = [];
        }
    } catch (err) {
        console.error('Failed to list documents:', err);
        _knowledgeDocs = [];
    }

    renderKnowledgeList();
}

function renderKnowledgeList() {
    const listEl = document.getElementById('knowledge-list');
    if (!listEl) return;

    if (!_knowledgeDocs || _knowledgeDocs.length === 0) {
        listEl.innerHTML = '<div class="knowledge-empty" style="color: var(--text-secondary); font-size: 13px; padding: 12px 0;">No documents uploaded yet</div>';
        return;
    }

    let html = '<div style="display: flex; flex-direction: column; gap: 8px;">';
    for (const doc of _knowledgeDocs) {
        const sizeKB = doc.SizeBytes ? (doc.SizeBytes / 1024).toFixed(1) : '?';
        const chunkCount = doc.ChunkCount || 0;
        html += `
            <div class="knowledge-item" style="display: flex; align-items: center; justify-content: space-between; padding: 8px 12px; background: rgba(0,0,0,0.2); border-radius: 6px; border: 1px solid var(--border-color, #333);">
                <div style="display: flex; flex-direction: column; min-width: 0;">
                    <span style="font-size: 14px; font-weight: 500; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">${escapeHtml(doc.Name)}</span>
                    <span style="font-size: 11px; color: var(--text-secondary);">${sizeKB} KB · ${chunkCount} chunks · ${doc.Scope}</span>
                </div>
                <button class="editor-btn editor-btn-danger" style="padding: 4px 10px; font-size: 12px;" onclick="deleteKnowledgeDocument('${doc.ID}')">Delete</button>
            </div>
        `;
    }
    html += '</div>';
    listEl.innerHTML = html;
}

// ---------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------

async function deleteKnowledgeDocument(docID) {
    if (!confirm('Delete this document from the knowledge base?')) return;

    try {
        showLoading('Deleting document...');
        if (window.go?.main?.App?.DeleteRAGDocument) {
            await window.go?.main?.App?.DeleteRAGDocument(docID);
        }
        await refreshKnowledgeList();
    } catch (err) {
        console.error('Failed to delete document:', err);
        alert('Delete failed: ' + err.message);
    } finally {
        hideLoading();
    }
}

async function toggleRAG(enabled) {
    try {
        if (window.go?.main?.App?.SetRAGEnabled) {
            await window.go.main.App.SetRAGEnabled(enabled);
        }
        console.log('RAG toggled:', enabled);
    } catch (err) {
        console.error('Failed to toggle RAG:', err);
    }
}

// ---------------------------------------------------------------------------
// Scope Resolution
// ---------------------------------------------------------------------------

function resolveScopeID(scope) {
    try {
        const currentState = JSON.parse(localStorage.getItem('howl-current-state') || '{}');
        switch (scope) {
            case 'world':
                return currentState.worldName || '';
            case 'scenario':
                return currentState.scenarioName || '';
            case 'chat':
                return currentState.chatName || '';
            default:
                return '';
        }
    } catch {
        return '';
    }
}

// ---------------------------------------------------------------------------
// Tab Integration
// ---------------------------------------------------------------------------

// switchMemoriesTab is owned by story-tools.js, which handles the knowledge
// panel itself. The wrapper that re-ran a captured copy is gone.

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function hideLoading() {
    const overlay = document.getElementById('loading-overlay');
    if (overlay) overlay.style.display = 'none';
}

// ---------------------------------------------------------------------------
// Init
// ---------------------------------------------------------------------------

document.addEventListener('DOMContentLoaded', function () {
    // Initialize RAG toggle from backend state
    if (window.go?.main?.App?.IsRAGEnabled) {
        window.go.main.App.IsRAGEnabled().then(enabled => {
            const toggle = document.getElementById('rag-enabled-toggle');
            if (toggle) toggle.checked = enabled;
        }).catch(() => {});
    }
});
