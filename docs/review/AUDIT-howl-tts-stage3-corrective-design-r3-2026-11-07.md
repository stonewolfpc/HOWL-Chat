# AUDIT — Stage 3 corrective design, round 3

**Date:** 2026-11-07
**Mode:** DOCUMENTATION REVIEW ONLY — read-only
**Subject:** corrected Stage 3 corrective sequencing proposal `1fb449ba…` and its second correction record `ccb28d57…`
**Governing contract:** `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md`
**Authorizes:** nothing. No stage, no code, no runtime.

This audit verifies whether the second correction pass documented the round-2 findings
(`a7384da9…`) faithfully, whether any earlier defect was reintroduced, and whether the stage
and decision structure remains internally consistent. It is not a review of the round-2 audit's
own conclusions except where those conclusions are tested against the contract again.

## 1. Verified inputs

Recomputed at the start of this review:

| Artifact | SHA-256 | Size | mtime |
|---|---|---|---|
| Governing contract | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f` | 2,586 lines / 182,904 B | 2026-11-06 20:57:07 |
| Corrective proposal (subject) | `1fb449ba2baef7b12e42a7b712d2370ac27e6f4dddb9bcb48d6aac3cad87630d` | 710 lines / 65,213 B | 2026-11-07 10:15:04 |
| Correction record (subject) | `ccb28d57f14f46988a0277df40676b9e67d5223b0a30ae33600ec75cff8c67d4` | 175 lines / 13,087 B | 2026-11-07 10:16:14 |
| Round-2 audit (evidence) | `a7384da9f3f38661fc28b3de48badb82245cfd41efb8e49964f17bcfaaf5a7e2` | 462 lines / 37,477 B | 2026-11-07 10:08:26 |
| Round-1 audit (evidence) | `71a17369e2389fae26e86783fc4df802ab91bad3daef512c0e7f743067e94ea2` | 249 lines / 25,962 B | 2026-11-07 09:31 |
| Readiness investigation | `a63f9f32f4d1bd047746a7bbbf5f509588e13ebd4272cffeb23f61fcfa891c6c` | 586 lines / 51,248 B | 2026-11-07 08:58 |

The record's §1 claims match disk for the pre-pass hash (`0f920b24…`, 703 lines / 59,634 B), the
post-pass hash, both prior audits, the contract, and the investigation. The proposal and record
were modified after the round-2 audit was written, so the round-2 report describes the superseded
`0f920b24…` text; its line references are stale and were re-derived here from the 710-line file.

## 2. Disposition

**PASS — the correction pass is accepted as documentation of the round-2 findings.**

All six addressed items are genuinely documented in the proposal text, not merely asserted in the
record. No earlier correction was narrowed or lost. One new **major** finding (M5) is a
misclassification in the decision ledger, not a defect in the correction; it must be repaired
before D9 is put to the owner, because a plausible owner answer would make C5 contradict §17.7.2.
Four **minor** findings are wording or coverage gaps with no acceptance consequence at C1.

Consequences for authorization, stated without authorizing anything:

- **C1 is free of blocking defects.** It may be authorized for implementation, but it cannot be
  accepted before D10 is dispositioned, because C1's acceptance explicitly refuses to freeze an
  unknown-snapshot result.
- **C2 is clear.** M1 is closed; the projection set is complete against A2-D5.
- **C3 remains gated on D1a**, which is now scheduled correctly.
- **C5 cannot be accepted** until D9 is dispositioned, and D9's framing must be corrected first (M5).

## 3. Finding-by-finding disposition

### B1 — cancellation scope — **Closed**

The proposal now separates the two call sites structurally:

> "A direct optional operation (`unload` or `Cancel`) whose capability is known absent returns
> `unsupported_capability` with no frame (A2-D1, A2-D7)." — proposal line 121

> "The path must use a path-scoped internal operation so the exemption cannot leak into the
> exported direct method." — proposal line 122

The regression test is now two-sided rather than blanket:

> "a `cancel`-less helper on the supervisor path produces no cancel frame and no
> `unsupported_capability` result, then exits through the existing grace/kill path; a direct
> `Cancel` on the same helper returns `unsupported_capability` with no frame." — proposal line 162

Acceptance, the §7 inventory row, and two separate §8 matrix rows carry the same split. This
matches A2-D1 ("The adapter returns `unsupported_capability` for an absent optional service
capability and MUST NOT write an operation frame", line 2502) and A2-D9's narrower grant
("absence of `cancel` is not an error in the supervisor's cancellation path and no
`unsupported_capability` is synthesized there", line 2510). The round-2 objection — a method-scoped
exemption that would let a direct call succeed while sending nothing — no longer applies.

**Evidence grade:** Observed (both texts read directly).

### B2 — unknown-snapshot unload — **Closed**

`provider_not_loaded` no longer appears on any unload branch. The proposal uses it only for
synthesis, at lines 268, 401, and 447, which is the code's registry meaning
("Synthesis before successful load", §17.9 line 1910, origin `wire`). The unknown branch is now
structural only:

> "Direct `unload` with unknown capability knowledge: C1 tests only the structural no-frame/no-spawn
> boundary; it does not assert a result code." — proposal line 156

The no-spawn rule is attributed correctly:

> "§17.7.2's lazy-start row limits starts to explicit TTS/catalog/status requests, not `unload` or
> `cancel`." — proposal line 123

The case is registered rather than assumed. D10 is an owner decision with no selected answer, no new
code, and candidate directions listed (line 683), gated "Before C1 acceptance", and C1's acceptance
states "no unknown-snapshot result code is claimed before D10" (line 176). The known-absent and
unknown branches are distinguished in rules, tests, and acceptance, so the two cannot be read as
one branch. The round-2 objection — a frozen wire code silently extended, and a decision resolved by
attribution — no longer applies.

**Evidence grade:** Observed; the registry meaning is Source-derived (§17.9 line 1910).

### M1 — C2 `Available` false set — **Closed**

> "It is false for `absent`, `starting`, `handshaking`, `stopping`, `crashed`, `incompatible`, and
> `backoff`, and for `unavailable` caused by a missing helper or invalid provider configuration."
> — proposal line 204

That is A2-D5's enumeration in the same order, including the seventh state (`backoff`) and the
non-legal `unavailable` cause the round-2 audit found missing. The matrix clause names the same
seven states (line 213), and the tests require "the seven aggregate states and the
missing-helper/configuration causes". Both legal `unavailable` rows stay `Available=true` and are
asserted separately. No transition ownership moved into C2; the forward-looking `backoff` row is
marked "not yet reachable", consistent with C5 owning it.

**Evidence grade:** Observed against A2-D5 (line 2506).

### M2 — C1 A2-D3 categories — **Closed, and correctly scoped**

C1 now names both missing semantic categories:

> "add A2-D3 semantic validation, including duplicate extension names and invalid handshake
> descriptors; do not reject wire capability arrays for order—only normalise them and order the
> `detail.missing` projection" — proposal line 141

The negative-injection list requires "duplicate extension name (must become `protocol_incompatible`,
not `protocol_violation`)" and "invalid handshake descriptor", and acceptance repeats the
classification. The earlier "stable capability-order validation" instruction is gone; a grep over
the whole proposal finds no occurrence of it.

This is a real correction, not a wording change. Current code returns `protocol_violation` for a
duplicate extension, and A2-D3 puts both cases under `protocol_incompatible` ("a duplicate provider
ID, capability ID, service ID, or extension name; an invalid handshake descriptor", line 2504), with
`protocol_violation` "reserved for actual protocol-sequencing violations".

One new expectation in the same list is also contract-accurate and worth recording as verified
rather than assumed: "duplicated JSON keys (must remain `invalid_frame`, not
`protocol_incompatible`)" matches A2-D3's explicit reservation of `invalid_frame` for
"unparseable bytes, duplicate JSON keys, oversized frames, or invalid UTF-8". The correction did not
overreach into a stricter classifier.

**Evidence grade:** Observed against A2-D3 (line 2504); current-code behaviour Source-derived from
`helper/protocol.go`.

### M3 — D1a gate placement — **Closed**

Every location now places the gate before C3:

| Location | Text |
|---|---|
| C3 prerequisites (252) | "C2 status model; D1a disposition." |
| C3 scope (250) | "only if D1a permits, the ordinary-load seam" |
| C3 design decision (254) | "D1a must be dispositioned before C3 adds any such seam" |
| C3 acceptance (287) | "D1a is dispositioned before any seam is added; if permitted … if denied … no new adapter/session file exists" |
| §7 inventory (575) | "D1a gate before C3; add only if permitted; D1b trigger" |
| §10 ledger (673) | "Before C3 introduces any loader seam; C3 is not authorizable until resolved" |
| §9.1 (630) | "Can begin only after C1/C2 review and D1a disposition: C3." |
| C7 prerequisites (432) | "D1a already decided before C3" |

The denial branch is now specified rather than left dangling: without D1a, C3 is limited to state
gating and the existing client `Load`, and C7 "records reachability as blocked rather than inventing
a surface". No surviving sentence schedules the seam's permissibility at C7. The round-2
contradiction — C3 creating the seam while its gate was scheduled for C7 — is gone.

**Evidence grade:** Observed.

### M4 — idle reset recorded as D9 — **Documented correctly, but the decision is mis-scoped (see M5)**

The question-form bullet is gone from C5's required behaviour and replaced by an explicit
non-assumption:

> "Idle-deadline reset behavior is **D9** and is not frozen here. … C5 must not assume either
> answer; D9 must be dispositioned before C5 acceptance." — proposal line 359

C5 holds D9 as a prerequisite (342) and acceptance (373), and the observation/respawn tests are
gated on the disposition (361). This is the correction the round-2 audit asked for, and F6 is no
longer an unrecorded question.

The content of D9 is nevertheless wrong in one half; see M5. The *documentation* of the finding is
correct, which is why M4 is closed and M5 is filed separately rather than as a partial closure.

**Evidence grade:** Observed.

### D7 — settled by citation — **Closed, citation verified**

> "**D7 is settled by §17.11.5:** after a successful in-place recovery, readiness resumes and the
> ledger resets only after ten continuous minutes of stable readiness; until then prior crashes
> remain counted. No owner decision or separate C6 gate remains." — proposal line 405

§17.11.5 reads "maximum 3 consecutive restarts per key in 10 minutes; reset after 10 minutes of
stable readiness" (line 2106). The two former alternatives therefore describe one rule, exactly as
the record states. No residual C6 gate survives, and the §10 row is closed (680). The proposal does
not reopen the contract's restart accounting.

**Evidence grade:** Observed; closure Source-derived and independently verified against the
contract, not accepted from the record.

## 4. Regression sweep (F1–F8)

| Correction | Status in the revised text | Note |
|---|---|---|
| F1 unload/cancel split | Preserved, refined by B1 | Path and direct branches are now separately implemented and separately tested. |
| F2 two legal `unavailable` rows | Preserved verbatim | Table at 206–209 matches the §17.7.4 additive rows, including `ServiceCapabilitiesKnown=false` for budget exhaustion. |
| F3 single-source projection table | Preserved, transitions still deferred | "C2 must not claim `backoff` or identity-replacement behavior before C5/C6" (line 230). |
| F4 five-part evidence rule | Preserved | C8 rules 1–5 intact, including "Unestablished" instead of "indirect". |
| F6 idle-reset question | Now recorded as D9 | Addressed — but see M5 for the decision's scope. |
| F7 C4 intent diagnostic | Preserved | "The intent value is diagnostic and does not override those transitions" (line 313). |
| F8 process-handle registry | Preserved and strengthened | Registry "keyed by the real process/job handle", entry removed only after `Wait()` is observed; cap counts "actual helper child processes" (C7). |

No prior correction was narrowed, widened, or contradicted. The §10 introduction now states the
split correctly ("D6 is applied in C4; D7 is settled by §17.11.5; D4 is folded into D2/D5"), and the
§1 gap/decision counts — three contract-surface gaps (D1a, D3, D10), four separate owner decisions
(D1b, D2, D5, D9) — agree with the ledger. The round-2 count mismatch is fixed.

## 5. New findings

### M5 — major — D9 offers an owner choice the contract has already made

D9 is recorded as a single product decision with two halves:

> "Does status/catalog/health observation reset the 120 s idle deadline, and may a status request
> respawn a helper after idle close makes it `absent`?" — proposal line 682

The first half is genuinely open. §17.8 supplies only "Idle close 120 s | Provisional | Stop an idle
helper instance" (line 1864) and "An idle instance may be stopped after a bounded idle period"
(line 521); neither defines whether observation counts as activity, so the deadline-reset rule is
undecided by the contract.

The second half is not open. §17.7.2's transition table contains:

> "| `absent` | explicit TTS/catalog/status request and helper resolved | `starting` | Lazy start |"
> — contract line 1773

A status request from `absent` *is* a lazy-start trigger. A2-D12's restriction is scoped to
recovery — "`synthesize`, `catalog`, voice-list, `health`, and status requests do not trigger that
recovery" (line 2513) — and says nothing about starting from `absent`. The proposal's own §4.4 cites
the same row in the same sense ("limits starts to explicit TTS/catalog/status requests", line 123),
and C7 commits the adapter to the behaviour: "`ProviderStatus`: may start/handshake; no load"
(line 453).

So the proposal asks the owner to decide something the contract settles and its own adapter already
implements. The concrete risk is not a silent freeze — D9 gates C5 acceptance — but a *wrong*
disposition: an owner answering "no, a status poll must not respawn" would oblige C5 to encode a
rule that contradicts §17.7.2 and would contradict C7's `ProviderStatus` contract. That is the
failure mode this audit exists to catch, and it is why M5 is graded major rather than minor.

**Recommended repair (documentation only):** narrow D9 to the deadline-reset question; record the
respawn half as contract-settled by §17.7.2 (a status request from `absent` starts the helper). If
the product intent is to avoid poll-driven respawn, that is an application-layer polling policy —
whether the shell issues status requests to an idle provider — and must be recorded as such rather
than as a C5 transition rule.

**Evidence grade:** Observed (both texts); the contract reading is Source-derived from §17.7.2 line
1773 and A2-D12 line 2513.

### m1 — minor — "exactly at 3" still overstates the budget rule

C5 acceptance reads "Budget exhausts exactly at 3 in 10 minutes" (line 373). §17.11.5 says
"maximum 3 consecutive restarts" (line 2106), so the third restart is permitted and exhaustion
occurs when a fourth is attempted. The accompanying test — "Three restarts then unavailable"
(line 361) — is correct and disambiguates the intent, so this is wording, not a behavioural defect.
It persists from round 2 unchanged.

### m2 — minor — no structural test for a direct `Cancel` with an unknown snapshot

The direct-`unload` unknown branch has an explicit structural test ("C1 tests only the
no-frame/no-spawn boundary", line 156). The direct-`Cancel` unknown branch does not; it appears only
inside the path test bullet (line 161) and in the acceptance sentence that refuses to claim a code.
Since D10 covers "direct `unload` and any direct `Cancel` outside the supervisor path" (line 683),
C1 should assert the same structural boundary for `Cancel` — no frame, no spawn, no result code —
or B2's guarantee holds for only half the surface it names.

### m3 — minor — same-day boundary omits the D10 acceptance gate

§9.1 says C1 and C2 "Can be implemented and fake-helper tested" (line 629) without noting that C1
cannot be *accepted* before D10 is dispositioned, which C1's own prerequisites and acceptance both
state. Implementation and acceptance are different acts, so this is not a contradiction, but the
same-day boundary is the sentence an owner is most likely to act on and should carry the gate.

### m4 — minor — "unsupported-operation frame" is an undefined term

C1 uses "No test emits an unsupported-operation frame" and "no unsupported-operation frame is
emitted" (lines 163, 176) plus the §4.4 sentence"No unsupported-operation frame is emitted in any path" (line 123). §17.9 has no such code or frame type. The intended meaning is A2-D1's "MUST NOT
write an operation frame" when the capability is absent. Aligning the wording with A2-D1 would
remove an ambiguity an implementer could read as a distinct wire artefact.

### m5 — minor — the path-scoped operation's call site is not named

The design requires the supervisor cancellation path to use a path-scoped internal operation
(line 122), but the `client.go` affected-files bullet (line 142) does not name the call site that
must switch to it — the existing cancellation watcher, which currently invokes the exported
`Cancel`. An implementer following only the file list could add the internal operation and leave the
watcher calling `Cancel`, recreating B1 in code while the test at line 162 catches it late. Naming
the call site costs one clause and removes the gap.

None of m1–m5 changes a stage's acceptance contract at C1. m2 and m5 are the two most likely to
reappear as implementation defects; m1 is the most likely to be read literally.

## 6. Stage authorization matrix

Stated as the audit's read of the *design*; nothing here authorizes work.

| Stage | Gate status after this pass | Note |
|---|---|---|
| C1 | Implementable; **acceptance blocked on D10** | B1 and B2 closed; M2 closed. No blocking design defect. The unknown-snapshot branch may be built structurally but not accepted with a code. |
| C2 | Clear | M1 closed; `Available` set complete against A2-D5. Transition ownership still deferred to C3–C6 by design. |
| C3 | **Blocked on D1a** | Correctly gated now; denial branch specified. M3 closed. |
| C4 | Clear of design findings | Depends on C3. D6 is applied, not open. Windows evidence owed at C8. |
| C5 | **Blocked on D9 acceptance** | M4 closed as documentation; M5 must be repaired before D9 is dispositioned. Budget wording (m1) should be tightened first. |
| C6 | Clear of design findings | D7 closed by citation; A2-D7/A2-D12 split is stated. |
| C7 | Blocked on D1b, D2, D3, and the remaining owner decisions | D3 remains a real error-surface gap; the proposal correctly refuses to invent a code. C7's `ProviderStatus` behaviour is entangled with M5. |
| C8 | Not reachable | Requires C1–C7 and separate runtime authorization. |

## 7. Decision ledger after review

| ID | Audit view |
|---|---|
| D1a | Genuine contract-surface gap; gate placement now correct. |
| D1b | Genuine product decision; correctly deferred to C7. |
| D2 | Genuine product decision; correctly deferred. |
| D3 | Genuine gap; correct refusal to repurpose `busy` / `helper_start_failed` / `unsupported_capability`. |
| D4 | Correctly folded. |
| D5 | Genuine product decision; correctly deferred to UI. |
| D6 | Correctly applied, not open. |
| D7 | Correctly closed by §17.11.5; verified. |
| D9 | **Partly misclassified** — the respawn half is contract-settled (§17.7.2) and already implemented by C7's `ProviderStatus`. Repair before disposition. |
| D10 | Correctly registered; correct scope (direct `unload` and direct `Cancel`); correctly gated before C1 acceptance; no code frozen. |
| D8, O-1–O-7 | Unchanged and out of this pass's scope. |

## 8. Verification performed

- Recomputed all five hashes and sizes; the record's §1 table matches disk.
- Read the proposal (710 lines) and record (175 lines) in full; no section was sampled.
- Verified every A2-D1, A2-D3, A2-D5, A2-D7, A2-D8, A2-D9, A2-D12, §17.7.2, §17.8, §17.9, and
  §17.11.5 quotation against the contract text at its stated line.
- Grepped for the two claimed removals ("stable capability-order validation", a frozen
  unknown-snapshot code) and confirmed both are absent from the proposal.
- Confirmed `provider_not_loaded` survives only on the three synthesis paths.
- No file was modified. No build, test, replay, or process launch was performed.

## 9. Limitations

- This review was performed by the author of the round-2 audit. It verifies documentation fidelity
  and re-derives contract citations independently, but it is not an independent third party for the
  findings it originated. A separate reviewer remains warranted before C7.
- The superseded proposal was overwritten in place, so no byte-level diff against `0f920b24…` was
  possible; the record's §11 section list was the only available change manifest, and each claimed
  change was located and read in the current text rather than assumed.
- Implementation facts cited (current `protocol_violation` classification, the cancellation watcher
  call site) are read-only observations of source, not executed behaviour.
