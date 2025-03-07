package main

import (
	"agent/internal/config"
	"agent/internal/logger"
	"agent/internal/worker"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger.New(false, "")
	defer func() {
		err := logger.Log.Sync()
		if err != nil {
			panic(err)
		}
	}()

	c := config.New()

	logger.Log.Infof("Worker started with URL: %s\n", c.ApiUrl)
	logger.Log.Infof("Workers: %d\n", c.ComputingPower)

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
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	for i := 0; i < c.ComputingPower; i++ {
		go worker.Work(taskCh, client, c.ApiUrl)
	}

	<-shutdownCh
	logger.Log.Info("Shutting down agent...")
}
