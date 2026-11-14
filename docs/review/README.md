# Review artifacts — project-wide review and evidence record

> **Directory scope:** This directory contains review, audit, decision, implementation-evidence, and
> historical artifacts for T2.9.5B audio, F3-C lorebook/world/scenario ownership, the Amendment 02
> cycle, requirements reconciliation, and ROOT-01. It is **not** limited to T2.9.5B audio. Start at
> `docs/README.md`; use `docs/DOCUMENT-REGISTRY.md` for every file.
>
> **Authority:** The sections below are the original T2.9.5B review index and remain historical
> evidence. They do not describe the whole directory and do not override later frozen contracts,
> owner decisions, or the master registry. For T2.9.5B, the frozen controlling pair is
> `docs/t2.9.5b-authorization-decision-register-r14.md` and its checksum, recorded by
> `docs/t2.9.5b-authorization-decision-register-r14.freeze.md`. The r14 file's own `DRAFT` header is
> the pre-freeze register header; the freeze record controls.
>
> Nothing here is an approved decision, authorization, measurement, or plan of record unless a
> specific record says so. See `docs/AMENDMENT-REGISTER.md` and `docs/FROZEN-ARTIFACTS.md`.

This directory contains research and review artifacts captured for external audit. Except for the
governing architecture contract identified in §10–§11 and the adopted §17.19 clarification recorded in
§21–§22, nothing here is an approved decision, an authorization, a measurement, or a plan of record.
The controlling T2.9.5B authorization register remains `docs/t2.9.5b-authorization-decision-register-r14.md`,
frozen by `docs/t2.9.5b-authorization-decision-register-r14.freeze.md`; these files neither amend nor
supersede it. Directory-wide authority is in `docs/README.md`.

## 1. `t2.9.5b-audio-research-correction-addendum-a1.md`

- **What it is:** Correction Addendum A1, adding corrections to the T2.9.5B audio research report
  (compiled evidence vs functional capability, route/handler establishment, per-family model
  compatibility, corrected capability matrix, version/upgrade resolution, revised tranche boundary,
  selection/resource rules, licensing).
- **Status:** UNREVIEWED RESEARCH ARTIFACT. Not owner-approved. Contains no implementation
  authorization. Its §6 explicitly declines to authorize any tranche.
- **Provenance:** produced in a Codebuff/Freebuff agent session working in
  `C:\Users\StoneWolfPC\HOWL_Chat` on 2026-11-04, in the same conversation as the research report it
  corrects. It was delivered as conversation text only and was never previously written to disk.
  This file is a verbatim transcription of that delivered text; no wording was reconciled, corrected,
  extended, or summarized in the capture.
- **Capture date:** 2026-11-04.

## 2. `AUDIT-t2.9.5b-audio-research-a1.md`

- **What it is:** claim-by-claim re-verification of Addendum A1 against primary sources (repository
  lines at the audited `HEAD`, llm.cpp files fetched at commit `a34de8618`, Hugging Face model/licence
  metadata, release notes).
- **Status:** UNREVIEWED AUDIT ARTIFACT. It records errata against A1; it does not amend or replace A1,
  and A1's hash is unchanged.
- **Result:** no central conclusion overturned; 3 wording findings (realtime branch condition, realtime
  output modality, licence grouping), 5 citation-range imprecisions, 3 unverified items.
- **Capture date:** 2026-11-04.

## 3. `AUDIT-t2.9.5b-full-reverification-muse.md`

- **What it is:** Muse's independent adversarial re-verification of A1 and of the DeepSeek audit,
  plus seven new application findings (single-image/history ceiling, JPEG MIME hardcode,
  `llama-connect.exe` timestamp outlier, Go directive vs toolchain, port-8080 lifecycle, dead
  router-mode client code, stale audio-architecture doc).
- **Status:** UNREVIEWED AUDIT ARTIFACT. Read-only; created no other file.
- **Capture date:** 2026-11-04.

## 4. `RECON-t2.9.5b-deepseek-muse-audio.md` (RECON-1)

- **What it is:** first reconciliation of A1, the DeepSeek audit and Muse's re-verification.
- **Status:** SUPERSEDED by RECON-2, preserved unchanged as evidence. It contains three defects
  corrected in RECON-2: wrong `/v1/realtime` line numbers, a false "tree drift" explanation for the
  status-flag citations, and a mischaracterised correction of Muse's dead-code finding.
- **Capture date:** 2026-11-04.

## 5. `RECON-t2.9.5b-deepseek-muse-audio-r2.md` (RECON-2, authoritative)

- **What it is:** consolidated reconciliation containing both auditors' positions explicitly, with
  source adjudication per claim; withdraws the DeepSeek audit's F11 as factually wrong.
- **Status:** AUTHORITATIVE RECONCILIATION, revision r2.1. Research/reporting only; authorizes nothing.
- **Result:** 22 claim rows — 11 confirmed, 7 corrected, 2 unestablished, 1 withdrawn (F11), 1
  downgraded to a non-debt observation — plus 4 auditor-level defects counted separately (3 in RECON-1,
  1 in the DeepSeek audit). No central conclusion overturned; A1, both audits and RECON-1 preserved
  byte-identical.
- **Revision r2.1:** integrity and evidence-grade pass. Technical-debt IDs are namespaced `DEBT-*` so
  they cannot be confused with authorization-register `T6`; the cosmetic `MMPProjPath` row left the
  debt register; the tally is reconciled and the D2 framing defect is accounted for; DEBT-T9 is
  regraded provisional (source mechanism cited, production reach unestablished, no fix prescribed);
  the licensing finding is narrowed to a `runtime/llama/` notice gap; the `llama-connect.exe` disposal
  recommendation is withdrawn pending identification; Vision carries four separate evidence grades;
  and owner decisions are split into owner approval, research prerequisites, and register-controlled
  items. See RECON-2 §16.
- **Capture date:** 2026-11-04.

## 7. `AUDIT-t2.9.5b-tranche-a-implementation-muse.md` (Tranche A implementation audit)

