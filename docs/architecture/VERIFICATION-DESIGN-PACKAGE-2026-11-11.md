# Verification record — corrected architecture design package

**Date:** 2026-11-11  
**Type:** Static documentation and integrity verification  
**Scope:** Owner-decision addendum and four architecture designs plus reconciliation  
**Status:** Targeted review corrections applied; no implementation, runtime, model, benchmark, migration, or network activity

## 1. Baseline integrity

| Artifact | Expected | Result |
|---|---|---|
| Corrected plan, Revision 3 | `d642039cd5ddcdc480001666f7a18ad36e22944f5540349056af1eee050ee874` | **Verified unchanged** |
| Requirements reconciliation | `7ad9d072d43c1243431a91cc304911b5dd8beeecc10923808115a306febda3ab` | **Verified unchanged** |
| Recovery owner addendum | `c9f8dc1c08ad5c65954e77c11270ae4d7de18449cd765aebd04cba6a1df7f056` | **Verified unchanged** |
| Independent architecture review | `f1756fb5e6c21b8dce5ed30e677426f581ecfc732463a30eb4a630d34d8d5eee` | **Verified; reviewed the pre-correction package** |

The Revision 3 plan was not edited by this package. The owner-decision addendum is the separate record of planning acceptance. The independent review's verdict remains authoritative: the package is a completed first architecture-design pass, **not approved architecture**.

## 2. Review corrections applied

| Review correction | Applied in | Status |
|---|---|---|
| **F-01** Character voice/TTS belongs to Chat; Image is next; Music/3D are future dedicated studios | Addendum `UD-05`/`DPR-05`/§3.1/§5; hardware §8.1; reconciliation `ADR-STU-01`, `DWP-08`, §7.1 | Applied |
| **F-02** Replace event-log absolutism with auditable append-oriented history and explicit correction/supersession/retraction/deletion/tombstone/as-of semantics | Memory §1, §3, §4.2; reconciliation `ADR-MEM-01` and conflict table | Applied |
| **F-03** Add a memory migration bridge from current JSON, candidates, continuity config, legacy ownership, and frontend caches | Memory §9.1; blocking `DA-M5`/`DA-M6` | Applied |
| **F-04** Complete resource lifecycle: atomic lease/commit, contention, expiry, rollback, release verification, live changes, existing allocations, and all consumers | Hardware §2, §8.4–§8.7, §11, §13 | Applied |
| **F-05** Define prompt-path cutover and retirement of existing contributor paths | Reconciliation §6.5; `ADR-INT-01` | Applied |
| **F-06** Add RAG migration, embedding/index identity, deletion receipts, stale citations, corpus separation, and retrieval seam hardening | RAG §2, §12.1, §14 | Applied |
| **F-07** Clarify DWP order as design-input order and preserve Chat → Image → expansion staging | Reconciliation §7, §7.1, §8; `DWP-08` | Applied |
| **F-08** Keep Android residual items attached to `DWP-07` | Reconciliation `DWP-07` row | Applied |
| **F-09** Separate current behavior from design targets and remove the verification self-hash | This record §4–§6; external `.sha256` sidecar | Applied |

No correction was treated as architecture approval. No correction authorized implementation or runtime work.

## 3. Deliverables

| Document | Purpose |
|---|---|
| `docs/review/OWNER-DECISION-ADDENDUM-REVISION-3-2026-11-11.md` | UD-01–UD-06, RAG, OA-09, hardware decisions, and delta register |
| `docs/architecture/DESIGN-MEMORY-CONTINUITY-2026-11-11.md` | Canonical state, auditable event history, memory retrieval, small-model context, migration bridge |
| `docs/architecture/DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md` | Hardware inventory, telemetry, 10% reserve, suitability, atomic leases, release/rollback |
| `docs/architecture/DESIGN-LOCAL-RAG-CREATIVE-2026-11-11.md` | Structure-aware local RAG, hybrid retrieval, citations, creative modes, migration/seam bridge |
| `docs/architecture/RESEARCH-ANDROID-REMOTE-SECURITY-2026-11-11.md` | Android threat model, transport comparison, pairing, E2E, revocation |
| `docs/architecture/RECONCILED-ARCHITECTURE-PLAN-2026-11-11.md` | Single-owner boundaries, design packages, prompt cutover, product staging |
| `docs/review/INDEPENDENT-ARCHITECTURE-REVIEW-2026-11-11.md` | Independent review and targeted correction list |

## 4. Coverage matrix

