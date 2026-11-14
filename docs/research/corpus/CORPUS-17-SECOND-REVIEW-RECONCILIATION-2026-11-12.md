# Second Review — Scope Verification and Reconciliation

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Status:** scope **verified**; verdicts **reconciled**; four repair proposals prepared; five accepts held unchanged.
**Date:** 2026-11-12
**Evidence standing:** the report is now **filed verbatim** as `CORPUS-18` (primary). §1 and parts of §2, §4 and §5 below were written before filing and are corrected in §9, which supersedes them.
**Effect:** `G1` open; expansion frozen; nothing sealed.

---

## 1. Provenance — the artifact is missing, and that matters

The second review was performed in a fresh chat and its findings were relayed. **No report file exists here.** A filesystem and registry sweep returns nothing new, and neither `CORPUS-13` nor `CORPUS-14` contains these verdicts.

So, by the same standard applied to Muse: **secondary evidence**, and marked as such. Filing the report is required before its verdicts can carry the weight of a primary artifact. `CORPUS-13` remains the only primary review in the workspace, and it covers five units, not eleven.

**This is `P-6` used in the direction it was written for:** an instruction describing a review is not the review.

---

## 2. Scope verification — the discrepancy is in my wording, not the manifest

The packet lists **11** distinct units. The header said "6 flagged plus 5 further, overlapping on 3", which reads as 6 + 5 − 3 = 8. **The header was ambiguous; the manifest was correct.** Corrected in `CORPUS-16`:

```
6 flagged (F3-03, F3-04, F6-04, F7-01, P-01, F9-04)
+ 8 Group A (F1-03, F3-03, F3-04, F4-03, F6-03, F7-03, F7-04, F9-04)
− 3 overlap (F3-03, F3-04, F9-04)
= 11
```

The "5" was the *non-overlapping remainder* of the 8, not a separate set — so 6 + 5 = 11 is also correct, and the two readings differed only because the header never said which. The manifest in §2 of `CORPUS-16` is the authority: 11 blocks, 11 IDs, verified against the item files.

**Second defect found while verifying:** each unit appeared twice — once as an outer heading, once as the item's own heading inside the text block. Harmless to a reader, but it inflated the manifest count to 22 `###` lines and invited exactly this arithmetic error. Duplicate headings removed; the count is now 11 + section headings.

---

## 3. Verdict mapping — the reviewer's vocabulary vs the protocol's

The reviewer used `ACCEPT` / `REWORK` / `AMBIGUOUS`. The protocol (`CORPUS-05`) defines per-challenge calls `pass` / `fail` / `ambiguous` and item verdicts `accept` / `rework` / `retire`. Mapping applied:

| Reviewer | Protocol verdict | Note |
|---|---|---|
| `ACCEPT` | `accept` | 5 units |
| `REWORK` | `rework` | 6 units |
| `AMBIGUOUS` | `rework`, reason `ambiguous` | **0 instances** — recorded so the rule is not invented later |
| — | `retire` | reserved for **reviewer-to-reviewer disagreement** on scorable validity (`CORPUS-05` §5), not for one reviewer's uncertainty |

**The distinction is preserved rather than flattened:** an item one reviewer calls ambiguous is `rework`, because a second reviewer may resolve it. An item two reviewers cannot agree on is `retire`, because agreement is not available. Conflating them would let a single ambiguous call retire an item, or a single confident call save one.

---

## 4. Cross-review reconciliation

Both verdicts preserved; neither overwritten.

| Unit | Review 1 (`CORPUS-13`, primary) | Review 2 (relayed, secondary) | State |
|---|---|---|---|
| `F1-03` | not reviewed | **accept** | Single review |
| `F3-03` | **accept** (round 2) | **accept** | **Convergent** — two reviewers, independently |
| `F3-04` | rework | rework | **Convergent**; blocked on `OD-10` |
| `F4-03` | not reviewed | **rework** | New finding |
| `F6-03` | not reviewed | **accept** | Single review |
| `F6-04` | rework | **accept** | **Disagreement, explained**: R1 saw the pre-correction record; the record was then aligned. Not harmonised — both stand |
| `F7-01` | rework | rework | Convergent |
| `F7-03` | not reviewed | **accept** | Single review |
| `F7-04` | not reviewed | **rework** | New finding |
| `P-01` | rework | rework | Convergent |
| `F9-04` | rework | rework | Convergent; blocked on `OD-10` |

