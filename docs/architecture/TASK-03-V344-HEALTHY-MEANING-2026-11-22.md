# TASK-03 — `V3-44`: what `Healthy` means at generation 1

| Field | Value |
|---|---|
| Purpose | Define exactly what a `Healthy` verdict asserts at generation 1, after the retargeting, creation-collision and identity findings, without adding a state |
| Status | **Read-only investigation. Outcome: Clarify. V3.1 unedited. Nothing adopted, nothing authorized.** |
| Method | Clause-by-clause reading of `V3-44`, rows 10/11/14–17, `V3-03`, `V3-06`, `V3-12`, `V3-50`; reconciliation of rows 572, 578 and 580 against what those rows actually observe |
| Recovery token | `TASK-03-V344-HEALTHY-MEANING-2026-11-22` |

Fifth of row 576's six baseline questions, after OQ-1 (row 577), OQ-2 (row 578), OQ-3b-b (row 579) and row 550 (row 580). It resolves the `V3-44` wording item (`OQR-10`) and does not reopen the rows themselves.

---

## 1. Scope, and a prerequisite elevated

**HM-01 — Scope.** This record answers one question: what may a `Healthy` verdict claim at generation 1, given that (a) row 572 shows correct content can be written into the wrong conversation, (b) row 580 shows a valid conversation can have a compromised identity and overwritten companions, and (c) row 578 shows physical location and logical identity are not interchangeable. It does not design creation, locking, identity binding or companion assessment.

**HM-02 — Outcome: Clarify.** The rows are sound; `V3-44`'s *justification* invites an inference the verdict does not support. `Healthy` establishes observed content availability and structural consistency **within the assessed conversation scope**. It does not certify authorship, exclusive creation, identity uniqueness, correct write targeting, or companion-data integrity. No state, reason or row is added.

**HM-03 — A prerequisite elevated, carried from the owner's review of row 580.** The creation-side requirement is a **foundational product correctness requirement**, not a reconciliation refinement: creating a conversation must never silently adopt, merge with, overwrite or redirect writes into an existing one. Two guarantees stay separate and must not be smuggled into each other:

| Guarantee | What it prevents | What it does not require |
|---|---|---|
| **Collision safety** | adoption, merge, overwrite, misdirected writes | staging, temp directories, atomic publication |
| **Crash-safe creation** | an interrupted creation leaving an apparently complete but partially initialized conversation | unique identity or exclusive claim by itself |

Staging and atomic publication may eventually be justified by the second; they must not be presented as necessary for the first. Neither is a reconciliation state.

---

## 2. What `V3-44` says, and what it actually establishes

**HM-04 — The sentence, verbatim:** *"Existing generation 1 conversation → **Healthy** under the original evidence rules (rows 14–15). No preceding generation exists, so nothing can contaminate it. Unchanged by the missing gate."*

**HM-05 — "Contaminate" has one meaning in this contract.** `V3-42` defines contamination as a mixed-generation log: an interrupted clear leaves G's records in place, a G+1 append lands after them, and the two are byte-indistinguishable from a legitimate G+1 history. `V3-43` names the harm: retracted content displayed as current. The generation-1 sentence is therefore true **about generation mixture**: with no preceding generation, no G can survive into G+1.

**HM-06 — The false inference is the leap from "no generation can contaminate it" to "nothing can be wrong with it".** Rows 14 and 15 require usable content, a current marker, and a settled transition. They do not read who wrote the content, which conversation it was intended for, whether the directory was created exclusively, whether another directory claims the same logical ID, or whether any companion file survived. The sentence's scope is narrower than its grammar suggests.

**HM-07 — Rows 10, 11, 14–17 share that narrowness.** Every `Healthy` row decides from the two history representations and the marker. None reads `memories.json`, `characters.json`, `world.json`, `scenario.json`, `lorebooks.json`, the directory name, the creation history, or any writer identity. Whatever else `Healthy` means, it cannot mean more than its inputs allow.

---

## 3. The three findings that break the broad reading

**HM-08 — Row 572: content can be written to the wrong conversation, with no generation involved.** `performSave` resolves its target at execution time by reading `LastChatID` (`app-chat.go:314-316`), so an autosave captured while X was selected can execute after selection moved to Y, load Y, and write X's in-memory state into Y. `O7-157` already requires deferred mutations to use captured identity and refuse rather than retarget, but nothing supplies that identity. The resulting directory is structurally valid: Y's log, X's content in the snapshot, Y's marker.

**HM-09 — Row 580: a valid conversation can be an adopted one, with its companions destroyed.** A creation in a colliding second overwrites `meta.json`, `history.json` and every companion, including `memories.json` — the probe observed `[{"text":"precious"}]` replaced by `[]` — while the old `history.ndjson` survives and is served under the new name. The directory loads cleanly. No generation transition occurred, and no record is corrupt.

