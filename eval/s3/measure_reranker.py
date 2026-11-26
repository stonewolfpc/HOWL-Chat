"""Complete the reranker resource envelope: CPU, partial GPU, full GPU.

Behavioral results are NOT re-measured here; they stay in
eval/s3/rerank-v12-receipt.json and eval/s3/prose-experiment-receipt.json.

Usage: PYTHONPATH=. python eval/s3/measure_reranker.py
Writes eval/s3/reranker-resources.json
"""
from __future__ import annotations

import json
import pathlib
import re
import statistics
import subprocess
import time
import urllib.request

EXE = r"H:/AI apps/llama.cpp-b10076-cuda12/runtime/llama-server.exe"
MODEL = "eval/s3/models/bge-reranker-v2-m3-Q8_0.gguf"
BASE = "http://127.0.0.1:8092"
OUT = pathlib.Path("eval/s3/reranker-resources.json")

CONFIGS = {
    "cpu": ["--device", "none", "-ngl", "0"],
    "partial_6": ["-ngl", "6"],
    "partial_12": ["-ngl", "12"],
    "full_gpu": ["-ngl", "99"],
}

DOCS = [
    "Sable trained_by: Corvin Ashgrave of Ossary, and quiet was not part of it",
    "Sable preference: tea more than I like surprises",
    "Sable never: trusted a map that was drawn by someone who had not walked the ground",
    "Sable keeps_ledger: of debts in the green cabinet now",
]
QUERY = "Who was your trainer, again?"


def vram() -> int:
    out = subprocess.run(["nvidia-smi", "--query-gpu=memory.used",
                          "--format=csv,noheader,nounits"], capture_output=True, text=True).stdout
    return int(out.strip().splitlines()[0])


def median_vram(n: int = 5) -> float:
    v = []
    for _ in range(n):
        v.append(vram()); time.sleep(1.0)
    return statistics.median(v)


def rss_mb() -> float | None:
    try:
        out = subprocess.run(["powershell", "-NoProfile", "-Command",
                              "(Get-Process llama-server | Select-Object -First 1).WorkingSet64"],
                             capture_output=True, text=True, timeout=30).stdout.strip()
        return round(int(out) / 1e6, 0) if out.isdigit() else None
    except Exception:
        return None


def health() -> bool:
    try:
        with urllib.request.urlopen(f"{BASE}/health", timeout=2) as r:
            return r.status == 200
    except Exception:
        return False


def rerank(q: str, docs: list[str]) -> list[float]:
    body = json.dumps({"query": q, "documents": docs}).encode()
    req = urllib.request.Request(f"{BASE}/v1/rerank", data=body,
                                 headers={"Content-Type": "application/json"})
    with urllib.request.urlopen(req, timeout=180) as r:
        payload = json.loads(r.read().decode())
    m = {d["index"]: d["relevance_score"] for d in payload["results"]}
    return [m[i] for i in range(len(docs))]


def load_seconds(log: str) -> float | None:
    start = end = None
    for line in log.splitlines():
        m = re.match(r"(\d+)\.(\d+)\.(\d+)\.(\d+)", line)
        if not m:
            continue
        # llama-server stamps look like "0.02.101.077" = 0 min, 2 s, 101 ms, 077 us
        ts = int(m.group(1)) * 60 + int(m.group(2)) + int(m.group(3)) / 1000 + int(m.group(4)) / 1e6
        if "load_model: loading model" in line:
            start = ts
        elif "model loaded" in line and end is None:
            end = ts
    return round(end - start, 2) if start and end else None


def placement(log: str) -> dict:
    d = {}
    for line in log.splitlines():
        if "offloaded" in line and "layers to GPU" in line:
            d["offload"] = line.split("I ")[-1].strip()
        elif "CUDA0 model buffer size" in line:
            d["cuda0_buffer"] = line.split("size =")[-1].strip()
        elif "CPU_Mapped model buffer size" in line:
            d["cpu_buffer"] = line.split("size =")[-1].strip()
    return d


