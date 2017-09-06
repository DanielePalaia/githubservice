package main

import (
	"net/http"
	"newrelic/request"
	"newrelic/servicelog"
	"time"
)

// Main entry point of the service
func main() {
	logger := servicelog.GetInstance()
	logger.Println(time.Now().UTC(), "Starting service")
	http.HandleFunc("/", request.HandleRequest)
	http.ListenAndServe(":8080", nil)
}
