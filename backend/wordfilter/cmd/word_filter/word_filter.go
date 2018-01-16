package main

import (
	"flag"
	word_filter "github.com/sniperkit/xfilter/backend/wordfilter"
)

var (
	configPath = flag.String("c", "", `config path`)
)

func main() {
	options := word_filter.NewOptions(*configPath)
	filter := word_filter.New(options)
	filter.Run(options.TCPAddr)
}
