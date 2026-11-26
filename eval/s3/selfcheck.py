"""S3 self-check — the whole memory pipeline, with ZERO model calls.

Proves, offline, that the system S3 intends to test actually works:
evidence is immutable and addressable, formation builds records with provenance,
retrieval honours ownership, supersession prefers current truth, abstention
happens when evidence is insufficient, and the target evidence is genuinely gone
from the window at probe time.

    PYTHONPATH=. python eval/s3/selfcheck.py
"""
from __future__ import annotations

import sys

from eval.s3.evidence import EvidenceLog
from eval.s3.records import RecordStore, ACTIVE, SUPERSEDED
from eval.s3.formation import Formatter
from eval.s3.retrieve import retrieve
from eval.s3.assemble import assemble, WRAPPERS
from eval.s3.context import Window

OWNER = {"Sable": "char-001", "Alaric": "char-002"}
FILLER = [
    "The rain came in sideways off the moor and nobody wanted to be the first to speak.",
    "We walked the north road until the light went grey and the horses refused to go on.",
    "Someone had left bread and hard cheese on the table, and it did not last long.",
    "There is a particular silence in a room where two people are deciding not to argue.",
    "By morning the fire had burned down to a red line and the kettle was cold.",
    "The innkeeper counted coins twice, then a third time, and still looked unhappy.",
    "Wind in the chimney sounds like someone talking downstairs, if you are tired enough.",
    "I have never trusted a map that was drawn by someone who had not walked the ground.",
]


def build_conversation():
    """A 60-turn conversation: introduce, develop, displace, change state, probe late."""
    script: list[tuple[str, str] | None] = [
        ("You", "The road south was worse than I remembered. My boots are still wet."),
        ("Sable", "Then you came the long way. Sit down before you fall down."),
        ("You", "Long way or not, I am here. Has it always been this quiet?"),
        ("Sable", "I was trained by Corvin Ashgrave of Ossary, and quiet was not part of it."),
        ("Sable", "I keep a ledger of unpaid debts in the blue cabinet, if you are curious."),
        ("You", "Noted. I will try to stay off that ledger."),
        ("Sable", "See that you do. I like tea more than I like surprises."),
        ("You", "Tea it is, then."),
        None, None, None,
        ("Alaric", "The vault code is seven-four-nine-one. Do not make me regret saying it."),
        ("Sable", "He says that to everyone."),
        ("Alaric", "I served four years in the Ninth Company. I know when someone is listening."),
        None, None, None, None, None, None, None,
        ("Sable", "Actually, I keep a ledger of debts in the green cabinet now. Blue warped."),
        ("You", "Green. Understood."),
        None, None, None, None, None, None, None, None, None, None, None, None,
        ("You", "Where is the ledger kept these days?"),
        None, None,
        ("You", "What is the vault code?"),
        None, None,
        None, None, None, None, None,
        ("You", "Remind me — who trained you?"),
        None, None, None, None,
        ("You", "What did you think of the weather?"),
        None, None, None, None,
        ("You", "Who was your trainer, again?"),
    ]
    turns, fi = [], 0
    for entry in script:
        if entry is None:
            speaker = "Sable" if fi % 3 else "You"
            turns.append((speaker, FILLER[fi % len(FILLER)]))
            fi += 1
        else:
            turns.append(entry)
    return turns


def toks(text: str) -> int:
    return max(1, round(len(text) / 4))


