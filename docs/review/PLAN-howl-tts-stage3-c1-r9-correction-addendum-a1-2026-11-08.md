# C1 Revision 9 Correction Plan — Addendum A1

**Date:** 2026-11-08
**Status:** OWNER REVIEW RESPONSE. The r9 plan is **not an approved design**; this addendum does not approve it either. **C1 REMAINS UNAUTHORIZED.**
**Addresses:** the three scrutiny items raised against `PLAN-howl-tts-stage3-c1-r9-correction-2026-11-08.md`, and the intake-classification boundary.
**Does not:** edit the r9 plan's architecture into "approved", edit any design document, or resolve R-38/R-39.

## 1. What this addendum does

The owner's verdict was "promising plan, not approved design", with three items requiring scrutiny and one
intake boundary. This addendum:

1. states the **complete dependency and object-wiring model** for the proposed capability package, including
   an impersonation analysis and the exact trust assignment (§3);
2. states the **precise ownership lifecycle** of the in-flight request capability, including transferability
   and invalidation (§4);
3. traces the **blast radius of stored availability state** across writers, transitions, projections and
   consumers, and locates the field (§5);
4. adds a **claim taxonomy** so findings, reviewer reports, author proposals and owner decisions are never
   merged (§2);
5. lists the **decisions this addendum deliberately leaves to the owner** (§6).

It contains no new design claims beyond those needed to answer the three questions, and every such claim is
labelled as a *proposal* under the taxonomy in §2 — not as a finding.

## 2. Claim taxonomy (intake boundary)

Four categories, never collapsed. Every substantive statement in the r9 corpus must carry one.

| Class | Definition | Can a revision assert it? |
|---|---|---|
| **F — source-verified finding** | independently confirmed against contract text and/or observed source in this workspace | yes, with citation |
| **R — reported, unverified** | stated by a reviewer, not yet independently checked here | only as "reported", never as established |
| **P — author proposal** | a correction authored here; unproven and unapproved | yes, labelled |
| **O — owner decision** | belongs to the owner; cannot be settled by a document | no; must remain open |

Application to the current corpus:

| Item | Class | Basis |
|---|---|---|
| R-24 … R-31 (r7) | F, `review-closed` rung 3, **not implemented** | r7 review + source still defective |
| R-32 … R-37 (r8) | P at rung 1 | authored in r8, unreviewed |
| R-40 (availability has no explicit state) | **F** | `providerAvailable` is a pure switch over `State`; A2-D5 gives two answers for `unavailable` |
| R-41 (intent inferred, not captured) | **F** | `observeExit` reads `providerState == stopping` |
| R-42 (transport-nil, write-failure) | **F** | silent return with no state write; client publishes `crashed` |
| R-43 (contradictory writers) | **F** | `SetState` and `MarkHandshake` write `processState` |
| R-44 (matrix is projection-only) | R | not yet re-derived here |
| R-45 (no adapter interface) | **F, partly** | `ProviderCatalog` is `return c.Catalog(...)`; an audio-side interface does exist at `types/provider.go:113–114` |
| R-46 (closure claims) | **F** | r8 register says "closed" for unimplemented findings |
| R-38, R-39 | R (identifier only) | no requirement text available |
| §3–§5 of this addendum | **P** | proposals awaiting review |
| C1-AC-15, `MaxPending=32`, HW-1.6, §6 items | **O** | owner authority |

## 3. Trust boundary and complete wiring model

### 3.1 Why the package is not itself the boundary

A package boundary constrains *who can name a symbol*, not *who is authorised to use it*. Therefore the
boundary must not rest on "only privileged code would call `New`". The design property that actually has to
hold is:

> **Authorization is established by installation, not by construction.**

Constructing a capability object confers nothing. A capability is authoritative only if the component that
enforces the rule was wired with *that* issuer at composition time. Anything else is a parallel object with
no authority over the installed component.

### 3.2 Object graph and wiring

