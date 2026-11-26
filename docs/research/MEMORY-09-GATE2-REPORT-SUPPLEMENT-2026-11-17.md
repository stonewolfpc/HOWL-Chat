# MEMORY-09 — Gate 2 report supplement: correction authorisation and review resolutions

**This document supersedes specific statements in
`MEMORY-09-GATE2-IMPLEMENTATION-REPORT-2026-11-17.md` and resolves the five items the owner withheld
completion on. The earlier report is retained unedited as the record of what was true under the first
authorisation; §3 below lists every statement it made that no longer holds.**

| Field | Value |
|---|---|
| Date | 2026-11-17 CST (−0600) |
| Trigger | Owner correction: the `authority.py` prohibition "was too strong" |
| Effect | `authority.py` changed; `author_ambient` implemented; four review items resolved |
| Still unauthorised | Testing, probing, compilation, linting, runtime verification |
| Executed | Nothing |

---

## 1. The correction instruction, recorded

> "The `authority.py` blocker is partly my fault. My implementation instructions explicitly prohibited
> modifying that file. That restriction was too strong. The accepted architecture already establishes
> `Grant.operations` as the distinction between authorship and promotion authority. If the existing
> implementation cannot represent that distinction, the correct response is to identify and implement the
> necessary change within the authorized experimental scope — not silently leave ambient authorship absent
> because of an extra restriction I introduced."

> "Ambient authorship: resolve the `Grant.operations` implementation gap without conflating promotion and
> authorship."

