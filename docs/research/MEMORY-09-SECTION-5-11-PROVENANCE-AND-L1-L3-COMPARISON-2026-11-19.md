# MEMORY-09 — §5 / §11 Provenance and L1–L3 Comparison

## One narrow comparison before any OD-b selection

**Status:** comparison record. **No selection is made. No schema is proposed.** No field list appears below.
Read-only: no canonical edit, schema change, implementation, test, probe, model call, registry change or ledger
entry.

| Field | Value |
|---|---|
| Date | **2026-11-19** (local, CST) |
| Answers | owner's seven-point request: provenance/chronology · vocabulary coexistence · preserve/lose · rumour · forgetting + statement-less · separations · no selection |
| Decides | **nothing** |
| OD-b | still **open**, still the owner's |
| C′ | still **separate**, untouched here |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 1. Provenance and revision chronology

**CLV-01 (Observed — no commit-level chronology exists).** `git log --follow` on
`MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md` returns **nothing**, and `git status` reports it **untracked**
(`??`). The file has **no git history**. Any chronology below comes from the document's own header, not from
commits.

**CLV-02 (Observed — the revision log is silent on both sections).** r33's header enumerates r1→r33 with a
one-line content summary each. Reading the whole block (`:1–34`), **not one entry names §5's status row or
§11's diagram.** The logged changes are: r1 ownership (superseded) · r2 D1–D5 · r3 edge model · r4
acquisition-vs-edge authority, contexts, classes · r5 projection purity · r6 the fold (`revoked_turn` removed)
· r7 summary-vs-account, continuation, recall · r8 materiality/exposure · r9 precondition tri-state · r10 I15 ·
r11 evaluation provenance · r12 replayability · r13 audit bounds/retention order · r14 deletion contract ·
r15 transfer provenance · r16 transfer hash · r17 conflict scoping · r18 conflict/temporal division · r19
three-valued comparison · r20 interval test · r21 comparison lifecycle · r22 commit-time check · r30/r31
detector identity · r32 labelling · r33 §8.4 classification.

**CLV-03 (Inferred — how much that silence proves).** The log **does** record schema-level edits (r6 removed
`revoked_turn`), so its silence about §11 is not vacuous. But the log is visibly **repetitive** — several
sentences recur verbatim within `:20–34` — so it has been edited or merged and is not a clean change ledger.
Net: silence is **moderate** evidence that §11 was not introduced in a logged revision, and **not proof**.

**CLV-04 (Unestablished).** **Which revision introduced §11's diagram.** It cannot be attributed to r33, to
r32, or to any logged revision. "§11 is the newer, more explicit statement" is therefore **not established** —
the explicit diagram may be as old as the document.

**CLV-05 (Source-derived — r33 did not touch either).** r33's recorded change is the header bump plus the
§8.4 note; §8.4's rule text is preserved unedited (Q-B record: `:49`). Neither §5's row nor §11 is in r33's
change set.

**CLV-06 (Observed — §11 is not marked as proposed).** §11 carries **no** `[P]` / "proposed" marker; only
**§11.1** (chat deletion, `:619`) is marked `[P] — proposed, not adopted`. The diagram and the promotion prose
are presented as design text of the document's own standing, not as a pending proposal.

