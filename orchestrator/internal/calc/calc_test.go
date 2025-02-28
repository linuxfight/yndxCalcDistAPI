package calc

import (
	"encoding/json"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  []string
		expectErr bool
	}{
		{
			name:      "simple addition",
			input:     "3+4",
			expected:  []string{"3", "+", "4"},
			expectErr: false,
		},
		{
			name:      "negative number",
			input:     "-3 + 4",
			expected:  []string{"-3", "+", "4"},
			expectErr: false,
		},
		{
			name:      "parentheses and multiplication",
			input:     "(2 + 3)*4",
			expected:  []string{"(", "2", "+", "3", ")", "*", "4"},
			expectErr: false,
		},
		{
			name:      "invalid character",
			input:     "3$4",
			expectErr: true,
		},
		{
			name:      "decimal number",
			input:     "2.5 * 3.8",
			expected:  []string{"2.5", "*", "3.8"},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tokenize(tt.input)
			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if !tt.expectErr && !reflect.DeepEqual(tokens, tt.expected) {
				t.Errorf("expected tokens %v, got %v", tt.expected, tokens)
			}
		})
	}
}

func TestConvertToPostfix(t *testing.T) {
	tests := []struct {
		name        string
		inputTokens []string
		expected    []string
		expectErr   bool
	}{
		{
			name:        "simple precedence",
			inputTokens: []string{"3", "+", "4", "*", "2"},
			expected:    []string{"3", "4", "2", "*", "+"},
			expectErr:   false,
		},
		{
			name:        "parentheses",
			inputTokens: []string{"(", "3", "+", "4", ")", "*", "2"},
			expected:    []string{"3", "4", "+", "2", "*"},
			expectErr:   false,
		},
		{
			name:        "mismatched parentheses",
			inputTokens: []string{"(", "3", "+", "4"},
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postfix, err := ConvertToPostfix(tt.inputTokens)
			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if !tt.expectErr && !reflect.DeepEqual(postfix, tt.expected) {
				t.Errorf("expected postfix %v, got %v", tt.expected, postfix)
			}
		})
	}
}

func TestGenerateTasks(t *testing.T) {
	tests := []struct {
		name      string
		postfix   []string
		expected  []Task
		expectErr bool
	}{
		{
			name:    "single addition",
			postfix: []string{"3", "4", "+"},
			expected: []Task{
				{ID: "00000000-0000-0000-0000-000000000000:1", Arg1: 3.0, Arg2: 4.0, Operation: "+", OperationTime: 0},
			},
			expectErr: false,
		},
		{
			name:    "multiple operations",
			postfix: []string{"3", "4", "2", "*", "+"},
			expected: []Task{
				{ID: "00000000-0000-0000-0000-000000000000:1", Arg1: 4.0, Arg2: 2.0, Operation: "*", OperationTime: 0},
				{ID: "00000000-0000-0000-0000-000000000000:2", Arg1: 3.0, Arg2: "00000000-0000-0000-0000-000000000000:1", Operation: "+", OperationTime: 0},
			},
			expectErr: false,
		},
		{
			name:      "invalid postfix (not enough operands)",
			postfix:   []string{"3", "+"},
			expectErr: true,
		},
		{
			name:      "invalid number",
			postfix:   []string{"3.2.1", "+", "4"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := GenerateTasks(tt.postfix, uuid.Nil)
			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if !tt.expectErr {
				if len(tasks) != len(tt.expected) {
					t.Fatalf("expected %d tasks, got %d", len(tt.expected), len(tasks))
				}
				for i, task := range tasks {
					expectedTask := tt.expected[i]
					if task.ID != expectedTask.ID || task.Operation != expectedTask.Operation || task.OperationTime != expectedTask.OperationTime {
						t.Errorf("task %d mismatch: expected %+v, got %+v", i, expectedTask, task)
					}
					if !reflect.DeepEqual(task.Arg1, expectedTask.Arg1) {
						t.Errorf("task %d Arg1 mismatch: expected %v (%T), got %v (%T)", i, expectedTask.Arg1, expectedTask.Arg1, task.Arg1, task.Arg1)
					}
					if !reflect.DeepEqual(task.Arg2, expectedTask.Arg2) {
						t.Errorf("task %d Arg2 mismatch: expected %v (%T), got %v (%T)", i, expectedTask.Arg2, expectedTask.Arg2, task.Arg2, task.Arg2)
					}
				}
			}
		})
	}
}

