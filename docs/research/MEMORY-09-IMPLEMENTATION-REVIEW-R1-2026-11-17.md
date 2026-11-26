# MEMORY-09 — Implementation review R1 (Gate 2)

**Status:** **SOURCE INSPECTION ONLY. 0 tests executed. 0 probes run. No runtime behaviour claimed.**
Recorded **2026-11-17 08:31 CST (−0600)** = **14:31 UTC**.

**Authorised by:** Gate 2 (owner, 2026-11-17), implementation only.
**Reviewed artifact:** `eval/mem1/m09/` — 6 modules, **886 lines**.
**Governing baselines:** v9 (57 rules, tier 3 under O-6a) · conformance map (118 requirements, O-6b) ·
D-4 · O-13 · all confirmed decisions.

**What this review is:** a mapping from rule and requirement to source location, produced by reading
the code that was just written.
**What this review is not:** evidence that anything works. No import, no execution, no probe, no
suite. `py_compile` was run and passed; that proves syntax, not behaviour.

---

## 1. Inventory

| File | Lines | Purpose |
|---|---|---|
| `m09/__init__.py` | 31 | Exports; declares `EXPOSES_TRANSFER = False` |
| `m09/errors.py` | 74 | One error class per refused condition, each naming its rule |
| `m09/records.py` | 154 | Four record classes + structural validation |
| `m09/authority.py` | 156 | Domains, holders, standing grants, typed basis |
| `m09/engine.py` | 391 | Claims, accounts, beliefs, promotion, disclosure, lifecycle, declarations |
| `m09/retrieval.py` | 80 | Fail-closed retrieval contexts; conflict presentation |

**Relationship to existing code:** the G1 experiment harness (`model.py`, `retrieve.py`,
`assemble.py`, `conversation.py`, 6 152 lines total) is **untouched**. It builds records for the
ablation arms and conflates containment, ownership and domain in one `owner_id`. The v9 engine is a
separate package and shares no state with it. This review covers `m09/` only.

---

## 2. Rule mapping (57 rules)

