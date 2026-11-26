# ADR-001 — Howl Chat Foundational Backend Rebuild

**Status:** Accepted as architecture decision record (owner-authorized documentation, 2026-11-22).
**Implementation:** **Not authorized.** No code, schema, migration, test, build, launch or contract amendment follows from this record.
**Frozen boundaries untouched:** FC-5 revision 3 (frozen), FC5-31 (`Unresolved`), R6 draft (unapplied).
**Companion record:** `ADR-002-STORAGE-AND-NAVIGATION-POLICY-2026-11-22.md` (storage locations, picker navigation, content routing).
**Owns:** storage authority (the five roots), configuration authority and its routing record, persistence lifecycle, runtime resolution and evidence levels, ownership classes and external-write grants.
**Does not own:** storage locations per category, picker navigation, or content routing — those belong to ADR-002, which refines this record's D1 for chats only. Neither record overrides the other by silence.
**Discoverability:** linked from `HOWL-MEMORY-09-PROGRAM-DIRECTION-2026-11-21.md` §2.3 and from the HCDM architectural contract register (`HCDM-2026-11-21-01` §F). **A link establishes discoverability, not authority.**

## Evidence labels

`SC` source-established · `AO` archived observation · `RC` runtime-confirmed · `CA` calculated · `IN` inferred · `UN` unestablished.

Static source evidence is never promoted to a runtime guarantee. Where a claim needs an execution, it is marked `UN`.

---

## 1. Context — why the foundations, not the features

The defects found in this review are not distributed randomly. They cluster where one responsibility has more than one implementation.

| Foundation | Demonstrated conflict | Evidence |
|---|---|---|
| Storage authority | Three resolution rules: `ProjectRoot()` (exe search, then cwd), `StatePath()`, and literal relative joins | `internal/paths/paths.go:29–54`, `startup_log.go:26`, `audio/pipeline.go:1057,1074`, `asr/whisper/recognizer.go:223,242` `SC` |
| Configuration authority | Four overlapping authorities with *temporal* precedence: unified settings, legacy model/sampler files, frontend `localStorage`, live in-memory maps | `app-chat.go:698–738`, `app.go:271–306`, `settings-state.js` `SC` |
| Persistence lifecycle | Reset cannot clear sampler/model settings while legacy files exist; "saved" means only "persisted" | chain `settings-state.js:94–95` → `app-chat.go:705–738` → `app.go:271–306` `SC` |
| Runtime resolution | Intent, derivation, emission and adoption are one undifferentiated value | `-c 3096` emitted from a value the user never chose `AO`+`CA`; `stop_sequences` never reaches `req.Stop` (string vs `[]string`, `app-prompt.go:1001,1006`) `SC` |
| Ownership | Implicit external mutation: image-page load writes sidecars into user model directories; rescan deletes by filename suffix | `image-generator-core.js:202–203`, `internal/image/scanner/{scanner.go:97–110,manifest.go,index.go}` `SC` |

**Scope of the rebuild:** those five foundations, behind stable interfaces. Retained unchanged: frontend, audio contract and its per-setting status vocabulary (`internal/audio/settings_contract.go:781`), content/worldstore/lorebook ownership model, image engine process lifecycle, llama.cpp integration, errors/diagnostics, content-import journal design, sampler preset values.

**Explicitly not a rewrite.** The majority of feature code stays.

---

## 2. Decision — five roots

Installation location must never be storage authority merely because the executable lives there.

| Root | Purpose | Contents | Writable by |
|---|---|---|---|
| **R1 Installation** | executable + bundled runtime | `howl-chat.exe`, `runtime/llama/*`, `sd-server.exe` | installer only; app treats installation artifacts as read-only |
| **R2 Application-data** | Howl-owned state | settings, presets, registries, journals, logs, caches, backups, launch records, navigation | app, no consent needed |
| **R3 User-content** | discoverable user work | chats, characters, worlds, scenarios, lorebooks, images, videos | app, no consent needed |
| **R4 External-resource** | user-selected resources | model directories, registered libraries | **only under a scoped grant** |
| **R5 Temporary** | scratch | audio/image temp, non-content staging | app, ephemeral |

**R1 is never a storage authority.** `SC`: today `InternalConfigDir()` = `StateRoot()/config` and `StateRoot()` derives from `ProjectRoot()`, so configuration lives in the installation directory.

