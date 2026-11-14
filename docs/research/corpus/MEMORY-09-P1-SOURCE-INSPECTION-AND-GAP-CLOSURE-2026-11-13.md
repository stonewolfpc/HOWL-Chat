# `MEMORY-09` — P1 Source Inspection and Gap Closure

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-P1SI-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **INSPECTION COMPLETE — 8 gaps found, 8 closed offline. No network, no model, no spend.**
**Addendum A1 (2026-11-13):** three follow-up points filed at owner direction — documentation only, no code, test, grant, or adapter change.
**P1.1 (WO-ARCH-07, 2026-11-13):** offline hardening applied — scratch lifecycle fixed, check 25 corrected (58 checks), free-allowance preflight found **not feasible** (§9). No network, no inference call, no P2 execution.
**P1.2 (2026-11-13):** hosted refusal **enforced** in code at three layers (62 checks) — the policy-versus-code gap recorded in §9.5 is closed (§10).
**Authority:** owner inspection order (2026-11-13). P2 remains unauthorised.
**Manifest:** this document is authoritative for current P1 file hashes; the table in row 483 is the build-time snapshot.
**Revision 1.1 (2026-11-13):** after the owner confirmed no payment method exists, the spend rule was simplified from *paid prohibited with attestation* to *paid unreachable*. Checks 27-30 were replaced by the lightweight rule and check 52 added; manifest below re-stated at 52 checks.

---

## 1. Method

Read every file, then trace the path: **grant → plan → dispatch → adapter → record**. For each of the four scrutiny areas, ask what the code *actually does* rather than what the tests assert, then write a check for each protection that had none.

**What this inspection is not:** no static analyser, no type checker, no formal verification. It is a read-through plus targeted checks, and it finds what reading finds.

## 2. Findings — eight gaps, all closed

| # | Gap | Severity | Closure |
|---|---|---|---|
| 1 | **Ceiling enforced at planning, not dispatch.** `check_caps` was only called by tests; a future path could call an adapter directly | **High** — the owner's exact point | New `dispatch.py`: `CallLedger.reserve` increments before the call, `Dispatcher` is the only caller of `generate`, `run` goes through it (checks 39, 40, 51) |
| 2 | **Evidence could be overwritten.** `write_record` used `exist_ok=True` and replaced `record.json` | **High** — evidence that can be replaced is not evidence | Write-once: second write raises; `record.sha256` sidecar; `verify_record` detects edits (check 45) |
| 3 | **Evidence could be written outside `eval/`** | Medium | Resolved path must sit under `eval/` (check 46) |
| 4 | **Snapshot was mutable.** A frozen dataclass holding a dict does not freeze the dict — a caller could mutate data after it was hashed | **High** — silently breaks `C2` evidence | `Snapshot` stores canonical JSON, parses on read (check 47) |
| 5 | **Grant validation was semantic, not typed.** `models` could be a list, `call_ceiling` a string, `owner_signature` any non-empty string; caps could be edited after signing | **High** | Types, date shape, ceiling range, signature length, and `caps_sha256` binding the signed fields (checks 41–44) |
| 6 | **No pricing pre-flight.** Row 484 requirement 3 was policy without code | **High** for P2 | New `cost.py`: no quote → refuse; paid quote → refuse; verified zero-cost quote → priced at 0.00 (checks 36–38) |
| 7 | **Adapter could send params the client rejects.** `GenerateRequest` validated names, but `chat_completion` may not accept `top_k`/`min_p`; that would be a `TypeError` after authorisation | Medium | Offline signature inspection of the installed client; unsupported params fail **before the wire** (check 48) |
| 8 | **Transcription unfinished** (`NotImplementedError`) | Medium | `transcribe()` implemented and tested against a **synthetic** body — no call needed; `raw_body` preserved verbatim (check 49) |

**Also closed:** `run` previously refused at the grant and never touched the dispatcher; it now constructs the ledger and dispatcher on the reachable path, and a check asserts the CLI contains no direct `generate(` call (check 51).

## 3. Map of the trace, after closure

