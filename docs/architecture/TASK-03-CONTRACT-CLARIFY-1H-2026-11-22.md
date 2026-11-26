# TASK-03 — clarification 1h: absence proof, ABA reuse, count deltas, volume continuity

| Field | Value |
|---|---|
| Purpose | Close four findings against Clarification 1g |
| Status | **Contract text. Nothing adopted. No implementation authorized. Deletion unavailable** |
| Relationship | **Additive.** None of rows 528–536 is edited |
| Supersedes | 1g S-06 (O1 as "not reopenable"), S-07/S-08 (`partial`), S-16 (under-claim proof), S-21 (clone limitation's resolution), S-30's `complete` scope |
| Recovery token | `TASK-03-CONTRACT-CLARIFY-1H-2026-11-22` |

Prefixes: `U-` normative, `UT-` acceptance test. Labels: **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Superseded**.

---

## 1. Failure is not absence

**Finding, accepted.** 1g made "identity not reopenable" (O1) the absent observation. That conflates a specific, established not-found result with every other failure — and the evidence says those are not distinguishable by the error code alone.

**U-01 — What the platform actually returns.** A filesystem filter developer's documentation records that opening by ID with an **ID that does not exist** fails with `STATUS_INVALID_PARAMETER` — the same status a malformed request produces (Source-derived: *Some Limitations Using Files Opened By ID*, 2012, a third-party engineering source, not vendor documentation). Windows' own `OpenFileById` page says only that failure returns `INVALID_HANDLE_VALUE` and that details come from `GetLastError`; it enumerates **no** error codes (Source-derived). No vendor source establishes a per-filesystem "not found" code for open-by-ID.

**Calculated consequence:** if the absent case and the malformed-request case share a status, then **a failed open by itself cannot establish absence**, and any clause treating it as absence is unsound.

**U-02 — `identity_not_found` requires five preconditions, all satisfied.**

> 1. The capability flag `FILE_SUPPORTS_OPEN_BY_FILE_ID` is set on the volume (Source-derived: `GetVolumeInformation` volume flags).
> 2. The descriptor is well-formed: correct `dwSize`, `ExtendedFileIdType`, the stored 128-bit ID.
> 3. The handle used satisfies the documented hint ("a handle to any file on a volume or share") and has the access the call needs.
> 4. A **control open** on the same volume, with the same descriptor shape and the same call parameters, succeeds on a known-live identity — proving the call form itself works here.
> 5. The target open then fails with a status that the control demonstrably cannot produce.

Only when all five hold is the result `identity_not_found`. Anything else is a weaker class.

**U-03 — Error taxonomy.**

| Class | Meaning | May support absence? |
|---|---|---|
| `identity_not_found` | U-02 satisfied | **yes** |
| `identity_reopened` | the call succeeded; identity present | no (present-side) |
| `access_denied` | required access unavailable | no |
| `volume_unavailable` | volume could not be reacquired | no |
| `volume_mismatch` | corroborating volume identity differs | no |
| `identity_unsupported` | capability flag absent, or the filesystem/API cannot perform the operation | no |
| `open_failure_unclassified` | any other failure, or the control open itself failed | no |

**U-04 — Unverifiable is its own outcome.** `access_denied`, `volume_unavailable`, `volume_mismatch`, `identity_unsupported` and `open_failure_unclassified` all yield **`unverifiable`**, with the failure class recorded, and are `indeterminate` for every state. They never produce absence, never produce a present-side outcome, and never advance the machine.

**U-05 — The control open is not optional.** Because the not-found status may be identical to the malformed-request status (U-01), a failed open without a successful control is `open_failure_unclassified` by definition. This is the difference between "the ID is gone" and "my call is wrong".

**U-06 — Capability is checked once, at acquisition.** `FILE_SUPPORTS_OPEN_BY_FILE_ID` is read when the volume is first opened and recorded. Absent → `identity_unsupported` at acquisition, so no journal entry is ever created for a volume where verification could never be performed. This supersedes 1g's `FileId == 0` heuristic, which was a proxy for the same fact.

---

## 2. The ABA reuse cycle

**Finding, accepted — and it breaks 1g's proof.** 1g's S-16 argued that reuse makes an ID *match*, hence a successful open, hence a present-side outcome, hence under-claiming. That covers A→B. It does not cover **A removed → ID reused by B → B removed by another actor → open fails**.

**U-07 — The scenario, stated plainly.** Our stored identity is gone; a replacement held the same ID for a while; a third party removed the replacement; verification now finds nothing. From that final failure the contract cannot tell:

- whether **our** removal removed A, or
- whether A was removed by someone else, B took the ID, and B was then removed by someone else.

**U-08 — Absence survives; attribution does not.** In every ABA ordering, the final observation `identity_not_found` establishes one thing: *no object currently answers to that identity*. What is not established is **who removed A**. So:

| Claim | Under ABA |
|---|---|
| "the acquired identity is not currently reopenable" | **established** (given U-02) |
| "this operation removed it" | **not established** |

**U-09 — The split.** `state: complete` and the attribution claim are separated:

| Field | Meaning |
|---|---|
| `state: complete` | absence established by U-02 |
| `own_removal_observed` | **fact**: this operation's removal call returned success in-process |
| `reuse_excluded` | **evidence claim**: no reuse of the identity is possible between removal and verification |
| `attribution` | `operation` only if `reuse_excluded`; otherwise `unknown` |

**U-10 — `reuse_excluded` requires evidence Phase 1 does not have.** Establishing it needs a per-volume, per-reference event history covering the window between acquisition and verification. The candidate is the USN change journal: a create or delete record for that file reference inside the window would prove reuse. Requirements — journal enabled, not truncated past our starting USN, and the file reference's records retained — are **Unestablished** (E-24). Without them, `reuse_excluded = false` and `attribution: unknown`.

**U-11 — Consequence, accepted deliberately.** In Phase 1 **no receipt attributes removal to the operation**. Every `complete` reads: the identity is gone, our removal call returned success, authorship unknown. This is weaker than the receipts in 1e–1g, and it is what the ABA finding requires. It is not a new owner decision — it follows from the rule that a receipt may not claim more than the evidence establishes.

**U-12 — Can ABA produce a false `complete`?** No. `complete` requires `identity_not_found`, which requires that no object answers to the identity. If A were still present, its own ID would resolve and the call would succeed. ABA can therefore only weaken **attribution**, never the absence claim. Stated as a proof sketch: *absence is a statement about the identity, and ABA changes who removed it, not whether it is there.*

**U-13 — Creation-time corroboration does not help here.** 1g's R-09/S-18 used creation time to suspect reuse. In ABA the replacement is **gone** at verification, so there is no object whose creation time could be read. Creation time can only downgrade a **present**-side observation; it cannot strengthen an absent-side one. 1g implied otherwise and is corrected.

**U-14 — Supported identity lifetime and recovery boundary.**

> An identity supports a **present-side** claim while the volume handle is continuously held and a fresh open matches ID, volume and creation-time corroboration. It supports an **absent-side** claim only as `identity_not_found` (U-02) with `reuse_excluded` unproven — i.e. absence without authorship. Beyond the process lifetime, no claim is supported at all: reconciliation outcomes are `absent_after_acquisition`, `absent_unverified`, `unverifiable` or `indeterminate`, and never `complete`.

---

## 3. Count deltas

**Finding, accepted.** `entries_at_acquisition` supports a count difference and nothing else. 1g's `partial` invited a reading it cannot support.

**U-15 — `partial` is renamed `count_delta_unattributed`.** Null hypothesis: the state name is the claim.

| Field | Meaning |
|---|---|
| `state: count_delta_unattributed` | a **complete** enumeration showed a child count different from `entries_at_acquisition` |
| `count_delta` | signed difference |
| `attribution` | **always `unknown`** |

**U-16 — What it establishes and what it does not.**

| Established | Not established |
|---|---|
| a count difference exists | which children changed |
| the difference was observed at a stated time | that the directory is otherwise intact |
| — | that this operation caused the change |
| — | that any child was removed by Howl rather than renamed, added, or replaced |

**U-17 — Equal counts are never "unchanged".** One child out and one child in leaves the count identical. `count_equal` is a count observation; it may not appear in a receipt as "intact", "unchanged", "no changes" or any synonym. `count_delta_unattributed` requires a complete enumeration (1g S-10 retained): with enumeration incomplete, the outcome is `not_completed` or `unverifiable`, never a delta.

**U-18 — No manifest, still.** This clarification does not introduce a child manifest or a persisted tree snapshot. 1g's assessment stands: P2–P5 remain unestablishable, and the count remains diagnostic only.

---

## 4. Volume continuity against D-10

**Finding, accepted.** 1g admitted clone-substitution and left D-10's fate implicit. Resolved here — not by weakening D-10, and not by redefining identity as "whichever volume the locator reaches".

**U-19 — Continuity is established by a held handle, not by metadata.**

> The volume is opened at acquisition and the handle is **held continuously** through removal and verification. The destructive act and the verification are performed against that held object. Metadata comparison (GUID, serial, filesystem) is a **corroborator at open time**, not the continuity mechanism.

**Calculated:** a volume handle refers to the volume object it was opened on. While it is held, the locator cannot be re-pointed at a different device for the purposes of these calls. That is why D-10 is satisfiable in-process and not across a restart.

**U-20 — D-10 is satisfied by construction for the destructive act.** Removal is performed through the handle obtained from the acquired identity, in the process that acquired it. Nothing re-resolves the locator between acquisition and destruction. So the deletion acts on the exact object acquired — which is precisely D-10's requirement.

**U-21 — Cross-restart: no completion claim at all.** After a restart the held handle is gone, continuity cannot be established, and therefore:

> reconciliation never writes `complete`. Its outcomes are confined to `absent_after_acquisition`, `absent_unverified`, `unverifiable`, `indeterminate`, or an abandoned entry. Only the operation itself, in-process, can write `removed` and thus `complete`.

This was already implicit (`removed` is operation-written); 1h makes it the continuity rule and bounds the clone exposure.

**U-22 — Clone cases, separated.**

| Case | Detectable? | Effect |
|---|---|---|
| Clone attached **simultaneously** with the original | **Yes in principle**: two volumes answering one GUID → `identity_ambiguous` (U-04) | no deletion, no claim |
| Original detached, clone attached **instead**, across a restart | **No** | cannot yield `complete` (U-21); at most `absent_after_acquisition`/`unverifiable` |
| Clone attached **mid-process** while the original handle is held | Not by metadata; the handle still refers to the original | removal is not redirected (U-20); corroboration may fail → `unverifiable` |
| Image-level restore reproducing old MFT records | **No** | same as detachment: no `complete` across restart |

**U-23 — Unsupported scenarios, declared.** Identity-bound deletion is **unsupported**, and therefore unavailable, where:

- the volume lacks `FILE_SUPPORTS_OPEN_BY_FILE_ID`;
- the filesystem is FAT/exFAT/SMB or otherwise lacks stable 128-bit identity (1f R-07 retained);
- the volume is expected to be cloned, imaged or replaced by another party while an operation is pending;
- continuity cannot be maintained (for example, the volume is removable and may be detached mid-operation).

Unavailability is reported as `identity_unsupported` / `deletion_unavailable`. It is not a degradation to path-based deletion.

**U-24 — What is not claimed.** No clause claims physical device continuity, clone detection, or that `complete` proves the acquired *device* was the one reached. `complete` is bounded by U-02 (absence), U-09 (attribution), and U-21 (in-process continuity).

---

## 5. Corrected observation matrix

**U-25 — Observation classes.**

| Code | Class | Source |
|---|---|---|
| **K1** | `identity_not_found` — U-02 satisfied | absent-side |
| **K2** | `identity_reopened` — present | present-side |
| **K3** | `unverifiable` — `access_denied`, `open_failure_unclassified`, failed control open | no claim |
| **K4** | `volume_unavailable` / `volume_mismatch` / `identity_ambiguous` | no claim |
| **K5** | `identity_unsupported` | no claim |

**U-26 — The matrix.** 4 states × K1 + (K2 with four sub-observations) + K3/K4/K5.

| State | K1 absent | K2 present: enum. incomplete | K2 present: count equal | K2 present: count differs | K3/K4/K5 unverifiable |
|---|---|---|---|---|---|
| `prepared` | `absent_after_acquisition` / `unknown` | `aborted` / — | `aborted` / — | `aborted` / — | `unverifiable` / `unknown` |
| `pending` | `absent_after_acquisition` / `unknown` | `aborted` / — | `aborted` / — | `aborted` / — | `unverifiable` / `unknown` |
| `removing` | `absent_unverified` / `unknown` | `not_completed` / `unknown` | `not_completed` / `unknown` | `count_delta_unattributed` / `unknown` | `unverifiable` / `unknown` |
| `removed` | `complete` / **`unknown`** (U-11) | `indeterminate` / `state_inconsistent` | same | same | `unverifiable` / `unknown` |

**U-27 — What changed from 1g.** `complete` no longer carries `attribution: operation`; `partial` became `count_delta_unattributed`; empty or unreadable volumes no longer look like absence (they are K3–K5); and `absent_after_acquisition` now requires K1, so a permission error is never reported as "the directory was already gone".

**U-28 — Corrections to the crash table (1f §3.1).** Its columns A–D assumed a binary absent/present observation. Replacing them:

| Condition | Now |
|---|---|
| A (absent) | requires **K1**; anything else is K3–K5 |
| B/C (present, counts) | K2 sub-observations; `aborted` at not_started, `not_completed`/`count_delta_unattributed` at possible |
| D (absent + replacement at path) | K1 with a path note; replacement irrelevant to the outcome |
| — | **new**: K3–K5 → `unverifiable`, all states |

Windows 9–10 keep their meaning, but window 9's `complete` now carries `attribution: unknown` and requires in-process continuity (U-21).

---

## 6. Preserved conclusions

**U-29 — Carried unchanged.** Count-only baseline (no manifest). The five-property separation (count, set, child identity, content, authorship) with only count establishable. FAT/exFAT refusal. Immutable terminal receipts and transient post-terminal anomalies (1g S-24–S-27). `aborted` asserting only "this operation performed no removal". `absence_unverified` as a result code. `absent_prior` gated and unreachable. D-10 mandatory; G-1 an engineering gate.

**U-30 — No new owner decision.** G-1 stands: deletion remains unavailable until the identity-bound mechanism is established and its evidence passes. U-11's loss of attribution is a consequence of the evidence rules, not a choice offered to the owner.

---

## 7. Acceptance tests

- **UT-01** Failed open with no successful control → `open_failure_unclassified` → `unverifiable`, never `absent_*`.
- **UT-02** Malformed descriptor and absent identity produce the same status; only the control-open procedure separates them.
- **UT-03** `ERROR_ACCESS_DENIED` → `unverifiable`, never absence.
- **UT-04** Volume lacks `FILE_SUPPORTS_OPEN_BY_FILE_ID` → `identity_unsupported` at acquisition; no journal entry created.
- **UT-05** Fabricated `identity_not_found` without U-02's five preconditions is rejected as invalid.
- **UT-06** No receipt outside K1 carries an absent-side state.
- **UT-07** ABA simulation (A removed, ID reused by B, B removed, then verify) → `identity_not_found` true, `reuse_excluded` false, `attribution: unknown`; no receipt states the operation removed A.
- **UT-08** `reuse_excluded = true` is unreachable while E-24 is unmet; any code path setting it fails the test.
- **UT-09** Creation-time corroboration is never consulted on an absent-side observation.
- **UT-10** `own_removal_observed` is a fact recorded independently of attribution.
- **UT-11** `complete` with `attribution: operation` is refused; the vocabulary no longer permits it.
- **UT-12** `count_delta_unattributed` carries `count_delta`, `attribution: unknown`, and no statement about which children changed.
- **UT-13** Equal counts never surface as "intact"/"unchanged" in any user-facing string.
- **UT-14** `count_delta_unattributed` is unavailable when enumeration is incomplete.
- **UT-15** Removal and verification both use the handle opened at acquisition; no locator re-resolution occurs between them (call-graph and handle-lifetime assertion).
- **UT-16** Reconciliation after a simulated restart never writes `complete`.
- **UT-17** Two volumes answering one GUID → `identity_ambiguous`; no deletion.
- **UT-18** Clone-substitution fixture across restart → at most `absent_after_acquisition`/`unverifiable`; no `complete`.
- **UT-19** Removable-volume detach mid-operation → `unverifiable`; entry preserved.
- **UT-20** Every cell of U-26 is reachable and produces the stated outcome; K3–K5 yield `unverifiable` in all four states.

---

## 8. Remaining unestablished claims

Carried: E-1…E-23.

| # | Claim | Evidence required |
|---|---|---|
| **E-24** | USN journal can prove `reuse_excluded` (journal ID continuity, retention window, per-reference records) | Destructive spike; if unmet, U-11 stands permanently |
| **E-25** | The exact status a deleted identity produces on NTFS, ReFS and UDF, and that a control open reliably distinguishes it from a malformed request | Per-filesystem matrix |
| **E-26** | Whether `identity_ambiguous` is actually produced when two volumes share a GUID, or whether the mount manager silently renames one | Spike with a cloned volume |
| **E-27** | That a held volume handle prevents locator re-pointing for the lifetime of the calls | Adversarial spike (detach/attach while held) |

**U-31 — The one proof that survives.** Absence is `identity_not_found` and nothing else; reuse can weaken attribution but cannot manufacture a false absence (U-12). Everything else in this document is a **claim deferred to evidence**, and labelled as such.

---

## 9. Verification and authorization

Rows 528–536 remain byte-identical; 1g (`1db25754a88b…`) unchanged. 1h registered as row 537; corpus counts updated; no historical entry altered. **Zero Go files modified; no ledger created; no migration; no deletion enabled.** No test executed: the failure taxonomy is source-derived (including one third-party engineering source, labelled), and every behavioural claim is `Unestablished` with its evidence named.

TASK-03 remains in contract preparation. Not accepted, not closed.
