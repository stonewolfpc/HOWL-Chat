# C1 Revision 9 — Correction Plan (narrow scope)

**Date:** 2026-11-08
**Status:** PLAN ONLY. **No design edit has been made. C1 REMAINS UNAUTHORIZED.**
**Prepared for:** owner review before any r9 design edit.
**Supersedes:** nothing. This is an intake and planning artifact.

## 1. Purpose, authority, and intake status

### 1.1 What this document is

A correction plan for revision 9, organized around findings **C1-R-38 … C1-R-46**. It states the
governing text and its location, the proposed package architecture with dependency direction, the state
ownership and transition model, the exit-intent and terminal-outcome corrections, the matrix
corrections, the disposition of inherited findings, and the unresolved owner decisions — so that the
plan can be reviewed before any design text is rewritten.

### 1.2 Intake disclosure — read this first

The r8 review artifact that would normally define C1-R-38 … C1-R-46 **is not present in this
workspace.** Verified at plan time:

| Check | Result |
|---|---|
| `docs/review/REVIEW-*` files present | r6 and r7 only; no r8 |
| Files anywhere containing `C1-R-38`, `C1-R-46`, `R-38`, or `R-46` | 0 |
| Newest file in `docs/review` | the r8 design/readiness/integrity set written by this task |
| `docs/review/*.md` count | 87 |

Consequence, stated plainly: this plan is built from **the owner's description of seven of the nine
findings** (the R-40 … R-46 table) plus the accompanying architectural analysis, not from the review
text itself. Two findings — **R-38** and **R-39** — are named only as endpoints of the range. Their
requirement text is `[Unestablished]` here.

Accordingly:

- §4 plans the authority-boundary correction at the architectural level the analysis describes, and
  reserves §4.6 for the R-38/R-39 slots.
- Every planned correction in §3, §5, §6, §7, §8 is derived from governing contract text and observed
  source, so it stands independently of the missing review text.
- **No finding is marked addressed by this document.** r9 may not be drafted as a "response to
  R-38 … R-46" until the review text is supplied; otherwise the revision would be answering a
  paraphrase, which is the failure mode this gate exists to prevent.

### 1.3 Evidence classes

`[Observed]` `[Source-derived]` `[Calculated]` `[Inferred]` `[Unestablished]` `[Superseded]` — as
established in r7/r8.

## 2. Findings intake

| Finding | Area (from owner material) | Planned correction | Text available |
|---|---|---|---|
| **R-38** | not described | §4.6 slot; authority-boundary class | `[Unestablished]` |
| **R-39** | not described | §4.6 slot; authority-boundary class | `[Unestablished]` |
| **R-40** | Availability: explicit state, writer, transition ownership | §5.2 explicit `availability` state + writer table | summary only |
| **R-41** | Exit classification: consistent captured-intent model | §6.1–§6.3 | summary only |
| **R-42** | Terminal events: missing transport-nil outcome and write-failure intent | §6.4 | summary only |
| **R-43** | Termination barrier: contradictory writer rules | §5.3 + §6.5 | summary only |
| **R-44** | Operation matrix: exact per-state behavior | §7 | summary only |
| **R-45** | C7 adapter: define actual interface, mark implementation pending | §8 | summary only |
| **R-46** | Closure claims: correct statuses and evidence language | §9 + §10 | summary only |

Nine findings, nine planned corrections, zero claimed closed. `[Calculated]`

## 3. Governing text and exact source locations

### 3.1 Contract decisions (authoritative for the disputed behavior)

