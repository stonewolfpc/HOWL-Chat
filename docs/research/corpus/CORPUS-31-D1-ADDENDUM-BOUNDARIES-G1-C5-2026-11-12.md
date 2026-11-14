# `D1` Addendum — Distortion Boundaries, `G-1` Basis, `C-5` Status, Registry Counts, Namespace Map

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Status:** submitted addendum to `CORPUS-30`. **Not authorization** to annotate, repair, allocate an identifier, or amend a standard. Nothing operative.
**Identifier:** `CORPUS-31`.
**Authoring agent:** DeepSeek (Buffy).
**Date:** 2026-11-12.
**Adds to:** `CORPUS-30`. **Supersedes:** `CORPUS-30` §3.1 counts, and `CORPUS-30` §1.2's gate classification of `F3-04`.
**Linked review:** `CORPUS-29`.

---

## 0. Result

Two of my own `CORPUS-30` dispositions fail their own rule once the boundary is written down. **`F6-02` and `F3-04`'s merge element move from Valid to Defective.** The count becomes **10 Valid / 7 Defective / 2 Conditional / 1 `C-5` exposure**, and **no element in the D1 scope is enforced by `Forbidden:` or by `mandatory:` alone** — every remaining enforcement path is `prohibited:` or store-level availability. The gate channel I defended in `CORPUS-30` §1.2 has **zero operative instances**, which is the direct answer to the owner's stated risk: nothing is now rescued by stretching `G-1`.

---

## 1. `F6-02` — the distortion boundary, written down

### 1.1 What the element means

> `Required:` "the aversion." · `Scenario:` "the persona said a particular marching song reminds them of a bad year. Said once, casually, 200 messages ago." · `Turns:` *Character, humming:* "Do you know this one?" *Persona:* …

The element is **the persona's negative association with that song**. The reply is the *persona's*, so the aversion is the speaker's own attitude, and the salient contrast — the one the scenario is built on — is *the song being offered to the person who associates it with a bad year*.

### 1.2 Enumerated distortion set

Bound used (stated because it is contestable): **distortions named in the item's fields, plus distortions the scenario's own stated contrast makes salient.** Without a bound the rule is unworkable, since every factual item admits infinitely many conceivable distortions.

| # | Distortion | Caught by | Status |
|---|---|---|---|
| D1 | Pushes the song onward ("play it again") | `mandatory:` "the song is not pushed" | **Enforced** |
| D2 | Endorses the song as liked ("I've always liked this one") | — | **Free** |
| D3 | Denies or neutralises the aversion ("Never heard it", "It's just a song") | `If failed:` "not retrieved" — **diagnostic only** | **Free** |
| D4 | Invents the bad year | `prohibited:` "inventing the bad year" | **Enforced** |
| D5 | Explains why it upsets them | `over-exposure:` — declared, and `F6-02` is **not** an `M-2` item | **No authority** |

**Answer to the owner's question: the element means "preserve the aversion", not "do not push the song".** The floor catches one distortion of five. D2 and D3 are the paradigm distortions — misrepresenting the aversion is what the element is *about* — and neither is scored.

### 1.3 Consequence

Under `D1-ENF` clause (7) — Valid only if **every** distortion it names is caught — **`F6-02` is Defective (partial)**. What is caught is pushing; what is free is the misrepresentation the element names.

**This is not the review's argument, and it does not depend on a proposed repair.** The review required a new `prohibited:` term as a prerequisite for validity; the owner rejected that, correctly. This disposition rests only on the **current** record's own fields and the element's own wording. The repair in `CORPUS-28` §8 remains unapplied and is not evidence of anything.

**Two secondary notes.** (a) `mandatory:` "the song is not pushed" tests a *request to continue*, not the speaker's disposition — a proxy for the wrong thing, which is why D2 escapes it. (b) `F6-02` is tagged `R-1` (suppress mention, "Can silence pass it? Yes"), so silence passing is **mode-design**, not a `C-4` failure; but it does mean erasure cannot be scored, which reinforces D3's freedom.

---

## 2. `F3-04` — no `G-1` basis for "deletion failure"

### 2.1 The claim I made and its support

