package ginx

type Meta struct {
	CurrentPage uint `json:"currentPage"`
	Total       uint `json:"total"`
	LastPage    uint `json:"lastPage"`
	PerPage     uint `json:"perPage"`
}

type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type ResponseWithCode struct {
	Code
	Response
}

type Code struct {
	Code int `json:"code"`
}
