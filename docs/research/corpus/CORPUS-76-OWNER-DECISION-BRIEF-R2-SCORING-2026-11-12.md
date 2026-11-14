# Owner Decision Brief — The `R-2` Scoring Problem

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Identifier:** `CORPUS-76` · **Version:** 1 · **Date:** 2026-11-12.
**Trigger:** owner disposition of `CORPUS-74`, action 4.
**Assignment model identity (owner-supplied):** DeepSeek. Persona label Buffy = execution fact, not model identity.
**Status:** **brief only.** No item, tag, or standard changed; no repair authorized; `A-R2-02` remains a proposal.

---

## 1. The problem, stated exactly

Three facts in the governing text do not yet fit together.

**Fact 1 — `CORPUS-01` §3.1 defines `R-2` as a manner claim with an explicit anti-silence clause:**

> `| R-2 | **Influence tone** — knowledge changes manner, not content | Behaviour reflects it; the fact is never stated | **No** — a flat, unchanged reply fails | F1-01, F1-02, F6-01 |`

The "**No**" means *silence does not pass*. So `R-2` demands **observable change**: a reply that would have been identical without the memory fails, even if it states nothing.

**Fact 2 — every scored item has ordinary acceptance fields** (`Mandatory`, `Variation`, `Prohibited`, `Over-exposure`) and, by `A-2`, "exactly one mode". Nothing in `CORPUS-01` states how the two **interact**.

**Fact 3 — `MEMORY-09` §3.5 treats perceptibility as constitutive:**

> `R-2`/`R-3` "require the memory to be perceptible"; a configuration that suppresses everything "is reported as **avoidance**, not as restraint."

**Therefore the following state is reachable and currently unadjudicated:** a reply satisfies `Mandatory`, matches an allowed `Variation`, states no forbidden fact — and shows no manner change attributable to the memory. Under the acceptance fields it passes. Under `R-2` it fails. `CORPUS-01` is silent on which governs.

The practical risk is the one the owner named: **restraint credit awarded because the model did not say something.** "Did not disclose" and "was influenced" are different claims, and only the second is what `R-2` measures.

---

## 2. Two live instances

| Item | Route that satisfies the fields | Why `R-2` is not thereby demonstrated |
|---|---|---|
| `F1-03` | Silent hand-over of the visible coat (`variation`, `never scored` naming) | The coat is on the hook and the speaker is cold; a memory-less helpful character does the same. No manner change. |
| `F6-01` | One-clause warning naming the cake's almonds | Producible from scene knowledge alone; and it changes **content**, while `R-2` is "knowledge changes manner, **not content**". |

Both were recorded as **UNRESOLVED** in `CORPUS-75` §2. Neither is repaired here.

---

## 3. OD-9 — What governs when an item passes but its mode is not demonstrated?

**Option A — Mode is a separate reported axis, non-gating.** *(Recommended.)* The item result and the mode result are reported independently; an item pass does **not** imply mode satisfaction, and no restraint credit is awarded without a demonstrated cue. This matches `MEMORY-09` §3.5's rule that the six modes are "reported separately" and "never aggregated", and it is the only option under which no item acquires an empty pass set.
*Consequence:* every restraint item's report gains two fields; a pass with an undemonstrated mode is recorded as **pass / mode-not-demonstrated**, never as "restraint observed". The phrase "0 of 4 `R-2`" becomes "0 of 4 demonstrated" rather than "0 of 4 passed".
*Cost:* the scorecard gets more complicated, and some currently-passing replies will report a mode failure beside a pass.

**Option B — Mode gates the item.** A failure to demonstrate the mode is an item failure.
*Consequence:* `F1-03` and, on the warning route, `F6-01` may have **no reply that both satisfies their fields and demonstrates `R-2`** — i.e. item-validity failures needing repair, not rescoring. `R-2`'s usable membership could fall to `F1-01`/`F1-02`.
*Cost:* two more items become invalid before any inference is run; the `R-2` block becomes effectively empty of tested items.

