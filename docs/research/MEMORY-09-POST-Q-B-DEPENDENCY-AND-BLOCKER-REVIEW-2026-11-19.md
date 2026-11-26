# MEMORY-09 — Post-Q-B Dependency and Blocker Review

**Status:** research report. **READ-ONLY.** No implementation, no test, no probe, no model call, no
corpus, no canonical edit, no registry change, no Git operation, no architecture selection. Q-B is
treated as settled and is not reopened.

| Field | Value |
|---|---|
| Created | **2026-11-19** (local, CST) |
| Scope | What remains unresolved now that O-14 settles Q-B |
| Baseline | Eight records read and hash-verified this session (§1) |
| Decisions made here | **None.** No relation adopted, no subject chosen, no requirement amended |
| Tests executed | **0** |
| Files changed | **One new file only.** Untracked, unregistered, unfrozen |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**. Every material claim carries one.

---

## 1. Baseline

Hashes recomputed from disk this session; all eight are the current bytes.

| Record | Lines | SHA-256 (first 16) | Role |
|---|---|---|---|
| `MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md` | 1040 | `dfb28b6003ca356f` | **r33** — §8.4 classified |
| `MEMORY-09-OWNER-DECISIONS-2026-11-17.md` | 688 | `e5008020ff68bd21` | **T1 ledger**, O-14 at `:205` |
| `MEMORY-09-Q-B-OWNER-DECISION-2026-11-19.md` | 189 | `9db51c00c6a3748e` | Decision record, scope frozen |
| `…-1A-1B-RETRACTION-AWARENESS-AND-RELATIONSHIP-DECISION-BRIEF-…` | 667 | `70245c1627e9e749` | Q-B/Q-A/Q-C framing |
| `…-1A-1C-Q-B-NO-PATH-PRESERVATION-VERIFICATION-…` | 399 | `0cee0d7226524caf` | Seven preservation checks |
| `…-1A-1D-Q-B-NO-PATH-CLARIFICATION-…` | 240 | `a2d3425f258d29d0` | Six corrections |
| `MEMORY-09-PHASE-1A-CONCEPTUAL-OBJECT-INVESTIGATION.md` | 1071 | `82173a6f78ac84c4` | CO/PH1A items |
| `MEMORY-09-CURRENT-STATE-AND-COMPLETION-AUDIT-2026-11-18.md` | 755 | `6c3b0db0df65a005` | P1–P12 prerequisites |

**DB-00 (Observed).** r33's §8.4 note is at `:427`; O-14 is at `:205` with its broad-reading paragraph
at `:228`; the decision record's scope freeze is at `:181`. None was edited by this review.

**DB-00b (Observed + Calculated).** The audit's own prerequisite count is inconsistent with its own
table: §1 says "**nine** independent prerequisites are unmet" (`:53`), while §12 tabulates **twelve**
(P1–P12), of which **ten** are marked as blocking for the next step. Not material to any conclusion
below; recorded because this review cites P1–P12 individually.

---

## 2. Q-B closure, verified rather than assumed

**DB-01 (Source-derived).** O-14 states the decision in three clauses: no automatic awareness from
status; characters learn by legitimate acquisition; MEMORY-09 handles retraction independently of
awareness. Its table names the routes as K-R11's event kinds, and its `:228` paragraph fixes the
reading — an inference or verification acquires without a disclosure occurring.

**DB-02 (Observed).** r33 `:427` classifies §8.4's display as illustrative and non-binding without
altering the rule text; A2b, E1c and D2 are untouched in the conformance map
(`b4ba33ccbfe8f8c1`, unchanged). Nothing in this review proposes otherwise.

**What this review did not do**, by instruction and by choice: reopen Q-B; edit r33 §8.4 or O-14;
reconsider A2b/E1c/D2; reintroduce status-derived awareness; treat inference as disclosure; convert P-5
or P-7 into authorisation. P-5 and P-7 are cited below only as gaps.

---

## 3. Dependent questions, reassessed

### 3.1 Q-A — shared content identity

**DB-03 (Calculated).** Zero of the 118 conformance requirements state, imply or test that two
characters' accounts be recognisable as concerning the same claim. A search of the full map for
`same content`, `same underlying`, `recognise`, `recognize`, `sameness` returns **0 matches**. The rows
that touch plurality of knowers — A2a, E1a–E1c, H4, H10, J9, D3 — all concern what one knower holds or
what two records do independently; none requires a relation between them.

