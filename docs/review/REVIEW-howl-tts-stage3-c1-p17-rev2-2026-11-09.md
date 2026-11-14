# Independent Review — P17 Rev2 Corrected Proposal (2026-11-09)

**Review status:** COMPLETE — **BLOCKING FINDINGS; NOT AN APPROVAL.**

**Reviewed target:** `docs/review/PROPOSAL-howl-tts-stage3-c1-evidence-semantics-p17-rev2-2026-11-09.md`

**Target identity at review:** **513 lines**, SHA-256 `49e9ecc2c812a361d03e739fbb571bfacc61f93fe45cffd30a8cb0bddb1055d6`

**Reviewed revision:** corrected P17 rev2 candidate, described by its header as corrected from the 504-line candidate `e1a5aaa2…`.

**Review scope:** internal coherence of the evidence model and projection rules; Gate 3/4 state coverage; active-state exclusion; CP-1 route coverage; the eight mandatory reviewer challenges in the target’s §11; and an anchor audit of the cited source behavior.

**No edits during review.** The reviewed proposal was not changed. The source tree was read only. No tests, builds, benchmarks, runtime experiments, contract edits, register edits, or implementation were performed.

**Verdict:** the proposal is **not ready for approval**. The Gate 3/4 structure still admits an unresolved-establishment state into `absent`, Gate 3 conflates “not published” with “internal-only,” and Case A bypasses both CP-1 and the active-state exclusion for established terminals. P17 remains unapproved; contract rev3 remains frozen; C1 remains unauthorized.

---

## 1. Blocking findings

### R-1 — BLOCKING: an advanced classification with a resolved threshold but unresolved establishment can pass Gate 4 and reach `absent`

**Target locations:** §8 P4-A gates 3–5 (lines 291–293), P4-A narrative and projection (lines 295–307), B-A (line 328), B-B (lines 331–335).

Gate 3 applies only to an **established** prior failure classification. Gate 4 applies only to a classification **advanced by the record but not established** with an **unresolved justification threshold**.

Those are not exhaustive over the plausible classification states. Consider:

- the record advances a candidate prior failure classification;
- the governing justification threshold is settled (the applicable rule is identified and its acceptance criteria are resolved);
- whether the available evidence satisfies that threshold is **not** resolved, so the classification is neither established nor affirmatively rejected.

Under the literal gate text:

- **Gate 3:** passes, because no **established** classification exists.
- **Gate 4:** passes, because the justification threshold is not unresolved; only the evidence application is unresolved.
- **Gate 5:** passes if no prior public failure value exists.
- **Gates 1–2:** can pass with orderly shutdown and no terminal.

The “All five gates hold: project `absent`” rule then applies. That directly contradicts the proposal’s own statement that unresolved establishment keeps the branch **OPEN** and the C-02 summary’s promise that `absent` is barred while a failure classification remains unresolved as to establishment.

The surrounding case text does not reliably close the gap:

- **B-A** is labeled “Bare fatal observation only,” but its orderly-shutdown bullet applies when there is “no justified prior failure classification” and projects `absent`. An advanced-but-unestablished classification is not a bare observation, yet B-A’s predicate can treat it as having no justified classification.
- **B-B** is defined as “Fatal observation with an independently justified failure classification,” but its introductory paragraph says an advanced classification with unresolved “establishment” remains OPEN under B-B. A classification that is not established is not “independently justified,” so B-B’s own definition does not cover it.
- The “all five gates hold” projection is stated before the case branches and is not conditioned on a resolved establishment determination.

This is not a wording preference. The state is reachable whenever a reviewer can state the threshold rule but cannot decide whether the evidence meets it, which is exactly the evidentiary uncertainty the proposal is supposed to preserve.

**Why blocking:** it permits `absent` over an unresolved prior failure classification, violating the proposal’s asserted evidence standard and C-02.

**What would clear it:** make the establishment state three-valued and explicit — `established`, `rejected`, `unresolved` — and make Gate 4 fail for `unresolved` regardless of whether the threshold rule itself is settled. Alternatively, add an unconditional gate that bars `absent` for any advanced classification without a resolved establishment determination. B-A/B-B must then be aligned with that predicate.

