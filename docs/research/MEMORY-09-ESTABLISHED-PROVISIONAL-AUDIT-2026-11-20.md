# MEMORY-09 — `established` and `provisional` Audit

## Neither existing entry faithfully carries the complete initial state

**Status:** audit. **No design selected. No schema, code, migration, test, requirement, map change or canonical
capture.** P-1 and S-3 remain closed; Decision B remains untouched.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Question | can either existing entry carry the initial `scope`, `knowledge_class` and `status` from the stream alone, without semantic distortion? |
| Finding | **neither** — `provisional` omits `scope` and forces N-b; `established` has no payload and conflates formation with a lifecycle transition; neither has an order key |
| Recommendation | document the gap; a **dedicated formation representation** is the least distorting, to be designed separately — no kind/name selected here |
| P-1 / S-3 | **closed**; not reopened |
| Decision B | **untouched** |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. Task and criteria

**EPA-01 (Observed — the task).** Before choosing a representation, audit the existing `established` and `provisional`
entries against the obligation adopted in Decision A: the initial `scope`, `knowledge_class` and `status` must be
reconstructible from the authoritative stream, with `MemoryRecord` fields as projections rather than baselines.

**EPA-02 (Observed — the five criteria).**

1. reconstructs all three initial values from the stream alone;
2. preserves the distinction between formation and later lifecycle transitions;
3. does not silently resolve N-b or formation authority;
4. supports deterministic ordering before subsequent P-1 changes;
5. allows `MemoryRecord` fields to be checked as projections rather than trusted as baselines.

**EPA-03 (Observed — the caution).** Do not select a new kind merely because it appears cleaner. If an existing entry
satisfies the obligation without semantic distortion, it is preferable. If none does, document the gap and recommend a
new representation separately.

---

## 2. The candidates

### 2.1 `established`

**EPA-04 (Source-derived — an event kind).** K-R12 lists `established | superseded | corrected | retracted` as
lifecycle **event kinds** (`:111`). `established` is produced by an event (A2-42).