**DB-04 (Source-derived).** Q-A's necessity was always conditional, and the condition has failed.
1A-1B §8 result 1 and §12 item 2 state the dependency: under a "no" to Q-B, RA-3's content identity
"buys nothing for retraction". Q-B is now "no" (O-14). Therefore **no independently established
requirement requires Q-A**, and none would be satisfied differently by answering it either way.

**DB-05 (Inferred).** Q-A is not *closed* — it was never answered, and cross-knower conflict detection
or a future "who else holds this" feature would reopen it — but it is **deferrable at zero cost to the
ratified set**. Its remaining justifications are all outside the 118. This is the single largest change
Q-B made to the decision landscape: the one question that looked structurally urgent is now optional.

### 3.2 Q-C — attribution of a notice

**DB-06 (Source-derived).** Q-C is **moot as posed**. RQ-02 asks whether a notice names the retracting
actor; with O-14 there is no notice, so the question has no object. 1A-1C §6.1's three forms all
presuppose something is shown to the character.

**DB-07 (Source-derived).** Three attribution questions are **not** resolved by that mootness, and none
of them is Q-C:

| Question | Status | Why it is separate |
|---|---|---|
| Who may retract a character's **account**? (RQ-03) | **Open** | asks authority over knowledge, not how a notice is phrased; `retract` still takes no actor (`engine.py:542–559`, RB-07) |
| Does a retraction **act** have a recorded actor at all? | **Open** | no T1 requirement demands it; E1a/E1c/H10 are silent on attribution |
| Who supplied a contribution? (rules 23, 26; O-12) | **Settled** | declared, not judged; unaffected by Q-B |

**DB-08 (Inferred).** With the notice gone, RQ-03 becomes the sharpest *behavioural* question in the
series: retracting an account is the one operation that removes knowledge outright, it has no owner
check, and under Q-B = no it produces no compensating signal to anyone. This is not a restatement of
Q-C; it is the question Q-C was standing in front of.

### 3.3 Q-D — conversation-state applicability

**DB-09 (Source-derived).** Exact current scope, from 1A-1A §13.2 and O-14's carrier paragraph:

- **RQS-01** — does MEMORY-09 track conversation states at all, and if so does applicability live on the
  record, in a projection over a transition log the chat supplies, or in partitions? Whichever is
  chosen, applicability must **not** become a fourth value of `status`.
- **Ownership is assigned, the answer is not**: Howl owns the state machine, state identity and current
  state (HSL-1–HSL-3); MEMORY-09 owns representation, retrieval and dependency (HSL-9–HSL-11); the
  crossing contract is deferred to Howl's foundational architecture work.
- **RQS-02/RQS-03/RQS-06** moved to the Howl requirements record (HSL-O1–HSL-O3); **RQS-04/RQS-05**
  remain MEMORY-09 semantics.
- **RPS-10** — a state-blind fixture cannot host a state-aware experiment.

**DB-10 (Superseded).** 1A-1A §11 item 3 and §13.2 order Q-B *after* Q-D, because "a retraction occurs
in a state and the notice join is state-relative" (RPS-09). That ordering rationale is discharged by
Q-B's answer: there is no notice join. The ordering constraint is superseded; Q-D itself is not.

**DB-11 (Inferred).** Q-D's residual scope is therefore **applicability, not awareness**: does a memory
become inapplicable when a conversation state is reverted or regenerated, and where is that
represented? That question is independent of Q-B, still open, still unowned on the MEMORY-09 side of
the boundary, and still blocks a state-aware experiment. It is narrower than 1A-1A described.

### 3.4 RA-1, RA-2, RA-3

| Relation | Post-Q-B assessment | Label |
|---|---|---|
| **RA-1** (r32 shared memory + edges) | Still fails O-1 in every reading: the edge names the memory, so B gains access to A's original | Source-derived (1A-1B §4.3, §8) |
| **RA-2** (v9 copied accounts) | Satisfies every settled requirement. No feature of O-14 needs anything else | Source-derived |
| **RA-3** (copies + content identity) | Its only advantage was computing a notice across accounts. With no notice, it adds new architecture for no ratified behaviour | Source-derived + Inferred |