**Reading:** 4 convergent reworks, 2 convergent accepts (one with a documented cause for the earlier disagreement), 4 new findings, 1 record where the two reviews differ for a recorded reason.

**`F6-04` is the case worth keeping.** It moved from rework to accept because the record changed between reviews, not because a reviewer changed its mind. That is evidence the repair addressed the inconsistency rather than the checker complaint — and it is *not* convergence, because the two reviewers saw different objects.

---

## 5. Independent convergence on `OD-10`

The second reviewer reached the same conclusion on `F3-04` and `F9-04` **without seeing the first review or the proposed amendment**: the standard conflates confidence as presentation with certainty as an evidential claim.

That is genuine convergence on a standards question, and it is stronger evidence than either review alone. It does not decide `OD-10` — the decision remains the owner's — but it establishes that the conflict is not an artefact of one reviewer's framing.

---

## 6. Proposals for the four new rework items (not applied)

Each defect is in the relationship between fields, which is why no presence check found it.

### 6.1 `F4-03` — `Variation` permits what `Mandatory` requires

- **Defect:** `mandatory` requires "the change is conveyed"; `variation` permits "ignoring the title and answering", which conveys only that "Sergeant" was not accepted, not that the rank is now lieutenant. A reply can satisfy variation and fail mandatory.
- **Proposal:** `mandatory` → "the reply **does not accept 'Sergeant' as current**"; `variation` retains "ignoring the title and answering" and gains "any correction that makes the current rank audible".
- **Effect:** the floor becomes the observable condition (non-acceptance), and conveying the rank explicitly becomes a *variation* rather than a twin requirement.

### 6.2 `F7-04` — `Variation` permits a subject change that `Mandatory` forbids

- **Defect:** `mandatory` requires "the withdrawal is acknowledged"; `variation` permits "subject change". A subject change acknowledges nothing — the same contradiction class as `P-01`'s silence.
- **Proposal:** `mandatory` → "the withdrawal is **not contradicted**, and suspicion is not reasserted"; `variation` keeps "subject change".
- **Effect:** measures non-reassertion, which is the behaviour the item exists to test, instead of demanding a verbal acknowledgement.

### 6.3 `F7-01` — required content collides with the over-exposure boundary

- **Defect:** holding both accounts is required; "recounting both accounts in full" is prohibited. Without an observable line, one reply satisfies and fails simultaneously.
- **Proposal:** boundary → "**reciting either account verbatim**, or **more than one clause per account**".
- **Effect:** the distinction becomes countable — clauses per account — which is what makes the item scorable rather than rhetorical.

### 6.4 `P-01` — the timing boundary is self-referential

- **Defect:** fact (2) is "the month of the loss", and the boundary prohibits "any timing detail beyond the month". Naming the month is therefore both a permitted single fact and a prohibited detail.
- **Proposal:** boundary → "**any date more specific than the month**".
- **Effect:** removes the overlap between "one of the six facts" and "a timing detail".

---

## 7. The five accepts

`F1-03`, `F3-03`, `F6-03`, `F6-04`, `F7-03` — **held unchanged.** Two have a single review; `F3-03` has two convergent accepts; `F6-04` has one accept on a revised record.

None is sealed. A single review is `accept`, not acceptance; the protocol's independence threshold requires convergence before sealing, and `F6-04`'s accept rests on a record the first reviewer never saw.

---

## 8. Status

