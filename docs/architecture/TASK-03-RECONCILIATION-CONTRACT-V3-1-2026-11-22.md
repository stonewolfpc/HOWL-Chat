# TASK-03 — reconciliation contract v3.1 (consolidation, self-contained)

| Field | Value |
|---|---|
| Status | **Contract v3.1 draft. Nothing adopted. No implementation authorized.** |
| Supersedes | v3 (row 559) wholly; v2 (row 554) and amendments 555–557 remain historical evidence, unedited |
| Consolidates | rows 548–560, with rows 547–560 retained as history |
| Governing principle | **One authoritative ordered decision procedure. No state may depend on a condition hidden in an amendment, footnote, transition narrative, or test expectation.** |
| Recovery token | `TASK-03-RECONCILIATION-CONTRACT-V3-1-2026-11-22` |

Prefixes: `V3-` normative, `VT3-` acceptance test. Labels: **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Superseded**.

**V3-01 — This document is self-contained.** §4 is the executable core; every condition that decides a state appears in §4's rows. Rows 547–560 are evidence and rationale, never normative text.

**V3-02 — What changed from v3.** Three things, all inline: the two authorities are separated (§3), `gate_enforced` gained a provenance rule (§3), and the generation ≥2 fail-closed outcome is justified as an *availability* finding rather than a bookkeeping block (§4.1). Clause numbers are stable; V3's text is carried forward except where marked.

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

**V3-05 — `LoadChat` is loadability, not validity.** It is lossy by design and its count never decides a state or a promotion.

**V3-06 — Displayability is not a state.** A conversation may be *shown* while its state is not Healthy. No string, badge or affordance may present Unavailable, Damaged or Pending content as confirmed.

**V3-07 — The two representations.** `history.ndjson` (log, append) and `history.json` (snapshot, a peer written by save, rewrite and duplication — Source-derived, rows 549/553), beside `meta.json` and the optional companions.

**V3-08 — Structural validity.**

| Format | Valid iff |
|---|---|
| `history.ndjson` | every non-blank line parses as a record, and a non-empty file ends with a newline |
| `history.json` | the whole document parses as an array of records |
| either, zero bytes | **valid and empty** — vacuous validity, no trailing newline required |

**V3-09 — Empty is defined.** Empty log = 0 bytes; empty snapshot = `[]` (2 bytes). A lone `"\n"` (1 byte) is neither valid nor empty, and is reported `malformed` (Observed, rows 548/551).

**V3-10 — Void** means no representation holds ≥1 record.

**V3-11 — Usable content** means a structurally valid representation holding ≥1 record, such that every other content-bearing representation stands to it in an admissible prefix relation (V3-45). A structurally invalid companion never contributes and never blocks; it is reported (V3-46). Where no usable content exists and something present is invalid, row 7 decides.

**V3-12 — Attribution** is the question *"can the generation of this content be established?"* It is established only when all three hold: the marker read is **current**; the generation is one the marker authorizes (§2); and the transition that opened that generation is **settled** (§3). If any fails, the content is **unattributed**. Attribution is never inferred from record content, IDs, counts, sizes or timestamps.

**V3-13 — `Inconsistent` is a reason, not a state.** It maps to Unavailable and fails closed.

---

## §2 — The expectation marker

**V3-14 — Shape.** One document, three fields, persisted through `DurableStore`:

| Field | Meaning |
|---|---|
| `state` | `none` · `pending` · `committed` · `cleared` |
| `generation` | the lineage the state describes |
| `updated_at` | bookkeeping only; never evidence |

**V3-15 — The marker stores intent, not state.** The six observed states of §4 are always computed.

**V3-16 — Generation semantics.** `0→1` on the first confirmed commit; `G→G+1` only when a clear of G has been recorded **and settled** and content is committed again. A clear alone closes G and never increments. A `Recovered` marker never increments anything.

**V3-17 — The marker has no history.** It holds the current pair only. **Consequence, and it is a finding:** after the marker advances past `(cleared,G-1)`, the scanner cannot observe that the clear was previously settled. "A generation 2 conversation with a known completed clear" is therefore not a describable input — see §4.1 case 3 and OQ-6.

