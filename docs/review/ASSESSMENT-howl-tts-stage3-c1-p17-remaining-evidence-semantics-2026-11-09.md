# C1 — P17 Remaining Evidence-Semantics Assessment (R09–R13, U1–U8)

**Date:** 2026-11-09
**Status:** RECONCILED READ-ONLY ASSESSMENT. **P17 REV2 NOT APPROVED. NO REV4, NO P18, NO CONTRACT AMENDMENT. C1 IMPLEMENTATION UNAUTHORIZED.**
**Subject:** the remaining evidence-semantics questions in P17 and the independent Mimo lifecycle findings.
**Baseline:** current P17 rev2 `651d38d3…` (538 lines); frozen original P17 `01d887c1…`; corrected r3 manifest response `60ca1cb5…`; lifecycle decision `aa136bc9…`; contract rev3 `ed0da1a7…`; Mimo review `d67a82da…` (547 lines); U1–U8 findings `33c776bf…` (640 lines).
**Method:** static source reading, hash verification, and line-by-line comparison of the Mimo review, U1–U8 findings, P17 rev2, contract A6, and the pinned helper source. No test, build, benchmark, or application execution was performed.

This document assesses what remains unestablished. It does not revise P17, propose a replacement schema, amend the lifecycle contract, or authorize implementation. Where the source cannot establish a proposition, the proposition is left unresolved rather than assigned an outcome.

## 1. Baseline identity and verification

### 1.1 Artifact hashes

| Artifact | SHA-256 | Observation |
|---|---|---|
| `PROPOSAL-howl-tts-stage3-c1-evidence-semantics-p17-2026-11-09.md` | `01d887c18700d63070f1a2649659281ec1e875a21b49375adea12f8f11d917bf` | Frozen P17 proposal, unchanged |
| `RESPONSE-howl-tts-stage3-c1-p17-review-r3-2026-11-09.md` | `60ca1cb563556d12163a644741348319de1f539aab96d9e2235b88013b234ec0` | Corrected r3 response, frozen after terminology cleanup |
| `DECISION-howl-tts-stage3-c1-lifecycle-disposition-2026-11-09.md` | `aa136bc9c3d5eec77e52be58662447f671a7350185bf045d6da3a1bd67758b7d` | Owner decision, unchanged |
| `CONTRACT-howl-tts-stage3-c1-installation-lifecycle-rev3-2026-11-09.md` | `ed0da1a75d4f62874f873979e6667f66f19541a7523ff96369f2b44e8391195e` | Contract rev3, frozen and unaccepted |
| `PROPOSAL-howl-tts-stage3-c1-evidence-semantics-p17-rev2-2026-11-09.md` | `651d38d3e2d97f981a606a39ea4a3117c2e748ff9f5941eb8bd9671b6e022f82` | Current P17 candidate, 538 lines; unapproved |
| `ASSESSMENT-howl-tts-stage3-c1-p17-u1-u8-scenario-findings-2026-11-09.md` | `33c776bf50b6a10467c319d61b7d0079815705b33e050991140d5f5611b71e95` | Authoritative U1–U8 working findings, 640 lines; all cases remain open |
| `REVIEW-howl-tts-stage3-c1-p17-independent-lifecycle-adversarial-mimo-2026-11-09.md` | `d67a82da2d0c1e8d3d2f35c70a4f40058e6267572cb1eca827dee5ffee155314` | Independent Mimo review, 547 lines; read-only |
| `DECISION-REGISTER-howl-tts-stage3-c1-p17-evidence-2026-11-09.md` | `3c64797bfdd48ddacca95e6a0def88d960e4104aa3fc3034a25dc539d97cbc99` | OD-1–OD-9 accepted; P17 approval gate remains unmet |
| `ATTESTATION-howl-tts-stage3-c1-independent-review-2026-11-08.md` | `598d40c53aed4ad73006277fe7711e2e54d21d7c92a65243e95b119911c88ca6` | Unsigned, unchanged |
| `RECONCILE-howl-tts-stage3-c1-p17-review-2026-11-09.md` | `b634cd20d99367cb…` | First reconciliation, unchanged |
| `REVIEW-howl-tts-stage3-c1-p17-independent-2026-11-09.md` | `c811e80eb577b010…` | Independent review r1, unchanged |
| `REVIEW-howl-tts-stage3-c1-p17-reconciliation-2026-11-09.md` | `cd8ddedd6f5bd5f6…` | Review r2, unchanged |

The r3 response's condition-pair terminology correction is confirmed at the declared-scope documentation level. This assessment does not reopen it.

### 1.2 Source identity

| Item | Value |
|---|---|
| Repository HEAD | `9c7d2a17dd7f702d7db6d95d9316e7640e34bf6f` (`9c7d2a1`) |
| Helper tracked by Git | **0 files**; `client.go` is not a Git path |
| Helper Go files | 19 |
| Helper Go lines | 4,992 |
| Helper content digest | `d0df8b2823389b305698fc06eda28e0e01f06259fdc69174cf7fb3851e271a51` |
| Working tree at assessment start | 142 entries, pre-existing; no helper Go file modified in this assessment |

The HEAD pin does **not** establish the helper's identity, because the helper directory is untracked. The content digest is the only source-identity claim made here. It binds relative paths and contents for `*.go`; it does not bind non-Go files, runtime state, or platform behavior.

Reproduction:

```bash
find internal/audio/tts/helper -type f -name '*.go' -print0 \
  | LC_ALL=C sort -z \
  | xargs -0 sha256sum \
  | sha256sum
# d0df8b2823389b305698fc06eda28e0e01f06259fdc69174cf7fb3851e271a51
```

### 1.3 What was not done

- No frozen artifact was edited.
- No Go source, test, configuration, or frontend file was modified.
- No test, build, benchmark, or application execution was run.
- No contract amendment, rev4, P18, manifest regeneration, or implementation was produced.
- No outcome was assigned to a case merely because the proposal asserts one.

## 2. R09–R13 individual assessment

### 2.0 Provenance limit and reading rule

The independent review of record identifies R09–R13 thematically. It does **not** provide one verbatim paragraph per identifier. Its exact available text is:

> R09–R13 are different: they challenge whether P17's replacement model is sufficiently specified.
>
> The most consequential is R09. An append-only log is not automatically a reliable evidence model. Without instance identity, sequence-allocation rules, causal linkage, action results, phase context, and retention rules, it can preserve entries while still producing an ambiguous or misleading history.
>
> R10 and R11 expose two more important gaps:
>
> - Precedence is currently a selection rule, not a demonstrated resolution of competing evidence.
> - `quarantined` is one of the five proposed projections, but P17 does not define a path that can produce it.
>
> R12 and R13 prevent a premature claim of compatibility with the existing application. Handshake classification and provider-state behavior still have unresolved semantics.

The per-ID headings in the reconciliation (`R09 — An append-only log is underspecified`, etc.) are the author's expansion, not original review prose. They are cited as such below. No per-ID review wording is invented.

### 2.1 Summary matrix

| ID | Classification | Source status | P17 treatment | Blocks approval | Minimum next action |
|---|---|---|---|---|---|
| R09 | Proposal/design omission | Source-verified absence of the required record fields; E1–E9 are source-verified losses | §3.1 record shape; §2 constraints; §8 P1/P3/P4/P9; allocator, retention, and linkage remain OPEN | **Yes** | Define the minimal record/allocator/retention and map E1–E9 |
| R10 | Proposal/design omission | Source-verified absence of ordering and causal linkage | §8 P3 rejects projection from a fatal observation alone; P4/P4-A route confirmed exits; P6 requires OD-8 A–E | **Yes** | Define a causality standard or leave combinations unresolved; no class-order winner |
| R11 | Design omission plus accepted owner decision | Source-verified: no `quarantined` value or transition | §8 P7/P10 makes it internal-only with no public projection; mechanics remain OPEN | **No**, not as a P17 text contradiction; blocks C1 mechanics | Define the internal quarantine record/transition and release conditions |
| R12 | Owner decision and design omission | Source-verified two projections for the same compatibility concept | O15/O16; P2 classifies missing mandatory capability as `incompatible` but leaves path-by-path classification OPEN | **Yes**, for the unresolved categories | Apply P2 to the source path and classify or explicitly leave the remaining exits open |
| R13 | Owner decision and design omission | Source-verified writes that conflate provider condition, process state, and severity | O18/O19; S9/P3 forbid `crashed` for a live process; source publishes it at `client.go:756` | **Yes**, for U6 | Remove the live-process write and define the no-exit projection |

### 2.2 R09 — An append-only log is underspecified

