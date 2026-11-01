"""Does sd-server actually emit TAESD step previews?

HOWL's chain is: engine passes --taesd/--taesd-preview-only, the job polls
/sdapi/v1/progress, the frame becomes job.preview_image, the UI paints it.
Every link is wired in code, but the middle one - sd-server decoding a frame
and exposing it as current_image - was never observed. This runs a real
generation against the configured SDXL bundle and records what the endpoint
actually returns while sampling is in progress.
"""

import json
import os
import subprocess
import threading
import time
import urllib.request
import sys

ROOT = r"C:\Users\StoneWolfPC\HOWL_Chat"
MODEL = r"E:\HOWL_Chat\models\image\unholyDesireMixSinister_v90.safetensors"
TAESD = r"E:\HOWL_Chat\models\image\taesdxl.safetensors"
PORT = 8091
BASE = "http://127.0.0.1:%d" % PORT

# Exactly what buildServerArgs produces for this bundle.
cmd = [
    os.path.join(ROOT, "sd-server.exe"),
    "-m", MODEL,
    "--taesd", TAESD,
    "--taesd-preview-only",
    "--listen-port", str(PORT),
    "--threads", "8",
]

log = open(os.path.join(ROOT, "build", "preview_probe_server.log"), "w",
           encoding="utf-8", errors="replace")
print("launch:", " ".join(cmd), flush=True)
proc = subprocess.Popen(cmd, cwd=ROOT, stdout=subprocess.PIPE,
                        stderr=subprocess.STDOUT, text=True,
                        encoding="utf-8", errors="replace")

tae_lines = []


def pump():
    for line in proc.stdout:
        log.write(line)
        log.flush()
        low = line.lower()
        if "tae" in low or "preview" in low:
            tae_lines.append(line.rstrip()[:160])


threading.Thread(target=pump, daemon=True).start()


def get(path, timeout=10):
    with urllib.request.urlopen(BASE + path, timeout=timeout) as r:
        return r.status, r.read()


try:
    ready = False
    for i in range(240):
        try:
            if get("/sdapi/v1/progress", timeout=3)[0] == 200:
                ready = True
                print("server ready after %.0fs" % (i * 0.5), flush=True)
                break
        except Exception:
            pass
        time.sleep(0.5)
    if not ready:
        print("SERVER NEVER CAME UP", flush=True)
        sys.exit(2)

    payload = json.dumps({
        "prompt": "a red fox sitting in snow, soft light",
        "negative_prompt": "blurry",
        "width": 512,
        "height": 512,
        "steps": 8,
        "cfg_scale": 7.0,
        "seed": 42,
        "sample_method": "euler_a",
    }).encode()

    outcome = {}

    def run_job():
        t0 = time.time()
        try:
            req = urllib.request.Request(
                BASE + "/sdapi/v1/txt2img", data=payload,
                headers={"Content-Type": "application/json"})
            with urllib.request.urlopen(req, timeout=900) as r:
                body = r.read()
            outcome["secs"] = round(time.time() - t0, 1)
            j = json.loads(body)
            imgs = j.get("images") or []
            outcome["images"] = len(imgs)
            outcome["first_len"] = len(str(imgs[0])) if imgs else 0
        except Exception as e:
            outcome["error"] = "%s: %s" % (type(e).__name__, e)

    th = threading.Thread(target=run_job, daemon=True)
    th.start()

    samples = []
    previews = []
    while th.is_alive():
        try:
            _, body = get("/sdapi/v1/progress", timeout=3)
            j = json.loads(body)
            ci = j.get("current_image") or ""
            samples.append((round(j.get("progress", 0), 3), len(ci)))
            if ci:
                previews.append((len(ci), ci[:8]))
        except Exception:
            pass
        time.sleep(0.25)
    th.join(5)

    print("poll samples: %d" % len(samples), flush=True)
    print("progress values seen:", sorted({p for p, _ in samples}), flush=True)
    print("preview frames returned: %d" % len(previews), flush=True)
    if previews:
        lens = [l for l, _ in previews]
        print("  sizes: min=%d max=%d first_prefix=%r" % (min(lens), max(lens), previews[0][1]), flush=True)
        print("  distinct frames: %d" % len({l for l, _ in previews}), flush=True)
        print("  looks like raw base64 (no data: prefix): %s" % (not previews[0][1].startswith("data:")), flush=True)
    else:
        print("  NO PREVIEW FRAMES - sd-server never decoded one", flush=True)
    print("job outcome:", outcome, flush=True)
    print("server lines mentioning tae/preview: %d" % len(tae_lines), flush=True)
    for l in tae_lines[:10]:
        print("   ", l, flush=True)
finally:
    try:
        proc.kill()
    except Exception:
        pass
    log.close()
    print("server stopped", flush=True)
