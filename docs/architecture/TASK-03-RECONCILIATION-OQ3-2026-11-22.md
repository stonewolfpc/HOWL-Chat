# TASK-03 — OQ-3 resolved: empty conversations, and a stale-history finding

| Field | Value |
|---|---|
| Purpose | Verify the creation and first-message writer lifecycle; fix the manifest's history rule |
| Status | **Contract addendum. Amends row 547 by addition; row 547 unedited** |
| Supersedes | row 547 RC-09 (alternatives group for all conversations); refines RC-08, RC-11, RC-15 |
| Artifacts | `build/tmp/spike8/probe8.go`, `build/tmp/spike8-oq3.log` |
| Recovery token | `TASK-03-RECONCILIATION-OQ3-2026-11-22` |

Prefixes: `RC2-` normative, `RT2-` test. Labels as before.

---

## 1. What was measured

Scratch root, `NewPersistenceManagerWithDir`, one `NewChat`, then one `AppendHistoryEntry`. No user data touched. **Observed:**

| Moment | Directory contents |
|---|---|
| after `NewChat` (no messages) | `meta.json` (421 b), **`history.json` = `[]` (2 b)**, `characters.json`, `lorebooks.json`, `memories.json` (each `[]`) — **no `history.ndjson`** |
| after first `AppendHistoryEntry` | `history.ndjson` appears (132 b); **`history.json` still `[]` (2 b)** |

`world.json` and `scenario.json` are absent at creation because `NewChat` leaves them `nil` and `SaveChat` writes them only when non-nil (Source-derived: `persistence.go:217,224`).

---

## 2. The decision you proposed, now grounded

**RC2-01 — Empty conversations are valid.** A conversation whose metadata is valid and which has never committed content is **Healthy**, not Damaged, with `history.ndjson` absent.

**RC2-02 — The three distinctions are separate, and the manifest must not merge them.**

| Distinction | Established by |
|---|---|
| conversation **existence** | valid `meta.json` |
| conversation **content** | committed history records |
| whether history was **expected** | a creation-time expectation recorded in the manifest, flipped when content is first committed |

**RC2-03 — The four conditions.**

| Condition | Expectation | Reconciliation |
|---|---|---|
| empty; no history expected yet | `history_expected = false` | Healthy with `meta.json` alone |
| content committed | `history_expected = true` | one accepted history representation must be present **and non-empty** |
| expected, all accepted representations absent or empty | `history_expected = true` | **Damaged** |
| cannot establish which applies (manifest missing/unreadable) | unknown | **Unavailable** |

---

## 3. The second finding: `history.json` presence is not evidence of content

**RC2-04 — `history.json` is a snapshot that lags.** After the first committed message the file was still `[]`. `AppendHistoryEntry` writes only the append log (`persistence.go:664-698`); `history.json` is rewritten by `SaveChat` and `RewriteHistory`. So a conversation with thousands of messages can carry a 2-byte `history.json`.

**Calculated consequence for row 547's rule:** "at least one of `history.ndjson` or `history.json`" is **wrong in both directions**.

| Case | What RC-09 would have said | What is true |
|---|---|---|
| empty conversation, both files `[]` | alternatives satisfied → Healthy | Healthy, but for a different reason: nothing was expected |
| conversation with content, `history.ndjson` deleted, stale `history.json` = `[]` still present | alternatives satisfied → **Healthy** | **Damaged** — the content is gone and nobody notices |

That second row is the failure this addendum exists to prevent: RC-09 would have reported a conversation whose messages were deleted as perfectly healthy, because an empty placeholder file still existed.

**RC2-05 — The accepted representation is version-scoped, not "either".**

| `contract_version` | Accepted content representation |
|---|---|
| legacy (JSON snapshot era) | `history.json` **with content** (larger than the empty-array form) |
| ndjson era (current) | `history.ndjson` **with at least one record** |

**RC2-06 — "Non-empty" is a property of content, not of existence.** The check is at minimum file size above the empty-array form, and ideally one successfully parsed record. A path that exists but contains nothing is not content.

---

## 4. Flipping the expectation

**RT2-01 — Order matters; append first, then mark.**

| Order | Crash between the two | Resulting state |
|---|---|---|
| **append, then mark** | content present, expectation still false | Healthy (benign); next successful append sets it |
| mark, then append | expectation true, no content | **Damaged** — a false alarm about data that never existed |

**RC2-07 — So the flip happens after the commit succeeds**, and a conversation that gains content without a manifest update is not damaged — it is merely under-described. The failure direction is chosen deliberately: never report damage that did not happen.

**RC2-08 — `UpdatedAt` is not the expectation.** `SaveChat` rewrites it on every save, including for an empty conversation, so it cannot distinguish "never had content" from "had content and lost it".

---

## 5. Changed rows in the state mapping

Row 547 RC-15 gains these rows; the rest stand:

| Observation | State |
|---|---|
| `meta.json` valid, `history_expected = false`, no history representation | **Healthy** (empty conversation) |
| `meta.json` valid, `history_expected = true`, accepted representation present and non-empty | Healthy |
| `meta.json` valid, `history_expected = true`, accepted representation absent **or empty** | **Damaged** |
| `meta.json` valid, expectation unknown | **Unavailable** |
| `meta.json` absent, directory present | **Damaged** (identity lost, data may remain) |
| directory absent | Missing |

**RC2-09 — `meta.json` absent is Damaged, not Missing.** The directory and possibly the user's messages still exist; calling that Missing would tell the UI to drop the conversation while its data is still on disk.

---

## 6. Tests added

- **RT2-02** A freshly created empty conversation scans Healthy with no `history.ndjson`.
- **RT2-03** A conversation with content whose `history.ndjson` is removed but whose stale `history.json` is `[]` scans **Damaged**.
- **RT2-04** "Non-empty" is checked, not mere existence.
- **RT2-05** Expectation flip after commit: simulated crash between append and mark leaves Healthy, never Damaged.
- **RT2-06** `UpdatedAt` is never consulted to infer expectation.
- **RT2-07** Legacy conversation with a populated `history.json` and no ndjson scans Healthy under the legacy version.
- **RT2-08** `world.json` / `scenario.json` absence never affects state (they are not written while nil).
- **RT2-09** Missing `meta.json` in a present directory scans Damaged, not Missing.

---

## 7. Still open

| # | Question |
|---|---|
| OQ-3a | Does any non-append path commit content (import, migration, restore)? If so it must set the expectation too. |
| OQ-3b | What is the exact "empty" size for each representation, and is a record-count check affordable per scan? |
| OQ-3c | Do older conversations exist whose `history.json` is populated *and* `history.ndjson` exists (dual state)? Which wins for the expectation? |

Row 547's OQ-3 is closed by RC2-01…RC2-09 above.

**RC2-10 — Still contract work only.** No implementation, no migration, no deletion enablement.
