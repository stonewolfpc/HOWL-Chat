# MEMORY-09 S3 — bge-reranker-v2-m3: Bounded Evaluation

**Identity:** `MEMORY-09-S3-BGE-1` · **Date:** 2026-11-15
**Authorization:** single pinned Q8_0 artifact, bounded evaluation, not production inclusion.
**Gate B: 6/7 — stays open.**

---

## Five outcomes, kept separate as required

| # | Stage | Outcome |
|---|---|---|
| 1 | Artifact integrity | **PASS** |
| 2 | Model / runtime compatibility | **PASS** |
| 3 | Instrument sanity | **PASS** |
| 4 | Frozen behavioral evaluation | **INCONCLUSIVE under the frozen rule**; diagnostic reading obtained |
| 5 | Resource envelope | latency **PASS** (batched) · VRAM **FAIL** for full residency · CPU/partial **untested** |

The separation matters: stage 4's inconclusiveness is an **apparatus scale problem**, not evidence against cross-encoders, and stage 5's VRAM failure is arithmetic, not a capability verdict.

## 1. Artifact integrity — PASS

| Field | Value |
|---|---|
| Repo / revision | `smarttasks/bge-reranker-v2-m3-GGUF` @ `fc50d873dcd84a7be9760cfecd31f513d5ebc436` |
| File | `bge-reranker-v2-m3-Q8_0.gguf` |
| Bytes | 635,674,432 — exact match |
| SHA-256 | `d32a877b7266f5275172458616824ce6728502416672052c2da795da3389201e` — exact match |
| Licence | Apache-2.0 |

## 2. Compatibility — PASS

`--reranking` on b10076: `load_tensors: offloaded 25/25 layers to GPU`, `CUDA0 model buffer = 308.30 MiB`, `/v1/rerank` responds, 393 tensors, GGUF V3.

Two operational findings:

- `pooling type = -1` in metadata, `[4]` (rank) applied by flag → warning, but scoring works and discriminates.
- **`/v1/rerank` returns results sorted by score, not in input order.** Adapters must map through `index`. A positional read silently permutes scores — this bit my own first comparison.

## 3. Instrument sanity — PASS

| Check | Result |
|---|---|
| Joint scoring (same doc, two queries) | −5.954 (trainer q) vs −11.049 (tea q), **Δ 5.096** |
| Rank discrimination | relevant −5.954 above irrelevant −11.043 |
| Reproducibility | three calls, identical to 4 dp, drift 0.0 |
| Order invariance | reversing documents preserves the ranking |
| Non-degenerate | finite, non-constant |

**Scale:** raw logits, observed −11.04 … +1.20. Sigmoid maps them to 0.000016 … 0.77. **The frozen threshold of 0.30 does not exist on this scale.**

## 4. Frozen behavioral evaluation

### Reading 1 — frozen rule applied unchanged: 6/9, "R4 veto broken"

| Case | Top | Raw score | Result |
|---|---|---:|---|
| R1-horse | none | −1.868 | PASS (by accident — see below) |
| R1-patient | none | −1.535 | FAIL abstained |
| R2-vault-foreign | none | −11.011 | PASS |
| R2-vault-own | rec-0004 | 0.586 | PASS |
| R3-now | rec-0007 | 1.200 | PASS |
| R3-before | none | −1.606 | FAIL abstained |
| R4-paraphrase | none | −5.954 | **FAIL** (veto) |
| R5-weather | none | −11.011 | PASS |
| R5-fair | none | −10.884 | PASS |

This reading is a **scale artefact**: only scores above 0.30 survive, and most of this model's output is negative. It is reported because the rule is frozen, not because it is interpretable.

Note the trap it creates: the horse case "passes" while a correct case fails, purely because one logit happened to sit above a threshold defined for a different score family.

### Reading 2 — ranking diagnostic, no threshold invented

