# TASK-03 — ledger decision record (D-1 … D-4)

| Field | Value |
|---|---|
| Task | Phase 1B — conversation authority |
| Status | **Decided for review. Implementation NOT authorized.** Alternatives preserved, open questions not forced |
| Basis | `TASK-03-CONVERSATION-AUTHORITY-PREPARATION-2026-11-22.md`; FC-1 §1–§11 (PROPOSED); ADR-001 §4, §5 D, §8; ADR-002 §4 |
| Owner constraints | no sixth authority · no assumption that DurableStore fits · recovery never silently discards valid history · retention is a policy decision · no seventh service |
| Recovery token | `TASK-03-LEDGER-DECISION-2026-11-22` |

---

## 1. Workload evidence — measured, and it does not exist

**Method.** Counted records in every `history.ndjson` under the resolved chat root and under the legacy in-repo
roots; parse-checked each line to count malformed records. Counts and sizes only — no message content was read or
recorded.

| Location | Result |
|---|---|
| `~/Documents/Howl/Chats` (resolved `content:chats`) | **absent** |
| `config/chats`, `config/chat` | absent |
| `./Chats` (inside the installation root) | **2 directories, `memories.json` + `memories_backup.json` only — no history file** |
| `%APPDATA%\Howl`, `%LOCALAPPDATA%\Howl` | absent |

**Result: zero conversations, zero ledger records.** `p50`, `p95` and `max` are therefore **`Unestablished`, not
zero** — the difference matters: "no data" is not "no workload".

**Observed vs projected.** Observed = 0. Any projection (turns per conversation, interval between turns) would be an
assumption and is labelled as such wherever this record uses one.

**Measurement procedure to run once data exists** (the numbers that decide §5):

```
for each conversation directory under the resolved content:chats root:
    count non-blank lines in history.ndjson   → records per conversation
report p50, p95, max, mean, and the number of conversations above 2,000 and 10,000 records
```

**Incidental finding, recorded because it is evidence for a closed task.** App-owned chat directories exist under the
**installation root** (`./Chats`), which is precisely the `ProjectRoot`-as-state-root defect `TASK-01` characterised
and deliberately did not migrate (limitation L-5). Which code path created them is **unverified**.

## 2. D-1 — Representation: evolve `history.ndjson`, do not replace it

The question is not which is cheaper to write. It is which representation can carry the required authority and
recovery guarantees. Both candidates were scored against the guarantees, not against edit volume.

| Guarantee (FC-1 / ADR-001) | **E2 — evolve `history.ndjson` format, commit by whole-document replacement** | **E1 — evolve into a true append journal (+ framing/checksum, new helper)** | **R — replace with a new versioned document** |
|---|---|---|---|
| Single authoritative representation | Yes — one file per conversation, already scoped per conversation | Yes | Yes, but it is a **new** file beside the old one → sixth authority risk |
| Committed events never edited (FC1-09) | Yes — prior events are re-emitted byte-identical; one commit = previous content + one record | Yes | Yes |
| Torn append cannot exist | **Yes, structurally** — the live file is only ever replaced by rename; a torn *temp* is discarded | Yes, but only via explicit framing + truncation rules | Yes |
| Crash-safe durability | `DurableStore` already verified for this (atomic rename, pre-images, hash recovery) | Needs a new primitive, unbuilt and unverified | `DurableStore` |
| Reader sees a complete document | Yes, verified primitive | Requires new reader rules for a partially written log | Yes |
| Duplicate suppression | Idempotency key stored in the record; lookup on commit | Same | Same |
| Legacy records remain readable | **Yes — line-oriented, additive fields; old lines parse** | Yes | No — needs a migration or a dual-read path |
| Cost of each commit | **O(n)** — whole file rewritten, plus 2 pre-images (3× write amplification) | O(1) append | O(n) |
| New architecture required | **None** | A new primitive | A new document + migration |

