# MEMORY-09 — Normative Contract (extract)

**Status:** normative extract of the design record, generated mechanically. Nothing here is
new, and nothing here is adopted beyond what the design record itself marks as decided.

- **Source:** `MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md`
- **Source revision:** r32
- **Source SHA-256:** `e88ef2d1697da1d904fd8791af0ccf1023e3bd577c84b26774a72b70b4ef598d`
- **Generator:** `docs/research/build_normative_extract.py` (verbatim copy; no paraphrase)
- **Rule count:** 90 · **Invariants:** 18 (counted from the
  invariant section; the rule index is counted from the whole source)

Every clause below is copied byte-for-byte from the source. Rationale, worked cases and the
normative tables remain in the design record and are referenced by section, not duplicated:
a second copy of the reasoning would create exactly the two-sources-of-truth problem this
contract exists to prevent (design record §13.1).

## 1. Adoption status

| Item | Status | Where |
|---|---|---|
| D1–D7 (ownership, existence≠knowledge, assets≠experience, knowledge-gated, decay, conflict, belief) | **adopted** | §2 of this document |
| D8 chat-deletion contract | proposed | design record §11.1; §2 here |
| Q13 transfer carries content, not knowledge | approved in principle | design record §11.2 |
| Materiality predicate revision 2 | accepted candidate, **not** approved | design record §10.1 |
| S24–S26 materiality probes | registered, **not runnable** | design record §12 |
| Schema versioning, evaluation provenance, replayability states | proposed | design record §10.2–10.2.2 |
| Temporal applicability conventions (§11.3) | proposed | design record §11.3 |
| Q1–Q15 | open | §4 of this document |
| E-SCOPE-CONTRACT suite | **not frozen** | design record §12 |

## 2. Decisions

Verbatim from the design record §5 (source lines 308+):

| Dimension | Mechanism | Notes |
|---|---|---|
| Ownership (chat) | field `owner_chat_id` | lifecycle/deletion boundary |
| Scope | field `scope` | applicability, not audience |
| Knowledge class | field `knowledge_class` | ambient / restricted / experiential / character_specific |
| Attribution | field `attribution{kind, by}` | observed / told / believed / asserted / authoritative |
| Audience | **derived** from edges + class + scope + context | never stored |
| Knowledge state | projection `KnowledgeEdge` | reconstructable from acquisition events only (I11, K-R11) |
| Knowledge history | append-only `AcquisitionEvent` | authoritative for knowledge |
| Memory lifecycle history | append-only memory events | authoritative for status/scope/class; status is projected, never a mutable field (K-R12) |
| Truth/lifecycle status | field `status` | provisional / active / superseded / retracted / corrected |
| Temporal provenance | fields `established_turn`, chain | |
| Confidence | field, separate from importance | |
| Importance | field | narrative salience |
| Surfacing priority | derived, decay-adjusted | never affects eligibility |
| Conflict state | field `conflict_group` | unresolved marker |
| Retrieval context | **call parameter** | character / world_state / narrator / inspection |

**Deliberate non-fields:** audience, `known_by`, and "current knowledge" are all derived.
Storing any of them alongside the event log creates the competing-source-of-truth failure
R1 exists to prevent.

D8 is proposed, not adopted; it appears in the design record §11.1 and in the open
questions listed in §4 below.

## 3. Invariants

Verbatim from the design record (source lines 283+):

