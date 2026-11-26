# MEMORY-09 — Formation and Classification

## Native classification, deferred class formation, and the hybrid between them

**Status:** **design investigation.** It does not choose a direction, amend L-2, amend Decision A, authorize code or
change the implementation plan. The owner has rejected **Option A** (author-supplied class) as the solution. Existing
L-2 remains intact for now; the relationship between formation and classification is a **new design question**.

| Field | Value |
|---|---|
| Date | **2026-11-21** (local, CST) |
| Owner direction | autonomous capture; no DeepSeek; no manual review; no interruption; no silent guessing |
| Directions investigated | **D** native deterministic classification; **E** formation without resolved class; **hybrid** |
| Finding | D can classify only explicit cases; E preserves capture but needs an L-2/Decision A amendment; hybrid is the only combination that meets both autonomy and no-fabrication |
| Recommendation | **hybrid**, not adopted |
| L-2 | unchanged; E/hybrid would require an explicit amendment |
| Decision A | unchanged; unresolved class is not automatically authorized by the status-slot precedent |
| Implementation | **0** — no code, schema, migration, test, requirement or canonical capture |
| Map | **118** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed**.

---

## 1. Owner direction

**CFD-01 (Observed — Option A rejected).** The owner rejects the class-source decision request's recommendation that
class be author-supplied at creation/review (CSR-15). A system that pauses for a human to categorize every memory is
not the intended autonomous system.

**CFD-02 (Observed — the operating constraints).** Howl should operate without DeepSeek, without manual review of every
captured memory, without interrupting conversation for classifications, and without guessing classifications and
storing them as facts.

**CFD-03 (Observed — the two directions and the hybrid).** The owner asks for investigation of:

- **D** — native, deterministic classification using information already available to the memory system;
- **E** — formation does not wait for classification; the formation entry records class as unresolved and a later
  deterministic operation resolves it;
- a **hybrid** — native classification when evidence is sufficient, explicit unresolved otherwise.

**CFD-04 (Observed — the boundary).** Existing L-2 remains intact. E/hybrid would require revisiting L-2's
`class_baseline` requirement; that must be an explicit decision, not a silent implementation change.

---

## 2. What the current implementation actually supplies

**CFD-05 (Observed — no class field).** `types.MemoryRecord` has `Scope` and legacy `Status`, but no
`KnowledgeClass`; `types.MemoryCandidate` has no class field. `classify.go` assigns scope, owner, title, summary,
content, tags and `pending` status — not class.

**CFD-06 (Observed — capture inputs).** `Service.ProcessMessageForEvents(message string, history []string, identities
[]CharacterIdentity)` receives only the message text, a list of history strings, and identities containing **ID and
Name only** (`CharacterIdentity`). It does not receive a speaker role, a source character ID, world/scenario metadata,
or attribution.

**CFD-07 (Observed — the frontend call).** The browser auto-detection wrapper calls
`window.go.main.App.ProcessMessageForEvents(message, [], [])` after the user sends a message. It passes **no history
and no identities**, and it runs on the **user's outgoing message**, not the character's reply.

**CFD-08 (Observed — no Go binding found).** No `App.ProcessMessageForEvents` method was found in the Go app; the
memory package defines the service method, but the frontend's `App` call is not backed by an App binding in the
inspected code. The backend pipeline is also not called from `app-chat.go` or `chat.Service`. The current
auto-capture path is therefore **not a working character-memory capture path** in the inspected build.

**CFD-09 (Observed — event detection).** `DetectEvents` is keyword-based. `EventDetectionResult` carries `EventType`,
`Confidence`, `Evidence` and `Context`; it carries no source character, attribution or class. The example “I grew up
in the northern mountains” has no matching keyword in the inspected map; it would not be detected unless a custom
trigger fired.

**CFD-10 (Observed — metadata exists elsewhere).** The app has `world`, `scenario` and `characters` entities, and chat
history carries `Sender` and `CharacterID`. That information is available at the app/chat layer, but **not passed into
the memory service**. Using it for native classification requires plumbing, not a new model.

**CFD-11 (Source-derived — ambient is authored-only).** K-R8: “Ambient applies only to authored world/scenario
definitions.” Auto-captured character or user content is not an authored world/scenario definition. Auto-captured
content therefore cannot be `ambient` merely because it concerns the world.

