# P17 — Third-Stage Review of the Site Manifest

**Date:** 2026-11-09
**Status:** REVIEW OF RECORD. **Preserved verbatim, unedited.**
**Subject:** `RESPONSE-howl-tts-stage3-c1-p17-review-r2-2026-11-09.md` (`57cc58f5…`)
**Disposition:** RC-03 and RC-04 partially resolved; CR-01–CR-08 require response.
**C1 IMPLEMENTATION REMAINS UNAUTHORIZED.**

The body below is transcribed without modification from the third-stage review delivered on 2026-11-09. Only this container header was added.

<!-- BEGIN VERBATIM REVIEW -->

I agree with the findings and the partial dispositions.

## The most important corrections

1\. O12/O13/O14: The corrected source mappings have not been propagated into all manifest associations. The `Synthesize` send, timeout, and cancellation paths need explicit associations rather than relying on their appearance under O18/O19.

2\. E8/E9: The manifest must represent all five ignored `Terminate` calls and the ignored `Cancel` call. E9 also needs the separate `started = false` operation at `client.go:671`.

3\. E7: The actual silent return is at `client.go:342`, not the preceding unlock at `:341`.

4\. O2/O3/O6/O17: State and evidence operations within those observation paths are missing. This matters because a manifest should represent the relevant behavior, not merely the most convenient line in each branch.

5\. The inclusion rule: This is the underlying accounting problem. Either include all state publications covered by the current wording or narrow the rule to explicitly selected evidence-bearing publications. Don't change the rule silently to preserve 67.

## One additional distinction

The reviewer is right to flag `S20` and `S25`: function declarations are mechanism anchors, not emission sites.

I'd maintain three separate counts going forward:

- Emission anchors: actual source lines where an event, state change, or action occurs.
- Mechanism anchors: declarations or shared implementations that explain how those operations work.
- Branch identities: distinct executable conditions and outcomes, independent of how many source lines each spans.

That avoids another situation where a technically reproducible count is mistaken for proof of coverage.

## Recommended disposition

| Area                       | Status                                          |
| -------------------------- | ----------------------------------------------- |
| RC-01, RC-02, RC-05, RC-06 | Resolved                                        |
| RC-03, RC-04               | Partially resolved                              |
| CR-01–CR-08                | Require targeted author response                |
| R01–R13                    | Preserve existing dispositions                  |
| E1–E9                      | Preserve individually; correct missing mappings |
| U1–U8                      | Open                                            |
| P17                        | Not approved                                    |
| C1 implementation          | Unauthorized                                    |

The next step should be a targeted manifest correction, followed by independent coverage verification. No need to rewrite P17 or reopen settled source facts.

The governing principle remains: reproducibility proves that someone else can regenerate the inventory; completeness requires demonstrating that the inventory includes everything covered by its declared scope.

<!-- END VERBATIM REVIEW -->
