package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

const (
	RequestTimeout   = 5 * time.Second
	WorkerPoolBuffer = 20
)

type Config struct {
	ApiUrl         string
	ComputingPower int
}

func New() *Config {
	apiUrl := os.Getenv("API_URL")
	if apiUrl == "" {
		apiUrl = "http://localhost:9090/internal/task"
	}

	powerStr := os.Getenv("POWER")
	power := 5
	if powerStr != "" {
		var err error
		if power, err = strconv.Atoi(powerStr); err != nil {
			log.Printf("Error converting POWER to int: %s\n", err)
		}
	}

	return &Config{
		ApiUrl:         apiUrl,
		ComputingPower: power,
	}
}
