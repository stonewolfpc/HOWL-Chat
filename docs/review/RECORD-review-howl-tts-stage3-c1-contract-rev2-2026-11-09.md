# C1 — Review Record: Contract rev2 (Installation Lifecycle)

**Date recorded:** 2026-11-09
**Status:** RECORD OF A REVIEW SUPPLIED BY THE OWNER. **Not a signed attestation.**
**Provenance:** findings and dispositions received in the owner's message; the review text itself is **not in the workspace** (no `Pasted*` file, no artifact containing `R2-0x`). The owner states the conclusions rest on the supplied review and its reported source inspection, **not** an independent inspection of this repository.
**Attestation:** remains unsigned; not populated from these findings.

## 1. Verdict

**NOT READY.** Rev2's ownership direction survives; several of its safety mechanisms can be bypassed, and one rule contradicts the documented `os/exec` contract.

## 2. Findings as received

### W2 — strongest finding

**Rev2's pipe-ordering rule is rejected.** `Wait` cannot be the EOF synchronization point for `StdoutPipe` and `StderrPipe`. The documented `os/exec` contract requires reads to complete before `Wait`, because `Wait` closes the pipes after seeing the command exit.

Three separate responsibilities:

1. reading stdout and stderr to completion;
2. reaping through exactly one `Wait` owner;
3. releasing remaining handles and reporting cleanup failures.

Also confirmed: the existing `Close()` closes `stderr` while its drain may still run, violating the ownership model.

Corrections supported: split the transport-facing interface from the lifecycle-facing `Process` interface so the client cannot call `Wait`, `Kill`, or `Close`; and quarantine must retain an actual reaping owner — a `sync.Once` prevents a second `Wait` but does not establish ownership or prevent an unsafe `Kill` afterward.

Refinement: specify pipe order against the real `exec.Cmd` implementation and process abstraction, not as a universal ordering for every possible process.

### W3 — publication is broader than `readLoop`

Handshake results, request delivery, failure delivery, and supervisor state updates all need lifecycle protection.

Required protocol, refined away from holding `s.mu` during delivery:

1. acquire publication authorization under the lifecycle lock;
2. register an in-flight publication;
3. release the lock;
4. complete or cancel the publication;
5. retirement closes admission to new publications and waits for or invalidates outstanding ones by a defined policy.

Abandonment transitions are required: a registration that never becomes a running reader must not stay authoritative indefinitely.

### W4 — the lease must govern every transport writer

A supervisor-side lease cannot protect a transport that client code can write to or close independently. The lease must cover ordinary request writes, the shutdown frame, transport closure, acknowledgment processing, expiry and cancellation, and late completions. Sender and observer need one arbitration mechanism, or the lease is advisory.

No-transport shutdown must be a first-class outcome, not a path that silently falls back to `Start`.

### Terminal projection — sequence is not causality

Sequence numbers order observations, not events. The crash-before-request case shows the difference. Recording an action before `Kill` and an exit at `Wait` is evidence, not proof the process was alive when the action began.

The next contract must distinguish actual event time where observable, observation time, action initiation, action completion, and unresolvable ambiguity. Where causality cannot be established it must preserve the uncertainty rather than manufacture a classification.

`failureObserved` is required: a reader-observed fatal protocol violation must survive a clean subsequent exit.

### X-1, X-2, and the three original open items

Confirmed open, and all three rev2 open items remain open.

## 3. Disposition summary

| Finding group | Disposition |
|---|---|
| W1 | direction accepted; existing-process/single-attachment API incomplete |
| W2 | pipe-order rule rejected; lifecycle interface and reaper ownership need revision |
| W3 | registration model accepted directionally; publication and abandonment incomplete |
| W4 | transaction accepted directionally; transport arbitration incomplete |
| Terminal projection | open; causality and fatal-failure recording incomplete |
| X-1 | confirmed, open |
| X-2 | confirmed, open |
| Three original open items | all remain open |

## 4. ID preservation

`R2-01`…`R2-05` were accepted by the owner, but their text is not in the workspace. They are **reserved**, not reinterpreted. The review is reported to contain six requested corrections; the sixth is not identified here and is not guessed.

Six revision areas, kept separate rather than collapsed:

1. process ownership — restricted transport interface, one reaping owner, no kill after reap;
2. pipe lifecycle — stdin closure, stdout/stderr draining, reaping, final release;
3. quarantine — bounded ownership and escalation without admission locks;
4. publication — authorization, retirement, abandonment, rotation;
5. transport lease — all writers, shutdown, closure, expiry, late acknowledgments;
6. terminal evidence — observed failures, actions, outcomes, unresolved causality.

## 5. Standing constraints

No r9 edit, no matrix regeneration, no implementation, no tests, no builds. Attestation unsigned. Frozen evidence unchanged. **C1 remains unauthorized.**
