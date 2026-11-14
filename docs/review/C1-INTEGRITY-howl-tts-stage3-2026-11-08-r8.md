# C1 Integrity Record — Revision 8

**Date:** 2026-11-08
**Subject:** Howl TTS Stage 3 — Capability C1, design revision 8
**Status:** documentation revision only; **C1 UNAUTHORIZED**

Convention: an integrity record does not contain its own hash. The SHA-256 of this file is recorded in
the revision report and must be captured immediately after issue.

## 1. Artifacts produced by this revision

| Artifact | SHA-256 | Lines | Bytes |
|---|---|---|---|
| `DESIGN-howl-tts-stage3-c1-revision-2026-11-08-r8.md` | `a8c40fac7d86f3cd34ad34e7f656668684830e76dfb605ae6817fa7119638e06` | 391 | 30,637 |
| `C1-READINESS-howl-tts-stage3-2026-11-08-r8.md` | `53c8751dcaaeaa92583041db26315f20f0cfbfa7dbce508796b259c05beaa40d` | 103 | 6,130 |
| `C1-INTEGRITY-howl-tts-stage3-2026-11-08-r8.md` | not self-hashed, by convention | — | — |

Manual edits were rejected after writes were verified; the design hash above is the post-edit hash and
is the one embedded in the readiness document.

## 2. Preserved inputs, re-hashed in this revision

| Artifact | SHA-256 | Lines | Bytes | Relation |
|---|---|---|---|---|
| `DESIGN-howl-tts-helper-architecture-2026-11-06.md` | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f` | 2586 | 182,904 | governing contract; unchanged |
| `REVIEW-howl-tts-stage3-c1-r7-2026-11-08.md` | `c025ba0a32c776996fd90a59cca0635da5532a43200a5f6b8dad39ff91730f9d` | 200 | 13,835 | the review r8 answers; unchanged |
| `DESIGN-howl-tts-stage3-c1-revision-2026-11-08-r7.md` | `0867acb7216a0e8ee2802b6f08275d997e4138f089ae51acac373c10eea4af54` | 1069 | 161,934 | normative for unchanged sections; unchanged |
| `C1-READINESS-howl-tts-stage3-2026-11-08-r7.md` | `8a42cf730a2ce35fad338c05746a36495fbf7ee52bb544bf5ff1e0ed4d29dfb2` | 127 | 14,684 | superseded by the r8 readiness pair only |
| `C1-INTEGRITY-howl-tts-stage3-2026-11-08-r7.md` | `65669ce771a5f173fb0c982c3cf3c86b723cf4a72dee3042646eb32c15d32f2b` | 84 | 7,368 | superseded by this record |
| `REVIEW-howl-tts-stage3-c1-r6-2026-11-08.md` | `9176f21512698a3d15ea073b84d39ae69c88fdcb8ef92ca37f0c99288889f0c9` | 252 | 28,668 | historical; unchanged |

All six match the values recorded during the r7 verification set `[Observed]`. No preserved artifact was
modified by this revision; the r7 documents remain on disk so the r6→r7→r8 chain is auditable.

The r9 requirements, owner-decision record, and D10 finalization hashes stated in design §1 are carried
from the r7 verification set; they were not re-derived here and are not re-asserted by this record.
`[Unestablished]` in this revision.

## 3. Supersession map

| r7 locus | r8 treatment | r8 locus |
|---|---|---|
| §6.13 update reset | **replaced in full** | §4 |
| §6.12 catalog selection | amended (carrier added) | §5 |
| §6.8 availability predicates | amended (three classes, A2-D5 citation) | §6 |
| §6.6 baseline labels | corrected, rows unchanged | §7 |
| §10 criteria 27, 31, 32 | replaced; criterion 33 added | §8 |
| §10.1 rows 27, 31, 32 | replaced; one row added | §9 |
| §11 test matrix | extended by 11 rows | §10 |
| everything else | carried, unmodified, at `0867acb7…af54` | §12 |

## 4. Prohibited-change findings

No Go source, implementation test, dependency, configuration, build, or runtime artifact was touched.
The pre-existing dirty working tree — including the 26 entries under `internal/audio` and
`app-audio.go` — was left exactly as found. Measured working-tree counts are in §7.

## 5. What this record does and does not establish

It establishes that the three r8 files exist with the stated identities, that preserved inputs are
byte-identical to their recorded values, and that the structural counts in §7 hold for the r8 text.

It does not establish that the design is correct, that C1-AC-15 or `MaxPending=32` is resolved, that any
test passes, or that implementation is authorised. Those remain open in the readiness document.

## 6. Limitations

- r8 is a delta document. Its normative corpus is r8 plus r7 §1–§12 exclusive of the sections in §3
  above, identified by hash. A reviewer needing one file must request the consolidated r9.
- All hashes are of files on disk at write time; no signature or timestamp authority is involved.
- Carried r9/owner/D10 identities are not re-derived here (see §2).

## 7. Validation log (appended at write time)

Measured 2026-11-08T20:27:58Z on the files as they exist on disk.

| Check | Result | Class |
|---|---|---|
| Distinct C1-R identifiers in the design | 37 | [Calculated] |
| Findings register rows | 37 | [Calculated] |
| Criteria accounted for in the §8.1 locus map | 26 + 3 + 4 = 33 (declared arithmetic, 6 map rows) | [Calculated] |
| Verification rows | 29 (r7 §10.1 carried) + 4 (r8 §9) = 33 | [Calculated] |
| Duplicate `## N.` headings in design | 0 | [Observed] |
| Duplicate `### N.N` headings in design | 0 | [Observed] |
| Exported ticket fields / `Approved` design surface | none; occurrences are historical references only | [Observed] |
| `UpdateDecision`, `UpdateGrant`, rereintroduced setters | absent as design surface | [Observed] |
| Readiness embeds the actual design hash | yes (1 match) | [Observed] |
| Design §12 embeds the actual r7 design hash | 2 match(es) | [Observed] |
| Implementation identifiers present in Go source | 0 files | [Observed] |
| Dirty tracked `.go` files in the working tree | 110 — all pre-existing, none touched by this revision | [Observed] |
| Dirty entries under internal/audio and app-audio.go | 26 (unchanged, pre-existing) | [Observed] |
| Total dirty working-tree entries | 142 | [Observed] |
| `docs/review/` presence in git index | 0 files tracked; the r6/r7/r8 records remain untracked as found | [Observed] |
| Commits, stages, pushes made | none | [Observed] |

Evidence classes: `[Observed]` read from disk, `[Calculated]` arithmetic or count.
