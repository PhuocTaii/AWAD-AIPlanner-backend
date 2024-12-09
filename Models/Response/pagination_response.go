package response

type PagingResponse struct {
	Data        interface{} `json:"data"`
	CurrentPage int         `json:"current_page"`
	PerPage     int         `json:"per_page"`
	TotalPage   int         `json:"total_page"`
	TotalItems  int         `json:"total_items"`
}