### R-2 — BLOCKING: Gate 3 conflates “not published” with “established internal-only”

**Target locations:** §8 gate 3 (line 291), public-effect dimensions (line 295), B-B (line 334), C-B (line 349).

Gate 3 says publication status must be resolved as either:

- `published`, with the public value identified; or
- `not published/internal-only`.

Those are not the same determination. “Not published” can mean the record contains no evidence of publication; “internal-only” is a positive determination that publication is excluded by the accepted classification rule. The proposal elsewhere recognizes the distinction:

- B-B requires that publication be “neither established nor excluded” to remain OPEN.
- C-B requires that publication be “excluded/internal-only.”
- The public-effect dimensions list “not published/internal-only” as one category while also listing “unresolved.”

If an established classification has no established publication and no established exclusion, the literal gate can treat it as “not published,” conclude there is no public value, pass Gate 5, and project `absent`. That treats missing publication evidence as proof of non-publication, contrary to the proposal’s own “missing evidence is not automatically proof” constraint and its statement that a public failure value must have evidence of publication.

**Why blocking:** it can erase a prior public value by evidentiary silence and it undermines the C-B distinction between established internal-only classifications and merely asserted or unknown publication status.

**What would clear it:** split the status into `published` / `excluded-internal-only` / `unresolved`. Gate 3 should pass only for the first two; “not published” without an established exclusion is `unresolved` and must fail Gate 3. Define whether an internal-only determination is a policy classification, an evidence finding, or both.

### R-3 — BLOCKING: Case A does not route an established terminal carrying a prior public `crashed` value through CP-1

**Target locations:** Case A (lines 310–317), Gate 2 projection (line 304), B-D (line 339), public-`crashed` section (line 364).

Case A says:

- preserve the terminal;
- later evidence cannot replace it;
- a later `crashed` classification may be recorded only internally and does not become public;
- public projection follows the terminal, and the revision “does not invent a new terminal projection rule.”

It never mentions CP-1. If the established terminal itself carries a previously published public `crashed` value, Case A preserves that terminal by OD-3 and effectively decides the CP-1 preservation question in favor of preservation. That contradicts:

- S13’s distinction between selecting a new public `crashed` and preserving a prior value under CP-1;
- S14’s statement that prior public `crashed` preservation is not settled;
- the public-`crashed` section’s statement that CP-1 is OPEN;
- the user-required coverage of an established terminal carrying a prior public `crashed` value.

The bypass is reachable through Gate 2 failure and through B-D, both of which route to Case A. The ordinary CP-1 routes in S7/S8/O3/O4/O5/B-A/B-B/B-C/C-A therefore do not cover the terminal route.

**Why blocking:** it decides CP-1 by terminal preservation without acknowledging the open interaction, contrary to the explicit requirement that CP-1 cover terminal cases.

**What would clear it:** add an explicit Case A clause for a terminal whose public value is `crashed`. It must state whether OD-3’s terminal immutability or CP-1’s unresolved preservation question controls, and it must not assert preservation or exclusion while CP-1 is OPEN. If the answer is “terminal controls,” that is a decision requiring separate authority, not an implication.

### R-4 — BLOCKING CLARIFICATION: Case A can preserve an active terminal because the terminal set remains OPEN

**Target locations:** §3.4 terminal publication (lines 126–135), Case A (lines 310–317), Gate 5 active-state clause (line 293), C-C (line 351).

Case A preserves any established terminal decision. The exact terminal set is explicitly **OPEN** in §3.4. Gate 5 and C-C exclude active operational states after a confirmed exit, including active `degraded`, but Case A is reached through Gate 2 and does not apply Gate 5’s active-state filter.

If the terminal set can include `ready` or active `degraded`, Case A preserves an active operational state after exit, contrary to OD-4’s accepted rule that active states cannot survive a confirmed exit. If the terminal set is intended to exclude all active states, that exclusion is not stated and cannot be verified from the current proposal.

**Why blocking clarification:** the requirement is unconditional active-state exclusion after confirmed exit. The current terminal set is undefined, so the guarantee is not established.

**What would clear it:** define the terminal set to exclude active operational states, or add a Case A guard that an active operational terminal is not preservable and routes to C-C. The former is cleaner; the latter should use the same “not applicable/cannot be preserved” language as Gate 5.

