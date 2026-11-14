# AUDIT — Howl thread/hardware/performance requirements, revision r8

**Reviewer:** Mimo (independent of the r8 author)
**Date:** 2026-11-07
**Mode:** Read-only independent review. No subject file was edited, repaired, reconciled, staged, committed, or executed. No build, test, replay, measurement, benchmark, or stage execution was run. No adoption, implementation, or runtime authorization is implied.
**Subject:** `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r8-2026-11-07.md`
**Subject SHA-256:** `b14dc4684605bed61eae83283578a37c34c979d88a362cdd99ff07821c24b4cf`
**Subject size:** 724 lines / 102,800 bytes, LF-only (0 CR bytes)
**Predecessor reviewed:** r7 (`58cd4656…`, 715 / 97,790), audited in `AUDIT-howl-requirements-thread-hardware-safety-r7-2026-11-07.md` (`8a1956bc…`)

## Verdict

**Not adoptable as written.** One major defect, three minor, eleven observations. No blocking defect.

The five r7 corrections are correctly applied. The health/status split is exactly right: health is state-invariant, status and catalog are barred from the crash-restart condition without asserting general state invariance, and the plan's only independent addition is the no-background-loop rule. C5's health and status/catalog evidence is now unconditional. The `stopping` map is complete. The design-incorporation condition points to §2.2 everywhere. The two non-reproducible inputs are labelled as such.

The three regression areas also hold: the TTS-only restart rule is qualified at every operative locus, A2-D12 is confined to `unavailable` at all eight mentions, and §7.3 reproduces AC-1–AC-7 verbatim. The matrix has 27 rows and six columns with no identifier collision.

The major defect is in the adoption checklist, not the engineering requirements: §14 still requires that all eight Stage 3 decisions "remain separate and unresolved" as a condition of adoption (L719), which contradicts §13 (L687), the preceding §14 item (L717), and the document's own A-2 disposition (L41). It is inherited from r7 L710 — where my r7 audit recorded it only as a double-count — so it is not a new regression by this author, but it is unresolved in r8's gate and must be fixed before adoption.

The fix is one line. The engineering corrections need no rework.

## Severity scale

- **Blocking** — would freeze a wrong interpretation into stage acceptance, or let a stage pass while a contract-settled behavior is unimplemented or unverified.
- **Major** — authority, scope, or allocation defect that does not by itself produce a wrong pass.
- **Minor** — citation, ambiguity, evidence-precision, or drafting defect.
- **Observation** — not a defect; recorded because it affected the review or confirms a pass.

---

## Major

### R8-1 — The §14 checklist makes "all eight decisions remain unresolved" an adoption condition, contradicting §13 and the document's own A-2 disposition

**Evidence.** §14 opens "Before owner adoption, all of the following must be true" and closes "Until every item is complete, this revision is **PROPOSED — NOT BINDING**." One item reads:

> "- [ ] D9, A-1, and all eight Stage 3 owner decisions remain separate and unresolved." (r8 L719)

Three other passages say the opposite:

> "This is a status report, not an adoption condition. **Adoption neither requires any decision to remain pending nor treats an answered decision as authority for this revision.**" (§13, r8 L687)

> "- [ ] This revision resolves no owner decision; **the adoption record does not require or imply that any decision remains pending**, and does not cite an answered decision as authority." (r8 L717)

> "§13 and §14 no longer make “all eight remain pending” an adoption condition." (§1.1 A-2, r8 L41)

**Why it matters.** §13 explicitly contemplates the owner answering a Stage 3 decision before adoption: "If any Stage 3 decision is answered before adoption, r8 does not incorporate that answer as authority" (r8 L688). On that path, L719 cannot be truthfully checked — the decision no longer "remains unresolved" — so the checklist's all-items rule blocks adoption, or the owner certifies a false statement. Either outcome defeats A-2, which r8 claims to have fixed. A checklist that is unsatisfiable on a path its own §13 permits is a defect in the adoption authority, not a drafting nicety.

**The double count is separate and smaller.** §13 lists the eight as D10(a), D10(b), D1a, D9, D1b, D2, D5, and D3 (r8 L686). D9 is therefore one of the eight, so "D9, A-1, and all eight" describes ten items. The same construction appears in the header (L9) and §1.5 (L114). This alone would be drafting; combined with the contradiction, it shows the item was carried forward unexamined.

