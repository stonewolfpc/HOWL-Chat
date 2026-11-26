# MEMORY-09 — Q-B Owner Decision and r32 §8.4 Reconciliation

**Status:** decision record. **Q-B is answered.** This file records the answer, the documentation changes
made to agree with it, and what it does **not** authorise. No implementation, test, probe or corpus work
was performed.

| Field | Value |
|---|---|
| Decision | **Q-B = Branch 2** — no status-derived character awareness |
| Recorded as | **O-14** in `MEMORY-09-OWNER-DECISIONS-2026-11-17.md` |
| Date | **2026-11-19** (local, CST) |
| Owner authority | Stated in session and recorded in the ledger; as with every prior confirmation, **no persisted signature exists** |
| Settles | Q-B only |
| Does not settle | Q-A, Q-C, Q-D, the `inspection` build, the disclosure path, V-5's semantics |

---

## 1. The decision

**Owner decision, verbatim:**

> Characters do not automatically become aware of retractions because application status changes.
>
> Characters retain what they acquired. They learn about a retraction through legitimate acquisition:
> witnessing, disclosure, reading, inference, verification, or another valid route.
>
> MEMORY-09 records and handles retraction independently of character awareness.

**In the question's own terms.** A character's usable context does not report a withdrawal that the
character did not acquire. The application continues to record the withdrawal, keep its history, and
stop treating the claim as current truth.

**What it agrees with.** A2b, E1c and D2 required no amendment, because the answer is the one they
already implied. This is the decision 1A-1B predicted would require no tier-1 change — the branch that
fits the ratified set without amendment.

**What it forecloses.** The §8.4 display as character-usable knowledge. That display was the only
behaviour in the 1A-1B matrix that satisfied §8.4 while violating A2b, E1c, D2 and K-R7. It is now
classified rather than implemented.

---

## 2. Documentation changes made

Three edits, all documentation. Nothing in `eval/` was touched.

| # | File | Change | Before | After |
|---|---|---|---|---|
| 1 | `MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md` | Header revision **r32 → r33**; O-14 note added beneath §8.4; §8.4 rule text preserved unedited | 1028 lines, `e88ef2d1697da1d9…` | **1040 lines, `dfb28b6003ca356f…`** |
| 2 | `MEMORY-09-OWNER-DECISIONS-2026-11-17.md` | **O-14** added at `:205`, with the broad-reading paragraph at `:228`; dated qualification on O-6b's "r32 is byte-identical" statement | 636 lines, `da3e72398ecbaa0a…` | **688 lines, `e5008020ff68bd21…`** |
| 3 | This file | Record of the decision and its consequences | — | new |

Full SHA-256 after the edits:

- r33 — `dfb28b6003ca356f237d5e9bf127ba978d8a35e9f15b258b509eb1beeee5f158`
- Ledger — `e5008020ff68bd21b41d03d09de4542f6fc2d0ecddfc1111cc5b48dae71dfc3c`

### 2.1 What the §8.4 note says

The original sentence — *"A's channel: 'you were told X; it has since been retracted; you still
believe it'"* — remains in the document exactly as written. A blockquote beneath it records the
classification: illustrative, not a requirement, not to be implemented as character-usable knowledge; a
retraction reaches a channel only through an `AcquisitionEvent` (K-R7); the character's knowledge state
in the **human** inspection view is unaffected.

The blockquote also fixes how that phrase is to be read. **Acquisition is broad, not disclosure-shaped:**
`witnessed`, `told`, `read`, `inferred`, `verified`, `narrator_grant` (K-R11). An inference or a
verification acquires without any disclosure occurring, and the event records it. The prohibition is
narrow and specific — deriving an acquisition from status — not a requirement that someone speak.

**Why a note and not a deletion.** Deleting the sentence would erase the evidence of the conflict that
RB-02 identified and would make the revision unauditable. The contract's own practice is to preserve
rule text and record the change around it.

### 2.2 Why the revision number moved

r32 was a labelling correction; this is a classification of a worked example. The header records r33
first and keeps r32's note as "Previous revision", so the history stays readable. No `K-R` rule, no
invariant, and no probe expectation was altered.

---

## 3. What deliberately did not change

| Item | Status | Reason |
|---|---|---|
| **A2b, E1c, D2** | unchanged | the decision agrees with them; amending them would now be wrong |
| **O-1, O-3, O-5, H10** | unchanged | untouched by the question |
| **RA-1 / RA-2 / RA-3** | none adopted | Q-A is still open; no content identity introduced |
| **K-R7, K-R13 + consequence** | unchanged | they were the side of the conflict the decision chose |
| **118 conformance requirements** | unchanged | O-14 creates no requirement and removes none |
| **`eval/mem1`** | unchanged | 0 files edited, 0 tests run |
| **`DOCUMENT-REGISTRY.md`** | unchanged | registry edits were not part of this direction and remain unauthorised |