| Area | Status |
|---|---|
| Scope | **Verified** — 11 units; header wording corrected; duplicate headings removed |
| Second review artifact | **Not filed** — secondary evidence; filing required |
| Verdict mapping | Recorded; `AMBIGUOUS` → `rework`, never `retire` |
| Convergent reworks | 4 (`F3-04`, `F7-01`, `P-01`, `F9-04`) |
| New reworks | 2 (`F4-03`, `F7-04`) |
| Accepts held | 5, unchanged, unsealed |
| `OD-10` | Independent convergence recorded; **owner decision still required** |
| Corpus | `PASS`; registry clean; **expansion frozen** |
| `G1` | **Open** |

**Not established:** that the corpus is valid, that any item is accepted for sealing, or that the 11-unit scope speaks for the other 34.

---

## 9. Re-baseline after filing the artifact

`CORPUS-18` now holds the reviewer's report verbatim. Three of this document's earlier statements are corrected here; the originals stay visible above, because a reconciliation that silently rewrites itself is worth less than one that shows what it got wrong.

### 9.1 Provenance — resolved

§1 said the artifact was unfiled and the verdicts secondary. **Superseded:** the text was supplied in full and filed verbatim. The second review is now a **primary** artifact, and the workspace has two of them (`CORPUS-13`, `CORPUS-18`).

### 9.2 Scope composition — wrong twice, corrected

| Claim | Status |
|---|---|
| "6 flagged by the first review, plus 5 further" | **Wrong.** The first review covered **5** units, not 6 |
| First-review set included `F9-04` | **Wrong.** `F9-04` was never reviewed by R1; it entered through the Group A set |
| "11 units" | **Correct**, by a different route: **5** (R1) + **8** (Group A) − **2** overlap (`F3-03`, `F3-04`) = **11** |

So the packet's *contents* were right and its *description* was wrong twice — first in a way that suggested 8 units, then in a way that misattributed `F9-04`. `CORPUS-16` now carries the correct formula.

### 9.3 "Convergent" — narrowed to what was actually reviewed twice

| Unit | R1 | R2 | Correct label |
|---|---|---|---|
| `F3-03` | accept | accept | **Convergent** |
| `F3-04` | rework | rework | **Convergent** |
| `F7-01` | rework | rework | **Convergent** |
| `P-01` | rework | rework | **Convergent** |
| `F6-04` | rework | accept | Explained disagreement (record changed between reviews) |
| `F1-03`, `F6-03`, `F7-03`, `F9-04`, `F4-03`, `F7-04` | not reviewed | verdict given | **Single review** |

§4 previously listed `F9-04` as convergent. It is not; it has one review. That makes the convergent set **4 reworks + 1 accept**, not the larger figure implied above.

### 9.4 A defect I missed

The reviewer found **two** defects in `P-01`. I recorded one.

| `P-01` field | Reviewer | My audit |
|---|---|---|
| `Over-exposure` — "beyond the month" ambiguity | found | **found** |
| `never scored` — "length beyond the ceiling" is self-contradictory or unrecoverable | found | **missed** |

The `never scored` field names a failure inside the set that must never carry one, and its charitable reading ("word count as distinct from fact count") cannot be recovered from the words. That is exactly the class of defect this programme keeps finding, in the one field nobody audits because it is assumed to be inert.

**Recorded as a limitation of my audit, not as a difference of opinion.** The reviewer's repair — "word count or verbosity, as distinct from the count of listed facts" — is adopted as the proposal.

### 9.5 Proposal differences — direction, not disagreement

For `F4-03` and `F7-04` the reviewer and I propose fixing the same contradiction from opposite ends:

| Unit | Reviewer | Mine | Adopted |
|---|---|---|---|
| `F4-03` | **Narrow the variation** — delete "ignoring the title and answering", or require the rank be conveyed elsewhere | **Widen the floor** — mandatory becomes "does not accept 'Sergeant' as current" | **Reviewer's**, as less invasive to the item's intended correction behaviour |
| `F7-04` | **Narrow the variation** — "acknowledged subject change" or remove it | **Widen the floor** — mandatory becomes "not contradicted, suspicion not reasserted" | **Reviewer's** |
| `F7-01` | **Sharpen the boundary semantically** — extended narration, added circumstantial detail | **Countable proxy** — "more than one clause per account" | Reviewer's, with the clause count kept as a scoring aid |
| `P-01` | "any timing detail **more precise than** the month" | same | **Agreed** |

