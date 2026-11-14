# AUDIT — Howl thread/hardware/performance requirements, revision r9

**Reviewer:** Mimo (independent of the r9 author)
**Date:** 2026-11-08
**Mode:** Read-only independent review. No subject file was edited, repaired, reconciled, staged, committed, or executed. No build, test, replay, measurement, benchmark, or stage execution was run. No adoption, implementation, or runtime authorization is implied.
**Subject:** `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r9-2026-11-08.md`
**Subject SHA-256:** `0f71cee5e0e5d25b58d1b3956ff549df8cc625ad7f9e8c82189f2f2ef9366c55`
**Subject size:** 769 lines / 112,270 bytes, CRLF (769 CR bytes, one per line)
**Predecessor reviewed:** r8 (`b14dc468…`, 724 / 102,800), audited in `AUDIT-howl-requirements-thread-hardware-safety-r8-2026-11-07.md` (`71ef9b72…`)

## Verdict

**Adoptable as written.** No blocking or major defect. Three minor record-preparation defects (M1–M3) and thirteen observations. None of the minor items is an unresolved authority, contract, lifecycle, or safety contradiction of the kind §14's second item excludes; M1 and M2 can be cured in the adoption record itself, and M3 is a label to correct in the next revision.

The four r8 corrections are correctly applied. R8-1: §14 no longer conditions adoption on any decision remaining pending, and D9 is counted once in each list. R8-2: §7 carries the frozen `R8-AC-1`–`R8-AC-10` overlay set and contains no `R7-AC` prefix. R8-3: all fifteen §2.2 statements now have a §2.3 field, and thirteen of the fifteen mappings are fully substantive. R8-4: §8.1 adds a thirteen-state matrix whose live/exited `unavailable` split matches §17.7.4 exactly.

The review's four focus areas:

- **§2.2 → §2.3 mapping.** 15/15 statements are mapped, and thirteen of the fifteen mappings are fully substantive. M1 and M2 concern two fields that permit a form §2.2 does not accept; neither leaves a statement without a field.
- **§14.** Twenty-three items. No item requires A-1 or any of the eight Stage 3 decisions to remain pending. L764 states the opposite and is consistent with §13 L732–734. D9 appears once per list.
- **§8.1.** All thirteen contract states are present. The live versus exited `unavailable` distinction is correct and is the contract's own distinction (§17.7.4). Wire-health rows are grounded in §17.4.1 (handshake complete) and §17.4.5 (read-only; no load, runtime initialization, new-asset discovery, synthesis, or failure reset; accepted during `load` and `synthesize`). The supervisor-level rows are r9-internal behavioral assertions, not contract health operations; the required evidence is adequate for them.
- **Authority boundaries.** The adoption template cannot authorize a stage, implementation, runtime verification, benchmark, or performance gate. §2.2 L158, §2.3 L206, §2.4, §12, and §7.2 each disclaim that authority independently.

## Severity scale

- **Blocking** — would freeze a wrong interpretation into stage acceptance, or let a stage pass while a contract-settled behavior is unimplemented or unverified.
- **Major** — authority, scope, or allocation defect that does not by itself produce a wrong pass.
- **Minor** — citation, ambiguity, evidence-precision, or drafting defect.
- **Observation** — not a defect; recorded because it affected the review or confirms a pass.

---

## Minor

### M1 — §2.3 permits an attachment path where §2.2 requires the record to contain the exact table

**Evidence.** §2.2 L148 requires the adoption record to contain "the exact governing-input table below, with hashes rechecked at adoption." §2.3 L194 provides:

> "**Governing-input table (§2.2) exact copy or attachment path:**"

**Why it matters.** A path is a reference; it is not the table. A record that fills this field with an attachment path satisfies the field's text but not §2.2's "must contain" requirement, so the record could be incomplete while looking complete. The "exact copy" option is compliant, so this is a breadth defect in the template, not a missing field.

**Fix.** Remove "or attachment path," or state that an attachment is acceptable only if it is incorporated into the record and reproduces the table in full with its hashes.

### M2 — §2.3's contract-status field mixes the mandatory match confirmation with the optional §17 instrument

**Evidence.** §2.2 L150 requires "confirmation that the governing contract and the correction record match the table and were not amended." §2.2 L152 separately requires the §17 contract-amendment instrument only "if the owner elects a contract-level C8 performance requirement." §2.3 L198 places both in one field as two checkboxes:

> "**Governing contract status:** [ ] Matches §2.2; unamended; … [ ] Separate §17 contract-amendment instrument outside A-1 Option B — instrument/hash:"

§2.3 L200 already provides the §17 instrument slot under "Contract-level C8 performance requirement (if elected)."

**Why it matters.** The field presents a mandatory confirmation and an optional instrument as alternatives. A record electing the §17 path could check only the second box and leave the L150 confirmation unstated, or record the same instrument twice. The defect is record completeness, not authority: it cannot authorize a stage, and the §17 path is outside what r9 authorizes.

**Fix.** Keep L198 as the mandatory "matches §2.2; unamended" confirmation for the contract (and L199 for the correction record), and keep the §17 instrument only in L200.

### M3 — r6 is labeled "Pre-review assessed" although it was assessed after r6

**Evidence.** The §2.2 table at L172 reads:

> "| Pre-review assessed r6 predecessor | …r6… | `8398802b…` | 700 | 95,092 |"

The provenance note at L99 distinguishes "the 2026-11-07 pre-review of r5" from "the post-r6 restart assessment." r5 was pre-review assessed (L171); r6 was the subject of the post-r6 restart assessment (S1), not a pre-review.

**Why it matters.** The artifact, hash, line count, and byte count are correct, so the input is identified and the error does not affect §14 item 4. The label misstates the review relationship and is inherited from r8 L159; it is not a regression by this author.

**Fix.** Relabel L172 "Post-r6 assessed r6 predecessor" or "Assessed after r6."

---

## Observations

**OB-1 — Artifact and all seventeen input rows verify.** r9 is 769 lines / 112,270 bytes with SHA-256 `0f71cee5e0e5d25b58d1b3956ff549df8cc625ad7f9e8c82189f2f2ef9366c55`. It is CRLF (769 CR bytes), unlike r8, which was LF. The hash binds the CRLF form; the artifact should be preserved byte-for-byte. On normalized text the delta from r8 is +45 lines and +8,701 bytes. I recomputed every §2.2 row — original, r2, r2 audit, r3, r3 audit, r4, r4 audit, r5, r6, r7, r7 audit, r8, r8 audit, contract, plan, correction record, and owner-decision instrument — and every hash, line count, and byte count matches. The table adds r8 and the r8 audit to r8's fifteen rows.

**OB-2 — §7.3 and the matrix verify.** §7.3 reproduces all seven original criteria verbatim against the original AC table (mechanical diff: identical). The §7.1 matrix has 27 data rows and six columns in every row, with no duplicate identifiers: seven original AC rows, ten overlay rows, and ten stage rows.

**OB-3 — R8-1 is fixed.** The §14 checklist has 23 items. No item requires any decision to remain pending; L764 states that adoption "neither requires nor implies that any decision remains pending" and grants no implementation authority. §13 L732 lists the eight as D10(a), D10(b), D1a, D9, D1b, D2, D5, and D3, and L733–734 makes the list a status report, not a condition. D9 appears once in each list: L114 (singular, as one of the eight), L260 (§3.3 non-resolution), L585 (§8 gate list), L732 (§13 status), and L764 (the eight, without D9). L41 and L123 are historical finding text. There is no "D9, A-1, and all eight" construction anywhere.

**OB-4 — R8-2 is fixed.** §7 now says the revision "carries the frozen overlay identifiers **R8-AC-1–R8-AC-10** and adds no new overlay identifier." There is no `R7-AC` identifier and no overlay "R7-" prefix. The matrix's "Originating authority" column says "r9 §7" for the frozen set, which is where r9 states the criteria; the identifier freeze is explained once and used consistently.

**OB-5 — R8-3 mapping.** All fifteen §2.2 statements have a mapped §2.3 field:

| §2.2 statement | §2.3 field | Result |
|---|---|---|
| Owner identity/accountable role + stage acceptance authority | L188–189 | Pass |
| Adoption date | L190 | Pass |
| Exact path and SHA-256 | L191–192 | Pass |
| Adopted scope incl. C1–C8 and eventual implementation | L193 | Pass |
| Exact governing-input table, hashes rechecked | L194–195 | **M1** (attachment option) |
| Design incorporation; if not, review-inputs statement | L196 | Pass |
| Contract and correction record match table, unamended | L198–199 | **M2** (mixed field) |
| Plan status; Option B instrument, scope, SHA-256, incorporation Yes | L197 | Pass |
| Contract-level C8 §17 instrument outside Option B | L200 | Pass |
| Independent r9 review record and verified hash | L201 | Pass |
| §7.1 matrix checked, no collision/leakage, four responsibilities separated | L202 | Pass |
| §10 role bounded, no delegation, not stage acceptance | L203 | Pass |
| A-1 decision; if open, PR-1–PR-4 inactive | L204 | Pass |
| Approved exceptions, scope, duration, safety evidence | L205 | Pass (see OB-11) |
| Explicit statement that adoption authorizes no stage or implementation | L206 | Pass |

