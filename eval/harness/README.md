# MEMORY-09 harness — P1 build

**Authority:** `BEA-1` · `P1 Build` authorisation (2026-11-13). **P2 execution is not granted.**
**Status:** instrument only. Nothing here has produced behavioural evidence, and nothing can until a separate owner decision.

## What P1 built

| File | Role |
|---|---|
| `adapter.py` | `GenerateRequest`/`GenerateResponse`, T-1/T-2/D-1 enforcement, provenance field list, `OfflineAdapter` (refuses), `HuggingFaceAdapter` (grant-gated, never constructed) |
| `state.py` | Entity model, hash-on-mutation, immutable snapshots, diffs, recorded removal, no silent normalisation |
| `fixtures.py` | `F-ST-01-S1` — one invented canonical fact, one probe, declared guards |
| `conditions.py` | Assisted/unassisted construction, frozen params, structural separation enforcement |
| `provenance.py` | Record assembly with `unavailable` markers, admissibility checks, record writer |
| `plan.py` | 18-call plan, 30-call ceiling, unassigned-slot barrier |
| `authorization.py` | P2 grant gate (`eval/runs/P2-GRANT.json` — must not exist during P1) |
| `budget.py` | Route policy: $0.00 paid ceiling, free-credit attestation, route-change refusal |
| `local_adapter.py` | Loopback-only KoboldCpp / SillyTavern / llama.cpp adapter, same request contract |
| `cost.py` | Pricing pre-flight: no quote, no spend; worst-case bounded against the $0.00 ceiling |
| `dispatch.py` | The dispatch boundary: `CallLedger.reserve` counts before the call, failures consume the budget |
| `local_surfaces.py` | KoboldCpp native + OpenAI-compatible serialisers, surface detection, transcribers, error handling |
| `local_fixtures.py` | Nine synthetic local reply bodies + one unrecognised body |
| `discover.py` | L2 identity GETs: loopback only, three fixed paths, no prompt, no server start |
| `l3.py` | L3 single-call local run: one dispatch, runtime provenance, requirement-by-requirement verification |
| `validate.py` | 88 offline checks under a socket guard, including scratch-lifecycle, hosted-refusal, local-surface, discovery and L3 checks |
| `cli.py` | `inspect`, `plan`, `validate`, `l3`, `l3-verify`, `discover`, `run` (refuses P2) |

## Commands

```
python -m eval.harness.validate           # offline checks, socket guard active
python -m eval.harness.cli inspect        # fixture, conditions, hashes
python -m eval.harness.cli plan           # caps and slot state
python -m eval.harness.cli run            # refuses; exit 2
python -m eval.harness.cli routes         # permission table and spend ceiling
python -m eval.harness.cli discover       # L2 identity GETs (add --record to file the block)
python -m eval.harness.cli l3 --record    # ONE local generation; needs eval/runs/L3-GRANT.json
python -m eval.harness.cli l3-verify --dir eval/runs/L3-<stamp>   # re-check a filed record, no calls
```

`l3` verifies the model sha256 against the grant before sending anything, refuses if no model is
resident, and dispatches exactly once through `CallLedger` — a failed call consumes the allowance
and there is no retry path. `l3-verify` re-reads filed evidence and makes no calls, so a corrected
checker can be applied to a call that has already been spent.

## Why execution cannot happen by accident

Two independent barriers, either of which alone is sufficient:

1. **Grant gate.** `HuggingFaceAdapter` raises unless handed an `ExecutionGrant`; `ExecutionGrant` can only come from `eval/runs/P2-GRANT.json` with owner signature, named models, pinned providers and caps within 30; that file does not exist. The CLI's `run` command is a refusal path.
2. **Plan gate.** `build_plan` refuses while any slot is `UNASSIGNED`, and P1 ships every slot unassigned.

A third, structural: no code path calls a provider. `HuggingFaceAdapter.generate` ends in `NotImplementedError` because response transcription is P2 work.