**CLV-07 (Observed — §5's row has two in-corpus witnesses).** The status row appears in r33 (`:322`) **and**
verbatim in `MEMORY-09-CONTRACT-NORMATIVE-2026-11-15.md` `:46`
(`provisional / active / superseded / retracted / corrected`).

**CLV-08 (Inferred — what that second witness is worth).** The extract is **derivative of the same source**, so
it evidences **wording continuity** across the extraction, not independent authorship or an earlier date. It
is a second *path*, not a second *origin*.

**CLV-09 (Observed — §11's promotion prose has one witness).** "Promotion per kind … rumour — stays
provisional" does **not** appear in the normative extract, which carries rules and invariants rather than §11's
prose and diagram. §11's promotion rules therefore rest on the r33 text alone within this corpus.

**CLV-10 (Source-derived — the other witness for §11's substance is K-R12).** §11's *stream* claim is
independently stated as a rule (K-R12, `:111`), and `MEMORY-09-SESSION-CORRECTIONS-2026-11-17.md` cites K-R12
to r32 `:109`. So the **stream principle** is attested in r32; the **five-state diagram** is not attested
anywhere but r33.

**CLV-11 (Observed — no intra-document precedence rule exists).** Searching r33 for precedence language returns
only `via` evidence ordering (`:155`, `:171`, `:181`, `:864`). There is **no rule stating which of §5, §11 or
K-R12 governs when they differ**, and no section numbering convention that would supply one implicitly.

**CLV-12 (Inferred — document order proves nothing).** §5 precedes §11 on the page. That is a **layout** fact.
It does not establish that §11 elaborates §5, that §5 constrains §11, or that either supersedes the other.
Equally, §11's greater explicitness does not confer precedence: an explicit restatement of an older idea is
not a later decision.

**CLV-13 (Unestablished — the normative relationship).** Whether §11 **elaborates**, **qualifies**, or
**conflicts with** §5 is **not determinable** from the corpus as it stands, because (i) §11's introduction
revision is unknown, (ii) no precedence rule exists, and (iii) on the status vocabulary itself they do not
differ (§2).

**CLV-14 (Calculated — the chronology verdict).** For the **status vocabulary**, §5 and §11 are **identical in
content** and differ only in structure, so no precedence question arises. For **`active` versus
`established`**, the three-way mismatch is **not** a §5-vs-§11 conflict at all — §5 and §11 agree, and K-R12 is
the outlier (§2). Treating the mismatch as a §5/§11 chronology problem would be a category error.

---

## 2. Can the three vocabularies coexist?

**CLV-15 (Observed — the three, side by side).**

| | §5 `:322` | §11 `:600–612` | K-R12 `:111` |
|---|---|---|---|
| Dimensions covered | status only ("Truth/lifecycle status") | status only | **status + `scope` + `knowledge_class`** |
| Names | provisional, active, superseded, retracted, corrected | same five, as entry + four exits | established, superseded, corrected, retracted |
| Entry state | listed, position unstated | **explicit: provisional** | **absent** |
| `active` | listed | explicit child of provisional | **absent** |
| `established` | absent | absent | **present** |
| `forgotten` | absent | absent | absent |

**CLV-16 (Source-derived — §5 and §11 coexist without contradiction).** Same five names. §11 adds structure
(entry, exits, per-kind conditions) and removes nothing. Reading §11 as a **refinement** of §5's row is
consistent with both texts. **Coexistence: yes.**

**CLV-17 (Source-derived — K-R12 is broader in dimensions, narrower in values).** K-R12 governs status **and**
scope **and** class changes on one stream; §5 and §11 govern status only. So K-R12 is not a competing status
list so much as a **statement about one stream carrying three dimensions**, illustrated with four values.

**CLV-18 (Calculated — `retracted` is the only fully stable term).** It appears in all three with the same
meaning. Everything else is either missing from one list or contested.

**CLV-19 (Unestablished — `established`'s referent).** Nothing maps `established` to `active`, to `provisional`,
or to a distinct state. Two readings remain live: (a) a **vocabulary variant** of `active`; (b) a **planned
state** never written into §5/§11.

**CLV-20 (Inferred — two readings of K-R12's scope, and they differ architecturally).**

- **K-R12a (omission).** K-R12 lists the transitions the stream carries and simply omits the entry state. Then
  a provisional record **is** represented on the lifecycle stream, and §11's entry is a transition like any
  other.
- **K-R12b (post-entry scope).** K-R12 describes only what happens **after** entry. Then `provisional` is the
  state **before any lifecycle event exists**, and a rumour that never exits generates **no lifecycle event at
  all** — which matches D4 ("retrievable … never promoted by retrieval", `:72`) and "stays provisional" without
  requiring an event.

**CLV-21 (Inferred — b explains both absences at once).** Under K-R12b, `provisional` **and** `active` are both
absent because the stream's values are the *exit* vocabulary and `established` is the exit target. Under
K-R12a, `established` is left unexplained. That is a point **in b's favour as an explanation**, not a finding
that b is correct.

**CLV-22 (Unestablished — which reading K-R12 bears).** No text chooses. The choice is exactly what OD-b's
representation decision would settle, which is why it belongs inside OD-b and not before it.

**CLV-23 (Calculated — coexistence verdict).** §5 + §11: **unconditional**. All three: **conditional** — K-R12
coexists under K-R12a or K-R12b, and no document says which. The mismatch is therefore **not** "three lists
disagree"; it is **one underspecified scope plus one unmapped term**.

---

## 3. What each representation preserves and loses

Stated at the level of stance, without fields. Common to all three: entry at formation, exit conditional on
evidence kind, and a terminal state must be expressible (ODB-40/ODB-41).

**CLV-24 (Inferred — L1, state-naming).**
- **Preserves:** the state **as an assertion** — a record's position is stated, not inferred; the terminal state
  is positively visible; a lost or unclassifiable event does **not** silently change a record's apparent state.
- **Loses:** the **ground** — a state name says nothing about why the record is there; scope/class changes are
  outside a status vocabulary; and every §11-adjacent concept not in the five names (forgetting above all) has
  no home without a sixth name.

**CLV-25 (Inferred — L2, transition-only).**
- **Preserves:** the change **and its type**; K-R12's projection principle natively; the terminal case as "no
  onward relation", which needs no artificial exit event; and the distinction between *record* transitions and
  *knower* transitions if relations are typed by target.
- **Loses:** the **ground**. A relation says *that* and *what kind of* change occurred, not *why*. It also
  cannot distinguish "never promoted" from "promotion not recorded" without an explicit closure marker — an
  I15/I18 exposure (absence read as information).

**CLV-26 (Inferred — L3, act-only).**
- **Preserves:** the **ground**, natively — the act carries its own kind, which is precisely what L2 lacks; and
  it keeps cause recorded as cause rather than as status.
- **Loses:** the state as an assertion; the derived state is only as good as the derivation rules, and a missing
  act is indistinguishable from no act (the same I15/I18 exposure, more acute because the act *is* the record).

**CLV-27 (Inferred — the owner's preserved question, answered).** **L2 does not by itself retain why a
transition occurred.** The relation carries the type of change, not its justification. Retaining "why" under L2
requires the transition to **reference** a ground, and that reference is a design choice — one that brings the
cause into the transition record. The inversion the owner flagged is therefore **live for L2**, not only L3:
a transition that *carries* its cause can drift into letting the cause authorise the status.

**CLV-28 (Inferred — the separation that avoids the inversion, stated as a constraint, not a design).** Cause
must be **referenced by** the transition, never **read as** the status. A representation that cannot state a
transition without also asserting its ground's authority has merged two of the six dimensions in CLV-30.

**CLV-29 (Calculated — what L2's fold commits to, before choosing it).** A total fold commits to all of:
1. a **closed, typed relation set** — an unlisted relation is undefined behaviour, not a default;
2. **explicit handling of unknown or unclassifiable relations** (fail closed, per I15/I18), not silent
   pass-through;
3. **determinism** — same stream, same order, same result, with order supplied by **A2**'s per-stream ordering
   and nothing else;
4. **no information recovery** — the fold can only project what the relation carries, so the relation vocabulary
   is effectively the specification;
5. **separation of closure from absence** — a record with no onward relation and a record whose successor event
   is missing must be distinguishable, or the fold reports confidence it does not have.

These are commitments of the *stance*, not fields. Item 5 is the one most easily overlooked and the one that
matters for the rumour case (§4).

**CLV-30 (Inferred — the six dimensions the representation must not merge).**

| Dimension | Question | Governing position |
|---|---|---|
| Transition | what changed | §11's progression; K-R12's projection |
| Cause | why | D-3b: actual cause retained, never fabricated |
| Evidence | on what basis | `evidence_span` non-nullable in four schemas, defined in none |
| Occurrence | when, in which order | A2 per-stream; cross-stream **OD-g** |
| Identity | which record | id defects independent of OD-b |
| Authority | who may cause a change | **ND-1**, open and widened |

**CLV-31 (Calculated — no dominance).** L1 uniquely asserts state positively; L2 uniquely matches K-R12's
projection principle; L3 uniquely preserves ground. Each leads on one axis and trails on another. There is no
option that dominates on the evidence, which is why the choice is the owner's and not derivable.

**CLV-32 (Inferred — verification is untouched by all three).** None of L1–L3 gives the verification dimension
a home; C′ survives regardless of selection (ODB-48/ODB-59). Choosing an option is **one** decision, not two.

---

## 4. The permanent-provisional rumour

§11: "rumour — stays provisional, repetition is not corroboration"; and `:612` rejects promotion by retrieval.
D4 (`:72`): provisional records are retrievable, attributed, never promoted by retrieval. This is the only
§11 case where **nothing happens after entry**, which makes it the sharpest discriminator.

**CLV-33 (Inferred — L1).** Expressible: the record's state is `provisional`, and no exit event is required. The
state is **positively asserted**, so later readers cannot mistake "still provisional" for "transition lost". Cost:
`provisional` must be a real, nameable state in the stream's vocabulary despite being absent from K-R12's list.

**CLV-34 (Inferred — L2).** Expressible and cheapest: one `entered` relation, nothing after. **But** the fold
must distinguish "entered, no onward relation" (rumour, complete) from "entered, successor unrecorded" (data
loss). Without that distinction the representation reports a terminal state it cannot actually vouch for
(CLV-29 item 5).

**CLV-35 (Inferred — L3).** Expressible only as "formation act, no exit act". Derivation must then read
**absence of an exit act** as the state — precisely the move I15 forbids ("absence of information is never
treated as information about absence"), and I18 requires every subsystem to name the evidence that would reveal
its own silent omission. L3 is **most exposed** here: a rumour and a lost promotion are the same event stream.

**CLV-36 (Calculated — rumour verdict).** All three express it. L2 is cheapest; L1 is safest against silent
loss; L3 is most exposed to absence-reasoning. The case **discriminates** rather than decides.

---

## 5. DN-4 forgetting, and D-3b statement-less events

**CLV-37 (Source-derived — forgetting is not a record status).** K-R13 (`:110`): knowledge that a claim is
retracted is a **read-time join** of edge × memory projection, never stored on the edge. K-R19/K-R11 govern the
edge fold over acquisition-family events. Forgetting is the knower-side absence of an edge; **retraction is the
record-side status**. The post-Q-B review states the same split (D-b: r32 distinguishes forgetting from
retraction; v9 collapses them). **No enumeration in §2 contains `forgotten`.**

**CLV-38 (Inferred — this is where L1's vocabulary cost becomes concrete).** A state-naming lifecycle stream
must record forgetting as **one of**: (a) a sixth record status — a **normative change** to §5 and §11, not a
representation choice; (b) not on the status stream at all — leaving forgetting unrecorded by the structure
OD-b is deciding; or (c) a second, knower-side stream. Option (c) is L2/L3's shape, adopted by implication.

So L1's "forces vocabulary reconciliation immediately" is **not mainly about `active`/`established`** — that
mismatch is K-R12's, not L1's — it is about **DN-4**, and it lands first.

**CLV-39 (Inferred — L2).** Forgetting is a relation whose **target is the knower, not the record**. If the
relation set is typed by target, L2 expresses it without touching the record's status and without a sixth name.
Cleanest of the three, on this test only.

**CLV-40 (Inferred — L3).** Forgetting is an act, which is natural — but a record-lifecycle act stream has no
place for a knower-side act unless the two streams' scopes are distinguished. The act is expressible; the
**ownership of the stream** is the open part.

**CLV-41 (Source-derived — the statement-less collision is orthogonal).** D-3b permits events that assert
nothing (revocation, forgetting, retraction-of-record). §6.2 marks nullability explicitly three times
(`from_character_id`, `source_memory_id`, `asset_ref`) and does **not** mark `turn` or `evidence_span`, so by
that notation they are non-nullable. **Every** lifecycle transition — entry, supersession, correction,
retraction, forgetting — is a candidate statement-less event.

**CLV-42 (Inferred — therefore it must not pick an option).** The §6.2/D-3b contradiction hits L1, L2 and L3
**equally**, because all three are event streams. It is a **schema-level conflict** and a **precondition** for
any of them, not a discriminator. Solving it inside OD-b would settle evidence and occurrence by accident.

**CLV-43 (Inferred — but it does constrain entry).** If entry at formation is a lifecycle event, it is the
**one** transition that always has an evidence turn, so entry can satisfy §6.2 today. Exits and forgetting
cannot. That asymmetry is worth recording: the entry event is the safest to represent, which is exactly the
reverse of K-R12's silence about it.

---

## 6. Separations held

**CLV-44 (Observed — held throughout, per owner instruction).**

| Dimension | Status here |
|---|---|
| Verification (C′) | **not** touched; all three fail to home it (CLV-32) |
| Cause | distinguished from status (CLV-27/28); D-3b governs |
| Evidence | flagged as a §6.2 conflict, **not resolved** (CLV-41/42) |
| Occurrence | A2 per-stream assumed; cross-stream **OD-g** left open |
| Identity | not examined; id defects independent of OD-b |
| Authority | **ND-1** open; no lifecycle actor named or implied |

**CLV-45 (Inferred).** None of the six was needed to answer the seven questions, which is itself evidence that
they are separable — and that OD-b can be decided without them if the decision is confined to *what the
structure records about the progression*.

---

## 7. What this record does not do

**CLV-46 (Observed).** It makes **no selection** among L1/L2/L3, proposes **no schema**, and names **no fields**.

**CLV-47 (Observed).** It does **not** declare §5 or §11 controlling. It establishes that the question of
precedence does not arise for the status vocabulary (CLV-14), and that no rule exists to decide it if it did
(CLV-11/CLV-13).

**CLV-48 (Observed).** It does **not** reopen OD-a, A2, D-3b, DN-1, DN-4 or Q-B, and does **not** settle
`active`/`established`, `forgotten`, C′ or OD-g.

**CLV-49 (Observed).** No canonical file was edited; the 118 requirements are unchanged; no ledger entry was
made; the Howl state machine was not touched.

---

## 8. Limits

**CLV-50 (Observed).** Evidence is r33's header and text, the normative extract, K-R12's citation to r32, and
the repository's git state. **Not inspected:** r32's own bytes (no copy survives in the workspace — only r33 and
the extract); any pre-r32 revision; the frozen preregistration's treatment of §11; `docs/DOCUMENT-REGISTRY.md`
contains **no row** for either contract (grep returns nothing, consistent with finding F18). The revision log's
repetitions were noted, not reconstructed.

**End of comparison record.**
