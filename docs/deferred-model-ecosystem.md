# Deferred · Model ecosystem and transparent connectivity

**Status:** deferred, unscheduled. Documentation only — this note authorizes no implementation,
no new branch, and no change to any current stage's acceptance criteria.

**Revision:** v1 (2026-11-02), written from the T2.5/T2.6/T2.7 audits.

**Scope:** what Howl Chat must eventually do so that selecting a model means selecting a
*complete, compatible package*, and so that every network action is explained before it
happens. Nothing here describes current behaviour unless it is labelled as an existing
limitation.

---

## 1. Dependency-aware model resolution

### 1.1 Requirement

Howl must identify the complete dependency chain for a model, not merely its primary weight
file. When a user selects a model and clicks Load, Howl resolves the package before enabling
any capability.

The resolution result must distinguish, per dependency:

| State | Meaning |
| --- | --- |
| Satisfied | Present, readable, and verified compatible |
| Missing | Required and not found in any searched location |
| Incompatible | Present but wrong architecture, version, quantisation, or variant |
| Ambiguous | Several candidates exist and none is provably the right one |
| Optional | Not required for the requested capability; would enable another one |
| Unknown | Metadata is insufficient to decide — must be surfaced, never guessed |

### 1.2 Identity

Identity is a tuple, not a filename:

- exact model and variant (base, finetune, abliterated, quantised release);
- architecture as declared by metadata, plus the format/container;
- **component role** — backbone/thinker, vision projector, audio projector, encoder,
  decoder, talker, vocoder, tokenizer, processor, text encoder, diffusion UNet, VAE;
- the capability the role enables, and the capabilities it does not.

Filenames, folder names, tags, and template contents are *hints*, never proof. A file named
`…-Omni-…` may be a thinker-only backbone; a file named `…vocoder…` may declare an
architecture string that is actually a guard message.

### 1.3 Dependency discovery

Sources, in order of authority:

1. the file's own metadata (architecture keys, tensor families, declared components);
2. a maintained architecture knowledge base shipped with Howl (component roles, required
   companions, version constraints, known-bad combinations);
3. authoritative upstream documentation and conversion tooling for that architecture;
4. user-supplied overrides.

The knowledge base must be data, not conditionals scattered through the loader, so it can be
corrected without a rebuild and reviewed like documentation.

Local search must cover: the model's own directory, all registered model directories, the
voice/TTS directory, any configured extra roots, and — only with explicit consent — removable
or network locations. Search must be bounded and must report which roots were checked.

### 1.4 Explanation

Before acquisition, Howl must be able to say, in plain language:

- what is missing and which capability it unlocks;
- why the model needs it (role in the architecture, not just "a file");
- how confident the match is, and what was checked;
- what happens if the user proceeds without it (which capabilities stay disabled).

### 1.5 Acquisition, placement, validation

Where a trustworthy, model-specific source exists, Howl may offer to fetch the component.
Acquisition must be user-authorized per operation, placed in the correct directory with a
name the runtime recognises, and validated before the capability is enabled:

- integrity (hash/size where published);
- architecture and role match;
- variant/quantisation compatibility with the backbone;
- a load-time check that the assembled package actually works together — a successful
  download is not evidence of usability.

### 1.6 Non-assumptions (normative)

Howl must not assume:

- a filename proves identity or compatibility;
- a model family supports every task in every variant;
- a model card lists every dependency;
- a Hugging Face collection proves technical compatibility;
- a successful download proves a component is usable;
- presence of a projector proves a specific modality is supported by that projector.

### 1.7 Experience target

A beginner selects a model and clicks Load; Howl resolves dependencies, explains anything
missing, offers a verified source, places and validates components, and enables exactly the
capabilities the package supports. Advanced users keep manual overrides: explicit projector
paths, role assignment, ignored warnings, and "load anyway" with the consequences stated.

### 1.8 Illustrative acceptance criteria (for the future stage, not now)

- Loading a thinker-only GGUF with no projector succeeds for text and clearly reports that
  vision and audio input remain unavailable, naming the missing component.
- Loading the same model with a valid projector enables the supported modalities and runs a
  real validation request before declaring success.
- A vocoder is never offered as a chat model, and a chat model is never offered as a vocoder.
- An ambiguous match (two candidate projectors) is never silently resolved.
- Disabling the knowledge base degrades to "unknown", not to a guess.

---

## 2. Curated Howl Chat collection on Hugging Face

### 2.1 Purpose

