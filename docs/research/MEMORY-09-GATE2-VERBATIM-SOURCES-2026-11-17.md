# MEMORY-09 — Gate 2: the verbatim implementing sources

**What this file is.** The exact text of every implementation delta row and of the operative
contract passages that governed the Gate 2 implementation, quoted from the accepted documents so
the report can be reviewed against the source rather than against a summary of it.

| Field | Value |
|---|---|
| Quoted from | `docs/research/MEMORY-09-ARCH-CONTRACTS-R2-FOUR-QUESTIONS-{REV-2…REV-13}-2026-11-17.md` |
| Nature | **Quotation only.** No paraphrase, no summary, no inference. Where a passage is elided it is marked `[…]` |
| Authority | The quoted files remain authoritative; this file carries no contract status of its own |
| Executed | Nothing |

Numbering note: the delta rows are not independent — each revision **replaces** specific rows of the
previous one ("REV 3 §7 stands, with these four rows amended"). The row actually implemented is
therefore the last revision to restate it, as recorded in the report's §3 table.

---

## 1. REV 3 §7.1 — corrections confirmed, awaiting implementation authorisation

> | # | File | Symbol | Change |
> |---|---|---|---|
> | 1 | `engine.py` | `promote` | Delete `supersede` and its displacement block (D-A, RAT-1) |
> | 2 | `engine.py` | `retract` | Refuse belief records |
> | 3 | `engine.py` | new `revoke_belief` | Model-gated, append-only revocation (D-B, D5) |
> | 4 | `engine.py` | `promote` | Refuse a belief citation unless the proposition is a `belief_ascription` of that owner |
> | 5 | `engine.py` | `acquire` | Evidence table; own-account exempt; inference basis owned; no truth reads |
> | 6 | `retrieval.py` | `world` | `unresolved_relation` for ≥2 live records; never `settled` |

## 2. REV 3 §7.2 — structural work, now with owner direction

> | # | File | Symbol | Change | Direction |
> |---|---|---|---|---|
> | 7 | `engine.py` | `_propositions` registry | Identity, `kind`, `label`, `about_character`, domain | U-1, U-4, R-E |
> | 8 | `records.py` | `Record` | `proposition_id`, required on truth classes | U-1 |
> | 9 | `engine.py` | `propositions(domain, target)` | Enumerated lookup; id-ordered | R-E |
> | 10 | `engine.py` | `possible_duplicates(domain, target)` | Multiplicity report with the non-evidence disclaimer | R-E |
> | 11 | `engine.py` | `promote` | Required `proposition` (id or mint); `coexisting_propositions` in the result | R-E |
> | 12 | `engine.py` | new `assert_incompatible` | Assertion with applicability, spanning propositions | U-7 |
> | 13 | `engine.py` | `_promotions`, `snapshot` | Persist groups, propositions, authorships | U-2, R-E |
> | 14 | `retrieval.py` | new `belief_facts` | Belief-fact view; world excludes ascriptions | U-4 |
> | 15 | `authority.py` | `Grant`, `require` | `operations` + `origin`; per-capability revocation | R-A |
> | 16 | `engine.py` | new `author_ambient` | Authoring capability; `claim` + `declaration` | U-5, U-8 |
> | 17 | `__init__.py` | `__all__` | Export `assert_incompatible`, `author_ambient`, `revoke_belief` | 3, 12, 16 |

> **Untouched:** `eval/mem1/model.py` and the G1 harness; production Howl; r32; plan v3; Amendment 1; v9;
> the map; the ledger; REV 1; REV 2.

---

## 3. REV 4 §10 — implementation list, deltas only

> REV 3 §7 stands, with these four rows amended. Nothing is applied.
>
> | # | File | Symbol | Change | Direction |
> |---|---|---|---|---|
> | 8′ | `records.py` | `Record` | `proposition_id` required on truth classes, **refused on `claim`/`account`/`belief`**, never assignable after construction | O-1, rule 5, §4.3 |
> | 9′ | `engine.py` | `propositions(domain, target, *, actor)` | `target` required; `authority.require(domain, actor, "promotion")` first; L1 + `_truth_projection`; `viewer`/`knowledge_conferred` in the result | §3, §5.1 |
> | 10′ | `engine.py` | `possible_duplicates(domain, *, actor, target=None)` | Same gate; ids and counts only; disclaimer in every response | §5.2 |
> | 11′ | `engine.py` | `promote` | `coexisting_propositions` as `{count, ids}`; no labels or texts | §5.3 |
>
> One new internal symbol, shared with the existing surface:
>
> | # | File | Symbol | Change |
> |---|---|---|---|
> | 17 | `engine.py` | `_truth_projection(domain, target)` | The single active-truth predicate, called by both `world()`'s backing `authoritative()` path and enumeration |

