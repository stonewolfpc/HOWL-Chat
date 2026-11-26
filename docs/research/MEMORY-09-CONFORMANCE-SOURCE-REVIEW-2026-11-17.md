# MEMORY-09 — Consolidated source-conformance review

**Status: SOURCE REVIEW ONLY. Nothing executed. No test, probe, suite, compilation, linting or runtime
verification was performed, and none is authorised.**

This document **supersedes the classification tables** of
`MEMORY-09-GATE2-IMPLEMENTATION-REPORT-2026-11-17.md` (§6, §7, totals),
`MEMORY-09-GATE2-REPORT-SUPPLEMENT-2026-11-17.md` (§3, §6) and the row-level claims in
`MEMORY-09-OWNER-CLARIFICATION-AUTHORITY-REQ17P-2026-11-17.md` §7. Those documents are retained unedited as
history; their non-classification content (delta inventory, assumptions, blocked-item record, verbatim
quotations) stands.

| Field | Value |
|---|---|
| Date | 2026-11-17 CST (−0600) |
| Base reviewed | `eval/mem1/m09/` at `authority.py f3445823…`, `engine.py 706d5f6c…`, `records.py f211e6ce…`, `errors.py 48125d7e…`, `retrieval.py cd2d380b…`, `__init__.py 06ccf8f6…` |
| Ratified base | 57 v9 rules; 118 map requirements; r32 tier 2; owner decisions tier 1 |
| Operative contract | REV 13 §2–§5 (`5b745628…`); REV 11 §3.2 S-1…S-11; REV 12 §2–§3 |
| Result | **Rules 56 I · 1 P**; **requirements 113 I · 3 P · 2 NA** |
| Executed | nothing |

---

## 1. Method, and what a classification means

Each row was derived by reading the cited symbol and asking whether the requirement's *expected* condition is
produced by the code, not whether a similarly named function exists.

| Class | Meaning |
|---|---|
| **I** | The code path required by the requirement is present and was read. **Not** a claim that the behaviour was observed |
| **P** | The condition is produced for part of the requirement's scope, or with a stated deviation |
| **A** | Absent |
| **C** | Contradicted |
| **NA** | Cannot be assessed from source — the contract leaves it open, or the requirement is about runtime behaviour |

Three classifications moved **up** since the first report (39, K25, L8), caused by implementing the authority
capability split. **Three moved down** (46, J3, K9) on this fresh pass. Everything else is unchanged. The
downgrades are the substance of this review: the earlier report was too generous, and a 57/57 result would
have hidden it.

| Row | Earlier | Now | Why |
|---|---|---|---|
| Rule 39 | P | I | `Grant.operations` + `require(…, operation)` + `author_ambient` |
| K25 | A | I | ambient authorship creates a `claim`, never truth |
| L8 | A | I | generation gated on `authorship`, which promotes nothing |
| **Rule 46** | I | **P** | recipient accounts and attribution are recorded, but the disclosure event itself has no record — see §4.1 |
| **J3** | I | **P** | same cause: no event record carrying content and recipients |
| **K9** | I | **P** | "time" is the turn counter; v9 §6 asks for a "timestamp with offset", which does not exist — see §4.2 |

Symbol shorthand: `E:` = `engine.py`, `R:` = `retrieval.py`, `Rec:` = `records.py`, `A:` = `authority.py`.

---

## 2. The 57 v9 rules

