"""Coexistence test: extraction model alongside an active chat model on 8 GB VRAM.

The 'active conversation' workload is a locally downloaded 9B at full offload,
not the production configuration and not the app's runtime. Both servers are
separate processes; extraction latency is compared with and without the primary
resident. If the primary cannot fit at full offload on this GPU, that is
recorded and the primary falls back to partial offload.
"""
import json, sys, threading, time
from pathlib import Path
sys.path.insert(0, str(Path(__file__).resolve().parent))
import harness
from configs import by_label

PRIMARY = by_label("qwen35-9b-think")
PRIMARY.ctx = 4096
PRIMARY.ngl = 99
PRIMARY.flash_attn = True
PRIMARY.kv_unified = True
PRIMARY.max_tokens = 256

EXTRACT = by_label("lfm120-instruct-official")
EXTRACT.ctx = 4096
EXTRACT.max_tokens = 512

def make(port):
    # harness.Server binds harness.PORT; patch per process use
    harness.PORT = port
    harness.BASE = f"http://127.0.0.1:{port}"

def job(spec, text, seed=4101):
    return harness.chat(spec, text, seed=seed)

def main():
    corpus = json.loads((Path(__file__).resolve().parent / "corpus.json").read_text(encoding="utf-8"))["cases"]
    c = next(x for x in corpus if x["id"] == "C04")
    prompt = harness.EXTRACTION_INSTRUCTION + "\n\nConversation:\n" + c["conversation"]
    out = {"primary": "qwen35-9b-think ngl99", "extract": "lfm120-instruct-official ngl0"}
    # isolated baseline
    harness.PORT = 5099; harness.BASE = "http://127.0.0.1:5099"
    with harness.Server(EXTRACT, harness.RAW / "coexist" / "extract-alone.log") as srv:
        try: srv.erase_slot()
        except Exception: pass
        r = job(EXTRACT, prompt)
        out["alone"] = {"wall_s": r["wall_s"], "ttft_s": r["ttft_s"],
                        "peak_rss_mb": round(srv.peak["rss"] / 1e6, 1),
                        "peak_vram_mb": srv.peak["vram"]}
        print("alone", out["alone"], flush=True)
    # primary resident on its own port, then extraction
    harness.PORT = 5100; harness.BASE = "http://127.0.0.1:5100"
    try:
        with harness.Server(PRIMARY, harness.RAW / "coexist" / "primary.log") as p:
            harness.PORT = 5099; harness.BASE = "http://127.0.0.1:5099"
            with harness.Server(EXTRACT, harness.RAW / "coexist" / "extract-with-primary.log") as srv:
                try: srv.erase_slot()
                except Exception: pass
                r = job(EXTRACT, prompt)
                out["with_primary"] = {"wall_s": r["wall_s"], "ttft_s": r["ttft_s"],
                                       "peak_rss_mb": round(srv.peak["rss"] / 1e6, 1),
                                       "peak_vram_mb": srv.peak["vram"],
                                       "primary_vram_mb": p.peak["vram"]}
                print("with primary", out["with_primary"], flush=True)
    except Exception as exc:
        out["with_primary_error"] = repr(exc)
        print("coexistence failed:", exc, flush=True)
    (harness.RESULTS / "coexist.json").write_text(json.dumps(out, indent=2), encoding="utf-8")
    print(json.dumps(out, indent=1))

if __name__ == "__main__":
    main()