---

## 4. REV 5 §9 — implementation list, deltas only

> REV 4 §7 and §10 stand, with these rows replaced. Nothing is applied.
>
> | # | File | Symbol | Change |
> |---|---|---|---|
> | 17′ | `engine.py` | `_eligible_truth(domain, target)` | REV 4's `_truth_projection`, renamed and **restricted to eligibility**; called by `authoritative()` and by enumeration |
> | 17a | `engine.py` | `_proposition_records(prop_id, domain, target)` | `_eligible_truth` filtered by exact `proposition_id`; no fallback, no widening |
> | 17b | `engine.py` | `_proposition_state(prop_id, domain, target)` | Per-identity state; `conflict` if a member is in an open group, else by count |
> | 9″ | `engine.py` | `propositions(domain, target, *, actor)` | `state` and `records` from 17b/17a; partition invariant checked; `viewer`/`knowledge_conferred` |
> | 10″ | `engine.py` | `possible_duplicates(domain, *, actor, target=None)` | Same gate; ids/counts only; `disclosure: "authorised_audit"`; no target in domain-wide mode |
>
> `authoritative()` keeps its semantics; it calls `_eligible_truth` instead of an inline comprehension. That
> is a refactor of one predicate, not a behaviour change.

---

## 5. REV 6 §8 — implementation list, deltas only

> REV 5 §9 stands, with these rows replaced. Nothing is applied.
>
> | # | File | Symbol | Change |
> |---|---|---|---|
> | 6′ (REV 3 §7.1) | `retrieval.py` | `world` | Returns the per-proposition envelope; `multiple`/`conflict` as §4; never flattens entries into one record list |
> | 17b′ | `engine.py` | `_accepts(view, kind)` | The query-kind predicate; unknown view false |
> | 17c | `engine.py` | `authoritative(domain, target, *, view)` | Required view; applies `_accepts` **inside**; groups by identity; top-level status per §4 |
> | 9‴ | `engine.py` | `propositions(...)` | Uses `view="propositions"` (both kinds); otherwise REV 5 §5.1 |
> | 18 | `retrieval.py` | new `belief_facts` | `authoritative(..., view="belief_facts")`; returns `about_character`; no record outside ascription identities |
> | — | `__init__.py` | `__all__` | Export `belief_facts` alongside the REV 3 §7.2 list |
>
> REV 3 §7.1 item 6 is narrowed to "≥2 live records **on one proposition** → `unresolved_relation`"; the
> cross-identity aggregate is withdrawn.

---

## 6. REV 7 §9 — implementation list, deltas only

> REV 6 §8 stands, with these rows replaced. Nothing is applied.
>
> | # | File | Symbol | Change |
> |---|---|---|---|
> | 17b″ (REV 6 17b′) | `engine.py` | `_accepts(view, kind)` | Unchanged in body; `KNOWN_VIEWS` named once so §5.1 step 2 and the dispatch cannot drift |
> | 17c′ (REV 6 17c) | `engine.py` | `authoritative(domain, target, *, view, actor=None)` | View required; **gate at step 3 for `propositions`**, before any read; uniform `no_domain_authority`; per-view dispatch; envelope per §4 |
> | 6″ (REV 6 6′) | `retrieval.py` | `world` | Unchanged shape; `open_groups` filtered to establishing groups; `integrity` codes only; no foreign id, text or count |
> | 9‴′ (REV 6 9‴) | `engine.py` | `propositions(...)` | Keeps its own check before building the envelope; delegates retrieval to `authoritative(..., view="propositions", actor=actor)`; reports `integrity` with `group_id` |
> | 18′ (REV 6 18) | `retrieval.py` | `belief_facts` | Unchanged shape; same `integrity` rule as `world` |
> | 17d | `engine.py` | `group_establishes(group, view, accepted_ids)`, `group_accepted`, `group_contained` | New; the only place §3.1's three-way relation is decided |
> | 17e | `engine.py` | `_proposition_state(prop_id, *, domain, target, view, accepted_ids)` | View-relative; conflict before counts; non-establishing groups ignored |
> | 12′ (REV 3 row 12) | `engine.py` | `assert_incompatible` | Refuses `cross_kind_group` before mutation; still spans propositions of one kind (REV 3 §5.4 unaffected) |
> | 19 | `records.py` / `engine.py` | kind of a proposition | One accessor used by `_accepts`, `group_accepted` and the §3.2 integrity check; no second derivation |

