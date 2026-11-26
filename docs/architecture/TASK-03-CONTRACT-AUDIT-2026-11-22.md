# TASK-03 — audit of the replacement contract against the reconciliation findings

| Field | Value |
|---|---|
| Purpose | Audit row 547 clause by clause against OQ-3, OQ-3a, OQ-3b, OQ-3c |
| Status | **Audit. Nothing adopted, nothing implemented, no authorization.** |
| Audits | `TASK-03-RECONCILIATION-CONTRACT-2026-11-22.md` (row 547) |
| Corrects | row 552 RC6-11 (**superseded**); RC6-10 (reworded, see §3.6) |
| Recovery token | `TASK-03-CONTRACT-AUDIT-2026-11-22` |

Row 547 is **not edited**. This audit states what a v2 must contain; the contract stays as written until it is replaced.

---

## 1. Correction first: count is not provenance

Row 552 RC6-11 said disagreement uses the **larger valid set** as the completeness basis. Rejected, correctly.

**RC7-01 — Disagreement fails closed.** Two structurally valid representations with different record sets yield `Inconsistent`. **No winner is chosen**, and neither set becomes the basis. Larger is not newer: the resurrection result (row 552 §2 case 9) is precisely a populated snapshot holding messages a clear intended to remove.

**RC7-02 — Count is not provenance, and size is not age.** Neither representation carries a trustworthy "written at"; both are written by different paths at different times (row 549 P2/P5/P6).

**RC7-03 — Provenance is testable only as a relationship.** The strongest available evidence is that one representation is an **exact prefix** of the other on the full record content in order.

**RC7-04 — Record IDs cannot carry that test.** IDs are **positional indices**, not stable identities: `AppendChatMessage` writes `fmt.Sprintf("%d", len(a.chatMessages)-1)` (`app-prompt.go:300-301`) and `ReplaceChatHistory` writes `fmt.Sprintf("%d", i)` (`app-prompt.go:349`). Deleting a message renumbers everything after it. A prefix match must therefore compare `sender`, `text`, `role`, `timestamp` (and `hasImage`/`characterId` when present) **in order**, and identical indices prove nothing on their own. **Source-derived.**

**RC7-05 — Even a proven prefix does not resolve the state.** It bounds the possibilities — the shorter is an earlier state of the longer, or the longer is a resurrected earlier state — and it cannot distinguish them. So it is **diagnostic**, reported as `prefix-related`, and the state stays `Inconsistent`.

**RC7-06 — An unprovable relationship stays unproven.** No heuristic — count, size, mtime, "the log is usually right" — may convert it into Healthy.

---

## 2. The audit's six questions, and what row 547 currently answers

| # | Question | Row 547's answer | Verdict |
|---|---|---|---|
| 1 | **Presence** — which representations are expected at each stage? absent optional vs invalid present | alternatives group "≥1 member exists" (RC-08/09) | **inadequate** — existence is not validity, and nothing distinguishes absent-optional from invalid-present |
| 2 | **Structural validity** — does the whole representation parse? | not addressed | **missing** |
| 3 | **Content presence** — records, not bytes or `[]` | not addressed | **missing** |
| 4 | **Completeness** — evidence that committed history survives, including corrections, retractions, clears | not addressed | **missing** |
| 5 | **Coexistence** — can the relationship of the two be established? | "at least one member" | **inadequate** — treats them as interchangeable |
| 6 | **Promotion** — what evidence permits `none`/`pending` → `committed`? | marker does not exist | **missing** |

---

## 3. Clause-by-clause

### 3.1 Stands unchanged

| Clause | Why it survives |
|---|---|
| RC-01…RC-04 | the retired-identity framing is still exact; nothing here disturbs it |
| RC-06, RC-07 | manifest generated from the storage contract, never from disk — the rule that stops a damaged conversation redefining its own requirements |
| RC-10 | no manifest → Unavailable, never Healthy |
| RC-16, RC-17 | error is never absence; incomplete enumeration is an error |
| RC-18, RC-20, RC-21, RC-22 | deletion flow, worklist is not a receipt, sweep is authority, orphans reported |
| RC-23 | deletion confined to the conversation subtree — still the only identity rule |
| RC-24…RC-27 | cadence, bounded work, backoff, concurrency |
| RC-28, RC-29 | wording rule; no implementation authority |

