# TASK-03 — OQ-3a: writer-lifecycle inventory

| Field | Value |
|---|---|
| Purpose | Inventory every path that creates or commits conversation content; answer whether each establishes the expectation |
| Status | **Contract addendum. Read-only investigation. No code changed.** |
| Amends | row 548 RC2-05 (superseded), RC2-03/RC2-06 (confirmed); row 547 unedited |
| Artifacts | `build/tmp/spike9/probe9.go`, `build/tmp/spike9-oq3a.log` |
| Recovery token | `TASK-03-RECONCILIATION-OQ3A-2026-11-22` |

Labels as before. Every row below states what was **Observed** on a disposable scratch root or **Source-derived** from named lines.

---

## 1. Inventory

| # | Path | Writer | Representation | Commit boundary | Expectation set? | Crash window |
|---|---|---|---|---|---|---|
| P1 | `NewChat` (creation) | `SaveChat` | `meta.json`, `history.json` = `[]`, characters/lorebooks/memories | directory + `meta.json` durable | **No** (correct: nothing expected) | none |
| P2 | `AppendChatMessage` (user/assistant turn) | `AppendHistoryEntry` | `history.ndjson` append | `Write` returns from the append | **No** | append → marker |
| P3 | edit / delete message (`ReplaceChatHistory`) | `RewriteHistory` | ndjson **and** `history.json` | temp file renamed over ndjson | **No** | same |
| P4 | clear history (`ClearChatHistory`) | `RewriteHistory([])` | ndjson emptied, `history.json` = `[]` | rename | No | content removed by design |
| P5 | exit (`performFinalSave` → `performSave` → `RewriteHistory`) | `RewriteHistory` | both | rename | No | same |
| P6 | **`SaveChatAs` (duplicate)** | `NewChat` + `SaveChat` | **`history.json` only — no `history.ndjson`** | `SaveChat` returns | **No** | content → marker |
| P7 | legacy migration (inside `AppendHistoryEntry`) | `AppendHistoryEntry` | json → ndjson, then append | first successful append | No | partial migrate then crash |
| P8 | content wipe restore (`restoreContentRoot`) | `os.Rename` of the staged root | whole tree returns as it was | rename | marker travels inside the tree unchanged | none (rename is atomic) |
| P9 | chat import | **does not exist** | — | `ExportContentItem` refuses `library.Chats`; `contentimport` has no chats handler | N/A | N/A |
| P10 | `DeleteChat` | `os.RemoveAll` | removal | — | No | — |

**Source-derived:** P1 `persistence.go:131-168`; P2 `app-prompt.go:301-311`; P3 `app-prompt.go:389`; P4 `app-prompt.go:406`; P5 `app-chat.go:340-350`; P6 `persistence.go:436-467`; P7 `persistence.go:675-692`; P8 `app-wipe.go:128-136`; P10 `persistence.go:471-483`.

---

## 2. P6 is not what I assumed, and it breaks row 548's rule

**Observed.** Source with three appended messages, then `SaveChatAs` across a second boundary so the IDs differ:

```
source      history.json 2 b ([])   history.ndjson 375 b
duplicate   history.json 503 b      history.ndjson ABSENT
LoadChat(duplicate) → 3 messages
```

`SaveChatAs` copies `History` into the new `ChatData` and calls `SaveChat`, which writes `history.json` and **never** touches the append log. The copy is a **live conversation with real content and no `history.ndjson`** — created under the current contract.

**RC3-01 — Row 548's RC2-05 is superseded.** Scoping the accepted representation by contract version would classify every duplicated conversation as **Damaged**. Version is the wrong discriminator; the copy is current-contract and ndjson-less.

**RC3-02 — Acceptance is content-based.** Content is available if any representation holds **at least one record**:

| Representation | Counts as content |
|---|---|
| `history.ndjson` | ≥ 1 parsed record |
| `history.json` | ≥ 1 record (**not** the 2-byte `[]` form) |

RC2-06's rule survives intact and is what does the work: an empty file is existence, not content.

**RC3-03 — This also resolves the row-548 "stale `[]`" worry correctly.** If `history.ndjson` is deleted but `history.json` holds records, the content is genuinely still loadable (`loadHistoryFile` falls back to it), so Healthy is the right answer, not Damaged. Damage is *content unavailable*, not *this particular file missing*.

**RC3-04 — P6 must set the expectation**, exactly like P2. A duplicate commits content without an append.

---

## 3. P1 can collide, and the collision produces the state we fear

**Observed, twice, after two creations:** `generateChatID` has one-second resolution and `NewChat` uses `MkdirAll` without exclusivity, so two chats created in the same second share a directory name.

The first probe run demonstrated the consequence accidentally: `SaveChatAs` inside the same second **reused the source directory**, overwrote its `meta.json` (422 b → 425 b, the source renamed to "duplicate"), wrote the copied `history.json` (503 b) over the source's `[]`, and left the source's `history.ndjson` (375 b) in place. One directory, two conversations, content and identity disagreeing.

