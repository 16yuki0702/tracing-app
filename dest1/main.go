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
	"github.com/yurishkuro/opentracing-tutorial/go/lib/http"
	"github.com/yurishkuro/opentracing-tutorial/go/lib/tracing"
)

const (
    serviceName = "dest1"
)

func handler(w http.ResponseWriter, r *http.Request) {
    tracer, closer := tracing.Init(serviceName)
	defer closer.Close()

    spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	span := tracer.StartSpan(serviceName, ext.RPCServerOption(spanCtx))
	defer span.Finish()

	resptStr := fmt.Sprintf("Hello from %s!!", serviceName)
	span.LogFields(
		otlog.String("event", serviceName),
		otlog.String("value", resptStr),
	)
	w.Write([]byte(resptStr))
}

func forwardDest2(w http.ResponseWriter, r *http.Request) {
    trace(w, r, serviceName, "http://dest2-svc:8080")
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

func trace(w http.ResponseWriter, r *http.Request, traceName, url string) {
    tracer, closer := tracing.Init(traceName)
	defer closer.Close()

    spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	span := tracer.StartSpan(traceName, ext.RPCServerOption(spanCtx))
	defer span.Finish()

    req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

    ext.SpanKindRPCClient.Set(span)
    ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	resp, err := xhttp.Do(req)
	if err != nil {
		panic(err.Error())
	}

	respStr := string(resp)

	span.LogFields(
		otlog.String("event", traceName),
		otlog.String("value", respStr),
	)

	w.Write([]byte(respStr))
}