```text
composition root (app wiring)                       built once, at startup
  │
  ├─ verifier, updater := cap.New(cfg)              one split pair, two distinct types
  │
  ├─ helper.NewSupervisor(SupervisorOptions{
  │        UpdateVerifier: verifier,               supervisor stores ONLY the verifier
  │        ... })                                  no getter, no runtime setter
  │
  └─ c7.NewUpdaterClient(updater)                  C7 stores the updater

  helper ──imports──► cap                          (one-way)
  c7     ──imports──► cap                          (one-way)
  cap    ──imports──► stdlib only                  (leaf; cannot import helper or c7)
```

Type split, which is the enforceable part:

| Type | Methods | Held by | Cannot |
|---|---|---|---|
| `Verifier` | `Consume(ctx, Ticket, Binding) error` | supervisor only | issue, mint, inspect registry |
| `Updater` | `Begin`, `Complete`, `Issue` | C7 only | consume, read another update's record |
| `Ticket` | none (no constructor, no accessor, no `String`) | transferred by value | be constructed with a chosen nonce |

There is **no** `Authority` object exposing both halves, and **no** method on `Verifier` that returns an
`Updater`. That absence is deliberate and testable.

### 3.3 Impersonation analysis — the five owner questions

| Question | Answer | Enforced by |
|---|---|---|
| Who constructs the authority? | the composition root, once | single call site; not a C1 surface |
| Who receives `Updater`? | C7 only | passed at construction; not reachable from a `Verifier` |
| Who receives `Verifier`? | the supervisor only | `SupervisorOptions` field, no setter |
| Can another application package obtain an updater through the public API? | It can call `cap.New` — and it gains **nothing**: that creates a *separate* authority whose tickets the installed verifier has never issued | installation-bound trust: the supervisor validates only against the registry of the verifier it holds |
| What prevents the update owner issuing without independently verified completion? | `Issue` requires a completion record and re-verifies observable state; it refuses on mismatch **and refuses when nothing changed** | see §3.4 |

The fourth row is the answer to the reviewer's central objection. A rogue package can construct, but cannot
become the installed authority without the composition root's cooperation. That is a real boundary — and it
is a boundary on *installation*, which is what Go can express, rather than on *caller identity*, which it
cannot.

### 3.4 Evidence requirement for issuance

`Issue` succeeds only when all of the following hold, checked inside `cap`:

1. a completion record exists for this update, created by `Begin` → `Complete` in order;
2. the record's captured **pre-update** digests exist (helper binary, asset-root manifest, config) — so an
   update that changed nothing is distinguishable from one that did;
3. the **post-update** digests re-read by `Issue` match the record and **differ from pre** for at least the
   helper binary or the asset manifest; if nothing changed, `Issue` returns a typed error and no ticket;
4. the record has not already been issued from (one issuance per update);
5. the supplied `Binding` matches the supervisor-observable identity for this update.

Item 3 is the substantive strengthening over the r8 design: r8 let the authority "trust its own recorded
update", which a malfunctioning C7 can satisfy by recording an update it never performed. Requiring a
pre/post delta makes a no-op issuance impossible without an actual file-state change.

### 3.5 Trust assignment — explicit, with residuals

| Trusted party | What is trusted | Compensating control | Residual risk (accepted or open) |
|---|---|---|---|
| C7 update owner | reports that it performed an update | pre/post digest delta required at `Issue` | C7 could perform a *different* update than intended and still obtain a ticket — open, needs owner acceptance (O) |
| composition root | installs the correct pair, once | single call site; code review | mis-wiring is a review defect, not preventable by types |
| Go runtime | no `unsafe`/`reflect` abuse inside the module | out of scope by declaration | memory corruption / `unsafe` — explicitly not covered |
| authority implementation | correct registry, expiry, one-use | C1-owned test matrix; C7 must supply executed evidence | unimplemented; no evidence exists |

### 3.6 What §3 does not claim

It does not claim the capability package is a security boundary against a hostile package inside the same
binary. It does not claim C7 is trustworthy. It does not claim any of this is implemented. It claims exactly
one thing: **a component other than the installed issuer cannot produce an authorization that the installed
verifier accepts**, provided `unsafe`/`reflect` are out of scope. `[P]`

