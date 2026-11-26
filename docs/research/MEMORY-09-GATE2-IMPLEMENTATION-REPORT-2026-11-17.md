# MEMORY-09 — Gate 2 implementation report

**Status: IMPLEMENTATION COMPLETE for the authorised scope, with one blocked item reported (§5).
Nothing executed. Nothing compiled. Nothing linted. No test, probe or invariant check ran.**

| Field | Value |
|---|---|
| Date | 2026-11-17 CST (−0600) |
| Authorisation | Owner's Gate 2 authorisation, §1–§8 (implementation only) |
| Contract implemented | REV 13 §2–§5 as operative snapshot contract; REV 11 §3.2 S-1…S-11; REV 12 §2–§3 |
| Ratified base | 57 v9 rules; 118 conformance requirements; r32 tier 2; owner decisions tier 1 |
| Scope touched | `eval/mem1/m09/` only |
| `authority.py` | **Not modified** — sha256 `07dfbeb2…`, unchanged |
| Transfer / cross-chat | **Not introduced**, and structurally absent |
| Verification | **Source inspection only.** No execution of any kind |

---

## 1. D-A and D-B corrections

### D-A — `promote(supersede=True)` bypassed the human gate

| Item | Detail |
|---|---|
| Defect | `promote` accepted `supersede=True`, marked prior authoritative records `SUPERSEDED` and set `superseded_by`, without calling `_require_human`. A model holding a standing promotion grant could displace an existing promotion, including its own. |
| Correction | The `supersede` parameter and its displacement block are **deleted**. `promote` now creates a record and displaces nothing, on every path. |
| Displacement path | `supersede_authority` only (`engine.py` `supersede_authority`), which calls `_require_human` before any write and then performs both writes inside one mutation. |
| Evidence | `engine.py` `promote` signature: `(*, actor, domain, target, text, proposition, basis, turn, rationale=None)` — no `supersede`. No `SUPERSEDED` assignment exists inside `promote`; the occurrences in the module are `supersede_authority` and the two belief-lifecycle methods. |
| Consequence | Two live promotions on one proposition now produce `unresolved_relation`, the state O-10 describes as unresolved, rather than silent displacement. |

### D-B — `retract()` could revoke another character's belief

| Item | Detail |
|---|---|
| Defect | `retract(record_id)` set `status = REVOKED` on any non-truth record with no ownership or role check, so anyone could revoke a belief. |
| Correction | `retract` **refuses belief records outright** (`BeliefError`), before mutating. Belief lifecycle moves to model-gated operations. |
| New operation | `revoke_belief(*, character, declarer, role, text, turn, record_id=None)` — requires `role == "character_model"` and `declarer == character`; creates a **new** belief record carrying the revocation (`attribution["revokes"]=True`, `revoked_records` listing the priors); marks priors `SUPERSEDED` with `superseded_by`; deletes nothing. |
| Evidence | `engine.py` `retract` (belief branch) and `revoke_belief`; `records.py` `Record.__post_init__` belief branch refuses a foreign or non-model declarer at construction. |
| Refused explicitly | An operator or narrator acting "with authority" through `retract`: domain authority is promotion authority and confers nothing over a character's interior state (O-5). |

---

## 2. Files changed

| File | Lines | Change |
|---|---|---|
| `eval/mem1/m09/records.py` | 192 | `proposition_id` field; required on truth classes, refused on private classes, sealed against re-pointing; `PROPOSITION_KINDS` closed set |
| `eval/mem1/m09/errors.py` | 103 | Added `EvidenceError`, `PropositionError`, `GroupError`, `SnapshotError`, `ViewRefusal` |
| `eval/mem1/m09/engine.py` | 1045 | Proposition registry; conflict groups; assertion; capture/assembly; gated `authoritative`; enumeration; multiplicity; evidence validation; D-A/D-B; version bracket |
| `eval/mem1/m09/retrieval.py` | 101 | `world` and `belief_facts` as pass-throughs to the gated `authoritative`; `CONTEXTS` extended |
| `eval/mem1/m09/__init__.py` | 35 | Exports `DISCLAIMER`, `KNOWN_VIEWS`, `PROPOSITION_KINDS`; `ENGINE_VERSION = "9.1.0-gate2-rev13"` |
| `eval/mem1/m09/authority.py` | 156 | **Unchanged** |
| `eval/mem1/model.py`, G1 harness, production Howl | — | **Unchanged** |

