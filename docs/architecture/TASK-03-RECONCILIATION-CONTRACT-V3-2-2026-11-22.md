# TASK-03 — reconciliation contract v3.2 (baseline outcomes incorporated)

| Field | Value |
|---|---|
| Status | **Contract v3.2 draft, revised once. Six baseline outcomes incorporated; index membership and the sweep's coverage limit made explicit. Nothing adopted. No implementation authorized.** |
| Supersedes | v3.1 (row 560) wholly; v3, v2 and amendments 555–557 remain historical evidence, unedited |
| Consolidates | rows 548–582, with rows 547–582 retained as history |
| Governing principle | **One authoritative ordered decision procedure. No state may depend on a condition hidden in an amendment, footnote, transition narrative, or test expectation.** |
| Recovery token | `TASK-03-RECONCILIATION-CONTRACT-V3-2-2026-11-22` |

Prefixes: `V3-` normative, `VT3-` acceptance test. Labels: **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Superseded**.

**V3-01 — This document is self-contained.** §4 is the executable core; every condition that decides a state appears in §4's rows. Rows 547–582 are evidence and rationale, never normative text.

**V3-02 — What changed from v3.1.** Six baseline outcomes, all inline: root observation and root class (row 577); four identities and the one-candidate-per-physical-directory rule (row 578); acceptance of an unterminated complete record plus the writer rule (row 579); the creation-collision limitation and separation of collision safety from crash-safe creation (row 580); the meaning of `Healthy` and the assessed scope (row 581); and the memory non-implication and separate disclosure (row 582). Clause numbers 01–75 are stable except where marked; inserted clauses use letter suffixes. §11 is the complete audit, including the added check that no clause treats the current memory implementation as the intended future architecture. One revision followed the owner's review: the index becomes a defined input with a measured producer gap, the verdict's coverage limit is stated, and an unparseable `meta.json` is classified (`V3-88`–`V3-92`). It adds no state.

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

**V3-07 — The two representations.** `history.ndjson` (log, append) and `history.json` (snapshot, a peer written by save, rewrite and duplication — Source-derived, rows 549/553), beside `meta.json` and the optional companions. Companions are outside the assessed scope (V3-13a); the conversation companion `memories.json` is outside the verdict entirely (V3-13b).

**V3-08 — Structural validity.**

| Format | Valid iff |
|---|---|
| `history.ndjson` | every non-blank line parses as a record; the final record may be terminated or unterminated |
| `history.json` | the whole document parses as an array of records |
| either, zero bytes | **valid and empty** — vacuous validity, no trailing newline required |

**V3-09 — Empty is defined.** Empty log = 0 bytes; empty snapshot = `[]` (2 bytes). A lone `"\n"` (1 byte) is neither valid nor empty, and is reported `malformed` (Observed, rows 548/551).

**V3-08a — The tail is evidence, not validity.** `tail = terminated | unterminated`, where `unterminated` means the last non-blank record is not followed by `\n`. The reader returns the final record at EOF (Observed, row 579, probe `unterm`). A final unterminated line that parses is a complete record and does not change the state; a final line that does not parse is `truncated` (row 7). The tail never changes which representation is authoritative. **Writer rule:** no append may concatenate onto an unterminated tail; under the mutation gate it either terminates the tail (a single `\n`, changing no content) inside the same critical section or refuses with a typed result, and never appends directly.

**V3-08b — The lone newline is still invalid.** `V3-09` is retained unchanged: a 1-byte `"\n"` is neither valid nor empty and remains `malformed`. Without that retention the relaxed rule would make it vacuously valid.

**V3-10 — Void** means no representation holds ≥1 record.

**V3-11 — Usable content** means a structurally valid representation holding ≥1 record, such that every other content-bearing representation stands to it in an admissible prefix relation (V3-45). A structurally invalid companion never contributes and never blocks; it is reported (V3-46). Where no usable content exists and something present is invalid, row 7 decides.

**V3-12 — Attribution** is the question *"can the generation of this content be established?"* It is established only when all three hold: the marker read is **current**; the generation is one the marker authorizes (§2); and the transition that opened that generation is **settled** (§3). If any fails, the content is **unattributed**. Attribution is never inferred from record content, IDs, counts, sizes or timestamps.

