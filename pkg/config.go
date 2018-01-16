package xfilter

import (
	"time"
)

const (
	DefaultTimeout         int = 120
	DefaultConcurrency     int = 4
	DefaultMemoryCacheSize int = 2048
)

var (
	DefaultReferencesTypes []string = []string{"list", "corpus", "dictionary", "lexicon", "thesaurus", "tag", "pattern", "regex"}
	DefaultEvents          []string = []string{"notify", "remove", "replace", "skip", "ignore", "notify", "annotate"}
)

/*
	Refs:
	- https://github.com/ymhhh/target-filter/blob/master/compare.go
*/

type Config struct {
	Backend  string            `yaml:"backend" json:"backend" toml:"backend" bson:"backend" csv:"backend" xml:"backend"`
	Timeout  time.Duration     `yaml:"timeout" json:"timeout" toml:"timeout" bson:"timeout" csv:"timeout" xml:"timeout" default:"120"`
	Encoding string            `yaml:"encoding" json:"encoding" toml:"encoding" bson:"encoding" csv:"encoding" xml:"encoding" default:"UTF-8"`
	Filters  map[string]Filter `yaml:"filters" json:"filters" toml:"filters" bson:"filters" csv:"filters" xml:"filters"`
	// FilterFuncs map[string]FilterFunc `yaml:"filter_funcs" json:"filter_funcs" toml:"filter_funcs" bson:"filter_funcs" csv:"filter_funcs" xml:"filter_funcs"`
}

type Filters map[string]interface{}

func (f Filters) NormalizeValues() Filters {
	for key, val := range f {
		switch v := val.(type) {
		default:
			f[key] = v
		}
	}
	return f
}

//FilterType is a type identifier for logger fields
type FilterType int8

//Filter is a struct to send paramaters to log messages
type Filter struct {
	key     string
	val     interface{}
	valType FilterType
}

// type FilterFunc func(req InputValues, cvs FilterValues) (filtered bool, err error)
