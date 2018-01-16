package cuckoofilter

import (
	"errors"
	"hash"
	"hash/fnv"
	"time"
)

type Filter struct {
	engine     *CuckooFilter
	concurency int
	timeout    *time.Duration
}

type Config struct {
	capacity uint
}

func New() (*Filter, error) {
	return &Filter{
		engine: NewDefaultCuckooFilter(),
	}, nil
}

func (f *Filter) Add(key []byte) {
	// inserts 'key' to the filter
	f.engine.InsertUnique(key)
}

func (f *Filter) Remove(key []byte) {
	// tries deleting 'key' from filter, may delete another element
	// this could occur when another byte slice with the same fingerprint
	// as another is 'deleted'
	f.engine.Delete(key)
}

func (f *Filter) Get(key []byte) {
	// Lookup a string (and it a miss) if it exists in the cuckoofilter
	f.engine.Lookup(key)
}

func (f *Filter) Count() {
	// returns 1 (given only 'buongiorno' was added)
	f.engine.Count()
}

func (f *Filter) Compare()    {}
func (f *Filter) Close()      {}
func (f *Filter) Sync() error { return errors.New("Not implemented yet") }