**What this does to the earlier report.** Its §5 ("Blocked item") is void. Its §10 confirmation 1 ("`authority.py`
was not modified") was true of that run and is false of the repository now. Its rule/requirement totals rise.
Nothing else in it is retracted.

---

## 2. What changed under the correction

| File | Lines | Change |
|---|---|---|
| `authority.py` | 236 | `Grant.operations: frozenset`; `Grant.revoked_operations`; `Grant.permits(operation)`; `revoked(operation=None)`; `OPERATIONS = ("promotion", "authorship")`; `DEFAULT_MODEL_OPERATIONS = {"promotion"}`; `grant(..., operations=("promotion",))`; `revoke_grant(..., operation=None)`; `active_holder(..., operation)`; **`require(domain, actor, operation)`** — no default |
| `engine.py` | 1118 | Every gate names its operation; `author_ambient` added; `possible_duplicates` version-bracketed via `_registry_read`; `_refusal` key changed to `code`; `_require_human` tests `permits("promotion")` |
| `records.py` | 210 | An authored ambient `claim` may be ownerless, and only with an authorship provenance block (REV 2 §6.1) |
| `__init__.py` | 35 | Exports `OPERATIONS`, `DEFAULT_MODEL_OPERATIONS`; version `9.2.0-gate2-rev13` |

### 2.1 The authority split, implemented

| Element | Before | Now |
|---|---|---|
| `Grant` | `grant_id, domain, holder, holder_kind, granted_by, granted_turn, revoked_turn` | + `operations`, + `revoked_operations`, + `permits(operation)` |
| Domain declaration | human with an implicit right to act | human with `frozenset({"promotion", "authorship"})` — one holder, two capabilities |
| Model grant | promotion only, by construction | defaults to `{"promotion"}`; `operations` must name `authorship` explicitly |
| Gate | `require(domain, actor)` | `require(domain, actor, operation)` — **no default**, so a promotion grant cannot silently authorise authorship |
| Revocation | whole grant | whole grant, or one capability (`revoke_grant(..., operation=...)`), leaving the other intact (L10) |

Three properties this makes structural rather than conventional: generation cannot promote (rule 24);
authorship cannot promote (rule 39 / K25); revoking one capability cannot disable the other (REV 2 §7, L10).

### 2.2 `author_ambient`

`author_ambient(*, actor, domain, target, text, turn, rationale=None, revises=None)`

| Rule from REV 2 §6.3 | Implementation |
|---|---|
| requires `authorship` capability in that domain | `self.authority.require(domain, actor, "authorship")`; a promotion-only grant is refused |
| creates a `claim` with kind `declaration`, contributor the actor | `contribution_kind="declaration"`, `contributor=actor` |
| never creates truth; never appears in a world answer | non-truth-apt class; `_eligible_truth` requires `truth_apt` |
| `revises` appends and marks the prior `superseded`, preserving text and history | prior `status=SUPERSEDED`, `superseded_by` set, `rec.supersedes` set |
| a lawful `record_cited` basis for a later promotion | `promote`'s citation loop admits it: no owner, so no ownership refusal |
| refusals: no capability · unknown domain · empty text · `revises` from another chat or a truth record | all four implemented, the last two as `OperationError` before mutation |

**The one design decision.** The authored claim has `owner_character = None`. That is not a relaxation of
rule 3 invented here: REV 2 §6.1 adopted "`claim` + `declaration` + provenance block `{"authorship": …}`"
precisely because ambient material is domain-scoped authorship and "a domain-scoped authored claim is not
private and is not owned by a character". The constructor therefore permits an ownerless non-truth record
**only** when `provenance["authorship"]` carries an actor and a domain, and only for `claim`; every other
private record still requires exactly one owner.

| Consequence | Effect |
|---|---|
| `perspective` never returns an authored claim | the pool is owner-matched, and no character owns it |
| Promotion may cite it | it is retrievable by the promotion operation, satisfying K14 |
| It never reaches a truth answer | class filter and `_eligible_truth` |
| It carries its author, domain, target, turn and rationale | rule 23 / K9 |

---

## 3. Statements in the earlier report that no longer hold

| Earlier report | Now |
|---|---|
| Status line: "with one blocked item reported (§5)" | No blocked item; §5 is void (§4 of this supplement) |
| §2 table: `authority.py` 156 lines, "**Unchanged**" | 236 lines, changed under this authorisation; sha256 in §7 |
| §2 table: `engine.py` 1045 | 1112 |
| §2 table: `records.py` 192 | 210 |
| §3 row 15: "`Grant.operations` … **NOT APPLIED — prohibited** (§5)" | **Applied** |
| §3 row 16: "`author_ambient` … **NOT APPLIED — blocked** (§5)" | **Applied** |
| §3 row 17p: "`require(domain, actor, operation)` … **NOT APPLIED — prohibited** (§5)" | **Applied**, with the operation argument now required rather than defaulted |
| §5 in full | Void; replaced by §4 of this supplement |
| §6, rule 39: **P** — "promotion present; `author_ambient` absent" | **I** |
| §7, K25: **A** | **I** |
| §7, L8: **A** | **I** |
| Totals: rules 56 I · 1 P; requirements 113 I · 1 P · 2 A · 2 NA | See §6 |
| §8 assumption 1 (refusal key `reason`) | **Resolved to `code`**, §5.3; code changed |
| §8 assumption 8 (`proposition_unregistered` invented) | **Resolved compatible**, §5.4; name retained |
| §8 assumption 10 (`possible_duplicates` unbracketed) | **Resolved**, §5.2; now version-bracketed, Class B on exhaustion |
| §10 confirmation 1 ("`authority.py` was not modified") | Superseded: it was not modified **then**; it is modified **now**, by this authorisation. The formal record for REV 13 §7 row 17p is `MEMORY-09-OWNER-CLARIFICATION-AUTHORITY-REQ17P-2026-11-17.md` |
| §9 residual rows 1 and 2 (blocked; absent) | Closed |

---

## 4. Review item 1 — the `Grant.operations` gap

**Disposition: closed by implementation, not by argument.**

The gap was never architectural. REV 2 §7 specifies exactly two changes — "one field on an existing object
plus a parameter on an existing check" — and both are now present. The briefing's prohibition, not the
contract, made `author_ambient` look like a limitation; with the prohibition lifted, rule 39's two-operation
regime is representable and represented.

What distinguishes the two capabilities mechanically:

| Path | Gate | What it can create |
|---|---|---|
| `promote` / `supersede_authority` / `revoke_promotion` / `correct` / `declare_available` / `withdraw_declaration` / `assert_incompatible` / `propositions` / `possible_duplicates` | `require(…, "promotion")` | authoritative records, coverage, assertions |
| `author_ambient` | `require(…, "authorship")` | a non-truth-apt `claim` |

No call site passes a computed or defaulted operation; each names its capability literally, so a reader can
check the split by reading the call sites (`grep -rn "authority.require(" eval/mem1/m09/*.py` → eight sites,
seven `"promotion"`, one `"authorship"`).

---

## 5. Review items 2–5

### 5.1 (item 5) The conformance report checked against operative text

Re-checked row by row against the quoted sources in
`MEMORY-09-GATE2-VERBATIM-SOURCES-2026-11-17.md`. Three classifications changed (39, K25, L8 — all upward,
all caused by the lifted restriction). The re-check also produced these findings, each now reflected in the
code or in the report:

| Finding | Action |
|---|---|
| Rule 24 was **I** only because promotion was the sole capability; a per-operation check is what rule 24's "separate" means in practice | still **I**, now for the right reason |
| The earlier §6/§7 totals line was internally garbled before delivery ("53 I · 2 P … precisely 55 I") | corrected in that file to 56 I · 1 P before delivery; superseded by §6 below |
| `_proposition_records` / `_proposition_state` exist only as blocks inside `_assemble`, not as named functions | recorded as assumption 9; behaviour identical, one fewer place to diverge |
| `proposition_kind_mismatch` is unreachable: the kind has exactly one carrier (the registry), so a record cannot disagree with its identity | recorded; the check REV 7 §3.2 asks for is satisfied by construction, and the code is not emitted because there is nothing to emit it for |
| `integrity` in canon views carries string codes; in the authority view it carries `{"code", "group_id"}` | matches REV 7 §4 / REV 9 §3.3 as quoted |

No other row's classification changed. Rows that depend on runtime behaviour (F7, Class B reachability,
concurrency) remain **NA** and are not upgraded by this re-check.

