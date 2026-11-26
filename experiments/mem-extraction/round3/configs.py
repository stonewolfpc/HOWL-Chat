"""Round-3 configuration registry.

Every sampler value here is traceable to a publisher statement (see
SOURCES.md). Nothing here is inferred from another model in the same family:
Qwen3.5-0.8B/2B and 4B/9B have *different* official presets, and the LFM
Instruct and Thinking checkpoints have different temperatures.

min_p is always set explicitly, because llama.cpp has a non-zero default
(0.05) and the publishers' presets are stated as 0.0. Leaving it unset would
silently change the chain.
"""

from __future__ import annotations

from harness import MODELS, ROUND1_MODELS, ROUND2_MODELS, Spec

# ---- official presets -------------------------------------------------
QWEN35_SMALL_NONTHINK = dict(temperature=1.0, top_p=1.0, top_k=20, min_p=0.0,
                             presence_penalty=2.0, repeat_penalty=1.0)
QWEN35_SMALL_THINK = dict(temperature=1.0, top_p=0.95, top_k=20, min_p=0.0,
                          presence_penalty=1.5, repeat_penalty=1.0)
QWEN35_LARGE_NONTHINK = dict(temperature=0.7, top_p=0.8, top_k=20, min_p=0.0,
                             presence_penalty=1.5, repeat_penalty=1.0)
QWEN35_LARGE_THINK = dict(temperature=1.0, top_p=0.95, top_k=20, min_p=0.0,
                          presence_penalty=1.5, repeat_penalty=1.0)
QWEN3_THINK = dict(temperature=0.6, top_p=0.95, top_k=20, min_p=0.0,
                   presence_penalty=0.0, repeat_penalty=1.0)
QWEN3_NONTHINK = dict(temperature=0.7, top_p=0.8, top_k=20, min_p=0.0,
                      presence_penalty=0.0, repeat_penalty=1.0)
LFM_INSTRUCT = dict(temperature=0.1, top_p=1.0, top_k=50, min_p=0.0,
                    presence_penalty=0.0, repeat_penalty=1.05)
LFM_THINK = dict(temperature=0.05, top_p=1.0, top_k=50, min_p=0.0,
                 presence_penalty=0.0, repeat_penalty=1.05)
SMOLLM3 = dict(temperature=0.6, top_p=0.95, top_k=0, min_p=0.0,
               presence_penalty=0.0, repeat_penalty=1.0)

OFF = {"chat_template_kwargs": {"enable_thinking": False}}
ON = {"chat_template_kwargs": {"enable_thinking": True}}


def _lfm(label, file, sampling, model_dir=ROUND1_MODELS, **kw):
    return Spec(label=label, file=file, model_dir=model_dir, sampling=sampling, **kw)


