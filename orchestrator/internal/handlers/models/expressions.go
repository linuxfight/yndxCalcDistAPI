package models

type ListAllExpressionsResponse struct {
	Expressions []Expression `json:"expressions"`
}

type GetByIdExpressionResponse struct {
	Expression Expression `json:"expression"`
}

type Expression struct {
	Id     string  `json:"id"`
	Result float64 `json:"result"`
	Status string  `json:"status"`
}
