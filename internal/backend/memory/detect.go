/**
 * Event Detection Engine
 *
 * Heuristic analysis of chat messages to detect memory-worthy events.
 * Uses keyword dictionaries per event type with confidence scoring.
 *
 * @package memory
 * @version 1.0.0
 */

package memory

import (
	"strings"

	"howl-chat/internal/backend/types"
)

// EventKeywordMap maps event types to weighted trigger keywords
type EventKeywordMap map[types.EventType][]KeywordRule

// KeywordRule defines a keyword pattern and its confidence weight
type KeywordRule struct {
	Word   string  // Keyword or phrase (lowercase)
	Weight float64 // Confidence contribution (0-100)
}

// defaultKeywordMap is the built-in heuristic dictionary
var defaultKeywordMap = EventKeywordMap{
	types.EventDeath: {
		{Word: "died", Weight: 80}, {Word: "dead", Weight: 70},
		{Word: "killed", Weight: 85}, {Word: "death", Weight: 75},
		{Word: "murdered", Weight: 80}, {Word: "slain", Weight: 70},
		{Word: "perished", Weight: 65}, {Word: "fallen", Weight: 50},
		{Word: "no longer alive", Weight: 70}, {Word: "lifeless", Weight: 60},
		{Word: "passed away", Weight: 65}, {Word: "executed", Weight: 70},
	},
	types.EventBetrayal: {
		{Word: "betrayed", Weight: 80}, {Word: "traitor", Weight: 75},
		{Word: "backstabbed", Weight: 70}, {Word: "deceived", Weight: 60},
		{Word: "lied to", Weight: 50}, {Word: "broken trust", Weight: 65},
		{Word: "turned against", Weight: 70}, {Word: "double-crossed", Weight: 75},
		{Word: "conspiracy", Weight: 55}, {Word: "scheme", Weight: 40},
	},
	types.EventRomance: {
		{Word: "kissed", Weight: 60}, {Word: "love", Weight: 50},
		{Word: "lovers", Weight: 55}, {Word: "embraced", Weight: 50},
		{Word: "married", Weight: 70}, {Word: "wedding", Weight: 70},
		{Word: "confession", Weight: 45}, {Word: "heart", Weight: 30},
		{Word: "romance", Weight: 55}, {Word: "intimate", Weight: 50},
	},
	types.EventDiscovery: {
		{Word: "discovered", Weight: 60}, {Word: "found", Weight: 45},
		{Word: "uncovered", Weight: 55}, {Word: "revealed", Weight: 50},
		{Word: "secret", Weight: 55}, {Word: "hidden", Weight: 40},
		{Word: "artifact", Weight: 50}, {Word: "treasure", Weight: 45},
		{Word: "mystery solved", Weight: 65}, {Word: "truth", Weight: 35},
	},
	types.EventConflict: {
		{Word: "attacked", Weight: 55}, {Word: "fought", Weight: 50},
		{Word: "battle", Weight: 55}, {Word: "war", Weight: 60},
		{Word: "argued", Weight: 40}, {Word: "hostile", Weight: 45},
		{Word: "threatened", Weight: 50}, {Word: "ambushed", Weight: 60},
		{Word: "siege", Weight: 55}, {Word: "duel", Weight: 50},
	},
	types.EventCharacterChange: {
		{Word: "changed", Weight: 30}, {Word: "transformed", Weight: 45},
		{Word: "became", Weight: 30}, {Word: "left", Weight: 35},
		{Word: "joined", Weight: 40}, {Word: "returned", Weight: 35},
		{Word: "promoted", Weight: 45}, {Word: "banished", Weight: 55},
		{Word: "exiled", Weight: 50}, {Word: "crowned", Weight: 55},
	},
	types.EventLocation: {
		{Word: "arrived at", Weight: 45}, {Word: "entered", Weight: 40},
		{Word: "reached", Weight: 40}, {Word: "journeyed to", Weight: 45},
		{Word: "castle", Weight: 25}, {Word: "dungeon", Weight: 25},
		{Word: "village", Weight: 20}, {Word: "temple", Weight: 25},
		{Word: "ruins", Weight: 25}, {Word: "crossed into", Weight: 40},
	},
	types.EventItem: {
		{Word: "obtained", Weight: 45}, {Word: "acquired", Weight: 40},
		{Word: "stole", Weight: 45}, {Word: "received", Weight: 30},
		{Word: "sword", Weight: 20}, {Word: "amulet", Weight: 25},
		{Word: "ring", Weight: 20}, {Word: "potion", Weight: 20},
		{Word: "weapon", Weight: 20}, {Word: "artifact", Weight: 30},
	},
	// ---- MICRO-MEMORIES ----
	types.EventPreference: {
		// Favorite things
		{Word: "favorite", Weight: 65}, {Word: "favourite", Weight: 65},
		{Word: "preferred", Weight: 55}, {Word: "prefer", Weight: 50},
		{Word: "always choose", Weight: 50}, {Word: "go-to", Weight: 45},
		{Word: "my favorite food", Weight: 70}, {Word: "favorite color", Weight: 70},
		{Word: "favorite music", Weight: 70}, {Word: "favorite song", Weight: 70},
		{Word: "favorite artist", Weight: 70}, {Word: "favorite band", Weight: 70},
		{Word: "favorite movie", Weight: 70}, {Word: "favorite show", Weight: 70},
		{Word: "favorite book", Weight: 70}, {Word: "favorite game", Weight: 70},
		{Word: "favorite place", Weight: 70}, {Word: "favorite season", Weight: 65},
		{Word: "favorite weather", Weight: 60}, {Word: "favorite time", Weight: 55},
		{Word: "favorite drink", Weight: 70}, {Word: "favorite dessert", Weight: 70},
		{Word: "best food", Weight: 55}, {Word: "best flavor", Weight: 55},
		{Word: "love the taste", Weight: 50}, {Word: "taste of", Weight: 40},
		{Word: "reminds me of home", Weight: 55}, {Word: "comfort food", Weight: 55},
	},
	types.EventLikeDislike: {
		{Word: "i love", Weight: 55}, {Word: "i like", Weight: 45},
		{Word: "i enjoy", Weight: 50}, {Word: "i adore", Weight: 55},
		{Word: "i hate", Weight: 55}, {Word: "i dislike", Weight: 50},
		{Word: "can't stand", Weight: 55}, {Word: "despise", Weight: 55},
		{Word: "detest", Weight: 50}, {Word: "sick of", Weight: 40},
		{Word: "tired of", Weight: 35}, {Word: "never liked", Weight: 45},
		{Word: "always loved", Weight: 50}, {Word: "grew to love", Weight: 50},
		{Word: "can't get enough", Weight: 50}, {Word: "obsessed with", Weight: 50},
		{Word: "not a fan", Weight: 40}, {Word: "could do without", Weight: 35},
		{Word: "love this song", Weight: 55}, {Word: "hate this song", Weight: 55},
		{Word: "love this picture", Weight: 55}, {Word: "love this painting", Weight: 55},
		{Word: "beautiful", Weight: 35}, {Word: "ugly", Weight: 30},
		{Word: "amazing", Weight: 30}, {Word: "terrible", Weight: 30},
	},
	types.EventEmotion: {
		{Word: "makes me happy", Weight: 55}, {Word: "makes me sad", Weight: 55},
		{Word: "makes me laugh", Weight: 55}, {Word: "makes me cry", Weight: 55},
		{Word: "makes me angry", Weight: 55}, {Word: "makes me nervous", Weight: 50},
		{Word: "i feel safe", Weight: 60}, {Word: "i feel scared", Weight: 55},
		{Word: "i feel happy", Weight: 50}, {Word: "i feel sad", Weight: 50},
		{Word: "i feel loved", Weight: 60}, {Word: "i feel alone", Weight: 50},
		{Word: "comforting", Weight: 45}, {Word: "terrifying", Weight: 50},
		{Word: "calming", Weight: 45}, {Word: "anxious", Weight: 40},
		{Word: "peaceful", Weight: 40}, {Word: "overwhelmed", Weight: 40},
		{Word: "grateful", Weight: 50}, {Word: "thankful", Weight: 45},
		{Word: "proud of", Weight: 45}, {Word: "ashamed", Weight: 40},
		{Word: "jealous", Weight: 45}, {Word: "envy", Weight: 40},
		{Word: "nostalgic", Weight: 45}, {Word: "homesick", Weight: 40},
		{Word: "lonely", Weight: 40}, {Word: "content", Weight: 35},
		{Word: "excited", Weight: 40}, {Word: "disappointed", Weight: 40},
	},
	types.EventAssociation: {
		{Word: "reminds me of", Weight: 60}, {Word: "reminds me", Weight: 50},
		{Word: "think of my", Weight: 45}, {Word: "makes me think of", Weight: 55},
		{Word: "brings back memories", Weight: 60}, {Word: "brings back", Weight: 50},
		{Word: "flashback", Weight: 50}, {Word: "like when i was", Weight: 45},
		{Word: "my mother used to", Weight: 50}, {Word: "my father used to", Weight: 50},
		{Word: "my sister used to", Weight: 50}, {Word: "my brother used to", Weight: 50},
		{Word: "my grandmother", Weight: 45}, {Word: "my grandfather", Weight: 45},
		{Word: "childhood", Weight: 45}, {Word: "growing up", Weight: 40},
		{Word: "when i was young", Weight: 45}, {Word: "when i was little", Weight: 45},
		{Word: "familiar", Weight: 35}, {Word: "deja vu", Weight: 40},
		{Word: "like home", Weight: 45}, {Word: "like family", Weight: 45},
		{Word: "my old", Weight: 35}, {Word: "used to have", Weight: 35},
	},
	types.EventAesthetic: {
		{Word: "beautiful", Weight: 45}, {Word: "stunning", Weight: 50},
		{Word: "gorgeous", Weight: 50}, {Word: "breathtaking", Weight: 55},
		{Word: "magnificent", Weight: 50}, {Word: "lovely", Weight: 40},
		{Word: "elegant", Weight: 45}, {Word: "charming", Weight: 40},
		{Word: "love the colors", Weight: 55}, {Word: "love the style", Weight: 50},
		{Word: "perfect design", Weight: 50}, {Word: "incredible art", Weight: 55},
		{Word: "masterpiece", Weight: 55}, {Word: "work of art", Weight: 50},
		{Word: "aesthetic", Weight: 40}, {Word: "visually", Weight: 35},
		{Word: "architecture", Weight: 35}, {Word: "landscape", Weight: 35},
		{Word: "scenery", Weight: 35}, {Word: "sunset", Weight: 30},
		{Word: "sunrise", Weight: 30}, {Word: "stars", Weight: 25},
		{Word: "moonlight", Weight: 30}, {Word: "candlelight", Weight: 30},
	},
	types.EventSensory: {
		{Word: "smells like", Weight: 50}, {Word: "smells of", Weight: 45},
		{Word: "tastes like", Weight: 50}, {Word: "tastes of", Weight: 45},
		{Word: "sounds like", Weight: 45}, {Word: "sounds of", Weight: 40},
		{Word: "feels like", Weight: 45}, {Word: "texture", Weight: 35},
		{Word: "that smell", Weight: 45}, {Word: "that sound", Weight: 40},
		{Word: "that taste", Weight: 45}, {Word: "that feeling", Weight: 40},
		{Word: "warmth", Weight: 35}, {Word: "chill", Weight: 30},
		{Word: "breeze", Weight: 30}, {Word: "fragrance", Weight: 35},
		{Word: "aroma", Weight: 35}, {Word: "scent", Weight: 30},
		{Word: "bitter", Weight: 25}, {Word: "sweet", Weight: 25},
		{Word: "sour", Weight: 25}, {Word: "spicy", Weight: 25},
		{Word: "salty", Weight: 25}, {Word: "savory", Weight: 25},
		{Word: "crunchy", Weight: 25}, {Word: "smooth", Weight: 25},
		{Word: "rough", Weight: 25}, {Word: "soft", Weight: 25},
	},
	types.EventSocial: {
		{Word: "inside joke", Weight: 60}, {Word: "our joke", Weight: 55},
		{Word: "remember when we", Weight: 60}, {Word: "that time we", Weight: 55},
		{Word: "we laughed", Weight: 45}, {Word: "we cried", Weight: 45},
		{Word: "thank you", Weight: 40}, {Word: "grateful for you", Weight: 50},
		{Word: "appreciate you", Weight: 50}, {Word: "you're the best", Weight: 45},
		{Word: "you mean a lot", Weight: 55}, {Word: "you matter to me", Weight: 60},
		{Word: "compliment", Weight: 40}, {Word: "you look", Weight: 35},
		{Word: "you're beautiful", Weight: 50}, {Word: "you're amazing", Weight: 50},
		{Word: "nickname", Weight: 40}, {Word: "call me", Weight: 35},
		{Word: "teasing", Weight: 35}, {Word: "joking with", Weight: 35},
		{Word: "apologize", Weight: 35}, {Word: "i'm sorry", Weight: 35},
		{Word: "forgive", Weight: 40}, {Word: "trust you", Weight: 50},
		{Word: "believe in you", Weight: 50}, {Word: "count on you", Weight: 50},
	},
	types.EventFirst: {
		{Word: "first time", Weight: 65}, {Word: "first time i", Weight: 65},
		{Word: "first time i've", Weight: 65}, {Word: "never done this", Weight: 55},
		{Word: "never seen this", Weight: 55}, {Word: "never been", Weight: 50},
		{Word: "never tried", Weight: 50}, {Word: "my first", Weight: 55},
		{Word: "the first", Weight: 45}, {Word: "first experience", Weight: 55},
		{Word: "new to me", Weight: 45}, {Word: "something new", Weight: 40},
		{Word: "discovered this", Weight: 45}, {Word: "just learned", Weight: 40},
		{Word: "beginning to", Weight: 35}, {Word: "starting to", Weight: 35},
		{Word: "virgin", Weight: 40}, {Word: "inexperienced", Weight: 35},
	},
	types.EventShared: {
		{Word: "you showed me", Weight: 70}, {Word: "you shared", Weight: 65},
		{Word: "you sent me", Weight: 60}, {Word: "you gave me", Weight: 55},
		{Word: "thanks for showing", Weight: 65}, {Word: "thanks for sharing", Weight: 65},
		{Word: "remember when you", Weight: 60}, {Word: "that picture you", Weight: 60},
		{Word: "that photo you", Weight: 60}, {Word: "that image you", Weight: 55},
		{Word: "that song you", Weight: 60}, {Word: "that video you", Weight: 55},
		{Word: "we did together", Weight: 55}, {Word: "we experienced", Weight: 50},
		{Word: "our moment", Weight: 50}, {Word: "our memory", Weight: 50},
		{Word: "together we", Weight: 45}, {Word: "with you", Weight: 40},
		{Word: "shared with", Weight: 45}, {Word: "showed me this", Weight: 65},
		{Word: "shared this with", Weight: 60}, {Word: "sent me this", Weight: 55},
	},
	types.EventPromise: {
		{Word: "i'll remember", Weight: 70}, {Word: "i will remember", Weight: 70},
		{Word: "won't forget", Weight: 65}, {Word: "never forget", Weight: 60},
		{Word: "always remember", Weight: 65}, {Word: "can't forget", Weight: 55},
		{Word: "means a lot", Weight: 60}, {Word: "means the world", Weight: 65},
		{Word: "so important to me", Weight: 55}, {Word: "cherish this", Weight: 60},
		{Word: "treasure this", Weight: 60}, {Word: "hold onto this", Weight: 55},
		{Word: "i promise", Weight: 50}, {Word: "i swear", Weight: 45},
		{Word: "my word", Weight: 45}, {Word: "vow to", Weight: 45},
		{Word: "committed to", Weight: 40}, {Word: "dedicated to", Weight: 40},
		{Word: "this moment", Weight: 35}, {Word: "forever", Weight: 35},
		{Word: "eternity", Weight: 35}, {Word: "always", Weight: 30},
	},
}

