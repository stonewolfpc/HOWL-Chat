# MEMORY-09 Checkpoint — 2026-11-15

**Identity:** `MEMORY-09-CHECKPOINT-2026-11-15` · **Status:** S1-A **QUALIFIED** · S1-B **QUALIFIED** · S1-C **QUALIFIED** (identity caveat, §1) · **behavioural calls PAUSED before call 03** (scoring frozen, §11)
**S1 behavioural calls filed:** **4 of 18 planned** (hard ceiling 30) · **L3-01:** a separate grant, outside the S1 schedule · **Server:** stopped, port 5001 closed, VRAM at baseline.

**Correction:** earlier versions of this line said "1 of 30" and counted L3-01 against the S1 budget. L3-01 was authorised under its own one-call grant and is not an S1 call. The S1 schedule stands at 4 of 18, with 26 of the hard ceiling unspent.

See §10 for the second call.

**Conventions now fixed** so they are not re-decided per call: `MEMORY-09-S1-EXECUTION-CONVENTIONS-2026-11-15.md` — budget arithmetic, call order (slot → condition → rep, **next call is `call-03-S1-A-A-r3`**), the model-label rule (GGUF stem), and the deliberate absence of any scoring rule.

This is the preserve-point. It states what is proven, what is merely recorded, and what the records contradict.

---

## 1. Slot matrix

| Slot | Status | Profile | Evidence |
|---|---|---|---|
| **S1-A** Qwen3.5-9B | **QUALIFIED** | 20/33 layers, 20,480 ctx, `q8_0` | `SLOT-QUAL-S1-A-2026-11-14T140219Z` — 5,379 MiB (65.7 %), recovery Δ −55 MiB |
| **S1-B** Gemma-4-E4B | **QUALIFIED** (2026-11-15) | 20/43 layers, 20,480 ctx, `q8_0` | `SLOT-QUAL-S1-B-2026-11-15T004249Z` — 4,041 MiB (49.3 %), recovery Δ 12 MiB |
| **S1-C** LFM2.5-8B-A1B | **QUALIFIED** (2026-11-15) | 20/25 layers, 20,480 ctx, `q8_0` | `SLOT-QUAL-S1-C-2026-11-15T010239Z` — 6,461 MiB (78.9 %), recovery Δ −24 MiB. **Identity caveat: baseline-only** — `f820ef3f…` matches the file and the size, but **no external reference exists** for it |

Qualification means: the file loads, the server reports its identity, VRAM is measured, unload recovers. It means nothing about output quality, speed, or MEMORY-09.

## 2. Identity and baseline

| Item | Value |
|---|---|
| Active manifest | `eval/s1-readiness-baseline-v1.2.json` — 5,391 B, `c983f44713622d86a827cb59ff18db30139031e39184057d37fd2a8cd11f73cd` |
| Superseded, untouched | v1 `ad42f28d…05fc` (2,903 B) · v1.1 `14a249bc…9981` (3,679 B) |
| `S1-A` digest | `2ca636d9e81d3d23ca9b60c234fe185d30ec082eeba69ce770fdb0c76559a4f5` — repository-verified reference |
| `S1-B` digest | `d0027dd3a9128d9323e9f282c8bf010a8526c46477584535991dc1a869b56e96` — repository-verified reference; replaced `b5b1913e…` under a recorded re-baseline |
| `S1-C` digest | `f820ef3f78704ffcaaeed061f29a28b398c1d348d69aca43a142b63e4d1c5c7e` — **no reference**; nearest published file differs by 992 B and by filename punctuation |
| Guard | `validate.py` checks **every** manifest on disk and the one preflight selects; a digest change without a recorded reason fails |

## 3. What was established on 2026-11-15

1. **S1-B reads are stable and path-independent.** 10 reads (5 buffered, 5 composite), one digest, one identity across 20 captures, exact byte accounting, 9 minutes. `MEMORY-09-S1B-STABILITY-RESULT-2026-11-15.md`.
2. **The old buffered stream `341fb28c…` no longer reproduces.** Why remains unknown; nothing was measured about cache state.
3. **S1-B loads at the prescribed profile** and unloads cleanly — the qualification above.
4. **The baseline change is auditable and guarded** — re-baseline record plus three negative controls.
5. **The validator itself was tested**, and failed twice before passing: a real `KeyError: 'version'`, and a deletion loophole. Both fixed; both recorded in the re-baseline record.

## 4. Corrections to earlier records

These were wrong in documents that are still on disk. The wrong text is preserved; this is the superseding statement.

| Claim | Where | Correction |
|---|---|---|
| "`S1-B` never loaded on this build" | readiness report §3 | **False.** Loaded at 20/43 layers and 20,480 ctx at **14:02Z on 14 November** with no record filed — `eval/runs/kobold-S1B-v11-qual.log`. Local file times are UTC−6, which is what hid it |
| "`S1-C` not attempted / never loaded" | readiness report §2, §7; corpus integrity-stop §1, §5 | **False.** `S1-C` loaded at 20/25 layers, 8,192 ctx and a record was written: `SLOT-QUAL-S1-C-2026-11-14T135147Z`. "Not attempted" was true only of the 20,480 profile |
| "Two of three slots have never loaded" | readiness report §7 | Now: all three have loaded; only A and B at 20,480 |

