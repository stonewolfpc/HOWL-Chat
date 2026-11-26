# TASK-03 — clarification 1g: observation space, identity lifetime, volume continuity

| Field | Value |
|---|---|
| Purpose | Close five consistency findings against Clarification 1f |
| Status | **Contract text. Nothing adopted. No implementation authorized. Deletion unavailable** |
| Relationship | **Additive.** None of rows 528–535 is edited |
| Supersedes | 1f's four-condition framing (A–D) and its "intact" wording; R-08's safety claim; R-21/R-22's scope; the `removing`+present outcome split; 1f's implicit "special volume handle" reading |
| Recovery token | `TASK-03-CONTRACT-CLARIFY-1G-2026-11-22` |

Prefixes: `S-` normative, `ST-` acceptance test. Labels: **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Superseded**.

---

## 1. The four conditions were not exhaustive, and one of them was a lie

**Finding, accepted.** 1f's matrix used four target conditions: absent, `entries == baseline` ("intact"), `entries < baseline` ("partial"), and absent-with-replacement. That is not exhaustive, and the word "intact" asserted something the acquisition record cannot support.

**S-01 — Five properties, and only one of them is establishable.**

| # | Property | Establishable from `entries_at_acquisition` + a live enumeration? |
|---|---|---|
| **P1** Child **count** | current count vs baseline | **Yes** |
| **P2** Entry **set** | same names present | **No** — no baseline set is stored |
| **P3** Child **identity** | same file IDs under those names | **No** — no baseline IDs are stored |
| **P4** **Content** | same bytes/hashes | **No** — and never in scope |
| **P5** **Attribution** | this operation caused the change | **No** — see S-03 |

**S-02 — "Intact" is retired.** Count equality does **not** mean unchanged. A directory can lose one child and gain another and report the same count. Every occurrence of "intact" in 1e/1f is superseded by `count_equal`, which is a count observation and nothing more.

**S-03 — Count change is not authorship.** `count_lower` does **not** mean Howl removed anything, and `count_higher` does not mean it failed to. With P2–P5 unestablishable, no count observation carries authorship. The only clause that attributes removal is the verification clause (S-12), and it attributes **absence**, not the deletion of individual children.

**S-04 — Enumeration can fail, and that is a state.** A directory may be present while its contents cannot be enumerated (access denied, transient lock, I/O error, adversarial churn). This is `enumeration_incomplete`, and it is **not** the same as an empty or unchanged directory: no count comparison is available at all.

**S-05 — What the contract actually needs to establish.** Only two things:

1. whether the **acquired identity** is reopenable (the deletion guarantee's subject), and
2. whether the acquired identity was verified absent by a **fresh** open after our removal.

Everything child-level is **diagnostic**: it may inform a receipt's `tree_observation` field, and it may never carry attribution. That is the honest consequence of storing a count instead of a manifest (S-07).

---

## 2. Exhaustive observation matrix

**S-06 — The observation space.** After reacquisition by identity:

| Code | Observation |
|---|---|
| **O1** | identity not reopenable |
| **O2** | identity reopenable, enumeration incomplete |
| **O3** | identity reopenable, enumeration complete, count equal to baseline |
| **O4** | identity reopenable, enumeration complete, count lower than baseline |
| **O5** | identity reopenable, enumeration complete, count higher than baseline |

Path occupancy is **orthogonal** and never part of the observation: a replacement identity at the original path is a different object (O1-with-replacement is still O1).

**S-07 — Outcomes, by removal evidence.** Three present-side outcomes replace 1f's `aborted`/`partial` split:

| Evidence | Present-side outcome | What it asserts |
|---|---|---|
| `not_started` (prepared/pending) | `aborted` | this operation performed no removal. **Only** this. Not that the tree is unchanged |
| `possible` (removing) | `not_completed` | removal may have been invoked and did not complete; no tree claim |
| `possible` + a complete enumeration showing count ≠ baseline | `partial` | the tree differs from baseline **and** this operation may have removed something. **Authorship unknown** |
| `verified` (removed) | `indeterminate` | see S-19 |

**S-08 — The matrix. 4 states × 5 observations = 20 cells.**

| State | O1 absent | O2 present, enum. incomplete | O3 present, count equal | O4 present, count lower | O5 present, count higher |
|---|---|---|---|---|---|
| `prepared` (not_started) | `absent_after_acquisition` / `unknown` | `aborted` / — | `aborted` / — | `aborted` / — | `aborted` / — |
| `pending` (not_started) | `absent_after_acquisition` / `unknown` | `aborted` / — | `aborted` / — | `aborted` / — | `aborted` / — |
| `removing` (possible) | `absent_unverified` / `unknown` | `not_completed` / `unknown` | `not_completed` / `unknown` | **`partial` / `unknown`** | **`partial` / `unknown`** |
| `removed` (verified) | `complete` / `operation` | `indeterminate` / `state_inconsistent` | `indeterminate` / `state_inconsistent` | `indeterminate` / `state_inconsistent` | `indeterminate` / `state_inconsistent` |

**S-09 — Three cells worth reading carefully.**

- **O4 at `prepared`/`pending`**: a smaller tree while no removal could have run is an **external** change. The outcome is `aborted` — the operation removed nothing — with `tree_observation: count_lower`. The receipt must not imply the operation deleted those children.
- **O4 at `removing`**: `partial` is permitted, attribution `unknown`. The operation may have done it; the evidence cannot say.
- **O5 at `removing`**: a *larger* tree is also `partial`. External writes during deletion are the likely cause, but "likely" is not a claim the receipt may make.

**S-10 — `partial` requires a complete enumeration.** With `enumeration_incomplete`, `partial` is unavailable: no difference was observed, so none can be reported.

**S-11 — Impossible cells.** `removed` with any present-side observation is structurally impossible (verification precedes the write); it is `state_inconsistent`, reported, never repaired. `partial` at `not_started` evidence is impossible.

---

## 3. Should the acquisition record store a tree manifest?

**S-12 — Assessed and declined for Phase 1.** A manifest (every child name + file ID) would make P2/P3 establishable. Costs:

| Dimension | Cost |
|---|---|
| Storage | one record per child; a conversation directory can hold hundreds of files; the journal is rewritten atomically, so size affects commit latency |
| Performance | full enumeration at acquisition, on the user's delete path |
| Privacy | child names are durable structural data; attachment names in particular can be content-adjacent |
| Recovery | reconciliation would need the manifest to be readable and versioned; a corrupt manifest becomes a second failure mode that can block deletion |
| Value | it would establish P2/P3 **at a moment in the past**, which still cannot attribute a later change to this operation without also proving exclusivity — which the platform does not give (E-3) |

**S-13 — Consequence, stated plainly.** With a count only, **historical tree integrity cannot be established**, and no clause may claim it. The count exists solely to distinguish "no observed difference" from "a difference was observed" for diagnostic purposes. If the owner later wants P2, it is a separate decision with the costs above — not a silent addition.

---

## 4. File-ID reuse and sequence wrap

**Finding, accepted.** 1f said a stale ID fails to match because the sequence increments (R-08), and also listed sequence wrap as an explanation for an identity reappearing (R-27). Both were true, but together they read as "safe" while admitting the exception. Resolved below.

**S-14 — Structure (Source-derived).**

| Filesystem | 128-bit `FileId` composition | Source |
|---|---|---|
| NTFS | 48-bit MFT index (low) + 16-bit sequence (next 16); high 64 bits zero | MS-FSA Appendix A, note 48 |
| ReFS | 64-bit parent-directory index (low) + 64-bit per-directory file index (high) | MS-FSA Appendix A, note 48 |
| FAT/exFAT | Not a stable identity; derived from directory-entry position; changes on defrag or rename | `BY_HANDLE_FILE_INFORMATION` Remarks |

**S-15 — Sequence wrap is real and its bound is arithmetic.** Freeing an MFT record increments its 16-bit sequence. A stored ID returns to matching a *different* object only after that specific record is freed and reused **2^16 = 65,536 times** and the sequence lands on the stored value (Source-derived structure, **Calculated** bound). That is not zero, so "stale IDs always fail" is wrong.

**S-16 — The direction of error, proved rather than asserted.** `complete` requires a **fresh `OpenFileById` to fail** (S-18). Reuse makes an ID *match*, which produces a successful open, which yields present-side outcomes (`indeterminate`, `not_completed`, `partial`) — never `complete`. Therefore:

> Reuse and sequence wrap can only cause **under-claiming**: a deletion that succeeded may be reported as `indeterminate`. They cannot produce a false `complete`.

This is the one asymmetry the contract can rely on, and it holds because absence is never inferred from a matching identity — only from a failed open.

**S-17 — ReFS.** IDs are described as monotonically increasing per directory and not reused for the volume's lifetime (Source-derived from a technical reference; vendor confirmation outstanding as E-15). Until confirmed, ReFS uses the same rule as NTFS: no assumption of permanence, and any corroboration failure downgrades to `indeterminate`.

**S-18 — Supported identity lifetime.**

> An identity is treated as continuous while: the volume locator reaches a volume whose serial and filesystem match the record; a fresh `OpenFileById` on the stored ID succeeds; and the reopened object's creation time is **not later** than the journal's acquisition timestamp.
>
> An open that succeeds but fails the creation-time corroboration is `identity_reuse_suspected` → `indeterminate`. Corroboration may only **lower** confidence, never raise it.

**S-19 — Cloned or restored volumes.** A byte-cloned volume preserves GUID, serial, filesystem and IDs, so the reopened object is indistinguishable from the original by any evidence this contract collects. There is **no** clause that detects it. Consequences:

- `complete` is scoped to "the volume reached by the stored locator, at verification time".
- A restore that re-creates the directory on the *same* volume produces a **new** identity, so it appears as O1 (absent) or as O3–O5 under a different ID — never as the same ID.
- An image-level restore that reproduces old MFT records reproduces the old ID; that is `identity_reuse_suspected` only if the creation time disagrees, otherwise indistinguishable — and is stated as a limitation rather than solved.

---

## 5. Volume identity, and the `OpenFileById` requirement

**S-20 — Five distinct identities, only four of which the contract establishes.**

| Layer | Meaning | Established here? |
|---|---|---|
| **Locator** | GUID path or drive letter — how to reach a volume | yes (stored) |
| **Volume identity** | GUID + serial + filesystem — does the reached volume claim to be the recorded one | yes (serial + fs compared) |
| **Filesystem identity** | the formatted filesystem instance: serial, format era | partially — a reformat changes it, which surfaces as a mismatch |
| **Directory identity** | file ID within that filesystem | yes |
| **Physical device identity** | the hardware (disk serial, device path) | **no** — not stored, not required, not obtainable without elevation |

**S-21 — The honest gap.** Because physical device identity is not established, a different physical volume presenting matching identifiers is **undetectable** by this contract. That is the cloned-volume limitation of S-19, restated at the volume layer. The contract does not claim device continuity and must not be described as if it did.

**S-22 — Correction to 1f's reading of `OpenFileById`.** The documented requirement is "a handle to any file on a volume or share on which the file to be opened is stored" (Source-derived). 1f implied a special volume handle; that is not required. What the contract requires is **some handle on the correct volume**, and in practice, with no path available after a restart, the volume GUID root is how one is obtained. Whether a volume-root handle satisfies the hint is E-12, not a documentation certainty.

**S-23 — Continuity failure classes, all fail-closed.** `volume_unavailable`, `volume_identity_changed` (serial or filesystem differs), `volume_unopenable`, `identity_ambiguous` (two volumes answering the same GUID), `identity_absent`, `identity_mismatch`, `identity_reuse_suspected`. Each yields `indeterminate`, preserves the entry, deletes nothing, and never consults a path.

---

## 6. Post-terminal inconsistency

**S-24 — The rule.**

> When reconciliation finds a terminal receipt and a fresh check contradicts it — the acquired identity is present, or the volume no longer matches — the terminal receipt is returned **unchanged**. Deleting nothing, re-verifying nothing, rewriting nothing.

**S-25 — What is reported, separately and without touching history.**

| Condition | Anomaly code |
|---|---|
| terminal `complete`, acquired identity reopens | `post_terminal_identity_present` |
| terminal exists, volume locator no longer matches its record | `post_terminal_volume_mismatch` |
| terminal exists, identity cannot be checked at all | `post_terminal_unverifiable` |

**S-26 — The anomaly is transient, not persisted.** It is returned to the caller and emitted as a diagnostic line. It is **not** appended to the deletion-evidence file, and it does not create a third evidence kind.

*Justification:* the terminal receipt is the record of what happened; an anomaly is an observation about the present. Persisting it would create a mutable second history adjacent to immutable receipts, and any consumer could come to treat the later file as an authority — the same failure mode the receipt's immutability exists to prevent. Adding a persistent anomaly store would also introduce a new evidence kind without a defect in the existing two.

*Cost, stated:* the anomaly has no history. Two reconciliations a week apart report it twice; nothing records that it was seen. If the owner wants durable anomaly history, that is a separate decision with its own justification — not a silent addition.

**S-27 — The anomaly changes no state.** It does not block deletion of other identities, does not mark the conversation unusable, does not alter attribution, and does not cause a retry. A later deletion of that directory is a **new operation** with a new identity and a new receipt.

---

## 7. Guarantees preserved

**S-28 — D-10 unchanged.** Deletion stays bound to the acquired physical directory identity; no path fallback; no weakening.

**S-29 — G-1 unchanged.** Deletion remains unavailable until the identity-bound mechanism is demonstrated and its acceptance evidence passes. No additional owner decision is created by this document.

**S-30 — `complete`, restated precisely.**

> `complete` means the acquired directory identity was no longer reopenable through the supported mechanism, by a fresh open, after this operation's removal returned — on the volume reached by the stored locator.

It does **not** mean: physical erasure; that all handles are closed; that no data stream remains accessible (POSIX semantics keeps data reachable through pre-existing handles — Source-derived); that no copy exists on another volume; or that the original path is empty.

---

## 8. Acceptance tests

- **ST-01** Count equal with one child swapped for another → `count_equal`, never "intact", never `aborted`-with-no-change language.
- **ST-02** `not_started` + `count_lower` → `aborted` with `tree_observation: count_lower`; receipt text makes no authorship claim.
- **ST-03** `removing` + `count_lower` → `partial`/`unknown`; `removing` + `count_higher` → `partial`/`unknown`.
- **ST-04** `removing` + enumeration incomplete → `not_completed`, **not** `partial`.
- **ST-05** All 20 cells of S-08 are reachable and produce the stated outcome; no cell carries `operation` attribution except `removed`+O1.
- **ST-06** Any `partial` or `aborted` receipt asserting that Howl deleted specific children fails a text assertion.
- **ST-07** No manifest is written; `entries_at_acquisition` is the only tree baseline in durable state.
- **ST-08** NTFS stale-ID simulation (record reuse with renewed sequence) → fresh open fails → `complete`.
- **ST-09** Sequence-wrap simulation (reused record landing on the stored sequence) → fresh open succeeds → `indeterminate`, **never** `complete`.
- **ST-10** Reopened identity with creation time later than acquisition → `identity_reuse_suspected`, downgraded to `indeterminate`, even when serial and ID match.
- **ST-11** ReFS and NTFS are both exercised under ST-08/ST-09; FAT is refused as `identity_unsupported`.
- **ST-12** Cloned-volume fixture: no code path can distinguish it; the limitation is documented and `complete`'s claim is asserted to be locator-scoped.
- **ST-13** `OpenFileById` is called with a handle on the volume obtained from the GUID root, not with a fabricated special handle type; the access mode used is recorded.
- **ST-14** Volume serial or filesystem mismatch → `volume_identity_changed`; nothing deleted.
- **ST-15** Terminal `complete` + identity reopens → receipt returned byte-identical, anomaly `post_terminal_identity_present` reported, no deletion, no rewrite.
- **ST-16** Terminal + volume mismatch → `post_terminal_volume_mismatch`; receipt untouched.
- **ST-17** Anomalies appear in no persistent file and create no third evidence kind.
- **ST-18** A new deletion of the same path after a `complete` receipt is a new operation with a new `operation_id`; the old receipt is not reused or amended.

---

## 9. Remaining unestablished claims

Carried: E-1…E-18. New:

| # | Claim | Evidence required |
|---|---|---|
| **E-19** | Whether USN journal ID or another filesystem "generation" marker could corroborate volume continuity | Spike; no documented guarantee assumed |
| **E-20** | ReFS non-reuse over a volume's lifetime (vendor confirmation of S-17) | Vendor statement or long-running fixture |
| **E-21** | Exact NTFS sequence-wrap reachability in practice (record reuse rate) | Fixture with controlled churn |
| **E-22** | Whether creation-time corroboration is meaningful on each filesystem (timestomping aside) | Spike, per filesystem |
| **E-23** | Whether an image-level restore reproducing old MFT records is distinguishable at all | Fixture; likely **no**, and then stated as a limitation |

**S-31 — No false `complete` is claimed as impossible for every mechanism.** It is **proved excluded** for reuse and wrap (S-16), and **openly admitted** for clones (S-19/S-21). Those are different statements and the contract keeps them apart.

---

## 10. Verification and authorization

Rows 528–535 remain byte-identical; 1f (`d386317a316f…`) is unchanged. 1g is registered as row 536; corpus counts updated; no historical entry altered. **Zero Go files modified; no ledger created; no migration; no deletion enabled.** No test executed — this document is contract text plus source-derived API structure, with S-16's direction-of-error as the only *proved* safety property.

TASK-03 remains in contract preparation. Not accepted, not closed.