| Source | Location | Governs | Key exact wording |
|---|---|---|---|
| A2-D4 | contract line 2505 | exit intent, `ProcessAlive`, incompatible | "`ProcessAlive` is derived solely from `TTSProcessState`: true for `starting`, `alive`, and `stopping`; false for `absent` and `exited`. It is never derived from aggregate `State`." · "The exit observer is intent-aware: an exit caused by supervisor termination after an incompatible handshake is expected and MUST NOT be reclassified as `crashed`; only an unexpected exit is `crashed`." |
| A2-D5 | line 2506 | `Available`/`Loaded`/`Ready` | "`Available` is the configured-provider flag … It MUST NOT be read as proof that the helper process is alive, that assets are present, or that any catalog voice exists." · "It is also true for `unavailable` when either (a) the failure is asset-specific and the helper, protocol, and provider configuration remain valid, or (b) … restart budget was exhausted …; in case (b), `ProcessAlive=false`" · "Thus `Available` and `State` are independent" |
| A2-D6 | line 2507 | service capability snapshot | "`ServiceCapabilitiesKnown` is true **iff** the current helper instance has completed a compatible handshake and remains authoritative: it has not entered `stopping`, `exited`, `crashed`, `incompatible`, or `backoff`, and no replacement has begun." |
| A2-D7 | line 2508 | replacement, budget exemption, exit classification | "(7) only after exit is observed, start a fresh process" · "classify by exit code: 0 → `absent`, nonzero → `crashed`" · "A transport write failure means failed shutdown, `KillTree`, observe `Wait()`, and `crashed`." · "expressly exempt from the §17.11.5 crash/restart budget" |
| A2-D12 | line 2513 | `unavailable` recovery | "no automatic retry occurs, and recovery is triggered only by a new explicit `load` request" · "`synthesize`, `catalog`, voice-list, `health`, and status requests do not trigger that recovery" |

### 3.2 Contract sections

| Source | Location | Note |
|---|---|---|
| §17.7.1 state vocabularies | lines 1743–1760 | 5 process states, 13 aggregate states; the `ProcessAlive` sentence here is **superseded by A2-D4** (line 2529) |
| §17.7.2 supervisor transitions | lines 1763–1794 | 22 transition edges; the `unavailable` row already carries "update-gated reset has precedence" |
| §17.7.4 health table | lines 1816–1837 | 7 base rows + 2 additive rows at 2556–2560 |
| §17.9 typed error registry | line 1876+ | `request_not_found` origin `wire`; `protocol_violation` origin `app` |
| §17.12.1 adapter behavior | lines 2230–2249 | `ProviderCatalog`/`CatalogCurrent` split is a C1 addition, not contract text — hence R-45 |
| §17.19.3 supersession map | lines 2515–2550 | establishes what A2 supersedes, inc. §17.7.1 `Available` |
| A2 consolidated bullets | lines 2554–2574 | 12 numbered behaviors; item 5 and 12 are the availability anchors |

### 3.3 Observed source (evidence that the defects are real, not hypothetical)

All `[Observed]` on `internal/audio/tts/helper/supervisor.go` and `client.go` as they exist now. The
files are untracked, so these are current working-tree reads, not committed history.

| Observation | Where | Relevance |
|---|---|---|
| `processAlive(process, provider)` returns false when provider state is `crashed` or `unavailable` | supervisor.go, `processAlive` | R-24 design-corrected; **implementation still defective** → R-46 |
| `providerAvailable(state)` is a pure switch over aggregate state | supervisor.go, `providerAvailable` | **cannot** express A2-D5 (a)+(b): `unavailable` maps to one answer for two causes → R-40 |
| `SetState` writes `providerState`, `activeRequestID`, `lastFailure`, **and** `processState=alive` | supervisor.go | generic writer crossing three ownership domains → R-43 |
| `setUnavailable` writes `processState=absent` even when a process is alive | supervisor.go | liveness/state contradiction → R-40, R-43 |
| `MarkHandshake` writes `processState=alive` plus identity and capabilities | supervisor.go | second `processState` writer outside the lifecycle owner → R-43 |
| `observeExit` classifies by `providerState == stopping` only | supervisor.go | intent inferred from aggregate state, not captured → R-41 |
| `Catalog` calls `c.sup.SetSelected(...)` after validation | client.go (~line 546) | r7 removed the setter in design only; source unchanged → R-46 |
| `ProviderCatalog` is exactly `return c.Catalog(ctx, req)` | client.go (~961–963) | adapter and direct path are the same function → R-35/R-45 |
| `transport == nil` returns `helper_start_failed` with **no state write** | client.go, synthesize | state remains `busy` with no terminal → **R-42 transport-nil** |
| `transport.Send` failure writes `SetState(crashed, …)` from the client | client.go, synthesize | lifecycle write by a non-owner, no intent record → R-42, R-43 |
| ~10 distinct `SetState` call sites inside `synthesize` | client.go | publication is scattered, not a closed terminal set → R-42 |

## 4. Proposed package architecture and dependency direction

### 4.1 The actual defect class