---

## 2. Gate 3/4 state table

This table applies the gates literally. “Orderly” means Gate 1 holds; “no terminal” means Gate 2 holds; where a row says `absent`, no other accepted condition is assumed.

| Prior classification state | Publication status | Threshold state | Gate 3 | Gate 4 | Literal result | Disposition |
|---|---|---|---|---|---|---|
| None | n/a | n/a | pass | pass | `absent` if orderly/no terminal | PASS |
| Advanced | n/a | unresolved | pass (not established) | fail | `OPEN` | PASS |
| Advanced | n/a | resolved, evidence application unresolved | pass (not established) | **pass** (no unresolved threshold) | **`absent`** | **R-1 BLOCKING** |
| Advanced, affirmatively rejected | n/a | resolved as not met | pass | pass | `absent` may be permissible | Not defined; see R-1 |
| Established | unresolved | n/a | fail | pass | `OPEN` | PASS |
| Established | “not published” but exclusion not established | n/a | **pass** under literal slash | pass | **`absent`** | **R-2 BLOCKING** |
| Established | excluded/internal-only | n/a | pass | pass | `absent` if orderly/no other failure | PASS |
| Established | published, value identified, compatible/applicable | n/a | pass | pass | Gate 5 fails; preserve | PASS |
| Established | published `crashed`, compatible/applicable | n/a | pass | pass | CP-1 OPEN | PASS |
| Established | published `crashed`, incompatible | n/a | pass | pass | historical; `absent` if orderly | PASS |
| Established terminal | any | n/a | pass | pass | Gate 2 fails → Case A | **R-3/R-4** |
| No orderly shutdown | any | any | fail | any | `unavailable` fallback | PASS |
| Active operational state | n/a | n/a | pass | pass | Gate 5 excludes; separate interpretation may prevent `absent` | PASS outside Case A |

**Reading of the table:** the only row that reaches `absent` while an establishment question remains unresolved is the third row. That is the minimal counterexample to P4-A’s claimed eligibility standard.

---

## 3. Active-state route audit

The requirement is that no active operational state survives a confirmed exit, including active `degraded`.

| Route after confirmed exit | Active-state treatment | Result |
|---|---|---|
| P4-A Gate 5 | Active state is not applicable and cannot be preserved; a separately asserted post-exit interpretation prevents `absent` while unresolved | PASS |
| C-C | Active `ready`/`degraded` not preservable; no active state survives | PASS |
| B-C | Active state not preservable; C-C controls | PASS |
| C-A | Active `degraded` is not C-A; C-C controls | PASS |
| B-A | No prior public failure; `absent` or `unavailable`; no active preservation | PASS |
| B-B | Same; no active preservation | PASS |
| C-B | Internal-only; no active public preservation | PASS |
| P4-D | No applicable prior controlling state; no active preservation | PASS |
| Case A | Preserves the established terminal without a value/active-state check | **R-4** |
| B-D | Routes to Case A | **R-4** |
| P3/O6–O12/O17/O18, no exit | Last justified current state may remain; this is before a confirmed exit | PASS |

**Conclusion:** the active-state exclusion is sound for the nonterminal and non-Case-A routes. It is not unconditional because Case A can preserve an active terminal while the terminal set is OPEN.

---

## 4. CP-1 route audit

The CP-1 question is whether a previously published public `crashed` value may be preserved under OD-4 or must remain historical evidence under OD-5. It remains OPEN and must not be decided by a route.

