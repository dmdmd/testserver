package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestHandlerNoParam(t *testing.T) {
	var (
		buf    bytes.Buffer
		logger = log.New(&buf, "logger: ", log.Lshortfile)
	)
	s := service{
		log: logger,
	}
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
	var (
		buf    bytes.Buffer
		logger = log.New(&buf, "logger: ", log.Lshortfile)
	)
	s := service{
		log: logger,
	}
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

// func BenchmarkHandler(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		rand.Int()
// 	}
// }

func BenchmarkCalculate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Calculate(2)
	}
}

// func TestRespondToClient(t *testing.T) {
// 	numbers := []int{1, 2}
// 	start := time.Now()
// 	w := httptest.NewRecorder()
// 	result := RespondToClient(w, start, numbers)
// 	if result != 1 {
// 		t.Errorf("Abs(-1) = %d; want 1", got)
// 	}
// }

// func respondToClient(w http.ResponseWriter, start time.Time, numbers []int) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)

// 	log.Printf("total duration %v\n", time.Now().Sub(start))
// 	json.NewEncoder(w).Encode(map[string]interface{}{"numbers": numbers})
// }

func TestMergeResultsEmpty(t *testing.T) {

	var m1 = map[int]bool{}
	var r1 []int
	var a1 = []int{1}

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