An official, curated collection that helps Howl recommend models, variants, companion
components, and supporting assets, and that the app can use for in-app discovery and
downloads.

### 2.2 Non-goal

It must not become a restrictive whitelist. Models from other legitimate sources remain
usable; the collection is a trusted default, not a gate.

### 2.3 Requirements

| Area | Requirement |
| --- | --- |
| Provenance | Every entry names its upstream source, conversion author, and the exact revision it was validated against |
| Compatibility metadata | Machine-readable: architecture, role, required companions, runtime minimum, quantisation, validated capabilities |
| Licensing | License and usage restrictions surfaced before download; redistribution terms respected; no relicensing by omission |
| Authenticity | Integrity data (hash/size) and, where available, signed or attested releases; source identity recorded |
| Versioning | Entries are immutable per revision; updates are new revisions with a changelog |
| Maintenance | A named owner, a review cadence, and a retirement path for stale entries |
| Staleness | Expired validation is displayed as expired; a stale compatibility claim is never presented as current |

Collection membership must never be treated as proof of technical compatibility by the
resolver; it is one input, and the local validation in §1.5 remains authoritative.

---

## 3. Offline-first operation and transparent connectivity

### 3.1 Offline guarantee

Howl remains fully usable offline for every capability whose models and dependencies are
already installed: chat, character chat, image generation, TTS/ASR with local models, and
model management. No core capability may require a network round trip.

### 3.2 Permitted online operations

Online access is optional, occasional, and authorized per operation:

- resolving model dependencies and fetching component metadata;
- downloading models or companion components;
- checking source metadata, revisions, or updates;
- other operations the user explicitly requests that inherently need the network.

No hidden or undisclosed network activity, including telemetry, analytics, "phone home"
checks, background update pings, or remote asset fetches.

### 3.3 Per-connection disclosure

Every connection must produce a clear, operation-specific explanation **before** it happens.
The disclosure must state:

1. exactly why this connection is needed, for this operation;
2. the destination service, host, and relevant source (repository/revision);
3. what request and what information will be sent — headers, identifiers, model name,
   file list, and whether any local user data or model content leaves the device;
4. what will be received — metadata, file list, specific files, sizes;
5. what Howl will do with the result once the connection succeeds;
6. security, licensing, and integrity considerations specific to this source.

A generic "this app uses the internet" notice does not satisfy this. The disclosure must
describe the actual connection being requested.

### 3.4 User control

- An explicit "Never show this again" preference per operation class, revocable.
- An accessible connection history: timestamp, operation, host, bytes, outcome, and the
  disclosure that was shown.
- A persistent network-activity indicator while a connection is in flight and after any
  recent activity, so the user can always inspect what Howl has done.
- A global offline switch that makes network operations fail closed with an explanation.

### 3.5 Security

- TLS only; certificate validation never disabled.
- Least-privilege requests: no credentials, tokens, or local paths beyond what the operation
  requires.
- Source verification and integrity validation before use.
- No silent substitution of an untrusted mirror, redirect, or "equivalent" source.
- No upload of local files or model content without explicit disclosure and authorization.
- Downloaded components are untrusted input: validated by role and structure before load.

### 3.6 Worked disclosure example (Qwen2.5-Omni projector)

> **Download the vision and audio projector for Qwen2.5-Omni-7B**
> Why: the selected model contains only the thinker backbone. The projector holds both the
> vision and audio encoders; without it, image and audio input stay disabled.
> Destination: huggingface.co, repository `unsloth/Qwen2.5-Omni-7B-GGUF`, revision pinned.
> Sent: the file request and standard HTTP headers. No local files, prompts, chats, or model
> data are sent.
> Received: `mmproj-F16.gguf`, about 2.6 GB.
> After: the file is placed beside the model, hash-checked against the published value, and
> verified by a single text+image+audio load test; only then are those capabilities enabled.
> Considerations: license terms apply; the file is third-party conversion output, not an
> official Qwen release.

---

## 4. In-app presentation

### 4.1 First-class, not buried

Dependency and connection dialogs are part of the primary experience, not settings pages or
secondary panels. They appear at the moment of the decision, in context, and they must be
readable without leaving the flow.

### 4.2 Visual language

These dialogs must look as though they have always belonged in Howl Chat. They should reuse
the application's existing theme tokens and treatment rather than introduce a new visual
system:

- surfaces and borders from the established surface/border tokens;
- the accent colour for emphasis, focus, and progress;
- the existing soft neon halo treatment around borders and cards, scaled to the dialog's
  importance rather than applied uniformly;
