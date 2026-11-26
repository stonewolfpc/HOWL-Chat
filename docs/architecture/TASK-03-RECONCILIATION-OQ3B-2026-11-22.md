# TASK-03 — OQ-3b: what counts as content, what is unreadable, and what a scan costs

| Field | Value |
|---|---|
| Purpose | Settle OQ-3b; state the content-based representation rule; qualify the marker |
| Status | **Contract addendum. Read-only investigation. No code changed.** |
| Amends | row 549 RC3-01…04 (confirmed, restated as RC4-); row 548 RC2-05 stays superseded |
| Evidence | `build/tmp/spike10/probe10.go`, `build/tmp/spike10b/probe10b.go`, `build/tmp/spike10-oq3b.log`, `build/tmp/spike10b.log` |
| Recovery token | `TASK-03-RECONCILIATION-OQ3B-2026-11-22` |

---

## 1. The representation rule, stated once

**RC4-01 — The question is not "does a file exist" but "can Howl still load this conversation's committed history".**

| Rule | Statement |
|---|---|
| RC4-02 | `history.json` = `[]` is not evidence of committed content. |
| RC4-03 | A populated `history.json` is valid even with no `history.ndjson`. |
| RC4-04 | A populated `history.ndjson` is valid without a populated snapshot. |
| RC4-05 | Missing one representation is not damage if another accepted representation preserves the content. |
| RC4-06 | No accepted representation containing expected content means **Damaged**, subject to the expectation marker. |

**RC4-07 — The historical correction stays explicit.** Row 548's version-scoped rule was wrong and is superseded, not quietly dropped; the reason it was wrong is observable (row 549 §2) and remains part of the record.

---

## 2. What empty looks like, per format — **Observed**

| Representation | Empty form | Bytes |
|---|---|---|
| `history.json` (`saveJSON` → `MarshalIndent`) | `[]` | **2** |
| `history.ndjson` (`RewriteHistory([])`) | zero bytes | **0** |
| `history.ndjson` with one stray newline | `\n` | 1 — neither empty nor content |

**RC4-08 — A threshold of "size > 2" for json and "size > 0 plus at least one parseable record" for ndjson** is the correct shape. The 1-byte newline case shows why size alone is insufficient: it is above the empty size and contains nothing.

---

## 3. Can malformed be told from empty? — **Observed, and the answer is no**

`loadHistoryFile` never errors. Loading `history.json` with these bodies:

| Body | Messages | Error |
|---|---|---|
| `[]` (empty) | 0 | nil |
| truncated `[{"id":"0","text":"tr` | **0** | **nil** |
| object, not array | **0** | **nil** |
| array of strings | **0** | **nil** |
| 1 valid entry | 1 | nil |
| 3 valid entries | 3 | nil |

For `history.ndjson`:

| Condition | Messages | Error |
|---|---|---|
| 3 records torn by 40 bytes mid-record | **2** | **nil** |
| 2 good + torn tail | **2** | nil |
| 2 good + corrupt line + 2 good | **2** | nil — stops at the bad line, later records vanish from view |
| 2 complete records, no final newline | **2** | nil |

**RC4-09 — The loader is lossy by design: it cannot distinguish empty, malformed, and truncated.** Three different truths collapse to "fewer messages, no error". It therefore **cannot be the scanner's observation primitive**.

**RC4-10 — Partial content is indistinguishable from complete content.** A history torn mid-record loads as a shorter healthy history. This is not a defect of the loader — fail-soft is right for the app — but it means reconciliation's `Healthy` means *"the loadable prefix is present"*, not *"every committed record is present"*. The contract must say so.

**RC4-11 — The one detectable tear is a missing final newline**, and only when the tear does not land on a record boundary. The boundary case (2 complete records, no trailing newline) reads as 2 records. **Unestablished** whether it can ever be distinguished; recorded as a limit, not solved.

---

## 4. Promotion must be grounded in a successful load

**RC4-12 — `pending` promotes to `committed` only when an accepted representation parses successfully and yields at least one record.** Not on file existence, not on size, not on mtime.

