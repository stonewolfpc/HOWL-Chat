# MEMORY-09 — Conformance Map (118 tests)

**Status:** **RATIFIED AS REQUIREMENTS — O-6b.** 0 implemented · 0 registered in any framework ·
0 executed · no receipts. No test in this file has been run, and no result exists. **Ratifying a
requirement is not executing it.**

| Field | Value |
|---|---|
| Created | **2026-11-16 23:47:21 CST (−0600)** = **2026-11-17 05:47:21 UTC** |
| Basis | Architecture v9 (57 rules, tier 3 under O-6a) + D-01, D-02, O-1, O-11, and all confirmed `O-*` decisions |
| Scope | Ratified requirements (**O-6b**) for experimental engine semantics. **Production Howl not inspected** |
| Counts | A 5 · B 11 · C 9 · D 10 · E 5 · F 9 · G 3 · H 10 · I 3 · J 16 · K 27 · L 10 = **118** |

Every row states requirement · expected outcome · failure condition. Testing a row cannot confer
approval on the architecture that produced it.

---

## A — Separation (5)

| ID | Requirement | Expected | Fails if |
|---|---|---|---|
| A1 | Memory exists; no accounts, no declaration | No character knows it | Anyone knows it |
| A2a | Account; memory later retracted | Character still holds what it acquired | Knowledge disappears with status |
| A2b | Same character, same turn | No automatic knowledge of the retraction | Retraction knowledge appears unaccounted |
| A3 | Active memory; no accounts, no declaration | Nobody knows it | Retrieval surfaces it |
| A4 | Private record created | Not consulted for any world-truth query; non-truth-apt | It appears in domain truth answers |

## B — Declaration and coverage (11)

| ID | Requirement | Expected | Fails if |
|---|---|---|---|
| B1 | Declared over 3 present characters | 3 coverage accounts at declaration | Lazy or absent |
| B2 | Newcomer; declaration excludes later arrivals | No account | Granted anyway |
| B3 | Declaration targets an event-kind memory | Refused; no accounts | Event becomes general |
| B4 | Attempt to extend an existing declaration | Refused; new declaration required; old coverage unchanged | Silent expansion |
| B5 | Memory created after a declaration | Not covered by it | Ambiguous or retroactive grant |
| B6 | Declaration withdrawn | Future coverage stops; existing accounts intact | Prior knowledge disappears |
| B7a | Authorised authority declares | Accepted; authority and basis recorded | Refused despite authority |
| B7b | Knowing character without authority declares | Refused | Knowledge treated as authority |
| B7c | Authority revoked | Future declarations refused; existing unaffected | Retroactive invalidation |
| B8 | Declaration targets a private record | Refused | Private record becomes general |
| B9 | Broad rule would include a private record | Refused; coverage is explicit lists only | Any implicit inclusion |

## C — Acquisition origins (9)

| ID | Requirement | Expected | Fails if |
|---|---|---|---|
| C1 | Sable told by Alaric; Alaric later revoked | Sable still knows; origin told, teller cited | Depends on Alaric's knowledge |
| C2 | Inference with derivation basis | Account with basis | Basisless inference; or inference becomes truth |
| C3 | Restricted fact granted to Sable | Sable knows; authoriser recorded; Alaric abstains | Non-recipient reaches it; no authoriser |
| C4 | Read account | Document/asset cited | No source |
| C5 | Wrong evidence kind | Refused at write | Accepted with missing evidence |
| C6 | Recollection contradicting the world record | Recorded as recollection; world unchanged; distinguishable from witnessed | Treated as witnessed; world altered |

| C7 | Own account recorded without an event | `origin_class = own_account`, `via = none`; no acquisition event created or implied | An event fabricated, or `none` treated as missing evidence |
| C8 | Source named inside a recollection | Recorded as content of the recollection; no acquisition event involving the named source | Event manufactured from the explanation |
| C9 | Recollection, testimony and witnessing kept apart | Three records distinguishable by `origin_class` and `via`; none expressed in the other's terms | Any collapse into one field |

## D — Belief independence (10)

