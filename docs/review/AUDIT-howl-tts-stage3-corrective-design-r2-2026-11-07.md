# AUDIT (round 2) — Stage 3 corrective design, corrected revision

**Date:** 2026-11-07
**Mode:** DOCUMENTATION REVIEW ONLY
**Reviewer position:** independent; not the author of the submitted proposal or correction record.

**Artifacts reviewed**

| Role | Path | SHA-256 (verified) |
|---|---|---|
| Governing contract | `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f` |
| Corrected proposal | `docs/review/DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md` | `0f920b247b8ba16bb78591916ae777d006e46190f5459cc8f738293527acc270` |
| Correction record | `docs/review/CORRECTIVE-howl-tts-stage3-corrective-design-2026-11-07.md` | `2eeaaf14ec9b80d0807453ca8029eac3c78c383de3b674cd8f193601b9565e57` |

**Status of this document:** review artifact. It authorizes nothing, edits nothing it reviews, and does not
alter the governing contract.

## 1. Verification of the supplied artifacts

All three supplied hashes match the values in the owner's brief exactly. Sizes and timestamps are consistent
with the stated revision sequence:

| Artifact | Lines | Bytes | mtime |
|---|---:|---:|---|
| Corrected proposal | 703 | 59,634 | 2026-11-07 09:43:07 |
| Correction record | 134 | 11,452 | 2026-11-07 09:44:13 |
| Contract | 2,586 | 182,904 | unchanged |

Preservation claims in correction record §9 verify:

| Claim | Result |
|---|---|
| Contract byte-identical at `ec3c789e…` | **Confirmed** (recomputed) |
| Readiness investigation byte-identical at `a63f9f32…` | **Confirmed** (recomputed) |
| Prior audit byte-identical at `71a17369…`, 25,962 bytes | **Confirmed** (recomputed) |
| No production source/test/fixture modified | **Consistent** with filesystem state: only four documentation files carry a 2026-11-07 modification time after the audit |
| Proposal corrected in place; record and README added | **Observed**: README §24 carries the new proposal hash `0f920b24…`; §26 is the new record entry; §25 (prior audit) is unchanged |

The correction record's §1 also states the prior audit hash was computed from disk rather than copied. The
recomputation agrees, so that provenance claim holds.

**Scope of my inspection.** Documentation review plus read-only source inspection for three specific claims
that the documents alone cannot settle: the current `Client.Unload`/`Client.Cancel` gating, the current
`validateHello` classification of duplicate extension descriptors, and whether `Client.Cancel` is reachable
outside the supervisor cancellation path. Those readings are graded **Observed (source)** below. No test,
build, execution, replay, or process launch was performed. No claim in this review depends on runtime
behavior.

## 2. Disposition

**PASS WITH BLOCKERS.**

Meaning precisely: the correction strategy is sound and five of the prior findings (F1 core, F2, F3, F4, F8)
are genuinely corrected rather than merely acknowledged; the eight-stage decomposition remains the right
shape. But two defects sit inside C1 and would be frozen into its acceptance tests, and one dependency
defect sits in the C1→C2→C3 ordering. "Pass with blockers" therefore means: accept the corrected design as
the design basis, but do not treat C1 as authorizable until B1 and B2 are repaired, and do not treat C3 as
authorizable until the D1a timing in M3 is repaired.

This is not a FAIL: nothing in the corrected proposal contradicts adopted §17.19 on the two matters that were
blocking before (F1's `unload`/`cancel` split and F2's two `unavailable` cases). It is not a clean PASS: the
replacement wording for C1 introduces new semantic surface that the record presents as settled.

## 3. Executive findings

1. **The two prior blocking findings are fixed, and one of them is fixed exactly.** C2's `unavailable` table
   now matches the contract's additive §17.7.4 rows field-for-field, including `ServiceCapabilitiesKnown=false`
   for the budget-exhausted case. C1 no longer assigns `unsupported_capability` to cancellation. These are
   real corrections, not rewording.
2. **B1 (blocking, C1): the cancellation exemption is written at method scope, but A2-D9 grants it at path
   scope.** `Client.Cancel` is an exported direct operation; the corrected C1 makes it skip the frame and
   return no error whenever `cancel` is absent. A2-D9 protects "the supervisor's cancellation path", while
   A2-D1's general rule still assigns `unsupported_capability` to an absent optional capability on a direct
   operation. A direct call would now report success without cancelling anything.
