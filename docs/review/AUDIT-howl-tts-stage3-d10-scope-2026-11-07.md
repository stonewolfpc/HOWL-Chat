# AUDIT — D10 scope split in the Stage 3 corrective design

**Date:** 2026-11-07
**Mode:** DOCUMENTATION REVIEW ONLY — read-only; nothing edited, repaired, or reconciled
**Subject:** `DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md` (corrective proposal) and
`CORRECTIVE-howl-tts-stage3-corrective-design-2026-11-07.md` (correction record)
**Governing contract:** `DESIGN-howl-tts-helper-architecture-2026-11-06.md`
**Focus:** the split of D10 into D10(a) and D10(b), the exclusion of the A2-D9 supervisor path, and
whether the new direct-`Cancel` branch weakens the supervisor-path behavior or contradicts A2-D1
**Reviewer position:** author of the round-1/2/3 audits; see §11, limitation 1

## 1. What this verification is

A delta check of one change to a document set already audited three times. It answers four questions
and nothing more:

1. Does the text actually implement the D10 split the record describes?
2. Are the two branches independently answerable, and can one be dispositioned without the other?
3. Does D10(b) — the direct `Cancel` unknown-snapshot branch — weaken the A2-D9 supervisor-path
   behavior or contradict A2-D1's direct-operation rule?
4. Did the pass reintroduce any earlier defect, and are the round-3 findings M5 and m1–m5 genuinely
   addressed in the text rather than only in the record?

It is not a re-audit of C1–C8. Where earlier findings are carried forward, they are spot-verified at
their current locations, not re-derived from scratch.

**Severity rule used:** *blocking* means the text would freeze a wrong contract interpretation into a
stage's acceptance tests, or would let a stage be accepted while a contract-settled behavior is
unimplemented. Everything else is minor.

## 2. Verified inputs

Recomputed from disk at review time:

