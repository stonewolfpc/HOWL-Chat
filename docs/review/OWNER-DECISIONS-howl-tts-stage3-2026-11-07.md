# OWNER DECISIONS — Howl TTS Stage 3 corrective design

**Date:** 2026-11-07; D10 decisions recorded 2026-11-08  
**Mode:** DOCUMENTATION ONLY — DECISION INSTRUMENT  
**Status:** OPEN — TWO DECISIONS RECORDED; SIX PENDING  
**Owner:** Stone Wolf Systems  
**D10 decisions recorded:** 2026-11-08  
**Governing contract:** `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md`  
**Frozen design inputs:** `docs/review/DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md` and `docs/review/CORRECTIVE-howl-tts-stage3-corrective-design-2026-11-07.md`

This document records the questions, evidence, constraints, and gates for eight owner decisions. It records answers for D10(a) and D10(b) only; it does not answer the remaining six, recommend an answer for them, modify the frozen design, authorize a stage, or authorize implementation. Blank fields are intentional. A blank field is not consent, a default, or an inferred selection.

## 1. Frozen source baseline and verification

The following hashes were recomputed before drafting and matched the supplied frozen baseline exactly.

| Artifact | SHA-256 | Size |
|---|---|---|
| Corrective proposal | `802ae8114993b7889140f5f1c216eca042f74e84f7dbcb1665f05434e355b151` | 711 lines / 69,168 bytes |
| Correction record | `8fe9356cc304f3135657d19a7f2ea6badd9cf2233585a03bd0bd5c4f6355e076` | 172 lines / 16,190 bytes |
| Governing contract | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f` | 2,586 lines / 182,904 bytes |

The review inputs were also hash-checked:

| Review input | SHA-256 |
|---|---|
| Round-1 audit | `71a17369e2389fae26e86783fc4df802ab91bad3daef512c0e7f743067e94ea2` |
| Round-2 audit | `a7384da9f3f38661fc28b3de48badb82245cfd41efb8e49964f17bcfaaf5a7e2` |
| Round-3 audit | `6681e11b7eb770d59f3c79294206f05c8d3a3d9d4cfbbd038ef103fc8c48fd15` |
| Focused D10-scope audit | `30ff59bfe89867824c0abed40c12e71a69e200b5a9c5beae067cf49a87f57e43` |

**Authority.** The governing contract controls. The proposal and correction record are frozen design evidence. The audits are review evidence, not contract text and not owner selections. This record does not amend any of them.

**Duplicate check.** `docs/review` contains no equivalent Stage 3 owner-decision record. The earlier `DECISION-howl-tts-capability-policy-c1-2026-11-06.md` is a separate C-1 capability-policy decision report; it does not consolidate these eight Stage 3 decisions.

**Non-authorization.** No stage is authorized by this document. C1 and C2 remain unauthorized. C8 still requires separate runtime-execution authorization.

## 2. How to read the entries

Every entry uses these fields:

- **Exact owner question** — the question the owner must answer.
- **Status** — `PENDING — OWNER DECISION REQUIRED` before selection; `DECIDED — 2026-11-08` after the owner records one.
- **Why this is necessary** — the design or contract consequence that makes the answer necessary.
- **Governing contract** — contract sections and exact line anchors.
- **Proposal evidence** — frozen proposal sections and exact line anchors.
- **Constraints** — rules the owner cannot override with this decision.
- **Credible options** — only options the source establishes. Where the source establishes no credible option set, the entry says `UNESTABLISHED`.
- **Consequences and gates** — what follows from each option, and which stage acceptance or reachability remains gated.
- **Owner fields** — selection, rationale, and date, left blank for the owner.

Evidence labels used below:

- **`[CONTRACT]`** — a rule stated in the governing contract.
- **`[PROPOSAL]`** — a design constraint or unresolved question stated in the frozen proposal or correction record.
- **`[AUDIT]`** — a review observation; not authority and not a decision.
- **`[UNRESOLVED]`** — the source identifies the question but does not answer it.
- **`[UNESTABLISHED]`** — the source does not establish a credible option set.
- **Recommendation:** none. This record contains no recommendation.
- **Assumption:** none. No unstated behavior should be inferred.

A scheduling deferral is not a behavioral result. This is especially important for D10(a) and D10(b): deferral may schedule when the owner answers, but it cannot satisfy C1 acceptance. The two D10 branches are independent and must not be collapsed into one automatic answer.

## 3. Decision summary

| ID | Decision | Status | Primary gate |
|---|---|---|---|
| D10(a) | Unknown-snapshot direct `unload` result | `DECIDED — A10a-2 — 2026-11-08` | C1 acceptance |
| D10(b) | Unknown-snapshot direct `Cancel` result outside the supervisor path | `DECIDED — A10b-2 — 2026-11-08` | C1 acceptance |
| D1a | Permission for a narrow internal loader seam | `PENDING — OWNER DECISION REQUIRED` | C3; C7 upstream |
| D9 | Whether status/catalog/health observation resets the 120 s idle deadline | `PENDING — OWNER DECISION REQUIRED` | C5 acceptance |
| D1b | Explicit load action and adapter lifetime | `PENDING — OWNER DECISION REQUIRED` | C7 production reachability |
| D2 | Helper-backed engine activation and O-4 relationship | `PENDING — OWNER DECISION REQUIRED` | C7 production activation |
| D5 | Renderer exposure and O-1 picker behavior | `PENDING — OWNER DECISION REQUIRED` | Production UI integration |
| D3 | Third-process admission rejection result | `PENDING — OWNER DECISION REQUIRED` | C7 cap-response acceptance |

No row is a partial approval. No row authorizes its stage.

---

## 4. D10(a) — Unknown-snapshot direct `unload`

**Exact owner question.** When the authoritative service-capability snapshot is unknown, what concrete caller-facing result does a direct `unload` return, while preserving the structural no-frame/no-spawn boundary?

**Status:** `DECIDED — A10a-2 — 2026-11-08`

**Why this was necessary.** The contract states the known-absent direct-`unload` rule, but not the unknown-snapshot rule. Before this selection, C1 could implement the known-present and known-absent branches and could test the unknown branch only at the structural no-frame/no-spawn boundary. C1 could not be accepted until this branch had a concrete caller-facing result. A deferral is scheduling, not a result; the owner has now supplied the result.

**Governing contract.**

- `[CONTRACT]` §17.19.2 A2-D7 (line 2508): if `howl.tts.unload` is absent, a direct `unload` returns `unsupported_capability` and writes no frame.
- `[CONTRACT]` §17.19.2 A2-D1 (line 2502): an absent optional service capability returns `unsupported_capability`; no operation frame is written.
- `[CONTRACT]` §17.9 (lines 1876–1920; `unsupported_capability` row line 1918): `unsupported_capability` means a requested optional capability is unavailable.
- `[CONTRACT]` §17.7.2 (line 1773): lazy start from `absent` is limited to an explicit TTS/catalog/status request; `unload` is not a start trigger.
- `[CONTRACT]` §17.19.2 A2-D9 (line 2510): the supervisor cancellation path is a separate settled path and is not D10(a).

**Proposal evidence.**

- `docs/review/DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md` §4.4 (lines 121–123): known-absent direct operations return `unsupported_capability` with no frame; unknown snapshot permits no frame and no spawn; D10(a) is the unresolved direct-`unload` result.
- Same proposal, C1 (lines 137, 156, 177): C1 is not accepted until D10(a) has a concrete caller-facing disposition; the unknown test asserts only no frame/no spawn.
- Same proposal, §10 (line 684): D10(a) is separately answerable and no result is frozen.
- `[AUDIT]` Focused D10-scope audit §4.3 (lines 119–149): the adopted text supplies known-absent and supervisor-path rules, but not the direct unknown-snapshot result.

**Constraints that cannot be changed by this decision.**

- No operation frame may be sent for the unknown-snapshot direct `unload`.
- No helper may be started merely to discover the capability.
- The known-absent direct-`unload` rule remains `unsupported_capability` with no frame.
- D10(a) is independent of D10(b); selecting one does not select the other.
- No public API, new Wails binding, or contract amendment is authorized by this record.
- C1 implementation remains unauthorized until a separate stage authorization.

**Credible options.** The proposal established three candidate directions. The owner selected A10a-2 below; the other options remain historical alternatives, not selected.

| Option | Source basis | Required owner specificity | Consequence if selected |
|---|---|---|---|
| A10a-1 — Explicit semantic extension of an existing code | `[PROPOSAL]` proposal line 684 | Name the existing code and define its caller-facing meaning for unknown-snapshot `unload`; confirm the extension is contract-compatible | C1 may assert that concrete result in the unknown test; C1 acceptance still requires D10(b) and separate authorization |
| A10a-2 — Treat the unverified capability as absent | `[PROPOSAL]` proposal line 684; known-absent result `[CONTRACT]` A2-D1/A2-D7 | State explicitly that unknown is being mapped to the known-absent policy; the known-absent caller-facing result is `unsupported_capability` with no frame | C1 may align the unknown result with the known-absent result; no frame/no spawn remains; C1 acceptance still requires D10(b) |
| A10a-3 — Another concrete caller-facing result | `[PROPOSAL]` proposal line 684 | Name the exact result/code and its caller-visible semantics; preserve no frame/no spawn | C1 may assert that result; contract compatibility and any new surface require separate review; C1 acceptance still requires D10(b) |

`UNESTABLISHED` for any other result. “Defer the branch” is not an option: it names no caller-facing behavior and cannot satisfy C1 acceptance.

**Consequences and gates.**

- The D10(a) selection clears the D10(a) half of C1’s acceptance gate.
- D10(b) is recorded separately below and clears the other half.
- C7 remains downstream of C1 and therefore cannot claim D10(a) acceptance transitively before C1 is accepted.
- No option authorizes C1, C7, or any later stage.

**Owner selection:** **A10a-2 — Treat the unverified capability as absent.** Accepted by Stone Wolf Systems on 2026-11-08. When the authoritative service-capability snapshot is unknown, a direct `Client.Unload` SHALL return `unsupported_capability` at the helper boundary and SHALL project it through the existing `AudioError` mapping (`ErrCodeSynthesisFailed`), preserving the precise code and bounded diagnostic detail `capability=howl.tts.unload` and `snapshot=unknown`. It SHALL NOT start, spawn, or restart a helper; SHALL NOT write an `unload` operation frame; SHALL NOT report provider state `unloaded`; and SHALL NOT claim that the provider was unloaded. “Unknown” means no current helper instance has an authoritative compatible-handshake snapshot under A2-D6; a stale `hello` after stop, exit, crash, incompatibility, backoff, or replacement is not known. This is a fail-closed mapping of an unverified capability to the known-absent policy. It does not alter the known-present or known-absent branches, A2-D1, A2-D7, or A2-D9. No new public error taxonomy is created.

**Owner rationale:** Accepted the engineering recommendation for fail-closed behavior; no additional owner rationale was supplied.

**Date decided:** 2026-11-08

**Future implementation/verification consequences — not authorization.** C1 may now assert the recorded result and must retain the no-frame/no-spawn assertions. C1 must also test the stale-after-exit case as unknown, not as a known-present capability. C1 remains unauthorized until a separate stage authorization.

---

## 5. D10(b) — Unknown-snapshot direct `Cancel`

**Exact owner question.** When the authoritative service-capability snapshot is unknown, what concrete caller-facing result does a direct `Cancel` return outside the supervisor cancellation path, independently of D10(a), while preserving A2-D1’s known-absent rule and not weakening A2-D9?

**Status:** `DECIDED — A10b-2 — 2026-11-08`

**Why this was necessary.** A2-D1 settles the known-absent direct operation. A2-D9 settles cancellation without `cancel` on the supervisor path. Neither settled the unknown-snapshot result of the exported direct `Cancel`. Before this selection, C1 could test the structural no-frame/no-spawn boundary but could not be accepted without a concrete result. This branch is not an automatic copy of D10(a); the owner selected it independently.

**Governing contract.**

- `[CONTRACT]` §17.19.2 A2-D1 (line 2502): absent optional service capability → `unsupported_capability`, no operation frame.
- `[CONTRACT]` §17.19.2 A2-D9 (line 2510): on the supervisor cancellation path, absent `cancel` skips the frame, synthesizes no `unsupported_capability`, and follows the existing grace/termination rules.
- `[CONTRACT]` §17.11.6 (lines 2113–2150): cancellation precedence and terminal-publication rules remain authoritative.
- `[CONTRACT]` §17.7.2 (line 1773): `cancel` is not a lazy-start trigger.
- `[CONTRACT]` §17.9 (`unsupported_capability` row line 1918): the code is scoped to an optional capability that is unavailable.

**Proposal evidence.**

- `docs/review/DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md` §4.4 (lines 121–123): the supervisor path skips the frame when `cancel` is absent or unknown; direct `Cancel` has a separate unknown-snapshot result recorded as D10(b).
- Same proposal, C1 (lines 137, 159, 177): the direct-`Cancel` unknown test is structural only; C1 is not accepted until D10(b) has a concrete caller-facing disposition.
- Same proposal, §10 (line 684): D10(b) is independently answerable and is not assumed to share unload’s disposition.
- `[AUDIT]` Focused D10-scope audit §4.2–§4.3 (lines 95–149): A2-D9 remains unmodified; A2-D1 does not state an unknown-snapshot result; the split is not a contradiction.

**Constraints that cannot be changed by this decision.**

- The A2-D9 supervisor cancellation path remains exactly as settled: record `cancel_requested`; send a cancel frame only when `cancel` is advertised; when absent or unknown, skip the frame, synthesize no `unsupported_capability`, and follow the existing grace/`KillTree`/`Wait()`/`cancelled` or `timeout` rules.
- The direct `Cancel` known-absent behavior remains `unsupported_capability` with no frame.
- The direct unknown branch sends no frame and starts no helper.
- D10(b) is independent of D10(a); the owner may select the same or a different result.
- The path-scoped internal cancellation operation must not inherit or leak the direct-operation result.
- No public API, new Wails binding, or contract amendment is authorized by this record.

**Credible options.** The proposal established three candidate directions. The owner selected A10b-2 below; the other options remain historical alternatives, not selected.

| Option | Source basis | Required owner specificity | Consequence if selected |
|---|---|---|---|
| A10b-1 — Explicit semantic extension of an existing code | `[PROPOSAL]` proposal line 684 | Name the existing code and define its caller-facing meaning for direct `Cancel` with an unknown snapshot; confirm contract compatibility | C1 may assert that result for direct `Cancel`; A2-D9’s supervisor path remains untouched; C1 acceptance still requires D10(a) |
| A10b-2 — Treat the unverified capability as absent | `[PROPOSAL]` proposal line 684; known-absent direct rule `[CONTRACT]` A2-D1 | State explicitly that unknown is mapped to the known-absent policy; known-absent direct `Cancel` is `unsupported_capability` with no frame | C1 may align direct `Cancel` with the known-absent result; the supervisor path still does not use that result; C1 acceptance still requires D10(a) |
| A10b-3 — Another concrete caller-facing result | `[PROPOSAL]` proposal line 684 | Name the exact result/code and caller-visible semantics; preserve no frame/no spawn | C1 may assert that result; contract compatibility and any new surface require separate review; A2-D9 remains unchanged; C1 acceptance still requires D10(a) |

`UNESTABLISHED` for any other result. “Defer the branch” is not an option: it names no caller-facing behavior and cannot satisfy C1 acceptance.

**Consequences and gates.**

- The D10(b) selection clears the D10(b) half of C1’s acceptance gate.
- D10(a) is recorded separately above and clears the other half.
- No option authorizes C1 or any later stage.
- No option changes A2-D9’s supervisor-path behavior.

**Owner selection:** **A10b-2 — Treat the unverified capability as absent.** Accepted by Stone Wolf Systems on 2026-11-08. When the authoritative service-capability snapshot is unknown, a direct `Client.Cancel` outside the supervisor cancellation path SHALL, after target-ID validation, return `unsupported_capability` at the helper boundary and SHALL project it through the existing `AudioError` mapping (`ErrCodeSynthesisFailed`), preserving the precise code and bounded diagnostic detail `capability=howl.tts.cancel` and `snapshot=unknown`. It SHALL NOT start, spawn, or restart a helper; SHALL NOT write a cancel operation frame; and SHALL NOT claim that cancellation succeeded. A malformed target SHALL continue to return `invalid_request` before the snapshot gate. “Unknown” means no current helper instance has an authoritative compatible-handshake snapshot under A2-D6; a stale `hello` after stop, exit, crash, incompatibility, backoff, or replacement is not known. This is a fail-closed mapping of an unverified capability to the known-absent policy. A2-D9’s supervisor cancellation path remains unchanged: it records `cancel_requested`, uses a path-scoped internal operation, skips the frame when `cancel` is absent or unknown, synthesizes no `unsupported_capability`, and follows the normal grace/`KillTree`/`Wait()` and cancellation-precedence rules. No new public error taxonomy is created.

**Owner rationale:** Accepted the engineering recommendation for fail-closed behavior; no additional owner rationale was supplied.

**Date decided:** 2026-11-08

**Future implementation/verification consequences — not authorization.** C1 may now assert the recorded direct-`Cancel` result and must separately assert that the supervisor path still skips the frame and emits no `unsupported_capability`. The direct and supervisor paths must remain separately implemented. C1 remains unauthorized until a separate stage authorization.

---

## 6. D1a — Loader-seam permission

**Exact owner question.** May the implementation add a narrow internal, non-Wails loader seam outside §17.12.1’s enumerated adapter behavior, or does explicit load reachability require a contract clarification/adoption-record decision?

**Status:** `PENDING — OWNER DECISION REQUIRED`

**Why this is necessary.** A2-D12 makes an explicit `load` the only recovery trigger from `unavailable`. The §17.12.1 adapter table has no `load`, `unload`, or `cancel` entry point. Without a disposition, C3 cannot know whether it may introduce the adapter/session seam, and C7 cannot claim explicit load reachability.

**Governing contract.**

- `[CONTRACT]` §17.12.1 (lines 2255–2265): the adapter behavior table lists `Initialize`, `Synthesize`, `SynthesizeStream`, `GetVoices`, `SupportsSSML`, `Release`, `GetConfig`, `ProviderStatus`, and `ProviderCatalog`; it has no `load`, `unload`, or `cancel` entry point.
- `[CONTRACT]` §17.19.2 A2-D12 (line 2513): recovery from `unavailable` is triggered only by a new explicit `load` request; ordinary `unloaded` synthesis retains its ensure-handshake/load path.
- `[CONTRACT]` §17.7.2 (line 1784): `unavailable` + explicit same-identity `load` transitions to `loading`.
- `[CONTRACT]` §17.12 (lines 2158–2174): `types.TTSSynthesizer` remains unchanged.

**Proposal evidence.**

- `docs/review/DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md` §1 (line 36): D1a asks whether a narrow internal loader seam outside the table is permitted or whether an adoption-record clarification is required.
- Same proposal, §4.4 (line 127): explicit `load` is the only recovery trigger, and §17.12.1 currently defines no production entry point.
- Same proposal, C3 (lines 249–265): C3’s design is conditional on D1a; if permitted, C3 introduces a narrow seam; if not, it is limited to state gating and an internal `Load` on the existing client, with no new adapter/session file.
- Same proposal, C7 (lines 458–463) and §10 (line 674): D1a gates C3 before any seam; C7 consumes the already-decided surface.

**Constraints that cannot be changed by this decision.**

- No public Wails API or recovery UI is authorized.
- The low-level `helper.Client` remains transport/state and must not infer application asset policy.
- D1a must be resolved before C3 introduces any seam.
- If the seam is denied, the source-established alternative is C3 state gating plus an internal `Load` on the existing client, with no new adapter/session file.
- C7 does not retroactively authorize a seam that C3 was denied.
- No contract amendment is made by this record.

**Credible options.**

| Option | Source basis | C3 consequence | C7 consequence / gate |
|---|---|---|---|
| A1a-1 — Permit a narrow internal, non-Wails loader seam outside §17.12.1 | `[PROPOSAL]` proposal lines 460, 674 | C3 may introduce the narrow seam; no manager or Wails wiring; C3 still requires separate stage authorization | C7 may consume the seam only after C1–C6 stability and its other gates; this selection alone does not authorize C7 |
| A1a-2 — Require a contract clarification/adoption-record decision before permitting a seam | `[PROPOSAL]` proposal lines 460, 674 | C3 cannot add the seam; C3 remains blocked until the clarification is recorded; if the clarification denies the seam, C3 is limited to state gating plus internal `Load` and creates no adapter/session file | C7 explicit load reachability remains blocked until the clarification and C3; no production surface is invented |

`UNESTABLISHED` for any other surface. The proposal does not establish a public loader API or a Wails binding.

**Consequences and gates.**

- D1a gates C3. C3 cannot add any seam before this is dispositioned.
- C7 remains downstream of C3 and C1–C6; D1a does not by itself authorize C7.
- If the seam is denied, the limited C3 path is the permitted alternative, but C7 explicit recovery remains blocked unless another authorized surface exists.
- No option authorizes implementation.

**Owner selection (owner to complete):**

**Owner rationale (owner to complete):**

**Date decided (owner to complete):**

**Future implementation/verification consequences — not authorization.** If a seam is permitted, C3 must prove the adapter/session is narrow and free of Wails/manager wiring, and C7 must test explicit recovery without changing `TTSSynthesizer`. If denied, C3 must prove no new adapter/session file exists and C7 must record reachability as blocked rather than inventing a surface.

---

## 7. D9 — Idle observation and deadline reset

**Exact owner question.** Does status/catalog/health observation reset the 120-second idle deadline?

**Status:** `PENDING — OWNER DECISION REQUIRED`

**Why this is necessary.** The contract supplies a 120-second idle-close trigger and A2-D8 supplies the shutdown mechanism, but it does not state whether observation resets the idle deadline. C5 cannot freeze observation-reset tests or be accepted without this answer.

**Governing contract.**

- `[CONTRACT]` §17.8 (line 1864): idle close is a provisional 120-second policy; it stops an idle helper instance.
- `[CONTRACT]` §17.19.2 A2-D8 (line 2509): idle close is a graceful `shutdown`, not an `unload`; it is available even when `howl.tts.unload` is absent; exit 0 → `absent`, nonzero/grace expiry → `crashed`.
- `[CONTRACT]` §17.7.2 (line 1773): an explicit TTS/catalog/status request from `absent` is a lazy-start trigger. This is settled and is not D9.
- `[CONTRACT]` §17.7.2 (line 1792): `ready`/`degraded`/`unloaded` + accepted `shutdown` → `stopping`.

**Proposal evidence.**

- `docs/review/DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md` §1 (line 36): D9 asks only whether status/catalog/health observation resets the 120 s idle deadline.
- Same proposal, C5 (line 360): the contract supplies the 120 s trigger but not the reset rule; C5 must not assume an answer; status-from-`absent` is contract-settled and is not part of D9.
- Same proposal, §10 (line 683): D9 is a genuine product decision; poll suppression is an application-layer policy, not a C5 lifecycle exception.

**Constraints that cannot be changed by this decision.**

- Status-triggered lazy startup from `absent` is not reopened; §17.7.2 settles it.
- Idle close remains a `shutdown`, never an `unload`.
- The 120-second value remains the provisional contract value; this decision does not select a new timeout.
- Any application-layer polling policy remains separate from the helper lifecycle rule.
- C5 cannot be accepted until D9 is dispositioned.
- D9 does not authorize C5 or any later stage.

**Credible options.** The source establishes a binary lifecycle question.

| Option | Source basis | Required owner specificity | Consequence and gate |
|---|---|---|---|
| A9-1 — Observation resets the 120 s idle deadline | `[PROPOSAL]` proposal lines 360, 683; `[CONTRACT]` §17.8 line 1864 and A2-D8 line 2509 | State whether “observation” means all three of status, catalog, and health, or a specified subset; the source groups them and does not establish partial treatment | C5 adds observation-reset tests and acceptance for the selected observations; C5 remains unauthorized until separately authorized |
| A9-2 — Observation does not reset the 120 s idle deadline | `[PROPOSAL]` proposal lines 360, 683; `[CONTRACT]` §17.8 line 1864 and A2-D8 line 2509 | State the same observation scope if the answer is not uniform | C5 tests that observations do not extend the deadline; idle close follows the fixed deadline; C5 remains unauthorized until separately authorized |

`UNESTABLISHED` for any more granular rule not stated by the owner. If different observation types receive different treatment, that is a new specificity the source does not supply and must be written explicitly.

**Consequences and gates.**

- D9 gates C5 acceptance only; it does not authorize C5.
- Status-from-`absent` lazy start remains contract-owned and is verified at C7, not as a D9 choice.
- Polling suppression, if desired, is an application-layer decision outside this lifecycle question.

**Owner selection (owner to complete):**

**Owner rationale (owner to complete):**

**Date decided (owner to complete):**

**Future implementation/verification consequences — not authorization.** C5 will add only the observation-reset tests that match the owner’s selected scope. It must separately preserve the §17.7.2 status-from-`absent` lazy-start behavior.

---

## 8. D1b — Explicit load action and adapter lifetime

**Exact owner question.** Which existing explicit product action invokes the permitted loader, and does the adapter survive that action for genuine same-identity in-place recovery?

**Status:** `PENDING — OWNER DECISION REQUIRED`

**Why this is necessary.** D1a answers whether a loader surface is permitted, not which product action reaches it. A2-D12 requires an explicit `load` for recovery. If the existing Apply/rebuild path destroys the adapter, the result is not in-place recovery and C7 must not describe it as such.

**Governing contract.**

- `[CONTRACT]` §17.19.2 A2-D12 (line 2513): recovery from `unavailable` is triggered only by a new explicit `load` request; same-identity recovery is limited to a live helper that accepts a load, and different identity uses A2-D7 stop-and-replace.
- `[CONTRACT]` §17.12.1 (lines 2255–2269): the adapter table defines no explicit load entry point and `GetTTSVoiceCatalog` remains a no-context Wails signature.
- `[CONTRACT]` §17.7.2 (line 1773): startup from `absent` is limited to explicit TTS/catalog/status requests; it is not a general loader trigger.

**Proposal evidence.**

- `docs/review/DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md` C7 (lines 458–463): D1b must identify the existing explicit action and whether the adapter survives it; loading must not be invoked from `Synthesize`, status, health, catalog, `GetVoices`, or startup.
- Same proposal, C7 (line 462): if Apply rebuilds the generation or destroys the adapter, do not describe that as in-place recovery; preserving the adapter or adding an internal retry command requires D1b authorization.
- Same proposal, §10 (line 675): the existing settings Apply/rebuild path is cited as context; D1b is still an open product/architecture decision.

**Constraints that cannot be changed by this decision.**

- No recovery UI and no public Wails binding are authorized.
- Loading must not be triggered from `Synthesize`, status, health, catalog, `GetVoices`, or startup.
- Same-identity in-place recovery requires a live helper that accepts a load; a destroyed/rebuilt adapter is not in-place recovery.
- D1b does not authorize a new action or command; any new internal retry command or adapter-preservation change requires explicit D1b authorization.
- D1b gates C7 production reachability, not C3.
- No stage is authorized by this decision.

**Credible options.**

| Element | Source basis | Owner must specify | Consequence and gate |
|---|---|---|---|
| Trigger action | `[UNESTABLISHED]` The proposal names the existing settings Apply/rebuild path only as context (line 675); it does not establish that this path invokes the loader or which other action does | Name the existing explicit product action, or state that none exists | Without a named action, C7 cannot claim explicit load reachability; C7 production reachability remains gated |
| Adapter lifetime | `[PROPOSAL]` lines 461–462 | State whether the adapter survives the action | If it survives, genuine same-identity in-place retry is possible; if it is rebuilt/destroyed, do not claim in-place recovery and C7 must record the reachability limit |

`UNESTABLISHED` for any new product action, public API, or recovery UI. The proposal does not define one.

**Consequences and gates.**

- D1b gates C7 production reachability.
- If no existing action and no adapter survival are established, C7 explicit recovery remains blocked.
- Preserving the adapter or adding an internal retry command is future work requiring D1b authorization; neither is selected here.
- No option authorizes C7.

**Owner selection — trigger action (owner to complete):**

**Owner selection — adapter lifetime (owner to complete):**

**Owner rationale (owner to complete):**

**Date decided (owner to complete):**

**Future implementation/verification consequences — not authorization.** C7 will need a caller-map test proving the loader is reachable only from the selected action and that the adapter is preserved only if the owner selected genuine in-place recovery. If the action rebuilds the adapter, C7 must not claim same-identity in-place recovery.

---

## 9. D2 — Helper-backed engine activation

**Exact owner question.** How should a helper-backed engine be activated in settings without changing `TTSSynthesizer`, and what is its relationship to O-4?

**Status:** `PENDING — OWNER DECISION REQUIRED`

**Why this is necessary.** C7 can build and fake-helper-test the adapter, but production activation is not authorized until D2 is decided. The activation choice also determines how the helper-backed engine relates to the still-open O-4 backend-switch question. O-4 must not be silently closed.

**Governing contract.**

- `[CONTRACT]` §17.12 (lines 2158–2174): `types.TTSSynthesizer` remains unchanged.
- `[CONTRACT]` §17.12.1 (lines 2255–2265): adapter methods and their v1 behavior; `Initialize` stores config without spawn or load.
- `[CONTRACT]` §17.15 (line 2351): O-4 — when the existing `kokoro` Python bridge becomes helper-backed — remains open.
- `[CONTRACT]` §17.19.2 A2-D11 (line 2512): O-1 through O-7 remain unresolved and are not decided by §17.19.

**Proposal evidence.**

- `docs/review/DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md` C7 (lines 429–478): production adapter and explicit load; no production activation until D2; adapter remains behind an internal default-off activation gate.
- Same proposal, §9.3 (line 664): O-4 default backend switching is not decided by “integration complete.”
- Same proposal, §10 (line 676): D2 concerns activation in settings without changing `TTSSynthesizer`, and its relationship to O-4; O-4 remains open.
- Same proposal, §11 (line 701): adapter activation stays behind an internal default-off gate; no Wails setting until D2.

**Constraints that cannot be changed by this decision.**

- `TTSSynthesizer` remains unchanged.
- No native provider, Kokoro/G2P/ORT asset, model download, or packaging work is authorized.
- O-4 remains open; D2 must not silently decide when Kokoro becomes helper-backed.
- O-1, O-3, O-5, O-6, and O-7 are not decided here.
- C7 production activation remains gated; C7 adapter construction and fake-helper testing do not equal production activation.
- No stage is authorized.

**Credible options.** `UNESTABLISHED`.

The source establishes constraints, not an activation design. The proposal says activation is a settings decision, keeps the adapter behind an internal default-off gate, and forbids production activation before D2; it does not establish a specific settings surface, default, rollout, or O-4 relationship. The owner must clarify:

1. Which existing settings path or explicit user action activates a helper-backed engine.
2. Whether activation is global, per-engine, or otherwise scoped.
3. How activation relates to O-4 without closing it.
4. What remains default-off until a later decision.

No option set is inferred from implementation convenience.

**Consequences and gates.**

- Without a D2 selection, C7 production activation remains blocked.
- C7 may still build and fake-helper-test the adapter behind the default-off gate, subject to separate stage authorization.
- O-4 remains open regardless of D2; D2 cannot be treated as O-4 closure.
- No option authorizes native provider work.

**Owner selection (owner to complete):**

**Owner rationale (owner to complete):**

**Date decided (owner to complete):**

**Future implementation/verification consequences — not authorization.** C7 will need an activation test proving `TTSSynthesizer` is unchanged, the default remains off until the selected activation, and no native provider is required. Any settings/UI surface remains future work and is not authorized here.

---

## 10. D5 — Renderer exposure and picker behavior

**Exact owner question.** What renderer exposure shape is approved, and may the Character voice picker start a helper for `catalog` (O-1), or must it show static/cached data until TTS is applied?

**Status:** `PENDING — OWNER DECISION REQUIRED`

**Why this is necessary.** C7 can construct an adapter without UI, but production renderer exposure and the O-1 picker question remain open. The frozen contract keeps the renderer-facing catalog shape unchanged and does not authorize eager helper startup from `GetTTSVoiceCatalog`. A UI decision cannot be inferred from implementation convenience.

**Governing contract.**

- `[CONTRACT]` §17.12.1 (lines 2269–2272): `GetTTSVoiceCatalog(engine string)` retains its existing no-context Wails signature; integration must keep the static profile fallback or add a bounded internal inspector path; whether the picker may start a helper is O-1; the renderer-facing `TTSVoiceCatalog`/`TTSVoiceOption` JSON shape remains unchanged.
- `[CONTRACT]` §17.15 (line 2348): O-1 asks whether the Character voice picker may start a helper for `catalog`, or must show static/cached data until TTS is applied.
- `[CONTRACT]` §17.19.2 A2-D11 (line 2512): O-1 through O-7 remain unresolved.

**Proposal evidence.**

- `docs/review/DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md` C7 (line 456): `GetTTSVoiceCatalog` uses the static fallback and does not start a helper; O-1 remains open.
- Same proposal, C7 (line 475): UI and O-1/O-3/O-4 activation are excluded from C7.
- Same proposal, §10 (line 679): D5 covers renderer exposure shape and O-1 picker behavior; both stay static until decided; D4 is folded into D2/D5 with no separate gate.

**Constraints that cannot be changed by this decision.**

- The renderer-facing `TTSVoiceCatalog`/`TTSVoiceOption` JSON shape remains unchanged unless a separate authorized decision changes it.
- `GetTTSVoiceCatalog` remains no-context and synchronous.
- The current C7 design uses the static fallback and does not start a helper from `GetTTSVoiceCatalog`.
- O-1 is not decided by this record; D5 cannot silently close it.
- No UI, Wails binding, or recovery surface is authorized.
- D5 gates production UI integration, not C7 adapter construction.

**Credible options.**

| Question | Credible option | Source basis | Consequence |
|---|---|---|---|
| O-1 picker startup | A5-1 — The Character voice picker may start a helper for `catalog` | `[CONTRACT]` §17.15 line 2348 | A future UI design must address the no-context/synchronous blocker; no UI work is authorized here |
| O-1 picker startup | A5-2 — The picker shows static/cached data until TTS is applied | `[CONTRACT]` §17.15 line 2348; current proposal line 456 | The existing static fallback remains; `GetTTSVoiceCatalog` does not start a helper; no UI work is authorized here |
| Renderer exposure shape | `UNESTABLISHED` beyond the frozen existing shape | `[CONTRACT]` §17.12.1 lines 2269–2272 | Any new exposure shape requires owner specificity and separate design/authorization; the existing shape remains the only frozen shape |

The source establishes the two O-1 alternatives; it does not establish a new renderer shape or a partial picker policy.

**Consequences and gates.**

- D5 gates production UI integration; it does not authorize UI work.
- Neither O-1 option changes C7’s adapter requirements or the frozen renderer JSON shape by itself.
- D4 has no separate gate; it is folded into D2/D5.
- No option authorizes C7 or C8.

**Owner selection — renderer exposure (owner to complete):**

**Owner selection — O-1 picker behavior (owner to complete):**

**Owner rationale (owner to complete):**

**Date decided (owner to complete):**

**Future implementation/verification consequences — not authorization.** Any future UI integration must show that the picker does or does not start a helper exactly as selected, that the existing JSON shape is preserved unless separately changed, and that no helper starts from an inspection path unless the owner selected that behavior.

---

## 11. D3 — Third-process admission rejection

**Exact owner question.** When a third real helper process would exceed the frozen two-process cap, what caller-facing result does the adapter return when rejecting it before `proc.Start`?

**Status:** `PENDING — OWNER DECISION REQUIRED`

**Why this is necessary.** The cap is frozen at two actual helper processes, but §17.9 contains no entry that exactly defines admission rejection. C7 may implement the counting/admission boundary, but it cannot claim a contract-compliant cap response until the caller-facing result is selected or the contract gap is clarified.

**Governing contract.**

- `[CONTRACT]` §17.8 (line 1857): two helper processes are a frozen hard cap — global plus one explicit engine.
- `[CONTRACT]` §17.19.1 (line 2496): the clarification does not change the two-process cap.
- `[CONTRACT]` §17.9 (lines 1904, 1906, 1918):
  - `busy` is wire-origin and means one non-control operation is active.
  - `helper_start_failed` means process/job/pipe start failed.
  - `unsupported_capability` is for a requested optional capability that is unavailable.
- `[CONTRACT]` §17.9 registry (lines 1876–1920): no entry exactly defines admission rejection for the third process.

**Proposal evidence.**

- `docs/review/DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md` C7 (lines 466–471): the cap counts actual helper child processes, not adapters, character voices, engine labels, or requests; admission is checked before `proc.Start`.
- Same proposal, C7 (line 471): `busy`, `helper_start_failed`, and `unsupported_capability` do not exactly define the admission rejection; D3 remains a narrow contract gap; do not invent or repurpose a code.
- Same proposal, §10 (line 677): C7 may implement the counting/admission boundary without asserting the error; the error surface remains open.

**Constraints that cannot be changed by this decision.**

- The cap remains two actual helper child processes per app session.
- The count is over real child processes, not adapters, engine labels, character voices, or requests.
- Admission is checked before `proc.Start`.
- `busy`, `helper_start_failed`, and `unsupported_capability` retain their existing contract meanings; none may be repurposed without explicit contract support.
- No new error code or contract amendment is authorized by this record.
- C7 may implement counting/admission but cannot claim a contract-compliant cap response until D3 is resolved.
- D3 gates C7 cap-response acceptance, not the entire C7 stage.

**Credible options.** `UNESTABLISHED`.

The source does not establish a contract-compliant caller-facing error. The following are not credible selections without further contract support:

- `busy` — its meaning is one active non-control operation, not process-cap admission.
- `helper_start_failed` — its meaning is a process/job/pipe start failure, not rejection before start.
- `unsupported_capability` — its meaning is an absent optional service capability, not the process cap.

**Clarification needed.** The owner must either identify an existing code whose normative meaning covers admission rejection before `proc.Start`, or decide that a contract clarification/new error surface is required. No code is selected by this record. If the existing codes remain a poor fit, the unresolved gap must stay visible rather than being hidden behind a reused code.

**Consequences and gates.**

- Without D3, C7 can implement and test the actual-process count and pre-start admission boundary, but must not assert a contract-compliant caller-facing rejection.
- C7 cap-response acceptance remains blocked.
- C7 production activation also remains gated by D1b and D2.
- No option authorizes C7 or a contract change.

**Owner selection (owner to complete):**

**Owner rationale (owner to complete):**

**Date decided (owner to complete):**

**Future implementation/verification consequences — not authorization.** C7 will need process-count tests covering shared processes, adapters, labels, and reaped slots, plus an admission test that runs before `proc.Start`. If no error is selected, the test must record the response as unestablished rather than asserting a reused code.

---

## 12. Dependency and gate matrix

| Decision | Gates | Upstream / sequencing | Effect |
|---|---|---|---|
| D10(a) | C1 acceptance | C1 design; independent of D10(b) | Decision recorded: A10a-2; the D10(a) gate is cleared; C1 remains unauthorized |
| D10(b) | C1 acceptance | C1 design; independent of D10(a) | Decision recorded: A10b-2; the D10(b) gate is cleared; C1 remains unauthorized |
| D1a | C3; C7 upstream | C2 status model; before any C3 seam | C3 cannot add a loader seam; C7 explicit reachability remains upstream-blocked |
| D9 | C5 acceptance | C4 exit/shutdown semantics | C5 observation-reset tests and acceptance cannot be finalized |
| D1b | C7 production reachability | D1a; C3; C1–C6 stability | C7 cannot claim a reachable explicit recovery path |
| D2 | C7 production activation | D1b; C7 adapter; O-4 remains open | C7 may build/test behind the default-off gate, but cannot activate production |
| D5 | Production UI integration | D2/D5 fold D4; O-1 open | No production UI exposure or picker startup may be claimed |
| D3 | C7 cap-response acceptance | C7 counting/admission boundary | C7 may count and reject before start, but cannot claim a compliant caller-facing error |

Stage summary:

| Stage | Design/owner gate | Authorization state |
|---|---|---|
| C1 | D10(a) A10a-2 and D10(b) A10b-2 recorded; C1 authorization still required | Unauthorized |
| C2 | No remaining design gate | Unauthorized |
| C3 | D1a before any loader seam; C2 sequencing | Unauthorized |
| C4 | No owner decision in this record; C1–C3 sequencing | Unauthorized |
| C5 | D9 for acceptance; C4 sequencing | Unauthorized |
| C6 | No owner decision in this record; C3–C5 sequencing | Unauthorized |
| C7 | D1a already before C3; D1b, D2, D3; D10 by C1 transitivity | Unauthorized |
| C8 | C1–C7 complete; separate runtime authorization | Unauthorized |

**No stage is authorized by a decision, a blank field, this record, or the frozen design.** C1 and C2 remain unauthorized. C8 requires separate runtime-execution authorization in addition to the design and stage gates.

## 13. Settled matters not reopened

The following are outside this decision instrument and must not be reopened by an owner selection:

- **D6** — forced termination is `stopping` then grace → `crashed`; request outcome remains `cancelled`/`timeout`; clean exit 0 → `absent`.
- **D7** — settled by §17.11.5: a load-only retry starts no process and consumes no process-restart budget; a retry that restarts a process consumes the budget and does not inherit A2-D7’s replacement exemption; stable readiness for ten minutes resets the ledger.
- **D4** — folded into D2/D5; no separate gate.
- **D8** — B12a/B12b remain deferred and are not decided here.
- **A2-D1 known-absent direct optional behavior** — `unsupported_capability`, no frame.
- **A2-D9 supervisor cancellation path** — remains settled and outside D10(a) and D10(b); it is not weakened by either branch.
- **A2-D12 ordinary `unloaded` synthesis** — the ensure-handshake/load path remains; only `unavailable` recovery requires an explicit `load`.
- **O-1, O-3, O-4, O-5, O-6, O-7** — remain open except where this record asks D5/D2 to address their relationship without closing them.

No selection in this record may rewrite these rules. A proposed answer that requires changing them is outside this instrument and must be returned as an explicit contract/design question.

## 14. Source ambiguity and missing evidence

The entries are complete as decision instruments. The table records residual source ambiguity; D10(a) and D10(b) are resolved and the remaining rows remain open in the frozen sources.

| Decision | Unresolved evidence | What the owner must clarify |
|---|---|---|
| D1a | The contract has no `load` entry point in §17.12.1, and the effect of an adoption-record clarification is not defined | Permit the narrow seam, or require a clarification and state its effect |
| D1b | No existing product action is established as the loader trigger; adapter survival is not established | Name the action and whether the adapter survives it |
| D2 | The settings activation design and its O-4 relationship are not established | Define activation scope/default and preserve O-4 as open |
| D5 | Only O-1’s two alternatives and the frozen existing renderer shape are established | Select O-1 and state whether any new shape is requested |
| D3 | No contract-compliant admission-rejection code exists in §17.9 | Identify a contract-supported existing code or require clarification/new surface |
| D9 | The contract supplies the 120 s idle close but not the observation-reset rule | Select reset or no reset, including observation scope |
| D10(a) | Resolved by owner selection A10a-2 on 2026-11-08; contract-silent branch now has a concrete result | None — recorded |
| D10(b) | Resolved by owner selection A10b-2 on 2026-11-08; contract-silent branch now has a concrete result | None — recorded |

No missing evidence prevented the eight entries from being stated. D10(a) and D10(b) are now selected; the remaining six decisions remain pending and must not be treated as selected.

## 15. Decision-record completion and authorization boundary

This is a decision instrument, not permission to implement. The owner may complete the fields in §§4–11; the six undecided decisions remain `PENDING — OWNER DECISION REQUIRED`. D10(a) and D10(b) are recorded with concrete results.

A completed selection must state the concrete caller-facing result where the question asks for one. The D10 selections satisfy that rule; the remaining six decisions must not be inferred from them.

No owner decision authorizes its stage. C1 and C2 remain unauthorized; C3–C7 require separate, stage-specific implementation authorization; C8 requires separate runtime-execution authorization. No build, test, replay, benchmark, launch, or runtime operation is authorized by this document.
