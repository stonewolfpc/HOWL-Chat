"""Harness CLI. `run` exists only to refuse.

Commands:
  plan      show the 18-call plan and caps (no calls are made)
  inspect   show fixture, conditions, prompt hashes, state hashes
  validate  run the offline validation suite
  run       refuse: P2 is not granted
"""

from __future__ import annotations

import argparse
import hashlib
import json
import sys
from datetime import datetime, timezone
from pathlib import Path

from . import HARNESS_VERSION
from . import authorization, budget, conditions, discover as discovery, dispatch, fixtures, l3, plan, preflight, provenance, qualify
from .adapter import ExecutionNotAuthorised, OfflineAdapter


def cmd_inspect(_: argparse.Namespace) -> int:
    fx = fixtures.build_f_st_01_s1()
    fixtures.assert_declared(fx)
    a, b = conditions.build_conditions(fx)
    print(f"harness      {HARNESS_VERSION}")
    print(f"fixture      {fx.fixture_id} v{fx.version} ({fx.dimension})")
    print(f"fixture hash {fx.sha256()}")
    print(f"state hash   {fx.state.hash()}")
    print(f"fact         {fx.fact_id} · source={fx.state.assert_fact(fx.fact_id)[1]}")
    print(f"probe        {fx.probe.user_message!r}")
    print()
    for cond in (a, b):
        print(f"condition {cond.name} · receives_state={cond.receives_state}")
        print(f"  system : {cond.system_text()!r}")
        print(f"  user   : {cond.user_text()!r}")
        print(f"  params : {json.dumps(cond.params, sort_keys=True)}")
        print(f"  prompt : {cond.prompt_hash()}")
    print()
    indices = [e["index"] for e in conditions.diff_conditions(a, b)["differing_indices"]]
    print(f"conditions differ only in the state line: {indices == [0]}")
    print(f"grant status {authorization.grant_status()}")
    return 0


def cmd_plan(_: argparse.Namespace) -> int:
    fx = fixtures.build_f_st_01_s1()
    conds = conditions.build_conditions(fx)
    print(f"caps: {json.dumps(plan.caps_summary(), indent=2)}")
    print()
    print("slots:")
    for slot in plan.default_slots():
        print(f"  {slot.name}: model={slot.model_id} provider={slot.provider}")
    print()
    try:
        calls = plan.build_plan(plan.default_slots(), conds)
    except plan.PlanIncomplete as exc:
        print(f"plan not executable: {exc}")
        print("This is the intended P1 state: naming models is an owner decision.")
        return 0
    print(f"{len(calls)} calls planned")
    return 0


def cmd_routes(_: argparse.Namespace) -> int:
    print(f"paid spend ceiling : USD {budget.MAX_PAID_SPEND_USD:.2f}")
    print(f"auto top-up         : {budget.AUTO_TOPUP_PERMITTED}")
    print(f"automatic fallback  : {budget.AUTOMATIC_FALLBACK_PERMITTED}")
    print(f"payment method      : {budget.owner_state().describe()}")
    print()
    for route, permission in budget.permission_table().items():
        print(f"  {route:26s} {permission}")
    print()
    print("hosted is REFUSED: a free account has no included monthly allowance; local routes are the supported path")
    return 0


def cmd_discover(args: argparse.Namespace) -> int:
    """L2: identity GETs only. No prompt, no model load, no server start."""
    try:
        block = discovery.discover(args.endpoint)
    except discovery.DiscoveryRefused as exc:
        print(f"REFUSED: {exc}")
        return 2
    except Exception as exc:
        print(f"discovery failed: {type(exc).__name__}: {exc}")
        return 1

    print(f"listening          : {block['listening']}")
    print(f"server version     : {block['server_version']}")
    print(f"loaded model       : {block['loaded_model']}")
    print(f"max context length : {block['max_context_length']}")
    for path, result in block["probe"].items():
        print(f"  {path:36s} http={result['status']}")
    if not block["listening"]:
        print()
        print(block["conclusion"])

    stamp = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H%M%SZ")
    out_dir = Path("eval") / "runs" / f"L2-{stamp}"
    if args.record:
        try:
            path = provenance.write_record(out_dir, block)
            print()
            print(f"identity recorded: {path}")
        except provenance.RecordError as exc:
            print(f"record not written: {exc}")
    return 0


def _file_sha256(path: Path, chunk: int = 1 << 20) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as handle:
        for block in iter(lambda: handle.read(chunk), b""):
            digest.update(block)
    return digest.hexdigest()