**One consequence to state plainly.** Every record that cites r32 as `e88ef2d1697da1d9…` now cites a
**superseded revision** — **30 files** under `docs/` at the time of this edit (`grep -rl e88ef2d1 docs/`).
Those records are not rewritten here: they were accurate at the hash they name, and the 1A-1 series is a
sequence of dated records, not a living document. Any future work needing current r32 text must read
**r33** (`dfb28b6003ca356f…`).

One further ledger mention is deliberately left as it stands: O-7's bounds line (`:584`, "does **not**
alter r32") records a **decision-time bound**, not a present-tense identity claim. Gate 2's "no changes
to r32" restriction is a bound of that authorisation in the same way (see §5). The only present-tense
identity claim — O-6b's "r32 is byte-identical" at `:527` — carries a dated qualification at `:529`.

---

## 4. What the decision closes, and what it does not

### 4.1 Closed

| ID | Was | Now |
|---|---|---|
| **Q-B** | open since 1A-1 | **answered** — Branch 2, O-14 |
| **GB-1** (1A-1B) | r32's three texts disagree with no stated priority | priority supplied by owner classification; the disagreement is recorded, not erased |
| **1A-1C §6 item 5** | §8.4 might need demoting | performed, as a classification, under O-14 |

### 4.2 Still open, now more load-bearing

| ID | Item | Why O-14 raises its importance |
|---|---|---|
| **P-5** | no `inspection` context in either implementation | the operator's ability to see a retraction is now the *only* place the system's knowledge is visible to a human |
| **P-7** | disclosure implemented but uncalled | acquisition is now the only route to character awareness, so the route must exist |
| **V-5** | v9 generates a retraction record; m09 does not; the record's meaning is Unestablished | whichever shape a notice design extends must not become a status-derived injection |
| **Q-A** | shared content identity | no longer forced by retraction; if pursued, it must be justified on other grounds |
| **Q-C** | attribution | moot for as long as no notice is shown; would return only with a new decision |
| **Q-D** | conversation-state applicability | unchanged; governs *from what* a notice would be computed, which O-14 now answers "not at all" |
| **GB-2** | A2b's "same character, same turn" | still unread, but no longer load-bearing for this question |

---

## 5. What this decision does not authorise

No implementation. No test, probe, selftest or benchmark execution. No corpus work. No change to the
118 requirements, to v9, to the experiment plan or to Amendment 1. No new relation, field, rule or
status value. No registry row. No persisted signature, and none implied.

**O-14 is a behavioural decision, not a build order.** The two build gaps it exposes — an inspection
surface and an exercised disclosure path — remain exactly as unbuilt as 1A-1C found them.

**A note on scope.** Gate 2's restriction ("no changes to r32") governed *implementation* work. The r33
classification was made under the owner's explicit Q-B documentation direction, which is a different
authority; the distinction is recorded here so the two are not confused later.

---

## 6. Plain statement of the resulting behaviour

For a character who was told something, and a later withdrawal of it:

- they keep what they were told;
- they may keep believing it;
- the application records the withdrawal and stops treating the claim as current;
- the operator can see the withdrawal and its history;
- the character learns of it only if they acquire it — told, witnessed, read, inferred, verified, or
  another valid route, each recorded as an acquisition event.

The character is not corrected behind the scenes. That is the whole of the decision.

---

## 7. Verification

| Check | Result |
|---|---|
| Files edited | 2 canonical + 1 new record |
| Implementation files touched | **0** |
| Tests / probes executed | **0** |
| r33 note present at §8.4 | yes, `:427` |
| Acquisition read broadly at §8.4 | yes — six K-R11 kinds; inference ≠ disclosure |
| O-14 broad-reading paragraph | yes, `:228` |
| Scope-freeze section | yes, `:176` |
| §8.4 original text preserved | yes |
| O-14 present once in the ledger | yes, `:205` |
| O-6b qualification present | yes, `:529` |
| Duplicate sections introduced | none |
| Hashes recomputed after edit | yes, §2 above |

## 8. Scope freeze

**The documentation scope of Q-B is closed as of this record.** No further edit to r33, to O-14, or to
this file is made under Q-B. The decision, its classification of §8.4 and its record of open gaps are
complete; anything beyond them is a different decision.

Q-A, Q-C, Q-D, content identity, the object model, the inspection surface and the disclosure path all
remain where this decision left them. None is a starting point implied by O-14, and none may be
treated as authorised because the documentation around Q-B is now tidy.

**End of record.**
