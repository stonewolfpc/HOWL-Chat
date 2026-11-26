# MEMORY-09 — Formation and Completeness Decision Scopes

## Two separate owner-level decisions: the initial state's representation and missing-latest assurance

**Status:** decision-scope analysis. **P-1 and S-3 remain closed as completed logical decisions. No implementation,
schema, code, migration, test, requirement, map change or canonical capture.** The formation-event question and the
completeness-assurance question are kept separate; neither is decided here.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Decision A | is a typed formation event **required** by an existing obligation, or is the record baseline sufficient? |
| Decision B | should a completeness-assurance mechanism be added, and which form detects the precise failure? |
| Finding A | **Unestablished** — K-R12/K-R14/§6.4 support the stronger reading, but no text explicitly names formation; the record-field design supports the weaker reading |
| Finding B | only a durable expected-version boundary detects a missing latest; every form has its own trust/atomicity cost; no form selected |
| Status of P-1/S-3 | **closed**; not reopened |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed**.

---

## 1. Task and boundary

**FCD-01 (Observed — the owner's direction).** Keep P-1 and S-3 closed. Open neither implementation nor canonical
capture. Treat formation-event representation and completeness assurance as **separate owner-level decisions** with
their own scope and authorization boundaries.

**FCD-02 (Observed — the sequence).** First determine whether the missing formation event is required by an existing
obligation, or whether the existing record baseline is sufficient. Then assess completeness-marker options against the
precise failure they detect.

**FCD-03 (Observed — the boundary).** No schema, code, tests, requirements, map change or canonical capture. The map
stays **118**. This record defines decision scopes; it does not decide them.

---

## 2. Decision A — the formation event

### 2.1 The question

**FCD-04 (Calculated — the exact question).** Does an existing obligation require the **initial** `scope` /
`knowledge_class` state to be represented on the lifecycle stream — and therefore a typed formation event (or
equivalent initial-state entry)? Or is the `MemoryRecord`'s formation field sufficient as the baseline, with the
stream governing only post-formation changes?

### 2.2 The evidence

**FCD-05 (Source-derived — the record-field design).** §4 Formation assigns `scope`, `knowledge_class`, `status`
(`:279`); §6.1 lists them as `MemoryRecord` fields (`:340`); §8.1 writes `scope=world`, `knowledge_class=experiential`,
`status=active` (`:404`). This supports the **weaker reading**: the initial state is a record field, and the lifecycle
stream governs changes after formation.

**FCD-06 (Source-derived — K-R12).** "A memory's current status, class and scope are a **projection of that stream**,
not mutable authoritative fields" (`:111`). The initial values are current values; if they are only record fields, the
projection is not purely from the stream. This supports the **stronger reading**.

**FCD-07 (Source-derived — K-R14).** "Each stream is rebuilt and checked independently" (`:113`). If the baseline is a
record field, the lifecycle stream alone cannot rebuild it. This supports the **stronger reading**.

**FCD-08 (Source-derived — §6.4).** "Memory events project to the memory's status/scope/class" (`:366`). If formation
is not an event, the initial values are not projected from memory events. This supports the **stronger reading**.

**FCD-09 (Source-derived — A2-43).** `provisional` has a producer — formation at §11 `:604` — but **no typed event
representation**. This observes the gap; it does not by itself settle which reading is binding.

**FCD-10 (Source-derived — ECI-17).** "The initial `provisional` **must be representable on the lifecycle side**. It
cannot be inferred from the acquisition structure." This is an **Inferred** consequence of K-R12/K-R14/§6.4; it
supports the stronger reading but is not itself a canonical requirement or owner decision.

### 2.3 The finding

**FCD-11 (Unestablished — no explicit obligation names formation).** No canonical text says "formation must be a typed
event." The stronger reading is a **derivation** from K-R12, K-R14 and §6.4; the weaker reading is supported by the
record-field design of §4/§6.1/§8.1. The existing material does not settle which reading is binding.

**FCD-12 (Calculated — what is settled for P-1).** P-1 can use the formation record's `scope`/`knowledge_class` as its
baseline without a new kind (P1F-09/10). P-1 governs **post-formation** changes; it does not depend on a typed
formation event to function.

**FCD-13 (Calculated — what is not settled).** Whether the **initial** state must be stream-derived is an
**interpretation question**. Under the stronger reading, the record baseline is insufficient and a formation event (or
equivalent) is required. Under the weaker reading, the record baseline is sufficient and the stream begins with the
first change.

**FCD-14 (Calculated — this is an owner-level interpretation, not a research gap).** The evidence supports the stronger
reading more naturally, but the text does not compel it. The owner must decide which reading governs the memory
contract. This record does not choose.

### 2.4 Options if the owner opens Decision A

**FCD-15 (Calculated).**

| Option | What it does | Cost |
|---|---|---|
| **Record baseline (status quo)** | formation fields are the baseline; the stream governs changes | the initial state is not stream-derived; the stronger K-R12/K-R14 reading is not satisfied |
| **Typed formation event** | a `memory_formed` (or equivalent) entry carries initial scope/class/status | a new kind; overlaps N-b/`provisional`; needs projection/check rules |
| **Equivalent initial-state entry** | an existing entry (if one can carry the initial values) represents the baseline | depends on the existing entry's payload; may not exist |
| **Initial checkpoint** | a separate baseline snapshot | adds a marker; does not by itself make the stream self-contained |

**FCD-16 (Calculated — what Decision A would decide).** Whether the initial state must be lifecycle-side, and if so,
which representation carries it. It would **not** decide N-b's `provisional` vocabulary, the completeness marker, the
physical schema, or implementation.

**FCD-17 (Calculated — what Decision A would not authorize).** No schema, code, migration, test, requirement capture or
canonical capture. A formation event is a representation decision, not an implementation authorization.

---

## 3. Decision B — completeness assurance

### 3.1 The precise failure

**FCD-18 (Calculated — the failure to detect).** A missing **latest** event. The recorded sequence `1,2,3` and the
sequence `1,2` with `3` lost are both contiguous as recorded. The stream alone cannot distinguish them
(P1D-17; S2F-10/11).

**FCD-19 (Calculated — the only detector).** A **durable expected-version boundary** that knows the highest version
that should exist. Internal gaps are already detectable from the sequence; the boundary adds **latest** detection.

### 3.2 Option assessment against that failure

**FCD-20 (Calculated).**

| Option | Detects missing latest? | Detects internal gap? | Trust/atomicity cost | New representation element? |
|---|---|---|---|---|
| **No marker** | **no** | yes, if a later version survives | none | no |
| **Mutable expected-max counter** | yes, if durable and ahead | yes | atomic with assignment; can be ahead (false gap), behind (missed latest), lost, corrupted, concurrently updated | yes |
| **Periodic checkpoint** | yes, up to the checkpoint | yes | checkpoint trust; changes after it are uncovered until the next | yes |
| **Write-ahead reservation** | yes, for the guarded change | yes | reservation integrity; an orphaned reservation is a false gap | yes |
| **Hash chain + checkpoint** | yes, with a checkpoint | yes, via a chain break | hash/checkpoint trust; no checkpoint ⇒ latest undetected | yes |
| **External manifest** | yes, if ahead | yes | external dependency; must not feed the projection (K-R14) | yes |

**FCD-21 (Calculated — the counter is not the default).** A mutable counter detects the latest loss but carries the
most failure modes: ahead/behind, lost, corrupted, concurrent update, and a false-complete window if updated
prematurely. A checkpoint is safer at the moment it is taken but periodic; a reservation is precise but adds an entry
per change. The form is a separate architectural decision.

**FCD-22 (Calculated — the marker is not P-1).** P-1 established `scope_changed` and `class_changed`; a counter,
checkpoint, reservation or manifest is an additional representation/operational element. Adding one requires a
separate owner decision, not an automatic extension of P-1.

**FCD-23 (Calculated — the boundary rule).** A marker is an integrity aid, not a value authority. If the projection
consults it, K-R14's independent rebuild is at risk. A separate check comparing the stream to the marker is safer
than feeding the marker into the projection.

### 3.3 "Recorded-complete", preserved

**FCD-24 (Calculated — the definition).** **Recorded-complete** means the recorded sequence is **contiguous and
internally consistent** — versions 1..N with no gaps or duplicates, in version order. It does **not** prove that no
event was ever lost. A missing latest may exist without any observable gap.

**FCD-25 (Calculated — the consequence).** "Recorded-complete" is authoritative **relative to the recorded stream**,
not completeness-verified. "Verified-complete" requires a trusted marker. This distinction must survive into any later
specification (P1C-15/16).

---

## 4. Why the decisions are separate

**FCD-26 (Calculated — different problems).**

| Decision A — formation event | Decision B — completeness assurance |
|---|---|
| how the stream represents its **starting point** | whether the stream is **missing entries** |
| an initial-state representation question | a latest-event detection question |
| affects K-R12/K-R14's projection base | affects the confidence that the recorded stream is complete |
| may overlap N-b/`provisional` | does not overlap N-b; adds an integrity aid |
| a new kind/entry if adopted | a counter/checkpoint/reservation/manifest if adopted |

**FCD-27 (Calculated — no bundling).** A formation event does not detect a missing latest, and a completeness marker
does not represent the initial state. Deciding one does not decide the other. Bundling them would repeat the
conflation the earlier S-2 work separated.

---

## 5. What remains closed

**FCD-28 (Observed).**

| Item | State |
|---|---|
| **P-1** | **closed** — established logical representation (P1D-01) |
| **S-3** | **closed** — passed the logical acceptance criteria (S3A-27) |
| **implementation** | **not authorized** |
| **schema / code / tests** | **0** |
| **requirements / map** | **0**; **118** |
| **canonical capture** | **0** |
| **Howl state machine / UX** | **0** |
| **Decision A** | scope defined; not decided |
| **Decision B** | scope defined; not decided |

---

## 6. What this record does not do

**FCD-29 (Observed).**

- It does not decide whether a formation event is required.
- It does not select a completeness-marker form.
- It does not reopen P-1 or S-3.
- It does not authorize implementation, schema, code, tests, requirements or canonical capture.
- It does not decide N-b, N-c, N-a, S-1, §11, O-8a or the v9 domain.

---

## 7. Traceability

**FCD-30 (Observed).**

| What | Where |
|---|---|
| Formation record fields | §4 `:279`; §6.1 `:340`; §8.1 `:404` |
| `provisional` producer, no typed event | A2-43; §11 `:604` |
| Initial state must be lifecycle-side (inferred) | ECI-17 |
| K-R12 / K-R14 / §6.4 | r33 `:111`–`:113`, `:366`–`:367` |
| P-1 formation resolution | P1F-09…P1F-17 |
| Missing-latest limit | P1D-17; S2F-10/11 |
| Recovery states / authority rule | P1C-14…P1C-22 |
| "Recorded-complete" | FCD-24/25; P1C-15/16 |
| P-1 and S-3 closed | P1D-01; S3A-27 |
| Map unchanged | 118 rows |

---

## 8. Record status

**FCD-31 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Decision A decided | **0** — scope only |
| Decision B decided | **0** — scope only |
| Implementation / tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **decision-scope analysis** — no canonical capture |

**End of scopes.** **P-1 and S-3 remain closed** (FCD-28). **Decision A** — whether the initial state must be
stream-derived — is **Unestablished**: K-R12/K-R14/§6.4 support the stronger reading, while §4/§6.1/§8.1 support the
record baseline; no text explicitly names formation, so the choice is an **owner-level interpretation** (FCD-11/14).
**Decision B** — completeness assurance — targets the **missing-latest** failure; only a durable expected-version
boundary detects it, and every form has its own trust/atomicity cost; no form is selected (FCD-18…FCD-23).
**"Recorded-complete" means contiguous and internally consistent relative to the recorded stream, not proof of no
loss** (FCD-24/25). **The two decisions are separate and must not be bundled** (FCD-26/27). **No implementation,
schema, requirement, map change or canonical capture follows.**
