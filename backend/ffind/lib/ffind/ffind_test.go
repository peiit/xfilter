// This file is part of ffind.
//
// Copyright (C) 2017  David Gamba Rios
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ffind

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func goToRootDir() {
	cwd, _ := os.Getwd()
	log.Printf("CWD: %s", cwd)
	if strings.HasSuffix(cwd, "lib/ffind") {
		// Go to the root of the repository
		_ = os.Chdir("../..")
	}
	cwd, _ = os.Getwd()
	log.Printf("CWD: %s", cwd)
}

func compareTestStringSlices(t *testing.T, expected []string, received []string) {
	if !reflect.DeepEqual(received, expected) {
		str := ""
		for i := range expected {
			switch {
			case len(received) <= i:
				str += fmt.Sprintf("%s\n", expected[i])
			case expected[i] != received[i]:
				str += fmt.Sprintf("%s != %s\n", expected[i], received[i])
			default:
				str += fmt.Sprintf("%s == %s\n", expected[i], received[i])
			}
		}
		if len(received) > len(expected) {
			for i := len(expected); i < len(received); i++ {
				str += fmt.Sprintf(" == %s\n", received[i])
			}
		}
		t.Fatalf("trees differ: expected (%d) received (%d)\n%s\n", len(expected), len(received), str)
	}
}

func TestNewFileError(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	// log.SetOutput(os.Stderr)
	goToRootDir()
	file := "./test_tree"
	fe, err := NewFileError(file)
	if err != nil {
		t.Fatalf("Unexpected error '%s': %s\n", file, err)
	}
	if fe.Path != "test_tree" {
		t.Fatalf("Paths don't match '%s': %s\n", file, "test_tree")
	}
	file = "non-existent"
	fe, err = NewFileError(file)
	if err == nil {
		t.Fatalf("Expected error, none received!\n")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("Expected IsNotExist error, received: %s\n", err)
	}
	log.Printf("%#v\n", fe)
}

// Make sure filepath.EvalSymlinks behaves as expected.
func TestEvalSymlink(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	goToRootDir()

	cases := []struct {
		file     string
		expected string
	}{
		{"./test_tree", "test_tree"},
		{"./test_tree/.A", "test_tree/.A"},
		{"./test_tree/.A/b/C/d/E", "test_tree/.A/b/C/d/E"},
		{"./test_tree/slnA", "test_tree/A"},
		{"./test_tree/slnB", "test_tree/A"},
	}
	for _, c := range cases {
		// Read given file information
		read, err := filepath.EvalSymlinks(c.file)
		if err != nil {
			t.Fatalf("Unexpected error '%s': %s\n", c.file, err)
		}
		if read != c.expected {
			t.Fatalf("Got %s != Expected %s\n", read, c.expected)
		}
	}
	_, err := filepath.EvalSymlinks("./test_tree/slnE")
	if err == nil {
		t.Fatalf("Expected error, none received! \n")
	}
	if !strings.Contains(err.Error(), "too many links") {
		t.Fatalf("Unexpected error: %s\n", err)
	}
}

func TestListOneLevel(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	// log.SetOutput(os.Stderr)
	goToRootDir()
	cases := []struct {
		file     string
		follow   bool
		expected []string
	}{
		{"./test_tree", false, []string{
			"test_tree/.A",
			"test_tree/.a",
			"test_tree/.hg",
			"test_tree/.svn",
			"test_tree/A",
			"test_tree/a",
			"test_tree/slnA",
			"test_tree/slnB",
			"test_tree/slnC",
			"test_tree/slnD",
			"test_tree/slnE",
			"test_tree/slnF",
			"test_tree/slnG",
			"test_tree/z",
		},
		},
		{"./test_tree/.A", false, []string{
			"test_tree/.A/b",
		},
		},
		{"./test_tree/.A/b/C/d/E", false, []string{
			"test_tree/.A/b/C/d/E",
		},
		},
		{"./test_tree/slnA", false, []string{
			"test_tree/slnA",
		},
		},
		{"./test_tree/slnA", true, []string{
			"test_tree/slnA/b",
		},
		},
		{"./test_tree/slnB", false, []string{
			"test_tree/slnB",
		},
		},
		{"./test_tree/slnB", true, []string{
			"test_tree/slnB/b",
		},
		},
	}
	for _, c := range cases {
		ch := ListOneLevel(c.file, c.follow, SortFnByName)
		tree := []string{}
		for e := range ch {
			if e.Error != nil {
				t.Fatalf("Unexpected error: %s\n", e.Error)
			}
			tree = append(tree, e.Path)
		}
		compareTestStringSlices(t, c.expected, tree)
	}
	file := "./test_tree/slnE"
	ch := ListOneLevel(file, true, SortFnByName)
	log.Printf("%d\n", len(ch))
	for e := range ch {
		if e.Error == nil {
			t.Fatalf("Expected error, none received! \n")
		}
	}
	file = "tmp-no-permissions"
	err := os.Mkdir(file, 0000)
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	ch = ListOneLevel(file, false, SortFnByName)
	tree := []string{}
	for e := range ch {
		if e.Error == nil {
			t.Errorf("Expected error, none received! \n")
		}
		if !os.IsPermission(e.Error) {
			t.Errorf("Expected IsPermission error, received: %s\n", e.Error.Error())
		}
		tree = append(tree, e.Path)
	}
	err = os.RemoveAll(file)
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	compareTestStringSlices(t, []string{file}, tree)
	file = "non-existent"
	ch = ListOneLevel(file, false, SortFnByName)
	for e := range ch {
		if e.Error == nil {
			t.Errorf("Expected error, none received! \n")
		}
		if !os.IsNotExist(e.Error) {
			t.Errorf("Expected IsNotExist error, received: %s\n", e.Error)
		}
	}
}

