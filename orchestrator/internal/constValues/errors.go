package constValues

import "errors"

var (
	NotFoundError          = errors.New("not found")
	ContentTypeError       = errors.New("invalid content type, must be application/json")
	InvalidJsonError       = errors.New("invalid json")
	InvalidExpressionError = errors.New("invalid expression")
	InvalidUuidError       = errors.New("invalid uuid")
)