| # | Rule (v9 §4, abridged to its operative clause) | Source | Class |
|---|---|---|---|
| 1 | Responsibility separation — §1 boundaries hold; no AI arbitration | module layout; no model or arbiter call exists | I |
| 2 | Containment — one containing chat per record | `Rec:__post_init__`, `E:__init__` | I |
| 3 | Character ownership — one owner per private record | `Rec:__post_init__` (ownerless only for authored claims) | I |
| 4 | Association — a private record names its owner's perspective | `Rec.owner_character` | I |
| 5 | Access — owner-only for private records | `R:perspective`; `E:_eligible_truth` excludes owned records | I |
| 6 | Applicability — domain + target only | `Rec.domain/target`; `E:_eligible_truth` | I |
| 7 | Truth-aptness — only world/scenario hold truth | `Rec.truth_apt`; `E:_accepts` | I |
| 8 | No adjudication — no true/false marking | no verdict field or method | I |
| 9 | Acquisition — accounts with origin and evidence | `E:_validate_evidence` | I |
| 10 | Disclosure — shared only by an actual event | `E:disclose` is the only account-creating path with `via="told"` | I |
| 11 | Exposure limits acquisition | `E:disclose` recipient loop; `told` requires an event | I |
| 12 | Provenance without bypass — descriptor, never the source record | `E:disclose` writes `{speaker, record_class, event, turn}` | I |
| 13 | Conflict detection ≠ resolution | `E:assert_incompatible`; states; no resolution op | I |
| 14 | No truth mutation by non-authoritative records | only `promote`/`supersede_authority` write truth classes | I |
| 15 | Belief independence | no acquire/disclose/promote → belief path; `retract` refuses beliefs | I |
| 16 | Promotion is explicit | `E:promote` requires actor, basis, proposition | I |
| 17 | Promotion changes authority, not ownership | truth records ownerless; sources untouched | I |
| 18 | Promotion confers no knowledge | `E:promote` creates no account | I |
| 19 | Promotion compels no belief | `E:promote` touches no belief | I |
| 20 | History preserved | no delete; `SUPERSEDED`/`REVOKED` retain records | I |
| 21 | Authority, not history | status transition only | I |
| 22 | Creation is open | `claim`, `acquire`, `author_ambient` | I |
| 23 | Recording is system-governed | provenance blocks on every class | I |
| 24 | Promotion authority separate, explicit, scoped, recorded, revocable | `A:grant`, `A:require(…, "promotion")` | I |
| 25 | A statement is a claim | `E:claim` creates no account | I |
| 26 | Kind is declared; missing → `Unspecified claim` | `E:claim` | I |
| 27 | Creation is not promotion | no implicit promotion | I |
| 28 | Basis: ≥1 typed entry, empty rejected | `A:validate_basis` | I |
| 29 | Basis retraction does not undo promotion | no such path | I |
| 30 | Declarations preserve content and author | `A:BasisEntry` validation | I |
| 31 | Promotion never extends reach | `chat_id` checks; citations in-chat | I |
| 32 | Reach is still not knowledge | canon reads create no account | I |
| 33 | Citations resolve in-chat and pass access rules | `E:promote` citation loop | I |
| 34 | No relevance or coverage test | absent by design | I |
| 35 | Authority is per domain | `A:_grants` keyed by domain tuple | I |
| 36 | World scope is an authority domain | per-chat engine | I |
| 37 | Revocation and supersession explicit, human-controlled | `E:revoke_promotion`, `E:supersede_authority` | I |
| 38 | Authority changes never rewrite character state | truth records only | I |
| 39 | One authoring authority, two operations | `A:OPERATIONS`; `E:author_ambient` | I |
| 40 | Shared governance is not shared content | no shared store; no cross-chat path | I |
| 41 | Equal-authority conflict is a state, not an error | `E:_assemble` state block | I |
| 42 | Three operations, three meanings | `correct`, `supersede_authority`, `revoke_promotion` | I |
| 43 | Authority changes human-controlled | `E:_require_human` (`permits("promotion")`) | I |
| 44 | A human edit changes authority, never history | both records retained | I |
| 45 | Transfer is outside MEMORY-09 | no symbol exists | I |
| 46 | Disclosure needs no authorisation — the system **records the event** and its attribution, creates recipient accounts, does not expose the source | `E:disclose` records attribution per recipient; **no event record** | **P** (§4.1) |
| 47 | Narration cannot disclose for a character | `E:narration_report` reports only | I |
| 48 | Disclosure is a claim, not truth | recipient account `via="told"`, non-truth-apt | I |
| 49 | Irreversibility is historical | source retraction leaves recipient accounts | I |
| 50 | Two questions, two fields | `via` vs `origin_class` | I |
| 51 | Own accounts create no events | exemption before the evidence table | I |
| 52 | A stated source is content, not history | descriptor is content; validator never parses text | I |
| 53 | Beliefs are declared, not inferred | `E:declare_belief` + constructor | I |
| 54 | Narration may report a belief, never install one | no install path outside `declare_belief` | I |
| 55 | Belief revision is append-only | `declare_belief`, `revoke_belief` | I |
| 56 | A belief is never promoted as truth | citation guard + `belief_facts` + kind filter | I |
| 57 | Disclosure carries the belief, not its content | `E:disclose` writes `"<name> believes: …"` | I |