3. **B2 (blocking, C1): unload with an unknown capability snapshot returns `provider_not_loaded`.** That is a
   frozen wire code whose registry meaning is "Synthesis before successful load" (§17.9 line 1910). The
   corrected design resolves an unlegislated case by reusing a code it elsewhere forbids repurposing (D3), and
   attributes the no-spawn rule to §17.19 rather than to the clause that actually establishes it (§17.7.2
   line 1773).
4. **M3 (major, sequencing): D1a gates C3, not C7.** C3 "introduces that seam" (lines 251–255), while §7 and
   §10 place the seam's addition in C7 (lines 572, 668). The decision that permits the seam is scheduled after
   the stage that creates it.
5. **M1 (major, C2): `Available`'s false set is incomplete.** A2-D5 enumerates `absent`, `starting`,
   `handshaking`, `stopping`, `crashed`, `incompatible`, `backoff` as false; C2 states only "start/configuration
   failure and incompatible" (line 202). The prior F2 fix repaired the true cases and left the false cases
   under-specified.
6. **M2 (major, C1): the test list omits two A2-D3 semantic categories.** Duplicate extension names and an
   invalid handshake descriptor are not enumerated, and current code classifies a duplicate extension as
   `protocol_violation` — the exact misclassification C1 exists to remove.
7. **M4 (major for C5): F6 was not addressed.** The `?`/`Define:` bullet is still in C5's *required behavior*
   (line 356), and the idle-reset question it encodes is a real undecided product rule absent from §10.
8. **D7 is smaller than the record claims.** Its residual collapses under §17.11.5's "reset after 10 minutes
   of stable readiness"; the two alternatives the record poses have the same effect. No owner decision is
   required.

## 4. Finding ledger

Severity is judged against the owner's gate: **blocking** means the defect would freeze a wrong contract
interpretation into a stage's acceptance tests. Grades: **Observed**, **Source-derived**, **Calculated**,
**Inferred**, **Unestablished**, **Superseded**.

| ID | Severity | Grade | References | Defect or confirmation |
|---|---|---|---|---|
| B1 | Blocking (C1) | Observed + Source-derived | Proposal C1 (142, 148–174) and §7; contract A2-D1 (2502), A2-D9 (2510), A2-D10 (2511); `client.go:622`; `client_test.go:407` | Cancellation skip is applied to the `Cancel` method, not only the cancellation path |
| B2 | Blocking (C1) | Observed + Source-derived | Proposal 123, 156, 174; contract §17.9 (1910), A2-D1/A2-D7, §17.7.2 (1773) | Unknown-snapshot unload reuses a frozen synthesis-scoped code and is presented as §17.19-required |
| M1 | Major (C2) | Source-derived | Proposal 202; contract A2-D5 (2506) | `Available` false-state enumeration incomplete |
| M2 | Major (C1) | Observed + Source-derived | Proposal C1 Tests (148–166); contract A2-D3 (2504); `protocol.go` `validateHello` | Two A2-D3 semantic categories missing from the C1 test list; current misclassification survives |
| M3 | Major (sequencing) | Observed | Proposal 251–255, 572, 668 | D1a is scheduled before C7 but the seam it governs appears in C3 |
| M4 | Major (C5) | Observed | Proposal 356; §10 | Unresolved idle-reset rule embedded in normative text; no decision entry (F6 unaddressed) |
| m1 | Minor | Observed | Proposal 36 vs §10 | Executive recommendation undercounts gaps and decisions |
| m2 | Minor | Observed + Source-derived | Proposal 141; contract A2-D3 (2504) | "stable capability-order validation" ambiguous; order is mandated only for `detail.missing` |
| m3 | Minor | Source-derived | Proposal 370; contract §17.8 | Acceptance fixes a provisional numeric as exact |
| m4 | Minor | Observed + Source-derived | Proposal 463; contract §17.8 (1857) | Cap is enforced as cardinality, but §17.8 names a composition ("Global + one explicit engine") |
| m5 | Minor | Observed + Source-derived | Proposal C5/C7; contract §17.7.2 (1773) | Idle close can be defeated by status polling; design does not state whether status may respawn |
| m6 | Minor (ledger) | Source-derived | Proposal §10 D7; contract §17.11.5 (2087–2088) | D7 residual is already settled by citation |
| m7 | Minor | Observed | Proposal §10 D1a; contract §17.12.1 (2255) | D1a under-describes the gap: §17.12.1 also has no `unload`/`cancel` rows |
| m8 | Minor | Observed | Proposal 154–160, C1 Tests | "No spawn" is asserted at a layer that owns no process |
| m9 | Minor | Source-derived | Proposal C5; contract A2-D12 | Budget-exhausted caller outcome derivable but not restated |
| m10 | Minor | Observed | Proposal C4; contract A2-D9 | Provider state after a cancellation terminal without kill is not stated |

