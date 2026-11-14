# C1 — State-Writer and Transition Audit (read-only)

**Date:** 2026-11-08
**Status:** SOURCE AUDIT. Not r9 design text. **C1 REMAINS UNAUTHORIZED.**
**Inputs preserved:** A1-rev1 `c061e9d5…`, A2-rev1 `26342c5e…`, A1 `92acde55…`, A2 `0a5ea047…`, r8 design `a8c40fac…`, r9 plan `99897422…`.
**Scope:** trace every provider-state and process-state writer in production code, plus callers, before specifying the intent-aware exit observer.

Classes: **F** source-verified, **R** reported, **P** proposal, **O** owner decision.

## 1. Method and limits

Searched: `internal/audio/**`, `app-*.go`, `--include=*.go`, production and test separated by filename.

Patterns: all 13 `TTSState*` and 5 `TTSProcess*` constants; `SetState(`, `setUnavailable(`, `MarkHandshake(`, `SetSelected(`, direct `providerState =` / `processState =`, struct-literal fields; `Terminate(` / `Kill(` for termination intents; `backoff` / `restart budget` for the recovery mechanism.

Limits, stated so the audit is not overread:

- static reading only; no execution, no race exercise;
- call sites inside vendor/build output not searched;
- test files examined for writers, not behaviour;
- "no writer" means no assignment found by these patterns in the searched paths (F for the search, not a proof about uninspected code).

## 2. Errata against A1-rev1 and A2-rev1 (originals preserved)

Three claims in the revisions I authored are wrong or overstated. They are corrected here, not silently edited, so the hashes stay meaningful.

| # | Claim in revision | Correction | Class |
|---|---|---|---|
| E1 | A1-rev1/A2-rev1: `setUnavailable` writes process `absent` **while `proc != nil`**, so a live process is reported dead | `setUnavailable` is called only at `supervisor.go:130` (discovery failure) and `:149` (spawn failure). Neither can have `proc != nil`. The defect is a **latent hazard** in the helper, not an actual live-process bug | F |
| E2 | A2-rev1 §4: `incompatible`, `busy`, `degraded` may have no writer; tracing outstanding | All three have production writers: `incompatible` at `client.go:251`, `busy` at `:737`, `degraded` at `:796,806,814,820,828` | F |
| E3 | A2-rev1 §4: `backoff` writer unknown | `backoff` has **no writer anywhere**; the string occurs once, as the constant definition. There is no restart budget implementation | F |

E3 is the largest: the contract's `backoff`/budget machinery does not exist in this build, so A2-D5 clause (b) and A2-D4's budget language currently describe an unreachable state.

## 3. Complete writer inventory (F)

### 3.1 Provider-state writers in `supervisor.go`

| Site | Writes | Trigger | Process effect |
|---|---|---|---|
| `:104` | `absent` | `NewSupervisor` | `absent` |
| `:137` | `starting` | `Start`, before spawn | `starting` |
| `:159` | `handshaking` | `Start`, spawn accepted | `starting` |
| `:182` | `absent` | `observeExit`, `providerState == stopping`, **both exit codes** | `absent` (0) / `exited` (nonzero) |
| `:185` | `crashed` | `observeExit`, any other state | `exited` |
| `:200` | `unavailable` | `setUnavailable` (§3.3) | `absent` |
| `:213` | `unloaded` | `MarkHandshake` | `alive` |
| `:236` | **caller-supplied** | `SetState` | `alive` if `proc != nil` and state ≠ `stopping` |
| `:258` | `stopping` | `Terminate` with `proc != nil` | `stopping` |

`SetState` is exported and accepts any of the thirteen states. It is the only generic writer, and it has **21 production call sites**, so "who publishes this transition" is currently "whichever caller last invoked it".

### 3.2 `SetState` call sites (`client.go`, F)

| Group | Lines | State |
|---|---|---|
| handshake failure | `:251`, `:253` | `incompatible`, `unavailable` |
| load | `:578`, `:581`, `:586`, `:590`, `:594` | `loading`, `unavailable` ×3, `ready` |
| unload | `:616` | `unloaded` |
| synthesis start | `:737` | `busy` |
| encode failure | `:748` | `ready` |
| transport write failure | `:756` | **`crashed`** |
| timeout / cancel | `:771`, `:783`, `:793`, `:803`, `:839` | `ready` |
| helper failure / bad result | `:796`, `:806`, `:814`, `:820`, `:828` | `degraded` |
| synthesis success | `:847` | `ready` |

### 3.3 `setUnavailable` call sites

`:130` discovery failure, `:149` spawn failure. Both with `proc == nil` (see E1).

### 3.4 Process-state writers

`:103` absent · `:136` starting · `:158` starting · `:178` absent · `:180` exited · `:184` exited · `:199` absent · `:212` alive · `:240` alive · `:257` stopping.

### 3.5 Test-only writers

None found. Tests assert on state; they do not assign it.

## 4. Per-state paths (F)

