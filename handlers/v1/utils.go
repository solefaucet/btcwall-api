package v1

type paginationPayload struct {
	Limit  int64 `form:"limit" binding:"min=0,max=100"`
	Offset int64 `form:"offset" binding:"min=0"`
}

type paginationResult struct {
	Count int64       `json:"count"`
	Data  interface{} `json:"data"`
}
