package dtos

type PaginatedData struct {
	Page       int64       `json:"page,omitempty" example:"1"`
	PerPage    int64       `json:"per_page,omitempty" example:"10"`
	Total      int64       `json:"total,omitempty" example:"100"`
	TotalPages int         `json:"total_pages,omitempty" example:"10"`
	Rows       interface{} `json:"rows,omitempty" swaggertype:"array,object"`
}
