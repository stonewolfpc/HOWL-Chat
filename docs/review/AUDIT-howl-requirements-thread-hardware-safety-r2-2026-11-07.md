# AUDIT — Thread and hardware safety requirements, r2 corrective revision

**Date:** 2026-11-07
**Mode:** DOCUMENTATION REVIEW ONLY — read-only; nothing edited, repaired, or reconciled
**Subject:** `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r2-2026-11-07.md`
(SHA-256 `c9bddd10608993ae73d86caa01a7adb6dd83220f98da8dab17e46c1aa3deb8ab`, 447 lines / 36,025 bytes)
**Audited predecessor:** `REQUIREMENTS-howl-thread-and-hardware-safety-2026-11-07.md`
(`ee38f0310042188a6821f8632e9e07bfa60dcac454a71ff2d972bb47439ac0e5`, 344 lines / 25,933 bytes)
**Governing contract:** `DESIGN-howl-tts-helper-architecture-2026-11-06.md` (`ec3c789e…`, 2,586 lines)
**Frozen design inputs referenced by the subject:**
`DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md` (`802ae811…`, 711 lines) and
`CORRECTIVE-howl-tts-stage3-corrective-design-2026-11-07.md` (`8fe9356c…`, 172 lines)
**Owner-decision instrument:** `OWNER-DECISIONS-howl-tts-stage3-2026-11-07.md` (`24d33445…`, 610 lines)
**Reviewer position:** author of the round-1/2/3 corrective-design audits and the focused D10-scope audit; not an author of this requirements set

## 1. Verdict

**Not adoptable as written. Substantively sound on the safety requirements; defective in scope boundaries and in its own adoption mechanism.**

The F3, F4, F5, F6, F7, F9, and F10 corrections are real and verifiable in the text, not merely asserted in the §1 table. F1 and F8 are mostly delivered. F2 is delivered in prose but left openable in practice by one under-specified phrase.

I found:

- **2 blocking findings** — one C8 scope conflict that requires an owner/design resolution rather than a text edit (A-1), and one C1 wording defect that would freeze the wrong interpretation into C1's acceptance evidence (A-5).
- **4 major findings** — three in the adoption mechanism itself (A-2, A-3, A-4) and one internal contradiction about when performance evidence is required (A-6).
- **5 minor findings** — scope narrowing (A-7), an exception/check contradiction (A-8), an unfalsifiable check (A-9), an over-broad N/A rule (A-10), and a circular approval field (A-11).

None of these weakens the safety requirements themselves. All are scoped, and four are one-sentence fixes. But A-1 cannot be fixed by editing this document, because it is a conflict of authority with the frozen design inputs the document says it does not amend.

**Adoption remains pending, and this review does not make it ready.** The owner-decision instrument is untouched: all eight Stage 3 decisions remain `PENDING — OWNER DECISION REQUIRED`.

## 2. Method

- Recomputed SHA-256 and line/byte counts for the subject, the predecessor, the contract, both frozen design inputs, and the owner-decision instrument.
- Read the subject and the predecessor in full.
- Checked every claim in the subject's §1 findings table against the requirement text it cites and against the contract/design rules that constrain it.
- Checked the §8 stage overlay row by row against the corresponding stage's scope, required behavior, exclusions, and acceptance in the frozen design.
- Checked every owner-decision statement against the instrument's header, §3 summary table, and §13.
- Searched the whole project for references to `AC-n`, `TH-1`, `HW-1`, `HW-2`, `PR-n` outside this document, to test whether the renumbering or the overlay collides with an existing citation.

## 3. F1–F10 dispositions