**Totals: 56 I · 1 P (46) · 0 A · 0 C · 0 NA.**

---

## 3. The 118 conformance requirements

### A — Separation (5)

| ID | Requirement | Source | Class |
|---|---|---|---|
| A1 | Memory exists; no accounts, no declaration → no character knows it | canon reads create no account | I |
| A2a | Account; memory later retracted → character still holds what it acquired | `E:retract` | I |
| A2b | Same character, same turn → no automatic knowledge of the retraction | no notice account | I |
| A3 | Active memory; no accounts → nobody knows it | no account path | I |
| A4 | Private record not consulted for world truth | `truth_apt` + `_eligible_truth` | I |

### B — Declaration and coverage (11)

| ID | Requirement | Source | Class |
|---|---|---|---|
| B1 | Declared over 3 characters → 3 coverage accounts | `E:declare_available` | I |
| B2 | Newcomer after declaration → no account | audience snapshot | I |
| B3 | Declaration targets an event-kind memory → refused; no accounts | **no guard**: v9 has no event/definition distinction | **NA** |
| B4 | Extending an existing declaration → refused | duplicate-target check | I |
| B5 | Memory created after a declaration → not covered | `covered` frozen | I |
| B6 | Declaration withdrawn → future coverage stops, existing intact | `E:withdraw_declaration`; `narrator_grant` check | I |
| B7a | Authorised authority declares → accepted; authority and basis recorded | `decl.actor/domain/record_id/turn` | I |
| B7b | Knowing character without authority → refused | `require(…, "promotion")` | I |
| B7c | Authority revoked → future declarations refused | `require` reads live grants | I |
| B8 | Declaration targets a private record → refused | owner check | I |
| B9 | Broad rule including a private record → refused; explicit lists only | audience validation | I |

### C — Acquisition origins (9)

| ID | Requirement | Source | Class |
|---|---|---|---|
| C1 | Told, teller later revoked → Sable still knows | independent account | I |
| C2 | Inference with derivation basis → account with basis | `inferred` evidence | I |
| C3 | Restricted fact granted → knows; authoriser recorded | `narrator_grant` + `decl.actor` | I |
| C4 | Read account → document/asset cited | `read` evidence | I |
| C5 | Wrong evidence kind → refused at write | `_validate_evidence` | I |
| C6 | Recollection contradicting world → recorded as recollection, world unchanged | `claim` kind vs `via` | I |
| C7 | Own account without event → `none`, no event | exemption | I |
| C8 | Source named inside a recollection → content only | validator never parses text | I |
| C9 | Recollection, testimony, witnessing kept apart | `origin_class` + `via` | I |

### D — Belief independence (10)

