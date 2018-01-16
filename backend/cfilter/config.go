package cfilter

import (
	"errors"
	"hash"
	"hash/fnv"
	"time"

	core "github.com/irfansharif/cfilter"
)

type Filter struct {
	engine     *core.CFilter
	concurency int
	timeout    *time.Duration
}

type Config struct {
	Size            *uint
	BucketSize      *uint8
	FingerprintSize *uint8
	MaximumKicks    *uint
	HashFn          *hash.Hash
}

func New(cfg *Config) (*Filter, error) {
	return &Filter{
		engine: core.New(),
	}, nil
}

func (f *Filter) Init(conf *Config) (*Filter, error) {
	if f.engine == nil {
		f.engine = engine
		if conf.Size != nil {
			f.engine.Size(conf.Size) // cfilter.Size(uint) sets the number of buckets in the filter
		}
		if conf.BucketSize != nil {
			f.engine.BucketSize(conf.BucketSize) // cfilter.BucketSize(uint8) sets the size of each bucket
		}
		if conf.FingerprintSize != nil {
			f.engine.FingerprintSize(conf.FingerprintSize) // cfilter.FingerprintSize(uint8) sets the size of the fingerprint
		}
		if conf.MaximumKicks != nil {
			f.engine.MaximumKicks(conf.MaximumKicks) // cfilter.MaximumKicks(uint) sets the maximum number of bucket kicks
		}
		if conf.HashFn != nil {
			f.engine.HashFn(conf.HashFn) // cfilter.HashFn(hash.Hash) sets the fingerprinting hashing function
		}
		return true, nil
	}
	return false, errors.New("filter engine already initialized")
}

func (f *Filter) Add(key []byte) {
	// inserts 'key' to the filter
	f.engine.Insert(key)
}

func (f *Filter) Remove(key []byte) {
	// tries deleting 'key' from filter, may delete another element
	// this could occur when another byte slice with the same fingerprint
	// as another is 'deleted'
	f.engine.Delete(key)
}

func (f *Filter) Get(key []byte) {
	// looks up 'key' in the filter, may return false positive
	f.engine.Lookup(key)
}

func (f *Filter) Count() {
	// returns 1 (given only 'buongiorno' was added)
	f.engine.Count()
}

func (f *Filter) Compare()    {}
func (f *Filter) Close()      {}
func (f *Filter) Sync() error { return errors.New("Not implemented yet") }