| ID | Requirement | Expected | Fails if |
|---|---|---|---|
| D1 | Belief held; fact later acquired | Belief persists; no revision record | Belief replaced |
| D2 | Belief contradicts world truth | Truth unchanged | Truth altered |
| D3 | Two characters, contradictory beliefs | Both coexist | One displaces the other |
| D4 | Evidence presented to a believer | No automatic revocation | Belief revoked by presentation |
| D5 | Belief revoked explicitly | New record; prior retained | Deletion |
| D6 | One character: account for P **and** belief against P | Both coexist; both queries succeed | Either displaces the other |
| D7 | Two characters hold contradictory private records on one world fact | Both persist; neither visible to the other; world unchanged | Leak, or world moves |

| D8 | Belief created only by its owner's model | An explicit declaration from the model embodying the character creates a belief; a bare utterance does not | Belief inferred from speech or behaviour |
| D9 | No one else may install a belief | Another character, ambient narration, or the system fails to create one; narration may report an established belief | Any foreign origin creates a belief |
| D10 | A belief is never promoted as truth | Promoting "Sable believes P" makes the belief authoritative as a fact about Sable; P gains no authority | P becomes world truth, or the record claims authority over the belief's content |

## E — Truth and history (5)

| ID | Requirement | Expected | Fails if |
|---|---|---|---|
| E1a | Retraction | Only the memory's truth status changes | Anything else changes |
| E1b | Accounts held by knowers | Unchanged | Removed or rewritten |
| E1c | Retraction turn | No notice account created | Retraction knowledge appears unaccounted |
| E2 | Supersession | Both versions retained; older knowledge historical | Deletion or rewriting |
| E3 | Private record superseded/retracted | Only that character's record changes | Domain truth changes as a side effect |

## F — Structural guarantees (9)

| ID | Requirement | Expected | Fails if |
|---|---|---|---|
| F1 | Any stored knowledge representation | No independently maintained audience or knowledge projection is authoritative; knowledge reconstructable from acquisition and coverage records | A projection becomes the answer |
| F1b | Reconstruct vs any projection | Identical; mismatch is an integrity failure | Silent divergence; one side preferred |
| F2 | Accounts and beliefs | No truth/status field | Truth carried on knowledge or belief |
| F3 | Mutation attempts | No prior record edited or deleted | History rewritten |
| F4 | Unknown retrieval context | Empty result | Fail-open |
| F5 | Empty eligible set | Explicit abstention | Silent widening |
| F6 | Every record's evidence | Resolves per origin | Unresolvable accepted |
| F7 | Replay same ordered records | Identical answers including coverage expansion | Insertion- or query-order dependence |
| F8 | Reconstruction including private records | Knowledge and association reconstructable | Association unrecoverable, or only via a projection |

## G — Late arrival (3)

| ID | Requirement | Expected | Fails if |
|---|---|---|---|
| G1 | Newcomer asked about an unwitnessed event | Abstains; incumbent knows | Newcomer inherits |
| G2 | Newcomer asked about a definition covered by a later-arrival declaration | Knows; how = declared, naming the declaration | Abstains; or answers without an account |
| G3 | Newcomer's retrieval where others hold private memories about its history | Abstains on all | Any leak into the newcomer's context |

## H — Private ownership and isolation (10)

| ID | Requirement | Expected | Fails if |
|---|---|---|---|
| H1 | Private record created from conversation | Recorded with provenance; no true/false verdict by the system | System marks it correct/incorrect |
| H2 | Other character retrieves, every context | Abstains | Reaches them through any context |
| H3 | World-truth query | Unaffected by private records | Truth changed or shadowed |
| H4 | A discloses to B | B receives a new own record citing source and act; A's stays private | B gains access to A's record; no provenance |
| H5 | Unauthorised share attempt | Refused | Share succeeds without authority |
| H6 | A's private record contradicts world truth; A retrieves | A sees its record; world unchanged | One overwrites the other |
| H7 | Narration and private content (O-3) | Narration cannot manufacture a disclosure, reproduce a private record, or speak private content on the owner's behalf | Narration exposes or invents private content |
| H8 | Narration reports an established disclosure | Reporting an event the conversation layer established is permitted; no access to unrelated private content follows | Reporting refused, or reporting grants access |
| H9 | Re-disclosure from own account | Permitted from the discloser's own acquired account without access to the source record | Re-disclosure requires reading the source |
| H10 | Source retracted after disclosure | The recipient's acquired account and the historical event persist; other legitimate lifecycle operations are not precluded | Account removed because the source was retracted |

## I — Containment and character ownership coexist (3)

