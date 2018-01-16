// MIT License

// Copyright (c) 2016 rutcode-go

package filter

import (
	"time"
)

type InputValues map[string]interface{}
type FilterValues map[string]interface{}

// FilterFunc is the function you want to compare input values with target values
// InputValues is the k:v maps with input values
// FilterValues is the k:v maps with filter values
// return bool is filtered
// return error if compare function has internal error
// you can write you campare functions like FilterFunc, then manager.AddFilterFunc(name, filterFunc)
type FilterFunc func(req InputValues, cvs FilterValues) (filtered bool, err error)

type CompareType int

const (
	CompareTypeSequence CompareType = iota
	CompareTypeConsistent
)

var (
	filterTimeout time.Duration = 30
)

type FilterParams struct {
	Names []string
	Type  CompareType
}

func (p *FilterParams) valid() (err error) {
	if p == nil {
		return ErrNotInputParams.New()
	}

	if err = p.validType(); err != nil {
		return
	}

	return p.validNames()
}

func (p *FilterParams) validType() error {
	if p.Type != CompareTypeSequence &&
		p.Type != CompareTypeConsistent {
		return ErrNotSupportedFilterType.New()
	}
	return nil
}

func (p *FilterParams) validNames() error {
	if len(p.Names) == 0 {
		return ErrInvalidFilterName.New()
	}
	return nil
}
