# P17 — Evidence Semantics Proposal (Frozen)

**Date:** 2026-11-09
**Status:** AUTHOR PROPOSAL FOR INDEPENDENT REVIEW. **Frozen on write.** Not a contract edit, not rev4.
**Source:** `internal/audio/tts/helper/**` at HEAD `9c7d2a1`; read-only.
**Scope:** lifecycle-affecting failure and termination observations only. Request validation, output validation, and discovery are inventoried for classification but are not lifecycle transitions.
**C1 REMAINS UNAUTHORIZED.**

No code, test, build, or contract file was changed. Every classification below cites a path and line; nothing is inferred from a code name alone.

## 1. Constraints

1. All observations are preserved in an append-only log.
2. Precedence selects a projection; it never deletes evidence.
3. Severity is assigned by the observing site, not looked up from the failure code.
4. An action is not evidence, and an exit fact is not an action.
5. Where the source cannot establish a classification, the case stays unresolved.

## 2. Model

```text
Observation { id, seq, site, code, class, severity, detail }
Action      { id, seq, site, intent, grace, completion }
EvidenceLog = append-only [Observation | Action | ExitFact]
```

`Failure` today carries `Code`, `Message`, `Retryable`, `Field`, `Detail` (`errors.go:15–21`) and no sequence, class, or severity. The proposal adds those at the observing site, not inside `Failure`.

## 3. Evidence classes

| Class | Definition | Fatal capable |
|---|---|---|
| handshake | pre-hello negotiation or hello validation | yes |
| protocol | frame, envelope, or response rule violation after transport establishment | yes |
| transport | write failure, stream closure, receive failure | yes |
| exit | result of `Wait` | fact, not severity |
| timeout | request, shutdown, or reaping deadline expiry | action outcome |
| request | validation, busy, cancellation | no |
| output | artifact validation after synthesis | no, but may trigger cleanup |
| discovery | helper path or presence | no |

## 4. Severity rules

- **Fatal**: the observation invalidates the session or requires termination; no further request can succeed on that instance.
- **Diagnostic**: recorded, no lifecycle consequence by itself.
- **Exit fact**: records what happened, never a severity.
- **Action outcome**: records what was requested and whether it completed.

Rules:

| # | Rule |
|---|---|
| S1 | severity is bound at the site; `timeout` is diagnostic at `client.go:439` and an action outcome at `supervisor.go:285` |
| S2 | one code may appear at more than one severity |
| S3 | fatal is not synonymous with crash: handshake fatal ≠ crash |
| S4 | diagnostics never select a projection |
| S5 | an action never overrides a pre-existing fatal observation |
| S6 | provider state is not derived from disposition or from an unobserved outcome |

## 5. Observation inventory

| ID | Site | Code | Class | Severity | Current handling |
|---|---|---|---|---|---|
| O1 | `discovery.go:39,42,61` | helper_missing | discovery | diagnostic | `setUnavailable` (`supervisor.go:132`) |
| O2 | `supervisor.go:147` | helper_start_failed | discovery | diagnostic | `setUnavailable`, no process |
| O3 | `supervisor.go:185–190` | helper_crashed | exit | fact | provider crashed, process exited |
| O4 | `supervisor.go:178–180` | none | exit | fact | absent, provider absent |
| O5 | `supervisor.go:179–181` | none, err retained in `waitErr` | exit | fact | provider absent despite nonzero exit |
| O6 | `client.go:299–303` | internal_error → helper_crashed | transport | fatal | failAll; no kill for this code |
| O7 | `client.go:304–306` | invalid_frame, frame_too_large, protocol_violation | protocol/transport | fatal | failAll + Kill |
| O8 | `client.go:309–313` | decode errors | protocol | fatal | failAll + Kill |
| O9 | `client.go:323` | protocol_violation | protocol | fatal | protocolViolation + Kill |
| O10 | `client.go:331` | protocol_violation | protocol | fatal | protocolViolation + Kill |
| O11 | `client.go:345` | protocol_violation | protocol | fatal | protocolViolation + Kill |
| O12 | `client.go:421` | helper_crashed | transport | fatal | returned only; no termination |
| O13 | `client.go:439` | timeout | timeout | diagnostic | request fails |
| O14 | `client.go:441` | cancelled | request | diagnostic | request fails |
| O15 | `client.go:248–253` | incompatible, minor_unsupported, others | handshake | fatal | Terminate then publish incompatible or unavailable |
| O16 | `protocol.go:502–545` | protocol_violation | handshake | fatal | falls to `unavailable` at `client.go:253` |
| O17 | `client.go:578–594` | wire failures | handshake/load | diagnostic | provider unavailable |
| O18 | `client.go:756` | helper_crashed | transport | fatal | provider crashed while process may be alive |
| O19 | `client.go:771,783,793,796` | timeout, cancelled, degraded | request | diagnostic | provider ready or degraded |
| O20 | `supervisor.go:283–286` | timeout | timeout | action outcome | returned; process state unchanged |