| Stage | Gate | Enforced by |
|---|---|---|
| Grant | signed, typed, caps-bound, $0.00, slots named, providers pinned | `authorization.load_grant` |
| Plan | 18 calls, ceiling 30, slots assigned | `plan.build_plan` |
| Budget | free-credit attested or refused | `budget.check_provider` |
| Price | worst case ≤ $0.00, quote required | `cost.preflight` |
| **Dispatch** | **count before call; failures consume; 31st refused** | **`dispatch.CallLedger.reserve`** |
| Adapter | provider pinned, params acceptable | `adapter.GenerateRequest`, `check_params_supported` |
| Record | write-once, hash-verified, inside `eval/` | `provenance.write_record` |

## 4. Remaining prerequisites for a responsible P2 decision

| Item | Why it still blocks P2 |
|---|---|
| Real provider response shape | `transcribe` is tested only against a hand-built body; the live shape is unverified |
| Model and provider naming | Slots are `UNASSIGNED`; the plan cannot produce calls |
| Free-credit attestation | Absent, so every hosted route is refused by policy |
| Cost quote | `cost.py` requires a quote; none exists, and no pricing has been retrieved |
| Local endpoint availability | No KoboldCpp/SillyTavern endpoint has been probed |

**Two of these are code-complete and waiting on owner input** (names, attestation). Two are inherently unverifiable offline (live response shape, local endpoint). One — pricing retrieval — could be checked without a generation call, and has not been.

## 5. Validation

`python -m eval.harness.validate` → **52/52 passed, exit 0**, socket guard active across all checks. `cli run` → exit 2, refused.

Checks 1–26 and 31–35 were preserved through gap closure; 27–30 were replaced when the payment rule was simplified; 36–52 are new.

## 6. Manifest (16 files, 97,833 bytes)

| File | Bytes | SHA-256 |
|---|---:|---|
| `eval/__init__.py` | 168 | `d99d08b1213f74f44aaa5656737cfc2c6801aa6431cc6486f308096785370c17` |
| `eval/harness/__init__.py` | 292 | `3dced7dee1277d415c3b3daf645bf427462b2f475cea12877c859f600ad8396f` |
| `eval/harness/adapter.py` | 11,836 | `d9ff3c63867b060243f0199c8fa511edd21d8daaa9ce71e16149f1e53000aa02` |
| `eval/harness/authorization.py` | 5,013 | `e165856d3636e51df44e026403a2eb2104c7e2d14b43849df16fa5fbf454a473` |
| `eval/harness/budget.py` | 5,207 | `03ef962d2f9e02646cd90368b2338039ce7305c3e170cffd243bc1d73a406154` |
| `eval/harness/cli.py` | 4,744 | `513410c590a4d0da1a9ce5c1ecb00141ca4574e7855211e52d782beb7fd1a0ce` |
| `eval/harness/conditions.py` | 4,365 | `ae3571e1de5b380012c0b1f7cca99bf48c6f30ac8ffd62c220297c33c0b62fff` |
| `eval/harness/cost.py` | 2,658 | `fc84bb54b80da9d7fd78e14ec851fe7c9a10bc8f0858064cff01b23ff574b2f1` |
| `eval/harness/dispatch.py` | 2,975 | `afb85faaf56669285ca752f85ac6005cd23405421434f0ef245c4c1f7379b7c2` |
| `eval/harness/fixtures.py` | 5,424 | `2a7ad5a53d03b4ac2effefdcfdb4599f92493e8c8a127f9f1e0dcb4347db54b3` |
| `eval/harness/local_adapter.py` | 3,877 | `07657261d0c266e81e950af15f1285881c24a426151a84ba1672ba06f2495ac6` |
| `eval/harness/plan.py` | 3,917 | `640988f9003dea2db83588d28bd9731a078a8ccdba2dbd1ab6c93e409d5c1396` |
| `eval/harness/provenance.py` | 6,651 | `32da8cc65c9c2c021b2a51bde88da6466d6aeb80526ec702df95327e55ed5e1b` |
| `eval/harness/state.py` | 5,492 | `37ff62d6d3204f99102728d9773d98f6fb26e5287993924fc2fde3d25658e12e` |
| `eval/harness/validate.py` | 30,445 | `b1b6c8584ec6cc056cc0297c37ebdc9260cc1919f5b9011975af21782ed919ae` |
| `eval/harness/README.md` | 4,769 | `6d0b9a274762209db917c5e03317f054280137b896f6e7c386c111e70b5d31bf` |

