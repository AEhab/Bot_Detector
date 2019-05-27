package main

// imported packages from Golang standard library
import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Dimension is struct to represent window dimensions
type Dimension struct {
	Width  string
	Height string
}

// Data is our main Struct
type Data struct {
	WebsiteUrl         string
	SessionId          string
	ResizeFrom         Dimension
	ResizeTo           Dimension
	CopyAndPaste       map[string]bool // map[fieldId]true
	FormCompletionTime int             // Seconds
}

// users is our temporary storage memory to store data of different users
var users map[string]Data

// mutex is an instance of sync.Mutex{} used to protect writing in the memory (users) at the same time (concurrent requests)
var mutex = &sync.Mutex{}

//insertData(sessionID string, data Data) is a function used to insert data in the memory (users)
func insertData(sessionID string, data Data) {
	// mutex.Lock() and mutex.Unlock() are used to protect map from writing at the same time
	mutex.Lock()
	users[sessionID] = data
	mutex.Unlock()
}

// getNewSessionID() is a function used to retrieve a new and unique SessionID not used before
func getNewSessionID() string {
	var sessionID string
	// Generate and check until find a new unique ID
	present := true
	for present {
		r := make([]byte, 21) //21 is the sessionID length which will allow us about 3.741444192e50 different user
		io.ReadFull(rand.Reader, r)
		sessionID = base64.URLEncoding.EncodeToString(r)
		_, present = users[sessionID]
	}
	return sessionID
}

// createNewUser(string) is a function used to initialize a user in the map using its sessionid as parameter
func createNewUser(SessionID string) {
	var data Data
	data.CopyAndPaste = make(map[string]bool)
	data.SessionId = SessionID
	insertData(SessionID, data)
}

//floatToStr(num float64) is a function used to floor a float and covert to string
func floatToStr(num float64) string {
	return strconv.Itoa(int(num))
}

// sumbmitinghandler(jsondata map[string]interface{}) is a function used as a handler for submitting the form
func submitingHandler(jsondata map[string]interface{}) {
	SessionID := jsondata["sessionId"].(string)
	var clientdata = users[SessionID]
	clientdata.WebsiteUrl = jsondata["websiteUrl"].(string)
	clientdata.FormCompletionTime = int(jsondata["time"].(float64))
	insertData(SessionID, clientdata)
}

// copyAndPastehandler(jsondata map[string]interface{}) is a function used as a handler for copy or paste in the form
func copyAndPasteHandler(jsondata map[string]interface{}) {
	SessionID := jsondata["sessionId"].(string)
	var clientdata = users[SessionID]
	clientdata.WebsiteUrl = jsondata["websiteUrl"].(string)
	clientdata.CopyAndPaste[jsondata["formId"].(string)] = jsondata["copiedorpaste"].(bool)
	insertData(SessionID, clientdata)
}

// resizinghandler(jsondata map[string]interface{}) is a function used as a handler for resizing the form's window
func resizingHandler(jsondata map[string]interface{}) {
	SessionID := jsondata["sessionId"].(string)
	var clientdata = users[SessionID]
	clientdata.WebsiteUrl = jsondata["websiteUrl"].(string)
	clientdata.ResizeFrom.Height = floatToStr(jsondata["beforeHeight"].(float64))
	clientdata.ResizeFrom.Width = floatToStr(jsondata["beforeWidth"].(float64))
	clientdata.ResizeTo.Height = floatToStr(jsondata["afterHeight"].(float64))
	clientdata.ResizeTo.Width = floatToStr(jsondata["afterWidth"].(float64))
	insertData(SessionID, clientdata)
}

// eventsHandler(w http.ResponseWriter, r *http.Request) is a function used as a handler fot all the event post requests sent by the client
func eventsHandler(w http.ResponseWriter, r *http.Request) {
	// handle only the POST requests
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		var jsonreq interface{}
		err = json.Unmarshal(body, &jsonreq)
		if err != nil {
			panic(err)
		}
		jsondata := jsonreq.(map[string]interface{})
		SessionID := jsondata["sessionId"].(string)
		//checking the existance of the sessionID
		_, present := users[SessionID]
		if !present {
			createNewUser(SessionID)
		}
		//checking the type of the event
		rightType := false
		if jsondata["eventType"] == "copyAndPaste" {
			rightType = true
			copyAndPasteHandler(jsondata)
		} else if jsondata["eventType"] == "resizingWindow" {
			rightType = true
			resizingHandler(jsondata)
		} else if jsondata["eventType"] == "timeTaken" {
			rightType = true
			submitingHandler(jsondata)
		}
		//Printing our struct if the event is a right type
		if rightType == true {
			log.Printf("%+v\n", users[SessionID])
		}
	}
}

//indexHandler(w http.ResponseWriter, r *http.Request) is a function used as a handler for any request (not event) by the client and return index.html
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("../client/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pagecookie, err := r.Cookie("sessionid")
		if err != nil {
			sessionidCookie := http.Cookie{Name: "sessionid", Value: getNewSessionID()}
			http.SetCookie(w, &sessionidCookie)
			createNewUser(sessionidCookie.Value)
			//	log.Printf("%+v\n", users[sessionidCookie.Value]) // printing for debugging
		} else {
			_, present := users[pagecookie.Value]
			if !present {
				createNewUser(pagecookie.Value)
				// log.Printf("%+v\n", users[pagecookie.Value])  // printing for debugging
			}
		}
		err = t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func main() {
	users = make(map[string]Data)
	http.HandleFunc("/event", eventsHandler)
	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