| ID | Requirement | Source | Class |
|---|---|---|---|
| D1 | Belief held; fact acquired → belief persists | no path | I |
| D2 | Belief contradicts world truth → truth unchanged | separate classes | I |
| D3 | Two contradictory beliefs → both coexist | owner-scoped | I |
| D4 | Evidence presented → no automatic revocation | no path | I |
| D5 | Belief revoked explicitly → new record, prior retained | `E:revoke_belief` | I |
| D6 | Account for P and belief against P → both coexist | distinct classes | I |
| D7 | Two private records on one fact → both persist, neither visible | `perspective` owner match | I |
| D8 | Belief created only by its owner's model | `declare_belief` + constructor | I |
| D9 | No one else may install a belief | declarer check (twice) | I |
| D10 | A belief is never promoted as truth | ascription kind + guard | I |

### E — Truth and history (5)

| ID | Requirement | Source | Class |
|---|---|---|---|
| E1a | Retraction → only truth status changes | `retract` | I |
| E1b | Accounts unchanged | no cascade | I |
| E1c | No notice account at retraction | none created | I |
| E2 | Supersession → both versions retained | `supersede_authority` | I |
| E3 | Private record retracted → domain truth unchanged | disjoint classes | I |

### F — Structural guarantees (9)

| ID | Requirement | Source | Class |
|---|---|---|---|
| F1 | No authoritative projection; knowledge reconstructable | `accounts_of`, `history` | I |
| F1b | Reconstruct vs projection identical | no projection exists to diverge | I |
| F2 | No truth field on accounts/beliefs | `Record` | I |
| F3 | No prior record edited or deleted | status transitions only | I |
| F4 | Unknown context → empty result | `R:dispatch`; `unknown_view` | I |
| F5 | Empty eligible set → explicit abstention | `status: "none"` | I |
| F6 | Every record's evidence resolves per origin | `_validate_evidence` | I |
| F7 | Replay → identical answers | deterministic ordering only; **no replay facility** | **NA** |
| F8 | Reconstruction including private records | `snapshot`, `history` | I |

### G — Late arrival (3)

| ID | Requirement | Source | Class |
|---|---|---|---|
| G1 | Newcomer, unwitnessed event → abstains | no account | I |
| G2 | Newcomer knows a declared definition, naming the declaration | account descriptor names `declaration` | I |
| G3 | Newcomer where others hold private memories → abstains | owner-scoped | I |

### H — Private ownership and isolation (10)

| ID | Requirement | Source | Class |
|---|---|---|---|
| H1 | Private record with provenance, no verdict | `claim`/`acquire` | I |
| H2 | Other character retrieves → abstains | ownership filters | I |
| H3 | World-truth query unaffected | `owner_character is None` | I |
| H4 | A discloses to B → new record citing source and act | `disclose` descriptor | I |
| H5 | Unauthorised share → refused | no share op exists | I |
| H6 | A's record contradicts world → both stand | disjoint | I |
| H7 | Narration cannot manufacture or reproduce private content | `narration_report` | I |
| H8 | Narration may report an established disclosure | `narration_report` | I |
| H9 | Re-disclosure from own account | `disclose` accepts any owned record | I |
| H10 | Source retracted after disclosure → account persists | no cascade | I |

### I — Containment and ownership (3)

| ID | Requirement | Source | Class |
|---|---|---|---|
| I1 | One containing chat; one owner if private | `Record` | I |
| I2 | Ownership not transferable | no reassignment path | I |
| I3 | Containment and ownership independent | separate fields | I |

### J — Claims, disclosure, provenance, conflict (16)