Go enforces exactly two things: field- and identifier-level privacy **per package**, and interface
satisfaction **per method set**. There is no friend declaration, no caller identity, no per-call-site
authorization. Therefore:

- an unexported field is package-private, not friend-private;
- an unexported function is callable by every file in its package;
- an exported issuer is an authorization surface regardless of how it is documented;
- a capability is only as trustworthy as the operation that created it, so a value that can be created
  outside that operation is not a capability at all.

r7 and r8 both treated "only this caller would do that" as enforcement. The correction is to move each
capability behind a **package boundary whose only minting path is the operation that performs the work**,
and to make the two handles on the authority — issue and verify — *different types*.

### 4.2 Proposed package

New package `internal/audio/tts/cap`, owning the update authority and every minted capability.

```text
internal/audio/tts/cap            (imports: stdlib only)
    Authority   → owns issuance registry, verification evidence, consumption state
    Updater     → handle with Issue/CompleteUpdate; held only by C7
    Verifier    → handle with Consume only; held by the supervisor
    Ticket      → exported type, unexported 32-byte nonce, no constructor, no accessor
    Binding     → provider/engine/generation/epoch/operation, built by the supervisor
    Sealed      → unexported marker method on each interface, so no external implementation exists
```

### 4.3 Dependency direction (must be enforced, not assumed)

```text
        app wiring
        │        │
        │        └──────────────► internal/audio/tts/cap   (no imports of helper)
        ▼                              ▲
internal/audio/tts/helper ─────────────┘   (helper imports cap; one-way)
        │
        ▼
internal/proc, internal/audio/types
```

Rules, each checkable by `go list -deps` / an import-scan test:

1. `cap` MUST NOT import `helper`. No back-edge, so nothing in `cap` can reach supervisor internals to
   reinterpret its own state. `[Source-derived]`
2. `helper` MUST import `cap`, never the reverse, and must not re-export `Ticket` construction.
3. C7 imports `cap` for `Updater`; the supervisor receives only `Verifier`. The two handles expose
   disjoint method sets, so holding the authority object does not confer issuance on the consumer side.
4. No import cycle; `cap` is a leaf.

### 4.4 Issuer and verifier ownership

| Decision | Owner | Enforced by |
|---|---|---|
| Minting a ticket | `cap.Authority`, via `Updater` | unexported minting path; `Ticket` has no exported field or constructor |
| Recording that an update completed | `cap.Authority` | `Updater.CompleteUpdate` captures observable evidence |
| Verifying | `cap.Verifier.Consume` | registry lookup, one-use, binding match, expiry, revocation |
| Applying the reset | supervisor | only after `Consume` returns nil and rechecks pass |
| Installing either handle | app wiring, once at construction | no runtime setter in C1 |

**Residual trust, stated rather than hidden:** C7 can lie about having performed an update. That is not
solvable by types, because C7 *is* the update owner. The plan shrinks it: `Issue` re-reads observable
evidence — helper binary digest, asset-root manifest digest, config digest — and refuses to issue when
those differ from what `CompleteUpdate` recorded. A fabricated completion then fails unless the
filesystem genuinely changed. The remaining assumption is "the update owner reports its own work
truthfully", which must be recorded as a trust assignment in the design and accepted by the owner, not
presented as enforcement. `[Inferred]`

### 4.5 Supervisor invocation path

`updateGatedReset(ticket)` → build `Binding` from supervisor fields under `mu` → `Verifier.Consume`
**outside** `mu` → on nil, reacquire, recheck (no live process, generation and epoch unchanged) → apply.
Recheck failure burns the ticket and performs no reset. No authority installed → reject. `[Source-derived]`

### 4.6 Slots for R-38 and R-39

Reserved. Both are expected to be authority-boundary findings, and §4.1–§4.5 is the planned correction
class for that boundary. The specific requirement text must be supplied before r9 asserts that either is
addressed.

## 5. State ownership and transition model

### 5.1 Why the model must be stated

R-40, R-41, and R-43 are one defect wearing three labels: **no single writer per state variable**, so
"who may set this" is answered per call site by convention. The plan replaces convention with a table
that each writer must satisfy, plus a test that greps writes.

### 5.2 State variables, sole writers, allowed transitions