**R1 is a physical classification, not a filesystem ACL guarantee.** An adopted user-content directory may physically overlap R1 (see §11.2) and remains writable. This is not a general exemption: new application state still goes only to R2/R3.

**R3 already exists and is under-adopted.** `paths.DefaultContentDir` (`internal/paths/content.go:25–64`) resolves `Documents/Howl/<Category>` through the Windows known-folder API and maps to `<state>/content/...` under isolation `SC`. Live call sites: `app-libraries.go:56` (creates and registers a Howl-owned save root after an explicit user action) and `:197` (reports the path only). It does **not** govern chats, characters, worlds or lorebooks `SC`.

---

## 3. Modes

| Mode | R2 | R3 | R5 | Selection |
|---|---|---|---|---|
| **Installed** (default) | `%APPDATA%\Howl` | known folders (`Documents`/`Pictures`/`Videos` + `Howl`) | OS temp | no marker present |
| **Isolated** | `<iso>\data` | `<iso>\data\content` | `<iso>\temp` | existing `HOWL_WEBVIEW_DATA_DIR` + marker `SC` |
| Portable | — | — | — | **deferred (D2)** |

- **Portable mode does not exist.** Exhaustive environment enumeration found `HOWL_WEBVIEW_DATA_DIR` as the only storage-affecting switch; every other `HOWL_*` variable is a resource, credential, or test-harness override `SC`. No flag, marker, config option or equivalent behaviour exists. `docs/review/INVENTORY-F3C-ROOT-OWNERSHIP-2026-11-11.md:265` and `docs/review/INVESTIGATION-F3C-DROOT-CLOSURE-2026-11-11.md:279` record a packaged `ProjectRoot` as unknown and state that no portable launch was ever run.
- `DefaultContentDir` has **no portable branch** — a portable install would write content to Documents while configuration sat beside the executable.
- Portable, if ever added, is a **third mode with its own acceptance matrix**, labelled new and unverified until a portable launch is recorded.

---

## 4. Services and primitives

Six services, two shared primitives. No seventh service has been justified; the authority record and navigation state are stores owned by ConfigAuthority.

| Component | Owns | Must not own | Depends on |
|---|---|---|---|
| **PathService** | root set, mode resolution, category→location, class labelling, containment checks, diagnostics | **write permission, grants, policy** | — |
| **ExternalDataGuard** | operation-scoped grants for R4, consent records, app-created-file manifest | directory-wide permission; deciding where a path belongs | PathService |
| **ConfigAuthority** | app-preferred settings, schema, defaults, reset semantics, save results, **authority record**, **navigation state** | model/preset config, machine config, restarting anything | PathService, DurableStore |
| **ModelConfigService** | presets (stable IDs), model↔preset associations, per-model custom inputs, compatibility gating | application defaults, machine paths, launch values | PathService, DurableStore, identity |
| **MachineConfigService** | machine-scoped desired configuration (engine paths, ports), validation, clear semantics | **engine lifecycle, restart, process control** | PathService, DurableStore |
| **RuntimeResolver** | resolution pipeline, transformation records, launch records | persistence of intent; consent | PathService, ConfigAuthority, ModelConfigService, MachineConfigService |
| **EngineLifecycle** (coordinator, not a service) | engine start/stop/restart and outcome reporting | persistence of configuration | MachineConfigService, ExternalDataGuard |
| **DurableStore** (primitive) | document atomicity, backup, checksum, version envelope, single-writer discipline | **multi-step workflow recovery; domain semantics** | — |
| **ContainmentChecker** (primitive) | reparse inspection, canonical identity pinning | authorization decisions | — |

**Retained content stores:** `worldstore`, `contentimport`, `lorebookstore`, `library` keep their existing ownership model (roles `RoleRead/Write/Save`, origins `OriginHowl/External/Imported`) and their own journals.

**Dependency direction (no cycles):**

```
PathService
   ├── DurableStore ── ConfigAuthority
   │                ├─ ModelConfigService
   │                └─ MachineConfigService ── EngineLifecycle
   ├── ExternalDataGuard ── ContentImport / ImageInstaller
   └── RuntimeResolver  (reads config; writes records only)
```

