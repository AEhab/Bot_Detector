# Bot Detector
## What I did
 I made a HTTP server that accepts any POST request (JSON) from multiple clients' websites. Each request forms part of a struct (for that particular visitor) that will be printed on the terminal in every POST request and when the struct is fully complete.

## Frontend (JS)

I used XMLHttpRequest (XHR) API for sending (JSON) requests in 3 cases:

1. In the first time the user attempts to resize the window.
2. Every time the user tries to copy or paste from any field in the form.
3. When the user submits the form.

See the frontend code here [index.html](client/index.html).

### Example of JSON Requests

```javascript
{
    "eventType": "resizingWindow",
    "websiteUrl": "https://ravelin.com",
    "sessionId":  "jGspdT-9cjpwBmqcfkKnkPn2DJS2",
    "beforeWidth": 1470 ,
    "beforeHeight": 830 ,
    "afterWidth": 830 ,
    "afterHeight": 830
}

{
  "eventType": "copyAndPaste",
  "websiteUrl": "https://ravelin.com",
  "sessionId": "jGspdT-9cjpwBmqcfkKnkPn2DJS2",
  "copiedorpaste": true,
  "formId": "inputCardNumber"
}

{
  "eventType": "timeTaken",
  "websiteUrl": "https://ravelin.com",
  "sessionId": "jGspdT-9cjpwBmqcfkKnkPn2DJS2",
  "time": 72, // seconds
}
```
## Backend (GO) 

At the backend I used GO standard packages to:
1. Serve `index.html` using `http.HandleFunc()` and `template.ParseFiles()` to capture the request and parse the html file.
2. Accept a POST request in JSON format using `json.Unmarshal()` function.
3. Map this JSON objects to our struct `Data` and insert it to a map, called `users`, where the key is the client sessionId and the value is our struct.
4. Print the struct at each stage of construction and at the end when it is completed.
5. Handling concurrent requests, since `http.HandleFunc()` is running in a goroutine functions so it can manage receiving multiple requests at the same time but we can't access out temporary database (our map `users`) so I created and a refrence instance from `mutex = &sync.Mutex{}` to protect the map from being accessed concurrency using `mutex.Lock()` and `mutex.Unlock()`.

See the code here [main.go](server/main.go)

### Testing

For more efficient and effective code and easy testing I created some unit tests and hole function test which can be run easliy through the terminal using `code-test/$ go test server/testname_test.go server/main.go` for separate tests or using `code-test/$ go test`. 

I have 4 different tests:
1. [basicFuncs_test.go](server/basicFuncs_test.go) which tests some basic functions like `inserData()` , `getNewSessionId()` and `floatToStr()`.
2. [handlerFuncs_test.go](server/handlerFuncs_test.go) which tests the functions used to handle different JSON objects like `submittingHandler()`, `copyAndPasteHandler()` and `resizingHandler()`.
3. [httpHandlerFuncs_test.go](server/httpHandlerFuncs_test.go) which tests the fucntions used to handle HTTP POST and GET requests like `eventHandler()` and `indexHandler()`
4. [overall_test.go](server/overall_test.go) which tests the over all functionality of the server in concurrent requests.



   
