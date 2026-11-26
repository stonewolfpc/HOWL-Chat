# MEMORY-09 — Experiment Plan Amendment 1: conformance probes

**Status:** **OWNER-APPROVED AMENDMENT — registered.** Numbered, dated and hashed under the
amendment-register convention of `MEMORY-09-PREREGISTRATION-FREEZE-2026-11-15.md` §4. Recorded
**2026-11-17 08:13:10 CST (−0600)** = **2026-11-17 14:13:10 UTC**.

**Amendment number:** 1
**Authority:** owner decision **O-7, Option B**, 2026-11-17 (see
`MEMORY-09-OWNER-DECISIONS-2026-11-17.md`).
**Amends:** the scope of one bullet of the approved experiment plan. **The plan is not edited.**
**This document's SHA-256** is recorded in the O-7 decision entry, not inside this file — a document
cannot contain its own digest.

| Field | Value |
|---|---|
| Basis document | `docs/research/MEMORY-09-EXPERIMENT-PLAN-2026-11-15.md`, **v3, approved at G0** |
| Basis identity | **21,616 B**, SHA-256 `6c4078bd8c1fd29f1c74691afd0e25273000ac8bbddf6cc5c7f162bcfbe6bbb9` |
| Basis condition | **Unmodified.** Verified byte-identical when this amendment was written |
| Requirement set | `MEMORY-09-CONFORMANCE-MAP-2026-11-17.md` — **118** requirements, ratified as requirements under O-6b |
| Contract | `MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md` r32, `e88ef2d1…` — **not amended** |

---

## 1. Amendment text

The amended scope of plan §10, bullet 2:

> Notwithstanding plan §10, bullet 2, conformance probes that operationalize only the 118
> requirements already ratified under O-6b are permitted.

The original bullet — *"No new architecture, rule, field or probe in the contract while this plan
executes"* — is **not deleted, not reworded and not replaced**. This amendment adds a narrow,
prospective permission alongside it.

## 2. What this amendment does not claim

- It does **not** establish the original intent of plan §10, bullet 2. That question was never
  settled and is not settled here.
- It does **not** attribute authorship of §10 or of the phrase **"in the contract"** to the owner.
  That phrasing arose in collaborative assistant/DeepSeek work; its authorship and original intent
  are **not** attributed to the owner.
- It does **not** reinterpret the plan retrospectively. The permission is **prospective only**: it
  governs probes defined from this date forward.
- It does **not** convert a probe definition into a result. Nothing here is executed.

## 3. What is permitted

Defining and running conformance probes that **correspond to the 118 requirements ratified under
O-6b** — those requirements, and **no others**. The set is closed at 118 as recorded in the
conformance map (`b4ba33cc…`).

## 4. What is not permitted

- **Any new requirement, architecture, rule or field.** Ratification under O-6b is closed. A
  behaviour not expressed in the 118 stays unwritten until a separate owner decision.
- **Any edit to r32**, or any change altering its hash `e88ef2d1…`.
- **Any edit to plan v3**, including its non-goals other than the narrow permission in §1.
- The plan's remaining §10 non-goals: no Howl Chat production integration, no UI, no runtime change;
  no re-run of the bge evaluation; no new reranker candidate; no threshold, case or prompt tuning
  against the evaluation set; no commit or push as part of the experiment.
- **Implementation of any `eval/mem1` code.** See §5.

## 5. Gate separation

| Gate | Question | Status | Instrument |
|---|---|---|---|
| **Gate 1** | May conformance probes for the ratified requirements be defined and run? | **Permitted** | This amendment + O-7 |
| **Gate 2** | May `eval/mem1` implementation code be written? | **UNANSWERED — not authorised** | A separate owner decision that does not yet exist |

Gate 1 does not authorise Gate 2. The G0 receipt states the same boundary independently:
*"A G0 pass authorizes nothing beyond scope: implementation still requires separate authorization."*

**Consequence, stated plainly:** with Gate 1 alone, the 118 requirements can be specified as probes
but cannot be executed against anything, because no engine exists to run them against. That is a
permitted state, not a defect.

## 6. Registration and instrument

Following `MEMORY-09-PREREGISTRATION-FREEZE-2026-11-15.md` §4: an amendment is numbered, dated,
hashed, and **registered in the decision record**. The freeze record's rule is adopted verbatim:

> An unregistered amendment does not bind.

This amendment is registered in the O-7 entry of
`MEMORY-09-OWNER-DECISIONS-2026-11-17.md` with its SHA-256. If the bytes of this file change, the
registered hash no longer matches and this amendment **stops binding** until re-registered.

## 7. Verification performed when writing this record

Observed, not inferred:

| Check | Result |
|---|---|
| Plan v3 bytes / hash | 21,616 B, `6c4078bd…` — **matches the G0-approved identity** |
| r32 hash | `e88ef2d1697da1d9…` — **unchanged** |
| Architecture rules | **57** in v9 |
| Conformance requirements | **118** — A 5 · B 11 · C 9 · D 10 · E 5 · F 9 · G 3 · H 10 · I 3 · J 16 · K 27 · L 10 |
| Tests implemented / executed | **0 / 0** |
| Source artifacts modified | **None** — plan, r32 and `eval/mem1` untouched |

## 8. Provenance note

O-7 was originally titled "O-7 — R-5a". `R-5a` is **not a source-defined identifier**: it appears
only in four internally authored records and in no source document. The authority for the
restriction is plan v3 §10, bullet 2, and O-7 has been re-anchored to it. Historical occurrences of
`R-5a` are preserved as historical references, never as independent authority. See session
correction **#24**.
