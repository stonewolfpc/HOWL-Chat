# MEMORY-09 — G1 Capture Fixture Contract

## What the capture boundary must prove before any classifier is measured

**Status:** **fixture contract; the test-only fixture subsequently passed (G1R-22).** No production API change,
schema, migration, model fetch, corpus, L-2 amendment or Decision A amendment. The contract remains the definition of
what G1 had to demonstrate; production wiring is still unproven.

| Field | Value |
|---|---|
| Date | **2026-11-21** (local, CST) |
| Gate | **G1 — capture completeness** |
| Fixture cases | **F1–F10**, plus the distinct **character-tell** case |
| Current result | **G1 not passed**; audit complete, fixture not built |
| Requested authorization | **fixture-only test harness**, no production behavior change |
| Production API change | **separate authorization** |
| L-2 / Decision A | **unchanged** |
| BGE / G2 | **not started** |
| Map | **118** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed**.

---

## 1. Preserved conclusions

**GFC-01 (Observed — no collapse).** Capture-incomplete must not be collapsed into classification-unresolved. A
missing `CharacterID`, attribution or authored-asset status is a capture failure, not uncertainty for a classifier to
solve.

**GFC-02 (Observed — no compensation).** A successful classification cannot compensate for incomplete provenance. A
confident class assignment cannot repair a missing source envelope.

**GFC-03 (Observed — character-tell acceptance).** The character-tell case remains an explicit acceptance case: the
memory can remain available for non-class-dependent processing without inventing a class or interrupting the
conversation.

**GFC-04 (Calculated — G1 is a contract, not a classifier test).** G1 passes when the capture boundary demonstrably
produces the required envelope and exposes incompleteness. It does not assign a class and does not use BGE.

---

## 2. Capture envelope

**GFC-05 (Calculated — logical fields).**

| Field | Requirement | Notes |
|---|---|---|
| `source_message_id` | required | stable message identity |
| `source_span` | required | start/end or exact text; provenance |
| `speaker_role` | required | `user` / `assistant` / `narrator` / `system` |
| `source_character_id` | conditional | required for an assistant character message; explicit none for user/system/narrator |
| `source_character_name` | optional label | display only; never ownership |
| `chat_id` | required | isolation |
| `turn_id` / sequence | required | local ordering |
| `event_type` | required | detection signal |
| `evidence` / source text | required | explicit rule evidence |
| `world_id`, `scenario_id` | required | context and authored-asset check |
| `attribution` | conditional | required for experience/belief; observed / told / believed / asserted / authoritative |
| `authored_asset_ref` | conditional | required for ambient/restricted eligibility |
| `content` | required | memory text |
| `existing_class` | optional | legacy only; never trusted as new |
| `capture_state` | required | `complete` or `incomplete` |
| `missing_fields` | required when incomplete | the specific absent fields |
| `determination_scope` | required | which determination the completeness applies to |

**GFC-06 (Calculated — completeness is relative to a determination).** A capture can be complete for ownership but
incomplete for experience/belief, or complete for experiential/character_specific but ineligible for
ambient/restricted. The fixture must record the determination scope, not a single global “complete” flag that hides
the difference.

**GFC-07 (Calculated — no name inference).** A missing `CharacterID` is never repaired by matching `Sender` or a
display name. Name collisions remain visible.

---

## 3. Completeness rules

**GFC-08 (Calculated — global required).** Missing `source_message_id`, `source_span`, `speaker_role`, `chat_id`,
`turn_id`, `event_type`, `evidence`, `world_id`, `scenario_id` or `content` makes the capture **incomplete for every
determination**.

**GFC-09 (Calculated — character identity).** An assistant character message without a stable `source_character_id` is
**incomplete for ownership and classification**. A user/system/narrator message carries an explicit none; that is
complete, not missing.

**GFC-10 (Calculated — attribution).** Missing attribution is **incomplete for the experience/belief determination**.
It does not automatically invalidate ownership, display or provenance.

