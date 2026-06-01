package emote

import (
	"github.com/sahilm/fuzzy"
)

// IsASCII reports whether every character in s is ASCII.
func IsASCII(s string) bool {
	for _, r := range s {
		if r > 127 {
			return false
		}
	}
	return true
}

// Emote is a standardized struct for an emote
type Emote struct {
	ID     string
	Name   string
	Type   string
	Width  int
	Height int
	URL    string
}

// emotes is a slice of emotes
type emotes []Emote

// Len returns the length of emotes
func (e emotes) Len() int {
	return len(e)
}

// String returns the name of the emote at the given index
func (e emotes) String(i int) string {
	return e[i].Name
}

// SearchEmotes fuzzy searches emotes from the query text.
// For non-ASCII queries (e.g. Cyrillic), the input slice is returned as-is
// since fuzzy matching against mostly-ASCII emote names would yield no results.
func SearchEmotes(query string, e []Emote) []Emote {
	if !IsASCII(query) {
		return e
	}
	matches := fuzzy.FindFrom(query, emotes(e))

	rankedEmotes := make([]Emote, len(matches))
	for i, match := range matches {
		rankedEmotes[i] = e[match.Index]
	}

	return rankedEmotes
}

// ExactSearchEmotes returns only emotes whose name exactly matches the query (case-sensitive).
func ExactSearchEmotes(query string, e []Emote) []Emote {
	var result []Emote
	for _, em := range e {
		if em.Name == query {
			result = append(result, em)
		}
	}
	return result
}
