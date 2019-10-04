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

    "github.com/16yuki0702/tracing-app/tracing"
	"github.com/gorilla/mux"
    opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/yurishkuro/opentracing-tutorial/go/lib/http"
	"github.com/yurishkuro/opentracing-tutorial/go/lib/tracing"
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
    util.Trace(w, r, serviceName, "http://dest1-svc:8080")
}

func toDest2(w http.ResponseWriter, r *http.Request) {
    util.Trace(w, r, serviceName, "http://dest1-svc:8080/forwarddest2")
}

func toDest3(w http.ResponseWriter, r *http.Request) {
    util.Ttrace(w, r, serviceName, "http://dest3-svc:8080")
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

	util.WaitForShutdown(srv)
}

//func waitForShutdown(srv *http.Server) {
//	interruptChan := make(chan os.Signal, 1)
//	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
//
//	<-interruptChan
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
//	defer cancel()
//	srv.Shutdown(ctx)
//
//	log.Println("Shutting down")
//	os.Exit(0)
//}
//
//func trace(w http.ResponseWriter, r *http.Request, traceName, url string) {
//    tracer, closer := tracing.Init(traceName)
//	defer closer.Close()
//
//    spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
//	span := tracer.StartSpan(traceName, ext.RPCServerOption(spanCtx))
//	defer span.Finish()
//
//    req, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		panic(err.Error())
//	}
//
//    ext.SpanKindRPCClient.Set(span)
//    ext.HTTPUrl.Set(span, url)
//	ext.HTTPMethod.Set(span, "GET")
//	span.Tracer().Inject(
//		span.Context(),
//		opentracing.HTTPHeaders,
//		opentracing.HTTPHeadersCarrier(req.Header),
//	)
//
//	resp, err := xhttp.Do(req)
//	if err != nil {
//		panic(err.Error())
//	}
//
//	respStr := string(resp)
//
//	span.LogFields(
//		otlog.String("event", traceName),
//		otlog.String("value", respStr),
//	)
//
//	w.Write([]byte(respStr))
//}
