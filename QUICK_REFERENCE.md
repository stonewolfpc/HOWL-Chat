# Howl Chat Lorebook - Quick Reference Guide

## Using Semantic Matching

### In the UI Form:
1. Create/edit a lore entry
2. Under **Trigger Mode**, select **"semantic"**
3. Enter trigger phrases that express concepts, not exact matches

### Example:
```
Entry: "Ancient Magic"
Trigger Phrase: "mystical force ancient"
Trigger Mode: semantic

When user writes: "magical power from long ago"
Result: ✓ MATCHES (semantic similarity ~60%)

When user writes: "what's for dinner?"
Result: ✗ NO MATCH (similarity ~5%)
```

### How It Works:
- Breaks text into meaningful tokens
- Removes common words (the, and, is, etc.)
- Weights longer/less common words higher
- Calculates cosine similarity (0-100)
- Threshold: 40% similarity = match

---

## Using Conflict Resolution

### Setup:
When you have **two entries that match the same trigger**, set `ConflictRule` on one or both:

```json
Entry 1: {
  "id": "world_lore",
  "triggerPhrases": ["dragon"],
  "conflictRule": "highest_priority",
  "priorityLevel": 1
}

Entry 2: {
  "id": "character_memory",
  "triggerPhrases": ["dragon"],
  "conflictRule": "use_both",
  "priorityLevel": 3
}
```

### Four Strategies:

#### 1. **highest_priority** (Default)
- Keeps the entry with **lowest priority number** (1 = highest)
- Use for: Truth hierarchy (world > scenario > character)
- Result: Only 1 entry returned

#### 2. **newer_wins**
- Keeps the entry with **most recent UpdatedAt timestamp**
- Use for: Character memories that update
- Result: Only 1 entry (the newest)

#### 3. **use_both**
- Keeps **all matching entries**
- Use for: Complementary perspectives
- Result: Multiple entries returned in order

#### 4. **merge_summaries**
- Combines all entries: `"Content 1" | "Content 2" | "Content 3"`
- Use for: Consolidated information
- Result: 1 merged entry

### Real-World Examples:

**Scenario 1: World vs Character Memory**
```
World Lore: "The war ended 50 years ago"
Character Memory: "The war ended when I was a child"
Conflict Rule: newer_wins
→ Uses whichever was updated more recently
```

**Scenario 2: Multiple Perspectives**
```
Observer 1: "The dragon was blue"
Observer 2: "The dragon was gold"
Conflict Rule: use_both
→ Includes both perspectives: "blue dragon" AND "gold dragon"
```

**Scenario 3: Consolidated Info**
```
Fact 1: "Dragons have wings"
Fact 2: "Dragons breathe fire"
Fact 3: "Dragons are rare"
Conflict Rule: merge_summaries
→ "Dragons have wings | Dragons breathe fire | Dragons are rare"
```

---

## Scope-Based Conflict Defaults

If you **don't specify ConflictRule**, it defaults by scope:

| Scope | Default Rule | Reason |
|-------|--------------|--------|
| World | highest_priority | One truth for the world |
| Scenario | use_both | Multiple things happen in a scene |
| Character | newer_wins | Memories change/update over time |

---

## Budget Edge Cases

### Scenario: Budget Too Small

```go
MaxCharacters: 50
Entry Content: 200 characters

Result: Entry is SKIPPED (not included)
```

**The system will:**
- ✓ Never error out
- ✓ Not break the prompt
- ✓ Return nothing rather than truncate unfairly
- ✓ Try other entries in priority order

### Budget Handling Examples:

```
MaxEntries: 4
MaxCharacters: 1000

Entry 1: 300 chars (included, used: 300)
Entry 2: 400 chars (included, used: 700)
Entry 3: 400 chars (SKIPPED - would exceed 1000)
Entry 4: 200 chars (included, used: 900)

Result: Entries 1, 2, 4 included
```

---

## Trigger Mode Comparison

| Feature | Exact | Loose | Semantic |
|---------|-------|-------|----------|
| Matches "the manager"? | No | Yes | No |
| Matches "mana" for trigger "mana"? | Yes | Yes | Yes |
| Matches "mana flow" for "flow of magical energy"? | No | No | Yes |
| Speed | Fastest | Fast | Fast (TF-IDF) |
| False Positives | None | Some | Rare |
| Best For | Specific keywords | Flexible phrases | Concept matching |

---

## Testing Entries

### Quick Test: Semantic vs Exact

Create two entries with same content but different trigger modes:

```
Entry A:
- Trigger: "wizardry"
- Mode: exact

Entry B:
- Trigger: "magical arts"
- Mode: semantic

Test 1: User says "wizardry"
→ A matches (exact), B matches (has "magical")

Test 2: User says "arcane arts"
→ A no match, B matches (semantic: "magical"+"arts" → "arcane"+"arts")
```

### Debugging Conflicts

If you're not seeing expected results:

1. **Check TriggerMode** - is it set to the right mode?
2. **Check ConflictRule** - are entries conflicting and being merged?
3. **Check Budget** - is content being skipped for size?
4. **Check Frequency** - did this entry already trigger in the chat?
5. **Check Enabled** - is the entry actually enabled?

---

## Performance Tips

1. **Semantic mode on important concepts**, exact mode for keywords
2. **Conflict rules only needed** when entries actually overlap
3. **Budget defaults work well** - 2400 chars total, 160 per entry
4. **Priority numbers**: 1-5 typical, lower = higher priority
5. **Keep trigger phrases to 2-5 words** for best semantic matching

---

## API Integration

### For Developers:

```go
import "howl-chat/internal/backend/lorebook"

// Create semantic embeddings
emb1 := lorebook.NewEmbedding("dragon attack")
emb2 := lorebook.NewEmbedding("dragon assault")
similarity := emb1.CosineSimilarity(emb2)  // Returns 0-100

// Resolve with conflict handling
entries := lorebook.Resolve(lorebook.ResolveRequest{
    Message: userMessage,
    MaxEntries: 8,
    MaxCharacters: 2400,
})

// Find conflicts
conflicts := lorebook.FindConflicts(entries)
```

---

## Checklist for QA Testing

- [ ] Create entry with `TriggerMode: "semantic"`
- [ ] Verify semantic match works (similar concepts match)
- [ ] Create two entries with same trigger phrase
- [ ] Set one with `ConflictRule: "newest_wins"`
- [ ] Verify only one returned (the newer one)
- [ ] Try `use_both` - verify both returned
- [ ] Try `merge_summaries` - verify content merged with ` | `
- [ ] Test with very small MaxCharacters (50 bytes)
- [ ] Verify no errors, graceful skip
- [ ] Test with disabled entry - verify it's skipped
- [ ] Test with empty content - verify it's skipped