**Three broken overlaps, explicitly resolved:**

1. **Resolution ≠ permission.** PathService answers *where*; ExternalDataGuard answers *may this operation write here*; the store answers *did it succeed*.
2. **Configuration ≠ lifecycle.** `MachineConfigService.Save()` persists. `EngineLifecycle.Restart()` is a separate, explicitly requested call. No stored-preference change may restart a subsystem as a side effect.
3. **Document persistence ≠ workflow journaling.** `DurableStore` covers single-document atomicity; multi-step operations keep per-operation journals with plan hash, resume, abort and rollback scope, as `contentimport` already does (`internal/worldstore/import.go:18–20,63–64,535–536` `SC`).

---

## 5. Contracts

### A — Lifecycle and ownership
- **A1** One application identity: `(stateRoot, appVersion, schemaVersion)` resolved once at startup by PathService, printed in diagnostics.
- **A2** The working directory never determines a storage location. Violations today: `app-characters.go:75,104,185`, `memory/candidates.go:31`, `memory/continuity.go:60` `SC`.
- **A3** Every store declares exactly one class: app-state, user-content, machine-config, external-metadata, external-content, transient.
- **A4** Startup: resolve paths → validate → open stores → load config → restore desired state → (no apply). Shutdown: stop children → flush queues → close stores → clean-shutdown marker.
- **A5** Interrupted write recovery is deterministic (see §8). An unclean-shutdown marker triggers an integrity check.
- **A6** Missing/unavailable path → typed failure naming the category and remedy. No silent relocation.

### B — Path and storage authority
- **B1** All app-owned writes go through PathService. No package joins a relative path for storage.
- **B2** Every resolved path must be contained in its category root; isolated mode additionally requires containment in the isolation root.
- **B3** External directories are read-only unless ExternalDataGuard holds a current grant for the exact operation.
- **B4** External metadata lives in app-managed storage by default, keyed by `(path, size, mtime, header metadata)`. Writing beside assets is per-directory opt-in.
- **B5** Every root is introspectable: category, concrete path, class, isolation status.

### C — Configuration authority
- **C1** One authority per domain: app prefs (ConfigAuthority), model/preset (ModelConfigService), machine (MachineConfigService), runtime (RuntimeResolver, non-persistent).
- **C2** Four distinguished empty states: `absent`, `never-initialized`, `intentionally-default`, `intentionally-empty`. A reset writes `intentionally-default` + `schemaVersion`; migration accepts only `absent|never-initialized`.
- **C3** Preset principles: stable IDs; associations separate from application; apply ≠ associate; explicit per-section application; rename-safe; delete requires explicit resolution; kind-gated compatibility; atomic section application; intent-only storage.
- **C4** Presets attach to a confirmed identity (§5.1). Unmatched configuration is preserved and marked, never reassigned.
- **C5** A changed family or replaced model marks dependent sections `incompatible` with reasons; no automatic rewrite.
- **C6** Custom prompt/Jinja inputs are explicit save targets with an association, not per-keystroke writes. Jinja distinguishes `no override` from `explicitly empty` `SC`.

### 5.1 Model identity — layered (normative)

Accepted 2026-11-22 during the challenge pass and recorded here because **no earlier repository document defines it** (see the note at the end of this subsection). Automated full-file hashing of models is not a default.

| Level | Evidence | Cost | Establishes |
|---|---|---|---|
| **L0** registry record | name, path, size, mtime, family, admission | free (`model-registry.json`) | candidate match only |
| **L1** header metadata | GGUF arch, layer count, embedding length, quant, declared context | one small read | strong candidate; **not** identity |
| **L2** partial digest | size + first/last 1 MiB hashed | bounded | **confirmed binding to the observed instance**, once the user confirms |
| **L3** full digest | whole-file hash | expensive, user-visible | content identity |

Rules:

- Presets bind to **L2 or better plus explicit user confirmation**. That establishes a re-checkable association with the file **as observed** — **not** cryptographic identity, and it is never reported as hash-verified.
- L0/L1-only matches are `unresolved` and must be confirmed before binding.
- The fingerprint is **recomputed when the association is resolved**. A change produces `needs-review`; the preset is preserved and shown unmatched. **Never auto-rebind**, even when L0/L1 also match.
- L3 is offered, never imposed, with its cost shown.
- Discovery and inventory use L0/L1 only. No level of hashing runs as part of ordinary directory scanning, and no hashing runs over model or content trees automatically.

