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

	"github.com/gin-gonic/gin"
)

type service struct {
}

const (
	timeout    = time.Duration(500 * time.Millisecond)
	timeoutMsg = "{\"numbers\":[]}"
)

func main() {
	listenAddr := flag.String("http.addr", ":8080", "http listen address")
	logFlag := flag.Bool("v", false, "Verbose")
	flag.Parse()

	if !*logFlag {
		log.SetOutput(ioutil.Discard)
	}

	s := service{}

	s.Listen(listenAddr)
}

/*
Listen starts the default http server endpoint
*/
func (s *service) Listen(listenAddr *string) {

	router := gin.New()
	router.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/pathsNotToLog/"),
		gin.Recovery(),
	)

	router.GET("/numbers", func(c *gin.Context) {
		s.Handler(c.Writer, c.Request)
	})

	log.Fatal(http.ListenAndServe(*listenAddr, router))
}

/*
Handler processes requests to our service
*/
func (s *service) Handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	result := []int{}
	param := "u"

	keys, ok := r.URL.Query()[param]

	if !ok || len(keys[0]) < 1 {
		log.Printf("Url Param '%v' is missing\n", param)
		s.RespondToClient(w, start, result)
		return
	}
	result = s.RequestToUrls(keys)

	s.RespondToClient(w, start, result)
}

/*
RequestToUrls requests the urls input and merges the return in a sorted array
*/
func (s *service) RequestToUrls(urls []string) []int {
	numbersMap := map[int]bool{}
	result := []int{}

	for _, url := range urls {
		arr, duration := s.handleURL(url)
		log.Printf("URL %v duration: %v and value %v\n", url, duration, arr)
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

/*
MergeResults Inserts arr elements that dont exist in results. It uses m1 to track
existing elements in "results"
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
RespondToClient Writes content type and status 200 on the header and the formated result
in the body
*/
func (s *service) RespondToClient(w http.ResponseWriter, start time.Time, numbers []int) time.Duration {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	dur := time.Now().Sub(start)
	log.Printf("total duration %v\n", dur)
	json.NewEncoder(w).Encode(map[string]interface{}{"numbers": numbers})

	return dur
}
