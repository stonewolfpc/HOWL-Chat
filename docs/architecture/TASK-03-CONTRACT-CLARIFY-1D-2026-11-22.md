# TASK-03 — clarification 1d: post-acquisition absence, removal guarantees, absent-target identity

| Field | Value |
|---|---|
| Purpose | Close the final three issues in the deletion authority |
| Status | **Contract text. Nothing adopted. No implementation authorized** |
| Relationship | **Additive.** None of rows 528–532 is edited |
| Supersedes | N-07, N-08 (directory handling), N-11, N-15; narrows the claim made by N-17 |
| Recovery token | `TASK-03-CONTRACT-CLARIFY-1D-2026-11-22` |

Prefixes: `P-` normative, `PT-` acceptance test.

---

## 1. `absent_prior` was a false historical claim

**Finding, accepted, and it is my error.** N-07 inferred `absent_prior` from `pending` + absent directory. But acquisition **observed the directory present** before writing the entry. So by the time the entry is `pending`, absence is not prior — it happened after acquisition and before our removal call. Labelling it `prior_absence` asserts a history that never occurred, which is precisely the failure mode this contract exists to prevent.

**P-01 — New outcome: `absent_after_acquisition`.**

> When the acquisition observed the directory present and the entry state is `prepared` or `pending`, an absent directory at reconciliation yields **`absent_after_acquisition`**, with `attribution: unknown` and `removal_attempted: false`.
>
> It states three facts and no more: the directory was present at acquisition; it is absent now; this operation's removal call was never made. It does **not** say who removed it, when, or why.

**P-02 — `prior_absence` is reserved for one case only.** The acquisition observation itself, before any entry or intent exists (K-09): the target was absent when the operation tried to begin. `absent_prior` may not be derived from any later observation.

**P-03 — Complete inference table. Replaces N-11.**

`removal_attempted` is derived from the entry state, never observed: `prepared`/`pending` → false; `removing`/`removed` → true.

| Entry | Intent | Removal attempted | Directory now | Outcome | Attribution |
|---|---|---|---|---|---|
| none | — | false | absent at acquisition | `absent_prior` | `prior_absence` |
| `prepared` | absent | false | any | `abandoned_no_intent`, entry removed, no receipt | — |
| `prepared` / `pending` | present | **false** | absent | **`absent_after_acquisition`** | `unknown` |
| `prepared` / `pending` | present | false | intact | `aborted` | — |
| `prepared` / `pending` | present | false | partly emptied | `partial`, `removal_attempted: false` | `unknown` |
| `removing` | present | true | absent | `complete` | `unknown` |
| `removing` | present | true | intact | `aborted` | — |
| `removing` | present | true | partly emptied | `partial` | `unknown` |
| `removed` | present | true | absent | `complete` | `operation` |
| `removed` | present | true | present | `indeterminate`, `state_inconsistent` | `unknown` |
| `pending`/`removing`/`removed` | absent | — | any | `indeterminate`, `receipt_missing`; entry preserved | `unknown` |
| any | any | — | — | terminal exists → return it, then cleanup | — |

**P-04 — The partial rows change too.** A partly emptied directory observed while the removal was never attempted is an **external** change: `partial` with `removal_attempted: false` and no authorship. The earlier table's silence on this case would have let a partial effect be read as the operation's work.

**Alternatives considered.**
- *Keep `absent_prior` and document the nuance.* Rejected: the receipt field is the claim; a footnote does not stop a consumer from reading `prior_absence` as "was already gone".
- *Report `indeterminate` for post-acquisition absence.* Rejected: it discards two facts the evidence does establish. The new outcome is strictly more informative *and* claims less than `prior_absence`.
- *Reuse `complete` with `attribution: unknown`.* Rejected: `complete` asserts the operation removed it, which the `removal_attempted: false` barrier contradicts.

**PT-01** — Present at acquisition, absent at reconciliation, entry `pending` → `absent_after_acquisition`/`unknown`, never `absent_prior`.
**PT-02** — Any receipt carrying `prior_absence` has no entry and no intent.
**PT-03** — `prepared`/`pending` with a partly emptied directory yields `partial` with `removal_attempted: false` and no authorship.
**PT-04** — The table is exhaustive: every (state, presence) pair maps to exactly one row; an unlisted pair fails the test.

