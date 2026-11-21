# FC-6 — Evidence and Acceptance Contract

## Contract paper `FC-6-2026-11-21-01`

| Field | Value |
|---|---|
| Date | **2026-11-21** (local, CST) |
| Type | Contract paper — review draft |
| Authority | `HF-CONTRACT-INTAKE-2026-11-20-01`; D7 acceptance requirement; `HRA-DB-2026-11-21-01` |
| Gate | All gates |
| Status | **PROPOSED — not approved** |
| Implementation authorization | **None** |
| Tests/models/process runs | **0** |

**Labels:** **Owner-approved** · **Proposed** · **Unresolved** · **Observed** · **Unestablished**.

---

## 1. Contract decision

**FC6-01 (Proposed).** Acceptance is evidence at a stated level, not a code diff, screenshot, demo, log line or historical report. A claim names what was tested, on which build, with which fixture, and what it does not establish.

**FC6-02 (Owner-approved — HRA-34).** “Ready” and “complete” are evidence claims. The same applies to persistence, state transitions, modality readiness and memory writes.

**FC6-03 (Calculated).** Every contract paper in this set exposes an acceptance surface. FC-6 defines the common receipt, claim ladder and fault matrix; it does not replace each paper's specific criteria.

---

## 2. Claim ladder

**FC6-04 (Proposed — required levels).**

| Level | Meaning | Permitted claim |
|---|---|---|
| **Requested** | Intent recorded | “The user asked” |
| **Admitted** | Policy/capacity allowed an attempt | “Howl decided to try” |
| **Attempted** | Operation dispatched | “The operation started” |
| **Observed** | A result or state change was seen | “Something appears to have happened” |
| **Verified** | A workload-appropriate check confirms the outcome | “This outcome is established for this scenario” |
| **Unknown** | Evidence is missing, insufficient or contradictory; the outcome is not established | “We do not know” |
| **Failed** | Evidence establishes that the operation did not succeed | “This failed” |

**FC6-05 (Proposed — no promotion by adjacency).** A process start is not readiness. An unload return is not release. A UI update is not a durable commit. A save call is not persistence. A log line is not proof.

**FC6-06 (Proposed — evidence provenance).** Every consequential claim records source, timestamp, build identity, fixture identity, confidence and scope. Measured, estimated and unknown remain distinguishable.

---

## 3. Receipt schema

**FC6-07 (Proposed — minimum receipt fields).**

| Field | Purpose |
|---|---|
| `receipt_id` | Stable evidence identity |
| `claim` | What is being asserted |
| `level` | Requested/admitted/attempted/observed/verified/unknown/failed |
| `contract` | FC paper and criterion ID |
| `build_identity` | Commit/build/artifact identity |
| `environment` | OS, hardware, runtime, model/tokenizer where relevant |
| `fixture` | Fixture ID and hash |
| `action` | What was done |
| `expected` | Pre-stated pass condition |
| `observed` | Actual result |
| `evidence_refs` | Logs, screenshots, hashes, records |
| `limitations` | What the evidence does not establish |
| `timestamp` | When observed |
| `producer` | Who/what produced the evidence |

**FC6-08 (Proposed).** A receipt without a pre-stated pass condition is not acceptance evidence. A receipt that cannot name its build/fixture is not reproducible.

---

## 4. Fault matrix

**FC6-09 (Proposed — required fault classes).**

| Class | Examples | Required response |
|---|---|---|
| Duplicate/retry | Repeated append, repeated request | Idempotent or explicit conflict |
| Partial write | Torn event, interrupted save | Recover or report unknown |
| Crash | Process death during transition | Reconcile; no phantom state |
| Missing dependency | Tokenizer/model/voice/profile absent | Explicit unavailable |
| Identity ambiguity | Name collision, missing ID | Refuse; no name fallback |
| Privacy | Cross-character context, snapshot exposure | Fail closed; redact evidence |
| Resource | OOM, telemetry loss, thermal, driver loss | Conservative policy; preserve work |
| Cancellation | User interrupts text/audio | Preserve durable text; separate audio state |
| Concurrency | Two instances, two requests | No duplicate promise within the authoritative reservation ledger and its atomicity domain; cross-instance safety follows the active D6 policy; a per-instance ledger establishes no global guarantee |
| Restoration | Missing revision, changed template | `degraded`/`unresolved`; no latest fallback |

**FC6-10 (Proposed).** Fault injection must name the exact failure point. A generic “error path tested” claim is insufficient.

---

## 5. D7 end-to-end acceptance

