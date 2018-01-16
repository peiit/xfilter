package keywords

import (
	"strings"
	"unicode"
)

// Keywords maps keywords to integers representing user IDs
type Keywords struct {
	kw map[string][]int64
}

// New creates a new set of keywords
func New() *Keywords {
	return &Keywords{kw: make(map[string][]int64)}
}

// Add adds a keyword for the given user
func (k *Keywords) Add(word string, user int64) {
	word = strings.ToLower(word)
	users := k.kw[word]
	for _, u := range users {
		if u == user {
			return
		}
	}
	k.kw[word] = append(k.kw[word], user)
}

// Remove removes a keyword for the given user
func (k *Keywords) Remove(word string, user int64) {
	for i, u := range k.kw[word] {
		if u == user {
			// Cut this user out of the slice
			k.kw[word] = append(k.kw[word][:i], k.kw[word][i+1:]...)
			// Remove the keyword if no other users are interested in it
			if len(k.kw[word]) == 0 {
				delete(k.kw, word)
			}
			return
		}
	}
}

// words returns all the words in a string as lowercase
func splitWords(text string) []string {
	text = strings.ToLower(text)
	words := strings.FieldsFunc(text, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})
	return words
}

// Find returns a slice of user IDs interested in at least one word
func (k *Keywords) Find(text string) []int64 {
	words := splitWords(text)
	matched := make(map[int64]bool)
	var users []int64
	for _, word := range words {
		wUsers := k.kw[word]
		for _, user := range wUsers {
			if matched[user] {
				// This user is already in our list
				continue
			}
			// This is a new user with a matching keyword
			matched[user] = true
			users = append(users, user)
		}
	}
	return users
}

// Match returns true if at least one word matches a keyword
func (k *Keywords) Match(text string) bool {
	words := splitWords(text)
	for _, word := range words {
		if _, ok := k.kw[word]; ok {
			return true
		}
	}
	return false
}
