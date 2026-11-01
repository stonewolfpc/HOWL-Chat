/**
 * Chat Thinking — reasoning display and settings
 *
 * Thinking models reason before they answer. Three modes:
 *   off    — the server is told not to think at all
 *   on     — reasoning streams into a collapsible spoilered block above the reply
 *   hidden — reasoning is used but never rendered
 *
 * In every mode the thought never becomes the reply and never enters history;
 * the backend keeps it out of the transcript and off the next prompt. This file
 * only decides whether a human sees it.
 *
 * @module chat-think
 */

'use strict';

let thinkText = '';
let thinkStart = null;

function getThinkMode() {
    return (window.settingsState?.current?.model?.think_mode) || 'hidden';
}

// The streaming assistant bubble is the only place reasoning may be shown.
function ensureThinkArea() {
    const bubble = document.querySelector('.chat-message.assistant.streaming .chat-bubble');
    if (!bubble) return null;
    let area = bubble.querySelector('.chat-think');
    if (!area) {
        area = document.createElement('div');
        area.className = 'chat-think';
        bubble.insertBefore(area, bubble.firstChild);
    }
    return area;
}

function onReasoningChunk(chunk) {
    if (!chunk || getThinkMode() !== 'on') return;
    if (!thinkStart) thinkStart = Date.now();
    thinkText += chunk;
    const area = ensureThinkArea();
    if (area && typeof updateThinkArea === 'function') {
        const escaped = thinkText.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
        updateThinkArea(area, escaped, thinkStart);
    }
}

function resetThink() {
    thinkText = '';
    thinkStart = null;
}

function initThinkSettings() {
    const sel = document.getElementById('think-mode');
    const budget = document.getElementById('think-budget');
    if (!sel) return;

    const modelSettings = () => {
        if (!window.settingsState?.current) return null;
        window.settingsState.current.model = window.settingsState.current.model || {};
        return window.settingsState.current.model;
    };

    const current = modelSettings();
    if (current) {
        sel.value = current.think_mode || 'hidden';
        if (budget && current.think_budget) budget.value = current.think_budget;
    }

    const save = () => {
        const m = modelSettings();
        if (!m) return;
        m.think_mode = sel.value;
        if (budget) {
            const n = parseInt(budget.value, 10);
            m.think_budget = Number.isFinite(n) && n >= 0 ? n : 2048;
        }
        if (window.settingsState?.saveSettings) window.settingsState.saveSettings();
    };

    sel.addEventListener('change', save);
    if (budget) budget.addEventListener('change', save);
}

document.addEventListener('DOMContentLoaded', () => {
    initThinkSettings();
    if (window.runtime?.EventsOn) {
        window.runtime.EventsOn('chat:reasoning', onReasoningChunk);
        window.runtime.EventsOn('chat:start', resetThink);
        window.runtime.EventsOn('chat:complete', resetThink);
        window.runtime.EventsOn('chat:aborted', resetThink);
    }
});
