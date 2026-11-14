"""S1 execution preflight (read-only).

Answers seven questions against the *filed* evidence and the frozen manifest.
It makes no request, starts nothing, and loads nothing; hashing a model file is
the only expensive thing it does, and it is the only way to confirm identity
rather than trust a record.

A FLAG is not a failure: it is a place where the evidence says something the
procedure has not yet accounted for, and a human must decide.
"""

from __future__ import annotations

import hashlib
import json
from pathlib import Path

EVAL = Path(__file__).resolve().parents[1]
RUNS = EVAL / "runs"
MANIFEST_DIR = EVAL


def manifest_path(preferred: str | None = None) -> Path:
    """Newest versioned manifest by default: v1.1 outranks v1, and v1 stays on disk."""
    if preferred:
        return Path(preferred)
    candidates = sorted(MANIFEST_DIR.glob("s1-readiness-baseline*.json"))
    if not candidates:
        return MANIFEST_DIR / "s1-readiness-baseline.json"

    def version_key(path: Path) -> tuple:
        try:
            data = json.loads(path.read_text(encoding="utf-8"))
            return tuple(int(x) for x in str(data.get("version", "0")).split("."))
        except Exception:
            return (0,)

    return max(candidates, key=version_key)
READINESS = EVAL.parent / "docs" / "research" / "corpus" / "MEMORY-09-LOCAL-S1-SEQUENTIAL-READINESS-2026-11-14.md"
PLAN_DOC = EVAL.parent / "docs" / "research" / "corpus" / "MEMORY-09-LOCAL-S1-EXECUTION-PLAN-2026-11-14.md"

PASS, FAIL, FLAG = "PASS", "FAIL", "FLAG"


def dig(path: Path, chunk: int = 1 << 20) -> tuple[str, int]:
    h = hashlib.sha256()
    size = 0
    with path.open("rb") as fh:
        for block in iter(lambda: fh.read(chunk), b""):
            h.update(block)
            size += len(block)
    return h.hexdigest(), size


def latest_qual(slot: str) -> Path | None:
    dirs = sorted((p for p in RUNS.glob(f"SLOT-QUAL-{slot}-*") if p.is_dir()), key=lambda p: p.name)
    return dirs[-1] if dirs else None


def latest_l3() -> Path | None:
    dirs = sorted((p for p in RUNS.glob("L3-*") if p.is_dir()), key=lambda p: p.name)
    return dirs[-1] if dirs else None


