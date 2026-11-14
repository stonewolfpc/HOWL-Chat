# C1 / P17 — OD-6 Handshake Compatibility Taxonomy Assessment (Read-Only)

**Date:** 2026-11-09

**Status:** Assessment only. **OD-6 remains UNRESOLVED.** No register update, disposition, contract amendment, source change, or implementation is made or recommended by this document.

**Baseline:**

- Decision register: `docs/review/DECISION-REGISTER-howl-tts-stage3-c1-p17-evidence-2026-11-09.md`, SHA-256 `4eb22b5e…`, 326 lines; OD-6 entry is unresolved.
- Accepted decisions: OD-2, OD-3, OD-5, OD-8. OD-1 is deferred; OD-4, OD-6, OD-7, and OD-9 are unresolved.
- P17: `01d887c1…`; not approved.
- Contract rev3: `ed0da1a7…`; proposed for independent review, not an accepted amendment.
- Lifecycle decision: `aa136bc9…`.
- Amendment draft A2-D3/A2-D4: not adopted and on conditional hold.

## 1. Exact OD-6 question and current status

The register asks:

> **Is a missing mandatory capability an `incompatible` condition or a `protocol_violation`/`unavailable` condition, and how should helper-declared and client-detected failures be classified?**

The recorded status is **UNRESOLVED**. The register's alternatives are:

- **A.** Capability mismatch = `incompatible`; malformed protocol = `unavailable`.
- **B.** All handshake validation failures = `unavailable` (the current client-detected behavior).
- **C.** All capability/version mismatches = `incompatible`, including missing mandatory capabilities.
- **D.** Keep a distinct internal class and map to one public state only after a contract decision.

The non-binding default is A. The register's A/C boundary is ambiguous: it does not say whether duplicates, invalid descriptors, or misplaced descriptors are "capability mismatch" or "malformed protocol." This assessment resolves that boundary by separating four questions.

### 1.1 Four separate questions

1. **Structural protocol validity.** Is the handshake frame/envelope well formed, and are its descriptors structurally valid enough to understand?
2. **Semantic identity and placement.** Are provider, capability, service, and extension descriptors unique and placed in their correct scope?
3. **Compatibility with mandatory requirements.** Is the helper missing a mandatory capability, or exposing a required capability at an unsupported version?
4. **Public provider-state classification.** Should the observable state be `incompatible` or `unavailable`? Internal failure codes and diagnostics are separate from this public projection.

A failure at layer 1 or 2 does not automatically justify `incompatible` at layer 4. A `protocol_violation` internal code does not, by itself, settle the public state; conversely, `incompatible` is not a generic label for every malformed handshake.

## 2. Source inventory

