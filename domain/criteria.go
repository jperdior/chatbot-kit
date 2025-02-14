package domain

type Criteria interface {
	Filters() []*Filter
	Sort() string
	SortDir() string
	Page() int
	PageSize() int
}

type Filter interface {
	Name() string
	Type() string
	Operation() string
	Value() interface{}
}
