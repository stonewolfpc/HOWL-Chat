# MEMORY-09 — P-1 Owner Decision

## The S-2 logical representation is established; implementation is not authorized

**Status:** **owner decision.** P-1 is selected. The S-2 **logical representation** is established under the B
authorization (S2B-01…S2B-17). **No physical schema, field definition, code, migration, test, deployment, requirement,
map change or canonical capture is authorized.** The Howl state machine and UX remain separate and unauthorized.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Owner decision | **P-1** |
| Establishes | two dedicated lifecycle change kinds; independent per-dimension version sequences; baseline formation value; contiguous versions; same-dimension serialization |
| Does not establish | physical schema/fields, code, migrations, tests, deployment, requirements, canon, Howl/UX |
| Acceptance test | **S-3** — O-KR12, K-R14, §6.4 |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed** ·
**Superseded**.

---

## 1. The decision

**P1D-01 (Observed — the owner's selection).** The owner selected **P-1** on 2026-11-20: post-formation `scope` and
`knowledge_class` changes are represented by **two dedicated lifecycle event kinds**, with independent version
sequences.

**P1D-02 (Observed — what the decision establishes).** The **logical representation** for S-2 is established. It is a
memory-contract decision, not an implementation decision.

**P1D-03 (Observed — what it does not establish).** Physical schema and field definitions; code, migrations and
implementation; tests, probes and deployment; requirement-map changes; canonical capture; Howl state-machine design
and UX. Those remain unauthorized (P1D-25).

---

## 2. The representation

**P1D-04 (Observed — the two kinds).**

| Kind | Records |
|---|---|
| `scope_changed` | a change to a memory's `scope` |
| `class_changed` | a change to a memory's `knowledge_class` |

These are **logical kind names**. Physical names, field definitions, storage and encoding are not established
(P1D-25).

**P1D-05 (Observed — independent sequences).** Each memory maintains:

| Sequence | Values |
|---|---|
| `scope_version` | 1, 2, 3 … |
| `class_version` | 1, 2, 3 … |

The sequences are independent. A scope change advances only `scope_version`; a class change advances only
`class_version` (S2C-05; S2F-04).

**P1D-06 (Observed — baseline).** The **formation value** is the baseline for each dimension. It has no change entry.
If no change entry exists for a dimension, the projection uses the formation value (S2C-04).

**P1D-07 (Observed — assignment).** The first accepted change to a dimension is version **1**; each subsequent
accepted change is **previous maximum + 1**. Versions are positive, **contiguous**, and unique within
`(memory, dimension)` (S2C-05).

**P1D-08 (Observed — serialization).** **Same-dimension writes are serialized** (or one is rejected and retried).
Cross-dimension writes are independent and do not contend (S2C-07/08; S2F-04).

**P1D-09 (Observed — projection).** For each memory and dimension: start from the formation value; apply the change
entry with the highest version for that dimension; if none, the formation value is current (S2I-20). The stored
current value is a cache, not a mutable authority (K-R12).

**P1D-10 (Observed — check).** Fold the lifecycle entries, apply the projection, compare to the stored value. A
mismatch is an integrity error surfaced to `inspection`, never silently resolved (§6.4). A duplicate version is an
integrity error/conflict (S2F-14/15).

**P1D-11 (Observed — no stored join).** The projection is reconstructed from the lifecycle stream alone; no acquisition
event, external operation or stored cross-stream join is consulted (K-R14; K-R13).

---

## 3. What the representation gives Howl

**P1D-12 (Observed — illustrative reconstruction).** The memory contract can reconstruct, independently:

| History | Illustrative sequence |
|---|---|
| **scope** | `character_private` (character-specific) → `scenario` (scope change 1) → `world` (scope change 2) |
| **knowledge class** | an independent progression (illustratively, rumour → belief → established knowledge); the exact `knowledge_class` values are the existing field vocabulary and are **not renamed** here |

**P1D-13 (Calculated — the two histories remain independent).** A scope change does not advance the class sequence,
and a class change does not advance the scope sequence. The state machine can consume either history without the other
(S2H-06…S2H-10).

**P1D-14 (Observed — not state-machine design).** These are memory-contract reconstructions, not rules for which
memories participate in a context. Howl's selection, switching, combination, activation, deactivation and restore
remain Howl's responsibility (S2H-11…S2H-14).

**P1D-15 (Observed — context inactivity is not a lifecycle change).** A memory removed from the active context is not
thereby retracted, superseded, corrected or otherwise changed in lifecycle `status` (S2H-12).

---

## 4. Preserved limitations

**P1D-16 (Observed — missing internal version).** A missing internal version is detectable as a gap. Its
classification depends on completeness knowledge:

| Condition | Classification |
|---|---|
| the sequence is **known complete** | **integrity error** |
| completeness is **not known** | **incomplete/untrusted**, not automatically corruption |

The surviving higher version identifies the dimension; it is held pending the missing predecessor (S2F-08/09).

**P1D-17 (Observed — missing latest event).** A missing **latest** event cannot be detected by sequence numbers alone:
no later version exists to expose the gap. The projection may be stale without being detectably incomplete
(S2F-10/11; S2C-12).

**P1D-18 (Observed — out-of-order arrival).** A higher version observed before a lower one is **recoverable**, not
automatically corruption. The higher version is held pending the missing predecessor; when the predecessor arrives,
the projection applies in version order. If it never arrives and completeness is known, it becomes an integrity error
(S2F-12/13).

**P1D-19 (Observed — duplicate version).** Two entries with the same `(memory, dimension, version)` are a **duplicate**
and an integrity error/conflict. Identical duplicates are still an error; the sequence invariant is violated
(S2F-14/15).

**P1D-20 (Observed — no silent repair).** There is no silent gap repair and no presentation of an untrusted projection
as authoritative. The affected dimension is surfaced to `inspection`; the last contiguous value may be shown as last
known, not as authoritative (S2C-11; S2F-06).

**P1D-21 (Observed — the adversarial correction is preserved).** The earlier claim that P-1 has an integrity advantage
over P-2 is **Narrowed/Superseded**: under identical per-dimension sequences and strict validation, P-1 and P-2 are
equivalent on internal removal, latest removal, out-of-order arrival, duplicates and rebuild (S2F-18…S2F-21). The P-1
selection rests on **narrow validation and self-description**, not on a detection advantage (S2F-22/23).

---

## 5. What is established

**P1D-22 (Observed — established logical representation).**

| Established | Content |
|---|---|
| event kinds | `scope_changed`, `class_changed` |
| sequences | `scope_version`, `class_version`; independent, contiguous |
| baseline | formation value |
| projection | baseline + highest version per dimension |
| check | fold/compare; mismatch and duplicate are integrity errors |
| same-dimension writes | serialized or rejected/retried |
| cross-dimension writes | independent |
| no stored join | projection from the lifecycle stream alone |
| no silent repair | untrusted projection surfaced, not authoritative |

**P1D-23 (Calculated — what "established" means).** The **logical** contract is decided. It is the contract the Howl
state machine will be designed around; it is not a schema, an implementation, a test result or a canonical edit
(P1D-25).

---

## 6. What remains unauthorized

**P1D-24 (Observed — the unauthorized set).**

| Unauthorized | State |
|---|---|
| **physical schema and field definitions** | not created |
| **code, migrations, implementation** | not created |
| **tests, probes, model calls, deployment** | **0** |
| **requirement-map changes** | **0**; map stays **118** |
| **canonical capture** (r33, v9, ledger, registry) | **0** |
| **Howl state-machine design and UX** | not designed |
| **authority / permission / sufficiency** | outside S-2; not decided |
| **validity / warrant / truth** | outside S-2; not decided |
| **causes of class/scope change** | D5/K-R9/K-R10/I13; not decided |
| **N-b / N-c / N-a fold** | separate; not decided |
| **§11 promotion / O-8a / v9 `domain`** | separate axes; not decided |

**P1D-25 (Observed — the decision is not implementation permission).** Selecting P-1 does not authorize a schema, a
field, a migration, code, a test, a deployment, a requirement or a canonical change. Each remains a separate
authorization (S2R-25).

---

## 7. What remains open

**P1D-26 (Calculated — logical details not settled by the selection).**

| Open | Why |
|---|---|
| **completeness marker** | how the system knows a sequence should be complete (checkpoint/terminal version); without it, a gap is incomplete/untrusted |
| **recovery policy** | how long to wait for a missing predecessor before escalating |
| **concurrency mechanism** | serialize vs compare-and-swap; the invariant is fixed, the mechanism is not |
| **exact kind/field names and physical schema** | implementation; unauthorized |
| **prior values for audit** | optional; not required by the projection |
| **S-1 fold scope** | separate downstream decision |
| **N-c** | status ordering; untouched |
| **requirement capture** | separate gate |
| **implementation / tests** | separate gate |

---

## 8. Acceptance conditions

**P1D-27 (Calculated — S-3 must verify).**

| Condition | Test |
|---|---|
| **recording** | every actual class/scope change has a lifecycle-stream record |
| **projection** | the current value is derived from the stream; the cache is not authority |
| **independent rebuild** | the value is reconstructable from the lifecycle stream alone (K-R14) |
| **check** | a mismatch or duplicate is an integrity error (§6.4) |
| **no stored join** | no cross-stream join is persisted (K-R13) |
| **no N-c spillover** | the version sequences order only the change entries; status ordering remains separate |
| **no context/status conflation** | context inactivity is not a lifecycle status (S2H-12) |
| **limits preserved** | the missing-latest limitation is stated; no silent repair |

---

## 9. Provenance and correction

**P1D-28 (Observed — the chain).**

| Step | Record |
|---|---|
| B authorization | S2B-01…S2B-17 |
| Ownership/dependency audit | S2A-01…S2A-52 |
| Representation investigation; P-1 recommended | S2I-01…S2I-39 |
| Versioning defined; P-1 vs P-2 compared | S2C-01…S2C-32 |
| Adversarial comparison; P-1 advantage corrected | S2F-01…S2F-28 |
| **P-1 selected** | this record |

**P1D-29 (Observed — the correction).** S2C-27's claim of a P-1 integrity advantage is **Narrowed/Superseded** by
S2F-18. The selection is consistent with the adversarial finding: P-1 is chosen on **narrow validation and
self-description** (S2F-22/23), not on a detection advantage over P-2.

---

## 10. What this decision does not do

**P1D-30 (Observed).**

| Untouched | State |
|---|---|
| physical schema / fields | **0** |
| code / migrations / implementation | **0** |
| tests / probes / deployment | **0** |
| requirements / map | **0**; **118** |
| canonical capture | **0** |
| Howl state machine / UX | **0** |
| authority / validity / causes | not decided |
| N-b / N-c / N-a / S-1 | not decided |
| §11 / O-8a / v9 domain | not decided |

---

## 11. Traceability

**P1D-31 (Observed).**

| What | Where |
|---|---|
| P-1 kinds and sequences | P1D-04…P1D-10 |
| Baseline and projection | S2C-04; S2I-20 |
| Versioning rules | S2C-05…S2C-13 |
| Failure taxonomy | S2F-06/07 |
| Adversarial equivalence | S2F-18…S2F-21 |
| Selection rationale | S2F-22/23 |
| Unauthorized gates | S2R-25; P1D-24/25 |
| Acceptance test | P1D-27; S2R-24 |
| Context inactivity ≠ lifecycle status | S2H-12 |
| Map unchanged | 118 rows |

---

## 12. Record status

**P1D-32 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Owner decision recorded | **P-1** |
| Implementation / tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **owner decision record** — no canonical capture |

**End of decision.** **P-1 is the owner's selection** (P1D-01). The S-2 **logical representation** is established:
`scope_changed` and `class_changed`, independent `scope_version`/`class_version` sequences, formation value as
baseline, contiguous versions, same-dimension serialization, projection from the lifecycle stream, and a check that
treats mismatch/duplicates as integrity errors (P1D-04…P1D-10). **The adversarial findings are preserved**: missing
internal versions are detectable as gaps; missing latest events are not detectable by sequence alone; out-of-order
arrivals are held; duplicates are errors; no silent repair (P1D-16…P1D-21). **The earlier P-1 integrity advantage is
Narrowed/Superseded**; the selection rests on narrow validation and self-description (P1D-21/29). **No physical
schema, code, migration, test, deployment, requirement, map change, canonical capture, state-machine design or UX is
authorized** (P1D-24/25).
