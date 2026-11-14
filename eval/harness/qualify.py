"""Slot qualification (WO-SLOT-QUAL-01): load, identify, release. Zero generations.

This module has no prompt path. It cannot call `generate`, it does not import the
adapter, and every request it makes is one of the three L2 identity GETs. That is
a property of the code, not a promise in a document: check 89 fails if a POST,
`generate(`, `dispatch(` or `Dispatcher` appears in this file.

A qualification record answers exactly one question — does this candidate work on
this machine — and is never evidence about MEMORY-09 behaviour.
"""

from __future__ import annotations

import hashlib
from datetime import datetime, timezone
from pathlib import Path

from . import discover as discovery
from .adapter import canonical, sha256_text

EVAL_ROOT = Path(__file__).resolve().parents[1]

# Identity GETs only. Qualification issues no other kind of request.
ALLOWED_REQUESTS = tuple(sorted(discovery.IDENTITY_PATHS))

GENERATIONS = 0


def normalise(text: str) -> str:
    """Strip case and punctuation: KoboldCpp sanitises '(' and ')' out of names.

    The file is `...Q4_K_M(1).gguf`; the server reports `...Q4_K_M1`. Comparing
    raw strings would report a false mismatch on a correct load.
    """
    return "".join(ch for ch in str(text).casefold() if ch.isalnum())


def utc_now() -> str:
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def file_digest(path: Path, chunk: int = 1 << 20) -> tuple[str, int]:
    digest = hashlib.sha256()
    size = 0
    with path.open("rb") as handle:
        for block in iter(lambda: handle.read(chunk), b""):
            digest.update(block)
            size += len(block)
    return digest.hexdigest(), size


def gpu_state() -> dict:
    """Observational: read the counters, change nothing."""
    import subprocess

    out: dict = {"captured_at": utc_now(), "available": False}
    try:
        proc = subprocess.run(
            ["nvidia-smi", "--query-gpu=memory.used,memory.free,utilization.gpu,temperature.gpu",
             "--format=csv,noheader"],
            capture_output=True, text=True, timeout=20,
        )
        line = (proc.stdout or "").strip()
        if proc.returncode == 0 and line:
            used, free, util, temp = [p.strip() for p in line.split(",")]
            out.update({"available": True, "vram_used_mib": used, "vram_free_mib": free,
                        "gpu_util_percent": util, "temperature_c": temp})
    except Exception as exc:
        out["reason"] = f"{type(exc).__name__}: {exc}"[:160]
    return out


def build_record(*, slot: str, model_path: str, expected_sha256: str, expected_bytes: int,
                 digest: str, size: int, identity_block: dict, context: int, gpu_layers: int,
                 quantkv: str, vram_before: dict, log_path: str | None = None) -> dict:
    reported = identity_block.get("loaded_model")
    stem = Path(model_path).stem
    matches = bool(reported) and normalise(stem) in normalise(reported)
    record = {
        "record_type": "slot_qualification",
        "report_id": "WO-SLOT-QUAL-01",
        "slot": slot,
        "generations": GENERATIONS,
        "requests_issued": list(ALLOWED_REQUESTS),
        "request_kinds": "identity GETs only",
        "model": {
            "path": model_path,
            "bytes": size,
            "sha256": digest,
            "expected_sha256": expected_sha256,
            "expected_bytes": expected_bytes,
            "digest_matches_baseline": digest == expected_sha256,
            "bytes_match_baseline": size == expected_bytes,
        },
        "configuration": {
            "context": context,
            "gpu_layers": gpu_layers,
            "quantkv": quantkv,
            "host": "127.0.0.1",
            "frozen_at_qualification": True,
        },
        "server": {
            "name": identity_block.get("server_name"),
            "version": identity_block.get("server_version"),
            "endpoint": identity_block.get("endpoint"),
            "identity_block_sha256": identity_block.get("block_sha256"),
        },
        "identity": {
            "reported": reported,
            "expected_stem": stem,
            "matches": matches,
            "match_rule": "normalised substring (casefold, alphanumerics only)",
        },
        "hardware": {"before_load": vram_before, "after_load": gpu_state()},
        "claim_ceiling": "compatibility of this file with this build on this machine; NOTHING about MEMORY-09",
        "timestamp": utc_now(),
    }
    if log_path:
        record["server"]["log_path"] = log_path
        record["server"]["log_evidence"] = extract_load_evidence(log_path)
    record["record_sha256"] = sha256_text(canonical(record))
    return record