A fourth, financial: `build_adapter` runs the budget policy **before constructing anything**. Hosted routes are refused without a free-credit attestation, and no attestation file exists, so even a forged grant could not reach a paid route.

A fifth, structural: `Dispatcher` is the only component that touches `adapter.generate`, and `CallLedger.reserve` increments **before** the call, so the 30-call ceiling binds at dispatch rather than at planning. Failed calls consume slots.

## Spend policy

**No payment method is on file, and none will be added.** So paid inference is not *prohibited* — it is *unreachable* (`PAYMENT_METHOD_ON_FILE = False`). The rule:

1. Hosted is **refused**. A free account has no included monthly allowance (HF docs change, 2026-10-07), so zero out-of-pocket spend cannot be guaranteed before dispatch. `budget.HOSTED_REFUSAL` is raised by the policy, the factory, and the adapter constructor (checks 27, 35, 59–62). The `HuggingFaceAdapter` remains in the tree as **unverified, unreachable code** — its existence is not authorisation to use it.
2. A payment demand (HTTP 402, "insufficient credits", "quota exceeded", "add a payment method") is a **stop**: no retry, no top-up, no automatic switch (`budget.PaymentDemanded`, raised at the dispatch boundary).
3. The fallback is local (`local:koboldcpp`, `local:sillytavern`, `local:llamacpp`, `local:ollama`), chosen by a human, requiring an explicit loopback endpoint, still gated by the P2 grant.

Verification is deliberately lightweight: the provider's own response is the signal. No balance API, no spend authorisation flow, no credential that could incur a charge. `cost.py` remains as belt-and-braces (no quote, no spend) but nothing depends on it once there is nothing to charge.

## Local surfaces (L1)

`local_surfaces.py` builds both request bodies from one `GenerateRequest` and transcribes both reply shapes into one record. Native requests are flattened to a prompt (template `system\n\nuser`), so a cross-surface comparison is a **different condition**, not a re-run. Parameter renaming (`max_tokens → max_length`, …) is recorded; a parameter with no native equivalent is refused rather than dropped. Truncation (`finish_reason == "length"`) is recorded as `incomplete`, never as a failure. Error bodies are handled **before** shape detection, because a local 500 has neither `choices` nor `results`.

**No server has been contacted.** Every fixture is documentation-derived; L3 corrects fixtures against a real body.

## L2 discovery

`discover.py` performs three GETs and nothing else: `/api/extra/version`, `/api/v1/model`, `/api/extra/true_max_context_length`. Loopback only; a non-loopback endpoint is refused. No prompt, no model load, no server start, no retries. Absent fields are recorded `unavailable` rather than guessed. If nothing answers, that is the recorded result — the harness does not start servers. Last run **established**: KoboldCpp **1.120**, model **inactive**, configured context **12288**, loopback only (2026-11-14). Identity blocks are filed under `eval/runs/`, one directory per attempt, write-once. Note the real `/api/extra/version` shape: name in `result`, version in `version` — the parser prefers `version` and keeps the name separately (check 82).

## Not done / not enforced

- `HuggingFaceAdapter` is **written but never executed**: not import-tested against a real client, not verified against provider response shapes. Treat it as a draft for P2, not tested code.
- `eval/runs/` does not exist; creating it does not grant anything.
- The socket guard blocks socket creation during validation only. It is a validation instrument, not a runtime sandbox.
- No `.gitignore`, packaging or CI wiring was added — out of P1 scope.

## Scratch artifacts

Test scratch lives in `eval/harness/_selftest-*` and `_grant-*` only. It is created by `scratch_dir()`, which owns the exact path and removes it in a `finally` — success, failure and early return all clean up. Cleanup never globs, so it cannot delete a directory it did not create, and it never touches `eval/runs/` (persistent evidence). Checks 25, 53–58 fail if either family is left behind.

## Reading a result later

`record.json` per call, with `raw_body.json` beside it. A record whose `raw_body` or `provider_used` is `unavailable` is inadmissible (`provenance.check_admissible`).
