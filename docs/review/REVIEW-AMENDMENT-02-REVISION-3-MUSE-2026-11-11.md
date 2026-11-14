# Muse review — Amendment 02, Revision 3

**Date:** 2026-11-11
**Reviewer:** Muse
**Review type:** Independent requirement-capture and planning review
**Package:** `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md`, revision 3
**Review request:** `docs/review/AMENDMENT-02-MUSE-REVIEW-REQUEST-2026-11-11.md`
**Status:** Review complete; no document or source file was modified

## 1. Verdict

**PASS**

The amendment accurately captures the owner-defined requirements, keeps the known ambiguities open, preserves the frozen inputs, and does not imply implementation or runtime authorization. It is ready for the affected research planning (R-07–R-11) after the owner reviews this result.

This verdict is a requirement-capture result only. It is not architecture approval, not an implementation gate, not acceptance, and not permission to begin research before the owner acts on it.

Two advisory findings are recorded in §6. Neither blocks research planning; both must be closed before architecture approval.

## 2. Artifacts inspected and integrity

| Artifact | Expected | Recomputed | Result |
|---|---|---|---|
| Amendment 02, revision 3 | `fa766f76779e2517b31bd55da993f872525a52941b348fc8f78bbeb1d9382c44` | same | Match |
| Amendment 02 sidecar | `sha256sum -c` | `OK` | Verified |
| Integration and acceptance plan | `1e92a54434e8e3269015c47374ad57bd6ae755e64e27a07954c083ab952510cf` | same | Match |
| Integration plan sidecar | `sha256sum -c` | `OK` | Verified |
| Amendment 01 | `f3a0f8b9c5d19129db85d559f1880f21451cc7aa70596d37a408952c85c3f4ea` | same | Match |
| Amendment 01 sidecar | `sha256sum -c` | `OK` | Verified |
| Corrected master plan, Revision 3 | `d642039cd5ddcdc480001666f7a18ad36e22944f5540349056af1eee050ee874` | same | Match |
| Requirements reconciliation | `7ad9d072d43c1243431a91cc304911b5dd8beeecc10923808115a306febda3ab` | same | Match |
| Recovery owner addendum | `c9f8dc1c08ad5c65954e77c11270ae4d7de18449cd765aebd04cba6a1df7f056` | same | Match |
| Owner-decision addendum | `d228c4ac8a1acfabc7164ab5b277e9c20e8e2c33522554ced4896ec02b41e5ac` | same | Match |
| Memory design | `2e8ec935843bfd6b110feb16a626ce856cff60dd9ec1f570df0d21b21ef57650` | same | Match |
| Hardware design | `9348fa0dc297db356f8391f66318e98f8310b3ee6900b7116137ac23f852b0c3` | same | Match |
| RAG design | `2a80de38267131c7c4557729b845ac0bd43b3e345ad395020b93deb9e738e7e6` | same | Match |
| Android research | `c25e45b511b883c26d44ac8d3b61d66d6397a12a13f120f035a73f75474e33a2` | same | Match |
| Reconciled plan | `58e33ed918e22628854a299380b9178b3cf61481b3a82cdebda94c459ee10040` | same | Match |
| Design package verification | `0871bd4ace4dc48567f2d338310a464aae13698719ab61dcefaae0244cbdd1d1` | same | Match |
| Prior independent review | `f1756fb5e6c21b8dce5ed30e677426f581ecfc732463a30eb4a630d34d8d5eee` | same | Match |
| Prior independent re-review | `88b12137baca5b27bb733de911c5fb69dc0fb21395123fe0e21dcbcc0fb413af` | same | Match |

Amendment size and line count also match the declared `61,111 B / 1,122 lines`.

Observed: every sidecar check was run from the repository root, because the sidecars record root-relative paths. A check run from `docs/architecture` reports "Failed open or read" for reasons unrelated to artifact integrity.

## 3. Review-question answers

| ID | Question | Answer | Evidence |
|---|---|---|---|
| `RG-01` | Are all owner-defined requirements captured? | **Yes** | §1, §5–§10, §16; see §4 below |
| `RG-02` | Any requirement weakened, reinterpreted, or turned into an implementation assumption? | **No** | §9.2, §10.2, §18, §15 |
| `RG-03` | Is Public/Private isolation covered across all channels? | **Yes**, with advisory finding `M2-01` | §10.3, §10.4, §10.5, `AC-C49` |
| `RG-04` | Are transcript and model input distinct and independently inspectable? | **Yes** | §3, §4.1, §5.3, `AC-C04`, `AC-C07` |
| `RG-05` | Are personal lorebooks, automatic lore, and memory distinct? | **Yes** | §4.2, §9.1, §9.2, `AC-C37`, `RQ-11.5` |
| `RG-06` | Is durable media identity/reuse/persistence/deletion/board association covered? | **Yes** | §5.2, §5.5, §8.3, §9.6, `AC-C01`–`AC-C11`, `AC-C41` |
| `RG-07` | Do Chat image generation and multi-character runs use Amendment-01 orchestration? | **Yes** | §3.1, §8.2, §8.4, `AC-C21` |
| `RG-08` | Is research scope adequate, evidence-led, explicit about the unestablished? | **Yes** | §12, §12.2, §13, `SRC-08` |
| `RG-09` | Are dependencies and acceptance criteria testable, including constrained systems? | **Yes**, via Amendment 01 for constrained-system evidence | §12.1, §14, §16, §13 |
| `RG-10` | Are original plan, Amendment 01, frozen documents, and digests intact? | **Yes** | §2 above; all sidecars `OK` |
| `RG-11` | Is implementation or runtime authorization implied? | **No** | §20, §21.1, §14, §24, §16 preamble |
| `RG-12` | Does the package support a clear verdict without implying architecture approval? | **Yes** | §24; review request §5–§6 |
| `RG-13` | Is `OD-01` preserved as settled without reopening the language debate? | **Yes** | §1.1, §5.5, §11, §14, §19, §21.1, `SRC-09`, `CH-07` |

