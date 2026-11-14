# C1 — P17 Review Reconciliation (R01–R13, E1–E9, U1–U8)

**Date:** 2026-11-09
**Status:** AUTHOR RECONCILIATION. **P17 IS NOT APPROVED.** Not a revision, not rev4, not P18.
**Responds to:** `REVIEW-howl-tts-stage3-c1-p17-independent-2026-11-09.md`, preserved verbatim.
**Subject:** P17 proposal `01d887c1…`, preserved unchanged.
**C1 IMPLEMENTATION REMAINS UNAUTHORIZED.**

No contract, source, test, or build change. Every claim below is tied to a line in the audited tree.

## 1. Provenance

| Item | Value |
|---|---|
| HEAD | `9c7d2a1` |
| helper tracked by Git | 0 files; `?? internal/audio/tts/helper/` |
| helper Go files | 19 |
| helper total lines | 4,992 |
| helper content digest | `d0df8b2823389b305698fc06eda28e0e01f06259fdc69174cf7fb3851e271a51` |
| P17 proposal | `01d887c1…`, unchanged |
| review of record | preserved verbatim, path above |
| working tree | 142 dirty entries, pre-existing |

The HEAD pin does not establish helper identity, because the directory is untracked. The content digest above is the only source-identity claim this reconciliation makes. It covers `helper/*.go` in sorted path order; it does not cover non-Go files or the supervisor's runtime state.

**Correction to P17, recorded here rather than silently edited:** P17 cited `O13`/`O14` at `client.go:439/441`; those lines are in `Synthesize`'s own select, while `callOn` has separate timeout/cancellation branches. P17's `O19` also omitted the degraded writes at `803,806,814,820,828`. P17 remains frozen; these defects are acknowledged, not repaired in place.

## 2. Verdict

The NOT APPROVED verdict is accepted. P17 is a partial inventory, not a dependable specification. No contract amendment, no rev4, no P18, no implementation. The reconciliation below does not seek approval.

## 3. R01–R08 — source-confirmed defects

### R01 — Operation context is missing

- **Source behavior:** `callOn` (`client.go:398–431`) returns generic `timeout` and `cancelled` failures without recording `op`, phase, or which deadline expired. The same helper serves `hello` with `HandshakeTimeout` (`:281`), `load` with `OperationTimeout` (`:580`), `shutdown` (`:652`), and `Synthesize` (`:767–771`).
- **Missing observation:** operation identity and deadline origin are absent from `Failure` and from every publication site.
- **Evidence required:** an observation must carry op, phase, deadline source, and instance identity before severity is assignable. A line number alone cannot distinguish handshake-fatal from request-diagnostic.

### R02 — Fatality is not preserved

- **Source behavior:** `failAll` stores only the first failure (`:373–376`) and delivers that same failure to every pending request (`:383–385`). `Synthesize` then publishes `TTSStateDegraded` for any `outcome.failure` without checking its class (`:796`).
- **Missing observation:** severity does not travel with the failure to its publisher.
- **Evidence required:** class and severity must be part of the delivered record, and `degraded` must be restricted to diagnostic failures. A fatal `protocol_violation` must not emerge as `degraded`.

### R03 — Fatal observations do not trigger enforcement

- **Source behavior:** the reader terminates only for `invalid_frame`, `frame_too_large`, or `protocol_violation` (`:304–306`). A stream close reclassified as `helper_crashed` (`:298–303`) does not terminate. A send failure in `callOn` returns `helper_crashed` without terminating (`:410–412`).
- **Missing observation:** recording a fatal failure and enforcing it are separate facts; only the first is recorded.
- **Evidence required:** every fatal observation needs an explicit enforcement outcome (action attempted, action result), or an explicit `not_enforced` disposition. Recording alone is insufficient.

### R04 — Reader completion is unobservable

- **Source behavior:** `waitForReader` (`:349–361`) has no return value; it returns silently on timeout. `Start` discards that silence (`:253`).
- **Missing observation:** whether the reader joined, and whether the wait timed out.
- **Evidence required:** a joined/timed-out result tied to the session and the termination action.

### R05 — Evidence-loss inventory was incomplete

- **Source behavior:** the three new losses are E7 (late terminal dropped), E8 (cleanup and cancellation results ignored), and E9 (shutdown state cleared after failed termination), traced in §5.
- **Missing observation:** the original six findings did not cover post-retirement terminals, action results, or post-failure state clearing.
- **Evidence required:** a loss register derived from every path that discards a value, not from the paths already suspected.

