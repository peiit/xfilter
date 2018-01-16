// MIT License

// Copyright (c) 2016 rutcode-go

package filter

import (
	"github.com/go-rut/errors"
)

const (
	namespace = "target-filter"
)

var (
	ErrInvalidFilterName        = errors.TN(namespace, 1000, "invalid filter name")
	ErrFilterFunctionEqualNil   = errors.TN(namespace, 1001, "filter function should not be nil")
	ErrNotSupportedFilterType   = errors.TN(namespace, 1002, "filter type not supported")
	ErrNotInputParams           = errors.TN(namespace, 1003, "not input params")
	ErrFailedExecFilterFunction = errors.TN(namespace, 1004, "failed exec filter function: {{.err}}")
	ErrFailedExecTimeout        = errors.TN(namespace, 1005, "exec filter function timeout")
	ErrTimeoutMustAboveZero     = errors.TN(namespace, 1006, "set timeout must above zero")
)
