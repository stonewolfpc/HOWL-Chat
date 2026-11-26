# MEMORY-09 — Owner clarification: authority implementation discrepancy (REV 13 §7 row 17p)

**Status: DATED, PROSPECTIVE OWNER CLARIFICATION. Effective 2026-11-17. It supersedes one
implementation instruction — REV 13 §7 row 17p's "no change" — and nothing else. REV 13's historical text
is preserved byte-identical and is not rewritten.**

| Field | Value |
|---|---|
| Date | 2026-11-17 CST (−0600) |
| Type | Owner clarification, prospective from this date |
| Supersedes | **REV 13 §7, row 17p — the "no change" instruction**, for the current experimental implementation only |
| Does not supersede | REV 13's contract text; its §2–§5 snapshot contract; its §7 rows other than 17p; any rule, requirement, r32 mechanism or owner decision |
| Affects | `eval/mem1/m09/authority.py` and its callers in `eval/mem1/m09/engine.py` |
| Grants testing | **No.** No test, probe, compilation, linting or runtime verification is authorised by this record |
| Executed | Nothing |

---

## 1. Owner direction, recorded verbatim

> "I explicitly supersede the **no-change instruction** in REV 13 §7, row 17p, for the current MEMORY-09
> experimental implementation."

> "The reason is a verified source discrepancy discovered during Gate 2 implementation:
> - REV 13 §7 row 17p assumes `authority.py` already provides `require(domain, actor, operation)`.
> - The actual pre-change implementation provided `require(domain, actor)` and did not implement
>   `Grant.operations`.
> - Consequently, the premise that the existing authority implementation could be consulted unchanged was
>   not satisfied. The no-change instruction is therefore superseded for this implementation."

> "The implementation must not treat promotion authority as implicit authorship authority."

> "This clarification applies only to the implementation discrepancy described above."

---

## 2. The original REV 13 wording, quoted unchanged

From `MEMORY-09-ARCH-CONTRACTS-R2-FOUR-QUESTIONS-REV-13-2026-11-17.md` §7, row 17p, verbatim:

> | 17p | `authority.py` | — | **No change.** `require(domain, actor, operation)` — REV 3's signature — is consulted as-is; the metadata it reads is not edited or extended |

That row rests on a premise about the source, not on an instruction about the design. The premise was false.
The row is preserved exactly as written; what is superseded is its **operative effect** for this
implementation, not its text.

---

## 3. The actual pre-change source

File `eval/mem1/m09/authority.py`, 156 lines, sha256 **`07dfbeb219221190…`**, unchanged from the R2
baseline through every revision REV 3–13.

### 3.1 The missing capability field

Verbatim, pre-change:

```python
@dataclass
class Grant:
    """Standing promotion authority for one actor in one domain (rule 24)."""

    grant_id: str
    domain: tuple[str, str]
    holder: str
    holder_kind: str
    granted_by: str
    granted_turn: int
    revoked_turn: int | None = None

    def revoked(self) -> bool:
        return self.revoked_turn is not None
```

There is no `operations`. A `Grant` therefore could not express *which* capability it carried, so promotion
and authorship were indistinguishable by construction — not by oversight in a caller.

### 3.2 The two-parameter gate

Verbatim, pre-change:

```python
    def require(self, domain: tuple[str, str], actor: str) -> None:
        if domain not in self._domains:
            raise AuthorityError(f"domain {domain!r} has not been declared")
        if not self.active_holder(domain, actor):
            raise AuthorityError(
                f"{actor!r} holds no promotion authority in {domain!r} (rule 35 / K8)")
```

Two parameters, and the message names "promotion authority" as the only thing being checked — accurate for
the object, because it was the only capability that existed.

### 3.3 What row 17p assumed

| Row 17p's premise | Pre-change reality |
|---|---|
| `require(domain, actor, operation)` exists | `require(self, domain, actor)` — two parameters |
| The operation argument can be passed "as-is" | No such parameter; no operation vocabulary |
| An authorship capability can be gated | No capability representation on `Grant` |
| The metadata "is not edited or extended" | It had to be extended for any capability split to exist at all |

**Conclusion recorded here for the file:** row 17p's premise was not satisfied, so "consulted as-is" was not
an available implementation. This is a source discrepancy, not a design disagreement.