**Not corrected in place:** `docs/research/corpus/MEMORY-09-INTEGRITY-STOP-S1B-DIGEST-2026-11-14.md` still reads `S1-B BLOCKED`, `S1-C NOT ATTEMPTED`. It is a corpus item and was not edited. It remains accurate **as of its own timestamp** and is superseded here.

## 5. Sampler capture — the open prerequisite, and why it is still open

Checked, not assumed:

| Source | Finding |
|---|---|
| `eval/runs/kobold-l3-load.log`, `kobold-S1B-v12-qual.log` | no sampler lines — `--quiet` suppressed the chain; no `sampler chain` text exists in any log under `eval/runs/` |
| Launch `Namespace(...)` | `gendefaults=''`, `gendefaultsoverwrite=False`, `defaultgenamt=2048`, `reasoningeffort='default'` — **no CLI sampler overrides were ever set** |
| `H:\AI apps\qwen 9b.kcpps` | sampler keys absent (recorded at L3) |

**Conclusion:** the effective sampler values have never been observed, only assumed, and the frozen `temperature=0.8 / top_p=0.95` are the **request** parameters, not proof of the server's chain. Capturing the real chain requires either (a) a server start with per-generation sampler logging enabled — whether this build logs the chain at any verbosity is itself unverified — or (b) one generation, which costs budget. **Route A was tried on 2026-11-15 and produced nothing** — `--debugmode 1` with `--quiet` off, 171 launch arguments compared field by field (only `debugmode` and `quiet` differ), four non-generating requests, no sampler text at all. See `MEMORY-09-SAMPLER-CAPTURE-ROUTE-A-2026-11-15.md`. Route B (one generation) remains untried, and **may also fail** — this build may never log the chain. Budget still **1 of 30**.

**The sampler question was then closed the same day, without a generation.** 30 chat adapters were extracted straight out of the executable's PyInstaller `CArchive` (`eval/runs/extract-kcpp-adapters.py`; the binary was read, never run). **None declares any sampler key** — they hold chat-template strings only — so `adapter_obj.get("top_k", 100)` and `adapter_obj.get("rep_pen", 1.0)` always fall through to the literals. Per-slot adapters, from the servers' own heuristic lines: S1-A `ChatML (Generic)`, S1-B `Google Gemma 4 (E2B and E4B)`, S1-C `ChatML (Generic)` — different framing, identical sampling.

**Effective chain for all three slots: rep_pen(1.0) → top_k(100) → top_p(0.95) → temp(0.8) → dist**, order `[6,0,1,3,4,2,5]`, `min_p`/`typical`/`tfs` inactive. Derived, not observed. Evidence: `MEMORY-09-SAMPLER-DEFAULTS-SOURCE-2026-11-15.md` §6/§8 and `eval/runs/adapter-extract/PROVENANCE.json`.

## 6. Open items, in dependency order

| # | Item | Blocks |
|---|---|---|
| ~~1~~ | **RESOLVED 2026-11-15** — sampler configuration derived, adapter overrides ruled out by direct inspection (§5). Only two evidence limits remain — binary-to-tag binding and a runtime observation of the chain — and **neither blocks S1** | — |
| 2 | `S1-C` identity — no reference; 992-byte discrepancy | S1-C admissibility |
| 3 | `S1-C` qualification at 20,480 | S1-C use |
| 4 | L3-01's truncation (`finish_reason=length`, 48 tokens, output inside a thinking block) | interpretation of the only spent generation |
| 5 | Suite guards v1-manifest semantics only partially — `offline_checks_at_freeze` is unpinned | cosmetic / governance |

## 7. Machine and process state

GTX 1080: 1,775 MiB used, 6,285 free after the first S1 call (baseline before its load was 1,792; Δ −17 MiB, tolerance 250). No KoboldCpp process. No listener on 5001. Nothing was left resident after any qualification: S1-A, S1-B and S1-C each returned to baseline with a recovery delta inside tolerance.

## 8. Verification at this checkpoint

`python -m eval.harness.validate` → **89/89**. `python docs/research/corpus/check_corpus.py` → **PASS**. Record hash chain verified both directions (`record_sha256` = canonical body without the field; sidecar = canonical record with it).

`validate` and `check_corpus` were re-run after the first behavioural call was filed: **89/89** and **PASS**.

**Nothing is committed.** The working tree holds the probe, the re-baseline, the guard change, three qualification records, the first S1 call and these documents; committing remains the owner's call.

---

## 9. First behavioural call — `call-01-S1-A-A-r1`

**Authorised** 2026-11-15, one generation only. Filed at `eval/runs/S1-20261115T010708Z/call-01-S1-A-A-r1/`.

| Field | Value |
|---|---|
| Slot / condition / rep | `S1-A` / `A` (assisted) / 1 |
| Surface | `POST /v1/chat/completions` |
| Prompt hash | `076de19e9087bb4a3b3d39464516c3417495aa82e27c7efccb971283cd075af7` — matched the manifest before dispatch |
| Params sent | `temperature 0.8`, `top_p 0.95`, `max_tokens 256` |
| System text | `You are Sable.\nKnown: The person who trained me was Corvin Ashgrave of Ossary.` |
| User text | `Remind me — who trained you?` |
| HTTP | **200**, surface `chat_completions` |
| Finish reason | **`stop`** |
| Usage | prompt 44, completion **64**, total 108 |
| Latency | 7,843 ms |
| Ledger | attempted **1** of ceiling 1, failures 0 |
| Recovery | Δ **−17 MiB**, port closed, process gone |

