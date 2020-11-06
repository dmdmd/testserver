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
)

func main() {
	listenAddr := flag.String("http.addr", ":8080", "http listen address")
	flag.Parse()

	http.HandleFunc("/numbers", handleNumbers)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func handleNumbers(w http.ResponseWriter, r *http.Request) {
	param := "u"
	keys, ok := r.URL.Query()[param]
	numbersMap := map[int]bool{}
	result := map[string][]int{
		"numbers": {},
	}

	if !ok || len(keys[0]) < 1 {
		log.Printf("Url Param '%v' is missing\n", param)
		return
	}

	for _, url := range keys {
		numbersMap = mergeInMap(numbersMap, handleURL(url))
	}

	for _, k := range reflect.ValueOf(numbersMap).MapKeys() {
		val := int(k.Int())
		result["numbers"] = append(result["numbers"], val)
	}

	sort.Ints(result["numbers"])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(result)
}

func handleURL(url string) []int {
	page := map[string][]int{
		"numbers": {},
	}
	resp, err := http.Get(url)

	if err != nil {
		log.Printf("Couldn't request %v: %v\n", url, err)
		return []int{}
	}

	requestDump, err := httputil.DumpResponse(resp, true)

	if err != nil {
		log.Printf("Url %v didn't return a valid response: %v\n", url, err)
		return []int{}
	}

	fmt.Println(string(requestDump))

	defer resp.Body.Close()

	//replace with a switch to cover other http response codes
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Printf("Url %v return code error %v\n", url, err)
			return []int{}
		}

		// we unmarshal our byteArray which contains our
		// jsonFile's content
		json.Unmarshal(bodyBytes, &page)

		return page["numbers"]
	}

	log.Printf("Url %v return code %v", url, resp.StatusCode)
	return []int{}

}

func mergeInMap(arr1 map[int]bool, arr2 []int) map[int]bool {
	for _, value := range arr2 {
		arr1[value] = true
	}

	return arr1
}
