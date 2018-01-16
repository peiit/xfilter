package bloom

import (
	"errors"
	"time"
)

type Filter struct {
	engine     *string
	concurency int
	timeout    *time.Duration
}

func New() (*Filter, error) {
	return &Filter{}, nil
}

func (f *Filter) Add()        {}
func (f *Filter) Remove()     {}
func (f *Filter) Get()        {}
func (f *Filter) Compare()    {}
func (f *Filter) Close()      {}
func (f *Filter) Sync() error { return errors.New("Not implemented yet") }