- **What it is:** Muse's independent, read-only audit of the completed Tranche A
  implementation against its authorization, RECON-2 r2.1, repository state, and the
  tranche implementation report. Not an authorization; implementation code unchanged.
- **Result:** PASS WITH CORRECTIONS — fail-closed admission throughout and genuine
  debt fixes, but the native-gate binary evidence is invalid (the stub
  `llama-server.exe` was scanned instead of `llama-server-impl.dll`, which contains
  the `/v1/audio/transcriptions` surface), the contract's auto-selection claims
  overstate enforcement (`AutoSelect`/`Fallbacks`/resolver unwired; auto-ASR bypasses
  dimensions), and saved engines are not re-applied at startup. Five ordered corrections.
- **Capture date:** 2026-11-04.

## 6. Original T2.9.5B audio research report — NOT RECOVERED

- **Status:** UNAVAILABLE. No copy exists in this workspace, and no verbatim copy is available to the
  capture session.
- **What was searched:** the whole repository (markdown, docs, build logs, root-level notes) for the
  document itself and for its distinctive content — `Qwen3-Omni`, `Voxtral`, `LFM3-Audio`, `5B.1`,
  `Correction Addendum`, `audio research`, `llama.cpp-native`, `b12015-a34de8618`, plus filename
  patterns for `research`/`report`/`addendum`. Zero matches outside this directory.
- **Why it cannot be transcribed:** the report body survives only as a condensed, paraphrased session
  summary, not as its original text. Reconstructing it from that summary would produce a different
  document and would misrepresent it as the original, so no file has been created in its place.
- **What is known about it:** it was delivered as conversation text earlier in the same session
  (2026-11-04) and covered the bundled llama.cpp audio capabilities, a capability matrix, the audio
  model families, and a proposed native-ASR tranche boundary ("5B.1"). Its claims are quoted and
  dispositioned individually in Addendum A1 §10.
- **To recover it:** it must be re-supplied from the source that holds the original conversation
  transcript. If re-supplied, place it here unchanged and record its hash in this note.

## 8. `AUDIT-kokoro-g2p-feasibility-challenge-2026-11-06.md` (Kokoro G2P research challenge)

- **What it is:** an independent challenge pass over the Kokoro G2P feasibility study in
  `build/tmp/kokoro-g2p-feasibility/`, plus the licensing/provenance investigation requested for the
  remaining gates. It re-verifies all six evidence hashes, re-executes the existing tokenizer probe,
  and records rights holders, contact channels and the three-permission split (use / redistribute /
  bundle) per component.
- **Status:** INVESTIGATIVE FINDING AND REVIEW ARTIFACT. Not an authorization, not a licensing
  opinion. It authorizes no download, spike, dependency, or production change; the register r14
  remains controlling.
- **Result:** the technical ranking survives and the MisakiSharp HOLD stands, but three recorded
  reasons are corrected — the v0.3.0 tokenizer fix is unproven and belongs to a different module than
  the one measured; shared ONNX Runtime is an unresolved ownership conflict and ORT is mandatory even
  without the neural fallback; MisakiSharp is not English-only (that is `misaki-rs`). Two gates are
  added (ORT ownership/version/packaging, and reference parity/OOV quality), and `misaki-rs`'s default
  build is recorded as statically vendoring GPL-3.0-or-later eSpeak NG. Also notes that
  `start-howl-chat.bat` and `build_howl_cgo.bat` disagree on `CGO_ENABLED`.
- **Capture date:** 2026-11-06.

## 9. `INVESTIGATION-go-kokoro-architecture-packaging-2026-11-06.md` (Go-native Kokoro/G2P hosting)

- **What it is:** read-only comparison of in-process Go/CGO/ONNX Runtime hosting versus a separately
  built helper process for a future Go-native Kokoro/G2P implementation. It examines Howl's build and
  launch paths, CGO discrepancy, audio activation boundary, ORT ownership, packaging, IPC options,
  lifecycle, failure containment, cancellation, and maintenance consequences.
- **Status:** INVESTIGATION ARTIFACT. Substantive findings accepted by the project owner on
  2026-11-06. It selects no architecture, approves no runtime DLL, authorizes no implementation, and
  does not amend any prior report. The controlling authorization register remains r14.
- **Result:** neither architecture is selected. In-process hosting requires a deliberate release-
  contract change (CGO, Go 1.25+, app-local ORT/CRT packaging, process-global ORT ownership, and an
  explicit native-failure policy). A helper preserves the main executable's `CGO_ENABLED=0` shape if
  built separately and not imported, but adds a second artifact, protocol, lifecycle, IPC security,
  and unmeasured cancellation/audio behavior. Ordinary TTS errors can remain per-capability; an
  in-process native fault can still terminate the app.
- **Documentation reconciliation (2026-11-06):** renamed to the requested filename from
  `INVESTIGATION-go-native-tts-architecture-packaging-2026-11-06.md`; exact working-tree source
  SHA-256 hashes added where practical; this index entry added. The report records that the worktree
  was already dirty and distinguishes source evidence from unestablished Windows execution.
- **Capture date:** 2026-11-06.

## 10. `DESIGN-howl-tts-helper-architecture-2026-11-06.md` (TTS helper architecture)

- **What it is:** design specification for the owner-selected separate-helper architecture for
  Howl's optional native Kokoro/G2P capability. It defines the application-facing contract, a
  provider-neutral versioned protocol, provider adapters, process lifecycle and failure handling,
  file-path audio ownership, packaging, versioning, character-voice integration, security, open
  questions, staged implementation, and acceptance criteria.
- **Status:** GOVERNING ARCHITECTURE CONTRACT; Stage 0 contract frozen in §17 and amended by adopted
  Clarification A2 (§17.19). Design only; it authorizes no implementation, dependency addition,
  build, launch, Windows spike, inference, model/voice operation, production edit, commit, or push.
  The controlling authorization register remains r14.
