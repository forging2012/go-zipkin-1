package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/raymasson/go-zipkin/config"
	"golang.org/x/net/context"
)

//Run the client
func Run(tracer opentracing.Tracer) {
	// nethttp.Transport from go-stdlib will do the tracing
	c := &http.Client{Transport: &nethttp.Transport{}}

	// create a top-level span to represent full work of the client
	span := tracer.StartSpan(config.Client)
	span.SetTag(string(ext.Component), config.Client)
	defer span.Finish()
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://localhost:%s/", *config.ServerPort),
		nil,
	)
	if err != nil {
		onError(span, err)
		return
	}

	req = req.WithContext(ctx)
	// wrap the request in nethttp.TraceRequest
	req, ht := nethttp.TraceRequest(tracer, req)
	defer ht.Finish()

	res, err := c.Do(req)
	if err != nil {
		onError(span, err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		onError(span, err)
		return
	}
	fmt.Printf("Received result: %s\n", string(body))
}

func onError(span opentracing.Span, err error) {
	// handle errors by recording them in the span
	span.SetTag(string(ext.Error), true)
	span.LogKV(otlog.Error(err))
	log.Print(err)
}
