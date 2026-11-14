# C1 Revision 9 — Addendum A2, Revision 1

**Date:** 2026-11-08
**Status:** CORRECTION PACKAGE, review input only. **C1 REMAINS UNAUTHORIZED.**
**Supersedes:** `PLAN-howl-tts-stage3-c1-r9-correction-addendum-a2-2026-11-08.md`, sha256 `0a5ea0479651e60bf48589dac912b50d01b9159f8af23ffa8ef9d13637b71ea0` (190 lines, 13,282 B). Original preserved.
**Companion:** `...addendum-a1-rev1-2026-11-08.md`.
**Unchanged inputs:** r8 design `a8c40fac…`; r9 plan `99897422…`.

Classes: **F** source-verified, **R** reported, **P** proposal, **O** owner decision.

## 1. What revision 1 changes

A2 correctly rejected the 48-cell matrix's coverage, then proposed 60 cells from the same foundation — a larger table with the same structural defect. This revision:

1. derives the operation inventory from the **actual exported API**, not from the existing table (§2);
2. derives the lifecycle dimension from the **actual state constants**, not from six behavioural classes (§3);
3. makes optional capabilities a **modifier**, not a third dimension (§4);
4. computes the matrix size from those dimensions and abandons both 48 and 60 (§5);
5. rebuilds the availability writer model on the actual transition code (§6).

## 2. Reconciliation — A2 findings

| ID | Finding | Evidence | Required correction | Author response | Independent disposition |
|---|---|---|---|---|---|
| **A2-01** | 60-cell minimum built on the wrong foundation; inventory must come from the exported API and planned boundary, lifecycle from actual states, optional capabilities separate, size computed after | A2 §6.3: "10 operations × 6 classes = 60", operations inherited from the r7 table rather than enumerated | derive dimensions from source, then compute | §2–§5 | Pending |
| **A2-02** | A2 committed the count-as-coverage error it had just identified, by asserting a total without per-dimension provenance | A2 §6.2 defect 1 vs §6.3 recommendation | state each dimension's source and recompute | §5 | Pending |
| **A2-03** | `PairID` startup assertion does not establish targeting; it only proves the two handles came from one authority | A2 §2.1 | target identity + generation; demote `PairID` to wiring check | A1-rev1 §3 | Pending |
| **A2-04** | Multi-supervisor independent consumption was asserted in a test (A6) with no mechanism behind it | A2 §3.4 test A6 | per-target registry keyed by target identity | A1-rev1 §3.2 | Pending |
| **A2-05** | A2 depended on availability transition semantics it never checked; the source contradicts the contract in five places | `supervisor.go:166–191`, `:197–202`, `:233–243`, `:325`, `:335`; contract A2-D4, A2-D5, A2-D8 | rebuild the writer model on the transition table | §6 | Pending |
| **A2-06** | "No production consumer" was labelled F although it is a bounded negative search; the search scope, not the conclusion, is what was verified | A2 §5.1 pattern list | restate as scope-bounded, name patterns and directories | §6.1 | Pending |

## 3. Operation inventory from the exported API (F)

`Client` exported methods, enumerated from source (`internal/audio/tts/helper/client.go`):

| # | Method | Line | Wire operation | Class |
|---|---|---|---|---|
| 1 | `Start` | 191 | `hello` (via handshake) | operation |
| 2 | `Health` | 501 | `health` | operation |
| 3 | `Catalog` | 517 | `catalog` | operation, optional-gated |
| 4 | `Load` | 551 | `load` | operation, mandatory |
| 5 | `Unload` | 599 | `unload` | operation, optional-gated |
| 6 | `Cancel` | 622 | `cancel` | operation, optional-gated |
| 7 | `Shutdown` | 645 | `shutdown` | operation, mandatory control |
| 8 | `Synthesize` | 679 | `synthesize` | operation, mandatory |
| 9 | `ProviderStatus` | 952 | none | operation |
| 10 | `ProviderCatalog` | 961 | `catalog` (delegates to `Catalog`) | operation, optional-gated |
| — | `String` | 965 | — | **excluded**: not an operation |

`ProviderCatalog` delegating to `Catalog` (`:962`) is what makes the C1-AC-12 adapter split necessary; they cannot share a row, because one may start a process and the other may not.

**Adapter-level operations** (not on `Client`, still state-changing, therefore still matrix rows):