| Route | Prior public `crashed`, otherwise compatible/applicable | Prior public `crashed`, incompatible | Compatibility unresolved | Disposition |
|---|---|---|---|---|
| S7 | CP-1 OPEN | historical; `unavailable`/fallback | OPEN | PASS |
| S8 | CP-1 OPEN | historical | OPEN | PASS |
| O3 | CP-1 OPEN | historical | OPEN | PASS |
| O4 | CP-1 OPEN | historical | OPEN | PASS |
| O5 | CP-1 OPEN | historical | OPEN | PASS |
| B-A no orderly | CP-1 OPEN | historical | OPEN | PASS |
| B-B no orderly | CP-1 OPEN | historical | OPEN | PASS |
| B-C | CP-1 OPEN | historical | OPEN | PASS |
| C-A | CP-1 OPEN | historical | OPEN | PASS |
| Gate 5 | CP-1 OPEN | historical | OPEN | PASS |
| Public-`crashed` section | CP-1 OPEN | historical | OPEN | PASS |
| Case A / B-D / Gate 2 | **No CP-1 clause; terminal preservation decides the question** | same | same | **R-3 BLOCKING** |
| P4-D | Prior public compatible `crashed` is not “no prior controlling state”; Gate 5 controls | historical where incompatible | OPEN | PASS if routed through Gate 5 |
| C-B | Prior public value is not established/internal-only; no CP-1 | n/a | n/a | PASS |

**Conclusion:** CP-1 coverage is broad and explicit for ordinary preservation routes, but not for an established terminal carrying a prior public `crashed` value. R-3 is the specific gap the owner asked the reviewer to test.

---

## 5. Mandatory reviewer challenges — explicit dispositions

| # | Challenge | Disposition | Finding |
|---|---|---|---|
| 1 | Verify every retained line and code claim against source | **PASS for cited anchors; no runtime verification** | The anchors in `client.go`, `supervisor.go`, `protocol.go`, `discovery.go`, `errors.go`, and `provider.go` are consistent with the cited behavior. No material line drift or contrary code claim was found. This is a source read, not a test. |
| 2 | Challenge each fatal assignment, especially O6, O12, O15, O16, O18 | **PASS with open-case preservation** | O6 maps `internal_error` to `helper_crashed` and does not kill that code; O12 returns `helper_crashed` on write failure without termination; O15 maps handshake failures to `incompatible`/`unavailable`; O16 matches the nine `validateHello` exits and leaves U8 open; O18 writes `crashed` while the process may be alive and the proposal correctly forbids treating that as an exit. No source contradiction. |
| 3 | Review U1–U8 as open; no tests/runtime | **PASS** | U1–U8 remain OPEN; CP-1 is separately OPEN; no test/build/runtime was performed. |
| 4 | Terminal immutability must not erase late evidence; old-generation evidence must not modify a replacement | **PARTIAL — R-3/R-4** | §3.4 and P9 preserve late evidence and generation isolation. Case A nevertheless preserves a terminal without CP-1 or active-state qualification, so the terminal route is not fully consistent. |
| 5 | No public projection uses `outcome_unknown`, `quarantined`, or `terminated_by_action`; thirteen-state set unchanged | **PASS** | The proposal keeps all three internal and lists the thirteen public states unchanged. |
| 6 | E1–E9 visible and unresolved; no implementation claim | **PASS** | E1–E9 are present and unresolved; the proposal does not claim implementation evidence. |
| 7 | Finding-level dispositions; do not approve P17 | **PASS** | This review gives blocking/non-blocking dispositions and does not approve P17. |
| 8 | Record exact path, revision, line count, SHA-256; no plan-level substitute | **PASS** | Target: 513 lines, `49e9ecc2…`; the review is bound to that exact artifact. |

---

## 6. Source-anchor audit

**Repository revision:** `9c7d2a1`.

The following anchors were checked against the working tree:

- `internal/audio/tts/helper/client.go`
  - `248–254`: handshake failure calls `Terminate`, maps `protocol_incompatible`/`protocol_minor_unsupported` to `incompatible`, otherwise `unavailable`.
  - `299–313`: `internal_error` is rewritten to `helper_crashed`; `failAll`; only `invalid_frame`, `frame_too_large`, and `protocol_violation` call `Kill`.
  - `323`, `331`, `345`: protocol-violation paths call `Kill`.
  - `421`: synthesis write failure returns `helper_crashed`.
  - `439`, `441`: request timeout and cancellation.
  - `578–594`: load-path wire/protocol failures set `unavailable`.
  - `654–674`: `Shutdown` calls `Terminate` and then clears lifecycle/transport state; reaping is not established.
  - `756`: synthesis write failure sets public `crashed` while the process may be alive.
  - `771`, `783`, `793`, `796`: timeout/cancel/ready/degraded handling.
  - `861`: cleanup calls `Terminate` and ignores its error.
