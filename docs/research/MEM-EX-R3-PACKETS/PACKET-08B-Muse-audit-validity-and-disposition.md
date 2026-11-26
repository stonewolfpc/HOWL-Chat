# Packet 08B — Muse audit: isolation, validity matrices, disposition

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** the Muse adversarial audit as delivered in session on
2026-11-19 (sections D–H). No repository copy exists; this is the only copy.
**Case / run identifiers covered:** all
**Raw model output in this packet:** none — this is a review document
**Omitted:** audit sections A–C are in Packet 08A; nothing else is omitted

**Labels used throughout:** `OBSERVED` · `SOURCE-DERIVED` · `CALCULATED` ·
`INFERRED` · `UNESTABLISHED` · `SUPERSEDED`.

**Whose words you are reading:** everything below the first rule is **Muse's**,
quoted as delivered, until the section marked "Buffy — finding register and
disposition". That last section is the correcting author's response, kept
separate so a reviewer can disagree with either independently.

---

## D. Context-isolation findings

Reset: 501, `erase_supported:false`, `cache_n:0` always, byte-identical A1/A2/fresh at fixed seed. Isolation: fresh/warm `None`, identical, no `Marrowgale/7712` leak into C01 extraction, on Qwen 2B+4B only. Proven: serial per-request isolation at `-np 1` on this build. Not proven: explicit reset, concurrency, cross-seed leakage, other families, other builds, MEMORY-09 records. Report's narrow claim is accurate; any "context flushing works" citation would be false.

---

## E. Result-validity matrix

| Stage | Verdict | Why |
|---|---|---|
| Triage 4B / LFM3 / Qwen3 / SmolLM3 | Trustworthy with limitations | n=3, high variance noted, protocol failures correctly excluded |
| Triage 2B/1.2B tuning ladder | Trustworthy with limitations | direction over magnitude; confounded vs round 2 |
| Corpus 4B-non-think (n=3) | Trustworthy with limitations | best measured, but atomicity unmeasured, 5 catastrophics incl. attribution/numeral |
| Corpus 4B-think-512 (n=1) | Questionable as ranking | single draw + worst-severity C17; "not better at 7× cost" only |
| Corpus LFM3-think-512 (n=2) | Trustworthy with limitations | 6 truncations are cap artifacts; C02/C12 inversions real |
| Corpus LFM2.5-T0.1/T0.8 | Questionable as quality, Trustworthy as schema probe | 39/54 unparseable; means are contract scores |
| LFM3-non-think corpus | Invalid (as condition) | off-switch ignored; single failed run correctly reported, not a mean |
| GPU 4B/1.2B | Trustworthy with limitations | flags differ from CPU; VRAM overstated by desktop use (disclosed) |
| Reset / Isolation | Trustworthy with limitations | narrow serial claim only |
| Coexistence | Trustworthy with limitations | 1 case, stand-in 9B, VRAM arithmetic only |
| TinyLlama | Trustworthy (as control) | genuine contract failure |
| Ablation 4B/2B | Trustworthy with limitations | n=2, noise correctly read |
| Ablation 1.2B | Trustworthy (as contract collapse) | not a capability delta |

---

## F. Claim audit

1. Scoring heuristic — **supported with disclosed defects**. §10/§21/§25 own the disagreements; audit confirms mechanism (F-01–F-04).
2. Thin-sample ranking — **overstated in §1, hedged in §9/§11/§17**. CIs overlap; "does not beat" needs n=1 qualifier.
3. Round-2 non-repro = environmental — **unsupported as stated**. Variance plausible, harness confound probable (F-06). Correct action (treat round 2 as single-session) retained, reason changed to unknown.
4. Reset narrow — **supported**. 501 + byte-identical + cache_n 0 = per-request isolation, not flush, not weights reset. §16 wording is precise.
5. Overwrite incident — **handled, spot-checked clean; full dedup unverified**. Per-label tags distinct; `all` sample tags match labels; stock-2B control is new data (`control_corpus_2b.json` present). Exhaustive run_code audit still open (§G).
6. Timing — **supported with flag asymmetry noted**. Wall excludes load; cold-start not in wall; GPU flags add FA/unified-KV; VRAM overstatement disclosed.
7. Rubric integrity — **supported as implemented; timing unverified**. Amendments 1–2 in code; C17 proposal correctly unapplied; pre-registration order needs VCS proof.
8. TinyLlama — **supported**. Genuine format failure + poor content; scorer not at fault.
9. Abliterated — **supported**. Clone chains verified in sample; 77.19 vs 77.09 noise correctly read; 1.2B collapse is contract, not understanding.

