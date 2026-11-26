"""Loud failures — one class per condition G1 must refuse.

The plan's rule is "no silent fallback, silent omission, or partial success"
(user G1 §5). Every failure mode therefore has its own type, so a test can prove
the *specific* refusal fired rather than matching on a message string.

All of them subclass `Mem1Error`, which subclasses `RuntimeError`: an uncaught
one aborts the run instead of degrading into a lower-quality result.
"""
from __future__ import annotations


class Mem1Error(RuntimeError):
    """Base class: the instrument refused to continue."""


class UnresolvedSpanError(Mem1Error):
    """A record or probe cited a span that does not exist in the evidence log."""


class DroppedSpanStillPresentError(Mem1Error):
    """A probe claims its evidence left the window, but the span is still present."""


class SeedMissingError(Mem1Error):
    """A cell has no seed, or seeds differ across arms for the same (conv, rep)."""


class SeedNotHonouredError(Mem1Error):
    """The runtime did not honour the seed; the cell is invalid and must be re-run."""


class TokenAccountingError(Mem1Error):
    """Token accounting is missing, non-positive, or exceeds the prompt capacity."""


class PlaceboFidelityError(Mem1Error):
    """Arm D could not be matched to A1 per turn (tokens, count, or shape)."""


class DecoyLeakError(Mem1Error):
    """A decoy answers or contradicts the probe: it is not a placebo."""


class TemporalVetoError(Mem1Error):
    """A superseded value was presented as current, or a current value as historical."""


class OwnershipVetoError(Mem1Error):
    """Cross-character leakage: a record outside the answering perspective was used."""


class HistoryErasedError(Mem1Error):
    """A superseded or revoked record became unrecoverable."""


class ProvenanceFabricationError(Mem1Error):
    """A formed record's source span does not belong to this conversation."""


class ReceiptFieldError(Mem1Error):
    """A required receipt field is missing or empty."""


class FrozenIdentityError(Mem1Error):
    """A frozen artifact's identity does not match the recorded one."""


class InvalidPairScoredError(Mem1Error):
    """An invalid or excluded conversation pair was included in scoring."""


class ProhibitedAggregateError(Mem1Error):
    """A composite score or a structurally undefined metric was requested."""


class ArmInertError(Mem1Error):
    """Arm A2 was requested while its proposals remain unapproved."""


class DesignViolationError(Mem1Error):
    """A frozen design parameter was violated (N, reps, probe quota, rollovers)."""


class ProtectedEvidenceError(DesignViolationError):
    """A write would reach protected historical evidence, or its identity is unknown."""


class ExecutionNotAuthorised(Mem1Error):
    """Something would have reached a model during G1. G4 is not authorized."""


__all__ = [name for name in dir() if name.endswith("Error") or name == "ExecutionNotAuthorised"]
