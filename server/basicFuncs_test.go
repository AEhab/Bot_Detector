package main

import (
	"log"
	"strconv"
	"sync"
	"testing"
)

//Checking for the ability of insertData() function to handel concurrency calls
func TestInsertData(t *testing.T) {

	i := 0
	//sync.WaitGroup is used to make the program wait untill all the goroutines finish
	var wt sync.WaitGroup
	wt.Add(1000)
	users = make(map[string]Data)
	for i < 1000 {
		go func(i int, wt *sync.WaitGroup) {
			var data Data
			insertData("123"+strconv.Itoa(i), data)
			wt.Done()
		}(i, &wt)
		i++
	}
	i = 0
	present := false
	wt.Wait()
	for i < 1000 {
		_, present = users["123"+strconv.Itoa(i)]
		if present == false {
			break
		}
		i++
	}
	if !present {
		t.Errorf("insertData() data not found in the map\n")
	}
}

//Checking the uniqueness of the NewSessionId created
func TestGetNewSessionID(t *testing.T) {
	i := 0
	var testmap map[string]bool
	present := false
	testmap = make(map[string]bool)
	for i < 1000 {
		newStr := getNewSessionID()
		_, present = testmap[newStr]
		if present {
			break
		}
		i++
	}
	if present {
		t.Errorf("getNewSessionID() returned a repeated value\n")
	}
}

func TestFloatToStr(t *testing.T) {
	x := 3.265485
	y := 5.956694
	ans1 := floatToStr(x)
	ans2 := floatToStr(y)
	if ans1 != "3" || ans2 != "5" {
		t.Errorf("Error in floatToStr ans1=%s ans2=%s\n", ans1, ans2)
	} else {
		log.Printf("convertion is right ans1 = %s , ans2 = %s\n", ans1, ans2)
	}
}
