# Muse addendum — OD-04 consistency across continuity surfaces

**Date:** 2026-11-11
**Reviewer:** Muse
**Review type:** Targeted policy-consistency audit (read-only)
**Subject:** `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md`, revision 6
**Requested by:** Owner, specifically to test whether `OD-04` is reflected consistently across six surfaces and whether any acceptance criterion contradicts the policy
**Status:** Complete; no document or source file modified

## 1. Verdict

**PASS WITH REQUIRED CHANGES**, scoped to `OD-04` enforcement.

The `OD-04` policy text is internally consistent and appears correctly in every structural location where the amendment could have drifted: the decision record (§1.4), the transition section (§7.7), the transition record fields (§7.2), the decision register (`DA-C4`), the open-question register (§15, where `O-11` is recorded as resolved rather than deleted), the conflict table (`A2-08`), and the research register (§12.2, `R-09`).

The failure mode you named did not occur in prose. No acceptance criterion states the opposite of `OD-04`, and no section quietly narrows it.

The failure mode that *did* occur is quieter and more dangerous at a gate: **`OD-04` has no acceptance criterion at all.** A sweep of §16.1–§16.7, `AC-C01` through `AC-C66`, finds no criterion that tests explicit-change preservation, scope visibility, deliberate promotion, or the historical record of a change. An implementation could satisfy all 66 criteria while promoting every chat-typed change to character scope — precisely what §1.4 forbids.

One further item sits on the edge of outright contradiction rather than omission: `AC-C23` requires scenario selection to "immediately" record and display the transition, and §7.1 states its eight-step transaction without a confirmation gate or a decline path. Read literally, a transition completes and reaches model context before §7.7's confirmation could apply.

The revision-3 requirement-capture PASS is unaffected. These are documentation changes and do not block `R-07`–`R-11` research, which already carries the affected questions.

## 2. Package identity

| Artifact | Value |
|---|---|
| Amendment 02, revision 6 | `ef2d5dd88fc30cf49424f6bcf20f50aca5c599ae6d9a8f76a33bc65673370433` |
| Size / lines | `71,428 B / 1,205 lines` |
| Sidecar | `sha256sum -c` → `OK` |
| Integration plan | `1e92a54434e8e3269015c47374ad57bd6ae755e64e27a07954c083ab952510cf` — unchanged |
| Amendment 01 | `f3a0f8b9c5d19129db85d559f1880f21451cc7aa70596d37a408952c85c3f4ea` — unchanged |
| Revision-3 review | `0c82874b8fc5c0b82aa3a0e42b83802a8e21b2a2384b453758a4c5a429b1acf3` — unchanged |

Revision 6 supersedes revision 3 as the current identity. The revision-3 review reviewed a document that no longer exists on disk under that content, which is expected: the package declares revision 3 as the independently reviewed identity and records revisions 4–6 as controlled deltas.

## 3. The policy being tested

§1.4 decomposes into six testable obligations:

| # | Obligation |
|---|---|
| `P1` | An explicit user change is a distinct operation, not an inference from typing. |
| `P2` | It must not be silently overwritten by a scenario reload or transition. |
| `P3` | It is preserved at the narrowest scope — the current conversation — unless the user explicitly selects a broader scope. |
| `P4` | A character-persistent change and a scenario edit are separate, explicit operations; neither may be inferred from Chat input. |
| `P5` | The user can see the scope of a change and promote it deliberately. |
| `P6` | The historical record of what changed and why is retained. |

`OD-03` supplies the adjacent obligation: replacing or discarding an explicit change requires confirmation, and undo is not a substitute.

## 4. Surface-by-surface result

