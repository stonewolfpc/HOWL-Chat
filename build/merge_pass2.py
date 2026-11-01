"""T2.2 final sub-pass, part 2: one owner per remaining shared name.

Removals use an async-aware, quote/comment-aware brace matcher and verify the
top-level function inventory afterwards: exactly the requested name may vanish.
"""
import io
import re

def rd(p):
    return io.open(p, encoding="utf-8", newline="").read()

def wr(p, s):
    io.open(p, "w", encoding="utf-8", newline="").write(s)

def must(cond, msg):
    if not cond:
        raise SystemExit("FAILED: " + msg)

def defs(text):
    return re.findall(r"^(?:async\s+)?function\s+([A-Za-z_$][\w$]*)\s*\(", text, flags=re.M)

def cut(text, name):
    """Remove the single top-level (async) definition of `name`."""
    m = list(re.finditer(r"^(?:async\s+)?function\s+%s\s*\(" % re.escape(name), text, flags=re.M))
    must(len(m) == 1, "%s: expected exactly one definition, found %d" % (name, len(m)))
    start = m[0].start()
    brace = text.index("{", m[0].end())
    depth, j, n = 0, brace, len(text)
    while j < n:
        c = text[j]
        if c in "\"'`":
            quote = c
            j += 1
            while j < n:
                if text[j] == chr(92):
                    j += 2
                    continue
                if text[j] == quote:
                    break
                j += 1
        elif c == "/" and j + 1 < n and text[j + 1] == "/":
            nl = text.find("\n", j)
            j = n if nl < 0 else nl
            continue
        elif c == "/" and j + 1 < n and text[j + 1] == "*":
            close = text.find("*/", j + 2)
            j = n if close < 0 else close + 1
        elif c == "{":
            depth += 1
        elif c == "}":
            depth -= 1
            if depth == 0:
                break
        j += 1
    return text[:start] + text[j + 1:]

# ---------------------------------------------------------------------------
# 1. sidebar-context-lorebooks.js — rename the chat-level toggle
# ---------------------------------------------------------------------------
p = "frontend/dist/js/sidebar-context-lorebooks.js"
s = rd(p)
if "async function toggleLorebookSelection(lorebookId) {" in s:
    s = s.replace("async function toggleLorebookSelection(lorebookId) {",
                  "// Chat-level toggle, by id. The world/scenario dropdowns own their own\n"
                  "// toggleLorebookSelection(lorebook, element); same name, different job,\n"
                  "// so this one is named for what it does.\n"
                  "async function toggleChatLorebook(lorebookId) {", 1)
    s = s.replace('onclick="toggleLorebookSelection(', 'onclick="toggleChatLorebook(')
    wr(p, s)
    print("sidebar-context-lorebooks.js: renamed to toggleChatLorebook")
else:
    must("async function toggleChatLorebook(lorebookId) {" in s,
         "chat lorebook toggle missing entirely")
    print("sidebar-context-lorebooks.js: already renamed")

# ---------------------------------------------------------------------------
# 2. story-tools.js — own switchMemoriesTab
# ---------------------------------------------------------------------------
p = "frontend/dist/js/story-tools.js"
s = rd(p)
old_tab = """function switchMemoriesTab(tab) {
    window.storyToolsState.currentTab = tab;
    document.querySelectorAll('.memories-tab').forEach(t => t.classList.remove('active'));
    document.getElementById(`memories-tab-${tab}`)?.classList.add('active');
    loadFullMemoriesList();
}"""
must(old_tab in s, "story-tools.js switchMemoriesTab not found")
new_tab = """function switchMemoriesTab(tab) {
    window.storyToolsState.currentTab = tab;

    // The knowledge tab swaps the memory list and editor for the RAG panel. This
    // logic used to live in rag-knowledge.js as a wrapper that captured an
    // "original" and re-ran it conditionally; one owner is enough.
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
print("story-tools.js: switchMemoriesTab owns the knowledge panel")

# ---------------------------------------------------------------------------
# 3. rag-knowledge.js — delete the wrapper
# ---------------------------------------------------------------------------
p = "frontend/dist/js/rag-knowledge.js"
s = rd(p)
start = s.find("const _origSwitchMemoriesTab = window.switchMemoriesTab;")
must(start != -1, "rag-knowledge wrapper start not found")
end = s.find("// ---------------------------------------------------------------------------\n// Helpers", start)
must(end != -1, "rag-knowledge wrapper end not found")
print("rag-knowledge.js: wrapper removed (%d bytes)" % (end - start))
s = s[:start] + ("// switchMemoriesTab is owned by story-tools.js, which handles the knowledge\n"
                 "// panel itself. The wrapper that re-ran a captured copy is gone.\n\n") + s[end:]
wr(p, s)

# ---------------------------------------------------------------------------
# 4. story-tools-data.js — delete the no-argument deleteMemory
# ---------------------------------------------------------------------------
p = "frontend/dist/js/story-tools-data.js"
s = rd(p)
before = defs(s)
s = cut(s, "deleteMemory")
after = defs(s)
must(set(before) - set(after) == {"deleteMemory"}, "story-tools-data inventory changed unexpectedly")
must(len(before) - len(after) == 1, "story-tools-data lost more than one function")
wr(p, s)
print("story-tools-data.js: deleteMemory removed (%d -> %d functions)" % (len(before), len(after)))

# ---------------------------------------------------------------------------
# 5. story-tools-ui.js — one deleteMemory that serves both call styles
# ---------------------------------------------------------------------------
p = "frontend/dist/js/story-tools-ui.js"
s = rd(p)
start = s.find("async function deleteMemory(memoryId, memoryTitle) {")
must(start != -1, "story-tools-ui deleteMemory not found")
end = s.find("function showSuccess(message)", start)
must(end != -1, "story-tools-ui deleteMemory end not found")
merged = """async function deleteMemory(memoryId, memoryTitle) {
    const st = window.storyToolsState;
    // Two call styles, one implementation: the editor's Delete button passes
    // nothing (it acts on the selected memory), list rows pass id and title.
    const id = memoryId || st?.currentMemoryId;
    if (!id) { showError('No memory selected to delete'); return; }

    const known = (st?.memoriesData || []).find(m => m.id === id);
    const title = memoryTitle || known?.title;
    if (!confirm(`Delete memory "${title || 'Untitled'}"? This cannot be undone.`)) return;

    try {
        const deleted = window.MemoryStore?.deleteMemoryEntry
            ? window.MemoryStore.deleteMemoryEntry(id)
            : false;

        // The backend copy stays authoritative for the chat; a failure here is
        // logged, not fatal (the local store already dropped it).
        if (st?.currentChatId && window.go?.main?.App?.DeleteMemory) {
            try { await window.go.main.App.DeleteMemory(st.currentChatId, id); }
            catch (err) { console.warn('Backend memory delete failed:', err); }
        }

        if (!deleted) { showError('Memory not found'); return; }

        renderRecentList(await getMemories({ limit: 5, sort: 'newest' }));
        renderFullList(await getMemories({}));
        if (st?.currentMemoryId === id) clearEditor();
        showSuccess('Memory deleted');
    } catch (error) {
        console.error('Failed to delete memory:', error);
        showError('Failed to delete memory');
    }
}

