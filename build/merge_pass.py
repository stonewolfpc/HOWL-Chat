"""T2.2 final sub-pass: collapse the remaining shadowed names onto one owner each.

Every edit asserts its anchor exists and that the result still parses, because a
silent no-op patch is exactly what let the last one survive a "successful" run.
"""
import io
import re
import subprocess

def rd(p):
    return io.open(p, encoding="utf-8", newline="").read()

def wr(p, s):
    io.open(p, "w", encoding="utf-8", newline="").write(s)

def must(cond, msg):
    if not cond:
        raise SystemExit("FAILED: " + msg)

# ---------------------------------------------------------------------------
# 1. sidebar-context.js — drop the loadWorlds stub, make the API lazy
# ---------------------------------------------------------------------------
p = "frontend/dist/js/sidebar-context.js"
s = rd(p)

stub = """function loadWorlds() {
    const worlds = JSON.parse(localStorage.getItem('howl-worlds') || '[]');
    console.log('Loaded worlds:', worlds);
}
"""
must(stub in s, "sidebar-context.js: loadWorlds stub not found")
s = s.replace(stub, "", 1)

must("    loadWorlds();\n" in s, "sidebar-context.js: loadWorlds() call not found")
s = s.replace("    loadWorlds();\n",
              "    // Worlds are loaded by world-grid.js (it rehydrates images and renders the\n"
              "    // grid). The log-only copy that lived here shadowed it.\n", 1)

api_start = s.index("window.sidebarContext = {")
api_end = s.index("};", api_start) + 2
s = s[:api_start] + """// Every panel entry is resolved at call time. This module is loaded before
// sidebar-panels.js, and on story-tools.html it loads without context-display.js
// at all, so a direct reference here was a ReferenceError that left
// window.sidebarContext undefined on that page.
function _callGlobal(name) {
    return function (...args) {
        const fn = window[name];
        return typeof fn === 'function' ? fn.apply(window, args) : undefined;
    };
}

window.sidebarContext = {
    initializeContextPanel,
    loadContextData,
    updateContextDisplay:  _callGlobal('updateContextDisplay'),
    updateWorldDisplay:    _callGlobal('updateWorldDisplay'),
    updateScenarioDisplay: _callGlobal('updateScenarioDisplay'),
    updateLorebookDisplay: _callGlobal('updateLorebookDisplay'),
    openWorldSelector:     _callGlobal('openWorldSelector'),
    selectWorld:           _callGlobal('selectWorld'),
    toggleWorldAutoSwitch: _callGlobal('toggleWorldAutoSwitch'),
    openScenarioSelector:  _callGlobal('openScenarioSelector'),
    filterScenarios:       _callGlobal('filterScenarios'),
    selectScenario:        _callGlobal('selectScenario'),
    toggleScenarioAutoSwitch: _callGlobal('toggleScenarioAutoSwitch'),
    createLorebookItem:    _callGlobal('createLorebookItem'),
    removeLorebookFromChat: _callGlobal('removeLorebookFromChat'),
    openLorebookSelector:  _callGlobal('openLorebookSelector'),
    filterLorebooks:       _callGlobal('filterLorebooks'),
    toggleChatLorebook:    _callGlobal('toggleChatLorebook'),
    openLorebookManager:   _callGlobal('openLorebookManager')
};""" + s[api_end:]
wr(p, s)
print("sidebar-context.js: stub removed, API resolved lazily")

# ---------------------------------------------------------------------------
# 2. sidebar-context-lorebooks.js — rename to describe what it actually toggles
#    (chat-level lorebooks by id, not a dropdown row)
# ---------------------------------------------------------------------------
p = "frontend/dist/js/sidebar-context-lorebooks.js"
s = rd(p)
must("async function toggleLorebookSelection(lorebookId) {" in s,
     "sidebar-context-lorebooks.js: function not found")
s = s.replace("async function toggleLorebookSelection(lorebookId) {",
              "// Chat-level toggle (by id). The world/scenario dropdowns have their own\n"
              "// toggleLorebookSelection(lorebook, element); same name, different job,\n"
              "// which is why this one is named for what it does.\n"
              "async function toggleChatLorebook(lorebookId) {", 1)
s = s.replace('onclick="toggleLorebookSelection(', 'onclick="toggleChatLorebook(')
wr(p, s)
must("toggleLorebookSelection" not in s, "sidebar-context-lorebooks.js: old name still present")
print("sidebar-context-lorebooks.js: renamed to toggleChatLorebook")

# ---------------------------------------------------------------------------
# 3. story-tools.js — own the memories tab switch, including the knowledge panel
# ---------------------------------------------------------------------------
p = "frontend/dist/js/story-tools.js"
s = rd(p)
old_tab = """function switchMemoriesTab(tab) {
    window.storyToolsState.currentTab = tab;
    document.querySelectorAll('.memories-tab').forEach(t => t.classList.remove('active'));
    document.getElementById(`memories-tab-${tab}`)?.classList.add('active');
    loadFullMemoriesList();
}"""
must(old_tab in s, "story-tools.js: switchMemoriesTab not found")
new_tab = """function switchMemoriesTab(tab) {
    window.storyToolsState.currentTab = tab;

    // The knowledge tab swaps the memory list and editor for the RAG panel.
    // This used to be a wrapper in rag-knowledge.js that captured an "original"
    // and re-ran it conditionally; one owner is enough.
    const knowledgePanel   = document.getElementById('knowledge-panel');
    const memoryListCard   = document.querySelector('.setting-card:nth-of-type(3)');
    const memoryEditorCard = document.querySelector('.setting-card:nth-of-type(4)');

    if (knowledgePanel)   knowledgePanel.style.display   = tab === 'knowledge' ? 'block' : 'none';
    if (memoryListCard)   memoryListCard.style.display   = tab === 'knowledge' ? 'none'  : 'block';
    if (memoryEditorCard) memoryEditorCard.style.display = tab === 'knowledge' ? 'none'  : 'block';

    document.querySelectorAll('.memories-tab').forEach(t => t.classList.remove('active'));
    document.getElementById(`memories-tab-${tab}`)?.classList.add('active');

    if (tab === 'knowledge') {
        if (typeof refreshKnowledgeList === 'function') refreshKnowledgeList();
    } else {
        loadFullMemoriesList();
    }
}"""
s = s.replace(old_tab, new_tab, 1)
wr(p, s)
print("story-tools.js: switchMemoriesTab now owns the knowledge panel")

# ---------------------------------------------------------------------------
# 4. rag-knowledge.js — delete the wrapper
# ---------------------------------------------------------------------------
p = "frontend/dist/js/rag-knowledge.js"
s = rd(p)
start = s.find("const _origSwitchMemoriesTab = window.switchMemoriesTab;")
must(start != -1, "rag-knowledge.js: wrapper start not found")
end = s.find("// ---------------------------------------------------------------------------\n// Helpers", start)
must(end != -1, "rag-knowledge.js: wrapper end not found")
removed = end - start
s = s[:start] + ("// switchMemoriesTab is owned by story-tools.js, which handles the knowledge\n"
                 "// panel itself; the wrapper that used to live here re-ran a captured copy.\n\n") + s[end:]
wr(p, s)
print("rag-knowledge.js: wrapper removed (%d bytes)" % removed)
