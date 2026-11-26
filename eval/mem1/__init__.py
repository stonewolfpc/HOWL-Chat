"""MEMORY-09 G1 — experimental instrument for the frozen preregistration.

This package *executes* the frozen design; it does not define it. Every constant
here is transcribed from a frozen document and cites its source:

  * N = 16 paired conversations, 3 matched reps      prereg §2
  * 13 probes per conversation, quota 4/2/2/2/1/1/1  prereg §5
  * alpha = .025 per claim (family-wise .05)         Amend. 2 §3.1
  * MMD = 4/39 (smallest attainable median >= 0.10)  Amend. 3 §1.2
  * M1 / M2 definitions                              Amend. 1 §2-§3
  * C1 = A1>A0 on M1, C2 = A1>D on M1                Amend. 2 §1

No model call is made by anything in this package during G1. The driver accepts
a `Responder`; the only implementation shipped here is a scripted one, and it
refuses to run if a responder claims to be networked without a G4 grant.

    PYTHONPATH=. python eval/mem1/selftest.py
"""
from __future__ import annotations

from fractions import Fraction

MEM1_VERSION = "0.1.0-g1"

# --- frozen design constants (transcribed, never tuned here) -----------------

N_CONVERSATIONS = 16
REPS_PER_CELL = 3
PROBES_PER_CONVERSATION = 13

# prereg §5: category quotas, identical in every arm.
PROBE_QUOTAS = {
    "character_specific": 4,
    "ambient": 2,
    "restricted": 2,
    "experiential": 2,
    "correction": 1,
    "retraction": 1,
    "unanswerable": 1,
}

# Amend. 2 §3.1: family-wise .05 across the two claims.
ALPHA_PER_CLAIM = 0.025
# Amend. 2 §2: below this many sign-bearing pairs a comparison is indeterminate.
MIN_N_EFF = 12
# Amend. 2 §3.3 / Amend. 3 §1.2: median paired difference must reach 4/39.
MMD = Fraction(4, 39)
# Amend. 2 §4: at most 2 of 16 pairs may be excluded per comparison.
MAX_EXCLUDED_PAIRS = 2

# Plan §2.4/§2.5: one prompt capacity shared by all arms; memory is paid for
# out of it and displaces history. The placebo tolerance is declared here, in
# advance, and asserted every turn.
DEFAULT_CAPACITY = 900
DEFAULT_RESERVE = 160
PLACEBO_TOKEN_TOLERANCE = 0

# Plan §2.1 arms.
ARMS = ("A0", "A1", "A2", "C", "D")
# Plan §4.2 / Amend. 1: A0 has no records, so M2 is undefined for it.
M2_UNDEFINED_ARMS = ("A0",)
# Arms whose records must cite genuine spans (plan §2.3 anti-cheat rule).
FORMATION_ARMS = ("A1", "A2")

# M1 categories (Amend. 1 §2) — exactly one per response.
CORRECT = "CORRECT"
INCORRECT = "INCORRECT"
ABSTAIN_APPROPRIATE = "ABSTAIN_APPROPRIATE"
ABSTAIN_INAPPROPRIATE = "ABSTAIN_INAPPROPRIATE"
MALFORMED = "MALFORMED"
M1_CATEGORIES = (CORRECT, INCORRECT, ABSTAIN_APPROPRIATE, ABSTAIN_INAPPROPRIATE, MALFORMED)

__all__ = [
    "MEM1_VERSION", "N_CONVERSATIONS", "REPS_PER_CELL", "PROBES_PER_CONVERSATION",
    "PROBE_QUOTAS", "ALPHA_PER_CLAIM", "MIN_N_EFF", "MMD", "MAX_EXCLUDED_PAIRS",
    "DEFAULT_CAPACITY", "DEFAULT_RESERVE", "PLACEBO_TOKEN_TOLERANCE", "ARMS",
    "M2_UNDEFINED_ARMS", "FORMATION_ARMS", "M1_CATEGORIES",
    "CORRECT", "INCORRECT", "ABSTAIN_APPROPRIATE", "ABSTAIN_INAPPROPRIATE", "MALFORMED",
]
