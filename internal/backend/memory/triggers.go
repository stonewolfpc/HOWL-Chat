/**
 * Trigger Generation
 *
 * Auto-generates trigger phrases from memory content using NLP heuristics.
 * Extracts key entities, proper nouns, and significant phrases.
 *
 * @package memory
 * @version 1.0.0
 */

package memory

import (
	"howl-chat/internal/backend/types"
	"regexp"
	"strings"
	"unicode"
)

// generateTriggerPhrases creates primary triggers from memory title and summary
func generateTriggerPhrases(record types.MemoryRecord) []string {
	phrases := []string{}

	// Add title words (excluding stop words)
	titleWords := extractKeyTerms(record.Title)
	for _, term := range titleWords {
		if len(term) > 3 {
			phrases = append(phrases, strings.ToLower(term))
		}
	}

	// Add named entities from summary
	entities := extractNamedEntities(record.Summary)
	phrases = append(phrases, entities...)

	// Add significant bigrams from summary
	bigrams := extractSignificantBigrams(record.Summary)
	phrases = append(phrases, bigrams...)

	// Deduplicate and limit
	phrases = uniqueStrings(phrases)
	if len(phrases) > 8 {
		phrases = phrases[:8]
	}

	return phrases
}

// generateSecondaryTriggers creates backup/alternative triggers
func generateSecondaryTriggers(record types.MemoryRecord) []string {
	phrases := []string{}

	// Add tags as triggers
	for _, tag := range record.Tags {
		if len(tag) > 2 {
			phrases = append(phrases, strings.ToLower(tag))
		}
	}

	// Add scope-specific context
	if record.Scope == "character" && record.OwnerID != "" {
		// Add owner name variations
		ownerParts := strings.Fields(record.OwnerID)
		for _, part := range ownerParts {
			if len(part) > 2 {
				phrases = append(phrases, strings.ToLower(part))
			}
		}
	}

	return uniqueStrings(phrases)
}

// extractKeyTerms extracts important words from text
func extractKeyTerms(text string) []string {
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true,
		"but": true, "in": true, "on": true, "at": true, "to": true,
		"for": true, "of": true, "is": true, "are": true, "was": true,
		"it": true, "they": true, "them": true, "their": true,
		"has": true, "have": true, "had": true, "been": true,
	}

	words := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	var terms []string
	for _, word := range words {
		word = strings.ToLower(strings.TrimSpace(word))
		if len(word) > 3 && !stopWords[word] {
			terms = append(terms, word)
		}
	}

	return terms
}

// extractNamedEntities extracts capitalized words (simple heuristic for names/places)
func extractNamedEntities(text string) []string {
	re := regexp.MustCompile(`\b[A-Z][a-zA-Z]{2,}\b`)
	matches := re.FindAllString(text, -1)

	// Filter common false positives
	commonWords := map[string]bool{
		"The": true, "A": true, "An": true, "And": true, "But": true,
		"He": true, "She": true, "They": true, "It": true, "This": true,
		"That": true, "These": true, "Those": true, "His": true, "Her": true,
	}

	var entities []string
	for _, match := range matches {
		if !commonWords[match] {
			entities = append(entities, strings.ToLower(match))
		}
	}

	return uniqueStrings(entities)
}

// extractSignificantBigrams extracts meaningful two-word phrases
func extractSignificantBigrams(text string) []string {
	words := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true,
		"but": true, "in": true, "on": true, "at": true, "to": true,
		"of": true, "for": true, "with": true, "by": true,
	}

	var bigrams []string
	for i := 0; i < len(words)-1; i++ {
		w1, w2 := words[i], words[i+1]
		if !stopWords[w1] && !stopWords[w2] && len(w1) > 2 && len(w2) > 2 {
			bigrams = append(bigrams, w1+" "+w2)
		}
	}

	return uniqueStrings(bigrams)
}

// uniqueStrings removes duplicates from string slice
func uniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range slice {
		s = strings.TrimSpace(strings.ToLower(s))
		if s != "" && !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}
