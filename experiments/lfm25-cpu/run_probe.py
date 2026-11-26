"""LFM2.5 CPU memory-triage probe.

One generation per model, CPU only (-ngl 0), default llama.cpp sampling,
llama-server's own chat template. Records prefill/decode speed, cold load
time, wall time, and peak RSS of the server process tree.

Not part of eval/ and not evidence for any MEMORY-09 requirement.
"""

from __future__ import annotations

import json
import os
import subprocess
import sys
import threading
import time
import urllib.error
import urllib.request
from pathlib import Path

import psutil

ROOT = Path(__file__).resolve().parent
REPO = ROOT.parent.parent
SERVER = REPO / "runtime" / "llama" / "llama-server.exe"
MODELS_DIR = ROOT / "models"
RAW = ROOT / "raw"
RESULTS = ROOT / "results"
PORT = 5099
BASE = f"http://127.0.0.1:{PORT}"

PROMPT = (
    "You are triaging candidate memories extracted from a roleplay session before "
    "they are written to long-term memory. For each item below, decide KEEP "
    "(durable, worth remembering), REJECT (transient or worthless), or HOLD (needs "
    "the user's confirmation), and give a reason of no more than twelve words. Watch "
    "for duplicates, contradictions, and anything marked as uncertain. Respond with "
    "exactly one line per item in the form `ID | DECISION | REASON`, no preamble and "
    "no extra commentary. Items: M1 — Sable was trained by Corvin Ashgrave of Ossary "
    '(stated as canon by the user). M2 — The user said "brb" and returned four '
    "minutes later. M3 — Sable is a fire mage who cannot swim. M4 — The user first "
    "said Sable has never met Corvin, then later said Corvin trained her. M5 — The "
    "user prefers short replies and dislikes long scene descriptions. M6 — A note "
    "says Sable's birthday is 3 Frostfell, but the user tagged it \"maybe canon\"."
)

MODELS = [
    ("LFM2.5-230M", "LFM2.5-230M-QAD-Q4_0.gguf"),
    ("LFM2.5-350M", "LFM2.5-350M-QAD-Q4_0.gguf"),
    ("LFM2.5-1.2B-Instruct", "LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf"),
]

# llama.cpp defaults, sent explicitly so the record shows them.
SAMPLING = {
    "temperature": 0.8,
    "top_p": 0.95,
    "top_k": 40,
    "repeat_penalty": 1.0,
    "max_tokens": 400,
}


def tree_rss(proc: psutil.Process) -> int:
    total = 0
    try:
        procs = [proc, *proc.children(recursive=True)]
    except psutil.NoSuchProcess:
        return 0
    for p in procs:
        try:
            total += p.memory_info().rss
        except psutil.NoSuchProcess:
            pass
    return total


def wait_health(timeout: float) -> float:
    start = time.perf_counter()
    while time.perf_counter() - start < timeout:
        try:
            with urllib.request.urlopen(f"{BASE}/health", timeout=2) as r:
                if r.status == 200:
                    return time.perf_counter() - start
        except (urllib.error.URLError, OSError):
            pass
        time.sleep(0.1)
    raise TimeoutError(f"server not healthy after {timeout}s")


def run_one(label: str, filename: str, draw: int = 1,
            extra_body: dict | None = None) -> dict:
    model_path = MODELS_DIR / filename
    outdir = RAW / f"{label}-d{draw}"
    outdir.mkdir(parents=True, exist_ok=True)
    logpath = outdir / "server.log"

    cmd = [
        str(SERVER),
        "-m", str(model_path),
        "-c", "2048",
        "-ngl", "0",
        "-np", "1",
        "--host", "127.0.0.1",
        "--port", str(PORT),
    ]
    ram_before = psutil.virtual_memory().available
    t0 = time.perf_counter()
    with open(logpath, "wb") as log:
        proc = subprocess.Popen(cmd, stdout=log, stderr=subprocess.STDOUT,
                                cwd=str(SERVER.parent))
    peak = {"rss": 0}
    stop = threading.Event()

    def sample():
        ps = psutil.Process(proc.pid)
        while not stop.is_set():
            peak["rss"] = max(peak["rss"], tree_rss(ps))
            stop.wait(0.1)

    sampler = threading.Thread(target=sample, daemon=True)
    sampler.start()

    record: dict = {
        "model": label,
        "draw": draw,
        "file": filename,
        "file_mb": round(model_path.stat().st_size / 1e6, 1),
        "command": cmd,
        "sampling": SAMPLING,
        "prompt": PROMPT,
    }
    try:
        record["load_s"] = round(wait_health(180), 2)
        body = json.dumps({
            "messages": [{"role": "user", "content": PROMPT}],
            "stream": False,
            **SAMPLING,
            **(extra_body or {}),
        }).encode()
        req = urllib.request.Request(
            f"{BASE}/v1/chat/completions", data=body,
            headers={"Content-Type": "application/json"}, method="POST")
        t1 = time.perf_counter()
        with urllib.request.urlopen(req, timeout=900) as r:
            payload = json.loads(r.read())
        record["wall_s"] = round(time.perf_counter() - t1, 2)
        (outdir / "response.json").write_text(
            json.dumps(payload, indent=2), encoding="utf-8")
        choice = payload["choices"][0]
        record["finish_reason"] = choice.get("finish_reason")
        record["content"] = choice["message"]["content"]
        if choice["message"].get("reasoning_content"):
            record["reasoning_content"] = choice["message"]["reasoning_content"]
        usage = payload.get("usage", {})
        record["prompt_tokens"] = usage.get("prompt_tokens")
        record["completion_tokens"] = usage.get("completion_tokens")
        timings = payload.get("timings", {})
        record["prefill_tok_s"] = round(timings.get("prompt_per_second", 0), 2)
        record["decode_tok_s"] = round(timings.get("predicted_per_second", 0), 2)
        record["prefill_ms"] = round(timings.get("prompt_ms", 0), 1)
        record["decode_ms"] = round(timings.get("predicted_ms", 0), 1)
    finally:
        elapsed = round(time.perf_counter() - t0, 2)
        stop.set()
        sampler.join(timeout=2)
        proc.terminate()
        try:
            proc.wait(timeout=20)
        except subprocess.TimeoutExpired:
            proc.kill()
            proc.wait(timeout=10)
        record["process_lifetime_s"] = elapsed
        record["peak_rss_mb"] = round(peak["rss"] / 1e6, 1)
        time.sleep(1)
        ram_after = psutil.virtual_memory().available
        record["ram_available_before_mb"] = round(ram_before / 1e6, 1)
        record["ram_available_after_mb"] = round(ram_after / 1e6, 1)
    return record


def main() -> int:
    missing = [f for _, f in MODELS if not (MODELS_DIR / f).exists()]
    if missing:
        print("missing model files:", missing, file=sys.stderr)
        return 2
    RESULTS.mkdir(parents=True, exist_ok=True)
    out = []
    for label, filename in MODELS:
        print(f"=== {label} ===", flush=True)
        try:
            rec = run_one(label, filename)
        except Exception as exc:  # one model failing must not lose the others
            rec = {"model": label, "file": filename, "error": repr(exc)}
            print("FAILED:", exc, file=sys.stderr, flush=True)
        print(json.dumps({k: v for k, v in rec.items() if k != "content"},
                         indent=2), flush=True)
        out.append(rec)
        (RESULTS / "results.json").write_text(
            json.dumps(out, indent=2), encoding="utf-8")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
