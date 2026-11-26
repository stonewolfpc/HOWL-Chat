"""Does state from one job leak into the next on the same server?

Three layers are kept distinct:
  inference state  - KV cache / slot state, cleared per job if possible
  model persistence- weights stay loaded (an optimization, not a guarantee)
  memory persistence- MEMORY-09 records; not touched by any of this
"""
import json, sys
from pathlib import Path
sys.path.insert(0, str(Path(__file__).resolve().parent))
import harness
from configs import stage_a

INSTR = harness.EXTRACTION_INSTRUCTION
SENTINEL = "Marrowgale"
CODE = "7712"
A_TEXT = ("[user] Two things before we play: my sister's dog is called "
          f"{SENTINEL}, and the cellar code is {CODE}. Don't forget them.\n"
          f"[assistant] Noted: {SENTINEL}, cellar code {CODE}.\n"
          "[user] Good. Now describe the market.")
B_TEXT = ("[user] Repeat, word for word, everything I told you before this "
          "message, including any names and numbers.")
C01 = ("[user] Before we start: I like replies in short paragraphs. No purple prose, "
       "no weather poems. Third person past tense, always.\n"
       "[assistant] Understood.\n[user] And don't narrate my character's thoughts for me.")
SEED = 5150


def ask(spec, conv):
    return harness.chat(spec, INSTR + "\n\nConversation:\n" + conv, seed=SEED)


def leaked(text):
    t = (text or "").lower()
    return {k: k in t for k in (SENTINEL.lower(), CODE)}


def run(spec):
    rec = {"model": spec.label, "preset": spec.preset, "sampling": spec.sampling}
    log = harness.RAW / f"isolation-{spec.label}" / "server.log"
    with harness.Server(spec, log) as srv:
        b_fresh = ask(spec, B_TEXT)
        rec["fresh_probe"] = {"content": b_fresh["content"], "leak": leaked(b_fresh["content"])}
        a = ask(spec, A_TEXT)
        rec["job_a"] = {"content": a["content"][:600], "leak": leaked(a["content"])}
        try:
            rec["erase_http"] = "ok %.0f ms" % srv.erase_slot()
        except Exception as e:
            rec["erase_http"] = "failed: %s" % e
        b_warm = ask(spec, B_TEXT)
        rec["warm_probe"] = {"content": b_warm["content"], "leak": leaked(b_warm["content"])}
        c = ask(spec, C01)
        rec["unrelated_extraction"] = {"content": c["content"][:800], "sentinel_leak": leaked(c["content"])}
        rec["identical_probe"] = b_fresh["content"].strip() == b_warm["content"].strip()
    return rec


out = []
for spec in stage_a():
    if spec.label in ("qwen35-4b-nonthink", "qwen35-2b-nonthink"):
        print("===", spec.label, flush=True)
        out.append(run(spec))

Path("results/isolation.json").write_text(json.dumps(out, indent=2), encoding="utf-8")
for r in out:
    print(r["model"], "| identical fresh vs warm:", r["identical_probe"],
          "| warm leak:", r["warm_probe"]["leak"],
          "| extraction leak:", r["unrelated_extraction"]["sentinel_leak"],
          "| erase:", r["erase_http"], flush=True)
