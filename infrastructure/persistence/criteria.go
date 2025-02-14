package persistence

import (
	"github.com/jperdior/chatbot-kit/domain"
	"gorm.io/gorm"
)

func ApplyCriteria(query *gorm.DB, criteria domain.Criteria) (*gorm.DB, error) {
	var value interface{}

	for _, filter := range criteria.Filters() {
		query = query.Where(filter.Name()+" "+filter.Operation()+" ?", value)
	}

	if criteria.Sort() != "" {
		query = query.Order(criteria.Sort() + " " + criteria.SortDir())
	}

	page := criteria.Page()
	pageSize := criteria.PageSize()
	if pageSize > 0 {
		query = query.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	return query, nil
}
