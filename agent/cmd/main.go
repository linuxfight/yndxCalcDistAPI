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
	c := config.New()

	log.Printf("Worker started with URL: %s\n", c.ApiUrl)
	log.Printf("Workers: %d\n", c.ComputingPower)

	client := &http.Client{
		Timeout: config.RequestTimeout,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: c.ComputingPower,
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

	for i := 0; i < c.ComputingPower; i++ {
		go worker.Work(taskCh, client, c.ApiUrl)
	}

	<-shutdownCh
	log.Println("Shutting down worker...")
}
