# TASK-03 — OQ-7: enforcement evidence

| Field | Value |
|---|---|
| Purpose | What build-bound evidence could prove the mutation boundary is enforced |
| Status | **Read-only design investigation. No code read was executed, no test run. Nothing authorized.** |
| Method | Source inspection of the lock primitive, the write paths, and the absence of any instance or self-test mechanism |
| Recovery token | `TASK-03-OQ7-ENFORCEMENT-2026-11-22` |

Prefix: `O7-` proposals. V3.1 (row 560) is unchanged; OQ-6 (row 562) is superseded in one clause, noted in §8.

---

## 1. The central issue: atomic is not truthful

**O7-01 — An atomically replaced record can still be an incorrect assertion.** DurableStore guarantees the marker is *whole*, never that it is *true*. Truth depends on three things it cannot see:

| Layer | Guarantee | Depends on |
|---|---|---|
| the record is whole | rename exclusivity | `DurableStore` (present) |
| the record is truthful | the writer held exclusivity across the whole transition | the gate (absent) |
| the writer's claim is complete | no mutation path bypasses the boundary | build-bound proof (absent) |

Break any link and the settlement field is well-formed evidence of nothing. **This is why OQ-7 precedes OQ-6 in practice, not only in sequence.**

---

## 2. Q1 — Authority: what proves every mutation path uses one boundary?

**O7-02 — The premise is currently false, so no certificate can exist.** Row 556 counted four writers inside persistence, at least twelve entry points, and three writers outside `PersistenceManager` entirely. There is no single boundary to certify.

**O7-03 — Three candidate proofs, and only two are proof:**

| Candidate | Strength | Failure mode |
|---|---|---|
| **Startup self-test** (exercise each path, assert refusal) | regression evidence only | blind to any path not exercised; a new bypass is invisible |
| **Static build assertion** (no code outside the boundary package opens a conversation file for writing) | mechanical, whole-subgraph | bypassable by whatever the checker does not model |
| **Unforgeable capability type** (a write token obtainable only under the lock, unexported field, single constructor in the boundary package) | compiler-enforced; cannot be constructed elsewhere | same-package code can forge; package boundary must be the enforcement unit |

**O7-04 — Authority must be the capability type, checked by the compiler, plus the static assertion as a second net.** A self-test cannot be a certificate.

**O7-05 — A successful startup self-test never certifies enforcement** while a later mutation path can bypass the tested boundary. This is the rule as stated, and it follows from the table above: the self-test's coverage is exactly the set of paths someone remembered.

**O7-06 — The claim must be scoped or it is false.** A capability token proves *history-representation* writes went through the boundary. It cannot prove anything about `memories.json`, `auto_capture_config.json`, `continuity_config.json` (three external writers, row 556) or `WipeAllData`'s root rename. Enforcement's public claim is therefore: **"no write to a history representation occurs outside the boundary, in this build"** — nothing broader.

---

## 3. Q2 — Exclusivity: does the boundary serialize instances and prevent check-then-write?

**O7-07 — The primitive is adequate (Source-derived).** `internal/durable/lock_windows.go`: `LockFileEx` with `LOCKFILE_EXCLUSIVE_LOCK|LOCKFILE_FAIL_IMMEDIATELY` on an open handle, polled to a deadline, `ErrBusy` on timeout. `lock_other.go`: `flock(LOCK_EX|LOCK_NB)`. Both are **OS-level, cross-process, fail-fast** locks.

**O7-08 — Two properties follow, and both matter for Q4:** an OS lock is released when the process dies, so there is **no stale-lock problem** and no lock-file reclamation logic to get wrong; and it cannot be withdrawn by another process while held.

**O7-09 — The composition is absent.** `PersistenceManager` takes no lock on any history write (row 556), and the marker CAS (V3-29) does not exist. The primitive would serialize two instances *if* both used it on the same path; today neither does.

**O7-10 — The lock is advisory, so exclusivity is only as good as O7-04.** A writer that skips the lock is unaffected by it. This is not a weakness of the primitive; it is why the capability type, not the lock, is what makes the property provable.

**O7-11 — The lock's scope is per document path, not per conversation.** A boundary spanning marker and both representations needs either one lock over the conversation directory or a documented order with a conversation-level lock. Not decided here; named as a design choice.

---

## 4. Q3 — Lifecycle: can enforcement be withdrawn, bypassed, or become unavailable?

| Question | Answer today |
|---|---|
| **Withdrawn while held?** | No. An OS lock on an open handle survives until release or process death (O7-08). |
| **Withdrawn between operations?** | Trivially — there is nothing to withdraw; no process holds a standing claim. Enforcement is per-critical-section, not a mode. |
| **Bypassed while running?** | Yes, completely. Any code path may write without the lock (O7-09, O7-10). |
| **Unavailable while running?** | Yes. Lock acquisition fails on `ErrBusy`, on open failure, on a read-only or disconnected volume. Mutations must then refuse. |

**O7-12 — Unavailability is the dangerous case, and the product currently converts it into silence.** Refusals propagate as ordinary errors, and at least four call sites discard them (`_ = a.performSave()` — `app-characters.go:242,257`, `app-chat.go:137,148`). An unavailable boundary today produces *silent non-persistence*, which is worse than either a refusal or a write.