**HM-10 — Row 578: one physical directory can be reached under several paths, and one logical ID can name several physical directories.** Junction and spelling aliases share a physical identity; a manual copy preserves `meta.ID` at a new location; `GetChatList` returns `meta.ID` while `LoadChat` treats its argument as the directory name. Identity is a label, not a proof of uniqueness.

**HM-11 — Calculated: each finding can produce a directory that every `Healthy` row accepts.** The content parses; the representations are prefix-compatible or equal; the marker is current and settled. The wrongness is in properties the rows never read: authorship, targeting, creation exclusivity, ID uniqueness, companion preservation.

**HM-12 — And none of it is detectable from the assessed bytes.** A record carries no writer identity, no operation identity and no conversation binding; `ChatHistory` has id, sender, text, role, timestamp, `hasImage`, `characterId` (`types.go:72-83`), none of which is provenance. `meta.ID` is a string the same writer can copy or overwrite. There is no per-record generation (`V3-18`). The scanner cannot tell a legitimate conversation from an adopted or mis-targeted one by reading it.

---

## 4. The four properties, separated

**HM-13 — The verdict may establish one of these and must not imply the other three:**

| Property | Does `Healthy` establish it? | Why |
|---|---|---|
| **Content is structurally available** — present, valid and mutually consistent, or legitimately empty/closed under the marker; attributable where content exists | **yes**, within the assessed scope | rows 10/11/14–17 read exactly these |
| **Content belongs to the intended logical conversation** | **no** | retargeting (row 572) and copied `meta.ID` (row 578) leave valid bytes |
| **Creation claimed a previously unused identity** | **no** | adoption and collision leave a coherent directory (row 580) |
| **User-authored companion data was preserved** | **no** | companions are never read; row 580 overwrites them |

**HM-14 — "Attribution" in HM-13 is generation attribution, not authorship.** `V3-12` defines it as marker currency, an authorized generation, and a settled transition. It answers "which generation does this content belong to", never "who wrote it" or "which conversation was it meant for". The record uses the word only in that sense.

**HM-15 — The assessed conversation scope, defined.** A verdict is about: the candidate directory selected by `meta.json` presence; its `history.ndjson`; its `history.json`; and its expectation marker. `meta.json` is used for candidacy and as the logical-ID label, never as evidence of authorship or uniqueness. Companions are outside the scope and are not assessed at all.

**HM-16 — Companions are excluded because the contract never parses them, not because they do not matter.** `V3-04`'s four properties are about the history representations. Row 573 already requires a `Healthy` verdict to state that memories were not assessed; row 580 shows the stronger fact that they can be silently replaced by a creation collision. A verdict that implied companion integrity would be asserting something no row observes.

**HM-17 — Identity is excluded because a label is not a proof.** `meta.ID` can be copied, overwritten or shared; the directory can be renamed; two paths can name one directory and one logical ID can name two. Row 578's separation is the reason the verdict may attach to a logical ID but may not certify its uniqueness.

**HM-18 — Targeting is excluded because content carries no destination.** There is no operation identity, no per-record conversation binding and no writer field (`V3-18`, row 574). A record that arrived through a retargeted save is indistinguishable from one written by the correct process.

---

## 5. What the verdict may say

**HM-19 — The `Healthy` statement, restated:** *observed content is present, structurally valid and mutually consistent within the permitted prefix relation, and attributable to the generation the marker authorizes; or the conversation is legitimately empty or closed under its marker; no row indicates loss.* Each of those clauses names evidence the scanner actually reads, and the empty/closed branch is what rows 10–11 establish.

**HM-20 — The non-claims, stated as part of the verdict's meaning:** `Healthy` does **not** certify authorship; exclusive creation; identity uniqueness; correct write targeting; or companion-data integrity. A UI, report or API that presents `Healthy` as any of those has misread the contract, exactly as `V3-24` says a product misreads it by treating `Healthy` as a write permit.

**HM-21 — The limitation applies to every `Healthy` row, not only to generation 1.** Rows 10, 11, 14, 15, 16 and 17 all decide from the same narrow inputs. `V3-44`'s generation-1 sentence is merely the place where the broad reading is most tempting, because "nothing can contaminate it" sounds absolute. Stating the limitation once, as the meaning of the verdict, avoids six copies of it.

**HM-22 — Nothing in the decision procedure changes.** No row gains a condition; no state, reason or transition is added; `V3-47` still forbids `Unavailable → Healthy` and automatic repair. The change is to what the word `Healthy` may be taken to mean, not to when it is produced.