| Finding | Subject's claimed resolution | Independent verdict |
|---|---|---|
| **F1 — authority** | §2 hash-bound adoption mechanism; non-binding until adopted | **Partly delivered.** The self-adoption analysis (§2.2, line 41) is correct and non-obvious: a document cannot carry its own final hash. But the mechanism has three defects (§5: A-2, A-3, A-4) and never states whether the predecessor's recorded requirements stay in force. |
| **F2 — C1 handshake contradiction** | No-spawn scoped to discovery; "contract-authorized handshake start" permitted (§1 line 18; §8 C1 line 347) | **Delivered in prose, openable in practice.** The authorizing trigger is unnamed, so the phrase can be stretched back to the loophole it was meant to close (§5: A-5). |
| **F3 — C6 overlap** | Replacement requires shutdown, kill if needed, observed `Wait()` before any new start; no concurrent residency | **Delivered and exact.** Matches the design's A2-D7 sequence word for word: "shutdown; grace; close stdin; wait; KillTree if needed; observe Wait; start fresh only after exit" (design C6). |
| **F4 — timeout/exit races** | Timeout and unexpected exit are first-class TH-1.2 events and C4/C8 evidence | **Delivered.** TH-1.2 now requires defined behavior for timeout and unexpected exit; TH-1.6 requires timeout-vs-cancel, timeout-vs-exit, cancel-vs-exit, and late-terminal-vs-retired tests. That is the right evidence set, not a gesture at one. |
| **F5 — 120 s idle grace** | Timer/helper may remain through the grace; premature shutdown not required | **Delivered.** Matches contract §5.5 ("An idle instance may be stopped after a bounded idle period (initial policy: 120 seconds)", line 521), the provisional idle-close row (line 1864), and A2-D8 (idle close is a graceful `shutdown`, available even when `unload` is absent). It also removes the predecessor's collision with TH-1.3's orphan rule. |
| **F6 — app-wide I/O evidence** | §9 mandatory per-stage I/O impact record; AC-7 | **Delivered.** The record's fields (path, operation, trigger/cadence, necessity, bound, class, measurement, N/A rationale and reviewer) are reviewable, and §9.3 requires measurement rather than assertion. One wording defect at AC-7 (§5: A-10). |
| **F7 — measurable budgets** | §10 resource-budget record; "bounded" requires value or contract-provisional value, method, threshold | **Delivered.** The provisional escape hatch has real referents in the contract's table at 1858–1866 (handshake 5 s, load/synthesis 120 s, cancel grace 2 s, shutdown 5 s, idle close 120 s, restart 3/10 min, backoff 1/2/5 s), so this is not an invented category. One circularity at the Approval field (§5: A-11). |
| **F8 — exception loophole** | §12 prohibits performance waivers of TH-1, HW-2, cleanup, lifecycle, data integrity, process ordering; HW-1 exceptions need owner decision, scope, expiry, evidence | **Mostly delivered.** The prohibition list is the right one and "no self-approval" is present. But an HW-1 exception can still create a state that fails PR-3, which lives under non-waivable HW-2 (§5: A-8). |
| **F9 — observability surface** | §11 internal diagnostics default; no renderer/public API/status-shape change without authority | **Delivered, and it matters.** §11 (line 414) does not accidentally authorize the pending D5 UI work, and it respects A2-D6 (service capabilities are not added to the provider-scoped status). |
| **F10 — C8 pass/fail** | §6.4 and §8 C8 require a pre-execution plan with hardware profile, workload, duration/iterations, budgets, thresholds | **Delivered as a mechanism, conflicting in scope.** The plan requirement and the refusal to invent thresholds (line 320: "This document does not invent them") are exactly right. But making PR-1–PR-4 C8 acceptance conditions collides with the design's explicit C8 exclusion (§5: A-1). |

## 4. PR-1 – PR-4 assessment

The four checks are well-formed as *measurements*, and the document correctly separates the definition of a check from the numeric threshold that a stage must later approve. That distinction — a requirement being well-defined is not the same as its numeric acceptance threshold having been established — is preserved here, and it is the right call.

Defects are in timing, not in the measurements:

- **PR-1** (normal-operation overhead) is a comparative claim: candidate versus baseline on latency, throughput, CPU, memory, disk, handles, process count. As written it is sound; the "no arbitrary ceiling or unnecessary throttle is active" clause is a policy assertion attached to a measurement, and both halves are reviewable. Scope caveat: the "GPU/accelerator use where applicable" dimension in the C8 row requires a native provider, which the design excludes at C8 (§5: A-1).
- **PR-2** (contention degradation) is sound. Its "safety behavior remains correct" clause overlaps PR-4, which is harmless.
- **PR-3** (adaptive control) is the weakest. It has no N/A rule when a stage implements no adaptive control, so it is either vacuous or unfalsifiable depending on the reviewer, and "actual contention" is not tied to a named signal (§5: A-9).
- **PR-4** (safety regression) is correctly absolute: any safety regression blocks acceptance regardless of PR-1/PR-2. This is the right precedence and should not be softened.

## 5. Findings

### A-1 — BLOCKING (scope, authority): the C8 performance overlay contradicts the frozen design's C8 exclusions

