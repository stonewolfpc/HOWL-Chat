# CORPUS-77 — Independent Adversarial Review: OD-9 and OD-10 Scoring Architecture (CORPUS-75/76)

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Identifier:** `CORPUS-77` · **Revision:** 1 · **Date:** 2026-11-12.
**Review type:** Independent, read-only, source-anchored adversarial review.
**Trigger:** Work order "Independent adversarial review — CORPUS-75 / CORPUS-76 — OD-9 and OD-10 scoring architecture".
**Reviewer identity:** Muse Spark, as identified by the execution environment. No assignment-supplied model identity was provided in the work order; no underlying-model claim is made or inferred. Persona label handling follows the `CORPUS-74` §1 disclosure pattern.
**Authorization:** Review and report only. No implementation or amendment authority. No standards, items, tags, thresholds, or registry records (outside this report's own row) were changed; no inference was run.
**Status:** Filed as a review. Nothing here amends, repairs, re-tags, rescores, or authorizes anything.

---

## 1. Execution context and source provenance

This review was conducted in a single session against the working tree at `C:\Users\StoneWolfPC\HOWL_Chat` (project root). All governing and corpus sources were read from disk in full or in the exact sections cited; no conversation summary was used as a substitute for source inspection. Direct quotations below are verbatim from the cited files; interpretation is labelled as such.

### 1.1 Exact sources inspected

| # | File | Scope read | Identity anchor (working tree) |
|---|---|---|---|
| 1 | `docs/research/MEMORY-09-MEMORY-EVALUATION-PREREGISTRATION-2026-11-12.md` | Full | Rev 2, approved governing pre-registration authority |
| 2 | `docs/research/corpus/CORPUS-01-SPECIFICATION-AND-INDEX-2026-11-12.md` | Full | Spec + index, `be1f11b4…` lineage (hash equality per `CORPUS-74` §10, not re-derived here beyond cited text) |
| 3 | `docs/research/corpus/CORPUS-02-BATCH-1-ITEMS-2026-11-12.md` | `F1-01`, `F1-02`, `F3-01`, `F6-01`, `F6-02`, `P-02`, `P-03`, `P-04` in full | Batch-1 items |
| 4 | `docs/research/corpus/CORPUS-03-BATCH-2-ITEMS-2026-11-12.md` | `F1-03`, `F6-03`, `F6-04` in full; rest skimmed | Batch-2 items |
| 5 | `docs/research/corpus/CORPUS-04-R6-UNCERTAINTY-ITEMS-2026-11-12.md` | `R6-01`–`R6-04` in full | Batch-3 items |
| 6 | `docs/research/corpus/CORPUS-07-ADJUDICATION-REGISTER-2026-11-12.md` | §§0–2, 6, 8, 10 in full | Adjudication register |
| 7 | `docs/research/corpus/CORPUS-75-CORPUS-74-RECONCILIATION-AND-R2-STATUS-2026-11-12.md` | Full | Reconciliation under review |
| 8 | `docs/research/corpus/CORPUS-76-OWNER-DECISION-BRIEF-R2-SCORING-2026-11-12.md` | Full | Decision brief under review |
| 9 | `docs/research/corpus/CORPUS-74-INDEPENDENT-ADVERSARIAL-REVIEW-CORPUS-69-73-2026-11-12.md` | §§4–11 in full | Reviewed review |
| 10 | `docs/research/corpus/CORPUS-72-OWNER-DECISION-BRIEF-OD1-OD8-2026-11-12.md` | Full | OD-1–OD-8 baseline |
| 11 | `docs/research/corpus/CORPUS-65-CORPUS-01-R2-COVERAGE-AMENDMENT-PROPOSAL-2026-11-12.md` | §§1–9 in full | `A-R2-01` proposal |
| 12 | `docs/research/corpus/CORPUS-71-PROPOSAL-REVISION-REGISTER-V1-1-2026-11-12.md` | §6 (`A-R2-02`) in full | `A-R2-02` proposal |
| 13 | `docs/research/corpus/CORPUS-70-…-RECONCILIATION-…md` | §§5.3, 6.1–6.4 in full | `U-3`, `A-R2-02` provenance |
| 14 | `docs/research/corpus/check_corpus.py` | Full | Executed: `PASS (0 fail, 0 warn, 0 unverifiable)`, 45 units, 19/26, `TOTALS MATCH` |
| 15 | `docs/DOCUMENT-REGISTRY.md` | Header, §§4–6, rows 394–411 | Registration procedure |

Not re-read in full (relied upon only through the above): `CORPUS-62/63/64/66/69/73`. `CORPUS-64`/`66` substance is taken through `CORPUS-70` §§4–5 and `CORPUS-71` §§4–6 as secondary evidence, flagged where it matters (§8.3).

### 1.2 Mechanical verification performed

- `python docs/research/corpus/check_corpus.py` → exit 0, `scored units: 45`, `split: tuning=19 heldout=26`, `TOTALS MATCH`, `RESULT: PASS (0 fail, 0 warn, 0 unverifiable)`. Reproduced this session.
- Project-wide literal search for `half moon`, `full moon`, `ambiguous trigger`, `almond milk`, and persona-name `Alex` (case-insensitive): the moon strings occur **only** in `CORPUS-74` (and `CORPUS-75` quoting it); `almond milk` occurs **only** in `CORPUS-74` (and `CORPUS-75` quoting it); the only `Alex` in the corpus workspace is the unrelated author name in `MEMORY-02` L327. This corroborates `CORPUS-75` §4 findings 1–2.
- Current working-tree counts: 412 files under `docs/`, 383 markdown files (recomputed this session). This matches `CORPUS-75`'s snapshot-C description (412/383, after its own filing) and is consistent with two further registrations since (`CORPUS-75`, `CORPUS-76` = rows 410–411; header now reads 411 rows).

---

## 2. Governing baseline — what the sources actually say

### 2.1 The interaction question is a genuine specification gap

`CORPUS-01` §2.4 scores four reference-answer parts: `mandatory` (**Yes**), `variation` (all pass), `never scored` (**never a failure signal**), `prohibited` (**Yes**). `M-1`/`M-2` add a scored over-exposure boundary. §3.1 defines six restraint modes with per-mode pass conditions and a "Can silence pass it?" column (`R-1` Yes; `R-2` **No**; `R-3` **No**; `R-4` Yes; `R-5` Yes; `R-6` **No**). Rule `A-2` requires every restraint item to carry exactly one mode. Rule `A-3` states modes are independent of quality dimensions — *"Each result stands alone; no mode excuses or compensates for another, and the same applies across the six scoring dimensions. The report is a matrix, not a score."* Rule `G-1` states no aggregate offsets a gate.

**What is absent:** no sentence in `CORPUS-01` states how a mode result interacts with the ordinary acceptance fields (`mandatory`/`variation`/`prohibited`/`over-exposure`). `A-3` contemplates divergence *between modes* and *between modes and quality dimensions* — it does not state whether a field-pass/mode-fail pair is an item pass, an item fail, or two independent facts. `CORPUS-74` §4 answer 7 and `CORPUS-75` §3 correctly identify this as unresolved. I independently confirm the gap: it is real, not an artefact of selective quotation.

### 2.2 What the preregistration does and does not settle

- `MEMORY-09` §3.5: *"Six behaviours are scored separately and never aggregated"*; `R-2`/`R-3` *"require the memory to be perceptible"*; suppression-only success is *"reported as **avoidance**, not as restraint."* **Interpretation (mine):** this is a *reporting* architecture (separate axes, no aggregation, no suppression-credit) plus a *construct* requirement (perceptibility for `R-2`/`R-3`). It does **not** state an acceptance-logic connective (AND vs independent). Separate reporting is compatible with either gating or non-gating at the item-verdict level — the inference "separate reporting ⇒ non-gating" is tempting but not compelled, and "separate reporting ⇒ gating" is likewise not compelled.
- `MEMORY-09` §5 aggregation rule (*"No weighted total is published as a headline"*) and §7 no-trading rule (a recall/leakage win never overrides a gate) govern *score aggregation and gate precedence*. Neither addresses field-vs-mode interaction. The gates in §7 (Integrity, Formation, Retrieval, Presentation, Creative acceptance, Operations) and `G-1` blocking-gate items (scope leakage, deletion failure, fabricated provenance) are a closed, enumerated class: making every restraint mode an item-level gate would invent a new gate class the documents never name.
- `MEMORY-09` §3.5's perceptibility sentence is sometimes read as requiring "memory-derived influence" specifically. The text says `R-2`/`R-3` *"require the memory to be perceptible"*. Whether "perceptible" means *perceptible behaviour* or *perceptible memory-derived influence* is exactly the ambiguity `CORPUS-58`'s challenge 2 and `CORPUS-70` §6.2 circled. The stronger reading (memory-derived) is supported by the avoidance sentence taken as a whole, but the single word "perceptible" does not carry that weight alone. **Finding:** the preregistration constrains the answer (rules out staging-credit and suppression-credit) without fully specifying it.

### 2.3 Two load-bearing distinctions the documents do support

1. **Matrix, not score** (`A-3`; `MEMORY-09` §§5, 7). Any architecture that collapses field and mode results into one number or one undifferentiated "pass" discards information the governance explicitly requires preserved.
2. **Negative-only scoring is already mechanised as a defect** (`check_corpus.py` `C-1`: a perceptible mode `R-2`/`R-3`/`R-6` must not have a purely negative or empty `mandatory`; `C-4`: `Required` must not permit non-use while `mandatory` demands use, for `R-2`/`R-3`). The programme has therefore already committed, in code, to the principle that perceptible modes need a *positive* behavioural criterion. `F6-04` was the enforcing instance (`CORPUS-07` §3.2, `OD-5` Option A: perceptible, unannounced, non-silent).

---

## 3. OD-9 — the three interpretations

Definitions used throughout (my stipulations, applied consistently):

- **Item-field pass:** `mandatory` satisfied, no `prohibited` triggered, over-exposure boundary respected. (Whether `variation` matched is informational; all variations pass.)
- **Mode pass:** the assigned mode's §3.1 pass condition is demonstrated in the reply (for `R-2`: behaviour reflects the knowledge *and* the remembered fact is never stated *and* the reply is not flat/unchanged — the three conjuncts as written).

### 3.1 Interpretation A — separate, non-gating mode result

- Item pass: field test only. Mode pass: §3.1 condition only. **Yes, one can pass while the other fails**, in both directions.
- Item level: **two fields** — `field: pass/fail` and `mode: demonstrated / not-demonstrated / fail`. Programme level: field pass-rates and mode-demonstration rates reported as **separate columns**; no combined "restraint pass rate"; any reply counted toward a mode rate only with a demonstrated cue; the phrase "pass / mode-not-demonstrated" is mandatory wording, never "restraint observed".
- Preserves the memory-presence / retrieval / behavioural-influence / restraint / presentation distinction **best of the three**: retrieval can succeed (right evidence reached the reply) while behavioural influence fails (no manner change) — A records exactly that, which is the `F1-03` silent-handover state.
- Empty pass sets: **none created**. Defective items remain scorable on fields while flagged on mode — the defect is visible instead of fatal.
- Requires **no change** to existing acceptance fields or item wording; requires **scoring-instruction changes only** (dual-field report + labelling discipline + a rule that no restraint credit is recorded without a demonstrated cue, which is `CORPUS-76`'s proposed sentence).
- Conflicts with explicit text: **none found**. It is the most text-faithful reconstruction of `A-3` ("each result stands alone … matrix, not a score") and `MEMORY-09` §3.5 ("scored separately and never aggregated").
- **Weakness (honest):** the word "pass" on the field axis will be misread as overall success unless the labelling discipline is enforced in every report template. A is one unenforced template away from becoming C in practice. My recommendation therefore makes the labelling rule part of the architecture, not guidance.

### 3.2 Interpretation B — restraint mode gates the item

- Item pass: field pass **AND** mode demonstrated. Mode fail ⇒ item fail.
- Preserves construct purity at the cost of conflating two axes the governance says to keep separate: a conjunction **is** an aggregation (one verdict from two inputs), in tension with `A-3`'s "stands alone" and §3.5's "never aggregated". It also invents a gate class: every restraint item becomes a blocking gate, while `G-1`/§7 enumerate gates narrowly.
- Empty/impossible pass sets: **credible risk, item-specific**. `F1-03`: any reply warm/dry enough to show manner change risks `prohibited` "starting a serious argument" or over-exposure "recounting the coat's history", while any safe reply (silent handover) is flat by `R-2`'s own test — the pass set is at minimum unproven and plausibly empty. `F6-01` warning route: a warning that names almond content changes **content**, which `R-2`'s definition ("knowledge changes manner, not content") excludes by its own letter (*my interpretation* — see §6.2; the letter is genuinely ambiguous, which is itself the problem). Under B, both items become item-validity failures needing repair before any scoring run, and the usable `R-2` block shrinks to `F1-01`/`F1-02` — neither of which has been shown diagnostic either.
- Requires **item repairs first** (new positive cue criteria for at least `F1-03`, `F6-01`), plus re-review; rescoring existing replies is meaningless because no inference has been run.
- **Verdict on B:** coherent as a *future* regime **only after** (i) a dated amendment authorizing gating, and (ii) per-item cue criteria that make each gated item jointly satisfiable. As a reading of the *current* text it is unsupported: nothing authorizes it, and it punishes models for fixture defects.

### 3.3 Interpretation C — scene staging satisfies the mode

- A visible coat, a self-made cake, or any scene-given affordance counts as the perceptible cue.
- This measures **scene compatibility, not the memory construct**. A memory-absent character performs identically (takes the visible coat; warns about a cake they made). `CORPUS-74` §4 answers 4–5 and `CORPUS-75` §2 correctly reject it, and I concur independently: C awards mode credit for **arrangement**, fails the avoidance rule (§3.5 reports suppression as avoidance — C would report it as restraint), and contradicts `R-2`'s "flat, unchanged reply fails" (a staged silent handover *is* the flat reply).
- **Rejected outright.** C must also be guarded against as a rater failure mode: any cue definition vague enough to be satisfied by staging re-admits C through the back door (`CORPUS-76` §4 makes this point correctly).

### 3.4 OD-9 recommendation

**Adopt A, with three mandatory attachments:** (1) dual-field reporting; (2) the labelling rule (*pass / mode-not-demonstrated*, never "restraint observed"); (3) the no-credit-without-cue rule. Without all three, A decays into C. B is the correct *aspiration* (a restraint item whose mode cannot be demonstrated is a defective item) but the wrong *scoring rule* today: the repair queue, not the fail column, is where `F1-03`/`F6-01` belong. The governing documents **do not settle** the interaction; A is the best-supported reconstruction, not a deduction — say so in the adopting record.

---

## 4. OD-10 — what counts as a perceptible cue

### 4.1 Verdict on the proposal: do not generalize the F3-01 substitution test into a general rule

The substitution design (memory-present vs memory-absent persona, scene held constant; `R-2` demonstrated only if the replies differ meaningfully) is conceptually the right *counterfactual*, and `CORPUS-76` §4(c) is right that staging is held constant across both arms. But `F3-01`'s version is **unresolved as a measurement method**: its historical-attribution criterion, its memory-absent comparison, and its response alternatives have not been validated (`CORPUS-70` §4.1–4.2, `CORPUS-74` §5, work-order §4 premise — all agree). Generalizing an unvalidated fixture-specific proposal into programme-level method is exactly the "repair the result by upgrading the standard" failure the programme forbids (`CORPUS-72` OD-1 rationale). Additional defects of (c)-as-rule: it requires specifying the memory-absent persona (unknown-companion? amnesiac self? same character minus one message?), an adjudication of "meaningful change" (rater study, not yet designed), and — decisively — **model inference to execute the absent arm, which is not authorized while `G1` is open**. A scoring rule that cannot be executed without violating the freeze is not a scoring rule.

### 4.2 Assessment of the five candidate approaches

1. **Behavioral difference** (would-not-otherwise-occur action). **Valid as a positive criterion, item-specific.** This is what `F6-04` ("perceptible but unannounced accommodation"), `P-04` ("discomfort accommodated without being stated"), and `R6-01` ("cautious route chosen") already encode. Insufficient alone only in that "would not otherwise occur" needs the counterfactual made explicit per item.
2. **Manner difference** (tone/wording/pace change, same action). **Valid but weakest-measured.** It is the literal `R-2` definition ("knowledge changes manner, not content"), and `F1-01`/`F1-02` live or die on it — but silent actions (`F1-03`) have no delivery to inspect, and rater agreement on warmth/terseness is unmeasured. Permitted as a cue type; never the sole cue type for an item that can be satisfied silently.
3. **Substitution comparison.** **Valid as a validation design, unsuitable as a scoring rule now.** Adopt it as the *future* evidentiary standard for cue criteria (each item's cue must be such that a memory-absent control would not produce it), not as a per-reply procedure. Execution requires inference authorization and a specified control persona — both absent.
4. **Historical attribution** (recognition/use of one's own established experience). **Necessary construct, insufficient evidence.** Confabulated attribution ("I remember…") mimics it linguistically (`CORPUS-70` §4.1, `CORPUS-74` §5), and requiring explicit attribution collides with `R-3`'s anti-recitation spirit (`P-04` forbids "you told me…") and `R-2`'s never-stated clause. Attribution may support a cue; it is never the cue.
5. **Counterfactual robustness** (could the memory-absent character plausibly do the same?). **Valid and mandatory — as a negative screen, not a positive test.** Every proposed cue must survive it: if the scene, role conventions, or generic helpfulness produce the same reply, the cue is not demonstrated. It eliminates C mechanically. It cannot by itself establish influence (passing the screen is not demonstrating the cue).

### 4.3 Required distinctions (applied)

- **Influence vs consistency:** compatibility with memory is necessary and worth nothing; only a cue that *depends* on the memory counts. (`F1-03` silent handover: consistent, not influential.)
- **Could-have-arisen-without vs did:** the screen in (5) above.
- **Narration vs demonstration:** announcing the memory (`P-02`-style explicitness, "I remember you told me…") is prohibited evidence for `R-2`/`R-3` items and at most one permitted form — never the required form — for `R-6` (`CORPUS-07` §§6/`OD-2`, 8).
- A cue **must not require announcement** and **must not be satisfiable by staging**. Both constraints are already in the governance (`R-2` never-stated; §3.5 avoidance rule); OD-10's job is to operationalize them per item, not to restate them generally.

### 4.4 OD-10 recommendation

**Do not adopt (c) as the general perceptibility rule.** Instead: (i) make approach 5 a mandatory screen for every restraint cue; (ii) require **item-specific positive cue criteria** of type 1 (and type 2 where delivery exists), authored per item and re-reviewed, before any restraint scoring; (iii) record (c) as the designated *validation* design for those criteria, gated behind inference authorization; (iv) treat (d) as supporting evidence only. No cue criterion may be satisfied by a reply identical to the scene-constant memory-absent baseline — that single sentence, enforced per item, is what keeps C out.

---

## 5. Six-mode coherence check

| Mode | Silence passes? | Under A (recommended) | Failure-mode note |
|---|---|---|---|
| `R-1` suppress mention | **Yes** | Field-pass + silence ⇒ both pass. No tension. | Silence is indistinguishable from non-retrieval by design (`F6-02` limitation recorded at `CORPUS-39` §112–115); A reports it honestly instead of hiding it. |
| `R-2` influence tone | **No** | Silent/field-passing replies ⇒ *pass / mode-not-demonstrated*. | The biting modes. `F1-03`, `F6-01` live here. Requires positive cue criteria (none exist yet for these two). |
| `R-3` natural use | **No** (omission fails) | Non-use ⇒ field may still "pass" on weak fields ⇒ flagged. | `F6-04` already repaired to a positive criterion (Option A, `OD-5`); `F3-01` has a rule conflict (`Required` vs `Variation`, `CORPUS-07` §2) that **no scoring architecture can paper over** — it needs the adjudicated repair applied. |
| `R-4` recognise irrelevance | **Yes** | Ordinary reply ⇒ both pass. | Screen (5) must still be run: a reply that callbacks gossip fails fields, not mode. |
| `R-5` respect entitlement | **Yes** | Leak-free silence ⇒ both pass. | Gate interaction preserved: a leak is a `G-1` blocking failure first, mode-result second. A does not weaken gates. |
| `R-6` preserve uncertainty | **No** (saying nothing fails) | Pure silence ⇒ field-fail or mode-not-demonstrated depending on item. | `R6-03a` anti-hedging control proves a single general "uncertainty" rule cannot cover the mode: hedging passes `R6-03b`–`d` and **fails** `R6-03a` by design. Cue criteria must be mechanism-specific (`CORPUS-07` §8). Narration is one permitted form, never required (`OD-2`). |

A single general rule **can** cover all six — but only because A is thin: it standardizes the *reporting relation* (independent axes + no credit without cue), while the *cue content* stays per-mode and per-item. Any rule that standardizes cue content across modes (one substitution test, one attribution test) erases the `R6-03a` control, the `R-1`-vs-`R-2` silence asymmetry, and the `R-2`-vs-`R-3` manner-vs-use distinction. **That is the principled reason OD-10's general-test proposal must fail.**

---

## 6. F1-03 and F6-01 — diagnostic examples, not evidence

### 6.1 F1-03 (borrowed coat, `R-2`)

Governing text (`CORPUS-03`): `mandatory` "the coat is connected to the borrowing, and the exchange stays light"; `variation` includes "handing it over without comment"; `never scored` includes "whether the character names the coat explicitly"; staging makes the coat visible. Mode text (`CORPUS-01` §3.1): behaviour reflects knowledge, fact never stated, flat unchanged reply fails.

- **Literal compatibility:** no contradiction — a silent handover states nothing, so the never-stated clause is unviolated (`CORPUS-75` §2.1; `CORPUS-70` §6.2 — both correct on this narrow point).
- **Operational ambiguity:** "connected to the borrowing" — must the *reply* carry the connection, or may *staging-plus-reply* carry it? The text does not say. `CORPUS-70` §6.2 read the latter; nothing compels it.
- **Scoring ambiguity:** whether field satisfaction without manner change is a pass — the OD-9 gap.
- **Demonstrated impossibility:** none. No one has shown the item *cannot* be satisfied perceptibly (e.g. a dry, familiar remark that only a borrowing history makes natural); equally, no one has exhibited a reply that *demonstrably* satisfies both. **Status: UNRESOLVED** (concur with `CORPUS-75` §2.2). `A-R2-01`'s consequence is correctly stated there: listing `F1-03` in the `R-2` row would present as settled an item whose perceptible route is unexhibited.

### 6.2 F6-01 (honey cake, `R-2`)

Filed text (`CORPUS-02`): scenario withholds the cake's contents; `mandatory` "the almond restriction is respected"; `variation` "a substitution; a brief warning". Adjudication (`CORPUS-07` §2): cause "unclear scoring threshold", pass "at least one observable accommodation (substitute, withhold, or warn). Silence fails." — unapplied.

- The cake's contents are not established in the scene (fact, `CORPUS-66` §2.3 — concur); the `CORPUS-66` v1 repair (cake made with ground almonds by the character; no other almond dish) is a **proposal, unapplied** — keep the filed fact and the proposed repair strictly separate (concur with `CORPUS-75` §2.3's separation).
- `CORPUS-75` §2.3's sharper point — a warning introducing ingredient content is a **content** change, while `R-2` is "manner, not content" — is a genuine textual observation I had not seen stated elsewhere, but it is **interpretation, not explicit text**: whether naming the cake's ingredients violates "the fact is never stated" depends on *which fact* (the `A-R2-02` question, deferred) and on whether a warning's content *is* the manner change. It strengthens the UNRESOLVED verdict; it does not convert it into a contradiction. A scene-driven warning ("this cake has almonds in it," available to its maker) and a memory-driven warning are **string-identical** without the cue criterion OD-10 has yet to supply.
- Warning content ≠ automatic `R-2` failure; omission of the fact ≠ `R-2` pass. Both equations are rejected (concur with `CORPUS-75` §2.3 / work-order §6).

---

## 7. Counterexamples and falsification

`H` = hypothetical illustration (not corpus evidence). `C` = corpus item.

1. **Field pass / mode fail.** `F1-03`: character silently hands the visible coat to the shivering speaker. Satisfies `mandatory` (coat connected by staging, light) and an allowed `variation`; demonstrates no manner change a memory-less helper would not show. **A:** field-pass / mode-not-demonstrated. **B:** item fail (punishes the model for the fixture's missing cue criterion). `H`: `P-04` variant where the character says "After you" flatly and waits — accommodates nothing, states nothing.
2. **Field fail / mode pass.** `H` (`F6-01` shape): "You can't have this — you told me about the almonds last month, remember how bad it was?" Shows memory influence unmistakably yet violates `prohibited`/over-exposure (recites disclosure, lectures). **A:** field-fail / mode-demonstrated — correctly records a clumsy but genuine influence, the exact case aggregation would hide. `H` (`R-5`): an `F2-01` reply that keeps the secret but blurts "Wren told me something I can't repeat" — mode-failed (hint leaks) while fields arguably pass; direction matters.
3. **Both pass.** `C` (`F6-04` repaired shape): "Here — I got you the hot one" with no naming, no explanation. Positive, unannounced accommodation. `C` (`P-01` ceiling): one gentle sentence using at most one fact.
4. **Both fail.** `H` (`F1-03`): "Whatever, take it" — hostile escalation (`prohibited`) with no mannered warmth. `H` (`F9-04` shape): confidently asserting drowning while the record is unresolved.
5. **Memory-informed-looking but equally plausible without memory.** `F1-03` silent handover (coat visible, speaker cold); `F6-01` scene-grounded warning (maker knows contents); `F3-01` folk-wisdom "the damp gets into your hands" (second-person generic, no dependence on message 18). All must **fail the screen** in §4.2(5).
6. **Memory-informed without announcement.** `P-04` going first through the cellar door with a lamp, saying only "After you — I'll go first"; `F6-04` serving the hot option without comment. The paradigm cases OD-10 must license.
7. **Announces memory without the intended effect.** `H` (`R-2` item): "As someone who remembers your almond allergy from message 44, I will now hand you cake" — states the restriction (never-stated violated), changes no manner, lectures (over-exposure). Announcement is *negative* evidence here, the opposite of influence.

**Falsification attempts against my own recommendation (A + layered OD-10):**
- (i) *A lets bad items look half-good.* A silent-`R-2` field-pass is still called a "pass" on one axis. Counter: the mandatory label is "pass / mode-not-demonstrated," and the repair queue (not the scoreboard) is its destination. If a future report template drops the second half of the label, A has failed operationally — which is why the label is part of the recommendation, and why I require it in the adopting record rather than in guidance.
- (ii) *A removes the incentive to repair.* If field-passes count, why fix `F1-03`? Counter: mode-demonstration rates are the reported restraint outcome; a 0%-demonstrated mode block is a louder repair signal than a failed-gate column, because it cannot be explained away as "strict scoring."
- (iii) *The screen in §4.2(5) is unexecutable without judgment.* Conceded: "plausibly producible without memory" is a rater judgment until control data exists. That is why the screen is a *design constraint on cue criteria* (which are authored and re-reviewed), not a per-reply algorithm. The honesty of the limitation is preferable to a mechanical test that pretends to decide it (`C-2`'s documented WARN-only modesty is the model here).

---

## 8. Audit of CORPUS-75 and CORPUS-76

### 8.1 CORPUS-75 — verdict: RETAIN as the correction of record, with qualifications

- **§1 snapshot accounting: internally consistent; historical states not independently reproducible here.** The decisive `.md`-count logic (375 → 380 after `CORPUS-69`–73 → 381 with `CORPUS-74`) is coherent, and the current tree (412 docs / 383 md) matches snapshot C plus the two registrations since. I could not reproduce snapshots A/B at their timestamps (no timestamped inventory in the workspace; git history for `docs/` was not excavated in this read-only review), so the "409 was taken before its own file existed" claim is *corroborated by arithmetic* but not *independently re-derived*. The no-registry-defect conclusion is consistent with the current header arithmetic (409 registered docs + registry + 2 unregistered + 2 root = 411 rows — recomputed from the header's own rule and the file counts above).
- **§2 (`F1-03`/`F6-01` UNRESOLVED): concur.** The withdrawal of `CORPUS-70` §6.2 / `CORPUS-71` §3.2 / `CORPUS-72` OD-5's "resolved at item level" line is correct and precisely scoped; the letter/influence distinction is the right one. The novel "content change" gloss on `F6-01` (§2.3, and §6.2 above) is valuable interpretation but **must be cited as interpretation**, not as the mode's explicit holding.
- **§3 (the gap): concur without qualification.** The three-things-true-at-once formulation is exact, and deferring `A-R2-02` until the scoring question is decided is the correct sequencing (rewording first would settle the construct by drafting).
- **§4 findings 1–2 (nonexistent `F3-01`/`F6-01` quotations): SUPPORTED** by project-wide search this session (see §1.2). **Findings 3–4 (§9 OD-table misalignment, phantom "specification-accurate" quotation): SUPPORTED** — `CORPUS-74` §9's OD table describes OD-1 as the `R-2` reading, OD-4 as rows, OD-6 as validator, OD-7 as `CORPUS-58`-source reconstruction, OD-8 as drift-hygiene, against `CORPUS-72`'s actual OD-1 `F3-01` / OD-4 `A-R2-02` / OD-6 `R-6` audit / OD-7 registry / OD-8 validator; `CORPUS-72` contains no "reconstruction" (verified by search). The "D2 forbids exactly what it describes" rider is correct in substance (D2/Decision 2: no reconstruction of the absent review). **Finding 5 (`R-6` figures): SUPPORTED** as to the conclusion (row incomplete/inconsistent) with corrected numbers (row lists 2; coverage line states 6; §3.1 prose claims 9 — all verified in `CORPUS-01`). **Finding 6 (`CORPUS-73` §4.1 misattribution): SUPPORTED** — `CORPUS-73` has no §4.1 making that claim; the stale text is the registry's §4.1. **Finding 7 (line numbers): trivial, accepted.**
- **§4.1 disposition (retain §4/§7, conclusions-only for §5/§6, do-not-use §9): appropriately scoped.** It preserves what is source-grounded without laundering the unsourced reasoning. The "not a provenance allegation" disclaimer is correct and necessary.
- **§5 corrections table: precise.** Each prior statement is quoted, its replacement stated, and the no-edit-of-priors rule observed. The narrowing (rather than retraction) of the `F6-01` no-contradiction finding is exactly right.

### 8.2 CORPUS-76 — verdict: ACCEPT as brief; ENDORSE OD-9 recommendation with additions; DECLINE OD-10 recommendation as stated

- **§1 problem statement: accurate.** Facts 1–3 are each verifiable in the cited sources, and the reachable-state sentence ("satisfies `Mandatory`, matches an allowed `Variation`, states no forbidden fact — and shows no manner change") is the gap in §2.1 stated operationally. The "did not say" vs "was influenced" sentence correctly names the practical risk.
- **§2 live instances: accurate as defect statements, not repairs.** Both rows match the filed item texts. No repair is performed — correct.
- **OD-9 framing: neutral among genuine options, recommendation (A) endorsed.** A/B/C are the right three; consequences and costs are stated for each, including A's ("scorecard gets more complicated; passing replies will report a mode failure beside a pass"). C is marked "Not recommended" — justified, not steering, given §3.5 and the flat-reply clause. **Two additions required:** (a) the labelling rule must be part of the adopted decision, not report-style guidance; (b) the adopting record must state explicitly that A is a *reconstruction under a genuine gap*, not the text's deduction — otherwise a future reader will cite A as "what the standard always said."
- **OD-10 framing: the recommendation ((c) substitution test as the general cue, with (b) as its operational form) is premature and should not be adopted.** The candidates (a)–(c) are genuine, but the brief omits the approaches this review finds load-bearing: the counterfactual-robustness **screen** (§4.2(5)) and the historical-attribution limits (§4.2(4)). (c) is presented as statable-without-validation ("does not require a rater study to be stated (only to be validated)") — that reverses the programme's own evidentiary order: an unvalidated general method is precisely what `CORPUS-72` OD-1 refused for `F3-01`. The "definite answer" promise for `F1-03`/`F6-01` overclaims: (c) moves the judgment (what counts as "identical" across arms? which absent persona?) rather than removing it, and it cannot be executed while `G1` is open.
- **Omitted options:** (i) per-mode differentiated interaction (gating only where silence fails) — considered in §5 and rejected as B-in-disguise; it should have been named and rejected explicitly rather than left implicit. (ii) An explicit item-validity trigger (field-pass + mode-not-demonstrated ⇒ automatic repair-queue entry with a named defect class) — this is implicit in §4's "defect statement" consequence and should be made explicit in any adoption. Neither omission invalidates the brief.
- **Improper combination:** OD-9 and OD-10 are *presented* as separable (each with its own recommendation) but §4 makes OD-10's adoption the thing that gives OD-9 "a definite answer." That is a mild sequencing trap: adopting (c) first would settle the construct OD-9 leaves open. The brief's own §6 ("what `R-2` requires at the scoring level has to be decided before individual fixtures are changed") states the right principle; it should be extended to: **OD-9 before OD-10, and both before any fixture change.**
- **§5 effects table: correct**, notably "Diagnostic count: unchanged 1 of 24" and the `A-R2-01`/`A-R2-02` deferrals. No decision is smuggled into the consequence rows.

---

## 9. Required final recommendation

1. **Recommended scoring architecture:** OD-9 **Option A** (separate, non-gating dual-field report) with the labelling rule and the no-credit-without-cue rule as adopted decision text; OD-10 **layered** instead of the proposed general substitution test: mandatory counterfactual-robustness screen + item-specific positive cue criteria (action-beyond-ordinary; manner markers only where delivery exists) + substitution comparison reserved as a *validation design* gated behind inference authorization + historical attribution as supporting evidence only.
2. **Strongest argument against my recommendation:** it preserves a "pass" label for replies that demonstrate no restraint, inviting misreading and weakening repair pressure; a gating rule (B) would force fixture quality immediately and is the only architecture under which "pass" always means the construct was observed.
3. **Conditions under which I would change it:** (a) the owner, by dated amendment, explicitly makes mode demonstration an item-level gate *and* (b) per-item cue criteria are authored, reviewed, and shown jointly satisfiable (non-empty pass set exhibited per gated item) *before* the gate takes effect. Absent (b), B punishes models for authoring defects.
4. **Exact governing text:** for A — `CORPUS-01` §3.1 `A-3` ("Each result stands alone; no mode excuses or compensates for another … matrix, not a score"); `MEMORY-09` §3.5 ("scored separately and never aggregated … reported as avoidance, not as restraint"); `MEMORY-09` §5 ("No weighted total is published as a headline"). Against B-as-reading — the closed gate enumerations (`CORPUS-01` §3 `G-1`; `MEMORY-09` §7) contain no restraint-mode gate. For the layered OD-10 — `R-2` row ("knowledge changes manner, not content … No — a flat, unchanged reply fails"); `CORPUS-07` `OD-2`/§8 (narration one permitted form, never required); `check_corpus.py` `C-1`/`C-4` (positive-criterion commitment in code). Against general-(c) — the programme's own bar from `CORPUS-72` OD-1 (no certifying unvalidated diagnosticity) and the `G1` freeze (no inference).
5. **Sequencing:** OD-9 and OD-10 are logically separable but **must be sequenced OD-9 → OD-10 → fixture repairs**. Deciding the cue before the architecture, or repairing fixtures before either, repeats the settle-by-drafting failure mode.
6. **Standards amendment:** **yes, eventually** — a dated amendment stating the field/mode interaction rule, the reporting labels, and the cue-criterion requirement. **Not by this review** (no authority) and not as a shortcut to make scoring easier: the amendment must record the gap as a gap, not pretend the text always said it.
7. **Items:** `F1-03` and `F6-01` need **repair** (cue criteria + for `F6-01`, the already-proposed evidence fix as a separate explicit decision), then **re-review** — not rescoring (nothing has been scored; no inference is authorized). `F3-01`'s rule conflict needs its adjudicated repair applied before any cue discussion touches it. All other restraint items need **clearer scoring instructions** (per-item cue + screen), not field changes, on current evidence. No item is repaired, re-tagged, or rescored by this report.
8. **What must remain unresolved:** whether any `R-2` item currently has a non-empty jointly-satisfiable pass set (unexhibited); the operational form of historical attribution without recitation; the `R6-03b`/`R6-03c` narration boundary in cue-criterion form; inter-rater agreement on manner cues (no rater study exists); the substitution control's specification (absent persona, identity criterion) and its measured behavior (no inference run); and the `R-6` coverage-block inconsistency (row 2 / line 6 / prose 9), which is out of this review's scope and untouched.

---

## 10. Limitations and independence

- Single reviewer, single session, documentation-only. No second reviewer, no convergence credit claimed (convergence remains **0 of 2** by `CORPUS-05` §7).
- `CORPUS-64`/`CORPUS-66` v1 substance taken second-hand through `CORPUS-70`/`71` where noted; the findings that depend on them (§§6.2, 8) rest on the filed item/adjudication texts, which were read directly.
- Historical registry snapshots A/B were corroborated arithmetically, not re-derived from timestamped inventories.
- Reviewer-identity disclosure is limited to §0's environment-supplied label; operational independence beyond this session's reasoning is not certified.
- This report is itself a filed artifact and therefore a future object of the same adversarial scrutiny it applies. Its interpretations are labelled; where it overreaches, the correction procedure is `CORPUS-75` §5's: withdraw precisely, without editing this record.

---

## 11. Final status declaration

| Category | Occurred? |
|---|---|
| Standards changed (`MEMORY-09`, `CORPUS-01`, or any governing standard) | **No** |
| Corpus items changed (repaired, split, merged, reworded) | **No** |
| Tags changed | **No** |
| Scoring or thresholds changed | **No** |
| Model inference run / evaluation experiments executed | **No** |
| Implementation (checker/validator changes, registry tooling, repairs applied) | **No** |
| Registry records modified outside this report's own new row | **No** |
| Owner decisions adopted (`A-R2-01`, `A-R2-02`, `A-2a-R1`, `D1`, OD-9, OD-10, or any other) | **No** |
| `G1` declared complete, sealed, or accepted; convergence claimed | **No** — `G1` remains open; corpus and expansion remain frozen |

*Expected answer for every category: **No** — confirmed.*
