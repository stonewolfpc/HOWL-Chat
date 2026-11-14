# AUDIT — Howl thread/hardware/performance requirements, revision r7

**Reviewer:** Mimo (independent of the r7 author)
**Date:** 2026-11-07
**Mode:** Read-only independent review. No subject file was edited, repaired, reconciled, staged, committed, or executed. No build, test, replay, measurement, or benchmark was run. No adoption, implementation, or runtime authorization is implied.
**Subject:** `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r7-2026-11-07.md`
**Subject SHA-256:** `58cd4656c35407691f4b590d095c28340eea6710cc28debb34abde36f387e772`
**Subject size:** 715 lines / 97,790 bytes, LF-only (0 CR bytes)

## Verdict

**Not adoptable as written.** One blocking defect, four minor, eight observations.

The three targeted corrections are substantially implemented, and two of the three are correct everywhere I checked. The crash-restart rule is now the contract's TTS-only rule (`§17.11.5`, line 2105) at every operative locus, and A2-D12 is confined to `unavailable` at all eight mentions. §7.3 reproduces the seven original criteria verbatim. The r4 findings are all addressed.

The blocking defect is an attribution error introduced by the same correction: r7 assigns the health/status/catalog observation bar to the proposed plan and states it is "not attributed to the contract" (r7 L106). The contract says otherwise in its own words — "A health response is never a transition source" (contract L1801–1802) — and the `§17.11.5` TTS-only restart rule entails the same bar for status and catalog requests. Because the only executed evidence for the bar is conditioned on plan C5 (r7 L564), a stage can be accepted without demonstrating a contract-fixed behavior.

This is a one-sentence correction in four places. It is not a design failure, and it does not recur the r6 restart-breadth defect.

## Severity scale

- **Blocking** — would freeze a wrong interpretation into stage acceptance, or let a stage pass while a contract-settled behavior is unimplemented or unverified.
- **Major** — authority, scope, or allocation defect that does not by itself produce a wrong pass.
- **Minor** — citation, ambiguity, or drafting defect.
- **Observation** — not a defect; recorded because it affected the review.

Observations are numbered `OB-n` to avoid collision with the contract's owner decisions O-1–O-7.

---

## Blocking

### R7-1 — The observation bar is contract-settled, but r7 attributes it to plan C5 and makes its evidence conditional

**Target:** plan-C5 attribution; restart breadth; acceptance evidence.

**Evidence.** r7's three-layer table states: "The plan's C5 observation bar and no-background-loop rule are therefore plan-level, not contract-level" and, in the layer row, "Retained only as a proposed-plan C5 condition. It binds only if the adoption record incorporates the plan as the stage plan. **It is not attributed to the contract or to A2-D12**" (r7 L106). The same attribution appears in C1: "not health, status, or catalog observation (that bar is proposed plan C5)" (r7 L560); in C5: "Proposed plan C5 separately bars health/status/catalog observation from restarting a helper … that is a plan-level condition, not A2-D12" (r7 L564); in §9.2: "status, health, and catalog do not restart under proposed plan C5" (r7 L591); and in the adoption checklist: "the health/status/catalog bar is proposed plan C5" (r7 L699).

The contract contradicts the "not attributed to the contract" claim directly:

> "`backoff` is a supervisor scheduling state; it is not a provider-loaded state. **A health response is never a transition source.** The process exit observer is the only authority for an unexpected exit." (contract L1801–1802)

And `§17.11.5`, which r7 itself designates controlling, says: "**Restart only on the next explicit TTS request.**" (contract L2105). A `status` or `catalog` request is not a TTS request, so it cannot produce a restart under that rule; a `health` request cannot produce any transition at all. The observation bar is therefore entailed by the contract twice over, not invented by the plan.

r7 cites L2105 correctly (r7 L105) but never cites L1801–1802, the sentence that settles the health half of the bar.

**Why it matters.** C5's required evidence is "no health/status/catalog restart **under proposed plan C5**" (r7 L564). If the adoption record does not incorporate the plan, §2.1 and §1.5 say the plan's criteria are review inputs only, so the only executed evidence for a contract-fixed behavior becomes inapplicable. The stage can then be accepted with no test that a health or status probe during `crashed` fails to restart the helper — even though the contract forbids it. §14's second adoption item ("no unresolved authority, contract, lifecycle, or safety contradiction") is therefore not satisfied.