## 4. Requirement fidelity (RG-01, RG-02)

Each required item is present, specific, and not softened:

| Requirement | Location | Assessment |
|---|---|---|
| Canonical visible history independent of model context | §3, §4, §5.1 | Explicit three-system separation; visible transcript is a projection of canonical events, not of the model context. |
| Durable media identity and lifecycle | §5.2, §5.5 | Seven named states replace `hasImage`; separate storage/visible/model-input identities. |
| Protected initial message, absolute top, editable, not deletable | §6, §6.1–§6.3, `AC-C08`–`AC-C10` | Protection survives editing, truncation, summarization, and model change. |
| Scenario first message and transition ordering | §7.1–§7.5 | Eight-step coordinated transaction, typed transition records, explicit eighth step: the model must not continue as if the old scenario were active. |
| Contextual Generate Image from an assistant message | §8.1, `AC-C12` | Three-dot menu → flyover → generate, with the full 12-step lifecycle. |
| Shared Image settings, job path, Amendment-01 orchestration | §8.2, §3.1, `AC-C21` | Prohibits a second loader, settings authority, gallery, registry, or readiness bypass. |
| Image association with conversation, message, operation, board | §8.3, `AC-C05`, `AC-C06`, `AC-C18` | Five-way association, surviving reload. |
| Public and Private routing including nonparticipants | §10.1, §10.2 | Two dropdown sections including nonparticipants; private character does not auto-join. |
| Strict privacy across derived data and later turns | §10.3, §10.4 | Nine leak channels named; seven derived-object classes with per-class treatment. |
| Deletion confirmation for memory and automatic lore | §9.2, `AC-C32`–`AC-C34` | Prompt must identify which types are affected. |
| Cancel deletes message, keeps entries | §9.2, `AC-C36` | Stated twice, including "Cancel does not cancel deletion of the message." |
| Personal lorebooks never included | §9.1, `AC-C37`, §21.3 | Also excluded from migration into automatic lore or memory. |
| Multi-source provenance and shared-source preservation | §9.3, §9.4, `AC-C38`, `AC-C56` | Free-text provenance explicitly rejected as a key; shared contributions preserved. |
| Go/Wails foundation under `OD-01` | §1.1, §11, §14, §19, §21.1 | Settled constraint, not a candidate; evidence-based reopening only. |
| ARES-informed maintainability/inspectability criteria, not rewrite authority | §1.1, `SRC-09` | States the experience informs criteria only. |

`RG-02` is satisfied because the places where the requirement could have been weakened are explicitly strengthened instead: `AC-C43` requires original input rather than another character's description, `AC-C51` states that making a response public does not publish its private source, and `AC-C36` removes the ambiguity in "Cancel."

## 5. Isolation scenario (RG-03)

`AC-C62` reproduces the required scenario line for line: D privately receives an image, receives the actual image, responds, the response enters shared history, A later receives D's response but not the image, a private description, or a revealing summary, and D does not become a participant.

Variants are present: `AC-C63` (audio), `AC-C64` (RAG, including that A cannot retrieve either the private material or a derived summary), and `AC-C65` (private-derived memory or automatic lore retrieved in public context).

Backend enforcement is required, not assumed: §10.5 makes the frontend checkbox/dropdown a request only, and `AC-C50` states the backend decides what each character may receive.

Observed: the current implementation matches the amendment's stated gap. Multi-character orchestration runs in `frontend/dist/js/chat-stream.js`, accumulates `priorResponses` into a frontend-composed message, and calls `GenerateCharacterResponse(message, characterId)` with no audience parameter; the backend accepts that composed message. The amendment's conflict entries `A2-09` and `A2-10` describe this accurately.

## 6. Findings

### M2-01 — Model-input receipts are absent from the audience-propagation table (advisory)

**Evidence class:** Source-derived from the package.