| Variable | Sole writer(s) | Allowed transitions | Forbidden writers |
|---|---|---|---|
| `processState` | `Start` (absent→starting), handshake publication (starting→alive), `beginTermination` (alive→stopping), `observeExit` (→absent/exited) | the 5-value vocabulary, §17.7.1 | `SetState`, `MarkHandshake`, `setUnavailable`, client code — all removed as `processState` writers |
| `providerState` | domain events only: handshake, `observeExit`, load outcome, synthesis outcome, budget exhaustion | 13-value vocabulary, §17.7.2 edges | generic `SetState` removed; no client-side lifecycle writes |
| `availability` | **new explicit field**, not derived | set true on compatible handshake; set by failure class on `unavailable`; unchanged by exit | derived-from-`providerState` functions (`providerAvailable`) deleted |
| `availabilityClass` | `assetSpecificConfigured` \| `budgetExhaustedConfigured` \| `helperOrConfigFailure` | one per `unavailable` entry; cleared on next attempt | client code |
| `loaded` | load/unload/exit/replacement events | A2-D5 sentence list | `SetState` |
| `serviceCapsKnown` | handshake (true); exit/stop/incompatible/replacement-begin (false) | A2-D6 invariant | any other writer |
| `generation` | `Start`/replacement | monotonic ++ | — |
| `epoch` | update-gated reset | monotonic ++ | — |
| `activeRequestID` | the request that owns it: set on accept, cleared on terminal publication | one active request | `SetState` on unrelated paths |
| `lastFailure` | the outcome writer for that domain | replaced per outcome | — |

Rationale for making `availability` explicit rather than derived `[Source-derived]`: A2-D5 assigns
**two different** `Available` answers to `State=unavailable` depending on cause, and `providerAvailable`
is a pure function of `State`. A pure function cannot represent A2-D5, so the value must be stored with
its class. This is the concrete content of R-40.

### 5.3 Termination barrier (R-43)

| Barrier rule | Owner | Contradiction being removed |
|---|---|---|
| `beginTermination` may write only `processState=stopping` | termination path | `SetState` also wrote `processState` (to `alive`) mid-termination |
| After `beginTermination`, no load/synthesis/catalog transition may be published | request paths | previously a late `SetState` could resurrect `alive` |
| Exit classification reads the captured intent, never `providerState==stopping` | exit observer | inferring intent from aggregate state is ambiguous when aggregate is `crashed` |
| Exactly one of {termination, exit observer} publishes the terminal process outcome | both | first-intent-wins replaces last-writer-wins |

## 6. Exit-intent and terminal-outcome corrections (R-41, R-42)

### 6.1 Captured intent, enumerated

`exitIntent` as a closed set, captured at the moment the exit is *caused*, first-intent-wins, readable by
the exit observer:

`intentNone` (unexpected) · `intentShutdownGraceful` · `intentShutdownForced` · `intentIncompatible` ·
`intentReplace` · `intentUpdateReset` · `intentCancelKill`

Mapping to the contract: `intentIncompatible` and `intentShutdownGraceful` are "expected" and MUST NOT
be reclassified `crashed` (A2-D4, A2-D7); `intentNone` is the only `crashed` source; `intentReplace` and
`intentUpdateReset` carry the budget exemption (A2-D7 "expressly exempt").

### 6.2 First-intent-wins and late exit

Rules: the intent is written once under `mu` and never overwritten; a second cause observes the recorded
intent; an exit observed after a *replacement* has begun is attributed to the old generation and MUST
NOT mutate the new generation's state — this is the replacement-overlap rule R-25 established, now
grounded on an explicit generation comparison rather than pointer identity.

### 6.3 Idempotence

`observeExit` applied twice for the same generation is a no-op, not a second classification. `[Inferred]`

### 6.4 Terminal outcomes — closed set with the two missing members

The client currently publishes outcomes ad hoc from ~10 sites. The plan replaces that with a closed
`terminalOutcome` set, published exactly once, with the intent recorded:

| Outcome | Trigger | Contract anchor | Present today? |
|---|---|---|---|
| `success` | validated result accepted | §17.11.4 step 7 | yes |
| `providerFailure` | wire error response | §17.9 | yes |
| `transportSendFailure` | `Send` error | A2-D7 "transport write failure" | partial — writes `crashed` from the client, no intent → **R-42** |
| **`transportNil`** | no current transport | — | **missing**: returns `helper_start_failed` with **no state write**, leaving `busy` → **R-42** |
| `timeout` | operation timeout | registry `timeout` | yes |
| `cancelAcknowledged` | cancel terminal received | §17.11.4 | yes |
| `cancelGraceExpired` | grace elapsed | §17.11.4 | yes |
| `decodeFailure` | result decode error | registry | yes |
| `validationFailure` | result/WAV validation | §17.11.4 | yes |
| `protocolViolation` | unmatched/duplicate ID, publication after terminal | §17.9 `protocol_violation` origin `app` | partial |

Two required corrections, both observable in source: (1) `transportNil` must have a defined outcome and a
defined state transition instead of returning silently; (2) `transportSendFailure` must record intent and
must not be the client's licence to publish a lifecycle state.

## 7. Operation matrix corrections (R-44)

The matrix currently specifies **projection** per cell; R-44 requires **behavior**. Each of the 48 cells
(6 operations × 8 aggregate states) must state, verbatim-specified:

1. whether the operation is legal, and the typed error + registry code if not;
2. the state the operation leaves behind on success and on each failure class;
3. which component publishes each resulting transition (single writer, per §5.2);
4. the effect on restart budget (consumed / exempt / none);
5. the `Available`/`Loaded`/`Ready`/`ServiceCapabilitiesKnown` projection, consistent with A2-D5/D6.

Row/column count stays 48 so the existing count is preserved; the change is cell content, not matrix
size. A cell that cannot be filled from governing text becomes an owner decision (§11), not an
invention.

## 8. C7 adapter interface (R-45)

### 8.1 What is actually there

`ProviderCatalog` is `return c.Catalog(ctx, req)` — one function, so there is no separate adapter path to
constrain, and no interface at all. The r7/r8 "adapter may start only from ordinary `absent`, then call
no-start `CatalogCurrent`" rule therefore constrains a distinction the code does not make. `[Observed]`

### 8.2 Planned definition

A named interface with an explicit method set, an explicit allowed-entry-point list, and an explicit
no-start guarantee, e.g.:

```text
type ProviderCatalogInspector interface {
    ProviderCatalog(ctx, req) (types.TTSProviderCatalog, error)  // may start from absent
    CatalogCurrent(ctx, req) (types.TTSProviderCatalog, error)   // MUST NOT start
}
```

plus: which supervisor entry points each may call; that neither may load, select, or mutate lifecycle;
that `CatalogCurrent` with no live process returns a typed error and starts nothing; and the ownership
carrier for the resulting selection (see §8.3).

Status: **design-defined, implementation pending**, and marked as such. The interface must not be
described as existing until a commit exists.

### 8.3 Selection capability — replace the handoff, don't secure it

r8 minted a run token via `beginCatalogSelection`, which any file in `helper` can call. Under the
principle in §4.1 (a capability is only as good as the operation that created it) the correction is to
**remove the standalone minting function**. The selection capability becomes the live in-flight request
entry itself — an unexported type created only by the operation's own `begin()` — so the only way to hold
it is to be inside that call stack, and the supervisor validates that the entry is still live and belongs
to the request being reported. If review finds any residual same-package reachability unacceptable, the
fallback is to move the catalog operation and the selection apply into `cap`-style package separation;
that costs an extra package and is a deliberate trade, not a default.

## 9. Closure-claim language standard (R-46)

A seven-step status ladder, one status per finding, no unqualified "closed":

| # | Status | Requires | r8 usage |
|---|---|---|---|
| 1 | `design-corrected` | text changed in a revision | R-32 … R-37 |
| 2 | `author-reconciled` | author checked against quoted governing text, citation recorded | R-24 … R-31 (r7) |
| 3 | `review-closed` | an independent review artifact (with hash) closes it | R-24 … R-31 (r7 review) |
| 4 | `owner-accepted` | owner decision record | none |
| 5 | `implemented` | code exists at a named path/revision | **none** |
| 6 | `tested` | named test executed with recorded result | **none** |
| 7 | `evidenced` | archived output/hash | **none** |

Rules: a finding carries exactly one status plus a reference; a design revision may only raise a status
to 2 (or cite an existing 3); statuses 5–7 can never be asserted by a document revision; and any sentence
using "closed", "resolved", "done", or "complete" must name the rung. `[Source-derived]` from the owner's
stated intent.

