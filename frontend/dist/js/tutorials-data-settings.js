/**
 * HOWL Chat — Tutorial Data: Sampler Settings (Exhaustive)
 *
 * Every slider, number, toggle, conflict, and effect documented.
 *
 * @module tutorials-data-settings
 * @version 2.0.0
 */

'use strict';

window.TUTORIALS_DATA_SETTINGS = {
    id: 'settings',
    title: 'Settings',
    topics: [
        {
            id: 'sampler-overview',
            title: 'Sampler Settings Overview',
            content: `# Sampler Settings Overview

Sampler settings control **how the model picks the next token** during text generation. After the model computes a probability score for every possible next word, the sampler applies filters and transformations to decide which token actually gets output.

## Where to Find Them

**Settings → Sampler Settings** (topbar menu) or **File → Advanced**.

## The Sampler Pipeline (Order of Operations)

1. **Temperature scaling** — applied first to all logits
2. **Repetition penalties** — DRY, repeat penalty, encoder repeat penalty
3. **Nucleus & threshold filters** — Top-P, Top-K, Min-P, Typical-P, Top-A, Epsilon Cutoff, Eta Cutoff
4. **Smoothing** — smoothing factor + curve
5. **Dynamic temperature** (if enabled) — adjusts temperature per-token
6. **Mirostat** (if enabled) — overrides temperature and nucleus samplers entirely
7. **Token selection** — sample from the filtered distribution

<info>
All sampler settings apply **per token**. The model recalculates everything for every single word it generates.
</info>`
        },
        {
            id: 'sampler-temperature',
            title: 'Temperature',
            content: `# Temperature

**Slider**: 0.0 to 2.0 | **Step**: 0.01 | **Default**: 0.7

Temperature scales the model's confidence before it picks a token. It is applied to the raw "logit" scores the model outputs.

## How It Works

- **Low temperature** (0.0–0.3) → The model strongly favors the most likely token. Output is focused, factual, and repetitive.
- **Medium temperature** (0.5–0.9) → Balanced creativity. Good for general chat and storytelling.
- **High temperature** (1.0–2.0) → The model explores unlikely tokens. Output becomes more random, creative, and sometimes incoherent.

## What the Numbers Mean

| Value | Effect |
|-------|--------|
| **0.0** | Greedy — always picks the highest-probability token. Deterministic. |
| **0.3** | Very focused. Good for code, math, and strict instruction following. |
| **0.7** | Balanced. Default for most creative writing and chat. |
| **1.0** | No scaling. The model uses its raw probabilities. |
| **1.5** | Highly creative. Unexpected word choices, surprising dialogue. |
| **2.0** | Maximum randomness. Often produces surreal or broken text. |

## Example

Prompt: *"The wizard opened the door and saw a"*

- **Temp 0.3**: "dragon sitting on a pile of gold."
- **Temp 0.7**: "strange glowing portal swirling with purple energy."
- **Temp 1.5**: "melon made of clockwork, ticking softly in the moonlight."

## Conflicts & Warnings

- **Overridden by Mirostat**: If Mirostat Mode is enabled (> 0), temperature is ignored entirely. Mirostat dynamically controls its own temperature.
- **Interacts with Dynamic Temperature**: If Dynamic Temperature Range is enabled, the base temperature is modified per-token based on probability.

<info>
Temperature appears in both **Core Settings** and **Advanced Controls**. Both sliders control the same backend value. The Core slider (step 0.01) gives finer control.
</info>`
        },
        {
            id: 'sampler-top-p',
            title: 'Top-P (Nucleus Sampling)',
            content: `# Top-P (Nucleus Sampling)

**Toggle**: Enabled by default | **Slider**: 0.0 to 1.0 | **Step**: 0.01 | **Default**: 0.9

Top-P filters tokens by **cumulative probability**. After temperature scaling, tokens are sorted by probability. Top-P keeps the smallest set of tokens whose probabilities add up to at least P.

## How It Works

Imagine the model's next-token probabilities look like this:

| Token | Probability |
|-------|-------------|
| "the" | 0.40 |
| "a" | 0.25 |
| "his" | 0.15 |
| "her" | 0.08 |
| ... | ... |

With **Top-P = 0.9**, the sampler adds probabilities until it reaches 90%:
- "the" (0.40) + "a" (0.25) + "his" (0.15) + "her" (0.08) = 0.88 (not enough)
- Include the next token to exceed 0.90
- Discard everything after that cutoff

This means Top-P **adaptively** adjusts: when the model is confident, very few tokens are kept. When uncertain, many are kept.

## What the Numbers Mean

| Value | Effect |
|-------|--------|
| **0.0** | Only the single most likely token is considered (extremely focused). |
| **0.5** | Top ~2-4 tokens kept. Very focused but allows some variation. |
| **0.9** | Default. Good balance. ~5-15 tokens typically kept. |
| **1.0** | All tokens considered. Top-P has no effect. |

## Conflicts

- **Conflicts with Typical-P**: Enabling Top-P auto-disables Typical-P.
- **Conflicts with Top-A**: Enabling Top-P auto-disables Top-A.
- **Does NOT conflict with Top-K or Min-P**: These can be used together safely.

<warning>
Top-P and Min-P serve similar purposes but work differently. Top-P is relative (cumulative). Min-P is absolute (probability threshold). You generally do not need both.
</warning>`
        },
        {
            id: 'sampler-top-k',
            title: 'Top-K',
            content: `# Top-K

**Slider**: 1 to 100 | **Step**: 1 | **Default**: 40

Top-K is the simplest sampler: it keeps only the **K most likely tokens** and discards everything else.

## How It Works

After temperature scaling, the model has a ranked list of all possible next tokens. Top-K chops off everything below rank K.

## What the Numbers Mean

| Value | Effect |
|-------|--------|
| **1** | Greedy decoding. Only the #1 most likely token. Deterministic. |
| **10** | Very focused. Only the top 10 tokens. Highly predictable. |
| **40** | Default. Standard for chat and creative writing. |
| **80** | More permissive. Allows less likely tokens through. |
| **100** | Maximum. All 100 most likely tokens kept. |

## When to Adjust

- **Lower K** (10-20): Coding, math, factual answers, structured output.
- **Default K** (40): General chat, storytelling, roleplay.
- **Higher K** (60-100): Brainstorming, poetry, experimental writing.

## Example

Prompt: *"The color of the sky was"*

- **K=1**: "blue" (every single time)
- **K=10**: "blue", "azure", "clear", "bright", "pale"...
- **K=40**: Adds "crimson", "golden", "swirling with", "a deep shade of"...

## Conflicts & Warnings

- **No direct conflicts**. Top-K works safely with Top-P, Min-P, and all other samplers.
- **Overlap with Top-P**: If you use both Top-K and Top-P, the most restrictive filter wins. With K=40 and P=0.9, usually P is the limiting factor when the model is confident, and K is the limiting factor when it's uncertain.

<info>
Top-K is computationally cheap and a good "safety net" even when combined with Top-P. K=40 with P=0.9 is the industry-standard default.
</info>`
        },
        {
            id: 'sampler-min-p',
            title: 'Min-P',
            content: `# Min-P

**Slider**: 0.0 to 1.0 | **Step**: 0.01 | **Default**: 0.05

Min-P is an **absolute probability floor**. It removes any token whose probability is less than (min_p × highest_token_probability).

## How It Works

If the most likely token has probability 0.60 and Min-P = 0.05:
- Floor = 0.05 × 0.60 = 0.03
- Any token below 3% probability is discarded

If the model is very confident and the top token has probability 0.95:
- Floor = 0.05 × 0.95 = 0.0475
- The cutoff is higher, keeping only very likely tokens

## What the Numbers Mean

| Value | Effect |
|-------|--------|
| **0.0** | Disabled. No minimum probability filter. |
| **0.01** | Very permissive. Only extreme outliers removed. |
| **0.05** | Default. Removes genuinely unlikely tokens while preserving creative options. |
| **0.10** | Moderate filter. Keeps only fairly likely tokens. |
| **0.50** | Aggressive. Very focused output. May feel rigid. |

## When to Use

- **Min-P = 0.05**: General chat. The sweet spot.
- **Min-P = 0.00**: When using Top-P alone and you want no absolute floor.
- **Min-P = 0.10+**: Technical writing, documentation, when you want very conservative output.

## Conflicts & Warnings

- **No direct conflicts** with other samplers.
- **Overlaps with Top-P**: Both filter by probability. If Top-P is already keeping a tight set, Min-P may have little additional effect.

<info>
Min-P was introduced in 2024 and is now considered superior to Top-K as a standalone filter because it adapts to the model's confidence, just like Top-P does.
</info>`
        },
        {
            id: 'sampler-typical-p',
            title: 'Typical-P',
            content: `# Typical-P

**Toggle**: Enabled by default | **Slider**: 0.0 to 1.0 | **Step**: 0.01 | **Default**: 1.0

Typical-P filters tokens based on **information-theoretic typicality**. A token is kept if its log-probability is close to the expected entropy of the distribution.

## How It Works

Instead of keeping the highest-probability tokens (like Top-P), Typical-P keeps tokens that are "typical" — neither too predictable nor too surprising.

- **Too predictable** → "the", "a", "is" — high probability, low information
- **Typical** → "swirling", "crimson", "hesitantly" — moderate probability, good information
- **Too surprising** → nonsense words, typos — very low probability, too much information

## What the Numbers Mean

| Value | Effect |
|-------|--------|
| **0.0** | Disabled. |
| **0.5** | Keeps only fairly typical tokens. Output becomes more natural but can lose specificity. |
| **0.9** | Most typical tokens kept. Highly natural-sounding text. |
| **1.0** | Default (disabled). All tokens pass the typicality check. |

## When to Use

- **Typical-P = 0.9**: When you want very human-like, natural prose. Great for creative writing.
- **Typical-P = 1.0 (disabled)**: For most use cases. The effect is subtle and often overridden by Top-P anyway.

## Conflicts

- **Conflicts with Top-P**: Enabling Typical-P auto-disables Top-P. The app enforces this.
- **Conflicts with Top-A**: Enabling Typical-P may conflict with Top-A depending on backend version.

<warning>
Typical-P is subtle. Many users leave it at 1.0 (disabled) and rely on Top-P + Temperature instead. If you enable it, start at 0.9 and adjust down if output feels too generic.
</warning>`
        },
        {
            id: 'sampler-repeat-penalty',
            title: 'Repeat Penalty',
            content: `# Repeat Penalty

**Slider**: 1.0 to 2.0 | **Step**: 0.01 | **Default**: 1.1

Repeat Penalty reduces the probability of tokens that have already appeared in the generated text, making repetition less likely.

## How It Works

For every token in the last N tokens (controlled by **Repeat Last-N**), the model divides that token's probability by the repeat penalty value.

- Penalty = 1.1 → token is ~9% less likely to repeat
- Penalty = 1.5 → token is ~33% less likely to repeat
- Penalty = 2.0 → token is 50% less likely to repeat

## What the Numbers Mean

| Value | Effect |
|-------|--------|
| **1.0** | No penalty. Repetition is not penalized at all. |
| **1.1** | Default. Mild penalty. Catches obvious loops like "very very very." |
| **1.3** | Moderate penalty. Good for long-form writing where repetition is distracting. |
| **1.5** | Strong penalty. Can make prose feel more varied but may break natural speech patterns. |
| **2.0** | Maximum. Very aggressive. May cause the model to avoid common words entirely. |

## Example

Without penalty: *"The castle was very very very tall and the towers were very very high."*

With penalty 1.1: *"The castle was very tall and the towers were quite high."*

With penalty 1.5: *"The castle was imposing and the towers rose dramatically into the clouds."*

## Conflicts & Warnings

- **No direct conflicts** with other samplers.
- **Interacts with Presence/Frequency Penalty**: These are different mechanisms but serve similar goals. You generally do not need all three active.
- **Too high breaks flow**: At 1.8+, the model may avoid pronouns like "he", "she", "the" because they already appeared, leading to awkward writing.

<info>
Repeat Penalty checks the **generated output only**. It does not penalize tokens that appear in your prompt or in the system prompt.
</info>`
        },
        {
            id: 'sampler-repeat-last-n',
            title: 'Repeat Last-N',
            content: `# Repeat Last-N

**Slider**: 0 to 2048 | **Step**: 1 | **Default**: 64

Repeat Last-N controls **how far back** the Repeat Penalty looks for repeated tokens.

## How It Works

The model maintains a sliding window of the last N generated tokens. Any token that appears in this window has its probability reduced by the Repeat Penalty.

- **N = 64**: Checks the last ~64 tokens (about 2-3 sentences).
- **N = 512**: Checks the last ~512 tokens (about a page of text).
- **N = 0**: Disabled. Repeat Penalty has no effect regardless of its value.

## What the Numbers Mean

| Value | Effect |
|-------|--------|
| **0** | Disabled. No repetition checking. |
| **32** | Very short window. Only catches immediate repetition. |
| **64** | Default. Catches obvious loops within a few sentences. |
| **256** | Long window. Prevents thematic repetition across paragraphs. |
| **512+** | Very long window. Prevents repetition across entire responses. Slightly slower. |

## Conflicts & Warnings

- **No direct conflicts**.
- **If N = 0**: Repeat Penalty slider is ignored entirely. Set N ≥ 32 for the penalty to have any effect.

<info>
64 tokens is roughly 40-50 words of English text. This is why the default catches loops like "very very very" within the same sentence but does not penalize using "the" again two paragraphs later.
</info>`
        },
        {
            id: 'sampler-presence-penalty',
            title: 'Presence Penalty',
            content: `# Presence Penalty

**Toggle**: Enabled by default | **Slider**: 0.0 to 2.0 | **Step**: 0.1 | **Default**: 0.0

Presence Penalty penalizes tokens that have **already appeared anywhere in the generated text**, regardless of how many times. It is designed to reduce **topic drift**.

## How It Works

For every unique token that has already been generated, its probability is reduced by a fixed amount (the presence penalty value). This encourages the model to introduce new vocabulary and concepts.

- Unlike Repeat Penalty (which scales with repetition count), Presence Penalty applies **once per unique token**.
- It reduces **topic repetition** — the tendency for the model to fixate on the same subjects.

## What the Numbers Mean

| Value | Effect |
|-------|--------|
| **0.0** | Disabled. Default. |
| **0.3** | Mild. Gently nudges the model toward new topics. |
| **0.6** | Moderate. Noticeably reduces topic fixation. Good for brainstorming. |
| **1.0** | Strong. The model actively avoids words it has already used. Can feel disjointed. |
| **2.0** | Maximum. Extremely aggressive. Output may feel fragmented. |

## Conflicts

- **Conflicts with Frequency Penalty**: Enabling Presence Penalty auto-disables Frequency Penalty. The app enforces this.

<warning>
Presence Penalty penalizes ALL previously used tokens, including common words like "the", "and", "a". At values above 0.8, the model may start using unusual synonyms or sentence structures to avoid repetition, which can hurt readability.
</warning>`
        },
        {
            id: 'sampler-frequency-penalty',
            title: 'Frequency Penalty',
            content: `# Frequency Penalty

**Toggle**: Enabled by default | **Slider**: 0.0 to 2.0 | **Step**: 0.1 | **Default**: 0.0

Frequency Penalty penalizes tokens **proportionally to how many times they have already appeared** in the generated text. It targets **word-level repetition**.

## How It Works

Every time a token is generated, its probability is reduced by (frequency_penalty × count_so_far). The more a word is used, the less likely it becomes.

## What the Numbers Mean

| Value | Effect |
|-------|--------|
| **0.0** | Disabled. Default. |
| **0.3** | Mild. Repeated words gradually become less likely. |
| **0.6** | Moderate. Good for preventing word tics like "very", "really", "quite". |
| **1.0** | Strong. Aggressively diversifies vocabulary. |
| **2.0** | Maximum. Each repeated word is heavily penalized. May cause odd synonyms. |

## Conflicts

- **Conflicts with Presence Penalty**: Enabling Frequency Penalty auto-disables Presence Penalty. The app enforces this.

## Presence vs Frequency

| | Presence Penalty | Frequency Penalty |
|---|---|---|
| **Punishes** | Any token that appeared once | Tokens that appeared multiple times |
| **Strength** | Fixed per token | Scales with repetition count |
| **Best for** | Topic diversity | Word-level repetition |

<info>
You generally only need **one** of Presence or Frequency Penalty active at a time. Frequency is usually preferred for chat because it specifically targets the "very very very" problem.
</info>`
        },
        {
            id: 'sampler-no-repeat-ngram',
            title: 'No-Repeat N-Gram Size',
            content: `# No-Repeat N-Gram Size

**Slider**: 0 to 10 | **Step**: 1 | **Default**: 0

No-Repeat N-Gram **completely blocks** any sequence of N tokens that has already appeared in the generated text. This is the most aggressive anti-repetition tool.

## How It Works

If N = 3, the model tracks every 3-token sequence it generates. If it is about to output a 3-token sequence that already exists, that sequence is forbidden.

## What the Numbers Mean

| Value | Effect |
|-------|--------|
| **0** | Disabled. Default. |
| **1** | Blocks repeating individual tokens. Same as a very strong repeat penalty. |
| **2** | Blocks repeating 2-word phrases. |
| **3** | Blocks repeating 3-word phrases. Good for preventing obvious repetition loops. |
| **5** | Blocks 5-word sequences. Very aggressive. |
| **10** | Maximum. Blocks 10-token sequences. Can severely constrain output. |

## When to Use

- **N = 0 (disabled)**: Most of the time. Repeat Penalty + Frequency Penalty handle repetition well enough.
- **N = 3**: When the model is stuck in a loop. A quick fix.
- **N = 5+**: For generated poetry or song lyrics where exact repetition should never occur.

## Conflicts & Warnings

- **Can break poetry and dialogue**: Poetry intentionally repeats lines (refrains). Dialogue often reuses phrases like "I don't know." N-gram blocking will prevent these.
- **Can cause generation failure**: If N is too high and the context is long, the model may run out of valid tokens and produce garbage or stop early.

<warning>
Start with N = 0. Only increase it if you see actual repetition loops. N = 3 is usually enough to break any stuck pattern.
</warning>`
        },
        {
            id: 'sampler-epsilon-cutoff',
            title: 'Epsilon Cutoff',
            content: `# Epsilon Cutoff

**Slider**: 0.0 to 1.0 | **Step**: 0.01 | **Default**: 0.0

Epsilon Cutoff removes any token whose **absolute probability** is below a fixed threshold.

## What the Numbers Mean

| Value | Effect |
|-------|--------|
| **0.0** | Disabled. Default. |
| **0.01** | Removes tokens below 1% probability. Very mild. |
| **0.05** | Removes tokens below 5% probability. Moderate filter. |
| **0.10** | Removes tokens below 10% probability. Aggressive. Output becomes very focused. |

## Conflicts & Warnings

- **Redundant with Min-P**: Epsilon Cutoff is an absolute threshold. Min-P is also an absolute threshold but relative to the top token. Min-P is generally preferred because it adapts to the model's confidence.

<info>
Epsilon Cutoff and Eta Cutoff are both absolute filters. You rarely need both. If you use either, start with very low values (0.01–0.03).
</info>`
        },
        {
            id: 'sampler-eta-cutoff',
            title: 'Eta Cutoff',
            content: `# Eta Cutoff

**Slider**: 0.0 to 1.0 | **Step**: 0.01 | **Default**: 0.0

Eta Cutoff removes tokens whose probability is too low **relative to the second-most-likely token**. It is a relative threshold filter.

## How It Works

If the top token has probability 0.50 and the second-best has 0.30:
- Eta = 0.2 → Cutoff = 0.2 × 0.30 = 0.06
- Any token below 6% is discarded

## What the Numbers Mean

| Value | Effect |
|-------|--------|
| **0.0** | Disabled. Default. |
| **0.1** | Mild. Removes tokens far below the runner-up. |
| **0.2** | Moderate. Good for keeping output focused without being rigid. |
| **0.5** | Aggressive. Only tokens competitive with the top few are kept. |
| **1.0** | Maximum. Only the top 2 tokens are typically kept. |

<info>
Eta Cutoff was introduced as an alternative to Epsilon Cutoff. Most users find Min-P to be a more intuitive and effective relative threshold.
</info>`
        },
        {
            id: 'sampler-mirostat-mode',
            title: 'Mirostat Mode',
            content: `# Mirostat Mode

**Toggle**: Enabled by default | **Slider**: 0 to 2 | **Step**: 1 | **Default**: 0

Mirostat is an **adaptive sampler** that dynamically adjusts temperature in real time to maintain a target "surprise level" (perplexity). When enabled, it overrides Temperature, Top-P, and Typical-P.

## What the Numbers Mean

| Value | Effect |
|-------|--------|
| **0** | Disabled. Default. Temperature and Top-P work normally. |
| **1** | Basic adaptive control. Good for chat. |
| **2** | Advanced adaptive control. Better for long stories and essays. |

## Conflicts

- **Overrides Temperature**: The temperature slider is ignored when Mirostat is active.
- **Overrides Top-P**: Top-P is ignored when Mirostat is active.
- **Overrides Typical-P**: Typical-P is ignored when Mirostat is active.
- **Conflicts with Dynamic Temperature**: Enabling Mirostat auto-disables Dynamic Temperature Range.

<warning>
Mirostat is powerful but can be unpredictable. If output quality suddenly changes mid-generation, Mirostat may be over-correcting. Try lowering Eta (to 0.05) or increasing Tau (to 6.0 or 7.0).
</warning>`
        },
        {
            id: 'sampler-mirostat-tau',
            title: 'Mirostat Tau',
            content: `# Mirostat Tau

**Slider**: 0.0 to 10.0 | **Step**: 0.1 | **Default**: 5.0

Mirostat Tau is the **target perplexity** (surprise level) that Mirostat tries to maintain. Only active when Mirostat Mode is set to 1 or 2.

## What the Numbers Mean

| Value | Effect |
|-------|--------|
| **2.0** | Very low target. The model stays highly predictable. Output is focused but may feel bland. |
| **5.0** | Default. Balanced surprise level. Good for general writing. |
| **7.0** | Higher target. More creative and varied output. |
| **10.0** | Maximum. Very high surprise. May feel chaotic. |

## When to Adjust

- **Tau = 3.0-4.0**: Technical writing, documentation, code explanation.
- **Tau = 5.0-6.0**: Storytelling, chat, general creative writing.
- **Tau = 7.0-8.0**: Brainstorming, experimental fiction, poetry.

<info>
Think of Tau as "how creative do I want the model to be on average?" It is the Mirostat equivalent of setting a fixed temperature, but it adapts automatically.
</info>`
        },
        {
            id: 'sampler-mirostat-eta',
            title: 'Mirostat Eta',
            content: `# Mirostat Eta

**Slider**: 0.01 to 1.0 | **Step**: 0.01 | **Default**: 0.1

Mirostat Eta is the **learning rate** for Mirostat's temperature adjustments. It controls how aggressively Mirostat reacts to changes in text predictability.

## What the Numbers Mean

| Value | Effect |
|-------|--------|
| **0.01** | Very slow adaptation. Temperature barely changes. Mirostat feels inactive. |
| **0.05** | Slow adaptation. Good for very long, stable generations. |
| **0.1** | Default. Standard response speed. |
| **0.3** | Fast adaptation. Good for chat where topic changes often. |
| **0.5+** | Very fast. Temperature swings wildly. Can feel unstable. |
| **1.0** | Maximum. Instant full correction. Often produces oscillating quality. |

## Conflicts & Warnings

- **Only active when Mirostat Mode > 0**.
- **Oscillation risk**: At Eta > 0.3, Mirostat may over-correct. The temperature goes too high → chaotic text → Mirostat drops it → boring text → Mirostat raises it → repeat.

<info>
If you see output quality swinging between "too creative" and "too boring" within a single response, lower Eta to 0.05 or 0.1.
</info>`
        },
        {
            id: 'sampler-dyn-temp-range',
            title: 'Dynamic Temperature Range',
            content: `# Dynamic Temperature Range

**Toggle**: Enabled by default | **Slider**: 0.0 to 5.0 | **Step**: 0.1 | **Default**: 0.0

Dynamic Temperature automatically adjusts the temperature **per token** based on how confident the model is about that specific position.

## How It Works

- When the model is very confident about the next token, temperature is **lowered** → more focused.
- When the model is uncertain, temperature is **raised** → more creative.

The Range value sets the **maximum temperature boost** that can be applied.

## What the Numbers Mean

| Value | Effect |
|-------|--------|
| **0.0** | Disabled. Default. Temperature is fixed. |
| **0.5** | Mild dynamic adjustment. Temperature varies by ±0.5 around the base. |
| **1.0** | Moderate. Good for creative writing where some positions need more freedom. |
| **2.0** | Strong. High uncertainty positions get very creative. |
| **5.0** | Maximum. Extreme variation. Can feel inconsistent. |

## Conflicts

- **Conflicts with Mirostat**: Enabling Dynamic Temperature auto-disables Mirostat Mode. The app enforces this.

<warning>
Dynamic Temperature is subtle. Most users leave it at 0.0 (disabled) and use a fixed temperature instead. If you enable it, pair it with a conservative base temperature (0.5–0.6) to prevent the peaks from becoming too wild.
</warning>`
        },
        {
            id: 'sampler-dyn-temp-exponent',
            title: 'Dynamic Temperature Exponent',
            content: `# Dynamic Temperature Exponent

**Slider**: 0.1 to 5.0 | **Step**: 0.1 | **Default**: 1.0

The Exponent controls the **shape** of the Dynamic Temperature curve. Only active when Dynamic Temperature Range is greater than 0.

## What the Numbers Mean

| Value | Effect |
|-------|--------|
| **0.1** | Extreme concave curve. Almost all positions get near-maximum temperature boost. |
| **0.5** | Concave. Uncertain positions get boosted heavily. |
| **1.0** | Linear. Default. Balanced scaling. |
| **2.0** | Convex. Only very uncertain positions get significant boosts. |
| **5.0** | Extreme convex. Dynamic Temperature barely has any effect. |

<info>
Think of Exponent as the "sensitivity dial" for Dynamic Temperature. Low exponent = sensitive (reacts to small uncertainty). High exponent = insensitive (only reacts to large uncertainty).
</info>`
        },
        {
            id: 'sampler-dry',
            title: 'DRY (Do Not Repeat Yourself)',
            content: `# DRY — Do Not Repeat Yourself

DRY is a modern repetition penalty that operates on **sequences**, not individual tokens. It is more effective than traditional repeat penalties for catching multi-word patterns.

## DRY Multiplier

**Slider**: 0.0 to 5.0 | **Step**: 0.1 | **Default**: 0.0

Controls the **strength** of the DRY penalty.

- **0.0**: Disabled. Default.
- **0.5**: Mild penalty. Catches obvious repeated phrases.
- **1.0**: Moderate. Good for general use.
- **2.0+**: Strong. Aggressively prevents any repeated sequences.

## DRY Allowed Length

**Slider**: 0 to 10 | **Step**: 1 | **Default**: 2

The **minimum sequence length** before DRY applies.

- **2**: DRY applies to 2-token sequences (word pairs). Default. Good balance.
- **3+**: DRY applies to 3+ token sequences. Very targeted. Only catches multi-word repetition.

## DRY Base

**Slider**: 0.0 to 5.0 | **Step**: 0.1 | **Default**: 1.0

The **base penalty value** before the multiplier is applied.

- **1.0**: Standard. The penalty is exactly (multiplier × base).

## When to Use DRY

- **DRY is disabled by default**. Start with traditional Repeat Penalty (1.1-1.3).
- Enable DRY if you see **multi-word repetition** that Repeat Penalty doesn't catch.
- **Recommended starting point**: Multiplier = 0.5, Allowed Length = 2, Base = 1.0.

<info>
DRY was added to llama.cpp in late 2024 and has become the preferred repetition control for long-form creative writing because it catches phrases, not just words.
</info>`
        },
        {
            id: 'sampler-smoothing',
            title: 'Smoothing Factor & Curve',
            content: `# Smoothing Factor & Curve

Smoothing applies a mathematical function to the probability distribution before sampling, making it more or less "peaked."

## Smoothing Factor

**Slider**: 0.0 to 5.0 | **Step**: 0.1 | **Default**: 0.0

- **0.0**: No smoothing. Default. The raw distribution is used.
- **0.5**: Mild smoothing. Less likely tokens get a small boost.
- **1.0**: Moderate. The distribution is noticeably flattened.
- **2.0+**: Strong smoothing. Highly unlikely tokens become much more probable.

## Smoothing Curve

**Slider**: 0.1 to 5.0 | **Step**: 0.1 | **Default**: 1.0

- **1.0**: Linear. Default. Even scaling across all probabilities.
- **0.5**: Concave. The effect is concentrated on the lowest-probability tokens.
- **2.0**: Convex. The effect is concentrated on mid-probability tokens.

## Conflicts & Warnings

- **At high values, output can become incoherent**: Factor > 2.0 can make the model choose genuinely bad tokens because they get boosted too much.

<info>
Smoothing is an alternative to raising temperature. Temperature scales logits before the softmax. Smoothing modifies the distribution after the softmax. They achieve similar "creativity" effects through different math.
</info>`
        },
        {
            id: 'sampler-top-a',
            title: 'Top-A',
            content: `# Top-A

**Toggle**: Enabled by default | **Slider**: 0.0 to 1.0 | **Step**: 0.01 | **Default**: 0.0

Top-A is an alternative to Top-P that uses an **absolute probability threshold** based on the most likely token.

## How It Works

A token is kept if its probability satisfies: token_prob >= top_a × (max_prob)^2

This creates a **quadratic cutoff** that is sharper than Min-P's linear cutoff.

## What the Numbers Mean

| Value | Effect |
|-------|--------|
| **0.0** | Disabled. Default. |
| **0.1** | Mild. Only tokens very close to the top are kept. |
| **0.2** | Moderate. More restrictive than Min-P = 0.05 would be. |
| **0.5** | Aggressive. Very focused output. |

## Conflicts

- **Conflicts with Top-P**: Enabling Top-A auto-disables Top-P.

<warning>
Top-A is largely superseded by Min-P in modern configurations. Most users disable Top-A and use Min-P = 0.05 instead, as Min-P is more intuitive to tune.
</warning>`
        },
        {
            id: 'sampler-encoder-repeat-penalty',
            title: 'Encoder Repetition Penalty',
            content: `# Encoder Repetition Penalty

**Slider**: 1.0 to 2.0 | **Step**: 0.01 | **Default**: 1.0

Encoder Repetition Penalty penalizes tokens that appear in the **input prompt** (the "encoder" side), not just the generated output.

## What the Numbers Mean

| Value | Effect |
|-------|--------|
| **1.0** | No penalty. Default. |
| **1.1** | Mild. Tokens from your prompt are 9% less likely to appear in the response. |
| **1.3** | Moderate. Good for preventing the model from echoing your phrasing. |
| **1.5** | Strong. The model actively avoids words you used. |
| **2.0** | Maximum. Very aggressive. Can break coherence. |

## When to Use

- **1.0 (disabled)**: Default. Most of the time this is fine.
- **1.1–1.2**: If the model frequently echoes your exact phrasing or parrots back your questions.
- **1.3+**: If the model feels like it is "agreeing with you" too much by reusing your terminology.

## Conflicts & Warnings

- **Too high breaks instruction following**: At 1.5+, the model may avoid key terms from your instructions, making it harder to follow complex prompts.

<info>
Encoder Repeat Penalty is especially useful in chat where the model might otherwise start every response with "As you mentioned..." or repeat phrases from your message.
</info>`
        },
        {
            id: 'sampler-seed',
            title: 'Seed',
            content: `# Seed

**Input**: 0 to 4294967295 | **Step**: 1 | **Default**: 0

The Seed controls the **random number generator** used during token sampling. With the same seed, same prompt, and same settings, the model produces **identical output**.

## What the Values Mean

| Value | Effect |
|-------|--------|
| **0** | Default. The backend uses a random seed each time. Output is not reproducible. |
| **1–4294967295** | Fixed seed. Same seed + same settings = same output. |

## When to Use

- **Seed = 0**: Default. For normal chat and creative writing where you want variety.
- **Fixed seed**: For testing sampler settings. Change one slider, keep the seed fixed, and compare outputs.
- **Fixed seed**: For sharing reproducible prompts.

## Important Notes

- **Seed only affects randomness**. If Temperature = 0.0 (greedy), the seed has no effect because there is no random choice.
- **Same settings required**: Changing ANY sampler setting (even by 0.01) will produce different output even with the same seed.
- **0 is not a special seed**: In this UI, 0 means "random." The backend generates a fresh random seed each time.

<info>
The maximum seed value (4294967295) is the limit of a 32-bit unsigned integer. You can enter any whole number between 0 and this value.
</info>`
        }
    ]
};