---

## 7. REV 8 §9 — implementation list, deltas only

> REV 7 §9 stands, with these rows replaced. Nothing is applied.
>
> | # | File | Symbol | Change |
> |---|---|---|---|
> | 17f | `engine.py` | `_accepted(domain, target, view)` | Computes the accepted record list once; every consumer takes it as an argument |
> | 17d′ (REV 7 17d) | `engine.py` | `group_establishes(group, view, accepted)` | `group_accepted` + `group_visible` (record ids) + `group_contained` (identities); unresolved member fails both |
> | 17e′ (REV 7 17e) | `engine.py` | `_proposition_state(prop_id, *, view, accepted)` | Records read from `accepted`; establishing groups only; conflict before counts |
> | 6‴ (REV 7 6″) | `retrieval.py` | `world` | `relation` per §4; `open_groups` establishing-only, all members inside `accepted`; `integrity` codes only |
> | 18″ (REV 7 18′) | `retrieval.py` | `belief_facts` | Same rules as `world` |
> | 9‴″ (REV 7 9‴′) | `engine.py` | `propositions(...)` | Registry-sourced identity set, gated; `open_groups` includes non-establishing groups with `establishes`/`cause` |
> | 17c″ (REV 7 17c′) | `engine.py` | `authoritative(...)` | Unchanged except that it computes `accepted` once and passes it; gate order per REV 7 §5.1 |

---

## 8. REV 9 §9 — implementation list, deltas only

> REV 8 §9 stands, with these rows replaced. Nothing is applied.
>
> | # | File | Symbol | Change |
> |---|---|---|---|
> | 6⁗ (REV 8 6‴) | `retrieval.py` | `world` | Emits `open_groups` (= establishing) and `establishing_groups` (identical); `relation` from the latter |
> | 18‴ (REV 8 18″) | `retrieval.py` | `belief_facts` | Same as `world` |
> | 9⁗ (REV 8 9‴″) | `engine.py` | `propositions(...)` | `open_groups` = all relevant with `establishes`/`cause`; `establishing_groups` = governing subset; `relation` from the latter |
> | 17g | `engine.py` | `_relevant_groups(view, …)` | Derives both collections from one pass; `establishing` computed once per group per call |
> | 17d″ | `engine.py` | `group_establishes(...)` | Unchanged from REV 8; called through one accessor so state and `establishing_groups` cannot diverge |

---

## 9. REV 10 §7 — implementation list, deltas only

> REV 9 §9 stands, with these rows replaced or added. Nothing is applied.
>
> | # | File | Symbol | Change |
> |---|---|---|---|
> | 17h | `engine.py` | `_snapshot(domain, target, view)` | Takes `accepted`, `resolution`, `groups`, `establishes` once; returns an immutable value; the only entry point for state reads on a response path |
> | 17d‴ (REV 9 17d″) | `engine.py` | `group_establishes(...)` | Pure over `S`; results memoised **in `S`**, not per call site |
> | 6⁵ (REV 9 6⁗) | `retrieval.py` | `world` | Consumes `S` only; no post-`S` read |
> | 18⁴ (REV 9 18‴) | `retrieval.py` | `belief_facts` | Same |
> | 9⁵ (REV 9 9⁗) | `engine.py` | `propositions(...)` | Consumes `S` only; gate still before `S` is taken (S-5) |

---

## 10. REV 11 §8 — implementation list, deltas only