### B1 — Cancellation: the exemption is path-scoped, the change is method-scoped

A2-D9 reads: "If it did not, the adapter skips the cancel frame; absence of `cancel` is not an error in **the
supervisor's cancellation path** and no `unsupported_capability` is synthesized there." The scoping phrase is
not incidental. A2-D1 states the general rule without any path qualifier: "The adapter returns
`unsupported_capability` for an absent optional service capability and MUST NOT write an operation frame."

C1 now says: "The cancellation path sends `cancel` only when the snapshot is known and advertises
`howl.tts.cancel`; when absent or unknown it skips the frame, synthesizes no `unsupported_capability`…"
(proposal line 142). That sentence is correctly path-scoped. But the machinery it mandates lives in
`Client.Cancel`: §7's inventory assigns "cancellation path / `Cancel`" to C1 with "skip-frame/no-error only",
and C1's tests assert "no cancel frame, no `unsupported_capability`" without distinguishing the path from a
direct call.

Source inspection confirms the distinction is live, not theoretical:

- `Client.Cancel` is exported and gates through `requireServiceCapability(CapabilityCancel)` (`client.go:622`),
  which today returns `unsupported_capability` when the capability is absent (`client.go:877`).
- The only production caller is `watchCancellation` (`client.go:871`), which discards the error; tests call
  `Cancel` directly and assert a nil error (`client_test.go:407`).