**Inheritance and reviewer error.** r7 L710 contains the identical item, and r7 L680 contains the identical §13 sentence. My r7 audit recorded this as OB-8, "double-counts D9", and did not identify it as a contradiction with §13. That was an under-classification on my part; the item was not in R7-1–R7-5, so the r8 author was not asked to change it. It is nevertheless unresolved in r8 and is the checklist item that fails outright.

**Severity.** Major, not blocking. It freezes a wrong adoption condition and can prevent a legitimate adoption, but it does not let a stage pass with unsafe or unimplemented behavior.

**Fix.** Replace L719 with a statement that the revision records no decision as resolved and cites none as authority — the substance already in L717 (merging the two items is sufficient) — and remove the D9 double count at L9, L114, and L719.

---

## Minor

### R8-2 — §7 still says the overlay prefix is "R7-"

> "This revision adds overlay criteria with an explicit **R7-** prefix. A bare `AC-n` always means the original criterion; an `R8-AC-n` is an r8 overlay criterion." (r8 L478)

Every overlay is actually `R8-AC-1`–`R8-AC-10`, and the matrix and checklist use `R8-AC-`. This is a rename residue: the sentence was correct when r7's overlays were `R7-AC`, and it was not updated when r8 renamed them. A reader checking the stated prefix finds no `R7-AC` identifier and may conclude the overlay set is missing. Minor; fix the prefix to `R8-AC-`.

### R8-3 — The §2.3 adoption block omits fields for five statements §2.2 requires the record to contain

§2.2 says the adoption record "must contain" the following, but the §2.3 template provides no field for them:

| §2.2 requirement | §2.3 status |
|---|---|
| Confirmation that the correction record matches the table and was not amended | **Missing** — only contract status is present |
| Statement that the four responsibilities are separated, not just that the matrix was checked | **Partial** — a yes/no checkbox only |
| Statement that the §10 budget-conformance role is bounded, delegates no owner authority, and is not stage acceptance | **Missing** |
| Statement that PR-1–PR-4 are inactive when A-1 is unresolved | **Missing** — the A-1 field assumes a gate decision |
| Explicit statement that adoption does not authorize any stage or implementation | **Missing** |

The substantive rules exist in §2.4, §7.2, and §10, so this is not authority leakage. But a record built from the template would fail §2.2's own "must contain" list while still satisfying the checklist item that only checks hash and scope. Minor; add the five fields (or one attestation block covering them).

### R8-4 — The health invariant is global, but its executed evidence is not state- or surface-scoped

The requirement is correctly stated as global: health "leaves state unchanged across the state machine" and "cannot start, restart, recover, or otherwise transition a helper" (§1.5, r8 L108), and C1 repeats "A health observation leaves state unchanged and is never a transition source" (r8 L567). C5's evidence item, however, is:

> "health observation state-unchanged evidence, including no transition to `degraded`" (r8 L571)

and C7's is a generic "no recovery from `unavailable`" (r8 L573). Neither names the states or API surfaces the invariant must span. A stage can satisfy C5 with one probe in a live state and never execute a health request from `absent`, `stopping`, or `incompatible`, where a health-triggered start or transition would equally violate the contract. The `degraded` clause forces a live-state probe; it does not force the start and recovery cases.

This is the same family as R7-1 — the requirement is right, the evidence is narrower than the requirement — but it does not reopen the attribution. Minor; fix by enumerating the states (at minimum `absent`, `ready`/`busy`, `crashed`, `stopping`, `incompatible`, `unavailable`) or requiring a request × state table with the observed result.

---

## Observations

**OB-1 — Artifact and all fifteen input rows verify.** r8 is 724 lines / 102,800 bytes with SHA-256 `b14dc4684605bed61eae83283578a37c34c979d88a362cdd99ff07821c24b4cf` and zero CR bytes. I recomputed every §2.2 row: original, r2, r2 audit, r3, r3 audit, r4, r4 audit, r5, r6, r7, r7 audit, contract, plan, correction record, and owner-decision instrument. Every hash, line count, and byte count matches. The delta from r7 is +9 lines and +5,010 bytes — consistent with a narrow correction, not a rewrite.