The subject makes PR-1–PR-4 mandatory C8 acceptance evidence:

- §6.4 (line 311): "A performance-regression check is mandatory for C8 and for any stage that changes scheduling, concurrency, throttling, caching, or normal-operation work."
- §8 C8 row (line 354): "Measure CPU, memory, disk I/O, GPU/accelerator use where applicable, handles, and process count at idle and under load; include the declared low-core/low-memory profile; execute PR-1 through PR-4 with approved thresholds."
- AC-8 (line 335): "C8 applies PR-1 through PR-4."

The frozen design says the opposite about C8's scope:

- Design C8, **Exclusions** (line 512): "No native provider, model download, licensing decision, **or performance claim**."
- Design C5, Exclusions (line 368): "No production adapter; no native measurements; numeric values remain provisional, **not performance claims**."
- Design §9, what closure does **not** mean (line 665): "TTS performance or model quality is established."
- Contract §5.5 (line 524): "The cap and idle period are policy values to validate in the implementation spike; **they are not performance claims**"; and the provisional table marks "Load 120 s … no performance claim" (line 1858).

Why this is blocking rather than a wording quibble: the subject claims to "supplement—and not replace—each stage's frozen functional acceptance criteria" (§8 preamble). Adding an acceptance condition that the design explicitly excludes is not a supplement; it is a de facto amendment of C8's scope. And §2.2 (line 47) requires the adoption record to certify "that the governing contract and frozen design inputs were not amended." Those two statements cannot both be true while PR-1–PR-4 gate C8: either the design's C8 scope is amended without authority, or the certification is false, or C8 can never be accepted.

Fair reading, stated so the finding is not overstated: a regression-preservation check is not the same as a claim that synthesis speed or model quality is established. If that is the intended distinction, the document has to say so and has to bound the checks accordingly — PR-1/PR-2 against the application's overhead with a fake helper, no GPU dimension, no throughput claim about synthesis. As written, the C8 row asks for accelerator measurements at a stage whose evidence classes exclude a native provider.

Minimal fix (choose one, and record which):
1. Scope PR-1–PR-4 as a **separately authorized performance-preservation gate** that does not gate C8 contract-closure acceptance; C8's acceptance stays the evidence matrix, and the performance gate is authorized separately with its own owner decision; or
2. Obtain an explicit owner decision amending the design's C8 exclusion, record it in the adoption record's exceptions section, and state that this is an amendment of C8's scope — not a supplement.

What must not happen is a silent third option: adopting the overlay and letting the design's exclusion stand unmentioned.

### A-2 — MAJOR (adoption mechanism): the checklist requires all eight decisions to remain pending

§14, line 443: "[ ] The owner-decision record remains unchanged and all eight decisions remain pending."
§13, line 432: "All eight Stage 3 owner decisions remain pending."

Why it matters: this makes a legitimate future state fail the adoption checklist. If the owner answers D10(a) — the natural next step, and C1's gate — adoption of the requirements document becomes impossible under its own checklist, and an honest reviewer following §14 must reject a valid adoption. It also couples an engineering-requirements adoption to a decision schedule that has nothing to do with it.

Fix: replace the item with a non-coupling statement: "This revision resolves none of the eight decisions; if any has been answered, the adoption record does not cite it as authority and the C1 gates remain as stated." The §13 phrasing ("remain pending") is a status report and can stay; the checklist item cannot.

### A-3 — MAJOR (authority): the predecessor's recorded status is left unresolved

The predecessor's header is not a proposal: "**Status:** RECORDED CROSS-CUTTING ENGINEERING REQUIREMENT — IMPLEMENTATION UNAUTHORIZED", "**Source:** Project-owner instruction recorded in session on 2026-11-07", and §1 states "A stage is not accepted if it does not satisfy the applicable TH-1, HW-1, and HW-2 requirements."

The subject says (line 29): "The predecessor remains the audited revision; this file is the proposed corrective revision"; and (line 35): "This document is **PROPOSED — NOT BINDING**. Until the owner completes a separate adoption record, it may be used as a review checklist but cannot be used to accept or reject a stage."

What is missing is one sentence of precedence. A reader who treats the newest document as the current statement of the requirements can reasonably conclude that no TH/HW acceptance rule applies at all until r2 is adopted — a regression against an owner-recorded instrument, and the exact opposite of what F1 was meant to fix. "Audited revision" describes review lineage, not force.

