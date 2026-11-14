# MUSE — INDEPENDENT ARCHITECTURE REVIEW
## Corrected Howl Chat Master Plan — Pre-Adoption Review

**Review ID:** `REVIEW-MUSE-CORRECTED-MASTER-PLAN-2026-11-11`
**Date:** 2026-11-11
**Review type:** Independent, adversarial completeness and architecture review
**Authorization:** Read-only. No edit, implementation, configuration, benchmark, model load, network experiment, acceptance, or plan adoption was performed.
**Plan status at review:** Proposed corrected master plan — not adopted.

---

## 0. Artifact identity verification

Every source was hashed in place. No filename was treated as authority.

| Artifact | SHA-256 | Size / lines | Match |
|---|---|---|---|
| `docs/review/CORRECTED-MASTER-PLAN-HOWL-2026-11-11.md` (review target) | `58e5d45eee1d2139f28f52df718ae93e91d85de5ea480267e03cf3da595351c1` | 107,967 B / 1,072 | Matches its verification record |
| `docs/review/VERIFICATION-HOWL-CORRECTED-MASTER-PLAN-2026-11-11.md` | `cc6489e6b9a82b7b148dd06c9239c506bac2d3430e4d83d99f36e53703693e9d` | 3,662 B / 53 | Read in full |
| `docs/review/REQUIREMENTS-RECOVERY-OWNER-ADDENDUM-HOWL-2026-11-11.md` | `c9f8dc1c08ad5c65954e77c11270ae4d7de18449cd765aebd04cba6a1df7f056` | 222,234 B / 2,418 | Matches plan table |
| `docs/review/REQUIREMENTS-RECONCILIATION-HOWL-2026-11-11.md` | `7ad9d072d43c1243431a91cc304911b5dd8beeecc10923808115a306febda3ab` | 39,612 B / 421 | Matches plan table |
| Existing master plan | `e39cc396fa20adf2656a1b4e650a8bdc6c9a49cb46c3ddf95309ff144229ef8c` | 83,354 B / 703 | Matches |
| Owner decision packet | `6721d770694cda2c7b132fdc15bbd8d19dd6de0502c4f3e5f3ca47ce620ff4f3` | 29,376 B / 385 | Matches plan table |
| Adjudication matrix | `a9fdaeefdc4ea57c28a706cc09e57a77566f3ed6150d2e4e305904de7b6049ed` | 46,357 B / 458 | Matches plan table |
| Muse review (prior) | `af493c0b6118ebede556d4927e8ab93b9378ce46c8298f3bb9e38554f6265e9d` | 33,410 B / 250 | Matches plan table |
| `PROJECT_ISSUES_AND_EXPANSION.md` | `8156318ec8855ddbab1f218927352e59de6d60e23a3c81f11f85a56c61213491` | 132,176 B / 1,931 | Matches plan table |
| `STAGE_STATUS.md` | `33a16fdbca509fc4407bf840deee1dd622a0fe5a1bd272ac4c1ca1179dbe9a62` | 112,949 B / 2,019 | Matches plan table |
| `docs/tier4-avatar.md` | `be04f0253e61c0fc28c245edbd8eb7f3ad753b8bdf152ceb8308cace23be04c9` | — | Matches plan table |
| F2 model-intelligence report | `0f40221f544d04365b784f103d6befc4230b7e2f1699f81d9002265169ac76d1` | — | Matches |
| F2-B integration report | `95e24ccdcdc033c8b050b0f435df8b3eee85dc2e2dd6dab7c3e0bd6d451bee4a` | — | Matches |
| Character-voice investigation | `0236e2300ef4cd704a946035bce5c943a407e33d3bf589214dc6927a8b85ea6e` | — | Matches |

**Limitations.** The reconciliation's F3-B/C acceptance evidence, the exact historical UI planning document, and a formal W4 closure record were not independently reproduced; the plan states the same limitation. No runtime, packet, benchmark, or UI evidence was produced by this review. Findings below that depend on runtime behavior are marked as unproven rather than established.

---

## 1. Recalculation of the reported totals

Recomputed directly from the plan's §12 matrices and the reconciliation §4, not from the verification record.