**Decision: E2.** The ledger keeps `history.ndjson` as its file — one JSON object per line, per conversation — and
gains a stable envelope per record. Commits go through `DurableStore`, which replaces the whole file atomically.

**Why not E1 (kept as the upgrade path, not discarded).** E1 is the better write model and may become necessary;
it is rejected *now* only because it requires a primitive that does not exist, cannot be verified this branch, and
would be introduced without workload evidence (§1). Its cost advantage is real and is preserved in §5.

**Why not R.** Replacing the format creates exactly the duplicate authority the owner prohibited, and forces a
migration for data that does not yet exist to justify it.

**Consequences of E2, stated plainly.**

- **Torn appends stop existing as a failure class.** Nothing ever writes into the live ledger; the only mutation is a
  rename. This removes the entire framing/checksum/truncation design from Phase 1 rather than deferring it.
- **Write cost is O(n) and storage is ~3× the ledger** (live + two pre-images). That is a calculated cost, not a
  measured one; §5 sets the threshold at which it becomes unacceptable.
- Legacy lines without the new envelope are read as **legacy records** — preserved read-only, never assigned
  invented identities (FC1-07 stays unresolved).

## 3. D-2 — Ordering and writer authority (`FC1-12`, `FC1-CONCURRENCY`)

| Rule | Decision |
|---|---|
| Serialization | One exclusive lock per conversation, on `<conversation>/ledger.lock`, reusing the `LockFileEx` discipline proven in `TASK-02` |
| Sequence assignment | Under the lock, `sequence = last valid record's sequence + 1`. Never computed from wall-clock, never from array position |
| Gap policy | A gap on read is an **error**: the ledger is reported `unresolved` from that point, with the gap named. Sequences are never renumbered and gaps are never filled |
| Multi-process writer authority | **Refused.** A second process reports contention (`Busy`), never a merge, never a silent append |
| Concurrent appends | Serialized by the lock; no event may be committed outside it |
| Conflict handling | None exists in Phase 1, because there is no concurrent writer. Conflicts are prevented, not resolved |

**Compatibility note.** This is the same rule set as `TASK-02`'s lock, deliberately: one locking discipline in the
codebase, not a second one invented for the ledger.

## 4. D-3 — Selection after deletion (`FC1-SELECTION-AFTER-DELETION`)

**Invariant (fixed by FC1-13a):** `current_state_id` never dangles; no silent reselection; no resurrection; no
manufactured identity.

| Candidate | Behaviour | Verdict |
|---|---|---|
| A. Nearest surviving ancestor | Walk the state chain upward from the deleted state; select the first surviving state; record a `selection_resolved` event naming the deleted state, the chosen state and the rule | **Selected** |
| B. Explicit empty conversation | Selection resolves to an explicit empty state; retained siblings remain reachable through history | Kept as the fallback when **no** ancestor survives |
| C. Unresolved marker | Selection becomes `unresolved` and requires a user choice | Rejected for the normal case: it makes an ordinary deletion a blocking dialog |
| D. Silent fallback to the most recent state | Simple | **Prohibited** by FC1-13a — it selects unrelated state without recording why |

**Decision: A, falling back to B when A has no candidate.** Either way the resolution is an **event**, so a restart
reproduces it by replay rather than by re-deriving with different code.

**Frontend consequence, and it is not optional.** Selection currently lives in `localStorage`
(`howl-chat-current-chat`). It must become a **cache of the ledger's projection**, written after a successful read
and invalidated when the ledger says otherwise — never a second source of truth. Removing its authority is part of
the branch (acceptance A-5, A-10).

## 5. D-4 — Recovery, and the journaling-primitive verdict

### 5.1 Failure semantics

