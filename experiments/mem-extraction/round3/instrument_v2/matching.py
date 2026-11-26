"""Stem-tolerant content matching for instrument v2.

Inherited semantics from the round-3 scorer (`../score_corpus.py::_alt_matches`),
kept deliberately so that v2 content recall is comparable with the historical
measurement. What is NOT inherited is the historical coupling of matching to the
decision token and to contract validity.

  v1: an item counted only if a *parseable* row carried the expected decision.
  v2: content is matched in any row that carries text, and the decision and
      contract status are recorded separately.

That separation is required by audit finding F-03: a model that states the right
fact in a broken format must be distinguishable from a model that misunderstands
the fact. It is not credited as production-usable; the contract score says so.
"""

from __future__ import annotations

import re

_ALT_PHRASE_CACHE: dict[str, list[str]] = {}


def norm(s: str) -> str:
    return re.sub(r"\s+", " ", (s or "").lower())


def words(s: str) -> list[str]:
    return re.findall(r"[a-z0-9']+", (s or "").lower())


def alt_matches(alt: str, toks: list[str]) -> bool:
    """Prefix match tolerant of plurals/participles.

    'narrating' matches 'narrate'; 'avoids' matches 'avoid'. Negations must be
    present as an explicit negation token, so 'no purple prose' cannot be
    satisfied by the word 'no' appearing anywhere unrelated.
    """
    parts = words(alt)
    if not parts:
        return False
    for p in parts:
        if p in ("no", "not", "never"):
            if not any(w in ("no", "not", "never", "without") for w in toks):
                return False
            continue
        stem = p[: max(4, len(p) - 2)]
        if not any(w.startswith(stem) for w in toks):
            return False
    return True


def group_matches(group: list[str], text: str) -> bool:
    toks = words(text)
    return any(alt_matches(alt, toks) for alt in group)


def item_in_text(item: dict, text: str) -> bool:
    """True when every keyword group of the required item is satisfied."""
    return all(group_matches(g, text) for g in item.get("keys", []))


def match_items(records: list, required: list[dict]) -> dict[str, dict]:
    """Map required-item id -> best match across records.

    `records` are outcomes.Record objects. Contract-invalid rows are searched
    too (record.text is still available); the result records which row matched
    and whether that row was contract-valid and carried the expected decision.

    A required item is matched at most once, to the first row that contains it,
    because duplicate rows are a duplication defect, not extra recall.
    """
    found: dict[str, dict] = {}
    for item in required:
        for idx, rec in enumerate(records):
            if item_in_text(item, rec.text):
                found[item["id"]] = {
                    "record_index": idx,
                    "contract_valid": bool(getattr(rec, "contract_valid", False)),
                    "decision": getattr(rec, "decision", None),
                    "decision_correct": getattr(rec, "decision", None) == item.get("decision"),
                    "raw": getattr(rec, "raw", ""),
                }
                break
    return found