> REV 10 §7 stands, with these rows replaced or added. Nothing is applied.
>
> | # | File | Symbol | Change |
> |---|---|---|---|
> | 17h′ (REV 10 17h) | `engine.py` | `_capture(domain, target, view, actor)` | Obtains one admissible `S` per S-7/S-8, including `version`, `returned_identities` and complete member resolution; gate runs before it is called |
> | 17i | `engine.py` | `_version` | Monotone counter advanced in the same step as each mutating call's publication; read only by `_capture` for admissibility, never published |
> | 17j | `engine.py` | `_relevant_groups(S)` | Selects by `S.returned_identities`; no registry access |
> | 6⁶ (REV 10 6⁵) | `retrieval.py` | `world` | Consumes `S` only |
> | 18⁵ (REV 10 18⁴) | `retrieval.py` | `belief_facts` | Consumes `S` only |
> | 9⁶ (REV 10 9⁵) | `engine.py` | `propositions(...)` | Consumes `S` only, including its identity set; enumeration reads `S.returned_identities` |

---

## 11. REV 12 §7 — implementation list, deltas only

> REV 11 §8 stands, with these rows replaced or added. Nothing is applied.
>
> | # | File | Symbol | Change |
> |---|---|---|---|
> | 17h″ (REV 11 17h′) | `engine.py` | `_capture(...)` | Follows A1–A11 explicitly; A3 takes the whole open group store; A10 returns a classified failure |
> | 17k | `engine.py` | `_resolve_members(refs)` | Total, three-valued map; `unresolved(record_absent)` for a missing record; never omits a key |
> | 17l | `engine.py` | `_relevant_groups(candidates, returned_identities, resolution)` | Relevance applied after resolution; no registry or record-store access |
> | 17m | `engine.py` | refusal classes | Two shapes sharing one envelope: `Refusal` (Class A) and `CaptureRefusal` (Class B); both carry only `code` |
> | 17j′ (REV 11 17j) | `engine.py` | `_relevant_groups` | Superseded by 17l |
> | 6⁷ (REV 11 6⁶) | `retrieval.py` | `world` | Propagates the refusing class unchanged; no field added to either |
> | 18⁶ (REV 11 18⁵) | `retrieval.py` | `belief_facts` | Same |
> | 9⁷ (REV 11 9⁶) | `engine.py` | `propositions(...)` | Same; gate remains A2, before every capture step |

## 12. REV 13 §7 — implementation list, deltas only

> REV 12 §7 stands, with these rows replaced or added. Nothing is applied.
>
> | # | File | Symbol | Change |
> |---|---|---|---|
> | 17h‴ (REV 12 17h″) | `engine.py` | `_capture(...)` | A1 now validates only `view`; an optional catalogue lookup is A6′; gate at A2 consults metadata and nothing else |
> | 17n | `engine.py` | `_validate_inputs(view)` | Total on the view name; no store access; the only input validation in the call |
> | 17o | `engine.py` | `_catalogue_lookup(...)` *(conditional)* | Only if a catalogue exists; called at A6′, inside the capture; empty result, never a refusal |
> | 17m′ (REV 12 17m) | `errors.py` | refusal classes | Both classes carry `code` only; Class A additionally may consult metadata before constructing the code |
> | 17p | `authority.py` | — | **No change.** `require(domain, actor, operation)` — REV 3's signature — is consulted as-is; the metadata it reads is not edited or extended |
> | 6⁸ / 18⁷ / 9⁸ | `retrieval.py`, `engine.py` | `world`, `belief_facts`, `propositions` | Unchanged from REV 12 apart from propagating the A1/A6′ split |

---

## 13. Operative passage — the refusal classes (REV 13 §2.2–§2.3)

> | | Class A — authorization | Class B — capture |
> |---|---|---|
> | When | Before any data read | After the gate admits the call; during acquisition |
> | Cause | The caller may not make this call, or the view does not exist | The caller may make it; the state could not be captured coherently |
> | Data read | **None.** Not the record store, not the group store, not the registry | Reads attempted; whatever was read is discarded |
> | `S` | Does not exist | No admissible `S` exists; partial state is dropped |
> | Codes | `no_domain_authority`, `unknown_view` | `inconsistent_snapshot` |
> | Views reachable | `propositions` (gate), any (`unknown_view`) | All three — each captures |
> | Envelope | Fixed refusal shape; no data-derived field | **Identical** fixed refusal shape; no data-derived field |