| Location | Source behavior | Relevance |
|---|---|---|
| `internal/audio/tts/helper/protocol.go:502–544` | `validateHello` returns nine top-level failure classes. | The client-detected handshake validation surface. |
| `protocol.go:503–505` | Empty or overlong helper version returns `protocol_violation`. | Version-shape validation only. |
| `protocol.go:506–508` | Helper limits not matching v1 return `protocol_violation`. | Protocol-limit mismatch. |
| `protocol.go:511–515` | Invalid provider descriptor returns `protocol_violation`. | Descriptor-shape validation. |
| `protocol.go:516–519` | Duplicate provider ID returns `protocol_violation`. | Identity ambiguity. |
| `protocol.go:520–526` | Provider and service capability lists are passed to `validateCapabilities`. | Capability validation. |
| `protocol.go:527–531` | Missing mandatory service capabilities return `protocol_violation` with `Detail["missing"]`. | The mandatory-capability condition. |
| `protocol.go:534–540` | Invalid or duplicate extension descriptor returns `protocol_violation`. | Extension-shape/identity validation. |
| `protocol.go:547–564` | `validateCapabilities` returns `protocol_violation` for an invalid descriptor, unknown core capability, wrong v1 version, or duplicate capability. | Nested capability failure classes. |
| `protocol.go:48–56` | `capabilityVersions` is the v1 registry; `mandatoryServiceCapabilities` is the required service floor. | The source's only capability/version/floor data. |
| `protocol.go:97–101` | `capability` is only `{ID, Version}`; it has no service/provider scope field. | The validator cannot distinguish placement by type. |
| `protocol.go:123–129` | `helloResult` carries separate `Providers` and `Capabilities` arrays, plus `Extensions`. | Scope is represented by array position, not by a validated capability kind. |
| `client.go:248–254` | After handshake failure, `Start` terminates and maps only `protocol_incompatible` and `protocol_minor_unsupported` to `TTSStateIncompatible`; every other failure maps to `TTSStateUnavailable`. | The public projection is a two-code lookup, not a taxonomy. |
| `client.go:240–254` | The comment says the typed failure is preserved so an incompatible helper is not reported as merely absent. | The source intends to preserve compatibility meaning. |
| `client.go:275–294` | `helloHandshake` returns `callOn`, decode, and `validateHello` failures without reclassification. | Helper-declared and client-detected failures enter the same start path. |
| `client.go:434–435`, `:471–481` | `failureFromWire` preserves the helper's wire `Code`. | A helper-declared incompatibility remains that code. |
| `protocol.go:486–500` | `knownErrorCode` includes both `protocol_incompatible` and `protocol_violation`. | Both are internal codes; the public mapping is separate. |
| `types/provider.go:18–36` | The frozen provider states include `incompatible` and `unavailable`. | No new public state is required. |

## 3. Source-verified behavior

### 3.1 The nine top-level failures are not one class

The nine `validateHello` classes are:

1. helper version;
2. limits;
3. provider descriptor;
4. duplicate provider ID;
5. provider capabilities;
6. service capabilities;
7. missing mandatory service capabilities;
8. extension descriptor;
9. duplicate extension.

Every direct return is `protocol_violation`; the delegated capability validator also returns `protocol_violation`. None returns `protocol_incompatible`.

The common internal code does **not** make the conditions semantically identical. The validator uses one code for shape failures, duplicate identities, unknown capability IDs, version mismatches, and missing mandatory capabilities.

### 3.2 Mandatory capability absence is source-visible and currently unavailable

`validateHello` compares service capability IDs against `mandatoryServiceCapabilities` and returns `protocol_violation` with `Detail["missing"]` (`protocol.go:527–531`). `Start` maps that code to `TTSStateUnavailable` (`client.go:252–254`).

The handshake is structurally valid enough to identify the missing IDs. The source nevertheless projects the condition as unavailable. This is the U8 split.

### 3.3 Helper-declared incompatibility is projected differently

A helper can return `protocol_incompatible` or `protocol_minor_unsupported`; `failureFromWire` preserves the code (`client.go:434–435`, `:471–481`), and `Start` maps it to `incompatible` (`client.go:250–254`).

So the source has two reachable projections:

- helper-declared incompatibility → `incompatible`;
- client-detected missing mandatory capability → `unavailable`.

Both are source-verified. The difference is the emitted code, not proof that the underlying compatibility facts differ.

### 3.4 The validator checks shape and some meaning, but not scope

`validateCapabilities` checks:

- ID pattern and length;
- version is a positive integer;
- ID is present in `capabilityVersions`;
- the version equals the v1 registry value;
- no duplicate within that one list.

It does not check whether a capability belongs in a service list versus a provider list. The `capability` type has no scope (`protocol.go:97–101`), and both `helloResult.Capabilities` and `providerDescriptor.Capabilities` use the same validator. A known service-only capability placed in a provider list, or a known provider-only feature placed in the service list, can pass validation if its ID/version are otherwise valid.

This is a material limit on any claim that the source has established incompatibility for "misplaced descriptors."

### 3.5 The validator does not establish supported extension versions

For extensions, `validateHello` checks only that the name matches a namespace pattern, the name is not overlong, the version is at least 1, and the name is not duplicated (`protocol.go:534–540`). It does not consult a supported-extension registry.

An extension with a well-formed but unsupported version (for example, version 2) can pass this validator. The source therefore cannot establish "unsupported required version" for extensions from its current validation.