var listRecursiveCases = []struct {
	file              string
	ignoreDirResults  bool
	ignoreFileResults bool
	ignoreCVSDirs     bool
	expected          []string
}{
	{"./test_tree/.A", false, false, false, []string{
		"test_tree/.A/b",
		"test_tree/.A/b/C",
		"test_tree/.A/b/C/d",
		"test_tree/.A/b/C/d/E",
	},
	},
	{"./test_tree", false, false, false, []string{
		"test_tree/.A",
		"test_tree/.A/b",
		"test_tree/.A/b/C",
		"test_tree/.A/b/C/d",
		"test_tree/.A/b/C/d/E",
		"test_tree/.a",
		"test_tree/.a/B",
		"test_tree/.a/B/c",
		"test_tree/.a/B/c/D",
		"test_tree/.a/B/c/D/e",
		"test_tree/.hg",
		"test_tree/.hg/E",
		"test_tree/.hg/e",
		"test_tree/.svn",
		"test_tree/.svn/E",
		"test_tree/.svn/e",
		"test_tree/A",
		"test_tree/A/b",
		"test_tree/A/b/C",
		"test_tree/A/b/C/d",
		"test_tree/A/b/C/d/E",
		"test_tree/a",
		"test_tree/a/B",
		"test_tree/a/B/c",
		"test_tree/a/B/c/D",
		"test_tree/a/B/c/D/e",
		"test_tree/slnA",
		"test_tree/slnA/b",
		"test_tree/slnA/b/C",
		"test_tree/slnA/b/C/d",
		"test_tree/slnA/b/C/d/E",
		"test_tree/slnB",
		"test_tree/slnB/b",
		"test_tree/slnB/b/C",
		"test_tree/slnB/b/C/d",
		"test_tree/slnB/b/C/d/E",
		"test_tree/slnC",
		"test_tree/slnD",
		"test_tree/slnE",
		"test_tree/slnF",
		"test_tree/slnG",
		"test_tree/z",
	},
	},
	{"./test_tree", true, false, false, []string{
		"test_tree/.A/b/C/d/E",
		"test_tree/.a/B/c/D/e",
		"test_tree/.hg/E",
		"test_tree/.hg/e",
		"test_tree/.svn/E",
		"test_tree/.svn/e",
		"test_tree/A/b/C/d/E",
		"test_tree/a/B/c/D/e",
		"test_tree/slnA/b/C/d/E",
		"test_tree/slnB/b/C/d/E",
		"test_tree/slnC",
		"test_tree/slnD",
		"test_tree/slnE",
		"test_tree/slnF",
		"test_tree/slnG",
		"test_tree/z",
	},
	},
	{"./test_tree", false, true, false, []string{
		"test_tree/.A",
		"test_tree/.A/b",
		"test_tree/.A/b/C",
		"test_tree/.A/b/C/d",
		"test_tree/.a",
		"test_tree/.a/B",
		"test_tree/.a/B/c",
		"test_tree/.a/B/c/D",
		"test_tree/.hg",
		"test_tree/.svn",
		"test_tree/A",
		"test_tree/A/b",
		"test_tree/A/b/C",
		"test_tree/A/b/C/d",
		"test_tree/a",
		"test_tree/a/B",
		"test_tree/a/B/c",
		"test_tree/a/B/c/D",
		"test_tree/slnA",
		"test_tree/slnA/b",
		"test_tree/slnA/b/C",
		"test_tree/slnA/b/C/d",
		"test_tree/slnB",
		"test_tree/slnB/b",
		"test_tree/slnB/b/C",
		"test_tree/slnB/b/C/d",
	},
	},
	{"./test_tree", false, false, true, []string{
		"test_tree/.A",
		"test_tree/.A/b",
		"test_tree/.A/b/C",
		"test_tree/.A/b/C/d",
		"test_tree/.A/b/C/d/E",
		"test_tree/.a",
		"test_tree/.a/B",
		"test_tree/.a/B/c",
		"test_tree/.a/B/c/D",
		"test_tree/.a/B/c/D/e",
		"test_tree/A",
		"test_tree/A/b",
		"test_tree/A/b/C",
		"test_tree/A/b/C/d",
		"test_tree/A/b/C/d/E",
		"test_tree/a",
		"test_tree/a/B",
		"test_tree/a/B/c",
		"test_tree/a/B/c/D",
		"test_tree/a/B/c/D/e",
		"test_tree/slnA",
		"test_tree/slnA/b",
		"test_tree/slnA/b/C",
		"test_tree/slnA/b/C/d",
		"test_tree/slnA/b/C/d/E",
		"test_tree/slnB",
		"test_tree/slnB/b",
		"test_tree/slnB/b/C",
		"test_tree/slnB/b/C/d",
		"test_tree/slnB/b/C/d/E",
		"test_tree/slnC",
		"test_tree/slnD",
		"test_tree/slnE",
		"test_tree/slnF",
		"test_tree/slnG",
		"test_tree/z",
	},
	},
}