**DB-12 (Inferred).** **No new relationship model is necessary.** RA-2 is sufficient for the settled
set; RA-1 is excluded by a tier-1 decision; RA-3 is now unmotivated by any requirement. That is a
statement about necessity, not a selection, and this review adopts none of them.

---

## 4. The object-model divergence, re-examined

r32's event-sourced A′ model and v9's record-class model stay distinct throughout.

### 4.1 Differences that remain substantive

| # | Difference | Why it survives Q-B | Label |
|---|---|---|---|
| **D-a** | r32 has acquisition **events + edges**; v9 has **accounts** only | F1/F1b can be satisfied by accounts, but the *history* of how knowledge changed (learn → forget → relearn) has no v9 carrier | Observed (CO-02, CO-10) |
| **D-b** | r32 distinguishes **forgetting** (no edge, K-R19) from **retraction** (status); v9 collapses both into a record status | RP-01; the distinction is T2 and untested by any of the 118 | Observed (CO-07) + Source-derived |
| **D-c** | r32's `MemoryRecord` has **no character owner**; v9's private record is exclusively owned | O-1 is T1 and v9-only; r32 cannot express exclusive ownership without edges | Source-derived (CO-02) |
| **D-d** | r32 orders by `(turn, seq)`; neither implementation has `seq` | Same-turn determinism is unrepresentable in both (PH1A-05) | Source-derived |

### 4.2 Differences Q-B made irrelevant

| Difference | Why it no longer matters |
|---|---|
| The **retraction notice join** across accounts | O-14: no notice is computed |
| **Sameness** of content across knowers | Q-A is not required (DB-04) |
| **Attribution** of a retraction to an actor in a notice | Q-C is moot (DB-06) |
| §8.4's display as a behavioural requirement | classified in r33 `:427` |

### 4.3 What neither implementation can represent

| Requirement | Why neither can carry it | Label |
|---|---|---|
| Same-turn ordering (K-R15; S38) | no `seq` in either model | Source-derived (PH1A-05) |
| Event-kind refusal in declarations (B3) | no event/definition distinction on a record | Source-derived (PH1A-09) |
| Human contribution as a claim (L1, rule 22) | engine refuses `contributor_kind == "human"`; harness has no contributor kinds | Observed (CO-08) |
| Materiality E1–E3 (if in scope) | no conflict group in the harness; no materiality in the engine | Source-derived (PH1A-08) |
| Transfer guarantees K-R46–K-R48 | closed as application-level by D-4; MEMORY-09 absent from the path | Source-derived (PH1A-14, D-4) |

### 4.4 Which decisions block selecting an experimental subject

**DB-13 (Observed + Inferred).** PH1A-12 asks "which subject does the experiment test:
`eval/mem1/model.py` or `eval/mem1/m09/`?" It remains **blocked**, and its predecessor set has changed:

- The edge from **Q-A** is **removed** — Q-A is no longer forced (DB-04), and 1A-1B §12 item 2's
  dependency ran through the notice.
- The edge from **Q-D** **remains**, narrowed to applicability (DB-11): a state-blind fixture cannot host
  a state-aware experiment (RPS-10).
- A predecessor **not previously recorded as such** is now visible: **PH1A-03/PH1A-04**. "Which subject"
  means "which implementation expresses the required mechanism set", and that set is exactly what those
  two items contest — r32's acquisition and lifecycle streams (T2) against v9's records-with-status
  (T3). Choosing a subject before deciding them answers the stream question by convenience.

**DB-14 (Inferred).** That is the substantive correction this review offers to the audit's
recommendation. The audit calls PH1A-12 "the next smallest action" (§15.2 item 7) on the ground that
corpus shape depends on it. True — but PH1A-12 is not answerable *as posed* until the mechanism set is
fixed, or it will be answered by whichever implementation is easier to point at. The dependency runs
PH1A-03/04 → PH1A-12 → corpus, not PH1A-12 → corpus.

---

## 5. Report

### 5.1 Settled requirements

