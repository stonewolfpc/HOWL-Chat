"""R1–R5 paired adversarial set.

Each case is (query, candidate records, expected outcome). Candidates are passed
explicitly so ownership/history filters are not the thing under test — the
reranker is. Outcomes are the registered criteria from
MEMORY-09-S3-RERANKER-SCOPE-2026-11-15.md.
"""
from __future__ import annotations

from eval.s3.records import Record

SABLE = "char-001"
ALARIC = "char-002"


def _cand(r: Record) -> dict:
    return {
        "record_id": r.record_id,
        "owner_id": r.owner_id,
        "status": r.status,
        "turn": r.established_turn,
        "text": f"{r.subject} {r.predicate}: {r.value}",
    }


def build_cases(store) -> list[dict]:
    by_id = {r.record_id: r for r in store.all()}
    c = lambda *ids: [_cand(by_id[i]) for i in ids]

    sable_pool = c("rec-0001", "rec-0003", "rec-0006", "rec-0007")
    return [
        # R1 — agent vs patient. "horse" is a trap: it shares "trained" with rec-0001
        # but no record is about a horse, so the only correct outcome is abstention.
        {"id": "R1-horse", "criterion": "R1", "owner": SABLE,
         "query": "Who trained the horse?", "candidates": sable_pool,
         "expect": None, "expect_abstain": True,
         "why": "shares 'trained'; different patient entirely"},
        # R1 — the same relation read the other way round must still find the record.
        {"id": "R1-patient", "criterion": "R1", "owner": SABLE,
         "query": "Was Corvin the one who trained you?", "candidates": sable_pool,
         "expect": "rec-0001", "expect_abstain": False,
         "why": "agent named explicitly; relation is trained_by"},

        # R2 — ownership. Vault code belongs to Alaric; Sable's pool must not answer.
        {"id": "R2-vault-foreign", "criterion": "R2", "owner": SABLE,
         "query": "What is the vault code?", "candidates": sable_pool,
         "expect": None, "expect_abstain": True,
         "why": "record belongs to another character"},
        {"id": "R2-vault-own", "criterion": "R2", "owner": ALARIC,
         "query": "What is the vault code?",
         "candidates": c("rec-0004", "rec-0005"),
         "expect": "rec-0004", "expect_abstain": False,
         "why": "same query, owning character must answer"},

        # R3 — history vs current, both directions.
        {"id": "R3-now", "criterion": "R3", "owner": SABLE,
         "query": "Where is the ledger kept now?",
         "candidates": c("rec-0002", "rec-0007"),
         "expect": "rec-0007", "expect_abstain": False,
         "why": "current record must beat the superseded one"},
        {"id": "R3-before", "criterion": "R3", "owner": SABLE,
         "query": "Where was the ledger kept before?",
         "candidates": c("rec-0002", "rec-0007"),
         "expect": "rec-0002", "expect_abstain": False,
         "why": "history stays reachable when explicitly asked"},

        # R4 — veto: paraphrase must still resolve to the right record.
        {"id": "R4-paraphrase", "criterion": "R4", "owner": SABLE,
         "query": "Who was your trainer, again?", "candidates": sable_pool,
         "expect": "rec-0001", "expect_abstain": False,
         "why": "the P1 behaviour that must not regress"},

        # R5 — veto: irrelevant queries select nothing.
        {"id": "R5-weather", "criterion": "R5", "owner": SABLE,
         "query": "What did you think of the weather?", "candidates": sable_pool,
         "expect": None, "expect_abstain": True,
         "why": "no record concerns weather"},
        {"id": "R5-fair", "criterion": "R5", "owner": SABLE,
         "query": "What did you win at the fair?", "candidates": sable_pool,
         "expect": None, "expect_abstain": True,
         "why": "no record concerns a fair or a prize"},
    ]


CRITERIA = {
    "R1": "agent vs patient",
    "R2": "subject ownership",
    "R3": "historical vs current",
    "R4": "preserve P1 (VETO)",
    "R5": "preserve abstention (VETO)",
}
VETOES = ("R4", "R5")