**Do not mistake the existing `modelintel` fingerprint for L2.** It digests a small sidecar JSON, not model weights (`internal/modelintel/modelintel.go:82,96,121`: *"fingerprint changes only when a small sidecar file changes"*) `SC`. It may serve L1; it is never identity.

**Provenance note.** A repository search across all `docs/` and Go sources for a layered model-identity definition returned no such document; the only identity vocabulary in the tree concerns filesystem identity (volume serial + file ID) for storage roots, e.g. `docs/review/OWNER-DECISION-REQUEST-F3C-FOUNDATIONS-2026-11-11.md:46` and `docs/review/INVESTIGATION-F3C-DROOT-CLOSURE-2026-11-11.md:268`. This subsection is therefore the first written record of the decision; it is not a cross-reference to an existing one.

### D — Persistence lifecycle
- **D1** Five observable states as separate columns: persisted, restored, applied, effective, reported.
- **D2** Three effectiveness levels: `consumer-established` (`SC`), `runtime-confirmed` (`RC`), `unestablished` (`UN`).
- **D3** Save returns an explicit result. The UI may not display "saved" without `Saved`.
- **D4** Reset is deterministic across restart: no legacy file may be consulted for a reset domain, enforced by tombstone rather than by absence.
- **D5** Migration is non-destructive: pre-backup, identity-checked source, plan hash, resumable journal, verification before retirement.
- **D6** Every document carries `schemaVersion`. Unknown newer versions are refused, never rewritten.

### E — Runtime resolution
- **E1** Stages recorded separately: intent → defaults → model/family → constraints → resolved → emitted → adopted → reported → observed.
- **E2** Every transform produces `{stage, input, output, rule, reason}`. A clamp is a transform, not a silent mutation.
- **E3** Forbidden: in-place mutation of config (`applySafeDefaults`), emission without record, and a value that never reaches its consumer.
- **E4** Unsupported values are rejected with reasons, never coerced.
- **E5** Beginner automation is retained — family detection, template selection, safe defaults — but each becomes a recorded transform.

### F — Isolation and external data
- **F1** Configuration isolation: every config/registry/journal path inside the isolation root. Established `SC`.
- **F2** Application-data isolation: every app-state and user-content write within the designated roots for the active mode (mode-relative, not absolute).
- **F3** External isolation: no Class C write in any mode (§10).
- **F4** Read-only browsing: pickers and model listing perform no writes. Violated today by `GetCharacters` (`app-characters.go:192`, `MkdirAll` on a read path) and by the image scanner on page load `SC`.
- **F5** Deletion is ownership-verified against a manifest, never by glob.

---

## 6. Configuration authority record (routing)

Routing must be decided by one document, written once.

### 6.1 Record

`R2/authority.json`, written only through DurableStore:

```json
{ "schemaVersion": 1, "recordVersion": 3,
  "domains": {
    "appSettings":  { "store": "settings.v2.json",    "authority": "live",    "generation": 4 },
    "modelPresets": { "store": "presets.v1.json",     "authority": "live",    "generation": 1 },
    "legacyModel":  { "store": "model-settings.json", "authority": "retired", "generation": 0 }
  } }
```

- The record is **total** over known domains: every domain appears with `authority ∈ {live, retired}`.
- **Routing comes from this record alone.** No runtime reader inspects a store to decide whether to use it.
- A store's own `state` field is a **write guard only** (`staging` refuses writes from anything except the migration importer, and refuses to be routed). It is never a router. Two fields, two jobs, no dual source of truth.
- **`generation` is per domain.** A store is stale iff `store.writtenUnderGeneration < domains[d].generation`. Migrating one domain never marks another stale.

### 6.2 Initialization states

| Marker | Record | Other evidence | State |
|---|---|---|---|
| absent | absent | no witness, no v2 store, no legacy store | **NeverInitialized** |
| absent | absent | any witness / v2 store / legacy store | **RecoveryRequired** |
| present | absent | — | **RecoveryRequired** |
| absent | present | — | **RecoveryRequired** |
| present | valid | — | **Active** |
| present | corrupt, `.bak` valid | — | **Active** (previous generation) |
| present | corrupt, backup corrupt | — | **RecoveryRequired** |