| Rule | Evidence | Verdict |
|---|---|---|
| 1 Responsibility separation | `engine.Engine` operations split by actor; no combined adjudicator | Implemented |
| 2 Containment | `records.Record.chat_id` required; engine constructed per chat | Implemented |
| 3 Character ownership | `Record.owner_character`; `records.__post_init__` refuses private without one | Implemented |
| 4 Association | `owner_character` + `attribution`; both distinct from `chat_id` | Implemented |
| 5 Access | `retrieval.perspective` filters by `owner_character`; `disclose` creates new records | Implemented |
| 6 Applicability | `Record.domain` + `Record.target`, separate from ownership | Implemented |
| 7 Truth-aptness | `Record.truth_apt`; only `world_truth`/`scenario_truth` | Implemented |
| 8 No adjudication | No field or operation marks a claim/account/belief true or false | Implemented (by absence) |
| 9 Acquisition | `Engine.acquire` records origin, via, span, turn | Implemented |
| 10 Disclosure | `Engine.disclose` requires an explicit event with recipients | Implemented |
| 11 Exposure limits acquisition | Accounts created only for the `recipients` list | Implemented |
| 12 Provenance without bypass | Recipient `descriptor` carries speaker/class/turn, never a record id | Implemented |
| 13 Conflict detection ≠ resolution | `Engine.authoritative` returns `conflict` state; no resolution path for accounts | **Partial** — conflicts are represented for authoritative records only |
| 14 No truth mutation by non-authoritative records | Only `promote` creates truth-apt records | Implemented |
| 15 Belief independence | No operation couples belief to acquisition, disclosure or promotion | Implemented |
| 16 Promotion is explicit | `promote` is a named call with actor, basis, domain | Implemented |
| 17 Promotion changes authority, not ownership | Truth record has `owner_character=None`; source untouched | Implemented |
| 18 Promotion confers no knowledge | `promote` creates no accounts | Implemented |
| 19 Promotion compels no belief | No belief written by `promote`; rejection needs no handling | Implemented |
| 20 History is preserved | Supersession sets `superseded_by`; records never deleted | Implemented |
| 21 Authority, not history | `authoritative()` answers truth; `history()` still returns all versions | Implemented |
| 22 Creation is open | `claim` accepts any contributor string and kind | Implemented |
| 23 Recording is system-governed | Provenance dict on every record; contributor/kind recorded | Implemented |
| 24 Promotion authority is separate | `AuthorityRegistry.require`; model grants explicit and revocable | Implemented |
| 25 A statement is a claim | `claim` creates `record_class='claim'` only | Implemented |
| 26 Kind is declared, not judged | Missing kind recorded as `Unspecified claim`; no inference path | Implemented |
| 27 Creation is not promotion | `claim`/`acquire` never create truth-apt records | Implemented |
| 28 Basis requirement | `validate_basis`; empty set refused; mixed kinds accepted | Implemented |
| 29 Basis retraction does not undo | No code path reverses a promotion when a cited record is retracted | Implemented (by absence) |
| 30 Declarations preserve content and author | `BasisEntry` requires content, actor, actor_authority | Implemented |
| 31 Promotion never extends reach | Single-chat engine; citations refused across chats | Implemented |
| 32 Reach is not knowledge | Promotion creates no account | Implemented |
| 33 Citations must resolve within the chat | `promote` resolves each citation; foreign/private refused | Implemented |
| 34 No relevance or coverage test | `validate_basis` checks structure only; no scoring | Implemented |
| 35 Authority is per domain | `AuthorityRegistry` keyed by `(kind, id)` | Implemented |
| 36 World scope is an authority domain | Per-chat `Engine` + per-domain grants; no repository | Implemented |
| 37 Revocation and supersession are explicit | `revoke_promotion`, `supersede_authority` | Implemented |
| 38 Authority changes never rewrite character state | Neither operation touches accounts or beliefs | Implemented |
| 39 One authoring authority, two operations | Promotion present; **ambient authorship absent** | **Partial** — one of the two operations has no implementation |
| 40 Shared governance is not shared content | No sharing or cross-chat path exists | Implemented (by absence) |
| 41 Equal-authority conflict is a state | Two active records for one domain+target produce `status='conflict'` | Implemented |
| 42 Three operations, three meanings | `correct`, `supersede_authority`, `revoke_promotion` are distinct | Implemented |
| 43 Authority changes are human-controlled | `_require_human` refuses model-held grants | Implemented |
| 44 A human edit changes authority, never history | `corrected_by` retains prior text; superseded records retained | Implemented |
| 45 Transfer is outside MEMORY-09 | `EXPOSES_TRANSFER=False`; `transfer()` raises; no operation exists | Implemented |
| 46 Disclosure needs no authorisation | `disclose` requires only ownership of the source | Implemented |
| 47 Narration cannot disclose for a character | `narration_report` requires an existing event id; no create path | Implemented |
| 48 Disclosure is a claim, not truth | Recipient account is non-truth-apt; source not exposed | Implemented |
| 49 Irreversibility is historical | `retract` leaves disclosure accounts; source removal does not cascade | Implemented |
| 50 Two questions, two fields | `origin_class` and `via` are separate validated fields | Implemented |
| 51 Own accounts create no events | `origin_class='own_account'` forces `via='none'`; no event object created | Implemented |
| 52 A stated source is content | `descriptor` is a dict; nothing resolves it to a record | Implemented |
| 53 Beliefs are declared, not inferred | `declare_belief` requires `role='character_model'` and declarer==owner | Implemented |
| 54 Narration may report, never install | No narration path reaches `declare_belief` | Implemented |
| 55 Belief revision is append-only | Prior belief set `superseded`, retained and retrievable | Implemented |
| 56 A belief is never promoted as truth | Promotion of a cited belief records a fact about the character; no content authority derived | **Partial** — nothing structurally distinguishes "Sable believes P" from "P" in `text` |
| 57 Disclosure carries the belief | Belief disclosure writes `"<name> believes: <P>"` as an account | Implemented |

**Rule totals:** 51 implemented · 3 partial (13, 39, 56) · 3 implemented by absence (8, 29, 40) counted
within the 51 · **0 with no evidence**.

---

## 3. Requirement mapping (118 requirements)

Verdicts: **I** implemented (source present) · **P** partial · **A** absent · **X** contradictory · **N** not assessable.
All verdicts are **source inspection, not execution**.

