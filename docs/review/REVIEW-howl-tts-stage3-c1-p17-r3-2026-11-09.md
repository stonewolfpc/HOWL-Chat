# P17 — Fourth-Stage Review of the r3 Manifest Correction

**Date:** 2026-11-09
**Status:** REVIEW OF RECORD. **Preserved verbatim, unedited.**
**Subject:** `RESPONSE-howl-tts-stage3-c1-p17-review-r3-2026-11-09.md`
**Disposition:** Scoped corrections required; verdict otherwise accepted.
**C1 IMPLEMENTATION REMAINS UNAUTHORIZED.**

The body below is transcribed without modification from the fourth-stage review delivered on 2026-11-09. Only this container header was added.

<!-- BEGIN VERBATIM REVIEW -->

This review is internally consistent, and I agree with the verdict. It found three substantive issues rather than reopening already corrected source facts.

The most important distinction is row-level association versus document-level presence. The §6 check passing did not establish that every finding was attached to every applicable emission row. O14 demonstrates exactly why that distinction matters.

My assessment:

1. O14: Add it to S59, S61, and S66 if O14 represents Synthesize cancellation generally. The source behaviors described are consistent with that interpretation.
2. M08: Correct `client.go:294` to `client.go:295`. This is a bad anchor, not a reason to change the mechanism inventory.
3. Branch count: Retain 61 only as the reproducible count of unique function/condition pairs. Don't call those verified branch identities. A proper branch count needs stable branch identifiers or source locations.
4. §6 verification: Restrict association checks to the actual ID cells in the relevant tables. Prose mentioning an ID must not satisfy a coverage check.

The rest of the disposition looks appropriately restrained:

- 68 emission anchors remain reproducible.
- O17's exclusion is explicit and consistent with the narrowed scope.
- E7–E9, E4, and O2/O3/O6 pass this review.
- RC-01 remains partial; RC-03/04 remain partial.
- R01–R13 and U1–U8 retain their existing dispositions.

I would keep the next author correction tightly scoped. Fix those O14 associations, M08, the branch-count label or definition, and the row-level verifier. Then request another independent coverage pass. Don't revise P17 or the lifecycle contract to accommodate manifest bookkeeping.

And importantly, this is still documentation/review progress—not C1 implementation authorization.

<!-- END VERBATIM REVIEW -->
