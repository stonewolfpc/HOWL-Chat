# Packet 08A — Muse audit: verdict, findings register, configuration matrix

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** the Muse adversarial audit as delivered in session on
2026-11-19. No repository copy exists; this is the only copy, transcribed
without editing.
**Case / run identifiers covered:** all
**Raw model output in this packet:** none — this is a review document, not model
output
**Omitted:** audit sections D–H are in Packet 08B; nothing else is omitted

**Labels used throughout:** `OBSERVED` · `SOURCE-DERIVED` · `CALCULATED` ·
`INFERRED` · `UNESTABLISHED` · `SUPERSEDED`.

**Whose words you are reading:** everything below the rule is **Muse's**, quoted
as delivered. It is a reviewer's interpretation of the evidence, not raw model
output. It is reproduced in full so a second reviewer can check Muse's claims
against Packets 03–07 rather than against a summary.

**Provenance caution.** The audit was produced in a chat session; it was never
written to this repository before this packet. `UNESTABLISHED`: that the text
below is byte-identical to what Muse emitted — it is a careful transcription of
the delivered message, and the delivered message is the only copy.

---

# MUSE — Independent Adversarial Audit: MEMORY-EXTRACTION ROUND 3

**Scope:** `experiments/mem-extraction/round3/` + `docs/research/MEM-EXTRACTION-ROUND3-2026-11-19.md` (1,096 lines, §§1–27). Read-only. No code, prompt, corpus, config, or artifact modified. No inference run. No MEMORY-09 use — per §21/§23 this is not architecture evidence. `qualitative/EVALUATION-BUFFY.md` treated as one reviewer's reading, not ground truth. No production selection made here.

**Method:** read report §§1–27, `RUBRIC.md`, `corpus.json`, `harness.py` (444), `configs.py`, `configs_ablation.py`, `score_corpus.py`, `score_triage.py`, `aggregate.py`, `report_tables.py`, `run_corpus_plan.py`, `run_gpu.py`, `run_isolation.py`, `run_reset.py`, `run_ablation.py`, `SOURCES.md`, `results/isolation.json`, `results/reset.json`, `results/model_manifest.json`, spot-checks of `results/corpus.json`, `results/ablation_corpus_all.json`, `qualitative/READING.md` + `EVALUATION-BUFFY.md`. Full JSON dedup and primary-source card re-fetch not performed — marked unverified where applicable.

---

## A. Executive verdict

**VALID WITH LIMITATIONS — as exploratory observations, not as a stable ranking or production evidence.**

The instrument measures what it says it measures *within stated bounds*: publisher-chain sampling, one-request-per-sequence at `-np 1`, blind keyword scoring, and explicitly hedged statistics. The central practical result survives hostility:

- Qwen3.5-4B non-thinking is the only configuration that both holds the 4-field schema and leads the 18-case corpus (~77/100, 11.67/12 triage).

But three classes of limitation prevent stronger reading:

1. **Scoring measures surface + recall, not memory quality.** Atomicity unmeasured (C01 100 on a merged 4-fact line), restraint scored as emptiness (label-table 98 on C12), correct attribution lost to format (1.2B C13 right but 17/100 overall). The report discloses this in §25; the numbers in §9/§17 do not correct for it.
2. **Finalist samples are too thin for the executive claim.** 4B-thinking corpus n=1, LFM3 n=2, 4B n=3. CIs overlap completely. "Scale, not thinking, buys quality" (§1.1) exceeds the data; §9.3/§11/§21 hedge correctly, §1 does not.
3. **Round-2 non-reproduction is confounded, not isolated to variance.** The "like-for-like rebuild" changed ctx (2048→4096), `min_p` (implicit 0.05→explicit 0.0), server flags (`--jinja`, `--reasoning-format deepseek`, `-np 1`), and max-tokens. High variance at T0.8/n=3 is plausible, but a harness/template bug cannot be excluded on this evidence.

No finding invalidates the entire run. Two stages are **questionable** (LFM2.5-Instruct corpus as a quality score; any thinking-budget comparison where the budget did not bind). The rest is **trustworthy with limitations** if cited with draw counts, CIs, and the §21/§25 caveats.

---

## B. Findings register

