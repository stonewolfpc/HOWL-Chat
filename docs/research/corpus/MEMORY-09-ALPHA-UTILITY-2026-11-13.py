"""MEMORY-09 alpha utility (nominal Krippendorff's alpha).

Status: VALIDATION CANDIDATE. Authorised by MEMORY-09-WO-ARCH-05 Stage One.
Authority: owner work order of 2026-11-13. This file is a certified *calculator*,
not a study. It must never be pointed at live corpus ratings under the current
authorisation.

Provenance note: the implementation used earlier in session was executed from an
unnamed heredoc and was never written to disk, so it cannot be recovered
byte-for-byte. This file is a RECONSTRUCTION and is identified as such in
MEMORY-09-ALPHA-UTILITY-VALIDATION-2026-11-13.md. It is certified only by passing
the full gate suite from this filed version.

Design constraints (from the work order):
  * nominal categories only, handled explicitly
  * missing ratings and dropped units handled exactly as the published method requires
  * coincidence matrix, marginals, Do and De all inspectable
  * deterministic for fixed input
  * no silent coercion of invalid values
  * no network access, no corpus access, no hidden defaults
  * does NOT invent a policy for not-assessable values
"""

from __future__ import annotations

from fractions import Fraction

__version__ = "1.0.1"

VALID_METHODS = ("coincidence", "pairs")
VALID_NA_POLICIES = ("missing", "category")


class AlphaError(ValueError):
    """Raised for invalid input. Never swallowed: no silent coercion."""


def _clean_units(units, categories, not_assessable, na_policy):
    """Validate and normalise units. Returns (units, categories)."""
    if na_policy not in VALID_NA_POLICIES:
        raise AlphaError(f"na_policy must be one of {VALID_NA_POLICIES}, got {na_policy!r}")

    cats = list(categories)
    if len(cats) < 2:
        raise AlphaError("at least two categories are required for alpha to be defined")

    normalised = []
    for idx, unit in enumerate(units):
        values = []
        for value in unit:
            if value is None:
                if na_policy == "category":
                    raise AlphaError(
                        "None cannot serve as a not-assessable category; pass an explicit "
                        "not_assessable label or use na_policy='missing'"
                    )
                continue
            if not_assessable is not None and value == not_assessable:
                if na_policy == "missing":
                    continue
                if not_assessable not in cats:
                    cats.append(not_assessable)
            elif not_assessable is None and value not in cats:
                raise AlphaError(
                    f"unit {idx} contains {value!r}, which is not a declared category; "
                    "declare it or set a not_assessable policy (no silent coercion)"
                )
            values.append(value)
        normalised.append(values)
    return normalised, cats


def coincidence_matrix(units, categories):
    """Krippendorff coincidence matrix; units with fewer than two values drop out."""
    index = {c: 0 for c in categories}
    for c in categories:
        index[c] = {}
    o = {c: {k: Fraction(0) for k in categories} for c in categories}
    kept, dropped = [], []
    for pos, values in enumerate(units):
        if len(values) < 2:
            dropped.append(pos)
            continue
        kept.append(values)
        m = len(values)
        for c in categories:
            for k in categories:
                numerator = values.count(c) * values.count(k) - (values.count(c) if c == k else 0)
                o[c][k] += Fraction(numerator, m - 1)
    return o, kept, dropped


def alpha_from_coincidence(o, categories):
    """alpha = 1 - Do/De, computed from a coincidence matrix."""
    n = sum(sum(o[c].values()) for c in categories)
    marginals = {c: sum(o[c].values()) for c in categories}
    if n < 2:
        raise AlphaError("fewer than two pairable values; alpha is undefined")
    do = sum(o[c][k] for c in categories for k in categories if c != k) / n
    de_num = sum(marginals[c] * marginals[k] for c in categories for k in categories if c != k)
    de = Fraction(de_num, n * (n - 1))
    if de == 0:
        return {"n": n, "marginals": marginals, "do": do, "de": de, "alpha": None,
                "reason": "expected disagreement is zero: every value has one category"}
    return {"n": n, "marginals": marginals, "do": do, "de": de, "alpha": 1 - do / de}


def alpha_pairs(units, categories):
    """Second path: enumerate unordered coder pairs within each unit directly.

    Defect fixed at v1.0.1 (found by Gate D): the first version counted
    DISAGREEING UNITS over UNITS, which is correct only when every unit has
    exactly two coders. With four coders a unit contains six unordered pairs,
    so the weight must be 2/(m-1) per disagreeing pair, summed over
    n = total padded values. The uncorrected version returned 1087/1672 on the
    canonical fixture against a published 113/152.
    """
    kept = [v for v in units if len(v) >= 2]
    if not kept:
        raise AlphaError("no unit has two or more values; alpha is undefined")
    values = [x for v in kept for x in v]
    n = len(values)
    marginals = {c: values.count(c) for c in categories}
    weighted = Fraction(0)
    for v in kept:
        m = len(v)
        disagreeing_pairs = sum(
            1 for i in range(m) for j in range(i + 1, m) if v[i] != v[j]
        )
        weighted += Fraction(2 * disagreeing_pairs, m - 1)
    do = weighted / n
    de_num = sum(marginals[c] * marginals[k] for c in categories for k in categories if c != k)
    de = Fraction(de_num, n * (n - 1))
    if de == 0:
        return {"n": n, "marginals": marginals, "do": do, "de": de, "alpha": None,
                "reason": "expected disagreement is zero"}
    return {"n": n, "marginals": marginals, "do": do, "de": de, "alpha": 1 - do / de}


def alpha(units, categories, *, na_policy="missing", not_assessable=None, method="coincidence"):
    """Nominal Krippendorff's alpha over units x coders.

    units: list of lists; None marks a missing rating.
    Returns a dict with n, marginals, do, de, alpha (Fraction or None), plus the
    coincidence matrix and dropped-unit indices when method='coincidence'.
    """
    if method not in VALID_METHODS:
        raise AlphaError(f"method must be one of {VALID_METHODS}, got {method!r}")
    cleaned, cats = _clean_units(units, categories, not_assessable, na_policy)
    if method == "pairs":
        return alpha_pairs(cleaned, cats)
    o, kept, dropped = coincidence_matrix(cleaned, cats)
    result = alpha_from_coincidence(o, cats)
    result["coincidence"] = o
    result["dropped_units"] = dropped
    result["pairable_units"] = len(kept)
    return result


def as_float(value):
    return None if value is None else float(value)