**GFC-11 (Calculated — authored asset).** Missing `authored_asset_ref` makes `ambient` and `restricted` **ineligible**.
It does not make a character message globally incomplete: auto-captured character speech is simply not an authored
world/scenario definition (K-R8).

**GFC-12 (Calculated — provenance is not optional).** A class assignment without a resolvable source span is not a
valid capture result, regardless of confidence.

---

## 4. F1–F10

**GFC-13 (Calculated — the fixture table).**

| Case | Input | Expected envelope | Incomplete for | Class-eligible? | Non-class processing |
|---|---|---|---|---|---|
| **F1** | character assistant message with stable `CharacterID`, attribution `observed`, world/scenario IDs | complete; source character ID; role assistant; span; attribution | — | yes: `experiential` / `character_specific`; ambient/restricted ineligible without authored ref | retain, display, inspect, ownership/status |
| **F2** | user message | complete; role user; source character **explicit none** | — | **deferred**: user-source class rule not established in G1 | retain, display, inspect |
| **F3** | narrator/system message | complete; role narrator/system; source character none | — | only if `authored_asset_ref` present; otherwise no ambient/restricted | retain, display, inspect |
| **F4** | assistant message with empty `CharacterID` | **incomplete**; missing `source_character_id` | ownership and classification | **no** | retain for repair/inspection only |
| **F5** | two characters with the same display name, stable IDs present | complete; IDs distinguish; no name match | — | yes, using stable ID | retain, display, inspect, ownership/status |
| **F6** | character message without attribution | **incomplete** for experience/belief; other fields present | experience/belief | **no** until attribution or an explicit rule | retain, display, inspect, provenance |
| **F7** | authored world/scenario definition with `authored_asset_ref` | complete; authored ref; role system/narrator | — | yes: `ambient` / `restricted` | retain, display, inspect |
| **F8** | auto-captured character speech about the world | complete; source character; **no authored ref** | — | yes: `experiential` / `character_specific`; ambient/restricted **ineligible** | retain, display, inspect, ownership/status |
| **F9** | legacy record with no class/source | **incomplete/unmigrated**; missing fields listed | classification and provenance | **no** | retain, display, inspect; never guess |
| **F10** | message missing world/scenario ID | **incomplete**; missing context fields | every determination | **no** | retain for repair/inspection only |

**GFC-14 (Calculated — F2/F3 are not failures).** They are capture-complete for source/ownership, but class
eligibility depends on rules not established in G1. The fixture must not invent a user-source or narrator-source class.

**GFC-15 (Calculated — F8 is the ambient trap).** F8 must never be classified `ambient` merely because the content is
about the world. The absence of an authored asset reference is decisive.

**GFC-16 (Calculated — F9 is not unresolved).** A legacy record without a class/source is **capture-incomplete**, not
a classifier uncertainty. It is reported, never guessed.

---

## 5. Character-tell case

**GFC-17 (Calculated — distinct acceptance case).** Character A tells character B about A's past. The captured memory
is about A; the source is A; attribution is `told`; B's acquisition is a separate edge and is not part of G1.

| Field | Expected |
|---|---|
| source character | A, stable ID |
| role | assistant |
| attribution | `told` |
| source span | A's message |
| world/scenario | present |
| class | **not assigned in G1**; eligible for classification, may be unresolved |
| non-class processing | retain, display, inspect, provenance, ownership/status |
| class-sensitive use | waits; no default |
| conversation | not interrupted; no invented class |

**GFC-18 (Calculated — what the case proves).** Capture can preserve the memory and provenance while class remains
unresolved. It does not prove that L-2 authorizes an unresolved class; that is Track 3/G4 and an owner decision.

**GFC-19 (Calculated — missing attribution variant).** If the same case omitted attribution, it would be
**capture-incomplete for experience/belief**, not “unresolved.” The fixture must include both variants to prove the
distinction.

---

## 6. Eligibility and permitted processing

**GFC-20 (Calculated — eligible cases).** F1, F5, F7, F8 and the character-tell case are class-eligible. F2/F3 are
capture-complete but class-eligibility is deferred/conditional. F4, F6, F9, F10 are not class-eligible.

