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
)

var waitG sync.WaitGroup

type req struct {
	Latency float64 `json:"latency"`
	Start   int64   `json:"start"`
	// StatusCode int32 `json:"status_code"`
}

var startingTime int64

var requestsArr []req

var rps = 25      // Num of requests per second
var duration = 15 // Num of seconds of the test

var progress int32
var errorsNum int32

func main() {

	http.HandleFunc("/load", loadTest)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", http.StripPrefix("/", fs))

	// http.Handle("/", http.StripPrefix(strings.TrimRight(path, "/"), http.FileServer(http.Dir(directory))))

	// http.HandleFunc("/", serveTemplate)
	http.ListenAndServe(":8080", nil)
}

func loadTest(w http.ResponseWriter, r *http.Request) {

	requestsArr = requestsArr[:0]
	progress = 0
	errorsNum = 0

	startingTime = time.Now().UnixNano() / int64(time.Millisecond)

	for i := 0; i < duration; i++ {
		// @ToDo : Handling if rps is > 1000
		period := 1000 / rps

		for j := 0; j < rps; j++ {
			waitG.Add(1)
			go mkReq()
			time.Sleep(time.Duration(period) * time.Millisecond)
		}
	}

	waitG.Wait()

	x, _ := json.Marshal(requestsArr)

	fmt.Fprint(w, string(x))
	fmt.Printf("\n%.1f%% Errors.", float32(errorsNum)/float32(rps*duration)*100)
	fmt.Printf("\n%.1f%% Success.", (float32(rps*duration)-float32(errorsNum))/float32(rps*duration)*100)
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

	// if resp.StatusCode != 200 {
	requestsArr = append(requestsArr, req{elapsedFloat, (start.UnixNano()/int64(time.Millisecond) - startingTime)})

	// }
	progress++
	fmt.Printf("\r%.1f%% Completed.", float32(progress)/float32(rps*duration)*100)
	waitG.Done()
	return string(body)
}