def run(verify_files: bool = True, manifest_file: str | None = None) -> dict:
    mpath = manifest_path(manifest_file)
    manifest = json.loads(mpath.read_text(encoding="utf-8"))
    profile = manifest.get("runtime_profile")
    checks: list[dict] = []

    def add(name: str, result: str, detail: str) -> None:
        checks.append({"check": name, "result": result, "detail": detail})

    # --- 1. digests and byte counts against the manifest -------------------
    rows = []
    for entry in manifest["models"]:
        slot = entry["slot"]
        qual = latest_qual(slot)
        if qual is None:
            l3 = latest_l3() if slot == "S1-A" else None
            if l3 is None:
                rows.append((slot, "no evidence directory", None))
                continue
            record = json.loads((l3 / "record.json").read_text(encoding="utf-8"))
            declared = record["runtime"]["model_file"]["sha256"]
            declared_bytes = record["runtime"]["model_file"]["bytes"]
            source = f"L3 record {l3.name} (not a qualification record)"
        else:
            record = json.loads((qual / "record.json").read_text(encoding="utf-8"))
            declared = record["model"]["sha256"]
            declared_bytes = record["model"]["bytes"]
            source = qual.name
        path = Path(entry["path"])
        rehashed = None
        if verify_files and path.exists():
            actual, size = dig(path)
            rehashed = (actual == entry["sha256"] and size == entry["bytes"])
        rows.append((slot, source, {
            "manifest_sha": entry["sha256"],
            "declared_sha": declared,
            "declared_matches": declared == entry["sha256"],
            "declared_bytes_matches": declared_bytes == entry["bytes"],
            "file_rehash_matches": rehashed,
        }))
    for slot, source, data in rows:
        if data is None:
            add(f"{slot}: digest vs manifest", FAIL, f"no evidence at all ({source})")
        elif data["declared_matches"] and data["declared_bytes_matches"] and data["file_rehash_matches"] is not False:
            add(f"{slot}: digest vs manifest", PASS,
                f"{data['manifest_sha'][:16]}… confirmed from {source}"
                + (" (file rehashed)" if data["file_rehash_matches"] else ""))
        else:
            add(f"{slot}: digest vs manifest", FAIL, json.dumps(data))

    # --- 2. recorded configuration vs the intended runtime profile --------
    configs: dict = {}
    for entry in manifest["models"]:
        slot = entry["slot"]
        qual = latest_qual(slot)
        if qual is None:
            l3 = latest_l3() if slot == "S1-A" else None
            if l3 is None:
                add(f"{slot}: runtime configuration", FAIL, "no qualification record and no L3 record")
                continue
            record = json.loads((l3 / "record.json").read_text(encoding="utf-8"))
            configs[slot] = {"context": record["runtime"].get("max_context_length"),
                             "gpu_layers": "not recorded", "source": l3.name}
            add(f"{slot}: runtime configuration", FAIL,
                f"only L3 evidence exists (context {configs[slot]['context']}, layers "
                f"{configs[slot]['gpu_layers']}): no qualification record for this slot")
            continue
        record = json.loads((qual / "record.json").read_text(encoding="utf-8"))
        cfg = record["configuration"]
        configs[slot] = {"context": cfg["context"], "gpu_layers": cfg["gpu_layers"],
                         "quantkv": cfg["quantkv"], "source": qual.name,
                         "vram_after_mib": record["hardware"]["after_load"].get("vram_used_mib")}
        if profile:
            ok = cfg["context"] == profile["context"] and cfg["gpu_layers"] == profile["gpu_layers"]
            add(f"{slot}: runtime configuration vs profile", PASS if ok else FAIL,
                f"recorded {cfg['gpu_layers']} layers / ctx {cfg['context']}; "
                f"profile requires {profile['gpu_layers']} / {profile['context']} "
                f"({qual.name})")
        else:
            add(f"{slot}: runtime configuration", PASS,
                f"context {cfg['context']}, GPU layers {cfg['gpu_layers']}, KV {cfg['quantkv']}")

    # --- 3. offload configuration as S1 will use it -----------------------
    for entry in manifest["models"]:
        slot = entry["slot"]
        qual = latest_qual(slot)
        if qual is None:
            add(f"{slot}: offload evidence", FAIL,
                "no qualification log: the server's own offload line was never captured for this slot")
            continue
        record = json.loads((qual / "record.json").read_text(encoding="utf-8"))
        evidence = record["server"].get("log_evidence") or {}
        if evidence.get("offloaded_layers"):
            add(f"{slot}: offload evidence", PASS,
                f"{evidence['offloaded_layers']} layers, GPU {evidence.get('gpu_buffer_mib')} MiB, "
                f"CPU {evidence.get('cpu_buffer_mib')} MiB")
        else:
            add(f"{slot}: offload evidence", FAIL, "log_evidence has no offload line")

    # --- 4. zero generation records in qualification directories ----------
    gens = 0
    for entry in manifest["models"]:
        qual = latest_qual(entry["slot"])
        if qual is None:
            continue
        record = json.loads((qual / "record.json").read_text(encoding="utf-8"))
        gens += int(record.get("generations", -1))
        if (qual / "raw_body.json").exists():
            add(f"{entry['slot']}: no generation artefacts", FAIL, "raw_body.json present in a qualification dir")
    add("qualification directories: zero generations", PASS if gens == 0 else FAIL,
        f"sum of generations across filed qualification records = {gens}")

    # --- 5. recovery baselines and abort thresholds in the procedure ------
    text = READINESS.read_text(encoding="utf-8")
    needs = {
        "abort 78 °C": "abort 78 °C" in text,
        "VRAM floor 400 MiB": "400 MiB" in text,
        "RAM floor 8 GB": "8 GB" in text,
        "recovery tolerance 250 MiB": "250 MiB" in text,
        "per-component thresholds": "per-component" in text.lower(),
    }
    add("abort thresholds carried into the procedure",
        PASS if all(needs.values()) else FAIL, json.dumps(needs))
    recov = []
    for entry in manifest["models"]:
        qual = latest_qual(entry["slot"])
        if qual is None:
            recov.append((entry["slot"], "no recovery file"))
            continue
        rpath = qual / "recovery.json"
        recov.append((entry["slot"], json.loads(rpath.read_text(encoding="utf-8"))["recovery_gate"]
                      if rpath.exists() else "missing"))
    add("recovery gates recorded", PASS if all(v == "pass" for _, v in recov) else FAIL,
        json.dumps(dict(recov)))

    # --- 6. schedule and ceiling unchanged --------------------------------
    from . import plan

    caps = plan.caps_summary()
    schema_ok = caps["planned"] == 18 and caps["hard_ceiling"] == 30 and caps["retry_allowance"] == 12
    doc_ok = "18" in text and "30" in text
    add("18-call schedule, 30-call ceiling unchanged",
        PASS if schema_ok and doc_ok else FAIL,
        f"plan reports planned={caps['planned']}, retry={caps['retry_allowance']}, hard={caps['hard_ceiling']}")

    # --- 7. truncation is data; no rerun path -----------------------------
    adapter_src = (EVAL / "harness" / "local_adapter.py").read_text(encoding="utf-8")
    plan_text = PLAN_DOC.read_text(encoding="utf-8") if PLAN_DOC.exists() else ""
    corpus = text + "\n" + plan_text
    no_retry = "for attempt" not in adapter_src and "range(2)" not in adapter_src
    trunc_phrases = [
        "never retried",
        "No retry of a truncated reply",
        "Truncation is data",
        "truncation is data",
    ]
    found = [phrase for phrase in trunc_phrases if phrase in corpus]
    mentions_length = "finish_reason=length" in corpus or "finish_reason" in corpus
    trunc_rule = bool(found) and mentions_length
    add("truncation is data, no reruns", PASS if (no_retry and trunc_rule) else FAIL,
        f"rule stated as {found[:2]}; finish_reason recorded={mentions_length}; "
        f"adapter has no retry loop={no_retry}")

    failed = [c for c in checks if c["result"] == FAIL]
    flagged = [c for c in checks if c["result"] == FLAG]
    return {
        "preflight": "S1-execution-preflight-1",
        "manifest": {"path": str(mpath), "version": manifest.get("version", "1"),
                     "profile": profile},
        "checks": checks,
        "counts": {"pass": len(checks) - len(failed) - len(flagged), "fail": len(failed), "flag": len(flagged)},
        "verdict": "FAIL" if failed else ("PASS WITH FLAGS" if flagged else "PASS"),
        "configs": configs,
    }
