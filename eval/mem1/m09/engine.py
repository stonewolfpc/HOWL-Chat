"""MEMORY-09 v9 engine: operations over records, authority and disclosure.

Scope discipline, enforced by absence rather than by a flag:

  * no transfer, share or carryover operation exists (rule 45 / K27);
  * no operation assigns truth to a claim, account or belief (rule 8);
  * no operation creates an account from promotion (rule 18 / K4);
  * no operation creates a belief from anything except that character's own model
    declaration (rule 53);
  * no operation displaces a promotion without a human edit (rule 43 / D-A).

Where the architecture says "explicit", the code takes an explicit argument rather
than inferring: a disclosure is an event with recipients, a belief names its declarer,
a promotion names or mints its proposition. Nothing is deduced from prose (rule 52).

Two boundaries from the accepted snapshot contract (REV 13 §2-§5) are structural here:

  * **Class A refusals** return before any protected read. The one thing they may
    consult is authorization metadata -- ``AuthorityRegistry.require`` -- which is the
    grant store, and which contains no proposition, record, belief or text.
  * **Class B refusals** happen only after the gate, discard the attempted capture and
    publish nothing derived from it. ``_capture`` returns ``None`` in that case.

A6' (the optional domain/target catalogue) is **not implemented**: no catalogue exists,
so no pre-gate lookup is possible. REV 13 §3.2 states that dropping it changes no
contract statement.
"""
from __future__ import annotations

import copy
from contextlib import contextmanager

from .authority import HUMAN, AuthorityRegistry, BasisEntry, validate_basis
from .errors import (AuthorityError, AuthorityStateError, BeliefError, CitationError,
                     ContributionKindError, ContributorKindError, DisclosureError,
                     EvidenceError, GroupError, OperationError, OwnershipError,
                     PropositionError, ProvenanceError, TruthAptnessError)
from .records import (ACTIVE, CONTRIBUTION_KINDS, CONTRIBUTOR_KINDS, NON_TRUTH_CLASSES,
                      PROPOSITION_KINDS, REVOKED, SUPERSEDED, VIAS, Record)

# REV 7 §5.1 step 2 / REV 6 §3.2: named once, so the view dispatch and the gate
# cannot drift apart.
KNOWN_VIEWS = ("world", "belief_facts", "propositions")

# S-7 capture attempts. A version test proves admissibility only if every mutating
# call publishes and advances the version in one step; `_mutation` is what makes that
# true here, and the bound is what turns exhaustion into a defined refusal rather
# than a loop.
CAPTURE_ATTEMPTS = 3

# REV 3 §5.5, required verbatim in every multiplicity response (REV 4 §5.2).
DISCLAIMER = ("identity multiplicity — unresolved. These identities share an "
              "applicability. The engine cannot determine whether they denote the same "
              "proposition. This report is not evidence that the propositions differ, "
              "and not evidence that a claim is missing.")

# B2 / D-E1: a disclosure is the character's own utterance, so the claim's
# contributor kind is fixed by the operation rather than supplied. Nothing about the
# *content* or the *recipients* is fixed here -- both are inputs (rule 46).
_DISCLOSER_CONTRIBUTOR_KIND = "character"

# Rule 45 / K27: no transfer capability exists in this module, and none is stubbed.
# Review R1 shipped a callable `transfer()` that always refused; R2 removed it, because
# an exported callable is an exposed interface regardless of what it does. The absence
# is now structural: there is no symbol to call. This is a removal, not a new rule.


