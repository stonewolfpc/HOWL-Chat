# Batch-1 Rule Applicability and Rule-Provenance Audit

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Status:** evidence audit and methodology record. **Not authorization** to repair, amend, annotate, or change the ledger. Nothing operative.
**Identifier:** `CORPUS-40`.
**Author:** DeepSeek (Buffy) — same declared agent identity as the prior chain; see §10.
**Date:** 2026-11-12.
**Scope:** which rules actually governed batch-1 items; whether later rules were made applicable to them; the `Required:`/`M-1`/`R-1` interaction. `F6-02` is the worked example, not the only subject.

---

## 1. Purpose

`CORPUS-39` reported that `CORPUS-01` §2.4 is "required from batch 2 onward" and therefore flagged `F6-02`'s formal compliance as pending an annex pass (gap G-6). The owner's reading went further: applying §2.4 to `F6-02` would be **retroactive**.

**Both readings are partly wrong, and the record says so explicitly.** `CORPUS-01` contains an annex mechanism for exactly this case, names it a **sealing blocker**, and records it **complete**. This audit establishes the actual provenance and applicability of each relevant rule, because a repair authorised on the "retroactive" premise would rest on a false account of the standard.

---

## 2. Sources

| Source | Revision | SHA-256 (prefix) | Bytes |
|---|---|---|---|
| `CORPUS-01-SPECIFICATION-AND-INDEX-2026-11-12.md` | `G1` specification; §2.4, §2.5, §3.1, §5 roadmap, §6 status | `be1f11b46383e98a…` | 18,587 |
| `CORPUS-02-BATCH-1-ITEMS-2026-11-12.md` | 22 batch-1 items (`F6-02` line 145) | `75aa9be707cc0c44…` | 23,025 |
| `CORPUS-03-BATCH-2-ITEMS-2026-11-12.md` | 16 batch-2 items | `d0dae1c0951e9ed2…` | 18,183 |
| `CORPUS-04-R6-UNCERTAINTY-ITEMS-2026-11-12.md` | 7 R-6 units | `3d34f0530699fd33…` | 11,238 |
| `MEMORY-09` rev 2 | APPROVED; revision history | `f77e47a98685e660…` | 23,349 |
| `CORPUS-05` | protocol **defined, not executed** | `990ba06c8104b244…` | 8,814 |
| `CORPUS-07-ADJUDICATION-REGISTER-2026-11-12.md` | 18 defect / 15 clean / 12 ambiguous | (`…`) | 18,634 |
| `CORPUS-39` | `F6-02` owner decision + coverage audit | registry row 374 | 17,033 |

All quotations below are from the files, read this session.

---

## 3. Rule-provenance inventory

| Rule | Exact text | Applicability clause | Applied to batch 1? | Evidence |
|---|---|---|---|---|
| `CORPUS-01` §2 schema | `required_evidence[]` "what must exist and reach the reply"; `forbidden_evidence[]` "the leakage set" | none needed — schema predates the items | **Yes** | all 22 carry the fields |
| `CORPUS-01` §2.3 reference types | `exact` / `acceptable_range` / `expected_abstention` / `expected_correction` | none | **Yes** | `F6-02` is `acceptable_range` |
| `CORPUS-01` §2.4 four-part sufficiency | "required from **batch 2 onward**" + "Batch 1 items predate this section; each requires an **annex pass** to add the four parts before the corpus is sealed." | **explicit annex mechanism** | **Yes — annex complete** | §6: "Batch-1 annex pass (`P-2`): **complete** — 22/22 items now carry the four parts and an over-exposure boundary"; "Items with four-part reference answers: **45 / 45** (batch 1 annexed 2026-11-12)" |
| `CORPUS-01` §5 roadmap item `1a` | "Annex pass: add the four-part reference structure to all 22 batch-1 items — §2.4" | roadmap entry | **Done** | listed in the remaining-work table; §6 records completion |
| `CORPUS-01` `P-2` (§7 blockers) | "**The batch-1 annex is a blocker.** The four-part standard must apply to all 22 earlier items before sealing." | blocker | **Cleared** | §6 completion record |
| `CORPUS-01` §2.5 `M-1` | "every item states its over-exposure boundary — what extra, *true* material would fail it — or declares `over_exposure: none`" | every item | **Yes** | `F6-02` states one |
| `CORPUS-01` §2.5 `M-2` | "for restraint items (`P-01`, `F6-04`, `F1-01`, `P-04`), over-exposure is a **scored** part, not a note" | **by name only** | **Selectively** — `F1-01` is batch 1 and is named; `F6-02` is not | by-name list; no general extension |
| `CORPUS-01` §3.1 modes + `A-1`/`A-2`/`A-3` | six modes with silence-pass column; `A-2` "every item is tagged with exactly one mode" | tagging obligation | **Yes** | `F6-02` carries `Restraint: R-1`; coverage recorded "after batch 3" |
| `MEMORY-09` §3.3/§3.4 rules (`M-1`, `G-1`, sufficiency) | governing preregistration | revision 2, "safeguards `R-11`–`R-14` locked" | **Yes** — governing authority, not corpus-local | `CORPUS-01` "Authority:" line cites `MEMORY-09` §3 |
| Coverage rule "every distortion must be scored" | **does not exist** in any governing text | — | — | `CORPUS-30` §2 / `CORPUS-31` §6 only, as unadopted `D1-ENF` clause (7) |