### R06 — Termination causes are conflated

- **Source behavior:** `Terminate` kills on `ctx.Done()` and returns `ctx.Err()` (`:274–276`), and kills on grace expiry then returns a `timeout` failure (`:277–286`). `Kill` is `Terminate` with the default grace (`:291–293`).
- **Missing observation:** neither branch records which cause fired, when, or with what grace.
- **Evidence required:** action records with cause, grace, start, completion, and result. Context cancellation and grace expiry are different outcomes and must not share a record shape without those fields.

### R07 — Rows combine incompatible phases

- **Source behavior:** P17 `O15` combined three handshake codes; `O16` collapsed nine `validateHello` failure exits (`protocol.go:502–545`); `O17` combined three load sites (`:581,586,590`); `O19` combined ten publication sites.
- **Missing observation:** phase and branch identity were lost at row level.
- **Evidence required:** one row per phase-and-branch emission, as expanded in §6.

### R08 — Completeness was overstated

- **Source behavior:** P17 asserted 20 observations as if they were 20 sites. At branch granularity the same rows expand to 45 emission sites (§6). The helper is untracked, so HEAD alone does not pin it.
- **Missing observation:** a site manifest and a source digest.
- **Evidence required:** the manifest in §6 plus the digest in §1, with the explicit statement that completeness is bounded to the audited tree and is not proven.

## 4. R09–R13 — model deficiencies

These identify missing definitions and compatibility decisions. No replacement mechanism is asserted as established.

### R09 — An append-only log is underspecified

- **Missing definition:** instance identity; sequence allocation across concurrent `failAll` and `SetState` writers; causal linkage between observation, action, and exit; action results; phase context; retention and drop policy.
- **Why it matters:** without these, appending preserves entries but not a reliable history. The log can be complete and still ambiguous.
- **Decision required:** define the record schema and allocator, or state that the log is diagnostic only and cannot support projection.

### R10 — Precedence is selection, not resolution

- **Missing definition:** an evidence rule that distinguishes correlation from causation. P17 §7 orders classes, but a protocol violation observed before a zero exit does not prove the violation caused the exit.
- **Why it matters:** class order is a policy choice. Presenting it as resolution would assert causality the source cannot show.
- **Decision required:** either define what independent evidence establishes causation, or keep the combination unresolved. No precedence level is proposed here.

### R11 — `quarantined` is unreachable

- **Missing definition:** no rule in P17 produces `quarantined`. The projection table omits it while the decision record defines it as an ownership condition (`DECISION…:33–41`).
- **Why it matters:** the natural path is a terminate timeout (`supervisor.go:284–285`) where the process may still be owned, but P17 neither projects nor records that state.
- **Decision required:** an explicit, amendment-gated projection for quarantine, or removal of `quarantined` from P17's projection claim. The author does not choose here.

### R12 — Handshake classification remains unresolved

- **Missing decision:** `validateHello` has nine failure exits, all `protocol_violation` (`protocol.go:502–545`), which `Start` maps to `unavailable` (`client.go:252–258`); a helper-declared incompatibility maps to `incompatible`. A missing mandatory capability is therefore projected as `unavailable`, not `incompatible`.
- **Why it matters:** the same compatibility concept has two projections, and one is arguably wrong. This is U8 and stays open.
- **Decision required:** a compatibility classification rule, not a code-name lookup.

### R13 — Provider-state compatibility is unresolved

- **Missing decision:** `Synthesize` publishes `TTSStateCrashed` on a write failure (`client.go:756`) while the process may still be alive, and publishes `Degraded` for a fatal pending failure (`:796`). The decision record separates provider state from disposition and forbids deriving one from the other (`DECISION…:43–47`).
- **Why it matters:** these writes conflate condition, disposition, and terminal evidence.
- **Decision required:** a compatibility decision for each write, or explicit recognition that the current writes are non-conforming. No change is proposed as fact.

## 5. Evidence-loss register E1–E9