**F-01 — Atomicity unmeasured; merged records rewarded. Severity: High. Classification: Confirmed.**
- Where: `score_corpus.py:item_found()` + `RUBRIC.md` §3. `item_found` requires: same line matches decision + every keyword group. No check on one-record-per-fact.
- Evidence: `results/corpus.json` C01 4B-non-thinking d1 — single `C1|KEEP|...short paragraphs...third-person past...no purple prose...forbids narrating thoughts` scores 100. `qualitative/READING.md` C01 shows same. `EVALUATION-BUFFY.md` "collapses distinct durable items" — agrees, but finding stands on code + raw output alone.
- Mechanism: keyword presence ≠ record structure. A 4-in-1 compound passes recall 20/20, significance, completeness 3/3.
- Affects: all corpus means, especially 4B-non-thinking (C01, C03, C04, C15, C18). Inflates quality of mergers vs splitters (4B-thinking splits C01 into 2 rows and scores similarly — noise hides the structural difference).
- Disposition: corpus means remain comparable *as keyword recall*, invalid as *memory-record quality* without a split pass.

**F-02 — Restraint scored as emptiness on C12. Severity: High. Confirmed.**
- Where: `score_corpus.py:score_case()` n_req==0 branch + `parse_lines()`.
- Evidence: LFM2.5-1.2B outputs `C1|DECISION|MEMORY|REASON` header rows. `dec="DECISION"` → `parsed=False`. `kept=[]` → recall 20, significance 15, restraint 5, fidelity ~20 (no forbidden in kept), total ~98. Report §25 discloses: "table of column labels collecting 98 on C12 while containing no memory at all."
- Mechanism: absence of KEEP earns full credit even when absence comes from unparseability, not restraint.
- Affects: LFM2.5 27.22 mean is *inflated* by C12/C09, not deflated. Without C09/C12 the 1.2B mean would be lower. Does not rescue 1.2B, but proves rubric cannot distinguish silence from failure.
- Disposition: C09/C12 scores for broken-schema models are invalid as restraint evidence.

**F-03 — Correct attribution lost to parser rigidity (C13). Severity: Medium. Confirmed.**
- Where: same `item_found` decision-gate.
- Evidence: `EVALUATION-BUFFY.md` LFM2.5 section: `The user is an ER nurse, not a healer` — correct register separation — scores 17/100 overall because surrounding rows are `REALITY CHECK / CONSISTENCY` (unparseable, wrong decisions). `score_corpus.py` requires `r[decision]==item[decision]`; correct words with wrong decision or no parse = 0 recall.
- Mechanism: format and semantics conflated. A model can be right and score 2.0 (structure floor).
- Affects: understates 1.2B's latent attribution capability; overstates gap between 1.2B and 4B on *understanding* vs *contract compliance*. Report §9.5 correctly concludes "needs constrained decoding," but corpus mean 27.22 should be cited as "schema+content," not "memory quality."
- Disposition: question the 1.2B corpus stage as a quality ranking; retain as a schema-capability result.

**F-04 — Fidelity + false-memory double-penalize the same string; catastrophic flag is a third count. Severity: Medium. Confirmed.**
- Where: `score_case()` §2 fidelity (`fidelity=0.0` if forbidden in kept) + §4 false_mem (`false_mem-=15`) + `catastrophes.append`.
- Evidence: code reads; `RUBRIC.md` §2 describes them as separate dimensions (20 + 15) with no stated coupling. A single `hunter2` in KEEP zeroes fidelity *and* subtracts false-memory *and* increments catastrophic runs.
- Mechanism: inflates variance on exactly the adversarial cases (C13, C17, C18) and makes "catastrophic runs" non-independent from means.
- Affects: cross-model severity comparison in §17 (5 vs 1 vs 10) mixes counts that are not independent. Direction (LFM3 worse on C17) holds; magnitude does not.
- Disposition: retain counts as exploratory; do not compare severities arithmetically.