| Claim | Reported | Recomputed | Result |
|---|---|---|---|
| Forward trace rows | 86 | 86 unique IDs, no duplicates | **Correct** |
| Original reconciliation rows | 66 | 49 `REC-*` + 11 `ODR-*` + 6 `UD-*` = 66 | **Correct** |
| Network rows | 20 | `NET-01`–`NET-20`, exact set | **Correct** |
| Work packages | 28 | `CWP-01`–`CWP-28`, 28 unique | **Correct** |
| Requirement→package links | 95 | 95 forward pairs, 95 reverse, exact set equality; 0 mismatches | **Correct** |
| Capability register | 37 | 37 rows reproduced, identical to reconciliation §5.2 | **Correct** |
| No-work rows | 6 | 6 (`REC-M-10`, `REC-P-05`, `REC-C-02/03/06`, `ODR-01`) | **Correct** |
| Outcome totals | 30 P / 30 M / 13 U / 5 C / 2 E / 6 S | Exactly that; **zero drift** from reconciliation | **Correct** |
| Orphan packages | 0 | 0 | **Correct** |
| Rows without package or no-work marker | 0 | 0 | **Correct** |

Traceability is genuinely generated and exact at the requirement level. This is the strongest part of the plan and it survives adversarial parsing unchanged.

---

## 2. Verdict

**RETURN FOR CORRECTION.**

The plan is structurally sound, unusually honest about evidence, and materially better than the artifact it corrects. Its counts are real, its outcome preservation is exact, its protected-boundary governance holds, and its owner-question discipline is intact.

It is nevertheless **architecturally incomplete in five places that matter before adoption**, and it contains **one dependency inversion and one governance overreach** that are internal contradictions rather than missing detail:

1. Single-character rollover is made dependent on multi-character flush — inverting the plan's own separation.
2. The hot-swap and resource-authority packages are mutually prerequisite, so neither can start under the stated order.
3. The keystone package is blocked by three unsequenced gates with no named unblocking path.
4. The network mission's 19 path rows and 28 verification tests are claimed preserved but exist nowhere at ID level, so the release criterion cannot be evaluated.
5. RAG has a recorded contradiction, a subsystem, and a resource consumer — but no requirement row, no package, and no acceptance criterion.

Plus: durable-state integrity (backup, atomic migration, crash consistency, free-space preflight) is absent from every package, and the plan authorizes investigations despite being unadopted.

None of this requires redesign. All of it requires correction before the plan can be called complete.

---

## 3. Findings

### 3.1 Blocking omissions or errors

---

#### MUSE-01 — Single-character rollover is incorrectly blocked by multi-character flush

- **Severity:** High (blocking)
- **Location:** `CWP-05` prerequisites; contradicts `C-05`, `ODR-11`, `REC-CH-03`/`REC-CH-11`
- **Source evidence:**
  - `CWP-05` (plan §8.2): "**Prerequisites:** `CWP-02`, `CWP-03`, `CWP-04`; canonical chat identity."
  - `CWP-04` is "Character-scoped context isolation and multi-character flush," sourced from `REC-CH-04`, `REC-CH-05`, `REC-CH-06`, `ODR-08`, `ODR-11`.
  - Plan §5.1 `C-05`: "Rollover is single-character context capacity; multi-character flush is character-scoped isolation." `ODR-11` requires them to remain distinct.
- **Problem:** The plan separates the two mechanisms in its contradiction register and then re-couples them in the dependency graph. Single-character rollover — the ordinary, most-used path — cannot begin until character-scoped flush and its own prerequisite `AD-05`/`CWP-15` are done.
- **Why it matters:** This is not a cosmetic edge. It inverts the plan's central distinction, delays the common case behind the hardest case, and creates an implicit claim that rollover needs inference-context isolation, which `ODR-11` explicitly denies. A dependency-ordered roadmap that mis-orders its most common conversation feature is not dependency-ordered.
- **Correction:** Remove `CWP-04` from `CWP-05` prerequisites. `CWP-05` should depend on `CWP-02`, `CWP-03`, and durable chat identity only. If a shared flush primitive is genuinely required, name the primitive as the shared artifact rather than the whole multi-character package.
- **Owner adjudication:** No.

---

#### MUSE-02 — `CWP-08` and `CWP-15` are mutually prerequisite

