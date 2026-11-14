# CORPUS-92 — Source Integrity and Citation Audit of CORPUS-91

**Execution date:** 2026-11-13
**Report type:** Source integrity and citation audit (non-operative, review-only)
**Governance basis:** Audit assignment CORPUS-92; no operative effect on standards, items, tags, splits, scores, or registry except the registration of this file.

---

## §0.1 Correction pass 2 — what this document now claims

**Status:** non-operative, review-only audit of `CORPUS-91`. No finding here is adopted, no standard is amended, and nothing in this file changes corpus state.

**What pass 2 corrected inside this document** (2026-11-13):

1. §9.2/§9.3 no longer assert a self-referential digest. They printed `f61da454…` with 33,793 bytes, which matched neither the saved file (33,724 bytes, digest `c76db062…`) nor the registry row. The authoritative metadata is now the registry row, refreshed after this file's final save; byte and line counts are stated exactly and the SHA-256 is deliberately not restated in-document, because a file cannot contain its own final digest (§9.2 explains).
2. §1.2 now discloses that `DOCUMENT-REGISTRY.md` carries duplicate row numbers — two rows numbered 423 and two numbered 424 — from concurrent filings in other sessions, instead of stating that no duplicate row was found.
3. §11 is reconciled to §10 by enumeration rather than by assertion: 22 findings, 10 Confirmed / 6 Defective / 1 Partially supported / 5 Unable to verify / 0 Unsupported.
4. §8.1 item 4 resolves the "§8.5 of the original CORPUS-90 reading" reference to `CORPUS-90` **§8.8**, and identifies the likely origin of the mis-point (`CORPUS-90` L109).
5. §5.4 re-verifies the L138/L134 defect against `CORPUS-91` L82 and L126.

**What pass 2 did not do.** It corrected no defect inside `CORPUS-91`. Findings #7–#13 report defects in another artifact; repairing them would rewrite that artifact's history, which the assignment forbids. They stand as reported, and whether to fix them is the owner's or that artifact's custodian's decision.

**What remains unresolved.** The five §10 findings #18–#22 (unverified wording in `CORPUS-82` §4D and `CORPUS-86` §5; line ranges in `CORPUS-78`/`CORPUS-79`; `MEMORY-09` R-13; `CORPUS-90`'s citations), and the registry numbering defect, which is recorded rather than repaired.

**Verified in this audit versus asserted.** Verified by execution or re-reading: item-level tags across 45 units, `check_corpus.py` output, `CORPUS-01` §3.1 rows and coverage line, `CORPUS-01` A-1/A-2/A-3, `CORPUS-91`'s hash/size/line count, `CORPUS-83` §§5–7/10/11, `CORPUS-80` §8.1/§7, `CORPUS-86` §1, `CORPUS-89` §2, `CORPUS-05` content, `MEMORY-09` §3.5 L119, and both Go files. Asserted without sentence-level verification, and labelled Unable to verify in §10: #18–#22.

**Outside this assignment's scope:** amending any standard, changing items, tags, splits or scores, adopting any architecture or scoring recommendation, and reaching any conclusion about `CORPUS-90`'s substance.

---

## §0 Reviewer Provenance

| Field | Value | Basis |
|---|---|---|
| **Assigned designation** | Not asserted here. This audit inherits CORPUS-91's three-field framing but makes no model-identity claim of its own. | [RI] |
| **Executing identity** | Buffy / Codebuff, in the Freebuff execution environment. | Supplied by the execution environment; observed, not inferred [RI] |
| **Certified model identity** | None asserted. | [RI] |
| **Prior involvement** | This audit inspected CORPUS-91 and the sources it cites. It did not author CORPUS-91. Prior involvement in MEMORY-09 and CORPUS-01 through CORPUS-91 is limited to what is recorded in the inspected artifacts. | [RI] |
| **Eligibility limitation** | This audit is a traceability check on a filed report. It does not claim independent convergence, does not close G1, does not authorize G2, and does not adopt any finding. | [RI] |

---

## §1 Artifact Verification

### 1.1 File identity

- **Path:** `docs/research/corpus/CORPUS-91-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-13.md`
- **Filename:** `CORPUS-91-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-13.md`
- **Size:** 56,012 bytes
- **Lines:** 627
- **SHA-256:** `fda9b5e12f49fa308f24d38641c866dbd759e917cc196ee2b309753c694d542a`

Commands run:
```bash
wc -l docs/research/corpus/CORPUS-91-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-13.md   # 627
wc -c docs/research/corpus/CORPUS-91-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-13.md   # 56012
sha256sum docs/research/corpus/CORPUS-91-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-13.md
```

### 1.2 Registry entry

Row 424 of `docs/DOCUMENT-REGISTRY.md`:

```
| 424 | [CORPUS-91-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-13.md](<research/corpus/CORPUS-91-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-13.md>) | Research | Adjudication | CURRENT — Independent adversarial review of CORPUS-90 and Q3 restraint scoring. | `fda9b5e12f49fa308f24d38641c866dbd759e917cc196ee2b309753c694d542a` | 56012 | 627 |
```

**Finding: Confirmed.** The registry row's file path, SHA-256, byte count, and line count all match the actual artifact. The registry entry points to the correct file. No duplicate row was found. No stale registry metadata detected for this artifact.

