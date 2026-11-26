# MEMORY-09 — P-1 Completeness and Recovery

## A durable expected-version boundary is the only latest-event detector; the authority rule is logical

**Status:** design analysis. **No option is selected. No schema, code, migration, test, requirement, map change or
canonical capture.** The completeness marker is not part of P-1; adding one would be a separate architectural
decision. The recovery **authority rule** is a logical integrity constraint; the recovery **behavior** is
implementation.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Questions | completeness-marker semantics and failure modes · logical recovery rule · concurrency invariant |
| Finding | a durable expected-version boundary is the only way to detect a missing latest; every marker shifts the trust problem; only a trusted projection is authoritative |
| Concurrency | invariant fixed; mechanism left to implementation |
| Selection | **none** |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed**.

---

## 1. Task and boundary

**P1C-01 (Observed — the owner's sequence).** Investigate completeness-marker semantics and failure modes; establish
recovery behavior at the logical level; leave the concurrency mechanism to implementation while preserving the
serialization invariant.

**P1C-02 (Observed — the correction).** Recovery policy may be an implementation detail in some respects, but the rule
**never to present an untrusted projection as authoritative** is a **logical integrity constraint**. S2O-23's
classification of recovery as purely operational is **Narrowed** (P1C-27).

**P1C-03 (Observed — the boundary).** No schema, code, tests, requirements, map change or canonical capture. The map
stays **118**. The completeness marker is **not** part of P-1; P-1 established only `scope_changed` and
`class_changed`.

---

## 2. Completeness marker

### 2.1 Why it is needed

**P1C-04 (Observed — the gap).** The sequence detects an internal gap when a higher version survives. It cannot detect
a missing **latest** event: `1,2,3` and `1,2` are both contiguous as recorded, but the second may have lost `3`
(P1D-17; S2F-10/11).

**P1C-05 (Calculated — the only way to close it).** A **durable expected-version boundary** — a record of the highest
version that should exist for `(memory, dimension)` — lets a rebuild compare the stream's maximum to the expected
maximum. Without some such boundary, the missing latest is undetectable.

### 2.2 Options

**P1C-06 (Calculated — the forms).**

| Form | How it detects | Cost | Failure modes |
|---|---|---|---|
| **No marker** | internal gaps only | none | missing latest undetectable |
| **Mutable expected-max counter** | stream max < counter ⇒ missing latest | a durable mutable counter; atomic update with assignment | stale/behind, ahead, lost, corrupted, concurrent update, false gap, false complete |
| **Periodic checkpoint** | snapshot of max/hash; stream max < checkpoint ⇒ missing | storage; periodic coverage | stale checkpoint; changes after it are uncovered until the next |
| **Write-ahead reservation** | reserve version N, then write; an unmatched reservation reveals a missing change | doubles entries; a new kind/entry | reservation lost/orphaned; false gap; ordering |
| **Hash chain + checkpoint** | chain break reveals internal loss; checkpoint head reveals latest loss | a hash field + checkpoint | checkpoint stale; hash/chain integrity; no checkpoint ⇒ latest undetected |
| **External manifest** | expected versions from outside | external dependency; trust | manifest stale/lost/corrupted; must not feed the projection (K-R14) |

### 2.3 The trust problem

**P1C-07 (Calculated — a marker is not proof).** Every marker shifts the completeness question to the marker itself. A
counter can be ahead or behind; a checkpoint can be stale; a manifest can be incomplete; a reservation can be
orphaned. None proves the stream complete; each raises confidence or detects a specific failure.

**P1C-08 (Calculated — the counter's specific hazards).**

| Hazard | Effect |
|---|---|
| counter updated **before** the change is durable | a failed change leaves the counter ahead ⇒ **false gap** |
| counter updated **after** the change | a crash between them leaves the counter behind ⇒ **missed latest** |
| counter lost/corrupted | falls back to no marker or reports falsely |
| concurrent counter update | skip/duplicate; needs atomicity independent of the change |
| counter as value authority | violates K-R12/K-R14 if the projection consults it |

**P1C-09 (Calculated — the checkpoint's trade-off).** A checkpoint derived from the stream is **safer** than a mutable
counter (it cannot be ahead of the stream at the moment it is taken), but it is **periodic**: a missing latest after
the checkpoint is undetected until the next checkpoint. It is a bounded-confidence detector, not a guarantee.

**P1C-10 (Calculated — the reservation's trade-off).** A write-ahead reservation is precise for the write it guards,
but it adds an entry per change and its own completeness problem. It is a new representation element, not part of P-1.

**P1C-11 (Calculated — the boundary rule).** A marker must be an **integrity aid**, not a value authority. If the
projection consults it, K-R14's independent rebuild is at risk. A separate check comparing the stream to the marker is
safer than feeding the marker into the projection.

### 2.4 Status

**P1C-12 (Observed — no selection).** A durable expected-version boundary is the **direction** the evidence supports
for latest detection; the **form** is not selected. A mutable counter is not automatically chosen: its atomicity and
consistency hazards (P1C-08) argue for examining a checkpoint or reservation before committing.

**P1C-13 (Calculated — it is a separate architectural decision).** The marker is not part of P-1. Adding a counter,
checkpoint or reservation kind/field would extend the representation and requires a separate owner decision, not an
implementation authorization alone.

---

## 3. Logical recovery

### 3.1 States

**P1C-14 (Calculated — the logical states).**

| State | Meaning | Authoritative? |
|---|---|---|
| **recorded-complete** | the stream is contiguous; no conflict; completeness beyond the stream is unknown | **yes, relative to the recorded stream** — the missing-latest limitation is stated |
| **verified-complete** | a trusted marker confirms no missing latest | **yes, completeness-verified** |
| **incomplete/untrusted** | a gap or marker mismatch | **no** |
| **pending** | a higher version arrived before a lower one | **no** for that dimension until contiguous |
| **conflict** | a duplicate version | **no** |
| **repaired** | an explicit re-record or verified restore restored contiguity | **yes**, back to recorded-complete/verified |

**P1C-15 (Calculated — the authority rule).** Only **recorded-complete**, **verified-complete** or **repaired** states
are authoritative. An incomplete, pending or conflict state may be surfaced as **last-known** or diagnostic, clearly
flagged, but **never as the current authoritative value**. This is the logical integrity constraint the owner
identified.

**P1C-16 (Calculated — the missing-latest caveat).** "Recorded-complete" is authoritative **relative to the recorded
stream**, not a claim that no event was lost. The system must not present it as completeness-verified. The distinction
between "complete as recorded" and "verified complete" is part of the contract.

### 3.2 Failure mapping

**P1C-17 (Calculated).**

| Failure | State | May show last-known? | Action |
|---|---|---|---|
| internal gap | incomplete/untrusted | yes, flagged | surface; repair/escalate |
| missing latest, marker detects it | incomplete/untrusted | yes, flagged | surface; repair/escalate |
| missing latest, no marker | recorded-complete (unknown) | yes, with the completeness caveat | serve as recorded; do not claim verified |
| out-of-order | pending | last contiguous only, flagged | hold higher; wait |
| duplicate | conflict | no | reject/conflict; untrust |
| repaired | repaired | yes | serve; return to recorded-complete/verified |

**P1C-18 (Calculated — no silent repair).** A missing value is never invented or guessed. A repair is an explicit
re-record or a verified restore; it creates a new entry or a verified state, not a silent patch.

**P1C-19 (Calculated — no lifecycle-status change).** An incomplete projection is an **integrity state**, not a new
lifecycle `status`. Context inactivity remains separate (S2H-12); an untrusted class/scope is not `retracted`,
`superseded` or otherwise status-changed.

**P1C-20 (Calculated — no authority or validity).** Recovery does not decide who may change class/scope or whether a
change was warranted. It surfaces or repairs history; authority and validity remain outside (S2R-17).

### 3.3 Logical rule versus implementation policy

**P1C-21 (Calculated — the split).**

| Layer | Content |
|---|---|
| **logical integrity constraint** | only a trusted state is authoritative; last-known is flagged; no silent repair; no lifecycle-status change |
| **implementation policy** | fail closed vs last-known+flag; checkpoint restore vs re-record; escalation thresholds; consumer-visible behavior |

**P1C-22 (Calculated — the owner's correction).** S2O-23 said recovery is operational. **Narrowed**: the **authority
rule** is logical; the **behavior** is implementation. The split is now explicit.

---

## 4. Concurrency

**P1C-23 (Observed — the invariant).** Two simultaneous same-dimension changes must not both receive the same version
or create a skipped version. The version assignment is atomic with the entry write; same-dimension writes serialize;
cross-dimension writes are independent (P1D-07/08; S2C-05/07).

**P1C-24 (Calculated — the mechanism options).**

| Mechanism | Guarantee | Trade-off |
|---|---|---|
| single writer/queue per `(memory, dimension)` | one assigner in order | simple; bottleneck; single point |
| compare-and-swap | read max, write with expected max; conflict ⇒ retry | throughput; retries; atomic conditional write needed |
| pessimistic lock | lock the pair during assignment | predictable; contention/deadlock; timeout |
| central sequencer | one component assigns versions | strong order; central dependency |

**P1C-25 (Calculated — the invariant is logical).** The mechanism is implementation. The implementation authorization
must state the guarantee and the conflict/retry behavior; P-1 does not select a mechanism.

**P1C-26 (Calculated — marker atomicity).** If a completeness marker exists, its update must be atomic with the change,
or the marker must be derived from the stream (checkpoint). A separately updated counter has the hazards in P1C-08.

---

## 5. What remains open

**P1C-27 (Calculated).**

| Open | Why |
|---|---|
| **completeness marker form** | counter vs checkpoint vs reservation vs manifest; separate architectural decision |
| **marker trust/atomicity** | how the boundary is kept consistent with the stream |
| **recovery default** | fail closed vs last-known+flag; implementation |
| **repair procedure** | explicit re-record vs verified restore; implementation/operational |
| **escalation** | thresholds and operator path |
| **concurrency mechanism** | queue vs CAS vs lock vs sequencer; implementation |
| **formation event** | separate representation decision (P1F-13/17) |
| **N-b / `provisional`** | separate |
| **schema / code / tests / requirements / canon** | separate gates |

---

## 6. What this record does not do

**P1C-28 (Observed).**

| Untouched | State |
|---|---|
| **completeness marker** | no form selected; not part of P-1 |
| **recovery behavior** | logical rule established; implementation policy open |
| **concurrency mechanism** | not selected |
| **schema / code / tests** | **0** |
| **requirements / map** | **0**; **118** |
| **canonical capture** | **0** |
| **Howl state machine / UX** | **0** |
| **authority / validity / causes** | not decided |

---

## 7. Traceability

**P1C-29 (Observed).**

| What | Where |
|---|---|
| Missing-latest limit | P1D-17; S2F-10/11 |
| Failure taxonomy | S2F-06/07 |
| P-1 decision | P1D-01…P1D-10 |
| S-3 acceptance | S3A-27 |
| Formation baseline | P1F-09…P1F-17 |
| Context ≠ lifecycle status | S2H-12 |
| K-R12 / K-R14 / §6.4 | r33 `:111`–`:113`, `:366`–`:367` |
| Recovery classification correction | P1C-02/22; S2O-23 |
| Map unchanged | 118 rows |

---

## 8. Record status

**P1C-30 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Completeness form selected | **0** |
| Recovery behavior selected | **0** — logical rule only |
| Concurrency mechanism selected | **0** |
| Implementation / tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **design analysis** — no canonical capture |

**End of analysis.** A **durable expected-version boundary** is the only way to detect a missing latest; every form
shifts the trust problem to the marker, and a mutable counter carries the most hazards (P1C-04…P1C-11). **No form is
selected**; the marker is a separate architectural decision, not part of P-1 (P1C-12/13). The **logical recovery
rule** is established: only **recorded-complete**, **verified-complete** or **repaired** projections are authoritative;
incomplete, pending or conflict states may be shown only as flagged last-known, never as current authority
(P1C-14…P1C-19). Recovery behavior is implementation; the authority rule is logical (P1C-21/22). **The concurrency
invariant is fixed; the mechanism is implementation** (P1C-23…P1C-26). **No schema, code, requirement, map change or
canonical capture follows.**
