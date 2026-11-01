"""Print the real-shadow duplicate functions side by side, with per-file page loads."""
import glob
import json
import os
import re

ROOT = "frontend/dist"
data = json.load(open("build/dupdefs.json"))

pages = {}
for p in sorted(glob.glob(os.path.join(ROOT, "*.html"))):
    s = open(p, encoding="utf-8", errors="replace").read()
    pages[os.path.basename(p)] = [os.path.basename(x.split("?")[0])
                                  for x in re.findall(r'<script[^>]+src="js/([^"]+)"', s)]

only = os.environ.get("ONLY", "").split(",") if os.environ.get("ONLY") else None

for rec in data["real"]:
    if only and rec["name"] not in only:
        continue
    print("=" * 78)
    print("%s   (%d variants)" % (rec["name"], rec["variants"]))
    for f, body in rec["bodies"].items():
        pgs = [pg for pg, srcs in pages.items() if f in srcs]
        print("-" * 78)
        print("%s   loaded by: %s" % (f, ", ".join(pgs)))
        print(body[:1400])
    print()
