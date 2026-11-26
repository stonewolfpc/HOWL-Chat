# MEMORY-09 — Snapshot contract: owner acceptance

**Status: ACCEPTED by the owner, 2026-11-17. Operative snapshot text is REV 13 §2–§5. No implementation, no
tests, no amendment to any ratified source.**
Recorded **2026-11-17 CST (−0600)**.

| Field | Value |
|---|---|
| Decision | Owner accepts the snapshot contract as specified |
| Operative text | `MEMORY-09-ARCH-CONTRACTS-R2-FOUR-QUESTIONS-REV-13-2026-11-17.md` §2–§5, sha256 `5b745628…` |
| Read with | REV 11 §3.2 (S-1…S-11) and REV 12 §2–§3 as amended by REV 13; REV 10 §3 as superseded |
| Binding constraint | A6′ catalogue lookups remain ancillary (§3 below), recorded verbatim |
| Implementation | **Not authorised.** Acceptance of a contract is not authorisation to build it |
| Testing | **Not authorised**, and not implied. The exact-rule source review precedes it |
| Amends | **Nothing** to r32, v9, the map, the ledger, plan v3 or Amendment 1 |
| Ratification of REV 8 §6 remainder | **Still flagged.** This record does not settle it either |

---

## 1. Owner acceptance, recorded verbatim

> "REV 13 addresses both remaining issues. I consider the snapshot contract architecturally coherent and
> ready for your acceptance. I don't see a reason to request another revision before moving forward."

> "**One implementation constraint I want preserved.** The optional A6′ catalogue lookup must remain
> ancillary. It may support indexing or presentation, but it must never independently determine whether a
> domain or target exists, introduce a refusal, or change the semantics of accepted records and returned
> identities. That follows from REV 13's own opaque-label rule. I would treat it as an interpretation
> constraint during implementation, not a reason to expand the architecture again."

> "Also, one important distinction remains: The contract specifies what a coherent capture must guarantee.
> The current implementation has not demonstrated that guarantee. A10 retry behaviour and Class B refusal
> remain untested and unreachable in the reported single-threaded implementation. Those are implementation
> and verification matters, not defects in the proposed contract."

> "My recommendation is to accept REV 13 §2–§5 as the operative snapshot contract, with the A6′ constraint
> above preserved in its implementation interpretation. Then move to the separate implementation-authorization
> decision. That would authorize implementation work only; it would not authorize tests. The exact-rule source
> review you requested still comes before testing authorization. And we should keep the REV 8 remainder question
> flagged. REV 13 correctly does not silently settle it."

> "**accepted.**"

---

## 2. What was accepted, and what it supersedes

| Instrument | Standing after this acceptance |
|---|---|
| REV 9 §3–§9 | Ratified within its declared sections (recorded in REV 10 §6). Unchanged. |
| REV 10 §3 (snapshot model) | **Superseded.** Its S-1…S-6 survive only as REV 11 §3.2 carries them forward |
| REV 10 §6 (REV 9 ratification record) | **Stands.** Including its flagged reading of REV 8's unamended sections |
| REV 11 §3.2, §3.3, §5 | **Accepted as amended** by REV 12 §2 (refusal classes) and §3 (acquisition) |
| REV 12 §2–§4 | **Accepted as amended** by REV 13 §2 (metadata boundary) and §3 (opaque labels) |
| REV 13 §2–§5 | **Operative snapshot contract** |
| REV 13 §6–§9 | Record of status, deltas, coherence and limits; accepted as the accompanying implementation list and limits |
| REV 1–12 | Byte-identical, unamended. Their contracts stand except where superseded above |

The snapshot contract as accepted requires: one coherent capture at a single committed version; a frozen
view-scoped identity set inside it; total three-valued member resolution with no representable "not
examined"; a whole-store candidate set with relevance applied after resolution; one establishment map
projected twice; refusal classes separated by whether protected memory was read; an authorization boundary
naming what the gate may read; and opaque `domain`/`target` labels that cannot produce an existence oracle.

**State vocabulary unchanged:** `none`, `settled`, `unresolved_relation`, `conflict`, `multiple`.
`unresolved(record_absent)` is a capture-resolution value, not a memory state or a knowledge state.

---

## 3. Binding implementation constraint — A6′

Recorded as an interpretation constraint on any implementation of the accepted contract. It is not a new
architectural rule and adds nothing to the ratified set.

