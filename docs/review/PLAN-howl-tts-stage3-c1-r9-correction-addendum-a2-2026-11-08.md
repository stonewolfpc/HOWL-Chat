# C1 Revision 9 Correction Plan — Addendum A2

**Date:** 2026-11-08
**Status:** ADVERSARIAL-REVIEW INPUT. Not approval. **C1 REMAINS UNAUTHORIZED.**
**Addresses:** verifier-to-authority binding, capability lock order, availability consumer audit, A2-D6 field disposition, and R-44 matrix verification from source.
**Does not:** approve A1, edit any design document, or resolve R-38/R-39.

## 1. Claim labels

Every statement below carries **F** (source-verified here), **R** (reported, unverified), **P** (author proposal), or **O** (owner decision). The labels are load-bearing: nothing marked P is a fact, and nothing marked R is evidence.

| Item | Class |
|---|---|
| §3.1–§3.3 verifier binding mechanism | P |
| §3.4 failure modes and tests | P |
| §4 preserved limitation | F (a property of what a digest can prove) |
| §5 capability lifecycle and lock order | P |
| §6.1–§6.4 consumer audit | **F** |
| §6.5 A2-D6 field gap | **F** |
| §7 R-44 verification | **F** |
| R-38, R-39 | R (identifier only, no text) |

## 2. Verifier binding — what establishes the installed verifier belongs to the intended authority

### 2.1 The mechanism, stated precisely

Three properties, each mechanical rather than stipulated:

1. **Pair identity is intrinsic.** `cap.New` mints a pair ID and returns a `Verifier` and an `Updater` that share it in an unexported field. `Ticket` validity is registry membership in the authority the verifier points at. A ticket from pair A presented to verifier B fails registry lookup — there is no "wrong authority" state to configure, because the verifier has no pluggable keyring.
2. **Mis-wiring fails at startup, not at first reset.** Both handles expose `PairID()`. The composition root asserts equality once:
   `if verifier.PairID() != updater.PairID() { refuse to start }`
   This converts the most likely composition error — C7 wired to one pair, supervisor to another — from a silent latent failure into a startup failure. Without this assertion, the failure appears hours later as "reset rejected", which is indistinguishable from a correct rejection.
3. **Absence is fail-closed and explicit.** Nil verifier means "no update authority installed"; every reset is rejected with a typed error. Invalid config (non-positive TTL, empty operation scope) makes `cap.New` return an error, so the composition root cannot start with a half-built authority.

### 2.2 Failure modes and required results

| # | Scenario | Required result | Class |
|---|---|---|---|
| A1 | Two independently constructed pairs | each verifier accepts only its own tickets | P |
| A2 | Ticket from pair A presented to verifier B | typed `ticket_unknown`; no reset | P |
| A3 | Accidental installation of the wrong verifier | startup assertion fires; no supervisor constructed | P |
| A4 | Missing verifier | every reset rejected; no panic; status unaffected | P |
| A5 | Invalid `cap.New` config | constructor error; composition root refuses to start | P |
| A6 | Multiple supervisors sharing one verifier | independent consumption state; one ticket cannot authorise two resets | P |
| A7 | Authority torn down while a supervisor is live | `Consume` returns `authority_closed`; no reset; no panic | P |
| A8 | Concurrent `Complete` + `Consume` | serialised by the authority mutex; no data race; exactly one consumption per ticket | P |

A6 and A7 are the two the review specifically named, and both are cheap to state but easy to leave unspecified — hence listed here rather than assumed.

## 3. Preserved limitation — digest delta proves change, not intent

Carried forward verbatim as a standing limitation, not a footnote:

> Requiring a binary or manifest digest to differ proves that **something** changed. It does not prove that the **intended** update occurred.

Consequences that must survive into every later revision:

- the compensating control in §3.4 of A1 is a *narrowing* of the trust in C7, not an elimination of it;
- C7 could perform an unintended-but-real file change and still obtain a ticket;
- therefore the trust assignment in A1 §3.5 remains an **O** item, and no document may describe issuance as "verified update completion" without the qualifier "verified as an observable state change".

## 4. Capability lifecycle — trusted-operation membership and lock order