- **Result:** main Howl remains `CGO_ENABLED=0`; the helper is a separately built optional provider
  using parent-owned stdio NDJSON, a single `Wait()` owner, job-object ownership, typed failures, and
  app-owned file-path audio. The design keeps `types.TTSSynthesizer` stable, adds an optional provider
  inspector, and keeps Kokoro/G2P behind a provider boundary. No architecture is implemented.
- **Open items:** resolved or classified in the Stage 0 entry below. Remaining product/legal
  decisions are catalog-start policy, output retention/cleanup, `auto` fallback, backend-switch
  timing, exact ORT/runtime licensing, and Kokoro/G2P asset layout/Windows build.
- **Capture date:** 2026-11-06.

## 11. Stage 0 — Frozen Protocol v1 Contract

- **What it is:** the project-owner-approved Stage 0 freeze added as §17 of
  `DESIGN-howl-tts-helper-architecture-2026-11-06.md`. It preserves the original proposal and freezes
  the exact NDJSON envelopes, version negotiation, operation semantics, reader/worker concurrency,
  provider/capability/extension rules, state machine, typed error registry, WAV output contract,
  process lifecycle and cancellation, compatibility rules, numeric-policy classification,
  implementation-ready interface summary, and remaining owner decisions.
- **Status:** GOVERNING CONTRACT — STAGE 0 FREEZE INCLUDING CLARIFICATIONS A1 AND A2; A2 ADOPTED
  2026-11-07. Design/specification only. It authorizes no Stage 1, Stage 2, Stage 3, implementation,
  protocol package, helper binary, dependency addition, production edit, build, launch, ONNX Runtime
  initialization, inference, model/voice operation, asset change, commit, or push. The controlling
  authorization register remains r14.
- **Result:** v1 is protocol 1.0 with `{major,minor}` versions, top-level namespaced extensions,
  strict core schemas, one protocol reader plus one operation worker, non-blocking control during
  synthesis, one active synthesis per helper, a two-process cap, parent-owned stdio, one `Wait()`
  owner, typed errors, and validated PCM16 mono WAV output in an app-owned cache directory. The main
  executable remains `CGO_ENABLED=0`; `types.TTSSynthesizer` is unchanged; the inspector is additive.
- **Stage 0 clarification A1:** §17.18 adds an explicit parent acceptance barrier, cancellation
  precedence, terminal publication barrier, kill/wait-before-cleanup ordering, fatal duplicate-key
  classification, and unmatched/retired/duplicate-terminal handling. A1-D7 records
  `protocol_violation` → `ErrCodeSynthesisFailed` as a public compatibility projection only, not
  semantic equivalence; the precise internal code and detail must survive for diagnostics. It is an
  additive amendment; O-1 through O-7 remain unresolved and no Stage 1 work is authorized.
- **Stage 0 clarification A2:** §17.19 adopts the owner-accepted Policy B capability floor,
  handshake semantic classification, service/provider status projection, `unavailable` recovery,
  and lifecycle clarifications. The accepted Revision 6 amendment remains standalone at
  `AMENDMENT-DRAFT-howl-tts-capability-policy-c1-17.19-2026-11-07.md`; the adoption record is
  `ADOPTION-howl-tts-capability-policy-c1-17.19-2026-11-07.md`. Wire DTOs remain unchanged;
  `TTSCapability` is an application projection. O-1–O-7 remain unresolved; B12a/B12b remain
  deferred.
- **Frozen:** envelope/schema strictness, 4 MiB frame and 64 KiB text limits, one active synthesis,
  two helper processes, one loaded provider, catalog-without-load, no unsolicited events, no HTTP or
  named pipes, no automatic replay, WAV format, path containment, and the process/readiness split.
- **Provisional:** 5 s handshake, 120 s load, 120 s synthesis, 2 s cancel grace, 5 s shutdown, 120 s
  idle close, 3 restarts/10 min, 1/2/5 s backoff, and 8 KiB stderr ring. These are implementation
  defaults requiring measurement, not evidence-backed guarantees.
- **Remaining owner decisions:** O-1 catalog startup policy; O-2 output retention/cleanup; O-3 `auto`
  fallback; O-4 backend-switch timing; O-5 ORT/runtime licensing and packaging; O-6 Kokoro/G2P asset
  layout and Windows build; O-7 provider/module split after v1. O-1/O-2 may use provisional behavior
  in fake-helper tests; O-4/O-5/O-6 must close before native production work.
- **Source corrections recorded:** the earlier integer protocol field and `params.extensions` examples
  are superseded by §17; the existing `GetTTSVoiceCatalog(engine string)` has no context and is not
  claimed to be provider-aware; cache eviction is metadata-only and does not delete audio files.
- **Capture date:** 2026-11-06.

## 12. `AUDIT-howl-tts-helper-stage1-adversarial-2026-11-06.md` (Stage 1 adversarial review)

- **What it is:** an independent adversarial review of the Stage 1 implementation report
  (`IMPLEMENTATION-howl-tts-helper-supervisor-stage1-2026-11-06.md`, SHA-256 verified). It inspects
  the implementation and tests directly, re-executes every reported verification, and answers the
  three scrutiny questions (cancellation races, process lifecycle, WAV cleanup) plus the
  restart/backoff/idle-close contract question.
- **Status:** INVESTIGATIVE FINDING AND REVIEW ARTIFACT. Stage 1 remains HOLD; Stage 2 is not
  authorized. No code was changed; one temporary probe was run and deleted.
- **Result:** all evidence reproduces (hashes, tests, race, vet, CGO build, chat regression). The
  acceptance barrier is sound; single-`Wait()` holds; no autonomous restart exists. One MEDIUM finding
  (proven by probe): `failAll` blocks forever on a full outcome channel, skipping `Kill` — must be
  fixed before the HOLD lifts. Two LOW findings (handshake-failure reader leak; cleanup-before-kill on
  the violation path), one contract non-conformance carried open (per-component reparse-point checks),
  and one Stage 1 boundary omission (golden fixtures). Recommendation: conditional acceptance.
