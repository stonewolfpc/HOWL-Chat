# MEMORY-09 — The `chat_id` Drop: A Record of Architectural History

**Status:** historical record. Read-only. No schema change, no restoration, no implementation, no Git
operation.

| Field | Value |
|---|---|
| Date | **2026-11-19** (local, CST) |
| Purpose | Record that the v1 design carried `chat_id` on events and r32 does not — **without** treating that as a defect or a restoration mandate |
| Baseline | `DESIGN-MEMORY-CONTINUITY-2026-11-11.md:107` · r33 `dfb28b6003ca356f` |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** ·
**Superseded**.

---

## 1. The observation

**CH-01 (Observed).** The v1 `Event` schema (`DESIGN-MEMORY-CONTINUITY-2026-11-11.md:96-112`) is the
following — fields verbatim and in original order, reformatted from one-per-line to a single list:

```
event_id · continuity_id · chat_id · turn_id · source_message_id · event_type ·
participants[] · summary · evidence_span · observed_at · effective_at ·
origin · confidence · status · supersedes_event_id?
```

**CH-02 (Observed).** r32's `AcquisitionEvent` (§6.2) is:

```
id · character_id · memory_id · turn · via · from_character_id · source_memory_id ·
confidence_in_source · evidence_span
```

**CH-03 (Calculated).** Four v1 **event** fields have no direct r32 counterpart on the event, and they are
not random:

| v1 field | r32 | Character |
|---|---|---|
| `chat_id` | **absent from the event, present on `MemoryRecord` as `owner_chat_id`** | event-level namespace, **relocated to the memory** |
| `source_message_id` | **absent** | message-level provenance |
| `observed_at`, `effective_at` | **absent** | wall clock, which K-R15 explicitly forbids |
| `continuity_id` | **absent** | MEMORY-09's scope is the chat |

**CH-04 (Inferred).** That is a **coherent simplification**, not an oversight: r32 removed the entire
temporal apparatus (`observed_at`, `effective_at`), moved chat identity onto the memory as the
"Ownership (chat)" dimension (§5), and dropped message-level provenance while keeping the span. Each
removal is consistent with K-R15's "no wall clock" and with a chat-owned record model.

**CH-04a (Source-derived).** Chat identity is therefore **recoverable per event by join**:
`event → memory_id → MemoryRecord.owner_chat_id`. The namespace exists; its direct carrier on the event
does not.

---

## 2. Why it looks like a regression

**CH-05 (Observed).** ARCH-03 **F6** records the absence as *"`AcquisitionEvent` has no `seq` (and no chat
field)"* — phrased as though the chat field had never been designed. It had.

**CH-06 (Inferred).** That phrasing overstates the loss. The namespace was not removed — it was moved to
the record that is chat-owned. F6 remains **materially correct** for the event schema (`seq` has no host
and chat is not on the event), but its parenthetical should be read as *"the field is not on the event"*,
not *"the field does not exist"*.

---

## 3. Why it should not be restored automatically

**CH-07 (Observed).** `MemoryStore` is constructed per conversation (`conversation_id`), one store per chat.

**CH-08 (Inferred).** Under a single-chat store, a `chat_id` on every event is **redundant**: the store's
identity *is* the chat. And even if a store held several chats, the chat is recoverable by join (CH-04a).
So K-R15's "monotonic **per-chat** counter" has a workable namespace either way; what it lacks is a
direct carrier and a value.

**CH-09 (Source-derived).** Cross-chat movement is out of MEMORY-09's scope: K-R46–K-R48 are outside its
authority under O-13/D-4, and no transfer path exists in either implementation.

**CH-10 (Inferred).** On the evidence available, the drop is **not a defect** — it is a relocation. Restoring
`chat_id` to the event would add a field whose value is derivable, to support a scope the project has
explicitly excluded. That is a reason to leave it, not to restore it.

---

## 4. What the record is for

**CH-11 (Inferred).** Three uses, none of them restoration:

1. **Accuracy of history.** K-R15's "per-chat" phrase now has no field-level carrier; reading it as "per
   store" is correct and should be stated rather than assumed.
2. **Accuracy of F6.** F6's parenthetical should be understood as *"the field was removed"*, not *"the field
   was never specified"*. The finding itself stands; only its provenance changes.
3. **A caution for D-1/D-2.** If `evidence_span` is redefined to include message identity, it would
   reintroduce `source_message_id` under another name — legitimately, but that should be a decision, not a
   side effect.

**CH-12 (Unestablished).** Whether the drop was deliberate. No record of the decision was found; the
pattern of the other three removals suggests it was, but that is inference from coherence, not evidence.

---

## 5. What this record does not do

**CH-13 (Observed).** It restores nothing, proposes nothing, and changes no artifact. It does not reopen
K-R15, does not affect DN-1/DN-4, and does not bear on the forgetting-provenance defect.

---

## 6. Verification

v1 design document **not modified**. r33 `dfb28b6003ca356f` unchanged. Tests executed: **0**.

**End of record.**
