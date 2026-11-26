# TASK-03 — OQ-3b-b: the missing final newline

| Field | Value |
|---|---|
| Purpose | Decide whether the final newline is part of the storage format's validity or only a writer invariant, and whether rejecting its absence misclassifies recoverable content |
| Status | **Read-only investigation. Outcome: Amend. V3.1 unedited. Nothing adopted, nothing authorized.** |
| Method | Source inspection of every `history.ndjson` writer and the loader; one disposable fixture probe exercising the real `LoadChat` against eight tail shapes (`build/tmp/oq3bbprobe`, log `build/tmp/oq3bbprobe.log`) |
| Recovery token | `TASK-03-OQ3BB-TERMINATOR-2026-11-22` |

Third of row 576's six baseline questions, after OQ-1 (row 577) and OQ-2 (row 578). Resolved under the owner's frame: **Amend**, **Clarify**, **Retain**. This record decides the terminator rule only; it does not design a repair path or a new storage format.

---

## 1. The question, split

**NL-01 — Scope.** OQ-3b-b asks whether a non-empty `history.ndjson` whose final record lacks a trailing newline is structurally valid. It does not decide the torn-line repair path (OQ-4), the row 570 window, or whether the two-representation store should be replaced.

**NL-02 — Five sub-questions, answered separately:** is the newline part of the format; does every supported writer emit it; can a legitimate file lack it; is an unterminated complete record distinguishable from a torn one; and would rejecting it declare recoverable content damaged.

**NL-03 — Answer in brief: the final newline is a writer invariant, not a format requirement.** The loader reads an unterminated final record in full, so the newline is not needed for availability. `V3-08`'s requirement is a scanner-added strictness, and it currently produces a false `Damaged` verdict while doing nothing to prevent the corruption a missing terminator actually causes. The outcome is **Amend**, with a writer-side rule that must be adopted alongside it.

---

## 2. The writer/reader pair

**NL-04 — Observed: all three writers of the log emit a terminator.**

| Writer | Site | Behaviour |
|---|---|---|
| `AppendHistoryEntry` | `persistence.go:688-702` | `json.Marshal`, `line = append(line, '\n')`, one `f.Write(line)` on an `O_APPEND` handle |
| `RewriteHistory` | `persistence.go:641-648` | each marshalled record, `buf.WriteByte(0x0A)`, whole buffer written to a temp file and renamed |
| first-append migration from `history.json` | `persistence.go:678-686` | each record, `append(line, '\n')`, direct append |

**NL-05 — Observed: the reader does not require the terminator.** `loadHistoryFile` uses `bufio.Scanner` over the file (`persistence.go:598-616`). `Scanner` returns the final token at EOF even when it is not newline-terminated; the loop `TrimSpace`es it, unmarshals it, and appends it to the history. A corrupt line breaks the loop and returns the prefix; a scanner error is logged and the partial slice is returned.

**NL-06 — Calculated: loadability is therefore unaffected by the final newline.** The format consumed by the product is "one JSON object per line, with a final line that may end at EOF". The newline separates records for the next append; it is not a validity marker for the last record.

**NL-07 — Observed: only one writer is atomic, and it cannot leave a partial file.** `RewriteHistory` writes a temp file and renames it, so its output is always whole. `AppendHistoryEntry` opens the real file with `O_APPEND|O_CREATE|O_WRONLY` and writes once, discarding the returned count and never syncing. The migration loop likewise ignores write errors and counts.

**NL-08 — Source-derived: there is no fourth writer.** No other production path writes `history.ndjson`; the remaining references are comments, the loader, and tests.

**NL-09 — Calculated: absence of the terminator therefore indicates an interrupted or foreign write, never normal completion.** Every supported writer terminates; a file that does not was either cut at the boundary, produced by an external tool, or migrated by a path that failed partway.

---

## 3. What the real loader does

**NL-10 — Method.** A disposable fixture created eight chat directories under `build/tmp/oq3bbfixture`, each with `meta.json`, and loaded them through `chat.NewPersistenceManagerWithDir` and the real `LoadChat`. No product file was touched.

**NL-11 — Results (Observed):**

| Fixture | Tail shape | Entries loaded | Records |
|---|---|---|---|
| `term` | three records, each newline-terminated | **3** | a,b,c |
| `unterm` | three records, last unterminated | **3** | a,b,c |
| `torn` | two complete records + partial third | **2** | a,b |
| `lone` | a single `"\n"` | 0 | — |
| `concat` | record 2 unterminated, record 3 appended directly | **1** | a |
| `crlf` | two CRLF-terminated records | 2 | a,b |
| `snapshotOnly` | no log, two-record `history.json` | 2 | a,b |
| `untermWithSnapshot` | unterminated two-record log, three-record snapshot | **2** | a,b |