def run_config(name: str, flags: list[str]) -> dict:
    base = median_vram()
    log_path = pathlib.Path(f"eval/s3/models/rr-{name}.log")
    t0 = time.perf_counter()
    with log_path.open("w", encoding="utf-8") as fh:
        p = subprocess.Popen([EXE, "--model", MODEL, "--reranking", "-c", "2048",
                              *flags, "-lv", "4", "--host", "127.0.0.1", "--port", "8092"],
                             stdout=fh, stderr=subprocess.STDOUT)
    ready = None
    while time.perf_counter() - t0 < 180:
        if health():
            ready = time.perf_counter() - t0
            break
        if p.poll() is not None:
            return {"config": name, "error": f"exited rc={p.returncode}"}
        time.sleep(0.2)
    if ready is None:
        p.terminate(); return {"config": name, "error": "never ready"}

    loaded = median_vram(3)
    mem = rss_mb()
    t = time.perf_counter(); rerank(QUERY, [DOCS[0]]); first = (time.perf_counter() - t) * 1000
    single = []
    for _ in range(12):
        t = time.perf_counter(); rerank(QUERY, [DOCS[0]]); single.append((time.perf_counter() - t) * 1000)
    batched = []
    for _ in range(12):
        t = time.perf_counter(); rerank(QUERY, DOCS); batched.append((time.perf_counter() - t) * 1000)
    single.sort(); batched.sort()

    p.terminate()
    try:
        p.wait(timeout=25)
    except subprocess.TimeoutExpired:
        p.kill(); p.wait(timeout=10)
    time.sleep(5)
    after = median_vram()
    log = log_path.read_text(encoding="utf-8", errors="replace")
    return {
        "config": name, "flags": flags,
        "spawn_to_ready_s": round(ready, 2), "model_load_s": load_seconds(log),
        "placement": placement(log),
        "first_call_ms": round(first, 1),
        "single_p50_ms": round(statistics.median(single), 1),
        "single_p95_ms": round(single[int(len(single) * 0.95)], 1),
        "batched4_p50_ms": round(statistics.median(batched), 1),
        "batched4_p95_ms": round(batched[int(len(batched) * 0.95)], 1),
        "vram_idle_mib": round(base), "vram_loaded_mib": round(loaded),
        "vram_delta_mib": round(loaded - base),
        "vram_residual_mib": round(after - base),
        "rss_mb": mem, "port_closed": not health(),
    }


def main() -> int:
    results = {}
    for name, flags in CONFIGS.items():
        print(f"--- {name} ---", flush=True)
        results[name] = run_config(name, flags)
        print(json.dumps(results[name], indent=2), flush=True)
        time.sleep(2)

    print("\n=== envelope table (budget: VRAM <=384 MiB, p50 <=40 ms, p95 <=120 ms, RSS <=1.0 GiB) ===")
    for name, r in results.items():
        if "error" in r:
            print(f"  {name:<11} ERROR {r['error']}")
            continue
        ok_lat = r["batched4_p50_ms"] <= 40 and r["batched4_p95_ms"] <= 120
        ok_vram = r["vram_delta_mib"] <= 384
        ok_rss = (r["rss_mb"] or 0) <= 1024
        print(f"  {name:<11} VRAM {r['vram_delta_mib']:>4} MiB {'PASS' if ok_vram else 'FAIL'} | "
              f"batched p50 {r['batched4_p50_ms']:>5} p95 {r['batched4_p95_ms']:>5} ms "
              f"{'PASS' if ok_lat else 'FAIL'} | RSS {r['rss_mb']} MB {'PASS' if ok_rss else 'FAIL'} | "
              f"ready {r['spawn_to_ready_s']} s | residual {r['vram_residual_mib']} MiB")

    OUT.write_text(json.dumps({"budgets": {"vram_mib": 384, "p50_ms": 40, "p95_ms": 120,
                                           "rss_mb": 1024, "residual_mib": 64},
                               "configs": results,
                               "behaviour_note": "behavioural results unchanged; see rerank-v12-receipt.json"},
                              indent=2) + "\n", encoding="utf-8")
    print(f"\nreceipt: {OUT} ({OUT.stat().st_size} B)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