- **I1** Eligibility precedes ranking. [D]
- **I2** Ranking cannot alter truth, membership or knowledge. [D, E §1.2]
- **I3** Status, ownership, scope, attribution, class, audience and knowledge are structured and out-of-band. [E §1.2]
- **I4** History is never destroyed; supersession and retraction are states. [D]
- **I5** Ownership, scope and knowledge are enforced at retrieval time. [E for owner; D for scope/knowledge]
- **I6** Empty eligible set → abstain; never widen the pool. [D]
- **I7** Retrievability does not confer authority. [X]
- **I8** Application owns conflict detection; the model owns the clarification request. [X]
- **I9** Memory existence does not imply character knowledge; presence does not imply observation. [X]
- **I10** Asset content may be context, never automatically experience. [X]
- **I11** Acquisition history is authoritative; knowledge edges are its projection; any edge must be reconstructable from events. [X, R1]
- **I12** Retrieval context is explicit and never inferred; the character channel is the default. [X, R2]
- **I13** Knowledge class is authored, not inferred from popularity or age. [X, R3]
- **I14** The edge is a summary of current knowledge; the event log is the account of how it was acquired. Questions about the account are never answered from the summary. [X, R1d]
- **I15** **Absence of information is never treated as information about absence.** An unpopulated field, an empty list, a missing record and an unwritten schema each mean "not established", never "established as none". A system that cannot safely determine an answer says so instead of manufacturing certainty. [X]

  **Applies uniformly to:** missing acquisition history, empty knowledge sets, unspecified permissions, undeclared action requirements, missing asset versions, and incomplete retrieval evidence. "Known to be absent" and "not established" are distinct states in every one of them.
- **I16** **Conflict determines whether claims disagree; temporal applicability determines which claims apply to a given point or interval.** Neither layer answers the other's question: a system that conflates them reads ordinary world evolution as contradiction. [X, R1d]
- **I17** **A result earns authority from the conditions under which it was established, not from when or where it was written.** Applies beyond temporal comparison to retrieval, formation and conflict evaluation: never "it was computed most recently", always "it was computed under these conditions". [X]
- **I18** **The system may not use absence of evidence as evidence that nothing happened.** Every subsystem must name the **independent evidence** that would reveal its own silent omission — evidence it does not produce itself. A subsystem that cannot answer "what would show I failed to act?" is incomplete, not merely untested. [X]

## 4. Rules

Verbatim, in document order, grouped by their source section. Continuation lines and fenced
blocks belonging to a rule are included; normative tables that follow a rule are not copied
and must be read in the source section named by the group heading.

### R1 — Acquisition history is authoritative; the edge is its projection [X]

- **K-R1** The edge for (character, memory) is the fold of that pair's acquisition events in turn order. No edge exists without at least one event.

- **K-R2** An edge carries `derived_from_event_ids`. A projection that cannot be reconstructed from its events is an **integrity failure**, not a policy question.

- **K-R3** Revocation is a new event (`via=revoked`), never a deleted history row. Forgetting changes current knowledge; it never rewrites what happened.

- **K-R4** Where provenance matters (assembly, audit), the *event chain* is rendered, not the edge alone: "told by B at turn 210, unverified" is history; "knows it" is state.

### R1b — Projection purity: what may enter the edge fold [X]

- **K-R11** The edge fold depends **only** on acquisition-family events for (character, memory): `witnessed | told | read | inferred | verified | narrator_grant | revoked`. Nothing else may influence it.

- **K-R12** Memory lifecycle is its **own authoritative stream** — `established | superseded | corrected | retracted`, plus `scope` and `knowledge_class` changes — and a memory's current status, class and scope are a **projection of that stream**, not mutable authoritative fields.

- **K-R13** "This character's knowledge includes a retracted claim" is a **read-time join** of edge × memory projection. It is never stored on the edge.

- **K-R14** Each stream is rebuilt and checked independently; a mismatch in either is an integrity failure. The cross-stream join has no stored form, so it cannot drift.

### R1c — The acquisition fold, specified [X]

- **K-R15** Acquisition events for a (character, memory) pair are totally ordered by `(turn, seq)`, where `seq` is a monotonic per-chat counter. **No wall clock and no insertion order participates.** Events sharing a turn are resolved by `seq` alone.

- **K-R16** Event ids are unique; a duplicate id is rejected at write.

- **K-R17** `revoked` is valid only when the pair's ordered log currently ends in a learn event. A revoke with no open episode is **rejected at write**, not silently ignored, because a silent no-op hides formation bugs.

- **K-R18** A learn event is valid at any time. A learn while an episode is open is a *continuation* (the character learns again / confirms), not a new episode.