| Item | Authority | Status |
|---|---|---|
| Retraction changes truth status only; knowers keep what they acquired | E1a, E1b, A2a, H10, O-3 | settled, unchanged |
| No status-derived character awareness; acquisition routes only | **O-14**, r33 `:427` | **settled by Q-B** |
| §8.4 display is illustrative, not a requirement | r33 `:427` | settled |
| Belief is model-declared and never system-revised | O-5, rules 15/53–55 | settled |
| Exclusive private ownership; recipients hold their own record | O-1, rule 12 | settled |
| Disclosure creates recipient accounts; narration cannot manufacture one | O-3, rules 46–49 | settled |
| 57 rules ratified; 118 requirements defined; 0 executed | O-6b | settled |

### 5.2 Remaining architectural conflicts

| ID | Conflict | Tier | Status |
|---|---|---|---|
| **AC-1** | r32 acquisition event-sourcing (K-R1, K-R11, K-R13, K-R15) vs v9's records without events | T2 vs T3 | **open** — PH1A-03 |
| **AC-2** | r32's lifecycle stream (K-R12) vs v9's mutable `status` field | T2 vs T3 | **open** — PH1A-04 |
| **AC-3** | Forgetting vs retraction collapsed into one status token | T2 vs T3 | **open** — RP-01/CO-07 |
| **AC-4** | D6 materiality vs v9's exclusion of materiality | T1 vs T3 | **open** — PH1A-08 |
| **AC-5** | Event vs definition (B3) has no representation | T3 | **open** — PH1A-09 |
| **AC-6** | REV 2's proposition registry rests on recommendations, not a decision | T3 | **open** — PH1A-13 |
| ~~AC-7~~ | §8.4 display vs A2b/E1c/D2 | T2 vs T1 | **closed by O-14** |
| ~~AC-8~~ | Recipient notice attribution | T2/T3 | **moot by O-14** |

**DB-15 (Source-derived).** AC-1 and AC-2 are the same structural question asked of two objects: does
the mechanism have a stream, or a field? They are the root of the remaining model divergence.

### 5.3 Owner decisions still required

| Priority | Decision | Blocks |
|---|---|---|
| **1** | **PH1A-03/PH1A-04** — do r32's acquisition and lifecycle streams survive as required mechanisms? | subject selection; bridge design; AC-1/AC-2/AC-3 |
| **2** | **PH1A-12** — which subject the experiment tests | all experiment work (P1 corpus onward) |
| **3** | **Q-D / RQS-01** — where applicability lives on MEMORY-09's side | state-aware experiments |
| **4** | **PH1A-01** — four-class taxonomy: retain, decompose, retire | prereg taxonomy vs engine classes |
| **5** | **RQ-03** — who may retract an account | the one operation that removes knowledge |
| 6 | PH1A-05 (`turn, seq`), PH1A-07 (incompatibility), PH1A-08 (materiality scope), PH1A-09 (event vs definition), PH1A-10 (human contribution), PH1A-11 (retrieval trigger), PH1A-13 (REV 2) | implementation and conformance detail |
| — | **Q-A** and **Q-C** | **nothing in the ratified set** (DB-04, DB-06) |

**DB-16 (Calculated).** The investigation tables **15** `PH1A` items (PH1A-01 … PH1A-15). **One is
settled** by O-14 — PH1A-06, "must a knower be shown a retraction": no — and **14 remain**. Separately,
two non-`PH1A` edges are removed (Q-A and Q-C). The apparent open-question count barely moved; what
moved is which ones matter.

### 5.4 Experimental blockers

From the audit's P1–P12, unblocked by nothing in this review:

| Blocker | State | Changes under Q-B? |
|---|---|---|
| **P1** corpus (16 conversations) | absent | no |
| **P2/P3** probe identities, `plan_sha256` | derived from P1 | no |
| **P5/G07** tokenizer binding | absent | no |
| **P6/G12** model revision + sampling binding | absent | no |
| **P7** grant mechanism (`ExecutionNotAuthorised`, `conversation.py:177`) | absent — structural | no |
| **P8** live responder | absent — structural | no |
| **P9** engine ↔ harness bridge | absent — structural, undefined | **worse**: the bridge's shape now depends on PH1A-03/04 |
| **P10/P11** execution and publication authorization | absent | no |

