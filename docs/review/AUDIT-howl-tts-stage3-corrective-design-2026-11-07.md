# Adversarial Design Review — Stage 3 Corrective Sequencing Proposal

**Review date:** 2026-11-07
**Reviewer:** Mimo (independent, read-only)
**Subject:** `DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md`
**Governing contract:** `DESIGN-howl-tts-helper-architecture-2026-11-06.md` as adopted (includes §17.19/A2)
**Status of subject:** DESIGN PROPOSAL. This review authorizes nothing.

---

## 1. Verification of the proposal's claims

Everything the proposal asserts about its own inputs and its own footprint was checked, not accepted.

| Claim | Result |
|---|---|
| Proposal SHA-256 `71e64a93…` | **Verified.** `71e64a93c328c71f4f35ab1bfa84e4abbbeb9ad0c74a0be8c2cb60e20f98b12b`, 641 lines / 44,056 bytes. |
| README §24 added, hash-matched | **Verified.** §24 at line 455 of the index, hash identical to the file. |
| Write set: exactly two documentation files | **Verified.** Only the proposal (`09:14:11`) and `README.md` (`09:14:24`) were modified after the readiness investigation. The readiness investigation's own mtime (`2026-11-06 23:32`) predates both. |
| "No production changes, no tests" | **Verified.** No file outside `docs/` has been modified on 2026-11-07 at all. |
| Governing architecture `ec3c789e…` | **Verified.** `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f`, 2,586 lines / 182,904 bytes (post-adoption: §17.19 at line 2447). |
| Readiness investigation `a63f9f32…` | **Verified.** |
| Stage 2 implementation `ec85b2b2…` | **Verified.** |
| Stage 2 capability correction `c9f498f6…` | **Verified** (`CORRECTIVE-howl-tts-helper-capabilities-no-overwrite-2026-11-06.md`). |
| Stage 2 adversarial audit `e73cef36…` | **Verified.** |
| Stage 2 correction audit `d6b2349e…` | **Verified.** |

All six input hashes are correct. This is the cleanest input-hash table in the review series so far.

**Inventory accuracy.** The proposal names specific functions and types rather than gesturing at files, so I checked them against the code:

- `validateHello`, `validateCapabilities`, `helloHandshake`, `requireServiceCapability`, `Client.Unload`, `Client.Cancel` all exist where stated.
- The "five mandatory capabilities" claim is exact: `protocol.go:53` defines `mandatoryServiceCapabilities` as `{catalog, load, unload, synthesize, cancel}`, and `validateHello` (line 527) rejects a missing one as `protocol_violation` — which is the defect C1 exists to fix.
- `TTSSynthesizer` (`interfaces.go:127`) has six methods; the C7 adapter covers all six. `GetConfig` is not on that interface, and the proposal is still right to list it: §17.12.1 includes a `GetConfig` row, and `pipeline.go:390` type-asserts `interface{ GetConfig() types.TTSConfig }`. Omitting it would have been a real defect.
- `TTSCapability` is today `{ID, Version int}`; C2's addition of `Supported` and `Detail` matches A2-D6 exactly, as does adding `Loaded`, `ServiceCapabilities`, `ServiceCapabilitiesKnown` to `TTSProviderStatus`.
- `ttsEngineKey`, `ttsVoicePath`, `loadTTSForRequest`, `ttsEngines`, and the character override path exist as described.
- "The helper packages are currently inert from production" holds: no non-test file imports `audio/tts/helper`.

---

## 2. Verdict

**CONDITIONAL PASS — accept the decomposition, correct two specifications before authorizing C1 and C2.**

The program structure is sound and the scope discipline is real. The eight-stage decomposition is genuinely dependency-ordered, the C1–C6 boundary ("inert helper packages") is accurately drawn, C7 is correctly isolated behind a default-off gate, and C8 is correctly separated from implementation. The proposal is honest about what it does not establish, and it does not smuggle policy into an implementation plan.

Two specification defects nevertheless sit inside the first batch. Neither is a policy question and neither requires a new stage; both are a few sentences of the design document:

- **F1 (blocking for C1):** C1 conflates `unload` and `cancel` optional-capability handling. Adopted A2-D9 *forbids* the error C1 would generate.
- **F2 (blocking for C2):** C2's `unavailable` projection contradicts A2-D5 and the proposal's own C5.

If C1 and C2 were authorized as written, their acceptance tests would encode one contract violation and one missing state case — the exact failure mode ("tests encode superseded policy", S3R-16) the program exists to eliminate.

---

## 3. Priority 1 — C1 and C2 against §17.19

### F1 — BLOCKING. C1 merges two operations the adopted contract treats differently

C1 says:

> "`Unload`/`Cancel` handshake before gating if capability knowledge is absent, and emit no frame when absent."

and its test list says:

> "Optional operation unavailable → `unsupported_capability`, no frame emitted (assert transport send count/bytes)."
> "Direct `unload`/`cancel` after handshake and before handshake; no false 'unsupported' when the helper advertises the operation."

The adopted contract splits these two cases deliberately:

- **A2-D7:** "If `howl.tts.unload` is absent, the adapter returns `unsupported_capability` for a **direct unload request** and writes no frame."
- **A2-D9:** "If it did not [advertise `howl.tts.cancel`], the adapter skips the cancel frame; absence of `cancel` is not an error in the supervisor's cancellation path and **no `unsupported_capability` is synthesized there**."

So `Cancel` is *not* a direct optional request: it is the cancellation path, and A2-D9 exempts it from the error that A2-D1/A2-D7 assign to direct requests. C1's combined phrasing and its blanket "optional operation unavailable → `unsupported_capability`" row would make the cancellation path return `unsupported_capability`, reversing A2-D9 — and the test that asserts it would make the violation permanent.

**Required correction (C1 scope, tests, and review checklist):**

1. Split the two operations. `Unload` absent → `unsupported_capability`, no frame (A2-D1, A2-D7). `Cancel` absent in the cancellation path → skip the frame, no error, continue to the cancel grace and the existing precedence rules (A2-D9).
2. Add a test that asserts a `cancel`-less helper produces **no** `unsupported_capability` and still kills after grace with the request ending `cancelled`/`timeout`.
3. Add A2-D9 to C1's independent-review checklist alongside A2-D1–A2-D3.

**Second, smaller point in the same sentence.** "Handshake before gating if capability knowledge is absent" means a direct `unload` on a not-yet-started helper *starts a helper* in order to discover whether `unload` is advertised. The contract authorizes starting for `ProviderStatus` and `ProviderCatalog` (§17.12.1), not for `unload`/`cancel`, and `Release` is "no new work". A2-D1 and A2-D10 forbid *frames*, not spawns. C1 should state which fail-closed reading it implements when no authoritative snapshot exists — no frame and a typed not-loaded result without spawning, or spawn-then-gate — because the two have opposite eager-start consequences and the contract does not choose between them.

### F2 — BLOCKING. C2's `unavailable` projection omits one case and forbids the other

C2 states:

> "`Available` is true for a configured, compatible provider in `unloaded`/`loading`/`ready`/`busy`/`degraded`, and for the contract-defined **asset-specific** `unavailable` **on a live helper**."

> "`ProcessAlive` is true only for `starting`, `alive`, `stopping`; false for `absent`, `exited`. **Never suppress a crashed/unavailable aggregate by saying a process is alive.**"

A2-D5 has two `unavailable` cases, and C2 covers only one of them:

> "`Available` is also true for `unavailable` when either **(a)** the failure is asset-specific and the helper, protocol, and provider configuration remain valid, or **(b)** the state was entered because the `crashed`/`backoff` restart budget was exhausted while the provider remains configured and compatible; in case (b), `ProcessAlive=false`, `Loaded=false`, `Ready=false`…"

Consequences:

