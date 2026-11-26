# MEMORY-09 — S-2 Authorization Addendum

## MEMORY-09 / Howl state-machine separation

**Status:** **binding addendum** to the S-2 authorization request (`MEMORY-09-S2-AUTHORIZATION-REQUEST-2026-11-20.md`,
S2R-01…S2R-35). It is part of the authorization package and binds any A/B/C/D selection. **It is not an authorization;
no carrier, schema, code, requirement, map change or canonical capture is proposed or authorized.** The read-only
boundary remains intact.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Adds | the boundary between the **MEMORY-09 memory contract** and the **Howl state machine** that consumes it |
| Binds | any A/B/C/D authorization of S-2 |
| Does not do | authorize S-2; select a carrier; design the state machine; add a requirement; change the map |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. The layering

**S2H-01 (Observed — the owner's framing, 2026-11-20).** The memory system should define what Howl can reliably
remember and retrieve. The Howl state machine should determine how those memories are selected, switched and
combined. The memory system significantly constrains and informs the state-machine design; it should not
independently dictate the entire state machine.

**S2H-02 (Calculated — the three layers).**

| Layer | Responsibility | Examples |
|---|---|---|
| **MEMORY-09 — memory foundation** | defines memory records, lifecycle, scope, knowledge classification, history and reconstruction; establishes the **memory contract/capabilities** | identity, scope, class, lifecycle, projection, history, rebuild |
| **Howl — state machine** | **consumes** the memory contract to select, switch, combine, activate, deactivate and restore context | which memories are in the current combination; switch and switch-back |
| **User experience** | the outcome the application must deliver | characters retain the right knowledge in the right world/scenario, without accidental leakage or loss |

**S2H-03 (Calculated — the direction of constraint).** Memory constrains and informs the state machine; the state
machine does not define memory semantics. The state machine is a **separate Howl responsibility**, built using the
memory contract (owner, 2026-11-20).

**S2H-04 (Calculated — why this matters to S-2).** S-2 is a **MEMORY-09 representation decision**. It investigates
the memory system's ability to represent and reconstruct changing `scope` and `knowledge_class`. It does not solve how
the state machine consumes those capabilities. Treating S-2 as a state-machine decision would put application
orchestration inside the memory foundation.

**S2H-05 (Calculated — what the separation is not).** It is not a prohibition on MEMORY-09 producing information the
future state machine needs. The memory contract must be **sufficient** for the state machine to answer identity,
scope, lifecycle, history and reconstruction questions. The separation concerns **who decides what**, not whether the
contract is informative enough.

---

## 2. What S-2 may investigate

**S2H-06 (Calculated — in scope for S-2).**

| In scope | Why |
|---|---|
| the memory system's ability to represent a post-formation class/scope change | the representation question |
| identity, scope, lifecycle, history and reconstruction as memory capabilities | the contract the state machine will consume |
| the projection rule and check that make the current value derivable | K-R12/§6.4 |
| the logical entry shape for a class/scope record | representation, not fold |
| whether the contract provides the information a future state machine would need | sufficiency of the memory contract, not its consumption |
| a reasoned conclusion that no admissible representation is currently justified | a legitimate process-level outcome (S2H-15) |

**S2H-07 (Calculated — the "no artificial restriction" principle).** S-2 must **not** reject a representation solely
on the ground that it would expose information the state machine does not yet consume. If a representation gives the
future state machine reliable answers about identity, scope, lifecycle, history or reconstruction, that is a virtue of
the memory contract, not a violation of the separation.

**S2H-08 (Calculated — what "sufficient" does not mean).** It does not mean S-2 designs the state machine's selection
algorithm, context-switching rules or UX. It means the memory contract must not be made artificially narrow for fear
that the state machine might use it.

**S2H-09 (Calculated — no new requirement).** The sufficiency principle is a **constraint on S-2's reasoning**, not a
new canonical requirement. It adds no row to the 118-row map and no obligation beyond K-R12/K-R14/§6.4. The memory
contract's capabilities are already the subject of those rules.

**S2H-10 (Calculated — the test of scope).** A question belongs to S-2 if it asks **what the memory system can
represent, record or reconstruct**. A question belongs to the state machine if it asks **which memories are selected,
switched, combined, activated, deactivated or restored in a given context**.

---

## 3. What S-2 must not solve

**S2H-11 (Calculated — out of scope for S-2).**

| Out of scope | Owner |
|---|---|
| which memories belong in the current world/scenario/character combination | Howl state machine |
| which memories become active or inactive during a switch | Howl state machine |
| which memories become active in a new combination | Howl state machine |
| what happens on switch-back | Howl state machine |
| how contexts are selected, combined or restored | Howl state machine |
| UX behavior and the user-visible result | application behavior |

**S2H-12 (Calculated — the critical distinction).** State-machine **activation/deactivation** is **not** the memory
lifecycle `status` (`provisional` / `active` / `superseded` / `corrected` / `retracted`). The memory status describes
the record's lifecycle; context activation describes whether a memory is included in the current scene. S-2 must not
conflate them, and must not treat "inactive during a switch" as a lifecycle status.

**S2H-13 (Calculated — the example is an illustration, not a requirement).** The owner's example — World A with
Scenarios 1 and 2, World B with Scenario 3, three characters in different combinations — illustrates **state-machine
orchestration**. It is not an S-2 requirement, a carrier design or a lifecycle status. It is recorded here only to
fix the boundary.

**S2H-14 (Source-derived — the example's division).**

| Question on switching from World A / Scenario 1 to World B / Scenario 3 | Layer |
|---|---|
| which memories belong to World A, Scenario 1, or an individual character? | **memory contract** — identity and `scope` |
| which memories should become inactive during the switch? | **state machine** — context activation |
| which memories should become active in the new combination? | **state machine** — context activation |
| what happens on switch-back? | **state machine** — restoration; the memory contract supplies history/reconstruction |
| what is the lifecycle/class of a record? | **memory contract** — lifecycle and class |

**S2H-15 (Calculated — process-level decline remains legitimate).** S-2 may conclude that no admissible representation
is currently justified. That conclusion is a **completed substantive decision**, not a failure to complete the
decision and not merely the owner declining to authorize further work. It records O-KR12 as unmet, preserves Q2 = Yes,
and leaves the representation question open. The future state machine's needs do not force a positive carrier.

---

## 4. Authorization conditions

**S2H-16 (Calculated — binding on every option).** Any A/B/C/D authorization of S-2 includes this separation. A and B
differ in whether A-6 applies to S-2; they do not differ in whether S-2 may design the state machine. Neither option
authorizes state-machine design, UX behavior or application orchestration.

**S2H-17 (Calculated — mutual non-dictation).** S-2's decisions about memory representation do not bind the state
machine's design; the state machine's design does not dictate memory representation beyond the contract's required
capabilities. The memory contract constrains and informs; it does not dictate the entire state machine.

**S2H-18 (Calculated — no implementation or capture).** This addendum authorizes no state-machine implementation, no
UX implementation, no application orchestration, no schema, no code, no migration, no requirement capture and no
canonical capture. The map stays at **118**.

**S2H-19 (Calculated — no implicit state-machine requirement).** The memory contract's sufficiency is not a new
conformance requirement. If a future state-machine design needs a memory capability the contract does not provide,
that is a **new question** for a later authorization, not a requirement created by this addendum.

---

## 5. What this addendum does not do

**S2H-20 (Observed).**

| Untouched | State |
|---|---|
| **S-2** | not opened; not authorized |
| **carrier / schema / kind / field** | none selected or proposed |
| **state-machine design / UX** | not authorized; not designed |
| **memory lifecycle status vs context activation** | distinguished, not merged |
| **authority / validity / causes** | fenced, not decided |
| **N-b / N-c / N-a fold** | not authorized |
| **implementation / requirement capture / canonical capture** | not authorized |
| **map / requirements** | **118**; no change |
| **code / schema / migration / tests** | **0** |

---

## 6. Traceability

**S2H-21 (Observed).**

| What | Where |
|---|---|
| Owner's MEMORY-09 / Howl / UX framing | owner, 2026-11-20 |
| MEMORY-09 as memory foundation; state machine as consumer | S2H-01…S2H-05 |
| In-scope memory questions | S2H-06…S2H-10 |
| Out-of-scope state-machine questions | S2H-11…S2H-14 |
| Process-level decline | S2H-15; S2R-27…29; S2V-16 |
| A/B/C/D subject to this addendum | S2H-16; S2R-31 |
| K-R12 / K-R14 / §6.4 | S2R-14; r33 `:111`–`:113`, `:366`–`:367` |
| Non-absorption of authority, validity, N-b, N-c, N-a | S2R-17; S2A-38 |
| Separate gates | S2R-25 |
| Map unchanged | S2R-33; 118 rows |

---

## 7. Record status

**S2H-22 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — the map still holds **118** rows |
| Carrier selected; S-2 opened | **0** |
| State-machine design / UX authorized | **0** |
| Authorization request updated to reference this addendum | **yes** (S2R header and S2R-17/30/31/33/34) |
| Tests / probes / model calls | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **binding addendum** — no canonical capture |

**End of addendum.** MEMORY-09 defines the **memory contract**; the Howl state machine **consumes** it; UX is the
application outcome. S-2 may investigate the memory system's ability to represent and reconstruct changing
`scope`/`knowledge_class`, and must not artificially narrow the contract to avoid future state-machine use
(S2H-06…S2H-10). S-2 must not solve selection, switching, combination, activation, deactivation, restore or UX
(S2H-11…S2H-14), and must not conflate context activation with lifecycle `status` (S2H-12). **Any A/B/C/D authorization
is subject to this addendum** (S2H-16). **Process-level decline remains a legitimate completed outcome** (S2H-15). **No
carrier, schema, requirement, map change, canonical capture, state-machine design or authorization follows from this
addendum.**