**r8 self-correction this standard forces:** r8's register says R-24 … R-31 are "closed in r7; retained".
Under this ladder they are `review-closed` (rung 3) **and not implemented** — the observed source still
contains `processAlive(process, provider)`, `SetSelected`, `SetState`, and `MarkHandshake`. The word
"closed" without rung and without the implementation caveat is exactly the inaccuracy R-46 identifies,
and r9 must restate all 37 findings under the ladder.

## 10. Disposition of inherited findings

| Findings | r9 status under §9 | Reference | Implementation |
|---|---|---|---|
| R-1 … R-23 | as recorded in r7 §2.5, restated with an explicit rung | r7 design/review | pending (never claimed) |
| R-24 … R-31 (8 findings) | rung 3 `review-closed` | r7 review `c025ba0a…` | **pending — source unchanged** |
| R-32 | rung 1 `design-corrected` | r8 §4 | pending |
| R-33 | rung 1; correction class replaced by §4 here | r8 §4 + this §4 | pending |
| R-34 | rung 2 `author-reconciled` | A2-D5 (a), line 2506 | pending |
| R-35 | rung 1; carrier design replaced by §8.3 | r8 §5 + this §8.3 | pending |
| R-36 | rung 1 | r8 §7 | n/a (documentation) |
| R-37 | rung 1 | r8 §6.3, C1-AC-33 | pending |
| R-38 … R-46 | **not dispositive until review text is supplied** | — | pending |

No finding in this plan is asserted at rung 3 or above except as an inherited r7 citation. `[Calculated]`

## 11. Unresolved owner decisions

Carried, unchanged by this plan:

1. **C1-AC-15** — resource determinations (13 rows) and I/O N/A entries (8 categories).
2. **`MaxPending=32`** — no accepted authorised source; R8-AC-4 disposition unrecorded.
3. **HW-1.6** independent review.

New, surfaced by this plan and requiring owner authority before r9 can settle them:

4. **Availability as stored state** — A2-D5 requires two answers for `unavailable`, which §17.12's
   projection cannot carry today. Whether `availabilityClass` is internal-only or must appear in the
   application status shape is an owner decision (A2-D6 makes §17.12 authoritative).
5. **Terminal-outcome set extension** — `transportNil` is not named in the contract; adding it may need an
   amendment under A2-D11's method, or may be implementation detail. Owner call.
6. **C7 trust assignment** — accepting that issuance is trusted to the update owner (§4.4 residual), with
   observable-evidence verification as the compensating control.
7. **Matrix behavior columns** — whether per-state operation legality becomes normative text in the
   contract or remains a C1 design table.
8. **R-38/R-39 requirement text** — needed before r9 claims coverage.

## 12. r9 execution sequence and non-goals

| Step | Work | Gate |
|---|---|---|
| S1 | Receive review text; dispose R-38/R-39 | owner provides artifact |
| S2 | §4 package architecture as normative text, with dependency rules and import-scan test | plan review |
| S3 | §5 state ownership table + `availability` explicit state | owner decision 4 |
| S4 | §6 captured intent + closed terminal set (incl. `transportNil`, write-failure intent) | owner decision 5 |
| S5 | §7 matrix cells to behavior level | owner decision 7 |
| S6 | §8 adapter interface + selection capability without a minting function | — |
| S7 | §9 rung labels on all 46 findings; ACs and verification rows reconciled | — |
| S8 | Structural validation + hash + readiness/integrity, then independent review | — |

Non-goals for r9: no implementation, no source or test edits, no budget or resource decisions, no
contract amendment without the A2-D11 method, no change to `C1-AC-15` or `MaxPending=32`.

## 13. Validation plan for r9 (to be executed, not asserted)

Counts to reconcile: 46 findings × one rung; criteria = r8's 33 plus whatever S3/S4/S5 add, each with one
verification method; matrix 48 cells × 5 behavior fields; zero duplicate headings; each ownership
restriction naming a mechanism (package boundary, unexported minting path, single-writer table); zero
unqualified "closed"; import-scan test proving §4.3 direction.

## 14. Authorization

**This plan authorizes nothing.** C1 remains unauthorized; no design text was edited; no source, test,
dependency, configuration, build, or runtime action was taken; the pre-existing dirty tree is untouched.
r9 drafting awaits review of this plan and the missing R-38/R-39 text.
