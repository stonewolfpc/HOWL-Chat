# Q-06 — Working-Tree Provenance Report
## Howl Chat / R-07 Evidence Follow-Through
### Date: 2026-11-11
### Status: Read-only baseline audit complete for the five R-07 review-critical files; repository-wide baseline remains dirty; no architecture or implementation authorized

## 0. Scope

Q-06 asked which five modified source files relative to Git HEAD were involved, what changed, and whether the changes can be shown to predate R-07 and the Muse independent review.

This report answers that question for the five files that carry the functions independently traced by Muse:

1. `app-prompt.go`
2. `frontend/dist/js/chat-stream.js`
3. `frontend/dist/js/memory-store.js`
4. `internal/backend/chat/persistence.go`
5. `internal/backend/lorebook/lorebook.go`

These are a **review-critical set**, not the total modified-file population. Git shows 91 modified tracked source-ish files (`.go`, `.js`, `.ts`, `.tsx`, `.jsx`), 107 tracked-modified files overall, and 258 dirty entries including untracked files. The premise “five modified source files” is therefore true only for this scoped set. — **Observed**

No source file was edited, staged, committed, reverted, or cleaned. No runtime experiment was performed. — **Observed**

## 1. Repository baseline

- Branch: `master`
- HEAD: `9c7d2a17dd7f702d7db6d95d9316e7640e34bf6f`
- HEAD commit date: `2026-11-05 14:39:18 -0600`
- HEAD subject: `feat(audio): wire the audio settings lifecycle and finish model discovery`
- Working tree: dirty, with no commits after HEAD in the inspected reflog.
- The five files are tracked, modified in the working tree, and unstaged: Git status type `1 .M`, with index blob equal to HEAD blob. — **Observed**

This establishes that the five files’ cited content is not the content of HEAD. It does **not** establish when the uncommitted edits were made or by whom. — **Observed**

## 2. Identity table

| File | Current SHA-256 | HEAD content SHA-256 | HEAD/index blob | Diff +/− |
|---|---|---|---|---|
| `app-prompt.go` | `93c9e6b8759936194a06ee5a49a3c2dd1f16054b20670738c45c42c5940d4128` | `a902c7cec753182b30c218e91cc713963cfef0677d6e55d2c4ab8c58ae6c2eaf` | `63c0a85b5b6842a5fdacda185e48d47a5adf6f83` | `+113 / −35` |
| `frontend/dist/js/chat-stream.js` | `2106be99e088476d58274b10d610034167bb3a6168a064371c5a2367656c377a` | `c6c6fc5557021b8d4cddf5e31c3b8dac418a2caada561434fa37b1d9478abe86` | `8bb0ac779dbf6390894e5890da9af0f96686dce8` | `+40 / −20` |
| `frontend/dist/js/memory-store.js` | `78a8f6b8c76d435ccdf4b58614b5692e7bad19d8c9cc6505436973d24176a4cd` | `65b7c03923a16105086a189756dcf8cdf624476a45a631500da9218c9f157c02` | `470255a91e78badef4cb3926994276ba442f5657` | `+7 / −1` |
| `internal/backend/chat/persistence.go` | `ad9663a48278cff0ee35eb48d28799702cdbfa7754f9f71b7541a3a6762a87c9` | `e7bee8807a3ba4053ddef29a96fdb23ff205a66a6da5c203307af9b1424d1b79` | `bc4c632311fab0ce083935f3ed285a8ecfa48131` | `+19 / −7` |
| `internal/backend/lorebook/lorebook.go` | `1f6b5d5fa18c9aeaa3e6bbe542e3c7dc5574d65ab872e94a0f965122527b069e` | `05030d74566119a19b12324edf17b186e1e115e3a6b11f100bd8756069045d5e` | `ed8aec2643a603eba6de9892c48ce77854f93443` | `+87 / −30` |

Combined diff: `266 insertions`, `93 deletions` across five files. — **Calculated**

## 3. Timestamp and Git-history evidence

### 3.1 Timestamps