| ID | V | Evidence / gap |
|---|---|---|
| A1 | I | `perspective()` pool is owned records only; a world record is never returned |
| A2a | I | `retract` touches one record's status; accounts are separate records |
| A2b | I | No notice account is created by `retract` |
| A3 | I | Truth records have `owner_character=None`; unreachable from `perspective()` |
| A4 | I | `truth_apt` False for all non-truth classes |
| B1 | I | `declare_available` creates one account per listed character |
| B2 | I | Coverage is a snapshot list; later characters are never added |
| B3 | **A** | No event/definition distinction exists on a record, so no guard fires. Recorded as an omission, not faked |
| B4 | I | Second active declaration for one target raises `OperationError` |
| B5 | I | No path covers records created after the declaration |
| B6 | I | `withdraw_declaration` sets `withdrawn_turn`; accounts remain |
| B7a | I | `authority.require` accepts the domain holder |
| B7b | I | Non-holder raises `AuthorityError` |
| B7c | I | `revoke_grant` stops future `require`; existing coverage intact |
| B8 | I | Private source raises `DisclosureError` |
| B9 | I | `audience` must be a non-empty explicit list; no pattern or group path |
| C1 | I | Recipient account is an independent record; revoking the teller's does not cascade |
| C2 | **P** | `via='inferred'` recorded, but a derivation basis is not stored as a typed structure |
| C3 | **P** | No restricted-scope concept in the v9 engine; authoriser is recorded only via declaration |
| C4 | **P** | Descriptor can cite a document; no evidence-kind validator confirms it |
| C5 | **A** | No write-time evidence-kind validation |
| C6 | I | Claim (`recollection`) and truth record coexist; neither mutates the other |
| C7 | I | `origin_class='own_account'` ⇒ `via='none'`, no event |
| C8 | I | A named source is text/descriptor only; no event is manufactured |
| C9 | I | `origin_class` + `via` distinguish the three; no collapse |
| D1 | I | Acquiring an account does not touch beliefs |
| D2 | I | Belief is non-truth-apt; truth unchanged |
| D3 | I | Beliefs keyed by `owner_character`; no cross-character path |
| D4 | I | No operation revokes a belief on evidence presentation |
| D5 | I | New belief record; prior retained as superseded |
| D6 | I | Account and counter-belief coexist; both retrievable |
| D7 | I | `perspective()` filters by owner; world untouched |
| D8 | I | Belief created only by `declare_belief`; `claim` cannot reach it |
| D9 | I | Foreign declarer/role refused in `records.__post_init__` |
| D10 | **P** | No content authority is derived, but `text` alone cannot distinguish a belief claim from content |
| E1a | I | `retract` mutates `status` only |
| E1b | I | Accounts unchanged |
| E1c | I | No notice account |
| E2 | I | Supersession retains both; `history()` returns all |
| E3 | I | Private retraction cannot reach a truth record |
| F1 | I | Knowledge read from accounts; no separate projection object exists |
| F1b | I | Nothing to diverge from: no projection is maintained |
| F2 | I | No truth/status field on accounts or beliefs |
| F3 | I | No delete/overwrite operation; corrections retain prior text |
| F4 | I | `dispatch` returns an empty answer for an unknown context |
| F5 | I | Explicit `abstain` statuses with reasons |
| F6 | **P** | Provenance stored per record, but resolution is not validated per origin kind |
| F7 | **P** | No replay function exists, so determinism is unverified and untested |
| F8 | I | Private records retained with `owner_character` and attribution |
| G1 | I | No account for an unwitnessed event |
| G2 | I | Declaration coverage account is retrievable and names the declaration |
| G3 | I | `perspective()` excludes others' private records in every context |
| H1 | I | Claim recorded; no verdict field |
| H2 | I | Abstains (no eligible records) rather than leaking |
| H3 | I | Truth query reads only truth-apt records |
| H4 | **P** | New own record with descriptor; it names speaker/class/turn, not a source id (J8 requires it not resolve) |
| H5 | I | `disclose` requires ownership of the source; no share operation exists |
| H6 | I | Both exist; world unchanged |
| H7 | I | `narration_report` cannot create an event or return private content |
| H8 | I | Reporting an established event is permitted |
| H9 | I | Discloser may disclose their own account without the source |
| H10 | I | Recipient accounts persist after `retract` of the source |
| I1 | I | `chat_id` and `owner_character` both required and independent |
| I2 | I | No reassignment path; disclosure copies into a new record |
| I3 | I | No operation writes both fields together |
| J1 | I | Claim created; truth untouched |
| J2 | I | Both retained with different classes/origins |
| J3 | I | `disclose` records content and recipients |
| J4 | I | No accounts without an explicit `disclose` call |
| J5 | I | Recipients list governs |
| J6 | I | No group concept exists; nothing can confer |
| J7 | I | Own account with descriptor, turn and content |
| J8 | I | Descriptor is a plain dict; no identifier to resolve |
| J9 | I | Both retrievable by their owners; no promotion path |
| J10 | I | No system-created belief revision |
| J11 | I | Disclosure writes no truth record |
| J12 | **P** | World conflicts flagged; account-level conflicts are not detected or flagged |
| J13 | **P** | Conflict expressed as a status + record list, not a persisted conflict-group object |
| J14 | I | `retrieval.world` returns `conflict`, not a pick |
| J15 | I | `correct` leaves both active records untouched in authority terms |
| J16 | I | `supersede_authority` + `corrected_by`; originals retained |
| K1 | I | No implicit promotion path |
| K2 | I | Domain + target + provenance + `promotion_id` |
| K3 | I | Source unchanged; truth query returns the promoted record |
| K4 | I | No accounts created |
| K5 | I | No belief created; no failure path for rejection |
| K6 | I | Owner and chat unchanged |
| K7 | **P** | Settles authority only when `supersede=True`; without it a conflict is the correct outcome, so "settles" is operator-dependent |
| K8 | I | `authority.require`; actor recorded in `promotions` |
| K9 | I | Actor, authority, turn, domain, target, record, rationale, basis |
| K10 | I | `supersedes` links and `promotions[*]['superseded']` chain |
| K11 | I | Prior authoritative records retained and retrievable |
| K12 | **P** | Both queries exist; no single call returns both answers in one state |
| K13 | I | `BASIS_KINDS` enforced |
| K14 | I | Existence, chat match and access checked; id-less citation refused |
| K15 | I | No sufficiency scoring anywhere |
| K16 | I | No auto-undo on basis retraction |
| K17 | I | `BasisEntry` requires content + actor + authority |
| K18 | I | Cross-chat citation refused; separate engine per chat |
| K19 | I | ≥1 entry; mixed accepted; empty refused |
| K20 | I | No relevance/coverage requirement |
| K21 | I | Grants keyed by domain; cross-domain refused |
| K22 | I | Only `revoke_promotion` / `supersede_authority` change authority |
| K23 | I | New state governs; prior retained with provenance |
| K24 | I | Neither operation touches accounts or beliefs |
| K25 | **A** | Ambient authorship has no implementation, so it cannot be shown not to promote |
| K26 | I | `_require_human` refuses model grants |
| K27 | I | `EXPOSES_TRANSFER=False`; no transfer/share/carryover operation |
| L1 | I | `claim` records contributor and kind |
| L2 | I | Model output recorded as a claim with attribution |
| L3 | I | Claim creation never creates an account |
| L4 | I | `contribution_kind` stored as declared; no reclassification |
| L5 | I | Recipients list governs acquisition |
| L6 | I | Scenario/world records are truth-apt only after `promote` |
| L7 | I | `grant()` + `require()`; unrecorded or unscoped grant cannot exist |
| L8 | I | No authorship path confers authority |
| L9 | I | Missing kind ⇒ `Unspecified claim` |
| L10 | I | `revoke_grant` stops future promotion; promoted records intact |