The start rule itself is not wrong: C1's operative sentence is "the next explicit TTS request under §17.11.5", which already excludes observation-triggered restarts. The defect is the source attribution and the conditional evidence, not the behavior.

**Fix.** Attribute the bar correctly in all four places: health to contract `§17.7.2` L1801–1802; status and catalog to `§17.11.5` L2105. Keep the plan's C5 sentence as a redundant restatement whose plan-level additions are the no-background-loop rule and the named observation surface, not the bar. Make C5's evidence item unconditional ("no health/status/catalog restart"), and correct the §1.5 layer row and the §14 checklist line.

---

## Minor

### R7-2 — C1's `stopping` row says "only exit 0 to `absent`", omitting the contract's path into `crashed`

C1 states: "From `stopping`: no direct start; only exit 0 to `absent`" (r7 L560). The contract has two further edges out of `stopping`: grace expiry → `crashed` ("`KillTree`; no replay", contract L1794), and nonzero exit → `crashed` (A2-D8 and the "any live state / unexpected `Wait()` exit" row, L1795). Since a `crashed` helper is exactly the state from which the restart rule operates, a stage reproducing C1's map could treat shutdown as a restart-dead end.

The start prohibition is right — no direct start from `stopping` — but "only" makes the map factually incomplete in the one document whose purpose is a complete state map. C4 covers exit classification, so this does not create a wrong pass on its own.

**Fix.** "From `stopping`: no direct start. Exit 0 → `absent`; nonzero exit or grace expiry → `crashed`, which then follows the crash-restart rule. No start before `absent` or `crashed` is observed."

### R7-3 — HW-1.8 requires a measurement that §6.5 does not authorize

HW-1.8 states: "The overhead of safety instrumentation and verification must be measured separately from contention-induced degradation" (r7 L401). §6.5 closes with: "No performance measurement or benchmark is authorized" (r7 L468), and §6.4 says a stage "does not execute the suite or use results as acceptance evidence while A-1 is open" (r7 L443).

Instrumentation-overhead measurement is a performance measurement. With A-1 open — the current state — the two requirements cannot both be satisfied, and HW-1.8's bullet has no authorized means of compliance. This is a scope conflict created by the A-1 deferral, not a safety defect; the safe reading (run nothing) is clear, but a stage report can cite either text.

**Fix.** Carve out safety-verification measurement from the prohibition: "No performance-gate measurement, benchmark, or PR-suite execution is authorized; instrumenting and measuring the overhead of safety verification for correctness is permitted and is not a performance claim." Alternatively make HW-1.8's bullet conditional on A-1 authorization.

### R7-4 — The design-incorporation condition is cross-referenced to two different sections

§2.1 says the plan's criteria apply "only to the extent the adoption record expressly incorporates the plan as the stage plan **under §2.2**" (r7 L120). Every §7.1 stage row says "(if incorporated under **§2.1**)" (r7 L513–522), and the §8 preamble repeats "only where the adoption record incorporates it **under §2.1**" (r7 L556). The operative mechanism is the adoption-record content in §2.2; §2.1 states the role, not the incorporation act.

No behavioral difference follows, but this is the same conditional that R4-4 and R4-6 turned on, and a reader checking whether the plan is incorporated should land on the adoption-record clause. Harmonize all references to §2.2.

### R7-5 — The two newest review inputs are the only ones not hash-bound

The header lists as review inputs "the 2026-11-07 pre-review assessment of r5, and the 2026-11-07 post-r6 restart assessment," and §1.4 summarizes "The 2026-11-07 pre-review of r5" (r7 L95–97) as the source of P1–P3, while §1.5 introduces "the post-r6 correction S1." Neither has a path or SHA-256, and no such artifact exists in `docs/review/`; a full-text search finds only r6 and r7 themselves.

Every other review input — the r2, r3, and r4 audits — is recorded in §2.2 with an exact hash, line count, and byte count. The evidence chain for r6's and r7's newest corrections is therefore the one part of the provenance that the owner cannot reproduce, which matters because §2.2 says "A hash mismatch in any row reopens review" and §2.4 forbids adoption before review is recorded.

**Fix.** Either record the two assessments as artifacts with hashes in §2.2, or label them explicitly as non-reproducible drafting-time inputs and state that S1's contract citations (which I verified independently) are the operative evidence.

---

## Observations

