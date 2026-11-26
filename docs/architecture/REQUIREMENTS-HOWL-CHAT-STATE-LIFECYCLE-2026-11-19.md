# HOWL Chat — Requirements: Chat-State Lifecycle (HSL)

**Status:** requirement record only. Records an owner-stated responsibility split (2026-11-19). **Not an
amendment** to `INTEGRATION-AND-ACCEPTANCE-PLAN-2026-11-11.md`, **not** architecture approval, **not**
frozen, **not** implementation authorization, and **not** a MEMORY-09 rule. No implementation exists for
any requirement below.

| Field | Value |
|---|---|
| Created | **2026-11-19** |
| Authority | The owner's instruction of this session: record chat-state lifecycle as a Howl requirement, keep it out of the MEMORY-09 experiment, defer the integration contract |
| Applies to | The Howl chat layer (state machine, branching, regeneration, retention, deletion) and the Howl↔MEMORY-09 boundary |
| Does not apply to | MEMORY-09 architecture, its experiment plan, its test corpus, its probes, its conformance rows, or its canonical records |
| Cross-references | `docs/research/MEMORY-09-PHASE-1A-1-RETRACTION-AND-KNOWLEDGE-PERSISTENCE-2026-11-19.md` §13; `docs/research/MEMORY-09-PHASE-1A-1A-CHAT-STATE-LIFECYCLE-2026-11-19.md` (findings RPS-01–RPS-12, questions RQS-01–RQS-06) |
| Current implementation | None. §5 records today's behaviour for information only; it authorizes no change. |
| Tests executed | **0.** Behavioural statements in §5 are static readings of source, not runs. |

---

## 1. Purpose

A conversation can stop being current for reasons that have nothing to do with what anyone asserted or
knows: the user picks another alternative, regenerates a reply, reverts to an earlier state, or deletes
messages. Those are **conversation-state lifecycle** events, and they belong to Howl's chat layer.

This note records that ownership as a requirement, so that:

- Howl's chat-state responsibilities are written down where Howl's architecture work will find them;
- MEMORY-09 is not asked to infer, model, or experimentally validate a state machine it does not own;
- the memory consequences of a state event are named now and specified later, at the boundary, rather
  than being discovered during integration.

It creates no experimental variable, no probe, no conformance row, and no MEMORY-09 status value.

---

## 2. The boundary this note draws

The near-miss this note exists to prevent: pulling a Howl responsibility into MEMORY-09's experimental
scope. A state machine is application behaviour — it is the user pressing a button, not a character
acquiring or losing knowledge — and an experiment that varied it would be measuring Howl, not memory.

| Responsibility | Owner |
|---|---|
| Chat state, branching, regeneration, deletion | **Howl** |
| Determining which state is currently authoritative | **Howl** |
| Retaining abandoned states and their material | **Howl** |
| Memory representation and retrieval within established boundaries | **MEMORY-09** |
| Determining which memories depend exclusively on deleted material | **MEMORY-09** (under its own rules — consequences only; Howl retains the deletion decision) |
| Applying memory lifecycle changes resulting from chat-state events | **Integration contract** between Howl and MEMORY-09 (deferred, §7) |

Two directions of the boundary are load-bearing, and both are requirements (HSL-2, HSL-11):

- MEMORY-09 must **not** infer a transition from message content, message absence, ordering, or turn
  numbers. It receives lifecycle facts or it does nothing.
- Howl must **not** write memory records, edit their status, or decide what a character knows. It
  declares state facts; MEMORY-09 applies them.

And a third line, because HSL-9 is easy to misread in the opposite direction: **MEMORY-09 does not own
deletion.** Determining dependency is a consequence calculation inside MEMORY-09's own domain; it is
not authority over the act. Howl decides what is deleted and which state exists; MEMORY-09 decides
what that means for its records under the contract. It cannot veto a deletion, reinterpret a retained
branch as deleted, or independently conclude that a chat state no longer exists.

---

## 3. Already established in MEMORY-09 — not reopened here

These behaviours exist in MEMORY-09's established design and frozen harness. They are recorded here
only to mark the line: **the state-lifecycle requirement below is not a rediscovery of them, and this
note adds nothing to them.**

| Established behaviour | Where it is established | What it already says |
|---|---|---|
| **World hot-swap** — memories deactivate/reactivate with the active world | `eval/mem1/model.py:102–110` `record_is_active`; `:42` `GATED_OWNER_KINDS = ("world", "scenario")`; `ScopeState` docstring `:55–62` | world-scoped records are eligible only while their world is active; activation changes eligibility only |
| **Scenario hot-swap** — memories deactivate/reactivate with the active scenario, under its parent world | `record_is_active` (scenario active **and** declared parent world active, explicit map); `eval/mem1/retrieve.py:138–140` (`requested_scenario` must match) | one scenario cannot see another's records even under a shared parent |
| **Character public/private context** — memories follow the applicable character context and access boundaries | `eval/mem1/retrieve.py:32–35` `VISIBILITY`; `:130–137` owner/restricted/audience conditions; `assert_selection_owned` | ownership and visibility are separate axes; a foreign record is never a candidate, and selection of one is a veto |