**NL-12 — The decisive comparisons.** `unterm` and `term` load identically, so the missing terminator costs no content. `torn` loses only the partial record. `concat` loses records 2 **and** 3, because two JSON objects on one line fail `json.Unmarshal`; this is the corruption a missing terminator causes on the *next* append. `untermWithSnapshot` shows the log winning over a longer snapshot, as `V3-45` already says.

**NL-13 — Calculated: an unterminated log is not a partial representation.** Its records are the basis; the snapshot is not consulted. The tail condition changes nothing about which representation is authoritative.

---

## 4. Terminated versus torn: what can and cannot be known

**NL-14 — A parseable final line is a complete record.** A JSON object ends at its closing brace; there is no valid object prefix that is missing fields. If `json.Unmarshal` accepts the final line, every field of that record is present and readable. The missing byte is a separator, not content.

**NL-15 — A write cut exactly at the boundary is byte-identical to a foreign unterminated record.** `f.Write(line)` can fail after writing `len(line)-1` bytes, leaving `{…}` without `\n`; an external editor can save the same bytes deliberately. Nothing in the file distinguishes them.

**NL-16 — But in both cases no content of that record is missing.** The record is complete; the uncertainty is about what the *writer intended next*, not about what this record contains. That is why the terminated/unterminated distinction is not an availability property.

**NL-17 — The real harm is the next append.** Appending at EOF onto an unterminated tail produces one line holding two JSON objects. The probe shows the loader then reads neither: `concat` yielded 1 entry instead of 3. The damage is caused by concatenation, not by the absence of a newline at rest.

**NL-18 — A partial JSON tail is different and remains detectable.** `torn` ends mid-object and fails to parse; that is real truncation, and row 7's `truncated` is the right finding.

**NL-19 — Calculated: the meaningful distinction is parseable versus unparseable, not terminated versus unterminated.** The newline adds no information about the last record's completeness; it is a mutation-safety property of the boundary.

---

## 5. What strict rejection does, and does not do

**NL-20 — The current rule.** `V3-08` requires a non-empty log to end with a newline. A file that does not is structurally invalid; `V3-11` then denies it usable content, and if no snapshot holds content, row 7 reports **Damaged** — `malformed`, or `truncated` when the file ends mid-record.

**NL-21 — That is a false damage claim for content the product reads in full.** The `unterm` fixture loads all three records; calling the conversation Damaged says content is unavailable when it is present and displayable. The owner's caution is exact here: the requirement must follow the data format and the tolerance for recoverable content, not the scanner's convenience.

**NL-22 — And the strict rule protects nothing, because the scanner is read-only.** `V3-55` forbids the scanner from repairing or rewriting; no writer consults a verdict. A Damaged classification cannot stop the next append, so "reject to be safe" buys no safety at all.

**NL-23 — Accepting without a writer rule would convert a false Damaged into silent loss.** If the state became Healthy and the app appended normally, `concat` is the result: two records become unreadable. The availability fix and the mutation rule must be adopted together.

**NL-24 — The two harms are separable.** Availability asks whether the records are readable; the answer is yes. Mutation safety asks whether appending is safe; the answer is no until the tail is terminated. A verdict cannot answer the second, and a writer rule cannot answer the first.

---

## 6. The amendment

These are proposed clauses, not edits.

**NL-25 — `V3-08` is restated for the log:** a non-empty `history.ndjson` is structurally valid iff every non-blank line parses as a record. **The final record may be terminated or unterminated.** A final line that parses is a complete record; a final line that does not is the truncation case.

**NL-26 — The tail is recorded as an observation, not a state:** `tail = terminated | unterminated`, where `unterminated` means the last non-blank record is not followed by `\n`. It never changes the six-state vocabulary and never changes which row applies except where the final line is unparseable.

**NL-27 — Row 7's `truncated` narrows to the unparseable tail.** A file whose last line is partial JSON is `truncated`; a file whose last record parses but lacks `\n` is not `truncated` and not `malformed`.

**NL-28 — A Healthy verdict carries `unterminated_tail` when the observation is `unterminated`.** Rows 14 and 15 keep their state; the reason column gains the tail reason. A user-facing surface must not present the tail condition as data loss, and must not present it as "nothing to see"; it is a mutation warning attached to an available conversation.

**NL-29 — `V3-58` gains `unterminated_tail`; `V3-09` is retained unchanged.** A lone `"\n"` (1 byte) is still neither valid nor empty and remains `malformed`. Without that explicit retention the relaxed rule would make a 1-byte file vacuously valid; the amendment states both.

