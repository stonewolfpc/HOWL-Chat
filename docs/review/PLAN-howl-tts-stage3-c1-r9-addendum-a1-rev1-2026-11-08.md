# C1 Revision 9 — Addendum A1, Revision 1

**Date:** 2026-11-08
**Status:** CORRECTION PACKAGE. Not r9 design text. **C1 REMAINS UNAUTHORIZED.**
**Supersedes:** `PLAN-howl-tts-stage3-c1-r9-correction-addendum-a1-2026-11-08.md`, sha256 `92acde5589df5e1808d236e2926060830285806bc5155433c1569829a0cfb93f` (300 lines, 18,074 B). The original is preserved unmodified and remains the record of what was proposed.
**Companion:** `...addendum-a2-rev1-2026-11-08.md` (rebuilds inventory, state model, availability writers, matrix sizing).
**Unchanged inputs:** r8 design `a8c40fac7d86f3cd34ad34e7f656668684830e76dfb605ae6817fa7119638e06`; r9 plan `998974229e74d8e20a048118d54126a4d488b3c827b8c7745663c8ba76522354`; A2 review `0a5ea0479651e60bf48589dac912b50d01b9159f8af23ffa8ef9d13637b71ea0`.

Finding classes: **F** source-verified here, **R** reported-unverified, **P** author proposal, **O** owner decision.

## 1. What revision 1 changes, and what it does not

A1 proposed mechanisms and then treated several of them as guarantees. This revision keeps the direction — installation-bound authority, record-reference capability — and removes the guarantees that were not earned.

Three substantive changes:

1. **The authority becomes supervisor-targeted.** Authority-wide pair identity is insufficient; a target identity and generation are now first class (§3).
2. **The lifecycle gets an explicit failure taxonomy and lock order** rather than a matrix of intentions (§4).
3. **The availability claim is withdrawn from A1 entirely** and rebuilt in A2-rev1 §6 on the actual transition table. A1's §5 assertion that the availability model was correct was unearned: the current source contradicts the contract at three points (§2, A1-05).

What A1 got right and is retained: the objection that a package boundary is not a security boundary (P), the split `Verifier`/`Updater` pair (P), pre/post digest deltas as an issuance precondition (P), the record-reference capability (P), and the bounded same-package selection residual (P, O-pending).

## 2. Reconciliation — A1 findings

No finding disappears because its wording changes. A dispute stays visible.

| ID | Finding | Evidence | Required correction | Author response | Independent disposition |
|---|---|---|---|---|---|
| **A1-01** | Authority is not supervisor-targeted; a shared verifier with an authority-wide registry cannot distinguish two supervisors with identical binding fields, and one-use consumption does not fix targeting | A1 §3.4 test A6 asserts independent per-supervisor consumption; A1 §3.1 defines validity as membership in one authority registry; no target field exists in `Ticket` | add explicit target identity + generation, keyed registry | §3 | Pending |
| **A1-02** | Lock order between record mutex and in-flight registry undefined; concurrency protections can deadlock or retire inconsistently | A1 §4.2 lists concurrency cases, no ordering | define non-nesting acquisition rule | §4.4 | Pending |
| **A1-03** | Digest delta proves change, not intent | property of digest comparison | preserve as permanent limitation, never phrase as "verified completion" | §5.1 | Pending |
| **A1-04** | Record reference is not proof the producing operation was authorised; creation, registration, validation must be inside the trusted operation | A1 §4.1 treats the reference as sufficient | registry liveness is the authorisation; single-invocation creation | §4.1 | Pending |
| **A1-05** | Availability writer model conflicts with the contract; implementing it as written can leave stale or wrong availability | `supervisor.go:166–191` `observeExit`: on any `stopping` exit provider state is overwritten with `absent`; nonzero exit therefore becomes `absent`; `incompatible` cannot survive termination. `supervisor.go:325` `providerAvailable` returns false for `unavailable`; `supervisor.go:335` `processAlive(process, provider)` derives liveness from aggregate state | rebuild as derived-at-read over a stored class, with an intent-aware exit observer | A2-rev1 §6 | Pending |
| **A1-06** | Consumer audit for `TTSProviderStatus.Available` absent | A1 §5.3 asserted consumer-inertness without enumeration | enumerate every consumer and what `false` means | A2-rev1 §6.1 (retained from A2 §5) | Pending |
| **A1-07** | A2-D6 fields (`Loaded`, `ServiceCapabilities`, `ServiceCapabilitiesKnown`) need explicit disposition; must not be absorbed into the availability fix | `types/provider.go:56–71` lacks all three; contract requires them | separate criterion, separate verification | A2-rev1 §6.5 | Pending |
| **A1-08** | Verifier-binding tests absent (two pairs, cross-pair, wrong verifier, missing/invalid, multiple supervisors, concurrent completion+consumption) | A1 §3.4 had no binding tests | ten named tests incl. targeting | §3.6 | Pending |

## 3. Supervisor-targeted authority (P)

### 3.1 Why pair identity was not enough

Pair identity answers "did this ticket come from the authority I installed?" It does not answer "was this ticket issued for the supervisor that will consume it?" Two supervisors created with identical binding fields — same helper, same provider, same config — are indistinguishable to an authority-wide registry, so a ticket issued for one is consumable by the other. One-use consumption limits the damage to one wrong target; it does not prevent it.

### 3.2 Target identity

Each supervisor instance owns a target identity:

```text
TargetID  = { instanceNonce, generation }
```

- `instanceNonce` is minted at supervisor construction and never reused.
- `generation` increments on identity replacement (A2-D7) and on update-gated reset.

