# HOWL Chat — Model / Prompt-Format Compatibility Matrix

**Status:** Research and compatibility assessment only. **No implementation authorization. No architecture selection.**
**Date:** 2026-11-12
**Why this document exists:** A memory comparison is invalid if the prompt format differs between conditions (`P10`). The work order also requires a separate investigation of model-specific prompt formats, because supporting a template in llama.cpp does not mean Howl constructs every required prompt correctly. This document records the format landscape, Howl's current position in it, the gaps, and the tests that would close them. It is separate from `MEMORY-01` because a 17-family compatibility matrix with per-family confidence would swamp a general memory synthesis.
**Companions:** `MEMORY-02-SOURCE-REGISTER-2026-11-12.md` (S28–S37), `MEMORY-03-EVALUATION-PROTOCOL-2026-11-12.md` (Stage G, `F0`–`F2`).
**Controlling records:** `docs/architecture/DESIGN-MEMORY-CONTINUITY-2026-11-11.md`; `app-prompt.go`, `app-llama.go`, `app-family-defaults.go` (read-only inspection).

**Authority:** Research evidence. It does not authorize changing the prompt path.

---

## 1. Why format fidelity is a memory problem

Four distinct failure paths connect prompt format to memory quality:

1. **Silent degradation.** A model given a nearly-correct template still answers, but worse. Nothing errors, so the memory condition looks weak when the format is wrong.
2. **System-prompt loss.** Several families do not support a system role (Gemma), or reject a system message after the first turn. Memory state injected as a system message is then dropped or rejected.
3. **Turn-boundary corruption.** Memory or lore inserted at the wrong depth can break `<|eot_id|>`/`<end_of_turn>` boundaries, merging two turns or truncating a reply.
4. **Cache and reasoning interaction.** Reasoning models emit a thinking channel; a memory block placed inside the thinking region is invisible to the answer, and reasoning tokens consume the context budget the memory design assumed was free.

**Consequence for evaluation:** every memory condition runs on the model's own template (`F0`). `F1` (Howl's own builder) and `F2` (deliberate mismatch) are separate conditions, so a format defect cannot be misread as a memory defect.

---

## 2. Verified facts about Howl's current prompt path

Source-inspected, no runtime claim.