**V3-18 — The generation bounds attribution; it does not identify records.** No per-record generation exists (**Unestablished**).

**V3-19 — First write is not a replacement.**

| Write | Outcomes | Death at the replacement point |
|---|---|---|
| **first** marker write | `Saved`, `Absent` | **Absent** — no document, no pre-image |
| **replacement** | `Saved`, `Recovered`, `NotSaved` | previous whole document, or nothing lost |

**V3-20 — `Recovered` is not authority.** Internally consistent, possibly stale (`Current()` is false — Source-derived). It never establishes settlement and never promotes; **row 5 stops it unconditionally.**

**V3-21 — The marker neither proves completeness nor authorises cleanup.** No marker value permits deleting any file.

**V3-22 — Marker ordering.**

| Operation | Order |
|---|---|
| first append | `(pending,1)` → append → `(committed,1)` |
| subsequent append | no marker write |
| correction with survivors | no marker write; both representations rewritten |
| clear | `(cleared,G)` → rewrite both → settled |
| append after clear | `(pending,G+1)` → append → `(committed,G+1)` |
| duplication | `(pending,1)` in the **new** directory → content → `(committed,1)` |

---

## §3 — The two authorities, and the gate

**V3-23 — Availability and mutation permission are different questions, decided by different mechanisms, and neither implies the other.**

| Authority | Question | Decided by | Can it be assumed? |
|---|---|---|---|
| **Availability** | is this conversation's data still there? | §4, from observation | no |
| **Mutation permission** | may this conversation be written to? | §3's gate | never by a state |

**V3-24 — A state never grants permission.** Healthy does not mean writable; Pending does not mean blocked. §4 classifies; §3 authorises. A product that treats a Healthy state as a write permit has misread this contract.

**V3-25 — The converse is the point of §4.1.** The absence of the gate does **not** automatically make existing content unavailable — and in generation 1 it does not. Where it does, the reason is stated in §4.1, not inferred from the missing mechanism.

**V3-26 — The gate exists in the contract; it does not exist in the product.** While a transition is unresolved, no history mutation may append to or rewrite a representation belonging to it: the writer completes the authorized operation or refuses with a typed, caller-visible failure. The scanner never completes it.

**V3-27 — Refusal is surfaced, never discarded.** A caller that swallows the failure is non-conforming (Source-derived: `_ = a.performSave()` at `app-characters.go:242,257`, `app-chat.go:137,148`).

**V3-28 — Read-then-write is not the gate.** Check, permission and write happen in **one critical section**.

**V3-29 — The marker needs compare-and-swap.** Replace `(state, generation)` only if it still equals the observed pair, inside the same lock hold. `Save` cannot express it; `Load` is lock-free. **Design requirement, unimplemented.**

**V3-30 — Serialization is cross-process.** A writer that cannot acquire the section reports `Busy` and refuses; an in-process mutex is insufficient. (Supersedes row 547 RC-27.)

**V3-31 — Settled means:**

| Transition | Settled iff |
|---|---|
| first commit (0→1) | marker `(committed,1)`, or no append happened and marker is `(none,0)` |
| clear closing G | marker `(cleared,G)` **and** both representations valid **and** void |
| append after clear (G→G+1) | the closing clear is settled, so the log was void when the first G+1 record was written |

**V3-32 — `cleared` alone is not settlement.** A populated snapshot is never dismissed as residue because the log is void.

**V3-33 — Clear retry is idempotent.** An interrupted clear that replaced the log but not the snapshot is finished by re-running the same clear; the retry never treats the void log as a new conversation and never permits an append in between.

**V3-34 — Deletion is explicitly outside the gate.** It destroys both representations, cannot corrupt the pair, and is an admissible exit from an unresolved transition. Exempt by statement; **not a resolution** — nothing is promoted afterwards.

**V3-35 — The inventory rule.** Every mutation entry point is gated except the deletion exemption; a new entry point requires a gate-coverage test.