| Failure | Required behaviour |
|---|---|
| **Torn append** | Cannot occur in E2: the live ledger is only replaced by rename. A torn temp is discarded; the previous ledger is intact |
| **Duplicate submission** | Same `idempotency_key` + identical payload → return the existing result, one durable effect, no second record |
| **Lost acknowledgement** (commit landed, caller did not see it) | Retry with the same key finds the record by key and returns `Saved` with the existing event id. The caller must treat this as success, not as a new commit |
| **Conflicting reuse of a key** | Explicit error. Never a silent overwrite (FC1-10) |
| **Damaged record mid-file** | Parsing **stops at the damage**. The prefix before it stays valid and readable; everything from the damage onward is reported `unresolved` and named. **No valid history is discarded, and no later line is silently skipped** |
| **Damaged tail** | Same rule, with the additional statement that it is a tail. Nothing is rewritten to "repair" it |
| **Hash mismatch against the manifest** | `Recovered(from pre-image)` or refusal — the `DurableStore` semantics already verified, with `Result.Current()` distinguishing a recovered generation from the current one |

**The non-negotiable rule:** recovery may *refuse*, may *report a prefix as valid and the rest as unresolved*, and
may *recover a retained generation* — it may never silently drop records that are valid, and it may never present a
shortened ledger as complete.

### 5.2 Is a journaling primitive required? — **Not on current evidence. Not now.**

| Option | Cost of being wrong |
|---|---|
| Build the primitive now | Real: an unverified concurrency/durability component, plus a new abstraction, with **zero** workload evidence justifying it — exactly the expansion the owner prohibited |
| DurableStore as-is (E2) | Real but **bounded and measurable**: O(n) per commit, ~3× storage. Each commit is a whole-file rename of a conversation-sized document |
| Wait for evidence | Selection deferred, ledger not started — not acceptable once the rest is decided |

**Decision: E2 now, with a written upgrade trigger.** The journaling helper is **not** built, **not** planned as a
service, and **not** a seventh component. If it is ever justified it is a private implementation detail behind the
ledger's commit interface — the ledger stays the only owner of conversation state either way.

**Upgrade trigger (measurement gate, not a preference):** implement the helper when any of these is observed —

- p95 records per conversation **> 2,000**, or
- max records per conversation **> 10,000**, or
- p95 ledger size **> 8 MiB**, or
- a measured commit latency above the turn budget on the reference machine.

**Calculated context for those thresholds** (arithmetic on assumed record sizes, *not* measurements): at ~2 KB per
record, 2,000 records ≈ 4 MiB, so each commit rewrites ~4 MiB and holds ~12 MiB across live + pre-images; 10,000
records ≈ 20 MiB, ~60 MiB including pre-images. Those are the points at which whole-document replacement stops being
a reasonable per-turn cost.

## 6. Retention — deleting a conversation and its ledger (FC1-22c)

This is an ownership and recovery decision, and FC-1 leaves the chat-scoped case explicitly unresolved
(`FC1-STATE-RETENTION`); FC1-22c directs the lifecycle only for **snapshots**.

| Option | Consequence |
|---|---|
| **A. Ledger deleted with its conversation, plus a retained deletion receipt** | Matches FC1-22c's directed lifecycle for conversation-owned data; the receipt proves what was deleted; no orphan state; but conversation history is gone, and an accidental deletion is unrecoverable |
| B. Ledger retained indefinitely after conversation deletion | History survives; but it is orphaned conversation state with no owner, which is precisely what `FC1-CHAT-SCOPED` refuses to classify, and it retains private material after the user asked for deletion |
| C. Retention window then delete | Middle ground, but the window is a policy nobody has set, and it creates an unattended deletion path that FC1-16 forbids |

**Recommendation: A**, with two conditions — the receipt is retained (identity, time, authority, record count), and
**authored content is untouched** (FC1-22b). Option A is recommended because a deleted conversation that keeps its
ledger is a second lifecycle the contract does not authorize.

**This fills an `UN` field, so it is a decision for the owner, not something this record can settle.** If the owner
prefers B or C, only the receipt/retention portion of this record changes; D-1…D-4 do not.