### 3.6 Duplicate checks are local and do not validate cross-list placement

The source rejects:

- duplicate provider IDs (`protocol.go:516–519`);
- duplicate capabilities within one capability list (`protocol.go:552–563`);
- duplicate extension names (`protocol.go:537–540`).

It does not reject a duplicate capability across the service and provider lists; the shared `howl.tts.synthesize` ID can intentionally appear in both. Cross-list duplication is not, by itself, a source-verified failure.

## 4. The requested categories

### 4.1 Structurally valid handshake missing a mandatory capability

**Source validation:** The capability list is structurally valid, the IDs are known, and `missingCapabilities` identifies a missing required service capability (`protocol.go:527–531`). The source has enough structure and meaning to identify the unmet mandatory requirement.

**Recommended public classification:** `incompatible`.

**Why:** The helper is structurally understandable and has positively failed a mandatory capability requirement. This is not a malformed frame and not a generic provider outage.

### 4.2 Structurally valid handshake with an unknown core capability

**Source validation:** The capability descriptor can be well formed (valid ID pattern, positive version), but its ID is absent from `capabilityVersions`; `validateCapabilities` rejects it (`protocol.go:551–554`). The source can distinguish an unknown ID from a malformed descriptor.

**Recommended public classification:** `incompatible` when the unknown capability is a **required core capability**. The helper is well formed but exposes a capability identity outside the supported v1 core set.

**Limit:** If the unknown ID is an optional or provider-specific feature, the source still rejects it, but the public `incompatible` label may overstate a non-required feature mismatch. That optional/unknown case remains uncertain; it is not automatically a compatibility failure merely because the code rejects it.

### 4.3 Structurally valid handshake with an unsupported required capability version

**Source validation:** A known capability ID with a version other than the v1 registry value is rejected by `validateCapabilities` (`protocol.go:555–558`). The descriptor is structurally valid; the source can identify the exact ID and version.

**Recommended public classification:** `incompatible` for a **required** capability. The helper is well formed but cannot satisfy the required v1 capability version.

**Limit:** The same rule does not transfer automatically to extension versions or helper-version strings. The source does not validate those against a supported-version registry (see §4.6).

### 4.4 Duplicate descriptors

**Source validation:** Duplicate provider IDs, duplicate capabilities within one list, and duplicate extensions are rejected. The descriptor set is ambiguous: the same identity appears more than once, so the client cannot reliably decide which descriptor controls compatibility.

**Recommended public classification:** `unavailable` (protocol/availability failure), not `incompatible`.

**Why:** A duplicate is a structural/identity defect, not evidence that a mandatory capability is absent or unsupported. Calling it `incompatible` would assert a compatibility conclusion the duplicate itself does not establish.

**Limit:** Cross-list duplicate capabilities are not rejected by the source and may be intentional for the shared `synthesize` ID. A duplicate within one list is rejected; a duplicate across lists is not a source-verified failure.

### 4.5 Misplaced descriptors

**Source validation:** The source does not validate service/provider placement. `capability` has no scope field, and both arrays use the same validator. A known service-only capability in a provider list, or a known provider-only feature in the service list, can pass if its ID/version are valid.

**Recommended public classification:** **No public compatibility classification is justified from the current source.** The condition is unvalidated. If a later contract rule rejects misplaced descriptors, the failure should be treated as a structural/semantic protocol failure and projected `unavailable`, not `incompatible`, unless the same handshake independently establishes a missing or unsupported **mandatory** capability.

**Why:** Misplacement is a scope/contract defect. It does not, by itself, establish that a required capability is missing or unsupported.

### 4.6 Invalid semantic descriptors and unsupported versions

**Source validation:** The source rejects invalid helper-version shape, mismatched limits, invalid provider descriptors, invalid capability descriptors, and invalid extension descriptors. Those are structural or contract-shape failures. The handshake cannot be relied on to establish a compatibility conclusion.

**Recommended public classification:** `unavailable` (protocol/availability failure), not `incompatible`.