**1. Exact available finding and provenance.** The independent review states that an append-only log is not automatically reliable without "instance identity, sequence-allocation rules, causal linkage, action results, phase context, and retention rules" (`REVIEW-howl-tts-stage3-c1-p17-independent-2026-11-09.md`, §2). The reconciliation expands this into its R09 heading and says the decision required is to "define the record schema and allocator, or state that the log is diagnostic only and cannot support projection" (`RECONCILE…`, §4). The latter is author text, not review text.

**2. Current P17 treatment.** P17 rev2 §3.1 proposes `Observation`, `Action`, and `ExitFact` records with generation/session binding and separate observation/event order; §2 constraints 2–5 and §8 P1/P3/P4/P9 make projection depend on retained, generation-bound evidence. The exact schema, allocator, retention, and storage remain **OPEN**. P17 does not define what allocates `seq`, how concurrent writers are ordered, how an observation is linked to an action or exit, what `completion` contains, or what happens when the log is bounded. Naming fields is not the same as defining their semantics.

**3. Source-verified behavior.**

- `Failure` carries only `Code`, `Message`, `Retryable`, `Field`, and `Detail` (`errors.go:15–21`). It has no instance, sequence, timestamp, phase, site, severity, action, or causal field.
- The client's `next` counter is a request counter, not an evidence sequence. It is incremented under `c.mu` and used to build `req-<session>-<counter>` (`client.go:135–139`, `:457–458`; `protocol.go:267–268`). It orders request IDs, not observations.
- `session` is a random client session (`protocol.go:259–265`) and is rotated when a crashed client is restarted (`client.go:206–223`). It is not a process-instance identifier, and it is not carried by `Failure`.
- The supervisor stores one `lastFailure` and one `waitErr` (`supervisor.go:166–195`, `:233–242`). `SetState` assigns `lastFailure = failure` on every publication (`:238`); `setUnavailable` does the same (`:201`); the unexpected-exit branch overwrites it with `helper_crashed` (`:190`). No prior failure is retained.
- `failAll` stores only the first read failure (`client.go:373–376`), clears the pending map (`:381`), and delivers that same failure to every waiter (`:383–385`). Later failures are not recorded.
- A terminal for a retired request is dropped without a record: the branch returns `true` after the comment "A late terminal is evidence only" but stores nothing (`client.go:339–342`). The `retired` map grows for the client's lifetime (`:464–468`).
- `Terminate` has no `intent`, `grace`, start, completion, or result record (`supervisor.go:246–288`). Its result is ignored at `client.go:248`, `:654`, `:661`, `:664`, and `:861`; the cancellation result is ignored at `:871`.
- `waitForReader` has no return value and returns silently on either join or timeout (`client.go:349–362`); `Start` discards that silence (`:249`).

**4. Classification.** Proposal/design omission, with source-verified absence of every field P17 needs. The source also contains the evidence losses E1–E9, so this is not only a documentation gap: the proposed log cannot recover facts that upstream paths already discard.

**5. What precise evidence or definition is missing.** The definition of the instance whose history is being recorded; the allocator and ordering domain for `seq`; the rule linking an observation to an action and an exit; the content and failure modes of `Action.completion`; the representation of phase; the retention/drop policy; and the authority of the log (diagnostic only or projection-supporting).

**6. Does it block P17 approval?** **Yes.** P17 §8 P3/P4/P9 makes public projection depend on retained, generation-bound evidence, and §9 says E1–E9 remain unresolved. A model that cannot order or link entries cannot support the publication rules.

**7. Minimum next action.** OD-1 has already accepted projection-supporting authority. Define the minimum record, allocator, and retention rule, and map every E1–E9 loss. This assessment does not choose the schema.

**8. What must not be inferred.** That append-only means complete; that request-ID order is evidence order; that `lastFailure` is a history; that `session` is a process instance; that a digest is a Git identity; that an unrecorded event did not occur.

### 2.3 R10 — Precedence is selection, not resolution

**1. Exact available finding and provenance.** The independent review says: "Precedence is currently a selection rule, not a demonstrated resolution of competing evidence" (`REVIEW…independent`, §2). The reconciliation's R10 heading says the missing definition is "an evidence rule that distinguishes correlation from causation" and that class order is a policy choice (`RECONCILE…`, §4). The heading is author text.

**2. Current P17 treatment.** P17 rev2 does not order evidence classes or supply a winning table. §8 P3 rejects projection from a fatal observation alone; P4/P4-A route confirmed exits, with P4-A still **PROPOSED/OPEN**; P6 requires OD-8's A–E standard for any causal attribution. P17 therefore does not define what evidence establishes that a fatal observation caused an exit, how to order observations from different writers, or what happens when an action and an exit both explain the outcome.

**3. Source-verified behavior.**

- `protocolViolation` calls `failAll(f)` and then `Kill` (`client.go:364–367`). `Kill` is `Terminate` with the default grace (`supervisor.go:291–293`). `Terminate` sets `providerState = stopping` before waiting (`:256–259`).
- The exit observer decides the outcome solely from whether `providerState == stopping` at the moment it acquires `s.mu` (`supervisor.go:176–182`). It does not consult `readFailure`, a sequence, a timestamp, or the action that set `stopping`.
- A fatal protocol observation is stored in `readFailure` (`client.go:374–375`) but is never published by `failAll`; `failAll` only delivers to pending requests. The supervisor's state is changed by `Terminate`/`observeExit`/`SetState`, not by `readFailure`.
- If the observer sees `stopping`, a zero exit becomes `absent` and a nonzero exit also becomes `absent` with `processState = exited` (`supervisor.go:176–182`). If it does not see `stopping`, the exit becomes `crashed` (`:183–190`). The same fatal protocol observation can therefore sit behind either projection, depending on interleaving.
- `SetState` overwrites `providerState` and `lastFailure` without consulting prior fatal evidence (`:233–242`).

**4. Classification.** Proposal/design omission, supported by source-verified absence of causal or ordering data. The source behavior is established; the claim that one class "wins" is not.

**5. What precise evidence or definition is missing.** A causal standard: what independent evidence shows that an action caused an exit, as opposed to merely preceding it; how observations from concurrent writers are ordered; how an exit fact is linked to the action that requested it; and which combinations remain unresolved rather than assigned a winner.

**6. Does it block P17 approval?** **Yes.** Every U1–U5 case depends on this rule. A rule that selects a winner without A–E evidence would encode causality the source does not record.

**7. Minimum next action.** Define the causal standard, or explicitly leave the combination unresolved. Do not add a precedence level to manufacture certainty. The reconciliation already declines to propose one; that restraint should be preserved.

**8. What must not be inferred.** That a fatal observation before an exit proves it caused the exit; that a kill explains an exit without an action record; that class order is chronology; that the first stored failure is the first event.

### 2.4 R11 — `quarantined` is unreachable

**1. Exact available finding and provenance.** The independent review says: "`quarantined` is one of the five proposed projections, but P17 does not define a path that can produce it" (`REVIEW…independent`, §2). The reconciliation's R11 heading identifies the natural path as a terminate timeout and notes that P17 neither projects nor records the state (`RECONCILE…`, §4). The heading is author text.

**2. Current P17 treatment.** P17 rev2 §8 P7 and P10 make `quarantined` internal-only, with no public projection and no enum change; the old five-projection claim is gone. P17 still does not define the internal representation, transition, or release conditions. The remaining defect is implementation mechanics, not a public-claim contradiction.

**3. Source-verified behavior.**

- The provider-state set has thirteen values and no `quarantined` (`internal/audio/types/provider.go:20–36`).
- The process-state set has five values — `absent`, `starting`, `alive`, `stopping`, `exited` — and no `quarantined` (`provider.go:8–16`).
- `Terminate` can return a `timeout` failure after a second grace period while the process may still be owned (`supervisor.go:277–286`). At that point `s.proc` may still be non-nil and `ProcessAlive` returns true if the process state is not `exited` (`:334–342`).
- `Shutdown` clears `started` and `transport` after `Terminate` regardless of its error (`client.go:668–674`). The supervisor still owns the process, but the client reports not started.
- No supervisor field records quarantine, ownership uncertainty, or a terminal-evidence latch. `SetState` can overwrite `stopping` (`supervisor.go:233–242`).

**4. Classification.** Design omission plus an accepted owner decision. OD-7 makes `quarantined` internal-only and not a provider state, but it does not define the record or transition P17 would need. The source confirms there is no such value or path today.

**5. What precise evidence or definition is missing.** The internal disposition record and transition; how the exact generation and ownership/reaping uncertainty are represented; how a timeout differs from quarantine; and when quarantine is released. The public status is already settled by OD-7.

**6. Does it block P17 approval?** **Not as a P17 text contradiction after rev2.** P17 no longer claims a public projection. It blocks C1 implementation because no internal representation or transition exists, and U2/U6/U7 remain open.

