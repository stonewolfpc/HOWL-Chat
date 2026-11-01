#!/usr/bin/env python3
"""
Kokoro TTS CLI Bridge
=====================
Standalone ONNX inference script for Kokoro TTS models.
Called by the Go backend to synthesize speech with per-character voice overrides.

Usage:
    python kokoro-tts.py --text "Hello world" --voice af_heart --speed 1.0 --output out.wav
    python kokoro-tts.py --text "Hello world" --voice af_heart --speed 1.0 --output out.wav \
        --model D:\\AI GGUF\\TTS\\kokoro-v1.0.onnx \
        --voices D:\\AI GGUF\\TTS\\voices-v1.0.bin \
        --config D:\\AI GGUF\\TTS\\hf_cache\\config.json

Dependencies: numpy, onnxruntime, kokoro (for G2P/phonemization)
"""

import argparse
import json
import os
import sys
import zipfile
from pathlib import Path

import numpy as np
import onnxruntime as ort

# Lazy-load kokoro G2P pipeline (only when needed)
_g2p_pipeline = None

def _get_g2p():
    global _g2p_pipeline
    if _g2p_pipeline is None:
        from kokoro import KPipeline
        _g2p_pipeline = KPipeline(lang_code="a", model=False)
    return _g2p_pipeline


# Standalone-CLI fallbacks only: the Go backend always passes --model, --voices
# and --config. Resolve them relative to the project instead of a fixed drive so
# a moved checkout still works, and honour the same env vars the Go side uses.
_PROJECT_DIR = Path(__file__).resolve().parent.parent
_DEFAULT_KOKORO_DIR = _PROJECT_DIR / "models" / "kokoro"

DEFAULT_MODEL = os.environ.get("HOWL_TTS_KOKORO_MODEL") or str(_DEFAULT_KOKORO_DIR / "kokoro-v1.0.onnx")
DEFAULT_VOICES = os.environ.get("HOWL_TTS_KOKORO_VOICES") or str(_DEFAULT_KOKORO_DIR / "voices-v1.0.bin")
DEFAULT_CONFIG = os.environ.get("HOWL_TTS_KOKORO_CONFIG") or str(_DEFAULT_KOKORO_DIR / "hf_cache" / "config.json")


def text_to_ipa(text):
    """Convert text to IPA phonemes using kokoro's G2P pipeline (no espeak-ng binary needed)."""
    pipeline = _get_g2p()
    phonemes = ""
    for result in pipeline(text):
        if hasattr(result, "phonemes") and result.phonemes:
            phonemes += result.phonemes
    return phonemes


def load_vocab(config_path):
    """Load phoneme → token ID mapping from config.json."""
    with open(config_path, "r", encoding="utf-8") as f:
        config = json.load(f)
    vocab = config.get("vocab", {})
    # Build char → id mapping
    return {k: int(v) for k, v in vocab.items()}


def ipa_to_tokens(ipa, vocab, max_len=510):
    """Convert IPA string to token ID array, padding to max_len."""
    tokens = []
    for ch in ipa:
        if ch in vocab:
            tokens.append(vocab[ch])
        else:
            # Unknown character, map to space or skip
            tokens.append(vocab.get(" ", 0))
    # Truncate/pad
    if len(tokens) > max_len:
        tokens = tokens[:max_len]
    tokens = tokens + [0] * (max_len - len(tokens))
    return np.array(tokens, dtype=np.int64).reshape(1, -1)


def load_voice_embedding(voices_zip, voice_name):
    """Extract a voice embedding .npy from the voices ZIP."""
    with zipfile.ZipFile(voices_zip, "r") as z:
        npy_name = f"{voice_name}.npy"
        if npy_name not in z.namelist():
            raise ValueError(f"Voice '{voice_name}' not found in {voices_zip}. Available: {z.namelist()}")
        with z.open(npy_name) as f:
            emb = np.load(f)
    # emb shape varies: (510, 1, 256), (2, 256), or (256,)
    if emb.ndim == 3:
        emb = emb[0]       # (1, 256)
    elif emb.ndim == 2:
        emb = emb[0]       # (256,)
    return emb.astype(np.float32).reshape(1, -1)


def synthesize(model_path, voices_zip, config_path, text, voice, speed, output_path):
    """Run Kokoro ONNX inference and save WAV."""
    vocab = load_vocab(config_path)
    ipa = text_to_ipa(text)
    tokens = ipa_to_tokens(ipa, vocab)
    style = load_voice_embedding(voices_zip, voice)
    speed_arr = np.array([float(speed)], dtype=np.float32)

    sess = ort.InferenceSession(model_path)
    input_names = {inp.name for inp in sess.get_inputs()}
    if "tokens" not in input_names:
        raise RuntimeError(f"Model missing 'tokens' input. Available: {input_names}")

    outputs = sess.run(None, {
        "tokens": tokens,
        "style": style,
        "speed": speed_arr,
    })
    audio = outputs[0].flatten()

    # Write WAV (24kHz mono float32)
    try:
        import scipy.io.wavfile as wavfile
        wavfile.write(output_path, 24000, audio)
    except ImportError:
        # Manual WAV header
        audio_int16 = (audio * 32767).astype(np.int16)
        with open(output_path, "wb") as f:
            import struct
            data = audio_int16.tobytes()
            f.write(b"RIFF")
            f.write(struct.pack("<I", 36 + len(data)))
            f.write(b"WAVE")
            f.write(b"fmt ")
            f.write(struct.pack("<I", 16))
            f.write(struct.pack("<HHIIHH", 1, 1, 24000, 48000, 2, 16))
            f.write(b"data")
            f.write(struct.pack("<I", len(data)))
            f.write(data)

    print(f"Synthesized: {output_path} ({len(audio)/24000:.1f}s)")
    return output_path


def list_voices(voices_zip):
    """List available voice names from the ZIP."""
    with zipfile.ZipFile(voices_zip, "r") as z:
        names = [n.replace(".npy", "") for n in z.namelist() if n.endswith(".npy")]
    return sorted(names)


def main():
    parser = argparse.ArgumentParser(description="Kokoro TTS ONNX Inference")
    parser.add_argument("--text", default="", help="Text to synthesize")
    parser.add_argument("--voice", default="af_heart", help="Voice name (e.g. af_heart, am_michael)")
    parser.add_argument("--speed", type=float, default=1.0, help="Speaking speed (0.5–2.0)")
    parser.add_argument("--output", default="", help="Output WAV file path")
    parser.add_argument("--model", default=DEFAULT_MODEL, help="Path to kokoro-v1.0.onnx")
    parser.add_argument("--voices", default=DEFAULT_VOICES, help="Path to voices-v1.0.bin ZIP")
    parser.add_argument("--config", default=DEFAULT_CONFIG, help="Path to config.json with vocab")
    parser.add_argument("--list-voices", action="store_true", help="List available voices and exit")
    args = parser.parse_args()

    if args.list_voices:
        for v in list_voices(args.voices):
            print(v)
        return

    if not args.text or not args.output:
        parser.error("--text and --output are required (unless using --list-voices)")

    synthesize(
        model_path=args.model,
        voices_zip=args.voices,
        config_path=args.config,
        text=args.text,
        voice=args.voice,
        speed=args.speed,
        output_path=args.output,
    )


if __name__ == "__main__":
    main()
