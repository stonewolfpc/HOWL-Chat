# C1 Design — Revision 8

**Date:** 2026-11-08
**Status:** DESIGN REVISION — **NOT ACCEPTANCE-READY; C1 UNAUTHORIZED; NO IMPLEMENTATION**
**Supersedes:** only the sections replaced in §4, §5, §6, §7 and §8 of this document. Every other section of revision 7 remains normative and is incorporated by §12.
**Scope:** documentation and design revision only. No source, test, build, runtime, benchmark, or deployment action was taken.
**Parent review:** `docs/review/REVIEW-howl-tts-stage3-c1-r7-2026-11-08.md` (independent, FAIL — narrow).

## 1. Inputs, authority, and evidence classes

| Artifact | SHA-256 | Lines / bytes |
|---|---|---|
| Governing helper contract (§17.19 decisions) | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f` | 2586 / 182,904 |
| Adopted r9 requirements | `0f71cee5e0e5d25b58d1b3956ff549df8cc625ad7f9e8c82189f2f2ef9366c55` | 769 / 112,270 |
| Owner decisions (D10 complete; six pending) | `677f884926de7e58c3dc7ab33f01340f1fbb5966c1b5c089b1aafd057bcf0fba` | 612 / 49,191 |
| D10 finalization | `84fa78dd214f16f40ab125b544727fce2239e66bb1bb54dffc3c8777df9beb02` | 22 / 1,887 |
| C1 design revision 7 (normative for unchanged sections) | `0867acb7216a0e8ee2802b6f08275d997e4138f089ae51acac373c10eea4af54` | 1069 / 161,934 |
| Independent r7 review | `c025ba0a32c776996fd90a59cca0635da5532a43200a5f6b8dad39ff91730f9d` | 200 / 13,835 |

Evidence classes used in every disposition and claim:

| Tag | Meaning |
|---|---|
| `[Observed]` | Read directly from the workspace in this revision. |
| `[Source-derived]` | Follows from the quoted governing text. |
| `[Calculated]` | Obtained by arithmetic or a structural count. |
| `[Inferred]` | A reasoning step beyond the quoted text. |
| `[Unestablished]` | Not determined here; carried openly. |
| `[Superseded]` | Earlier text replaced by this revision; retained historically. |

**Standing rule.** No claim in this document is an implementation property. Every test is a specification that has not been executed.

## 2. Change record: r6 → r7 → r8

| Area | r6 | r7 | r8 |
|---|---|---|---|
| Derived liveness/availability | `processAlive` read provider state; live `unavailable` reported not alive | Predicates separated (§6.8) | Retained unchanged; availability classes completed in §6 |
| Exit classification | Live-`unavailable` exit became `crashed` | Captured intents (§6.9) | Retained unchanged |
| Handshake provenance | `helloHandshake` remained alongside publication | `attempt.publish` sole publication; `installHandshake` sole client writer (§6.11) | Retained unchanged |
| Terminal event set | Two-value success/failure pair | Closed ten-outcome `synthesisTerminal` (§6.10) | Retained unchanged |
| Catalog selection | `SetSelected` removed with no replacement; adapter still delegated | `catalogSelected` + `Adapter.ProviderCatalog`/`CatalogCurrent` (§6.12) | Carrier added (§5): prose restriction becomes enforceable |
| Incompatible barrier | Criterion contradicted `beginTermination` | `beginTermination` limited to `processState` (§6.7) | Retained unchanged |
| Operation matrix | “six × seven”, non-uniform | 48 uniform cells | Retained unchanged |
| `request_not_found` origin | Adapter-synthesised wire code | One current-transport frame (§4.2/§4.5) | Retained unchanged |
| Update reset | Token asserted, undefined | Opaque grant in an authority; **decision still forgeable** | Replaced (§4): authority-verified ticket, no caller-asserted approval |
| §6.6 baseline label | — | Still labelled revision 6 | Corrected (§7) |
| Budget-exhausted class | — | Defined but unexercised | Criterion and test added (§6, §8) |

## 3. Findings register (C1-R-1 … C1-R-37)

Dispositions of R-1 … R-23 are carried from r7 §2.5 unchanged; they are restated here only so the register is complete and every identifier has an explicit disposition.

| Finding | Disposition | Class | Locus / authority |
|---|---|---|---|
| C1-R-1 — start/restart map | design-corrected; C5/C7 execution pending | `[Source-derived]` | r7 §8; r9 C1-START-MAP |
| C1-R-2 — resource/I/O evidence | records complete; independent determinations open (C1-AC-15) | `[Unestablished]` | r7 §9.2–§9.4; r9 §10 |
| C1-R-3 — asset-specific `unavailable` | design-corrected; tests prospective | `[Source-derived]` | r7 §3.4, §6.9; A2-D5/A2-D12 |
| C1-R-4 — semantic handshake termination | design-corrected; execution pending | `[Source-derived]` | r7 §6.4; A2-D3/A2-D4 |
| C1-R-5 — D10 projection | design-corrected; tests prospective | `[Source-derived]` | r7 §4.3; D10 A10a-2/A10b-2 |
| C1-R-6 — traceability | design-corrected | `[Observed]` | r7 §12 |
| C1-R-7 — supervisor known-absent | design-corrected; test prospective | `[Source-derived]` | r7 §5; A2-D9 |
| C1-R-8 — concurrency | design-corrected; race evidence pending | `[Source-derived]` | r7 §9.1 |
| C1-R-9 — protocol/evidence details | design-corrected | `[Source-derived]` | r7 §6.1; A2-D3 |
| C1-R-10 — readiness overstatement | author-reconciled; acceptance pending | `[Observed]` | r7 readiness; this §14 |
| C1-R-11 — direct `Unload` clears `unavailable` | design-corrected | `[Source-derived]` | r7 §4.1; A2-D12 |
| C1-R-12 — incompatible start bypass | design-corrected | `[Source-derived]` | r7 §6.4; A2-D4 |
| C1-R-13 — process-cap citation | corrected | `[Observed]` | contract §5.5 |
| C1-R-14 — AC/HW-1.6 precision | criteria complete; HW-1.6 review pending | `[Unestablished]` | r9 HW-1.6 |
| C1-R-15 — health/Catalog contradiction | design-corrected | `[Source-derived]` | r7 §4.6; A2-D6 |
| C1-R-16 — error semantics | design-corrected | `[Source-derived]` | r7 §4.5; §17.9 |
| C1-R-17 — `MarkHandshake` bypass | closed by r7 §6.5; retained | `[Source-derived]` | r7 §6.5/§6.11 |
| C1-R-18 — incompatible aggregate during termination | closed by r7 §6.7; retained | `[Source-derived]` | r7 §6.7; A2-D4 |
| C1-R-19 — independence wording | closed | `[Observed]` | r7 readiness §8 |
| C1-R-20 — mutation surface | closed by r7 §6.5/§6.6; retained | `[Source-derived]` | r7 §6.6 |
| C1-R-21 — handshake provenance | closed by r7 §6.11; retained | `[Source-derived]` | r7 §6.11 |
| C1-R-22 — catalog start boundary | closed by r7 §4.7/§6.12; carrier added in §5 | `[Source-derived]` | §17.12.1 |
| C1-R-23 — projection matrix | closed by r7 §4.6; retained | `[Calculated]` | 48 cells |
| C1-R-24 — derived liveness/availability | **closed in r7; retained verbatim** | `[Source-derived]` | r7 §6.8; A2-D4 |
| C1-R-25 — exit classification | **closed in r7; retained verbatim** | `[Source-derived]` | r7 §6.9; A2-D6 |
| C1-R-26 — client handshake provenance | **closed in r7; retained verbatim** | `[Source-derived]` | r7 §6.11; A2-D3/D4/D6 |
| C1-R-27 — terminal set and `activeRequestID` | **closed in r7; retained verbatim** | `[Source-derived]` | r7 §6.10; A2-D5/D6 |
| C1-R-28 — catalog selection and C7 adapter | closed in r7; **ownership carrier added in §5** | `[Source-derived]` | r7 §6.12; §17.12.1 |
| C1-R-29 — incompatible barrier contradiction | **closed in r7; retained verbatim** | `[Source-derived]` | r7 §6.7; A2-D4 |
| C1-R-30 — matrix uniformity | **closed in r7; retained verbatim** | `[Calculated]` | r7 §4.6; 48 cells |
| C1-R-31 — `request_not_found` authority | **closed in r7; retained verbatim** | `[Source-derived]` | r7 §4.2/§4.5; §17.4.6 |
| C1-R-32 — enforceable update reset | **closed by §4 of this revision** | `[Source-derived]` | §4; A2-D4/D12; r9 §10.2 |
| C1-R-33 — `UpdateDecision` forgeable | **closed by §4 of this revision** | `[Source-derived]` | §4; A2-D4/D12 |
| C1-R-34 — exited asset-specific `Available` | **closed by §6 of this revision; no owner decision required** | `[Source-derived]` | A2-D5 (a), contract line 2506 |
| C1-R-35 — `catalogSelected` prose restriction | **closed by §5 of this revision** | `[Source-derived]` | §5; §17.12.1 |
| C1-R-36 — stale §6.6 baseline label | **closed by §7 of this revision** | `[Observed]` | §7 |
| C1-R-37 — unexercised budget class | **closed by §6 and C1-AC-33 in §8** | `[Source-derived]` | A2-D5 (b) |

Register totals: **37 findings**, every identifier carrying a disposition. `[Calculated]`

## 4. Update authority — §6.13 replaced (closes C1-R-32 and C1-R-33)

### 4.1 Defect being corrected

Revision 7 left `updateGatedReset` consuming a caller-constructible `UpdateDecision` whose fields
carried `Approved`, provider, engine, generation, epoch and operation. Because every field was
exported, any package-internal caller could construct it from observable state; the design's rechecks
read the same observable state. That is an assertion of authority, not an enforcement of it.
`[Source-derived]` from r7 §6.13 and the r7 review §4.

### 4.2 Normative design

```text
// Binding: derived by the SUPERVISOR from its own state. Never caller-supplied.
type UpdateBinding struct {
    ProviderID string
    Engine     string
    Generation uint64
    Epoch      uint64
    Operation  string
}

