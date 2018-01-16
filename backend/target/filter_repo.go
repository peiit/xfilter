package target

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