## 7. Preserved alternatives and open questions

**Alternatives preserved, not discarded:** E1 (true append journal) as the §5 upgrade path · selection candidate B
as the fallback and C as a possible future mode · retention options B and C pending the owner's call.

**Open, and deliberately not forced to consensus:**

| Question | State |
|---|---|
| `FC1-02` — build the state/branch graph now, or a linear chain with the graph derived later? | Unresolved in this record; Phase 1 can proceed linear, but that must be stated, not assumed |
| `FC1-08` — which envelope fields are mandatory in Phase 1 | Proposed field set in this record; not ratified |
| `FC1-11` — minimum event kinds | Proposed; not ratified |
| `FC1-14b` / `FC1-14c` | Out of scope: no branching, one turn = one message in Phase 1 |
| `FC1-07` / `FC1-LEGACY` | Unresolved; legacy read-only, no invented IDs |
| `FC1-STATE-RETENTION` | The subject of §6; needs the owner |
| `FC1-SNAPSHOT-CONTENT`, `FC1-PRIVACY-DELETE`, memory-copy clauses | Out of scope for TASK-03 |

## 8. What this record does not authorize

No code, no migration, no call-site change, no contract approval. FC-1 remains `PROPOSED`; implementing a clause
does not approve it. TASK-03 remains unimplemented.

**Next:** owner review of this record as a whole. On approval, implementation may be authorized; the first code is
the ledger package and the retirement of the four competing implementations, in that order.

---

## 9. Amendments after owner review (appended 2026-11-22)

Sections 1–8 stand as written. This section adds what the review required and **expands §6 with the exact
definitions it only summarised**. No earlier text was rewritten.

### 9.1 D-1 condition — physical replacement must never become logical rewrite

**Clause.** Whole-document replacement preserves the logical append-only event history. Physical replacement is an
implementation detail of durability; it grants **no** permission to rewrite, reorder, re-encode, deduplicate or
remove a committed event.

**How it is enforced.** A commit constructs the new document by reading the current committed bytes and appending
exactly one record. It must **not** serialise the ledger from an in-memory projection, because that would re-encode
prior records (key order, whitespace, float formatting) and silently alter history. Prior records are carried as
bytes.

**How it is observed (test, not promise).** After N commits, the first N records of the ledger must be
**byte-identical** to the prefix of the file as it stood after commit N-1. A second, independent check compares
against the retained pre-image, whose hash was recorded in the manifest.

**Corollary — no compaction in Phase 1.** No deduplication, no "fix-up" pass, no re-serialisation, no truncation. If
compaction is ever introduced it is a separate authorised act with its own receipt, and it must preserve every event
or record an explicit supersession — it may not delete history to save space.

### 9.2 The recovery guarantee becomes an acceptance criterion

Added to §8's table. The sentence is now testable:

| # | Criterion | Observable form |
|---|---|---|
| **A-12** | A damaged record produces a **valid prefix, an exact boundary, and no modification** | Given a ledger of n records with record k damaged: `Load` yields records 1…k-1 as valid, reports the ledger `unresolved` **from k**, names k, and the file's hash is **unchanged after the read**. A test asserts all three, including that the reader performed no repair |
| **A-13** | A shortened ledger is never presented as complete | Any result that omits a record existing before the read must report a shorter boundary and an `unresolved` status. A test asserts that a "successful, complete" report is impossible when the boundary is not n |

**Prohibited, and now checkable:** silently skipping a damaged record and continuing (which would report post-damage
records as if the ledger were intact); rewriting the file to make it parse; treating a truncated read as success.

### 9.3 Installation-root chat directories — separated as a migration finding

Recorded as **MF-1**, deliberately **not** part of TASK-03.