### 5.2 (item 2) A6′ and `possible_duplicates`

**Disposition: resolved in the code; the earlier "merely stale-able" claim is withdrawn as insufficient.**

Two separate questions were conflated in the earlier report. Separated:

| Question | Answer |
|---|---|
| Is A6′ compliant? | A6′ is **not implemented**: no catalogue exists, so no pre-gate lookup is possible. REV 13 §3.2 states that dropping the catalogue entirely "must not change contract semantics", and the interpretation constraint (`never independently determine existence`, `never introduce/select/suppress a refusal`, `return an empty result for an unknown label`, `never change states/relations/collections/identities`) is therefore satisfied by absence. |
| Is `possible_duplicates` compliant? | It is a **non-refused response**, so S-1, S-2 and S-7 apply to it. It now takes one version-bracketed read (`_registry_read`), projects every field from that one list, and returns `refused / inconsistent_snapshot` on exhaustion instead of reporting a stale or blended result. Reading the registry only means no *cross-store* blend was ever possible; that was not enough, because a report assembled across two registry states is itself a blend. |

The stronger statement now available: every non-refused response in the module — the three canon/authority
views **and** the multiplicity audit — is a function of exactly one coherent read, or is a Class B refusal.

### 5.3 (item 3) The refusal envelope: `reason` or `code`

**Disposition: `code`. The earlier report's choice of `reason` was wrong, and the code is changed.**

The argument, on operative text only:

| Instrument | What it says | Weight |
|---|---|---|
| REV 7 §5.1 | `{"status": "refused", "reason": "no_domain_authority"}` | the only place the envelope is printed with `reason` |
| REV 12 §2.2, §7 row 17m | column "Codes"; "Two shapes sharing one envelope … both carry only `code`" | the revision that **introduces the class split**, and therefore restates the envelope |
| REV 12 §5 example | `{"status": "refused", "code": "inconsistent_snapshot"}` | illustrative |
| REV 13 §7 row 17m′ | "Both classes carry `code` only; Class A additionally may consult metadata before constructing the code" | the operative revision, restating REV 12's row |
| REV 12/13 S-9, S-10, S-11 | "no data-derived field", "A Class B code", "the refusal class and code" | accepted contract language in the operative snapshot text |

`reason` appears once, before the class split existed. `code` is stated by both revisions that define the
refusal classes, including the operative one, and the accepted S-rules consistently speak of codes. A
conflict between an earlier general statement and two later specific restatements resolves to the later
specific text; that is also the only reading under which REV 12 §7 row 17m is implementable at all.

Implementation: `_refusal(code)` returns `{"status": "refused", "code": code}`. Both classes use it, so the
"identical fixed refusal shape" of S-9 holds, and the value set is unchanged
(`no_domain_authority`, `unknown_view`, `inconsistent_snapshot`, `target_required`,
`not_a_perspective_surface`).

### 5.4 (item 4) `proposition_unregistered`

**Disposition: compatible with the contract's orphan requirement; retained, with one caveat.**