**7. Minimum next action.** Define the internal quarantine record, transition, and release conditions. No provider-state enum change is authorized.

**8. What must not be inferred.** That `quarantined` is a synonym for `unavailable`; that a terminate timeout proves the process is dead; that supervisor ownership is a provider state; that a later exit clears quarantine without evidence.

### 2.5 R12 — Handshake classification remains unresolved

**1. Exact available finding and provenance.** The independent review says: "Handshake classification and provider-state behavior still have unresolved semantics" (`REVIEW…independent`, §2). The reconciliation's R12 heading identifies nine `validateHello` exits, all `protocol_violation`, mapped to `unavailable`, while a helper-declared incompatibility maps to `incompatible` (`RECONCILE…`, §4). The heading is author text; the source behavior below is verified independently.

**2. Current P17 treatment.** P17 rev2 §8 P2 (accepted OD-6) classifies a missing mandatory service capability as `incompatible`, and gives categories for several other handshake failures. O15/O16 keep path-by-path classification OPEN. P17 no longer preserves the old code-name split as policy, but it does not classify all nine `validateHello` exits.

**3. Source-verified behavior.**

- `validateHello` has nine failure returns: helper version (`protocol.go:503–505`), limits (`:506–508`), provider descriptor (`:511–515`), duplicate provider ID (`:516–519`), provider capabilities (delegated, `:520–522`), service capabilities (delegated, `:524–526`), missing mandatory service capabilities (`:527–531`), extension descriptor (`:534–536`), and duplicate extension (`:537–540`). Every direct return is `protocol_violation`; the delegated capability validators also return `protocol_violation`.
- `Start` maps only `protocol_incompatible` and `protocol_minor_unsupported` to `TTSStateIncompatible`; every other handshake failure maps to `TTSStateUnavailable` (`client.go:250–254`).
- A helper-declared error reaches the client through `failureFromWire`, which preserves the wire code (`client.go:434–435`, `:471–481`). Thus a helper that explicitly says `protocol_incompatible` is projected `incompatible`, while a helper that omits a mandatory capability is projected `unavailable`.
- The lifecycle decision says "handshake fatal → `incompatible` (existing state is legitimate)" (`DECISION…`, §5, P15). That sentence is ambiguous between all handshake fatals and the narrower incompatibility case. It does not supply a classification rule.

**4. Classification.** Source-verified behavior with an accepted policy boundary. OD-6/P2 settle missing mandatory capability at the category level; the remaining `validateHello` paths and the source's code mapping are not fully classified.

**5. What precise evidence or definition is missing.** The path-by-path classification for the nine `validateHello` exits, including duplicate identities, invalid descriptors, malformed version/limits, extension and optional-requiredness cases; and the mapping from wire code/condition to the P2 categories.

**6. Does it block P17 approval?** **Yes, but narrower than originally stated.** P2 resolves the missing-mandatory-capability case; the source still maps it to `unavailable`, and O16 leaves several categories unclassified. U8 remains OPEN for those categories.

**7. Minimum next action.** Apply P2 to the missing-capability path and classify or explicitly leave the remaining `validateHello` exits open. No source change is authorized here.

**8. What must not be inferred.** That a code name states the helper's intent; that a missing capability is necessarily `incompatible`; that helper-declared and client-detected incompatibility are the same fact; that `unavailable` is a safe default without a decision.

### 2.6 R13 — Provider-state compatibility is unresolved

**1. Exact available finding and provenance.** The independent review says: "R12 and R13 prevent a premature claim of compatibility with the existing application. Handshake classification and provider-state behavior still have unresolved semantics" (`REVIEW…independent`, §2). The reconciliation's R13 heading identifies the `Synthesize` write-failure `crashed` publication at `client.go:756` and the fatal-pending `degraded` publication at `:796`, and says the decision record separates provider state from disposition (`RECONCILE…`, §4). The heading is author text; the writes are source-verified.

**2. Current P17 treatment.** P17 rev2 S9/O18/P3 says a live process must not be classified `crashed`, and a fatal transport observation does not establish an exit, a terminal, or a public state. P17 no longer projects fatal transport to `crashed + helper_crashed`; O18 requires preserving the last justified state and leaves the no-justified-state case U6/C-01 OPEN.

**3. Source-verified behavior.**

- On `transport.Send` failure, `Synthesize` immediately calls `SetState(TTSStateCrashed, "", newFailure("helper_crashed", ...))` (`client.go:755–758`). The process may still be alive at that point.
- `cleanupRejected` then checks `sup.ProcessAlive()` and may call `Terminate` (`client.go:859–864`). That check uses the supervisor's process state (`supervisor.go:334–342`), while `SetState` sets `processState = alive` whenever `s.proc != nil` and the new state is not `stopping` (`:239–240`). The publication therefore does not itself prove the process exited.
- `Status` derives `ProcessAlive` from both dimensions (`supervisor.go:324`); the helper at `:351–356` returns false for `crashed` and `unavailable`, so a live OS process can be reported as not alive while `crashed` is the provider state.
- On a fatal pending outcome, `Synthesize` publishes `TTSStateDegraded` (`client.go:790–797`). The same path handles a `protocol_violation` delivered by `failAll`, so a fatal protocol observation can emerge as `degraded`, not `crashed`.
- Wire errors, result decode errors, result validation errors, and WAV validation errors all publish `degraded` (`client.go:799–829`). `SetState` overwrites `lastFailure` each time (`supervisor.go:233–242`).
- No latch prevents a later `observeExit` or `Terminate` from changing a `crashed` or `degraded` state (`supervisor.go:176–190`, `:256–259`).

**4. Classification.** Source-verified writes with an accepted semantic boundary. OD-5 is now accepted: `crashed` means an observed exit without established orderly shutdown, and a live process must not be `crashed`. The source is nonconforming, and the replacement projection for a send failure with no exit remains OPEN.

**5. What precise evidence or definition is missing.** The post-failure current-state projection when no exit is established; a typed cleanup outcome and actual liveness; whether a live process may retain an active or `degraded` state; and how U6 resolves without a public `crashed` write.

**6. Does it block P17 approval?** **Yes.** P17 rev2 forbids `crashed` for a live process (S9/O18/P3), while the source publishes it at `client.go:756` and can later publish `degraded` or `ready`. U6 remains OPEN, and the model and writers do not agree.

**7. Minimum next action.** Remove the live-process `crashed` write and resolve the post-failure/no-exit projection under OD-5/P3. No source change is authorized here.

**8. What must not be inferred.** That a send failure proves a process crash; that `degraded` proves a failure is nonfatal; that provider state is the process disposition; that a later exit confirms the earlier state was correct.

## 3. U1–U8 scenario assessments

Each case is kept separate. "Current implementation" describes source paths, not an observed run. Where two paths are possible, the winner depends on goroutine interleaving and is not recorded; static reading cannot establish which path occurs in a particular execution. No test was run to settle any interleaving.

### 3.1 U1 — Orderly EOF racing the reader's crash classification

**Initial state and identity.** A helper process is alive and handshaken; the client has a live transport and reader. `Shutdown` and the reader can run concurrently. The process instance is not named by any stable identifier in the evidence; `session` belongs to the client.

**Event/observation sequence.** `Shutdown` sends `shutdown` and waits for a terminal (`client.go:645–652`). The case assumes the helper closes stdout after its acknowledgement or exit. The reader's `Receive` returns an error; `internal_error` is reclassified to `helper_crashed` (`client.go:297–303`). `failAll` stores it as the first read failure and delivers it to pending requests (`:369–386`). Because the code is `helper_crashed`, the reader does not call `Kill` (`:304–306`). `Shutdown` sees `helper_crashed`, skips its early non-timeout branch, and calls `Terminate` (`:653–668`). `Terminate` sets `providerState = stopping` (`supervisor.go:253–259`).

**Concurrent or delayed observations.** The exit observer and `Shutdown` race for `s.mu`. If the observer runs after `stopping` is set, a zero or nonzero exit becomes `absent` (`supervisor.go:176–182`). If it runs before, it becomes `crashed` (`:183–190`). The reader's `failAll` and the observer are separate writers; no sequence or timestamp orders them.

**Evidence retained or lost.** `readFailure` retains only the first failure (the reclassified `helper_crashed`); later failures are discarded (`client.go:373–376`). `waitErr` retains the exit error (`supervisor.go:174`). `lastFailure` may be overwritten by a later `SetState` or not written at all in the stopping branch. The reader-join result is discarded (`client.go:349–362`). There is no record of whether EOF was natural or caused by shutdown.

**Current implementation.** The source contains both outcomes: `absent` if `stopping` is observed, `crashed` otherwise. It does not publish a terminal on the `failAll` path itself; the observer's view of `providerState` decides. This is a real race, not a hypothetical.

