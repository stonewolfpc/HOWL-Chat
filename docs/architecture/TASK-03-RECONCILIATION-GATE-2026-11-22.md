# TASK-03 — v2 amendment: the unresolved-transition mutation gate

| Field | Value |
|---|---|
| Purpose | Challenge the six v2 transitions against "can G+1 write while G is open?" |
| Status | **Amendment to row 554. Adds V2-51…V2-60. Contract work only.** |
| Supersedes | row 554 V2-16 (partly), V2-27 (qualified), V2-30 (extended) |
| Evidence | `build/tmp/spike12/probe12.go`, `build/tmp/spike12-gate.log` |
| Recovery token | `TASK-03-RECONCILIATION-GATE-2026-11-22` |

---

## 1. The question, answered against the writer, not the tables

**V2-59 — Observed: there is no gate, and the mixture is real.** Generation G with three committed records, an interrupted clear (marker written, replacement never performed), process restart, then an ordinary append:

```
AppendHistoryEntry error: <nil>
log: 0=gen-G-0  1=gen-G-1  2=gen-G-2  3=gen-2-first-message
LoadChat -> 4 messages, err=<nil>
```

Two generations share one file. The writer was never asked, and the marker file was untouched (mtime unchanged) — nothing reads it.

**Source-derived corroboration:** `AppendHistoryEntry` opens the log with `O_APPEND|O_CREATE` and performs no state check (`persistence.go:669-700`); there is no `pending`/`unresolved`/recovery logic anywhere in `internal/backend/chat` (search: no matches); and `App.startup` does nothing about chat persistence (`app.go:326-343`) — no reconcile, no recovery step before mutation is possible.

**V2-60 — The consequence for the contract.** The six transitions describe what the *scanner reports*. They never say the *writer* may not run, and v2 assumed the authorized operation completes or does not begin. It can begin and be abandoned mid-flight, which is exactly the observed state.

---

## 2. Why this cannot be caught afterwards

**V2-61 — Mixed-generation content is undetectable by inspection.** v2-16 already admitted the generation does not identify records. A log holding G's three records followed by one G+1 record is, byte for byte, indistinguishable from a single-generation conversation of four messages. There is no per-record stamp, no boundary, and no ordering evidence that separates them.

**V2-62 — Therefore prevention, not detection.** The contract's `Healthy` guarantee must rest on a gate that stops the mixture from arising. If the mixture arises, the scanner will call it Healthy and be right about everything except which generation it is reading.

**V2-63 — This makes V2-OQ-1 a correctness question.** Per-representation generation stamps would make attribution exact and the mixture detectable. Until then, correctness depends entirely on §3.

---

## 3. The gate

**V2-51 — While a generation transition is unresolved, no history mutation may append to or rewrite a representation belonging to it.** The writer must either (a) complete the authorized operation, or (b) refuse the mutation with a typed, user-visible failure. Refusal is not queueing: a queued mutation would be applied to a base whose generation is unknown.

**V2-52 — "Resolved" is defined per operation, not felt:**

| Operation | Resolved iff |
|---|---|
| first commit (0 → 1) | marker `(committed,1)`, or the append demonstrably never happened and marker is `(none,0)` |
| clear (closes G) | marker `(cleared,G)` **and** both representations are structurally valid **and** empty |
| append after clear (begins G+1) | the preceding clear is resolved (§ above), so the log was empty when the first G+1 record was written |

**V2-53 — G+1 may accept writes only after G is closed.** No append, rewrite, clear, save or duplication may touch a conversation whose clear is unresolved.

**V2-54 — The application finishes or refuses; the scanner does neither.** `RewriteHistory` is a full replace, so retrying an interrupted clear on startup is idempotent and safe to attempt. If it cannot complete, mutations stay refused until the user acts or the cause is removed.

**V2-55 — The scanner never resolves a transition.** It may report `Pending reconciliation` and promote a marker under V2-56; it may not write content, re-run a rewrite, or clear residue.