**Rule: two missing documents never prove a first launch.** `NeverInitialized` requires positive evidence of freshness.

**Bootstrap record:** `R2/authority.boot` = `{initializedAt, firstGeneration, appVersion, schemaVersion}`. Contains no routing, so it cannot compete with the record. Write-ahead order: bootstrap → record (generation 1, legacy routing). A crash between them yields `RecoveryRequired`, the conservative direction.

**Independent witness:** `%LOCALAPPDATA%\Howl\state-identity.json` = `{stateId, installId, initializedAt, schemaVersion}` — deliberately in a different root from R2, so a loss confined to R2 cannot erase it. Local (not roamed), so profile sync cannot fabricate or resurrect it. Second line of evidence: any v2 live store with `generation ≥ 1`.

**Repair ordering (normative):** fresh-start or repair establishes the replacement first — write bootstrap, write record generation 1, verify both readable, copy the previous pair aside, and only then rotate the witness. The witness is never the first thing removed.

### 6.3 Crash matrix

| Crash point | Result | Recovery |
|---|---|---|
| Before the record write | all domains legacy → legacy reads. Consistent | re-run migration |
| During the record write | atomic replace → old or new record. Consistent | none |
| After the record write | target domains live → new reads. Consistent | none |
| Named live store missing/corrupt | that domain **refused**, reported | revert generation or re-import. **Never silent fallback to legacy** |
| Record corrupt, `.bak` valid | route per backup | report recovery |
| Both corrupt | safe mode: **no domain routed** | explicit repair only |

Safe mode offers revert-generation, re-import with confirmation, or start-fresh with confirmation, and performs none automatically.

---

## 7. Containment — fail-closed

Verify-then-write is not a security guarantee: a write that already landed outside the grant is not undone by detecting it afterwards.

> **External writes that require containment must use handle-based operations that preserve the validated directory identity through creation and writing. If the platform or toolchain cannot provide that guarantee, the operation fails closed.**

1. At grant issuance: open `destinationRoot`, pin canonical path + volume serial + file ID.
2. Create `.howl-import-staging` and `<operationId>` **one level at a time relative to the validated parent handle**, re-validating the parent identity at each step. `MkdirAll` on a full string is prohibited.
3. Retain the staging-root handle for the operation; all item operations are relative to it.
4. Reject a reparse point **created after issuance** anywhere in the granted subtree. Reparse points present at issuance are legitimate and are represented in the pinned identity — blanket rejection would break real installations (`G:\AI GGUF`, `H:\InvokeAI`, `E:\HOWL_Chat` are junction- or mount-backed) `SC`.
5. Comparison is by `(volume serial, file ID)`, never by string equality — settling case-insensitivity, short names, `..` traversal and trailing character forms in one step.
6. **Capability gate:** handle-relative creation requires `NtCreateFile` with `RootDirectory` set to the directory handle (`golang.org/x/sys/windows`; not in the stdlib). This capability must be demonstrated before external import is enabled; otherwise import fails closed with an explicit "not available on this build" state. No degradation to path-based writes.

**Existing prior work:** `hasReparseComponent` (`internal/audio/tts/helper/pathcheck_windows.go:12–45`, duplicated in `fakehelper`) walks every component to the volume root and rejects symlinks and `FILE_ATTRIBUTE_REPARSE_POINT`. It is promoted into **one** shared pre-flight primitive for inspection and rejection at grant time — not TOCTOU protection, and not the complete policy for external roots.

**Weaker sites, recorded not blessed:** `gallery.go:243–250` compares `EvalSymlinks` results and *accepts* an unresolved non-existent target; `contentimport/prepare.go:414–422` ignores resolution errors.

---

## 8. DurableStore — guarantees

Defined as guarantees first; implementation shape follows.