## 6. Termination actions

| ID | Site | Intent (proposed) | Grace | Intent recorded today |
|---|---|---|---|---|
| T1 | `client.go:248` | handshakeFailure | ShutdownTimeout | no |
| T2 | `client.go:305` | protocol | default | no |
| T3 | `client.go:312` | protocol | default | no |
| T4 | `client.go:366` | protocol | default | no |
| T5 | `client.go:654` | shutdown | ShutdownTimeout | no |
| T6 | `client.go:661` | shutdown | ShutdownTimeout | no |
| T7 | `client.go:664` | shutdown | ShutdownTimeout | no |
| T8 | `client.go:668` | shutdown | ShutdownTimeout | no |
| T9 | `client.go:861` | cleanup | ShutdownTimeout | no |
| T10 | `supervisor.go:275,278` | shutdown-grace-expiry | remaining grace | no |

`Kill` is an alias for `Terminate` with the default grace (`supervisor.go:291–293`), so T2–T4 are indistinguishable from T10 in evidence.

## 7. Precedence

Selection order, applied to the preserved log:

1. fatal handshake observation → `incompatible` for `protocol_incompatible`/`protocol_minor_unsupported`; `unavailable` for handshake `protocol_violation`.
2. fatal protocol observation → `crashed` + `protocol_violation`.
3. fatal transport observation → `crashed` + `helper_crashed`.
4. exit fact with no fatal observation:
   - zero + completed shutdown action → `absent`;
   - zero + no action → `outcome_unknown`;
   - nonzero + known action cause → internal `terminated_by_action` (public projection undecided, see U5);
   - nonzero + ambiguous cause → `outcome_unknown`.
5. diagnostics never change the projection.

| Evidence combination | Projection |
|---|---|
| fatal handshake + any exit | incompatible or unavailable |
| fatal protocol + any exit | crashed + protocol_violation |
| fatal transport + any exit | crashed + helper_crashed |
| exit zero, shutdown completed | absent |
| exit zero, no action | outcome_unknown |
| exit nonzero, ambiguous | outcome_unknown |
| exit nonzero, action caused, no fatal | terminated_by_action (internal) |
| diagnostics only | unchanged |

## 8. Evidence-loss findings

| ID | Finding | Source |
|---|---|---|
| E1 | `failAll` keeps only the first read failure; later failures are discarded | `client.go:369–386` |
| E2 | `SetState` overwrites `lastFailure`, so only the last published failure survives | `supervisor.go:233–241` |
| E3 | the stopping branch discards a nonzero exit from classification | `supervisor.go:176–187` |
| E4 | termination intent is not a parameter, and `Kill` aliases `Terminate` | `supervisor.go:246,291–293` |
| E5 | `waitForReader` returns nothing on timeout | `client.go:352–367` |
| E6 | stream closure is marked crashed but does not terminate for that code | `client.go:299–308` |

These violate constraint 1 and must be fixed before precedence can be trusted.

## 9. Unresolved cases

| ID | Case | Why unresolved |
|---|---|---|
| U1 | orderly EOF during shutdown vs `helper_crashed` from the reader | the reader maps generic stream errors to crash while the observer publishes absent |
| U2 | shutdown timeout racing a protocol violation | two fatal/action outcomes with no defined winner |
| U3 | handshake incompatible followed by a process crash | first terminal vs sequence order |
| U4 | exit zero after a fatal protocol violation | fatal-wins rule may overstate certainty |
| U5 | nonzero exit after a deliberate kill | crash vs terminated_by_action is not established |
| U6 | provider `crashed` while the process is alive (`client.go:756`) | provider state conflates condition with process outcome |
| U7 | `Terminate` timeout after the terminal projection is published | terminal immutability is undefined |
| U8 | missing mandatory capability yields `protocol_violation` to `unavailable`, while helper-declared incompatibility yields `incompatible` | same concept, two projections |

No classification is asserted for U1–U8.

## 10. Reviewer challenges

1. Verify every line and code claim above against source.
2. Challenge each fatal assignment, especially O6, O12, O15, O16.
3. Test U1–U8, including shutdown timeout plus protocol violation and handshake incompatibility plus crash.
4. Check that precedence never deletes evidence and that no projection is unreachable.
5. Confirm the five required projections (`absent`, `crashed`, `incompatible`, `unknown`, `quarantined`) are reachable or declare them not.
6. Return finding-level dispositions; do not accept the design on the author's behalf.

## 11. Status

Frozen proposal. No contract edit, no rev4, no matrix, no implementation, no tests, nothing built or run. Frozen evidence unchanged. **C1 remains unauthorized.**
