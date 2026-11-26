# TASK-03 — D-2 and D-3: programmer's technical assessment and recommendation

| Field | Value |
|---|---|
| Purpose | Resolve what can be resolved of D-2 from code, state what cannot, and give the owner technical input for D-3 |
| Status | **Programmer assessment. Authorizes nothing. Code, gate, locks and index are all unbuilt and unauthorized** |
| Basis | Rows 584–586; the code read in this session; rows 567–569, 571, 572, 574 evidence already recorded |
| Edits | None. No code, contract, row or test changed |
| Recovery token | `TASK-03-D2-D3-PROGRAMMER-ASSESSMENT-2026-11-22` |

Labels: **[Observed]** read in the workspace this session · **[Source-derived]** stated by a cited row or clause · **[Calculated]** arithmetic on observed counts · **[Inferred]** engineering judgement from those · **[Unestablished]** not established here.

---

## §1 Role, and what this document is not

**PA-01 — Scope of the programmer's input.** D-2 (gate implementability) is an engineering question and this document resolves as much of it as evidence allows. D-3 (coverage acceptance) is a product statement; the input here is its technical consequence, not a decision.

**PA-02 — Nothing here authorizes work.** A resolved design, a cost estimate and a recommendation are not permission. Every code change named in §6 needs explicit implementation authorization, and none is granted by this record, by row 586, or by the disposition it records.

**PA-03 — No test was executed for this assessment.** All counts below are greps over source; no build, package test or probe was run for this document. The one verification available — that the tree still compiles — was last observed at `CP-2026-11-22-74` and is not re-claimed here.

---

## §2 D-2 is a work package, not a decision analysis can close

**PA-04 — The framing correction.** D-2 cannot be "resolved" the way OQ-1 and OQ-2 were. Those were questions about what the code does; answering them ended the question. Two of D-2's three blockers are questions about what the code should *become*, and their only honest resolution is a change to code. So the decision actually available is not "resolve or accept" but **"authorize a bounded slice, or leave D-2 open with its scope now precise"**.

**PA-05 — The two blockers that are code, and the one that is scope.** OQ-7a (path centralisation) and the enforcement self-test are code changes. V2-OQ-2 (marker schema) is neither: it is a prerequisite for a component that does not exist in the product at all. **[Observed]** no marker is written or read anywhere in `internal/`; a case-insensitive search for "marker" finds only unrelated audio test fixtures and an ASR transcript comment.

---

## §3 Blocker B-1 — path centralisation (OQ-7a)

**PA-06 — What exists now.** `PersistenceManager` builds conversation paths inline. The same literal appears **nine** times: `filepath.Join(pm.rootDir, "chats", …)` at `persistence.go` lines 85, 149, 193, 257, 472, 487, 556, 636 and 669. **[Observed]** **[Calculated]** nine sites.

**PA-07 — No constructor exists.** A search for any directory-constructor function in `internal/` returns only the reader helper `loadHistoryFile(chatDir string)`, which takes a path rather than producing one. **[Observed]** So there is no single place to assert containment against, which is precisely why `V3-13b`-style containment cannot be enforced today.

**PA-08 — Three conventions, confirmed.** The conversation tree is `pm.rootDir/chats/<id>`; the shadow tree is `paths.StateRoot()/Chats/<id>` via `chat/service.go:148,515` and `memory/service.go:43,57`; and resource resolution uses `paths.ProjectRoot()`. **[Observed]** These are three distinct derivations, not one root with three spellings, so a constructor cannot be written against "the root" without first naming which root it belongs to.

**PA-09 — Resolution: one constructor, one containment point.** Add `conversationDir(id string) (string, error)` and `chatsRoot()` to the chat package, returning a path that has passed `paths.CheckContained` against the history root. PathService already supplies both primitives — `CheckContained` and `FileIdentity` — and `CheckContained` is already proven through a junction alias (row 578). **[Source-derived]** The constructor therefore consumes an existing capability rather than inventing one.

**PA-10 — Feasibility: high, and mechanical.** Nine call sites in one package; one external consumer of `RootDir()` (`app-wipe.go:40`); the change is compile-checked, touches no document format, and adds no file. **[Observed]** The risk is not difficulty, it is scope creep: a constructor that also "tidies" the shadow tree would silently become the OQ-7d ownership change, which is out of scope and separately owned (row 573).

**PA-11 — Why this cannot be decided on paper.** The constructor's signature is a design choice, but its *effect* — that no write can be issued against an uncontained path — is a property of compiled code. Until it is built, `containment` remains an assertion about a helper that has no caller, which is exactly the gap row 578 recorded for `identityOf`. **[Source-derived]**

