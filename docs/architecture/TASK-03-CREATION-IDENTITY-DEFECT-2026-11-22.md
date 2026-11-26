# TASK-03 — creation identity defect (separate from reconciliation)

| Field | Value |
|---|---|
| Purpose | Record the creation-collision defect as its own finding, with its intended invariant |
| Status | **Defect record. No repair authorized. No code changed.** |
| Found | OQ-3a inventory, row 549 |
| Evidence | `build/tmp/spike9/probe9.go`, `build/tmp/spike9-oq3a.log` |
| Recovery token | `TASK-03-CREATION-IDENTITY-DEFECT-2026-11-22` |

This is **not** a reconciliation concern and must not be treated as one. Reconciliation can *detect* a collision; it cannot make creation safe.

---

## 1. Observed

`NewChat` → `generateChatID` → `MkdirAll(chatDir)` → `SaveChat`.

- `generateChatID` formats the current UTC time at **one-second resolution** (`persistence.go:118-125`).
- `NewChat` calls `os.MkdirAll(chatDir, 0755)` — which **succeeds on an existing directory** (`persistence.go:131-135`).
- `SaveChat` then overwrites `meta.json` and `history.json` unconditionally (`persistence.go:200-215`).
- `SaveChat` sets the global `LastChatID` to the chat it just wrote (`persistence.go:239`).

**Observed twice, after exactly two creations**, in disposable fixtures: two `NewChat` calls in the same second returned the same ID.

**Observed consequence**, from probe 9's first run: `SaveChatAs` inside the same second reused the source directory — `meta.json` 422 b → 425 b (the source renamed to "duplicate"), `history.json` `[]` → 503 b (the copy's content), and the source's own `history.ndjson` (375 b) left in place. One directory, two conversations, content and identity disagreeing.

---

## 2. Impact

| Aspect | Consequence |
|---|---|
| Identity | Two conversations share one directory; the later one silently adopts the earlier one's ID. |
| Data | The later creation's empty `history.json` overwrites the earlier one's snapshot; if the earlier one had appended, its `history.ndjson` survives under the new metadata. |
| Ownership | Because every write resolves through the global `LastChatID`, subsequent messages, edits, clears and final saves follow the pointer to whichever chat was created or saved last. |
| User visibility | None at creation. The damage appears later as a renamed or wrong conversation. |

**RC5-1 — Silent adoption is worse than failure.** The user sees no error at the moment the identity is lost, so nothing prompts them to stop.

---

## 3. Intended creation invariant

**RC5-2 — Every newly created conversation receives a unique ID.**

**RC5-3 — Directory creation must be exclusive.** A create that finds the directory already present is a collision, not a reuse.

**RC5-4 — An existing directory must never be silently adopted as a newly created conversation.**

**RC5-5 — A collision must fail loudly and leave the existing conversation untouched.**

**RC5-6 — No implementation authorization follows from identifying this.** These are the invariants a future repair task must satisfy, not a plan for performing it.

---

## 4. What reconciliation owes this defect

| Rule | Statement |
|---|---|
| RC5-7 | A scan reports a collision-shaped directory (metadata describing one conversation, history belonging to another, or an ID that two creations attempted) as an **anomaly**. |
| RC5-8 | A scan never repairs, renames, splits or re-creates it. |
| RC5-9 | A collision is not `Damaged` in the repair sense and must not be offered the same remedy; no reconciliation action may delete either party's data. |
| RC5-10 | The state table cannot distinguish "one conversation, partly damaged" from "two conversations, one directory" by content alone. **Unestablished** whether it can be distinguished at all; recorded as an open question rather than assumed detectable. |

---

## 5. Suggested separation

A distinct foundational task (proposed as **TASK-04**, not started, not authorized) covering: unique-ID generation, exclusive directory creation, collision failure semantics, and the ownership question of what happens to a directory whose metadata and log disagree.

**RC5-11 — Reconciliation is not a substitute for this fix.** A scanner that reports collisions every launch is a symptom report, not a repair.

---

## 6. Unestablished

| # | Item |
|---|---|
| RC5-A | Real-world frequency; the mechanism is confirmed, the exposure is not measured. |
| RC5-B | Whether two windows or a retry can reach the same second in practice. |
| RC5-C | Whether any existing conversation on disk was already created this way. Not investigated — that would require reading user data, which this investigation does not do. |
