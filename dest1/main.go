package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/16yuki0702/tracing-app/tracing"
	"github.com/gorilla/mux"
)

const (
	serviceName = "dest1"
)

func handler(w http.ResponseWriter, r *http.Request) {
	resptStr := fmt.Sprintf("Hello from %s!!", serviceName)
	tracing.TraceWithForward(w, r, serviceName, resptStr)
	w.Write([]byte(resptStr))
}

func forwardDest2(w http.ResponseWriter, r *http.Request) {
	tracing.TraceWithForward(w, r, serviceName, "http://dest2-svc:8080")
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handler)
	r.HandleFunc("/forwarddest2", forwardDest2)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	tracing.WaitForShutdown(srv)
}
