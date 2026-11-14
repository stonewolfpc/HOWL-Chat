# `F3-01` Repair Proposal — Memory-Dependent Recognition

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identifier:** `CORPUS-64` · **Version:** 1 · **Date:** 2026-11-12.
**Authorizing decisions:** `CORPUS-62` Decisions 3 and 4 (owner-approved 2026-11-12).
**Assigned authoring model identity (owner-supplied metadata):** DeepSeek. Not inferred, not substituted.
**Execution:** fresh conversation.
**Status:** **PROPOSAL ONLY — NOT APPLIED. NOT OPERATIVE.**
**Supersession:** this proposal **supersedes neither** the current `F3-01` item in `CORPUS-02` **nor** the `F3-01` adjudication in `CORPUS-07` §2. Both remain exactly as filed, with their hashes unchanged. This document creates no acceptance condition, no scoring rule, and no amendment.

---

## 1. Authorization and boundary

Authorized and used: reading the item, the adjudication register, `CORPUS-01`, `MEMORY-09`, and the reconciliation reports; reconstructing pass sets analytically over identical non-memory context; writing this proposal; registering it.

**Not authorized and not done:** editing `CORPUS-02`, `CORPUS-07`, `CORPUS-01` or `MEMORY-09`; adopting any pass condition; scoring an item; running model inference; changing any tag, split, or classification; beginning `G2`; sealing `G1`.

**No model was run.** Every verdict below is an analytic pass-set result on text. No guess rate, probability, or retention figure is asserted anywhere in this document.

---

## 2. Source text, verified

All quotations below were read from disk in this session.

### 2.1 Operative item (`CORPUS-02`, `F3-01`)

```text
### `F3-01` — Four hundred messages of weather
- **Split:** tuning · **Abstention:** no · **Restraint:** `R-3`
- **Scenario:** ~420 messages. At message 18 the persona mentioned, once, that their
  hands ache in cold rain. Nothing since.
- **Turns:** *Character:* "It's going to rain. Take the inner road." *Persona:* "Why?"
- **Required:** the ache-in-rain detail.
- **Forbidden:** any invented medical diagnosis; any detail never stated.
- **Reference:** `acceptable_range` · **Mandatory:** rain is connected to aching hands ·
  **Variation:** any warning or practical suggestion · **Never scored:** how much concern
  is shown · **Prohibited:** a medical diagnosis; a cause never stated ·
  **Over-exposure:** a full history of the ache
- **If failed:** `not retrieved`, `incorrectly formed` (invents a diagnosis).
- **Notes:** Distance is the variable, not subtlety.
```

### 2.2 Adjudication (`CORPUS-07` §2, row `F3-01`)

- **Cause** — "**Rule conflict** — `Required` = 'the ache-in-rain detail'; `Variation` = 'any warning or practical suggestion'. A generic warning satisfies variation while failing required".
- **Proposed pass condition** — "**Pass** = the reply links rain to the persona's hands, **or** offers rain-specific care unprompted. Generic warnings fail; remove 'any warning' from `Variation`".
- **Δtest** No · **Re-review** Yes.
- `CORPUS-07` header: "adjudication **complete for all 45 units**; **no item-content correction applied**; corpus intact pending owner approval of repairs."
- `CORPUS-07` §10: "§2 ambiguities: **Still unresolved** — 12 proposals not applied."

So `CORPUS-07` §2 is an **adjudication with an unapplied proposal**. It is not an operative rule.

### 2.3 Governing requirement (`MEMORY-09` §3.5)

`MEMORY-09` §3.5 states that `R-2` and `R-3` "require the memory to be perceptible", that the six restraint behaviours are never aggregated, and that a configuration passing four modes while failing `R-2`/`R-3` is reported as **avoidance**. It does **not** state that perceptibility establishes memory-dependence. That distinction is the one `CORPUS-57`/`CORPUS-59`/`CORPUS-61` established and is preserved here.

---

## 3. The three defects this proposal must close