func TestListRecursive(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	// log.SetOutput(os.Stderr)
	goToRootDir()
	for _, c := range listRecursiveCases {
		ch := ListRecursive(
			c.file,
			true,
			&BasicFileMatch{
				IgnoreDirResults:  c.ignoreDirResults,
				IgnoreFileResults: c.ignoreFileResults,
				IgnoreVCSDirs:     c.ignoreCVSDirs,
			},
			SortFnByName)
		tree := []string{}
		for e := range ch {
			if e.Error != nil {
				if !strings.Contains(e.Error.Error(), "too many links") {
					t.Fatalf("Unexpected error: %s\n", e.Error)
				}
			}
			tree = append(tree, e.Path)
		}
		compareTestStringSlices(t, c.expected, tree)
	}
	file := "tmp-no-permissions"
	err := os.Mkdir(file, 0000)
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	ch := ListRecursive(file, false, &BasicFileMatch{}, SortFnByName)
	tree := []string{}
	for e := range ch {
		if e.Error == nil {
			t.Errorf("Expected error, none received! \n")
		}
		if !os.IsPermission(e.Error) {
			t.Errorf("Expected IsPermission error, received: %s\n", e.Error.Error())
		}
		tree = append(tree, e.Path)
	}
	err = os.RemoveAll(file)
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	compareTestStringSlices(t, []string{file}, tree)
	file = "non-existent"
	ch = ListRecursive(file, false, &BasicFileMatch{}, SortFnByName)
	for e := range ch {
		if e.Error == nil {
			t.Errorf("Expected error, none received! \n")
		}
		if !os.IsNotExist(e.Error) {
			t.Errorf("Expected IsNotExist error, received: %s\n", e.Error)
		}
	}
	file = "test_files"
	ch = ListRecursive(
		file,
		false,
		&BasicFileMatch{
			MatchFileTypeList: []string{"ruby"},
		},
		SortFnByName)
	tree = []string{}
	for e := range ch {
		if e.Error != nil {
			t.Fatalf("Unexpected error: %s\n", err)
		}
		tree = append(tree, e.Path)
	}
	compareTestStringSlices(t, []string{"test_files/file.rb"}, tree)
}

func TestListRecursiveWalk(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	// log.SetOutput(os.Stderr)
	goToRootDir()
	for _, c := range listRecursiveCases {
		tree := []string{}
		myWalkFn := func(path string, info os.FileInfo, err error) error {
			// walkFn is called with the root itself.
			// Ignore the root unless it is a file.
			if path == c.file {
				if info.IsDir() {
					return nil
				}
			}
			tree = append(tree, path)
			if err != nil {
				if !strings.Contains(err.Error(), "too many links") {
					t.Fatalf("Unexpected error: %s\n", err)
					return err
				}
			}
			return nil
		}
		s := &BasicFileMatch{
			IgnoreDirResults:  c.ignoreDirResults,
			IgnoreFileResults: c.ignoreFileResults,
			IgnoreVCSDirs:     c.ignoreCVSDirs,
		}
		err := ListRecursiveWalk(c.file, true, s, myWalkFn)
		if err != nil {
			if !strings.Contains(err.Error(), "too many links") {
				t.Fatalf("Unexpected error: %s\n", err)
			}
		}
		compareTestStringSlices(t, c.expected, tree)
	}
}
