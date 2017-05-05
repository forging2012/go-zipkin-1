package main

import (
	"log"

	"github.com/raymasson/go-zipkin/client"
	"github.com/raymasson/go-zipkin/config"
	"github.com/raymasson/go-zipkin/server"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport/zipkin"
)

func init() {
	config.Get()
}

func main() {
	if *config.ActorKind != config.Server && *config.ActorKind != config.Client {
		log.Fatal("Please specify '-actor server' or '-actor client'")
	}

	// Jaeger tracer can be initialized with a transport that will
	// report tracing Spans to a Zipkin backend
	transport, err := zipkin.NewHTTPTransport(
		*config.ZipkinURL,
		zipkin.HTTPBatchSize(1),
		zipkin.HTTPLogger(jaeger.StdLogger),
	)
	if err != nil {
		log.Fatalf("Cannot initialize HTTP transport: %v", err)
	}
	// create Jaeger tracer
	tracer, closer := jaeger.NewTracer(
		*config.ActorKind,
		jaeger.NewConstSampler(true), // sample all traces
		jaeger.NewRemoteReporter(transport, nil),
	)

	if *config.ActorKind == config.Server {
		server.Run(tracer)
		return
	}

	client.Run(tracer)

	// Close the tracer to guarantee that all spans that could
	// be still buffered in memory are sent to the tracing backend
	closer.Close()
}
