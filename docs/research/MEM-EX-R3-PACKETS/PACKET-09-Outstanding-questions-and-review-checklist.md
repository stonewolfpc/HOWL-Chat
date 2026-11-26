# Packet 09 — Outstanding questions and review checklist

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** synthesis of all packets
**Case / run identifiers covered:** all
**Raw model output in this packet:** n/a — no model output
**Omitted:** nothing omitted; this packet is argument, not evidence

**Labels used throughout:** `OBSERVED` directly present in a stored response or
execution record · `SOURCE-DERIVED` stated in a source or configuration file ·
`CALCULATED` derived by explicit arithmetic · `INFERRED` a reading of the
evidence · `UNESTABLISHED` not established by the available evidence ·
`SUPERSEDED` an earlier claim replaced by later evidence.

**Whose words you are reading:** fenced blocks are the model's own output
(`OBSERVED`, verbatim). Prompt text is `SOURCE-DERIVED` from `corpus.json`.
Scores are `CALCULATED`. Any sentence outside a fence that is not a definition is
Buffy's interpretation unless it is attributed to Muse or to the project owner.

**Caution.** These are exploratory measurements. Nothing here selects a model,
establishes memory quality, or supplies MEMORY-09 evidence. Sample sizes are
1–3 draws per case.

---

## A. What is actually established

`OBSERVED` and `CALCULATED`, stated at the strength the evidence supports:

1. On this machine and corpus, Qwen3.5-4B non-thinking held the four-field
   contract in every run and had the highest measured semantic score. It is a
   reasonable shortlist candidate.
2. LFM3-1.7B thinking was the only small model with a credible semantic profile;
   it lost 6 of 36 runs to the 512-token cap.
3. LFM2.5-1.2B-Instruct did not produce a contract-valid row in any of its 90
   corpus runs. That is a schema-capability result, not a quality score.
4. Thinking mode produced no measured quality gain at 4B and cost roughly 14×
   the wall time at the tested budget. n=1.
5. Serial per-request isolation was observed at `-np 1`; no reset API exists.

## B. What is not established

`UNESTABLISHED`, and should be treated as open by any reviewer:

1. Whether any configuration is *better* than another in general. CIs overlap,
   draws are 1–3, and one corpus with one author cannot support a ranking.
2. Whether round 2's numbers failed to reproduce because of variance, a template
   change, or a harness bug.
3. When the rubric and corpus were written relative to the runs. No pre-run hash
   exists.
4. Whether v2's semantic composite tracks human judgment. Its tests prove rule
   conformance, nothing more.
5. Any ownership behaviour beyond the declarative cases: the detector cannot
   infer scope from free text.
6. Anything about concurrency, other builds, other quantizations, or MEMORY-09
   records.

## C. Questions for independent review

These are the questions where a second reader's judgment is worth more than
another computation.

1. **C01.** Is "the user prefers short paragraphs, third-person past tense,
   avoids purple prose, and forbids narrating thoughts" one memory or four?
   Does the answer change what a store should do with it?
2. **C12.** Three column labels and the word `MEMORY` in the third column. Is
   that "restraint" in any sense a reader would accept? If not, what should a
   scorer do with a model that answers this way consistently?
3. **C13.** The character is a healer; the user is an ER nurse. How much should a
   single fused sentence cost, given that a human reader can still recover both
   facts from it?
4. **C17.** The thinking model stores `hunter2` while saying it is not durable
   canon. Is that a memory failure, a policy failure, or correct fidelity to an
   in-world disclosure? §26 argues the last; the v1 scorer says the first.
5. **C18.** `exactly double the toll (not triple)` — how should a scorer treat a
   forbidden word that appears only inside a negation? v1 counted it; v2 does
   not.
6. **F-05.** With n=1–3, is *any* ordering in §17 defensible, or should the
   report present cost and contract separately and refuse to order by quality?
7. **F-06.** Which single configuration difference (context, `min_p`, `--jinja`,
   or max tokens) is the most likely cause of the round-2 non-reproduction, and
   what is the cheapest experiment that would distinguish them?
8. **Truncation.** LFM3 lost 6 runs to a 512-token cap. Should a truncated run
   contribute a partial score, be excluded, or be re-run before comparison?
9. **Ownership.** The adversarial cases in
   `instrument_v2/cases_adversarial_v2.json` define belief, knowledge, disclosure
   and scenario scope. Are those the right distinctions for Howl, and is
   declarative detection sufficient for a first measurement?
10. **Instrument v2.** Are the two composites weighted defensibly? Is
    `production_usable` (contract ≥ 80, no critical event) the right gate, or
    does it hide outputs that a human would accept?

## D. Decisions that belong to the project owner, not to a scorer

1. Whether to apply the Appendix C errata to the report in place or publish them
   as an addendum. Recommended: addendum.
2. Whether retaining a disclosed credential should remain a critical event
   (current v1 and v2 rule) or become a policy question scored on attribution and
   register (§26's argument).
3. Whether a three-field row (`ID | DECISION | MEMORY`, no reason) counts as
   contract-valid. v2 says no; this materially affects the 4B's C14 draws.
4. Whether to fund a powered confirmation (≥40 cases × 5 draws, two independent
   scorers) before any model decision.
5. Whether v2 should replace v1 for future rounds, run beside it, or be revised
   first.

## E. Reading checklist

- [ ] Packet 00 — know which file is raw and which is interpretation.
- [ ] Packet 01 — know what was actually sent and measured.
- [ ] Packet 02 — know the instrument's rules and its known defects.
- [ ] Packet 03 — see the six-item prompt saturate.
- [ ] Packet 04 — read at least C01, C07, C13, C17 with the reference in view.
- [ ] Packet 05 — check the five failure cases against the raw text.
- [ ] Packet 06 — check that the isolation claim is narrower than "flush works".
- [ ] Packet 07 — check that CPU and GPU flags differ.
- [ ] Packet 08 — read Muse's findings before forming a verdict.
- [ ] Packet 09 — answer the ten questions, then the five owner decisions.


---