**OB-1 — The artifact and all thirteen input rows verify.** r7 is 715 lines / 97,790 bytes with SHA-256 `58cd4656c35407691f4b590d095c28340eea6710cc28debb34abde36f387e772`, and the file is pure LF (0 CR bytes). I recomputed every §2.2 row: original 344/25,933; r2 447/36,025; r2 audit 221/29,130; r3 513/47,710; r3 audit 226/30,143; r4 603/68,813; r4 audit 258/30,900; r5 686/88,436; r6 700/95,092; contract 2,586/182,904; plan 711/69,168; correction record 172/16,190; owner-decision instrument 610/45,891. Every hash, line count, and byte count matches.

**OB-2 — The 715-byte upload difference is newline conversion, confirmed.** Adding one CR per line to the workspace artifact yields exactly 98,505 bytes — the uploaded size — and a different digest (`9d1d5e41…`). The workspace file contains no CR bytes. So the upload is an encoding variant of the same text, and neither its size nor its digest can establish the artifact hash. The workspace artifact at `58cd4656…` is the one the adoption record must bind.

**OB-3 — The TTS-only restart rule is consistent at every operative locus.** I checked all "next explicit" occurrences. The crash-restart trigger is TTS-qualified in C1 (L560), C5 (L564), §9.2 (L591), the F2/A-5/R3-5/R4-1 dispositions (L21, L44, L62, L77), the §1.5 layer table (L105), and the checklist (L699). The unqualified phrase survives only where it is explicitly historical (r6's withdrawn P3, L97) or descriptive of `§17.7.2`'s shorthand (L105–L106). The `§17.11.5` quote is exact (contract L2105), including 1/2/5 s, three restarts per key in 10 minutes, and reset after 10 minutes of stable readiness. The r6 restart-breadth defect does not recur.

**OB-4 — A2-D12 is confined to `unavailable` at all eight mentions.** L21, L44, L62, L106, L107, L521, L564, and L709 all treat the explicit-load exception as `unavailable`-only, and L107 states the limitation explicitly. The ordinary `unloaded` path is not restricted: A2-D12's own carve-out that `Synthesize`'s ensure-handshake/load may load an ordinary `unloaded` provider is preserved, and C1's note that a `load` transition in an already-live helper "is not a process start" matches the plan's budget rule (design L355). This target passes.

**OB-5 — §7.3 is a verbatim reproduction.** I compared the seven rows mechanically against the original AC table (`REQUIREMENTS-howl-thread-and-hardware-safety-2026-11-07.md` L248–254): identifier, title, and body text are identical, including AC-2's static-review sentence and AC-3's repeatability sentence. Combined with the non-normative pointers (R7-AC-1, -2, -3, -5, -6, -10) and §2.1's "original governs on any difference," R4-5 is closed as intended.

**OB-6 — Every r4 finding has a corresponding, verified correction.** R4-1: full state map including the crash-restart edge and an explicit "omits the `crashed`/`backoff` edge fails C1" failure example (L560). R4-2: duty to determine, written element/source/cure conditions, owner escalation for source meaning, replacement reviewer, and "not a veto" (§10). R4-3: §7.2 owner acceptance authority with the separate-acceptor exception, plus a sixth matrix column (§7.1 header, L494). R4-4: one coherent design role (L120). R4-5: pointers plus verbatim §7.3. R4-6: conditional design-non-amendment statement, design-status fields in §2.3, and the §2.5 A/B/open consistency table. R4-7: explicit per-category I/O N/A. R4-8: C1 handshake tests retained, trigger-map execution assigned to C5/C7. R4-9: C3 forbids only unsolicited starts and cross-refers the explicit status start. R4-10: O-1 recorded in §13 with the `ProviderCatalog`/`GetTTSVoiceCatalog` distinction. R4-11: provisional value as working ceiling.

**OB-7 — Adoption is expressible with A-1 open and all eight decisions pending.** §2.4 requires A-1 before PR results have effect but permits adoption as a non-performance revision with PR-1–PR-4 inactive. §2.5's three paths are internally consistent; Option B requires a plan-amendment instrument, scope, and hash, and cannot be paired with an unqualified non-amendment statement; Option A leaves the plan unamended; the open path records PR inactive. No adoption record exists, and none should be created on this review.

**OB-8 — §14 double-counts D9.** The checklist line "D9, A-1, and all eight Stage 3 owner decisions remain separate and unresolved" (r7 L710) lists D9 separately from the eight of which it is one, implying nine Stage 3 decisions. §13 states the eight correctly. Drafting only.

---

## Targeted consistency check

| Target | Result |
|---|---|
| **TTS-only restart rule** | **Consistent.** Contract `§17.11.5` L2105 quoted exactly; applied in C1, C5, §9.2, F2/A-5/R3-5/R4-1, §1.5, and §14. `§17.7.2` L1796 is subordinated explicitly and only appears as shorthand. No recurrence of r6's breadth defect. |
| **Plan-C5 attribution** | **Defect (R7-1).** The bar is attributed to the plan and declared "not attributed to the contract" (L106), but the contract bars health transitions outright (L1801–1802) and bars non-TTS restarts (L2105). C5's only executed evidence for the bar is consequently conditional on plan incorporation. |
| **A2-D12 unavailable-only scope** | **Consistent.** All eight mentions limit recovery to `unavailable`; no generalization to `crashed`/`backoff`; the ordinary `unloaded` load path is preserved; `load`-only retries are correctly excluded from the process-restart budget. |

## r4 finding disposition

| r4 finding | Status |
|---|---|
| R4-1 blocking crash-restart prohibition | **Closed** — complete map plus failure example; verified against contract L1793–1799 |
| R4-2 unbounded reviewer withholding | **Closed** — duty to determine, cure conditions, escalation, replacement; withholding is no longer a veto |
| R4-3 no stage acceptance authority | **Closed** — §7.2 plus sixth matrix column; author and reviewer excluded |
| R4-4 design role | **Closed in substance** — one role in §2.1; residual cross-reference in R7-4 |
| R4-5 restatements narrowing originals | **Closed** — pointers plus verbatim §7.3 |
| R4-6 Option B record contradiction | **Closed** — conditional confirmation, §2.3 fields, §2.5 check |
| R4-7 "likely N/A" | **Closed** |
| R4-8 C1 start-test wording | **Closed** — handshakes retained; execution split to C5/C7 |
| R4-9 C3 status-spawn ambiguity | **Closed** — "unsolicited" plus C7 cross-reference |
| R4-10 O-1 omitted | **Closed** — §13 entry |
| R4-11 provisional-limit comparison | **Closed** — provisional value as working ceiling |

## What adoption requires

Textual changes only; no redesign:

1. Correct the observation bar's source in §1.5, C1, C5, §9.2, and the §14 checklist; make C5's evidence unconditional (R7-1).
2. Complete C1's `stopping` row with the `crashed` edges (R7-2).
3. Reconcile HW-1.8's overhead measurement with §6.5's prohibition (R7-3).
4. Harmonize the design-incorporation cross-reference to §2.2 (R7-4).
5. Record or label the r5 pre-review and post-r6 assessment as review inputs (R7-5).

Until then, r7 is **PROPOSED — NOT BINDING**, the original requirements remain in force, r5 and r6 remain preserved and non-adopted, A-1 remains open with PR-1–PR-4 inactive, and the eight Stage 3 decisions remain untouched.

## Limitations

- **L1 — Review scope.** This is a full independent review of the text as written, not an implementation, runtime, or measurement review. No build, test, replay, benchmark, or stage execution occurred.
- **L2 — Upstream provenance.** F1–F10, P1–P3, and S1 exist only as r7's (and r6's) summaries. I verified every contract and plan citation they rely on, but not that the findings themselves were recorded as stated by their originators (R7-5 records the gap for P1–P3/S1).
- **L3 — Design status.** The corrective sequencing document is a proposal, not an adopted artifact. Where r7 and the plan differ, I record the difference; I do not resolve the plan's status.
- **L4 — Reviewer position.** I authored the r2, r3, and r4 audits this lineage responds to. This review is independent of the r7 author but not of the review program.
- **L5 — Hash scope.** I verified r7 and all thirteen §2.2 inputs. I did not verify r7's claims about files outside `docs/review/`.

## Review record

- Read in full: r7 (715 lines).
- Mechanically compared: §7.3 against original AC-1–AC-7; all 13 §2.2 rows against recomputed hashes and `wc`; all "next explicit" occurrences; all `A2-D12` occurrences; the 27-row, six-column §7.1 matrix.
- Contract passages verified: §17.7.2 transition table (L1773, L1793–L1799, L1801–1802), §17.11.5 (L2105), §17.12.1 adapter rows, the two-process cap (L2366), A2-D7 (L2508), A2-D12 (L2513), line 1844, and the zero-`C8` scan.
- Plan passages verified: L348, L354–355, L362, L374, L512.
- Files changed by this review: this file only. No subject was modified.
