/**
 * HOWL Chat — Tutorial System UI
 *
 * 2-panel renderer: left navigation tree, right content.
 * Supports search, section expand/collapse, hash-based deep linking.
 *
 * Depends on: tutorials-data.js (window.TUTORIALS_DATA)
 *
 * @module tutorials-ui
 * @version 1.0.0
 */

'use strict';

// ---------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------

let _activeTopicId = null;
let _expandedSections = new Set();
let _searchQuery = '';

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

function openTutorials(topicId) {
    const overlay = document.getElementById('tutorials-overlay');
    if (!overlay) { console.warn('Tutorials overlay not found'); return; }

    overlay.classList.add('show');
    document.body.style.overflow = 'hidden';

    // Inject mobile nav toggle if missing
    const header = overlay.querySelector('.tutorials-header');
    if (header && !document.getElementById('tutorials-nav-toggle')) {
        const toggle = document.createElement('button');
        toggle.id = 'tutorials-nav-toggle';
        toggle.className = 'tutorials-nav-toggle';
        toggle.innerHTML = '<span class="icon icon-menu"></span>';
        toggle.title = 'Toggle navigation';
        header.insertBefore(toggle, header.children[1]);
    }

    if (topicId) {
        _activeTopicId = topicId;
        const sectionId = _findSectionForTopic(topicId);
        if (sectionId) _expandedSections.add(sectionId);
    } else if (!_activeTopicId) {
        const firstSection = window.TUTORIALS_DATA?.sections?.[0];
        if (firstSection) {
            _activeTopicId = firstSection.topics?.[0]?.id;
            _expandedSections.add(firstSection.id);
        }
    }

    _renderNav();
    _renderContent();
    _bindNavToggle();

    const searchInput = document.getElementById('tutorials-search-input');
    if (searchInput) {
        searchInput.value = '';
        searchInput.focus();
    }
}

function closeTutorials() {
    const overlay = document.getElementById('tutorials-overlay');
    if (overlay) overlay.classList.remove('show');
    document.body.style.overflow = '';
}

function navigateToTutorial(topicId) {
    if (!topicId || !window.TUTORIALS_DATA) return;
    _activeTopicId = topicId;
    const sectionId = _findSectionForTopic(topicId);
    if (sectionId) _expandedSections.add(sectionId);
    _renderNav();
    _renderContent();
}

// ---------------------------------------------------------------------------
// Navigation Renderer
// ---------------------------------------------------------------------------

function _renderNav() {
    const navEl = document.getElementById('tutorials-nav');
    if (!navEl || !window.TUTORIALS_DATA) return;

    if (_searchQuery) {
        _renderSearchResults(navEl);
        return;
    }

    navEl.innerHTML = window.TUTORIALS_DATA.sections.map(section => {
        const isExpanded = _expandedSections.has(section.id);
        const topicsHtml = section.topics.map(topic => {
            const isActive = topic.id === _activeTopicId;
            return `<div class="tutorials-nav-item ${isActive ? 'active' : ''}" data-topic-id="${topic.id}">${topic.title}</div>`;
        }).join('');

        return `
            <div class="tutorials-nav-section ${isExpanded ? 'expanded' : ''}" data-section-id="${section.id}">
                <div class="tutorials-nav-section-header" onclick="_toggleSection('${section.id}')">
                    <span>${section.title}</span>
                    <span class="icon icon-chevron-right"></span>
                </div>
                <div class="tutorials-nav-items">${topicsHtml}</div>
            </div>
        `;
    }).join('');

    _bindNavItems(navEl);
}

function _renderSearchResults(navEl) {
    const q = _searchQuery.toLowerCase();
    const results = [];
    window.TUTORIALS_DATA.sections.forEach(section => {
        section.topics.forEach(topic => {
            const titleMatch = topic.title.toLowerCase().includes(q);
            const contentMatch = topic.content.toLowerCase().includes(q);
            if (titleMatch || contentMatch) {
                results.push({ topic, section, titleMatch, contentMatch });
            }
        });
    });

    if (results.length === 0) {
        navEl.innerHTML = `<div class="tutorials-search-results"><div class="tutorials-search-result" style="color:rgba(212,212,212,0.4)">No results for "${_searchQuery}"</div></div>`;
        return;
    }

    navEl.innerHTML = `<div class="tutorials-search-results">
        <div class="tutorials-nav-section-header" style="padding:8px 16px;cursor:default">
            <span>${results.length} result${results.length !== 1 ? 's' : ''}</span>
        </div>
        ${results.map(r => `
            <div class="tutorials-search-result" data-topic-id="${r.topic.id}">
                ${r.titleMatch ? `<span class="result-match">${r.topic.title}</span>` : r.topic.title}
                <span style="color:rgba(212,212,212,0.35);font-size:11px"> — ${r.section.title}</span>
            </div>
        `).join('')}
    </div>`;

    navEl.querySelectorAll('.tutorials-search-result[data-topic-id]').forEach(el => {
        el.addEventListener('click', () => {
            _searchQuery = '';
            document.getElementById('tutorials-search-input').value = '';
            navigateToTutorial(el.dataset.topicId);
            if (window.innerWidth <= 560) {
                navEl.classList.remove('show');
            }
        });
    });
}

