package calc

import "errors"

var (
	divisionByZero   = errors.New("division by zero")
	invalidCharacter = errors.New("invalid character in expression")
)
