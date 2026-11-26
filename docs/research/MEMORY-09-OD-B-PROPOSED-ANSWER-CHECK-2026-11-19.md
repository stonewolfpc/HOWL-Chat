# MEMORY-09 — OD-b: Proposed-Answer Check (`provisional`)

**Status:** verification record. Narrow scope as instructed. **No canon, schema, code, test or registry
change.** No decision taken; OD-b remains open.

| Field | Value |
|---|---|
| Date | **2026-11-19** (local, CST) |
| Instruction | Inspect every proposed-answer section containing `provisional`; determine whether the term has a distinct meaning there; report exact locations and wording; reassess A/B/C; change the brief only if findings justify a documented correction |
| Outcome | **One finding is decisive and contradicts Reading C as stated; one prior finding of mine (A2-43) is corrected** |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 1. What was searched, exactly

**ODBC-01 (Observed — files and result).**

| Target | `provisional` hits | Substantive? |
|---|---|---|
| r33 §9 "proposed answers to the r3 questions" (Q1–Q7), `:241–252` | 1 (`:241`, meta) | **No** |
| r33 `:283` Assembly row | 1 | **Yes** — an assembly prohibition |
| r33 §11 lifecycle, `:600–612` | 1 | **Yes — decisive** |
| r33 §10 conflict semantics, `:459` | 1 | **Yes** — materiality |
| r33 `:72` D4, `:446` C4, `:941` S12, `:956` S26, `:1014` | 5 | 4 substantive, 1 meta (`:1014`) |
| `MEMORY-09-ARCH-03-R3-INDEPENDENT-REVIEW-RECORD-2026-11-17.md` | 1 (`:166`, F14) | No — provenance (`via`/`attribution`), not status |
| `MEMORY-09-REVIEW-SUPPLEMENT-EVENTS-AND-TIME-2026-11-17.md` | 0 | — |
| `MEMORY-09-DOCUMENTATION-REVIEW-2026-11-15.md` (r4 review) | 0 | — |
| `MEMORY-09-REVIEW-LEDGER-2026-11-15-REV-2.md` | 0 | — |
| `MEMORY-09-CONFORMANCE-SOURCE-REVIEW-2026-11-17.md` | 0 | — |
| `MEMORY-09-POST-Q-B-DEPENDENCY-AND-BLOCKER-REVIEW-2026-11-19.md` | 0 | — |

**ODBC-02 (Observed — the r3 proposed answers do not redefine the term).** The section at `:241` opens
"Your provisional answers to the r3 questions, recorded as proposals. None is a decision" — **meta-sense
only**. Q1–Q7 (`:243–252`) never use `provisional` as a status value; Q6 uses `retracted/corrected`
(`:251`). **The owner's concern that the proposed answers might carry a distinct meaning is not borne out.**

---

## 2. The decisive finding

**ODBC-03 (Source-derived — verbatim, r33 `:600–612`, §11 "Lifecycle and retraction").**

```
evidence turn
   └─> provisional (status, class, attribution set)
          ├─> active        (episodic, preference, validated fact, belief)
          ├─> superseded    (replaced; history kept)
          ├─> corrected     (earlier statement wrong; correction chain)
          └─> retracted     (asserted, now known false; acquisitions retained)
```

followed by: "Promotion per kind: episodic — on witnessing; preference — immediate for holder; belief — never
promotable to fact; **rumour — stays provisional**, repetition is not corroboration; world fact — authority or
corroboration; absence — attributed. Rejected: blanket three-source rule, promotion by retrieval, decay
altering truth, deleting history."

**ODBC-04 (Inferred — what this establishes).** `provisional` is **the lifecycle's entry state**, reached from
"evidence turn" at formation, and `active` is its **first child**. A third enumeration of the status
vocabulary therefore exists, and it includes both values K-R12 omits.

**ODBC-05 (Observed — `provisional` has a producer).** The producer is **formation itself** (`evidence turn →
provisional`), with `class` and `attribution set` written at the same step. **A2-43's statement that
`provisional` "has no producer anywhere" is superseded** — see §5.

**ODBC-06 (Inferred — a state with no exit).** "rumour — stays provisional" makes `provisional` **terminal for
one kind**: a state that may never be left. A lifecycle structure must be able to represent a state with no
outgoing transition; this is a structural requirement, not a vocabulary nit.

**ODBC-07 (Source-derived — provisional must not be presented as verified).** r33 `:283`, Assembly row,
prohibited behaviours include "**present provisional as verified**". This is a *presentation* constraint
attached to the value — evidence that the value carries a verification meaning at assembly, independent of its
lifecycle position.

---

## 3. Reassessment of A, B and C

**ODBC-08 (Inferred — Reading C is not viable as stated).** C claimed `provisional` is a verification
property, **not** a lifecycle state, and that K-R12's omission is therefore a category boundary. `:604`
**makes it the lifecycle's entry state and the parent of `active`**. C's premise fails.

