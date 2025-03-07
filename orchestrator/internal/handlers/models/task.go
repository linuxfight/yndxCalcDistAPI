package models

type TaskResponse struct {
	ID            string  `json:"id" example:"928b303f-cfcc-46f4-ae24-aabb72bbb7d9"`
	Arg1          float64 `json:"arg1" example:"1"`
	Arg2          float64 `json:"arg2" example:"1"`
	Operation     string  `json:"operation" example:"+"`
	OperationTime int     `json:"operation_time" example:"1000"`
	ExpressionId  string  `json:"expression_id" example:"928b303f-cfcc-46f4-ae24-aabb72bbb7d9"`
}

type InternalTask struct {
	ID        string      `json:"id" example:"928b303f-cfcc-46f4-ae24-aabb72bbb7d9"`
	Arg1      interface{} `json:"arg1" example:"1"`
	Arg2      interface{} `json:"arg2" example:"928b303f-cfcc-46f4-ae24-aabb72bbb7d9"`
	Operation string      `json:"operation" example:"-"`
	Result    interface{} `json:"result" example:"0"`
}

type TaskRequest struct {
	ID     string      `json:"id" example:"928b303f-cfcc-46f4-ae24-aabb72bbb7d9"`
	Result interface{} `json:"result"`
}
