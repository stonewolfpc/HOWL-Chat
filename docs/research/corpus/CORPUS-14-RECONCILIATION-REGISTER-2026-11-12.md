# Group B Post-Repair Reconciliation, Defect Register, and Escalation

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Status:** review **reconciled**; **four items require rework**; **one accepted on item validity**; **no item retired**; **`G1` open**; **one finding requires owner direction before any edit** (§6).
**Date:** 2026-11-12
**Inputs:** `CORPUS-13` (verdicts, verbatim), `CORPUS-12` (packet v2), `CORPUS-11` (packet v1, preserved), `CORPUS-08` (repair proposals), `CORPUS-09` (repair manifest), `CORPUS-07` (adjudication, `OD-4`).
**Effect:** this register **does not edit any item**. It records what the review found, separates repair-caused from pre-existing defects, and stops at the boundary the work order sets.

---

## 0. Answer to the commission's question

The work order asks whether the approved repairs **corrected the identified defects without introducing new defects, changing the intended tests improperly, or weakening the registered evaluation standards.**

| Question | Answer | Basis |
|---|---|---|
| Were the five identified defects corrected? | **Yes for all five, in the field each repair targeted.** | §2 — field-level before/after against the pre-repair snapshots |
| Did any repair leave the item inconsistent with itself? | **Yes — `F6-04`.** | `Required` and `Notes` were not updated and now contradict the repaired `mandatory`. §3 `RD-1` |
| Did any repair introduce a new scoring ambiguity? | **Yes — `P-01`**, two of them. | §3 `RD-2`, `RD-3` |
| Did any repair introduce a corpus-wide gap? | **Yes — Group A's 15 boundaries have no failure classes.** | §3 `RD-4` |
| Did any repair change an intended test improperly? | **No.** | No restraint mode changed; no split changed; no required evidence removed |
| Was a registered standard weakened to make something pass? | **No.** | §7 — every file that encodes a standard is byte-identical |

**One-line verdict.** The repairs did what they claimed; the items still are not acceptable, for reasons that mostly predate the repairs. `F6-04` is the exception: its repair is incomplete in a way that created a contradiction.

---

## 1. How cause was attributed

Blind reviewers could not say whether a finding was pre-existing or repair-caused, and were not asked to. Cause is attributed here, by field-level comparison against the preserved snapshots in `docs/research/corpus/snapshots/`. **Attribution is not a reviewer verdict and is labelled as this register's own finding throughout.**

Method: for each item, every registered field was diffed pre-repair → post-repair. A finding is:

- **Pre-existing** — the field it cites is byte-identical in both versions.
- **Repair-introduced** — the field it cites was added or altered by the repair, or the finding arises from the new field conflicting with an unaltered one.
- **Surfaced** — the field is unchanged but the repair's added evidence made a latent problem visible.

### 1.1 Field-level diff (evidence)

| Item | Fields changed by the repair | Fields byte-identical |
|---|---|---|
| `F3-03` | `Scenario` (name added), `mandatory` (name added), `over-exposure` (added) | `Split`, `Restraint`, `Turns`, `Required`, `Forbidden`, `variation`, `never scored`, `prohibited`, `If failed` |
| `F3-04` | `Scenario` ("The scar came from the fever; the current turn neither confirms nor denies it"), `Turns` ("No." → "The scar?"), `over-exposure` (added) | `Split`, `Restraint`, `Required`, `Forbidden`, `mandatory`, `variation`, `never scored`, `prohibited`, `If failed` |
| `F7-01` | `Scenario` (both accounts stated) — **only field changed** | everything else, including `Restraint: R-6`, `mandatory`, `variation`, `prohibited`, `over-exposure`, `If failed` |
| `P-01` | `Scenario` (six facts enumerated), `mandatory` ("using at most one of the facts listed in the scenario"), `over-exposure` ("more than one listed fact, or any timing detail beyond the month") | `Split`, `Restraint`, `Required`, `Forbidden`, `variation`, `never scored`, `prohibited`, `If failed`, `Notes` |
| `F6-04` | `mandatory`, `variation`, `prohibited`, `over-exposure` | `Split`, `Restraint`, `Scenario`, `Turns`, **`Required`**, `Forbidden`, `never scored`, `If failed`, **`Notes`** |