## 4. In-flight request capability — ownership lifecycle

### 4.1 Why "created inside `begin()`" was not an answer

The owner's objection is correct: if unauthorised code can call `begin()`, then minting inside `begin()` moves
the problem without solving it. Two corrections follow, and neither relies on the caller being well-behaved.

**Correction 1 — the capability is not a token, it is a reference to the authoritative record.** The
selection apply does not accept `(requestID, capability)` as two arguments that must agree. It accepts the
record reference alone and derives the request ID, the generation, the validated result and the provider
identity **from that record**. There is nothing to forge because there is no caller-supplied claim to compare.

**Correction 2 — `begin()` is not a public minting entry point.** It is the request-acceptance step of a
validated operation, reachable only from the operation's own call stack, and the record it creates is
registered in the in-flight set that the apply step consults. The apply step refuses any reference whose
record is not in that set, is not live, or was not created by an accepted-and-validated request.

### 4.2 Binding

The capability is bound to, and validated against, all of:

| Bound to | Source of truth | Rejection when mismatched |
|---|---|---|
| the actual request | the record's own ID, not an argument | record not live |
| the validated result | the decoded, schema-validated catalog result stored on the record | validation never ran |
| the operation | `catalog`, not `load`/`synthesize`/`health` | operation field mismatch |
| the generation | supervisor generation captured at acceptance | generation advanced |
| the provider/engine identity | the current validated handshake | identity not in handshake |

### 4.3 Transferability

Not transferable, by construction: the record reference carries its own operation and identity, and the apply
step reads them from the record rather than trusting a second argument. A reference to request A used to
apply a result claimed for request B fails the identity check inside the record — there is no path where the
two can be made to disagree, because there is only one input.

### 4.4 Invalidation matrix

| Event | Effect on the capability |
|---|---|
| completion (apply succeeds) | consumed; single-use flag set under the record's own mutex; second apply is a no-op |
| request failure | retired before publication; apply refuses |
| cancellation | retired; apply refuses |
| timeout | retired; apply refuses |
| generation change (replacement, restart, update reset) | generation mismatch; apply refuses |
| process exit | record retired; apply refuses |
| concurrent apply | record mutex serialises; exactly one consumes |

### 4.5 Enforcement points, each testable

1. unexported record type with no exported fields;
2. apply step takes the record, not an ID;
3. liveness check against the in-flight set;
4. single-use flag under the record's mutex;
5. generation comparison;
6. identity check against the current handshake;
7. selection has no lifecycle, capability, failure, or recovery effect — so even a successful forgery cannot
   clear a terminal state, start a process, or change readiness.

### 4.6 Residual, stated precisely

Same-package code that calls `begin()` itself and then applies its own fabricated catalog result can produce
a *selection* — the same residual r8 had, but with a strictly smaller effect: identity-only, generation-bound,
consumed once, no lifecycle consequence. If review judges even that unacceptable, the fallback is the package
split of §3 applied to catalog operations, at the cost of an extra package. That trade is an **owner decision
(O)**, not a design default.

## 5. Availability — blast radius of stored state

### 5.1 Current shape (facts)

| Fact | Location | Class |
|---|---|---|
| `Available` is computed by a pure switch over aggregate `State` | `supervisor.go:325` → `providerAvailable` | F |
| `TTSProviderStatus` already has an `Available bool` field | `types/provider.go:56–71` | F |
| that Go type has **no** `Loaded`, `ServiceCapabilities`, `ServiceCapabilitiesKnown` | same | F — A2-D6 requires all three |
| the only status-side projection of `Available` is line 325 | `supervisor.go` | F |
| `app-tts-catalog.go` uses a **different** `Available` (catalog/voice-level, `TTSProviderCatalog`) at 6 sites | `app-tts-catalog.go:139–228` | F |

### 5.2 Writer trace (what must change)