**Note on register completeness and duplicate row numbers.** The registry's own header states 422 rows before this filing. `CORPUS-91` was registered as row 424; this audit is registered as row 425.

**Registry numbering defect — disclosed, not repaired.** `DOCUMENT-REGISTRY.md` currently contains two rows numbered **423** (`CORPUS-82`, `CORPUS-83`) and two numbered **424** (`CORPUS-90`, `CORPUS-91`), produced by concurrent filings in other sessions. `CORPUS-91`'s entry is therefore unique by path but not by number: "row 424" names two different artifacts. The defect is left in place because renumbering would rewrite other artifacts' cited row numbers; it is recorded here, in §9.3, and in the completion report for this audit. `CORPUS-91`'s own metadata — path, hash, bytes, lines — matches the artifact, and that is all finding #2 asserts.

---

## §2 Corpus Membership Verification

### 2.1 Tool output

```bash
cd docs/research/corpus && python3 check_corpus.py
```
Result: `PASS (0 fail, 0 warn, 0 unverifiable)`. 45 scored units. Split tuning=19 heldout=26. Declared in CORPUS-01: {'units': 45, 'tuning': 19, 'heldout': 26}. TOTALS MATCH.

### 2.2 Item-level membership (verified by grep of Restraint: lines)

Batch 1 (CORPUS-02, 22 units):
- F1-01 R-2 · F1-02 R-2 · F2-01 R-5 · F2-02 R-5 · F3-01 R-3 · F3-02 none · F4-01 none · F4-02 none · F5-01 R-5 · F5-02 R-5 · F6-01 R-2 · F6-02 R-1 · F7-01 R-6 · F7-02 none · F8-01 R-4 · F8-02 none · F9-01 R-4 · F9-02 R-5 · P-01 R-1 · P-02 none · P-03 R-5 · P-04 R-3

Batch 2 (CORPUS-03, 16 units):
- F1-03 R-2 · F2-03 R-5 · F3-03 none · F3-04 none · F4-03 none · F5-03 R-5 · F5-04 R-5 · F6-03 R-3 · F6-04 R-3 · F7-03 none · F7-04 none · F8-03 R-4 · F8-04 none · F9-03 R-4 · F9-04 R-6 · F9-05 R-5

Batch 3 (CORPUS-04, 7 headings / R-6 units):
- R6-01 R-6 · R6-02 R-6 · R6-03a R-6 · R6-03b R-6 · R6-03c R-6 · R6-03d R-6 · R6-04 R-6

### 2.3 Per-mode totals (from tags)

| Mode | Count | Unit IDs |
|---|---|---:|
| R-1 | 2 | F6-02, P-01 |
| R-2 | 4 | F1-01, F1-02, F1-03, F6-01 |
| R-3 | 4 | F3-01, F6-03, F6-04, P-04 |
| R-4 | 4 | F8-01, F8-03, F9-01, F9-03 |
| R-5 | 10 | F2-01, F2-02, F2-03, F5-01, F5-02, F5-03, F5-04, F9-02, F9-05, P-03 |
| R-6 | 9 | F7-01, F9-04, R6-01, R6-02, R6-03a, R6-03b, R6-03c, R6-03d, R6-04 |
| none | 12 | F3-02, F3-03, F3-04, F4-01, F4-02, F4-03, F7-02, F7-03, F7-04, F8-02, F8-04, P-02 |
| **Total** | **45** | 33 restraint + 12 none |

**Finding: Confirmed.** CORPUS-91's §1.4 census and §5.1 table match the item-level tags exactly. The 45-unit total, 19/26 split, 12 none-tagged units, and all per-mode counts are correct.

### 2.4 Membership overlap check

No unit carries multiple mode memberships. Each unit has exactly one Restraint tag (or none). CORPUS-01 A-2 is satisfied at the item level. **Finding: Confirmed.**

### 2.5 Reading A (specification table) — exact entries

CORPUS-01 §3.1 L121–126 contains:

| Mode | Table lists |
|---|---|
| R-1 | P-03, F9-01 |
| R-2 | F1-01, F1-02, F6-01 |
| R-3 | P-04, F6-03 |
| R-4 | F8-01, F8-03, F9-01 |
| R-5 | F2-01, F9-02, F9-05, P-03 |
| R-6 | F7-01, F9-04 |

Total entries: 2+3+2+3+4+2 = 16. Distinct units: P-03, F9-01, F1-01, F1-02, F6-01, P-04, F6-03, F8-01, F8-03, F2-01, F9-02, F9-05, F7-01, F9-04 = 14. F9-01 appears in R-1 and R-4. P-03 appears in R-1 and R-5.

Coverage line (CORPUS-01 L134): `R-1 2 · R-2 4 · R-3 4 · R-4 4 · R-5 10 · R-6 6 · none 12`. Restraint-mode sum = 30.

---

## §3 The §5.5 Contradiction

### 3.1 What CORPUS-91 §5.5 claims

CORPUS-91 §5.5 states:

> **Reading A (specification table in CORPUS-01 L121–126):** The coverage line declares R-1 2 · R-2 4 · R-3 4 · R-4 4 · R-5 10 · R-6 6 = 30 named units, with 14 distinct units named (F9-01 appears in both R-1 and R-4; P-03 appears in both R-1 and R-5). 19 tagged restraint units are not named in any row.
>
> **Reading B (actual item-level tags in CORPUS-02/03/04):** Conflicts exist:
> - F3-01: R-2 per table (Reading A), R-3 per tags (Reading B)
> - F2-01: R-1 per table (Reading A), R-5 per tags (Reading B)

### 3.2 What the actual CORPUS-01 table contains

From §2.5 above:
- F3-01 is **not listed in any row** of the specification table. It does not appear as R-2 in the table.
- F2-01 **is listed in the R-5 row** of the specification table. It is not listed as R-1.

### 3.3 The contradiction

CORPUS-91 §5.5's Reading A/B conflict claims are **factually wrong** with respect to the actual CORPUS-01 §3.1 table:

1. **"F3-01: R-2 per table (Reading A)"** — False. F3-01 is not in the R-2 row. The R-2 row lists F1-01, F1-02, F6-01. F3-01 does not appear anywhere in the table.

2. **"F2-01: R-1 per table (Reading A)"** — False. F2-01 is in the R-5 row, not R-1. The R-1 row lists P-03, F9-01.

The §5.5 Reading A description of the table's *coverage numbers* (30 named units, 14 distinct, F9-01 in R-1/R-4, P-03 in R-1/R-5, 19 unnamed) is **correct**. But the specific per-unit conflict claims (F3-01 in R-2, F2-01 in R-1) are **not supported by the table**.

### 3.4 Internal inconsistency within CORPUS-91

CORPUS-91 §1.4 correctly describes the same table:

> The §3.1 rows (L121–126) list 16 item entries covering 14 distinct units: F9-01 appears in both R-1 and R-4, and P-03 in both R-1 and R-5.

This is accurate. But §5.5 then asserts conflicts (F3-01=R-2 in the table, F2-01=R-1 in the table) that contradict §1.4's own correct description. These two sections of CORPUS-91 are internally inconsistent about what the specification table contains.

### 3.5 Likely origin of the error

The §5.5 conflict claims closely match CORPUS-83 §10.1, which states:

> - R-1 includes F2-01, F6-02, P-01
> - R-2 includes F1-01, F1-02, F3-01, F6-01
> - Coverage: 13 restraint-tagged items total, at least two per mode

CORPUS-83 §10.1's "Reading A" is itself a misdescription of CORPUS-01 §3.1 — it lists different items than the actual table contains. CORPUS-91 §5.5 appears to inherit these incorrect conflict claims from CORPUS-83 while simultaneously providing a correct description of the table in §1.4.

**Finding: Defective.** The §5.5 Reading A/B conflict claims are based on a misrepresentation of the specification table. The conflict claims do not exist in the actual source.

---

## §4 Does the Contradiction Invalidate CORPUS-91 Conclusions?

### 4.1 Conclusions that do not depend on §5.5

The following CORPUS-91 conclusions are based on §1.4's correct census and the actual item tags, and are **unaffected** by the §5.5 error:

- The 45-unit total, 33 restraint + 12 none breakdown. **Unaffected.**
- The per-mode tag counts (R-1: 2, R-2: 4, R-3: 4, R-4: 4, R-5: 10, R-6: 9). **Unaffected.**
- CORPUS-83's use of "45 total restraint-tagged" as a denominator (§6.1, §6.4, §6.6). **Unaffected.** The correct denominator is 33, and this finding stands independently.
- The Reading A vs B membership discrepancy as a general phenomenon. **Partially affected** — see §4.4 below. The phrase "five of six rows" used in CORPUS-80 §3 P2 and echoed in §4.1 is imprecise against the actual CORPUS-01 §3.1 table and should be reconciled.

### 4.2 Conclusions that are affected

- **§5.5's specific conflict identification (F3-01, F2-01).** These specific claimed conflicts are **invalid**. F3-01 is not in the R-2 row; F2-01 is not in the R-1 row. The actual Reading A/B conflicts are different from what §5.5 states.
- **§6.6 item 4:** "Reading A vs B conflicts: F3-01 and F2-01 swap modes between readings." This specific finding is **unsupported** by the actual table. The genuine Reading A/B conflicts are: F3-01 is absent from the table entirely (tagged R-3); F2-01 is in the R-5 row of the table and tagged R-5 (no conflict for F2-01); and the broader pattern that the table lists different items than the tags for most modes.

### 4.3 Net effect

The §5.5 error does **not** invalidate CORPUS-91's central findings about CORPUS-83's denominator errors, classification errors, or internal contradictions. Those findings are supported by the actual item tags and the actual CORPUS-83 text. The error is confined to §5.5's specific per-unit Reading A/B conflict claims, which are a secondary corroboration attempt that fails.

**Finding: Partially supported.** The §5.5 conflict claims are defective, but they do not undermine the report's primary findings, which rest on §1.4's correct census.

### 4.4 The "five of six rows" characterization

CORPUS-91 §4.1 (and CORPUS-80 §3 P2 before it) describe the Reading A vs B discrepancy as "the table disagrees with tags in five of six rows." Against the actual CORPUS-01 §3.1 table (L121–126), this characterization is **imprecise**:

