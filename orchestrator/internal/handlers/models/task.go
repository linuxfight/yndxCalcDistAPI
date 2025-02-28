package models

type TaskResponse struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
	ExpressionId  string  `json:"expression_id"`
}

type InternalTask struct {
	ID            string      `json:"id"`
	Arg1          interface{} `json:"arg1"`
	Arg2          interface{} `json:"arg2"`
	Operation     string      `json:"operation"`
	OperationTime int         `json:"operation_time"`
	Result        interface{} `json:"result"`
}

type TaskRequest struct {
	ID     string      `json:"id"`
	Result interface{} `json:"result"`
}