- **Capture date:** 2026-11-06.

## 13. `AUDIT-howl-tts-helper-stage2-adversarial-2026-11-06.md` (Stage 2 adversarial review)

- **What it is:** an independent adversarial review of the Stage 2 real-process interoperability
  report (`IMPLEMENTATION-howl-tts-helper-fake-process-stage2-2026-11-06.md`, SHA-256 verified). It
  verifies all 22 declared hashes, reads the new fake-helper and interop sources directly, re-runs
  every reported command, and answers the seven scrutiny questions (protocol parity, process
  ownership, terminal ordering, WAV ownership, failure isolation, test independence, scope).
- **Status:** INVESTIGATIVE FINDING AND REVIEW ARTIFACT. Stage 2 accepted for its milestone;
  Stage 3 not authorized. No code changed; one temporary parity probe was run and deleted.
- **Result:** the real-process boundary is genuine and reproducible (11 real-child tests, ownership
  enforced by construction, real fault injection at the process boundary, clean scope). A parity
  probe added here proved the fake codec accepts all eight contract-native golden requests — a
  direction the submission never tested. One MODERATE gap: capability negotiation is unimplemented
  on the parent side and `hello.response.json` advertises non-registry capability IDs, so a helper
  can silently omit `cancel`/`synthesize` and still pass handshake. Minor: the "renames without
  overwriting" claim is false for `os.Rename` on Windows; kill-before-delete ordering is inferred, not
  observed at the OS level. The disclosed full-suite flake was independently reproduced with the same
  signature and confirmed pre-existing by mtime.
- **Capture date:** 2026-11-06.

## 14. `AUDIT-howl-tts-helper-capability-correction-2026-11-06.md` (correction adversarial review)

- **What it is:** an independent adversarial review of `CORRECTIVE-howl-tts-helper-capabilities-no-overwrite-2026-11-06.md`
  and the updated Stage 2 report (both SHA-256 verified), scoped to S2-1, S2-2, S2-3, and S2-5. It
  verifies all declared hashes, confirms the eleven declared files changed and nothing else, audits
  every operation gate above the only two frame-emission sites, and tests the hard-link publication
  path.
- **Status:** INVESTIGATIVE FINDING AND REVIEW ARTIFACT. No blocking discrepancy; correction gate may
  close. Stage 3 remains unauthorized. No code changed; one temporary probe was run and deleted.
- **Result:** capability registry, membership, version, and both-side gates are upheld and the golden
  fixture is now validated rather than merely decoded, closing the drift hole that produced S2-1.
  Publication is atomic and leak-free. Three findings: the mandatory five-service rule is stricter than
  §17.6.1 and needs a contract line (C-1, non-blocking); `Unload`/`Cancel` gate before the lifecycle and
  so report a false capability diagnosis on a never-started client (C-2, latent); hard links are
  NTFS/same-volume and the pattern should not be inherited blindly by a real helper (C-3).
- **Capture date:** 2026-11-06.

## 15. `AUDIT-howl-tts-capability-policy-c1-2026-11-06.md` (C-1 capability policy review)

- **What it is:** an independent adversarial review of `DECISION-howl-tts-capability-policy-c1-2026-11-06.md`
  (SHA-256 verified) against the frozen architecture contract (SHA-256 verified). It rechecks every cited
  clause, reconstructs the ambiguity from the source text, and assesses the proposed §6 wording clause by
  clause.
- **Status:** INVESTIGATIVE FINDING AND REVIEW ARTIFACT. PASS WITH FINDINGS. No contract or report edit,
  no implementation, no Stage 3. Read-only.
- **Result:** the owner's Policy B with a `load`+`synthesize` floor is the correct decision and is better
  supported by the contract than the report argues — §5.3 allocates only `protocol_incompatible` and
  `helper_start_failed` as `hello` failure states, which makes the current `protocol_violation`
  classification non-conformant under either policy. The wording cannot be adopted as drafted: the
  transition table has no row reaching `incompatible` for a capability-short hello (blocking), and the
  proposed reservation of `protocol_violation` contradicts four remaining handshake sites (blocking).
  Also identified: §17.11.6 rule 4 still mandates an unconditional `cancel`, no parent rule exists for
  identity change without `unload`, `Available` is defined in terms of catalog-derived knowledge, and the
  status projection conflates process-scoped with provider-scoped capabilities.
- **Capture date:** 2026-11-06.

## 16. `AUDIT-howl-tts-capability-policy-c1-rev2-2026-11-07.md` (amendment revision 2 review)

- **What it is:** an independent adversarial review of `AMENDMENT-DRAFT-howl-tts-capability-policy-c1-17.19-2026-11-07.md`
  (SHA-256 verified, 229 lines) against the frozen architecture and the C-1 audit. It verifies all 34
  line references, traces every B1–B12 disposition into the normative text, and checks the 29 impact rows
  against the supersession map.
- **Status:** INVESTIGATIVE FINDING AND REVIEW ARTIFACT. PASS WITH FINDINGS — conditional. Read-only.
- **Result:** the policy substance is adoptable; the text is not. Two blocking defects: the impact list
  prescribes in-place rewrites for five clauses that A2-D11 promises to preserve (§2 rows 4, 7, 12, 15,
  19), and §17.7.4 cannot represent the new `Available=true` + `State=unavailable` + live-process
  combination while impact row #17 says no change is needed. Six moderate findings cover a missing
  `unavailable` retry transition, an over-narrow `ProcessAlive` supersession, a duplicated idle-close
  transition, a rejected-shutdown mapping that contradicts the `shutdown` error code, an A2-D7 step that
  presupposes pre-empting an active synthesis despite the serialization lock, and incomplete
  `ServiceCapabilitiesKnown` invalidation. Minor items include a wrong bound citation (§17.2.4 instead of
  §17.2.5), two off-by-two line numbers, and two dispositions asserted only in the traceability table.
- **Capture date:** 2026-11-07.

