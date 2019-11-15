package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	rng "github.com/leesper/go_rng" // go get -u -v github.com/leesper/go_rng
)

var waitG sync.WaitGroup

type req struct {
	Latency float64 `json:"latency"`
	Start   int64   `json:"start"`
	// StatusCode int32 `json:"status_code"`
}

type graph struct {
	Success       int64 `json:"success"`
	Failed        int64 `json:"failed"`
	TotalRequests int64 `json:"total_requests"`
	Progress      int64 `json:"progress"`
	// StatusCode int32 `json:"status_code"`
}

var startingTime int64

var requestsArr []req

var rps = 22            // Num of requests per second
var duration = 1        // Num of seconds of the test
var ticksPerSecond = 10 // how many ticks per second

var totalRequests int64
var successRequests int64
var failedRequests int64
var progress int64
var errorsNum int64

func main() {
	// start := time.Now()
	// loadTest()
	// fmt.Println(time.Since(start))

	http.HandleFunc("/load", loadTest)
	http.HandleFunc("/graph", showRequestsGraph)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", http.StripPrefix("/", fs))
	http.ListenAndServe(":8080", nil)
}

func showRequestsGraph(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, string(totalRequests))
	g := graph{successRequests, failedRequests, totalRequests, progress}
	x, _ := json.Marshal(g)
	fmt.Fprintf(w, string(x))
}

func loadTest(w http.ResponseWriter, r *http.Request) {
	// func loadTest() {

	requestsArr = requestsArr[:0]
	totalRequests = 0
	progress = 0
	errorsNum = 0
	successRequests = 0
	failedRequests = 0
	var i int64
	var tickCounter = 0
	var lambda = float64(rps) / float64(ticksPerSecond)
	ticker := time.NewTicker(time.Duration(1000/ticksPerSecond) * time.Millisecond)
	startingTime = time.Now().UnixNano() / int64(time.Millisecond)
	for _ = range ticker.C {
		prange := rng.NewPoissonGenerator(time.Now().UnixNano())
		requestsNum := prange.Poisson(lambda)       // Number of requests every tick
		totalRequests = totalRequests + requestsNum // current Total Num of requests
		for i = 0; i < requestsNum; i++ {
			// fmt.Println(totalRequests)
			waitG.Add(1)
			go mkReq()
		}

		tickCounter++
		fmt.Println(requestsNum, totalRequests)
		// if tickCounter == (duration * (ticksPerSecond)) {
		// 	break
		// }
		if (time.Now().UnixNano()/int64(time.Millisecond) - startingTime) > int64(duration*1000) {
			break
		}
	}
	waitG.Wait()
	ticker.Stop()
	x, _ := json.Marshal(requestsArr)

	fmt.Fprint(w, string(x))
	fmt.Println(" ")
	// fmt.Printf("\n%.1f%% Errors.", float32(errorsNum)/float32(totalRequests)*100)
	// fmt.Printf("\n%.1f%% Success.", (float32(totalRequests)-float32(errorsNum))/float32(totalRequests)*100)
}

func mkReq() string {
	start := time.Now()
	request, _ := http.NewRequest("GET", "http://localhost:8081/sleep", nil)
	client := &http.Client{}
	resp, _ := client.Do(request)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	elapsed := time.Since(start)
	elapsedFloat, _ := strconv.ParseFloat(strings.TrimSuffix(elapsed.String(), "s"), 64)
	if resp.StatusCode != 200 {
		errorsNum++
	} // else {
	// 	requestsArr = append(requestsArr, req{elapsedFloat, (start.UnixNano()/int64(time.Millisecond) - startingTime)})
	// }

	if resp.StatusCode != 200 {
		failedRequests++
	} else {
		successRequests++
	}
	requestsArr = append(requestsArr, req{elapsedFloat, (start.UnixNano()/int64(time.Millisecond) - startingTime)})

	// }
	progress++
	// fmt.Println("")
	waitG.Done()
	// fmt.Println(progress, totalRequests)

	// fmt.Printf("\r%.1f%% Completed.", float32(progress)/float32(totalRequests)*100)
	return string(body)
}