**V2-56 — Promotion requires generation-safe provenance.** `pending,G+1` may promote only when the transition that opened G+1 was resolved, so that the observed content is attributable to G+1. With the preceding clear unresolved, present bytes may belong to G — **never promote; report Pending**.

**V2-57 — Prefix rule, qualified.** `snapshot ⊂ log` counts as expected lag and supports Healthy **only when the log is structurally valid and the generation transition is settled** (V2-52). Otherwise the relationship is `Inconsistent` or the state is Pending, and the shorter representation is never treated as the base.

**V2-58 — The marker's own crash behaviour is tested independently, including generation-pair consistency.** After a death mid-save at `testHookBeforeReplace`, the reloaded record must be the previous **or** the new `(state, generation)` pair — never a state from one generation with the generation of another.

---

## 4. The six verification areas, answered

| Area | v2 as written | Verdict now |
|---|---|---|
| **Generation increment** | v2-15 says when the number changes; says nothing about write permission | **was incomplete** → closed by V2-52/V2-53 |
| **Clear crash recovery** | v2 reports `retracted-residue`; the app had no way to distinguish "clear finished" from "clear abandoned with stale snapshot" | **was report-only** → resolution is now defined by evidence (both empty), and the writer retries |
| **Append-after-clear** | v2's T6 assumed an empty log | **not guaranteed** — mixture observed → V2-53 + V2-56 |
| **Promotion** | v2-30 required validity + coexistence, not generation safety | **insufficient** → V2-56 |
| **Prefix relationship** | v2-27 allowed Healthy on prefix alone | **unqualified** → V2-57 |
| **DurableStore** | v2-13 covers previous-or-new value | **incomplete** → V2-58 adds pair consistency |

---

## 5. Tests added

- **VT-21** An append attempted while a clear is unresolved is refused (typed error), and the log is unchanged.
- **VT-22** A rewrite attempted while a clear is unresolved is refused.
- **VT-23** After a completed clear (both empty), the next append is accepted and the log contains only G+1 records.
- **VT-24** `pending,G+1` with an unresolved preceding clear → Pending, never promoted, even with valid non-empty bytes.
- **VT-25** `snapshot ⊂ log` with a valid log and a settled transition → Healthy; with an unsettled transition → not Healthy.
- **VT-26** Marker death at `testHookBeforeReplace` never yields a mixed `(state, generation)` pair.
- **VT-27** Retried clear on startup is idempotent: running it twice leaves both representations empty and the marker unchanged.
- **VT-28** A mixed-generation log (constructed as in the probe) is reported as indistinguishable — the test asserts the contract's limit, not a detection.

**V2-64 — VT-28 is deliberately a failing-by-design expectation**: it records that mixing cannot be detected, so a future per-record generation stamp changes that test rather than silently passing it.

---

## 6. Architecturally, where the gate must sit (not authorized)

Row 549 enumerated every content path — P2 append, P3 edit/delete, P4 clear, P5 final save, P6 duplication. A gate enforced in only some of them is not a gate. Whether that means one shared private helper or one predicate consulted by all five is an implementation question; **this document authorizes neither**, and states only that the five paths must agree.

**V2-65 — Contract work only.** No code, no gate implemented, no marker written, no migration, no deletion enablement. Row 554 remains a draft; this amendment is part of its review, not its adoption.

---

## 7. What closing V2-OQ-1 would take

| Option | What it buys | What it costs |
|---|---|---|
| **Strict gate only** (V2-51…V2-55) | prevents the mixture for the current writer architecture | unproved for any future writer that writes history outside these five paths |
| **Per-representation generation stamps** | makes attribution exact, mixing detectable after the fact | changes the on-disk format of both representations |

**V2-66 — The gate is sufficient for the current architecture only if every mutation path is enumerated and gated.** That inventory exists (row 549); whether it is closed is **Unestablished** — a sixth path added later would open the hole again. Per-representation stamps are the only option that does not depend on that inventory staying complete.
