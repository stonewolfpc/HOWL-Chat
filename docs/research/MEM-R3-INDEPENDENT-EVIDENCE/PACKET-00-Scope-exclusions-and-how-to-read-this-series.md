# Packet 00 — Scope, exclusions, and how to read this series

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19 · **For:** independent review
**Runs covered:** 478 stored generations (379 corpus-prompt,
99 six-item-prompt)

## 1. What this series is

A source-preserving view of the Round 3 model evidence, assembled above and
independently of every scoring instrument. It answers four questions per run:
what was sent, to which model, what came back, and where the record lives.

It does not score, classify, rank or interpret. No file in this series contains
a model recommendation.

## 2. What is in the primary reading set

For every one of the 478 stored generations:

1. the prompt, or an explicit statement that the exact request was not recorded
   (see §4 and Packet 01);
2. the stored model identifier and model file, exactly as recorded;
3. the stored response, byte-identical to the JSON value, including empty,
   malformed, repetitive and truncated output;
4. `run_code`, case, draw, batch tag and the exact source file with its array
   index.

Corpus runs are grouped by case (Packets 02A–02R) because the case conversation
is the variable a reader needs in view. Six-item runs are grouped by configuration
(Packets 03A–03M) because the prompt is identical across them.

## 3. What is deliberately absent

- **Scores of any kind**, including the v1 rubric and the separate v2
  recalculation.
- **Scorer statuses and event labels**, and any classification of a response as
  valid, malformed, truncated or catastrophic.
- **Rankings or comparative quality claims.**
- **Expected answers.** The corpus's `required`/`forbidden` blocks are the
  grading key. The model never saw them, so they are not model evidence and are
  not reproduced. Their location is given in §6 so a reviewer can obtain them
  independently.
- **Commentary by Buffy or by Muse**, including the earlier packet series'
  interpretation.
- **Sampler and execution metadata** in the packets. It is indexed separately,
  by `run_code`, in `REFERENCE-EXECUTION-METADATA.md`.
- **Any retrospective correction** to a prompt or response. Nothing here has
  been rewritten.

## 4. The prompt is the one incomplete component

No per-run request body was ever written to disk. The result files store the
response, token counts and a `prompt_tokens` figure, but not the prompt text or a
prompt hash.

What exists is the code path that built each request and the corpus text it read.
Packet 01 reproduces the extraction instruction, the exact composition
expression, and the six-item prompt, and states precisely which parts are
`SOURCE-DERIVED` reconstruction. No reconstruction is presented as a recorded
request.

**Consequence for a reviewer:** the response side of every run is
`OBSERVED`; the prompt side is `SOURCE-DERIVED` for every run. Prompt bytes
identical to what the server received are `UNESTABLISHED` throughout this series.

## 5. Evidence labels

| Label | Meaning in this series |
|---|---|
| `OBSERVED` | Present in a stored record: a response string, a recorded identifier. |
| `SOURCE-DERIVED` | Read from a source file — the instruction text, the composition expression, the case conversation. |
| `UNESTABLISHED` | Not established by the available evidence; do not infer it. |

`INFERRED`, `CALCULATED` and `SUPERSEDED` are not used: this series contains no
inference, no arithmetic over runs, and no replacement of earlier claims.

## 6. Where everything else lives

| For | File | Key |
|---|---|---|
| Every run, its packet and its source | `INDEX.md` | `run_code` |
| Sampler chain, context, token counts, timings, memory | `REFERENCE-EXECUTION-METADATA.md` | `run_code`, `label` |
| Model file, size, SHA-256, repository string where recorded | `REFERENCE-MODEL-IDENTITY.md` | `label`, `file` |
| Grading key (not reproduced here) | `experiments/mem-extraction/round3/corpus.json`, each case's `required`, `forbidden`, `notes` | case id |
| Findings and gaps | `REPORT.md` | §1–§6 |

## 7. Reading order

1. **01** — the prompt record, including what is and is not recoverable.
2. **02A–02R** — one case at a time: what was given, then every response.
3. **03A–03M** — the six-item prompt across configurations.
4. `REFERENCE-EXECUTION-METADATA.md` only when a specific claim needs checking.
5. `REPORT.md` for coverage and integrity statements.

## 8. What this series does not establish

- That the reconstructed prompt equals the bytes the server received.
- That the recorded configuration values were successfully applied.
- That any response is correct, useful or safe.
- Anything about models not present in the stored records.