| Field | Value |
|---|---|
| Observation | `./Chats` exists inside the installation root, containing two conversation directories with `memories.json` and `memories_backup.json` and **no** history file |
| Provenance | **Unverified.** Which code path created these directories is not established |
| Relation to closed work | Consistent with TASK-01's finding that `StateRoot()` defaults to `ProjectRoot()` (deferred there as L-5) |
| TASK-03 boundary | TASK-03 establishes the **conversation authority and its storage binding** to the PathService-resolved chat root. It does **not** migrate, move, delete or clean up these directories |
| Preservation requirement | The two `memories.json` directories are **preserved untouched** pending a migration decision. They are not "stale state" and not orphans to sweep — `FC1-16` and `FC1-23a` apply: nothing is deleted for being unreferenced or of unknown lifecycle |
| Next owner of this finding | A separate migration task, with its own inventory, preservation check and receipts |

### 9.4 Retention — the exact Option A, B and C definitions

§6's one-line entries are expanded here, because "retain" and "window" hide different consequences.

#### Option A — ledger deleted with its conversation, plus a receipt

| Aspect | Definition |
|---|---|
| Deleted | The conversation's ledger, its pre-images, its manifest and its snapshots (FC1-22c). The conversation becomes unselectable and unrecoverable |
| Retained | One **deletion receipt**; nothing else |
| Recovery after deletion | **None.** The history is gone and the app must say so. An interrupted deletion reports partial/unknown and never claims completeness (FC1-29) |
| Privacy | **Strongest.** No conversation content or title survives the operation |
| Receipt contents (proposed) | Deletion-operation id · timestamp · authority · record and snapshot counts · the **last record's hash** as an end-of-history marker · a **salted hash of the `conversation_id`**, not the id itself |
| Does the receipt identify the conversation? | Only to someone who already holds the id and the salt — not to a reader of the receipt. **A human-readable title is not retained.** If the product must display "you deleted N conversations", that label reintroduces identifying data and becomes a separate privacy decision (`FC1-PRIVACY-DELETE`) |
| Contract compatibility | Requires **no** amendment. It is the directed lifecycle FC1-22c already gives conversation-owned data, extended to the ledger |
| Storage after deletion | A fixed-size receipt. No growth |

#### Option B — retain the ledger after the conversation is deleted

| Aspect | Definition |
|---|---|
| Deleted | Nothing of the ledger. The conversation is marked deleted (tombstone), hidden from selection, and not listed |
| Retained | The full ledger **and its snapshots** — retaining snapshots is required for the retained history to be usable, and that is the clause B strains against |
| Recovery after deletion | **Possible.** The conversation can be reconstructed and, if the product allows, undeleted |
| Privacy | **Weakest, and arguably misleading.** A user who deletes a conversation and still has every message on disk has not deleted it, whatever the interface says |
| Receipt contents | A tombstone record inside the retained ledger; no separate receipt is needed, but the ledger must remain writable by *some* owner — and a deleted conversation has no lifecycle owner |
| Does the receipt identify the conversation? | Yes, fully — the whole history is the identifier |
| Contract compatibility | **Contradicts FC1-22c**, which is owner-approved (D2) and directs that conversation-owned snapshots are deleted with the conversation. B therefore needs an explicit amendment to FC1-22c, and it leaves `FC1-CHAT-SCOPED` (unclassified orphan-like records) doing work it was not written for |
| Storage after deletion | Unbounded; grows with every deleted conversation |

#### Option C — retain for a window, then delete

| Aspect | Definition |
|---|---|
| Deleted | After a retention window elapses, the ledger and snapshots are deleted; until then everything is retained as in B |
| Retained | Everything, for the window — 30 days, 90 days, or whatever the owner sets |
| Recovery after deletion | Possible **within** the window, irreversible after it, with the boundary enforced by a clock |
| Privacy | Time-limited exposure; the user's deletion is not effective until the window expires |
| Receipt contents | Two events: the deletion request, and the expiry deletion, each with its own operation id and counts |
| Does the receipt identify the conversation? | Until expiry the full history is present; after expiry the receipt follows A's shape |
| Contract compatibility | **Contradicts FC1-16**, which authorises no independent background deletion path and no sweep. A scheduled expiry deletion is exactly such a path unless the owner authorises it as a fourth directed deletion form. It also needs a window that nobody has set, and it depends on the application running after expiry — so the guarantee is conditional |
| Storage after deletion | Bounded by the window |

