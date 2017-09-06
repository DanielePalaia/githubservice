package request

import (
	"fmt"
	"net/http"
	"newrelic/engine"
	"newrelic/servicelog"
	"newrelic/utility"
	"strconv"
	"strings"
	"time"
)

// Receive and manage the request (city and numbers of info)
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	logger := servicelog.GetInstance()
	logger.Println(time.Now().UTC(), "New request")

	// Load auth info (username and password)
	utility.LoadAuth()
	queryPath := r.URL.Path[1:]
	// Check for validity
	if !strings.Contains(queryPath, "location") {
		logger.Println(time.Now().UTC(), "Error 401: Bad request")
		fmt.Fprint(w, "Error 401: Bad request")
		return
	}

	query := r.URL.Path[len("/location="):]
	fmt.Printf("query %s\n", query)

	plusIndex := strings.Index(query, "+")
	if plusIndex < 0 {
		logger.Println(time.Now().UTC(), "Error 401: Bad request")
		fmt.Fprint(w, "Error 401: Bad request")
		return
	}

	location := query[0:plusIndex]
	limit := query[plusIndex+1:]

	limitcount, err := strconv.Atoi(limit)
	if err != nil {
		logger.Println(time.Now().UTC(), "Error 500: Internal server error")
		fmt.Fprint(w, "Error 500: Internal server error")
		return
	}
	if limitcount != 50 && limitcount != 100 && limitcount != 150 {
		logger.Println(time.Now().UTC(), "Error 401: Bad request")
		fmt.Fprint(w, "Error 401: Bad request")
		return
	}

	err = service.CallGitApi(w, location, limitcount)
	if err != nil {
		logger.Println(time.Now().UTC(), "Error 500: Internal server error")
		fmt.Fprint(w, "Error 500: Internal server error")
		return
	}
}
