"""MEMORY-09 alpha utility fixtures and gate runner.

Status: VALIDATION SUPPORT. Authorised by MEMORY-09-WO-ARCH-05 Stage One.
Every gate runs against the FILED utility module (no in-session copy).

Published fixtures and their sources:
  Gate C1  Krippendorff, "Computing Krippendorff's Alpha-Reliability" (UPenn).
           4 observers x 12 units, 5 nominal categories, 7 missing. Published
           alpha = .743, n = 40 pairable values, marginals 9/13/10/5/3.
  Gate C2  Freelon, "Worked Examples for Nominal Intercoder Reliability" (ReCal).
           Published coincidence matrix, n = 20, published alpha = .85.
"""

from __future__ import annotations

import random
import statistics
from fractions import Fraction

import importlib.util
import pathlib

_UTIL = pathlib.Path(__file__).with_name("MEMORY-09-ALPHA-UTILITY-2026-11-13.py")
_spec = importlib.util.spec_from_file_location("m09_alpha_utility", _UTIL)
alpha_util = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(alpha_util)

# ---------------------------------------------------------------- fixtures

# rows = units, columns = observers A-D; None = missing rating
KRIPPENDORFF_UNITS = [
    [1, 1, None, 1],
    [2, 2, 3, 2],
    [3, 3, 3, 3],
    [3, 3, 3, 3],
    [2, 2, 2, 2],
    [1, 2, 3, 4],
    [4, 4, 4, 4],
    [1, 1, 2, 1],
    [2, 2, 2, 2],
    [None, 5, 5, 5],
    [None, None, 1, 1],
    [None, 3, None, None],
]
KRIPPENDORFF_CATS = [1, 2, 3, 4, 5]
KRIPPENDORFF_EXPECTED = {
    "n": 40,
    "marginals": {1: 9, 2: 13, 3: 10, 4: 5, 5: 3},
    "alpha": Fraction(113, 152),  # 0.74342..., published .743
}

RECAL_MATRIX = [[8, 1, 0], [1, 6, 0], [0, 0, 4]]
RECAL_CATS = [0, 1, 2]
RECAL_EXPECTED = {"n": 20, "alpha": Fraction(108, 127)}  # 0.85039..., published .85


def perfect_fixture(units=45, seed=3):
    rng = random.Random(seed)
    cats = ["D", "ND", "NA", "IE"]
    return [[(v := rng.choice(cats)), v] for _ in range(units)], cats


def random_fixture(n_units, missing_rate, seed):
    rng = random.Random(seed)
    cats = ["D", "ND", "NA", "IE"]
    units = []
    for _ in range(n_units):
        a = rng.choice(cats)
        b = rng.choice(cats)
        if rng.random() < missing_rate:
            b = None
        units.append([a, b])
    return units, cats


# ---------------------------------------------------------------- gates

def gate_a():
    units, cats = perfect_fixture()
    res = alpha_util.alpha(units, cats)
    ok = res["alpha"] == 1
    return {"gate": "A perfect agreement", "alpha": str(res["alpha"]), "pass": ok}


def gate_b(reps=400, n_units=45):
    rows = []
    for missing in (0.0, 0.10, 0.25):
        vals = []
        for rep in range(reps):
            units, cats = random_fixture(n_units, missing, seed=1000 + rep)
            res = alpha_util.alpha(units, cats)
            if res["alpha"] is not None:
                vals.append(float(res["alpha"]))
        vals.sort()
        rows.append({
            "missing": missing,
            "reps": len(vals),
            "mean": round(statistics.fmean(vals), 4),
            "p05": round(vals[int(0.05 * len(vals))], 3),
            "p95": round(vals[int(0.95 * len(vals))], 3),
        })
    ok = all(abs(r["mean"]) < 0.05 for r in rows)
    return {"gate": "B random simulations", "conditions": rows, "pass": ok}