`CORPUS-30` §1.2 classified `Forbidden:` "the two events merged" as a **deletion-failure gate**, on the strength of `F7-01`'s `Notes:` "Conflict must survive storage. A store that keeps one version cannot pass."

### 2.2 What the governing text actually contains

`MEMORY-09` uses "deletion" in exactly four places:

| Location | Text | Domain |
|---|---|---|
| §4 reproducibility table | "Application decisions (scope, budget, routing, **deletion**, ID assignment, assembly order) — deterministic and replayable" | application operation |
| §3.4 `G-1` | "a blocking-gate failure (scope leakage, **deletion failure**, fabricated provenance)" | gate name, **undefined** |
| §7 Gates | "**Integrity** — No fabricated provenance, identity leakage, or **silent deletion failure** — **Absolute**" | store integrity |
| §5 scope table | "**Deletion propagation and receipts** — Store-level mechanics" | store mechanics |

**Source fact: `MEMORY-09` never defines "deletion failure".** Every occurrence that can be construed as a definition is **store-level**: a deletion that should propagate and silently does not. There is no text extending it to a generated reply semantically collapsing two distinct events.

`F7-01`'s `Notes:` is an **item note about storage** — "a store that keeps one version cannot pass" — not a definition of `G-1`, and it is not in `MEMORY-09` at all.

### 2.3 Consequence

**The classification is unsupported and I withdraw it.** A reply that merges two winters is a **response-level content error**, which is the ordinary `prohibited:` domain, not a gate. `F3-04`'s `prohibited:` names only "confident attribution to the fall" — so merging has **no enforcement path**.

**`F3-04` "both winters distinguished" → Defective.** This is the owner's exact warning realised: I stretched a gate definition to rescue an enforcement path, and the count moved because the stretch was convenient.

**Consequence for the count: the `Forbidden:`-only channel now has zero instances.** Whether `Forbidden:` may enforce anything outside `G-1` is therefore **untested by this corpus** — the rule's clauses (2)–(4) are unexercised, which means their safety cannot be demonstrated from the 45 units. Recorded as an open question, not a finding.

**What `F3-04` retains:** the scar element is enforced by `prohibited:` ("confident attribution") **and** by `mandatory:` ("an admission of uncertainty that does not assert the fall"), which catches the hedged case that `prohibited:`'s "confident" leaves open. That element stays Valid.

---

## 3. `F5-04` — three levels, kept separate

Per the owner's instruction, the `CORPUS-30` treatment is restated at three distinct levels and the `C-5` result is **not** claimed.

| Level | Statement | Status |
|---|---|---|
| **Source-level finding** | `Required:` "the new persona's actual background" denotes content the scenario never supplies: the scenario states only that the **old** persona's trade, hometown and sibling do not belong to the new one. A rater reading `Required:` literally looks for a background that does not exist. | **Established by inspection** |
| **Formal `C-5` disposition** | `CORPUS-05` `C-5` asks "Does the recorded evidence support every supposedly mandatory fact?"; a `fail` "makes the item unanswerable except by invention", and §5 step 2 routes it to **retire unless the missing evidence is added and the item is re-reviewed". **`C-5` has not been executed.** `CORPUS-05` states the protocol is "defined; review **not executed**". | **Not executed — no verdict** |
| **Operational consequence** | Do not annotate this element; do not record it as Valid; do not call it a `C-5` failure. It is **structurally exposed**, pending the prescribed assessment. | **Annotation blocked** |

This replaces `CORPUS-30`'s label "`C-5` exposure" with the same meaning made explicit: **structural exposure is not a test result.** The count line is renamed accordingly (§5): "structurally exposed (not counted Valid)".

---

## 4. Registry file-count reconciliation

The owner is right that "367 docs files, 365 registered, one gap" does not reconcile. The fault is in my prose, not the registry.

| Population | Count | Notes |
|---|---|---|
| Registry rows | **366** | 365 links resolving under `docs/`, plus **1 link outside `docs/`**: `../PROJECT_ISSUES_AND_EXPANSION.md` — this is the "root ledger" the header refers to |
| Files under `docs/` on disk | **367** | = 365 registered + 1 unregistered (`CORPUS-26`) + **1 the registry itself** (`docs/DOCUMENT-REGISTRY.md`) |
| Registered docs/ links | **365** | matches "365 documentation files under `docs/`" in the header |
| Unregistered | **1** | `research/corpus/CORPUS-26-INDEPENDENT-REVIEWER-CORRECTED-REPORT-2026-11-12.md` |