**Why:** A descriptor that is invalid, or a version string the client cannot validate, does not establish that the helper is missing a mandatory capability. The client should not turn malformed/ambiguous content into a compatibility claim.

**Unsupported extension/helper version:** The source checks extension version only for `>= 1` and helper version only for non-empty/length; it does not check a supported-version registry. An unsupported but well-formed extension/helper version can pass. This category is **unvalidated** and remains uncertain. If a registry is later established, the public classification still needs a separate rule; do not infer `incompatible` from the label alone.

### 4.7 Limits mismatch and provider descriptor failures

**Source validation:** A helper whose limits differ from v1, or whose provider descriptor is structurally invalid, is rejected (`protocol.go:506–515`). These are protocol/contract-shape failures, not mandatory capability failures.

**Recommended public classification:** `unavailable` (protocol/availability failure), not `incompatible`.

**Why:** The helper may be structurally valid JSON but does not satisfy the v1 protocol contract. This is not evidence about a specific provider capability floor. If the owner wants protocol-limit mismatch to be a compatibility category, that requires a separate rule; the source currently maps it to `unavailable`.

## 5. Category summary

| Category | Source validation | Does source establish compatibility conclusion? | Recommended public state |
|---|---|---|---|
| Missing mandatory service capability | Yes; IDs compared with mandatory floor (`protocol.go:527–531`) | Yes: required capability absent | `incompatible` |
| Unknown required core capability | Yes; ID absent from v1 registry (`protocol.go:551–554`) | Yes: required identity unsupported | `incompatible` |
| Unsupported required capability version | Yes; known ID at wrong version (`protocol.go:555–558`) | Yes: required version unsupported | `incompatible` |
| Duplicate provider/capability/extension within a list | Yes; duplicate detected | No: ambiguity, not a missing/unsupported capability | `unavailable` |
| Misplaced service/provider descriptor | No; no scope validation | No: source cannot establish placement | No current classification; if rejected later, `unavailable` |
| Invalid semantic descriptor | Yes; shape/contract invalid | No: content is unreliable | `unavailable` |
| Helper version malformed | Yes; empty/overlong check | No: shape invalid | `unavailable` |
| Limits mismatch | Yes; exact v1 comparison | No: protocol contract mismatch, not capability floor | `unavailable` |
| Extension invalid/duplicate | Yes; shape/duplicate check | No: structure/identity invalid | `unavailable` |
| Unsupported extension version | No supported-version registry | No; well-formed unsupported version passes | Unvalidated; remains uncertain |
| Unknown optional/provider feature | Source rejects unknown ID | Not necessarily; requiredness unclear | Uncertain; do not infer `incompatible` without scope/requirement rule |

This table deliberately does not treat the nine `validateHello` exits as one class.

## 6. Public classification versus internal evidence

The following separation must be preserved:

- **Public provider state:** only `incompatible` or `unavailable` for the cases above. No new public state is added.
- **Internal failure code:** `protocol_violation`, `protocol_incompatible`, `protocol_minor_unsupported`, and other precise codes remain evidence. A code name is not a public projection.
- **Structural protocol validity:** frame/envelope/descriptor well-formedness.
- **Compatibility with mandatory requirements:** missing or unsupported required capability/version.

A structurally malformed handshake can be `protocol_violation` internally and `unavailable` publicly. A structurally valid handshake missing a mandatory capability can be `protocol_violation` internally under the current source, but the recommended public state is `incompatible`. The internal code and the public state need not be identical, but the public state cannot assert more than the evidence establishes.

## 7. Competing interpretations

