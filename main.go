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
	resptStr := fmt.Sprintf("Hello from %s!!!!!!", serviceName)
	w.Write([]byte(resptStr))
}

func generatePropagateFunc(envName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tracing.Propagate(w, r, os.Getenv(envName))
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
	_, closer := tracing.InitTracing(os.Getenv("SERVICE_NAME"))
	defer closer.Close()

	r := mux.NewRouter()

	r.HandleFunc("/", handler)
	for i := getServiceNum(); i > 0; i-- {
		r.HandleFunc(fmt.Sprintf("/propagate%d", i), generatePropagateFunc(fmt.Sprintf("PROPAGATE%d", i)))
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
