# Image failure diagnosis

Date: 2026-11-26
Isolated root: `C:\HOWL-Test\hra-image-20261125` (`.howl-isolated-root`)

M5 established that HiDream is blocked by architecture identification, not by
memory. Every engine failure used to reach the user as one undifferentiated
string, so that distinction existed only in the log. This slice makes the
category part of what the user and any later reader see.

## Categories

`internal/image/diagnosis` owns one classification per failure:

| Category | Meaning | Next step it gives |
|---|---|---|
| `unsupported_architecture` | The engine does not know the architecture (Safetensors/checkpoint input) | Use a supported model file or update the engine; memory will not help |
| `unsupported_weight_format` | The architecture may be supported, but this file's format or quantized variant is not | Use the original Safetensors release for this family, or update the engine |
| `invalid_model_file` | Missing, truncated, unparseable, or unsupported format | Check the file, re-scan the folder |
| `insufficient_resources` | Real memory shortfall against the post-reserve budget | Smaller/more quantized model, lower resolution, unload another workload |
| `placement_failure` | Backend / params-backend / flag rejected | Reset placement to Auto and retry |
| `engine_startup_failure` | The process never became usable | Check engine binary and model folder; the engine was restarted |
| `generation_failure` | Engine accepted the job and the job failed | Retry; check the engine log if it repeats |

Classification reads the engine's own words (the 400 body, process-exit text,
admission reason) and keeps those words attached as `engine said: …`. Nothing
is discarded, and a diagnosis is never re-labelled on its way up: `Wrap` is
idempotent and preserves the stage of the first classification.

Format and architecture share one engine message ("cannot identify model family
from '<path>'"), so the file the engine named decides the category: `.gguf`,
`.pt`/`.pth` and `.bin` are format/variant rejections, `.safetensors`/`.sft`/
`.ckpt` are architecture rejections. The category is a statement about the file
the engine was handed, never about a family in general.

Extension is a practical discriminator for the failures observed here, not proof
that every file with that extension is incompatible: a malformed Safetensors
file can fail for its own reasons and would be reported as an architecture
rejection by this rule. `invalid_model_file` covers the parse/IO failures the
engine names explicitly; the extension rule only resolves the generic
"cannot identify" message.

`retryable` is set per category so a UI does not offer "try again" for an
architecture the engine will never identify.

## Wiring

- `internal/image/job` wraps the prepare failure as `startup`, a dead-engine
  health check as `generation`, and the generation result as `generation`, so
  `job.Error` carries the category. It also emits
  `[ImageDiagnosis] category=… stage=… retryable=… raw=…`.
- The synchronous `GenerateImage` path returns the same classified error.
- Admission refusals keep their existing detailed text; when they reach a job
  they classify as `insufficient_resources`.

## Evidence (browser dev host, real engine)

HiDream selection, bounded placement admitted and started inside the cap, then
the engine refused the GGUF file. Before this slice the job error was the raw
400. With the first classification it read `unsupported_architecture`; after the
correction below it reads the narrower category:

```
[ImageDiagnosis] category=unsupported_weight_format stage=generation retryable=false
  raw="… sd-server txt2img error 400: {\"error\":\"cannot identify model family from
  'H:\\\\AI GGUF\\\\Image models\\\\hidream-i1-fast-Q2_K.gguf'\"}"

job.error = [unsupported_architecture] This engine build cannot identify the
model's architecture, so it cannot run it. Choose a model file this engine
supports, or update the engine; changing VRAM settings will not help.
(engine said: … cannot identify model family from '…' …)
```

The resource case is covered by unit tests over the admission text
(`needs 11761 MB VRAM; 6233 MB available` → `insufficient_resources`), and the
placement/dead-engine cases by the same table.

## Owners

- `diagnosis` classifies; it decides nothing and retries nothing.
- The job manager and `GenerateImage` remain the only places a failure becomes
  job state or a returned error.
- Admission keeps its own detailed refusal text; classification wraps it rather
  than replacing it. No second decision path was added.

## Verified limitations

- **Classification is textual.** It reads stable engine phrases; a new engine
  build with different wording falls back to `generation_failure` /
  `engine_startup_failure` rather than guessing.
- **Only the image path is classified.** Chat and audio failures still surface
  raw; the same treatment is a separate slice.
- **The packaged-host image UI was not clicked**; evidence is from the browser
  dev host over the same bindings.
- **The UI does not yet render categories specially.** The category is present
  in the job error and the log; a tailored panel is not built.
