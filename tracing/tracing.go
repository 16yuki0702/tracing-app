package tracing

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/zipkin"
)

func WaitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-interruptChan

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}

func Propagate(w http.ResponseWriter, r *http.Request, url string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	incomingHeaders := []string{
		"x-request-id",
		"x-b3-traceid",
		"x-b3-spanid",
		"x-b3-parentspanid",
		"x-b3-sampled",
		"x-b3-flags",
		"x-datadog-trace-id",
		"x-datadog-parent-id",
		"x-datadog-sampled",
	}

	for _, header := range incomingHeaders {
		req.Header.Set(header, r.Header.Get(header))
	}
	req.Header.Set("user-agent", r.Header.Get("user-agent"))

	resp, err := do(req)
	if err != nil {
		panic(err.Error())
	}

	w.Write([]byte(string(resp)))
}

func Trace(w http.ResponseWriter, r *http.Request, traceName, url string) {
	tracer := opentracing.GlobalTracer()

	spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil {
		panic(err)
	}

	span := tracer.StartSpan(traceName, opentracing.ChildOf(spanCtx), ext.SpanKindRPCServer)
	defer span.Finish()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	tracer.Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	incomingHeaders := []string{
		"x-request-id",
		"x-datadog-trace-id",
		"x-datadog-parent-id",
		"x-datadog-sampled",
	}
	for _, header := range incomingHeaders {
		req.Header.Set(header, r.Header.Get(header))
	}
	req.Header.Set("user-agent", r.Header.Get("user-agent"))

	resp, err := do(req)
	if err != nil {
		panic(err.Error())
	}

	w.Write([]byte(string(resp)))
}

func InitTracing(serviceName string) (opentracing.Tracer, io.Closer) {
	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	injector := jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, zipkinPropagator)
	extractor := jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, zipkinPropagator)
	zipkinSharedRPCSpan := jaeger.TracerOptions.ZipkinSharedRPCSpan(true)

	tracer, closer := jaeger.NewTracer(
		serviceName,
		jaeger.NewConstSampler(true),
		jaeger.NewNullReporter(),
		injector,
		extractor,
		zipkinSharedRPCSpan,
	)

	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}

func do(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("StatusCode: %d, Body: %s", resp.StatusCode, body)
	}

	return body, nil
}