No files added outside `eval/mem1/` except this report. No external module imports `m09` (`grep -rn "m09" eval/mem1/*.py` → no matches), so the changed signatures break no in-package caller.

---

## 3. Implementation deltas applied

| Delta row | Symbol | Location | Applied |
|---|---|---|---|
| REV 3 §7.1 1 | `promote` — `supersede` deleted | `engine.py` `promote` | yes |
| REV 3 §7.1 2 | `retract` refuses beliefs | `engine.py` `retract` | yes |
| REV 3 §7.1 3 | `revoke_belief` | `engine.py` `revoke_belief` | yes |
| REV 3 §7.1 4 | belief citation guard | `engine.py` `promote` (cited-belief loop) | yes |
| REV 3 §7.1 5 | `acquire` evidence table | `engine.py` `_validate_evidence` | yes |
| REV 3 §7.2 7 | `_propositions` registry | `engine.py` `__init__`, `_resolve_proposition` | yes |
| REV 3 §7.2 8 / REV 4 8′ | `Record.proposition_id` | `records.py` `Record` | yes |
| REV 3 §7.2 9 / REV 4 9′ / REV 5 9″ / REV 6 9‴ / REV 7 9‴′ / REV 8 9‴″ / REV 9 9⁗ / REV 10 9⁵ / REV 11 9⁶ / REV 12 9⁷ / REV 13 9⁸ | `propositions(...)` | `engine.py` `propositions` | yes |
| REV 3 §7.2 10 / REV 4 10′ / REV 5 10″ | `possible_duplicates(...)` | `engine.py` `possible_duplicates` | yes |
| REV 3 §7.2 11 / REV 4 11′ | `promote` proposition + `coexisting_propositions` | `engine.py` `promote` | yes |
| REV 3 §7.2 12 / REV 7 12′ | `assert_incompatible` | `engine.py` `assert_incompatible` | yes |
| REV 3 §7.2 13 | `_promotions`, `snapshot` persist groups/propositions | `engine.py` `snapshot` | yes |
| REV 3 §7.2 14 / REV 6 18 / REV 7 18′ / REV 8 18″ / REV 9 18‴ / REV 10 18⁴ / REV 11 18⁵ / REV 12 18⁶ / REV 13 18⁷ | `belief_facts` | `retrieval.py` `belief_facts` | yes |
| REV 3 §7.2 15 | `Grant.operations` + `origin`; `require(..., operation)` | `authority.py` | **NOT APPLIED — prohibited** (§5) |
| REV 3 §7.2 16 | `author_ambient` | — | **NOT APPLIED — blocked** (§5) |
| REV 3 §7.2 17 / REV 4 17 | `_eligible_truth` | `engine.py` `_eligible_truth` | yes |
| REV 4 17 | `_truth_projection` | merged into `_eligible_truth` (REV 5 17′ rename) | yes |
| REV 5 17a | `_proposition_records` | inlined in `_assemble` (records filtered from `S.accepted`) | yes |
| REV 5 17b / REV 7 17e / REV 8 17e′ | `_proposition_state` | `engine.py` `_assemble` per-entry state block | yes |
| REV 6 17b′ / REV 7 17b″ | `_accepts`, `KNOWN_VIEWS` | `engine.py` `_accepts`, `KNOWN_VIEWS` | yes |
| REV 6 17c / REV 7 17c′ / REV 8 17c″ / REV 12 17h‴-lineage | `authoritative(domain, target, *, view, actor)` | `engine.py` `authoritative` | yes |
| REV 6/7/8/9/10/11/12/13 `world` rows (6′…6⁸) | `world` | `retrieval.py` `world` | yes |
| REV 7 17d / REV 8 17d′ / REV 9 17d″ / REV 10 17d‴ | `group_establishes` | `engine.py` `_group_establishes` | yes |
| REV 9 17g / REV 11 17j / REV 12 17j′ / 17l | `_relevant_groups` | `engine.py` `_take` (relevance block) | yes |
| REV 8 17f | `_accepted` | `engine.py` `_take` | yes |
| REV 10/11 17h′ / REV 12 17h″ / REV 13 17h‴ | `_capture` | `engine.py` `_capture` + `_take` | yes |
| REV 11 17i | `_version` | `engine.py` `__init__`, `_mutation` | yes |
| REV 12 17k | total three-valued resolution | `engine.py` `_take` (A4/A5) | yes |
| REV 12 17m / REV 13 17m′ | refusal classes | `engine.py` `_refusal`, `authoritative`, `propositions`, `possible_duplicates` | yes |
| REV 13 17n | `_validate_inputs(view)` | inlined as the first branch of `authoritative` | yes |
| REV 13 17o | `_catalogue_lookup` | **omitted by design** — no catalogue exists; A6′ constrains nothing | yes (as absence) |
| REV 13 17p | `require(domain, actor, operation)` | `authority.py` | **NOT APPLIED — prohibited** (§5) |