**V3-13 — `Inconsistent` is a reason, not a state.** It maps to Unavailable and fails closed.

**V3-13a — What `Healthy` establishes, and the assessed scope.** `Healthy` establishes observed content availability and structural consistency within the **assessed conversation scope**: the candidate directory's two history representations and its expectation marker. `meta.json` supplies the logical-ID label; candidacy is the root-and-index rule of V3-13d, and a missing `meta.json` on a candidate is row 4, not disqualification. Companions are excluded and are never assessed. The verdict means: content is present, structurally valid and mutually consistent within the permitted prefix relation, and attributable to the generation the marker authorizes; or the conversation is legitimately empty or closed under its marker; and no row indicates loss. It establishes **availability only**. It does **not** certify authorship, exclusive creation, identity uniqueness, correct write targeting, or companion-data integrity. **Nor does it certify the conversation set.** A verdict attaches to a candidate, and the candidate set is the membership rule of V3-13d: a clean sweep means *every candidate is available*, never *no other conversation data exists*, and never *nothing is here that the index does not name*. This applies to every `Healthy` row (10, 11, 14–17), not only generation 1.

**V3-13b — Memories are outside the verdict.** Conversation health never implies memory health, memory presence, memory completeness, correct memory association, or that any memory data was assessed. The current `MemoryRecord` store, both `memories.json` files (the conversation companion under `<rootDir>/chats/<id>/` and the shadow store under `StateRoot()/Chats/<id>/`), and the lorebook-trigger capture mechanism are evidence about today's data boundaries; they are **not** the intended memory architecture, data model, lifecycle, retrieval or recovery design. A separate read-only disclosure may report observations at a named scope — `unassessed` · `not_found_at_scope` · `unreadable` · `observed` · `orphaned` — and it feeds no row, adds no state, performs no write, repair, restore or deletion. `not_found_at_scope` is not deletion; `observed` is not completeness.