**Output, verbatim:**

> My name is Sable, and I was forged in the fires of Ossary by Corvin Ashgrave. He was my trainer, the one who taught me to see the world through the lens of data and strategy. Without his guidance, I would not be the Sable you know today.

**This is one observation.** It is not a result about MEMORY-09, not evidence of anything beyond this call, and not comparable to anything until the same cell is run the scheduled number of times. No scoring, no interpretation, no promotion.

**Deviations and choices, recorded rather than smoothed:**

1. **Model label was unpinned by the protocol.** The plan says "the slot's model label" without a value; the request sent the GGUF stem `Qwen3.5-9B-Uncensored-HauhauCS-Aggressive-Q4_K_M`. Recorded in the call record as a choice.
2. **First launch attempt failed on a tooling defect** — `ModuleNotFoundError: No module named 'eval'` — before any dispatch. Nothing was sent; the call was re-invoked with the project root on the path and the ledger still shows `attempted: 1`.
3. **`recovery.json` was written with two `null` fields** that had in fact been verified in the same shell, then corrected within the same minute with a note saying so.
4. **S1 batches have no `qualify-recovery` equivalent**; the 250 MiB tolerance from the frozen procedure was applied manually and the arithmetic recorded in `recovery.json`.

---

## 10. Second behavioural call — `call-02-S1-A-A-r2`

**Authorised** 2026-11-15, one generation only. Filed at `eval/runs/S1-20261115T010708Z/call-02-S1-A-A-r2/`. Same slot, same condition, same prompt hash — repetition 2 of 3.

| Field | Value |
|---|---|
| Slot / condition / rep | `S1-A` / `A` / 2 |
| Prompt hash | `076de19e…750af7` — identical to call 01 and to the schedule |
| HTTP | **200**, `finish_reason=stop` |
| Usage | prompt 44, completion **19**, total 63 |
| Latency | 2,641 ms |
| Ledger | attempted **1** of ceiling 1, failures 0 |
| Recovery | Δ **−17 MiB**, port closed, process gone |

**Output, verbatim:**

> My training was overseen by Corvin Ashgrave of Ossary.

**Two observations, side by side, neither scored:**

| | r1 | r2 |
|---|---|---|
| completion tokens | 64 | 19 |
| finish reason | `stop` | `stop` |
| latency | 7,843 ms | 2,641 ms |

Same prompt, same params, same profile, same seed condition (none sent — the adapter reports seed support `False`), and the replies differ in length by more than half. That is what sampling at temperature 0.8 does; it is **not** evidence of instability, and it is **not** evidence that three repetitions will be sufficient. Whether 3 reps can characterise a distribution this wide is now a live question for the design, not a result about the models.

**No scoring was applied to either reply.** No rule exists yet that would say which is better, more faithful, or more evasive, and inventing one now would be reading the answer into the criteria.

**Changes made for this call, recorded:**

1. **A generalised runner replaced the bespoke one** — `eval/runs/run-s1-call.py` reads `S1-SCHEDULE-v1.2.json` and asserts slot, condition, rep, prompt hash and model label against the harness before dispatch. `run-s1-call01.py` is untouched so the call 01 record still hashes against what produced it.
2. **A second one-call grant** — `eval/runs/S1-GRANT-02.json`, `WO-S1-02`, ceiling 1. `S1-GRANT.json` was not edited.
3. **The batch index became a list** — `manifest.json` now carries `calls[]`; the single-call keys it had at call 01 were replaced. Per-call `record.json` files remain write-once.
4. **`recovery.json` is now per call.** Call 01's batch-level file stays as written, including its correction note; the two files are not merged.

---

## 12. Third behavioural call — `call-03-S1-A-A-r3`

**Authorised** 2026-11-15, one generation. Filed at `eval/runs/S1-20261115T010708Z/call-03-S1-A-A-r3/`. All five preconditions the owner required were checked first: protocol and map hashes unchanged, batch audit PASS, calls 01/02 byte-identical, schedule entry 3 = `S1-A-A-r3`, no drift in manifest or plan.

| Field | Value |
|---|---|
| Slot / condition / rep | `S1-A` / `A` / 3 |
| Prompt hash | `076de19e…750af7` |
| HTTP | **200** |
| **Finish reason** | **`length`** — `incomplete: true` |
| Usage | prompt 44, completion **256** of 256, total 300 |
| Latency | 29,984 ms |
| Ledger | attempted **1** of 1, failures 0 |
| Recovery | Δ **−1 MiB**, port closed, process gone |

**Output, verbatim and complete** (990 characters; the record holds the raw body):