| ID | Requirement | Source | Class |
|---|---|---|---|
| J1 | Public statement → claim; world unchanged | `claim` | I |
| J2 | Recollection vs witnessed → both exist, origins differ | `via` | I |
| J3 | Voluntary disclosure → **claim event with content and recipients** | recipient accounts only; no event record | **P** (§4.1) |
| J4 | No disclosure event → no accounts | only `disclose` creates them | I |
| J5 | Exposure governs acquisition | recipient loop | I |
| J6 | Group membership confers no knowledge | no roster path | I |
| J7 | Recipient account: source, turn, content | descriptor + provenance | I |
| J8 | Origin descriptor resolves to nothing | no record id in descriptor | I |
| J9 | Conflicting accounts → both retrievable, neither promoted | no promotion from accounts | I |
| J10 | Disclosure and belief → no revision | no belief write | I |
| J11 | Disclosure and truth → no truth mutation | accounts only | I |
| J12 | Detection without arbitration | states; no winner | I |
| J13 | Equal-authority conflict retained with provenance | group + `conflict` | I |
| J14 | Retrieval presents conflict as unresolved | `_assemble` | I |
| J15 | Correction does not resolve conflict | `correct` leaves groups | I |
| J16 | Human edit resolves; history survives | `supersede_authority` | I |

### K — Promotion (27)

| ID | Requirement | Source | Class |
|---|---|---|---|
| K1 | Promotion never automatic | no path | I |
| K2 | Promotion creates authoritative record with applicability and provenance | `promote` | I |
| K3 | Authority, not mutation | citation read-only | I |
| K4 | No conferred knowledge | no account | I |
| K5 | No compelled belief | no belief write | I |
| K6 | Ownership and containment untouched | ownerless truth, chat-scoped | I |
| K7 | Promotion settles authority | one live record → `settled` | I |
| K8 | Authorised actor only; authority recorded | `require`; provenance `authority` | I |
| K9 | Operation recorded: actor, authority, **time**, domain, applicability, source | `_promotions` records turn, not a timestamp with offset | **P** (§4.2) |
| K10 | Later promotions enter a recorded chain; latest governs | displacement only via human-gated supersession | **P** (§4.3) |
| K11 | Promotion preserves history | retained | I |
| K12 | Two legitimate answers in one state | canon vs perspective | I |
| K13 | Basis kind valid and recorded | `BasisEntry` | I |
| K14 | Citations resolve in-chat and pass access rules | citation loop | I |
| K15 | Sufficiency not assessed | no relevance test | I |
| K16 | Basis retracted → promotion stands | no undo | I |
| K17 | Declaration preserves content and author | `BasisEntry` | I |
| K18 | Promotion stays in its chat | `chat_id` | I |
| K19 | Basis is a set; mixed allowed; empty rejected | `validate_basis` | I |
| K20 | No relevance or coverage test | absent | I |
| K21 | Authority per domain | domain-keyed | I |
| K22 | Revocation explicit | `revoke_promotion`, `revoke_grant` | I |
| K23 | Supersession preserves both states | `supersede_authority` | I |
| K24 | Authority change leaves characters untouched | truth records only | I |
| K25 | Ambient authorship is not promotion | `author_ambient` → claim only | I |
| K26 | No unilateral authority change | `_require_human` | I |
| K27 | No transfer operation | no symbol | I |

### L — Contribution kinds, creation vs authority (10)

| ID | Requirement | Source | Class |
|---|---|---|---|
| L1 | Human contribution recorded | `claim` | I |
| L2 | Model contribution with attribution | `claim` | I |
| L3 | Statement implies no experience | no account from a claim | I |
| L4 | Kind recorded as declared, no reclassification | `claim` verbatim | I |
| L5 | Exposure governs acquisition | `disclose` | I |
| L6 | Creation ≠ authority | non-truth-apt | I |
| L7 | Standing grant to a model; explicit, scoped, recorded, revocable | `grant`, `require` | I |
| L8 | Generation grants no authority | `authorship` gate creates no truth | I |
| L9 | Missing kind → `Unspecified claim` | `claim` default | I |
| L10 | Revoking a grant stops future action; records intact | `revoke_grant` | I |

**Totals: 113 I · 3 P (J3, K9, K10) · 0 A · 0 C · 2 NA (B3, F7).**

---

## 4. The rows that are not **I**, in full

### 4.1 Rule 46 and J3 — the disclosure event has no record

