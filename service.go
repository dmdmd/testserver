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

const (
	timeout    = time.Duration(500 * time.Millisecond)
	timeoutMsg = "{\"numbers\":[]}"
)

func main() {
	listenAddr := flag.String("http.addr", ":8080", "http listen address")
	flag.Parse()

	// http.HandleFunc("/numbers", handler)
	http.Handle("/numbers", http.TimeoutHandler(http.HandlerFunc(handler), timeout, timeoutMsg))

	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	// c := make(chan []int)
	result := []int{}
	param := "u"

	// go func() {
	keys, ok := r.URL.Query()[param]

	if !ok || len(keys[0]) < 1 {
		log.Printf("Url Param '%v' is missing\n", param)
		respondToClient(w, start, result)
		return
	}
	result = handleRequest(keys)
	// c <- aux
	// }()

	// Listen on our channel AND a timeout channel - which ever happens first.
	// select {
	// case res := <-c:
	// 	result = res
	// 	log.Println("finished on time")
	// 	respondToClient(w, start, result)

	// case <-time.After(500 * time.Millisecond):
	// 	result = <-c
	// 	log.Println("out of time :(")
	// 	respondToClient(w, start, result)
	// 	return
	// }

	respondToClient(w, start, result)
}

func handleRequest(keys []string) []int {
	numbersMap := map[int]bool{}
	result := []int{}

	for _, url := range keys {
		arr, duration := handleURL(url)
		log.Printf("URL %v duration: %v and value %v\n", url, duration, arr)
		numbersMap, result = mergeResults(numbersMap, result, arr)
	}

	// for _, k := range reflect.ValueOf(numbersMap).MapKeys() {
	// 	val := int(k.Int())
	// 	result = append(result, val)
	// }

	sort.Ints(result)
	return result
}

func handleURL(url string) ([]int, time.Duration) {
	start := time.Now()
	page := map[string][]int{
		"numbers": {},
	}
	resp, err := http.Get(url)

	if err != nil {
		log.Printf("Couldn't request %v: %v\n", url, err)
		return []int{}, time.Now().Sub(start)
	}

	requestDump, err := httputil.DumpResponse(resp, true)

	if err != nil {
		log.Printf("Url %v didn't return a valid response: %v\n", url, err)
		return []int{}, time.Now().Sub(start)
	}

	fmt.Println(string(requestDump))

	defer resp.Body.Close()

	//replace with a switch to cover other http response codes
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Printf("Url %v return code error %v\n", url, err)
			return []int{}, time.Now().Sub(start)
		}

		// we unmarshal our byteArray which contains our
		// jsonFile's content
		json.Unmarshal(bodyBytes, &page)

		return page["numbers"], time.Now().Sub(start)
	}

	log.Printf("Url %v return code %v", url, resp.StatusCode)
	return []int{}, time.Now().Sub(start)
}

func mergeResults(m1 map[int]bool, results []int, arr []int) (map[int]bool, []int) {
	for _, value := range arr {
		if !m1[value] {
			results = append(results, value)
		}
		m1[value] = true
	}

	return m1, results
}

func respondToClient(w http.ResponseWriter, start time.Time, numbers []int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	log.Printf("total duration %v\n", time.Now().Sub(start))
	json.NewEncoder(w).Encode(map[string]interface{}{"numbers": numbers})
}