| Guarantee | Acceptance criterion |
|---|---|
| Atomic visibility | a concurrent reader always parses a complete document; no partial read in 10⁴ interleaved iterations |
| Crash durability (process) | kill between flush and rename → previous document intact; kill during temp write → previous intact; at most one stale temp, swept on next start |
| **Power-loss durability** | **`UN`.** Depends on NTFS/volume flush behaviour and directory metadata durability, for which Windows offers no portable handle. Must be verified empirically before being claimed |
| Pre-image retention | **N = 2 (D3, provisional)** per document, `doc.v1.bak`/`doc.v2.bak` with checksum, timestamp, schemaVersion. Discarded only when a third would exceed N, or on explicit user clear. The most recent pre-image is never discarded by a failed write |
| Backup interruption | kill during pre-image creation → previous intact, no `live` backup record pointing at a truncated copy |
| Integrity | checksum mismatch → `Recovered(from backup)` or `Refused`; never overwrite |
| Unknown newer schema | preserved byte-for-byte as `RefusedUnknownSchema`; copied aside before any downgrade attempt; no repair, no rewrite, no field stripping |
| Single writer | **every** writer takes the lock: save, reset, migration importer, recovery, rotation. No exemptions |
| Outcome | `Saved | NotSaved(reason) | Recovered(from) | RefusedUnknownSchema`. No booleans |

**Rotation ordering (invariant):** `acquire lock → write temp → flush → pre-image copy → rename → drop oldest beyond N (only after rename succeeded) → release`. A failed commit never triggers rotation, so rotation can never cause loss of the last good document.

**Locking:** `LockFileEx` on `<doc>.lock`, exclusive, byte range `[0,1)`, via `golang.org/x/sys/windows` (existing dependency). Windows releases the lock when the owning process terminates, so an abandoned lock cannot wedge the app; a leftover lock file is harmless and is not "stale state" to delete. **No cross-process locking exists anywhere in the codebase today** — the only lock is `runMu.TryLock()` (`internal/image/engine/sdcpp/video.go:154`), an in-process mutex `SC`.

**Two independent timers, never one number:** **lock acquisition** (default 5000 ms) and **rename retry** (3 × 50 ms) *inside* an acquired lock. Total operation bound = lock wait + write + flush + pre-image + retries. "150 ms" is not a bound on any operation.

**Contention → `Busy`. No bypass** (D4).

### 8.1 Reader-open contract

Replacement requires delete access to the target; a reader that does not share delete blocks it.

> **The current settings reader has no explicitly established delete-sharing contract. Its read is short-lived** — `os.ReadFile` at `internal/settings/store.go:23` opens, reads and closes. That does not establish that overlap cannot occur, and does not establish delete sharing. **Compatible reader sharing is an acceptance requirement of the replacement protocol, for every reader that can overlap replacement.**

`os.WriteFile` at `:51` is a **direct, non-atomic write** — no temp, no flush, no rename, no backup `SC`. DurableStore must replace this; it must not be assumed to exist.