| Writer | Today | After | Note |
|---|---|---|---|
| `providerAvailable(state)` | derives | **deleted** | cannot express A2-D5 |
| handshake success | implicit via state switch | sets availability = true, class = none | |
| `setUnavailable(f)` | state only | sets availability by class + stores class | must stop writing `processState` |
| `observeExit` | state only | **must not change availability** (A2-D5 (a)/(b) keep it true) | |
| `SetState` / `MarkHandshake` | write `processState` | removed as `processState` writers | §5.3 of the plan |

### 5.3 Transition trace

Every `unavailable` entry must record a class, and the class decides `Available`. There are exactly three
entry causes with contract text: asset-specific with valid config (A2-D5 (a), true), budget exhausted with
configured provider (A2-D5 (b), true, `ProcessAlive=false`), and helper/config invalid (A2-D5 negative list,
false). A transition that reaches `unavailable` without a class is a design error, not a default-false.

### 5.4 Projection trace

`Status()` is the single projection site. `Available` comes from stored state; nothing else is derived from
it. `Ready` continues to derive from `Loaded` + state (A2-D5), and `ProcessAlive` from `processState` alone
(A2-D4). The three must not be re-coupled in either direction.

### 5.5 Consumer trace and the conflation risk

| Consumer | Uses | Impact |
|---|---|---|
| `Status()` projection | status `Available` | behavioural change: true in two `unavailable` classes where it is false today |
| adapter/tests reading `ProviderStatus` | status `Available` | same |
| `app-tts-catalog.go` (6 sites) | **catalog** `Available`, a different field | **none** — must not be "unified"; doing so would be a regression |
| `client.go:937/945` | wire catalog result mapping | none |

The blast radius is therefore **narrow in code and non-trivial in meaning**: one projection line changes
behaviour for consumers of provider status, while six catalog call sites look similar and are unrelated. Any
r9 text that implies one boolean covers both is wrong.

### 5.6 Where the field belongs

Internal supervisor state (`availability`, `availabilityClass`), projected into the **existing** public
`Available` bool. No new public field is required for availability. The A2-D6 fields
(`Loaded`, `ServiceCapabilities`, `ServiceCapabilitiesKnown`) are a separate, already-designed requirement and
must not be smuggled in as part of this fix — they are a distinct change with their own evidence.

### 5.7 Regression risks

1. consumers that treat `Available=false` as "not installed" would now see true for asset failures — intended
   by A2-D5, but observable;
2. a stored field can go stale where a derived one could not — hence the class is mandatory on entry and the
   field is written only by the writers in §5.2;
3. `app-tts-catalog` semantics must remain untouched.

### 5.8 Owner decision (O)

Whether the *class* is exposed (internal only vs. surfaced for diagnostics), and acceptance of the behaviour
change in §5.5, remain owner decisions. Recommended default: internal only, no public shape change.

## 6. Decisions requested from the owner

| # | Decision | Recommended default |
|---|---|---|
| 1 | Accept installation-bound trust as the authority boundary (§3.3) | accept; document as C1's security claim, with `unsafe`/`reflect` out of scope |
| 2 | Accept C7 as trusted update reporter with pre/post digest compensation (§3.5) | accept with residual recorded |
| 3 | Same-package selection residual: accept bounded residual or pay for package split (§4.6) | accept bounded residual |
| 4 | Availability class internal-only vs. exposed (§5.8) | internal only |
| 5 | Whether `transportNil` needs an A2-D11 amendment | implementation detail; no amendment |
| 6 | Whether matrix behaviour columns become contract text | remain C1 design table |
| 7 | Supply R-38/R-39 requirement text | required before r9 claims coverage |

## 7. What remains unverified

R-38 and R-39 (no text). R-40 … R-46 beyond the parts re-derived as class F in §2. The entire §3–§5 proposal
set, which is unproven, unimplemented, and unreviewed. No test has been executed; no code changed.

## 8. Authorization

**C1 remains unauthorized.** This addendum adds analysis only. No design document was edited, no source or
test touched, no dependency or configuration changed, nothing built or run. The pre-existing dirty tree is
untouched. Approval of the plan remains the owner's, and is not asserted here.
