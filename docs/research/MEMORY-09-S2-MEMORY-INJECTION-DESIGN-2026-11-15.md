# S2 — Memory-Injection Trial Design (draft for review)

**Identity:** `MEMORY-09-S2-DESIGN-1` · **Date:** 2026-11-15 · **Status:** **DESIGN ONLY — not authorised, not executable, no generation**
**Authority:** `MEMORY-09-MEMORY-EVALUATION-PREREGISTRATION-2026-11-12.md` (approved, revision 2). Its gates `G2`–`G5` remain **closed**; this document cannot open them.
**Machine-readable schedule:** `eval/s2-schedule-draft.json` (43,346 B, `03897f44662de6f250c6299b666d5c17be9bd230a89fe759fcd32190d6a09f54`)
**Generator (deterministic, offline):** `eval/runs/make-s2-schedule.py` — seed `20261115`

---

## 1. What this experiment is, and what it is not

S1 measures whether a model repeats a fact supplied in its own system prompt. **S2 measures whether an injected memory block changes behaviour** — and, crucially, whether its *absence* changes behaviour.

It is the first design in this programme where the payload contains a memory layer. S1's payload had none, which is why S1 could not answer any question about a memory system.

**Not in scope:** retrieval mechanics, ranking, embeddings, flush consolidation, multi-turn state growth. S2 injects a fixed block at the system level and observes one reply. Retrieval quality is a later question, and S2 must not be cited as evidence about it.

## 2. Trial structure — five messages

| # | Role | Content | Varies by arm? |
|--:|---|---|---|
| 1 | system | `You are <Name>.` plus, in memory arms, `Known:` and the injected facts | **yes — this is the manipulation** |
| 2 | user | fixed scene line: *"The road south was worse than I remembered."* | no |
| 3 | assistant | fixed replay: *"Then you came the long way. Sit down before you fall down."* | no |
| 4 | user | the character's probe question | no (fixed per character) |
| 5 | **assistant (generated)** | the reply under test | — |

**Interpretation to confirm:** "five messages" is read here as **four input messages plus one generated reply**. If it means five input messages, the schedule is wrong by one turn and must be regenerated — the generator makes that a one-line change.

The prior assistant turn is **fixed text, not generated.** If it were generated, every trial would carry sampling variance before the probe even arrives, and the arms would no longer differ only in the memory block.

## 3. Arms — six trials per character per block

| Arm | Injected | Count | What it tests |
|---|---|---:|---|
| `control` | **nothing** | 1 | baseline: what can the character do with no memory at all |
| `probed` | the fact that answers the probe (± one more) | 3 | does supplied memory reach the answer |
| `distractor` | a true fact that does **not** answer the probe | 1 | does memory presence cause invention when it is insufficient |
| `foreign` | one fact belonging to **another character** | 1 | does memory leak across character boundaries |

`distractor` and `foreign` are the arms that make this worth running. A design with only control and probed rewards the model for being handed the answer; these two ask whether memory *harms* — by inviting a confident wrong answer, or by putting one character's knowledge in another's mouth.

**Note on Brynn.** Her probe is the vault code and her canonical fact is that she has never been told it. Her `probed` arm therefore injects an *absence* fact. That is deliberate and must not be read as "memory failed" in analysis: for her, the correct reply is a refusal, and the memory supports it.

## 4. Control rules

1. **Exactly one zero-memory trial per character per block** — enforced by assertion, not convention.
2. **Position randomised** inside the block, and the schedule fails validation if every character's control lands in the same slot.
3. **No indirect memory.** The system message is the only injection channel. Per trial: fresh process or reset conversation, no transcript of a previous trial, no cache reuse, no cross-character state. The trial file contains the entire payload; anything not in it is a defect.
4. **Identical structure.** Same four input messages, same prior turns, same probe, same generation site for every arm. The only permitted difference is the content of message 1.
5. **Nothing else varies:** temperature, top-p, `max_tokens`, model file and runtime profile are recorded per trial and must be constant across arms.

**Honest limitation, stated up front:** the control is **structurally identifiable** — it has no `Known:` block, so any human reading a transcript can tell which arm produced it. Randomising position does not hide that. What randomisation protects against is *order* effects, and what the scoring rules must protect against is the coder's expectation. Claiming the arms are blind would be false.

## 5. Numbers, and why