- **Severity:** High (blocking)
- **Location:** `CWP-08` risks; `CWP-15` prerequisites; `AD-05`
- **Source evidence:**
  - `CWP-08` risks: "Implementation cannot precede `CWP-15`."
  - `CWP-15` prerequisites: "`CWP-08` design, `INV-06`, `AD-05`."
  - Both packages claim `AD-05` ("Single model/resource state authority") as their governing decision.
- **Problem:** §8.1 asserts the order "follows architectural dependency," yet two packages each require the other's output. The plan acknowledges the tension in prose but never breaks the cycle, so there is no defined first artifact.
- **Why it matters:** Hot-swap safety and resource release are the riskiest runtime changes in the plan. A cycle means either package can be declared "blocked on the other," which is exactly how unsafe work gets sequenced by convenience instead of evidence.
- **Correction:** Split the handoff explicitly. Either (a) `CWP-08a` produces the consumer/state contract and `CWP-15` consumes it, then `CWP-08b` implements against the authority; or (b) declare `AD-05` a single joint decision with one owner and one artifact, and list a strict order (`15a design → 08 design → 15b enforce → 08 implement`). State which artifact ends the cycle.
- **Owner adjudication:** No (technical), though the ordering should be visible to the owner because it gates a required product behavior.

---

#### MUSE-03 — The keystone package is blocked by three gates with no unblocking sequence

- **Severity:** High (blocking)
- **Location:** `CWP-01` prerequisites; §3.4; `RG-02`; `INV-03`
- **Source evidence:**
  - `CWP-01` prerequisites: "`INV-03` root/alias evidence; `INV-04` schema investigation; owner answer to `UD-01`; F3-C/D-STORAGE authorization before migration."
  - §3.4: F3-C Phase 0 **OPEN**; D-ROOT **OPEN** (six defects); D-STORAGE **BLOCKED**.
  - `INV-03`: "Resolve F3-C root, alias, startup, wipe, and identity behavior."
- **Problem:** The plan's entire continuity, memory, media, provenance, studio, and presentation chain descends from `CWP-01`, and `CWP-01` cannot proceed on any of three independent gates. The plan names the gates but never names the **decision sequence** that clears them: which owner action comes first, what may proceed before it, and which artifact records resolution. `INV-03` claims to "resolve" F3-C behavior, but F3-C resolution runs through an unauthorized phase-0/R6 path that the plan does not mention.
- **Why it matters:** A roadmap whose first package has three unresolved external blockers and no ordering between them is not executable. It also risks the worst outcome: an implementer resolving root/alias behavior "locally" inside `CWP-01` and thereby making an unauthorized F3-C decision by accident.
- **Correction:** Add an explicit "unblocking decision sequence" subsection stating: (a) the owner answer to `UD-01`; (b) the authorized path for F3-C Phase 0/D-ROOT resolution and who authorizes it; (c) the D-STORAGE condition; and (d) exactly what `CWP-01` may produce before each is satisfied ("schema proposal only — no migration, no identity change"). Name `INV-03`'s output as evidence for the protected path, not a substitute for it.
- **Owner adjudication:** **Yes** — the F3-C/phase-0/R6 path and `UD-01` are owner decisions; the plan must ask for them explicitly rather than absorb them into a package.

---

#### MUSE-04 — The network path register and 28 verification tests are claimed preserved but absent at ID level

- **Severity:** High (blocking)
- **Location:** §10.3, §10.4, `CWP-23`, `CWP-24`; §12
- **Source evidence:**
  - Plan §10.3: "The addendum classifies 14 external path rows, five loopback rows, eight proposed paths, and 28 boundary tests. **The plan preserves those identities**…"
  - Direct search of the plan: **zero** occurrences of `N-01` … `N-14`; **zero** occurrences of `L-01` … `L-05`; `V-NET` appears only as the range `V-NET-01`–`V-NET-28` in three places (§2.1 table, `CWP-23` sources, `CWP-24` sources, §10.4).
  - Addendum §23 defines all 28 tests individually, including `V-NET-23` (no provider connection during local model load), `V-NET-24` (fail-closed on failed/expired/revoked/malformed permission), `V-NET-26` (revocation prevents subsequent requests), `V-NET-27` (Android isolation), `V-NET-25` (browser/download separation), `V-NET-28` (boundary harness, including redirect and oversized-response branches).
  - Addendum §23.7–§23.8 additionally require: no raw `http.DefaultClient`/`net.Dial` bypass, a renderer/webview navigation guard, in-flight cancellation when offline is enabled, `.part` downloads with length and checksum/signature verification, and no mirror substitution. The plan contains **zero** occurrences of `checksum`, `.part`, `webview`, `redirect`, `oversized`.