**RC4-13 — Malformed never promotes.** A representation that exists but does not parse leaves the marker where it is and is reported as its own condition; it must not be silently promoted to "content is expected and present".

**RC4-14 — The scanner parses for itself.** It may share the parser, but it must not adopt the loader's conflation of failure with emptiness — otherwise a corrupt conversation is promoted to `committed` (or read as empty) with no evidence.

---

## 5. Cost — **Observed** on this machine, best of five

| Messages | Bytes | Full load | Throughput |
|---|---|---|---|
| 100 | 14.9 KB | **0.50 ms** | 28 MB/s |
| 1,000 | 150 KB | **1.54 ms** | 93 MB/s |
| 10,000 | 1.51 MB | **10.2 ms** | 141 MB/s |
| 100,000 | 15.2 MB | **104.7 ms** | 138 MB/s |

Per-conversation non-parsing checks: **`stat` ≈ 36–39 µs**, **final-byte read ≈ 2–3 µs**.

**RC4-15 — Availability does not require parsing.** `stat` plus a final-byte probe answers "is there content" in microseconds; a full parse costs ~0.5 ms at 100 messages and ~105 ms at 100,000. Counting records *does* require the parse.

**RC4-16 — So the scan is tiered.**

| Tier | What | Cost |
|---|---|---|
| every pass | existence, empty-form size, final-byte probe | µs per conversation |
| on a marker decision | one parse of the accepted representation | ms |
| on a count needed for display | full parse | ms |

**RC4-17 — Numbers are machine- and format-specific, and are recorded as a measurement, not a budget.** They are not a performance target and must not become an optimization exercise.

---

## 6. Can the scanner reuse the normal loader? — **Observed: yes for safety, no for truth**

Repeated `LoadChat` calls left the directory byte-identical (hash `19ac31e7cf63e048` before and after three loads).

**RC4-18 — `LoadChat` does not mutate conversation state**, so using it cannot damage data.

**RC4-19 — It still must not be the scanner's observer**, per RC4-09: safety and fidelity are different questions, and the loader answers only the first.

---

## 7. Marker persistence qualification

**RC4-20 — The marker is persisted through the established `DurableStore`** (`internal/durable`), not a hand-rolled file: it already provides an exclusive lock, temp write under lock, manifest commit with the hash and size, pre-images, and a deterministic death injection point (`testHookBeforeReplace`) for scheduling a crash between manifest commit and replacement.

**RC4-21 — The marker's own crash behaviour must be tested**, not inferred from the store's: after a simulated death mid-save, the reloaded marker must be either the previous value or the new one, never a third state and never unreadable.

**RC4-22 — Otherwise the ambiguity merely moves into another file**, which is exactly what this clause exists to prevent.

---

## 8. Tests added

- **RT4-01** `history.json` = `[]` → not content.
- **RT4-02** ndjson of one newline → not content.
- **RT4-03** populated `history.json`, no ndjson → Healthy, promote on successful parse.
- **RT4-04** populated ndjson, `history.json` = `[]` → Healthy.
- **RT4-05** truncated `history.json` → not promoted, reported as unreadable, never Healthy-by-promotion.
- **RT4-06** corrupt ndjson line mid-file → reported as unreadable/partial, never promoted.
- **RT4-07** availability pass does not parse (bounded by the µs tier on a 100,000-message fixture).
- **RT4-08** marker save through `DurableStore` with death injected at `testHookBeforeReplace` → reload yields previous or new value only.
- **RT4-09** scan leaves a conversation directory byte-identical.

---

## 9. Still open

| # | Question |
|---|---|
| OQ-3b-a | Exact per-format empty threshold constants as contract numbers rather than measurements. |
| OQ-3b-b | Whether a missing final newline should be Damaged, a warning, or nothing — it is detectable but harmless when the content loads. |
| OQ-3b-c | Where the marker document lives (one per conversation, or one for all) and how it scales. |

**RC4-23 — Contract work only.** No code, no repair, no migration. OQ-3c (dual representation) is next, then the audit of the replacement contract against all three findings.
