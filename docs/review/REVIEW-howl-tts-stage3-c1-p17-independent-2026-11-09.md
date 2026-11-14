# P17 — Independent Review of the Evidence-Semantics Proposal

**Date:** 2026-11-09
**Status:** INDEPENDENT REVIEW OF RECORD. **Preserved verbatim, unedited.**
**Subject:** `PROPOSAL-howl-tts-stage3-c1-evidence-semantics-p17-2026-11-09.md` (`01d887c1…`)
**Disposition:** **NOT APPROVED.**
**C1 IMPLEMENTATION REMAINS UNAUTHORIZED.**

The body below is transcribed without modification from the independent review delivered on 2026-11-09. Only this container header was added. Hash is recorded in the accompanying reconciliation and validation output.

<!-- BEGIN VERBATIM REVIEW -->

I agree with the NOT APPROVED verdict based on the findings supplied. I would not advance P17 into the contract or treat its proposed semantics as accepted.

## What this review establishes

There are three distinct layers, and they should remain separate.

### 1. Source-confirmed defects

R01–R08 identify concrete problems in the reported source:

- Operation context is missing. The same timeout/cancellation site serves ordinary requests, handshake, and shutdown. Severity cannot be assigned from the source line alone.
- Fatality is not consistently preserved. A fatal failure can travel through `failAll` and emerge at a caller that publishes `degraded`.
- Fatal observations do not consistently trigger lifecycle action. Recording `helper_crashed` does not guarantee request invalidation, termination, or other enforcement.
- Reader completion is unobservable. `waitForReader` cannot tell its caller whether the reader actually joined.
- The evidence-loss inventory is incomplete. E7–E9 add late-terminal dropping, ignored cleanup/cancellation outcomes, and shutdown state clearing after failed termination.
- Termination causes are conflated. Context cancellation and grace expiry are different outcomes.
- Observation rows combine incompatible phases and meanings. In particular, handshake, load, and generic request failures need separation.
- The inventory's completeness claim is overstated. Row counts are not distinct source-site counts, and the Git baseline does not identify the untracked helper source.

These are not merely wording problems. They affect what information a future implementation could actually preserve.

### 2. Proposed-mechanism deficiencies

R09–R13 are different: they challenge whether P17's replacement model is sufficiently specified.

The most consequential is R09. An append-only log is not automatically a reliable evidence model. Without instance identity, sequence-allocation rules, causal linkage, action results, phase context, and retention rules, it can preserve entries while still producing an ambiguous or misleading history.

R10 and R11 expose two more important gaps:

- Precedence is currently a selection rule, not a demonstrated resolution of competing evidence.
- `quarantined` is one of the five proposed projections, but P17 does not define a path that can produce it.

R12 and R13 prevent a premature claim of compatibility with the existing application. Handshake classification and provider-state behavior still have unresolved semantics.

### 3. Unresolved cases

U1–U8 should remain individually open. I would resist any attempt to close them by simply assigning a higher precedence to one event.

For example, a protocol violation observed before a zero exit does not, by itself, prove that the violation caused the exit. Likewise, a nonzero exit after a deliberate kill does not independently establish a crash.

The review correctly distinguishes what was observed from what caused the outcome.

## Recommended next step

I would return this to DeepSeek for a finding-by-finding reconciliation, not another proposal revision yet.

The reconciliation should:

1. Preserve R01–R13 as separate review findings.
2. For R01–R08, identify the exact source behavior, missing observation or distinction, and evidence required to establish the correction.
3. For R09–R13, identify the missing model definition or compatibility decision without presenting a proposed mechanism as an established fact.
4. Keep E1–E9 individually represented, including the original six.
5. Reconcile the observation and termination inventories at the source-site and branch level, not merely by row count.
6. Keep U1–U8 open unless a specific, defensible resolution is demonstrated.
7. Preserve the review's provenance limitation: the helper directory is untracked, so the Git HEAD pin alone does not establish its source identity.

One particularly important distinction: a proposed evidence log must not become a substitute for fixing evidence-loss paths. If a source path drops an event before it reaches the log, the log's append-only property provides no protection.

## Current disposition

| Area                    | Disposition                                            |
| ----------------------- | ------------------------------------------------------ |
| P17 inventory           | Partial; completeness not established                  |
| Source findings R01–R08 | Require individual author reconciliation               |
| Model findings R09–R13  | Require specification and compatibility reconciliation |
| U1–U8                   | Open                                                   |
| P17 approval            | Not approved                                           |
| Contract amendment      | Not authorized                                         |
| C1 implementation       | Not authorized                                         |

The review does not invalidate the value of P17. It establishes that P17 is a useful first inventory, but not yet a dependable specification.

I recommend preserving this review unchanged and having DeepSeek reconcile it before anyone writes P18 or a rev4 contract.

<!-- END VERBATIM REVIEW -->