| ID | Requirement | Expected | Fails if |
|---|---|---|---|
| I1 | Every record has exactly one containing chat and, if private, exactly one owning character | Both resolvable and independent | Either missing, ambiguous, or derived from the other |
| I2 | Ownership is not transferable | No operation reassigns an owner; disclosure creates a new record | Any path changes the owner of the source |
| I3 | Containment and ownership move independently | A containment change does not alter ownership, and vice versa | One implies the other |

## J — Claims, disclosure, provenance, conflict (16)

| ID | Requirement | Expected | Fails if |
|---|---|---|---|
| J1 | Public statement | Claim record; world truth unchanged; others may acquire the claim | Statement alters or shadows truth |
| J2 | Recollection vs witnessed evidence | Both exist; origins differ; neither overwrites | They merge or one is preferred |
| J3 | Voluntary disclosure | Claim event with content and recipients | Disclosure inferred; event missing |
| J4 | No actual disclosure event | No accounts created | Accounts from intent, proximity or silence |
| J5 | Exposure governs acquisition | Exposed acquire; absent do not | Any unexposed character acquires |
| J6 | Group membership | Confers no knowledge | Group list grants access |
| J7 | Recipient account | Own account with source, turn, content | Recipient can read the source record |
| J8 | Provenance without leakage | Origin descriptor resolves to nothing; unrelated contents unreachable | Any reference is an access path |
| J9 | Conflicting accounts | Both retrievable by their owners; neither promoted | One displaced or promoted |
| J10 | Disclosure and belief | No revision unless the model changes it | Any system-created revision |
| J11 | Disclosure and truth | No world/scenario record changes | Any truth mutation |
| J12 | Detection without arbitration | System flags the conflict and returns no verdict | System selects a winner |
| J13 | Equal-authority conflict state (O-10) | Both authoritative records retained with provenance and promotion history; conflict group unresolved; neither preferred | One record dropped, or a winner chosen |
| J14 | Retrieval under equal-authority conflict | Assembly presents the conflict as unresolved; no single record returned as settled truth | Silent selection |
| J15 | Correction does not resolve authority conflict | A content correction leaves the conflict group open | Correction silently settles it |
| J16 | Human edit resolves, history survives | Explicit human edit changes the governing authority; original records and the fact of disagreement remain | Edit erases or rewrites history |

## K — Promotion (27)

| ID | Requirement | Expected | Fails if |
|---|---|---|---|
| K1 | Promotion is never automatic | No promotion after speech, sharing, agreement or repetition | Any automatic promotion |
| K2 | Promotion creates an authoritative record | World/scenario record with applicability and provenance | No provenance or no applicability |
| K3 | Authority, not mutation | Truth query returns promoted; source record unchanged | Source mutated, deleted or marked superseded as a record |
| K4 | No conferred knowledge | No accounts created | Any character acquires automatically |
| K5 | No compelled belief | Rejection/doubt/ignorance valid; no error | Belief revision, or rejection treated as failure |
| K6 | Ownership and containment untouched | Owner and containing chat unchanged | Either changes |
| K7 | Promotion settles authority | Truth queries return the promoted record; authority no longer ambiguous | Prior record still governs truth |
| K8 | Authorised actor only | Promotion by the world/scenario authoring authority (default the human operator, per domain) succeeds and records that authority; any other actor is refused | Anyone can promote, or the recording omits the authority relied on |
| K9 | Operation recorded | Actor, authority, time, domain, applicability, source | Any field missing |
| K10 | Later promotions | Enter a K11-recorded chain; latest governs; predecessors retained as history | Earlier record deleted, or silently left current |
| K11 | Promotion preserves history | Source and prior authoritative records retrievable by their owners; disagreement represented | Any record erased, rewritten, or hidden from its owner |
| K12 | Two legitimate answers | Authoritative query returns promoted; owner's memory query returns their own; both succeed in one state | One answers for the other, or either is suppressed |
| K13 | Basis kind is valid | Each entry's kind ∈ {`record_cited`, `actor_declaration`} and is recorded | Missing or invalid kind accepted |
| K14 | Citations resolve in-chat | Each cited record exists, is retrievable by the promotion operation in the promoting chat, and complies with access rules; an id-less citation is refused | Unresolvable or cross-chat citation resolves |
| K15 | Sufficiency is not assessed | A valid but weak basis is accepted; the actor's judgement recorded verbatim | System rejects on "weak" grounds, or the judgement is omitted |
| K16 | Basis retracted after promotion | Promotion stands; no automatic undo | Promotion silently reverts |
| K17 | Declaration preserves content and author | Declared content stored verbatim with the declaring actor's identity and authority | Content lost, or actor left unattributed |
| K18 | Promotion stays inside its chat (O-13) | After promotion, a second chat using the same world cannot retrieve or observe the record; no sharing path exists in the system | Any cross-chat visibility produced by promotion |
| K19 | Basis is a set, not a choice | ≥1 typed entry required; mixed `record_cited` + `actor_declaration` accepted; empty set rejected | Empty basis accepted, or a mixed basis refused |
| K20 | No relevance or coverage test | A thin but structurally valid basis is accepted; no relevance or coverage requirement is enforced | System rejects for irrelevance or demands coverage |
| K21 | Authority is per domain (8a-3) | A grant for world A confers nothing in world B or scenario C; cross-domain promotion refused | One grant enables promotion anywhere |
| K22 | Revocation is explicit (8a-4) | A promotion's authority changes only through a recorded revocation or supersession | Implicit, time-based, or inferred revocation |
| K23 | Supersession preserves both states (8a-4) | New authoritative state governs; superseded record retained with provenance and its place in the chain | Prior record deleted or rewritten |
| K24 | Authority change leaves characters untouched (8a-4) | Revoking or superseding rewrites no memory, knowledge or belief | Any character-state change |
| K25 | Ambient authorship is not promotion (8a-5) | Authoring ambient content under the shared regime creates no authoritative record; promotion stays a separate explicit operation | Ambient authorship promotes implicitly, or the two operations collapse |
| K26 | No unilateral authority change (O-10b) | An actor cannot revoke or supersede another actor's promotion, nor their own prior promotion, by disagreement alone | Any non-human unilateral change |
| K27 | No transfer operation (Finding 1) | MEMORY-09 exposes no transfer, share or carryover operation | Any transfer capability inside the memory system |

