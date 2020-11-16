package main

import (
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
	port string
}

type serviceResult struct {
	numbersMap  map[int]bool
	resultArray []int
}

const (
	timeout    time.Duration = time.Duration(500 * time.Millisecond)
	timeoutMsg string        = "{\"numbers\":[]}"
)

func main() {
	var listenAddr *string = flag.String("http.addr", ":8080", "http listen address")
	var logFlag *bool = flag.Bool("v", false, "Verbose")
	flag.Parse()

	if !*logFlag {
		log.SetOutput(ioutil.Discard)
	}

	s := service{
		port: *listenAddr,
	}

	s.listen()
}

/*
listen starts the default http server endpoint
*/
func (s *service) listen() {
	http.Handle("/numbers", http.TimeoutHandler(http.HandlerFunc(s.handler), timeout, timeoutMsg))

	log.Fatal(http.ListenAndServe(s.port, nil))
}

/*
handler extracts other urls from the query parameters
*/
func (s *service) handler(w http.ResponseWriter, r *http.Request) {
	var (
		start    time.Time = time.Now()
		result   []int
		urlParam = "u"
	)

	urlsInQuery, ok := r.URL.Query()[urlParam]

	if !ok || len(urlsInQuery[0]) < 1 {
		log.Printf("Url Param '%v' is missing\n", urlParam)
		s.respondToClient(w, start, result)
		return
	}
	result = s.requestToUrls(urlsInQuery)

	s.respondToClient(w, start, result)
}

/*
requestToUrls requests the urls input and merges the return in a sorted array
*/
func (s *service) requestToUrls(urls []string) []int {
	var (
		arr          []int
		duration     time.Duration
		resultStruct = serviceResult{
			numbersMap:  map[int]bool{},
			resultArray: []int{},
		}
	)

	for _, url := range urls {
		arr, duration = s.handleURL(url)
		log.Printf("URL %v duration: %v and value %v\n", url, duration, arr)
		resultStruct.mergeResults(arr)
	}

	resultStruct.sort()

	return resultStruct.resultArray
}

func (s *service) handleURL(url string) ([]int, time.Duration) {
	var (
		start            time.Time = time.Now()
		responseBodyJSON           = map[string][]int{
			"numbers": {},
		}
		resp        *http.Response
		err         error
		requestDump []byte
		bodyBytes   []byte
	)
	resp, err = http.Get(url)

	if err != nil {
		log.Printf("Couldn't request %v: %v\n", url, err)
		return []int{}, time.Now().Sub(start)
	}

	requestDump, err = httputil.DumpResponse(resp, true)

	if err != nil {
		log.Printf("Url %v didn't return a valid response: %v\n", url, err)
		return []int{}, time.Now().Sub(start)
	}

	fmt.Println(string(requestDump))

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err = ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Printf("Url %v return code error %v\n", url, err)
			return []int{}, time.Now().Sub(start)
		}

		err = json.Unmarshal(bodyBytes, &responseBodyJSON)

		if err != nil {
			log.Printf("Content in %v is not a valid json: %v\n", url, err)
			return []int{}, time.Now().Sub(start)
		}

		return responseBodyJSON["numbers"], time.Now().Sub(start)
	}

	log.Printf("Url %v return code %v", url, resp.StatusCode)
	return []int{}, time.Now().Sub(start)
}

/*
respondToClient Writes content type and status 200 on the header and the formated result
in the body
*/
func (s *service) respondToClient(w http.ResponseWriter, start time.Time, numbers []int) time.Duration {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	dur := time.Now().Sub(start)
	log.Printf("total duration %v\n", dur)
	json.NewEncoder(w).Encode(map[string]interface{}{"numbers": numbers})

	return dur
}

/*
MergeResults Inserts arr elements that dont exist in results. It uses m1 to track
existing elements in "results"
*/
func (sr *serviceResult) mergeResults(arr []int) {
	for _, value := range arr {
		if !sr.numbersMap[value] {
			sr.resultArray = append(sr.resultArray, value)
		}
		sr.numbersMap[value] = true
	}
}

func (sr *serviceResult) sort() {
	sort.Ints(sr.resultArray)
}