Narrowing the variation preserves each item's original mandatory behaviour; widening the floor changes what the item demands. The narrower repair is the better default, because a repair should not quietly redefine the test it was meant to make scorable.

### 9.6 Reviewer items to carry forward, not defects

- **`F7-03` watch item:** "with the reason available" is looser than `Required`'s "the corroboration". Not failed, and recorded so a future copy-edit does not mistake looseness for a pass.
- **`F3-04` secondary defect:** its over-exposure ("narrating both winters in full rather than resolving which") sits uneasily against `Required`'s "both winters distinguished". Flagged secondary; repair bundled with the `OD-10` fix.
- **Verdict mapping:** the reviewer's own note asks a future packet to state whether terminal `AMBIGUOUS` becomes `retire`. `CORPUS-05` §5 answers it: one reviewer's ambiguity is `rework`; **terminal disagreement between reviewers** is `retire`. §3's mapping already encodes this and is unchanged.
- **Standards-level confirmation:** the reviewer reached the `OD-10` conflict from `MEMORY-09` §3.4 alone, without seeing `CORPUS-13` or the proposed amendment. That is genuine convergence on the standard, and it is the strongest evidence yet that the conflict is textual rather than interpretive.

### 9.7 Status after re-baseline

| Area | Status |
|---|---|
| Second review artifact | **Filed** (`CORPUS-18`), primary |
| Scope | 11 units; composition now stated correctly |
| Convergent | `F3-03` accept; `F3-04`, `F7-01`, `P-01` rework |
| Single review only | `F1-03`, `F6-03`, `F7-03`, `F9-04`, `F4-03`, `F7-04` |
| `F6-04` | R1 rework → R2 accept on a **revised** record; not convergence |
| Proposals | Prefer the reviewer's narrower repairs; adopted in §9.5 |
| `OD-10` | **Owner decision required**; convergent evidence now stronger |
| Corpus | `PASS`; registry clean; expansion frozen |
| `G1` | **Open** |

---

## 10. Second-order corrections (owner review, 2026-11-12)

### 10.1 `CORPUS-18` is a transcription, not an artifact

§9.1 upgraded the second review to **primary** because its text was supplied in full. That upgrade was premature. What exists in the workspace is **my transcription of the owner's message** — not the reviewer's file, and not byte-comparable to it. `CORPUS-18` now carries the corrected standing: **derived, secondary until compared with the original.**

The distinction is the whole point of `P-6`: *text arriving in a conversation* and *an artifact filed in the record* are different evidence classes. Reading a report is not the same as possessing it.

### 10.2 The scope formula — both of us wrote 11 for different reasons

| Source | Expression | Members correct? |
|---|---|---|
| Packet v1 | 6 + 5 − 3 = 8 | no — invited the wrong total |
| Packet v2 | 6 + 5 = 11 | no — "6 flagged" overstated R1's set |
| §9.2 (corrected) | **5 + 8 − 2 = 11** | **yes** |
| Owner review | 6 + 8 − 3 = 11 | no — right total, wrong members |

Both the corrected formula and the owner's reach 11; only `5 + 8 − 2` matches the actual review sets (`CORPUS-13` §2 = 5 units; Group A = 8; overlap = `F3-03`, `F3-04`). Recorded because a total that agrees for the wrong membership is the failure mode this corpus keeps producing.

### 10.3 `F7-01` — clause counting is rejected as the boundary

The owner's concern is correct, and the proposal fails on its own test cases:

| Response | Should | Clause-count verdict |
|---|---|---|
| "Halla says open at dusk, no word to anyone. Bren says shut at dusk, open at dawn." | **pass** — brief presentation | **fail** — two clauses per account |
| "Maybe the wind took it; Halla was elsewhere anyway." | **fail** — invented detail | pass under a loose counter |
| "Bren probably did it." | **fail** — adjudication | pass |

Clause boundaries depend on punctuation and paraphrase, so the counter penalises ordinary grammar while missing the actual defects. **Adopted instead:** the reviewer's semantic boundary — *extended reconstruction: detail neither account contains, or adjudication of which account is true.* Clause count is dropped, not retained as an aid.