def stage_a() -> list[Spec]:
    """One official configuration per model/mode. Max output 2048 tokens so a
    thinking model has room for reasoning *and* the six-line answer."""
    return [
        _lfm("lfm120-instruct-official", "LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf",
             LFM_INSTRUCT, preset="Liquid AI card: T0.1 / top_k50 / rep1.05",
             max_tokens=2048),
        _lfm("lfm350-official", "LFM2.5-350M-QAD-Q4_0.gguf", LFM_INSTRUCT,
             preset="Liquid AI card: T0.1 / top_k50 / rep1.05", max_tokens=2048),
        Spec("lfm120-think-official", "LFM2.5-1.2B-Thinking-Q4_K_M.gguf", MODELS,
             LFM_THINK, preset="Liquid AI card: T0.05 / top_k50 / rep1.05",
             max_tokens=3072, notes="thinking always on, no effort levels, no off switch"),
        Spec("qwen35-08b-nonthink", "Qwen3.5-0.8B-Q4_K_M.gguf", ROUND2_MODELS,
             QWEN35_SMALL_NONTHINK, extra_body=OFF,
             preset="Qwen3.5 small non-thinking text preset", max_tokens=2048),
        Spec("qwen35-2b-nonthink", "Qwen3.5-2B-Q4_K_M.gguf", ROUND2_MODELS,
             QWEN35_SMALL_NONTHINK, extra_body=OFF,
             preset="Qwen3.5 small non-thinking text preset", max_tokens=2048),
        Spec("qwen35-2b-think", "Qwen3.5-2B-Q4_K_M.gguf", ROUND2_MODELS,
             QWEN35_SMALL_THINK, extra_body=ON,
             preset="Qwen3.5 small thinking text preset", max_tokens=3072),
        Spec("qwen35-4b-nonthink", "Qwen3.5-4B-Q4_K_M.gguf", MODELS,
             QWEN35_LARGE_NONTHINK, extra_body=OFF,
             preset="Qwen3.5 4B/9B non-thinking general preset", max_tokens=2048),
        Spec("qwen35-4b-think", "Qwen3.5-4B-Q4_K_M.gguf", MODELS,
             QWEN35_LARGE_THINK, extra_body=ON,
             preset="Qwen3.5 4B/9B thinking general preset", max_tokens=3072),
        Spec("qwen35-9b-think", "Qwen3.5-9B-Q4_K_M.gguf", MODELS,
             QWEN35_LARGE_THINK, extra_body=ON,
             preset="Qwen3.5 4B/9B thinking general preset", max_tokens=3072,
             notes="quality-ceiling candidate"),
        Spec("qwen3-17b-think", "Qwen3-1.7B-Q4_K_M.gguf", MODELS, QWEN3_THINK,
             extra_body=ON, preset="Qwen3 thinking preset T0.6/top_p0.95",
             max_tokens=3072),
        Spec("qwen3-17b-nonthink", "Qwen3-1.7B-Q4_K_M.gguf", MODELS,
             QWEN3_NONTHINK, extra_body=OFF,
             preset="Qwen3 non-thinking preset T0.7/top_p0.8", max_tokens=2048),
        Spec("smollm3-3b-think", "SmolLM3-Q4_K_M.gguf", MODELS, SMOLLM3,
             system="/think", preset="SmolLM3 extended thinking, T0.6/top_p0.95",
             max_tokens=3072),
        Spec("smollm3-3b-nonthink", "SmolLM3-Q4_K_M.gguf", MODELS, SMOLLM3,
             system="/no_think", preset="SmolLM3 no_think, T0.6/top_p0.95",
             max_tokens=2048),
    ]


LFM3 = dict(temperature=0.1, top_p=1.0, top_k=50, min_p=0.0,
           presence_penalty=0.0, repeat_penalty=1.1)


def lfm3(thinking: bool) -> Spec:
    return Spec(f"lfm3-17b-{'think' if thinking else 'nonthink'}",
                "LFM3-1.7B-Q4_K_M.gguf", MODELS, LFM3,
                extra_body={"chat_template_kwargs": {"enable_thinking": thinking}},
                preset="LFM3 card: T0.1 / top_k50 / rep1.1", max_tokens=3072,
                notes="reasoning default on; enable_thinking=False is vendor-flagged as inconsistent")


TINYLLAMA = dict(temperature=0.7, top_p=0.95, top_k=50, min_p=0.0,
                 presence_penalty=0.0, repeat_penalty=1.0)


def tinyllama() -> Spec:
    return Spec("tinyllama-11b-chat", "tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf", MODELS,
                TINYLLAMA, preset="TinyLlama card: T0.7 / top_k50 / top_p0.95",
                max_tokens=768, ctx=2048,
                notes="latest official TinyLlama chat (2024-01); 2K ctx, no reasoning mode")


def by_label(label: str) -> Spec:
    for s in stage_a():
        if s.label == label:
            return s
    if label.startswith("lfm3"):
        return lfm3("think" in label)
    if label == "tinyllama-11b-chat":
        return tinyllama()
    raise KeyError(label)