| Artifact | mtime / ctime |
|---|---|
| `app-prompt.go` | `2026-11-10 00:58:15 -0600` |
| `frontend/dist/js/chat-stream.js` | `2026-11-10 00:58:16 -0600` |
| `frontend/dist/js/memory-store.js` | `2026-11-10 01:23:17 -0600` |
| `internal/backend/chat/persistence.go` | `2026-11-10 09:20:06 -0600` |
| `internal/backend/lorebook/lorebook.go` | `2026-11-10 09:21:06 -0600` |
| R-07 report | `2026-11-11 16:37:14 -0600` |
| Muse review | `2026-11-11 16:44:02 -0600` |

All five source mtimes precede R-07 and the review by roughly 31–39 hours. This **supports** the claim that the reviewed file versions predate both artifacts. It does not prove authorship, exact edit time, or that timestamps were not altered. — **Observed / Inferred**

### 3.2 Git history

- No commit exists after HEAD `9c7d2a1` in the inspected reflog.
- Last commits touching the five files:
  - `app-prompt.go`: `9c7d2a1`, `2026-11-05 14:39:18 -0600`
  - `frontend/dist/js/chat-stream.js`: `b2fa517`, `2026-11-01 09:40:02 -0600`
  - `frontend/dist/js/memory-store.js`: `b0fe410`, `2026-05-19 18:56:41 -0500`
  - `internal/backend/chat/persistence.go`: `b2fa517`, `2026-11-01 09:40:02 -0600`
  - `internal/backend/lorebook/lorebook.go`: `b0fe410`, `2026-05-19 18:56:41 -0500`

The last-touching commits are all older than the source mtimes. This proves that the current modifications are uncommitted and not represented by those commits. It does not prove when the uncommitted edits occurred. — **Observed**

### 3.3 What the timestamp evidence can and cannot establish

**Established:**

- The five files differ from HEAD.
- Their current mtimes and ctimes precede the creation of R-07 and the Muse review.
- There is no commit after HEAD that could contain the reviewed changes.

**Inferred:**

- The working-tree edits were probably made before R-07 and the review.
- Neither artifact appears to have edited these five files, because their mtimes remained on 2026-11-10.

**Unestablished:**

- Who made each edit.
- Whether all hunks in a file were made in one operation.
- Whether the current bytes are identical to the bytes read by R-07 or Muse; neither artifact recorded per-file hashes before Q-06.
- Whether mtime/ctime values were preserved or altered by copying, sync tools, or external editors.
- Whether the five files are the same five the owner intended by “five modified source files.”

## 4. Exact diffs

Exact Git diff against HEAD for the five review-critical files is preserved separately:

- `docs/research/Q-06-baseline-diffs-2026-11-11.patch`
- SHA-256: `832d08e8857599fbe3b21965ceb75ada2c4f7d722aa30a70c37d35c346307119`
- Size: `26,771 B / 634 lines`
- Format: `git diff --no-ext-diff --binary HEAD -- <five files>`

Validation:

- `git apply --check --reverse` succeeds against the current working tree. — **Observed**
- `git diff --check HEAD -- <five files>` reports no whitespace errors. — **Observed**

The patch is Git’s text diff and is not the same artifact as the full working-tree bytes; the current-file SHA-256 values in §2 identify those bytes. — **Observed**

## 5. What changed in the review-critical five

### 5.1 `app-prompt.go`

Observed changes include:

- `ensureAcceptingWork` guards on send and character-response paths.
- Stream cancellation moves from `streamCancel` to `beginStream`/`endStream` tokens.
- Character responses persist the generating `characterID` through `appendChatMessage`.
- Public `AppendChatMessage` delegates to a character-aware private append.
- History timestamps become RFC3339; `HasImage` and `CharacterId` survive replacement.
- Lorebook activation can attach a canonical entry created after chat save.
- Lore resolution changes from a character name to stable `characterID`.
- `BuildCharacterLoreBlock` is added for frontend character turns.
- `resolveLoreBlock` centralizes budget, trigger bookkeeping, and routing.
- `SetLorebooks` overlays canonical text through `mergeCanonicalLore`.

These changes contain the exact stable-ID lore routing the Muse review traced. They do not add conversation/chat scope to lore or memory. — **Observed**

### 5.2 `frontend/dist/js/chat-stream.js`

Observed changes include:

- Lorebook repository hydration before lore sync.
- New `requestCharacterLoreBlock(characterID, userMessage)` calling backend `BuildCharacterLoreBlock`.
- Local lore reconstruction removed from `buildCharacterSpecificMessage`.
- The composer now accepts a backend `loreBlock` argument.
- Multi-character streaming requests the lore block; on failure it reports the error, breaks the queue loop, and skips the turn.
- Frontend `MemoryStore.getCharacterMemories(character.id)` injection remains.
- `[RECENT EVENTS]` and inactive-character notes remain part of the composed message.

This is the composer that carries the M-01/M-02 path. — **Observed**

### 5.3 `frontend/dist/js/memory-store.js`

Observed changes include:

- New records carry `ownerKind`.
- Backend sync payload carries `ownerKind`.
- `getCharacterMemories` now requires `ownerKind === 'character_id'` and matching `ownerId`.

The changed predicate has no conversation or chat term; this diff does not add conversation scoping. — **Observed**

### 5.4 `internal/backend/chat/persistence.go`

Observed changes include:

- Adds `RootDir()`.
- Isolated mode redirects the legacy chat root to `StatePath("persisted-chats")`.
- Legacy migration no longer invents world/scenario IDs from display names.

No transaction, commit marker, history reconciliation, or backup change appears in this diff. — **Observed**

### 5.5 `internal/backend/lorebook/lorebook.go`

Observed changes include:

- Adds visibility constants, `OriginLegacyGeneratedUnattributed`, `OwnerName`, `OwnerKind`, `Stale`, and `OriginKind`.
- `canConsider` excludes stale and hidden entries.
- `RouteResolvedEntries` changes from character-name matching to stable `characterID`.
- Unowned character entries pass only when `shareableWithoutOwner`.
- World and scenario cases remain unconditional.
- `scenarioName` and `worldName` remain unused in the function body.

This is the exact inert active-world/scenario check the Muse review identified. — **Observed**

## 6. Provenance verdict

| Question | Answer | Class |
|---|---|---|
| Are there exactly five modified source files in the repository? | No. There are 91 modified tracked source-ish files; five is the review-critical set. | Observed |
| Do the five files differ from HEAD? | Yes, all five are tracked and unstaged. | Observed |
| Are the current bytes identified? | Yes, by SHA-256 in §2. | Observed |
| Are exact diffs preserved? | Yes, in the patch artifact. | Observed |
| Do the five modifications predate R-07 and Muse review? | Timestamp evidence supports yes; git history cannot prove it; byte identity at review time is unestablished. | Observed / Inferred / Unestablished |
| Was the working tree clean? | No. | Observed |
| Can the five files be treated as HEAD behavior? | No. | Observed |
| Can future reports inherit a clean-baseline claim? | No. | Observed |

## 7. What remains unresolved

- Q-06 does not establish who authored the uncommitted changes.
- Q-06 does not establish that all changes in each file belong to one coherent change set.
- Q-06 does not establish that R-07 and Muse read byte-identical files; no pre-Q-06 per-file hashes were recorded.
- Q-06 does not enumerate the other 86 modified source-ish files or the 107 tracked modifications.
- Q-06 does not authorize Gate B, runtime experiments, or implementation.

## 8. Artifact identity

This report’s own SHA-256 and size are reported in the delivery message, not embedded here, because embedding them changes the content.

| Artifact | SHA-256 | Size |
|---|---|---|
| Exact five-file diff | `832d08e8857599fbe3b21965ceb75ada2c4f7d722aa30a70c37d35c346307119` | `26,771 B / 634 lines` |
| R-07 original | `e3860e3d5d6ca71724d31d0021897819c968afc1b8e1aad09fc011f8d9a67134` | `46,107 B / 876 lines` |
| Muse review | `4bdb399cfcb56c75e32d3103a97fd57b06db8895ea3fd8dfef16feba5a98d792` | `15,423 B / 180 lines` |

**Q-06 conclusion:** the five review-critical files are modified, unstaged, and timestamped before R-07 and the review, but the repository is not a clean baseline and byte identity at review time is not independently established. Future evidence must cite these hashes and the dirty-tree caveat, not a clean-HEAD assumption.