**FC6-11 (Owner-approved — D7 requirement).** TTS is not complete until the actual multi-character experience works. Two voice selectors in the UI do not prove the requirement.

**FC6-12 (Proposed — required scenario `D7-E2E-01`).**

| Step | Requirement |
|---|---|
| Setup | Two characters with distinct, locally generated voice profiles; multi-character conversation; no cloud TTS endpoint configured |
| 1 | Character A produces a committed assistant turn; TTS uses A's assigned voice |
| 2 | Speaker changes to Character B; hard flush occurs; B's response uses B's distinct voice |
| 3 | Speaker returns to A; A's voice assignment is unchanged |
| 4 | Interrupt A's TTS playback mid-utterance; A's committed text remains durable; audio is explicitly cancelled |
| 5 | Restart the app and restore the conversation; both character IDs and voice assignments resolve |
| 6 | Replay/continue the conversation; each character still uses its assigned voice |
| 7 | Change one character's voice setting; only that character changes |
| 8 | Remove/disable one voice profile; that character reports unavailable; no silent substitution |

**FC6-13 (Proposed — required evidence).** Voice profile IDs per character, audio artifacts or engine logs, assignment table before/after, switch/restart/restore timestamps, interruption record, no-cloud-request evidence, and the exact build/fixture identity. A screenshot of selectors is not sufficient.

**FC6-14 (Proposed — failure semantics).** TTS failure does not roll back the committed text turn. Text cancellation does not leave an orphan audio job. A missing voice does not fall back to another character's voice or a cloud voice.

**FC6-15 (Unresolved).** The exact audio-artifact retention, playback ownership and privacy policy for generated voice output require an owner decision before this test is finalized.

---

## 6. Gate evidence

**FC6-16 (Proposed).**

| Gate | Minimum evidence |
|---|---|
| M0A | A1–A9 from `HCDM-M0A-AR-2026-11-21-01`, separately authorized |
| M1 | Replay/restart/branch/duplicate/deletion/voice-persistence evidence |
| M2 | Real tokenizer/template plan, fail-closed over-budget, reproducible plan hash |
| M3 | Two-character private-context isolation, hard flush, distinct voices |
| M4 | Logical reservation, admission, drain/release/readiness, unknown handling |
| M0B/M5 | Canonical revision binding, browser import/quarantine, no name fallback |
| M8 | Delegated D8 decision record plus lifecycle-fact evidence |

**FC6-17 (Calculated).** M0A is separate. Nothing in FC-1–FC-6 authorizes M0A or closes its criteria. M0A uses its own acceptance review and owner authorization.

---

## 7. No-claim rules

**FC6-18 (Proposed).**

- A historical test, worklog or DeepSeek report is not a current result.
- A passing unit test is not an end-to-end result.
- A screenshot is not persistence evidence.
- A model reply is not proof of context isolation.
- A process start is not readiness.
- An unload return is not release.
- A UI selector is not voice-assignment persistence.
- A code diff is not acceptance.
- An unavailable check is reported as unavailable, never as passed.

---

## 8. Acceptance design for FC-6 itself

**FC6-20 (Proposed — reproducible receipt).** A receipt can be reproduced from build + fixture + action.

**FC6-21 (Proposed — failure is not pass).** A failed or unavailable check is visibly not passed.

**FC6-22 (Proposed — gate evidence level).** Every gate names its evidence level.

**FC6-23 (Proposed — D7 scenario not run here).** D7-E2E-01 is the required evidence for distinct voices, speaker switch, interruption and restoration; this paper does not run it.

**FC6-24 (Proposed — T0.5 separation).** T0.5 evidence is separate from D6 admission evidence.

**FC6-25 (Proposed — D4 evidence identity).** D4 evidence names the actual tokenizer/template, not a library choice.

**FC6-26 (Proposed — D8 evidence identity).** D8 evidence is the delegated decision record, not a unilateral paper.

---

## 9. Traceability and non-authorization

| Link | Artifact |
|---|---|
| Contract register | `README.md` (this set) |
| M0A acceptance | `HCDM-M0A-AR-2026-11-21-01` |
| Resource evidence | `HRA-DB-2026-11-21-01` |
| Phase 1 acceptance | `HF1-PLAN-2026-11-21-01` |
| D7 requirement | `HF-CONTRACT-INTAKE-2026-11-20-01` |

**FC6-19 (Observed — no authorization).** No test, fault injection, model run, TTS run or acceptance claim was performed. This paper defines evidence requirements only.

**End of FC-6.** Acceptance is a reproducible receipt at a named claim level. The two-character TTS scenario is part of the contract, not a UI demonstration.
