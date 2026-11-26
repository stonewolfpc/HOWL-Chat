"""Fixed-sample stopping, decided without looking at outcomes (prereg §7).

The run ends when the preregistered cells are complete. There is no interim
analysis. The only permitted halts are **instrument** failures, and this module
is deliberately blind to `m1_category` / `m2_status`: it reads structure only,
so a halt cannot be triggered by an inconvenient result.

A halt is not a result. The instrument is repaired and the whole fixed sample is
restarted; no partial data is reused, and no extension is allowed after outcomes
are visible.
"""
from __future__ import annotations

from dataclasses import dataclass

from . import N_CONVERSATIONS, REPS_PER_CELL
from .errors import DesignViolationError

MISSING_PROBE_LIMIT = 0.10        # prereg §7 condition 1


@dataclass
class HealthReport:
    expected_cells: int
    observed_cells: int
    expected_probes: int
    observed_probes: int
    missing_fraction: float
    halts: list[str]

    @property
    def halted(self) -> bool:
        return bool(self.halts)

    def as_dict(self) -> dict:
        return {"expected_cells": self.expected_cells, "observed_cells": self.observed_cells,
                "expected_probes": self.expected_probes, "observed_probes": self.observed_probes,
                "missing_fraction": self.missing_fraction, "halts": list(self.halts)}


def evaluate_instrument_health(receipts: list, arms: tuple[str, ...] = ("A0", "A1", "D"),
                               probes_per_conversation: int = 13) -> HealthReport:
    """Outcome-blind: uses presence, arm, turn and token fields, never correctness."""
    expected_cells = N_CONVERSATIONS * REPS_PER_CELL * len(arms)
    cells = {(r.conversation_id, r.rep_index, r.arm) for r in receipts}
    expected_probes = expected_cells * probes_per_conversation
    observed = len(receipts)

    halts: list[str] = []
    missing_fraction = 1.0 - (observed / expected_probes) if expected_probes else 1.0
    if missing_fraction > MISSING_PROBE_LIMIT:
        halts.append(f"missing_probe_receipts:{missing_fraction:.3f}")
    if len(cells) < expected_cells:
        halts.append(f"incomplete_cells:{len(cells)}/{expected_cells}")
    for r in receipts:
        tokens = getattr(r, "tokens", None) or {}
        if not tokens or tokens.get("total") is None:
            halts.append(f"token_accounting_missing:{r.probe_id}:{r.arm}")
        if getattr(r, "selected_ids", None) and not getattr(r, "block_sha256", ""):
            halts.append(f"citation_without_block:{r.probe_id}:{r.arm}")
    return HealthReport(expected_cells=expected_cells, observed_cells=len(cells),
                        expected_probes=expected_probes, observed_probes=observed,
                        missing_fraction=missing_fraction, halts=halts)


def assert_no_extension(planned_n: int, requested_n: int) -> None:
    """No post-hoc sample growth: a larger N is a new preregistration, not an extension."""
    if requested_n > planned_n:
        raise DesignViolationError(
            f"planned N={planned_n}; refusing to extend to {requested_n}: "
            f"an indeterminate result requires a new preregistration"
        )


def assert_fixed_sample(cells: list[dict]) -> None:
    expected = N_CONVERSATIONS * REPS_PER_CELL
    if len(cells) != expected:
        raise DesignViolationError(f"fixed sample is {expected} cells; got {len(cells)}")


__all__ = ["MISSING_PROBE_LIMIT", "HealthReport", "evaluate_instrument_health",
           "assert_no_extension", "assert_fixed_sample"]
