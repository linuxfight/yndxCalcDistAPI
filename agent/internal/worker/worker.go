package worker

import (
	"agent/internal/logger"
	"agent/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Work is a main worker method
func Work(taskCh <-chan struct{}, client *http.Client, apiUrl string) {
	for range taskCh {
		processTask(client, apiUrl)
	}
}

// processTask is a method for processing task in a worker
func processTask(client *http.Client, apiUrl string) {
	task, err := getTask(client, apiUrl)
	if err != nil {
		logger.Log.Infof("Error getting task: %v\n", err)
		return
	}
	if task == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(task.OperationTime)*time.Millisecond)
	defer cancel()

	resultChan := make(chan float64, 1)
	errorChan := make(chan error, 1)

	go func() {
		res, err := calculateResult(task)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- res
	}()

	select {
	case <-ctx.Done():
		logger.Log.Infof("Task %s timed out", task.ID)
		if err := sendResult(client, apiUrl, task.ID, models.ERROR); err != nil {
			logger.Log.Infof("Error sending timeout result: %v\n", err)
		}
	case err := <-errorChan:
		logger.Log.Infof("Error calculating result: %v\n", err)
		if err := sendResult(client, apiUrl, task.ID, models.ERROR); err != nil {
			logger.Log.Infof("Error sending error result: %v\n", err)
		}
	case result := <-resultChan:
		if err := sendResult(client, apiUrl, task.ID, result); err != nil {
			logger.Log.Infof("Error sending result: %v\n", err)
			_ = sendResult(client, apiUrl, task.ID, models.ERROR)
		}
	}
}

// getTask is a method for getting a task from the API
func getTask(client *http.Client, apiUrl string) (*models.TaskResponse, error) {
	resp, err := client.Get(apiUrl)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer safeClose(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}

	var task models.TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}

	return &task, nil
}

// calculateResult is a method for calculating result of task
func calculateResult(task *models.TaskResponse) (float64, error) {
	switch task.Operation {
	case "+":
		return task.Arg1 + task.Arg2, nil
	case "-":
		return task.Arg1 - task.Arg2, nil
	case "*":
		return task.Arg1 * task.Arg2, nil
	case "/":
		if task.Arg2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return task.Arg1 / task.Arg2, nil
	default:
		return 0, fmt.Errorf("unknown operation: %s", task.Operation)
	}
}

// sendResult is a method for sending calculation result to the API
func sendResult(client *http.Client, apiUrl, taskID string, result interface{}) error {
	data := models.TaskRequest{
		ID:     taskID,
		Result: result,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer safeClose(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// safeClose is a method for safely closing io
func safeClose(closer io.Closer) {
	if err := closer.Close(); err != nil {
		logger.Log.Infof("Error closing resource: %v\n", err)
	}
}
