# MEMORY-09 — Implementation review R2: focused classification

**Status:** **SOURCE INSPECTION ONLY. 0 tests executed. 0 probes run. No runtime behaviour claimed.**
Recorded **2026-11-17, R2** — second review, requested by the owner after R1.

**Scope:** the five findings the owner raised, plus the three absent requirements. Each is classified
as **A** implementable under the ratified architecture · **B** requires owner clarification ·
**C** cannot be implemented without an architecture change.

**Constraints honoured:** no new rules, fields or requirements · no probing · no test execution ·
no change to any ratified artifact.

Line anchors below are from the current `eval/mem1/m09/` build (post-R2 corrections). Function names
are the stable reference; line numbers shift with edits.

| Anchor | Location |
|---|---|
| `Engine` class | `engine.py:29` |
| `claim` / `acquire` | `engine.py:62` / `82` |
| `declare_belief` | `engine.py:102` |
| `promote` | `engine.py:131` |
| `authoritative` | `engine.py:191` |
| `disclose` | `engine.py:203` |
| `narration_report` | `engine.py:242` |
| `retract` / `correct` | `engine.py:255` / `268` |
| `supersede_authority` / `revoke_promotion` | `engine.py:278` / `292` |
| `_require_human` | `engine.py:301` |
| `declare_available` / `withdraw_declaration` | `engine.py:315` / `355` |
| `snapshot` | `engine.py:385` |
| `perspective` / `world` / `dispatch` | `retrieval.py:38` / `57` / `69` |

---

## 1. Rule 41 — false conflicts

**Execution path.** `promote` (`engine.py:131`) collects every record with the same `record_class`,
`domain` and `target` whose `status` is `ACTIVE`, regardless of content. When `supersede=False` it
leaves them all active and adds the new record. `authoritative` (`engine.py:191`) then filters the same
way and returns `status="conflict"` whenever more than one record survives. `retrieval.world`
(`retrieval.py:57`) reports that as `equal_authority_conflict`.

**What that means.** The conflict predicate is *target equality*, not *incompatibility*. Consequences:

- two promotions of identical text produce a conflict;
- two genuinely compatible statements about one target produce a conflict;
- nothing ever compares content, because rule 34 and v9 §9 forbid the system judging relevance,
  persuasion or truth.

**Why it exists.** "Incompatible" is never given an operational meaning in v9. r32's materiality
predicate is itself open — recorded in the O-6a register as **D-2** (D6 versus §10.1 rev 2, still
`[P]`). The implementation substituted the only predicate available without adjudication.

**Classification: B — requires owner clarification.** The unresolved question is *who establishes
incompatibility, and how is it represented*: a declared incompatibility by the promoting actor, a
declared applicability overlap, or something else. Any comparator the engine invents would be the
semantic adjudication rule 34 forbids, so there is no implementation-only answer.

---

## 2. Rule 56 / D10 — belief versus truth

**Execution path.** R1 shipped this branch inside `promote` (`engine.py:131`):

```python
if cited.record_class == "belief" and not supersede:
    pass
```

It tested for a cited belief and did nothing. **R2 removed it** — dead code that reads as a guard is
worse than a documented gap.

After removal, the path is: basis validated (`authority.validate_basis`) → citations resolved → a
truth record is created with `text` supplied by the caller → `authoritative` (`engine.py:191`) returns
it as settled for `(domain, target)` → `retrieval.world` reports `settled`.

**What that means.** The distinction between *"Sable believes P"* and *"P is true"* exists only in the
promoting actor's prose. The engine records basis (including a cited belief) in
`promotions[pid]["basis"]`, but never consults it when answering truth queries. If an actor promotes a
belief citation with `target=P` and `text="P"`, the record answers as truth about P.

**Why it cannot be fixed in place.** The engine would need a structural carrier: a way to say *this
record's subject is a character's belief, not the world*. That is a record class, a subject-kind field,
or an applicability dimension that does not exist in v9. Adding one is precisely what Gate 2 prohibits.

**Classification: C — cannot be implemented without an architecture change.**

---

## 3. Rule 39 / K25 — ambient authorship

**Execution path.** There is none. `claim` (`engine.py:62`) records contributions; `declare_available`
(`engine.py:314`) creates coverage accounts; neither authors ambient content under a shared regime, and
`promote` (`engine.py:131`) is the only path to authority. Nothing in the package references Q7.

**What that means.** "Ambient authorship does not promote" is currently true **vacuously** — because
ambient authorship does not exist. The separation is not demonstrated; it is absent along with the
feature.

**Why it cannot be closed by implementation alone.** v9 rule 39 names "Q7 ambient authorship" but the
ratified set never defines it: not its inputs, not its outputs, not whether it creates records at all,
not how it shares the per-domain regime with promotion. Implementing it would mean inventing that
contract.

**Classification: B — requires owner clarification** (or a decision to keep it out of scope, which is
equally a clarification).

---

## 4. Rule 45 / K27 — the transfer stub

**Execution path, R1.** `engine.py` exported `EXPOSES_TRANSFER = False`, listed `"transfer"` in
`__all__`, and defined `def transfer(*_args, **_kwargs): raise OperationError(...)`.

