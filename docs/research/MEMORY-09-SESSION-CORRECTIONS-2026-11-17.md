# MEMORY-09 — Session Corrections and Reasoning Trail

**Status:** record of corrections made during the architecture work. Documentation only. Creates no
authority, decides nothing, and does not modify any prior record.

| Field | Value |
|---|---|
| Created | **2026-11-16 23:47:21 CST (−0600)** = **2026-11-17 05:47:21 UTC** |
| Purpose | Preserve what was claimed, what was wrong, and what replaced it |
| Authority | None |

This file exists because the corrections matter as much as the conclusions: without them, the same
mistakes get made again by the next reader.

---

## 1. Sequence

1. r3 recovery investigation closed as unsuccessful within a documented search boundary.
2. ARCH-03 review retargeted from r3 to r32 under owner authorisation (record amended).
3. Ownership audit of `eval/mem1` performed read-only.
4. Architecture board produced; corrected twice.
5. Architecture v4 → v9 built with the owner; decisions recorded and open items tracked.

---

## 2. Corrections in the r3 recovery phase

| # | Correction | What replaced it |
|---|---|---|
| 1 | Claimed "public/private" was not contract vocabulary, based on the normative extract | **Wrong.** r32 §2 D1 uses `character_public` and `character_private` as scope names. The extract omits §2 prose entirely |
| 2 | Implied a separate r3 architecture record was likely | **Simpler reading:** r3 was a revision of the same growing file; the byte-size curve is monotonic. Recorded as *likely*, not established |
| 3 | Chat-derived figures (r1/r2 sizes, ~25,817 bytes, Q2–Q6 acceptance) were treated as evidence | **Corrected:** they came from a third-party assistant transcript, unverifiable in this workspace. Recorded as reported, not observed |
| 4 | A name-based search result was reported as package-wide absence | **Corrected:** a recursive sweep of `eval/mem1` was run; the broad claim was scoped to what was actually searched |

---

## 3. Corrections in the ownership audit

| # | Correction | What replaced it |
|---|---|---|
| 5 | Quoted `retrieve.py:136` as `owner_id != owner_id` | **Transcription error.** Source is `rec.owner_id != owner_id` (record field vs requester parameter). Conclusion unchanged, evidence corrected |
| 6 | Used a stored `audience` field to argue groups "must" be pseudo-characters | **Withdrawn.** A group can be a derived retrieval-time audience with no stored list. The stored-audience defect stands independently |
| 7 | Said "every world-scoped record behaves as if ambient" | **Imprecise.** Accurate: world-scoped records carry no requester condition, so access is ambient-*like behaviour* produced by absence of a class — not an authored class |
| 8 | Proposed conformance as a separate work class | **Precondition missing:** R-5a (who may amend the experiment plan) is unanswered, so no authorised channel exists yet |

---

## 4. The finding that survived every correction

**Static observation:** the experimental package models three different questions through one
identity field — record ownership (`rec.owner_id`), private-access identity (`rec.owner_id !=
requester`), and world/scenario activation identity (`record_is_active`). Renaming the field would
not fix it, because the coupling is functional.

**Consequence:** the contract's separation of memory ownership, applicability, and character
knowledge could not be represented. That is why the knowledge-edge work resisted reconciliation: the
edge was being asked to sit on an identity layer that had not separated those questions.

**Scope caveat:** this describes `eval/mem1` only. Production Howl was not inspected, and no claim
is made about it.

---

## 5. Corrections in the architecture phase

