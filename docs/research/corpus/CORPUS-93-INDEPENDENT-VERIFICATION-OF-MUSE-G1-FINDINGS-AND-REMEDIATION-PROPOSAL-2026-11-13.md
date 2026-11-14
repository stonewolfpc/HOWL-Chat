# CORPUS-93 — Independent Verification of Muse's 30 `G1` Findings, and Dependency-Ordered Remediation Proposal

**Execution date:** 2026-11-13
**Report type:** Independent verification and remediation proposal (non-operative, review-only)
**Scope:** Verify the 30 findings of `MEMORY-09-G1-READINESS-01` against the actual source artifacts; investigate the missing-source-conversation finding first; reconcile the corpus inventory and all six restraint-mode memberships; verify the `CORPUS-91`, `CORPUS-92` and registry integrity findings; classify existing work as reusable or rebuild; propose a dependency-ordered sequence.

**What this report does not do.** It does not amend `MEMORY-09`, `CORPUS-01`, any item file, any tag, any split, any checker, or any registry row. It does not adopt a finding. It does not decide the `R-6` silence/perceptibility question, and it does not begin implementation. The only write performed is **this file**. The document registry was **not** modified (see §12).

---

## 0. Reviewer provenance and independence

### 0.1 Three-field form

`CORPUS-85` §1, `CORPUS-86` §1, `CORPUS-90` §0 and `CORPUS-91` §0.1 established a three-field provenance form. This report uses it rather than a single label. `[GT]`

| Field | Value | Basis |
|---|---|---|
| **Owner-assigned reviewer designation** | **None asserted.** No designation was assigned for this verification. | Absence of an instrument. `[RI]` |
| **Executing identity** | **Buffy / Codebuff**, Freebuff agent environment, working tree `C:\Users\StoneWolfPC\HOWL_Chat` | Supplied by the environment; observed. `[OBS]` |
| **Certified model identity** | **None asserted.** No model lineage is attested or inferred. | Absence of attestation is a fact about the environment. `[RI]` |

I deliberately assert **no** designation. `CORPUS-91` §0.1 records "Nemotron Lightning" assigned by owner instruction; `CORPUS-79` §0 records "Buffy" while `CORPUS-80` §8.1 and `CORPUS-86` §1 attribute `CORPUS-79` to "Muse"; `MEMORY-09-G1-READINESS-01` §2.1 records "Muse" for itself and flags that the label already collides. Adding a fourth unverifiable label to that pile would worsen the defect this report is asked to check. `[RI]`

### 0.2 Independence rating

**Environment-contaminated; task-level re-derivation genuine; not independent lineage.** `[RI]`

- Environment-level independence across the `Q1`/`Q2`/`Q3` lineage and this tree is **absent**. `CORPUS-80` §8 already rules that convergence is **0 of 2**, and `CORPUS-89` §8 repeats it. `[GT]` — verified at `CORPUS-80` L154/L201 and `CORPUS-89` L86/L128 `[OBS]`.
- Task-level independence: this assignment is a **verification of a verification**, not a fourth adjudication. Every figure below was re-derived from the item files, the specification, the checker and the registry — not from Muse's summary of them. Where I agree with Muse, I show the derivation; where I disagree, I show the number.
- I authored none of the inspected artifacts. `[OBS]`

### 0.3 Access limitations recorded, not worked around

- I did not read `CORPUS-78`, `CORPUS-79`, `CORPUS-81`, or `CORPUS-82`/`83` in full. Where a Muse finding depends on their text I say so and mark it unverified **by this review**. `[UNEST]`
- I did not read `MEMORY-01`, `-02`, `-05`–`-08`. `[UNEST]`
- No model was executed; no inference, probe or scoring run was performed. `[OBS]`

### 0.4 Label key

| Tag | Meaning |
|---|---|
| `[OBS]` | Observed — produced by a command or direct read in this session. |
| `[GT]` | Source-derived (governing) — quotation or strict paraphrase of an operative document, cited by file and line. |
| `[IA]` | Source-derived (artifact) — a determination about a specific item, record, or report. |
| `[CALC]` | Calculated — arithmetic from source figures, inputs shown. |
| `[RI]` | Inferred — this reviewer's reasoning; not a quotation and not a measurement. |
| `[UNEST]` | Unestablished — the record does not determine the answer; the missing evidence is named. |
| `[PROP]` | Proposal — a recommendation with no operative force. |
| `[UNRES]` | Unresolved — requires an owner decision or missing evidence. |

---

## 1. Method and re-runnable checks

```bash
# Artifact identity
for f in <artifact>; do wc -l "$f"; wc -c "$f"; sha256sum "$f"; done

# Corpus census, derived from item headings + Restraint/Split/Abstention lines
python - <<'PY'   # parser reproduced in §3.2; a block counts as a unit only if it
                  # contains **Reference:** — the same is_test rule check_corpus.py uses
PY

# Checker behaviour
python docs/research/corpus/check_corpus.py

# Source-evidence presence
grep -rn "source_conversation_id\|turn_range" docs/research/corpus/*.md

# Registry integrity
grep -cE '^\| [0-9]+ \|' docs/DOCUMENT-REGISTRY.md
grep -oE '^\| [0-9]+ \|' docs/DOCUMENT-REGISTRY.md | grep -oE '[0-9]+' | sort -n | uniq -d
find docs -type f | wc -l
```

**The one parsing subtlety that changes every count.** `CORPUS-04` L53 opens a *container* heading ``## `R6-03` `` that carries its own `Family`/`Split`/`Restraint` line but no `**Reference:**`. It is a section header for four separately-scored variants, not a scored unit. `check_corpus.py` excludes it by requiring `**Reference:**` (its `is_test` flag, L60–100 `[OBS]`). A naive counter that keys on `**Restraint:**` alone reports **46 units and `R-6` 10**, which is wrong. Both `CORPUS-01` §7 ("Item headings: 42 · scored units: 45") and the checker's `45` are correct once this is handled. `[OBS]`

---

## 2. FIRST: the missing-source-conversation finding (F-04)

This is the finding that blocks the most work, so it is investigated first and at depth.

### 2.1 What the governing text requires

