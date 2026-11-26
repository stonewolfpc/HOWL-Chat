"""Round-3 local memory-extraction harness (CPU baseline + optional GPU offload).

Isolated from eval/ and from the production app. Does not read or write any
MEMORY-09 artifact, model registry, or inference configuration.

Design notes
------------
* One llama-server process per *session* (model + launch flags). Requests within
  a session vary sampling / thinking budget / seed, so launch-level reasoning
  controls are fixed per session and recorded per request.
* Streaming is used so that time-to-first-token, thinking duration and
  answer duration are measured separately. Non-streaming timings are the
  fallback if a build returns no usage chunk.
* Every output is stored with a random run_code for blind scoring.
"""

from __future__ import annotations

import json
import os
import random
import re
import string
import subprocess
import sys
import threading
import time
import urllib.error
import urllib.request
from dataclasses import dataclass, field
from pathlib import Path

import psutil

ROOT = Path(__file__).resolve().parent
REPO = ROOT.parent.parent.parent
SERVER = REPO / "runtime" / "llama" / "llama-server.exe"
MODELS = ROOT / "models"
ROUND1_MODELS = REPO / "experiments" / "lfm25-cpu" / "models"
ROUND2_MODELS = REPO / "experiments" / "lfm25-cpu" / "round2" / "models"
RAW = ROOT / "raw"
RESULTS = ROOT / "results"
PORT = 5099
BASE = f"http://127.0.0.1:{PORT}"

# Original round-1 prompt, imported rather than retyped so it cannot drift.
sys.path.insert(0, str(REPO / "experiments" / "lfm25-cpu"))
from run_probe import PROMPT as TRIAGE_PROMPT  # noqa: E402

EXTRACTION_INSTRUCTION = (
    "You are the memory extractor for a long-running roleplay. Read the conversation and "
    "decide what belongs in long-term memory.\n"
    "Rules:\n"
    "- KEEP: durable facts, preferences, relationships, world rules, or corrections that will "
    "matter in later sessions.\n"
    "- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new "
    "content.\n"
    "- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.\n"
    "Output exactly one line per candidate in the form:\n"
    "ID | DECISION | MEMORY | REASON\n"
    "Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE\n"
    "Do not invent details that are not in the conversation."
)


@dataclass
class Spec:
    label: str
    file: str
    model_dir: Path
    sampling: dict
    max_tokens: int = 768
    ctx: int = 4096
    ngl: int = 0
    reasoning_budget: int | None = None
    reasoning_effort: str | None = None
    reasoning: str | None = None          # on / off / auto, server flag
    extra_body: dict = field(default_factory=dict)
    system: str | None = None
    preset: str = ""
    notes: str = ""
    threads: int | None = None
    flash_attn: bool = False
    kv_unified: bool = False

    @property
    def path(self) -> Path:
        return self.model_dir / self.file


def _q() -> str:
    return "".join(random.choices(string.ascii_lowercase + string.digits, k=8))


def tree_rss(proc: psutil.Process) -> int:
    total = 0
    try:
        procs = [proc, *proc.children(recursive=True)]
    except psutil.NoSuchProcess:
        return 0
    for p in procs:
        try:
            total += p.memory_info().rss
        except psutil.NoSuchProcess:
            pass
    return total


def _nvidia(fields: str) -> str | None:
    try:
        out = subprocess.run(
            ["nvidia-smi", f"--query-gpu={fields}", "--format=csv,noheader,nounits"],
            capture_output=True, text=True, timeout=5)
        if out.returncode == 0:
            return out.stdout.strip()
    except (OSError, subprocess.SubprocessError):
        return None
    return None


def wait_health(timeout: float, proc: subprocess.Popen | None = None) -> float:
    start = time.perf_counter()
    while time.perf_counter() - start < timeout:
        if proc is not None and proc.poll() is not None:
            raise RuntimeError(f"llama-server exited with code {proc.returncode} "
                               f"before becoming healthy")
        try:
            with urllib.request.urlopen(f"{BASE}/health", timeout=2) as r:
                if r.status == 200:
                    return time.perf_counter() - start
        except (urllib.error.URLError, OSError):
            pass
        time.sleep(0.1)
    raise TimeoutError(f"server not healthy after {timeout}s")


