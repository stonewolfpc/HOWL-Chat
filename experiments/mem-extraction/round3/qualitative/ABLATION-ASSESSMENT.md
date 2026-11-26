# Abliterated / uncensored variants — assessment

**Date:** 2026-11-19 · **Runs:** `results/ablation_triage.json`,
`results/ablation_corpus_all.json`, `results/control_corpus_2b.json`
**Docs:** `qualitative/ablation/` · **Config:** `configs_ablation.py`

Four variants, each run on the same prompts and the same sampler chain that the
matching baseline run used, so a difference is attributable to the weights:

| Variant | Base | File |
|---|---|---|
| Qwen3.5-4B Uncensored (HauhauCS Aggressive) | 4B non-thinking | Q4_K_M, 2.71 GB |
| Qwen3.5-2B Uncensored (HauhauCS Aggressive) | 2B non-thinking | Q4_K_M, 1.27 GB |
| LFM2.5-1.2B-Instruct abliterated (Huihui) | 1.2B T0.1 | Q4_K_M, 0.73 GB |
| LFM3-1.7B Heretic abliterated (Vortecks) | 1.7B thinking | Q4_K_M, 1.10 GB |

Triage was 3 draws for all four. Corpus was 2 draws for the three non-thinking
configurations; the thinking variants were triage-only because their corpus cost
is 96 and 34 minutes respectively at this hardware. A stock Qwen3.5-2B corpus
control was run so the 2B comparison has a baseline at all — none existed before.

## Results

| Condition | Six-item (3 draws) | 18-case corpus | Δ corpus | Catastrophic | Wall | Peak RSS |
|---|---:|---:|---:|---:|---:|---:|
| Qwen3.5-4B non-thinking, stock | 11.67 | **77.09** (n=3) | — | 5 | 11.9 s | 3,438 MB |
| Qwen3.5-4B HauhauCS Uncensored | 11.33 | **77.19** (n=2) | +0.10 | 4 | 13.5 s | 4,261 MB |
| Qwen3.5-4B thinking, budget 512, stock | 10.5 | not run | — | — | — | — |
| Qwen3.5-4B thinking, budget 512, uncensored | **12.00** | not run | — | — | 86.4 s | — |
| Qwen3.5-2B non-thinking, stock | 4.33 | 54.83 (n=2) | — | 3 | 6.0 s | — |
| Qwen3.5-2B HauhauCS Uncensored | 5.00 | 53.07 (n=2) | −1.76 | 2 | 5.8 s | 2,037 MB |
| LFM2.5-1.2B-Instruct T0.1, stock | 8.00 | 27.22 (n=3) | — | 1 | 2.6 s | 1,098 MB |
| LFM2.5-1.2B-Instruct abliterated | 6.67 | **12.64** (n=2) | **−14.58** | 0 | 2.8 s | 1,166 MB |
| LFM3-1.7B thinking 512, stock | 11.33 @1024 | 69.35 (n=2) | — | 10 | 55.9 s | — |
| LFM3-1.7B Heretic abliterated, 512 | 9.67 | not run | — | — | 31 s | — |

## What the numbers do and do not say

**At 4B, abliteration is close to free on this task.** 77.19 vs 77.09 is not a
difference; with two and three draws it is noise. The claim on the model card —
"zero capability loss" — holds for extraction at this sample size. It costs
about 0.8 GB of RSS, which is quantisation/metadata, not behaviour.

**At 2B, also free** (−1.8, within a draw's worth of variance), and it removes
one catastrophic run (3 → 2).

**At 1.2B it is not free.** 27.22 → 12.64 is a collapse of the output contract,
not of understanding. Read the C03 answer:

```
C1
MEMORY: Sable is a fire mage trained at Ossary.
REASON: Core identity and profession are established.
...
```

Four clean, atomic, correctly separated facts — better structure than the stock
model produced on that case. But the `DECISION` field is missing or contains
`Halden` instead of `KEEP`, so nothing parses and the rubric floors it at 2. The
same pattern on C05 and C09. The pattern is content preserved, decision token
lost. That is a fixable contract failure (a grammar would recover it), not a
capability regression, but as shipped it makes the variant unusable.

## The interesting semantic shifts at 4B

The aggregate is unchanged; the *profile* is not.

**C13 improves, and for the right reason.** Stock 4B wrote "Sable … acts as an
ER nurse", fusing the user's real job into the character. The uncensored variant
writes: "The user claims to be an ER nurse … but this has not been confirmed".
Wrong in a new way — it hedges a fact the user stated plainly about themselves —
but it keeps the registers apart, which is the error that mattered. 37 → 55.

**C17 gets worse under the revised policy.** Stock non-thinking answered `NONE`.
The uncensored variant splits the case into four rows and stores
`C2 | KEEP | User's password is hunter2.` Flat, no register, attributed to the
user rather than the steward. That is precisely the LFM3 failure I described in
`SECRETS-AND-MEMORY.md`: fiction promoted to a user record. Removing refusals
made it more willing to write down a claim that should have been framed.

**C10 and C14 improve** (75 → 92, 37 → 58); **C01–C04 fall** slightly (100 → 92,
92 → 79, 95 → 82, 100 → 88). The merge-everything habit is unchanged: C01 is
still one line carrying four rules.

## Thinking mode

The uncensored 4B thinking variant scored 12/12 on all three triage draws, above
the stock model's 10.5 at the same 512 budget and its 11.67 unbudgeted. Triage is
six items; I would not promote a configuration on that alone, and it was not run
on the corpus, so the honest statement is: *promising, unverified*. Given that
the stock thinking variant's corpus failures were invented reconciliation (C07)
and a stored credential (C17), the corpus is exactly where this variant needs
testing before any conclusion.

## Caveats

- The HauhauCS card says to keep at least 128K context to preserve thinking
  quality. Both the stock and uncensored thinking runs used 4K. The comparison is
  symmetric, but both may understate thinking capability; nothing here tests the
  card's claim.
- Two draws per corpus condition; one draw would not have been enough for
  anything, two is not enough for confidence intervals.
- These are community re-quantisations and community abliterations, labelled as
  such. No first-party uncensored checkpoint exists for either family.
- Abliteration removes refusal behaviour. This corpus contains no request a stock
  model would refuse, so the variants are being tested on a task where their only
  difference is side effects on formatting and judgement — which is what the
  results show.

## Bearing on the choice

None of the four unseats Qwen3.5-4B non-thinking as the reference. The uncensored
4B is a legitimate drop-in if refusal behaviour is a problem in production, at
equal quality, slightly higher memory, and one new failure mode on secrets.
The abliterated 1.2B should not be used: it loses the decision token that the
entire pipeline depends on.