| Surface | Policy coverage | Enforcement coverage | Result |
|---|---|---|---|
| Scenario transitions | §1.4 `P2`; §7.7 final paragraph; §7.2 records "the change reason and actor" | `AC-C23`–`AC-C31` contain no explicit-change criterion; `AC-C23` asserts "immediately" | **Gap** — `M3-01`, `M3-02` |
| Character state and memory | §1.4 `P4`; §7.6 "must not invent personal experience"; `OD-02` cross-chat rule | No criterion constrains the scope of a change that enters memory; §9.3 lists corrections without a scope rule | **Gap** — `M3-03` |
| Chat-specific lifecycle events | §1.2; §7.6 (Step Out, return, death, ghost are chat-scoped) | `AC-C61` covers only participant-list side effects of private input | **Gap** — `M3-04` |
| Explicit user corrections | §1.4; §9.3 lists "corrections and contradictory information" | None | **Gap** — `M3-03` |
| Private information and derived memories | §1.4; §5.3; §10.3; §10.4; §15 `O-11` | `AC-C49`, `AC-C54`, `AC-C65`, `AC-C66`, `AC-C62`–`AC-C64` | **Sound** |
| Historical records versus active state | §1.2; §4.1 invariants; §7.2 reason and actor | `AC-C28`, `AC-C31` cover only "old scenario state is not current" | **Partial** — `M3-05` |

### 4.1 Private information and derived memories — why this surface passes

`OD-04`'s most probative intersection is the promotion rule: a conversation-scoped change becoming character- or scenario-persistent. §10.4 and `AC-C65` handle it correctly — a private-derived memory retrieved in public context "remains inaccessible or is transformed only under an approved visibility change," and §15 states that `O-11`'s resolution "may not be used to narrow the owner policy in §1.4."

The audience model and the change-scope model therefore meet at a consistent point rather than contradicting. `AC-C66` and §5.3's receipt rule also correctly inherit the strictest audience of referenced content, which is the `M2-01` correction now closed.

## 5. Findings

### M3-01 — `OD-04` has no acceptance criterion (required)

**Evidence class:** Source-derived, complete sweep of §16.

I read every criterion in §16.1–§16.7. The word "explicit" appears once, in `AC-C20`: "The originating assistant message is unchanged unless explicitly edited" — which governs image generation, not change scope. No criterion references preservation across a transition, scope display, promotion, or change history.

§16 exists to make requirements testable. `OD-04` is a resolved owner decision and the only resolved decision with no criterion attached. `OD-01` likewise has none, but it is a technology constraint expressed in §11 and §21 rather than a behavior a test would exercise; `OD-04` describes observable product behavior.

**Consequence if unfixed:** an implementation that silently resets user edits on every scenario selection can pass acceptance.

**Required:** add criteria covering `P2`–`P6`. Suggested placement is §16.3, since transitions are where the conflict occurs, plus one criterion in §16.4 for scope of derived entries.

### M3-02 — The transition transaction has no confirmation gate or decline path (required)

**Evidence class:** Source-derived, §7.1 and §7.7 read against `AC-C23`.

§7.1 states that selecting a new scenario "must, as one coordinated operation" record, persist, insert the first message, display, record the transition, persist atomically, and — step 8 — "make the transition available to the next model interaction." No step consults `OD-03`, and no step defines what happens when the user declines.

§7.7 separately requires confirmation before "replacing or discarding an explicit user change." The two sections are individually correct but their ordering is contradictory when a conflict exists: step 8 makes the transition reachable by the model before §7.7's confirmation is described as occurring.

`AC-C23` sharpens this by requiring the transition be recorded and displayed "immediately" on selection.

Note that `AC-C29` already handles the analogous problem for in-flight generation ("A transition cannot silently apply to an in-flight generation from the old scene"), which is exactly the shape of gate `OD-04` needs but does not have.

**Required:** either add an explicit confirmation step to §7.1 between steps 1 and 6, or state that the transaction executes only after §7.7 resolves, and define the decline outcome — in particular whether the selection reverts, remains recorded but unapplied, or is discarded. Then align `AC-C23`'s wording so "immediately" does not read as "before confirmation."

### M3-03 — Chat-typed corrections have no defined scope, and memory is character-scoped today (required)

**Evidence class:** Source-derived from §9.3, §11, and the current implementation.

§9.3 requires the system to support "corrections and contradictory information" as a provenance case. It says nothing about scope. §11 records that automatic lore is currently stored in "character-scoped files," and `OD-02` records that character memory persists across chats.

