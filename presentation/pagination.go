package presentation

// PaginationDTO represents the pagination metadata and the results
type PaginationDTO struct {
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalRows  int64       `json:"totalRows"`
	TotalPages int         `json:"totalPages"`
	Data       interface{} `json:"data"`
}