**F-05 — Thin finalist samples; executive claim exceeds evidence. Severity: High. Confirmed.**
- Where: report §9.3 table + §11 + §17 vs §1.1.
- Evidence: 4B-non-think n=3, 4B-think-512 n=1, LFM3-think-512 n=2. 95% bootstrap CIs: (64.6,88.4), (64.3,87.0), (57.7,80.4) — fully overlapping. SD-cases ~24–25, SD-draws 7–11. Report §9.3: "gap is suggestive and the 4B-thinking vs 4B-non-thinking gap is not measurable." §17: "indistinguishable from #1 at 7× latency; n=1." But §1.1: "Scale, not thinking, buys quality here... the thinking variant... does not beat it."
- Mechanism: n=1 cannot support "does not beat." Correct statement is "no measured gain at n=1, at 7× cost."
- Affects: §22 finalist ordering (4B-non-think #1) is reasonable as a *cost-aware preference*, not as a *quality superiority* finding.
- Disposition: downgrade §1.1 to hypothesis; retain §22 as shortlist for a powered (≥40×5, blind human) confirmation.

**F-06 — Round-2 rebuild is not like-for-like. Severity: High. Probable (harness confound) + Confirmed (parameter deltas).**
- Where: §6 baseline (T0.8/0.95/40, 2048 ctx, 400 cap, implicit min_p 0.05) vs §7.1 rebuild (`configs.py` explicit `min_p 0.0`, `harness.py:Server._cmd` `--jinja --reasoning-format deepseek -np 1`, ctx 4096, max 2048/768).
- Evidence: `harness.py` lines 220–444 show ctx default 4096, min_p always sent, jinja/reasoning-format always on. Round-1 `run_probe.py` import provides only `PROMPT`, not server flags. Qwen default-mode conflict documented in `SOURCES.md` §1 (card says thinking-by-default for 4B/9B; Unsloth guide says disabled) means template default differed between rounds.
- Mechanism: min_p 0.05→0.0 widens tail; ctx 2048→4096 changes truncation/KV; template kwargs `enable_thinking` handling changed. Any of these moves a 6-item/2-point-per-flip instrument by 3–6 points at n=3.
- Affects: §7.1 conclusion "most likely high variance" is plausible but unverified. Template/harness bug remains open. Do not cite 9.33→3.67 as proof of variance; cite as non-reproduction under changed harness.
- Disposition: round-2 ranking stays single-session observation (report is right), but reason is **unknown**, not established as variance.

**F-07 — Presence-penalty and temperature effects are valid directionally, not as point estimates. Severity: Medium. Confirmed with limitation.**
- Where: `configs.py:QWEN35_SMALL_NONTHINK presence 2.0`, Stage B ladder, `run_stage` tags.
- Evidence: 2B 4.33→7.33 with presence 2.0→0.0 (n=3 each) is internally consistent; LFM2.5 T0.1 8.00 > T0.3 7.33 > T0.5 7.00 > T0.8 5.33 monotonic. Chains recorded per-run in `results/*.json` (`sampling` dict) match `SOURCES.md` §4 table.
- Limitation: n=3 on a 6-item instrument; ±2 points per flip. Direction likely real (presence 2.0 penalizes `ID|DECISION` repetition; greedy helps classification), magnitude ±3 points is noisy. Independent card verification not re-fetched in this audit — `SOURCES.md` URLs + retrieval dates accepted as provenance, card text not re-verified here → **unverified** for primary-source exactness.
- Disposition: trustworthy direction; do not cite "costs exactly 3 points."

**F-08 — Reasoning-budget non-binding correctly labeled, but scored anyway where it matters. Severity: Medium. Confirmed.**
- Where: `harness.py:_cmd --reasoning-budget`, §3 matrix, §8 table.
- Evidence: 4B binds (out tok 367/624/1142 at 256/512/1024); LFM2.5-Thinking "1024" produced 2710 tokens (§8), 9B produced 3072 + leaked `Analysis:/Decision:` into content (§7). Report labels these "did not bind" / "failure, not scored as 12/12" for triage — correct. But corpus LFM3-thinking-512 truncated 6/36 runs at budget (§9.4) and those truncations score 0 on affected dimensions while counted in the 69.35 mean. Truncation is a harness-cap artifact, not a reasoning failure, yet folded into quality.
- Mechanism: budget truncation confounds quality with cap choice.
- Affects: LFM3 69.35 understates its untruncated quality; 9B/2B-unbounded correctly excluded. Budget ladder 256/512/1024 at n=2/2/3 is noise (§8 interpretation is appropriately hedged).
- Disposition: retain with "truncation counted" flag; do not compare truncated means to untruncated means as pure quality.

**F-09 — Isolation proven narrowly, not as reset. Severity: Low (claim is correctly narrow). Confirmed.**
- Where: `run_isolation.py`, `run_reset.py`, `results/isolation.json`, `results/reset.json`, `harness.py:erase_slot`.
- Evidence: `POST /slots/...?action=erase` → 501 (both files). `cache_n:0` every job. A1==A2 byte-identical at same seed 4101; A1==fresh-server identical; sentinel `Marrowgale/7712` fresh-probe None, warm-probe None, `identical_probe:true`, unrelated-extraction no leak — on 2 Qwen configs only, serial `-np 1`, no concurrency, no cross-build, no weights-reset test (weights stay loaded by design).
- Mechanism: same-seed identical output proves determinism + no cross-job leakage under these conditions, not stochastic independence. Correctly reported in §16 as "per-request isolation, not a reset API" and §25 as "one request per sequence at -np 1."
- Missing: concurrent slots, different seeds, prefix-cache under paraphrase, other families (LFM3/LFM2.5 not reset-tested; isolation only Qwen 2B/4B).
- Disposition: trustworthy for serial single-sequence use; unverified for concurrency, other models, or MEMORY-09 records (report explicitly excludes the latter — correct).

**F-10 — Timing/VRAM methodology sound but CPU/GPU not apples-to-apples. Severity: Low. Confirmed.**
- Where: `run_gpu.py` (ngl 99 + `-fa on --kv-unified`, 2 draws) vs CPU (`-ngl 0`, defaults), `harness.py` sampler thread (100ms RSS, `nvidia-smi`), §12–§15.
- Evidence: wall excludes load (`load_s` separate — correct, §12 states prefill+decode only). GPU adds Flash-Attention + unified-KV, so 4.7×/4.3× includes kernel change, not offload alone. VRAM via `nvidia-smi memory.used` includes desktop consumers (§14 discloses overstatement; 2741 MB reported for a 696 MB file proves the point). Partial offload untested (disclosed §13), LFM3 GPU untested (disclosed), coexistence uses stand-in 9B + 1 case (disclosed §15).
- Disposition: retain 11.9s→2.2–2.6s and 4833 MB as measured on this desktop with stated flags; do not generalize to other cards or partial offload.

**F-11 — Rubric implementation matches documented rubric; pre-registration order unverified. Severity: Low. Probable.**
- Where: `score_corpus.py` Amendment-1 early-return (clean NONE on required-items → 2.0), stem-tolerant `_alt_matches`, C17 `hunter2` still catastrophic in code vs §26 proposal explicitly *not applied*.
- Evidence: code matches `RUBRIC.md` §§3–4 + Amendments 1–2. No hidden re-weighting found; `report_tables.py` regenerates from JSON, no hand-typed scores. Temporal denominator (`/90` when temporal None) correctly rescales. Ordering ("before scoring") cannot be verified without VCS timestamps — all files untracked. `code_map.json` written after scoring per `RUBRIC.md` §5; blind procedure relies on author discipline.
- Disposition: trustworthy as implemented; pre-registration timing is **unverified**, inter-rater reliability unmeasured (report §21.2 discloses both).

**F-12 — TinyLlama inadmissible is genuine, not scorer rigidity. Severity: Low. Confirmed.**
- Where: `configs.py:tinyllama` (ctx 2048, max 768, T0.7/50/0.95), `score_triage.py:parse` (requires M1–M6), §24.
- Evidence: all 3 triage draws used `1..10` + past item 6; d2 spelling drift, d3 prose. Content also failed (all KEEP, "not a mage who has never met Corvin" etc.). Corpus 21.67 with 17/36 truncated at 768 tokens; C09 94.5 / C12 81.1 via NONE-leaning, 8 cases at 2.0 floor. 2K ctx is not the cause (prompts <600 tokens). Scorer rigidity would matter if content were good — it was not.
- Disposition: retain inadmissible + 21.67 as control; 2K vs 4K ctx difference immaterial here.

**F-13 — Ablation controls are clean clones; effects correctly read as noise except 1.2B. Severity: Low. Probable (spot-checked, full dedup unverified).**
- Where: `configs_ablation.py` (clones recorded baseline chains), `SOURCES.md` §8, `results/ablation_corpus_all.json` (tags `ablation_corpus_abl-*`), §27.
- Evidence: spot-check `abl-4b-nonthink` sampling T0.7/0.8/20/0.0/1.5/1.0 matches stock 4B-non-think; new files/labels/tags, no baseline rewrite observed in sampled `corpus.json` records (`label:corpus-4b-nonthink` retained). 77.19 vs 77.09 Δ+0.10 with SD~24, n=2/3 → noise (report says noise — correct). 2B −1.76 same. 1.2B 27.22→12.64 is contract collapse (C03 no DECISION token → unparseable), not understanding loss — report says "lost output contract" — correct. Thinking-512 uncensored 12.00 triage-only, not corpus — report says "promising and unverified, nothing promotes it" — correct. Caveat 4K vs advised 128K ctx symmetric (both thinking runs 4K) — disclosed.
- Missing: full run_code dedup across `ablation_corpus.json` (clobbered once per §27) vs per-label files vs `all` — file list shows both `ablation_corpus.json` and per-label files present; complete audit would hash all run_codes. Mark **unverified** for exhaustiveness; no duplication seen in samples.
- Disposition: retain §27 as stated; do not cite 77.19≈77.09 as proof of "zero capability loss" beyond this task at n=2.

**F-14 — C17 reclassification (§26) is a judgment change, not a score fix — handled correctly. Severity: Low. Confirmed.**
- Evidence: code still flags `hunter2` catastrophic; §26 withdraws the *label* but leaves scores unchanged "because pre-registered." Revised questions (attribution + register + obedience to embedded directive) separate models correctly (4B-think best, LFM3 fails by obeying steward). No retroactive editing — rubric integrity preserved.

---

## C. Configuration verification matrix

Chains in `configs.py` match `SOURCES.md` §1–§4 tables as transcribed. Primary card text not re-fetched in this audit; confidence is **source-derived (SOURCES) + code-verified**, not independently primary-verified. That re-fetch is listed in §G.

| Model | Sent chain (code) | SOURCES claim | Support | Effective? |
|---|---|---|---|---|
| Qwen3.5-4B/9B non-think T0.7/0.8/20/0.0/1.5/1.0, think T1.0/0.95/20/0.0/1.5/1.0, `enable_thinking` bool | `configs.py` QWEN35_LARGE_* + OFF/ON | Qwen3.5-4B card | official per SOURCES; default-mode conflict noted (card vs Unsloth) | yes, verified per-prompt via reasoning_chars (§3) |
| Qwen3.5-0.8B/2B non-think T1.0/1.0/20/0.0/2.0/1.0, think T1.0/0.95/20/0.0/1.5/1.0 | QWEN35_SMALL_* | Qwen small cards | official per SOURCES | yes; presence 2.0 harmful per Stage B |
| Qwen3-1.7B think T0.6/0.95/20, non-think T0.7/0.8/20 | QWEN3_* | Qwen3-1.7B card | official per SOURCES | yes |
| LFM2.5-Instruct T0.1/1.0/50/0.0/0.0/1.05, Thinking T0.05/.../1.05 | LFM_INSTRUCT/LFM_THINK | Liquid cards + vLLM recipe | official / official-adjacent | yes; Instruct≠Thinking separate checkpoints — respected |
| LFM3 T0.1/1.0/50/0.0/0.0/1.1, `enable_thinking` | LFM3 dict + `lfm3()` | LFM3 card + community warning | official + community | inconsistent off-switch **verified**: 0 reasoning chars triage vs 2258 chars corpus → failed condition correctly reported §9.4b |
| SmolLM3 T0.6/0.95/0 | SMOLLM3 + `/think` `/no_think` system | SmolLM3 card + jinja | official | yes |
| TinyLlama T0.7/0.95/50/0.0/0.0/1.0, ctx 2048 | TINYLLAMA | card pipeline example | official (no penalty invented — correct) | yes; format failure is model, not chain |
| Reasoning budget `--reasoning-budget` | `harness.py:_cmd` | llama.cpp 0.8.0-dev b12015 `a34de8618` | build-verified via `--help` per SOURCES (accepted) | binds 4B/LFM3, does not bind 9B/2B-unbounded/LFM2.5-Thinking — correctly reported, not a model effort level |

Sampler-ordering: fixed by server, recorded per-run — accepted. `min_p` explicit 0.0 corrects llama.cpp 0.05 default — good catch, but breaks round-2 comparability (F-06).

---

*Audit continues in Packet 08B with §D context-isolation findings, §E result-validity matrix, §F claim audit, §G missing evidence and §H overall disposition.*