## 17. `AUDIT-howl-tts-capability-policy-c1-rev3-2026-11-07.md` (amendment revision 3 review)

- **What it is:** an independent adversarial review of revision 3 of
  `AMENDMENT-DRAFT-howl-tts-capability-policy-c1-17.19-2026-11-07.md` (SHA-256 verified, 246 lines)
  against the frozen architecture, the C-1 decision, and the revision 2 audit. It re-verifies 32 cited
  line ranges, disposes of R1–R14, recomputes the 12 decisions / 29 rows / 14 dispositions, and checks
  every legal combination of `State`, `ProcessState`, `Available`, `Loaded`, `Ready`, and
  `ServiceCapabilitiesKnown` for representability.
- **Status:** INVESTIGATIVE FINDING AND REVIEW ARTIFACT. PASS WITH FINDINGS — conditional. Read-only.
- **Result:** both prior blockers are resolved and all fourteen dispositions now appear in normative text,
  so the additive method, the state model, and the capability-validity invariant are coherent. One blocking
  defect remains: A2-D6 requires `ServiceCapabilities` to report `Supported:true`/`Supported:false`, but
  §17.12's implementation-ready `TTSCapability` is `{ID, Version int}` with no `Supported` field, while
  impact row 3 says no type change is needed (and `Version` is `string` in §3.2, so the same new field
  would carry two element types). Four moderate findings: A2-D12's retry guard rests on a helper
  transition §17.7.3 does not define (§17.7.3 is absent from the whole impact map), the retry trigger "next
  explicit TTS operation" could load a model on a voice-list query inverting §5.3 and colliding with O-1,
  A2-D9 leaves the observed-terminal cancellation branch without a cleanup order while §17.11.4 step 7 is
  uncited, and §3.2 cannot carry the `ProcessState` distinction the new §17.7.4 row branches on. Minor:
  §17.19.4 item 6 drops `busy` from the unload row, §17.7.1 should be 1753–1755 not 1755, §17.11.5 should
  be 2089–2090 not 2088–2090, four feature capability IDs are left unclassified, and the new row omits
  `ServiceCapabilitiesKnown`.
- **Capture date:** 2026-11-07.

## 18. `AUDIT-howl-tts-capability-policy-c1-rev4-2026-11-07.md` (amendment revision 4 review)

- **What it is:** an independent adversarial review of revision 4 of
  `AMENDMENT-DRAFT-howl-tts-capability-policy-c1-17.19-2026-11-07.md` (SHA-256 verified, 278 lines)
  against the frozen architecture, the C-1 decision report, and the revision 3 audit. It disposes of
  F1–F5 and F6–F13, re-verifies 34 cited line ranges, recomputes the 12 decisions / 31 impact rows /
  14 dispositions / 6 F rows, and re-checks all thirteen `State` values against the seven status fields.
- **Status:** INVESTIGATIVE FINDING AND REVIEW ARTIFACT. PASS WITH FINDINGS — recommend adoption after one
  clarification. Read-only.
- **Result:** the blocking type defect is resolved without touching the wire DTO — §17.12 is named the
  authoritative status shape, one application `TTSCapability` is defined, and impact rows 3 and 25 record
  the change instead of denying it. F2–F5 are closed in normative text: the recovery rule returns a typed
  non-fatal load failure to `unloaded`, only an explicit `load` triggers it, both cancellation branches
  are stated with step 7 cited, and §3.2 no longer carries authority. One moderate defect remains: A2-D7
  and A2-D12 both claim a differing-identity explicit load while `unavailable` and prescribe opposite
  actions (reuse the live `unloaded` helper versus stop-and-replace), and A2-D12 carries no identity
  qualifier. Minor: §17.7.3 is cited as 1792–1808 though its heading is 1795, §4's summary of owner
  decision 11 says only a missing terminal triggers a kill, contradicting A2-D9 and §17.11.6 rule 4, and
  the §17.7.3 entry is labelled a clarification while adding a `loading` → `unloaded` edge to a table
  headed "Legal transitions". Editorial residues: A2-D2's provider-side `Supported:false` projection has
  no field to live in, A2-D5 omits the budget-exhausted route into `unavailable`, `detail.missing` has no
  defined encoding, and §17.7.4 row 1818 still over-describes the helper/config case.
- **Capture date:** 2026-11-07.

## 19. `AUDIT-howl-tts-capability-policy-c1-rev5-2026-11-07.md` (amendment revision 5 review)

- **What it is:** an independent adversarial review of revision 5 of
  `AMENDMENT-DRAFT-howl-tts-capability-policy-c1-17.19-2026-11-07.md` (SHA-256
  `8c728f2b…`, 302 lines) against the frozen architecture and the C-1 decision report. It answers the
  six focus areas the owner nominated — identity-specific recovery, budget-exhausted availability, the
  `loading` → `unloaded` edge, the `detail["missing"]` encoding, capability projection, and N/F
  traceability — and checks the thirteen owner decisions against O-1–O-7 and B12a/B12b. Read-only.
- **Status:** INVESTIGATIVE FINDING AND REVIEW ARTIFACT. PASS WITH FINDINGS — recommend adoption after
  one one-line scoping correction. Read-only.
- **Result:** all six focus areas pass. The identity collision is closed in both directions: A2-D7 now
  counts a *failed* identity as a held one, and A2-D12 limits in-place recovery to the same
  `provider_id`/`engine`/`model_identity` while routing any difference to stop-and-replace, so no
  differing-identity `load` reaches a helper that holds or has failed on an identity (the one in-place
  case, a cleanly unloaded helper, holds nothing and is what §17.4.4 requires). Budget-exhausted
  `unavailable` is fully representable (`Available=true`, `ProcessAlive=false`,
  `ServiceCapabilitiesKnown=false`) in four agreeing places, and the thirteen-state matrix closes
  without an unstated decision. `detail["missing"]` stays inside the frozen `map[string]string` and is
  unambiguous by construction (closed five-ID vocabulary, no commas, fixed order, 31 of 256 bytes).
  One defect remains, N9: A2-D12's “`synthesize` … never trigger a load” is unscoped and, read
  literally, collides with the only reachable load path in the adapter (§17.12.1 and §17.11.4 step 2,
  neither of which the supersession map touches). It is inherited from revision 4, not introduced by
  revision 5, and the precedence rule keeps it from superseding the frozen clause; the fix is one
  sentence scoping it to `State=unavailable`. One minor item (item 11's exhaustive “only” omits the
  identity condition) and four editorial notes (unenumerated `helper_start_failed` availability, the
  `accepting_work=true` clause on a state-keyed column, the missing O-1–O-7 sentence in §17.19.5, and
  the borrowed A2-D7 sequence without its budget exemption).