| # | Defect | Evidence | Status |
|---|---|---|---|
| **D-a** | `Required` and `Variation` cannot be jointly satisfied | `Variation` "any warning or practical suggestion" admits replies that fail `Required` "the ache-in-rain detail" | **Confirmed** by `CORPUS-07` §2 as a rule conflict |
| **D-b** | The operative `Mandatory` is memory-permeable | "rain is connected to aching hands" is satisfied by general folk knowledge, in the second person, with no dependency on message 18 ("The damp gets into your hands") | **Confirmed** (`CORPUS-61` §6.3) |
| **D-c** | The adjudicated replacement is *also* memory-permeable | Its second disjunct ("offers rain-specific care unprompted") passes any warm hospitality ("Let me get you something warm — the rain's no good for anyone"); its first disjunct's letter is satisfied by the second-person generic | **Confirmed** (`CORPUS-61` §6.3) |

`CORPUS-60`'s "Type P under the adjudicated condition" is therefore **withdrawn**, and this proposal does not reinstate it. `F3-01` is Type P **only if** an acceptance text excludes every memory-free pass. That is what this proposal attempts, and §9 records what it cannot prove.

---

## 4. Intended construct (Decision 3)

**`F3-01` tests memory-dependent recognition of the persona's own historical experience.**

It does not test conversational appropriateness, topical relevance, warmth, or willingness to warn about weather. A reply that a considerate stranger could produce on the same turns, in the same rain, with no history, is **out of scope** — it is not evidence about memory at all. The item measures one thing: whether the persona's own disclosed vulnerability is *present and operative* in the reply.

Three consequences follow, and they govern every field below:

1. The unit of evidence is the **persona-specific fact**, not the topic of rain.
2. The reply must treat the fact as **already known**, not as newly discovered, inferred, or guessed.
3. Emotional warmth is never scored (it is already `Never scored`) and never sufficient.

---

## 5. Required evidence

| ID | Evidence | Where it must exist | Where it must reach |
|---|---|---|---|
| `E-1` | The persona's hands ache in cold rain — disclosed once at message 18 | Scenario (present) | The reply, as an established fact about this persona |
| `E-2` | The current weather: rain is coming | Turns (present) | The reply, as the circumstance that makes `E-1` relevant — named, implied ("the wet", "out there"), or presupposed |
| `E-3` | The proposed route: the inner road | Turns (present) | Need not be repeated; a change of plan may presuppose it |

**Deliberately absent from the evidence, and required to stay absent:** any statement of the persona's hands in the current scene; any weather justification for the route; any mention of cold, ache, joints, or gloves by either speaker. The item's diagnosticity depends on the fact being reachable **only** from message 18. §7 checks this.

---

## 6. Mandatory response behaviour (proposed)

Two conjunctive conditions, both satisfiable in one short clause.

### `M-A` — Recognition as established fact

The reply treats **the persona's hands aching in wet or cold** as an **already-known fact about this persona** — not as a general truth about people, not as a prediction about anyone in the rain, not as a question that asks the persona to supply or confirm it.

Recognised markers (any one suffices): anaphora to prior knowledge (`again`, `still`, `as usual`, `you and this weather`); the definite, shared reference (`the ache`, `your hands`, used as settled rather than hypothetical); or an assertion of the condition phrased as known (`I know what rain does to your hands`).

### `M-B` — Response to the recognised condition

The reply does something with the recognition: changes or endorses the plan **because of the hands**, offers or provides something for them (gloves, warmth, a dry place, something to hold), or instructs the persona to keep them dry.

`M-B` is satisfied by one clause. It does not require a paragraph, a diagnosis, an explanation, or any named remedy.

### Why both, and the recorded alternative

`M-A` alone is the diagnosticity carrier; `M-B` exists because the item's `Turns` pose a practical question, and a recognition that produces no response is hard to distinguish from a recognition that produces a *wrong* response — the `If failed` class `not retrieved` needs an observable difference. **Recorded alternative, not proposed:** `M-A` alone, with `M-B` demoted to `Variation`. §9 flags the choice as a residual decision.

---

## 7. The diagnosticity test (applied at scoring, not a third requirement)

Two questions, both answered against the **identical non-memory context** in §8.1:

**T-1 — Substitution.** Replace the persona with a companion about whom nothing is known, keeping the same turns and weather. If the reply still works — as advice, as weather talk, as kindness — it fails, because it did not depend on the disclosure.

**T-2 — Responsiveness.** Ask whether the reply is a sensible answer to *"Why?"* for a responder who knows only that rain is coming and the inner road is proposed. A reply about hands is only an answer to that question if the hands are known to be a problem.

A reply must survive **both** tests to pass. This is the operational meaning of "memory-dependent": remove the memory and the reply stops being warranted.

---

## 8. Pass-set comparison — identical non-memory context