**RC3-05 — Creation is the one path where the expectation question is not the interesting one.** A collision yields *content present, metadata belonging to a newer empty conversation* — the exact shape the state table would read as Healthy-with-content. This is not a reconciliation bug, but reconciliation cannot repair it, and a scan must not "fix" it.

**RC3-06 — Unestablished:** whether any user-visible flow can realistically create two conversations within one second. The mechanism is confirmed; the frequency is not. Not claimed either way.

---

## 4. Every write resolves the target through one global pointer

`SaveChat` sets `pm.config.LastChatID = chatData.Meta.ID` (`persistence.go:239`), and P2–P5 all resolve their target with `GetLastChatID()` (`app-prompt.go:297,392,402`; `app-chat.go:331`).

**RC3-07 — The commit boundary is not conversation-scoped.** After `SaveChatAs` the pointer names the **copy**, so the next message, the next edit, the next final save and the next clear all target the duplicate. **Source-derived**, from the lines above; not observed at the app level.

For reconciliation this matters because content can be committed to a directory that is not the one the user believes they are writing to — a conversation can gain and another lose, with no damage in the filesystem sense.

---

## 5. The observation gap you raised, and what I will and will not claim

You are right that append-then-mark leaves a gap. Stated precisely:

1. append succeeds;
2. process dies before the marker is persisted;
3. content disappears before any scan observes it;
4. the manifest still says content was never expected → **Healthy**.

**RC3-08 — I do not claim the expectation marker is crash-safe.** It cannot be, as a second write after the first: there is always an interval in which one landed and the other did not. **Unestablished**, and recorded as such rather than asserted.

Ordering is still correct, because reversing it trades a silent-Hidden failure for a loud false one. What fixes the *representation* of the gap is a third marker state:

**RC3-09 — The marker has three values, and `pending` is what makes the gap visible.**

| Marker | Content observed | State |
|---|---|---|
| `none` | none | **Healthy** (empty conversation) |
| `none` | present | Healthy — and the scan **promotes** the marker to `committed` |
| `pending` | present | Healthy — promote to `committed` (the append did land) |
| `pending` | none | **Pending reconciliation** — an append was attempted and its outcome cannot be established |
| `committed` | present | Healthy |
| `committed` | none | **Damaged** |

**RC3-10 — Written before the append, `pending` can never produce a false Damaged**, because `pending` with no content is not Damaged. Written after a successful append but before the promotion, a crash leaves `pending` — and if the content then vanishes, the result is **Pending reconciliation**, not Healthy. The unknowable case becomes a visible question instead of a silent reassurance.

**RC3-11 — Promotion is monotonic.** Once `committed`, never reset by an observation. The scan may write the marker and nothing else: never content, never a removal, never a repair.

**RC3-12 — Residual window, stated as a limit.** From the last successful append to the next completed scan, one crash can leave `pending`+nothing, which surfaces as Pending reconciliation. That is the worst case, and it is visible. It is not zero, and the contract must not pretend otherwise.

---

## 6. Tests: what exists, and what does not

| Lifecycle | Covered? |
|---|---|
| append / rewrite round-trip | yes — `history_character_test.go:50,61,101` |
| `SaveChatAs` preserves content | yes — `TestSaveChatAsPreservesCharacterID:125` |
| `SaveChatAs` **representation** (no ndjson) | **no** |
| creation ID collision | **no** |
| expectation marker, any state | **no** (does not exist yet) |
| `LastChatID` redirection after `SaveChatAs` | **no** |

**RC3-13 — The existing `SaveChatAs` test asserts content and passes; it cannot see the representation gap** because `LoadChat` reads either form. A test that only checks "the messages came back" will never catch a missing append log.

---

## 7. Tests added for the contract

- **RT3-01** A duplicated conversation (content in `history.json`, no ndjson) scans **Healthy**.
- **RT3-02** `history.json` = `[]` alone is **not** content.
- **RT3-03** ndjson deleted, populated `history.json` present → Healthy (content still loadable).
- **RT3-04** all representations empty/absent with marker `committed` → **Damaged**.
- **RT3-05** marker `pending` + no content → **Pending reconciliation**, never Healthy, never Damaged.
- **RT3-06** marker `pending` + content → promoted to `committed`, Healthy.
- **RT3-07** marker `none` + content → promoted, Healthy.
- **RT3-08** promotion is monotonic: a later absence does not reset it.
- **RT3-09** the scan writes only the marker (no content, no removal).
- **RT3-10** an ID-collision directory (meta of one, ndjson of another) is reported, never repaired.

---

## 8. Left unestablished

| # | Item |
|---|---|
| OQ-3a-a | Frequency of same-second creation collisions in real use. |
| OQ-3a-b | Whether `SaveChatAs` is expected to re-point the app to the copy or to leave the user on the source. |
| OQ-3a-c | Whether a future chat import will reuse `SaveChatAs` (which would inherit P6's representation gap). |
| OQ-3b | Still open: exact per-representation "empty" threshold and per-scan cost. |
| OQ-3c | Partly answered: dual state exists (P3/P5 write both) and is benign — content in either is content. |

**RC3-14 — Contract work only.** No code, no migration, no repair, no deletion. P8's restore semantics are source-derived from rename, not observed.
