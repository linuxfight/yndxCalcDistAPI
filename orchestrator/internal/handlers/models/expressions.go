package models

type ListAllExpressionsResponse struct {
	Expressions []Expression `json:"expressions"`
}

type GetByIdExpressionResponse struct {
	Expression Expression `json:"expression"`
}

type Expression struct {
	Id     string  `json:"id" example:"928b303f-cfcc-46f4-ae24-aabb72bbb7d9"`
	Result float64 `json:"result"`
	Status string  `json:"status" example:"DONE"`
}