- **K-R19** A revoked pair has **no edge**, not an edge with a "revoked" flag. "Forgot" is the absence of the current relation; the fact that it was once known lives in the log, not in a half-state row.

- **K-R20** `revoked_turn` is therefore **not a field**. It is derived on demand — the timestamp of the last revoke event — and appears only in the `knowledge_history` inspection view.

### R1d — Summary versus account; continuation; recall [X]

- **K-R21** `via` reports the strongest evidence in the current episode. It does **not** describe the full history: a character told a rumour and later witnessing the same event has `via=witnessed`, while both events remain in the log.

- **K-R90** Any question of the form "how did they learn this?" is answered from **acquisition events**, never from `via`. Rendering and audit read the log; eligibility reads the edge.

- **K-R22** A learn event while an episode is open may raise `via` along the precedence order and may set `verified`, but **must not change `acquired_turn`**; the episode's start is fixed by its opening event.

- **K-R23** Within an episode, `via` never downgrades. A weaker later event (e.g. `told` after `witnessed`) is appended to the log and does not lower the edge's summary.

- **K-R26** Precedence is a **fixed total order over kinds**: `witnessed > verified > read > told > inferred > narrator_grant`. Comparison is on kind only, so two events can never tie on `via` — no alphabetical, insertion, locale or hashing order participates in the result.

- **K-R24** Recall difficulty is a **surfacing** property, never an eligibility or truth property. It must not hide a record from a direct question, and it must never be represented by removing an edge.

- **K-R25** Forgetting (no edge) and difficulty recalling (edge, low surfacing) are distinguishable in every view that reports either.

### R2 — Retrieval is context-explicit [X]

- **K-R5** Every retrieval call names its context. **Default is `character_perspective`** (fail closed).

- **K-R6** A scene in which a character discovers something requires the `narrator` context to know what the character does not — i.e. R2 exists so that knowledge gating does not make discovery impossible.

- **K-R7** Moving information from a broader context into the character channel requires an `AcquisitionEvent`, never a context switch alone.

### R3 — Knowledge classes, not a single flag [X]

- **K-R8** `knowledge_class` is a field on memories **and** assets. Ambient applies only to authored world/scenario definitions.

- **K-R9** No record becomes ambient because many characters know it. Popularity is not authorship. **[D]** (needs review)

- **K-R10** An experiential event never downgrades to ambient by repetition or by age; decay does not change class (D5).

### 10.1 Candidate mechanical predicate, revision 2 — proposed, awaiting approval [P]

- **K-R27** Absence of a precondition list is not a claim that no preconditions exist. Only an explicit `declared_none` may be treated as "nothing required"; `undeclared` must raise an evaluability error rather than silently yielding "immaterial".

- **K-R28** `undeclared` is a property of the action schema, not of any memory record. It is declared once, reviewed once, and does not migrate into per-query state.

### 10.2 Schema versioning for reproducible evaluations — proposed, not adopted [P]

- **K-R29** Action-schema declarations are **immutable once published**; any change creates a new version. This mirrors D3: a schema is a reusable asset, not chat memory.

- **K-R30** Every materiality evaluation records `schema_version` alongside its inputs and verdict. Replay uses the recorded version, never the current one.

- **K-R31** Re-evaluating the same query under a newer schema is a **new** decision with its own record; it never overwrites or reinterprets the earlier one.

### 10.2.1 What an evaluation record must carry — proposed, not adopted [P]

- **K-R32** Replay must verify **every** recorded governing version. If any cannot be matched, the replay **fails loudly**; it may not emit a verdict from partial agreement.

- **K-R33** The evaluation log is its own append-only artifact. It adds no fields to memory records and no new state to the knowledge model — it is a record *of* decisions, not part of them.

### 10.2.2 Replayability states — proposed, not adopted [P]

- **K-R34** A replay result carries exactly one of these states. `auditable_only` and `expired` **must never** be reported as a successful reproduction.

- **K-R35** Deleting an artifact referenced by a live evaluation record requires an **explicit, recorded retention action**; a tombstone names the artefact and the policy decision that authorised removal. Silent garbage collection may not achieve it.

