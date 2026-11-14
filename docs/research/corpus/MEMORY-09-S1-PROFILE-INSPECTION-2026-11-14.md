# MEMORY-09 — S1 Operating Profile vs Frozen Configuration (inspection)

**Status:** **INSPECTION ONLY — no load, no generation, no tuning, no baseline edit.**
**Date:** 2026-11-14 · **Trigger:** owner's established operating practice (20 layers / 20K context, VRAM < 90%).
**Evidence:** read-only GGUF metadata + filed qualification records + `nvidia-smi` baseline.

---

## 1. The practice being respected

| Model class | GPU layers | Context |
|---|---|---|
| 9B | **20** | **20,480** |
| 12B | **20** | **20,480** |

Target: **VRAM used below 90%** of the card, so the machine stays usable — SillyTavern, Spotify, a browser — while a model runs. This is the profile already used during Pure Byte experiments with SillyTavern running alongside.

**On this card:** 90% of 8,192 MiB = **7,373 MiB used**. Desktop baseline measured at load time: **1,229 MiB**. So the profile's real budget is ~6,144 MiB for model + KV.

**This is a different objective from maximum offload.** The qualification runs were aiming at "loads and identifies", not "leaves the machine usable". They answered the first question correctly and were never asked the second.

---

## 2. What is frozen or filed today

| Slot | Config in filed evidence | Context | Source |
|---|---|---|---|
| `S1-A` | **30 layers** (launch args), *not recorded in the record* | 20,480 | L3 run, not a qualification |
| `S1-B` | **36** layers of 43 | 8,192 | `SLOT-QUAL-S1-B-…135031Z` |
| `S1-C` | **20** layers of 25 | 8,192 | `SLOT-QUAL-S1-C-…135147Z` |

---

## 3. Discrepancies against the profile

| Slot | Layers: filed → profile | Context: filed → profile | Discrepancy |
|---|---|---|---|
| `S1-A` | 30 → **20** | 20,480 → 20,480 | layer count differs; no recorded config to compare against |
| `S1-B` | 36 → **20** | 8,192 → **20,480** | both differ |
| `S1-C` | 20 → 20 | 8,192 → **20,480** | context differs |

**Nothing has been changed.** §10.2 of the readiness record freezes per-model runtime configuration at qualification, so adopting the profile is an **explicit versioned decision**, not an adjustment made in place. The filed qualification records stay valid as evidence of what those loads actually used.

---

## 4. Does the profile fit? (estimate, not measurement)

Method: parse GGUF attention metadata → KV bytes per token; subtract measured KV from the observed post-load usage to infer GPU-resident weight bytes; scale for 20 layers; add back 20K KV (+6% for `q8_0` scale overhead). **Estimates. Only a load confirms them.**

| Slot | KV @8K | KV @20K | Est. total at 20 layers / 20K | % of card | vs 90% (7,373) |
|---|---|---|---|---|---|
| `S1-A` | 512 MiB | **1,280 MiB** | ~5,295 MiB | 65% | **fits**, 2.1 GiB spare |
| `S1-B` | 672 MiB | **1,680 MiB** | ~4,411 MiB | 54% | **fits**, 2.9 GiB spare |
| `S1-C` | 96–192 MiB* | 240–480 MiB* | ~6,000–6,300 MiB | 73–77% | **fits**, 1.1–1.4 GiB spare |

\* `S1-C` is MoE; its KV head count is a per-layer array, so the range is a bound rather than a point estimate.

**All three fit under 90% with 1.1–2.9 GiB spare** — which is the margin that makes the profile worth having. Note that even the *qualification* configurations were under 90% (A 80%, B 54%, C 72%); the profile is about keeping that true with other applications running, not about avoiding a hard limit.

---

## 5. What the profile costs

Fewer layers on the GPU means more weights in CPU-mapped memory, so generation is slower — particularly `S1-B`, which already kept 3,173 MiB on the CPU at 36 layers and will keep more at 20. That is an accepted cost of a machine that stays usable, and it is a property of the **runtime**, not of the experiment. It must be recorded per slot before any latency comparison across slots is attempted, and it must not become a reason to raise layers mid-batch.

---

## 6. Two gates, not one

| Gate | Value | Status |
|---|---|---|
| **Operating target** (owner) | used < 90% of VRAM (7,373 MiB) | profile chosen to satisfy it |
| **Hard safety abort** | free < 400 MiB | unchanged, separate, not revised by this document |

The 90% target is how the machine is used. The 400 MiB floor is where the run stops. They answer different questions and both remain in force.

---

## 7. Decision required (versioned)

Adopting the profile means amending the frozen runtime configuration, so it needs a decision recorded as a version, not an edit:

1. **`S1-v1.1 — runtime profile`**: 20 GPU layers and 20,480 context for all three slots; VRAM < 90% as an operating target; 400 MiB abort unchanged. This supersedes §10.2's "config fixed at qualification" for *configuration*, while leaving prompts, fixture, state, params, repetitions and the 18/30 ceiling untouched.
2. **Re-qualification at the profile is then required** — the filed records describe 36/8192 and 20/8192, so none of them proves the profile loads. Three zero-generation load cycles would produce config + offload + recovery evidence at the profile, and a re-run of the preflight would return without flags.
3. If the profile is instead adopted **without** re-qualification, the preflight flags stay open and S1 would run on an unverified configuration. That is a legitimate owner choice, but it should be a stated one.

**Recommendation:** version the profile, re-qualify all three at it (zero generations), re-run `preflight-s1`, then decide on S1. The whole sequence is three load cycles and no prompts.

---

## 8. Boundaries observed

No model was loaded, no server started, no prompt sent, no threshold revised, and `eval/s1-readiness-baseline.json` is untouched (`ad42f28d…`). Everything above is either read from filed records, parsed from file headers, or labelled an estimate.