> | Rule | Statement |
> |---|---|
> | S-5 (narrowed) | **Authorization refusals take no snapshot and read no data.** The gate, and view validation, complete before any store access; a Class A refusal has no `S` and its envelope contains no id, count, text or version. |
> | S-9 (new) | **Capture refusals publish nothing derived.** A Class B refusal may have read state while attempting capture; it must discard it and return the same fixed refusal shape, with no field derived from any candidate snapshot, no partial collections, and no `version`. |
> | S-10 (new) | **Ordering.** Class B is unreachable before the gate. An unauthorized caller therefore only ever receives a Class A code, and a Class B code can never disclose whether a domain holds data. |
> | S-11 (new) | **Content independence.** The choice of refusal class and code must not depend on the content, size or emptiness of the stores. `inconsistent_snapshot` reports an acquisition outcome, not a fact about the domain. |

## 14. Operative passage — the authorization boundary (REV 13 §2.2)

> | Class | Contents | May Class A read it? |
> |---|---|---|
> | **Authorization metadata** | Grant records keyed by `(domain, actor, operation)`; delegation edges and their origins; the operation names a grant may carry | **Yes** — this is what the gate is for |
> | **Protected memory data** | The record store and its record contents; the group store and its members; the proposition registry and its identities; and every value derived from them (`accepted`, `returned_identities`, `resolution`, `groups`, `establishes`) | **No** — unreachable until the gate passes |

> | Rule | Statement |
> |---|---|
> | S-5 (restated) | **Authorization refusals read no protected memory data and take no snapshot.** They may consult authorization metadata, and only that, to evaluate the gate. Their envelope contains no protected-memory field and no authorization-metadata field. |
> | S-12 (new) | **Metadata boundary.** Protected memory data is read only after the gate admits the call, and only inside a capture attempt. Authorization metadata is read only to evaluate the gate and to choose the refusal code; it is never published, summarised, counted or hinted at in any response, refused or served. |
> | S-11 (extended) | **Content independence** now covers two repositories: neither the refusal class nor the code may vary with the content of the stores **or** with the content of the authorization metadata beyond the boolean outcome of the gate. |

## 15. Operative passage — A1–A11 (REV 12 §3.1, amended by REV 13 §3.1)

> | Step | Action | On failure |
> |---|---|---|
> | A1 | Validate `view`; resolve `(domain, target)` | Class A: `unknown_view` |
> | A2 | For `view == "propositions"`: `require(domain, actor, "promotion")` | Class A: `no_domain_authority` |
> | A3 | **Candidate set:** `candidate_groups = every group in the open group store` — unfiltered | — |
> | A4 | **Member references:** `member_refs = ⋃ members(g) for g ∈ candidate_groups`, as a set of record ids | — |
> | A5 | **Total resolution:** `resolution = { m ↦ resolve(m) for every m ∈ member_refs }`, three-valued (§3.2) | — |
> | A6 | **Accepted records:** `accepted` = eligible ∧ kind-accepted for this view | — |
> | A7 | **Returned identities:** canon from `accepted`; authority from the registry, target-scoped | — |
> | A8 | **Relevance:** `groups = [g ∈ candidate_groups : ∃ m ∈ members(g), pid(resolution[m]) ∈ returned_identities]` | — |
> | A9 | **Establishment:** `establishes[g]` for every `g ∈ groups`, against `accepted` and `resolution` | — |
> | A10 | **Version recheck** | Changed → discard everything, retry; attempts exhausted → Class B: `inconsistent_snapshot` |
> | A11 | Assemble the response from `S` | — |

REV 13 §3.1 amends A1 to:

> | Step | REV 12 | REV 13 |
> |---|---|---|
> | A1 | "Validate `view`; resolve `(domain, target)`" | **Validate `view`** only — a name check against the three defined views, no store access. Class A: `unknown_view` |

## 16. Operative passage — resolution values (REV 12 §3.2, restated by REV 13 §4)

> | Value | Meaning | Effect |
> |---|---|---|
> | `resolved(pid)` | Examined, found, identity known | Usable by A8 and by `group_visible` / `group_contained` |
> | `unresolved(record_absent)` | **Examined, not found** — a divergence between the group store and the record store | `group_visible` and `group_contained` both fail (REV 8 §3.2); authority view reports `establishes: false`, `cause: "member_unresolved"`; canon views surface `group_member_unresolved` in `integrity` only |

> There is no third value, and in particular **no absence**. "Not examined" is unrepresentable because an
> incomplete capture is a failed capture (A10/S-9), not a partially-filled one.

## 17. Operative passage — S-1…S-8 (REV 11 §3.2)

