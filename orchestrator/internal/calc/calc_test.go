package calc

import (
	"encoding/json"
	"errors"
	"orchestrator/internal/handlers/models"
	"testing"
)

func TestParseExpression(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		wantTasks   int
		wantErr     error
		checkResult func(t *testing.T, tasks []models.InternalTask)
	}{
		{
			name:       "simple addition",
			expression: "2 + 3",
			wantTasks:  1,
			checkResult: func(t *testing.T, tasks []models.InternalTask) {
				assertTask(t, tasks[0], "+", 2.0, 3.0)
			},
		},
		{
			name:       "operator precedence",
			expression: "2 + 3 * 4",
			wantTasks:  2,
			checkResult: func(t *testing.T, tasks []models.InternalTask) {
				assertTask(t, tasks[0], "*", 3.0, 4.0)
				assertTask(t, tasks[1], "+", 2.0, tasks[0].ID)
			},
		},
		{
			name:       "parentheses",
			expression: "(2 + 3) * 4",
			wantTasks:  2,
			checkResult: func(t *testing.T, tasks []models.InternalTask) {
				assertTask(t, tasks[0], "+", 2.0, 3.0)
				assertTask(t, tasks[1], "*", tasks[0].ID, 4.0)
			},
		},
		{
			name:       "unary minus",
			expression: "-5 + 3",
			wantTasks:  2,
			checkResult: func(t *testing.T, tasks []models.InternalTask) {
				assertTask(t, tasks[0], "-", 0.0, 5.0)
				assertTask(t, tasks[1], "+", tasks[0].ID, 3.0)
			},
		},
		{
			name:       "division by zero",
			expression: "5 / 0",
			wantTasks:  0,
			wantErr:    divisionByZeroError,
		},
		{
			name:       "invalid character",
			expression: "2 + a",
			wantTasks:  0,
			wantErr:    errors.New("unsupported node type"),
		},
		{
			name:       "floating point numbers",
			expression: "3.5 * 2",
			wantTasks:  1,
			checkResult: func(t *testing.T, tasks []models.InternalTask) {
				assertTask(t, tasks[0], "*", 3.5, 2.0)
			},
		},
		{
			name:       "complex expression",
			expression: "((3 + 5) * 2 - 4) / 2",
			wantTasks:  4,
			checkResult: func(t *testing.T, tasks []models.InternalTask) {
				assertTask(t, tasks[0], "+", 3.0, 5.0)
				assertTask(t, tasks[1], "*", tasks[0].ID, 2.0)
				assertTask(t, tasks[2], "-", tasks[1].ID, 4.0)
				assertTask(t, tasks[3], "/", tasks[2].ID, 2.0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := ParseExpression(tt.expression)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if !errors.Is(err, tt.wantErr) && !contains(err.Error(), tt.wantErr.Error()) {
					t.Errorf("unexpected error: %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(tasks) != tt.wantTasks {
				t.Errorf("expected %d tasks, got %d", tt.wantTasks, len(tasks))
			}

			if tt.checkResult != nil {
				tt.checkResult(t, tasks)
			}
		})
	}
}

func TestGetTasksJSON(t *testing.T) {
	expr := "2 + 3 * 4"
	jsonStr, err := GetTasksJSON(expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var tasks []models.InternalTask
	if err := json.Unmarshal([]byte(jsonStr), &tasks); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func assertTask(t *testing.T, task models.InternalTask, op string, arg1, arg2 interface{}) {
	t.Helper()

	if task.Operation != op {
		t.Errorf("expected operation %q, got %q", op, task.Operation)
	}

	if !compareArgs(task.Arg1, arg1) {
		t.Errorf("arg1 mismatch: expected %v (%T), got %v (%T)", arg1, arg1, task.Arg1, task.Arg1)
	}

	if !compareArgs(task.Arg2, arg2) {
		t.Errorf("arg2 mismatch: expected %v (%T), got %v (%T)", arg2, arg2, task.Arg2, task.Arg2)
	}
}

func compareArgs(a, b interface{}) bool {
	switch v1 := a.(type) {
	case float64:
		if v2, ok := b.(float64); ok {
			return v1 == v2
		}
	case string:
		if v2, ok := b.(string); ok {
			return v2 == v1
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) >= len(substr) && s[:len(substr)] == substr)
}
