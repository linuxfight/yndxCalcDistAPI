package models

type TaskResponse struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

type TaskRequest struct {
	ID     string `json:"id"`
	Result string `json:"result"`
}