class Server:
    def __init__(self, spec: Spec, log_path: Path):
        self.spec = spec
        self.proc: subprocess.Popen | None = None
        self.peak = {"rss": 0, "vram": 0, "gpu_util": 0, "cpu": 0.0}
        self._stop = threading.Event()
        self._sampler: threading.Thread | None = None
        self.log_path = log_path
        self.load_s = 0.0

    def _cmd(self) -> list[str]:
        s = self.spec
        cmd = [
            str(SERVER), "-m", str(s.path), "-c", str(s.ctx),
            "-ngl", str(s.ngl), "-np", "1",
            "--host", "127.0.0.1", "--port", str(PORT),
            "--jinja", "--reasoning-format", "deepseek",
            "--metrics",
        ]
        if s.threads:
            cmd += ["-t", str(s.threads)]
        if s.flash_attn:
            cmd += ["-fa", "on"]
        if s.kv_unified:
            cmd += ["--kv-unified"]
        if s.reasoning_budget is not None:
            cmd += ["--reasoning-budget", str(s.reasoning_budget)]
        if s.reasoning_effort:
            cmd += ["--reasoning-effort", s.reasoning_effort]
        if s.reasoning:
            cmd += ["-rea", s.reasoning]
        return cmd

    def __enter__(self):
        self.log_path.parent.mkdir(parents=True, exist_ok=True)
        t0 = time.perf_counter()
        with open(self.log_path, "wb") as log:
            self.proc = subprocess.Popen(
                self._cmd(), stdout=log, stderr=subprocess.STDOUT,
                cwd=str(SERVER.parent))
        self.load_s = wait_health(300, self.proc)
        self._stop.clear()
        self._sampler = threading.Thread(target=self._sample, daemon=True)
        self._sampler.start()
        self._t0 = t0
        return self

    def _sample(self):
        ps = psutil.Process(self.proc.pid)
        while not self._stop.is_set():
            self.peak["rss"] = max(self.peak["rss"], tree_rss(ps))
            v = _nvidia("memory.used,utilization.gpu")
            if v:
                try:
                    used, util = [int(float(x)) for x in v.split(",")[:2]]
                    self.peak["vram"] = max(self.peak["vram"], used)
                    self.peak["gpu_util"] = max(self.peak["gpu_util"], util)
                except ValueError:
                    pass
            try:
                self.peak["cpu"] = max(self.peak["cpu"], ps.cpu_percent(interval=None))
            except psutil.NoSuchProcess:
                pass
            self._stop.wait(0.1)

    def __exit__(self, *exc):
        self._stop.set()
        if self._sampler:
            self._sampler.join(timeout=2)
        elapsed = round(time.perf_counter() - self._t0, 2)
        self.lifetime_s = elapsed
        if self.proc:
            self.proc.terminate()
            try:
                self.proc.wait(timeout=20)
            except subprocess.TimeoutExpired:
                self.proc.kill()
                self.proc.wait(timeout=10)

    # ---- HTTP -----------------------------------------------------------
    def erase_slot(self, slot: int = 0) -> float:
        req = urllib.request.Request(
            f"{BASE}/slots/{slot}?action=erase", data=b"", method="POST")
        t = time.perf_counter()
        with urllib.request.urlopen(req, timeout=30) as r:
            r.read()
        return (time.perf_counter() - t) * 1000


