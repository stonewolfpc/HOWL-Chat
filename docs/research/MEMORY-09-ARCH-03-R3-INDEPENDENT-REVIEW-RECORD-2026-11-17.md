# MEMORY-09 — Independent Architecture Review Record (ARCH-03/r3)

> ## ⚠ RETARGETED — THIS IS NOT A REVIEW OF r3
>
> **Original target:** ARCH-03/**r3**, SHA-256 prefix `f8bbe6…`. It was not present in the
> workspace and was not recovered by a bounded read-only investigation, now closed in
> `MEMORY-09-R3-RECOVERY-CLOSURE-2026-11-17.md`.
>
> **Replacement target:** **r32**, SHA-256
> `e88ef2d1697da1d904fd8791af0ccf1023e3bd577c84b26774a72b70b4ef598d`, 91,311 bytes.
>
> **On 2026-11-17 the owner authorized retargeting this review from ARCH-03/r3 to r32** (§1.4, §5).
> Every finding in this record is attributable to examination of **r32 only**. No finding is
> evidence of a defect in r3, and no finding establishes r3's exact contents.
>
> The filename retains `ARCH-03/r3` for continuity with the work order. It is not a claim about the
> reviewed target.

---

## 0. Reviewer identity and review time

| Field | Value |
|---|---|
| **Agent role** | Buffy (Codebuff coding agent) |
| **Model name (as supplied by the owner)** | **DeepSeek V4.1 Flash** |
| **Model identity basis** | **Artifact-reported, not independently verifiable.** The reviewer cannot read its own serving model from this environment. This matches the project's standing convention that model labels are *provenance data, not lineage proof* (`MEMORY-09-OWNER-DECISION-RECORD-OD-M09` §0; `CORPUS-69` §0). If the served model differs from the label above, the label is wrong and this record inherits that error |
| **Review date** | **2026-11-17** (UTC) / 2026-11-16 (local) |
| **Review time** | **02:23:34 UTC** = **2026-11-16 20:23:34 CST (−0600)** |
| **Review location** | `C:\Users\StoneWolfPC\HOWL_Chat` (win32) |
| **Review type** | Independent, adversarial, **read-only**, document-level |
| **Authorization** | Owner work order `MEMORY-09` — independent architecture review of ARCH-03/r3 |

The filename carries the UTC date because that is the session date. Both stamps are given above
so the record cannot be misdated by timezone.

---

## 1. Target identity — the finding that governs everything else

| Item | Requested | Found | Result |
|---|---|---|---|
| Path | `docs/research/MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md` | same path | **match** |
| Revision | **r3** | **r32** | **MISMATCH** |
| SHA-256 | prefix `f8bbe6…` | `e88ef2d1697da1d904fd8791af0ccf1023e3bd577c84b26774a72b70b4ef598d` | **MISMATCH** |
| Size / lines | — | 91,311 bytes / 1,028 lines | recorded |
| Git history for the file | — | **none** — untracked (`??`), `git log --` returns nothing | **no recoverable revision** |
| Archive / sidecar copy of r3 | — | **not found** | **unavailable** |

**1.1 Where `f8bbe6…` does occur.** Twice, only inside the r32 document's own header changelog
(lines 34 and 36), as the fragment `r3 (f8bbe6aa…) introduced the edge model`. A self-referential
changelog line is an **assertion by the document**, not evidence of an artifact. No file in this
workspace can be hashed to that prefix.

**1.2 The on-disk file is the r32 that was formally preserved.** `MEMORY-09-PREREGISTRATION-FREEZE-2026-11-15.md`
§3 records this path at `e88ef2d1…` as "design contract r32". The computed digest matches, so the
surviving artifact has not drifted since that freeze. This is the one integrity anchor available.

**1.3 Corrected wording — a self-correction of the first report.** The review initially described
r32 as containing "r3's edge model plus 29 later revisions". **That containment claim is not
established.** What is established is that *r32 asserts* r3 introduced the edge model. Whether any
given rule in r32 is identical to r3's text, was amended by a later revision, or was introduced
after r3, **cannot be determined without r3**. Every finding below is a finding about r32's text
only.

**1.4 Target change — owner decision, 2026-11-17.** The requested deliverable — an independent
architecture review of ARCH-03/r3 — **could not be produced**, because r3 was unavailable: the file
is untracked, has no Git history, carries no sidecar, and no artifact in the workspace or its
object store reproduces `f8bbe6…`. A bounded recovery investigation was opened, conducted
read-only, and closed as unsuccessful within its documented scope
(`MEMORY-09-R3-RECOVERY-CLOSURE-2026-11-17.md`).

On 2026-11-17 the owner **authorized changing the review target from ARCH-03/r3 to r32**, identified
by full digest `e88ef2d1697da1d904fd8791af0ccf1023e3bd577c84b26774a72b70b4ef598d` (91,311 bytes),
the revision formally preserved under that identity in
`MEMORY-09-PREREGISTRATION-FREEZE-2026-11-15.md` §3.

Three statements bind every use of this record:

1. All findings below are findings about **r32's text**, which is what was actually examined.
2. **No finding is evidence of a defect in r3**, and **no finding establishes r3's exact contents**.
3. Retargeting **validates nothing**: F2, F3, F5, F6, F7 and F15 remain open, and none may be closed
   by this change.

What follows is therefore a review of r32 under the requested protocol, and is citable as such.

---

## 2. Accepted owner qualifications and reclassifications

The owner's post-review assessment raised three qualifications. All three are accepted. Two of them
**change the severity classification** of findings in the original report; the third is a boundary
it already observed and is restated here as binding on this record.

### Qualification A — missing approval records are not proof that approval never happened

Accepted, and the original report's framing is corrected. Three statements must be kept separate:

| # | Statement | Status in this review |
|---|---|---|
| **A1** | The owner decision exists in the project's established decision history | **Not contested.** This review found no basis to doubt it |
| **A2** | Its authoritative record is missing from the inspected workspace | **Established** (Observed) |
| **A3** | The contract provides no traceable reference to it | **Established** (Observed) — ARCH-STAGE-CONTRACT §2, lines 67–73, cites no record, date or instrument |

**Reclassification.** The original report treated "D1–D7 has no locatable authority" as a *Critical
contradiction of an approved decision*. That was too strong. It is a **traceability defect**: the
contract asserts owner decisions without a citation, and the repository cannot supply one. The
decisions themselves are not thereby invalidated, and this review does **not** find that approval
never occurred.

**Retained, unchanged:** the identifier collision. Four unrelated sets use the label `D1–D7`
(`CORPUS-69` = corpus/registry/F3-01 decisions; ARCH-STAGE-CONTRACT §2 = architecture decisions;
`docs/review/AUDIT-howl-tts-stage3-corrective-design-2026-11-07.md` §6; `RECON-t2.9.5b-deepseek-muse-audio.md` §4),
plus `MEMORY-09-G0-RECONCILIATION` §1 using `D1`/`D2` for temporal provenance labels and `Q-MF-01`.
This is a live navigation hazard independent of whether approval occurred.

### Qualification B — chat deletion is outside MEMORY-09's architectural ownership

Accepted, and binding on how the finding may be resolved.

The established boundary is that **Howl's state machine owns chat lifecycle and deletion.**
MEMORY-09 must not independently design chat-deletion memory destruction.

**Reclassification.** The `I4` vs `K-R43` finding is a **real internal inconsistency inside r32**
and separately a **scope/boundary defect** — it is *not* an invitation for MEMORY-09 to decide
deletion semantics. It may only be resolved by deferring to the subsystem that owns the lifecycle,
or by removing the deletion-design content from this contract. `K-R43`'s protocol detail
(tombstone ordering, retention actions, `expired` transitions) is itself open to the same objection:
it specifies deletion behaviour for an artifact this document does not own.

### Qualification C — contradictions are identified, runtime failures are not

Accepted and unchanged. **No probe was executed; no implementation was inspected.** Nothing in this
record is a demonstrated software defect. These are document-level findings about a design record.

---

## 3. Findings register (r32)

Severity key: **C** = architecture-breaking at document level · **H** = likely to produce incorrect
ownership/knowledge/truth/retrieval if built as written · **M** = underspecified, permits divergent
implementations · **L** = documentation/terminology.

### 3.1 Architecture-breaking at document level

| ID | Finding | Evidence | Basis | Classification after qualification |
|---|---|---|---|---|
| **F1** | **Target revision unavailable.** r3 is absent; attested only by an 8-hex prefix in a changelog that duplicates itself with divergent wording. This review cannot be attributed to r3 | §1 above | Observed | **Critical — target compliance failed** |
| **F2** | **S20 contradicts K-R46.** S20 expects the same character to retrieve a memory *after a valid transfer*; K-R46 requires that no events or edges are copied and that later knowledge be a **new** event. S54/S56 agree with K-R46 | §12 v1.1/v1.2 probe table; §11.2 | Observed | **Critical — unresolved** |
| **F3** | **Five-condition materiality replaced by a three-condition predicate.** §10 invokes "D6's five conditions"; §10.1 defines (1) eligibility, (2) dependency (a)/(b)/(c), (3) unresolved. Condition 4 — which the work order says remains undefined — is not identifiable in §10.1 at all. S25 tests "all three" | §2 D6; §10; §10.1; S25 | Observed | **Critical — unresolved.** Condition 4 was **not** invented by this review |
| **F4** | **I15 vs empty acquisition log.** §6.4 calls a missing edge "a determinate consequence of the last event" — true when a log exists and ends in revoke. For a pair with **no events**, the system cannot distinguish "never learned" from "event never written", yet treats it as established ignorance. I15 forbids exactly that | §4 I15; §6.4 | Source-derived | **High — unresolved** |

### 3.2 High — likely to produce incorrect behaviour if built as written

| ID | Finding | Evidence | Basis |
|---|---|---|---|
| **F5** | **Status filtering can suppress retraction display.** §4's eligibility stage excludes by `status` "for the named context" with no character-channel exception; §8.4 requires the knower to be shown "you were told X; it has since been retracted". A strict reading of §4 excludes the record before the character channel sees it | §4; §8.4; S34; S3 | Source-derived |
| **F6** | **`AcquisitionEvent` has no `seq`** (and no chat field). §6.2 lists `id, character_id, memory_id, turn, via, from_character_id, source_memory_id, confidence_in_source, evidence_span`. K-R15 orders by per-chat `(turn, seq)`, and S38 asserts same-turn resolution "by `seq`". Same-turn ordering is **not representable** as specified | §6.2 vs K-R15/K-R16, S38 | Observed |
| **F7** | **Restricted-authored assets may be ungated.** K-R8 puts `knowledge_class` on assets; edges are `(character, memory)` and `AcquisitionEvent.memory_id` is memory-only. A `restricted_authored` *asset* has no acquisition mechanism: either it leaks into `character_perspective` or no character can ever obtain one | K-R8; §6.2; §6.3; R2 | Source-derived |
| **F8** | **`status` is both field and projection.** K-R12: status/class/scope are "a projection of that stream, not mutable authoritative fields". §5 lists `status` as "field `status`"; §6.1 includes it in `MemoryRecord`; §4 Formation lists "assign … status". Two authorities; no projection check specified for the field form | K-R12; §5; §6.1; §4 | Observed |
| **F9** | **`conflict_group` is stored derived state with no stream and no rebuild check** — the only derived value exempted from §3.6's own rubric, and a second authority beside §10's detection | §6.1; §5; §3.6 | Source-derived |
| **F10** | **Open questions pre-answered by adopted rules.** Q2 asks whether `narrator_grant` is an acceptable `via` while K-R11 `[X]` already lists it; Q4 asks whether revocation is event-based while K-R3/K-R19/K-R20 `[X]` decide it | §14 Q2/Q4; §3.5 | Observed |
| **F11** | **S29 contradicts ambient access.** S29: no edge ⇒ abstain. R3/S32/§8.2: ambient records need no edge. Two probes disagree about the same filter | §12 S29, S32; R3 | Observed |
| **F12** | **S8's mandatory gate rests on `[D]` invariants.** I1 and I2 are tagged "design detail proposed here, not yet reviewed", yet §12 makes S8 mandatory and §15 gates reranker reconsideration on S8 passing | §4 I1/I2; §12; §15 | Observed |
| **F13** | **T10's frozen expectation is the defect under repair.** T10 expects `abstain` on the world-scope fact; S1 expects reach. No rule states T10 is a baseline the contract intends to invalidate | §1; §12 S1; `eval/s3/temporal-spec.json` | Source-derived |
| **F14** | **Three provenance representations, no stated relation:** `AcquisitionEvent.via`, `memory.attribution{kind, by}`, `KnowledgeEdge.via`; and C4's rule keys on an edge field while "provisional" is a memory field | §6.2; §6.3; §5; §9 C4 | Source-derived |

### 3.3 The I4 / K-R43 finding (boundary-reclassified)

| ID | Finding | Classification |
|---|---|---|
| **F15** | I4 states "History is never destroyed"; K-R43 states deletion removes "their dependent knowledge relations and **acquisition history**". The two cannot both hold. **This is a genuine internal inconsistency in r32** | **Boundary defect + internal inconsistency — not a MEMORY-09 decision.** Per Qualification B, it may not be resolved by this document designing deletion semantics. Resolution paths are (i) defer the deletion contract to the owning subsystem, or (ii) remove the deletion-design content from this contract — an owner choice, not a reviewer choice. `K-R43`'s tombstone/retention detail is open to the same objection |

### 3.4 Medium

| ID | Finding |
|---|---|
| **F16** | `via` has a **7-value domain** on events (incl. `revoked`) and a **6-value precedence order** on edges (excl. it), with no statement that `revoked` is out of band. K-R26's "two events can never tie on `via`" covers only the 6 |
| **F17** | The revision changelog is a **duplicated, self-divergent splice**: r4, r5, r6, r7, r8, r9, r11, r12, r13, r16, r17, r18, r21 each appear twice with different wording (e.g. r6 "specified the fold" vs "…and removed `revoked_turn`"). A reader cannot tell which description is authoritative — and this is the only place r3 is attested |
| **F18** | The document is **unregistered** in `docs/DOCUMENT-REGISTRY.md`, which states it covers every file under `docs/`. No registry row, therefore no registry hash |
| **F19** | §12's coverage paragraph claims "**I1–I14** are each covered" while the document declares **I1–I18**; no coverage claim exists for I15–I18 |
| **F20** | **Channel × context mapping undefined.** §9 defines C1–C5 with no context parameter; §4 assembles "channels C1–C5 for the named context" without saying which contexts may emit which |
| **F21** | The `[E]/[L]/[X]/[P]/[D]` tag set has **no tag for "accepted candidate, not approved"** — the exact status of §10.1 rev2, on which the whole S24–S26 block rests |
| **F22** | **Asset ownership is never defined.** §11.1 "retains" assets with no owner field, while D1's ownership clause is the document's most emphatic decision |

### 3.5 Low

| ID | Finding |
|---|---|
| **F23** | §13.1 says the record is "~51 KB across thirteen revisions"; it is 91,311 bytes at r32 |
| **F24** | Rule numbering out of order in R1d (`K-R23 → K-R26 → K-R24`); `K-R90` mid-section |
| **F25** | `acquired_turn` is defined as `(turn, seq)` but named and asserted as a turn (`acquired_turn=40`, S37) |
| **F26** | Revisions cited by 8-hex prefixes only, while §13.1 promises "immutable references to prior revisions and their SHA-256" |
| **F27** | `I3` lists "audience" as structured and out-of-band; §5 makes audience a deliberate non-field, derived, never stored |

### 3.6 Positive findings — recorded so the defects are not mistaken for the whole

| Finding | Basis |
|---|---|
| **Rule registry mechanically complete:** 90 rule definitions, **no duplicate definitions**, **no missing** `K-R1`…`K-R90` | Calculated |
| **Probe registry complete:** 94 probes plus `S14n` (95 rows), **no duplicate rows**, **no missing** `S1`…`S94` | Calculated |
| r32's `I14` → `K-R90` relabelling is correctly applied: `I14` appears once as an invariant, `K-R90` once in R1d | Calculated |
| The existence / knowledge / acquisition separation (R1, R1b, R1c, R1d) is genuinely strong; §8.1's fresh-character case demonstrates D2 correctly | Source-derived |
| Fail-closed retrieval default (K-R5) and the summary-vs-account split (I14/K-R90) are the strongest parts of the design | Source-derived |
| **T10 evidence and version boundaries intact:** spec `43d81e07…93a`, receipt `caed2a2f…ca2`, case set `43d744e9…`, all matching `MEMORY-09-G0-RECONCILIATION` §1.2 | Calculated |

### 3.7 What was verified about S8, S16–S29 and T10

- **S8** is preserved verbatim; the data model can represent it; its authority basis is `[D]` (F12). It was **not** weakened and **not** promoted.
- **S16–S29:** all fourteen present. S20, S29 are contradicted as written (F2, F11). S24–S25 correctly blocked on Q1. **S26's block is over-broad** — unequal authority is non-material under §10 without any predicate evaluation, so the cited dependency does not cover it. S28 depends on D8, which is `[P]` proposed. S19/S27 inherit F7.
- **S24–S26 runnable status:** the previously recorded "not runnable" limitation **still applies to S24 and S25**; it does **not** apply to S26 on its cited ground. S26 was **not** silently promoted.
- **T10:** not a versioning probe (it is `scope_gap_known`); the receipt records **no candidate** for the queried fact, so it demonstrates a pool-construction gap, not a demonstrated filter rejection. Evidence hashes intact (F13).

---

## 4. Quality classification of this review

Supplied by the owner, recorded as received.

| Dimension | Verdict |
|---|---|
| Revision verification | **Strong** |
| Evidence discipline | **Strong** |
| Contradiction detection | **Strong** |
| Scope discipline | **Good** |
| Target compliance | **Failed — r3 was unavailable** |
| Runtime validation | **Not performed** |

**Overall (owner's wording): promising.** Recorded here as the owner's assessment of a first
substantial assignment on this model. The reviewer concurs that the one unambiguous success is
procedural: **the mismatch was surfaced before any conclusion was drawn, and r32 was never
substituted silently.**

---

## 5. Next steps — options only, no decision taken

This record resolves nothing. It records what exists and what must be decided by the owner.

1. **Attempt recovery of the r3 artifact from an external backup or another preserved source.**
   Nothing in this workspace can produce it: the file is untracked, has no history, and carries no
   sidecar. If r3 is recovered, the review must be re-run against it — this record does **not**
   transfer to it.
2. **CLOSED 2026-11-17 — owner decision recorded.** The owner authorized changing the review target
   from ARCH-03/r3 to r32 (SHA-256 `e88ef2d1…`, 91,311 bytes); see §1.4. The state is now "review
   complete against r32". It is **not**, and may never be represented as, a review of r3. The
   findings listed in item 3 remain open and are untouched by this change.
3. **Do not treat F2, F3, F5, F6, F7 or F15 as settled by this record.** They are document-level
   findings requiring an owner decision or a corrected revision. None may be closed by a later
   revision appearing to address it, and none may be closed by observing current implementation
   behaviour (the standing rule in `MEMORY-09-OPEN-DECISIONS-2026-11-15-REV-2` §2).

**Not authorized, not performed:** any edit to ARCH-03/r32 or any other project document; creation
of `E-SCOPE-CONTRACT v1.2` or any later version; changes to tests, fixtures, schemas, code or
implementation; model calls, Kobold conversations, corpus generation; G4 or any frozen experiment;
changes to S-cases, probes, preregistration or frozen evidence; registry reconciliation.

---

## 6. Provenance and boundaries of this record

| Field | Fact |
|---|---|
| Record type | Independent review record — documentation and classification only |
| **Target change** | **2026-11-17** — owner authorization: review retargeted from ARCH-03/r3 to r32 (§1.4). Recorded here; no finding validated, superseded or closed |
| Authority created | **None.** This record approves no architecture, freezes no suite, and authorizes no work |
| Files modified | **None during the review.** This record was amended once, on 2026-11-17, under owner authorization, to record the target change (§1.4, §5, provenance row) |
| Files created | **Exactly one** — this record |
| Commands used | `sha256sum`, `grep`, `wc`, `ls`, `git log`, `git status`, `date` (all read-only) |
| Experiments / probes run | **None** |
| Model calls made | **None** |
| Runtime inspected | **None** |
| Revision of the reviewed document | **r32** — *not* the requested r3 |

### 6.1 Declared consequence: registry incompleteness

This record is a new file under `docs/`, which `docs/DOCUMENT-REGISTRY.md` states covers exhaustively.
**No registry row was added**, because this work order authorizes no registry edit — the same
posture taken by `MEMORY-09-REV-2-SUPERSESSION-NOTICE-2026-11-16` §3, item 9. The registry's
completeness claim is therefore now understated by at least two files (this record and
`MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md` itself, finding F18). This is declared, not repaired.

### 6.2 Evidence discipline applied

Every finding above is labelled **Observed**, **Source-derived** or **Calculated**. No finding is
labelled **Inferred** without saying so, none is **Unestablished** yet stated as fact, and nothing
is labelled **Superseded** on a later revision's mere appearance. Where two sources disagree
(Qualification A; F15; F13), the disagreement is preserved and the competing authorities named,
rather than one being silently selected.

### 6.3 Stop condition

The review is complete and returns no fixes. Nothing was repaired, no contract was drafted, no
implementation was begun, and no experiment was run — as instructed.
