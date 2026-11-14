# C1 — Clause Disposition, Ownership Evaluation, Inventory Completion

**Date:** 2026-11-08
**Status:** READ-ONLY REVIEW ARTIFACT. Not r9 design text. **C1 REMAINS UNAUTHORIZED.**
**Frozen predecessors (not overwritten):** reconcile `479ac26c…`, audit `1541332e…`, A1-rev1 `c061e9d5…`, A2-rev1 `26342c5e…`, A1 `92acde55…`, A2 `0a5ea047…`, r8 `a8c40fac…`, plan `99897422…`.

Three workstreams, deliberately kept separate: **A** clause disposition and supersession, **B** ownership evaluation, **C** inventory completion. No workstream closes another's item.

Classes: **V** verified existing · **C** required change · **N** specified, not implemented · **X** normative conflict.

## Part A — Clause disposition

Disposition vocabulary: **ACCEPT** (V, becomes review baseline) · **CHANGE** (C, C1 acceptance criterion) · **DEFER** (N, owner must confirm out of scope) · **DECIDE** (owner choice required) · **SUPERSEDE** (X, amendment required).

| Row | Class | Disposition | Authority |
|---|---|---|---|
| R-1 restart schedule 1/2/5 s | N | DEFER | O-13 |
| R-2 three per 10 min, reset after 10 min | N | DEFER | O-13 |
| R-3 `crashed → backoff` | N | DEFER | O-13 |
| R-4 `backoff → starting` | N | DEFER | O-13 |
| R-5 budget exhausted → `unavailable` | N | DEFER | O-13 |
| R-6 `Available` for budget-exhausted | N | DEFER | O-13 |
| R-7 `stopping` exit 0 → `absent` | V | ACCEPT | — |
| R-8 `stopping` grace/nonzero → `crashed` | C | CHANGE | — |
| R-9 exit codes 0/2/3/4/other | C | CHANGE | — |
| R-10 unexpected `Wait()` → `crashed` | V | ACCEPT | — |
| R-11 hello failure → `incompatible`, no auto-restart | V | ACCEPT (automatic only) | O-14 |
| R-11b no explicit retry after `incompatible` | N | CHANGE | O-14 |
| R-12 `starting` failure → `crashed` or `unavailable` | C | CHANGE | — |
| R-13 `ProcessAlive` from process state only | C | CHANGE | — |
| R-14 §17.7.1 old definitions | X | SUPERSEDE | §A.2 |
| R-15 `Available` for unloaded/loading/ready/busy/degraded | V | ACCEPT | — |
| R-16 `Available` for asset-specific `unavailable` | C | CHANGE | — |
| R-17 `Loaded` false on exit/unload/identity change | N | DEFER | O-19 |
| R-18 `Loaded`/`ServiceCapabilities`/`…Known` fields | N | DEFER | O-19 |
| R-19 identity replacement sequence | N | DECIDE | O-20 |
| R-20 idle close | N | DECIDE | O-20 |
| R-21 supervisor cancellation path | N | DECIDE | O-20 |
| R-22 explicit-load recovery | V/C | ACCEPT + CHANGE (gate) | O-20 |
| R-23 ensure-load inside synthesize | N | CHANGE | — |
| R-24 direct `Catalog` must not start | C | CHANGE | O-17 |
| R-25 `ProviderCatalog` read-only, identity-neutral | C | CHANGE | O-17 |
| R-26 cancel grace → kill | V | ACCEPT | — |
| R-27 delete only after observed termination | V | ACCEPT | — |

### A.1 The N rows are a scope list, not a defect list

R-1..R-6, R-17..R-21, R-23 were never implemented. Only O-13, O-19 and O-20 can decide whether they belong to Stage 3 C1. Writing acceptance criteria for them before that decision is how the previous revision ended up specifying a budget that does not exist.

### A.2 Supersession record — proposed amendment (P)

**Older clauses:** §17.7.1, which defines `ProcessAlive` as true for `starting`, `handshaking`, `unloaded`, `loading`, `ready`, `busy`, `degraded`, `stopping` "while the process is alive", and `Available` as false for `unavailable`.

**Newer requirements:** A2-D4 (line 2505) — `ProcessAlive` derives solely from `TTSProcessState`, true for `starting`/`alive`/`stopping`. A2-D5 (line 2506) — `Available` is a stored class, true for asset-specific `unavailable`.

