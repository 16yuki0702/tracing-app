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
	serviceName = "gateway"
)

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	msg := query.Get("message")
	if msg == "" {
		msg = "Nothing"
	}
	log.Printf("Received message %s\n", msg)
	w.Write([]byte(fmt.Sprintf("You sent message %s\n", msg)))
}

func toDest1(w http.ResponseWriter, r *http.Request) {
    tracing.TraceWithForward(w, r, serviceName, "http://dest1-svc:8080")
}

func toDest2(w http.ResponseWriter, r *http.Request) {
	tracing.TraceWithForward(w, r, serviceName, "http://dest1-svc:8080/forwarddest2")
}

func toDest3(w http.ResponseWriter, r *http.Request) {
    tracing.TraceWithForward(w, r, serviceName, "http://dest3-svc:8080")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handler)
	r.HandleFunc("/todest1", toDest1)
	r.HandleFunc("/todest2", toDest2)
	r.HandleFunc("/todest3", toDest3)
	r.HandleFunc("/health", healthHandler)
	r.HandleFunc("/readiness", readinessHandler)

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