**The inconsistency the owner identified was real.** R1 said both *"exists only to raise"* and *"no
transfer operation is exposed"*. A module-level callable in `__all__` **is** an exposed interface; that
it always raises is a behaviour, not an absence. An integrator could import it, type-check against it,
or wrap it.

**Correction (R2).** The constant, the function and both `__all__` entries were deleted. The only
remaining occurrences of the word are comments explaining the absence. `grep -rn transfer
eval/mem1/m09/` now returns comments and stale `.pyc` files only.

**Classification: A — implementable, and done.** Removing capability-shaped surface is not an
architecture change; it is the removal of a mistake.

---

## 5. F7 — replay determinism

**Execution path.** `snapshot` (`engine.py:385`) serialises current state. There is **no operation log
and no replay entry point**, so F7 has nothing to run against.

**Distinction the review must keep.** *Formation determinism* — identical ordered inputs produce
identical records — is a property of the write paths. *Replay determinism* — reconstructing state from
a recorded history reproduces the same snapshot — is a property of a replay facility. The engine has
neither demonstrated; having no model call is not evidence of either.

**One concrete order-dependence found and fixed.** `retrieval.perspective` sorted candidates by score
alone. Python's sort is stable, so equal scores were resolved by dictionary insertion order — the same
logical state could answer differently depending on creation order. **R2 replaced the sort key with
`(-score, record_id)`**, a canonical tie-break. This removes an order-dependence; it does not
demonstrate determinism.

**Classification: A — implementable**, in two parts: (i) canonical ordering, done; (ii) a replay
function over a recorded operation list, not yet written. Neither closes F7 by itself, and **no claim
of determinism is made** until a replay exists and is executed — which testing alone can establish.

---

## 6. The three absent requirements

| Req | Missing capability | Classification | Reasoning |
|---|---|---|---|
| **B3** | Event/definition distinction | **C** | v9 carries no event-versus-definition marker on a record. Enforcing "an event-kind memory cannot be declared" needs a representation that does not exist. **Remains unresolved by the owner's instruction** — no field invented, no heuristic substituted |
| **C5** | Write-time evidence-kind validation | **B** | The *structure* is implementable: `acquire` (`engine.py:82`) can refuse an acquisition whose evidence is absent. What is missing is the **per-`via` evidence requirement table** — what counts as evidence for `witnessed`, `told`, `read`, `inferred`, `narrator_grant`. That mapping is not ratified, and inventing it would be architecture by accident |
| **K25** | Ambient authorship separation | **B** | Same finding as §3: Q7 is unspecified (see above) |

---

## 7. Classification summary

| # | Item | Rules / reqs | Class |
|---|---|---|---|
| 1 | Conflict predicate uses target equality | 41, J13, J14 | **B** |
| 2 | Belief-versus-truth distinction unrepresented | 56, D10 | **C** |
| 3 | Ambient authorship absent | 39, K25 | **B** |
| 4 | Transfer exposed as a callable stub | 45, K27 | **A — done** |
| 5 | No replay facility; tie-break order-dependence | F7 | **A — partly done** |
| 6 | Event/definition distinction missing | B3 | **C** |
| 7 | No evidence-kind validation | C5 | **B** |

**A: 2 · B: 3 · C: 2.** No item was closed by assumption.

---

## 8. Changes made in R2 (three, all corrections)

| Change | File | Why it is not an architecture change |
|---|---|---|
| Deleted `transfer()` and `EXPOSES_TRANSFER` | `engine.py`, `__init__.py` | Removes an interface; adds nothing |
| Deleted the dead rule-56 branch | `engine.py:131` | Removes code that did nothing |
| Canonical tie-break `(-score, record_id)` | `retrieval.py:38` | Removes order-dependence |

A fourth pass removed one more dead branch in `promote`: `if active and not supersede: pass`. It had
the same defect as the rule-56 branch — a placeholder that reads as handling while doing nothing. The
state it deferred to (`authoritative` reporting a conflict) is unchanged.

No rule, field, requirement or probe was added. The 57 rules and 118 requirements are unchanged.

## 9. What is not claimed

- No test, probe, import-execution or suite was run. `py_compile` passed; that is syntax.
- R1's verdicts stand except where R2 refines them: rule 56 moves from *partial* to **cannot be
  completed without an architecture change**; F7 gains a fixed order-dependence and still has no replay.
- "Implemented" continues to mean *a code path appears to satisfy the requirement*, not *the
  requirement is met*.

## 10. Clarification requests (no answers assumed)

| # | Question | Blocks |
|---|---|---|
| **Q-R2-1** | How is *incompatibility* established and represented — declared by the promoting actor, derived from declared applicability overlap, or otherwise? | Rule 41, J13, J14 |
| **Q-R2-2** | Is a structural carrier for "record is about a character's belief" authorised as an architecture change? | Rule 56, D10 |
| **Q-R2-3** | What evidence does each `via` kind require at write time? | C5, rule 50 |
| **Q-R2-4** | Is Q7 ambient authorship in scope, deferred, or out of scope? Either answer closes K25 legitimately | Rule 39, K25 |