**OB-2 — §7.3 is verbatim against the original.** I extracted the seven text cells from r8's §7.3 and from the original AC table (original L248–254) and diffed them: identical. AC-2's static-review sentence and AC-3's repeatability sentence are intact. R4-5 remains closed.

**OB-3 — The TTS-only restart rule is intact at every operative locus.** All "next explicit" occurrences are TTS-qualified or explicitly historical: F2/A-5/R4-1 (L21, L44, L77), the withdrawn r6 P3 (L97), the §1.5 contract row (L108), the plan row (L109), C1 (L567), C5 (L571), §9.2 (L598), and the checklist (L706). L1796 is labelled shorthand that "must be read consistently with §17.11.5" and "does not broaden the eligible request class." No regression.

**OB-4 — A2-D12 remains `unavailable`-only at all eight mentions.** L21, L44, L62, L109, L110, L528, L571, and L718 all confine the explicit-load exception to `unavailable`; none generalizes it to `crashed` or `backoff`. The ordinary `unloaded` load path is preserved, and C1's note that a `load` transition in an already-live helper "is not a process start" remains.

**OB-5 — Matrix and adoption paths verify.** The §7.1 matrix has 27 data rows, six columns in every row, and no duplicate identifiers: seven original AC rows, ten R8-AC rows, and ten stage rows. §2.5's Option A, Option B, and open paths remain internally consistent, and Option B still requires a plan-amendment instrument, scope, and hash without amending the contract.

**OB-6 — R7-1 is correctly fixed.** The health bar is now attributed to the contract, not the plan: §1.5 cites §17.7.2 L1801–1802 and states health cannot start, restart, recover, or transition; C1 and C5 repeat the rule; §9.2 and the checklist cite it. Status and catalog are deliberately narrower — barred from satisfying §17.11.5, with "other applicable contract transition rules remain in force" — and C5 explicitly disclaims general state invariance for them. The plan's independent addition is now only the no-background-loop rule. C5's health and status/catalog evidence items are unconditional; no "under proposed plan C5" qualifier survives. The r7 audit's blocking finding is closed.

**OB-7 — R7-2 is correctly fixed.** C1's `stopping` row now records both edges: exit 0 → `absent`; nonzero exit or grace expiry → `crashed`, which then follows the crash-restart rule. The F2, A-5, and R3-5 dispositions and the checklist repeat it. This matches contract L1793–1795 and A2-D8.

**OB-8 — R7-3 is resolved, with one residual observation.** HW-1.8 now requires instrumentation overhead to be measured "separately from contention-induced degradation when that correctness instrumentation is executed" and classifies it as "not a performance-gate measurement, benchmark, or PR-suite execution," creating "no performance claim" (r8 L408). §6.5 repeats the permission and the no-PR-result rule (r8 L475), and §6.4, §12, and the checklist repeat the firewall. A performance *gate* measurement still requires an owner-authorized threshold, contention profile, run count, and warm-up period, none of which exists while A-1 is open; the measurement therefore cannot gate or accept anything. Residual: the method — measuring overhead separately from contention — is structurally close to PR-1/PR-2, and the exclusion is established by purpose and classification rather than by method. That does not create an authorization path, but if the owner wants the boundary tighter, the sentence could limit the measurement to the instrumentation delta and forbid reporting it as a latency/throughput result or reusing it as PR evidence.

**OB-9 — R7-4 and R7-5 are correctly fixed.** Every plan-incorporation reference now points to §2.2: §1.5, §2.1, §2.2, R3-3, the ten stage rows in §7.1, the §8 preamble, and the checklist. The r5 pre-review and post-r6 assessment are labelled "non-reproducible drafting-time inputs" in the header, §1.4, §2.2, and the checklist, with S1's contract citations identified as the operative evidence.

**OB-10 — The health claim matches the contract on independent re-check.** The §17.7.2 event column contains no health event; `degraded` is reached only by `busy | synthesis fails`; the supervisor health surface is read-only and leaves state unchanged (§17.x line 490: "No new I/O … `degraded`/`last_error`"). r8's state-invariance claim is therefore accurate, and no `degraded` carve-out is needed.

