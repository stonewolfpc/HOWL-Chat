# MEMORY-09 — Future integration testing note: Howl Chat, post-experiment

**Status:** **OWNER PLANNING NOTE.** It records a future phase and its boundaries. It is **not a
decision**, **not a gate**, and **not an authorisation**. It changes no ratified artifact.

**Recorded:** **2026-11-17 08:28:48 CST (-0600)** = **2026-11-17 14:28:48 UTC**
**Basis:** owner statement, 2026-11-17, reproduced verbatim in §1.
**Affects:** nothing currently in force. The architecture (57 rules), the contract (r32), the
conformance requirements (118) and the experiment plan (v3) are unchanged by this record.

---

## 1. Owner's note (recorded verbatim)

> ## Future Howl Chat integration testing — post-experiment
>
> MEMORY-09's current experiments are deliberately conducted outside Howl Chat, using Kobold and
> Python as an independent experimental environment.
>
> This separation is intentional.
>
> Howl Chat is not yet sufficiently complete to provide a reliable environment for testing
> MEMORY-09's underlying mechanics. Completing additional application functionality merely to reach
> those mechanics would unnecessarily couple experimental validation to unfinished product
> development.
>
> The external environment provides a controlled playground for testing and evaluating memory
> behaviour independently of Howl Chat.
>
> **After the experiment has run and the memory architecture has been evaluated and explicitly
> decided for adoption, a separate phase of integration testing must take place inside Howl Chat.**
>
> That future phase will establish whether the adopted memory architecture integrates correctly with
> Howl's actual application state, chat ownership, worlds, scenarios, characters, lore, and
> operational workflows.
>
> ### Boundaries
>
> - External experiments validate the memory mechanics independently of Howl Chat.
> - External experimental success does not establish successful Howl integration.
> - Howl integration testing begins only after the experiment and an explicit architecture adoption
>   decision.
> - Integration testing is a separate phase, with its own scope and authorization.
> - This note does not authorize implementation, production integration, or changes to the current
>   experiment.
> - MEMORY-09's existing scope and ratified requirements remain unchanged.
>
> The purpose is to avoid finishing an application layer prematurely just to test an experimental
> subsystem, while ensuring that successful experimental results are not mistaken for proof of
> production integration.

## 2. Why the separation exists

The experiment tests memory mechanics. Howl Chat is an application whose state, ownership model,
worlds, scenarios, characters and lore are still being built. Testing an experimental subsystem
inside an unfinished host would make the experiment's validity depend on the host's progress, and
would make the host's progress depend on the experiment. Neither dependency is wanted.

The external environment (Kobold and Python) is therefore not a convenience — it is the **controlled
playground** in which the mechanics are validated on their own terms.

## 3. Standing inference limits

Two inferences this note forbids:

| Forbidden inference | Why |
|---|---|
| "The experiment passed, therefore Howl integration works" | External success establishes mechanics only. Integration is untested by construction |
| "Howl must add functionality so the experiment can be tested in place" | That inverts the dependency the separation exists to prevent |

The second is the one that quietly re-couples the phases. An integration gap discovered early is
**not** a reason to move the experiment inside Howl.

## 4. Relationship to the current gates

This note names a **third, later phase**. It adds no gate today and authorises nothing.

| Phase | Question | Status |
|---|---|---|
| **Gate 1** — conformance (O-7 / Amendment 1) | May probes for the 118 ratified requirements be defined and run? | **Permitted**, prospective, probes only |
| **Gate 2** — implementation (`eval/mem1`) | May implementation code be written? | **UNANSWERED — not authorised** |
| **Phase 3** — Howl integration (this note) | Does the adopted architecture integrate with real application state? | **Future. Not started. Requires the experiment to run, an explicit adoption decision, and its own separate authorisation** |

Phase 3 has **three** preconditions, all unmet:

1. the experiment has run;
2. the memory architecture has been evaluated;
3. adoption has been **explicitly decided** by the owner.

A separate scope and a separate authorisation are then required. Nothing here schedules them or
implies they are approved.

## 5. What this note does not do

- It does **not** change the architecture (57 rules), the contract (r32), the conformance map
  (118 requirements) or the experiment plan (v3).
- It does **not** add a requirement, probe, rule or field.
- It does **not** authorise implementation, production integration, or any change to the current
  experiment.
- It does **not** claim the experiment will succeed, that adoption will follow, or that integration
  will be achievable.
- It does **not** reopen O-2, Gate 2, or any confirmed decision.

## 6. Verification performed when writing this record

Observed, not inferred:

| Check | Result |
|---|---|
| Architecture rules | **57** — unchanged |
| Conformance requirements | **118**, 12/12 group headings match — unchanged |
| Contract r32 | `e88ef2d1697da1d9…` — unchanged |
| Experiment plan v3 | 21,616 B, `6c4078bd…` — unchanged |
| Amendment 1 | `6afd2a0d8c19a5cd…` — unchanged |
| `eval/mem1` files modified on 2026-11-17 | **0** |
| Tests implemented / executed | **0 / 0** |
