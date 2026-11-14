# MEMORY-09 — S1 Execution Preflight (read-only)

**Status:** **PASS WITH FLAGS.** 11 pass, 0 fail, 3 flags. The flags are all one item: `S1-A` has no qualification record.
**Date:** 2026-11-14 · **Mode:** read-only. No server started, no model loaded, no prompt sent, no baseline change.
**Evidence:** `eval/runs/PREFLIGHT-S1-2026-11-14T135649Z/preflight.json` · re-runnable: `python -m eval.harness.cli preflight-s1`

---

## 1. The seven confirmations requested

| # | Question | Result |
|---|---|---|
| 1 | Each qualified model's digest and byte count against the manifest | **PASS** — all three files re-hashed on disk, not merely trusted from records |
| 2 | Context and GPU-layer settings match the frozen configuration | **PASS for B and C** · **FLAG for A** |
| 3 | Recorded offload configuration is what S1 will use | **PASS for B and C** · **FLAG for A** |
| 4 | Qualification directories contain zero generation records | **PASS** — sum of `generations` across filed qualification records = 0 |
| 5 | Recovery baselines and abort thresholds carried into the procedure | **PASS** — 78 °C / 400 MiB / 8 GB / 250 MiB tolerance, per-component |
| 6 | 18-call schedule and 30-call ceiling unchanged | **PASS** — planned 18, retry 12, hard 30 |
| 7 | Truncation is data, no reruns | **PASS** — rule stated in both documents; adapter has no retry loop |

**Verdict: no failure. Three flags, all against `S1-A`.**

---

## 2. The flags, precisely

`S1-A` (Qwen3.5-9B) was never *qualified* — it was *demonstrated* at L3, which is a different act with a different output.

| Flag | What the evidence says | Why it matters for S1 |
|---|---|---|
| Runtime configuration | L3 ran at **context 20480**, and its record **does not record GPU layer count** at all | `S1-B` and `S1-C` are qualified at **8192**. Running A at 20480 would mean three slots with three different context settings, only one of which is documented |
| Offload evidence | No server log exists for A with the `offloaded N/M layers` line | The other two slots have their actual split (36/43, 20/25). A's is unknown, so "the same configuration" cannot be asserted for it |
| Recovery gate | No `recovery.json` for A | B and C each returned to baseline with Δ 0 MiB. A's release was never measured in this form |

Note what the L3 record *does* establish: the file hash (re-confirmed here on disk), one successful generation, and clean termination. It does **not** establish that A is configured as S1 will configure it.

**The zero-generation check is not contradicted:** it counts qualification records, and A has none. A's single generation belongs to L3, which was authorised as a generation.

---

## 3. What this does not license

The preflight does not authorise S1. It confirms that the *evidence* is internally consistent for two of three slots, and that the procedure, ceiling and truncation policy are intact. It cannot confirm that six consecutive calls per model are safe, that latency is comparable across slots, or that `max_tokens=256` behaves alike on three architectures.

Agreed and recorded: **Gemma's CPU-heavy residency and LFM's lower headroom are not reasons to tune the experiment.** They are properties of the population and go into the manifest as observations, not adjustments.

---

## 4. Closing the gap

One zero-generation load of `S1-A` at the S1 configuration (context 8192, its own layer count, `q8_0` KV) would produce the same three artefacts B and C have — configuration, offload line, recovery file — and clear all three flags. It sends no prompt and cannot produce a MEMORY-09 observation.

Two honest options:

| Option | Consequence |
|---|---|
| **Qualify A at 8192** (recommended) | All three slots uniform and equally documented; one load cycle, zero generations |
| **Run A at 20480 as in L3** | Defensible only if the differing context is recorded per slot and no cross-slot comparability is claimed for context-sensitive behaviour; leaves the offload and recovery gaps open |

The second option is not wrong, but it is weaker, and it would be weaker for a reason unrelated to the experiment.

---

## 5. A correction made to the checker during this pass

The first run reported a **FAIL** on "truncation is data, no reruns". The rule is stated plainly in the plan ("never retried") and in the readiness record ("No retry of a truncated reply"), but the check searched one document for one exact phrase. The check was wrong, not the procedure — it now reads both documents and accepts the rule where it is stated, and still fails if the adapter ever grows a retry loop.

Recorded here because a checker that fails on phrasing is as misleading as one that passes on nothing.

---

## 6. Bottom line

**Two slots are ready with evidence; one slot is ready on hash and generation but not on documented configuration.** That is the entire gap between here and an authorisable S1. It is one load cycle wide.
