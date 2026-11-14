# Decision Report — C-1 Capability Policy

**Date:** 2026-11-06  
**Scope:** Read-only investigation and recommendation for C-1 only  
**Governing contract:** `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` §17  
**Status:** DECISION ANALYSIS. NO CONTRACT AMENDMENT, CODE CHANGE, TEST CHANGE, STAGE 3 WORK, COMMIT, OR PUSH.  
**Stage boundary:** Stage 2 and its correction remain accepted and closed. Stage 3 is not authorized.

---

## 1. Bottom line

The frozen contract **does not decide C-1**. It contains:

- a hello example with five service capabilities;
- no normative statement that those five are mandatory;
- a per-request `unsupported_capability` remedy for an unavailable capability;
- a statement that a false capability “is not silently omitted”; and
- a wire `Capability` DTO with no `supported` field, so absence is the only wire encoding of a false capability.

That is not enough to declare either Policy A or Policy B contract-conformant by derivation.

Recommendation: **Policy B, with a narrow functional floor.** Require `howl.tts.load` and `howl.tts.synthesize`; make `howl.tts.catalog`, `howl.tts.unload`, and `howl.tts.cancel` optional. Missing optional operations return `unsupported_capability` without sending a frame. Missing the two required service capabilities is `protocol_incompatible` at handshake, not `protocol_violation`.

This preserves the provider-neutral boundary, keeps fail-closed behavior for the operations that make a TTS helper useful, and leaves process termination as the lifecycle fallback for unload/cancel. The current implementation remains Policy A; it must not be changed until the owner approves the contract wording in §6.

**Owner decision required:** choose A or B-with-floor, and approve the wire meaning of absence. This report does not resolve the ambiguity by assumption.

---

## 2. Normative record and exact references

All line references are to `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` unless stated otherwise.

### 2.1 Capability registry and versions

| Subject | Normative text | Consequence |
|---|---|---|
| Core IDs | §17.6.1, lines 1683–1699 | Nine v1 IDs: five service operations and four provider features (`speed`, `language`, `ssml`, `streaming`). |
| Wire capability shape | §17.5.1, lines 1429–1439 | `Capability` is exactly `{id, version}`. No other field is permitted. Provider `capabilities` may be an empty array. |
| Version rule | §17.6.1, lines 1685–1695; §17.9, lines 1880–1901 | Each registry ID has version 1. The codec rejects wrong versions and unknown IDs. |
| Service/provider split | §17.6.1, lines 1697–1699 | `hello` advertises **service** capabilities; provider descriptors and `catalog` advertise **provider** capabilities. |
| False capability | §17.6.1, lines 1698–1699 | The text says a capability with `Supported:false` is false, not silently omitted. The wire DTO has no `supported` field. |

### 2.2 Hello negotiation

| Subject | Normative text | Consequence |
|---|---|---|
| Hello result shape | §17.3.2, lines 1244–1281 | The example advertises exactly the five service IDs at lines 1265–1271. The surrounding text describes the example but states no mandatory minimum. |
| Provider descriptors | §17.3.2, lines 1277–1279 | `hello` describes what the binary can expose; `catalog` describes what the selected provider/assets actually support. |
| Handshake validation | §17.11.3, line 2056 | The supervisor validates “providers, capabilities, extensions, and limits.” It does not name a required service set. |
| Compatibility rules | §17.3.4, lines 1300–1317 | Major/minor, provider identity, catalog schema, and extensions are defined. Service capability presence is not. |

### 2.3 Unsupported capability behavior

| Subject | Normative text | Consequence |
|---|---|---|
| Strict unsupported matrix | §17.6.3, lines 1719–1732 | “Capability false” maps to `unsupported_capability`. |
| Error meaning | §4.6, lines 383–407; §17.9, line 1900 | `unsupported_capability` is the per-request error for an unavailable **optional** capability. |
| Error registry | §17.9, lines 1874–1907 | `protocol_violation` is reserved for unmatched IDs, duplicate terminals, and publication after terminal; `protocol_incompatible` is for no common major or major mismatch. |

### 2.4 Service operations