**ODBC-09 (Inferred — C's residual insight survives, as C′).** What remains true: the value carries a
verification dimension that assembly reads (`:283`, C4 `:446`), and that conflict materiality reads
(`:459` "provisional vs established"). So the corrected form is:

> **C′** — Reading A holds, **plus**: a record also carries a verification/qualification attribute that
> assembly and conflict semantics consume, and which `status` alone may not be the right home for.

C′ is **not** a competing reading. It is a dimension observation on A.

**ODBC-10 (Inferred — Reading A is strengthened, and refined).** A is now the only reading consistent with
§5 `:322`, §11 `:604`, D4 `:72`, C4 `:446` and `:612`. Refinement: `provisional` is not merely "a state" — it
is the **entry state**, entered by **formation** and left by **evidence-kind-driven promotion** (witnessing /
immediacy / authority or corroboration), **never** by retrieval, repetition or approval.

**ODBC-11 (Inferred — Reading B remains unsupported).** §11's exit conditions are **evidence kinds**, and the
rejected list explicitly bars "blanket three-source rule" and "promotion by retrieval". No reviewer, approver
or approval appears. B would still have to invent an actor.

**ODBC-12 (Calculated — revised viability).**

| Reading | Before | After | Basis |
|---|---|---|---|
| **A** — truth/lifecycle state | strong | **viable, strengthened and refined** | §5 + §11 + D4 + C4 |
| **B** — review state | weak | **not supported** | §11 exit conditions are evidence-based; no reviewer exists |
| **C** — independent qualification | analytic | **not viable as stated**; survives as **C′** | §11 `:604` contradicted its premise |
| **C′** — A + verification dimension | — | **viable** as an addition to A, not an alternative | `:283`, `:446`, `:459` |

---

## 4. The owner's additional question — answered

**ODBC-13 (Source-derived — question).** "Does Reading C actually explain the term's relationship to `active`,
or does it merely explain why `provisional` is absent from K-R12?"

**ODBC-14 (Inferred — neither, and the distinction is now moot).** C explained **only** the K-R12 omission,
and even that is better explained otherwise (ODBC-15). C said nothing about `active` — and could not have:
under C, `provisional` was not a lifecycle state, so it could not be `active`'s predecessor. §11 `:604` makes
it exactly that.

**ODBC-15 (Inferred — the better explanation of the omission).** K-R12 `:111` is an **enumeration of the
transitions it cares about** (`established | superseded | corrected | retracted`), not a complete state list.
It omits the entry state. That is **incompleteness**, which is a P-1/P-2 finding about the enumeration — not a
category boundary and not evidence of a third dimension.

---

## 5. Correction to a prior finding

**ODBC-16 (Superseded).** **A2-43(i)**, "`provisional` has no producer anywhere", is **wrong** and is
corrected here: formation produces it at `:604`, together with `class` and `attribution set`.

**ODBC-17 (Inferred — what the real gap is).** The gap is **not** a missing producer but a **missing typed
representation**: §11 describes the transitions in a diagram with no event vocabulary, and K-R12's stream
cannot record the entry state at all. Under **OD-a Option B**, that reproduction is precisely what the
lifecycle structure must carry. The finding is smaller and more actionable than the one it replaces.

**ODBC-18 (Observed — consequence for records already filed).** A2-43 is cited in
`MEMORY-09-A2-ORDERING-AND-EVENT-SCHEMA-RECONCILIATION-2026-11-19.md` `:297` and `:355`, and the OD-a record
inherits the wording. **The A2 report needs the same one-line correction** so the record does not carry a
finding this check has refuted.

---

## 6. A term conflated, newly identified

**ODBC-19 (Observed — "promotion" in two senses).**

| Sense | Meaning | Source |
|---|---|---|
| §11 record promotion | a record leaves `provisional` by an evidence kind | r33 `:610–612` |
| v9 authority promotion | an actor confers world/scenario authority, under a grant | v9 rules 24, 27–38 (`:111–126`) |

**ODBC-20 (Inferred).** These are different operations with different authorities and different effects. The
same word is used for both, and §11's "world fact — authority or corroboration" sits at the seam. This is a
**P-2 inadequate-definition** finding, not a new requirement. It matters to OD-b because Reading B would be
tempted to read §11's promotion as an authority act.

---

## 7. What this changes in the OD-b brief — and what it does not

**ODBC-21 (Observed — renumbered after brief revision 2).** Corrected in the brief, each marked with a pointer
to this record: **ODB-04** (third enumeration), **ODB-06** (producer exists — was ODB-08 in revision 1),
**ODB-09** (Reading B retired), **ODB-10** (Reading C retired; supersedes revision 1's ODB-17 and ODB-22, whose
numbers no longer exist), **ODB-05** (enumeration incompleteness), **ODB-35** (r3/r4 check closed). Revision 2
also reframes the decision itself (§5, L1/L2/L3).

**ODBC-22 (Observed — unchanged by this check).** `active` vs `established` still has **no** mapping, and the
question is now sharper: §11 names `active` as provisional's child and never uses `established` as a state;
§5 lists `active` and never lists `established`; K-R12 lists `established` and never lists `active`.
`forgotten` remains in **none of the three** enumerations (v1's `Event` carries it; DN-4 requires it). `evidence_span`,
OD-g, domain identity, and event cause/evidence are untouched.

**ODBC-23 (Observed).** OD-b remains **open and undecided**. This record reports what canon says; it does not
select a reading, and it authorises nothing. No canonical file, schema, code, test or registry was changed.

---

## 8. Limits

**ODBC-24 (Observed).** Search was term-based over `docs/research` markdown plus the previously read Go types
file. The r3/r4 answer sections were located via r33's in-document reference (`:241`); I did not exhaustively
enumerate every review round's attachment set, and §11's authority level was read as the same "design record,
not implementation" status declared at r33 `:3`.

**ODBC-25 (Unestablished).** Whether §11's diagram predates or postdates §5's list, and which the owner regards
as controlling if they conflict; whether the r4 review discussed `provisional` in an attachment not present in
this repository; whether any unratified requirement depends on provisional-as-review.

**End of check.**