| A6′ may | A6′ may not |
|---|---|
| Support indexing | Determine, alone or jointly, whether a domain or target exists |
| Support presentation of already-authorized data | Introduce, select or suppress a refusal |
| Be consulted inside the capture, within the version bracket | Alter the semantics of `accepted` records or `returned_identities` |
| Return an empty result for an unknown label | Change any entry `state`, `relation`, collection or identity set |

Consequences for the implementation list: the optional `_catalogue_lookup` (REV 13 row 17o) is not on any
refusal path and not a source of any published value; if an implementation drops it entirely, no contract
statement changes.

---

## 4. What acceptance does not do

| Not done | Note |
|---|---|
| No implementation authorisation | Implementation is a separate decision, still pending. Acceptance of a contract is not permission to build it. |
| No testing authorisation | Not implied, and not granted. The exact-rule source review (R1/R2 and any further review the owner requires) precedes any testing decision. |
| No code changed | `eval/mem1/m09/*.py` byte-identical to R2, including `authority.py`. |
| No execution | 0 tests, 0 probes, 0 suites, nothing compiled. |
| No amendment | r32 `e88ef2d1…`, v9 `698254ca…`, map `b4ba33cc…`, ledger `da3e7239…`, plan v3 `6c4078bd…`, Amendment 1 `6afd2a0d…` all unmoved. |
| No tally change | Rules 39 U · 13 M · 5 N · 0 R · 0 X; requirements 86 U · 22 M · 8 N · 1 R · 1 X. |
| No closure of open items | B3, belief-voiced dialogue, replay/F7 (still **R**), R-B, R-C, R-D remain open. |
| No correction of D-A / D-B | Confirmed and **uncorrected**. They are implementation work, not contract work. |
| No settlement of the REV 8 remainder question | Whether REV 8's unamended sections are inside the ratified chain remains the owner's call, flagged in REV 10 §6 and untouched here. |

**Concurrency, stated plainly:** the accepted contract specifies a guarantee. The current implementation has
not demonstrated it. A10's retry path and the Class B refusal are unreachable in the single-threaded store
described by R1/R2, which is vacuity, not evidence. That is a verification matter for a later gate, not a
defect in the contract — and this record does not upgrade it.

---

## 5. Next decisions, in order

| # | Decision | Status | What it would authorise |
|---|---|---|---|
| 1 | Accept snapshot contract (REV 13 §2–§5) | **Accepted, this record** | Nothing but the text |
| 2 | Implementation authorisation | **Pending — owner's** | Implementation work only. Would begin with D-A and D-B, then the delta rows of REV 3 §7.1/§7.2, REV 4 §10, REV 5 §9, REV 6 §8, REV 7 §9, REV 8 §9, REV 9 §9, REV 10 §7, REV 11 §8, REV 12 §7, REV 13 §7 |
| 3 | Implementation | Blocked on #2 | Files in `eval/mem1/m09/` |
| 4 | Conformance review against the ratified rules | Blocked on #3 | A review, not execution |
| 5 | Exact-rule source review, then testing authorisation | Blocked; precedes any execution | Separate decision, separately granted |

This record requests #2 and does not presume its outcome.

---

## 6. Verification at the time of acceptance

| Item | Value |
|---|---|
| Snapshot contract | REV 13 `5b745628…` |
| Contract chain intact | `98447423…`, `74ac6f77…`, `de4f68b6…`, `1da08f2c…`, `a7e7c013…`, `41421dce…`, `a281f3e5…`, `b38f1062…`, `5b3a9008…`, `7d703b77…`, `40850419…`, `0592d689…` |
| Ratified sources | v9 `698254ca…`, map `b4ba33cc…`, ledger `da3e7239…`, r32 `e88ef2d1…`, plan v3 `6c4078bd…`, Amendment 1 `6afd2a0d…` |
| Implementation | `__init__.py` `6dcba8ac…`, `authority.py` `07dfbeb2…`, `engine.py` `205acbbe…`, `errors.py` `a96f6e50…`, `records.py` `72b7c405…`, `retrieval.py` `d33e54b4…` |
| Counts | rules 57; map rows 118 |
| Executed | nothing |

**Sequence remains:** ARCHITECTURAL CONTRACTS → OWNER REVIEW → IMPLEMENTATION AUTHORIZATION → IMPLEMENTATION
→ CONFORMANCE REVIEW → SEPARATE OWNER TESTING AUTHORIZATION. Step 2 is the next gate and is the owner's to
open.