Fix: "The predecessor's recorded requirements remain in force, unamended, until this revision is adopted. Adoption supersedes the predecessor only where the two texts differ; nothing in this revision suspends the predecessor in the interim."

### A-4 — MAJOR (evidence): the adoption record cannot verify the non-amendment it certifies

§2.2 (line 47) requires "confirmation that the governing contract and frozen design inputs were not amended," and §2.3 includes a "Contract-amendment confirmation" field. But the document records no hashes for those artifacts. It binds only its own path and hash, and even leaves its own hash blank in the block by design.

The gap is not theoretical: the contract is `ec3c789e…`, the corrective proposal `802ae811…` (711 lines / 69,168 bytes, revised 12:05), the correction record `8fe9356c…` (172 lines / 16,190 bytes, revised 12:06), and the owner-decision instrument `24d33445…` — i.e., the "frozen" inputs were revised ninety minutes before this document was written, and one of them includes the D10 wording C1 depends on. Without a recorded hash table, "not amended" is an assurance about whichever revision the signer had in mind.

This is the same failure mode F1 targets, one level up: a hash-bound claim whose reference hashes are absent.

Fix: add a verified-input table in §1 or §2 with the four paths and hashes above, and make the adoption record's confirmation reference that table rather than an unqualified sentence. The owner-decision instrument already does this (its §1 table) and is the model to copy.

### A-5 — BLOCKING for C1 (F2): "contract-authorized handshake start" names no trigger

§1 (line 18) and §8 C1 (line 347) both say: "No spawn solely to discover an optional capability; a contract-authorized handshake start remains permitted."

The design's C1 boundary is narrower than the phrase suggests. Line 123: "no optional operation may start a helper merely to discover it; §17.7.2's lazy-start row limits starts to explicit TTS/catalog/status requests, not `unload` or `cancel`." The contract does authorize starts, but from specific *request* paths — `ProviderStatus` "May handshake/start, but must not load a model solely for status" (line 2266), `ProviderCatalog` "May start/handshake and call `catalog`" (line 2267) — and it separately forbids model load during handshake (line 2075). Nothing in the contract authorizes a start whose purpose is the handshake.

So the phrase can be read as licensing a spawn from a direct `unload` or `Cancel` with unknown capabilities, provided the helper handshakes. That is precisely the behavior C1's structural no-spawn tests assert must not occur, and precisely the ambiguity F2 was raised to remove. As written, the requirement would let an implementer record "authorized-handshake start evidence" for a spawn the design forbids — a wrong interpretation frozen into C1's acceptance evidence.

Fix: "A start authorized by an explicit TTS, catalog, or status request remains permitted (contract §17.7.2; adapter rows 2266–2267). Capability knowledge obtained by that request's handshake is incidental. No start is authorized by the need to complete a handshake or to discover a capability."

### A-6 — MAJOR (internal consistency): three texts give three answers on when PR evidence is required

- §6.4 (line 311): mandatory for C8 **and for any stage** that changes scheduling, concurrency, throttling, caching, or normal-operation work.
- AC-8 (line 335): the stage **records whether** its changes affect normal-operation performance; "C8 applies PR-1 through PR-4."
- §8 overlay table: PR-1–PR-4 appear **only** in C8's row.
- AC-10: "C8 may aggregate evidence but cannot substitute for missing C1–C7 evidence."

Consequence: either every qualifying stage (C2's status caching, C5's timers, C6's replacement, C7's catalog caching) must run the full PR suite — expensive, and in tension with the design's "no native measurements" at C5 — or a stage that changes caching or scheduling is accepted on a recorded yes/no while AC-10 forbids C8 from closing the gap. Both readings are defensible; the document must pick one, because otherwise the first reviewer and the first implementer will disagree and acceptance becomes negotiable.

Fix: "PR-1–PR-4 run at C8 against the integrated build. A stage that changes scheduling, concurrency, throttling, caching, or normal-operation work records an AC-8 impact statement naming the PR dimension its change affects, for C8 to measure. A stage may run a PR subset voluntarily; C8 cannot substitute for a stage's safety evidence (AC-10)."

### A-7 — MINOR (scope): §3.1 narrows the inherited "relies on" scope

Predecessor §2.1 applies to "all code those stages rely on for shared state, process ownership, scheduling, cleanup, or user-visible status."
Subject line 87: "'All code those stages rely on' is limited to code the stage changes or **directly integrates**."

