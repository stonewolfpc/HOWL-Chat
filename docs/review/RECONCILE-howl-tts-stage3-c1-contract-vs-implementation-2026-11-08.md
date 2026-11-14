# C1 — Contract-to-Implementation Reconciliation

**Date:** 2026-11-08
**Status:** READ-ONLY RECONCILIATION. Not r9 design text. **C1 REMAINS UNAUTHORIZED.**
**Prerequisites addressed:** (1) contract vs actual functionality, (2) state ownership and semantics, (3) operation inventory.
**Not done here, deliberately:** exit-observer specification, matrix regeneration, any design edit.
**Preserved inputs:** A1-rev1 `c061e9d5…`, A2-rev1 `26342c5e…`, audit `1541332e…`, A1 `92acde55…`, A2 `0a5ea047…`, r8 `a8c40fac…`, plan `99897422…`.

Classes: **V** verified existing behaviour · **C** required change to existing implementation · **N** specified but not implemented (proposed functionality) · **X** superseded text still present · **U** unresolved.

## 1. Method

Cross-read the governing contract (`DESIGN-howl-tts-helper-architecture-2026-11-06.md` §17.7.1, §17.7.2, §17.11.4, §17.11.5, §17.12.1, §17.19 A2-D4..D12) against the writer audit and the helper source. Each row states which of V/C/N/X applies and the exact evidence.

Limits: static reading; no execution; no race exercise; the classification is about *this build*, not about intent.

## 2. The restart-budget finding, stated exactly

The contract does specify a budget — §17.11.5: *"Backoff is 1 s, 2 s, then 5 s; maximum 3 consecutive restarts per key in 10 minutes; reset after 10 minutes of stable readiness."* §17.7.2 specifies `crashed → backoff`, `backoff → starting`, and `budget exhausted → unavailable`.

**None of it exists.** `backoff` has no writer, there is no counter, no schedule, no per-key bookkeeping, and no stable-readiness clock. `"backoff"` occurs once in the production tree, as the constant.

Consequence: A2-D4's "does not enter `backoff`, spend restart budget" is true only because nothing enters `backoff`, and A2-D5 clause (b) describes a state that cannot be reached. Neither may be implemented by inventing the subsystem inside C1; the scope decision is O-13.

## 3. Clause reconciliation

| # | Clause | Required | Current | Class | Note |
|---|---|---|---|---|---|
| R-1 | §17.11.5 restart schedule 1/2/5 s | budget + schedule | absent entirely | **N** | scope decision O-13 |
| R-2 | §17.11.5 max 3 per key / 10 min, reset after 10 min | counter | absent | **N** | same |
| R-3 | §17.7.2 `crashed → backoff` | transition | no writer | **N** | unreachable |
| R-4 | §17.7.2 `backoff → starting` | transition | no writer | **N** | unreachable |
| R-5 | §17.7.2 budget exhausted → `unavailable` | transition | absent | **N** | same |
| R-6 | A2-D5(b) `Available=true` for budget-exhausted | derivation | depends on R-1..R-5 | **N** | cannot be implemented alone |
| R-7 | §17.7.2 `stopping` exit 0 → `absent` | transition | `observeExit:178–182` | **V** | — |
| R-8 | §17.7.2 `stopping` grace expires → `crashed` | transition | nonzero/expired → `absent` | **C** | `:176–182` |
| R-9 | §17.11.5 exit codes 0/2/3/4/other | classify by code | only `err == nil` vs non-nil | **C** | codes never inspected |
| R-10 | §17.7.2 unexpected `Wait()` → `crashed` | transition | `:184–185` | **V** | — |
| R-11 | §17.7.2 semantic hello failure → `incompatible`, terminate, no restart | transition | `client.go:248–253` | **V** | achieved by post-hoc restore, not ownership |
| R-12 | §17.7.2 `starting` start/job failure → `crashed` **or** `unavailable` | either | always `unavailable` (`:199`) | **C** | `crashed` branch unused |
| R-13 | A2-D4 `ProcessAlive` solely from `TTSProcessState` | single source | `:325`, `:348–356` subtract provider state | **C** | contradicts R-13's own text |
| R-14 | §17.7.1 old `ProcessAlive`/`Available` definitions | — | still present in the contract | **X** | A2-D4/A2-D5 supersede, text not removed |
| R-15 | A2-D5 `Available` true for unloaded/loading/ready/busy/degraded | derivation | `providerAvailable:358–365` | **V** | — |
| R-16 | A2-D5 `Available` true for asset-specific `unavailable` | derivation | false | **C** | needs the stored class |
| R-17 | A2-D5 `Loaded` false on exit/crash/unload/identity change | field | field absent | **N** | depends on A2-D6 |
| R-18 | A2-D6 `Loaded`, `ServiceCapabilities`, `ServiceCapabilitiesKnown` | fields | all absent (`types/provider.go:56–71`) | **N** | separate criterion (O-7) |
| R-19 | A2-D7 identity replacement sequence | 7-step sequence | not implemented | **N** | one textual `replacement` match, no sequence |
| R-20 | A2-D8 idle close | trigger + timer | no idle mechanism | **N** | `Shutdown` exists; no trigger |
| R-21 | A2-D9 cancel without `cancel` capability | supervisor path | `cancelForSupervisor` absent from source | **N** | `Client.Cancel` requires the capability (**C**) |
| R-22 | A2-D12 explicit-load recovery | `unavailable`→`loading` | `Load` sends `load` without a state gate (`:551–594`) | **V/C** | works incidentally; update-gate and identity checks absent |
| R-23 | §17.12.1 ensure-handshake **and load** on synthesize | ensure-load | `Synthesize` calls `Start`, never `Load` | **N** | `:679–737` |
| R-24 | C1-AC-12 direct `Catalog` starts nothing | no start | `:517–520` calls `Start` | **C** | hidden side effect |
| R-25 | C7 `ProviderCatalog` read-only, no identity change | no side effects | `:961–962` → `:546` `SetSelected` | **C** | starts too |
| R-26 | §17.11.4 step 6 cancel grace → `KillTree` | kill bound | `watchCancellation` + `cleanupRejected` terminate | **V** | path exists |
| R-27 | §17.11.4 step 7 delete only after observed termination | ordering | `cleanupRejected:856–864` | **V** | — |

