package calc

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"

	"github.com/google/uuid"
	"orchestrator/internal/handlers/models"
)

var (
	divisionByZeroError  = fmt.Errorf("division by zero")
	unsupportedNodeError = fmt.Errorf("unsupported node type")
)

// ParseExpression parses a mathematical expression into a sequence of tasks
func ParseExpression(expression string) ([]models.InternalTask, error) {
	exprAst, err := parser.ParseExpr(expression)
	if err != nil {
		return nil, fmt.Errorf("parsing error: %w", err)
	}

	var tasks []models.InternalTask
	_, err = processNode(exprAst, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// processNode recursively processes AST nodes and creates tasks
func processNode(node ast.Node, tasks *[]models.InternalTask) (interface{}, error) {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		return processBinaryExpr(n, tasks)
	case *ast.UnaryExpr:
		return processUnaryExpr(n, tasks)
	case *ast.BasicLit:
		return processBasicLit(n)
	case *ast.ParenExpr:
		return processNode(n.X, tasks)
	default:
		return nil, unsupportedNodeError
	}
}

func processBinaryExpr(expr *ast.BinaryExpr, tasks *[]models.InternalTask) (interface{}, error) {
	left, err := processNode(expr.X, tasks)
	if err != nil {
		return nil, err
	}

	right, err := processNode(expr.Y, tasks)
	if err != nil {
		return nil, err
	}

	// Check for division by zero with literal values
	if expr.Op == token.QUO {
		if rval, ok := right.(float64); ok && rval == 0 {
			return nil, divisionByZeroError
		}
	}

	return createTask(tasks, left, right, expr.Op.String())
}

func processUnaryExpr(expr *ast.UnaryExpr, tasks *[]models.InternalTask) (interface{}, error) {
	if expr.Op != token.SUB {
		return nil, fmt.Errorf("unsupported unary operator: %v", expr.Op)
	}

	operand, err := processNode(expr.X, tasks)
	if err != nil {
		return nil, err
	}

	return createTask(tasks, 0.0, operand, token.SUB.String())
}

func processBasicLit(lit *ast.BasicLit) (float64, error) {
	switch lit.Kind {
	case token.INT, token.FLOAT:
		value, err := strconv.ParseFloat(lit.Value, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number: %w", err)
		}
		return value, nil
	default:
		return 0, fmt.Errorf("unsupported literal type: %v", lit.Kind)
	}
}

func createTask(tasks *[]models.InternalTask, left, right interface{}, operation string) (string, error) {
	taskID := uuid.New().String()
	*tasks = append(*tasks, models.InternalTask{
		ID:            taskID,
		Arg1:          left,
		Arg2:          right,
		Operation:     operation,
		OperationTime: 0,
		Result:        "",
	})
	return taskID, nil
}

// GetTasksJSON returns tasks as JSON string
func GetTasksJSON(expression string) (string, error) {
	tasks, err := ParseExpression(expression)
	if err != nil {
		return "", fmt.Errorf("error parsing expression: %w", err)
	}

	jsonData, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		return "", fmt.Errorf("error marshaling tasks: %w", err)
	}

	return string(jsonData), nil
}