---

## §4 Blocker B-2 — the enforcement self-test

**PA-12 — Why a startup self-test is the wrong instrument.** A startup test proves what it happened to exercise; row 563's finding stands. **[Source-derived]** A path-dependent test cannot establish a universal claim, and the claim the gate needs is universal: *no write happens outside the boundary*.

**PA-13 — What the code shape is today.** `internal/backend/chat` contains **eight** exported mutators (`NewChat`, `SaveChat`, `SaveChatAs`, `DeleteChat`, `SetLastChatID`, `RewriteHistory`, `AppendHistoryEntry`, `RenameChat`) over **fourteen** direct filesystem write calls — eleven in `persistence.go` (`MkdirAll` ×5, `WriteFile` ×3, `Rename` ×2, `RemoveAll` ×1) and three in `service.go` (`WriteFile` ×3). **[Observed]** **[Calculated]** 8 mutators, 14 write sites, 2 files.

**PA-14 — Resolution: make the boundary structural, then assert the structure.** Route all fourteen write sites through unexported helpers behind one exported entry point, and replace the startup test with a **build-failing analysis test**: parse `internal/backend/chat`, and fail if any function other than the designated writer references `os.WriteFile`, `os.Rename`, `os.MkdirAll`, `os.Mkdir`, `os.Remove` or `os.RemoveAll`. That is deterministic, path-independent, and a positive build-bound result — the shape `V3-36` permits. **[Source-derived]** `V3-36`; **[Inferred]** feasibility.

**PA-15 — The exception list is the whole difficulty.** `service.go:531,541,605` write the *shadow* tree, in the same package but a different root and a different owner. A single assertion over the package must either (a) name those three as listed exceptions with a stated reason — which is a checklist with extra steps — or (b) move the shadow-tree writer out of the package, which is an ownership change (row 573) and therefore a separate, larger decision. **[Observed]** **[Inferred]** Recommendation: take (a) for the slice, with the exception list as data in the test and a stated reason per entry, and record (b) as the honest end state rather than smuggling it in.

**PA-16 — Blast radius of the shape change.** The App layer calls these mutators from **four** files at **thirteen** sites (`app-chat.go`, `app-prompt.go`, `app-characters.go`, `app-wipe.go`), and **seven** probe programs under `build/tmp/` call them directly. **[Observed]** **[Calculated]** The package-internal change is small; the call-site migration is the bulk of the mechanical work, and the probe programs are throwaway and should be rebuilt or retired rather than adapted.

---

## §5 Blocker B-3 — V2-OQ-2 and the marker half

**PA-17 — The marker does not exist.** No code writes or reads a marker, a generation, a settlement field or a retraction record. **[Observed]** So the half of the gate that V2-OQ-2 blocks has no implementation to block; it is blocked on a schema decision for a component that would have to be built from nothing.

**PA-18 — Scope exclusion, stated as an exclusion.** The refusal-only slice of row 585 `RQ-12` needs per-conversation locks, the tail rule, containment and a boundary — no marker field. **[Source-derived]** Within that slice V2-OQ-2 is *outside*, not resolved, and OQ-6 is likewise not pursued. The accepted costs are the ones the contract already states: rows 16–17 stay unreachable, `gate_enforced` stays false, and generation ≥2 content stays `Pending`.

**PA-19 — What that exclusion is not.** It is not a claim that the marker half is unimportant, and it is not permission to add a marker later without OQ-6 and V2-OQ-2. It is a statement that the gate's demonstrable part is separable from its schema-dependent part. **[Inferred]**

---

## §6 Slice S-1 — the bounded work an authorization would cover

**PA-20 — Scope, stated so it can be authorized or refused as one unit.** S-1 is:

1. `conversationDir(id)` / `chatsRoot()` with containment against the history root; nine call sites replaced; `RootDir()` kept for `app-wipe.go`.
2. One exported mutation entry with a typed refusal result; the eight mutators reduced to wrappers or unexported.
3. Per-conversation lock in a sibling namespace outside `chats/` (row 567), refusing on contention rather than waiting.
4. Tail rule in the append path: terminate the unterminated tail in the same critical section, or refuse (row 579).
5. The build-failing analysis test of PA-14, with an explicit, reasoned exception list.
6. App-layer refusal surfacing as a typed error at the thirteen call sites, so `V3-27`'s "silence is non-conforming" is met without adopting the seven-field result.