**Conflict:** the two definitions cannot both be satisfied. The implementation matches the **older** text exactly (`supervisor.go:325`, `:348–356`, `:358–365`), which is why the code looks non-conformant only relative to the newer decision.

**Governing provision after acceptance:** A2-D4 and A2-D5. `ProcessAlive` is a function of process state alone; `Available` is a stored class plus derivation, never inferred from the provider state alone.

**Rationale to preserve:** the older text conflated "process exists" with "helper is usable", which is what allowed a live process under `unavailable` to be reported dead, and a dead process under `ready` to be reported alive.

**Resulting behaviour:**

| Process state | `ProcessAlive` |
|---|---|
| `absent`, `exited` | false |
| `starting`, `alive`, `stopping` | true |

**Proposed replacement text for §17.7.1:** replace both derivation paragraphs with: *"`ProcessAlive` is derived solely from `TTSProcessState`: true for `starting`, `alive`, and `stopping`; false for `absent` and `exited`. `Available` is the stored availability class, true for `unloaded`, `loading`, `ready`, `busy`, `degraded`, and asset-specific `unavailable`; false for `absent`, `starting`, `handshaking`, `stopping`, `crashed`, `incompatible`, and configuration-level `unavailable`."*

**Also required:** the `backoff` rows in §17.7.2 must move to a deferred appendix under O-13. Superseding availability while leaving `crashed → backoff` normative re-creates the same conflict.

## Part B — Ownership evaluation (proposal, not approved)

### B.1 Q1 — Who validates each event

Two distinct validations, currently absent:

- **Event legality — `Supervisor`.** Only the supervisor knows process state, active request, and intent. A `synthesizeStarted` while `processState == exited` must be rejected by the supervisor, not by the client.
- **Request admissibility — `Client`.** Capability presence, identity well-formedness, output path, speed bounds are request concerns and stay in the client (`client.go:688–696` today).

The split matters because the client currently performs both by naming states, which is why 21 call sites can each write any state.

### B.2 Q2 — Rejected or normalized

**Rejected, with a typed failure and no state write.** Normalizing makes an illegal transition indistinguishable from a legal one and would erase the exact evidence W-11 depends on. Two exceptions, both explicit:

- idempotent event (`startRequested` while `alive`) returns success with no transition;
- termination events are always accepted and are the only events allowed from any state.

### B.3 Q3 — `LastError` and `ActiveRequestID`

| Event family | `LastError` | `ActiveRequestID` |
|---|---|---|
| `startRequested`/`handshakeSucceeded` | cleared | cleared |
| `loadStarted` | cleared | cleared |
| `loadFailed` | set | cleared |
| `synthesizeStarted` | cleared | set |
| `synthesizeSucceeded` | cleared | cleared |
| `synthesizeFailed` | set | cleared |
| `unloadCompleted` | preserved | cleared |
| termination (any intent) | preserved for handshake/protocol; set for nonzero exit | cleared |
| identity change | cleared | cleared |

Open point: whether a preserved `LastError` should survive the next successful load. Recommended: clear on success, preserve on unload only. O-14.

### B.4 Q4 — Serialization of termination against requests

Verified mechanics:

- `s.mu` serializes individual transitions: `observeExit`, `SetState`, `MarkHandshake`, `SetSelected`, `Terminate` (`supervisor.go:66`, `:170–192`, `:216–243`, `:246–259`).
- `startMu` serializes startup (`supervisor.go:117`, `client.go:198`).
- `workMu` serializes synthesis and catalog work (`client.go:142`, `:518`, `:685`).
- **`Shutdown` holds no `workMu`** (`client.go:645–674`), and `Terminate` does not know about it. A shutdown or cleanup termination can therefore complete while a synthesis is in flight, and the in-flight terminal `SetState(ready)` (`:771`, `:783`, `:793`, `:803`, `:839`) then overwrites the exit result.

So transitions are individually atomic, but **lifecycle completion is not serialized against requests**. The required fix is not a bigger lock — it is a lifecycle epoch: every event carries the epoch it was issued under, and events from a stale epoch are dropped. This is the same mechanism the TargetID generation question needs (§B.6), which is why the two must be decided together.

### B.5 Q5 — Unambiguous tuples

Ambiguous today: protocol violation (exit code 4) versus unexpected crash both currently land on `crashed`; cleanup when the process is already absent; `Terminate` with `proc == nil` returns without writing state (`:252–254`), so a stale `ready` with an absent process is not corrected by shutdown.

