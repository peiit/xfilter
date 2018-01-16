package main

import (
	"fmt"
	"sync"

	"golang.org/x/text/unicode/norm"

	"github.com/sniperkit/xfilter/backend/shamoji"
	"github.com/sniperkit/xfilter/backend/shamoji/filter"
	"github.com/sniperkit/xfilter/backend/shamoji/tokenizer"
)

var (
	o sync.Once
	s *shamoji.Serve
)

func main() {
	yes, word := Contains("No regret of one piece in my life")
	fmt.Printf("Result: %v, Word: %s", yes, word)
}

func Contains(sentence string) (bool, string) {
	o.Do(func() {
		s = &shamoji.Serve{
			Tokenizer: tokenizer.NewKagomeSimpleTokenizer(norm.NFKC),
			Filer:     filter.NewCuckooFilter("In the long run", "Repentance"),
		}
	})
	return s.Do(sentence)
}
