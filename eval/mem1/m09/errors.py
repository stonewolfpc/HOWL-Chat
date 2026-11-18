"""Failure classes for the MEMORY-09 v9 engine.

One class per refusal, for the same reason the G1 harness has one per condition:
a review must be able to say *which* rule refused, not parse a message. Every
class names the rule or rule family it enforces.
"""
from __future__ import annotations


class M09Error(RuntimeError):
    """Base: the engine refused an operation."""


class ContainmentError(M09Error):
    """Rule 2 — a record must be contained by exactly one chat."""


class OwnershipError(M09Error):
    """Rule 3 — a private record must have exactly one owning character."""


class TruthAptnessError(M09Error):
    """Rule 7 — only world and scenario records may hold truth."""


class AdjudicationError(M09Error):
    """Rule 8 — the system never marks a claim, account or belief true or false."""


class ProvenanceError(M09Error):
    """Rule 12 / 18 / 52 — unresolvable evidence, or an origin that was invented."""


class ContributionKindError(M09Error):
    """Rule 26 / L9 — kind is declared, never inferred; missing means Unspecified."""


class ContributorKindError(M09Error):
    """A' — a contributor identity category outside the closed set.

    ``contributor_kind`` names *who* supplied the content, never the technology that
    generated it (A' contract correction). Anchored to that contract rather than to a
    numbered v9 rule: no v9 rule gives contributor identity its own number, and
    inventing one would misstate the baseline. Distinct from ``ContributionKindError``
    (rule 26, content kind) and from ``ProvenanceError`` (rules 12 / 18 / 52, evidence
    origin): an invented contributor category is its own refusal.
    """


class AuthorityError(M09Error):
    """Rule 24 / 35 / K8 / K21 — no authority in this domain for this actor."""


class BasisError(M09Error):
    """Rule 28 / K13 / K19 / K20 — malformed basis. Never a relevance judgement."""


class CitationError(M09Error):
    """Rule 33 / K14 — a citation that does not resolve inside the promoting chat."""


class DisclosureError(M09Error):
    """Rules 10 / 11 / 47 — disclosure must be an event, and only exposure acquires."""


class BeliefError(M09Error):
    """Rules 15 / 53 / 54 / 56 — belief creation, installation and promotion."""


class AuthorityStateError(M09Error):
    """Rule 43 / K26 — a non-human actor tried to change shared authority state."""


class OperationError(M09Error):
    """Rule 42 — the wrong operation for the intended meaning."""


class HistoryError(M09Error):
    """Rules 20 / 44 / K23 — an operation would erase retained history."""


class RetrievalError(M09Error):
    """Rule F4 / F5 — unknown context, or a silent widening of an empty set."""


class EvidenceError(M09Error):
    """Contract C-3 — an acquisition whose own kind's evidence is absent.

    Raised at write time and never from a truth read: evidence establishes that
    information was acquired, never that it is correct (rule 9 / C5).
    """


class PropositionError(M09Error):
    """REV 3 §5.2-§5.4 — proposition identity: mint, reference, kind, scope."""


class GroupError(M09Error):
    """U-7 / REV 7 §3.4 — an incompatibility assertion that cannot be made."""


class SnapshotError(M09Error):
    """S-7 / S-9 — no admissible snapshot could be captured (Class B refusal).

    Capture is attempted and discarded; nothing derived from the partial read is
    published. Distinct from an authorization refusal, which reads no protected
    memory at all (S-5 / S-10).
    """


class ViewRefusal(M09Error):
    """REV 4 §3 — enumeration is not a character-perspective surface."""


__all__ = [name for name in dir() if name.endswith("Error")]
