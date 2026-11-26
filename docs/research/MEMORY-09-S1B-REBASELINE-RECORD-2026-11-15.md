# S1-B Re-baseline Record — prepared, not applied

**Identity:** `MEMORY-09-S1B-REBASELINE-1` · **Date:** 2026-11-15
**Status:** **PREPARED. NOT APPLIED. NOT COMMITTED.**
**Proposal:** option **B-1** — re-baseline `S1-B` to `d0027dd3…6e96` as `S1-v1.2`.

This document is the "explicit record of why the baseline changed" that B-1 was made conditional on. It records the intended change; it does not make it. `eval/s1-readiness-baseline.json` and `eval/s1-readiness-baseline-v1.1.json` are byte-identical to their state before this record was written.

---

## 1. Three claims, kept separate

B-1 rests on the middle one only. The other two are not upgraded by it.

| Level | Claim | Status | Evidence |
|---|---|---|---|
| **L-1** | `d0027dd3…6e96` is the content hash the **repository serves** for that exact filename and byte size | **repository-verified reference** | Git LFS pointer; content-addressed, machine-generated at upload. Attests to the served blob — **not** to the artifact's construction, and **not** to the history of the local copy |
| **L-2** | the **current local file** returns that digest repeatably, through both tested read paths | **demonstrated 2026-11-15** | `MEMORY-09-S1B-STABILITY-RESULT-2026-11-15.md`: 10 reads, 10 identical digests, one file identity across 20 captures, zero short reads, complete byte accounting |
| **L-3** | **why** the 13–14 November reads differed | **unresolved** | seven non-matching digests, cause unidentified; `341fb28c…` no longer reproducible |

**L-1 + L-2 do not produce L-3.** Nothing in this record explains the historical discrepancy, and nothing here diagnoses memory, controller, driver or storage health. A healthy result today is not a clean bill of health for the hardware, and the earlier disagreement is not evidence of failing hardware either.

**Wording standard:** "repository-verified reference", never "publisher-verified provenance". The pointer authenticates bytes; it says nothing about who built the model or in what state their local copy sat.

---

## 2. What the record preserves

| Preserved | Where it survives |
|---|---|
| The **original** manifest digest `b5b1913e44aa64d2c12b6159bb782774dd98362528211ac830860595d293ede0` | `s1-readiness-baseline-v1.1.json` stays on disk **untouched**; v1.2 supersedes rather than edits it; §4 below repeats the value |
| The **ten-read probe** and all raw evidence | `eval/runs/S1B-STAB-20261115T000423Z/` (10 × `read-nn.json`, `summary.tsv`), driver `eval/runs/probe-s1b-stability.py`, log `eval/runs/probe-s1b-stability.log` |
| The **repository reference** and its verification level | §1 L-1, and `WO-STORAGE-EXT-01-EXTERNAL-HASH-RESEARCH-2026-11-14.md` §2.6 (`cd7366fe…`), which already carries the correct class label |
| The **historical discrepancy as unresolved** | §1 L-3, `MEMORY-09-INTEGRITY-STOP-S1B-DIGEST-2026-11-14.md`, and the errata file `6ce6ea0e…` |
| **No S1-B behavioural generation has occurred** | qualification `generations: 0` for all three slots; the single spent generation (L3-01) was `S1-A`, temperature 0.0, 48 tokens |

---

## 3. The exact delta

Applied to a copy of `s1-readiness-baseline-v1.1.json`, producing `s1-readiness-baseline-v1.2.json`. **Four additions, one replacement.** Nothing else changes.

```diff
- "baseline_id": "S1-readiness-baseline-1.1",
- "version": "1.1",
- "frozen_at": "2026-11-14",
- "supersedes": "S1-readiness-baseline-1",
+ "baseline_id": "S1-readiness-baseline-1.2",
+ "version": "1.2",
+ "frozen_at": "2026-11-15",
+ "supersedes": "S1-readiness-baseline-1.1",
+ "rebaseline": {
+   "slot": "S1-B",
+   "field": "models[1].sha256",
+   "from": "b5b1913e44aa64d2c12b6159bb782774dd98362528211ac830860595d293ede0",
+   "to": "d0027dd3a9128d9323e9f282c8bf010a8526c46477584535991dc1a869b56e96",
+   "reason": "The superseded value was recorded on 2026-11-14 from a buffered read whose result is no longer reproducible. Ten independent reads on 2026-11-15 across two read paths returned the new value with one file identity, and that value equals the repository-hosted LFS pointer. The original value was never a reference-independent truth; it was one reading.",
+   "evidence": [
+     {"doc": "MEMORY-09-S1B-STABILITY-RESULT-2026-11-15.md", "sha256": "2a9acb6ec307d5a0fd514de4d122828dde557b12b87bd015c8fcac758990ac30", "claim": "L-2"},
+     {"doc": "S1B-READ-STABILITY-PROBE-2026-11-14.md", "sha256": "4266500135e25f6712a16e21f9adf8e41218817325d95ac48c3b831befcb42e1", "claim": "pre-registration"},
+     {"runs": "eval/runs/S1B-STAB-20261115T000423Z/", "reads": 10, "distinct_digests": 1}
+   ],
+   "reference": {"class": "repository-verified", "source": "git-lfs pointer, HauhauCS/Gemma-4-E4B-Uncensored-HauhauCS-Aggressive", "publisher_verified": false},
+   "unresolved": "No explanation for the six non-matching digests or for the disappearance of the 341fb28c buffered stream. Not attributable to any component.",
+   "generations_at_rebaseline": 0,
+   "rebaseline_does_not_establish": ["history of the local copy", "cause of the earlier digests", "hardware health"]
+ }
```

