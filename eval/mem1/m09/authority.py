"""Authority: domains, holders, bases and standing grants (v9 rules 24, 28-39).

Three separate ideas, deliberately not merged:

  * a **domain** — one world or one scenario (rule 35). Authority is held per domain.
  * a **holder** — the actor authorised to act there. The default holder in every
    domain is the human operator (K8). A model may hold *standing authority* only
    through an explicit, domain-scoped, recorded, revocable grant (rule 24 / L7).
  * a **basis** — a set of typed entries (rule 28). The registry validates structure
    and references. It never judges sufficiency (rule 34 / K15 / K20).

**Two capabilities, one regime** (rule 39, U-8, REV 2 §7). A grant carries a set of
operations drawn from ``OPERATIONS``:

  * ``promotion`` — may create authoritative world/scenario records;
  * ``authorship`` — may author ambient, non-authoritative contributions.

The regime is shared in structure (one holder source, one per-domain scoping, one grant
form, one revocation path) while capability is per operation. Generation never confers
promotion (rule 24), authorship never promotes (rule 39 / K25), and revoking one
capability leaves the other untouched (REV 2 §7, L10).

Nothing here decides whether a claim is true. It decides who may perform an operation,
which is a different question (rule 21).
"""
from __future__ import annotations

from dataclasses import dataclass, field

from .errors import AuthorityError, BasisError

HUMAN = "human"
MODEL = "model"

BASIS_KINDS = ("record_cited", "actor_declaration")

# Rule 39 / U-8: one authoring regime, two operations. The set is closed: a third
# capability is an owner decision, not a parameter.
OPERATIONS = ("promotion", "authorship")

# A model grant defaults to promotion only; authorship must name it explicitly
# (REV 2 §7). The human operator installed by a domain declaration holds both.
DEFAULT_MODEL_OPERATIONS = frozenset({"promotion"})


@dataclass(frozen=True)
class BasisEntry:
    """One typed basis entry.

    ``record_cited`` names a record id that must resolve inside the promoting chat.
    ``actor_declaration`` carries the declared content verbatim plus the identity and
    authority of the declaring actor (rule 30 / K17) -- never a bare string.
    """

    kind: str
    record_id: str | None = None
    content: str | None = None
    actor: str | None = None
    actor_authority: str | None = None
    rationale: str | None = None

    def __post_init__(self) -> None:
        if self.kind not in BASIS_KINDS:
            raise BasisError(f"basis entry kind {self.kind!r} is not one of {BASIS_KINDS}")
        if self.kind == "record_cited":
            if not self.record_id:
                raise BasisError("record_cited entry has no record id (K14)")
        else:
            if not self.content or not self.actor or not self.actor_authority:
                raise BasisError(
                    "actor_declaration must carry content, actor identity and authority (K17)")


@dataclass
class Grant:
    """Standing authority for one actor in one domain (rule 24 / REV 2 §7).

    ``operations`` is the capability set. ``granted_by`` is the delegation origin: who
    conferred the standing authority, retained so a grant can always be attributed
    rather than assumed.
    """

    grant_id: str
    domain: tuple[str, str]
    holder: str
    holder_kind: str
    granted_by: str
    granted_turn: int
    operations: frozenset = DEFAULT_MODEL_OPERATIONS
    revoked_turn: int | None = None
    revoked_operations: frozenset = frozenset()

    def revoked(self, operation: str | None = None) -> bool:
        """Whole-grant revocation when no operation is named.

        With an operation, this answers "is that capability withdrawn?" -- which is
        either the whole grant being revoked or that one capability being revoked
        (L10: revoking one leaves the other untouched).
        """
        if self.revoked_turn is not None:
            return True
        if operation is None:
            return False
        return operation in self.revoked_operations

    def permits(self, operation: str) -> bool:
        return operation in self.operations and not self.revoked(operation)