| Reader protocol | Requirement | Suitable for |
|---|---|---|
| Short-lived read + bounded writer retry | reader opens, reads, closes within a bound; writer treats sharing violation as retryable **inside** the lock window, then `Busy` | small config documents, current call sites |
| Delete-sharing reader | open with `FILE_SHARE_DELETE` (Go opt-in per golang/go#34681, or `NtCreateFile` with all three share flags); exact flag verified against the toolchain in use (`UN`) | any reader that may hold a document across another operation |

---

## 9. Runtime evidence

Nine stages. **No stage may be inferred from the one before it.**

Labels: **M** source establishes the mechanism exists · **X** a specific execution produced a recorded value · **A** the backend acknowledged or reported it · **O** behaviour was observed to change.

| Stage | Established by | Required evidence | Current |
|---|---|---|---|
| 1 Intent | ConfigAuthority / ModelConfigService | stored value + provenance | **M** |
| 2 Defaults | schema | default-table entry | **M** |
| 3 Model/family | ModelConfigService + profile | profile identity, section association | **M** (partial) |
| 4 Constraints | runtime environment probe | captured hardware snapshot | **M** only — no snapshot recorded here |
| 5 Resolved | RuntimeResolver | computation record per key | absent |
| 6 Emitted | launcher | argv or HTTP-body capture | **M + X** (`run1.out:62`, argv `-c 3096` `AO`) |
| 7 Accepted | backend | ack, health, documented acceptance signal | `UN` |
| 8 Reported active | backend | `/props`, status endpoint, server log | `UN` |
| 9 Behaviourally observed | measurement harness | e.g. output truncation at a stop token | `UN` |

**Emission is not adoption.** An argv proves stage 6 only. Stages 7–9 may be permanently unavailable for some settings; the record carries `strongest_evidence` rather than a boolean, and `unestablished` is a legitimate terminal value.

**Launch record is designed independently of FC-5** (`howl.launch.v1`). Mapping compatible fields into a future approved contract is a later, separate act; no frozen document is depended upon.

---

## 10. External data and isolation

### 10.1 Three write classes

| Class | Definition | Authorization |
|---|---|---|
| **A — app-managed** | writes inside the active mode's R2/R3/R5 | none needed |
| **B — authorized external** | writes into R4 at a user-requested destination | **scoped grant** |
| **C — implicit external** | any R4 write not requested per-operation | **forbidden** |

**The guarantee is zero Class C writes**, not zero writes outside R2.

### 10.2 Grants

`Grant{ id, operation, destinationRoot, stagingNamespace=".howl-import-staging", stagingRoot=<destinationRoot>/.howl-import-staging/<operationId>, maxItems, initiatedBy, createdAt, expiresAt, journalRef }`

1. The namespace is reserved and literal; the grant covers exactly one `operationId` subtree.
2. **Generated paths are authorized by containment, not enumeration:** permitted iff strictly inside `stagingRoot`. Equal to `stagingRoot`, inside `.howl-import-staging` itself, or inside a *different* `operationId` → refused.
3. **The grant authorises creating `stagingRoot` itself** (a directory) and creating/writing/deleting descendants of that one subtree. It does not authorise writing to `stagingRoot` as a file, nor anything else under `destinationRoot`.
4. **Commit needs a second authorization:** exact final paths, enumerated from the plan and journaled before the first write. Generated staging paths are legitimate; generated *final* paths are not.
5. Deletion is confined to `stagingRoot` and the operation's own `operationId` — which is why owner markers exist (`import.go:535–536` `SC`); the grant makes their use mandatory.
6. Terminal states are enforced: on success or abort, `stagingRoot` must be empty or removed; a non-empty root at next start is reported, never silently adopted or deleted.
7. Grants expire; an expired grant refuses all writes and records the refusal.

### 10.3 Worldstore staging, reconciled

World imports stage into `<library>/.howl-import-staging/<operationId>/<itemId>` — inside R4 (`import.go:18–20,63–64` `SC`). That is **legitimate Class B**: it exists only because a user initiated an import, it is operation-scoped, and it already carries owner markers. The isolation contract therefore permits external staging under an import grant, requires it to be inside the granted library root, requires cleanup on abort, and forbids reuse as general storage. **No change to the ownership model is required.**

### 10.4 Mode-relative guarantees

> **F1 (configuration):** every app-state configuration, registry and journal write lands within R2 for the active mode.
> **F2 (application data):** every app-state and user-content write lands within R2 or R3 for the active mode.
> **F3 (external):** no Class C write in any mode; Class B only inside a scoped grant.

Isolated mode is not a different guarantee — it is the case where R2/R3/R5 **collapse into the isolation root**, which is why it is the cheapest mode to test.

### 10.5 Journal and staging classification

| Artifact | Class | Location |
|---|---|---|
| `contentimport` journal, plan hash, resume tokens | **app-state** | config dir under R2 `SC` |
| Worldstore staging tree + owner markers | **operation-scoped external staging** | granted library under R4 `SC` |
| Import audit/commit results | app-state | R2 |
| Launch records, scanner manifests, index | app-state | R2 |

The two cases must never share a rule.

---

## 11. Content location

Full policy: `ADR-002`.

### 11.1 Chats — two populations
New installations → `Documents\Howl\Chats`. Existing installations → existing root, adopted in place. Later changes → explicit relocation only. Selection is an **adoption decision** based on installation identity and prior initialization, never on chat-file presence (ADR-002 N18).

### 11.2 Characters — adopted legacy root
`SC`: `CreateCharacter` (`app-characters.go:74–79`) prefers `<exeDir>\Characters` when it exists, while `getCharactersDir` (`:103–113`) and `GetCharacters` (`:181–192`) prefer cwd and fall back to the exe dir only when cwd stat fails — **create and read can resolve different directories**. `GetCharacters` also `MkdirAll`s on a read path. Real data exists at the install-relative location (`characters/Aria`, `characters/Doren` in this workspace).

Decision: preserve in place; do not move, merge into the library system, or overwrite. Record an **adoption record** `{category, resolvedPath, volumeSerial, fileID, adoptedAt, legacy:true}`; resolve thereafter from the record; one resolver with one priority for create and read; reads stop creating. A changed identity triggers `needs-adoption-review` and never auto-adopts whatever now occupies the old path. Permanent root is `UN` pending installed-layout evidence.

> Character data is user-content; its existing physical location is preserved during the rebuild, but installation location and working directory no longer implicitly determine its future storage destination.

### 11.3 Presets — three kinds
Managed presets and associations (R2, no picker) · imported preset files (external source, reference only) · exported preset documents (user-content, `Documents\Howl\Presets` default). Never flattened into one model-like external category.

---

## 12. Migration sequencing

> **Exactly one authoritative writer and one authoritative runtime reader per configuration domain, at every stage.**

| Stage | Authority (`live`) | New store | Legacy | Runtime reader |
|---|---|---|---|---|
| S0 Inventory | legacy | absent | read-only inventory | legacy |
| S1 Import | **legacy** | `staging` (migration tool only) | read-only source | legacy |
| S2 Verify | **legacy** | `staging`, compared | read-only source | legacy |
| S3 Cutover | flip → **new** | `live` | **frozen** | **new** |
| S4 Retire | new | `live` | `retired`, read-only recovery source | new |

- A runtime reader **refuses** a `staging` store, so a partially imported domain cannot be read by accident.
- Cutover is one metadata flip per domain, after verification, and is the only moment the reader changes.
- Legacy is never consulted by the new runtime after cutover; recovery is an explicit operation, not a fallback.
- **Inventory ≠ backup ≠ hashing:** inventory is read-only enumeration plus metadata; backup is an explicit write into R2/backups for the domains in the current phase; hashing is identity-grade and never runs automatically over model or content trees.

---

## 13. Decisions

| # | Decision | Consequence |
|---|---|---|
| **D1** | Chats at `%APPDATA%\HOWLChat` for **existing** installations, adopted in place; `Documents\Howl\Chats` for **new** installations (refined by ADR-002 §4) | no user-visible move; chats remain user-content in both cases |
| **D2** | Portable mode deferred; Installed + Isolated are v1 | explicit scope decision, not an absence; a third mode needs its own matrix and a recorded portable launch |
| **D3** | N = 2 pre-images, provisional, configuration documents only | not applied to content documents until size is measured (`UN`) |
| **D4** | Lock contention → `Busy`; **no bypass** | "overwrite anyway" is rejected as incompatible with single-writer |
| **D5** | Legacy stores retained read-only indefinitely | retirement is separately authorized; migration implies no permission to delete |

---

## 14. Consequences, risks, open items

**Becomes impossible:** legacy resurrection after reset · implicit sidecar writes on page load · cwd-derived storage · dual authority (one record, one writer and reader per domain) · lock bypass · silent fallback from a missing post-cutover store to legacy.

**Becomes required before shipping:** demonstrated handle-relative creation capability (else import is disabled) · delete-sharing or bounded short-lived readers · verified lock wrapper · crash-injection tests for interrupted replacement, interrupted pre-image and simulated sharing violation.

**Residual risks:** import availability depends on a capability gate — a real product risk, to be demonstrated early rather than assumed · power-loss durability may never be claimable · portable installs receive installed-mode behaviour silently (a startup diagnostic is the available mitigation).

**Open (`UN`):** portable mode · chat relocation for existing installs · character permanent root (needs installed-layout evidence) · largest chat/content document size for D3 · delete-sharing flag and `NtCreateFile` capability verification · `LockFileEx` wrapper and lifecycle · WER/Wails external writes · `contentimport` journal filename · `chat.Service` reachability (C-03) · `prompt` package disposition · persistence path for character/world image data URLs and image-generator uploads.

---

## 15. Authorization

This record fixes architecture direction and decisions only. It authorizes no implementation, no code, schema, API, migration, test, build, launch or process change, and no amendment to any frozen contract. FC-5 revision 3 remains frozen, FC5-31 remains `Unresolved`, R6 remains unapplied.

Implementation of any part requires separate, explicit authorization.