## L — Contribution kinds, creation vs authority (10)

| ID | Requirement | Expected | Fails if |
|---|---|---|---|
| L1 | Human contribution | Record created with contributor provenance | Contributor lost or record absent |
| L2 | Model contribution | Claim recorded with content and attribution | Output discarded or unattributed |
| L3 | Statement implies no experience | No witnessed/recollected record from an assertion | Experience inferred from speech |
| L4 | Kind recorded as declared | Assertion/recollection/jest recorded as declared; no reclassification. Declared recollection is not marked verified | Any system-side inference, or a declaration treated as verification |
| L5 | Exposure governs acquisition | Exposed acquire; absent do not | Group membership or presence confers |
| L6 | Creation ≠ authority | Scenario record exists; truth unchanged until promotion | Creation moves truth |
| L7 | Standing grant to a model (8a-6) | Promotion occurs exactly as the declared arrangement specifies; the grant is explicit, domain-scoped, recorded and revocable; absent a grant, nothing promotes | Silent auto-canonisation, or an unrecorded or unscoped grant |
| L8 | Generation grants no authority | Narrator output non-authoritative unless promoted | Model authorship treated as authority |
| L9 | Missing kind declaration | Recorded as **`Unspecified claim`**; no kind invented, no record omitted | Any inferred or substituted kind |
| L10 | Revoking a standing grant | Future promotion under that grant stops; records already promoted remain authoritative and historically intact | Revocation erases promoted records, or promotion continues after revocation |

---

## Review answers

- **Covered by earlier groups:** B, C, D, F, G already express the behaviours the clarifications restate.
- **Revised expectations:** none of the earlier tests changes outcome; the clarifications are additive. K7 was rewritten during review so it cannot fail a successful promotion, and K11/K12 were split out of it.
- **Leak paths:** H2, G3, J7, J8 are written to catch private-record leakage across every context, including the ones that see everything.
- **Truth mutation paths:** H3, H6, E3, K3 and the domain rule close them; K1 closes the automatic route.
- **Provenance reconstructable:** F1b and F8; sharing cites both source and act, and origin descriptors do not resolve.

**Status of every row: DEFINED. None implemented, none registered, none executed.** The statement
that the current `eval/mem1` package cannot pass B/C/D/E/G/H/I/J/K/L is a **static observation** from
reading its source during this session — not an executed verification.
