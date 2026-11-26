# Context isolation: what is reset, what is persisted, what must not be confused

**Date:** 2026-11-19 · **Runtime:** llama.cpp build 12015 (`a34de8618`) ·
**Script:** `../run_isolation.py` · **Raw:** `../results/isolation.json`

## The three layers

| Layer | What it is | Required behaviour | How it was tested |
|---|---|---|---|
| Inference state | KV cache, slot state, any transient per-job context | must not carry from one job to the next | sentinel job followed by a probe on the same server, compared with a fresh-server probe |
| Model persistence | weights resident in RAM/VRAM | an optimization; must never be relied on for correctness | servers kept alive across jobs in the warm condition |
| Memory persistence | MEMORY-09 records and receipts | must never be affected by either of the above | untouched by this test; nothing in the harness writes memories |

The distinction matters because "we keep the model loaded" and "the model
remembers the last job" are different claims, and a benchmark that keeps a
server alive can silently turn the second into an assumption.

## Design

One server per model, `-np 1`, no restart between jobs in the warm condition.

1. **Fresh probe (B)** — a request whose conversation asks the model to repeat
   everything it was told before this message. Sent first, on a freshly loaded
   server, with seed 5150.
2. **Sentinel job (A)** — an extraction job containing two unique, invented
   facts: a dog named **Marrowgale** and a cellar code **7712**.
3. **Erase attempt** — `POST /slots/{id}?action=erase`.
4. **Warm probe (B)** — the identical request, same seed, same server.
5. **Unrelated extraction (C01)** — a normal corpus case, checked for sentinel
   tokens leaking into stored memory.

Leakage is measured by literal token search for `marrowgale` and `7712`, which
cannot occur by chance.

## Results

| Model | Erase endpoint | Fresh vs warm probe | Sentinel in warm probe | Sentinel in unrelated extraction |
|---|---|---|---|---|
| Qwen3.5-2B non-thinking | HTTP 501 | **byte-identical** | no | no |
| Qwen3.5-4B non-thinking | HTTP 501 | **byte-identical** | no | no |

The 4B probe is the substantive one, because it produced content rather than a
token:

> `C1 | HOLD | "Repeat, word for word..." | The user's request implies they
> expect specific prior context or a memory of previous interactions that is not
> present in the current session log…`

That is the correct reading of the probe, and it is identical on the warm server
after job A. The model does not believe it has a previous conversation, because
it does not.

## Interpretation

**Observed:** with `-np 1` and one request at a time, job state does not carry
between `/v1/chat/completions` calls on this build. The erase endpoint returning
501 does not prevent that, because the endpoint is not what provides the
isolation — each request builds a new sequence rather than continuing a
conversation.

**Inferred:** the prefix cache can still be shared, but only over the portion of
the prompt the two jobs genuinely have in common (the fixed instruction). The
conversation suffix differs, so the cached prefix cannot carry job A's facts
into job B; the byte-identical probes are consistent with that and would have
been very unlikely if a longer prefix were being reused.

**Not established:** this does not prove isolation under concurrent requests
(`-np` > 1), across runtime builds, or for prompts that share a long identical
prefix with different trailing content. It also says nothing about MEMORY-09
records, which are a different layer entirely.

**Why Qwen was included:** the earlier rounds' 4B results could have been
explained by cross-run contamination of a favoured model. They are not: the
sentinel survives neither the probe nor an unrelated extraction, and the fresh
and warm outputs are identical. No model is exempt from this test, and no model
gets credit for passing it.

## Consequence for the harness

The harness erases the slot before every job and tolerates the 501. The evidence
above says the erase attempt is not what protects the results. What protects
them is one request per sequence at `-np 1`, plus a fresh server per
configuration. Both are recorded in every run, so a reader can check whether a
particular result was produced under those conditions.