// DetectEvents analyzes a message for memory-worthy events
func DetectEvents(message string, history []string, config types.AutoCaptureConfig) []types.EventDetectionResult {
	if !config.Enabled {
		return nil
	}

	lowerMsg := strings.ToLower(message)
	var results []types.EventDetectionResult

	for eventType, rules := range defaultKeywordMap {
		confidence, evidence := matchKeywords(lowerMsg, rules)
		if confidence >= config.MinConfidence {
			// Boost for multi-word phrases (more specific = higher confidence)
			confidence = boostForPhrases(evidence, confidence)

			// Add sentiment intensity
			sentimentBoost := analyzeSentimentIntensity(lowerMsg)
			confidence += sentimentBoost
			if confidence > 100 {
				confidence = 100
			}

			context := extractContext(message, history)
			results = append(results, types.EventDetectionResult{
				EventType:  eventType,
				Confidence: confidence,
				Evidence:   evidence,
				Context:    context,
			})
		}
	}

	// Also check custom trigger words from config
	if len(config.TriggerWords) > 0 {
		customConfidence, customEvidence := matchCustomTriggers(lowerMsg, config.TriggerWords)
		if customConfidence >= config.MinConfidence {
			results = append(results, types.EventDetectionResult{
				EventType:  types.EventCustom,
				Confidence: customConfidence,
				Evidence:   customEvidence,
				Context:    message,
			})
		}
	}

	// Apply cross-category suppression: death suppresses conflict, etc.
	results = suppressSubordinateEvents(results)

	return limitResults(results, config.MaxPerMessage)
}

