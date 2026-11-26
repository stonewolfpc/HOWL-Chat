'use strict';

const test = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');

class FakeClassList {
    constructor() {
        this.values = new Set();
    }
    contains(name) {
        return this.values.has(name);
    }
}

class FakeElement {
    constructor(tagName) {
        this.tagName = tagName;
        this.dataset = {};
        this.classList = new FakeClassList();
        this.innerHTML = '';
        this.children = [];
        this.scrollTop = 0;
        this.scrollHeight = 0;
    }
    appendChild(child) {
        this.children.push(child);
    }
}

function loadRenderer(characters) {
    const container = new FakeElement('div');
    const sandbox = {
        console,
        document: {
            querySelector(selector) {
                return selector === '.chat-messages' ? container : null;
            },
            createElement(tagName) {
                return new FakeElement(tagName);
            },
        },
        currentChatData: { characters },
        ICONS: {
            regenerate: 'regenerate',
            previous: 'previous',
            next: 'next',
            continue: 'continue',
            edit: 'edit',
            copy: 'copy',
            undo: 'undo',
            delete: 'delete',
            speaker: 'speaker',
        },
        updateNavigationButtons() {},
    };
    sandbox.globalThis = sandbox;
    vm.createContext(sandbox);
    const source = fs.readFileSync(
        path.join(__dirname, '..', '..', 'frontend', 'dist', 'js', 'chat-render.js'),
        'utf8'
    );
    vm.runInContext(source, sandbox, { filename: 'chat-render.js' });
    return { sandbox, container };
}

test('explicit character identity selects the owning character', () => {
    const { sandbox, container } = loadRenderer([
        { id: 'char_aria', name: 'Aria' },
        { id: 'char_doren', name: 'Doren' },
    ]);

    sandbox.addAssistantMessage('Aria', 'First line', 'char_doren');

    const message = container.children[0];
    assert.equal(message.dataset.characterId, 'char_doren');
    assert.match(message.innerHTML, /Doren/);
    assert.doesNotMatch(message.innerHTML, />Aria</);
});

test('missing character identity does not fall back to the first character', () => {
    const { sandbox, container } = loadRenderer([
        { id: 'char_aria', name: 'Aria' },
        { id: 'char_doren', name: 'Doren' },
    ]);

    sandbox.addAssistantMessage('Legacy sender', 'Old line');

    const message = container.children[0];
    assert.equal(message.dataset.characterId, '');
    assert.match(message.innerHTML, /Legacy sender/);
    assert.doesNotMatch(message.innerHTML, /Aria/);
});

test('unknown character identity is preserved and does not select another character', () => {
    const { sandbox, container } = loadRenderer([
        { id: 'char_aria', name: 'Aria' },
        { id: 'char_doren', name: 'Doren' },
    ]);

    sandbox.addAssistantMessage('Stored sender', 'Old line', 'char_missing');

    const message = container.children[0];
    assert.equal(message.dataset.characterId, 'char_missing');
    assert.match(message.innerHTML, /Stored sender/);
    assert.doesNotMatch(message.innerHTML, /Aria/);
});