| Requested area | Coverage |
|---|---|
| Preserve Revision 3 | Plan not edited; acceptance recorded separately |
| Record UD-01–UD-06 | Addendum §3; delta register §4 |
| RAG commitment | Addendum `ODA-07`/`DPR-07`; RAG design |
| OA-09 clarification | Addendum `ODA-08`; Android/hardware/reconciliation boundaries |
| Automatic editable memory | Memory design §5–§8 |
| Clothing, equipment, location, activity, relationships, events | Memory design §4 |
| History distinct from canonical state | Memory design §4.2, §6 |
| Correction, supersession, retraction, deletion, forget, as-of | Memory design §4.2, §6–§7 |
| No invented state/silent overwrite | Memory design §5.2, §10 |
| Small-model retrieval/context | Memory design §7 |
| Howl owns authoritative state | Memory design §3, §5; reconciliation §4 |
| Memory migration/atomicity/stale cache | Memory design §9.1 |
| Hardware detection/suitability | Hardware design §4, §7 |
| CPU/RAM/VRAM/context/KV/offload/storage | Hardware design §4–§7 |
| 10% free VRAM reserve | Hardware design §6 |
| Vendor-neutral telemetry | Hardware design §5 |
| Atomic lease/commit and contention | Hardware design §8.4–§8.5 |
| Rollback, expiry, verified release, live changes | Hardware design §8.6–§8.7, §11 |
| Conservative monitor fallback | Hardware design §9.1 |
| Power-user overrides | Hardware design §9.3 |
| Local RAG indexing | RAG design §4–§7 |
| Hybrid keyword/semantic retrieval | RAG design §8 |
| Structure-aware retrieval | RAG design §5, §8 |
| Creative modes | RAG design §8.1 |
| Source references/evidence | RAG design §9 |
| Reference separated from canon | RAG design §3, §10, §14 |
| RAG migration, index identity, deletion receipts | RAG design §2, §12.1, §14 |
| Android desktop host | Android design §2, §5 |
| No phone model/cloud inference | Android design §2 |
| P2P/NAT/rendezvous/tunnel/relay tradeoffs | Android design §6, §8 |
| Pairing/mutual auth/revocable credentials/E2E | Android design §7 |
| No public listener/third-party chat storage | Android design §4, §8–§10 |
| No transport chosen prematurely | Android design §6, §11 |
| One prompt path and cutover | Reconciliation §6.5 |
| Chat/TTS destination and Image-next staging | Addendum `UD-05`; reconciliation §7.1 |
| DWP order is not build order | Reconciliation §7, §8 |
| Android residual security items | Reconciliation `DWP-07` |
| Reconcile designs | Reconciliation §4–§9 |

## 5. Contradiction and duplication checks

- One canonical state/event owner: memory design; no second memory authority.
- Frontend memory cache is explicitly non-authoritative; migration must prevent stale-cache resurrection.
- RAG cannot write canonical state; reference evidence is labeled and separately sourced.
- Lorebooks remain distinct from RAG and character state.
- **Design target:** one resource authority covers chat, image, TTS/ASR, embedding, and future studios. Current code remains advisory and incompletely registered.
- **Design target:** one prompt assembler consumes state, memory, lore, and RAG evidence. Current code still has several prompt paths; §6.5 defines their retirement.
- Character voice/TTS is a Chat capability, not a separate studio. TTS/ASR is still a resource consumer.
- History is auditable and append-oriented; no event is treated as an unchangeable record. Deletion/forget semantics are explicit and blocking.
- Android reuses desktop policy and resource authority; it does not create a parallel model host.
- Product staging is Chat → Image → expansion; DWP numbering is design-input order only and is not modality authorization.
- UD-01 resolves per-chat model scope without making hot-swap a release blocker.
- RAG remains outside the 66-row reconciliation as `DPR-07`; no hidden 67th row.
- OA-10 product-scope commitment is separated from implementation/lifecycle authorization.

## 6. Current behavior versus design target