---

## G. Missing evidence (needed for VALID)

- Full run_code uniqueness + generation-parameter audit across `corpus.json`, `ablation_corpus.json`, per-label files, `ablation_corpus_all.json`.
- Prompt-token accounting per model tokenizer (actual vs recorded `prompt_tokens`) for C10-long + TinyLlama-2K edge.
- Primary-source re-fetch of all sampler/effort claims (Qwen, Liquid, SmolLM, Granite cards at pinned revisions) + llama.cpp b12015 `--reasoning-budget` semantics doc.
- Per-model chat-template rendering dump (jinja) for Qwen3.5 default-mode conflict and LFM3 off-switch inconsistency.
- Truncation-free LFM3 corpus re-run (higher cap) to separate quality from budget artifact.
- Powered confirmation: ≥40 cases × 5 draws, blind human scoring, inter-rater reliability, seeded variance decomposition.
- Concurrency/isolation matrix (parallel slots, cross-seeds, other families) if isolation is ever cited beyond serial use.

---

## H. Overall disposition

**Retain as exploratory observations with the report's own caveats; exclude from any MEMORY-09 or production decision; do not re-rank on these means.**

Specifically: retain 4B-non-thinking and LFM3-thinking as confirmation finalists (§22) — not as selected models. Retain LFM2.5 corpus means only as schema-failure evidence. Treat the 4B vs 4B-thinking tie as unresolved, the 4B/LFM3 gap as suggestive, round-2 numbers as non-reproduced for unknown reasons, and all safety-relevant behaviors (C13 attribution, C17 directive-obedience, numeral reversal) as testable hypotheses for a harder adversarial set — not as closed measurements.

No redesign prescribed; no rerun initiated.

---

# Buffy — finding register and disposition after the correction pass

**This section is not Muse's.** It is the correcting author's response, written
2026-11-19, so a reviewer can compare the audit with what was actually changed
and disagree with either. The corrected instrument is
`experiments/mem-extraction/round3/instrument_v2/` (MEMX-EVAL 2.0.0); the
correction report is
`docs/research/MEM-EXTRACTION-ROUND3-INSTRUMENT-CORRECTIONS-2026-11-19.md`.

| Finding | Muse severity | Disposition | Evidence in v2 |
|---|---|---|---|
| F-01 atomicity | High | Corrected | `atomicity.py`; four-facts-in-one-row recall 30 → 15 |
| F-02 restraint as emptiness | High | Corrected | `outcomes.py`; C12 label table 97.8 → `MALFORMED`, 0 |
| F-03 semantics vs format | Medium | Corrected | `matching.py`; recall kept, contract separated, `production_usable` false |
| F-04 double penalty | Medium | Corrected | one event, one primary dimension |
| F-05 thin samples | High | Tooling corrected; report errata recorded, not rewritten | `reporting.compare` returns `UNESTABLISHED` |
| F-06 round-2 rebuild | High | Documented; cause left **unknown** | correction report §4 |
| F-07 sampler point estimates | Medium | Tooling corrected | `reporting.effect` |
| F-08 truncation | Medium | Corrected | `quality_eligible`, eligible/all-runs means |
| F-09 isolation | Low | Narrowed, preserved | `reporting.ISOLATION_ALLOWED` |
| F-10 timing | Low | Qualified, preserved | `reporting.TIMING_ALLOWED` |
| F-11 pre-registration | Low | **UNESTABLISHED**, stated as such | `reporting.PREREGISTRATION` |
| F-12 TinyLlama | Low | Preserved, not rescored | triage scorer untouched |
| F-13 ablation controls | Low | Preserved; dedup audit still open | correction report §7 |
| F-14 C17 judgment | Low | Preserved; event marked `policy_sensitive` | `rubric_v2.py` |

**Two defects Muse did not name, found in the correction pass:**

- **D-15.** A forbidden word inside a negation was counted as stored.
  `C18`'s `"exactly double the toll (not triple)"` was recorded by v1 as a
  catastrophic numeric reversal. v2 does not count it; 16 of 28 previously
  counted stored-forbidden incidents disappear on the preserved corpus.
- **D-16.** One credential disclosure produced two events in the first v2 build.
  Events are now aggregated per incident.

**One audit sentence corrected against the stored record:** the 1.2B's C13 case
mean is 17.3, but draw 1 is 47.8 and draws 2–3 are 2.0. "Scores 17" is true of
the case mean, not of a run.