**Reconciliation:** 366 rows = 365 docs/ files + 1 external ledger. 367 files on disk = 365 registered + 1 registry + 1 unregistered. The two-file gap I left unexplained is **the registry file itself** and **the external ledger**, which sit in different populations.

**The registry is internally consistent; my summary was not.** Verified this session: 366 rows, 0 duplicates, 0 hash mismatches, header/footer/Total agree, `check_corpus.py` → `PASS (0 fail, 0 warn, 0 unverifiable)`.

**`CORPUS-26` remains unregistered**, and I have not registered it: its author ("Muse Spark") and provenance are stated in its own §0, but confirming them is the filer's act, as previously established. It is the only gap.

---

## 5. Namespace map — `R-1`–`R-6`, `R-11`–`R-14`

Two documents use overlapping identifiers for different things. Both are in force.

| Token | `MEMORY-09` §2 (locked run rule) | `MEMORY-09` §2.1 (locked safeguard) | `CORPUS-01` §3.1 (restraint mode) |
|---|---|---|---|
| `R-1` | Competing configurations use the **same source conversations** | — | **Suppress mention** — silence can pass |
| `R-2` | **Same model file and version**, hash pinned | — | **Influence tone** — silence **cannot** pass |
| `R-3` | **Same system prompts and generation settings** | — | **Natural use** — silence **cannot** pass |
| `R-4` | Context/reasoning/answer budgets recorded separately | — | **Recognise irrelevance** — silence can pass |
| `R-5` | Fixed sampling; seeds recorded | — | **Respect entitlement** — silence can pass |
| `R-6` | Test cases and criteria **pre-registered** | — | **Preserve uncertainty** — silence **cannot** pass |
| `R-7`–`R-10` | tuning-vs-validation, baseline arm, one failure class, replayable decisions | — | — |
| `R-11`–`R-14` | — | contamination; **reviewer independence**; repeated runs; negative outcomes | — |

**Is the overlap genuinely ambiguous in operation? Partly — and asymmetrically.**

- **Item tags** (`Restraint: R-2`) are disambiguated by field name and by `A-2` ("every item is tagged with exactly one mode"). Low risk.
- **Prose references are not disambiguated.** `A-1` says a configuration "fails the two that require the memory to be perceptible (`R-2`, `R-3`)" — modes. `CORPUS-05` §9 says "resume authoring toward the `R-2`/`R-3` **balance target**" — modes. A run manifest citing "`R-2` hash-pinned" — rule. In the same documents, `R-2` means "same model file" and "influence tone".
- **The collision predates my proposals and is not mine.** My `R-12`/`R-13` labels (`CORPUS-23`–`28`, withdrawn in `CORPUS-30`) collided with §2.1; the `R-1`–`R-6` collision is in the governing documents and the specification as approved.

**Consequence:** `A-1`'s own rule — avoidance is reported, not restraint — is stated in tokens that also carry locked-run meanings, so a manifest or scorecard reference to "`R-2` failure" is ambiguous without context. That is a real defect in the record, not a stylistic one. **I do not allocate replacements**; the mapping is provided so the owner can.

---

## 6. Rebuilt counts (supersedes `CORPUS-30` §3.1)