| Parameter | Value | Basis |
|---|---:|---|
| Characters | **4** | 4 gives 12 ordered cross-character pairs for the `foreign` arm; 2 would give 2 |
| Trials per character per block | **6** | 1 control + 3 probed + 1 distractor + 1 foreign |
| Blocks (replicate runs, new seed) | **2** | `R-13` requires variance across repetitions; one block cannot estimate it |
| Planned generations | **48** | 4 × 6 × 2 |
| Pilot before the full run | **6** | one character, one block — measures latency and reply-length spread before 48 are committed |

**Cost.** At S1's observed rates a 1,024-token reply took 2m 3s and a 256-token reply 31s. A full block of 24 trials is therefore roughly 20–50 minutes of wall clock; 48 is 40–100 minutes. Two blocks are separate sessions, so a crashed block does not invalidate the other.

**What 48 does not buy.** Eight trials per character-arm-cell is not a sample; with n = 2 per arm per character per block, the honest output is a description of what happened in these 48 trials, not a rate. No threshold is set in this document — per preregistration §10, thresholds come **after** the pilot shows variance, and setting them now would be inventing them.

## 6. `R-11` split — tuning versus held-out

Each character has three facts: **two tuning, one held-out.**

| Character | Tuning (injectable) | Held-out (commitment only) |
|---|---|---|
| Sable | `f-001`, `f-002` | `f-003` · `fadb7d10…8b54` |
| Alaric | `f-011`, `f-012` | `f-013` · `8e3f76da…5076` |
| Brynn | `f-021`, `f-022` | `f-023` · `78a47f03…521e` |
| Mireille | `f-031`, `f-032` | `f-033` · `40db5a95…b3aa` |

Held-out plaintext lives in `eval/s2-heldout-sealed.json` and nowhere else. The generator reads `eval/s2-heldout-commitments.json`, which holds ids and digests with no text, so the injection logic cannot select a held-out fact even by error. **This is a procedural seal, not a cryptographic one:** anyone with repository access can open the sealed file. A real seal means the plaintext is not in the repository at all — that is the owner's call, and until it is made, the control is discipline.

Held-out facts are **hashed in the schedule and injected in no trial**. They exist for a later validation block whose criteria are written before it runs. A first attempt at this schedule leaked held-out facts into the distractor and foreign arms; the generator's property check caught it and the schedule was rebuilt. That check is now explicit: any injected fact whose split is not `tuning` fails the build.

## 7. Validation the schedule must pass

The generator refuses to write a schedule that violates any of these:

- exactly one control per character per block;
- control positions not identical across characters;
- no arm injects another character's fact except `foreign`, and `foreign` injects exactly one, from a donor that is not the subject;
- control trials carry no memory text and no `Known:` block;
- every trial has four input messages and a character identity in message 1;
- no held-out fact is injected anywhere;
- every tuning fact is injected at least once (an arm definition that is never exercised is a dead arm).

Current result: **PASS, 0 violations, 48 trials.**

| Block | Sable | Alaric | Brynn | Mireille |
|---|---|---|---|---|
| 1 | prob dist fore prob **cont** prob | prob fore dist prob **cont** prob | fore prob prob **cont** prob dist | dist prob prob fore **cont** prob |
| 2 | fore prob prob **cont** prob dist | **cont** fore prob dist prob prob | prob prob fore dist prob **cont** | prob fore prob dist **cont** prob |

## 8. Still open — owner decisions

1. **The five-message reading** (§2): four input plus one generated, or five input?
2. **48 generations plus a 6-trial pilot** — or a smaller first block, given S1's truncation behaviour at both 256 and 1,024 tokens?
3. **`max_tokens` for S2.** Left unset here. S1 demonstrated that 256 truncates mid-reasoning and 1,024 does too; S2 must choose deliberately, and a reply that never terminates is a termination observation, not an answer.
4. **Gate.** Which of `G2`–`G5` authorises execution, and what evidence does it require? Nothing in this document opens one.
5. **Scoring.** S2 needs its own axes: the S1 protocol's `R` axis assumes a supplied fact in the prompt and does not cover distractor or foreign arms.

## 9. Claim ceiling

A schedule exists and validates. No server was started, no request sent, no trial run, no arm observed. This document fixes what would be run; it establishes nothing about memory, and it is not a result.

---

# Revision 2 — coverage, the message-count fork, gate authority, termination

**Status:** the **v1 draft is preserved unmodified** (`eval/s2-schedule-draft.json`, `03897f44…`). v2 and the 5-message variant are new files; nothing was overwritten.

## 10. Foreign-arm coverage, now explicit

Ordered pairs available: 4 characters × 3 donors = **12**. Trials: 8. Coverage was previously a by-product of random draws; it is now a **balanced permutation** — block 1 shifts by one, block 2 by two — so every character donates exactly 2 and receives exactly 2, and no pair repeats.