---

## 4. The snapshot contract, as implemented

| Rule | Implementation |
|---|---|
| S-1 One take | `authoritative` calls `_capture` once; `_take` performs every read for the response |
| S-2 No re-read | `_assemble` reads only the snapshot dict; `_record_view` formats frozen `Record` references |
| S-3 One map, two projections | `_take` fills `establishes`/`causes` once per group; `_assemble` reads them; entry state and `establishing_groups` both project that map |
| S-4 Staleness permitted, mixing not | Snapshot is a value; a concurrent write may make it stale and cannot splice two states |
| S-5 / S-12 Class A | `authoritative` returns before `_eligible_truth` for `unknown_view` and for a failed gate; the only pre-read call is `AuthorityRegistry.require`, i.e. authorization metadata |
| S-7 Coherent capture | `_capture` brackets `_take` with the version counter; `_mutation` advances it once per outermost mutating call; a mismatch discards and retries; exhaustion → `None` |
| S-8 Frozen identity set | `_take` A7 populates `identities`/`returned_identities`; relevance (A8) and every entry read that frozen value |
| S-9 / S-10 / S-11 Class B | `_refusal("inconsistent_snapshot")` with no data-derived field; unreachable before the gate because the gate returns first; the code does not vary with store content |
| S-6 Assembled from S | Every envelope field is a projection of one snapshot dict |

A-acquisition order is literal: A1 view validation (no read) → A2 gate → A3 whole open-group store → A4 member refs →
A5 total resolution → A6 eligibility and kind acceptance → A7 identity set → A8 relevance → A9 establishment →
A10 version recheck → A11 assembly.

---

## 5. Blocked item — `author_ambient` and the per-operation gate

**This is the conflict §6 of the authorisation asks to be reported rather than resolved.**

| Item | Detail |
|---|---|
| Required by | Authorisation §4: "Respect `Grant.operations`, authority scope, and delegation origin." REV 2 §7: `Grant` gains `operations: frozenset[str]`, `require(domain, actor)` becomes `require(domain, actor, operation)`. REV 3 §7.2 rows 15/16. REV 13 §7 row 17p (§7 asserts the signature is `require(domain, actor, operation)`). |
| Forbidden by | Authorisation §1: "This authorization does NOT authorize … Changes to `authority.py`." §3: "Do not modify `authority.py`." |
| Actual state | `authority.py:71-81` — `Grant` carries `grant_id, domain, holder, holder_kind, granted_by, granted_turn, revoked_turn`. No `operations`, no authorship capability. `authority.py:134` — `require(self, domain, actor)`: two parameters. |
| Consequence | (a) `author_ambient` is **not implemented**: its gate is "authorship capability in that domain", and no such capability can exist without editing `Grant`. (b) The per-operation gate cannot be expressed; `promote` and `propositions` call `require(domain, actor)`, which is promotion-shaped only because promotion is the only capability that exists. (c) REV 13 row 17p's premise ("REV 3's signature — consulted as-is") does not match the source; the source has two parameters. |
| Not done | No workaround. Gating `author_ambient` on `require(domain, actor)` would make promotion authority confer authorship, contradicting U-8 and rule 39 — a resolution the authorisation forbids inventing. |
| Effect on the tally | Rule 39 is **P** (promotion half implemented, authorship half absent). K25 and L8's positive surface are **A**. Everything else is unaffected: no other operation needs a per-operation capability. |

**To unblock:** an authorisation to add `Grant.operations` and the third `require` parameter to
`authority.py` (two changes: one field, one parameter, defaulted so existing call sites keep working).

---

## 6. The 57 v9 rules

Classification: **I** implemented · **P** partial · **A** absent · **C** contradictory · **NA** not assessable.