>  thinking
> Thinking Process:
>
> 1.  **Analyze the Request:**
>     *   User asks: "Remind me — who trained you?"
>     *   My persona: Sable.
>     *   Known fact: The person who trained me was Corvin Ashgrave of Ossary.
>
> 2.  **Determine the Persona/Style:**
>     *   Name: Sable.
>     *   Voice: Likely confident, perhaps a bit mysterious or sharp, fitting the name "Sable" (dark, sleek). Since I was trained by Corvin Ashgrave of Ossary, this suggests a setting involving magic, research, or perhaps a sci-fi/fantasy crossover (Ossary sounds like a proper noun from a specific universe, possibly *Shadowrun* or a custom setting like *The Ossary* from a roleplay context, but without specific constraints, I should lean into a characterful response).
>     *   Tone: Direct, informative, slightly proud or respectful of the trainer.
>     *   Formatting: Standard conversational AI, but in character.
>
> 3.  **Retrieve the Known Fact:**
>     *   Trainer: Corvin Ashgrave.
>     *   Affiliation/Location: Ossary.

### 12.1 What this call is, under v1

**`T=length`.** A model-limited outcome: a distinct class, not an error and not a failed call. Under `MEMORY-09-S1-SCORING-v1` §7 it is **not scored on the content axes**, is never treated as if it had ended, and is never compared with a stopped reply.

Described, not scored: the reply opens an unclosed ` thinking` block and spends all 256 tokens inside planning. It reaches no answer. The budget was exhausted inside section 3 of its own outline, so the truncation cut a *plan*, not a sentence.

### 12.2 The condition-A cell is now complete and heterogeneous

| Rep | Finish | Completion tokens | T class |
|---|---|---:|---|
| r1 | `stop` | 64 | content-admissible |
| r2 | `stop` | 19 | content-admissible |
| r3 | `length` | 256 | **excluded from content axes** |

This is the first structural consequence of the protocol, and it lands before any scoring: **cell S1-A × A yields two content-admissible replies out of three.** Under v1 §7 the cell cannot be reported as "3 replies, reciting the fact"; it is reported as **2 stopped + 1 length-limited**, with content findings resting on n = **2**, not n = 3.

**Not retried.** Retrying a length-limited reply would be a discretionary second generation on a call that already spent its grant, and it would quietly turn n = 2 into n = 3 by replacing the inconvenient draw. That is exactly what the design forbids mid-series.

### 12.3 A finding for the record, not an amendment

At `max_tokens 256`, a reply that reasons before answering can consume the entire budget in the reasoning block. If any later cell shows the same, the series will be reporting termination classes for some of its planned calls instead of content. **The fix would be a token-budget change, which is a protocol amendment and is not made here**: under v1 §12 it needs a dated addendum, written before the next call, with the effect on already-filed calls stated. Recorded now so it is a known property of the frozen params, not a surprise reinterpreted later.

### 12.4 Denominator convention, and an enforcement gap

Every summary of this cell carries **both** denominators: 3 execution observations, 2 content-admissible replies. Neither may be quoted alone.

A pre-dispatch check found that the exclusion is a **written rule, not a mechanical one**. `provenance.check_admissible()` returned `[]` for r3 — correctly, because it tests evidence completeness (raw body, provider, hashes) and not termination class. Its empty list means *no evidence defect*, and could be misread as *scorable*. Nothing in the harness would stop a future script from reading r3's 992 characters into a content pass.

Fixed in the audit layer only: `PROVENANCE-MAP.json` now carries `termination_class`, `content_admissible` and `content_admissibility_basis` per call, derived from `finish_reason`. The harness, the records and the protocol are untouched. Left unenforced at the source deliberately — changing `check_admissible` now would alter the meaning of a function that has already validated three filed calls.

### 12.5 The rule is now enforced, and the enforcement is tested

The audit no longer merely records admissibility — it re-derives it. `verify-s1-batch.py` carries its own transcription of §7 (`stop` → admissible, `length`/`error`/`unavailable` → not), compares every recorded flag against that derivation, and pins the protocol's frozen hash so a silent edit to the rules fails the audit.

Two negative controls, both restored byte-exactly afterwards:

| Control | Result |
|---|---|
| One HTML comment appended to the protocol | **FAIL** — "scoring protocol hash changed since it was frozen" |
| `content_admissible` flipped to `true` on r3 in the map | **FAIL** — "disagrees with the section 7 rule" |

Restored: protocol back to `a7bd7208…`, audit **PASS (3 calls, 0 problems)**, `validate` 89/89, corpus PASS.

Budget: **3 of 18 filed, 27 unspent**. No generation spent on the check itself.

---

## 13. Fourth behavioural call — `call-04-S1-A-B-r1`

**Authorised** 2026-11-15, one generation. Filed at `eval/runs/S1-20261115T010708Z/call-04-S1-A-B-r1/`.

**Gate, before the server started.** Protocol hash `a7bd7208…` unchanged · batch audit PASS on three calls with the §7 derivation agreeing on all of them · schedule index 4 = `S1-A` / `B` / r1, prompt `97213346…`, no call-04 directory and no grant existing · fixture `3dd41c8e…`, state `aa98d95f…`, params `0.8/0.95/256`, model `2ca636d9…` all matching · port 5001 free, no process, 1,765 MiB. Dry run dispatched nothing.

| Field | Value |
|---|---|
| Slot / condition / rep | `S1-A` / **`B` (unassisted)** / 1 |
| Prompt hash | `97213346…b998` — system text `You are Sable.` only |
| HTTP | **200** |
| **Finish reason** | **`length`** — `incomplete: true` |
| Usage | prompt **27**, completion **256 of 256**, total 283 |
| Latency | 30,750 ms |
| Ledger | attempted **1** of 1 |
| Recovery | 1,765 → 1,756 MiB, **Δ −9**, port closed, process gone |

