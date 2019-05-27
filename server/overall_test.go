package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"
)

// wt is used to make the main function wait for all the goroutine to finish
var wt sync.WaitGroup

// makeRequest(url string, sessionid string, i int) isfunction used to make HTTP POST requests to url
func makeRequest(t *testing.T, sessionid string, i int) {
	// Example for copy or paste request
	req1 := map[string]interface{}{
		"eventType":     "copyAndPaste",
		"websiteUrl":    "https://ravelin.com",
		"sessionId":     sessionid,
		"copiedorpaste": bool(true),
		"formId":        "inputCardNumber"}
	// Example for resizing window request
	req2 := map[string]interface{}{
		"eventType":    "resizingWindow",
		"websiteUrl":   "https://ravelin.com",
		"sessionId":    sessionid,
		"beforeWidth":  float64(700),
		"beforeHeight": float64(1600),
		"afterWidth":   float64(750),
		"afterHeight":  float64(1650)}
	// Example for Submitting form request
	req3 := map[string]interface{}{
		"eventType":  "timeTaken",
		"websiteUrl": "https://ravelin.com",
		"sessionId":  sessionid,
		"time":       float64(2500)}

	requests := []map[string]interface{}{req1, req2, req3}
	//Encoding the request ro json format
	jsonValue, err := json.Marshal(requests[int(math.Mod(float64(i), 3))])
	if err != nil {
		t.Errorf("Problem in Marshal " + sessionid)
	}
	//Sending the http POST request to the server
	res, err := http.Post("http://localhost:8080/event", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Errorf("Problem in Requesting " + sessionid)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Bad status code = %d at session id %s", res.StatusCode, sessionid)
	}
	log.Println("done " + sessionid)
	wt.Done()
}

//Testing the hole functionality of the code especially during concurrent requests
func TestHoleserverfunctionality(t *testing.T) {
	// assume we have 20 client at the same time
	clientCount := 20
	// start listen and serve and setting the handlers
	go func() {
		users = make(map[string]Data)
		http.HandleFunc("/event", eventsHandler)
		http.HandleFunc("/", indexHandler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// waiting for the serve to start listening
	time.Sleep(1 * time.Second)

	wt.Add(clientCount)
	i := 0
	for i < clientCount {
		i++
		str := strconv.Itoa(i)
		// go requesting concurrent requests
		go makeRequest(t, "SessionID"+str, i)
	}
	//Wait untill all the goroutines finish
	wt.Wait()
}
