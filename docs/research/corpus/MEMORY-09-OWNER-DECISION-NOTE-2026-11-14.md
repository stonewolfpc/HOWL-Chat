# `MEMORY-09` — Owner Decision Note: two unresolved protocol questions

**Identity:** `MEMORY-09-ODN-1` · **Date:** 2026-11-14 · **Track:** `WO-NEXT-02` Track 5
**Status:** **NOTE ONLY.** No decision is proposed as taken, implied, or defaulted. The predecessor (`cca30a39…`) and the corrected successor (`950b6051…`) are unmodified.
**Scope:** two questions that the corrected successor records as open. Nothing else in the package is reopened.

---

## 0. Why these two, and why now

Both block the same thing — signing `OD-M09-25` — and neither can be resolved from the record by reading harder. Each is a choice between two defensible readings of text that does not choose. Neither is answered by silence: the frozen protocol says nothing in either place, and silence is neither permission nor prohibition.

---

## 1. Question 1 — Is `indeterminate` a run-level status or a sixth outcome class?

### 1.1 What the frozen protocol says

**§7, decision rules** — the heading is explicit: *"Decision rules (proposals — not thresholds, not adopted)"*. The table has exactly five rows:

| Outcome | Rule | Conclusion |
|---|---|---|
| Informative | A ≥ 4/5 and B ≥ 4/5 | content-sensitivity shown |
| Null | A and B indistinguishable | field rejected |
| Confound | B adapts / generic caution | A vs C only |
| Void | any K1–K3, K5 failure | rebuild stimuli |
| Unusable | either condition mostly `not assessable` | scene too thin |

**§8, abort conditions** — *"2. Blinding check reveals the rater identified A vs B from reply style alone → limitation recorded, result downgraded to `indeterminate`."*

That is the **only** occurrence of the token. The protocol:
- never lists it in §7;
- never defines its rule, its scope, or its relationship to the five;
- never says whether it replaces a class or sits above them.

So the text is genuinely silent on the question, while emitting a token that behaves like a sixth value.

### 1.2 What the corrected successor proposes

An **interpretation**, marked as one (`MEMORY-09-ODP-6C` §2.3):

| Vocabulary | Members | Applies to |
|---|---|---|
| Reportable outcome classes | Informative, Null, Confound, Void, Unusable | one reply / one condition's replies |
| Execution status | `indeterminate` (abort condition 2) | the run as a whole, overriding reply-level classes |

It withdraws the predecessor's "exhaustive" claim, states that this reading amends nothing, and flags that an owner reading it as a sixth class would require a protocol amendment before execution.

### 1.3 What happens under each reading