**V3-36 — `gate_enforced` is derived, never declared.** This supersedes the "declared property" wording of v3 (row 559, its own clause V3-21) — a cross-document citation made only to record the change, and never needed to execute this table:

| Rule | Text |
|---|---|
| **default** | false, in every build that has not proved otherwise |
| **source** | a positive result from an enforcement self-test bound to the build: the single write boundary exists, the CAS path is exercised, and a refusal is observed when a transition is unresolved |
| **never** | a marker field, a configuration flag, a settings toggle, an environment variable, or a manually asserted boolean — all are writable or forgeable, and a configuration mistake must not make rows 16–17 reachable |
| **on doubt** | false. A failed, skipped or unavailable self-test is false, never true |
| **no inference from absence** | the absence of observed contamination is not evidence of enforcement |

**V3-37 — Enforcement does not certify the past.** A gate that begins operating today cannot establish that content written before it was contamination-free. Rows 16–17 therefore apply only to content that cannot predate enforcement; where a conversation's content may predate it and no enforcement record exists, the reason is `attribution_unestablished` (§4 row 18/19). **Consequence to state plainly:** shipping the gate will not retroactively clear existing generation ≥2 conversations; it stops new contamination, and existing ones remain unresolved until an enforcement record exists. Recording one is **OQ-7**, not invented here.

**V3-38 — Ownership boundary.** Creation identity collision (row 550), loader precedence (OQ-3c-a) and `RewriteHistory` resurrection (OQ-3c-b) are separately owned. V3 reports their shapes and never compensates.

---

## §4 — The decision procedure

**V3-39 — Applied in order; the first matching row decides.** Rows are mutually exclusive by construction; a row that cannot be reached is a defect of this table. The only input that is not an observation is `gate_enforced`, whose provenance is fixed by V3-36.

**One stated exception:** rows 16–17 are unreachable while `gate_enforced` is false or no enforcement record exists. That unreachability is the fail-closed outcome required by V3-37, not a defect — an input that cannot exist is a defect; a capability that has not been built is a state. The distinction is part of the table, not an excuse for it.

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
| 16 | usable content, marker `pending,G≥2`, `gate_enforced`, content cannot predate enforcement | **Healthy**, promote `(committed,G)` | — |
| 17 | usable content, marker `committed,G≥2`, `gate_enforced`, content cannot predate enforcement | **Healthy** | — |
| 18 | usable content, marker `pending,G≥2`, otherwise | **Pending reconciliation** | `attribution_unestablished` |
| 19 | usable content, marker `committed,G≥2`, otherwise | **Pending reconciliation** | `attribution_unestablished` |
| 20 | directory present, otherwise decided above as Healthy, with no index entry | **Orphaned** | — |

**V3-40 — Six states, all derived.** Healthy · Pending reconciliation · Missing · Damaged · Unavailable · Orphaned. `inconsistent` is a reason on Unavailable.

**V3-41 — Row 6 deliberately precedes rows 12–19.** Where two valid representations disagree without a permitted prefix relation, the disagreement is the more specific finding, and the marker cannot resolve it. Both outcomes fail closed; neither becomes Healthy without new evidence.

### 4.1 — Why generation ≥2 content is not Healthy

**V3-42 — The justification, and it is about availability, not write permission.** Contamination and legitimate lag are byte-identical:

- A legitimate generation ≥2 history is: clear settled (both representations void) → append → log holds G+1, snapshot void or a prefix.
- A contaminated generation ≥2 history is: clear interrupted before the rewrite → append into a log that still holds G → log holds G followed by G+1.

Both yield a log whose snapshot is a prefix, both parse as records, both are indistinguishable in structure, size, ordering and content. Nothing in the filesystem marks the boundary (V3-18), and the marker's current pair is the same in both cases (V3-17). **Therefore availability itself is unestablished**, and the fail-closed state is Pending — not because mutation is forbidden, but because the scanner cannot tell the user whether the retracted messages they asked to delete are still present.

