# Audit — Kokoro G2P Feasibility: Independent Challenge and Licensing Investigation

**Date:** 2026-11-06
**Auditor:** Muse (independent challenge pass)
**Subject artifacts:** `build/tmp/kokoro-g2p-feasibility/` (`G2P-TECHNICAL-FEASIBILITY.md`, `CANDIDATE-COMPARISON.md`, `EVIDENCE.md`, `tokencheck.go`, `evidence/`) and the prior `build/tmp/native-kokoro-spike/` results.

**Status:** INVESTIGATIVE FINDING ONLY. This document is not an implementation authorization, not a licensing opinion, and not a legal determination. It authorizes no download, no Windows spike, no dependency addition, and no production change. The controlling authorization register remains `docs/t2.9.5b-authorization-decision-register-r14.md`.

**Evidence grades used:** `[Verified locally]` reproduced in this audit on this machine; `[Source]` read from the pinned upstream repository or published terms; `[Inference]` reasoned from those two; `[Unverified]` asserted but not established by any artifact or source reviewed here.

---

## 1. Verdict

**The technical ranking survives; two of the three reasons given for it do not.**

`github.com/androiddrew/go-g2p v0.2.0` remains the strongest *technical* candidate for a native-Go Kokoro frontend. The two decisions under review are both correct:

- **HOLD on MisakiSharp** — but it is currently held for partly inaccurate reasons (see §6). The hold should be re-recorded on the actual grounds before it is cited again.
- **No implementation authorization** — correct and independently confirmed.

Three corrections are material enough that they should be recorded before any further work cites this research:

| # | Recorded claim | Audit finding |
|---|---|---|
| **C-01** | "The newer Kokoro tokenizer fixes the silent phoneme-loss problem" | **Not established.** The defect was measured in a *different module* (`chasedputnam/go-kokoro-tts v1.0.1`) than the one tagged v0.3.0 (`androiddrew/go-kokoro`). The v0.2.0→v0.3.0 change is a new `text` subpackage, not a documented tokenizer fix. No v0.3.0 artifact exists in this workspace. `[Verified locally + Source]` |
| **C-02** | "The shared ONNX Runtime design gives us a plausible way to manage the G2P and Kokoro inference runtime together" | **This is a conflict, not a design.** The two components have two competing ORT environment owners, and ORT is **mandatory** even without the neural fallback. `[Source]` |
| **C-03** | "MisakiSharp … The English-only label doesn't eliminate those questions" | **MisakiSharp is not English-only** in 2.2.0 (en, zh, ja, es, fr, hi, it, pt). The English-only scope belongs to **misaki-rs**, a different project. The hold is still correct, but for different reasons. `[Source]` |

The four-gate table in the review request is the right shape. It is missing two gates (§5), and one of the four is closer to closure than recorded (§4).

---

## 2. What this audit re-verified locally

Read-only verification of the existing investigation directory; no downloads, no builds, no dependency fetches.

| Check | Result |
|---|---|
| Artifact hashes vs `EVIDENCE.md` §5 | **All six match exactly** (`tokencheck.go`, `go.mod`, `go.sum`, `tokencheck.exe`, `evidence/build.log`, `evidence/tokencheck.log`). `[Verified locally]` |
| `tokencheck.exe` re-execution | Exit 0; output byte-identical to `evidence/tokencheck.log`, including the `ᵊ` U+1D4A drop and the `❓` U+2753 drop. `[Verified locally]` |
| Tokenizer under test | `github.com/chasedputnam/go-kokoro-tts v1.0.1`, `pkg/text` — proven by the import in `tokencheck.go` and the `require` in the feasibility `go.mod`. `[Verified locally]` |
| Module cache | Contains `go-kokoro-tts@v1.0.1` only. **Neither `go-g2p` nor `androiddrew/go-kokoro` is present**, so no offline build of either candidate is possible. `[Verified locally]` |
| Installed toolchain | `go1.27.3 windows/amd64`; `CGO_ENABLED=1` in the shell environment; mingw64 `gcc` on `PATH`. `[Verified locally]` |
| Release build invariant | `build_howl_cgo.bat` sets `CGO_ENABLED=0` and states: "Pure Go: llama.cpp runs as an external llama-server.exe, so no cgo, no MSVC/clang toolchain and no DLL bundling are required." `[Verified locally]` |
| Module Go directive | `go.mod`: `go 1.23.8`. `[Verified locally]` |

