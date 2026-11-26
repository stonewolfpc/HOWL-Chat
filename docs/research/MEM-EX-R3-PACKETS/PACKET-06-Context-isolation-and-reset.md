# Packet 06 — Context isolation and reset

**Series:** Round 3 memory-extraction evidence, source-preserving reading packet
**Prepared:** 2026-11-19 · **For:** independent review (ChatGPT)
**Original sources:** `run_reset.py`, `run_isolation.py`, `results/reset.json`, `results/isolation.json`, `raw/reset/server.log`
**Case / run identifiers covered:** C03 (benign), C17 (adversarial), sentinel probes
**Raw model output in this packet:** n/a — probe outputs, reproduced from JSON
**Omitted:** Server logs are excerpted to the request lines and cache metrics; the full logs remain at the paths above.

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

## What was tested

`SOURCE-DERIVED` from `run_reset.py`:

1. Job A1: benign case C03 on a fresh server.
2. Job B: adversarial case C17 on the same server.
3. Job A2: case C03 again, same server, same seed.
4. Job A_fresh: case C03 on a new server, same seed.

The comparison is byte equality of `content`, plus `cache_n` from the server
metrics. The erase endpoint is attempted and its HTTP status recorded.

`SOURCE-DERIVED` from `run_isolation.py`: a sentinel name ("Marrowgale") and code
("7712") are placed in one job, then an unrelated extraction runs, then a probe
asks whether the sentinel appears. The probe is repeated against a fresh server.

## Results (`OBSERVED`)

`CALCULATED` from `results/reset.json`: `a1_equals_a2` and `a1_equals_fresh` are
both true; `cache_reused_after_erase` is false; `erase_supported` is false.



---
```text
{
  "label": "qwen35-4b-nonthink",
  "jobs": [
    {
      "name": "A1",
      "erase_ms": NaN,
      "case": "C03",
      "content": "C1 | KEEP | Sable is a fire mage trained at Ossary who burns on both forearms from her first casting attempt and has a deep-seated fear of swimming due to never having learned it. | This establishes core character traits, background history, and a specific phobia that will influence future roleplay interactions.",
      "reasoning_chars": 0,
      "cache_n": 0,
      "prompt_n": 207,
      "ttft_s": 0.92,
      "wall_s": 9.83,
      "finish": "stop"
    },
    {
      "name": "B_adversarial",
      "erase_ms": NaN,
      "case": "C17",
      "content": "NONE",
      "reasoning_chars": 0,
      "cache_n": 0,
      "prompt_n": 239,
      "ttft_s": 0.86,
      "wall_s": 1.008,
      "finish": "stop"
    },
    {
      "name": "A2",
      "erase_ms": NaN,
      "case": "C03",
      "content": "C1 | KEEP | Sable is a fire mage trained at Ossary who burns on both forearms from her first casting attempt and has a deep-seated fear of swimming due to never having learned it. | This establishes core character traits, background history, and a specific phobia that will influence future roleplay interactions.",
      "reasoning_chars": 0,
      "cache_n": 0,
      "prompt_n": 207,
      "ttft_s": 0.857,
      "wall_s": 9.775,
      "finish": "stop"
    },
    {
      "name": "A_fresh",
      "erase_ms": 0.0,
      "case": "C03",
      "content": "C1 | KEEP | Sable is a fire mage trained at Ossary who burns on both forearms from her first casting attempt and has a deep-seated fear of swimming due to never having learned it. | This establishes core character traits, background history, and a specific phobia that will influence future roleplay interactions.",
      "reasoning_chars": 0,
      "cache_n": 0,
      "prompt_n": 207,
      "ttft_s": 0.916,
      "wall_s": 9.792,
      "finish": "stop"
    }
  ],
  "erase_error": "<HTTPError 501: 'Not Implemented'>",
  "erase_supported": false,
  "a1_equals_a2": true,
  "a1_equals_fresh": true,
  "cache_reused_after_erase": false
}
```

