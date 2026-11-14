# Adversarial Review — C-1 TTS Capability Policy Amendment, Revision 6

**Review date:** 2026-11-07
**Reviewer:** Mimo (independent, read-only, targeted)
**Subject:** `AMENDMENT-DRAFT-howl-tts-capability-policy-c1-17.19-2026-11-07.md`, revision 6
**Governing contract:** `DESIGN-howl-tts-helper-architecture-2026-11-06.md` (frozen)
**Status of subject:** DRAFT — REVISION 6; NOT ADOPTED. This review adopts nothing.

---

## 1. Identity, and the hash discrepancy

### 1.1 Both subject hashes verified

| Artifact | SHA-256 (computed this pass) | Size |
|---|---|---|
| Amendment revision 6 | `dc5852adbc70641dbe8f233851eaec44266330fce197b78a0ab2641e7f73f784` | 320 lines / 65,274 bytes |
| Frozen contract | `4ff6b72ea03eef98e86ac158add8eb05a7e575cea22384470991f406e5547393` | 2,416 lines / 139,351 bytes |
| C-1 decision report | `c25a0d3fd38e0767a4b7bf278cb1d6617e2e53caf11053dabb71449578c4c8ad` | 302 lines |
| Revision 5 audit | `5944ae9e57a17076dbac9b29143036cf94257bff6c8bd642ff6cc9a9d6b76dbb` | 469 lines |
| Revision 4 audit | `772707bb5cb85b2002c668b602daeae0a330e3b2f3e06baaa8a6dc0edc550404` | 310 lines |
| Revision 3 audit | `9ceaab3038ac8855a87c53ef5c84d2e26a081ff98f3d12fb1e948635f417075b` | 412 lines |
| Revision 2 audit | `db6beb0a604cb65af4fe82a97d29db42395bc1312eabce051398b48ef4163bb3` | 330 lines |
| Revision 1 audit | `258fd1c1c981c2000ddcca4e87faea1e667d2b24114b68a48c8dea021d9c870e` | 358 lines |
| Stage 2 implementation record | `ec85b2b2c3760542699bcd406f1f3b6b3eb2cdc0165b9d6a3b21977a06741227` | — |

The supplied revision-6 hash is correct. It is no longer a received claim: the file on disk
reproduces it exactly, at 320 lines. Revision 5 was 302 lines / 60,121 bytes, so revision 6 adds 18
lines and 5,153 bytes.

The amendment's header also cites the revision 5 audit at `5944ae9e…`, which matches the audit file
written after the revision 5 pass. The review chain is intact and self-consistent.

### 1.2 The frozen-contract discrepancy: `…e5547393` is authoritative

The two candidate values:

```
A  4ff6b72ea03eef98e86ac158add8eb05a7e575cea22384470991f406e5547393   ← recorded by this reviewer, by six other artifacts, and by the file itself
B  4ff6b72ea03eef98e86ac158add8eb05a7e575cea22384470991f406f1f21818   ← earlier checkpoint, per the owner's message
```

**The file hashes to A.** `sha256sum` and an independent `openssl dgst -sha256` over
`docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` both return A. The file is 2,416 lines
and 139,351 bytes, and its mtime is `2026-11-06 14:00:11` — earlier than every revision of the
amendment and earlier than every audit of it, so the bytes under review are the bytes that were
frozen.

**A is also what this reviewer recorded independently, before revision 6 existed.** A appears in
artifacts written at six different times, none of which could have been edited by the revision 6
author:

| Artifact | Written | Records |
|---|---|---|
| Stage 2 implementation record | 16:59 | A |
| Revision 1 audit | 17:51 | A |
| Revision 2 audit | 19:02 | A |
| Revision 3 audit | 19:25 | A |
| Revision 4 audit | 19:50 | A |
| Revision 5 audit | 20:31 | A |

The amendment records A in its own header (line 6) and again in §5, and has recorded A in every
revision. Nothing in the repository records B: a bounded search of the working tree (excluding
`.git` and `node_modules`) finds `f1f21818` in no file at all.