**Current P17 rev2.** P3 retains the reader's fatal transport observation but does not infer an exit, a terminal, or a public state; P4/P4-A route a confirmed exit, with P4-A still **PROPOSED/OPEN**. P17 does not decide whether an orderly EOF during shutdown is independent orderly shutdown, a fatal transport observation, or both. It lists U1 as unresolved.

**Unique, justified outcome?** **No.** The model does not determine a unique outcome, and the source does not record the order needed to choose one.

**Needed evidence, rule, or owner decision.** Stream-close provenance (natural EOF vs forced close vs shutdown), an action start/completion record, a stable process-instance identity, and a rule stating whether a completed shutdown action explains an orderly EOF. The owner must decide whether that explanation is sufficient to prefer `absent`, or whether the fatal transport observation remains controlling.

**Prohibited inference.** EOF is not automatically a crash; `absent` is not proof of a clean provider outcome; a completed shutdown is not proof that shutdown caused the EOF.

### 3.2 U2 — Shutdown timeout racing a protocol violation

**Initial state and identity.** A helper is alive, possibly busy; `Shutdown` is in flight. The protocol reader can observe a frame-level violation concurrently. The client session and supervisor process are distinct identities, and neither is recorded with the failure.

**Event/observation sequence.** `Shutdown` waits on a `shutdown` terminal (`client.go:645–652`). A violation in `dispatch` or `readLoop` calls `protocolViolation`, which calls `failAll` and then `Kill` (`client.go:321–367`). `Kill` is `Terminate` with the default grace (`supervisor.go:291–293`); it sets `stopping` and closes stdin (`:253–264`). The pending shutdown request receives the protocol failure. `Shutdown`'s early branch calls `Terminate` again and returns (`client.go:653–655`), or, on a timeout path, falls through to `Terminate` and then clears `started`/`transport` regardless of the result (`:668–674`).

**Concurrent or delayed observations.** The protocol reader, the `Kill` caller, the `Shutdown` caller, and the exit observer can all run concurrently. Two `Terminate` calls can overlap; both set `stopping` and wait on the same `done`. One can return `nil`, the other can return `timeout` after the second grace (`supervisor.go:277–286`). No record says which action completed or which failure preceded which.

**Evidence retained or lost.** `readFailure` keeps the first failure, usually the protocol violation. The `Terminate` timeout is returned to the `Shutdown` caller and ignored by it except as the return value; it is not stored in the supervisor. `waitErr` retains an exit error if one arrives. `lastFailure` is not set by `failAll`, so the provider's `LastError` can remain stale. The ownership result after a timeout is not recorded.

**Current implementation.** If an exit is observed while `stopping` is set, the provider becomes `absent`, including a nonzero exit (`supervisor.go:176–182`). If no exit arrives, the provider remains `stopping` while the supervisor may still own the process; the client clears its transport anyway (`client.go:668–674`). The fatal protocol observation is not published as a terminal state. There is no quarantine marker.

**Current P17 rev2.** P3 retains the fatal protocol observation without publishing `crashed`; P4/P4-A route any confirmed exit, and P10 treats quarantine as internal ownership/reaping uncertainty with mechanics **OPEN**. P17 has no rule for an action timeout where no exit fact exists, and U2 remains OPEN.

**Unique, justified outcome?** **No.** The model supplies no winner, and the source does not record the ordering or the action result needed to construct one.

**Needed evidence, rule, or owner decision.** Arbitration between a pre-existing fatal observation and a shutdown action; a rule for ownership after a terminate timeout; an exit fact or an explicit `outcome_unknown`; and a public projection for the timeout case (or a decision that none is published). The owner must decide whether a pre-existing fatal survives the action.

**Prohibited inference.** A timeout is not a crash; a kill is not proof of the exit's cause; `absent` is not proof of clean termination; a fatal observation is not erased by a later action.

### 3.3 U3 — Incompatibility followed by a crash

**Initial state and identity.** The client is starting a helper and is in the handshake phase. `providerState` is `handshaking` before `MarkHandshake`. The process instance has no stable identity in the failure record.

**Event/observation sequence.** `helloHandshake` returns a failure: either a wire `protocol_incompatible` preserved by `failureFromWire` (`client.go:434–435`, `:471–481`), or a `validateHello` `protocol_violation` (`protocol.go:502–545`). `Start` calls `Terminate` with the shutdown grace and ignores its result (`client.go:248`), then waits for the reader and ignores whether it joined (`:249`), then publishes `incompatible` for the two incompatibility codes or `unavailable` otherwise (`:250–254`), then clears `started`/`transport` (`:255–258`).

**Concurrent or delayed observations.** The exit observer can run before, during, or after the `SetState`. If it runs while `providerState` is `stopping`, it writes `absent` (`supervisor.go:176–182`), and the later `SetState(incompatible/unavailable)` wins. If it runs after the `SetState` and the process was not stopped, it sees a non-`stopping` state and writes `crashed` (`:183–190`), overwriting the handshake projection. If `Terminate` times out, `Start` still publishes the handshake state while the process may be alive; a later exit then becomes `crashed`.

**Evidence retained or lost.** The typed handshake failure is returned to the caller and stored in `lastFailure` by `SetState`. `readFailure` may hold a `helper_crashed` from the reader. `waitErr` may hold an exit error. The `Terminate` result and the reader-join result are discarded. No sequence records whether the handshake failure or the exit was observed first.

**Current implementation.** The source admits `incompatible`/`unavailable` or `crashed`, depending on the observer's timing. The final write wins; there is no terminal latch. A crash after the handshake projection can replace it; a handshake projection after an observer's `absent` can replace that.

**Current P17 rev2.** P2 classifies a missing mandatory service capability as `incompatible` but leaves the other `validateHello` exits unclassified; O15/O16 keep path-by-path classification OPEN. P5 and S3 say a fatal handshake observation does not establish an exit, a terminal, or a public state; P3 supplies the no-exit evidence boundary, and P4 governs a later confirmed exit. Rev2 does not assert that "any exit" leaves the handshake projection controlling. U3 remains OPEN for path classification and terminal/exit ordering.

**Unique, justified outcome?** **No.** Rev2 intentionally leaves path classification and terminal/exit ordering open; the source has no terminal latch. A handshake failure alone does not determine the later projection.

**Needed evidence, rule, or owner decision.** A terminal-publication ordering rule (first terminal vs last write), a stable process-instance identity, the P2 path classification, and a rule for whether an exit after a handshake failure is a new outcome or evidence about the old one. The owner must decide terminal immutability and whether an exit after a handshake failure is a new outcome or merely evidence about the old one.

**Prohibited inference.** A handshake failure did not cause the crash; a later exit is not automatically the correct final state; `incompatible` is not proof that the process exited cleanly.

### 3.4 U4 — Zero exit after a fatal protocol violation

**Initial state and identity.** A helper is alive and may have pending requests. A fatal protocol violation is observed. The process identity is not recorded with the violation.

**Event/observation sequence.** `dispatch` detects an empty ID, duplicate terminal, or unknown ID (`client.go:321–346`), or `readLoop` detects a frame/decode violation (`:309–313`). `protocolViolation` calls `failAll` and `Kill` (`:364–367`). `Kill`/`Terminate` sets `stopping` and closes stdin (`supervisor.go:253–264`). The helper exits with status zero. The observer sees `stopping` and writes `processState = absent`, `providerState = absent` (`:176–182`). The `protocol_violation` remains in `readFailure`; it is not published by `failAll`.

**Concurrent or delayed observations.** The violation, the kill, and the exit observer can be concurrent. If a pending synthesis receives the violation through `failAll`, it publishes `degraded` (`client.go:790–797`), not `crashed`. The provider state and the pending request's projection can therefore disagree.

**Evidence retained or lost.** `readFailure` retains the protocol violation only if it is first. `waitErr` retains `nil` for a zero exit. The action that killed the process is not recorded with intent or result. There is no sequence or timestamp linking the violation to the exit.

**Current implementation.** The source can produce `absent` (observer sees `stopping`), `degraded` (pending synthesis receives the fatal failure), or a later overwrite, while the fatal protocol evidence remains only in `readFailure`. The source never publishes `crashed + protocol_violation` on this path.

**Current P17 rev2.** P3 retains a fatal protocol observation as internal evidence but does not publish `crashed`, infer an exit, or establish a terminal. P4/P4-A govern a confirmed exit; P4-A is **PROPOSED/OPEN**. S3 and S9 distinguish fatal from crash, and OD-5 defines `crashed` only for an observed exit without orderly shutdown. Rev2 therefore does not map fatal + zero exit to `crashed`. The open issue is the classification threshold and orderly-exit determination, not a rule-2 contradiction.