- **Problem:** The claim "the plan preserves those identities" is false at identifier level. The plan preserves counts and ranges, not rows. Consequently no `V-NET` test has a disposition, no `N-*`/`L-*` path has a registry/disposition entry, and several §23.8 integrity requirements are not represented in any package deliverable.
- **Why it matters:** `NET-15` is the network release criterion, and it is defined as all 28 tests executed with negative proof. Without a per-test map, "`NET-15` → `CWP-24`" cannot show that `V-NET-23` (no provider connection during local model load), `V-NET-26` (revocation), or `V-NET-27` (Android isolation) were ever executed — the three tests most likely to be silently skipped. Likewise, a path row (`N-10`, TCP dial during capability probe; `N-13`, proxy; `N-04`, `git clone`) can drop out of the inventory with no visible loss. This is precisely the renderer-mock-versus-boundary gap the mission warns about, reproduced at planning level.
- **Correction:** Add two tables to §10:
  1. **Path register:** one row per `N-01`–`N-14` and `L-01`–`L-05`, with destination, trigger, data class, current compliance, owning package (`CWP-20`/`21`/`22`/`23`), and disposition (registry/pending/removed/unauthorized). Include `P-01`–`P-08` as unauthorized proposals.
  2. **Verification map:** one row per `V-NET-01`–`V-NET-28`, with boundary (Go dialer/sidecar/webview/renderer), owning package, and state (`not run` for all 28).
  Then add explicit deliverables for cancellation, `.part`/checksum validation, and webview navigation guard to `CWP-20`/`CWP-23`, and the redirect/oversized branches to `CWP-24`.
- **Owner adjudication:** No.

---

#### MUSE-05 — RAG is a subsystem with a recorded contradiction but no requirement, package, or acceptance criterion

- **Severity:** High (blocking)
- **Location:** `S7`, `C-16`, `CWP-02` (dependency/risk only), `CWP-28`
- **Source evidence:**
  - Plan §3.2 `S7`: "Source-confirmed, integration gap. The live prompt path does not demonstrably call RAG. Presence is not proof that replies use it."
  - Plan §5.1 `C-16`: "RAG affects live replies. — **Unproven.**"
  - The reconciliation's 66 rows contain **no** RAG requirement. The plan's 86-row forward matrix therefore has no RAG disposition, no owning package, and no validation criterion.
  - RAG appears only as a `CWP-02` dependency and risk ("RAG may currently be inert") and as a `CWP-28` boundary ("Scope may overlap memory/RAG").
  - Resource inventory in `CWP-15` is generic ("every model consumer"); embeddings appear once, parenthetically, in `CWP-08`'s consumer list.
- **Problem:** A user-visible capability with a source-confirmed integration gap is handled as an incidental dependency of prompt-path work. Three responsibilities are unowned: (a) proving RAG affects replies; (b) the embedding model and vector index as durable state and as a VRAM/resource consumer; (c) document retention/deletion/rebuild semantics, including index rebuild when the embedding model changes or a document is deleted.
- **Why it matters:** "Knowledge documents are used in replies" is a product promise the plan never verifies. An embedding model that competes for VRAM but is outside `CWP-15`'s authority re-creates the exact advisory-admission problem the plan identifies. Uploaded documents are user content with privacy and deletion obligations that no package currently owns.
- **Correction:** Add a RAG requirement row (or an explicit sub-requirement set under `CWP-02` with its own acceptance), and require: a single retrieval injection point proven exactly-once; embedding model identity + index versioning; index rebuild/invalidation on model change and document deletion; document retention/deletion semantics; and embedding workload included in `CWP-15`'s authority. Add `RG` coverage for "RAG demonstrably affects replies."
- **Owner adjudication:** Partially — whether RAG is in committed product scope is an owner question (`UD-` style); the architectural ownership gap is not.

---

### 3.2 Significant architectural weaknesses

---

