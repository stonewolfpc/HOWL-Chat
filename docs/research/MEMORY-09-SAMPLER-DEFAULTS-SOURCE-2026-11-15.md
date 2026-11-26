# Sampler Defaults — Source Inspection (KoboldCpp 1.120)

**Identity:** `MEMORY-09-SAMPLER-SRC-1` · **Date:** 2026-11-15 · **Mode:** read-only; no server launch, no generation, no model read
**Question:** what sampler configuration does KoboldCpp 1.120 actually use, and is it what MEMORY-09 intended to test?

---

## 1. The binary under test

| Field | Value |
|---|---|
| Path | `H:\AI apps\koboldcpp.exe` |
| Bytes | 635,698,825 |
| SHA-256 | `6544239ab2747ee84e1ca265a702772fe63719774b829ec9e829f4029ac013db` |
| mtime | 2026-09-15T07:51 (local) |
| Self-reported version | `KoboldCpp 1.120`, from `/api/extra/version` and the startup banner |

**Limit stated up front:** the version is self-reported by the bundle. No hash ties this exact binary to the `v1.120` tag, so everything below is *source-of-that-tag*, not *measured on this process*.

## 2. What the source says (tag `v1.120`, `koboldcpp.py`)

Defaults for parameters a request omits:

```python
temperature = tryparsefloat(genparams.get('temperature', adapter_obj.get("temperature", 0.7)), 0.7)
top_k       = tryparseint  (genparams.get('top_k',       adapter_obj.get("top_k", 100)), 100)
top_p       = tryparsefloat(genparams.get('top_p',       adapter_obj.get("top_p", 0.9)), 0.9)
min_p       = tryparsefloat(genparams.get('min_p',       0.0), 0.0)
typical_p   = tryparsefloat(genparams.get('typical',     1.0), 1.0)
tfs         = tryparsefloat(genparams.get('tfs',         1.0), 1.0)
rep_pen     = tryparsefloat(genparams.get('rep_pen',     adapter_obj.get("rep_pen", 1.0)), 1.0)
sampler_order = genparams.get('sampler_order', [6, 0, 1, 3, 4, 2, 5])
```

**The adapter is consulted before the numeric default.** `adapter_obj` is the ChatCompletions adapter chosen by `--chatcompletionsadapter AutoGuess` — for these models it resolves to a chat template adapter. So "no override" does not mean "0.7 / 0.9 / 100"; it means *the adapter's value if it declares one, otherwise 0.7 / 0.9 / 100*. The adapter's own JSON was not inspected, so for `top_k` and `rep_pen` the effective value is **adapter-or-default, unresolved**.

## 3. What the experiment actually sent

Frozen params: `temperature=0.8`, `top_p=0.95`, `max_tokens=256`. **No `top_k`, no `rep_pen`, no `min_p`, no `sampler_order`.**

Under the logic above, the effective values for an S1 call would be:

| Sampler | Request sends | Effective (source-derived) | In the frozen record? |
|---|---|---|---|
| temperature | 0.8 | **0.8** | yes |
| top_p | 0.95 | **0.95** | yes |
| top_k | — | **100** (adapter may override) | **no** |
| min_p | — | 0.0 → inactive | no |
| typical | — | 1.0 → inactive | no |
| tfs | — | 1.0 → inactive | no |
| rep_pen | — | **1.0** (adapter may override) → arithmetic no-op, sampler still in chain | **no** |
| sampler_order | — | `[6,0,1,3,4,2,5]` | **no** |

Chain construction order in the C++ path (`gpttype_adapter.cpp`): rep-pen → logit-bias → `top_k` (if > 0) → `top_p` (if 0 < p < 1) → `min_p` (if > 0) → typical (if 0 < p < 1) → temperature, then a distribution sampler seeded; if `temperature == 0`, greedy instead.

So the real chain for S1 is **rep_pen(1.0) → top_k(100) → top_p(0.95) → temp(0.8) → dist**, and the frozen parameter table names only two of its five stages.

## 4. The finding that matters for the next decision

`--debugmode` does **not** print a resolved sampler chain. It prints the received request:

```python
if args.debugmode >= 1:
    printablegenparams = truncate_long_json(genparams, trunc_len)   # before transform_genparams
```

That is why Route A produced nothing — the print lives inside the generation handler, not at startup — and it means **Route B would echo our own request parameters back, not the effective chain.** Spending a generation for that would buy confirmation of what we already know and would not resolve `top_k`.

**Route B is therefore withdrawn as a sampler-capture method.** It retains value only as a behavioural observation, which is a different experiment.

