"""Apply the non-trivial duplicate fixes that need more than a deletion."""
import io
import re

changed = []


def read(p):
    return io.open(p, encoding="utf-8", newline="").read()


def write(p, s):
    io.open(p, "w", encoding="utf-8", newline="").write(s)
    changed.append(p)


# ---------------------------------------------------------------------------
# 1. theme.js owns setTheme again
# ---------------------------------------------------------------------------
p = "frontend/dist/js/theme.js"
s = read(p)
new_set_theme = (
    "function setTheme(theme) {\n"
    "    // Single owner. menu-view.js used to define its own setTheme, which\n"
    "    // shadowed this one on six pages; it resolved 'system' and passed the\n"
    "    // result to setThemeMode. That resolution lives here now, because\n"
    "    // setThemeMode only ever accepts light or dark.\n"
    "    const mode = theme === 'system'\n"
    "        ? (window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark')\n"
    "        : theme;\n"
    "    setThemeMode(mode);\n"
    "}"
)
pattern = re.compile(r"function setTheme\(theme\)\s*\{\s*setThemeMode\(theme\);\s*\}")
if pattern.search(s):
    write(p, pattern.sub(new_set_theme.replace(chr(92), chr(92) * 2), s, count=1))
    print("theme.js: setTheme now resolves 'system' before delegating")
else:
    print("theme.js: setTheme shape unexpected, left alone")

# ---------------------------------------------------------------------------
# 2. context-display.js owns every panel shim
# ---------------------------------------------------------------------------
p = "frontend/dist/js/context-display.js"
s = read(p)
anchor = "window.updateCharacterPanel = () => contextDisplay.updateCharacterPanel();"
if anchor in s and "window.updateLorebookDisplay" not in s:
    addition = anchor + (
        "\nwindow.updateWorldDisplay    = () => contextDisplay.updateWorldDisplay();"
        "\nwindow.updateScenarioDisplay = () => contextDisplay.updateScenarioDisplay();"
        "\nwindow.updateLorebookDisplay = () => contextDisplay.updateLorebookDisplay();"
    )
    write(p, s.replace(anchor, addition, 1))
    print("context-display.js: world/scenario/lorebook shims added")
else:
    print("context-display.js: shims already present or anchor missing")

# ---------------------------------------------------------------------------
# 3. html-utils.js on every page that renders escaped text
# ---------------------------------------------------------------------------
tag = '<script src="js/html-utils.js?v=16"></script>'
for page in ["index.html", "story-tools.html", "lorebooks.html"]:
    path = "frontend/dist/" + page
    s = read(path)
    if "html-utils.js" in s:
        print("%-18s already loads html-utils.js" % page)
        continue
    m = re.search(r'[ \t]*<script src="js/[^"]+"></script>', s)
    if not m:
        print("%-18s no js script tag found, skipped" % page)
        continue
    indent = re.match(r"[ \t]*", m.group(0)).group(0)
    write(path, s[:m.start()] + indent + tag + "\n" + s[m.start():])
    print("%-18s html-utils.js inserted before first script" % page)

print("\nfiles changed:", len(set(changed)))