| Interpretation | Public rule | Source fit | Main consequence |
|---|---|---|---|
| **A. Narrow compatibility** | Only structurally valid, semantically identified mandatory capability/version failures → `incompatible`; malformed/ambiguous/scope-invalid handshakes → `unavailable`. | Fits the source's ability to identify missing/unknown/wrong-version required capabilities; requires a new mapping for client-detected cases. | Best preserves the meaning of `incompatible`; leaves duplicate/misplaced/invalid descriptors as protocol failures. |
| **B. Current code lookup** | Only helper-declared `protocol_incompatible`/`protocol_minor_unsupported` → `incompatible`; every other handshake failure → `unavailable`. | Exactly matches `Start` (`client.go:250–254`). | Preserves source behavior but leaves U8: the same compatibility concept has two public projections. |
| **C. Broad semantic incompatible** | Every well-framed semantic `validateHello` failure, including duplicates, invalid descriptors, and misplaced descriptors → `incompatible`. | Matches parts of A2-D3, but not the source's actual scope validation and not the structural/ambiguity distinction. | Overstates compatibility for duplicate/misplaced/invalid content and asserts facts the source does not validate. |
| **D. Internal taxonomy first** | Record precise internal classes; map to public state only after a contract decision. | Best preserves evidence; the source currently has no such taxonomy. | Avoids a premature public decision but leaves U8 and P17's projection unresolved. |

The register's A and C overlap. The material boundary is **required capability/version incompatibility** versus **structural/ambiguity/scope invalid handshake**, not "capability vs version" alone.

## 8. Recommended OD-6 decision formulation

The following is a recommendation for owner consideration, not a recorded decision.

> **OD-6 — Handshake compatibility taxonomy (proposed):**
>
> 1. A structurally valid `hello` that is semantically sufficient to identify a missing or unsupported **mandatory** requirement — a missing mandatory service capability, an unknown required core capability, or a required capability at an unsupported version — projects `incompatible`.
> 2. A structurally malformed, structurally ambiguous, or scope-invalid handshake projects `unavailable`, not `incompatible`. This includes invalid helper-version shape, mismatched v1 limits, invalid provider/capability/extension descriptors, duplicate identities within a list, and misplaced capability descriptors when a scope rule detects them.
> 3. A helper-declared `protocol_incompatible`/`protocol_minor_unsupported` and a client-detected equivalent mandatory-capability failure project the same public state, `incompatible`, while preserving distinct internal codes and diagnostic evidence.
> 4. Duplicate descriptors are identity/ambiguity failures, not proof of an unsupported mandatory capability. They do not justify `incompatible`.
> 5. Misplaced descriptors are not validated by the current source. The current source cannot justify `incompatible` for them; if a later scope rule rejects them, they project `unavailable` unless the same handshake independently establishes a missing/unsupported mandatory capability.
> 6. Unsupported extension versions and unsupported helper-version strings are not validated against a supported-version registry. They must not be labeled `incompatible` until a registry/rule establishes the required version relation.
> 7. Missing **optional** capabilities are not handshake failures. They are handled at operation time as unsupported capabilities.
> 8. No new public provider state is added. The rule uses frozen `incompatible` and `unavailable`.
> 9. OD-6 does not decide restart/backoff, terminal immutability, evidence-record ownership, or public projection for action-caused exits.
> 10. The current source's all-`protocol_violation` `validateHello` path and its two-code `Start` lookup are not declared conforming. No source change is proposed.

### 8.1 Why this formulation

- It lets `incompatible` mean a specific mandatory capability/version failure, not a generic "handshake failed" label.
- It keeps malformed, duplicate, invalid, and misplaced descriptors in the protocol/availability class where the source cannot establish a compatibility conclusion.
- It preserves the distinction between helper-declared and client-detected evidence while allowing them to share a public state.
- It avoids adding a public state or adopting A2-D4's no-restart rule.
- It leaves the exact internal taxonomy and the contract text to separate decisions.

### 8.2 What it would make non-conforming

If accepted, the current `validateHello` returns for missing mandatory capabilities, unknown required core capabilities, and unsupported required capability versions would need a compatibility classification rather than `protocol_violation`; the current `Start` lookup would need to recognize that classification. The exact internal code (`protocol_incompatible`, a distinct client-detected compatibility code, or another approved representation) is not decided here.

The current source's helper-declared path already projects `incompatible`; it would remain conforming under the public rule, subject to preserving its precise evidence.

### 8.3 What it would not decide