## 5. What this establishes, and what it does not

**Establishes (source-derived):** the default set at tag 1.120; that the adapter is consulted first; the chain construction order; that `debugmode` logs received params, not resolved samplers; that the frozen parameter record omits `top_k`, `rep_pen` and `sampler_order` — with `top_k=100` materially active.

**Does not establish:** that this binary equals the tag (self-reported only), or the chain as a runtime observation. **The adapter question has since been settled and is no longer in this list** — see §8.

## 6. Declared configuration — the documentation correction (applied)

Per the owner's decision, the experiment's record now distinguishes three evidence classes instead of asserting one blurry set of numbers. **No frozen parameter, prompt, fixture or protocol value was changed.**

| Class | Values | Basis |
|---|---|---|
| **Request-derived** | `temperature=0.8`, `top_p=0.95`, `max_tokens=256` | what the S1 payload sends |
| **Source-derived** | `top_k=100`, `rep_pen=1.0`, `sampler_order=[6,0,1,3,4,2,5]`; `min_p=0`, `typical=1.0`, `tfs=1.0` (inactive) | fallback branches in `koboldcpp.py` at tag `v1.120` |
| **Unresolved** | binary-to-tag binding; the chain as a direct runtime observation | version self-reported; no sampler logging exists in this build |

**Adapter overrides are now ruled out** (§8). **The configuration is still not fully specified, and improving the description does not create observability.** `top_k=100` follows from source plus the absence of any adapter sampler key — it is a strong derivation, not a measurement of what the server did.

## 7. Recommendation

Do not spend a generation to "capture samplers" — the mechanism cannot show them.

The defensible options, in order of honesty:

1. **Record the effective chain as source-derived**, with the version caveat and the adapter caveat both stated, and add `top_k=100`, `rep_pen=1.0`, `sampler_order=[6,0,1,3,4,2,5]` to the experiment's declared configuration — this changes nothing about what runs; it stops the record from under-describing it.
2. **Make the chain explicit in the request** — send `top_k`, `rep_pen`, `sampler_order` in the S1 payload so the effective values stop depending on an adapter — but that is a **protocol change** requiring a version bump, not an edit.
3. Behavioural inference, later, deliberately.

Option 1 is a documentation correction and is what I would recommend; option 2 changes the experiment and belongs to the owner.

**No generation was spent. Budget remains 1 of 30.**

## 8. Adapter recovery — done, and it closes the question

**Method.** `koboldcpp.exe` is a PyInstaller one-file bundle: a launcher followed by a `CArchive` whose table of contents sits at the end of the file. `eval/runs/extract-kcpp-adapters.py` parses that TOC (magic `MEI\014\013\012\013\016`, 88-byte cookie, 1,311 entries, `pyvers 308`) and decompresses the adapter entries with zlib. **Nothing was executed** — the file was only read.

**Result: 30 adapters extracted**, 160–12,474 B each, into `eval/runs/adapter-extract/`.

**Every one contains template strings only** — `system_start/end`, `user_start/end`, `assistant_start/end/gen`. A recursive key scan across all 30 files for `temperature`, `top_k`, `top_p`, `rep_pen`, `min_p`, `typical`, `tfs`, `sampler_order`, `dynatemp_range` returned **no matches**. `AutoGuess.json` is a 32-entry candidate list keyed by prompt-marker `search` strings; `AutoGuess-NoThink.json` is the 30-entry variant — neither carries sampler values.

**Which adapter each model resolved to** (from the servers' own logs, `Chat completion heuristic:`):

| Slot | Adapter chosen | Declares samplers? |
|---|---|---|
| `S1-A` Qwen3.5-9B | `ChatML (Generic)` | no |
| `S1-B` Gemma-4-E4B | `Google Gemma 4 (E2B and E4B)` | no |
| `S1-C` LFM2.5-8B-A1B | `ChatML (Generic)` | no |

**Consequence.** In `adapter_obj.get("top_k", 100)` and `adapter_obj.get("rep_pen", 1.0)`, the adapter never supplies a value, so the fallback always applies. The effective chain for an S1 call is therefore:

**rep_pen(1.0) → top_k(100) → top_p(0.95) → temp(0.8) → dist**, order `[6,0,1,3,4,2,5]`

with `min_p`, `typical` and `tfs` inactive. All three slots share the same sampler behaviour, because their adapters differ only in chat framing.

**What still is not established:** that this binary is the `v1.120` tag (version is self-reported), and any direct observation of the chain at runtime. Those are limits of evidence, not open questions about the adapters.
