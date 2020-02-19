package tracing

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	opentelemetry "go.opentelemetry.io/otel/api/trace"
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
	propagator := opentelemetry.B3{}
	ctx := propagator.Extract(context.Background(), r.Header)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	propagator.Inject(ctx, r.Header)

	incomingHeaders := []string{
		"x-request-id",
		"user-agent",
	}
	for _, header := range incomingHeaders {
		req.Header.Set(header, r.Header.Get(header))
	}

	resp, err := doRequest(req)
	if err != nil {
		panic(err.Error())
	}

	w.Write([]byte(string(resp)))
}

func doRequest(req *http.Request) ([]byte, error) {
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
