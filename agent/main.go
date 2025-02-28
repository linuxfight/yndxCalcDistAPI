package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	apiUrl := os.Getenv("API_URL")

	if apiUrl == "" {
		apiUrl = "http://localhost:9090/internal/task"
	}

	println("worker started with url: " + apiUrl)

	for {
		time.Sleep(time.Second * 10)
		go func() {
			settings := &http.Transport{
				DisableKeepAlives: true,
			}
			client := &http.Client{Transport: settings}

			get, err := client.Get(apiUrl)
			if err != nil {
				return
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					return
				}
			}(get.Body)
			if get.StatusCode == 200 {
				var task taskResponse
				body, err := io.ReadAll(get.Body)
				if err != nil {
					return
				}
				err = json.Unmarshal(body, &task)
				if err != nil {
					return
				}

				var result float64
				switch task.Operation {
				case "+":
					result = task.Arg1 + task.Arg2
				case "-":
					result = task.Arg2 - task.Arg1
				case "*":
					result = task.Arg1 * task.Arg2
				case "/":
					result = task.Arg1 / task.Arg2
				}

				data := taskRequest{
					ID:     task.ID,
					Result: strconv.FormatFloat(result, 'f', -1, 64),
				}

				marshal, err := json.Marshal(&data)
				if err != nil {
					return
				}

				post, err := client.Post(apiUrl, "application/json", bytes.NewBuffer(marshal))
				if err != nil {
					return
				}
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
						return
					}
				}(post.Body)
			}
		}()
	}
}

type taskResponse struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

type taskRequest struct {
	ID     string `json:"id"`
	Result string `json:"result"`
}

func runWithTimeout(timeout time.Duration, f func() int) (int, error) {
	resultCh := make(chan int, 1)

	// Run the function in a separate goroutine.
	go func() {
		resultCh <- f()
	}()

	// Wait for either the result or a timeout.
	select {
	case result := <-resultCh:
		return result, nil
	case <-time.After(timeout):
		return 0, fmt.Errorf("function execution exceeded timeout of %v", timeout)
	}
}