---

## 2. Removal cannot be made TOCTOU-safe by re-checking, and the contract must say so

**Finding, accepted.** I described identity re-checking as closing the window. It narrows it. A check followed by a path-based `RemoveAll` re-resolves the path at call time, so a replacement between the second check and the call operates on whatever now sits at that path. That is a correctness boundary, not a polish item.

**P-05 — The authority boundary, restated.** PathService provides path classification and canonical identity (TASK-01). It does **not** provide authorised safe deletion, and resolving identity does not create a deletion capability. Nothing in this contract may be read as promoting one into the other.

**P-06 — Two removal bases, and only one of them supports an identity claim.**

| Basis | Mechanism | What a `complete` receipt may claim |
|---|---|---|
| **`handle_identity`** | Every entry removed through a handle opened on the acquired identity; children enumerated and removed by handle; no path re-resolution after the identity is opened | The acquired physical directory was removed |
| **`path_recheck`** | Canonical path, identity re-checked immediately before the call, then `RemoveAll` by path | **The path is absent and our removal call returned success.** Not that the acquired directory was the one removed |

**P-07 — Phase 1 records its basis.** The receipt carries `removal_basis`. Where the basis is `path_recheck`, no field, log line, UI string or document may state that the acquired directory was removed.

**P-08 — What `identity_changed` proves.** Only that a mismatch was observed at a check. It is not evidence that replacement races are prevented, and it may not be described as such anywhere in the UI, docs or diagnostics.

**P-09 — The residual risk, stated as a limitation rather than mitigated.** Under `path_recheck` a directory removed and replaced between the check and the call can cause the deletion to operate on a different directory that now occupies the same path. The receipt is honest about the basis but honesty does not undo the effect. Mitigation requires `handle_identity`, which Phase 1 does not have and this contract does not claim.

**P-10 — Owner decision required.** Accept `path_recheck` for Phase 1 with this limitation recorded, or require `handle_identity` before any Phase 1 deletion is enabled. This is the tenth owner decision and it is not mine to take.

**Alternatives considered.**
- *Assert that re-checking is sufficient.* Rejected: it is false, and the review caught it.
- *Specify `handle_identity` as the Phase 1 mechanism.* Rejected as a unilateral choice: it is unbuilt, unproven here, and materially larger than the rest of the deletion path — but it is named as the only basis that supports the stronger claim.
- *Drop deletion from Phase 1.* Rejected: retention Option A requires it; the honest path is to bound the claim, not to remove the feature.

**PT-05** — No `path_recheck` receipt or message asserts removal of the acquired directory (string and structured assertions).
**PT-06** — Every `complete` receipt carries a `removal_basis`.
**PT-07** — `identity_changed` is documented and surfaced as an observation, not a guarantee.
**PT-08** — A simulated replacement between the final check and the call is recorded as a known limitation, and the resulting receipt still reports only path-level facts.

---

## 3. Identity for a target that does not exist

**Finding, accepted.** N-15 let an absent target yield `absent_prior` without defining what identity means when there is nothing to resolve. A directory that does not exist has no final physical identity.

**P-11 — Two identity representations.**

| Target | Identity | Established by |
|---|---|---|
| **Present** | Physical identity of the directory itself | PathService final-path resolution by handle (existing TASK-01 capability) |
| **Absent** | **Composite: (canonical parent identity, exact child name)** | Parent resolved by PathService; child name validated as a single component |

**P-12 — Composite identity rules.**

- The parent must resolve; if it does not, the target's identity is unresolved.
- The child name is exactly one component: no separators, no `.` or `..`, no trailing dot or space, no reserved device name on platforms that have them.
- Comparison follows the filesystem's own case rule: on a case-insensitive filesystem two spellings of the same child are one identity; on a case-sensitive one they are two.
- The composite identity is used for **classification only** — distinguishing genuine absence from other failures. It is never stored, and it never becomes an ownership key, because an absent target takes no ownership and writes no entry.

**P-13 — Failure taxonomy. Replaces N-15.**