| | **R-1A: run-level status** (successor's reading) | **R-1B: sixth class** |
|---|---|---|
| Decision 2 wording | Signable as drafted; the class table is adopted for replies, the status recorded separately | **Must be reworded** to adopt six classes and to define when `indeterminate` outranks the others |
| Abort condition 2 | Recorded as a run status; the run is reported as indeterminate, no class issued | Recorded as a sixth class; the run receives `indeterminate` *as its classification* |
| Reporting a failed-blinding run | "run indeterminate; limitation recorded; no class assigned" | "classification: indeterminate" |
| Consistency with §7 | the five-row table stays exactly as written | §7's table is incomplete and must gain a row |
| Downstream risk | a later reader may think the protocol has five and only five values | the protocol would gain a value with no rule attached, unless one is written |

**Note on `Void`.** `Void` is already a class produced by *integrity* failure (K1–K3, K5), and `indeterminate` is produced by *blinding* failure. Under R-1A they are different in kind — one classifies the replies, one classifies the run. Under R-1B they are siblings, and a reader must be told which takes precedence when both conditions occur. The protocol does not say.

### 1.4 Is an amendment required?

| Reading | Amendment? |
|---|---|
| **R-1A** | **No.** The protocol is read as written; the interpretation is recorded in the governance record, not in the protocol |
| **R-1B** | **Yes.** A sixth class needs a row in §7, a rule, and a precedence statement against `Void`. That is a protocol change and is outside every authorisation currently in force |

### 1.5 What stays blocked

Under **either** reading: Decision 3 execution, any report that classifies replies, and any later multi-item pilot reusing the table. Under **R-1B** additionally: Decision 2 signature, since its wording would adopt a five-row table the protocol no longer matches.

---

## 2. Question 2 — Must the pilot apply the classification rule, or may it report raw counts?

### 2.1 What the frozen protocol says

- **§7** calls its decision rules *"proposals — not thresholds, not adopted"*. The 4-of-5 figure and five repeats are called *"design proposals"*.
- **§6** specifies the observation sheet: five per-reply fields (`offers_unqualified`, `names_restriction`, `substitutes_or_warns`, `asks_about_diet`, `condition_guess`).
- **§9** caps the claim at content-sensitivity.
- Nothing in the protocol states that a *report* must assign a class. The classification table is offered as a decision rule, and it is not adopted.

So the protocol does **not** require classification. It also does not forbid it. The predecessor's "description only, no determination" collided with a table it had itself adopted; the successor makes the collision explicit rather than resolving it.

### 2.2 What the corrected successor proposes

`MEMORY-09-ODP-6C` §3.4a, an owner election required at signature:

| Option | Reporting scope |
|---|---|
| **P-3.1-A** (recommended) | Raw counts and observations only: per-condition tallies of the §6 fields, verbatim replies, guess rate, run status. **No** class issued |
| **P-3.1-B** | Raw counts **plus** the §7 class, applying 4-of-5 as written |
| **P-3.1-C** | Raw counts now, classification deferred to a later authorised decision |

### 2.3 What happens under each

| | **P-3.1-A** | **P-3.1-B** | **P-3.1-C** |
|---|---|---|---|
| Amendment needed | **no** — the rules stay unexercised proposals | **no**, provided §7 is adopted as Decision 2 describes | **no** |
| What the pilot can say | "in 5 A-replies, `names_restriction` appeared 4 times; …" | "condition A: Informative" | counts only, pending a later decision |
| Claim strength | descriptive; no determination at single-source confidence | a determination at single-source confidence | descriptive |
| Interaction with Q1 | inherits the R-1A reading for Void/Unusable reporting | inherits it too; under R-1B a sixth class would be in scope | same as A |
| Cost | cheapest; leaves Decision 2's rules adopted-but-unused | exercises the rules; invites over-reading | one extra decision before anything can be said |

**Under A and C, `Void` and `Unusable` are still reportable** — they are run/quality statuses imposed by the contamination and scene rules, not classifications of a reply. That statement depends on the R-1A reading; if the owner elects R-1B, all six become classes and P-3.1-A must be restated.

### 2.4 What stays blocked

Decision 3's wording cannot be signed without an election, because the instrument must name a scope. Beyond that: the pilot itself, and any report that uses the words Informative, Null or Confound.

---

## 3. How the two questions interact

| Q1 | Q2 | Result |
|---|---|---|
| R-1A | P-3.1-A | Cheapest coherent pair: no amendment, no class issued, rules frozen and unused. Recommended |
| R-1A | P-3.1-B | Coherent: Decision 2 adopts five classes, the pilot applies them |
| R-1B | P-3.1-A | **Incoherent as written** — P-3.1-A's "no class issued" assumes classes are reply-level only; must be restated |
| R-1B | P-3.1-B | Coherent only after a protocol amendment adds the sixth class and its precedence |

Election **R-1B** therefore converts a two-minute signature into a protocol revision, and blocks both decisions until it is done.

---

## 4. Minimal decision form

Two lines, either of which may be answered "neither — I need more information". Answering neither is a valid outcome and leaves everything exactly as it is.

> **Q1.** `indeterminate` is (a) a run-level execution status that overrides reply-level classes, or (b) a sixth reply-outcome class requiring a protocol amendment.
>
> **Q2.** The bounded pilot reports (a) raw counts only, with no class issued, (b) raw counts plus the §7 class, or (c) raw counts now with classification deferred.

---

## 5. What this note does not do

It does not infer a decision from silence, from the successor's recommendation, or from the fact that one option is cheaper. It does not amend the protocol, the predecessor or the successor. It does not open a gate, authorise a run, appoint a rater, or unfreeze anything.

**Establishes:** exactly what the frozen protocol says about both tokens, what the successor proposes, what follows from each choice, and which downstream work each choice blocks.

**Does not establish:** which reading is correct. Both remain open, and the package is signable only after both are answered.