That drops code that is depended on but not directly integrated — for example, a shared settings writer or a cache the new path calls but does not modify. The predecessor's own anti-refactor sentence already prevented scope creep, so the narrowing buys nothing and quietly reduces reach; §3.2 meanwhile keeps "all paths touched or **directly depended on**," so the two scope paragraphs now define the boundary differently.

Fix: restore "relies on … or is directly integrated with" and keep the anti-refactor guard.

### A-8 — MINOR (exception/check contradiction): an HW-1 exception can make PR-3 impossible to pass

§12 (line 421) permits "a HW-1 budget exception … through an explicit owner decision … with scope, duration, safety evidence, and a reason other than speed alone," and excludes only TH-1 and HW-2 from waivers.
PR-3's pass condition (line 317) is "no permanent normal-operation throttle" — a HW-2.4 pass condition, which cannot be waived.

An owner exception to HW-1.8 (for example, permitting a standing concurrency reduction) is a HW-1 exception, so §12 allows it; but PR-3 then cannot pass, and HW-2 cannot be waived to accommodate it. The document can produce an authorized state that its own mandated check must fail.

Fix: add "An exception may not create a condition that fails a HW-2 pass condition; for exception purposes, HW-1.8's prohibition on normal-operation throttling is treated as part of HW-2.1."

### A-9 — MINOR (well-definedness): PR-3 has no N/A rule and no defined contention signal

PR-3 (line 317) requires controls that "activate only after actual contention, are reversible, and deactivate when contention clears," and §5.8 requires adaptive controls to respond to "actual, measured contention." But nothing says what a stage with no adaptive control records, and "actual contention" is never tied to a named indicator.

Result: PR-3 is vacuously passable for a stage with no controls, or unfalsifiable for one that has them, depending on the reviewer. HW-1.7 offers only qualitative "resource pressure."

Fix: add "A stage that implements no adaptive control records PR-3 as N/A with review," and name at least one signal (queue depth, rejection count, measured latency versus the declared profile) plus the detection latency that makes "after actual contention" checkable.

### A-10 — MINOR (over-breadth): AC-7's "untouched paths" is unbounded

AC-7 (line 334) requires "explicit N/A for untouched paths." §9.1 requires the record for "every application path it changes or directly depends on," with N/A only when the stage has no such path. The two rules scope differently: read literally, AC-7 invites an N/A entry for every path in the application, which no reviewer will examine, and which dilutes the entries that matter.

Fix: "explicit N/A for each in-scope path category it did not touch, per §9.1," or reference §9.1 directly.

### A-11 — MINOR (circularity): the budget's Approval field and the unnamed approver

§10 (line 402) defines Approval as "the stage reviewer/owner record that accepted the budget," while AC-4 makes the budget an input to that same acceptance. As written, the approval can only be filled by the gate that depends on it. Separately, §6.4 says thresholds "must be recorded and approved before C8 execution" without naming who approves, and §12 only says the change author may not self-approve.

Fix: name the approver (independent reviewer per §12, or owner for exceptions) and state that budget approval occurs at the stage acceptance gate and is recorded there; for C8 thresholds, name the same approver and require the approval before execution, not at acceptance.

## 6. Answers to the two questions in the brief

**Does r2 contradict the frozen contract?** No direct contract contradiction was found. On the specific points tested, it agrees: idle grace (contract 521, 1864, A2-D8), handshake/no-model-load (2075, 2266–2267), process cap described as "contract-fixed" (the contract calls it a "Frozen v1 hard cap of two helper processes", line 2366, so the wording is accurate), status surface (A2-D6), and provisional values (1858–1866). The one contract-adjacent risk is A-5, where a permissive phrase could be read against the design's C1 boundary; the contract itself does not authorize a handshake-purposed start.

**Does r2 contradict the frozen corrective sequencing?** Yes, in one place: A-1. The design excludes performance claims from C8 (line 512) and native measurements from C5 (368), and states that closure does not establish TTS performance (665); the subject makes PR-1–PR-4 mandatory C8 acceptance evidence. The contradiction is about scope and authority, not about any safety rule, and it cannot be resolved by editing only the requirements document.

Everything else aligns with the design: C6 replacement ordering matches A2-D7 exactly; C5 idle behavior matches; C1's no-frame/no-spawn structural boundary is preserved (the subject's C1 evidence row asks for "structural no-frame/no-spawn checks," not a result code, so it does not unfreeze D10); C7 default-off/no-startup-load matches; C8's evidence-class separation and "cannot be satisfied by fake-helper-only evidence where the claim is about real hardware" match the design's binding evidence rule.