**Confirmed by this audit:** the "CGO-disabled, pure-Go release build" premise in the review request is **correct and repo-evidenced**, and it is a stronger constraint than the framing suggests, because it is documented as a *release rationale* ("no toolchain, no DLL bundling"), not merely a flag.

**New defect found while checking it:** `start-howl-chat.bat` auto-builds with `GOFLAGS=-buildvcs=false` only — it does **not** set `CGO_ENABLED=0`. The two build paths therefore disagree about the project's build characteristics. Any CGO-bearing dependency would be silently exposed on the launcher path while remaining invisible to the release script's stated rationale. This is a pre-existing inconsistency, not introduced by the G2P work, but it must be resolved before "CGO is a controlled property of the build" can be relied on in any feasibility argument. `[Verified locally]`

---

## 3. Correction C-01 — the tokenizer claim is the weakest link, and it carries the most weight

The review request credits the candidate with three advantages. The first — "the newer Kokoro tokenizer fixes the silent phoneme-loss problem identified in the earlier spike" — is load-bearing: it is the reason the `ᵊ` finding stops being a defect and becomes evidence of compatibility. It is not established.

**What is established:** the `ᵊ` drop is real, reproducible, and it is a defect in `chasedputnam/go-kokoro-tts v1.0.1`. This audit reconfirmed it. The feasibility report correctly refuses to blame `go-g2p` for it. That reasoning stands.

**What is not established:**

1. The module tagged **v0.3.0** is `github.com/androiddrew/go-kokoro` — the author's *other*, newer project. The module whose tokenizer was measured is `github.com/chasedputnam/go-kokoro-tts`. They are distinct modules with distinct licenses (Apache-2.0 vs MIT). `[Source]`
2. The published v0.2.0→v0.3.0 comparison documents a new `text` subpackage — English text preparation, sentence splitting, G2P, fallback pronunciation, chunking, streaming. **A tokenizer fix is not listed.** `[Source]`
3. `androiddrew/go-kokoro`'s README asserts an embedded Kokoro v1.0 vocabulary and that unknown phonemes fail. That is a *claim about behaviour*, not a demonstrated fix, and it was not executed against the seven-sample set. `[Unverified]`
4. No v0.3.0 source, binary, or fixture exists anywhere in this workspace. `[Verified locally]`

So the correct statement of record is: **the phoneme-loss problem has a known, cheap, unexecuted fix path, and the candidate stack plausibly already implements it.** That is materially weaker than "fixed", and it matters because a tokenizer that silently drops `ᵊ` converts a correct frontend into a subtly wrong one without any error surface — exactly the failure mode the study set out to remove.

**This is the cheapest thing in the whole program to falsify.** The existing `tokencheck.go` is 40 lines and its seven samples are pinned. Re-pointing it at the v0.3.0 tokenizer and re-running produces a binary answer for near-zero cost. Doing that requires fetching the module, which is outside this audit's authorization (§9).

---

## 4. Correction C-02 — the ORT design is an unresolved conflict, and ORT is not optional

The review request treats shared ORT as "a plausible way to manage the G2P and Kokoro inference runtime together". The underlying report is more careful — it lists ORT ownership as unresolved risk. The elevated framing should be withdrawn.

Three source facts:

1. `go-g2p` depends on `github.com/androiddrew/ortenv v0.2.0` and `github.com/yalue/onnxruntime_go v1.22.0`. `[Source]`
2. The existing `go-kokoro-tts` engine initializes **and destroys** its own ORT environment. Two owners in one process is a lifecycle conflict, not a sharing strategy. `[Source: prior spike + go-g2p source]`
3. **The POS tagger uses ORT even when the neural fallback is disabled.** ORT's native library is not bundled — the application must supply it. `[Source: go-g2p README / ASSETS.md]`

Fact 3 is the one that changes the architecture conversation. The candidate's selling point was the eSpeak-free neural fallback. But dropping the fallback does **not** drop CGO or ORT, because the tagger needs both regardless. There is no configuration of `go-g2p` that is pure-Go:

| Configuration | CGO | ORT native lib | eSpeak |
|---|---|---|---|
| Dictionary + POS + neural fallback | required | required | not required |
| Dictionary + POS, fallback disabled | **required** | **required** | not required |
| Any pure-Go variant | — | — | — (does not exist) |

