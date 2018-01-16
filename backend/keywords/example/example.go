package main

import (
	"fmt"

	"github.com/transitorykris/go-keywords"
)

func main() {
	// Create our new keywords
	kw := keywords.New()

	// Add some users
	kw.Add("hello", 1)
	kw.Add("Hello", 1)
	kw.Add("hello", 2)
	kw.Add("world", 2)
	kw.Add("WORLD", 3)

	// See who matches this
	fmt.Println(kw.Find("Hello, world!"))            // [1 2 3]
	fmt.Println(kw.Find("This is an example World")) // [2 3]

	kw.Remove("world", 3)
	fmt.Println(kw.Find("This is an example World")) // [2]

	fmt.Println(kw.Match("Thanks for all the fish")) // false
	fmt.Println(kw.Match("So long, World!"))         // true
}