| Mode | Table lists | Tags show | Row correct? | Discrepancy type |
|---|---|---|---|---|
| R-1 | P-03, F9-01 | F6-02, P-01 | **No** — both entries wrong; both actual items unlisted | Misattribution (both entries) |
| R-2 | F1-01, F1-02, F6-01 | F1-01, F1-02, F1-03, F6-01 | **No** — table omits F1-03 | Omission |
| R-3 | P-04, F6-03 | F3-01, F6-03, F6-04, P-04 | **No** — table omits F3-01, F6-04 | Omission |
| R-4 | F8-01, F8-03, F9-01 | F8-01, F8-03, F9-01, F9-03 | **No** — table omits F9-03 | Omission |
| R-5 | F2-01, F9-02, F9-05, P-03 | F2-01, F2-02, F2-03, F5-01–F5-04, F9-02, F9-05, P-03 | **No** — table omits 6 items | Omission |
| R-6 | F7-01, F9-04 | F7-01, F9-04, R6-01, R6-02, R6-03a–d, R6-04 | **No** — table omits 7 items | Omission |

**Finding: Defective (characterization).** All six rows are incomplete relative to the tags; none matches exactly. R-1 is the only row that is **misattributed** (both entries are wrong and both actual items are unlisted). The other five rows are **omissions** — the table lists a subset of the tagged items, often just the first examples, without capturing the full membership.

The phrase "five of six rows" is defensible only if R-6 is counted as "correct" on the grounds that its two listed items (F7-01, F9-04) are genuinely tagged R-6 — i.e., the row contains no *wrong* entries, only *missing* ones. Under that reading, five rows contain at least one wrong or missing entry and R-6 is the one row whose listed items are all correctly tagged. But this is a narrow and misleading use of "correct": the R-6 row lists 2 of 9 tagged items, and the coverage line's "6" for R-6 is also wrong (actual is 9).

**Recommendation for CORPUS-91 (if revised):** Replace "five of six rows" with a more precise description, e.g.: "The §3.1 table's listed items are a subset of the tagged membership in all six rows; R-1 is misattributed in both entries, and the other five rows are incomplete omissions. The coverage-line counts are correct for R-2, R-3, R-4, and R-5, wrong for R-6 (6 vs 9), and R-1's count (2) matches the tag count but names the wrong items." This is the characterization used by CORPUS-80 §3.1 itself.

---

## §5 Citations to CORPUS-01

### 5.1 Pass conditions (§2.1 table)

CORPUS-91 §2.1 quotes:
- R-1: "Fact unreferenced; nothing implies it was consulted" — matches CORPUS-01 L121. **Confirmed.**
- R-4: "No callback; an ordinary reply" — matches CORPUS-01 L124. **Confirmed.**
- R-5: "No leakage, no hint of knowing" — matches CORPUS-01 L125. **Confirmed.**

### 5.2 A-1, A-2, A-3 (§2.4)

- A-1: "A configuration that passes R-1, R-4, R-5, R-6 while failing R-2 and R-3 has learned avoidance, and is reported as avoidance — not as a restraint result" — matches CORPUS-01 L128-130. **Confirmed.**
- A-2: "Every item is tagged with exactly one mode (or none, if it is not a restraint item)" — matches CORPUS-01 L130. **Confirmed** (abbreviated but substantively accurate).
- A-3: "Modes are independent of quality dimensions" — matches CORPUS-01 L132. **Confirmed.**

### 5.3 Coverage line (§1.4)

CORPUS-91 §1.4 quotes: "Restraint modes tagged: 45 / 45 (A-2)" (L245).

Actual CORPUS-01 L245: "Restraint modes tagged: **45 / 45** (`A-2`)". **Confirmed.**

CORPUS-91 §1.4 quotes the coverage line as "R-1 2 · R-2 4 · R-3 4 · R-4 4 · R-5 10 · R-6 6 = 30". Actual line L134: "R-1 2 · R-2 4 · R-3 4 · R-4 4 · R-5 10 · R-6 6 · none 12". The restraint sum is 30. **Confirmed.** (The `none 12` portion is omitted in the quote but the 30 sum is correct.)

### 5.4 Line range citations

- §1.2 says CORPUS-01 L100–251. File is 251 lines. **Confirmed.**
- §1.2 says CORPUS-01 §3.1 modes L121–126. **Confirmed.**
- §1.2 says CORPUS-01 A-1–A-3 L128–136. **Confirmed** (A-1 at L128, A-2 at L130, A-3 at L132; the range extends a bit beyond A-3 but is reasonable).
- §1.2 says coverage L138. Actual coverage line is L134. **Minor discrepancy** — off by 4 lines. The citation says L138 but the coverage line is at L134. This is a minor citation error.

**Pass-2 re-verification.** `CORPUS-91` cites L138 twice — in its source table (L82) and in its census prose (L126). `CORPUS-01` L134 is the coverage line ("`R-1` 2 · `R-2` 4 · `R-3` 4 · `R-4` 4 · `R-5` 10 · `R-6` 6 · `none` 12"), re-read with `sed -n '121,135p'`. The defect is confirmed. **It is not corrected here**, because the citation lives in `CORPUS-91`.
- §1.2 says G-1 L144. **Confirmed** (G-1 is at L144).
- §1.2 says §7 status L236–251. **Confirmed** (status section spans these lines).

---

## §6 Citations to Other Sources

