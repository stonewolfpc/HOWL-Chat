/**
 * Memory Candidate Review UI
 *
 * Displays detected memory events awaiting user approval.
 * Allows approve/reject of auto-captured candidates.
 *
 * @module memory-candidates
 * @version 1.0.0
 * @requires memory-store.js
 */

'use strict';

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

window.MemoryCandidates = {
    showReviewModal,
    refreshPendingCount,
    getPendingCount,
    addLocalCandidate
};

// ---------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------

const CANDIDATES_KEY = 'howl-chat-candidates';

// ---------------------------------------------------------------------------
// UI Rendering
// ---------------------------------------------------------------------------

function showReviewModal() {
    const candidates = getPendingCandidates();
    renderReviewModal(candidates);
}

function renderReviewModal(candidates) {
    const existing = document.querySelector('.memory-candidates-modal');
    if (existing) existing.remove();

    const modal = document.createElement('div');
    modal.className = 'memory-candidates-modal modal-overlay';
    modal.style.cssText = `
        position: fixed; top: 0; left: 0; right: 0; bottom: 0;
        background: rgba(0,0,0,0.8); display: flex;
        align-items: center; justify-content: center;
        z-index: 10000; backdrop-filter: blur(4px);
    `;

    const approved = getCandidatesByStatus('approved');
    const rejected = getCandidatesByStatus('rejected');

    modal.innerHTML = `
        <div class="modal-container" style="
            background: var(--shell-panel-bg); border-radius: 16px;
            width: min(800px, 95%); max-height: 85vh;
            display: flex; flex-direction: column;
            border: 1px solid rgba(255,255,255,0.12);
        ">
            <div class="modal-header" style="
                padding: 20px 24px; border-bottom: 1px solid rgba(255,255,255,0.08);
                display: flex; align-items: center; justify-content: space-between;
            ">
                <div>
                    <div style="font-size: 20px; font-weight: 700; color: #fff;">
                        Memory Review Queue
                    </div>
                    <div style="font-size: 13px; color: var(--surface-hover); margin-top: 4px;">
                        ${candidates.length} pending •
                        ${approved.length} approved •
                        ${rejected.length} rejected
                    </div>
                </div>
                <button onclick="closeReviewModal()" style="
                    background: var(--surface-hover); border: none; color: #fff;
                    width: 36px; height: 36px; border-radius: 10px;
                    cursor: pointer; font-size: 18px;
                ">✕</button>
            </div>

            <div class="modal-body" style="
                padding: 20px 24px; overflow-y: auto; flex: 1;
            ">
                ${candidates.length === 0 ? renderEmptyQueue() : candidates.map(c => renderCandidateCard(c)).join('')}
            </div>
        </div>
    `;

    document.body.appendChild(modal);
}

