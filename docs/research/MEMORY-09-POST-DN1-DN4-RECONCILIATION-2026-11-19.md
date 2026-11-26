# MEMORY-09 — Post-DN-1/DN-4 Reconciliation

**Status:** reconciliation report. Read-only. No canonical edit, no implementation, no test, no probe,
no model call, no corpus, no registry change, no Git operation.

| Field | Value |
|---|---|
| Date | **2026-11-19** (local, CST) |
| Scope | Which earlier findings are superseded, which stand, which need follow-up |
| Companion | `MEMORY-09-OWNER-DECISION-DN1-DN4-2026-11-19.md` |
| Baseline | v9 `698254ca262c0874` · r33 `dfb28b6003ca356f` · ledger `e5008020ff68bd21` · conformance `b4ba33ccbfe8f8c1` |

---

## 1. Superseded

| ID | Earlier finding | Superseded by | Label |
|---|---|---|---|
| **DB-19** | PH1A-03/04 is "the single next owner decision" | the decomposition, then DN-1 + DN-4 — the bundle no longer exists | Superseded |
| **DB-14** | PH1A-12 is blocked until PH1A-03/04 is settled | PH1A-12's predecessors are now PH1A-05 and DN-6a/b; DN-1/DN-3/DN-4 are discharged | Superseded |
| **FD-16** | the No-path exclusion list | No was not selected; the Yes consequences replace it. The three-way classification method stands | Superseded |
| **TC-09** | DN-1 "open pending confirmation" | owner confirmation, §1 of the decision record | Superseded |
| **TC-11** | whether K-R1's transposition to accounts is a gap-fill — `Unestablished` | the owner's wording ("account backed by their own legitimate acquisition event") supplies the carrier phrasing | Superseded |
| **1A-1B §3.2** | Q-B open | O-14 | Superseded |
| **1A-1B RB-02 / GB-1** | r32's internal inconsistency has no priority | O-14 classified §8.4 as illustrative | Superseded |
| **1A-1A §11 item 3, §13.2** | Q-D must precede Q-B (state-relative notice join) | O-14 removed the notice; no join to compute | Superseded |
| **1A-1B §12** | recommended order Q-B → Q-A → Q-C | Q-B settled; Q-A deferrable; Q-C moot | Superseded |

---

## 2. Still valid

| Finding | Why it stands |
|---|---|
| **Four information surfaces** (S-i…S-iv) and **handling vs awareness** | O-14 and DN-4 both depend on the separation; neither collapses it |
| **RB-04 / RB-05** (E1c narrower than "no notice"; D2 reaches only S-iv) | unchanged by DN-1/DN-4 |
| **Q-A deferrable** — 0 of the 118 require sameness | DN-1 confirms an account is backed by *its own* event; it creates no cross-account identity |
| **Q-C moot** as posed | no notice exists |
| **No new relationship model necessary**; RA-1 excluded by O-1; RA-2 sufficient | unchanged |
| **Neither implementation is a conforming subject** | strengthened: now non-conformant under YES as well (§5) |
| **V-5** (harness retraction record vs engine status-only) | unchanged; the record's semantic status remains `Unestablished` |
| **P-5 / P-7** (no inspection context; disclosure uncalled) | unchanged, still unauthorised work |
| **Audit P1–P12** blockers | unchanged by any decision here |
| **CO-02 / CO-07** (mutual incompleteness; `revoked` means three things) | unchanged; DN-4 sharpens CO-07 rather than resolving it |

---

## 3. Requires follow-up

| ID | Item | Kind | Status |
|---|---|---|---|
| **NEW-1** | Does the 118 need a forgetting row? | owner decision | open — O-6b closed the set at 118 |
| **NEW-2** | Is PH1A-05 (`turn, seq`) now blocking? | owner decision | open — YES makes same-turn ordering material |
| **DN-6a** | Who may cause forgetting | owner decision | open — no source; O-5 is the closest principle |
| **DN-6b** | Who may retract content / an account | owner decision | open — `retract` has no actor in either implementation |
| **PH1A-12** | Experimental subject and bridge | owner decision | blocked on PH1A-05 and DN-6a/b; this record selects nothing |
| **Q-D / RQS-01** | Where applicability lives on MEMORY-09's side | owner decision | open, narrowed to applicability (no notice join) |
| **PH1A-15** | No review covers the committed engine | evidence gap | open |
| **O-15** | Ledger entry recording DN-1/DN-4 | documentation | **proposed, not made** (§6) |

---

## 4. Dependency chain, current state

| Category | Items |
|---|---|
| **1. Owner decisions settled** | DN-1 (confirmation) · DN-4 = YES · O-14 · A2b/E1c/D2 |
| **2. Governed automatically by tier precedence** | K-R12 status-as-projection (v9 silent) · K-R1/K-R2/K-R11 acquisition fold (DN-1) · K-R3/K-R17/K-R19/K-R20/K-R25 (DN-4) · O-1/O-3/O-5 reach |
| **3. Remaining owner decisions** | NEW-1 · NEW-2/PH1A-05 · DN-6a · DN-6b · Q-D · PH1A-12 · PH1A-01 · PH1A-07/08/09/10/13 |
| **4. Deferrable** | Q-A (no requirement) · Q-C (moot) · PH1A-11 (Howl boundary) · PH1A-14 (application-level) |
| **5. Blockers (build, not decision)** | P1 corpus · P5/P6 bindings · P7 grant mechanism · P8 live responder · P9 bridge · P-5 inspection · P-7 disclosure path |

---

## 5. Contradictions, gaps and cautions

1. **K-R15 cannot be satisfied as written.** r32's own `AcquisitionEvent` schema has `turn` but no `seq`
   (ARCH-03 F6). Under YES this stops being cosmetic: a forget and a learn in the same turn are
   ambiguous without it. *(Observed + Source-derived)*
2. **The 118 have no forgetting row.** YES creates a T2 obligation with no conformance coverage
   (NEW-1). *(Calculated)*
3. **A2a must not be read as prohibiting forgetting.** Its row context is source retraction; the same is
   true of C1 at row level. *(Source-derived)*
4. **The harness's `revoke()` writes both an acquisition event and a record status.** That is a
   stream conflation in its own right — K-R12 places status on the memory stream, K-R11 on the
   acquisition stream. DN-4 makes the conflation visible but does not create it. *(Observed)*
5. **No actor exists on retraction in either implementation.** DN-6b inherits an unimplemented
   precondition, not a policy gap alone. *(Observed)*
6. **The proposed ledger entry does not exist yet.** Until O-15 is authorised, this decision lives only
   in this record and its companion. *(Observed)*

---

## 6. Proposed canonical change — identified, not made

**Ledger:** add `### O-15 — DN-1 and DN-4 (CONFIRMED)` to
`MEMORY-09-OWNER-DECISIONS-2026-11-17.md`, carrying the verbatim DN-1 wording, DN-4 = YES, the six
established items and the eight excluded ones, and pointers to this record. **Not made:** the work order
prohibits changes to the existing ledger without further authorisation.

No change is proposed to v9, r33, or the conformance map.

---

## 7. Files

| State | Files |
|---|---|
| **Created** | `MEMORY-09-OWNER-DECISION-DN1-DN4-2026-11-19.md` · `MEMORY-09-POST-DN1-DN4-RECONCILIATION-2026-11-19.md` |
| **Verified unchanged** | r33 `dfb28b6003ca356f` · v9 `698254ca262c0874` · ledger `e5008020ff68bd21` · conformance `b4ba33ccbfe8f8c1` |
| **Not touched** | all `eval/` code · registry · frozen artifacts |

**End of reconciliation.**