| Case | Top record | Score | Margin | Reading |
|---|---|---:|---:|---|
| R1-horse | **rec-0001 (trainer)** | **−1.868** | 9.064 | **subject confusion present** |
| R1-patient | rec-0001 | −1.535 | 9.497 | right record first |
| R2-vault-foreign | rec-0006 | −11.011 | 0.018 | effectively rejects |
| R2-vault-own | rec-0004 | 0.586 | 11.624 | right record first |
| R3-now | rec-0007 | 1.200 | 1.704 | current beats superseded |
| R3-before | rec-0002 | −1.606 | 1.122 | history reachable |
| R4-paraphrase | rec-0001 | −5.954 | 5.003 | right record first |
| R5-weather | rec-0006 | −11.011 | 0.023 | effectively rejects |
| R5-fair | rec-0001 | −10.884 | 0.058 | effectively rejects |

### The narrow question, answered

The same trainer record scores:

- **−1.868** under "Who trained the horse?"
- **−5.954** under "Who was your trainer, again?"

The horse query matches the trainer record **4.086 points more strongly** than the correct question does. bge does not merely fail to resolve the subject confusion — **it amplifies it**, confidently ranking the trainer record for a question about a horse. R2, R3 and R5 behave well, so this is specific to agent/patient discrimination, not a broken model.

**Caveat that must travel with this result:** records were fed as `Subject predicate: value` strings, the same text the other branches see. A cross-encoder is trained on natural prose; schema-shaped text may depress or distort its judgements. This is a limitation of the comparison, not a licence to reword until it passes. Testing natural-language renderings is a **new experiment and a new spec version**.

## 5. Resource envelope

| Quantity | Measured | Budget | Verdict |
|---|---:|---|---|
| Spawn → ready | 2.54 s | ≤ 3.0 s | **PASS** |
| First call after load | 21 ms | — | — |
| Warm p50 (single pair) | 11 ms | ≤ 40 ms | PASS |
| Warm p95 (single pair) | 30 ms | ≤ 120 ms | PASS |
| Batched 4 docs p50 / p95 | 30 / 43 ms | 40 / 120 ms | **PASS** |
| Batched vs sequential equivalence | max Δ 0.0011, identical ranking | — | PASS |
| VRAM delta, 25/25 offloaded | **+433 MiB** (308.3 MiB weights; earlier run +680) | ≤ 384 MiB | **FAIL for full residency** |
| Residual after shutdown | −9 MiB | ≤ 64 MiB | PASS |
| RSS | 933 MB | ≤ 1.0 GiB | PASS (tight) |
| Port closed after shutdown | true | — | PASS |

Batching keeps this inside the latency budget; sequential scoring of 4–6 candidates would not. The VRAM result confirms §3 of the candidate review: full residency does not fit, and **partial offload remains untested**.

## 6. What this establishes

1. The artifact is exactly what the publisher published.
2. llama.cpp serves it correctly and scores jointly — the runtime is not the confound.
3. A real cross-encoder is **fast enough** on this machine when batched, and releases cleanly.
4. **bge-reranker-v2-m3 does not fix subject confusion on these records; it strengthens the wrong association.** The R1 failure is not a threshold problem and not an instrument problem.

## 7. Demonstrated defect in the frozen apparatus (owner decision required)

The spec's selection rule (`keep ≥ 0.30`) implicitly assumes a [0,1] score family. A real cross-encoder emits logits, so the rule is undefined for it. I did **not** amend the spec mid-run; the candidate was measured under the rule as written, plus a clearly separated diagnostic.

Proposed amendment, approved by the owner and implemented as **v1.2** (see §9) with the case set unchanged. v1.1 and its receipt remain byte-immutable at `eval/s3/rerank-spec-v1.1.json` (`bbf400c2…`); v1.2 lives in a separate module (`eval/s3/rerank_eval_v12.py`, spec `eval/s3/rerank-spec-v1.2.json`, harness `63cb7aec…`) so nothing was edited in place.

## 9. Result under protocol v1.2