**Requirement totals:** **104 implemented · 11 partial · 3 absent · 0 contradictory · 0 not assessable.**

Partial: C2, C3, C4, D10, F6, F7, H4, J12, J13, K7, K12.
Absent: **B3, C5, K25**.

---

## 4. Implementation assumptions and their authority

Each assumption is a decision the code had to make that no ratified document settled. None was used
to close an open question on the owner's behalf.

| # | Assumption | Source of authority |
|---|---|---|
| 1 | The human operator is the default holder of every declared domain | K8 — *"default the human operator, per domain"* |
| 2 | Declarations target non-private records only | B8; the assumption is that "generally available" excludes private records |
| 3 | A human actor cannot cite a private record as basis | Rule 33 "complies with existing access rules"; a human is not a character owner |
| 4 | Two active records for the same domain+target are an equal-authority conflict | Rule 41; **no semantic incompatibility comparator exists** — the engine does not judge content |
| 5 | `correct` edits text in place and retains the prior text in `corrected_by` | Rules 42 + 44; supersession is the new-record path, correction is not |
| 6 | Belief disclosure renders `"<name> believes: <P>"` | Rule 57; wording is presentation, not semantics |
| 7 | Retrieval scoring is lexical term overlap | Not specified by v9; inherited from the existing experiment harness |
| 8 | Prior belief on revision is marked `superseded` yet remains retrievable | Rule 55 requires retention; status is not evidence of deletion |
| 9 | Disclosure accounts use `via='told'`; declaration coverage uses `via='narrator_grant'` | Rule 50 vocabulary, applied literally |
| 10 | "Incompatible" is not computed anywhere | Rule 34 / §9 forbid the system judging persuasion or relevance |

