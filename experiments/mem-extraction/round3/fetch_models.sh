#!/usr/bin/env bash
set -u
cd "$(dirname "$0")/models"
dl() { # repo file
  local repo="$1" f="$2"
  if [ -s "$f" ]; then echo "SKIP $f"; return 0; fi
  echo "GET $repo/$f"
  curl -fL --retry 3 --retry-delay 2 -o "$f.part" \
    "https://huggingface.co/$repo/resolve/main/$f?download=true" && mv "$f.part" "$f" \
    && echo "OK $f $(stat -c%s "$f")" || echo "FAIL $f"
}
dl LiquidAI/LFM2.5-1.2B-Thinking-GGUF LFM2.5-1.2B-Thinking-Q4_K_M.gguf
dl unsloth/Qwen3-1.7B-GGUF Qwen3-1.7B-Q4_K_M.gguf
dl unsloth/Qwen3.5-4B-GGUF Qwen3.5-4B-Q4_K_M.gguf
dl ggml-org/SmolLM3-3B-GGUF SmolLM3-Q4_K_M.gguf
dl unsloth/Qwen3.5-9B-GGUF Qwen3.5-9B-Q4_K_M.gguf
echo DONE