Consequence: after C1, a direct `Cancel` against a `cancel`-less helper returns nil and writes nothing. That is
a successful no-op for an absent capability, which A2-D10 forbids ("An absent capability is never treated as
supported"), and it removes the `unsupported_capability` that A2-D1 assigns to the direct operation. The
supervisor path is entitled to swallow the error; the exported method is not, unless it stops being a direct
operation.

Required resolution: distinguish the two call sites in C1's specification. The cancellation path records
`cancel_requested` and proceeds through grace/`KillTree`/`Wait` with no frame and no capability error
(A2-D9). A direct `Cancel` invocation with the capability known absent returns `unsupported_capability` with
no frame (A2-D1). If the intent is that `Cancel` is never called directly, say so and make the path-only
nature structural (separate internal function or an explicit flag), and keep a test that a direct call still
reports the contract's direct-operation result. C1's tests must assert the two separately, exactly as F1
required for `unload`.

### B2 — Unload with an unknown snapshot reuses a frozen code

The corrected behaviour is: "Direct `unload` before handshake/with unknown capability knowledge:
`provider_not_loaded`, no frame, no helper start" (line 156), accepted at line 174.

Three problems, in increasing order of consequence.

First, the *no-spawn* rule is correct but mis-attributed. §4.4 presents it as "New behavior required by
§17.19". §17.19 does not address unknown snapshots for optional operations. The rule is actually grounded in
§17.7.2's lazy-start row (line 1773), which limits starts to "explicit TTS/catalog/status request"; `unload`
and `cancel` are not triggers. Citing the right clause makes the rule defensible rather than invented.

Second, `provider_not_loaded` is a frozen **wire**-origin code with a narrow meaning: "Synthesis before
successful load" (§17.9 line 1910). No synthesis is attempted in this path. The design would have the adapter
mint a wire-origin code locally for a non-synthesis condition — and it does so while D3 argues, correctly,
that `busy`, `helper_start_failed`, and `unsupported_capability` cannot be repurposed for admission
rejection. The same principle applies here. Note that the *legitimate* uses of this code elsewhere in the
proposal (C3 line 265, C6 line 398, C7 line 444, for synthesis while `unavailable`) are contract-explicit in
A2-D12 and therefore fine.

Third, and most important: this is a fourth caller-facing semantic gap, and it is not registered as one. §10
lists D1a and D3 as gaps; the unknown-snapshot unload disposition is a decision with the same character —
the contract is silent, so whatever the adapter returns becomes the de facto contract. A2-D7 assigns
`unsupported_capability` only when the capability is *absent*; A2-D9's cancellation rule is about
*advertised* versus not. Neither says what "unknown" means for a direct unload.

Required resolution (owner decision; this review does not choose): register the case as a decision entry of
the same class as D3, and either (a) accept `provider_not_loaded` with a recorded note that its registry
meaning is being extended for adapter-origin fail-closed use, which is a semantic amendment and should be
acknowledged as such; (b) return `unsupported_capability` on the argument that an unverified capability is
never treated as present, acknowledging that this asserts absence it has not established; or (c) remove the
unknown branch from C1's acceptance and let C1 implement only the known-present and known-absent cases. Until
that is recorded, C1's acceptance row freezes an unlegislated choice.

### M1 — `Available` false enumeration

C2 line 202 states: true for configured/compatible in `unloaded`/`loading`/`ready`/`busy`/`degraded` and both
`unavailable` cases; "false for start/configuration failure and for incompatible."

A2-D5 enumerates the false set explicitly: "false for `absent`, `starting`, `handshaking`, `stopping`,
`crashed`, `incompatible`, and `backoff`, and for `unavailable` caused by missing helper or invalid provider
configuration." C2's sentence covers one of those seven states plus one `unavailable` cause. The omitted
states are not exotic: `crashed` and `backoff` are exactly the states C4/C5 create, and C2 is the stage that
freezes the projection.

This is an omission, not a contradiction, so it is not blocking in the same sense as B1/B2 — but C2's
acceptance ("every status field has the single documented source") cannot compensate for a rule that does not
say what the field equals in six states. Required: restate A2-D5's full enumeration in C2 and add the
corresponding matrix rows (at minimum `absent`, `handshaking`, `stopping`, `crashed`).

### M2 — C1's test list does not cover all of A2-D3

A2-D3 enumerates the semantic `hello` failures that must become `protocol_incompatible`: missing required
service capability; unknown core capability ID; wrong capability or extension version; duplicate provider ID,
capability ID, service ID, or extension name; invalid handshake descriptor; service-only capability in a
provider descriptor; provider-only feature in the service array.

C1's Tests (lines 148–166) cover missing required, unknown service ID, wrong version, duplicate capability,
duplicate provider ID, duplicate service ID, the two scope errors, and the shared `synthesize` acceptance.
They do not cover **duplicate extension name** or the generic **invalid handshake descriptor**. Source
inspection shows the current `validateHello` returns `protocol_violation` for a duplicate extension, which is
precisely the classification A2-D3 reverses.

Consequence: C1 could pass its stated tests while leaving one frozen reclassification unimplemented, because
nothing in the stage observes it. Required: enumerate every A2-D3 category in C1's test list, including the
extension duplicate and one malformed-descriptor case.

### M3 — D1a is sequenced after the artifact it governs

- C3's design decision: "C3 introduces that seam in a narrow, testable form; C7 makes it the production
  `TTSSynthesizer`" (line 251); affected files include "a new internal adapter/session file … limited to
  `Start`, `Load`, `Synthesize`, `State`" (line 255); required behavior includes "`Load` is explicit".
- §7 inventory: "Explicit loader seam | internal non-Wails loader interface | Add only if D1a permits … | C7"
  (line 572).
- §10 D1a: "Before any C7 loader interface or seam is added" (line 668).

D1a's own question is whether "the implementation [may] add a narrow internal, non-Wails loader seam outside
§17.12.1's enumerated adapter behavior". That is C3's artifact, not C7's. As written, C3 creates the object
whose permissibility is undecided, and the decision that would authorize it is scheduled two stages later.

Consequence: C3 is not honestly authorizable while D1a is open, and C7's gating statement is misleading.
Required: either (a) D1a's needed-by moves to "before C3 introduces the seam", with C3 authorised only after
the decision; or (b) C3 is narrowed to state gating and an internal `Load` on the existing client, with the
"seam" language and the adapter/session file moved into C7. The design should pick one and make §7, §10 and
C3 agree.

### M4 — F6 remains, and it hides a real decision

C5 still contains, inside **Required behavior**: "Status/catalog/health do not reset idle indefinitely?
Define: status/health do not reset the idle deadline unless contract explicitly says; accepted work/load
resets it" (line 356).

The correction record's dispositions cover F1, F2, F3, F4, F8; F6 is not addressed, and the record does not
claim it is. The text is not merely editorial: it is the only place the idle-reset rule appears, it is phrased
as a question with an instruction to the author, and §10 contains no decision entry for it.

The underlying question is real. §17.8 supplies the 120 s value and A2-D8 supplies the trigger, but neither
says whether observation resets the deadline. The contract's §17.7.2 lazy-start row compounds it: after idle
close the helper is `absent`, and an explicit status request may start it again — so "does status reset idle"
and "may status respawn after idle close" are two halves of one product rule. Required: replace the bullet
with a stated rule and register it as a decision if the rule is not derivable (see D9 below).

### Minor findings

- **m1.** §1 (line 36) says the design identifies "one contract-completeness gap and two genuine
  product/architecture decisions". §10 lists two gaps (D1a, D3) and three decisions (D1b, D2, D5). The
  executive summary should match the ledger, because it frames what the owner is being asked to decide.
- **m2.** C1's `protocol.go` change list includes "add stable capability-order validation" (line 141). A2-D3
  mandates §17.6.1 table order only for the `detail.missing` string. Nothing requires the wire capability
  arrays to be ordered; a validator that rejects an unordered array would reject a conformant helper. State
  that this normalises the projection (or orders `missing`) rather than validating wire order.
- **m3.** C5's acceptance says "Budget exhausts exactly at 3 in 10 minutes" (line 370). §17.8 classifies the
  budget as **Provisional** and says such values "MAY be changed"; C5's own exclusions say the numerics are
  provisional. Also ambiguous: whether three restarts are permitted and the fourth refused, or only two are
  permitted. Reword as a behavioural rule tied to the configured value.
- **m4.** C7 counts processes (line 463), which fixes F8. §17.8's row reads "Two helper processes | Frozen |
  Global + one explicit engine; hard cap". C7 enforces cardinality ≤ 2 but never addresses the composition —
  whether two explicit-engine processes without a global one are conformant, or whether the second process
  must be the explicit engine. One sentence, plus a test if the composition is normative.
- **m5.** Idle close and lazy start interact: A2-D8 closes an idle helper, and §17.7.2 line 1773 lets an
  explicit status request start one. A renderer that polls `ProviderStatus` therefore closes and reopens the
  process. This is not a contract violation, but the design should state the intended behaviour, because it
  determines whether idle close is effective in production.
- **m6.** D7's residual is already settled. §17.11.5: "maximum 3 consecutive restarts per key in 10 minutes;
  reset after 10 minutes of stable readiness" (lines 2087–2088). After a successful in-place recovery,
  readiness resumes and the ledger resets when ten continuous minutes of readiness have elapsed; until then
  prior crashes remain counted. "Starts a fresh ten-minute window" and "leaves the prior ledger until ten
  continuous minutes" describe the same rule. The residual can be closed by citation, not by decision.
- **m7.** D1a describes §17.12.1 as lacking a `load` entry point. It also lacks `unload` and `cancel` rows
  (the table spans lines 2257–2265). A2-D7's "direct `unload` request" therefore has no adapter surface; the
  only unload surface is internal (`Client.Unload`). Say so, so D1a is not read as "the only missing method is
  `load`".
- **m8.** C1's "no spawn" assertions are made at the client layer, which owns no process. "No helper start"
  is guaranteed by construction there and can pass while a supervisor-level caller spawns. Either assert at
  the supervisor/adapter boundary or state that the property is structural.
- **m9.** C5 defines the budget-exhausted state and flags but not the caller-facing outcome of the next
  synthesis. A2-D12 supplies it (`provider_not_loaded`, no retry); C5's acceptance should name it so the test
  is determinate rather than incidental.
- **m10.** For cancellation racing a success frame, the request outcome is fully defined (A2-D9: success
  rejected, never cached) and cleanup is defined. What is not stated is the provider's aggregate state after a
  cancellation terminal that does not trigger `KillTree`. C4 should state it, since C4 owns exit/status truth.

## 5. Correction coverage matrix

Each row maps a record finding to the proposal text and to this review's assessment. "Complete" means the
corrective intent is implemented and survives contract comparison.

| Finding | Proposal location | Record claim | Assessment |
|---|---|---|---|
| F1 — unload/cancel conflation | §4.4 (123); C1 (142, 148–174) | §2 | **Partially corrected.** `unload` and `cancel` are now separate; A2-D9 is respected for the path; C1 tests separate the cases. Residual **B1**: method versus path scope. |
| F2 — one of two `unavailable` cases | C2 (186–215, 229) | §3 | **Complete.** The table matches contract §17.7.4 additive rows verbatim, including `ServiceCapabilitiesKnown`. The old "never suppress a live-helper unavailable" sentence is gone. |
| F3 — `Loaded` double ownership | C2 (186–227); single-source table (216–227); C3 (251) | §4 | **Complete.** C2 projects only; transitions are attributed to C3/C4/C5/C6; C2's acceptance defers backoff/replacement. Prerequisite cycle removed. |
| F4 — Windows evidence cannot close indirectly | C8 (500–511), matrix preamble (§8) | §5 | **Complete.** Five numbered rules; "Unestablished" replaces "indirect"; every row must name environment and boundary; C1 inheritance explicitly scoped. |
| F6 — C5 question-form rule | C5 (356) | not claimed | **Unaddressed.** Still present; see M4. |
| F7 — `timeout-cancel` intent | C4 (310) | not claimed (folded into "D6 applied") | **Substantively corrected.** "The intent value is diagnostic and does not override those transitions." Not recorded as its own disposition, so the record under-reports the change. |
| F8 — cap counted over adapters | C7 (463–473) | §6 | **Complete.** Registry keyed by process/job handle, entry removed after `Wait()`, admission before `proc.Start`, never derived from adapter/engine/character. Tests cover shared adapter, shared process, shared label, zero-process adapter, reaped slot, third process. Residual **m4** (composition). |
| D1a — load surface | §4.4 (123); C7 (456); §10 (668) | §7 | **Classification correct, timing wrong.** See **M3**; see **m7** for the unload/cancel omission. |
| D1b — trigger and lifetime | C7 (457); §10 | §7 | **Complete as a decision.** Correctly left to the owner, correctly coupled to the Apply/rebuild path. |
| D2 — helper-backed activation | §10 | §7 | **Complete as a decision.** Kept distinct from O-4; no assumption made. |
| D3 — admission rejection code | C7 (463–473); §10 | §7 | **Accurate.** `busy` is wire-origin and operation-scoped; `helper_start_failed` is a start failure; `unsupported_capability` is optional-capability-only. The "do not invent or repurpose" instruction is correct — and is the reason **B2** is a defect. |
| D4 — inspector consumer | §10 | §7 | **Correctly folded** into D2/D5. |
| D5 — renderer/O-1 | §10 | §7 | **Complete as a decision.** Static fallback preserved; no eager start. |
| D6 — exit classification | C4 (310) | §7 | **Settled and applied.** Matches §17.7.2 and A2-D8/A2-D9. |
| D7 — stable-readiness reset | C6; §10 | §7 | **Over-deferred.** See **m6**; contract supplies the rule. |
| (no record entry) | C1 (141) capability-order validation | — | **Unaccounted change** if newly added in this revision; **Unestablished** whether new, since the prior proposal version was replaced in place. |

Proposal changes not accounted for in the record: F7's diagnostic sentence (the record folds it under "D6
applied" without naming it); the C1 capability-order clause (see above). No other unaccounted change was
found; the record's §8 section list matches the corrected document.

## 6. C1–C8 stage matrix

"Ready" here means *design-ready for an authorization question*, never that authorization is granted.

| Stage | Prerequisites genuinely met by earlier stages | Scope boundary | Exit criteria observable | Depends on later stage | Ready |
|---|---|---|---|---|---|
| C1 | Design acceptance only — **currently blocked** by B1/B2 | Capability policy + gating; no status fields | Yes, except "no spawn" (m8) | No | **No** |
| C2 | C1 snapshot/scope separation | Projection only; transitions attributed to C3–C6 | Yes after M1 | No (backoff/replacement correctly deferred) | **After C1 + M1** |
| C3 | C2 status model | Load seam + state gate | Yes | **Yes — depends on D1a (M3)** | **No** until M3 fixed |
| C4 | C2 fields, C3 `Loaded` transitions | Exit intent, shutdown, cancellation truth | Yes; Windows evidence legitimately deferred to C8 | No | After m10 |
| C5 | C4 intent-aware exits | Budget/backoff/idle | Not fully — D9 open (M4) | No | **No** until D9 |
| C6 | C3 seam, C4 termination, C5 budget | Identity/recovery/replacement | Yes; D7 no longer blocks | No | After C5 |
| C7 | C1–C6; D1a/D1b/D2/D3 | Adapter + reachability | Partly — D3 code undecided; D2/D5 open | No | **No** (owner decisions) |
| C8 | C1–C7 + execution authorization | Evidence only | Yes by construction | No | **No** |

Notes:

- **C3's later-stage dependency is the material sequencing defect.** Everything else in the chain is
  dependency-ordered: C4's intent truth precedes C5's timers, C5's budget precedes C6's accounting, C6's
  behavior precedes C7's reachability, and C7 precedes C8. The "do not combine" list remains correct,
  particularly "exit truth must precede timers/budgets".
- **C2's forward references are now honest.** Its `backoff` and replacement rows are marked not-yet-reachable
  and deferred, which removes the earlier cycle without weakening the target invariant.
- **C1 cannot authorize implementation or runtime activity outside its boundary** as written: it changes
  protocol/client/supervisor storage and fixtures only, with no production importer. B1 and B2 are semantic
  defects, not scope leaks.
- **C7 is the broadest stage** (adapter, pipeline, factory, app hook, cap registry) and the only one that
  touches production wiring. Its default-off gate and its explicit "no production activation until D2" are
  adequate containment, but it must not be authorized while D2/D5/D3 are open.

## 7. Unresolved-decision ledger

**Owner decisions (must not be assumed)**

| ID | Decision | Why it is the owner's | Stage |
|---|---|---|---|
| D1b | Which existing explicit action invokes the loader; whether the adapter survives it | Product/architecture; depends on the Apply/rebuild lifecycle | Before C7 |
| D2 | How a helper-backed engine is activated without changing `TTSSynthesizer`; relation to O-4 | Product policy | Before C7 |
| D5 | Renderer exposure and O-1 picker behaviour | Product policy | Before UI integration |
| D9 (new) | Does status/catalog/health observation reset the idle deadline, and may status respawn a helper after idle close? | The contract supplies the 120 s trigger but not the reset rule | Before C5 |
| D10 (new) | What a direct `unload` returns when the capability snapshot is unknown | Contract is silent; the answer becomes de facto contract | Before C1 acceptance |

**Contract gaps (disposition required; not product policy)**

| ID | Gap | Status |
|---|---|---|
| D1a | §17.12.1 enumerates neither `load` nor (also) `unload`/`cancel`; A2-D12 requires an explicit load trigger | Correctly identified; **mis-scheduled** (M3) |
| D3 | No §17.9 code exactly defines two-process admission rejection | Correctly identified; must not be papered over |
| D10 | Unknown-snapshot optional-operation semantics | **Newly identified**; currently resolved by assumption in C1 |

**Implementation questions (decidable inside a stage, with the contract as the tiebreak)**

- Capability-array order: normalisation versus validation (m2). Contract requires order only for `detail.missing`.
- Post-cancellation provider state when no kill occurs (m10).
- Observability of "no spawn" at the correct layer (m8).

**Closed by citation — no decision needed**

- **D6** — §17.7.2 `stopping` + grace → `crashed`; A2-D8/A2-D9 for outcomes. Correctly applied.
- **D7** — §17.11.5's reset rule resolves the residual (m6).
- **D4** — folded into D2/D5.

**Evidence gates**

- C8's Windows rows: real child process, Windows environment, named commands and results. Fake-helper and
  CI results support only. Correctly specified.
- Native provider: out of scope; absence must be recorded **Unestablished**, not substituted. Correctly
  specified.
- C1's Windows inheritance: scoped to C1 and cannot close C4–C8. Correctly specified.

## 8. Failure-path audit

"Unique?" asks whether the documents define exactly one outcome. Citations are to the corrected proposal and
the contract.

| # | Case | State | `ProcessAlive` | Caps known | Request outcome | Next allowed action | Unique? |
|---|---|---|---|---|---:|---|---|
| 1 | Direct `unload`, `unload` known absent | unchanged | unchanged | true | `unsupported_capability`, no frame | none | **Yes** (A2-D1/A2-D7) |
| 2 | Direct `unload`, snapshot unknown | unchanged | unchanged | false | proposal: `provider_not_loaded`, no frame, no spawn | handshake only via TTS/catalog/status | Proposal says yes; **semantics contested (B2)** |
| 3 | Cancellation path, `cancel` known absent | request cancelled | unchanged | true | no frame, no capability error; grace → `KillTree`/`Wait`; `cancelled`/`timeout` | cleanup after observation | **Yes** (A2-D9) |
| 4 | Direct `Cancel` call, `cancel` absent | unchanged | unchanged | true | proposal implies nil (no-op); A2-D1 implies `unsupported_capability` | — | **No (B1)** |
| 5 | Asset-specific load failure, live helper | `unavailable` | true | true | typed load failure; `Loaded=false`, `Ready=false`, `Available=true` | explicit same-identity load only (A2-D12) | **Yes** |
| 6 | Budget exhausted after process exit | `unavailable` | false | false | next synthesis: `provider_not_loaded`, no retry (A2-D12) | no automatic restart; update-gated reset | **Yes**, though C5 does not restate it (m9) |
| 7 | Explicit same-identity recovery | `unavailable` → `loading` → `ready` | true | true | one load; repeated failure stays `unavailable` | explicit load again only | **Yes** (A2-D12) |
| 8 | Different-identity replacement | `unavailable`/ready → stop → new process | true → false → true | invalidated then re-established | no differing-identity load; failed start → `helper_start_failed` | A2-D7 order; budget exempt | **Yes** |
| 9 | Cancellation racing a success frame | request retired | unchanged | true | success rejected; never returned or cached | cleanup per §17.11.4 step 7 | Outcome **yes**; post-terminal provider state **not stated (m10)** |
| 10 | Grace expiry while `stopping` | `crashed` | false | invalidated | `cancelled`/`timeout` | none | **Yes** (§17.7.2; D6) |
| 11 | Third helper-process admission | unchanged | unchanged | unchanged | proposal: rejected before `proc.Start`; **code undefined** (D3) | none | **No — acknowledged gap** |
| 12 | Helper exit during startup / load / synthesis / shutdown | `crashed` / `crashed` / `crashed` / `0→absent`, `≠0→crashed`, grace→`crashed` | false | invalidated | failed start; load failure; `helper_crashed`; shutdown classification | restart per budget (not for incompatible) | **Yes** |
| 13 | Late output after retirement | retired | — | — | never returned or cached | exact-path cleanup after termination/observation | **Yes** (§17.11.4 step 7; A2-D9) |
| 14 | Idle close, then status poll | `absent` → `starting` | false → true | false → re-established | status may start (§17.7.2) | — | Contract-defined; **product intent unaddressed (m5)** |

Cases 2, 4, 9, 11, and 14 are where the documents do not give a single clean answer; 2 and 4 are the ones
that touch C1.

## 9. Evidence and limitations

- The three supplied hashes, sizes, and the preservation claims were **recomputed** from disk. They match.
- Contract quotations in this review are from the `ec3c789e…` text; §17.19 controls where §17.19.3 lists a
  supersession, and I checked the specific supersession rows relied on here (A2-D4 for `ProcessAlive`, A2-D5
  for `Available`/`Loaded`/`Ready`, A2-D6 for the status shape, A2-D12 for §17.12.1).
- Three facts were checked in source, read-only: `Client.Unload`/`Client.Cancel` gating
  (`client.go:599`, `622`, `877`), `validateHello`'s duplicate-extension classification, and the direct
  `Cancel` call in `client_test.go:407`. Nothing was executed.
- I did **not** verify runtime behavior, Windows process semantics, or any test result. No claim here rests
  on execution.
- I could not compare the corrected proposal to its predecessor text: the file was replaced in place and the
  prior version is not retained in the workspace. Any finding phrased "newly added" is therefore limited to
  what the correction record claims, except where the record's own section list asserts the change.
- The correction record's claim that no non-docs file changed is consistent with filesystem timestamps; I
  did not audit repository history (the review directory is untracked).

## 10. Recommendation

**Accept the corrected design as the design basis.** The decomposition, the invariant preservation, the C8
evidence rules, and the D-ledger reclassification are sound, and the F1/F2/F3/F4/F8 corrections genuinely
land.

**Do not treat C1 as authorizable yet.** Two repairs, both narrow and both inside C1's text:

1. **B1** — separate the cancellation *path* from the direct `Cancel` operation. Path: skip, no capability
   error, grace/kill (A2-D9). Direct call: `unsupported_capability`, no frame (A2-D1). Test both.
2. **B2** — register the unknown-snapshot `unload` disposition as a decision entry (D10) instead of freezing
   `provider_not_loaded` in acceptance; cite §17.7.2 for the no-spawn rule rather than §17.19.

**C2's prior blockers are cleared.** Complete M1 (the `Available` false set) as part of C2's own text; it does
not need a separate gate.

**Before C3:** fix M3, because D1a governs the artifact C3 creates.

**Before C5:** replace F6's bullet with a rule and record D9.

**Before C7:** D1b, D2, D5, D3, and D10 must be decided by the owner; none may be assumed by the implementer.

**Explicit exclusions for the next step:** no C3–C8 authorization, no production wiring, no runtime evidence
collection, no contract amendment (D10 and D3 are disposition questions for the owner, not edits this review
may recommend making unilaterally), and no reopening of D6 or D7, which the contract settles.
