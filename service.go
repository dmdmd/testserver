package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"reflect"
	"sort"
	"time"
)

func main() {
	listenAddr := flag.String("http.addr", ":8080", "http listen address")
	flag.Parse()

	http.HandleFunc("/numbers", handler)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	// c := make(chan int)
	param := "u"
	keys, ok := r.URL.Query()[param]

	result := map[string][]int{
		"numbers": {},
	}

	if !ok || len(keys[0]) < 1 {
		log.Printf("Url Param '%v' is missing\n", param)
		respondToClient(w, start, result)
		return
	}

	result["numbers"] = handleRequest(keys)

	respondToClient(w, start, result)
}

func respondToClient(w http.ResponseWriter, start time.Time, result map[string][]int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	log.Printf("total duration %v\n", time.Now().Sub(start))
	json.NewEncoder(w).Encode(result)
}

func handleRequest(keys []string) []int {
	numbersMap := map[int]bool{}
	result := []int{}

	for _, url := range keys {
		arr, duration := handleURL(url)
		log.Printf("URL %v duration: %v and value %v\n", url, duration, arr)
		numbersMap = mergeInMap(numbersMap, arr)
	}

	for _, k := range reflect.ValueOf(numbersMap).MapKeys() {
		val := int(k.Int())
		result = append(result, val)
	}

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

func mergeInMap(arr1 map[int]bool, arr2 []int) map[int]bool {
	for _, value := range arr2 {
		arr1[value] = true
	}

	return arr1
}