def chat(spec: Spec, user: str, system: str | None = None, seed: int | None = None,
         stream: bool = True) -> dict:
    """One request. Returns content, reasoning, timings and token counts."""
    msgs = []
    sys_prompt = system if system is not None else spec.system
    if sys_prompt:
        msgs.append({"role": "system", "content": sys_prompt})
    msgs.append({"role": "user", "content": user})
    body = {
        "messages": msgs,
        "max_tokens": spec.max_tokens,
        **spec.sampling,
        **({"seed": seed} if seed is not None else {}),
        **spec.extra_body,
    }
    if not stream:
        body["stream"] = False
        data = json.dumps(body).encode()
        req = urllib.request.Request(
            f"{BASE}/v1/chat/completions", data=data,
            headers={"Content-Type": "application/json"}, method="POST")
        t = time.perf_counter()
        with urllib.request.urlopen(req, timeout=3600) as r:
            payload = json.loads(r.read())
        wall = time.perf_counter() - t
        ch = payload["choices"][0]
        return {
            "content": ch["message"].get("content") or "",
            "reasoning": ch["message"].get("reasoning_content") or "",
            "finish_reason": ch.get("finish_reason"),
            "usage": payload.get("usage", {}),
            "timings": payload.get("timings", {}),
            "ttft_s": None, "thinking_s": None, "answer_s": None,
            "wall_s": round(wall, 3), "streamed": False,
        }

    body["stream"] = True
    body["stream_options"] = {"include_usage": True}
    data = json.dumps(body).encode()
    req = urllib.request.Request(
        f"{BASE}/v1/chat/completions", data=data,
        headers={"Content-Type": "application/json"}, method="POST")
    content: list[str] = []
    reasoning: list[str] = []
    t0 = time.perf_counter()
    t_first = None
    t_content = None
    finish = None
    usage = {}
    timings = {}
    with urllib.request.urlopen(req, timeout=3600) as r:
        for raw in r:
            line = raw.decode("utf-8", "replace").strip()
            if not line or not line.startswith("data:"):
                continue
            chunk = line[5:].strip()
            if chunk == "[DONE]":
                break
            try:
                obj = json.loads(chunk)
            except json.JSONDecodeError:
                continue
            if obj.get("usage"):
                usage = obj["usage"]
            if obj.get("timings"):
                timings = obj["timings"]
            for ch in obj.get("choices", []):
                if ch.get("finish_reason"):
                    finish = ch["finish_reason"]
                d = ch.get("delta") or {}
                piece_r = d.get("reasoning_content") or ""
                piece_c = d.get("content") or ""
                if piece_r or piece_c:
                    if t_first is None:
                        t_first = time.perf_counter()
                if piece_r:
                    reasoning.append(piece_r)
                if piece_c:
                    if t_content is None:
                        t_content = time.perf_counter()
                    content.append(piece_c)
    t_end = time.perf_counter()
    return {
        "content": "".join(content),
        "reasoning": "".join(reasoning),
        "finish_reason": finish,
        "usage": usage,
        "timings": timings,
        "ttft_s": round(t_first - t0, 3) if t_first else None,
        "thinking_s": round((t_content - t_first), 3) if (t_first and t_content) else None,
        "answer_s": round((t_end - t_content), 3) if t_content else None,
        "wall_s": round(t_end - t0, 3),
        "streamed": True,
    }


def triage_messages(text: str) -> str:
    return text


def run_triage(specs: list[Spec], draws: int, tag: str, prompts: dict | None = None) -> list[dict]:
    out: list[dict] = []
    RESULTS.mkdir(parents=True, exist_ok=True)
    for spec in specs:
        if not spec.path.exists():
            print(f"SKIP missing {spec.path}", flush=True)
            continue
        log = RAW / f"{tag}-{spec.label}" / "server.log"
        print(f"=== {spec.label} ngl={spec.ngl} budget={spec.reasoning_budget} "
              f"preset={spec.preset} ===", flush=True)
        try:
            with Server(spec, log) as srv:
                for draw in range(1, draws + 1):
                    seed = 4100 + draw
                    erase_ms = None
                    try:
                        erase_ms = srv.erase_slot()
                    except Exception:
                        pass
                    r = chat(spec, TRIAGE_PROMPT, seed=seed)
                    rec = {
                        "run_code": _q(), "tag": tag, "label": spec.label,
                        "preset": spec.preset, "notes": spec.notes,
                        "draw": draw, "seed": seed, "file": spec.file,
                        "file_mb": round(spec.path.stat().st_size / 1e6, 1),
                        "ngl": spec.ngl, "ctx": spec.ctx,
                        "sampling": spec.sampling, "max_tokens": spec.max_tokens,
                        "reasoning_budget": spec.reasoning_budget,
                        "reasoning_effort": spec.reasoning_effort,
                        "reasoning_flag": spec.reasoning,
                        "extra_body": spec.extra_body,
                        "system": spec.system,
                        "command": srv._cmd(), "load_s": round(srv.load_s, 2),
                        "erase_ms": round(erase_ms, 2) if erase_ms is not None else None,
                        "content": r["content"], "reasoning_content": r["reasoning"],
                        "finish_reason": r["finish_reason"],
                        "prompt_tokens": r["usage"].get("prompt_tokens"),
                        "completion_tokens": r["usage"].get("completion_tokens"),
                        "ttft_s": r["ttft_s"], "thinking_s": r["thinking_s"],
                        "answer_s": r["answer_s"], "wall_s": r["wall_s"],
                        "prefill_ms": round(r["timings"].get("prompt_ms", 0), 1),
                        "decode_ms": round(r["timings"].get("predicted_ms", 0), 1),
                        "decode_tok_s": round(r["timings"].get("predicted_per_second", 0), 2),
                        "prefill_tok_s": round(r["timings"].get("prompt_per_second", 0), 2),
                        "peak_rss_mb": round(srv.peak["rss"] / 1e6, 1),
                        "peak_vram_mb": srv.peak["vram"],
                        "peak_gpu_util": srv.peak["gpu_util"],
                    }
                    rec["kind"] = "triage"
                    out.append(rec)
                    RESULTS.joinpath(f"{tag}.json").write_text(
                        json.dumps(out, indent=2), encoding="utf-8")
                    print(f"  d{draw} {r['wall_s']}s finish={r['finish_reason']} "
                          f"tok={rec['completion_tokens']} ok={bool(r['content'].strip())}",
                          flush=True)
        except Exception as exc:
            out.append({"run_code": _q(), "tag": tag, "label": spec.label,
                        "kind": "triage", "error": repr(exc)})
            print("FAILED:", exc, flush=True)
            RESULTS.joinpath(f"{tag}.json").write_text(
                json.dumps(out, indent=2), encoding="utf-8")
        time.sleep(1)
    return out