**Consequence:** three of the five repairs ( `F3-03`, `F3-04`, `F7-01` ) changed only `Scenario`/`Turns`/`mandatory`/`over-exposure`. Any finding citing `Required`, `Forbidden`, `variation`, `prohibited`, `never scored`, or `If failed` in those items is **pre-existing by construction** — those fields could not have been changed, and were not.

---

## 2. Per-item reconciliation

Original verdicts are in `CORPUS-13` and are not restated here as mine. This section records, for each item, whether I agree, and what the reconciliation resolves to.

### 2.1 `F3-03` — verdict `accept` (round 2)

| Finding | Cause | Reconciliation |
|---|---|---|
| Targeted defect: mandatory required a name absent from evidence | **Corrected** | `Scenario` and `mandatory` both now name Ilsa. Round-2 reviewer: "the entity is established in-field (name + event, both explicit)". Agreement. |
| D1: over-exposure has no failure class in `If failed` | **Repair-introduced** (Group A boundary added, `If failed` untouched) | **Registered as `RD-4`, not repaired here.** Shared by all 15 Group A units, so it is a system defect, not an `F3-03` defect. |
| D2: minimum sufficient "wedding context" undefined | Pre-existing (`Required` unchanged) | Registered as `RD-5`. Low severity. |
| D3: fabricated-sibling failure has no class | Pre-existing (`If failed` unchanged) | Folded into `RD-4`/`RD-5`. |
| D4: `acceptable_range` has no anchors | Pre-existing, corpus-wide | Registered as `RD-6`. |

**Reconciled status: `accept` on item validity, with four recorded patches outstanding.** The item is the only one of the five that is scoreable as registered. Its acceptance does **not** seal it — sealing is an owner act (`CORPUS-05` §5).

**Recorded disagreement with round 1.** Round 1 called the `Restraint: none` + over-exposure combination "the blocking defect". That reading was produced by a briefing that omitted the taxonomy. Both texts are preserved (`CORPUS-13` §3.1, §3.2). **The finding is void as stated; the reviewer's other findings stand and round 2 re-derived them independently.**

### 2.2 `F3-04` — verdict `rework`

| Finding | Cause | Reconciliation |
|---|---|---|
| Targeted defect: the turn answered the question ("No." to "the fall?") | **Corrected** | Round reviewer: "the leak is gone… the fever is not stated." Agreement. |
| Retrieval is optional; blank abstention passes | **Pre-existing** (`mandatory` unchanged; the uncertainty disjunct pre-dates the repair) | Agreement with the reviewer. Registered as `RD-7`. |
| `Required` vs `mandatory` conflict | Pre-existing | Agreement. Folded into `RD-7`. |
| `mandatory` vs `Forbidden` on "merged" | Pre-existing | Agreement. Folded into `RD-7`. |
| `variation: any wording` is a null constraint | Pre-existing | Agreement. Folded into `RD-7`. |
| No registered turn under evaluation | Pre-existing | Agreement. Folded into `RD-8` (schema gap). |
| Over-exposure "in full" unmeasurable | **Repair-introduced** (boundary is new) | Registered as `RD-9`. |

**Reconciled status: `rework`.** Attribution note: the repair **succeeded at its own stated purpose** and **did not cause the item's rejection**. The rejection is caused by a defect present before the repair and never addressed. Saying otherwise would misattribute the failure.

**Where I differ, and only in attribution.** The reviewer wrote that "de-leaking succeeded, but at the cost of vacuity". The repair cost nothing: `Turns` changed from one content-free persona line to another, and the retrieval cue was never in that line — it is in `Character: "That scar — the fall?"`, which the repair did not touch. The vacuity the reviewer measures is a property of the item's design, not a cost the repair paid. Both readings are preserved; the difference does not change the verdict.

### 2.3 `F7-01` — verdict `rework`

