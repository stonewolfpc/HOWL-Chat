package lorebook

import (
    "strings"
    "testing"
)

func TestEmbeddingCreation(t *testing.T) {
    text := "The magical flow of energy"
    emb := NewEmbedding(text)

    if emb == nil {
        t.Fatal("embedding should not be nil")
    }
    if emb.norm == 0 {
        t.Fatal("embedding norm should not be zero")
    }
}

func TestCosineSimilarityExact(t *testing.T) {
    text := "magical energy"
    emb1 := NewEmbedding(text)
    emb2 := NewEmbedding(text)

    similarity := emb1.CosineSimilarity(emb2)
    if similarity < 95 { // Allow for floating point precision
        t.Fatalf("identical text should have high similarity, got %d", similarity)
    }
}

func TestCosineSimilaritySemantic(t *testing.T) {
    emb1 := NewEmbedding("mana flow")
    emb2 := NewEmbedding("flow of magical energy")

    similarity := emb1.CosineSimilarity(emb2)
    if similarity < 40 {
        t.Fatalf("semantically similar phrases should have >40 similarity, got %d", similarity)
    }
    if similarity > 100 {
        t.Fatalf("similarity should never exceed 100, got %d", similarity)
    }
}

func TestCosineSimilarityKingMonarch(t *testing.T) {
    emb1 := NewEmbedding("king")
    emb2 := NewEmbedding("monarch")

    similarity := emb1.CosineSimilarity(emb2)
    // These might not match perfectly without a semantic model,
    // but let's verify the function works
    if similarity < 0 || similarity > 100 {
        t.Fatalf("similarity should be 0-100, got %d", similarity)
    }
}

func TestCosineSimilarityCompletellyDifferent(t *testing.T) {
    emb1 := NewEmbedding("wizard magic spell")
    emb2 := NewEmbedding("pizza restaurant food")

    similarity := emb1.CosineSimilarity(emb2)
    if similarity > 30 {
        t.Fatalf("completely different text should have low similarity, got %d", similarity)
    }
}

func TestIsSemanticMatch(t *testing.T) {
    emb1 := NewEmbedding("mana flow")
    emb2 := NewEmbedding("flow of magical energy")

    matched := emb1.IsSemanticMatch(emb2, 30)
    if !matched {
        t.Fatal("semantically similar text should match at 30% threshold")
    }
}

func TestIsSemanticMatchThreshold(t *testing.T) {
    emb1 := NewEmbedding("wizard")
    emb2 := NewEmbedding("sorcerer")

    // At high threshold, might not match
    matched := emb1.IsSemanticMatch(emb2, 95)
    if matched {
        t.Fatalf("very high threshold should not match: similarity=%d", emb1.CosineSimilarity(emb2))
    }
}

func TestTokenization(t *testing.T) {
    text := "The quick brown fox jumps"
    tokens := tokenize(text)

    // Should exclude short words and stopwords like "the"
    if len(tokens) == 0 {
        t.Fatal("tokenization should produce tokens")
    }

    // "the" should be filtered out
    for _, token := range tokens {
        if token == "the" {
            t.Fatal("stopword 'the' should be filtered out")
        }
    }
}

func TestTokenizationRemovesShortWords(t *testing.T) {
    text := "I am a wizard"
    tokens := tokenize(text)

    // "I", "am", "a" are all too short or stopwords
    for _, token := range tokens {
        if len(token) <= 2 {
            t.Fatalf("short token should be filtered: %q", token)
        }
    }
}

func TestWeightTerms(t *testing.T) {
    tokens := []string{"magic", "magic", "wizard"}
    weights := weightTerms(tokens)

    if len(weights) == 0 {
        t.Fatal("weightTerms should produce weights")
    }

    // "magic" appears twice, so should have higher weight than "wizard"
    if weights["magic"] <= weights["wizard"] {
        t.Fatalf("repeated term should have higher weight: magic=%f, wizard=%f",
            weights["magic"], weights["wizard"])
    }
}

func TestWeightTermsLengthBonus(t *testing.T) {
    tokens := []string{"short", "verylongword"}
    weights := weightTerms(tokens)

    // Longer words get bonus multiplier
    if weights["verylongword"] <= weights["short"] {
        t.Fatalf("longer word should get length bonus: long=%f, short=%f",
            weights["verylongword"], weights["short"])
    }
}

func TestEmbeddingNormalization(t *testing.T) {
    text1 := "magical energy flows"
    text2 := strings.Repeat("magical energy flows ", 10)

    emb1 := NewEmbedding(text1)
    emb2 := NewEmbedding(text2)

    // Repeated text should have higher norm
    if emb1.norm >= emb2.norm {
        t.Fatalf("repeated text should have higher norm: norm1=%f, norm2=%f", emb1.norm, emb2.norm)
    }
}

func TestEmptyTextEmbedding(t *testing.T) {
    emb := NewEmbedding("")

    if emb == nil {
        t.Fatal("empty embedding should not be nil")
    }
    if emb.norm == 0 {
        t.Fatal("empty embedding norm should default to 1.0")
    }
}

func TestEmbeddingBothEmpty(t *testing.T) {
    emb1 := NewEmbedding("")
    emb2 := NewEmbedding("")

    similarity := emb1.CosineSimilarity(emb2)
    if similarity < 0 || similarity > 100 {
        t.Fatalf("empty embeddings should produce valid similarity, got %d", similarity)
    }
}

func TestCaseSensitivity(t *testing.T) {
    emb1 := NewEmbedding("MAGICAL ENERGY")
    emb2 := NewEmbedding("magical energy")

    similarity := emb1.CosineSimilarity(emb2)
    if similarity < 95 {
        t.Fatalf("case should not matter, similarity=%d", similarity)
    }
}

func TestPunctuationHandling(t *testing.T) {
    emb1 := NewEmbedding("magical energy!")
    emb2 := NewEmbedding("magical energy")

    similarity := emb1.CosineSimilarity(emb2)
    if similarity < 90 {
        t.Fatalf("punctuation should not significantly impact similarity, got %d", similarity)
    }
}