// matchKeywords scans text against a rule set and returns aggregate confidence
func matchKeywords(text string, rules []KeywordRule) (float64, []string) {
	var evidence []string
	var totalWeight float64
	seen := make(map[string]bool)

	for _, rule := range rules {
		if strings.Contains(text, rule.Word) && !seen[rule.Word] {
			evidence = append(evidence, rule.Word)
			totalWeight += rule.Weight
			seen[rule.Word] = true
		}
	}

	// Confidence is weighted sum capped at 100
	confidence := totalWeight
	if confidence > 100 {
		confidence = 100
	}

	// Boost confidence slightly for multiple independent matches
	if len(evidence) > 1 {
		confidence += float64(len(evidence)) * 3
		if confidence > 100 {
			confidence = 100
		}
	}

	return confidence, evidence
}

// matchCustomTriggers checks user-defined trigger words
func matchCustomTriggers(text string, triggers []string) (float64, []string) {
	var evidence []string
	for _, t := range triggers {
		if strings.Contains(text, strings.ToLower(t)) {
			evidence = append(evidence, t)
		}
	}
	if len(evidence) == 0 {
		return 0, nil
	}
	// Base confidence + bonus for multiple triggers
	confidence := 50.0 + float64(len(evidence))*10
	if confidence > 100 {
		confidence = 100
	}
	return confidence, evidence
}

