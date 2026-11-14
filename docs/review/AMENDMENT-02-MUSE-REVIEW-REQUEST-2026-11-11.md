# Amendment 02 — independent review request

**Date:** 2026-11-11  
**Reviewer:** Muse  
**Review type:** Independent requirement-capture and planning review  
**Not requested:** architecture approval, implementation authorization, runtime investigation, or research execution  
**Decision requested:** PASS / PASS WITH REQUIRED CHANGES / FAIL on whether the amendment accurately captures the owner-defined requirements and is ready for research planning.

## 1. Package under review

| Artifact | SHA-256 | Size |
|---|---|---|
| `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md` | `fa766f76779e2517b31bd55da993f872525a52941b348fc8f78bbeb1d9382c44` | `61,111 B / 1,122 lines` |
| `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md.sha256` | Sidecar verifies the amendment. | — |
| `docs/architecture/AMENDMENT-02-VERIFICATION-RECORD-2026-11-11.md` | Supporting verification record. | — |
| `docs/architecture/INTEGRATION-AND-ACCEPTANCE-PLAN-2026-11-11.md` | `1e92a54434e8e3269015c47374ad57bd6ae755e64e27a07954c083ab952510cf` | Preserved |
| `docs/architecture/AMENDMENT-01-UNIFIED-HARDWARE-INTELLIGENCE-2026-11-11.md` | `f3a0f8b9c5d19129db85d559f1880f21451cc7aa70596d37a408952c85c3f4ea` | Preserved |

The original plan and Amendment 01 must remain unchanged. Amendment 01 remains controlling for hardware intelligence, model residency, scheduling, readiness, release verification, and manual resource controls. `OD-01` is a resolved owner constraint: Go/Wails remains the technology foundation, and review must not treat the amendment as permission to reopen a language debate or propose a wholesale rewrite.

## 2. Review questions

Muse should answer each question explicitly and cite the amendment section or line where the evidence is found.

| ID | Review question |
|---|---|
| `RG-01` | Is every owner-defined requirement captured accurately? |
| `RG-02` | Has any requirement been weakened, silently reinterpreted, or converted into an implementation assumption? |
| `RG-03` | Is Public/Private isolation addressed across prompts, media, memory, automatic lore, summaries, RAG, persistence, and subsequent turns? |
| `RG-04` | Are canonical transcript history and model input kept distinct and independently inspectable? |
| `RG-05` | Are personal lorebooks, automatic lore, and automated memory kept distinct? |
| `RG-06` | Are durable media identity, reuse, persistence, deletion, and Chat/Image-board association covered? |
| `RG-07` | Do Chat image generation and multi-character runs use Amendment-01 orchestration? |
| `RG-08` | Is research scope adequate, evidence-led, and explicit about unestablished behavior? |
| `RG-09` | Are dependencies and acceptance criteria testable, including constrained systems? |
| `RG-10` | Are the original plan, Amendment 01, frozen documents, and prior digests intact? |
| `RG-11` | Is implementation or runtime authorization accidentally implied anywhere? |
| `RG-12` | Does the package support a clear PASS/FAIL review verdict without treating review as architecture approval? |
| `RG-13` | Does the package preserve `OD-01`—Go/Wails as the foundation, no speculative rewrite, and concrete evidence as the only reopening condition—without turning review into permission to reopen the language debate? |

## 3. Required checks

### 3.1 Requirement fidelity

Confirm that the amendment captures:

- canonical visible history independent of model context;
- durable media identity and lifecycle;
- the protected initial message at the absolute top, editable but not deletable;
- scenario first message and transition ordering;
- contextual Generate Image from an assistant message;
- shared Image settings, job path, and Amendment-01 orchestration;
- image association with conversation, message, operation, and Image board;
- Public and Private routing, including nonparticipant private recipients;
- strict privacy across derived memory, automatic lore, summaries, RAG, attachments, and subsequent turns;
- message deletion confirmation for automated memory and automatic lore;
- Cancel means “delete the message, keep associated system-generated entries”;
- personal lorebooks are never included;
- multi-source provenance and shared-source preservation;
- Go/Wails remains the technology foundation under `OD-01`, with no language or backend rewrite based on preference or speculative advantage and concrete engineering evidence as the only reopening condition;
- AI-assisted maintainability and inspectability remain explicit evaluation criteria, informed by the owner's ARES experience, without being treated as a rewrite authorization.

### 3.2 Isolation scenario

Verify that the amendment requires this exact scenario:

1. Character D is privately sent an image.
2. D receives the actual image and responds.
3. D's response appears in shared history.
4. Character A subsequently receives D's response.
5. A does not receive the original image, a private description, or a summary that reveals its contents.
6. D does not become a group participant merely because it received private content.

Confirm the corresponding audio and RAG variants and private-derived memory/summary tests are present.

### 3.3 Evidence discipline

Confirm that the amendment distinguishes:

- owner-defined requirements;
- proposals;
- observed static implementation;
- source-derived findings;
- calculated or inferred conclusions;
- unestablished behavior;
- frozen and superseded material.

Confirm that it does not present R-07–R-11 as completed research or claim any acceptance result.

### 3.4 Authorization discipline

Confirm that the amendment does not authorize or imply:

- source changes;
- migrations;
- model loads;
- benchmarks;
- runtime probes;
- network or Android pairing;
- Public/Private implementation;
- contextual image implementation;
- production behavior changes;
- a wholesale Go/Wails language or backend rewrite;
- architecture approval.

## 4. Known ambiguities that must remain open

These are recorded as `O-01`–`O-10` in the amendment and must not be silently resolved by the reviewer or implementation. `OD-01` is not one of these open questions; it is a resolved owner decision and must remain settled:

- `O-01` automatic assistant-message prompt extraction;
- `O-02` Public delivery across character/model configurations;
- `O-03` private-derived memory/lore when independently established;
- `O-04` private-event representation in shared visible history;
- `O-05` scenario transitions versus active/queued generation;
- `O-06` historical edits and derived state;
- `O-07` message deletion versus Image-board assets;
- `O-08` protected initial-message role and content model;
- `O-09` nonparticipant access to shared scene history;
- `O-10` scenarios without a first message.

## 5. Required verdict

Use one of:

- **PASS** — requirements are accurately captured and the package is ready for the affected research planning.
- **PASS WITH REQUIRED CHANGES** — list exact changes; research may not begin until they are applied and rechecked.
- **FAIL** — identify the missing, weakened, or contradicted requirements; research may not begin.

A PASS is a requirement-capture review result only. It is not architecture approval, implementation authorization, or acceptance completion.

## 6. Stop condition

After issuing the verdict, stop. Do not begin R-07–R-11 research, implementation, migration, runtime investigation, model loading, benchmarking, network testing, or Android pairing until:

1. the review verdict is PASS or required changes are closed;
2. the owner reviews the result;
3. architecture approval and the relevant authorization gate are explicitly granted.

No source, configuration, model, database, or runtime changes were made in preparing this package.