| # | Path | Source | Class |
|---|---|---|---|
| 11 | supervisor cancellation | `cancelForSupervisor` (§5 of r7 design) | operation |
| 12 | identity replacement | A2-D7 | compound operation |
| 13 | update-gated reset | A2-D12 / §17.7.2 | compound operation |
| 14 | idle close | A2-D8 | compound operation |

**Supervisor methods** (`Start`, `MarkHandshake`, `SetSelected`, `SetState`, `Terminate`, `Kill`, `Stderr`, `Status`, `ProcessAlive`, `Done`) are **state writers, not caller operations**. They appear in §6 as publishers, not as rows.

This is the difference between revision 0 and revision 1: 14 operations exist because 14 are in the source, not because a table had 8 rows and two were missing.

## 4. Lifecycle dimension from the actual constants (F)

Provider states (`types/provider.go:23–35`), thirteen:

`absent`, `starting`, `handshaking`, `unloaded`, `loading`, `ready`, `busy`, `degraded`, `stopping`, `crashed`, `incompatible`, `unavailable`, `backoff`.

Process states (`types/provider.go:11–15`), five: `absent`, `starting`, `alive`, `stopping`, `exited`.

Observations that matter for the matrix:

- `incompatible` has **no writer in `supervisor.go`**; the only provider-state writer there is `SetState`. Who writes it, and whether it survives termination, is **R** until the handshake path is traced.
- `backoff` likewise has no writer in `supervisor.go` (**R**).
- `busy` and `degraded` are reachable only through `SetState` from callers (**R** for exact call sites).
- The six behavioural classes (`A`, `K0`, `K1`, `L`, `X`, `U`) compress these thirteen states by dispatch behaviour. That compression is legitimate for *dispatch*, but it is not a lifecycle model, and it hides the transitions where availability goes stale.

## 5. Matrix sizing, computed (P)

Dimensions, each with a source:

| Dimension | Values | Source |
|---|---|---|
| Operations | 10 | §3, exported `Client` methods |
| Adapter operations | 4 | §3, contract-mandated paths |
| Provider states | 13 | §4, constants |
| Optional-capability modifier | present / absent | A2-D1, governing capability per operation |

Arithmetic:

- base cells = (10 + 4) × 13 = **182**
- optional-absence dispositions = 4 optional-gated client operations × 13 = **52**
- total dispositions = **234**

The modifier is a column, not a multiplied dimension: a helper advertises `catalog`, `unload`, and `cancel` independently, but each operation is gated only by its own capability, so multiplying by 2³ would produce 7/8 meaningless cells.

48 and 60 are both withdrawn. 234 is not a target to hit; it is what the dimensions yield, and it changes if any dimension does.

## 6. Availability writer model, rebuilt (F)

### 6.1 Consumer audit, scope stated (A2-06 correction)

Searched, whole workspace, `--include=*.go`, excluding `.git`: `status.Available`, `Status().Available`, `.ProviderStatus(`, `TTSProviderStatus{`, `TTSVoiceCatalog`, `catalog.Available`.

| Finding | Status |
|---|---|
| only projection site is `supervisor.go:325` | F |
| only construction sites are `supervisor.go:308`, `:312` | F |
| no production consumer of provider-status `Available` outside the helper package | **F for the searched patterns only**; a consumer using a different accessor would not be caught |
| `app-tts-catalog.go` reads `TTSVoiceCatalog.Available`, meaning installed voice models | F |

The absence claim is therefore scope-bounded, not absolute.

### 6.2 Five defects where source contradicts contract

| # | Defect | Source | Contract violated |
|---|---|---|---|
| D1 | `observeExit` overwrites provider state with `absent` on **any** `stopping` exit, regardless of exit code or reason | `:166–191` | A2-D8: nonzero exit MUST NOT silently become `absent`; A2-D4: `incompatible` must remain after an expected termination |
| D2 | `providerAvailable` returns false for every `unavailable` | `:358–365` | A2-D5: true for asset-specific `unavailable`, and for budget-exhausted `unavailable` with `ProcessAlive=false` |
| D3 | `ProcessAlive` is computed from aggregate state as well as process state | `:335`, `:325` | A2-D4: "derived solely from `TTSProcessState` … never derived from aggregate `State`" |
| D4 | Two writers produce `State=unavailable` with opposite process state: `setUnavailable` writes process `absent`; `Load`'s failure path calls `SetState(unavailable)`, which writes process `alive` | `:197–202` vs `:233–243`, `:576–592` | A2-D5/A2-D4: one state, one liveness |
| D5 | `setUnavailable` writes process `absent` while `proc != nil` | `:199` | A2-D4: liveness reflects an owned process |

