# MEMORY-09 — Autonomous Classification Experiment Plan

## Capture completeness, classification capability, and the unresolved-class lifecycle

**Status:** **experiment plan; G0 accepted; G1 next.** No L-2 amendment, no Decision A amendment, no schema, code,
migration, test, model fetch, corpus creation or canonical capture. The owner's direction is to establish technical
achievability and the shape of unresolved cases **before** amending L-2. The gates G0–G5 are fixed.

| Field | Value |
|---|---|
| Date | **2026-11-21** (local, CST) |
| Owner direction | capture before classification; hybrid is an investigation direction, not adopted; do not amend L-2 yet |
| Three questions | capture completeness · classification capability · unresolved-class lifecycle |
| BGE | **embedding candidate-ranking experiment only**, not class authority; the evaluated BGE reranker is a different artifact |
| Current corpus | none for class labels; the 0/144 experiment state is unchanged |
| Current capture | not a functioning backend character-memory path |
| L-2 / Decision A | **unchanged** |
| G0 | **accepted** (2026-11-21) — scope fixed |
| G1 | **fixture contract passed**; production wiring not established |
| G3 | **not authorized** — no BGE fetch |
| Implementation | **0** |
| Map | **118** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed**.

---

## 1. The sequencing decision

**ACX-01 (Observed — the owner's direction).** The capture problem comes before the classification problem. A
classifier cannot be evaluated against information the capture pipeline does not supply.

**ACX-02 (Observed — separate responsibilities).**

| Operation | Question |
|---|---|
| **capture** | what happened, who experienced it, and where the information came from |
| **classification** | which knowledge category the captured information belongs to |
| **formation** | how the captured information becomes an authoritative memory |

**ACX-03 (Calculated — no compensation by guessing).** Missing capture information must not be repaired by a classifier
inference. The correct outcome is **capture-incomplete** or **unresolved**, not a default class.

**ACX-04 (Observed — L-2 stays intact).** This plan does not amend L-2's one formation entry or its `class_baseline`.
The unresolved-class lifecycle is investigated as a **decision question**, not implemented.

---

## 2. What is established

### 2.1 The capture gap

**ACX-05 (Observed — the backend path).** `Service.ProcessMessageForEvents(message, history, identities)` receives only
message text, history strings, and `CharacterIdentity{ID, Name}`. `CharacterIdentity` contains no role, metadata or
attribution.

**ACX-06 (Observed — the frontend path).** The browser wrapper calls
`window.go.main.App.ProcessMessageForEvents(message, [], [])` after the user sends a message. It passes empty history
and identities, and it runs on the **user's outgoing message**.

**ACX-07 (Observed — no App binding).** `frontend/wailsjs/go/main/App.d.ts` contains no
`ProcessMessageForEvents` function. The frontend check therefore falls through to its frontend-only keyword fallback.
The backend auto-capture pipeline is not exposed to the frontend in the inspected build.

**ACX-08 (Observed — the detector).** `DetectEvents` is keyword-based and returns `EventType`, `Confidence`, `Evidence`
and `Context`. It has no source character, attribution, world/scenario metadata or class. “I grew up in the northern
mountains” has no matching keyword in the inspected map.

**ACX-09 (Observed — what exists elsewhere).** `ChatHistory` carries `ID`, `Sender`, `Text`, `Role`, `Timestamp`,
`HasImage` and `CharacterID`. `ChatData` carries `World`, `Scenario`, `Characters` and `Meta` with stable world,
scenario and character IDs. This information is available at the app/chat layer but is not passed into the memory
service.

### 2.2 BGE is two different things

**ACX-10 (Observed — the evaluated BGE artifact).** `MEMORY-09-S3-BGE-EVALUATION-2026-11-15.md` evaluates
`bge-reranker-v2-m3` — a **cross-encoder reranker**. It was found technically retainable but behaviourally limited,
failed the temporal product gate, and is **not shipped**. It is not an embedding model and must not be reused as one.

**ACX-11 (Calculated — the proposed BGE use is a bi-encoder).** Comparing a memory against class descriptions requires
an **embedding model** that produces vectors for both texts and cosine similarity. That is a different artifact from
the reranker. No BGE embedding model has been evaluated or pinned in this project.

**ACX-12 (Calculated — similarity is not authority).** A BGE top-1 similarity is a **candidate ranking**. It cannot
establish authorship, ownership, audience or eligibility. Two classes can be semantically close while having different
consequences. BGE may suggest; it may not assign.

### 2.3 Corpus and experiment state

**ACX-13 (Observed — no class-labelled corpus).** There is no corpus of memories labelled with the four knowledge
classes and the capture fields needed to evaluate classification. The existing memory-evaluation corpus is not a
class-label corpus.

**ACX-14 (Observed — no experiment).** 0/144 sessions, no `plan_sha256`, G07/G12 unbound, `ExecutionNotAuthorised` at
`conversation.py:177`. The BGE reranker evaluation is a prior, separate experiment and is not re-run.

---

## 3. The three questions

**ACX-15 (Calculated).**

| Track | Question | Produces | Gate |
|---|---|---|---|
| **1** | what must capture provide so classification has adequate context? | a minimum capture envelope and fixture coverage | no classifier measurement without it |
| **2** | can deterministic rules, metadata and optionally BGE classify enough memories without external AI or human review? | coverage/accuracy/abstention and BGE added value | a capability screen, not production accuracy |
| **3** | can a memory remain useful while class is unresolved without violating Decision A, L-2 or class-sensitive retrieval? | a lifecycle/retrieval analysis and an owner decision point | no L-2 amendment before it |

**ACX-16 (Calculated — order).** Track 1 first. Track 2 only on captured fixtures that satisfy Track 1. Track 3 can
run in parallel because it is a lifecycle/decision analysis, not a classifier measurement.

---

## 4. Track 1 — capture completeness

### 4.1 Minimum capture envelope

**ACX-17 (Calculated — logical fields).**

| Field | Why classification needs it | Current source |
|---|---|---|
| `source_message_id` / span | provenance; replay | `ChatHistory.ID` exists; not passed |
| `speaker_role` | user / assistant / narrator / system | `ChatHistory.Role` exists; not passed |
| `source_character_id` | who experienced or asserted it | `ChatHistory.CharacterID` exists; not passed |
| `source_character_name` | display only; never ownership | `Sender` exists; not passed |
| `chat_id` | isolation | service has it |
| `turn_id` / sequence | local ordering | history order exists; not passed |
| `event_type` | detection signal | `EventDetectionResult` |
| `evidence` / source text | explicit rule evidence | `EventDetectionResult.Evidence`, `Context` |
| `world_id`, `scenario_id` | authored-context check | `ChatMeta` exists; not passed |
| `attribution` (observed/told/believed/asserted/authoritative) | experience vs belief | **absent** |
| `authored_asset_ref` | ambient/restricted eligibility | **absent**; auto-capture is not an authored asset |
| `existing_class` | legacy/migration only | **absent**; no class field |
| `content` | the memory text | message |

**ACX-18 (Observed — availability matrix).**

| Field | Available in app/chat | Passed to memory service | Status |
|---|---|---|---|
| message content | yes | yes | available |
| chat ID | yes | yes | available |
| speaker role | yes (`ChatHistory.Role`) | **no** | **missing** |
| source character ID | yes (`ChatHistory.CharacterID`) | **no** | **missing** |
| source character name | yes (`Sender`) | **no** | **missing** |
| world/scenario IDs | yes (`ChatMeta`) | **no** | **missing** |
| world/scenario descriptions | yes (`ChatData`) | **no** | **missing** |
| character metadata | yes (`ChatData.Characters`) | **no** | **missing** |
| attribution | **no** | **no** | **absent** |
| authored asset reference | **no** | **no** | **absent** |
| class | **no** | **no** | **absent** |

**ACX-19 (Calculated — the minimum fix).** For classification to be evaluable, the capture call must carry at least:
source message/span, speaker role, stable source character ID (or explicit none), chat/turn identity, event type and
evidence, world/scenario IDs, and attribution. Without attribution, the classifier must return unresolved for the
experience/belief distinction; without an authored asset reference, it must not assign ambient/restricted.

**ACX-20 (Calculated — capture-incomplete is a state).** A candidate missing a required field is marked
**capture-incomplete**, not classified. This is distinct from a classifier abstention. The two must not collapse.

### 4.2 Fixture design

**ACX-21 (Calculated — controlled transcript fixture).** A synthetic transcript with:

- a character assistant message carrying a stable `CharacterID`;
- a user message;
- a narrator/system message;
- two characters with an ambiguous name collision;
- a missing-`CharacterID` legacy message;
- a message that is captured but has no attribution;
- a message with explicit authored world/scenario context.

**ACX-22 (Calculated — assertions).**

| Assertion | Expected |
|---|---|
| character message | source character ID is the stable ID, not a display-name guess |
| user message | source character is explicitly none; not a character |
| ambiguous name | capture is incomplete; no guess |
| missing metadata | capture is incomplete; no default |
| world/scenario | IDs and authored asset reference are present when applicable |
| attribution | preserved when supplied; absent is a state, not “observed” |

**ACX-23 (Calculated — no classifier in Track 1).** Track 1 uses no class labels and no model. It proves only that the
capture envelope is adequate and that incompleteness is visible.

### 4.3 Track 1 stop conditions

**ACX-24 (Calculated).** Stop the classification experiment if:

- source character cannot be stable for character messages;
- user/narrator/character roles are not distinguishable;
- attribution is unavailable and the experiment still claims experience/belief classification;
- authored asset status is unavailable and the experiment still claims ambient/restricted classification;
- `capture-incomplete` is treated as unresolved class rather than a separate state.

---

## 5. Track 2 — classification capability

### 5.1 Arms

**ACX-25 (Calculated).**

| Arm | Condition | Purpose |
|---|---|---|
| **R** | deterministic rules/metadata only | baseline without any embedding model |
| **B** | BGE embedding top-1 only | diagnostic; no authority |
| **B+margin** | BGE top-1 with a calibrated margin; abstain otherwise | candidate ranking with abstention |
| **R+B** | rules gate; BGE ranks candidates only when rules abstain; conflicts/low margin → unresolved | the proposed hybrid classifier |

**ACX-26 (Calculated — R must work without BGE).** BGE is optional infrastructure. If BGE is absent, R still runs and
the system reports a degraded classifier, not a default class.

### 5.2 Vocabulary and class descriptions

**ACX-27 (Source-derived — fixed vocabulary).** The four classes are `ambient`, `restricted`, `experiential`,
`character_specific`. “Preference” is an example of `character_specific`, not a fifth class. “World knowledge” is
`ambient` only when it is an authored world/scenario definition (K-R8).

**ACX-28 (Calculated — class descriptions for BGE).** Freeze a small set of authored descriptions per class, for
example:

| Class | Description family |
|---|---|
| `ambient` | publicly available authored world/scenario definition; a newcomer may be given it without acquisition |
| `restricted` | authored secret/hidden world/scenario knowledge; requires discovery or grant |
| `experiential` | a character's witnessed or lived event/experience |
| `character_specific` | a character's preference, belief, opinion, emotion or personal association |

**ACX-29 (Calculated — the adversarial case).** “The northern kingdom is ruled by Queen Elara” is semantically close
to a world fact, but if it is auto-captured character speech it is **not ambient** (K-R8). The experiment must include
this case and treat an ambient prediction as a critical error, not a tolerable miss.

### 5.3 Corpus

**ACX-30 (Calculated — controlled class corpus).** Minimum pilot: 40 examples per class (160), plus 40 ambiguous or
capture-incomplete cases, split before the run into development/calibration and held-out evaluation. Examples must
carry the Track 1 envelope; examples that lack it are capture-incomplete, not class examples.

**ACX-31 (Calculated — labels).** Labels may be authored for the experiment; this is not production manual review.
Labels must be frozen before scoring; no threshold or example may be tuned on the held-out set.

**ACX-32 (Observed — corpus limitation).** Synthetic examples are not real conversations. The result is a **capability
screen**, not a production accuracy claim. Real-session validation remains future work.

### 5.4 BGE experiment

**ACX-33 (Calculated — model candidate).** A BGE **embedding** model, not the evaluated reranker. Candidate families
include BGE-M3 (multilingual, larger) and BGE-small-en-v1.5 (small, 384-d). The owner must name the candidate; no
model is fetched by this plan.

**ACX-34 (Calculated — procedure).**

1. Embed each memory text and each frozen class description.
2. Cosine similarity; rank classes.
3. Record top-1, top-2 margin, and the full ranking.
4. Calibrate the abstention margin on the held-out calibration split, never on the evaluation split.
5. If the top two classes are within the calibrated margin, return unresolved.
6. BGE never overrides an explicit rule or an authorship/ownership exclusion.

**ACX-35 (Calculated — what BGE may decide).** Only “these class descriptions are closer/farther” as a candidate
signal. It may not decide ambient/restricted without authored status; it may not decide ownership; it may not turn a
capture-incomplete case into a class.

### 5.5 Metrics

**ACX-36 (Calculated).**

| Metric | Why |
|---|---|
| coverage | fraction resolved, overall and per class |
| accuracy on resolved | correctness among resolved cases; never hide abstention |
| per-class precision/recall | especially ambient and restricted |
| abstention rate | the honest-uncertainty rate |
| confusion matrix | which classes are confused |
| BGE top-1 / top-3 | candidate-ranking quality |
| margin distribution | whether similarity is decisive or flat |
| **false-confidence rate** | confidently assigned **incorrect** class; per class and per stream |
| **separate streams** | rules-only · BGE recommendations · hybrid recommendations · human-established reference; never merged |
| R vs B+margin vs R+B | added value, not a single headline |
| latency / memory | whether it fits the formation path |

**ACX-37 (Calculated — vetoes).** The following are vetoes, not aggregate misses:

- auto-captured content assigned `ambient`;
- a character's private knowledge assigned `restricted` without authored status;
- a user/narrator memory assigned a character's ownership;
- a capture-incomplete case assigned any class;
- a BGE similarity used to override an explicit rule or authored-status exclusion.

**ACX-38 (Calculated — success criterion).** The experiment supports the hybrid only if R or R+B resolves a useful
share of cases, BGE adds measurable value over R, the abstention rate is honest, and no veto fires. “Useful share” is
declared before the run; it is not chosen after seeing results.

**ACX-39 (Calculated — falsification).** BGE adds no value over R; BGE is confidently wrong on the adversarial ambient
case; rules cannot exceed a trivial baseline; or abstention is so high that the system cannot classify enough to be
useful. Any of these is a result.

### 5.6 Resource and artifact controls

**ACX-40 (Calculated — resource envelope).** Reuse the existing budget: VRAM delta ≤ 384 MiB, warm p50 ≤ 40 ms,
p95 ≤ 120 ms, cold load ≤ 3 s, residual ≤ 64 MiB. BGE-M3 may not fit GPU residency; CPU/partial modes must be
measured. No budget is relaxed silently.

**ACX-41 (Calculated — artifact pinning).** Pin repository, revision, file, bytes, SHA-256 and licence before use.
The evaluated `bge-reranker-v2-m3` artifact is not reused. No model fetch without owner authorization.

**ACX-42 (Calculated — no external AI).** No DeepSeek or network classification. BGE runs locally; rules run locally;
absence of either produces an explicit state, not a default.

---

## 6. Track 3 — unresolved-class lifecycle

### 6.1 The exact question

**ACX-43 (Calculated — the contradiction).** Decision A requires the initial `knowledge_class` to be reconstructible
from the authoritative stream. A formation entry that records **unresolved** and a later resolution entry that supplies
the first class reconstructs an unresolved state, not a class value. The owner must decide whether that satisfies
Decision A or requires an amendment.

**ACX-44 (Calculated — the status-slot precedent is not transferable automatically).** N-b authorized an unresolved
**status** slot; class was not included. The class slot must be explicitly authorized or excluded.

**ACX-45 (Calculated — the first resolution is not a class change).** The first class resolution completes the
unresolved baseline. It must not be represented as a P-1 `class_changed`; subsequent changes are P-1 changes. This
distinction is required whether or not the owner adopts the hybrid.

### 6.2 Lifecycle scenarios

**ACX-46 (Calculated).**

| Scenario | Required behavior |
|---|---|
| resolved at formation | class baseline is the value; no resolution entry |
| unresolved at formation, resolved later | unresolved is preserved until a resolution entry supplies the first value |
| unresolved forever | class remains unresolved; no default |
| character tells another about their past | capture/provenance preserved; non-class-dependent use available; class-sensitive use waits; class unresolved; no interruption or invented class |
| resolution conflicts with a later change | resolution completes baseline; change versions after it |
| legacy record without class | reported as unmigrated; never guessed |
| capture-incomplete | not formed/classified; distinct from unresolved |

**ACX-47 (Calculated — retrieval).** Unresolved class cannot be treated as any class. Class-sensitive audience and
eligibility must fail closed or explicitly show unresolved. Non-class operations (content, ownership, status, display,
inspection) may remain available. “Useful while unresolved” must be defined against this boundary, not assumed.

**ACX-48 (Calculated — no L-2 amendment yet).** Track 3 produces the evidence for the owner's decision. It does not
edit L-2, Decision A, or the implementation plan. If the owner rejects an unresolved class, the hybrid is blocked and
only an always-resolving classifier (with fabrication risk) or an external source (Option A, rejected) remains.

---

## 7. Sequencing and gates

**ACX-49 (Calculated).**

| Gate | Content | Exit |
|---|---|---|
| **G0** | this plan accepted | **accepted 2026-11-21**; scope fixed; G1 next |
| **G1** | Track 1 fixture | source/role/attribution gaps visible; capture-incomplete distinct |
| **G2** | Track 2 corpus + rules baseline | labels frozen; R measured without BGE |
| **G3** | BGE candidate-ranking run | artifact pinned; margins/abstention measured; vetoes reported |
| **G4** | Track 3 lifecycle analysis | unresolved representation and retrieval consequences stated |
| **G5** | owner decision | class source/representation; no amendment before this |

**ACX-50 (Calculated — no stage skipping).** A good BGE score cannot excuse a capture-incomplete fixture. A good
classifier cannot excuse an unresolved lifecycle that violates Decision A.

**G1 result (2026-11-21).** The test-only fixture passed all F1–F10 and character-tell cases; production capture wiring
is still unestablished and requires separate authorization. G2 has not started. See `MEMORY-09-G1-CAPTURE-FIXTURE-RESULTS-2026-11-21.md`.

---

## 8. What would count as success

**ACX-51 (Calculated).**

- Track 1: every character memory has a stable source character and role; every missing required field is
  capture-incomplete, not guessed.
- Track 2: R or R+B reaches a predeclared coverage/accuracy level; BGE adds measurable value; abstention is calibrated;
  no veto fires.
- Track 3: unresolved class can be represented and retrieved without violating Decision A/L-2, **or** the owner
  explicitly amends them; the first resolution is distinct from a class change.

**ACX-52 (Calculated — falsification).** Capture cannot supply the envelope; BGE similarity is not discriminative or
is confidently wrong on ambient/restricted; rules cannot beat a trivial baseline; unresolved lifecycle cannot be made
compatible without an amendment the owner rejects. Any of these blocks the hybrid.

---

## 9. Decision requested

**ACX-53 (Observed).** Authorize the next investigation in this order:

1. Track 1 capture-completeness fixture (no model, no class labels).
2. Track 2 rules baseline, then the BGE embedding candidate-ranking experiment (model/corpus authorization separate).
3. Track 3 unresolved-class lifecycle analysis and owner decision point.

No L-2 amendment is requested now. No implementation, schema, migration, test or model fetch follows from this plan.

---

## 10. What this plan does not do

**ACX-54 (Observed).**

| Untouched | State |
|---|---|
| L-2 / Decision A | **unchanged** |
| P-1 / S-3 | closed; not reopened |
| Decision B | untouched |
| N-b / N-c / authority | untouched |
| schema / code / migration / tests | **0** |
| model fetch / BGE run | **0** |
| corpus creation | **0** |
| requirements / map | **0**; **118** |
| canonical capture | **0** |
| implementation plan | **not revised** |

---

## 10A. Amendment 1 — G0 acceptance and measurement additions

**ACX-57 (Observed — G0 accepted).** The owner accepts the plan and the G0–G5 sequencing. G0 is closed; G1 is next.
No L-2 amendment is authorised; no BGE model is fetched; G3 remains a separate authorisation.

**ACX-58 (Calculated — false confidence is a first-class metric).** Accuracy and coverage are not enough. The
experiment must measure **false confidence**: a confidently assigned **incorrect** class. Confidence is the
classifier's calibrated confidence or calibrated BGE margin, not a raw cosine similarity. Report it overall and per
class, with special attention to `ambient` and `restricted`. A classifier that abstains frequently may be preferable
to one that appears productive while silently assigning inappropriate eligibility.

**ACX-59 (Calculated — separate result streams).** Report four streams separately:

| Stream | What it is |
|---|---|
| **rules-only** | deterministic rules/metadata; no BGE |
| **BGE recommendations** | top-1/top-k and margins; candidate evidence only |
| **hybrid recommendations** | R+B where BGE ranks only after rules abstain; conflicts/low margin → unresolved |
| **human-established reference** | the frozen reference labels; not a classifier arm |

No single headline may merge them. Report where BGE changes a rules result, where the hybrid changes a BGE result,
and whether either change is correct, incorrect, or abstained.

**ACX-60 (Calculated — no aggregate hiding).** A high aggregate cannot hide false confidence, a veto, or a class with
poor coverage. Report abstention, false confidence, per-class results and vetoes before any aggregate.

**ACX-61 (Calculated — the character-tell scenario).** A character tells another character about their past. Capture
preserves the source, span and provenance. The memory is available for non-class-dependent processing (content,
ownership, status, display, inspection). Its class remains unresolved; class-sensitive audience/eligibility waits; the
conversation is not interrupted and no class is invented. This is an investigation scenario, not an authorisation. It
must be reconciled with Decision A and L-2 before any implementation.

**ACX-62 (Observed — no L-2 change, no BGE fetch).** This amendment adds measurement requirements only. L-2 and
Decision A remain unchanged; G1 is the next gate; G3 requires a separate model/corpus authorisation.

---

## 11. Traceability

**ACX-55 (Observed).**

| What | Where |
|---|---|
| Capture gap | CFD-06…10; ACX-05…09 |
| No App binding | `frontend/wailsjs/go/main/App.d.ts`; ACX-07 |
| Chat history fields | `ChatHistory`; ACX-09 |
| BGE reranker evaluation | `MEMORY-09-S3-BGE-EVALUATION-2026-11-15.md`; ACX-10 |
| Embedder envelope | `MEMORY-09-S3-EMBEDDER-RESOURCE-ENVELOPE-2026-11-15.md`; ACX-40 |
| Class vocabulary | K-R8; §5; ACX-27 |
| Ambient authored-only | K-R8; ACX-29 |
| Decision A | DAO-05; ACX-43 |
| L-2 | L2D-01/08/10; ACX-04/44/48 |
| First resolution not a change | CFD-31; ACX-45 |
| No experiment | 0/144; ACX-14 |
| Map | **118** |

---

## 12. Record status

**ACX-56 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — map **118** |
| L-2 / Decision A amended | **0** |
| Direction selected | **0** — plan only |
| Model fetched / BGE run | **0** |
| Tests executed | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **experiment plan** — no canonical capture |

**End of plan.** Capture comes before classification (ACX-01…04). The current backend path lacks speaker, stable source
character, world/scenario metadata and attribution, and is not exposed as an App binding; the frontend passes empty
history and identities (ACX-05…09). BGE is proposed only as an **embedding candidate-ranker**; the evaluated
`bge-reranker-v2-m3` is a different artifact and is not reused (ACX-10…12). Track 1 establishes the capture envelope
and marks capture-incomplete separately from unresolved (ACX-17…24). Track 2 measures rules, BGE and R+B under frozen
labels, calibration and vetoes (ACX-25…42). Track 3 analyses whether an unresolved class can satisfy Decision A/L-2
without amendment, and keeps the first resolution distinct from a class change (ACX-43…48). **No L-2 amendment, model
fetch, code, test, corpus or implementation follows** (ACX-54/56).