// extractContext builds a short context string from message and recent history
func extractContext(message string, history []string) string {
	// Include up to 2 previous messages for context
	start := 0
	if len(history) > 2 {
		start = len(history) - 2
	}

	var parts []string
	for i := start; i < len(history); i++ {
		parts = append(parts, truncate(history[i], 120))
	}
	parts = append(parts, truncate(message, 200))

	return strings.Join(parts, " | ")
}

// limitResults caps results at maxPerMessage, keeping highest confidence
func limitResults(results []types.EventDetectionResult, max int) []types.EventDetectionResult {
	if max <= 0 {
		max = 3
	}
	if len(results) <= max {
		return results
	}

	// Simple bubble sort by confidence descending
	for i := 0; i < len(results)-1; i++ {
		for j := 0; j < len(results)-i-1; j++ {
			if results[j].Confidence < results[j+1].Confidence {
				results[j], results[j+1] = results[j+1], results[j]
			}
		}
	}

	return results[:max]
}

// truncate shortens text to maxLen with ellipsis
func truncate(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}

// boostForPhrases gives bonus confidence for multi-word phrase matches
// Longer, more specific phrases indicate stronger intent
func boostForPhrases(evidence []string, baseConfidence float64) float64 {
	boost := 0.0
	for _, e := range evidence {
		words := strings.Fields(e)
		if len(words) >= 3 {
			boost += 5.0 // Strong phrase match
		} else if len(words) == 2 {
			boost += 2.0 // Moderate phrase match
		}
	}
	if baseConfidence+boost > 100 {
		return 100
	}
	return baseConfidence + boost
}

