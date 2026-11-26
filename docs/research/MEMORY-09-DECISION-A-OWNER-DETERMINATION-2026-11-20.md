# MEMORY-09 — Decision A Owner Determination

## Yes: the initial state must be reconstructible from the authoritative stream

**Status:** **owner decision recorded** at the interpretation level. It is the stronger reading of the existing
requirements, now adopted by the owner — **not** a new requirement, requirement capture, map change or canonical
capture. No schema, code, migration, test or implementation.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Decision | **Yes** |
| Basis | K-R12 + K-R14 + §6.4; §4/§6.1/§8.1 describe the current representation, not an exception |
| Nature | stronger **interpretation** of existing requirements, adopted by the owner |
| Effect | initial `scope`, `knowledge_class` and `status` must be reconstructible from the authoritative stream |
| Record fields | may remain stored projections/caches; **not** the independent source of initial truth |
| Not decided | formation kind, payload, authority, relationship to `provisional`, N-b, §8.1 collision |
| P-1 / S-3 | **closed**; not reopened |
| Map | **118** rows |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. The determination

**DAO-01 (Observed — the owner’s decision).** The owner accepts **Yes**: K-R12, read with K-R14 and §6.4, requires the
initial `scope`, `knowledge_class` and `status` to be reconstructible from the authoritative event stream.

**DAO-02 (Observed — the precision).** This is the **stronger interpretation** of the existing requirements. It was a
recommendation before the owner’s agreement; it is now the owner’s adopted interpretation. It is **not** a new
requirement and does not become one by being adopted.

**DAO-03 (Observed — the basis).** The decisive combination is:

- K-R12 makes current status, class and scope **projections of the authoritative stream**.
- K-R14 requires that stream to be **independently rebuildable and checkable**.
- If reconstruction needs formation-time fields from `MemoryRecord`, those fields function as an **external baseline
  authority**, weakening the independence requirement.
- §4/§6.1/§8.1 describe the existing representation; they do **not** expressly establish an exception to K-R12.

**DAO-04 (Calculated — the rejected qualification).** The counter-reading is plausible, but it adds an unstated
qualification: the stream records changes after formation, while formation fields supply the initial state. The
unqualified wording of K-R12 and K-R14 supports rejecting that qualification.

---

## 2. What Yes preserves

**DAO-05 (Calculated — the obligations).**

| Obligation | Meaning |
|---|---|
| **initial state reconstructible** | initial `scope`, `knowledge_class` and `status` must be derivable from the authoritative stream |
| **record fields are caches** | `MemoryRecord.scope` / `knowledge_class` / `status` may remain stored projections, but they cannot be the independent source of initial truth |
| **no external baseline authority** | reconstruction must not depend on reading formation-time record fields as the authority |
| **stream remains authoritative** | the memory stream is the source of truth for the initial state and subsequent changes |

**DAO-06 (Calculated — what this changes).** The current record-field baseline is **insufficient as the baseline
authority**. A formation representation is required. The record fields may continue as the stored projection/cache,
validated against the stream.

**DAO-07 (Calculated — what this does not change).** P-1’s post-formation representation (`scope_changed` /
`class_changed`) is unaffected. S-3 remains passed. The record fields still exist; their role is narrowed to cache.

---

## 3. What Yes does not decide

**DAO-08 (Calculated — the boundaries).**

| Not decided | Why |
|---|---|
| formation entry **kind name** | a design choice |
| formation entry **payload** | a design choice (scope/class/status/attribution/turn) |
| whether it is a **new kind or an existing entry** | a design choice |
| **formation authority** (who may form) | an authority question, separate |
| **relationship to `provisional`** | N-b remains separate |
| **§8.1 status collision** (`active` vs `provisional`) | N-b/§8.1 remain separate |
| **completeness marker** | Decision B, untouched |
| **schema / code / tests** | implementation |
| **requirement capture / canonical capture** | separate gates |

**DAO-09 (Calculated — representation obligation, not design).** Yes settles only that the initial state must be
stream-reconstructible. It does **not** authorize a particular event shape, payload, authority or implementation.

**DAO-10 (Calculated — no reopening).** P-1 and S-3 remain closed. This determination does not reopen their acceptance
or revise their limitations.

---

## 4. What remains open

**DAO-11 (Calculated).**

| Open | State |
|---|---|
| **formation representation** | a separate, separately bounded decision — see FRS |
| **N-b / `provisional`** | open; not decided by Yes |
| **§8.1 status collision** | open; not decided by Yes |
| **formation authority** | open; representation does not decide it |
| **Decision B — completeness** | untouched; its own decision |
| **implementation / schema / code / tests** | not authorized |
| **requirements / canonical** | not authorized; map **118** |

---

## 5. What this record does not do

**DAO-12 (Observed).**

| Untouched | State |
|---|---|
| **P-1 / S-3** | closed; not reopened |
| **formation design** | not selected |
| **N-b / §8.1** | not decided |
| **Decision B** | untouched |
| **schema / code / tests** | **0** |
| **requirements / map** | **0**; **118** |
| **canonical capture** | **0** |
| **Howl state machine / UX** | **0** |

---

## 6. Traceability

**DAO-13 (Observed).**

| What | Where |
|---|---|
| K-R12 | `:111` |
| K-R14 | `:113` |
| §6.4 | `:366`–`:367` |
| Formation writes record fields | §4 `:279`; §6.1 `:340`; §8.1 `:404` |
| `provisional` producer, no typed event | A2-43; §11 `:604` |
| Initial state must be lifecycle-side (inferred) | ECI-17 |
| Decision A brief | DA-01…DA-36 |
| Decision scopes | FCD-04…FCD-17 |
| P-1 / S-3 closed | P1D-01; S3A-27 |
| Map unchanged | 118 rows |

---

## 7. Record status

**DAO-14 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Decision A | **resolved: Yes**, owner-adopted interpretation |
| Formation design selected | **0** |
| Implementation / tests | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **owner determination** — no canonical capture |

**End of determination.** The owner adopts **Yes** (DAO-01/02): the initial `scope`, `knowledge_class` and `status`
must be **reconstructible from the authoritative stream**; record fields may remain stored projections/caches, but not
the independent source of initial truth (DAO-05/06). This settles the **representation obligation only** (DAO-08/09).
The formation entry’s kind, payload, authority, relationship to `provisional`, N-b and the §8.1 collision remain
**separate** (DAO-08/11). **P-1 and S-3 remain closed; Decision B is untouched; no schema, code, requirement, map change
or canonical capture follows** (DAO-10/12).
