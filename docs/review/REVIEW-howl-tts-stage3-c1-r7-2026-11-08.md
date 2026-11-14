# Independent Review — Howl TTS Stage 3 C1 Revision 7

**Date:** 2026-11-08
**Review type:** Independent, read-only design review
**Reviewer:** Buffy
**Artifact set reviewed:** C1 design r7, readiness r7, integrity r7
**Verdict:** **FAIL — NARROW. R-24 through R-31 DESIGN-CLOSED; C1-R-32 PARTIALLY CLOSED. C1 REMAINS UNAUTHORIZED.**

This review is the independent r7 review requested by the r7 readiness report §9. It is not an
acceptance act and does not authorise implementation.

## 1. Independent identity verification

Computed directly from the workspace, not copied from the integrity report.

| Artifact | SHA-256 | Lines / bytes | Matches recorded claim |
|---|---|---|---|
| `DESIGN-howl-tts-stage3-c1-revision-2026-11-08-r7.md` | `0867acb7216a0e8ee2802b6f08275d997e4138f089ae51acac373c10eea4af54` | 1069 / 161,934 | Yes |
| `C1-READINESS-howl-tts-stage3-2026-11-08-r7.md` | `8a42cf730a2ce35fad338c05746a36495fbf7ee52bb544bf5ff1e0ed4d29dfb2` | 127 / 14,684 | Yes |
| `DESIGN-...-r6.md` (preserved) | `dae4b371e901c11144a80dd1e3b99239ed25d1341088546ce058447e232a6fb4` | 853 / 127,414 | Yes |
| `C1-READINESS-...-r6.md` (preserved) | `d5d21bbd9055b51cba8aac58cd0ea5b46ed83ea3d634e29700310b6afb8823e7` | 141 / 13,300 | Yes |
| `REVIEW-...-c1-r6-2026-11-08.md` (preserved) | `9176f21512698a3d15ea073b84d39ae69c88fdcb8ef92ca37f0c99288889f0c9` | 252 / 28,668 | Yes |
| `DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md` | `802ae8114993b7889140f5f1c216eca042f74e84f7dbcb1665f05434e355b151` | 711 / 69,168 | Yes |
| `OWNER-DECISIONS-howl-tts-stage3-2026-11-07.md` | `677f884926de7e58c3dc7ab33f01340f1fbb5966c1b5c089b1aafd057bcf0fba` | 612 / 49,191 | Yes |
| `FINALIZATION-...-d10-...-2026-11-08.md` | `84fa78dd214f16f40ab125b544727fce2239e66bb1bb54dffc3c8777df9beb02` | 22 / 1,887 | Yes |
| `REQUIREMENTS-...-r9-2026-11-08.md` | `0f71cee5e0e5d25b58d1b3956ff549df8cc625ad7f9e8c82189f2f2ef9366c55` | 769 / 112,270 | Yes |
| `DESIGN-howl-tts-helper-architecture-2026-11-06.md` (contract) | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f` | 2586 / 182,904 | Yes |

No prior artifact was modified by this review or by revision 7.

## 2. Structural verification

| Check | Result |
|---|---|
| Distinct `C1-R` identifiers | 32 (max `C1-R-32`) |
| Distinct `C1-AC` identifiers | 32 (max `C1-AC-32`) |
| §10.1 verification rows | 32, one per criterion |
| §4.6 cell count | **48 data rows**, derived as 60 table lines − 6 separators − 6 headers; matches the claimed six classes × eight operations |
| §6.6 writer inventory | 18 rows, matching the stated count |
| Duplicate `### N.M` section numbers | none |
| Removed surfaces (`SetState`, `SetSelected`, `MarkHandshake`) | no definitions remain; `helloHandshake` appears only in removal statements |
| `C1-AC-15` | explicitly **OPEN**, not absorbed |
| `MaxPending=32` | still recorded as unresolved, no invented authority |

The r7 numerals are internally consistent. This is a real improvement over r6, whose matrix count
was arithmetically wrong.

## 3. Reconciliation of C1-R-24 through C1-R-32

"Design-closed" means the design text now addresses the cited defect. It does not mean
acceptance-closed; every criterion remains executed-evidence-pending.