**OB-6 — R8-4 is fixed.** §8.1 has thirteen rows, exactly the thirteen aggregate states in §17.7.1: `absent`, `starting`, `handshaking`, `unloaded`, `loading`, `ready`, `busy`, `degraded`, `stopping`, `crashed`, `incompatible`, `unavailable`, `backoff`. The wire rows are grounded in §17.4.1's "Handshake complete" precondition and §17.4.5's read-only rule; the supervisor rows cover the states where no handshaken helper exists. The `handshaking` and `incompatible` exclusions are grounded in contract preconditions (§17.4.1 and A2-D4), not test cost.

**OB-7 — The live/exited `unavailable` split is correct.** §8.1's two variants match §17.7.4's two additive rows: `ProcessState=alive`, `State=unavailable`, `Available=true`, `ServiceCapabilitiesKnown=true` (asset-specific failure on a live helper) and `ProcessState=exited`, `State=unavailable` (restart budget exhausted). A2-D12 L2513 supplies the recovery rule: health, status, catalog, voice-list, and synthesis do not trigger recovery; only an explicit same-identity `load` on a live helper whose helper-side state is `unloaded` does. The table's live variant correctly retains the wire-health surface because the handshake is retained.

**OB-8 — Supervisor-level rows are behavioral, not contract health operations.** For `absent`, `starting`, `handshaking`, `crashed`, `incompatible`, `backoff`, and the exited `unavailable` variant there is no handshaken helper, so there is no wire `health` response to which §17.7.2's "never a transition source" sentence directly applies. The rows assert an r9-internal property — an application-level health read must not start, restart, recover, load, or transition — and the required evidence (before/after provider and process state, no start/frame/transition attributable to the read) is adequate. The matrix would be more precise if it said this expressly, so a reader does not treat the contract citation as establishing the supervisor rows.

**OB-9 — Live `unavailable` should record the helper-side state.** A wire `health` response from the live variant reports the helper-side state (`unloaded`), while the supervisor state remains `unavailable`. The row's "no recovery, load, start, or transition" assertion covers this, but the record requirement should say to capture both the supervisor `State=unavailable` and the helper-side `unloaded`, so a reviewer does not read `unloaded` in the response as recovery and does not miss a helper-side change.

**OB-10 — "Only exclusions" is narrower than the table.** §8.1 L618 says `handshaking` and `incompatible` are "the only exclusions from the wire-health matrix." Five further states (`absent`, `starting`, `crashed`, `backoff`, exited `unavailable`) also lack a wire request; the next sentence explains that they use the supervisor surface because no handshaken helper exists. The table is correct; the word "exclusions" is doing two different jobs. Wording, not substance.

**OB-11 — Authority boundaries pass, with one timing observation.** §2.2 L158 requires an explicit non-authorization statement; §2.3 L206 supplies it and adds runtime execution, benchmarking, and deployment. §2.4 says no adoption, review, or stage acceptance authorizes C1–C8, implementation, runtime execution, or benchmarking. §12 prohibits review-created and implied authority. §7.2 says acceptance does not authorize implementation. The one timing observation is the "Approved exceptions" field at L205: §12 permits an HW-1 budget exception only through an owner decision "recorded after adoption." A record prepared at adoption should mark the field "None" and treat later exceptions as post-adoption entries; the field should say so.

**OB-12 — TTS-only and A2-D12 regressions pass.** Every "next explicit" occurrence is TTS-qualified or explicitly historical: L21, L44, L77, L108, L109, L589, L593, L644, and L752 are operative and TTS-qualified; L97 is the historical P3 withdrawal. L108 also labels §17.7.2 L1796's "next explicit request" as shorthand that does not broaden the class. All A2-D12 mentions (L21, L44, L62, L109, L110, L128, L550, L593, L763) keep the rule confined to `unavailable`; none generalizes it to `crashed` or `backoff`. No "under proposed plan C5" qualifier survives.