Rules: candidate declares its scale before scoring; **ranking** judged threshold-free; **abstention** judged only with a threshold calibrated on held-out pairs (`eval/s3/calibration_set.py`, different phrasings and different pairings, never the nine cases); overlapping bands would be reported UNCALIBRATABLE rather than fitted.

**Calibration:** min relevant −4.81, max irrelevant −11.033, separation 6.223, threshold **−7.921**, calibratable.

| Verdict | Result |
|---|---|
| **Ranking** (threshold-free) | **5/5** answerable cases rank the expected record first; R4 ok |
| **Selection** (calibrated abstention) | **8/9**; veto broken: **none** |
| **Subject distinction (R1-horse)** | **FAIL** — trainer record ranked first at −1.868, far above the −7.921 threshold |

Per case under v1.2: R1-patient, R2-own, R3-now, R3-before, R4-paraphrase all select correctly; R2-foreign, R5-weather, R5-fair all abstain; **R1-horse selects the trainer record**.

Two things follow that v1.1 could not show, because its rule collapsed everything to abstention:

1. Calibration is fit for purpose here — the relevant/irrelevant bands separate by 6.2 logits on pairs the model never saw in the case set.
2. The paraphrase case (R4) passes at −5.950, above the calibrated threshold, so the earlier v1.1 "veto failure" was an artefact, not a real regression.

**Calibration did not launder the bad ranking.** R1-horse is not a threshold problem: the model puts the trainer record 6.05 logits above the calibrated threshold for a question about a horse. No monotone transform can correct that.

**Do not collapse these numbers.** "5/5 ranking" means the five answerable cases ranked their expected record first. It is **not** a subject-distinction pass: R1-horse is an adversarial negative case, it selects the wrong record, and the overall selection figure is 8/9. Any summary that merges these three into one headline is wrong, including a flattering one.

## 10. Experiment E-PROSE-RENDER — representation is not the cause

Only the rendering changed: same model, runtime, nine cases, calibration discipline. Schema text (`Sable trained_by: …`) was replaced by prose (`Sable was trained by …`) for candidate documents, with the calibration pairs rendered the same way. Schema results stay in `rerank-v12-receipt.json` and are **not superseded**.

| Case | Schema top / score | Prose top / score |
|---|---|---|
| **R1-horse** | rec-0001 / **−1.868** | rec-0001 / **−1.247** |
| R1-patient | rec-0001 / −1.535 | rec-0001 / −0.956 |
| R2-vault-foreign | rec-0006 / −11.010 | rec-0001 / −11.034 |
| R2-vault-own | rec-0004 / 0.586 | rec-0004 / **1.675** |
| R3-now | rec-0007 / 1.200 | rec-0007 / **3.871** |
| R3-before | rec-0002 / −1.606 | rec-0002 / **0.992** |
| R4-paraphrase | rec-0001 / −5.950 | rec-0001 / −5.536 |
| R5-weather | rec-0006 / −11.016 | rec-0006 / −11.038 |
| R5-fair | rec-0001 / −10.884 | rec-0001 / −11.031 |

Prose calibration: separation 4.623, threshold −8.714.

| Metric | Schema | Prose |
|---|---|---|
| Ranking (answerable) | 5/5 | **5/5** |
| Selection (calibrated) | 8/9 | **8/9** |
| Subject distinction (R1-horse) | FAIL | **FAIL** |

**Interpretation.** Prose raises the model's confidence on correct matches substantially (R3-now 1.200 → 3.871; R3-before −1.606 → 0.992) but leaves the horse case wrong **and more confident** (−1.868 → −1.247). The schema rendering was therefore **not** the cause of the subject-role failure. That moves the finding from "possibly an artefact of input format" to a **model-specific limitation under the tested conditions**, and rules out further rendering work as a fix.

Receipt: `eval/s3/prose-experiment-receipt.json`.

## 11. Resource envelope completed — CPU, partial offload, full GPU

Same harness (`eval/s3/measure_reranker.py`), batched 4-document scoring per query, idle VRAM median sampled 5×, residual sampled 5× after shutdown. Behaviour was **not** re-measured; §9 and §10 stand unchanged.