**NL-30 — The writer rule that makes acceptance safe: no append may concatenate onto an unterminated tail.** Under the mutation gate, an append to such a log either **terminates the tail inside the same critical section** (a single `\n` write before the record, which changes no content) or **refuses** with a typed result. It may never append directly. `RewriteHistory` already terminates every record, so any rewrite repairs the tail; row 570's final save is such a rewrite. Which of the two the gate chooses is a §3 design decision, not this record's.

**NL-31 — Coexistence is unchanged.** Row 6's prefix relation is over records, not bytes; an unterminated two-record log beside a three-record snapshot is still `log ⊂ snapshot` and still Unavailable `inconsistent`. The probe's `untermWithSnapshot` confirms the log is the basis, so the tail observation cannot be used to prefer the snapshot.

**NL-32 — Totality is unchanged.** No row is added and no state is added; `unterminated_tail` is a reason on an existing Healthy row, and `tail` is evidence recorded beside the verdict. Every input still reaches exactly one row.

---

## 7. What remains open, and what is not decided

**NL-33 — Repair ownership remains OQ-4's.** Terminating a tail is a write; the scanner may never perform it (`V3-47`, `V3-55`). Whether the product exposes a user-initiated repair, and whether the gate terminates or refuses, are separate decisions.

**NL-34 — Legacy prevalence is not measured.** Whether real user files without a final newline exist cannot be established without reading user data, which this investigation does not do. The amendment accepts them by rule rather than by measurement; that is a deliberate tolerance, not a claim about frequency.

**NL-35 — Torn-at-boundary frequency is Unestablished.** The mechanism is source-derived and the byte-identity is Observed; how often a write stops exactly at the boundary is not measured and is not needed for the rule.

**NL-36 — No storage redesign is decided or excluded.** If the two-representation store is replaced, the same rule applies to any line-delimited format: parseability decides the record, termination decides whether appending is safe. The record does not defend the current layout.

---

## 8. Outcome

**NL-37 — Outcome: Amend.** `V3-08`'s trailing-newline requirement is contradicted by the product's own reader (`unterm` loads in full) and does not protect against the harm it appears to address (the scanner cannot stop an append). The bounded amendment is:

| Clause | Change |
|---|---|
| `V3-08` | validity is "every non-blank line parses"; final record may be unterminated |
| §1 evidence | add `tail = terminated | unterminated`, a diagnostic not a state |
| `V3-39` row 7 | `truncated` only when the final line does not parse |
| `V3-39` rows 14–15 | reason `unterminated_tail` when the tail is unterminated; state unchanged |
| `V3-58` | add `unterminated_tail`; retain `V3-09` |
| §3 gate | no append onto an unterminated tail; terminate-in-section or refuse, never concatenate |
| `V3-73` | OQ-3b-b closes, subject to owner adoption |

`V3-09`, row 6, `V3-45`, `V3-47`, `V3-55` and the six states are unchanged.

**NL-38 — The dependency is stated rather than hidden.** Adopting the classification without the writer rule is worse than the current strictness: it would report Healthy and then lose records on the next append. The two clauses are one change.

---

## 9. Tests this would require

- **NL-T1** A three-record log whose last record is unterminated loads three records and scans Healthy `unterminated_tail`, not Damaged.
- **NL-T2** A log ending in partial JSON scans Damaged `truncated` when no usable content exists; the partial record is never counted.
- **NL-T3** A lone `"\n"` remains `malformed`, not valid-and-empty (`V3-09` retained).
- **NL-T4** An append onto an unterminated tail is refused, or terminates the tail first; it never concatenates. The concatenated fixture demonstrates the loss it prevents (1 entry instead of 3).
- **NL-T5** A rewrite terminates every record and clears `unterminated_tail`.
- **NL-T6** CRLF and trailing whitespace on the final record still parse and do not change the verdict except for the tail observation.
- **NL-T7** An unterminated two-record log beside a three-record snapshot still reports row 6 `inconsistent`; the log remains the basis.
- **NL-T8** The negative assertion: no verdict is derived from the presence or absence of a trailing newline alone; parseability decides the record and the gate decides append safety.

---

## 10. Nothing here is authorized

**No implementation follows.** This record authorizes no change to `loadHistoryFile`, `AppendHistoryEntry`, `RewriteHistory`, the gate, the scanner or any UI copy; no newline write, no repair, no migration, no new state, no deletion enablement and no storage replacement. It changes what the contract may call valid and states the writer rule that must accompany it.