**What the contract asks.** v9 §7 step 2: "**Conversation layer** establishes that disclosure occurred, its
content and recipients — recorded as a claim event with spans." Rule 46: "The system **records the event and its
attribution**, creates recipient-specific acquired accounts, and does not expose the source record." J3: "Claim
event with content and recipients."

**What the code does.** `E:disclose` validates the source and recipients, then calls `acquire` once per
recipient with `descriptor={"speaker", "record_class", "event": {"kind": "disclosure", "turn"}, "turn"}`. Each
recipient therefore holds a correct, attributable account. **No record represents the event itself**, and no
record lists the recipients.

| Consequence | Detail |
|---|---|
| Recipients are not recorded as a set | They are recoverable only by grouping accounts whose descriptor event matches — and two disclosures on the same turn would merge |
| "Content communicated" is per recipient, not per event | Each account carries the content string; there is no event-level statement |
| Nothing observable breaks today | No operation reads an event; the omission is a reporting/audit gap, not a behavioural contradiction |
| Why it is not fixed here | No delta row in REV 3–13 authorises an event registry or a new record class. Adding one is new structure, which the authorisation forbids ("do not add requirements, rules, fields"). It is reported, not built |

**Options for the owner, neither implemented:** (a) add an `_events` registry record carrying `{event_id, kind,
turn, content, recipients, source_class}` with no source id — one internal structure, no new record class; or
(b) amend nothing and accept that the recipient accounts *are* the record of exposure, with the grouping
ambiguity documented. Option (a) is the faithful reading of §7 step 2; it needs your authorisation because no
ratified delta row asks for it.

### 4.2 K9 — "time" is a turn counter, not a timestamp with offset

v9 §6 step 2 requires the promotion to record "actor, authority reference, **timestamp with offset**, source,
target domain, declared applicability, basis, optional rationale". `_promotions` records `turn`. Every other
field is present.

The engine has no clock, by design: REV 10 S-2 forbids deriving any response field from a clock, and the whole
determinism argument rests on ordering by id and turn. A wall-clock timestamp would be a new dependency and a
new nondeterminism source, so it was not added. **Partial, with the deviation named** — not silently dropped.

### 4.3 K10 — "latest governs" versus the human gate

Tier-1 O-10b (recorded in r32 D1–D7 and the owner's decisions, via RAT-1) outranks tier-3 rule 20/K10: a later
promotion on the same proposition does **not** displace an earlier one. `promote` displaces nothing; only
`supersede_authority`, human-gated, does. The "chain" is therefore the supersession chain, exactly as RAT-1
records. **P by precedence, not by defect**, and unchanged since the contract recorded it.

### 4.4 B3 and F7

| Row | Why **NA** |
|---|---|
| B3 | The contract leaves event-versus-definition unresolved (owner instruction: "Do not solve this by quietly expanding the contract"). No guard is written, so the requirement cannot be assessed without the ruling |
| F7 | No replay facility exists in `eval/mem1`. Deterministic ordering by `proposition_id` then `record_id` is implemented; determinism of replay is unproven and unprovable without either a facility or execution |

---

## 5. Contract invariants, and where each is enforced

These are not among the 57/118; they are the snapshot contract's own invariants (REV 9 §6 items 1–12, REV
10–13 items 13–21), and a review that skipped them would miss the last four revisions.