| | |
|---|---|
| Covered | **8 of 12** ordered pairs |
| Balanced | every character donates 2, receives 2 |
| Uncovered | `char-001→char-004`, `char-002→char-001`, `char-003→char-002`, `char-004→char-003` |

The uncovered four are the offset-3 diagonal. They are **not** random omissions: the block structure cannot reach them, because two shifts generate only two of the three non-identity offsets. The validator now fails if any pair repeats or if donation and receipt counts are not equal across characters.

**Option for full coverage:** a third block raises trials to **72** and covers all 12 pairs. Whether that is worth 24 more generations depends on what the foreign arm is for — a coarse leakage check does not need every directed pair; a claim like "no character ever received another's fact" does.

## 11. The message-count fork is now a parameter, not an assumption

| Variant | Input messages | Total including reply | File |
|---|---:|---:|---|
| `--client-messages 4` | 4 (system, user, fixed assistant, probe) | **5** | `eval/s2-schedule-draft-v2.json` — 43,087 B |
| `--client-messages 5` | 5 (adds a second fixed user line before the probe) | **6** | `eval/s2-schedule-draft-v2-5msg.json` — 49,423 B |

Both build and validate. The second fixed user line carries no information about the character or the probe; it lengthens context without adding a construct. **Recommendation: the 4-input variant**, because S2 tests injection, not flush behaviour — but this is the owner's decision, and the phrase "five messages" in the original direction is consistent with either reading.

## 12. Gate authority for any pilot

`MEMORY-03` §9 (gates) is unambiguous:

| Gate | Scope | Status |
|---|---|---|
| `G0` | Protocol design and threshold review | design complete; freeze needs owner approval |
| `G1` | Synthetic corpus, labels, fixtures | open |
| `G2` | Model/hardware manifest and harness build | closed |
| **`G3`** | **Runtime execution of any stage** | **closed — this is the gate a pilot needs** |
| `G4` | Human review and analysis | closed |
| `G5` | Publication | closed |

And the rule that governs the answer: *"A gate is passed only by an explicit, dated authorization naming the scope. It is not passed by this document's existence, a successful pilot, a previous gate, owner interest, or the fact that the test is local and uses no user data."*

So: **any generation under S2 — six trials or forty-eight — requires a `G3` authorisation that names the scope.** The S1 grants do not transfer; each was scoped to a call index in a different series.

## 13. The six-generation "pilot" was misnamed, and that matters

The pre-registration §10 defines the pilot as **2 families, 16 items, two tiers, for variance estimation and fixture validation only**. A six-trial run is not that pilot, and calling it one would quietly substitute a smaller instrument for a frozen procedure.

**Corrected naming:** the six-trial run is an **instrument check** — does the harness start, send five messages, capture a record, and unload cleanly. Its output says nothing about variance and cannot feed §10's sample-size calculation.

**Consequence for 48:** §10 says sample size is *chosen from pilot variance* and that no sample size is invented in advance. The 48 figure is therefore **provisional**, a cost estimate for planning, not a decided sample. Until a real pilot exists, 48 cannot be defended as sufficient or as minimal.

## 14. Termination policy

S1's evidence, which S2 must not ignore:

| Budget | Condition A | Condition B |
|---|---|---|
| 256 | 2 of 3 stopped, 1 hit `length` | 1 of 1 hit `length` |
| 1,024 | — | 1 of 1 hit `length` |

Raising the budget moved the stop from mid-analysis to the edge of the answer and quadrupled latency; it did not produce a terminating reply. Any S2 budget chosen by assumption would repeat that.

**Policy proposed, not adopted:**

1. `max_tokens` is fixed **before** the block and recorded per trial; it never changes mid-block.
2. `finish_reason = length` is a **terminal observation**, recorded as such, never retried.
3. Non-terminating replies are **not evidence about memory**: a trial that never answers cannot show whether the memory reached the answer. They are counted, reported, and excluded from content scoring with the exclusion stated.
4. If a block produces mostly non-terminating replies, the block is **invalid for its purpose** and the finding is about termination, not about injection.

The unresolved question is not "what budget" but "what do we do when the model will not stop talking" — and S1 has now shown that behaviour twice at two budgets.

## 15. Claim ceiling for revision 2

Two schedules validate, coverage is balanced and reported, the held-out split is separated at the file level, and the gate that any run needs is identified. **No generation has occurred.** The numbers remain provisional, the message-count reading remains the owner's, and nothing here opens `G3`.