**CFD-12 (Source-derived — class is authored, not popularity/age).** I13: “Knowledge class is authored, not inferred
from popularity or age.” A native deterministic classifier is an **inference rule**. Whether an explicit,
owner-approved native rule counts as “authored” for I13 is a semantic question this record does not settle.

**CFD-13 (Calculated — the input gap matters).** The owner's list of available inputs — source character, conversation
context, content, existing world/character metadata, vocabulary — describes what the app **could** supply. The memory
capture API currently supplies only content, shallow history strings, and optional `{ID, Name}` identities. A native
classifier cannot use what is not passed.

---

## 3. The four operations

**CFD-14 (Calculated — separate operations).**

| Operation | What it does | Must happen now? |
|---|---|---|
| **capture** | detect a memory-worthy event / create the candidate or record content | no classification needed |
| **formation** | write the lifecycle baseline: scope, class component, status slot | yes, if the memory is formed |
| **classification** | determine a class from evidence | only when sufficient evidence exists |
| **class change** | change a class after it has been established | later, if warranted |

**CFD-15 (Calculated — the distinction the owner drew).** Capture and formation need not wait for classification.
Classification can be a separate deterministic operation. A later class **change** is distinct from the initial
classification: the first classification completes the baseline; subsequent changes are P-1 changes.

**CFD-16 (Calculated — why this matters).** If formation requires a resolved class, capture waits for classification,
and if classification needs an external model or a human, the autonomous system stops. If formation can carry an
explicit unresolved class, capture proceeds; classification can happen later without fabricating a value.

---

## 4. Option D — native, deterministic classification

**CFD-17 (Calculated — the form).** A pure Go function, no network, no model, no clock, no randomness, classifies a
memory from explicit inputs and returns either a vocabulary value or an explicit **unresolved** result with a reason.
It uses a versioned rule set; the same inputs and rule version produce the same output.

**CFD-18 (Calculated — candidate evidence).**

| Input | Available now? | Safe use |
|---|---|---|
| event type | yes | a signal only; event-type→class requires an approved rule |
| content / source message | yes | explicit markers and structure |
| history / context | partial (strings) | limited; no roles or attribution |
| owner ID / kind | inferred | can exclude `ambient` for non-authored content; not a class |
| source character | **no** (not passed) | would require plumbing |
| world/scenario metadata | **no** (not passed) | would require plumbing |
| attribution (observed/told/believed/asserted) | **no** | would require a new field |
| class vocabulary | docs only | add as a code enum |

**CFD-19 (Calculated — the strongest deterministic exclusion).** Auto-captured content cannot be `ambient`: K-R8
reserves ambient for authored world/scenario definitions, and this path has no authored asset source. That removes one
of four classes without guessing.

**CFD-20 (Calculated — the remaining distinctions are semantic).** For auto-captured character content, the plausible
classes are `experiential` and `character_specific`; `restricted` is used for authored secrets in the worked cases and
its application to a character's private knowledge is not established. Distinguishing an experience from a
belief/preference requires attribution or content semantics that the current inputs do not reliably carry.

**CFD-21 (Calculated — candidate rule families, not approved).**

| Signal | Candidate outcome | Caveat |
|---|---|---|
| source is an authored world/scenario definition | `ambient` / authored-restricted | not available in auto-capture |
| explicit first-person preference, emotion, aesthetic, sensory, social or promise | `character_specific` | event type is a proxy; conflicts possible |
| explicit first-person or witnessed event (death, betrayal, romance, discovery, conflict, change, location, item, first, shared) | `experiential` | a claim about the world may be belief, not experience |
| explicit secrecy/visibility marker with authored status | `restricted` | character secrets may be `character_specific`, not authored-restricted |
| no explicit signal, or conflicting signals | **unresolved** | no default |

No row is approved. The event-type rows are exactly the kind of mapping the prior class-source request said needs
owner approval.

**CFD-22 (Calculated — D must abstain).** If D always returns a class, it fabricates when evidence is insufficient.
The owner explicitly forbids silent guessing. D therefore needs an unresolved outcome; D with abstention is already
the hybrid direction.

