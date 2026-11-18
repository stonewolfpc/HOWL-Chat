"""Record classes for the v9 architecture.

Four classes, one per concept in v9 §2:

  * ``claim``     — what someone said, or a recollection they voiced. Not truth-apt.
  * ``account``   — what a character holds, with origin and evidence. Not truth-apt.
  * ``belief``    — what a character accepts. Not truth-apt, never promoted as truth.
  * ``world_truth`` / ``scenario_truth`` — the only truth-apt classes.

Every record carries its containing chat (rule 2), and every private record carries
exactly one owning character (rule 3). Those are separate fields because v9 §5 makes
them separate dimensions; the G1 ``owner_id`` field conflates them, which is why the
v9 engine does not reuse it.

Proposition identity (REV 3 §5.2, REV 4 §4.3) is a single field, ``proposition_id``:

  * **required** on ``world_truth`` / ``scenario_truth`` — every authoritative record
    belongs to exactly one declared identity;
  * **refused** on ``claim``, ``account`` and ``belief`` — a private record cannot be
    attached to a proposition, so enumeration has nothing to traverse. This absence,
    not a filter, is what keeps private content out of every canon surface;
  * **fixed at construction** — no operation re-points a record, so there is no late
    attachment, backfill or migration path.

Belief is listed under ``record_class`` as ``belief`` and carries an ``attribution``
naming the declaring character and the declaring role. Only the character's own model
may create one (rule 53). The engine enforces that; this module only refuses
malformed shapes, so a bad record cannot be constructed even directly.
"""
from __future__ import annotations

import copy
from dataclasses import dataclass, field

from .errors import (ContainmentError, DisclosureError, OwnershipError, TruthAptnessError,
                     ContributionKindError, ProvenanceError, BeliefError)

# v9 §3 truth-aptness table.
TRUTH_CLASSES = ("world_truth", "scenario_truth")
NON_TRUTH_CLASSES = ("claim", "account", "belief")
RECORD_CLASSES = TRUTH_CLASSES + NON_TRUTH_CLASSES

# Rule 50: `via` answers how information entered knowledge. `none` exists only for
# own accounts (rule 51) and is never a synonym for "missing".
VIAS = ("witnessed", "verified", "read", "told", "inferred", "narrator_grant", "none")
ORIGIN_CLASSES = ("acquired", "own_account")

# Rule 26 / L9: declared contribution kinds. `Unspecified claim` is the recorded
# value when the model declared nothing -- it is data, not an inference.
CONTRIBUTION_KINDS = (
    "assertion", "recollection", "jest", "question", "instruction",
    "declaration", "Unspecified claim",
)

# A': contributor identity categories. These name *who* supplied the content, never
# the technology that generated it -- a model producing a character's words does not
# make the contributor a model. The set is closed: a third category is an owner
# decision, not a parameter.
CONTRIBUTOR_KINDS = ("character", "human")

# Rule 53: whose declaration counts. Only the model embodying the character.
BELIEF_DECLARERS = ("character_model",)

# REV 2 §3.1 / REV 4 §4.1: a proposition's kind is a closed set of two. There is no
# third value and no minting path for one, so `kind` cannot become a general-purpose
# truth category (C-2's caveat).
PROPOSITION_KINDS = ("claim_about_world", "belief_ascription")

ACTIVE = "active"
SUPERSEDED = "superseded"
REVOKED = "revoked"
STATUSES = (ACTIVE, SUPERSEDED, REVOKED)