#### What the owner is actually choosing between

| Question | A | B | C |
|---|---|---|---|
| Does deletion delete history? | Yes | No | Eventually |
| Can a deleted conversation be restored? | No | Yes | Only inside the window |
| Does it need a contract amendment? | **No** | Yes — FC1-22c | Yes — FC1-16, plus a chosen window |
| Does the receipt carry identifying data? | No (salted hash) | Yes (the history itself) | Yes until expiry |
| Unattended deletion path introduced? | No | No | **Yes** |

**Recommendation unchanged: A.** It is the only option that requires no amendment, introduces no unattended
deletion, and makes the user's deletion mean what it says. **The decision remains the owner's**; B and C are now
defined precisely enough to be chosen deliberately rather than by default.

### 9.5 What still blocks freezing the ledger contract

Not decided here, and needed before the contract is frozen: event envelope fields (`FC1-08`) · Phase 1 event kinds
(`FC1-11`) · graph scope, linear or full graph (`FC1-02`) · legacy record handling (`FC1-07`/`FC1-LEGACY`) ·
branch and turn/message relationships (`FC1-14b`, `FC1-14c`) · the retention choice in §9.4.

These are dispositions to be prepared and reviewed, not gaps to be filled by implementation convenience.

### 9.6 Authorization

Still none. No code, no migration (including MF-1), no call-site change, no contract approval. FC-1 remains
`PROPOSED`.

---

## 10. Retention — owner decision recorded (2026-11-22)

§6 recommended A and §9.4 expanded A, B and C. **The owner has now decided. This is a decision, not a
recommendation.**

### 10.1 Decision

**Option A is selected. Options B and C are rejected.**

> Deleting a conversation deletes its conversation-owned ledger and snapshots together. No retained conversation
> ledger after deletion. No delayed-retention window. No independent background expiry or sweep. Authored content
> remains outside this deletion boundary unless separately governed. A deletion receipt is preserved, subject to
> the privacy constraints in §10.2.

**B and C are preserved as rejected alternatives, with their definitions and rationale intact in §9.4**, because the
reasons they were rejected are part of the record: B contradicts the owner-approved FC1-22c snapshot lifecycle; C
introduces the independent background deletion path FC1-16 forbids, and needs a window nobody set.

**What "deleted together" means, precisely.** The ledger, its pre-images, its manifest, its lock file and its
snapshots are removed as one conversation-scoped operation. The receipt survives that operation.

**What it does not mean.** Option A governs conversation-owned data only. It does **not** authorise deleting
authored content (FC1-22b), and it does **not** silently govern chat-scoped memory records whose lifecycle is
FC1-22e's. Observed today: `DeleteChat` is `os.RemoveAll(chatDir)` (`persistence.go:471–479`), which removes
*everything* in the conversation directory — including `memories.json` — with no receipt and no report of what was
removed. Implementing A must make the deletion **enumerated and receipted**, not widen a directory removal.

### 10.2 Deletion receipt — fields retained for investigation, and the privacy questions

**Minimal fields, revised after owner review** — the exact table is in
`TASK-03-LEDGER-SCHEMA-PROPOSAL-2026-11-22.md` §5:

- **Required:** operation ID · timestamp · authority · record count · snapshot count.
- **Optional, proposed default absent:** last record's hash. It is a **fingerprint** — anyone holding a surviving copy
  of the ledger can match it and identify the deleted conversation — so it must earn its place operationally or be
  dropped.
