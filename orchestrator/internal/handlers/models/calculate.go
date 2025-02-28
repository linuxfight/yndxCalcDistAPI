package models

type CalculateRequest struct {
	Expression string `json:"expression,required" validate:"expression,required"`
}

type CalculateResponse struct {
	Id string `json:"id,required"`
}