| Requirement | Text | Consequence |
|---|---|---|
| REV 5 §4, item 3 | "The union of all entries' `records` equals `_eligible_truth(domain, target)` … A violation of (3) … is likewise an integrity failure, because `proposition_id` is required … The response must not quietly omit it, or enumeration would under-report canon while looking complete." | An orphan **must** be reported. |
| REV 5 §4 | "A violation of (1), (2) or (4) is an **integrity failure reported as such**, not silently filtered." | The report is the `integrity` channel. |
| REV 7 §4 | "`integrity` — Codes only in `world` / `belief_facts`: `"group_spans_query_kinds"`, `"proposition_kind_mismatch"`. **No group id, no foreign record id, no text, no count**" | The binding constraint is on **payload**, in the context of group faults; the list is not an exhaustive vocabulary for every integrity condition, and REV 5 §4 already requires a fault REV 7 §4 does not name. |

So the contract demands a report and supplies no name for this one. Emitting
`{"code": "proposition_unregistered"}` — no id, no text, no count — satisfies the payload rule and the
reporting rule simultaneously, and the canon view emits the bare string, exactly as it does for the two named
codes.

**Caveat, stated rather than buried:** the *name* is implementation-supplied. If the owner wants the code
vocabulary closed by contract, that is a one-line addition to a future revision, not an implementation change.
Until then the code is emitted only when the registry and record store have actually diverged — a state no
operation in this module can produce, since `promote` mints before writing and both land inside one mutation.

---

## 6. Corrected totals

| Set | Result |
|---|---|
| Rules | **57 I · 0 P · 0 A · 0 C · 0 NA** |
| Requirements | **115 I · 1 P (K10) · 0 A · 0 C · 2 NA (B3, F7)** |

K10 remains **P** by design, not by defect: tier-1 O-10b outranks "latest governs", so displacement exists
only through the human-gated supersession chain (RAT-1). B3 is **NA** because the contract leaves event
versus definition unresolved and no guard was invented. F7 is **NA** because no replay facility exists and
nothing was executed.

**These remain source-inspection classifications.** 57/57 "I" means "the code implementing this rule is
present and was read", not "the rule was observed to hold". No requirement is promoted on the strength of
code existing.

---

## 7. Verification at the time of this supplement

| Item | Value |
|---|---|
| `authority.py` | changed under this authorisation — hash in the final message |
| Ratified sources | v9 `698254ca…`, map `b4ba33cc…`, ledger `da3e7239…`, acceptance `9a25ded0…`, r32 `e88ef2d1…`, plan v3 `6c4078bd…`, Amendment 1 `6afd2a0d…` — unmoved |
| Rules / map rows | 57 / 118 |
| Transfer, share, carryover | no symbol exists (`def transfer` → no match) |
| `author_ambient` | present, gated on `"authorship"` |
| Executed | **nothing** — no import, no `py_compile`, no lint, no test, no probe |

### Confirmations required by the original authorisation, restated

1. **`authority.py` is now modified**, under the owner's correction instruction and no other. The change is
   the one REV 2 §7 specifies: one field, one capability set, one parameter, one per-capability revocation.
2. **No transfer, destination-write, sharing or cross-chat carryover** was introduced.
3. **No test, probe, suite, compilation or linting was executed.** Verification is source inspection.
4. **Source claims and executed verification remain separated**: every classification here is a claim about
   code that was read, and the three runtime-dependent rows say **NA** rather than pass.

---

## 8. Remaining residuals, unchanged

| # | Item | Status |
|---|---|---|
| 1 | B3 — event versus definition | **NA**, unresolved by contract; no guard invented |
| 2 | F7 — replay | **NA**; ordering determinism only |
| 3 | Concurrency | **Not demonstrated.** The version brackets in `_capture` and `_registry_read` are sound only because `_mutation` publishes atomically in one step; no concurrent scenario was run and Class B is unreachable in this implementation |
| 4 | `assert_incompatible` reads the two cited records before its gate, because its domain is a property of those records | Structural consequence of a signature that carries no domain; limited to the two records the caller named |
| 5 | Belief-voiced dialogue (R-C), definitions (R-B), replay (R-D) | Outside these interfaces |

**Sequence:** implementation and this supplement → owner review of the verbatim sources and the corrected
report → exact-rule source review → separate testing authorisation. This supplement requests neither the
outcome of that review nor the testing gate.
