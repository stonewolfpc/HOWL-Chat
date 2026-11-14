# `MEMORY-09` — `OD-M09-17` Human-Rater Admissibility (Stage Three)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-HRA-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **ADMISSIBILITY ANALYSIS — no decision adopted, no rater recruited, no session held.**
**Authority:** `MEMORY-09-WO-ARCH-02` §5.
**Question, narrowly:** does `OD-M09-17` permit qualified human raters, prohibit them, or leave it unresolved?

---

## 1. Answer

**Unresolved — by silence.**

The governing record **never names a rater type.** It contains no instance of "human", "person", "reviewer" or "model rater" in the study provisions. It says "rater" throughout (L104, L106), "two-rater study" (L88, L90), and "another qualified rater" (decision record L86).

Silence is not permission, and it is not prohibition. **A clarification is required before any human rater is appointed**, and this report states exactly what must be clarified.

## 2. The exact text

`MEMORY-09-OWNER-DECISION-ADDENDUM-1-OD-M09-15-18-2026-11-13.md`:

| L | Text (verbatim) |
|---|---|
| 88 | "**Disposition: TWO-RATER STUDY REQUIRED BEFORE PROMOTION BEYOND SINGLE-SOURCE CONFIDENCE.**" |
| 90 | "Cue-level results may be recorded and reported at **single-source confidence**. They may not support claims beyond it until a properly designed two-rater reproducibility study exists." |
| 97 | "Repeated judgements by the same source become independent evidence — **Denied**" |
| 98 | "A second reading of the same material constitutes a second independent source — **Denied by default**" |
| 99 | "Different chat sessions, or conflicting model labels, establish independent model lineage — **Denied**" |
| 100 | "Independence and provenance may be inferred — **Denied** — they must be established and documented" |
| 104 | "1. Rater qualification and identity/provenance requirements." |
| 106 | "3. Blinding and the information available to each rater." |

`MEMORY-09-OWNER-DECISION-RECORD-OD-M09-2026-11-13.md` **L86**: "**Reproducible** — another qualified rater applies the same rule consistently | **Unresolved validation requirement.**"

## 3. Classification of the textual status

| Category | Present? | Evidence |
|---|---|---|
| **Explicit permission** | **No** | No sentence admits humans |
| **Explicit prohibition** | **No** | No sentence restricts raters to models |
| **Silence** | **Yes** (primary) | The type is simply never specified |
| **Ambiguity** | **Yes** (secondary) | L99 speaks of "conflicting model labels" and "independent model lineage" — modelling language that *implies* models are among the admissible sources, without excluding others |
| **Interpretive support** | **Partial** | "Source" (L97) and "qualified rater" (L86) are type-neutral words; the record's own distinction between *source* and *session* (L97–L99) is about provenance, not species |

**The strongest available reading:** the record regulates **provenance**, not rater type. L99 denies that sessions or labels create lineage; L100 requires independence to be *documented*. On that reading a human rater is admissible **if and only if** their independence and provenance are documented to the same standard — which is a higher bar than "two people who read the docs".

**This remains an interpretation, not a permission.** It is offered because the record must be read somehow, and it is labelled so it cannot be mistaken for authority.

## 4. Requirement-by-requirement assessment

| Requirement | Satisfiable by humans? | Note |
|---|---|---|
| **Rater qualification** (E-1) | **Yes** | Qualification can be defined (read the mode definitions, pass a calibration set) and evidenced |
| **Independence** (E-2) | **Yes, conditionally** | Independent *work* is easy to arrange; independent *judgement* is not — see §5 |
| **Documented provenance** (E-1, L100) | **Yes** | Easier for humans than models: name, role, background, exposure to the corpus, declared |
| **Consistent instructions** (E-3) | **Yes** | A written packet is the same for any rater |
| **Reproducibility** (E-6) | **Yes in principle** | Nothing about α requires machine raters |
| **Shared-bias control** | **Not specified anywhere** | **Gap.** The record requires independence, never defines what would defeat it for humans |

## 5. The shared-bias objection

**Two humans are not automatically two sources.** They may share professional training, the same reading of an ambiguous acceptance field, the same exposure to prior reports in this repository, and the same instinct about what "counts" as restraint.

Under L100 this is not a technicality: independence "must be established and documented", and L97's denial of repeated same-source judgement applies to *any* source. A pair of humans who both read `CORPUS-78` and both rate "restraint" the way that report frames it are closer to one source than two.

**Practical mitigations, all documentary:** declare prior corpus exposure; exclude raters who authored or reviewed the artifacts under test; supply only the frozen packet; record disagreements verbatim so a shared-prior pattern is visible.

**None of these is required by the record today**, which is precisely why clarification is needed rather than an assumption.

## 6. Options and recommendation

| Option | Consequence | Assessment |
|---|---|---|
| **A · Declare human raters admissible**, subject to documented independence and a declared prior-exposure rule | Study becomes executable without new model capacity; independence rests on documentation, not species | **Recommended** — lowest cost, satisfies every stated requirement, and the shared-bias risk is manageable by disclosure |
| **B · Restrict raters to models** | Study blocked until the pool grows; the constraint becomes permanent rather than temporary | Not recommended — no textual basis for the restriction |
| **C · Proceed on the silence**, appointing humans without clarification | Fastest; leaves the authority question unresolved and could invalidate a future claim | Not recommended — converts an ambiguity into an unrecorded assumption |
| **D · Defer** | Nothing lost; nothing gained | Acceptable fallback |

**Recommended: A**, framed as a one-line clarification rather than an amendment to the eight requirements:

> *Owner clarification: "rater" in `OD-M09-17` is not restricted by type. A qualified human rater satisfies E-1 provided independence and provenance are documented under §5.1, and prior exposure to the artifacts under test is declared.*

**This is an owner decision, not a finding.** Nothing in this report authorises recruiting anyone, and no rating session occurred.