| Subject | Normative text | Consequence |
|---|---|---|
| Operation table | §17.4.1, lines 1323–1335 | All operations require “Handshake complete.” `catalog` needs no load; `load` requires no other non-control operation active; `synthesize` requires the matching provider loaded and idle; `cancel` targets an active synthesis. |
| Reader/worker | §17.4.2, lines 1337–1354 | `catalog`, `load`, `unload`, and `synthesize` run on one worker; `health`, `cancel`, and `shutdown` are handled by the reader during synthesis. |
| `load`/`unload` | §17.4.4, lines 1368–1379 | `unload` is idempotent, releases resources when idle, returns `busy` during synthesis, and **never terminates the process to force an unload**. |
| `cancel` | §17.4.6, lines 1389–1395 | Cancel is a cooperative attempt. Acknowledgement is not proof that native work stopped. The supervisor kills after the grace if no terminal arrives. |
| `shutdown` | §17.4.7, lines 1397–1402 | Shutdown stops new work, cancels active synthesis, acknowledges, releases resources, and exits; the supervisor kills after the grace. |

### 2.5 Lifecycle, cancellation, and readiness

| Subject | Normative text | Consequence |
|---|---|---|
| Discovery | §17.11.1, lines 2018–2028 | Helper resolution is explicit and does not search `PATH`. |
| Ownership | §17.11.2, lines 2030–2047 | The main supervisor owns the job, pipes, and the sole `Wait()` observer. |
| Cancellation path | §17.11.4, lines 2064–2080 | The parent records cancellation, sends `cancel`, waits the grace, then kills and observes exit before deleting output. |
| Acceptance barrier | §17.11.6, lines 2100–2138 | Cancellation before `accepted` wins; a late success cannot revive a cancelled request. |
| State vocabulary | §17.7.1, lines 1736–1760 | `Available` means helper/protocol/provider/assets are present; `Ready` means loaded and able to synthesize. |
| Transitions | §17.7.2, lines 1762–1790 | `unload` moves `ready`/`degraded` to `unloaded`; process exit moves to `absent`/`crashed`. |
| Provider surface | §5.2, lines 454–478 | The provider has `Discover`, `Load`, `Synthesize`, `Unload`, and `Close`; the helper owns process shutdown. |
| Provider instances | §5.5, lines 502–518 | One loaded provider per process; at most two helper processes; no silent replacement. |

### 2.6 The controlling ambiguity

The contract cannot yield a single answer from these clauses:

1. The §17.3.2 example shows five service capabilities, but an example is not a mandatory minimum.
2. §17.6.1 says `hello` “advertises” service capabilities and that false capabilities are not omitted, but the wire DTO (§17.5.1) has no `supported` field. The only wire encoding of a false capability is absence.
3. §17.6.3 says a false capability produces a per-request `unsupported_capability`; that is Policy B behavior, not Policy A handshake rejection.
4. The app-level `TTSCapability` has `Supported bool` (§3.2, lines 192–196), so the “not silently omitted” sentence may be describing the app projection rather than the wire. That is not stated.
5. No clause says missing service capabilities make a hello malformed, incompatible, or a `protocol_violation`.

This is the exact unresolved decision: **does omission from `hello.result.capabilities` mean “not supported” (B), or is every service capability mandatory despite the per-request remedy (A)?**

---

## 3. Policy comparison

### 3.1 Definitions used here

- **Policy A:** all five service capabilities are mandatory. A hello missing any one fails negotiation.
- **Policy B:** a helper may omit a registered service capability. The operation requiring it returns `unsupported_capability`; other operations remain usable.

For the recommendation, **B-floor** means B with `load` and `synthesize` required. It is not pure “zero-capability helper” acceptance.

### 3.2 Comparison table

