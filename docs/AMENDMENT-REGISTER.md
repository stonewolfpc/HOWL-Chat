# HOWL Chat — Amendment Register

**Status:** Administrative status and placement overlay only. This register creates no decision, contract, freeze, acceptance, or implementation authorization.
**Last reconciled:** 2026-11-12
**Rule:** The amendment bodies remain byte-identical. This register resolves their **current status and home**; it does not rewrite their historical headers.

---

## 1. Purpose

The amendment trail contains current amendments, archived revisions, adopted drafts, review requests, proposals, and closure records. Several headers are historical and now misleading if read alone:

- the T2.9.5B r14 register header says `DRAFT — NOT FROZEN`, but a later freeze record freezes it;
- the TTS §17.19 draft header says `NOT ADOPTED`, but the adoption record adopted Revision 6 content;
- revisions 6, 7, and 8 of Amendment 02 still exist after revision 9 closed the cycle;
- F3C-TARGET-2 revisions 1–4 still exist after R5 became the current proposal.

This register gives each of them one status, one home, and one predecessor/successor. It exists so no one has to infer authority from a filename or a stale header.

---

## 2. Status vocabulary

| Status | Meaning |
|---|---|
| **Current planning amendment** | Current requirements/planning delta. Not architecture approval and not implementation authorization. |
| **Closed baseline** | The revision cycle is closed at this revision. It remains the requirements baseline, not an approval. |
| **Archived revision** | Preserved historical revision. Superseded for current reading. |
| **Historical draft** | Preserved draft whose content was adopted elsewhere. It is not a second governing authority. |
| **Current proposal** | Current proposal for an amendment. Proposal only; not frozen, not accepted. |
| **Historical proposal** | Superseded for review by a later proposal. Preserved as history. |
| **Decision/closure record** | Records a decision or closure. Does not create implementation authorization. |
| **Review/evidence** | Review, audit, or verification evidence. Not an amendment. |

---

## 3. Primary amendments

### 3.1 Amendment 01 — Unified Hardware Intelligence

| Field | Entry |
|---|---|
| Artifact | `docs/architecture/AMENDMENT-01-UNIFIED-HARDWARE-INTELLIGENCE-2026-11-11.md` |
| SHA-256 | `f3a0f8b9c5d19129db85d559f1880f21451cc7aa70596d37a408952c85c3f4ea` |
| Status | **Current planning amendment.** Not frozen. No architecture approval, implementation, migration, runtime, model, benchmark, or network authorization. |
| Amends | `docs/architecture/INTEGRATION-AND-ACCEPTANCE-PLAN-2026-11-11.md` |
| Home | `docs/architecture/` — hardware/resource branch of the integrated plan |
| Successor | None. It remains the controlling hardware/resource delta until a later owner decision says otherwise. |
| Cross-links | `DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md`, `INTEGRATION-AND-ACCEPTANCE-PLAN-2026-11-11.md` |

**Placement decision:** Amendment 01 is not an F3-C amendment and does not authorize F3-C work. It governs hardware intelligence, model residency, workload orchestration, readiness, release verification, and manual resource controls.

### 3.2 Amendment 02 — Conversation Continuity and Multi-Character Routing

| Field | Entry |
|---|---|
| Artifact | `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md` |
| SHA-256 | `073d0d4be768918841763a4ed714877ea8db2412053a1671fef2007a1327a640` |
| Status | **Closed requirements baseline at revision 9.** Not frozen. Not architecture approval. No revision 10. |
| Amends | `docs/architecture/INTEGRATION-AND-ACCEPTANCE-PLAN-2026-11-11.md` |
| Preserves | Original plan, Amendment 01, reviews, design documents |
| Home | `docs/architecture/` — conversation continuity branch |
| Cycle closure | `docs/review/AMENDMENT-02-REVISION-CYCLE-CLOSURE-2026-11-11.md` |
| Verification | `docs/architecture/AMENDMENT-02-VERIFICATION-RECORD-2026-11-11.md` |
| Open obligations | Muse `M6-01`–`M6-03`, recorded in the closure record; they are evidence obligations, not a revision 10. |

**Placement decision:** Revision 9 is the current requirements baseline for conversation continuity. It is **not** an architecture freeze and does not replace Amendment 01. Where it says "controlling" for its affected Chat capabilities, that is requirements control within its scope, not implementation authorization.

### 3.3 Amendment 02 archived revisions

