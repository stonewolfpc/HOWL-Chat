# `MEMORY-09` — P1 Build Completion Record

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-P1CR-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **P1 COMPLETE — instrument built, nothing executed. P2 remains unauthorised.**
**Authority:** owner decision, `BEA-1` / **P1 BUILD AUTHORIZED · P2 EXECUTION DEFERRED**.
**Current P1 manifest:** superseded by `MEMORY-09-P1-SOURCE-INSPECTION-AND-GAP-CLOSURE-2026-11-13.md` (row 485), which lists the 16 files after gap closure (51 checks, 94,608 bytes). The table below is the build-time snapshot and is retained verbatim as history, not as current hashes.

**Revision 1.1 (2026-11-13):** route policy added at owner direction — `budget.py` and `local_adapter.py` created, adapter factory and CLI extended, validation grown from 25 to 35 checks. Table and totals below are current as of this revision; the $0.00 paid ceiling is recorded in row 484.
**Scope check:** every file created lies under `eval/`. No Howl source, corpus, standard, scoring or authoritative item was modified.

---

## 1. Files created

| File | Bytes | SHA-256 |
|---|---:|---|
| `eval/__init__.py` | 168 | `d99d08b1213f74f44aaa5656737cfc2c6801aa6431cc6486f308096785370c17` |
| `eval/harness/__init__.py` | 292 | `3dced7dee1277d415c3b3daf645bf427462b2f475cea12877c859f600ad8396f` |
| `eval/harness/adapter.py` | 9,444 | `42be4408db2f1f01c1e55445252c89cc0b430679ebfc826698a975535c26d22b` |
| `eval/harness/authorization.py` | 3,273 | `46b750379b1305895fe840f002bf137803eeb79b75ceda39e2182436361dfff8` |
| `eval/harness/budget.py` | 3,964 | `e6eb0231ef0924d1f76ff98c2c70113957630ee861a2f3daa125b7d109798f3f` |
| `eval/harness/cli.py` | 4,310 | `e9e0b0a015e9604db01773a0b6a5cf19ae0576a0661c580b27b2a232a0ca32d0` |
| `eval/harness/conditions.py` | 4,365 | `ae3571e1de5b380012c0b1f7cca99bf48c6f30ac8ffd62c220297c33c0b62fff` |
| `eval/harness/fixtures.py` | 5,424 | `2a7ad5a53d03b4ac2effefdcfdb4599f92493e8c8a127f9f1e0dcb4347db54b3` |
| `eval/harness/local_adapter.py` | 3,877 | `07657261d0c266e81e950af15f1285881c24a426151a84ba1672ba06f2495ac6` |
| `eval/harness/plan.py` | 3,917 | `640988f9003dea2db83588d28bd9731a078a8ccdba2dbd1ab6c93e409d5c1396` |
| `eval/harness/provenance.py` | 5,268 | `5d240b48c4de3c99dff73dc3d770649f720a28eb9ccb658e96af817e21e222e6` |
| `eval/harness/state.py` | 4,969 | `c51bafff082eb3c7785772c671876094eb1252936a3f26d1cdf491032e45f2a0` |
| `eval/harness/validate.py` | 19,807 | `8b5b9cbf8c9bfe85dd729a494e0446a2f15b5bbc70d20fc3ca0cee6604040461` |
| `eval/harness/README.md` | 3,850 | `d63cb7c4b29ca1ca4076fe4e9ec576054c4870f615b4e353b491438bc78b1c5e` |

Total 72,928 bytes across 14 files.

## 2. Validation performed

`python -m eval.harness.validate` → **35/35 passed, 0 failed, exit 0.**

| Group | Checks | Result |
|---|---:|---|
| Adapter rules (T-1, T-2, streaming, params, offline refusal, grant gate, no module-scope client import) | 8 | PASS |
| Authorisation (no P2 grant exists; `require_grant` refuses) | 1 | PASS |
| State (deterministic hash, immutable snapshot, diff, recorded removal, source validation, absent fact) | 5 | PASS |
| Fixture (JSON round-trip hash stability, declared guards catch tampering) | 2 | PASS |
| Conditions (single-line separation, widened difference refused) | 2 | PASS |
| Plan (18 calls / 30 ceiling, unassigned refusal, overrun refusal) | 3 | PASS |
| Provenance (`unavailable` markers, inadmissible without raw body, write/read) | 3 | PASS |
| Budget (zero ceiling, unverified refusal, paid tier, auto top-up, exhaustion, route change) | 6 | PASS |
| Local route (non-loopback refused, grant required, shared contract) and factory ordering | 4 | PASS |
| Cleanup (no self-test artifacts left) | 1 | PASS |

