# TASK-03 — reconciliation contract v3 (consolidation, self-contained)

| Field | Value |
|---|---|
| Status | **Contract v3 draft. Nothing adopted. No implementation authorized.** |
| Supersedes | v2 (row 554) and amendments 555–557 as normative text; they remain historical evidence, unedited |
| Consolidates | rows 548–558, with rows 547 and 554–558 retained as history |
| Governing principle | **One authoritative ordered decision procedure. No state may depend on a condition hidden in an amendment, footnote, transition narrative, or test expectation.** |
| Recovery token | `TASK-03-RECONCILIATION-CONTRACT-V3-2026-11-22` |

Prefixes: `V3-` normative, `VT3-` acceptance test. Labels: **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Superseded**.

**V3-01 — This document is self-contained.** It does not require row 554, 555, 556, 557 or 558 to be read to be executed. §4 is the executable core; every condition that decides a state appears in §4's rows. The older documents are evidence and rationale, not normative text.

**V3-02 — What changed and why.** Rows 555–558 added rules to a table they could not reach. V3 folds them in: eleven contradictions (row 558 §2) are resolved here, each named at the clause that resolves it.

---

## §1 — Terms and the evidence properties

**V3-03 — The question.** Whether a conversation's data is still *available*. Not who deleted it, not that it is unrecoverable, not that the filesystem cooperated.

**V3-04 — Four properties, never interchangeable.**

| Property | Question | Evidence |
|---|---|---|
| **Presence** | is the artefact there at all? | `stat` / listing |
| **Structural validity** | does the *whole* representation conform to its format? | whole-file parse |
| **Content presence** | does a valid representation hold ≥1 record? | parse result |
| **Completeness** | is the committed history still there? | the decision, from both representations, the marker, and attribution |

**V3-05 — `LoadChat` is loadability, not validity.** It is lossy by design: truncated, wrong-shaped and interior-corrupt histories return fewer messages with no error (Observed, row 551). Its count never decides a state or a promotion.

**V3-06 — The two representations.**

```
<rootDir>/chats/<chatID>/
    meta.json          required — identity
    history.ndjson     log      — append representation
    history.json       snapshot — peer representation, not legacy
    world.json · scenario.json · characters.json · lorebooks.json · memories.json
```

`world.json` and `scenario.json` are written only when non-nil; `history.json` is written by save, rewrite and duplication (Source-derived, rows 549/553).

**V3-07 — Structural validity.**

| Format | Valid iff |
|---|---|
| `history.ndjson` | every non-blank line parses as a record, and a non-empty file ends with a newline |
| `history.json` | the whole document parses as an array of records |
| either, zero bytes | **valid and empty** — vacuous validity, no trailing newline required |

**V3-08 — Empty is defined, not inferred.** Empty log = 0 bytes. Empty snapshot = `[]` (2 bytes). A lone `"\n"` (1 byte) is neither valid nor empty and is reported `malformed` (Observed, rows 548/551). *Resolves row 558 CA-02.*

**V3-09 — Void** means no representation holds ≥1 record. A void directory may hold empty files, absent files, or both.

**Usable content** means a structurally valid representation holding ≥1 record, such that every other content-bearing representation stands to it in an admissible prefix relation (V3-39). A structurally invalid companion never contributes and never blocks; it is reported (V3-40). Where no usable content exists and something present is invalid, row 7 decides.

**V3-10 — Attribution** is the question *"can the generation of this content be established?"* It is established only when all three hold:

1. the marker read is **current** (`Saved`, not `Recovered`, not absent);
2. the content's generation is one the marker authorizes (§2);
3. the transition that opened that generation is **settled** (§3).

If any fails, the content is **unattributed**. Attribution is never inferred from record content, record IDs, counts, sizes or timestamps.

**V3-11 — `IsInconsistent` is a reason, not a state.** It maps to Unavailable and fails closed.

---

## §2 — The expectation marker

**V3-12 — Shape.** One document, three fields, persisted through `DurableStore` (exclusive lock, temp write, manifest commit with hash and size, pre-images, `testHookBeforeReplace`):

| Field | Meaning |
|---|---|
| `state` | `none` · `pending` · `committed` · `cleared` |
| `generation` | the lineage the state describes |
| `updated_at` | bookkeeping only; never evidence |

