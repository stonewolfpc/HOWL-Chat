"""Abliterated / uncensored variants, matched to the existing presets.

Each spec clones the sampler chain actually recorded for the corresponding
baseline corpus run, so a difference is attributable to the weights, not the
configuration. New files, new labels, new result tags; no existing artifact is
touched.
"""
from __future__ import annotations

from harness import MODELS, Spec
from configs import (QWEN35_LARGE_NONTHINK, QWEN35_LARGE_THINK,
                     QWEN35_SMALL_NONTHINK, LFM_INSTRUCT, OFF, ON)

Q4 = MODELS


def ablation() -> list[Spec]:
    return [
        Spec("abl-4b-nonthink", "abl-4b-q4km.gguf", Q4, QWEN35_LARGE_NONTHINK,
             extra_body=OFF, max_tokens=768,
             preset="HauhauCS Aggressive; Qwen3.5-4B non-thinking preset (T0.7/top_p0.8/pres1.5)",
             notes="HauhauCS/Qwen3.5-4B-Uncensored-HauhauCS-Aggressive Q4_K_M"),
        Spec("abl-4b-think-512", "abl-4b-q4km.gguf", Q4, QWEN35_LARGE_THINK,
             extra_body=ON, max_tokens=1024, reasoning_budget=512,
             preset="HauhauCS Aggressive; Qwen3.5-4B thinking preset, runtime budget 512",
             notes="same weights as abl-4b-nonthink, thinking on"),
        Spec("abl-2b-nonthink", "abl-2b-q4km.gguf", Q4, QWEN35_SMALL_NONTHINK,
             extra_body=OFF, max_tokens=768,
             preset="HauhauCS Aggressive; Qwen3.5-2B non-thinking preset (T1.0/top_p1.0/pres2.0)",
             notes="HauhauCS/Qwen3.5-2B-Uncensored-HauhauCS-Aggressive Q4_K_M"),
        Spec("abl-lfm120", "abl-lfm120-q4km.gguf", Q4, LFM_INSTRUCT,
             max_tokens=512,
             preset="Huihui abliterated; Liquid preset T0.1/top_k50/rep1.05",
             notes="Huihui-LFM2.5-1.2B-Instruct-abliterated Q4_K_M"),
        Spec("abl-lfm3-think-512", "abl-lfm3-q4km.gguf", Q4,
             dict(temperature=0.1, top_p=1.0, top_k=50, min_p=0.0,
                  presence_penalty=0.0, repeat_penalty=1.1),
             max_tokens=1024, reasoning_budget=512,
             preset="Heretic abliterated; LFM3 T0.1/top_k50/rep1.1, runtime budget 512",
             notes="Vortecks/LFM3-1.7B-Heretic-Abliterated Q4_K_M"),
    ]


CORPUS_SPECS = {"abl-4b-nonthink", "abl-2b-nonthink", "abl-lfm120"}
