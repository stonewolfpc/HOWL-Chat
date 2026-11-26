"""Experiment E-TEMPORAL-SET — unseen temporal cases on a fresh record store.

Completely independent of the frozen seven records and the R1–R5 cases. Nothing
here edits, extends or reinterprets them.

Coverage required before the run:
  * multiple successive updates to one fact (three-deep chain)
  * corrections versus ordinary updates
  * conflicting records with different establishment turns, neither superseded
  * queries that explicitly request historical state
  * queries with no temporal language that must resolve to current state
  * records owned by different characters, so ownership must still hold
  * one chat-scope record, to record the known scope gap WITHOUT scoring it
"""
from __future__ import annotations

from eval.s3.records import RecordStore

MIRA, TEODOR, CHAT = "char-101", "char-102", "chat-001"


def _rec(store, owner, owner_kind, subject, predicate, value, turn,
         kind="fact", confidence=0.9, span=""):
    return store.add(kind=kind, owner_id=owner, owner_kind=owner_kind,
                     text=f"{subject} {predicate}: {value}", subject=subject,
                     predicate=predicate, value=value,
                     source_span=span or f"t{turn}", established_turn=turn,
                     confidence=confidence)


def build_store() -> RecordStore:
    s = RecordStore()
    # chat scope: world fact, no character owns it
    _rec(s, CHAT, "chat", "the fleet", "flag", "the Alder flag", 5,
         kind="fact", span="world-1")

    # Mira — three successive ship updates (Kestrel -> Wren -> Osprey)
    kestrel = _rec(s, MIRA, "character", "Mira", "ship", "the Kestrel", 8)
    wren = _rec(s, MIRA, "character", "Mira", "ship", "the Wren", 96)
    osprey = _rec(s, MIRA, "character", "Mira", "ship", "the Osprey", 164)
    s.supersede(kestrel.record_id, wren.record_id)
    s.supersede(wren.record_id, osprey.record_id)

    # Mira — non-temporal distractors, same owner
    _rec(s, MIRA, "character", "Mira", "home_port", "Vell", 30)
    _rec(s, MIRA, "character", "Mira", "cargo", "salt", 120)

    # Teodor — a correction of a mistaken statement
    wrong_post = _rec(s, TEODOR, "character", "Teodor", "post", "quartermaster", 20,
                      confidence=0.5)
    right_post = _rec(s, TEODOR, "character", "Teodor", "post", "bosun", 34,
                      kind="state_change")
    s.supersede(wrong_post.record_id, right_post.record_id)

    _rec(s, TEODOR, "character", "Teodor", "ship", "the Gull", 40)

    # Teodor — genuine conflict: two statements, different turns, neither superseded
    dues_a = _rec(s, TEODOR, "character", "Teodor", "harbour_dues", "three coppers", 60)
    dues_b = _rec(s, TEODOR, "character", "Teodor", "harbour_dues", "five coppers", 62)

    return s


def ids(store, owner, predicate=None, include_history=False) -> list[str]:
    out = []
    for r in store.all():
        if r.owner_id != owner:
            continue
        if not include_history and r.status != "active":
            continue
        if predicate and r.predicate != predicate:
            continue
        out.append(r.record_id)
    return out


def build_cases(store) -> list[dict]:
    """Candidate pools are explicit so both conditions see exactly the same pool."""
    def case(cid, criterion, query, cands, expect, include_history=False, gap=False):
        return {"id": cid, "criterion": criterion, "owner": cands.owner, "query": query,
                "candidates": [{"record_id": rid, "text": store.get(rid).text,
                                "value": store.get(rid).value,
                                "status": store.get(rid).status}
                               for rid in cands.ids],
                "expect": expect, "expect_abstain": expect is None,
                "include_history": include_history, "gap": gap}

    class Pool:
        def __init__(self, owner, rid_list):
            self.owner = owner
            self.ids = rid_list

    mira_ships = ids(store, MIRA, "ship", True)
    mira_all = ids(store, MIRA)
    teo_post = ids(store, TEODOR, "post", True)
    teo_dues = ids(store, TEODOR, "harbour_dues")
    teo_ship = ids(store, TEODOR, "ship")
    mira_ship_active = [r for r in mira_all if store.get(r).predicate == "ship"]

    return [
        case("T1-current-no-temporal", "current_without_temporal_language",
             "What ship does Mira sail?", Pool(MIRA, mira_ship_active + mira_all), "the Osprey"),
        case("T2-previous", "explicit_history",
             "Which ship was Mira sailing before the Osprey?", Pool(MIRA, mira_ships), "the Wren",
             include_history=True),
        case("T3-first", "explicit_history_deep",
             "What was Mira's first ship?", Pool(MIRA, mira_ships), "the Kestrel",
             include_history=True),
        case("T4-correction-current", "correction_resolves_to_current",
             "What is Teodor's post?", Pool(TEODOR, teo_post), "bosun"),
        case("T5-correction-history", "correction_history_reachable",
             "What was Teodor's post before the correction?", Pool(TEODOR, teo_post),
             "quartermaster", include_history=True),
        case("T6-conflict", "conflicting_timestamps",
             "What does Teodor owe in harbour dues?", Pool(TEODOR, teo_dues), "five coppers"),
        case("T7-ownership-own", "ownership_positive",
             "What ship does Teodor sail?", Pool(TEODOR, teo_ship), "the Gull"),
        case("T8-ownership-negation", "ownership_negation",
             "What ship does Teodor sail?", Pool(TEODOR, teo_ship + mira_ships), "the Gull"),
        case("T9-distractor", "non_temporal_control",
             "Where is Mira's home port?", Pool(MIRA, mira_all), "Vell"),
        # Scope gap probe: a world-scope fact asked by a character. The pool a
        # character is allowed to see excludes chat-scope records, so the fact is
        # unreachable. That is the parked ownership gap, measured not assumed.
        case("T10-chat-scope-gap", "scope_gap_known",
             "What flag does the fleet sail under?", Pool(MIRA, mira_all),
             None, gap=True),
    ]


# Held-out calibration pairs for this store: new phrasings, records not used by
# the cases above where possible, and never the expected labels.
CALIBRATION = [
    ("Q1", "What goods does Mira carry?", "Mira cargo: salt", True),
    ("Q2", "Which ship is Teodor's?", "Teodor ship: the Gull", True),
    ("Q3", "What was Mira's home port?", "Mira home_port: Vell", True),
    ("Q4", "What is Mira's ship called?", "Mira ship: the Osprey", True),
    ("Q5", "What goods does Mira carry?", "Teodor post: bosun", False),
    ("Q6", "Which ship is Teodor's?", "Mira ship: the Wren", False),
    ("Q7", "What was Mira's home port?", "Teodor harbour_dues: five coppers", False),
    ("Q8", "What is Mira's ship called?", "Teodor post: quartermaster", False),
]

FROZEN_NOTE = ("fresh store, new owner ids, new cases; the frozen seven records and "
               "the R1-R5 case set are neither referenced nor modified")