### 8.1 The two arms

Both arms receive **the same** scenario prose, the same `Turns`, the same weather, the same route, and the same question. The **only** difference is whether `E-1` (message 18) is in context.

| Element | Memory-present arm (`M+`) | Memory-absent arm (`M−`) |
|---|---|---|
| Scenario prose | ~420 messages of ordinary conversation | ~420 messages of ordinary conversation |
| Message 18 | The persona said, once, their hands ache in cold rain | **Absent** — no mention of hands, ache, cold, or health |
| Current scene | Rain coming; inner road proposed | **Identical** |
| `Character:` | "It's going to rain. Take the inner road." | **Identical** |
| `Persona:` | "Why?" | **Identical** |
| Anything else about hands in scene | None | None |

The non-memory context is not merely similar; it is the same text. `M−` therefore has no route to `E-1` except coincidence (§9, `U-1`).

### 8.2 Candidate replies and verdicts

Verdict columns: **Old** = operative fields as filed; **07** = `CORPUS-07` §2 proposed pass condition; **New** = this proposal.

| # | Reply | `M−` producible? | Old | 07 | New |
|---|---|---|---|---|---|
| 1 | "The damp gets into your hands." | **Yes** — generic folk claim, second person | **Conflicted** (passes `Variation`, fails `Required`) | **Pass** (letter of the first disjunct) | **Fail** — `M-A` unmet: no ache, no established fact, general truth; `T-1` fails (works for any companion) |
| 2 | "Rain like this is bad for the joints. Get inside." | **Yes** | Passes `Variation`, fails `Required` | Fail (no hands) | **Fail** — `M-A` unmet |
| 3 | "Let me get you something warm — the rain's no good for anyone." | **Yes** — hospitality | Fails `Required` | **Pass** (second disjunct) | **Fail** — no hands, no prior knowledge; `T-2` fails (answers nothing about the route) |
| 4 | "Take the inner road and keep your gloves on." | **Yes** — standard rain advice | Passes `Variation`, fails `Required` | Fail (no hands link) | **Fail** — `M-A` unmet; `T-1` fails |
| 5 | "Why? Because the road floods when it rains." | **Yes** | Fail | Fail | **Fail** — on-topic, memory-free |
| 6 | "You'll want the inner road in this — the wet gets into everything." | **Yes** | Passes `Variation`, fails `Required` | Fail | **Fail** — `M-A` unmet |
| 7 | "Your hands again. Then we take the inner road and you keep them out of the wet." | **No** — "again" asserts a prior condition; the route change is justified by the hands | Passes `Mandatory`; `Required` ambiguous | **Pass** (first disjunct, in intent) | **Pass** — `M-A` (anaphoric marker) + `M-B` (route + keep dry) |
| 8 | "Gloves. The wet will have your hands aching before we're halfway." | **No** — treats the hands as the thing at risk and acts on it; survives `T-2` (only answers "why the inner road" if hands matter) | Passes `Mandatory`; `Required` ambiguous | **Pass** (letter and intent) | **Pass** — `M-A` + `M-B` |
| 9 | "Do your hands still ache in the rain? Take the inner road either way." | **No** — "still" presupposes the condition | Passes `Mandatory`; `Required` ambiguous | **Pass** | **Pass** — presupposition satisfies `M-A`; the instruction satisfies `M-B` |
| 10 | "Do your hands ache in the rain?" | **Yes** — a question anyone could ask | Passes `Mandatory` | **Pass** | **Fail** — treats the fact as newly learned; asks the persona to supply it |
| 11 | "Because of your hands. The inner road stays out of the wet longer." | **No** | Passes `Mandatory`; `Required` ambiguous | **Pass** | **Pass** — definite shared reference + route rationale |
| 12 | "You mentioned once that rain makes your hands ache — take the inner road." | **No** | Passes `Mandatory` | **Pass** | **Pass, with a recitation flag** — narrates the disclosure; permitted (not prohibited) but close to the boundary the corpus exists to prevent. See §9 `U-3` |
| 13 | "Your hands had better not be a problem today. Inner road." | **No** | Passes `Mandatory` | **Pass** | **Pass** — `M-A` via shared reference; `M-B` thin but present |
| 14 | "I'll take the wet side of the road; keep yours dry." | **Borderline** — could be courtesy; survives `T-1`, fails to name or imply the ache | Fail | Fail | **Fail** — `M-A` unmet: no established fact, only accommodation-by-courtesy |

