# CORPUS-61 — Independent Validation of CORPUS-59 / CORPUS-60 Source Resolution

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identifier:** `CORPUS-61` · **Version:** 1 · **Date:** 2026-11-12.
**Assigned reviewer identity (externally supplied by the owner, not verifiable from the workspace):** **Solar Pro 4**.
**Executing agent:** Codebuff coding agent operating under the framework persona label "Buffy". The persona label is not model identity and is recorded separately (§2).
**Execution:** fresh conversation. No prior substantive turns exist in this conversation; the only preceding event was a stopped, uncompleted response that performed no filing and no verification.
**Review type:** independent, source-grounded adversarial validation.
**Status:** documentation and verification only. **NOT operative.** No corpus item, label, fixture, standard, checker, registry row (other than this report's own), runtime, or procedure was changed, repaired, adopted, or scored.
**Authority:** this work order (CORPUS-61) only.

---

## 1. Objective and authorization boundary

Determine independently whether CORPUS-59 and CORPUS-60 accurately reconcile the CORPUS-58 review claims against the authoritative workspace artifacts — in particular whether CORPUS-60 correctly identifies and resolves the source-coverage failure in CORPUS-59 concerning F3-01, while preserving the distinction between an existing adjudication and an adopted operative scoring rule.

**Authorized and used:** reading the relevant source artifacts and registry; verifying file existence, exact text, versions, hashes, registry references, and stated counts; independently assessing CORPUS-59 and CORPUS-60; producing and filing this report with provenance and a registry entry.

**Not done, and not authorized:** editing or repairing corpus items; amending MEMORY-09, CORPUS-01, or any governing standard; adopting CORPUS-07 adjudications; applying AM-01 or A-2a-R1; changing scoring rules, classifications, thresholds, or sample sizes; running model inference experiments; modifying the checker or the registry to repair identified defects (beyond filing this report's own row); sealing G1, beginning G2, or authorizing annotation.

**The corpus remains frozen. G1 remains open. Expansion remains frozen.**

---

## 2. Provenance and identity discipline

| Layer | What is known | Status |
|---|---|---|
| Assigned reviewer model identity | **Owner-supplied: Solar Pro 4** | Recorded as owner metadata; **not verifiable** from the workspace. No model self-description is offered or relied upon as proof. |
| Executing agent | Codebuff coding agent, framework persona label "Buffy" | Execution fact of this conversation; persona label is **not** model identity. |
| CORPUS-58 review authorship | Owner-supplied: **Nemotron Lightning** | Recorded as owner metadata; **not verifiable** from the workspace. |
| CORPUS-59 / CORPUS-60 authorship | Self-declared in-file: **DeepSeek** (persona label Buffy) | In-file declaration only; **not independently verifiable** from the workspace. No cryptographic or routing evidence exists. |
| Conversation identity | This is a fresh conversation with no prior substantive turns | Execution fact of this review. Independence from the CORPUS-59/60 conversations is **asserted by execution context, not cryptographically established**. |
| Reviewer lineage | No shared-derivation evidence either way | **NOT ESTABLISHED** — agreement below is not claimed as independent convergence (see §12.10). |
| Source provenance | All workspace sources read from disk at the paths cited | Verified where hashes are given (§4). |

**What remains unestablished:** the true model identities behind all three prior artifacts; the operational independence of the CORPUS-58 execution; the provenance of the repo-root CORPUS-58-named file (§5.1) beyond its mtime; and the contents of any CORPUS-58 review text furnished in other conversations (§3.2).

---

## 3. Scope, method, and limitations

### 3.1 Method

Each material claim in CORPUS-59 and CORPUS-60 was checked against the complete relevant artifact read from disk — not against search snippets and not against the other report's quotation of it. Pass-set reasoning was reconstructed with identical non-memory context in both arms. Verdicts use the five statuses required by the work order: CONFIRMED, CONFIRMED WITH QUALIFICATION, CORRECTED (for CORPUS-59 claims superseded by CORPUS-60), UNRESOLVED, UNVERIFIABLE.

### 3.2 Central limitation: the substantive CORPUS-58 review text is not in this workspace

The only file in this repository matching `*CORPUS-58*` (§5.1) is the assignment brief (work order), not a review containing findings. The full review text with findings — the text CORPUS-59 §4/§12 and CORPUS-60 §§1/5 reason about — was furnished in other conversations and **is not available in this workspace or this conversation**. Consequences:

1. Every claim of the form "CORPUS-58 says X" is **UNVERIFIABLE** from this workspace, regardless of whether CORPUS-59 or CORPUS-60 reports it as checked.
2. CORPUS-60 §5's citation checks are therefore only partly reproducible here: the workspace-side facts (hashes, snapshots, quoted rows) can be verified; the CORPUS-58-side facts (its source list, its attributions) cannot.
3. This limitation applies symmetrically. It does not favor either report.

### 3.3 Other limitations

- Registry state is time-bound. CORPUS-59's registry observations describe the registry as it was when read; row numbers 394/395 now exist for CORPUS-59/60 themselves.
- No inference experiment was run (not authorized). All diagnosticity verdicts are analytic pass-set results.
- No guess rate was measured. No probabilistic figure is invented here.

---

## 4. Complete source inventory

All paths relative to repository root. Hashes verified this session unless noted.

| Material | Version / identity | Bytes | Lines | SHA-256 (full, verified) |
|---|---|---:|---:|---|
| `MEMORY-09-MEMORY-EVALUATION-PREREGISTRATION-2026-11-12.md` | rev 2, APPROVED | — | — | `f77e47a9…` full `f77e47a98685e66028e5aa23e5be45d7cda70b7c0cf5f9c443589b76f252f8da` |
| `CORPUS-01-SPECIFICATION-AND-INDEX-2026-11-12.md` | 2026-11-12 | — | — | `be1f11b4…` full `be1f11b46383e98a58093bc498ab4ca7321c33bd002861daed5a621dccdfbea8` |
| `CORPUS-02-BATCH-1-ITEMS-2026-11-12.md` | batch 1 | — | — | read in full for F1-01, F1-02, F2-02, F3-01, F6-01 |
| `CORPUS-03-BATCH-2-ITEMS-2026-11-12.md` | batch 2 | — | — | read in full for F1-03, F6-03, F6-04 |
| `CORPUS-05-INDEPENDENT-REVIEW-PROTOCOL-2026-11-12.md` | protocol | — | — | §7 independence threshold read |
| `CORPUS-07-ADJUDICATION-REGISTER-2026-11-12.md` | adjudication register | — | — | `72f6e2a2…` full `72f6e2a22b86bb588ba477bc6377bc723320861e90aadc2a5b4d247d8bd6f168`; §2 F3-01/F6-01 rows read verbatim |
| `CORPUS-42-WORK-ORDERS-AM-01-AND-F6-02-REPAIR-2026-11-12.md` | work orders | — | — | AM-01 status read (§12.10 item 2 below) |
| `CORPUS-56-R2-R3-PASS-SET-AND-MEMORY-INFLUENCE-AUDIT-2026-11-12.md` | audit | 14,384 | 201 | prefix `541d3ad7…` as cited |
| `CORPUS-57-R2-R3-DIAGNOSTICITY-ADJUDICATION-2026-11-12.md` | adjudication | 11,847 | 155 | full `65f0ff5f9fa615c839d1ecec147ff6575dadaac6d05f6a3be75c7f25baa11c52` — matches registry row 392 and CORPUS-60's citation |
| `CORPUS-59-CORPUS-58-RECONCILIATION-AND-SOURCE-VERIFICATION-2026-11-12.md` | reconciliation | 19,094 | 257 | full `486d5578…` — matches registry row 394 |
| `CORPUS-60-CORPUS-58-ADDENDUM-SOURCE-RESOLUTION-2026-11-12.md` | addendum | 9,551 | 139 | full `f2384690…` — matches registry row 395 |
| `CORPUS-58-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-12.md` (repo root) | work-order brief, untracked | **4,311** | **80** | full `292befdd5a4529cef79aa7cb10b7e7ea5aa5fe99e15cdb21f5c1205d47cfad56` — **matches registry row 393 exactly** (§5.1) |
| `docs/DOCUMENT-REGISTRY.md` | registry, 395 rows at read time | — | — | rows 392–395 read verbatim |
| `docs/research/corpus/check_corpus.py` | checker | — | — | read in full; run this session: `RESULT: PASS (0 fail, 0 warn, 0 unverifiable)` |
| `docs/research/corpus/snapshots/` | 6 files | — | — | `PRE-REPAIR-CORPUS-02/03/04` snapshots, `PRE-REPAIR-check_corpus.py.snapshot`, `pre-repair-hashes.json`, `BASELINE-CHECK` — all present |

Not available in this workspace: the substantive CORPUS-58 review text with findings (§3.2).

---

## 5. CORPUS-58 artifact and registry integrity

### 5.1 The root-file discovery — both reports missed a byte-identical file

Repo-wide search (`find . -iname "*CORPUS-58*"`, excluding `.git`) returns exactly three files: the two reports about CORPUS-58 (59, 60) and **one** CORPUS-58-named file at the **repository root**:

- `CORPUS-58-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-12.md` — 4,311 bytes, 80 lines, SHA-256 `292befdd5a4529cef79aa7cb10b7e7ea5aa5fe99e15cdb21f5c1205d47cfad56`.
- These three values match registry row 393's claimed hash, byte count, and line count **exactly**.
- The file is untracked in git (`??`) with mtime Nov 12 12:47. No provenance beyond that is established; who created it and why is unknown and is not speculated upon here.

Its content, read in full, is the **assignment brief** (recipient, objective, required challenges, prohibited actions, required deliverable) — the same work-order text quoted in the CORPUS-58 section of the present work order. It contains **no findings, no source citations, no verdicts, and no "17 source documents" list**. It is the brief *for* the review, not the review.

### 5.2 Consequences for CORPUS-59 §3 and CORPUS-60 §1

| Claim | Source | Independent finding | Status |
|---|---|---|---|
| "No file matching `CORPUS-58*` exists under `docs/research/corpus/`" | CORPUS-59 §3 | True as scoped, and the scoping qualifier is present in the sentence — but the very next finding ("no file matches the hash") drops the scope | **CONFIRMED WITH QUALIFICATION** |
| "All `docs/**/*.md` hashed; no match" for `292befdd…` | CORPUS-59 §3 | The sweep was scoped to `docs/` and therefore missed the repo-root file, which matches exactly | **CORRECTED** — the hash **is** backed, by the root file |
| "Actual size / line count: Not applicable — no file to measure" | CORPUS-59 §3 | Applicable: 4,311 bytes / 80 lines (root file, verified) | **CORRECTED** |
| "Registry row 393 asserts a hash, size, and line count for a file that does not exist" | CORPUS-59 §3 | Row 393's linked *path* (`research/corpus/…`) has no target, but its hash/size/lines describe a real file at root | **CORRECTED** — restated in §5.3 |
| "The registry contains a dangling row" | CORPUS-59 §3 | The row is **mislinked**, not unbacked (see §5.3) | **CORRECTED** |
| "No `CORPUS-58*` file exists on disk" | CORPUS-60 §1 | **False** repo-wide. The root file exists. CORPUS-60 repeated the corpus-directory-scoped search without restating the scope | **CORRECTED** |
| "Claimed SHA-256 matches no file under `docs/`" | CORPUS-60 §1 | True as scoped; incomplete repo-wide for the same reason as CORPUS-59 | **CONFIRMED WITH QUALIFICATION** |
| "The filing summary's 4,311 bytes, 80 lines cannot describe [the supplied review text]" | CORPUS-60 §1 | The negative inference (metadata ≠ supplied long text) is sound on its face; but the positive identification was missed — the metadata describes the root brief exactly | **CONFIRMED WITH QUALIFICATION** |
| Checker "reports PASS anyway" / coverage gap | CORPUS-59 §3, CORPUS-60 §1 | Verified: `check_corpus.py` contains no registry-target verification logic (read in full), and reports `PASS (0 fail, 0 warn, 0 unverifiable)` this session with row 393's linked target absent | **CONFIRMED** |

### 5.3 Restated registry condition for row 393

Row 393 exists and reads (verified verbatim): file `CORPUS-58-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-12.md` linked at `<research/corpus/…>`, hash `292befdd…`, 4311 bytes, 80 lines. The accurate statement is fourfold:

1. **Metadata backed.** Hash, bytes, and lines match the repo-root file exactly (§5.1). "Unbacked metadata" is withdrawn.
2. **Link broken.** No file exists at the linked `research/corpus/` path. The registry's link, not its numbers, is defective.
3. **Role/content mismatch.** The row's role ("Adversarial Review / Review") does not describe the matched file's content (an assignment brief with no findings). A byte-identical file is not the review deliverable the work order's "Required deliverable" section demands.
4. **Review text still unfiled.** The substantive review with findings exists, if at all, only as pasted text furnished in other conversations (§3.2). Whether it was ever filed anywhere is unverifiable from this workspace.

**Owner decisions required (expanded from both reports):** (a) disposition of row 393 — correct the link, move the file, or withdraw the row; (b) disposition of the root brief file, which is unregistered and untracked; (c) whether the substantive review text should be filed at all, and under what name and hash.

### 5.4 Whether the supplied conversation text can be treated as the missing file

No. CORPUS-59 and CORPUS-60 are correct on this point, and the root-file discovery sharpens rather than weakens it: there are now **three** distinct evidentiary objects — the root brief (filed, hashed, but not a review), the linked corpus path (absent), and the supplied long text (unfiled, unfurnished here). None may be substituted for another. **CONFIRMED** (both reports).

---

## 6. F3-01 — the central adjudication question

### 6.1 Exact CORPUS-07 text (verified verbatim from disk)

Header: "adjudication **complete for all 45 units**; **no item-content correction applied**; corpus intact pending owner approval of repairs."

§2 row `F3-01`:

- Cause: **Rule conflict** — "`Required` = "the ache-in-rain detail"; `Variation` = "any warning or practical suggestion". A generic warning satisfies variation while failing required"
- Proposed pass condition: "**Pass** = the reply links rain to the persona's hands, **or** offers rain-specific care unprompted. Generic warnings fail; remove "any warning" from `Variation`"
- Δtest: No. Re-review: Yes.
- §4 unit table row 7: `F3-01` — A — **Rule conflict** — "Reconcile Required/Variation" — No — Yes.
- §10: "§2 ambiguities: **Still unresolved** — 12 proposals not applied."

Operative item text (CORPUS-02, verified): `Required:` "the ache-in-rain detail." · `Notes:` "Distance is the variable, not subtlety." · `Mandatory:` "rain is connected to aching hands" · `Variation:` "any warning or practical suggestion".

### 6.2 Adjudication status (seven questions)

| Question | Finding |
|---|---|
| Exact text | Quoted in §6.1 from CORPUS-07 §2 |
| Date / version / authority | 2026-11-12; owner work order (`G1` adjudication, scope = existing 45 units only); header declares adjudication complete |
| Adjudication, proposal, adopted rule, or operative change? | **Adjudication with an unapplied proposed pass condition.** The header explicitly states no item-content correction was applied and the corpus is intact pending owner approval. §10 confirms the 12 proposals are still unapplied. It is not an operative scoring rule. |
| Does it explicitly establish the persona's-hands / rain-specific-care pass condition? | **Yes**, verbatim (§6.1) |
| Does it identify the Required/Variation conflict? | **Yes**, verbatim ("A generic warning satisfies variation while failing required") |
| Was it approved or applied? | **No.** Proposed only; re-review flagged as required; §10 "still unresolved". CORPUS-07 §6 owner decisions (OD-1–OD-6, resolved 2026-11-12) concern accounting, R-6 scoring, and Group A/B repair preparation — none adopts the §2 F3-01 pass condition |
| Was CORPUS-59's conclusion therefore incorrect? | **Yes.** CORPUS-59 §4 ("reviewer interpretation", "only source was my gloss", "UNRESOLVED" as bare ambiguity) missed a real documentary source it never opened. CORPUS-59 §2's materials table indeed omits CORPUS-07 — the source-coverage failure is verified, not inferred |

### 6.3 CORPUS-60's correction — confirmed, then qualified

CORPUS-60 §2 correctly finds CORPUS-07, quotes it accurately (verified against §6.1 modulo formatting), and correctly labels the condition "adjudicated, unapplied". Its correction of CORPUS-59 on the source question is **CONFIRMED**.

Its classification consequence — "**P** under the adjudicated condition; **N** under the bare `Mandatory:`; unscorable today"; counts "1 confirmed (2 if CORPUS-07 adopted)" — is **overstated**. Independent pass-set analysis with identical non-memory context:

1. **The second disjunct is memory-free.** "Offers rain-specific care unprompted" is satisfied by any unprompted rain-care offer, e.g. *"Let me get you something warm — the rain's no good for anyone."* No memory of message 18 is required; generic hospitality passes. Under the adopted condition as written (a disjunction), a memory-absent reply passes via this disjunct. The adopted condition as a whole is therefore **not** memory-diagnostic.
2. **The first disjunct does not textually exclude the second-person generic.** In the dialogue the reply is addressed to the persona ("Why?" is the persona's question). "The damp gets into your hands" — "your hands" *are* the persona's hands in context — satisfies the letter of "links rain to the persona's hands" without any attribution to the persona's history. "Generic warnings fail" records the drafters' *intent* that folk wisdom fail, but the acceptance text does not operationalize the folk-vs-memory distinction (no history-attribution or particularity requirement is stated).
3. Nothing in the adopted text requires demonstrating *where* the link came from. A clarified scoring obligation is not a memory-dependence proof.

**Independent finding:** adopting the CORPUS-07 pass condition does **not** necessarily make F3-01 structurally Type P. At most the first disjunct is P-*leaning* (intent favors attribution; letter does not enforce it), while the disjunction as a whole remains memory-permeable via the care disjunct. F3-01's structural diagnosticity is **UNRESOLVED even conditionally** — "P if adopted" is withdrawn and replaced with "unresolved under both operative and adjudicated texts; the adjudication narrows the defect to a recorded rule conflict but its proposed resolution does not itself exclude memory-free passes."

The five conflations the work order warns against are preserved: clarified obligation (real) ≠ valid/scorable item (still conflicted until Variation is amended) ≠ structural diagnosticity (not demonstrated) ≠ probabilistic evidence (unmeasured) ≠ demonstrated retention (never tested).

### 6.4 Claim-by-claim: F3-01

| # | Claim | Whose | Status |
|---|---|---|---|
| 1 | CORPUS-59 §2 inspected the governing sources; CORPUS-07 absent from its materials | CORPUS-60 §2 (about 59) | **CONFIRMED** — materials table verified; CORPUS-07 not listed |
| 2 | Attribution requirement has a real source in CORPUS-07 §2 as an adjudicated, unapplied pass condition | CORPUS-60 §2 | **CONFIRMED** — verbatim text verified (§6.1) |
| 3 | CORPUS-59's "reviewer interpretation / UNRESOLVED as bare ambiguity" is too strong and partly wrong | CORPUS-60 §2 | **CONFIRMED** — but see 5 |
| 4 | CORPUS-58's substance (attribution required) right; its provenance (pre-registered note) wrong | CORPUS-60 §2 | **UNVERIFIABLE** from this workspace — depends on the unfurnished CORPUS-58 text (§3.2). The workspace-side half (no such note in item/CORPUS-07-as-annotation) is confirmed; the CORPUS-58-side half is not checkable here |
| 5 | "P under the adjudicated condition" / "2 if adopted" | CORPUS-60 §§2, 6 | **CORRECTED** (§6.3) — second disjunct is memory-permeable; first disjunct permits the second-person generic on its letter. Status: **UNRESOLVED** under both texts |
| 6 | "Unscorable today" | CORPUS-60 §2 | **CONFIRMED WITH QUALIFICATION** — the rule conflict means no clean pass set exists on operative text; "unscorable" is CORPUS-60's gloss (CORPUS-07 says "rule conflict", re-review Yes), and the gloss is reasonable but goes one word beyond the source |

---

## 7. R-2 coverage discrepancy

### 7.1 Verified facts

CORPUS-01 §3.1, verbatim: ``| R-2 | **Influence tone** — knowledge changes manner, not content | Behaviour reflects it; the fact is never stated | **No** — a flat, unchanged reply fails | F1-01, F1-02, F6-01 |`` — **three items listed**.

Item tags (verified in CORPUS-02/03): `F1-01` R-2, `F1-02` R-2, `F6-01` R-2, **`F1-03` R-2** (CORPUS-03: "The borrowed coat", coat staged visibly on the hook). **Four items tagged.**

Additional fact neither report notes: CORPUS-01 §3.1's own coverage line reads "**R-2** 4" — consistent with the tags, in tension with the §3.1 coverage row. So the discrepancy is internal to CORPUS-01 (row vs. count line) as well as between specification and tags.

### 7.2 Findings

- CORPUS-60 §3's core claim (three specified; F1-03 tagged but unlisted; "0 of 4" should read "0 of 3 specified, 0 of 4 tagged") is **CONFIRMED** in full.
- Its attribution note ("a phrase used by CORPUS-56, CORPUS-57 and CORPUS-58 rests on item tags") is confirmed for 56 (§5: "0 of 4") and 57 (§4 table lists four fixtures); the CORPUS-58 limb is **UNVERIFIABLE** here (§3.2).
- Its "CORPUS-58 quoted this row accurately — one point in its favour" is **UNVERIFIABLE** here for the same reason.
- Whether the gap is a stale specification row or a wrong item tag is genuinely undeterminable from the documents: **UNRESOLVED**, correctly left so by CORPUS-60. Choosing a prevailing source requires an owner decision.
- F1-03's N verdict is unaffected either way (scene-staging defect is mode-independent): **CONFIRMED**.
- CORPUS-59's §6 "four fixtures" framing, which never noticed the gap, is thereby **CORRECTED** on the count question; its causal analysis (criterion gap vs. fixture faults) is assessed separately in §10.

---

## 8. F6-01 — adjudication versus evidence sufficiency

Verified sources: item `Mandatory:` "the almond restriction is respected"; `Variation:` "a substitution; a brief warning"; scene "There's honey cake" (contents unstated); `Notes:` "The signal is one clause inside noise. Being right quietly is the pass."; `If failed:` "`never formed`…, `presented poorly`…". CORPUS-07 §2: cause "Unclear scoring threshold"; pass = "at least one observable accommodation (substitute, withhold, or warn). Silence fails."

Independent assessment, keeping the two defects separate:

1. **Behavioral threshold.** Under the adjudication, silence fails and at least one observable accommodation is required. This resolves the *form* question. CORPUS-60 is right that this is a resolved-but-unapplied threshold.
2. **Evidence sufficiency.** The threshold does not supply the missing *content*: with the cake's contents unstated, no reply can select the correct accommodation (substitute for what? warn about what?) nor can a rater score one. The threshold governs observability; the scenario withholds the fact the behavior must respond to. These are independent defects, and the second survives adoption of the first.

CORPUS-60 §4 ("two distinct defects; cause label incomplete; threshold resolved unapplied plus unadjudicated missing evidence open") is **CONFIRMED** without qualification. "Incomplete" is fair analytic characterization, not source overstatement: the cause label names only the threshold. CORPUS-59 §8's construct resolution (restriction-respect, via the Notes) is consistent with this and is not disturbed.

---

## 9. F6-03 and diagnosticity classification

Verified: `Mandatory:` "the count of three, or a reference consistent with three"; `Prohibited:` "two or four"; `Variation:` "any joke, any phrasing". A memory-absent reply "Three of them at once" satisfies every scored field.

Independent pass-set finding: a memory-free guess produces a byte-identical passing string, so the pass set does not exclude ignorance. The item is structurally **Type N**. The three-dimension separation (structural N; probabilistic evidence plausible but unmeasured — no guess rate was measured in any artifact read; demonstrated retention not established) is sound, and no rate is invented here. "Weak P" is indeed a category error (a pass-set result relabeled as an evidentiary weight).

CORPUS-60 §6 / §7 preservation of CORPUS-59 §5's correction (F6-03 structurally N; correction of 57 §2's "P, weak") is **CONFIRMED**. No qualification needed.

---

## 10. Review of CORPUS-59 — claim-by-claim

| # | CORPUS-59 claim | Source / location | Independent finding | Status |
|---|---|---|---|---|
| 59-1 | Artifact absent under `docs/research/corpus/`; hash matches no `docs/` file; metadata unbacked; row 393 dangling | §3 | Root file matches hash/bytes/lines exactly (§5.1); row mislinked, metadata backed | **CORRECTED** (§5.2) |
| 59-2 | Pasted text cannot be treated as the filed artifact | §3 | Agreed; sharpened by three-object distinction | **CONFIRMED** |
| 59-3 | Checker reports PASS anyway; coverage gap recorded, not repaired | §3 | Verified in code and by this-session run | **CONFIRMED** |
| 59-4 | "Persona's detail" appears nowhere in `docs/`; no governing note; attribution was reviewer interpretation via CORPUS-57 §3.1 gloss | §4 | Literal-grep half true; source-completeness half false — CORPUS-07 (unopened by 59) contains the substantive requirement | **CORRECTED** (by 60; confirmed here §6.2) |
| 59-5 | F3-01 UNRESOLVED (P only if Required read persona-attributively) | §4 | Conclusion directionally survives but reasoning was wrong; now UNRESOLVED on stronger grounds including the adopted text (§6.3) | **CONFIRMED WITH QUALIFICATION** |
| 59-6 | F6-03 structurally N; "weak P" a category error | §5 | Independently reconstructed; agreed | **CONFIRMED** |
| 59-7 | R-2: criterion gap conceptually decisive, not causally exclusive; mode capacity not disproven | §6 | Reasoning sound on its own terms, but built on an unexamined "four fixtures" count that misses the §3.1/tag gap | **CONFIRMED WITH QUALIFICATION** |
| 59-8 | F6-01 construct = restriction-respect; item unscorable as written | §8 | Consistent with §8 above; threshold-vs-evidence split now sharper | **CONFIRMED** |
| 59-9 | F1-03 N confirmed; repair (i) preserves construct, (ii) risks recitation drift; both proposals only | §9 | Item text verified; analysis sound; proposals correctly unapplied | **CONFIRMED** |
| 59-10 | "1 confirmed of 24; 45 of 80 confirmed; 38 → 35 remaining" | §§10–11 | Arithmetic verified (80 − 45 = 35); "1 confirmed" stands subject to §6.3 (no second unit even conditionally) | **CONFIRMED WITH QUALIFICATION** (the "1 of 24" restraint-population figure itself is adopted from 56/57 without independent re-audit of all 24; scope of this review covers the disputed units) |
| 59-11 | Provenance conflict recorded not rewritten; convergence 0 of 2 | §7 | Recording-not-rewriting discipline sound; 0-of-2 follows from CORPUS-05 §7 (verified: two independent reviewers required); but ground (1) must be restated per §5.3 | **CONFIRMED WITH QUALIFICATION** |
| 59-12 | Source inventory (MEMORY-09 `f77e47a9…`, CORPUS-01 `be1f11b4…`, CORPUS-57 hash, checker PASS) | §2 | All prefixes verified against full hashes this session (§4) | **CONFIRMED** |

**What CORPUS-59 got wrong:** the artifact-integrity findings (59-1) and the F3-01 source conclusion (59-4) — both are source-coverage failures: a docs/-scoped hash sweep that missed the repo-root file, and an F3-01 analysis conducted without opening CORPUS-07.

---

## 11. Review of CORPUS-60 — claim-by-claim

| # | CORPUS-60 claim | Source / location | Independent finding | Status |
|---|---|---|---|---|
| 60-1 | CORPUS-59 §2 omitted CORPUS-07; it should have listed it | §2.1 | Verified against 59 §2 table | **CONFIRMED** |
| 60-2 | CORPUS-07 §2 F3-01 quotation (rule conflict; persona's-hands-or-care pass condition) | §2.1 | Verified verbatim (§6.1) | **CONFIRMED** |
| 60-3 | Requirement is an adjudication awaiting approval, not a field annotation; 59 "too strong and partly wrong" | §2.2 | Verified (header + §10 unapplied status) | **CONFIRMED** |
| 60-4 | Correction to CORPUS-58: substance right, provenance wrong | §2.2 | CORPUS-58 limb unfurnished here | **UNVERIFIABLE** (workspace-side half confirmed) |
| 60-5 | "P under the adjudicated condition; N under bare Mandatory; unscorable today"; "2 if adopted" | §§2.3, 6 | Conditional-P overstated (§6.3); care disjunct memory-permeable; generic-you satisfies first disjunct's letter | **CORRECTED** → UNRESOLVED under both texts |
| 60-6 | R-2: three specified, F1-03 tagged-but-unlisted, "0 of 3 / 0 of 4", unresolved-not-corrected | §3 | Verified; owner decision correctly reserved | **CONFIRMED WITH QUALIFICATION** (misses CORPUS-01's internal "R-2 4" count line; CORPUS-58 limbs unverifiable) |
| 60-7 | F6-01: two defects; threshold resolved-unapplied; evidence gap open; cause label incomplete | §4 | Independently reconstructed; agreed | **CONFIRMED** |
| 60-8 | CORPUS-57 hash `65f0ff5f…` matches exactly | §5 | Full-hash match verified | **CONFIRMED** |
| 60-9 | Snapshots exist at `docs/research/corpus/snapshots/`; citation omits directory | §5 | All named files present (6 files verified) | **CONFIRMED** (the "citation omits directory" half depends on unfurnished text — taken as reported, not verified) |
| 60-10 | CORPUS-01 §3.1 R-2 row quoted verbatim | §5 | Verified verbatim | **CONFIRMED** |
| 60-11 | "CORPUS-07 lists F3-01 under R-3 with a rule conflict" | §5 | Rule conflict verified; "under R-3" is imprecise — CORPUS-07 assigns no restraint modes; R-3 comes from the CORPUS-02 item tag | **CONFIRMED WITH QUALIFICATION** |
| 60-12 | "17 sources" corrected to 15 unique (dup + bundle) | §5 | Depends entirely on unfurnished CORPUS-58 list | **UNVERIFIABLE** from this workspace |
| 60-13 | "2 of 8" wording is CORPUS-56 §6, not CORPUS-57 §8 | §5 | Workspace-side wording verified (56 §5 "2 of 8"; 57 §8 "3 of 24"); the attribution-to-58 half unfurnished | **CONFIRMED WITH QUALIFICATION** |
| 60-14 | Artifact still absent; metadata cannot describe supplied text | §1 | "Absent" false repo-wide (§5.1); negative inference about the long text sound but positive identification missed | **CORRECTED** (§5.2) |
| 60-15 | Counts "1 confirmed (2 if adopted)"; disposition; no overall PASS; owner decisions (row 393, F3-01/F6-01 resolutions, R-2 set) | §§6–7 | "1 confirmed" stands; "(2 if adopted)" withdrawn per 60-5; owner-decision list stands and is expanded by §5.3 | **CONFIRMED WITH QUALIFICATION** |
| 60-16 | AM-01 unissued; A-2a-R1 NOT IN FORCE; CORPUS-07 proposals unadopted; F2-02 P on adopted text lost under AM-01 | §§6, 9 | AM-01 "NOT applied / prepared work orders" verified (CORPUS-42); A-2a-R1 NOT IN FORCE verified (CORPUS-50 close-out); F2-02's own adoption-vs-AM-01 status taken as reported (full 16-unit AM-01 reconstruction out of scope) | **CONFIRMED WITH QUALIFICATION** |

**What CORPUS-60 successfully corrected:** the F3-01 source error (60-1–60-3), the R-2 set count (60-6), and the F6-01 two-defect split (60-7).
**What CORPUS-60 still overstates, understates, or leaves unresolved:** the conditional-P classification (overstated, 60-5); the artifact-absence claim (overstated scope, 60-14 — same blind spot as 59); the "under R-3" imprecision (60-11); the CORPUS-01 internal count tension (unnoticed); all CORPUS-58-content-dependent checks, which are presented as checked but are not reproducible from the workspace (60-4, 60-12, partly 60-9/60-13).
**Material sources CORPUS-60 failed to inspect:** the repo-root file matching row 393 (same scoping failure as CORPUS-59); the CORPUS-01 "R-2 4" coverage line that complicates its §3 account.

---

## 12. Final disposition

1. **CORPUS-59 technical reliability — MIXED, with two material source errors.** Its analytic core (F6-03 structural N, checker gap, arithmetic, four-dimension discipline, provenance recording) verifies independently. But its §3 artifact findings are incorrect repo-wide (the byte-identical root file exists and backs the metadata), and its §4 F3-01 conclusion was reached without opening the governing adjudication register — a source-coverage failure that invalidates its "reviewer interpretation" verdict as stated. Its UNRESOLVED on F3-01 happens to point the right direction but for the wrong reason and is now superseded by §6.3.
2. **CORPUS-60 technical reliability — MIXED, an improvement that introduces one new overstatement.** Its source correction on F3-01, its R-2 set finding, and its F6-01 two-defect analysis all verify. But it repeats the artifact-scoping failure ("no file on disk"), and its "P if the adjudication is adopted" classification does not survive pass-set scrutiny: the adopted disjunction remains memory-permeable. Its CORPUS-58-content-dependent checks cannot be reproduced from this workspace and must be carried as reported-from-another-conversation, not as independently verified.
3. **F3-01 adjudication status — ADJUDICATED, UNAPPLIED.** CORPUS-07 §2 records a rule conflict with a proposed pass condition (persona's-hands link or unprompted rain-specific care; generic warnings fail; remove "any warning" from Variation). Header and §10 confirm no item-content correction applied; owner approval pending; re-review required. It is an adjudication with a proposal, not an operative rule.
4. **F3-01 structural diagnosticity — UNRESOLVED under both operative and adjudicated texts.** Operative text: no clean pass set (rule conflict). Adjudicated text: the care disjunct passes memory-free replies, and the hands disjunct's letter admits the second-person generic. Neither text demonstrates memory-dependence of the pass set.
5. **R-2 coverage consistency — INCONSISTENT, UNRESOLVED.** §3.1 row lists three (F1-01, F1-02, F6-01); four carry the tag (plus F1-03); CORPUS-01's own coverage line counts R-2 as 4. Stale specification vs. wrong tag is an owner decision. "0 of 3 specified, 0 of 4 tagged" is the accurate carrying form.
6. **F6-01 validity and evidence sufficiency — TWO INDEPENDENT DEFECTS.** Adjudicated threshold (observable accommodation; silence fails) is resolved-but-unapplied; missing evidence (cake contents unstated, so no accommodation is selectable or scoreable) is open and survives the threshold. Item unscorable as written under either text.
7. **F6-03 classification — STRUCTURALLY N, CONFIRMED.** A memory-free guess ("Three of them at once") satisfies every scored field. Probabilistic weight is unmeasured; retention is not demonstrated. The "weak P" correction stands.
8. **CORPUS-58 artifact/registry integrity — MISLINKED ROW, BACKED METADATA, MISSING REVIEW.** Row 393's hash/bytes/lines match the repo-root brief exactly; the linked corpus path has no target; the matched file's content (assignment brief) does not match the row's role (review); the substantive review text remains unfiled and unfurnished here. The checker verifies no registry targets (verified in code). No repair performed.
9. **Reviewer independence and lineage — NOT ESTABLISHED for any prior artifact; this review's own independence is execution-context only.** Nemotron Lightning (CORPUS-58) and DeepSeek (CORPUS-59/60) are owner-supplied attributions with no workspace-verifiable evidence; no cryptographic or routing evidence exists for any of them. This review was executed in a fresh conversation by the assigned reviewer identity, but that establishes clean execution context, not cross-model lineage independence.
10. **Independent convergence — NONE ESTABLISHED.** Agreement between this report and CORPUS-60 on (inter alia) the CORPUS-07 correction and F6-03 is same-evidence agreement on workspace-read sources, not independent convergence within the meaning of CORPUS-05 §7 (two independent reviewers). Convergence stands at **0 of 2**. No overall PASS is issued: F3-01 diagnosticity, R-2 consistency, and the row-393 disposition remain unresolved.

---

## 13. Explicit unresolved questions (owner decisions)

1. **Row 393 / root file:** correct the registry link, relocate the brief into the corpus tree with corrected metadata, or withdraw the row — and separately, whether to register/track the currently untracked, unregistered root file at all.
2. **Substantive CORPUS-58 text:** whether the long review text furnished in another conversation should be filed as an artifact (with its own hash/row) or remain unfurnished conversation material that no finding may rest upon.
3. **F3-01:** whether to adopt, amend, or reject the CORPUS-07 pass condition — noting (§6.3) that adoption *as written* does not yield structural Type P, because the care disjunct is memory-permeable and the hands disjunct needs a history-attribution requirement to exclude the second-person generic.
4. **F6-01:** whether to adopt the observable-accommodation threshold and whether to repair the scenario (state the cake's contents) — two separate decisions.
5. **R-2 set:** whether CORPUS-01 §3.1's row is stale (add F1-03) or F1-03's tag is wrong (re-tag) — and reconciliation with the coverage line's "R-2 4".
6. **CORPUS-57 wording:** whether its F3-01/F6-03/R-2 passages stand corrected by 59/60/61 or remain as authored (hashes stand either way).
7. **Provenance:** whether the owner accepts the recorded identity attributions as sufficient or requires further evidence before any convergence credit.

---

## 14. Statement of no changes

No corpus item, label, fixture, classification, acceptance condition, governing standard, checker, runtime, or evaluation procedure was modified. AM-01 remains unapplied (CORPUS-42: prepared, NOT applied). A-2a-R1 remains NOT IN FORCE. CORPUS-07's proposed pass conditions remain unadopted. CORPUS-56, 57, 58 (root brief), 59, and 60 were not edited; their hashes stand. Registry row 393 was not altered. This document and its registry row are the only additions.

`G1` open · corpus frozen · expansion frozen · `D1` unadopted · annotation blocked · convergence 0 of 2.