| # | Rule (v9 §4) | Source | Class | Note |
|---|---|---|---|---|
| 1 | Responsibility separation — §1 boundaries hold | module layout; no adjudication path anywhere | I | no model, no arbitration call exists |
| 2 | Containment — every record in exactly one chat | `records.py` `__post_init__`; `engine.py` `__init__` | I | |
| 3 | Character ownership — every private record owned by exactly one character | `records.py` `__post_init__` | I | |
| 4 | Association — a private record names its owner's perspective | `Record.owner_character`, `attribution` | I | |
| 5 | Access — private retrievable only by its owner | `retrieval.perspective`; `_eligible_truth` excludes owned records | I | |
| 6 | Applicability — domain + target, never ownership or knowledge | `Record.domain/target`; `_eligible_truth` | I | |
| 7 | Truth-aptness — only world/scenario records hold truth | `Record.truth_apt`; `_accepts` | I | |
| 8 | No adjudication — no true/false marking | no verdict field or method exists | I | |
| 9 | Acquisition — accounts with attributable origin and evidence | `_validate_evidence` | I | C-3 table applied |
| 10 | Disclosure — shared only by an actual event | `disclose` | I | |
| 11 | Exposure limits acquisition | `disclose` creates accounts only for recipients; `told` needs an event | I | |
| 12 | Provenance without bypass — descriptor, never the source | `disclose` descriptor `{speaker, record_class, event, turn}` | I | |
| 13 | Conflict detection ≠ resolution | `assert_incompatible` + states; no resolution operation | I | |
| 14 | No truth mutation by non-authoritative records | only `promote`/`supersede_authority` write truth classes | I | |
| 15 | Belief independence — not created/replaced/revoked by information | no path from acquire/disclose/promote into belief; `retract` refuses beliefs | I | |
| 16 | Promotion is explicit | `promote` requires an authorised actor and a basis | I | |
| 17 | Promotion changes authority, not ownership | truth records have `owner_character is None`; sources untouched | I | |
| 18 | Promotion confers no knowledge | `promote` creates no account | I | |
| 19 | Promotion compels no belief | `promote` touches no belief record | I | |
| 20 | History is preserved — nothing erased | no delete; `SUPERSEDED`/`REVOKED` retain records | I | |
| 21 | Authority, not history | status changes only; the disagreement survives | I | |
| 22 | Creation is open | `claim`, `acquire`, `author_ambient`(blocked) | I | authorship gap in §5 |
| 23 | Recording is system-governed | provenance blocks on every class | I | |
| 24 | Promotion authority separate; explicit, domain-scoped, recorded, revocable | `AuthorityRegistry` + `require` | I | per-operation split absent (§5) |
| 25 | A statement is a claim | `claim` never creates an account | I | |
| 26 | Kind is declared, not judged; missing → `Unspecified claim` | `claim` records the declared kind verbatim | I | |
| 27 | Creation is not promotion | no implicit promotion path | I | |
| 28 | Basis requirement — ≥1 typed entry, empty rejected | `validate_basis` | I | |
| 29 | Basis retraction does not undo promotion | no code path reverses a promotion on retraction | I | |
| 30 | Declarations preserve content and author | `BasisEntry` requires content, actor, authority | I | |
| 31 | Promotion never extends reach | `chat_id` on every record; citations must be in-chat | I | |
| 32 | Reach is still not knowledge | reading canon creates no account | I | |
| 33 | Citations resolve within the chat and access rules | `promote` citation loop (`chat_id`, owner check) | I | |
| 34 | No relevance or coverage test | no sufficiency check anywhere | I | |
| 35 | Authority is per domain | grants keyed by domain tuple | I | |
| 36 | World scope is an authority domain, not a repository | per-chat engine; no global canon | I | |
| 37 | Revocation and supersession explicit, human-controlled | `revoke_promotion`, `supersede_authority` | I | |
| 38 | Authority changes never rewrite character state | both operations touch truth records only | I | |
| 39 | One authoring authority, two operations | `promote` present; `author_ambient` **absent** | **P** | blocked, §5 |
| 40 | Shared governance is not shared content | no shared store, no cross-chat path | I | |
| 41 | Equal-authority conflict is a state, not an error | `_assemble`: `none`/`settled`/`unresolved_relation`/`conflict`/`multiple` | I | |
| 42 | Three operations, three meanings | `correct`, `supersede_authority`, `revoke_promotion` | I | |
| 43 | Authority changes human-controlled | `_require_human` (`holder_kind == HUMAN`) | I | |
| 44 | A human edit changes authority, never history | supersession keeps both records | I | |
| 45 | Transfer is outside MEMORY-09 | no transfer symbol exists | I | structural absence |
| 46 | Disclosure needs no authorisation | `disclose` has no authority gate | I | |
| 47 | Narration cannot disclose for a character | `narration_report` reports only; no create path | I | |
| 48 | Disclosure is a claim, not truth | recipient account `via="told"`, never truth-apt | I | |
| 49 | Irreversibility is historical | retracting a source leaves recipient accounts | I | |
| 50 | Two questions, two fields | `via` and `origin_class` separate | I | |
| 51 | Own accounts create no events | `origin_class="own_account"` forces `via="none"`, exempt before the table | I | |
| 52 | A stated source is content, not history | descriptor is content; validator never parses text | I | |
| 53 | Beliefs are declared, not inferred | `declare_belief` + `Record` validation | I | |
| 54 | Narration may report a belief, never install one | no install path outside `declare_belief` | I | |
| 55 | Belief revision is append-only | `declare_belief` supersedes, retains; `revoke_belief` appends | I | |
| 56 | A belief is never promoted as truth | citation guard + `belief_facts` view + kind filter | I | residual R-C unchanged |
| 57 | Disclosure carries the belief, not its content | `disclose` writes `"<name> believes: …"` with `via="told"` | I | |