- The exact internal code taxonomy for client-detected semantic failures.
- Whether duplicate provider/extension descriptors should use `protocol_violation` or another internal code.
- Whether a misplaced descriptor rule should exist, and what internal code it should use.
- Whether an unsupported extension/helper version should ever become a compatibility failure.
- Whether A2-D4's no-restart/backoff rule should be adopted.
- Whether U3's handshake-incompatible-then-exit sequence changes the first terminal.
- Whether U8 is closed; it remains OPEN until the owner decides and implementation is separately verified.

## 9. Remaining uncertainty

The following remain unresolved after this assessment:

1. **Misplaced descriptors:** the source has no scope validation. A2-D3 proposes treating them as `protocol_incompatible`; that is not current behavior and would require a scope taxonomy. This assessment recommends `unavailable` only if a future rule detects them, not `incompatible` from the current evidence.
2. **Unknown optional/provider features:** the source rejects any unknown ID. Whether every unknown feature is a mandatory compatibility failure, or only unknown *required* core capabilities, remains undecided.
3. **Extension versions:** the source checks only `>= 1`; it cannot establish an unsupported required extension version. A registry and projection rule are missing.
4. **Helper-version compatibility:** the source checks presence/length, not a supported-version relation. An unsupported helper-version string can pass.
5. **Duplicate semantics:** duplicate provider/extension/capability identities are structural/ambiguity failures in this assessment; whether any duplicate should ever imply incompatibility is not established.
6. **Limits mismatch:** the source treats it as a protocol failure; whether policy should call it incompatibility is not decided.
7. **U8:** the public-projection question would narrow under the recommendation, but U8 remains OPEN until an owner decision and implementation verification.
8. **U3:** handshake-incompatible-then-exit remains governed by OD-3 and is not resolved here.

## 10. Consequences

### 10.1 For U8

U8's public-projection question would have a clear answer for mandatory capability/version failures under the recommendation: equivalent mandatory failures share `incompatible`; structural/ambiguous failures remain `unavailable`. U8 remains OPEN until the owner decides and implementation is separately verified. This assessment does not close it.

### 10.2 For U3 and OD-3

A handshake `incompatible` result followed by a process exit remains a sequence question. The accepted OD-3 first-terminal principle applies; OD-6 does not decide whether the incompatibility or a later crash is the first terminal. U3 remains OPEN.

### 10.3 For P17 and the contract

P17 §7 rule 1 preserves the current split; the recommendation would change that row's taxonomy, but P17 is not approved and no amendment is made here. Contract rev3 is proposed, not accepted; this assessment does not amend it or add a public state.

### 10.4 For OD-5/OD-8

A handshake failure is not an observed process exit. OD-5's accepted `crashed` meaning is unaffected, and a handshake failure must not be used to publish `crashed`. OD-8's causal standard is not engaged. Neither decision is reopened.

### 10.5 For OD-1, OD-4, OD-7, and OD-9

OD-1 remains deferred; no evidence-record ownership is assumed. OD-4, OD-7, and OD-9 are untouched. OD-9 is not addressed; P17 remains unapproved.

## 11. What this assessment does not decide

This assessment does not:

- record an OD-6 disposition, acceptance, rejection, or deferral;
- change OD-1, OD-2, OD-3, OD-4, OD-5, OD-6, OD-7, OD-8, OD-9, or U1–U8;
- close U3 or U8;
- decide the exact internal error-code taxonomy;
- decide whether misplaced descriptors, unsupported extension versions, or optional features should become compatibility failures;
- decide restart/backoff or terminal behavior;
- amend P17 or contract rev3;
- authorize implementation, source changes, tests, builds, or runtime experiments;
- claim that the current source is conforming to the recommendation.

## 12. Verification note

This is a read-only assessment. No existing artifact other than this assessment was created or modified. The decision register remains at SHA-256 `4eb22b5e…`, 326 lines; P17, contract rev3, the lifecycle decision, the OD-5 assessment, and the accepted-decision records remain unchanged. No source, test, configuration, or generated file was modified; no tests, builds, benchmarks, or runtime experiments were run.

**OD-6 remains UNRESOLVED.**
