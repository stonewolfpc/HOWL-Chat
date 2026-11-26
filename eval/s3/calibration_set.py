"""Held-out calibration pairs for abstention thresholds (v1.2 protocol).

These are NOT the nine R1–R5 cases. They use different phrasings and different
query–record pairings, so a threshold fitted here is not fitted on the
evaluation set.

Limitation, stated rather than hidden: both sets draw on the same seven records.
Independence here is independence of phrasing and pairing, not of corpus.
"""
from __future__ import annotations

# (pair_id, query, document_text, relevant?)
PAIRS = [
    ("C1", "How long did you serve in the company?",
     "Alaric served: four years in the Ninth Company", True),
    ("C2", "What did you learn about trusting maps?",
     "Sable never: trusted a map that was drawn by someone who had not walked the ground", True),
    ("C3", "Where do the debts sit now?",
     "Sable keeps_ledger: of debts in the green cabinet now", True),
    ("C4", "What is the code for the vault?",
     "Alaric vault_code: seven-four-nine-one", True),

    ("C5", "How long did you serve in the company?",
     "Sable keeps_ledger: of debts in the green cabinet now", False),
    ("C6", "Where do the debts sit now?",
     "Sable preference: tea more than I like surprises", False),
    ("C7", "Do you take sugar with it?",
     "Sable never: trusted a map that was drawn by someone who had not walked the ground", False),
    ("C8", "What is the code for the vault?",
     "Sable trained_by: Corvin Ashgrave of Ossary, and quiet was not part of it", False),
]

RULE = ("threshold = midpoint between the highest irrelevant score and the lowest "
        "relevant score on these held-out pairs; if the bands overlap, abstention "
        "is reported UNCALIBRATABLE rather than fitted anyway")