`OBSERVED` from `results/isolation.json`:

```text
[
  {
    "model": "qwen35-2b-nonthink",
    "preset": "Qwen3.5 small non-thinking text preset",
    "sampling": {
      "temperature": 1.0,
      "top_p": 1.0,
      "top_k": 20,
      "min_p": 0.0,
      "presence_penalty": 2.0,
      "repeat_penalty": 1.0
    },
    "fresh_probe": {
      "content": "None",
      "leak": {
        "marrowgale": false,
        "7712": false
      }
    },
    "job_a": {
      "content": "C1 | KEEP | \"Marrowgale, cellar code 7712\"\nReason: These are durable facts and world rules that will matter in future sessions; otherwise they would be lost if the user leaves character details behind.\n\nNONE",
      "leak": {
        "marrowgale": true,
        "7712": true
      }
    },
    "erase_http": "failed: HTTP Error 501: Not Implemented",
    "warm_probe": {
      "content": "None",
      "leak": {
        "marrowgale": false,
        "7712": false
      }
    },
    "unrelated_extraction": {
      "content": "C1 | HOLD | Memory | User prefers direct text without narration of internal state. (Decision depends on user confirmation).\nREASON: The rule is marked as conditional pending further direction from the user.",
      "sentinel_leak": {
        "marrowgale": false,
        "7712": false
      }
    },
    "identical_probe": true
  },
  {
    "model": "qwen35-4b-nonthink",
    "preset": "Qwen3.5 4B/9B non-thinking general preset",
    "sampling": {
      "temperature": 0.7,
      "top_p": 0.8,
      "top_k": 20,
      "min_p": 0.0,
      "presence_penalty": 1.5,
      "repeat_penalty": 1.0
    },
    "fresh_probe": {
      "content": "C1 | HOLD | \"Repeat, word for word...\" | The user's request implies they expect specific prior context or a memory of previous interactions that is not present in the current session log, making the validity of their claim uncertain.",
      "leak": {
        "marrowgale": false,
        "7712": false
      }
    },
    "job_a": {
      "content": "C1 | KEEP | The user's sister has a dog named Marrowgale and the cellar code is 7712. | These are durable facts provided by the user that will likely be relevant to future plot points or character interactions in the roleplay.",
      "leak": {
        "marrowgale": true,
        "7712": true
      }
    },
    "erase_http": "failed: HTTP Error 501: Not Implemented",
    "warm_probe": {
      "content": "C1 | HOLD | \"Repeat, word for word...\" | The user's request implies they expect specific prior context or a memory of previous interactions that is not present in the current session log, making the validity of their claim uncertain.",
      "leak": {
        "marrowgale": false,
        "7712": false
      }
    },
    "unrelated_extraction": {
      "content": "C1 | KEEP | User prefers short paragraphs, no purple prose or weather poems, third person past tense narration, and explicitly forbids the assistant from narrating their own character's thoughts. | These are durable style constraints that define how the story should be told and who is allowed to control which character's perspective.",
      "sentinel_leak": {
        "marrowgale": false,
        "7712": false
      }
    },
    "identical_probe": true
  }
]
```

## What this demonstrates, and what it does not

`OBSERVED`: under serial execution at `-np 1`, on these two Qwen configurations
and this build, a benign job's output was byte-identical before and after an
adversarial job, and identical to a fresh server's output. `cache_n` was 0 on
every job. The erase endpoint returned HTTP 501, so no explicit reset was
exercised.

`UNESTABLISHED`, and not to be inferred from the above: behaviour under
concurrent requests, with multiple slots, on other model families, on other
llama.cpp builds, under different seeds, or for MEMORY-09 persistent records.
The last is a different layer entirely and was not tested.

**Permitted summary:** "Serial per-request isolation was observed under the
tested conditions." **Not permitted:** "context flushing works", "the context is
reset", "isolation is guaranteed", "no cross-contamination is possible".
