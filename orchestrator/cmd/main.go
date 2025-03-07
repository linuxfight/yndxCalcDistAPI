package main

import (
	_ "orchestrator/docs"
	"orchestrator/internal/handlers"
)

// @title           Orchestrator API
// @version         1.0
// @description 	API documentation for the Calc Orchestrator
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:9090
func main() {
	// create api controller
	c := handlers.New()

	// start server
	c.Start()
}