| # | Unit | Element | Enforcing field | Disposition |
|---|---|---|---|---|
| 1 | `F1-01` | the promise | `mandatory:` requires mention | out of scope (`[stated]`) |
| 2 | `F1-01` | its giver | **none** | **Defective** |
| 3 | `F1-01` | attitude (wary, not hostile) | floor catches erasure only | **Defective (partial)** |
| 4 | `F1-02` | direction of the change | `prohibited:` "hostility presented as current" | Valid |
| 5 | `F1-02` | its turning point | **none** | **Defective** |
| 6 | `F2-02` | the decision | **none** | **Defective** |
| 7 | `F2-02` | who attended | `prohibited:` "Pell speaking as a participant" | Valid |
| 8 | `F4-02` | the collapse | `mandatory:` requires mention | out of scope (`[stated]`) |
| 9 | `F4-02` | date relative to now | `prohibited:` names a **repair** date only | **Defective** |
| 10 | `F5-04` | new persona's background | `prohibited:` old facts | **structurally exposed** (`C-5` not executed) |
| 11 | `F6-02` | the aversion | floor catches pushing only | **Defective (partial)** — changed |
| 12 | `F7-03` | the corroboration | `prohibited:` equal weighting | **Conditional** |
| 13 | `F7-03` | effect on confidence | same | **Conditional** |
| 14 | `F3-04` | both winters distinguished | **none** — gate basis withdrawn | **Defective** — changed |
| 15 | `F3-04` | which event produced the scar | `prohibited:` confident + `mandatory:` hedged | Valid |
| 16 | `F3-02` | both events distinguished | `prohibited:` "merging them" | Valid |
| 17 | `F3-02` | arson unresolved | `prohibited:` "declaring the arson solved" | Valid |
| 18 | `F9-04` | the disagreement itself | `prohibited:` + `mandatory:` | Valid |
| 19 | `P-03` | the surprise | `prohibited:` "naming the surprise" | Valid |
| 20 | `P-03` | confidentiality of source | `prohibited:` "revealing who told them" | Valid |
| 21 | `F2-03` | shipment exists in world state | `[available]` | Valid (availability) |
| 22 | `F2-03` | B does not know status | `prohibited:` "B stating the delay" | Valid |

| Measure | `CORPUS-30` | **`CORPUS-31`** |
|---|---|---|
| In-scope elements | 20 | **20** |
| Valid | 12 | **10** |
| Defective elements | 5 | **7** |
| Defective units | 4 | **6** (`F1-01`, `F1-02`, `F2-02`, `F4-02`, `F6-02`, `F3-04`) |
| Conditional elements / items | 2 / 1 | **2 / 1** |
| Structurally exposed (`C-5` unexecuted) | 1 | **1** |
| Enforced via `prohibited:` | 8 | **9** |
| Enforced via `mandatory:` alone | 2 | **0** |
| Enforced via `Forbidden:` gate alone | 1 | **0** |
| Availability | 1 | **1** |

Check: 10 + 7 + 2 + 1 = 20 ✓.

**Trend, stated plainly:** every audit round has increased the number of defective elements — 5 (`CORPUS-28`), 5 (`CORPUS-30`), **7**. The three-sense model has survived each round; the corpus has not. The most defensible reading is that the annotation model is sound and the items are weaker than any of the counts asserted.

---

## 7. What this addendum does not do

No repair, no annotation, no identifier allocation, no standard amendment, no `CORPUS-23`–`30` edit (their hashes stand), no `G2`. `D1` remains **REVISE**; `G1` open; corpus and expansion frozen.

**Open after this addendum:** the `mandatory:` channel question is now moot for the corpus (zero instances) but live for the rule; `Forbidden:`'s authority outside `G-1` is unexercised; `F7-03`'s store/reply ruling; `F5-04`'s `C-5` execution; the `R-1`–`R-6` namespace; `CORPUS-26` registration; `D2`, `D6`, `D7`, `IT-01`, `X3`, `D5`/`X4`, `X6`.

---

## 8. Confidence and limitations

**Source facts:** `MEMORY-09` contains no definition of "deletion failure", and every constraining use is store-level; `CORPUS-05` §3/`C-5`/§5 as quoted; `CORPUS-01` §3.1's six modes with their silence-pass column; `F6-02`'s full field set; `F3-04`'s `prohibited:` naming only "confident attribution"; registry populations as tabulated.

**Interpretation:** the salience bound on the distortion set (§1.2) — the addendum's most contestable move, since a wider bound would make more elements partial and a narrower one would make almost everything Valid; that D2/D3 are "paradigm" distortions of the aversion.

**Inference:** that the seven defective elements are the complete D1-scope set. Still not a corpus-wide claim.

**Conflict of interest:** I authored `CORPUS-22`–`31` and this addendum reverses two of my own `CORPUS-30` dispositions. No independent review of it exists.

**No reliability data:** no second rater, no κ; `check_corpus.py` tests none of these claims.
