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

func ApplyCriteriaWithCount(query *gorm.DB, criteria domain.CriteriaInterface) (*gorm.DB, int64, error) {
	var total int64

	// Clone the query for counting total records
	countQuery := query.Model(query.Statement.Model)

	// Apply only filters to count query, not pagination or sorting
	for _, filter := range criteria.Filters() {
		value := filter.Value()
		countQuery = countQuery.Where(filter.Name()+" "+filter.Operation()+" ?", value)
	}

	// Get total count without pagination
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply all criteria including pagination to the data query
	query, err := ApplyCriteria(query, criteria)
	if err != nil {
		return nil, 0, err
	}

	return query, total, nil
}