> | Rule | Statement |
> |---|---|
> | S-1 | **One take.** `S` is obtained once per call, before the first read that contributes to the response. |
> | S-2 | **No re-read.** Every response field is a projection of `S`. No field is computed from a read made after `S` — including any read of the registry. |
> | S-3 | **One map, two projections.** `establishes` is computed once per group against `S.accepted` and `S.resolution`; entry `state`s and `establishing_groups` are read from it, never recomputed. |
> | S-4 | **Staleness permitted, mixing not.** A concurrent write may make `S` stale; the response must still be a function of exactly one `S`. No read-your-writes guarantee is offered within a call. |
> | S-5 | **Refusals take no snapshot.** A refused call stops at the gate before any data read, so it has no `S` and cannot leak one. |
> | S-6 | **Assembled means assembled from `S`.** "Response `R` is consistent" means: there exists an `S` such that `R = assemble(S)`. |
> | **S-7** | **Coherent capture.** Every field of `S` must have co-existed at one committed state. `S` is admissible only if there is a single committed version `v` with `accepted = accepted@v`, `returned_identities = identities@v`, `resolution = resolution@v` and `groups = groups@v`. Atomic snapshot, transactionally consistent read, and version-check-and-retry all satisfy this; the contract prescribes none of them. |
> | **S-8** | **Frozen identity set.** `returned_identities` is captured inside the same coherent read as the rest of `S`, populated by the rule for the selected view (§4.1). Enumeration, group relevance and scope checks read the frozen value; none may consult the registry afterwards. |

> **S-7's soundness condition, stated because a version test is easy to get subtly wrong.** A
> compare-then-recheck on a version counter proves admissibility only if every mutating call publishes its
> changes and advances the version in **one step** — no reader may observe a partial write under an unchanged
> version.

## 18. Operative passage — group establishment (REV 8 §3.2–§3.3)

> ```
> group_accepted(g, view)        = all(_accepts(view, kind_of(m))   for m in g.members)
> group_visible(g, accepted_ids) = all(m in accepted_ids            for m in g.members)   # NEW, record-level
> group_contained(g, prop_ids)   = all(proposition_of(m) in prop_ids for m in g.members)
> group_establishes(g, view, accepted) =
>         group_accepted(g, view) and group_visible(g, ids(accepted)) and group_contained(g, props(accepted))
> ```

> | Cause | What it means | Fault? | Canon view | Authority view |
> |---|---|---|---|---|
> | Kind straddle | One relation claims two questions | **Integrity fault** | `integrity: ["group_spans_query_kinds"]`, no ids | `open_groups` entry, `establishes: false`, `cause: "spans_query_kinds"` |
> | Member not current or not in scope | Superseded, revoked, or another applicability | **Not a fault** — this is history (rules 20, 44) | Silent | `open_groups` entry, `establishes: false`, `cause: "member_not_accepted"` |
> | Member id unresolved | The registry and record store have diverged | **Integrity fault** | `integrity: ["group_member_unresolved"]`, no ids | same, `cause: "member_unresolved"` |

## 19. Operative passage — `relation` and the two collections (REV 9 §4, §3.2)

> ```
> "asserted"        iff establishing_groups(view) is non-empty
> "not_established" otherwise
> ```

> | Biconditional | Holds in | Why |
> |---|---|---|
> | `relation == "asserted"` ⟺ `establishing_groups` non-empty | **every view** | This is the definition |
> | `relation == "asserted"` ⟺ some entry has `state == "conflict"` | **every view** | Establishing members are in `accepted`, so some entry is in conflict; a conflicting entry has an establishing group (REV 8 §3.4) |
> | `relation == "asserted"` ⟹ `open_groups` non-empty | every view | `establishing_groups ⊆ open_groups` |
> | `open_groups` non-empty ⟹ `relation == "asserted"` | **canon only** | In the authority view this implication is false by design (§3.3) |

> Canon views: the two collections coincide
> ```
> open_groups(world)         == establishing_groups(world)
> open_groups(belief_facts)  == establishing_groups(belief_facts)
> ```
> Both are literal, identical lists.

## 20. Operative passage — the gate order (REV 7 §5.1)

> ```
> authoritative(self, domain, target, *, view, actor=None):
>
>   1  view == "propositions"       → requires actor; no default
>   2  view not in KNOWN_VIEWS      → refuse("unknown_view")            # fail closed, F4
>   3  view == "propositions":
>   3a     actor is None            → refuse("no_domain_authority")
>   3b     authority.require(domain, actor, "promotion")                # R-A item 4
>   3c     AuthorityError           → refuse("no_domain_authority")     # uniform, detail withheld
>   4  eligible = _eligible_truth(domain, target)                       # first data read
>   5  accepted = [r for r in eligible if _accepts(view, kind_of(r))]
>   6  identity set, state, envelope as §3–§4
> ```