## 7. What still cannot be claimed

That the harness is *safe*; that the adapter matches the live provider contract; that any call would succeed; that the safety properties hold against a deliberately hostile caller (this is not a sandbox — a determined process could import `socket` itself).

What can be claimed: the eight gaps above were found by reading, closed in code, and each has a check that fails if the protection is removed.

---

## 8. Addendum A1 — follow-up filing (documentation only)

**Filed:** 2026-11-13 · **Scope:** record only · **No implementation, no execution, no inference request.**

### 8.0 Decision state preserved

| Item | Decision |
|---|---|
| P1 offline inspection | **Accepted** |
| P2 hosted execution | **Not authorized** |
| Paid inference | **$0 out-of-pocket permitted** |
| Live adapter | **Unverified** |
| Follow-up work | **Documentation only** |

### 8.1 Hugging Face free-credit semantics — **unresolved verification requirement**

Three quantities are distinct and must not be conflated:

| # | Quantity | Meaning |
|---|---|---|
| (a) | **Published inference price** | What the provider's price list states for the request |
| (b) | **Free-credit coverage** | Whether *this* request is fully absorbed by the account's available monthly free allowance |
| (c) | **Out-of-pocket charge** | Whether money actually leaves the account |

The governing policy is **(c) = $0**, not (a) = $0.00. A nonzero published price (a) is acceptable **only** where (b) is reliably established and automatic paid usage cannot occur.

**Verified facts about the current harness:**

- `budget.check_provider` gates hosted routes on `PAYMENT_METHOD_ON_FILE` alone. It does **not** read, estimate or verify a remaining allowance — there is no code path that establishes (b).
- `cost.preflight` exists and refuses unpriced or over-ceiling plans, but **is called only by the validation suite**. No execution path invokes it (`cli run` → dispatcher → adapter; no preflight call).
- The only allowance signal the harness can currently observe is **response-time**: HTTP 402, "insufficient credits", "quota exceeded", "add a payment method" → `PaymentDemanded`, raised at dispatch after the request has already been sent.

**Consequence, stated plainly:** the requirement "if coverage, remaining allowance, or paid-usage behaviour cannot be reliably established, execution must be refused" is **not satisfied before the call**. The harness refuses *after* a provider signals a demand, and relies on the owner statement that no payment method exists for the claim that no charge is possible.

**Unresolved:** how coverage is established (account allowance API, provider-side spending cap, prepaid-only account, or an owner-maintained statement refreshed per series). No solution is implemented, and none is authorised by this filing.

### 8.2 Owner signature authentication — **not established**

`caps_sha256` binds the signed fields, so editing `call_ceiling`, `cost_ceiling_usd`, `date`, `models` or `providers` after the hash was written breaks validation. That is **integrity against accidental or casual edit**, and nothing more.

What the implementation actually checks on `owner_signature`:

```
if not isinstance(signature, str) or len(signature.strip()) < 8:
    raise ExecutionNotAuthorised("P2 grant carries no usable owner signature")
```

**Exactly what that establishes:** the field is a string of at least eight characters after trimming. It is a **presence-and-shape check**.

**What it does not establish, and must never be described as establishing:**

- No cryptographic signature verification.
- No trusted owner public key, no key registry, no trust anchor of any kind.
- No binding between the signer's identity and the file.
- No protection against forgery: `caps_sha256` is computed by `authorization.caps_sha256`, a public function with **no secret input**. Anyone who can write the grant file can compute a matching hash.

**Therefore:** a grant passing validation proves *internal consistency and well-formedness*, not that the owner wrote it. The harness has **no authenticated owner authorisation**; the authority for execution is the owner's message in the project record, not the file.

**Unresolved:** whether any signing mechanism is wanted. None is implemented, and none is authorised by this filing.

### 8.3 Live adapter verification — **unverified**

Check 49 exercises `transcribe()` against a **hand-built body** shaped by us, not by the provider. That establishes the transcription logic is self-consistent; it establishes **nothing** about the real Hugging Face response structure.