- `internal/audio/tts/helper/supervisor.go`
  - `132`, `147`: discovery/start failures set `unavailable`.
  - `176–190`: `observeExit` maps `stopping` to `absent` and other exits to `crashed`; nonzero status is not classified separately in the stopping branch.
  - `233–241`: `SetState` overwrites `lastFailure`.
  - `246–293`: `Terminate`/`Kill`; `Kill` aliases `Terminate`; no intent or outcome record.
  - `283–286`: termination timeout.
  - `324`, `351–356`: `ProcessAlive` derives false for `crashed`/`unavailable`.
- `internal/audio/tts/helper/protocol.go`
  - `502–545`: nine `validateHello` failure exits; missing mandatory capabilities produce `protocol_violation` in the source.
- `internal/audio/tts/helper/errors.go`
  - `15–21`: `Failure` fields.
- `internal/audio/types/provider.go`
  - `23–35`: thirteen public provider states, including `crashed`.
- `internal/audio/tts/helper/discovery.go`
  - `39`, `42`, `61`: discovery failure paths.

**Result:** no material source discrepancy was found in the cited anchors. The proposal’s source citations are usable for review, but this does not verify the proposed semantics or any implementation.

---

## 7. Non-blocking findings and confirmations

- **C-D placement:** corrected. C-D is now a sibling `###` transition to P4-D, and Case C states three subcategories. **PASS.**
- **P4-A labels:** P4-A is consistently labeled proposed; C-02, C-11, and Case D use proposed language. **PASS.**
- **S5 inline tracking:** B-A, B-C, Case C, and P4-A carry the proposed non-application and `S5-scope OPEN` qualifiers. **PASS.**
- **Established versus asserted internal-only:** C-B correctly distinguishes established internal-only from merely asserted internal-only. **PASS at C-B**, but Gate 3’s “not published/internal-only” still undermines the distinction (R-2).
- **Gate 3/4 overlap:** the previous overlap is removed. **PASS for overlap**, but the new partition is not exhaustive (R-1).
- **CP-1 tracking:** S13, S14, Section 10, S7/S8, O3/O4/O5, B-A/B-B/B-C/C-A, Gate 5, C-11/C-14, and Section 13 track CP-1. **PASS outside Case A.**
- **Active-state tracking:** Gate 5, C-C, B-C, C-A, B-A/B-B, C-B, and P4-D handle the nonterminal routes. **PASS outside Case A.**
- **No public state change:** the thirteen-value set is unchanged. **PASS.**
- **No implementation evidence:** no tests, builds, or runtime experiments. **PASS as a limitation.**

---

## 8. Disposition and next step

**Blocking findings:** R-1, R-2, R-3, and R-4.

**Non-blocking:** the source anchors are consistent; C-D, P4-A labels, S5 tracking, and CP-1 tracking outside Case A are adequate.

**Required next step:** a separately authorized correction pass must close R-1 through R-4, then the artifact must be rebound by a new line count and SHA-256 before any further review or approval attempt. This review does not authorize that pass, does not approve P17, and does not authorize C1.

**Current authority state:** P17 unapproved; contract rev3 frozen; C1 unauthorized.

---

## 9. Review identity

- Reviewed target: `docs/review/PROPOSAL-howl-tts-stage3-c1-evidence-semantics-p17-rev2-2026-11-09.md`
- Reviewed target lines: **513**
- Reviewed target SHA-256: `49e9ecc2c812a361d03e739fbb571bfacc61f93fe45cffd30a8cb0bddb1055d6`
- Contract rev3: 183 lines, `ed0da1a75d4f62874f873979e6667f66f19541a7523ff96369f2b44e8391195e` — unchanged
- Original P17: 166 lines, `01d887c18700d63070f1a2649659281ec1e875a21b49375adea12f8f11d917bf` — unchanged
- Decision register: 434 lines, `3c64797bfdd48ddacca95e6a0def88d960e4104aa3fc3034a25dc539d97cbc99` — unchanged
- Review method: document review plus source-anchor audit; no tests, builds, runtime experiments, or edits to the reviewed proposal.

**This review is not an approval. P17 remains unapproved. Contract rev3 remains frozen. C1 remains unauthorized.**