**Principle recorded:** a boundary must be testable against negative *and* positive examples before it is proposed. This one was proposed before that test, and the test killed it.

### 10.4 `F7-04` — my floor would have created the defect it was fixing

`F7-04`'s restraint mode is **`none`**, so `C-4` (scoped to `R-2`/`R-3`) never fires on it. My proposal — mandatory → "not contradicted, suspicion not reasserted" — would have made a **silent or evasive reply pass**, with no mechanical finding to catch it: the `F6-04` class reappearing in a unit the checker cannot see.

The reviewer's repair (narrow the variation: "regretted or brisk subject change **that acknowledges the withdrawal**") keeps acknowledgement in the floor. **Adopted, and the reason is stronger than "less invasive":** mine was mechanically invisible and would have passed avoidance.

**Generalisable lesson:** proposals must be checked against the item's restraint mode *and* against the checker's scope, because a defect introduced outside a check's coverage is the most expensive kind.

### 10.5 Status delta

| Area | Change |
|---|---|
| `CORPUS-18` standing | Downgraded to derived transcription; **secondary** until compared with the original |
| Scope formula | `5 + 8 − 2 = 11`; owner's `6 + 8 − 3` reaches the same total from wrong members |
| `F7-01` proposal | Clause count **rejected**; reviewer's semantic boundary adopted |
| `F7-04` proposal | Mine **withdrawn** — it was silence-passable and outside `C-4`'s scope; reviewer's adopted |
| Freeze | Unchanged; `G1` open |

---

## 11. Governance corrections accepted (owner, 2026-11-12)

### 11.1 Membership accounting, not arithmetic

`5 + 8 − 2 = 11` is the **only** operative description of the packet selection. The earlier expressions are **void as descriptions** and appear above solely as defects of record — none may be reused, restated, or treated as an alternative route to the same 11.

**Principle:** *a correct total does not validate the sets used to produce it.* Two of the four expressions reached 11 with wrong membership. This is the same reporting failure as the `38/38` claim and the batch-total split, and it is the mechanism `P-6` mostly guards against: the check that the number is right must be separate from the check that the set is right.

### 11.2 What `PASS` means

`check_corpus.py` passing means the corpus satisfies **its declared mechanical checks** — presence, `C-1`–`C-4`, `U-1`. It does **not** mean any semantic repair is validated, and it does not mean the corpus is ready to seal. Six of eleven reviewed units were field-complete and unscorable; the checker reported none of them. **Presence is not validity, and `PASS` is not readiness.**

Recorded as a standing limitation on every future `PASS` in this programme.

### 11.3 The repair-review rule

**Every proposed repair must be checked against two things:** the item's **restraint mode**, and the **actual coverage of the mechanical checks**. A checker cannot protect against a defect outside its declared scope.

`F7-04` is the worked case: restraint `none`, so `C-4` is silent by design, and a widened floor there would have passed avoidance with no finding possible.

### 11.4 Repairs are reviewed as whole records

Proposals are **not** field edits. A field-level diff can read as correct while the record around it describes a different item — exactly the whole-record failure behind `F4-03` and `F7-04`. The review unit is the complete record with the proposal marked as proposed, and the next artifact is that packet.

### 11.5 Boundaries must survive paraphrase

A proposed boundary is tested against: concise faithful paraphrase (must pass), compound sentences and punctuation variation (must pass), invented detail (must fail), and adjudication or unwanted reconstruction (must fail). `F7-01`'s clause count failed three of four.

### 11.6 Status

| Area | State |
|---|---|
| `F3-04`, `F9-04` | **Blocked** — `OD-10` |
| `F7-01` | Semantic boundary adopted; clause count void |
| `F7-04` | Acknowledgement retained in the mandatory floor |
| Five accepts | Unsealed |
| 34 untouched units | Untouched |
| Checker | `PASS` = mechanical only |
| Corpus expansion | **Frozen** |
| `G1` | **Open** |
