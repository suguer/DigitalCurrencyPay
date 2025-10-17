package mdb

type Pagination struct {
	Current  int64 `json:"current" form:"current,default=1"`
	PageSize int64 `json:"pageSize" form:"pageSize,default=10"`
	Total    int64 `json:"total" form:"total"`
}
