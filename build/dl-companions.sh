#!/bin/bash
# One-off fetch of the companion files the dependency catalogue names, into the
# folder the model scanner already indexes. Sizes are checked against the
# publisher's Content-Length so a truncated transfer cannot masquerade as done.
set -u
cd "/g/AI GGUF/Image models" || exit 1

dl() {
  name="$1"; url="$2"; want="$3"
  echo "[start] $name"
  if curl -L --fail --retry 3 --retry-delay 3 -C - -sS -o "$name" "$url"; then
    got=$(stat -c %s "$name" 2>/dev/null || echo 0)
    if [ "$got" = "$want" ]; then
      echo "[ok] $name $got bytes"
    else
      echo "[size-mismatch] $name got $got want $want"
    fi
  else
    echo "[FAIL] $name"
  fi
}

dl clip_l.safetensors \
  "https://huggingface.co/comfyanonymous/flux_text_encoders/resolve/main/clip_l.safetensors" \
  246144152 &
dl qwen_image_vae.safetensors \
  "https://huggingface.co/circlestone-labs/Anima/resolve/main/split_files/vae/qwen_image_vae.safetensors" \
  253806246 &
dl qwen_3_4b.safetensors \
  "https://huggingface.co/Comfy-Org/Z-Image-Turbo/resolve/main/split_files/text_encoders/qwen_3_4b.safetensors" \
  8044982048 &
wait

dl t5xxl_fp8_e4m3fn.safetensors \
  "https://huggingface.co/comfyanonymous/flux_text_encoders/resolve/main/t5xxl_fp8_e4m3fn.safetensors" \
  4893934904 &
dl qwen_2.5_vl_7b_fp8_scaled.safetensors \
  "https://huggingface.co/Comfy-Org/Qwen-Image/resolve/main/split_files/text_encoders/qwen_2.5_vl_7b_fp8_scaled.safetensors" \
  0 &
wait
echo "[all downloads finished]"
ls -la --time-style=+%H:%M *.safetensors | awk '{printf "%s %s\n", $5, $7}'