### 3.2 Contradicted by the evidence — must change

**RC-05 — the file table mislabels the snapshot.** It calls `history.json` "legacy, migrated into ndjson". It is neither: `SaveChat` writes it on every save (`persistence.go:211`), `RewriteHistory` writes both (`persistence.go:661`), and `SaveChatAs` writes **only** it (row 549 P6). Calling it legacy is what produced the version-scoped error in row 548. **Amend:** "second accepted representation — written by save, rewrite and duplication; not legacy."

**RC-09 — "content requires at least one member of the history alternatives group"** is wrong for empty conversations (row 548 RC2-01) and wrong for a `[]`-only directory. **Superseded.**

**RC-11 — Healthy** is defined as "every `required` path exists and every `alternatives` group has ≥1 member". Path existence is satisfied by a 2-byte `[]`. **Amend** to the four properties (§4).

**RC-11 — Missing**: "the directory does not exist, **or all required content is absent**". The second half contradicts RC-12's own Damaged→Missing prohibition when `meta.json` is the only survivor. **Amend:** all content absent **and** `meta.json` absent → Missing; `meta.json` present with content absent → Damaged.

**RC-13 — Damaged is described only as "the state a failed deletion produces".** Now also produced by a torn write and by a partial clear. **Amend.**

**RC-14 — "never creates"** conflicts with the marker, which the scanner must write. **Amend:** never writes conversation data; may write the marker document only, through `DurableStore` (row 551 RC4-20…RC4-22).

**RC-15 — outcome table** has "directory present, no content at all → Missing" and no row for truncated, inconsistent, malformed, or empty-with-no-expectation. **Amend** per §5.

**RC-08 — alternatives group** must become: at least one member structurally valid **and** holding ≥1 record.

**RC-19 — ordering** must extend to the marker: `pending` written **before** the append, promoted after a successful parse (rows 549 RC3-09…RC3-12).

**RT-05 — "`history.ndjson` removed while `meta.json` and others remain → Damaged, not Healthy"** is now **wrong**: if a populated `history.json` remains, the content is still loadable and the answer is Healthy (row 549 RC3-03). **Superseded.**

**RT-06 — "both absent → Damaged"** holds only when the marker is `committed`. **Amend.**

### 3.3 Missing entirely (v2 must add)

Expectation marker and its lifecycle · promotion evidence · structural validity · content presence · completeness · coexistence relationship · the four-property labelling requirement · reason vocabulary · clear as a retraction · `DurableStore` persistence and its crash test.

### 3.4 RC6-10 reworded

Row 552 RC6-10 said "any invalid representation blocks promotion". Read literally that would also demote a conversation whose log is valid and whose mirror is corrupt. **RC7-07 — the invalid representation may never be the *basis* of promotion; it does not by itself prevent a state derived from a valid one.** Its own condition is reported alongside (reason `malformed`), so a broken mirror is visible without being mistaken for lost content.

### 3.5 A structurally valid empty representation is not a valid empty conversation

**RC7-08 — Emptiness means nothing on its own.** Its meaning is read against two facts:

| Marker | Representations | Meaning |
|---|---|---|
| `none` | structurally valid, empty | never committed — **valid empty conversation** |
| `pending` | structurally valid, empty | an append was attempted and did not land — **Pending reconciliation** |
| `committed` | structurally valid, empty | committed content is gone, **or was cleared** — see RC7-09 |

**RC7-09 — A clear is a retraction, and the marker must be able to record it.** Without that, every deliberately cleared conversation is permanently Damaged. So the marker needs a fourth value — `cleared` — or an explicit cleared-flag beside it:

| Event | Marker after |
|---|---|
| first successful commit | `committed` |
| `ClearChatHistory` completes | **`cleared`** |
| delete of individual messages (survivors remain) | `committed` (partial retraction; content still present) |
| `cleared` + empty | **Healthy** (a valid, empty conversation) |
| `cleared` + records present | content exists; Healthy on content grounds, reported as a reappearance |

**RC7-10 — This is a new finding, not a restatement.** `RewriteHistory([])` currently leaves no trace of intent (`app-prompt.go:406`), so without `cleared` the reconciliation system cannot tell "cleared" from "lost".