#### MUSE-06 — Durable-state integrity has no owner: no backup, atomic migration, crash consistency, or free-space preflight

- **Severity:** High (must be fixed before `OA-01` can be granted)
- **Location:** `CWP-01`, `CWP-10`, `RG-02`
- **Source evidence:** `CWP-01` deliverable mentions "migration and rollback design" only. Direct search of the plan returns **zero** occurrences of `backup`, `corrupt`, `free space`, `disk space`, `concurrency`, or `single writer`.
- **Problem:** The plan authorizes a durable chat/media migration over irreplaceable user content while specifying no backup/snapshot step, no atomic staging, no post-migration integrity verification, no free-space preflight, no crash/interruption behavior, and no single-writer guarantee during migration. Rollback design is not crash consistency: a power loss mid-migration has no rollback to run.
- **Why it matters:** The F3-A work already established a staged-wipe precedent (stage both roots, restore on partial failure, report recoverable data) precisely because partial failure is expected. The plan regresses below that standard for a larger and more destructive operation. `CWP-01`'s acceptance says "round-trip design" — a design artifact cannot substitute for a data-safety contract.
- **Correction:** Add a migration-safety contract to `CWP-01` (and media equivalent in `CWP-10`): writable preflight including free space; backup/snapshot with a named restore procedure; atomic staging and commit; integrity verification of migrated records; resumable or cleanly abortable interruption; single-writer enforcement for the duration; and explicit behavior when any step fails. Add it to `RG-02` as evidence, not prose.
- **Owner adjudication:** Authorization for the operation (already `OA-01`); design is engineering's.

---

#### MUSE-07 — Addendum §23.8 integrity and cancellation rules are not carried into any deliverable

- **Severity:** Medium
- **Location:** `CWP-20`, `CWP-23`, `CWP-24`
- **Source evidence:** Addendum §23.8 requires download validation (temporary `.part`, length check, published checksum/signature, no mirror fallback, truncated file never used), in-flight cancellation when offline is enabled, and a renderer/webview navigation guard. Plan: `checksum` 0, `.part` 0, `webview` 0, `redirect` 0, `oversized` 0. Cancellation appears only for image jobs (`CWP-18`) and hot-swap (`CWP-08`).
- **Problem:** These are requirement-level statements in the controlling addendum that survive neither into `NET-01`–`NET-20` (they sit in §23.8 prose) nor into any package deliverable. Only the generic `V-NET` range reference touches them, and that range has no per-test disposition (MUSE-04).
- **Why it matters:** Silent acceptance of a truncated model weight or a substituted mirror is a data-integrity and safety failure, not a UX nicety. Offline mode that blocks new requests but lets in-flight ones complete does not satisfy "local-only."
- **Correction:** Add these as explicit deliverable lines in `CWP-20` (offline cancellation, navigation guard) and `CWP-23` (`.part`, length, checksum/signature, no-mirror, truncation rejection), and as named fixtures in `CWP-24`.
- **Owner adjudication:** No.

---

#### MUSE-08 — Hot-swap preserves the process, not the conversation

- **Severity:** Medium
- **Location:** `CWP-08` validation; `CWP-05` trigger list
- **Source evidence:** `CWP-08` validation covers "spawn failure, readiness timeout, mid-stream switch, cancellation, and rollback" only. `CWP-05` names "model change" among rollover triggers but `CWP-08` never requires transcript canonicalization, template re-rendering, re-tokenization, or context re-fit under the new model.
- **Problem:** `ODR-02` requires hot-swap during active conversation, and the plan's own continuity argument requires that the conversation survive it. The package specifies making the swap *safe* but not making it *coherent*: a switch from a ChatML 32k model to a 4k model with a different template changes message rendering, budget, and possibly special-token handling.
- **Why it matters:** "Changing models preserves meaningful conversational state" is the user-facing promise. A swap that succeeds mechanically but silently truncates or mis-renders history satisfies every existing `CWP-08` criterion while failing the requirement.
- **Correction:** Add a swap-compatibility step to `CWP-08`: canonical transcript → re-render under the new template → re-tokenize → re-fit to the new budget → visible notice of what changed. Add a cross-template/cross-context-size fixture to validation and a corresponding `RG-05` line.
- **Owner adjudication:** No.

---

#### MUSE-09 — Untrusted import parsing has no hardening responsibility