**HM-23 — This follows an existing precedent rather than inventing one.** `V3-50` already separates ownership from availability: row 20 reports `Orphaned` after the conversation would otherwise be `Healthy`, and orphaning "changes ownership, not availability". Creation integrity, targeting and companions are further properties in the same position: reportable elsewhere or not at all, never folded into the availability verdict.

**HM-24 — Display is unaffected and must not fill the gap.** `V3-06` permits a conversation to be shown while not `Healthy`; it equally forbids a healthy display from implying the non-claims. A conversation may be displayed, be `Healthy`, and still have been written into the wrong directory.

---

## 6. Why this is not a seventh state

**HM-25 — A state needs a detectable condition.** Every state in `V3-40` is decided from observed evidence. "Compromised identity", "mis-targeted content" and "overwritten companions" are not observable from the assessed bytes (HM-12). A state named after them would be a guess dressed as a classification.

**HM-26 — `Damaged` would be false.** The probe's adopted conversation loads cleanly: valid meta, valid log, readable message. `Damaged` asserts a malformed or missing representation; none is.

**HM-27 — `Pending reconciliation` would be false as a default.** Pending means an unresolved transition or unattributed content. These directories have a current marker and settled transition; there is nothing unresolved to reconcile. Using Pending would invent a transition the evidence does not show.

**HM-28 — Therefore the correct instrument is a limitation.** It constrains the meaning of an existing verdict without changing its production. This is the same choice `V3-50` made for ownership and row 573 made for memories.

---

## 7. The amendment

**HM-29 — `V3-44`'s first case is restated:** *"Existing generation 1 conversation → Healthy under rows 14–15. No preceding generation can produce a mixed-generation log, so generation contamination is impossible. This establishes availability only; it does not certify authorship, exclusive creation, identity uniqueness, correct write targeting or companion integrity (§…)."*

**HM-30 — A new §1 term, `assessed scope`:** the candidate directory's two history representations and its expectation marker; `meta.json` for candidacy and logical-ID attachment only; companions excluded.

**HM-31 — A new verdict-meaning clause:** the `Healthy` statement of HM-19 and the five non-claims of HM-20, stated once and applying to every `Healthy` row.

**HM-32 — `V3-73` closes the `V3-44` wording item, and row 576's `OQR-10` is answered.** No row changes; the item was a meaning question, not a classification question.

**HM-33 — What the amendment does not do.** It does not repair creation, bind writes to conversations, assess companions, deduplicate identities, add provenance, change the marker, add a state or reason, or authorize any implementation. It states the ceiling of what the observed bytes can prove.

**HM-34 — And it is not an argument for the current two-file design.** If persistence is replaced by a transactional store with explicit conversation identity, the verdict still means "observed content is available within the assessed scope" and still must not certify the properties the store has not checked. The limitation is about evidence, not about preserving the layout.

---

## 8. Outcome

**HM-35 — Outcome: Clarify.** `V3-44`'s generation-1 conclusion is correct; its justification is narrower than it reads, and the verdict's meaning needs the five non-claims stated once. The bounded amendment is three clauses (`V3-44` wording, `assessed scope`, verdict meaning) plus the `V3-73` closure. No state, row, reason or transition changes.

---

## 9. Tests this would require

- **HM-T1** A generation-1 conversation with valid content and a current committed marker is `Healthy`, and the verdict's meaning lists the five non-claims.
- **HM-T2** A directory whose snapshot was written from another conversation's in-memory state, with prefix-compatible representations, remains `Healthy`; the verdict does not claim targeting.
- **HM-T3** An adopted conversation (meta renamed, companions emptied, old log retained) is `Healthy`; companions are reported as not assessed, never as intact.
- **HM-T4** A manual copy at a second physical location with the same `meta.ID` yields two `Healthy`-eligible directories; the verdict does not claim uniqueness and the two are not merged.
- **HM-T5** Two non-prefix representations remain row 6 `Unavailable`; the non-claims are never used to excuse an inconsistency.
- **HM-T6** `Orphaned` beside `Healthy` still reports ownership separately and leaves availability unchanged.
- **HM-T7** No seventh state or new reason exists, and no `V3-39` row gained a condition.
- **HM-T8** The negative assertion: no verdict reads companions, `meta.ID` uniqueness, writer identity, conversation binding or operation identity.

---

## 10. Nothing here is authorized

**No implementation follows.** This record authorizes no creation repair, identity binding, companion assessment, provenance field, lock, gate, dedup pass, migration, state or reason, and no change to `V3-39`, `V3-40`, `V3-44`'s row or any verdict producer. It defines what the word `Healthy` may mean and leaves every mechanism unimplemented.
