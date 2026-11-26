# TASK-03 — OQ-3c: coexistence of the two representations, and strict validity

| Field | Value |
|---|---|
| Purpose | Establish what coexistence actually means; replace "at least one record loaded" as a promotion test |
| Status | **Contract addendum. Read-only investigation. No code changed.** |
| Amends | row 551 RC4-12…RC4-14 (tightened); row 549 RC3-03 (refined); row 547 unedited |
| Evidence | `build/tmp/spike11/probe11.go`, `build/tmp/spike11-oq3c.log` |
| Recovery token | `TASK-03-RECONCILIATION-OQ3C-2026-11-22` |

---

## 1. The four distinctions, normative

**RC6-01** — Four properties, never interchangeable:

| Property | Question | Who can answer it |
|---|---|---|
| **Loadability** | can the app display it? | `LoadChat` |
| **Structural validity** | does the whole representation conform, with no silently ignored interior or tail corruption? | a whole-file validator |
| **Content presence** | does a structurally valid representation hold committed records? | the validator's record count |
| **Completeness** | does the evidence establish the expected content is still there? | the reconciliation decision, using both representations and the marker |

**RC6-02 — `LoadChat` output is not a validity test.** RC4-09 showed why; §2 below shows it is worse than "needs caution": the loader can prefer the *lossier* representation and return a clean-looking count that omits records that still exist on disk.

**RC6-03 — A scanner may call the loader as a read-only observation tool** (row 551 RC4-18 confirms it mutates nothing), **but must not equate its output with structural validity or completeness.**

**RC6-04 — "Healthy" may never mean "the app can display whatever survived."** If that is all the evidence supports, the state is not Healthy.

---

## 2. What coexistence actually is — **Observed**

Selection rule (source: `loadHistoryFile`): use `history.ndjson` when it exists **and is larger than zero bytes**; otherwise fall back to `history.json`.

| # | ndjson | json | LoadChat | Verdict |
|---|---|---|---|---|
| 1 | 2 records | 2 records, same | 2 | agree |
| 2 | 1 record `from-ndjson` | 1 record `from-json` | 1, **`from-ndjson`** | **distinct content; ndjson wins** |
| 3 | 1 record | **3 records** | **1** — two records on disk never shown | **precedence loses data** |
| 5 | 1 record | `[]` | 1 | ndjson only |
| 6 | 1 good + torn tail | `[]` | 1, silently | torn log reads fine |
| 7 | 1 good + torn tail | **2 valid** | **1** | **identity of interest: the valid file is ignored** |
| 8 | 0 bytes | `[]` | 0 | both empty |
| 9 | **0 bytes** | **2 records** | **2 — `deleted-a` returns** | **resurrection** |
| 10 | absent | 1 record | 1 | legacy fallback |

**RC6-05 — Coexistence is not benign.** It is benign only when the two agree. Cases 3, 7 and 9 are the counterexamples:

- **case 3** — both structurally valid, different record counts; the loader takes the shorter one and the extra records are invisible;
- **case 7** — the log is torn and the snapshot is valid and fuller; the loader still takes the torn log;
- **case 9** — a zero-byte log with a populated snapshot **reloads content that was deleted**. `RewriteHistory` renames the new log into place and then writes the snapshot, so a crash between those two writes (or a clear whose snapshot lagged) resurrects the old history on next load.

**RC6-06 — "Which representation is fresher" is not answerable by inspection of the files.** During normal appends the snapshot lags the log (row 549 P2). After a torn append the log lags the snapshot. After a clear the log is authoritative and the snapshot may be stale. No timestamp or size comparison decides it.

**RC6-07 — Precedence is a display rule, not an evidence rule.** Reconciliation must treat each representation as its own witness and never inherit the loader's ordering.

---

## 3. Tightened promotion — the criterion must be impossible to satisfy by truncation

Row 551's RC4-12 said "parses successfully and yields at least one record". **That is insufficient**: a truncated log with two valid records followed by corruption satisfies it while silently losing the rest (case 6 satisfies it outright).

**RC6-08 — Structural validity, per format:**