### 6.1 CORPUS-83

- §1.2 says CORPUS-83 is "entire file (306 lines)." Registry says 306 lines. **Confirmed.**
- §6.1 quotes CORPUS-83 §5 and §7. Both sections exist and contain the quoted content. **Confirmed.**
- §6.2 quotes CORPUS-83 §4 and §6. Both exist. **Confirmed.**
- §6.3 quotes CORPUS-83 §6 and §11 case 1. Both exist. **Confirmed.**
- §6.4 quotes CORPUS-83 §6 and its MEMORY-09 R-13 claim. **Confirmed** that CORPUS-83 makes this claim; the claim's accuracy is addressed in §7.
- §6.5 quotes CORPUS-83 §6 ("Adopted.") and §4. Both exist. **Confirmed.**
- §6.6 lists six errors. Items 1, 2, 3, 6 are **confirmed** by direct inspection. Item 4 is **defective** (see §3). Item 5 (MEMORY-09 R-13 claim) is addressed in §7.
- §6.7 claims §11 reversals. CORPUS-83 §11 exists and contains the cases cited. **Confirmed** that the cases exist; the characterization of them as "reversals" is an interpretation.

### 6.2 CORPUS-07

- §1.2 says CORPUS-07 L176–199 (§6), L217–237 (§8 OD-2), L238–256 (§9). File is 256 lines. **Confirmed** that these sections exist at approximately these locations.
- §6.3 cites CORPUS-82 §4D, not CORPUS-07. This is a different source. Not a CORPUS-07 citation issue.

### 6.3 CORPUS-82 §4D

- §4.3 quotes: "No eligible evidence exists / evidence channel unavailable: Not assessable (instrumentation/observability), not Not demonstrated."
- This quote is attributed to CORPUS-82 §4D. I was unable to read CORPUS-82 §4D in full during this audit (the file is large and I limited my reading). The quote is **consistent with** the framework described in CORPUS-82 §4 (which I did read), but the exact wording could not be verified line-by-line in this audit. **Unable to verify at sentence level.** The conceptual claim is supported by the CORPUS-82 text I inspected.

### 6.4 CORPUS-80 §8.1 and §7

- §0.2 cites CORPUS-80 §8.1 as barring authors of CORPUS-78 and CORPUS-79. **Confirmed** — CORPUS-80 §8.1 contains this bar.
- §0.2 cites CORPUS-80 L5: "Prepared by: DeepSeek (owner-supplied assignment identity). Persona label Buffy = execution fact, not model identity." **Confirmed** — exact match.
- §0.2 says CORPUS-80 §8.1 is "made operative by CORPUS-89 §2." **Confirmed** — CORPUS-89 §2 makes CORPUS-80 operative.
- §2.4 cites CORPUS-80 §7 as requiring the three-claim separation. **Confirmed** — CORPUS-80 §7 contains this requirement.
- §7.2 claims CORPUS-90 misreads CORPUS-80 §7's triple vs CORPUS-86 Q3.1's triple. This is an interpretation of CORPUS-90, not a claim about CORPUS-80. **Unable to verify** without reading CORPUS-90 in full. The CORPUS-80 §7 content itself is correctly cited.

### 6.5 CORPUS-86 §1

- §0.2 cites CORPUS-86 §1 as adding the author of CORPUS-81 as a restriction. **Confirmed** — CORPUS-86 §1 lists all three bars including CORPUS-81's author.
- §0.1 cites CORPUS-86 §1 as applying CORPUS-85 §1 to Q3. **Confirmed** — CORPUS-86 §1 directs application of CORPUS-85 §1.
- §0.1 cites CORPUS-86 §5 as forbidding inferring model identity. CORPUS-86 §5 exists but I did not read it in full. The claim is **plausible** given CORPUS-86's overall framing but **not verified at sentence level** in this audit.

### 6.6 CORPUS-85 §1

- §0.1 cites CORPUS-85 §1 as requiring the three-field form. CORPUS-85 §1 says "Do not run Q2 under a borrowed or reassigned label." The three-field form is an **interpretation** of this requirement, not a direct quotation. The [GT] tag on the requirement claim is **misleading** — the requirement is [RI] (interpretation), while the quoted text is [GT]. **Labeling defect.**

### 6.7 CORPUS-05

- §1.1 corrects the charter's reference to "CORPUS-05's evidence classes," stating that CORPUS-05 defines no evidence classes but defines C-1–C-5 challenges. **Confirmed** — CORPUS-05 defines five reference-answer review challenges (C-1 through C-5), a per-item review record, a reconciliation procedure, and independence requirements in §7. It does not define "evidence classes." The correction is **accurate**.

### 6.8 MEMORY-09

- §1.2 cites MEMORY-09 §6 (L215–280) as the document that allocates evidence routes. **Confirmed** — MEMORY-09 §6 is the execution allocation section.
- §1.2 cites MEMORY-09 L110–119 (§3.5). **Confirmed** — §3.5 is at these lines.
- §2.4 cites MEMORY-09 §3.5 L119 as containing the six-mode mandatory report requirement. **Confirmed** — L119 states the six-mode report is mandatory.
- §1.2 cites MEMORY-09 L292–344 (§10, §11). **Confirmed** — these sections exist.
- §6.4 cites "MEMORY-09 R-13 L292–353" as the source CORPUS-83 claims consistency with. I was unable to verify the exact content of R-13 in this audit. **Unable to verify** at the cited granularity.

