package models

type CalculateRequest struct {
	Expression string `json:"expression,required" validate:"expression,required" example:"2+2"`
}

type CalculateResponse struct {
	Id string `json:"id,required" example:"928b303f-cfcc-46f4-ae24-aabb72bbb7d9"`
}
