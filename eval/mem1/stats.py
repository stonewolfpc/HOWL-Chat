"""The frozen decision rule (Amend. 2 §1–§4, Amend. 3 §2).

Two claims on M1, each evaluated at alpha = .025:

    C1  A1 > A0   net capability improvement
    C2  A1 > D    selection-specific benefit
    full claim = C1 and C2

Ties are excluded from the sign count (`n_eff`) but **included** in the median.
Invalid pairs are excluded outcome-blind, at most two per comparison, and are
never zero-filled. A median below 4/39 is not a capability claim even when the
sign test is significant.
"""
from __future__ import annotations

from dataclasses import dataclass, field
from fractions import Fraction
from math import comb

from . import ALPHA_PER_CLAIM, MAX_EXCLUDED_PAIRS, MIN_N_EFF, MMD
from .errors import InvalidPairScoredError


@dataclass
class Pair:
    conversation_id: str
    a1: float | None
    comparator: float | None
    valid: bool = True
    invalid_reason: str = ""

    @property
    def d(self) -> float | None:
        if not self.valid or self.a1 is None or self.comparator is None:
            return None
        return self.a1 - self.comparator

    def as_dict(self) -> dict:
        return {"conversation_id": self.conversation_id, "a1": self.a1,
                "comparator": self.comparator, "d": self.d, "valid": self.valid,
                "invalid_reason": self.invalid_reason}


def build_pairs(a1_values: dict[str, float], comparator_values: dict[str, float],
                exclusions: dict[str, str] | None = None) -> list[Pair]:
    """One pair per conversation. Mood/median are computed over valid pairs only."""
    exclusions = exclusions or {}
    ids = sorted(set(a1_values) | set(comparator_values) | set(exclusions))
    pairs: list[Pair] = []
    for cid in ids:
        reason = exclusions.get(cid)
        if reason:
            pairs.append(Pair(cid, a1_values.get(cid), comparator_values.get(cid),
                              valid=False, invalid_reason=reason))
        elif cid not in a1_values or cid not in comparator_values:
            pairs.append(Pair(cid, a1_values.get(cid), comparator_values.get(cid),
                              valid=False, invalid_reason="missing_arm_value"))
        else:
            pairs.append(Pair(cid, a1_values[cid], comparator_values[cid], valid=True))
    return pairs


def assert_pairs_clean(pairs: list[Pair]) -> None:
    """Guard against an excluded pair reaching an average.

    An invalid pair that still carries both arm values could silently contribute a
    difference to any caller that ignored the flag, so its presence is a hard error.
    """
    bad = [p.conversation_id for p in pairs
           if not p.valid and p.a1 is not None and p.comparator is not None]
    if bad:
        raise InvalidPairScoredError(f"invalid pairs still carry arm values: {bad}")


def sign_test(k: int, n: int, alpha: float = ALPHA_PER_CLAIM) -> float:
    """Exact two-sided binomial sign test (p = .5), ties already excluded from n."""
    if n <= 0:
        return 1.0
    if not 0 <= k <= n:
        raise ValueError(f"k={k} outside 0..{n}")
    upper = sum(comb(n, i) for i in range(k, n + 1)) / 2**n
    return min(1.0, 2 * upper)


def min_concordant(n: int, alpha: float = ALPHA_PER_CLAIM) -> int | None:
    """Fewest favourable pairs that reach `alpha` two-sided at this n."""
    for k in range(1, n + 1):
        if sign_test(k, n, alpha) < alpha:
            return k
    return None


def median(values: list[float]) -> float | None:
    """Standard median; even counts average the two central order statistics.

    Ties are included: they carry no direction, but they are part of the
    magnitude distribution (Amend. 3 §2).
    """
    if not values:
        return None
    ordered = sorted(values)
    n = len(ordered)
    mid = n // 2
    if n % 2:
        return ordered[mid]
    return (ordered[mid - 1] + ordered[mid]) / 2