function renderCandidateCard(candidate) {
    const eventColor = getEventColor(candidate.event_type);
    const confidenceBar = Math.round(candidate.confidence);

    return `
        <div class="candidate-card" data-id="${candidate.id}" style="
            background: var(--surface-hover); border-radius: 12px;
            padding: 16px; margin-bottom: 12px;
            border: 1px solid rgba(255,255,255,0.06);
        ">
            <div style="display: flex; align-items: flex-start; gap: 12px;">
                <div style="
                    width: 40px; height: 40px; border-radius: 10px;
                    background: ${eventColor}22; color: ${eventColor};
                    display: flex; align-items: center; justify-content: center;
                    font-size: 18px; flex-shrink: 0;
                ">${getEventIcon(candidate.event_type)}</div>

                <div style="flex: 1; min-width: 0;">
                    <div style="display: flex; align-items: center; gap: 8px; margin-bottom: 4px;">
                        <span style="font-weight: 600; color: #fff; font-size: 15px;">${escapeHtml(candidate.title)}</span>
                        <span style="
                            font-size: 11px; padding: 2px 8px; border-radius: 6px;
                            background: ${eventColor}22; color: ${eventColor};
                            text-transform: uppercase; font-weight: 600;
                        ">${candidate.event_type}</span>
                    </div>
                    <div style="color: var(--text-secondary); font-size: 13px; margin-bottom: 8px;">
                        ${escapeHtml(candidate.summary)}
                    </div>
                    <div style="
                        background: rgba(0,0,0,0.3); border-radius: 6px;
                        padding: 8px 12px; font-size: 12px;
                        color: rgba(255,255,255,0.5); margin-bottom: 10px;
                        max-height: 80px; overflow-y: auto;
                    ">
                        <div style="color: var(--text-tertiary); margin-bottom: 4px; font-size: 11px;">Source:</div>
                        ${escapeHtml(candidate.source_message || candidate.content).substring(0, 200)}
                    </div>
                    <div style="display: flex; align-items: center; gap: 8px; margin-bottom: 10px;">
                        <div style="flex: 1; height: 4px; background: var(--surface-hover); border-radius: 2px; overflow: hidden;">
                            <div style="width: ${confidenceBar}%; height: 100%; background: ${confidenceBar >= 75 ? '#4ade80' : confidenceBar >= 50 ? '#fbbf24' : '#f87171'}; border-radius: 2px;"></div>
                        </div>
                        <span style="font-size: 11px; color: var(--surface-hover);">${confidenceBar}%</span>
                    </div>
                    <div style="display: flex; gap: 8px;">
                        <button onclick="MemoryCandidates.approve('${candidate.id}')" style="
                            background: rgba(74,222,128,0.15); border: 1px solid rgba(74,222,128,0.3);
                            color: #4ade80; padding: 6px 16px; border-radius: 8px;
                            cursor: pointer; font-size: 13px; font-weight: 500;
                        ">Approve</button>
                        <button onclick="MemoryCandidates.reject('${candidate.id}')" style="
                            background: rgba(248,113,113,0.15); border: 1px solid rgba(248,113,113,0.3);
                            color: #f87171; padding: 6px 16px; border-radius: 8px;
                            cursor: pointer; font-size: 13px; font-weight: 500;
                        ">Reject</button>
                    </div>
                </div>
            </div>
        </div>
    `;
}

function renderEmptyQueue() {
    return `
        <div style="text-align: center; padding: 60px 20px; color: var(--surface-hover);">
            <div style="font-size: 48px; margin-bottom: 16px;">📝</div>
            <div style="font-size: 16px; font-weight: 600; color: var(--surface-hover); margin-bottom: 8px;">
                No Pending Memories
            </div>
            <div style="font-size: 13px;">
                Auto-captured events will appear here for your review.<br>
                Chat naturally and the system will detect notable moments.
            </div>
        </div>
    `;
}

// ---------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------

window.MemoryCandidates.approve = async function(candidateId) {
    const card = document.querySelector(`.candidate-card[data-id="${candidateId}"]`);
    if (card) {
        card.style.opacity = '0.5';
        card.style.pointerEvents = 'none';
    }

    try {
        // Try backend first
        if (window.go?.main?.App?.ApproveCandidate) {
            await window.go.main.App.ApproveCandidate(candidateId);
        } else {
            // Fallback: convert localStorage candidate to memory
            approveLocalCandidate(candidateId);
        }
        showNotification('Memory approved', 'success');
        refreshModal();
        refreshPendingCount();
    } catch (err) {
        console.error('Approve failed:', err);
        showNotification('Failed to approve memory', 'error');
        if (card) {
            card.style.opacity = '1';
            card.style.pointerEvents = '';
        }
    }
};

window.MemoryCandidates.reject = async function(candidateId) {
    const card = document.querySelector(`.candidate-card[data-id="${candidateId}"]`);
    if (card) card.style.opacity = '0.3';

    try {
        if (window.go?.main?.App?.RejectCandidate) {
            await window.go.main.App.RejectCandidate(candidateId);
        } else {
            rejectLocalCandidate(candidateId);
        }
        showNotification('Memory rejected', 'info');
        refreshModal();
        refreshPendingCount();
    } catch (err) {
        console.error('Reject failed:', err);
        showNotification('Failed to reject memory', 'error');
        if (card) card.style.opacity = '1';
    }
};

