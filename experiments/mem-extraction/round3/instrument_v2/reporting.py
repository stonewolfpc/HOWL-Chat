"""Reporting language rules — instrument v2 (findings F-05, F-07, F-09, F-10, F-11).

The defect these rules address is not arithmetic; it is wording. A number can be
correct and its sentence false. This module encodes what a claim is allowed to
say given its sample size and interval, and provides a scanner so an existing
document can be checked rather than trusted.

`scan()` returns rule hits with the exact matched text, so an erratum can quote
the source instead of paraphrasing it.
"""

from __future__ import annotations

import re

#: Sentences that are never permitted, whatever the number says.
FORBIDDEN_PATTERNS = [
    (r"\bproves?\b", "PROOF", "measurement never proves; say 'is consistent with'"),
    (r"\bproven\b", "PROOF", "measurement never proves; say 'observed under these conditions'"),
    (r"\bconfirms? that\b", "PROOF", "say 'is consistent with' or 'supports'"),
    (r"\bdoes not beat\b", "NEGATIVE_CLAIM", "n<5 cannot establish no difference; say 'no measured gain'"),
    (r"\bno difference\b", "NEGATIVE_CLAIM", "say 'no measured difference at this n'"),
    (r"\bequivalent\b", "EQUIVALENCE", "overlapping intervals do not establish equivalence"),
    (r"\bthe same as\b", "EQUIVALENCE", "say 'not distinguishable at this sample size'"),
    (r"\bwinner\b", "SELECTION", "no model is selected by this instrument"),
    (r"\branking\b", "RANKING", "use 'ordering by mean score' and state the n"),
    (r"\bscales? linearly\b", "GENERALISATION", "one card, one config; do not generalise"),
    (r"\bcontext flush(ing)? (works|is proven)\b", "ISOLATION",
     "only serial per-request isolation was observed"),
]

QUALIFIERS = {
    "single_draw": "single-draw observation (n=1); not a ranking",
    "underpowered": "underpowered (n<5); direction only",
    "adequate": "n>=5; still corpus-specific",
}


def sample_qualifier(draws: float) -> str:
    if draws < 2:
        return QUALIFIERS["single_draw"]
    if draws < 5:
        return QUALIFIERS["underpowered"]
    return QUALIFIERS["adequate"]


def compare(a_mean: float | None, a_ci: tuple[float, float] | None, a_n: float,
            b_mean: float | None, b_ci: tuple[float, float] | None, b_n: float) -> dict:
    """Classify what may be said about two conditions."""
    if a_mean is None or b_mean is None:
        return {"claim": "UNESTABLISHED", "basis": "missing measurement",
                "allowed": "not measured; say so"}
    delta = round(a_mean - b_mean, 2)
    thin = a_n < 5 or b_n < 5
    overlap = None
    if a_ci and b_ci:
        overlap = not (a_ci[1] < b_ci[0] or b_ci[1] < a_ci[0])

    if thin:
        claim = "UNESTABLISHED"
        allowed = (f"delta {delta:+} points at n={a_n:g}/{b_n:g}: no measured "
                   "difference at this sample size; this is not evidence that a "
                   "difference does not exist")
    elif overlap:
        claim = "NOT_ESTABLISHED"
        allowed = (f"delta {delta:+} points; intervals overlap, so the difference "
                   "is not established and equivalence is not established either")
    else:
        claim = "SUGGESTIVE"
        allowed = (f"delta {delta:+} points with non-overlapping intervals on this "
                   "corpus; still not a production selection")
    return {"claim": claim, "delta": delta, "thin": thin,
            "intervals_overlap": overlap, "allowed": allowed,
            "forbidden": "proves / no difference / equivalent / winner"}


def effect(label: str, n: int, deltas: list[float]) -> dict:
    """A sampler or ablation delta is a directional observation, not an effect size."""
    return {
        "label": label,
        "n": n,
        "point_estimate": round(sum(deltas) / len(deltas), 2) if deltas else None,
        "observed_range": (round(min(deltas), 2), round(max(deltas), 2)) if deltas else None,
        "statement": (f"{label}: directional observation at n={n}; point estimate is "
                      "not a causal effect and the configuration differs from any "
                      "baseline it is compared with"),
    }


ISOLATION_ALLOWED = (
    "Under serial, one-request-per-sequence execution at -np 1 on this build, a "
    "benign job produced byte-identical output before and after an adversarial "
    "job, and a fresh server produced the same output. The erase endpoint "
    "returned HTTP 501, so no reset API was exercised."
)
ISOLATION_PROHIBITED = [
    "context flushing works",
    "the context is reset",
    "isolation is guaranteed",
    "no cross-contamination is possible",
    "holds for concurrent requests",
    "holds for other model families",
    "holds for MEMORY-09 records",
]

TIMING_ALLOWED = (
    "Wall time excludes model load and includes prefill and decode. CPU and GPU "
    "runs do not share launch flags (GPU adds full offload, Flash Attention and "
    "unified KV). Peak VRAM is nvidia-smi memory.used and includes other desktop "
    "consumers on this machine."
)
TIMING_PROHIBITED = [
    "the model uses X GB",
    "scales to other cards",
    "partial offload is unnecessary",
]

PREREGISTRATION = {
    "status": "UNESTABLISHED",
    "why": ("the corpus and rubric files are untracked; no hash was published "
            "before generation, so the order of writing cannot be proven from "
            "the repository"),
    "what_would_establish_it": ("a signed or committed hash of corpus.json and "
                                "RUBRIC.md dated before the first run_code, plus "
                                "a commit for the scoring code"),
}


def scan(text: str) -> list[dict]:
    """Return rule hits with the exact matched phrase and its sentence."""
    hits = []
    for pattern, code, advice in FORBIDDEN_PATTERNS:
        for m in re.finditer(pattern, text or "", re.I):
            start = max(0, text.rfind(".", 0, m.start()) + 1)
            end = text.find(".", m.end())
            sentence = re.sub(r"\s+", " ", text[start:end + 1 if end != -1 else len(text)]).strip()
            hits.append({"code": code, "match": m.group(0),
                         "sentence": sentence[:400], "advice": advice})
    return hits