**Unique, justified outcome?** **No.** Rev2 removes the fatal-to-`crashed` rule, but the classification threshold, P4-A, and the orderly-shutdown evidence remain unresolved; the source does not establish that the violation caused the exit.

**Needed evidence, rule, or owner decision.** The failure-classification threshold, an orderly-shutdown predicate, the exit status/signal, and terminal authority. The owner must decide whether a zero exit after a fatal observation can ever satisfy `crashed`, or whether the outcome remains OPEN.

**Prohibited inference.** A fatal diagnostic is not a crashed process; a zero exit is not proof of a clean session; the kill is not proof that the violation caused the exit.

### 3.5 U5 — Nonzero exit after a deliberate kill

**Initial state and identity.** A helper is alive. A deliberate `Kill`/`Terminate` is issued by a protocol path, cleanup, or shutdown. The process identity and the action identity are not recorded together.

**Event/observation sequence.** `Terminate` sets `stopping`, closes stdin, and waits (`supervisor.go:246–264`). If the grace expires, it calls `p.Kill()` and waits again (`:277–286`). The process exits nonzero. The observer sees `stopping`, sets `processState = exited`, and publishes `providerState = absent` (`:176–182`). `waitErr` retains the nonzero exit error (`:174`).

**Concurrent or delayed observations.** If the exit is observed before `Terminate` sets `stopping`, the observer publishes `crashed` (`:183–190`). If the action and exit overlap, the order is not recorded. A nonzero exit after a kill is indistinguishable from a nonzero exit that was already in progress.

**Evidence retained or lost.** `waitErr` holds the exit error, but it is not separated into an exit code and a cause. The action's intent, grace, start, completion, and result are absent. `Terminate` results are ignored at several call sites (`client.go:248`, `:654`, `:661`, `:664`, `:861`). The supervisor does not store `terminated_by_action`.

**Current implementation.** The source publishes `absent` for a nonzero exit observed while `stopping`; it does not publish `crashed`. P17's internal `terminated_by_action` has no source representation.

**Current P17 rev2.** OD-4 is accepted: a confirmed exit without orderly shutdown projects `unavailable` unless an applicable prior justified failure is preserved; an independently established orderly shutdown projects `absent`; `terminated_by_action` is internal with no public projection; and no active state survives exit. P17 S7/S8/S11 and T5–T10 encode those boundaries. P4-A remains **PROPOSED/OPEN**. The missing element is the OD-8 A–E evidence that the action caused the exit, not a decision on whether a public projection exists.

**Unique, justified outcome?** **No.** The public policy is now accepted, but orderly shutdown and action causation are not established by the source's `stopping` observation. P4-A is also unaccepted.

**Needed evidence, rule, or owner decision.** OD-8 A–E evidence, action identity and outcome, exit status/signal, and an independently established orderly shutdown. This is an evidence gap, not an open public-projection choice.

**Prohibited inference.** A nonzero exit is not automatically a crash; a kill is not automatically the cause of the exit; `absent` is not proof that the termination was clean or intended.

### 3.6 U6 — Provider `crashed` while the process remains alive

**Initial state and identity.** A helper is alive and the provider is ready or busy. A synthesis write fails. The process and provider are distinct identities, but the state model does not keep them distinct at publication.

**Event/observation sequence.** `transport.Send` returns an error (`client.go:755`). `Synthesize` publishes `TTSStateCrashed` with a `helper_crashed` failure (`:756`) before calling `cleanupRejected` (`:757`). `cleanupRejected` checks `sup.ProcessAlive()` (`:859–861`); if true, it calls `Terminate` and ignores the result (`:861`). `SetState` sets `processState = alive` when `s.proc != nil` and the state is not `stopping` (`supervisor.go:239–240`), so the publication itself does not mark the process exited.

**Concurrent or delayed observations.** The exit observer may subsequently see `stopping` and publish `absent` (`supervisor.go:176–182`), overwriting `crashed`; or it may see `crashed` and publish `crashed` again (`:183–190`). If `Terminate` times out, the process can remain alive while the provider is `crashed`.

**Evidence retained or lost.** `lastFailure` holds the `helper_crashed` failure. There is no separate process-outcome record, no action result, and no ownership/uncertainty field. `Status` derives `ProcessAlive` from both state dimensions (`supervisor.go:324`); a live OS process can be reported as not alive because the provider state is `crashed`.

**Current implementation.** The source demonstrably publishes `crashed` before establishing that the process exited. It may terminate the process afterward, but the publication is not evidence of the exit. The later state can be `absent` or `crashed`, depending on timing.

**Current P17 rev2.** S9/O18/P3 says a live process must not be classified `crashed`; a fatal transport observation alone does not establish an exit or a public state. The source's `SetState(crashed)` at `client.go:756` is nonconforming on this point. With no exit, P17 requires preserving the last justified state; if none is justified, it selects no projection and leaves U6/C-01 OPEN.

**Unique, justified outcome?** **No.** Rev2 resolves whether a live process may be `crashed`; it does not resolve what current state is justified after a send failure with no established exit.

**Needed evidence, rule, or owner decision.** The post-failure projection when no exit is established, a typed cleanup outcome, and actual process liveness. The owner does not need to revisit whether a live process can be `crashed`; OD-5 answers that.

**Prohibited inference.** A write failure is not a process crash; `crashed` is not proof that the process is dead; a later `Terminate` is not proof of reaping.

### 3.7 U7 — Terminal state published before the termination timeout

**Initial state and identity.** A terminal provider state has been published (for example, `incompatible` after a handshake failure, or `absent` after an exit). A `Terminate` is still outstanding and may time out. The terminal's authority is not recorded.

**Event/observation sequence.** `SetState` writes a terminal state and `lastFailure` (`supervisor.go:233–242`). A later `Terminate` sees `s.proc != nil`, sets `providerState = stopping` (`:253–259`), and may time out after the second grace (`:277–286`). If the process later exits, the observer writes `absent` or `crashed` based on whether it sees `stopping` (`:176–190`). In the handshake path, `Start` publishes `incompatible`/`unavailable` after `Terminate` returns, even if it returned a timeout (`client.go:248–254`). In `Shutdown`, the client clears `started`/`transport` after a timeout (`:668–674`).

**Concurrent or delayed observations.** The terminal write, the `Terminate` state write, and the observer can interleave. `Terminate` unconditionally overwrites a terminal with `stopping` if a process is still owned; the observer can then overwrite `stopping` with `absent`/`crashed`. No latch prevents any of these writes.

**Evidence retained or lost.** `lastFailure` is overwritten by `SetState` (`supervisor.go:238`). There is no terminal timestamp, terminal flag, or revision record. The client's `closing` flag (`client.go:670`) is not a supervisor terminal latch. The `Terminate` result is often ignored.

**Current implementation.** A published terminal is mutable. It can be replaced by `stopping`, `absent`, or `crashed`; the final value depends on write order. The source does not define terminal immutability.

**Current P17 rev2.** P9 accepts an immutable first terminal, and P4/Case A preserve it, but the terminal set, linearization point, CP-1, and late-evidence effects remain **OPEN**. P17 does not define every value that can be terminal or whether a later action may revise a terminal state.

**Unique, justified outcome?** **No.** The model does not define terminal authority, so it cannot determine whether the terminal or the later action wins.

**Needed evidence, rule, or owner decision.** A terminal-authority rule: once published, may a terminal be revised, and by what evidence? How is a timeout represented? Does ownership uncertainty override a terminal? The owner must decide terminal immutability and the relationship to quarantine.

**Prohibited inference.** Append-only evidence does not make a projection immutable; a later exit is not automatically a correction; a timeout is not a terminal state.

### 3.8 U8 — Two handshake projections for one compatibility concept

**Initial state and identity.** The client is handshaking with a helper. The helper's capabilities are being validated. The process may be alive or may exit after the rejection.

**Event/observation sequence.** Case A: the helper returns a wire error with code `protocol_incompatible`; `failureFromWire` preserves it (`client.go:434–435`, `:471–481`); `Start` publishes `incompatible` (`:250–251`). Case B: the helper omits a mandatory service capability; `validateHello` returns `protocol_violation` with a `missing` detail (`protocol.go:527–531`); `Start` publishes `unavailable` (`client.go:252–253`). Both describe a helper that cannot satisfy the required compatibility set.

**Concurrent or delayed observations.** This case does not require a race. The projection follows deterministically from which path produced the failure. A later exit can still overwrite either projection as in U3.

**Evidence retained or lost.** The failure code and `Detail["missing"]` are retained in `lastFailure` and the returned error. The reason the code was chosen is not recorded beyond the code itself. There is no compatibility classification field.

