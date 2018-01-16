# target-filter
target filter

## Build

* [![Build Status](https://travis-ci.org/go-rut/target-filter.png)](https://travis-ci.org/go-rut/target-filter)

## a filter tool for customer's filter functions 


## Usage


### repo functions

```golang
type FilterRepo interface {
	// Set filter timeout: is second, and must be above zero.
	SetFilterTimeout(int) error
	// add a filter function : FilterFunc
	AddFilterFunc(name string, cf FilterFunc)
	// remove a filter funcation by name
	RemoveFilterFunc(name string)
	// get filter function
	GetFilterFunc(name string) FilterFunc
	// compare input values and filter values
	Compare(*FilterParams, InputValues, FilterValues) (bool, error)
}

type FilterParams struct {
	// filter names
	Names []string
	// compare type: Sequence (0) or Consistent (1)
	Type  CompareType
}
```


### new a filter

```golang
	manager := filter.New()
	manager.AddFilterFunc(filterName, filterNameFunction)
```


### do sequence filter

```golang
fParams := &filter.FilterParams{Names: []string{"filterName"}
manager.Compare(fParams, nil, nil)
```

Or

```golang
fParams := &filter.FilterParams{Type: filter.CompareTypeSequence, Names: []string{"filterName"}
manager.Compare(fParams, nil, nil)
```

### do consistent filter

```golang
fParams := &filter.FilterParams{Type: filter.CompareTypeConsistent, Names: []string{"filterName"}
manager.Compare(fParams, nil, nil)
```

### Get more

[test_file](filter_test.go)