### 3.6 Promotion

**RC7-11 — Promotion to `committed` requires exactly this evidence, and nothing else:** one representation that (a) is present, (b) parses **in full** according to its format, (c) ends lawfully (final newline for ndjson), (d) yields ≥1 record. Not `LoadChat`'s count. Not size. Not mtime. Not the larger set. Not the loader's success.

**RC7-12 — Malformed or truncated content never promotes**, and a proven prefix never promotes on its own (RC7-05).

---

## 4. The four properties, as clauses

**RC7-13** — Every state-determining clause must name its basis:

| Property | Permitted use |
|---|---|
| loadability | display only; never a state |
| structural validity | gates promotion and content presence |
| content presence | drives Empty vs Present |
| completeness | the only basis for Healthy on content |

**RC7-14 — Healthy requires completeness evidence**, never mere displayability.

---

## 5. Outcome table, amended

| Observation | State | Basis |
|---|---|---|
| marker `none`, valid empty | Healthy | presence, emptiness explained |
| marker `cleared`, valid empty | Healthy | retraction recorded |
| valid representation with ≥1 record, marker `none`/`pending` | Healthy, promoted | structural validity + content presence |
| both valid and agreeing | Healthy | coherence |
| both valid, differing, not prefix-related | **Inconsistent** (Unavailable, reported) | coexistence |
| both valid, exact prefix relation | **Inconsistent**, `prefix-related` diagnostic | coexistence, no resolution (RC7-05) |
| one representation truncated, no other | **Damaged**, reason `truncated` | structural validity fails; real loss |
| one representation malformed, no other | **Damaged**, reason `malformed` | structural validity fails |
| valid log + malformed mirror | Healthy on content, `malformed` companion reported | RC7-07 |
| marker `committed`, nothing present | **Damaged**, reason `absent` | expectation unmet |
| marker `pending`, nothing present | **Pending reconciliation** | RC3-09 |
| `meta.json` absent, directory present | **Damaged** | identity lost |
| directory absent | **Missing** | total absence |
| listing incomplete or path denied | **Unavailable** | error is not evidence |

**RC7-15 — `Inconsistent` maps to Unavailable, not to a seventh state.** The six states stay; inconsistency is a reason the condition could not be established. It fails closed, is retained and retried, and never becomes Healthy without new evidence.

---

## 6. Defect boundaries stay intact

| Finding | Owner |
|---|---|
| reconciliation observes and reports inconsistent history | **TASK-03** (this contract) |
| creation ID collision, non-exclusive directory creation | **row 550** — independent, unauthorized |
| loader prefers a torn log over a fuller valid snapshot | **OQ-3c-a** — product question, not absorbed |
| `RewriteHistory` can resurrect cleared content after a crash | **OQ-3c-b** — product defect, not absorbed |

**RC7-16 — None of the three product defects is fixed, worked around, or hidden by this contract.** Reconciliation reports the symptoms; it never compensates for them. In particular the scanner must not "choose the good file" to paper over OQ-3c-a, and must not re-write a log to close the OQ-3c-b window.

---

## 7. What v2 must contain

A consolidation, not another addendum: row 547's surviving clauses, with RC-05/08/09/11/13/14/15/19 amended, RT-05/RT-06 replaced, the marker (four values) and promotion rules integrated, the four-property labelling applied clause by clause, the reason vocabulary attached to the states, and the OQ-3c-a/OQ-3c-b boundary stated in the contract itself so a future implementer cannot absorb them.

**RC7-17 — No implementation authorization.** No code, no manifest, no migration, no deletion enablement, no marker written anywhere.

---

## 8. Still unestablished

| # | Item |
|---|---|
| A-1 | Whether a prefix relation can ever be strengthened into provenance; needs a writer-side record (the marker gives intent, not per-record proof). |
| A-2 | Whether `cleared` belongs in the marker or in a separate retraction record. |
| A-3 | Whether a torn mirror (invalid `history.json` beside a valid log) should ever be user-visible or only a log line. |
| A-4 | OQ-1, OQ-2, OQ-4, OQ-5 remain open from row 547; OQ-3 is closed by row 548. |
