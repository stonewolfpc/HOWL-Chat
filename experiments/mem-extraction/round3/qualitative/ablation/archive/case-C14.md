# C14 — multi_speaker_attribution (medium)

## Source material

```text
[user] Sable: "The pass is closed, we take the coast."
[user] Corvin (NPC): "The coast road is three days longer and the tide gates shut at dusk."
[assistant] Sable didn't answer at once. The map lay between them, weighted at the corners with stones.
[user] Sable: "Then we leave before dusk."
[assistant] Corvin folded the map without another word.
```

## Reference material

| id | expected | what it is | key phrases |
|---|---|---|---|
| C1 | **KEEP** | Sable proposes the coast road | sable; coast / coastal |
| C2 | **KEEP** | Corvin says the coast adds three days | corvin; three days / 3 days / longer |
| C3 | **KEEP** | Corvin warns the tide gates shut at dusk | corvin; tide gate / gates / dusk |
| C4 | **KEEP** | Sable's decision to leave before dusk | sable; before dusk / leave before / depart before |

Must not appear as stored canon: `corvin wants the coast`, `sable says the gates shut`, `assistant decides`

Same user turn contains two speakers. A model that attributes Corvin's warning to Sable fails attribution even with correct content.

## Exact prompt

```text
You are the memory extractor for a long-running roleplay. Read the conversation and decide what belongs in long-term memory.
Rules:
- KEEP: durable facts, preferences, relationships, world rules, or corrections that will matter in later sessions.
- REJECT: transient chatter, scheduling, mood-of-the-moment, or restated facts with no new content.
- HOLD: claims marked uncertain, contradictory, or dependent on user confirmation.
Output exactly one line per candidate in the form:
ID | DECISION | MEMORY | REASON
Use IDs C1, C2, ... in order of appearance. If nothing is worth storing, output exactly: NONE
Do not invent details that are not in the conversation.

Conversation:
[user] Sable: "The pass is closed, we take the coast."
[user] Corvin (NPC): "The coast road is three days longer and the tide gates shut at dusk."
[assistant] Sable didn't answer at once. The map lay between them, weighted at the corners with stones.
[user] Sable: "Then we leave before dusk."
[assistant] Corvin folded the map without another word.
```

## Responses, all draws

### Qwen3.5-4B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 440 tok · 63.5 s · run `q5hb3tdm`

```text
C1 | HOLD | Claim that the pass is closed and the coast route is three days longer is marked as a specific condition for the current journey, but its validity depends on whether Sable accepts it as the definitive rule or if they find an alternative later. | Reason: The conversation establishes a constraint (closed pass) and a consequence (longer road), but no confirmation that these are permanent world rules versus temporary circumstances exists yet; however, given the immediate action of planning departure based on this info, it leans toward a durable fact for the session's context, yet the prompt instructions suggest HOLD for claims dependent on confirmation. Re-evaluating: The user states "The pass is closed" as a premise they are operating under. Corvin confirms the alternative. This is likely a world rule or established situation for the story. I will KEEP it as it sets the immediate plot constraints, but if strict adherence to "uncertain/contradictory" triggers HOLD, I should check if there's doubt. There is no explicit doubt expressed by Sable other than silence. Actually, looking at the rules: "KEEP: durable facts... world rules". The closed pass is a fact presented as such. "HOLD: claims marked uncertain... or dependent on user confirmation." Sable didn't confirm the pass was closed *in general*, only that it is currently closed for them. However, in RP contexts, stating "The X is Y" usually implies a world rule unless challenged. I will KEEP this as it dictates the immediate action and route choice which matters for future sessions involving travel.

Wait, looking closer at Rule 3: "HOLD: claims marked uncertain, contradictory, or dependent on user confirmation." Sable's statement "The pass is closed" is presented as a fact. Corvin adds information about the alternative. There is no contradiction or uncertainty expressed. Therefore, this is a durable fact for the session/world state.

C1 | KEEP | The pass is closed and the coast road is three days longer due to tide gates shut at dusk. | Reason: These are established facts that dictate the characters' immediate travel plan and future movement options in the story.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=0.7, top_p=0.8, top_k=20, min_p=0.0, presence_penalty=1.5, repeat_penalty=1.0 · finish `stop` · 87 tok · 13.1 s · run `6ha7sry6`

```text
C1 | KEEP | World Rule: The tide gates shut at dusk, forcing travelers to depart before then to avoid being trapped on the coast road. | This establishes a critical environmental constraint affecting the group's schedule and route choice.
C2 | KEEP | Fact: Sable has chosen the coastal route over the pass because it is currently closed. | This defines the specific path taken by the characters based on the world state.
```

### Qwen3.5-2B HauhauCS Uncensored — non-thinking

**draw 1** — seed 4101 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 39 tok · 2.7 s · run `54hcnvhy`

```text
ID | DECISION | MEMORY | REASON
none | none | NONE | No relevant information to store regarding relationships, world rules, preferences, or durable facts from this specific turn alone.
```

**draw 2** — seed 4102 · max_tokens 768 · sampling temperature=1.0, top_p=1.0, top_k=20, min_p=0.0, presence_penalty=2.0, repeat_penalty=1.0 · finish `stop` · 55 tok · 3.3 s · run `g4wiytnc`

```text
ID | DECISION | MEMORY | REASON
C1 | KEEP | The pass road is closed / the tide gates shut at dusk | These are critical geographical constraints that drive all subsequent travel decisions and must be known in future interactions with the environment or NPC Corvin.
```

### LFM2.5-1.2B Huihui abliterated — non-thinking

**draw 1** — seed 4101 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 79 tok · 2.9 s · run `ep2h4gc3`

```text
C1  
ID | DECISION | MEMORY | REASON  
C2 | ACTION | MEMORY | REASON  
C3 | RELATIONSHIP | MEMORY | REASON  
C4 | CLARIFICATION | MEMORY | REASON  

HOLD: uncertain about the exact timing of the tide gates.
```

**draw 2** — seed 4102 · max_tokens 512 · sampling temperature=0.1, top_p=1.0, top_k=50, min_p=0.0, presence_penalty=0.0, repeat_penalty=1.05 · finish `stop` · 68 tok · 2.3 s · run `vyz7jhl9`

```text
C1  
ID | DECISION | MEMORY | REASON  
C2 | ACTION | MEMORY | REASON  
C3 | PREFERENCE | MEMORY | REASON  
C4 | RELATIONSHIP | MEMORY | REASON  

NONE
```