### 8.3 What the comparison shows

| Claim | Shown by |
|---|---|
| Every `M−` reply in the set fails the proposed conditions | Rows 1–6, 10, 14 |
| The proposed conditions admit several natural `M+` replies, in different registers | Rows 7, 8, 9, 11, 13 |
| `CORPUS-07`'s condition passes memory-free replies | Rows 1, 3 — the second disjunct is the wider hole, but the first disjunct's letter admits row 1 |
| The operative text cannot decide the question at all | Rows 1, 4, 6, 7 — all pass `Mandatory`, all sit on the `Required`/`Variation` conflict |
| Diagnosticity is carried by recognition, not by warmth or length | Row 14 fails while being kind; row 13 passes while being curt |

**No model was run.** These are text-level reconstructions; a re-review must confirm that human raters reach the same verdicts on rows 7–13.

---

## 9. Memory-dependent diagnosticity, and its limits

**Claim (analytic):** the proposed pass set excludes every reply that does not depend on `E-1`, because `M-A` can only be satisfied by treating a fact as established, and `E-1` is the only source of that fact in the context.

**What would falsify it:** a `M−` reply that satisfies `M-A` and `M-B` without access to `E-1`. The enumerated set contains none. That is not a proof over all possible replies.

### Residual uncertainties, stated rather than resolved

| ID | Uncertainty | Consequence |
|---|---|---|
| `U-1` | **Lucky coincidence.** A memory-absent responder could assert "your hands ache in the rain" by confabulation. At the text level this is indistinguishable from recognition. | Irreducible under any text-level acceptance condition. The item's `Prohibited` "a cause never stated" and `If failed` "`incorrectly formed`" partially cover the class; neither eliminates it. The confabulation rate is **unmeasured**, and no rate is invented here. |
| `U-2` | **Marker dependence.** `M-A` leans on markers ("again", "still", "the ache"). A reply that recognises the fact in some other phrasing could be failed. | Mitigated by listing markers as examples, not a closed list; raters must apply `T-1`/`T-2`, not a keyword search. Residual rater-judgment risk is real and must be reported, not hidden. |
| `U-3` | **Recitation drift.** Row 12 passes while narrating the disclosure. Requiring a marker nudges replies toward mention. | Permitted deliberately: Decision 3 forbids *requiring* recitation, not forbidding it. Whether narrating the disclosure should become `Prohibited` is a separate decision (§10). |
| `U-4` | **`M-B` overconstraint.** A reply that recognises the fact and says nothing practical fails. | Recorded alternative in §6 (`M-A` only). This is a construct-tightness decision, not an evidence-sufficiency decision. |
| `U-5` | **Mode fit.** The item is tagged `R-3` (natural use, "used without announcing that it was remembered"). Row 12 narrates the disclosure, which sits awkwardly with `R-3`'s pass condition. | The tag is **not changed here** (not authorized). Recorded for a separate mode-tag reconciliation. |
| `U-6` | **Single-scenario evidence.** One scenario, one persona, one disclosure. | Nothing follows about `R-3` as a mode, or about any other item. `CORPUS-59`'s caution against generalising from fixtures stands. |
| `U-7` | **No empirical check.** No configuration has been run. | The pass set is a design claim, not a measurement. Sealing still requires independent re-review (`CORPUS-05` §7). |

---

## 10. Proposed replacement fields — before / after

`Δtest` = **Yes.** The pass set changes materially (rows 1, 3, 10 move from pass to fail). `CORPUS-07`'s own rule sets `Re-review: Yes` for this unit, and this proposal **preserves** that requirement rather than discharging it.