@dataclass
class AuthorityRegistry:
    """Per-domain holders. A grant in one domain permits nothing in another (K21)."""

    _grants: dict[tuple[str, str], list[Grant]] = field(default_factory=dict)
    _domains: set[tuple[str, str]] = field(default_factory=set)
    _seq: int = 0

    def declare_domain(self, domain: tuple[str, str], turn: int,
                       operator: str = "human_operator") -> None:
        """Declaring a world or scenario installs the human operator as its holder.

        K8 names the world/scenario authoring authority as the promoting actor and the
        human operator as its default. That default is *recorded* here as a grant, not
        assumed at call sites. The operator holds **both** capabilities (REV 2 §7):
        the declaration installs one authority with two operations.
        """
        if not isinstance(domain, tuple) or len(domain) != 2 or not all(domain):
            raise AuthorityError("a domain is (kind, id), both non-empty")
        self._domains.add(domain)
        self._seq += 1
        self._grants.setdefault(domain, []).append(Grant(
            grant_id=f"grant-{self._seq:04d}", domain=domain, holder=operator,
            holder_kind=HUMAN, granted_by=operator, granted_turn=turn,
            operations=frozenset(OPERATIONS)))

    def grant(self, domain: tuple[str, str], holder: str, holder_kind: str,
              granted_by: str, turn: int, operations=("promotion",)) -> Grant:
        """Rule 24 / L7: a model grant must be explicit and domain-scoped to exist.

        ``operations`` defaults to promotion only. Authorship never follows from a
        promotion grant; it must be named (REV 2 §7).
        """
        if domain not in self._domains:
            raise AuthorityError(f"domain {domain!r} has not been declared")
        if holder_kind not in (HUMAN, MODEL):
            raise AuthorityError(f"unknown holder kind {holder_kind!r}")
        requested = frozenset(operations)
        if not requested:
            raise AuthorityError("a grant must carry at least one operation")
        unknown = requested.difference(OPERATIONS)
        if unknown:
            raise AuthorityError(
                f"unknown operation(s) {sorted(unknown)!r}; the set is {OPERATIONS} "
                "(a new capability is an owner decision)")
        self._seq += 1
        g = Grant(grant_id=f"grant-{self._seq:04d}", domain=domain, holder=holder,
                  holder_kind=holder_kind, granted_by=granted_by, granted_turn=turn,
                  operations=requested)
        self._grants.setdefault(domain, []).append(g)
        return g

    def revoke_grant(self, grant_id: str, turn: int,
                     operation: str | None = None) -> Grant:
        """L10: revocation stops future action; it erases nothing already published.

        With no ``operation`` the whole grant is revoked. With one, only that
        capability is withdrawn and the other continues (REV 2 §7, L10).
        """
        for grants in self._grants.values():
            for g in grants:
                if g.grant_id == grant_id:
                    if operation is None:
                        if g.revoked_turn is not None:
                            raise AuthorityError(f"grant {grant_id} is already revoked")
                        g.revoked_turn = turn
                        return g
                    if operation not in OPERATIONS:
                        raise AuthorityError(
                            f"{operation!r} is not an operation ({OPERATIONS})")
                    if operation not in g.operations:
                        raise AuthorityError(
                            f"grant {grant_id} does not carry {operation!r}")
                    if g.revoked(operation):
                        raise AuthorityError(
                            f"grant {grant_id} has already withdrawn {operation!r}")
                    g.revoked_operations = frozenset(g.revoked_operations | {operation})
                    return g
        raise AuthorityError(f"unknown grant {grant_id}")

    def active_holder(self, domain: tuple[str, str], actor: str,
                      operation: str = "promotion") -> bool:
        """Is this actor currently authorised for this operation here (K8/K21)?"""
        for g in self._grants.get(domain, ()):
            if g.holder == actor and g.permits(operation):
                return True
        return False

    def require(self, domain: tuple[str, str], actor: str, operation: str) -> None:
        """The gate. Reads authorization metadata only; never protected memory (S-12).

        ``operation`` is required: there is no default, because a default would let a
        promotion grant silently authorise authorship (U-8, rule 39).
        """
        if operation not in OPERATIONS:
            raise AuthorityError(f"{operation!r} is not an operation ({OPERATIONS})")
        if domain not in self._domains:
            raise AuthorityError(f"domain {domain!r} has not been declared")
        if not self.active_holder(domain, actor, operation):
            raise AuthorityError(
                f"{actor!r} holds no {operation!r} authority in {domain!r} "
                "(rule 35 / K8 / U-8)")

    def grants_for(self, domain: tuple[str, str]) -> list[Grant]:
        return list(self._grants.get(domain, ()))


def validate_basis(entries: list[BasisEntry]) -> list[BasisEntry]:
    """Rule 28 / K13 / K19: a non-empty set of typed entries, mixed kinds allowed.

    Sufficiency is never assessed here (rule 34). A one-entry basis is accepted for
    the same reason a twenty-entry basis is: the authorised actor judged it.
    """
    if entries is None:
        raise BasisError("no basis supplied; an empty basis is rejected (K19)")
    if not isinstance(entries, (list, tuple)):
        raise BasisError("basis must be a set of entries, not a single choice")
    if len(entries) == 0:
        raise BasisError("empty basis rejected (rule 28 / K19)")
    for e in entries:
        if not isinstance(e, BasisEntry):
            raise BasisError(f"basis entry {e!r} is not a typed BasisEntry")
    return list(entries)


__all__ = ["HUMAN", "MODEL", "BASIS_KINDS", "OPERATIONS", "DEFAULT_MODEL_OPERATIONS",
           "BasisEntry", "Grant", "AuthorityRegistry", "validate_basis"]