| Artifact | SHA-256 | Size |
|---|---|---|
| Governing architecture | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f` | 2,586 lines / 182,904 bytes |
| Corrective proposal (this revision) | `911fbb83a7fe238d9a59c4a29157a5d9bdd8cbb6158481c934dbaacf46150a84` | 711 lines / 68,383 bytes |
| Correction record (this revision) | `573f303396ded8270e19b816f84fd9f7139cf1a1194ea9e014b3e6df5521de4f` | 165 lines / 14,339 bytes |
| Stage 3 readiness investigation | `a63f9f32f4d1bd047746a7bbbf5f509588e13ebd4272cffeb23f61fcfa891c6c` | 586 lines / 51,248 bytes |
| Audit round 1 | `71a17369e2389fae26e86783fc4df802ab91bad3daef512c0e7f743067e94ea2` | 249 lines / 25,962 bytes |
| Audit round 2 | `a7384da9f3f38661fc28b3de48badb82245cfd41efb8e49964f17bcfaaf5a7e2` | 462 lines / 37,477 bytes |
| Audit round 3 | `6681e11b7eb770d59f3c79294206f05c8d3a3d9d4cfbbd038ef103fc8c48fd15` | 349 lines / 21,676 bytes |

The record's §1 table (lines 11–20) matches disk on every row, including the audit and investigation
hashes it did not author. The contract is byte-identical to the hash the previous three rounds
reviewed, so no contract drift is possible in this pass.

The record also lists an intermediate proposal revision, `7b73635d…` at 711 lines / 67,037 bytes,
taken between the initial third pass and the D10 clarification (line 15). That file no longer exists
because the proposal was corrected in place. Its size sits correctly between the previous revision
(710 / 65,213) and this one (711 / 68,383) with the same 711-line count, so the claimed chain is
plausible, but it is not independently verifiable. See §11, limitation 2.

Modification times confirm the pass was confined to two files: proposal 11:24:05, record 11:24:42.
The round-3 audit (10:29:08) and `README.md` (10:29:45) retain their earlier mtimes and hashes; no
other file in the tree changed after the previous review.

## 3. Disposition

**PASS, with two minors.** The D10 split is implemented as described, the supervisor path is not
weakened, and A2-D1 is not contradicted. Both branches gate C1 acceptance in every place that
acceptance is defined or summarized. The round-3 major M5 and minors m1–m5 are addressed in the
proposal text itself, not merely asserted in the record.

**C1 remains blocked for acceptance by D10(a) and D10(b)** — both, independently, as the record and
proposal require. This audit does not disposition either branch and does not authorize any stage.

No prior correction was narrowed, weakened, or lost. Two minor documentation defects are recorded in
§7: an owner-decision count in §1 that no longer sums, and one owner-decision candidate phrase that
does not name a result.

## 4. The requested checks

### 4.1 Does the split exist in the text, with two separately answerable branches?

Yes. The proposal states it in §1 (36), in the §4.4 policy bullets (121–123), in C1's prerequisites
(137), in C1's `client.go` bullet (142), in three separate tests (156, 159, 162), in the independent
review instruction (175), in acceptance (177), in C7's prerequisites (433), in the requirements
matrix (598), in §9.1 (630, 633), in §10's introduction (670), and in the D10 ledger row (684).

The branches are defined as:

- **D10(a):** unknown-snapshot **direct `unload`** result.
- **D10(b):** unknown-snapshot **direct `Cancel`** result outside the supervisor path.
- **Excluded:** the A2-D9 supervisor cancellation path, "because its behavior is already specified."

Independence is explicit rather than implied. The ledger row reads: "The owner may select the same or
different result for D10(a) and D10(b); neither is frozen. C1 implements known-present and
known-absent only" (684). §1 adds that the direct-`Cancel` branch "is separately answerable and is not
assumed to receive the same disposition" (36). That is the property required for partial closure not
to be mistaken for acceptance: the two answers cannot be derived from one another, and neither is
pre-filled by an unstated assumption.

### 4.2 The supervisor path: A2-D9 is unmodified in every location

The contract settles this path, and the design reproduces it rather than restating it loosely. A2-D9
(2510) reads, in part: "If the helper advertised `howl.tts.cancel`, it sends the cancel frame using
the existing path. If it did not, the adapter skips the cancel frame; absence of `cancel` is not an
error in the supervisor's cancellation path and no `unsupported_capability` is synthesized there."

The proposal's supervisor bullet (122) preserves all four of A2-D9's elements verbatim in meaning:
record `cancel_requested`, skip the frame when the capability is absent **or unknown**, synthesize no
`unsupported_capability`, and continue the existing grace/termination rules. C1's supervisor-path
test for the unknown case (162) adds only that no helper may be started solely to discover the
capability, which follows from §17.7.2's lazy-start row rather than from A2-D9. The two-sided
regression test (163) keeps the direct and supervisor behaviors asserted against the same helper, so
a future change cannot silently merge them.

The internal path-scoped operation (122) exists for one reason: to keep the exemption from leaking
into the exported method. Nothing in the text authorizes the path to consume a direct-operation
result, and the review instruction (175) makes that an explicit review question: "confirm the path
uses an internal path-scoped operation and direct `Cancel` still returns the direct-operation result."

**No weakening found.** The unknown-snapshot supervisor branch is described in strictly the same
terms as before the split; the only change is that the direct branch beside it now has its own
decision ID instead of borrowing the unload branch's unresolved slot.

### 4.3 Direct `Cancel` and A2-D1: no contradiction

A2-D1 (2502) scopes the direct-operation rule to capability **absence**: "The adapter returns
`unsupported_capability` for an absent optional service capability and MUST NOT write an operation
frame; §4.6 line 397 already scopes that code to an optional capability."

The adopted text therefore supplies two rules and no third:

1. **Known absent, direct operation** → `unsupported_capability`, no frame.
2. **Supervisor path, absent or unknown** → no frame, no `unsupported_capability` (A2-D9).

It says nothing about the **unknown** snapshot for a direct operation. A2-D12 (2513) does not close
that gap either: it restricts `catalog`, `voice-list`, `health`, and status requests from triggering
*`unavailable` recovery*, which is a different transition from starting an `absent` helper.

The proposal does not fill this silence by assumption. It freezes the two behaviors the contract does
supply and defers only the result: line 121 keeps known-absent `unsupported_capability` for both
direct operations, and line 123 keeps "no frame, no spawn" for the unknown case while recording
"D10(a) and D10(b) respectively" as the unresolved results. C1's `client.go` bullet (142) restates
the same split per operation.

This is the conservative reading. Sending a frame to a helper whose capability snapshot is unknown
risks writing a service operation the helper may not have advertised, which is what A2-D1's "MUST NOT
write an operation frame" forbids on the absent path; suppressing the frame is safe under either
possible disposition. Deferring the *result code* is legitimate because no adopted text selects one.

**No contradiction found.** D10(b) does not change the direct known-absent rule, does not extend the
A2-D9 exemption to the exported method, and does not assume that "unknown" behaves like "absent" —
treating an unverified capability as absent is listed as one unselected candidate direction (684),
not as the rule.

### 4.4 Do both branches gate C1 acceptance, in every place acceptance is defined?

Yes. The dependency is stated at:

- 137 — prerequisites: "D10(a) and D10(b) must both be dispositioned before C1 acceptance";
- 156 — direct-`unload` unknown test: "C1 is not accepted until D10(a) is dispositioned";
- 159 — direct-`Cancel` unknown test: "C1 is not accepted until D10(b) is dispositioned";
- 177 — acceptance: "D10(a) and D10(b) must both be dispositioned before C1 acceptance; this list
  does not freeze an unknown-snapshot code";
- 433 — C7 prerequisites carry both forward;
- 630 — §9.1: "C1 cannot be accepted until D10(a) and D10(b) are dispositioned";
- 633 — §9.1: C7 depends on "the D3/D10(a)/D10(b) dispositions";
- 684 — ledger: "Before C1 acceptance (both branches)".

The requirements matrix (598) records the same state without freezing a result: "Unknown capability
knowledge → no optional-operation spawn; D10(a)/D10(b) results not frozen." The record asserts the
same boundary in four places (49, 103, 133, 165); the §2 scope note (51–52) separately states that
the owner "may choose the same or different result for each."

Consistency check on the two test sites: the tests assert structure only — no frame, no spawn, no
asserted result code. Because C1's acceptance is gated on both dispositions, no acceptance run can
occur while either branch is open. Requiring both, rather than either, is what prevents a partial
answer from being read as closure.

One redundancy, not a defect: C7's prerequisites repeat the D10 gate even though C1 acceptance
precedes C7 in the critical path, so the gate is already satisfied by construction. Listing it again
is harmless and arguably clearer than relying on transitivity.

### 4.5 Does the structural test over-constrain the owner's eventual answer?

No. This matters because D10(b) must remain genuinely answerable for any reasonable disposition. The
unknown-snapshot tests (156, 159) and the regression test (163) assert:

- no operation frame is written,
- no helper is started to discover the capability,
- **no result code is asserted.**

Every candidate direction in the ledger (684) — an explicit semantic extension of an existing code,
treating an unverified capability as absent, or deferring the branch — is compatible with those three
assertions. Only a disposition that required writing a frame would conflict, and A2-D1 forecloses
that on the absent path. The test surface is therefore answer-agnostic, which is the correct posture
for a branch whose result is deferred to the owner.

## 5. Round-3 findings in the current text

Verified against the proposal itself, not the record's claims about it.

- **M5 (major, D9 mis-scoped) — closed.** D9 now asks only "whether status/catalog/health observation
  resets the 120 s idle deadline" (36, 360, 684). Status-from-`absent` start is recorded as
  contract-settled by §17.7.2 and excluded from D9 (36, 360, 684), and suppression of poll-driven
  startup is assigned to the application layer rather than to a lifecycle exception (360, 684).
  C7's `ProviderStatus` bullet (454) still reads "may start/handshake; no load", which is now
  consistent with D9's narrowed scope instead of contradicting it. C5's observation-reset tests and
  acceptance remain gated on D9 (343, 362, 374), and status-from-`absent` is verified at C7 as
  contract behavior, "not gated on D9" (374). This is the repair the round-3 audit recommended.
- **m1 (budget wording) — closed.** C5 acceptance (374) now says "Budget permits three restarts within
  10 minutes; a fourth attempt exhausts the budget," matching §17.11.5's "maximum 3 consecutive
  restarts." The phrase "exactly at 3" is gone from the proposal; it survives only in the record's
  description of what was changed (101).
- **m2 (no structural direct-`Cancel` unknown test) — closed.** The test now exists at 159 with the
  same structural-only boundary as 156, and the independent review (175) requires the two direct
  results to be kept separate from the A2-D9-settled path.
- **m3 (§9.1 omitted C1's D10 gate) — closed.** §9.1 (630) now reads "C1 cannot be accepted until
  D10(a) and D10(b) are dispositioned; C2 has no open design gate."
- **m4 (undefined "unsupported-operation frame") — closed.** The term is removed from the proposal
  and replaced by a stated absence: "No separate unsupported-capability frame or wire message exists;
  capability absence is expressed, if at all, by the A2-D1 result code and never by a frame" (123;
  also 122, 163, 175, 177).
- **m5 (watcher call site unnamed) — closed.** C1's `client.go` bullet (142) now names it: "the
  existing cancellation watcher—the call site that currently invokes exported `Cancel`—must call that
  internal operation instead." The record's §2 (36) explains why that naming matters, which is the
  right way to record a fix to a review finding.

## 6. Regression sweep

Spot-verified at current line numbers, because the pass rewrote the proposal in place and shifted
every line number:

| Item | Status | Where |
|---|---|---|
| F1 direct `unload`/cancellation split | Preserved, strengthened by the named watcher call site | 121–123, 142, 163 |
| F2 both `unavailable` cases true for `Available` | Preserved | 205, 214, 233 |
| F3 `Loaded` is a projection, single writer | Preserved | 227 ("Project only; do not write"), 568 |
| F4 Windows evidence only from a real child | Preserved | 506, 642 |
| F7 C4 intent diagnostic | Preserved | 304, 296–298 |
| F8 cap counts actual processes; registry keyed by process/job handle | Preserved | 467–471, 575, 583 |
| D7 settled by §17.11.5 | Preserved | 357, 684 (D7 row) |
| D1a gates C3 | Preserved, consistent at all sites | 253, 261, 288, 433, 570, 576, 631, 674 |
| C1–C8 order and inert C1–C6 boundary | Preserved | 16–27, §11 |
| Write set | Two files only; contract and all prior artifacts byte-identical | §2 above, record §13 |

Gap and decision counts: the record's ledger (§11) and the proposal's §10 table agree on the set
{D1a, D1b, D2, D3, D4 folded, D5, D6 applied, D7 closed, D8 deferred, D9 open, D10 open as two
branches} plus O-1–O-7 unchanged. See §7 below for one arithmetic defect in the proposal's §1
summary of that set.

## 7. New findings

### n1 — minor: the §1 owner-decision count no longer sums

Proposal line 36 reads: "It identifies three contract-surface gaps (D1a, D3, D10) and **four separate
owner decisions** (D1b, D2, D5, D9); **D10 is also an owner decision** because the contract is silent."

The two clauses are individually correct and jointly inconsistent. D10 is classified both as a
contract-surface gap and as an owner decision, so the owner-decision set is {D1b, D2, D5, D9, D10} —
five, not four. The record agrees with the five-item reading: its ledger row calls D10 an "Open owner
decision" (133), and its §2 states the same expanded scope (49, 51).

Not blocking: no gate, test, or acceptance boundary depends on the count, and every individual DECISION
citation is correct. It matters because §1 is the summary a reader will quote into a decision brief,
and "four owner decisions" will understate the gate count by one.

Repair, if wanted: "four further owner decisions (D1b, D2, D5, D9), with D10 as a fifth because the
contract is silent" — or drop "four" and say "five owner decisions (D1b, D2, D5, D9, D10), one of
which is also a contract-surface gap."

### n2 — minor: "deferring the branch" is not a result

The D10 ledger row lists three unselected candidate directions: "an explicit semantic extension of an
existing code, treating an unverified capability as absent, or deferring the branch" (684). The first
two name a caller-facing outcome. The third names a scheduling choice, not a behavior: if the owner
selects it, C1's unknown-snapshot direct branches would still have no defined result, and C1's
acceptance text (177) would be satisfied by a disposition that leaves the runtime behavior undefined.

Not blocking, and possibly intentional — "deferring the branch" may be shorthand for "keep the
structural behavior and ship no result," which is coherent with A2-D1. But a disposition that only
defers should still name what the code returns in the interim, or should explicitly forbid C1
acceptance until a result is chosen. Worth resolving when the owner answers, not now.

## 8. Stage matrix after this pass

| Stage | Status |
|---|---|
| C1 | Implementable and fake-helper testable. **Acceptance blocked by D10(a) and D10(b).** No design defect remains. |
| C2 | No open design gate; clear. |
| C3 | Blocked by D1a. |
| C4 | No open design gate; depends on C1–C3 for sequencing only. |
| C5 | Blocked by D9 (observation-reset half). |
| C6 | No open design gate. |
| C7 | Blocked by D1a (already required before C3), D1b, D2, D3, and the D10 branches by transitivity. |
| C8 | Blocked by C1–C7 and requires separate runtime-execution authorization. |

Nothing in this pass changes the authorization state: no stage is authorized, C1 and C2 remain
unauthorized, and Stage 3 implementation still requires explicit stage-specific authorization.

## 9. Decision ledger status

| ID | State after this pass |
|---|---|
| D1a | Open contract-surface gap; gates C3. |
| D1b | Open; gates C7. |
| D2 | Open; gates C7. |
| D3 | Open (caller-facing rejection code); gates C7's cap-response claim. |
| D5 | Open; gates UI integration. |
| D7 | Closed by §17.11.5. |
| D8 | Deferred (B12a/B12b). |
| D9 | Open, narrowed to observation-driven deadline reset. Status-from-`absent` start is contract-settled and excluded. |
| D10(a) | Open owner decision; gates C1 acceptance. |
| D10(b) | Open owner decision; gates C1 acceptance. Independent of D10(a). |
| O-1–O-7 | Unchanged. |

No decision was resolved by assumption in this pass, and this audit resolves none.

## 10. What was checked, and how

- Recomputed SHA-256 and line/byte counts for both subjects, the contract, all three prior audits, and
  the readiness investigation; compared against the record's §1 table row by row.
- Read both revised documents in full.
- Re-read the three contract decisions the split depends on at their adopted lines: A2-D1 (2502),
  A2-D9 (2510), A2-D12 (2513).
- Enumerated every `D10(a)`/`D10(b)` reference in both documents and checked each for a consistent
  branch assignment and a consistent gate.
- Searched for residual instances of the m1 and m4 defects ("exactly at 3", "unsupported-operation")
  in the proposal; both are absent there and appear only in the record's account of the repair.
- Checked for bare `D10` references that could let the two branches be collapsed; all are collective
  ("D10", "outside D10") and none assigns a result.
- Verified the write set from modification times and hashes: only the proposal and the record changed.

## 11. Limitations

1. **Not independent.** The reviewer wrote the round-1, round-2, and round-3 audits and recommended
   the M5 repair and the m1–m5 fixes that this pass implements. Verifying another author's
   implementation of your own recommendation is the intended role here, but it is not a third-party
   check. A separate reviewer remains warranted before C7, as stated in the round-3 audit.
2. **No byte-level diff.** Both documents were corrected in place and the intermediate revision
   `7b73635d…` is gone. Changes were located in the current text and checked there; the claimed
   revision chain is consistent but its middle element is unverifiable.
3. **Documentation only.** No build, test, replay, benchmark, or process launch was run, and none of
   the cited code paths were executed. Statements about C1's behavior are statements about what the
   proposal specifies, not about what any implementation does.
4. **The index was not updated, by instruction.** `README.md` still ends at §28, which describes the
   round-3 audit of `1fb449ba…`/`ccb28d57…`. That entry remains accurate for the revision it covers,
   but the index does not yet describe this pass. An index entry is a separate maintenance task; this
   audit intentionally leaves both `README.md` and the round-3 audit untouched.

## 12. Non-authorization

This document is a review artifact. It authorizes no implementation stage, resolves no owner
decision, and does not disposition D10(a) or D10(b). C1 remains gated on both branches; C3 on D1a;
C5 on D9; C7 on D1b, D2, D3, and the branches; C8 on separate runtime-execution authorization.
