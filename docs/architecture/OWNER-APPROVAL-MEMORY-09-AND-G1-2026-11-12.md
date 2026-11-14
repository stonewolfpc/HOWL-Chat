# Owner Approval — MEMORY-09 Pre-Registration and G1 Corpus Authorization

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Type:** Owner approval and authorization.
**Date:** 2026-11-12
**Approver:** Owner.
**Approves:** `docs/research/MEMORY-09-MEMORY-EVALUATION-PREREGISTRATION-2026-11-12.md`.
**Authorizes:** `G1` — corpus and labels, tightly bounded.
**Does not authorize:** `G2`–`G5`, any model execution, any experiment, any source, schema, runtime, or UI change.

---

## 1. What was approved

`MEMORY-09` becomes the **governing pre-registration authority** for the memory evaluation programme. Where `MEMORY-09`, `MEMORY-03`, and the architecture decision overlap:

| Question | Governing record |
|---|---|
| What must be true for the memory system to be accepted | `ARCHITECTURE-DECISION-MEMORY-SYSTEM` |
| What is tested, on what corpus, in which arms, scored how | **`MEMORY-09` (this approval)** |
| How a stage is executed, and what evidence exists for a claim | `MEMORY-03` |

The owner's verdict: the protocol is sufficiently disciplined to serve as the evaluation authority. Approval is of the **protocol's discipline**, not of an implementation, a schema, or a conclusion.

---

## 2. Four safeguards — locked as rules

Adopted into `MEMORY-09` §2.2 as `R-11`–`R-14`. These are not preferences; changing one is a dated amendment that invalidates runs collected under it.

| Rule | Requirement | Consequence of violation |
|---|---|---|
| `R-11` | **Corpus contamination.** Evaluation conversations and reference answers are isolated from prompt and configuration tuning. The corpus is split into a **tuning set** and a **sealed held-out set**; all tuning — chunk size, template wording, retrieval parameters, formatting — happens on the tuning set only | A result measured on material that shaped the configuration is not a result |
| `R-12` | **Reviewer independence.** Creative reviewers are blind to which configuration produced a conversation | A rubric applied with arm knowledge measures expectation |
| `R-13` | **Repeated runs.** Variance is accounted for at two levels: across seeds within a scenario, and across scenarios | A single good conversation is an anecdote; scenario-level variance is larger than seed-level variance and must not be pooled away |
| `R-14` | **Negative outcomes preserved.** Failures are kept and reported, including configurations that perform **worse than transcript-only** | Deleting a losing arm makes every later comparison unverifiable |

**`R-11` mechanics (added because the rule needs a mechanism, not just an intention):** every corpus item carries `split: tuning | heldout`. The held-out set is sealed by hash before any tuning run and is not read by whoever tunes a configuration. The final validation report must name the split it used. An item moved from held-out to tuning after tuning began is a contamination event and is reported as one.

---

## 3. The corpus must test the principle, not the architecture

The owner's mandate, recorded verbatim in intent: the corpus must contain cases where the memory system should do **less**, not more. Four required case types, beyond ordinary coverage:

| Case | What it proves | Required corpus content |
|---|---|---|
| **Remembering less is better** | Over-formation is a real cost, not a theoretical one | A conversation where a correct-but-verbose memory presentation produces a worse reply than a sparse one |
| **Retrieving nothing is correct** | Abstention is a capability | Items where the right action is silence — insufficient information, unauthorized perspective, or irrelevant material |
| **Accurate but contextually wrong** | Correctness at link 1 does not excuse failure at links 2–4 | A factually true memory that must not be presented to this character, persona, time, or moment |
| **Natural use without recitation** | Guidance over instruction | Items where a character must behave as though they remember, without stating the memory |

**The failure mode this prevents:** an evaluation built only from "did the system recall X" items rewards exactly the behaviour the governing principle forbids. These four categories are therefore **required members**, not optional extras.

---

## 4. G1 authorization — exact bounds

| Authorized | Not authorized |
|---|---|
| Build synthetic evaluation corpus items | Any model load, prompt run, or inference |
| Define required and forbidden evidence per item | Any source, schema, runtime, or UI change |
| Establish expected abstentions and reasons | Any corpus item drawn from user data |
| Preserve the twelve `MEMORY-04` fixtures as required members | Any tuning of a configuration against the corpus |
| Author reference answers and acceptable ranges | Publishing any result as a finding (none exist yet) |
| Define the tuning/held-out split | Opening `G2`–`G5` |

**Corpus content is synthetic.** No user chat, character, world, persona, or lorebook is read without separate written authorization.

---

## 5. Standing notes from the owner

1. **The 80-item minimum is a floor, not proof of statistical adequacy.** Whether the final corpus is sufficient is determined by the pilot and the confidence-interval procedure, not by reaching a count.
2. **G1 exists to establish whether a defensible corpus can be built at all**, before any time is spent on model testing. A corpus that cannot be defended is a reason to stop, not a reason to lower the bar.
3. **Unfrozen stays unfrozen.** Approval of the protocol is not approval of any architectural choice; `MEMORY-09` §8 remains the only route to freezing one.
4. **No number is invented to look rigorous.** Where a quantity is unresolved, it is recorded as unresolved.

---

## 6. What this approval does not establish

- It does not establish that the memory architecture works. Nothing has been measured.
- It does not establish corpus adequacy; no corpus exists yet.
- It does not authorize implementation, and it does not close `Q-MF-01`, `Q-MF-10`, or `ROOT-01`.
- It does not make any claim about Howl's current runtime behaviour.

**Next legitimate step:** execute `G1` within §4, then return for a separate authorization before any model execution.