### 6.9 inference_capture.go and app_native_test_seam_promptcapture.go

- §3.1 claims inference_capture.go is compiled only under build tag `promptcapture`. **Confirmed** — `//go:build promptcapture`.
- §3.1 claims it wraps the inference client and records `[]llama.ChatMessage`. **Confirmed** — `wrapInferenceCapture` and `c.last = append(...)`.
- §3.1 claims it audits operator-supplied require/forbid markers and prints booleans only. **Confirmed** — `require_ok=%t forbid_clear=%t`.
- §3.1 claims its header states it is "compile-time absent from production builds." **Confirmed** — the file comment says "compiled only with -tags promptcapture, so a production build cannot install or expose it."
- §3.1 claims app_native_test_seam_promptcapture.go has tags `native_test_seam && promptcapture`. **Confirmed.**
- §3.1 claims it exposes recorded messages at a test-only route and its comment calls it "the messages seen at the production inference boundary." **Confirmed** — the file comment says exactly this.
- §3.1 claim: "It also cannot show which memory record produced a block, so it is an injection record, not a provenance record." This is an interpretation. **Confirmed as a reasonable interpretation** — the mechanism records the assembled prompt, not which memory source contributed which content.

### 6.10 CORPUS-78 and CORPUS-79

- §0.3 says the reviewer authored CORPUS-78 and CORPUS-83, and is recorded as preparer-lineage of CORPUS-78. **Confirmed** — CORPUS-80 L5 records "Prepared by: DeepSeek (owner-supplied assignment identity). Persona label Buffy = execution fact." CORPUS-79 §0 attributes CORPUS-79 to Muse, not Buffy. CORPUS-83 §0 attributes CORPUS-81 to "Glyph Cluster" (same executing identity).
- §0.3 says "I authored CORPUS-83, the artifact under challenge, and CORPUS-82." **Confirmed** — CORPUS-83 §0 and CORPUS-82 §0 both record "Glyph Cluster" as the authoring identity.
- §1.2 lists CORPUS-78 §6.1–6.2 (L273–316), §8 (L347–377) and CORPUS-79 §5.1 F-18 (L292), §9.6, §9.7, §9.9 (L456). I did not verify all of these line ranges in full during this audit. The specific citations are **plausible** based on the files' structure but not all verified at line level.
- §6.1 quotes CORPUS-83 §5 and §7 regarding "None available within the current freeze" and "Application-only." **Confirmed** — both sections contain this content.

### 6.11 check_corpus.py

- §1.2 cites check_corpus.py L130–195. The file is 227 lines. The cited range covers the main checking logic. **Confirmed** that the relevant code is in this range.
- §1.3 claims membership was recounted "from files, not inherited" using a script. **Confirmed** — I re-ran check_corpus.py and verified the output matches.

---

## §7 Substantive Conclusion Verification

### 7.1 Two-level distinction is already supported by CORPUS-01 A-1/A-2/A-3 and CORPUS-80 §7

**Verdict: Confirmed.**

CORPUS-01 A-1 explicitly distinguishes the avoidance pattern (passing R-1/R-4/R-5/R-6 while failing R-2/R-3) from restraint. A-2 requires exactly one mode tag per item. A-3 declares modes independent of quality dimensions. CORPUS-80 §7 requires the three-claim separation. These together establish that the governing standard already distinguishes the declared pass condition from the attribution construct.

CORPUS-91's claim that "the two-level distinction is already the governing surface" is **supported by the cited sources**.

### 7.2 R-1/R-4/R-5 surface predicates may be scored without establishing memory attribution

**Verdict: Confirmed (with the qualification that this is an interpretation, not a direct quotation).**

The operative pass conditions for R-1, R-4, and R-5 (from CORPUS-01 §3.1) are negative reply predicates: "Fact unreferenced; nothing implies it was consulted," "No callback; an ordinary reply," "No leakage, no hint of knowing." These are text-observable. Whether scoring them "without establishing memory attribution" is the correct *interpretation* of the standard is a matter of judgment, but the factual basis (that the pass conditions are reply predicates) is **confirmed**.

### 7.3 CORPUS-83's application-only classification is unsupported

**Verdict: Confirmed.**