@dataclass
class Record:
    record_id: str
    record_class: str
    chat_id: str
    text: str
    created_turn: int
    # Rule 3. None only for domain truth, which belongs to a world or scenario.
    owner_character: str | None = None
    # Rule 6. Applicability is domain + target, never ownership or knowledge.
    domain: tuple[str, str] | None = None
    target: str = ""
    # REV 3 §5.2 / REV 4 §4.3. Required on truth classes, refused on private
    # classes, fixed once the instance is constructed.
    proposition_id: str | None = None
    provenance: dict = field(default_factory=dict)
    # Accounts only (rules 50-52).
    origin_class: str | None = None
    via: str | None = None
    # Claims only (rules 25-26).
    contribution_kind: str | None = None
    contributor: str | None = None
    # Claims only (rule 46 / J3 / B2). A disclosure event's established recipient set,
    # supplied by the conversation layer as structured metadata
    # ``{"recipients": (...), "source_class": ...}``. ``None`` on every record that is
    # not a disclosure claim, so its presence *is* the marker: no reader has to infer
    # from text or provenance whether a claim recorded an event.
    disclosure: dict | None = None
    # Beliefs only (rules 53, 55, 57).
    attribution: dict | None = None
    # Authoritative lifecycle (rules 20, 37, 42).
    status: str = ACTIVE
    supersedes: str | None = None
    superseded_by: str | None = None
    promotion_id: str | None = None
    corrected_by: list = field(default_factory=list)
    _sealed: bool = field(default=False, repr=False, compare=False)

    def __post_init__(self) -> None:
        if self.record_class not in RECORD_CLASSES:
            raise TruthAptnessError(f"unknown record class: {self.record_class!r}")
        if not self.chat_id:
            raise ContainmentError("record has no containing chat (rule 2)")
        if self.status not in STATUSES:
            raise OperationErrorFactory(self.status)

        if self.record_class in TRUTH_CLASSES:
            # Rule 17/24: authority records belong to a domain, not a character.
            if self.owner_character is not None:
                raise OwnershipError("a truth record cannot be owned by a character")
            if self.domain is None:
                raise TruthAptnessError("a truth record must declare its domain")
            if not self.promotion_id:
                raise TruthAptnessError("a truth record exists only through a promotion")
            # REV 3 §5.2: identity is declared, never implied by applicability.
            if not self.proposition_id:
                raise TruthAptnessError(
                    "a truth record must name its proposition (REV 3 §5.2)")
        else:
            # REV 2 §6.1: an ambient contribution is domain-scoped authorship, not a
            # private record -- `claim` + `declaration` + an authorship provenance
            # block. It is therefore owned by nobody, and this is the only shape in
            # which a non-truth record may lack an owner.
            authorship = (self.provenance.get("authorship")
                          if isinstance(self.provenance, dict) else None)
            if isinstance(authorship, dict) and authorship.get("actor"):
                if self.record_class != "claim":
                    raise OwnershipError(
                        "only a claim may be domain-authored (REV 2 §6.1)")
                if self.owner_character is not None:
                    raise OwnershipError(
                        "an authored ambient record is domain-scoped, not owned by a "
                        "character (REV 2 §6.1)")
                if not authorship.get("domain"):
                    raise OwnershipError(
                        "an authored ambient record must name its domain (rule 23)")
            else:
                # Rules 3, 5: private records name exactly one owning character.
                if not self.owner_character:
                    raise OwnershipError("a private record must name its owning character")
            # REV 4 §4.3: the edge does not exist, so no canon surface can traverse it.
            if self.proposition_id is not None:
                raise OwnershipError(
                    "a private record cannot be attached to a proposition (REV 4 §4.3)")

        if self.record_class == "account":
            if self.origin_class not in ORIGIN_CLASSES:
                raise ProvenanceError("account has no origin_class (rule 50)")
            if self.origin_class == "own_account":
                # Rule 51.
                if self.via != "none":
                    raise ProvenanceError("an own account must be recorded via=none")
            elif self.via not in VIAS or self.via == "none":
                raise ProvenanceError("an acquired account needs a real via (rule 50)")

        if self.record_class == "belief":
            if not self.attribution or not self.attribution.get("declarer"):
                raise BeliefError("a belief must name its declaring character")
            if self.attribution.get("declarer") != self.owner_character:
                raise BeliefError("a belief belongs to the character whose model declared it")
            if self.attribution.get("role") not in BELIEF_DECLARERS:
                raise BeliefError("only the character's own model may declare a belief (rule 53)")

        if self.record_class == "claim":
            if self.contribution_kind is None:
                raise ContributionKindError("claim has no recorded contribution kind")
            if self.contribution_kind not in CONTRIBUTION_KINDS:
                raise ContributionKindError(
                    f"kind {self.contribution_kind!r} is not a declared contribution kind")

        # J3 / B2: a disclosure claim cannot exist without the recipient set it was
        # given. The refusal is here, not only in the operation, so a malformed event
        # cannot be constructed directly either.
        if self.disclosure is not None:
            if self.record_class != "claim":
                raise DisclosureError(
                    "only a claim may carry disclosure metadata (rule 46)")
            named = self.disclosure.get("recipients")
            if not isinstance(named, tuple) or not named:
                raise DisclosureError(
                    "a disclosure claim must name the recipient set it was given (J3)")
            for who in named:
                if not isinstance(who, str) or not who:
                    raise DisclosureError("a recipient must be a named character")
            if len(set(named)) != len(named):
                raise DisclosureError("the recipient set names the same character twice")
            if self.owner_character in named:
                raise DisclosureError("the discloser cannot be their own recipient")
            if self.disclosure.get("source_class") not in NON_TRUTH_CLASSES:
                raise DisclosureError(
                    "a disclosure records which private class was disclosed")

        # Sealed after validation: `proposition_id` cannot be re-pointed later.
        object.__setattr__(self, "_sealed", True)

    def __setattr__(self, name: str, value) -> None:
        if (name == "proposition_id" and getattr(self, "_sealed", False)
                and value != getattr(self, "proposition_id", None)):
            raise TruthAptnessError(
                "proposition_id is fixed at construction (REV 4 §4.3: no late attachment)")
        object.__setattr__(self, name, value)

    @property
    def truth_apt(self) -> bool:
        """Rule 7: only world and scenario records answer truth queries."""
        return self.record_class in TRUTH_CLASSES

    def as_dict(self) -> dict:
        return {
            "record_id": self.record_id, "record_class": self.record_class,
            "chat_id": self.chat_id, "text": self.text,
            "created_turn": self.created_turn, "owner_character": self.owner_character,
            "domain": list(self.domain) if self.domain else None, "target": self.target,
            "proposition_id": self.proposition_id,
            # S5: deep-copied, so a caller editing a serialised record cannot rewrite the
            # stored descriptor or the prior text retained by rules 20/44.
            "provenance": copy.deepcopy(self.provenance), "origin_class": self.origin_class,
            "via": self.via, "contribution_kind": self.contribution_kind,
            "contributor": self.contributor,
            "disclosure": ({"recipients": list(self.disclosure["recipients"]),
                            "source_class": self.disclosure["source_class"]}
                           if self.disclosure else None),
            "attribution": dict(self.attribution or {}),
            "status": self.status, "supersedes": self.supersedes,
            "superseded_by": self.superseded_by, "promotion_id": self.promotion_id,
            "corrected_by": copy.deepcopy(self.corrected_by),
        }


def OperationErrorFactory(status: str):
    from .errors import OperationError
    return OperationError(f"unknown record status: {status!r}")


__all__ = ["TRUTH_CLASSES", "NON_TRUTH_CLASSES", "RECORD_CLASSES", "VIAS",
           "ORIGIN_CLASSES", "CONTRIBUTION_KINDS", "CONTRIBUTOR_KINDS",
           "BELIEF_DECLARERS", "PROPOSITION_KINDS", "ACTIVE", "SUPERSEDED", "REVOKED",
           "STATUSES", "Record"]