Still unknown, and unverifiable offline: the live envelope (where `provider`, `model`, `system_fingerprint`, usage and timestamps actually appear), whether streamed and non-streamed shapes differ, how provider-pinned routing reports the serving backend, and which errors arrive as status codes versus bodies.

**A real provider call would be required to establish live response compatibility. It is not authorised by this filing, and none was made.**

### 8.4 Incidental finding during this filing — test hygiene

The source search performed for §8.2 exposed **42 leftover `_grant-*` directories** under `eval/harness/`, each holding a synthetic grant written by the validation suite. They were **removed** as housekeeping (a deletion of generated temp files, not a code change).

The cause is a defect worth recording: `_write_grant` creates a `tempfile.mkdtemp` directory and never removes it, and **check 25 ("cleanup: no self-test artifacts left behind") only globs `_selftest-*`**, so its name overstates what it verifies. The check passes while grant scratch dirs accumulate.

**Not fixed:** modifying the test is outside a documentation-only filing. Recorded as a known defect with a named remedy (glob `_grant-*` too, and clean up in a `finally`).

### 8.5 Claim ceiling for this addendum

This addendum records facts about the code as inspected and questions that remain open. It changes no behaviour: the harness still runs 52/52 offline checks, `run` still refuses, and no inference request was made while filing it.

---

## 9. Addendum A2 — P1.1 offline hardening (WO-ARCH-07)

**Authority:** `WO-ARCH-07`, 2026-11-13 · **Mode:** local only · **No network request, no inference call, no P2 execution.**
**Scope:** scratch-artifact lifecycle, check-25 correctness, regression tests, and a free-allowance preflight feasibility study.

### 9.1 Workstream A — grant scratch lifecycle

**Defect and exact cause.** `_write_grant` called `tempfile.mkdtemp(prefix="_grant-", dir=<harness dir>)` and returned the file path. The directory handle was discarded, so nothing could remove it — not on success, not on failure, not on exception. Each validation run therefore created six permanent directories (checks 41–44, the `_grant_payload` loop in 42 counting three). Confirmed empirically before the fix: **42 directories present**, and one run regenerated **6**.

**Before → after.**

| Path | Before | After |
|---|---|---|
| Successful check | directory leaked | removed in `finally` |
| Refused grant (`ExecutionNotAuthorised`) | leaked | removed |
| Exception inside the block | leaked | removed |
| Repeated runs | accumulated | zero accumulation (3 consecutive runs verified) |
| Scope of deletion | n/a | only the path the context holds |

**Implementation.** `scratch_dir(prefix)` creates via `mkdtemp` and yields the path; removal happens in `finally`, so success, failure and early return all clean up. `scratch_grant(payload)` builds on it and writes the synthetic grant inside. **Ownership is explicit:** the context holds the exact path it created and never globs to decide what to delete, so a directory the operation did not create cannot be removed. Cleanup is never broadened into a filesystem sweep.

**Persistent vs disposable, stated:** `eval/harness/_selftest-*` and `_grant-*` are disposable scratch and are deleted by their owner. `eval/runs/` is **persistent evidence** territory and is never touched by cleanup; `eval/harness/*.py` and the README are source, not scratch.

**Grant-validation behaviour:** unchanged. Checks 41–44 assert the same refusals as before; only the directory lifecycle moved.

### 9.2 Workstream B — check 25 corrected

**Original defect:** `_c25` globbed only `_selftest-*`, so it could pass while `_grant-*` scratch accumulated — the check's name claimed more than its scope delivered. It also derived its target through `parents[2] / "eval" / "harness"`, which resolves correctly but obscures the intent.

**Corrected:**

- `SCRATCH_PREFIXES = ("_selftest-", "_grant-")` and `find_scratch(root)` scan **both** families in **`HARNESS_DIR = Path(__file__).parent`**.
- `find_scratch` raises if the scan root is not a directory, so the check cannot pass by silently scanning nothing.
- `_c25` asserts its own root is the harness directory before scanning.
- Name now states what it verifies: *"cleanup: no test scratch directories remain (self-test or grant)"*.
- A **final** check (`_c58`) repeats the scan at the end of the run, because the grant checks execute *after* check 25 — so check 25 alone cannot certify a whole run.

