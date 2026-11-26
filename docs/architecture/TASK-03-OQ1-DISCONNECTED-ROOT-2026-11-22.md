# TASK-03 — OQ-1: disconnected-root classification

| Field | Value |
|---|---|
| Purpose | Decide whether V3.1 can distinguish an absent conversation from a location the scanner could not observe, and what row 1 may report |
| Status | **Read-only investigation. Outcome: Amend. V3.1 unedited. Nothing adopted, nothing authorized.** |
| Method | Source inspection of root resolution, the chat store's error predicates and the Windows path APIs; one disposable probe of this toolchain's error classification (`build/tmp/oq1probe`, log `build/tmp/oq1probe.log`) |
| Recovery token | `TASK-03-OQ1-DISCONNECTED-ROOT-2026-11-22` |

First of the six baseline-changing questions named by row 576 (`OQR-06`…`OQR-11`), resolved under the owner's three-outcome frame: **Amend**, **Clarify**, **Retain**. This record answers OQ-1 only; OQ-2, OQ-3b-b, the row 550 collision, the `V3-44` wording and the memories limitation remain open and are not touched.

---

## 1. A correction carried, and three guardrails not reopened

**UNC-01 — Correction to row 576's prerequisite wording.** `OQR-13`, `OQR-14` and `OQR-26` say that no prerequisite can change a row, and that sentence is imprecise. The correct statement is: **a prerequisite cannot change a decision row's defined rule, but it can change capability availability and therefore which rows are reachable.** Rows 16–17 are unreachable while `gate_enforced` is false (`V3-36`, `V3-37`); OQ-6 would make `V3-44` case 3 an observable input; an enforcement record would make rows 16–17 reachable. None of them changes what a reached row decides. Row 576's summary is superseded to that extent, and the distinction is carried into this record's own outcome language. Rows 575 and 576 are otherwise unchanged; the correction is recorded here rather than by rewriting an accepted investigation record.

**UNC-02 — Three guardrails from row 575's review are carried, not reopened.** `applied`, `stored` and `durable` remain three different claims: a successful write call is not proof of crash durability, and an in-memory capture is not a stored message (`O7-248`, `O7-249`). A failed shutdown may not convert displayed work into presumed saved work (`O7-266`, `O7-268`). A discard site is a conformance failure regardless of how correct the API is (`O7-272`). None is a reconciliation input and none is re-decided here.

**UNC-03 — One question, one outcome.** OQ-1 asks what a verdict may say when the scanner cannot observe the location it is judging. It does not decide OQ-2's physical identity, OQ-3b-b's terminator rule, repair, retry or messaging.

---

## 2. What V3.1 says, and the assumption it does not state

**UNC-04 — The rule as written.** `V3-39` row 1 is: *conversation directory absent* → **Missing**, reason `absent`. Row 2 is: *enumeration incomplete, path denied, or I/O error* → **Unavailable**, reason `unreadable`. `V3-60` adds that incomplete enumeration is an error, never a partial inference.

**UNC-05 — The unstated assumption: "conversation directory absent" is treated as a decidable observation.** Row 1 has no precondition about the root that contains the conversation. It does not say who observed the root, whether that observation succeeded, or what class of volume the root is on. The row reads as if a failed `stat` of the child were independent evidence of deletion.

**UNC-06 — The open question was recorded as a guess, not a rule.** Row 547 §9 lists OQ-1 as "Unestablished; likely always Unavailable"; `V3-73` carries it forward. "Likely" is not a condition a scanner can execute, and `V3-39`'s first-match table has no branch for it. That is precisely the state the governing principle forbids: a condition that exists in prose but not in the table.

---

## 3. Observed: this toolchain classifies an unreachable root as "not exist"

**UNC-07 — Method.** A disposable probe called `os.Stat` and `os.ReadDir` on six paths and printed the error, `os.IsNotExist`, `errors.Is(err, fs.ErrNotExist)`, the `syscall.Errno`, `filepath.VolumeName` and `GetDriveTypeW`. It created nothing and wrote nothing outside its own console output. Toolchain: `go1.27.3 windows/amd64`. No product code was executed or changed.

**UNC-08 — Results, verbatim from the log:**

| Path | errno | Message | `os.IsNotExist` | `errors.Is(ErrNotExist)` |
|---|---|---|---|---|
| `\\nonexistent-host-howloq1\share` | **53** | The network path was not found | **true** | **true** |
| `\\nonexistent-host-howloq1\share\chats\abc` | **53** | The network path was not found | **true** | **true** |
| `\\127.0.0.1\howl_oq1_no_such_share` | **67** | The network name cannot be found | **false** | **false** |
| `\\127.0.0.1\howl_oq1_no_such_share\chats\abc` | **67** | The network name cannot be found | **false** | **false** |
| `Q:\HOWL-oq1-nonexistent` | 3 | The system cannot find the path specified | true | true |
| `C:\HOWL-oq1-nonexistent-dir` | 2 | The system cannot find the file specified | true | true |

