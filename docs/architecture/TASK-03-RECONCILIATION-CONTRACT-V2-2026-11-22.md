# TASK-03 — reconciliation contract v2 (consolidation)

| Field | Value |
|---|---|
| Status | **Contract v2 draft. Nothing adopted. No implementation authorized.** |
| Supersedes | `TASK-03-RECONCILIATION-CONTRACT-2026-11-22.md` (row 547) — retained unedited as the historical contract |
| Consolidates | rows 548, 549, 550, 551, 552, 553 |
| Excludes | row 550's creation-identity defect, OQ-3c-a, OQ-3c-b — separately owned, see §11 |
| Recovery token | `TASK-03-RECONCILIATION-CONTRACT-V2-2026-11-22` |

Prefixes: `V2-` normative, `VT-` acceptance test. Labels: **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Superseded**.

This is a consolidation, not an addendum: a clause here replaces its predecessors rather than stacking on them.

---

## 1. Scope

**V2-01 — The question.** Whether a conversation's data is still *available*. Not who deleted it, not that it is unrecoverable, not that the filesystem cooperated.

**V2-02 — What is not claimed.** No authorship, no destruction proof, no deletion receipt, no physical identity. The retired 528–546 series remains historical evidence; rows 528–553 remain byte-identical.

**V2-03 — The scanner's authority is observation.** It reads, parses and reports. It never repairs, rewrites, chooses a winner between representations, imports, deletes, or compensates for a writer defect. The one write it may perform is the expectation record (§4).

---

## 2. The four properties

**V2-04 — Every state-determining clause names its basis.** These four are never interchangeable:

| Property | Question | Evidence |
|---|---|---|
| **Presence** | is the artefact there at all? | `stat` / listing |
| **Structural validity** | does the *whole* representation conform to its format? | whole-file parse |
| **Content presence** | does a structurally valid representation hold ≥1 record? | parse result |
| **Completeness** | is there evidence the committed history is still there? | the decision, from both representations + the expectation record |

**V2-05 — `LoadChat` is a loadability mechanism, not a validity test.** It is safe to call (Observed: repeated loads leave a conversation directory byte-identical) and it is lossy by design: truncated, wrong-shaped and interior-corrupt histories all return fewer messages with no error (Observed, row 551 §3). Its count may never decide a state or a promotion.

**V2-06 — A state may never mean "the app can display whatever survived."** If displayability is all the evidence supports, the state is not Healthy.

---

## 3. Storage and representations

**V2-07 — Layout (Source-derived).**

```
<rootDir>/chats/<chatID>/
    meta.json          required — the conversation's identity
    history.ndjson     accepted representation — append log
    history.json       accepted representation — snapshot (NOT legacy)
    world.json         optional — written only when non-nil
    scenario.json      optional — written only when non-nil
    characters.json    written on save
    lorebooks.json     written on save
    memories.json      written on save
```

**V2-08 — `history.json` is a peer, not a legacy remnant** (Superseded: row 547 RC-05). `SaveChat`, `RewriteHistory` and `SaveChatAs` all write it; a duplicated conversation may have *only* it (Observed, row 549 §2).

**V2-09 — Structural validity, per format.**

| Format | Valid iff |
|---|---|
| `history.ndjson` | every non-blank line parses as a record **and** the file ends with a newline |
| `history.json` | the whole document parses as an array of records |

**V2-10 — Empty forms.** `history.json` empty = `[]` (2 bytes). `history.ndjson` empty = 0 bytes. A lone newline (1 byte) is neither empty nor content. (Observed, row 551 §2.)

**V2-11 — An empty representation is not evidence of an empty conversation**; §5 decides what it means.

---

## 4. The expectation record

**V2-12 — Shape.** A marker holding three fields, persisted as one document through `DurableStore` (`internal/durable`: exclusive lock, temp write, manifest commit with hash and size, pre-images, `testHookBeforeReplace`):

| Field | Meaning |
|---|---|
| `state` | `none` · `pending` · `committed` · `cleared` |
| `generation` | which history lineage the state describes (§4.1) |
| `updated_at` | bookkeeping only; never evidence |

**V2-13 — The marker's crash behaviour is tested, not inherited**: after a simulated death mid-save the reloaded marker must be the previous or the new value, never a third state and never unreadable.

### 4.1 Generation — why four scalars are not enough

**V2-14 — A clear retracts one lineage and the next commit begins another.** Without a generation, a retracted snapshot and a new conversation's first records are the same evidence, and the scanner could read the retracted content as the new history or treat new content as cleared.

| `state` | Meaning |
|---|---|
| `none` | no content committed, no clear recorded — generation 0 |
| `pending` | a content-changing operation began in this generation and is not confirmed |
| `committed` | content is committed in this generation |
| `cleared` | the named generation is intentionally retracted; the next commit begins `generation + 1` |