| ID | Site | Loss | Consequence |
|---|---|---|---|
| E1 | `client.go:373–376` | only the first read failure is stored | later failures are discarded; all later calls return the first |
| E2 | `supervisor.go:233–243` | `lastFailure` overwritten | only the last publication survives |
| E3 | `supervisor.go:176–182` | nonzero exit in stopping branch not classified | unexpected failure appears as clean `absent` |
| E4 | `supervisor.go:246,291–293` | no intent; `Kill` aliases `Terminate` | ten actions are indistinguishable |
| E5 | `client.go:349–361` | reader join result discarded | timeout and join are indistinguishable |
| E6 | `client.go:298–306` | stream close marked crashed without kill | fatal recorded, process possibly alive |
| E7 | `client.go:339–343` | late terminal returns true with no record | comment claims evidence; nothing is stored |
| E8 | `client.go:248,654,661,664,861,871` | terminate and cancel results discarded | failed cleanup is indistinguishable from success; output is still deleted |
| E9 | `client.go:668–674` | shutdown clears `started`/`transport` regardless of terminate error | a timed-out process can remain owned while the client reports not started |

E7 is the sharpest: the code comment says "A late terminal is evidence only" but the branch neither stores nor publishes it. E8 covers five ignored `Terminate` calls and one ignored `Cancel`. E9 is the post-failure state clearing the review identified.

## 6. Inventory reconciliation at site and branch level

| Row | Phase | Contained emission sites | Note |
|---|---|---|---|
| O1 | discovery | 3 | `discovery.go:39,42,61` |
| O2 | start | 1 | factory failure |
| O3 | exit | 1 | unexpected-exit branch |
| O4 | exit | 1 | stopping, zero |
| O5 | exit | 1 | stopping, nonzero |
| O6 | transport/exit | 1 | stream close reclassified |
| O7 | protocol | 1 | three code values, one branch |
| O8 | protocol | 1 | decode failure |
| O9 | protocol | 1 | empty ID |
| O10 | protocol | 1 | duplicate terminal |
| O11 | protocol | 1 | unknown ID |
| O12 | transport | 2 | `callOn` send, `Synthesize` send |
| O13 | request | 2 | `callOn`, `Synthesize` |
| O14 | request | 2 | `callOn`, `Synthesize` |
| O15 | handshake | 1 | three outcome values, one branch |
| O16 | handshake | 9 | nine `validateHello` failure exits |
| O17 | load | 3 | call, decode, validation |
| O18 | transport | 1 | write failure with provider crash |
| O19 | request | 10 | publication sites `771–839` |
| O20 | shutdown | 2 | ctx and grace branches |

**Total: 45 emission sites across 20 rows.** Code variants and outcome values are not multiplied; this is a minimum branch-granularity count, not a completeness proof.

**Termination reconciliation:** ten actions (T1–T10) reduce to two mechanisms (`Terminate`; `Kill` as its default-grace alias) and nine call sites. `Terminate` has three select outcomes and three post-kill outcomes (`supervisor.go:271–286`). Intent is recorded at none of them. Row counts therefore overstated mechanism diversity while understating branch diversity.

## 7. U1–U8 remain open

| ID | Case | Status |
|---|---|---|
| U1 | orderly EOF vs reader `helper_crashed` | open; no resolution demonstrated |
| U2 | shutdown timeout racing protocol violation | open; two outcomes, no winner |
| U3 | incompatible handshake followed by crash | open; ordering is not causation |
| U4 | zero exit after fatal protocol violation | open; fatal-wins would overstate certainty |
| U5 | nonzero exit after deliberate kill | open; kill does not establish crash |
| U6 | provider `crashed` while process alive | open; dimensions conflated |
| U7 | terminal published before terminate timeout | open; immutability undefined |
| U8 | two handshake projections for one concept | open; compatibility rule missing |

No precedence assignment in this reconciliation closes any of them.

## 8. What this reconciliation does not claim

- It does not approve P17, revise it, or stand in for reviewer confirmation.
- It does not assert a new evidence mechanism, log schema, or projection.
- It does not claim the inventory is complete beyond the audited tree.
- It does not close R01–R13 or U1–U8.
- It does not authorize rev4, P18, implementation, tests, or builds.
- It does not treat the content digest as a Git identity or a substitute for tracking the helper.

## 9. Status

Review preserved unchanged. Reconciliation complete at finding level. P17 remains **NOT APPROVED**. Contract amendment and C1 implementation remain **UNAUTHORIZED**. Frozen artifacts untouched. No tests or builds were run.