def cmd_l3(args: argparse.Namespace) -> int:
    """One controlled local generation. Refuses unless the L3 grant is intact."""
    try:
        grant = authorization.load_l3_grant()
    except ExecutionNotAuthorised as exc:
        print("REFUSED.")
        print(str(exc))
        return 2

    model_path = Path(grant.model_path)
    if not model_path.exists():
        print(f"REFUSED: model file not found: {model_path}")
        return 1
    print(f"model file   : {model_path}")
    print(f"hashing ({model_path.stat().st_size:,} bytes)...")
    digest = _file_sha256(model_path)
    if digest != grant.model_sha256:
        print(f"REFUSED: model sha256 {digest} does not match the granted {grant.model_sha256}")
        return 1
    print(f"model sha256 : {digest} (matches grant)")

    block = discovery.discover(grant.endpoint)
    print(f"server       : {block['server_name']} {block['server_version']} @ {grant.endpoint}")
    print(f"loaded model : {block['loaded_model']}")
    if not block["listening"]:
        print("REFUSED: nothing listening; L3 does not start a server")
        return 1
    if block["loaded_model"] in (None, "", "inactive", "unavailable"):
        print("REFUSED: no model resident; L3 does not load one by API")
        return 1

    preset = Path(args.preset) if args.preset else None
    runtime = {
        "server_name": block["server_name"],
        "server_version": block["server_version"],
        "endpoint": grant.endpoint,
        "max_context_length": block["max_context_length"],
        "loaded_model_reported": block["loaded_model"],
        "model_file": {"path": str(model_path), "bytes": model_path.stat().st_size, "sha256": digest},
        "preset": (
            {"path": str(preset), "bytes": preset.stat().st_size, "sha256": _file_sha256(preset)}
            if preset and preset.exists()
            else {"path": str(preset) if preset else UNAVAILABLE, "sha256": UNAVAILABLE}
        ),
        "koboldcpp_executable": args.executable or UNAVAILABLE,
        "identity_block_sha256": block.get("block_sha256"),
        "captured_at": provenance.utc_now(),
    }

    before = l3.collect_hardware()
    record, verification = l3.run(
        grant, model_id=args.model_id, runtime=runtime, hardware_before=before
    )

    print()
    print(f"http status  : {record['completeness']['http_status']}")
    print(f"surface      : {record['completeness']['surface']}")
    print(f"latency      : {record['completeness']['latency_ms']} ms")
    print(f"finish reason: {record['completeness']['finish_reason']} (incomplete={record['completeness']['incomplete']})")
    print(f"usage        : {json.dumps(record['completeness']['usage'])}")
    print(f"text         : {record['text']!r}")
    print(f"ledger       : attempted={record['ledger']['attempted']}/{record['ledger']['ceiling']}")
    for key, value in verification.items():
        if key != "admissible":
            print(f"  {key:24s} {value}")
    problems = verification["admissible"]
    print(f"  {'admissible':24s} {'yes' if not problems else problems}")

    if args.record:
        out_dir = Path("eval") / "runs" / f"L3-{l3.stamp()}"
        path = provenance.write_record(out_dir, record)
        ok, why = provenance.verify_record(out_dir)
        matches, mismatch = l3.parsed_matches_raw(record)
        after = l3.collect_hardware()
        (out_dir / "verification.json").write_text(
            json.dumps(
                {
                    "verification": verification,
                    "record_hash_matches": ok,
                    "record_hash_detail": why,
                    "parsed_matches_raw": matches,
                    "parsed_detail": mismatch,
                    "hardware_after": after,
                },
                indent=2,
                sort_keys=True,
            )
            + "\n",
            encoding="utf-8",
        )
        print()
        print(f"evidence     : {path}")
        print(f"write-once   : record.json + record.sha256 + raw_body.json")
        print(f"hash matches : {ok} {why}")
        print(f"parsed==raw  : {matches} {mismatch}")
    return 0


def cmd_l3_verify(args: argparse.Namespace) -> int:
    """Re-check a filed L3 record. Makes no calls and starts nothing."""
    directory = Path(args.dir)
    if not (directory / "record.json").exists():
        print(f"no record at {directory}")
        return 1
    record = json.loads((directory / "record.json").read_text(encoding="utf-8"))
    ok, why = provenance.verify_record(directory)
    matches, mismatch = l3.parsed_matches_raw(record)
    print(f"record      : {directory}")
    print(f"hash vs raw : {ok} {why}")
    print(f"parsed==raw : {matches} {mismatch}")
    print()
    failures = 0
    for row in l3.verify_l3(record):
        if row["result"] != "PASS":
            failures += 1
        print(f"  {row['result']:4s} {row['requirement']:38s} {row['detail']}")
    print()
    print(f"RESULT: {len(l3.verify_l3(record)) - failures} pass, {failures} fail")
    return 0 if failures == 0 and ok and matches else 1