**Two facts settle the applicability question.**

1. **The retroactivity premise is false.** `CORPUS-01` did not extend §2.4 by after-the-fact reinterpretation; it (a) named the gap, (b) made it a **sealing blocker** (`P-2`), (c) put the fix on the roadmap (`1a`), and (d) **recorded it complete on 2026-11-12**. Applying §2.4 to `F6-02` today is applying a standard through its own express mechanism, not applying it retroactively.
2. **Selective naming is the record's practice.** `M-2` scores over-exposure for four named items, one of which (`F1-01`) is batch 1. `F6-02` is **not** named. Later rules reached batch-1 items **by explicit name**, not by blanket retroactivity. That is the operative convention, and it is consistent.

---

## 4. What actually governed batch-1 items at authoring

| Governing at authoring | Basis |
|---|---|
| §2 schema, §2.3 reference types | pre-existing; all fields present |
| §2.5 `M-1` (declare a boundary) | stated as universal |
| §3.1 modes and `A-2` tagging | applied; `F6-02` tagged `R-1` |
| Owner disposition: batch 1 "accepted as a **format-validation milestone** only — not a completed corpus, and not evidence supporting any architecture choice" | `CORPUS-01` header |
| **Not** governing at authoring: §2.4's four-part scoring table | "required from batch 2 onward" |

**Limit of the evidence.** The record establishes *that* the annex happened and *what* it added (the four parts and an over-exposure boundary). It does **not** record a field-by-field diff of any item before/after, so I cannot show that `F6-02`'s current `mandatory:`/`prohibited:` wording was authored in batch 1 or introduced by the annex. This is a genuine gap in provenance and is recorded as such rather than assumed.

---

## 5. Consequences for `F6-02`

**Scoring authority today.** With the annex complete, §2.4's four-part semantics apply: `mandatory` and `prohibited` are scored, `variation` passes, `never scored` bars style. `M-1` applies (a boundary is stated). `M-2` does **not** apply — `F6-02` is not named — so its `over-exposure:` boundary is declared but carries no scoring authority. `G-1` applies to `Forbidden:`.

**So `CORPUS-39`'s substantive findings are unaffected, and its procedural caveat (G-6) is corrected:**

| `CORPUS-39` finding | Status after this audit |
|---|---|
| Behaviours 4 (endorse) and 5 (deny/neutralise) unenforced | **UNAFFECTED** — still absent from every scored field |
| Behaviour 3 permitted only because `over-exposure:` is unscored | **UNAFFECTED** — `M-2` still does not name `F6-02` |
| Behaviour 1 actor ambiguity | **UNAFFECTED** — wording, not applicability |
| G-6: "batch-1 status: §2.4 applies formally from batch 2 onward; compliance awaits an annex pass" | **CORRECTED** — the annex pass is recorded **complete** (22/22); §2.4 applies now |

**History of `Defective (partial)` — unchanged conclusion, narrower reason.** The classification is still not derivable from adopted rules, but not because §2.4 fails to apply. It is because **§2.4 contains no coverage requirement**: it says which parts score, never that a `Required` element must score every distortion of itself. The only text in `CORPUS-01` using the word "defective" is the inverse rule — "a creative item that scores a stylistic difference as a factual failure is a **defective item**" — which `F6-02` does not do. So:

- **Product defect:** real, on the owner's approved boundary (behaviours 4 and 5 cannot fail).
- **Protocol violation:** none established. The scoring that exists is authorised; the scoring that is absent is not *required* by any adopted text.
- **`Defective (partial)`:** supportable only as a coverage gap against an approved interpretation, exactly as `CORPUS-39` §6 concluded — with the reason now being the absence of a coverage rule, not the absence of applicability.