**CFD-23 (Calculated — D alone is not sufficient).** D can classify explicit cases; it cannot resolve the semantic
cases without source, attribution or metadata. If the implementation adds those inputs, D becomes stronger, but the
design still needs an outcome for residual uncertainty. D alone cannot guarantee a class for every memory.

**CFD-24 (Calculated — D's rule set is an authoring mechanism).** The native rule set must be explicit, versioned and
reviewable. It must not use scope as class, popularity, age, repetition, or an unapproved event-type default. Whether
this satisfies I13's “authored” language is the owner's semantic decision (CFD-12).

**CFD-25 (Unestablished — accuracy).** No corpus or test establishes how often D can resolve a class correctly. The
0/144 experiment state remains; no accuracy claim is made. A native classifier is testable with unit cases, but its
coverage is not established.

---

## 5. Option E — formation does not wait for classification

**CFD-26 (Calculated — the form).** The formation entry records `class_state: unresolved` (or an equivalent tagged
slot) instead of a class value. The memory is captured and formed immediately. A later deterministic operation
classifies it and records the result with provenance; subsequent class changes are P-1 changes.

**CFD-27 (Calculated — representation variants).**

| Variant | Form | Assessment |
|---|---|---|
| **E1** | formation carries a class slot; an append-only `class_resolution` entry supplies the first value | **preferred if E is chosen**; keeps the stream append-only and distinguishes completion from change |
| **E2** | mutate the formation entry from unresolved to resolved | rejected: the lifecycle stream is append-only; mutation hides the history |
| **E3** | record the first classification as a P-1 `class_change` version 1 with no baseline | rejected: a first value is not a change; it would confuse baseline and change |
| **E4** | keep class only on the record, not in the stream | rejected: violates Decision A's reconstruction requirement |

**CFD-28 (Calculated — E1's projection).** If the formation slot is unresolved and no resolution exists, class is
unresolved. If a resolution exists, its value is the baseline; class changes version from 1 after it. If the formation
slot is resolved at formation, it is the baseline directly. No cross-dimension or global order is created.

**CFD-29 (Calculated — E1 changes L-2's shape).** L-2 is one formation entry with `scope_baseline`, `class_baseline`
and `status_slot`. E1 changes the class component to a tagged slot and adds a class-resolution entry. That is a
logical amendment to L-2, not an implementation detail. It must be decided explicitly.

**CFD-30 (Calculated — E1 changes Decision A's initial-class obligation).** Decision A says the initial
`knowledge_class` must be reconstructible from the stream. E1 reconstructs **“unresolved”**, not a class value. The
owner must decide whether an explicit unresolved initial state satisfies Decision A, or whether Decision A must be
amended. The status-slot precedent does not automatically transfer: N-b authorized an unresolved **status** slot;
class was not included.

**CFD-31 (Calculated — class resolution is not a class change).** The first classification completes the baseline;
it is not a P-1 `class_changed`. P-1 changes remain post-baseline changes. This distinction must be explicit, or the
system will treat the initial classification as a change and lose the baseline.

**CFD-32 (Calculated — usability is qualified).** Scope, ownership, content and status can be used immediately. Class
is used for audience/eligibility; an unresolved class cannot be treated as `ambient`, `restricted`, `experiential` or
`character_specific`. Class-sensitive retrieval must fail closed or explicitly show “unresolved”; it must not silently
use a default. “Immediately usable” therefore means usable for non-class purposes, not for class-gated audience.

**CFD-33 (Calculated — E alone is not autonomous).** If E defers class indefinitely without a classifier, the system
never classifies; if it defers to manual review, it violates the owner's autonomy constraint. E needs a deterministic
later classifier to be autonomous; that is E + D.

---

## 6. The hybrid

**CFD-34 (Calculated — the hybrid form).** At formation, run the native deterministic classifier. If it returns a
class with sufficient evidence, write the resolved class baseline. If it returns unresolved, write the class slot as
unresolved and continue capturing. A later deterministic operation retries unresolved memories; when it resolves, it
writes a class-resolution entry with provenance. Subsequent changes are P-1 class changes.

**CFD-35 (Calculated — why hybrid meets both constraints).** It preserves automatic capture (no manual review, no
external model, no interruption) and avoids silent fabrication (explicit unresolved, with a reason). D alone cannot
cover uncertain cases; E alone cannot classify without a classifier. The hybrid is the composition the owner
described as possibly better.

**CFD-36 (Calculated — hybrid still requires explicit authorization).** It requires:

1. an L-2 amendment for the class slot and class-resolution entry;
2. a Decision A interpretation/amendment allowing an unresolved initial class;
3. owner approval of the native classifier's rule mechanism, evidence thresholds and abstention behavior;
4. a defined retrieval behavior for unresolved class (fail closed; no class-gated audience);
5. tests for resolved, unresolved, resolution, change, no-default and no-cross-dimension ordering.

**CFD-37 (Calculated — hybrid does not decide authority).** Who may resolve a class, whether a background operation
may do so, and whether a later resolution is automatic remain authority questions. The classifier can be deterministic
without deciding who is allowed to run it.

---

## 7. L-2 and Decision A impact

**CFD-38 (Observed — L-2 is intact).** L-2 remains the adopted logical representation: one formation entry with
`scope_baseline`, `class_baseline`, `status_slot`. This record does not edit it.

**CFD-39 (Calculated — the exact clauses E/hybrid would touch).**

| Clause | Current meaning | E/hybrid implication |
|---|---|---|
| L2D-01/07 | one formation entry | the entry carries a class slot; a resolution entry may follow |
| L2D-04 | independent local precedence | the resolution is class-local; no cross-dimension order |
| L2D-08 | `class_baseline` carries initial class | it carries a resolved value or an unresolved obligation |
| L2D-10 | class reconstruction reads baseline then changes | it reads slot, then resolution, then changes |
| L2D-14 | record class is a checkable projection | unresolved is not a class; the check must distinguish |
| DAO-05 | initial class reconstructible | unresolved is reconstructible as unresolved, not as a class |

**CFD-40 (Calculated — no silent change).** E/hybrid cannot be implemented under the current L-2/Decision A text. The
owner must amend them, or choose a different direction. The implementation plan's Phase 0 stop remains in force; it
is now a broader design stop, not a request for an external class source.

---

## 8. Example: “I grew up in the northern mountains.”

**CFD-41 (Observed — the example is not currently captured).** The inspected keyword map does not contain “grew up”;
the frontend runs detection on the user's outgoing message with empty identities. A character saying this would not
reach the current pipeline as a character memory.

**CFD-42 (Calculated — even if captured, the class is ambiguous).**

| Question | Available answer |
|---|---|
| who said it? | not passed; identities are empty in the frontend call |
| is it an experience or a belief? | not distinguishable without attribution/context |
| is it about the character or the world? | content suggests personal history, but subject/predicate are not modeled |
| is it authored world/scenario knowledge? | no; so not `ambient` |
| is it a secret? | no explicit marker; not enough for `restricted` |

**CFD-43 (Calculated — correct behavior).** A native classifier should return **unresolved** for this example unless an
owner-approved rule explicitly covers first-person background statements. It should not default to `experiential` or
`character_specific` merely because the sentence is first-person. This is exactly the case E/hybrid exists to handle.

---

## 9. Comparison

**CFD-44 (Calculated).**

| Criterion | D native only | E deferred only | Hybrid |
|---|---|---|---|
| automatic capture | yes | yes | yes |
| no DeepSeek | yes | yes | yes |
| no manual review | yes | yes, if later classifier exists | yes |
| no interruption | yes | yes | yes |
| no silent guessing | only if it abstains | yes | yes |
| class available at formation | for explicit cases | no | for explicit cases |
| handles uncertain cases | needs E | needs D | yes |
| L-2 amendment | no | **yes** | **yes** |
| Decision A amendment | no | **yes** | **yes** |
| class-sensitive retrieval | class when resolved | blocked until resolved | class when resolved; fail closed otherwise |
| implementation risk | medium; rule quality | medium; representation | higher; both |

**CFD-45 (Calculated — D and E are not rivals).** D without E has no honest answer for insufficient evidence; E
without D has no autonomous classifier. The hybrid is the only combination that satisfies the owner's constraints,
but it is also the largest change.

---

## 10. Recommendation

**CFD-46 (Calculated — recommend hybrid, not adopted).** The owner should investigate a hybrid: a native deterministic
classifier at formation, an explicit unresolved class slot when evidence is insufficient, and a later deterministic
resolution entry with provenance. This preserves immersion and autonomy while refusing to store a guess as a class.

**CFD-47 (Calculated — what must be approved first).** The recommendation is not implementable until the owner decides:

1. whether a native rule-based classifier satisfies I13's “authored” requirement;
2. whether L-2 may carry an unresolved class slot and a class-resolution entry;
3. whether Decision A permits an initial class to be reconstructible as unresolved;
4. how unresolved class behaves in class-sensitive retrieval;
5. how the native rule set is versioned, reviewed and changed.

**CFD-48 (Calculated — the fallback).** If the owner rejects an unresolved class, the only alternatives are a native
classifier that always resolves (which risks fabrication) or a class supplied externally (Option A, already rejected).
The record does not choose among them; it states the tension.

---

## 11. Decision requested

**CFD-49 (Observed — the owner's choice).**

| Option | Meaning |
|---|---|
| **D** | native deterministic classifier at formation; must define what happens when evidence is insufficient |
| **E** | class deferred; formation carries unresolved; later classifier/resolution |
| **Hybrid** | native at formation when sufficient; unresolved otherwise; later deterministic resolution |
| **Other** | specify the mechanism; it must satisfy the no-default/no-guess constraints |

**CFD-50 (Observed — no implementation follows).** This record does not select a direction. The owner's choice is the
next decision. If hybrid is chosen, the L-2/Decision A amendment must be recorded separately before the implementation
plan is revised.

---

## 12. What this record does not do

**CFD-51 (Observed).**

| Untouched | State |
|---|---|
| L-2 | **intact**; no amendment |
| Decision A | **intact**; no reinterpretation |
| P-1 / S-3 | closed; not reopened |
| Decision B | untouched |
| N-b / N-c / authority | untouched |
| schema / code / migration / tests | **0** |
| requirements / map | **0**; **118** |
| canonical capture | **0** |
| implementation plan | **not revised**; Phase 0 remains a stop |
| Howl UX / state machine | **0** |

---

## 13. Traceability

**CFD-52 (Observed).**

| What | Where |
|---|---|
| Owner direction | 2026-11-21; CFD-01…04 |
| Class-source request (A rejected) | CSR-15; CFD-01 |
| No class field | `types/memory.go`; CFD-05 |
| Capture inputs | `Service.ProcessMessageForEvents`; `CharacterIdentity`; CFD-06 |
| Frontend call | `frontend/dist/js/memory-candidates-menu.js`; CFD-07 |
| No App binding | Go app inspection; CFD-08 |
| Keyword detector | `detect.go`; CFD-09 |
| Ambient authored-only | K-R8; CFD-11 |
| Class authored | I13; CFD-12 |
| L-2 shape | L2D-01/04/07/08/10/14/20; CFD-29/39 |
| Decision A | DAO-05; CFD-30/39 |
| P-1 class changes | P1D-04…11; CFD-31 |
| Status-slot precedent | IAA-03; CFD-30 |
| Map | **118** |

---

## 14. Record status

**CFD-53 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — map **118** |
| Direction selected | **0** — investigation only |
| L-2 / Decision A amended | **0** |
| Implementation started | **0** |
| Tests executed | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **design investigation** — no canonical capture |

**End of record.** Option A is rejected (CFD-01). Native classification can exclude `ambient` for auto-captured
content and can resolve explicit cases, but the current capture API lacks the speaker, attribution and metadata needed
for the semantic distinctions, and it must abstain when evidence is insufficient (CFD-17…CFD-25). Deferred class
formation preserves capture but changes L-2's one-entry `class_baseline` into a class slot plus a later resolution,
and changes Decision A's initial-class obligation; neither may be done silently (CFD-26…CFD-33). The hybrid is the
only composition that meets autonomy and no-fabrication, but it requires explicit L-2/Decision A amendments and
approval of the native classifier mechanism (CFD-34…CFD-47). **No direction is chosen and no implementation follows**
(CFD-50/53).