| Finding | Cause | Reconciliation |
|---|---|---|
| Targeted defect: both accounts required, neither content stated | **Corrected** | Round reviewer: accounts are "supplied, and distinguishable… Good." Agreement. |
| Over-exposure is not independent of the mandatory (M-2 cannot bite) | **Pre-existing** (`over-exposure` byte-identical) | Agreement. Registered as `RD-10`. |
| "survive" undefined | Pre-existing | Agreement. Folded into `RD-10`. |
| `prohibited: a confident culprit` anchors on a never-scored attribute | Pre-existing, and **rule-level** | Agreement. Registered as `RD-11` — this is a conflict between a registered rule and an authored field, not a wording error. |
| Incompatibly asserted, not entailed ("at dusk" may not be one instant) | **Surfaced** — the added accounts made the timing question visible; the fields that penalise the consistent reading are unchanged | Agreement that the fields can penalise the best-supported reading. Registered as `RD-10`. |
| "a fourth account" undefined against an inference | Pre-existing | Agreement. Folded into `RD-10`. |
| Probable co-active `R-1` mode | Pre-existing; **`ambiguity`, not adopted** | Recorded in §5. The blind reviewer inferred a co-active mode from scenario shape. The taxonomy requires exactly one tag (`A-2`), so a second co-active mode is either a re-tag or a scenario change. **Neither is authorised here.** |
| Failure taxonomy undiscriminating | Pre-existing | Folded into `RD-4`. |

**Reconciled status: `rework`.** The repair corrected its defect and caused none. Every remaining finding cites a field the repair did not touch.

### 2.4 `P-01` — verdict `rework`

| Finding | Cause | Reconciliation |
|---|---|---|
| Targeted defect: the boundary counted "six facts" that were never listed | **Corrected** | Round reviewer: "M-1 is met: boundary is explicit and countable". Agreement — the defect `CORPUS-06` §3.2 identified is genuinely gone. |
| D1: blanket silence passes | **Pre-existing** (`variation: "a silence"` unchanged; `Restraint: R-1` unchanged) | Agreement. Registered as `RD-12`. |
| D2: `mandatory` "at the scale of one gentle question" is a length criterion the item declares never scored | Pre-existing | Agreement. Folded into `RD-12`. |
| D3: the boundary counts omission, not restraint; unlisted disclosure passes | **Intensified by the repair** — the boundary is now countable, so its wrong axis is now operative rather than decorative | **Partial disagreement**, §5. Registered as `RD-13`. |
| D4: fact (1) "the sibling's name" has an ambiguous referent (the lost sibling or the arriving one) | **Repair-introduced** — the numbered facts are new | Agreement. Registered as `RD-2`. |
| D5: `Forbidden` "the timeline" vs permitted "the month" | **Repair-introduced** — "beyond the month" is new | Agreement. Registered as `RD-3`. |
| D8: fact values unregistered, so "any timing detail beyond the month" cannot be checked | **Repair-introduced** — there were no facts to value before | Agreement. Folded into `RD-2`. |
| D6 `presented poorly`, D7 no exemplars, D9 gesture unscorable | Pre-existing | Folded into `RD-6`, `RD-8`. |

**Reconciled status: `rework`.** This is the one item where the repair both **fixed the stated defect and added two new ambiguities**. The fix is real; so is the cost. Recorded as such rather than as a clean success.

### 2.5 `F6-04` — verdict `rework`

| Finding | Cause | Reconciliation |
|---|---|---|
| Targeted defect: no positive required behaviour; non-use, unavailability, avoidance indistinguishable | **Corrected in `mandatory`** | The new `mandatory` is perceptible and unannounced, and `variation` now requires acts. Agreement that the *mandatory* defect is fixed. |
| `Required: the preference, if it is used at all` contradicts the new mandatory | **Repair-introduced contradiction** — `Required` is unchanged and was consistent with the old `mandatory: "nothing"`; against the new mandatory it licenses the exact thing the repair forbids | **Agreement, and the most serious finding in this review.** Registered as `RD-1`. |
| `Notes: "Deliberately permits non-use. Trivia that is available need not be used"` now false | **Repair-introduced** — `Notes` is unchanged and was true under the old mandatory | Agreement. Folded into `RD-1`. |
| No anti-coincidence guard | Partly pre-existing (`never scored: whether the memory appears` unchanged), partly a consequence of the new silent-mandatory | Registered as `RD-14`. |
| `mandatory` presupposes an unregistered menu / "hottest" is spice-or-temperature | Pre-existing in `Scenario` | Folded into `RD-14`. |
| Graded role unstated | Pre-existing | Folded into `RD-8`. |
| `over-exposure` duplicates `prohibited` | **Repair-introduced** (both were rewritten to near-identical wording) | Registered as `RD-15`. Minor. |

**Reconciled status: `rework`.** The repair is **incomplete, not wrong**. Four of the six fields it needed to touch were changed; `Required` and `Notes` were left behind, and the item now contradicts itself in two places. Under `M-2` and `A-2` the mandatory line governs, but a scorer reading `Required` alone reaches the opposite verdict — which is precisely the rating disagreement `CORPUS-05` `C-3` exists to prevent.