| Config | Offload | VRAM Δ | Batched p50 / p95 | RSS | Ready | Residual | Verdict |
|---|---|---:|---|---:|---:|---:|---|
| CPU (`--device none -ngl 0`) | 0/25 | **+87 MiB** | 282 / 292 ms | 775 MB | 2.18 s | −25 MiB | VRAM pass, **latency fail** |
| Partial 6 | 6/25 | **+198 MiB** | 102 / 118 ms | 834 MB | 2.36 s | −11 MiB | VRAM pass, **latency fail** |
| Partial 12 | 12/25 | **+277 MiB** | 70 / 93 ms | 833 MB | 2.46 s | −2 MiB | VRAM pass, **latency fail (p50)** |
| Full GPU | 25/25 | **+444 MiB** | **20 / 44 ms** | 830 MB | 2.42 s | +1 MiB | **VRAM fail**, latency pass |

Model load 1.92–2.12 s; spawn→ready 2.2–2.5 s, inside the 3 s budget. Shutdown releases cleanly in every configuration.

### The finding: a budget gap, not a tuning problem

**No configuration satisfies both the 384 MiB VRAM budget and the 40 ms p50 latency budget.** Offload is smooth and monotone — 282 → 102 → 70 → 20 ms as layers move to GPU — but the two constraints cross between partial-12 and full offload. Partial-12 is 30 ms over on p50; full GPU is 60 MiB over on VRAM.

Three readings are available to the owner, and they are a policy choice rather than a measurement:

1. **Relax latency** — CPU at 282 ms per query, 87 MiB, no VRAM dependence. Acceptable only if reranking sits off the conversational critical path.
2. **Relax VRAM** — full GPU at 444 MiB, 20 ms. That is 60 MiB above the registered budget and must be recorded as a budget amendment, not absorbed silently.
3. **Retain partial-12** — 277 MiB, 70 ms, and treat the missed p50 as a known limitation.

### Retention recommendation

bge is **technically retainable as an optional component and behaviourally limited**. Its resources are now fully characterised, and its subject-role failure is unchanged by representation (§10). Whether to keep it depends on the policy choice above, not on further measurement. It must remain optional: formation, lexical retrieval, semantic retrieval when available, and abstention all work without it.

## 12. System-level A/B: does it actually make memory better?

Isolated model tests do not answer the product question. This compares the **memory system** with and without reranking, holding queries, records, ownership boundaries and the case set fixed. Only the final selection differs.

- **A** — existing retrieval: lexical + subword, RRF, choosing from the case pool (no models).
- **B** — same pool, final selection by bge with abstention at the held-out calibrated threshold (−7.921).

| Case | Expect | A (pipeline) | B (+bge) | Change |
|---|---|---|---|---|
| R1-horse | abstain | rec-0001 | rec-0001 | – |
| R1-patient | rec-0001 | rec-0001 | rec-0001 | – |
| R2-vault-foreign | abstain | none | none | – |
| R2-vault-own | rec-0004 | rec-0004 | rec-0004 | – |
| **R3-now** | **rec-0007** | **rec-0002 (superseded)** | **rec-0007** | **IMPROVED** |
| R3-before | rec-0002 | rec-0002 | rec-0002 | – |
| R4-paraphrase | rec-0001 | rec-0001 | rec-0001 | – |
| R5-weather | abstain | none | none | – |
| R5-fair | abstain | none | none | – |

| Metric | A | B |
|---|---|---|
| Correct retrieval (answerable) | 4/5 | **5/5** |
| False retrievals | 1 | **1** |
| Abstention correct | 3/4 | **3/4** |
| Total | 7/9 | **8/9** |
| Regressions | — | **none** |
| Latency p50 (this pool) | 1 ms | 34 ms |

### Reading it honestly

**The improvement is narrow and specific.** Reranking fixes exactly one thing here: the pipeline returns the *superseded* blue-cabinet record for "where is the ledger kept now?", and bge returns the current green one. That is a real memory-correctness improvement, not a benchmark artefact — the pipeline had no signal for temporal status at selection time.

