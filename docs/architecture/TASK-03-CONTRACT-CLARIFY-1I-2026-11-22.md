# TASK-03 — clarification 1i: execution context, and the crash window that loses continuity

| Field | Value |
|---|---|
| Purpose | Resolve the `removed` + K1 matrix contradiction; tighten three claims |
| Status | **Contract text. Nothing adopted. No implementation authorized. Deletion unavailable** |
| Relationship | **Additive.** None of rows 528–537 is edited |
| Supersedes | 1h U-26 (`removed` row, context-free), U-21 (as a statement about reconciliation's *outcome* rather than a note), U-31 (scope), U-19/U-20 (as proof rather than design intent) |
| Recovery token | `TASK-03-CONTRACT-CLARIFY-1I-2026-11-22` |

Prefixes: `V-` normative, `VT-` acceptance test. Labels: **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Superseded**.

---

## 1. The contradiction, stated exactly

**V-01 — 1h contradicts itself.** 1h said two things that cannot both hold without a third dimension:

- U-26's matrix row `removed` + K1 → `complete`, unconditionally;
- U-21: after a restart the acquisition handle is gone, continuity cannot be established, and **reconciliation never writes `complete`**.

**V-02 — The unresolved case.** The operation writes the durable `removed` marker, then dies before writing the terminal receipt. On restart the journal reads `removed`; the target open may well produce K1. The matrix as written permits `complete`. U-21 forbids it. Both were presented as settled. They are not.

**V-03 — Why it is not a wording problem.** `complete` asserts absence *for the acquired identity, observed under continuity*. The `removed` marker records a **past** observation made while a handle was held. It is evidence that a removal was attempted and that its call returned; it does not recreate the handle, and it cannot make the identity continuously identified across the crash. Treating the marker as sufficient would let a recovery path claim a stronger continuity than the live operation ever held. That is precisely the failure mode this contract exists to prevent.

**V-04 — The fix is a dimension, not a softened sentence.** Every observation is now qualified by **execution context**:

| Context | Meaning |
|---|---|
| `original` | the process that opened the volume at acquisition, still holding that handle, verifying in the same process lifetime |
| `reconciliation` | any later process; the acquisition handle is gone |

The matrix is indexed by `(state, context, observation)`. A state alone no longer determines an outcome.

---

## 2. Corrected matrix

**V-05 — 4 states × 2 contexts × observation classes.**

| State | Context | K1 absent | K2 present (enum. incomplete) | K2 present (count equal) | K2 present (count differs) | K3/K4/K5 unverifiable |
|---|---|---|---|---|---|---|
| `prepared` | `original` | `absent_after_acquisition`/`unknown` | `aborted` | `aborted` | `aborted` | `unverifiable`/`unknown` |
| `prepared` | `reconciliation` | `absent_after_acquisition`/`unknown` | `aborted` | `aborted` | `aborted` | `unverifiable`/`unknown` |
| `pending` | `original` | `absent_after_acquisition`/`unknown` | `aborted` | `aborted` | `aborted` | `unverifiable`/`unknown` |
| `pending` | `reconciliation` | `absent_after_acquisition`/`unknown` | `aborted` | `aborted` | `aborted` | `unverifiable`/`unknown` |
| `removing` | `original` | `absent_unverified`/`unknown` | `not_completed`/`unknown` | `not_completed`/`unknown` | `count_delta_unattributed`/`unknown` | `unverifiable`/`unknown` |
| `removing` | `reconciliation` | `absent_unverified`/`unknown` | `not_completed`/`unknown` | `not_completed`/`unknown` | `count_delta_unattributed`/`unknown` | `unverifiable`/`unknown` |
| `removed` | `original` | **`complete`**/`unknown` | `indeterminate`/`state_inconsistent` | same | same | `unverifiable`/`unknown` |
| `removed` | `reconciliation` | **`completion_unconfirmed`**/`unknown` | `indeterminate`/`state_inconsistent` | same | same | `unverifiable`/`unknown` |

**V-06 — Only one cell in the table can produce `complete`.** `removed` + K1 + `original`. Every other K1 cell produces a non-complete outcome, and no cell in `reconciliation` produces `complete`.

**V-07 — `prepared`/`pending` are context-independent.** Their K1 outcome, `absent_after_acquisition`, claims only that the identity was not found by an operation that never started removing. Context adds nothing there, which is why the two rows are identical — recorded rather than collapsed, because a future change to either context must not silently alter them.

---

## 3. Crash window 11

**V-08 — The window.** 1f defined ten crash windows; 1h updated its conditions without adding one. The case in V-02 was not among them:

> **Window 11 — `removed` durable, terminal not yet written, process dies.**

Durable state at the crash: `removed` marker present, `own_removal_observed: true` recorded, no terminal receipt, handle gone.

**V-09 — Reconciliation outcome for window 11.**

| Observation at restart | Outcome | Terminal? |
|---|---|---|
| K1 `identity_not_found` (U-02 satisfied) | **`completion_unconfirmed`** | yes |
| K2 present | `indeterminate` / `state_inconsistent` | yes |
| K3/K4/K5 unverifiable | `unverifiable` | **no** — entry stays open, may be re-examined |

**V-10 — `completion_unconfirmed` defined.**

| Field | Value |
|---|---|
| Asserts | the identity was **not found** by a post-restart observation, and this operation's removal call returned success before the crash |
| Does **not** assert | continuous identification across the crash; that the absence was caused by this operation; attribution of any kind |
| `own_removal_observed` | `true` — from the durable marker, a past fact |
| `attribution` | `unknown` |
| `continuity` | `none_after_restart` |
| Mutability | **terminal and immutable**; never upgraded to `complete` |

**V-11 — Why a new outcome rather than reusing `indeterminate`.** `indeterminate` already means "state and observation are inconsistent" and is paired with `state_inconsistent`. Window 11 is not inconsistent: everything observed agrees with the record. It is *insufficient* — absence was seen, continuity was not. Collapsing it into `indeterminate` would lose the distinction between "the evidence conflicts" and "the evidence is incomplete", and would leave a future reader unable to tell whether the volume was ever verified at all.

**V-12 — This is not a new evidence kind.** `completion_unconfirmed` is a terminal outcome value within the existing receipt. The five ledger kinds and two evidence kinds are unchanged; no second history is introduced; the receipt stays immutable.

**V-13 — Reconciliation must not write `removed`.** Only the operation, in-process, writes `removed`. Reconciliation may read it and must never overwrite it, and must never re-run removal from it. The marker enables a *report*, not an *action*.

**V-14 — Existing terminal receipts.** If window 11 is entered and a terminal already exists — impossible in the defined order, but reachable if a receipt is written by an earlier run — the receipt is returned **unchanged** under the 1g post-terminal rules: `post_terminal_identity_present`, `post_terminal_volume_mismatch`, `post_terminal_unverifiable`; nothing deleted, nothing rewritten, anomaly transient.

**V-15 — The cost, stated.** A crash between the `removed` marker and the terminal permanently costs the completion claim for that operation, even when the identity is provably absent. There is no recovery that restores it, because the handle cannot be restored. This is accepted: the alternative is a recovery path asserting continuity it never had.

---

## 4. Tightened claims

**V-16 — U-31's scope, bounded.** 1h's U-31 read as a general property. It is not. Restated:

> Within **one continuously held volume object** and **the exact acquired directory identity**, a reuse of that identity can weaken attribution but cannot manufacture a false absence.

Out of scope, and not covered by that claim: a **reopened locator**, a **replacement volume**, a **post-restart observation**, or any identity re-derived by path. Those make no absence claim at all (V-05) — because the claim was never available to them, not because the failure mode is absent.

**V-17 — U-19/U-20 are design intent, not proof.** 1h described handle continuity in a way that read as established behavior. It is not:

- That the destructive call reaches the object identified through the acquired handle is **design intent**, gated on E-27 (Unestablished).
- **A held volume handle alone guarantees nothing if the directory is later resolved by path.** The directory must be opened *and identified* through that same volume handle — by-ID open on it — and the removal must act on that opened object. Holding a volume handle while resolving `C:\...\conversations\x` gives the same guarantee as re-resolving a path: none.

**V-18 — No claim about Windows behavior.** No clause asserts that a held handle prevents locator re-pointing, or that the mount manager treats a duplicate GUID in any particular way. Both are E-26/E-27, unestablished, and the contract is written to be correct **either way**: if E-27 fails, `complete` becomes unreachable and `completion_unconfirmed`/`unverifiable` carry the whole space. That is the fail-closed direction.

**V-19 — `count_delta_unattributed` restrictions, made explicit.** Retained verbatim from 1h and extended:

| Observation | Must not be read as |
|---|---|
| equal counts | unchanged contents |
| lower count | deletion by this operation |
| higher count | that nothing was removed |
| incomplete enumeration | any count-delta result (it yields `not_completed`/`unverifiable`) |
| `count_delta_unattributed` | "partially deleted by Howl" |

**V-20 — The UI wording rule is normative.** No user-facing string may translate `count_delta_unattributed` (or any state in V-05) into an authorship claim: not "partially deleted by Howl", not "Howl removed some files", not "deleted (partial)". This is a contract requirement because a precise state name can be technically correct while the surface gives it a misleading reading. Enforcement is a lint on the string table, specified in VT-10.

---

## 5. Preserved

**V-21 — Carried unchanged.** Count-only baseline; no manifest; the five-property separation; FAT/exFAT/SMB refusal; `identity_unsupported`; U-02's five preconditions and the control open; the six failure classes; `absent_after_acquisition`; `absent_unverified`; `aborted` asserting only that this operation performed no removal; immutable terminal receipts; transient post-terminal anomalies; D-10 unchanged; G-1 an engineering gate.

**V-22 — E-24…E-27 remain unestablished**, unchanged in content:

| # | Status |
|---|---|
| E-24 | USN proof of `reuse_excluded` — unmet, so no Phase 1 attribution |
| E-25 | The exact target status per filesystem and the control procedure — **the third-party filter note is source evidence, not empirical validation of Howl's call sequence** |
| E-26 | Duplicate-GUID behavior / `identity_ambiguous` |
| E-27 | Held-handle continuity for the call sequence — now also covers: the directory must be identified *through* that handle |

**V-23 — No new owner decision.** `completion_unconfirmed` is a consequence of the evidence rule, not a choice offered to the owner. No weaker path-based fallback is introduced anywhere.

---

## 6. Acceptance tests

- **VT-01** Window 11 fixture (`removed` durable, no terminal, simulated restart) with K1 → `completion_unconfirmed`; `complete` is refused.
- **VT-02** `removed` + K1 + `original` context with a live acquisition handle → `complete` permitted; `attribution: unknown` required.
- **VT-03** Context gating: any path producing `complete` outside `removed` + K1 + `original` fails.
- **VT-04** Reconciliation never writes `complete` for any state or observation (exhaustive over V-05).
- **VT-05** Reconciliation never writes `removed` and never re-runs removal from a marker.
- **VT-06** `completion_unconfirmed` is terminal, immutable, and is never upgraded to `complete` by a later run.
- **VT-07** `completion_unconfirmed` carries `own_removal_observed: true`, `attribution: unknown`, `continuity: none_after_restart`.
- **VT-08** Window 11 with K2 → `indeterminate`/`state_inconsistent`; with K3–K5 → `unverifiable`, no terminal.
- **VT-09** Existing terminal at window 11 entry → returned byte-identical; no post-terminal rewrite.
- **VT-10** String-table lint: every state maps to wording that makes no authorship claim; "partially deleted by Howl" and equivalents rejected for `count_delta_unattributed`.
- **VT-11** The four count observations of V-19 each produce no authorship reading, tested against the real string table.
- **VT-12** No ledger or evidence kind is added by V-10's outcome (kinds remain five and two).

---

## 7. Verification and authorization

Rows 528–537 remain byte-identical; 1h (`c247c9d34972…`) unchanged. 1i registered as row 538; counts updated; no historical entry altered. **Zero Go files modified; no ledger created; no migration; no deletion enabled.** No test executed.

TASK-03 remains in contract preparation. Not accepted, not closed. Deletion remains **unavailable**.