- **K-R36** Expiry affects **replayability only**. A verdict does not become false, withdrawn or revised because its evidence was lawfully deleted — and an `expired` record may never be quietly dropped from the log.

- **K-R37** **A hash proves integrity, not truth.** `auditable_only` claims exactly two things: the archived content matches what was hashed, and the chain of identifiers is traceable. It never claims the original evidence was accurate, complete or sufficient, and no downstream decision may cite it as evidence of correctness.

- **K-R38** **Retention is atomic from the audit's point of view.** For any observer at any instant, the state "artifact absent and no tombstone" must be unreachable. Ordering is therefore tombstone-first:

### 11.1 Chat deletion contract — proposed, not adopted [P]

- **K-R39** Ownership is scope-independent: every experiential memory of the chat is deleted with it, whatever its scope. `scope=world` does not exempt a memory from the chat's lifecycle. [links D1]

- **K-R40** **No automatic transfer.** Memories never silently migrate to another chat, world, scenario or character. Manual transfer is the only preservation mechanism, and it names both source and destination.

- **K-R41** **Warnings are computed, never estimated.** Any deletion prompt reports counts and categories enumerated from the actual deletion set at confirmation time; a stale or approximated count is a defect, not a rounding issue.

- **K-R42** **Explicit confirmation is required**, and the text states what is deleted, what is retained, and that the action cannot be undone. Example: *"Deleting this chat will permanently delete its associated world, scenario and character memories that have not been manually transferred. Reusable character, world and scenario assets will remain. This action cannot be undone."*

- **K-R43** Deletion removes the experiential records **and their dependent knowledge relations and acquisition history**. Dangling edges and orphaned events are forbidden. Deletion is recorded as a retention action under K-R35, so evaluations referencing those artifacts transition to `expired` by the same tombstone protocol — never to a silent gap.

- **K-R44** Reusable assets and transferred destinations are untouched. Deleting a chat must not delete a lorebook, a persona, or a record that now lives in another chat.

### 11.2 Manual transfer — Q13 approved in principle [X]

- **K-R45** Transfer is an explicit event naming actor, source chat, source memory and destination. The destination record is marked `transferred_from` and its origin fields are **copied references, never rewritten**.

- **K-R46** No acquisition events or edges are copied. Any later knowledge in the destination is a **new** event with its own `via` (`read | told | narrator_grant`), attributable to that chat.

- **K-R47** The chain `original memory → transfer event → destination memory` must remain auditable even if the source chat is later deleted: the transfer record carries a **content hash and provenance stub**, so the destination's history does not dissolve into an unverifiable claim. Under **K-R37**, this hash proves *which* content was transferred and *what provenance was recorded* — never that the original memory was accurate, complete or still available.

- **K-R48** Both relations are **additive only**. No transfer may delete, merge, overwrite or re-date either record, and both provenance chains survive association.

- **K-R49** Equivalence that cannot be mechanically decided is recorded as **`unresolved_identity`** — an explicit state requiring a normalisable value or inspection. It may never default silently to "distinct", consistent with **I15**.

- **K-R50** Incompatibility is evaluated over **temporal applicability**, not over a bare subject/predicate/scope key. Two values conflict only where their applicability intervals **intersect**; an empty intersection is temporal progression, never contradiction.

- **K-R51** A **partial** overlap scopes the conflict to the intersecting interval. Outside that window the records are independent facts; the conflict may be neither widened to the whole span nor narrowed to nothing.

- **K-R52** Disjoint intervals are recorded as **succession** — ordinary story progression. Succession is not contradiction and not `unresolved_identity`. Which value applies *now* is the temporal layer's question (D5), not the conflict layer's.

### 11.3 Temporal applicability conventions — proposed, not adopted [P]

- **K-R53** Intervals are **half-open**, `[start, end)`. The instant a segment begins belongs to that segment, so adjacent intervals never overlap and never leave a gap.

- **K-R54** An interval is one of four tiers, and the tiers are not interchangeable:

- **K-R55** An unknown endpoint is **not** infinity. Unknown start does not mean "since the beginning of time"; unknown end does not mean "still current". A `none` interval yields **`unresolved_temporal_applicability`** — distinct from `unresolved_identity` (K-R49), because the uncertainty is *when*, not *which record this is*. Reusing the identity state would conflate two different questions in one field.

- **K-R56** Comparing two intervals yields exactly one of three results: **`disjoint`** (evidence establishes no overlap), **`overlapping`** (evidence establishes an intersection), or **`indeterminate`** (evidence establishes neither). `partial` and `none` tiers can produce any of the three depending on what the known bounds prove.

- **K-R57** The result is computed from **provable** interval relations over known bounds, never from the tier label. For `A = [a₁, a₂)` and `B = [b₁, b₂)`:

- **K-R58** **Indeterminate is never coerced.** It is not conflict and not succession: it opens no conflict group, it also permits no "no conflict" conclusion, and a material query over an indeterminate pair must gather temporal evidence or abstain — never guess a value.

- **K-R59** Every evaluated pair records its outcome, and four states are distinct:

- **K-R60** An absent comparison may **never** be read as established non-conflict. `not_compared` is recorded explicitly; a consumer needing the relation must either perform the comparison or abstain. This is E3's rule applied to comparison: a retrieval stage must not mistake its own omission for proven agreement.

- **K-R61** `compared: indeterminate` is persisted, not discarded when no conflict group forms. Dropping it would make the pair indistinguishable from one never examined — the exact collapse K-R58 forbids.

- **K-R62** A comparison is bound to **the evidence it was computed from**: interval tiers, known bounds, and derivation basis. That input fingerprint is recorded with the result. A comparison is a judgement about specific evidence, not a permanent property of the pair.

- **K-R63** When the underlying evidence changes — a missing endpoint established, a `partial` becoming `explicit` — the existing comparison becomes **`stale`** and is never silently authoritative. A consumer must recompute before relying on the relation.

- **K-R64** Re-evaluation **supersedes without erasing**: the new comparison is authoritative, the previous one is retained as history of what was knowable at the time. `superseded` and `stale` are distinct — stale means no replacement exists yet and the relation cannot be answered from old results.

- **K-R65** **Commit is a compare-and-set on the input fingerprint.** A recomputed result may become authoritative only if the evidence fingerprint still matches the current inputs at commit time. If evidence changed while the comparison was being computed, the result is **discarded**, never published as current.

- **K-R66** A discarded attempt is logged with its own fingerprint and is distinguishable from `stale` and `superseded`; it never answers a query. Discarding is not a failure of the computation — the computation was correct for evidence that no longer applies.

- **K-R67** **No last-writer-wins.** Authority derives from fingerprint match, never from commit recency; a later commit against older evidence cannot displace an authoritative result computed against current evidence.

- **K-R68** The fingerprint is **content-addressed over the complete relevant input state**, not a version counter: interval bounds and tiers for both records, derivation bases, the governing policy versions, and the **implementation identity** of the comparison itself. A version number that returns to an earlier value must produce the same fingerprint only if the content genuinely matches.

- **K-R69** A prior result is **reusable** exactly when the full fingerprint matches: same evidence content, same governing policies, same implementation identity. Then the commit check may accept it without recomputation. Match on evidence alone is insufficient.

- **K-R70** Same evidence under a **different implementation identity** or policy version is **not reusable** — the earlier result is retained as history under its own fingerprint and the pair is recomputed. Reuse never crosses a change in the rules that produced the answer.

- **K-R71** Fingerprint inputs are serialised **canonically**: fixed field order, fixed encoding, normalised numeric representation, sets sorted, and no locale-, map-iteration- or serialiser-dependent ordering. Equivalent inputs must hash identically regardless of how the fields happened to be assembled.

- **K-R72** Canonicalisation preserves **semantic distinctions**: absent, empty and null are three different states and must encode differently. Collapsing them would make `not established` and `established as none` hash alike — the I15 failure at the representation layer.