§5.3 makes every model request produce a durable receipt containing ordered event IDs, ordered media IDs, and payload hashes. §10.4 enumerates derived objects that require privacy treatment: memory, automatic lore, summary, RAG chunk/index, attachment reference, generated description, scenario transition. Model-input receipts are not in that list, and `AC-C49`'s leak-channel enumeration does not name them either.

A receipt for a private character request necessarily references private event and media identities. Because receipts are durable and per-request, a future implementation could treat them as shared conversation metadata without violating any stated rule.

**Required before architecture approval:** add the model-input receipt to §10.4 and to `AC-C49`, with the rule that a receipt is retained for user inspection and retained audience, is never supplied to another character's context, and never by itself authorizes media resolution.

**Why this does not block research:** `RQ-07.4` already asks what receipts must contain to prove what a model received. The finding narrows that question rather than opening a new one.

### M2-02 — `A2-08` requires owner resolution but has no decision ID (advisory)

**Evidence class:** Source-derived from the package.

§17 records `A2-08`: the scenario comment says scenarios never override chat state, which conflicts with transition requirements, and it "needs owner-approved resolution." §15 lists `O-01`–`O-10` and does not include it; the review request presents `O-01`–`O-10` as the complete set of owner decisions that must remain open.

`DA-C4` blocks implementation of the transition model, so the risk is not unblocked code. The risk is that `DA-C4` could later be answered as a design question while the product question — whether a scenario transition may override chat-level state — is never put to the owner.

**Required before architecture approval:** either add an `O-11` for scenario-versus-chat state authority, or record in §15 that `A2-08` is tracked under `DA-C4` with an explicit statement that it needs no separate owner decision.

**Why this does not block research:** `R-09` already covers world/scenario/character switching semantics, and the conflict is recorded rather than hidden.

### M2-03 — Sequencing diagram understates cross-workstream dependencies (informational)

§19's diagram shows `R-08`, `R-09`, `R-10`, and `R-11` as siblings under `R-07`. §12.1 additionally makes `R-08` depend on `R-03`, and `R-10` depend on `R-01`, `R-02`, `R-03`, and `DA-I1`. §12.1 is authoritative and the diagram is not wrong so much as incomplete. Research planning should treat §12.1 as the dependency source.

## 7. Evidence and authorization discipline (RG-08, RG-11, RG-12)

`RG-08`: research scope is adequate and explicitly unfinished. Five workstreams carry concrete questions, 31 stable `RQ-*` identifiers, an evidence-class taxonomy, and an explicit statement that no question in the register has been researched. `SRC-08` records the dossiers as "not produced." Marketing claims and frontier-model transfer are expressly rejected as evidence.

`RG-11`: no authorization is implied. §20 enumerates twelve prohibited activities, §14 marks every `DA-C*` decision as blocking implementation, §16 opens by calling all criteria "unvalidated design targets," and §24 stops at the review gate. Requirement verbs such as "must" appear in requirement sections and in §21 as required properties of a later migration proposal, not as permission.

`RG-12`: the review request and §24 define the verdict space and both state that a PASS is a requirement-capture result only. The verification record likewise claims no approval.

One accuracy note, informational: the amendment's current-behavior claims are static observations and are presented as such. My sampling supports them. `frontend/dist/js/chat-session.js` still carries a stale header comment about "Image handling with thumbnails," but its `rebuildChatFromLoadedHistory` restores text and `hasImage` only, which supports `A2-14` rather than contradicting it. No change is requested.

## 8. `OD-01` (RG-13)

`OD-01` is preserved as a settled constraint at seven distinct points: the decision statement (§1.1), the media-ownership split (§5.5), the integration map row (§11), the decision register (§14), the sequencing constraint (§19), governance and migration (§20, §21.1), and the source register (`SRC-09`, `CH-07`).

It is not treated as permission to reopen the language question. §1.1 restricts reopening to concrete engineering evidence of a limitation that cannot reasonably be resolved within Go/Wails, §14 places any such proposal outside the amendment's scope and behind a separate owner decision, and §19 forbids any workstream or migration phase from treating a rewrite as an available option, fallback, or optimization. Ordinary Go/Wails refactoring is explicitly preserved.

`OD-01` is also correctly distinguished from `O-01`–`O-10`: it is not one of the open questions and is not counted among them.

## 9. What this verdict does not authorize

- No source, configuration, model, database, or runtime change.
- No migration or cutover.
- No model load, benchmark, or runtime probe.
- No network, Android pairing, or packet capture.
- No Public/Private implementation.
- No contextual-image implementation.
- No production behavior change.
- No architecture approval and no acceptance result.
- No language or backend rewrite.

Research planning for R-07–R-11 may begin only after the owner reviews this verdict and explicitly authorizes it.

## 10. Recommended next step

The owner reviews this PASS. On approval, R-07–R-11 planning proceeds with:

1. `M2-01` carried into `RQ-07.4` and the `R-10` leakage matrix.
2. `M2-02` carried into `R-09`, with the owner asked to decide `A2-08` before architecture approval.
3. §12.1 treated as the dependency authority, not the §19 diagram.
