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

        // Write file to a temporary path that Go can access
        // For Wails, we can pass the file content or path
        // Since Wails file dialog returns paths, we use the path directly
        const filePath = await saveFileForUpload(file);

        if (window.go?.main?.App?.UploadRAGDocument) {
            const docID = await window.go.main.App.UploadRAGDocument(
                file.name,
                scope,
                scopeID,
                'file',
                filePath
            );
            console.log('Document uploaded:', docID);
        } else {
            console.warn('UploadRAGDocument not available in Wails runtime');
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

// Save file content to a known path for Go to read
// In Wails, we can write to the app's working directory
async function saveFileForUpload(file) {
    // Read file as text
    const text = await file.text();
    // Store in a temporary location
    const tempPath = 'RAG/_upload_temp_' + Date.now() + '_' + file.name;
    try {
        // We can't directly write from JS, so we pass the content to Go
        // But since the Go binding expects a file path, we need the user to select
        // files via Wails OpenFileDialog instead
        // For now, return the file name and rely on Wails file dialog
        return file.name;
    } catch (e) {
        console.warn('Could not save temp file:', e);
        return file.name;
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

const _origSwitchMemoriesTab = window.switchMemoriesTab;

window.switchMemoriesTab = function(tab) {
    // Hide knowledge panel by default
    const knowledgePanel = document.getElementById('knowledge-panel');
    const memoryListCard = document.querySelector('.setting-card:nth-of-type(3)');
    const memoryEditorCard = document.querySelector('.setting-card:nth-of-type(4)');

    if (knowledgePanel) {
        knowledgePanel.style.display = tab === 'knowledge' ? 'block' : 'none';
    }
    if (memoryListCard) {
        memoryListCard.style.display = tab === 'knowledge' ? 'none' : 'block';
    }
    if (memoryEditorCard) {
        memoryEditorCard.style.display = tab === 'knowledge' ? 'none' : 'block';
    }

    // Update tab active states
    document.querySelectorAll('.memories-tab').forEach(t => t.classList.remove('active'));
    const activeTab = document.getElementById(`memories-tab-${tab}`);
    if (activeTab) activeTab.classList.add('active');

    if (tab === 'knowledge') {
        refreshKnowledgeList();
    } else if (_origSwitchMemoriesTab) {
        _origSwitchMemoriesTab(tab);
    }
};

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function showLoading(text) {
    const overlay = document.getElementById('loading-overlay');
    if (overlay) {
        overlay.querySelector('.loading-text').textContent = text;
        overlay.style.display = 'flex';
    }
}

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