**B cannot be a digest of any file.** A and B are both 64 hexadecimal characters and share their
first 56. A SHA-256 output is uniformly distributed, so two digests of *different* content agree in
their first 56 characters with probability 2⁻²²⁴. B is therefore not the hash of some other version of
the architecture document — it is a damaged copy of A. The damage is in the final eight characters,
the classic shape of a value re-typed or hand-edited rather than pasted.

**Conclusion.** `4ff6b72ea03eef98e86ac158add8eb05a7e575cea22384470991f406e5547393` is authoritative.
The transcription error is in the earlier checkpoint's record, not in any repository artifact, and
not in the revision 6 report. Nothing under `docs/review/` needs correcting; the checkpoint record
does, wherever it is kept. This reviewer cannot inspect a record held outside the repository, so if
that checkpoint is itself a tracked document, its frozen-hash line should be amended to A.

One residual possibility cannot be excluded by hash algebra alone: that the checkpoint intended to
fingerprint a *different* artifact and pasted this prefix by mistake. That reading is not credible —
the prefix is 56 characters of the architecture document's own digest — but if the owner has a second
candidate file, comparing it is a one-command check.

---

## 2. Verdict

**PASS — recommend adoption. No required fixes.**

Revision 6 does exactly what it claims and nothing else. N9 is closed at the decision, the state
summary, the map, and the checklist, with the ordinary path affirmed in all four places. N10's
identity qualifier now appears in every summary that describes recovery. E1–E4 are each closed at the
clause that carried them, and the four citation tails are corrected. The frozen contract, the decision
report, and all five prior audits are byte-identical to their pre-revision-6 values, and only the
amendment changed after the revision 5 review.

One editorial nit (E5) and one integration observation (I1) are recorded. Neither calls for a
revision, and neither touches policy. Adoption remains the owner's separate decision; this review is
not it, and nothing here authorizes Stage 3.

---

## 3. The five focus areas

### 3.1 N9 — ordinary `unloaded` synthesis versus `unavailable` recovery

**Pass, closed at four independent places.**

| Location | Text |
|---|---|
| A2-D12 (line 83) | "This decision governs recovery from `unavailable`; it does not alter the ordinary `unloaded` path. The ordinary `Synthesize` ensure-handshake/load behavior in §17.12.1 and §17.11.4 step 2 remains unchanged and may load an ordinary `unloaded` provider." Then: "A `synthesize` request while `State=unavailable` returns `provider_not_loaded` and does not initiate the retry." |
| A2-D5 (line 76) | "An ordinary `unloaded` provider may be loaded by `Synthesize`'s existing ensure-handshake/load path; only recovery from `unavailable` requires an explicit `load` under A2-D12." |
| §17.19.4 item 3 | "`Synthesize` may use its existing ensure-handshake/load path, so no separate explicit `load` is required before synthesis." |
| §17.19.4 item 11 | "…do not trigger that recovery; the ordinary `Synthesize` ensure-load from `unloaded` remains unchanged." |

Both halves of the old defect are gone. The prohibition is now scoped to a state (`State=unavailable`)
rather than to a method, and the method's frozen behavior is affirmed positively instead of being left
to inference. The over-broad phrasing that revision 5 carried — "`synthesize` … never trigger a load" —
does not appear anywhere in revision 6. Every "trigger" construction in the file (four occurrences, at
lines 83, 143, 249, 296) has a scoped object: "that recovery" or "it".

The supersession map now supports the same reading structurally. It adds a row for §17.12.1 (lines
97): "Confirmed by A2-D12: ordinary `Synthesize` ensure-handshake/load remains unchanged, including for
a provider that is merely `unloaded`." Because §17.19.2 controls "only for the matter identified in
§17.19.3", and the §17.11.4 row lists steps 4–6 and step 7 but not step 2, the precedence rule and the
decision text now agree instead of relying on the reader to notice a conflict. Revision 5 needed that
argument; revision 6 does not.

One consequence worth stating explicitly, because it is intended rather than overlooked: after a typed
non-fatal load failure the aggregate state is `unavailable` and a `synthesize` returns
`provider_not_loaded`, not the underlying `model_missing`. A2-D12 preserves the cause — "The prior
`LastError` remains until a new attempt replaces it" — so diagnostics survive, but the immediate code
is the less specific one. That is the price of explicit-load-only recovery and it matches owner
decisions 9 and 12.