**V2-15 — Generation changes only at these points:** 0→1 on the first confirmed commit; G→G+1 when a clear is recorded and content is committed again. A clear alone does not increment: it closes generation G.

**V2-16 — The generation bounds attribution; it does not identify records.** It tells the scanner which lineage the marker describes. It cannot say which bytes in a representation belong to which lineage — that distinction is not recorded anywhere today (**Unestablished**), and V2 must not pretend otherwise.

### 4.2 Marker writes, in order

| Operation | Order |
|---|---|
| first append | `(pending,1)` → append → `(committed,1)` |
| subsequent append | no write (state unchanged) |
| correction / retraction with survivors | no state change; content rewrites both representations |
| clear | `(cleared,G)` → rewrite → done |
| append after clear | `(pending,G+1)` → append → `(committed,G+1)` |

**V2-17 — The marker precedes the content change.** Intent is durable before content moves, so a crash leaves a *record of an attempt*, never a silent gap. This Supersedes the 549 append-then-mark ordering.

**V2-18 — `pending` can never produce a false Damaged.** `pending` with nothing present is Pending reconciliation, not loss.

---

## 5. States, derived

**V2-19 — Six states**, all computed per observation. No state is stored.

| State | Definition | UI |
|---|---|---|
| **Healthy** | the current generation's expectations are met by evidence (§7) | normal |
| **Pending reconciliation** | an operation began and is not confirmed; or a worklist entry is unresolved | normal, unobtrusive marker |
| **Missing** | the directory is absent | removed from the active list |
| **Damaged** | expected content is absent, or a representation is invalid where it should be valid | shown as damaged, never normal, never auto-repaired |
| **Unavailable** | the condition could not be established — error, denied, incomplete listing, missing manifest, or `Inconsistent` | retained, marked unresolved, retried |
| **Orphaned** | a directory with no index entry | reported; never imported, never deleted |

**V2-20 — `Inconsistent` is a reason, not a seventh state.** It maps to **Unavailable**: the condition could not be established. It fails closed and never becomes Healthy without new evidence.

**V2-21 — Forbidden transitions.** Unavailable→Healthy; Unavailable→Missing; Damaged→Missing; Damaged→Healthy; Orphaned→Healthy; `Inconsistent`→Healthy; any automatic repair.

