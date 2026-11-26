# TASK-03 — DurableStore pair consistency, closed with evidence

| Field | Value |
|---|---|
| Purpose | Close the last open item in the row 556 assessment table |
| Status | **Adds V2-83…V2-88. Contract work only. No tests executed, none written.** |
| Evidence | `internal/durable/store.go`, `internal/durable/store_test.go` (existing, passing test read, not re-run) |
| Recovery token | `TASK-03-DURABLE-PAIR-2026-11-22` |

---

## 1. The four requirements, mapped

Row 556 was the call-graph audit. Each requirement in the review now has a clause and a test:

| Requirement | Clause | Test |
|---|---|---|
| Check and write inseparable; one enforcement point | V2-74, V2-75, V2-76 | VT-29, VT-30, VT-31 |
| Clear retry idempotent | V2-79 | VT-32, VT-24 |
| `cleared` plus both valid plus both empty | V2-80 | VT-33 |
| Inventory maintenance rule | V2-81 | VT-30, VT-36 |

---

## 2. Pair consistency: demonstrated, not assumed

**V2-83 — Observed (source + existing test).** `TestCrashBeforeReplaceLeavesPreviousDocument` spawns a real child that dies at `testHookBeforeReplace` (exit 9), i.e. between the manifest commit and the replacement, then asserts:

```
document after crash = "previous"          // intact bytes
Load -> "previous", Recovered(interrupted_commit_rolled_back)
```

The write sequence is: exclusive lock → write temp → stage generations → commit manifest → **die** → rename → promote. The live document's only mutation is a rename, so no reader ever sees a partial file, and the crashed commit is reported as a rollback rather than as a clean save or as corruption.

**V2-84 — Both crash windows are safe.** Death *before* the rename leaves the previous document intact and `Load` reports `Recovered` (§2). Death *after* the rename but before promotion leaves the new document installed with a matching manifest, so `Load` reports `Saved`, and the un-promoted staged files are swept by the next write (`store.go:533-538`). Neither window produces a document that is half of one commit and half of another.

**V2-85 — Therefore a mixed `(state, generation)` pair is impossible *by construction* — if, and only if, the marker is one document.** A whole document swapped by rename cannot contribute its `state` from one generation and its `generation` from another. The hazard row 555's V2-58 guarded against only exists if `state` and `generation` are written as separate records, separate keys or separate files. **Constraint: the marker is a single document; the pair is never split.** This is a schema rule, not a store limitation.

**V2-86 — `Recovered` is not authoritative, and the gate must not treat it as such.** `Result.Current()` is false for `Recovered` by design: the bytes came from a retained generation and may predate the last successful write. So a marker read that returns `Recovered` describes a pair that is internally consistent but possibly **stale relative to the content on disk** — for example `(committed,G)` while a clear's rewrite had already emptied the representations.

Consequence: **a `Recovered` marker read never counts as evidence that a transition is settled.** The gate refuses mutations; reconciliation reports `Pending reconciliation`, not `Healthy`, until the marker is re-established as current.

**V2-87 — The boundary needs a conditional replace, not read-then-save.** `Store.Load` is lock-free (it re-reads document and manifest up to three times to avoid reporting a commit-in-flight as corruption), and `Store.Save` takes the lock internally. A gate built as "Load the marker, decide, then Save" is exactly the caller-side check V2-74 forbids: another writer can commit between the two calls.

So the boundary requires **compare-and-swap on the marker under one lock hold** — "replace `(state, generation)` with the new pair only if it still equals the observed pair" — with the permitted/refused decision made inside that hold. `Save` as it exists cannot express that. This is a **design requirement on the eventual implementation**, not a defect in the store.

---

## 3. Assessment table, updated

| Finding | Before | Now |
|---|---|---|
| Cross-generation mixture | Reproduced defect | unchanged |
| Prevention by mutation gate | Specified, not implemented | specified; boundary requirements now include CAS (§2 V2-87) |
| Detection after mixing | Not reliably possible | unchanged |
| Per-representation generation stamps | Unestablished alternative | unchanged — still the only option independent of the inventory |
| **DurableStore pair consistency** | **test requirement added; not yet demonstrated** | **demonstrated for whole documents (§2), with the single-document constraint made explicit (V2-85)** |

---

## 4. Tests added

- **VT-37** A marker read returning `Recovered` causes the gate to refuse, and reconciliation to report `Pending` rather than `Healthy`.
- **VT-38** Conditional marker replace under the lock: a replace attempted against a stale observed pair is refused, with no interleaving window.
- **VT-39** After a crash at the replacement point, the marker is either the previous whole pair or the new whole pair, never a synthesis of both.

VT-39 is the marker-level counterpart of an existing store test; it is not a new proof of rename atomicity.

---

## 5. Where the contract stands

**V2-88 — Row 555 remains the amendment; row 556 the call-graph audit; this row closes the last open assessment item.** Remaining items in the table are unchanged and each names its own reason: the mixture is a reproduced defect, the gate is specified and unimplemented, detection is not reliably possible, and per-representation stamps are unestablished.

**V2-89 — Contract work only.** No gate, lock, marker, CAS or recovery implemented; no test written or executed beyond reading an existing one. The store's passing status is **Source-derived** from its test source, not re-observed in this session.