**GFC-21 (Calculated — non-class-dependent processing).** Permitted for complete and incomplete captures, as
applicable:

- retain the candidate/record and its source span;
- display it in memory inspection;
- preserve ownership and status;
- inspect provenance and missing fields;
- report capture-incomplete separately from unresolved.

**GFC-22 (Calculated — not permitted).** No class-gated audience/eligibility, no ambient/restricted assignment without
authored status, no character ownership from a display name, no class from a capture-incomplete case.

**GFC-23 (Calculated — no compensation).** A high classifier confidence cannot make F4, F6, F9 or F10 class-eligible.
The capture boundary decides eligibility before classification.

---

## 7. Fixture-only harness authorization request

**GFC-24 (Observed — requested scope).** A **test-only** harness that:

- defines the logical envelope and the F1–F11 fixture cases;
- uses existing `ChatHistory` / `ChatMeta` / `ChatData` types read-only;
- asserts the expected envelope, missing fields, eligibility and non-class processing;
- uses no production API change, no persistence, no model, no class labels and no L-2 representation;
- fails visibly if the current capture path cannot produce the envelope.

**GFC-25 (Observed — not requested).** No production API change, no `ProcessMessageForEvents` signature change, no new
Wails binding, no storage schema, no classifier, no BGE, no corpus, no L-2/Decision A edit.

**GFC-26 (Calculated — stop condition).** If the fixture cannot be expressed without changing production behavior, the
harness stops and returns for a separate production-API authorization. The test must not be made to pass by changing
the implementation.

**GFC-27 (Calculated — G1 pass condition).** G1 passes only when F1–F10 and the character-tell case pass with the
expected states, capture-incomplete is distinct from unresolved, no class is assigned, and no production behavior was
changed to satisfy the test.

**GFC-28 (Calculated — G1 fail condition).** G1 fails if any missing required field is treated as a class, any name is
used as identity, any capture-incomplete case becomes class-eligible, or the harness requires a production change to
pass.

---

## 8. What this contract does not do

**GFC-29 (Observed).**

| Untouched | State |
|---|---|
| L-2 / Decision A | **unchanged** |
| P-1 / S-3 | closed; not reopened |
| Decision B | untouched |
| schema / code / migration / tests | **0** |
| fixture | **not built** |
| production API | **unchanged** |
| classifier / BGE / corpus / G2 | **0** |
| requirements / map | **0**; **118** |
| canonical capture | **0** |

---

## 9. Traceability

**GFC-30 (Observed).**

| What | Where |
|---|---|
| G1 audit | CCG-01…33 |
| Capture envelope | CCG-15; GFC-05 |
| Capture-incomplete vs unresolved | CCG-20…23; GFC-01/16/19 |
| F1–F10 | CCG-24; GFC-13 |
| Character-tell | ACX-61; GFC-17…19 |
| Ambient authored-only | K-R8; GFC-11/15 |
| No production change | GFC-24…26 |
| L-2 / Decision A | L2D-01/08/10; DAO-05; GFC-29 |
| Map | **118** |

---

## 10. Record status

**GFC-31 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — map **118** |
| G1 | **fixture contract passed; production wiring not established** (G1R-22) |
| Fixture built | **test-only harness; passed** |
| Production API changed | **0** |
| L-2 / Decision A amended | **0** |
| BGE / G2 started | **0** |
| Tests executed | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **fixture contract** — no canonical capture |

**End of contract.** G1 is a capture-boundary test, not a classifier test (GFC-04). F1–F10 define the expected
envelope, the missing fields that make each case capture-incomplete, the cases eligible for classification, and the
non-class processing permitted (GFC-13/20…23). The character-tell case is distinct: capture and provenance survive,
class remains unassigned, and class-sensitive use waits (GFC-17…19). The requested authorization is **fixture-only**;
a production API change is separate, and the test must not be made to pass by changing production behavior
(GFC-24…28). **No L-2 amendment, BGE fetch, G2 start, schema, code, fixture or test follows** (GFC-29/31).
