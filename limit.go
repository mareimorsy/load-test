package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

func maxClients(h http.Handler, n int) http.Handler {
	sema := make(chan struct{}, n)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(sema) >= 10 {
			w.WriteHeader(http.StatusServiceUnavailable)
			io.WriteString(w, "Rate limit exceded")
		} else {
			sema <- struct{}{}
			defer func() { <-sema }()
			h.ServeHTTP(w, r)
		}
	})
}

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1000 * time.Millisecond)
		io.WriteString(w, "It works exactly as expected ya Awadi :)")
	})

	http.Handle("/sleep", maxClients(handler, 10))

	log.Fatal(http.ListenAndServe(":8081", nil))
}