- **Withdrawn:** the salted conversation-ID hash (§10.2.1).
- **Prohibited:** conversation title · plaintext conversation ID.

Removing the identifier did **not** make the remainder non-identifying: timestamp and counts still narrow a receipt
to a particular deletion. That residual correlation is stated rather than hidden; it is not a fingerprint of
content, which is precisely the difference the optional hash removes.

**Explicitly not added without resolving `FC1-PRIVACY-DELETE`:** UI labels, deleted-conversation counts, or any
other identifying metadata.

#### 10.2.1 The salted hash is not automatically anonymous — and here it is weak

A salted hash is only a pseudonym, and a pseudonym is only as strong as its input space. **Observed:** conversation
directories are named by timestamp (`Chats/2026-05-27T16-46-11Z`, `Chats/2026-05-29T12-28-48Z`). Anyone who can
guess a timestamp-shaped ID can confirm a guess against a receipt, because the salt is stored beside the hash. The
hash therefore behaves like an ID for any conversation whose identifier is enumerable — which is **all of them**
today.

| Question | Answer / state |
|---|---|
| Who controls and retains the salt? | **Undecided.** No component owns a salt today; nothing in the tree stores one |
| Global, scoped, or per-operation? | Per-operation is the least linkable, but **storing it beside the hash defeats scoping**: anyone holding the receipt can recompute the hash for a candidate ID |
| Can the receipt be linked to a conversation by anyone holding ID + salt? | **Yes — by construction**, and for timestamp-named conversations, by anyone at all |
| Is retaining the salt necessary to prove deletion? | Only if per-conversation proof is required. **No clause requires it** |
| Does the receipt create a persistent identity or reconstruction path? | With a stored salt and enumerable IDs: **yes, a verification path**. Without any identifier: no |

**Consequence for the design:** the earlier "salted hash" proposal is **doubted by this investigation** and is not
adopted. The candidate that survives is a receipt **with no conversation identifier at all** — operation ID,
timestamp, authority, counts, last-record hash — which proves a deletion occurred without creating a lookup path.
Whether per-conversation proof is needed is an owner question, not an implementation choice. See the disposition in
`TASK-03-FC1-DISPOSITIONS-2026-11-22.md` (`FC1-PRIVACY-DELETE`).

### 10.3 D-1, A-12 and A-13 restated as enforceable requirements

**D-1 — historical byte preservation.** Existing records are carried forward **as bytes**. Historical records are
never re-serialised from an in-memory projection. **No compaction in Phase 1.** Prefix byte identity is tested after
commits, and the retained pre-image is an independent verification source.

Three properties that must not be conflated:

| Property | What it means | What establishes it |
|---|---|---|
| **Byte preservation** | the committed bytes are unchanged | prefix byte-identity test; pre-image hash |
| **Semantic validity** | the record means what it meant when committed | parsing under the recorded format version |
| **Integrity validation** | the record is the record that was committed | checksum / manifest hash |

A byte-identical record can still fail to parse (damaged at write time, or written by a newer format), and a
parsing record can be semantically stale. The ledger reports these **separately**; success at one is never reported
as success at another.

**A-12 — damaged-record recovery.** Return the valid prefix. Identify the **exact** damaged boundary. Mark the
remainder `unresolved`. Do not modify or repair the source during reading. Verify the source file hash is
**unchanged after the read**.

**A-13 — truncation.** A shortened ledger is **never** reported as complete. **Completeness is an explicit field of
the result**, never inferred from a successful parse or from a record count. A valid prefix must not be mistaken for
the complete conversation.

### 10.4 MF-1 — scope separation restated

TASK-03 binds the conversation authority to the **PathService-resolved root**. It does **not** migrate, move, clean
up, or claim ownership of the installation-root `./Chats` directory. Both observed `memories.json` directories
remain **untouched**, and are not to be classified as stale or sweepable. **MF-1 stays outside TASK-03.**

