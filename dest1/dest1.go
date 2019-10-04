package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
    opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/yurishkuro/opentracing-tutorial/go/lib/tracing"
)

func handler(w http.ResponseWriter, r *http.Request) {
    tracer, closer := tracing.Init("dest1")
	defer closer.Close()

    spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	span := tracer.StartSpan("dest1", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	resptStr := fmt.Sprintf("Hello from dest1!!")
	span.LogFields(
		otlog.String("event", "dest1"),
		otlog.String("value", resptStr),
	)
	w.Write([]byte(resptStr))
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handler)

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

	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-interruptChan

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}
