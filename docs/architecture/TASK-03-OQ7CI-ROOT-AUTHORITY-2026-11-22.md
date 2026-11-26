# TASK-03 — OQ-7c-i: root authority

| Field | Value |
|---|---|
| Purpose | Which root is authoritative, and whether two instances can resolve one conversation to different physical paths |
| Status | **Read-only investigation. Design only. V3.1 unchanged. Nothing authorized.** |
| Method | Source inspection of every root derivation, the isolation switch, and launch-time configuration |
| Recovery token | `TASK-03-OQ7CI-ROOT-AUTHORITY-2026-11-22` |

Continues OQ-7's numbering (`O7-116`…).

---

## 1. Every root derivation

**O7-116 — Three derivations, two cached, one recomputed:**

| Root | Derivation | Cached? | Scope |
|---|---|---|---|
| `paths.ProjectRoot()` | search **upward from the executable's directory** for `wails.json` or `go.mod`, max depth 6; else CWD; else exe dir | yes, `cachedRoot` | installation |
| `getAppDataDir()` | `os.UserHomeDir()/AppData/Roaming/HOWLChat`, or `StatePath("persisted-chats")` when isolated | **no**, per `NewPersistenceManager` | user profile |
| `paths.StateRoot()` | isolation `StateRoot` when active, else `ProjectRoot()` | via isolation `sync.Once` | run |

**O7-117 — `ProjectRoot()` is launch-dependent, which is a root-authority defect in its own right.** Two behaviours follow from the code:

| Situation | Resolved root |
|---|---|
| shipped build under a tree containing `go.mod`/`wails.json` within 6 levels | **that enclosing project** — the app adopts a repository as its data root |
| no marker found anywhere up the chain | **the current working directory** — launching from a different directory silently relocates every installation-tree file |

The practical consequence: launching the same binary from two different working directories, or placing it in two different folders, yields **two different installation roots** — so `Chats/`, `memories.json`, configs and caches are split, while history is not.

**O7-118 — Isolation is resolved once and cannot change while running.** `ResolveIsolation` uses `sync.Once`; `main()` validates before `NewApp` and exits with code 2 if the root is unusable. An isolated run uses `<root>/data` as its state root and `<root>/data/persisted-chats` for the legacy chat store, and the isolation root is required to be outside the project tree with a marker file.

**O7-119 — No configuration can relocate a root at runtime.** `AppConfig` carries exactly one field, `LastChatID`. There is no user-facing setting for a data directory, no environment override for the chat root, and no reload path for settings that would move it.

---

## 2. Can two instances resolve one conversation to different physical roots?

**O7-120 — For history, only in three cases, and two of them are intentional:**

| Case | Same physical root? |
|---|---|
| two production instances, same user, same installation | **yes** — `%APPDATA%\HOWLChat` is recomputed identically |
| two different Windows users | no — different profiles; they are different conversations in every sense |
| production vs isolated run | no — by design, and validated to be disjoint |
| roaming-redirected profile on two machines | **no**, even though the *logical* path is identical |

**O7-121 — The roaming case is the one that matters.** The history tree lives under `AppData\Roaming`. Roaming means the folder may be synchronised between machines, so two machines can hold the same conversation at the same logical path with **no shared filesystem**, hence no shared lock, and no mutual exclusion of any kind. A sibling lock namespace under a roaming root therefore proves exclusivity **within one machine's view of the profile** and nothing more.

**O7-122 — For installation-tree data, divergence is ordinary, not exotic.** Per O7-117, two launches can resolve different installation roots. Every non-history file inventoried in row 566 lives there, so the *shadow* tree is the least stable part of the system.

---

## 3. Verdict

**O7-123 — The authoritative conversation root is the history root**, `%APPDATA%\HOWLChat` (`chats/`), with this precise scope:

| Property | Holds? |
|---|---|
| identical across production instances of one user on one machine | yes |
| stable within a process | yes |
| stable across launches | yes — no launch-location dependence |
| shared with isolated runs | no, deliberately |
| shared across machines via roaming | **no** |

**O7-124 — No lock may be derived from `ProjectRoot()`.** It is launch-location dependent (O7-117), so a lock namespace under it would not identify one conversation across two launches — the failure the whole branch is trying to prevent.

**O7-125 — Therefore: lock identity derives from the history root only, and the guarantee is machine-scoped per user.** Any enforcement claim must say "two instances on this machine, this user, this root" — not "two instances". A claim that spans machines would require a coordination mechanism the product does not have and this branch does not propose.

**O7-126 — Isolation is the one case where a different root is correct**, and it is correctly validated: the marker file prevents adopting a directory holding real data, and the root is required to be outside the project tree. An isolated run is not a second writer to production.

---

## 4. Consequences for the branch

**O7-127 — OQ-7b's sibling namespace stays valid, with a narrowed claim** (O7-125): per-conversation locking can be trusted for same-user, same-machine production instances. For roaming profiles it is trustworthy as a local serialiser only, and the contract must not imply more.

**O7-128 — The install-tree instability is not this contract's problem, but it must not be mistaken for a lock problem.** Row 566's shadow tree already sits outside reconciliation's scope (OQ-7d-i); O7-117 now shows it is also outside any stable identity, which reinforces the exclusion rather than weakening it.

**O7-129 — Nothing here is authorized.** No root consolidation, no relocation setting, no roaming mitigation, no lock namespace, no isolation change, no pointer change, no gate, CAS, marker, migration or deletion change.

---

## 5. Tests this would require

- **O7-T31** Two production instances of one user resolve the identical history root and contend for the same lock.
- **O7-T32** An isolated run and a production run are proven disjoint, and neither can see the other's conversations.
- **O7-T33** The same binary launched from two working directories resolves the same **history** root — found by test, since it currently resolves the same root for history by construction.
- **O7-T34** The same binary launched from two working directories is shown to resolve **different** installation roots, documenting O7-117 as a known limitation rather than a defect to fix here.
- **O7-T35** A lock path derived from `ProjectRoot()` is rejected by the boundary API.
- **O7-T36** The enforcement claim's scope statement names user and machine, and fails review if it says "two instances" unqualified.