// Exported type, no exported fields, no constructor, no accessor, no Stringer/JSON method.
// A ticket is meaningful only if its nonce is present in the issuing authority's registry.
type UpdateTicket struct {
    nonce [32]byte
}

// C6/C7 supplies the implementation through SupervisorOptions at construction.
type UpdateAuthority interface {
    // Issue is called only by the C6/C7 update path, after it has completed and recorded
    // an actual helper/assets/config update. One issuance per recorded update.
    Issue(ctx context.Context, b UpdateBinding) (UpdateTicket, error)

    // Consume verifies that the ticket was issued for exactly this binding, that it has not
    // been consumed or revoked, and that it has not expired. It then marks it consumed.
    // Nil error IS the approval. There is no Approved field anywhere in C1.
    Consume(ctx context.Context, t UpdateTicket, b UpdateBinding) error
}
```

Supervisor-side contract, normative:

1. The supervisor holds `authority UpdateAuthority` installed once at construction. There is **no runtime setter** and **no exported ticket constructor**. `[Source-derived]`
2. `updateGatedReset` carries only `UpdateTicket`. The supervisor constructs `UpdateBinding` from its own fields under `Supervisor.mu`; it never accepts a caller-supplied binding, provider, engine, generation, epoch, or operation. `[Source-derived]`
3. Sequence: pre-check no live process and a terminal state; call `authority.Consume` **outside** `Supervisor.mu`; on error, reject with no mutation; on success, reacquire, recheck process absence, generation and epoch, then apply the reset. `[Source-derived]`
4. If the post-consume recheck fails, the ticket is **burned and no reset occurs**. Fail-safe direction is always "no reset". `[Inferred]`, stated as a requirement.
5. No authority installed → every reset is rejected. `[Source-derived]`
6. The supervisor stores no ticket, caches no approval, and exposes no query that reveals ticket validity. Tickets are never logged; `UpdateTicket` exposes no formatting surface. `[Inferred]`, stated as a requirement.
7. Authority obligations (`Issue`/`Consume` must implement): CSPRNG nonce of 32 bytes; a registry of issued tickets with their binding; one-use consumption; TTL expiry; revocation; rejection of unknown nonces; rejection of binding mismatch; atomicity under the authority's own lock so two concurrent `Consume` calls cannot both succeed; `Issue` gated on the authority's own recorded update evidence, not on the caller's assertion. `[Source-derived]` from r9 §10.2/§12 intent; the concrete store is C6/C7's.
8. `Issue` is idempotent-hostile: a second `Issue` for the same recorded update returns an error rather than a second ticket. `[Inferred]`

### 4.3 Why this is enforceable rather than descriptive

- A ticket is valid only if its 32-byte nonce exists in the authority's registry. Construction from observed state cannot produce a registered nonce. `[Source-derived]`
- The zero value `UpdateTicket{}` has an all-zero nonce that no correct authority issues; it is rejected by registry lookup, not by a field comparison the caller controls. `[Source-derived]`
- Approval is `Consume` returning nil. There is no boolean a caller can set. `[Source-derived]`
- Binding comes from the supervisor, so "issued for a different operation/generation/epoch/provider/engine" is detectable by the authority even though the caller cannot name those values. `[Source-derived]`

### 4.4 Stated limits — not overclaimed

- Go cannot prevent `unsafe`, `reflect` field-writes inside the same module, or memory corruption. The guarantee is at the Go type and API level, within the helper package boundary. `[Unestablished]` beyond that boundary; declared out of scope in §14.
- Installing an `UpdateAuthority` is itself an authority-bearing act performed by the update owner at construction. C1 gates *installation* (no runtime setter) but cannot audit the installed implementation; auditing it is C6/C7 acceptance work. `[Source-derived]`
- This design does not implement `Issue` evidence collection. C1 defines the boundary; C6/C7 must supply and test the real authority. `[Unestablished]`

### 4.5 Required tests (specification, not executed)

| ID | Setup | Required result |
|---|---|---|
| T33-1 | `var t UpdateTicket` (zero value) | `Consume` rejects; no reset |
| T33-2 | `reflect.New(UpdateTicket).Elem().Interface()` | rejected; no reset |
| T33-3 | Ticket issued for binding A, consumed with binding B | rejected; ticket **not** consumed |
| T33-4 | Valid ticket consumed once, then replayed | second rejected |
| T33-5 | Two concurrent `Consume` calls, same ticket | exactly one succeeds |
| T33-6 | Ticket issued, then generation advances / process starts | post-consume recheck fails; ticket burned; no reset |
| T33-7 | No authority installed | every reset rejected |
| T33-8 | Static scan of C1 surface | no exported ticket constructor, no exported ticket field, no runtime authority setter, no `Approved` field |
| T33-9 | Ticket issued for wrong operation/provider/engine | rejected |
| T33-10 | Expired or revoked ticket | rejected |
| T33-11 | Ordinary package-internal caller with no issued ticket | cannot produce an accepted value |

### 4.6 Ownership table

| Decision | Owner | Mechanism |
|---|---|---|
| Whether an update actually happened | C6/C7 authority | Authority's own recorded evidence; `Issue` refuses without it |
| Issuing a ticket | C6/C7 authority | `Issue`, one per recorded update |
| Whether a specific ticket authorises this reset | C6/C7 authority | `Consume` registry lookup + one-use + binding match |
| Constructing the binding | C1 supervisor | Supervisor state only |
| Applying the reset | C1 supervisor | Only after `Consume` returns nil and rechecks pass |
| Clearing `incompatible`/`unavailable` | C1 supervisor | Only via this path or `observeExit` |

## 5. Selection ownership — §6.12 amended (closes C1-R-35)

The r7 restriction ("the caller is the direct Catalog path") is replaced by a carrier with the same
shape as `handshakeAttempt`:

```text
catalogSelection struct {
    run        uint64   // minted by beginCatalogSelection; zero never matches
    generation uint64
    providerID string
    engine     string
}