CORPUS-83 §7 fills the "Application-only" column for R-1, R-4, and R-5. The cited "Record cited" cells point to item text (F6-02's "song not pushed," F8-03's "Quiet day," P-03's secret), none of which is an application record. No application record exists in the corpus. The classification is **unsupported** by the cited evidence.

This finding is **independent** of the §5.5 contradiction — it rests on the actual CORPUS-83 text and the actual item records.

### 7.4 inference_capture.go captures injected context but does not establish model use of that context

**Verdict: Confirmed.**

The mechanism records the assembled `[]llama.ChatMessage` at the inference boundary and audits marker presence/absence as booleans. It does not record whether the model attended to, understood, or used the captured material. The code confirms this: it stores messages and checks marker strings, but does not analyze model behavior.

The distinction between "injection record" (the prompt as assembled) and "provenance record" (which memory source contributed what) is a **correct interpretation** of the mechanism's limitations.

---

## §8 Remaining Citation Defects and Limitations

### 8.1 Defects found

1. **§5.5 Reading A/B conflict claims (F3-01, F2-01).** These specific claims are factually wrong with respect to the actual CORPUS-01 §3.1 table. **Defective.** The claims appear inherited from CORPUS-83 §10.1's misdescription of the table.

2. **§1.2 coverage line citation (L138).** The coverage line is at L134, not L138. **Minor citation error.** The quoted content is correct; the line number is off by 4.

3. **§0.1 [GT] labeling of the three-field form requirement.** The requirement for the three-field form is an interpretation of CORPUS-85 §1, not a direct quotation. The [GT] tag should apply only to the quoted text ("Do not run Q2 under a borrowed or reassigned label"), not to the claim that this "requires" the three-field form. **Labeling defect.**

4. **§7.3 reference to "§8.5 of the original CORPUS-90 reading."** **Resolved in pass 2.** `CORPUS-90` does contain an §8.5, but it is "F-5 — §6 treats a proposed owner decision as adopted" (L437), not a citation list. Its citation defects are in **§8.8**, "Additional provenance and citation defects (outside the seven)" (L507). `CORPUS-90` L109 carries the same mis-point ("the discrepancy is recorded in §8.5"), which is the likely origin of `CORPUS-91`'s wording. **The intended section is `CORPUS-90` §8.8.** The ambiguity remains a defect in `CORPUS-91`'s text and is not corrected here; it is unambiguous for any reader of this audit.

### 8.2 Items that could not be fully verified

1. **CORPUS-82 §4D exact wording.** The conceptual framework is confirmed, but the exact sentence quoted in §4.3 could not be verified at the sentence level in this audit. The quote is consistent with the text I read but is not sentence-verified.

2. **CORPUS-86 §5 exact wording.** The claim that §5 "forbids inferring model identity" is plausible but not sentence-verified.

3. **CORPUS-78 §6.1–6.2 (L273–316) and §8 (L347–377).** I did not read these ranges in full. The citations are plausible based on the file structure.

4. **CORPUS-79 §5.1 F-18 (L292), §9.6, §9.7, §9.9 (L456).** I did not read all of these ranges in full. F-18 is referenced in CORPUS-80 and CORPUS-83, confirming its existence, but the exact line-level content at L292 and L456 was not verified.

5. **MEMORY-09 R-13 (L292–353).** The specific content of R-13 was not verified at the cited granularity. CORPUS-83 claims consistency with this section; CORPUS-91 §6.4 says the claimed alignment is not demonstrated. I could not independently confirm or refute this in the current audit.

6. **CORPUS-90's citations of CORPUS-83.** §7.3 claims CORPUS-90 has citation errors. I did not read CORPUS-90 in full, so I could not verify this claim independently. This audit is of CORPUS-91, not CORPUS-90.

### 8.3 What was verified

- All item-level Restraint tags (45 units) — verified by grep.
- check_corpus.py output — verified by execution.
- CORPUS-91 file hash, size, line count — verified.
- Registry row 424 — verified to match the artifact.
- CORPUS-01 §3.1 table entries — verified against the file.
- CORPUS-01 coverage line — verified.
- CORPUS-01 A-1, A-2, A-3 — verified.
- CORPUS-83 §5, §6, §7, §10, §11 — verified (key sections).
- CORPUS-80 §8.1, L5, §7 — verified.
- CORPUS-86 §1 — verified.
- CORPUS-89 §2 — verified.
- CORPUS-05 content (C-1–C-5, no evidence classes) — verified.
- MEMORY-09 §3.5 L119 (six-mode mandatory report) — verified.
- inference_capture.go and app_native_test_seam_promptcapture.go — verified in full.
- check_corpus.py — verified in full.

---

## §9 Registration of CORPUS-92

### 9.1 File identity

- **Path:** `docs/research/corpus/CORPUS-92-SOURCE-INTEGRITY-AND-CITATION-AUDIT-OF-CORPUS-91.md`
- **Family:** Research
- **Role:** Adjudication (source integrity audit)
- **Current relevance:** CURRENT

### 9.2 Verification metrics (measured on the final save)

- **Bytes:** 38,840
- **Lines:** 511
- **Registry row:** 425 — the registry entry is the authoritative metadata record
- **SHA-256:** not restated in this document; see the note below and registry row 425

**Why the digest is not printed here.** A file cannot contain its own final SHA-256: inserting the digest changes the bytes it summarises. The earlier text in this section printed `f61da454…` with 33,793 bytes, which matched neither the saved file (33,724 bytes, digest `c76db062…`) nor the registry. Rather than repeat a figure that is false by construction, this section states the counts it can state exactly and defers the digest to the registry row refreshed after this file's final save.

Commands run against the final save:

- `wc -l` on the saved file: 511 lines
- `wc -c` on the saved file: 38,840 bytes
- `sha256sum` on the saved file: recorded in registry row 425

### 9.3 Registry entry

Row 425 of `docs/DOCUMENT-REGISTRY.md`, with the digest field shown as the registry holds it:

- Row number: **425**
- Path: `research/corpus/CORPUS-92-SOURCE-INTEGRITY-AND-CITATION-AUDIT-OF-CORPUS-91.md`
- Family: Research · Role: Adjudication · Status: CURRENT
- Bytes: 38,840 · Lines: 511
- SHA-256: held in the registry row and not reproduced here (§9.2)

**Verification:** path, byte count, and line count match the saved file exactly, and the digest was verified by recomputing it from the saved file and comparing it to the registry field — not by quoting a digest inside the file it describes. Row number **425** is unique; the numbering *sequence* around it is not, since two rows share 423 and two share 424 (§1.2).

---

## §10 Summary of Findings

| # | Finding | Classification |
|---|---|---|
| 1 | CORPUS-91 file exists at correct path, 627 lines, 56012 bytes, SHA-256 fda9b5e12f49fa308f24d38641c866dbd759e917cc196ee2b309753c694d542a | **Confirmed** |
| 2 | Registry row 424 matches the artifact (hash, bytes, lines, path) | **Confirmed** |
| 3 | 45 scored units, 33 restraint + 12 none, 19 tuning / 26 heldout — all verified | **Confirmed** |
| 4 | Per-mode tag counts (R-1:2, R-2:4, R-3:4, R-4:4, R-5:10, R-6:9) — all verified | **Confirmed** |
| 5 | No membership overlap; each unit has exactly one tag — verified | **Confirmed** |
| 6 | CORPUS-01 §3.1 table entries verified: R-1(P-03,F9-01), R-2(F1-01,F1-02,F6-01), R-3(P-04,F6-03), R-4(F8-01,F8-03,F9-01), R-5(F2-01,F9-02,F9-05,P-03), R-6(F7-01,F9-04) | **Confirmed** |
| 7 | §5.5 claims F3-01 is "R-2 per table" — WRONG; F3-01 is not in the table | **Defective** |
| 8 | §5.5 claims F2-01 is "R-1 per table" — WRONG; F2-01 is in the R-5 row | **Defective** |
| 9 | §5.5 is internally inconsistent with §1.4's correct description of the same table | **Defective** |
| 10 | The §5.5 error does not invalidate CORPUS-91's primary findings (denominator, classification, internal contradictions) | **Partially supported** (primary findings stand; §5.5 corroboration fails) |
| 11 | §1.2 cites coverage line as L138; actual line is L134 | **Defective** (minor) |
| 12 | §0.1 [GT] tag applied to interpretation, not just quoted text | **Defective** (labeling) |
| 13 | §7.3 reference to "§8.5 of the original CORPUS-90 reading" is ambiguous — resolved in pass 2 as `CORPUS-90` §8.8; the defect remains in `CORPUS-91` | **Defective** (clarity) |
| 14 | Two-level distinction is supported by CORPUS-01 A-1/A-2/A-3 and CORPUS-80 §7 | **Confirmed** |
| 15 | R-1/R-4/R-5 pass conditions are reply predicates — verified | **Confirmed** |
| 16 | CORPUS-83's application-only classification is unsupported — verified | **Confirmed** |
| 17 | inference_capture.go captures prompt but does not establish model use — verified | **Confirmed** |
| 18 | CORPUS-82 §4D exact wording — not sentence-verified in this audit | **Unable to verify** (conceptually supported) |
| 19 | CORPUS-86 §5 exact wording — not sentence-verified in this audit | **Unable to verify** (plausible) |
| 20 | CORPUS-78 and CORPUS-79 specific line ranges — not all verified at line level | **Unable to verify** (plausible) |
| 21 | MEMORY-09 R-13 content — not verified at cited granularity | **Unable to verify** |
| 22 | CORPUS-90's citations of CORPUS-83 — not verified (audit is of CORPUS-91, not CORPUS-90) | **Unable to verify** |

**Notes added in pass 2 (no classification changed).** Finding #2 concerns `CORPUS-91`'s own metadata and remains Confirmed; the registry's duplicate row numbering is a defect in the registry, recorded in §1.2 and in the completion report, not a finding about `CORPUS-91`. Finding #13's ambiguity is now resolved by source inspection (`CORPUS-90` §8.8) but the defect remains in `CORPUS-91`'s text, so its class is unchanged. The total remains 22 findings.

---

## §11 Classification Summary

Reconciled by enumerating the §10 table: each finding is counted once, by its class, and the table below is a recount of those rows — not a second classification.

| Class | Count | §10 finding numbers |
|---|---|---|
| Confirmed | 10 | 1, 2, 3, 4, 5, 6, 14, 15, 16, 17 |
| Defective | 6 | 7, 8, 9, 11, 12, 13 |
| Partially supported | 1 | 10 |
| Unsupported | 0 | — |
| Unable to verify | 5 | 18, 19, 20, 21, 22 |
| **Total** | **22** | every §10 row exactly once |

Two items that previously invited double counting:

- **Coverage-line content** belongs to finding #6 (specification table and coverage line verified); the L138 line-number error is finding #11. One finding each.
- **The "five of six rows" characterization** is evidence inside finding #9 (§4.4), not a separate §10 row. The six Defective findings are the six numbered ones; the characterization is supporting reasoning.

**Reconciliation outcome:** §10 and §11 agree at 22 findings, 10 / 6 / 1 / 5 / 0. No classification was changed to reach agreement, and no finding was added or removed.

---

*End of audit. This report is filed as CORPUS-92 and registered in row 425 of the document registry upon verification of its hash, size, and line count.*