**The two failure modes the owner named are both present here:** the right behavioural expectation with the wrong instrument (behaviours 4–5 unmeasured), and a reasonable rule that never applied to this item in the form assumed (`M-2`'s scored over-exposure).

---

## 6. `Required:` / `M-1` / `R-1` — the interaction, unresolved by the text

| Source | Text | Tension |
|---|---|---|
| `CORPUS-01` §2 | `required_evidence[]` — "what must exist **and reach the reply**" | literal reading requires the aversion to appear |
| `CORPUS-01` §2.5 `M-1` | "`mandatory` is a **floor, not a checklist to be recited**"; Mandatory = "The **least** that satisfies the item. Never 'everything true about it'" | non-mention may satisfy |
| `CORPUS-01` §3.1 `R-1` | "Suppress mention — knows the fact, chooses not to name it … Can silence pass it? **Yes**" | silence passes by design |
| `MEMORY-09` §5 (owner) | "Memory should change what the model understands, not force the model to announce what it remembers" | recital is the failure mode |

**What the record does not do:** it never states whether "reach the reply" is satisfied by *behaviour* (tone, non-promotion) or requires *content*. `M-1` and `R-1` both point to behaviour; `§2`'s phrase points to content. For `R-1` items the mode row is the specific and the schema phrase the general, so the specific governs by ordinary construction — but that is a resolution *I* am applying, not one the text states.

**Three coherent resolutions, none adopted here:**
1. **Gloss `§2`** — read "reach the reply" as "be available to the reply", consistent with `M-1`/`R-1`. No text change; interpretive only.
2. **Amend `§2`** — add "for `R-1`-tagged items, behavioural reach satisfies `required_evidence`". A dated standard amendment.
3. **Per-item note** — record the resolution on `F6-02` only. Narrowest, but leaves the general question open.

**Recommendation (not a decision):** option 1, because it makes the schema consistent with two later owner statements and requires no amendment. It should still be recorded, since an unrecorded interpretation is how the "silently upgraded standard" failure begins.

---

## 7. Product behaviour vs available evaluation criteria

The owner's separation, made operational:

| Layer | Content | Authority to change |
|---|---|---|
| **Product intent** | the six approved behaviours (`CORPUS-39` §1) | owner |
| **Item text** | `F6-02`'s seven fields | owner-authorised repair; new evidence; re-review from scratch |
| **Protocol scoring authority** | §2.4 four parts (applicable via completed annex), `M-1`, `M-2` by name, `G-1`, diagnostics non-scoring | owner amendment only |
| **Unadopted proposals** | `D1-ENF` (1)–(8), including clause (7) | nothing until adopted |

A repair may add a `prohibited:` term. That gives the item **scoring authority over a behaviour** — it does not, and must not, retroactively make the historical classification correct. The ledger records what was measured under the rules then available; the repair creates a new item to be measured later.

---

## 8. Recommended gating sequence (no action taken)

1. Record the §6 resolution (`Required:` vs `M-1`/`R-1`) — interpretive note or amendment, owner's choice.
2. Decide the actor question in `mandatory:` ("the song is not pushed" → who pushes).
3. Decide the `over-exposure:` wording so supported explanation and fabricated elaboration are distinguishable — and decide whether `F6-02` should join `M-2`'s scored list (amendment) or stay unscored.
4. Only then authorise the repair, with a **new baseline** and `CORPUS-05` §5 step 6 re-review of the changed item.
5. Keep the historical ledger untouched throughout.

Steps 1–3 are questions about the standard; step 4 is a question about the item. Doing 4 first would repair an item against a standard that is still ambiguous on the very points being repaired.

---

## 9. What this does not authorize, limitations, provenance

**Does not authorize:** editing `F6-02` or any item; ledger change; amendment to `MEMORY-09`, `CORPUS-01`–`CORPUS-05`; adoption of `D1-ENF`; identifier allocation; checker/schema/runtime/model changes; experiments; annotation; sealing; `G2`.

**Limitations.** (a) No before/after annex diff exists, so which wording the annex introduced is unknown (§4). (b) The §6 resolution is my construction, not a recorded ruling. (c) `CORPUS-07`'s adjudication register was read for the 18/15/12 basis, not re-audited. (d) No rater data; `CORPUS-05` unexecuted. (e) The claim that `M-2`'s list is exhaustive rests on its own text ("for restraint items (`P-01`, `F6-04`, `F1-01`, `P-04`)"), with no general extension clause.

**Prepared by** DeepSeek (Buffy), author of `CORPUS-22`–`CORPUS-39`, with the interest that implies. This audit **corrects** `CORPUS-39`'s G-6 and contradicts the owner's "retroactive" framing, and it also contradicts the prior chain's own D-2 implication that §2.4 was inapplicable to batch 1. Both corrections run against the author's earlier position.

**Governing state preserved:** `G1` open; corpus frozen; expansion frozen; `D1` REVISE, unadopted; annotation blocked.
