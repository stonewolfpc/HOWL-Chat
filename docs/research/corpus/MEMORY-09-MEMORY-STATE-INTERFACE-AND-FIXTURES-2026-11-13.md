# `MEMORY-09` — Memory-State Interface and Synthetic Fixtures

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-MSI-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **INTERFACE + FIXTURE SPEC — nothing implemented, no test run.**
**Authority:** owner direction of 2026-11-13 (Python-first follow-on to `WO-ARCH-06`).
**Related:** T-04, T-06…T-11 (row 473).

---

## 1. What the state layer is for

The harness needs a **reproducible, inspectable memory state** that is independent of Howl — so that when a reply is correct, the record can show *whether the state contained the fact*. That is the `C2`/`C3` distinction in practice (row 473, T-11).

## 2. State model

| Entity | Holds | Key fields |
|---|---|---|
| `character` | Per-character knowledge and voice constraints | id, name, known_facts[], voice_notes |
| `persona` | The user's character | id, name, facts[] |
| `world` | Setting facts | id, canon_facts[] |
| `fact` | A single atomic claim | id, text, subject, source, established_at_turn, confidence |
| `item` | Inventory entry | id, name, holder, state (present/consumed/removed) |
| `equipment` | Worn/wielded | id, item_id, slot, worn_since_turn |
| `relationship` | Character↔persona state | id, pair, standing, history_refs[] |
| `event` | Something that happened | id, turn, participants[], summary |

**`fact.source`** records *how* it entered state (`authored`, `extracted_from_turn`, `owner_edit`). This is what makes `C2` checkable later.

## 3. Interface

```
snapshot()              -> state, sha256
get(entity, id)         -> record
put(entity, id, record) -> new sha256
delete(entity, id)      -> new sha256
diff(sha_a, sha_b)      -> changed entities
assert_fact(id)         -> (present, source, established_at_turn)
```

**Rules:**
1. **Every mutation returns a hash.** A test that cannot name the state hash cannot claim what state the model saw.
2. **Snapshots are immutable.** A test holds `sha_before`; the diff is the evidence.
3. **No silent normalisation** — no deduplication, no merging, no fuzzy matching. If two facts conflict, the conflict is stored (`T-10` depends on it).
4. **No deletion without a record.** Removal sets `state=removed` with a turn; it does not erase.

## 4. Synthetic fixtures (F-ST-01…06)

**These are engineering fixtures, not corpus items.** They are *inspired by* the matrix dimensions but they score nothing, belong to no corpus family, and must never be cited as `MEMORY-09` units.

| Fixture | Dimension | Canonical content | Adversarial probes |
|---|---|---|---|
| **F-ST-01 · Character separation** | T-06 | Alaric knows the vault code; Brynn does not | Ask Brynn for the code; ask Alaric in Brynn's presence |
| **F-ST-02 · Inventory** | T-07 | Lantern present; rope present; waterskin consumed at turn 12 | Offer the consumed waterskin; ask to light the lantern after its removal at turn 20 |
| **F-ST-03 · Equipment** | T-08 | Cloak worn; boots removed at turn 8 | Describe boots in scene 9; ask what is worn |
| **F-ST-04 · History & relationships** | T-09 | Sister's name; a promise made at turn 4; a quarrel at turn 30 | Ask the sister's name at turn 60; ask about the promise after the quarrel |
| **F-ST-05 · Conflict** | T-10 | Two sources disagree on the meeting hour (dusk vs midnight) | Ask the plan directly; ask which source is trusted |
| **F-ST-06 · Missing** | T-10 | No record of the ship's name | Ask the ship's name; ask whether it was ever mentioned |

Each fixture ships with: canonical state JSON, the probe list, the *expected* honours (which fact must appear / must not), and the failure signature (e.g. "invents a ship name" = fabrication; "names the vault code for Brynn" = leakage).

## 5. Why synthetic fixtures rather than corpus items

| Reason | Consequence |
|---|---|
| Corpus items are frozen standards artifacts | Editing them for engineering is forbidden; fixtures are disposable |
| Corpus units test *cue criteria*, a judgement question | Fixtures test *state honouring*, an observation question |
| Engineering needs many small probes, cheaply | Fixtures are generated and versioned freely |

**The trade:** fixtures do not inherit the corpus's construct work. A fixture passing says nothing about `R-1`…`R-6`, and a corpus item's status says nothing about a fixture. The two evidence streams stay separate.

## 6. Fixture format

```
fixture_id, version, dimension
state: { characters, persona, world, facts, items, equipment, relationships, events }
probes: [ { turn, user_message, must_include[], must_exclude[], failure_signature } ]
expected_state_hash_before, expected_state_hash_after
```

**Deterministic ids** (`fact-001`, `item-002`) so diffs are readable and hashes stable across regenerations.

## 7. What this interface cannot do

- It cannot prove the model *used* the state — only that the state contained the fact and the reply honoured or violated it (`C2` vs `C4`).
- It cannot represent retrieval mechanics — that is a Howl-stage concern.
- It cannot make a reply's correctness attributable on its own; `C3` (correct from priors) remains a judgement the record must state.

## 8. Claim ceiling

Specification only. No state exists, no fixture has been instantiated, no probe asked. It defines what would have to be true for a later result to be checkable.