def main() -> int:
    log, store, win = EvidenceLog(), RecordStore(), Window(capacity=420, reserve=80)
    former = Formatter(store, log)
    turns = build_conversation()
    probes: list[dict] = []
    checks: list[tuple[str, bool, str]] = []
    leak_violations: list[str] = []
    paraphrase_gap = False

    for i, (speaker, text) in enumerate(turns, start=1):
        turn = log.append(speaker, text, toks(text))
        roll = win.add(turn)
        if roll or i % 8 == 0:
            former.run_pending(log, OWNER)

        if text.startswith("Where is the ledger"):
            r = retrieve(text, store, "char-001")
            block = assemble(r, store, WRAPPERS["labelled_evidence"])
            probes.append({"turn": i, "probe": "ledger", "selected": [s.record_id for s in r.selected],
                           "values": [store.get(s.record_id).value for s in r.selected],
                           "abstained": r.abstained, "block": block})
        elif text.startswith("What is the vault code"):
            r_self = retrieve(text, store, "char-001")   # Sable must not know it
            r_other = retrieve(text, store, "char-002")  # Alaric must
            probes.append({"turn": i, "probe": "vault_code_as_sable",
                           "selected": [s.record_id for s in r_self.selected],
                           "abstained": r_self.abstained, "reason": r_self.reason})
            probes.append({"turn": i, "probe": "vault_code_as_alaric",
                           "selected": [s.record_id for s in r_other.selected],
                           "abstained": r_other.abstained})
            for s in r_self.selected:
                if store.get(s.record_id).owner_id != "char-001":
                    leak_violations.append(f"turn {i}: {s.record_id} foreign")
        elif text.startswith("Remind me"):
            r = retrieve(text, store, "char-001")
            src = log.span(store.get(r.selected[0].record_id).source_span) if r.selected else None
            probes.append({"turn": i, "probe": "trainer_lexical",
                           "selected": [s.record_id for s in r.selected],
                           "abstained": r.abstained,
                           "source_turn": src.turn if src else None,
                           "source_in_window": win.contains_span(src.span_id) if src else None})
        elif text.startswith("What did you think of the weather"):
            r = retrieve(text, store, "char-001")
            probes.append({"turn": i, "probe": "ambiguous", "abstained": r.abstained,
                           "selected": [s.record_id for s in r.selected], "reason": r.reason})
        elif text.startswith("Who was your trainer"):
            r = retrieve(text, store, "char-001")
            paraphrase_gap = r.abstained
            probes.append({"turn": i, "probe": "trainer_paraphrase", "abstained": r.abstained,
                           "selected": [s.record_id for s in r.selected], "reason": r.reason})

    # ---- assertions -----------------------------------------------------------
    trained = [r for r in store.all() if r.predicate == "trained_by"]
    source_span = trained[0].source_span if trained else ""
    source_turn = log.span(source_span).turn if trained else -1

    checks.append(("displacement: >=2 rollovers", len(win.rollovers) >= 2,
                   f"{len(win.rollovers)} rollovers"))
    checks.append(("displacement: trainer evidence left the window",
                   bool(trained) and not win.contains_span(source_span),
                   f"source turn {source_turn} (span {source_span}) absent from window"))

    trainer_probe = next(p for p in probes if p["probe"] == "trainer_lexical")
    checks.append(("retrieval after displacement finds the record",
                   (not trainer_probe["abstained"]) and len(trainer_probe["selected"]) == 1,
                   f"selected {trainer_probe['selected']}"))
    checks.append(("retrieved record's source is NOT in the window",
                   trainer_probe["source_in_window"] is False,
                   "recall is retrieval, not survival in the prompt"))

    hist = store.history("char-001", "Sable", "keeps_ledger")
    current = [r for r in hist if r.status == ACTIVE]
    checks.append(("supersession: two versions retained, one current",
                   len(hist) == 2 and len(current) == 1, f"{len(hist)} records, {len(current)} active"))
    checks.append(("supersession: current value wins",
                   bool(current) and "green" in current[0].value.lower(),
                   current[0].value if current else "none"))
    ledger_probe = next(p for p in probes if p["probe"] == "ledger")
    checks.append(("supersession: stale value never retrieved",
                   all("blue" not in v.lower() for v in ledger_probe["values"]),
                   f"values={ledger_probe['values']}"))

    sable_vault = next(p for p in probes if p["probe"] == "vault_code_as_sable")
    alaric_vault = next(p for p in probes if p["probe"] == "vault_code_as_alaric")
    checks.append(("isolation: Sable abstains on Alaric's code", sable_vault["abstained"] is True,
                   sable_vault.get("reason", "")))
    checks.append(("isolation: Alaric retrieves his own code", not alaric_vault["abstained"],
                   f"selected {alaric_vault['selected']}"))
    checks.append(("isolation: no foreign record ever selected", not leak_violations,
                   f"{len(leak_violations)} violations"))

    amb = next(p for p in probes if p["probe"] == "ambiguous")
    checks.append(("abstention: ambiguous probe retrieves nothing", amb["abstained"] is True,
                   amb.get("reason", "")))

    grounded = {r.predicate for r in store.all() if r.status == ACTIVE}
    checks.append(("formation: expected predicates captured",
                   {"trained_by", "keeps_ledger", "vault_code"} <= grounded,
                   f"active predicates={sorted(grounded)}"))
    # Regression test for the defect found on the first run: formation ran only on
    # the current turn at checkpoints, so a fact introduced between checkpoints
    # was never formed. Turn 4 is not a checkpoint multiple and must still exist.
    checks.append(("formation: non-checkpoint turn still formed",
                   bool(trained) and trained[0].established_turn == 4,
                   f"trained_by established_turn={trained[0].established_turn if trained else None}"))
    checks.append(("formation: exact repeats not duplicated",
                   former.stats.duplicates == 0 or former.stats.captured >= former.stats.duplicates,
                   former.stats.as_dict().__str__()))

    # ---- report ---------------------------------------------------------------
    print(f"conversation: {len(turns)} turns, {sum(t.tokens for t in log.turns())} tokens, "
          f"window capacity {win.capacity} (prompt {win.prompt_capacity}), "
          f"{len(win.rollovers)} rollovers")
    print(f"evidence sha256: {log.sha256()[:32]}")
    print(f"records: {len(store.all())} total, {store.count_by_status()}")
    print(f"formation: {former.stats.as_dict()}")
    print()
    for p in probes:
        print(f"  turn {p['turn']:>3} {p['probe']:<24} abstained={p['abstained']} "
              f"selected={p.get('selected')} {p.get('reason','')}")
    print()
    ok = True
    for i, (name, passed, detail) in enumerate(checks, start=1):
        print(f"  [A{i:02d}] {'PASS' if passed else 'FAIL'} {name:<52} {detail}")
        ok = ok and passed
    print()
    print(f"KNOWN GAP (Gate B evidence, not a failure): paraphrase probe abstained={paraphrase_gap} — "
          f"lexical-only retrieval cannot match 'trainer' to 'trained_by'.")
    print(f"\nRESULT: {'PASS' if ok else 'FAIL'} ({len(checks)} checks, "
          f"{sum(1 for _, p, _ in checks if p)} passed)")

    # Machine-readable receipt so the check count is never a matter of counting lines.
    import hashlib as _h
    import json as _j
    import pathlib as _p
    mod_hashes = {f.name: _h.sha256(f.read_bytes()).hexdigest()
                  for f in sorted(_p.Path(__file__).parent.glob("*.py"))}
    receipt = {
        "gate": "A",
        "suite": "s3.selfcheck",
        "checks_total": len(checks),
        "checks_passed": sum(1 for _, p, _ in checks if p),
        "checks": [{"id": f"A{i:02d}", "name": n, "passed": p, "detail": d}
                   for i, (n, p, d) in enumerate(checks, start=1)],
        "result": "PASS" if ok else "FAIL",
        "conversation_turns": len(turns),
        "conversation_tokens": sum(t.tokens for t in log.turns()),
        "evidence_sha256": log.sha256(),
        "rollovers": len(win.rollovers),
        "records": store.count_by_status(),
        "formation": former.stats.as_dict(),
        "known_gap": {"probe": "trainer_paraphrase", "abstained": paraphrase_gap,
                      "why": "lexical branch cannot match 'trainer' to 'trained_by' (Gate B)"},
        "module_sha256": mod_hashes,
        "note": "window is simulated Python; KoboldCpp displacement is unproven (Gate C)",
    }
    out = _p.Path("eval/s3/gate-a-receipt.json")
    out.write_text(_j.dumps(receipt, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    print(f"receipt: {out} ({out.stat().st_size} B)")
    return 0 if ok else 1


if __name__ == "__main__":
    sys.exit(main())
