# C1 — Review Record: Contract rev3

**Date recorded:** 2026-11-09
**Status:** RECORD OF A REVIEW SUPPLIED BY THE OWNER. **Not a signed attestation.**
**Provenance:** the review report was supplied to the owner; its text and its `R3-01`…`R3-21` IDs are **not in the workspace**. This record captures the owner's summary only. Reviewer identity, date, and method were not stated.
**Attestation:** remains unsigned; not populated from these findings.

## 1. Verdict

**NOT READY.** The review identifies contradictions between rev3's mechanisms and rev3's own invariants, not merely missing explanation.

## 2. Findings as received

### R3-01…R3-07 — process and pipe ownership (connected lifecycle problem)

- A refused `pending` process has no legal transition into reaping.
- The forced-drain path does not explicitly join the drain goroutines before `Wait`.
- Natural EOF and forced closure are not distinguishable.
- Stdin closure has no consistent owner in the proposed interfaces.
- A3's "Wait again" conflicts with the one-owner invariant.
- Quarantine has no explicit terminal projection.

**Key correction:** a timeout must never transfer or duplicate ownership of `cmd.Wait()`. One owner retains the operation; other participants wait on its completion signal. Quarantine retains that ownership even after the caller's deadline expires.

### R3-08…R3-11 — publication and attachment

- The pre-retirement rule cannot simultaneously permit delayed delivery and prohibit post-retirement delivery without stating which wins.
- A terminal projection needs an explicit authorization path.
- Ordinary publications must not overwrite a terminal projection.
- W1/W3 cannot be considered closed: the interface split is useful, but the registration and attachment capability that gives those interfaces lifecycle authority has not been restated.

### R3-12…R3-15 — transport arbitration

- **Acknowledgment deadlock:** a shutdown sender holding the lease while awaiting an acknowledgment that also needs the lease deadlocks.
- The lease needs separate treatment of frame writing, acknowledgment delivery, ordinary admission, shutdown priority, expiry/cancellation, late completion and late acknowledgment, transport closure, and evidence recording after closure.
- One arbiter is necessary, but requiring everyone to use it does not establish a correct policy.

### R3-16…R3-21 — terminal evidence

- Exit code zero alone does not establish an expected exit.
- A nonzero exit with ambiguous causality does not establish a crash.
- `protocol_incompatible` is not equivalent to a process crash.
- A deliberate kill is not automatically an unexpected crash.
- Required distinct projections: `absent`, `crashed`, `incompatible`, `unknown`, `quarantined`.
- The contract must define fatal versus diagnostic evidence, precedence among multiple failures, and the effect of each termination action.

## 3. Disposition summary

| Finding | Assessment |
|---|---|
| R3-01–R3-07 | substantive; revise |
| R3-08–R3-11 | substantive; revise |
| R3-12–R3-15 | substantive; revise |
| R3-16–R3-21 | substantive; revise |
| W1–W4 | remain open |
| X-1 / X-2 | remain open |
| R2-01…R2-05 | reserved; source text unavailable |
| Rev2 pipe rule | withdrawal accepted; replacement not resolved |
| Crash-before-request | open |
| Lease-expiry winner | open |
| C1 | unauthorized |

## 4. Governance

- `R3-01`…`R3-21` are preserved individually. They are not collapsed into "pipe fixes", "lifecycle fixes", or any other bucket; each must be shown addressed on its own.
- Rev4 does **not** start by editing the contract. The next artifact is a finding-to-mechanism reconciliation using real source paths and existing contract sections, exposing dependencies between corrections.
- An independent reviewer examines the proposed mechanisms before a new contract is written.
- No matrix regeneration, implementation, tests, or builds. Attestation unsigned. Frozen evidence unchanged. **C1 remains unauthorized.**