---

## 4. The decision

The no-change instruction in row 17p is **superseded, prospectively and narrowly**, for the current
experimental implementation in `eval/mem1/`. 

| Aspect | Disposition |
|---|---|
| Row 17p's text | **Preserved.** Not edited, not reinterpreted as having said something else |
| Row 17p's premise | Recorded as false at the time it was written |
| Row 17p's operative effect | **Superseded** for this implementation only |
| Every other REV 13 row | Untouched and still in force |
| REV 13 §2–§5 | Untouched; remains the operative snapshot contract |

The correction is prospective: it governs implementation from this date forward. It makes no claim about what
was authorised before 2026-11-17, and it does not retroactively authorise anything already done.

---

## 5. Authority requirements that remain operative and unchanged

The clarification does not weaken the architecture. These nine requirements were operative before it and
remain operative after it:

| # | Requirement |
|---|---|
| 1 | `Grant.operations` distinguishes promotion and authorship capabilities |
| 2 | The permitted operation vocabulary is `{"promotion", "authorship"}` |
| 3 | Domain declaration installs the human operator with both capabilities |
| 4 | Model grants default to promotion only |
| 5 | Authorship must be explicitly granted |
| 6 | `require(domain, actor, operation)` checks the requested capability |
| 7 | Revocation may remove one capability without removing the other |
| 8 | Authority remains scoped to its applicable domain |
| 9 | Grant origin and delegation provenance remain represented as required by the accepted architecture |

Plus the constraint stated in the direction itself: **promotion authority is never implicit authorship
authority**, and no caller may obtain one capability by presenting the other.

---

## 6. Final authority API and capability semantics

File `eval/mem1/m09/authority.py`, 236 lines, sha256 **`f344582350af88f1…`**.

| Symbol | Line | Semantics |
|---|---|---|
| `OPERATIONS = ("promotion", "authorship")` | 39 | Closed vocabulary. A third capability is an owner decision, not a parameter; `grant` and `require` both refuse an unknown operation |
| `DEFAULT_MODEL_OPERATIONS = frozenset({"promotion"})` | 43 | The default for a model grant (requirement 4) |
| `Grant.operations: frozenset` | 89 | The capability set carried by the grant |
| `Grant.revoked_operations: frozenset` | 91 | Capabilities withdrawn individually (requirement 7) |
| `Grant.revoked(operation=None)` | 93 | Whole-grant revocation with no argument; with an argument, "is this capability withdrawn?" |
| `Grant.permits(operation)` | 106 | `operation in operations and not revoked(operation)` |
| `grant(..., operations=("promotion",))` | 136 | Explicit grant; refuses an empty set, an unknown operation, an undeclared domain, an unknown holder kind |
| `revoke_grant(grant_id, turn, operation=None)` | 162 | Whole grant, or one capability only; refuses an operation the grant does not carry, and double revocation |
| `active_holder(domain, actor, operation="promotion")` | 190 | Domain-scoped capability check |
| `require(domain, actor, operation)` | 198 | **No default.** A missing argument is a `TypeError`, not a silent promotion check |
| `grants_for(domain)` | 213 | Read-only, used by `_require_human` to establish that an authority change comes from a human holding promotion capability |
| `declare_domain(...)` | — | Installs the human operator with `frozenset(OPERATIONS)` — one holder, both capabilities (requirement 3) |
| `Grant.granted_by` | — | Delegation origin, retained (requirement 9) |

Domain scoping (requirement 8) is unchanged: grants are keyed by the `(kind, id)` domain tuple, and a grant in
one domain authorises nothing in another.

---

## 7. Affected implementation file and callers

| File | Effect |
|---|---|
| `eval/mem1/m09/authority.py` | **Changed** under this clarification: 156 → 236 lines; `07dfbeb2…` → `f3445823…` |
| `eval/mem1/m09/engine.py` | **Changed** — every gate call now names its operation |
| `eval/mem1/m09/__init__.py` | **Changed** — exports `OPERATIONS`, `DEFAULT_MODEL_OPERATIONS` |
| `eval/mem1/m09/records.py` | Changed for the authored-claim shape, not for this clarification (see §8) |
| `eval/mem1/model.py`, the G1 harness, production Howl Chat | **Untouched** |

