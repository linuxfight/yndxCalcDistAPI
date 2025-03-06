package main

import (
	"agent/internal/config"
	"agent/internal/worker"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	apiUrl := os.Getenv("API_URL")
	if apiUrl == "" {
		apiUrl = "http://localhost:9090/internal/task"
	}

	log.Printf("Worker started with URL: %s\n", apiUrl)

	client := &http.Client{
		Timeout: config.RequestTimeout,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: config.MaxWorkers,
			DisableKeepAlives:   false,
		},
	}

	taskCh := make(chan struct{}, config.WorkerPoolBuffer)
	defer close(taskCh)

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-shutdownCh:
				return
			default:
				taskCh <- struct{}{}
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	for i := 0; i < config.MaxWorkers; i++ {
		go worker.Work(taskCh, client, apiUrl)
	}

	<-shutdownCh
	log.Println("Shutting down worker...")
}
