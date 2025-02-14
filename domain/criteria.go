package domain

type CriteriaInterface interface {
	GetFilters() []*Filter
	GetSort() string
	GetSortDir() string
	GetPage() int
	GetPageSize() int
}

type Criteria struct {
	filters  []*Filter
	sort     string
	sortDir  string
	page     int
	pageSize int
}

func NewCriteria(filters []*Filter, sort, sortDir string, page, pageSize int) *Criteria {
	return &Criteria{
		filters:  filters,
		sort:     sort,
		sortDir:  sortDir,
		page:     page,
		pageSize: pageSize,
	}
}

func (c *Criteria) GetFilters() []*Filter {
	return c.filters
}

func (c *Criteria) GetSort() string {
	return c.sort
}

func (c *Criteria) GetSortDir() string {
	return c.sortDir
}

func (c *Criteria) GetPage() int {
	return c.page
}

func (c *Criteria) GetPageSize() int {
	return c.pageSize
}

type Filter struct {
	Name      string
	Type      string
	Operation string
	Value     interface{}
}

func NewFilter(name, t, operation string, value interface{}) *Filter {
	return &Filter{
		Name:      name,
		Type:      t,
		Operation: operation,
		Value:     value,
	}
}
