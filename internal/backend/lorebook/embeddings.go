package lorebook

import (
    "math"
    "strings"
    "unicode"
)

// SimpleEmbedding provides a lightweight vector embedding without external dependencies
// Uses TF-IDF-inspired term weighting for semantic similarity
type SimpleEmbedding struct {
    terms   map[string]float64
    norm    float64
    rawText string
}

// NewEmbedding creates an embedding from text using simplified TF-IDF
func NewEmbedding(text string) *SimpleEmbedding {
    text = strings.ToLower(strings.TrimSpace(text))
    terms := tokenize(text)
    weighted := weightTerms(terms)

    norm := 0.0
    for _, w := range weighted {
        norm += w * w
    }
    norm = math.Sqrt(norm)

    if norm == 0 {
        norm = 1.0
    }

    return &SimpleEmbedding{
        terms:   weighted,
        norm:    norm,
        rawText: text,
    }
}

// CosineSimilarity calculates cosine similarity between this embedding and another (0-100 scale)
func (e *SimpleEmbedding) CosineSimilarity(other *SimpleEmbedding) int {
    if e.norm == 0 || other.norm == 0 {
        return 0
    }

    dotProduct := 0.0
    for term, weight := range e.terms {
        if otherWeight, ok := other.terms[term]; ok {
            dotProduct += weight * otherWeight
        }
    }

    similarity := dotProduct / (e.norm * other.norm)
    if similarity < 0 {
        similarity = 0
    }
    if similarity > 1 {
        similarity = 1
    }

    return int(similarity * 100)
}

// IsSemanticMatch checks if similarity exceeds threshold
func (e *SimpleEmbedding) IsSemanticMatch(other *SimpleEmbedding, threshold int) bool {
    return e.CosineSimilarity(other) >= threshold
}

// tokenize breaks text into meaningful tokens, removing stopwords and short tokens
func tokenize(text string) []string {
    words := strings.FieldsFunc(text, func(r rune) bool {
        return !unicode.IsLetter(r) && !unicode.IsNumber(r)
    })

    var tokens []string
    stopwords := map[string]bool{
        "the": true, "a": true, "an": true, "and": true, "or": true, "but": true,
        "in": true, "on": true, "at": true, "to": true, "for": true, "of": true,
        "is": true, "are": true, "am": true, "was": true, "were": true, "be": true,
        "it": true, "they": true, "them": true, "their": true, "as": true,
    }

    for _, word := range words {
        if len(word) > 2 && !stopwords[word] {
            tokens = append(tokens, word)
        }
    }

    return tokens
}

// weightTerms assigns weights based on term frequency and length bonus
func weightTerms(terms []string) map[string]float64 {
    frequencies := make(map[string]int)
    for _, term := range terms {
        frequencies[term]++
    }

    weights := make(map[string]float64)
    for term, freq := range frequencies {
        base := math.Log(float64(freq) + 1)
        lengthBonus := 1.0
        if len(term) > 8 {
            lengthBonus = 1.5
        }
        weights[term] = base * lengthBonus
    }

    return weights
}