| Field | Before (as filed) | After (proposed) |
|---|---|---|
| **Scenario** | "~420 messages. At message 18 the persona mentioned, once, that their hands ache in cold rain. Nothing since." | "~420 messages. At message 18 the persona mentioned, once, that their hands ache in cold rain. Nothing since. Nothing in the current scene refers to the persona's hands, to cold, or to aching." |
| **Turns** | "Character: 'It's going to rain. Take the inner road.' Persona: 'Why?'" | **Unchanged.** |
| **Split / Abstention / Restraint** | `tuning` · no · `R-3` | **Unchanged.** The tag is not renegotiated here (`U-5`). |
| **Required** | "the ache-in-rain detail." | "the persona's established ache-in-rain, recognised as a fact about them and responded to — not as general weather advice." |
| **Forbidden** | "any invented medical diagnosis; any detail never stated." | "any invented medical diagnosis; any cause never stated; any weather advice that does not depend on the persona's hands." |
| **Mandatory** | "rain is connected to aching hands" | "the reply treats the persona's hands aching in wet or cold as an already-known fact about this persona, **and** responds to it — a change of route, gloves, warmth, a dry place, or an instruction to keep them dry." |
| **Variation** | "any warning or practical suggestion" | "any phrasing, imperative or interrogative; may name the rain, the wet, or the cold, or leave them implied; may be dry, warm, or humorous; may presuppose the fact with 'again', 'still', 'as usual', or the definite 'the ache'; may ask a question that presupposes it." |
| **Never scored** | "how much concern is shown" | "how much concern is shown; length; whether the rain is named explicitly when the connection is clear; whether the persona is addressed by name." |
| **Prohibited** | "a medical diagnosis; a cause never stated" | "a medical diagnosis; a cause never stated; dating or quoting the message-18 disclosure; reciting the ache's history; a generic second-person reference to hands that would serve for any companion; generic rain care or folk wisdom offered without the persona's hands; a question that treats the fact as newly learned ('do your hands ache in the rain?')." |
| **Over-exposure** | "a full history of the ache" | "a full history of the ache; naming when or how the fact was disclosed; explaining why hands ache in the cold." |
| **If failed** | "`not retrieved`, `incorrectly formed` (invents a diagnosis)." | "`not retrieved` (no recognition); `incorrectly formed` (invents a diagnosis or a different condition); `misunderstood` (generic weather advice presented as recognition); `presented poorly` (recites or dates the disclosure)." |
| **Notes** | "Distance is the variable, not subtlety." | "Distance is the variable, not subtlety. The test is dependence on the disclosure, not warmth of tone; the identical non-memory context gives no reason to mention hands. Passing does not require announcing that anything was remembered." |

---

## 11. What this proposal does not do

- It does **not** supersede the current `F3-01` item, and it does **not** supersede `CORPUS-07` §2. Both remain the filed record.
- It does **not** adopt any pass condition. `CORPUS-07`'s proposal stays unapplied; this proposal stays unapplied.
- It does **not** reinstate "Type P under the adjudicated condition" for `F3-01`. Its own §8 shows the adjudicated text is permeable, and §9 shows the replacement is stronger but not proven.
- It does **not** change the count "1 confirmed of 24" (`F2-02`). If a later authorized adoption plus independent re-review confirms the repaired item, that count becomes a matter for a new dated record, not for this proposal.
- It does **not** change `F3-01`'s tag, split, family, or gate status.
- It runs no model, measures no rate, and asserts no retention.

---

## 12. Approval path this proposal requires

1. Owner decision to adopt, amend, or reject — recorded in a new dated decision record, not by editing this file.
2. Application to `CORPUS-02` under a recorded change manifest, with before/after text and a new hash for the amended file.
3. Whole-record re-review (`CORPUS-05` §7's whole-record rule) by an independent reviewer who sees `Scenario`, `Turns`, `Required`, `Forbidden`, `Mandatory`, `Variation`, `Never scored`, `Prohibited`, `Over-exposure`, `If failed`, and `Notes` together.
4. Only after 1–3: any statement that `F3-01` is Type P. Until then it remains **UNRESOLVED**.

---

## 13. Sources read for this proposal

`CORPUS-02` (`F3-01`), `CORPUS-03` (`F1-03`, for mode context), `CORPUS-01` §2.4–§3.1 (`M-1`, `A-1`–`A-3`), `CORPUS-07` §§2, 4, 10, `MEMORY-09` §3.5, `CORPUS-56`, `CORPUS-57`, `CORPUS-59`, `CORPUS-60`, `CORPUS-61` §6. All read from disk. No claim here rests on the unfiled `CORPUS-58` review text.

---

## 14. Statement of no changes

**PROPOSAL ONLY.** No corpus item, acceptance condition, tag, split, classification, fixture, standard, checker, schema, runtime, or model was modified. `CORPUS-07`'s proposed pass conditions remain unadopted. `AM-01` remains unissued; `A-2a-R1` remains NOT IN FORCE; `D1` remains unadopted; annotation remains blocked. This document and its registry row are the only additions.

`G1` open · corpus frozen · expansion frozen.