**Current implementation.** The source demonstrably has two projections for the same concept: helper-declared incompatibility → `incompatible`; client-detected missing mandatory capability → `unavailable`. Both are reachable.

**Current P17 rev2.** P2 (accepted OD-6) classifies a missing mandatory service capability as `incompatible`; the source still maps `validateHello`'s `protocol_violation` to `unavailable`. P2 also classifies several other conditions and leaves misplaced descriptors, unsupported extension versions, optional/feature requiredness, and unlisted cases unclassified; O16/U8 remain OPEN. P17 no longer preserves the old split as policy, but source and proposal disagree on the missing-capability path.

**Unique, justified outcome?** **No for the unclassified paths; yes for the missing-mandatory-capability policy.** P2 selects `incompatible`, but the source does not implement it, so U8 remains OPEN rather than satisfied.

**Needed evidence, rule, or owner decision.** Apply P2 to the source path, and classify or explicitly leave the remaining `validateHello` exits open. Record the condition/code mapping; a code-name lookup is not a classification rule.

**Prohibited inference.** A code name is not the helper's intent; a missing capability is not automatically `incompatible`; helper-declared and client-detected failures are not the same fact; `unavailable` is not a neutral default without a decision.

## 4. Mimo reconciliation

The 547-line independent review was checked against the pinned helper source, P17 rev2, the accepted decisions, and contract A6. Classifications describe source support, not observed runtime. No interleaving was executed; a source-reachable schedule is not a measured occurrence.

| ID | Mimo finding | Classification | Basis and engineering consequence |
|---|---|---|---|
| F-01 | OD-9 gate unmet; U1–U8 open | **CONFIRMED** | Decision register records OD-9 Option A; P17 §10 leaves all cases OPEN. P17 cannot be approved while this remains. |
| F-02 | No terminal latch; terminal set undefined | **CONFIRMED** | `SetState`, `Terminate`, and `observeExit` can all overwrite provider state (`supervisor.go:176–195`, `:233–242`, `:246–288`). There is no terminal flag, set, or linearization point. |
| F-03 | No generation isolation; old reader can affect a replacement | **CONFIRMED** | The reader captures process `p` and calls `failAll`/`Kill` without a generation check (`client.go:295–318`, `:364–386`); `Start` resets session and pending state (`:191–270`). |
| F-04 | `Shutdown` can start a replacement | **CONFIRMED** | `Shutdown` calls `c.call` (`client.go:645–652`), which calls `Start` (`:398–443`); `closing` is set too late to prevent start (`:668–674`). |
| F-05 | Contract A6 conflicts with OD-5/P17 | **CONFIRMED** | A6 rule 3 and its table project `crashed + kind` for a fatal observation even at exit 0; S9/P3 and OD-5 prohibit `crashed` without an observed exit and no orderly shutdown. |
| F-06 | `MarkHandshake` can resurrect a dead process | **CONFIRMED** | `MarkHandshake` sets `alive`/`unloaded` without checking `s.proc`, generation, or terminal state (`supervisor.go:206–220`). A late successful handshake can overwrite exit evidence. |
| F-07 | `ready`/`degraded` after termination or exit | **CONFIRMED** | `Synthesize` writes `ready` after cleanup (`client.go:771`, `:783`, `:849`) and `degraded` after fatal pending failures (`:790–797`); no post-exit publication guard exists. Source-reachable, not observed. |
| F-08 | `ProcessAlive` is not independent | **PARTIALLY CONFIRMED** | The public `Status.ProcessAlive` uses `processAlive(process, provider)` and returns false for `crashed`/`unavailable` (`supervisor.go:305–342`, `:351–359`). The `ProcessAlive()` method itself checks ownership and process state (`:335–342`); the defect is the public derivation, not every liveness check. |
| F-09 | Reader is never joined; replacement race | **CONFIRMED** | `waitForReader` returns nothing on join or timeout (`client.go:349–362`); `Start` discards the result (`:249`). A replacement can begin while the old reader is still acting. |
| F-10 | Kill without a reaped guard; repeated `Terminate` | **CONFIRMED** | `Terminate` has no reaped flag and may select the timer branch and call `Kill` while `done` is closing (`supervisor.go:277–286`); repeated callers share the same `done`. |
| F-11 | `Shutdown`/`Start` have no common admission barrier | **CONFIRMED** | Client `startMu` and supervisor `startMu` are distinct; `Shutdown` checks `currentTransport()` without the start lock (`client.go:645–652`). A concurrent `Start` can complete after `Shutdown` returns. |
| F-12 | Cleanup outcome is not retained | **CONFIRMED** | `cleanupRejected` ignores `Terminate`'s result and retains no typed outcome (`client.go:859–864`). Retention can be indistinguishable from successful cleanup. |
| F-13 | U assessment is stale on U3/U4/U5/U6 | **CONFIRMED** | The prior text cited the frozen original P17's rule 1/rule 2/rule 3 and an undecided public projection. This pass rebinds those sections to P17 rev2 and OD-4/OD-5. |
| F-14 | P17 `absent` depends on unaccepted P4-A | **CONFIRMED** | P4-A is **PROPOSED/OPEN**; S8/O4 and the no-orderly-shutdown branches depend on it. OD-4's orderly-shutdown rule is accepted, but the shared eligibility predicate is not. |
| F-15 | `outcome_unknown` is not assigned to open cases | **CONFIRMED** | P17 §10 leaves U1–U8 OPEN; P7 permits but does not require an internal unknown record. No explicit per-case unknown rule satisfies OD-9. |
| F-16 | No action/exit status evidence for OD-8 | **CONFIRMED** | OD-8 is accepted, but the source has no action identity/outcome or exit status/signal record; `Terminate` results are discarded and `waitErr` is not decomposed (`supervisor.go:246–288`). |
| F-17 | No quarantine transition | **CONFIRMED** | P7/P10 describe quarantine internally; neither process nor provider state has it, and a timeout leaves no ownership/reaping record (`supervisor.go:277–286`). |
| F-18 | `retired` is unbounded; status starts the helper | **PARTIALLY CONFIRMED** | The unbounded `retired` map is present (`client.go:339–342`, `:464–468`) and `ProviderStatus` calls `Start` (`:952–958`). The full catalog call path was not re-traced in this pass; the finding is low severity and not a blocker by itself. |

No material Mimo finding was **NOT CONFIRMED**. F-08 and F-18 are narrowed rather than rejected. No finding was promoted from unresolved to resolved; the table records what source supports, not a runtime result.

## 5. Engineering transition handoff

This handoff is planning input, not authorization. Each area states the behavioral correction required; it does not select an implementation mechanism. The accepted decisions constrain the plan, while the open items below remain gates.

### 5.1 Generation and process ownership

**Confirmed problem.** The client has session state but no process-generation token. `readLoop` captures the `Process` it was started with and calls `failAll`, `protocolViolation`, and `Kill` without confirming that the process is still current (`client.go:295–318`, `:321–347`, `:364–386`). `Start` resets `retired` and session state on restart (`:191–270`), so a delayed response from the old stream can be classified as an unknown ID and kill the replacement. `Terminate` and the exit observer identify the process only by the mutable `s.proc`/`s.procDone` fields (`supervisor.go:166–195`, `:246–288`).

**Required behavior.** Bind every reader, action, failure, exit fact, and publication to an immutable generation. A stale generation may append evidence but must not fail, kill, or publish against the current generation. Route responses by generation and request identity, not by a reused request ID alone. Retire and join the old reader before installing a replacement, or make the old reader's operations no-ops. Keep one ownership boundary for writers and process handles.

**Blocking boundary.** OD-2 and OD-3 are accepted and sufficient for planning. The allocator, join mechanism, and storage representation remain design choices, not missing evidence. The terminal set and P4-A still affect what a stale publication is allowed to do.

**Verification.** U1, U2, U3, and U7, plus regression cases for late frames after restart, old-reader failure after replacement, request-ID reuse across generations, and `Shutdown` racing `Start`. Verify that old-generation evidence is retained without changing the current status.

### 5.2 Lifecycle and terminal correctness

**Confirmed problem.** Provider state is a mutable last-writer field. `SetState` can mark a process alive after an exit (`supervisor.go:233–242`), `MarkHandshake` can restore `alive`/`unloaded` without a process check (`:206–220`), and `observeExit` chooses `absent` or `crashed` from whether `stopping` happened to be set (`:166–195`). `Shutdown` can start a helper through `c.call` and then clear the transport, and `Synthesize` can publish `ready` or `degraded` after cleanup has terminated the process (`client.go:645–674`, `:755–797`, `:849`).