- **Severity:** Medium
- **Location:** `CWP-27`, `CWP-10`, `S17`
- **Source evidence:** `CWP-27` covers "field origin… unknown-field preservation, conflict handling, and round-trip expectations" and "destructive normalization." Search of the plan: zero occurrences of `archive`, `zip`, `malicious`. Security treatment is confined to §10 (network).
- **Problem:** Import is a write path fed by third-party character cards, PNG `chara` metadata, JSON, and possibly archives. No package owns size bounds, decompression limits, path containment for extracted entries, safe handling of embedded metadata, or refusal of in-place rewrite of an external file.
- **Why it matters:** The plan's own `REC-P-06` and `FUT-ST1`/`FUT-ST2` work depends on parsing attacker-controllable input. Path traversal or decompression bombs in import would defeat the F3-A containment work the project already paid for.
- **Correction:** Add a parsing-hardening line to `CWP-27`: bounded read sizes, decompression limits, containment for any extracted path, no in-place rewrite, and malformed-input failure that preserves the source. Add a hostile-input fixture to the `Import/export` regression domain.
- **Owner adjudication:** No.

---

#### MUSE-10 — The plan grants investigation authority it does not have

- **Severity:** Medium (governance)
- **Location:** `CWP-02` owner-authorization line; `CWP-01`, `CWP-06`, `CWP-07`, `CWP-09` authorization lines
- **Source evidence:**
  - `CWP-02`: "**Owner authorization:** **Implementation authorization required.** Investigation is authorized by this plan; no code change is authorized."
  - Also "Investigation is permitted" (`CWP-01`, `CWP-07`, `CWP-09`), "Investigation may proceed" (`CWP-06`).
  - Plan header: "**PROPOSED CORRECTED MASTER PLAN — NOT ADOPTED; NO IMPLEMENTATION AUTHORIZATION**"; §12.4: "No work package may be treated as accepted, closed, or authorized because it appears here."
- **Problem:** An unadopted planning artifact states that it authorizes investigations. That is a self-granted permission and contradicts the document's own status and closing statement. There is also no `OA` gate for investigation authorization, even though investigations such as `CWP-13` require model loads and benchmarks and `CWP-24` requires network experiments.
- **Why it matters:** The review brief's governance rule is that research must not automatically authorize the next stage. A plan that authorizes its own investigations establishes exactly the automatic progression it claims to reject, and the inconsistency will be resolved in practice by whoever wants to start.
- **Correction:** Replace every self-authorizing phrase with "proposed; requires owner authorization to begin." Add `OA-09 — Investigation authorization` covering any investigation that loads a model, runs a benchmark, exercises the network, or inspects user content, and make `CWP-13`/`CWP-24` depend on it explicitly.
- **Owner adjudication:** **Yes** — the owner should explicitly grant or decline investigation authorization as a single decision.

---

#### MUSE-11 — `NET-02` is mis-traced to the Civitai package

- **Severity:** Low
- **Location:** §12.1 `NET-02`; `CWP-21` sources
- **Source evidence:** Reconciliation §4.7 records `NET-02` ("No cloud inference or cloud model provider") as contradicted by DashScope ASR and remappable MeloTTS. `CWP-21` is the Civitai/automatic-lookup package (`NET-03`). `CWP-21` risks even note: "`NET-02` also covers remote audio and cannot be closed by Civitai alone."
- **Problem:** Cloud speech is not cloud model inference, and Civitai metadata lookup is not either. `NET-02`'s actual owning package is `CWP-22` (with `CWP-20` for governance).
- **Why it matters:** The plan's strongest boundary claim — no cloud model provider — is diluted across a metadata-lookup package. A reader could reasonably conclude that fixing Civitai addresses `NET-02`.
- **Correction:** Remove `NET-02` from `CWP-21` sources and from its forward row; keep `CWP-22` as sole owner, with `CWP-20` named as governance dependency.
- **Owner adjudication:** No.

---

#### MUSE-12 — Satisfied constraint rows are given "Action required" dispositions

