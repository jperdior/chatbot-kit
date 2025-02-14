package domain

type CriteriaInterface interface {
	Filters() []Filter
	Sort() string
	SortDir() string
	Page() int
	PageSize() int
}

type Criteria struct {
	filters  []Filter
	sort     string
	sortDir  string
	page     int
	pageSize int
}

func NewCriteria(filters []Filter, sort, sortDir string, page, pageSize int) *Criteria {
	return &Criteria{
		filters:  filters,
		sort:     sort,
		sortDir:  sortDir,
		page:     page,
		pageSize: pageSize,
	}
}

func (c *Criteria) Filters() []Filter {
	return c.filters
}

func (c *Criteria) Sort() string {
	return c.sort
}

func (c *Criteria) SortDir() string {
	return c.sortDir
}

func (c *Criteria) Page() int {
	return c.page
}

func (c *Criteria) PageSize() int {
	return c.pageSize
}

type FilterInterface interface {
	Name() string
	Operation() string
	Value() interface{}
}

type Filter struct {
	name      string
	operation string
	value     interface{}
}

func NewFilter(name, operation string, value interface{}) *Filter {
	return &Filter{
		name:      name,
		operation: operation,
		value:     value,
	}
}

func (f *Filter) Name() string {
	return f.name
}

func (f *Filter) Operation() string {
	return f.operation
}

func (f *Filter) Value() interface{} {
	return f.value
}