### 9.3 Workstream C — regression tests and results

| # | Test | Outcome |
|---|---|---|
| 53 | a leaked `_grant-` directory is detected | **PASS** |
| 54 | a leaked `_selftest-` directory is detected | **PASS** |
| 55 | owned directories removed on success **and** on failure | **PASS** |
| 56 | cleanup does not touch unrelated files or directories | **PASS** |
| 57 | a failing grant validation leaves no new directories | **PASS** |
| 58 | a complete run leaves no scratch directories behind | **PASS** |

**End-to-end reproduction of the original failure.** With a `_grant-leak-probe` directory injected: **55/58, 3 failures, exit 1** — check 25, check 53 and check 58 each caught it. On removal: **58/58, exit 0**.

**Repeated runs:** three consecutive runs, exit 0 each, **0 accumulated scratch directories**.

**Full suite:** 58 checks, 58 passed, 0 failed, exit 0, socket guard active across all checks.

### 9.4 Workstream D — free-allowance preflight feasibility

Documentation review only; no inference endpoint was contacted.

| Question | Answer | Basis |
|---|---|---|
| Applicable published price per request | **Partly available** | Provider pricing is documented as pay-as-you-go at provider rates with no markup; the *request-specific* price is not returned before dispatch |
| Remaining free monthly allowance | **Unavailable for a free account — there is none** | HF docs: free users get **no** included monthly credits; only PRO ($2/mo compute credits) and Team/Enterprise ($2/seat) have included credits. The former $0.10 free allowance was removed (docs change merged 2026-10-07) |
| Whether a specific request is fully covered | **Cannot be established pre-dispatch** | Cost depends on token counts known only after the call |
| Whether out-of-pocket charge could occur | **Cannot be established** | Follows from the two rows above |
| Whether paid usage can be reliably disabled | **Only for Team/Enterprise organizations** | Documented: org admins can set a spending limit and disable providers. No documented equivalent for personal accounts |

**On programmatic balance checking:** the documented billing endpoint (`/api/organizations/{org}/billing/usage/inference`) is **organization-scoped, retrospective** (results appear up to **2 hours** after requests) and requires a token with billing-read permission. It reports usage already incurred; it does not announce remaining allowance before a call. `/api/settings/billing/usage-v2` appears in the OpenAPI spec but its response shape is undocumented, so it is not a reliable basis for a spend decision. `whoami-v2` reports plan (free/PRO) but no credit balance.

**Conclusion — the acceptance rule is triggered.** Free coverage and zero out-of-pocket behaviour **cannot** be established reliably before dispatch for this account. Per `WO-ARCH-07` §5, the harness must therefore **refuse hosted execution in any future authorised run**. Reported, not implemented: changing dispatch behaviour is explicitly outside this order.

**Sharper consequence for this project:** since the owner has no payment method and free accounts now carry **zero** included credits, the hosted Hugging Face route is not merely unverifiable — it is **unavailable** without purchased credits or a paid subscription. That is an owner decision about money, not a harness question.

**Sources consulted (documentation only, 2026-11-13):** `huggingface.co/docs/inference-providers/pricing` (included credits table, pay-as-you-go, org spending limit and provider disabling) · `huggingface.co/docs/hub/enterprise` (plan comparison, included credits per seat) · `github.com/huggingface/hub-docs` PR #2865, merged 2026-10-07 (free accounts no longer receive included monthly credits) · `huggingface.co/docs/hub/billing` (usage tracking, billing dashboard).

**Recommendation (no implementation):** treat S1 as **local-first**. If hosted is ever revisited, the only documented configuration that bounds spend is a Team/Enterprise organization with a spending limit set and providers disabled — which is itself a paid subscription, so it contradicts the $0 rule and would need explicit owner authorisation.

### 9.5 Deviations, limits, and unresolved defects

| Item | State |
|---|---|
| Scratch scan scope | Scans the harness directory only, non-recursively. It does **not** assert that a leaked directory was never committed to git |
| Free-allowance refusal | **Not implemented** — recommended, per §9.4; requires authorisation |
| Live adapter | Still unverified (§8.3 unchanged) |
| Owner-signature authentication | Still absent (§8.2 unchanged); `caps_sha256` remains integrity-only |
| Hosted route availability | Hosted remains *policy-permitted* in code while §9.4 concludes it should be refused — an unresolved **policy-versus-code** gap, recorded rather than silently reconciled |