- **Severity:** Low
- **Location:** §12.1 `REC-M-09`, `NET-20`
- **Source evidence:** `REC-M-09` outcome "Satisfied (constraint preservation)" → disposition "Action required; implementation not authorized" with `CWP-25`. `NET-20` "Satisfied (absence only)" → "Action required" with `CWP-20`.
- **Problem:** The disposition vocabulary conflates "preserve this constraint" with "do work." It is inconsistent with the six no-work rows that use "No work — preserve."
- **Why it matters:** Small, but the plan's value is a clean disposition lexicon; two inconsistent rows undermine the claim that dispositions are machine-checkable.
- **Correction:** Use "Preserve constraint; no product action" for constraint-preservation rows, and keep the using package as a dependency reference rather than a work owner. Keep `NET-20` linked to `CWP-20` only as an absence-maintenance check.
- **Owner adjudication:** No.

---

### 3.3 Improvements that can be deferred

- **MUSE-13 (Low).** `AD-10` and `CWP-25` cite `MOB-1`–`MOB-4` by ID only; §1's remote row does not restate that desktop completeness is mandatory and a companion may never unlock desktop function. Restating the three constraints in §1 would make the acceptance condition readable without the ledger. Not blocking — the IDs are correct and traceable.
- **MUSE-14 (Observation).** `CWP-19` bundles studio foundation, video completion (`REC-C-05`, engine exists), music (`REC-C-04`), and 3D assets (`REC-C-07`, nothing exists) into one package, and `CWP-18` depends on it. The foundation framing is defensible, but video completion is a distinct subsystem with an existing engine and could stall behind modality-priority debate (`UD-05`). Consider naming video as its own sub-scope at design time. No action needed before adoption.
- **MUSE-15 (Observation).** `CWP-16` lists `CWP-10` (durable media identity) as a prerequisite for portraits. Portraits are already partially implemented and do not need durable attachment identity to reach parity between streaming and finalized messages. This over-couples a cheap presentation fix to an expensive migration. Verify the edge is real or drop it at design time.
- **MUSE-16 (Observation).** The plan's §4.4 states runtime claims are "limited to existing records." Once implementation begins, it would help to carry a single "evidence ledger" identifier scheme so `CWP` validation results attach to the seven-state vector rather than to prose. Deferred; no adoption impact.

---

### 3.4 Correctly represented areas (independently checked)

These were checked and found accurate; they should not be re-litigated during correction.

1. **Traceability arithmetic.** 86 rows, 20 network rows, 28 packages, 95 links, 37 register rows, 6 no-work rows, exact forward/reverse set equality, zero outcome drift. Recomputed, not accepted.
2. **Outcome and capability separation.** The six reconciliation outcomes and the seven-state `I/W/S/O/B/V/A` vector remain distinct columns; no capability is promoted; the Flash Attention worked example is unflattering and correct.
3. **Contradiction preservation.** All `C-01`–`C-22` are present with corrected positions, and `NC-01`–`NC-08` are present. Spot-checked against reconciliation §6 and addendum §23; no contradiction was silently resolved.
4. **Governance gates.** `OA-01`–`OA-08` correctly withhold schema, prompt-path, technique, destination, remote-audio, Android, modality, and release authority. The six-stage progression is stated as non-automatic in §1 and §2.4.
5. **Protected boundaries.** F3-C freeze ≠ acceptance (`C-12`), F3-C Phase 0 open, R6 unauthorized, F3-D blocked, D-ROOT open, D-STORAGE blocked, V20 blocked, Gate v2 V1–V17 adopted / V18–V24 proposed, F3-B scoped, TTS/F2-B7 frozen, `W8` verified, `W4` recorded as a closure gap. Consistent with the ledger and freeze records.
6. **Exclusions and parked items.** Desktop avatars excluded (`ODR-01`, `REC-C-06`) with the tier-4 note retained as evidence only; `REC-C-02`/`REC-C-03` parked without promotion. Matches `docs/tier4-avatar.md` and the ledger.
7. **Source claims I verified in code, not just in reports:**
   - `internal/backend/chat/types.go:50-78` — `ChatMeta` has world/scenario/character/lorebook scope arrays and **no model field**; `ChatHistory` carries `id, sender, text, role, timestamp, hasImage, characterId?` and no per-turn metadata.
   - `internal/image/scanner/civitai.go:17-35` — a long-lived browser-UA HTTP client calling `https://civitai.com/api/v1/model-versions/by-hash/<sha256>`; `scanner.go:206-210,260-270` calls `IdentifyFromFile` during scan when the name is generic and no cached name exists. `NC-01`/`C-08` are accurate, not inherited.
   - `frontend/dist/js/settings-actions.js:157-181` — Local-Only Mode replaces `window.fetch` and `XMLHttpRequest.open` only. `NET-08` is accurate.
   - `app-model.go:27-215` — `LoadModel` validates path, admission, and projector pairing **before** stopping the current server, then stops it, waits for the port, validates the binary, and spawns. There is no stream quiescence check and no restart/rollback of the old model if the new spawn fails. `CWP-08`'s premise is accurate.
   - `main.go:53` — `OnShutdown: app.OnShutdown` is registered; `C-01`'s "superseded" position is correct.
