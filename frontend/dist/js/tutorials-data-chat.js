/**
 * HOWL Chat — Tutorial Data: Chat System
 *
 * @module tutorials-data-chat
 */

'use strict';

window.TUTORIALS_DATA_CHAT = {
    id: 'chat',
    title: 'Chat',
    topics: [
        {
            id: 'chat-overview',
            title: 'Chat Overview',
            content: `# Chat System

HOWL Chat is the heart of the application. This is where you interact with AI characters in real time. Everything in the app — worlds, scenarios, characters, lorebooks — feeds into the chat to create a cohesive, immersive experience.

## What You Control

- **Messages**: Type, send, edit, delete, regenerate, and continue messages.
- **Attachments**: Attach images for vision-capable models.
- **Session Lifecycle**: Create new chats, save, load, and delete sessions.
- **Model Management**: Load and unload LLM models on demand.
- **Real-Time Settings**: Adjust sampler and model settings without restarting.

## How It Works

When you send a message, HOWL builds a **context window** from:
- The active **world** description and attached lorebooks
- The active **scenario** and its lorebooks
- All loaded **characters** with their personas and memories
- Your **chat history**
- Any **attachments** (vision images)

The model receives this context and generates a response. The response is streamed to the UI word by word.

<info>
You do not need to load a world or scenario to chat. A bare chat with no attachments works fine — but loading content makes the experience far richer.
</info>`
        },
        {
            id: 'chat-new-chat',
            title: 'Starting a New Chat',
            content: `# Starting a New Chat

A fresh chat begins with a clean slate. No history, no characters, no world context.

## How to Start

1. Click **File** in the top menu bar.
2. Select **New Chat**.
3. The backend creates a new session with a timestamped name.
4. The chat area clears, and a new session appears in the sidebar.

## What Gets Created

- A unique session ID
- An empty message history
- Default sampler and model settings (your saved preferences)
- No world, scenario, or characters attached (you add these manually)

<info>
New chats are automatically saved to the backend. If the app crashes, your chat history persists as long as the backend is running.
</info>

## Loading Content Into a Chat

After creating a new chat, load a world, scenario, or characters from the sidebar panels. The content is injected into the context window for the next message.`
        },
        {
            id: 'chat-send-message',
            title: 'Sending a Message',
            content: `# Sending a Message

Type your message in the input bar at the bottom of the chat area and press **Enter** to send.

## Keyboard Shortcuts

| Action | Key |
|--------|-----|
| Send message | <key>Enter</key> |
| New line (no send) | <key>Shift</key> + <key>Enter</key> |
| Abort generation | <key>Escape</key> or click the Stop button |

## The Send Button

The circular button to the right of the input bar sends your message. While the model is generating, this button becomes a **Stop** button. Click it to abort generation immediately.

## What Happens When You Send

1. Your message is added to the chat history.
2. The context window is assembled (world + scenario + characters + history + attachments).
3. The model begins generating a response.
4. Tokens stream to the UI in real time.
5. The full response is saved to the backend.

<info>
Generation speed depends on your hardware. GPU offloading (via Model Settings > GPU Layers) dramatically improves token throughput.
</info>`
        },
        {
            id: 'chat-message-actions',
            title: 'Message Actions',
            content: `# Message Actions

Every message in the chat has a toolbar with actions. Hover over a message bubble to reveal the buttons.

## Available Actions

- **Copy** — Copies the message text to your clipboard.
- **Edit** — Opens the message in an inline editor. You can modify the text and save or cancel. Editing a user message prunes all messages after it and regenerates the response.
- **Delete** — Removes the message from the chat. Deleting an assistant message shows a **Generate** banner so you can regenerate a replacement.
- **Regenerate** — Re-rolls the assistant response with the same context but a different seed. Useful when a response is off-target.
- **Continue** — Appends more text to the last assistant message. The model picks up where it left off.
- **Prev / Next** — Navigates through generation variations if you have regenerated a message multiple times.

## Editing Messages

When you edit a user message, HOWL prunes everything after it and re-sends the modified message to the model. This is powerful for correcting typos or rephrasing without starting over.

When you edit an assistant message, the text is updated directly without re-generation.

<warning>
Deleting a message is permanent. There is no undo. Use Prev/Next to cycle variations instead of deleting if you want to keep options.
</warning>`
        },
        {
            id: 'chat-attachments',
            title: 'Image Attachments (Vision)',
            content: `# Image Attachments (Vision)

HOWL supports vision-capable models. Attach an image to your message, and the model can "see" and describe it.

## How to Attach an Image

1. Click the **image icon** in the input bar (or drag and drop an image).
2. The image is resized to a maximum of **1024px** on the longest edge.
3. It is encoded as **JPEG** for compatibility with llama-server.
4. A thumbnail preview appears above the input bar.
5. Type your message and send.

## Supported Vision Models

Vision support depends on the model architecture. Look for models tagged with **vision** or **multimodal** on HuggingFace.

| Model Family | VRAM (FP16) | Notes |
|-------------|-------------|-------|
| LLaVA / LLaVA-NeXT | 12–16 GB | Popular open vision models |
| Qwen-VL / Qwen2-VL | 10–14 GB | Strong multilingual vision |
| InternVL2 | 12–16 GB | High-resolution vision support |
| MiniCPM-V | 6–8 GB | Lightweight vision, good for lower VRAM |

## What the Model Sees

The model receives your text prompt plus the image. It does not see previous images unless they are still attached. Each message can have **one** image attachment.

<info>
Not all GGUF models support vision. Check the model card on HuggingFace before downloading. Vision models typically have **-vision** or **-vl** in the name.
</info>

## Removing an Attachment

Click the **X** on the image preview to remove it before sending.`
        },
        {
            id: 'chat-save-load',
            title: 'Saving & Loading Chats',
            content: `# Saving & Loading Chat Sessions

Your chat sessions persist automatically in the backend, but you can also save and load them manually.

## Save Chat

**File → Save** or **File → Save As**

- Saves the current chat session to the backend storage.
- Includes: message history, loaded world/scenario/characters, sampler settings, and model settings.
- The session name defaults to the first few words of your first message.

## Load Chat

**File → Load**

1. A list of saved sessions appears.
2. Select a session.
3. The chat history, world, scenario, characters, and settings are restored.
4. The model is **not** automatically reloaded — load it manually if needed.

## Delete Chat

**File → Delete**

1. Select the session to delete.
2. Confirm the deletion.
3. The session and all its data are permanently removed.

<warning>
Deleting a chat session is irreversible. Consider creating a backup first via Settings.
</warning>

## Where Data Lives

- **Backend storage**: Primary persistence (SQLite or file-based, depending on build).
- **localStorage fallback**: If the backend is unavailable, minimal metadata is cached in the browser.`
        },
        {
            id: 'chat-model-load',
            title: 'Loading a Model',
            content: `# Loading a Model

Before chatting, you need to load a GGUF model file. HOWL uses **llama.cpp** under the hood for inference.

## How to Load

1. Click **Model Settings** in the top menu (Settings → Model Settings).
2. Click **Load Model**.
3. A file picker opens. Select a **.gguf** file.
4. The backend reads metadata (layer count, context size, architecture).
5. GPU Layers slider auto-populates with the model's layer count.
6. Adjust settings if needed, then click **Load**.

## What Gets Read

The backend extracts:
- **Architecture** (Llama, Mistral, Qwen, etc.) → auto-detects model family
- **Layer count** → sets GPU Layers max value
- **Context size** → sets default context window
- **Quantization** → displayed in model info

## Supported Formats

| Format | Supported | Notes |
|--------|-----------|-------|
| GGUF (Q4_K_M, Q5_K_M, Q8_0) | Yes | Recommended |
| GGUF (fp16 / bf16) | Yes | Best quality, largest size |
| GGML (legacy) | No | Convert to GGUF |
| PyTorch (.bin, .safetensors) | No | Convert to GGUF first |

<info>
Download models from [HuggingFace](https://huggingface.co/models?search=gguf). Search for "GGUF" to find quantized versions. TheBloke and bartowski are popular quantizers.
</info>`
        },
        {
            id: 'chat-model-unload',
            title: 'Unloading a Model',
            content: `# Unloading a Model

When you unload a model, all VRAM and RAM used by the model is freed.

## How to Unload

1. Open **Model Settings** (Settings → Model Settings).
2. Click **Unload Model**.
3. The model is removed from memory.
4. The model name display in the top bar shows **"No model loaded"**.

## When to Unload

- Before loading a different model (prevents memory conflicts).
- Before generating images (frees VRAM for the image generator).
- When closing the app (automatic on exit).
- If the model becomes unresponsive.

<info>
Unloading does not delete the model file from disk. It simply frees the loaded instance from GPU/CPU memory.
</info>`
        },
        {
            id: 'chat-supported-models',
            title: 'Supported Models & Recommendations',
            content: `# Supported Models & Hardware Recommendations

HOWL runs any GGUF model compatible with **llama.cpp**. Here are recommendations for 2026, organized by VRAM tier.

## Model Format: GGUF

GGUF is the standard format for quantized large language models. Each model is compressed using different quantization methods:

| Quant | Size | Quality | Speed | Best For |
|-------|------|---------|-------|----------|
| Q4_K_M | ~4.5 GB | Good | Fastest | Low VRAM, long context |
| Q5_K_M | ~5.5 GB | Very Good | Fast | Balanced quality/speed |
| Q8_0 | ~7.5 GB | Excellent | Moderate | Maximum quality |
| FP16 | ~15 GB | Perfect | Slowest | High-end GPUs only |

## Low-End (8–12 GB VRAM / CPU + RAM)

Perfect for laptops, older GPUs, or CPU-only setups.

| Model | Size (Q4) | Use Case |
|-------|-----------|----------|
| **Mistral 7B Instruct** | ~4.5 GB | Fast, capable general chat |
| **Llama 3.1 8B Instruct** | ~4.8 GB | Strong reasoning, coding |
| **Gemma 2 9B** | ~5.2 GB | Multilingual, safe defaults |

**CPU Mode**: Any model under your total RAM works. Expect 2–8 tokens/sec depending on CPU cores.

## Mid-Range (12–16 GB VRAM)

The sweet spot for most users. Can run larger models or higher quants.

| Model | Size (Q5) | Use Case |
|-------|-----------|----------|
| **Mixtral 8x7B (Q4)** | ~14 GB | MoE architecture, very capable |
| **Qwen2.5 14B (Q4)** | ~9 GB | Excellent reasoning, long context |
| **Llama 3.1 70B (Q4)** | ~14 GB | Best quality in this tier |

## High-End (16–32 GB VRAM)

Run the best models at high quality with room for image generation.

| Model | Size (Q8) | Use Case |
|-------|-----------|----------|
| **Llama 3.1 70B (Q8)** | ~28 GB | Near-uncompressed quality |
| **Qwen2.5 72B (Q4)** | ~18 GB | Massive reasoning capability |
| **Mistral Large 2 (Q4)** | ~16 GB | Enterprise-grade performance |

<info>
VRAM requirements scale with context length. A 32K context on a 70B model needs significantly more memory than a 4K context. Use **KV Cache Quantization** (Model Settings) to reduce memory at the cost of minor quality loss.
</info>

## Where to Download

Search [HuggingFace GGUF](https://huggingface.co/models?search=gguf&sort=downloads) for the latest models. Recommended uploaders: **bartowski**, **TheBloke**, **QuantFactory**.

## CPU + RAM Mode

HOWL fully supports CPU inference. Set **GPU Layers to 0** in Model Settings. The model loads entirely into system RAM. More CPU cores = faster generation. Enable **mmap** and **mlock** for optimal CPU performance.`
        },
        {
            id: 'character-step-out-death',
            title: 'Character Step Out & Death',
            content: `# Character Step Out & Death

Characters can be temporarily removed from the chat ("stepped out") or permanently killed ("death"). These mechanics allow for dramatic storytelling and character management without losing the character from your library.

## How to Access

1. In the chat sidebar, **right-click** a character or click the **⋮** menu button.
2. Select **Step Out** (for active characters) or **Step Back In / Revive** (for inactive characters).

## Step Out Dialog

When you select "Step Out," a dialog appears with two modes:

### Mode: Stepped Out
- The character is marked as inactive.
- They no longer participate in the chat.
- Their status changes to "stepped-out."
- A reason is recorded (e.g., "went to investigate the cave").

### Mode: Death
- The character is marked as dead.
- They no longer participate in the chat.
- Their status changes to "dead."
- A reason is recorded (e.g., "fell from the tower during the storm").
- **Automatically creates a memory entry** documenting the death.

## What Happens When You Step Out or Kill a Character

1. Character status changes to "stepped-out" or "dead"
2. Character active flag is set to false
3. stepOutReason and stepOutMode are saved
4. Bridge state is updated for panel refresh
5. Chat data is persisted (survives reloads)
6. A narrator event is queued: "[Name] has stepped out: [reason]" or "[Name] has died: [reason]"
7. For death only: A memory entry is auto-created via MemoryStore

## Death Memory Entry

When a character dies, HOWL automatically creates a memory entry with:

- **Title**: "Death of [character name]"
- **Summary**: "[character name] died on [date] at [time]. Reason: [reason]"
- **Content**: Full details including:
  - Date and time of death
  - Death reason
  - World name
  - Scenario name
  - Other characters present at the time
  - Active lorebooks
- **Scope**: character (linked to that specific character)
- **Tags**: ['death', 'event', character-name]
- **Source**: auto
- **Trigger**: death
- **Confidence**: 100

## Returning Characters (Step Back In / Revive)

### Stepped Out Characters
- Right-click the character in the sidebar.
- Select **Step Back In**.
- Enter a return reason (e.g., "returned from the cave with news").
- Character status changes to "ready" and active = true.
- stepOutReason and stepOutMode are cleared.
- Narrator event: "[Name] has stepped back in: [reason]"

### Dead Characters
- Right-click the character in the sidebar.
- Select **Revive / Resurrect**.
- Enter a revival method/reason (e.g., "resurrected by the ancient ritual").
- Character status changes to "ready" and active = true.
- stepOutReason and stepOutMode are cleared.
- Narrator event: "[Name] has been revived: [reason]"

## Editing Death Memories

To view or edit the memory entry created by a death:

1. Right-click the character in the sidebar.
2. Select **View Memories**.
3. The memory viewer opens showing all memories for that character.
4. Find the "Death of [name]" memory.
5. Use the memory viewer to edit the summary, content, tags, or pin status.

Alternatively:
1. Click **Story Tools** in the top menu.
2. Select **Memories**.
3. Filter by the character's name.
4. Edit the death memory entry.

## Important Notes

- **Characters are not deleted**: Stepping out or killing a character only affects the current chat. The character remains in your character library.
- **Persistence**: All status changes are saved to the chat data and persist across app reloads.
- **Narrator events**: Step out/death/return events are added to the chat as system messages, creating a narrative record.
- **Memory integration**: Death memories are automatically injected into future context via the memory system, so the model "remembers" the death.

<info>
Step Out is perfect for temporary character absences (scouting missions, sleep, travel). Death is permanent within the story context, but you can always revive characters if the story requires it.
</info>`
        }
    ]
};
