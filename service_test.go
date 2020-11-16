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

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/numbers", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.Handler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"numbers":[]}`
	if !strings.HasPrefix(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHandlerSuccess(t *testing.T) {
	s := service{}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/numbers?u=http://localhost:8090/primes&u=http://localhost:8090/fibo", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.Handler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
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

	dur := s.RespondToClient(rr, start, numbers)

	expected := "{\"numbers\":[1,2]}"

	if dur.Milliseconds() > 500 {
		t.Errorf("Handler took too long to return: %d; want 500 ms", dur)
	}

	if !strings.HasPrefix(rr.Body.String(), expected) {
		t.Errorf("Test response: got \"%v\"; expected %v", rr.Body.String(), expected)
	}
}

func TestMergeResultsEmpty(t *testing.T) {
	var (
		m1 = map[int]bool{}
		r1 []int
		a1 = []int{1}
	)

	m1, r1 = MergeResults(m1, r1, a1)

	expectedMap := map[int]bool{
		1: true,
	}
	expectedArr := []int{1}

	if !reflect.DeepEqual(m1, expectedMap) {
		t.Errorf("Method returned unexpected map: got %v want %v",
			m1, expectedMap)
	}
	if !reflect.DeepEqual(r1, expectedArr) {
		t.Errorf("Method returned unexpected array: got %v want %v",
			r1, expectedArr)
	}
}

func TestMergeResultsWithResults(t *testing.T) {
	var (
		m1 = map[int]bool{
			1: true,
		}
		r1 = []int{1}
		a1 = []int{2, 5, 3}
	)

	m1, r1 = MergeResults(m1, r1, a1)

	sort.Ints(r1)

	expectedMap := map[int]bool{
		1: true,
		2: true,
		3: true,
		5: true,
	}
	expectedArr := []int{1, 2, 3, 5}

	if !reflect.DeepEqual(m1, expectedMap) {
		t.Errorf("Method returned unexpected map: got %v want %v",
			m1, expectedMap)
	}
	if !reflect.DeepEqual(r1, expectedArr) {
		t.Errorf("Method returned unexpected array: got %v want %v",
			r1, expectedArr)
	}
}