def extract_load_evidence(log_path: str) -> dict:
    """Pull the load-time facts the server states about itself."""
    import re

    out: dict = {"offloaded_layers": None, "gpu_buffer_mib": None, "cpu_buffer_mib": None,
                 "n_layer": None, "tensors": None}
    try:
        text = Path(log_path).read_text(encoding="utf-8", errors="replace")
    except Exception as exc:
        out["error"] = f"{type(exc).__name__}: {exc}"
        return out
    m = re.search(r"offloaded (\d+)/(\d+) layers to GPU", text)
    if m:
        out["offloaded_layers"] = f"{m.group(1)}/{m.group(2)}"
    m = re.search(r"CUDA0 model buffer size =\s*([\d.]+) MiB", text)
    if m:
        out["gpu_buffer_mib"] = float(m.group(1))
    m = re.search(r"CPU_Mapped model buffer size =\s*([\d.]+) MiB", text)
    if m:
        out["cpu_buffer_mib"] = float(m.group(1))
    m = re.search(r"print_info: n_layer\s*=\s*(\d+)", text)
    if m:
        out["n_layer"] = int(m.group(1))
    m = re.search(r"loaded meta data with \d+ key-value pairs and (\d+) tensors", text)
    if m:
        out["tensors"] = int(m.group(1))
    return out


def write_qualification(folder: Path, record: dict) -> Path:
    """Write-once, inside eval/ only."""
    folder = folder.resolve()
    if not str(folder).startswith(str(EVAL_ROOT.resolve())):
        raise ValueError(f"refusing to write evidence outside eval/: {folder}")
    target = folder / "record.json"
    if target.exists():
        raise FileExistsError(f"qualification evidence is write-once: {target} exists")
    folder.mkdir(parents=True, exist_ok=True)
    import json

    target.write_text(json.dumps(record, indent=2, sort_keys=True, ensure_ascii=False) + "\n",
                      encoding="utf-8")
    (folder / "record.sha256").write_text(
        sha256_text(canonical(record)) + "  record.json\n", encoding="utf-8")
    return target


def write_recovery(folder: Path, *, vram_reference: dict, port_closed: bool) -> Path:
    """Recovery evidence after the server is terminated. Never rewrites the record."""
    import json

    folder = folder.resolve()
    now = gpu_state()
    try:
        ref_used = int(str(vram_reference.get("vram_used_mib", "0")).split()[0])
        now_used = int(str(now.get("vram_used_mib", "0")).split()[0])
        delta = now_used - ref_used
        recovered = abs(delta) <= 250
    except Exception:
        delta, recovered = None, None
    payload = {
        "record_type": "qualification_recovery",
        "after_unload": now,
        "baseline_vram_used_mib": vram_reference.get("vram_used_mib"),
        "delta_mib": delta,
        "recovery_gate": "pass" if recovered else ("fail" if recovered is False else "unavailable"),
        "recovery_tolerance_mib": 250,
        "port_closed": port_closed,
        "generations": GENERATIONS,
        "timestamp": utc_now(),
    }
    target = folder / "recovery.json"
    if target.exists():
        raise FileExistsError(f"recovery evidence is write-once: {target} exists")
    target.write_text(json.dumps(payload, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    return target


def stamp() -> str:
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H%M%SZ")