| Observation | Result |
|---|---|
| Parent resolves; child absent | `absent_prior` — a claim about the observation instant |
| Parent resolves; child absent, then present on re-observation | re-enter as a present target (one retry) |
| Still oscillating after the retry | `identity_race`; fail closed; nothing written |
| Parent unresolvable | `identity_unresolved` (`parent_unresolved`); no receipt |
| Parent present but inaccessible | `identity_unresolved` (`parent_inaccessible`); no receipt — **not** `absent_prior` |
| Target present but its own resolution fails | `identity_unresolved` (`target_unresolved`); no receipt, no ownership |

**P-14 — No string fallback.** If neither representation can be established, the request fails closed. Ownership is never acquired on a raw path string, and no comparison falls back to string equality.

**P-15 — Ownership is unaffected by composite identity.** Because an absent target takes no entry, two concurrent requests for the same absent path may each write `absent_prior`. Both claim only "absent at the instant I observed"; neither claims removal, so there is no ambiguity to resolve. This is deliberate: absence needs no owner.

**Alternatives considered.**
- *Resolve the absent path's longest existing prefix and hash it.* Rejected: a hash is not an identity, and the salt problem returns.
- *Treat every failed resolution as absence.* Rejected: it converts an inaccessible parent into a deletion claim.
- *Give absent targets ownership entries.* Rejected: it creates ownership for an operation that performs no removal, and leaves an entry that nothing can justify ending.

**PT-09** — Genuine absence, inaccessible parent and unresolvable parent produce three distinct results.
**PT-10** — A target appearing during classification is retried once, then `identity_race`.
**PT-11** — Case-variant spellings of an absent child compare equal on a case-insensitive filesystem and distinct on a case-sensitive one.
**PT-12** — No absent-target identity is persisted; no code path compares ownership keys as strings.
**PT-13** — Two `absent_prior` receipts for one absent path are permitted and neither is treated as a duplicate deletion.

---

## 4. Delta consistency matrix

| Clause | 1d disposition | Coherent with |
|---|---|---|
| **N-07** post-acquisition absence | **superseded** by P-01 | N-05 barrier |
| **N-08** `prepared`+intent | **superseded** by P-03 rows | P-01 |
| **N-11** inference table | **superseded** by P-03 | P-04 |
| **N-15** absent vs failure | **superseded** by P-13 | P-11, P-14 |
| N-12 / N-13 / N-14 | unchanged; the authority boundary is clarified by P-05 | P-05 |
| **N-16** canonical form stored | unchanged; `removal_basis` and `removal_attempted` added to the receipt | P-07 |
| **N-17** identity re-check | **narrowed**: detection only, per P-08 | P-06 |
| N-05 barrier | unchanged; it is what makes P-01 sound | P-03 |
| N-06 abandonment | unchanged | P-03 row 2 |
| N-10 release | unchanged | — |
| L-03 (all rows) | fully superseded by P-03 | — |
| M-01 / M-02 attribution | unchanged; P-01 adds `unknown` for a new state | P-03 |
| M-11 acquisition | refined by P-13's classification | P-11 |
| C-33 receipt fields | **extended**: `removal_attempted`, `removal_basis`, `absent_after_acquisition` | P-07 |
| C-34 / C-35 / C-36 | unchanged | P-01 |
| C-31 byte preservation | untouched | ✓ |

**Event-kind count: unchanged — five ledger kinds plus two evidence kinds.** `absent_after_acquisition`, `identity_race`, `parent_inaccessible`, `parent_unresolved`, `target_unresolved` are **receipt states and result codes**, not new event kinds.

### Residual items

1. Under `path_recheck`, a replaced directory can be affected — recorded, not mitigated (P-09, P-10).
2. `selection_resolved` and D-3's ancestor rule remain deferred.
3. A permanently corrupt journal leaves an irremovable stale entry.
4. Partial deletions and post-acquisition absences carry no authorship claim.
5. Power-loss durability `UN`; workload `Unestablished`.

### Owner decisions, now ten

The nine prior, plus: **accept `path_recheck` for Phase 1 with its recorded limitation, or require `handle_identity` before Phase 1 deletion is enabled.**

### Authorization

No implementation, migration, ledger creation, persistence adoption, or FC-1 approval. The contract is consistent as text; no test of any prefix has been executed.