**O7-13 — Therefore enforcement is not a mode that can be switched off; it is a property of each critical section.** `gate_enforced` is true iff every mutation in this build passes through such a section — a build-level fact, not a runtime flag (V3-36).

---

## 5. Q4 — Crash behaviour across the whole transition

Holding exclusivity from the first content byte to the settlement CAS, the death points are:

| # | Death after | Marker | Content | Row | State |
|---|---|---|---|---|---|
| 1 | acquiring the boundary, before any write | unchanged | unchanged | — | unchanged; lock released by death, no stale lock |
| 2 | rewriting the log | `(cleared,G)`, unsettled | log new, snapshot old | 6 | Unavailable `inconsistent` |
| 3 | rewriting the snapshot | `(cleared,G)`, unsettled | both void | 13 | Pending `unattributed` |
| 4 | settlement CAS | `(cleared,G,settled=G)` | both void | 11 | Healthy `closed` |
| 5 | mid-append | `(pending,G+1,settled=G)` | torn log | 7 | Damaged `truncated` |
| 6 | append, before `(committed,G+1)` | `(pending,G+1,settled=G)` | valid content | 16/18 | Healthy (gated) or Pending |
| 7 | `(committed,G+1)` | committed | valid | 17/19 | as evidence permits |

**O7-14 — Every death point is fail-closed, and none reads as Healthy-but-wrong.** Point 2 is the only one producing `Inconsistent`; point 3 is the case that matters most: **both representations are void and valid, yet the clear is not recorded as settled.** A later observer must not infer settlement from that condition (§8).

**O7-15 — Retry is safe at every point** because `RewriteHistory` is a full replace: re-running the clear from point 2 or 3 finishes it, and the settlement CAS from point 3 is idempotent (V3-33).

**O7-16 — The OS lock releases on death (O7-08), so the retry never waits on a ghost holder** and no lock-file timeout heuristic is needed. This is the one part of the enforcement story that is already sound.

---

## 6. What enforcement can and cannot certify

**O7-17 — It certifies:** for content written by **this build**, through **this boundary**, while **holding exclusivity**, the recorded settlement was asserted under the hold; and no history write in this build occurred outside the boundary (O7-06).

**O7-18 — It cannot certify:** content written before the boundary existed (V3-37, no backfill); non-history files in the conversation directory (O7-06); root-level operations; or anything about a build whose capability proof is absent.

**O7-19 — `settled_at` never influences a state.** A timestamp proves neither settlement, nor ordering, nor writer exclusivity. It is audit metadata and no row of §4 may read it.

---

## 7. Correction to OQ-6: the precondition is the hold, not the observation

**O7-20 — This supersedes OQ-6's S-06.** Settlement must not advance because both representations are empty and structurally valid. That condition establishes their *observed state*, not that a competing writer could not have modified them during the operation. The precondition for writing `settled_generation` is: **the writer has held the boundary continuously from before the first content byte through the settlement CAS.** Observation alone never authorises it — and a scanner observing void, valid representations may not record settlement either (V3-56, S-13).

**O7-21 — Consequence:** settlement is a statement *about a writer's hold*, recorded in the same CAS as the pair (S-07). Its truth is exactly O7-04's truth — no more.

---

## 8. Verdict: can the single-marker settlement fields ever become trustworthy?

**O7-22 — Conditionally yes, and only under all five:**

1. one boundary, enforced by an unforgeable capability type (O7-04);
2. a build-time assertion that no other code opens a history representation for writing (O7-06 scoped);
3. OS-lock exclusivity held continuously across the whole transition (O7-11, O7-16);
4. settlement written only under that hold, never by observation (O7-20);
5. unavailability surfaced as refusal, never discarded (O7-12).

**O7-23 — Absent (1) and (2), the settlement field is an assertion by an unverified writer** — atomic, recoverable, and evidentially worthless. Absent (5) it is worse than worthless, because the failure is silent.

**O7-24 — Therefore OQ-6 is not adoptable yet**, and not because its storage design is weak: because storage design cannot supply authority. OQ-7's five conditions are the price of the field meaning anything.

**O7-25 — Nothing here is authorized.** No capability type, no lock composition, no marker field, no self-test, no static assertion, no code change. This investigation is read-only and its conclusions are proposals.

---

## 9. Tests enforcement would require

- **O7-T1** A history write attempted without the capability token does not compile (fixture that must fail to build).
- **O7-T2** The static assertion fails when a new file opens a history representation for writing outside the boundary.
- **O7-T3** Two processes: the second reports `Busy` and refuses; neither writes across an unresolved transition.
- **O7-T4** Death at each of §5's seven points leaves the row named there, and none yields Healthy-with-wrong-content.
- **O7-T5** Retry from points 2 and 3 completes the clear; settlement is recorded once.
- **O7-T6** Settlement is never written from an observation without the hold — instrumented assertion.
- **O7-T7** Lock unavailability surfaces as a caller-visible refusal at every entry point; no `_ =` discard remains.
- **O7-T8** A self-test that passes while a bypassing path exists does not set `gate_enforced` (the rule of O7-05, as a test).
