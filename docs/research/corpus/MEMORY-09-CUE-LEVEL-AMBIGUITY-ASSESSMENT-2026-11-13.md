# `MEMORY-09` — "Cue-Level" Operational Ambiguity Assessment

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-CLA-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **ASSESSMENT — concludes no instrument is needed yet.**
**Trigger:** owner instruction: establish whether the undefined phrase creates real operational ambiguity **before** expanding standards work.

---

## 1. Question

The audit (row 476) found that `OD-M09-17` never defines **"cue-level"**, so its scope boundary rests on an undefined term. Does that undefined term actually prevent any decision we need to make?

**Test applied:** for each claim type we might need to classify, is the classification decidable *without* a definition?

## 2. Classification test

| Claim type | Is it "cue-level"? | Decidable without a definition? | Why |
|---|---|---|---|
| A rater scores a corpus item against an acceptance field | **Yes** | **Yes** | It is a judgement about a cue's satisfaction — the addendum's own subject (§1 L15) |
| Two raters compare labels on the same unit (α) | **Yes** | **Yes** | Same subject, with the two-rater condition met |
| Owner's qualitative impression of a reply | **No** | **Yes** | Not a judgement against a cue criterion; L3 |
| Vendor or public benchmark result | **No** | **Yes** | Not produced under this corpus or its standard |
| Engineering observation (inventory honoured, clothing persisted) | **No** | **Yes** | Observation of state honouring, not judgement of a cue |
| Mechanical state check (hash, diff) | **No** | **Yes** | Deterministic, no judgement at all |
| **"The model complied *because* the cue was present"** | **Ambiguous** | **No** | Mixes a cue judgement with a behavioural observation |

**Six of seven classify cleanly.** The ambiguity is real but confined to one claim type.

## 3. The one genuinely ambiguous case

A **bridging claim** — one that uses a cue judgement to explain a behavioural observation — cannot be classified without the definition. It is also the claim the architecture already refuses to make:

- Row 470: attribution requires a contrast, never a single reply.
- Row 473 (`C2` vs `C3`): a correct reply does not establish that state caused it.
- Row 470 §6: `S-5` forbids statements about internal state.

**So the ambiguous case is the case already excluded by the claim ceiling.** The undefined term does not block any permitted decision; it blocks a claim we have already decided not to make.

## 4. Conclusion

**No instrument (`OD-M09-19`) is justified now.** The ambiguity is:
- **narrow** — one claim type of seven;
- **already excluded** by the architecture's claim ceiling;
- **cheap to avoid** by an operational rule rather than a standards amendment.

Expanding standards work to define a term whose only ambiguous application is already prohibited would add governance overhead without changing any decision.

## 5. What to do instead — the no-bridge rule

> **Results from the cue stream and the engineering stream may not be combined into a single claim.** Each stands alone, at its own confidence level, under its own authorisation. A statement that explains one by the other is a bridging claim and is **out of scope until explicitly authorised**.

This is a **record-keeping rule**, not a standards change. It costs one sentence per report and preserves the distinction the audit established.

## 6. Conditions under which this conclusion should be revisited

| Trigger | Why it would change the answer |
|---|---|
| A report needs a bridging claim to be useful | The ambiguous case becomes live |
| A reviewer disputes a classification in practice | The boundary is being used, and inconsistently |
| A test design requires attributing a behaviour to a cue | Same |
| The owner wants engineering results citable alongside cue results | Requires a stated bridging rule, hence a definition |

**Until one of those occurs, clarification is premature** — and if one does occur, the instrument should be written then, against the actual disputed case rather than in the abstract.

## 7. Claim ceiling

This assessment establishes that six claim types classify cleanly and one is ambiguous, and that the ambiguous one is already excluded. It does **not** establish that the term is well defined, nor that no future dispute will arise.
