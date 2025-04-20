package gorm

import (
	"github.com/jperdior/chatbot-kit/domain"
	"gorm.io/gorm"
)

func ApplyCriteria(query *gorm.DB, criteria domain.CriteriaInterface) (*gorm.DB, error) {
	for _, filter := range criteria.Filters() {
		value := filter.Value()
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

func ApplyCriteriaWithoutPagination(query *gorm.DB, criteria domain.CriteriaInterface) (*gorm.DB, error) {
	for _, filter := range criteria.Filters() {
		value := filter.Value()
		query = query.Where(filter.Name()+" "+filter.Operation()+" ?", value)
	}

	if criteria.Sort() != "" {
		query = query.Order(criteria.Sort() + " " + criteria.SortDir())
	}
	return query, nil
}