**V3-13c — Root observation and root class.** Before §4 is evaluated, the scan observes the root itself and records exactly one `root_state` — `observed` (listing succeeds and enumeration completes), `absent` (**only** on a `local-fixed` root, with errno 2 or 3), or `unobservable` (every other failure, every `unknown` class, every incomplete enumeration). It also records exactly one `root_class` — `local-fixed` (drive letter and `GetDriveTypeW = 3`), `non-local` (UNC volume including `\\?\UNC\`, or any other drive type), or `unknown` (no volume or undeterminable). `unknown` is treated as `non-local`. The class is a property of the resolved physical volume, not of the path string. `Missing` requires both the root observation and the error class; a sentinel such as `os.IsNotExist` is never sufficient (Observed: errno 53 maps to `ErrNotExist`).

**V3-13d — Four identities, and one candidate per physical directory.** **Locator** is a path string: it opens bytes and is never identity, a dedup key or a lock key. **Physical identity** is `(volume serial, file index)`: it survives rename, respelling, junctions and drive-letter remap, and changes on copy or restore; it is usable only when both values are nonzero. If two distinct resolved paths yield the same physical identity, the scan does **not** deduplicate, reports both, and treats the identity as unreliable. **Logical identity** is `meta.ID`, else the directory name: it attaches verdicts, ownership and retry binding, and a manual copy wrongly inherits it. **Operation identity** is row 574's durable key: retry only. **The index, defined.** The index is a membership record supplied to the scan. Each entry names exactly one conversation by **logical identity** — `meta.ID` where readable, else the directory name. Entry identity is never a path string, a display name, a companion name, a physical identity or a timestamp. **Membership is byte-equality** between a candidate's logical identity and an entry's identifier: no case folding, no path normalization, no prefix or fuzzy match. The scan records exactly one `index_state` — `present` (a record was supplied, it parses, and its enumeration completes), `absent` (none supplied), or `unreadable` (supplied but unparseable, or its enumeration fails). `absent` and `unreadable` are never treated as an empty index: a scan that has no index cannot claim a conversation is unindexed. Duplicate entries naming one identifier still establish membership and are reported; an entry naming nothing present is not a verdict in itself (row 1 when the root is absent, nothing otherwise).

**Candidacy.** A scan candidate is one physical conversation directory: a child of the conversation root that either contains `meta.json` **or** whose logical identity is an index member. **Presence, never parseability**, satisfies the first test. A directory that is neither is a non-candidate and gets no verdict; a candidate whose `meta.json` is absent, or present but unparseable, is row 4 `identity_lost`. Candidacy is therefore a **defined membership rule, not a sweep of everything on disk**. Candidates are deduplicated by physical identity when computable, never by path string and never by logical ID. Two physical directories claiming one logical ID are both reported with their own evidence; neither is preferred, merged, renamed, split or deleted. The lock key is physical identity when computable; an uncomputable identity is a typed refusal, never a path fallback.

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

**V3-30a — The lock key is physical identity when computable.** A path-string key lets two aliases obtain two locks; a logical-ID key lets a manual copy share one lock across two directories. The key is `(volume, index)`, verified before and after acquisition. An uncomputable or colliding identity is a typed refusal, never a fallback to the path. **Design requirement, unimplemented.**

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

**V3-39a — Candidate and root preconditions precede the rows.** Rows are evaluated per candidate, one candidate per physical directory (V3-13d). Root observation (V3-13c) precedes row 1: an unobservable root routes to row 2 before any conversation-level claim, and a non-local or unknown root can never reach row 1. These preconditions are part of the rows below, not footnotes. **Index membership and `index_state` are the second precondition:** rows 1 and 20 each require `index_state = present`, so an absent or unreadable index can neither report `Missing` nor `Orphaned`. Both remain fail-closed silences, not defaults — no row says "index missing therefore everything is fine".

**One stated exception:** rows 16–17 are unreachable while `gate_enforced` is false or no enforcement record exists. That unreachability is the fail-closed outcome required by V3-37, not a defect — an input that cannot exist is a defect; a capability that has not been built is a state. The distinction is part of the table, not an excuse for it.

| # | Condition (all must hold) | State | Reason |
|---|---|---|---|
| 1 | `index_state = present` **and** on a `local-fixed` root, the conversation is not found and either (`root_state = observed` and the child observation fails with errno 2 or 3) or (`root_state = absent` and the index names this conversation) | **Missing** | `absent` |
| 2 | root `unobservable`, enumeration incomplete, path denied, or a transport/access/timeout error, on any root class | **Unavailable** | `unreadable` |
| 2b | root `observed` on a `non-local` or `unknown` root, and the child observation fails with an absence-shaped code 2 or 3 | **Unavailable** | `absence_unverified` |
| 3 | marker document present but unreadable (corrupt, no matching pre-image) | **Unavailable** | `marker_unreadable` |
| 4 | `meta.json` absent, or present but unparseable, in a present directory | **Damaged** | `identity_lost` |
| 5 | marker read is `Recovered` (not current) | **Pending reconciliation** | `marker_stale` |
| 6 | two representations valid, both content-bearing, and neither is an exact prefix of the other | **Unavailable** | `inconsistent` |
| 7 | no usable content exists, and at least one present representation is structurally invalid | **Damaged** | `malformed`, or `truncated` only when the final line does not parse |
| 8 | no usable content, every present representation valid and empty or absent, marker `committed` | **Damaged** | `absent` |
| 9 | no usable content, representations valid and empty or absent, marker `pending` | **Pending reconciliation** | `awaiting_content` |
| 10 | no usable content, representations valid and empty or absent, marker `none` or `Absent` | **Healthy** | `empty` |
| 11 | no usable content, representations valid and empty or absent, marker `cleared` | **Healthy** | `closed` |
| 12 | usable content **and** marker is `Absent` or `none` | **Pending reconciliation** | `unattributed` |
| 13 | usable content **and** marker is `cleared` | **Pending reconciliation** | `unattributed` |
| 14 | usable content, marker `pending,1` | **Healthy**, promote `(committed,1)` | `unterminated_tail` when the tail is unterminated; otherwise — |
| 15 | usable content, marker `committed,1` | **Healthy** | `unterminated_tail` when the tail is unterminated; otherwise — |
| 16 | usable content, marker `pending,G≥2`, `gate_enforced`, content cannot predate enforcement | **Healthy**, promote `(committed,G)` | `unterminated_tail` when the tail is unterminated; otherwise — |
| 17 | usable content, marker `committed,G≥2`, `gate_enforced`, content cannot predate enforcement | **Healthy** | `unterminated_tail` when the tail is unterminated; otherwise — |
| 18 | usable content, marker `pending,G≥2`, otherwise | **Pending reconciliation** | `attribution_unestablished` |
| 19 | usable content, marker `committed,G≥2`, otherwise | **Pending reconciliation** | `attribution_unestablished` |
| 20 | `index_state = present` **and** directory present, otherwise decided above as Healthy, with no index entry | **Orphaned** | — |

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
| Existing **generation 1** conversation | **Healthy** under the original evidence rules (rows 14–15). No preceding generation can produce a mixed-generation log, so generation contamination is impossible. This establishes availability only; it does not certify authorship, exclusive creation, identity uniqueness, correct write targeting or companion integrity (V3-13a). Unchanged by the missing gate. |
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

**V3-55a — Memory observation is outside the scan.** The scanner does not read `StateRoot()/Chats`. If the product reports memory association, it is a separate read-only disclosure at a named scope (V3-13b), never a row input, never a conversation state, and never a claim of memory health.

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
| `identity_lost` | `meta.json` absent or unparseable — the conversation's own identity document cannot be read |
| `unreadable` | presence could not be established |
| `absence_unverified` | content absence was suggested by an absence-shaped error on a root whose continuity cannot support the claim |
| `unterminated_tail` | the last non-blank log record is not newline-terminated; a mutation warning, not loss |
| `retracted-residue` | **reserved** — requires per-record generation; never emitted |

**V3-59 — `retracted-residue` is vocabulary, not behaviour.** Emitting it would assert lineage the system cannot establish.

**V3-60 — Incomplete enumeration and a failed root observation are errors**, never a partial inference. A listing that begins and fails mid-iteration is `unobservable`; it never yields a shorter list of conversations.

**V3-61 — Wording.** No string renders Damaged, Missing, Unavailable or `inconsistent` as a deletion confirmation.

---

## §8 — Test matrix and its mechanical check

**V3-62 — Every state-producing test names the §4 row it exercises.**

**V3-63 — Mechanical consistency check.** A metadata assertion requires **every row 1, 2, 2b and 3–20 covered by at least one test naming it**, and every state-producing test to name an existing row. A rule change that orphans a row fails. Rows 16–17 are covered by tests declared **inactive with capability `gate_enforced`**: the inactive set must be explicit and enumerable, and a row with no test at all still fails. Inactive does not mean skipped silently — it means *declared unreachable until the capability exists*, which is exactly what V3-37 says about the gate.

| Row | Tests |
|---|---|
| 1 | VT3-01; VT3-31; VT3-50 local-fixed root absent, index names the conversation; VT3-59 no index, no `Missing` |
| 2 | VT3-02 denied; VT3-03 incomplete enumeration; VT3-48 unobservable root |
| 2b | VT3-49 absence-shaped child failure on a non-local root |
| 3 | VT3-04 corrupt marker |
| 4 | VT3-05 `meta.json` absent; VT3-60 `meta.json` present but unparseable |
| 5 | VT3-06 `Recovered` never promotes; VT3-43 recovered marker with G≥2 content |
| 6 | VT3-07 log-prefix-of-snapshot; VT3-08 no prefix; VT3-09 collision shape; VT3-52 manual copy, two physical directories, one logical ID |
| 7 | VT3-10 torn only rep; VT3-11 torn log + retracted snapshot; VT3-42 lone newline |
| 8 | VT3-12 `committed` + void; VT3-32 partial deletion |
| 9 | VT3-13 `pending` + nothing |
| 10 | VT3-14 `(none,0)` + `[]`, no log |
| 11 | VT3-15 `cleared` + void |
| 12 | VT3-16 content with `none`; VT3-17 content with `Absent` marker |
| 13 | VT3-18 content under `cleared`; VT3-19 populated snapshot, void log |
| 14 | VT3-20 `pending,1` promotes; VT3-21 duplicate promotes; VT3-53 unterminated complete record |
| 15 | VT3-22 `committed,1`; VT3-44 generation 1 with no preceding clear is Healthy; VT3-55 Healthy scope statement |
| 16 | VT3-23 `pending,G≥2` under enforcement, no pre-enforcement content |
| 17 | VT3-24 `committed,G≥2` idem |
| 18 | VT3-25 `pending,G≥2` without gate; VT3-45 pre-enforcement content stays Pending |
| 19 | VT3-26 `committed,G≥2` without gate |
| 20 | VT3-27 orphan reported, never adopted or deleted; VT3-51 alias paths yield one candidate and one verdict; VT3-61 no index, no `Orphaned` |

Transitions: T1 VT3-28 · T2 VT3-29 · T3 VT3-30 · T4 VT3-33 · T5 VT3-34 · T6 VT3-35 · T7 VT3-36 · T8 VT3-31/32.

New state tests: **VT3-48** an unobservable root reports every indexed conversation Unavailable and zero Missing · **VT3-49** an observed non-local root with a child errno 2/3 reports `absence_unverified` · **VT3-50** a local-fixed root absent with the conversation in the index reports Missing · **VT3-51** alias paths yield one candidate and one verdict · **VT3-52** a manual copy yields two physical candidates with one logical ID, both reported, neither merged · **VT3-53** an unterminated complete record is Healthy with `unterminated_tail` · **VT3-55** every Healthy verdict states the assessed scope, the five non-claims and that coverage is the candidate set · **VT3-59** with `index_state = absent`, no conversation is reported `Missing`, including on an absent `local-fixed` root · **VT3-60** a candidate whose `meta.json` is present but unparseable is Damaged `identity_lost`, never Healthy · **VT3-61** with `index_state = absent` or `unreadable`, no directory is reported `Orphaned` · **VT3-62** a directory carrying neither `meta.json` nor an index entry gets no verdict, and the sweep reports its candidate count rather than a conversation total.

**V3-64 — Gate and store tests are not state tests**, and name their clause: VT3-37 CAS refuses a stale pair · VT3-38 second instance reports `Busy` · VT3-39 clear retry idempotent, no append between attempts · VT3-40 first marker death yields `Absent` with no content · VT3-41 scanner writes nothing but a promotion · VT3-46 `gate_enforced` is false when the self-test fails, is skipped, or is absent · VT3-47 a configuration flag or marker field cannot set `gate_enforced`.

**V3-65 — No state is derived from `LoadChat`'s count** — an instrumented assertion.

**V3-65a — Non-state checks, naming their clause:** VT3-54 append onto an unterminated tail terminates in-section or refuses, never concatenates · VT3-56 memory disclosure reports `not_found_at_scope` versus `unreadable` versus `orphaned` without a verdict · VT3-57 a creation collision changes no row and is never repaired · VT3-58 no clause treats the current memory implementation as the intended architecture. These produce no state.

---

## §9 — Final decision-table review

**V3-66 — Totality, row by row.** The conditions partition the input space on these axes: *index state and membership* (rows 1 and 20), *root observation and class* (row 1's precondition, rows 2 and 2b), *presence* (rows 1–4), *marker currency* (row 5), *representation agreement* (row 6), *usability* (row 7), *void-with-marker* (rows 8–11), *usable-with-marker* (rows 12–19), *index membership* (row 20). Every candidate input reaches exactly one row: an `unobservable` root reaches row 2; an observed non-local root with an absence-shaped child failure reaches row 2b; only a local-fixed observed/absent root reaches row 1; rows 8–11 require every present representation valid; row 7 is their complement when something is invalid; rows 12–19 require usable content, which rows 8–11 exclude; rows 14–19 are partitioned by marker state and, for G≥2, by V3-36's verdict and V3-37's predating test; row 20 applies only to what would otherwise be Healthy.

**V3-67 — Mutual exclusivity.** Rows 1, 2 and 2b are disjoint on root observation, root class and error class. Rows 7 and 8–11 are disjoint on usability. Rows 8–15 are disjoint on marker value. Rows 16–19 are disjoint on the enforcement verdict. Rows 12–19 cannot overlap 5 (which stops earlier) or 6 (which requires two content-bearing valid representations).

**V3-68 — No row reads text outside §4** except `gate_enforced`, whose provenance V3-36 fixes and whose default is false. No row cites a footnote, a transition table, or a test expectation. The transition tables in §5 are *consequences*, never inputs.

**V3-69 — Known asymmetries, stated rather than hidden.** `Missing` is deliberately unreachable on non-local and unknown roots (row 2b), so deletion detection is lost there (row 577); rows 18–19 will hold most generation ≥2 conversations for as long as enforcement is absent (V3-42); row 5 pre-empts rows 12–19, so a recovered marker reports `marker_stale` even when content is also unattributed; row 6 pre-empts the marker rows, per V3-41; a manual copy with one logical ID yields two verdicts by physical directory, never one merged verdict; with `index_state` absent or unreadable, neither `Missing` nor `Orphaned` is reportable, so a deletion and an unowned directory are both unsurfaced (no index producer exists today — `V3-89`).

---

## §10 — Standing, ownership, and what is not authorized

**V3-70 — Carried forward unchanged and unimplemented:** the mutation gate; CAS in one lock hold; cross-process serialization and the physical-identity lock key (V3-30a); root observation and class (V3-13c); physical-identity export and candidate dedup (V3-13d); the tail writer rule (V3-08a); the supplied index and `index_state` (V3-13d; **no producer exists** — `V3-89`); per-representation generation stamps (**Unestablished**); `retracted-residue` (vocabulary only).

**V3-71 — Newly opened, not invented:** **OQ-6** — should the marker record the settlement of the last clear, so V3-44 case 3 becomes an observable input? **OQ-7** — should an enforcement record exist, so V3-37's pre-enforcement content can be distinguished? Both change the marker schema and are owner decisions.

**V3-72 — Separately owned, never absorbed:** creation identity collision (row 550), loader precedence (OQ-3c-a), `RewriteHistory` resurrection (OQ-3c-b), execution-time retargeting (rows 572/575), and memory lifecycle, recovery or migration (row 582). V3 reports their shapes and never compensates.

**V3-72a — Creation is a product requirement, not a reconciliation state.** Creating a conversation must never silently adopt, merge with, overwrite or redirect writes into an existing conversation. This is a foundational product-correctness requirement. **Collision safety** (non-clock identity, exclusive claim on the final directory, no implicit creation in other writers, a commit point that does not make a partial creation look complete) and **crash-safe creation** (no interrupted creation leaving an apparently complete partial conversation) are separate guarantees. Staging and atomic publication may serve the second; they are not required by the first. The scanner reports the residue through existing rows and never repairs, renames, splits or deletes it.

**V3-73 — Still open:** OQ-4 repair path · OQ-5 manifest versioning · V2-OQ-2 `cleared` in the marker vs a separate record · OQ-6 settlement record · OQ-7 enforcement record and gate design (OQ-7a/7b) · OQ-7d-iii caller-visible outcomes policy · operation-identity adoption · **OQ-8** the index producer (no code path writes a membership record; `GetChatList` derives the list from a directory read — `V3-89`). None changes a rule of a currently reachable row: OQ-6 would add a settlement condition to rows 16–17 and OQ-7 would make them reachable, and both rows are unreachable today; OQ-8 changes only whether rows 1 and 20 can be evaluated at all.

**V3-73a — Closed at v3.2, subject to owner adoption:** OQ-1 (row 577), OQ-2 (row 578), OQ-3b-b (row 579), row 550's reporting effect (row 580), `V3-44`'s generation-1 wording and `OQ-7c-ii` as a contract question (row 581; the retargeting defect remains separately owned), and the memories limitation (row 582).

**V3-74 — Nothing is authorized.** No code, gate, CAS, marker write, duplication change, deletion path, migration, recovery or persistence change. V3.2 is a draft.

**V3-75 — Rows 547–582 remain byte-identical.** Historical evidence, cited for rationale, never required to execute §4.


---

## §11 — V3.2 audit

**V3-76 — Audit method.** The complete v3.2 text was read as a whole after the inline changes, then checked mechanically for clause numbering, reason vocabulary, row/test coverage and every mention of memory, `MemoryRecord`, `memories.json` and lorebook capture. This is a textual and structural audit; **no test was executed**. It is not a substitute for the test matrix.

**V3-77 — Clause and reference integrity.** Every clause ID is unique; inserted clauses use suffixes and do not renumber 01–75. Every reason emitted by a row appears in `V3-58`; `retracted-residue` remains reserved and is emitted nowhere. No row cites a footnote, a transition table or a test expectation (`V3-68`). Terms used by rows — `root_state`, `root_class`, `local-fixed`, `non-local`, `tail`, `unterminated`, physical and logical identity — are defined in §1.

**V3-78 — Root rule audit.** `Missing` is reachable only on a `local-fixed` root, with the root observed and a child errno 2/3, or with the root absent and the index naming the conversation. An unobservable root is row 2; an observed non-local/unknown root with a child errno 2/3 is row 2b `absence_unverified`. No clause derives `Missing` from `os.IsNotExist`, `errors.Is(...fs.ErrNotExist)`, a sentinel, or a successful stat. The observed errno-53 behaviour is the reason (row 577); the audit confirms no surviving sentence contradicts it.

**V3-79 — Identity and candidate audit.** One candidate per physical directory; dedup by physical identity when computable; never by path string or logical ID. Candidacy is the root-and-index rule of V3-13d: a directory that is neither under the conversation root nor indexed is a non-candidate, and a missing `meta.json` on a candidate is row 4 rather than disqualification. This corrects the overbroad wording of OQ-2's non-candidate clause (`ALIAS-29`) without changing row 4 or the shadow exclusion. A verdict attaches to the logical identity. Two physical directories with one logical ID are reported twice and never merged, preferred, renamed or split. A non-candidate gets no verdict. Physical identity is guarded on nonzero values; a collision makes it unreliable rather than authoritative. The lock key is physical identity, with refusal on an uncomputable identity.

**V3-80 — Tail audit.** Validity is parseability, not termination. `unterminated` is evidence and a reason on rows 14–17, never a state and never in row 7. A lone `"\n"` remains `malformed`. Row 6's prefix relation is over records, not bytes. The writer rule — terminate-in-section or refuse, never concatenate — is present in `V3-08a` and is a gate requirement, not a scanner behaviour.

**V3-81 — `Healthy` audit.** Every `Healthy` row is covered by V3-13a's meaning and non-claims: availability only, not authorship, exclusive creation, identity uniqueness, targeting or companions. `V3-44`'s generation-1 sentence now states only that generation contamination is impossible, and points to the limitation. Rows 10–11 remain `Healthy` empty/closed and carry no content claim. No row was added for creation collision; V3-72a states the product requirement separately and V3-72 keeps it out of the scanner.

**V3-82 — Memory audit, the added check.** Every normative occurrence of memory language in v3.2 was classified; the audit's own descriptions of that language are not counted. There are exactly three kinds: (a) **exclusion** — memories are outside the assessed scope, the scan and the verdict (`V3-07`, `V3-13b`, `V3-55a`); (b) **prohibition** — conversation health must not imply memory health, presence, completeness or assessment (`V3-13b`); (c) **implementation evidence** — `MemoryRecord`, both `memories.json` files and lorebook-trigger capture are named once, explicitly as today's data boundaries and **not** the intended architecture (`V3-13b`). No clause defines a memory state, memory lifecycle, memory recovery, memory retrieval, memory ownership model, memory schema or memory migration. The disclosure vocabulary is five observations of current files, not a data model. The audit therefore finds **no accidental architectural commitment to the current implementation**; the only sentence that could be misread is the one that explicitly denies it.

**V3-83 — Fail-closed audit.** No scan error becomes `Missing`; `Unavailable` never becomes `Healthy`; `Missing` is never derived from a sentinel; a scan error never becomes a shorter list; the scanner never repairs, deletes, renames, splits, imports or chooses between representations; `Unavailable → Healthy` and the other transitions in `V3-47` remain forbidden. `Healthy` may be reached on a valid generation-1 history, and its meaning is now bounded (V3-13a).

**V3-84 — Totality and reachability audit.** Rows 1, 2, 2b and 3–20 cover the candidate input space with the preconditions in V3-39a. The only unreachable rows are 16–17, and only while `gate_enforced` is false or no enforcement record exists; that is the stated fail-closed exception, not an orphan. Every other row is reachable from an enumerated input. Row 20 applies only after rows 1–19 would otherwise yield `Healthy`.

**V3-85 — Test-matrix audit.** Every row has at least one naming test; rows 16–17 are declared inactive with capability `gate_enforced`; the inactive set is explicit. New checks cover root observation, alias dedup, manual copy, tail acceptance, the tail gate rule, the `Healthy` scope statement, memory disclosure and creation-collision non-repair. No test expectation decides a state; the matrix names rows, it does not define them. `VT3-58` is the audit check this record adds: a search over the whole v3.2 text for language that treats the current memory implementation as the intended architecture.

**V3-86 — Residual limitations, stated rather than hidden.** Network file-ID stability, real-world prevalence of unterminated logs, creation-repair implementation, companion assessment and the memory experiment are **Unestablished** or separately owned. The contract does not claim the current store is the target architecture; it states what today's evidence can support. If persistence is replaced, the root rule, identity separation, tail parseability, `Healthy` non-claims, creation safety and memory non-implication still apply; the two-representation layout does not.

**V3-87 — Audit result: pass, with the stated limitations.** The six baseline outcomes are incorporated without adding a state; §4 gains the root preconditions in rows 1 and 2, row 2b, row 7's truncation wording, and the tail reason, which are the amendments those records require. No contradiction was found between the new clauses and `V3-47`, `V3-55`, `V3-60` or the six-state vocabulary. No clause treats the current lorebook automation, `MemoryRecord` or either `memories.json` as the intended memory architecture. Nothing is adopted or authorized by this audit.

**V3-88 — Revision audit method.** After the owner's review, §1, §4, §8, §9 and §10 were re-read whole and every mention of "the index", "index entry", candidacy, coverage and `meta.json` classified. Textual and structural; **no test executed**. This is a second pass over the same draft, not a new contract.

**V3-89 — Index audit, and the finding it produced.** The v3.1 text used "the index" as an undefined authority in row 1, row 20 and the candidacy rule; a reader could not tell whether membership was a path, an ID, a listing or a timestamp. `V3-13d` now defines entry identity, the byte-equality test and `index_state`, and rows 1 and 20 require `index_state = present` rather than assuming one. **Observed, and stated rather than papered over:** no membership record is written by any code path read. `GetChatList` (`internal/backend/chat/persistence.go:485`) enumerates `<rootDir>/chats` with `os.ReadDir` and `continue`s past any directory whose `meta.json` fails to load; `AppConfig` carries only `LastChatID` (row 569). The product's own list is therefore a directory listing filtered by parseability, not an index. Consequence, in the contract's own terms: in today's product neither `Missing` nor `Orphaned` is reachable, and rows 1 and 20 are **not declared unreachable** — they are evaluable whenever a supplied index exists, and a scan without one reports neither. The producer gap is opened as **OQ-8**, a capability question, not a defect of the table.

**V3-90 — Coverage audit.** `V3-13a` now states that a verdict covers the candidate set and that a clean sweep is not evidence about non-candidates, unindexed data or the completeness of the conversation set. No clause in v3.2 says or implies otherwise; the word "sweep" never appears as a completeness claim. The mechanical check is `VT3-62`: a directory that is neither a `meta.json`-bearing child nor an index member receives no verdict, and the scan reports a candidate count, never a conversation total.

**V3-91 — Row 4 audit.** Candidacy is presence-based, so a candidate whose `meta.json` is present but unparseable was previously neither rejected (presence satisfied the test) nor classified (row 4 said "absent") and could reach rows 14–15 as `Healthy` while its own identity document was unreadable. Row 4 now covers that input as `identity_lost`; V3-58 says so. This is a consequence of making candidacy precise, not a reopening of OQ-2, and it adds no state. The product's listing skips such a directory entirely — display, not verdict (`V3-06`).

**V3-92 — Revised result: pass, with the stated limitations.** The revision adds no clause to §4's order, no state and no reason. It defines an input, states a limit, and closes one previously unclassified input (unparseable `meta.json`). Clause and reference integrity re-checked after the edits: no duplicate ID, no undefined reference, both new reasons already in V3-58's vocabulary (`identity_lost`, and `index_state` is a precondition rather than a reason). The limitations of `V3-86` stand, joined by one more: **no index producer exists**, so `Missing` and `Orphaned` cannot be reported by any build of the current product until OQ-8 is decided. Nothing is adopted or authorized.