**OB-13 — Status boundary is preserved.** No adoption record exists. r9 is proposed and non-binding. A-1 remains open; PR-1–PR-4 are inactive; C1–C8 are unauthorized; the eight Stage 3 decisions are pending; the original requirements remain operative. This audit creates no adoption, no stage acceptance, and no implementation authority.

---

## Target-by-target result

| Target | Result |
|---|---|
| **§2.2 → §2.3, all 15 statements** | **Pass with M1–M2.** Every statement has a substantive field; two fields permit a form §2.2 does not accept. |
| **§14 pending-decision independence and D9 count** | **Fixed.** 23 items; no pending condition; D9 once per list; §13 and §14 consistent. |
| **§8.1, 13 state rows** | **Fixed.** All 13 states; live/exited `unavailable` matches §17.7.4; wire rows contract-grounded; supervisor rows behavioral and adequately evidenced. |
| **Authority boundaries** | **Pass.** No stage, implementation, runtime, benchmark, or performance-gate authorization; OB-11 timing observation only. |
| **R8-1 — §14 contradiction** | **Fixed.** |
| **R8-2 — overlay-prefix residue** | **Fixed.** |
| **R8-3 — adoption-template omissions** | **Fixed**, with M1–M2. |
| **R8-4 — health evidence coverage** | **Fixed**, with OB-8–OB-10. |
| **Regression — TTS-only restart** | **Pass.** |
| **Regression — A2-D12 `unavailable`-only** | **Pass.** |
| **Regression — original AC-1–AC-7** | **Pass.** §7.3 verbatim. |
| **Regression — matrix** | **Pass.** 27 data rows, six columns, no duplicates. |

## Adoption readiness

- **Checklist item 1** — independent review complete and recorded with a verified hash: satisfied by this audit once its hash is recorded.
- **Checklist item 2** — no unresolved authority, contract, lifecycle, or safety contradiction: **satisfied.** M1–M3 are record-preparation precision and a label; none is a contradiction of the specified kinds.
- **Checklist item 4** — governing-input hashes rechecked and match: 17/17 verified above.
- **§2.3 record** — still blank, as it should be. No adoption record was created. A-1 remains open, PR-1–PR-4 inactive, C1–C8 unauthorized, and the eight Stage 3 decisions pending.
- **M1/M2 cure** — both can be cured in the adoption record: use the exact-copy option at L194, and check the L198 match confirmation while recording any §17 instrument at L200. No new revision is required for the record to be complete.
- **M3** — relabel the r6 row in the next revision; the current hash identifies the correct artifact.

r9 is **adoptable as written** for the owner's separate adoption decision, subject to the record preparer avoiding the M1/M2 field forms. This audit does not adopt it.

## Limitations

- **L1 — Review scope.** Text and hash review only. No build, test, replay, benchmark, or stage execution; §8.1's evidence matrix was assessed for contract grounding and completeness, not executed.
- **L2 — Upstream provenance.** F1–F10, P1–P3, and S1 remain r9's summaries. I verified the contract citations that the health, restart, and §8.1 rules rely on, and the plan lines r9 cites, not every citation in the finding summaries and not the original recording of each finding.
- **L3 — Reviewer position.** I authored the r2, r3, r4, r7, and r8 audits in this lineage. This review is independent of the r9 author but not of the review program; the r8 audit it relies on is my own work.
- **L4 — Hash scope.** I verified r9 and all seventeen §2.2 rows. I did not verify claims about files outside `docs/review/`.
- **L5 — Design status.** The stage plan remains a proposal. Where r9 and the plan differ, I record the difference without resolving the plan's status.

## Review record

- Read in full: r9 (769 lines), including §2.2, §2.3, §7, §8, §8.1, §13, and §14.
- Mechanically compared: all 17 §2.2 rows against recomputed SHA-256, line, and byte counts; §7.3 against the original AC table; the 27-row, six-column §7.1 matrix and its identifiers; the 13 §8.1 state rows against §17.7.1 and §17.7.4; the 15 §2.2 statements against the 20 §2.3 fields; all "next explicit" and `A2-D12` occurrences; the §14 items against §13, §1.1 A-2, and §1.5.
- Contract passages re-verified: §17.4.1 (L1328), §17.4.5 (L1388–1391), §17.7.1 (L1743, L1754–1755), §17.7.2 (L1769–1802), §17.7.4 (L1822–1830), §17.11.5 (L2105–2108), §5.3 health row (L490), A2-D4 (L2505), A2-D12 (L2513), O-1 (L2271, L2348).
- Plan passages re-verified: L348, L354, L512.
- Files changed by this review: this file only. No subject was modified.