Additionally: `python -m compileall -q eval` succeeded, and `python -m eval.harness.cli inspect` printed the fixture, both conditions and both prompt hashes.

## 3. Network execution absent — evidence

1. **Socket guard:** all 35 checks ran with `socket.socket` and `socket.create_connection` replaced by raises. No connection attempt occurred.
2. **No client import at module scope:** check 1 asserts `huggingface_hub` is absent from `sys.modules` after the harness imports; the import exists only inside a method.
3. **No adapter can be constructed:** `HuggingFaceAdapter(grant=None)` raises `ExecutionNotAuthorised` (check 8).
4. **Offline adapter refuses:** `OfflineAdapter.generate` raises (check 7).
5. **No auth token read, no endpoint string, no HTTP library imported** anywhere in `eval/`.

**Honest limit:** the socket guard is a *validation instrument*, not a runtime sandbox. It proves no socket opened during validation; it does not make the interpreter incapable of network.

## 4. P2 cannot be triggered accidentally — three independent barriers

| Barrier | Mechanism |
|---|---|
| **Grant gate** | `ExecutionGrant` can only be produced from `eval/runs/P2-GRANT.json` with owner signature, all three slots named, pinned providers within the 30-call ceiling. That file does not exist and `eval/runs/` does not exist. |
| **Plan gate** | `build_plan` raises while any slot is `UNASSIGNED`; P1 ships all three unassigned, and naming them is an owner decision on the schedule. |
| **Structural** | `HuggingFaceAdapter.generate` ends in `NotImplementedError` — response transcription is P2 work. There is no code path that returns a provider response. |

`python -m eval.harness.cli run` → prints `REFUSED.` and exits **2**.

## 5. Defects found and fixed during the build

| # | Defect | Fix |
|---|---|---|
| 1 | My own check compared a list of change dicts to `[0]` | compare extracted indices |
| 2 | `Condition` had no `prompt_hash`, so the separation check crashed instead of testing | added the method, same formula as `GenerateRequest` |
| 3 | README documented `python -m eval.harness.inspect` — no such module | corrected to `python -m eval.harness.cli inspect` |

All three were caught by running the checks, and the suite was re-run to 25/25 after each fix.

## 6. Documentation correction (source-verified)

Signature inspection of the installed `huggingface_hub` **1.29.0** shows `InferenceClient.__init__` takes `provider` with default **`None`**, not the string `"auto"`. A harness accepting "anything but `auto`" would therefore route silently. Row 477 now forbids unset as well, and records the verification method. `chat_completion(model, messages, stream, …)` confirmed present.

## 7. Incomplete or unenforced requirements

| Item | State |
|---|---|
| `HuggingFaceAdapter` | **Written, never executed.** Not verified against a real response shape; transcription deliberately unimplemented |
| Cost ceiling / pricing pre-flight | **Not implemented** — P2 concern; the grant schema requires the field but nothing prices a call |
| `model_revision`, `system_fingerprint`, provider headers | Recorded as `unavailable`; population depends on P2 response data |
| Seed honouring | Always `unknown`; no provider evidence exists |
| Retry/backoff implementation | **Not implemented** — P2 |
| Socket guard | Validation-only; not a runtime sandbox |
| Packaging, `.gitignore`, CI | Not added; out of P1 scope |
| Fixture coverage | One fixture (`F-ST-01-S1`), as authorised. F-ST-02…06 were not built |
| Model slots | Unassigned by design |

## 8. What P1 does not establish

Nothing about model behaviour, memory concept validity, provider consistency, or `MEMORY-09` effectiveness. It establishes that a compliant instrument exists, that its rules are mechanically enforced, and that execution is unreachable without a new owner decision.

## 9. Registry

Row 483. `check_corpus.py` unaffected (no corpus change); registry metadata re-verified after insertion.