**EPA-05 (Source-derived — no status of that name).** §5's status field values are `provisional | active | superseded |
retracted | corrected` (`:322`). `established` is **not** a status value. A2-43(ii): `established` (event) and `active`
(status) are the same fact under two names with **no stated mapping**.

**EPA-06 (Observed — no payload).** No lifecycle-event schema exists (A2-11/A2-12; ECI-11; KRO-13). `established` is a
**name without a schema**: no text states that it carries `scope`, `knowledge_class` or `status`.

### 2.2 `provisional`

**EPA-07 (Source-derived — a status value and the §11 entry state).** §5 lists `provisional` among the status values
(`:322`). §11's diagram starts at `provisional (status, class, attribution set)` (`:604`). It is the **entry state**
from which promotion to `active`, `superseded`, `corrected` or `retracted` proceeds.

**EPA-08 (Source-derived — a producer, but no typed event).** `provisional` has a producer — formation at `:604` — but
**no typed event representation** (A2-43, part (i) superseded). This is the pre-existing gap.

**EPA-09 (Source-derived — the payload omits `scope`).** The §11 parenthetical is `(status, class, attribution set)`.
It names status, class and attribution — **not `scope`**. §4 and §8.1 assign `scope` at formation, but the `provisional`
entry's stated payload does not carry it.

### 2.3 `active`

**EPA-10 (Source-derived — a status with no named event).** §5 lists `active` (`:322`); §8.1 formation writes
`status=active` (`:404`). A2-42: `active` has **no named event**. NAE-12/NAE-25: a fold over lifecycle events must
derive `active` from `established` plus promotion conditions, needs a missing kind, or treats `active` as outside the
fold.

---

## 3. Audit against the five criteria

**EPA-11 (Calculated — the table).**

| Criterion | `provisional` | `established` |
|---|---|---|
| **1. all three values from the stream** | **fails** — carries status/class/attribution, **not scope** | **fails / unestablished** — no payload/schema; cannot show it carries any of the three |
| **2. formation vs later transition** | **passes** — §11 starts at `provisional`, then promotes | **fails** — it is a lifecycle event kind; using it at formation conflates formation with a transition |
| **3. no silent N-b / authority** | **fails on N-b** — choosing `provisional` fixes the initial status; authority is orthogonal | **fails on N-b** — `established` ↔ `active` unmapped; choosing it fixes the status toward established/active |
| **4. deterministic ordering** | **fails** — no typed event, no order key | **fails** — no schema, no order key |
| **5. record fields checkable** | **fails** — `scope` is not carried, so it cannot be checked from the stream | **fails** — no payload to check against |

**EPA-12 (Calculated — `provisional` passes only criterion 2).** It has the right **shape** — an entry state before
later transitions — but its payload is incomplete and its name fixes the N-b choice.

**EPA-13 (Calculated — `established` passes none).** It has no defined payload, its relation to `active` is unstated,
and using it at formation makes formation look like a later lifecycle transition.

---

## 4. Criterion-by-criterion

### 4.1 Criterion 1 — all three values

**EPA-14 (Calculated — `provisional`).** Status ✔, class ✔, attribution ✔; **scope ✘**. A rebuild from the stream
would have to obtain `scope` from the record, which is exactly the external baseline Decision A rejects. Adding
`scope` to the parenthetical is a payload extension, not a faithful use of the entry as specified.

**EPA-15 (Calculated — `established`).** No text states its payload. It could be **defined** to carry the three values,
but that definition does not exist. As specified, the claim that it carries them is **Unestablished**.

### 4.2 Criterion 2 — formation vs later transition

**EPA-16 (Calculated — `provisional` preserves the distinction).** §11 places `provisional` at the start and promotes
later. Formation is the entry state; later transitions are separate. This is the one criterion it satisfies.

**EPA-17 (Calculated — `established` does not).** `established` is one of K-R12's lifecycle event kinds, alongside
`superseded`, `corrected` and `retracted`. Using it as the formation entry makes formation one of the lifecycle
transitions, unless the kind is redefined. That redefinition is a semantic change, not a faithful reuse.

### 4.3 Criterion 3 — N-b and authority

**EPA-18 (Calculated — `provisional` fixes N-b).** Choosing `provisional` as the formation entry says the initial
status is `provisional`. §8.1 says formation writes `active`. That is the §8.1 collision; choosing `provisional`
resolves it one way without an N-b decision.

**EPA-19 (Calculated — `established` fixes N-b differently).** `established` (event) and `active` (status) have no
stated mapping. Using `established` at formation implies the initial status is established/active, resolving the
collision the other way. It is still a silent N-b choice.

**EPA-20 (Calculated — authority is orthogonal).** Neither entry decides who may form. Formation authority remains
open and separate; the criterion fails on N-b, not on authority.

### 4.4 Criterion 4 — deterministic ordering

**EPA-21 (Calculated — neither supplies an order).** The lifecycle stream has no order key and no fold rule (A2-41;
DS-09), and no lifecycle schema exists (A2-11/A2-12; ECI-11). `provisional` has no typed event to position;
`established` has a name but no order or version. Neither can place itself before P-1's version-1 change.

**EPA-22 (Calculated — this is a general gap, not entry-specific).** The formation entry needs an order/identity
before the first `scope_changed`/`class_changed`. That requirement applies to any option; neither existing entry
supplies it as specified.

### 4.5 Criterion 5 — record fields checkable

**EPA-23 (Calculated — `provisional`).** Status and class can be checked; `scope` cannot, because the entry does not
carry it. The record's `scope` would remain an unchecked baseline.

**EPA-24 (Calculated — `established`).** With no payload, there is nothing to check the record against. The projection
check would have no stream input for the initial values.

---

## 5. Conclusion

**EPA-25 (Calculated — neither is a faithful carrier).** `provisional` fails criteria 1, 3, 4 and 5; it passes only 2.
`established` fails 1, 2, 3, 4 and 5. The gap is real: **no existing entry, as specified, carries the complete initial
state while preserving its meaning and leaving N-b open.**

**EPA-26 (Calculated — repairing `provisional` is a redesign).** Adding `scope` to its payload and neutralizing its
status would make it a **generic formation entry with a provisional-sounding name**. That is a new representation, not
a faithful reuse of `provisional`. It also does not solve the order key.

**EPA-27 (Calculated — repairing `established` is a redesign).** Defining a payload and an `established`/`active`
mapping would change what `established` means and would still conflate formation with a lifecycle transition. That is
a semantic redesign, not a faithful reuse.

**EPA-28 (Observed — the gap is documented).** The absence is not a reason to force an existing entry. It is the
finding the audit was asked to produce.

---

## 6. Recommendation

**EPA-29 (Calculated — a dedicated formation representation is the least distorting).** Because both existing entries
fail on semantics or payload, a **dedicated formation entry** is the cleanest carrier. This recommendation rests on
the five-criterion failures, **not** on a preference for a new kind.

**EPA-30 (Calculated — what the dedicated entry must provide).**

| Requirement | Why |
|---|---|
| carry `scope` and `knowledge_class` | Decision A; the record fields must be checkable |
| carry the initial `status` **structurally**, with its value/vocabulary explicitly deferred if N-b is unresolved | avoid silently choosing `provisional` or `active` |
| have a deterministic order/identity before the first P-1 change | criterion 4 |
| not be a lifecycle transition kind | preserve formation vs later transitions |
| allow the `MemoryRecord` fields to be checked against it | criterion 5 |
| not decide formation authority | authority remains separate |

**EPA-31 (Calculated — N-b deferral is explicit).** The dedicated entry can satisfy `scope` and `knowledge_class` now.
The `status` value is blocked by N-b; it should be carried as a field whose value/vocabulary is **explicitly deferred**
to N-b, not filled with `provisional` or `active` by implication. Decision A's status obligation is then conditional on
N-b, and that condition is stated rather than hidden.

**EPA-32 (Calculated — no kind, name, payload or schema selected).** The dedicated entry's kind name, exact payload,
order key and schema are for the formation-representation decision, not this audit. This record recommends the
**direction**, not the design.

**EPA-33 (Calculated — if the owner prefers an existing entry).** The only path is to redefine its payload and
semantics explicitly; that redefinition should be evaluated in the same decision, not assumed to be a faithful reuse.
The audit does not forbid that path; it records that neither entry passes as specified.

---

## 7. What remains open

**EPA-34 (Calculated).**

| Open | State |
|---|---|
| **formation representation** | dedicated entry recommended; design not selected |
| **N-b / `provisional`** | open; the status value is deferred, not decided |
| **§8.1 collision** | open; not decided |
| **formation authority** | open; orthogonal |
| **order key / lifecycle schema** | open; required by any option |
| **Decision B** | untouched |
| **schema / code / tests / requirements / canonical** | not authorized |

---

## 8. What this record does not do

**EPA-35 (Observed).**

| Untouched | State |
|---|---|
| **P-1 / S-3** | closed; not reopened |
| **Decision A** | resolved Yes; not reopened |
| **formation design** | not selected |
| **N-b / §8.1** | not decided |
| **Decision B** | untouched |
| **schema / code / tests** | **0** |
| **requirements / map** | **0**; **118** |
| **canonical capture** | **0** |
| **Howl state machine / UX** | **0** |

---

## 9. Traceability

**EPA-36 (Observed).**

| What | Where |
|---|---|
| `established` is an event kind | K-R12 `:111`; A2-42 |
| `provisional` is a status; §11 entry state | `:322`; `:604` |
| `provisional` payload omits `scope` | `:604` |
| `active` has no named event; §8.1 writes `active` | `:322`; `:404`; A2-42; NAE-12/25 |
| No lifecycle schema | A2-11/A2-12; ECI-11; KRO-13 |
| No order key / fold rule | A2-41; DS-09 |
| Entry must be lifecycle-side | ECI-17 |
| Entry architectures | NAB-16/18 |
| Decision A owner determination | DAO-01…DAO-14 |
| Formation resolution | P1F-09…P1F-17 |
| Map unchanged | 118 rows |

---

## 10. Record status

**EPA-37 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Existing entry selected | **0** — neither passes |
| New representation selected | **0** — direction only |
| N-b / §8.1 decided | **0** |
| Implementation / tests | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **audit** — no canonical capture |

**End of audit.** Neither existing entry is a faithful carrier: `provisional` omits `scope` and fixes the N-b choice;
`established` has no payload and conflates formation with a lifecycle transition; neither supplies an order key
(EPA-11…EPA-25). The gap is documented, and a **dedicated formation representation** is recommended as the least
distorting direction — not selected, not named, not schematized (EPA-29…EPA-33). The status value is **explicitly
deferred** to N-b, not filled by implication (EPA-31). **P-1 and S-3 remain closed; Decision B is untouched; no schema,
code, requirement, map change or canonical capture follows** (EPA-35/37).
