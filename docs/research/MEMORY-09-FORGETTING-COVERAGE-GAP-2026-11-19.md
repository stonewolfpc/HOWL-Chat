# MEMORY-09 — Forgetting Coverage Gap (NEW-1)

**Status:** coverage-gap record. Read-only. No requirement added, no conformance-map change, no registry
edit, no implementation, no test, no probe, no model call, no corpus, no Git operation.

| Field | Value |
|---|---|
| Date | **2026-11-19** (local, CST) |
| Decision recorded | **NEW-1** — forgetting is an established obligation with **no dedicated conformance row**; the 118 stay closed |
| Authority | Owner work order, 2026-11-19; DN-4 = YES remains authoritative |
| Set size | **118**, unchanged |
| Files changed | This record only (plus two companion records) |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 1. The obligation being recorded

**CG-01 (Source-derived).** DN-4 = YES established six things: forgetting is **distinct from
retraction**; it is a **new event**; the **account and its history remain preserved**; the character
**no longer retains the knowledge as usable knowledge**; the **historical fact of acquisition remains
represented**; and it **must not** be implemented as an ordinary content-status change or conflated with
retraction.

**CG-02 (Source-derived).** These are T2 obligations — K-R3, K-R19, K-R17, K-R20, K-R25 plus K-R1/K-R2
— stated against requirements that exist and are not in doubt. DN-4 settled the operation's
**existence**; T2 had already settled its **representation**.

---

## 2. The gap

**CG-03 (Calculated).** **0 of the 118** ratified requirements test character knowledge-loss. Verified
by search of `MEMORY-09-CONFORMANCE-MAP-2026-11-17.md`: `forget` returns no row; every `revoke` match is
a different subject — authority (B7c, K26), belief (D4, D5), and C1, where a *teller's* revocation
leaves the recipient knowing.

**CG-04 (Observed).** The closest available rows bound the question from the outside:

| Row | What it says | Why it is not coverage |
|---|---|---|
| **A2a** | "Character still holds what it acquired" after a *memory* is retracted | governs source retraction, not forgetting |
| **C1** | a teller's later revocation leaves Sable knowing | the same distinction at prose level; not a test of forgetting |
| **E1b** | accounts unchanged by retraction | retraction path only |
| **F1 / F1b** | reconstruction from events; mismatch is an integrity failure | **would** detect a missing forget event, but only if one is written and the compare is run against a projection that models it |

**CG-05 (Inferred).** F1/F1b are the nearest thing to indirect coverage, and they are **not**
equivalent: they compare a rebuild to a stored projection. If neither models forgetting, both rebuild
identically and the check passes while the behaviour is absent. Indirect detection is not coverage.

---

## 3. Why the set remains closed

**CG-06 (Source-derived).** O-6b closed the requirement set at **118**. Adding a row would require an
owner decision the work order explicitly withholds ("Do not add a requirement row").

**CG-07 (Inferred).** Closing the set does not narrow the obligation. The obligation comes from T2 and
from DN-4, not from the 118. The 118 are the experiment's *requirements*; T2 is the architecture the
experiment is supposed to test. A gap in the test set is a gap in **coverage**, not a licence.

**CG-08 (Source-derived, owner's wording).** The absence of a dedicated row is **not** permission to
omit the behaviour. Recorded as the governing sentence of this gap.

---

## 4. How the gap is represented without modifying the set

**CG-09 (Observed).** The gap is carried **in documentation only**: this record, the DN-1/DN-4 decision
record, and the dependency summary. The conformance map, the registry, v9, r33 and the ledger are
unchanged.

**CG-10 (Inferred).** A reader of the conformance map alone cannot discover this gap. That is the cost
of not editing the map, and it should be stated wherever the map is relied on: **"no row" means "not
tested by the 118", never "not required".**

### 4.1 Relationship to the D-4 precedent — same form, different species

The owner asked to use K-R46–K-R48 as a reference without assuming it dictates the details. It does not.

| | **D-4 (K-R46–K-R48)** | **NEW-1 (forgetting)** |
|---|---|---|
| Why no row | the guarantees are **outside MEMORY-09's authority** (O-13: application-level) | the obligation is **inside** MEMORY-09's authority (T2, DN-4 = YES) |
| Ledger wording | "not enforceable by MEMORY-09 under the ratified experiment scope" | **not applicable** — it *is* enforceable by MEMORY-09 |
| Species | **out-of-scope guarantee** | **in-scope obligation with no ratified test** |
| Consequence if ignored | the application may fail a guarantee MEMORY-09 never claimed | MEMORY-09 **violates** a requirement it holds |
| Future route | separate application-level requirements (O-6b closed) | a new MEMORY-09 test, which also needs new authorisation |

**CG-11 (Inferred).** Borrowing D-4's *form* — record the gap, keep the set closed, name the consequence
— is sound. Borrowing its *reasoning* would be wrong: forgetting is not an outside-scope guarantee, and
a record that says so would misstate the obligation.

---

## 5. Consequences for conformance assessment

**CG-12 (Inferred).** Four states must now be distinguished where three were used before:

| State | Meaning | Forgetting today |
|---|---|---|
| **met** | a row exists and passes | — |
| **unmet** | a row exists and fails | — |
| **untestable** | out of scope, no row possible (D-4 species) | — |
| **untested** | in scope, obligation live, **no row exists** | **this is forgetting** |

**CG-13 (Inferred).** A conformance verdict of "118/118 met" would be **true and misleading**: it says
nothing about forgetting. Any future summary must state that the 118 do not cover it, or the verdict
overstates what was established.

**CG-14 (Inferred).** The gap also affects *negative* claims. "The system never removes knowledge
without an event" cannot be asserted from the 118; C1 and A2a bound only the retraction path.

---

## 6. Consequences for future test coverage

**CG-15 (Source-derived).** T2 already specifies three probes that would test the obligation, and none
is among the 118:

| Probe | What it would establish | Status |
|---|---|---|
| **S37** | multi-cycle learn → revoke → relearn; `acquired_turn` re-dated; "did A ever know" still yes; no `revoked_turn` field | never executed |
| **S38** | a revoke with no open episode rejected at write; same-turn resolution by `seq`, never insertion order | never executed |
| **S41** | forgetting (no edge) and recall difficulty (edge, reduced surfacing) remain distinguishable | never executed |

**CG-16 (Inferred).** S37 and S41 are directly runnable as behaviour tests once an implementation models
forgetting. **S38 is not** — it asserts resolution "by `seq`", and `seq` does not exist in the schema
(PH1A-05, companion record). So the coverage gap is not uniform: two probes are blocked on
implementation, one is blocked on the schema question as well.

**CG-17 (Inferred).** Any forgetting test added to the experiment is **new authorisation** under O-6b.
This record does not request it.

---

## 7. What this record does not do

**CG-18 (Observed).** It adds no requirement, changes no row, edits no map, registry or ledger, and
makes no implementation proposal. It does not decide PH1A-05, DN-6a/b, or any schema question. It does
not claim the 118 are deficient as a *set* — only that forgetting is not in it.

---

## 8. Verification

| Check | Result |
|---|---|
| Conformance map hash | `b4ba33ccbfe8f8c1` — unchanged |
| Requirement count | **118** — unchanged |
| Tests / probes run | **0** |

**End of gap record.**
