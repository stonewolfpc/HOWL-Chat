# P17 — Second-Stage Review of the Reconciliation

**Date:** 2026-11-09
**Status:** REVIEW OF RECORD. **Preserved verbatim, unedited.**
**Subject:** `RECONCILE-howl-tts-stage3-c1-p17-review-2026-11-09.md` (`b634cd20…`)
**Disposition:** **NOT CONFIRMED AS COMPLETE.**
**C1 IMPLEMENTATION REMAINS UNAUTHORIZED.**

The body below is transcribed without modification from the second-stage review delivered on 2026-11-09. Only this container header was added.

<!-- BEGIN VERBATIM REVIEW -->

I agree with the NOT CONFIRMED AS COMPLETE disposition, based on the review you supplied.

## The three corrections that matter most

1\. R01 — restore the correct source mapping. The original P17 references at `client.go:439/441` were `callOn` timeout/cancellation branches. The reconciliation incorrectly reassigned them to `Synthesize`. The underlying finding remains valid, but the correction must be withdrawn.

2\. R03 — correct the transport-send citation. `:410–412` identifies encode failure; `:420–421` identifies the transport-send failure. The finding survives, but the source trace must be exact.

3\. R08 — stop treating 45 as a reproducible site count. The arithmetic is correct, but the categories overlap and some rows lack branch-level mappings. The next inventory needs unique site IDs, exact lines, branch conditions, and emission types. A site appearing in two conceptual categories must not silently become two distinct sites.

## What I would preserve

The review also correctly distinguishes several levels of acceptance:

- R02, R04–R06: source behavior confirmed.
- R07: conceptual finding confirmed, inventory still needs reproducible mapping.
- R09–R13: correctly left as unresolved model and compatibility questions.
- E1–E9: retained individually, with E7–E9 confirmed as additional loss paths.
- U1–U8: still open; no precedence rule has established causality.

The helper digest result is also useful. The review has identified the exact manifest-style hashing procedure that reproduces `d0df8b28…`. That resolves the digest-value discrepancy, but the procedure and inclusion of test files need to be documented in the reconciliation itself.

## Recommended gate

I would make the next step a targeted author response to this review, not another P17 proposal revision.

It should address:

- RC-01 and RC-02: correct the source mappings.
- RC-03 and RC-04: produce a unique, branch-level site manifest and recompute the count from that manifest.
- RC-05: document the exact digest command, ordering, path scope, and whether test files are included.
- RC-06: correct the `callOn` range and branch references.
- Preserve every R01–R13, E1–E9, and U1–U8 identifier and its existing disposition unless new evidence justifies a change.

One additional safeguard: the reviewer should be able to regenerate the site count and helper digest from the documented procedure without relying on the author's narrative.

Status remains unchanged: P17 not approved; reconciliation not independently confirmed as complete; no contract amendment, rev4, P18, or C1 implementation authorized.

The important progress here is not a higher finding-closure count. It is that the review process is now catching errors in the reconciliation itself instead of allowing an inaccurate correction to become the next frozen baseline.&#x20;

<!-- END VERBATIM REVIEW -->