function _bindNavItems(navEl) {
    navEl.querySelectorAll('.tutorials-nav-item[data-topic-id]').forEach(el => {
        el.addEventListener('click', () => {
            navigateToTutorial(el.dataset.topicId);
            // Close mobile nav panel after selection
            if (window.innerWidth <= 560) {
                navEl.classList.remove('show');
            }
        });
    });
}

function _bindNavToggle() {
    const toggle = document.getElementById('tutorials-nav-toggle');
    const nav = document.getElementById('tutorials-nav');
    if (toggle && nav) {
        toggle.onclick = () => nav.classList.toggle('show');
    }
}

function _toggleSection(sectionId) {
    if (_expandedSections.has(sectionId)) {
        _expandedSections.delete(sectionId);
    } else {
        _expandedSections.add(sectionId);
    }
    _renderNav();
}

// ---------------------------------------------------------------------------
// Content Renderer
// ---------------------------------------------------------------------------

function _renderContent() {
    const contentEl = document.getElementById('tutorials-content');
    if (!contentEl || !window.TUTORIALS_DATA) return;

    const topic = _findTopic(_activeTopicId);
    if (!topic) {
        contentEl.innerHTML = `<div style="text-align:center;padding:60px 20px;color:rgba(212,212,212,0.4)"><div style="font-size:48px;margin-bottom:16px" class="icon icon-book"></div><p>Select a topic from the left panel</p></div>`;
        return;
    }

    const section = _findSectionForTopic(_activeTopicId);
    const breadcrumb = section
        ? `<div class="tutorials-breadcrumb"><span>${section.title}</span><span style="color:rgba(212,212,212,0.3)">/</span><span class="active">${topic.title}</span></div>`
        : '';

    contentEl.innerHTML = breadcrumb + _parseContent(topic.content);

    // Scroll to top
    contentEl.scrollTop = 0;
}