Field-by-field, unchanged and confirmed against `s1-readiness-baseline-v1.1.json`:

`runtime_profile` (20 layers, 20,480 ctx, 90 %, 400 MiB floor, 250 MiB recovery tolerance, 74/78 °C, `q8_0` KV, 8 threads, batch 512, `127.0.0.1`) · `fixture` `3dd41c8e…` · `state_sha256` `aa98d95f…` · `prompts` `076de19e…` / `97213346…` · `params` 0.8 / 0.95 / 256 · `schedule` 18 / 12 / 30 / 3 reps · `plan` `d15ebace…` 12,074 B · `models[0]` and `models[2]` untouched · `preset_reference` `9ed540a5…`

**Two fields beyond that diff, disclosed rather than folded in:** `note` and the three `models[*].status` strings name the profile version, so they were updated to `S1-v1.2`. Left at `v1.1` the manifest would misdescribe itself — the exact failure mode this project keeps correcting. No check reads either field, nothing normative lives in them, and reverting them changes no behaviour.

---

## 4. Applying is not inert — how activation actually works

`eval/harness/preflight.py` selects a manifest by **highest `version`** among `s1-readiness-baseline*.json`:

```python
candidates = sorted(MANIFEST_DIR.glob("s1-readiness-baseline*.json"))
return max(candidates, key=version_key)   # newest version wins, older files stay
```

So **creating `s1-readiness-baseline-v1.2.json` immediately becomes the manifest preflight uses** — no other edit is required for it to take effect. That is why the file was not written as part of preparing this record.

**That gap is now closed.** The check previously pinned `eval/s1-readiness-baseline.json` — the v1 file — so a newly active manifest would have run unguarded. It now:

- globs `eval/s1-readiness-baseline*.json` and applies every existing assertion to **every** manifest on disk (full 64-char digests, fixture/state/prompt re-derivation, frozen params, plan bytes and digest, three named slots, zero generations);
- resolves the active manifest through `preflight.manifest_path()` itself, so a future v1.3 is covered by adding a file rather than by editing the check again;
- walks the supersession chain: a predecessor named in `supersedes` must still be on disk, and a manifest that changes any model digest **must** carry a `rebaseline` block naming that slot, the digest replaced, the digest installed, and zero generations.

Still **89** checks — the guard was extended in place, so `offline_checks_at_freeze` remains accurate and no field in v1.2 diverges.

**Rollback:** delete `s1-readiness-baseline-v1.2.json`; `max()` falls back to v1.1. No history is destroyed, because v1 and v1.1 are never edited.

---

## 5. What applying would not do

No qualification run. No model load. No server. No generation — budget stays at **1 spent of 30**. No change to thresholds, prompts, fixture, state, plan or schedule. No claim that S1-B is usable: after a re-baseline it still needs a qualification run at the v1.1/v1.2 profile, and `gemma4` has **never been loaded on this build** (42 layers, 20 offloaded, 22 on CPU).

---

## 6. Disposition — applied 2026-11-15

| Step | Result |
|---|---|
| `eval/s1-readiness-baseline-v1.2.json` created | 5,391 B, `c983f44713622d86a827cb59ff18db30139031e39184057d37fd2a8cd11f73cd` |
| `s1-readiness-baseline.json` (v1) | untouched — `ad42f28d9402e34624e0b3bf4bd7fed169099a188174b9c12accf5b99f7b05fc` |
| `s1-readiness-baseline-v1.1.json` | untouched — `14a249bce1673cf7186b963566fb71f2b4314e75e1abb31253a1f3a136129981` |
| `preflight.manifest_path()` | selects **v1.2** |
| `python -m eval.harness.validate` | **89/89**, count unchanged |
| `python docs/research/corpus/check_corpus.py` | **PASS** |
| `eval/harness/validate.py` (guard change) | 66,133 B, `6a957f700e54b66fcafcff7ea43b33bce9ec1541cf9271584fed7dcd7e5eb52b` |
| model load / server / generation | none — budget still **1 of 30** |

Three negative controls, each applied to a temporary edit and each restored **byte-exactly** (`c983f447…` before and after):

| Control | Expected | Observed |
|---|---|---|
| corrupt `models[1].sha256` | fail | `rebaseline.to d0027dd3… is not the digest now installed` |
| delete the `rebaseline` block | fail | `S1-readiness-baseline-1.2 changes ['S1-B'] but records no rebaseline reason` |
| remove v1.1 from disk | fail | `supersedes S1-readiness-baseline-1.1, which is not on disk` |
| (all three together) | no effect on the real file | digests above re-verified after restore |

The second control **initially passed**, showing the record requirement was optional: a manifest could have silently dropped the reason-of-change. That hole was closed before this section was written. The first run of the extended check also failed on real data — `KeyError: 'version'`, because the v1 manifest has no `version` field — which is why the guard now mirrors preflight's ordering with a default instead of assuming the key.