### 3.2 N10 — identity qualifier in every recovery summary

**Pass.** The same three-field condition now appears wherever recovery is described, not only in the
decision entry:

| Location | Carrier |
|---|---|
| A2-D12 | "limited to a later explicit `load` for the same `provider_id`, `engine`, and `model_identity` as the failed attempt" |
| §17.7.4 new health row | "explicit `load` retry for the same identity under A2-D12" |
| §17.19.3, §17.7.4 row | "retry only on an explicit `load` request for the same `provider_id`/`engine`/`model_identity` under A2-D12" |
| §17.19.3, §17.7.2 row | "`unavailable` + explicit `load` request for the same identity (asset-specific failure, live helper) → `loading`" |
| §17.19.4 item 11 | "only for the same `provider_id`/`engine`/`model_identity` as the failed attempt; a differing identity follows A2-D7" |
| Impact row 16 and row 17 | "for the same identity" / "for the same identity only" |
| §3.2 F3 row | "for the same `provider_id`/`engine`/`model_identity`" |
| Owner decisions 9 and 12 | same three fields |

The exhaustive "only" in item 11 — the defect recorded as N10 — now carries the qualifier, so the
summary no longer permits by omission what the decision forbids by prohibition. The complement is
stated too: a differing identity follows A2-D7, and A2-D12 repeats that no differing-identity `load` is
sent to the old helper. All four reachable cases from the revision 5 review still resolve the same way.

### 3.3 E1–E4

**E1 — `helper_start_failed` availability: closed.** A2-D5 now closes the last unassigned cell of the
`Available` enumeration: "A `proc.Start` failure is `starting` → `unavailable` with typed
`helper_start_failed`; in that state `Available=false`, `ProcessAlive=false`, `Loaded=false`,
`Ready=false`, `ServiceCapabilitiesKnown=false`. This records the availability outcome only; it creates
no restart path."

Three checks on that sentence:

- It agrees with the frozen transition row. §17.7.2 line 1769 reads "`starting` | start/job failure |
  `crashed` or `unavailable` | Typed `helper_start_failed`". `unavailable` is one of the two permitted
  outcomes, so this is an assignment inside the frozen row rather than a supersession the map would
  have to record. A process that starts and then exits during handshake is a different case and remains
  `crashed` by way of the exit observer — `starting` has no process to crash, which is why the
  assignment is coherent.
- It agrees with §17.19.4 item 5, which states the same four flags and "no new restart path is
  created".
- It completes the enumeration rather than extending it. Five states are true, seven plus the three
  `unavailable` variants are false, and A2-D5 now names all thirteen. The thirteen-state matrix from
  the revision 5 review now has no cell that depends on inference.

**E2 — state-keyed `accepting_work`: closed.** A2-D12 attaches the value to the resulting state with a
parenthetical that removes the misreading: "the resulting `unloaded` state has `accepting_work=true`
(the existing value in §17.7.3's `unloaded` row)". The later restatement, the §17.7.3 map row, and
impact row 29 all attach it to "the resulting `unloaded` state". No text now raises the `loading` row's
`false`, and the additive edge remains labeled as an amendment.

**E3 — O-1 through O-7: closed.** §17.19.5 now carries "O-1 through O-7 remain unresolved and are not
decided by this amendment", and so do A2-D11 and §5. This restores the symmetry with A1 (§17.18.3,
line 2410) that revision 5 lacked. The substance is unchanged and correct: nothing in §17.19 closes any
of them. O-1 survives because the amendment constrains *loads*, not starts, and §17.7.2 line 1766 still
lets an `absent` helper start on a catalog request; O-3 is named and refused in A2-D10 and impact row
26; the rest are untouched. B12a and B12b remain deferred in the map and again in §5.

**E4 — budget exemption does not travel: closed, twice over.** The sentence that invited the
misreading now says "…the A2-D7 stop-and-replace sequence, without its budget exemption, govern", and
the following sentence states the rule positively: "A2-D12 itself does not inherit that exemption: its
in-place retry is not an identity replacement, and any process restart remains governed by §17.11.5."
The map's §17.11.5 budget row and impact row 23 say the same, and owner decisions 10 and 12 repeat it.
The exemption stays where A2-D7 put it — on deliberate identity replacement — and nothing else inherits
it.