| Format | Valid iff |
|---|---|
| `history.ndjson` | **every** non-blank line parses as a record **and** the file ends with a newline (a record without one is a torn write) |
| `history.json` | the **whole** document parses as an array of records |

**RC6-09 — Promotion to `committed` requires a structurally valid representation holding ≥1 record.** A prefix is not a record set.

**RC6-10 — Any invalid representation blocks promotion**, even when another is valid, until the invalid one is explained. Reporting "Healthy" because the good file covered for the bad one is exactly the silent-partial failure this contract exists to stop.

**RC6-11 — When the two disagree, the larger valid set is the completeness basis and the disagreement is reported.** The smaller set is never assumed to be the truth (case 3: it was not).

**RC6-12 — Case 9 is reported as `Inconsistent`, never as content-present.** A zero-byte log beside a populated snapshot is evidence of a torn rewrite, and the snapshot's records are of **unestablished** validity as current content — they may be resurrected deletions. Reconciliation reports the condition and changes nothing.

---

## 4. Condition vocabulary

Row 547's states are unchanged; these are the **reasons** a state carries.

| Reason | Meaning |
|---|---|
| `empty` | structurally valid, zero records, and nothing expected |
| `malformed` | does not conform at all (wrong shape, unparseable document) |
| `truncated` | conforms up to a point, then ends mid-record or without a final newline |
| `inconsistent` | both representations valid but disagreeing, or a zero-byte log beside a populated snapshot |
| `absent` | no accepted representation present at all |

**RC6-13 — A torn tail is real content loss.** The app cannot recover the record that was being written, so with the marker `committed` the conversation is **Damaged (reason `truncated`)** — not Healthy, not merely "partial display".

---

## 5. Cost of strict validation

**RC6-14 — Whole-representation validation costs the same order as the full parse already measured** (row 551 §5: 0.50 ms at 100 messages, 10.2 ms at 10,000, 104.7 ms at 100,000). It is therefore an **on-decision** cost, not a per-pass one: the µs tier decides *whether* to validate, the ms tier performs it when a marker or a state change depends on the answer.

**RC6-15 — The measurements remain evidence.** They are not a budget, a target, or a licence to optimize.

---

## 6. Tests added

- **RT6-01** both valid and agreeing → Healthy.
- **RT6-02** both valid, different counts → reported `inconsistent`; the larger set is the basis; never silently the smaller.
- **RT6-03** torn log + valid fuller snapshot → **not** Healthy-by-loadability; `truncated` reported, snapshot not used to excuse it.
- **RT6-04** zero-byte log + populated snapshot → `inconsistent`, promotion refused, nothing deleted or rewritten.
- **RT6-05** valid document that is not an array → `malformed`, never promoted.
- **RT6-06** log whose every line parses but which lacks a final newline → `truncated`.
- **RT6-07** absent log + populated snapshot (legacy) → valid content, Healthy, promoted.
- **RT6-08** promotion blocked while any representation is invalid.
- **RT6-09** the scan makes no write except the marker (byte-identical directory otherwise).

---

## 7. Audit carried forward

The replacement contract (row 547) must be re-checked against three findings before it can be finalized:

| Finding | Clauses to reconcile |
|---|---|
| content-based representation | row 549 RC3-01…04, row 551 RC4-02…RC4-07 |
| expectation lifecycle | row 548 RC2-01…RC2-08, row 551 RC4-12…RC4-14, RC6-08…RC6-13 |
| unique conversation identity | row 550 RC5-1…RC5-11 — retained as an independent defect, no authorization |

**RC6-16 — Additional audit requirement:** every clause that determines a state must name which of the four properties it relies on, so no clause can be read as "the loader could display it, therefore Healthy".

**RC6-17 — Contract work only.** No code, no repair, no migration. Creation identity stays separate and unauthorized.

---

## 8. Still open

| # | Question |
|---|---|
| OQ-3c-a | Should the app itself stop preferring a torn log over a valid fuller snapshot? A product question, not a reconciliation one. |
| OQ-3c-b | Is `RewriteHistory`'s two-write window (log then snapshot) a defect to record separately, given it can resurrect deleted content? It looks like one; not opened here. |
| OQ-3c-c | Whether `inconsistent` should ever be user-actionable, or report-only. |
