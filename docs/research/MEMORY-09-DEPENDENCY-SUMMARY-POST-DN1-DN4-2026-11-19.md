# MEMORY-09 — Dependency Summary After DN-1/DN-4, NEW-1 and PH1A-05

**Status:** dependency summary. Read-only. No decision is made here. No canonical edit, no
implementation, no test, no probe, no model call, no corpus, no Git operation.

| Field | Value |
|---|---|
| Date | **2026-11-19** (local, CST) |
| Supersedes | §4 of `MEMORY-09-POST-DN1-DN4-RECONCILIATION-2026-11-19.md` (the dependency table only) |
| Companions | `…-FORGETTING-COVERAGE-GAP-…` · `…-PH1A-05-ORDERING-INVESTIGATION-…` |

---

## 1. Settled — do not re-decide

| ID | Decision | Consequence |
|---|---|---|
| **O-14** | Q-B / Branch 2 | no status-derived awareness; §8.4 classified, r33 |
| **DN-1** | v9 §8 applies to the source record | PH1A-03 divergence removed; K-R1/K-R2/K-R11 govern over the account |
| **DN-4** | forgetting = **YES** | distinct operation, distinct event, account and history preserved, no usable knowledge after; six established, eight excluded |
| **O-1 / O-3 / O-5** | ownership, retraction, belief | govern automatically |

---

## 2. Governed automatically — no owner decision needed

**DS-01 (Inferred).** These follow from tier precedence and the decisions above; applying them is not a
choice:

- **K-R12** — status is a projection (v9 silent → T2 governs).
- **K-R1 / K-R2 / K-R11** — the acquisition fold backs the account (DN-1).
- **K-R3 / K-R17 / K-R19 / K-R20 / K-R25** — the forgetting path (DN-4).
- **K-R26** — `via` precedence is a fixed total order; no ties.
- **K-R13** — read-time join, inspection-only under O-14.

---

## 3. Coverage gap — recorded, set stays closed

**DS-02 (Observed).** **NEW-1** is recorded: forgetting is an established obligation with **no dedicated
row in the 118**. The requirement set remains **118** and the conformance map is unchanged
(`b4ba33ccbfe8f8c1`). The gap is documented, not repaired. "No row" means "not tested by the 118",
never "not required".

**DS-03 (Inferred).** A `118/118 met` verdict is now known to be silent about forgetting and must be
qualified wherever it is reported.

---

## 4. Investigated, still open

**DS-04 (Observed).** **PH1A-05** was investigated, not decided. Findings: K-R15 demands a **total order
per (character, memory) pair** keyed by `(turn, seq)` with `seq` monotonic per chat; §6.2 has no `seq`;
no existing field can substitute; **the same order-stripped event multiset yields two different
conforming states**, so the fold is underdetermined; K-R17's write-time check degrades to turn
granularity; and `acquired_turn` is stated three different ways across R1c, S37 and K-R20.

**DS-05 (Inferred).** PH1A-05 remains **open**, with six options listed (O-A…O-F) and none selected.
O-A/O-B/O-F are schema changes; **O-C/O-D are normative amendments** to K-R15 and would need an owner
decision beyond PH1A-05.

---

## 5. Open owner decisions

| ID | Question | Blocking what |
|---|---|---|
| **PH1A-05** | how `(turn, seq)` is carried — or whether K-R15 is amended | any conforming subject; S38 |
| **DN-6a** | **who may cause a character to forget?** | conformance of any forgetting implementation |
| **DN-6b** | **who may retract content or another character's account?** | conformance of the retraction path (no actor exists today) |
| **Q-D** | where applicability lives on MEMORY-09's side | presentation decisions |
| **PH1A-12** | which implementation is the experimental subject, and where the bridge lives | the experiment |

**DS-06 (Source-derived).** **DN-6a and DN-6b are OPEN.** Neither is decided, and neither is implied by
DN-4. O-5 concerns **belief** ownership; it must not be extended into knowledge-removal authority by
inference.

**DS-07 (Observed).** The authority analysis is **not yet ready**. Relevant sources identified for it:
O-5 (belief, model-owned), O-1 (ownership), O-10b (authority changes are human), K-R12 (two streams), and
the absence of any actor parameter on `retract` in either implementation. Recommendations are deferred.

---

## 6. PH1A-12 status

**DS-08 (Inferred).** **PH1A-12 remains blocked.** Its actual predecessors are:

1. **PH1A-05** — the ordering question, because no subject can carry K-R15 as written;
2. **DN-6a / DN-6b** — for conformance, though a scripted experiment could run without an authority
   policy;
3. **ND-1** (§7) — if the memory stream turns out to need an ordering rule.

DN-1, DN-3 and DN-4 are discharged. No implementation or bridge was selected.

---

## 7. Newly discovered dependency (undecided)

**DS-09 (Unestablished).** **ND-1 — the memory/status stream has no ordering rule.** K-R12/K-R14 require
it to be independently rebuildable, but r32 specifies no order key and no fold rule for it. Two status
events in one turn are indistinguishable, exactly as two acquisition events are. Whether this matters
depends on whether the experiment can produce two status events in one turn — not established. Recorded
only; no requirement added, nothing decided.

---

## 8. Blockers that are not decisions

**DS-10 (Observed).** Unchanged: no corpus (P1), unbound G07/G12 (P5/P6), no grant mechanism (P7), no
live responder (P8), no bridge (P9), no inspection context (P-5), no disclosure call site (P-7), and
`ExecutionNotAuthorised` at `conversation.py:177`.

---

## 9. Verification

Canonicals unchanged: v9 `698254ca262c0874` · r33 `dfb28b6003ca356f` · ledger `e5008020ff68bd21` ·
conformance `b4ba33ccbfe8f8c1`. Tests executed: **0**.

**End of summary.**