**V3-13 — The marker stores intent, not state.** The six observed states of §4 are always computed. The marker's `state` field is an *intent* record, and calling it state was the confusion behind v2-19. *Resolves row 558 CA-05.*

**V3-14 — Generation semantics.**

| Value | Meaning |
|---|---|
| `(none,0)` | nothing committed, no clear recorded |
| `(pending,G)` | an operation in generation G began and is unconfirmed |
| `(committed,G)` | content is committed in G |
| `(cleared,G)` | G is intentionally retracted; the next commit begins G+1 |

`0→1` on the first confirmed commit. `G→G+1` only when a clear of G has been **recorded and settled** and content is committed again. **A clear alone never increments** — it closes G. A `Recovered` marker never increments anything.

**V3-15 — The generation bounds attribution; it does not identify records.** No per-record generation exists (**Unestablished**). V3 must not pretend otherwise, and no clause may infer lineage from bytes (*Resolves CA-04's root*).

**V3-16 — First write is not a replacement.** A marker document that has never existed has no previous generation to recover:

| Write | Possible outcomes | After death at the replacement point |
|---|---|---|
| **first** marker write | `Saved`, `Absent` | **Absent** — no document, no pre-image, nothing to roll back |
| **replacement** | `Saved`, `Recovered`, `NotSaved` | previous whole document (`Recovered`) or nothing lost (`NotSaved`) |

*Resolves CA-10 and row 557's overstatement: pair safety of §2 is Source-derived and applies to documents that already existed.*

**V3-17 — `Recovered` is not authority.** A `Recovered` marker is internally consistent but may predate the last successful write (`Result.Current()` is false for it — Source-derived). It never establishes that a transition is settled and never authorises promotion. **Row 5 of §4 makes this unconditional.** *Resolves CA-03.*

**V3-18 — The marker neither proves completeness nor authorises cleanup.** No marker value permits deleting any file, including content that is provably retracted.

**V3-19 — Marker ordering.**

| Operation | Order |
|---|---|
| first append | `(pending,1)` → append → `(committed,1)` |
| subsequent append | no marker write |
| correction with survivors | no marker write; both representations rewritten |
| clear | `(cleared,G)` → rewrite both → settled |
| append after clear | `(pending,G+1)` → append → `(committed,G+1)` |
| **duplication** | `(pending,1)` in the **new** directory → write content → `(committed,1)` |

*Resolves CA-09: duplication now has an order, and the copy never inherits the source's generation.*

---

## §3 — The mutation gate

**V3-20 — The gate exists in the contract; it does not exist in the product.** This section specifies behaviour that no current code implements. Nothing here is authorised.

**V3-21 — `gate_enforced` is a declared property of the writer.** The decision table consults it explicitly (§4 rows 10–13). V3 never assumes the gate exists; where the gate's absence makes attribution impossible, the state is Pending with reason `promotion_requires_gate`, named in the table rather than hidden in prose.

**V3-22 — While a generation transition is unresolved, no history mutation may append to or rewrite a representation belonging to it.** The writer either completes the authorized operation or refuses with a typed, caller-visible failure.

**V3-23 — Refusal is surfaced, never discarded.** A caller that swallows the failure is non-conforming. (Source-derived: `_ = a.performSave()` appears at `app-characters.go:242,257` and `app-chat.go:137,148`.)

**V3-24 — Read-then-write is not the gate.** The check, the permission decision and the write happen in **one critical section**. Load-then-Save is the caller-side race.

**V3-25 — The marker needs compare-and-swap.** "Replace `(state, generation)` with the new pair only if it still equals the observed pair", decided inside the same lock hold. `Save` cannot express that; `Load` is lock-free. **A design requirement, unimplemented.**

**V3-26 — Serialization is cross-process.** Two instances may run. A writer that cannot acquire the critical section reports `Busy` and refuses; it never proceeds on a stale observation. An in-process mutex is insufficient. (Supersedes row 547 RC-27.)

**V3-27 — A transition is settled when:**

| Transition | Settled iff |
|---|---|
| first commit (0→1) | marker `(committed,1)`, or the append never happened and the marker is `(none,0)` |
| clear closing G | marker `(cleared,G)` **and** both representations valid **and** void |
| append after clear (G→G+1) | the closing clear is settled, so the log was void when the first G+1 record was written |

**V3-28 — `cleared` alone is not settlement.** A populated snapshot is never dismissed as harmless residue because the log is void. *Preserves row 555 V2-80.*

**V3-29 — Clear retry is idempotent.** An interrupted clear that replaced the log but not the snapshot is finished by re-running the same clear; the retry never treats the void log as a new conversation and never permits an append in between. *Preserves V2-79.*

**V3-30 — Deletion is explicitly outside the gate.** Deleting a conversation destroys both representations, cannot corrupt the pair, and is an admissible exit from an unresolved transition. It is exempt by statement, not by omission, and it is not a resolution: nothing is promoted afterwards. *Resolves CA-07.*

**V3-31 — The inventory rule.** Every mutation entry point is enumerated and gated, except the deletion exemption. A new entry point requires a gate-coverage test. *Preserves V2-81.*

**V3-32 — Ownership boundary.** Creation identity collision (row 550), loader precedence (OQ-3c-a) and `RewriteHistory` resurrection (OQ-3c-b) are separately owned. V3 reports their shapes and never compensates for them.

---

## §4 — The decision procedure

**V3-33 — Applied in order; the first matching row decides.** Rows are mutually exclusive by construction. A row that cannot be reached is a defect of this table.

Inputs, all observed: directory presence · enumeration completeness · marker read outcome (`Saved`/`Recovered`/`Absent`/unreadable) · marker `(state, generation)` · per representation: absent / empty / valid-with-content / invalid · coexistence relationship · index membership · `gate_enforced`.

| # | Condition (all must hold) | State | Reason |
|---|---|---|---|
| 1 | conversation directory absent | **Missing** | `absent` |
| 2 | enumeration incomplete, path denied, or I/O error | **Unavailable** | `unreadable` |
| 3 | marker document present but unreadable (corrupt, no matching pre-image) | **Unavailable** | `marker_unreadable` |
| 4 | `meta.json` absent in a present directory | **Damaged** | `identity_lost` |
| 5 | marker read is `Recovered` (not current) | **Pending reconciliation** | `marker_stale` |
| 6 | two representations valid, both content-bearing, and neither is an exact prefix of the other | **Unavailable** | `inconsistent` |
| 7 | no usable content exists, and at least one present representation is structurally invalid | **Damaged** | `malformed`, or `truncated` when a file ends mid-record |
| 8 | no usable content, every present representation valid and empty or absent, marker `committed` | **Damaged** | `absent` |
| 9 | no usable content, representations valid and empty or absent, marker `pending` | **Pending reconciliation** | `awaiting_content` |
| 10 | no usable content, representations valid and empty or absent, marker `none` or `Absent` | **Healthy** | `empty` |
| 11 | no usable content, representations valid and empty or absent, marker `cleared` | **Healthy** | `closed` |
| 12 | usable content **and** marker is `Absent` or `none` | **Pending reconciliation** | `unattributed` |
| 13 | usable content **and** marker is `cleared` | **Pending reconciliation** | `unattributed` |
| 14 | usable content, marker `pending,1` | **Healthy**, promote `(committed,1)` | — |
| 15 | usable content, marker `committed,1` | **Healthy** | — |
| 16 | usable content, marker `pending,G≥2`, `gate_enforced` | **Healthy**, promote `(committed,G)` | — |
| 17 | usable content, marker `committed,G≥2`, `gate_enforced` | **Healthy** | — |
| 18 | usable content, marker `pending,G≥2`, **not** `gate_enforced` | **Pending reconciliation** | `promotion_requires_gate` |
| 19 | usable content, marker `committed,G≥2`, **not** `gate_enforced` | **Pending reconciliation** | `promotion_requires_gate` |
| 20 | directory present, otherwise decided above as Healthy, with no index entry | **Orphaned** | — |

**V3-34 — Six states, all derived: Healthy · Pending reconciliation · Missing · Damaged · Unavailable · Orphaned.** `Inconsistent` is reason `inconsistent` on Unavailable.

**V3-35 — Row 20 is evaluated only after rows 1–19 would otherwise yield Healthy.** A directory with no index entry and no content is `Healthy`-shape and `Orphaned`; a directory with no index entry and damaged content is Damaged first. Orphaning changes ownership, not availability.

**V3-36 — Row 12 is the fail-closed answer to duplication and to bypassing writers.** `SaveChatAs` writes content into a fresh directory; under V3-19 it writes its own marker, so a conforming duplicate never reaches row 12. A non-conforming duplicate — and probe 12's appended mixture — does, and is Pending, never Healthy. *Resolves CA-01: VT-18's expectation is replaced, not reinterpreted.*

**V3-37 — Rows 18 and 19 are the honest consequence of an unimplemented gate.** For G≥2, content may contain bytes of the closed generation and no evidence distinguishes them (§2 V3-15). While `gate_enforced` is false, those conversations are Pending — visible, unresolved, retried — rather than Healthy on an assumption. **When the gate is implemented and `gate_enforced` is declared true, rows 16 and 17 become reachable; V3 requires no other change.**

**V3-38 — Row 8 versus row 12.** `committed` with void content is Damaged; `none`/`Absent` with no content is Healthy-empty. The marker is what separates "content was committed and is gone" from "nothing was ever committed".

**V3-39 — Coexistence, direction from the writer lifecycle** (Source-derived: `AppendHistoryEntry` writes only the log; `SaveChat` and `RewriteHistory` write both):

| Relationship | Meaning | Effect |
|---|---|---|
| equal | coherent | rows 14–19 as applicable |
| snapshot is an exact prefix of the log | expected lag; the log is the basis | coherent; the shorter set is never claimed complete |
| log is an exact prefix of the snapshot | not expected (torn log, resurrected snapshot) | row 6 |
| no prefix relation | incoherent | row 6 |

Prefix means exact prefix on full record content in order (`sender`, `text`, `role`, `timestamp`, plus `hasImage`/`characterId` when present). **Record IDs never establish it** — they are positional indices (`len-1`, `i`). A prefix is evidence about order, never proof of provenance.

**V3-40 — An invalid companion does not decide a state.** A malformed snapshot beside a valid content-bearing log does not make the conversation Damaged and does not prevent Healthy: the log's availability is established, and the companion is reported with reason `malformed`. The invalid representation is never the basis of a promotion. **Row 7 applies only when no usable content exists at all** — a lone `"\n"`, a torn-only log, or two invalid representations. A void directory holding one malformed empty-ish file is therefore Damaged, not `Healthy empty`: rows 8–11 require every present representation to be valid. *Resolves v2 VT-09's policy and closes the malformed-empty hole.*

**V3-41 — Forbidden transitions.** Unavailable→Healthy; Unavailable→Missing; Damaged→Missing; Damaged→Healthy; Orphaned→Healthy; `inconsistent`→Healthy; any automatic repair.

---

## §5 — Lifecycle transitions

Read-only throughout. "Confirmed" means the marker write returned `Saved`.

### T1 — Creation

| Step | Crash after | Rows | State |
|---|---|---|---|
| directory created | anything | 4 | Damaged `identity_lost` |
| `meta.json` written, `(none,0)`, nothing else | anything | 10 | Healthy `empty` |
| snapshots written (`history.json = []`) | anything | 10 | Healthy `empty` |

### T2 — First append

| Step | Crash after | Rows | State |
|---|---|---|---|
| `(pending,1)` written | — | 9 | Pending `awaiting_content` |
| record appended | torn mid-record | 7 | Damaged `truncated` |
| record appended | intact | 14 | Healthy, promoted `(committed,1)` |
| `(committed,1)` written | — | 15 | Healthy |

### T3 — Subsequent append

| Step | Crash after | Rows | State |
|---|---|---|---|
| record appended | intact | 14/15 | Healthy (snapshot is a prefix) |
| record appended | torn, snapshot intact | 7 | Damaged `truncated`; the snapshot is reported alongside, never used to excuse it |
| nothing written | — | unchanged | unchanged |

### T4 — Correction / retraction with survivors

| Step | Crash after | Rows | State |
|---|---|---|---|
| log rewritten | before snapshot | 6 | Unavailable `inconsistent` — not resolved by choosing |
| snapshot rewritten | — | 14/15 | Healthy |
| nothing written | — | unchanged | unchanged |

### T5 — Clear

| Step | Crash after | Rows | State |
|---|---|---|---|
| `(cleared,G)` written | before log rewrite | 13 | Pending `unattributed` |
| log emptied | before snapshot rewrite | 13 | Pending `unattributed` (snapshot populated) |
| both emptied | — | 11 | Healthy `closed` |

### T6 — Append after clear

| Step | Crash after | Rows | State |
|---|---|---|---|
| `(pending,G+1)` written, nothing appended | — | 9 | Pending `awaiting_content` |
| record appended, gate enforced | intact | 16 | Healthy, promoted `(committed,G+1)` |
| record appended, gate **not** enforced | intact | 18 | Pending `promotion_requires_gate` |
| record appended, gate not enforced | old snapshot still populated | 18, then 6 if both content-bearing and non-prefix | Pending; `inconsistent` reported when it applies |
| append torn | — | 7 | Damaged `truncated`; retracted content never stands in for the new record |

### T7 — Duplication

| Step | Crash after | Rows | State |
|---|---|---|---|
| new directory + `(pending,1)` | — | 9 | Pending `awaiting_content` |
| snapshot written (log absent) | — | 14 | Healthy, promoted `(committed,1)` — one representation is sufficient |
| `(committed,1)` | — | 15 | Healthy |
| source identity collided (row 550 shape) | — | 6 or 12 | reported as collision; never repaired, renamed or split |

### T8 — Deletion

| Step | Crash after | Rows | State |
|---|---|---|---|
| deletion requested during an unresolved transition | permitted (§3 V3-30) | — | no marker write, no promotion |
| directory removed | — | 1 | Missing |
| removal partially completed | directory present, content gone | 8 | Damaged `absent` |

---

## §6 — Crash points and store behaviour

**V3-42 — Marker crash outcomes, by write kind** (V3-16): first write → `Saved` or `Absent`; replacement → `Saved`, `Recovered`, or previous document intact. Never a partial document, never a third state. Source-derived from `TestCrashBeforeReplaceLeavesPreviousDocument`, which is **not re-executed here**.

**V3-43 — `Recovered` reaches row 5 and stops.** It never promotes, never satisfies a settlement, never raises a state. Its pair is whole and internally consistent — that is all it establishes.

**V3-44 — Content crash points never route through the marker.** A torn append, an interrupted rewrite and an aborted clear are decided by rows 6, 7 and 13 from observed bytes. The marker contributes only what §4's rows ask of it.

**V3-45 — `Absent` marker with content is row 12, not row 10.** A first marker write that died leaves no marker; the content that may have been written afterwards is unattributed, not a valid empty conversation.

---

## §7 — Scanner boundaries and reasons

**V3-46 — The scanner reads, parses, reports and may promote a marker under §4 rows 14 and 16.** It never repairs, rewrites, chooses between representations, imports, deletes, clears residue, or re-runs a writer operation. Its only writes are marker promotions.

**V3-47 — Marker promotion is marker resolution, not content settlement.** Promoting `pending→committed` records what the evidence supports; it never completes an interrupted writer operation, and the application remains responsible for that (§3 V3-22, V3-29). *Resolves CA-06.*

**V3-48 — Absence of a marker is not absence of a conversation.** Creation before the first marker write is Healthy `empty`; content before the first marker write is Pending `unattributed`.

**V3-49 — Reason vocabulary.**

| Reason | Meaning |
|---|---|
| `empty` | valid, void, nothing committed |
| `closed` | generation intentionally retracted and settled |
| `awaiting_content` | `pending` with nothing written yet |
| `unattributed` | content whose generation cannot be established |
| `promotion_requires_gate` | content in G≥2 without an enforced gate |
| `marker_stale` | marker recovered, not current |
| `marker_unreadable` | marker document corrupt with no valid pre-image |
| `inconsistent` | representations disagree without a permitted prefix relation |
| `malformed` | representation fails its format |
| `truncated` | representation ends mid-record — real loss |
| `absent` | expected content is missing |
| `identity_lost` | `meta.json` absent |
| `unreadable` | presence could not be established |
| `retracted-residue` | **reserved** — requires per-record generation, which does not exist; not derivable today and therefore never emitted |

**V3-50 — `retracted-residue` is retained as vocabulary, not as behaviour.** Emitting it would assert lineage the system cannot establish. *Resolves CA-04.*

**V3-51 — Incomplete enumeration is an error**, never a partial inference.

**V3-52 — Wording.** No string renders Damaged, Missing, Unavailable or `inconsistent` as a deletion confirmation.

---

## §8 — Test matrix and its mechanical check

**V3-53 — Every state-producing test names the §4 row it exercises.** The test matrix is keyed by row number; a test that cannot name its row is not a state test.

**V3-54 — Mechanical consistency check.** A test-metadata assertion walks the suite and requires that **every row 1–20 is covered by at least one test that names it**, and that every state-producing test names a row that exists. A rule change that leaves a row unreferenced fails the check. This is the mechanism that prevents a test from preserving an obsolete expectation.

| Row | Tests |
|---|---|
| 1 | VT3-01 Missing; VT3-31 deletion leaves Missing |
| 2 | VT3-02 denied path; VT3-03 incomplete enumeration |
| 3 | VT3-04 corrupt marker, no pre-image |
| 4 | VT3-05 `meta.json` absent |
| 5 | VT3-06 `Recovered` never promotes |
| 6 | VT3-07 log-prefix-of-snapshot; VT3-08 no prefix; VT3-09 collision shape |
| 7 | VT3-10 torn only rep; VT3-11 torn log with retracted snapshot; VT3-42 lone newline is Damaged, not `Healthy empty` |
| 8 | VT3-12 `committed` + void; VT3-32 partial deletion |
| 9 | VT3-13 `pending` + nothing |
| 10 | VT3-14 `(none,0)` + `[]` and no log |
| 11 | VT3-15 `cleared` + void |
| 12 | VT3-16 content with `none`; VT3-17 content with `Absent` marker |
| 13 | VT3-18 content under `cleared`; VT3-19 populated snapshot, void log |
| 14 | VT3-20 `pending,1` promotes; VT3-21 duplicate promotes |
| 15 | VT3-22 `committed,1` |
| 16 | VT3-23 `pending,G≥2` with gate |
| 17 | VT3-24 `committed,G≥2` with gate |
| 18 | VT3-25 `pending,G≥2` without gate |
| 19 | VT3-26 `committed,G≥2` without gate |
| 20 | VT3-27 orphan reported, never adopted or deleted |

Transition tests: T1 VT3-28 · T2 VT3-29 · T3 VT3-30 · T4 VT3-33 · T5 VT3-34 · T6 VT3-35 · T7 VT3-36 · T8 VT3-31/32.

**V3-55 — Gate and store tests are not state tests.** CAS, cross-process serialization, idempotent clear retry and marker crash behaviour are tested against §3 and §6, and name their clause, not a row: VT3-37 CAS refuses a stale pair · VT3-38 second instance reports `Busy` · VT3-39 clear retry is idempotent with no append between attempts · VT3-40 first marker write death yields `Absent` with no content · VT3-41 scanner writes nothing but a promotion.

**V3-56 — No state is derived from `LoadChat`'s count** — an instrumented assertion, not a convention.

---

## §9 — Standing, ownership, and what is not authorized

**V3-57 — Carried forward unchanged and unimplemented:** the mutation gate (specified, absent); CAS in one lock hold (design requirement); cross-process serialization (specified); per-representation generation stamps (**Unestablished** — the only option independent of the entry-point inventory); `retracted-residue` (vocabulary only).

**V3-58 — Separately owned, never absorbed:** creation identity collision (row 550), loader precedence (OQ-3c-a), `RewriteHistory` resurrection (OQ-3c-b).

**V3-59 — Open:** OQ-1 network roots · OQ-2 case sensitivity and junctions · OQ-4 repair path · OQ-5 manifest versioning · V2-OQ-2 `cleared` in the marker vs a separate retraction record · the row 550 collision's effect on rows 6 and 12.

**V3-60 — Nothing is authorized.** No code, no gate, no CAS, no marker write, no duplication change, no deletion path, no migration, no recovery, no persistence change. V3 is a draft for review; row 558's proposed clauses are consolidated here as **proposals**, and each is marked by where it landed, not adopted by inclusion.

**V3-61 — Row 547 and rows 554–558 remain byte-identical.** They are historical evidence, cited for rationale, never required to execute §4.