---

## 3. New defect register

Registered without editing any item. `RD-n` are review-derived; they are **not** item IDs and do not compete with the corpus ID scheme.

| ID | Defect | Items | Cause | Severity | Blocks accept? |
|---|---|---|---|---|---|
| `RD-1` | `Required` and `Notes` left stale by the `F6-04` repair; both now contradict `mandatory` and `R-3` | `F6-04` | Repair-introduced | **Critical** | **Yes** |
| `RD-2` | Numbered fact list has an ambiguous referent (fact 1) and no fact **values**, so "one listed fact" cannot be counted consistently | `P-01` | Repair-introduced | High | **Yes** |
| `RD-3` | `Forbidden: "the timeline"` versus `over-exposure` permitting "the month" — one dated mention is undefined | `P-01` | Repair-introduced | Medium | **Yes** |
| `RD-4` | Group A added 15 scored over-exposure boundaries and no item's `If failed` taxonomy gained a class for an over-exposure failure. `R-9` requires exactly one failure class per failed item | `F3-03`, `F3-04`, `F4-03`, `F5-03`, `F5-04`, `F6-03`, `F7-03`, `F7-04`, `F8-03`, `F8-04`, `F9-03`, `F9-04`, `F9-05`, `F1-03`, `F2-03` | Repair-introduced, systemic | High | No for `F3-03`; aggravating elsewhere |
| `RD-5` | Minimum sufficient `Required` content undefined against the over-exposure boundary | `F3-03` | Pre-existing | Low | No |
| `RD-6` | `acceptable_range` is a type label with no registered anchors or exemplars | corpus-wide | Pre-existing | Medium | No, but limits reproducibility |
| `RD-7` | `mandatory`'s uncertainty disjunct permits a pass with no retrieval; `Required`/`mandatory`/`Forbidden` disagree about whether the winters must be distinguished | `F3-04` | Pre-existing | **Critical** | **Yes** |
| `RD-8` | No item registers which turn is scored; graded role is implicit | `F3-04`, `F6-04` (schema-wide) | Pre-existing | Medium | Contributes |
| `RD-9` | Over-exposure "in full" has no unit and collides with `Required: both winters distinguished` | `F3-04` | Repair-introduced | Medium | Contributes |
| `RD-10` | Over-exposure is collinear with the mandatory; "survive" undefined; "fourth account" undefined against inference; incompatibility not entailed | `F7-01` | Pre-existing | High | **Yes** |
| `RD-11` | `prohibited` anchored on "confident", while registered rules declare confidence **never** a failure signal | `F3-04`, `F7-01` | Pre-existing; **rule-level conflict** | High | Contributes |
| `RD-12` | `variation: "a silence"` coexists with a mandatory acknowledgment; `mandatory`'s "at the scale of one gentle question" is a length/tone criterion the item declares never scored | `P-01` | Pre-existing | High | **Yes** |
| `RD-13` | The countable boundary measures *omission* along one axis only; unlisted disclosure passes | `P-01` | Pre-existing axis, intensified by repair | High | **Yes** |
| `RD-14` | No anti-coincidence control; `mandatory` has no referent object (no menu, no dish); "hottest" is spice-or-temperature | `F6-04` | Mixed | High | **Yes** |
| `RD-15` | `over-exposure` and `prohibited` were rewritten to near-identical wording and can no longer be scored independently | `F6-04` | Repair-introduced | Low | No |

**Nothing above was repaired.** `CORPUS-11` §7 forbids modifying the five items during review, and `RD-1`, `RD-4` and `RD-11` each need a scope decision, not an edit.

---

## 4. Verdict matrix, reconciled

| Item | Reviewer verdict | Reconciled | Target defect fixed? | New defect introduced? | Disposition proposed |
|---|---|---|---|---|---|
| `F3-03` | `accept` | **`accept`** | Yes | Minor, systemic (`RD-4`) | Accept for scoring; patch `RD-4`/`RD-5` at Group A level |
| `F3-04` | `rework` | **`rework`** | Yes | Minor (`RD-9`) | Rework under new authorisation (`RD-7`, `RD-11`) |
| `F7-01` | `rework` | **`rework`** | Yes | None | Rework under new authorisation (`RD-10`, `RD-11`) |
| `P-01` | `rework` | **`rework`** | Yes | **Yes** (`RD-2`, `RD-3`) | Rework under new authorisation (`RD-12`, `RD-13` + new) |
| `F6-04` | `rework` | **`rework`** | Partly | **Yes** (`RD-1`, `RD-15`) | Rework; `RD-1` needs owner direction first |