The precedent matters for §4. `ScopeState`'s own docstring already states the semantics the new
requirement needs:

> "State is a declaration, not a record: applying it changes eligibility only. It never creates,
> deletes, supersedes or revokes anything, and it never touches truth status. **Preservation is
> retention; restoration is renewed eligibility.**"

That is exactly the shape HSL-10 asks for, applied to a different trigger. **What is new is the
trigger — a chat-state transition — not the semantics.**

Two things this note deliberately does not do:

- It does not reopen world/scenario/character gating as an experimental requirement, and it does not
  treat any of the three as a variable in the MEMORY-09 experiment.
- It does not audit whether `m09/` implements the gating. `m09/` retrieval has no active-state gate;
  that is a **pre-existing implementation gap recorded elsewhere**, and it is not reopened by this
  note. The citations above are to the harness, which is where the behaviour is established.

---

## 4. The requirement

Normative statements about Howl's chat layer. No format, identifier scheme, storage layout, or API
signature is chosen here (§7 defers that).

**HSL-1 — Howl owns the state machine.**
Chat state identity, branching, alternative selection, regeneration, retention and deletion are Howl
responsibilities. MEMORY-09 neither models the state machine nor decides when a state exists.

**HSL-2 — Authoritative notification.**
Every operation that changes which state is current, or changes the set of retained states, emits an
explicit lifecycle notification to MEMORY-09, naming the operation, the state(s) affected and the state
that is now current. Absence of notification means no transition; MEMORY-09 must never infer one.

**HSL-3 — Stable state identity.**
A state's identifier is stable across save and load, message renumbering, in-place edits, and
regeneration. It is not derived from a message's position in the conversation, and it does not change
when earlier messages are removed.

**HSL-4 — Alternative selection ("Next").**
Producing another response produces a **distinct state**. The state that produced the previous
response is retained with its material; the two are alternatives, and neither replaces the other in
history.

**HSL-5 — Regeneration.**
Replacing a response creates a state whose supporting material is the new response. The replaced
response remains identifiable as material of the state that produced it. Regeneration is not deletion,
and it does not silently overwrite the previous state's material.

**HSL-6 — Revert.**
Reverting selects a retained state as current. The act of reverting creates, deletes and rewrites
nothing; it changes which retained state is authoritative.

**HSL-7 — Deletion is the only destructive operation.**
No operation other than explicit deletion removes conversation material or the memory that depends on
it. Deletion is explicit, names what it will remove, requires confirmation, and reports the set it
removed — counts and categories computed from the actual deletion set at confirmation time, never
estimated or stale.

**HSL-8 — Deletion names its set.**
A deletion notification names the material removed and the state(s) it belonged to. Howl does not
decide which memories that implies.

**HSL-9 — Dependency is determined by MEMORY-09; deletion authority is not.**
Given a deletion set, MEMORY-09 determines which records depend on it under its own dependency rules,
and removes **only** records that depend exclusively on the deleted material. A record that also
depends on retained material is retained. What "depends" means is a MEMORY-09 question, not a Howl one.

This is a consequence calculation, not authority over the act. MEMORY-09 may not veto or defer Howl's
deletion, may not treat a retained state as deleted, and may not decide that a chat state no longer
exists. Where it cannot determine dependency, it reports that under the contract (§7 item 7) and the
deletion proceeds as Howl decided; an unresolved dependency is not a reason to invent a retention or to
widen a deletion.

**HSL-10 — Retention across non-destructive transitions.**
A transition under HSL-4, HSL-5 or HSL-6 never removes memory. A memory that is inapplicable in the
current state is not deleted, not forgotten, and not retracted; it remains available to history and to
restoration. (This is `ScopeState`'s "preservation is retention; restoration is renewed eligibility",
with a chat-state trigger.)

**HSL-11 — No unilateral decisions, in either direction.**
MEMORY-09 never decides that a branch has been abandoned, restored or replaced; it applies what Howl
declares. Howl never writes, edits, suppresses or un-suppresses a memory record; it declares state
facts and leaves their memory consequences to MEMORY-09 under the contract.

---

## 5. Current Howl behaviour (informational, authorizes nothing)

Observed at the hashes in the MEMORY-09 1A-1A record §3, except where marked Inferred. Recorded so the
gap between requirement and reality is not discovered later as a surprise; **not** a defect report, not
a work order, and not a statement that these must be fixed now.