| Revision | Artifact | SHA-256 | Status | Relationship |
|---|---|---|---|---|
| 6 | `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-REVISION-6-2026-11-11.md` | `ef2d5dd88fc30cf49424f6bcf20f50aca5c599ae6d9a8f76a33bc65673370433` | **Archived revision** | Historical; superseded by rev 7/8/9 |
| 7 | `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-REVISION-7-2026-11-11.md` | `4cd760643d57f039ecb7fc11503608c8041b8e824b1aa5258e1ed8729fb42dd1` | **Archived revision** | Historical; superseded by rev 8/9 |
| 8 | `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-REVISION-8-2026-11-11.md` | `7f97d3594f03c88dc2ca9148345ad7b7433b6cdf4ce61f6831abe2478911890f` | **Archived revision** | Historical; superseded by rev 9 |

**Placement decision:** These are preserved revision identities, not current alternatives. Read revision 9 for current requirements. Revision 3 remains the independently reviewed identity, not the current baseline.

### 3.4 TTS §17.19 capability-policy draft

| Field | Entry |
|---|---|
| Artifact | `docs/review/AMENDMENT-DRAFT-howl-tts-capability-policy-c1-17.19-2026-11-07.md` |
| SHA-256 | `dc5852adbc70641dbe8f233851eaec44266330fce197b78a0ab2641e7f73f784` |
| Header status | `DRAFT FOR OWNER REVIEW — REVISION 6; NOT ADOPTED; CONTRACT AMENDMENT ON CONDITIONAL HOLD` |
| Current status | **Historical draft. Revision 6 content was adopted on 2026-11-07 and incorporated into the governing TTS design at §17.19.** |
| Adoption record | `docs/review/ADOPTION-howl-tts-capability-policy-c1-17.19-2026-11-07.md` (`9da36691054df40d98c772f7ba9643bc959416e39a7f2d30681d9404cd16f043`) |
| Governing design | `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` (`ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f`) |
| Home | `docs/review/` — T2.9.5B / TTS evidence |

**Placement decision:** The draft's own `NOT ADOPTED` header is its pre-adoption state. The later adoption record controls. The standalone draft is preserved as history; the governing text is the design's §17.19, not a second copy of the draft.

**Why the body was not edited:** The draft is cited by adoption/review hashes. Rewriting its header would break those identities. The status overlay resolves the contradiction without changing bytes.

### 3.5 F3C-TARGET-2 proposal chain

| Revision | Artifact | SHA-256 | Status | Relationship |
|---|---|---|---|---|
| R1 | `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-2026-11-11.md` | `29f942ebd6c9f688af80d238dc7804f9a954b46c8c952e32d564b28343809353` | **Historical proposal** | Superseded for review by R2 |
| R2 | `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R2-2026-11-11.md` | `ba989dcc87d8788421e5cb554dc893863644a28f121f30dc1d5ede3352e82fd4` | **Historical proposal** | Superseded for review by R3 |
| R3 | `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R3-2026-11-11.md` | `5013eb6fad7e0ec6f067d0fdbd350fd7f7c68ff9e263a7967b37f39e6daaf42c` | **Historical proposal** | Superseded for review by R4 |
| R4 | `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R4-2026-11-11.md` | `7a2c221147621706e299924c20de1fe29080b7d649fac25264ce212d11e7417a` | **Historical proposal** | Superseded for review by R5 |
| R5 | `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R5-2026-11-11.md` | `bbe9263a3251168e4ba22d8c327e665b7dbaf50cf4692f2ddf93e62de206275e` | **Current proposal only.** Not frozen, not accepted, not implementation. | Current review target; supersedes R1–R4 for review |

**Placement decision:** R5 is the only current proposal in this chain. It does not amend the frozen F3-C contract. Its reviews and dispositions remain evidence; the frozen master contract remains controlling.

Relevant R5 review/disposition artifacts:

- `docs/review/REVIEW-ADVERSARIAL-F3C-TARGET-2-R2-2026-11-11.md`
- `docs/review/REVIEW-INDEPENDENT-F3C-TARGET-2-R3-2026-11-11.md`
- `docs/review/REVIEW-INDEPENDENT-F3C-TARGET-2-R4-2026-11-11.md`
- `docs/review/REVIEW-INDEPENDENT-F3C-TARGET-2-R5-2026-11-11.md`
- `docs/review/REVIEW-INDEPENDENT-INTEGRATION-REVIEW-F3C-TARGET-2-R5-2026-11-11.md`
- `docs/review/DISPOSITION-MATRIX-F3C-TARGET-2-R5-MUSE-2026-11-11.md`
- `docs/review/ADDENDUM-PRECISION-DISPOSITION-MATRIX-F3C-TARGET-2-R5-2026-11-11.md`
- `docs/review/ADDENDUM-2-PRECISION-RECOVERY-F3C-TARGET-2-R5-2026-11-11.md`
- `docs/review/INTEGRATION-REVIEW-F3C-TARGET-2-R5-2026-11-11.md`