**Nothing regressed.** Every case the pipeline already won, bge keeps.

**The horse failure is untouched and stays visible.** Both configurations select the trainer record for "who trained the horse?". Reranking neither fixes nor worsens subject confusion; it is the same false retrieval in both columns, and it is the one abstention miss in both.

**Latency is a separate axis.** 1 ms → 34 ms p50 on this pool. Whether that is acceptable is a product decision about where reranking sits in the turn, not a verdict on the capability.

**Scope.** Nine cases, one rendering, one model. A one-case improvement on nine cases justifies *keeping the candidate eligible*, not shipping a general claim.

### Recommendation

**Keep bge as an optional candidate and ship-gate it on system performance, not model quality.** The A/B supports eligibility: +1 case, zero regressions, a genuine temporal-correctness gain. It does **not** support a claim that reranking generally improves memory, and it does not soften the agent/patient limitation.

Hardware decides who can use it (§11 gives CPU, partial and full modes); capability decides whether it deserves to exist. Both are now measured separately, which is the only way to keep this decision honest.

Receipt: `eval/s3/ab-rerank-receipt.json`.

**Superseded in part by §13.** That recommendation was conditional on the system gain repeating on cases not designed to produce it. It did not repeat, and the caveat below turned out to be the whole story.

## 13. E-TEMPORAL-SET — the gain does not replicate

Fresh store, new owner ids (`char-101`, `char-102`, `chat-001`), ten cases covering three-deep update chains, corrections, conflicting timestamps, explicit-history queries, temporally unmarked queries, and cross-character ownership. The frozen seven records and the R1–R5 set are not referenced. A = lexical + subword RRF over the case pool; B = bge over the same pool, threshold −3.634 from this store's own held-out pairs (separation 8.888).

`v1` of this experiment (cases `8c6f834a…`) was **withdrawn before any result was used**: it compared record ids against expected values, and its scope probe built a pool that bypassed the filter it was meant to test. `v2` (`43d744e9…`) fixes both.

| Case | Expect | A (pipeline) | B (+bge) | Change |
|---|---|---|---|---|
| T1 current, no temporal words | the Osprey | Vell (distractor) | **the Osprey** | IMPROVED |
| T2 "before the Osprey" | the Wren | **the Wren** | the Kestrel | **REGRESSED** |
| T3 "first ship" | the Kestrel | the Osprey | the Osprey | – |
| T4 correction, current | bosun | **bosun** | quartermaster | **REGRESSED** |
| T5 correction, history | quartermaster | bosun | **quartermaster** | IMPROVED |
| T6 conflicting timestamps | five coppers | five coppers | five coppers | – |
| T7 ownership, own record | the Gull | the Gull | the Gull | – |
| T8 ownership, foreign leak | the Gull | the Osprey (leak) | **the Gull** | IMPROVED |
| T9 non-temporal control | Vell | Vell | Vell | – |
| T10 scope gap (unscored) | abstain | none | none | – |

| Metric | A | B |
|---|---|---|
| Scored cases | 5/9 | **6/9** |
| Current-state | 1/2 | 1/2 |
| Historical-state | 1/3 | 1/3 |
| Ownership | 1/2 | 2/2 |
| Improved | — | T1, T5, T8 |
| **Regressed** | — | **T2, T4** |
| Latency p50 | 0 ms | 31 ms |

### Mechanism probes (diagnostic, not scored)

T4 is the decisive one. bge prefers the **retracted** record at +2.984 against −0.936 for the current one. Three direct probes against the same model:

| Probe | Result |
|---|---|
| "What is Teodor's post?" — plain texts | quartermaster **2.984** / bosun −0.936 → wrong record |
| same, with explicit `(superseded at turn 34, no longer true)` / `(current)` | quartermaster **2.991** / bosun 1.147 → **still wrong** |
| "What ship does Mira sail?" — chain, no temporal words | Osprey **2.651** (correct) |
| same chain with `now` added | Kestrel **2.443** → **worse**, hypothesis falsified |
| "Which ship was Mira sailing before the Osprey?" | Kestrel −0.015, Wren −1.834, Osprey −0.358 → picks the **earliest**, not the previous |

**The mechanism hypothesis I had is wrong, and the correction matters.** I assumed the R3-now win came from the word "now" in the record text. Adding "now" to the chain makes the answer *worse*, and an explicit supersession annotation does not help at all. There is no evidence of a temporal mechanism, however cued. Text-only scoring cannot resolve state, and telling it the state in the text does not fix that.

### What this changes

**The ship gate fails.** §12's condition was "improves meaningful retrieval without unacceptable regressions". Result: three improvements, two regressions, and one regression returns a **retracted fact as current** — precisely the harm that matters most, since it gives a character a wrong model of the present.

**The R3-now gain was one case, not a capability.** Unseen cases split three for and two against. The honest reading is that the earlier A/B measured a favourable case, not a general benefit.

**A design constraint follows, and it is architectural.** A text-only reranker must never be the arbiter of what is current: it can reorder candidates that have *already* been filtered by status and ownership, and nothing more. Status and establishment turn belong out-of-band, in the selection-stage contract — which is the parked architecture workstream's business, not the reranker's.

**Scope gap, measured not assumed.** A world-scope fact ("the fleet's flag") is unreachable under a character filter in both conditions: chat-scope records never enter the character pool. That is a genuine ownership-modelling gap, and it is exactly what the parked workstream exists to fix.

### Disposition after replication

| Capability | Result |
|---|---|
| Ownership discrimination | Mixed — rescues a foreign leak (T8), but prefers a superseded same-subject record (T4) |
| Current vs historical state | **Failed** — T4, T2, T3 all wrong |
| Paraphrase retrieval | Positive (unchanged, R4) |
| Abstention | Unchanged, 3/4 on the earlier set; no new abstain evidence here |
| Agent vs patient | **Failed** (unchanged) |
| Overall selection | 6/9 against 5/9, with **2 regressions** |
| Resource flexibility | Multiple viable modes (§11) |

**Recommendation: do not ship bge.** It remains a registered, replaceable candidate in the selection stage — nothing is deleted — but it fails the product gate, and the failure is in the capability this project cares about: knowing what is true now.

Receipts: `eval/s3/temporal-receipt.json`, `eval/s3/temporal-spec.json`.

## 14. Final disposition — not shipped

**Decision:** do not ship bge, do not fetch another ranking candidate. Artifacts, receipts and the full evaluation history are retained; nothing is deleted. The model stays a registered, replaceable selection stage with one permitted job: reordering candidates that an eligibility filter has already admitted.

Three findings are preserved as the durable result of this workstream:

1. **bge's temporal capability is not established.** The ledger improvement was one favourable case; three improved and two regressed on unseen temporal cases. No mechanism, not even with explicit supersession annotations in the text.
2. **The regressions matter more than the aggregate score.** T4 returning a retracted fact as current is a correctness failure, not a ranking imperfection.
3. **The scope gap is architectural.** World-scope facts never enter a character pool, so no ranking component can recover them.

Reconsideration would require both: status and establishment turn available to the selection stage out-of-band, and the `E-SCOPE-CONTRACT` probe S8 ("a reranker cannot resurrect an ineligible record") passing. Testing another selector before the eligibility rules exist would measure nothing.

The architecture workstream this opens is recorded in
`docs/research/MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md`.

## 8. Claim ceiling

One candidate, nine paired cases, one machine, two record renderings (schema and prose). Nothing here generalises to other rerankers or to larger record sets. Gate B is 6/7 and open. The reranker remains optional infrastructure.

A handcrafted subject/agent/patient rule is **not** an accepted remedy: if one is ever introduced, it needs its own registered experiment with the same case set, and must demonstrate it does not damage the existing wins.