> | Property | Consequence |
> |---|---|
> | Step 3 precedes step 4 | No record, no identity, no count is touched before authorisation |
> | Step 3c is uniform | "Domain not declared" and "no live `promotion` grant" and "`actor=None`" are the same refusal: `{"status": "refused", "reason": "no_domain_authority"}` and nothing else. |

## 21. Operative passage — enumeration refusals and layers (REV 4 §4.1, §4.3, §7)

> | Field | Returned | Notes |
> |---|---|---|
> | `id` | yes | `prop-00007`; opaque, never derived from text |
> | `label` | yes | The minter's identity string, verbatim. Actor-authored text; see §4.4 |
> | `kind` | yes | `claim_about_world` / `belief_ascription` / … (closed set, REV 2) |
> | `about_character` | yes, when the proposition is an ascription | The subject of an **authoritative** record the caller may read anyway |
> | `domain`, `target` | yes | Domain is part of the identity (REV 3 §5.2) |
> | `created_by`, `created_turn` | yes | Provenance already required of promotions (K9) |
> | `state` | yes | `settled` / `unresolved_relation` / `conflict` / `none`, computed as today |
> | **any record id** | **no** | L2 only |
> | **any record text** | **no** | L2 only |

> | Rule | Effect |
> |---|---|
> | `Record.proposition_id` is **required** on `world_truth` / `scenario_truth` | Every authoritative record is attached to one identity (REV 3 item 8) |
> | `Record.proposition_id` is **refused** on `claim`, `account` and `belief` | A private record cannot be attached to a proposition, so enumeration has nothing to traverse |
> | No operation sets `proposition_id` after construction | No late attachment, no backfill, no migration path |
> | The proposition registry stores `label` and `kind`, never a record reference | The registry cannot be used as an index into private records |

> | Situation | Result |
> |---|---|
> | No live grant in `domain` | `refused` / `no_domain_authority`; **no identities, no counts, no ids** |
> | Grant in a different domain (rule 35) | `refused` / `no_domain_authority` |
> | Caller asks for a character view of enumeration | `refused` / `not_a_perspective_surface` |
> | `target` omitted on §5.1 | `refused` / `target_required` |
> | Unknown target, valid authority | `ok`, empty `propositions` (abstention is not an error) |
> | Proposition whose records are all superseded/revoked | L1 identity present, `state: none`, `records: []` |

## 22. Operative passage — the partition invariant (REV 5 §4)

> **For any `(domain, target)` with authorised enumeration:**
>
> 1. Every record in every entry's `records` satisfies `r.proposition_id == entry.id`.
> 2. No record id appears in two entries — guaranteed by (1) and single-valued `proposition_id`.
> 3. The union of all entries' `records` equals `_eligible_truth(domain, target)`.
> 4. No entry has a record whose `domain` or `target` differs from the query.
> 5. `state` for each entry is computed from that entry's `records` only.
>
> A violation of (1), (2) or (4) is an **integrity failure reported as such**, not silently filtered.
> A violation of (3) — an eligible truth record attached to no enumerated identity — is likewise an
> integrity failure, because `proposition_id` is required: an orphan means the registry and the record
> store have diverged. The response must not quietly omit it, or enumeration would under-report canon
> while looking complete.

## 23. Operative passage — the canon views (REV 6 §3.2, §4)

> ```
> _accepts(view, kind):
>     "world"        -> kind == "claim_about_world"
>     "belief_facts" -> kind == "belief_ascription"
>     "propositions" -> kind in {"claim_about_world", "belief_ascription"}   # authoring surface
>     otherwise      -> False            # fail-closed
> ```

> | Situation | `status` | Note |
> |---|---|---|
> | No accepted identity | `none` | Abstention, not error |
> | One identity | That identity's state: `settled`, `unresolved_relation`, `conflict`, or `none` | Same as today for the single-identity case |
> | Several identities, no open group spanning them | `multiple` | **The query is underspecified as to identity**, not in unresolved conflict |
> | Several identities, open group spanning them | `conflict` | The relationship was asserted; `open_groups` names it |