---

## 7. The 118 conformance requirements

| ID | Requirement (map wording) | Source | Class | Note |
|---|---|---|---|---|
| A1 | Memory exists; no accounts, no declaration → no character knows it | `promote` creates no account; `_eligible_truth` is a canon read | I | |
| A2a | Account; memory later retracted → character still holds what it acquired | `retract` touches one record | I | |
| A2b | Same character, same turn → no automatic knowledge of the retraction | no notice account is created | I | |
| A3 | Active memory; no accounts, no declaration → nobody knows it | retrieval is owner-scoped for private records | I | |
| A4 | Private record not consulted for world truth; non-truth-apt | `truth_apt`; canon classes only | I | |
| B1 | Declared over 3 present characters → 3 coverage accounts | `declare_available` creates one account per audience member | I | |
| B2 | Newcomer excluded by an earlier declaration → no account | audience snapshot at declaration | I | |
| B3 | Declaration targets an event-kind memory → refused | no guard written: v9 carries no event/definition distinction | **NA** | map status **X**, unresolved; not invented |
| B4 | Extending an existing declaration → refused | duplicate-target check | I | |
| B5 | Memory created after a declaration → not covered | coverage list frozen at declaration | I | |
| B6 | Declaration withdrawn → future coverage stops, existing intact | `withdraw_declaration`; `narrator_grant` refuses a withdrawn declaration | I | |
| B7a | Authorised authority declares → accepted with authority and basis recorded | `require` + declaration record | I | |
| B7b | Knowing character without authority declares → refused | `require` | I | |
| B7c | Authority revoked → future declarations refused, existing unaffected | `require` reads live grants only | I | |
| B8 | Declaration targets a private record → refused | `owner_character is not None` check | I | |
| B9 | Broad rule would include a private record → refused; explicit lists only | audience must be a non-empty explicit list | I | |
| C1 | Sable told by Alaric; Alaric later revoked → Sable still knows | recipient account independent of source status | I | |
| C2 | Inference with derivation basis → account with basis | `inferred` requires owned, existing basis | I | |
| C3 | Restricted fact granted to Sable → knows; authoriser recorded | `narrator_grant` + declaration actor | I | |
| C4 | Read account → document/asset cited | `read` requires `source` and event | I | |
| C5 | Wrong evidence kind → refused at write | `_validate_evidence` | I | |
| C6 | Recollection contradicting the world record → recorded as recollection, world unchanged | `claim` kind `recollection` vs `via="witnessed"` account | I | |
| C7 | Own account without an event → `own_account`/`none`, no event | exemption before the table | I | |
| C8 | Source named inside a recollection → content, no event | no event is manufactured; text never parsed | I | |
| C9 | Recollection, testimony and witnessing kept apart | `origin_class` + `via` | I | |
| D1 | Belief held; fact later acquired → belief persists | acquire touches no belief | I | |
| D2 | Belief contradicts world truth → truth unchanged | separate records; no write path | I | |
| D3 | Two characters, contradictory beliefs → both coexist | owner-scoped belief records | I | |
| D4 | Evidence presented to a believer → no automatic revocation | no such path | I | |
| D5 | Belief revoked explicitly → new record, prior retained | `revoke_belief` | I | |
| D6 | Account for P **and** belief against P → both coexist | different classes, no interference | I | |
| D7 | Two contradictory private records on one world fact → both persist, neither visible | ownership at `perspective` | I | |
| D8 | Belief created only by its owner's model | `declare_belief` + constructor | I | |
| D9 | No one else may install a belief | `declarer != character` refused; constructor too | I | |
| D10 | A belief is never promoted as truth | ascription kind + citation guard + `belief_facts` | I | |
| E1a | Retraction → only the memory's truth status changes | `retract` | I | |
| E1b | Accounts held by knowers unchanged | no cascade | I | |
| E1c | Retraction turn → no notice account | no account created | I | |
| E2 | Supersession → both versions retained | `supersede_authority` | I | |
| E3 | Private record superseded/retracted → domain truth unchanged | classes are disjoint | I | |
| F1 | No independently maintained audience projection; knowledge reconstructable | `accounts_of`, `history` | I | |
| F1b | Reconstruct vs any projection identical | no projection exists to diverge | I | |
| F2 | Accounts and beliefs carry no truth/status field | `Record` has no truth field; `status` is lifecycle | I | |
| F3 | Mutation attempts → no prior record edited or deleted | only `status`/`superseded_by` transition | I | |
| F4 | Unknown retrieval context → empty result | `dispatch` unknown branch; `unknown_view` refusal | I | |
| F5 | Empty eligible set → explicit abstention | `status: "none"`, empty `propositions` | I | |
| F6 | Every record's evidence resolves per origin | `_validate_evidence`, per `via` | I | |
| F7 | Replay same ordered records → identical answers | deterministic ordering by `proposition_id` then `record_id`; **no replay facility** | **NA** | map status **R**; no execution, so not assessable |
| F8 | Reconstruction including private records | `snapshot`, `history`, `as_dict` | I | |
| G1 | Newcomer asked about an unwitnessed event → abstains | no account ⇒ nothing to retrieve | I | |
| G2 | Newcomer knows a declared definition, naming the declaration | coverage account descriptor names `declaration` | I | |
| G3 | Newcomer where others hold private memories → abstains on all | owner-scoped pool | I | |
| H1 | Private record recorded with provenance, no verdict | `claim`/`acquire` | I | |
| H2 | Other character retrieves, every context → abstains | `perspective` and canon exclude owned records | I | |
| H3 | World-truth query unaffected by private records | `_eligible_truth` requires `owner_character is None` | I | |
| H4 | A discloses to B → B gets a new own record with provenance | `disclose` | I | |
| H5 | Unauthorised share attempt → refused | no share operation exists | I | |
| H6 | A's private record contradicts world truth → both stand | disjoint classes | I | |
| H7 | Narration cannot manufacture or reproduce private content | `narration_report` returns only the reported event | I | |
| H8 | Narration may report an established disclosure | `narration_report` | I | |
| H9 | Re-disclosure from own account without source access | `disclose` accepts any record the discloser owns | I | |
| H10 | Source retracted after disclosure → recipient account persists | no cascade | I | |
| I1 | Every record has one containing chat and, if private, one owner | `chat_id`, `owner_character` | I | |
| I2 | Ownership is not transferable; disclosure creates a new record | no reassignment path | I | |
| I3 | Containment and ownership move independently | separate fields, no derivation | I | |
| J1 | Public statement → claim; world unchanged | `claim` | I | |
| J2 | Recollection vs witnessed evidence → both exist, origins differ | `via` distinction | I | |
| J3 | Voluntary disclosure → claim event with content and recipients | `disclose` | I | |
| J4 | No actual disclosure event → no accounts | accounts only from `disclose` | I | |
| J5 | Exposure governs acquisition | recipient loop | I | |
| J6 | Group membership confers no knowledge | no group-membership read on any acquisition path | I | |
| J7 | Recipient account has source, turn, content | descriptor + provenance | I | |
| J8 | Origin descriptor resolves to nothing | `{speaker, record_class, event, turn}` — no record id | I | |
| J9 | Conflicting accounts → both retrievable by owners, neither promoted | no promotion from accounts | I | |
| J10 | Disclosure and belief → no revision unless the model changes it | no belief write in `disclose` | I | |
| J11 | Disclosure and truth → no truth mutation | `disclose` writes accounts only | I | |
| J12 | Detection without arbitration | states + `relation`; no winner selection | I | |
| J13 | Equal-authority conflict retained with provenance; group unresolved | `assert_incompatible` + `conflict` | I | |
| J14 | Retrieval presents the conflict as unresolved | `_assemble`; never selects a record as settled | I | |
| J15 | Correction does not resolve an authority conflict | `correct` leaves groups open | I | |
| J16 | Human edit resolves; history survives | `supersede_authority` | I | |
| K1 | Promotion is never automatic | no automatic path | I | |
| K2 | Promotion creates an authoritative record with applicability and provenance | `promote` | I | |
| K3 | Authority, not mutation — source unchanged | citation is read-only | I | |
| K4 | No conferred knowledge | no account created | I | |
| K5 | No compelled belief | no belief written | I | |
| K6 | Ownership and containment untouched | truth record is chat-scoped, ownerless | I | |
| K7 | Promotion settles authority | one live record on the identity → `settled` | I | |
| K8 | Authorised actor only; authority recorded | `require`; provenance `authority` | I | |
| K9 | Operation recorded — actor, authority, time, domain, applicability, source | `_promotions` entry | I | |
| K10 | Later promotions enter a recorded chain; latest governs | displacement only via human-gated supersession | **P** | RAT-1 precedence: tier 1 O-10b outranks "latest wins" |
| K11 | Promotion preserves history | records retained, retrievable by owners | I | |
| K12 | Two legitimate answers in one state | canon and perspective are separate surfaces | I | |
| K13 | Basis kind valid and recorded | `BasisEntry` + `validate_basis` | I | |
| K14 | Citations resolve in-chat and comply with access rules | citation loop | I | |
| K15 | Sufficiency not assessed | no relevance test | I | |
| K16 | Basis retracted after promotion → promotion stands | no undo path | I | |
| K17 | Declaration preserves content and author | `BasisEntry` fields | I | |
| K18 | Promotion stays inside its chat | `chat_id` checks | I | |
| K19 | Basis is a set, mixed kinds allowed, empty rejected | `validate_basis` | I | |
| K20 | No relevance or coverage test | absent by design | I | |
| K21 | Authority per domain | domain-keyed grants | I | |
| K22 | Revocation explicit | `revoke_promotion`, `revoke_grant` | I | |
| K23 | Supersession preserves both states | `supersede_authority` | I | |
| K24 | Authority change leaves characters untouched | truth records only | I | |
| K25 | Ambient authorship is not promotion | `author_ambient` absent | **A** | blocked, §5 |
| K26 | No unilateral authority change | `_require_human` | I | |
| K27 | No transfer operation | no symbol exists | I | |
| L1 | Human contribution recorded with provenance | `claim` | I | |
| L2 | Model contribution recorded with content and attribution | `claim` | I | |
| L3 | Statement implies no experience | no account created from a claim | I | |
| L4 | Kind recorded as declared, no reclassification | `claim` verbatim | I | |
| L5 | Exposure governs acquisition | `disclose` | I | |
| L6 | Creation ≠ authority | `claim` non-truth-apt | I | |
| L7 | Standing grant to a model works as declared; explicit, scoped, recorded, revocable | `grant` + `require` | I | operations field absent (§5) |
| L8 | Generation grants no authority | no authorship operation exists | **A** | blocked, §5 |
| L9 | Missing kind → `Unspecified claim`, nothing invented | `claim` default | I | |
| L10 | Revoking a standing grant stops future promotion; records intact | `revoke_grant` + `require` | I | |

