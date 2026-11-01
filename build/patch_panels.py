"""Insert the panel-less-page no-op shims into sidebar-panels.js (CRLF safe)."""
import io
import re

path = "frontend/dist/js/sidebar-panels.js"
s = io.open(path, encoding="utf-8", newline="").read()
nl = chr(13) + chr(10) if chr(13) + chr(10) in s else chr(10)

if "window.updateContextDisplay  = function () {}" in s:
    print("already patched")
    raise SystemExit(0)

lines = [
    "// context-display.js owns updateContextDisplay / updateLorebookDisplay and the",
    "// three panel renderers: it defines them on window and instantiates ContextDisplay.",
    "// story-tools.html loads these sidebar modules without that orchestrator and has",
    "// no world, scenario or lorebook panel at all, so its shared session code used to",
    "// call a function that was either undefined or threw on an undefined contextDisplay.",
    "// Panels that do not exist on the page get an explicit no-op.",
    "if (typeof window.contextDisplay === 'undefined') {",
    "    window.updateContextDisplay  = function () {};",
    "    window.updateLorebookDisplay = function () {};",
    "    window.updateWorldDisplay    = function () {};",
    "    window.updateScenarioDisplay = function () {};",
    "}",
]

m = re.search(r"^function saveChatData\(\).*$", s, flags=re.M)
if not m:
    print("ANCHOR MISSING: saveChatData definition not found")
    raise SystemExit(1)

s = s[:m.end()] + nl + nl + nl.join(lines) + s[m.end():]
io.open(path, "w", encoding="utf-8", newline="").write(s)
print("sidebar-panels.js: no-op shims inserted after saveChatData")
print("contextDisplay guard present:", "typeof window.contextDisplay === 'undefined'" in s)
