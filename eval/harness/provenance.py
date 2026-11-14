"""Provenance and evidence records (row 477 v1.1 §8; row 482 §7).

Every field the provider does not supply is stored as `unavailable` — never
inferred, never omitted. A record missing a field *without* that marker is
inadmissible.
"""

from __future__ import annotations

import json
from datetime import datetime, timezone
from pathlib import Path

from . import HARNESS_VERSION
from .adapter import (
    HARNESS_ADAPTER_NAME,
    HARNESS_ADAPTER_VERSION,
    PROVENANCE_FIELDS,
    GenerateRequest,
    GenerateResponse,
    UNAVAILABLE,
    canonical,
    sha256_text,
)

# Fields whose absence must be marked, not assumed.
MARKABLE = tuple(f for f in PROVENANCE_FIELDS if f not in ("adapter_name", "adapter_version"))

REQUIRED_NON_MARKABLE = ("adapter_name", "adapter_version", "test_id", "fixture_id", "prompt_hash")


EVAL_ROOT = Path(__file__).resolve().parents[1]


class RecordError(ValueError):
    """Raised when a record cannot be admitted, overwritten, or written outside eval/."""


def utc_now() -> str:
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def build_record(
    *,
    test_id: str,
    fixture_id: str,
    fixture_version: str,
    request: GenerateRequest,
    model_revision: str | None = None,
    response: GenerateResponse | None = None,
    state_hash_before: str | None = None,
    state_hash_after: str | None = None,
    provider_used: str | None = None,
    client_version: str | None = None,
    provider_headers: dict | None = None,
    seed_honoured: str | None = None,
) -> dict:
    """Assemble a record; `response` is absent during P1 by definition."""
    prompt_hash = request.prompt_hash()
    record = {
        "adapter_name": HARNESS_ADAPTER_NAME,
        "adapter_version": HARNESS_ADAPTER_VERSION,
        "client_version": client_version,
        "model_requested": request.model_id,
        "model_revision": model_revision,
        "provider_requested": request.provider,
        "provider_used": provider_used,
        "returned_model": (response.model_reported if response else None),
        "system_fingerprint": (response.system_fingerprint if response else None),
        "provider_headers": provider_headers,
        "params_sent": dict(request.params),
        "seed": request.seed,
        "seed_honoured": seed_honoured,
        "prompt_hash": prompt_hash,
        "response_hash": (sha256_text(canonical(response.raw_body)) if response and response.raw_body else None),
        "state_hash_before": state_hash_before,
        "state_hash_after": state_hash_after,
        "fixture_id": fixture_id,
        "fixture_version": fixture_version,
        "test_id": test_id,
        "harness_version": HARNESS_VERSION,
        "timestamp": (response.timestamp if response and response.timestamp else utc_now()),
        "raw_body": (response.raw_body if response else None),
        "text": (response.text if response else None),
        "request_echo": request.echo(),
        "error": (response.error if response else None),
    }
    if response is not None:
        record["provider_used"] = provider_used or response.provider_used
    return mark_unavailable(record)


def mark_unavailable(record: dict) -> dict:
    """Replace absent/None markable fields with the literal marker."""
    out = dict(record)
    for field_name in MARKABLE:
        if field_name not in out or out[field_name] is None:
            out[field_name] = UNAVAILABLE
    return out


def check_admissible(record: dict) -> list[str]:
    """Return the list of reasons this record cannot be admitted."""
    problems: list[str] = []
    for field_name in PROVENANCE_FIELDS:
        if field_name not in record:
            problems.append(f"missing field: {field_name}")
    for field_name in REQUIRED_NON_MARKABLE:
        if record.get(field_name) in (None, "", UNAVAILABLE):
            problems.append(f"required field is {record.get(field_name)!r}: {field_name}")
    if record.get("raw_body") in (None, UNAVAILABLE):
        problems.append("raw_body is missing or unavailable: text alone is not evidence")
    if record.get("provider_used") in (None, "", UNAVAILABLE):
        problems.append("provider_used not recorded")
    if record.get("provider_requested") in (None, "", UNAVAILABLE):
        problems.append("provider_requested not recorded")
    for key in ("prompt_hash", "state_hash_before", "state_hash_after"):
        value = record.get(key)
        if value in (None, UNAVAILABLE):
            problems.append(f"{key} not recorded")
    if record.get("seed") not in (None, UNAVAILABLE) and record.get("seed_honoured") in (None, UNAVAILABLE):
        problems.append("seed sent but seed_honoured not stated")
    return problems


def write_record(out_dir: str | Path, record: dict) -> Path:
    """Write evidence **once**, and only inside `eval/`.

    Evidence that can be overwritten is not evidence. A second write to the same
    directory raises instead of replacing the earlier record; the record's own
    hash goes to a sidecar so later readers can detect edits.
    """
    directory = Path(out_dir).resolve()
    if not str(directory).startswith(str(EVAL_ROOT.resolve())):
        raise RecordError(f"refusing to write evidence outside eval/: {directory}")
    record_path = directory / "record.json"
    if record_path.exists():
        raise RecordError(
            f"evidence is write-once: {record_path} already exists and will not be replaced"
        )
    directory.mkdir(parents=True, exist_ok=True)
    payload = json.dumps(record, indent=2, sort_keys=True, ensure_ascii=False) + "\n"
    record_path.write_text(payload, encoding="utf-8")
    (directory / "record.sha256").write_text(
        sha256_text(canonical(record)) + "  record.json\n", encoding="utf-8"
    )
    raw = record.get("raw_body")
    if raw not in (None, UNAVAILABLE):
        (directory / "raw_body.json").write_text(
            json.dumps(raw, indent=2, sort_keys=True, ensure_ascii=False) + "\n", encoding="utf-8"
        )
    return record_path


def verify_record(out_dir: str | Path) -> tuple[bool, str]:
    """Re-hash a written record against its sidecar."""
    directory = Path(out_dir)
    record_path = directory / "record.json"
    sidecar = directory / "record.sha256"
    if not record_path.exists() or not sidecar.exists():
        return (False, "record or sidecar missing")
    record = json.loads(record_path.read_text(encoding="utf-8"))
    expected = sidecar.read_text(encoding="utf-8").split()[0]
    actual = sha256_text(canonical(record))
    return (actual == expected, "" if actual == expected else f"record edited: {actual} != {expected}")
