# MEMORY-09 — S-2 P-1 / P-2 Comparison

## Versioning defined, and the three issues that decide between dedicated and generic change entries

**Status:** S-2 logical comparison. **No carrier is adopted. No schema, code, migration, test, requirement, map change
or canonical capture.** The recommendation is a decision input, not an adoption.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Question | before adopting P-1, compare it with P-2 on reconstruction reliability, integrity detection/recovery, and close/concurrent changes |
| Versioning | **per memory, per dimension**, contiguous, monotonically assigned; baseline is the formation value |
| Finding | **P-1 is stronger on all three issues**; P-2 is viable only with per-dimension versions and strict dimension/value validation |
| Recommendation | **P-1 remains the leading candidate, not adopted**; the owner's next decision settles the logical representation |
| Not decided | physical schema, implementation, requirement capture, canonical capture, S-1 fold, N-c |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. Task and boundary

**S2C-01 (Observed — the owner's request).** Define the local versioning precisely — assignment, concurrent changes,
what a missing version means — then compare P-1 and P-2 on: (i) reliability of reconstruction; (ii) integrity
detection and recovery from incomplete history; (iii) complexity when multiple changes occur close together or
concurrently. Do not settle physical schema, implementation or requirement capture.

**S2C-02 (Observed — carried state).** Under B, P-1 (dedicated `scope`/`class` change entries) and P-2 (one generic
change entry with a dimension field) are viable; P-4 is excluded by K-R14; P-6 cannot represent pure changes; P-3 is
a decline; P-5 is storage (S2I-04…S2I-15). P-1 was recommended, not adopted (S2I-28…30).

**S2C-03 (Observed — the boundary).** This record compares logical representations. It does not create a kind, a
field, a key, a schema, code, a test, a requirement or a canonical change. The map stays at **118**.

---

## 2. Versioning, defined

### 2.1 The baseline

**S2C-04 (Calculated — baseline).** The **formation value** is the baseline for each dimension. It has no change
entry; it is the value written at formation (r33 `:279`, `:404`, `:604`). If no change entry exists for a dimension,
the projection uses the baseline.

### 2.2 Assignment

**S2C-05 (Calculated — the assignment rule).**

- Each memory has **two independent sequences**: `scope_version` and `class_version`.
- The first accepted change to a dimension is **1**; each subsequent accepted change is **previous maximum + 1**.
- Versions are **positive integers**, **contiguous**, and **unique within (memory, dimension)**.
- A change entry carries its dimension's version; it does not carry the other dimension's version.
- The baseline is not versioned; it is the absence of a change entry.

**S2C-06 (Calculated — who assigns).** The **write path** assigns the next version when it accepts the change. The
representation states the invariant; it does not prescribe the physical mechanism (counter, compare-and-swap, or a
serialized writer). Assignment must be atomic with respect to the `(memory, dimension)` pair: two accepted entries
cannot share a version.

### 2.3 Concurrency

**S2C-07 (Calculated — same dimension).** Two concurrent changes to the same `(memory, dimension)` must be
**serialized** or one must be **rejected/retried**. The first accepted write takes the next version; the second
re-reads the maximum and takes the following version, or is surfaced as a conflict. A write is never silently dropped.

**S2C-08 (Calculated — different dimensions).** `scope` and `class` changes use **independent sequences**, so a
concurrent scope change and class change do not contend. This is a property of the logical model, not a physical
implementation.

**S2C-09 (Calculated — close together).** Two changes to different dimensions in the same turn each take version 1
in their own sequence. Two changes to the same dimension in the same turn take versions 1 and 2 in acceptance order;
the order is the version order, not the log append order.

### 2.4 Gaps, duplicates and recovery

**S2C-10 (Calculated — what a gap means).** A **missing integer** between the baseline and the highest version for a
dimension is an **integrity error**: an accepted change entry is absent. A **duplicate** version is a conflict and an
integrity error. An entry whose version is out of order relative to the append-only log is an integrity error.

**S2C-11 (Calculated — recovery).** On a gap or duplicate, the projection is **not trusted** for that dimension. The
system surfaces the error to `inspection` (§6.4); it may reconstruct up to the last contiguous version, but it must
not present a "current" value as authoritative for the incomplete dimension. Recovery requires the missing entry or a
re-record; silent repair is forbidden.

**S2C-12 (Calculated — the completeness limit).** A version gap detects a missing **internal** entry only because
assignment is contiguous. It **cannot** detect a change that was never recorded at all: if no later version exists,
there is no gap to observe. The check verifies the projection against the **recorded** stream; it is not a
completeness detector for the write path. This limit is shared by P-1 and P-2.

**S2C-13 (Calculated — illustrative sequence).**

| Event | Scope value | Scope version | Class value | Class version |
|---|---|---|---|---|
| Formation | `character_private` | baseline | `experiential` | baseline |
| Character-only → scenario-wide | `scenario` | 1 | — | — |
| Scenario-wide → world-wide | `world` | 2 | — | — |
| World-wide → scenario-wide | `scenario` | 3 | — | — |
| Class reclassification (illustrative; cause not implied) | — | — | `character_specific` | 1 |

The two sequences are independent. The class example is a logical illustration; it does not assert that the cause is
permitted (K-R9/K-R10/D5 still govern causes).

---

## 3. Issue 1 — reliability of reconstruction

**S2C-14 (Calculated — P-1).** Reconstruction filters the lifecycle stream by the two dedicated kinds, takes the
highest version per dimension, and falls back to the formation value. The kind itself identifies the dimension, so
there is no polymorphic `dimension` field to validate or misinterpret. The value domain is determined by the kind.
Reliability is **high**, subject to S2C-12's completeness limit.

**S2C-15 (Calculated — P-2, per-dimension versions).** Reconstruction filters by the generic kind and reads the
`dimension` field, then takes the highest version for that dimension. Reliability is **adequate** if the
`dimension`/`value` pairing is strictly validated: the dimension field selects the allowed value domain. A malformed
or missing dimension invalidates the entry.

**S2C-16 (Calculated — P-2, single global version).** If P-2 uses one version sequence for both dimensions, a gap is
ambiguous: the missing version could have belonged to either dimension. Reconstruction can still take the highest
per-dimension entry, but it cannot tell which dimension is incomplete from the gap alone. Reliability is **lower**.

**S2C-17 (Calculated — comparison).** P-1 has a **slight reliability advantage** from explicit kinds and narrower
validation. P-2 is adequate with per-dimension versions and strict validation; the global variant is weaker. Neither
solves S2C-12.

---

## 4. Issue 2 — integrity detection and recovery

**S2C-18 (Calculated — P-1).** Each dimension has its own contiguous sequence, so a gap identifies **which**
dimension is incomplete. A duplicate is detectable within the dimension. Recovery can name the missing version and
the affected dimension; the projection is untrusted for that dimension only. The other dimension remains usable.
This is a **strong** integrity story.

**S2C-19 (Calculated — P-2, per-dimension versions).** The same detection is possible **if** the `dimension` field is
reliable and the version is scoped by dimension. A gap identifies the dimension via the field of the missing entry —
but the missing entry is absent, so the dimension must be inferred from the sequence context or a separate index. If
the version is stored per dimension, the gap's dimension is known; if the version is a single field with a dimension
tag, the missing tag is unknown. This is a **slightly weaker** recovery story than P-1.

**S2C-20 (Calculated — P-2, single global version).** A gap is **ambiguous** between dimensions. Recovery cannot tell
which dimension lost an entry without additional metadata. The check may have to treat both dimensions as untrusted,
which over-constrains recovery. This is **weaker** than P-1.

**S2C-21 (Calculated — comparison).** P-1 is **stronger** because the kind carries the dimension even when an entry is
missing from the sequence. P-2 is equivalent only with per-dimension sequences and a reliable dimension tag; the
global variant is weaker. The difference matters most in the failure case the owner asked about: incomplete history.

---

## 5. Issue 3 — close or concurrent changes

**S2C-22 (Calculated — P-1).** Scope and class changes have independent sequences, so concurrent changes to different
dimensions do not contend. Same-dimension changes serialize or retry. The entry kind makes each change self-describing;
the projection applies each dimension independently. Complexity is **low** at the logical level; the cost is two
kinds.

**S2C-23 (Calculated — P-2, per-dimension versions).** The independent-sequence behavior is the same, but every entry
carries a polymorphic `dimension` field and a value whose domain depends on that field. Close-together changes require
branching on the dimension during validation and projection. Complexity is **moderate**.

**S2C-24 (Calculated — P-2, single global version).** All changes contend on one per-memory counter. Concurrent scope
and class changes can conflict even though they are logically independent; retries increase. The gap ambiguity of
S2C-20 also applies. Complexity is **higher** and concurrency is worse.

**S2C-25 (Calculated — comparison).** P-1 is **stronger** on concurrency because its sequences are naturally
per-dimension. P-2 with per-dimension versions matches the concurrency behavior but adds polymorphic validation; P-2
with a global version is weaker on both counts.

---

## 6. Summary

**S2C-26 (Calculated).**

| Issue | P-1 dedicated | P-2 per-dimension | P-2 global |
|---|---|---|---|
| **Reconstruction reliability** | **high** — kind identifies dimension | adequate — strict dimension/value validation | lower — ambiguous gaps |
| **Integrity detection** | **strong** — dimension known from kind | adequate — dimension must be inferred | weak — gap ambiguous |
| **Recovery** | **strong** — untrust only the affected dimension | adequate — needs reliable dimension tag | weak — may untrust both dimensions |
| **Close/concurrent changes** | **low complexity; independent sequences** | moderate — polymorphic validation | higher — single counter contention |
| **Vocabulary** | two kinds | one kind + dimension field | one kind + dimension field |
| **Validation surface** | narrow | broad | broad |

**S2C-27 (Calculated — conclusion).** **P-1 is stronger on all three issues the owner specified.** Its advantage is
not merely stylistic: the dedicated kind preserves the dimension in the **failure** case, when the entry that would
carry a dimension tag is the one that is missing. P-2 is viable only with per-dimension sequences and strict
dimension/value validation; that configuration erodes its compactness advantage and still has a weaker recovery story
for a missing entry.

**S2C-28 (Observed — not adopted).** P-1 remains the **leading candidate, not adopted**. The comparison supports it;
it does not select it. The owner's next decision settles the logical representation (S2I-30; S2B-05).

---

## 7. What remains open

**S2C-29 (Calculated).**

| Open | Why |
|---|---|
| P-1 vs P-2 | the owner's logical decision; this record supports P-1 |
| version vs predecessor link | both avoid N-c; not settled |
| exact kind names and fields | logical detail; no physical schema |
| concurrency mechanism (serialize vs compare-and-swap) | write-path design; not S-2's schema |
| recovery procedure for a gap | operational; the representation only requires surfacing and no silent repair |
| whether to carry prior values for audit | optional; not required by the projection |
| S-1 fold scope | separate downstream decision |
| N-c | status ordering; untouched |
| requirement capture / implementation | separate gates |

---

## 8. What this does not do

**S2C-30 (Observed).**

| Untouched | State |
|---|---|
| **carrier adoption** | **0** — P-1 supported, not adopted |
| **schema / kind / field / key** | none created or implemented |
| **physical concurrency mechanism** | not selected |
| **authority / validity / causes** | fenced, not decided |
| **N-b / N-c / N-a** | not decided |
| **state machine / UX** | not designed |
| **implementation / tests** | **0** |
| **requirement capture** | **0**; map **118** |
| **canonical capture** | **0** |

---

## 9. Traceability

**S2C-31 (Observed).**

| What | Where |
|---|---|
| B authorization; baseline | S2B-01…S2B-17 |
| P-1/P-2/P-6/P-4/P-3/P-5 | S2I-04…S2I-15 |
| P-1 recommendation, not adopted | S2I-28…S2I-30 |
| Versioning rules | S2C-04…S2C-13 |
| Three-issue comparison | S2C-14…S2C-27 |
| Context inactivity ≠ lifecycle status | S2H-12 |
| K-R12 / K-R14 / §6.4 | r33 `:111`–`:113`, `:366`–`:367` |
| Map unchanged | 118 rows |

---

## 10. Record status

**S2C-32 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Carrier adopted | **0** — P-1 supported, not adopted |
| Implementation / tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **comparison** — no canonical capture |

**End of comparison.** Versioning is defined as **per memory, per dimension**, contiguous and monotonically assigned,
with the formation value as baseline; same-dimension writes serialize or retry, different dimensions are independent,
and a gap/duplicate is an integrity error with the projection untrusted for the affected dimension (S2C-04…S2C-13).
A gap detects a missing **internal** entry but **cannot** detect a change that was never recorded (S2C-12). **P-1 is
stronger on reconstruction reliability, integrity detection/recovery, and close/concurrent changes**, because the
dedicated kind preserves the dimension even when the entry carrying a dimension tag is the missing one (S2C-14…S2C-27).
**P-1 remains the leading candidate, not adopted** (S2C-28). **No schema, code, requirement, map change or canonical
capture follows.**
