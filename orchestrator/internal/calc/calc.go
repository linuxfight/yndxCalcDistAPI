package calc

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"orchestrator/internal/handlers/models"
	"strconv"
	"strings"
)

// tokenize splits the expression into tokens.
func tokenize(expression string) ([]string, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	expression = strings.ReplaceAll(expression, ",", ".")
	if strings.Contains(expression, "/0") {
		return nil, divisionByZero
	}
	var tokens []string
	var current string
	for i, r := range expression {
		if r == '+' || r == '-' || r == '*' || r == '/' || r == '(' || r == ')' {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
			// Handle negative numbers.
			if r == '-' && (i == 0 || strings.ContainsRune("+-*/(", rune(expression[i-1]))) {
				current += string(r)
			} else {
				tokens = append(tokens, string(r))
			}
		} else if (r >= '0' && r <= '9') || r == '.' {
			current += string(r)
		} else {
			return nil, invalidCharacter
		}
	}
	if current != "" {
		tokens = append(tokens, current)
	}
	return tokens, nil
}

// precedence returns the precedence of the given operator.
func precedence(op rune) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	case '(', ')':
		return 0
	default:
		return -1
	}
}

// ConvertToPostfix converts the infix tokens to postfix notation.
func ConvertToPostfix(tokens []string) ([]string, error) {
	var output []string
	var stack []rune

	for _, token := range tokens {
		if len(token) == 0 {
			continue
		}
		if token == "(" {
			stack = append(stack, '(')
		} else if token == ")" {
			for len(stack) > 0 && stack[len(stack)-1] != '(' {
				output = append(output, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, fmt.Errorf("mismatched parentheses")
			}
			stack = stack[:len(stack)-1]
		} else if isOperator(token) {
			op := rune(token[0])
			for len(stack) > 0 && precedence(op) <= precedence(stack[len(stack)-1]) {
				output = append(output, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, op)
		} else {
			output = append(output, token)
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == '(' {
			return nil, fmt.Errorf("mismatched parentheses")
		}
		output = append(output, string(stack[len(stack)-1]))
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

func isOperator(token string) bool {
	if len(token) != 1 {
		return false
	}
	c := token[0]
	return c == '+' || c == '-' || c == '*' || c == '/'
}

// GenerateTasks processes the postfix tokens to generate tasks.
func GenerateTasks(postfix []string) ([]models.InternalTask, error) {
	var stack []interface{}
	var tasks []models.InternalTask

	for _, token := range postfix {
		if isOperator(token) {
			if len(stack) < 2 {
				return nil, fmt.Errorf("invalid expression")
			}
			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			id := uuid.New().String()

			task := models.InternalTask{
				ID:            id,
				Arg1:          left,
				Arg2:          right,
				Operation:     token,
				OperationTime: 0, // TODO: get time from env
				Result:        "",
			}
			tasks = append(tasks, task)
			stack = append(stack, id)
		} else {
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid number %s: %v", token, err)
			}
			stack = append(stack, num)
		}
	}

	if len(stack) != 1 {
		return nil, fmt.Errorf("invalid expression")
	}

	return tasks, nil
}

// ParseExpression parses the expression into a list of tasks.
func ParseExpression(expression string) ([]models.InternalTask, error) {
	tokens, err := tokenize(expression)
	if err != nil {
		return nil, err
	}
	postfix, err := ConvertToPostfix(tokens)
	if err != nil {
		return nil, err
	}
	tasks, err := GenerateTasks(postfix)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetTasksJSON returns the tasks as a JSON array.
func GetTasksJSON(expression string) (string, error) {
	tasks, err := ParseExpression(expression)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