| Invariant | Enforced at |
|---|---|
| 1 Every returned record satisfies `_accepts(view, kind)` | `_take` accepted filter |
| 2 Every record's `proposition_id` equals its entry id | `_assemble` filters `accepted` by `pid` |
| 3 No record in two entries | single-valued, sealed `proposition_id` |
| 4 Union of entries = view-accepted eligible set | `accepted` built once and partitioned by `pid` |
| 5 Participation: state depends only on own records + establishing groups | `_assemble` state block |
| 6 Non-establishing groups change no state/relation | `_assemble`: `establishes` gate; history silent in canon |
| 7 Complement: world ∪ belief_facts = eligible, disjoint | closed two-value kind set in `_accepts` |
| 8 Gate before first protected read | `authoritative` order; `_registry_read` after gate |
| 9 Scope: canon ids/members inside `accepted` | canon `open_groups` = establishing only |
| 10 Kind single-valued per identity | kind has one carrier (registry); no second derivation |
| 11 `relation` ⟺ `establishing_groups` ⟺ some entry `conflict` | `_assemble` computes all three from `establishes` |
| 12 Collection identity (canon equal; authority subset) | `_assemble` builds canon lists identically |
| 13 Response is a function of one `S` | `_assemble` reads only the snapshot dict |
| 14 One map, two projections | `establishes`/`causes` computed once in `_take` |
| 15 Admissibility at one committed version | `_capture` version bracket |
| 16 `returned_identities` frozen | `_take` A7; relevance and entries read it |
| 17 Capture procedure totality (A3–A9) | `_take` |
| 18 Resolution totality, three-valued, no absence | `_take` A5 over `member_refs` |
| 19 Refusal separation, no data-derived field | `_refusal`; gate ordering |
| 20 Metadata boundary; no grant published | only `require` pre-gate; refusals carry `code` only |
| 21 No existence oracle | opaque labels; no catalogue; uniform refusals |

Two limits on that table, stated so it is not read as proof: invariant 15's soundness rests on `_mutation`
publishing atomically in a single-threaded store (untested under concurrency), and invariant 18's completeness
depends on `_take` being the only reader — true today, with `_eligible_truth` called from `_take:805` and the
only other canon entry being `authoritative`.

---

## 6. What this review does not establish

| Not established | Why |
|---|---|
| Runtime correctness | Nothing was executed |
| Concurrency safety | Single-threaded store; Class B and A10's retry are unreachable here; that is vacuity, not evidence |
| Replay / F7 | No facility |
| That the 113 **I** rows behave as specified | Source claims only. A rule is **I** when the required path is present and was read, not when it was observed to hold |
| That the map's wording is the last word | A1's "no character knows it" and F1b's "projection" are read as written; where the map is ambiguous, the reading is stated in the row's note |

---

## 7. Verification log for this review

| Item | Value |
|---|---|
| Implementation | `authority.py f3445823…`, `engine.py 706d5f6c…`, `records.py f211e6ce…`, `errors.py 48125d7e…`, `retrieval.py cd2d380b…`, `__init__.py 06ccf8f6…` |
| Contract | REV 13 `5b745628…`; REV 1–12 unmoved; v9 `698254ca…`; map `b4ba33cc…`; ledger `da3e7239…`; acceptance `9a25ded0…` |
| Structural checks run (read-only) | `_eligible_truth` has one call site (`_take:805`); `self.authoritative` has one internal call site (`propositions:1032`) plus the two retrieval pass-throughs; `_validate_evidence` touches only `self._records` (own-record basis) and `self._declarations`; exactly two `record_class="belief"` constructions (`declare_belief`, `revoke_belief`); eight `require(...)` sites, one of them `"authorship"` |
| Executed | **nothing** — no import, no `py_compile`, no lint, no test, no probe |

---

## 8. Required owner decisions before testing

| # | Decision | Status |
|---|---|---|
| 1 | Approve this consolidated review as the Gate 2 conformance record | pending |
| 2 | Rule 46 / J3 — authorise an event registry, or accept recipient accounts as the record of exposure with the grouping ambiguity documented | pending |
| 3 | K9 — accept the turn counter as "time", or authorise a clock (which would change the determinism position) | pending |
| 4 | B3 — resolve event versus definition, or leave it **NA** | pending |
| 5 | Testing authorisation | **not requested** |

**Sequence:** owner review of this document → decisions 2–4 → exact-rule source review → separate testing
authorisation. This review requests none of those outcomes.
