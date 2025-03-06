package config

import "time"

const (
	MaxWorkers       = 5
	RequestTimeout   = 5 * time.Second
	WorkerPoolBuffer = 20
)