None of these is an amendment to the frozen contract unless a later owner decision says so.

---

## 4. Amendment-cycle records (not amendments)

These records control or explain the amendment cycle. They are not alternative amendment texts.

| Artifact | Role | Status | Home |
|---|---|---|---|
| `docs/architecture/AMENDMENT-02-VERIFICATION-RECORD-2026-11-11.md` | Verification record | Current evidence | `docs/architecture/` |
| `docs/review/AMENDMENT-02-REVISION-CYCLE-CLOSURE-2026-11-11.md` | Closure record | Closed at rev 9; no rev 10 | `docs/review/` |
| `docs/review/AMENDMENT-02-OWNER-DECISION-REQUEST-A2-08-2026-11-11.md` | Decision record | Resolved by OD-02/03/04 | `docs/review/` |
| `docs/review/AMENDMENT-02-MUSE-PASS-ACCEPTANCE-2026-11-11.md` | PASS acceptance / advisory disposition | Review evidence | `docs/review/` |
| `docs/review/AMENDMENT-02-MUSE-REVIEW-REQUEST-2026-11-11.md` | Review request | Review evidence | `docs/review/` |
| `docs/review/REVIEW-AMENDMENT-02-REVISION-3-MUSE-2026-11-11.md` | Independent requirement-capture review | PASS; not architecture approval | `docs/review/` |
| `docs/review/REVIEW-AMENDMENT-02-REVISION-6-OD-04-CONSISTENCY-ADDENDUM-2026-11-11.md` | Delta review | Historical review evidence | `docs/review/` |
| `docs/review/REVIEW-AMENDMENT-02-REVISION-7-MUSE-DELTA-2026-11-11.md` | Delta review | Historical review evidence | `docs/review/` |
| `docs/review/REVIEW-AMENDMENT-02-REVISION-8-MUSE-DELTA-2026-11-11.md` | Delta review | Historical review evidence | `docs/review/` |
| `docs/review/REVIEW-AMENDMENT-02-REVISION-9-MUSE-DELTA-2026-11-11.md` | Delta review | PASS WITH REQUIRED CHANGES; M6 obligations | `docs/review/` |
| `docs/review/AMENDMENT-02-REVISION-7-MUSE-DELTA-REVIEW-REQUEST-2026-11-11.md` | Review request | Historical | `docs/review/` |
| `docs/review/AMENDMENT-02-REVISION-8-MUSE-DELTA-REVIEW-REQUEST-2026-11-11.md` | Review request | Historical | `docs/review/` |
| `docs/review/AMENDMENT-02-REVISION-9-MUSE-DELTA-REVIEW-REQUEST-2026-11-11.md` | Review request | Historical | `docs/review/` |

---

## 5. Amendment status resolution table

This is the table to read when a header and a later record disagree.

| Question | Controlling answer | Basis |
|---|---|---|
| Is Amendment 01 frozen? | **No.** It is a current planning amendment. | Its header; no freeze record exists. |
| Is Amendment 02 revision 9 frozen? | **No.** It is a closed requirements baseline. | Closure record; no freeze record exists. |
| Are Amendment 02 revisions 6/7/8 current? | **No.** Archived revisions. | Revision 9 and closure record. |
| Is the TTS §17.19 draft the governing text? | **No.** Its Revision 6 content was adopted into the governing design. | Adoption record; design §17.19. |
| Is F3C-TARGET-2 R5 frozen? | **No.** Proposal only. | R5 header; no freeze record. |
| Are F3C-TARGET-2 R1–R4 current? | **No.** Historical proposals. | R5 supersedes for review. |
| Does a PASS review approve architecture? | **No.** Review evidence only. | Review and closure records. |
| Does an amendment register create authorization? | **No.** Administrative status only. | This register. |

---

## 6. Preservation and maintenance

- No amendment body was edited to create this register.
- Every amendment retains its current SHA-256; the registry records the full set.
- Frozen artifacts are listed separately in `docs/FROZEN-ARTIFACTS.md` and must not be edited in place.
- A future amendment must state what it amends, what it supersedes, and which record controls. Add it to this register and the document registry in the same change.
- A stale header must be resolved by a later specific record, not by silently rewriting the historical file.

---

## 7. Bottom line

Every amendment now has one of five homes: current planning amendment, closed baseline, archived revision, historical draft, or current proposal. None is abandoned. None is silently promoted. None of the preserved bodies was rewritten. The frozen contracts and owner decisions remain the actual authority.