- **Capture date:** 2026-11-07.

## 20. `AUDIT-howl-tts-capability-policy-c1-rev6-2026-11-07.md` (amendment revision 6 review, with hash reconciliation)

- **What it is:** a targeted, independent, read-only review of revision 6 of
  `AMENDMENT-DRAFT-howl-tts-capability-policy-c1-17.19-2026-11-07.md` (SHA-256
  `dc5852ad…`, 320 lines), covering N9, N10, E1–E4, cross-instrument consistency, scope
  integrity, and the frozen-contract hash discrepancy raised at the revision 6 checkpoint.
- **Status:** INVESTIGATIVE FINDING AND REVIEW ARTIFACT. PASS — recommended adoption, no required
  fixes. Accepted for adoption on 2026-11-07; see §21–§22. Read-only.
- **Result:** both supplied hashes verified against the bytes. The frozen architecture hashes to
  `4ff6b72ea03eef98e86ac158add8eb05a7e575cea22384470991f406e5547393`; the competing `…f1f21818`
  variant appears in no file in the working tree and cannot be a digest of anything, since it shares
  its first 56 hex characters with the real value — it is a damaged copy, and the error is in the
  earlier checkpoint record, not in any artifact. N9 is closed in four places: A2-D12 governs recovery
  from `unavailable` only and affirms the frozen `Synthesize` ensure-handshake/load (§17.12.1,
  §17.11.4 step 2), A2-D5 and items 3 and 11 say the same, and the supersession map now lists
  §17.12.1 as confirmed rather than superseded, so the precedence rule and the decision agree without
  argument. N10's identity qualifier now appears in all eight recovery summaries. E1 closes the last
  unassigned `Available` cell (`starting` → `helper_start_failed` → `unavailable`, false, inside
  frozen row 1769's permitted outcomes); E2 attaches `accepting_work=true` to the resulting `unloaded`
  state with the "existing value" parenthetical; E3 restores the O-1–O-7 sentence in §17.19.5, A2-D11,
  and §5; E4 states twice that A2-D12 does not inherit A2-D7's budget exemption. Citation tails
  corrected (179–197, 192–197, 2354–2369, 2087–2088). One editorial nit (E5, singular "additive row"
  over a two-row table) and one integration observation (I1: the frozen adapter surface has no method
  that sends a `load`, so in-session retry depends on the update-gated reset or instance re-creation;
  not a contract defect). Scope verified: contract, decision report, and all five prior audits
  byte-identical; only the amendment changed after the revision 5 audit. Evidence limitation recorded:
  `docs/review/` is untracked and revisions overwrite one path, so no byte diff against revision 5 was
  possible.
- **Capture date:** 2026-11-07.

## 21. `AMENDMENT-DRAFT-howl-tts-capability-policy-c1-17.19-2026-11-07.md` (accepted Revision 6)

- **What it is:** the standalone Revision 6 amendment that adds §17.19, the owner-accepted Policy B
  capability floor, handshake semantic classification, service/provider status projection,
  `unavailable` recovery, and lifecycle clarifications. It preserves the pre-adoption contract and
  records its supersession map.
- **Status:** ACCEPTED AUTHORITY FOR ADOPTION — 2026-11-07. The file remains a standalone amendment;
  adoption did not rewrite it. It authorizes no implementation, Stage 3, code, caller wiring, native
  provider, asset, build, test, fixture, commit, or push.
- **Verified SHA-256:** `dc5852adbc70641dbe8f233851eaec44266330fce197b78a0ab2641e7f73f784`
  (320 lines, 65,274 bytes).
- **Independent review:** §20, Revision 6 audit, verified PASS — recommend adoption, no required fixes.
- **Result:** 12 decisions (A2-D1–A2-D12), a complete supersession/clarification map, and a resulting
  state/status ruleset. O-1–O-7 remain unresolved; B12a/B12b remain deferred.
- **Capture date:** 2026-11-07.

## 22. `ADOPTION-howl-tts-capability-policy-c1-17.19-2026-11-07.md` (A2 adoption record)

- **What it is:** the verifiable record of incorporating accepted Revision 6 into
  `DESIGN-howl-tts-helper-architecture-2026-11-06.md`. It records the verified pre- and post-adoption
  hashes, exact sections changed, decision-to-section mapping, checks and exit codes, deviations, and
  unresolved items.
- **Status:** ADOPTION RECORD. It records the governing architecture contract; it is not an
  implementation authorization. Stage 3 and all implementation work remain unauthorized.
- **Result:** §17.19 appended; §17.1 item 11, §17.7.2, §17.7.3, §17.7.4, §17.12, and §17.17
  updated additively; the review index qualified. The accepted amendment, C-1 decision, Revision 6
  audit, and all historical review artifacts remain unchanged.
- **Capture date:** 2026-11-07.

## 23. `INVESTIGATION-howl-tts-stage3-readiness-2026-11-07.md` (Stage 3 readiness gap analysis)

- **What it is:** the read-only contract-to-code gap analysis of the accepted Stage 2 helper/fake-helper
  implementation against the subsequently adopted governing contract, especially §17.19 (A2-D1–A2-D12).
  It covers the capability floor, handshake/error classification, thirteen-state lifecycle, status
  projection, recovery/identity/backoff/idle rules, process/file ownership, production reachability,
  and test evidence.
- **Status:** READ-ONLY INVESTIGATION — 2026-11-07. **Disposition: NOT READY** as a contract-compliant
  Stage 3 baseline. The real-process foundation is structurally sound, but current source contradicts
  adopted capability and lifecycle requirements. Stage 3 remains unauthorized; no implementation,
  test, fixture, wiring, build, run, or prior artifact was changed.
- **Verified SHA-256:** `a63f9f32f4d1bd047746a7bbbf5f509588e13ebd4272cffeb23f61fcfa891c6c`
  (586 lines, 51,248 bytes).
- **Key result:** five mandatory service capabilities must become two required plus three optional;
  semantic handshake invalidity must become terminal `protocol_incompatible`; status must expose
  `Loaded`/service-capability fields; ordinary unloaded synthesis must ensure-load; `unavailable`
  recovery needs identity-gated explicit load; restart budget/backoff, idle close, and correct
  shutdown classification are missing; no production caller or helper `TTSSynthesizer` adapter exists.
- **Capture date:** 2026-11-07.

## 24. `DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md` (Stage 3 corrective design and sequencing)

- **What it is:** the design-only corrective program for resolving the Stage 3 readiness gaps against
  adopted §17.19. It maps the investigation findings to eight dependency-ordered stages, preserves the
  Stage 1/Stage 2 process and output invariants, and defines scope, tests, acceptance evidence,
  rollback, affected interfaces, and same-day limits for each stage.
- **Status:** DESIGN PROPOSAL — corrected 2026-11-07 after adversarial review. Not implementation
  authorization. It selects C1 capability negotiation, C2 status projection, C3 lifecycle/load seam, C4
  exit/shutdown consistency, C5 budget/backoff/idle, C6 identity/recovery/replacement, C7 production
  adapter/explicit load, and C8 separately authorized contract closure and runtime verification. Stage 3
  remains unauthorized.
- **Verified SHA-256:** `0f920b247b8ba16bb78591916ae777d006e46190f5459cc8f738293527acc270`
  (703 lines, 59,634 bytes).
- **Key result:** C1 now separates direct `unload` from cancellation absence; C2 preserves both legal
  `unavailable` flag sets and single-source projection ownership; C3 owns `Loaded` transitions; C8 has
  binding Windows evidence rules; C7 counts real helper processes. C1–C6 remain inside the inert helper
  packages; C7 is the first production wiring and stays behind an internal default-off gate. D6 is
  settled, D7 is narrowed, D4 is folded, and D1 is split into D1a/D1b. D1a/D1b, D2, D5, B12a/B12b,
  O-1–O-7, and native provider work remain unresolved.
- **Capture date:** 2026-11-07.

## 25. `AUDIT-howl-tts-stage3-corrective-design-2026-11-07.md` (adversarial review of the Stage 3 corrective sequencing design)

- **What it is:** an independent, read-only design review of
  `DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md` against the adopted contract (§17.19/A2),
  covering C1–C2 capability semantics, C3–C6 dependency ordering, C7 production-adapter and
  character-voice compatibility, the C8 evidence boundary, and the D1–D7 decision inventory.
- **Status:** INVESTIGATIVE FINDING AND REVIEW ARTIFACT. **CONDITIONAL PASS** — accept the decomposition;
  two specifications must be corrected before C1 and C2 are authorized. Read-only.
- **Result:** proposal hash `71e64a93…` verified (641 lines / 44,056 bytes), §24 hash-matched, write set
  confirmed as exactly two documentation files with no non-docs file touched on 2026-11-07, and all six
  input hashes recomputed and correct. Inventory verified against source: the named functions exist, the
  five-mandatory set is real (`protocol.go:53`), `TTSSynthesizer` has six methods, `GetConfig` is in
  §17.12.1 and type-asserted at `pipeline.go:390`, and no production file imports the helper package.
  **F1 (blocking for C1):** C1 merges `unload` and `cancel` handling; A2-D9 forbids `unsupported_capability`
  in the cancellation path while A2-D1/A2-D7 assign it to a direct `unload`, so C1's blanket test row
  would freeze a contract violation, and its "handshake before gating" can start a helper for an
  operation the contract does not authorize starting for. **F2 (blocking for C2):** C2's `unavailable`
  rule covers only A2-D5 case (a) and forbids A2-D5 case (a) with its process-liveness sentence, while C5
  requires case (b) with `Available=true`; the rule and C2's own test list disagree. **F3:** `Loaded` has
  two owners and C2's acceptance cannot be met before C3. **F4:** C8's "direct, indirect, or explicitly
  unestablished" wording lets a fake-helper pass be recorded as indirect Windows evidence; an explicit
  four-part evidence rule is required. **D-table:** D6 is settled by frozen §17.7.2 (`stopping` + grace →
  `crashed`); D7 is substantially settled by A2-D12 (residual: ledger reset after in-place recovery);
  D3 is bounded by the frozen cap and §17.9 registry; D4 folds into D2/D5; D1 is genuinely open and must
  be split into seam-permissibility (D1a) and trigger/lifetime (D1b), because §17.12.1 has no `load` entry
  point at all.
- **Capture date:** 2026-11-07.

## 26. `CORRECTIVE-howl-tts-stage3-corrective-design-2026-11-07.md` (correction record for the Stage 3 corrective design)

- **What it is:** the documentation-only correction record for the adversarial review of
  `DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md`. It records the original and corrected
  proposal hashes, the audit hash, dispositions for F1/F2/F3/F4/F8, the D1–D7 reclassification, exact
  changed sections, and preservation of the governing contract and readiness investigation.
- **Status:** CORRECTION RECORD — 2026-11-07. Documentation only; authorizes no implementation. C1 and C2
  remain unauthorized.
- **Verified SHA-256:** `2eeaaf14ec9b80d0807453ca8029eac3c78c383de3b674cd8f193601b9565e57`
  (134 lines, 11,452 bytes).
- **Key result:** direct `unload` absence returns `unsupported_capability` with no frame, cancellation
  absence skips the frame with no capability error, and unknown snapshots never spawn for an optional
  operation; both `unavailable` cases are preserved; `Loaded` has one writer per stage; Windows-specific
  evidence cannot close indirectly; the cap counts actual helper processes. D6 is settled, D7 retains
  only the stable-readiness reset question, D4 is folded, D3 remains a narrow error-surface gap, and D1
  is split into D1a/D1b.
- **Capture date:** 2026-11-07.

## 27. `AUDIT-howl-tts-stage3-corrective-design-r2-2026-11-07.md` (second-round adversarial review of the corrected Stage 3 design)

- **What it is:** an independent, read-only re-review of the corrected Stage 3 corrective design and its
  correction record against the adopted contract, covering every A2 decision, the two legal
  `unavailable` states, required-versus-optional capability absence, direct unload versus cancellation,
  C1–C8 sequencing and authorization boundaries, the D1–D7/O-1–O-7 ledger, and a fourteen-case
  failure-path audit.
- **Status:** INVESTIGATIVE FINDING AND REVIEW ARTIFACT. **PASS WITH BLOCKERS** — the correction
  strategy and eight-stage decomposition are accepted as the design basis, but C1 is not authorizable
  until two defects are repaired and C3 is not authorizable until a sequencing defect is repaired.
  Read-only; authorizes no stage.
- **Verified SHAs:** contract `ec3c789e…`, corrected proposal `0f920b24…` (703 lines / 59,634 bytes),
  correction record `2eeaaf14…` (134 lines / 11,452 bytes). All three match the supplied values; the
  contract, readiness investigation `a63f9f32…`, and prior audit `71a17369…` (25,962 bytes) were
  recomputed and are unchanged. Only documentation files carry a 2026-11-07 mtime after the prior audit.
- **Result:** F1's core, F2, F3, F4, and F8 are genuinely corrected; F2's two-case table matches the
  contract's additive §17.7.4 rows field-for-field. **B1 (blocking, C1):** the cancellation exemption is
  written at `Cancel` method scope although A2-D9 grants it only to "the supervisor's cancellation path",
  so a direct call would return success without sending a frame, against A2-D1/A2-D10. **B2 (blocking,
  C1):** unknown-snapshot direct `unload` returns `provider_not_loaded`, a frozen wire code meaning
  "Synthesis before successful load", while the same document forbids repurposing codes for D3; the
  no-spawn rule is right but belongs to §17.7.2, not §17.19, and the case is an unregistered decision.
  **M1:** C2's `Available` false set omits `absent`, `starting`, `handshaking`, `stopping`, `crashed`,
  and `backoff`. **M2:** C1's tests omit duplicate-extension and malformed-descriptor A2-D3 categories,
  where current code still returns `protocol_violation`. **M3:** D1a is scheduled before C7 while C3
  creates the loader seam it governs. **M4:** F6 was never addressed — C5 still carries the
  question-form idle-reset rule, which needs decision D9. D7's residual is settled by §17.11.5 and needs
  no owner decision.
- **Capture date:** 2026-11-07.

## 28. `AUDIT-howl-tts-stage3-corrective-design-r3-2026-11-07.md` (third-round audit of the Stage 3 corrective design)

- **What it is:** a read-only verification that the second correction pass documented the round-2
  findings faithfully in the proposal text (not merely in its record), that no earlier correction was
  narrowed or lost, and that the stage and decision structure remains internally consistent. It
  re-derives every A2/§17 citation from the governing contract rather than accepting the record's
  claims.
- **Status:** REVIEW ARTIFACT — **PASS**. The correction pass is accepted as documentation. One new
  major finding (M5) and five minors are recorded. Read-only; authorizes no stage.
- **Subject hashes:** proposal `1fb449ba…` (710 lines / 65,213 bytes, 10:15:04), record `ccb28d57…`
  (175 lines / 13,087 bytes, 10:16:14), contract unchanged at `ec3c789e…` (2,586 lines / 182,904
  bytes). The record's §1 table matches disk, including the pre-pass proposal `0f920b24…`.
- **Result:** B1, B2, M1, M2, M3, and D7 are **closed**; M4's documentation is correct. B1's
  exemption is now path-scoped with a two-sided regression test; B2's unknown branch asserts no
  result code and is registered as D10; the C2 false set matches A2-D5's seven states exactly; C1
  names duplicate-extension and invalid-descriptor categories and correctly keeps `invalid_frame`
  for duplicate JSON keys; D1a is gated before C3 in all eight places that mention it; D7's closure
  by §17.11.5 is independently verified. F1–F8 all survive. **M5 (major, new):** D9's second clause
  asks whether a status request may respawn an idle-closed helper, but §17.7.2's lazy-start row
  makes a status request from `absent` a permitted start, A2-D12 restricts status only from
  `unavailable` recovery, and C7's own `ProviderStatus` contract already answers yes — an owner
  answer of "no" would put C5 in conflict with the contract. Repair D9's scope before it is
  dispositioned. Minors: "exactly at 3" overstates §17.11.5's "maximum 3"; no structural no-frame
  test for a direct `Cancel` with an unknown snapshot; §9.1 omits C1's D10 acceptance gate;
  "unsupported-operation frame" is an undefined term; the path-scoped operation's call site is not
  named in C1's file list.
- **Index note:** §24, §26, and §27 above cite the superseded `0f920b24…` / `2eeaaf14…` hashes and
  describe the prior revision. They were deliberately left unedited by the correction pass. An index
  re-sync is a separate maintenance task and has not been performed.
- **Capture date:** 2026-11-07.

## Provenance caveat shared by all documents

All artifacts here were produced in agent sessions whose tool results are not a build attestation. Addendum A1
grades its own evidence (§0) and caps every audio capability at source-verified + inferred binary
linkage; no audio path was load-tested or run end-to-end. Treat the artifact's stated evidence grades
as binding on how its findings may be cited.