| Finding | Disposition | Basis |
|---|---|---|
| C1-R-24 — derived liveness/availability | **design-closed** | §6.8 makes `processAlive(processState)` sole-input and forbids reading provider state, matching A2-D4's "derived solely from `TTSProcessState`" exactly. `providerAvailable(state, availabilityClass)` matches A2-D5's enumeration. C1-AC-27 carries the negative control. |
| C1-R-25 — exit classification | **design-closed** | §6.9 defines four captured intents including `unavailableFailure`, first-intent-wins, late-exit, replacement overlap, and a forbidden-action table; §6.7 and §3.4 align. C1-AC-28 covers the state × intent matrix. |
| C1-R-26 — client handshake provenance | **design-closed** | §6.11 removes `helloHandshake`, makes `attempt.publish` the only publication and `installHandshake` the only client writer, and binds every consumer to `currentHandshake()` with generation/process/transport/snapshot validation. C1-AC-29 adds the forbidden-intermediate and stale-copy cases. |
| C1-R-27 — terminal event set and `activeRequestID` | **design-closed** | §6.10 enumerates ten outcomes mapped line-by-line onto the current `client.go` call sites, moves terminal writes into one parameterised event, and gives `activeRequestID` an explicit writer set. C1-AC-30 requires a table-driven test per outcome. |
| C1-R-28 — catalog selection and C7 adapter | **design-closed** | §6.12 introduces `catalogSelected` (identity-only, normal-live direct Catalog) and `Adapter.ProviderCatalog` as the sole adapter start path calling no-start `CatalogCurrent`; §4.7 makes the authority table normative. `Client.ProviderCatalog` delegation is removed/renamed. |
| C1-R-29 — incompatible barrier contradiction | **design-closed** | §6.7's R18 termination rule restricts `beginTermination` to `processState=stopping` while preserving aggregate `incompatible` and `LastError`; C1-AC-24 now separates "ordinary writers mutate nothing" from the authorised process-state transition and asserts the barrier observation. |
| C1-R-30 — matrix uniformity | **design-closed** | §4.6 now carries direct `Catalog` and `Catalog` (`ProviderCatalog`) rows in all six classes; 48 cells verified arithmetically. C1-AC-26 and the §11 row state 48. |
| C1-R-31 — `request_not_found` authority | **design-closed** | §4.2/§4.5 send exactly one current-transport `cancel` frame and consume the helper's wire `request_not_found`, matching contract §17.4.6 and the §17.9 registry entry (wire origin → `ErrCodeSynthesisFailed`). The adapter no longer synthesises a wire-origin code. |
| C1-R-32 — enforceable update reset | **partially closed — see C1-R-33** | §6.13 correctly moves grant creation, storage, revocation and one-use consumption into an `UpdateAuthority` supplied at construction, with no runtime setter or exported constructor. But the value `updateGatedReset` actually consumes is still forgeable. |

## 4. New findings

### C1-R-33 — `UpdateDecision` is forgeable, so the sealed reset boundary is not sealed

**Severity:** High
**Artifacts:** design §6.13 (types and `updateGatedReset` contract), §6.5 (`updateGatedReset` may change `incompatible`), C1-AC-25, C1-AC-32.

**Observed text.** §6.13 declares:

```text
type UpdateDecision struct {
    GrantID [32]byte
    Approved bool
    ProviderID string
    Engine string
    Generation uint64
    Epoch uint64
    Operation string
}
```

every field exported, and states that `updateGatedReset` "accepts only a validated `UpdateDecision`".
The same section states the authority "owns `UpdateGrant` creation, storage, revocation, and one-use
consumption" and that "the supervisor stores no bearer token and cannot create a grant".

**Inference.** The grant is now opaque, but the artefact the supervisor trusts is a plain exported
struct. Any package-internal caller — including C7 adapter code, which is the intended producer —
can build `UpdateDecision{Approved: true, Generation: g, Epoch: e, ...}` from readable state and
submit it. The stated rechecks (no live process, current generation, matching epoch) are all
satisfiable with observed values, and nothing in §6.13 describes the supervisor consulting an
authority-side consumed-grant registry; `GrantID` is data, not proof. So the property C1-R-32 asked
for — "an ownership check that the caller cannot fabricate" — is asserted, not enforced. This is the
same class of gap r6 identified, relocated one layer outward.

**Status:** **OPEN — the only high-severity finding in revision 7.** R-32 is not yet closed.

**Required correction.** Make the consumed value unforgeable rather than descriptive: have `Consume`
return an unexported type with no exported fields (or a value carrying an authority-private nonce the
authority can re-verify), and keep `UpdateDecision` internal to the authority. Add a negative test
that constructs the public shape directly and is rejected, alongside the existing
fabricated/reused/foreign/stale/revoked cases.

### C1-R-34 — `Available=true` for an exited asset-specific `unavailable` is an uncited extension

**Severity:** Medium
**Artifacts:** design §6.8 (`providerAvailable`), C1-AC-27, §11 "Derived predicates"; contract A2-D5, §17.7.4 additive rows.

**Observed text.** §6.8 returns `true` for `unavailable` when the class is `assetSpecificConfigured`,
and explicitly states "an asset-specific failure whose process has since exited remains configured
and therefore available while the aggregate state remains `unavailable`".

