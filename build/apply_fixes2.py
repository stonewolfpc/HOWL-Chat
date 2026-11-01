"""Second structural pass: panel-less page shims + system theme listener."""
import io

changed = []


def read(p):
    return io.open(p, encoding="utf-8", newline="").read()


def write(p, s):
    io.open(p, "w", encoding="utf-8", newline="").write(s)
    changed.append(p)


# ---------------------------------------------------------------------------
# 1. sidebar-panels.js: panel-less pages get working no-ops
# ---------------------------------------------------------------------------
p = "frontend/dist/js/sidebar-panels.js"
s = read(p)

stale = "// Re-establish canonical panel shims after sidebar modules load (they overwrite context-display.js shims)\n"
s = s.replace(stale, "")

fallback = (
    "// context-display.js owns updateContextDisplay / updateLorebookDisplay and the\n"
    "// three panel renderers: it defines them on window and instantiates ContextDisplay.\n"
    "// story-tools.html loads these sidebar modules without that orchestrator and has\n"
    "// no world, scenario or lorebook panel at all, so its shared session code used to\n"
    "// call an undefined function (or one that threw on an undefined contextDisplay).\n"
    "// Panels that do not exist get an explicit no-op instead.\n"
    "if (typeof window.contextDisplay === 'undefined') {\n"
    "    window.updateContextDisplay  = function () {};\n"
    "    window.updateLorebookDisplay = function () {};\n"
    "    window.updateWorldDisplay    = function () {};\n"
    "    window.updateScenarioDisplay = function () {};\n"
    "}\n"
)

anchor = "function saveChatData() { return window.sidebarSession.saveChatData(); }\n"
if anchor in s and "window.updateContextDisplay  = function () {}" not in s:
    s = s.replace(anchor, anchor + fallback, 1)
    write(p, s)
    print("sidebar-panels.js: no-op panel shims added for panel-less pages")
else:
    print("sidebar-panels.js: anchor missing or already patched")

# ---------------------------------------------------------------------------
# 2. theme.js: follow the OS while the preference is 'system'
# ---------------------------------------------------------------------------
p = "frontend/dist/js/theme.js"
s = read(p)

init_block = (
    "if (document.readyState === 'loading') {\n"
    "    document.addEventListener('DOMContentLoaded', loadThemePreference);\n"
    "} else {\n"
    "    loadThemePreference();\n"
    "}\n"
)
listener = (
    "\n// Follow the OS while the preference is 'system'. This listener used to sit in\n"
    "// menu-view.js and referenced a variable that no longer existed, so switching\n"
    "// the OS theme never reached the app. theme.js owns the preference, so it owns\n"
    "// the reaction to it too.\n"
    "if (window.matchMedia) {\n"
    "    window.matchMedia('(prefers-color-scheme: light)').addEventListener('change', function () {\n"
    "        if (currentThemeMode !== 'system') return;\n"
    "        const resolved = resolveThemeMode('system');\n"
    "        if (resolved) setThemeMode(resolved, false);\n"
    "    });\n"
    "}\n"
)
if init_block in s and "Follow the OS while the preference" not in s:
    s = s.replace(init_block, init_block + listener, 1)
    write(p, s)
    print("theme.js: system-mode OS change listener added")
else:
    print("theme.js: init block missing or already patched")

print("\nfiles changed:", len(set(changed)))