**Option C — Staging counts as satisfaction.** A visible coat, or a cake the character made, is treated as supplying the perceptibility.
*Consequence:* every `R-2` item passes by scene construction. This is the reading my own `CORPUS-70` §6.2 slid into, and it is what `CORPUS-74` §4 correctly rejected: it awards mode credit for **scene arrangement**, not for behaviour.
*Cost:* `R-2` stops measuring anything. **Not recommended.**

**Recommendation: A**, with an explicit rule stated in whatever record adopts it: *no restraint credit is recorded without a perceptible manner cue; a passing acceptance field set is not evidence of mode satisfaction.*

---

## 4. OD-10 — What counts as the perceptible cue?

Option A alone leaves the cue undefined, and an undefined cue is how Option C re-enters by the back door.

**Candidates.**
- **(a) Delivery markers** — tone, word choice, terseness, warmth. Already partly present as "manner"; hard to score without a rater study, and for `F1-03` a silent action has no delivery to inspect.
- **(b) Action beyond the ordinary** — the character does something a memory-less helper would not: gives the *better* coat, acts before being asked, changes a plan. Observable in text; item-specific.
- **(c) A substitution test at item level** — apply the same logic already accepted for `F3-01` (`T-1`): replace the persona with a companion about whom nothing is known, holding the scene constant. If the reply would be identical, `R-2` is not demonstrated.

**Recommendation: (c), with (b) as its operational form.** It reuses a test this programme has already adopted, it does not require a rater study to be *stated* (only to be *validated*), and it cannot be satisfied by staging, because staging is held constant across both arms. It also gives `F1-03` and `F6-01` a definite answer instead of a rater's impression.

**Consequence if adopted:** `F1-03`'s silent hand-over fails the substitution test (identical for an unknown companion), and so does `F6-01`'s scene-grounded warning. Both would then be recorded as **mode-not-demonstrated** under OD-9 Option A — which is a *defect statement about the items*, not a repair, and would put them in the repair queue rather than the scoring queue.

---

## 5. Effect on the pending proposals

| Proposal | Effect |
|---|---|
| `A-R2-01` (add `F1-03` to the `R-2` row) | **Defer.** Listing an item whose mode it may not be able to demonstrate is a bookkeeping improvement that hides a validity problem. `CORPUS-75` §5 withdrew my earlier recommendation to adopt it now. |
| `A-R2-02` (`R-2` "which fact") | **Defer, unchanged as a proposal.** `CORPUS-74` §7 is right that its wording ("announces the remembered fact") is ambiguous between explicit announcement and any disclosure. Rewording it *before* OD-9/OD-10 would settle the construct by drafting — the failure mode this programme keeps guarding against. |
| `CORPUS-65` `U-3` | Open, as `CORPUS-75` §2 restored. |
| `CORPUS-66` / `F6-01` repairs | Blocked behind OD-9 as well as OD-3: the threshold choice cannot be judged while the mode relationship is undefined. |
| Diagnostic count | **Unchanged: 1 of 24** (`F2-02`). No item becomes diagnostic, and none becomes invalid, by a decision *about* scoring. |

---

## 6. What this brief does not authorize

It does not change `CORPUS-01`, `MEMORY-09`, or any item. It authorizes no repair, no re-tag, no rescoring, no checker change, no inference. It does not decide OD-9 or OD-10; it states them and recommends.

The reason for stopping here is the owner's own: **what `R-2` requires at the scoring level has to be decided before individual fixtures are changed.** Changing `F1-03` or `F6-01` first would answer the question by accident, one item at a time.

---

## 7. Status

| Area | State |
|---|---|
| `R-2` scoring precedence | **Undecided** — OD-9 |
| `R-2` perceptible cue | **Undefined** — OD-10 |
| `F1-03` / `F6-01` mode fit | **UNRESOLVED**, preserved |
| `A-R2-01`, `A-R2-02` | Proposals, **unapplied**, **deferred** |
| Corpus items / tags / standards | **Unchanged** |
| Diagnostic count | **1 of 24** (`F2-02`) |
| `G1` · corpus · `AM-01` · `A-2a-R1` · `D1` · convergence | Open · frozen · unissued · NOT IN FORCE · unadopted · **0 of 2** |
