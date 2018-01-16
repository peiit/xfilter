package keywords

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMatchedUsers(t *testing.T) {
	Convey("When creating a new Keywords", t, func() {
		Convey("And it is successful", func() {
			kw := New()
			So(*kw, ShouldNotBeNil)
		})
	})

	Convey("When adding a user to a keyword", t, func() {
		Convey("And it's a new keywords", func() {
			kw := New()
			kw.Add("hello", 1)
			So(kw.kw["hello"], ShouldResemble, []int64{1})
		})
		Convey("And there is already another user", func() {
			kw := New()
			kw.Add("hello", 1)
			kw.Add("HELLo", 2)
			So(kw.kw["hello"], ShouldResemble, []int64{1, 2})
		})
		Convey("And the user already has that keyword", func() {
			kw := New()
			kw.Add("hello", 1)
			kw.Add("hELlO", 1)
			So(kw.kw["hello"], ShouldResemble, []int64{1})
			So(kw.kw["hELlO"], ShouldBeNil)
		})
	})

	Convey("When removing a user from a keyword", t, func() {
		Convey("And the keyword does not exist", func() {
			kw := New()
			kw.Remove("nothere", 1)
			So(kw.kw["nothere"], ShouldBeNil)
		})
		Convey("And the keyword exists but the user is not in it", func() {
			kw := New()
			kw.Add("hello", 1)
			kw.Remove("hello", 2)
			So(kw.kw["hello"], ShouldResemble, []int64{1})
		})
		Convey("And they are the only ones in the list", func() {
			kw := New()
			kw.Add("hello", 1)
			kw.Remove("hello", 1)
			So(kw.kw["hello"], ShouldBeNil)
		})
		Convey("And there are others in the list too", func() {
			kw := New()
			kw.Add("hello", 1)
			kw.Add("hello", 2)
			kw.Remove("hello", 1)
			So(kw.kw["hello"], ShouldResemble, []int64{2})
		})
	})

	Convey("Given a line of text", t, func() {
		Convey("And there is no matching user", func() {
			kw := New()
			users := kw.Find("This line does not match anything")
			So(users, ShouldBeNil)
		})
		Convey("And the text is empty", func() {
			kw := New()
			users := kw.Find("")
			So(users, ShouldBeNil)
		})
		Convey("And there is a single matching user", func() {
			kw := New()
			kw.Add("hello", 1)
			kw.Add("Keywords", 1)
			users := kw.Find("Hello, Keywords!")
			So(users, ShouldResemble, []int64{1})
		})
		Convey("And there are multiple matching users", func() {
			kw := New()
			kw.Add("hello", 1)
			kw.Add("keywords", 2)
			kw.Add("keywords", 3)
			users := kw.Find("Hello, Keywords!")
			So(users, ShouldResemble, []int64{1, 2, 3})
		})
	})

	Convey("When checking if a line matches or not", t, func() {
		Convey("And it does not match", func() {
			kw := New()
			matched := kw.Match("Nothing will match this")
			So(matched, ShouldBeFalse)
		})
		Convey("And it does match", func() {
			kw := New()
			kw.Add("match", 1)
			matched := kw.Match("This will match!")
			So(matched, ShouldBeTrue)
		})
		Convey("And the string is empty", func() {
			kw := New()
			matched := kw.Match("")
			So(matched, ShouldBeFalse)
		})
	})
}

// Adds the id 1 to the keyword _hello_
func ExampleKeywords_Add() {
	kw := New()
	kw.Add("hello", 1) // kw has now associated _hello_ with id 1
}

// Removes the id 1 from the keyword _hello_
func ExampleKeywords_Remove() {
	kw := New()
	kw.Add("hello", 1)
	kw.Remove("hello", 1) // kw is empty again
}

// Finds all ids interested in keywords in the given text
func ExampleKeywords_Find() {
	kw := New()
	kw.Add("hello", 1)
	kw.Add("world", 2)
	fmt.Println(kw.Find("Hello World"))
	// Output:
	// [1 2]
}

// Match returns true if at least one id is interested in a keyword
func ExampleKeywords_Match() {
	kw := New()
	kw.Add("hello", 1)
	fmt.Println(kw.Match("Hello World"))
	fmt.Println(kw.Match("So long"))
	// Output:
	// true
	// false
}