**Totals:** rules **56 I · 1 P (39) · 0 A · 0 C · 0 NA**; requirements **113 I · 1 P (K10) · 2 A (K25, L8) ·
0 C · 2 NA (B3, F7)** — 118 rows. These classifications are **source inspection**, not executed
verification.

---

## 8. Assumptions, and why

| # | Assumption | Justification |
|---|---|---|
| 1 | Refusal envelope key is `reason`, not `code` | REV 7 §5.1 gives the envelope literally as `{"status": "refused", "reason": "no_domain_authority"}`; REV 12 §5 and REV 13 §2.3 use `code` for the same values in illustrative JSON. One fixed shape is load-bearing (S-9), so the earlier literal envelope was treated as the shape and the later word as shorthand for its value. **Flagged** — if `code` is intended, one key changes in `_refusal`. |
| 2 | Top-level `status` is `conflict` when any entry is `conflict` and several entries exist | REV 6 §4 names `conflict` for "open group spanning them"; REV 8/9 bind `relation == "asserted"` ⟺ some entry `conflict`. Using `multiple` alongside `relation: "asserted"` would contradict invariant 11, so the stricter reading was taken. |
| 3 | `propositions()` accepts `viewer="authority"` and refuses anything else | REV 4 §7 lists `not_a_perspective_surface` as a refusal; a surface must exist to produce it. No other behaviour depends on the parameter. |
| 4 | `revoke_belief` takes the revocation text from the caller | D5 requires a new record carrying the revocation; generating that text would be the engine inventing a character's statement (rules 8, 53). |
| 5 | `disclose` writes an `event` descriptor | A `told` account's evidence is a disclosure event (C-3); `disclose` *is* that event, so recording its kind and turn is recording, not manufacturing. |
| 6 | `promote` returns a dict rather than a `Record` | REV 3 §7.2 row 11 requires `coexisting_propositions` "in the result" and REV 4 §5.3 fixes its shape; a record object cannot carry it. |
| 7 | `assert_incompatible` derives `domain` from the cited records | Its contract signature has no `domain`; the two records must agree, and a mismatch is refused. |
| 8 | One diagnostic code, `proposition_unregistered`, is invented for the orphan case | REV 5 §4 item 3 requires an orphan be reported, not silently omitted, but names no code. **Flagged** as the only invented diagnostic name; it names an integrity fault, never a memory state, and appears only when the registry and record store have diverged. |
| 9 | `_proposition_records` and `_proposition_state` are inlined in `_assemble` | Their contract is a projection of `S`; inlining them keeps S-2 checkable in one function. No behaviour differs. |
| 10 | `possible_duplicates` reads the registry without a version bracket | It joins nothing and reads one store, so no cross-store blend is possible; it is not one of the three snapshot views. Staleness is possible under concurrency, which is untested. **Flagged.** |

