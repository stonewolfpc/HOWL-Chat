# C1 — Owner Decision: Lifecycle Disposition vs Provider State

**Date:** 2026-11-09
**Status:** OWNER DECISION RECORDED. **Not a contract amendment.**
**Decision owner:** repository owner. **Author role:** proposal and recording only.
**Resolves:** the P18 taxonomy question in reconciliation `28615d9d…`.
**C1 REMAINS UNAUTHORIZED.**

## 1. Decision

1. The thirteen-value provider state set stays **frozen**: `absent`, `starting`, `handshaking`, `unloaded`, `loading`, `ready`, `busy`, `degraded`, `stopping`, `crashed`, `incompatible`, `unavailable`, `backoff`.
2. Provider state and process-lifecycle disposition are **two dimensions**, not one enum.
3. `unknown` and `quarantined` are **not** provider states and are not added to that contract.
4. The lifecycle dimension is internal, conceptually `ProcessDisposition { active, reaping, reaped, quarantined, outcome_unknown }`. These names are illustrative, not an implementation instruction.
5. `TerminalEvidence` is a third, separate record: exit status, observed failures, termination actions, and causality.
6. Public exposure of any new distinction requires a **separately approved contract amendment** with a defined projection. It cannot appear implicitly.

## 2. Why the dimensions are separate

| Concept | Describes | Owner |
|---|---|---|
| `unknown` | terminal outcome cannot be established with sufficient evidence | terminal evidence / lifecycle projection |
| `quarantined` | supervisor still owns a process whose reaping or cleanup is unconfirmed | supervisor lifecycle / resource ownership |
| provider state | what the provider reports about operational condition | existing frozen contract |

`unknown` is epistemic. `quarantined` is an ownership condition. Neither states provider availability, so neither belongs in the availability enum. Expanding the enum to hold them would conflate three questions that fail independently.

## 3. Prohibited mappings

Silent substitution is forbidden because it discards the distinction the model exists to preserve:

| Forbidden | Why |
|---|---|
| `quarantined → unavailable` | reports a provider condition where the fact is unconfirmed reaping ownership |
| `unknown → crashed` | asserts a failure mode that was not established |
| `unknown → absent` | asserts a clean outcome that was not established |

The prohibition is on **deriving** provider state from disposition. A provider state that was already justified earlier (for example, `unavailable` recorded when startup failed) is not a mapping and stays as it is.

## 4. Exposure and amendment gate

1. Preserve the thirteen provider states.
2. Keep process ownership and terminal certainty in separate internal fields.
3. Define how those fields are exposed to consumers, if at all.
4. Require a separately approved contract amendment before adding public provider states.

Until step 4 is approved: no terminal provider state is written from a disposition; the provider state remains at its last operationally justified value, and the disposition is internal diagnostic evidence. An internal distinction does not become a public one by accident.

## 5. Effect on open mechanisms

| Mechanism | Effect of this decision | Remaining work |
|---|---|---|
| P18 | **decided** — no enum change; internal dimension | rev4 must name the fields and transitions |
| P6 | unblocked — quarantine has an internal disposition | no provider state; specify publication path |
| P13 | zero exit without an action → `outcome_unknown`, not `absent` | P17 severity/class |
| P14 | ambiguous nonzero → `outcome_unknown`, not `crashed` | P17 |
| P15 | handshake fatal → `incompatible` (existing state is legitimate); disposition `reaped` | post-handshake class case |
| P16 | a kill explains the exit; it does not become the crash | P17 precedence for pre-existing failure |
| P17 | unchanged — still open | severity, class, precedence for `failureObserved` |
| P8 / P12 | terminal evidence publishes internally; public projection gated by §4 | terminal latch and post-closure buffer |

R3 IDs remain open individually; this decision removes the taxonomy blocker, not the findings.

## 6. What this decision does not do

- It does not add, rename, or remove any provider state.
- It does not authorize a UI, API, or consumer-visible change.
- It does not close any R3 finding, W1–W4, X-1, or X-2.
- It does not authorize rev4, implementation, tests, or builds.
- The author cannot amend the frozen contract; only the owner can, by separate decision.

## 7. Status

Dependency order step 1 (taxonomy) is complete at decision level. Next gate is P17 — severity, class, and precedence — before any rev4 text. No matrix regeneration, no implementation, no tests. Frozen evidence unchanged. **C1 remains unauthorized.**
