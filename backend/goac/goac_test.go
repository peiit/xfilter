package goac

import (
	"fmt"
	"testing"
)

func TestAhoCorasick(t *testing.T) {
	content := "Aho-Corasick is a deterministic finite state automata algorithm based on Trie tree"
	ac := NewAhoCorasick()
	ac.AddPattern("tree")
	ac.AddPattern("Automaton")
	ac.AddPattern("Trie")
	ac.AddPattern("based on")
	ac.Build()
	results := ac.Scan(content)
	fmt.Println("content: " + content)
	fmt.Println("Matching words: ")
	for _, result := range results {
		fmt.Println(string([]rune(content)[result.Start : result.End+1]))
	}
}