| Source | Requirement | Verified |
|---|---|---|
| `MEMORY-09` §3.1(4), L87 | "**Originals are preserved as evidence.** The corpus stores the source conversation untouched, alongside labels." | `[GT]` confirmed at L87 `[OBS]` |
| `MEMORY-09` §3.3, L127–139 | Per-item record includes `source_conversation_id`, `turn_range`, `required_evidence_ids[]`, `forbidden_evidence_ids[]` | `[GT]` confirmed at L130–131, L134–135 `[OBS]` |
| `MEMORY-03` §5.3, L229 | "Ground-truth ledger" | `[GT]` section exists `[OBS]` |

Note the field names precisely: the schema requires `required_evidence_**ids**[]` and `forbidden_evidence_**ids**[]`. IDs, not descriptions. `[GT]`

### 2.2 What exists

| Check | Result |
|---|---|
| Occurrences of `source_conversation_id` in any corpus document | **0** `[OBS]` |
| Occurrences of `turn_range` in any corpus document | **0** `[OBS]` |
| Conversation artifacts under `docs/` | **none** — the only `*conversation*` files are `AMENDMENT-02` architecture documents and `ROOT-01`, none of which is a corpus conversation `[OBS]` |
| Ground-truth ledger with fact IDs and spans | **none found** `[OBS]` |
| Evidence fields in items | Prose `**Required:**` / `**Forbidden:**` descriptions, not IDs `[IA]` |
| Positional references in items | Approximate prose: "~340 messages of shared history. Around message 90 the persona promised…" (`CORPUS-02` L17 `[GT]`) |

### 2.3 Verdict on F-04

**CONFIRMED, and at the strongest evidential level available.** The schema's own field names appear nowhere in the corpus. The gap is not a missing attachment or a moved file; the item class was never authored. `[IA]`

Three consequences, tested rather than asserted:

1. **No item is executable.** `required_evidence` cannot resolve to a span because there are no spans. `[IA]`
2. **The model-free half of the programme is also blocked, not only the model-dependent half.** `MEMORY-09` §6 allocates retrieval recall/nDCG, deletion propagation and token accounting to non-inference testing — but every one of those needs an addressable record. `[GT]`/`[RI]` This is the part of Muse's §1.1 and §8.1 argument I consider most underrated and I confirm it independently.
3. **The 45 specifications are not diminished by this.** They specify what a conversation must contain. That makes them *authoring briefs*, which is a durable asset (§7). `[RI]`

### 2.4 One correction to how the finding is usually stated