- **Case (b) is missing.** C2's rule sets `Available=false` on budget exhaustion, while the proposal's own C5 requires "budget exhaustion → `unavailable`, process exited, **`Available=true`**". A literal C2 implementation contradicts A2-D5(b) and C5. C2's test list even includes "budget-exhausted unavailable", so the rule text and the test matrix disagree.
- **Case (a) is forbidden by the second sentence.** A process that is alive with an asset-specific `unavailable` is exactly the legal pair A2-D5(a) permits, and the pair C3 and C6 depend on ("live helper" recovery, `unavailable` → `loading`). As written the instruction reads as "an `unavailable` aggregate cannot have `ProcessAlive=true`", which inverts the case C3's load seam is built on.

**Required correction (C2 projection rules and tests):** state both cases explicitly with their flag sets — (a) `Available=true`, `ProcessAlive=true`, `Loaded=false`, `Ready=false`, `ServiceCapabilitiesKnown=true`; (b) `Available=true`, `ProcessAlive=false`, `Loaded=false`, `Ready=false`, `ServiceCapabilitiesKnown=false` — and narrow the process-liveness sentence to "`ProcessAlive` is never derived from aggregate `State`" (A2-D4's actual rule). Add both rows to the status matrix, since C3 and C6 depend on them.

### What C1–C2 get right

The rest of the capability correction preserves the §17.19 distinctions:

- Required-two/optional-three floor, terminal `protocol_incompatible` for a missing required capability with no backoff or restart (A2-D1, A2-D4).
- Semantic-vs-frame classification with `detail.missing` limited to required IDs, table order, no spaces or duplicates (A2-D3), including the correct carve-outs: duplicate JSON keys stay `invalid_frame`; `invalid_request` stays for the `hello` request itself.
- Service/provider scope separation, shared `synthesize` accepted in both arrays, no `supported` field on the wire, `Supported:false` only in the application projection (A2-D2).
- "ServiceCapabilities contains all five IDs only when `ServiceCapabilitiesKnown=true`; absent optional IDs appear with `Supported=false`… Do not synthesize five false values" (A2-D6) — this is the exact inversion the contract demands, and it is stated correctly.
- `Loaded` false on begin/failure/unload/exit/replacement and true through `stopping` until exit is observed (A2-D5) — correct in both C2 and C3.
- No silent downgrade or no-op (A2-D10); `health`/`shutdown` correctly excluded from capability gating; idle close kept as `shutdown`, not `unload` (A2-D8).

---

## 4. Priority 2 — dependency ordering in C3–C6

**The chain is genuinely acyclic except for one seam, and it is the seam the proposal half-acknowledges.**

Sound: C1 → C2 (wire semantics before projection) → C3 → C4 (exit intent before timers) → C5 (budget) → C6 (identity/replacement consumes budget and termination) → C7 → C8. The stated prohibitions against combining stages are correct and well-reasoned, especially "exit truth must precede timers/budgets; otherwise a timer can launder a crashed process into absent" and "C3 and C6… A2-D12 specifically distinguishes them".

### F3 — MODERATE. `Loaded` has two owners and one prerequisite cycle

- C2 defines the `Loaded` rule in full ("false before load, during load, on load failure, on unload, on exit/crash, and on replacement").
- C3 defines it again ("`Loaded` becomes false when a load begins, fails, unloads, exits, or replacement begins").
- C2's acceptance is "Every status field has a single documented source" — but `Loaded`'s writer (the load seam) arrives in C3, and C2 lists C3 as its successor. Meanwhile C3 lists C2 as its prerequisite.

This is a small cycle: C2 cannot satisfy its own acceptance criterion for `Loaded`, and two stages own one rule. Two clean resolutions, either acceptable:

1. C2 ships the fields and projects `Loaded` from the existing provider snapshot, with C3 owning the transition rewrite and restating the rule as a change; or
2. Move the `Loaded`/`Ready` acceptance row into C3 and have C2 state only the type and the "never derived from `Available`/`Ready`" invariant.

Also note C2's invalidation list covers `backoff` and replacement, which are C5 and C6 states. That is harmless as forward-declaration but cannot be accepted at C2; C2's matrix should mark those rows "not yet reachable" rather than claim coverage.

### Ordering points that are correct

- C4's intent enum precedes C5's budget, and C5's budget precedes C6's replacement accounting — the right order, because A2-D7's exemption and A2-D12's non-exemption are only distinguishable once exit intent and the ledger both exist.
- C6's ordering mirrors A2-D7 step by step, including "no start before exit is observed", the `internal_error` no-cancel rule, and the shutdown-error/transport-failure classifications.
- C4's classifications match the frozen rows: exit 0 → `absent`, nonzero → `crashed`, grace expiry/kill → `crashed`, `incompatible` stays `incompatible` with the process exited and no budget spent (A2-D4), and `unloaded` is never entered merely because a process stopped (A2-D7, A2-D8).
- C5's numbers are the frozen ones: 3-in-10-minutes budget, 1/2/5 s backoff, 120 s idle (§17.8 line 1864).
- C5's "crash/backoff only on the next explicit TTS request; no background loop" is the correct reading of the no-automatic-restart rule, and its "status/catalog/health do not reset the budget" matches the frozen avoid-autonomy intent.

### F7 — MINOR. The `timeout-cancel` intent name invites the wrong classification

C4's intent enum includes `timeout-cancel`, and its behavior text is correct ("after kill/wait, state is `crashed` under the existing stopping/grace rule"). But intent-aware exit exists in A2-D4 for exactly one case — termination after an incompatible handshake, which "MUST NOT be reclassified as `crashed`". A frozen §17.7.2 row says `stopping` + grace expiry → `crashed`. A future implementer reading "intent: timeout-cancel" may reasonably suppress `crashed` as "expected". State in C4 that the intent value records cause for diagnostics only and never overrides the frozen transition; only `incompatible` is intent-exempt.

---

## 5. Priority 3 — C7 against the real interface and character architecture

**Compatible on the points that matter, with one accounting ambiguity and one under-scoped decision.**

Verified against the code rather than assumed:

- **No accidental eager load.** §17.12.1's rows are reproduced faithfully: `Initialize` stores config without spawn or load; `GetVoices` returns the last snapshot and never spawns or blocks; `ProviderStatus` may handshake/start but must not load; `ProviderCatalog` may start, handshake, and call `catalog`, but must not load; `Release` is idempotent with no new work. The proposal's explicit "no helper start from `GetTTSVoiceCatalog`; O-1 remains open" is correct: that path stays a static fallback.
- **Character architecture.** Voices are per-request (`resolvedTTS.voice` flows into `SynthesisRequest`), while the load identity is `provider_id`/`engine`/`model_identity`. So a character switching voice on the already-loaded helper engine does **not** create an identity change and does **not** trigger A2-D7 churn — which is the failure mode worth worrying about, and the design avoids it by keeping voice out of the identity tuple (correctly, per A2-D7's definition). Piper's per-voice model path stays keyed in `ttsEngineKey`; non-piper keys by engine alone, so the helper adapter is shared rather than duplicated per character.
- **No fallback drift.** "Explicit engine failure does not fall back" matches A2-D10 and O-3's open status.
- **`TTSSynthesizer` unchanged** is achievable: six methods, all covered, with `GetConfig` supplied on the concrete type (the pipeline's type assertion, not the interface).

### F8 — MINOR. The process cap is counted over adapters, but the contract counts processes

C7 says "at most one request-scoped explicit engine adapter, preserving the two-process cap" and "A third concurrent helper-backed engine is rejected before process creation". Because `ttsEngineKey` returns the bare engine name for non-piper, two requests naming the same helper engine share one adapter and therefore one process — so the cap must be counted over **live helper processes**, with the adapter count derived from it. As written, an implementation could satisfy "one explicit adapter" while two adapters each own a process, or reject a second *request* that should have reused the first adapter. State the counting basis.

### D1 deserves more weight than the proposal gives it

The proposal is right that reachability is unresolved, but it frames D1 as "which button". It is larger than that. In the adopted §17.12.1 table there is **no `load` entry point at all** — `Synthesize`'s ensure-load is the only path that sends a `load`, and A2-D12 excludes it while `State=unavailable`. So "the contract requires explicit `load` to be reachable" is true, and the contract simultaneously defines no interface method that can produce that request. Choosing a settings action does not close that: the action still needs a method to call, and adding one to the adapter/interface surface is not described by §17.12.1.

That makes D1 two decisions, not one:

- **D1a — surface permissibility.** May the implementation add an internal loader seam outside §17.12.1's enumerated adapter behavior (the proposal's recommendation), or does reachability require a contract clarification/amendment? This is a contract-completeness question and should be answered in the adoption record, not left to implementation.
- **D1b — trigger and adapter lifetime.** Which explicit action invokes it, and whether the adapter survives that action so the retry is genuinely in-place rather than a rebuild (the proposal already states this half correctly).

The proposal's claim that it found "no new contradiction in the adopted contract beyond B12a/B12b and O-1–O-7" is therefore slightly overstated: the load-surface gap is not in that list, even though the document does surface it.

---

## 6. Priority 4 — which of D1–D7 are genuine decisions

The table mixes three different kinds of item. Sorting them is worth more than debating them case by case, because four of the seven can be closed by citation.

| ID | Classification | Basis |
|---|---|---|
| **D1** | **Genuine, but split** | Not a product preference: the adopted §17.12.1 surface has no load entry point. Split into D1a (is an out-of-table internal seam permitted, or is a contract clarification needed) and D1b (trigger + adapter lifetime). |
| **D2** | **Genuine** | Activation in settings without changing `TTSSynthesizer`, and its relationship to O-4, is a product/architecture choice. |
| **D3** | **Genuine, but constrained** | The cap is frozen ("Process pool cap — Frozen v1 hard cap of two helper processes"), the error surface is not. Answer it **from the existing §17.9 registry**; inventing a new code would be a contract amendment, not a design choice. |
| **D4** | **Fold into D2/D5** | Whether the inspector has a production consumer is wiring, not policy; it is already implicitly answered by "O-1 stays static". |
| **D5** | **Genuine** | Renderer exposure shape and O-1 picker behavior are owner decisions, and the proposal correctly keeps both static. |
| **D6** | **Already settled — close by citation** | Frozen §17.7.2: `stopping` + grace expiry → `crashed` (`KillTree`, no replay). A2-D9 governs the request outcome (`cancelled`/`timeout`) and A2-D8 governs clean shutdown (exit 0 → `absent`). C4 already implements this; D6 should record it as applied, not open. |
| **D7** | **Substantially settled — narrow the residual** | A2-D12 is explicit: "A2-D12 itself does not inherit that exemption: its in-place retry is not an identity replacement, and any process restart remains governed by §17.11.5." A load-only retry starts no process, so it consumes no restart; a retry that restarts a process consumes budget and never inherits A2-D7's exemption. The only genuinely open part is the ledger **reset** rule: whether a successful in-place recovery restarts the 10-minute stable-readiness window. Decide that, and D7 becomes a one-line statement. |

For the owner's framing — "which should be resolved by applying the already-adopted contract rather than reopening settled policy" — the answer is: **D6 and D7 are settled** (D7 modulo reset semantics), **D3 is bounded by the frozen cap and the frozen error registry**, **D4 is not a policy question**, and **D1 is a real decision the contract does not answer**. Only D1, D2, and D5 remain genuinely open, and D1 must be answered before C7 — not before C1.

---

## 7. C8 — your evidence-separation concern is only half addressed

The proposal does distinguish environments in the matrix (`Windows`, `Runtime`, `Recheck`, `Preserve`) and its §9.2 table assigns "Requires Windows boundary" per stage. C4's and C5's runtime-evidence lines correctly say the Windows boundary "is required in C8", and C1's says fake-helper real-process tests are sufficient *for C1* because Windows job behavior is inherited there. Those are the right calls.

### F4 — MODERATE. C8 lacks the rule that makes the distinction binding

C8's acceptance is "Every matrix row has **direct, indirect, or explicitly unestablished** evidence". That wording permits a fake-helper or Linux-CI pass to be recorded as "indirect" evidence for a row labeled Windows — which is precisely the substitution you want excluded. The proposal never states that a Windows-boundary row cannot be discharged by a fake-helper run.

**Required addition to C8:** an evidence rule, stated once and applied per row —

1. A row labeled Windows (suspended start/job ownership, `KillTree`/wait, nonzero-exit classification, forced kill, final symlink/junction handling) is satisfied only by execution on Windows against a real child process.
2. A fake-helper or cross-platform test may support such a row but can never close it; if the Windows run did not occur, the row is recorded **unestablished**, not "indirect".
3. Every C8 row in the evidence report names its environment and command; "passed" without an environment is not evidence.
4. C1's "Windows behavior is inherited" shortcut is scoped to C1 and may not be reused for C4–C7, which change the very paths involved.

That is a four-sentence change and it converts your concern from an instruction in this conversation into a rule the next implementer must follow.

---

## 8. Remaining minor defects

### F6 — MINOR (editorial, but normative text with a question mark)

C5 contains: "Status/catalog/health do not reset idle indefinitely? **Define:** status/health do not reset the idle deadline unless contract explicitly says; accepted work/load resets it."

A stage specification cannot carry the reviewer's own open question and a "Define:" instruction. The contract supplies the trigger (A2-D8) but not the reset semantics, so this is a real (small) decision; record it as a one-line rule in C5, or move it to the decision table. The current wording would be copied into an implementation note verbatim.

### Minor accuracy nits

- C7 lists `app-audio.go` `applyAudioSettings` as the hook. I verified that function exists; the proposal should also record that the Apply path's existing behaviour (rebuilding the audio generation) is what makes D1b nontrivial — it already says this, which is good, but the two facts should sit together.
- §7's inventory lists `helper/client.go` `Unload`, `Cancel` under C1 and again under C1 for gating; if they are touched only for gating, say so, so a reviewer does not expect behavior changes in C1.
- The stage table's S3R-17 mapping ("C4, C5, C8") is right, but the C8 row of §9.2 says Windows is required "Yes" while the S3R-17 row implies C8 is where it closes — consistent, but worth one explicit sentence that C4/C5/C6 produce the *code* and C8 produces the *evidence*.

---

## 9. What this review does and does not establish

**Established by direct inspection:** the proposal's hash, its input hashes, its write set, its inventory accuracy against the current source, the absence of a production importer, the interface shapes, the character-voice path and its identity implications, the frozen numeric policies, and the adopted-contract text of A2-D1, A2-D5, A2-D6, A2-D7, A2-D8, A2-D9, A2-D12, §17.7.2, §17.12.1, and §17.9.

**Not established:** that any stage is implementable in the estimated time. The same-day boundary in §9.1 is a judgment, not a measurement, and this review offers no opinion on it beyond agreeing that C4–C6 as a unit cannot be honestly closed without Windows evidence.

**Not done:** no code, test, fixture, build, run, or artifact was modified. No stage was authorized.

---

## 10. Recommendation

1. **Accept the decomposition as the Stage 3 plan.** The staging, the critical path, the rollout gates, and the scope boundaries are correct.
2. **Fix F1 and F2 in the design document** before authorizing C1 and C2. Both are textual, both are contract-accuracy issues, and both would otherwise be frozen into the first batch's acceptance tests.
3. **Record F3, F4, and the D-table reclassification** (D6/D7 closed by citation, D3 bounded by the frozen registry, D1 split into D1a/D1b, D4 folded) in the proposal or the adoption record. Only F4 should be treated as a precondition for C8 authorization; F3 and the decision reclassification are preconditions for C3 and C7 respectively.
4. **Then authorize C1 and C2 as the first batch**, with A2-D9 added to C1's review checklist and both `unavailable` cases added to C2's status matrix.

With those two corrections, this proposal is a sound basis for Stage 3. It does not authorize anything by itself, and this review does not either.

**End — Stage 3 corrective design review.**