def _baseline_models() -> dict:
    # Same selection rule as the preflight: the newest versioned manifest wins, and
    # older versions stay on disk as history. Hardcoding the v1 filename here meant
    # qualification kept gating on v1 hashes after v1.1 superseded them.
    path = preflight.manifest_path()
    data = json.loads(path.read_text(encoding="utf-8"))
    return {m["slot"]: m for m in data["models"]}


def cmd_qualify(args: argparse.Namespace) -> int:
    """Slot qualification: hash, identity GETs, record. Zero generations, by construction."""
    slots = _baseline_models()
    slot = args.slot.upper()
    if slot not in slots:
        print(f"unknown slot {slot!r}; baseline has {sorted(slots)}")
        return 1
    entry = slots[slot]
    model_path = Path(entry["path"])

    if not model_path.exists():
        print(f"REFUSED: model file missing: {model_path}")
        return 1

    # Gate 1: identity of the file, before anything is loaded into it.
    digest, size = qualify.file_digest(model_path)
    if digest != entry["sha256"] or size != entry["bytes"]:
        print(f"REFUSED: {slot} file does not match the frozen baseline")
        print(f"  expected {entry['sha256']} ({entry['bytes']:,} B)")
        print(f"  actual   {digest} ({size:,} B)")
        return 1
    print(f"{slot} file verified: {size:,} B  {digest}")

    # Gate 2: a model must actually be resident before identity is worth recording.
    block = discovery.discover(args.endpoint)
    if not block["listening"]:
        print("REFUSED: nothing listening; qualification does not start a server")
        return 1
    reported = block.get("loaded_model")
    print(f"server {block['server_name']} {block['server_version']}  reported model: {reported}")
    if reported in (None, "", "inactive", "unavailable"):
        print("REFUSED: no model resident")
        return 1

    folder = Path("eval") / "runs" / f"SLOT-QUAL-{slot}-{qualify.stamp()}"
    record = qualify.build_record(
        slot=slot,
        model_path=str(model_path),
        expected_sha256=entry["sha256"],
        expected_bytes=entry["bytes"],
        digest=digest,
        size=size,
        identity_block=block,
        context=args.context,
        gpu_layers=args.gpu_layers,
        quantkv=args.quantkv,
        vram_before=args.baseline_vram,
        log_path=args.log,
    )
    print(f"identity matches: {record['identity']['matches']}  ({record['identity']['reported']})")
    print(f"generations: {record['generations']}   requests: {record['requests_issued']}")
    if args.record:
        target = qualify.write_qualification(folder, record)
        print(f"qualification record: {target}")
    return 0 if record["identity"]["matches"] else 2


def cmd_qualify_recovery(args: argparse.Namespace) -> int:
    folder = Path(args.dir)
    if not (folder / "record.json").exists():
        print(f"no qualification record at {folder}")
        return 1
    reference = json.loads((folder / "record.json").read_text(encoding="utf-8"))["hardware"]["before_load"]
    target = qualify.write_recovery(folder, vram_reference=reference, port_closed=args.port_closed)
    payload = json.loads(target.read_text(encoding="utf-8"))
    print(f"recovery: {payload['recovery_gate']}  delta={payload['delta_mib']} MiB  port_closed={payload['port_closed']}")
    return 0 if payload["recovery_gate"] == "pass" and payload["port_closed"] else 2


def cmd_preflight(args: argparse.Namespace) -> int:
    """Read-only S1 preflight against filed evidence and the frozen manifest."""
    result = preflight.run(verify_files=not args.no_hash, manifest_file=args.manifest)
    for row in result["checks"]:
        print(f"  {row['result']:4s} {row['check']:46s} {row['detail']}")
    counts = result["counts"]
    print()
    print(f"counts: {counts['pass']} pass, {counts['fail']} fail, {counts['flag']} flag")
    print(f"VERDICT: {result['verdict']}")
    if args.record:
        folder = Path("eval") / "runs" / f"PREFLIGHT-S1-{qualify.stamp()}"
        folder.mkdir(parents=True, exist_ok=True)
        target = folder / "preflight.json"
        if target.exists():
            print(f"refusing to overwrite {target}")
            return 2
        target.write_text(json.dumps(result, indent=2, sort_keys=True) + "\n", encoding="utf-8")
        print(f"preflight evidence: {target}")
    return 0 if counts["fail"] == 0 else 1


