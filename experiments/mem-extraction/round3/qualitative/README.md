# Qualitative evaluation track

This directory holds the human-readable evidence for the memory-extraction
model choice. It is deliberately separate from the scoring pipeline: the scores
explain the measurements, the evaluation decides whether the memories are any
good.

## Three independent tracks

| Track | Owner | Artifact | Status |
|---|---|---|---|
| Semantic evaluation | Buffy | `EVALUATION-BUFFY.md` (this directory) | done |
| Experimental dossier | DeepSeek | not produced here — separate track | pending |
| Adversarial review | project owner | reviews the dossier without reading the evaluation first | pending |

The evaluation was written from the raw outputs alone, before any dossier
existed. That ordering is deliberate: it keeps my reading independent of
DeepSeek's analysis, and it keeps the owner's review independent of mine. If a
later dossier contradicts something here, the contradiction is evidence, not an
error to reconcile quietly.

## Evidence labels

Every claim in this track is tagged, so a reader can tell what is measured from
what is believed.

- **Observed** — read directly in a stored response.
- **Source-derived** — taken from a publisher's model card, paper or repository,
  with the URL and retrieval date in `../SOURCES.md`.
- **Calculated** — computed from stored measurements (means, ratios, deltas).
- **Inferred** — my reading of why a behaviour occurred; may be wrong.
- **Unestablished** — asserted by nobody, not tested, or tested too thinly.
- **Superseded** — true of an earlier round or configuration and no longer current.

## What is preserved for every case

Per the evaluation format, each case file in `archive/` contains:

1. the exact prompt the model received (instruction + conversation, byte for
   byte as the harness built it);
2. every response verbatim, including `reasoning_content` where the model
   produced it, for every draw;
3. the reference material — required items, expected decisions, traps and
   forbidden strings — written before any output existed;
4. the other configurations' responses to the same case;
5. the configuration: model, quantization, sampler chain, thinking mode, token
   budget, context size, seed, finish reason, token count and run code.

The run identifier (`run_code`) joins every document to the raw record in
`../results/corpus.json`, so any claim can be checked against the JSON.

Nothing was discarded. Unusual, malformed, truncated and empty outputs are all
present; the ones that look like noise are often the most informative (TinyLlama
looping, LFM3 leaking its deliberation into the answer).

## File map

```
qualitative/
  READING.md              compact side-by-side: draw 1 of every model per case
  EVALUATION-BUFFY.md     the semantic assessment (the authority for this track)
  SECRETS-AND-MEMORY.md   C17 revision: extractor fidelity vs MEMORY-09 authority
  CONTEXT-ISOLATION.md    the three-layer reset test and its results
  archive/case-C01.md …   full verbatim evidence per case
  models/<config>.md      one model at a time, all 18 cases
```

`READING.md` is the file to read first. `archive/` is the file to cite.

## Regenerating

```
python qual_extract.py     # rebuild everything in this directory from results/corpus.json
```

Regeneration is deterministic and reads only the preserved runs; it does not
call a model. `EVALUATION-BUFFY.md` and `CONTEXT-ISOLATION.md` are hand-written
and are not overwritten.

## Scope boundary

This track adds interpretation. It changes no score, no ranking, no finalist and
no requirement mapping. It does not touch production Howl, MEMORY-09 records,
ratified architecture or the frozen 118-requirement mapping.
