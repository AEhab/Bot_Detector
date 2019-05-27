package main

import (
	"testing"
)

//checking on inserting the JSON request data in the right place for the 3 events
func TestSubmitingHandler(t *testing.T) {
	submitReq := map[string]interface{}{
		"eventType":  "timeTaken",
		"websiteUrl": "https://ravelin.com",
		"sessionId":  "ihugyfvygbhunj",
		"time":       float64(2500)}

	users = make(map[string]Data)
	createNewUser("ihugyfvygbhunj")
	submitingHandler(submitReq)
	data := users["ihugyfvygbhunj"]
	cond := (data.SessionId == "ihugyfvygbhunj") && (data.FormCompletionTime == 2500) && (data.WebsiteUrl == "https://ravelin.com")
	if !cond {
		t.Errorf("Error in data %+v\n", data)
	}
}

func TestCopyAndPasteHandler(t *testing.T) {
	copyAndPasteReq := map[string]interface{}{
		"eventType":     "copyorpaste",
		"websiteUrl":    "https://ravelin.com",
		"sessionId":     "ihugyfvygbhunj",
		"copiedorpaste": bool(true),
		"formId":        "inputCardNumber"}

	users = make(map[string]Data)
	createNewUser("ihugyfvygbhunj")
	copyAndPasteHandler(copyAndPasteReq)
	data := users["ihugyfvygbhunj"]
	cond := (data.SessionId == "ihugyfvygbhunj") && (data.WebsiteUrl == "https://ravelin.com") && (data.CopyAndPaste["inputCardNumber"])
	if !cond {
		t.Errorf("Error in data %+v\n", data)
	}
}

func TestResizingHandler(t *testing.T) {
	resizingReq := map[string]interface{}{
		"eventType":    "resizingWindow",
		"websiteUrl":   "https://ravelin.com",
		"sessionId":    "ihugyfvygbhunj",
		"beforeWidth":  float64(700),
		"beforeHeight": float64(1600),
		"afterWidth":   float64(750),
		"afterHeight":  float64(1650)}
	users = make(map[string]Data)
	createNewUser("ihugyfvygbhunj")
	resizingHandler(resizingReq)
	data := users["ihugyfvygbhunj"]
	cond := (data.SessionId == "ihugyfvygbhunj") && (data.WebsiteUrl == "https://ravelin.com")
	cond2 := (data.ResizeFrom == Dimension{"700", "1600"}) && (data.ResizeTo == Dimension{"750", "1650"})
	if !(cond && cond2) {
		t.Errorf("Error in data %+v\n", data)
	}
}
