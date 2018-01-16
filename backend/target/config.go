package target

import (
	"errors"
	"time"
)

type Filter struct {
	backend    *FilterRepo
	concurency int
	timeout    *time.Duration
}

func NewFilter() (*Filter, error) {

	factory := &Filter{}
	factory.backend = New()

	return factory, nil
}

func (f *Filter) Init(engine *FilterRepo) (*Filter, error) {
	if f.engine == nil {
		f.engine = engine
		return true, nil
	}
	return false, errors.New("filter engine already initialized")
}

func (f *Filter) Add()        {}
func (f *Filter) Remove()     {}
func (f *Filter) Get()        {}
func (f *Filter) Compare()    {}
func (f *Filter) Close()      {}
func (f *Filter) Sync() error { return errors.New("Not implemented yet") }
