"""Reproducible embedder resource measurement.

Keeps the four phases separate, as required:
  spawn -> model load -> first query -> warm queries -> shutdown release

Usage: PYTHONPATH=. python eval/s3/measure_embedder.py cpu|gpu
Writes eval/s3/measure-<config>.json
"""
from __future__ import annotations

import json
import pathlib
import re
import statistics
import subprocess
import sys
import time
import urllib.request

EXE = r"H:/AI apps/llama.cpp-b10076-cuda12/runtime/llama-server.exe"
MODEL = "eval/s3/models/nomic-embed-text-v1.5.Q8_0.gguf"
PORT = 8091
BASE = f"http://127.0.0.1:{PORT}"
CONFIGS = {
    "cpu": ["--device", "none", "-ngl", "0"],
    "gpu": ["-ngl", "99"],
}


def vram() -> int:
    out = subprocess.run(["nvidia-smi", "--query-gpu=memory.used",
                          "--format=csv,noheader,nounits"], capture_output=True, text=True).stdout
    return int(out.strip().splitlines()[0])


def idle_vram(samples: int = 5) -> float:
    vals = []
    for _ in range(samples):
        vals.append(vram())
        time.sleep(1.0)
    return statistics.median(vals)


def embed(text: str) -> None:
    body = json.dumps({"input": text, "model": "nomic-embed-text-v1.5"}).encode()
    req = urllib.request.Request(f"{BASE}/v1/embeddings", data=body,
                                 headers={"Content-Type": "application/json"})
    with urllib.request.urlopen(req, timeout=120) as r:
        r.read()


def health() -> bool:
    try:
        with urllib.request.urlopen(f"{BASE}/health", timeout=2) as r:
            return r.status == 200
    except Exception:
        return False


def parse_load(log: str) -> dict:
    times = {}
    for line in log.splitlines():
        m = re.match(r"(\d+\.\d+\.\d+\.\d+)", line)
        if not m:
            continue
        ts = m.group(1)
        if "load_model: loading model" in line:
            times["load_start"] = ts
        elif "model loaded" in line:
            times["load_end"] = ts
        elif "offloaded" in line and "layers to GPU" in line:
            times["offload"] = line.split("I ")[-1].strip()
        elif "CPU_Mapped model buffer size" in line:
            times["cpu_buffer"] = line.split("size =")[-1].strip()
        elif "CUDA0 model buffer size" in line:
            times["gpu_buffer"] = line.split("size =")[-1].strip()
    return times


def main() -> int:
    config = sys.argv[1] if len(sys.argv) > 1 else "cpu"
    if config not in CONFIGS:
        print(f"unknown config {config}")
        return 2
    log_path = pathlib.Path(f"eval/s3/models/measure-{config}.log")
    idle = idle_vram()
    print(f"[{config}] idle VRAM median: {idle:.0f} MiB")

    cmd = [EXE, "--model", MODEL, "--embedding", "--pooling", "mean",
           "--ctx-size", "2048", "--threads", "8", *CONFIGS[config],
           "-lv", "4", "--host", "127.0.0.1", "--port", str(PORT)]
    t0 = time.perf_counter()
    with log_path.open("w", encoding="utf-8") as fh:
        proc = subprocess.Popen(cmd, stdout=fh, stderr=subprocess.STDOUT)
    ready_after = None
    deadline = t0 + 90
    while time.perf_counter() < deadline:
        if health():
            ready_after = time.perf_counter() - t0
            break
        if proc.poll() is not None:
            print(f"[{config}] process exited early rc={proc.returncode}")
            print(log_path.read_text(encoding="utf-8", errors="replace")[-1500:])
            return 1
        time.sleep(0.25)
    if ready_after is None:
        proc.terminate()
        print(f"[{config}] never became ready")
        return 1

    during = idle_vram(3)
    t = time.perf_counter(); embed("search_query: who trained you"); first = (time.perf_counter() - t) * 1000
    lat = []
    for _ in range(20):
        t = time.perf_counter(); embed("search_query: where is the ledger kept now")
        lat.append((time.perf_counter() - t) * 1000)
    lat.sort()

    proc.terminate()
    try:
        proc.wait(timeout=20)
    except subprocess.TimeoutExpired:
        proc.kill(); proc.wait(timeout=10)
    time.sleep(5)
    after = idle_vram(5)
    log = log_path.read_text(encoding="utf-8", errors="replace")
    placement = parse_load(log)

    result = {
        "config": config, "flags": CONFIGS[config],
        "spawn_to_ready_seconds": round(ready_after, 2),
        "model_load_seconds": _dur(placement),
        "placement": placement,
        "first_query_ms": round(first, 1),
        "warm_p50_ms": round(statistics.median(lat), 1),
        "warm_p95_ms": round(lat[int(len(lat) * 0.95)], 1),
        "warm_max_ms": round(lat[-1], 1),
        "vram_idle_mib": round(idle, 0),
        "vram_loaded_mib": round(during, 0),
        "vram_delta_mib": round(during - idle, 0),
        "vram_after_shutdown_mib": round(after, 0),
        "vram_residual_mib": round(after - idle, 0),
        "port_closed_after_shutdown": not health(),
    }
    out = pathlib.Path(f"eval/s3/measure-{config}.json")
    out.write_text(json.dumps(result, indent=2) + "\n", encoding="utf-8")
    print(json.dumps(result, indent=2))
    return 0


def _dur(placement: dict) -> float | None:
    # llama-server timestamps look like "0.00.194.855": h.mm.<ms>.<us>.
    def secs(ts: str) -> float:
        h, m, ms, us = (int(x) for x in ts.split("."))
        return h * 3600 + m * 60 + ms / 1000.0 + us / 1_000_000.0
    if "load_start" in placement and "load_end" in placement:
        return round(secs(placement["load_end"]) - secs(placement["load_start"]), 2)
    return None


if __name__ == "__main__":
    sys.exit(main())
