package google

import (
	"net/http"
	"net/http/httptrace"

	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
)

// We will talk about this later
var tracer opentracing.Tracer

func AskGoogle(ctx context.Context) error {
	// retrieve current Span from Context
	var parentCtx opentracing.SpanContext
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan != nil {
		parentCtx = parentSpan.Context()
	}

	// start a new Span to wrap HTTP request
	span := tracer.StartSpan(
		"ask google",
		opentracing.ChildOf(parentCtx),
	)

	// make sure the Span is finished once we're done
	defer span.Finish()

	// make the Span current in the context
	ctx = opentracing.ContextWithSpan(ctx, span)

	// now prepare the request
	req, err := http.NewRequest("GET", "http://google.com", nil)
	if err != nil {
		return err
	}

	// attach ClientTrace to the Context, and Context to request
	trace := NewClientTrace(span)
	ctx = httptrace.WithClientTrace(ctx, trace)
	req = req.WithContext(ctx)

	// execute the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// Google home page is not too exciting, so ignore the result
	res.Body.Close()
	return nil
}