- the same radius, spacing, type scale, and shadow language as the surrounding panels;
- both dark and light themes, with contrast verified in each;
- real icons from the theme's icon set — never ASCII or emoji substitutes;
- motion consistent with the app's existing transitions, and suppressed under reduced-motion.

### 4.3 Creative latitude

The exact composition is deliberately unspecified. The implementer should inspect the
running application, its CSS tokens, and neighbouring dialogs, then design within that
language. "Matchy-matchy" is the goal: cohesive, professional, and native-feeling. Do not
import an external design system or force a predetermined layout.

### 4.4 Content hierarchy

1. one-line action and affected capability;
2. why it is needed (role in the architecture);
3. what leaves the device / what arrives;
4. what happens after success, and what remains unavailable;
5. risks and licence notes, expandable;
6. primary action, cancel, and the per-operation "never show again" control.

Technical detail (hosts, hashes, sizes, revisions) must be available without being required
reading; the plain-language summary comes first.

### 4.5 Accessibility and integrity

- Keyboard operable, focus trapped while modal, Escape cancels, focus returns on close.
- Screen-reader labels for every state; status changes announced.
- No dark patterns: cancelling is as easy as proceeding; "never show again" never hides the
  fact that a connection is happening, only the repeated explanation.
- Failure states explain what happened and what was not changed.

---

## 5. Evidence motivating these requirements

### 5.1 Qwen2.5-Omni thinker GGUFs (T2.6, T2.7)

Two files, `Qwen2.5-Omni-7B-Q4_K_M.gguf` (4,683,073,088 B, Unsloth attribution) and
`Qwen2.5-Omni-7B-abliterated.Q4_K_M.gguf` (4,683,072,672 B, provenance unconfirmed), both
declare `general.architecture = qwen2vl`, carry 339 tensors, and contain only the thinker's
text decoder. Their chat templates retain image, audio, and video branches, which makes them
look complete to a name- or template-based check.

The missing component is a single projector (`mmproj-…gguf`, roughly 1.5–2.6 GB) that
carries **both** the vision and audio encoders. The runtime supports that arrangement; the
package on disk simply does not include it. This is the canonical case for role-aware
resolution: correct architecture, correct family, incomplete package.

It also shows why template content is not evidence: the template describes what the
checkpoint family was trained for, not what the local files can do.

### 5.2 LFM2.5-Audio vocoder (T2.5)

`vocoder-LFM2.5-Audio-1.5B-Q4_0.gguf` (108,986,560 B) is an output-only vocoder. Its
`general.architecture` value is a guard message instructing the user to load it through a
vocoder interface; the shipped runtime exposes no such interface. It is not a speech model,
needs a compatible backbone, and cannot replace Whisper, Piper, Kokoro, or MeloTTS.

### 5.3 Scanner misclassification (observed)

Howl's generic GGUF scanner classifies by architecture and extension. `qwen2vl` is treated as
a chat architecture, so an incomplete Omni backbone is offered as a complete chat model. The
vocoder's guard-message architecture is unknown, and unknown architectures are admitted to
chat directories by default; the TTS scanner accepts only `.onnx`, so a GGUF vocoder has no
correct home at all. Extension plus architecture label is insufficient; role and companion
relationships are the missing dimension.

### 5.4 Video (out of scope, recorded for completeness)

The Omni GGUF route does not support video input, regardless of packaging. A future video
capability needs its own architecture decision and cannot be inferred from the presence of
video tokens in a template.

---

## 6. Deferred status and sequencing

- This note records a **future product requirement**. It is not authorization to implement.
- The current release-completion sequence and the outstanding audio investigation, model
  compatibility work, and audit series continue unchanged.
- No new feature branch, stage, or acceptance criterion is created by this note. Existing
  stage criteria are not modified.
- The work becomes eligible only after the current tiers close and the dependency knowledge
  base (§1.3) has an owner.
- Pickup prerequisites: a reviewed component-role taxonomy; a decision on where the
  knowledge base lives; a source-verification policy; and a design pass against the running
  application's theme, not a mock-up in isolation.

### Open questions for the future owner

1. Does the knowledge base ship with the app, update independently, or both — and how are
   corrections distributed without network dependence?
2. Which sources are trusted for automatic acquisition, and who signs off on adding one?
3. How should ambiguous matches be presented: block, or allow an explicit user choice with
   consequences?
4. What is the minimum validation that justifies enabling a modality (load only, or a real
   inference smoke test)?
5. How long is a compatibility claim valid before it is shown as stale?
