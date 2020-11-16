package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestHandlerNoParam(t *testing.T) {
	s := service{}

	req, err := http.NewRequest("GET", "/numbers", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"numbers":[]}`
	if !strings.HasPrefix(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHandlerSuccess(t *testing.T) {
	s := service{}

	req, err := http.NewRequest("GET", "/numbers?u=http://localhost:8090/primes&u=http://localhost:8090/fibo", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"numbers":[1,2,3,5,7,8,11,13,21]}`
	if !strings.HasPrefix(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestRespondToClient(t *testing.T) {
	s := service{}

	numbers := []int{1, 2}
	start := time.Now()
	rr := httptest.NewRecorder()

	dur := s.respondToClient(rr, start, numbers)

	expected := "{\"numbers\":[1,2]}"

	if dur.Milliseconds() > 500 {
		t.Errorf("Handler took too long to return: %d; want 500 ms", dur)
	}

	if !strings.HasPrefix(rr.Body.String(), expected) {
		t.Errorf("Test response: got \"%v\"; expected %v", rr.Body.String(), expected)
	}
}

func TestMergeResultsNew(t *testing.T) {
	var (
		arr          = []int{1}
		resultStruct = serviceResult{
			numbersMap:  map[int]bool{},
			resultArray: []int{},
		}
	)

	resultStruct.mergeResults(arr)

	expectedMap := map[int]bool{
		1: true,
	}
	expectedArr := []int{1}

	if !reflect.DeepEqual(resultStruct.numbersMap, expectedMap) {
		t.Errorf("Method returned unexpected map: got %v want %v",
			resultStruct.numbersMap, expectedMap)
	}
	if !reflect.DeepEqual(resultStruct.resultArray, expectedArr) {
		t.Errorf("Method returned unexpected array: got %v want %v",
			resultStruct.resultArray, expectedArr)
	}
}

func TestMergeResultsWithResults(t *testing.T) {
	var (
		a1           = []int{2, 5, 3}
		resultStruct = serviceResult{
			numbersMap: map[int]bool{
				1: true,
			},
			resultArray: []int{1},
		}
	)

	resultStruct.mergeResults(a1)

	sort.Ints(resultStruct.resultArray)

	expectedMap := map[int]bool{
		1: true,
		2: true,
		3: true,
		5: true,
	}
	expectedArr := []int{1, 2, 3, 5}

	if !reflect.DeepEqual(resultStruct.numbersMap, expectedMap) {
		t.Errorf("Method returned unexpected map: got %v want %v",
			resultStruct.numbersMap, expectedMap)
	}
	if !reflect.DeepEqual(resultStruct.resultArray, expectedArr) {
		t.Errorf("Method returned unexpected array: got %v want %v",
			resultStruct.resultArray, expectedArr)
	}
}