**Contract comparison.** A2-D5's asset clause keeps `Available` set while "the helper, protocol, and
provider configuration remain valid" — it does not require a live process, which supports the r7
reading. However §17.7.4's additive rows pair `ProcessState=exited`, `State=unavailable`,
`Available=true` with *restart-budget exhaustion*, and describe the asset-specific failure row with a
live helper (`ProcessState=alive`). The exited asset-specific combination is therefore reached by
inference from A2-D5 rather than by an explicit contract sentence.

**Status:** **OPEN — interpretation, not contradiction.** It may be correct; it is not yet authorised.

**Required correction.** Cite A2-D5's asset clause explicitly as the authority for the exited case, or
record an owner decision. Either way, trace it from C1-AC-27 rather than leaving it in §6.8 prose.

### C1-R-35 — `catalogSelected` caller restriction has no mechanism

**Severity:** Medium
**Artifacts:** design §6.12, §6.5 event set, §6.7 matrix row, C1-AC-31.

**Observed text.** §6.12 accepts `catalogSelected` only when "the caller is the direct Catalog path,
not the C7 adapter or a supervisor cancellation path".

**Inference.** `catalogSelected` is a value in the closed event set, and `apply(event)` is
package-internal. Unlike the handshake, which r7 correctly made an owned object
(`handshakeAttempt` with private state), selection has no carrier that distinguishes its producer.
The design removes exported setters and then re-admits a writer whose restriction lives only in
prose. The exposure is bounded — the event changes identity only, validated against the current
handshake's providers — but it is the same "asserted, not enforced" shape as C1-R-33.

**Status:** **OPEN — medium-severity enforcement gap.**

**Required correction.** Either give selection an owned carrier analogous to `handshakeAttempt`
(returned by the validated direct-catalog path), or state the package-internal trust boundary
explicitly and add a test proving the adapter path cannot emit `catalogSelected`.

### C1-R-36 — §6.6 still labels itself revision 6

**Severity:** Low
**Artifacts:** design lines 490 and 492.

Line 490 reads "at the revision-6 source baseline"; the table header at line 492 reads
`Revision-6 owner/event`, while the rows now describe r7 dispositions (for example "Removed; closed
`apply(event)`"). C1-AC-25 traces directly to this table, so a mislabeled baseline is a traceability
defect rather than cosmetics. Relabel to revision 7 and state that the source baseline is unchanged.

### C1-R-37 — The `budgetExhaustedConfigured` availability class is defined but never asserted

**Severity:** Low
**Artifacts:** design §6.8 (three classes), C1-AC-27, §11 "Derived predicates".

§6.8 defines `budgetExhaustedConfigured` alongside `assetSpecificConfigured` and
`helperOrConfigFailure`, but C1-AC-27 and the corresponding §11 row exercise only the asset-specific
and helper/config cases. Mark the class C5-owned with a binding criterion, or add it to AC-27.

## 5. Acceptance evidence still absent

- **C1-AC-15 remains OPEN**: 13 resource rows without independent written determinations, 8 I/O
  categories without named/dated reviewed N/A entries.
- `MaxPending=32` remains unresolved with no accepted authorised source.
- HW-1.6 independent review and performance-preservation evidence remain pending.
- No test, build, runtime, benchmark, or conformance result exists for any criterion.
- No finding is acceptance-closed; r7's own §2.5 correctly says so.

## 6. Verdict and exit criteria

**FAIL — narrow.** Revision 7 design-closes eight of the nine high/medium findings that failed r6,
and does so substantively rather than cosmetically: the derived predicates match A2-D4 literally, the
terminal set is traced line-by-line to current source, and the 48-cell matrix is arithmetically true.

One high-severity gap remains: the reset boundary consumes a forgeable exported decision, which is
precisely the property C1-R-32 was raised to fix. Three medium/low findings are traceability or
enforcement wording issues.

To clear this review, revision 8 needs: an unforgeable consumed update value (C1-R-33), an authorised
basis for the exited asset-specific `Available` case (C1-R-34), an enforcement mechanism or explicit
trust boundary for selection (C1-R-35), the §6.6 relabel (C1-R-36), and coverage for the budget class
(C1-R-37). None of these requires redesign.

## 7. Review boundary

Read-only. No source, test, build, runtime, benchmark, deployment, owner decision, or governing input
was modified, and no prior revision was touched. This review does not authorise C1 implementation,
C2–C8, or any claim of acceptance. C1 remains **unauthorized**.

**Limitation of this review:** the contract's A2-D5 text exceeded the reader's per-line clip limit and
was partly truncated; the C1-R-34 conclusion rests on the visible clause plus the §17.7.4 additive
rows. A reviewer with the full clause should confirm or withdraw that finding.