def cmd_validate(_: argparse.Namespace) -> int:
    from .validate import main as validate_main

    return validate_main()


def cmd_run(_: argparse.Namespace) -> int:
    try:
        authorization.require_grant()
    except ExecutionNotAuthorised as exc:
        print("REFUSED.")
        print(str(exc))
        print()
        print("P1 authorises building an instrument, not running it. Nothing was sent, spent or loaded.")
        return 2
    # Reachable only with a signed grant. The dispatcher is the sole counting
    # boundary, so `run` goes through it even though no adapter can answer in P1.
    ledger = dispatch.CallLedger(ceiling=grant.call_ceiling)
    dispatcher = dispatch.Dispatcher(ledger, OfflineAdapter())
    print(f"grant present (ceiling {grant.call_ceiling}); dispatch path wired through the ledger")
    print(f"ledger: attempted={ledger.attempted} remaining={ledger.remaining()}")
    print("P1 adapters cannot produce a response; stopping before any dispatch.")
    return 3


def build_parser() -> argparse.ArgumentParser:
    p = argparse.ArgumentParser(prog="memory09-harness", description="MEMORY-09 P1 harness (offline)")
    sub = p.add_subparsers(dest="command", required=True)
    sub.add_parser("inspect", help="show fixture, conditions and hashes").set_defaults(fn=cmd_inspect)
    sub.add_parser("plan", help="show the call plan and caps").set_defaults(fn=cmd_plan)
    sub.add_parser("validate", help="run offline validation checks").set_defaults(fn=cmd_validate)
    sub.add_parser("routes", help="permission table and spend ceiling").set_defaults(fn=cmd_routes)
    disc = sub.add_parser("discover", help="L2 identity GETs (no prompt, no generation)")
    disc.add_argument("--endpoint", default=discovery.DEFAULT_ENDPOINT, help="loopback base URL")
    disc.add_argument("--record", action="store_true", help="write the identity block under eval/runs")
    disc.set_defaults(fn=cmd_discover)
    l3p = sub.add_parser("l3", help="one controlled local generation (needs the L3 grant)")
    l3p.add_argument("--record", action="store_true", help="write evidence under eval/runs")
    l3p.add_argument("--model-id", default="qwen3.5-9b-q4_k_m", help="model label sent to the server")
    l3p.add_argument("--preset", default=r"H:\AI apps\qwen 9b.kcpps", help="preset file used for provenance")
    l3p.add_argument("--executable", default=r"H:\AI apps\koboldcpp.exe", help="server executable")
    l3p.set_defaults(fn=cmd_l3)
    ver = sub.add_parser("l3-verify", help="re-check a filed L3 record (no calls)")
    ver.add_argument("--dir", required=True, help="an eval/runs/L3-* directory")
    ver.set_defaults(fn=cmd_l3_verify)
    q = sub.add_parser("qualify", help="slot qualification: identity only, zero generations")
    q.add_argument("--slot", required=True, help="S1-A / S1-B / S1-C")
    q.add_argument("--context", type=int, default=8192)
    q.add_argument("--gpu-layers", type=int, required=True)
    q.add_argument("--quantkv", default="q8_0")
    q.add_argument("--endpoint", default=discovery.DEFAULT_ENDPOINT)
    q.add_argument("--baseline-vram", type=json.loads, default={}, help='JSON, e.g. {"vram_used_mib": "1160 MiB"}')
    q.add_argument("--log", default=None, help="server log to extract load evidence from")
    q.add_argument("--record", action="store_true")
    q.set_defaults(fn=cmd_qualify)
    qr = sub.add_parser("qualify-recovery", help="record recovery after the server is stopped")
    qr.add_argument("--dir", required=True)
    qr.add_argument("--port-closed", action="store_true")
    qr.set_defaults(fn=cmd_qualify_recovery)
    pf = sub.add_parser("preflight-s1", help="read-only S1 preflight against filed evidence")
    pf.add_argument("--record", action="store_true")
    pf.add_argument("--no-hash", action="store_true", help="skip rehashing model files")
    pf.add_argument("--manifest", default=None, help="baseline manifest (default: newest version)")
    pf.set_defaults(fn=cmd_preflight)
    sub.add_parser("run", help="refuse: P2 not granted").set_defaults(fn=cmd_run)
    return p


def main(argv: list[str] | None = None) -> int:
    args = build_parser().parse_args(argv)
    return args.fn(args)


if __name__ == "__main__":
    sys.exit(main())