- **K-R73** Canonicalisation is part of **implementation identity** (K-R70). Changing the canonical form invalidates reuse exactly as changing the comparison does, even if the computation is untouched. Input that cannot be encoded canonically is an **error**, never a silent fallback to a looser encoding.

- **K-R74** **Numeric domain is declared, not inferred.** Within the declared domain, numerically equal values encode identically (`1`, `1.0` and `1e0` are one value); values outside it — NaN, ±Infinity, and any undeclared sentinel — are **rejected**, not encoded.

- **K-R75** **Text is UTF-8 and Unicode-normalised to a declared form** (proposed: NFC). Case is preserved: no implicit case folding, no locale-dependent collation. Normalisation is applied explicitly by the protocol, never left to the runtime.

- **K-R76** **Ordered sequences and unordered sets are distinct types** and encode differently. Sets sort by their **canonically encoded bytes**, never by locale-sensitive text comparison.

- **K-R77** The fingerprint protocol records its **digest algorithm and canonicalisation version**. A hash migration or canonical-form change is a new protocol version and **invalidates reuse** — old fingerprints are history, not authority.

- **K-R78** **No partial fingerprints.** If serialisation fails partway, no fingerprint is produced and no comparison result may become authoritative; a truncated digest is never a usable key.

- **K-R79** **Implementation identity is a declared identifier with its own version**, not a source-code hash. Irrelevant edits must not invalidate reuse, and meaningful changes must not silently preserve it — the identity is asserted deliberately and recorded with the result.

- **K-R80** **Three identities must all match for reuse**, and they answer different questions:

- **K-R81** **Identity change is a release-governance act, not a refactor side effect.** A material change in comparison behaviour requires a reviewed, recorded identity increment with its rationale; the identity registry is the record. An identifier that no process is required to move becomes dangerously stable.

- **K-R82** Behaviour-affecting change without a corresponding identity increment is a **defect**, not a judgement call. It is detected and reported, never silently accepted — otherwise reuse authority would survive a change nobody noticed making.

- **K-R83** **Identity completeness is established outside the registry.** Detecting an unincremented change requires independent evidence that behaviour changed — controlled behaviour-change fixtures and reviewed change declarations — because a registry cannot prove its own completeness. A missing entry is not evidence that nothing changed (I15, applied to provenance).

- **K-R84** Every subsystem specification states an **omission detector**: the independent evidence that would expose a failure to act. Omitting that statement is an incompleteness defect at review, not a documentation gap. [links I18]

- **K-R85** **Detector coverage is bound to the subsystem contract it was written against.** A change to a subsystem's responsibilities triggers re-evaluation of its detector; coverage that no longer matches the current contract is **stale coverage** and is treated as no detector at all until re-established. This reuses K-R63's staleness model rather than inventing a second one.

- **K-R86** Detector validity is **demonstrated, not asserted**: it must have fired on a seeded omission for the **current** contract, and must not fire on a no-op. A detector validated against a superseded contract is stale coverage, however impressive its history.

- **K-R87** **A detector has its own implementation identity** (K-R79). Changing how a detector works — even with the subsystem contract and its coverage label untouched — invalidates that detector's demonstrated validity until it is re-seeded. A current label over a changed mechanism is stale coverage by another route.

- **K-R88** Detector implementation identity includes its **material dependency closure** — the parsers, fixture runners, evidence sources and downstream libraries that can affect either detection behaviour or the validity of the evidence it consumes. The criterion is capability, not change: a dependency belongs in the closure if it *can* influence the outcome, regardless of whether it moved.