| Concern | Policy A — mandatory five | Policy B — degrade per capability | Assessment |
|---|---|---|---|
| **Negotiation** | One fail-fast check; every accepted helper has the full surface. | Accepts a partial helper; each operation must be gated. | A is simpler; B is closer to the frozen per-request remedy. |
| **Future provider with no independent unload** | Must implement `unload` or be rejected. A no-op unload would be false advertising; process termination is not a legal unload under §17.4.4. | Omits `unload`; process termination remains a supervisor lifecycle action. | **B** avoids forcing a provider to claim a capability it does not have. |
| **Provider with no native cancellation** | Must implement `cancel`, even if only to acknowledge; failure to stop is still handled by the kill fallback. | Omits `cancel`; parent skips the frame and uses the normal termination fallback. | **B** is honest; A is not unsafe, but it adds an implementation burden that does not improve the fallback. |
| **Provider with no catalog** | Rejected even if it can synthesize with a configured voice. | `catalog` returns `unsupported_capability`; static/fallback metadata must be explicit. | **B** supports narrow providers; A forces catalog support for every provider. |
| **Provider with no load** | Rejected. | A helper with `synthesize` but no `load` is not usable under the frozen state machine, so B-floor still rejects it. | The floor is necessary; `synthesize` requires a loaded provider (§17.4.1, 1330; §17.4.4, 1379). |
| **Provider replacement/backward compatibility** | A newer helper that drops any service operation is incompatible; older helpers with less surface are rejected. | A limited helper remains usable for what it advertises; unsupported operations are typed failures. | **B** is more compatible with independent helper updates (§9, lines 821–836). |
| **Error classification** | Requires a handshake error for missing capabilities. Current code uses `protocol_violation`, which the frozen registry does not define for this condition. | Uses the existing `unsupported_capability` per request. Missing B-floor is `protocol_incompatible` at handshake. | **B** reuses the defined per-request error and avoids overloading `protocol_violation`. |
| **User-visible status** | `incompatible`/unavailable for any missing service. | `Available`/`Ready` plus a truthful capability list; unsupported operation returns a typed error. | B requires per-capability status discipline; A is simpler but less truthful for partial helpers. |
| **Security** | No material gain. Job ownership, one `Wait()`, path validation, output validation, and frame limits are independent of this policy. | No material loss if unsupported operations fail closed and no fallback provider is selected silently. | **Tie** on security; do not justify A as a security requirement. |
| **Lifecycle** | All helpers can unload/cancel through the protocol; process kill remains the backstop. | `unload`/`cancel` may be absent; supervisor termination remains the backstop. | **B** needs explicit state and cancellation wording; A needs a fake or forced unload. |
| **Implementation complexity** | One mandatory-set check; fewer paths. | Per-operation availability, status projection, optional unload/cancel paths, more tests. | A is cheaper; B is more provider-neutral. |
| **Contract fit** | Requires an explicit amendment to §17.6.1; the current text is silent. | Also requires an amendment, because absence-as-false and the mandatory floor are not explicit. | **Neither is currently derivable.** |

### 3.3 What mandatory capabilities actually protect

They protect one useful property: **a successful handshake means the helper can perform the complete v1 operation set.** That is a fail-fast integration contract. It does not protect process ownership, output-path safety, cancellation bounds, or orphan prevention; those are enforced by §17.11 and are independent of capability negotiation. It also does not protect audio correctness.

The cost is real: a provider that can synthesize but cannot independently release its model is rejected even though the supervisor can terminate the process. A provider that cannot acknowledge cancellation is rejected even though the kill fallback already bounds the request. That is a policy choice, not a security necessity.

### 3.4 Why not pure Policy B

Pure B permits a helper with no `load` or `synthesize` to complete the handshake. That helper is not a TTS provider, yet the §17.7.1 definition of `Available` would not distinguish it from a usable one without additional capability-aware status rules. It also makes the initial state `unloaded` and `Ready=false` for an object that can never become ready.

The narrow floor resolves that without returning to mandatory-five:

- `synthesize` is the purpose of the helper.
- `load` is the only frozen path to a loaded provider before synthesis.
- `catalog` can degrade to static/fallback metadata.
- `unload` can degrade to supervisor process termination.
- `cancel` can degrade to the existing grace/kill fallback.

This is the smallest policy that keeps “accepted TTS helper” meaningful while allowing provider variation.

---

## 4. Consequences in detail

### 4.1 Future providers with different capabilities

Policy A would make the helper protocol a full-surface contract: every provider must expose catalog, load, unload, synthesize, and cancel even when the provider's native runtime has no independent unload or cancellation. That tends to produce no-op implementations, which are worse than absence because they make status and failure diagnosis untrue.

Policy B lets a provider advertise only what it can do. The app must show that truth: a capability absent from the wire is `Supported:false` in `TTSProviderStatus.Capabilities`; an operation returns `unsupported_capability`; it is not treated as a successful no-op. This is consistent with the additive inspector in §3.2 and the provider-neutral boundary in §2.2.

### 4.2 Providers that cannot unload independently

`unload` is not equivalent to “stop the process.” §17.4.4 says it releases model/runtime resources when idle and **never terminates the process to force an unload** (lines 1377–1378). A provider that cannot release its model must therefore omit `unload`; advertising it and doing nothing would violate the frozen semantics.

Under B, the supervisor may still stop the process for idle close, provider replacement, or shutdown. That is a lifecycle transition, not a successful unload:

- it must not report `unloaded` merely because the process exited;
- it must transition through `stopping` to `absent` (or `crashed` on failure), preserving the existing `Wait()` authority;
- a later request starts a new process and reloads.

This is a contract change beyond the hello array: §17.7.2 needs an explicit “process stop is not unload” note.

### 4.3 Providers that cannot cancel an in-progress native operation

`cancel` is already defined as a cooperative attempt, not proof of interruption (lines 1391–1395). The supervisor's 2-second grace and `KillTree` fallback bound a hung native call. If `cancel` is absent:

- the parent must not send a `cancel` frame;
- it records `cancel_requested` before any termination;
- it may wait the normal grace for a terminal, then kills and observes `Wait()` before cleanup;
- the request still resolves as `cancelled` or `timeout` under §17.11.6.

This preserves the cancellation guarantee. If `cancel` is present, the existing acknowledgement path remains. The absence of `cancel` must not cause the parent to hang or to report a successful cooperative stop.

### 4.4 Provider replacement and backward compatibility

The versioning rules permit independent helper updates within a major, with optional additions advertised (§9, lines 821–836). Policy A makes a helper that drops any one of the five operations incompatible with the same major. Policy B lets the app continue to use the remaining operations, but the app must not silently route synthesis to a different provider: provider selection and fallback behavior remain explicit (O-3 remains open; §17.15, line 2325).

This is the main compatibility advantage of B. It is also the main operational risk: partial providers need truthful status, or users will see repeated typed failures instead of a clear “unavailable” state.

### 4.5 Error classification and user-visible status

Policy B should use:

| Situation | Error | Status |
|---|---|---|
| Optional service operation omitted | `unsupported_capability` at request time; no frame sent | `Available`/`Ready` may remain true for other operations; capability entry is false |
| Required `load` or `synthesize` omitted | `protocol_incompatible` at hello, with bounded `detail.missing` | `incompatible`, `Available=false`, no automatic retry |
| Provider feature omitted (`speed`, `language`, `ssml`, `streaming`) | `unsupported_capability` before synthesis | Provider remains loaded and usable for other requests |
| `unload` omitted | `unsupported_capability`; no frame | Process may be stopped by supervisor; not reported as unloaded |
| `cancel` omitted | `unsupported_capability` from the public method; internal cancellation path skips the frame | Request still ends `cancelled`/`timeout` after the bounded fallback |

Using `protocol_violation` for a missing service capability is not supported by the frozen registry. §17.9 line 1878 defines `protocol_violation` as an unmatched ID, duplicate terminal, or publication after terminal. A capability omission is a compatibility/policy condition, not a malformed frame. The current code uses `protocol_violation` at `internal/audio/tts/helper/protocol.go:527-531`; that is an implementation choice requiring contract text, not a settled classification.

`unsupported_capability` currently maps to `ErrCodeSynthesisFailed` (line 1900). That is the established compatibility projection and is acceptable for per-request optional features. The precise internal code and detail must remain visible in status/diagnostics per §17.9 lines 1863–1872.

### 4.6 Security and lifecycle guarantees

Neither policy changes the security boundary:

- the helper remains a same-user, job-owned process;
- the parent owns pipes and the single `Wait()` observer;
- output paths, frame sizes, text size, and process count remain validated;
- cancellation still ends in termination when cooperative cancellation fails;
- no missing capability authorizes a silent fallback, a new process, or a path outside the app-owned cache.

Policy B must add one fail-closed rule: an unsupported operation never becomes a no-op success, and an absent capability never causes the parent to send the operation anyway. That is the same gate principle already established for provider features; it is not a new security model.

---

## 5. Recommendation

**Adopt Policy B with the narrow floor: require `load` and `synthesize`; make `catalog`, `unload`, and `cancel` optional.**

Reasons:

1. It follows the contract's own per-request `unsupported_capability` remedy instead of inventing a handshake-wide rule.
2. It preserves provider neutrality for providers that cannot independently unload or cancel.
3. It keeps the lifecycle guarantee via the already-frozen supervisor kill/`Wait()` path.
4. The floor prevents a zero-capability helper from being reported as an available TTS provider.
5. It requires explicit, testable status and error semantics rather than silent degradation.

**Do not choose A merely because it is fail-fast.** A is defensible if the owner wants every accepted helper to have the complete surface, but it must be stated as a contract amendment, it forces honest-but-limited providers out, and its handshake error must be reclassified from `protocol_violation` to `protocol_incompatible` (or a new registry code). The initial provider may implement all five; that does not justify making five the protocol law for every future provider.

**No code change is recommended by this report.** The current implementation remains Policy A until the owner decides. When authorized, the changes would be limited to the capability policy and its tests, not the accepted process implementation.

---

## 6. Proposed contract wording

The following is a **proposal only**. It does not amend the architecture document.

### 6.1 Replace §17.6.1, lines 1697–1699

