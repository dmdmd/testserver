package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"sort"
	"time"
)

type service struct {
	log *log.Logger
}

const (
	timeout    = time.Duration(500 * time.Millisecond)
	timeoutMsg = "{\"numbers\":[]}"
)

func main() {
	listenAddr := flag.String("http.addr", ":8080", "http listen address")
	flag.Parse()

	var (
		buf    bytes.Buffer
		logger = log.New(&buf, "logger: ", log.Lshortfile)
	)

	s := service{
		log: logger,
	}

	s.Handle(listenAddr)

}

/*

 */
func (s *service) Handle(listenAddr *string) {
	http.Handle("/numbers", http.TimeoutHandler(http.HandlerFunc(s.Handler), timeout, timeoutMsg))

	s.log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

/*
Handler ...
*/
func (s *service) Handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	result := []int{}
	param := "u"

	keys, ok := r.URL.Query()[param]

	if !ok || len(keys[0]) < 1 {
		s.log.Printf("Url Param '%v' is missing\n", param)
		s.RespondToClient(w, start, result)
		return
	}
	result = s.RequestToUrls(keys)

	s.RespondToClient(w, start, result)
}

/*
RequestToUrls ...
*/
func (s *service) RequestToUrls(keys []string) []int {
	numbersMap := map[int]bool{}
	result := []int{}

	for _, url := range keys {
		arr, duration := s.handleURL(url)
		s.log.Printf("URL %v duration: %v and value %v\n", url, duration, arr)
		numbersMap, result = MergeResults(numbersMap, result, arr)
	}

	sort.Ints(result)
	return result
}

func (s *service) handleURL(url string) ([]int, time.Duration) {
	start := time.Now()
	page := map[string][]int{
		"numbers": {},
	}
	resp, err := http.Get(url)

	if err != nil {
		s.log.Printf("Couldn't request %v: %v\n", url, err)
		return []int{}, time.Now().Sub(start)
	}

	requestDump, err := httputil.DumpResponse(resp, true)

	if err != nil {
		s.log.Printf("Url %v didn't return a valid response: %v\n", url, err)
		return []int{}, time.Now().Sub(start)
	}

	fmt.Println(string(requestDump))

	defer resp.Body.Close()

	//replace with a switch to cover other http response codes
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			s.log.Printf("Url %v return code error %v\n", url, err)
			return []int{}, time.Now().Sub(start)
		}

		// we unmarshal our byteArray which contains our
		// jsonFile's content
		json.Unmarshal(bodyBytes, &page)

		return page["numbers"], time.Now().Sub(start)
	}

	s.log.Printf("Url %v return code %v", url, resp.StatusCode)
	return []int{}, time.Now().Sub(start)
}

/*
MergeResults ...
*/
func MergeResults(m1 map[int]bool, results []int, arr []int) (map[int]bool, []int) {
	for _, value := range arr {
		if !m1[value] {
			results = append(results, value)
		}
		m1[value] = true
	}

	return m1, results
}

/*
RespondToClient ...
*/
func (s *service) RespondToClient(w http.ResponseWriter, start time.Time, numbers []int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	s.log.Printf("total duration %v\n", time.Now().Sub(start))
	json.NewEncoder(w).Encode(map[string]interface{}{"numbers": numbers})
}

// Calculate returns x + 2.
func Calculate(x int) (result int) {
	result = x + 2
	return result
}