def gate_c():
    res = alpha_util.alpha(KRIPPENDORFF_UNITS, KRIPPENDORFF_CATS)
    c1_ok = (res["n"] == KRIPPENDORFF_EXPECTED["n"]
             and {k: int(v) for k, v in res["marginals"].items()} == KRIPPENDORFF_EXPECTED["marginals"]
             and res["alpha"] == KRIPPENDORFF_EXPECTED["alpha"])

    o = {c: {k: Fraction(RECAL_MATRIX[i][j]) for j, k in enumerate(RECAL_CATS)}
         for i, c in enumerate(RECAL_CATS)}
    recal = alpha_util.alpha_from_coincidence(o, RECAL_CATS)
    c2_ok = recal["n"] == RECAL_EXPECTED["n"] and recal["alpha"] == RECAL_EXPECTED["alpha"]

    return {
        "gate": "C published worked examples",
        "c1": {"alpha": str(res["alpha"]), "float": round(float(res["alpha"]), 5),
               "n": res["n"], "marginals": {k: int(v) for k, v in res["marginals"].items()},
               "expected": str(KRIPPENDORFF_EXPECTED["alpha"]), "pass": c1_ok},
        "c2": {"alpha": str(recal["alpha"]), "float": round(float(recal["alpha"]), 5),
               "n": recal["n"], "expected": str(RECAL_EXPECTED["alpha"]), "pass": c2_ok},
        "pass": c1_ok and c2_ok,
    }


def gate_d():
    checks = []
    for name, units, cats in (
        ("krippendorff", KRIPPENDORFF_UNITS, KRIPPENDORFF_CATS),
        ("perfect", *perfect_fixture()),
    ):
        a = alpha_util.alpha(units, cats, method="coincidence")
        b = alpha_util.alpha(units, cats, method="pairs")
        checks.append({"fixture": name, "coincidence": str(a["alpha"]), "pairs": str(b["alpha"]),
                       "pass": a["alpha"] == b["alpha"] and a["do"] == b["do"]})
    return {"gate": "D independent paths", "checks": checks,
            "note": "ReCal is matrix-only input; the pairs path does not apply",
            "pass": all(c["pass"] for c in checks)}


def gate_e():
    res = alpha_util.alpha(KRIPPENDORFF_UNITS, KRIPPENDORFF_CATS)
    return {
        "gate": "E inspectable intermediates",
        "n": res["n"],
        "dropped_units": res["dropped_units"],
        "marginals": {k: int(v) for k, v in res["marginals"].items()},
        "do": str(res["do"]),
        "de": str(res["de"]),
        "alpha": str(res["alpha"]),
        "coincidence_diagonal": {k: str(res["coincidence"][k][k]) for k in KRIPPENDORFF_CATS},
        "pass": res["dropped_units"] == [11],
    }


def gate_f():
    """Negative controls: the utility must refuse invalid input.

    Fixture defects fixed at v1.0.1. Two expectations were wrong, not the utility:
      * a single unit holding three values is a VALID alpha input (n=3), so it
        must not be rejected; the insufficient-data case needs units that all
        drop out (each with fewer than two values);
      * zero expected disagreement is a defined degeneracy, not an exception:
        the utility returns alpha=None with a stated reason, which is the
        behaviour the order asks for (flag, never coerce).
    """
    cases = []
    try:
        alpha_util.alpha([[1], [2]], [1, 2])
        cases.append(("insufficient data rejected", False))
    except alpha_util.AlphaError:
        cases.append(("insufficient data rejected", True))
    try:
        alpha_util.alpha([[1, 2], [1, 3]], [1, 2])
        cases.append(("undeclared category rejected", False))
    except alpha_util.AlphaError:
        cases.append(("undeclared category rejected", True))
    try:
        alpha_util.alpha([[1, None], [None, 1]], [1, 2], na_policy="category")
        cases.append(("undefined NA policy rejected", False))
    except alpha_util.AlphaError:
        cases.append(("undefined NA policy rejected", True))
    same = alpha_util.alpha([[1, 1], [1, 1]], [1, 2])
    cases.append(("zero expected disagreement flagged, not invented",
                  same["alpha"] is None and "reason" in same))
    single = alpha_util.alpha([[1, 1, 1]], [1, 2])
    cases.append(("valid three-value unit accepted", single["n"] == 3))
    return {"gate": "F negative controls", "cases": cases, "pass": all(v for _, v in cases)}


def run():
    results = [gate_a(), gate_b(), gate_c(), gate_d(), gate_e(), gate_f()]
    return results


if __name__ == "__main__":
    import json
    for r in run():
        print(json.dumps(r, default=str, sort_keys=False))
    print("ALL GATES PASS:", all(r["pass"] for r in run()))
