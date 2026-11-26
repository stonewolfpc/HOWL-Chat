import json, sys
from pathlib import Path
sys.path.insert(0, str(Path(__file__).resolve().parent))
from harness import ROUND2_MODELS, RAW, Spec, Server, chat, TRIAGE_PROMPT
s = Spec("smoke-qwen35-08b-think", "Qwen3.5-0.8B-Q4_K_M.gguf", ROUND2_MODELS,
         dict(temperature=1.0, top_p=0.95, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0),
         max_tokens=1024, ctx=8192, reasoning_budget=512,
         extra_body={"chat_template_kwargs": {"enable_thinking": True}})
with Server(s, RAW/"smoke"/"server.log") as srv:
    r = chat(s, TRIAGE_PROMPT, seed=4101)
    print(json.dumps({k: v for k, v in r.items() if k not in ("content","reasoning")}, indent=1))
    print("CONTENT:\n", r["content"][:600])
    print("REASONING chars:", len(r["reasoning"]))
    print("REASONING tail:", r["reasoning"][-300:])