### 3.1 What the classification rules out

- No restart-budget subsystem may be specified as "current behaviour". It is **N**.
- No matrix may treat `backoff` as reachable. It is **N**.
- `Available`-for-`unavailable` and `Loaded` are **C**/**N** and depend on A2-D6; they cannot be verified from the current status shape.
- The contract contains **two** `ProcessAlive`/`Available` definitions (§17.7.1 vs A2-D4/D5). Until supersession is explicit in the normative text, both remain quotable, which is how the implementation came to match the older one.

## 4. State ownership and semantics

### 4.1 Publishers today (F)

Nine write sites in `supervisor.go`, one of them `SetState(state, …)` accepting any of the thirteen states, with 21 production call sites in `client.go`. No caller outside the helper package (**F**, 0 matches), so the generic setter is exported without an external consumer.

The handshake path's post-termination restore (`client.go:248–253`) is the only place a caller compensates for generic exit behaviour.

### 4.2 Proposed single owner (P)

One unexported transition entry point on `Supervisor`, taking an event rather than a state:

```text
Supervisor.apply(event)   // event ∈ enumerated set; caller never names a state
```

- `SetState` is removed from the exported surface and becomes unreachable outside the package (mechanically safe: 0 external callers, **F**).
- Each event carries the tuple it produces: provider state, process state, class, `LastError` disposition, active request ID.
- The client asks for events (`loadStarted`, `loadFailed`, `synthesizeStarted`, `synthesizeFailed`, `synthesizeDone`, …), never for states.

This is a proposal (P). It is not r9 text and does not specify the exit observer.

### 4.3 `crashed` must be split (P, O-15)

Two different facts currently write `crashed`: process death (`:185`) and transport write failure (`:756`, process possibly alive). The public vocabulary is frozen at thirteen states (A2-D6), so adding a state requires an amendment. Options:

| Option | Effect | Cost |
|---|---|---|
| (a) transport failure keeps state, records `LastError` | no vocabulary change; state stops lying | loses "crashed" signal for I/O |
| (b) internal flag `ioFailed` alongside `crashed` | vocabulary change, accurate | new field, new projections |
| (c) amend vocabulary with `io_failed` | most accurate | contract amendment |

### 4.4 Five termination intents — required resulting tuples (P)

| Intent | Provider | Process | LastError | Loaded |
|---|---|---|---|---|
| handshake failure | `incompatible` / `unavailable` | `exited` (absent if never spawned) | preserved typed failure | false |
| protocol violation | `crashed` | `exited` | `protocol_violation` | false |
| shutdown / idle close | exit 0 → `absent`; nonzero → `crashed` | `absent` / `exited` | on nonzero only | false |
| request cleanup | `absent` (process gone) | `absent` / `exited` | request failure | false |
| supervisor kill | exit 0 → `absent`; nonzero → `crashed` | as observed | on nonzero | false |

The first row is what the client already achieves by patching. The other four are the behaviour to be **defined**, not currently implemented (**N**). This table is a semantic requirement; the observer algorithm is still out of scope.

### 4.5 Liveness invariant (C)

There must be exactly one answer to "is the helper process alive", derived from `processState` alone, and both `Supervisor.ProcessAlive()` and `TTSProviderStatus.ProcessAlive` must return it. Today they disagree whenever provider state is `crashed` or `unavailable` with a live process — and internal cleanup acts on the method while the UI would read the field.

## 5. Operation inventory reconciled

Fourteen operations; side effects marked from source.

| Op | Spawns | Wire frames | Identity effect | Provider-state writes | Budget |
|---|---|---|---|---|---|
| `Start` | yes | `hello` | none | `starting`, `handshaking`, `unloaded`, `unavailable`, `incompatible` | none |
| `Health` | **yes** (via `call`→`Start`) | `health` | none | via start only | consumes |
| `Catalog` | **yes** (`:520`) | `catalog` | **`SetSelected`** (`:546`) | via start only | consumes |
| `Load` | yes (`:554`) | `load` | **`SetSelected`** (`:593`) | `loading`, `unavailable`, `ready` | consumes |
| `Unload` | yes (via `call`) | `unload` | none | `unloaded` | consumes |
| `Cancel` | yes (via `call`) | `cancel` | none | none | consumes |
| `Shutdown` | no | `shutdown` | none | `stopping`→`absent` via terminate | consumes |
| `Synthesize` | yes (`:697`) | `synthesize` | none | `busy`, `ready`, `degraded`, `crashed` | consumes |
| `ProviderStatus` | **yes** (`:953`) | none | none | via start only | none |
| `ProviderCatalog` | **yes** (delegates) | `catalog` | **`SetSelected`** | via start only | consumes |
| supervisor cancel | n/a | `cancel` | none | none | none — **not implemented** |
| identity replacement | n/a | `shutdown`,`load` | replaces identity | full sequence | **not implemented** |
| update-gated reset | n/a | none | none | `absent` | **not implemented** |
| idle close | n/a | `shutdown` | none | `stopping`→`absent` | **not implemented** |

Corrections this forces on the previous inventory: `Health`, `Catalog`, and `ProviderStatus` are **not** observational — all three can spawn. `Catalog` and `ProviderCatalog` mutate selected identity. `Voice-list` is an application-level operation, not a `Client` method, and should leave the client matrix. `Load` has no production caller, so its row is API coverage, not an exercised path.

C7: no type implements `ProviderStatus`/`ProviderCatalog` from `types/provider.go:113–114` other than `Client` itself (**F**), so every C7-adapter row is design-only (**N**).

## 6. Findings W-1..W-14 — all still open

| Finding | Status | Needs |
|---|---|---|
| W-1 `backoff` unreachable | open | O-13 scope decision |
| W-2 nonzero exit → `absent` | open | R-8/R-9 code change |
| W-3 intent-blind observer | open | §4.4 semantics + O-14 |
| W-4 `Health` spawns | open | R-24 style classification fix |
| W-5 `Catalog` spawns | open | same |
| W-6 `ProviderCatalog` spawns + selects | open | C7 boundary decision |
| W-7 `ProcessAlive` disagreement | open | §4.5 |
| W-8 `Available` excludes `unavailable` | open | R-16 |
| W-9 two liveness values for one state | open | §4.2, §4.5 |
| W-10 `crashed` ambiguity | open | O-15 |
| W-11 state resurrection on timeout/cancel | open | §4.2 |
| W-12 cleanup termination undeclared | open | §4.4 |
| W-13 `Start`/`Load`/`Shutdown` absent from matrix | open | §5 |
| W-14 `Cancel` writes no state | open | §4.2 |

No finding is closed by this document. Each needs its source evidence reviewed independently.

## 7. Owner decisions

| # | Decision | Recommended |
|---|---|---|
| O-13 | restart budget: remove from Stage 3 C1 scope, or implement | remove from C1; amendment records it as future |
| O-14 | one transition owner, `SetState` unexported | accept |
| O-15 | split `crashed` from transport failure: option (a)/(b)/(c) | (a) now, (b) if a signal is required |
| O-16 | five termination intents' tuples (§4.4) | accept |
| O-17 | `Health`/`Catalog`/`ProviderStatus` classified as may-spawn | accept |
| O-18 | C7 adapter is not implemented; all C7 rows are design-only | accept |
| O-8 | R-38/R-39 original text | still required |

## 8. What this does not do

No exit-observer specification. No matrix regeneration. No design edit. No source, test, dependency or configuration change. Nothing built or run.

## 9. Authorization

**C1 remains unauthorized.** The read-only boundary holds; 142 pre-existing dirty entries untouched.