| State | Created by | Exits via | On normal exit | On nonzero exit | On termination | On restart |
|---|---|---|---|---|---|---|
| `absent` | ctor, exit0, reset | Start | — | — | — | → `starting` |
| `starting` | `Start` | spawn ok/fail | n/a | → `unavailable` | → `stopping`→`absent` | — |
| `handshaking` | `Start` | handshake ok/fail | n/a | → `incompatible`/`unavailable` | → `stopping`→`absent` | — |
| `unloaded` | `MarkHandshake`, `Unload` | Load, Synthesize | → `loading` | → `stopping`→`absent` | → `absent` | → `starting` |
| `loading` | `Load` | load result | → `ready` | → `unavailable` (process alive) | → `absent` | — |
| `ready` | `Load` ok, synth terminal | Synthesize, Unload | → `absent` | → **`absent`** (defect D1) | → `absent` | → `starting` |
| `busy` | Synthesize | terminal | → `ready`/`degraded` | → `stopping`→`absent` | → `absent` | — |
| `degraded` | synth/result failure | Synthesize, Unload | → `absent` | → `absent` | → `absent` | — |
| `stopping` | `Terminate` | exit observer | → `absent` | → **`absent`** (must be `crashed`) | — | — |
| `crashed` | `observeExit`, Send failure | Start | — | — | → `stopping`→`absent` | → `starting` |
| `incompatible` | handshake failure | Start | — | → **`absent`** unless the client restores it | → `absent` | → `starting` |
| `unavailable` | `setUnavailable`, load/handshake failure | Load | — | → `absent` | → `absent` | → `starting` |
| `backoff` | **nothing** | — | — | — | — | — |

## 5. Termination intents (F)

Five distinct reasons call `Terminate`/`Kill`, and only one restores a typed state afterward:

| Intent | Site | Post-exit state |
|---|---|---|
| failed handshake | `client.go:248` then `:251/:253` | `incompatible`/`unavailable` **restored by the client** |
| protocol violation | `client.go:305/312/366` | `absent` |
| shutdown / idle close | `client.go:654–668` | `absent` (nonzero → `absent`, defect) |
| synthesis cleanup | `client.go:861` (`cleanupRejected`) | `absent` |
| supervisor kill | `Terminate`/`Kill` | `absent` |

The handshake path carries an explicit comment about preserving the typed failure. It is a **per-path patch**, not an ownership rule — which is why the other four intents lose state.

## 6. Findings (F unless marked)

| # | Finding | Evidence |
|---|---|---|
| W-1 | `backoff` is unreachable; no restart budget exists | constant only; `"backoff"` literal count 1 |
| W-2 | Nonzero exit during `stopping` becomes `absent`, violating A2-D8 | `:176–182` |
| W-3 | Exit observer is intent-blind; four of five intents erase the prior state | §5 |
| W-4 | `Health` starts the helper, contradicting the matrix row that says "no start" | `client.go:501` → `:389` → `Start` |
| W-5 | Direct `Catalog` starts the helper, contradicting C1-AC-12 and the §4.7 table | `client.go:520` |
| W-6 | `ProviderCatalog` delegates to `Catalog`, therefore starts **and** mutates selection | `:961–962`, `:546` |
| W-7 | `ProcessAlive` has two inconsistent definitions: method ignores aggregate state (`:341`), status field subtracts `crashed`/`unavailable` (`:325`, `:348–356`) | same fact, two answers |
| W-8 | `providerAvailable` excludes every `unavailable` | `:358–365` |
| W-9 | Same provider state, two liveness values: load failure → `unavailable` with process `alive`; discovery/spawn failure → `unavailable` with process `absent` | `:578–590` vs `:199` |
| W-10 | `crashed` has two meanings: unexpected exit (`:185`) and transport write failure (`:756`, process possibly alive) | both write `crashed` |
| W-11 | State can be resurrected: timeout/cancel writes `ready` even after a failure path set `degraded`/`crashed` | `:771,783,793,803,839` |
| W-12 | `cleanupRejected` terminates the helper as an undeclared state operation; a failed request can erase `ready`/`degraded` | `:856–864` |
| W-13 | Three exported operations are absent from the matrix: `Start`, `Load`, `Shutdown` | §3 of A2-rev1 vs API list |
| W-14 | `Cancel` never writes provider state; a cancel with no owner terminal leaves `busy` until the owner finishes | `:622–644` |

## 7. What the exit observer needs before it can be specified

Not a design. The prerequisites the trace establishes:

1. **An intent enum** with one value per §5 row, passed to the termination path; the observer branches on intent, not on the current provider state.
2. **A single owner for the exit tuple** — provider state, process state, liveness — so the client's post-hoc `SetState` restore is removed rather than replicated.
3. **A rule for exit codes per intent**: handshake → typed failure; cleanup → prior state or `degraded`; shutdown → `absent`/`crashed` by code; protocol violation → its own terminal state; unexpected → `crashed`.
4. **A decision on `crashed` split** (W-10): process death versus transport failure are different facts and cannot share one state.
5. **Resolution of W-4/W-5/W-6** before any matrix is authoritative: current read paths spawn and mutate selection.

## 8. TargetID lifecycle — questions the trace cannot answer

The trace found no target identity in source; `TargetID` exists only in A1-rev1. All five owner questions are therefore **P**, not F:

| Question | State |
|---|---|
| who creates the instance nonce | undecided; must be supervisor construction |
| who owns and advances generation | undecided; must be the component that replaces identity, not the issuer |
| can a generation change race with issuance or consumption | must be answered with the §4.4 lock order plus a generation read at consume time |
| does `Consume` validate against live supervisor identity, or only ticket data | **only ticket data today by design**; A1-rev1 does not yet require a live check — this is an open gap |
| can a supervisor be replaced while an old ticket stays valid | yes unless generation is revalidated at consume; the current proposal does not close it |

The fourth row is the one that matters: comparing a ticket's `TargetID` to a stored record is not the same as validating against the live identity. A1-rev1 should not be read as having settled it.

## 9. Still open

R-38/R-39 (no text). W-1 through W-14 disposition (none accepted). 234-disposition matrix (O-9, not approved). The exit-observer intent model (prerequisites only). TargetID live validation (§8). No r9 text written.

## 10. Authorization

**C1 remains unauthorized.** Read-only audit: no design, source, test, dependency or configuration changed; nothing built or run; 142 pre-existing dirty entries untouched.