class Engine:
    """One chat's memory engine. Containment is the constructor, not a field check."""

    def __init__(self, chat_id: str, registry: AuthorityRegistry | None = None) -> None:
        if not chat_id:
            raise OperationError("an engine must be constructed for one chat (rule 2)")
        self.chat_id = chat_id
        self.authority = registry or AuthorityRegistry()
        self._records: dict[str, Record] = {}
        self._propositions: dict[str, dict] = {}
        self._promotions: dict[str, dict] = {}
        self._declarations: dict[str, dict] = {}
        self._groups: dict[str, dict] = {}
        self._seq = 0
        self._seq_decl = 0
        self._seq_group = 0
        self._seq_prop = 0
        # S-7 diagnostic: monotone, bumped once per mutating call, never published.
        self._version = 0
        self._depth = 0

    # -- internals ---------------------------------------------------------
    def _id(self, prefix: str) -> str:
        self._seq += 1
        return f"{prefix}-{self._seq:05d}"

    @contextmanager
    def _mutation(self):
        """Publish-then-advance: one version step per outermost mutating call.

        S-7's soundness condition requires that no reader observes a partial write
        under an unchanged version. Nested calls (``disclose`` -> ``acquire``,
        ``supersede_authority`` -> ``promote``) therefore bump once, at the outermost
        exit, after every write in the call has landed.
        """
        self._depth += 1
        try:
            yield
        finally:
            self._depth -= 1
            if self._depth == 0:
                self._version += 1

    def _add(self, rec: Record) -> Record:
        self._records[rec.record_id] = rec
        return rec

    def get(self, record_id: str) -> Record:
        try:
            return self._records[record_id]
        except KeyError:
            raise CitationError(f"record {record_id!r} does not exist in this chat") from None

    def all_records(self) -> list[Record]:
        return list(self._records.values())

    @staticmethod
    def _refusal(code: str) -> dict:
        """One fixed refusal envelope, both classes (S-9 / S-12).

        ``code`` names which boundary refused. The key is `code`, not `reason`: the
        refusal envelope's operative definition is the class split (REV 12 §2.2 and §7,
        restated as REV 13 §7 row 17m′ — "both carry only `code`"), and the accepted
        S-9/S-10/S-11 language speaks of codes throughout. REV 7 §5.1's earlier `reason`
        spelling is not restated by any later instrument.

        Nothing data-derived is added: no identity, no count, no record id, no version,
        and no grant.
        """
        return {"status": "refused", "code": code}

    @staticmethod
    def _record_view(rec: Record) -> dict:
        return {"record_id": rec.record_id, "text": rec.text,
                "created_turn": rec.created_turn,
                "promotion_id": rec.promotion_id}

    def _kind_of(self, rec: Record) -> str | None:
        """The proposition's kind, read from the registry -- the only carrier.

        REV 7 §3.2 requires every record of an identity to carry that identity's kind.
        Here the kind exists in exactly one place, so a record cannot disagree with its
        identity: the check is satisfied by construction rather than by comparison.
        """
        if rec.proposition_id is None:
            return None
        return self._propositions.get(rec.proposition_id, {}).get("kind")

    @staticmethod
    def _accepts(view: str, kind: str | None) -> bool:
        """Which question a record answers (REV 6 §3.2). Unknown view accepts nothing."""
        if view == "world":
            return kind == "claim_about_world"
        if view == "belief_facts":
            return kind == "belief_ascription"
        if view == "propositions":
            return kind in PROPOSITION_KINDS
        return False

    def _eligible_truth(self, domain: tuple[str, str], target: str) -> list[Record]:
        """Access only, kind-agnostic (REV 5 §3.1).

        Called by ``authoritative`` and by nothing else, so canon access has exactly
        one definition.
        """
        found = [r for r in self._records.values()
                 if r.truth_apt and r.owner_character is None
                 and r.domain == domain and r.target == target and r.status == ACTIVE]
        return sorted(found, key=lambda r: r.record_id)

    # -- claims (rules 25-27) ----------------------------------------------
    def claim(self, *, contributor: str, contributor_kind: str, text: str,
              turn: int, kind: str | None = None, span: str | None = None) -> Record:
        """Record what was said. Never what it implies.

        ``kind`` is the model's declaration. A missing declaration is recorded as
        ``Unspecified claim`` -- the system does not supply or infer a kind (rule 26).

        ``contributor_kind`` is the identity category of the contributor (A'). A human
        contributor is refused: a claim is a character-private record, and a human has
        no private memory stream for one to belong to (rule 3).
        """
        if kind is not None and kind not in CONTRIBUTION_KINDS:
            raise ContributionKindError(f"{kind!r} is not a declared contribution kind")
        # A': identity, not mechanism. Vocabulary is checked first, so an unknown
        # category is never reported as an ownership problem. The literal "human" is
        # used deliberately: authority.HUMAN is the same string on a different axis,
        # and referring to it here would couple contributor identity to holder kind.
        if contributor_kind not in CONTRIBUTOR_KINDS:
            raise ContributorKindError(
                f"{contributor_kind!r} is not a contributor identity category "
                f"{CONTRIBUTOR_KINDS} (A')")
        if contributor_kind == "human":
            raise OwnershipError(
                f"{contributor!r} contributes as a human, who has no character-private "
                "memory stream; a claim would make them a character owner (rule 3 / A')")
        recorded_kind = kind if kind is not None else "Unspecified claim"
        with self._mutation():
            rec = Record(
                record_id=self._id("claim"), record_class="claim", chat_id=self.chat_id,
                text=text, created_turn=turn, owner_character=contributor,
                contribution_kind=recorded_kind, contributor=contributor,
                provenance={"contributor": contributor, "contributor_kind": contributor_kind,
                            "span": span, "turn": turn},
            )
            return self._add(rec)

    # -- accounts (rules 9, 50-52) -----------------------------------------
    def acquire(self, *, character: str, text: str, turn: int, via: str,
                span: str | None = None, origin_class: str = "acquired",
                descriptor: dict | None = None) -> Record:
        """Record what a character holds, after checking that kind's evidence.

        ``origin_class='own_account'`` forces ``via='none'`` and is exempt from the
        evidence table (rule 51) -- the exemption is checked before the table, so a
        validator can never demand evidence for an absence that is the record's
        meaning. A descriptor names the source *as content*; it is never a resolvable
        pointer to another record (rule 12). The validator never reads domain truth.
        """
        desc = dict(descriptor or {})
        if origin_class == "own_account":
            via = "none"
        self._validate_evidence(character=character, via=via, origin_class=origin_class,
                                descriptor=desc, turn=turn)
        with self._mutation():
            rec = Record(
                record_id=self._id("acct"), record_class="account", chat_id=self.chat_id,
                text=text, created_turn=turn, owner_character=character,
                origin_class=origin_class, via=via,
                provenance={"descriptor": desc, "span": span, "turn": turn},
            )
            return self._add(rec)

    def _validate_evidence(self, *, character: str, via: str, origin_class: str,
                           descriptor: dict, turn: int) -> None:
        """Contract C-3: required evidence per acquisition kind, at write time.

        The only things read are: the arguments, this character's own records, and the
        declaration registry. ``authoritative`` is never called, which is what keeps
        provenance validation from becoming truth adjudication.
        """
        if origin_class == "own_account":
            return  # rule 51 exemption, checked before the table (C7)
        if via not in VIAS or via == "none":
            raise ProvenanceError(f"{via!r} is not an acquisition kind (rule 50)")

        if via == "witnessed":
            if not descriptor.get("event"):
                raise EvidenceError("a witnessed account needs the acquisition event it "
                                    "witnessed (C-3)")
            named = descriptor.get("character")
            if named is not None and named != character:
                raise EvidenceError("the witnessed event is not this character's (C-3)")
            return

        if via == "told":
            speaker = descriptor.get("speaker")
            if not speaker:
                raise EvidenceError("a told account needs the speaker (C-3)")
            if speaker == character:
                raise EvidenceError("a character cannot be told by itself (C-3)")
            if not descriptor.get("event"):
                raise EvidenceError("a told account needs the disclosure event (rule 11)")
            return

        if via == "read":
            if not descriptor.get("source"):
                raise EvidenceError("a read account needs its document or asset (C4)")
            if not descriptor.get("event"):
                raise EvidenceError("a read account needs the reading event (C-3)")
            return

        if via == "inferred":
            basis = descriptor.get("basis")
            if not isinstance(basis, (list, tuple)) or not basis:
                raise EvidenceError("an inference needs a derivation basis (C2)")
            held = {r.record_id for r in self._records.values()
                    if r.owner_character == character}
            for ref in basis:
                if ref not in held:
                    raise EvidenceError(
                        f"basis {ref!r} is not held by {character!r}; an inference cannot "
                        "manufacture access (C-3)")
            return

        if via == "verified":
            if not descriptor.get("verification"):
                raise EvidenceError("a verified account needs its verification event (C-3)")
            if not descriptor.get("basis"):
                raise EvidenceError("a verified account needs its basis (C-3)")
            return

        if via == "narrator_grant":
            decl_id = descriptor.get("declaration")
            decl = self._declarations.get(decl_id) if decl_id else None
            if decl is None:
                raise EvidenceError("a coverage account needs a declaration (B1 / B7)")
            if decl.get("withdrawn_turn") is not None:
                raise EvidenceError("a withdrawn declaration confers no coverage (B6)")
            return

        raise ProvenanceError(f"{via!r} has no evidence rule (C-3)")

    # -- beliefs (rules 15, 53-57) -----------------------------------------
    def declare_belief(self, *, character: str, declarer: str, role: str,
                       text: str, turn: int) -> Record:
        """Create a belief. Only the character's own model may (rule 53 / D8).

        ``role`` must be ``character_model`` and ``declarer`` must be that character.
        An utterance never reaches this method: claims and beliefs are separate
        operations, so speech cannot become psychology by accident (D8).
        """
        if role != "character_model":
            raise BeliefError("only the model embodying the character may declare a belief")
        if declarer != character:
            raise BeliefError(
                f"{declarer!r} cannot install a belief for {character!r} (rule 53 / D9)")
        with self._mutation():
            prior = [r for r in self._records.values()
                     if r.record_class == "belief" and r.owner_character == character
                     and r.status == ACTIVE]
            rec = Record(
                record_id=self._id("belief"), record_class="belief", chat_id=self.chat_id,
                text=text, created_turn=turn, owner_character=character,
                attribution={"declarer": declarer, "role": role},
                provenance={"turn": turn},
            )
            for old in prior:
                # Rule 55: revision is append-only; the prior belief stays retrievable.
                old.status, old.superseded_by = SUPERSEDED, rec.record_id
                rec.supersedes = old.record_id
            return self._add(rec)

    def revoke_belief(self, *, character: str, declarer: str, role: str,
                      text: str, turn: int, record_id: str | None = None) -> Record:
        """Explicit belief revocation, model-gated and append-only (D-B, D5).

        The revocation is a *new* belief record, so nothing is deleted and the prior
        belief remains retrievable by its owner. Authority in a domain is promotion
        authority; it confers nothing over a character's interior state (O-5).
        """
        if role != "character_model":
            raise BeliefError("only the model embodying the character may revoke a belief")
        if declarer != character:
            raise BeliefError(
                f"{declarer!r} cannot revoke a belief of {character!r} (rule 53 / D-B)")
        with self._mutation():
            prior = [r for r in self._records.values()
                     if r.record_class == "belief" and r.owner_character == character
                     and r.status == ACTIVE]
            if record_id is not None:
                prior = [r for r in prior if r.record_id == record_id]
            if not prior:
                raise BeliefError("no active belief of that character to revoke")
            rec = Record(
                record_id=self._id("belief"), record_class="belief", chat_id=self.chat_id,
                text=text, created_turn=turn, owner_character=character,
                attribution={"declarer": declarer, "role": role, "revokes": True,
                             "revoked_records": [r.record_id for r in prior]},
                provenance={"turn": turn},
            )
            for old in prior:
                old.status, old.superseded_by = SUPERSEDED, rec.record_id
            rec.supersedes = prior[-1].record_id
            return self._add(rec)

    # -- promotion (rules 24, 28-37, 41-44) --------------------------------
    def promote(self, *, actor: str, domain: tuple[str, str], target: str, text: str,
                proposition: dict, basis: list[BasisEntry], turn: int,
                rationale: str | None = None) -> dict:
        """Create an authoritative world/scenario record under a recorded authority.

        Displacement is **not** available here (D-A): every promotion creates a record
        and displaces nothing. ``supersede_authority`` is the human-gated operation
        that changes which record governs. The promotion must name or mint its
        proposition (REV 3 §5.2); the result reports what else exists on that
        applicability, without labels or texts (REV 4 §5.3).
        """
        with self._mutation():
            self.authority.require(domain, actor, "promotion")
            entries = validate_basis(basis)
            cited: list[Record] = []
            for entry in entries:
                if entry.kind == "record_cited":
                    record = self.get(entry.record_id)
                    if record.chat_id != self.chat_id:
                        raise CitationError("a citation cannot resolve across chats (rule 33)")
                    if record.owner_character and record.owner_character != actor:
                        raise CitationError(
                            f"{entry.record_id} is private to {record.owner_character}; "
                            "the promotion cannot retrieve it (rule 33 / K14)")
                    cited.append(record)

            if domain[0] not in ("world", "scenario"):
                raise TruthAptnessError(f"domain kind {domain[0]!r} holds no truth")
            record_class = "world_truth" if domain[0] == "world" else "scenario_truth"

            pid, spec = self._resolve_proposition(domain, target, actor, proposition, turn)
            for record in cited:
                # Rule 56 / REV 3 §5.4: a cited belief may only be promoted as a fact
                # about its owner, through an ascription identity naming that owner.
                if record.record_class == "belief":
                    if (spec.get("kind") != "belief_ascription"
                            or spec.get("about_character") != record.owner_character):
                        raise BeliefError(
                            "a cited belief may only be promoted as a belief_ascription "
                            f"about {record.owner_character!r} (rule 56 / C-2)")

            coexisting = sorted(
                other for other, entry in self._propositions.items()
                if entry.get("domain") == domain and entry.get("target") == target
                and other != pid)

            if pid not in self._propositions:
                self._propositions[pid] = spec

            promo_id = self._id("promo")
            rec = Record(
                record_id=self._id("truth"), record_class=record_class, chat_id=self.chat_id,
                text=text, created_turn=turn, domain=domain, target=target,
                proposition_id=pid, promotion_id=promo_id,
                provenance={"actor": actor, "authority": list(domain), "turn": turn,
                            "rationale": rationale},
            )
            self._promotions[promo_id] = {
                "promotion_id": promo_id, "actor": actor, "domain": domain,
                "target": target, "turn": turn, "record_id": rec.record_id,
                "proposition_id": pid, "rationale": rationale,
                "basis": [dict(vars(e)) for e in entries], "superseded": [],
            }
            self._add(rec)
            return {"record": rec, "promotion_id": promo_id, "proposition_id": pid,
                    "coexisting_propositions": {"count": len(coexisting), "ids": coexisting},
                    "superseded": []}

    def _resolve_proposition(self, domain: tuple[str, str], target: str, actor: str,
                             proposition: dict, turn: int) -> tuple[str, dict]:
        """Reference an existing identity or mint one. Never inferred (REV 3 §5.2)."""
        if not isinstance(proposition, dict) or not proposition:
            raise PropositionError(
                "a promotion must declare its proposition: {'id': ...} or {'new': {...}}")
        if "id" in proposition:
            pid = proposition["id"]
            spec = self._propositions.get(pid)
            if spec is None:
                raise PropositionError(f"unknown proposition {pid!r}")
            if spec.get("chat_id") not in (None, self.chat_id):
                raise PropositionError("a proposition belongs to one chat (rule 31)")
            if spec.get("domain") != domain:
                raise PropositionError(
                    "a proposition belongs to one domain; it cannot be referenced from "
                    "another (rule 35)")
            if spec.get("target") != target:
                raise PropositionError(
                    "the promotion's applicability disagrees with the proposition's target")
            return pid, spec

        if "new" in proposition:
            spec_in = dict(proposition.get("new") or {})
            label = spec_in.get("label")
            kind = spec_in.get("kind")
            if not isinstance(label, str) or not label:
                raise PropositionError("minting requires a label, recorded verbatim")
            if kind not in PROPOSITION_KINDS:
                raise PropositionError(
                    f"{kind!r} is not a declared proposition kind ({PROPOSITION_KINDS}); a "
                    "new kind is an owner decision, not a parameter")
            about = spec_in.get("about_character")
            if kind == "belief_ascription":
                if not about:
                    raise PropositionError(
                        "a belief_ascription requires about_character (C-2)")
            elif about is not None:
                raise PropositionError(
                    "about_character applies only to a belief_ascription (C-2)")
            self._seq_prop += 1
            pid = f"prop-{self._seq_prop:05d}"
            spec = {"id": pid, "chat_id": self.chat_id, "domain": domain, "target": target,
                    "kind": kind, "about_character": about, "label": label,
                    "created_by": actor, "created_turn": turn}
            return pid, spec

        raise PropositionError(
            "a promotion must declare its proposition: {'id': ...} or {'new': {...}}")

    def supersede_authority(self, *, record_id: str, new_text: str, actor: str,
                            turn: int, basis: list[BasisEntry],
                            rationale: str | None = None) -> dict:
        """Explicit supersession by a human operator; both states survive (K23).

        This is the only displacement path (D-A). The write pair -- new record, then
        the prior record's status -- completes inside one mutation, so no reader can
        capture the half-applied state that REV 11 §2 traced.
        """
        with self._mutation():
            old = self.get(record_id)
            if not old.truth_apt:
                raise OperationError("supersession applies to authoritative records")
            self._require_human(old, actor)
            if old.domain is None:
                raise TruthAptnessError("authoritative record has no domain")
            if old.proposition_id is None:
                raise PropositionError("authoritative record has no proposition")
            result = self.promote(actor=actor, domain=old.domain, target=old.target,
                                  text=new_text, proposition={"id": old.proposition_id},
                                  basis=basis, turn=turn, rationale=rationale)
            new = result["record"]
            old.status, old.superseded_by = SUPERSEDED, new.record_id
            new.supersedes = old.record_id
            self._promotions[new.promotion_id]["superseded"].append(old.record_id)
            result["superseded"] = [old.record_id]
            return result

    def correct(self, *, record_id: str, new_text: str, actor: str, turn: int) -> Record:
        """Correction changes content and never settles an authority conflict (J15)."""
        with self._mutation():
            rec = self.get(record_id)
            if not rec.truth_apt:
                raise OperationError("correction applies to authoritative records (rule 42)")
            self._require_human(rec, actor)
            rec.corrected_by.append({"actor": actor, "turn": turn, "prior_text": rec.text})
            rec.text = new_text
            return rec

    def revoke_promotion(self, *, record_id: str, actor: str, turn: int) -> Record:
        """Explicit revocation; history and provenance are retained (rules 37, 43)."""
        with self._mutation():
            rec = self.get(record_id)
            if not rec.truth_apt:
                raise OperationError("revocation applies to a promotion")
            self._require_human(rec, actor)
            rec.status = REVOKED
            return rec

    def retract(self, record_id: str, turn: int) -> Record:
        """Retraction changes only this record's status (rule 42 / E1a).

        Accounts held by other characters are untouched (A2a), no notice account is
        created (E1c). Beliefs are refused outright: a belief's lifecycle belongs to
        the character's own model, so retraction cannot reach it (D-B / O-5).
        """
        with self._mutation():
            rec = self.get(record_id)
            if rec.record_class == "belief":
                raise BeliefError(
                    "a belief is not retractable here; belief lifecycle is model-gated "
                    "(D-B) -- use revoke_belief")
            if rec.record_class in ("world_truth", "scenario_truth"):
                raise OperationError("retract applies to non-authoritative records; "
                                     "use revoke_promotion for authority (rule 42)")
            rec.status = REVOKED
            return rec

    def _require_human(self, rec: Record, actor: str) -> None:
        """Rule 43 / K26: no actor changes authority state by disagreement."""
        if rec.domain is None:
            raise AuthorityError("authority state requires a domain")
        holders = [g for g in self.authority.grants_for(rec.domain)
                   if g.holder == actor and g.permits("promotion")]
        if not holders:
            raise AuthorityError(f"{actor!r} holds no authority in {rec.domain!r}")
        if all(g.holder_kind != HUMAN for g in holders):
            raise AuthorityStateError(
                "authority changes are human-controlled; a model grant cannot "
                "revoke or supersede (rule 43 / K26)")

    # -- disclosure (rules 10-12, 46-49, 57) -------------------------------
    def disclose(self, *, discloser: str, content: str, recipients: list[str],
                 source_record_id: str, turn: int, span: str | None = None,
                 kind: str | None = None) -> tuple[Record, tuple[Record, ...]]:
        """Record an established disclosure event, then give each exposed character
        its own account (rules 11-12, 46-48; J3).

        B2 (owner, 2026-11-17): the conversation layer establishes every fact written
        here. ``content`` and ``recipients`` are **inputs** -- neither is derived from
        the source record, from its owner's holdings, from group membership or from
        proximity. The recipient set is stored once, on the claim, as structured
        metadata, and that same supplied sequence creates the accounts: there is no
        second set and no reconstruction of one (§5.3).

        The source record is read only to check that the discloser may disclose it
        (rule 46: only its owner) and to record which private class was disclosed. Its
        id is never stored on the claim and never reaches a recipient (rule 12 / J8):
        the recipient's descriptor is speaker, class, event and turn.

        Every refusal is decided before the first write, so a refused disclosure leaves
        neither a claim nor an account behind. The nested ``acquire`` calls ride the
        one enclosing ``_mutation``, which advances the version once (S-7).
        """
        if kind is not None and kind not in CONTRIBUTION_KINDS:
            raise ContributionKindError(f"{kind!r} is not a declared contribution kind")
        if not isinstance(content, str) or not content.strip():
            raise DisclosureError(
                "a disclosure needs the communicated content the conversation layer "
                "established (rule 46 / B2)")

        supplied = list(recipients)
        if not supplied:
            raise DisclosureError("a disclosure needs at least one recipient (rule 11)")
        for who in supplied:
            if not isinstance(who, str) or not who:
                raise DisclosureError("a recipient must be a named character (rule 11)")
        if len(set(supplied)) != len(supplied):
            raise DisclosureError("the recipient set names the same character twice")
        if discloser in supplied:
            raise DisclosureError("the discloser cannot be their own recipient")

        src = self.get(source_record_id)
        if src.chat_id != self.chat_id:
            raise DisclosureError("a disclosure cannot leave the chat (rule 31)")
        if src.owner_character != discloser:
            raise DisclosureError(
                f"{discloser!r} does not hold {source_record_id}; only its owner may disclose")
        if src.record_class not in NON_TRUTH_CLASSES:
            raise DisclosureError("domain truth is not disclosed from a private record")

        event = {"kind": "disclosure", "turn": turn}

        def _account_descriptor() -> dict:
            """S3: a fresh descriptor per recipient.

            Each account owns its own event object, so mutating one recipient's
            provenance cannot reach another's. ``declare`` builds its descriptors this
            way already; this is the same shape.
            """
            return {"speaker": discloser, "record_class": src.record_class,
                    "turn": turn, "event": dict(event)}

        # Pre-flight: each recipient's evidence is checked before anything is written,
        # so the account loop cannot refuse after the claim has landed. This is what
        # makes the operation all-or-nothing without a rollback journal.
        for who in supplied:
            self._validate_evidence(character=who, via="told", origin_class="acquired",
                                    descriptor=_account_descriptor(), turn=turn)

        recorded_kind = kind if kind is not None else "Unspecified claim"
        recipients_set = tuple(supplied)
        with self._mutation():
            event_claim = self._add(Record(
                record_id=self._id("claim"), record_class="claim", chat_id=self.chat_id,
                text=content, created_turn=turn, owner_character=discloser,
                contribution_kind=recorded_kind, contributor=discloser,
                provenance={"contributor": discloser,
                            "contributor_kind": _DISCLOSER_CONTRIBUTOR_KIND,
                            "span": span, "turn": turn},
                disclosure={"recipients": recipients_set,
                            "source_class": src.record_class},
            ))
            accounts = tuple(self.acquire(
                character=who, text=content, turn=turn, via="told",
                descriptor=_account_descriptor(),
            ) for who in recipients_set)
            return event_claim, accounts

    def narration_report(self, *, established_event: str, text: str,
                         turn: int) -> dict:
        """Narration may report a disclosure already established (rule 47 / H8).

        There is deliberately no narration path that *creates* a disclosure, and none
        that returns another character's private record.
        """
        if not established_event:
            raise DisclosureError(
                "narration cannot manufacture a disclosure event (rule 47)")
        return {"reported": established_event, "text": text, "turn": turn}

    # -- declaration coverage (rules 10-11, B1-B9) -------------------------
    def declare_available(self, *, actor: str, domain: tuple[str, str], target: str,
                          audience: list[str], turn: int, record_id: str) -> dict:
        """Make one record generally available to an explicit list of characters.

        Coverage is a snapshot at declaration time: a newcomer after the declaration
        acquires nothing (B2), and a record created later is not covered (B5). A
        second declaration for the same target is refused rather than merged (B4).
        """
        with self._mutation():
            self.authority.require(domain, actor, "promotion")
            if not isinstance(audience, (list, tuple)) or not audience:
                raise DisclosureError("coverage is an explicit list of characters (B9)")
            if any(not isinstance(x, str) or not x for x in audience):
                raise DisclosureError("coverage list contains an invalid character id")
            src = self.get(record_id)
            if src.owner_character is not None:
                raise DisclosureError("a private record cannot be declared generally (B8)")
            # B3 (event-kind memories cannot be declared) has no representation here:
            # v9 carries no event/definition distinction on a record, so no guard is
            # written rather than a guard that can never fire. Recorded as an omission
            # in the implementation review. The owner's instruction not to expand the
            # contract quietly stands.
            for decl in self._declarations.values():
                if decl["target"] == target and decl["withdrawn_turn"] is None:
                    raise OperationError(
                        "a declaration for this target already exists; withdraw it or "
                        "create a new declaration (B4)")

            self._seq_decl += 1
            decl_id = f"decl-{self._seq_decl:04d}"
            # The declaration is recorded before the accounts it authorises, because a
            # `narrator_grant` account's evidence is that declaration resolving (C-3).
            decl = {"declaration_id": decl_id, "actor": actor, "domain": domain,
                    "target": target, "audience": list(audience), "turn": turn,
                    "record_id": record_id, "withdrawn_turn": None, "covered": []}
            self._declarations[decl_id] = decl
            accounts = [self.acquire(character=who, text=src.text, turn=turn,
                                     via="narrator_grant",
                                     descriptor={"declaration": decl_id, "target": target,
                                                 "event": {"kind": "declaration",
                                                           "turn": turn}})
                        for who in audience]
            decl["covered"] = [a.record_id for a in accounts]
            # S1: the caller gets a copy. `decl` stays the stored dict, so withdrawal
            # and coverage keep working on the engine's own object.
            return dict(decl, audience=list(decl["audience"]),
                        covered=list(decl["covered"]))

    def withdraw_declaration(self, *, declaration_id: str, actor: str,
                             turn: int) -> dict:
        """Future coverage stops; existing accounts remain (B6)."""
        with self._mutation():
            decl = self._declarations.get(declaration_id)
            if decl is None:
                raise OperationError(f"unknown declaration {declaration_id!r}")
            self.authority.require(decl["domain"], actor, "promotion")
            if decl["withdrawn_turn"] is not None:
                raise OperationError("declaration is already withdrawn")
            decl["withdrawn_turn"] = turn
            # W1: the caller gets a copy. The stored declaration keeps the withdrawal,
            # so `_validate_evidence`'s narrator_grant branch still sees it (B6).
            return dict(decl, audience=list(decl["audience"]),
                        covered=list(decl["covered"]))

    # -- ambient authorship (rule 39, U-5, U-8, REV 2 §6) ------------------
    def author_ambient(self, *, actor: str, contributor_kind: str,
                       domain: tuple[str, str], target: str,
                       text: str, turn: int, rationale: str | None = None,
                       revises: str | None = None) -> Record:
        """Author a domain-scoped ambient contribution. Never promotes.

        The shape is the one REV 2 §6.1 adopted: ``claim`` + declared kind
        ``declaration`` + an authorship provenance block. The record is owned by no
        character -- it is a contribution in a domain, not anyone's private memory --
        so it is not a private-record leak and cannot appear in a truth answer
        (rules 17, 18, 39 / K25). Authorship authority is a separate capability from
        promotion (U-8): holding one confers nothing of the other.

        ``contributor_kind`` is declared explicitly (A'): it is never derived from the
        actor's authority ``holder_kind``, because authority kind and contributor
        identity are separate axes.
        """
        if contributor_kind not in CONTRIBUTOR_KINDS:
            raise ContributorKindError(
                f"{contributor_kind!r} is not a contributor identity category "
                f"{CONTRIBUTOR_KINDS} (A')")
        if not isinstance(text, str) or not text.strip():
            raise OperationError("an authored contribution needs text (C-4 / U-5)")
        self.authority.require(domain, actor, "authorship")
        with self._mutation():
            prior = None
            if revises is not None:
                prior = self.get(revises)
                if prior.chat_id != self.chat_id:
                    raise OperationError("a revision cannot cross chats (rule 31)")
                if prior.truth_apt:
                    raise OperationError(
                        "a revision cannot target an authoritative record; changing "
                        "authority is supersession, not authorship (rule 42)")
            rec = Record(
                record_id=self._id("claim"), record_class="claim",
                chat_id=self.chat_id, text=text, created_turn=turn,
                owner_character=None, contribution_kind="declaration",
                contributor=actor, domain=domain, target=target,
                provenance={"contributor": actor, "contributor_kind": contributor_kind,
                            "turn": turn,
                            "authorship": {"actor": actor, "domain": list(domain),
                                           "target": target,
                                           "capability": "authorship", "turn": turn,
                                           "rationale": rationale}},
            )
            if prior is not None:
                # Append, preserve, and point the prior record at its successor.
                prior.status, prior.superseded_by = SUPERSEDED, rec.record_id
                rec.supersedes = prior.record_id
            return self._add(rec)

    # -- incompatibility assertion (U-7) -----------------------------------
    def assert_incompatible(self, *, record_a: str, record_b: str, applicability,
                            statement: str, actor: str, turn: int) -> dict:
        """Record an authorised assertion that two authoritative records cannot both stand.

        The engine records the assertion; it never invents one, and never promotes it
        into a verdict (rule 34). A cross-kind pair is refused before mutation: a
        relation that claims a world claim and a belief ascription answer one question
        would make a canon view adjudicate a belief (REV 7 §3.4).
        """
        if record_a == record_b:
            raise GroupError("an assertion needs two distinct records")
        first = self.get(record_a)
        second = self.get(record_b)
        if first.chat_id != self.chat_id or second.chat_id != self.chat_id:
            raise GroupError("an assertion cannot span chats (rule 31)")
        if not (first.truth_apt and second.truth_apt):
            raise GroupError("an assertion names authoritative records only")
        if first.domain is None or first.domain != second.domain:
            raise GroupError("an assertion applies within one domain (rule 35)")
        # The gate needs the domain, which is a property of the cited records, so this
        # one read precedes authorisation. It is limited to the two records the caller
        # named by id; nothing else is touched before `require` returns.
        self.authority.require(first.domain, actor, "promotion")
        targets = (applicability,) if isinstance(applicability, str) else tuple(applicability)
        if not targets or any(not isinstance(t, str) or not t for t in targets):
            raise GroupError("an assertion must name its shared applicability")
        if first.target != second.target or first.target not in targets:
            raise GroupError("both records must lie on the asserted applicability")
        if first.status != ACTIVE or second.status != ACTIVE:
            raise GroupError("an assertion names active records")
        if self._kind_of(first) != self._kind_of(second):
            raise GroupError(
                "cross_kind_group: a world claim and a belief ascription cannot be "
                "asserted incompatible (REV 7 §3.4)")
        for group in self._groups.values():
            if (group.get("status") == "open"
                    and {record_a, record_b} <= set(group.get("members", []))):
                raise GroupError("these records already share an open group")
        with self._mutation():
            self._seq_group += 1
            group_id = f"group-{self._seq_group:04d}"
            group = {"group_id": group_id, "domain": first.domain,
                     "targets": list(targets), "members": [record_a, record_b],
                     "propositions": [first.proposition_id, second.proposition_id],
                     "statement": statement, "asserted_by": actor,
                     "authority_ref": list(first.domain), "turn": turn,
                     "status": "open"}
            self._groups[group_id] = group
            # S2: the caller gets a copy. The stored group is the one every `status ==
            # "open"` read consults, so a returned handle must not be able to resolve it.
            return dict(group, targets=list(group["targets"]),
                        members=list(group["members"]),
                        propositions=list(group["propositions"]),
                        authority_ref=list(group["authority_ref"]))

    # -- capture (S-1...S-11, REV 12 A1-A11, REV 13 §2-§3) -----------------
    def _capture(self, domain: tuple[str, str], target: str, view: str) -> dict | None:
        """Obtain one admissible snapshot, or nothing (S-7 / S-9).

        Returns ``None`` when no admissible snapshot can be established; the caller
        turns that into a Class B refusal. `domain` and `target` are opaque labels
        throughout: no catalogue is consulted, and an unmatched label simply matches
        nothing (REV 13 §3).
        """
        for _ in range(CAPTURE_ATTEMPTS):
            before = self._version
            snapshot = self._take(domain, target, view)
            # Sound because `_mutation` advances the version exactly once per outermost
            # mutating call: a reader can never observe a partial write under an
            # unchanged version. Under concurrency this must be re-established; nothing
            # here demonstrates it.
            if before == self._version:
                return snapshot
        return None

    def _take(self, domain: tuple[str, str], target: str, view: str) -> dict:
        """A3-A9: whole-store candidates, total resolution, then relevance."""
        version = self._version

        # A3: the candidate set is every open group, unfiltered. Filtering first would
        # need resolution, and resolution is part of the capture (REV 12 §3.1).
        candidate_groups = [g for g in self._groups.values() if g.get("status") == "open"]

        # A4/A5: total, three-valued resolution over every member reference.
        member_refs = sorted({m for g in candidate_groups for m in g.get("members", [])})
        resolution: dict[str, dict] = {}
        for ref in member_refs:
            record = self._records.get(ref)
            if record is None:
                resolution[ref] = {"resolution": "unresolved", "reason": "record_absent",
                                   "proposition_id": None, "kind": None}
            else:
                resolution[ref] = {
                    "resolution": "resolved",
                    "proposition_id": record.proposition_id,
                    "kind": self._propositions.get(record.proposition_id or "", {}).get("kind"),
                    "record_class": record.record_class,
                }

        # A6/A6': eligibility and kind acceptance. No catalogue exists, so A6' is
        # omitted; dropping it changes no contract statement (REV 13 §3.2).
        eligible = self._eligible_truth(domain, target)
        accepted = [r for r in eligible if self._accepts(view, self._kind_of(r))]
        accepted_ids = {r.record_id for r in accepted}
        accepted_pids = {r.proposition_id for r in accepted}

        # A7: the frozen identity set, scoped by view.
        if view == "propositions":
            identity_ids = sorted(pid for pid, spec in self._propositions.items()
                                  if spec.get("domain") == domain
                                  and spec.get("target") == target)
        else:
            identity_ids = sorted(p for p in accepted_pids if p is not None)

        identities = []
        for pid in identity_ids:
            spec = self._propositions.get(pid) or {}
            item = {"id": pid, "kind": spec.get("kind")}
            if view in ("belief_facts", "propositions"):
                item["about_character"] = spec.get("about_character")
            if view == "propositions":
                item["label"] = spec.get("label")
                item["created_by"] = spec.get("created_by")
                item["created_turn"] = spec.get("created_turn")
                item["domain"] = list(domain)
                item["target"] = target
            identities.append(item)

        # A8: relevance, applied after resolution, against the frozen identity set.
        returned = set(identity_ids)
        groups = []
        for group in candidate_groups:
            members = list(group.get("members", []))
            pids = [resolution[m]["proposition_id"] for m in members]
            if any(p in returned for p in pids if p is not None):
                groups.append({"group_id": group.get("group_id"), "members": members,
                               "propositions": pids,
                               "targets": list(group.get("targets", [])),
                               "statement": group.get("statement"),
                               "asserted_by": group.get("asserted_by"),
                               "turn": group.get("turn")})

        # A9: one establishment map, computed once, against this capture.
        establishes: dict[str, bool] = {}
        causes: dict[str, str] = {}
        for group in groups:
            ok, cause = self._group_establishes(group, view, accepted_ids,
                                                accepted_pids, resolution)
            establishes[group["group_id"]] = ok
            if not ok and cause is not None:
                causes[group["group_id"]] = cause

        return {"version": version, "domain": domain, "target": target, "view": view,
                "accepted": accepted, "identities": identities,
                "returned_identities": identity_ids, "resolution": resolution,
                "groups": groups, "establishes": establishes, "causes": causes}

    @staticmethod
    def _group_establishes(group: dict, view: str, accepted_ids: set,
                           accepted_pids: set, resolution: dict) -> tuple[bool, str | None]:
        """REV 8 §3.2: accepted, visible, contained -- all three, in this view.

        An unresolved member fails visibility and containment alike, which is the
        fail-closed answer: a group naming a record that does not exist cannot be
        evaluated at all (REV 12 §3.2).
        """
        members = group["members"]
        if any(resolution[m]["resolution"] != "resolved" for m in members):
            return False, "member_unresolved"
        if not all(Engine._accepts(view, resolution[m]["kind"]) for m in members):
            return False, "spans_query_kinds"
        if not all(m in accepted_ids for m in members):
            return False, "member_not_accepted"
        if not all(resolution[m]["proposition_id"] in accepted_pids for m in members):
            return False, "member_not_accepted"
        return True, None

    def _assemble(self, snapshot: dict) -> dict:
        """Project one snapshot into the envelope. No read happens here (S-2)."""
        view = snapshot["view"]
        accepted = snapshot["accepted"]
        establishes = snapshot["establishes"]
        causes = snapshot["causes"]
        groups = snapshot["groups"]

        establishing = [g for g in groups if establishes.get(g["group_id"])]
        relation = "asserted" if establishing else "not_established"

        entries = []
        for ident in snapshot["identities"]:
            pid = ident["id"]
            records = sorted((r for r in accepted if r.proposition_id == pid),
                             key=lambda r: r.record_id)
            in_conflict = any(establishes.get(g["group_id"]) and pid in g["propositions"]
                              for g in groups)
            if in_conflict:
                state = "conflict"
            elif len(records) >= 2:
                state = "unresolved_relation"
            elif len(records) == 1:
                state = "settled"
            else:
                state = "none"
            entry = {"proposition_id": pid, "kind": ident.get("kind"), "state": state,
                     "records": [self._record_view(r) for r in records]}
            if ident.get("about_character") is not None:
                entry["about_character"] = ident["about_character"]
            if view == "propositions":
                entry["label"] = ident.get("label")
                entry["created_by"] = ident.get("created_by")
                entry["created_turn"] = ident.get("created_turn")
            entries.append(entry)

        integrity: list = []
        seen: set[str] = set()
        for group in groups:
            gid = group["group_id"]
            if establishes.get(gid):
                continue
            cause = causes.get(gid)
            if cause == "member_not_accepted":
                # History, not a fault: a retired member is silent in canon and listed
                # with `establishes: false` in the authority view (REV 8 §3.3).
                continue
            code = {"member_unresolved": "group_member_unresolved",
                    "spans_query_kinds": "group_spans_query_kinds"}.get(cause)
            if code is None:
                continue
            if view == "propositions":
                integrity.append({"code": code, "group_id": gid})
            elif code not in seen:
                seen.add(code)
                integrity.append(code)

        known = {ident["id"] for ident in snapshot["identities"]}
        if any(r.proposition_id not in known for r in accepted):
            # REV 5 §4 item 3: an eligible record attached to no enumerated identity is
            # a registry/record-store divergence, reported rather than quietly omitted.
            if view == "propositions":
                integrity.append({"code": "proposition_unregistered"})
            elif "proposition_unregistered" not in seen:
                seen.add("proposition_unregistered")
                integrity.append("proposition_unregistered")

        if view == "propositions":
            open_groups = []
            for group in groups:
                item = {"group_id": group["group_id"], "members": list(group["members"]),
                        "propositions": list(group["propositions"]),
                        "targets": list(group["targets"]),
                        "establishes": bool(establishes.get(group["group_id"]))}
                if not item["establishes"]:
                    item["cause"] = causes.get(group["group_id"])
                open_groups.append(item)
            establishing_view = [
                {"group_id": item["group_id"], "members": list(item["members"]),
                 "propositions": list(item["propositions"]),
                 "targets": list(item["targets"])}
                for item in open_groups if item["establishes"]]
        else:
            open_groups = [{"group_id": g["group_id"], "members": list(g["members"]),
                            "propositions": list(g["propositions"])}
                           for g in establishing]
            # Literal, identical lists in canon views (REV 9 §3.2): every canon group
            # establishes, so the two collections coincide by construction.
            establishing_view = [dict(item) for item in open_groups]

        if view == "propositions":
            status = "ok"
        elif not entries:
            status = "none"
        elif len(entries) == 1:
            status = entries[0]["state"]
        elif any(e["state"] == "conflict" for e in entries):
            status = "conflict"
        else:
            status = "multiple"

        envelope = {"status": status, "relation": relation, "view": view,
                    "open_groups": open_groups,
                    "establishing_groups": establishing_view,
                    "integrity": integrity, "propositions": entries}
        if view == "propositions":
            envelope["viewer"] = "authority"
            envelope["knowledge_conferred"] = False
        elif len(entries) == 1:
            # Present only for a single entry, so a partial set cannot be read as the
            # whole answer (REV 6 §4).
            envelope["records"] = entries[0]["records"]
        return envelope

    def authoritative(self, domain: tuple[str, str], target: str, *,
                      view: str, actor: str | None = None) -> dict:
        """Truth answer for one view. A conflict is returned as a conflict (J14).

        Order is the gate's invariant (REV 7 §5.1): unknown view, then authorisation
        for the one gated view, then -- and only then -- the first protected read.
        """
        if view not in KNOWN_VIEWS:
            return self._refusal("unknown_view")
        if view == "propositions":
            if actor is None:
                return self._refusal("no_domain_authority")
            try:
                self.authority.require(domain, actor, "promotion")  # metadata only
            except AuthorityError:
                # Uniform: undeclared domain, absent grant and actor=None are one
                # answer, and it carries no identity, count or id (REV 7 §5.1 / D-02 b8).
                return self._refusal("no_domain_authority")
        snapshot = self._capture(domain, target, view)
        if snapshot is None:
            return self._refusal("inconsistent_snapshot")
        return self._assemble(snapshot)

    # -- enumeration (REV 4 §5, gated) -------------------------------------
    def propositions(self, domain: tuple[str, str], target: str = "", *,
                     actor: str | None = None, viewer: str = "authority") -> dict:
        """Identities on one applicability. Authority audience only."""
        if viewer != "authority":
            return self._refusal("not_a_perspective_surface")
        if not target:
            return self._refusal("target_required")
        if actor is None:
            return self._refusal("no_domain_authority")
        try:
            self.authority.require(domain, actor, "promotion")
        except AuthorityError:
            return self._refusal("no_domain_authority")
        return self.authoritative(domain, target, view="propositions", actor=actor)

    def possible_duplicates(self, domain: tuple[str, str], *, actor: str | None = None,
                            target: str | None = None) -> dict:
        """Identity multiplicity report: ids and counts, never a verdict (REV 5 §5.2).

        Both forms read the registry, so both are gated (a count is information,
        D-02 behaviour 8) **and** version-bracketed: a non-refused response is a
        function of one coherent read (S-1, S-2, S-7), so a report assembled across
        two registry states would be a blend. If none can be established the answer is
        a Class B refusal, never a partial report (REV 12 §2.3 S-9).
        """
        if actor is None:
            return self._refusal("no_domain_authority")
        try:
            self.authority.require(domain, actor, "promotion")
        except AuthorityError:
            return self._refusal("no_domain_authority")
        rows = self._registry_read(domain, target)
        if rows is None:
            return self._refusal("inconsistent_snapshot")
        grouped: dict[str, list[str]] = {}
        for spec in rows:
            grouped.setdefault(spec.get("target", ""), []).append(spec.get("id", ""))
        groups = []
        for where in sorted(grouped):
            pids = sorted(grouped[where])
            if len(pids) < 2:
                continue
            item = {"ids": pids, "count": len(pids)}
            if target is not None:
                # Domain-wide mode names no applicability (REV 5 §5.2).
                item["target"] = where
            groups.append(item)
        return {"status": "ok", "viewer": "authority",
                "disclosure": "authorised_audit", "disclaimer": DISCLAIMER,
                "groups": groups}

    def _registry_read(self, domain: tuple[str, str],
                       target: str | None) -> list[dict] | None:
        """One coherent read of the identity registry (S-1, S-2, S-7).

        Consumes nothing else and returns frozen rows, so no field of a multiplicity
        report can be derived from a second read. The version test has the same
        soundness condition as `_capture`: `_mutation` advances the counter once per
        outermost mutating call, so no partial registry state is observable under an
        unchanged version. Exhaustion returns None, which the caller turns into a
        Class B refusal rather than a stale answer.
        """
        for _ in range(CAPTURE_ATTEMPTS):
            before = self._version
            rows = [dict(spec) for spec in self._propositions.values()
                    if spec.get("domain") == domain
                    and (target is None or spec.get("target") == target)]
            rows.sort(key=lambda s: (s.get("target", ""), s.get("id", "")))
            if before == self._version:
                return rows
        return None

    # -- reads -------------------------------------------------------------
    def accounts_of(self, character: str) -> list[Record]:
        """Knowledge is reconstructable from accounts alone (F1 / F1b)."""
        return [r for r in self._records.values()
                if r.record_class == "account" and r.owner_character == character]

    def beliefs_of(self, character: str) -> list[Record]:
        return [r for r in self._records.values()
                if r.record_class == "belief" and r.owner_character == character]

    def history(self, *, character: str | None = None,
                domain: tuple[str, str] | None = None, target: str = "") -> list[Record]:
        out = [r for r in self._records.values()
               if (character is None or r.owner_character == character)
               and (domain is None or r.domain == domain)
               and (not target or r.target == target)]
        return sorted(out, key=lambda r: (r.created_turn, r.record_id))

    def snapshot(self) -> dict:
        return {"chat_id": self.chat_id, "version": self._version,
                "records": [r.as_dict() for r in self._records.values()],
                "propositions": {k: dict(v) for k, v in self._propositions.items()},
                # S6: deep-copied. `propositions` needs no change -- its values are
                # scalars and one immutable tuple.
                "groups": {k: copy.deepcopy(v) for k, v in self._groups.items()},
                "promotions": {k: copy.deepcopy(v) for k, v in self._promotions.items()},
                "declarations": {k: copy.deepcopy(v) for k, v in self._declarations.items()}}


__all__ = ["Engine", "KNOWN_VIEWS", "CAPTURE_ATTEMPTS", "DISCLAIMER"]