Together these permit a path that `OD-04` forbids in terms: the user types a correction in one conversation, the memory or automatic-lore system records it as a character fact, and it becomes visible in every other conversation with that character — a `P4` violation arrived at without any explicit operation.

`AC-C54` does not help; it constrains privacy, not scope.

**Required:** state in §9 that an entry derived from Chat input is conversation-scoped by default, and that making it character-persistent is the explicit promotion `OD-04` requires. Add the matching criterion.

### M3-04 — Chat-scoped lifecycle has no cross-chat isolation criterion (required)

**Evidence class:** Source-derived, §7.6 read against §16.

§7.6 states the rule twice and states it well: Step Out, return, death, and ghost participation "are chat-scoped lifecycle states" that "must not change the character's status in other chats," and a departing character "remains part of the chat's history."

No criterion tests it. `AC-C61` covers a different mechanism — a nonparticipant's private response must not add them to the participant list — and `AC-C27` covers scenario state persisting across reload.

This matters because §7.6 also says the existing Step Out and Death mechanics are "preserved" and "reconciled," which in practice means retaining code whose current scope behavior is unverified. The requirement is the one most likely to be assumed satisfied by preservation rather than tested.

**Required:** add criteria that death, departure, and ghost participation in one conversation leave the character's state in another conversation unchanged, and that the departed character's history remains visible in the conversation where it happened.

### M3-05 — The change record is required but not tested (partial, required)

**Evidence class:** Source-derived, §1.4 `P6` and §7.2 read against §16.

§7.2 requires the transition record to carry "the change reason and actor." §1.4 requires retaining "the historical record of what changed and why." No criterion requires that record to be written, retrievable, or that a superseded value remain inspectable. `AC-C28` establishes only that old scenario state is not treated as current, which is the opposite direction — it tests that history does not leak into active state, not that active state's history survives.

**Required:** one criterion that a transition, its reason, its actor, and the superseded value survive reload and remain retrievable.

### M3-06 — `AC-C66` breaks the criterion sequence (informational)

`AC-C66` sits inside §16.5 between `AC-C51` and §16.6, so the visible sequence runs `C50`, `C51`, `C66`, `C52`. The ID is correct and stable, and §16 forbids reuse. The only consequence is that coverage tooling or review counting by contiguous range will report `C52`–`C65` out of order. No change required; noting it so a later automated check is not mistaken for a gap.

## 6. What is not wrong

To be explicit about the limits of this audit, since the question was whether a section is contradicted elsewhere:

- §7.7's final paragraph states `OD-04` correctly and completely — preserve at conversation scope, show scope, require explicit promotion, never overwrite silently during a transition.
- §7.2 includes "the change reason and actor" in the transition record, which is the mechanism `P6` needs.
- §15 records `O-11` as resolved by `OD-04` rather than deleting the ID, and adds a non-narrowing sentence. That is the correct treatment and it is the safeguard that keeps `R-09` design work from eroding the policy.
- §17's `A2-08` entry now attributes the resolution correctly and confines the remainder to field taxonomy.
- §12.2's `R-09` question list already carries "Explicit user changes, scenario defaults, and field-level conflict rules" and "Promotion of an explicit change between conversation, character, and scenario scope." Research is therefore scoped to answer this, independent of the criteria gap.
- §10.4's model-input receipt row and `AC-C66` close `M2-01` as claimed.

## 7. Effect on the next stage

No effect on research authorization. `R-09` is already charged with the transition model, explicit-change conflict rules, and promotion mechanics; `M3-01` through `M3-05` change what the acceptance gate must eventually verify, not what research must investigate.

The required changes are documentation-only and can be applied as a controlled revision at any point before architecture approval. They should be closed before that gate, because §16 is the instrument by which `OD-04` becomes enforceable, and an unenforced resolved decision is indistinguishable from an open one at acceptance time.

## 8. Authorization boundary

This addendum authorizes nothing. It is a read-only finding against revision 6. No source, configuration, model, database, runtime, migration, network, Android, or production change follows from it. Architecture approval and acceptance remain not passed. Research for `R-07`–`R-11` requires the owner's explicit authorization.