`GetDriveTypeW` returned **3** (`DRIVE_FIXED`) for `C:\` and **1** (`DRIVE_NO_ROOT_DIR`) for `Q:\`. `filepath.VolumeName` returned the full `\\host\share` for a UNC path (Observed).

**UNC-09 — The decisive fact: errno 53 is `ERROR_BAD_NETPATH`, and Go maps it to `fs.ErrNotExist`.** `os.IsNotExist(err)` is **true** for an unresolvable network root. A scanner written with the obvious Go idiom — `if os.IsNotExist(err) { Missing }`, or `errors.Is(err, fs.ErrNotExist)` — reports `Missing` for a conversation it never observed. The error is not "the data is gone"; it is "the network path was not found".

**UNC-10 — And the same predicate is false for a different network failure.** errno 67 (`ERROR_BAD_NET_NAME`) is not mapped to `ErrNotExist`. So two unreachable UNC roots, both unobservable for the same practical reason, produce **opposite verdicts** from the same predicate: one `Missing`, one `Unavailable`. Neither verdict follows from evidence; both follow from an error-mapping detail. That asymmetry is the finding, and it is stronger than the original "absence may mean disconnected" phrasing.

**UNC-11 — The local cases are correctly classified, which is why the sentinel looks trustworthy.** errno 2 and 3 on a fixed local volume are genuine absence and `IsNotExist` is true. Any test written against a local temporary directory will pass while the network case is inverted. This is a test-environment blind spot, not a code style issue.

**UNC-12 — Root class is detectable, but not from the path string alone.** `filepath.VolumeName` identifies a UNC root by shape, which covers `\\server\share` and `\\?\UNC\server\share`. A mapped network drive (`Z:`) has a drive-letter shape; its class comes from `GetDriveTypeW`, where 3 is `DRIVE_FIXED` and 4 is `DRIVE_REMOTE` (the constant is Source-derived; no mapped drive was available to observe). A drive letter with no mounted volume returned 1, observed. Therefore path shape under-detects network roots exactly where enterprises most often use them — mapped drives and redirected profiles.

**UNC-13 — `DRIVE_FIXED` is evidence, not proof.** A fixed volume can be backed by iSCSI, a VHD on a share, or a filter driver that returns transport errors. The class can decide whether absence is *assertable*; it cannot override an error code that says the volume is unreachable. Both inputs are required (§6).

---

## 4. Can this product actually have a network root?

**UNC-14 — Yes, conditionally, and the condition is ordinary.** `getAppDataDir` returns `os.UserHomeDir()` + `AppData\Roaming\HOWLChat` (`persistence.go:62-75`). Windows folder redirection can place `AppData\Roaming` on a UNC path, in which case the history root is a network root with no configuration change by the app. This is **Source-derived** from the resolution code and Windows redirection behaviour; it was **not reproduced** in this session and no user profile was inspected.

**UNC-15 — Isolation can also supply one, but only as an operator setup.** `HOWL_WEBVIEW_DATA_DIR` accepts any absolute path, including a UNC path, and validates it with a marker file and a create-temp probe (`isolation.go:18,67-120`). That is a test/development path, and a root unreachable at startup fails validation. It is not a production configuration.

**UNC-16 — A root unreachable at startup stops the app before any scan.** `initialize` calls `os.MkdirAll(rootDir)` and `os.MkdirAll(chatsDir)` before anything else (`persistence.go:77-89`); on a disconnected share both fail and `NewPersistenceManager` returns an error. So the realistic case is not "the app starts on a dead share" — it is **the share disappearing mid-session**, after the root existed, while a scan is running or about to run.

**UNC-17 — The same sentinel already appears in product code, and its reachability there is not established.** `loadConfig` uses `if _, err := os.Stat(configPath); os.IsNotExist(err)` (`persistence.go:103`), which is true for errno 53. Whether a mid-session disconnect can reach that branch and what a subsequent `saveConfig` would do is **Unestablished**; the startup `MkdirAll` usually fails first. This is recorded as corroboration that the predicate is load-bearing in the package, not as a second defect finding.

**UNC-18 — No reconciliation scanner exists yet.** The rule is being written, not repaired: no `Scan` implementation, no verdict producer and no root observation exist in the repository (observed by absence). OQ-1 is therefore a **contract defect**, and this is the point at which it is cheapest to correct.

---

## 5. Why row 1 is wrong as written

**UNC-19 — A failed child `stat` is not independent evidence when the parent was never observed.** Row 1 evaluates "conversation directory absent" from a per-conversation call. On errno 53 that call returns the same answer for *deleted conversation*, *share offline*, *server renamed*, *DNS failure* and *credential expired*. The row then asserts the most specific and least supported of those.

**UNC-20 — This is the exact inversion of the contract's central rule.** "A scan error must never become `Missing`" is stated for the scanner, and `V3-47` forbids `Unavailable → Missing`. But `V3-47` governs transitions between verdicts; it says nothing about a Missing verdict that was reached in one step from an unobservable root. The prohibition is satisfied vacuously while the harm occurs.

**UNC-21 — `V3-60` does not cover it either, and row 2 is ambiguous.** `V3-60` makes *incomplete enumeration* an error. Row 2's third condition, "I/O error", could be read as covering a failed child `stat` — but then rows 1 and 2 overlap on the same input, contradicting `V3-67`'s mutual-exclusivity claim, and `V3-39`'s first-match rule resolves the overlap in favour of row 1. Either reading ends at Missing; the table's ordering makes the unsafe reading authoritative.

**UNC-22 — Calculated: the harm is a false deletion claim, not an uncertainty.** `Missing` is the state a user reads as "this conversation no longer exists". Reporting it because a laptop left the office network is a stronger error than reporting `Unavailable` on a conversation that really was deleted, because the second is cautious and the first invites the user to conclude their data is gone.

---

## 6. The amendment: observe the root before judging the conversation

These are proposed clauses, not edits to V3.1 and not implementation.

**UNC-23 — New term: root observation.** Before any row of §4 is evaluated, the scan observes the root directory itself and records exactly one of:

| `root_state` | Assigned when |
|---|---|
| `observed` | a directory listing of the root succeeds and the enumeration completes |
| `absent` | **only** on a `local-fixed` root, and only when the root observation fails with errno 2 or 3 |
| `unobservable` | every other failure — every transport or access code, every timeout, every `unknown` class, and any incomplete enumeration |

**UNC-24 — New term: root class.** Exactly one of:

| `root_class` | Assigned when |
|---|---|
| `local-fixed` | the root's volume is a drive letter and `GetDriveTypeW` returns 3 (`DRIVE_FIXED`) |
| `non-local` | the volume is a UNC volume (including `\\?\UNC\`), or `GetDriveTypeW` returns anything other than 3 |
| `unknown` | there is no volume, or the class could not be determined |

`unknown` is treated as `non-local` everywhere below. The class is a property of the **resolved physical volume**, not of the path string; a junction or symlink is OQ-2's question, not a licence to reclassify.

**UNC-25 — Row 1 gains a precondition.** Row 1 applies only when **all** of the following hold: `root_state = observed`; `root_class = local-fixed`; and the per-conversation observation fails with errno 2 or 3. Anything else that cannot observe a conversation is Unavailable, never Missing.

**UNC-26 — On a non-local or unknown root, `Missing` is unreachable.** A child observation that fails is Unavailable: reason `unreadable` for a transport, access or timeout code; reason `absence_unverified` for an absence-shaped code (2 or 3) that cannot be trusted because the share may have dropped between the root observation and the child observation. This is a deliberate loss of deletion detection on network roots, and it is the price of never claiming deletion from an unobservable location.

**UNC-27 — A root that is genuinely absent is not the same as an unobservable root.** With `root_state = absent` on a `local-fixed` root: if the index names conversations, each is Missing (`absent`), because the volume was observable and the root is not there; if the index names none, the scan reports no conversation verdicts. `root_state = absent` is not assignable on a non-local or unknown root, because the same failure code can mean a dropped connection (Observed: 53).

**UNC-28 — Error codes override the class.** Any transport or access code — observed 53 and 67; Source-derived 21 `ERROR_NOT_READY`, 55 `ERROR_DEV_NOT_EXIST`, 59 `ERROR_UNEXP_NET_ERR`, 64 `ERROR_NETNAME_DELETED`, 121 `ERROR_SEM_TIMEOUT`, 1231 `ERROR_NETWORK_UNREACHABLE`, 1232 `ERROR_HOST_UNREACHABLE` — is Unavailable on every class, including `local-fixed`. Absence is assertable only from codes 2 and 3, on `local-fixed`, with the root observed.

**UNC-29 — The negative rule: absence is never derived from a sentinel alone.** `os.IsNotExist(err)`, `errors.Is(err, fs.ErrNotExist)` and any similar mapping are insufficient, because 53 satisfies both (Observed). A verdict of Missing requires the root observation, the class and the code, each recorded. This is the same shape as the contract's other negative assertions: the conclusion must be derived from the named evidence, not from a convenience predicate.

**UNC-30 — The rejected alternative: confirm absence with a second root observation.** Re-observing the root after a child absence would narrow the window, but it cannot close it — the disconnect can occur between the two observations — and on a non-local root it would add two network round trips per missing conversation while still not proving a negative. It is recorded as considered and **not adopted**, because it buys the appearance of verification without the property. `V3-47` already prefers the cautious verdict.

**UNC-31 — `V3-60` extends to the root.** A failed root observation is an error, never a partial inference. A listing that begins and then fails mid-iteration is `unobservable`, not a shorter list of conversations.

**UNC-32 — One reason is added.** `absence_unverified` — Unavailable; content absence was suggested by an absence-shaped error but the root's class or continuity cannot support the claim. `unreadable` keeps its existing meaning. No other reason changes.

**UNC-33 — Totality survives the precondition.** `V3-66`'s partition gains one prior axis: root observation. Every input still reaches exactly one row, because `unobservable` routes to row 2 before rows 1 and 3–20 are considered, and `absent` on a local-fixed root routes to row 1 or to no verdict. No row is added and no state is added.

**UNC-34 — What the amendment deliberately does not do.** It does not probe the network beyond the scan's own root, does not retry, does not wait for a reconnect, does not distinguish a server rename from a deletion, does not add a seventh state, does not change `V3-47`, and does not authorize any user-facing text. It changes what the scanner is allowed to conclude, not what it does about it.

---

## 7. Outcome

**UNC-35 — Outcome: Amend.** OQ-1 is not a clarification of an existing rule and not a case for retaining it. V3.1's row 1 is contradicted by observed toolchain behaviour: an unobservable UNC root yields `IsNotExist = true` at both the root and the child, so the table's first row asserts deletion from a location the scanner never observed. The bounded amendment is:

| Clause | Change |
|---|---|
| §1 terms | add `root_state` (`observed` / `absent` / `unobservable`) and `root_class` (`local-fixed` / `non-local` / `unknown`) |
| `V3-39` row 1 | add the three-part precondition of UNC-25 |
| `V3-58` | add `absence_unverified` |
| `V3-60` | extend the error rule to the root observation |
| `V3-66` | add the root-observation axis to the totality statement |
| `V3-73` | OQ-1 closes, subject to owner adoption |

`V3-47`, row 2, the six states and the reason vocabulary's existing entries are unchanged. The amendment is one term pair, one precondition, one reason and two clarifications — not an extension.

**UNC-36 — If the owner prefers no new reason,** `absence_unverified` may be folded into `unreadable` without weakening the precondition; the root observation and the class rule are the load-bearing parts and are not optional. The cost of folding is that a user cannot tell "the share is gone" from "the file is unreadable", which is the distinction this record exists to preserve.

---

## 8. Tests this would require

- **UNC-T1** Lock the toolchain fact: `os.Stat` on an unresolvable UNC root returns errno 53 and `errors.Is(err, fs.ErrNotExist)` is true, so a future Go change cannot silently restore the unsafe predicate.
- **UNC-T2** Lock the asymmetry: a reachable-host/missing-share UNC path returns errno 67 and is not `ErrNotExist`.
- **UNC-T3** Lock the local baseline: errno 2 and 3 on a fixed local volume are `ErrNotExist`.
- **UNC-T4** An unobservable root (injected transport error) reports every indexed conversation as Unavailable, and zero Missing — including when the child observation itself returns 53.
- **UNC-T5** A non-local root observed successfully, with a child observation failing 2 or 3, reports Unavailable `absence_unverified`, never Missing.
- **UNC-T6** A local-fixed root observed successfully, with a child absent, reports Missing `absent`.
- **UNC-T7** A disconnect injected between the root observation and the child observation reports Unavailable, never Missing, and the scan does not retry.
- **UNC-T8** Root class detection: UNC volume → `non-local`; `DRIVE_FIXED` → `local-fixed`; `DRIVE_NO_ROOT_DIR`, `DRIVE_REMOTE`, `DRIVE_REMOVABLE`, `DRIVE_CDROM` and an unknown class → `non-local`.
- **UNC-T9** The negative assertion: no code path produces Missing from `os.IsNotExist`, `errors.Is(fs.ErrNotExist)` or any other sentinel alone; every Missing carries a recorded root observation, class and code.

---

## 9. Nothing here is authorized

**No implementation follows.** This record authorizes no scanner, no root probe, no `GetDriveType` call, no new reason constant, no change to any verdict, no retry, no reconnect wait, no repair, no migration, no marker change, no lock, no gate, no deletion enablement and no adoption of the refusal-only retry policy. It closes one open question as an amendment proposal; the owner decides, and the remaining five baseline questions stay open.
