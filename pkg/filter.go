package xfilter

type Filter interface {
	// Init() (*Filter, error)
	Add()
	Remove()
	Get()
	Compare()
	Close()
	Sync() error
}
