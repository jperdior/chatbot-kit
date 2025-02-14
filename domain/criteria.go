package domain

type Criteria interface {
	GetFilters() []*Filter
	GetSort() string
	GetSortDir() string
	GetPage() int
	GetPageSize() int
}

type Filter interface {
	GetName() string
	GetType() string
	GetOperation() string
	GetValue() interface{}
}
