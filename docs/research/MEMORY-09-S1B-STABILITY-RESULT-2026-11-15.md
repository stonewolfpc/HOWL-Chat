# S1-B Read-Stability Probe — Result

**Identity:** `MEMORY-09-S1B-STAB-1` · **Executed:** 2026-11-15T00:04:23Z · **Authorised by:** owner ("run the S1-B stability probe")
**Specification:** `docs/research/S1B-READ-STABILITY-PROBE-2026-11-14.md` (`42665001…42e1`), unchanged before execution.
**Evidence:** `eval/runs/S1B-STAB-20261115T000423Z/` (10 × `read-nn.json` + `summary.tsv`) · log `eval/runs/probe-s1b-stability.log` · driver `eval/runs/probe-s1b-stability.py`.

---

## 1. Outcome

**Row 1 of the pre-registered decision table: all `B` equal, all `C` equal, `B = C` — stable and path-independent, in this session.**

Ten reads, ten identical digests:

```
d0027dd3a9128d9323e9f282c8bf010a8526c46477584535991dc1a869b56e96
```

| Read | Kind | Path | Digest | Elapsed |
|---:|---|---|---|---:|
| 1 | B | buffered (`sha256sum`) | `d0027dd3…6e96` | 48.419 s |
| 2 | C | unbuffered body + buffered tail | `d0027dd3…6e96` | 49.084 s |
| 3 | B | buffered | `d0027dd3…6e96` | 48.946 s |
| 4 | C | composite | `d0027dd3…6e96` | 48.717 s |
| 5 | B | buffered | `d0027dd3…6e96` | 48.738 s |
| 6 | C | composite | `d0027dd3…6e96` | 48.740 s |
| — | — | **60 s idle** | — | — |
| 7 | B | buffered | `d0027dd3…6e96` | 48.689 s |
| 8 | C | composite | `d0027dd3…6e96` | 48.749 s |
| 9 | B | buffered | `d0027dd3…6e96` | 48.236 s |
| 10 | C | composite | `d0027dd3…6e96` | 48.697 s |

**Distinct digests: 1.** No stop condition fired. Exit 0, total ≈ 9 min, ≈53.4 GB of I/O.

## 2. Identity, pinned at all 20 capture points

| Field | Value | Distinct values across 10 reads × before/after |
|---|---|---|
| size | 5,335,285,728 B | **1** |
| mtime_ns | 1788573228448965700 (2026-09-04) | **1** |
| volume serial | `000000000ecb9918` | **1** |
| file id | `0003000000014f91` | **1** |
| in use | false at every capture | 1 |

One file identity throughout. No write occurred; `mtime` is unchanged from the 4 September value recorded before this investigation.

## 3. Byte accounting for the composite reads

All five composite reads reported identical structure:

| Field | Value |
|---|---|
| `unbuffered_bytes` | 5,335,281,664 (1,273 × 4 MiB sectors) |
| `buffered_bytes` / `tail_len` | 4,064 |
| `read_calls` | 1,273 |
| `short_sector_aligned_reads` | 0 |
| IO delta | 1,274 read ops, 5,335,285,728 bytes = file size exactly |

Sum of unbuffered + buffered = 5,335,285,728 = size. No byte unaccounted for, no byte read twice.

## 4. What this establishes

1. **The input is stable today.** Same bytes, same identity, both paths, across a 9-minute window including a 60 s idle gap.
2. **The path dependence recorded on 14 November does not reproduce.** That day a buffered read returned `341fb28c…` while a cache-bypassing read returned `d0027dd3…6e96`. Today both paths return `d0027dd3…6e96`.
3. **`b5b1913e…ede0` — the manifest's expected hash — was not produced by any of the ten reads.** 0/10.
4. **`d0027dd3…6e96` — the content hash in the repository-hosted LFS pointer — was produced by all ten.** 10/10, and it equals the size in that pointer (5,335,285,728 B). That pointer is a **repository-verified reference**: it attests to the bytes the repository serves for that exact filename, and to nothing about the history of the local copy.

## 5. What this does not establish

- **Why the old buffered stream disappeared.** Cache eviction, a reboot, or a state change in the observing stack are all consistent with the data; none was measured. No cache query was made and no purge was attempted.
- **What the file contained on 13–14 November.** Today's stability says nothing retroactively.
- **That `b5b1913e` was a bad reading.** It was not reproduced today; the cause of the error is still unidentified.
- **That the file is stable for any reader.** Ten reads from two paths in one session on one machine. Independence was independent *process start*, not independent cache state — the system cache was live throughout.

**Timing note, marked as inference:** buffered and composite reads both ran at ≈110 MB/s (5.34 GB / 48.7 s). A cache-served buffered read would be expected to be substantially faster. That is *consistent with* the buffered reads being served from the volume rather than from resident pages, but it is a throughput observation, not a cache measurement, and it was not verified.

## 6. Consequence for the decision — not the decision

The owner's objection was that re-baselining would calibrate the experiment around an input whose stability was unknown. **That objection is now answered by measurement rather than argument: the input is stable, and its digest equals the repository-verified reference.** The concern itself was right, and the probe was the correct response to it.

That makes option **B-1** (re-baseline `S1-B` to `d0027dd3…6e96`) evidence-backed rather than premature — but **this report does not adopt it**. Whether to re-baseline, keep the blocked hash, or drop the slot remains the owner's decision, as does the version label (`S1-v1.2` under B-1).

Unchanged and untouched: `eval/s1-readiness-baseline.json` (`ad42f28d…05fc`) and `eval/s1-readiness-baseline-v1.1.json` (`14a249bc…9981`). No qualification run. No model load. No server. No generation — budget still 1 spent.

## 7. Machine state after the probe

GTX 1080: 1,735 MiB used / 8,192 total (21.2 %), 6,325 MiB free. No KoboldCpp, no listener on 5001. Probe was CPU and disk only.