**V3-43 — The harm is not hypothetical.** A contaminated log displays retracted content as current (Observed: probe 12's log returned four messages after a clear). That is the resurrection class owned by OQ-3c-b, reached through *reading*, not writing. Classifying such a conversation Healthy would present deleted content as confirmed.

**V3-44 — The five cases, answered.**

| Case | Answer |
|---|---|
| Existing **generation 1** conversation | **Healthy** under the original evidence rules (rows 14–15). No preceding generation exists, so nothing can contaminate it. Unchanged by the missing gate. |
| Existing **generation ≥2**, no unresolved transition | **Pending**, `attribution_unestablished`. Pending because availability cannot be established (V3-42), not merely because the gate is absent. |
| Generation ≥2 with a **known completed clear** | **Not a describable input.** The marker has no history (V3-17): once it advanced to `G`, the earlier settlement is unobservable. The scanner cannot receive "known". Making this input observable requires an enforcement or settlement record — **OQ-6**. |
| Generation ≥2 **after a successful append** | **Displayable, state Pending.** Display is allowed and lossy (V3-06); no state grants write permission either way (V3-24), so nothing is gained by calling it Healthy. |
| Generation ≥2 with a **recovered marker** | **Row 5**, Pending `marker_stale`, unchanged. Row 5 precedes rows 18–19 and is the stricter finding. |

**V3-45 — Coexistence direction** (Source-derived: `AppendHistoryEntry` writes only the log; `SaveChat` and `RewriteHistory` write both): equal = coherent; snapshot an exact prefix of the log = expected lag, log is the basis, the shorter set never claimed complete; log an exact prefix of the snapshot = row 6; no prefix = row 6. Prefix means exact prefix on full record content in order (`sender`, `text`, `role`, `timestamp`, plus `hasImage`/`characterId` when present). Record IDs never establish it — they are positional indices (`len-1`, `i`).

**V3-46 — An invalid companion does not decide a state.** A malformed snapshot beside a valid content-bearing log does not make the conversation Damaged and does not prevent Healthy; the companion is reported `malformed`. The invalid representation is never the basis of a promotion. **Row 7 applies only when no usable content exists at all** — a lone `"\n"`, a torn-only log, or two invalid representations.

**V3-47 — Forbidden transitions.** Unavailable→Healthy; Unavailable→Missing; Damaged→Missing; Damaged→Healthy; Orphaned→Healthy; `inconsistent`→Healthy; any automatic repair.

**V3-48 — Row 8 versus row 12.** `committed` with no usable content is Damaged; `none`/`Absent` with no content is Healthy-empty. The marker separates "content was committed and is gone" from "nothing was ever committed".

**V3-49 — Row 12 is the fail-closed answer to duplication and to bypassing writers.** A conforming duplicate writes its own marker (V3-22), so it never reaches row 12. A non-conforming duplicate — and probe 12's appended mixture — does, and is Pending, never Healthy.

**V3-50 — Row 20 is evaluated only after rows 1–19 would otherwise yield Healthy.** Orphaning changes ownership, not availability.

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
| record appended, gate enforced, no pre-enforcement content | intact | 16 | Healthy, promoted `(committed,G+1)` |
| record appended, otherwise | intact | 18 | Pending `attribution_unestablished` |
| record appended, otherwise | old snapshot still populated | 18, then 6 if both content-bearing and non-prefix | Pending; `inconsistent` reported when it applies |
| append torn | — | 7 | Damaged `truncated`; retracted content never stands in for the new record |

### T7 — Duplication

| Step | Crash after | Rows | State |
|---|---|---|---|
| new directory + `(pending,1)` | — | 9 | Pending `awaiting_content` |
| snapshot written (log absent) | — | 14 | Healthy, promoted `(committed,1)` — one representation suffices |
| `(committed,1)` | — | 15 | Healthy |
| source identity collided (row 550 shape) | — | 6 or 12 | reported as collision; never repaired, renamed or split |

### T8 — Deletion

| Step | Crash after | Rows | State |
|---|---|---|---|
| deletion requested during an unresolved transition | permitted (V3-34) | — | no marker write, no promotion |
| directory removed | — | 1 | Missing |
| removal partially completed | directory present, content gone | 8 | Damaged `absent` |

---

## §6 — Crash points and store behaviour

**V3-51 — Marker crash outcomes** (V3-19): first write → `Saved` or `Absent`; replacement → `Saved`, `Recovered`, or previous document intact. Never a partial document, never a third state. Source-derived from `TestCrashBeforeReplaceLeavesPreviousDocument`, **not re-executed**.

**V3-52 — `Recovered` reaches row 5 and stops.** It never promotes, never satisfies settlement, never raises a state.

**V3-53 — Content crash points never route through the marker.** A torn append, an interrupted rewrite and an aborted clear are decided by rows 6, 7 and 13 from observed bytes.

**V3-54 — `Absent` marker with content is row 12, not row 10.**

---

## §7 — Scanner boundaries and reasons

**V3-55 — The scanner reads, parses, reports and may promote a marker under rows 14 and 16.** It never repairs, rewrites, chooses between representations, imports, deletes, clears residue, or re-runs a writer operation.

**V3-56 — Marker promotion is marker resolution, not content settlement.**

**V3-57 — Absence of a marker is not absence of a conversation.** Creation before the first marker write is Healthy `empty`; content before it is Pending `unattributed`.

**V3-58 — Reason vocabulary.**

| Reason | Meaning |
|---|---|
| `empty` | valid, void, nothing committed |
| `closed` | generation retracted and settled |
| `awaiting_content` | `pending` with nothing written |
| `unattributed` | content whose generation cannot be established |
| `attribution_unestablished` | content in G≥2 whose purity cannot be established (V3-42) |
| `marker_stale` | marker recovered, not current |
| `marker_unreadable` | marker corrupt with no valid pre-image |
| `inconsistent` | representations disagree without a permitted prefix relation |
| `malformed` | representation fails its format |
| `truncated` | representation ends mid-record — real loss |
| `absent` | expected content missing |
| `identity_lost` | `meta.json` absent |
| `unreadable` | presence could not be established |
| `retracted-residue` | **reserved** — requires per-record generation; never emitted |

**V3-59 — `retracted-residue` is vocabulary, not behaviour.** Emitting it would assert lineage the system cannot establish.

**V3-60 — Incomplete enumeration is an error**, never a partial inference.

**V3-61 — Wording.** No string renders Damaged, Missing, Unavailable or `inconsistent` as a deletion confirmation.

---

## §8 — Test matrix and its mechanical check

**V3-62 — Every state-producing test names the §4 row it exercises.**

**V3-63 — Mechanical consistency check.** A metadata assertion requires **every row 1–20 covered by at least one test naming it**, and every state-producing test to name an existing row. A rule change that orphans a row fails. Rows 16–17 are covered by tests declared **inactive with capability `gate_enforced`**: the inactive set must be explicit and enumerable, and a row with no test at all still fails. Inactive does not mean skipped silently — it means *declared unreachable until the capability exists*, which is exactly what V3-37 says about the gate.

| Row | Tests |
|---|---|
| 1 | VT3-01; VT3-31 |
| 2 | VT3-02 denied; VT3-03 incomplete enumeration |
| 3 | VT3-04 corrupt marker |
| 4 | VT3-05 `meta.json` absent |
| 5 | VT3-06 `Recovered` never promotes; VT3-43 recovered marker with G≥2 content |
| 6 | VT3-07 log-prefix-of-snapshot; VT3-08 no prefix; VT3-09 collision shape |
| 7 | VT3-10 torn only rep; VT3-11 torn log + retracted snapshot; VT3-42 lone newline |
| 8 | VT3-12 `committed` + void; VT3-32 partial deletion |
| 9 | VT3-13 `pending` + nothing |
| 10 | VT3-14 `(none,0)` + `[]`, no log |
| 11 | VT3-15 `cleared` + void |
| 12 | VT3-16 content with `none`; VT3-17 content with `Absent` marker |
| 13 | VT3-18 content under `cleared`; VT3-19 populated snapshot, void log |
| 14 | VT3-20 `pending,1` promotes; VT3-21 duplicate promotes |
| 15 | VT3-22 `committed,1`; VT3-44 generation 1 with no preceding clear is Healthy |
| 16 | VT3-23 `pending,G≥2` under enforcement, no pre-enforcement content |
| 17 | VT3-24 `committed,G≥2` idem |
| 18 | VT3-25 `pending,G≥2` without gate; VT3-45 pre-enforcement content stays Pending |
| 19 | VT3-26 `committed,G≥2` without gate |
| 20 | VT3-27 orphan reported, never adopted or deleted |

Transitions: T1 VT3-28 · T2 VT3-29 · T3 VT3-30 · T4 VT3-33 · T5 VT3-34 · T6 VT3-35 · T7 VT3-36 · T8 VT3-31/32.

**V3-64 — Gate and store tests are not state tests**, and name their clause: VT3-37 CAS refuses a stale pair · VT3-38 second instance reports `Busy` · VT3-39 clear retry idempotent, no append between attempts · VT3-40 first marker death yields `Absent` with no content · VT3-41 scanner writes nothing but a promotion · VT3-46 `gate_enforced` is false when the self-test fails, is skipped, or is absent · VT3-47 a configuration flag or marker field cannot set `gate_enforced`.

**V3-65 — No state is derived from `LoadChat`'s count** — an instrumented assertion.

---

## §9 — Final decision-table review

**V3-66 — Totality, row by row.** The conditions partition the input space on three axes: *presence* (rows 1–4), *marker currency* (row 5), *representation agreement* (row 6), *usability* (row 7), *void-with-marker* (rows 8–11), *usable-with-marker* (rows 12–19), *index membership* (row 20). Every input reaches exactly one row: rows 8–11 require every present representation valid; row 7 is their complement when something is invalid; rows 12–19 require usable content, which rows 8–11 exclude; rows 14–19 are partitioned by marker state and, for G≥2, by V3-36's verdict and V3-37's predating test; row 20 applies only to what would otherwise be Healthy.

**V3-67 — Mutual exclusivity.** Rows 7 and 8–11 are disjoint on usability. Rows 8–15 are disjoint on marker value. Rows 16–19 are disjoint on the enforcement verdict. Rows 12–19 cannot overlap 5 (which stops earlier) or 6 (which requires two content-bearing valid representations).

**V3-68 — No row reads text outside §4** except `gate_enforced`, whose provenance V3-36 fixes and whose default is false. No row cites a footnote, a transition table, or a test expectation. The transition tables in §5 are *consequences*, never inputs.

**V3-69 — Known asymmetries, stated rather than hidden.** Rows 18–19 will hold most generation ≥2 conversations for as long as enforcement is absent (V3-42); row 5 pre-empts rows 12–19, so a recovered marker reports `marker_stale` even when content is also unattributed; row 6 pre-empts the marker rows, per V3-41.

---

## §10 — Standing, ownership, and what is not authorized

**V3-70 — Carried forward unchanged and unimplemented:** the mutation gate; CAS in one lock hold; cross-process serialization; per-representation generation stamps (**Unestablished**); `retracted-residue` (vocabulary only).

**V3-71 — Newly opened, not invented:** **OQ-6** — should the marker record the settlement of the last clear, so V3-44 case 3 becomes an observable input? **OQ-7** — should an enforcement record exist, so V3-37's pre-enforcement content can be distinguished? Both change the marker schema and are owner decisions.

**V3-72 — Separately owned, never absorbed:** creation identity collision (row 550), loader precedence (OQ-3c-a), `RewriteHistory` resurrection (OQ-3c-b).

**V3-73 — Still open:** OQ-1 network roots · OQ-2 case sensitivity and junctions · OQ-4 repair path · OQ-5 manifest versioning · V2-OQ-2 `cleared` in the marker vs a separate record · the row 550 collision's effect on rows 6 and 12.

**V3-74 — Nothing is authorized.** No code, gate, CAS, marker write, duplication change, deletion path, migration, recovery or persistence change. V3.1 is a draft.

**V3-75 — Rows 547–560 remain byte-identical.** Historical evidence, cited for rationale, never required to execute §4.