Consequences that should be recorded now, before anyone reasons from "the fallback gives us an eSpeak-free route":

- Adopting `go-g2p` in-process **invalidates the documented release invariant** in `build_howl_cgo.bat` (no CGO, no toolchain, no DLL bundling). That is a direct contradiction of a stated packaging property, not a build-flag detail.
- The helper-process option is therefore more attractive than the request implies — but for a reason not yet stated: a separate module also **contains the Go directive bump** (`go 1.23.8` → `≥1.25`), keeping the main module's floor intact. The request frames the helper as a choice "we shouldn't make prematurely", which is right, but the choice is more consequential *and* more tractable than presented, because it is simultaneously a CGO, ORT-ownership, **and** toolchain-floor decision.
- ORT version is a fourth variable: tested baseline is ORT 1.22.0 / C API 22, and a second ORT version must not end up loaded beside the Kokoro engine's.

---

## 5. Revised gate table

| Gate | Recorded status | Audit status | Basis |
|---|---|---|---|
| G-1 Source-code licensing | "Favorable declarations; complete notice verification needed" | **PARTIALLY OPEN — worse than recorded.** Permissive declarations exist, but the `LICENSE` files for `go-g2p`, `MisakiSharp`, and Kokoro use the **unfilled Apache placeholder** `Copyright [yyyy] [name of copyright owner]`. No legal rights holder is named in any of them. | `[Source]` |
| G-2 Embedded dictionaries/models | "Provenance and exact asset mapping need closure" | **OPEN — and a specific commercial-source risk exists.** See §7. | `[Source]` |
| G-3 Kokoro model compatibility | "Float32 target differs from existing INT8 asset" | **NOT TESTED.** Also unresolved: whether the integration tokenizer uses the official 178-token vocabulary (C-01). | `[Inference]` |
| G-4 Windows integration feasibility | "Unestablished" | **NOT TESTED, but the premise needs correcting.** A spike is feasible on this machine (1.27.3 toolchain, mingw, CGO on). The real blockers are the release invariant, the Go directive, and ORT packaging — not toolchain availability. | `[Verified locally]` |
| **G-5 (new) ORT ownership, version, packaging** | *absent from the four-gate table* | **OPEN.** ORT is mandatory, unbundled upstream, has two competing owners, and is pinned to a tested 1.22.0 baseline. | `[Source]` |
| **G-6 (new) Reference parity and OOV quality** | *absent* | **OPEN — no measurement exists.** No eSpeak baseline, no IPA corpus, no PAR/WER, no OOV coverage figure for either candidate. The candidate's entire value proposition is "matches Misaki", and that is currently a project claim. | `[Source + Inference]` |

G-6 is not scope creep. Without it, the program can clear every legal gate and still ship a frontend whose pronunciation quality is unmeasured — and unmeasurable against the thing it claims parity with.

---

## 6. Correction C-03 — MisakiSharp is held for the wrong reasons

The HOLD is correct. The stated grounds should be replaced, because as written they will not survive contact with the maintainer.

| Recorded premise | Actual position |
|---|---|
| "LGPL-derived number conversion" | **Correct, and precisely locatable.** A `Num2Words` class in `EnglishG2P.Lexicon.cs` is documented as a behavioural port of Python num2words (LGPL-2.1): `SplitNum`, `Clean`, `Merge`, and `PyFloatPrecision` replicating Python `repr(float)`. The repo's own `NOTICE` discloses it. `[Source]` |
| "GPL-associated eSpeak behavior" | **Misclassified risk.** The package contains **no** `libespeak-ng`; P/Invoke and subprocess paths are opt-in. The GPL question attaches to the **replayed corpus** (`enus_espeak.tsv.gz`, `engb_espeak.tsv.gz`, etc.) derived from eSpeak NG 1.52.0 output. That is a data-derivative question, not a linking question — and no one, upstream or maintainer, has answered whether phoneme *output tables* are a derivative work. `[Source]` |
| "The English-only label" | **Does not apply to MisakiSharp.** 2.2.0 ships English (US/GB), Chinese, Japanese, plus es/fr/hi/it/pt replay. English-only describes **misaki-rs**. `[Source]` |
| Version | Not stated. The licensing question is filed against **2.1.1** and remains open against **2.2.0**; the number-conversion code's history between them is unconfirmed. Any request must name the version. `[Source]` |
| Rights holder | `Lyrcaxis` — a handle, declared via `<Authors>`; no legal name, no email. `LICENSE` has the same unfilled placeholder as `go-g2p`. Channels: GitHub Issues, NuGet "Contact owners", `lyrcaxis.gr`. `[Source]` |

