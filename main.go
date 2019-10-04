package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/16yuki0702/tracing-app/tracing"
	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	serviceName := os.Getenv("SERVICE_NAME")
	resptStr := fmt.Sprintf("Hello from %s!!", serviceName)
	tracing.Trace(w, r, serviceName, resptStr)
	w.Write([]byte(resptStr))
}

func generateForwardFunc(envName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tracing.TraceWithForward(w, r, os.Getenv("SERVICE_NAME"), os.Getenv(envName))
	}
}

func getServiceNum() int {
	serviceNum := os.Getenv("SERVICE_NUM")
	v, err := strconv.Atoi(serviceNum)
	if err != nil {
		panic("Unexpected Error!!")
	}
	return v
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handler)
	for i := getServiceNum(); i > 0; i-- {
		r.HandleFunc(fmt.Sprintf("/todest%d", i), generateForwardFunc(fmt.Sprintf("DEST%d", i)))
		r.HandleFunc(fmt.Sprintf("/forward%d", i), generateForwardFunc(fmt.Sprintf("FORWARD%d", i)))
	}

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
