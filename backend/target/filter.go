// MIT License

// Copyright (c) 2016 rutcode-go

package filter

import (
	"sync"
	"time"

	"github.com/go-rut/errors"
)

type filter struct {
	filterFuncs map[string]FilterFunc

	sync.RWMutex
}

// New generate an new target filter
func New() FilterRepo {
	return &filter{
		filterFuncs: make(map[string]FilterFunc),
	}
}

type filterMessage struct {
	Name     string
	Filtered bool
	Err      error
}

// AddFilterFunc add filter functions
func (p *filter) AddFilterFunc(name string, cf FilterFunc) {
	p.Lock()
	defer p.Unlock()

	if name != "" && cf != nil {
		p.filterFuncs[name] = cf
	}
}

// RemoveFilterFunc remove a filter function
func (p *filter) RemoveFilterFunc(name string) {
	p.Lock()
	defer p.Unlock()

	if name != "" {
		delete(p.filterFuncs, name)
	}
}

// GetFilterFunc get a filter funcation by name
func (p *filter) GetFilterFunc(name string) FilterFunc {
	p.RLock()
	defer p.RUnlock()
	return p.filterFuncs[name]
}

// SetFilterTimeout set filter function timeout
func (p *filter) SetFilterTimeout(timeout int) error {
	p.Lock()
	defer p.Unlock()
	if timeout <= 0 {
		return ErrTimeoutMustAboveZero.New()
	}
	filterTimeout = time.Duration(timeout)
	return nil
}

// Compare
func (p *filter) Compare(params *FilterParams, req InputValues, cvs FilterValues) (filtered bool, err error) {

	if err = params.valid(); err != nil {
		return
	}

	fChan := make(chan filterMessage)

	// sequence to compare filter functions
	if params.Type == CompareTypeSequence {
		defer close(fChan)

		go func() {
			for _, name := range params.Names {
				fMessage := p.dofilter(name, req, cvs)

				if fMessage.Err != nil || fMessage.Filtered {
					fChan <- fMessage
					return
				}
			}
			fChan <- filterMessage{}
			return
		}()
	} else {
		// consistent to compare filter functions
		var wg sync.WaitGroup
		wg.Add(len(params.Names))
		for _, name := range params.Names {
			go func(n string) {
				defer wg.Done()
				fChan <- p.dofilter(n, req, cvs)
			}(name)
		}

		go func() {
			wg.Wait()
			close(fChan)
		}()

	}

	for i := 0; i < len(params.Names); i++ {
		select {
		case filter := <-fChan:
			if filter.Err != nil {
				return filter.Filtered, ErrFailedExecFilterFunction.New(errors.Params{"err": err.Error()})
			} else if filter.Filtered {
				return filter.Filtered, nil
			}
		case <-time.After(time.Second * filterTimeout):
			return true, ErrFailedExecTimeout.New()
		}
	}
	return
}

// dofilter execute filter function
func (p *filter) dofilter(name string, req InputValues, cvs FilterValues) filterMessage {
	fc := p.GetFilterFunc(name)
	if fc == nil {
		return filterMessage{Name: name}
	}

	f, e := fc(req, cvs)
	return filterMessage{Name: name, Filtered: f, Err: e}
}