## 24. Operative passage — proposition identity (REV 3 §5.2, §5.4)

> | Rule | Consequence |
> |---|---|
> | `proposition` is a **required** parameter | There is no default identity; omission is a refusal, not an inference |
> | Minting requires a `label` | The actor's own identity string, recorded verbatim. It is **never compared**, never normalised, never used for similarity |
> | `kind` and `about_character` are fixed at mint | A proposition cannot change kind later; a belief-ascription cannot become a world claim |
> | Ids are monotonic and never reused | `prop-00007` remains `prop-00007`; there is no delete, no re-mint, no id recycling |
> | Domain is part of the identity | A proposition belongs to one world or scenario; referencing it from another domain is refused |

> Refused when: the id is unknown; the id belongs to another chat (rule 31); the id belongs to another
> domain (rule 35); or the promotion's declared applicability disagrees with the proposition's `target`.

## 25. Operative passage — acquisition evidence (REV 2 §5.3)

> | Kind | Evidence required at write | Refused when |
> |---|---|---|
> | `witnessed` | acquisition event reference (`span`) attributable to the character | no span; character absent from the event |
> | `told` | disclosure event reference and `descriptor.speaker` ≠ owner | no speaker; self-told; no event |
> | `read` | `descriptor.source` **and** a reading event | no source; no event |
> | `inferred` | `descriptor.basis` = ids of **the same character's** existing accounts | empty basis; any id not owned by that character; any id that does not exist |
> | `verified` | `descriptor.verification` event **and** its basis | no verification event; no basis |
> | `narrator_grant` | `descriptor.declaration` resolving to a live declaration in this chat | unknown or withdrawn declaration |
> | `none` (`own_account`) | **nothing** — the absence is the record's meaning | never |

> **Structural invariant.** The validator may read only: the record being written, records already owned by
> that character, the declaration registry, and the event reference. It may never call `authoritative()`.
> Evidence establishes *acquisition*, never *correctness*.

## 26. Operative passage — ambient authorship (REV 2 §6.1, §6.3)

> **Adopted:** `claim` + `declaration` + provenance block `{"authorship": {actor, domain, grant_id, turn}}`.
>
> **The failure condition, named:** reuse holds **because** ambient material is assertoric. If authored
> material is later required to be a **definition** … the claim class cannot represent it without
> collapsing the distinction B3 is waiting on.

> `author_ambient(actor, domain, target, text, turn, rationale=None, revises=None)`
>
> - requires `authorship` capability in that domain;
> - creates a `claim` with kind `declaration`, contributor the actor;
> - never creates truth; never appears in a world answer;
> - `revises` appends a new record and marks the prior `superseded`, preserving text and history;
> - is a lawful `record_cited` basis for a later promotion, which remains the only path to authority.
>
> **Refusals:** no authorship capability · unknown domain · empty text · `revises` naming a record from
> another chat or a truth record · attempting to author in a domain where the actor holds only promotion
> capability.

## 27. Operative passage — the authority model (REV 2 §7)

> | Element | Change |
> |---|---|
> | `Grant` | gains `operations: frozenset[str]` over `{"promotion", "authorship"}` |
> | Domain declaration | installs the human operator with **both** capabilities |
> | Model grant (8a-6) | defaults to `{"promotion"}`; authorship requires an explicit grant naming it |
> | `require(domain, actor)` | becomes `require(domain, actor, operation)` |
> | Revocation | per grant; revoking promotion capability leaves authorship capability untouched, and vice versa |

> This is **one field on an existing object** plus a parameter on an existing check.

---

## 28. What is not in this file

| Not quoted here | Where it is |
|---|---|
| The 57 rules and 118 requirements | `MEMORY-09-ARCHITECTURE-v9-2026-11-17.md` §4; `MEMORY-09-CONFORMANCE-MAP-2026-11-17.md` |
| The owner's acceptance and its verbatim terms | `MEMORY-09-SNAPSHOT-CONTRACT-ACCEPTANCE-2026-11-17.md` §1 |
| REV 11 §2's tear trace and worked illustrations | REV 11 §2, §6; summarised in the report's §4 |
| The prose of REV 3 §5.5–§5.7 (multiplicity) | REV 3 §5.5–§5.7 |

Every quotation above is checkable against its source section; if any transcription differs from the
source, the source governs and this file is in error.