**PA-21 — What S-1 explicitly does not include.** The marker, CAS, settlement, generation, index producer, repair of any kind, the memories architecture, creation repair (`V3-72a`), the seven-field caller result, retry policy, and any change to a document format. **[Source-derived]** rows 584–586 scope statements.

**PA-22 — The one storage side-effect, named rather than denied.** S-1 creates lock files in a sibling namespace. They are transient, they change no document format, and they add a new file class outside `chats/` that a future scanner must be told about. Calling S-1 "no storage change" would be false; "no document-format change, one new transient file class" is accurate. **[Observed]** by construction from row 567's placement finding.

**PA-23 — Acceptance criteria, not executed tests.** S-1 would be demonstrable when: the analysis test fails on a deliberately introduced bypass; two processes contending on one conversation produce a typed refusal and no interleaved write; an append against an unterminated tail either terminates or refuses, never concatenates; a junction alias inside the root is refused a path outside it; `go build ./...` and the chat package tests pass. None of these has been run.

**PA-24 — Size, honestly.** One package plus one calling layer; nine path sites, fourteen write sites, thirteen call sites, one lock mechanism, one analysis test. **[Calculated]** No migration, no format change, no data touched. It is small enough to authorize in one piece and large enough that it should not be spread across several.

---

## §7 D-3 — the technical consequence of accepting the coverage boundary

**PA-25 — Mechanically, acceptance costs nothing.** Accepting the boundary requires no code, closes no row, and contradicts nothing in v3.2. It makes an existing behaviour explicit: with no index producer, `index_state = present` can never hold, so `Missing` and `Orphaned` never fire. **[Source-derived]** `V3-13d`, rows 1 and 20, row 586 `DP-10`.

**PA-26 — Acceptance does not weaken the gate.** The gate's value is that no new damage is written outside the boundary. That value does not depend on `Missing` or `Orphaned` existing: damage detected inside candidates is still detected, and damage prevented is still prevented. **[Inferred]** So D-3 is not a technical precondition of D-2, in either direction.

**PA-27 — The index producer is a separate feature with its own prerequisites.** It needs a write path, an owner, and membership by logical identity only (`V3-13d`; row 578). Row 573's precedent is directly relevant: an index written by an ungoverned writer would be a second authority over the same fact, which is the defect that row removed from settlement. **[Source-derived]** **[Inferred]** A trustworthy index is cheaper to build *after* S-1 than before it.

**PA-28 — The one thing acceptance should not be read as.** It is not a statement that unindexed data does not exist, and not a statement that a clean scan proves health. `V3-13a` already forbids both readings; acceptance only makes the limit official. **[Source-derived]**

---

## §8 Recommendation

**PA-29 — For D-3: accept the boundary, as a product statement, and do it first.** It is free, it is reversible by a later explicit act, it removes an ambiguity the product currently carries silently, and it does not commit anyone to building an index. This is technical input to a product decision; the decision remains the owner's.

**PA-30 — For D-2: authorize S-1 as one bounded unit, or leave D-2 open.** D-2 cannot be closed by more analysis. If S-1 is not authorized, the correct statement is "D-2 open, scope now precise", not "resolved". If it is authorized, the three blockers resolve as: B-1 by construction (§3), B-2 by capability shape plus build-failing analysis (§4), B-3 by explicit exclusion (§5).

**PA-31 — Order between them, and why it is not a dependency.** Accepting D-3 first costs nothing and changes no code. S-1 then proceeds without waiting on an index decision, and the index question stays open on its own merits. **[Inferred]** If the owner prefers to authorize first and accept later, nothing in S-1 depends on the acceptance; the recommendation is ordering, not coupling.

**PA-32 — What stays outside regardless.** Repair, the writer defects of rows 570, 572, 574 and 580, the memories architecture, and any marker work. S-1 prevents new damage through the boundary; it repairs nothing and makes no existing verdict stronger. **[Source-derived]**

---

## §9 Checks this assessment implies

Not executed.

- **PA-T1** If S-1 is authorized, the analysis test must fail on an introduced bypass before it is trusted; a test that has never failed has not been shown to detect anything.
- **PA-T2** The exception list in that test must carry a reason per entry; an unexplained exception is the checklist the blocker was opened about.
- **PA-T3** No S-1 change may alter a document format, a verdict rule, a state name or a row; a change that does is a different authorization.
- **PA-T4** D-3's acceptance, whenever it is recorded, must not be read as adoption of v3.2 — row 586 `DP-18` governs.

---

**Scope of this document.** It resolves the design of two blockers, excludes the third by scope, sizes the slice that would demonstrate implementability, and gives technical input on D-3. It builds nothing, tests nothing, adopts nothing, and authorizes nothing.
