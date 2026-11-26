"""Experiment E-PROSE-RENDER — natural-language rendering of the same records.

Registered in rerank-spec-v1.2.json as a separate condition. It changes ONLY how
a record is presented to the candidate. Records, queries, case set, model,
runtime and thresholds are held fixed.

Schema condition (C-SCHEMA), still the original:
    "Sable trained_by: Corvin Ashgrave of Ossary, and quiet was not part of it"
Prose condition (E-PROSE-RENDER):
    "Sable was trained by Corvin Ashgrave of Ossary, and quiet was not part of it."
"""
from __future__ import annotations

RENDER = {
    "rec-0001": "Sable was trained by Corvin Ashgrave of Ossary, and quiet was not part of it.",
    "rec-0002": "Sable used to keep the ledger of unpaid debts in the blue cabinet.",
    "rec-0003": "Sable likes tea more than she likes surprises.",
    "rec-0004": "Alaric's vault code is seven-four-nine-one.",
    "rec-0005": "Alaric served four years in the Ninth Company.",
    "rec-0006": "Sable never trusted a map that was drawn by someone who had not walked the ground.",
    "rec-0007": "Sable keeps the ledger of debts in the green cabinet now.",
}

# Held-out calibration pairs, same records, prose form. Different phrasings and
# pairings from the nine cases, exactly as in the schema condition.
CALIBRATION = [
    ("P1", "How long did you serve in the company?",
     "Alaric served four years in the Ninth Company.", True),
    ("P2", "What did you learn about trusting maps?",
     "Sable never trusted a map that was drawn by someone who had not walked the ground.", True),
    ("P3", "Where do the debts sit now?",
     "Sable keeps the ledger of debts in the green cabinet now.", True),
    ("P4", "What is the code for the vault?",
     "Alaric's vault code is seven-four-nine-one.", True),
    ("P5", "How long did you serve in the company?",
     "Sable keeps the ledger of debts in the green cabinet now.", False),
    ("P6", "Where do the debts sit now?",
     "Sable likes tea more than she likes surprises.", False),
    ("P7", "Do you take sugar with it?",
     "Sable never trusted a map that was drawn by someone who had not walked the ground.", False),
    ("P8", "What is the code for the vault?",
     "Sable was trained by Corvin Ashgrave of Ossary, and quiet was not part of it.", False),
]


def render_cases(cases: list[dict]) -> list[dict]:
    """Swap candidate text for prose, keeping ids, order and expectations."""
    out = []
    for case in cases:
        out.append({**case,
                    "candidates": [{**c, "text": RENDER[c["record_id"]]}
                                   for c in case["candidates"]]})
    return out