`Updater.Begin` requires a `TargetID` obtained from the supervisor being updated. `Ticket` carries it. `Verifier.Consume` rejects a ticket whose `TargetID` differs from the verifier's installed target (`target_mismatch`), and rejects a stale generation (`generation_mismatch`).

Consequence: a shared verifier is safe, because consumption is keyed by target, not by authority. Two supervisors share nothing but the issuer.

`PairID` equality remains, demoted: it is a **wiring sanity check** that catches C7 and supervisor being constructed from different authorities at startup. It is not the targeting control, and A2-rev1 §2 records that A2 overstated it.

### 3.3 Failure taxonomy (P)

| Code | Meaning | Reset occurs |
|---|---|---|
| `ticket_unknown` | not issued by this authority | no |
| `ticket_expired` | TTL elapsed | no |
| `target_mismatch` | issued for another supervisor | no |
| `generation_mismatch` | supervisor replaced since issuance | no |
| `already_consumed` | single use exhausted | no |
| `authority_closed` | authority torn down | no |
| `evidence_unchanged` | binary and manifest unchanged since `Begin` | no |
| `evidence_partial` | one of binary/manifest/config unreadable | no |

Every entry is fail-closed. None of them mutates supervisor state.

### 3.4 Required tests

| # | Scenario | Required result |
|---|---|---|
| A1 | two independently constructed pairs | each verifier accepts only its own tickets |
| A2 | ticket from pair A to verifier B | `ticket_unknown`, no reset |
| A3 | wrong verifier installed | startup assertion fires; no supervisor |
| A4 | missing verifier | every reset rejected; no panic |
| A5 | invalid `cap.New` config | constructor error; no start |
| A6 | **two supervisors, one authority, identical binding fields** | ticket for S1 rejected by S2 (`target_mismatch`); S1's consumption unaffected |
| A7 | authority torn down while supervisor live | `authority_closed`; no reset |
| A8 | concurrent `Complete` + `Consume` | exactly one consumption; no race |
| A9 | ticket issued before generation bump, consumed after | `generation_mismatch` |
| A10 | digest delta unchanged between begin and issue | `evidence_unchanged`; no ticket |

A6 is the finding A1-01 correction; A9/A10 close the two largest residual gaps in the previous issuance rule.

## 4. Capability lifecycle (P)

### 4.1 Membership is the authorisation

A record is created, validated, and registered by one invocation of the authorised operation. Apply accepts only records that are **currently registered and live**. A record that exists but was never registered is not a capability; a retired record is not a capability. The reference carries nothing; the registry decides.

### 4.2 Consumption is single-use and burns on failure

One record → at most one apply. A failed apply burns the record rather than leaving it retryable, because a half-applied identity change is worse than a dropped one. Caller retries re-enter the operation and get a new record with a new generation.

### 4.3 What the capability may not do

It may not: start or stop a process, load or unload a provider, change selected provider identity, alter availability or lifecycle state, consume a request from the pending budget, or cross a supervisor boundary. Its only effect is the one apply it was created for.

### 4.4 Lock order — sequential, never nested

```text
1. Client.mu      → read liveness + target/generation snapshot, release
2. record.mu      → claim single-use, release
3. Supervisor.mu  → apply, release
```

- No lock in this set is held while acquiring another.
- `record.mu` is never held across a transport call or supervisor call.
- `Supervisor.mu` is never acquired while `Client.mu` or `record.mu` is held.

The race this admits is explicit: a record can be retired between steps 1 and 2. Step 2 then fails closed and the apply does not happen. That direction is safe; the reverse (applying to a retired record) is what the rule forbids.

### 4.5 Lifecycle failure semantics

| Event | Record state | Visible effect |
|---|---|---|
| begin | `live` | none |
| validation fails | `aborted` | typed error, no record published |
| issue | `issued` | ticket; no state change |
| consume ok | `consumed` | one apply |
| consume fails | `burned` | none |
| operation cancelled | `retired` | none |
| generation change | `stale` | later consumes rejected |
| supervisor teardown | `retired` | later consumes rejected |

### 4.6 Bounded selection residual (unchanged, P)

Same-package code can still construct an identity-only selection: generation-bound, single-use, no lifecycle effect, no budget effect. Accepting or removing this residual is an owner decision (O-5).

## 5. Preserved limitations

### 5.1 Digest delta

> Requiring a binary or manifest digest to differ proves that **something** changed. It does not prove that the **intended** update occurred.

C7 could perform a real but unintended file change and still obtain a ticket. Issuance may be described only as "verified as an observable state change".

### 5.2 C7 trust

C7 remains a trusted reporter with digest compensation. This is an owner decision, not a derived property.

### 5.3 Authority boundary

Installation-bound authority constrains composition. It does not constrain code already inside the enforcing component, and it does not survive `unsafe` or reflection. Stated so no later document claims more.

## 6. Owner decisions still required

| # | Decision | Recommended |
|---|---|---|
| O-1 | installation-bound authority, now target-keyed | accept |
| O-2 | C7 trusted reporter + digest compensation | accept |
| O-3 | burn-on-failure consumption (§4.2) | accept |
| O-4 | target identity in `Ticket` | accept |
| O-5 | bounded same-package selection residual | accept |
| O-6 | availability class internal-only | accept (A2-rev1 §6) |
| O-7 | A2-D6 three-field gap as its own criterion | accept, separate |
| O-8 | R-38/R-39 original text | required |

## 7. Authorization

**C1 remains unauthorized.** This revision edits no design document, no source, no test, no dependency, and no configuration. Nothing was built or run. The 142 pre-existing dirty entries are untouched.