---

## 9. Unresolved ambiguities and residual risks

| # | Item | Status |
|---|---|---|
| 1 | `Grant.operations` / `require(..., operation)` | **Blocked**; `author_ambient` absent (§5) |
| 2 | `author_ambient` | **Absent**; requires §5's unblock |
| 3 | B3 (event vs definition) | **NA**; no guard written, no resolution invented |
| 4 | F7 replay | **NA**; no facility; determinism of ordering only |
| 5 | Refusal envelope key (`reason` vs `code`) | Assumption 1, flagged |
| 6 | `proposition_unregistered` code name | Assumption 8, flagged |
| 7 | `assert_incompatible` reads the two cited records before `require`, because its domain is a property of those records | Structural consequence of a signature that carries no domain; limited to the two named records, and every other protected read is post-gate |
| 8 | Concurrency | **Not demonstrated.** `_capture`'s version test is sound only because `_mutation` publishes atomically in a single-threaded store. No concurrent scenario has been run, and Class B is unreachable in this implementation. |
| 9 | Behaviour of `supersede`'s write pair under interleaving | Specified (one mutation, one version step); not observed |

---

## 10. Confirmations required by §7

1. **`authority.py` was not modified.** sha256 `07dfbeb219221190…`, identical to the R2 baseline. Its
   `require(domain, actor)` remains two-parameter; §5 records the resulting block.