| Topic | Current inspection finding | Design target |
|---|---|---|
| Memory writes | In-place JSON edits/deletes; read-modify-write without a transaction or revision check; placeholder backup; second-resolution IDs | Event/projection store with atomic staging, unique IDs, CAS, tombstones, and stale-cache guards |
| Event history | No general correction/supersession/as-of history | Auditable append-oriented history with explicit correction, retraction, deletion, and as-of semantics |
| Prompt assembly | Multiple contributors and call sites | One assembler with declared precedence, exactly-once injection, and retired legacy paths |
| Resource admission | Advisory registry, 90% free heuristic, fixed release wait, force-load bypass, unregistered consumers | Atomic lease/commit, 10% reserve, verified release, rollback, and one authority for all consumers |
| RAG | Flat/in-memory index, no embedding identity, read-path counter mutation, no empty-history guard, no durable deletion receipt | Versioned hybrid index, persisted identity, read-only retrieval, deletion receipts, and labeled stale evidence |
| RAG migration | No safe bridge from the current index | Staged manifest migration with rollback and corpus/rights preservation |
| TTS product placement | Addendum previously sequenced TTS as a separate studio; current catalog is already Character-tab functionality | Character voice/TTS in Chat; Image next; Music/3D future dedicated studios |
| Android | No transport selected; security design open | Desktop-hosted, off-by-default, paired, mutually authenticated, E2E, revocable remote access |
| DWP package order | Design packages can be misread as a build sequence | Design-input order; product staging is separate and gated |

These are **static inspection and design statements**, not runtime evidence. No behavior in this table was exercised.

## 7. Technical research status

Android research consulted:

- WireGuard known limitations and protocol documentation;
- RFC 8656 / RFC 5128 on TURN and NAT traversal;
- libp2p hole-punching and Circuit Relay v2 documentation;
- Tailscale NAT-traversal discussion.

Hardware research consulted:

- NVIDIA NVML API documentation, including WDDM memory-reporting caveats;
- AMD ADL/PDH and Windows GPU-memory counter discussion;
- Intel Level Zero Sysman as a provider candidate.

These are design inputs, not runtime evidence. No transport was selected. No telemetry provider was tested.

## 8. Verification limitations

- No source code was implemented or changed by this package.
- No build, test, migration, schema change, model load, benchmark, packet capture, network listener, pairing, or Android run was performed.
- No runtime performance, memory, security, or compatibility claim is made.
- No design is accepted as implementation-ready.
- The independent architecture review examined the pre-correction package. The F-01–F-09 corrections have not received a further independent review or architecture approval.
- The earlier MUSE review applies to `5ff5e1d4…`; Revision 3 includes documentary O-1–O-3 changes and was not re-reviewed by MUSE here. This record does not claim that review.
- The verification record's own digest is externalized in `VERIFICATION-DESIGN-PACKAGE-2026-11-11.md.sha256`; it is not embedded in this file.

## 9. Final authorization state

| Question | Answer |
|---|---|
| Planning baseline accepted? | Yes, by separate owner addendum |
| Review corrections applied? | **Yes, documentation only** |
| Architecture approved? | **No** |
| Independent re-review completed? | **No** |
| Implementation authorized? | **No** |
| Runtime investigation authorized? | **No** |
| Android implementation authorized? | **No** |
| RAG implementation authorized? | **No** |
| Memory migration authorized? | **No** |
| Resource enforcement implemented? | **No** |
| Any gate passed? | **No** |
| Any CWP accepted? | **No** |
| Ready for integration/acceptance planning? | **Yes, after re-review of these corrections** |

## 10. Final hashes

| Deliverable | SHA-256 | Size / lines |
|---|---|---|
| `docs/review/OWNER-DECISION-ADDENDUM-REVISION-3-2026-11-11.md` | `d228c4ac8a1acfabc7164ab5b277e9c20e8e2c33522554ced4896ec02b41e5ac` | 11939 B / 141 lines |
| `docs/architecture/DESIGN-MEMORY-CONTINUITY-2026-11-11.md` | `2e8ec935843bfd6b110feb16a626ce856cff60dd9ec1f570df0d21b21ef57650` | 20764 B / 421 lines |
| `docs/architecture/DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md` | `9348fa0dc297db356f8391f66318e98f8310b3ee6900b7116137ac23f852b0c3` | 20240 B / 392 lines |
| `docs/architecture/DESIGN-LOCAL-RAG-CREATIVE-2026-11-11.md` | `2a80de38267131c7c4557729b845ac0bd43b3e345ad395020b93deb9e738e7e6` | 16394 B / 361 lines |
| `docs/architecture/RESEARCH-ANDROID-REMOTE-SECURITY-2026-11-11.md` | `c25e45b511b883c26d44ac8d3b61d66d6397a12a13f120f035a73f75474e33a2` | 15641 B / 318 lines |
| `docs/architecture/RECONCILED-ARCHITECTURE-PLAN-2026-11-11.md` | `58e33ed918e22628854a299380b9178b3cf61481b3a82cdebda94c459ee10040` | 22008 B / 354 lines |

The verification record's own SHA-256 is recorded in `docs/architecture/VERIFICATION-DESIGN-PACKAGE-2026-11-11.md.sha256`. This avoids a self-referential hash.