### 4.1 The objection and its answer

The review's point: a record reference is not proof that the operation producing it was authorised; creation, registration, and validation must all be inside the trusted operation. Agreed, and it changes the design in one specific way:

**Membership in the in-flight registry is the authorisation.** A record is created, validated, and registered by the same operation invocation, and the apply step accepts only records that are currently *registered and live*. A record that exists but was never registered is not a capability; a record that has been retired is not a capability. Nothing about the reference itself confers anything.

### 4.2 Lock order — defined to make deadlock impossible by construction

Three locks are involved: `Client.mu` (in-flight registry), the record's own `mu`, and `Supervisor.mu`. The order is **sequential, never nested**:

```text
1. Client.mu        → read liveness + generation snapshot, release
2. record.mu        → claim single-use, release
3. Supervisor.mu    → apply identity change, release
```

Rules:

- `Client.mu` MUST NOT be held while acquiring `record.mu`, and vice versa. Nesting is what creates the cycle, so the design forbids nesting rather than ordering it.
- `record.mu` MUST NOT be held across any transport call or any supervisor call.
- `Supervisor.mu` is never acquired while `Client.mu` or `record.mu` is held.
- Retirement (`Client.mu`) and consumption (`record.mu`) are therefore independent; the only shared fact is the record's `consumed` flag, which is the single source of truth for "already applied".

Consequence, stated: a record can be retired between step 1 and step 2. That is not a lost update — step 2 fails closed on a retired record, and the reset/selection simply does not occur. Fail-safe direction is preserved.

### 4.3 Required tests

| # | Scenario | Required result |
|---|---|---|
| C1 | apply before validation | rejected; nothing published |
| C2 | apply after retirement | rejected |
| C3 | apply twice | second is a no-op; exactly one effect |
| C4 | concurrent apply | one wins under `record.mu`; no duplicate effect |
| C5 | apply after generation replacement | rejected on generation mismatch |
| C6 | record from another supervisor | rejected (registry mismatch) |
| C7 | record from another provider/engine | rejected on identity check |
| C8 | failure during consumption | no partial state; capability remains or is burned per rule, never half-applied |
| C9 | cancellation racing completion | exactly one terminal; the loser cannot publish |

## 5. Availability consumer audit (class F)

### 5.1 Every consumer of provider-status `Available`

Searched the whole workspace for `status.Available`, `Status().Available`, `.ProviderStatus(`, and `TTSProviderStatus{`:

| Site | Uses | What `false` means today | Impact of change |
|---|---|---|---|
| `supervisor.go:312–325` | **the only projection** (`st.Available = providerAvailable(st.State)`) | state is outside the five available states | becomes class-driven |
| `supervisor.go:308` | zero-value status (`Absent`) | no supervisor | none |
| helper package tests | assertions | as above | tests must be updated |
| **production code outside the helper package** | **none found** | — | **no current consumer** |

### 5.2 The conflation risk, confirmed concretely

`app-tts-catalog.go` reads `TTSVoiceCatalog.Available` — a different type entirely, meaning "these voice model files are installed on disk". Its own detail strings say so: *"Piper voices are installed. Piper binary availability is checked when synthesis starts."* (lines 130–160). Six sites, all voice/model presence. `client.go:937/945` maps the **wire catalog** result. Neither is provider status.

Conclusion: one boolean now carries "is this provider usable" while a similarly named boolean carries "are these files installed". Any r9 text that unifies them is wrong, and any reviewer reading `Available` in the catalog path must not read provider status into it.

### 5.3 What this means for risk

The change is **contractually correct and currently consumer-inert**: no production code reads provider-status `Available` today. The real risk is **latent** — C7 or UI wiring will read it later and inherit whichever semantics we ship. That argues for adding the class-driven test matrix *now*, while there is nothing to break, rather than after a consumer exists.

### 5.4 A2-D6 field disposition (class F gap, separate change)

| Field | Contract requires (§17.12 shape, design lines 2187–2205) | Present in `types/provider.go:56–71` | Disposition |
|---|---|---|---|
| `Loaded` | yes | **no** | separate criterion, separate verification, **not** part of the availability fix |
| `ServiceCapabilities` | yes | **no** | same |
| `ServiceCapabilitiesKnown` | yes | **no** | same |