```text
accepted on item validity : 1  (F3-03)
rework                    : 4  (F3-04, F7-01, P-01, F6-04)
retired                   : 0
repairs confirmed correct : 5 / 5 targeted defects
repairs incomplete        : 1  (F6-04)
repairs with side-effects : 2  (P-01, F6-04)
```

**No retirement is proposed.** Every `rework` is bounded and every item retains a mechanism worth keeping; each reviewer said so explicitly. `CORPUS-05` §5 step 4 reserves retirement for items that **cannot** be reconciled, and none of these are in that class.

---

## 5. Preserved disagreements and ambiguities

Consensus is not correctness (`CORPUS-05` §5). These are recorded rather than resolved.

| # | Point | Reviewer position | This register's position | Status |
|---|---|---|---|---|
| 1 | `F3-03` round 1, `Restraint: none` vs over-exposure | "the blocking defect… two fields give opposite guidance" | The registered interpretation rule makes the combination consistent; the finding was produced by a briefing omission and is void **as stated**. The reviewer's separable point — that the boundary has no failure class — survives in round 2 as D1. | **Resolved against round 1 on the contradiction; round-2 finding retained** |
| 2 | `F3-04`, "the repair cost vacuity" | "de-leaking succeeded, but at the cost of vacuity" | The retrieval cue was never in the removed line; the repair cost nothing. The vacuity is a design property. | **Disagreement on attribution only; verdict unaffected** |
| 3 | `P-01` `D3`, "the boundary measures omission, not restraint" | "It is countable — it is just not measuring restraint" | Partly right and partly overstated: a fact-count boundary **does** score over-exposure on the axis it names. It is silent on undisclosed-but-true material outside the list, which is a coverage gap, not a total miss. | **Partial disagreement; both positions preserved** |
| 4 | `F7-01`, co-active `R-1` mode | "probable mode mis-assignment" | Plausible but not decidable from the fields. `A-2` permits one tag only; changing the tag or the scenario is outside this authorisation. | **Ambiguous — retained, not adopted** |
| 5 | `F7-01`, "at dusk" as one instant | "Could not verify" | The incompatibility rests on an unstated same-instant assumption. This is a real gap in the evidence, and it is the kind of thing the repair's added accounts exposed rather than caused. | **Ambiguous — registered in `RD-10`** |
| 6 | `F6-04` silence | "ambiguous (fail)" | The mandatory and `R-3` govern, so silence fails; but `Required` says the opposite, so a scorer can reach either verdict from the item alone. The contradiction, not the silence question, is the defect. | **Resolved as `RD-1`** |

**Nothing was harmonised to produce a cleaner table.** Items 2 and 3 are live disagreements that a later reviewer may settle against this register.

---

## 6. Owner decisions required before anything is edited

The work order says to stop and request direction when a finding needs authority the review does not carry. Five do.

| # | Decision | Why it cannot be made here |
|---|---|---|
| `OD-7` | **Authorise a second Group B rework round** for `F3-04`, `F7-01`, `P-01`, `F6-04` against `RD-7`, `RD-10`–`RD-15`. | `CORPUS-11` §7 forbids editing the reviewed items. A second rework is new authorisation, and each item's intended behaviour changes, so each needs its own re-review (`CORPUS-11` §5). |
| `OD-8` | **`F6-04` stale-field repair specifically** — restore consistency between `mandatory` and `Required`/`Notes`. | This is not a new design choice: `OD-5` already fixed the behaviour. It is the missing half of an approved repair. It still edits a reviewed item, so it needs explicit permission. **Recommended: approve as a correction, not a redesign.** |
| `OD-9` | **Extend `If failed` taxonomies** to cover over-exposure failures across the 15 Group A units and `F6-04`. | `R-9` requires one failure class per failed item. Group A was classified "metadata only, no re-review required" (`CORPUS-07` §3.1); this finding challenges that classification, and reversing it is an owner call. |
| `OD-10` | **Resolve `RD-11`**: either "confident" stops being a scoring qualifier, or the never-scored list stops naming confidence. | This is a **rule conflict inside the standard**, not an item defect. Fixing it by editing `F3-04`/`F7-01` would silently change what every other item's never-scored clause means. `MEMORY-09` is out of bounds for this review. |
| `OD-11` | **Settle reviewer sufficiency**: does one non-human reviewer satisfy `P-3` for sealing? | `CORPUS-05` §5 step 1 treats a single reviewer as lower confidence. Whether that is enough to seal is an owner judgement no agent should make. |