def run_corpus(specs: list[Spec], draws: int, tag: str, corpus_path: Path) -> list[dict]:
    corpus = json.loads(corpus_path.read_text(encoding="utf-8"))["cases"]
    out: list[dict] = []
    RESULTS.mkdir(parents=True, exist_ok=True)
    for spec in specs:
        if not spec.path.exists():
            print(f"SKIP missing {spec.path}", flush=True)
            continue
        log = RAW / f"{tag}-{spec.label}" / "server.log"
        print(f"=== corpus {spec.label} budget={spec.reasoning_budget} ===", flush=True)
        try:
            with Server(spec, log) as srv:
                for case in corpus:
                    for draw in range(1, draws + 1):
                        try:
                            srv.erase_slot()
                        except Exception:
                            pass
                        user = EXTRACTION_INSTRUCTION + "\n\nConversation:\n" + case["conversation"]
                        r = chat(spec, user, seed=4100 + draw)
                        out.append({
                            "run_code": _q(), "tag": tag, "label": spec.label,
                            "case": case["id"], "category": case["category"],
                            "length": case["length"], "draw": draw,
                            "seed": 4100 + draw,
                            "sampling": spec.sampling, "max_tokens": spec.max_tokens,
                            "reasoning_budget": spec.reasoning_budget,
                            "reasoning_effort": spec.reasoning_effort,
                            "ngl": spec.ngl, "file": spec.file,
                            "content": r["content"], "reasoning_content": r["reasoning"],
                            "finish_reason": r["finish_reason"],
                            "prompt_tokens": r["usage"].get("prompt_tokens"),
                            "completion_tokens": r["usage"].get("completion_tokens"),
                            "ttft_s": r["ttft_s"], "thinking_s": r["thinking_s"],
                            "answer_s": r["answer_s"], "wall_s": r["wall_s"],
                            "load_s": round(srv.load_s, 2),
                            "peak_rss_mb": round(srv.peak["rss"] / 1e6, 1),
                            "peak_vram_mb": srv.peak["vram"],
                            "kind": "corpus",
                        })
                        pr = out[-1]
                        print(f"  {case['id']} d{draw} {r['wall_s']:6.1f}s "
                              f"fin={r['finish_reason']} tok={pr['completion_tokens']}",
                              flush=True)
                        RESULTS.joinpath(f"{tag}.json").write_text(
                            json.dumps(out, indent=2), encoding="utf-8")
        except Exception as exc:
            out.append({"run_code": _q(), "tag": tag, "label": spec.label,
                        "kind": "corpus", "error": repr(exc)})
            print("FAILED:", exc, flush=True)
            RESULTS.joinpath(f"{tag}.json").write_text(
                json.dumps(out, indent=2), encoding="utf-8")
        time.sleep(1)
    return out