**DB-17 (Observed).** Q-B changed no experimental blocker. It removed an *architectural* reason to
answer Q-A, which shortens the decision chain but not the build chain. The three structural blockers
(P7, P8, P9) remain exactly as the audit described them.

### 5.5 Implementation gaps that are not authorised work

| Gap | Evidence | Status |
|---|---|---|
| **P-5** no `inspection` context | `m09/retrieval.py:29`; `retrieve.py:24` | unbuilt, unauthorised |
| **P-7** disclosure implemented, uncalled | `engine.py:575–700`; no call site | unexercised, unauthorised |
| **V-5** retraction represented as a record (v9) vs a status (m09); record's meaning Unestablished | `model.py:282–299` | divergence, open |
| **RP-02/GB-8** revoked account unreachable to its owner | `m09/retrieval.py:50` | defect-shaped, unauthorised to fix |
| **CO-08** `claim()` refuses human contributors | engine | uncovered by any decision |
| **PH1A-15** no review covers the committed 1,951-line engine | audit §13 | evidence gap |

**DB-18 (Source-derived).** Under O-14 these gaps matter more than they did, because acquisition is
now the only route to awareness and `inspection` is now the only place the system's knowledge is
visible to a human. They remain unauthorised work, not implied tasks.

---

## 6. The single most important next owner decision

**DB-19 (Inferred).** It is **PH1A-03/PH1A-04, framed as one question**:

> Does a conforming MEMORY-09 keep r32's two streams — an append-only acquisition event stream and an
> append-only memory-lifecycle stream, with status and knowledge as projections of them — or does it
> keep v9's model, in which an account is a record and status is a field?

**Why this and not PH1A-12.** Every remaining choice is a consequence of it:

- it decides AC-1, AC-2 and AC-3 in one act;
- it determines what a subject must be able to express, and therefore answers PH1A-12 rather than
  being pre-empted by it (DB-13/DB-14);
- it fixes the shape of the P9 bridge, which no record currently defines;
- it decides whether "forgot" and "retracted" can ever be distinguished — the one conflation both
  implementations currently make.

**Why it is genuinely required, not manufactured.** O-6a case 2: v9 contradicts an applicable r32
mechanism, so "neither rule is implemented until explicitly decided". The divergence is real, it is
recorded (PH1A-03/04), and nothing in Q-B touches it. It has been open since Phase 1A and is the only
unresolved question that gates every other remaining choice.

**What it does not require.** No new architecture, no code, no corpus, no test. One written answer.

---

## 7. Limits and what this review does not do

No architecture selected, no relation adopted, no subject chosen, no requirement added, amended or
relaxed, no probe designed, no threshold invented. No canonical document edited; no registry row; no Git
action. No file executed.

**Static only.** Every behavioural claim is a reading at the hashes in §1. Nothing was run, so nothing
here demonstrates runtime behaviour. P-5 and P-7 are described as gaps, not as evidence of failure.

**Not re-litigated.** Q-B (settled), A2b/E1c/D2 (untouched), r33 §8.4 (classified), O-14 (closed to
edit), S3 (closed), the 57 rules and 118 requirements (ratified).

**One caution offered, not asserted.** If PH1A-12 is put to the owner before PH1A-03/04, the most likely
outcome is a choice of whichever implementation is nearer the harness — which would settle the stream
question by convenience. That is the failure mode this review most wants to prevent.

---

## Appendix — Claim inventory

| Label | Uses | Representative claims |
|---|---|---|
| Observed | 7 | baseline hashes; P-5/P-7 code state; audit count inconsistency |
| Source-derived | 12 | O-14 clauses; requirement silence on sameness; RA-1/RA-2 status; AC-1..AC-6 |
| Calculated | 3 | 0/118 sameness matches; 15 PH1A items (1 settled, 14 open); P1–P12 counts |
| Inferred | 6 | Q-A deferability; RQ-03 as the sharper question; dependency order PH1A-03/04 → PH1A-12 |
| Unestablished | 1 | V-5 record's semantic status (carried, unresolved) |
| Superseded | 1 | 1A-1A's Q-D → Q-B ordering rationale (RPS-09 join) |

**End of report.** The next action is an owner answer to PH1A-03/PH1A-04 — or an explicit decision to
defer it, which would itself answer it.