**V2-22 — `meta.json` absent in a present directory is Damaged**, not Missing: the messages may still be on disk. (Supersedes row 547 RC-11's "all required content absent → Missing".)

**V2-23 — Reason vocabulary:** `empty` · `malformed` · `truncated` · `inconsistent` · `retracted-residue` · `absent` · `unreadable`. Reasons accompany states; they never replace them.

---

## 6. Coexistence

**V2-24 — Count is not provenance.** No clause may select a representation by record count, size, timestamp, or "the log is usually right". (Supersedes row 552 RC6-11.)

**V2-25 — Disagreement fails closed.** Two structurally valid representations with different record sets are `Inconsistent` unless the relationship is established as an **exact prefix on full record content in order** (`sender`, `text`, `role`, `timestamp`, plus `hasImage`/`characterId` when present).

**V2-26 — Record IDs cannot establish it.** They are positional indices (`AppendChatMessage`: `len-1`; `ReplaceChatHistory`: `i` — Source-derived), so deleting a message renumbers the tail.

**V2-27 — Direction is interpreted against the writer lifecycle, not assumed.** `AppendHistoryEntry` writes only the log; `SaveChat` and `RewriteHistory` write both. Therefore:

| Relationship | Interpretation |
|---|---|
| equal | coherent |
| snapshot is an exact prefix of the log | **expected lag** — coherent; the log is the basis; the shorter set is not claimed complete |
| log is an exact prefix of the snapshot | **not expected** — `Inconsistent` (torn log, or resurrected snapshot) |
| no prefix relation | `Inconsistent` |

**V2-28 — A prefix relation is evidence about order, never proof of provenance or completeness.** Reverse-prefix and non-prefix cases are never resolved by choosing.

**V2-29 — When the marker records a generation the content cannot belong to**, the content is `retracted-residue`: reported, never used to decide, never deleted (§10).

---

## 7. Decision table

Applied in order; the first matching row wins.

| # | Observation | State | Basis |
|---|---|---|---|
| 1 | directory absent | **Missing** | presence |
| 2 | listing incomplete, path denied, I/O error | **Unavailable** | presence (failed) |
| 3 | manifest absent/unreadable | **Unavailable** | presence |
| 4 | `meta.json` absent, directory present | **Damaged** (identity lost) | presence |
| 5 | marker `none`/`cleared`, all representations structurally valid and empty | **Healthy** (valid empty) | content presence = none expected |
| 6 | marker `cleared`, content present in any representation | **Pending reconciliation**, reason `retracted-residue` | completeness: clear unconfirmed |
| 7 | marker `pending`, nothing present | **Pending reconciliation** | V2-18 |
| 8 | marker `pending`/`committed`, ≥1 structurally valid representation with ≥1 record, coexistence coherent per V2-27 | **Healthy**, promoted if `pending` | structural validity + content presence + coherence |
| 9 | coexistence `Inconsistent` per V2-27 | **Unavailable**, reason `inconsistent` | completeness |
| 10 | only representation is invalid; content expected | **Damaged**, reason `malformed` or `truncated` | structural validity |
| 11 | marker `committed`, nothing present, no retraction recorded | **Damaged**, reason `absent` | completeness |
| 12 | directory found in full sweep with no index entry | **Orphaned** | presence |

**V2-30 — Row 8 is the only promotion path**, and its evidence is exactly: presence, whole-representation validity, ≥1 record, coherent coexistence. Not `LoadChat`'s count, not size, not mtime, not the larger set.

**V2-31 — Row 10's `truncated` is real loss**: the record being written is unrecoverable.

---

## 8. Transitions, crash points, and what the scanner must report

Read-only throughout. "Confirmed" means the marker write returned success.

### T1 — Creation

| Step | Crash after this step | Scanner reports |
|---|---|---|
| directory created | empty directory, no `meta.json` | **Damaged** (identity incomplete) — or an **anomaly** if it matches a collision shape (row 550) |
| `meta.json` written | `meta.json` only, marker `(none,0)` | **Healthy** — valid empty conversation |
| snapshots written (`history.json` = `[]`, companions) | all empty files, marker `(none,0)` | **Healthy** |

**V2-32 — Creation is the only transition where a missing history representation is normal**, and only while the marker is `(none,0)`.

### T2 — First append

| Step | Crash after this step | Scanner reports |
|---|---|---|
| `(pending,1)` written | pending, nothing present | **Pending reconciliation** |
| record appended to the log | pending, ≥1 valid record | **Healthy**, promoted to `(committed,1)` |
| `(committed,1)` written | committed, ≥1 valid record | **Healthy** |
| append torn mid-record | pending or committed, log invalid, no other representation | **Damaged** (`truncated`) |

### T3 — Subsequent append

| Step | Crash after this step | Scanner reports |
|---|---|---|
| record appended | log = snapshot + 1 record (snapshot is a prefix) | **Healthy** — expected lag (V2-27) |
| append torn mid-record | log invalid or ends without newline; snapshot is a valid prefix | **Damaged** (`truncated`) if the log is the only content; the snapshot is reported alongside, not used to excuse it |
| nothing written | unchanged | unchanged |

**V2-33 — The expected-lag case is why `snapshot ⊂ log` is coherent while `log ⊂ snapshot` is not.** The lifecycle predicts one direction only.

### T4 — Correction / retraction (survivors remain)

| Step | Crash after this step | Scanner reports |
|---|---|---|
| log rewritten (rename) | log = new set, snapshot = old set; old is not a prefix of new when an interior message was removed | **Unavailable** (`inconsistent`) — not resolved by choosing |
| snapshot rewritten | both = new set | **Healthy** |
| nothing written | unchanged | unchanged |

**V2-34 — The scanner does not reconstruct the intended list.** Which set is "correct" is the writer's question, and OQ-3c-b owns the resurrection risk.

### T5 — Clear

| Step | Crash after this step | Scanner reports |
|---|---|---|
| `(cleared,G)` written | marker says retracted, content still fully present | **Pending reconciliation**, reason `retracted-residue` (row 6) |
| log emptied (rename) | log 0 bytes, snapshot still populated | **Pending reconciliation**, `retracted-residue` |
| snapshot emptied | both empty, marker `(cleared,G)` | **Healthy** — valid empty conversation |

**V2-35 — A clear that recorded intent but did not finish is Pending, never Damaged and never Healthy-with-content.** The user asked for empty; the state says the request is unresolved.

### T6 — Append after clear

| Step | Crash after this step | Scanner reports |
|---|---|---|
| `(pending,G+1)` written | pending, nothing present, `(cleared,G)` history may still be on disk | **Pending reconciliation**; residue reported |
| record appended | pending, log has ≥1 valid record; old snapshot may still hold retracted content | **Healthy**, promoted to `(committed,G+1)`; the snapshot is `retracted-residue`, reported and not compared |
| `(committed,G+1)` written | committed, log ≥1 record | **Healthy** |
| append torn | pending/committed G+1, log invalid, snapshot holds retracted content | **Damaged** (`truncated`) — the retracted content must not stand in for the new record |

**V2-36 — This table is the answer to "old history vs new history".** The generation is what stops the retracted snapshot from being read as the new conversation's content; content cannot make that distinction on its own.

---

## 9. Scan, worklist, cadence, deletion

**V2-37 — Read-only, with one exception.** The scanner writes the expectation record (§4) and nothing else. No repair, no rewrite, no re-pointing, no cleanup.

**V2-38 — Bounded passes**: record cap and time budget per pass; unprocessed entries carry forward.

**V2-39 — Worklist is an optimization; the full sweep is authority.** An entry written before removal is harmless if removal never happens; a lost entry is rediscovered by the sweep.

**V2-40 — Retry with per-record backoff and a persistent attempt counter**; the counter is diagnostic and never changes a state.

**V2-41 — Cadence**: startup, periodic with jitter, resume-from-sleep, manual.

**V2-42 — Deletion scope**: exactly `<rootDir>/chats/<chatID>`, never a parent, never an uncontained assembled path. The only identity rule kept, constraining what may be deleted.

**V2-43 — Orphans are reported, never adopted or deleted.**

**V2-44 — `Incomplete enumeration is an error`**: a stopped listing yields Unavailable for that conversation, never a partial inference.

**V2-45 — Wording**: no string may render Damaged, Missing, Unavailable or `Inconsistent` as a deletion confirmation.

---

## 10. What the marker may not do

**V2-46 — `cleared` grants no cleanup authority.** It records intent and changes interpretation only. Residual files are reported; removing them is a separate, user-initiated operation outside this contract.

**V2-47 — No marker value authorises deletion of any file**, including residue that is provably retracted.

**V2-48 — The marker never becomes a proof of completeness.** It bounds attribution between generations; it cannot certify which records survive.

---

## 11. Defect boundaries

| Finding | Owner | This contract's obligation |
|---|---|---|
| Creation ID collision, non-exclusive directory creation | **row 550** | report the collision shape; never repair, rename, split |
| Loader prefers a torn log over a fuller valid snapshot | **OQ-3c-a** | do not choose the good file to hide it |
| `RewriteHistory` resurrects cleared content after a crash | **OQ-3c-b** | report `retracted-residue`; never rewrite a log to close the window |

**V2-49 — None of these is fixed, worked around, or absorbed here.** A scanner that reports them is a symptom report; a scanner that compensates for them has exceeded its authority.

---

## 12. Acceptance tests

- **VT-01** Valid empty `(none,0)` → Healthy.
- **VT-02** `(cleared,G)` + empty → Healthy.
- **VT-03** `(cleared,G)` + content → Pending reconciliation, `retracted-residue`.
- **VT-04** `(pending,G)` + nothing → Pending, never Damaged.
- **VT-05** `(pending,G)` + valid record → Healthy, promoted `(committed,G)`.
- **VT-06** snapshot exact prefix of log → Healthy, log is the basis, shorter not claimed complete.
- **VT-07** log exact prefix of snapshot → Unavailable (`inconsistent`).
- **VT-08** neither prefix → Unavailable (`inconsistent`).
- **VT-09** valid log + malformed snapshot → Healthy on the log, `malformed` companion reported; promotion never rests on the malformed file.
- **VT-10** torn log, no other representation → Damaged (`truncated`).
- **VT-11** torn log, retracted snapshot present → Damaged; retracted content does not excuse it.
- **VT-12** `(committed,G)` + nothing → Damaged (`absent`).
- **VT-13** `meta.json` absent → Damaged, not Missing.
- **VT-14** denied path → Unavailable, never Missing.
- **VT-15** marker save killed at `testHookBeforeReplace` → previous or new value only.
- **VT-16** scanner writes nothing but the marker (byte-identical directory otherwise).
- **VT-17** `cleared` authorises no deletion; residue present after a clear is reported and left.
- **VT-18** a conversation whose log is duplicated (`SaveChatAs`, snapshot only) scans Healthy.
- **VT-19** no state is derived from `LoadChat`'s count (instrumented assertion).
- **VT-20** deletion outside the conversation subtree is refused.

---

## 13. Open

| # | Question |
|---|---|
| OQ-1 | UNC/network roots, where absence may mean disconnected — likely always Unavailable |
| OQ-2 | Case sensitivity and junctions inside a conversation directory |
| OQ-4 | Whether Damaged gets a user-initiated repair path, and what repair may do |
| OQ-5 | Manifest versioning across storage changes |
| V2-OQ-1 | Whether the generation should be recorded per representation rather than per conversation (would make attribution exact; today it is not) |
| V2-OQ-2 | Whether `cleared` belongs in the marker or in a separate retraction record |
| V2-OQ-3 | Whether `retracted-residue` should ever be user-visible or only logged |

**V2-50 — No implementation authorization.** No code, no manifest, no marker, no migration, no deletion enablement follows from this document.
