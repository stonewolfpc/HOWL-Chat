# MEMORY-09 — S-2 Adversarial Comparison

## P-1 and P-2 under identical per-dimension sequences

**Status:** S-2 adversarial comparison. **No carrier is adopted. No schema, code, migration, test, requirement, map
change or canonical capture.** The comparison corrects the earlier claim that P-1 has an integrity advantage over P-2.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Question | with identical per-dimension sequences and strict validation, does P-1 detect anything P-2 cannot? |
| Scenarios | internal removal · latest removal · out-of-order arrival · duplicate version · rebuild after each |
| Finding | **P-1 and P-2 are equivalent on all five**; the earlier P-1 integrity advantage is **not supported** |
| Consequence | the choice rests on **simplicity and maintainability**, not detection |
| Selection | **none**; P-1/P-2 remain open |
| Tests executed | **0** — the scenarios are logical, not executed |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed** ·
**Superseded**.

---

## 1. Task and the correction

**S2F-01 (Observed — the owner's check).** The earlier comparison claimed P-1 can identify which dimension is
incomplete because the event kind identifies it. The owner's counterexample: scope versions 1,2,3; class versions 1,2;
scope v2 missing. The surviving scope v3 identifies the gap in both P-1 and a P-2 entry carrying
`dimension=scope, scope_version=3`. The claimed advantage needs an explicit test.

**S2F-02 (Observed — the adversarial request).** Give P-1 and P-2 **identical per-dimension sequences**; remove an
internal event; remove the latest event; deliver events out of order; introduce a duplicate; rebuild the projection
after each failure. Do not select a carrier.

**S2F-03 (Observed — the earlier claim).** S2C-27 said P-1 is stronger on all three issues because the dedicated kind
preserves the dimension in the failure case. **That claim is Narrowed/Superseded for the equivalent-sequence case**:
the surviving entry carries the dimension in both designs, and if no entry survives, neither detects the gap. The
earlier comparison was effectively against P-2's **global** variant, not the per-dimension variant the owner
specified.

---

## 2. Identical rules for both designs

**S2F-04 (Calculated — common rules).** Both designs use:

| Rule | Content |
|---|---|
| **baseline** | the formation value; no change entry |
| **sequences** | independent `scope_version` and `class_version` per memory |
| **assignment** | first accepted change = 1; subsequent = previous max + 1; contiguous, unique |
| **projection order** | **version order**, not arrival order |
| **strict validation** | the dimension and value must be valid for the entry |
| **same-dimension writes** | serialized or one rejected/retried |
| **cross-dimension writes** | independent |
| **silent repair** | forbidden |
| **completeness** | a gap is **incomplete/untrusted** unless the system independently knows the sequence should be complete |

**S2F-05 (Calculated — the only difference).** P-1 encodes the dimension in the **kind** (`scope_changed` /
`class_changed`); P-2 encodes it in a **field** (`dimension=scope|class`) on a generic change entry. Both carry a
per-dimension version. Everything else is held constant.

---

## 3. Failure taxonomy

**S2F-06 (Calculated — the owner's distinctions).**

| Observation | Classification | Action |
|---|---|---|
| two entries with the same `(memory, dimension, version)` | **duplicate** — integrity error | reject/conflict; projection untrusted |
| a gap when the sequence is **known complete** | **integrity error** | surface to inspection; no silent repair |
| a higher version observed while a lower version is absent | **incomplete/untrusted** — not yet corruption | hold the higher value; wait for the lower one |
| arrival out of order, all versions eventually present | **recoverable** — not corruption | order by version; apply when contiguous |
| a change never recorded, with no later version | **undetectable by sequence** | the completeness limit (S2C-12) |

**S2F-07 (Calculated — completeness is orthogonal).** Whether a gap is an integrity error or "incomplete/untrusted"
depends on whether the system **knows** the sequence should be complete (e.g., a checkpoint or a known terminal
version). That knowledge is the same for P-1 and P-2; it is not a property of the kind/dimension encoding.

---

## 4. Scenario matrix

### 4.1 Internal event removed

**S2F-08 (Calculated — scope v2 missing; v1 and v3 present).**

| Step | P-1 | P-2 per-dimension |
|---|---|---|
| surviving entries | `scope_changed` v1, v3 | `dimension=scope` v1, v3 |
| gap detected? | **yes** — v2 absent | **yes** — v2 absent |
| dimension known? | **yes** — both surviving entries are scope-kind | **yes** — both surviving entries carry `dimension=scope` |
| projection | v1 trusted; v3 held; scope incomplete | v1 trusted; v3 held; scope incomplete |
| rebuild | from surviving scope entries | from surviving scope entries |
| difference | **none** | **none** |

**S2F-09 (Calculated — conclusion).** The surviving v3 carries the dimension in both designs. The missing v2's kind
or tag is irrelevant to detecting the gap. **P-1 and P-2 are equivalent.** If completeness is known, both classify
the gap as an integrity error; if not, both classify it as incomplete/untrusted.

### 4.2 Latest event removed

**S2F-10 (Calculated — scope v3 missing; v1 and v2 present).**

| Step | P-1 | P-2 per-dimension |
|---|---|---|
| gap? | **no** — no later version exists | **no** — no later version exists |
| projection | v2 (stale but undetectable) | v2 (stale but undetectable) |
| detection | **none** | **none** |

**S2F-11 (Calculated — conclusion).** Neither detects the missing latest change. This is the completeness limit
(S2C-12). A completeness marker or a later version is required to detect it; the encoding does not help. **Equivalent.**

### 4.3 Out-of-order arrival

**S2F-12 (Calculated — v3 arrives before v2).**

| Step | P-1 | P-2 per-dimension |
|---|---|---|
| classification | **recoverable**, not corruption | **recoverable**, not corruption |
| action | hold v3; mark scope incomplete; wait for v2 | hold v3; mark scope incomplete; wait for v2 |
| when v2 arrives | apply v2 then v3; projection v3 | apply v2 then v3; projection v3 |
| if v2 never arrives, completeness unknown | remains incomplete/untrusted | remains incomplete/untrusted |
| if v2 never arrives, completeness known | integrity error | integrity error |

**S2F-13 (Calculated — conclusion).** Both order by version, not arrival. The dimension is known from the surviving
v3 in both designs. **Equivalent.** The key rule is "do not publish a higher version as current while a lower version
is absent," not the kind/dimension encoding.

### 4.4 Duplicate version

**S2F-14 (Calculated — two scope v2 entries).**

| Step | P-1 | P-2 per-dimension |
|---|---|---|
| detection | **yes** — duplicate `(memory, scope, 2)` | **yes** — duplicate `(memory, scope, 2)` |
| classification | integrity error/conflict | integrity error/conflict |
| action | reject/conflict; projection untrusted | reject/conflict; projection untrusted |
| difference | **none** | **none** |

**S2F-15 (Calculated — conclusion).** Identical duplicates are still an error: the sequence invariant is violated.
The write path must reject duplicates. **Equivalent.**

### 4.5 Rebuild after each failure

**S2F-16 (Calculated — reconstruction from surviving entries).**

| Failure | P-1 rebuild | P-2 rebuild |
|---|---|---|
| internal removal | filter scope-kind; v1 trusted; v3 held | filter `dimension=scope`; v1 trusted; v3 held |
| latest removal | v2; undetectable | v2; undetectable |
| out-of-order | v1 trusted; v3 held; apply on v2 | v1 trusted; v3 held; apply on v2 |
| duplicate | untrusted; surface conflict | untrusted; surface conflict |

**S2F-17 (Calculated — conclusion).** The surviving entries are the same; only the filter differs (kind vs field).
The projection results are identical. **Equivalent.**

---

## 5. Why the earlier advantage does not hold

**S2F-18 (Calculated — the error in S2C-27).** The earlier claim assumed that the missing entry's dimension tag would
be lost, so P-1's kind would preserve it. But the **surviving** entries carry the dimension in both designs. The
missing entry's tag is not available in either design; the gap is identified by the surviving sequence. If the missing
entry is the **only** entry for that dimension, there is no surviving entry to reveal the gap — and neither design
detects it. **The kind/field distinction is not a detection advantage.**

**S2F-19 (Calculated — what remains different).**

| Difference | P-1 | P-2 |
|---|---|---|
| vocabulary | two kinds | one kind + dimension field |
| validation | narrow: kind fixes the value domain | broad: dimension/value must be validated together |
| entry self-description | kind is structural | dimension field is data; it can be malformed |
| extensibility | a new dimension needs a new kind | a new dimension needs a new allowed value, no new kind |
| failure detection | **equivalent** under per-dimension sequences | **equivalent** under per-dimension sequences |

**S2F-20 (Calculated — a minor robustness difference, not the claimed advantage).** A P-2 entry with a missing or
invalid `dimension` field is malformed and must be rejected; P-1's kind is structural, so a malformed kind is a
different error. This is a **validation-surface** difference, not a sequence-detection advantage. It does not change
S2F-08…S2F-17.

---

## 6. Conclusion

**S2F-21 (Calculated — the finding).** With identical per-dimension sequences and strict validation, **P-1 and P-2 are
equivalent on all five adversarial scenarios**:

| Scenario | P-1 | P-2 | Difference |
|---|---|---|---|
| internal removal | gap detected; incomplete | gap detected; incomplete | none |
| latest removal | undetectable | undetectable | none |
| out-of-order | recoverable; hold higher | recoverable; hold higher | none |
| duplicate | integrity error | integrity error | none |
| rebuild | same surviving entries | same surviving entries | none |

**S2F-22 (Calculated — the choice shifts).** The earlier recommendation of P-1 rested on an integrity advantage that
does not survive the equivalent-sequence test. The choice now rests on **simplicity and maintainability**:

| Criterion | P-1 | P-2 |
|---|---|---|
| **vocabulary** | two kinds | **one kind + field** — more compact |
| **per-entry validation** | **narrow** — kind fixes the domain | broad — dimension/value pairing |
| **extensibility** | new dimension = new kind | **new dimension = new allowed value** |
| **auditability** | **self-describing kind** | dimension is data |
| **failure detection** | equivalent | equivalent |

**S2F-23 (Calculated — a preference, not a selection).** P-1 is **slightly preferable on narrow validation and
self-description**; P-2 is **slightly preferable on vocabulary compactness and extensibility**. Neither difference is
an integrity advantage. The owner's choice should be based on which simplicity matters more for this system. **No
carrier is selected.**

**S2F-24 (Observed — no adoption).** P-1 and P-2 remain open. The owner's next decision settles the logical
representation (S2I-30; S2B-05).

---

## 7. What remains open

**S2F-25 (Calculated).**

| Open | Why |
|---|---|
| P-1 vs P-2 | simplicity/maintainability choice; no integrity difference |
| version vs predecessor link | both avoid N-c; not settled |
| completeness marker | how the system knows a sequence should be complete |
| recovery policy | how long to wait for a missing lower version before escalating |
| concurrency mechanism | serialize vs compare-and-swap; write-path design |
| exact kind/field names | logical detail; no physical schema |
| S-1 fold scope | separate downstream decision |
| N-c | status ordering; untouched |
| implementation / requirement capture | separate gates |

---

## 8. What this does not do

**S2F-26 (Observed).**

| Untouched | State |
|---|---|
| **carrier selection** | **0** |
| **schema / kind / field / key** | none created or implemented |
| **completeness marker** | not implemented; an open design question |
| **authority / validity / causes** | fenced, not decided |
| **N-b / N-c / N-a** | not decided |
| **state machine / UX** | not designed |
| **implementation / tests** | **0** |
| **requirement capture** | **0**; map **118** |
| **canonical capture** | **0** |

---

## 9. Traceability

**S2F-27 (Observed).**

| What | Where |
|---|---|
| Earlier P-1 advantage claim | S2C-27 |
| Earlier versioning rules | S2C-04…S2C-13 |
| Identical-sequence comparison | S2F-04…S2F-17 |
| Failure taxonomy | S2F-06/07 |
| Correction | S2F-03, S2F-18…S2F-20 |
| Simplicity trade-off | S2F-22/23 |
| B authorization; baseline | S2B-01…S2B-17 |
| P-1/P-2 candidates | S2I-04…S2I-15 |
| Context inactivity ≠ lifecycle status | S2H-12 |
| Map unchanged | 118 rows |

---

## 10. Record status

**S2F-28 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Carrier selected | **0** |
| Implementation / tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **adversarial comparison** — no canonical capture |

**End of comparison.** With **identical per-dimension sequences and strict validation**, P-1 and P-2 are **equivalent**
on internal removal, latest removal, out-of-order arrival, duplicate versions and rebuild after each failure
(S2F-08…S2F-21). The earlier claim that P-1 preserves the dimension in the failure case is **Narrowed/Superseded**:
the surviving entry carries the dimension in both designs, and if no entry survives, neither detects the gap
(S2F-18). A missing **latest** event remains undetectable by sequence alone (S2F-10/11). The choice now rests on
**simplicity and maintainability** — P-1 narrow validation and self-description, P-2 compact vocabulary and
extensibility (S2F-22/23). **No carrier is selected, no schema/code/requirement/map change or canonical capture
follows.**