| # | Correction | What replaced it |
|---|---|---|
| 9 | v1 rule: knowledge = accounts ∪ *currently active* declarations | **Withdrawn.** It made knowledge depend on a declaration's current status and on when you asked. Replaced by materialised **coverage accounts** (K-R1 also forbids an edge without an event) |
| 10 | v1: belief memory "owned by the character" | **Withdrawn.** Contradicted D-01. Belief records are chat-contained, character-scoped, with the character relation held separately |
| 11 | v1: one source span per record | **Withdrawn.** Too strict — an inference has evidence without a single span stating the conclusion. Replaced by **evidence appropriate to origin** |
| 12 | v2: `believed` listed as an origin of knowledge | **Withdrawn.** It made belief count as knowledge. Belief is now a **separate relation**, and knowledge queries never consult it |
| 13 | v5 table: `character_public` described as truth-apt | **Withdrawn.** Only world and scenario records are truth-apt. `character_public` is a visibility marker. The file contradicted itself; the error is recorded rather than quietly fixed |
| 14 | v4/v5 test counts stated as 32 and 42 | **Wrong.** Enumerated totals were **38** and **52**. Undercounted twice; the corrected count for the current map is **87** |
| 15 | K7 could fail a successful promotion | **Rewritten.** Split into K7 (authority settled) and K11 (history preserved), with K12 for the dual-answer requirement |
| 16 | Concluded disclosure could not be evidence for promotion | **Too strong.** Disclosure may be a *basis*; what it cannot be is automatic authority. Evidence and authority are distinct (O-9) |
| 17 | Framed O-8a as "who may create world memories" | **Reframed.** Creation is open to any participant; recording is system-governed; **promotion authority** is the separate, gated permission |
| 18 | Did not distinguish a statement from an experience | **Added:** a character's utterance is a claim; it never establishes that they experienced what they described. Contribution kind must be declared, not inferred (O-12) |
| 19 | Listed O-8a-1 as **open** in the published ledger after the owner had approved the default-and-delegation arrangement | **Corrected in `MEMORY-09-OWNER-DECISIONS-2026-11-17.md`:** 8a-1 is CONFIRMED (world/scenario authoring authority, default the human operator, per domain). The remaining 8a-2…8a-6 stay open. The board's status line was stale when written |
| 20 | O-12 left open after confirmation, and the missing-declaration case was unspecified | **Corrected:** O-12 CONFIRMED. Kind is model-declared and recorded; a missing declaration is recorded as **`Unspecified claim`** and the system invents nothing; a declared recollection is not verified truth. Test **L9** added (map now 88) |
| 21 | Recorded O-8a-2's basis kinds as **CONFIRMED** on the strength of "I support the two-basis proposal" and "I recommend accepting" | **Corrected:** marked **PROPOSED**, pending explicit confirmation. Support and recommendation are not decisions; the owner then asked for the kinds to be checked *as stated*, which showed the label was premature. This is the same error class as correction 19, in the opposite direction. **Resolved:** the owner confirmed 8a-2 with all three amendments; kinds are now `record_cited` and `actor_declaration` |
| 22 | 8a-4 was recorded as "the authority holder may explicitly revoke or supersede", which implies the promoting actor executes it | **Qualification, not an error in `eval/mem1`.** The owner later established that execution is **human-controlled**: no actor may revoke or supersede another's promotion, or their own prior promotion, by disagreement alone (O-10b, rule 43). The 8a-4 wording is retained unedited and the qualification is recorded alongside it |
| 23 | Two discrepancies between confirmed decisions and r32 normtext were found unreconciled and are reported, not fixed | **(1)** K-R45–R47 specify transfer as a MEMORY-09 responsibility, contradicting O-13; recorded as a scope correction owned by the application. **(2)** D6 describes a "five-condition" predicate while §10.1 revision 2 has three conditions plus E1–E3 and remains **[P] proposed**. Neither document was altered |

---

## 6. Boundary failures avoided

Recorded because each was a live temptation during the work:

| Temptation | Why it was refused |
|---|---|
| Treating implementation behaviour as architecture | Code does not acquire authority by existing; D-01 governs |
| Reading a passing test as proof of semantics | No test has been executed; defined tests prove nothing |
| Promoting a recommendation into an owner decision | Only the owner decides; recommendations are labelled |
| Inferring approval from wording in a contract | Q2–Q7 remain open in r32 regardless of conversation claims |
| Silently resolving the D-01 vs character-ownership conflict | Recorded as O-1, then separated into two dimensions with the owner's confirmation |
| Reporting a document as written when it was not | Persistence required explicit authorisation, and did not happen until it was given |
| 24 | O-7 cited `R-5a` as its authority for the experiment-plan restriction | **No source document defines `R-5a`.** It appears only in four internally authored records. Re-anchored to **plan v3 §10, bullet 2**, which is the source-defined restriction. Historical occurrences of `R-5a` are preserved as historical references only, never as independent authority. The phrase **"in the contract"** was **not** written by the owner — it arose in collaborative assistant/DeepSeek work; authorship and original intent are **not** attributed to the owner. Amendment 1 therefore claims **no** prior intent for §10 |

---

## 7. Standing evidence boundaries

- **Defined ≠ executed.** All **118** tests are defined only; **0 implemented, 0 executed**.
- **Ratified ≠ implemented.** v9 is **ratified as a design baseline** (O-6b) as tier 3 under O-6a.
  Every decision in the ledger is confirmed **except** **Gate 2** (implementation — unanswered,
  not authorised) and **O-2** (parked).
- **Experimental ≠ production.** `eval/mem1` is the only implementation examined.
- **Reported ≠ observed.** Chat-derived and transcript-derived figures are labelled as reported.
- **Recorded ≠ signed.** Owner confirmations live in conversation; no persisted signature exists.