**OB-11 — No adoption record exists and none should be created on this review.** A-1 remains open, PR-1–PR-4 are inactive, C1–C8 remain unauthorized, and the eight Stage 3 decisions are pending. The §2.2 "A hash mismatch in any row reopens review" rule is satisfied: none mismatched.

---

## Target-by-target result

| Target | Result |
|---|---|
| **R7-1 — observation bar attribution** | **Fixed.** Health is contract-level (§17.7.2 L1801–1802); status/catalog are barred from the §17.11.5 condition without overextension; plan addition is only the no-background-loop rule; C5 evidence unconditional. |
| **R7-2 — `stopping` map** | **Fixed.** Exit 0 → `absent`; nonzero/grace expiry → `crashed`; matches the contract. |
| **R7-3 — HW-1.8 vs §6.5** | **Fixed.** Correctness instrumentation measurement is permitted and firewalled from performance-gate effects; no acceptance path while A-1 is open. Residual method overlap noted in OB-8. |
| **R7-4 — incorporation cross-references** | **Fixed.** All point to §2.2. |
| **R7-5 — unhash-bound inputs** | **Fixed.** Labelled non-reproducible drafting-time inputs in four places. |
| **Regression — TTS-only restart** | **Pass.** Qualified at every operative locus; L1796 labelled shorthand. |
| **Regression — A2-D12** | **Pass.** Eight mentions, all `unavailable`-only. |
| **Regression — original AC-1–AC-7** | **Pass.** §7.3 text identical to the original (mechanical diff). |
| **Regression — matrix and adoption template** | **Pass on structure** (27 rows, six columns, no duplicates); §2.3 template gaps are R8-3. |

## Adoption readiness

- **Checklist item 1** — independent review complete and recorded with a verified hash: satisfied by this audit once its hash is recorded.
- **Checklist item 2** — "no unresolved authority, contract, lifecycle, or safety contradiction": **fails** because of R8-1.
- **Hashes, matrix, §7.2, §7.3, C1 map, C5 evidence, C3, design incorporation, budget, §2.5, O-1, HW-1.8, C8, contract/plan status, provenance:** verified above; no other failure.
- **§2.3 record template:** would be incomplete under §2.2 until R8-3 is fixed.

R8 is therefore **PROPOSED — NOT BINDING** and not adoptable as written. The engineering content needs no redesign: fix the §14 item, the prefix residue, the template fields, and the health evidence scoping, then re-review those four loci and the unchanged regression areas.

## Limitations

- **L1 — Review scope.** Text review only. No build, test, replay, benchmark, or stage execution; the HW-1.8/PR framework was assessed for authority and wording, not measured.
- **L2 — Upstream provenance.** F1–F10, P1–P3, and S1 remain r8's summaries; I verified the contract and plan citations they rely on, not the original recording of each finding.
- **L3 — Reviewer position.** I authored the r2, r3, and r4 audits and the r7 audit this revision responds to. This review is independent of the r8 author but not of the review program, and one of its findings corrects my own earlier under-classification.
- **L4 — Hash scope.** I verified r8 and all fifteen §2.2 rows. I did not verify claims about files outside `docs/review/`.
- **L5 — Design status.** The stage plan remains a proposal; where r8 and the plan differ, I record the difference without resolving the plan's status.

## Review record

- Read in full: r8 (724 lines).
- Mechanically compared: all 15 §2.2 rows against recomputed SHA-256, line, and byte counts; §7.3 against the original AC table; the 27-row, six-column §7.1 matrix and its identifiers; all "next explicit" occurrences; all `A2-D12` occurrences; all plan-incorporation references; the §14 checklist items against §13 and §1.1 A-2.
- Contract passages re-verified: §17.7.2 transition table (L1773, L1793–1799, L1801–1802), the health surface row (L490), §17.11.5 (L2105), line 1844, A2-D7 (L2508), A2-D12 (L2513), O-1 (L2271, L2348).
- Plan passages re-verified: L348, L354–355, L362, L374, L512.
- Files changed by this review: this file only. No subject was modified.