Required: one tuple per (intent, observed exit) pair, with no fall-through. The five intents are enumerated in the reconcile §4.4; the observer algorithm remains out of scope.

### B.6 `crashed` options — assess, do not choose

| Option | Caller impact | UI impact | Recovery | Compatibility |
|---|---|---|---|---|
| (a) keep state, set `LastError` | none | no new state; I/O failure looks like `degraded` | explicit request | none |
| (b) internal `ioFailed` flag | none | needs a projection decision | explicit request | new field, additive |
| (c) new `io_failed` state | none (no external implementers) | new badge | explicit request | vocabulary amendment |

No external implementer of the status interface exists (§C.3), so (c) is cheaper than it looks; (a) is cheapest. Remains O-15 and is **not chosen** here.

## Part C — Inventory completion

### C.1 API coverage, corrected

Ten exported `Client` methods. Direct production call sites outside the package: **none verified** for any of them. The earlier counts were false positives on receiver names:

| Apparent hit | Actual receiver |
|---|---|
| `.Load(` ×8 | `atomic.Bool.Load`, `store.Load`, settings store |
| `.Start(` ×9 | `cmd.Start`, `jobManager.Start`, `proc.Start` |
| `.Shutdown(` ×4 | `imageEngine.Shutdown`, `Pipeline.Shutdown` |
| `.Catalog(` ×2 | `capability.Catalog` |
| `.Health(` | `httpClient.Health` |
| `.Unload(` ×2 | `loader.Unload` |
| `.Cancel(` | `jobManager.Cancel` |

Reachability was then checked directly: **no file outside `internal/audio/tts/helper` imports the package** (0 matches). The entire exported surface is reachable only from the package's own tests. This is not a defect claim — C1 implementation has not been authorized — but it does mean no operation in this inventory is currently exercised in production, and every side-effect classification is a statement about code paths, not about observed runtime behaviour. O-21 resolved as **not integrated**.

### C.2 Side effects must be recorded per surface

The matrix previously mixed direct and C7 surfaces, which is how `Catalog` got classified as observational. Side effects differ by surface, so each row needs a surface column.

| Operation | Surface | Spawns | Identity | State writes |
|---|---|---|---|---|
| `Start` | direct | yes | none | `starting`→`handshaking`→`unloaded`/`incompatible`/`unavailable` |
| `Health` | direct | **yes** (`:501`→`:389`) | none | via start |
| `Catalog` | direct | **yes** (`:520`) | **yes** (`:546`) | via start |
| `Load` | direct | yes (`:554`) | **yes** (`:593`) | `loading`/`unavailable`/`ready` |
| `Unload` | direct | yes (`call`) | none | `unloaded` |
| `Cancel` | direct | yes (`call`) | none | none |
| `Shutdown` | direct | no | none | `stopping`→`absent` via terminate; no `workMu` |
| `Synthesize` | direct | yes (`:697`) | none | `busy`/`ready`/`degraded`/`crashed` |
| `ProviderStatus` | direct | **yes** (`:953`) | none | via start |
| `ProviderCatalog` | direct | **yes** (delegates) | **yes** | via start |

`Start`, `Load`, `Shutdown` are now in the inventory. `Voice-list` stays out: it is not a `Client` method.

### C.3 C7 is design-only

`types/provider.go:113–114` declares `ProviderStatus`/`ProviderCatalog`; only `Client` implements them, and no adapter type does. Every C7 row is therefore unimplemented by construction, including the read-only guarantees at issue in R-24/R-25.

## Part D — Decisions and status

| # | Decision | Recommendation |
|---|---|---|
| O-13 | restart budget / `backoff` scope | defer; amendment records future functionality |
| O-14 | event owner, `SetState` unexported, `LastError` rules | accept semantics, defer implementation |
| O-15 | `crashed` split | (a) for C1, (b) if signal needed |
| O-16 | five termination tuples | accept |
| O-17 | per-surface side-effect classification | accept |
| O-18 | C7 design-only | accept |
| O-19 | `Loaded`/capability fields scope | defer to the matrix revision |
| O-20 | replacement / idle close / supervisor cancel scope | decide before r9 |
| O-21 | helper integration | resolved: no external import; unintegrated |

All fourteen W findings remain open. No finding is closed by a disposition row. R-38/R-39 remain absent. No matrix regenerated, no exit-observer specification, no design edit, no code change, nothing built or run.

**C1 remains unauthorized.**
