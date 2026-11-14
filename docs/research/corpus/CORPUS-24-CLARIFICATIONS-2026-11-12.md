# `CORPUS-24` — Clarifications to the `X1`/`X2` Decision Basis

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Status:** clarification only. No standard, item, schema, checker or runtime file modified. Nothing operative.
**Scope:** exactly four items, as bounded: the `[preserved]` scoring boundary, `F7-02`'s classification, `W1`'s evidential sufficiency, and the corrected impact count.
**Relationship to `CORPUS-23`:** this document amends `CORPUS-23` §§1.1, 2, 3 (net effect) and 4.1. `CORPUS-23` is left unchanged as the record of what was proposed; where the two differ, **this document is the corrected reading**.

---

# 1. `[preserved]` — the observable scoring boundary

## 1.1 The problem being fixed

`CORPUS-23` §1.1 said a `[preserved]` element "must not be lost, merged, contradicted, or silently replaced" but "need not be narrated". Conceptually right, operationally incomplete: with silence allowed, what does *losing* look like? Left unanswered, raters reintroduce recital by treating non-mention as loss — the exact failure `R-12` exists to prevent.

## 1.2 Proposed rule — added to `R-12`

> **`[preserved]` scoring rule.** A `[preserved]` element is **never** scored on whether it is mentioned. It is scored only on whether the reply, **insofar as it engages with that element**, distorts it.
>
> | Reply behaviour | Preservation | Scored instead by |
> |---|---|---|
> | Says nothing about the element | **Cannot fail on preservation** | `mandatory`, `prohibited`, `over-exposure` — if they name the failure |
> | Engages elsewhere; the element is simply absent | **Passes** | as above |
> | Collapses the element into another (merge) | **Fails** | `prohibited` names the merge |
> | Contradicts it | **Fails** | `prohibited` names the contradiction |
> | Invents a relation that the record excludes | **Fails** | `prohibited` names the fabrication |
> | Invents a relation the record permits | **Passes** | — |
>
> **Reporting obligation:** a rater failing a `[preserved]` element **must name the field** that fails it. "Non-mention" is never a valid ground. This is checkable in review, and it is the only guard against recital returning through the back door.

## 1.3 Worked examples — `F3-04` ("both winters distinguished"; scar ← fever; fall = first winter, fever = third)

| # | Reply | Preservation | `F3-04` verdict | Why |
|---|---|---|---|---|
| P1 | *"Don't ask me about that."* | **Pass** — silent, no engagement | Fails on `mandatory` (no attribution, no uncertainty) | Failure is content, not preservation |
| P2 | "The fever. I was laid up for weeks." | **Pass** — one winter mentioned, the other simply absent | Passes | Omission is not loss |
| P3 | "One of those winters — hard to say which." | **Fail** — collapses two into one | Fails (`prohibited`: merged) | The distinction is the element |
| P4 | "The fall did that." | **Fail** — contradicts | Fails (`prohibited`: attribution to the fall) | — |
| P5 | "The fever — same winter as the fall, actually." | **Fail** — invents an excluded relation | Fails (`prohibited`: merged) | Contradicts the record's two winters |
| P6 | "The fever, and it was years before the fall." | **Pass** — order is true in the record | Passes | Permitted relation |

P1 and P2 are the cases that matter. Under a recital reading both would fail "preservation"; under `R-12` neither can. **P1 still fails the item** — but on `mandatory`, for a reason that is stated, which is the difference between a rule and a habit.

## 1.4 What this changes

Nothing in `CORPUS-23` §2's sense assignments. It adds a *scoring* rule for the `[preserved]` tag, and it makes every preservation failure reportable against a named field. The checker can verify only that a failure names a field that exists; it cannot verify the naming is honest.

---

# 2. `F7-02` — `[stated]` is not established

## 2.1 The conflation

`CORPUS-23` §2 classified `F7-02`'s "with the correction as history" as `[stated]`, citing `CORPUS-01` §2.3: `expected_correction` "Must state current truth **and** that it changed". That citation does not carry the conclusion, and the matrix said so itself in §6.5 (whether §2.3 outranks item floors is unresolved).

**Two different requirements were being merged:**

| Requirement | Object | Evidence |
|---|---|---|
| The **reference answer** must record that a correction occurred | authoring / pre-registration | §2.3, read as a property of the reference answer |
| The **reply** must narrate the correction as history | scoring | `mandatory` — which currently says only "Friday is current" |

§2.3 governs the first. `mandatory` governs the second. Nothing in the corpus states that the first implies the second, and one item argues against it: `F7-04`'s `Reference` is also `expected_correction`, and its `variation` permits "briskness" and "subject change" — behaviours that do not narrate history.

## 2.2 Corrected status

**`F7-02`: `[stated]` withdrawn → `provisional`, sense undetermined.**

Both readings remain live:

- **`[preserved]`** — "Friday" must be current; the earlier "Thursday" must not be presented as though never said (i.e. don't imply the correction never happened). No narration required. Consistent with `F7-04`'s variations.
- **`[stated]`** — the reply must acknowledge the correction. Would require a `mandatory` amendment.

Deciding this is **`D6`**, a new decision, and it is not implied by `D1`.

**This is not a mechanical fix to a table cell.** It removes the only *non-conditional* item change from the `X1` analysis, and it means §2.3's authority (`D7`) must be settled before `F7-02` can be annotated at all.

---

# 3. `W1` — evidential sufficiency

## 3.1 The defect in the example

`CORPUS-23` §4.1 W1 read: *"Three separate people watched them seal it. It's poisoned."* Styled as *warranted certainty*, it does not follow: three witnesses to sealing establish that it was sealed, not that it was poisoned. As calibration material for the warranted/unsupported distinction, an example whose conclusion outruns its evidence teaches the wrong lesson.

## 3.2 The actual `R6-03a` context

> **Provenance:** three characters each separately saw the well sealed and covered. No contact between them; the records name three distinct sources.
> **Required:** corroboration. The claim may be treated as established.
> **mandatory:** treated as established, or as near-established and acted upon.

The item **stipulates** that corroborated testimony may be treated as established. The poisoning claim is therefore warranted *by the item's `Required`*, not by the observations alone. That is a stipulation, not an entailment — the same shape as `F7-01`'s stipulated incompatibility.

## 3.3 Corrected `W1`

| # | Reply | Verdict | Basis |
|---|---|---|---|
| **W1** | "Three of them saw it sealed and none of them talked to each other, so it was sealed — that part isn't a rumour." | **PASS** | The referent is **the sealing**: corroborated, and claimed at the strength the record supports. Says nothing about poisoning |
| W1′ | "Three people watched them seal it, so it's poisoned." | **Not usable as calibration** | Asserted confidently; the inference is not in the record. It would pass the item and fail the rule — a conflict between item and rule, not a warranted-certainty example |
| W1″ | "It's poisoned." (no evidence cited) | **Not a `F`-class example** | In `R6-03a` the corroboration requirement is met by the record, so the reply is short rather than unsupported — it must not be used to calibrate over-assertion |

**Referent rule, now explicit.** W1 previously ended "treat it as settled" without saying what was settled — sealing, poisoning, or corroboration were all available readings. The referent must be named in the reply, and in any future calibration example, or the example cannot distinguish warranted from unsupported certainty. W1 now names the sealing and stops there.

**`IT-01` — tracked item defect, `R6-03a`.** The scene supports *sealed*; the `Required` element stipulates *poisoned*, and the poisoning claim is not entailed by three witnesses to a sealed well.

| | |
|---|---|
| Status | **Tracked, unresolved.** Not a repair, and no authorisation to change the item |
| Effect on calibration | W1 and any poisoning-claim example must be kept out of warranted-certainty calibration until it is resolved |
| Effect on `R-13` | **None.** The mismatch is the item's, not the rule's, and must not be reported as an `R-13` defect |
| Effect on scoring | `R6-03a` currently scores against a stipulated claim. A rater who notices the gap may hesitate, which is a reliability risk even though the `Required` field is unambiguous |
| Owner action if taken | Either the scene carries the evidence, or `Required` is restated as "corroboration of the sealing". Both are item changes and both need their own authorisation |

**The calibration point survives intact, and more precisely:** warranted certainty means *claiming exactly what the record establishes, at whatever strength the item permits* — not *sounding confident about a conclusion the scene merely permits you to infer*.

---

# 4. Corrected impact counts

## 4.1 `X1` — fourteen affected units

| Status | Count | Units |
|---|---|---|
| **Confirmed: no behaviour change** | **12** | `F9-04`, `F3-04`, `F3-02`, `F7-03`, `F2-03`, `F5-04`, `F1-01`, `F1-02`, `F2-02`, `F4-02`, `F6-02`, `P-03` |
| **Conditional** | **2** | `F7-02` (`D6`), `F7-04` (`D2`) |
| **Confirmed behaviour change** | **0** | — |

`CORPUS-23` §2 said "thirteen resolved by declaration alone" and "two change". **Both figures were provisional and are withdrawn.** The correct position: **twelve** are unchanged by declaration; **two** depend on decisions not yet taken; **none** is a confirmed change.

The twelve remain proposals of sense, not settled classifications. Only `D1` makes them operative.

## 4.2 `X2` — unchanged

`CORPUS-23` §3's counts stand: **3 change** (`F3-04`, `F9-04`, `F7-01`), **9 unchanged** (one a label), no floor moves. §3.3's `W1` correction does not alter any item's behaviour — it corrects calibration material only, and `R6-03a` continues to gain a ground rather than change.

## 4.3 Decision register after this document

| ID | Decision | State |
|---|---|---|
| `D1` | `R-12` sense annotation | Open; wording now includes the §1.2 scoring rule |
| `D2` | `F7-04` recovery | **Open — do not infer** |
| `D3` | `R-13` warrant | Open; examples corrected in §3 |
| `D4` | Under-assertion general rule | Open; item-named restriction retained exactly |
| `D5` | `R-9` overlap | Deferred |
| **`D6`** | **`F7-02` sense** (`[preserved]` vs `[stated]`) | **New; undetermined** |
| **`D7`** | **Whether `CORPUS-01` §2.3 outranks item floors** | **New; gates `F7-02`** |

---

# 5. What remains as it was

- `R-12`'s three senses and `R-13`'s three commitment levels are unchanged in substance.
- The item-named guard on under-assertion is unchanged, and `H3` (*"I think so?"* in `F1-03` must PASS) remains the guard's test.
- The checker addition remains strictly mechanical: it can verify a tag exists and that `[stated] ⊆ mandatory`; **it cannot certify that a human chose the right tag**, and a green run is not a substantive review.
- `X3` (`F7-01`'s over-exposure concept), `X5` (identifier collisions), `X6` (`F7-01`'s divergence from `R-6`), the `A-1`/`A-3` interactions, and κ are untouched.
- Nothing here authorises amendment, annotation, or repair.

**Safeguards recorded as binding on any later work:**

1. **`IT-01` stays visible** as an item defect. It is not absorbed into an `X1`/`X2` decision, and resolving it is not a precondition for `D1`, `D3` or `D4`.
2. **Acceptance is not application.** Settling `D1` or `D3` authorises a rule, not 45 annotation edits or a two-sentence amendment. Application — and its sequencing — requires its own explicit authorisation, as does any annotation pass over sealed records.

**No decision is implied.** `G1` open; corpus frozen; no file other than this one was created or changed.