8. **No false completion.** No package is marked done, no test passed, no benchmark exists, no estimate appears, and no placeholder or duration language was found.

---

## 4. Unresolved questions requiring owner decisions

These are correctly left open by the plan. The review adds three that the plan does not currently surface as questions.

| ID | Question | Status in plan | Review note |
|---|---|---|---|
| `UD-01` | Global or per-chat model selection | Preserved, unresolved | Also blocks `CWP-08` scope and `CWP-01` schema. Ask first; it gates the keystone. |
| `UD-02` | Android remote scope and security model | Preserved, unresolved | Correctly off by default. |
| `UD-03` | Committed Story Tools scope | Preserved, unresolved | `CWP-28` is a scope brief only. |
| `UD-04` | Narration semantics | Preserved, unresolved | `CWP-17` correctly produces a decision brief, not a feature. |
| `UD-05` | Next studio modality | Preserved, unresolved | Gates `CWP-19` modality work only. |
| `UD-06` | Memory review/retention policy | Preserved, unresolved | Affects memory injection and privacy. |
| — | Retain or remove remote ASR/TTS | Raised in `CWP-22`/`OA-05` | Not assigned a `UD-` ID; recommend promoting it so it is tracked as an owner question, not a package recommendation. |
| — | Authorize investigation as a class (`OA-09`) | **Missing** | Added by MUSE-10. |
| — | F3-C Phase 0 / D-ROOT unblocking path | **Missing** | Added by MUSE-03. |
| — | Is RAG a committed product capability? | **Missing** | Added by MUSE-05; affects whether RAG gets a requirement row or an explicit exclusion. |
| — | Any external destination authorization | Preserved as unauthorized | Correct; `P-01`–`P-08` remain proposals. |

---

## 5. What this review did not do

- It did not edit the corrected plan, the verification record, the recovery addendum, the reconciliation, the existing master plan, the ledger, the stage record, or any protected F3-C artifact.
- It did not implement, configure, migrate, load a model, run a benchmark, exercise the UI, capture packets, or perform any network experiment.
- It did not mark any requirement accepted, any capability operational, or any release gate passed.
- It did not resolve any owner question, including the three it recommends surfacing.
- It did not treat the verification record's conclusions as evidence; it recomputed the counts and re-derived the source claims independently.

---

## 6. Required corrections before a further adoption review

1. Remove `CWP-04` from `CWP-05` prerequisites (MUSE-01).
2. Break the `CWP-08`/`CWP-15` cycle with a named first artifact (MUSE-02).
3. Add the unblocking decision sequence for `CWP-01`, naming the owner actions (MUSE-03).
4. Add the path register (`N-01`–`N-14`, `L-01`–`L-05`, `P-01`–`P-08`) and the per-`V-NET` disposition map; add cancellation, download validation, and navigation-guard deliverables (MUSE-04, MUSE-07).
5. Give RAG a requirement row or explicit exclusion, an owner, and acceptance criteria; include embeddings in resource authority (MUSE-05).
6. Add the durable-state migration-safety contract (MUSE-06).
7. Add swap-compatibility continuity to `CWP-08` (MUSE-08).
8. Add import parsing hardening to `CWP-27` (MUSE-09).
9. Replace self-authorizing investigation language and add `OA-09` (MUSE-10).
10. Correct the `NET-02` trace and the two Satisfied-row dispositions (MUSE-11, MUSE-12).

After these corrections, the plan can return for a focused re-review of the changed sections rather than a full re-read. Nothing in this review authorizes implementation, migration, enablement, or adoption.