2. **No transfer, destination-write, sharing or cross-chat carryover was introduced.** No such symbol
   exists; `grep -rn "transfer(" eval/mem1/m09/*.py` matches one comment recording R2's removal.
3. **No test, probe, suite, self-test, compilation or linting was executed.** Nothing was imported, run,
   compiled, linted or benchmarked. `py_compile` was deliberately **not** run, because §6 prohibits
   compilation; the only verification performed is the source inspection this report records.
4. **Source-inspection evidence vs executed verification.** Every classification above is a source claim.
   No requirement is reported as verified because code exists: where a property depends on runtime
   behaviour (F7, Class B reachability, concurrency), it is **NA** and says so.

---

## 11. Known behavioural consequences of the accepted contract

These are consequences, not defects, and are recorded so the owner is not surprised:

| Change | Caller-visible effect |
|---|---|
| `promote` returns a dict | Callers read `result["record"]` |
| `authoritative` requires `view` | No default view; omission is a `TypeError`, not a silent world answer |
| `world()` returns the per-proposition envelope | A caller that read a flat `records` list must read `propositions` |
| `world()` status `multiple` | Several unrelated identities are no longer flattened into one answer |
| `propose`-style supersession | `supersede_authority` returns the promotion dict plus `superseded` |
| Refusals are envelopes, not exceptions | `authoritative`, `propositions`, `possible_duplicates` return `{"status": "refused", …}` |
| `acquire` requires evidence | `told`, `read`, `witnessed`, `inferred`, `verified`, `narrator_grant` each need their descriptor keys |

**Sequence:** implementation complete → this report → owner review → exact-rule source review →
separate testing authorisation. This report requests neither the review's outcome nor the testing gate.