**Required behavior.** Define the terminal set, one linearization point, and one immutable first terminal per generation. Retain later evidence without allowing it to replace the terminal. No active state may survive a confirmed exit, and no handshake or cleanup path may publish after the terminal. `Shutdown` must close admission and must never call `Start` or create a replacement. Liveness must come from the actual process handle and reaping state, not from the aggregate provider label.

**Blocking boundary.** OD-3, OD-4, and OD-5 are accepted. The terminal set, CP-1, and the `absent` eligibility predicate P4-A remain open. Those are policy/design gates; they are not reasons to repeat the source review.

**Verification.** U1, U2, U3, U5, U6, and U7, plus `Shutdown`/`Start` overlap, handshake completion after exit, cleanup completion versus `ready`, repeated `Terminate`, and kill-after-reap. Assert first-terminal stability and no post-exit active projection.

### 5.3 Protocol and compatibility

**Confirmed problem.** A missing mandatory service capability reaches `validateHello` as `protocol_violation` (`protocol.go:527–531`) and `Start` maps it to `unavailable` (`client.go:250–254`), while helper-declared `protocol_incompatible` maps to `incompatible`. P17 P2 and OD-6 require the missing-capability case to be `incompatible`. Other `validateHello` exits remain unclassified in O16/U8.

**Required behavior.** Apply P2's categories at the observation site. Missing mandatory capability becomes `incompatible`; duplicate identities, invalid descriptors, and malformed version/limits follow their assigned categories; unlisted, extension, and optional-requiredness cases remain explicitly OPEN rather than falling through to a generic state. Do not publish `crashed` for a live process, and do not publish `ready`/`degraded` after a confirmed exit.

**Blocking boundary.** OD-6 settles the missing-capability policy. The path-by-path classification of the remaining nine `validateHello` exits is an open taxonomy question; it can be planned as an explicit matrix, but it cannot be silently decided by implementation.

**Verification.** U6 and U8, with a table of every `validateHello` exit and its condition, wire code, P2 category, and public projection. Include helper-declared incompatibility, missing mandatory capability, duplicate identity, malformed descriptor, and an unlisted/extension case. Confirm no live-process `crashed` write remains.

### 5.4 Contract and evidence alignment

**Exact A6 conflict.** Contract rev3 A6 rule 3 says a fatal `failureObserved` projects `crashed` with that kind even when `exitStatus == 0`; its table repeats `0 | fatal | any | crashed + kind`. OD-5/S9 and P17 P3 say `crashed` requires an observed exit without established orderly shutdown, and a fatal observation alone establishes neither an exit nor a terminal. A6 therefore cannot be used as the implementation basis while it retains that rule.

**Required behavior.** Correct or supersede A6 before treating contract rev3 as an implementation authority. Record the accepted OD-4 and OD-5 boundaries, retain the fatal observation as evidence, and do not infer a process outcome from it. P4-A, the terminal set, CP-1, and the per-case unknown rule must be resolved or explicitly carried as OPEN before P17 approval; they may be planned now but not assumed.

**Safe to resolve during implementation.** Lock placement, generation allocator, reader-join mechanism, record storage, sequence allocation, retention windows, and internal quarantine representation do not need another governance cycle if they preserve the accepted behavioral boundaries.

**Verification.** Cross-check A6 against S7/S8/S9 and P3/P4; test U3, U4, and U5 with exit 0 and nonzero exits, with and without an established orderly shutdown. Confirm A6 is not silently implemented as written.

## 6. Exact source evidence

The following table consolidates the citations used above. Line numbers are from the audited helper tree and are stable for the recorded digest.

| ID | Source location | Exact observed fact |
|---|---|---|
| SRC-01 | `internal/audio/tts/helper/errors.go:15–21` | `Failure` has `Code`, `Message`, `Retryable`, `Field`, `Detail`; no instance, sequence, timestamp, phase, severity, action, or cause. |
| SRC-02 | `internal/audio/tts/helper/client.go:130–140` | Client state includes `transport`, `started`, `session`, `readerDone`, `next`, `pending`, `retired`, `readFailure`, `closing`; `next` is a request counter, not an evidence sequence. |
| SRC-03 | `client.go:191–259` | `Start` rotates the session on restart, starts the reader goroutine, and on handshake failure calls `Terminate` (result ignored), `waitForReader` (result discarded), then `SetState(incompatible/unavailable)` and clears `started`/`transport`. |
| SRC-04 | `client.go:295–318` | Reader reclassifies `internal_error` to `helper_crashed`, calls `failAll`, and kills only for `invalid_frame`, `frame_too_large`, or `protocol_violation`. |
| SRC-05 | `client.go:321–347` | Dispatch reports empty ID, duplicate terminal, and unknown ID as `protocol_violation`; a retired ID is silently dropped after a comment claiming it is evidence. |
| SRC-06 | `client.go:349–362` | `waitForReader` returns nothing on either join or timeout. |
| SRC-07 | `client.go:364–386` | `protocolViolation` calls `failAll` then `Kill`; `failAll` stores only the first read failure and delivers it to all pending requests. |
| SRC-08 | `client.go:398–443` | `callOn` has generic timeout and cancellation branches with no operation, phase, or deadline-origin record; send failure returns `helper_crashed` without termination. |
| SRC-09 | `client.go:445–469` | `begin` increments `next` and builds a request ID; `retire` moves the ID to an unbounded `retired` map. |
| SRC-10 | `client.go:645–674` | `Shutdown` sends the frame, may return early after `Terminate`, otherwise calls `Terminate` and clears `started`/`transport` regardless of the result. |
| SRC-11 | `client.go:755–758` | Synthesis send failure publishes `TTSStateCrashed` before `cleanupRejected`; the process may still be alive. |
| SRC-12 | `client.go:790–797` | A pending failure publishes `TTSStateDegraded`; a fatal `protocol_violation` can therefore emerge as `degraded`. |
| SRC-13 | `client.go:799–829` | Wire, decode, result-validation, and WAV-validation failures all publish `degraded` and overwrite `lastFailure`. |
| SRC-14 | `client.go:859–864` | `cleanupRejected` terminates only if `ProcessAlive()` and either force or output exists; the `Terminate` result is ignored. |
| SRC-15 | `supervisor.go:166–195` | `observeExit` decides `absent` vs `crashed` solely from `providerState == stopping`; nonzero exit in the stopping branch becomes `absent` with `processState = exited`; `lastFailure` is overwritten only in the unexpected branch. |
| SRC-16 | `supervisor.go:197–203` | `setUnavailable` overwrites `lastFailure` and sets process state `absent`. |
| SRC-17 | `supervisor.go:233–242` | `SetState` overwrites `providerState` and `lastFailure`; it sets `processState = alive` whenever a process is owned and the new state is not `stopping`. |
| SRC-18 | `supervisor.go:246–288` | `Terminate` has no intent, grace, start, completion, or result record; it sets `stopping`; it kills on context cancellation or grace expiry; after a second grace it returns `timeout` while the process may remain owned. |
| SRC-19 | `supervisor.go:291–293` | `Kill` is `Terminate` with the default grace; it is not an independently recorded action. |
| SRC-20 | `supervisor.go:305–342`, `:351–356` | `Status` derives `ProcessAlive` from both process and provider state; `ProcessAlive()` checks only process ownership and `processState != exited`; the helper returns false for `crashed` or `unavailable`. |
| SRC-21 | `protocol.go:502–545` | `validateHello` has nine failure exits, all `protocol_violation`; missing mandatory capabilities is one of them with a `missing` detail. |
| SRC-22 | `protocol.go:259–268` | `newSessionID` is a random session; `newRequestID` formats `req-<session>-<counter>`. |
| SRC-23 | `internal/audio/types/provider.go:8–36` | The process-state set is `absent`, `starting`, `alive`, `stopping`, `exited`; the provider-state set has thirteen values; neither contains `quarantined`. |
| SRC-24 | `client.go:191–270` | `Start` has no generation token; it can clear `started`/`transport`, reset `retired`, rotate the session, start a reader for a new process, and republish after a failed handshake. |
| SRC-25 | `client.go:645–674`; `:398–443` | `Shutdown` acquires no admission barrier, calls `c.call`, and can therefore invoke `Start`; `closing` is set only late in the shutdown path. |
| SRC-26 | `supervisor.go:206–220` | `MarkHandshake` unconditionally sets process `alive` and provider `unloaded`; it checks neither the current process, generation, nor terminal state. |
| SRC-27 | `client.go:755–797`, `:849` | `Synthesize` can publish `crashed` before exit, `ready` after timeout/cancellation cleanup, and `degraded` for fatal pending failure, with no post-exit publication guard. |
| SRC-28 | `supervisor.go:305–342`, `:351–359` | Public `ProcessAlive` is derived from both process and provider state and is false for `crashed`/`unavailable`; `ProcessAlive()` itself is ownership/process-state based. |

