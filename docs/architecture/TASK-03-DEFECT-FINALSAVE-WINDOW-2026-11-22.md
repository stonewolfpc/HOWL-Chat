# TASK-03 — defect record: the final-save representation window

| Field | Value |
|---|---|
| Purpose | Record the shutdown write-ordering defect as its own defect, separate from re-entrancy and lock scope |
| Status | **Defect record only. Not repaired. Not authorized.** |
| Discovered | OQ-7c call-graph inspection (row 568, §O7-109) |
| Owner | **writer ordering and crash consistency** — not reconciliation, not lock scope |
| Recovery token | `TASK-03-DEFECT-FINALSAVE-WINDOW-2026-11-22` |

Prefix: `FD-`.

---

## 1. Shape

**FD-01 — Observed (Source-derived).** `App.performFinalSave` (`app-chat.go:314-352`) runs, in order:

1. `performSave` → `LoadChat` → `SaveChat`, which writes `meta`, **`history.json` from app state**, companions, and the active-chat pointer;
2. `LoadChat` again;
3. `RewriteHistory`, which writes **`history.ndjson` and then `history.json`** again.

**FD-02 — Therefore every clean shutdown passes through a state where the snapshot is newer than the log.** After step 1 the snapshot holds the final history while the log still holds the pre-save log. That is `log ⊂ snapshot` — **row 6, `inconsistent`** — which the decision table reports as Unavailable.

**FD-03 — Normally the window is closed by step 3 within milliseconds.** It becomes persistent only if the process dies between steps 1 and 3.

**FD-04 — Why it is still a defect, not a tolerated transient:**

- it occurs on **every** shutdown, so any crash during shutdown has a real chance of landing inside it;
- it is not a race — it is the *designed* order of two writes that both rewrite the same pair;
- recovery is not automatic: on the next launch the loader prefers a non-empty log, so the conversation silently reverts to the older log while the newer snapshot sits beside it, and reconciliation reports `Inconsistent` for a conversation the user did not touch.

---

## 2. What reconciliation must do — and must not

**FD-05 — The contract's behaviour is correct and unchanged.** Row 6 reports `Unavailable / inconsistent`. It must **not** pick the newer representation, prefer the log because the loader does, or reconstruct the intended history (V3-41, V3-47).

**FD-06 — This defect is not a reason to amend the table.** It is a writer-ordering defect that manifests as a state the table already handles. Absorbing it into reconciliation would repeat the mistake row 558 found: making the scanner compensate for what only the writer can prevent.

**FD-07 — The repair direction is a single write of both representations in one operation** (OQ-7c §O7-113's decomposition 1), not a reconciliation rule.

---

## 3. Boundaries

**FD-08 — Not the same as OQ-3c-b.** `RewriteHistory` resurrection concerns a clear followed by a crash. This window concerns a *save* followed by a *rewrite* during shutdown. They share the log/snapshot pair and nothing else.

**FD-09 — Not a lock-scope issue.** Taking a per-conversation lock does not close the window; it only prevents *other* writers from interleaving while the window is open. Ordering is the fix; exclusion is a separate requirement.

**FD-10 — Not repaired here, and no repair is authorized.** No change to `performFinalSave`, `SaveChat`, `RewriteHistory`, or their ordering follows from this record.

**FD-11 — The finding stays four-valued on reachability:** the path is production, reached by the ordinary shutdown sequence, and observed by reading the code — **not** reproduced in this session.