D1 is the one that makes the previous availability model unsafe to implement as written: a derived-from-state flag will be computed from a state the exit observer has just falsified.

### 6.3 The rebuild: one stored class, derived at read

Add one internal field, `availabilityClass`:

`none` · `assetSpecific` · `budgetExhausted` · `configInvalid`

`Available` is **never stored**. It is a pure function evaluated in `Status()`:

| provider state | available |
|---|---|
| `unloaded`, `loading`, `ready`, `busy`, `degraded` | true |
| `unavailable` + class `assetSpecific` or `budgetExhausted` | true |
| `unavailable` + class `configInvalid` or `none` | false |
| `absent`, `starting`, `handshaking`, `stopping`, `crashed`, `incompatible`, `backoff` | false |

`ProcessAlive` becomes `processState ∈ {starting, alive, stopping}` with no provider argument.

Nothing can go stale, because there is no second write to miss: the only input that can drift is `availabilityClass`, and it is written exactly where the condition is created and cleared on every exit from `unavailable`.

### 6.4 Event table — before/after tuples with a named publisher

| Event | Provider | Process | Class | Publisher |
|---|---|---|---|---|
| construct | absent→absent | absent→absent | none | `NewSupervisor` |
| discovery failure | absent→unavailable | absent→absent | configInvalid | `Start` via `setUnavailable` |
| spawn failure | absent→unavailable | absent→absent | configInvalid | `Start` via `setUnavailable` |
| spawn accepted | absent→starting | absent→starting | none | `Start` |
| handshake begins | starting→handshaking | starting→starting | none | `Start` |
| handshake ok | handshaking→unloaded | starting→alive | none | `MarkHandshake` |
| handshake incompatible | handshaking→incompatible | starting→alive→exited | none | handshake path + exit observer |
| load begins | unloaded/unavailable→loading | alive→alive | none (clear) | `Load` |
| load ok | loading→ready | alive→alive | none | `Load` |
| load failure (asset) | loading→unavailable | alive→alive | assetSpecific | `Load` |
| unload ok | ready/busy/degraded→unloaded | alive→alive | none | `Unload` |
| synthesis | ready→busy→ready | alive→alive | none | `Synthesize` |
| terminate begins | any→stopping | alive→stopping | none | `Terminate` |
| exit expected, code 0 | stopping→absent | stopping→absent | none | exit observer |
| exit expected, nonzero | stopping→**crashed** | stopping→exited | none | exit observer |
| exit unexpected | any→crashed | alive→exited | none | exit observer |
| budget exhausted | crashed/backoff→unavailable | exited→absent | budgetExhausted | restart path |
| identity replacement | any→absent | →exited | none | adapter (A2-D7) |
| update-gated reset | unavailable/incompatible→absent | as observed | none | adapter (§17.7.2) |

Two rows are contract-mandated and **not** what the source does today: the expected-nonzero row (today `absent`) and the incompatible row (today `absent`). Both are corrections for r9, not descriptions of current behaviour.

### 6.5 Invariants

- **I1** `Available` has exactly one publisher: the pure projection in `Status()`.
- **I2** `ProcessAlive` reads only `processState`.
- **I3** `State=unavailable` with `proc != nil` reports `ProcessAlive=true`.
- **I4** `availabilityClass` is non-`none` only while `State=unavailable`.
- **I5** every exit from `unavailable` clears the class.
- **I6** no event sets provider state without a process state in the same tuple.

### 6.6 A2-D6 fields — separate disposition (unchanged)

`Loaded`, `ServiceCapabilities`, `ServiceCapabilitiesKnown` are required by §17.12 and absent from `types/provider.go:56–71` (**F**). They get their own criterion, tests, and review; they are not part of this fix and must not be claimed under it.

## 7. Owner decisions

| # | Decision | Recommended |
|---|---|---|
| O-9 | accept the 234-disposition model as the r9 matrix basis | accept |
| O-10 | accept derived-at-read `Available` + `availabilityClass` | accept |
| O-11 | accept the two contract-mandated corrections (nonzero exit → `crashed`; `incompatible` survives termination) | accept |
| O-12 | trace `incompatible`/`backoff`/`busy`/`degraded` writers (currently R) before r9 text | required |

## 8. Authorization

**C1 remains unauthorized.** No design text, source, test, dependency or configuration changed. Nothing built or run. Original A1/A2 preserved; 142 pre-existing dirty entries untouched.