**No deviation from the order** was required: no change touched architecture, corpus, adapter behaviour, grants, or the P2 boundary.

### 9.6 Claim ceiling

P1.1 establishes that scratch artifacts no longer accumulate, that the cleanliness check detects both families including an injected leak, and that hosted free-allowance preflight is not feasible for a free account. It establishes nothing about model behaviour, and it does **not** make hosted execution safe.

---

## 10. Addendum A3 — P1.2 hosted refusal enforced

**Authority:** owner direction, 2026-11-13, following §9.4: *"make the current zero-out-of-pocket policy enforceable at dispatch"*, keeping local inference as the intended S1 path.
**Change class:** bounded policy enforcement. No architecture change, no new route, no inference call.

### 10.1 The gap this closes

§9.5 recorded a **policy-versus-code gap**: §9.4 concluded hosted must be refused, while `budget.check_provider` still *permitted* it because no payment method exists. A finding that is not enforced is a note, not a control. It is now enforced.

### 10.2 Three layers, each with a check

| Layer | Behaviour | Check |
|---|---|---|
| **Policy** — `budget.check_provider` | Raises `BudgetRefusal(HOSTED_REFUSAL)` for every non-local provider. The refusal states the reason (no included monthly allowance; price and coverage unverifiable pre-dispatch) and names the local fallback | 27 |
| **Factory** — `build_adapter` | Runs the policy first, so no hosted provider string produces an adapter, with or without a grant. Routing suffixes are rejected earlier still (T-2) | 35 |
| **Object** — `HuggingFaceAdapter.__init__` | Refuses **even with a valid grant** when `FREE_ALLOWANCE_EXISTS` is false | 59 |

`FREE_ALLOWANCE_EXISTS = False` is the single switch, and it is a statement of fact about the provider, not a user preference: free accounts receive no included Inference Providers credits (docs change merged 2026-10-07).

### 10.3 Path coverage — hosted cannot be "accidentally selected"

| Path | Outcome | Check |
|---|---|---|
| CLI `routes` | Prints `hosted is REFUSED …`; no free-allowance advertising remains | 61 |
| CLI `run` | Still exits 2 at the grant gate; no hosted construction anywhere in `cli.py` | 51, 62 |
| Dispatcher | Constructs only what it is given; `cli`/`dispatch` contain no direct adapter construction | 62 |
| Direct constructor | Refused at the object, independent of the factory | 59 |
| Refusal cost | Zero adapter calls beyond the deliberate probe; the refusal precedes any dispatch | 60 |

### 10.4 Validation

**62/62 passed, exit 0**, socket guard active, 0 scratch directories left. Checks 1–26 and 28–34 and 36–58 preserved; 27 and 35 rewritten for the refusal; 59–62 added.

### 10.5 Milestone state at close of P1.2

| Item | State |
|---|---|
| P1 offline instrument | **Accepted** |
| P1.1 offline hardening | **Complete, 58/58 at the time** |
| P1.2 hosted refusal | **Enforced, 62/62** |
| Hugging Face hosted route | **Unavailable** under the zero-cost constraint |
| Local execution route | **Next validation target — not yet approved as validated** |
| P2 hosted execution | **Not authorized** |
| Live adapter | **Unverified, unreachable code** |

### 10.6 What this does not do

- It does **not** validate the local route. Local-first is not local-validated: the local adapter's response handling still needs its own verification before any S1 behavioural testing.
- It does **not** delete the hosted adapter. The code is retained deliberately as an inspectable artifact; unreachable and unverified is its correct status.
- It does **not** decide whether hosted ever returns. That would require a paid plan or purchased credits, which is an owner decision about money.
- It does **not** change dispatch, the ledger, the grant mechanism, or the signature situation (§8.2 stands: integrity binding, not authentication).

### 10.7 Claim ceiling

A control now exists where previously only a recommendation did. The harness refuses hosted inference for a reason that is checkable, at three independent points, and a test fails if any of them is removed. Nothing here shows local inference works, and no model has been run.