function _parseContent(content) {
    // Simple markdown-like parser for our tutorial format
    let html = content;

    // Escape HTML entities first
    html = html.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');

    // Un-escape our custom tags
    html = html.replace(/&lt;info&gt;(.*?)&lt;\/info&gt;/gs, '<div class="tutorials-info-box"><div class="info-title"><span class="icon icon-info"></span>Note</div><p>$1</p></div>');
    html = html.replace(/&lt;warning&gt;(.*?)&lt;\/warning&gt;/gs, '<div class="tutorials-warning-box"><div class="warning-title"><span class="icon icon-alert"></span>Important</div><p>$1</p></div>');

    // Code blocks
    html = html.replace(/```(\w+)?\n(.*?)```/gs, (match, lang, code) => {
        return `<pre><code>${code.trim()}</code></pre>`;
    });

    // Inline code (after code blocks to avoid double-processing)
    html = html.replace(/`([^`]+)`/g, '<code>$1</code>');

    // Headers
    html = html.replace(/^### (.*$)/gm, '<h3>$1</h3>');
    html = html.replace(/^## (.*$)/gm, '<h2>$1</h2>');
    html = html.replace(/^# (.*$)/gm, '<h1>$1</h1>');

    // Bold and italic
    html = html.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');
    html = html.replace(/\*(.*?)\*/g, '<em>$1</em>');

    // Links [text](url)
    html = html.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank" rel="noopener">$1</a>');

    // Keyboard keys <key>Name</key>
    html = html.replace(/&lt;key&gt;(.*?)&lt;\/key&gt;/g, '<span class="tutorials-key">$1</span>');

    // Tables
    html = html.replace(/\|(.+)\|\n\|[-:\s|]+\|\n((?:\|.+)\n?)+/g, (match) => {
        const lines = match.trim().split('\n').filter(l => l.trim());
        if (lines.length < 2) return match;
        let tableHtml = '<table class="tutorials-table">';
        lines.forEach((line, i) => {
            const cells = line.split('|').map(c => c.trim()).filter(c => c);
            const tag = i === 0 ? 'th' : 'td';
            tableHtml += '<tr>' + cells.map(c => `<${tag}>${c}</${tag}>`).join('') + '</tr>';
        });
        tableHtml += '</table>';
        return tableHtml;
    });

    // Lists
    html = html.replace(/(?:^|\n)((?:- .+\n?)+)/g, (match, listBlock) => {
        const items = listBlock.trim().split('\n').map(l => l.replace(/^- /, '')).filter(l => l);
        return '<ul>' + items.map(item => `<li>${item}</li>`).join('') + '</ul>';
    });

    // Numbered lists
    html = html.replace(/(?:^|\n)((?:\d+\. .+\n?)+)/g, (match, listBlock) => {
        const items = listBlock.trim().split('\n').map(l => l.replace(/^\d+\. /, '')).filter(l => l);
        return '<ol>' + items.map(item => `<li>${item}</li>`).join('') + '</ol>';
    });

    // Paragraphs (split on double newlines)
    const blocks = html.split(/\n\n+/);
    html = blocks.map(block => {
        const trimmed = block.trim();
        if (!trimmed) return '';
        // Skip if already wrapped in a tag
        if (trimmed.startsWith('<h') || trimmed.startsWith('<ul') || trimmed.startsWith('<ol') ||
            trimmed.startsWith('<table') || trimmed.startsWith('<div') || trimmed.startsWith('<pre')) {
            return trimmed;
        }
        return '<p>' + trimmed.replace(/\n/g, ' ') + '</p>';
    }).join('\n');

    return html;
}

// ---------------------------------------------------------------------------
// Search
// ---------------------------------------------------------------------------

function _initSearch() {
    const searchInput = document.getElementById('tutorials-search-input');
    if (!searchInput) return;

    let debounceTimer;
    searchInput.addEventListener('input', (e) => {
        clearTimeout(debounceTimer);
        debounceTimer = setTimeout(() => {
            _searchQuery = e.target.value.trim();
            _renderNav();
        }, 200);
    });

    searchInput.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') {
            _searchQuery = '';
            searchInput.value = '';
            _renderNav();
        }
    });
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function _findTopic(topicId) {
    if (!window.TUTORIALS_DATA) return null;
    for (const section of window.TUTORIALS_DATA.sections) {
        const topic = section.topics.find(t => t.id === topicId);
        if (topic) return topic;
    }
    return null;
}

function _findSectionForTopic(topicId) {
    if (!window.TUTORIALS_DATA) return null;
    for (const section of window.TUTORIALS_DATA.sections) {
        if (section.topics.some(t => t.id === topicId)) return section;
    }
    return null;
}

// ---------------------------------------------------------------------------
// Init
// ---------------------------------------------------------------------------

document.addEventListener('DOMContentLoaded', () => {
    _initSearch();

    // Close on overlay click (but not container click)
    const overlay = document.getElementById('tutorials-overlay');
    if (overlay) {
        overlay.addEventListener('click', (e) => {
            if (e.target === overlay) {
                // On mobile, close nav panel first if open
                if (window.innerWidth <= 560) {
                    const nav = document.getElementById('tutorials-nav');
                    if (nav?.classList.contains('show')) {
                        nav.classList.remove('show');
                        return;
                    }
                }
                closeTutorials();
            }
        });
    }

    // Close on Escape
    document.addEventListener('keydown', (e) => {
        const overlayEl = document.getElementById('tutorials-overlay');
        if (e.key === 'Escape' && overlayEl?.classList.contains('show')) {
            if (window.innerWidth <= 560) {
                const nav = document.getElementById('tutorials-nav');
                if (nav?.classList.contains('show')) {
                    nav.classList.remove('show');
                    return;
                }
            }
            closeTutorials();
        }
    });

    // Clicking content area on mobile closes the nav panel
    const contentEl = document.getElementById('tutorials-content');
    if (contentEl) {
        contentEl.addEventListener('click', () => {
            if (window.innerWidth <= 560) {
                const nav = document.getElementById('tutorials-nav');
                if (nav) nav.classList.remove('show');
            }
        });
    }
});

// ---------------------------------------------------------------------------
// Exports
// ---------------------------------------------------------------------------

window.TutorialsUI = {
    open: openTutorials,
    close: closeTutorials,
    navigate: navigateToTutorial
};
