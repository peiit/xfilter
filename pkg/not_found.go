package xfilter

import (
	"errors"
)

func (nf *NotFound) Add()     { return nf }
func (nf *NotFound) Remove()  {}
func (nf *NotFound) Get()     {}
func (nf *NotFound) Compare() {}
func (nf *NotFound) Close()   {}
func (f *Filter) Sync() error { return errors.New("Not implemented yet") }
