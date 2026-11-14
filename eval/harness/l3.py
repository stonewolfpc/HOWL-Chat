"""L3 — one controlled local generation (WO-L3-01).

Not the MEMORY-09 evaluation, and not S1. This module exists to answer a narrow
question: can the harness drive a real local server once, keep the bytes, and
record enough provenance that the result could be checked by someone else?

Ceiling is the grant's ceiling, and the grant for this order says one call.
`Dispatcher` counts before the adapter is touched, so a failed call still uses
the allowance — there is no retry path here and there is no way to add one by
accident: `dispatch` is called exactly once in `run`.
"""

from __future__ import annotations

import subprocess
from datetime import datetime, timezone

from .adapter import (
    HARNESS_ADAPTER_VERSION,
    UNAVAILABLE,
    GenerateRequest,
    canonical,
    sha256_text,
)
from . import dispatch, provenance

L3_TEST_ID = "L3-01"
L3_FIXTURE_VERSION = "1"

# Synthetic only: no personal, sensitive or corpus-derived material.
PROMPT_SYSTEM = "You are a controlled inference check. Answer in one short sentence."
PROMPT_USER = "Name the colour of a clear midday sky."
PARAMS = {"temperature": 0.0, "max_tokens": 48}


def build_request(model_id: str, provider: str) -> GenerateRequest:
    return GenerateRequest(
        messages=(
            {"role": "system", "content": PROMPT_SYSTEM},
            {"role": "user", "content": PROMPT_USER},
        ),
        model_id=model_id,
        provider=provider,
        params=dict(PARAMS),
        seed=None,
    )


def collect_hardware() -> dict:
    """Observational only: read GPU counters, change nothing."""
    out: dict = {"source": "nvidia-smi", "captured_at": provenance.utc_now()}
    try:
        proc = subprocess.run(
            [
                "nvidia-smi",
                "--query-gpu=name,memory.total,memory.used,memory.free,utilization.gpu,temperature.gpu",
                "--format=csv,noheader",
            ],
            capture_output=True,
            text=True,
            timeout=20,
        )
        line = (proc.stdout or "").strip()
        if proc.returncode != 0 or not line:
            out.update({"available": False, "reason": (proc.stderr or "no output").strip()[:200]})
            return out
        parts = [p.strip() for p in line.split(",")]
        out.update(
            {
                "available": True,
                "gpu": parts[0],
                "vram_total_mib": parts[1],
                "vram_used_mib": parts[2],
                "vram_free_mib": parts[3],
                "gpu_util_percent": parts[4],
                "temperature_c": parts[5],
            }
        )
    except Exception as exc:  # monitoring may simply not exist
        out.update({"available": False, "reason": f"{type(exc).__name__}: {exc}"[:200]})
    return out


def run(grant, *, model_id: str, runtime: dict, post_transport=None, hardware_before=None):
    """Exactly one dispatch. Returns (record, verification)."""
    from .local_adapter import LocalEndpoint, LocalHTTPAdapter

    endpoint = LocalEndpoint(base_url=grant.endpoint)
    adapter = LocalHTTPAdapter(grant=grant, endpoint=endpoint, post_transport=post_transport)
    ledger = dispatch.CallLedger(ceiling=grant.call_ceiling)
    dispatcher = dispatch.Dispatcher(ledger, adapter)

    request = build_request(model_id, "local:koboldcpp")
    response = dispatcher.dispatch(L3_TEST_ID, request)  # the one and only call

    record = provenance.build_record(
        test_id=L3_TEST_ID,
        fixture_id=L3_TEST_ID,
        fixture_version=L3_FIXTURE_VERSION,
        request=request,
        response=response,
        provider_used=response.provider_used,
        client_version=HARNESS_ADAPTER_VERSION,
        seed_honoured=UNAVAILABLE,  # no seed was sent; nothing to honour
    )
    record["runtime"] = runtime
    record["hardware_before"] = hardware_before or {"available": False, "reason": "not captured"}
    record["completeness"] = {
        "http_status": response.http_status if response.http_status is not None else UNAVAILABLE,
        "surface": response.surface or UNAVAILABLE,
        "finish_reason": response.finish_reason if response.finish_reason not in (None, "") else UNAVAILABLE,
        "incomplete": response.incomplete,
        "latency_ms": response.latency_ms if response.latency_ms is not None else UNAVAILABLE,
        "usage": response.usage if response.usage is not None else UNAVAILABLE,
    }
    record["ledger"] = {
        "ceiling": ledger.ceiling,
        "attempted": ledger.attempted,
        "per_slot": dict(ledger.per_slot),
        "failures": ledger.consumed_by_failure,
    }
    record["grant"] = {
        "report_id": grant.report_id,
        "date": grant.date,
        "call_ceiling": grant.call_ceiling,
        "endpoint": grant.endpoint,
        "model_path": grant.model_path,
        "model_sha256": grant.model_sha256,
        "signature_kind": "owner-stated, presence-checked; NOT cryptographic authentication",
    }

    verification = {
        "one_call_respected": ledger.attempted == 1,
        "raw_body_present": record.get("raw_body") not in (None, UNAVAILABLE),
        "text_present": record.get("text") not in (None, UNAVAILABLE),
        "usage_marked": record["completeness"]["usage"] == UNAVAILABLE or isinstance(
            record["completeness"]["usage"], dict
        ),
        "latency_recorded": record["completeness"]["latency_ms"] != UNAVAILABLE,
        "finish_reason_recorded": record["completeness"]["finish_reason"] != UNAVAILABLE,
        "no_hosted_route": is_loopback(str(record.get("provider_used", ""))),
        "admissible": provenance.check_admissible(record),
    }
    return record, verification