- **K-R89** The material set is **declared and reviewed** as part of the detector specification. A change to any declared material dependency invalidates demonstrated validity until re-seeded; updates outside the declared set do not. Discovering an undeclared material dependency is itself a defect (K-R83's rule, applied to detectors): the list cannot be certified by the list.

## 5. Open questions

Verbatim from the design record §14 (source lines 1000+):

**All items below remain open and provisional.** Accepting the r4 review corrections into
this revision does not approve the architecture those corrections now describe.

Awaiting approval (the seven directions in §3.5 are proposals):
1. **Q1 materiality** — approve a mechanism, or keep S24–S26 blocked?
2. **Q2 edge creation** — explicit acquisition events only; is narrator-granted knowledge an acceptable `via`, and how is it marked?
3. **Q3 crowds** — group-level knowledge first: who owns a group, and when does a member individualise?
4. **Q4 forgetting** — allow explicit revocation events (representable today), or is forgetting surfacing-only?
5. **Q5 asset versioning** — pin each chat to the asset version it started with?
6. **Q6 retraction** — R4 recorded as direction; confirm the "A may still believe" clause.
7. **Q7 ambient facts** — approve the four classes, and decide who may author an `ambient_authored` record.

New in r4:
8. **Group knowledge identity** (Q3 follow-on): is a group a pseudo-character for edges, or a scope-like audience? Affects whether S-probes need a group actor.
9. **Ambient and beliefs:** may a character hold a belief that contradicts an ambient authored fact (a flat-earther in a spherical world)? Proposed: yes, as `character_specific`, never mutating the ambient record.
10. **Inspection exposure:** does `inspection` show the full knowledge map to the player, or is some of it withheld to preserve mystery?
11. **Schema versioning (§10.2) and asset pinning (Q5):** decide these together — immutable published schemas with replay against the recorded version, or mutable schemas with evaluations treated as point-in-time only?
12. **Evaluation-log retention and minimisation (§10.2.1–10.2.2):** keep full input snapshots for exactness, or identifiers plus watermarks at lower storage cost? And crucially: what is the authorised path for deleting an artifact a live evaluation record references, so that lawful retention never gets reported as provenance failure? Outcomes must be distinguished as `reproducible`, `auditable_only` or `expired` — never collapsed into "replay worked" or "replay broke". Also: approve the retention-action ordering protocol (K-R38) and decide who may authorise deletion of artifacts referenced by live evaluation records.
13. **Deletion contract (D8, §11.1):** approve the chat-deletion lifecycle. **Q13 is approved in principle** (§11.2): content and provenance move, knowledge does not.
14. **Q14a equivalence on transfer:** what declared normalisation decides that two records assert the same fact (making them `possible_duplicate` rather than in conflict)?
15. **Q14b incompatibility on transfer:** what applicability window (time, context) decides that competing values cannot both be true — and therefore open a conflict group? The interval algebra is now fixed (K-R50–K-R52): empty intersection means succession, identical instant means conflict, partial overlap is scoped to the intersection. Open: how an interval is derived for a record that was never explicitly dated. §11.3 proposes four tiers (`explicit`, `derived`, `partial`, `none`), half-open boundaries and a distinct `unresolved_temporal_applicability` state, all pending approval — including which event-ordering sources qualify as trustworthy enough to licence a `derived` interval. Neither merge nor overwrite is acceptable in any case.

## 6. Measured baseline at the time of this extract

Established by running the existing apparatus unchanged; no measurement was made for this
document. Receipts live in `eval/s3/`.

| Measurement | Result |
|---|---|
| Gate A | 14/14 PASS |
| Gate B | **6/7 PARTIAL** — P5 ("Who trained the horse?") still open |
| Formation | frozen at Gate A baseline (7 records); hashes match |
| Harness self-checks | PASS |
| S1 audit | PASS (5 calls, 0 problems) |
| `validate` | 89/89 |
| Corpus | PASS |
| Reranker (bge-reranker-v2-m3 Q8_0) | evaluated, **not shipped**; no further candidate fetched |
| Resource envelope | no configuration meets VRAM ≤384 MiB **and** p50 ≤40 ms (87/198/277/444 MiB; 282/102/70/20 ms) |

Gate B 6/7 is **not** a completed baseline: P5 is an open failure, and the proposed
architecture may not inherit it as an assumption.

## 7. What this extract does not claim

- It is **not** behaviourally validated. No probe in the registry has been run.
- It is **not** a suite freeze. The registry is not frozen and cannot be until an independent
  adversarial review of its expectations has taken place.
- It does **not** authorise implementation, or integration into Howl Chat.