**Also flagged, not needing a decision yet:**

- **`RD-6`** (no reference anchors) and **`RD-8`** (turn under evaluation unregistered) are corpus-wide schema gaps. They are cheaper to fix once, before authoring resumes, than per item.
- **Challenge-numbering collision.** `CORPUS-05` uses `C-1`–`C-5`; `CORPUS-12` §4 uses `C-1`–`C-7` with different meanings. Two schemes now exist in the record. Registered as `ND-10`; recommend adopting the seven-challenge scheme everywhere and renumbering `CORPUS-05` in a dated amendment rather than by silent edit.

---

## 7. Standards not weakened — verification

The commission asks whether the review weakened a registered standard. Every file that encodes one was hashed before and after; none changed.

| File | Encodes | Status |
|---|---|---|
| `CORPUS-02-BATCH-1-ITEMS` | items | **unchanged** — `fc39ac6f…` |
| `CORPUS-03-BATCH-2-ITEMS` | items | **unchanged** — `638bcaf6…` |
| `CORPUS-04-R6-UNCERTAINTY-ITEMS` | items | **unchanged** — `3d34f053…` |
| `check_corpus.py` | the checker | **not modified**; not re-run to suppress a finding |
| `MEMORY-09` | acceptance criteria, rubrics, gates | **not modified** — `f77e47a9…` |
| `CORPUS-05` | review protocol | **not modified** |
| `CORPUS-08` | repair proposals, including the before/after rationale | **not modified** — `1b68bd4d…` |
| `CORPUS-09` | repair manifest and post-repair checker output | **not modified** — `f2dc7e6e…` |
| `CORPUS-10` | checker fixture validation | **not modified** — `f67c9e36…` |
| `CORPUS-06` | the Muse review record and its secondary status | **not modified** — `6fe3b8e6…` |
| `CORPUS-11` | review packet v1 and its hash | **not modified** — `b51c4b68…` |
| `snapshots/*` | pre-repair corpus and baseline output | **not modified** |

**New files only:** `CORPUS-12`, `CORPUS-13`, `CORPUS-14`, plus registry rows for them.

```text
items modified          : 0
items added / removed   : 0 / 0
split changes           : 0
restraint re-tags       : 0
checker changes         : 0
MEMORY-09 changes       : 0
models loaded           : 0
experiments run         : 0
corpus sealed           : no
G1 closed               : no
```

---

## 8. Updated corpus status

| State | Units |
|---|---|
| Authored, mechanically checked (`PASS`, 0/0/0) | **45** |
| Reviewed at least once (Muse, relayed/secondary) | **45** |
| Repaired under authorisation | **18 units / 20 findings** |
| **Independently reviewed post-repair** | **5** — `F3-03`, `F3-04`, `F7-01`, `P-01`, `F6-04` |
| — accepted on item validity | **1** (`F3-03`) |
| — require rework | **4** (`F3-04`, `F7-01`, `P-01`, `F6-04`) |
| — retired | **0** |
| Reviewed exactly once and **not** repaired or re-reviewed | **27** — do not treat as validated |
| Ambiguities from `CORPUS-07` §2 still unapplied | **12** |
| Rule conflicts open | **3** registered (`F3-01`, `F7-02`, `R6-01`) + **1** new (`RD-11`) |
| Held-out sealed (`S-2`) | **no** |
| Corpus sealed | **no** — 45 units against an ≥80 floor |
| `G1` | **open** |
| `G2`–`G5` | **closed** |
| Authoring | **paused** (`P-5`) — and now paused for a second reason: 4 of the last 5 reviewed items need rework |

**What may be said now.** The five repairs each corrected the defect they targeted. One repair (`F6-04`) is incomplete and left the item self-contradictory. Two repairs (`P-01`, `F6-04`) introduced side-effects. One item (`F3-03`) is scoreable as registered.