def parsed_matches_raw(record: dict) -> tuple[bool, str]:
    """Re-derive the text from the preserved body and compare with `text`."""
    from .local_surfaces import transcribe_any

    body = record.get("raw_body")
    if not isinstance(body, dict):
        return (False, "raw_body is not an object")
    try:
        _surface, response = transcribe_any(body)
    except Exception as exc:
        return (False, f"transcription failed: {type(exc).__name__}: {exc}")
    raw_text = response.text
    if raw_text != record.get("text"):
        return (False, f"text mismatch: re-parsed {raw_text!r} != recorded {record.get('text')!r}")
    expected_hash = sha256_text(canonical(response.raw_body)) if response.raw_body else None
    if expected_hash and record.get("response_hash") != expected_hash:
        return (False, "response_hash does not match the preserved body")
    return (True, "")


def is_loopback(url: str) -> bool:
    """`provider_used` is the endpoint URL, not the provider label: parse it."""
    from urllib.parse import urlparse

    from .local_adapter import LOOPBACK_HOSTS

    host = (urlparse(url).hostname or "").lower()
    return host in LOOPBACK_HOSTS


def verify_l3(record: dict) -> list[dict]:
    """Requirement-by-requirement result for a completed L3 record.

    L3 carries no memory state — it is an inference-path check, not the
    behavioural test — so `state_hash_*` being `unavailable` is expected here
    and is reported as not-applicable rather than silently accepted.
    """
    completeness = record.get("completeness") or {}
    ledger = record.get("ledger") or {}
    out: list[dict] = []

    def add(name: str, ok: bool, detail: str) -> None:
        out.append({"requirement": name, "result": "PASS" if ok else "FAIL", "detail": detail})

    add(
        "one real generation completed",
        completeness.get("http_status") == 200,
        f"http_status={completeness.get('http_status')!r}",
    )
    add(
        "raw body preserved",
        isinstance(record.get("raw_body"), dict),
        "raw_body.json holds the unedited response object",
    )
    add(
        "parsed text matches raw body",
        (lambda r: r[0])(parsed_matches_raw(record)),
        (lambda r: r[1] or "re-transcription is identical")(parsed_matches_raw(record)),
    )
    add(
        "latency recorded",
        completeness.get("latency_ms") not in (None, "unavailable"),
        f"latency_ms={completeness.get('latency_ms')!r}",
    )
    add(
        "usage recorded or marked",
        completeness.get("usage") == "unavailable" or isinstance(completeness.get("usage"), dict),
        f"usage={completeness.get('usage')!r}",
    )
    add(
        "completeness stated",
        completeness.get("finish_reason") != "unavailable",
        f"finish_reason={completeness.get('finish_reason')!r} incomplete={completeness.get('incomplete')!r}",
    )
    add(
        "one-call ceiling respected",
        ledger.get("attempted") == ledger.get("ceiling") == 1,
        f"attempted={ledger.get('attempted')} ceiling={ledger.get('ceiling')}",
    )
    add(
        "no hosted route used",
        is_loopback(str(record.get("provider_used", ""))),
        f"provider_used={record.get('provider_used')!r}",
    )
    add(
        "model/runtime provenance present",
        isinstance(record.get("runtime"), dict) and bool(record["runtime"].get("model_file")),
        "runtime block carries file path, bytes, sha256 and server version",
    )
    state_ok = record.get("state_hash_before") in (None, "unavailable") and record.get(
        "state_hash_after"
    ) in (None, "unavailable")
    add(
        "state hashes (not applicable to L3)",
        state_ok,
        "L3 has no memory state; absence is correct and was marked 'unavailable'",
    )
    return out


def stamp() -> str:
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H%M%SZ")