func TestParseExpression(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  []Task
		expectErr bool
	}{
		{
			name:  "simple expression",
			input: "3 + 4 * 2",
			expected: []Task{
				{ID: "00000000-0000-0000-0000-000000000000:1", Arg1: 4.0, Arg2: 2.0, Operation: "*", OperationTime: 0},
				{ID: "00000000-0000-0000-0000-000000000000:2", Arg1: 3.0, Arg2: "00000000-0000-0000-0000-000000000000:1", Operation: "+", OperationTime: 0},
			},
			expectErr: false,
		},
		{
			name:  "expression with parentheses",
			input: "(2 + 3) * 4",
			expected: []Task{
				{ID: "00000000-0000-0000-0000-000000000000:1", Arg1: 2.0, Arg2: 3.0, Operation: "+", OperationTime: 0},
				{ID: "00000000-0000-0000-0000-000000000000:2", Arg1: "00000000-0000-0000-0000-000000000000:1", Arg2: 4.0, Operation: "*", OperationTime: 0},
			},
			expectErr: false,
		},
		{
			name:      "invalid expression",
			input:     "3 + * 4",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := ParseExpression(tt.input)
			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if !tt.expectErr {
				if len(tasks) != len(tt.expected) {
					t.Fatalf("expected %d tasks, got %d", len(tt.expected), len(tasks))
				}
				for i, task := range tasks {
					expectedTask := tt.expected[i]
					if task.ID != expectedTask.ID || task.Operation != expectedTask.Operation || task.OperationTime != expectedTask.OperationTime {
						t.Errorf("task %d mismatch: expected %+v, got %+v", i, expectedTask, task)
					}
					if !reflect.DeepEqual(task.Arg1, expectedTask.Arg1) {
						t.Errorf("task %d Arg1 mismatch: expected %v (%T), got %v (%T)", i, expectedTask.Arg1, expectedTask.Arg1, task.Arg1, task.Arg1)
					}
					if !reflect.DeepEqual(task.Arg2, expectedTask.Arg2) {
						t.Errorf("task %d Arg2 mismatch: expected %v (%T), got %v (%T)", i, expectedTask.Arg2, expectedTask.Arg2, task.Arg2, task.Arg2)
					}
				}
			}
		})
	}
}

func TestGetTasksJSON(t *testing.T) {
	input := "(2 + 3) * 4"
	expected := []Task{
		{ID: "00000000-0000-0000-0000-000000000000:1", Arg1: 2.0, Arg2: 3.0, Operation: "+", OperationTime: 0},
		{ID: "00000000-0000-0000-0000-000000000000:2", Arg1: "00000000-0000-0000-0000-000000000000:1", Arg2: 4.0, Operation: "*", OperationTime: 0},
	}

	jsonStr, err := GetTasksJSON(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var tasks []Task
	if err := json.Unmarshal([]byte(jsonStr), &tasks); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if len(tasks) != len(expected) {
		t.Fatalf("expected %d tasks, got %d", len(expected), len(tasks))
	}

	for i, task := range tasks {
		expTask := expected[i]
		if task.ID != expTask.ID || task.Operation != expTask.Operation || task.OperationTime != expTask.OperationTime {
			t.Errorf("task %d mismatch: expected %+v, got %+v", i, expTask, task)
		}
		if !reflect.DeepEqual(task.Arg1, expTask.Arg1) {
			t.Errorf("task %d Arg1 mismatch: expected %v (%T), got %v (%T)", i, expTask.Arg1, expTask.Arg1, task.Arg1, task.Arg1)
		}
		if !reflect.DeepEqual(task.Arg2, expTask.Arg2) {
			t.Errorf("task %d Arg2 mismatch: expected %v (%T), got %v (%T)", i, expTask.Arg2, expTask.Arg2, task.Arg2, task.Arg2)
		}
	}
}

func TestSingleNumber(t *testing.T) {
	input := "42"

	tasks, err := ParseExpression(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}
}

func TestEmptyInput(t *testing.T) {
	val, err := ParseExpression("")
	if err == nil {
		t.Errorf("expected error, got %v", val)
	}
}

func TestPrecedence(t *testing.T) {
	tests := []struct {
		op       rune
		expected int
	}{
		{'+', 1},
		{'-', 1},
		{'*', 2},
		{'/', 2},
		{'(', 0},
		{')', 0},
		{'^', -1},
	}

	for _, test := range tests {
		got := precedence(test.op)
		if got != test.expected {
			t.Errorf("precedence(%c) = %d; want %d", test.op, got, test.expected)
		}
	}
}