// analyzeSentimentIntensity scans for emotional intensity markers
// Returns a confidence boost (0-15) based on intensity of language
func analyzeSentimentIntensity(text string) float64 {
	intensity := 0.0

	// High intensity markers
	highIntensity := []string{
		"absolutely", "completely", "utterly", "totally", "incredibly",
		"so much", "so happy", "so sad", "so angry", "so scared",
		"really love", "really hate", "truly", "deeply", "profoundly",
	}
	for _, marker := range highIntensity {
		if strings.Contains(text, marker) {
			intensity += 3.0
		}
	}

	// Exclamation intensity
	if strings.Count(text, "!") >= 2 {
		intensity += 2.0
	}

	// Repeated letters (e.g., "sooooo")
	if strings.Contains(text, "ooo") || strings.Contains(text, "aaa") || strings.Contains(text, "hhh") {
		intensity += 2.0
	}

	if intensity > 15 {
		return 15
	}
	return intensity
}

// suppressSubordinateEvents removes lower-confidence overlapping events
// e.g., if death (85%) and conflict (60%) both fire, suppress conflict
func suppressSubordinateEvents(results []types.EventDetectionResult) []types.EventDetectionResult {
	suppressionRules := map[types.EventType][]types.EventType{
		types.EventDeath:    {types.EventConflict, types.EventCharacterChange},
		types.EventBetrayal: {types.EventConflict},
		types.EventRomance:  {types.EventLikeDislike, types.EventEmotion},
		types.EventPromise:  {types.EventEmotion, types.EventSocial},
		types.EventShared:   {types.EventLikeDislike, types.EventPreference},
		types.EventFirst:    {types.EventDiscovery, types.EventLocation},
	}

	// Build set of dominant event types that fired
	dominant := make(map[types.EventType]bool)
	for _, r := range results {
		dominant[r.EventType] = true
	}

	var filtered []types.EventDetectionResult
	for _, r := range results {
		suppressed := false
		for domType, subTypes := range suppressionRules {
			if !dominant[domType] || domType == r.EventType {
				continue
			}
			for _, sub := range subTypes {
				if r.EventType == sub {
					suppressed = true
					break
				}
			}
			if suppressed {
				break
			}
		}
		if !suppressed {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