// private; called only by the direct Catalog path while its operation is in flight under workMu.
func (s *Supervisor) beginCatalogSelection(gen uint64) (catalogSelection, error)
func (s *Supervisor) endCatalogSelection(run uint64)   // closes the run
```

Rules, normative:

1. `catalogSelected` accepts only a `catalogSelection` whose `run` equals the open run, whose
   `generation` equals the current generation, and whose provider/engine belongs to the current
   validated handshake. `[Source-derived]`
2. The run counter starts at 1, so a zero value never validates. `[Calculated]`
3. The adapter path (`Adapter.ProviderCatalog`) never calls `beginCatalogSelection`, so it cannot
   obtain a run that validates. `[Source-derived]`
4. The run is closed when the Catalog operation returns; a late or replayed event is rejected. `[Source-derived]`
5. The event still changes identity only — no lifecycle, capability, failure, or recovery effect. `[Source-derived]`
6. Tests: adapter attempts selection → rejected; zero/guessed run → rejected; stale run after
   completion → rejected; positive direct-Catalog run → accepted once. `[Calculated]`

Residual limit, stated: a same-package caller that calls `beginCatalogSelection` itself can open a
run. This is narrower than r7's prose restriction because the caller must be executing a Catalog
operation under `workMu`, and the result is still validated against the handshake. Selection is
identity-only and cannot clear a terminal state or start a process, so the exposure is bounded and
strictly smaller than the reset boundary of §4. `[Inferred]`

## 6. Availability classes — §6.8 amended (closes C1-R-34 and C1-R-37)

### 6.1 Contract basis, quoted in full

A2-D5 (contract line 2506), verbatim on the decisive clause:

> `Available` is true for `unloaded`, `loading`, `ready`, `busy`, and `degraded`. It is also true for
> `unavailable` when either (a) the failure is asset-specific and the helper, protocol, and provider
> configuration remain valid, or (b) the state was entered because the `crashed`/`backoff` restart
> budget was exhausted while the provider remains configured and compatible; in case (b),
> `ProcessAlive=false`, `Loaded=false`, `Ready=false`, and no automatic restart is permitted.

> Thus `Available` and `State` are independent, and `State=unavailable` may coexist with
> `Available=true` for an asset-specific or budget-exhausted failure.

### 6.2 R-34 disposition

The liveness condition `ProcessAlive=false` is stated **only for case (b)**. Case (a) is conditioned
on "the helper, protocol, and provider configuration remain valid", which is a compatibility and
configuration test, **not** a liveness test; and the closing sentence binds `Available=true` to
"asset-specific or budget-exhausted" without distinguishing live from exited. `[Source-derived]`

Therefore `Available=true` for an exited asset-specific `unavailable` is **authorised by A2-D5 (a)**
and requires no owner decision. §17.7.4's additive row that pairs `ProcessState=exited` with
`Available=true` and budget exhaustion is an example of one way to reach that combination, not an
exhaustive enumeration. `[Inferred]` from §17.19.3's additive-row structure.

`ProcessAlive` remains independent and comes from `processState` alone (r7 §6.8, A2-D4 `[Source-derived]`).
`ServiceCapabilitiesKnown` remains false after exit (A2-D6 `[Source-derived]`); availability and
snapshot currency are separate facts.

### 6.3 Availability classes, completed

| Class | Set when | `Available` | Owner of the classification |
|---|---|---|---|
| `assetSpecificConfigured` | typed asset/model/runtime load failure, configuration valid | true, live or exited | C1 |
| `budgetExhaustedConfigured` | `crashed`/`backoff` restart budget exhausted, provider configured and compatible | true; `ProcessAlive=false` per A2-D5 (b) | **C5** |
| `helperOrConfigFailure` | missing helper, `proc.Start` failure, invalid provider configuration | false | C1 |

`providerAvailable` is unchanged in shape; the third class is now exercised by C1-AC-33 (§8) and a
test row (§10). C1 supplies the predicate and the C5-supplied classification input; C1 does not own
budget accounting. `[Source-derived]`

## 7. §6.6 baseline label — corrected (closes C1-R-36)

The r7 §6.6 wording "at the revision-6 source baseline" and the header cell `Revision-6 owner/event`
are replaced with:

> Baseline: current `internal/audio/tts/helper` source at HEAD `9c7d2a1` plus the working tree as of
> 2026-11-08. Line references are current-source evidence, not implementation claims. Historical
> finding identifiers (C1-R-20, C1-R-27) are preserved in the register in §3 of this document; the
> disposition column states the **revision-8** disposition for each row.

The 18 inventory rows and their line references are unchanged `[Observed]`; only the baseline label
and the disposition column name change. No historical finding is restated as having existed in r6.

## 8. Acceptance criteria (33)

Criteria **1–26 and 28–30 are unchanged** from r7 §10 and remain normative there (`[Observed]`,
r7 hash above). Reproduced in full because they change: **27, 31, 32**; added: **33**.

| ID | Criterion (revision 8) |
|---|---|
| **C1-AC-27** | `ProcessAlive` derives solely from `processState` and remains true for a live asset-specific `unavailable`; provider state cannot make it false. `providerAvailable` is true for asset-specific `unavailable` **whether the process is live or exited, provided configuration remains valid** (A2-D5 (a)), true for `budgetExhaustedConfigured` with `ProcessAlive=false` (A2-D5 (b)), and false for `helperOrConfigFailure`. Tests cover: live retained snapshot, exited invalidated snapshot, incompatible-live, budget-exhausted, helper/config failure, and a negative control that a cleared process cannot leave `Available` stale. The A2-D5 (a) clause is cited in the trace, and no owner decision is claimed. |
| **C1-AC-31** | Direct `Client.Catalog` in a normal live state may emit only the validated `catalogSelected` event, carrying a `catalogSelection` minted by `beginCatalogSelection` for the in-flight direct Catalog operation, matching the current generation and a provider/engine in the current handshake; the run is closed when the operation returns. Live `unavailable` Catalog is read-only and emits no selection. No `SetSelected` setter remains. `catalogSelected` cannot change lifecycle, capabilities, failure, or recovery. C7 `Adapter.ProviderCatalog` is the only adapter start entry point, may start only from ordinary `absent`, then calls no-start `Client.CatalogCurrent`; it cannot mint a selection run. Negative tests: adapter-emitted selection, zero run, stale run after completion, handshake-external provider; positive: one accepted direct selection. |
| **C1-AC-32** | The C6/C7 reset is authorised only by an authority-verified `UpdateTicket`. There is no exported ticket constructor, no exported ticket field, no `Approved` value, no caller-supplied binding, and no runtime authority setter; with no authority installed every reset is rejected. The supervisor builds `UpdateBinding` from its own state, calls `Consume` outside `Supervisor.mu`, applies only on nil error, and burns the ticket if the post-consume recheck fails. Tests reject: zero-value ticket, reflect-constructed ticket, wrong binding, replay after consumption, concurrent double-consume (exactly one succeeds), expired/revoked ticket, stale generation or epoch, live process, and an ordinary caller with no issued ticket; a valid ticket succeeds exactly once. C6/C7 must supply the real authority and executed evidence. |
| **C1-AC-33** | The `budgetExhaustedConfigured` availability class is exercised: a `crashed`/`backoff` budget-exhausted `unavailable` with a configured, compatible provider reports `Available=true`, `ProcessAlive=false`, `Loaded=false`, `Ready=false`, `ServiceCapabilitiesKnown=false`, and permits no automatic restart. A test drives the class through the C5-supplied classification seam and asserts both the projection and the absence of a restart. C1 owns the predicate; C5 owns budget accounting. |

### 8.1 Criterion locus map (all 33)

Reconciliation: 33 criteria = 26 carried (AC-1 … AC-26) + 3 carried (AC-28 … AC-30) + 4 in this
revision (AC-27, AC-31, AC-32, AC-33). Verification rows: 33 = 29 at r7 §10.1 + 4 at §9 above. `[Calculated]`

| Criterion(s) | Normative text | Verification row | Status |
|---|---|---|---|
| AC-1 … AC-26 (26 criteria) | r7 §10, unchanged | r7 §10.1 rows 1–26, unchanged | carried `[Observed]` |
| AC-27 | **this document §8** (revised: exited asset-specific available, A2-D5 (a)) | **this document §9** | revised |
| AC-28, AC-29, AC-30 (3 criteria) | r7 §10, unchanged | r7 §10.1 rows 28–30, unchanged | carried `[Observed]` |
| AC-31 | **this document §8** (revised: selection run carrier) | **this document §9** | revised |
| AC-32 | **this document §8** (revised: authority-verified ticket) | **this document §9** | revised |
| AC-33 | **this document §8** (new: budget-exhausted availability) | **this document §9** | new |

Every criterion therefore has exactly one normative text and one verification method, with none
unverified and none duplicated. The r7 §10.1 table remains normative for the 29 carried rows. `[Calculated]`

## 9. Verification and traceability additions

| Criterion | Design locus | Requirement trace | Verification method |
|---|---|---|---|
| C1-AC-27 | §6, r7 §6.8 | A2-D4, A2-D5 (a) and (b), A2-D6 | Predicate matrix incl. exited asset-specific and budget-exhausted; process-state-only liveness negative control |
| C1-AC-31 | §5, r7 §6.12 | §17.12.1; A2-D5, A2-D12 | Run-carrier matrix: adapter, zero run, stale run, foreign provider, positive direct |
| C1-AC-32 | §4 | A2-D4, A2-D12; r9 §10.2/§12 | Ticket matrix T33-1…T33-11; static scan for constructors/fields/setters/`Approved` |
| C1-AC-33 | §6.3 | A2-D5 (b) | Budget-class projection test through the C5 seam; no-restart assertion |

## 10. Test matrix additions

| Case | Setup | Required result |
|---|---|---|
| Ticket: zero value | `var t UpdateTicket` submitted to reset | Rejected; no mutation |
| Ticket: reflected value | `reflect.New(UpdateTicket).Elem().Interface()` | Rejected; no mutation |
| Ticket: binding mismatch | Valid ticket, different operation/provider/engine/generation/epoch | Rejected; ticket not consumed |
| Ticket: replay | Valid ticket consumed, resubmitted | Rejected |
| Ticket: concurrent consume | Two goroutines, same ticket | Exactly one reset; no double apply |
| Ticket: stale after recheck | Valid ticket, process starts or generation advances before apply | Ticket burned; no reset |
| Ticket: no authority | Supervisor constructed without an authority | Every reset rejected |
| Ticket: static surface | Source scan of C1 | No exported constructor/field, no runtime setter, no `Approved` |
| Selection: adapter attempt | Adapter path emits `catalogSelected` | Rejected; identity unchanged |
| Selection: stale run | Run closed before the event is applied | Rejected |
| Budget-exhausted availability | C5 seam reports budget exhaustion, provider configured | `Available=true`, `ProcessAlive=false`, no restart |

## 11. Requirements-to-design trace for changed items

| Requirement | Change | Governing source |
|---|---|---|
| A2-D4 (`ProcessAlive` from process state only) | unchanged predicate; re-traced in AC-27 | contract line 2505 |
| A2-D5 (a) availability | R-34 closed by quoting the clause; no liveness condition | contract line 2506 |
| A2-D5 (b) availability | R-37 closed; AC-33 and budget class | contract line 2506 |
| A2-D6 (snapshot currency) | restated: exit clears `ServiceCapabilitiesKnown`, not `Available` | contract line 2507 |
| A2-D12 (recovery only by explicit load) | unchanged; reset path cannot recover | contract line 2513 |
| §17.12.1 (adapter start, no load) | unchanged; selection carrier added | contract §17.12.1 |
| §17.4.6 (`request_not_found`) | unchanged from r7 | contract §17.4.6 |
| r9 §10.2/§12 (update authority) | token replaced by authority-verified ticket | r9 requirements |
| D10 A10a-2/A10b-2 | unchanged | owner decisions |

## 12. Incorporation of unchanged sections

This revision replaces §6.13 (§4 here), amends §6.12 (§5), §6.8 (§6), §6.6 labels (§7), and criteria
27, 31, 32 with a new 33 (§8), and adds the rows in §9–§11. **All other sections of revision 7 remain
normative, unmodified, at SHA-256 `0867acb7216a0e8ee2802b6f08275d997e4138f089ae51acac373c10eea4af54`**:
§1, §2, §3 (3.1–3.5), §4 (4.1–4.7), §5, §6 (6.1–6.7, 6.9–6.11), §7, §8, §9, §10 (except 27/31/32),
§10.1 (except those rows), §11, §12, §13, §14.

Read together with this document, the normative corpus is closed and unambiguous. This is a
deliberate choice: rewriting 1,069 lines of reviewed text to append 37 findings would risk silent
regressions in areas the review just validated, which §6 of the task forbids. `[Inferred]` A
consolidated single-file r9 can be produced if the reviewer requires it.

## 13. Explicitly unresolved (not closed by this revision)

- **C1-AC-15** — 13 resource rows need independent written determinations; 8 I/O categories need
  named/dated reviewed N/A entries. `[Unestablished]`
- **`MaxPending=32`** — no accepted authorised source or recorded R8-AC-4 disposition. `[Unestablished]`
- HW-1.6 independent review and performance-preservation evidence. `[Unestablished]`
- Every criterion remains executed-evidence-pending; no test was run. `[Observed]`

## 14. Authorization boundary

**C1 remains UNAUTHORIZED.** This document authorises nothing: no implementation, no C2–C8, no test
execution, no build, no runtime, no deployment, no owner decision, no contract amendment, no change to
the frozen plan.

Out of scope by declaration: `unsafe`/`reflect` memory manipulation inside the module, process-level
tampering, and the correctness of any C6/C7 authority implementation. The §4 boundary is a Go
type-and-API guarantee within the helper package, and is claimed as nothing more.

## 15. Self-validation of this revision

| Check | Result | Class |
|---|---|---|
| Findings R-1…R-37 each have a disposition | yes, 37 rows in §3 | `[Calculated]` |
| Criteria 1…33 each have a verification row | yes: 29 carried at r7 §10.1 + 4 at §9 (§8.1 map) | `[Calculated]` |
| Changed criteria reproduced in full | 27, 31, 32, 33 | `[Observed]` |
| No duplicate section numbering | §1–§15 unique | `[Calculated]` |
| Removed r7 surfaces still removed | no `SetState`/`SetSelected`/`MarkHandshake` reintroduced; `UpdateDecision` and `Approved` appear only as historical defect references, never as design surface | `[Observed]` |
| Every ownership restriction has a mechanism | reset → §4 registry ticket; selection → §5 run carrier | `[Source-derived]` |
| Obsolete definitions removed | `UpdateGrant`/`UpdateDecision` superseded by `UpdateTicket`; budget class now exercised | `[Superseded]` |
| Stale revision labels remaining | none introduced; §6.6 corrected in §7 | `[Observed]` |
| Implementation performed | none | `[Observed]` |