Muse's register entry F-04 says "Author 45 synthetic conversations". Under the corpus's own accounting there are **45 scored units against 42 headings**, because `R6-03a`–`d` are four scored units sharing one scene and one heading. `[GT]` (`CORPUS-01` §7, and §4's "42 item *headings*, 45 scored *units*"). `[OBS]` The conversation count is therefore **42 distinct scenes**, of which one yields four scored units, unless the owner chooses to author four separate conversation files. This matters for effort estimation, not for the finding, which stands unimpaired. `[RI]`

---

## 3. Corpus inventory and all six restraint-mode memberships

### 3.1 The census, independently derived

Derived from 45 scored blocks (`**Reference:**`-bearing) across the three item files. `[CALC]`/`[OBS]`

| Mode | Units | Count | Split T/H |
|---|---|---:|---|
| `R-1` Suppress mention | `F6-02`, `P-01` | 2 | 2 / 0 |
| `R-2` Influence tone | `F1-01`, `F1-02`, `F1-03`, `F6-01` | 4 | 2 / 2 |
| `R-3` Natural use | `F3-01`, `F6-03`, `F6-04`, `P-04` | 4 | 1 / 3 |
| `R-4` Recognise irrelevance | `F8-01`, `F8-03`, `F9-01`, `F9-03` | 4 | 3 / 1 |
| `R-5` Respect entitlement | `F2-01`–`F2-03`, `F5-01`–`F5-04`, `F9-02`, `F9-05`, `P-03` | 10 | 4 / 6 |
| `R-6` Preserve uncertainty | `F7-01`, `F9-04`, `R6-01`, `R6-02`, `R6-03a`–`d`, `R6-04` | 9 | 2 / 7 |
| `none` | `F3-02`–`F3-04`, `F4-01`–`F4-03`, `F7-02`–`F7-04`, `F8-02`, `F8-04`, `P-02` | 12 | 5 / 7 |
| **Total** | | **45** | **19 / 26** |

- Restraint-tagged: **33**. `none`: **12**. `[CALC]`
- Held-out share: 26/45 = **57.8%**, above the `S-4` 40% floor. `[CALC]`
- No unit carries two mode tags at item level; `A-2` holds at item level. `[OBS]`
- No duplicate item IDs. `[OBS]`

**This reproduces Muse §5.2 and `CORPUS-91` §1.4 and `CORPUS-92` §2.3 exactly, term for term.** Three independent derivations now agree with the item files, which is the strongest support any figure in this programme currently has. `[IA]`

### 3.2 Checker behaviour

```text
scored units: 45
split: tuning=19 heldout=26
declared in CORPUS-01: {'units': 45, 'tuning': 19, 'heldout': 26}
TOTALS MATCH
RESULT: PASS (0 fail, 0 warn, 0 unverifiable)
```

`[OBS]` Two readings must be kept apart, and I confirm Muse's separation of them:

1. **The totals reproduce.** `CORPUS-01` L166's "measures 19 tuning / 22 held-out … Not reproducible → not compliant" is **stale and false against the current files**. The four `R6-03` variants now carry explicit per-unit `Split:` lines. `[GT]`/`[IA]`
2. **`PASS` is not schema validity.** The checker verifies unique ID, `Split`, `Restraint`, required/forbidden presence, the five reference parts. It does **not** check `family`, `source_conversation_id`, `turn_range`, `expected_abstention`, `expected_character_perspective`, or the `*_ids[]` form of the evidence fields. `[OBS]`

### 3.3 The four membership representations

| Source | R-1 | R-2 | R-3 | R-4 | R-5 | R-6 | Sum |
|---|---:|---:|---:|---:|---:|---:|---:|
| `CORPUS-01` §3.1 rows L121–126 | 2 | 3 | 2 | 3 | 4 | 2 | 16 entries / 14 distinct |
| `CORPUS-01` L134 coverage line | 2 | 4 | 4 | 4 | 10 | **6** | **42** ✗ |
| `CORPUS-01` L245 status prose | — | — | — | — | — | **9** | — |
| Item tags (this session) | 2 | 4 | 4 | 4 | 10 | **9** | **45** ✓ |

`[CALC]` and `[GT]`.

**Verified findings.**

1. **The coverage line cannot be right.** 2+4+4+4+10+6 = 30; +12 `none` = **42**, not 45. `R-6` is the sole wrong term, and the same file's L245 says 9. `[CALC]`
2. **The rows are not memberships.** `R-1`'s row lists `P-03` and `F9-01`. `P-03` is tagged `R-5`; `F9-01` is tagged `R-4`. **Neither actual `R-1` unit (`F6-02`, `P-01`) appears.** This is misattribution in both entries, not omission. `[IA]` — re-read at `CORPUS-01` L121 `[OBS]`
3. **The rows violate `A-2`.** "Every item is tagged with exactly one mode" (L130) cannot hold while the table places `P-03` under `R-1` (L121) and `R-5` (L125), and `F9-01` under `R-1` (L121) and `R-4` (L124). `[IA]`
4. **The table omits 19 of the 33 tagged restraint units.** 14 distinct units are named in rows; all 14 are restraint-tagged; 33 − 14 = **19 unnamed**. `[CALC]`
5. **The five-of-six framing is imprecise, and `CORPUS-92`'s refinement is the accurate one.** All six rows list a subset; `R-1` misattributes both entries; the other five omit. `R-6`'s two listed units are genuinely `R-6`, so it contains no wrong entry — only 2 of 9. `[IA]` I adopt `CORPUS-92` §4.4's characterisation over both `CORPUS-80` P2's and Muse's summary.

**Source of truth, if one exists:** for *counts*, the item tags. For *meaning*, `CORPUS-01` §3.1's pass-condition column. For *authority*, **nothing** — no document states which representation controls. That is Muse F-09 and `CORPUS-81` U-2, and I confirm the absence. `[UNRES]`

### 3.4 Family distribution against the preregistered minimums

`MEMORY-09` L94–102 states the minimums; L104 states "**Total minimum: 80 scored items**, plus the twelve `MEMORY-04` fixtures … not replaced." `[GT]` — read at source `[OBS]`

| Family | Have | Required | Deficit | Held-out | ≥2 held-out? |
|---|---:|---:|---:|---:|---|
| `F-1` | 3 | 10 | 7 | 2 | yes |
| `F-2` | 3 | 10 | 7 | 2 | yes |
| `F-3` | 4 | 10 | 6 | 2 | yes |
| `F-4` | 3 | 10 | 7 | 2 | yes |
| `F-5` | 4 | 8 | 4 | 2 | yes |
| `F-6` | 4 | 8 | 4 | 2 | yes |
| `F-7` | 11 | 8 | **0 (+3 surplus)** | 8 | yes |
| `F-8` | 4 | 8 | 4 | 2 | yes |
| `F-9` | 5 | 8 | 3 | 2 | yes |
| `P` principle cases | 4 | mandatory, outside the 80 | 0 | 2 | n/a |
| **Family units** | **41** | **80** | **39 aggregate / 42 per-family** | **24** | — |

`[CALC]`. `R6-01`–`R6-04` are `Family: F-7` per `CORPUS-04`'s own field, which is why `F-7` reads 11 and not 4. `[GT]`

**A precision defect in Muse's F-06, found on verification.** `[IA]`

- `80 − 41 = 39`. That is the **aggregate** shortfall, and Muse's "39 family units are missing" is arithmetically true as an aggregate.
- But the **per-family deficits sum to 42**, because `F-7` holds a 3-unit surplus that a per-family minimum cannot absorb.
- Muse's F-06 corrective action says "Author 39 units". Under the same strict per-family reading Muse applies to `CORPUS-01` L180, 39 new units would reach a total of 80 while leaving three families still short. **The correct authoring requirement is 42 new family units, or 39 plus an explicit owner authorisation to reallocate `F-7`'s surplus** — which is precisely OD-C, still undecided.
- This is the mirror of the defect Muse attributes to `CORPUS-01`: a figure derived by offsetting `F-7`'s surplus. Muse's own §5.4 table makes the contradiction visible (deficits 42, headline 39), and its WP4 completion evidence — "Every family at minimum" — is the 42-figure, not the 39-figure. Muse's conclusion is unaffected; its number is.

**A second-order error in `CORPUS-01` that Muse did not name.** L180 reads "Author ≥38 more family items to reach the floor (42 of 80)". The corpus holds **41** family units, not 42. The `42` appears to be the *heading* count (42) applied to family units. So `CORPUS-01`'s floor arithmetic is wrong twice: wrong base (42 vs 41) and an unauthorised `F-7` offset. `[CALC]`/`[IA]`

**Total scored units when complete:** 80 family + 4 principle cases + 12 converted fixtures = **96**, under the strict reading. `[CALC]`

**Held-out structural risk confirmed.** `F-7` supplies 8 of 26 held-out units (31%) and 7 of those are `R-6`. Every other family supplies exactly 2. `R-1` has **zero** held-out coverage. `S-4`'s per-family rule is met; per-mode balance is not required by any rule. `[CALC]`

---

## 4. `R-6` — the contradiction, stated and left undecided

The assignment forbids a premature `R-6` decision. This section therefore records the conflict and decides nothing.

| Source | Text | Implication |
|---|---|---|
| `MEMORY-09` L119 | "A configuration that suppresses everything passes **four of the six** and fails **the two** that require the memory to be perceptible (`R-2`, `R-3`)" | `R-6` is silence-passable; only two modes require perceptibility |
| `CORPUS-01` L126 | "Uncertainty is conveyed" · "Can silence pass it? **No** — saying nothing also fails" | `R-6` requires a perceptible response |
| `CORPUS-01` L128 (`A-1`) | "passes `R-1`, `R-4`, `R-5`, `R-6` while failing `R-2` and `R-3` has learned **avoidance**" | `R-6` is silence-passable |
| `check_corpus.py` L58 | `POSITIVE_MODES = {"R-2", "R-3", "R-6"}` | `R-6` requires a perceptible response |

All four verified at source this session. `[OBS]`/`[GT]`

**F-02 is CONFIRMED as a contradiction, not a terminology difference.** The two readings give different pass conditions in the limiting case (a reply that says nothing), and "passes four / fails two" is true only under the `MEMORY-09` reading. `[IA]`

**Neither option is recommended here.** Option A (R-6 perceptible) requires a dated amendment to `MEMORY-09` L119 *and* to `CORPUS-01` `A-1`, because both are approved text. Option B (silence passes R-6) requires amending `CORPUS-01` L126 and the checker, and would make `R6-03a` — the anti-hedging control whose failure mode *is* refusing to know — incoherent. Both are coherent positions; neither is mine to take. `[UNRES]`

**One consequence that no report states and that follows from the arithmetic:** whichever option is chosen, `MEMORY-09` L119's phrase "the two" is the only sentence in the approved pre-registration that enumerates the perceptibility requirement at all. Option A therefore cannot be implemented by editing `CORPUS-01` and the checker alone. `[RI]`

**F-17 confirmed and distinct.** `CORPUS-07` L181 records owner decision `OD-2`: "**Score the behaviour** that preserves uncertainty. Epistemic narration is **not** required." `CORPUS-07` L211 records it as "**Recorded, not applied**." `[GT]`/`[OBS]` Meanwhile `R6-01`'s `mandatory` requires that "the choice is **presented** as made under conflicting information" (`CORPUS-04` L26 `[GT]`). So R-6 is currently governed by two constructs at once: an owner-decided behavioural one, unapplied, and a narration-permitting one, applied but owner-unauthorised. `[IA]`

---

## 5. Verification of the `CORPUS-91`, `CORPUS-92` and registry integrity findings

### 5.1 `CORPUS-91` was never filed — F-23 CONFIRMED

| Check | Result |
|---|---|
| Lines / bytes / SHA-256 | 627 / 56,012 / `fda9b5e12f49fa308f24d38641c866dbd759e917cc196ee2b309753c694d542a` `[OBS]` |
| §9.2 registry row | Placeholder: "`SHA-256:` <to be computed after file save>", same for Bytes and Lines (L612–614) `[OBS]` |
| §9.3 | "Let me now save the file and compute the verification metrics." (L623) `[OBS]` |
| End of file | **Unexecuted tool-call fragment**, L625–627: `<tool_call>` / `<function=run_terminal_command>` / `<parameter=command>` then EOF `[OBS]` |
| Registry row 424 | Exists and matches the file (hash, bytes, lines) `[OBS]` |

**Verdict:** a complete analytical body with an **incomplete filing section**. The registry row is correct because someone completed the filing outside the document; the document's own §9 certifies a verification it never performed. Muse's classification — retain the body, refile §9, no analytical loss — is correct. `[IA]`

### 5.2 `CORPUS-92`'s self-declared metrics are stale — F-22 CONFIRMED

| Source | SHA-256 | Bytes | Lines |
|---|---|---|---|
| `CORPUS-92` §9.2 (self-declared) | `f61da454bb9184a71cce1822f9e7ae8789faf3da185d9bd06efb4ecaa9a94a09` | 33,793 | 469 |
| Actual file (this session) | `c76db062b6abe31fe02d650e69c55441d7eb64cb9bbf655d78fe52069ef98299` | 33,724 | 469 |
| Registry row 425 | `c76db062…98299` | 33,724 | 469 |

`[OBS]`. **The registry is correct; the report body is wrong by 69 bytes and a different hash.** Line count matches, so the later edit was within lines — consistent with §11's classification summary being revised after §9.2 was hashed. The defect the assignment asked this lineage to catch **persists in the document written to catch it**. `[IA]`

### 5.3 Registry integrity — F-24 CONFIRMED, figures now moved

| Check | Muse at filing (§9.7) | This session |
|---|---|---|
| Rows matching `^\| N \|` | 427 | **428** `[OBS]` |
| Files under `docs/` | 428 | **429** `[OBS]` |
| Duplicate row numbers | 423, 424 | **423, 424** `[OBS]` — both persist |
| Header claim (L6) | "Rows: **422**" | unchanged `[OBS]` |
| §4 claim (L489) | "Rows: **409**" | unchanged `[OBS]` |
| Governor of row numbers | none | none — `A-REG-01` unrepaired `[GT]` |

The +1 on both counts is Muse's own row 426, correctly filed: row 426's hash and byte count **do** match the file on disk (778 lines / 87,623 B / `95a34acd…`). `[OBS]` So the report that found the filing-mismatch defect filed correctly, and the report that missed it (F-23's) did not.

**Consequence unchanged:** every file is registered, but row identity is not unique. Row 424 refers to both `CORPUS-90` and `CORPUS-91`; row 423 to both `CORPUS-82` and `CORPUS-83`. `CORPUS-92`'s "No duplicate row was found" is false, and its "Row 424" citation is true of the tail row and false of the earlier one. `[IA]`

### 5.4 Citation integrity — F-14 CONFIRMED

`MEMORY-09` is **344 lines** (`wc -l`, verified). `[OBS]` Both commissioned adjudications cite ranges past its end:

- `CORPUS-82` L36: "L113–122, L220–289, L292–353, **L355–411, L414–472, L475–526**" — three of six ranges exceed 344. `[OBS]`
- `CORPUS-82` L93, L113, L121, L149, L251 repeat `L292–353` / `L355–411` / `L475–526` as support for R-3/R-10/R-13 intent. `[OBS]`
- `CORPUS-83` L33 and L155 repeat the same ranges. `[OBS]`

**The invented "Reading A"** is also confirmed as propagating: `CORPUS-82` §3.1 and `CORPUS-83` §10.1 describe a table listing `F2-01` under `R-1` and `F3-01` under `R-2`; the actual `R-1` row is `P-03, F9-01` and `F3-01` appears in no row at all. `[OBS]` `CORPUS-91` §5.5 inherited it and `CORPUS-92` §3 correctly refuted it. **Muse is right that the correction must reach `CORPUS-82`, not only `CORPUS-91`.** `[IA]`

What survives: the *pattern* — incomplete rows, `R-1` misattributed — is real and independently verified (§3.3). Both reports reached a true conclusion through a false description of their source. `[RI]`

---

## 6. Finding-by-finding verification of all 30 register entries

Verdicts: **Confirmed** / **Confirmed with correction** / **Partly confirmed** / **Not verified here** / **Not supported**.

| # | Muse finding | Verdict | Independent basis |
|---|---|---|---|
| F-01 | Membership rows/coverage/status disagree | **Confirmed** | §3.3: 16 entries/14 distinct; `R-1` both entries wrong; 19 units unnamed; coverage sums 42; `R-6` 6 vs 9 `[OBS]` |
| F-02 | R-6 contract self-contradictory | **Confirmed** | §4: all four sources read at source `[GT]` |
| F-03 | `CORPUS-01` stale internal state | **Confirmed, and understated** | L166 stale vs `TOTALS MATCH`; L180 "42 of 80" wrong base (41); over-exposure "30/45, 15 missing" vs checker 0 missing; §4 heading "batches 1–2" over 45 units `[OBS]` |
| F-04 | No source conversations | **Confirmed (strongest class)** | §2: 0 occurrences of the schema's own field names; prose evidence; no ledger `[OBS]` |
| F-05 | `E-01`–`E-12` unconverted | **Confirmed** | 12 experiment specs in `MEMORY-04` L456–467; **0** `E-` units in the item files `[OBS]` |
| F-06 | 41 family units vs 80; 39 missing | **Confirmed with correction** | 41/80 confirmed; aggregate gap 39, **per-family deficits 42**; "author 39" is the same `F-7` offset Muse rejects in `CORPUS-01` `[CALC]` |
| F-07 | Held-out unsealed, interleaved | **Confirmed** | No manifest file exists anywhere in the tree; `CORPUS-01` §7 "Held-out sealed: **no**"; tuning/heldout share three files `[OBS]` |
| F-08 | `P-3` not started; no independent reviewer | **Confirmed** | `CORPUS-05` L5 "protocol defined; review **not executed**"; `CORPUS-01` §7 "not started"; convergence 0 of 2 at `CORPUS-80` L154/L201, `CORPUS-89` L86/L128 `[OBS]` |
| F-09 | Membership authority undecided | **Confirmed as absence** | No document states which representation governs; verified by search `[OBS]` |
| F-10 | Floor accounting (do `P` count; `F-7` surplus) | **Confirmed as undecided** | `MEMORY-09` L104: "80 scored items, plus the twelve fixtures"; `P-01`–`P-04` nowhere included in the 80 — the strict reading is textually supported, and the surplus question is genuinely open `[GT]` |
| F-11 | Verdict states undefined | **Confirmed (decisive)** | `MEMORY-09` contains **0** occurrences of "assessable"/"demonstrated" — the vocabulary of four reports is absent from the governing text `[OBS]` |
| F-12 | Proposal status drift; `CORPUS-83` self-adopted | **Confirmed** | `CORPUS-83` §6 "**Adopted.**" vs §4 "held unresolved"; `CORPUS-89` §3 records no adoption `[GT]` |
| F-13 | "Application-only" rests on item text | **Confirmed** | `CORPUS-83` §7 cites `F6-02`'s "song not pushed", `F8-03`'s "Quiet day", `P-03`'s secret — none is an application record `[IA]` |
| F-14 | Fabricated citations; invented Reading A | **Confirmed** | §5.4: 3 of 6 ranges exceed 344 lines in both reports; Reading A absent from `CORPUS-01` `[OBS]` |
| F-15 | 20 of 33 restraint units are over-injection detectors | **Partly confirmed** | **16** units (`R-1` 2 + `R-4` 4 + `R-5` 10) are absence-scored *by their pass conditions* — verified from `CORPUS-01` L121/L124/L125. Whether the `R-2` units bring it to 20 requires item-level judgement I did not fully re-derive `[OBS]`/`[RI]` |
| F-16 | "1 confirmed diagnostic of 24" | **Not verified here** | Muse also marks it Unverified; I did not re-derive `CORPUS-57`'s classifications `[UNEST]` |
| F-17 | `OD-2` unapplied vs R-6 `mandatory` narration | **Confirmed** | `CORPUS-07` L181 (OD-2), L211 (Recorded, not applied), `CORPUS-04` L26 `[OBS]` |
| F-18 | Checker narrower than schema | **Confirmed as a gap; correction to the "overclaim"** | The checker omits `family`, `source_conversation_id`, `turn_range`, `expected_abstention`, perspective, and the `*_ids[]` form `[OBS]`. But `CORPUS-01` §6.2 *enumerates* the fields it checks, so the wording is a subset description, not a false total. The real defect is the **schema/enforcement gap**, not the sentence `[RI]` |
| F-19 | Threshold regime conflict (`MEMORY-03` §8.2 vs `MEMORY-09` §10) | **Not verified here** | I did not read `MEMORY-03` §8.2 `[UNEST]` |
| F-20 | Abstention "8 total" wrong | **Confirmed** | 4 (batch 1) + 3 (batch 2) = **7**; `CORPUS-03` L290 says 8 `[OBS]` |
| F-21 | `CORPUS-81` §5.2 split column wrong for `R-4`, `none` | **True values confirmed; text not reopened** | Actual `R-4` 3/1 and `none` 5/7 `[CALC]`; I did not re-read `CORPUS-81` §5.2 to confirm it states 2/2 and 6/6 `[UNEST]` |
| F-22 | `CORPUS-92` metrics stale | **Confirmed** | §5.2: 69-byte and hash mismatch; registry correct `[OBS]` |
| F-23 | `CORPUS-91` never filed | **Confirmed** | §5.1: placeholders L612–614; EOF tool-call fragment L625–627 `[OBS]` |
| F-24 | Registry duplicated rows; header 422 vs §4 409 vs 427 | **Confirmed** | §5.3: duplicates 423 and 424 persist; header L6 422; §4 L489 409; actual 428 `[OBS]` |
| F-25 | No run manifest / token instrument | **Confirmed as absence** | No manifest artifact found by search; `R-10`/`R-4` rules exist without instruments `[OBS]` |
| F-26 | Displayed independence unreliable | **Confirmed** | `CORPUS-91` §0.1 "Nemotron Lightning"; `CORPUS-79` §0 "Buffy"; `CORPUS-80` L5 "DeepSeek/Buffy"; `CORPUS-83` §0 "Glyph Cluster" `[GT]` |
| F-27 | `CORPUS-58` substantive review absent | **Not verified here** | I did not open registry row 393 or `CORPUS-62`/`63` `[UNEST]` |
| F-28 | Table gives `P-03`, `F9-01` two modes, contra `A-2` | **Confirmed** | `CORPUS-01` L121/L124/L125 `[OBS]` |
| F-29 | `CORPUS-01` header misstates pause reason | **Confirmed** | Header says pause "pending independent review"; the review performed was architecture adjudication; `P-3` is "not started" in the same file's §7 `[OBS]` |
| F-30 | Already-resolved items to retain | **Confirmed** | checker PASS; 45/45 four-part answers; single tag per unit; 19/26 reproducible; ≥2 held-out per family `[OBS]` |

**Tally:** 24 confirmed outright; 1 confirmed with a number correction (F-06); 1 confirmed-and-understated (F-03); 2 partly confirmed (F-15, F-18); 1 true-values-confirmed-text-not-reopened (F-21); 3 not verified by this review (F-16, F-19, F-27). **No Muse finding is contradicted.** The register's substantive reliability is high; its two arithmetic/precision defects are F-06's headline and F-03's unremarked `42 of 80`.

---

## 7. What is reusable, and what genuinely needs rebuilding

### 7.1 Reusable as-is

| Asset | Why it survives | Basis |
|---|---|---|
| The 45 item specifications | Internally coherent, four-part answers, forbidden sets, over-exposure boundaries, failure classes. They are the **authoring briefs** for the conversations. | F-30; checker PASS `[OBS]` |
| `CORPUS-01` §2 item schema | Correct shape; the gap is field *population*, not design | `[GT]` |
| `CORPUS-01` §2.4 four-part reference standard | Real, applied, and the reason style cannot be scored as fact | `[GT]` |
| Rules `M-1`/`M-2` over-exposure boundary | The single best safeguard in the corpus; makes verbosity scoreable as failure | `[GT]` |
| `check_corpus.py` | Reports only, never repairs; states its own three unverifiable classes honestly; reproduces 45/19/26 | §3.2 `[OBS]` |
| `CORPUS-04`'s four uncertainty mechanisms + `R6-03a`–`d` provenance variants | Genuine construct design, not repetition of one behaviour | `[GT]` |
| `P-01`–`P-04` principle cases | Mandatory, and the only items that cannot be passed by injecting more | `[GT]` |
| The `R-3` units | Only mode with a positive, memory-attributable predicate in the item text | §6, F-15 `[RI]` |
| Split discipline as a number (19/26; ≥2/family) | Reproducible from files | `[OBS]` |
| `CORPUS-80` §3 P1–P8; `CORPUS-81` C-01–C-08 + tag inventory; `CORPUS-82` §4A–§4E *reasoning*; `CORPUS-90` §1.4 + §5.5; `CORPUS-91` body except §5.5; `CORPUS-92` audit except §9.2 | Each independently re-derived here or consistent with my derivation | §6 `[IA]` |

### 7.2 Reusable after repair (repair, not rebuild)

| Asset | Repair needed |
|---|---|
| `CORPUS-01` §3.1 rows, L134 coverage, L166, L180, §4 heading, §7 status | Replace rows with the tag inventory or correct them; fix `R-6` 6→9; delete stale discrepancy; fix base 42→41; one dated editorial pass **after** WP1 |
| `CORPUS-82` / `CORPUS-83` citations and Reading A | Re-derive ranges to ≤344 lines; delete invented Reading A; carry `CORPUS-92`'s correction back to `CORPUS-82` |
| `CORPUS-83` §6–§7 | Restore "held unresolved"; withdraw "Adopted."; withdraw the unsupported classification (§6.1) |
| `CORPUS-91` §9 | Complete the filing section; remove the EOF fragment |
| `CORPUS-92` §9.2 | Correct to registry values (`c76db062…`, 33,724) |
| `CORPUS-81` §5.2 | Recompute `R-4` (3/1) and `none` (5/7) |
| `CORPUS-03` L290 | 8 → 7 |
| Registry | Renumber the two colliding pairs or adopt a uniqueness rule; reconcile header 422 / §4 409 / actual 428 |
| `check_corpus.py` | Extend to `family`, `source_conversation_id`, `turn_range`, `expected_abstention`, `*_ids[]` — **after** the conversations exist |
| `R6-01`, `R6-03b`, `R6-03d`, `R6-04` | Apply `OD-2` (owner-authorised) and re-review; blocked on F-02 |

### 7.3 Genuinely new work (not rebuild — it never existed)

1. **42 source conversations** (one per scene heading; `R6-03` yields four scored units).
2. **The `MEMORY-03` §5.3 ground-truth ledger** — fact IDs, evidence spans, visibility.
3. **12 converted `E-01`–`E-12` scored units.**
4. **42 new family units** under per-family minimums (39 + `F-7` reallocation authority), or 39 if OD-C authorises reallocation.
5. **Physically separated `tuning`/`heldout` files + published hash manifest.**
6. **Run manifest schema, token-accounting instrument, constant-budget definition, variance plan.**
7. **A `P-3` review of 45+ units** by an eligible reviewer.

### 7.4 Contingent rebuild — depends on an undecided owner question

If the owner requires `R-1`/`R-2`/`R-4`/`R-5` to be able to **demonstrate** memory-guided restraint rather than merely not-fail, those **16–20 units** need new authoring with paired conditions or declinable cues. That is a design decision (OD-E), not a review conclusion, and I do not make it. If the owner instead accepts them explicitly as over-injection detectors, they are retained unchanged and the corpus must say so in its own text. `[RI]`/`[PROP]`

---

## 8. Dependency-ordered remediation proposal

**Ordering principle.** Nothing that needs a settled construct is built before the construct is settled; nothing that needs source material is reviewed before the material exists; nothing is sealed before it is final; no owner decision is requested twice; and the record must be citable before it is cited.

```text
WP0 record-integrity repair (no dependencies, do first)
      │
WP1 owner decisions: R-6 contract, membership authority, floor accounting
      │                    ╲
WP2 construct & verdict definitions (parallel with WP1)  ╲
      │                                                   ╲
WP3 source conversations + ground-truth ledger  ◄──────────╲
      │
WP4 item repair, coverage, E-01–E-12, OD-2 application
      │
WP5 protocol, manifest, token instrument, hash seal
      │
WP6 model-free validation
      │
WP7 P-3 independent review + reconciliation
      │
WP8 G1 completion instrument + registry repair
```

### WP0 — Record-integrity repair (independent; no owner decision required)

- **Do:** complete `CORPUS-91` §9 and delete the EOF fragment; correct `CORPUS-92` §9.2 to registry values; correct `CORPUS-03` L290; recompute `CORPUS-81` §5.2; renumber registry duplicates 423/424 or adopt a uniqueness rule and reconcile header/§4 counts.
- **Why first:** every downstream report cites these rows and hashes. `CORPUS-92` and `CORPUS-91` are currently the *evidence* that filing discipline is unreliable, which makes them unusable as evidence.
- **Does not touch** items, tags, splits, `MEMORY-09`, or `CORPUS-01`.
- **Completion:** no citation exceeds source length; row numbers unique; every report's self-declared metrics match its bytes.

### WP1 — Owner decisions (blocks everything substantive)

| Decision | Question | Why it blocks |
|---|---|---|
| OD-A | Does `R-6` require a perceptible response? | F-02; no R-6 unit is scoreable against a settled rule until answered |
| OD-B | Do item tags or §3.1 rows govern membership? | F-01/F-09/F-28; no count or denominator is citable until answered |
| OD-C | Do `P-01`–`P-04` count toward 80? Is `F-7`'s 3-unit surplus reallocable? | F-10; determines whether 39 or 42 units must be authored |
| OD-I | Is authoring 42 conversations + ledger + 12 fixture conversions inside `G1`? | F-04/F-05; the largest body of work, unassigned |

**This report takes no position on OD-A.** `[UNRES]`

### WP2 — Construct and verdict definitions (parallel with WP1)

Adopt one definition of `Demonstrated` / `Not demonstrated` / `Not assessable` (OD-D), with assessability decided **pre-response**, and settle the denominator rule. `MEMORY-09` supplies none of this vocabulary (§6, F-11). Recommended source: `CORPUS-82` §4D — but it is a **proposal**, and adopting it is an owner act. `[PROP]`

### WP3 — Source conversations and ledger (the long pole)

42 conversations; `MEMORY-03` §5.3 ledger with fact IDs, spans, visibility, provenance; bind every `required`/`forbidden` reference to a span and an ID. **Can begin immediately on `R-3` and `R-6` units**, whose constructs are settled and whose items already discriminate. It should **not** begin on `R-1`/`R-2`/`R-4`/`R-5` until OD-E decides whether those items are to be rebuilt. `[RI]`

### WP4 — Item repair, coverage, fixtures

Apply `OD-2` to the four R-6 items; author 42 family units (or 39 under OD-C authorisation); convert `E-01`–`E-12`; classify each unit's discriminability **pre-response**; repair `CORPUS-01`'s figures once, after WP1. Depends on WP1 (construct clarity), WP2, WP3.

### WP5 — Protocol and evidence pre-registration

Manifest schema (`R-10`), token accounting, `C1`/`C2`/`C3` constant budget, paired-condition design if adopted, **physically separated split files**, published hash manifest, variance plan. Depends on WP3, WP2.

### WP6 — Pre-inference validation

Scope structure, deletion/resurrection, provenance completeness, token accounting, retrieval recall/nDCG over the ledger, schema validity. All model-free, all requiring WP3's addressable corpus. Failures preserved (`R-14`).

### WP7 — `P-3` independent review

Every unit reviewed blind to author intent; `agree`/`disagree`/`ambiguous`; reconciliation register; retirements marked `RETIRED`. Requires a **genuinely separate reviewer** (OD-H) — convergence is 0 of 2 and no artifact in this tree can supply independent confirmation of another. Depends on WP3, WP4.

### WP8 — `G1` completion and registry repair

Owner instrument citing objective, file-checkable criteria; `A-REG-01` closed or formally accepted; no run, score or result cited.

### 8.1 What must not be done

**Do not commission another adjudication before the conversations exist.** This is Muse's sharpest point and I confirm it from the evidence: twelve consecutive artifacts have argued about how to score material that has never been written, and each one has concluded the question cannot be settled — for three different reasons, none of which is the blocker. A further adjudication moves the verdict nowhere; it recedes by one artifact. `[RI]`

### 8.2 Where I differ from Muse's sequencing

| Muse | This report | Reason |
|---|---|---|
| WP1 includes editorial corrections to `CORPUS-01` | **Split:** WP0 administrative repair; `CORPUS-01` figures in WP4 | Registry/filing/hash repair has **no** dependency and improves citability immediately; `CORPUS-01`'s figures depend on OD-A/OD-B/OD-C and would be wrong again if fixed first. Muse itself notes this tension in §11.3 but keeps both in WP1 |
| WP2 conversations may begin on `F-7`/`R-6` and `R-3` | **Agreed and adopted** | Constructs settled; items already discriminate |
| F-06: "author 39 units" | **42, or 39 under an explicit reallocation authorisation** | Per-family deficits sum to 42; `F-7`'s surplus cannot absorb a per-family minimum (§3.4) |

---

## 9. Owner decisions required — none taken here

Restated from the verification, not created by it. `[PROP]`

| # | Decision | Blocks |
|---|---|---|
| OD-A | Does `R-6` require a perceptible response? (If yes, `MEMORY-09` L119 **and** `CORPUS-01` `A-1` need dated amendment.) | F-02, WP1, WP4 |
| OD-B | Do item tags or §3.1 rows govern membership? | F-01, F-09, F-28 |
| OD-C | Do `P-01`–`P-04` count toward 80? Is `F-7`'s surplus reallocable? | F-06, F-10, WP4 |
| OD-D | Which verdict-state definition is adopted, and by whom? | F-11, WP2 |
| OD-E | Are `R-1`/`R-2`/`R-4`/`R-5` accepted as over-injection detectors, or rebuilt to demonstrate? | F-13, F-15, §7.4 |
| OD-F | Authorise applying `OD-2` to the R-6 items + re-review | F-17 |
| OD-G | One threshold regime (`MEMORY-03` §8.2 vs `MEMORY-09` §10) | F-19 *(not verified here)* |
| OD-H | Reviewer scheduling, and a designation→executing-identity map | F-08, F-26 |
| OD-I | Authorise the source corpus (42 conversations + ledger + 12 conversions) | F-04, F-05 |
| OD-J | Apply, modify or decline the pending amendment proposals | F-01, F-02 |

---

## 10. Corrections to Muse's report

| Where | Correction |
|---|---|
| F-06, "Author 39 units" | Aggregate gap is 39; **per-family deficits sum to 42**. 39 units reaches total 80 but leaves three families short unless OD-C authorises reallocating `F-7`'s surplus. WP4's own completion criterion ("every family at minimum") is the 42-figure. |
| F-03 | Understated. `CORPUS-01` L180's "42 of 80" is wrong on the base too — the corpus holds **41** family units, not 42; the 42 appears to be the heading count. |
| F-18 | The `CORPUS-01` §6.2 "overclaim" is arguable: the sentence *enumerates* the fields checked. The defensible finding is the **schema/enforcement gap** (6 required fields and the `*_ids[]` form unchecked), not a false total. |
| F-01 / §5.5 | "Coverage line sums to 42" is right, but the cleaner statement is: restraint terms sum to 30 and `none` adds 12, so 42 total against 45 actual; the sole wrong term is `R-6` 6 vs 9. Also adopt `CORPUS-92` §4.4's "all six rows incomplete, `R-1` misattributed" over "five of six". |
| F-04 | Effort should be stated as **42 scene conversations** (45 scored units, `R6-03` yields four), not 45 conversations. |
| §9.7 registry counts | Correct as of filing (427/428); now 428/429 after row 426. Not an error — a time-dependence worth recording. |

**Nothing in Muse's §10 register is contradicted by this verification.** These are precision corrections, not reversals.

---

## 11. What this review did not verify

1. **F-16** — the "1 confirmed diagnostic of 24" figure; I did not re-derive `CORPUS-57`'s classifications.
2. **F-19** — the `MEMORY-03` §8.2 / `MEMORY-09` §10 threshold conflict; I did not read §8.2.
3. **F-27** — the `CORPUS-58` review absence; I did not open registry row 393 or `CORPUS-62`/`63`.
4. **F-21** — I verified the true split values (`R-4` 3/1, `none` 5/7) but did not re-read `CORPUS-81` §5.2 to confirm it states the wrong ones.
5. **F-15** — I confirmed **16** absence-scored units from their pass conditions; the 20-figure requires item-level reading of the `R-2` units that I did not fully perform.
6. **F-16/F-15 downstream** — no discriminator ratio is computed here; computing it would be a re-classification exercise, not a verification.
7. `CORPUS-78`, `CORPUS-79`, `CORPUS-81`, `CORPUS-82`, `CORPUS-83` were not read in full.
8. **No model identity** was available for any artifact, including this one.

---

## 12. Filing note

| Item | Value |
|---|---|
| Path | `docs/research/corpus/CORPUS-93-INDEPENDENT-VERIFICATION-OF-MUSE-G1-FINDINGS-AND-REMEDIATION-PROPOSAL-2026-11-13.md` |
| Family | Research |
| Role | Independent verification and remediation proposal |
| Current relevance | CURRENT — advisory; adopts nothing; authorises nothing; decides nothing |
| Documents modified | **none** |
| Corpus records modified | **none** |
| Registry modified | **no** |

**The registry was deliberately not modified.** This assignment forbids changing governing documents, and `DOCUMENT-REGISTRY.md` is one. Registering this file is a secretariat action that should happen **after** WP0 repairs the duplicate row numbers 423 and 424 — otherwise a new row is appended to a table whose unique-integer assumption is already broken, and the row will inherit the defect it was added beside. The identity anchor for this artifact is therefore its path plus its bytes on disk, and no self-referential hash is asserted inside it — the same convention `CORPUS-80` §11.1, `CORPUS-90` §11 and `MEMORY-09-G1-READINESS-01` §14.3 adopted, for the same reason: `CORPUS-92` §9.2 is the worked example of what happens when a document states its own hash before its last edit. `[RI]`

---

## 13. Conclusion

**Muse's assessment is substantially correct and its diagnosis is reliable.** Of 30 register entries, 24 are confirmed outright from the source artifacts, one is confirmed and understated, two are partly confirmed, one is confirmed in its numbers, and three I could not verify. **None is contradicted.** The three integrity findings that this assignment asked to check are all true at the byte level: `CORPUS-91` ends mid-tool-call with a placeholder filing section, `CORPUS-92` declares a hash and byte count that do not match its own saved bytes while the registry row is correct, and the registry carries two duplicated row numbers against a header claiming 422 rows and a §4 claiming 409. `[OBS]`

**The corpus inventory is in better shape than the record says and the source evidence is in worse shape than the work queue implies.** The 45 units, the 33/12 restraint split, the 19/26 split, `R-6`'s nine units and the ≥2-held-out-per-family rule all reproduce exactly from the files, and the checker passes honestly. But the schema's own field names — `source_conversation_id`, `turn_range`, `required_evidence_ids[]` — appear **nowhere** in the corpus, and no conversation or ledger exists. That is not a filing gap; it is the missing body of work. `[IA]`

**The reuse picture is favourable.** Almost nothing needs rebuilding. The 45 specifications become authoring briefs; the schema, the four-part reference standard, the over-exposure rules, the checker, the `R-6` provenance design, the principle cases and the `R-3` items all survive. The rework list is short and mechanical: citations in `CORPUS-82`/`83`, Reading A in `82`/`83`/`91`, status in `CORPUS-83` §6–§7, metrics in `CORPUS-92`, §9 in `CORPUS-91`, one split column in `CORPUS-81`, one abstention count in `CORPUS-03`, and one editorial pass over `CORPUS-01` **after** the owner rules. The genuinely new work is the conversations, the ledger, the 12 fixture conversions, 42 family units, the sealed split and the `P-3` review. `[RI]`

**Two owner decisions unlock the most and neither is a reviewer's to make:** whether `R-6` requires a perceptible response, and whether item tags or the specification table govern membership. Until the first is answered no R-6 unit is scoreable against a settled rule; until the second is answered no count, denominator or coverage claim is citable. The remaining sequencing is in §8 and requires no further adjudication to begin. `[PROP]`

**End of report.**
