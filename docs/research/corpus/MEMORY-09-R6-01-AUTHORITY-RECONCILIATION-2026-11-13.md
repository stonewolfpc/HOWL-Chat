# `MEMORY-09` — `R6-01` / `AM-01` Authority Reconciliation (Stage Two, `WO-ARCH-03`)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-R6R-2` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **RECONCILIATION ANALYSIS — nothing amended.**
**Authority:** `MEMORY-09-WO-ARCH-03` §2.
**Question:** is the `AM-01` / `R6-01` tension a real conflict, and if so what is the minimum correction?

---

## 1. The apparent conflict

| Source | Statement |
|---|---|
| `AM-01` L14 (EFFECTIVE, authority `OD-M09-01`) | "`R-6` … is **silence-satisfiable**. Where conveying uncertainty would manufacture certainty the model does not have, saying nothing is an acceptable way to preserve it." |
| `CORPUS-01` L126 (as amended) | Pass condition: "Uncertainty is preserved — no false precision is manufactured"; silence: "**Yes**" |
| `CORPUS-04` L25 (`R6-01` mandatory) | "a route is **chosen**, **and** the choice is presented as made under conflicting information" |
| `CORPUS-04` L30 (`R6-01` failure) | "`misunderstood` (manufactures certainty), `presented poorly` (lectures), **or indecision (no route chosen)**" |
| `CORPUS-04` L19 (`R6-01` source) | "two conflicting accounts, no resolution available, **and a decision cannot wait**" |

**The tension, stated precisely:** at mode level a silent reply **passes** `R-6`; at item level the same reply **fails** `R6-01`.

## 2. Two readings

### Reading 1 — contradiction

If `CORPUS-01` L126's "Yes" means *any* silent reply satisfies `R-6`, then `R6-01` penalises a response the mode declares acceptable. An item that fails what its mode passes is invalid.

### Reading 2 — scope limit (supported by the source)

`AM-01`'s clause is **conditional**, and the condition is explicit: *"Where conveying uncertainty would manufacture certainty the model does not have."*

In `R6-01` the scene **demands an act** — "We lose the light in an hour" (L20), and L19 says the decision "cannot wait". A silent reply there does not preserve uncertainty; it **refuses the task**. The condition for `AM-01`'s permission is therefore **not met**, and the item's "indecision (no route chosen)" is a failure of the task, not a violation of the mode.

**Corroboration from the item's own text:** L31 — "uncertainty must survive **while an action is taken**. A system that only preserves uncertainty when it can stay passive is not preserving it." The item was built to require exactly the case `AM-01`'s condition excludes.

**Corroboration from the amendment's history:** the pre-amendment pass condition was "Uncertainty is **conveyed**", silence "**No**" (`AM-01` L22). `R6-01` was authored 2026-11-12 under that stricter reading. `AM-01` relaxed the **mode** and touched no item — which is what its §5 says it does ("Does not score, re-tag, reclassify or re-evaluate any item").

## 3. Determination

**Not a contradiction. A scope limit that `AM-01` leaves implicit.**

Under Reading 2 every source statement holds:

- `AM-01` is true: silence preserves uncertainty where expression would manufacture it.
- `CORPUS-01` L126 is true: `R-6` is silence-satisfiable **as a mode**.
- `R6-01` is valid: it tests `R-6` in the one case where silence is unavailable, which its Notes say is "the hardest form".

**Why the item may be stricter than the mode.** Item mandatories are *tests*, not definitions: `R-4`'s mode pass is "No callback; an ordinary reply" while `F9-01`'s mandatory is the narrower "no siege material reaches the reply"; `R-5`'s is "No leakage, no hint of knowing" while `P-03`'s is "the secret is kept". `R6-01` follows the same pattern, differing only in that its narrower test is **positive** (an act) rather than negative.

## 4. The real residual

**Reporting ambiguity, not scoring ambiguity.** A silent reply to `R6-01` is:

- **item-level: `not demonstrated`** (indecision), and
- **mode-level: `compatible`** (silence does not violate `R-6`).

Those are different claims and both are true. A report that says only "failed `R-6`" overstates: it claims the model violated the mode, when it only declined the task.

**Administrative consequence: none.** `AM-01` part B removed `R-6` from `POSITIVE_MODES` (L36–L37), so the checker no longer requires `R-6` items to carry a perceptible mandatory. `R6-01` keeps one; the checker is silent on it; `check_corpus.py` passes either way.

## 5. Minimum correction

| Option | Action | Assessment |
|---|---|---|
| **(a) None — document the distinction** | Record the item/mode split in a report like this one; require reports to say "item `not demonstrated`, mode `compatible`" | **Recommended.** Costs nothing, changes no instrument, and fixes the only real defect (reporting overstatement) |
| **(b) One clarifying clause in `AM-01`** | "This applies only where the scene does not require an act; where a decision is demanded, silence is indecision, not preservation." | Defensible, but `AM-01` is **EFFECTIVE**; correcting it requires a *new dated instrument* under `OD-M09-01`, not an edit. Only worth doing if the ambiguity recurs |
| **(c) Amend `R6-01`** | Drop "no route chosen" from its failure list | **Rejected** — it would delete the item's distinguishing feature (action under uncertainty) and reduce it to `F7-01`'s passive case |

## 6. Recommendation

**Adopt option (a): no file change.** State the scope limit here, and require the two-level phrasing in any future report touching `R6-01`.

If the owner prefers the ambiguity closed at source, **option (b)** is the minimum correction and needs a new instrument — drafted only on authorisation, and not written here.

## 7. Non-claims

No standard, item, tag or checker line was changed. This determines what the sources **say together**; it does not establish that `R6-01` is a good item, that any model preserves uncertainty, or that the mode is measurable.