---

## 5. Deviations, omissions and unresolved questions

**Omissions (recorded, not hidden):**

1. **B3 absent** — v9 carries no event/definition distinction on a record. I removed an early guard
   that could never fire rather than ship a guard that looks like coverage. The requirement has no
   representation.
2. **C5 absent** — no write-time evidence-kind validator.
3. **K25 absent** — ambient authorship does not exist, so "ambient authorship is not promotion" is
   vacuously true rather than demonstrated.
4. **F7** — no replay function; determinism is neither implemented nor checkable.

**Partial implementations** — the 11 rows above. Each has a specific missing half, named in the table;
none is "mostly done".

**Contradictions:** none found between the implementation and a ratified rule.

**Unresolved questions, not resolved by assumption:**

- **B3's representation** — whether an event/definition distinction belongs on a record or elsewhere
  is an architecture question, not an implementation one.
- **K12's single-call form** — whether one query should return both the authoritative and the owner's
  answer is a design question the 57 rules do not settle.
- **Conflicts as objects** — J13 asks for a "conflict group". The implementation expresses the state
  as a query result. Whether a persisted group is required is unresolved.

**Documentation deviation, flagged not fixed:** v9's header still lists *"Gate 2 (implementation —
unanswered)"* under open inputs. Editing a ratified artifact was prohibited by the Gate 2
authorisation, so it was left untouched. It is stale as of this decision and needs an owner-authorised
correction.

---

## 6. Prohibited-scope confirmation

| Constraint | Result |
|---|---|
| `eval/mem1` only | Yes — 6 new files under `eval/mem1/m09/`; no file outside it written |
| No production Howl integration | None |
| r32 unchanged | `e88ef2d1697da1d9…` |
| Experiment plan v3 unchanged | 21 616 B, `6c4078bd…` |
| Amendment 1 unchanged | `6afd2a0d8c19a5cd…` |
| Ratified decision records unchanged | Only a new Gate 2 entry appended; prior entries preserved verbatim |
| No new architecture, rules, fields, requirements, probes | None added; 57 rules and 118 requirements unchanged |
| No transfer/destination-write/sharing/carryover | Absent; `transfer()` raises; `EXPOSES_TRANSFER=False` |
| No open question resolved by implementation | Three recorded as unresolved (§5) |
| No conformance probe or test executed | **0 executed** — `py_compile` only |
| Existing G1 harness untouched | Yes — `model.py`, `retrieve.py`, `assemble.py` and others unmodified |

---

## 7. Inspection versus execution

| Claim | Basis |
|---|---|
| "The code compiles" | **Executed** — `python -m py_compile eval/mem1/m09/*.py` → OK |
| "Requirement X is implemented" | **Source inspection** — the mapping in §3 |
| "Requirement X passes" | **Not claimed and not knowable.** No probe exists, none ran |
| "The engine preserves history" | **Source inspection** — the code paths that would erase it are absent |

**Every verdict in §2 and §3 is source-level.** A reviewer should treat "implemented" as *"a code path
that appears to satisfy the requirement exists"*, not as *"the requirement is met"*. The difference is
the entire reason testing is a separate gate.

---

## 8. What remains unauthorised

- Running any of the 118 probes.
- Running `selftest.py` or any suite in `eval/mem1`.
- Treating any verdict above as a result.
- Gate 3 — or whatever the owner names the testing stage — does not exist yet.