## 7. Owner decisions and remaining boundaries

OD-1 through OD-9 are accepted. This table records what each decision settles and what remains open; it does not reopen or reinterpret an accepted decision.

| ID | Accepted decision | Remaining boundary |
|---|---|---|
| OD-1 | Projection-supporting evidence; supervisor authority; client observations retained | Record schema, allocator, retention, and conflict mechanics remain open. |
| OD-2 | Composite session + supervisor generation + PID; generation is process-lifetime key | Generation allocator, join mechanism, and storage representation remain open. |
| OD-3 | Immutable first terminal per process generation; later evidence retained | Terminal set, linearization mechanism, and late-evidence effects on non-terminal projections remain open. |
| OD-4 | Confirmed exit without orderly shutdown → `unavailable` unless prior justified failure; orderly shutdown → `absent`; `terminated_by_action` internal | P4-A remains proposed; action-caused nonzero-exit evidence and prior-public-value interactions remain open. |
| OD-5 | `crashed` = observed exit without established orderly shutdown; live process must not be `crashed`; `ProcessAlive` independent | Five-state implementation mapping and the post-failure/no-exit projection remain open. |
| OD-6 | Missing mandatory capability → `incompatible`; structural ambiguity/unestablished meaning → `unavailable` | Path-by-path classification of the nine `validateHello` exits and extension-version scope remain open. |
| OD-7 | `quarantined` internal-only, exact generation, explicit ownership/reaping condition | Internal record, transition, and release rules remain open. |
| OD-8 | Six-dimension A–E evidence standard; association is not causation | Source action/exit evidence and U5's causal attribution remain unimplemented. |
| OD-9 | P17 is not approved while U1–U8 are OPEN without determined outcomes or explicit unknown rules | Gate remains unmet; U1–U8 all remain OPEN. |

## 8. Missing source or runtime evidence

These are facts the current source does not retain, or runtime facts that static reading cannot establish.

| ID | Missing evidence | Consequence |
|---|---|---|
| ME-1 | No evidence log, record type, or sequence allocator exists. | R09 cannot be evaluated as implemented; P17's `seq` has no source counterpart. |
| ME-2 | No stable process-instance or supervisor-generation identity is stored with failures, actions, or exits. | U1–U7 cannot attribute an observation to a specific process lifetime. |
| ME-3 | No timestamps or monotonic order accompany failures, actions, or exits. | Correlation cannot be distinguished from causation (R10). |
| ME-4 | No action intent, grace, start, completion, or result is recorded; `Terminate` results are ignored at five sites and `Cancel` at one. | U2, U5, U6, and U7 cannot establish what action occurred or whether it completed. |
| ME-5 | No exit code or cause is separated from `waitErr`; `exec.ExitError` is retained only as an error. | U4 and U5 cannot distinguish zero/nonzero, kill-induced exit, and independent crash. |
| ME-6 | No stream-close provenance distinguishes natural EOF, shutdown-induced EOF, forced close, and read error. | U1 cannot decide whether the reader's `helper_crashed` is a crash or an expected consequence. |
| ME-7 | `waitForReader` discards whether the reader joined or timed out. | U1, U3, and U7 cannot establish whether the reader was still active when state was published. |
| ME-8 | A late terminal for a retired request is dropped without a record (`client.go:339–342`). | The evidence claim in the comment is not implemented; E7 remains. |
| ME-9 | No terminal latch or revision record exists in the supervisor. | U3 and U7 cannot establish terminal immutability. |
| ME-10 | No quarantine or ownership-uncertainty state exists in either enum. | R11 and U2 have no internal representation, let alone a public one. |
| ME-11 | No runtime interleaving evidence was collected; no tests were run. | U1–U7 outcomes are source-admissible paths, not observed executions; the winner remains unverified. |
| ME-12 | No platform-specific evidence for `KillTree`, exit-status mapping, or pipe-close behavior was collected. | Windows-specific behavior of the kill/exit paths remains unverified. |
| ME-13 | No helper-side diagnostic record is retained beyond the bounded stderr ring. | The helper's own view of a violation or shutdown is not available to the supervisor. |
| ME-14 | No current-process or generation guard in `MarkHandshake`. | A late handshake can restore `alive`/`unloaded` after exit or replacement. |
| ME-15 | No post-exit publication guard in `Synthesize`. | `ready` or `degraded` can overwrite the result of termination. |
| ME-16 | No shared admission barrier between `Start` and `Shutdown`. | A shutdown can return successfully while a concurrent start leaves a live helper. |
| ME-17 | No reader join/retirement ownership. | An old reader can fail or kill the current generation; late responses can be treated as protocol violations. |
| ME-18 | No reaped flag or action-outcome record in `Terminate`. | A kill can follow reaping, and repeated termination cannot be distinguished from a single completed action. |

## 9. Dependency ordering

The questions cannot be evaluated in arbitrary order. The following is the minimum ordering implied by the dependencies, not a work plan.

| Order | Must settle first | Because |
|---|---|---|
| 1 | Apply accepted OD-1 and OD-2: log authority and instance identity | R09's record schema and sequence domain depend on both; without them, R10 has no ordered evidence to resolve. |
| 2 | R09: minimal record, allocator, causal link, action result, phase, retention | R10 and U1–U7 consume these fields; a projection rule over undefined records is not evaluable. |
| 3 | Apply OD-3, OD-5, OD-6, OD-7: terminal authority, `crashed` semantics, compatibility taxonomy, quarantine status | R11–R13 and U3, U6, U7, U8 cannot be projected into a state set whose authority and meaning are undecided. |
| 4 | R10 and OD-8: causality standard and evidence | U1–U5 depend on whether an action or fatal observation can be said to explain an exit. |
| 5 | U1–U8: individual outcomes or explicit `outcome_unknown` | Each case needs the prior rules; no case should be closed by assigning a higher precedence to one event. |
| 6 | P17 approval gate (OD-9) | Only after 1–5 can the proposal be judged; implementation remains unauthorized regardless. |

R12/U8 can be decided with the taxonomy at step 3; R13/U6 can be decided with the `crashed` semantics at step 3. The scenarios still depend on the evidence rules from steps 1–2 for a defensible outcome.

## 10. Verdict and recommendation

**READY FOR ENGINEERING PLANNING.**

This recommendation means the evidence is sufficient to produce a concrete implementation plan. It is not P17 approval and it does not authorize C1 implementation. The current source establishes real lifecycle hazards, and the remaining open items are bounded behavioral decisions, not missing evidence requiring another review cycle.

### 10.1 What is settled

- The Mimo findings are materially confirmed: generation cross-talk, shutdown-initiated replacement, missing terminal immutability, post-exit active projections, non-independent public liveness, missing action/exit evidence, and the A6 contradiction are source-verified.
- OD-1 through OD-9 are accepted. The supervisor is the evidence/publication authority, generation is the process-lifetime identity, the first terminal is immutable, `crashed` is an observed exit without orderly shutdown, missing mandatory capability is `incompatible`, quarantine is internal-only, and OD-8's A–E standard governs causation.
- R09–R13 remain located and reproducible. U1–U8 remain OPEN; U6 and U8 have partial policy corrections but are not closed.
- The engineering transition is organized into four concrete work areas: generation/process ownership, lifecycle/terminal correctness, protocol/compatibility, and contract/evidence alignment.

### 10.2 What remains gated

- **P17 remains unapproved.** OD-9 Option A is unmet because U1–U8 are OPEN without determined outcomes or explicit unknown-case rules.
- **P4-A, the terminal set, and CP-1 remain open.** `absent` eligibility cannot be implemented as accepted policy until P4-A is accepted or replaced; terminal preservation of a prior public `crashed` remains CP-1.
- **Contract A6 remains conflicting.** It must not be used as written. A6 rule 3 and its table project `crashed` from a fatal observation at exit 0, contrary to OD-5/S9 and P3.
- **No runtime evidence exists.** All interleavings are source-reachable schedules. The implementation plan must include the verification cases above; it cannot claim current runtime behavior.

### 10.3 Why this is not another evidence pass

The reviews now agree on the mechanisms. Another independent reviewer would not supply a missing source file or a missing policy decision; it would repeat the same findings. The remaining work is to turn the four handoff areas into a bounded design and verification plan, while carrying the named open decisions as explicit gates.

No source change, test, build, benchmark, runtime experiment, contract amendment, rev4, P18, P17 approval, or C1 implementation was performed or is authorized by this assessment. The thirteen-value provider-state set remains frozen. `outcome_unknown`, `terminated_by_action`, and `quarantined` remain internal classifications and must not be translated automatically into a public provider state.