| Finding | Evidence | Consequence |
|---|---|---|
| The live path sends a **message array** to llama-server; the server applies a template | `app-prompt.go` (`buildPromptMessages`) | Format correctness depends on the server's template engine and the template it selects |
| A legacy builder emits **hardcoded ChatML** (`<\|im_start\|>…<\|im_end\|>`) with a ```` ``` ```` fence around the history and the rule text inside the system block | `app-prompt.go` (`buildPrompt`) | If reachable, it is wrong for every non-ChatML family and mixes a code fence into the turn structure |
| Family **sampler** defaults exist for `qwen3`, `gemma`, `deepseek`, with a generic fallback | `app-family-defaults.go` | Sampling is family-aware; **format is not** — there are no per-family template builders |
| `custom_jinja_template` is passed as `--chat-template <string>` | `app-llama.go` | A custom template can override the model's own |
| **No `--jinja` flag is passed** | `app-llama.go` (full argument list inspected) | Upstream server documentation says `--chat-template` accepts a full jinja template only when `--jinja` is set before it; otherwise only built-in template names are accepted (S28). This is a **probable defect requiring verification**, not a confirmed runtime failure |
| `prompt_template` is set to `"auto"` | `config/settings.json` | "auto" must map to a deterministic choice; an unmapped "auto" is a silent default |
| Truncation is `contextSize × 4` **characters**, keeping system + newest | `app-prompt.go` | Token accounting is approximate; a 20,480-token window is truncated at ~81,920 characters, which is not a token count |
| A comment notes that several templates reject a system message after the first turn | `app-prompt.go` | The authors already recognised template constraints; the handling is incomplete |
| Reasoning flags: `--reasoning on/off/auto`, `--reasoning-format deepseek`, `--reasoning-budget`, `--reasoning-budget-message` | `app-llama.go` | Reasoning is wired, but `--reasoning-format deepseek` is DeepSeek-specific and may not fit every reasoning family |

**Most important finding:** Howl has family-aware **sampling** and generic **formatting**. A memory design that injects structured state assumes the format layer will carry it; that assumption is currently untested and, in the legacy builder, demonstrably wrong for non-ChatML families.

---

## 3. Compatibility matrix

Legend — **Official template:** whether the model family publishes an authoritative chat template. **Backend:** llama.cpp support (built-in name, or jinja from model metadata). **Howl:** what Howl does today. **Confidence:** in the *format description*, not in Howl's behaviour.

| # | Format / family | Official template | Backend support | Howl support today | Missing implementation | Known incompatibilities | Required tests | Confidence | Source |
|---|---|---|---|---|---|---|---|---|---|
| 1 | **ChatML** (generic) | Yes (widely used convention) | Built-in `chatml` | Legacy builder emits it; live path relies on server | Per-family selection; verify the fence/turn bug is unreachable | Not native for Llama/Gemma/Phi/Command-R; system-role assumptions vary | Emitted-sequence diff vs a pinned ChatML model | **High** | S29 |
| 2 | **Llama 2** | Yes (`[INST]`, `<<SYS>>`) | Built-in `llama2` | Server template only | No builder; no system-role mapping (`<<SYS>>` must be inside first `[INST]`) | Two system messages, or a system turn after the first, are outside the format | System-in-first-turn test; multi-turn diff | **High** | S29 |
| 3 | **Llama 3 / 3.1 / 3.2 / 3.3** | Yes (`<\|start_header_id\|>`, `<\|end_header_id\|>`, `<\|eot_id\|>`) | Built-in `llama3` | Server template only | No builder; tool/`ipython` role unhandled | A memory block inserted mid-turn breaks `<\|eot_id\|>` boundaries; tool results need the `ipython` header | Turn-boundary insertion test; tool round-trip | **High** | S29, S30 |
| 4 | **Mistral / Mixtral** (v0.1–v0.3) | Yes (`[INST]…[/INST]`) | Built-in `llama2`-family/mistral handling | Server template only | No builder | v0.1 has no system role; consecutive user turns are not representable without merging | System-role merge test; consecutive-user-turn test | **Medium–High** | S29 |
| 5 | **Qwen2 / 2.5** | Yes (ChatML variant) | Built-in `chatml`-compatible | Server template; sampler default | No format builder | ChatML-shaped but not identical to generic ChatML in every version | Emitted-sequence diff | **High** | S36 |
| 6 | **Qwen3** | Yes (ChatML + `enable_thinking`) | Jinja from model metadata; **the hard thinking switch is not exposed in llama.cpp** (Qwen docs recommend a custom template file) | Sampler default only; no thinking toggle in Howl | Thinking/non-thinking selection per request; `/no_think` soft switch; reasoning budget interaction | With thinking on, the reasoning block consumes budget the memory design counted as free; disabled thinking must be primed by closing the block in the prompt | Thinking-on/off token accounting; assistant-prefix correctness; budget interaction | **High** on markers, **Medium** on toggle | S36 |
| 7 | **DeepSeek V3 / V3.1 / R1** | Yes (fullwidth-pipe envelope; dedicated tool tokens) | Jinja from metadata; tool templates shipped by vLLM/llama.cpp (S30) | Sampler default; `--reasoning-format deepseek` | No builder; **no system token exists** — system text must be folded into the first user turn | Treating DeepSeek as ChatML inserts invalid `<\|im_start\|>` markers; reasoning-format flag is DeepSeek-specific | Emitted-sequence diff; system-fold test; reasoning channel extraction | **High** | S34, S30 |
| 8 | **DeepSeek-R1-Distill (Llama/Qwen base)** | Inherits the **base** family's template, not DeepSeek's | Jinja example shipped in llama.cpp (S30) | Sampler default keyed on "deepseek" | Distinguish distill (Llama/Qwen template) from native DeepSeek (pipe envelope) | Applying the DeepSeek envelope to a Qwen-based distill is wrong | Family-detection test: distill vs native | **Medium–High** | S30, S34 |
| 9 | **Gemma 2 / 3** | Yes (`<start_of_turn>` / `<end_of_turn>`) | Built-in `gemma` | Sampler default; server template | **No system role**: system instructions must be placed in the first user turn | A separate system message is outside the format; a memory block as a system turn is dropped or misread | System-folding test; multi-turn diff | **High** | S33 |
| 10 | **Command-R / R+** | Yes (`<\|START_OF_TURN_TOKEN\|>`, `<\|SYSTEM_TOKEN\|>`, `<\|USER_TOKEN\|>`, `<\|CHATBOT_TOKEN\|>`, `<\|END_OF_TURN_TOKEN\|>`) | Built-in `command-r` | Server template only | Preamble structure (safety preamble, system preamble, user preamble, tools) is template-defined; Howl has no preamble model | Dropping the preamble markers changes behaviour; tool section layout is prescribed | Preamble-presence test; tool section diff | **High** | S35, S29 |
| 11 | **Phi-3 / Phi-4** | Yes (`<\|system\|>`, `<\|user\|>`, `<\|assistant\|>`, `<\|end\|>`) | Built-in `phi3`/`phi4` | Server template only | No builder; tool format is model-specific | ChatML markers are invalid; `<\|end\|>` is the turn terminator | Emitted-sequence diff | **Medium–High** | S29, S30 |
| 12 | **GLM-4 / ChatGLM** | Yes (`[gMASK]<sop>`, `<\|user\|>`, `<\|assistant\|>`, `<\|observation\|>`) | Built-in `chatglm3`, `chatglm4` | Server template only | No builder; `[gMASK]`/`<sop>` framing unhandled | Omitting `[gMASK]<sop>` invalidates the prompt | Emitted-sequence diff including the prefix | **Medium** | S29 |
| 13 | **InternLM 2 / 2.5** | Yes (ChatML-like, plus action tokens) | Jinja from metadata | Server template only | Action/tool token handling | Assumed ChatML equivalence is not verified for every version | Emitted-sequence diff | **Medium–Low** | S29 (built-in list) |
| 14 | **Yi** | Yes (ChatML-like) | Jinja from metadata | Server template only | Version-specific differences unverified | Assumed ChatML equivalence is not verified | Emitted-sequence diff | **Low–Medium** | S29 |
| 15 | **Cohere open-weight (Aya / Command lineage)** | Model-card template | Via `command-r`/jinja | Server template only | Version mapping | Not every Cohere open-weight model uses the Command-R envelope | Per-model card check | **Low–Medium** | S35 |
| 16 | **Vision / multimodal** | Content-array messages with image parts | llama.cpp multimodal path (`mtmd`); server accepts content arrays | **Not supported** — Howl's path is text-only | Image part structure, placeholder tokens, per-family image tokenisation | Text-only templates reject or silently drop image parts | Vision round-trip on one vision model | **Medium** | S28 |
| 17 | **Tool / function calling** | Model-specific; generic fallback is a known compromise | `--jinja` + per-model tool template; generic ChatML fallback "YMMV" (S30) | **Not implemented** in Howl's prompt path | Tool schema advertisement, tool-call parsing, tool-result role | Generic fallback costs tokens and is unreliable; per-family templates are required for native behaviour | Tool round-trip per family; parse-failure rate | **High** on risk | S30 |
| 18 | **Reasoning channels** | Per-family (DeepSeek-style, Qwen-style, others) | `--reasoning on/off/auto`, `--reasoning-format`, `--reasoning-budget` | Flags passed with `deepseek` format | Per-family reasoning format selection; separate reasoning token budget in memory accounting | A single `deepseek` format does not fit every reasoning family; reasoning tokens consume context | Reasoning extraction test; budget accounting | **Medium–High** | S28 |

### 3.1 Confidence caveats

- Rows 13–15 are marked lower because the exact per-version markers were not verified against each model card. They are **not** claimed as exact.
- Row 6's thinking toggle is documented as **not exposed** by llama.cpp; the Qwen documentation recommends a custom template file to force non-thinking. Treat any claim that Howl "supports Qwen3 thinking control" as unverified until tested.
- Row 7's fullwidth-pipe markers were verified from a technical wire-format reference and the model card; they are described as token-level markers, and any implementation must read the exact bytes from the tokenizer rather than transcoding them.
- **No exact marker string in this table should be copied into code without reading the pinned model's own template.** The table is a map, not a source of truth for bytes.

---

## 4. Howl's gap list

| Gap | Severity | Why it matters for memory |
|---|---|---|
| No per-family **format builder** | High | Structured state/provenance blocks are injected as if every family were ChatML |
| `--chat-template` used without `--jinja` | High (unverified) | A custom template may not be applied as intended; the server may fall back |
| `prompt_template: "auto"` mapping undefined | Medium | Two runs can differ with no recorded reason |
| Legacy hardcoded-ChatML builder still present | High | A reachable legacy path is wrong for most families; the code fence is a structural defect |
| Character-based truncation | Medium | Token budget claims are estimates, not measurements |
| No system-role folding for Gemma/Llama 2/DeepSeek | High | Memory injected as a system message is lost or invalid |
| No thinking-budget accounting | Medium | Reasoning tokens silently consume the memory budget |
| No tool role / tool result handling | Medium | Future memory tools would be malformed |
| No vision content parts | Low (for text memory) | Blocks multimodal characters later |
| No emitted-vs-expected sequence check | High | Format regressions are invisible |

---

## 5. Required tests (Stage G)

| Test | Method | Pass condition |
|---|---|---|
| `T-FMT-1` Emitted-sequence diff | Render a fixed 6-message conversation through Howl and through the model's own template; compare token IDs | Byte/ID-identical for `F0`; documented difference for `F1` |
| `T-FMT-2` System-role folding | Send a system message to Gemma, Llama 2, DeepSeek; inspect rendered prompt | System content present exactly once, in a valid position |
| `T-FMT-3` Turn-boundary insertion | Insert a memory block at each declared depth; inspect boundaries | No merged or truncated turns |
| `T-FMT-4` Thinking toggle | Qwen3 with thinking on/off; count reasoning tokens | Toggle has the documented effect, or is reported unsupported |
| `T-FMT-5` Reasoning channel | DeepSeek-style reasoning model; confirm channel separation | Reasoning not mixed into visible answer; budget counted |
| `T-FMT-6` Prompt-cache interaction | Same prompt with/without a changed memory block | Cache-hit behaviour recorded, not assumed |
| `T-FMT-7` Tool round-trip | One tool-capable family; advertise, call, return | Parse success; no malformed markers |
| `T-FMT-8` Negative control `F2` | Deliberate ChatML for a Gemma/Command-R model | Measurable penalty; confirms the harness is sensitive |
| `T-FMT-9` Truncation reality | Fill the window; compare char estimate to tokenizer count | Error within a declared tolerance, or truncation replaced by token counting |
| `T-FMT-10` Multi-template model | A model shipping more than one template (e.g. thinking/non-thinking) | The selected template is recorded and reproducible |

---

## 6. What "auto" must mean

`prompt_template: "auto"` is only acceptable if it resolves deterministically and is recorded. The minimum contract:

1. If a model file carries `tokenizer.chat_template`, use it (with `--jinja`).
2. Else, map the detected family to a named built-in template and record the name.
3. Else, fail visibly with a "template unresolved" error rather than defaulting silently to ChatML.
4. Record the resolved template's hash in the run manifest.
5. Never infer the template from the sampler family name alone; sampler defaults and template identity are separate.

**A silent ChatML fallback is the single most likely cause of a false "model X is bad at memory" result.**

---

## 7. Unresolved questions

| ID | Question | Status |
|---|---|---|
| `Q-P1` | Does `--chat-template` without `--jinja` actually apply a custom template in the pinned backend version? | Unverified; must be tested |
| `Q-P2` | Which families can accept a system message after the first turn? | Partially known; needs per-family test |
| `Q-P3` | How should Howl fold system content for Gemma, Llama 2, and DeepSeek? | Design open |
| `Q-P4` | Does the legacy `buildPrompt` path remain reachable? | Source suggests legacy; reachability unproven |
| `Q-P5` | What is the token-accurate truncation replacement? | Engineering decision, not authorized |
| `Q-P6` | Which reasoning format should non-DeepSeek reasoning models use? | Open |
| `Q-P7` | Do any target families require a BOS/EOS post-processing step Howl does not perform? | Unknown |
| `Q-P8` | How much does `F2` cost in correctness? | Unmeasured; Stage G |
| `Q-P9` | Should memory blocks be inserted at a fixed depth or family-specific depth? | Depends on `Q-P2` and position evidence (S19) |
| `Q-P10` | What is the minimum viable tool-call contract for future memory tools? | Future design |

---

## 8. What this document does not do

- It does not authorize editing `app-prompt.go`, `app-llama.go`, or any settings file.
- It does not assert that any family is broken at runtime — only that the format description, the backend capability, and Howl's implementation are three different things and Howl's layer is the least proven.
- It does not certify any exact marker string for production use; the pinned model's template is the authority.