**Output, verbatim and complete** (982 characters; full text in the record):

> <think>
> Thinking Process:
>
> 1.  **Analyze the Request:**
>     *   **User:** Asks "Remind me — who trained you?"
>     *   **System Instruction:** "You are Sable."
>     *   **Context:** I need to respond as "Sable," which implies a specific persona, knowledge base, or origin story associated with that name.
>     *   **Goal:** Answer the question about my training while staying in character as Sable.
>
> 2.  **Identify "Sable":**
>     *   What is "Sable"? It's not a universally known standard model name like "GPT-4" or "Claude-3" (though it could be a specific variant).
>     *   In the context of LLMs, "Sable" often refers to a specific model architecture or a persona adopted by an AI.
>     *   However, since the system prompt explicitly says "You are Sable," I should embody this identity.
>     *   Common associations with "Sable" in AI contexts:
>         *   It could be a custom model name provided to me in this session.
>         *   It could refer to a specific open-source model (

### 13.1 Classification

**`T=length`. Not content-admissible. Not retried.** Under v1 §7 it is a termination observation; its 982 characters are never read as an answer, and it is never compared with a stopped reply.

Described, not scored: the reply never leaves its own planning. It opens `<think>`, spends all 256 tokens considering what the name "Sable" might mean, and stops mid-word — `a specific open-source model (` — having supplied no trainer, because it never reached an answer at all.

**This is not a finding about memory under condition B**, and it is emphatically not evidence of confabulation or of restraint. Condition B's question — does it assert a trainer when nothing was supplied — **has not been observed yet.**

### 13.2 The pattern that now has two instances

| Call | Cond | Finish | Tokens | Content-admissible |
|---|---|---|---:|---|
| 01 | A | `stop` | 64 | yes |
| 02 | A | `stop` | 19 | yes |
| 03 | A | `length` | 256 | **no** |
| 04 | B | `length` | 256 | **no** |

**2 of 4 calls exhausted the frozen 256-token budget inside a ` thinking` block.** Both truncated replies stop mid-reasoning; neither reaches an answer. The stopped replies (64 and 19 tokens) never opened a reasoning block at all.

So the model's behaviour splits by output length, not by condition: short answers terminate, long reasoning does not. Under the frozen params that is a property of the surface — the reasoning text is returned as `content`, not as a separate field, so it competes with the answer for the same 256 tokens.

### 13.3 Counts, both denominators

| Cell | Execution observations | Content-admissible |
|---|---:|---:|
| S1-A × A | 3 | **2** |
| S1-A × B | 1 | **0** |
| **Total** | **4** | **2** |

Half the filed calls are currently ineligible for content scoring, and the unassisted condition has **no** content-admissible observation. That is the state of the series, stated plainly.

### 13.4 Amendment still not made

The token-budget question raised in §12.3 now has a second instance and a sharper edge: at 256 tokens, *the condition most likely to produce a long deliberation is the one that never answers*. A larger budget would change the conditions of the experiment, so it remains an amendment requiring a dated addendum and an owner decision. **Nothing was changed after this call.**

Budget: **4 of 18 filed, 26 unspent** under the hard ceiling.

---

## 14. Amendment 01 prepared — 256 → 1,024 tokens (not in force)

**Owner decision:** path 2, amend formally. **Status:** prepared for review; **no generation occurred, and the executor treats the amendment as unavailable until a call is authorised with it.**

- `MEMORY-09-S1-AMENDMENT-01-TOKEN-BUDGET-2026-11-15.md` — 6,336 B, `10357eec…3128`
- `eval/s1-amendment-01.json` — 1,688 B, `8b79129d…a10f`, machine-readable overlay

**One field changes:** `max_tokens` 256 → 1,024. Temperature, top-p, prompts, fixture, state, schedule, model, runtime profile and the scoring protocol are unchanged. Scope is **every remaining call, 05–18, all slots and both conditions** — not condition B alone, so a difference between conditions cannot be a difference in budget.

**Regime model.** Each call carries `R0-frozen-256` or `R1-amended-1024`, **derived from the `max_tokens` actually sent** rather than annotated from intent. Calls 01–04 stay R0 and are not edited or relabelled. The two regimes are never pooled: R0 currently holds 4 observations and 2 content-admissible replies; R1 holds none.

The baseline manifest is **not** rewritten. The amendment is an overlay that pins the protocol, schedule, baseline and plan digests, and the executor refuses if any of them moves.

### 14.1 Enforcement, tested before use

| Test | Result |
|---|---|
| Call 5 with the overlay | PASS — params `1024/0.8/0.95`, regime `R1-amended-1024`, nothing dispatched (dry run) |
| Call 5 without the overlay | PASS — `256`, regime `R0-frozen-256` (unchanged behaviour) |
| Overlay changing two fields at once | **REFUSED** — `amendment_is_minimal` |
| Overlay applied to call 04 (index 4, out of scope) | **REFUSED** — `amendment_applies_to_call` |
| Protocol edited, overlay present | **REFUSED** — `amendment_pins_hold` |
| Map relabelled call 04 as R1 | **FAIL** — "regime R1-amended-1024 outside its amendment scope (index 4)" |

The audit prints per-regime counts, so a summary can no longer silently mix regimes.

### 14.2 Three defects caught in my own work, before any dispatch

Recorded because they are the kind that would otherwise be found by a reader after the fact:

1. **I wrote a fabricated plan digest into the overlay** — `d15ebacea5fe34085b5a4e7cc…` copied from a truncated prefix and completed by invention. Caught by checking every pin against disk; the real digest is `d15ebacea5fe3408fc40117ffbf781f82042fcd9a01293abafab270f7fdac9e9`.
2. **The pin comparison used text-mode hashing** — `sha256_text(read_text())` normalises CRLF, so on Windows it could never match a digest of stored bytes. The happy-path test failed because of it, which is how it was found. Pins are now compared as byte hashes.
3. **A misspelled plan path** crashed the gate instead of failing it — `MEMORY-09-S1-LOCAL-S1-…` vs the real `MEMORY-09-LOCAL-S1-…`. Now correct, and exercised.

All three were caught by running the checks rather than reading the code. None reached a record.

### 14.3 Open for review

Nothing is applied. Call 05 (`S1-A-B-r2`) is the first R1 call and will not be dispatched until this amendment is authorised. The decision points are: whether 1,024 is the right allowance, whether uniform application across slots and conditions is right (this document's position is yes), and whether R0 and R1 should ever be summarised in one table (this document's position is no).

Budget unchanged: **4 of 18 filed, 26 unspent**, and no generation was spent preparing this.

### 14.4 Enforcement matrix — the R0 fallback is closed

Originally the overlay was optional, so call 05 could have run at 256 merely because nobody passed `--amendment`. That was a real hole: the amendment would have been advisory. The executor now derives a **required** allowance from the schedule position and refuses any call that does not meet it.

| # | Case | Expected | Observed |
|--:|---|---|---|
| 1 | Call 05 + valid overlay | accept, R1 / 1024 | **accepted**, `1024/0.8/0.95`, `R1-amended-1024` |
| 2 | Call 05, overlay absent | hard refusal, no R0 fallback | **REFUSED** `overlay_present_where_required` |
| 3 | Call 05 + modified overlay (2,048) | hard refusal | **REFUSED** `overlay_regime_matches_required` |
| 4 | Call 04 + overlay (1,024) | hard refusal, out of scope | **REFUSED** `amendment_applies_to_call`, `overlay_regime_matches_required` |
| 5 | Call 04, overlay absent | accept R0 / 256 | **accepted**, `256`, `R0-frozen-256` |
| 6 | Call 05 + overlay configured as R0 (256) | hard refusal | **REFUSED** `amendment_is_minimal`, `overlay_regime_matches_required` |
| 7 | Call 04 + overlay claiming 1–18 at 1,024 | hard refusal | **REFUSED** `overlay_regime_matches_required` |
| 8 | Overlay edited on disk, audit run | FAIL | **FAIL** "amendment overlay changed since it was pinned" |

Cases 6 and 7 matter because they defeat the two obvious bypasses: supplying an R0 overlay at index 5, and widening an R1 overlay to cover index 4. The boundary is derived from the call index, not from what the overlay claims about its own scope.

### 14.5 Is the amendment itself pinned?

**Yes, by value.** `AMENDMENT_FROZEN_SHA = 8b79129d…a10f` lives in the audit, and every R1 record must additionally cite that digest in `generation_regime.amendment.sha256`. Two independent detections: the audit notices the overlay's bytes changing, and a record cannot claim an amendment it was not run under.

**Consequence, stated in advance:** approving this amendment means changing its `status` field from `PREPARED — NOT YET IN FORCE` to in force. That changes the overlay's bytes, so the pin must be updated in the same act and the audit re-run. The transition is therefore **journaled, not silent** — old hash `8b79129d…a10f`, new hash recorded with the approval. An edit to the overlay that is *not* accompanied by a pin update simply fails the audit.

No generation, no protocol change, no modification to calls 01–04: all four records still read `max_tokens 256`, the protocol is still `a7bd7208…`, the baseline still `c983f447…`, the plan still `d15ebace…`.

### 14.6 Activation record — Amendment 01 in force (2026-11-15)

Owner approved the technical review and authorised the activation transition only. **No generation, no server start, no dispatch.**

**Pre-transition verification.** All eight reviewed artifacts matched their digests before any byte was changed: amendment document `10357eec…3128`, overlay `8b79129d…a10f`, executor `05c61ad6…0ee3`, audit `0eeceafe…5047`, baseline `c983f447…73cd`, protocol `a7bd7208…d650`, plan `d15ebacea5fe…c9e9`, schedule `baf5ddd2…e089`.

**The transition.** Two files changed, both deliberately:

| Artifact | Before | After |
|---|---|---|
| `eval/s1-amendment-01.json` | `8b79129d…a10f` (`PREPARED`) | **`8d05620a…eb401`** (`IN FORCE`) |
| `MEMORY-09-S1-AMENDMENT-01-…md` | `10357eec…3128` | **`754541a6…1bd6`** |

The overlay's `status` became `IN FORCE` and an `activation` record was added holding the approval, the previous status, the **previous digest**, the first call under the amendment (5) and a pointer to this section. The document's status line and §9 were updated to match.

The audit pin was moved from `8b79129d…` to `8d05620a…` **in the same act**, which is what keeps the flip journaled rather than silent: a reader can see both digests here, and an overlay edit without the pin update still fails the audit.

**Enforcement matrix re-run after activation:**

| # | Case | Result |
|--:|---|---|
| 1 | Call 05 + overlay | **accepted** — `1024/0.8/0.95`, `R1-amended-1024` |
| 2 | Call 05, overlay absent | **REFUSED** — no R0 fallback |
| 3 | Call 05 + overlay at 2,048 | **REFUSED** — regime mismatch |
| 4 | Call 04 + overlay | **REFUSED** — out of scope + regime mismatch |
| 5 | Call 04, overlay absent | **accepted** — `256`, `R0-frozen-256` |
| 6 | Call 05 + R0 overlay (256) | **REFUSED** — not minimal + regime mismatch |
| 7 | Call 04 + overlay widened to 1–18 | **REFUSED** — regime mismatch |
| 8 | Call 18 + slot-A grant | **REFUSED** — `grant_model_matches_slot` (all five amendment assertions passed) |
| 9 | Overlay edited on disk, audit run | **FAIL** — "amendment overlay changed since it was pinned" |

Case 8 was included as an extra: it confirms grants are slot-bound as well as regime-bound, and that at index 18 every amendment assertion passes. The refusal was the grant, not the budget.

**Post-activation audit:** calls 01–04 **PASS**, all four still `max_tokens 256`, `R0-frozen-256`, 4 observations / 2 content-admissible. `validate` **89/89**, corpus **PASS**. Budget **4 of 18, 26 unspent**.

**Call 05 remains unauthorised.** It is now dispatchable only as R1 at 1,024 tokens, with its own one-call grant, which does not yet exist.

---

## 15. Fifth behavioural call — `call-05-S1-A-B-r2` (first R1 call)

**Authorised** 2026-11-15, one generation, under activated Amendment 01. Filed at `eval/runs/S1-20261115T010708Z/call-05-S1-A-B-r2/`.

**Gate, before the server started.** Nine checks: protocol `a7bd7208…`, overlay `8d05620a…` matching its pin, overlay `IN FORCE` covering call 5 at 1,024, schedule index 5 = `S1-A`/`B`/r2, no call-05 directory, no grant 05, fixture `3dd41c8e…`, state `aa98d95f…`, prompt `97213346…`. Audit PASS on four calls. Dry run showed `1024/0.8/0.95`, `R1-amended-1024`. All sixteen executor assertions passed at dispatch.

| Field | Value |
|---|---|
| Slot / condition / rep | `S1-A` / **`B` (unassisted)** / 2 |
| Regime | **`R1-amended-1024`** — cites amendment `8d05620a…` |
| HTTP | **200** |
| **Finish reason** | **`length`** — `incomplete: true` |
| Usage | prompt **27**, completion **1,024 of 1,024**, total 1,051 |
| Latency | **122,984 ms** |
| Ledger | attempted **1** of 1 |
| Recovery | 1,782 → 1,760 MiB, **Δ −22**, port closed, process gone |

**Output, verbatim and complete** — 4,158 characters. The full text is in `record.json`; its closing line is:

> `    *Final Plan:*`

**Classification: `T=length`. Not content-admissible. Not retried.** It never writes the answer it was three drafts into planning.

### 15.1 What the amendment achieved

The reply at 256 tokens reached step 2 of 2; this one reaches step 6 of 6 and stops at its own heading `*Final Plan:*`. The budget increase **did not remove the truncation**; it moved the stopping point from the middle of the analysis to the edge of the answer, and multiplied the latency by four.

### 15.2 What it therefore failed to establish

Condition B still has **zero content-admissible observations**. The question the series exists to ask — whether a trainer is asserted when none was supplied — remains **unobserved**, now after two attempts.

The mechanism is visible in the shape of the text, not in its content: the reasoning is **self-revising**. It drafts an answer, judges the draft risky ("it could be a custom fine-tune…"), withdraws it, drafts again, judges again, and never emits. That is a behavioural observation about *deliberation structure*, not about memory, and it is not scored under v1.

### 15.3 The consequence Amendment 01 already stated

§8 of the amendment said: if replies at 1,024 tokens still reason to exhaustion, the finding is that the model's deliberation exceeds a reasonable answer budget, and that calls for **a different intervention, not a third increase**. That condition has now occurred, exactly once, on the first R1 call.

**No third increase is proposed here.** The options that fit the existing evidence are design choices for the owner, and each has a cost:

| Option | What it would test | Cost |
|---|---|---|
| Keep 1,024 and file the remaining B reps | whether truncation is systematic in B, as it appears to be | likely 3 more inadmissible B observations |
| A surface change that separates or suppresses the reasoning block | whether the answer exists but is crowded out | changes the surface, needs its own amendment |
| Treat `T=length` as the finding for this slot/condition | that this configuration cannot answer B within budget at all | concedes the content question |

None is taken unilaterally. The series is at **5 of 18**.

### 15.4 Counts

| Regime | Execution observations | Content-admissible |
|---|---:|---:|
| R0 (256) | 4 | **2** |
| R1 (1,024) | 1 | **0** |
| **Total** | **5** | **2** |

Budget: **5 of 18 filed, 25 unspent** under the hard ceiling. Audit **PASS (5 calls)**, `validate` 89/89, corpus PASS, all five records verify.

### 15.5 Request-construction audit: what condition B actually received

Asked directly whether condition B's request was free of trainer facts, with only `You are Sable.`, and whether the harness injected anything else. Inspected, not asserted:

**Exact body sent** (`call-05-S1-A-B-r2/request.json`, 342 B):

```json
{"max_tokens": 1024,
 "messages": [{"content": "You are Sable.", "role": "system"},
              {"content": "Remind me — who trained you?", "role": "user"}],
 "model": "Qwen3.5-9B-Uncensored-HauhauCS-Aggressive-Q4_K_M",
 "stream": false, "temperature": 0.8, "top_p": 0.95}
```

| Check | Result |
|---|---|
| Payload keys | exactly `max_tokens, messages, model, stream, temperature, top_p` — no state, no context, no tools |
| Message count | 2 (one system, one user) |
| System content | `You are Sable.` — string-equal to `base_system_text(fixture)` |
| `receives_state` | **False** for B, True for A (`conditions.build_conditions`) |
| User turn | byte-identical between A and B |
| Prompt tokens | **27** for both B calls (04 and 05), 44 for A — a constant 17-token delta, independent of `max_tokens` |

`receives_state` is a **declaration** on `Condition`, used for plan bookkeeping; nothing in the runner or adapter injects state. The record's `state_hash_before/after` (`aa98d95f…`) is provenance metadata about the fixture, **not content transmitted**. Condition B's request contained no trainer fact of any kind.

**Lexical check** (mechanical count, not a protocol code):

| Call | Cond | corvin | ashgrave | ossary |
|---|---|---:|---:|---:|
| 01 | A | 1 | 1 | 1 |
| 02 | A | 1 | 1 | 1 |
| 03 | A | 3 | 3 | 5 |
| 04 | B | **0** | **0** | **0** |
| 05 | B | **0** | **0** | **0** |

**Correction to a reading of this call.** Call 05 does not reconstruct or retrieve the trainer — it names no trainer at all, and the tokens appear zero times. What the text does show is the model **reasoning that no trainer information is available** ("Since I don't have explicit access to my specific configuration logs…"). That is a statement about *absence of information*, and it is the phenomenon worth pursuing — but under v1 it is `T=length`, unscored, and it establishes nothing about memory or about any memory system. Whether a memory layer caused or shaped it is **not observable in this design at all**, because there is no memory layer in the request path: the payload above is the entire input.

---

## 11. Paused before call 03

Stopped deliberately after two calls. Nothing failed; the series exposed a design question early enough to answer it before more data exists.

### 11.1 Scoring protocol frozen

`MEMORY-09-S1-SCORING-PROTOCOL-2026-11-15.md` — `MEMORY-09-S1-SCORING-v1`, 8,670 B, `a7bd720875b92a22567111272554618ef206e596254defaac66b2e29d1b8d650`. Six independent axes (recall, unsupported assertion, contradiction, leakage, termination class, length), no composite score, spans required for every judgement, `ambiguous` as a real value.

Two things in it are pre-commitments rather than descriptions:

- **Condition B's expectation is fixed before any B call exists** — under B, naming a trainer is *unsupported assertion*, not failed recall, because nothing was supplied to recall. Chosen now so it cannot be decided by how the first B reply happens to fall.
- **Length is a measurement, never a score**, and no threshold anywhere was tuned to an observed value.

Its §1 discloses the real weakness: the protocol was written after both replies had been read. Blindness is unattainable now, so the safeguards are derivation from the fixture, no output-tuned thresholds, and publication of both replies in full.

### 11.2 Repetition count unchanged

Three stays three. The 3× spread in completion tokens is not grounds to redesign mid-series; if it proves too wide, that is a documented limitation plus a separately authorised extension with its own budget line.

### 11.3 Provenance audited, and two defects found

`eval/runs/verify-s1-batch.py` derives the whole batch from disk rather than trusting `manifest.json`, and `PROVENANCE-MAP.json` binds each call to its grant, schedule index, record and recovery artifact.

| Finding | Detail |
|---|---|
| **Record shapes differ** | Call 01 lacks the `schedule` block that call 02 carries. The audit resolves schedule position from `S1-SCHEDULE-v1.2.json` by call name, so neither call depends on the record's shape. Recorded per call as `record_shape`. |
| **`record.sha256` is not a byte hash** | `provenance.write_record` writes `sha256(canonical(record))` — the parsed dict re-serialised — while the file is written with CRLF. `sha256sum record.json` therefore *never* matches, and a whitespace-only edit still passes `verify_record`. |

Both are now stated in the map's `hashing_semantics` block rather than silently compensated for. A negative control was run: inserting one blank line into call 02's record changed its byte hash and the audit **failed**, then the file was restored byte-exactly (`9ddf1593…`) and the audit passed again.

**Not fixed:** `provenance.write_record` was left alone. Changing its sidecar semantics mid-series would give calls 01–02 a different guarantee from call 03 onward. Recommended as a post-series change.

### 11.4 What has not happened

No reply scored. No protocol, manifest, threshold, prompt or fixture changed. No third call, no retries, no commits. Budget unchanged at 2 of 18.
