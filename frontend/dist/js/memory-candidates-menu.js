/**
 * Memory Candidate Sidebar Integration
 *
 * Adds "Review Memories" button to the sidebar with pending badge.
 * Hooks into chat message flow for auto-detection.
 *
 * @module memory-candidates-menu
 * @version 1.0.0
 * @requires memory-candidates.js
 */

'use strict';

// ---------------------------------------------------------------------------
// Sidebar Button Injection
// ---------------------------------------------------------------------------

function injectMemoryReviewButton() {
    const container = document.querySelector('.sidebar-tools') ||
                      document.querySelector('.sidebar-actions') ||
                      document.getElementById('sidebar-menu');
    if (!container) return;

    // Avoid duplicates
    if (container.querySelector('.memory-review-btn')) return;

    const btn = document.createElement('button');
    btn.className = 'memory-review-btn sidebar-btn';
    btn.title = 'Review Detected Memories';
    btn.innerHTML = `
        <span class="btn-icon">📝</span>
        <span class="btn-label">Review Memories</span>
        <span class="memory-candidate-badge" style="
            display: none; margin-left: auto; background: #f87171;
            color: #fff; font-size: 10px; font-weight: 700;
            padding: 2px 6px; border-radius: 10px; min-width: 18px;
            text-align: center;
        "></span>
    `;
    btn.style.cssText = `
        display: flex; align-items: center; gap: 8px;
        width: 100%; padding: 10px 14px;
        background: rgba(255,255,255,0.04);
        border: 1px solid rgba(255,255,255,0.06);
        border-radius: 10px; color: rgba(255,255,255,0.8);
        cursor: pointer; font-size: 14px; margin-top: 8px;
        transition: background 0.15s;
    `;
    btn.addEventListener('mouseenter', () => {
        btn.style.background = 'rgba(255,255,255,0.08)';
    });
    btn.addEventListener('mouseleave', () => {
        btn.style.background = 'rgba(255,255,255,0.04)';
    });
    btn.addEventListener('click', () => {
        if (window.MemoryCandidates) {
            window.MemoryCandidates.showReviewModal();
        }
    });

    container.appendChild(btn);

    // Initial badge refresh
    if (window.MemoryCandidates) {
        window.MemoryCandidates.refreshPendingCount();
    }
}

// ---------------------------------------------------------------------------
// Chat Message Hook
// ---------------------------------------------------------------------------

const ORIGINAL_SEND_MESSAGE = window.sendMessage;

function hookChatSend() {
    if (!window.sendMessage || window._memoryHookInstalled) return;

    window._memoryHookInstalled = true;
    const original = window.sendMessage.bind(window);

    window.sendMessage = async function(message, ...args) {
        // Call original first
        const result = await original(message, ...args);

        // Then run auto-detection in background
        try {
            await runAutoDetection(message);
        } catch (err) {
            console.warn('Auto-detection error:', err);
        }

        return result;
    };
}

async function runAutoDetection(message) {
    if (!message || typeof message !== 'string') return;

    // Try backend detection first
    if (window.go?.main?.App?.ProcessMessageForEvents) {
        try {
            const count = await window.go.main.App.ProcessMessageForEvents(message, [], []);
            if (count > 0 && window.MemoryCandidates) {
                window.MemoryCandidates.refreshPendingCount();
                showNotification(`${count} new memory candidate(s) detected`, 'info');
            }
            return;
        } catch (err) {
            console.warn('Backend detection failed, falling back:', err);
        }
    }

    // Frontend fallback: simple keyword detection
    const lower = message.toLowerCase();
    const keywords = {
        death: ['died', 'dead', 'killed', 'death', 'murdered'],
        betrayal: ['betrayed', 'traitor', 'deceived', 'backstabbed'],
        romance: ['kissed', 'love', 'married', 'lovers'],
        discovery: ['discovered', 'found', 'secret', 'uncovered'],
        conflict: ['attacked', 'battle', 'fought', 'war'],
        character_change: ['became', 'transformed', 'promoted', 'exiled'],
        location: ['arrived', 'entered', 'reached'],
        item: ['obtained', 'stole', 'acquired']
    };

    for (const [eventType, words] of Object.entries(keywords)) {
        const found = words.some(w => lower.includes(w));
        if (found && window.MemoryCandidates) {
            const cd = window.sidebarState?.currentChatData || {};
            const worldName = cd.world?.name || 'Unknown';
            const scenarioName = cd.scenario?.name || 'Unknown';

            window.MemoryCandidates.addLocalCandidate({
                id: 'cand_' + Date.now(),
                event_type: eventType,
                title: eventType.replace('_', ' ').replace(/\b\w/g, c => c.toUpperCase()),
                summary: `Detected ${eventType} event from message`,
                content: `Auto-captured ${eventType} event.\n\nWorld: ${worldName}\nScenario: ${scenarioName}`,
                source_message: message,
                scope: 'chat',
                owner_id: cd.sessionId || 'default',
                tags: [eventType, 'auto-captured'],
                confidence: 60,
                status: 'pending',
                created_at: new Date().toISOString()
            });
            break; // Only one candidate per message in fallback
        }
    }
}

// ---------------------------------------------------------------------------
// Initialization
// ---------------------------------------------------------------------------

function initMemoryCandidateMenu() {
    injectMemoryReviewButton();
    hookChatSend();
}

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initMemoryCandidateMenu);
} else {
    // Delay to ensure other scripts have set up
    setTimeout(initMemoryCandidateMenu, 500);
}