| Operation | Today | Requirement it does not meet |
|---|---|---|
| State identity | none — history is one mutable list; message ids are positional (`app-chat.go:405`; `app-prompt.go:327–390`) | HSL-1, HSL-3 |
| Notification | none — no interface exists | HSL-2 |
| Next / alternative | no alternative set; Prev/Next only scrolls (`chat-actions.js:329–348`) despite documentation promising variations (`tutorials-data-chat.js:117`) | HSL-4 |
| Regenerate | removes the reply from the DOM and re-streams without syncing history (`chat-actions.js:273–304`); the removed reply may survive in the append log and reload beside its replacement (Inferred) | HSL-5, HSL-6 |
| Revert | no revert action exists | HSL-6 |
| Undo / delete / edit | prune and atomically rewrite the history log (`app-prompt.go:389`; `persistence.go:631–661`); the previous state is unrecoverable | HSL-7, HSL-10 |
| Deletion reporting | none — the file rewrite is silent | HSL-7, HSL-8 |

---

## 6. What this note does not authorize

- No implementation, migration, build, launch, or runtime change.
- No change to Howl's current behaviour, including the destructive undo. §5 is description.
- No architecture approval for a branching chat model, and no choice among retention designs.
- No integration format, protocol, schema or identifier scheme (§7).
- No MEMORY-09 change of any kind.

If the owner later wants HSL folded into `INTEGRATION-AND-ACCEPTANCE-PLAN-2026-11-11.md`, that is a
separate amendment act with its own revision and register entry. This note does not presume it.

---

## 7. The integration contract — deferred, not specified

To be addressed at Howl's foundational architecture work, as the owner directed. It must eventually
carry, at minimum:

1. the identity of the current state, and the ordered history of transitions;
2. the operation kind for each transition (alternative, regeneration, revert, deletion);
3. the material each transition affects, including an explicit deletion set;
4. the durability guarantee for retained states (what Howl promises to keep, and for how long);
5. refusal and error semantics — e.g. a transition naming an unknown or unreachable state;
6. ordering and idempotency of notifications, including replay after restart;
7. what MEMORY-09 returns for a deletion set (its dependency determination under HSL-9), including an
   explicit failure result when dependency cannot be determined — Howl's deletion decision does not
   depend on that result and is not conditional on it;
8. what MEMORY-09 may assume about a state it has never been told about.

None of these is decided here. Item 4 in particular is the one that binds Howl: MEMORY-09 cannot
satisfy HSL-10 over material Howl has already discarded.

---

## 8. Explicit exclusions

This note and the requirements in §4 are **not**:

- an addition to `MEMORY-09-EXPERIMENT-PLAN-2026-11-15.md` or its amendments;
- a corpus, fixture, probe or scoring change;
- a row in the 118-requirement conformance map;
- a MEMORY-09 rule, invariant, record class, field, or status value;
- an input to PH1A-12 (experimental subject), which remains downstream of Q-A and Q-D in the research
  record and is unaffected by this note.

The research record's findings about what a state-aware memory representation would require
(`RPS-03`, `RPS-04`, `RPS-11`) remain research observations inside MEMORY-09's boundary. This note does
not convert them into requirements.

---

## 9. Open items handed to Howl's foundational architecture work

| Item | Question | Where it came from |
|---|---|---|
| **HSL-O1** | What is the unit of state: a single message alternative, a branch, or a lineage? The three have different retention and query consequences. | RQS-02 |
| **HSL-O2** | Must abandoned states remain queryable, and by whom? Retention for restoration and retention for inspection are different promises. | RQS-03 |
| **HSL-O3** | What must Howl retain for how long, and what happens to memory when Howl's retention obligation ends? | RQS-06 |
| **HSL-O4** | Does a state transition affect a character's ownership of a record, or only its applicability? | RQS-05 |
| **HSL-O5** | Where applicability lives — on the record, in a projection, or in partitions — is a joint contract decision, not a Howl-internal one. | RQS-01, RS-1–RS-3 |

**HSL-O5 is the one item that cannot be closed by Howl alone**, and it is why the contract is a
contract: Howl can guarantee a durable state identity and an ordered transition history, but whether
MEMORY-09 tags records, projects eligibility, or partitions its store remains MEMORY-09's side of the
boundary.

One MEMORY-09 question stays wholly inside MEMORY-09 and is **not** handed over: whether a record that
outlives its branch — a promoted `world_truth`/`scenario_truth`, or standing character knowledge formed
before the branch point — remains authoritative after the branch is abandoned (research RQS-04). That is
a question about MEMORY-09's own authority semantics, and it is unaffected by who owns the state
machine.

---

## 10. What this note does not decide, and its limits

No state model is chosen. No retention design, identifier scheme, storage layout or protocol is
specified. No MEMORY-09 rule or status is created. The §4 requirements are recorded as the owner's
stated division of responsibility; they are not an acceptance criterion, not a frozen contract, and not
an implementation authorization.

**Limits.** §5's behavioural statements are static readings of the working tree at the hashes in the
1A-1A record; nothing was executed and no interface was rendered. The regenerate-duplication cell is
Inferred from the append-log path and is not demonstrated. The §3 citations are to the frozen harness,
not to `m09/`, which does not implement active-state gating; that gap is acknowledged, not measured
here.

---

**End of requirement record.** No code, test, corpus, MEMORY-09 artifact, or canonical document was
modified. The next step is Howl's foundational architecture work, where §7's contract is specified.