These cannot disappear into an availability change. They are an application-status-shape gap with its own AC, its own tests, and their own review. Recording it here prevents it from being silently absorbed and later claimed as "done as part of R-40".

## 6. R-44 matrix verification (class F)

### 6.1 What the matrix actually is

| Measure | Result |
|---|---|
| State classes | 6: `A` (ordinary absent), `K0`, `K1`, `L` (live asset-specific unavailable), `X` (exited/no-instance), `U` (other) |
| Operation rows per class | 8: `Health`, `Catalog` direct, `Catalog` `ProviderCatalog`, `Unload`, `Cancel`, `Synthesize`, `ProviderStatus`, `Voice-list` |
| Cells | 48 — confirmed (the 42 backticked rows plus 6 unbackticked `Voice-list` rows) |
| Columns per cell | helper/transport · frame · start-handshake-load · result/failure · caller `AudioError` · aggregate/`LastError` · owner |

The count is correct. **The count is not coverage**, which is exactly the review's point, and the gaps are specific:

### 6.2 Coverage defects found

| # | Defect | Evidence | Severity |
|---|---|---|---|
| 1 | **`Load` has no row.** `Client.Load` is exported, state-changing, and absent from all six tables | 0 matches for a `Load` matrix row; `client.go:551` defines it | high |
| 2 | **Supervisor cancellation has no row.** `cancelForSupervisor` exists only as prose in §5 | 1 prose mention; 1 I/O-inventory mention; no cell | high |
| 3 | **No cell states the pending-request budget effect.** `MaxPending` appears 3× in the whole design, never in a cell | grep counts | medium — and it is one of the five properties the owner requires |
| 4 | **Compound operations are unrepresented.** Identity replacement (A2-D7), update-gated reset, and ensure-load-inside-`Synthesize` appear in prose only | `identity replacement` 0, `compound` 0, `ensure-load` 10 (prose only) | medium |
| 5 | **"Owner" is a decision owner, not a single publisher.** Cells read "C7 adapter; C1 read path" — two parties, no publication split | cells in `K1`, `L`, `X` | medium — conflicts with the single-writer requirement of R-43 |
| 6 | **Legality is implicit.** No cell says "legal/illegal" as a field; it is inferred from the result column | all cells | low |

### 6.3 Required matrix shape (proposal, class P)

Minimum coverage is **not** 48. Adding `Load` and supervisor cancellation gives **10 operations × 6 classes = 60 cells**; representing identity replacement and update-gated reset as first-class operations adds 12 more if they are not folded into an existing row. Each cell must carry the five requested properties explicitly:

1. legality (yes/no + typed error);
2. exact resulting state (all of aggregate, process, loaded, availability class);
3. **single** publishing component;
4. pending-budget effect (none / consumes / rejects at 32 / exempt);
5. externally visible projection.

Recommendation: keep the 6 state classes, expand to 10 operation rows, and treat compound operations as explicit rows rather than prose — because prose is where the current defects hide.

## 7. Decisions requested

| # | Decision | Recommended |
|---|---|---|
| 1 | Accept startup pair-ID assertion as the mis-wiring control (§2.1) | accept |
| 2 | Accept sequential-not-nested lock order (§4.2) | accept |
| 3 | Accept 60-cell minimum matrix (§6.3) | accept |
| 4 | Schedule A2-D6 field gap as its own criterion (§5.4) | accept, separate |
| 5 | Accept bounded same-package selection residual (A1 §4.6, unchanged) | accept |
| 6 | Availability class internal-only (A1 §5.8, unchanged) | accept |
| 7 | Supply R-38/R-39 text | required |

## 8. Still open

R-38/R-39 (class R, no text). R-44 remediation approval (the 60-cell shape is P). The A2-D6 three-field gap. The authority-boundary architecture remains **proposed, not approved**. No test executed; no code changed.

## 9. Authorization

**C1 remains unauthorized.** This addendum is analysis only: no design document edited, no source or test touched, no dependency or configuration changed, nothing built or run, pre-existing dirty tree untouched.