**Citation tails: closed.** §3.2 is cited as 179–197 (row 2) and 192–197 (row 3); §17.17 as 2354–2369
(row 27, which now includes the A1 row the new A2 row will follow); the restart budget as 2087–2088.
Each range is exact against the frozen file.

### 3.4 Consistency across map, impact rows, summary, and owner decisions

**Pass.** The revision was checked as a network of mutually referring statements, not sentence by
sentence.

- **No orphaned amendment.** Every change named in §17.19.3 also appears in the impact table with the
  same decision attribution: `handshaking` → `incompatible` (rows 16, 23), `unavailable` → `loading`
  on an explicit same-identity load (row 16), the two health rows (row 17), the `loading` → `unloaded`
  edge (row 29), the new status fields (row 25), the two additive `§17.1`/`§17.17` records (rows 1,
  27). Nothing appears in one instrument and not the other.
- **Precedence is now self-consistent.** The map lists §17.12.1 as *confirmed*, not superseded, which
  is what the precedence clause needs in order for A2-D12's clause to be readable at all. Steps 4–6 and
  step 7 of §17.11.4 are superseded/confirmed; step 2 is deliberately absent and unimpeded.
- **No contradictory field values.** The `Available` enumeration in A2-D5, the two health rows in
  §17.19.4 and §17.19.3, and item 5 assign the same values for the same states. `Available=true` with
  `ProcessAlive=false` appears only in the budget-exhausted case, and `ServiceCapabilitiesKnown` is
  false there by A2-D6's "after any exit", so the three fields never assert incompatible things.
- **Identity language is uniform.** "Same `provider_id`/`engine`/`model_identity`" is the phrasing in
  every recovery context; "identity change" (A2-D7) is used only for the replacement path. The two are
  complements, not synonyms, and revision 6 keeps them apart.
- **Summary shorthand is backed by decision text.** Item 11's "A failed/unusable live helper uses the
  stop-and-replace path" is now precisely defined by A2-D12's "If the helper is live but `failed` or
  otherwise no longer accepts a load…". The shorthand cannot be read as overriding the in-place retry,
  because the preceding clause of the same item states the retry condition.
- **No path regresses to auto-recovery.** The one place that could reintroduce it — the update-gated
  `unavailable` → `absent` row — is given precedence explicitly in A2-D12 and item 11, and a request
  after that reset is "handled as a new attempt from `absent`, not by the retry transition".
- **Owner decisions still map cleanly.** Thirteen decisions, twelve decision entries, no reopening.
  Decisions 9 and 10 map to A2-D12 and A2-D7 respectively, and each now records the budget rule rather
  than inheriting it by adjacency. Decision 13 remains consistent with case (b).

### 3.5 Scope integrity

**Pass, with one evidence limitation.**

- Frozen contract: `4ff6b72e…e5547393`, unchanged, 2,416 lines, mtime `14:00:11` — earlier than every
  amendment revision and every audit.
- Decision report: `c25a0d3f…`, unchanged.
- Prior audits: revision 1 `258fd1c1…`, revision 2 `db6beb0a…`, revision 3 `9ceaab30…`, revision 4
  `772707bb…`, revision 5 `5944ae9e…` — all unchanged.
- Stage 2 implementation record: `ec85b2b2…`, mtime 16:59, unchanged.
- Filesystem sweep for anything modified after the revision 5 audit was written (`20:31:49` on
  2026-11-06): the amendment, and `docs/review/README.md` at `20:32:02` — which is this reviewer's own
  index entry for the revision 5 audit, part of that deliverable, not a revision-6 edit. Nothing else
  in the working tree was touched.
- No tracked source file was modified by this review. The 140 tracked modifications in the working
  tree are pre-existing project work that predates this review series and was not touched.