## 7. Owner decisions and authorization

- The subject resolves no decision. §13's eight-item list matches the owner-decision instrument's status (`OPEN — ALL EIGHT DECISIONS PENDING`) and its §3 table entry for entry: D10(a), D10(b), D1a, D9, D1b, D2, D5, D3.
- The statement that "a conversational preference for a D10(a) result is not a recorded owner decision" agrees with the instrument's own rule that a blank field is not consent, and with the current instrument, which records no selection.
- The subject's C1 evidence row does not assert a D10 result, so it does not narrow or pre-empt either branch.
- No stage authorization leaked: the document states repeatedly that adoption, review, and stage acceptance authorize nothing; C1 and C2 remain unauthorized; C8 still needs separate runtime-execution authorization. §3.3's non-goals list is complete.
- The predecessor is preserved: its hash matches the subject's header claim, and the contract is unchanged at `ec3c789e…`.

## 8. Observations (not defects)

- **O-1:** F1–F10 exist only in this document's §1 table. No audit of the predecessor is recorded anywhere in the project, so the findings themselves cannot be checked independently — only the resolutions, against the text. This also weakens §2.2's requirement that the adoption record carry "the independent Mimo review record": there is no finding record for this requirements line, only this review. Recording the findings (or citing the audit that produced them) would close the gap.
- **O-2:** §1's F5 phrasing "premature shutdown is not required" is weaker than the operative rule. Contract 521 allows stopping "after a bounded idle period," and §8's C5 row correctly says idle close "occurs after the contract-defined 120 s grace"; a helper closed before the grace is not merely unnecessary, it is outside the contract's idle-close trigger. The §1 summary could mislead a reader who does not reach §8.
- **O-3:** The AC-7 → AC-10 renumbering (I/O evidence inserted at AC-7, no-substitution moved to AC-10) is citation-safe: no other artifact references `AC-n`. Verified across the project.
- **O-4:** The design and correction record were revised at 12:05/12:06, after the 11:24 revision I last audited. My n1 (owner-decision count) and n2 ("deferring the branch" is not a result) from that pass are now fixed in the design: line 36 reads "five owner decisions (D1b, D2, D5, D9, D10)," and C1's prerequisite/acceptance/matrix now require "a concrete caller-facing result." The subject binds to `802ae811…`/`8fe9356c…`, which match the owner-decision instrument's baseline — but it does not record those hashes (A-4).
- **O-5:** `README.md` still ends at §28. Neither requirements document, nor the focused D10-scope audit, is indexed. Maintenance gap, not a defect in the subject.

## 9. Limitations

1. **Not independent of the design line.** I authored the audits that produced the corrective design's current wording. This review is independent of the requirements author and of the requirements text, but it is not a third-party check of the program as a whole.
2. **F1–F10 unverifiable as findings.** See O-1. I verified each claimed resolution against the requirement text and the governing rules; I could not verify that the findings were accurately characterized in the first place.
3. **No execution.** Documentation review only; no build, test, benchmark, or process launch. Statements about measurement feasibility (including PR-1/PR-2 on a fake helper) are judgments from the design's evidence rules, not measurements.
4. **Terminology note.** "Frozen design inputs" is the owner-decision instrument's own term for the proposal and correction record, so the subject's use of it is consistent with established nomenclature even though those artifacts are proposals. A-1 is about substance (C8 scope), not about that word.
5. **This review does not adopt, reject, or correct anything.** It creates no adoption record and touches no subject artifact.

## 10. Recommended path

1. Fix A-5 and A-6 in the requirements text (one clause each, no authority needed).
2. Fix A-2, A-3, A-4 in the adoption mechanism (one sentence, one table).
3. Take A-1 to the owner as a scope question, not a wording question: either the PR gate is separate from C8 acceptance, or C8's exclusion is formally amended. Record the choice.
4. Then re-review the changed sections only, and only then consider adoption.

The safety requirements themselves do not need to change. That is the strongest thing that can be said for this revision: F3, F4, F5, F7, F9 and F10 are real repairs, and the document's instinct to keep numeric thresholds out of the requirement text is correct. The defects are in how far the overlay reaches and in whether its own adoption gate can be satisfied honestly.