**The decisive fact:** [issue #3](https://github.com/Lyrcaxis/MisakiSharp/issues/3) — "Licensing clarification for num2words-derived code in MisakiSharp 2.1.1", opened 2026-10-14, asking exactly the four questions we would ask — **is open with zero replies.** So the licensing question is not merely unresolved; it has been asked and is unanswered. That is a stronger justification for the hold than the request's version, and it should replace it.

**Additional engineering reason for the hold, independent of licensing:** MisakiSharp is C#/.NET (netstandard2.0). Howl is Go + Wails. Shipping it means a .NET runtime dependency — a second runtime, alongside the ONNX Runtime question, in a project whose release script advertises "no DLL bundling". Even fully cleared, its realistic role is a **reference oracle** for parity testing, not a shipped component. That conclusion holds regardless of what the maintainer answers, which makes it the more durable basis for the hold.

**Note in the report's favour:** the report's own framing that MisakiSharp's problems concern "LGPL-derived number conversion and GPL-associated eSpeak behavior" was accurate to the repo's NOTICE. The English-only gloss is the review layer's addition, and it is the part that is wrong.

---

## 7. Correction C-04 — misaki-rs carries a GPL trap the report did not record

The feasibility report calls `misaki-rs 0.8.0` "self-contained, MIT-licensed, needs no Python or ONNX when built without its default eSpeak feature." Every clause is true. The framing is incomplete in a way that would cause real harm if acted on.

- `misaki-rs` is MIT (`Copyright (c) 2026 Michele Yin`). `[Source]`
- `[features] default = ["espeak"]`, and `espeak = ["dep:espeak-rs"]`. `[Source]`
- `espeak-rs` declares MIT, but depends on `espeak-rs-sys`, which **vendors the eSpeak NG C sources and data** and whose `build.rs` runs CMake with `BUILD_SHARED_LIBS = OFF` by default — i.e. **static linkage**. Its public header grants GPLv3-or-later. `[Source]`

So the **default** `cargo build` produces a binary that statically contains GPL-3.0-or-later eSpeak NG. Wrapper-crate metadata is not clearance. Building `--no-default-features` removes it — and then unknown words become `❓`, which is not a Kokoro phoneme and which the current tokenizer silently deletes.

Worse for the eSpeak question specifically: `misaki-rs`'s README states its pronunciation dictionary was partly corrected **"using eSpeak"** and defers underlying-dictionary licensing to "the original repository". So the no-eSpeak build still rests on eSpeak-derived data. **Both candidates inherit the same unresolved data-derivative question; neither escapes it by dropping a binary.**

One item in the candidate's favour, for precision: `misaki-rs` depends on the **Rust** `num2words` crate (MIT OR Apache-2.0), not the LGPL Python library. It does **not** inherit MisakiSharp's LGPL problem. `[Source]`

---

## 8. Licensing investigation — rights holders, channels, and the three permissions

The request's addition of three separate permissions is correct and is adopted as the structure below. One further correction to the framing: **price is irrelevant to the analysis.** The operative constraint is copyleft, not cost. A free app can be proprietary; GPL-3 does not forbid charging or free distribution — it forbids added restrictions and requires source for the GPL parts. So "Howl is free" is neither a defence nor an obstacle, and it should not appear in any permission request as though it were.

### 8.1 Rights holders and contact channels

| Asset / component | Declared licence | Named rights holder | Contact channel | Status |
|---|---|---|---|---|
| `androiddrew/go-g2p` v0.2.0 (source) | Apache-2.0 + NOTICE | **None in LICENSE** (placeholder). Tag metadata: Drew Bednar, `drew@runcible.io`, tag **unsigned** | GitHub Issues; email per tag metadata. Discussions returns 404 | Permissive declaration, unattributed grant |
| `go-g2p` embedded US/UK dictionaries | Apache-2.0 via Misaki, commit `fba1236`, unchanged | hexgrad (Misaki); no per-entry origin statement | — | Origin not independently traceable |
| `go-g2p` POS tagger / tokenizer rules | spaCy `en_core_web_sm` 3.8.0, MIT, notices retained | Explosion | — | **See 8.2 — commercial training source disclosed** |
| `go-g2p` neural fallback models | Apache-2.0 per model cards (PeterReid HF, pinned revisions); converted, not retrained | PeterReid / per model card | HF repo | **No training data documented on the cards** |
| `go-g2p` digit table | Unicode 15.0 licence, notice retained | Unicode Consortium | — | Shippable with notice |
| Kokoro model weights | Apache-2.0 (model card); LICENSE is boilerplate, no named holder | hexgrad (project); architecture Li et al./`yl4579`; trained by `@rzvzn` | [hexgrad/kokoro](https://github.com/hexgrad/kokoro), [HF model page](https://huggingface.co/hexgrad/Kokoro-82M), Discord from model card | Permissive declaration |
| `kokoro-v1.0.onnx` + `voices-v1.0.bin` | Model Apache-2.0; wrapper MIT | **thewh1teagle/kokoro-onnx** (re-host/conversion) — *not* the current hexgrad repo, which ships `.pth` | [kokoro-onnx](https://github.com/thewh1teagle/kokoro-onnx) issues | **Exact asset mapping now identified** |
| Voice training audio | CC BY 3.0 (Koniwa), CC BY 4.0 (SIWIS), flagged per voice in `VOICES.md` | Original speakers/corpora | — | Attribution obligation; no separate licence grant found for the `.bin` archive |
| `chasedputnam/go-kokoro-tts` v1.0.1 | MIT | Chase Putnam (`cdputnam.com`) | GitHub | Clean, but this is the module with the tokenizer defect |
| `androiddrew/go-kokoro` v0.3.0 | Apache-2.0 | Drew Bednar, `drew@runcible.io` | GitHub Issues | Clean declaration; tokenizer fix unproven (C-01) |
| `Lyrcaxis/MisakiSharp` 2.2.0 | Apache-2.0 declared, **LGPL/GPL portions conceded in NOTICE** | `Lyrcaxis` (handle only) | Issues, NuGet contact, `lyrcaxis.gr` | **Internally inconsistent; issue #3 unanswered** |
| `MicheleYin/misaki-rs` 0.8.0 | MIT crate; **default feature statically vendors GPL-3.0-or-later eSpeak NG** | Michele Yin | GitHub | Metadata ≠ clearance (C-04) |
| `espeak-ng` | GPL-3.0-or-later | eSpeak NG project; original author Jonathan Duddington (unreachable) | espeak-ng issues | Excluded by boundary; maintainers state they cannot add a linking exception |

### 8.2 The specific commercial-source exposure in G-2

This is the finding most likely to be missed by a notice-completeness checklist, because every top-level label looks permissive.

The spaCy model notice retained by `go-g2p` (`third_party/licenses/en-core-web-sm-LICENSES_SOURCES.txt`) lists **OntoNotes 5** as **"commercial (licensed by Explosion)"**, alongside WordNet 3.0 terms. The package is MIT-labelled and the notice is complete — but a complete notice naming a commercial, licence-gated training corpus is not the same as permission to redistribute a model derived from it. `[Source]`

Whether that restricts redistribution of the *trained tagger* is a question for the rights holders, not for this audit. It should be asked explicitly, because "MIT model, notice included" is exactly the answer that closes G-2 by inspection and would be wrong here.

### 8.3 The three permissions, applied

| Component | (1) Use source | (2) Redistribute model/voice assets | (3) Bundle inside Howl's distribution |
|---|---|---|---|
| `go-g2p` code | Permitted (Apache-2.0) | n/a | Permitted with LICENSE + NOTICE |
| `go-g2p` dictionaries | Permitted | Apache-2.0 declaration, but **no per-entry origin** — cannot be asserted, only relied upon | **Unresolved** |
| `go-g2p` POS tagger | Permitted | MIT + notices, **OntoNotes 5 commercial source disclosed** | **Unresolved — ask Explosion** |
| `go-g2p` neural models | Permitted | Apache-2.0 per card, **no documented training data** | **Unresolved — provenance blank** |
| Kokoro weights | Permitted | Apache-2.0 | Permitted with licence + attribution |
| Kokoro voices `.bin` | Permitted | Apache-2.0 declaration; CC BY 3.0/4.0 training-audio credits required | **Probably permitted — attribution must ship** |
| Kokoro ONNX re-host | Permitted (MIT wrapper) | Model terms follow Apache-2.0; verify by checksum, not filename | Permitted |
| MisakiSharp | Permitted (Apache-2.0) | **Unresolved** — LGPL port + eSpeak-derived tables | **Unresolved — issue #3** |
| `misaki-rs` | Permitted (MIT) | **Unresolved** — eSpeak-corrected dictionary | **Only with `--no-default-features`, and still unresolved** |

Two operational notes for whoever drafts the requests:

- `kokoro-v1_0.onnx` should **not** be assumed identical to `kokoro-v1.0.onnx`. The located release uses the latter; identify assets by source and checksum. The existing spike's INT8 assets are already hash-pinned (`6e742170…`, `bca610b8…`), which is the right practice and should extend to the float32 pair.
- Any request to Explosion or to the Misaki/Kokoro maintainers should say plainly that Howl is **free but not open source**, and ask separately about source use, asset redistribution, and bundling. Do not let "free" imply "open source" — that is precisely the misreading the request warns against, and it is the one most likely to produce a permissive-sounding answer that does not actually cover bundling.

---

## 9. Disposition

| Item | Disposition |
|---|---|
| Technical candidate | **Go route leads** — unchanged, on the merits |
| "Tokenizer problem is fixed" | **Withdrawn to unproven.** Cheapest open item in the program |
| "Shared ORT is a plausible design" | **Withdrawn.** It is an unresolved ownership conflict; ORT is mandatory |
| MisakiSharp | **HOLD** — re-recorded on the num2words LGPL port (issue #3 unanswered), the eSpeak-derived corpus, and the .NET-runtime cost |
| `misaki-rs` fallback | **Qualified hold** — default build statically vendors GPL-3.0-or-later; `--no-default-features` leaves OOV unresolved |
| Kokoro asset redistribution | **Unresolved** — permissive declarations, unattributed grants, CC BY attribution not yet packaged |
| Windows feasibility | **Not tested** — and the real blockers are the release invariant and Go directive, not toolchain availability |
| G-5 ORT ownership/version | **New — open** |
| G-6 reference parity / OOV quality | **New — open, unmeasured** |
| Implementation authorization | **NONE.** No spike, no download, no integration, no production change |

**Proposed next steps, in ascending cost — none authorized by this document:**

1. **Re-point `tokencheck.go` at `androiddrew/go-kokoro` v0.3.0's tokenizer** and re-run the pinned seven samples. Binary answer to C-01. Requires a module fetch.
2. **Answer G-2's specific question** — put the OntoNotes 5 disclosure to Explosion and to the `go-g2p` maintainer before any notice-completeness work, since a complete notice is not a grant.
3. **Re-file the MisakiSharp question against 2.2.0** rather than relying on issue #3, naming the version and the three permissions separately.
4. **Only then** consider a Windows spike — and scope it as an ORT-ownership and packaging experiment, not a compile check, since compilation is the least risky part of G-4.
5. Fix the `start-howl-chat.bat` / `build_howl_cgo.bat` CGO disagreement before any argument depends on the build invariant.

**If permission is refused:** the fallback the request describes — link users to the official project and publish a recorded installation walkthrough — is sound, and it is cheaper than any of the above. It should be stated as a real outcome, not a concession.

---

## 10. Boundaries of this audit

**Performed:** read-only inspection of `build/tmp/kokoro-g2p-feasibility/`, `build/tmp/native-kokoro-spike/`, `go.mod`, `wails.json`, and the two build scripts; hash verification of all six evidence artifacts; re-execution of the pre-existing `tokencheck.exe`; module-cache and toolchain inspection; published-source review of upstream licences, notices, tags, issues, and release comparisons via delegated web research.

**Not performed:** no module or asset download; no `go-g2p`, `go-kokoro`, or `misaki-rs` build; no Windows spike; no model inference; no audio; no production, routing, cache, or bridge change; no commit or push.

**Limits of this document:** the licensing sections record published terms and locate the question. They are not a legal opinion, and three items are explicitly unresolved by any source reviewed: whether eSpeak phoneme *output tables* are a derivative work; whether the OntoNotes-derived tagger may be redistributed; and whether the voice archive carries a grant beyond the model card's Apache-2.0 declaration. Each requires the rights holder or counsel.

**On the prior report:** its technical work holds up. Its hashes verify, its tokenizer observation reproduces, and its central caution — that no candidate is demonstrated well enough for adoption — is correct. The corrections above are to the review layer's three reasons, not to the feasibility study's evidence.