### 7.1 Every caller now passes the operation explicitly

`grep -n 'authority.require(' eval/mem1/m09/engine.py` → eight call sites, no two-argument call remaining:

| Line | Caller | Operation |
|---|---|---|
| 364 | `promote` | `"promotion"` |
| 615 | `declare_available` | `"promotion"` |
| 658 | `withdraw_declaration` | `"promotion"` |
| 679 | `author_ambient` | **`"authorship"`** |
| 731 | `assert_incompatible` | `"promotion"` |
| 1008 | `authoritative` (the gated view) | `"promotion"` |
| 1029 | `propositions` | `"promotion"` |
| 1047 | `possible_duplicates` | `"promotion"` |

One further capability check exists outside `require`: `_require_human` (engine.py:542-543) filters
`grants_for(domain)` by `g.permits("promotion")`, so a model grant cannot change authority state even if it
holds authorship.

### 7.2 Reading of the public surface

| Operation | Grant it requires | Can create |
|---|---|---|
| Promotion, supersession, correction, revocation, coverage, assertion, enumeration, multiplicity | `"promotion"` | authoritative records, coverage accounts, groups |
| Ambient authorship | `"authorship"` | a non-truth-apt `claim` only |

No call site derives or defaults its operation; each names it literally, so the split is auditable by reading
the call sites rather than by trusting a summary.

---

## 8. Interaction with the earlier correction instruction

The Gate 2 implementation was first performed under an instruction that forbade changing `authority.py`. The
owner then lifted that restriction (2026-11-17, recorded in
`MEMORY-09-GATE2-REPORT-SUPPLEMENT-2026-11-17.md` §1), and this clarification supplies the formal record for
the row that instruction contradicted.

| Change | Authority under which it was made |
|---|---|
| `Grant.operations`, `require(..., operation)`, per-capability revocation | This clarification (§4–§6) |
| `author_ambient` and its ownerless authored-claim shape | The owner's earlier correction instruction, implementing REV 2 §6.1 — **not** this clarification, which is confined to the authority file |
| `possible_duplicates` version bracket, refusal key `code` | The earlier correction instruction and the report supplement §5 |

This clarification adds no architecture, rule, field, requirement or probe beyond the accepted authority
design, and it authorises no operation that did not already exist in the accepted documents.

---

## 9. What this clarification does not do

| Not done | Note |
|---|---|
| No change to the 57 v9 rules | Count remains 57 |
| No change to the 118 conformance requirements | Count remains 118; no status moves |
| No amendment to r32, v9, the conformance map, plan v3, Amendment 1 | All byte-identical (§10) |
| No amendment to the accepted REV 13 snapshot contract | REV 13 is byte-identical; only row 17p's operative effect is superseded, and its text stands |
| No new architecture, field, requirement or probe beyond the accepted authority design | `OPERATIONS` is the accepted vocabulary, not an addition to it |
| No production Howl Chat integration | Untouched |
| **No testing authorisation** | No test, probe, self-test, suite, compilation, linting or runtime verification is authorised by this record, and none was performed |
| No transfer, destination-write, sharing or cross-chat carryover | No such symbol exists |

---

## 10. Verification at the time of this record

| Item | Value |
|---|---|
| REV 13 (preserved, not rewritten) | `5b745628…` |
| REV 1–12 | `98447423…`, `74ac6f77…`, `de4f68b6…`, `1da08f2c…`, `a7e7c013…`, `41421dce…`, `a281f3e5…`, `b38f1062…`, `5b3a9008…`, `7d703b77…`, `40850419…`, `0592d689…` — all unmoved |
| v9 / map / ledger / acceptance instrument | `698254ca…` / `b4ba33cc…` / `da3e7239…` / `9a25ded0…` — all unmoved |
| r32 / plan v3 / Amendment 1 | `e88ef2d1…` / `6c4078bd…` / `6afd2a0d…` — all unmoved |
| `authority.py` before / after | `07dfbeb2…` (156 lines) → `f3445823…` (236 lines) |
| Rules / map rows | 57 / 118 |
| Executed | **nothing** |

**Sequence:** this record → owner review → exact-rule source review → separate testing authorisation. This
record requests neither.