function closeReviewModal() {
    const modal = document.querySelector('.memory-candidates-modal');
    if (modal) modal.remove();
}

function refreshModal() {
    const existing = document.querySelector('.memory-candidates-modal');
    if (existing) {
        showReviewModal();
    }
}

// ---------------------------------------------------------------------------
// Data Layer (localStorage fallback)
// ---------------------------------------------------------------------------

function getPendingCandidates() {
    return getCandidatesByStatus('pending');
}

function getCandidatesByStatus(status) {
    const all = loadCandidates();
    return all.filter(c => c.status === status);
}

function loadCandidates() {
    try {
        const raw = localStorage.getItem(CANDIDATES_KEY);
        return raw ? JSON.parse(raw) : [];
    } catch {
        return [];
    }
}

function saveCandidates(candidates) {
    localStorage.setItem(CANDIDATES_KEY, JSON.stringify(candidates));
}

function addLocalCandidate(candidate) {
    const all = loadCandidates();
    // Avoid duplicates by source message
    const exists = all.some(c => c.source_message === candidate.source_message);
    if (!exists) {
        all.push(candidate);
        saveCandidates(all);
        refreshPendingCount();
        showNotification('New memory detected - review pending', 'info');
    }
}

function approveLocalCandidate(candidateId) {
    const all = loadCandidates();
    const idx = all.findIndex(c => c.id === candidateId);
    if (idx === -1) return;

    const candidate = all[idx];
    candidate.status = 'approved';
    candidate.reviewed_at = new Date().toISOString();

    // Convert to memory entry
    if (window.MemoryStore) {
        window.MemoryStore.createMemoryEntry({
            title: candidate.title,
            summary: candidate.summary,
            content: candidate.content,
            scope: candidate.scope,
            ownerId: candidate.owner_id,
            tags: candidate.tags,
            source: 'auto',
            confidence: candidate.confidence
        });
    }

    saveCandidates(all);
}

function rejectLocalCandidate(candidateId) {
    const all = loadCandidates();
    const idx = all.findIndex(c => c.id === candidateId);
    if (idx !== -1) {
        all[idx].status = 'rejected';
        all[idx].reviewed_at = new Date().toISOString();
        saveCandidates(all);
    }
}

// ---------------------------------------------------------------------------
// Badge / Notifications
// ---------------------------------------------------------------------------

function getPendingCount() {
    return getPendingCandidates().length;
}

function refreshPendingCount() {
    const count = getPendingCount();
    updateBadge(count);
    return count;
}

function updateBadge(count) {
    // Update any existing badge elements
    document.querySelectorAll('.memory-candidate-badge').forEach(el => {
        el.textContent = count > 0 ? String(count) : '';
        el.style.display = count > 0 ? 'flex' : 'none';
    });
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function getEventColor(eventType) {
    const colors = {
        // Big events
        death: '#f87171', betrayal: '#fb923c', romance: '#f472b6',
        discovery: '#60a5fa', conflict: '#fbbf24',
        character_change: '#a78bfa', location: '#34d399',
        item: '#facc15', custom: '#9ca3af',
        // Micro-memories
        preference: '#818cf8', like_dislike: '#f472b6',
        emotion: '#fbbf24', association: '#a78bfa',
        aesthetic: '#34d399', sensory: '#60a5fa',
        social: '#fb923c', first: '#f87171',
        shared: '#818cf8', promise: '#f472b6'
    };
    return colors[eventType] || '#9ca3af';
}

function getEventIcon(eventType) {
    const icons = {
        // Big events
        death: '💀', betrayal: '🗡️', romance: '💕',
        discovery: '🔍', conflict: '⚔️',
        character_change: '🎭', location: '🏰',
        item: '💎', custom: '📌',
        // Micro-memories
        preference: '⭐', like_dislike: '👍',
        emotion: '💭', association: '🔗',
        aesthetic: '🎨', sensory: '👃',
        social: '🤝', first: '🌟',
        shared: '📤', promise: '🤞'
    };
    return icons[eventType] || '📌';
}

// Auto-refresh badge on load
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', refreshPendingCount);
} else {
    refreshPendingCount();
}