> `hello.result.capabilities` advertises the service capabilities the helper implements. Service capabilities use the five core service IDs in the table above. Presence means the helper implements that operation at v1; absence means it does not. A conforming TTS helper MUST advertise `howl.tts.load` and `howl.tts.synthesize`. A hello missing either is `protocol_incompatible` with a bounded `detail.missing` naming the absent IDs, and the provider is not available. `howl.tts.catalog`, `howl.tts.unload`, and `howl.tts.cancel` are optional. A request for an absent optional service capability returns `unsupported_capability` and MUST NOT write an operation frame.
>
> `ProviderDescriptor.capabilities` and `catalog` advertise provider capabilities. A provider capability that is not supported is absent; the wire `Capability` DTO has no `supported` field. The application projection sets `TTSCapability.Supported=false` for an absent provider capability. A request requiring an absent provider capability returns `unsupported_capability` before synthesis. A capability is never silently treated as supported, and no no-op implementation may substitute for an absent capability.

### 6.2 Add to §17.3.2, after line 1281

> The five service IDs in the example are the complete v1 service vocabulary, not a requirement that every helper advertise all five. The required floor and the optional set are defined in §17.6.1. Provider descriptors may contain only provider capabilities; the service array may contain only service capabilities.

### 6.3 Add to §17.4.4, after line 1378

> If the helper does not advertise `howl.tts.unload`, `unload` returns `unsupported_capability` and no frame is written. The supervisor may stop the process for idle close, replacement, or shutdown, but that is a lifecycle action, not a successful unload. A process stop MUST NOT be reported as the `unloaded` provider state; the transition follows the process lifecycle (`stopping` to `absent`/`crashed`). A provider that advertises `unload` MUST actually release the loaded resources; a no-op is non-conformant.

### 6.4 Amend §17.11.4, lines 2069–2073

> 4. On context cancellation, it first records `cancel_requested` for that request. If the helper advertised `howl.tts.cancel`, it sends `cancel` with its own ID and the target request ID. If it did not, it skips the cancel frame; absence of `cancel` is not an error in the cancellation path.
> 5. If a cancel frame was sent, the helper reader acknowledges and signals the provider context; native code may not stop in place.
> 6. If no terminal response arrives within the cancel grace, the supervisor calls `proc.KillTree`, waits for the process to exit, and returns `cancelled` or `timeout`.

### 6.5 Amend §17.9, `protocol_incompatible` row (line 1881)

> `protocol_incompatible` | wire/app | false | No common major, major mismatch, or a hello missing a required service capability | `ErrCodeModelLoadFailed`

The `detail.missing` value is bounded to the absent required IDs. `protocol_violation` remains reserved for unmatched IDs, duplicate terminals, and publication after terminal; a capability omission is not a protocol violation.

### 6.6 Add to §17.7.1 or §17.7.4

> A helper that advertises only a subset of optional service capabilities remains `Available` for the capabilities it advertises; `Ready` still means loaded and able to synthesize. The status capability list marks absent IDs as unsupported. `unload` and `cancel` absence must not be represented as successful operations, and the process lifecycle remains authoritative for termination.

### 6.7 No wire-shape change

Do not add a `supported` field to the wire `Capability` DTO. That would break the exact v1 shape and the frozen fixtures. Absence is the wire encoding of an unsupported capability; `Supported` remains an application projection.

---

## 7. Exact owner decision required

The owner must choose one:

**Option A — Mandatory five**
- Approve §17.6.1 wording making all five service capabilities mandatory.
- Approve `protocol_incompatible` (recommended) or a new registry error for a missing capability; do not keep `protocol_violation`.
- Accept that honest providers without independent unload/cancel are non-conformant.
- Update the current implementation and its tests before Stage 3.

**Option B — B-floor (recommended)**
- Approve `load` + `synthesize` as the mandatory floor and `catalog`/`unload`/`cancel` as optional.
- Approve absence-as-false on the wire and `Supported=false` for the application projection.
- Approve the unload-as-lifecycle and cancel-skip fallback wording.
- Add focused tests for optional capability omission, required-floor rejection, no-frame gates, unload state, and cancellation without `cancel`.
- Keep `protocol_violation` reserved for malformed protocol.

The contract decision belongs in the architecture document before Stage 3, not only in code or tests. The current implementation is Policy A and is not being changed by this report.

---

## 8. Read-only scope confirmation

- No architecture document was edited.
- No code, test, fixture, golden file, or accepted Stage 2 artifact was modified.
- No `Unload` or `Cancel` caller was wired.
- No native provider, ONNX Runtime, model asset, download, packaging, commit, or push occurred.
- Stage 2 and its correction remain accepted and closed.
- Stage 3 remains unauthorized.

**End — decision report.**