def evaluate_claim(name: str, pairs: list[Pair], alpha: float = ALPHA_PER_CLAIM,
                   mmd: Fraction = MMD) -> dict:
    """One claim's verdict. Never fuses with the other claim."""
    valid = [p for p in pairs if p.valid]
    excluded = [p for p in pairs if not p.valid]
    diffs = [p.d for p in valid if p.d is not None]
    nonzero = [d for d in diffs if d != 0]
    ties = len(diffs) - len(nonzero)
    n_eff = len(nonzero)
    favourable = sum(1 for d in nonzero if d > 0)

    out = {"claim": name, "n_valid": len(diffs), "n_pairs": len(pairs),
           "excluded": [p.as_dict() for p in excluded], "excluded_count": len(excluded),
           "ties": ties, "n_eff": n_eff, "favourable": favourable,
           "median": median(diffs) if diffs else None,
           "mmd": float(mmd), "alpha": alpha, "p": None, "threshold": None,
           "significant": False, "mmd_met": False, "disposition": "indeterminate",
           "diffs": [p.as_dict() for p in valid]}

    if len(excluded) > MAX_EXCLUDED_PAIRS:
        out["reason"] = f"{len(excluded)} excluded pairs exceed the cap of {MAX_EXCLUDED_PAIRS}"
        return out
    if n_eff < MIN_N_EFF:
        out["reason"] = f"n_eff={n_eff} below the interpretability floor of {MIN_N_EFF}"
        return out

    threshold = min_concordant(n_eff, alpha)
    p = sign_test(favourable, n_eff, alpha)
    med = out["median"]
    significant = p < alpha and favourable * 2 > n_eff
    mmd_met = med is not None and med >= float(mmd)

    out.update({"p": p, "threshold": threshold, "significant": significant,
                "mmd_met": mmd_met, "reason": ""})

    if significant and mmd_met:
        out["disposition"] = "supported"
    elif significant:
        out["disposition"] = "significant_but_negligible"
    elif favourable * 2 < n_eff and sign_test(n_eff - favourable, n_eff, alpha) < alpha:
        out["disposition"] = "reversed"
    else:
        out["disposition"] = "not_supported"
    return out


def decide(c1: dict, c2: dict, vetoes: list[str] | None = None) -> dict:
    """Combine the two claims exactly as Amend. 2 §1 specifies."""
    vetoes = vetoes or []
    if vetoes:
        return {"disposition": "run_failed", "vetoes": list(vetoes),
                "statement": "a veto fired; no aggregate or claim is offered",
                "c1": c1, "c2": c2}

    d1, d2 = c1["disposition"], c2["disposition"]
    if "indeterminate" in (d1, d2):
        disposition, statement = "indeterminate", (
            "at least one comparison is indeterminate (too few evaluable pairs or too many "
            "exclusions); no capability claim is made"
        )
    elif d1 == "reversed" or d2 == "reversed":
        disposition, statement = "regression", (
            "a comparison is significantly reversed: reported as a regression, not a null result"
        )
    elif d1 == "supported" and d2 == "supported":
        disposition, statement = "full_capability_claim", (
            "memory improved capability (C1) and relevance drove it (C2)"
        )
    elif d1 == "supported":
        disposition, statement = "net_only", (
            "capability improved, but the gain is not shown to be selection-specific: "
            "tokens or presentation could explain it (C2 not supported)"
        )
    elif d2 == "supported":
        disposition, statement = "selection_only", (
            "relevance beats decoys, but not enough to beat no memory (C1 not supported)"
        )
    else:
        disposition, statement = "no_claim", (
            "neither claim is supported; significant-but-negligible results are not claims"
        )
    return {"disposition": disposition, "statement": statement,
            "vetoes": [], "c1": c1, "c2": c2}


__all__ = ["Pair", "build_pairs", "assert_pairs_clean", "sign_test", "min_concordant",
           "median", "evaluate_claim", "decide"]
