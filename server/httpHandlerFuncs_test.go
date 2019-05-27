package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

//checking indexHandler() function respond and setting the cookies
func TestIndexHandler(t *testing.T) {
	users = make(map[string]Data)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Using Http recoder to recoder the respond
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(indexHandler)
	handler.ServeHTTP(rr, req)

	// checking for respond status code
	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	} else {
		log.Printf("Status respond = %d ", status)
	}
	///Testing for cookie presents
	// Copy the Cookie over to a new Request
	request := &http.Request{Header: http.Header{"Cookie": rr.HeaderMap["Set-Cookie"]}}

	// Extract the dropped cookie from the request.
	cookie, err := request.Cookie("sessionid")
	if err != nil {
		t.Errorf("Expected Cookie named 'test'")
	} else {
		log.Printf("Cookie is found = %+v \n", cookie)
	}

}

//Checking eventHandler() function respond and checking building the structure
func TestEventHandler(t *testing.T) {
	// Example for copy or paste request
	req1 := map[string]interface{}{
		"eventType":     "copyAndPaste",
		"websiteUrl":    "https://ravelin.com",
		"sessionId":     "test1",
		"copiedorpaste": bool(true),
		"formId":        "inputCardNumber"}
	// Example for resizing window request
	req2 := map[string]interface{}{
		"eventType":    "resizingWindow",
		"websiteUrl":   "https://ravelin.com",
		"sessionId":    "test2",
		"beforeWidth":  float64(700),
		"beforeHeight": float64(1600),
		"afterWidth":   float64(750),
		"afterHeight":  float64(1650)}
	// Example for Submitting form request
	req3 := map[string]interface{}{
		"eventType":  "timeTaken",
		"websiteUrl": "https://ravelin.com",
		"sessionId":  "test3",
		"time":       float64(2500)}

	requests := []map[string]interface{}{req1, req2, req3}
	users = make(map[string]Data)

	var data [3]Data
	i := 0
	for i < 3 {
		jsonValue, err := json.Marshal(requests[i])
		if err != nil {
			t.Errorf("Error in json.Marshal()")
		}
		req, err := http.NewRequest("POST", "/event", bytes.NewBuffer(jsonValue))
		if err != nil {
			t.Errorf("Couldn't establish new connection")
		}
		// creating http recorder to examin the respond of the function
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(eventsHandler)
		handler.ServeHTTP(rr, req)
		status := rr.Code
		if status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		} else {
			log.Printf("Status respond = %d ", status)
		}
		data[i] = users["test"+strconv.Itoa(i+1)]
		i++
	}
	// Testing the structure already built
	cond := (data[0].SessionId == "test1") && (data[0].WebsiteUrl == "https://ravelin.com") && (data[0].CopyAndPaste["inputCardNumber"])
	if !cond {
		t.Errorf("Error in data %+v\n", data[0])
	}
	cond = (data[1].SessionId == "test2") && (data[1].WebsiteUrl == "https://ravelin.com")
	cond2 := (data[1].ResizeFrom == Dimension{"700", "1600"}) && (data[1].ResizeTo == Dimension{"750", "1650"})
	if !(cond && cond2) {
		t.Errorf("Error in data %+v\n", data[1])
	}
	cond = (data[2].SessionId == "test3") && (data[2].WebsiteUrl == "https://ravelin.com") && (data[2].FormCompletionTime == 2500)
	if !cond {
		t.Errorf("Error in data %+v\n", data[2])
	}
}