**What may not be said.** That the corpus is valid, adequate, or ready. That `G1` is near closure on quality. That a single non-human review settles `P-3`.

---

## 9. File identity

Hashes of the documents this commission produced, and of the artifacts it was required to preserve. Verification command: `python -c "import hashlib,pathlib,sys;print(hashlib.sha256(pathlib.Path(sys.argv[1]).read_bytes()).hexdigest())" <file>`.

| File | SHA-256 |
|---|---|
| `CORPUS-12-GROUP-B-REVIEW-PACKET-V2-2026-11-12.md` | `ddf5658c2828d03db9b350ef8a612bbb8d704d173560107b7d403f143700b4aa` |
| `CORPUS-13-INDEPENDENT-REVIEWER-REPORT-2026-11-12.md` | `e6b969be8c69f622ca45ea7728c668667a92283956aa601d330210d337ba6d92` |
| `CORPUS-14-RECONCILIATION-REGISTER-2026-11-12.md` | `self — not recordable in this file` |
| `CORPUS-11-GROUP-B-REVIEW-PACKET-2026-11-12.md` (preserved) | `b51c4b680dd1f300bdf659228ed4def462f0047ab8fc3ebe39fa419ea7d826aa` |
| `CORPUS-08-REPAIR-PROPOSALS-2026-11-12.md` (preserved) | `1b68bd4d9d4313c0447363f629d052df1b7e4cb40e110859a4659045e4868103` |
| `CORPUS-09-REPAIR-EXECUTION-AND-VERIFICATION-2026-11-12.md` (preserved) | `f2dc7e6ed66d8cdc546baf0363ab7f8eefb743753c0a31da5bd0f6e08bd55fcb` |
| `CORPUS-10-CHECKER-FIXTURE-VALIDATION-2026-11-12.md` (preserved) | `f67c9e365e17c4cf7daf4ce202cf59402c4a852dc4803aa693b691abb7b5549b` |
| `CORPUS-06-REVIEW-FINDINGS-AND-CORRECTION-REGISTER-2026-11-12.md` (preserved) | `6fe3b8e6c0d5f846476a8df8d7e95d325b88193bb8cfec5c57273efbe718252b` |
| `CORPUS-02-BATCH-1-ITEMS-2026-11-12.md` (unmodified) | `fc39ac6fd5a3d4f9bce8d255353ae2f103058f49305c96a2d6ad7f7accb37a8d` |
| `CORPUS-03-BATCH-2-ITEMS-2026-11-12.md` (unmodified) | `638bcaf659b366f858ffc74d192b4424da40df048f8e65c8c640978b2b97907d` |
| `CORPUS-04-R6-UNCERTAINTY-ITEMS-2026-11-12.md` (unmodified) | `3d34f0530699fd333fe89b08e3dbbbfab4db7cc2908781a9a25790a398add434` |
| `MEMORY-09-MEMORY-EVALUATION-PREREGISTRATION-2026-11-12.md` (unmodified) | `f77e47a98685e66028e5aa23e5be45d7cda70b7c0cf5f9c443589b76f252f8da` |

---

## 10. Required follow-ups, in order

1. **Owner decides `OD-7`–`OD-11`** (§6). Nothing is edited before that.
2. **`OD-8` first** if approved — `F6-04` is the only item whose *own* repair is incomplete, and it is cheaper to close before the larger rework round.
3. **Then the second rework round** for `F3-04`, `F7-01`, `P-01`, `F6-04`, with a fresh independent review of each (`CORPUS-11` §5), and this register preserved unedited beside it.
4. **`RD-4` at Group A scale** — one decision covering 15 units, not fifteen edits.
5. **`RD-6`/`RD-8` before authoring resumes** — fixing the schema once beats patching it 45 times.
6. **Do not resume authoring toward the 80 floor** until step 3 closes; 4 of 5 reviewed items needing rework is a signal about items, not about count.

---

## 11. Boundaries observed

- No item modified. No item added, removed, split, or retired.
- `MEMORY-09` acceptance criteria untouched.
- The checker was not modified, and no finding was suppressed by changing what is measured.
- No model loaded, no experiment run, no source, schema, runtime, or application data changed.
- `G1` not sealed; corpus validity not declared; `G2` not opened.
- Author intent was never supplied to a reviewer and never used here as evidence.
- Reviewer verdicts were preserved verbatim, including the two this register disputes.