"""
s = s[:start] + merged + s[end:]
wr(p, s)
print("story-tools-ui.js: deleteMemory merged for both call styles")

# ---------------------------------------------------------------------------
# 6. error-system.js — call the recovery hook here, not from a patched copy
# ---------------------------------------------------------------------------
p = "frontend/dist/js/error-system.js"
s = rd(p)
anchor = """    incrementCircuitBreaker(code);
}"""
must(anchor in s, "error-system handleBackendError tail not found")
s = s.replace(anchor, """    incrementCircuitBreaker(code);

    // Recovery hooks (error-handler.js) run from here now. It used to wrap this
    // function after load, which only worked while the load order held.
    if (typeof onErrorReported === 'function') onErrorReported(payload);
}""", 1)
wr(p, s)
print("error-system.js: invokes onErrorReported itself")

# ---------------------------------------------------------------------------
# 7. error-handler.js — remove the patch
# ---------------------------------------------------------------------------
p = "frontend/dist/js/error-handler.js"
s = rd(p)
old_patch = """// Patch ErrorSystem to trigger recovery
const originalHandleBackendError = window.handleBackendError;
window.handleBackendError = function(payload) {
    if (originalHandleBackendError) originalHandleBackendError(payload);
    onErrorReported(payload);
};"""
must(old_patch in s, "error-handler patch block not found")
s = s.replace(old_patch, """// ErrorSystem calls onErrorReported() itself now, so handleBackendError is not
// monkey-patched. Recovery happens even if this file loads after error-system.""", 1)
wr(p, s)
print("error-handler.js: patch removed")

# ---------------------------------------------------------------------------
# 8. sidebar-characters.js — context-display owns the panel; drop the dead copy
# ---------------------------------------------------------------------------
p = "frontend/dist/js/sidebar-characters.js"
s = rd(p)
before = defs(s)
s = cut(s, "updateCharacterPanel")
after = defs(s)
must(set(before) - set(after) == {"updateCharacterPanel"}, "sidebar-characters inventory changed unexpectedly")
# every remaining bare call must survive a page without the panel
s = s.replace("    await updateCharacterPanel();", "    if (window.updateCharacterPanel) await window.updateCharacterPanel();")
s = s.replace("                updateCharacterPanel();", "                if (window.updateCharacterPanel) window.updateCharacterPanel();")
s = s.replace("            updateCharacterPanel();", "            if (window.updateCharacterPanel) window.updateCharacterPanel();")
s = s.replace("        updateCharacterPanel();", "        if (window.updateCharacterPanel) window.updateCharacterPanel();")
s = s.replace("        if (typeof updateCharacterPanel === 'function') updateCharacterPanel();",
              "        if (typeof window.updateCharacterPanel === 'function') window.updateCharacterPanel();")
s = s.replace("    updateCharacterPanel,\n", "    updateCharacterPanel: (...args) => window.updateCharacterPanel?.(...args),\n", 1)
wr(p, s)
print("sidebar-characters.js: renderer removed, calls routed to the single owner")