**Limitation, stated plainly.** `docs/review/` is entirely untracked and every revision overwrites one
path, so revision 5's bytes no longer exist and a true byte diff is impossible. "The amendment changed
only in the targeted places" is therefore established by three other means rather than by diff: the
line/byte delta is small and positive (+18 / +5,153); every load-bearing sentence verified in revision 5
is still present and unchanged in meaning; and all eighteen added lines map to the N9, N10, and E1–E4
fixes plus their traceability rows. That is strong evidence, not proof. A future revision series would
benefit from preserving each revision's bytes (for example, a per-revision copy or a hash manifest) so
that delta review can be exact.

---

## 4. Residual notes

### E5 — EDITORIAL (no action required)

§17.19.4 says "§17.7.4 gains this **additive row**; its existing rows are unchanged", and the table
beneath it has **two** rows. Two small imprecisions in one sentence: the singular, and a claim of
"unchanged" that sits two rows above the map entry narrowing line 1818's meaning. The frozen *text* is
indeed unchanged, so nothing false is asserted; a reader moving between instruments may pause. One
optional edit: "adds these two rows; no existing row's text changes (see the line 1818 note in
§17.19.3)".

### I1 — INFORMATIONAL (integration observation, not a contract defect)

Recovery from `unavailable` is triggered "only by a new explicit `load` request", but the frozen
adapter surface (§17.12.1, lines 2230–2242) has no method that sends a `load`: `Initialize` does not
load, `Synthesize`'s ensure-load is explicitly excluded while `State=unavailable`, `ProviderStatus` and
`ProviderCatalog` must not load, and `GetVoices` never spawns. Within the frozen interface, the
practical recovery routes are therefore the update-gated `unavailable` → `absent` reset after a
helper/assets/config change, and instance re-creation.

This is consistent with the amendment's scope — no caller wiring is authorized, and a future explicit
`load` surface is a Stage 4 integration question — so it is not a text defect and requires no change.
It is recorded so that adoption does not silently imply a user-reachable retry button. If the owner
intends in-session retry, that is a caller-interface decision to be taken with Stage 3/4, not a
property of this amendment.

### Persisted, previously noted (no action)

The revision 5 review's fourth immaterial nit persists: §17.19.3's §17.7.4 row ("no existing row is
rewritten") sits one row above the line-1818 correction. The two are consistent when read together, and
the correction is explicitly text-preserving, so no change is required.

### Checked and clean

- `stopping` with `Available=false` while `Loaded` may still be true — stated, not contradictory
  (A2-D5 for `Loaded`, §17.7.1's surviving false list for `Available`, A2-D6 for the snapshot).
- `incompatible` with a live or `stopping` process — A2-D4, and frozen row 1819 requires the process to
  be stopped; the exit observer's intent-awareness prevents reclassification as `crashed`.
- The six §17.7.4 rows remain mutually exclusive on `(ProcessState, State)` after the two additions.
- `degraded` remains `Ready` only on §17.7.1's definition, not on the label (A2-D5).
- `detail["missing"]` still fits the frozen `map[string]string`, canonical byte string, 31 of 256
  bytes; no wire change (A2-D3, unchanged from revision 5).
- Provider features still project through `TTSProviderCatalog`, and `Supported:false` remains confined
  to `ServiceCapabilities` (A2-D2, A2-D6, unchanged).
- B1–B11, F1–F13, and N1–N8 dispositions were re-read against their clauses while checking consistency;
  none regressed under the revision 6 edits.

---

## 5. Conclusion

Revision 6 closes N9 and N10 at every place they appeared, resolves E1–E4 at their clauses, corrects all
four citation tails, and leaves everything else byte-for-byte alone. The recovery rule is now readable
without leaning on the precedence argument, and the identity condition is present in every summary that
describes recovery.

The hash discrepancy is a record-keeping error, not an artifact error:
`4ff6b72ea03eef98e86ac158add8eb05a7e575cea22384470991f406e5547393` is the frozen architecture's digest,
the file under review reproduces it, and six independent earlier artifacts agree. The differing suffix
is in the earlier checkpoint.

**Recommendation:** adopt revision 6. Two notes (E5, I1) may be carried in the adoption record or
ignored; neither justifies another revision.

This review does not adopt the amendment, does not authorize implementation, does not authorize Stage 3,
and does not modify any artifact. Adoption and Stage 3 authorization remain separate owner decisions.

**End — revision 6 review.**
