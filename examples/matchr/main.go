package main

import (
	"github.com/sniperkit/xfilter/backend/matchr"
)

var mtchr *matchr.String

func main() {
	mtchr = matchr.NewString("contents")
}
